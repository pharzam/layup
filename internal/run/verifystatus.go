package run

// This file holds the status layup/verify (docs/spec/loop.md, layup/verify and
// The independence of the reviewers; task T-q0kr, #362, row 66f of the plan of M2e).

import "context"

// A VerifyInput is what the caller knows of the verification of a head: the task,
// the base of its change, the head, the verifier's session (none before it starts),
// whether it waits for its turn, and the harness and model registers.
type VerifyInput struct {
	Task, Base, Head, Session string
	Waiting                   bool
	Harnesses, Models         [][]string
}

// VerifyAt gives the SHA that layup/verify is posted on and its status.
func (s *Sessions) VerifyAt(context.Context, VerifyInput) (string, Status, error) {
	return "", Status{}, nil
}
