package records

import (
	"regexp"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// The block markers of milestone M2d (docs/spec/records.md, NFR-001 — The
// records of the Scaffold; task T-tetz, #270): the schema of
// intake/markers.tsv, which layup run writes at the step markers, and the rules
// that its block gives in words, read in the committed form. Each "— when" and
// "— only for" is read as "exactly when", as for the records of Intake.
//
// The file that the session reads, input/markers.tsv, has value, start, end
// and session as — (scaffold.md, The step markers); the reader here refuses
// its session, and that form is its writer's (row 53). What reads another file
// is the writer's too: the quote of a proposal in the problem statement, which
// gives start and end, and that the rows are the markers of the pinned tree.
//
// The id is read by its form, not recomputed. It is the key of S10: M- and the
// first 8 hexadecimal characters of the SHA-256 of the file, a tab and the
// marker, over the raw text (setup.MarkerID). The row holds the marker after
// the field rule, which makes a tab or a carriage return one space, so a reader
// that recomputed the hash would refuse the row of a marker that holds one. The
// writer computes the id from the scan.
//
// A table with no row is a scan that found no marker, and is read.

// MarkersSchema is the form of intake/markers.tsv: the block markers.
var MarkersSchema = tsv.Schema{Name: "markers", Location: "records:intake/markers.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(M-xxxxxxxx)", Key: true},
	{Name: "file", Type: "path"},
	{Name: "line", Type: "int"},
	{Name: "marker", Type: "text"},
	{Name: "value", Type: "text"},
	{Name: "start", Type: "int"},
	{Name: "end", Type: "int"},
	{Name: "session", Type: "text"},
}}

// The two angle quotes of a marker, U+2039 and U+203A, by their code points: a
// source file holds no marker, which the check markers of setup-check.sh reads.
const (
	markerOpen  = string(rune(0x2039))
	markerClose = string(rune(0x203a))
)

// markerIDForm is the form of the id: the type id(M-xxxxxxxx) takes any
// lowercase letter or digit, and the block's rule says hexadecimal.
var markerIDForm = regexp.MustCompile(`^M-[0-9a-f]{8}$`)

// markersRequired are the columns whose rule has no clause for the empty value,
// which never hold it (docs/spec/README.md, The types); the id is a key.
var markersRequired = []string{"file", "line", "marker", "session"}

// CheckMarker checks the rules of the block markers on one row: the form of the
// id, a line of 1 or more, the two angle quotes of the marker, value and start
// and end each the empty value exactly when the one before it is, start below
// end, and the session's form.
func CheckMarker(r []string) error {
	f := fields(MarkersSchema, r)
	if c := missing(f, markersRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	switch {
	case !markerIDForm.MatchString(f["id"]):
		return &RowError{"id", "an id is M- and 8 lowercase hexadecimal characters"}
	case !atLeastOne(f["line"]):
		return &RowError{"line", "a line is counted from 1"}
	case !strings.HasPrefix(f["marker"], markerOpen) || !strings.HasSuffix(f["marker"], markerClose):
		return &RowError{"marker", "a marker has its two angle quotes"}
	case (f["value"] == "") != (f["start"] == ""):
		return &RowError{"start", "start is the empty value exactly when value is"}
	case (f["start"] == "") != (f["end"] == ""):
		return &RowError{"end", "end is the empty value exactly when start is"}
	case f["start"] != "" && !below(f["start"], f["end"]):
		return &RowError{"end", "end is above start"}
	case !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	}
	return nil
}

// CheckMarkers checks the order of intake/markers.tsv, which reads more than
// one row: the rows are in the order of the stop table of S10, by file and then
// by line. The file of a row is not below the file before it, in the bytes of
// the path, and in one file the line is not below the line before it; two
// markers on one line are two rows with one line. A broken rule is a
// *tsv.Error.
func CheckMarkers(rows [][]string) error {
	for i := 1; i < len(rows); i++ {
		before, f := fields(MarkersSchema, rows[i-1]), fields(MarkersSchema, rows[i])
		switch {
		case f["file"] < before["file"]:
			return &tsv.Error{Line: i + 2, Column: "file", Reason: "the rows are in the order of the stop table of S10: the file is not below the file before it"}
		case f["file"] == before["file"] && below(f["line"], before["line"]):
			return &tsv.Error{Line: i + 2, Column: "line", Reason: "the rows are in the order of the stop table of S10: in a file the line is not below the line before it"}
		}
	}
	return nil
}

// ReadMarkers reads intake/markers.tsv by its schema, its row rules and the
// order of its rows.
func ReadMarkers(data []byte) ([][]string, error) {
	rows, err := read(data, MarkersSchema, CheckMarker)
	if err != nil {
		return nil, err
	}
	if err := CheckMarkers(rows); err != nil {
		return nil, err
	}
	return rows, nil
}
