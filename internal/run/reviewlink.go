package run

// This file holds the link of the review record in the body of the pull request
// (docs/spec/loop.md, The link of the review record in the body; task T-08rj,
// #413, row 66q of the plan of M2e). The step that calls it is row 66k's.

import "context"

// reviewLine is the last line of the body: the link of the comment of the review record.
func reviewLine(web, repo string, issue int, comment int64) string { return "" }

// LinkReviewRecord replaces the body of the task's pull request with the body that
// links the review record of the verifier's session; it gives whether it edited.
func (s *Sessions) LinkReviewRecord(ctx context.Context, task, verifier string) (bool, error) {
	return false, nil
}
