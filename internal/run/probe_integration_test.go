//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// probeHarness gives a harness register row whose scripts are the test's own:
// version prints v, and the harness reads its prompt on stdin and writes
// probe.tsv with the prompt's token and AGENTS.md, or with the token the test
// gives (bad) when it is not "".
func probeHarness(t *testing.T, id, v, bad string) []string {
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\n"+v+"\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte(`#!/bin/sh
p=$(cat)
result=$(printf '%s\n' "$p" | sed -n 's/^Write one file, \(.*\), of tab-separated values, with the header line$/\1/p')
token=$(printf '%s\n' "$p" | sed -n 's/.*one row "token<TAB>\([0-9a-f]*\)".*/\1/p')
[ -n "`+bad+`" ] && token=`+bad+`
printf 'kind\tvalue\ntoken\t%s\nfile\tAGENTS.md\n' "$token" > "$result"
`), 0o755)
	return []string{id, "", "1", filepath.Join(bin, "harness"), "stdin", filepath.Join(bin, "version"), "", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""}
}

var probeModels = [][]string{{"fake", "m1", "1000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""}}

// The rows of harnesses.tsv of the records branch of the target.
func harnessRows(t *testing.T, w *world) [][]string {
	var out [][]string
	for i, l := range strings.Split(strings.TrimSpace(show(t, w.bare, "layup-records", "harnesses.tsv")), "\n") {
		if i > 0 {
			out = append(out, strings.Split(l, "\t"))
		}
	}
	return out
}

func TestTheStepProbe(t *testing.T) {
	w, r, f, a := sessionWorld(t)
	ctx := context.Background()
	harnesses := [][]string{probeHarness(t, "fake", "echo 1.0.0", ""), {"none", "", "1", "x", "stdin", "x", "", "", "R.md", "", "none", "api", ""}}
	routing := [][]string{{"developer", "execution", "1", "fake", "m1"}}
	passed, failed, skipped, err := r.ProbeStep(ctx, harnesses, probeModels, routing, a.base, r.Store.(*recordsStore).Base)
	if err != nil || passed != 1 || failed != 0 || skipped != 1 {
		t.Fatalf("the step: %d, %d, %d, %v; want 1 passed, 1 skipped", passed, failed, skipped, err)
	}
	rows := harnessRows(t, w)
	if len(rows) != 1 || rows[0][1] != "fake" || rows[0][2] != "1.0.0" || rows[0][3] != "m1" || rows[0][4] != "passed" || rows[0][5] != "—" || rows[0][6] != "AGENTS.md" {
		t.Fatalf("harnesses.tsv %q", rows)
	}
	id := rows[0][0]
	if s := show(t, w.bare, "layup-records", "sessions.tsv"); !strings.Contains(s, "\n"+id+"\tT-") || !strings.Contains(s, "\t1\tprobe\tfake\t1.0.0\tm1\t"+a.base+"\t") {
		t.Errorf("the start row of the probe:\n%s", s)
	}
	// The records column of the start row is the last pushed records commit at
	// the probe's start: the parent of the commit that adds the row, after the
	// copy of the routing register.
	added := strings.Fields(gitOut(t, w.bare, "log", "--reverse", "--format=%P", "-S", id+"\tT-", "layup-records", "--", "sessions.tsv"))
	if len(added) == 0 || !strings.Contains(show(t, w.bare, "layup-records", "sessions.tsv"), "\t"+a.base+"\t"+added[0]+"\t") {
		t.Errorf("the records column of the probe's start row: want %q", added)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\t") || !strings.Contains(tel, "\tprobe\tfake\tm1\tsubscription\t") {
		t.Errorf("the telemetry row of the probe:\n%s", tel)
	}
	if got := show(t, w.bare, "layup-records", "routing.tsv"); got != "role\ttier\tposition\tharness\tmodel\ndeveloper\texecution\t1\tfake\tm1\n" {
		t.Errorf("routing.tsv %q", got)
	}
	// A probe has no events: its task has no events.tsv.
	if names := gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records", "tasks"); strings.Contains(names, "tasks/T-") && !strings.Contains(names, "tasks/T-ab12/") || strings.Count(names, "events.tsv") != 1 {
		t.Errorf("the probe has events:\n%s", names)
	}
	if c := f.posted[2]; len(c) != 1 || c[0] != id+": probe of fake 1.0.0: passed\n" {
		t.Errorf("the comments %q", c)
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the probe's directory after it: %v", err)
	}
	// Again at the same version: no probe, no new record, no copy of an equal
	// routing register.
	commits := gitOut(t, w.bare, "rev-list", "--count", "layup-records")
	passed, failed, skipped, err = r.ProbeStep(ctx, harnesses, probeModels, routing, a.base, r.Store.(*recordsStore).Base)
	if err != nil || passed != 0 || failed != 0 || skipped != 1 || gitOut(t, w.bare, "rev-list", "--count", "layup-records") != commits {
		t.Errorf("the step again: %d, %d, %d, %v; want none probed and no commit", passed, failed, skipped, err)
	}
	// The probe that ended at its version check left no directory.
	if left, err := os.ReadDir(filepath.Join(r.Host, "sessions")); err != nil || len(left) != 0 {
		t.Errorf("the session directories after the step again: %v, %v; want none", left, err)
	}
}

func TestAProbeThatFailsAndARefusedProbe(t *testing.T) {
	w, r, f, a := sessionWorld(t)
	ctx := context.Background()
	harnesses := [][]string{probeHarness(t, "fake", "echo 1.0.0", "ffffffffffffffff"), probeHarness(t, "gone", "exit 1", "")}
	models := append(append([][]string{}, probeModels...), []string{"gone", "m2", "1000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""})
	passed, failed, skipped, err := r.ProbeStep(ctx, harnesses, models, nil, a.base, r.Store.(*recordsStore).Base)
	if err != nil || passed != 0 || failed != 2 || skipped != 0 {
		t.Fatalf("the step: %d, %d, %d, %v; want 2 failed", passed, failed, skipped, err)
	}
	rows := harnessRows(t, w)
	if len(rows) != 2 || rows[0][4] != "failed" || rows[0][5] != "token" || rows[1][1] != "gone" || rows[1][2] != "—" || rows[1][5] != "version" {
		t.Errorf("harnesses.tsv %q; want token, then a refused start (version, its version —)", rows)
	}
	// One comment: the probe that ran; a refused start posts none.
	if c := f.posted[2]; len(c) != 1 || !strings.HasSuffix(c[0], ": probe of fake 1.0.0: failed token\n") {
		t.Errorf("the comments %q", c)
	}
}

// A task session on a version with no probe runs the probe first; one whose
// probe fails is refused (probe).
func TestATaskSessionProbesItsVersionFirst(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	ctx := context.Background()
	for _, c := range []struct {
		bad, want string
	}{{"", ""}, {"ffffffffffffffff", "probe"}} {
		spec := harness(t, a, work(helloSum))
		h := probeHarness(t, "fake", "echo "+map[string]string{"": "2.0.0", "ffffffffffffffff": "3.0.0"}[c.bad], c.bad)
		spec.Pair.Harness, spec.Pair.VersionCommand = "fake", h[regVersion]
		spec.Admit = r.Admit(ctx, ProbeSpec{Harness: h, Models: probeModels, Base: a.base, Records: a.records})
		id, err := r.TaskSession(ctx, spec)
		switch c.want {
		case "":
			if err != nil || len(refusals(t, w, id)) != 0 {
				t.Errorf("a version with no probe that passes: %v, %q", err, refusals(t, w, id))
			}
		default:
			if got := refusals(t, w, id); len(got) != 1 || got[0] != "probe" {
				t.Errorf("a probe that fails: %v, the refusals %q; want probe", err, got)
			}
		}
	}
	rows := harnessRows(t, w)
	if len(rows) != 2 || rows[0][2] != "2.0.0" || rows[0][4] != "passed" || rows[1][2] != "3.0.0" || rows[1][4] != "failed" {
		t.Errorf("harnesses.tsv %q; want the probe of each version", rows)
	}
}
