package psb

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"strings"
	"testing"
)

// testdata holds the goldens, so the unit tests read no file at run time.
//
//go:embed testdata/*.md testdata/*.tsv
var testdata embed.FS

// golden runs Check on input and compares the table, byte for byte, with want.
// One mechanism, one test, every rule (issue #37).
func golden(t *testing.T, name string, input, want []byte) {
	t.Helper()
	var got bytes.Buffer
	if err := WriteTSV(&got, Check(string(input))); err != nil {
		t.Fatal(err)
	}
	if got.String() != string(want) {
		t.Fatalf("%s: the table differs\n--- got ---\n%s--- want ---\n%s", name, got.String(), want)
	}
}

func read(t *testing.T, name string) []byte {
	t.Helper()
	b, err := testdata.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestGolden(t *testing.T) {
	for _, name := range []string{"triggers", "clean", "edge", "values"} {
		t.Run(name, func(t *testing.T) { golden(t, name, read(t, name+".md"), read(t, name+".tsv")) })
	}
}

// A line feed after a carriage return is a line end; the table is the same.
func TestCRLFGivesTheSameTable(t *testing.T) {
	golden(t, "triggers with CRLF", []byte(strings.ReplaceAll(string(read(t, "triggers.md")), "\n", "\r\n")), read(t, "triggers.tsv"))
}

const header = "id\trule\tline\texcerpt\tquestion\n"

// A lone carriage return inside a line becomes one space in the excerpt, by
// the field rule of the records. A Go literal keeps the raw byte out of a
// tracked file.
func TestALoneCarriageReturnBecomesASpace(t *testing.T) {
	golden(t, "a lone CR", []byte("Technology stack: Go\nThe API\ris fast.\n"),
		[]byte(header+"Q-001\tG4\t2\tThe API is fast.\tWhich number or threshold does \"fast\" stand for here?\n"))
}

const g1 = "Q-001\tG1\t0\t—\tWhich technology stack does the product use (languages, frameworks, tools)?\n"

// O-131: the value of a named stack is a character that is neither white
// space (Unicode White_Space) nor *, so an empty label of a template is a G1
// gap. After the colon, white space other than a tab, a form feed or a
// carriage return counts as a space (review round 1 of #83).
func TestG1ReadsTheValueOfAStack(t *testing.T) {
	for _, c := range []struct{ input, want string }{
		{"**Technology stack:**\n", header + g1},
		{"Technology stack: *\n", header + g1},
		{"Technology stack:\n", header + g1},
		{"Technology stack\t: Go\n", header + g1},
		{"Technology stack: \u00a0\n", header + g1},
		{"Technology stack: \u2003\n", header + g1},
		{"Technology stack: \v\n", header + g1},
		{"Technology stack:\tGo\n", header + g1},
		{"**Technology stack:** Go\n", header},
		{"TECHNOLOGY STACK: Go\n", header},
		{"Technology stack: *Go*\n", header},
		{"Technology stack:\u00a0Go\n", header},
		{"Technology stack:\vGo\n", header},
	} {
		golden(t, fmt.Sprintf("%q", c.input), []byte(c.input), []byte(c.want))
	}
}

// The edge cases that the rule table does not settle and that need a file of
// their own, with the present code as their reference (docs/spec/psb-check.md,
// The rules): G3 reads the terms table of the whole file. The other edge cases
// are in testdata/edge.md.
func TestEdgeCases(t *testing.T) {
	const stack = "Technology stack: Go\n"
	for _, c := range []struct{ name, input, want string }{
		{"G3: a heading terms in lower case is not a terms heading", stack + "## terms\n| **API** | x |\nThe SLA holds.\n", header},
		{"G3: the first heading that holds Terms makes the terms table", stack + "## Payment Terms\n| **SLA** | x |\n\n## Terms\n| **API** | y |\nThe API holds.\n",
			header + "Q-001\tG3\t6\t| **API** | y |\tWhat does \"API\" mean? The terms table does not define it.\n"},
		{"G3: a byte-order mark before a first-line # Terms stops the heading", "\uFEFF# Terms\n| **API** | x |\n\n" + stack + "The SLA holds.\n", header},
		{"G3: an abbreviation inside the terms table", stack + "## Terms\n| **API** | the SLA |\n",
			header + "Q-001\tG3\t3\t| **API** | the SLA |\tWhat does \"SLA\" mean? The terms table does not define it.\n"},
		{"a # line inside a code fence is a heading, and G4 applies there", stack + "```\n# Terms\n| **API** | x |\nthe cache is fast\n```\nThe SLA holds.\n",
			header + "Q-001\tG4\t5\tthe cache is fast\tWhich number or threshold does \"fast\" stand for here?\n" +
				"Q-002\tG3\t7\tThe SLA holds.\tWhat does \"SLA\" mean? The terms table does not define it.\n"},
	} {
		golden(t, c.name, []byte(c.input), []byte(c.want))
	}
}

// failingWriter is an output that refuses each write, as a closed pipe does.
type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("broken pipe") }

func TestWriteTSVGivesTheErrorOfItsOutput(t *testing.T) {
	if err := WriteTSV(failingWriter{}, Check("No stack.\n")); err == nil {
		t.Fatal("WriteTSV to an output that refuses each write: no error")
	}
}
