//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The built binary prints its version and exits 0, two runs give the same
// bytes (NFR-005), and it needs no program on the PATH for it.
func TestVersion(t *testing.T) {
	r := repeat(t, "version")
	if r.code != 0 || r.stdout != "layup 0.1.0-dev\n" || r.stderr != "" {
		t.Fatalf("exit %d, stdout %q, stderr %q; want 0, %q and nothing", r.code, r.stdout, r.stderr, "layup 0.1.0-dev\n")
	}
	if e := layupWith(t, []string{"PATH="}, "version"); e != r {
		t.Fatalf("with an empty PATH: %+v; want %+v", e, r)
	}
}

// Each usage error of the built binary gives exit code 2, the reason and the
// usage on standard error, and nothing on standard output (the demo of #80).
func TestUsageErrors(t *testing.T) {
	file := filepath.Join(t.TempDir(), "psb.md")
	if err := os.WriteFile(file, []byte("**Technology stack:** Go.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{
		{}, {"frobnicate"}, {"-h"}, {"--help"},
		{"version", "extra"}, {"version", "--x"},
		{"psb"}, {"psb", "check"}, {"psb", "chek", file}, {"psb", "check", file, file}, {"psb", "check", "--x", file},
	} {
		r := layup(t, args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: ") || !strings.Contains(r.stderr, "\nusage: layup ") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing, and the reason with the usage", args, r.code, r.stdout, r.stderr)
		}
	}
}
