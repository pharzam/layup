package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// draftSpans is a span table of one heading and four facts, one of each class.
const draftSpans = "start\tend\tfact\tclass\treason\ttext\n" +
	"0\t3\t—\t—\ta heading\t# A\n" +
	"4\t10\t1\tneed\t—\tneed 1\n" +
	"11\t17\t2\tconstraint\t—\tcons 2\n" +
	"18\t24\t3\tmeasure\t—\tmeas 3\n" +
	"25\t31\t4\tcontext\t—\tcont 4\n"

// The check of the draft (docs/spec/intake.md, The draft): each covers number
// is a fact of spec/spans.tsv, each need and constraint fact is covered, and
// each measure fact has a success criterion; Keep gives spec/draft.tsv with
// R-NNN and the session, or refuses the file.
func TestKeepOfDraft(t *testing.T) {
	spans, err := records.ReadSpans([]byte(draftSpans))
	if err != nil {
		t.Fatal(err)
	}
	keep := draftKeep(spans)
	head := "kind\tstatement\tcovers\tcriterion\n"
	good := head + "requirement\tThe tool does 1.\t1\tIt does 1.\n" +
		"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.\n" +
		"success-criterion\tThe tool meets 3.\t3 1\tIt meets 3.\n"
	files, err := keep("S-1a2b3c4d", []byte(good))
	if err != nil {
		t.Fatalf("Keep: %v", err)
	}
	want := "id\tkind\tstatement\tcovers\tcriterion\tsession\n" +
		"R-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS-1a2b3c4d\n" +
		"R-002\tnon-functional\tThe tool keeps 2.\t2\tIt keeps 2.\tS-1a2b3c4d\n" +
		"R-003\tsuccess-criterion\tThe tool meets 3.\t3 1\tIt meets 3.\tS-1a2b3c4d\n"
	if got := string(files["spec/draft.tsv"]); got != want || len(files) != 1 {
		t.Errorf("the records %q:\n%s\nwant:\n%s", keys(files), got, want)
	}
	if rows, err := records.ReadDraft(files["spec/draft.tsv"]); err != nil || len(rows) != 3 {
		t.Errorf("ReadDraft of the record: %d rows, %v", len(rows), err)
	}
	for _, c := range []struct{ name, file, refused string }{
		{"a row that covers nothing", good + "requirement\tMore.\t—\tMore.\n", "row 4: covers names no fact"},
		{"a number above the facts", good + "requirement\tMore.\t5\tMore.\n", "row 4: covers names 5, which is no fact"},
		{"a need covered by no row", head + "non-functional\tK.\t2\tK.\nsuccess-criterion\tM.\t3\tM.\n", "the need fact 1 is covered by no row"},
		{"a constraint covered by no row", head + "requirement\tD.\t1\tD.\nsuccess-criterion\tM.\t3\tM.\n", "the constraint fact 2 is covered by no row"},
		{"a measure with no success criterion", head + "requirement\tD.\t1 3\tD.\nnon-functional\tK.\t2\tK.\n", "the measure fact 3 has no success criterion"},
		{"a file with the column id", "id\t" + head + "R-001\trequirement\tD.\t1\tD.\n", "header"},
	} {
		if files, err := keep("S-1a2b3c4d", []byte(c.file)); err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
}
