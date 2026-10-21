package run

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// The scheduler of the step loop (loop.md, The step loop; task T-9nbf, row 66k) over injected reads and acts.

func tk(id, kind string, pred ...string) Task { return Task{ID: id, Kind: kind, Predecessors: pred} }

func reopenRows(ids ...string) [][]string {
	var rows [][]string
	for _, id := range ids {
		row := make([]string, reoCopy+1)
		row[reoTask] = id
		rows = append(rows, row)
	}
	return rows
}

func idsOf(tasks []Task) []string {
	var out []string
	for _, t := range tasks {
		out = append(out, t.ID)
	}
	return out
}

func TestTheMergeOrderIsSpecificationReopenedShapeArchitectureThenBetArchitecture(t *testing.T) {
	register := []Task{tk("S", "specification"), tk("A", "architecture"), tk("B1", "architecture", "S", "A"), tk("R2", "specification", "S"), tk("R1", "specification", "S"), tk("B2", "architecture", "S", "A")}
	got, err := mergeOrder(register, reopenRows("R1", "R2"))
	if want := []string{"S", "R1", "R2", "A", "B1", "B2"}; err != nil || !reflect.DeepEqual(idsOf(got), want) {
		t.Errorf("the order = %v, %v; want %v (the order of the rows of the reopenings, then Shape's architecture, then the bet's in register order)", idsOf(got), err, want)
	}
	if got, err := mergeOrder(register[:3], nil); err != nil || !reflect.DeepEqual(idsOf(got), []string{"S", "A", "B1"}) {
		t.Errorf("with no reopening: %v, %v", idsOf(got), err)
	}
}

func TestAKindThatM2eDoesNotRunFailsClosed(t *testing.T) {
	for _, kind := range []string{"batch", "plan", "build", "rendered"} {
		_, err := mergeOrder([]Task{tk("S", "specification"), tk("X", kind, "S")}, nil)
		if err == nil || !strings.Contains(err.Error(), "X: the kind "+kind+" is not run in M2e") || (kind == "rendered") != strings.Contains(err.Error(), "66t") {
			t.Errorf("kind %s: %v; want the refusal, naming row 66t for rendered only", kind, err)
		}
	}
}

func TestATaskWithAnUnmergedPredecessorIsNotStarted(t *testing.T) {
	merged := map[string]bool{"S": true}
	order := []Task{tk("S", "specification"), tk("A", "architecture", "S", "Z")}
	if got, ok, err := pick(order, func(id string) bool { return merged[id] }); err == nil || err.Error() != "A: its predecessor Z has not merged" || ok || got.ID != "" {
		t.Errorf("pick = %+v, %v, %v; want the failure that names the predecessor", got, ok, err)
	}
	if got, ok, err := pick(order[:1], func(id string) bool { return false }); err != nil || !ok || got.ID != "S" {
		t.Errorf("the first task with no predecessor: %+v, %v, %v", got, ok, err)
	}
	if _, ok, err := pick(order[:1], func(id string) bool { return true }); err != nil || ok {
		t.Errorf("every task merged: %v, %v; want none", ok, err)
	}
}

// fakeLoop is a world for runLoop: the register as a function of the acts, the merged tasks, and a script of the answers of the acts.
type fakeLoop struct {
	register   func(f *fakeLoop) []Task
	rows       []string
	merged     map[string]bool
	acts       []string
	opened     int
	boundaries int
	answer     func(f *fakeLoop, t Task) (Advance, error)
	stamps     map[string]int
}

func (f *fakeLoop) env() loopEnv {
	return loopEnv{
		register:       func(context.Context) ([]Task, error) { return f.register(f), nil },
		reopenings:     func(context.Context) ([][]string, error) { return reopenRows(f.rows...), nil },
		merged:         func(_ context.Context, id string) (bool, error) { return f.merged[id], nil },
		openReopenings: func(context.Context, string) error { f.opened++; return nil },
		boundary:       func(context.Context) error { f.boundaries++; return nil },
		position:       func(_ context.Context, t Task) (string, error) { return fmt.Sprint("stamp ", f.stamps[t.ID]), nil },
		act: func(_ context.Context, t Task, _ bool) (Advance, error) {
			f.acts = append(f.acts, t.ID)
			adv, err := f.answer(f, t)
			if adv.Merged {
				f.merged[t.ID] = true
			}
			if adv.Verified || adv.Merged || adv.Replaced != "" || adv.Attempt != 0 {
				f.stamps[t.ID]++
			}
			return adv, err
		},
	}
}

func newFakeLoop(register ...Task) *fakeLoop {
	return &fakeLoop{register: func(*fakeLoop) []Task { return register }, merged: map[string]bool{}, stamps: map[string]int{}}
}

func TestOnlyTheFirstUnmergedTaskIsActedOn(t *testing.T) {
	f := newFakeLoop(tk("S", "specification"), tk("A", "architecture", "S"))
	f.answer = func(f *fakeLoop, t Task) (Advance, error) {
		if t.ID == "S" && f.stamps["S"] == 0 {
			return Advance{Task: t.ID, Verified: true}, nil // a verified head is no merge: the same task goes on
		}
		return Advance{Task: t.ID, Merged: true}, nil
	}
	got := runLoop(context.Background(), f.env())
	if want := []string{"S", "S", "A"}; !slices.Equal(f.acts, want) || got.Result != "done" || got.Detail != "the loop merged S, A" {
		t.Errorf("acts %v and the row %+v; want %v, never A before S has merged, and the merged tasks named", f.acts, got, want)
	}
}

func TestTheRegisterIsReadAgainAfterAReplacementOrAReopening(t *testing.T) {
	f := newFakeLoop()
	f.register = func(f *fakeLoop) []Task {
		out := []Task{tk("S", "specification")}
		if f.opened > 0 {
			f.rows = []string{"R"}
			out = append(out, tk("R", "specification", "S"))
		}
		if f.stamps["Aold"] == 0 {
			return append(out, tk("Aold", "architecture", "S"))
		}
		return append(out, tk("Anew", "architecture", "S"))
	}
	f.answer = func(f *fakeLoop, t Task) (Advance, error) {
		if t.ID == "Aold" {
			return Advance{Task: t.ID, Replaced: "Anew"}, nil
		}
		return Advance{Task: t.ID, Merged: true}, nil
	}
	got := runLoop(context.Background(), f.env())
	if want := []string{"S", "R", "Aold", "Anew"}; !slices.Equal(f.acts, want) || f.opened != 1 || got.Result != "done" {
		t.Errorf("acts %v after %d openings, row %+v; want %v: the reopened task is picked after the merge of S, the replaced ID never again, the new one", f.acts, f.opened, got, want)
	}
	// The call is made whenever the event of the specification exists, also at the start of the invocation.
	g := newFakeLoop(tk("S", "specification"))
	g.merged["S"] = true
	if runLoop(context.Background(), g.env()); g.opened != 1 {
		t.Errorf("openReopenings called %d times for a specification that had merged before the invocation; want 1", g.opened)
	}
}

func TestDoneBeforeOnlyWhenThisInvocationActedOnNoTask(t *testing.T) {
	f := newFakeLoop(tk("S", "specification"), tk("A", "architecture", "S"))
	f.merged["S"], f.merged["A"] = true, true
	f.answer = func(*fakeLoop, Task) (Advance, error) { t.Fatal("an act on a merged task"); return Advance{}, nil }
	if got := runLoop(context.Background(), f.env()); got != (Step{"loop", "done", "done before"}) {
		t.Errorf("every task merged: %+v; want done before", got)
	}
	g := newFakeLoop(tk("S", "specification"), tk("A", "architecture", "S"))
	g.merged["S"] = true
	g.answer = func(g *fakeLoop, t Task) (Advance, error) { return Advance{Task: t.ID, Merged: true}, nil }
	if got := runLoop(context.Background(), g.env()); got != (Step{"loop", "done", "the loop merged A"}) {
		t.Errorf("one task merged by the invocation: %+v; want the loop merged A", got)
	}
}

func TestAWaitIsTheLastRowAndAStopIsAFail(t *testing.T) {
	for _, c := range []struct {
		adv  Advance
		err  error
		want Step
	}{
		{Advance{Wait: "no word that says it"}, nil, Step{"loop", "done", "no word that says it"}},
		{Advance{Stop: "waiting for the Operator: A: stops, by its field"}, nil, Step{"loop", "fail", "waiting for the Operator: A: stops, by its field"}},
		{Advance{}, errors.New("a read failed"), Step{"loop", "fail", "a read failed"}},
	} {
		f := newFakeLoop(tk("S", "specification"), tk("A", "architecture", "S"))
		f.answer = func(f *fakeLoop, t Task) (Advance, error) {
			if t.ID == "S" {
				return Advance{Task: "S", Merged: true}, nil
			}
			return c.adv, c.err
		}
		if got := runLoop(context.Background(), f.env()); got != c.want || len(f.acts) != 2 {
			t.Errorf("%+v, %v: row %+v after %d acts; want %+v after 2", c.adv, c.err, got, len(f.acts), c.want)
		}
	}
}

func TestAnActThatMakesNoProgressFails(t *testing.T) {
	f := newFakeLoop(tk("S", "specification"))
	f.answer = func(f *fakeLoop, t Task) (Advance, error) { return Advance{Task: t.ID, Attempt: 0}, nil } // nothing moves
	got := runLoop(context.Background(), f.env())
	if got.Result != "fail" || got.Detail != "S: the loop made no progress at stamp 0" || len(f.acts) != 2 {
		t.Errorf("row %+v after %d acts; want the failure at the second act, and exactly 2 acts", got, len(f.acts))
	}
	g := newFakeLoop(tk("S", "specification"))
	n := 0
	g.answer = func(g *fakeLoop, t Task) (Advance, error) { // an act that moves the records on is no failure, however often it repeats
		if n++; n < 5 {
			return Advance{Task: t.ID, Verified: true}, nil
		}
		return Advance{Task: t.ID, Merged: true}, nil
	}
	if got := runLoop(context.Background(), g.env()); got.Result != "done" || len(g.acts) != 5 {
		t.Errorf("progressing acts: %+v after %d acts; want done after 5", got, len(g.acts))
	}
}

func TestTheRestartNamesEndWithLoopOnlyAfterTheScaffold(t *testing.T) {
	if got := restartNames(true, true); !slices.Equal(got[len(got)-3:], []string{"shape", "bet", "loop"}) {
		t.Errorf("the last rows after the Scaffold = %v; want shape, bet, loop", got[len(got)-3:])
	}
	for _, c := range [][2]bool{{false, false}, {false, true}, {true, false}} {
		if got := restartNames(c[0], c[1]); slices.Contains(got, "loop") {
			t.Errorf("restartNames(%v, %v) = %v; want no loop row", c[0], c[1], got)
		}
	}
	if !strings.Contains(RunStepsSchema.Columns[0].Type, "|bet|loop)") {
		t.Errorf("the enum of run-steps = %s; want loop after bet", RunStepsSchema.Columns[0].Type)
	}
}

func TestTheBoundaryIsSkippedOnlyForAnAbsentTable(t *testing.T) {
	for _, c := range []struct {
		parameters, tasks []byte
		skip, fails       bool
	}{{nil, nil, true, false}, {nil, []byte("task\n"), false, true}, {[]byte("p"), nil, false, false}, {[]byte("p"), []byte("t"), false, false}} {
		if skip, err := boundaryRule(c.parameters, c.tasks); skip != c.skip || (err != nil) != c.fails {
			t.Errorf("parameters %q, tasks %q: skip %v, err %v; want skip %v, failure %v", c.parameters, c.tasks, skip, err, c.skip, c.fails)
		}
	}
	var calls []string
	perm := func() Step { calls = append(calls, "permission"); return done("forge", "ok") }
	step := func(context.Context) Step { calls = append(calls, "step"); return done("shape", "ran") }
	bound := func(context.Context) error { calls = append(calls, "boundary"); return nil }
	if got := wrapBoundary("shape", perm, bound, step)(context.Background()); got != (Step{"shape", "done", "ran"}) || !slices.Equal(calls, []string{"permission", "boundary", "step"}) {
		t.Errorf("the wrapped step = %+v after %v; want the permission, the boundary, then the step", got, calls)
	}
	calls = nil
	missing := func() Step {
		calls = append(calls, "permission")
		return failf("forge", "the installation lacks the permissions of the task loop: statuses: write")
	}
	if got := wrapBoundary("bet", missing, bound, step)(context.Background()); got != (Step{"bet", "fail", "forge: the installation lacks the permissions of the task loop: statuses: write"}) || !slices.Equal(calls, []string{"permission"}) {
		t.Errorf("a missing permission: %+v after %v; want the row of the step, forge: …, and neither boundary nor step", got, calls)
	}
	calls = nil
	broken := func(context.Context) error {
		calls = append(calls, "boundary")
		return errors.New("parameters.tsv: line 2: broken")
	}
	if got := wrapBoundary("loop", perm, broken, step)(context.Background()); got != (Step{"loop", "fail", "parameters: parameters.tsv: line 2: broken"}) || slices.Contains(calls, "step") {
		t.Errorf("a boundary that fails: %+v after %v; want parameters: …, and no step", got, calls)
	}
}

func TestTheActDispatchesByTheRecords(t *testing.T) {
	ev := func(kinds ...string) loopView {
		v := loopView{Task: "T", Kind: "architecture"}
		for i, k := range kinds {
			v.Events = append(v.Events, []string{fmt.Sprint(i + 1), k, "1", "", "", "", "", ""})
		}
		return v
	}
	for _, c := range []struct {
		name     string
		v        loopView
		terminal bool
		want     actKind
	}{
		{"merged", ev("attempt", "merged"), false, actNone},
		{"closed", ev("attempt", "closed"), false, actReread},
		{"closed, then a pull request that is merged", ev("attempt", "closed"), true, actReread},
		{"a pull request that is merged or closed", ev("attempt"), true, actFinish},
		{"no attempt", ev(), false, actDrive},
		{"an attempt that is not verified", ev("attempt"), false, actDrive},
		{"a handback: a repair's new attempt", ev("attempt", "refused", "attempt"), false, actDrive},
	} {
		if got := dispatch(c.v, c.terminal); got != c.want {
			t.Errorf("%s: dispatch = %d; want %d", c.name, got, c.want)
		}
	}
	verified := upTo("architecture", 15).v
	if got := dispatch(verified, false); got != actBring {
		t.Errorf("a verified view: dispatch = %d; want BringUpToDate then FinishTask", got)
	}
	// The same view after its close-out is bound is still verified: the cut is the last handoff of a verifier (F1).
	closed := upTo("architecture", 15).start("S-c1", stepRole("close-out"), "1", shH).ev("push", "1", "S-c1", "", shG, "task/T-sp01/1").ev("bound", "1", "S-c1", "", shG, "task/T-sp01/1").v
	if p := loopPosition(closed); p.Pos != posStop {
		t.Fatalf("the driver's position of a view with its close-out bound = %d; want the stop that F1 shows", p.Pos)
	}
	if got := dispatch(closed, false); got != actBring {
		t.Errorf("a view whose close-out is bound: dispatch = %d; want BringUpToDate, not a stop of the driver", got)
	}
}

// A verified view is the view of the latest attempt: a later attempt after the cut is not a verified head (the check of FinishTask, kept when the cut is shared).
func TestAVerifiedViewIsOfTheLatestAttempt(t *testing.T) {
	b := upTo("architecture", 15)
	if _, ok := verifiedView(b.v); !ok {
		t.Fatal("the view of a verified head is not verified")
	}
	if st, ok := verifiedView(b.ev("attempt", "2", "", shB, "", "").v); ok {
		t.Errorf("a view whose latest attempt starts after the cut = %+v, verified; want it not verified", st)
	}
}

func TestAReopenedSpecificationIsNotTheOriginal(t *testing.T) {
	reopened := reopenedTasks(reopenRows("R"))
	if originalSpec(tk("S", "specification"), reopened) != true || originalSpec(tk("R", "specification"), reopened) != false || originalSpec(tk("A", "architecture"), reopened) != false {
		t.Errorf("originalSpec: S %v, R %v, A %v; want true, false, false", originalSpec(tk("S", "specification"), reopened), originalSpec(tk("R", "specification"), reopened), originalSpec(tk("A", "architecture"), reopened))
	}
}
