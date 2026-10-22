package run

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"slices"
	"strings"
	"testing"
)

// The demo of task T-0xr7 (#493): the records of the architecture task T-am8b of pharzam/layup-uat-m2e2, where the verifier found the bet's decisions
// material in two rounds (the comments 6347751487 and 6347893425) and the task stood at the review-cycle cap of 1 with no way out.
//
// Provenance: testdata/t-am8b.txtar holds, each under a line `-- <path> --`, the files of the branch layup-records of
// https://github.com/pharzam/layup-uat-m2e2.git (public) at the commit e65db3e30027cee5e74f466c9fc88897fd60a8b3, byte for byte: approvers.tsv and
// tasks.tsv; the 28 events of tasks/T-am8b/events.tsv; the rows of sessions.tsv of the task; its seven results and the three kept comments; and the rows
// of copies.tsv of the plan review and the two review records, with their copies. Nothing of the target is repaired, and the findings stay
// open as the verifier wrote them. The test adds one thing: the Operator's comment, as the boundary would copy it.
//
//go:embed testdata/t-am8b.txtar
var am8bArchive string

const (
	am8bBase      = "d03c8650303fbdff49986ef5478ed53d11219964"
	am8bGrantID   = "6347900000"
	am8bGrantBody = "cycle T-am8b because the findings of rounds 1 and 2 are real\n"
)

func am8bFiles(t *testing.T) map[string][]byte {
	t.Helper()
	files, name := map[string][]byte{}, ""
	for _, l := range strings.SplitAfter(am8bArchive, "\n") {
		if p, ok := strings.CutSuffix(strings.TrimPrefix(strings.TrimSuffix(l, "\n"), "-- "), " --"); ok && strings.HasPrefix(l, "-- ") {
			name = p
			files[name] = nil
		} else if name != "" {
			files[name] = append(files[name], l...)
		}
	}
	if len(files) != 18 {
		t.Fatalf("the archive holds %d files; want 18", len(files))
	}
	return files
}

func TestTheDemoRecordsOfTAm8bWithAGrantStartAttempt4AtTheRecordedBaseWithTheOpenFindings(t *testing.T) {
	ctx := context.Background()
	files := am8bFiles(t)
	original := map[string][]byte{}
	for p, d := range files {
		original[p] = slices.Clone(d)
	}
	store := &standinStore{files: files, holder: "0123456789abcdef"}
	s := newSessions(store)
	s.Control = 2

	// The records as the target left them: 28 events, two rounds with findings, the wait at the cap with the command named.
	v, err := s.shapeView(ctx, "T-am8b")
	if err != nil {
		t.Fatal(err)
	}
	if n := len(v.Events); n != 28 || reviewRounds(v) != 2 {
		t.Fatalf("%d events and %d rounds with findings; want the 28 events and 2 rounds of the target", n, reviewRounds(v))
	}
	want := "waiting for the Operator: T-am8b: the review-cycle cap 1 is reached with findings open; the Operator may grant one more cycle with `cycle T-am8b because <reason>` on the control issue"
	if got := loopPosition(v); got.Pos != posWait || got.Msg != want {
		t.Fatalf("the position of the records: %+v; want the wait %q", got, want)
	}

	// The Operator's comment on the control issue, after the comment of the cap round (6347893425), is copied by the boundary: one row, one body.
	sum := sha256.Sum256([]byte(am8bGrantBody))
	row := strings.Join([]string{am8bGrantID, "1", "2", "1675602", "pharzam", "—", "2026-10-22T18:30:00Z", "2026-10-22T18:30:30Z", hex.EncodeToString(sum[:]), "copies/" + am8bGrantID + "-1.md"}, "\t") + "\n"
	files["copies.tsv"] = append(files["copies.tsv"], row...)
	files["copies/"+am8bGrantID+"-1.md"] = []byte(am8bGrantBody)
	if v, err = s.shapeView(ctx, "T-am8b"); err != nil {
		t.Fatal(err)
	}
	got := loopPosition(v)
	if got.Pos != posNextAttempt || got.Attempt != 3 || got.Base != am8bBase || got.Session != "S-0465b007" {
		t.Fatalf("the position with the grant: %+v; want the next attempt after the verdict of S-0465b007, at the base %s", got, am8bBase)
	}

	// The driver's act is StartAttempt at that base: attempt 4, and the old records are an unchanged prefix.
	if n, err := s.StartAttempt(ctx, "T-am8b", got.Base); err != nil || n != 4 {
		t.Fatalf("StartAttempt = %d, %v; want attempt 4", n, err)
	}
	events := string(files["tasks/T-am8b/events.tsv"])
	lines := strings.Split(strings.TrimSuffix(events, "\n"), "\n")
	if f := strings.Split(lines[len(lines)-1], "\t"); len(lines) != 30 || f[0] != "29" || f[1] != "attempt" || f[2] != "4" || f[4] != am8bBase {
		t.Errorf("the new event %q (%d lines); want the event 29, attempt 4, base %s", lines[len(lines)-1], len(lines), am8bBase)
	}
	for p, was := range original { // the two tables that grew keep their old rows as a prefix, every other record is unchanged
		if grown := p == "tasks/T-am8b/events.tsv" || p == "copies.tsv"; !bytes.Equal(files[p], was) && (!grown || !bytes.HasPrefix(files[p], was)) {
			t.Errorf("%s is not an unchanged prefix of the old records", p)
		}
	}

	// Attempt 4 is a developer's, at the base, and the developer is given the findings of the attempt-3 record, both of them, still open.
	if v, err = s.shapeView(ctx, "T-am8b"); err != nil {
		t.Fatal(err)
	}
	if now := loopPosition(v); now.Pos != posDeveloper || now.Attempt != 4 || now.Base != am8bBase {
		t.Errorf("the position after the attempt: %+v; want the developer of attempt 4 at the base %s", now, am8bBase)
	}
	rec, ok := ParseReviewRecord(original["tasks/T-am8b/results/S-0465b007/review-record.md"])
	findings, gate := findingsOf(v, 4)
	if !ok || len(rec.Findings) != 2 || !slices.Equal(findings, rec.Findings) || gate != "" {
		t.Errorf("the findings of attempt 4: %+v, gate head %q; want the two findings of the record of attempt 3: %+v", findings, gate, rec.Findings)
	}
}
