package run

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"slices"
	"strings"
	"testing"
)

// The demo of task T-0xr7 (#493): the records of the architecture task T-am8b of the target pharzam/layup-uat-m2e2, where the verifier found the bet's
// decisions material in two rounds (the comments 6347751487 and 6347893425) and the task stood at the review-cycle cap of 1, with no way out.
//
// Provenance: the files below testdata/t-am8b are the files of the branch layup-records of https://github.com/pharzam/layup-uat-m2e2.git (public) at the
// commit e65db3e30027cee5e74f466c9fc88897fd60a8b3, taken as they stand, unedited: approvers.tsv and tasks.tsv; the 28 events of
// tasks/T-am8b/events.tsv; the rows of sessions.tsv of the task; the seven results (the three kept comments among them); and the rows of copies.tsv
// of the two review records and of the plan and the plan review, with the copies they name. Nothing of the target is repaired by this test, and the
// findings stay open as the verifier wrote them. The test adds one thing, the Operator's comment, as the boundary would copy it.
//
//go:embed all:testdata/t-am8b
var am8bRecords embed.FS

const (
	am8bRoot      = "testdata/t-am8b"
	am8bBase      = "d03c8650303fbdff49986ef5478ed53d11219964"
	am8bOperator  = "1675602"
	am8bGrantID   = "6347900000"
	am8bGrantBody = "cycle T-am8b because the findings of rounds 1 and 2 are real\n"
)

func am8bFiles(t *testing.T) map[string][]byte {
	t.Helper()
	files := map[string][]byte{}
	err := fs.WalkDir(am8bRecords, am8bRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := am8bRecords.ReadFile(p)
		files[strings.TrimPrefix(p, am8bRoot+"/")] = data
		return err
	})
	if err != nil {
		t.Fatal(err)
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
	row := strings.Join([]string{am8bGrantID, "1", "2", am8bOperator, "pharzam", "—", "2026-10-22T18:30:00Z", "2026-10-22T18:30:30Z", hex.EncodeToString(sum[:]), "copies/" + am8bGrantID + "-1.md"}, "\t") + "\n"
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
	events := files["tasks/T-am8b/events.tsv"]
	if !bytes.HasPrefix(events, original["tasks/T-am8b/events.tsv"]) {
		t.Errorf("the old events are not an unchanged prefix of the new events.tsv")
	}
	lines := strings.Split(strings.TrimSuffix(string(events), "\n"), "\n")
	if f := strings.Split(lines[len(lines)-1], "\t"); len(lines) != 30 || f[0] != "29" || f[1] != "attempt" || f[2] != "4" || f[4] != am8bBase {
		t.Errorf("the new event %q (%d lines); want the event 29, attempt 4, base %s", lines[len(lines)-1], len(lines), am8bBase)
	}
	for p, was := range original {
		if p != "tasks/T-am8b/events.tsv" && p != "copies.tsv" && !bytes.Equal(files[p], was) {
			t.Errorf("%s changed; want every old record unchanged", p)
		}
	}
	if !bytes.HasPrefix(files["copies.tsv"], original["copies.tsv"]) {
		t.Errorf("copies.tsv: the old rows are not an unchanged prefix")
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
	if len(findings) == 2 && (findings[0].File != "docs/tasks/T-am8b.md" || findings[0].Line != 49 || findings[1].File != "docs/adr/README.md" || findings[1].Line != 20) {
		t.Errorf("the findings %+v; want the two the verifier wrote at docs/tasks/T-am8b.md:49 and docs/adr/README.md:20", findings)
	}
}
