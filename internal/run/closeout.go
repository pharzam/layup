package run

// This file holds the close-out commit and the carry-over of layup/verify
// (docs/spec/loop.md, The merge at the head SHA and the close-out commit; task
// T-xvvd, #363, row 66g).

import (
	"context"

	"github.com/pharzam/layup/internal/rules"
)

// A RenderedInput is the rendered commit above the verified head: the commit and
// the records that code rendered into it.
type RenderedInput struct {
	Head   string
	Render []rules.Record
}

// A CarryInput is what the caller knows of the verdict that a head carries: the
// task, the verifier's session, the registers that decide its independence, and
// the rendered commit, nil for none.
type CarryInput struct {
	Task, Verifier    string
	Harnesses, Models [][]string
	Rendered          *RenderedInput
}

// closeOutPaths gives the two paths that the close-out commit may change.
func closeOutPaths(task string) []string { return nil }

// PushCloseOut is the push hook of the close-out session.
func (s *Sessions) PushCloseOut(ctx context.Context, id string, in CarryInput) error { return nil }

// CheckRendered tells whether the rendered commit is a permitted child of the
// verified head.
func (s *Sessions) CheckRendered(ctx context.Context, in CarryInput) (bool, error) {
	return false, nil
}

// CarryVerify posts layup/verify on the head that the verdict carries to, and
// gives the head.
func (s *Sessions) CarryVerify(ctx context.Context, in CarryInput, session string) (string, error) {
	return "", nil
}
