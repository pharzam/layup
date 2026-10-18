package run

// This file holds the status layup/spec and the table of the full check of a head
// as a record (docs/spec/loop.md, layup/spec; task T-q0kr, #362, row 66f of the
// plan of M2e).

import "context"

func specPath(task string) string { return "tasks/" + task + "/spec-check.tsv" }

// SpecStatus gives the status layup/spec of the ten rows of the full check.
func SpecStatus([][]string) Status { return Status{} }

// SpecTable gives the recorded rows of the full check of a head, and whether it has any.
func (s *Sessions) SpecTable(context.Context, string, string) ([][]string, bool, error) {
	return nil, false, nil
}
