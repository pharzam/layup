package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// The Keep of the completeness review (docs/spec/intake.md, The completeness
// review): the session writes the block coverage without id and session; each
// fact is a fact of the span table and each requirement an ID of the draft.
func TestKeepOfCompleteness(t *testing.T) {
	spans, err := records.ReadSpans([]byte(draftSpans))
	if err != nil {
		t.Fatal(err)
	}
	draft, err := records.ReadDraft([]byte("id\tkind\tstatement\tcovers\tcriterion\tsession\n" +
		"R-001\trequirement\tD.\t1\tD.\tS-1a2b3c4d\nR-002\tnon-functional\tK.\t2\tK.\tS-1a2b3c4d\n"))
	if err != nil {
		t.Fatal(err)
	}
	keep := completenessKeep(spans, draft)
	head := "kind\tfact\trequirement\tquestion\n"
	good := head + "uncovered\t3\t—\tWhat measures 3?\nunsupported\t—\tR-002\tWhy R-002?\nclass\t4\t—\tIs 4 a need?\n"
	files, err := keep("S-0a0b0c0d", []byte(good))
	if err != nil {
		t.Fatalf("Keep: %v", err)
	}
	want := "id\tkind\tfact\trequirement\tquestion\tsession\n" +
		"C-001\tuncovered\t3\t—\tWhat measures 3?\tS-0a0b0c0d\n" +
		"C-002\tunsupported\t—\tR-002\tWhy R-002?\tS-0a0b0c0d\n" +
		"C-003\tclass\t4\t—\tIs 4 a need?\tS-0a0b0c0d\n"
	if got := string(files["intake/coverage.tsv"]); got != want || len(files) != 1 {
		t.Errorf("the records %q:\n%s\nwant:\n%s", keys(files), got, want)
	}
	if files, err := keep("S-0a0b0c0d", []byte(head)); err != nil || string(files["intake/coverage.tsv"]) != "id\tkind\tfact\trequirement\tquestion\tsession\n" {
		t.Errorf("a file of no item: %q, %v; want the header alone", files, err)
	}
	for _, c := range []struct{ name, file, refused string }{
		{"a fact that is no fact", head + "uncovered\t9\t—\tQ?\n", "row 1: the fact 9 is no fact of spec/spans.tsv"},
		{"a requirement of no draft row", head + "unsupported\t—\tR-007\tQ?\n", "row 1: the requirement R-007 is no row of spec/draft.tsv"},
		{"an unsupported item with a fact", head + "unsupported\t1\tR-001\tQ?\n", "fact"},
		{"a file of the record's form", "id\t" + head + "C-001\tclass\t4\t—\tQ?\n", "header"},
	} {
		if files, err := keep("S-0a0b0c0d", []byte(c.file)); err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
}
