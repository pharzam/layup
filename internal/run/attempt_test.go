package run

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/rules"
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

// The findings of layup/spec and layup/rules (task T-tnxn, #427, row 66aa; loop.md, The attempts).

func specTableRows(results map[int]string) [][]string {
	var rows [][]string
	for i, rule := range specRuleList() {
		result, detail := "pass", "—"
		if r, ok := results[i]; ok {
			result, detail = r, "detail of "+rule
		}
		rows = append(rows, []string{rule, result, detail})
	}
	return rows
}

// U1: the fail rows, in the order of the rules, the row ran included.
func TestTheSpecFindingsAreTheFailRowsInTheOrderOfTheRules(t *testing.T) {
	list := specRuleList()
	got := SpecFindings(specTableRows(map[int]string{1: "fail", 4: "clear", 9: "fail"}))
	want := []Finding{{Rule: "spec:" + list[1], Text: "detail of " + list[1]}, {Rule: "spec:" + list[9], Text: "detail of " + list[9]}}
	if !reflect.DeepEqual(got, want) || list[9] != "ran" {
		t.Errorf("SpecFindings = %+v; want %+v (the rule ran last)", got, want)
	}
}

// U2: the finding of the basic rule is its reason with its path, the text of the status.
func TestTheRulesFindingIsItsReasonWithItsPath(t *testing.T) {
	got := RulesFindings(rules.Result{State: rules.Failure, Reason: rules.ReasonRulePath, Path: ".github/workflows/ci.yml"})
	want := []Finding{{Rule: "rules:rule path", File: ".github/workflows/ci.yml", Text: "rule path: .github/workflows/ci.yml"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RulesFindings = %+v; want %+v", got, want)
	}
}

// U3, a green negative at the stub (its red controls are U1 and U2): a pass, a clear row, a success and an error give no finding.
func TestAPassClearSuccessOrErrorGivesNoFinding(t *testing.T) {
	if got := SpecFindings(specTableRows(map[int]string{0: "clear", 3: "clear"})); len(got) != 0 {
		t.Errorf("SpecFindings of pass and clear rows = %+v; want none", got)
	}
	for _, r := range []rules.Result{{State: rules.Success, Reason: rules.ReasonNoRulePath}, {State: rules.Error, Reason: rules.ReasonRegister}, {}} {
		if got := RulesFindings(r); len(got) != 0 {
			t.Errorf("RulesFindings(%+v) = %+v; want none", r, got)
		}
	}
}

// U4: the gate's findings, then the given ones; the given ones alone start an attempt; neither is ErrNoFailure.
func TestTheNextAttemptTakesTheGateFindingsThenTheGivenOnes(t *testing.T) {
	st := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	s := newSessions(st)
	given := []Finding{{Rule: "spec:criteria", Text: "x"}, {Rule: "rules:rule path", File: "a", Text: "y"}}
	failed := tableOf(gate.Row{Kind: "exits", State: "active", Result: "fail", Reason: "exit 3"})
	a, err := s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: aSHA, Table: failed, Findings: given})
	want := append([]Finding{{Rule: "gate:exits", Text: "exit 3"}}, given...)
	if err != nil || a.N != 1 || !reflect.DeepEqual(a.Findings, want) {
		t.Fatalf("NextAttempt = %+v, %v; want attempt 1 with %+v", a, err, want)
	}
	passed := tableOf(gate.Row{Kind: "build", State: "active", Result: "pass", Reason: "—"})
	if a, err = s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: aSHA, Table: passed, Findings: given}); err != nil || a.N != 2 || !reflect.DeepEqual(a.Findings, given) {
		t.Errorf("NextAttempt with a passing table and findings = %+v, %v; want attempt 2 with the given findings alone", a, err)
	}
	if _, err = s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: aSHA, Table: passed}); !errors.Is(err, ErrNoFailure) {
		t.Errorf("NextAttempt with neither = %v; want ErrNoFailure", err)
	}
}

// U5 (plan review condition 5): the expected failures of a rule batch in M2e. The control is every other kind: the same failures are findings.
func TestAnExpectedFailureOfABatchTaskIsNoFindingAndAnActiveRowStaysOne(t *testing.T) {
	failedRules := rules.Result{State: rules.Failure, Reason: rules.ReasonRulePath, Path: "rule-paths.tsv"}
	pending := gate.Row{Kind: "docs", State: "pending", Result: "fail", Reason: "pending: product path changed: a.go"}
	active := gate.Row{Kind: "exits", State: "active", Result: "fail", Reason: "exit 3"}
	rows := specTableRows(nil)
	for _, c := range []struct {
		kind              string
		gates, rules      int // the findings of the retry boundary
		wantRules, wantAt string
	}{
		{"batch", 0, 0, "", ""},
		{"architecture", 1, 1, "rules:rule path", "docs"},
		{"specification", 1, 1, "rules:rule path", "docs"},
	} {
		check := checkResultOf(c.kind, rows, true, failedRules)
		if len(check.Findings) != c.rules || check.Stop != "" || (c.rules == 1 && check.Findings[0].Rule != c.wantRules) {
			t.Errorf("%s: the rules failure gives %+v; want %d finding", c.kind, check, c.rules)
		}
		got := RetryGateFindings(c.kind, tableOf(pending))
		if len(got) != c.gates || (c.gates == 1 && got[0].Rule != "gate:"+c.wantAt) {
			t.Errorf("%s: a pending row that failed gives %+v; want %d finding", c.kind, got, c.gates)
		}
		// The active row is a finding for every kind, the pending one only where it is not expected.
		both := RetryGateFindings(c.kind, tableOf(pending, active))
		if len(both) != 1+c.gates || both[len(both)-1].Rule != "gate:exits" {
			t.Errorf("%s: a pending and an active row that failed give %+v; want the active one kept", c.kind, both)
		}
	}
	// The pending rows alone start no attempt of a batch task, with the rules failure beside them.
	s := newSessions(&standinStore{holder: "0123456789abcdef", files: map[string][]byte{}})
	if _, err := s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Kind: "batch", Base: aSHA, Table: tableOf(pending)}); !errors.Is(err, ErrNoFailure) {
		t.Errorf("NextAttempt of a batch task with a pending row that failed = %v; want ErrNoFailure", err)
	}
	if a, err := s.NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Kind: "batch", Base: aSHA, Table: tableOf(pending, active)}); err != nil || len(a.Findings) != 1 || a.Findings[0].Rule != "gate:exits" {
		t.Errorf("NextAttempt of a batch task with an active row that failed = %+v, %v; want attempt 1 with the active finding alone", a, err)
	}
}

// U6: a stop comes before the findings of the other check; a head with no row of the spec table is a stop, not a pass.
func TestACheckThatCouldNotRunStopsBeforeAnyFinding(t *testing.T) {
	failedSpec := specTableRows(map[int]string{1: "fail"})
	failedRules := rules.Result{State: rules.Failure, Reason: rules.ReasonRulePath, Path: "a"}
	for _, c := range []struct {
		name  string
		rows  [][]string
		found bool
		r     rules.Result
		stop  string
		n     int
	}{
		{"both failed", failedSpec, true, failedRules, "", 2},
		{"a rules error beside a spec failure", failedSpec, true, rules.Result{State: rules.Error, Reason: rules.ReasonRegister}, "layup/rules is error: register", 0},
		{"a rules result with no state", failedSpec, true, rules.Result{}, "layup/rules is error: error: no result", 0},
		{"no row of the spec table beside a rules failure", nil, false, failedRules, "layup/spec is error: no row of the spec check for the head", 0},
		{"the passes", specTableRows(nil), true, rules.Result{State: rules.Success, Reason: rules.ReasonNoRulePath}, "", 0},
	} {
		got := checkResultOf("specification", c.rows, c.found, c.r)
		if got.Stop != c.stop || len(got.Findings) != c.n {
			t.Errorf("%s: %+v; want the stop %q and %d findings", c.name, got, c.stop, c.n)
		}
	}
}
