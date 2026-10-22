package run

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/tsv"
)

// The task register (task T-cak3, #357, row 66a of the plan of M2e; docs/spec/loop.md, The task register;
// docs/spec/records.md, the block tasks and its owner checks).

var opened = time.Date(2026, 10, 16, 7, 0, 0, 0, time.UTC)

func aTask(id string, preds ...string) Task {
	return Task{ID: id, Issue: 10, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Predecessors: preds, Opened: opened}
}

// refusedAt asserts that err is a *tsv.Error at line and column whose reason holds each of words.
func refusedAt(t *testing.T, what string, err error, line int, column string, words ...string) {
	t.Helper()
	var e *tsv.Error
	if !errors.As(err, &e) {
		t.Errorf("%s: %v; want a *tsv.Error naming line %d, column %s", what, err, line, column)
		return
	}
	if e.Line != line || e.Column != column {
		t.Errorf("%s: line %d column %q; want line %d column %q (%v)", what, e.Line, e.Column, line, column, err)
	}
	for _, w := range words {
		if !strings.Contains(e.Reason, w) {
			t.Errorf("%s: the reason %q does not hold %q", what, e.Reason, w)
		}
	}
}

func TestARegisterIsWrittenAndReadBackWithItsDashes(t *testing.T) {
	first := Task{ID: "T-aaaa", Issue: 10, Kind: "specification", Requirements: []string{"REQ-001"}, Opened: opened}
	second := Task{ID: "T-bbbb", Issue: 11, Kind: "build", Requirements: []string{"REQ-005", "NFR-001"}, Size: "large", Predecessors: []string{"T-aaaa"}, PR: 12, Opened: opened.Add(time.Hour)}
	var b bytes.Buffer
	if err := WriteRegister(&b, []Task{first, second}); err != nil {
		t.Fatal(err)
	}
	want := "task\tissue\tkind\trequirements\tsize\tpredecessors\tpr\topened\n" +
		"T-aaaa\t10\tspecification\tREQ-001\t—\t—\t—\t2026-10-16T07:00:00Z\n" +
		"T-bbbb\t11\tbuild\tREQ-005 NFR-001\tlarge\tT-aaaa\t12\t2026-10-16T08:00:00Z\n"
	if b.String() != want {
		t.Errorf("the register:\n%q\nwant\n%q", b.String(), want)
	}
	got, err := ReadRegister(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []Task{first, second}) {
		t.Errorf("the register read back:\n%+v\nwant\n%+v", got, []Task{first, second})
	}
}

func TestARegisterThatReadWouldRefuseIsNotWritten(t *testing.T) {
	for _, c := range []struct {
		name   string
		task   Task
		column string
	}{
		{"no requirement", func() Task { x := aTask("T-aaaa"); x.Requirements = nil; return x }(), "requirements"},
		{"the text — as a requirement reads back as none", func() Task { x := aTask("T-aaaa"); x.Requirements = []string{"—"}; return x }(), "requirements"},
		{"a requirement that holds a space is two", func() Task { x := aTask("T-aaaa"); x.Requirements = []string{"REQ 005"}; return x }(), "requirements"},
		{"a kind outside the enum", func() Task { x := aTask("T-aaaa"); x.Kind = "chore"; return x }(), "kind"},
		{"a task ID outside its form", aTask("T-AAAA"), "task"},
	} {
		var b bytes.Buffer
		err := WriteRegister(&b, []Task{c.task})
		var e *tsv.Error
		if !errors.As(err, &e) || e.Column != c.column || e.Line != 2 {
			t.Errorf("%s: %v; want a refusal at line 2, column %s", c.name, err, c.column)
		}
		if b.Len() != 0 {
			t.Errorf("%s: %q was written", c.name, b.String())
		}
	}
}

func TestARowWithNoRequirementIsRefusedNamingItsLineAndColumn(t *testing.T) {
	table := "task\tissue\tkind\trequirements\tsize\tpredecessors\tpr\topened\n" +
		"T-aaaa\t10\tbuild\tREQ-005\t—\t—\t—\t2026-10-16T07:00:00Z\n" +
		"T-bbbb\t11\tbuild\t—\t—\t—\t—\t2026-10-16T07:00:00Z\n"
	_, err := ReadRegister([]byte(table))
	refusedAt(t, "a table with a row of no requirement", err, 3, "requirements", "requirement")
	none := aTask("T-bbbb")
	none.Requirements = nil
	_, err = AddTask([]Task{aTask("T-aaaa")}, none)
	refusedAt(t, "a task added with no requirement", err, 3, "requirements", "requirement")
}

func TestAPredecessorThatIsNotATaskIsRefusedWhereverTheTaskIsInTheTable(t *testing.T) {
	// A predecessor may come later in the table: the check is of the table, not of the order.
	if err := CheckRegister([]Task{aTask("T-aaaa", "T-cccc"), aTask("T-bbbb"), aTask("T-cccc")}); err != nil {
		t.Errorf("a predecessor that is a later row: %v", err)
	}
	err := CheckRegister([]Task{aTask("T-aaaa"), aTask("T-bbbb", "T-aaaa", "T-zzzz")})
	refusedAt(t, "a predecessor that is no row", err, 3, "predecessors", "T-zzzz")
	_, err = AddTask([]Task{aTask("T-aaaa")}, aTask("T-bbbb", "T-yyyy"))
	refusedAt(t, "a task added with a predecessor that is no task", err, 3, "predecessors", "T-yyyy")
	table := "task\tissue\tkind\trequirements\tsize\tpredecessors\tpr\topened\n" +
		"T-aaaa\t10\tbuild\tREQ-005\t—\tT-qqqq\t—\t2026-10-16T07:00:00Z\n"
	_, err = ReadRegister([]byte(table))
	refusedAt(t, "a table read", err, 2, "predecessors", "T-qqqq")
}

func TestACycleOfPredecessorsIsRefusedAndADiamondIsNot(t *testing.T) {
	for _, c := range []struct {
		name  string
		tasks []Task
		line  int // the first row, in table order, that is on a cycle
	}{
		{"a task that is its own predecessor", []Task{aTask("T-aaaa"), aTask("T-bbbb", "T-bbbb")}, 3},
		{"two tasks", []Task{aTask("T-aaaa", "T-bbbb"), aTask("T-bbbb", "T-aaaa")}, 2},
		{"three tasks", []Task{aTask("T-aaaa", "T-cccc"), aTask("T-bbbb", "T-aaaa"), aTask("T-cccc", "T-bbbb")}, 2},
		{"a cycle behind a tail", []Task{aTask("T-dddd", "T-aaaa"), aTask("T-aaaa", "T-bbbb"), aTask("T-bbbb", "T-aaaa")}, 3},
	} {
		refusedAt(t, c.name, CheckRegister(c.tasks), c.line, "predecessors", "cycle")
	}
	diamond := []Task{aTask("T-aaaa"), aTask("T-bbbb", "T-aaaa"), aTask("T-cccc", "T-aaaa"), aTask("T-dddd", "T-bbbb", "T-cccc")}
	if err := CheckRegister(diamond); err != nil {
		t.Errorf("a diamond of predecessors: %v", err)
	}
	var b bytes.Buffer
	if err := WriteRegister(&b, []Task{aTask("T-aaaa", "T-bbbb"), aTask("T-bbbb", "T-aaaa")}); err == nil {
		t.Error("a register with a cycle was written")
	}
}

func TestAddTaskLeavesTheTableItWasGivenAlone(t *testing.T) {
	base := make([]Task, 1, 4) // spare room, where an append that shares the array would write
	base[0] = aTask("T-aaaa")
	got, err := AddTask(base, aTask("T-bbbb", "T-aaaa"))
	if err != nil || len(got) != 2 || got[1].ID != "T-bbbb" {
		t.Fatalf("AddTask = %+v, %v", got, err)
	}
	if len(base) != 1 || base[:2][1].ID != "" {
		t.Errorf("the table given was changed: %+v, next %+v", base, base[:2][1])
	}
	_, err = AddTask(got, aTask("T-bbbb"))
	refusedAt(t, "a task with the ID of a row", err, 4, "task", "T-bbbb")
	if len(got) != 2 {
		t.Errorf("a refused task changed the table: %+v", got)
	}
}

func TestTheFilesOfTheFirstTaskAreTheRegisterAndTheDefaultTable(t *testing.T) {
	first := aTask("T-aaaa")
	files, err := registerFiles(nil, nil, first)
	if err != nil {
		t.Fatal(err)
	}
	var def bytes.Buffer
	if err := handoff.WriteTable(&def, handoff.Default()); err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || !bytes.Equal(files["transitions.tsv"], def.Bytes()) {
		t.Errorf("the files of the first task: %d files, transitions.tsv %q; want tasks.tsv and the default table", len(files), files["transitions.tsv"])
	}
	rows, err := ReadRegister(files["tasks.tsv"])
	if err != nil || !reflect.DeepEqual(rows, []Task{first}) {
		t.Errorf("tasks.tsv holds %+v, %v; want the task", rows, err)
	}
	// Another task: the register has the two rows and the table is not written again.
	later, err := registerFiles(files["tasks.tsv"], files["transitions.tsv"], aTask("T-bbbb", "T-aaaa"))
	if err != nil {
		t.Fatal(err)
	}
	if _, has := later["transitions.tsv"]; has || len(later) != 1 {
		t.Errorf("the files of a later task: %d, with transitions.tsv %v; want tasks.tsv alone", len(later), has)
	}
	rows, err = ReadRegister(later["tasks.tsv"])
	if err != nil || len(rows) != 2 || rows[1].ID != "T-bbbb" {
		t.Errorf("tasks.tsv holds %+v, %v; want the two tasks", rows, err)
	}
}

// transitionBytes gives the bytes of a transition table.
func transitionBytes(t *testing.T, rows []handoff.Row) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := handoff.WriteTable(&b, rows); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestAnExistingTransitionTableIsKeptAndAnInvalidOneRefusesTheTask(t *testing.T) {
	// The Operator's table: the default and a pair of its own.
	custom := append(handoff.Default(), handoff.Row{From: "verifier", To: "asking role", Artifacts: []string{"question"}, Items: []string{"schema"}})
	files, err := registerFiles(nil, transitionBytes(t, custom), aTask("T-aaaa"))
	if err != nil {
		t.Fatal(err)
	}
	if _, has := files["transitions.tsv"]; has || len(files) != 1 {
		t.Errorf("an existing table: the files are %d, with transitions.tsv %v; want tasks.tsv alone", len(files), has)
	}
	missing := handoff.Default()[1:] // no row for the first default pair
	badItem := handoff.Default()
	badItem[0].Items = []string{"nonsense"}
	for _, c := range []struct {
		name string
		data []byte
		want string
	}{
		{"a table with no row for a default pair", transitionBytes(t, missing), "task plan"},
		{"a table with an item outside the list", func() []byte {
			return []byte(strings.Replace(string(transitionBytes(t, handoff.Default())), "schema artifacts plan-fields", "schema nonsense plan-fields", 1))
		}(), "nonsense"},
		{"a file that is no table", []byte("not a table"), "transitions.tsv"},
	} {
		files, err := registerFiles(nil, c.data, aTask("T-aaaa"))
		if err == nil || !strings.Contains(err.Error(), "transitions.tsv") || !strings.Contains(err.Error(), c.want) || files != nil {
			t.Errorf("%s: %v, files %v; want a refusal that names transitions.tsv and %q, and no file", c.name, err, files, c.want)
		}
	}
}

// Task T-t94t (#483), condition 4 of the plan review: registration reads a table written before the pair plan review to task plan had a row
// (InForce), keeps its bytes, and still refuses a table that lacks an original pair.
func TestARegistrationKeepsATableWrittenBeforeALaterPair(t *testing.T) {
	old := transitionBytes(t, original13())
	files, err := registerFiles(nil, old, aTask("T-aaaa"))
	if err != nil {
		t.Fatalf("the table of 13 rows: %v", err)
	}
	if _, has := files["transitions.tsv"]; has || len(files) != 1 {
		t.Errorf("the table of 13 rows: %d files, with transitions.tsv %v; want tasks.tsv alone, the table's bytes kept", len(files), has)
	}
	if files, err := registerFiles(nil, transitionBytes(t, original13()[1:]), aTask("T-aaaa")); err == nil || !strings.Contains(err.Error(), "task plan -> plan review") || files != nil {
		t.Errorf("a table of 12 rows lacking an original pair: %v, files %v; want the pair named and no file", err, files)
	}
	for _, rows := range [][]handoff.Row{append(original13(), handoff.Row{From: "plan review", To: "task plan", Artifacts: []string{"plan-review"}, Items: []string{"schema", "artifacts"}})} {
		if _, err := registerFiles(nil, transitionBytes(t, rows), aTask("T-aaaa")); err != nil {
			t.Errorf("a customized row 14: %v; want it kept", err)
		}
	}
}

func TestAnExistingRegisterThatReadRefusesStopsTheWrite(t *testing.T) {
	cycle := handWritten(t, []Task{aTask("T-aaaa", "T-bbbb"), aTask("T-bbbb", "T-aaaa")})
	for _, c := range []struct {
		name string
		data []byte
		want string
	}{
		{"a register with a cycle", cycle, "cycle"},
		{"a file that is no register", []byte("not a register"), "tasks.tsv"},
	} {
		files, err := registerFiles(c.data, nil, aTask("T-cccc"))
		if err == nil || !strings.Contains(err.Error(), "tasks.tsv") || !strings.Contains(err.Error(), c.want) || files != nil {
			t.Errorf("%s: %v, files %v; want a refusal that names tasks.tsv and %q, and no file", c.name, err, files, c.want)
		}
	}
}

// handWritten gives the bytes of a register with no check, as a hand-written one is.
func handWritten(t *testing.T, tasks []Task) []byte {
	t.Helper()
	var b bytes.Buffer
	b.WriteString("task\tissue\tkind\trequirements\tsize\tpredecessors\tpr\topened\n")
	for _, x := range tasks {
		pred := "—"
		if len(x.Predecessors) > 0 {
			pred = strings.Join(x.Predecessors, " ")
		}
		b.WriteString(x.ID + "\t10\tbuild\t" + strings.Join(x.Requirements, " ") + "\t—\t" + pred + "\t—\t2026-10-16T07:00:00Z\n")
	}
	return b.Bytes()
}

func TestTheStepForgeOfStartAndOfTheRestartRecordsTheLoopPermissionsAndDoesNotFailOnThem(t *testing.T) {
	ctx := context.Background()
	f := goodForge() // the permissions of M2a alone, and statuses read
	r := &state{cfg: startConfig(f)}
	if s := r.forgeStep(ctx, true); s.Result != "done" {
		t.Fatalf("Start with an App that lacks the permissions of the loop: %+v", s)
	}
	if !reflect.DeepEqual(r.granted, f.perms) {
		t.Errorf("the grant %v; want the permissions of the installation %v", r.granted, f.perms)
	}
	f.repo.Branches, f.repo.HasCommit = []string{"layup-records", "main"}, true
	r = &state{cfg: startConfig(f)}
	if s := r.forgeStep(ctx, false); s.Result != "done" || !reflect.DeepEqual(r.granted, f.perms) {
		t.Errorf("the restart with an App that lacks them: %+v, grant %v", s, r.granted)
	}
	if strings.Contains(r.permissions, "pull_requests") {
		t.Errorf("app.permissions %q names a permission the installation lacks", r.permissions)
	}
	// The grant is a copy: a later change of the installation's map does not change it.
	f.perms["pull_requests"] = "write"
	if _, has := r.granted["pull_requests"]; has {
		t.Error("the grant shares the map of the installation")
	}
}

func TestTheLoopPermissionsStepNamesEachMissingPermission(t *testing.T) {
	const prefix = "the installation lacks the permissions of the task loop: "
	full := map[string]string{"metadata": "read", "contents": "write", "issues": "write", "pull_requests": "write", "statuses": "write", "checks": "read"}
	without := func(drop string, set map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range full {
			if k != drop {
				m[k] = v
			}
		}
		for k, v := range set {
			m[k] = v
		}
		return m
	}
	m2a := map[string]string{"metadata": "read", "contents": "write", "issues": "write"}
	for _, c := range []struct {
		name    string
		granted map[string]string
		want    []string // none: done
	}{
		{"all three", full, nil},
		{"checks write holds read", without("checks", map[string]string{"checks": "write"}), nil},
		{"no pull requests", without("pull_requests", nil), []string{"pull_requests: write"}},
		{"no statuses", without("statuses", nil), []string{"statuses: write"}},
		{"no checks", without("checks", nil), []string{"checks: read"}},
		{"pull requests read is not write", without("pull_requests", map[string]string{"pull_requests": "read"}), []string{"pull_requests: write"}},
		{"statuses read is not write", without("statuses", map[string]string{"statuses": "read"}), []string{"statuses: write"}},
		{"the permissions of M2a alone", m2a, []string{"pull_requests: write", "statuses: write", "checks: read"}},
		{"a state that read no installation", nil, []string{"pull_requests: write", "statuses: write", "checks: read"}},
	} {
		s := (&state{granted: c.granted}).loopPermissionsStep()
		if s.Name != "forge" {
			t.Errorf("%s: the step is named %q; want forge", c.name, s.Name)
		}
		switch {
		case c.want == nil && s.Result != "done":
			t.Errorf("%s: %+v; want done", c.name, s)
		case c.want != nil && (s.Result != "fail" || s.Detail != prefix+strings.Join(c.want, ", ")):
			t.Errorf("%s: %+v; want fail with %q", c.name, s, prefix+strings.Join(c.want, ", "))
		}
	}
}

// A table of the Operator whose row of code's own handoff lists schema can satisfy no rendered task, so the check of
// a table refuses it where the register reads it (the comment of #358; task T-dx1h, row 66b).
func TestAnExistingTransitionTableThatListsSchemaForTheRenderedTaskRefusesTheTask(t *testing.T) {
	custom := handoff.Default()
	for i := range custom {
		if custom[i].From == "code (rendered task)" {
			custom[i].Items = []string{"schema", "artifacts", "rendering-equal"}
		}
	}
	files, err := registerFiles(nil, transitionBytes(t, custom), aTask("T-aaaa"))
	if err == nil || !strings.Contains(err.Error(), "code (rendered task)") || files != nil {
		t.Errorf("registerFiles: %v, %v; want a refusal that names the row and no file", files, err)
	}
}
