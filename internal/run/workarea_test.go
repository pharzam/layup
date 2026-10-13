package run

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The work area refuses, before it touches the disk, a file of the records
// whose path leaves inputs/files (docs/spec/scaffold.md, The work area): the
// records are LAYUP's own, and a path with .. is no path of the target.
func TestTheWorkAreaRefusesAFilePathThatLeavesItsDirectory(t *testing.T) {
	for _, path := range []string{"..", "../escape.md", "docs/../../escape.md", "/etc/passwd", "docs//glossary.md", "", ".", "docs/./glossary.md"} {
		dir := filepath.Join(t.TempDir(), "work", "acme", "target")
		err := buildWorkArea(dir, workSource{Brief: []byte("b"), Files: map[string][]byte{path: []byte("x")}}, "/no/such/clone", "main", "https://github.invalid/acme/target.git")
		if err == nil || !strings.Contains(err.Error(), "setup/inputs/files/") {
			t.Errorf("the path %q: %v; want an error that names the path in the records", path, err)
		}
		if _, statErr := os.Stat(dir); statErr == nil {
			t.Errorf("the path %q: the work area exists after the refusal; want no change of the disk", path)
		}
	}
}
