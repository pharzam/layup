package run

// This file holds the question in the loop (docs/spec/loop.md, A question during the work, "The question in the loop"; task
// T-r7b9, #426, row 66z of the plan of M2e): the stages that take a session that ended needs_context to its owner role's answer and
// to the next attempt, as a position of the driver of a Shape task.

import "context"

// A qStage is the first stage of the question that the records show not done.
type qStage int

const (
	qNone          qStage = iota // no asking session in the latest attempt
	qHandoffAsk                  // P1: record the handoff asking role to owner role of an asking session
	qPostQuestion                // P2: post and copy the comment of a question row
	qWait                        // P3: the questions are not of one owner role: wait for the Operator
	qRunOwner                    // P4: run the answer session
	qPostAnswer                  // P6: post and copy the comment of an answer row
	qHandoffAnswer               // P7: record the handoff owner role to asking role
	qNextAttempt                 // P8: start the next attempt
	qStop                        // a stop: P1, P5, P7, or a row that cannot be composed
)

// A qStep is a stage with what its act needs: the session it concerns (the asking session of P1 and P2, the answer session of P6 and P7),
// the row and, for an answer, the question it answers, the owner role and the base of the answer session, the stop or the wait.
type qStep struct {
	Stage    qStage
	Session  string
	Row      []string
	Question string
	Owner    string
	Base     string
	Msg      string
}

// questionPosition tells the stage of the question of the attempt from the view alone.
func questionPosition(v loopView, attempt int) qStep { return qStep{} }

// askStage does the act of a stage; the driver reads the records again after it. It gives the stop or the wait of the stage, if any.
func (s *Sessions) askStage(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (stop, wait string, err error) {
	return "", "", nil
}
