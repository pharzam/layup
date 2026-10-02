//go:build integration

package work

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go schemas of the two records, and of the open gaps of a target (#87),
// equal their blocks of docs/spec/setup.md (K9: this package is their one
// home).
func TestTheSchemasEqualTheirBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"setup-answers": AnswersSchema, "setup-record": RecordSchema, "open-gaps": OpenGapsSchema} {
		block, ok := blocks[name]
		if !ok {
			t.Errorf("docs/spec/ has no block %s", name)
			continue
		}
		if err := tsv.Compare(block, s); err != nil {
			t.Error(err)
		}
	}
}

// The readers read the two files of a work area, and an error names the file
// from the root of the work area.
func TestReadTheTwoFilesOfAWorkArea(t *testing.T) {
	dir := t.TempDir()
	for name, text := range map[string]string{RecordPath: record, AnswersPath: answers} {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if r, err := ReadRecord(dir); err != nil || len(r) != 3 {
		t.Errorf("ReadRecord: %d rows, %v; want 3", len(r), err)
	}
	if a, err := ReadAnswers(dir); err != nil || len(a) != 1 {
		t.Errorf("ReadAnswers: %d rows, %v; want 1", len(a), err)
	}
	if err := os.WriteFile(filepath.Join(dir, "out", "record.tsv"), []byte(answers), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadRecord(dir); err == nil || !strings.HasPrefix(err.Error(), "out/record.tsv: line 1") {
		t.Errorf("ReadRecord of a file of another form: %v; want an error that starts with out/record.tsv: line 1", err)
	}
	if _, err := ReadAnswers(t.TempDir()); err == nil || !strings.HasPrefix(err.Error(), "inputs/answers.tsv: ") {
		t.Errorf("ReadAnswers with no file: %v; want an error that starts with inputs/answers.tsv: ", err)
	}
}
