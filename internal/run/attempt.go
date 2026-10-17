package run

// This file holds the new attempt that a failure of step 5 starts, with its
// findings (docs/spec/loop.md, The attempts; task T-6zj3, #361, row 66e of the
// plan of M2e). The call that notices the failure, the statuses that report it
// (#362) and the durability of the findings across a restart (gates.tsv, #362)
// are later wiring: the caller gives the table and the findings.

import (
	"context"
	"errors"

	"github.com/pharzam/layup/internal/gate"
)

// A Finding is one thing that goes back to the developer, in the fixed form of
// a rule, a file and a line (architecture.md §11). A finding with no file or no
// line, as a gate's row has none, has "" and 0.
type Finding struct {
	Rule, File string
	Line       int
	Text       string
}

// ErrNoFailure is the refusal of an attempt that no failure asks for: no row of
// the gate's table failed.
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

// NextAttemptInput is what a failure of step 5 gives: the task, the base commit
// the new attempt starts from, and the table of the gate of the failing head.
type NextAttemptInput struct {
	Task, Base string
	Table      gate.Table
}

// An Attempt is the attempt that a failure started: its number and the findings
// that its prompt file carries.
type Attempt struct {
	N        int
	Findings []Finding
}

// NextAttempt starts the attempt that follows a gate failure: the event attempt
// with the base (StartAttempt), after the findings of the table, never before
// them. A table with no failed row starts no attempt (ErrNoFailure). The new
// attempt starts from the base, so its commits have one author harness.
func (s *Sessions) NextAttempt(ctx context.Context, in NextAttemptInput) (Attempt, error) {
	found := GateFindings(in.Table)
	if len(found) == 0 {
		return Attempt{}, ErrNoFailure
	}
	n, err := s.StartAttempt(ctx, in.Task, in.Base)
	if err != nil {
		return Attempt{}, err
	}
	return Attempt{N: n, Findings: found}, nil
}
