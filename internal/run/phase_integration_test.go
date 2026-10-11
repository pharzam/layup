//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/route"
)

// headsForge is the stand-in forge that keeps, for each comment, the head of
// layup-records at the time it was posted, so a test sees the order of the
// records commit and the comment.
type headsForge struct {
	*standInForge
	t     *testing.T
	bare  string
	heads []string
}

func (f *headsForge) Comment(ctx context.Context, issue int, body string) (int64, error) {
	f.heads = append(f.heads, strings.TrimSpace(gitOut(f.t, f.bare, "rev-parse", "layup-records")))
	return f.standInForge.Comment(ctx, issue, body)
}

// phaseHarness gives a harness register row whose program answers a probe
// as probeHarness does, and runs body for a session of a phase, with $in the
// input file and $out the output file that the test's prompt names.
func phaseHarness(t *testing.T, body string) []string {
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte(`#!/bin/sh
p=$(cat)
case "$p" in
"This is a probe"*)
	result=$(printf '%s\n' "$p" | sed -n 's/^Write one file, \(.*\), of tab-separated values, with the header line$/\1/p')
	token=$(printf '%s\n' "$p" | sed -n 's/.*one row "token<TAB>\([0-9a-f]*\)".*/\1/p')
	printf 'kind\tvalue\ntoken\t%s\nfile\tAGENTS.md\n' "$token" > "$result" ;;
*)
	in=$(printf '%s\n' "$p" | sed -n 's/^input //p')
	out=$(printf '%s\n' "$p" | sed -n 's/^output //p')
	`+body+` ;;
esac
`), 0o755)
	return []string{"fake", "", "1", filepath.Join(bin, "harness"), "stdin", filepath.Join(bin, "version"), "", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""}
}

// phaseWorld gives the world of the sessions with the routing register of
// one pair for systems-architect, probed and passed by the step probe, and a
// forge that keeps the head of the records at each comment.
func phaseWorld(t *testing.T, body string) (*world, *Sessions, *headsForge, worldAt, []string) {
	w, r, f, a := sessionWorld(t)
	hf := &headsForge{standInForge: f, t: t, bare: w.bare}
	r.Forge = hf
	h := phaseHarness(t, body)
	routing := [][]string{{"systems-architect", "reasoning", "1", "fake", "m1"}}
	if p, _, _, err := r.ProbeStep(context.Background(), [][]string{h}, probeModels, routing, a.base, r.Store.(*recordsStore).Base); err != nil || p != 1 {
		t.Fatalf("the step probe: %d passed, %v", p, err)
	}
	hf.heads, f.posted = nil, nil
	return w, r, hf, worldAt{a.base, r.Store.(*recordsStore).Base()}, h
}

// phaseSpec gives the session of the step meaning on the pair that PhasePair
// gives from the harness row, with the problem statement of the records as its input; Keep refuses
// a file that holds "bad", and gives intake/phase.tsv, the session ID and the
// file.
func phaseSpec(t *testing.T, w *world, r *Sessions, a worldAt, harness []string) PhaseSpec {
	t.Helper()
	h, m, err := r.PhasePair(context.Background(), "systems-architect", [][]string{harness}, probeModels)
	if err != nil {
		t.Fatalf("PhasePair: %v", err)
	}
	return PhaseSpec{Step: "meaning", Role: "systems-architect", Harness: h, Model: m, Models: probeModels, Base: a.base, Records: a.records,
		Input:  map[string][]byte{"problem-statement.md": []byte(show(t, w.bare, "layup-records", "start/problem-statement.md"))},
		Output: "meaning.tsv",
		Prompt: func(in map[string]string, out string) []byte {
			return []byte("input " + in["problem-statement.md"] + "\noutput " + out + "\n")
		},
		Keep: func(id string, data []byte) (map[string][]byte, error) {
			if bytes.Contains(data, []byte("bad")) {
				return nil, errors.New("a bad row")
			}
			return map[string][]byte{"intake/phase.tsv": append([]byte(id+"\t"), data...)}, nil
		}}
}

// tree gives the paths of layup-records.
func tree(t *testing.T, w *world) []string {
	return strings.Fields(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records"))
}

// The demo of row 42a: a session of a phase on a scripted harness reads its
// input from input/ and commits its records and its telemetry row in one
// records commit; it has its own task ID, a start row and no event; then its
// comment; its directory is removed.
func TestASessionOfAPhaseCommitsItsRecordsAndItsTelemetryRow(t *testing.T) {
	seen := filepath.Join(t.TempDir(), "seen")
	w, r, f, a, h := phaseWorld(t, `cp "$in" `+seen+` && printf 'kind\tvalue\nnote\tread\n' > "$out"`)
	before := show(t, w.bare, "layup-records", "sessions.tsv")
	spec := phaseSpec(t, w, r, a, h)
	id, err := r.PhaseSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("PhaseSession: %q, %v", id, err)
	}
	if got, err := os.ReadFile(seen); err != nil || !bytes.Equal(got, spec.Input["problem-statement.md"]) {
		t.Errorf("the input the harness read: %q, %v; want the record byte for byte", got, err)
	}
	if got := show(t, w.bare, "layup-records", "intake/phase.tsv"); got != id+"\tkind\tvalue\nnote\tread\n" {
		t.Errorf("the record kept: %q", got)
	}
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "layup-records", "--", "intake/phase.tsv"))
	if changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", commit)); !slices.Equal(changed, []string{"intake/phase.tsv", "telemetry.tsv"}) {
		t.Errorf("the commit of the records changed %q; want the records and the telemetry row in one", changed)
	}
	// Its own task ID: T- and four characters, held by no row before it.
	sessions := show(t, w.bare, "layup-records", "sessions.tsv")
	var task string
	for _, l := range strings.Split(sessions, "\n") {
		if c := strings.Split(l, "\t"); c[0] == id {
			task = c[1]
			if c[2] != "1" || c[3] != "systems-architect" || c[4] != "fake" || c[6] != "m1" || c[7] != a.base {
				t.Errorf("the start row %q", c)
			}
		}
	}
	if !regexp.MustCompile(`^T-[0-9a-z]{4}$`).MatchString(task) || strings.Contains(before, "\t"+task+"\t") {
		t.Errorf("the task ID %q: want T- and four characters, held by no earlier row of sessions.tsv:\n%s", task, before)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\t"+task+"\t—\tsystems-architect\tfake\tm1\tsubscription\t") {
		t.Errorf("telemetry.tsv has no row of %s:\n%s", id, tel)
	}
	// No event: no events.tsv but the one of T-ab12.
	for _, p := range tree(t, w) {
		if strings.HasSuffix(p, "events.tsv") && p != "tasks/T-ab12/events.tsv" {
			t.Errorf("a session of a phase wrote %s", p)
		}
	}
	if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of meaning: done\n" {
		t.Errorf("the comments %q", f.posted[2])
	}
	if len(f.heads) != 1 || f.heads[0] != commit {
		t.Errorf("the head of the records at the comment %q; want the records commit %s", f.heads, commit)
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

// Row 3 of Input states: an output file that is missing, a link, over 1 MiB,
// or that its check refuses: the step's error names the session; the end
// commits the telemetry row alone, and nothing of the file.
func TestAnOutputFileThatIsRefusedIsNoRecord(t *testing.T) {
	for _, c := range []struct{ name, body, comment, reason string }{
		{"missing", `true`, "no-result", "no-result"},
		{"a link", `ln -s /etc/hostname "$out"`, "no-result", "no-result"},
		{"over 1 MiB", `head -c 1048577 /dev/zero | tr '\0' x > "$out"`, "no-result", "no-result"},
		{"refused by its check", `printf 'kind\tvalue\nnote\tbad\n' > "$out"`, "done, refused file", "refused file"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, f, a, h := phaseWorld(t, c.body)
			id, err := r.PhaseSession(context.Background(), phaseSpec(t, w, r, a, h))
			var pe *PhaseError
			if !errors.As(err, &pe) || pe.Session != id || id == "" || !strings.HasPrefix(pe.Reason, c.reason) {
				t.Fatalf("PhaseSession: %q, %v; want a PhaseError of the session, %s", id, err, c.reason)
			}
			if slices.Contains(tree(t, w), "intake/phase.tsv") {
				t.Errorf("the refused file is a record")
			}
			last := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", "layup-records"))
			if !slices.Equal(last, []string{"telemetry.tsv"}) || !strings.Contains(show(t, w.bare, "layup-records", "telemetry.tsv"), "\n"+id+"\t") {
				t.Errorf("the last records commit changed %q; want the telemetry row of %s alone", last, id)
			}
			if len(f.posted[2]) != 1 || f.posted[2][0] != id+": systems-architect session of meaning: "+c.comment+"\n" {
				t.Errorf("the comments %q", f.posted[2])
			}
		})
	}
}

// Row 1 of Input states: with no admitted pair for the role, no session
// starts: ErrNoPair, and no records commit.
func TestNoAdmittedPairStartsNoSession(t *testing.T) {
	w, r, _, _, h := phaseWorld(t, `true`)
	commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
	for _, c := range []struct {
		name   string
		role   string
		models [][]string
	}{
		{"a role with no row", "product-owner", probeModels},
		{"a model of use no", "systems-architect", [][]string{{"fake", "m1", "1000", "https://example.invalid", "2026-10-09T12:00:00Z", "no", ""}}},
	} {
		if _, _, err := r.PhasePair(context.Background(), c.role, [][]string{h}, c.models); !errors.Is(err, route.ErrNoPair) {
			t.Errorf("%s: %v; want ErrNoPair", c.name, err)
		}
	}
	if got := gitOut(t, w.bare, "rev-list", "--count", "layup-records"); got != commits {
		t.Errorf("the records commits: %s, then %s; want none", commits, got)
	}
	// A harness with no row of harnesses.tsv is not admitted.
	_, r2, _, _ := sessionWorld(t)
	if err := r2.commit(context.Background(), map[string][]byte{"routing.tsv": []byte("role\ttier\tposition\tharness\tmodel\nsystems-architect\treasoning\t1\tfake\tm1\n")}, "the routing"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := r2.PhasePair(context.Background(), "systems-architect", [][]string{h}, probeModels); !errors.Is(err, route.ErrNoPair) {
		t.Errorf("a harness never probed: %v; want ErrNoPair", err)
	}
}

// A refused start of a session of a phase (here its version check) writes
// no record and removes its directory; the step's error names the session
// and the refusal.
func TestARefusedStartOfASessionOfAPhaseIsNoRecord(t *testing.T) {
	w, r, f, a, h := phaseWorld(t, `true`)
	spec := phaseSpec(t, w, r, a, h)
	spec.Harness = append([]string{}, h...)
	spec.Harness[regVersion] = "false"
	commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
	id, err := r.PhaseSession(context.Background(), spec)
	var pe *PhaseError
	if !errors.As(err, &pe) || pe.Session != id || id == "" || !strings.HasPrefix(pe.Reason, "refused version") {
		t.Fatalf("PhaseSession: %q, %v; want a PhaseError of the session, refused version", id, err)
	}
	if got := gitOut(t, w.bare, "rev-list", "--count", "layup-records"); got != commits {
		t.Errorf("the records commits: %s, then %s; want none", commits, got)
	}
	if len(f.posted[2]) != 0 {
		t.Errorf("the comments %q; want none", f.posted[2])
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

// A comment that fails after a refused file keeps the session's reason: the
// step's error names the session and the forge (round 1 of #215, note 5).
func TestACommentThatFailsKeepsTheReasonOfTheSession(t *testing.T) {
	w, r, f, a, h := phaseWorld(t, `printf 'kind\tvalue\nnote\tbad\n' > "$out"`)
	f.fail = "Comment"
	id, err := r.PhaseSession(context.Background(), phaseSpec(t, w, r, a, h))
	var pe *PhaseError
	if !errors.As(err, &pe) || pe.Session != id || !strings.HasPrefix(pe.Reason, "refused file") || !strings.Contains(err.Error(), "Comment") {
		t.Errorf("PhaseSession: %q, %v; want the PhaseError of the session, refused file, with the error of the comment", id, err)
	}
}
