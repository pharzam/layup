package run

// This file holds the step shape (docs/spec/loop.md, REQ-005 — The steps of the loop after the Scaffold; task T-3f94, #364, row
// 66h of the plan of M2e): the questions left open, the options, the stop, and the two Shape tasks to their verified heads.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
)

// noSurvivorHeader is the first line of the comment of the round with no progress.
const noSurvivorHeader = "LAYUP Shape: no survivor"

// archStart gives the base of the architecture task arch: the verified head of the specification task (plan amendment 1 of #364,
// #432), and a stop when the task already has an attempt at another base. A specification that is not verified gives no base.
func archStart(specTask Advance, arch string, events [][]string) (base, stop string, err error) {
	if !specTask.Verified {
		return "", "", nil
	}
	if specTask.Head == "" {
		return "", "", fmt.Errorf("%s: the task is verified and the driver gave no head", specTask.Task)
	}
	for _, e := range events {
		if e[1] == "attempt" && e[4] != specTask.Head {
			return "", fmt.Sprintf("%s: its base %s is not the verified head %s of %s", arch, e[4], specTask.Head, specTask.Task), nil
		}
	}
	return specTask.Head, "", nil
}

// shapeStep is the step shape: the questions left open, the options with the filter, the one rewrite and the stop, and then the
// specification task and the architecture task, each to its verified head. It reads the records and the forge at each run, stops at
// the first stage that waits or fails, and is done before when it wrote nothing.
func (r *state) shapeStep(ctx context.Context) Step {
	const name = "shape"
	before := r.store.Base()
	wait, err := r.shapeQuestions(ctx)
	if err != nil {
		return fail(name, err)
	} else if len(wait) > 0 {
		return done(name, "waiting: the intent form: "+strings.Join(wait, " "))
	}
	results, stop, err := r.optionsStage(ctx)
	if err != nil {
		return fail(name, err)
	} else if stop != "" {
		return done(name, stop)
	}
	tasks, wait2, err := r.tasksStage(ctx, results)
	if err != nil {
		return fail(name, err)
	} else if wait2 != "" {
		return done(name, wait2)
	}
	if r.store.Base() == before {
		return done(name, "done before")
	}
	return done(name, fmt.Sprintf("%d decisions, each with a survivor; %s and %s verified", len(results), tasks[0], tasks[1]))
}

// tasksStage is the third stage: the specification task to its verified head, then the architecture task, opened at that head. It
// gives the two task IDs, or the wait of a task, or the stop of one as an error.
func (r *state) tasksStage(ctx context.Context, results map[string]spec.Result) (tasks [2]string, wait string, err error) {
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return tasks, "", err
	}
	s := r.sessions()
	input := func(task, base string, options map[string]spec.Result) ShapeTaskInput {
		return ShapeTaskInput{Task: task, Base: base, Records: r.store.Base, Harnesses: r.cfg.Harnesses, Models: r.cfg.Models, Options: options}
	}
	advance := func(in ShapeTaskInput) (Advance, string, error) {
		adv, err := s.AdvanceShapeTask(ctx, in)
		if err == nil && adv.Stop != "" {
			err = errors.New(adv.Stop)
		}
		return adv, adv.Wait, err
	}
	specTask, err := r.openShapeTask(ctx, "specification")
	if err != nil {
		return tasks, "", err
	}
	specAdv, wait, err := advance(input(specTask.ID, base, nil))
	if err != nil || wait != "" {
		return tasks, wait, err
	}
	archTask, err := r.openShapeTask(ctx, "architecture")
	if err != nil {
		return tasks, "", err
	}
	events, err := s.readTable(ctx, eventsPath(archTask.ID), records.ReadEvents)
	if err != nil {
		return tasks, "", err
	}
	archBase, stop, err := archStart(specAdv, archTask.ID, events)
	if err != nil || stop != "" {
		if err == nil {
			err = errors.New(stop)
		}
		return tasks, "", err
	}
	if _, wait, err = advance(input(archTask.ID, archBase, results)); err != nil || wait != "" {
		return tasks, wait, err
	}
	return [2]string{specTask.ID, archTask.ID}, "", nil
}
