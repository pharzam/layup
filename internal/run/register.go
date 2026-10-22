package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/tsv"
)

// A Task is one row of the task register, tasks.tsv (docs/spec/records.md,
// the block tasks; task T-cak3, #357, row 66a of the plan of M2e).
type Task struct {
	ID           string
	Issue        int
	Kind         string
	Requirements []string
	Size         string // "" is — (a task with no size class)
	Predecessors []string
	PR           int // 0 is — (until the pull request is opened)
	Opened       time.Time
}

// TasksSchema is the Go value of the block tasks of docs/spec/records.md; its
// owner is this package, and its test compares it with the block.
var TasksSchema = tsv.Schema{Name: "tasks", Location: "records:tasks.tsv", Columns: []tsv.Column{
	{Name: "task", Type: "id(T-xxxx)", Key: true},
	{Name: "issue", Type: "int"},
	{Name: "kind", Type: "enum(specification|architecture|batch|rendered|plan|build)"},
	{Name: "requirements", Type: "list(text)"},
	{Name: "size", Type: "enum(small|large)"},
	{Name: "predecessors", Type: "list(id(T-xxxx))"},
	{Name: "pr", Type: "int"},
	{Name: "opened", Type: "time"},
}}

// The paths of the register and the table in the records.
const (
	tasksPath       = "tasks.tsv"
	transitionsPath = "transitions.tsv"
)

// row gives the fields of a task by TasksSchema, or the refusal of the column
// that the task cannot be written in; line is the line of the row.
func (t Task) row(line int) ([]string, error) {
	list := func(column string, values []string) (string, error) {
		s, err := tsv.JoinList(values)
		if err != nil {
			return "", &tsv.Error{Line: line, Column: column, Reason: err.Error()}
		}
		return s, nil
	}
	reqs, err := list("requirements", t.Requirements)
	if err != nil {
		return nil, err
	}
	preds, err := list("predecessors", t.Predecessors)
	if err != nil {
		return nil, err
	}
	pr := ""
	if t.PR != 0 {
		pr = strconv.Itoa(t.PR)
	}
	return []string{t.ID, strconv.Itoa(t.Issue), t.Kind, reqs, t.Size, preds, pr, stamp(t.Opened)}, nil
}

// taskOf reads the fields of a row that TasksSchema accepted.
func taskOf(f []string) (Task, error) {
	issue, err := strconv.Atoi(f[1])
	if err != nil {
		return Task{}, err
	}
	pr := 0
	if f[6] != "" {
		if pr, err = strconv.Atoi(f[6]); err != nil {
			return Task{}, err
		}
	}
	opened, err := time.Parse(timeForm, f[7])
	if err != nil {
		return Task{}, err
	}
	return Task{ID: f[0], Issue: issue, Kind: f[2], Requirements: words(f[3]), Size: f[4], Predecessors: words(f[5]), PR: pr, Opened: opened}, nil
}

// words gives the values of a list field: none for the empty list, nil.
func words(field string) []string {
	if field == "" {
		return nil
	}
	return strings.Fields(field)
}

// WriteRegister writes the register by TasksSchema. It writes nothing for a
// register that ReadRegister would refuse.
func WriteRegister(w io.Writer, tasks []Task) error {
	rows := make([][]string, len(tasks))
	for i, t := range tasks {
		row, err := t.row(i + 2)
		if err != nil {
			return err
		}
		rows[i] = row
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, TasksSchema, rows); err != nil {
		return err
	}
	if _, err := ReadRegister(b.Bytes()); err != nil {
		return err
	}
	_, err := w.Write(b.Bytes())
	return err
}

// ReadRegister reads a register by TasksSchema and then makes the checks of
// its owner (CheckRegister).
func ReadRegister(data []byte) ([]Task, error) {
	rows, err := tsv.Read(data, TasksSchema)
	if err != nil {
		return nil, err
	}
	tasks := make([]Task, len(rows))
	for i, f := range rows {
		if tasks[i], err = taskOf(f); err != nil {
			return nil, &tsv.Error{Line: i + 2, Reason: err.Error()}
		}
	}
	if err := CheckRegister(tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// CheckRegister makes the checks that docs/spec/records.md gives the owner of
// tasks, after the types: requirements is never —, each predecessor is a task of
// the register, and no cycle of predecessors exists. The refusal names the line
// and the column; a cycle is named at the first row, in the order of the table,
// that is on one.
func CheckRegister(tasks []Task) error {
	index := map[string]int{}
	for i, t := range tasks {
		index[t.ID] = i
	}
	for i, t := range tasks {
		if len(t.Requirements) == 0 {
			return &tsv.Error{Line: i + 2, Column: "requirements", Reason: "a task serves at least one requirement; " + t.ID + " names none"}
		}
		for _, p := range t.Predecessors {
			if _, ok := index[p]; !ok {
				return &tsv.Error{Line: i + 2, Column: "predecessors", Reason: p + " is not a task of the register"}
			}
		}
	}
	for i, t := range tasks {
		if reaches(tasks, index, t.Predecessors, t.ID) {
			return &tsv.Error{Line: i + 2, Column: "predecessors", Reason: "the predecessors of " + t.ID + " form a cycle"}
		}
	}
	return nil
}

// reaches reports whether the task goal is among from or among the
// predecessors, at any depth, of the tasks in from.
func reaches(tasks []Task, index map[string]int, from []string, goal string) bool {
	seen := map[string]bool{}
	for stack := slices.Clone(from); len(stack) > 0; {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if id == goal {
			return true
		}
		if !seen[id] {
			seen[id] = true
			stack = append(stack, tasks[index[id]].Predecessors...)
		}
	}
	return false
}

// AddTask gives the register with t added as its last row, or the refusal of
// CheckRegister for that row; the table given is not changed. A task whose ID is
// a row's is refused: the ID is the key, and a row is replaced by a later commit
// (docs/spec/loop.md, The task register), which is not this call.
func AddTask(tasks []Task, t Task) ([]Task, error) {
	for _, x := range tasks {
		if x.ID == t.ID {
			return nil, &tsv.Error{Line: len(tasks) + 2, Column: "task", Reason: t.ID + " is already a task of the register"}
		}
	}
	next := append(slices.Clone(tasks), t)
	if err := CheckRegister(next); err != nil {
		return nil, err
	}
	return next, nil
}

// registerFiles gives the files of the records commit that opens t: tasks.tsv
// with t added, and the default transition table when the records have none.
// An existing transition table is read and checked, and kept; one that fails
// the check, or a register that ReadRegister refuses, refuses the task and no
// file is given (docs/spec/loop.md, The transition table; the table is written
// in the commit that opens the first task and never again from the default).
func registerFiles(tasks, transitions []byte, t Task) (map[string][]byte, error) {
	var register []Task
	if tasks != nil {
		var err error
		if register, err = ReadRegister(tasks); err != nil {
			return nil, fmt.Errorf("%s: %w", tasksPath, err)
		}
	}
	files := map[string][]byte{}
	if transitions == nil {
		var def bytes.Buffer
		if err := handoff.WriteTable(&def, handoff.Default()); err != nil {
			return nil, err
		}
		files[transitionsPath] = def.Bytes()
	} else {
		rows, err := handoff.ReadTable(transitions)
		if err == nil {
			err = handoff.CheckTable(rows)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", transitionsPath, err)
		}
	}
	register, err := AddTask(register, t)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tasksPath, err)
	}
	var out bytes.Buffer
	if err := WriteRegister(&out, register); err != nil {
		return nil, fmt.Errorf("%s: %w", tasksPath, err)
	}
	files[tasksPath] = out.Bytes()
	return files, nil
}

// loopPermissionsStep checks the permissions of the loop on the installation
// that step forge read: the loop's caller makes it before its first write, the
// issue of a task included (docs/spec/forge.md, The calls of M2e). A missing
// permission is a fail of the step forge that names it. Start and the restart
// record the permissions and do not fail on them, but the steps shape, bet and
// loop of the restart call this check before they do anything and fail their
// own row with the detail forge: … (atBoundary, task T-9nbf, row 66k; D6). A
// state that read no installation holds none.
func (r *state) loopPermissionsStep() Step {
	const name = "forge"
	if missing := forge.CheckLoopPermissions(r.granted); len(missing) > 0 {
		return failf(name, "the installation lacks the permissions of the task loop: %s", strings.Join(missing, ", "))
	}
	return done(name, "the permissions of the task loop")
}

// registerTask opens t in the register: one records commit of tasks.tsv and,
// when the records have none, of the default transition table and of
// parameters.tsv with its defaults (row 66b). It checks the
// permissions first, and writes nothing when the task, the register or the
// table is refused. The step of the caller (row 66k) is not named here.
func (r *state) registerTask(ctx context.Context, t Task) error {
	return r.registerTaskWith(ctx, t, nil)
}

// registerTaskWith is registerTask whose commit also holds the extra files (task T-2vva: the binding of an architecture task to its
// option); a path that the commit already writes is refused.
func (r *state) registerTaskWith(ctx context.Context, t Task, extra map[string][]byte) error {
	if s := r.loopPermissionsStep(); s.Result == "fail" {
		return fmt.Errorf("%s", s.Detail)
	}
	tasks, err := r.store.ReadFile(ctx, tasksPath)
	if err != nil {
		return err
	}
	transitions, err := r.store.ReadFile(ctx, transitionsPath)
	if err != nil {
		return err
	}
	parameters, err := r.store.ReadFile(ctx, parametersPath)
	if err != nil {
		return err
	}
	t.Opened = r.cfg.Clock.Now().UTC().Truncate(time.Second)
	files, err := registerFiles(tasks, transitions, t)
	if err != nil {
		return err
	}
	// parameters.tsv is written in the same commit, when the records have none (task T-dx1h, row 66b).
	more, err := parameterFiles(parameters, r.value, t.Opened)
	if err != nil {
		return err
	}
	maps.Copy(files, more)
	for p, data := range extra {
		if _, clash := files[p]; clash {
			return fmt.Errorf("the extra file %s is a file of the commit that registers %s", p, t.ID)
		}
		files[p] = data
	}
	return r.fenced(ctx, func(ctx context.Context) error {
		return r.store.Commit(ctx, files, "tasks: register "+t.ID)
	})
}
