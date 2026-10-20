//go:build integration

package run

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestMain gives the package one temporary HOME: t.Setenv panics under
// t.Parallel, and one HOME keeps the host's ~/.gitconfig out, as the per-test
// HOME did. It reads go env GOCACHE before changing HOME so the sub-runs of go
// keep the host's build cache whether or not GOCACHE was exported.
func TestMain(m *testing.M) {
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go env GOCACHE: %v\n", err)
		os.Exit(1)
	}
	if err := os.Setenv("GOCACHE", strings.TrimSpace(string(cache))); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	home, err := os.MkdirTemp("", "layup-run-home-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.Setenv("HOME", home); err != nil {
		os.RemoveAll(home)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(home)
	os.Exit(code)
}
