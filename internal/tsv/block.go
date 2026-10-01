package tsv

import (
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"
)

// ReadBlocks returns the schema of each tsv-schema block of the Markdown files
// at the root of fsys, by name; for example ReadBlocks(os.DirFS("docs/spec")).
// It refuses a block that does not have the form of docs/spec/README.md (The
// schema block), and two blocks with one name. An error names the file and the
// line.
func ReadBlocks(fsys fs.FS) (map[string]Schema, error) {
	files, _ := fs.Glob(fsys, "*.md") // the pattern is valid
	schemas, at := map[string]Schema{}, map[string]string{}
	for _, f := range files {
		text, err := fs.ReadFile(fsys, f)
		if err != nil {
			return nil, err
		}
		blocks, err := parseBlocks(string(text))
		if err != nil {
			return nil, fmt.Errorf("%s: %v", f, err)
		}
		for _, b := range blocks {
			here := fmt.Sprintf("%s:%d", f, b.line)
			if there, ok := at[b.Name]; ok {
				return nil, fmt.Errorf("%s: the block %s has the name of the block at %s", here, b.Name, there)
			}
			schemas[b.Name], at[b.Name] = b.Schema, here
		}
	}
	return schemas, nil
}

// block is a tsv-schema block and the line of its opening fence.
type block struct {
	Schema
	line int
}

// parseBlocks returns the tsv-schema blocks of a Markdown text: each fenced
// code block whose info string starts with the word tsv-schema. A fence closes
// at a fence of the same character that is at least as long, so a fence inside
// a longer fence is text: the README's example of the form, inside a fence of
// four backticks, is not a block.
func parseBlocks(text string) ([]block, error) {
	var blocks []block
	var open string // the fence that is open, or ""
	var b *block    // the block that is open, or nil
	for i, l := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		n, err := i+1, error(nil)
		run, rest := fence(l)
		switch {
		case open == "" && run != "" && !(run[0] == '`' && strings.Contains(rest, "`")):
			open = run
			if words := strings.Fields(rest); len(words) > 0 && words[0] == "tsv-schema" {
				b = &block{line: n}
				err = blockHead(&b.Schema, words, l)
			}
		case open != "" && run != "" && run[0] == open[0] && len(run) >= len(open) && strings.TrimSpace(rest) == "":
			if b != nil && b.Columns == nil {
				err = fmt.Errorf("the block %s has no column", b.Name)
			} else if b != nil {
				blocks = append(blocks, *b)
			}
			open, b = "", nil
		case b != nil:
			err = columnLine(&b.Schema, l)
		}
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", n, err)
		}
	}
	if b != nil {
		return nil, fmt.Errorf("line %d: the block %s has no closing fence", b.line, b.Name)
	}
	return blocks, nil
}

// fence returns the run of three or more backticks or tildes that starts line
// l after at most three spaces, and the rest of the line; or "" and "".
func fence(l string) (run, rest string) {
	t := strings.TrimLeft(l, " ")
	if len(l)-len(t) > 3 || t == "" || (t[0] != '`' && t[0] != '~') {
		return "", ""
	}
	rest = strings.TrimLeft(t, t[:1])
	if len(t)-len(rest) < 3 {
		return "", ""
	}
	return t[:len(t)-len(rest)], rest
}

// blockHead reads the words of the opening line l of a block: tsv-schema, the
// name, the location, and no-header for a record with no header row.
func blockHead(s *Schema, w []string, l string) error {
	if strings.Contains(l, "\t") {
		return errors.New("a tab; a block holds no tab")
	}
	if len(w) < 3 || len(w) > 4 || (len(w) == 4 && w[3] != "no-header") {
		return fmt.Errorf("%q: a block opens with tsv-schema, a name, a location, and no-header or nothing", strings.Join(w, " "))
	}
	s.Name, s.Location, s.NoHeader = w[1], w[2], len(w) == 4
	prefix, path, found := strings.Cut(s.Location, ":")
	switch {
	case s.Location == "stdout":
	case !found || !slices.Contains([]string{"records", "target", "layup", "host"}, prefix):
		return fmt.Errorf("the location %q is not stdout and has no prefix records:, target:, layup: or host:", s.Location)
	case path == "." || !fs.ValidPath(path):
		return fmt.Errorf("the location %q: the path is not relative to the root in one form", s.Location)
	}
	return nil
}

// columnLine reads a line of a block: the column name, its type, key or -, and
// the rest of the line, which is the rule. One or more spaces separate them.
func columnLine(s *Schema, l string) error {
	var f [3]string
	rest := strings.TrimLeft(l, " ")
	for k := range f {
		f[k], rest, _ = strings.Cut(rest, " ")
		rest = strings.TrimLeft(rest, " ")
	}
	c := Column{Name: f[0], Type: f[1], Key: f[2] == "key", Rule: strings.TrimRight(rest, " ")}
	_, err := parseType(c.Type)
	switch {
	case strings.Contains(l, "\t"):
		return errors.New("a tab; a block holds no tab")
	case strings.TrimSpace(l) == "":
		return errors.New("an empty line; a block has one line per column")
	case c.Rule == "":
		return fmt.Errorf("%q: a column line has a name, a type, key or -, and a rule", l)
	case err != nil:
		return fmt.Errorf("column %s: %v", c.Name, err)
	case f[2] != "key" && f[2] != "-":
		return fmt.Errorf("column %s: %q is not key or -", c.Name, f[2])
	case slices.ContainsFunc(s.Columns, func(d Column) bool { return d.Name == c.Name }):
		return fmt.Errorf("the column %s repeats", c.Name)
	}
	s.Columns = append(s.Columns, c)
	return nil
}

// Compare returns an error that names each difference between a block of
// docs/spec/ and the Go schema s of its record: the name, the location, the
// no-header mark, the number of columns, and the name, the type and the key
// mark of each column. It does not compare the rules, which are for a person.
func Compare(block, s Schema) error {
	pairs := [][3]string{
		{"the name", block.Name, s.Name},
		{"the location", block.Location, s.Location},
		{"no-header", fmt.Sprint(block.NoHeader), fmt.Sprint(s.NoHeader)},
		{"the number of columns", fmt.Sprint(len(block.Columns)), fmt.Sprint(len(s.Columns))},
	}
	for j := range min(len(block.Columns), len(s.Columns)) {
		b, g := block.Columns[j], s.Columns[j]
		at := fmt.Sprintf("column %d (%s): ", j+1, b.Name)
		pairs = append(pairs, [3]string{at + "the name", b.Name, g.Name}, [3]string{at + "the type", b.Type, g.Type},
			[3]string{at + "key", fmt.Sprint(b.Key), fmt.Sprint(g.Key)})
	}
	var diffs []error
	for _, p := range pairs {
		if p[1] != p[2] {
			diffs = append(diffs, fmt.Errorf("%s: the block has %q; the Go schema has %q", p[0], p[1], p[2]))
		}
	}
	if diffs == nil {
		return nil
	}
	return fmt.Errorf("schema %s: %w", block.Name, errors.Join(diffs...))
}
