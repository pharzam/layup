package run

// This file holds the repairs of a task session at a restart (docs/spec/loop.md,
// REQ-005 — The steps of the loop after the Scaffold; row 66v of the plan, task
// T-syxm, #422): a push with no result, and the comment of a session on the
// task's issue.

import "context"

// Repair is what RepairTask wrote: the sessions it bound, the sessions whose push
// it refused as push-lost, the comments it posted, and the new attempt it
// started (0 for none).
type Repair struct {
	Bound, Lost, Posted []string
	Attempt             int
}

// RepairTask repairs the sessions of task before a driver reads its position.
// This stub repairs nothing yet.
func (s *Sessions) RepairTask(ctx context.Context, task string) (Repair, error) {
	return Repair{}, nil
}

// RepairPushes resolves the pushes of task that have no result. This stub repairs
// nothing yet.
func (s *Sessions) RepairPushes(ctx context.Context, task string) (Repair, error) {
	return Repair{}, nil
}

// RepostComments posts the comment of each session of task that is not on the
// task's issue. This stub posts nothing yet.
func (s *Sessions) RepostComments(ctx context.Context, task string) ([]string, error) {
	return nil, nil
}
