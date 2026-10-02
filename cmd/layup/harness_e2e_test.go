//go:build e2e

package main

// The end-to-end harness: TestMain builds the layup binary once, and each
// scenario runs it as a user would.

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var binary string // the layup binary that TestMain builds

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "layup-e2e-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binary = filepath.Join(dir, "layup")
	build := exec.Command("go", "build", "-o", binary, ".")
	build.Env = append(os.Environ(), "GOFLAGS=", "GOENV=off", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "go build: %v\n%s", err, out)
		os.RemoveAll(dir)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(dir)
	os.Exit(code)
}

// A result is what one run of the binary gave.
type result struct {
	stdout, stderr string
	code           int
}

// layup runs the binary as a user would, with the fixed environment of
// layupWith.
func layup(t *testing.T, args ...string) result {
	t.Helper()
	return layupWith(t, nil, args...)
}

// layupWith runs the binary in a new temporary directory, with no standard
// input, and with a fixed environment: HOME and TMPDIR in temporary
// directories, the PATH of the host, and LC_ALL=C. An entry NAME=value of env
// replaces the entry of the same name, or adds one. No other variable of the
// host reaches the binary.
func layupWith(t *testing.T, env []string, args ...string) result {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = t.TempDir()
	cmd.Env = []string{"HOME=" + t.TempDir(), "TMPDIR=" + t.TempDir(), "PATH=" + os.Getenv("PATH"), "LC_ALL=C"}
	for _, e := range env {
		name, _, _ := strings.Cut(e, "=")
		cmd.Env = append(slices.DeleteFunc(cmd.Env, func(f string) bool { return strings.HasPrefix(f, name+"=") }), e)
	}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("layup %q: %v", args, err)
		}
		code = exit.ExitCode()
	}
	return result{stdout.String(), stderr.String(), code}
}

// repeat runs the binary twice with the same arguments, and fails when the
// two exit codes or the two standard outputs differ by one byte (NFR-005).
func repeat(t *testing.T, args ...string) result {
	t.Helper()
	return repeatWith(t, nil, args...)
}

// repeatWith is repeat with the entries of layupWith.
func repeatWith(t *testing.T, env []string, args ...string) result {
	t.Helper()
	first, second := layupWith(t, env, args...), layupWith(t, env, args...)
	if err := sameBytes([]byte(first.stdout), []byte(second.stdout)); err != nil || first.code != second.code {
		t.Fatalf("layup %q twice: exit %d and %d; standard output: %v", args, first.code, second.code, err)
	}
	return first
}
