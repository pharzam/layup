//go:build integration

package run

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/verify"
)

// The step scaffold-push on the real step runner (docs/spec/scaffold.md, The commands for the
// Operator and its paragraph "Decided here" of task T-hfn1; the four input states of the step;
// row 56 of the plan, #279).

// The outputs of the Operator's applies are read by their names; a name that is a directory, or
// that cannot be read, is an error that names it; a directory that does not exist has none.
func TestTheAppliesOfTheOperatorAreReadByTheirNames(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "applies", "acme", "target")
	if got, err := readApplies(dir); err != nil || len(got) != 0 {
		t.Errorf("a directory that does not exist: %q, %v; want no output", got, err)
	}
	writeFile(t, filepath.Join(dir, "default.json"), "d")
	writeFile(t, filepath.Join(dir, "other.json"), "x")
	if got, err := readApplies(dir); err != nil || len(got) != 1 || string(got["default.json"]) != "d" {
		t.Errorf("default.json alone: %q, %v; want it alone, as other names are not the Operator's", got, err)
	}
	writeFile(t, filepath.Join(dir, "records.json"), "r")
	if got, err := readApplies(dir); err != nil || len(got) != 2 || string(got["records.json"]) != "r" {
		t.Errorf("both: %q, %v; want both", got, err)
	}
	if err := os.Remove(filepath.Join(dir, "records.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "records.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, err := readApplies(dir); err == nil || !strings.Contains(err.Error(), "records.json") {
		t.Errorf("a directory named records.json: %q, %v; want an error that names it", got, err)
	}
}

// notes gives the lines of progress that the world has seen so far.
func notes(w *world) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return append([]string(nil), w.lines...)
}

// The demo of row 56: the step prints the whole rule of the Operator and each line of commands.sh,
// waits (the last row of the restart, which releases the lease) while the head of the forge is not the
// last commit of setup/tree.tsv or an output is missing, builds a missing work area again to the same
// commits in either case, fails for an output that is a directory, not an object of the repository, or
// of another repository, and, once the Operator has pushed and applied, commits setup/rulesets.txt in
// one commit of that path alone; a later restart finds it done before.
func TestTheStepScaffoldPushWaitsAndRecordsTheAppliesOfTheOperator(t *testing.T) {
	saved := verifyRun
	t.Cleanup(func() { verifyRun = saved })
	w := scaffoldWorld(t, proseFiles())
	verifyRun = func(string, func(int, int, string) func(), io.Writer) (verify.Table, error) {
		return verify.Table{Rows: []verify.Row{{Check: "pin", Result: "pass", Reason: "—"}}}, nil
	}
	cfg := w.config("0.1.0-dev")
	waiting := Step{"scaffold-push", "done", scaffoldWaiting}
	rows := toSetup(t, w, "go")
	if r := rows[len(rows)-1]; r != waiting || rowOf(rows, "verify") != (Step{"verify", "done", verifyDoneDetail}) {
		t.Fatalf("the last row of the first restart: %+v; want scaffold-push, waiting, after the step verify\n%s", r, steps(rows))
	}
	area := filepath.Join(w.dir, "work", "acme", "target")
	target := filepath.Join(area, "target")
	head := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup"))
	commands, err := os.ReadFile(filepath.Join(area, "out", "commands.sh"))
	if err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, n := range notes(w) {
		seen[n] = true
	}
	for _, line := range strings.Split(strings.TrimRight(string(commands), "\n"), "\n") {
		if !seen[line] {
			t.Errorf("commands.sh has the line %q, which was not printed", line)
		}
	}
	if !seen["scaffold-push: run "+filepath.Join(area, "out", "commands.sh")+" with your own login, from a clone with no hooks installed"] {
		t.Errorf("no note names the file and says whose login and what clone:\n%s", strings.Join(notes(w), "\n"))
	}
	if lease := recordsFile(t, w, "lease.tsv"); !strings.Contains(lease, "\treleased") {
		t.Errorf("the lease after a restart that waits: %q; want it released", lease)
	}
	if out, err := gitRunOK(w.bare, "show", "refs/heads/layup-records:setup/rulesets.txt"); err == nil {
		t.Errorf("setup/rulesets.txt is in the records while the step waits: %q", out)
	}
	sessions := setupSessions(w)
	applies := filepath.Join(w.dir, "applies", "acme", "target")
	defaultOut, recordsOut := "{\"source\":\"acme/target\",\"bypass_actors\":[]}\n", "{\"source\":\"acme/target\",\"bypass_actors\":[{\"actor_id\":42}]}"
	restart := func(want Step) []Step {
		t.Helper()
		rows := Restart(context.Background(), cfg)
		if r := rowOf(rows, "scaffold-push"); r.Name != want.Name || r.Result != want.Result || !strings.Contains(r.Detail, want.Detail) {
			t.Fatalf("the step scaffold-push: %+v; want %+v\n%s", r, want, steps(rows))
		}
		return rows
	}
	// Both outputs, but the forge's default branch is at the root commit: waiting; a work area that is there is not built again.
	writeFile(t, filepath.Join(applies, "default.json"), defaultOut)
	writeFile(t, filepath.Join(applies, "records.json"), recordsOut)
	writeFile(t, filepath.Join(area, "stays.txt"), "x")
	restart(waiting)
	if _, err := os.Stat(filepath.Join(area, "stays.txt")); err != nil {
		t.Errorf("a work area that was there was built again: %v", err)
	}
	// The Operator's push; an output missing; the work area gone: built again to the same commits, and the step waits.
	gitOut(t, target, "push", "-q", w.bare, "layup-setup:main")
	if err := os.Remove(filepath.Join(applies, "records.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
		t.Fatal(err)
	}
	restart(waiting)
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup")); got != head || setupSessions(w) != sessions {
		t.Errorf("the work area built again: layup-setup %s (was %s), %d sessions (were %d); want the same head and no session", got, head, setupSessions(w), sessions)
	}
	// An output that is a directory, then one of another repository: fail, naming the file.
	if err := os.Mkdir(filepath.Join(applies, "records.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	restart(Step{"scaffold-push", "fail", "records.json"})
	if err := os.Remove(filepath.Join(applies, "records.json")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(applies, "records.json"), `{"source":"acme/other"}`)
	restart(Step{"scaffold-push", "fail", "records.json"})
	// A work area built again whose loop fails (a text of the records whose citation resolves to no fact), and one whose
	// head is not the last commit of setup/tree.tsv: fail, the second naming both; the records are put right after each.
	writeFile(t, filepath.Join(applies, "records.json"), recordsOut)
	glossary := "setup/inputs/files/docs/glossary.md"
	pushRecords(t, w.bare, map[string]string{glossary: "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository (`F-0099#1`) |\n"}, "a glossary whose citation resolves to no fact")
	if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
		t.Fatal(err)
	}
	restart(Step{"scaffold-push", "fail", "glossary"})
	pushRecords(t, w.bare, map[string]string{glossary: proseFiles()["docs/glossary.md"]}, "the glossary again")
	tree := recordsFile(t, w, "setup/tree.tsv")
	lines := strings.Split(strings.TrimRight(tree, "\n"), "\n")
	cells := strings.Split(lines[len(lines)-1], "\t")
	cells[1] = strings.Repeat("1", 40)
	lines[len(lines)-1] = strings.Join(cells, "\t")
	pushRecords(t, w.bare, map[string]string{"setup/tree.tsv": strings.Join(lines, "\n") + "\n"}, "a tree whose last commit is not the head")
	if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
		t.Fatal(err)
	}
	for _, r := range restart(Step{"scaffold-push", "fail", strings.Repeat("1", 40)}) {
		if r.Name == "scaffold-push" && !strings.Contains(r.Detail, head) {
			t.Errorf("the detail of a head that is not the last commit: %q; want the head %s and the last commit", r.Detail, head)
		}
	}
	pushRecords(t, w.bare, map[string]string{"setup/tree.tsv": tree}, "the tree again")
	// Ready, with the work area gone: built again before the outcome; one commit of the one path.
	if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
		t.Fatal(err)
	}
	restart(Step{"scaffold-push", "done", scaffoldDone})
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup")); got != head || setupSessions(w) != sessions {
		t.Errorf("the work area built again on the ready path: layup-setup %s (was %s), %d sessions (were %d); want the same head and no session", got, head, setupSessions(w), sessions)
	}
	want := "== default.json\n" + defaultOut + "\n== records.json\n" + recordsOut + "\n"
	if got := recordsFile(t, w, "setup/rulesets.txt"); got != want {
		t.Errorf("setup/rulesets.txt:\n%q\nwant\n%q", got, want)
	}
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "--grep=^layup run: the apply of the rulesets$", "refs/heads/layup-records"))
	if changed := strings.Fields(gitOut(t, w.bare, "diff", "--name-only", commit+"~1", commit)); commit == "" || strings.Join(changed, " ") != "setup/rulesets.txt" {
		t.Errorf("the records commit of the applies %q changes %q; want setup/rulesets.txt alone", commit, changed)
	}
	// A later restart: done before, and nothing printed again.
	before := len(notes(w))
	restart(Step{"scaffold-push", "done", "done before"})
	for _, n := range notes(w)[before:] {
		if strings.HasPrefix(n, "scaffold-push:") {
			t.Errorf("a step that is done before printed %q", n)
		}
	}
}
