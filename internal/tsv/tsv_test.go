package tsv

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

// The schemas of these tests have the form of blocks of docs/spec/.
var (
	gaps = Schema{Name: "psb-gaps", Location: "stdout", Columns: []Column{
		{Name: "id", Type: "id(Q-NNN)", Key: true}, {Name: "rule", Type: "enum(G1|G2|G3|G4|G5)"},
		{Name: "line", Type: "int"}, {Name: "excerpt", Type: "text"}, {Name: "question", Type: "text"},
	}}
	// stalls has a key of two columns: one stall has three rows.
	stalls = Schema{Name: "stalls", Location: "records:stalls.tsv", Columns: []Column{
		{Name: "stall", Type: "id(ST-NNN)", Key: true}, {Name: "kind", Type: "enum(stall|diagnosis|outcome)", Key: true},
		{Name: "note", Type: "text"},
	}}
	openGaps = Schema{Name: "open-gaps", Location: "target:docs/setup/open-gaps.tsv", NoHeader: true, Columns: []Column{
		{Name: "file", Type: "path", Key: true}, {Name: "marker", Type: "text", Key: true}, {Name: "question", Type: "text"},
	}}
	manifest = Schema{Name: "gate-manifest", Location: "target:docs/gates.tsv", Columns: []Column{
		{Name: "kind", Type: "id(<word>)", Key: true}, {Name: "config", Type: "list(path)"},
	}}
)

const gapsHeader = "id\trule\tline\texcerpt\tquestion\n"

func write(t *testing.T, s Schema, rows ...[]string) string {
	t.Helper()
	var b strings.Builder
	if err := Write(&b, s, rows); err != nil {
		t.Fatalf("Write: %v", err)
	}
	return b.String()
}

// errorAt fails unless err is an *Error at line and column whose text holds part.
func errorAt(t *testing.T, err error, line int, column, part string) {
	t.Helper()
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("error %v; want an *Error at line %d", err, line)
	}
	if e.Line != line || e.Column != column || !strings.Contains(e.Error(), part) {
		t.Fatalf("error %q (line %d, column %q); want line %d, column %q, and %q in the text", e, e.Line, e.Column, line, column, part)
	}
}

func TestWriteAppliesTheFieldRule(t *testing.T) {
	got := write(t, gaps,
		[]string{"Q-001", "G1", "0", "", "Which stack?"},
		[]string{"Q-002", "G4", "12", "a\tb\nc\rd\r\ne", "—"},
	)
	want := gapsHeader + "Q-001\tG1\t0\t—\tWhich stack?\n" + "Q-002\tG4\t12\ta b c d  e\t—\n"
	if got != want {
		t.Fatalf("got %q\nwant %q", got, want)
	}
	if got := write(t, gaps); got != gapsHeader {
		t.Fatalf("no rows: got %q; want the header row only", got)
	}
}

func TestTheNoHeaderFormHasRowsOnly(t *testing.T) {
	rows := [][]string{{"docs/a.md", "m1", "Which value?"}, {"docs/b.md", "m1", ""}}
	out := write(t, openGaps, rows...)
	if want := "docs/a.md\tm1\tWhich value?\ndocs/b.md\tm1\t—\n"; out != want {
		t.Fatalf("got %q; want %q", out, want)
	}
	if back, err := Read([]byte(out), openGaps); err != nil || !reflect.DeepEqual(back, rows) {
		t.Fatalf("read back %q, error %v; want %q", back, err, rows)
	}
	if out := write(t, openGaps); out != "" {
		t.Fatalf("no rows: got %q; want nothing", out)
	}
	if back, err := Read(nil, openGaps); err != nil || len(back) != 0 {
		t.Fatalf("an empty record: rows %q, error %v; want no rows", back, err)
	}
}

// Note 3: the writer makes the checks of the reader, and then writes nothing.
func TestWriteRefusesARowThatReadWouldRefuse(t *testing.T) {
	for _, c := range []struct {
		name         string
		s            Schema
		rows         [][]string
		line         int
		column, part string
	}{
		{"invalid UTF-8 (D5)", gaps, [][]string{{"Q-001", "G1", "1", "\xff\xfe fast", "q"}}, 2, "excerpt", "UTF-8"},
		{"a value of the wrong type", gaps, [][]string{{"Q-001", "G9", "1", "x", "q"}}, 2, "rule", "enum"},
		{"a tab that the field rule makes a space", gaps, [][]string{{"Q-001", "G1", "1\t", "x", "q"}}, 2, "line", "int"},
		{"too few fields", gaps, [][]string{{"Q-001", "G1"}}, 2, "line", "2 fields"},
		{"too many fields", gaps, [][]string{{"Q-001", "G1", "1", "x", "q", "y"}}, 2, "", "6 fields"},
		{"an empty key", gaps, [][]string{{"", "G1", "1", "x", "q"}}, 2, "id", "key"},
		{"— as a key (note 1)", gaps, [][]string{{"—", "G1", "1", "x", "q"}}, 2, "id", "key"},
		{"a key that repeats (D6)", stalls, [][]string{{"ST-001", "stall", "a"}, {"ST-001", "stall", "b"}}, 3, "", "line 2"},
		{"no header row: a key that repeats", openGaps, [][]string{{"a.md", "m", "q"}, {"a.md", "m", "q2"}}, 2, "", "line 1"},
		{"no header row: the column names as the first row", openGaps, [][]string{{"file", "marker", "question"}}, 1, "", "header row"},
		{"a byte-order mark at the start", openGaps, [][]string{{"\uFEFFa.md", "m", "q"}}, 1, "file", "byte-order mark"},
	} {
		t.Run(c.name, func(t *testing.T) {
			var b strings.Builder
			errorAt(t, Write(&b, c.s, c.rows), c.line, c.column, c.part)
			if b.Len() != 0 {
				t.Fatalf("wrote %q; want nothing", b.String())
			}
		})
	}
}

// D4, D5, note 1, note 9, condition 5: each refusal names the line, and the
// column of the field that holds the fault.
func TestReadRefusesARecordThatDoesNotMatch(t *testing.T) {
	row := "Q-001\tG1\t0\t—\tWhich stack?\n"
	for _, c := range []struct {
		name, input  string
		s            Schema
		line         int
		column, part string
	}{
		{"a carriage return before a line feed", gapsHeader + strings.TrimSuffix(row, "\n") + "\r\n", gaps, 2, "question", "line-feed endings"},
		{"a lone carriage return", gapsHeader + "Q-001\tG1\t3\ta\rb\tq\n", gaps, 2, "excerpt", "carriage return"},
		{"a carriage return in a field with no column", gapsHeader + "Q-001\tG1\t3\tx\tq\ta\rb\n", gaps, 2, "", "carriage return"},
		{"no header row: the column names as the first line", "file\tmarker\tquestion\ndocs/a.md\tm\tq\n", openGaps, 1, "", "header row"},
		{"a byte-order mark", "\uFEFF" + gapsHeader + row, gaps, 1, "", "byte-order mark"},
		{"invalid UTF-8", gapsHeader + "Q-001\tG1\t3\t\xff\xfe fast\tq\n", gaps, 2, "excerpt", "UTF-8"},
		{"an empty line", gapsHeader + "\n" + row, gaps, 2, "", "empty line"},
		{"an empty last line", gapsHeader + row + "\n", gaps, 3, "", "empty line"},
		{"no line feed after the last row", gapsHeader + strings.TrimSuffix(row, "\n"), gaps, 2, "", "no line feed"},
		{"no header row", "", gaps, 1, "", "header"},
		{"a header row with another name", "id\trules\tline\texcerpt\tquestion\n" + row, gaps, 1, "rule", "header"},
		{"a header row with a column missing", "id\trule\tline\texcerpt\n" + row, gaps, 1, "question", "header"},
		{"a header row with an extra field", "id\trule\tline\texcerpt\tquestion\tx\n" + row, gaps, 1, "", "6 fields"},
		{"too few fields", gapsHeader + "Q-001\tG1\t0\t—\n", gaps, 2, "question", "4 fields"},
		{"too many fields", gapsHeader + "Q-001\tG1\t0\t—\tq\textra\n", gaps, 2, "", "6 fields"},
		{"a field of the wrong type", gapsHeader + "Q-001\tG1\t07\t—\tq\n", gaps, 2, "line", "int"},
		{"an empty string as a field", gapsHeader + "Q-001\tG1\t0\t\tq\n", gaps, 2, "excerpt", "empty field"},
		{"— in a key column", gapsHeader + "—\tG1\t0\t—\tq\n", gaps, 2, "id", "key"},
		{"a key of one column that repeats", gapsHeader + row + row, gaps, 3, "", "line 2"},
		{"a key of two columns that repeats", "stall\tkind\tnote\nST-001\tstall\ta\nST-001\toutcome\tb\nST-001\tstall\tc\n", stalls, 4, "", "line 2"},
	} {
		t.Run(c.name, func(t *testing.T) {
			rows, err := Read([]byte(c.input), c.s)
			if rows != nil {
				t.Errorf("rows %q; want none", rows)
			}
			errorAt(t, err, c.line, c.column, c.part)
		})
	}
}

// Condition 1: the key is the tuple of all key columns; one of them may repeat.
func TestReadTakesAKeyOfTwoColumnsAsOneTuple(t *testing.T) {
	in := "stall\tkind\tnote\nST-001\tstall\ta\nST-001\tdiagnosis\tb\nST-001\toutcome\tc\nST-002\tstall\td\n"
	if rows, err := Read([]byte(in), stalls); err != nil || len(rows) != 4 {
		t.Fatalf("rows %q, error %v; want 4 rows", rows, err)
	}
}

// D1: — is valid in each column that is not a key column, whatever its type.
func TestReadGivesTheEmptyMarkAsTheEmptyValue(t *testing.T) {
	s := Schema{Name: "all", Location: "stdout", Columns: []Column{{Name: "id", Type: "id(P-NNN)", Key: true}}}
	header, row, want := "id", "P-001", []string{"P-001"}
	for _, typ := range []string{"text", "int", "decimal", "bool", "time", "sha1", "sha256", "path", "enum(a|b)", "id(<word>)", "list(path)"} {
		s.Columns = append(s.Columns, Column{Name: typ, Type: typ})
		header, row, want = header+"\t"+typ, row+"\t—", append(want, "")
	}
	rows, err := Read([]byte(header+"\n"+row+"\n"), s)
	if err != nil || !reflect.DeepEqual(rows, [][]string{want}) {
		t.Fatalf("rows %q, error %v; want %q", rows, err, want)
	}
}

// D7, note 4: a list value holds no space, and an empty list is —.
func TestAListValueHoldsNoSpaceAndAnEmptyListIsTheEmptyMark(t *testing.T) {
	for _, v := range []string{"docs/my file.md", "", "a\tb", "a\nb", "a\rb"} {
		if _, err := JoinList([]string{"go.mod", v}); err == nil {
			t.Errorf("JoinList with %q: no error; want one", v)
		}
	}
	none, _ := JoinList(nil)
	two, err := JoinList([]string{".golangci.yml", "docs/gates.tsv"})
	out := write(t, manifest, []string{"static", none}, []string{"lint", two})
	if want := "kind\tconfig\nstatic\t—\nlint\t.golangci.yml docs/gates.tsv\n"; err != nil || out != want {
		t.Fatalf("got %q, error %v; want %q", out, err, want)
	}
	rows, err := Read([]byte(out), manifest)
	if err != nil || len(strings.Fields(rows[0][1])) != 0 || !reflect.DeepEqual(strings.Fields(rows[1][1]), []string{".golangci.yml", "docs/gates.tsv"}) {
		t.Fatalf("rows %q, error %v; want no values, then two", rows, err)
	}
}

// NFR-005: the same rows give the same bytes, and a read gives the rows back.
func TestWriteThenReadGivesTheSameRowsAndBytes(t *testing.T) {
	rows := [][]string{
		{"Q-001", "G1", "0", "", "Which technology stack?"},
		{"Q-002", "G3", "14", "the PSB says so", `What does "PSB" mean?`},
		{"Q-1000", "G4", "200", "it is fast", "Which number?"},
	}
	first := write(t, gaps, rows...)
	back, err := Read([]byte(first), gaps)
	if err != nil || !reflect.DeepEqual(back, rows) {
		t.Fatalf("read back %q, error %v; want %q", back, err, rows)
	}
	if second, third := write(t, gaps, rows...), write(t, gaps, back...); second != first || third != first {
		t.Fatalf("the writes differ:\n%q\n%q\n%q", first, second, third)
	}
	// Note 1, the one lossy case: the text — is written as the empty value is.
	mark := write(t, gaps, []string{"Q-001", "G1", "0", "—", "q"})
	if empty := write(t, gaps, []string{"Q-001", "G1", "0", "", "q"}); mark != empty {
		t.Fatalf("the text — gives %q; the empty value gives %q", mark, empty)
	}
	if back, err := Read([]byte(mark), gaps); err != nil || back[0][3] != "" {
		t.Fatalf("the text — reads back as %q, error %v; want the empty value", back, err)
	}
}

// The error must be the error of the schema, which starts "schema x: ", and not
// an error of the header row or of the field count.
func TestWriteAndReadRefuseASchemaThatIsNotValid(t *testing.T) {
	for _, cols := range [][]Column{
		nil, {{Name: "a", Type: "float"}}, {{Name: "a", Type: "text"}, {Name: "a", Type: "int"}},
		{{Name: "", Type: "text"}}, {{Name: "a\tb", Type: "text"}}, {{Name: "a\xff", Type: "text"}},
	} {
		s := Schema{Name: "x", Location: "stdout", Columns: cols}
		if err := Write(&strings.Builder{}, s, nil); err == nil || !strings.HasPrefix(err.Error(), "schema x: ") {
			t.Errorf("Write with columns %v: error %v; want the error of the schema", cols, err)
		}
		if _, err := Read([]byte("a\n"), s); err == nil || !strings.HasPrefix(err.Error(), "schema x: ") {
			t.Errorf("Read with columns %v: error %v; want the error of the schema", cols, err)
		}
	}
}
