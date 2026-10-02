//go:build integration

package psb

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/tsv"
)

// The real PSB (fact F-0001) must yield the golden batch, which holds the G1
// row: the gap that cost a question in the manual setup.
func TestGoldenRealPSB(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "docs", "facts", "problem-statement-brief.md"))
	if err != nil {
		t.Fatal(err)
	}
	golden(t, "the real PSB", src, read(t, "psb.tsv"))
}

// questions gives the question text of each rule, with the part that a gap
// fills, as docs/spec/psb-check.md (The rules) writes it.
var questions = map[string]*regexp.Regexp{
	"G1": regexp.MustCompile(`^Which technology stack does the product use \(languages, frameworks, tools\)\?$`),
	"G2": regexp.MustCompile(`^How is this metric measured\? The row gives no measurement method\.$`),
	"G3": regexp.MustCompile(`^What does "[A-Z][A-Z0-9]{1,5}" mean\? The terms table does not define it\.$`),
	"G4": regexp.MustCompile(`^Which number or threshold does "(fast|robust|soon|clean|better|handle)" stand for here\?$`),
	"G5": regexp.MustCompile(`^Which start value does the pilot use for this target, and who sets it\?$`),
}

// The Go schema of the gap table equals its block, and each golden table, the
// one of the real PSB included, reads back by the block and keeps the rules of
// the table that code can check (plan review of #83, condition 3).
func TestEveryGoldenIsARecordOfTheBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["psb-gaps"]
	if !ok {
		t.Fatal("docs/spec/ has no block psb-gaps")
	}
	if err := tsv.Compare(block, GapsSchema); err != nil {
		t.Error(err)
	}
	for _, name := range []string{"triggers", "clean", "edge", "values", "psb"} {
		rows, err := tsv.Read(read(t, name+".tsv"), block)
		if err != nil {
			t.Errorf("%s.tsv: %v", name, err)
			continue
		}
		for i, r := range rows {
			at := fmt.Sprintf("%s.tsv, row %d", name, i+1)
			if want := fmt.Sprintf("Q-%03d", i+1); r[0] != want {
				t.Errorf("%s: id %s, want %s", at, r[0], want)
			}
			if n := utf8.RuneCountInString(r[3]); n > 80 {
				t.Errorf("%s: the excerpt has %d characters, more than 80", at, n)
			}
			if (r[2] == "0") != (r[3] == "") {
				t.Errorf("%s: line %s with the excerpt %q; the excerpt is — exactly when the line is 0", at, r[2], r[3])
			}
			if q, ok := questions[r[1]]; !ok || !q.MatchString(r[4]) {
				t.Errorf("%s: the question of %s is not the text of its rule: %q", at, r[1], r[4])
			}
		}
	}
}
