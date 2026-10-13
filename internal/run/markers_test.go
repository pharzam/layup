package run

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/setup"
)

// The unit tests of the step markers (docs/spec/scaffold.md, The step markers;
// the acceptance test "The step markers"; row 53 of the plan, task T-4c5c,
// #276): the list that the session reads, and the keep of the file that it
// writes. The angle quotes of a marker are written by their code points, as a
// source file holds none (check markers of setup-check.sh).

const (
	angleOpen  = string(rune(0x2039))
	angleClose = string(rune(0x203a))
)

func angled(s string) string { return angleOpen + s + angleClose }

// markersPSB is a problem statement of five lines: the heading, an empty line,
// and three lines that state a value, two of them with the same word.
const markersPSB = "# The brief\n\nThe owner is Carol Jones.\nThe port is 8080 on every host.\nA host is a host.\n"

// markerPlaces gives the places of three markers: one on a line of a file, one
// over two lines of another file, and one more on a later line of the first.
func markerPlaces() []setup.MarkerPlace {
	return setup.MarkerPlaces([]setup.Marker{
		{File: "docs/a.md", Line: 2, Col: 1, Text: angled("owner")},
		{File: "docs/a.md", Line: 7, Col: 1, Text: angled("port")},
		{File: "docs/b.md", Line: 1, Col: 1, Text: angled("two\nlines")},
	})
}

// psbAt gives the byte offsets of text in the problem statement.
func psbAt(text string) (int, int) {
	i := strings.Index(markersPSB, text)
	return i, i + len(text)
}

// The list that the session reads is the block markers with value, start, end
// and session as the empty value, one row per place, in the order of S10's stop
// table. It is the writer's form: the reader of the block refuses its session.
func TestTheListOfMarkers(t *testing.T) {
	places := markerPlaces()
	data, err := markersList(places)
	if err != nil {
		t.Fatal(err)
	}
	row := func(p setup.MarkerPlace) string {
		return fmt.Sprintf("%s\t%s\t%d\t%s\t—\t—\t—\t—\n", p.ID, p.File, p.Line, p.Marker)
	}
	want := "id\tfile\tline\tmarker\tvalue\tstart\tend\tsession\n" + row(places[0]) + row(places[1]) + row(places[2])
	if string(data) != want {
		t.Errorf("the list:\n%q\nwant:\n%q", data, want)
	}
	if got := places[2].Marker; got != angled("two lines") {
		t.Errorf("the key of a marker over two lines is %q; want its line break as a space", got)
	}
	if _, err := records.ReadMarkers(data); err == nil || !strings.Contains(err.Error(), "session") {
		t.Errorf("ReadMarkers of the list: %v; want it refused at the session, the list's form being the writer's", err)
	}
}

// The prompt names the problem statement, the list and the output file, each on
// a line of its own, and the three columns of the file.
func TestTheTextOfTheMarkersPrompt(t *testing.T) {
	p := string(markersText(map[string]string{"problem-statement.md": "/s/input/problem-statement.md", "markers.tsv": "/s/input/markers.tsv"}, "/s/result/markers.tsv"))
	for _, want := range []string{"\nThe problem statement: /s/input/problem-statement.md\n", "\nThe list of markers: /s/input/markers.tsv\n",
		"\nThe output file: /s/result/markers.tsv\n", `"id<TAB>value<TAB>quote"`} {
		if !strings.Contains(p, want) {
			t.Errorf("the prompt has no %q:\n%s", want, p)
		}
	}
	if strings.Contains(p, "{") {
		t.Errorf("the prompt keeps a placeholder:\n%s", p)
	}
}

// Keep reads the session's file by its form, finds each quote once, and gives
// intake/markers.tsv: one row per place, in the scan's order, with the value and
// the offsets of a proposal, or the empty value in the three columns; the file's
// own order does not matter.
func TestKeepOfMarkers(t *testing.T) {
	places := markerPlaces()
	keep := markersKeep([]byte(markersPSB), places)
	head := "id\tvalue\tquote\n"
	ownerQuote, portQuote := "The owner is Carol Jones.", "The port is 8080 on every host."
	os, oe := psbAt(ownerQuote)
	ps, pe := psbAt(portQuote)
	hs, he := psbAt("A host is a host.")
	row := func(p setup.MarkerPlace, value string, start, end int) string {
		s, e := "—", "—"
		if value != "" {
			s, e = fmt.Sprint(start), fmt.Sprint(end)
		} else {
			value = "—"
		}
		return fmt.Sprintf("%s\t%s\t%d\t%s\t%s\t%s\t%s\tS-1a2b3c4d\n", p.ID, p.File, p.Line, p.Marker, value, s, e)
	}
	table := "id\tfile\tline\tmarker\tvalue\tstart\tend\tsession\n"
	for _, c := range []struct{ name, file, want string }{
		{"two proposals", head + places[0].ID + "\tCarol Jones\t" + ownerQuote + "\n" + places[1].ID + "\t8080\t" + portQuote + "\n",
			table + row(places[0], "Carol Jones", os, oe) + row(places[1], "8080", ps, pe) + row(places[2], "", 0, 0)},
		{"the rows in the reverse order", head + places[1].ID + "\t8080\t" + portQuote + "\n" + places[0].ID + "\tCarol Jones\t" + ownerQuote + "\n",
			table + row(places[0], "Carol Jones", os, oe) + row(places[1], "8080", ps, pe) + row(places[2], "", 0, 0)},
		{"one proposal, for the marker over two lines", head + places[2].ID + "\ta host\tA host is a host.\n",
			table + row(places[0], "", 0, 0) + row(places[1], "", 0, 0) + row(places[2], "a host", hs, he)},
		{"no proposal", head, table + row(places[0], "", 0, 0) + row(places[1], "", 0, 0) + row(places[2], "", 0, 0)},
	} {
		files, err := keep("S-1a2b3c4d", []byte(c.file))
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got := string(files["intake/markers.tsv"]); got != c.want || len(files) != 1 {
			t.Errorf("%s: the records %q:\n%s\nwant:\n%s", c.name, keys(files), got, c.want)
		}
		if rows, err := records.ReadMarkers(files["intake/markers.tsv"]); err != nil || len(rows) != 3 {
			t.Errorf("%s: ReadMarkers of the record: %d rows, %v", c.name, len(rows), err)
		}
	}
}

// Each refusal of the session's file, alone: a key that the list lacks, a key
// twice, an empty value or quote, a quote that occurs nowhere, twice or holds a
// tab, and a file of the whole block. A refused file gives no record.
func TestKeepOfMarkersRefusesEachBrokenRule(t *testing.T) {
	places := markerPlaces()
	keep := markersKeep([]byte(markersPSB), places)
	head := "id\tvalue\tquote\n"
	id := places[0].ID
	for _, c := range []struct{ name, file, refused string }{
		{"a key that the list lacks", head + "M-00000000\tx\tThe owner is Carol Jones.\n", "row 1: the key M-00000000 is not in the list"},
		{"a key twice", head + id + "\tx\tThe owner is Carol Jones.\n" + id + "\ty\tThe port is 8080 on every host.\n", "repeats the key of line 2"},
		{"an empty value", head + id + "\t—\tThe owner is Carol Jones.\n", "row 1: the value is empty"},
		{"an empty quote", head + id + "\tx\t—\n", "the quote is empty"},
		{"a quote that occurs nowhere", head + id + "\tx\tThe owner is Dave.\n", "occurs nowhere"},
		{"a quote that occurs twice", head + id + "\tx\thost\n", "occurs 3 times"},
		{"a quote that holds a tab", head + id + "\tx\tThe\towner\n", "line 2"},
		{"a file of the whole block", "id\tfile\tline\tmarker\tvalue\tstart\tend\tsession\n", "header"},
	} {
		files, err := keep("S-1a2b3c4d", []byte(c.file))
		if err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
}

// The table of the record is read by the reader of the block before it is kept: a
// place whose marker lacks its angle quotes (which the scan never gives) refuses
// the record as a whole.
func TestKeepOfMarkersReadsTheTableByItsBlock(t *testing.T) {
	bad := []setup.MarkerPlace{{ID: markerPlaces()[0].ID, File: "docs/a.md", Marker: "no quotes", Line: 2}}
	files, err := markersKeep([]byte(markersPSB), bad)("S-1a2b3c4d", []byte("id\tvalue\tquote\n"))
	if err == nil || files != nil || !strings.Contains(err.Error(), "angle quotes") {
		t.Errorf("a place with no angle quotes: %q, %v; want the record refused by the block's rule", files, err)
	}
}
