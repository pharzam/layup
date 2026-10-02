//go:build e2e

package main

import (
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
	for _, c := range []struct {
		args          []string
		reason, usage string
	}{
		{[]string{"setup"}, "missing argument WORK", "\nusage: layup "},
		{[]string{"setup", "--x", "y", "w"}, `unknown flag "--x"`, "\nusage: layup "},
		{[]string{"setup", missing}, "the work area " + missing + " is not a directory", ""},
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
