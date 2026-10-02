//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The usage and input errors of layup setup give exit 2, nothing on standard
// output, and the reason on standard error; a new work area gives exit 1, with
// S01 not built yet and each later step not run (D9 of #85), until rows 9, 13
// and 15 build the steps.
func TestSetup(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	unasked := t.TempDir() // an answer to a question that S01 does not ask (finding 4 of round 1)
	if err := os.MkdirAll(filepath.Join(unasked, "inputs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(unasked, "inputs", "answers.tsv"),
		[]byte("question\tanswer\tby\tsource\tquestion_text\nS01-wrong\tx\toperator\tu\t\u2014\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		args          []string
		reason, usage string
	}{
		{[]string{"setup"}, "missing argument WORK", "\nusage: layup "},
		{[]string{"setup", "--x", "y", "w"}, `unknown flag "--x"`, "\nusage: layup "},
		{[]string{"setup", missing}, "the work area " + missing + " is not a directory", ""},
		{[]string{"setup", unasked}, "inputs/answers.tsv: line 2, S01-wrong: S01 does not ask this question", ""},
	} {
		r := layup(t, c.args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: "+c.reason+"\n") || (c.usage == "") == strings.Contains(r.stderr, "usage:") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing and %q", c.args, r.code, r.stdout, r.stderr, c.reason)
		}
	}
	r := repeat(t, "setup", t.TempDir())
	if r.code != 1 || !strings.HasPrefix(r.stdout, "step\tactor\tresult\tevidence\nS01\tlayup-setup\tnot-active\tnot built yet\nS02\tlayup-setup\tnot-active\tnot run: S01 did not pass\n") ||
		strings.Count(r.stdout, "\n") != 16 || !strings.Contains(r.stderr, "layup setup: [1/15] S01\n") {
		t.Errorf("a new work area: exit %d, stdout %q, stderr %q; want 1, S01 not built yet, and the progress line", r.code, r.stdout, r.stderr)
	}
}
