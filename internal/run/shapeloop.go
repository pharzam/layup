package run

// This file holds the driver of a Shape task (docs/spec/loop.md, The driver of a Shape task; task
// T-jy6f, #415, row 66r of the plan of M2e).

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// ShapeTaskInput is what the driver is given for one task: the task, the base of the run, the
// records commit, the registers, and the results of the filter for an architecture task.
type ShapeTaskInput struct {
	Task, Base        string
	Records           func() string
	Harnesses, Models [][]string
	Options           map[string]spec.Result
}

// Advance is where the driver stopped: a verified head, a wait for the Operator or a stop.
type Advance struct {
	Task       string
	Verified   bool
	Wait, Stop string
}

// shapeCycleCap is the review-cycle cap of a Shape task.
const shapeCycleCap = 1

type position int

const (
	posNone position = iota
	posAttempt
	posPlan
	posPlanFinish
	posReview
	posReviewFinish
	posDeveloper
	posToVerifier
	posShip
	posVerifierFinish
	posVerified
	posNextAttempt
	posWait
	posStop
)

// loopView is what the position is told from.
type loopView struct {
	Task, Kind string
	Events     [][]string
	Sessions   [][]string
	Status     map[string]string
	Records    map[string]ReviewRecord
}

// loopStep is a position with what its act needs.
type loopStep struct {
	Pos          position
	Session, Msg string
	Attempt      int
	Base, Head   string
}

// loopPosition gives the first step of the task table that the records show not done.
func loopPosition(v loopView) loopStep { return loopStep{} }

// reviewRounds counts the rounds with findings.
func reviewRounds(v loopView) int { return 0 }

// findingsOf gives the findings of the record of the attempt before, and the bound head of that
// attempt whose gate table gives them when it has none.
func findingsOf(v loopView, attempt int) (found []Finding, gateHead string) { return nil, "" }

// AdvanceShapeTask takes one Shape task to its verified head.
func (s *Sessions) AdvanceShapeTask(ctx context.Context, in ShapeTaskInput) (Advance, error) {
	return Advance{}, nil
}

// readOnlyStep reports whether the session of a step only reads.
func readOnlyStep(step string) bool { return false }
