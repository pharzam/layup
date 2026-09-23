//go:build e2e

package main

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// Builds the real binary and runs it, as a user would.
func TestBinaryPrintsItsVersion(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "layup")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	out, err := exec.Command(bin, "version").Output()
	if err != nil {
		t.Fatalf("layup version: %v", err)
	}
	if string(out) != "layup 0.1.0-dev\n" {
		t.Fatalf("stdout %q, want %q", out, "layup 0.1.0-dev\n")
	}
}
