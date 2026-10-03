package verify

import (
	"errors"
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

const fixtureCommit = "f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1f1"

// gateRecord is a record with the stack and pin.time that a fixture commit
// reads.
var gateRecord = work.Record{{"S01", "stack", "go", "answer", "S01-stack"}, {"S02", "pin.time", "2026-09-30T08:00:00Z", "computed", "the clock"}}

// standInGate gives the result of the kind static on the clean head and on
// the fixture commit, and logs each run.
func standInGate(t *testing.T, clean, fixture gate.Row, log *[]string) {
	t.Helper()
	saved, savedFixture := gateRun, fixtureOf
	t.Cleanup(func() { gateRun, fixtureOf = saved, savedFixture })
	gateRun = func(repo, base, head string, _ func(int, int, string) func(), _ io.Writer) (gate.Table, error) {
		*log = append(*log, "gate "+repo+" "+base+" "+head)
		r := clean
		if head == fixtureCommit {
			r = fixture
		}
		if r.Result == "" {
			return gate.Table{}, &gate.InputError{Err: errors.New("sh not found\nmore")}
		}
		return gate.Table{Rows: []gate.Row{r, {Kind: "layout", State: "pending", Result: "clear", Reason: "pending: no product path"}}}, nil
	}
	fixtureOf = func(stack, kind string) ([]byte, error) {
		if stack != "go" || kind != "static" {
			return nil, errors.New("catalog: go has no active kind \"" + kind + "\"")
		}
		return []byte("the patch"), nil
	}
}

// The rows gate:<kind> follow the table of setup.md, the first rule that
// matches deciding; a pending kind is clear with no run; the clean run is made
// once per run; a fixture that does not apply is not-active (D3 of #92, with
// condition 1 of its plan review).
func TestTheGateRows(t *testing.T) {
	row := func(result, reason string) gate.Row {
		return gate.Row{Kind: "static", State: "active", Result: result, Reason: reason}
	}
	for _, c := range []struct {
		name           string
		clean, fixture gate.Row
		applyErr       error
		want           Row
		runs           int
	}{
		{"detected", row("clear", "no product path"), row("fail", "exit 1"), nil, Row{"gate:static", "pass", ""}, 2},
		{"a pass, then detected", row("pass", "—"), row("fail", "exit 1"), nil, Row{"gate:static", "pass", ""}, 2},
		{"a clean fail", row("fail", "exit 1"), row("not-active", "tool not found: go"), nil, Row{"gate:static", "fail", "the clean run: exit 1"}, 1},
		{"a clean not-active", row("not-active", "tool not found: go"), row("fail", "exit 1"), nil, Row{"gate:static", "not-active", "the clean run: tool not found: go"}, 1},
		{"a fixture not-active", row("clear", "no product path"), row("not-active", "tool not found: go"), nil,
			Row{"gate:static", "not-active", "the fixture run: tool not found: go"}, 2},
		{"not detected", row("clear", "no product path"), row("pass", "—"), nil, Row{"gate:static", "fail", "fixture not detected"}, 2},
		{"not detected, clear", row("pass", "—"), row("clear", "no product path"), nil, Row{"gate:static", "fail", "fixture not detected"}, 2},
		{"a clean run with no table", gate.Row{}, row("fail", "exit 1"), nil, Row{"gate:static", "not-active", "the clean run: sh not found"}, 1},
		{"a fixture run with no table", row("clear", "no product path"), gate.Row{}, nil, Row{"gate:static", "not-active", "the fixture run: sh not found"}, 2},
		{"a fixture that does not apply", row("clear", "no product path"), row("fail", "exit 1"),
			&git.FailedError{Args: []string{"apply", "--", "/tmp/x/fixture.patch"}, Code: 1, Stderr: "error: gatefixture/static.go: already exists in working directory\nmore\n", Err: errors.New("exit status 1")},
			Row{"gate:static", "not-active", "fixture does not apply: error: gatefixture/static.go: already exists in working directory"}, 1},
		{"an apply with no output", row("clear", "no product path"), row("fail", "exit 1"), errors.New("/tmp/x: no such file"),
			Row{"gate:static", "not-active", "fixture does not apply: git apply failed"}, 1},
	} {
		var runs []string
		g := goodGit()
		g.applyErr = c.applyErr
		standIn(t, g, nil, nil, nil)
		standInGate(t, c.clean, c.fixture, &runs)
		rows := &gateRows{target: "/w/target", head: commitA, record: gateRecord, out: io.Discard}
		if got := rows.row(gate.Kind{Name: "static", State: "active"}); !reflect.DeepEqual(got, c.want) || len(runs) != c.runs {
			t.Errorf("%s: %q after %d runs; want %q after %d", c.name, got, len(runs), c.want, c.runs)
		}
		if c.runs == 2 && len(runs) == 2 && runs[1] != "gate /w/target "+commitA+" "+fixtureCommit {
			t.Errorf("%s: the fixture run %q; want the fixture commit as the head, and the setup head as the base", c.name, runs[1])
		}
	}

	var runs []string
	g := goodGit()
	removed, _ := standIn(t, g, nil, nil, nil)
	standInGate(t, row("clear", "no product path"), row("fail", "exit 1"), &runs)
	rows := &gateRows{target: "/w/target", head: commitA, record: gateRecord, out: io.Discard}
	if got := rows.row(gate.Kind{Name: "layout", State: "pending"}); got != (Row{"gate:layout", "clear", "pending: fixture not run"}) || len(runs) != 0 || len(g.calls) != 0 {
		t.Errorf("a pending kind: %q, runs %q, calls %q; want clear with no run", got, runs, g.calls)
	}
	rows.row(gate.Kind{Name: "static", State: "active"})
	rows.row(gate.Kind{Name: "static", State: "active"})
	if strings.Count(strings.Join(runs, "\n"), commitA+" "+commitA) != 1 {
		t.Errorf("the runs %q; want one clean run for the run", runs)
	}
	want := []string{"worktree add /w/target /stand-in/scratch/tree " + commitA, "apply /stand-in/scratch/tree /stand-in/scratch/fixture.patch",
		"commit /stand-in/scratch/tree \"test: the known-bad fixture of the kind static\" by layup setup verify <verify@layup.invalid> at 2026-09-30T08:00:00Z",
		"rev-parse /stand-in/scratch/tree HEAD^{commit}", "worktree remove /w/target /stand-in/scratch/tree"}
	if len(g.calls) < 5 || !reflect.DeepEqual(g.calls[:5], want) || !slices.Equal(*removed, []string{"/stand-in/scratch", "/stand-in/scratch-2"}) {
		t.Errorf("the calls of git %q, removed %q; want the fixture commit in a scratch tree that is removed: %q", g.calls, *removed, want)
	}
	for _, c := range []struct {
		name   string
		record work.Record
		want   string
	}{
		{"no stack", gateRecord[1:], "record: no value at S01 stack"},
		{"no pin.time", gateRecord[:1], "record: no value at S02 pin.time"},
		{"a kind with no fixture", append(work.Record{{"S01", "stack", "rust", "answer", "S01-stack"}}, gateRecord[1:]...), "fixture: catalog: go has no active kind"},
	} {
		standIn(t, goodGit(), nil, nil, nil)
		standInGate(t, row("clear", "no product path"), row("fail", "exit 1"), new([]string))
		rows := &gateRows{target: "/w/target", head: commitA, record: c.record, out: io.Discard}
		if got := rows.row(gate.Kind{Name: "static", State: "active"}); got.Result != "not-active" || !strings.HasPrefix(got.Reason, c.want) {
			t.Errorf("%s: %q; want not-active, %s", c.name, got, c.want)
		}
	}
}
