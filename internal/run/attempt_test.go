package run

import (
	"context"
	"errors"
	"testing"

	"github.com/pharzam/layup/internal/gate"
)

func tableOf(rows ...gate.Row) gate.Table { return gate.Table{Base: aSHA, Head: aSHA, Rows: rows} }

// loop.md, The attempts: a failure of a check goes back to the developer as a
// finding. The rows of the gate's table that did not pass are the findings:
// fail and not-active (NFR-004: never a success), not pass or clear.
func TestAGateFailureIsTheRowsThatDidNotPass(t *testing.T) {
	got := GateFindings(tableOf(
		gate.Row{Kind: "build", State: "active", Result: "pass", Reason: "—"},
		gate.Row{Kind: "static", State: "active", Result: "fail", Reason: "command failed"},
		gate.Row{Kind: "boundary", State: "pending", Result: "clear", Reason: "no product path changed"},
		gate.Row{Kind: "tests", State: "active", Result: "not-active", Reason: "scratch tree: add failed"},
	))
	want := []Finding{{Rule: "gate:static", Text: "command failed"}, {Rule: "gate:tests", Text: "scratch tree: add failed"}}
	if len(got) != len(want) {
		t.Fatalf("GateFindings = %+v, want %+v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("finding %d = %+v, want %+v", i, got[i], want[i])
		}
	}
	if len(GateFindings(tableOf(gate.Row{Kind: "build", State: "active", Result: "pass", Reason: "—"}))) != 0 {
		t.Error("a table that passed has findings")
	}
	if len(GateFindings(gate.Table{})) != 0 {
		t.Error("a table with no row has findings")
	}
}

func TestNoFailureStartsNoAttempt(t *testing.T) {
	w := newHandoffWorld(t)
	_, err := w.s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: aSHA, Table: tableOf(gate.Row{Kind: "build", State: "active", Result: "pass", Reason: "—"})})
	if !errors.Is(err, ErrNoFailure) {
		t.Errorf("NextAttempt with no failed row: %v; want ErrNoFailure", err)
	}
	if len(w.store.commits) != 0 {
		t.Errorf("%d records commits for a table that passed; want none", len(w.store.commits))
	}
}

func TestTheNextAttemptIsTheNextNumberAtTheBaseWithTheFindings(t *testing.T) {
	st := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	s := newSessions(st)
	failed := tableOf(gate.Row{Kind: "static", State: "active", Result: "fail", Reason: "command failed"})
	for want := 1; want <= 2; want++ {
		a, err := s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: aSHA, Table: failed})
		if err != nil || a.N != want || len(a.Findings) != 1 || a.Findings[0].Rule != "gate:static" {
			t.Fatalf("NextAttempt = %+v, %v; want attempt %d with the finding of the gate", a, err, want)
		}
	}
	rows := events(t, st, "T-ab12")
	if len(rows) != 2 || rows[1][1] != "attempt" || rows[1][2] != "2" || rows[1][4] != aSHA {
		t.Errorf("the events %v; want an event attempt 2 with the base %s", rows, aSHA)
	}
	// A base that is no commit ID is refused by the records' own check, and nothing is committed.
	n := len(st.commits)
	if _, err := s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: "main", Table: failed}); err == nil || len(st.commits) != n {
		t.Errorf("a base of %q: %v, %d commits more", "main", err, len(st.commits)-n)
	}
}
