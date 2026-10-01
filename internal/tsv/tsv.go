// Package tsv reads and writes the records of LAYUP, each by its schema
// (docs/spec/README.md, Records). It also parses the tsv-schema blocks of
// docs/spec/ and compares a block with a Go schema, so that the package that
// owns a record can test the two. It imports no package of this module.
package tsv

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	emptyMark = "—"      // the field of an empty value (U+2014), never ""
	bom       = "\uFEFF" // the byte-order mark of UTF-8
)

// Schema is the form of one record: its tsv-schema block, or the same form as
// a Go value in the package that owns the record.
type Schema struct {
	Name     string   // unique in docs/spec/
	Location string   // stdout, or records:, target:, layup: or host: and a path
	NoHeader bool     // the record has no header row
	Columns  []Column // in the column order
}

// Column is one column of a record.
type Column struct {
	Name string // its field in the header row
	Type string // a type of the closed list, for example int or id(Q-NNN)
	Key  bool   // part of the record's key
	Rule string // the rule in words, for a person; Compare does not compare it
}

// Error is a line of a record that does not match its schema. Line counts from
// 1, and a header row is line 1. Column names the column of the field that does
// not match, or is "" when the error is about the whole line.
type Error struct {
	Line   int
	Column string
	Reason string
}

func (e *Error) Error() string {
	if e.Column == "" {
		return fmt.Sprintf("line %d: %s", e.Line, e.Reason)
	}
	return fmt.Sprintf("line %d, column %q: %s", e.Line, e.Column, e.Reason)
}

// fieldRule makes each tab, line feed and carriage return of a value one space.
var fieldRule = strings.NewReplacer("\t", " ", "\n", " ", "\r", " ")

// Write writes rows as a record of schema s: the header row, unless s has
// none, then one line per row, each with a line feed at its end. It applies the
// field rule and writes an empty value as —. It refuses a value that is not
// valid UTF-8 (an error, not a repair), and it makes each check that Read
// makes, so it never writes a record that Read refuses. On an error it writes
// nothing.
func Write(w io.Writer, s Schema, rows [][]string) error {
	types, err := s.types()
	if err != nil {
		return err
	}
	var b strings.Builder
	n, keys := 0, map[string]int{} // n is the line of the row
	if !s.NoHeader {
		names := make([]string, len(s.Columns))
		for j, c := range s.Columns {
			names[j] = c.Name
		}
		b.WriteString(strings.Join(names, "\t") + "\n")
		n = 1
	}
	for _, row := range rows {
		n++
		fields := make([]string, len(row))
		for j, v := range row {
			if fields[j] = fieldRule.Replace(v); fields[j] == "" {
				fields[j] = emptyMark
			}
		}
		if err := checkRow(n, fields, s.Columns, types, keys); err != nil {
			return err
		}
		b.WriteString(strings.Join(fields, "\t") + "\n")
	}
	if strings.HasPrefix(b.String(), bom) {
		return &Error{Line: 1, Column: s.Columns[0].Name, Reason: "the record would start with a byte-order mark, which the reader refuses"}
	}
	_, err = io.WriteString(w, b.String())
	return err
}

// Read reads a record of schema s and returns its rows, one value per column;
// — gives the empty value. It refuses a byte-order mark, a carriage return
// anywhere, an empty line, no line feed after the last line, a header row that
// is not the column names, and each row that Write refuses. The error is an
// *Error for the first line that does not match, so a command can give exit 2.
func Read(data []byte, s Schema) ([][]string, error) {
	types, err := s.types()
	if err != nil {
		return nil, err
	}
	text := string(data)
	switch {
	case strings.HasPrefix(text, bom):
		return nil, &Error{Line: 1, Reason: "a byte-order mark; save the file as UTF-8 with no byte-order mark"}
	case text == "" && s.NoHeader:
		return nil, nil
	case text == "":
		return nil, &Error{Line: 1, Reason: "no header row"}
	}
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	var rows [][]string
	keys := map[string]int{}
	for i, l := range lines {
		n, fields := i+1, strings.Split(l, "\t")
		switch {
		case strings.Contains(l, "\r"):
			err = &Error{Line: n, Reason: "a carriage return; save the file with line-feed endings only, and with no carriage return in a field"}
		case l == "":
			err = &Error{Line: n, Reason: "an empty line"}
		case n == 1 && !s.NoHeader:
			err = checkHeader(fields, s.Columns)
		default:
			if err = checkRow(n, fields, s.Columns, types, keys); err == nil {
				for j := range fields {
					if fields[j] == emptyMark {
						fields[j] = ""
					}
				}
				rows = append(rows, fields)
			}
		}
		if err != nil {
			return nil, err
		}
	}
	if !strings.HasSuffix(text, "\n") {
		return nil, &Error{Line: len(lines), Reason: "no line feed after the last line"}
	}
	return rows, nil
}

// JoinList gives the field of a list(<type>) column: the values, joined by one
// space. It refuses a value that is empty or holds a space, a tab, a line feed
// or a carriage return, because it would read back as more values or as none.
// No values give the empty string, which Write writes as —. strings.Fields
// gives the values of a field back.
func JoinList(values []string) (string, error) {
	for i, v := range values {
		if v == "" || strings.ContainsAny(v, " \t\n\r") {
			return "", fmt.Errorf("value %d of the list, %q: a list value is not empty and holds no space, tab, line feed or carriage return", i+1, v)
		}
	}
	return strings.Join(values, " "), nil
}

// types parses the type of each column. It refuses a schema with no column,
// and a column name that is empty, repeats, is not valid UTF-8, or holds a tab,
// a line feed or a carriage return.
func (s Schema) types() ([]fieldType, error) {
	if len(s.Columns) == 0 {
		return nil, fmt.Errorf("schema %s: no column", s.Name)
	}
	seen := map[string]bool{}
	types := make([]fieldType, len(s.Columns))
	for j, c := range s.Columns {
		if !oneLine(c.Name) || seen[c.Name] {
			return nil, fmt.Errorf("schema %s: column %d: the name %q is empty, repeats, is not valid UTF-8, or holds a tab, a line feed or a carriage return", s.Name, j+1, c.Name)
		}
		seen[c.Name] = true
		var err error
		if types[j], err = parseType(c.Type); err != nil {
			return nil, fmt.Errorf("schema %s: column %s: %v", s.Name, c.Name, err)
		}
	}
	return types, nil
}

// checkHeader checks the header row against the column names.
func checkHeader(fields []string, cols []Column) error {
	for j, c := range cols {
		switch {
		case j >= len(fields):
			return &Error{Line: 1, Column: c.Name, Reason: "the header row has no field for this column"}
		case fields[j] != c.Name:
			return &Error{Line: 1, Column: c.Name, Reason: fmt.Sprintf("the header row has %q in its place", fields[j])}
		}
	}
	if len(fields) > len(cols) {
		return &Error{Line: 1, Reason: fmt.Sprintf("the header row has %d fields; the schema has %d columns", len(fields), len(cols))}
	}
	return nil
}

// checkRow checks the fields of the row on line n, as the record holds them,
// and its key: the tuple of the fields of all key columns, which no earlier row
// may have. keys maps each key to its line.
func checkRow(n int, fields []string, cols []Column, types []fieldType, keys map[string]int) error {
	switch {
	case len(fields) < len(cols):
		return &Error{Line: n, Column: cols[len(fields)].Name, Reason: fmt.Sprintf("the row has no field for this column: %d fields; the schema has %d columns", len(fields), len(cols))}
	case len(fields) > len(cols):
		return &Error{Line: n, Reason: fmt.Sprintf("a field with no column: the row has %d fields; the schema has %d columns", len(fields), len(cols))}
	}
	var key, shown []string
	for j, f := range fields {
		c, reason := cols[j], ""
		switch {
		case !utf8.ValidString(f):
			reason = "not valid UTF-8"
		case f == "":
			reason = "an empty field; write — for an empty value"
		case f == emptyMark && c.Key:
			reason = "— in a key column; a key is never empty"
		case f != emptyMark:
			if err := types[j].check(f); err != nil {
				reason = err.Error()
			}
		}
		if reason != "" {
			return &Error{Line: n, Column: c.Name, Reason: reason}
		}
		if c.Key {
			key, shown = append(key, f), append(shown, c.Name+" "+strconv.Quote(f))
		}
	}
	k := strings.Join(key, "\t") // a field holds no tab
	if m, ok := keys[k]; ok && key != nil {
		return &Error{Line: n, Reason: fmt.Sprintf("the key (%s) repeats the key of line %d", strings.Join(shown, ", "), m)}
	}
	keys[k] = n
	return nil
}
