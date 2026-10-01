package tsv

import (
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

// spec has the form of a file of docs/spec/: the README's example of the form
// in a fence of four backticks, a text fence, a tilde fence, and two blocks.
const spec = "# Records\n\n" +
	"````text\n```tsv-schema <name> <location>\n<column> <type> <key> <rule>\n```\n````\n\n" +
	"```text\nlayup psb check FILE\n```\n\n" +
	"```tsv-schema psb-gaps stdout\n" +
	"id        id(Q-NNN)          key  the row number, from Q-001\n" +
	"rule      enum(G1|G2|G3|G4|G5)  -  the rule that found the gap  \n" +
	"```\n\n" +
	"~~~text\n```tsv-schema hidden stdout\nx text - a fence inside a fence\n```\n~~~\n\n" +
	"  ```tsv-schema open-gaps target:docs/setup/open-gaps.tsv no-header\n" +
	"file      path  key  the file that holds the marker\n" +
	"question  text  -    never `—`\n" +
	"  ```\n"

func TestParseBlocksReadsTheFormOfTheREADME(t *testing.T) {
	got, err := parseBlocks(spec)
	want := []block{
		{Schema{Name: "psb-gaps", Location: "stdout", Columns: []Column{
			{Name: "id", Type: "id(Q-NNN)", Key: true, Rule: "the row number, from Q-001"},
			{Name: "rule", Type: "enum(G1|G2|G3|G4|G5)", Rule: "the rule that found the gap"},
		}}, 13},
		{Schema{Name: "open-gaps", Location: "target:docs/setup/open-gaps.tsv", NoHeader: true, Columns: []Column{
			{Name: "file", Type: "path", Key: true, Rule: "the file that holds the marker"},
			{Name: "question", Type: "text", Rule: "never `—`"},
		}}, 24},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("got %+v, error %v\nwant %+v", got, err, want)
	}
}

func TestParseBlocksRefusesABlockThatDoesNotHaveTheForm(t *testing.T) {
	const open, col = "```tsv-schema a stdout\n", "x text - r\n"
	for _, c := range []struct{ name, text, want string }{
		{"a tab in a column line", open + "x\ttext - r\n```\n", "line 2: a tab"},
		{"a tab in the opening line", "```tsv-schema\ta stdout\n" + col + "```\n", "line 1: a tab"},
		{"an unknown type", open + "x float - r\n```\n", `line 2: column x: type "float"`},
		{"a key word that is not key or -", open + "x text yes r\n```\n", `line 2: column x: "yes"`},
		{"a column that repeats", open + col + "x int - r\n```\n", "line 3: the column x repeats"},
		{"a column line with no rule", open + "x text -\n```\n", "line 2: \"x text -\""},
		{"an empty line", open + col + "\n```\n", "line 3: an empty line"},
		{"no column", open + "```\n", "line 2: the block a has no column"},
		{"no closing fence", open + col, "line 1: the block a has no closing fence"},
		{"no location", "```tsv-schema a\n" + col + "```\n", "line 1: \"tsv-schema a\""},
		{"a fourth word that is not no-header", "```tsv-schema a stdout header\n" + col + "```\n", "line 1: \"tsv-schema a stdout header\""},
		{"another location prefix", "```tsv-schema a file:x.tsv\n" + col + "```\n", "has no prefix"},
		{"stdout with a path", "```tsv-schema a stdout:x\n" + col + "```\n", "has no prefix"},
		{"a location with no path", "```tsv-schema a records:\n" + col + "```\n", "not relative"},
		{"a location with ..", "```tsv-schema a target:../x.tsv\n" + col + "```\n", "not relative"},
		{"a location with a leading /", "```tsv-schema a host:/x.tsv\n" + col + "```\n", "not relative"},
	} {
		if _, err := parseBlocks(c.text); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q in it", c.name, err, c.want)
		}
	}
}

func TestReadBlocksRefusesTwoBlocksWithOneName(t *testing.T) {
	a := "```tsv-schema a stdout\nx text - r\n```\n"
	fsys := fstest.MapFS{
		"a.md":       {Data: []byte(a)},
		"b.md":       {Data: []byte("text\n\n" + strings.Replace(a, " a ", " b ", 1))},
		"c.txt":      {Data: []byte(a)},
		"sub/d.md":   {Data: []byte(a)},
		"README.tsv": {Data: []byte("x")},
	}
	got, err := ReadBlocks(fsys)
	if err != nil || len(got) != 2 || got["a"].Location != "stdout" || got["b"].Columns[0].Name != "x" {
		t.Fatalf("got %+v, error %v; want the blocks a and b", got, err)
	}
	for _, c := range []struct{ name, b, want string }{
		{"in two files", a, "b.md:1: the block a has the name of the block at a.md:1"},
		{"in one file", "```tsv-schema b stdout\nx text - r\n```\n" + strings.Replace(a, " a ", " b ", 1), "b.md:4: the block b has the name of the block at b.md:1"},
		{"a block that does not have the form", "```tsv-schema b stdout\nx text\n```\n", "b.md: line 2: "},
	} {
		fsys["b.md"] = &fstest.MapFile{Data: []byte(c.b)}
		if _, err := ReadBlocks(fsys); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q in it", c.name, err, c.want)
		}
	}
}

func TestCompareNamesEachDifference(t *testing.T) {
	b := gaps
	if err := Compare(b, gaps); err != nil {
		t.Fatalf("the same schema: %v", err)
	}
	ruled := Schema{Name: gaps.Name, Location: gaps.Location, Columns: append([]Column{}, gaps.Columns...)}
	ruled.Columns[0].Rule = "a rule in other words"
	if err := Compare(b, ruled); err != nil {
		t.Fatalf("another rule: %v; the rules are not compared", err)
	}
	for _, c := range []struct {
		name   string
		change func(s *Schema)
		want   string
	}{
		{"the name", func(s *Schema) { s.Name = "psb-gap" }, `the name: the block has "psb-gaps"; the Go schema has "psb-gap"`},
		{"the location", func(s *Schema) { s.Location = "records:gaps.tsv" }, "the location: "},
		{"no-header", func(s *Schema) { s.NoHeader = true }, `no-header: the block has "false"`},
		{"a column less", func(s *Schema) { s.Columns = s.Columns[:4] }, `the number of columns: the block has "5"; the Go schema has "4"`},
		{"the order", func(s *Schema) { s.Columns[1], s.Columns[2] = s.Columns[2], s.Columns[1] }, `column 2 (rule): the name: the block has "rule"; the Go schema has "line"`},
		{"a type", func(s *Schema) { s.Columns[2].Type = "text" }, `column 3 (line): the type: the block has "int"; the Go schema has "text"`},
		{"a key mark", func(s *Schema) { s.Columns[0].Key = false }, `column 1 (id): key: the block has "true"`},
	} {
		g := Schema{Name: gaps.Name, Location: gaps.Location, Columns: append([]Column{}, gaps.Columns...)}
		c.change(&g)
		if err := Compare(b, g); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q in it", c.name, err, c.want)
		}
	}
}
