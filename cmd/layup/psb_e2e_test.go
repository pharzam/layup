//go:build e2e

package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The built binary on the real problem statement of LAYUP prints the golden
// table of internal/psb byte for byte and exits 1 (REQ-001), two runs give
// the same bytes (NFR-005), and a run with no environment variable and no
// standard input gives the same result (the demo of #83).
func TestPSBCheckOnTheRealProblemStatement(t *testing.T) {
	file, err := filepath.Abs(filepath.Join("..", "..", "docs", "facts", "problem-statement-brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile(filepath.Join("..", "..", "internal", "psb", "testdata", "psb.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	r := repeat(t, "psb", "check", file)
	if r.code != 1 || r.stderr != "" {
		t.Fatalf("exit %d, stderr %q; want 1 and nothing", r.code, r.stderr)
	}
	if err := sameBytes([]byte(r.stdout), want); err != nil {
		t.Fatalf("the table is not internal/psb/testdata/psb.tsv: %v", err)
	}
	if b := layupBare(t, "psb", "check", file); b != r {
		t.Fatalf("with no environment variable: exit %d, stderr %q; want the result of a normal run", b.code, b.stderr)
	}
}

// A problem statement with no gap gives the header row only and exit 0.
func TestPSBCheckOnAStatementWithNoGap(t *testing.T) {
	file := filepath.Join(t.TempDir(), "clean.md")
	if err := os.WriteFile(file, []byte("**Technology stack:** Go.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if r := repeat(t, "psb", "check", file); r.code != 0 || r.stdout != "id\trule\tline\texcerpt\tquestion\n" || r.stderr != "" {
		t.Fatalf("exit %d, stdout %q, stderr %q; want 0, the header only, and nothing", r.code, r.stdout, r.stderr)
	}
}

// A FILE that is missing, a directory, or not valid UTF-8 gives exit 2,
// nothing on standard output, and the reason without the usage on standard
// error; for UTF-8 the reason names the file and the line (K32).
func TestPSBCheckInputErrors(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.md")
	if err := os.WriteFile(bad, []byte("Technology stack: Go\n\nThe API is f\xffst.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ name, file, reason string }{
		{"a missing file", filepath.Join(dir, "missing.md"), "no such file or directory"},
		{"a directory", dir, "is a directory"},
		{"a file that is not valid UTF-8", bad, "layup: " + bad + ": line 3 is not valid UTF-8\n"},
	} {
		r := layup(t, "psb", "check", c.file)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: ") || !strings.Contains(r.stderr, c.reason) || strings.Contains(r.stderr, "usage:") {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 2, nothing, and the reason %q with no usage", c.name, r.code, r.stdout, r.stderr, c.reason)
		}
	}
}

// A table that the binary cannot write gives exit 2 and the error on standard
// error: here its standard output is a file that is open for reading only, so
// each write fails. (A closed pipe stops a Go program with SIGPIPE instead.)
func TestPSBCheckWithAReadOnlyStandardOutput(t *testing.T) {
	dir := t.TempDir()
	in, out := filepath.Join(dir, "in.md"), filepath.Join(dir, "out")
	for _, f := range []string{in, out} {
		if err := os.WriteFile(f, []byte("No stack is named here.\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	readOnly, err := os.Open(out)
	if err != nil {
		t.Fatal(err)
	}
	defer readOnly.Close()
	var stderr bytes.Buffer
	cmd := exec.Command(binary, "psb", "check", in)
	cmd.Dir, cmd.Env, cmd.Stdout, cmd.Stderr = dir, []string{}, readOnly, &stderr
	var exit *exec.ExitError
	if err := cmd.Run(); !errors.As(err, &exit) || exit.ExitCode() != 2 || !strings.HasPrefix(stderr.String(), "layup: write /dev/stdout: ") {
		t.Fatalf("%v, stderr %q; want exit 2 and the error of the write", err, stderr.String())
	}
}
