package run

// This file holds the new attempt that a failure of step 5 starts, with its
// findings (docs/spec/loop.md, The attempts; task T-6zj3, #361, row 66e of the
// plan of M2e). The call that notices the failure is later wiring (#367): the
// caller gives the table, which GatesTable (gates.tsv, task T-q0kr) rebuilds after
// a restart, and the findings.

import (
	"context"
	"errors"
	"slices"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/rules"
)

// A Finding is one thing that goes back to the developer, in the fixed form of
// a rule, a file and a line (architecture.md §11). A finding with no file or no
// line, as a gate's row has none, has "" and 0.
type Finding struct {
	Rule, File string
	Line       int
	Text       string
}

// ErrNoFailure is the refusal of an attempt that no failure asks for: no finding of
// the gate's table or of the given checks.
var ErrNoFailure = errors.New("no row of the gate's table failed")

// GateFindings gives the findings of a gate's table (docs/spec/gate.md): each
// row whose result is neither pass nor clear, in the order of the manifest, as
// the rule gate:<kind> and the row's reason. A not-active row is a finding, as
// it never counts as a pass (NFR-004).
func GateFindings(t gate.Table) []Finding {
	var out []Finding
	for _, r := range t.Rows {
		if r.Result != "pass" && r.Result != "clear" {
			out = append(out, Finding{Rule: "gate:" + r.Kind, Text: r.Reason})
		}
	}
	return out
}

// SpecFindings gives the findings of layup/spec from the rows (rule, result, detail) of a
// head's table: each fail row, in the order of the rules, as spec:<rule> with its detail
// (loop.md, The attempts; task T-tnxn, #427, row 66aa). A clear row is none.
func SpecFindings(rows [][]string) []Finding {
	var out []Finding
	for _, r := range rows {
		if r[1] == "fail" {
			out = append(out, Finding{Rule: "spec:" + r[0], Text: r[2]})
		}
	}
	return out
}

// RulesFindings gives the finding of layup/rules: the failure of the basic rule, as
// rules:<reason> with its path as the file and the text of the status. A success and an
// error give none.
func RulesFindings(r rules.Result) []Finding {
	if r.State != rules.Failure {
		return nil
	}
	return []Finding{{Rule: "rules:" + r.Reason, File: r.Path, Text: RulesStatus(r).Description}}
}

// expectedFailure says whether the failure of a context is one of the two that
// architecture.md §8 step 5 does not send back, in M2e: for a task of the kind batch, the
// failure of layup/rules and the pending kinds of layup/gates.
func expectedFailure(kind, context string) bool {
	return kind == "batch" && (context == RulesContext || context == GatesContext)
}

// RetryGateFindings is GateFindings for the retry boundary: the rows of a pending kind are
// no finding for a task whose failure of layup/gates is expected, an active row always is.
func RetryGateFindings(kind string, t gate.Table) []Finding {
	if expectedFailure(kind, GatesContext) {
		t.Rows = slices.DeleteFunc(slices.Clone(t.Rows), func(r gate.Row) bool { return r.State == "pending" })
	}
	return GateFindings(t)
}

// A CheckResult is what the records of layup/spec and layup/rules give for a head: the
// findings to retry with, or a stop that comes before any retry.
type CheckResult struct {
	Findings []Finding
	Stop     string
}

// checkResultOf classifies the rows of a head's spec table (found: it has any) and the
// result of the basic rule, for a task of the kind: a stop for a head without a row and for
// a rules error, else the findings of both.
func checkResultOf(kind string, rows [][]string, found bool, r rules.Result) CheckResult {
	if !found {
		return CheckResult{Stop: SpecContext + " is error: no row of the spec check for the head"}
	}
	if st := RulesStatus(r); st.State == "error" {
		return CheckResult{Stop: RulesContext + " is error: " + st.Description}
	}
	out := SpecFindings(rows)
	if !expectedFailure(kind, RulesContext) {
		out = append(out, RulesFindings(r)...)
	}
	return CheckResult{Findings: out}
}

// CheckFindings reads the records of a head whose statuses were posted (the spec table, the
// register commit of its bound session) and classifies them. A read that fails is an error.
func (s *Sessions) CheckFindings(ctx context.Context, in HeadInput) (CheckResult, error) {
	t, err := s.taskRow(ctx, in.Task)
	if err != nil {
		return CheckResult{}, err
	}
	rows, found, err := s.SpecTable(ctx, in.Task, in.Head)
	if err != nil {
		return CheckResult{}, err
	}
	r, err := s.RulesResult(ctx, RulesInput{Task: in.Task, Base: in.Base, Head: in.Head, PreRev: in.PreRev, RegisterRev: in.RegisterRev, Records: in.Records})
	if err != nil {
		return CheckResult{}, err
	}
	return checkResultOf(t.Kind, rows, found, r), nil
}

// NextAttemptInput is what a failure of step 5 gives: the task and its kind, the base
// commit the new attempt starts from, the table of the gate of the failing head, and the
// findings of the other checks.
type NextAttemptInput struct {
	Task, Kind, Base string
	Table            gate.Table
	Findings         []Finding
}

// An Attempt is the attempt that a failure started: its number and the findings
// that its prompt file carries.
type Attempt struct {
	N        int
	Findings []Finding
}

// NextAttempt starts the attempt that follows a failure of step 5: the event attempt
// with the base (StartAttempt), after the findings of the table and the given ones,
// never before them, the gate's first. No finding starts no attempt (ErrNoFailure).
// The new attempt starts from the base, so its commits have one author harness.
func (s *Sessions) NextAttempt(ctx context.Context, in NextAttemptInput) (Attempt, error) {
	found := append(RetryGateFindings(in.Kind, in.Table), in.Findings...)
	if len(found) == 0 {
		return Attempt{}, ErrNoFailure
	}
	n, err := s.StartAttempt(ctx, in.Task, in.Base)
	if err != nil {
		return Attempt{}, err
	}
	return Attempt{N: n, Findings: found}, nil
}
