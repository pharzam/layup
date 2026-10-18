package run

// This file holds the ready and merge rule and the merge at the head SHA
// (docs/spec/loop.md, The ready and merge rule; task T-xvvd, #363, row 66g).

import (
	"context"

	"github.com/pharzam/layup/internal/forge"
)

// A Readiness tells whether a head may be marked ready and merged, and the first
// reason when it may not.
type Readiness struct {
	Ready  bool
	Reason string
}

// requiredJobs gives the names of the required checks of GitHub Actions that the
// rules hold, or an error for a malformed rule.
func requiredJobs(rules []forge.Rule) ([]string, error) { return nil, nil }

// Ready reads the forge and tells whether the head may be marked ready and merged.
func (s *Sessions) Ready(ctx context.Context, head string) (Readiness, error) {
	return Readiness{}, nil
}

// MarkReady marks the pull request of the task ready when its head is the head
// given and Ready holds.
func (s *Sessions) MarkReady(ctx context.Context, task, head string) (Readiness, error) {
	return Readiness{}, nil
}

// MergeTask merges the pull request of the task at the head SHA it read and
// records the event merged. It gives the merge commit.
func (s *Sessions) MergeTask(ctx context.Context, task, head string) (string, Readiness, error) {
	return "", Readiness{}, nil
}

// RecordMerged writes the event merged of a merged pull request, once.
func (s *Sessions) RecordMerged(ctx context.Context, task string, pr forge.PullRequest) error {
	return nil
}
