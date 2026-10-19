package run

// This file holds the opening of a Shape task (docs/spec/loop.md, The task register; task T-jy6f,
// #415, row 66r of the plan of M2e).

import "context"

// finalID gives the stable requirement ID of a row R-NNN of the draft by its kind: REQ-NNN for a
// requirement, NFR-NNN for a non-functional row; a success criterion has none.
func finalID(id, kind string) (string, bool) { return "", false }

// finalIDs gives the IDs of every row of the draft that has one, in the order of the draft.
func finalIDs(draft [][]string) []string { return nil }

// shapeIssueTitle is the title of the issue of a Shape task.
func shapeIssueTitle(id, kind string) string { return "" }

// shapeIssueText is the body of the issue of a Shape task: the block of its kind in loop.md.
func shapeIssueText(id, kind string, ids []string, sum, records string) string { return "" }

// openShapeTask opens the Shape task of the kind once (loop.md, The opening of a Shape task).
func (r *state) openShapeTask(ctx context.Context, kind string) (Task, error) { return Task{}, nil }
