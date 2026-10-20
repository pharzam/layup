package run

// This file holds the close-out and the merge of a verified task (docs/spec/loop.md, The merge at the head SHA and the close-out
// commit; task T-ehb4, #425, row 66y of the plan of M2e): FinishTask takes a verified, up-to-date task to a recorded merge.

import (
	"context"

	"github.com/pharzam/layup/internal/rules"
)

// A RenderedBet is the typed input of the three checks on a head that holds the rendered commit of the first bet: the copy of the answer and
// the rendered records.
type RenderedBet struct {
	Copy    []byte
	Records []rules.Record
}

// A FinishInput is what FinishTask is given for one task: the task, the run's records commit as it is now, the harness and model registers,
// and, for the specification task, the inputs of its rendered head (D10).
type FinishInput struct {
	Task              string
	Records           func() string
	Harnesses, Models [][]string
	Bet               func(ctx context.Context, verified string) (RenderedBet, error)
}

// FinishTask takes a verified task to its recorded merge (D1 to D10).
func (s *Sessions) FinishTask(ctx context.Context, in FinishInput) (Advance, error) {
	return Advance{}, nil
}

// closePos is a position of the close-out session.
type closePos int

const (
	closeStart  closePos = iota // no session of the position, or one with no result: start one
	closeResume                 // a session ended with no push: its checked push
	closeStop                   // refused, or ended not completed
	closeBound                  // bound at Head
)

// A closeStep is the position of the close-out with the session it concerns, the message of a stop and the bound head.
type closeStep struct {
	Pos           closePos
	Session, Head string
	Msg           string
}

// finishCut is the index of the last handoff of a session of the verification role in the events, -1 for none.
func finishCut(v loopView) int { return 0 }

// closeOutPosition reads the close-out session of the attempt at the carried head from the events after the cut (D5).
func closeOutPosition(v loopView, cut, attempt int, carried string) closeStep { return closeStep{} }

// closeOutStep resolves the close-out at the tier of the plan and the position (D4).
func closeOutStep(plan Plan, pos Position) (Resolved, error) { return Resolved{}, nil }

// closedWait is the wait of a pull request closed and not merged (D2).
func closedWait(task, branch string, number int) string { return "" }

// notReady is the result of a merge call that did not merge (D6).
func notReady(task, head string, r Readiness) Advance { return Advance{} }

// upToDate tells whether the default branch's head is the attempt's base or an ancestor of head, the head that carries the verdict (D9).
func (s *Sessions) upToDate(ctx context.Context, base, head string) (bool, error) { return false, nil }

// renderedBet rebuilds the typed inputs of the rendered head of the first bet from the records and the verified head (D10).
func (r *state) renderedBet(ctx context.Context, verified string) (RenderedBet, error) {
	return RenderedBet{}, nil
}
