package run

// This file holds the architecture tasks that the bet opens (docs/spec/loop.md, The step bet; task T-2vva).

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/spec"
)

// settledDecisions gives the decisions that this bet settles, in name order: those of the answer's option lines that have more than one
// survivor (a decision with one survivor was written as decided before the bet).
func settledDecisions(options map[string]string, survivors map[string][]string) []string {
	var out []string
	for _, d := range slices.Sorted(maps.Keys(options)) {
		if len(survivors[d]) > 1 {
			out = append(out, d)
		}
	}
	return out
}

var (
	finalRequirement = regexp.MustCompile(`^(REQ|NFR)-[0-9]{3,}$`)
	draftRow         = regexp.MustCompile(`^R-[0-9]{3,}$`)
)

// archRequirements gives the requirement IDs of a chosen option from the constraint values of its rows: REQ-NNN and NFR-NNN kept when the
// final IDs of the draft name them, R-NNN mapped by finalID with its kind in the draft, any other value none. Sorted, once each; none is an error.
func archRequirements(decision, option string, rows, draft [][]string) ([]string, error) {
	final := finalIDs(draft)
	var out []string
	for _, r := range rows {
		if r[0] != option {
			continue
		}
		v := r[1]
		if draftRow.MatchString(v) {
			for _, d := range draft {
				if id, ok := finalID(v, d[1]); d[0] == v && ok {
					v = id
				}
			}
		}
		if finalRequirement.MatchString(v) && slices.Contains(final, v) && !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("the decision %s: its option %s names no requirement", decision, option)
	}
	slices.Sort(out)
	return out, nil
}

// bindingPath is the binding of an architecture task to the option of its decision.
func bindingPath(task, decision string) string {
	return "tasks/" + task + "/options/" + decision + ".tsv"
}

// bindingBytes is the binding of a decision: the rows of the chosen option, in the block options.
func bindingBytes(rows [][]string, option string) ([]byte, error) {
	var chosen [][]string
	for _, r := range rows {
		if r[0] == option {
			chosen = append(chosen, r)
		}
	}
	return table(spec.OptionsSchema, chosen)
}

// archIssueBlock is the body of the issue of an architecture task of a decision: the block bet-issue-architecture of loop.md.
const archIssueBlock = `LAYUP task {task}: the architecture of the decision {decision}.

Requirement IDs: {ids}
Chosen option: {option}, recorded in tasks/{task}/options/{decision}.tsv of the records.

Definition of Done: docs/spec/ and docs/adr/ hold the decision {decision} as decided with the option {option}, and nothing of the other options; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
`

// decisionTasks gives, for each decision that has a binding, the task that holds it.
func (r *state) decisionTasks(ctx context.Context) (map[string]string, error) {
	tasks, err := r.readRegisterRows(ctx)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, t := range tasks {
		if t.Kind != "architecture" {
			continue
		}
		files, err := r.store.ReadFiles(ctx, "tasks/"+t.ID+"/options")
		if err != nil {
			return nil, err
		}
		for name := range files {
			if d, ok := strings.CutSuffix(name, ".tsv"); ok {
				out[d] = t.ID
			}
		}
	}
	return out, nil
}

// archStage opens the architecture task of each settled decision that has no binding yet, in name order: its row and its binding in one
// records commit, after the preflight of both and before nothing else is written but the issue.
func (r *state) archStage(ctx context.Context, settled []string, answer spec.Answer, tables map[string][][]string, requirements map[string][]string, predecessors []string) error {
	have, err := r.decisionTasks(ctx)
	if err != nil {
		return err
	}
	for _, d := range settled {
		if _, ok := have[d]; ok {
			continue
		}
		option, ids := answer.Options[d], requirements[d]
		binding, err := bindingBytes(tables[d], option)
		if err != nil {
			return err
		}
		if _, err := spec.ReadOptions(binding); err != nil {
			return fmt.Errorf("the binding of %s: %w", d, err)
		}
		if _, err := r.openTask(ctx, openSpec{kind: "architecture", ids: ids, predecessors: predecessors,
			title: func(id string) string { return "LAYUP task " + id + ": the architecture of the decision " + d },
			body: func(id string, ids []string, _, _ string) string {
				return strings.NewReplacer("{task}", id, "{decision}", d, "{ids}", strings.Join(ids, " "), "{option}", option).Replace(archIssueBlock)
			},
			extra:      func(id string) map[string][]byte { return map[string][]byte{bindingPath(id, d): binding} },
			afterIssue: func() error { return betStopAt("issue") }}); err != nil {
			return err
		}
		if err := betStopAt("task"); err != nil {
			return err
		}
	}
	return nil
}
