package verify

import (
	"testing"
	"testing/fstest"

	"github.com/pharzam/layup/internal/work"
)

// Check sources: each value row of the setup record has a source that
// resolves (D4 of #87, with note 4 of its plan review).
func TestTheSourcesOfARecord(t *testing.T) {
	saved := catalogHas
	t.Cleanup(func() { catalogHas = saved })
	catalogHas = func(stack, ref string) bool { return stack == "go" && ref == "go/files/go.mod" }
	tree := fstest.MapFS{
		"docs/facts/F-0001-setup-answers.md": {Data: []byte("1. `S01-stack` go\n")},
		"docs/a.md":                          {Data: []byte("one\nthe \u2039port\u203a here\n")},
		"docs/setup/open-gaps.tsv":           {Data: []byte("docs/a.md\t\u2039port\u203a\tWhich port?\n")},
		"docs/x.txt":                         {Data: []byte("x")},
	}
	area := fstest.MapFS{"inputs/answers.tsv": {Data: []byte("x")}}
	answers := work.Answers{{"S01-stack", "go", "operator", "u", ""}}
	base := [][]string{{"S01", "stack", "go", "answer", "S01-stack"}}
	for _, c := range []struct {
		name string
		row  []string
		want []string
	}{
		{"an answer", []string{"S04", "x", "go", "answer", "S01-stack"}, nil},
		{"an answer with no row", []string{"S04", "x", "go", "answer", "S01-name"},
			[]string{"source: S04 x: the answer S01-name is not a row of inputs/answers.tsv"}},
		{"a catalog file", []string{"S12", "go.mod", "x", "catalog", "go/files/go.mod"}, nil},
		{"a catalog ref that names no file", []string{"S12", "x", "x", "catalog", "go/files/nope"},
			[]string{"source: S12 x: the catalog entry go has no file go/files/nope"}},
		{"a fact", []string{"S05", "x", "x", "fact", "F-0001#1"}, nil},
		{"a fact that does not resolve", []string{"S05", "x", "x", "fact", "F-0001#9"},
			[]string{"source: S05 x: fact: F-0001#9 is not a fact of the F-0001 record"}},
		{"a gap", []string{"S11", "marker:docs/a.md:2", "\u2039port\u203a", "gap", "docs/setup/open-gaps.tsv"}, nil},
		{"a gap whose line has no marker", []string{"S11", "marker:docs/a.md:1", "\u2039port\u203a", "gap", "docs/setup/open-gaps.tsv"},
			[]string{"source: S11 marker:docs/a.md:1: line 1 of docs/a.md does not hold \u2039port\u203a"}},
		{"a gap with no row of open gaps", []string{"S11", "marker:docs/a.md:2", "\u2039other\u203a", "gap", "docs/setup/open-gaps.tsv"},
			[]string{"source: S11 marker:docs/a.md:2: line 2 of docs/a.md does not hold \u2039other\u203a",
				"source: S11 marker:docs/a.md:2: docs/setup/open-gaps.tsv has no row for docs/a.md \u2039other\u203a"}},
		{"a gap row of another name", []string{"S11", "x", "\u2039port\u203a", "gap", "docs/setup/open-gaps.tsv"},
			[]string{"source: S11 x: a gap row is not named marker:<file>:<line>"}},
		{"a command", []string{"S02", "pin.commit", "abc", "computed", "git ls-remote x HEAD"}, nil},
		{"a hash of a file of the tree", []string{"S06", "x", "abc", "computed", "sha256 docs/x.txt"}, nil},
		{"a hash of rows of a file of the work area", []string{"S01", "answers.sha256", "abc", "computed", "sha256 inputs/answers.tsv S01- Q-"}, nil},
		{"a hash of no file", []string{"S06", "x", "abc", "computed", "sha256 docs/nope"},
			[]string{"source: S06 x: sha256 names docs/nope, which is not a file of the tree or of the work area"}},
		{"a computed row with no ref", []string{"S02", "x", "abc", "computed", ""}, []string{"source: S02 x: no ref"}},
		{"a value row with the source step", []string{"S03", "x", "abc", "step", ""}, []string{"source: S03 x: a value row with the source step"}},
		{"a done row is no value row", []string{"S03", "done", "ok", "step", ""}, nil},
	} {
		in := input{fsys: tree, area: area, answers: answers, record: append(work.Record{base[0]}, c.row)}
		same(t, c.name, checkSources(in), c.want)
	}
	catalogHas = func(stack, ref string) bool { return false }
	in := input{fsys: tree, area: area, answers: answers, record: work.Record{{"S12", "x", "x", "catalog", "go/files/go.mod"}}}
	same(t, "a catalog row with no stack row", checkSources(in), []string{"source: S12 x: the catalog entry  has no file go/files/go.mod", "record: no value at S01 stack"})
}

// Check sources resolves a catalog ref in the entry of the binary (D10 of
// #91): the seam of row 10 reads the embedded Go entry, and the binary has no
// test entry.
func TestACatalogRefOfTheBinary(t *testing.T) {
	for _, c := range []struct {
		stack, ref string
		want       bool
	}{
		{"go", "go/kinds.tsv", true},
		{"go", "go/files/go.mod.tmpl", true},
		{"go", "go/files/.github/workflows/gates.yml.tmpl", true},
		{"go", "go/gaps.tsv", true},
		{"go", "go/fixtures/test.patch", true},
		{"go", "go/files/none.tmpl", false},
		{"go", "go/go.mod", false},
		{"go", "test/kinds.tsv", false},
		{"test", "test/kinds.tsv", false},
	} {
		if got := catalogHas(c.stack, c.ref); got != c.want {
			t.Errorf("catalogHas(%q, %q) = %v, want %v", c.stack, c.ref, got, c.want)
		}
	}
}
