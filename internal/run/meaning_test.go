package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// The quote check of the review of meaning (docs/spec/intake.md, The review of
// meaning; Input states): a quote is text of one line of the problem
// statement, byte for byte, found once; its offsets are code's.
func TestQuoteAt(t *testing.T) {
	psb := []byte("# The problem\n\nThe service is fast.\nIt is fast and safe.\n")
	for _, c := range []struct {
		quote      string
		start, end int
		refused    string
	}{
		{"# The problem", 0, 13, ""},
		{"service is", 19, 29, ""},
		{"safe.", 51, 56, ""},
		{"The service is fast.", 15, 35, ""},
		{"slow", 0, 0, "occurs nowhere"},
		{"fast", 0, 0, "occurs 2 times"},
		{"s", 0, 0, "occurs 6 times"},
		{"is\tfast", 0, 0, "holds a tab"},
		{"fast.\nIt", 0, 0, "holds a line break"},
		{"problem\r", 0, 0, "holds a line break"},
		{"", 0, 0, "is empty"},
	} {
		start, end, err := quoteAt(psb, c.quote)
		switch {
		case c.refused == "" && (err != nil || start != c.start || end != c.end || string(psb[start:end]) != c.quote):
			t.Errorf("%q: %d, %d, %v; want %d, %d", c.quote, start, end, err, c.start, c.end)
		case c.refused != "" && (err == nil || !strings.Contains(err.Error(), c.refused)):
			t.Errorf("%q: %v; want it refused, %q", c.quote, err, c.refused)
		}
	}
	// Occurrences that overlap are each an occurrence: "--" is twice in
	// "---", so it names no one place (round 1 of #218).
	for _, c := range []struct{ psb, quote, refused string }{
		{"A rule:\n---\n", "--", "occurs 2 times"},
		{"....\n", "...", "occurs 2 times"},
		{"it is is is so\n", "is is", "occurs 2 times"},
	} {
		if _, _, err := quoteAt([]byte(c.psb), c.quote); err == nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%q in %q: %v; want it refused, %q", c.quote, c.psb, err, c.refused)
		}
	}
}

// The file of a session of the review of meaning is the block meaning without
// the columns that code adds; Keep gives intake/meaning.tsv with the IDs, the
// offsets and the session, or refuses the file.
func TestKeepOfMeaning(t *testing.T) {
	psb := []byte("# The problem\n\nThe service is fast.\nIt is fast and safe.\n")
	keep := meaningKeep(psb)
	head := "kind\tquote\tquestion\n"
	files, err := keep("S-1a2b3c4d", []byte(head+"unmeasurable\tThe service is fast.\tHow fast?\ntwo-readings\tsafe.\tSafe from what?\n"))
	if err != nil {
		t.Fatalf("Keep: %v", err)
	}
	want := "id\tkind\tstart\tend\tquote\tquestion\tsession\n" +
		"M-001\tunmeasurable\t15\t35\tThe service is fast.\tHow fast?\tS-1a2b3c4d\n" +
		"M-002\ttwo-readings\t51\t56\tsafe.\tSafe from what?\tS-1a2b3c4d\n"
	if got := string(files["intake/meaning.tsv"]); got != want || len(files) != 1 {
		t.Errorf("the records %q:\n%s\nwant:\n%s", keys(files), got, want)
	}
	if rows, err := records.ReadMeaning(files["intake/meaning.tsv"]); err != nil || len(rows) != 2 {
		t.Errorf("ReadMeaning of the record: %d rows, %v", len(rows), err)
	}
	if files, err := keep("S-1a2b3c4d", []byte(head)); err != nil || string(files["intake/meaning.tsv"]) != "id\tkind\tstart\tend\tquote\tquestion\tsession\n" {
		t.Errorf("a file of no gap: %q, %v; want the header alone", files, err)
	}
	for _, c := range []struct{ name, file, refused string }{
		{"a quote that occurs twice", head + "unmeasurable\tfast\tHow fast?\n", `the quote "fast" occurs 2 times`},
		{"a quote that occurs nowhere", head + "fit\tslow\tWhy?\n", `the quote "slow" occurs nowhere`},
		{"a kind of no list", head + "vague\tsafe.\tWhy?\n", "kind"},
		{"a file of the whole block", "id\tkind\tstart\tend\tquote\tquestion\tsession\n", "header"},
	} {
		if files, err := keep("S-1a2b3c4d", []byte(c.file)); err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
}

func keys(m map[string][]byte) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}

// The schema of a session's file is its block without the columns that code
// adds, in the block's order; a key column that is kept stays a key.
func TestSessionSchema(t *testing.T) {
	for _, c := range []struct {
		block      tsv.Schema
		drop       []string
		want, keys string
	}{
		{records.MeaningSchema, []string{"id", "start", "end", "session"}, "kind quote question", ""},
		{records.SpansSchema, []string{"text"}, "start end fact class reason", "start"},
	} {
		var names, keys []string
		for _, col := range sessionSchema(c.block, c.drop...).Columns {
			names = append(names, col.Name)
			if col.Key {
				keys = append(keys, col.Name)
			}
		}
		if strings.Join(names, " ") != c.want || strings.Join(keys, " ") != c.keys {
			t.Errorf("%s: the columns %q, the keys %q; want %q and %q", c.block.Name, names, keys, c.want, c.keys)
		}
	}
	if len(records.MeaningSchema.Columns) != 7 {
		t.Error("the block meaning changed")
	}
}
