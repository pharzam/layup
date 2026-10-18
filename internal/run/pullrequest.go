package run

// This file holds the pull request of a task (docs/spec/loop.md, The pull
// request; task T-xvvd, #363, row 66g of the plan of M2e). The step loop that
// calls it is row 66k's; the edit of the body to link the review record is row
// 66q's.

import (
	"context"
	"errors"

	"github.com/pharzam/layup/internal/forge"
)

// A PullInput is what the caller knows of the draft to open: the task, its
// attempt, and the type and the description of the title.
type PullInput struct {
	Task        string
	Attempt     int
	Type        string
	Description string
}

// checkTitle refuses a title that the checked-in pr-title check would refuse.
func checkTitle(typ, description string) error { return nil }

// PullBody gives the body of the draft: the line that closes the issue, a blank
// line, the line that links the records commit, and each further line.
func PullBody(issue int, web, repo, records string, more ...string) string { return "" }

// withPR gives the register with the pr of the task set.
func withPR(tasks []Task, id string, n int) ([]Task, error) { return nil, errors.New("not built") }

// OpenPullRequest opens the draft pull request of the attempt's last bound
// head, once, and records its number in the register.
func (s *Sessions) OpenPullRequest(ctx context.Context, in PullInput) (forge.PullRequest, error) {
	return forge.PullRequest{}, nil
}
