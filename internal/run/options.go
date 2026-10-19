package run

// This file holds the options of the step shape (docs/spec/loop.md, REQ-012 — Shape; task T-3f94, #364, row 66h of the plan of
// M2e): the constraints that apply, the options session and its output, and the filter of every decision.

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// adrsAt gives the decision records of the files of docs/adr/ (by base name) with their status lines.
func adrsAt(files map[string][]byte) ([]spec.ADR, error) { return nil, nil }

// constraintsAt gives the constraints that apply, from the records and from the target at base.
func (r *state) constraintsAt(ctx context.Context, base string) ([]spec.Constraint, error) {
	return nil, nil
}

// splitOptions checks the output of an options session and splits it into the table of each decision.
func splitOptions(data []byte) (map[string][]byte, error) { return nil, nil }

// optionTables splits the files of spec/options/ into the tables of the current round and those of round 1.
func optionTables(files map[string][]byte) (current, round1 map[string][]byte) { return nil, nil }

// filterAll applies the filter to the table of each decision.
func filterAll(tables map[string][]byte, ids []string) (map[string]spec.Result, error) {
	return nil, nil
}

// rewriteFiles checks the output of the rewrite session against the tables it replaces and gives the records to keep.
func rewriteFiles(old map[string][]byte, data []byte) (map[string][]byte, error) { return nil, nil }

// optionsPrompt is the prompt of the options session, with the constraints that apply written into it.
func optionsPrompt(cs []spec.Constraint) func(input map[string]string, output string) []byte {
	return func(map[string]string, string) []byte { return nil }
}

// rewritePrompt is the prompt of the rewrite session.
func rewritePrompt(cs []spec.Constraint) func(input map[string]string, output string) []byte {
	return func(map[string]string, string) []byte { return nil }
}
