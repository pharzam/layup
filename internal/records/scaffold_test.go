package records

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The block markers of the Scaffold (task T-tetz, #270): one valid file, one
// case per rule that its block gives in words, one case per column whose rule
// has no clause for the empty value, and a wrong header. Each refusal must name
// the line, the column and the reason of its rule (refusedBy), because the form
// of a column refuses its empty value too, so a case that another rule refuses
// at the same column fails no mutation of its own. A marker is written with the
// code points of its two angle quotes, never the characters, so the check
// markers of setup-check.sh reads no marker in this file.

const (
	angleOpen     = string(rune(0x2039))
	angleClose    = string(rune(0x203a))
	markersHeader = "id\tfile\tline\tmarker\tvalue\tstart\tend\tsession"
)

// angled gives the text of a marker: the text between its two angle quotes.
func angled(text string) string { return angleOpen + text + angleClose }

// Three markers in the order of the stop table of S10: a marker with a
// proposal, a marker with none, and a marker of the next file.
var (
	markerA = "M-1a2b3c4d\tdocs/a.md\t3\t" + angled("the name") + "\tLAYUP\t10\t15\tS-0a1b2c3d"
	markerB = "M-5e6f7a8b\tdocs/a.md\t9\t" + angled("the owner") + "\t—\t—\t—\tS-0a1b2c3d"
	markerC = "M-9c0d1e2f\tdocs/b.md\t1\t" + angled("two words") + "\t—\t—\t—\tS-0a1b2c3d"
)

// refusedBy reports a test error unless ReadMarkers fails on a *tsv.Error of
// line and column whose reason holds the text reason.
func refusedBy(t *testing.T, name string, data []byte, line int, column, reason string) {
	t.Helper()
	_, err := ReadMarkers(data)
	var te *tsv.Error
	switch {
	case err == nil:
		t.Errorf("%s: read, want an error of line %d, column %q", name, line, column)
	case !errors.As(err, &te):
		t.Errorf("%s: %v, want a *tsv.Error of line %d, column %q", name, err, line, column)
	case te.Line != line || te.Column != column || !strings.Contains(te.Reason, reason):
		t.Errorf("%s: %v, want line %d, column %q and a reason with %q", name, err, line, column, reason)
	}
}

func TestValidMarkersAreRead(t *testing.T) {
	sameLine := field(markersHeader, strings.Replace(markerB, "M-5e6f7a8b", "M-0f1e2d3c", 1), "line", "3")
	// The id is the key of S10 over the raw text; the field rule writes the tab
	// of a marker as one space, so the reader checks the form of the id, and
	// does not compute the hash again.
	sum := sha256.Sum256([]byte("docs/a.md\t" + angled("one\ttwo")))
	tabbed := field(markersHeader, field(markersHeader, markerA, "id", fmt.Sprintf("M-%x", sum[:4])), "marker", angled("one two"))
	for name, data := range map[string][]byte{
		"a marker with a proposal, one with none, and one of the next file": table(markersHeader, markerA, markerB, markerC),
		"a scan that found no marker":                                       table(markersHeader),
		"two markers on one line":                                           table(markersHeader, markerA, sameLine),
		"a marker that held a tab, under the id of its raw text":            table(markersHeader, tabbed),
	} {
		if _, err := ReadMarkers(data); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestMarkersRefuseEachBrokenRule(t *testing.T) {
	const never = "never holds the empty value"
	row := func(change ...string) []byte { // markerA with the changes, a column and a value each
		r := markerA
		for i := 0; i < len(change); i += 2 {
			r = field(markersHeader, r, change[i], change[i+1])
		}
		return table(markersHeader, r)
	}
	refusedBy(t, "the empty value in id", row("id", "—"), 2, "id", "key")
	for _, col := range []string{"file", "line", "marker", "session"} {
		refusedBy(t, "the empty value in "+col, row(col, "—"), 2, col, never)
	}
	refusedBy(t, "an id with a letter that is not hexadecimal", row("id", "M-1a2b3c4z"), 2, "id", "8 lowercase hexadecimal")
	refusedBy(t, "a line 0", row("line", "0"), 2, "line", "counted from 1")
	refusedBy(t, "a marker with no open quote", row("marker", "the name"+angleClose), 2, "marker", "two angle quotes")
	refusedBy(t, "a marker with no close quote", row("marker", angleOpen+"the name"), 2, "marker", "two angle quotes")
	refusedBy(t, "a value with no start and no end", row("start", "—", "end", "—"), 2, "start", "exactly when value is")
	refusedBy(t, "a start and an end with no value", row("value", "—"), 2, "start", "exactly when value is")
	refusedBy(t, "an end with no start", table(markersHeader, field(markersHeader, markerB, "end", "15")), 2, "end", "exactly when start is")
	refusedBy(t, "start equal to end", row("end", "10"), 2, "end", "above start")
	refusedBy(t, "start above end", row("end", "9"), 2, "end", "above start")
	refusedBy(t, "a session that is not S- and 8 hexadecimal characters", row("session", "S-0a1b2c3z"), 2, "session", "S- and 8 lowercase hexadecimal")
	refusedBy(t, "a file that goes down", table(markersHeader, markerC, markerA), 3, "file", "not below the file")
	refusedBy(t, "a line that goes down in a file", table(markersHeader, markerB, markerA), 3, "line", "not below the line")
	refusedBy(t, "the second row broken", table(markersHeader, markerA, field(markersHeader, markerC, "line", "0")), 3, "line", "counted from 1")
}

func TestMarkersRefuseAWrongHeader(t *testing.T) {
	if _, err := ReadMarkers(table(strings.Replace(markersHeader, "\t", "\tx", 1), markerA)); err == nil {
		t.Error("markers: a wrong header: read, want an error")
	}
}
