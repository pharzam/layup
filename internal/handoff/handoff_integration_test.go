//go:build integration

package handoff

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

func spec(t *testing.T, file string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("..", "..", "docs", "spec", file))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// The Go value of the block transitions equals its block of
// docs/spec/records.md (task T-h17s, #348).
func TestTheSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["transitions"]
	if !ok {
		t.Fatal("docs/spec/ has no block transitions")
	}
	if err := tsv.Compare(block, Schema); err != nil {
		t.Error(err)
	}
}

var span = regexp.MustCompile("`([^`]+)`")

func words(cell string) []string {
	var out []string
	for _, m := range span.FindAllStringSubmatch(cell, -1) {
		out = append(out, m[1])
	}
	return out
}

// The default rows are the thirteen of the table "The checks of a handoff" of
// docs/spec/loop.md, in its order, with its kinds and its items.
func TestTheDefaultRowsAreTheTableOfTheSpecification(t *testing.T) {
	text := spec(t, "loop.md")
	i := strings.Index(text, "| From, to | Artifacts | Items |")
	if i < 0 {
		t.Fatal("loop.md has no table with the columns From, to; Artifacts; Items")
	}
	var rows []Row
	for _, line := range strings.Split(text[i:], "\n")[2:] {
		if !strings.HasPrefix(line, "|") {
			break
		}
		cells := strings.Split(strings.Trim(line, "| "), " | ")
		if len(cells) != 3 {
			t.Fatalf("a row of the table has %d cells: %q", len(cells), line)
		}
		from, to, _ := strings.Cut(cells[0], ", ")
		// "task plan, plan review": the first ", " parts the two steps; no step
		// holds a comma ("developer (build task), verifier").
		rows = append(rows, Row{From: from, To: to, Artifacts: words(cells[1]), Items: words(strings.NewReplacer("(`M2f`)", "", "(`M2g`)", "").Replace(cells[2]))})
	}
	got := Default()
	if len(rows) != len(got) {
		t.Fatalf("the specification has %d rows, Default %d", len(rows), len(got))
	}
	for k, r := range rows {
		if got[k].From != r.From || got[k].To != r.To || !slices.Equal(got[k].Artifacts, r.Artifacts) || !slices.Equal(got[k].Items, r.Items) {
			t.Errorf("row %d: Default %+v, loop.md %+v", k+1, got[k], r)
		}
	}
}

// The lists of kinds and of items are those of docs/spec/loop.md, "The
// transition table".
func TestTheListsAreThoseOfTheSpecification(t *testing.T) {
	text := spec(t, "loop.md")
	i := strings.Index(text, "**The kinds of artifact** that a row may name are:")
	j := strings.Index(text[i:], "\n- **The names of the items**")
	if i < 0 || j < 0 {
		t.Fatal("loop.md has no list of the kinds of artifact")
	}
	if got := words(text[i : i+j]); !slices.Equal(got, Kinds()) {
		t.Errorf("kinds %v, loop.md %v", Kinds(), got)
	}
	var items []string
	k := strings.Index(text, "| Item | What code checks |")
	for _, line := range strings.Split(text[k:], "\n")[2:] {
		if !strings.HasPrefix(line, "|") {
			break
		}
		items = append(items, words(strings.SplitN(line, " | ", 2)[0])...)
	}
	if !slices.Equal(items, Items()) {
		t.Errorf("items %v, loop.md %v", Items(), items)
	}
}

// The package may import only what the table of M2e says (internal/tsv and
// internal/records), and the table gives it its row.
func TestThePackageHasItsRowInTheTableOfM2e(t *testing.T) {
	text := spec(t, "packages.md")
	i := strings.Index(text, "### The table of M2e")
	if i < 0 {
		t.Fatal("packages.md has no table of M2e")
	}
	var row string
	for _, line := range strings.Split(text[i:], "\n") {
		if strings.HasPrefix(line, "| `internal/handoff`") {
			row = line
			break
		}
	}
	if row == "" {
		t.Fatal("the table of M2e has no row of internal/handoff")
	}
	cells := strings.Split(row, " | ")
	if got := words(cells[2]); !slices.Equal(got, []string{"internal/tsv", "internal/records"}) {
		t.Errorf("May import %v", got)
	}
}
