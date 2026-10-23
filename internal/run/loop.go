package run

// This file holds the step loop (docs/spec/loop.md, The step loop; task T-9nbf, #367, row 66k of the plan of M2e): the merge order, the scheduler that acts on one
// task at a time, the permission check and the parameter boundary of the steps shape, bet and loop, and the act on a task.

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
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

// mergeOrder gives the tasks of the register in the order the loop merges them: the specification task of Shape, the tasks of the reopenings in the order of their
// rows, the architecture task of Shape (the first of its kind), then the architecture tasks of the bet in register order. A kind that M2e does not run is an error.
func mergeOrder(register []Task, reopenings [][]string) ([]Task, error) {
	reopened := reopenedTasks(reopenings)
	var specTask, shape *Task
	var bet []Task
	for i, t := range register {
		switch {
		case t.Kind == "rendered":
			return nil, fmt.Errorf("%s: the kind rendered is not run in M2e (row 66t's, M2g)", t.ID)
		case t.Kind != "specification" && t.Kind != "architecture":
			return nil, fmt.Errorf("%s: the kind %s is not run in M2e", t.ID, t.Kind)
		case t.Kind == "specification" && reopened[t.ID]:
		case t.Kind == "specification" && specTask != nil:
			return nil, fmt.Errorf("%s: a second specification task that %s does not name", t.ID, reopeningsPath)
		case t.Kind == "specification":
			specTask = &register[i]
		case shape == nil:
			shape = &register[i]
		default:
			bet = append(bet, t)
		}
	}
	if specTask == nil {
		return nil, fmt.Errorf("the register holds no specification task")
	}
	order := []Task{*specTask}
	for _, row := range reopenings {
		if i := slices.IndexFunc(register, func(t Task) bool { return t.ID == row[reoTask] && t.Kind == "specification" }); i >= 0 {
			order = append(order, register[i])
		}
	}
	if shape != nil {
		order = append(order, *shape)
	}
	return append(order, bet...), nil
}

// pick gives the first task of the order with no recorded merge, and whether there is one; a task whose predecessor has not merged is an error.
func pick(order []Task, merged func(id string) bool) (Task, bool, error) {
	for _, t := range order {
		if merged(t.ID) {
			continue
		}
		for _, p := range t.Predecessors {
			if !merged(p) {
				return Task{}, false, fmt.Errorf("%s: its predecessor %s has not merged", t.ID, p)
			}
		}
		return t, true, nil
	}
	return Task{}, false, nil
}

// originalSpec tells the specification task of Shape from a task that a changed class reopened.
func originalSpec(t Task, reopened map[string]bool) bool {
	return t.Kind == "specification" && !reopened[t.ID]
}

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
func dispatch(v loopView, terminal bool) actKind {
	for _, e := range v.Events {
		if e[1] == "merged" {
			return actNone
		}
	}
	for _, e := range v.Events {
		if e[1] == "closed" {
			return actReread
		}
	}
	if terminal {
		return actFinish
	}
	if _, ok := verifiedView(v); ok {
		return actBring
	}
	return actDrive
}

// boundaryRule tells whether the parameter boundary is skipped, from the bytes of parameters.tsv and tasks.tsv (nil for an absent file): only before the first task.
func boundaryRule(parameters, tasks []byte) (skip bool, err error) {
	switch {
	case parameters != nil:
		return false, nil
	case tasks == nil:
		return true, nil
	}
	return false, fmt.Errorf("%s is not in the records, and %s holds tasks", parametersPath, tasksPath)
}

// wrapBoundary is atBoundary over the permission check, the boundary and the step it wraps.
func wrapBoundary(name string, permission func() Step, boundary func(context.Context) error, step func(context.Context) Step) func(context.Context) Step {
	return func(ctx context.Context) Step {
		if s := permission(); s.Result == "fail" {
			return failf(name, "forge: %s", s.Detail)
		}
		if err := boundary(ctx); err != nil {
			return failf(name, "parameters: %s", strings.SplitN(err.Error(), "\n", 2)[0])
		}
		return step(ctx)
	}
}

// runLoop is the scheduler of the step loop over its environment (loop.md, The step loop, items 4, 6 and 7).
func runLoop(ctx context.Context, env loopEnv) Step {
	const name = "loop"
	var merged []string
	acted, opened, stuck := 0, false, 0
	stuckTask := ""
	for {
		if acted > maxLoopActs {
			return failf(name, "the loop did not end in %d acts", maxLoopActs)
		}
		register, err := env.register(ctx)
		if err != nil {
			return fail(name, err)
		}
		rows, err := env.reopenings(ctx)
		if err != nil {
			return fail(name, err)
		}
		order, err := mergeOrder(register, rows)
		if err != nil {
			return fail(name, err)
		}
		has := map[string]bool{}
		for _, t := range order {
			if has[t.ID], err = env.merged(ctx, t.ID); err != nil {
				return fail(name, err)
			}
		}
		if !opened && has[order[0].ID] { // whenever the event exists, also at the start of the invocation
			if err := env.openReopenings(ctx, order[0].ID); err != nil {
				return fail(name, err)
			}
			opened = true
			continue
		}
		t, found, err := pick(order, func(id string) bool { return has[id] })
		switch {
		case err != nil:
			return fail(name, err)
		case !found && acted == 0:
			return done(name, "done before")
		case !found:
			return done(name, "the loop merged "+strings.Join(merged, ", "))
		}
		if acted > 0 {
			if err := env.boundary(ctx); err != nil {
				return failf(name, "parameters: %s", strings.SplitN(err.Error(), "\n", 2)[0])
			}
		}
		before, err := env.position(ctx, t)
		if err != nil {
			return fail(name, err)
		}
		adv, err := env.act(ctx, t, originalSpec(t, reopenedTasks(rows)))
		acted++
		switch {
		case err != nil:
			return fail(name, err)
		case adv.Stop != "":
			return failf(name, "%s", adv.Stop)
		case adv.Wait != "":
			return done(name, adv.Wait)
		case adv.Merged:
			merged, stuck = append(merged, t.ID), 0
			continue
		case adv.Replaced != "":
			stuck = 0
			continue
		}
		after, err := env.position(ctx, t)
		switch {
		case err != nil:
			return fail(name, err)
		case after != before || t.ID != stuckTask:
			stuck, stuckTask = 0, t.ID
		}
		if stuck++; after == before && stuck >= 2 {
			return failf(name, "%s: the loop made no progress at %s", t.ID, after)
		}
	}
}

// loopStep is the step loop: after the bet it merges every task of M2e in the merge order, one at a time.
func (r *state) loopStep(ctx context.Context) Step {
	return r.atBoundary("loop", r.runLoopStep)(ctx)
}

func (r *state) runLoopStep(ctx context.Context) Step {
	const name = "loop"
	if ok, err := r.betted(ctx); err != nil {
		return fail(name, err)
	} else if !ok {
		return failf(name, "bets.tsv holds no bet: the step bet comes first")
	}
	s := r.loopSessions()
	return runLoop(ctx, loopEnv{
		register: r.readRegisterRows,
		reopenings: func(ctx context.Context) ([][]string, error) {
			data, err := r.store.ReadFile(ctx, reopeningsPath)
			if err != nil || data == nil {
				return nil, err
			}
			rows, err := ReadReopenings(data)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", reopeningsPath, err)
			}
			return rows, nil
		},
		merged: func(ctx context.Context, task string) (bool, error) {
			events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
			return slices.ContainsFunc(events, func(e []string) bool { return e[1] == "merged" }), err
		},
		openReopenings: func(ctx context.Context, specTask string) error {
			_, err := r.openReopenings(ctx, specTask)
			return err
		},
		boundary: r.boundary,
		act:      r.actOn,
		position: func(ctx context.Context, t Task) (string, error) {
			events, err := s.readTable(ctx, eventsPath(t.ID), records.ReadEvents)
			n, _ := latestAttempt(events)
			return fmt.Sprintf("attempt %d with %d events", n, len(events)), err
		},
	})
}

// atBoundary wraps a step: it checks the permissions of the loop, calls the parameter boundary and then runs the step (loop.md, The step loop, item 5).
func (r *state) atBoundary(name string, step func(context.Context) Step) func(context.Context) Step {
	return wrapBoundary(name, r.loopPermissionsStep, r.boundary, step)
}

// boundary calls the parameter boundary of the records, once per call of a step; an absent parameters.tsv is skipped before the first task only.
func (r *state) boundary(ctx context.Context) error {
	params, err := r.store.ReadFile(ctx, parametersPath)
	if err != nil {
		return err
	}
	tasks, err := r.store.ReadFile(ctx, tasksPath)
	if err != nil {
		return err
	}
	if skip, err := boundaryRule(params, tasks); skip || err != nil {
		return err
	}
	if _, _, err = r.parametersBoundary(ctx); err != nil {
		return err
	}
	_, err = r.cycleBoundary(ctx)
	return err
}

// betted reports whether bets.tsv holds a row with the verdict bet; from then on every task belongs to the step loop (D1).
func (r *state) betted(ctx context.Context) (bool, error) {
	data, err := r.store.ReadFile(ctx, "bets.tsv")
	if err != nil || data == nil {
		return false, err
	}
	bets, err := ReadBets(data)
	if err != nil {
		return false, fmt.Errorf("bets.tsv: %w", err)
	}
	return slices.ContainsFunc(bets, func(b []string) bool { return b[betVerdict] == "bet" }), nil
}

// loopSessions gives the Sessions of the step loop.
func (r *state) loopSessions() *Sessions {
	s := r.sessions()
	if r.sessionsHook != nil {
		r.sessionsHook(s)
	}
	return s
}

// reopenedTask reports whether spec/reopenings.tsv names the task: a specification task that a changed class reopened is not the original one (D9).
func (s *Sessions) reopenedTask(ctx context.Context, task string) (bool, error) {
	data, err := s.Store.ReadFile(ctx, reopeningsPath)
	if err != nil || data == nil {
		return false, err
	}
	rows, err := ReadReopenings(data)
	if err != nil {
		return false, fmt.Errorf("%s: %w", reopeningsPath, err)
	}
	return reopeningOf(rows, task) != nil, nil
}

// actOn takes one task one act on (loop.md, The step loop, item 3). original tells the specification task of Shape.
func (r *state) actOn(ctx context.Context, t Task, original bool) (Advance, error) {
	s := r.loopSessions()
	v, err := s.shapeView(ctx, t.ID)
	if err != nil {
		return Advance{}, err
	}
	attempt, _ := latestAttempt(v.Events)
	terminal := false
	if attempt > 0 {
		pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: fmt.Sprintf("task/%s/%d", t.ID, attempt)})
		if err != nil {
			return Advance{}, err
		}
		terminal = found && (pr.Merged || pr.State == "closed")
	}
	kind := dispatch(v, terminal)
	if original && kind != actNone && kind != actReread && kind != actFinish && attempt > 0 && !renderedPushed(v.Events, attempt) {
		return Advance{Task: t.ID, Stop: fmt.Sprintf("waiting for the Operator: %s: an attempt after the bet has no rendered commit", t.ID)}, nil
	}
	fin := FinishInput{Task: t.ID, Records: r.store.Base, Harnesses: r.cfg.Harnesses, Models: r.cfg.Models, Copy: r.betCopy}
	if original {
		fin.Bet = r.renderedBet
	}
	switch kind {
	case actNone:
		return Advance{Task: t.ID}, nil
	case actReread:
		for _, e := range v.Events {
			if e[1] == "closed" {
				return Advance{Task: t.ID, Replaced: e[6]}, nil
			}
		}
	case actFinish:
		return s.FinishTask(ctx, fin)
	case actBring:
		// The pushes of the task's sessions are bound or refused before the tip for the base movement is chosen: BringUpToDate takes its tip from the
		// recorded bounds, and a close-out push that the forge took and whose answer was lost would otherwise be found on the branch where it did not expect it.
		if _, err := s.resolvePushes(ctx, t.ID); err != nil {
			return Advance{}, err
		}
		if v, err = s.shapeView(ctx, t.ID); err != nil {
			return Advance{}, err
		}
		attempt, _ = latestAttempt(v.Events)
		st, ok := verifiedView(v)
		if !ok {
			return Advance{Task: t.ID}, nil // the next pass reads the records again
		}
		cin := CarryInput{Task: t.ID, Verifier: st.Session, Harnesses: r.cfg.Harnesses, Models: r.cfg.Models}
		if original {
			rendered, err := s.renderedHead(ctx, fin, v, attempt, st.Head)
			if err != nil {
				return Advance{}, err
			}
			cin.Rendered = &RenderedInput{Head: rendered.Head, Render: rendered.Records}
		}
		if adv, err := s.BringUpToDate(ctx, cin); err != nil || !adv.Verified {
			return adv, err
		}
		return s.FinishTask(ctx, fin)
	}
	in, err := r.driverInput(ctx, s, t)
	if err != nil {
		return Advance{}, err
	}
	return s.AdvanceShapeTask(ctx, in)
}

// driverInput is what the driver is given for a task of the step loop: the head of the default branch for a first attempt, the first copy of the answer, and the
// options of an architecture task (D10).
func (r *state) driverInput(ctx context.Context, s *Sessions, t Task) (ShapeTaskInput, error) {
	base, err := s.Store.FetchBranch(ctx, s.Branch)
	if err != nil {
		return ShapeTaskInput{}, err
	}
	in := ShapeTaskInput{Task: t.ID, Base: base, Records: r.store.Base, Harnesses: r.cfg.Harnesses, Models: r.cfg.Models}
	if in.Copy, err = r.betCopy(ctx); err != nil {
		return in, err
	}
	if t.Kind == "architecture" {
		if in.Options, err = r.bindingOptions(ctx, t.ID); err == nil && len(in.Options) == 0 { // no binding: the task of Shape, or its redo
			in.Options, err = r.optionResults(ctx)
		}
	}
	return in, err
}

// bindingOptions gives the options of an architecture task of the bet: every option of each decision's table, the one it is bound to the only survivor, so that the one
// reader of the prompt file checks the result against the whole table. A task with no binding has none.
func (r *state) bindingOptions(ctx context.Context, task string) (map[string]spec.Result, error) {
	files, err := r.store.ReadFiles(ctx, "tasks/"+task+"/options")
	if err != nil {
		return nil, err
	}
	out := map[string]spec.Result{}
	for name, data := range files {
		decision, ok := strings.CutSuffix(name, ".tsv")
		if !ok || strings.Contains(decision, "/") {
			continue
		}
		bound, err := spec.ReadOptions(data)
		if err != nil || len(bound) == 0 {
			return nil, fmt.Errorf("%s: %v", bindingPath(task, decision), err)
		}
		table, err := r.store.ReadFile(ctx, "spec/options/"+decision+".tsv")
		if err != nil {
			return nil, err
		} else if table == nil {
			return nil, fmt.Errorf("the records have no spec/options/%s.tsv", decision)
		}
		rows, err := spec.ReadOptions(table)
		if err != nil {
			return nil, fmt.Errorf("spec/options/%s.tsv: %w", decision, err)
		}
		var res spec.Result
		for _, row := range rows {
			if !slices.ContainsFunc(res.Options, func(o spec.Option) bool { return o.Name == row[0] }) {
				res.Options = append(res.Options, spec.Option{Name: row[0]})
			}
		}
		if i := slices.IndexFunc(res.Options, func(o spec.Option) bool { return o.Name == bound[0][0] }); i >= 0 {
			res.Survivors = []spec.Option{res.Options[i]}
		}
		out[decision] = res
	}
	return out, nil
}
