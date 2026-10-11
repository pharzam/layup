package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// The Keep of the step numbering (docs/spec/intake.md, The numbering; The span
// check): the session writes the block spans without text, and code adds each
// span's text from the file, with each tab, CR and LF as one space.
func TestKeepOfNumbering(t *testing.T) {
	psb := []byte("# A\n\nThe b\tc.\n")
	keep := numberingKeep(psb)
	head := "start\tend\tfact\tclass\treason\n"
	files, err := keep("S-1a2b3c4d", []byte(head+"0\t3\t—\t—\ta heading\n5\t14\t1\tneed\t—\n"))
	if err != nil {
		t.Fatalf("Keep: %v", err)
	}
	want := "start\tend\tfact\tclass\treason\ttext\n0\t3\t—\t—\ta heading\t# A\n5\t14\t1\tneed\t—\tThe b c. \n"
	if got := string(files["spec/spans.tsv"]); got != want || len(files) != 1 {
		t.Errorf("the records %q:\n%q\nwant:\n%q", keys(files), got, want)
	}
	if rows, err := records.ReadSpans(files["spec/spans.tsv"]); err != nil || len(rows) != 2 {
		t.Errorf("ReadSpans of the record: %d rows, %v", len(rows), err)
	}
	// A span partly out of the file keeps the text of its part in the file;
	// the rule Bounds of the check reports it.
	files, err = keep("S-1a2b3c4d", []byte(head+"0\t3\t—\t—\ta heading\n5\t40\t1\tneed\t—\n"))
	if err != nil || !strings.HasSuffix(string(files["spec/spans.tsv"]), "\n5\t40\t1\tneed\t—\tThe b c. \n") {
		t.Errorf("a span partly out of the file: %q, %v; want its text in the file", files["spec/spans.tsv"], err)
	}
	for _, c := range []struct{ name, file, refused string }{
		{"a span with no byte in the file", head + "20\t25\t1\tneed\t—\n", "text"},
		{"a start that repeats", head + "0\t3\t—\t—\ta heading\n0\t14\t1\tneed\t—\n", "start"},
		{"a class with no fact", head + "0\t3\t—\tneed\ta heading\n", "class"},
		{"a file with the column text", head[:len(head)-1] + "\ttext\n0\t3\t—\t—\ta heading\t# A\n", "header"},
	} {
		if files, err := keep("S-1a2b3c4d", []byte(c.file)); err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
}
