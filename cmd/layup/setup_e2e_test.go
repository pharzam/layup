//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
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
