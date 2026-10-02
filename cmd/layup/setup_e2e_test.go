//go:build e2e

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// write writes a file of a work area, with its directories.
func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The usage and input errors of layup setup give exit 2, nothing on standard
// output, and the reason on standard error: a work area with no problem
// statement too (D5 of #86).
func TestSetup(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	unasked := t.TempDir() // an answer to a question that S01 does not ask (finding 4 of round 1)
	write(t, filepath.Join(unasked, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	write(t, filepath.Join(unasked, "inputs", "answers.tsv"), "question\tanswer\tby\tsource\tquestion_text\nS01-wrong\tx\toperator\tu\t—\n")
	nobrief := t.TempDir()
	for _, c := range []struct {
		args          []string
		reason, usage string
	}{
		{[]string{"setup"}, "missing argument WORK", "\nusage: layup "},
		{[]string{"setup", "--x", "y", "w"}, `unknown flag "--x"`, "\nusage: layup "},
		{[]string{"setup", missing}, "the work area " + missing + " is not a directory", ""},
		{[]string{"setup", nobrief}, "inputs/briefs/problem-statement.md: the problem statement is absent", ""},
		{[]string{"setup", unasked}, "inputs/answers.tsv: line 2, S01-wrong: S01 does not ask this question", ""},
	} {
		r := layup(t, c.args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: "+c.reason+"\n") || (c.usage == "") == strings.Contains(r.stderr, "usage:") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing and %q", c.args, r.code, r.stdout, r.stderr, c.reason)
		}
	}
}

// The demo of #86: on a new work area, S01 stops with one table of its four
// questions (exit 3), the same on each run; with the answers, the next run does
// S01 to S04 on the stand-in baseline by its file:// URL, and S05 is not built
// yet (exit 1); a third run gives the same table and starts at S05.
func TestSetupRunsS01ToS04(t *testing.T) {
	tmp := t.TempDir()
	url, _, err := standin.Baseline(filepath.Join(tmp, "baseline"))
	if err != nil {
		t.Fatal(err)
	}
	w := filepath.Join(tmp, "work")
	write(t, filepath.Join(w, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	stop := "step\tquestion\task\twhere\n"
	for _, q := range work.S01Questions {
		stop += "S01\t" + q.ID + "\t" + q.Text + "\t—\n"
	}
	if r := repeat(t, "setup", w); r.code != 3 || r.stdout != stop || !strings.Contains(r.stderr, "layup setup: [1/15] S01\n") {
		t.Fatalf("a new work area: exit %d, stdout %q, stderr %q; want 3 and the stop table\n%s", r.code, r.stdout, r.stderr, stop)
	}
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	var b bytes.Buffer
	if err := tsv.Write(&b, work.AnswersSchema, [][]string{{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", url, "operator", at, ""}}); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(w, "inputs", "answers.tsv"), b.String())
	r := layup(t, "setup", w)
	if r.code != 1 || !strings.HasPrefix(r.stdout, "step\tactor\tresult\tevidence\n"+
		"S01\tlayup-setup\tdone\tevery answer present; the stack has a catalog entry\n"+
		"S02\tlayup-setup\tdone\tthe commit and the tree\n"+
		"S03\tlayup-setup\tdone\troot tree = pin.tree\n"+
		"S04\tlayup-setup\tdone\tchecks pin and facts\n"+
		"S05\tlayup-setup\tnot-active\tnot built yet\n"+
		"S06\tlayup-setup\tnot-active\tnot run: S05 did not pass\n") || !strings.Contains(r.stderr, "layup setup: [4/15] S04\n") {
		t.Fatalf("the run with the answers: exit %d, stdout %q, stderr %q; want 1, S01 to S04 done and S05 not built yet", r.code, r.stdout, r.stderr)
	}
	again := layup(t, "setup", w)
	if again.code != 1 || again.stdout != r.stdout || !strings.HasPrefix(again.stderr, "layup setup: [1/11] S05\n") {
		t.Errorf("a third run: exit %d, stdout %q, stderr %q; want the same table, from S05", again.code, again.stdout, again.stderr)
	}
}
