package run

// This file holds the step bet (docs/spec/loop.md, The step bet; task T-2vva, #366, row 66j of the plan of M2e).

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// betInput is what the answer to a brief is parsed against: the requirement IDs, the phase tags of the target, the fact IDs and the
// survivors of each decision, all of the kept basis of the brief.
type betInput struct {
	Required, Tags, Facts []string
	Survivors             map[string][]string
}

// phaseTagsOf gives the first column of the table under the heading "## 9." of a PRD: the phase tags of the target.
func phaseTagsOf(prd []byte) []string { return nil }

// optionsColumn gives the column options of bets.tsv: the chosen options as <decision>=<option>, sorted by decision, "" for none.
func optionsColumn(options map[string]string) []string { return nil }

// betAnswer gives the answer to the brief: the earliest copy, by numeric comment ID, that is of the issue, the first copy of its comment,
// no App's, by the approver's account, created strictly after the brief's copy, and that ParseAnswer takes without a finding. It gives
// the parsed answer and the row of the copy.
func betAnswer(copies []answerCopy, brief, approver []string, in betInput) (spec.Answer, []string, bool) {
	return spec.Answer{}, nil, false
}

// betStep is the step bet.
func (r *state) betStep(ctx context.Context) Step { return done("bet", "not wired") }
