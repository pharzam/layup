package records

import (
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// payload is the SHA-256 of a payload of a stall row or a diagnosis row.
var payload = strings.Repeat("ab", 32)

// stallRows gives three stalls of stalls.tsv (D2 and D3 of #95): ST-001 is
// closed, with a diagnosis and an outcome; ST-002 is a stall of the
// orchestrator, whose diagnosis failed; ST-003 is open, with its stall row
// only. The rows of ST-001 and ST-002 are mixed, as a writer appends each row
// when its step ends.
func stallRows() []map[string]string {
	return []map[string]string{
		{"stall": "ST-001", "kind": "stall", "time": "2026-10-03T01:00:00Z", "task": "T-dgy7", "trigger": "no-progress", "evidence": payload},
		{"stall": "ST-002", "kind": "stall", "time": "2026-10-03T01:05:00Z", "task": "project", "trigger": "orchestrator", "evidence": payload},
		{"stall": "ST-001", "kind": "diagnosis", "time": "2026-10-03T01:10:00Z", "task": "T-dgy7", "evidence": payload, "cause": "disagreement", "rung": "retry", "examiner": "S-0123abcd"},
		{"stall": "ST-002", "kind": "diagnosis-failed", "time": "2026-10-03T01:11:00Z", "task": "project", "note": "no admitted harness"},
		{"stall": "ST-001", "kind": "outcome", "time": "2026-10-03T01:30:00Z", "task": "T-dgy7", "outcome": "closed-without-human", "note": "the retry passed"},
		{"stall": "ST-002", "kind": "outcome", "time": "2026-10-03T02:00:00Z", "task": "project", "outcome": "closed-by-operator", "note": "the Operator restarted the run"},
		{"stall": "ST-003", "kind": "stall", "time": "2026-10-03T02:10:00Z", "task": "T-dgy7", "trigger": "hang", "evidence": payload},
	}
}

// with gives the rows of stallRows with the row i changed; "" is the empty
// value.
func with(i int, change map[string]string) []map[string]string {
	rows := stallRows()
	for k, v := range change {
		rows[i][k] = v
	}
	return rows
}

// A file of closed, open and orchestrator stalls passes (D2 and D3 of #95).
func TestStallRowsThatPass(t *testing.T) {
	rows, err := ReadStalls([]byte(fileOf(StallsSchema, stallRows()...)))
	if err != nil || len(rows) != 7 {
		t.Errorf("three stalls: %d rows, %v; want the seven rows", len(rows), err)
	}
	if rows, err := ReadStalls([]byte(fileOf(StallsSchema, stallRows()[:1]...))); err != nil || len(rows) != 1 {
		t.Errorf("one open stall: %d rows, %v", len(rows), err)
	}
}

// Each broken rule of a row or of the order is an error that names the line
// and the column (D2 and D3 of #95, with the conditions and the notes of its
// plan review).
func TestStallRowsThatBreakARule(t *testing.T) {
	for _, c := range []struct {
		name   string
		rows   []map[string]string
		line   int
		column string
	}{
		// The columns that never hold the empty value.
		{"no time", with(0, map[string]string{"time": ""}), 2, "time"},
		{"no task", with(0, map[string]string{"task": ""}), 2, "task"},
		// The columns of each kind.
		{"a stall row with no trigger", with(0, map[string]string{"trigger": ""}), 2, "trigger"},
		{"a stall row with no evidence", with(0, map[string]string{"evidence": ""}), 2, "evidence"},
		{"a diagnosis with no evidence", with(2, map[string]string{"evidence": ""}), 4, "evidence"},
		{"a diagnosis with no cause", with(2, map[string]string{"cause": ""}), 4, "cause"},
		{"a diagnosis with no rung", with(2, map[string]string{"rung": ""}), 4, "rung"},
		{"a diagnosis with no examiner", with(2, map[string]string{"examiner": ""}), 4, "examiner"},
		{"a failed diagnosis with no note", with(3, map[string]string{"note": ""}), 5, "note"},
		{"an outcome row with no outcome", with(4, map[string]string{"outcome": ""}), 6, "outcome"},
		{"an outcome row with no note", with(4, map[string]string{"note": ""}), 6, "note"},
		{"a trigger on a diagnosis", with(2, map[string]string{"trigger": "hang"}), 4, "trigger"},
		{"evidence on an outcome", with(4, map[string]string{"evidence": payload}), 6, "evidence"},
		{"a cause on a stall row", with(0, map[string]string{"cause": "other"}), 2, "cause"},
		{"a rung on an outcome", with(4, map[string]string{"rung": "panel"}), 6, "rung"},
		{"an examiner on a stall row", with(0, map[string]string{"examiner": "S-0123abcd"}), 2, "examiner"},
		{"an outcome on a stall row", with(0, map[string]string{"outcome": "task-stopped"}), 2, "outcome"},
		{"a note on a stall row", with(0, map[string]string{"note": "a reason"}), 2, "note"},
		{"a note on a diagnosis", with(2, map[string]string{"note": "a reason"}), 4, "note"},
		{"an examiner that is not a session ID", with(2, map[string]string{"examiner": "claude"}), 4, "examiner"},
		{"an examiner of a failed diagnosis that is not a session ID", with(3, map[string]string{"examiner": "S-XYZ"}), 5, "examiner"},
		// The order, the IDs, the time and the task of each stall.
		{"a diagnosis before its stall row", []map[string]string{stallRows()[2], stallRows()[0]}, 2, "stall"},
		{"both diagnosis rows of one stall", append(stallRows()[:3], map[string]string{"stall": "ST-001", "kind": "diagnosis-failed",
			"time": "2026-10-03T01:12:00Z", "task": "T-dgy7", "note": "a second diagnosis"}), 5, "kind"},
		{"an outcome before the second row", []map[string]string{stallRows()[0], stallRows()[4]}, 3, "kind"},
		{"ST-002 first", []map[string]string{stallRows()[1]}, 2, "stall"},
		{"ST-001 then ST-003", []map[string]string{stallRows()[0], stallRows()[6]}, 3, "stall"},
		{"a row before the time of the row before it", with(2, map[string]string{"time": "2026-10-03T00:59:00Z"}), 4, "time"},
		{"a row of a stall with another task", with(2, map[string]string{"task": "T-other"}), 4, "task"},
		{"an orchestrator stall of a task", with(1, map[string]string{"task": "T-dgy7"}), 3, "task"},
		{"a stall of project with another trigger", with(0, map[string]string{"task": "project"}), 2, "task"},
	} {
		_, err := ReadStalls([]byte(fileOf(StallsSchema, c.rows...)))
		var e *tsv.Error
		if !errors.As(err, &e) || e.Line != c.line || e.Column != c.column {
			t.Errorf("%s: %v; want an error of line %d, column %q", c.name, err, c.line, c.column)
		}
	}
}
