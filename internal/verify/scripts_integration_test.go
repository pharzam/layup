//go:build integration

package verify

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// writeScript writes a file of a tree, with its directories.
func writeScript(t *testing.T, dir, path, text string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The baseline's scripts run with sh in a tree, as a gate command (D5 of
// #87): pass, fail with the exit or the signal, not-active with no script;
// their output goes to standard error. The call of K10 reads the lines of
// link-lint.sh on its standard error (D6).
func TestTheScriptsOnATree(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	writeScript(t, dir, discTests, "echo the discipline tests ran\n")
	writeScript(t, dir, linkLint, "echo 'FAIL  L2: docs/b.md:3 links x.md, but that file has no heading with that anchor' >&2\necho 'FAIL  L3: docs/a.md:1 links #y, but this file has no heading with that anchor' >&2\nexit 1\n")
	var out bytes.Buffer
	if got := scriptRow(dir, "discipline-tests", discTests, &out); got != (Row{"discipline-tests", "pass", ""}) || !strings.Contains(out.String(), "the discipline tests ran") {
		t.Errorf("an OK script: %q, the output %q", got, out.String())
	}
	if got := scriptRow(dir, "link-lint", linkLint, &out); got != (Row{"link-lint", "fail", "exit 1"}) || !strings.Contains(out.String(), "FAIL  L2") {
		t.Errorf("a failed script: %q, the output %q", got, out.String())
	}
	if files, err := BrokenLinks(dir); err != nil || !slices.Equal(files, []string{"docs/a.md", "docs/b.md"}) {
		t.Errorf("BrokenLinks: %q, %v", files, err)
	}
	writeScript(t, dir, discTests, "kill -TERM $$\n")
	if got := scriptRow(dir, "discipline-tests", discTests, &out); got != (Row{"discipline-tests", "fail", "signal terminated"}) {
		t.Errorf("a killed script: %q", got)
	}
	if err := os.Remove(filepath.Join(dir, filepath.FromSlash(linkLint))); err != nil {
		t.Fatal(err)
	}
	if got := scriptRow(dir, "link-lint", linkLint, &out); got != (Row{"link-lint", "not-active", "missing: docs/links/link-lint.sh"}) {
		t.Errorf("no script: %q", got)
	}
	if _, err := BrokenLinks(dir); err == nil {
		t.Error("BrokenLinks with no script: no error")
	}
}

// LAYUP's own two scripts pass on a clone of its HEAD, whose checkout keeps
// the eol=crlf fixtures of run-discipline-tests.sh (open question 2 of
// verify-baseline-scripts).
func TestTheScriptsOfLAYUP(t *testing.T) {
	isolate(t)
	root, _ := filepath.Abs(filepath.Join("..", ".."))
	clone := filepath.Join(t.TempDir(), "clone")
	env := os.Environ()
	for _, args := range [][]string{{"-c", "core.hooksPath=/dev/null", "clone", "-q", "--no-hardlinks", root, clone}} {
		cmd := exec.Command("git", args...)
		cmd.Env = env
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
	for name, script := range map[string]string{"discipline-tests": discTests, "link-lint": linkLint} {
		var out bytes.Buffer
		if got := scriptRow(clone, name, script, &out); got != (Row{name, "pass", ""}) {
			t.Errorf("%s on a clone of LAYUP: %q\n%s", name, got, out.String())
		}
	}
}
