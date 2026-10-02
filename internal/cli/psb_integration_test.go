//go:build integration

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, path, text string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// psb check on real files: 0 with the header only, 1 with a gap, and 2 with
// nothing on standard output and the reason without the usage for a file that
// is missing, a directory, or not valid UTF-8.
func TestPSBCheckExitCodes(t *testing.T) {
	dir := t.TempDir()
	clean := filepath.Join(dir, "clean.md")
	gaps := filepath.Join(dir, "gaps.md")
	bad := filepath.Join(dir, "bad.md")
	writeFile(t, clean, "**Technology stack:** Go.\n")
	writeFile(t, gaps, "No stack is named here.\n")
	writeFile(t, bad, "Technology stack: Go\nThe API is f\xffst.\n")

	if code, out, _ := run("psb", "check", clean); code != 0 || out != "id\trule\tline\texcerpt\tquestion\n" {
		t.Fatalf("clean: exit %d, stdout %q; want 0 and the header only", code, out)
	}
	if code, out, _ := run("psb", "check", gaps); code != 1 || !strings.Contains(out, "\tG1\t") {
		t.Fatalf("gaps: exit %d, stdout %q; want 1 and a G1 row", code, out)
	}
	for _, c := range []struct{ name, file, reason string }{
		{"a missing file", filepath.Join(dir, "missing.md"), "no such file or directory"},
		{"a directory", dir, "is a directory"},
		{"a file that is not valid UTF-8", bad, bad + ": line 2 is not valid UTF-8"},
	} {
		code, out, errOut := run("psb", "check", c.file)
		if code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: ") || !strings.Contains(errOut, c.reason) || strings.Contains(errOut, "usage:") {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 2, nothing, and the reason %q with no usage", c.name, code, out, errOut, c.reason)
		}
	}
}
