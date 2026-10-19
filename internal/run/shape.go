package run

// This file holds the step shape (docs/spec/loop.md, REQ-005 — The steps of the loop after the Scaffold; task T-3f94, #364, row
// 66h of the plan of M2e): the questions left open, the options, the stop, and the two Shape tasks to their verified heads.

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// noSurvivorHeader is the first line of the comment of the round with no progress.
const noSurvivorHeader = "LAYUP Shape: no survivor"

// noSurvivorComment gives the comment that names, for each decision, each dropped option of both rounds with its reasons.
func noSurvivorComment(decisions []string, round1, round2 map[string]spec.Result) string { return "" }

// archStart gives the base of the architecture task arch: the verified head of the specification task, and a stop when the
// task already has an attempt at another base.
func archStart(specTask Advance, arch string, events [][]string) (base, stop string, err error) {
	return "", "", nil
}

// shapeStep is the step shape.
func (r *state) shapeStep(ctx context.Context) Step { return Step{Name: "shape"} }
