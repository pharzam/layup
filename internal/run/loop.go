package run

// This file holds the step loop (docs/spec/loop.md, The step loop; task T-9nbf, #367, row 66k of the plan of M2e): the merge order, the scheduler that acts on one
// task at a time, the permission check and the parameter boundary of the steps shape, bet and loop, and the act on a task.

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
)

// maxLoopActs bounds the acts of one invocation of the step: each act moves a record on, so a step that does not end is a defect.
const maxLoopActs = 200

// A loopEnv is what the scheduler reads and does, so that it is told apart from the records and the forge.
type loopEnv struct {
	register       func(ctx context.Context) ([]Task, error)
	reopenings     func(ctx context.Context) ([][]string, error) // the rows of spec/reopenings.tsv, none when the record is absent
	merged         func(ctx context.Context, task string) (bool, error)
	openReopenings func(ctx context.Context, specTask string) error
	boundary       func(ctx context.Context) error
	act            func(ctx context.Context, t Task, original bool) (Advance, error)
	position       func(ctx context.Context, t Task) (string, error)
}

// mergeOrder gives the tasks of the register in the order the loop merges them.
func mergeOrder(register []Task, reopenings [][]string) ([]Task, error) { return nil, nil }

// pick gives the first task of the order with no recorded merge, and whether there is one; a task whose predecessor has not merged is an error.
func pick(order []Task, merged func(id string) bool) (Task, bool, error) { return Task{}, false, nil }

// originalSpec tells the specification task of Shape from a task that a changed class reopened.
func originalSpec(t Task, reopened map[string]bool) bool { return t.Kind == "specification" }

// An actKind is what the act on a task is, told from its records and the pull request of its latest attempt.
type actKind int

const (
	actDrive  actKind = iota // AdvanceShapeTask
	actNone                  // the task has merged
	actReread                // the task was redone: the register is read again
	actFinish                // the pull request is merged or closed: FinishTask records or waits
	actBring                 // a verified head: BringUpToDate, then FinishTask
)

// dispatch tells the act on a task from its view; terminal is whether the pull request of its latest attempt is merged or closed.
func dispatch(v loopView, terminal bool) actKind { return actDrive }

// boundaryRule tells whether the parameter boundary is skipped, from the bytes of parameters.tsv and tasks.tsv (nil for an absent file): only before the first task.
func boundaryRule(parameters, tasks []byte) (skip bool, err error) { return false, nil }

// wrapBoundary is atBoundary over the permission check, the boundary and the step it wraps.
func wrapBoundary(name string, permission func() Step, boundary func(context.Context) error, step func(context.Context) Step) func(context.Context) Step {
	return step
}

// runLoop is the scheduler of the step loop over its environment.
func runLoop(ctx context.Context, env loopEnv) Step { return Step{} }

// loopStep is the step loop.
func (r *state) loopStep(ctx context.Context) Step { return Step{} }

// atBoundary wraps a step: it checks the permissions of the loop, calls the parameter boundary and then runs the step.
func (r *state) atBoundary(name string, step func(context.Context) Step) func(context.Context) Step {
	return step
}

// bindingOptions gives the options of an architecture task of the bet: every option of each decision's table, the one it is bound to the only survivor.
func (r *state) bindingOptions(ctx context.Context, task string) (map[string]spec.Result, error) {
	return nil, nil
}
