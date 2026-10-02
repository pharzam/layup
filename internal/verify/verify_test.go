package verify

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/work"
)

const twoKinds = "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
	"static\tactive\tgo\tgo vet ./...\t./*.go\t—\n" +
	"layout\tpending\tgo\tgo test ./layout/\t./*.go\t—\n"

// standInGit answers the calls of a run, and records them.
type standInGit struct {
	version           string // "" is no git
	head              string // the commit of layup-setup; "" is no such branch
	manifest          []byte // nil is no docs/gates.tsv at the head
	addErr, removeErr error
	calls             []string
}

func (g *standInGit) Version() (string, error) {
	if g.version == "" {
		return "", errors.New("git not found")
	}
	return g.version, nil
}

func (g *standInGit) RevParse(dir, rev string) (string, error) {
	g.calls = append(g.calls, "rev-parse "+dir+" "+rev)
	if g.head == "" {
		return "", errors.New("fatal: Needed a single revision\nmore text")
	}
	return g.head, nil
}

func (g *standInGit) Show(dir, rev, path string) ([]byte, error) {
	g.calls = append(g.calls, "show "+dir+" "+rev+":"+path)
	if g.manifest == nil {
		return nil, errors.New("fatal: path 'docs/gates.tsv' does not exist in '" + rev + "'")
	}
	return g.manifest, nil
}

func (g *standInGit) WorktreeAdd(dir, path, rev string) error {
	g.calls = append(g.calls, "worktree add "+dir+" "+path+" "+rev)
	return g.addErr
}

func (g *standInGit) WorktreeRemove(dir, path string) error {
	g.calls = append(g.calls, "worktree remove "+dir+" "+path)
	return g.removeErr
}

func (g *standInGit) history(dir string) history { return oneRoot }

var standInRecord = work.Record{{"S01", "name", "acme", "answer", "S01-name"}}

// standIn puts stand-ins for git, the records, the scratch directory and the
// built checks in place until the test ends. The built checks keep their
// order and give findings; it gives the removed paths and the record that
// each check got.
func standIn(t *testing.T, g *standInGit, answersErr, recordErr error, findings map[string][]string) (*[]string, map[string]work.Record) {
	t.Helper()
	savedAPI, savedAnswers, savedRecord, savedTemp, savedRoot, savedRemove, savedChecks := repoAPI, readAnswers, readRecord, tempDir, tempRoot, removeAll, checks
	t.Cleanup(func() {
		repoAPI, readAnswers, readRecord, tempDir, tempRoot, removeAll, checks = savedAPI, savedAnswers, savedRecord, savedTemp, savedRoot, savedRemove, savedChecks
	})
	tempRoot = func() string { return "/stand-in" }
	repoAPI = g
	readAnswers = func(dir string) (work.Answers, error) { return work.Answers{}, answersErr }
	readRecord = func(dir string) (work.Record, error) { return standInRecord, recordErr }
	tempDir = func() (string, error) { return "/stand-in/scratch", nil }
	removed := &[]string{}
	removeAll = func(path string) error { *removed = append(*removed, path); return nil }
	got := map[string]work.Record{}
	var table []check
	for _, c := range savedChecks {
		if c.run != nil {
			name, f := c.name, findings[c.name]
			c.run = func(in input) []string { got[name] = in.record; return f }
		}
		table = append(table, c)
	}
	checks = table
	return removed, got
}

func goodGit() *standInGit {
	return &standInGit{version: "2.54.0 (Apple Git-157)", head: commitA, manifest: []byte(twoKinds)}
}

// steps records the progress calls of a run: each start, and each end.
func steps(log *[]string) func(i, n int, check string) func() {
	return func(i, n int, check string) func() {
		*log = append(*log, fmt.Sprintf("%d/%d %s", i, n, check))
		return func() { *log = append(*log, "end") }
	}
}

func notBuiltRow(check string) Row { return Row{check, "not-active", "not built yet"} }

// The rows of layup setup verify are the table of setup.md in its order, then
// one row gate:<kind> per kind of the manifest at the setup head (D4 of #84).
// A check that this version does not have is not-active, so the command does
// not give exit 0 (NFR-004); a failed check gives its first finding.
func TestRunGivesEachRowInTheOrderOfTheTable(t *testing.T) {
	g := goodGit()
	removed, got := standIn(t, g, nil, nil, map[string][]string{"kit-history": {"orphan: x", "kit-link: y"}})
	var log []string
	tbl, err := Run("/w", steps(&log), io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	want := []Row{notBuiltRow("discipline-tests"), {"pin", "pass", ""}, {"kit-history", "fail", "orphan: x"},
		notBuiltRow("facts"), notBuiltRow("onboarding"), notBuiltRow("glossary"), notBuiltRow("guardrails"),
		notBuiltRow("markers"), notBuiltRow("adapted"), {"identity", "pass", ""}, notBuiltRow("link-lint"),
		notBuiltRow("sources"), notBuiltRow("jobs"), notBuiltRow("gate:static"), notBuiltRow("gate:layout")}
	if !reflect.DeepEqual(tbl.Rows, want) {
		t.Errorf("the rows\n got %q\nwant %q", tbl.Rows, want)
	}
	if len(log) != 30 || log[0] != "1/15 discipline-tests" || log[1] != "end" || log[28] != "15/15 gate:layout" || log[29] != "end" {
		t.Errorf("the progress calls %q; want a start and an end for each of the 15 rows", log)
	}
	wantCalls := []string{"rev-parse /w/target refs/heads/layup-setup^{commit}", "show /w/target " + commitA + ":docs/gates.tsv",
		"worktree add /w/target /stand-in/scratch/tree " + commitA, "worktree remove /w/target /stand-in/scratch/tree"}
	if !reflect.DeepEqual(g.calls, wantCalls) {
		t.Errorf("the calls of git\n got %q\nwant %q", g.calls, wantCalls)
	}
	if !reflect.DeepEqual(*removed, []string{"/stand-in/scratch"}) {
		t.Errorf("removed %q; want the scratch directory", *removed)
	}
	for _, name := range []string{"pin", "kit-history", "identity"} {
		if !reflect.DeepEqual(got[name], standInRecord) {
			t.Errorf("the check %s got the record %q; want the record of the work area", name, got[name])
		}
	}
}

// The one-check call runs the named checks in the order of the table; gates
// is each row gate:<kind>, and only it reads the manifest (D5 of #84).
func TestCheckRunsTheNamedChecks(t *testing.T) {
	for _, c := range []struct {
		names []string
		rows  []string
		show  bool
	}{
		{[]string{"identity", "pin"}, []string{"pin", "identity"}, false},
		{[]string{"pin", "pin"}, []string{"pin"}, false},
		{[]string{"gates", "jobs"}, []string{"jobs", "gate:static", "gate:layout"}, true},
		{[]string{"gates"}, []string{"gate:static", "gate:layout"}, true},
	} {
		g := goodGit()
		standIn(t, g, nil, nil, nil)
		tbl, err := Check("/w", c.names, steps(new([]string)), io.Discard)
		var rows []string
		for _, r := range tbl.Rows {
			rows = append(rows, r.Check)
		}
		if err != nil || !reflect.DeepEqual(rows, c.rows) {
			t.Errorf("Check(%q): %q, %v; want %q", c.names, rows, err, c.rows)
		}
		if read := strings.Contains(strings.Join(g.calls, "\n"), "show "); read != c.show {
			t.Errorf("Check(%q) read the manifest: %v, want %v", c.names, read, c.show)
		}
	}
	for _, names := range [][]string{{"pin", "frobnicate"}, {"gate:static"}, nil} {
		g := goodGit()
		standIn(t, g, nil, nil, nil)
		var in *InputError
		if _, err := Check("/w", names, steps(new([]string)), io.Discard); err == nil || errors.As(err, &in) || len(g.calls) != 0 {
			t.Errorf("Check(%q): %v, calls %q; want an error that is not an input error, before any call of git", names, err, g.calls)
		}
	}
}

// Each input error comes before the scratch tree, with no table (exit 2).
func TestTheInputErrors(t *testing.T) {
	recordErr := errors.New("out/record.tsv: line 2, column \"source\": not a value of enum(answer|catalog|fact|computed|gap|step)")
	for _, c := range []struct {
		name                  string
		g                     *standInGit
		answersErr, recordErr error
		want                  string
	}{
		{"an old git", &standInGit{version: "2.31.0", head: commitA, manifest: []byte(twoKinds)}, nil, nil, "git 2.32.0 or newer is needed"},
		{"no git", &standInGit{head: commitA, manifest: []byte(twoKinds)}, nil, nil, "git 2.32.0 or newer is needed: git not found"},
		{"no answers", goodGit(), errors.New("inputs/answers.tsv: open w/inputs/answers.tsv: no such file or directory"), nil, "inputs/answers.tsv: open"},
		{"a record of another form", goodGit(), nil, recordErr, recordErr.Error()},
		{"no branch layup-setup", &standInGit{version: "2.54.0", manifest: []byte(twoKinds)}, nil, nil,
			"target: no commit at the branch layup-setup: fatal: Needed a single revision"},
		{"no manifest", &standInGit{version: "2.54.0", head: commitA}, nil, nil,
			"docs/gates.tsv at the head of layup-setup: fatal: path 'docs/gates.tsv' does not exist"},
		{"a manifest of another form", &standInGit{version: "2.54.0", head: commitA, manifest: []byte("kind\tstate\n")}, nil, nil,
			"docs/gates.tsv at the head of layup-setup: line 1"},
		{"a manifest with no row", &standInGit{version: "2.54.0", head: commitA, manifest: []byte("kind\tstate\ttool\tcommand\tscope\tconfig\n")}, nil, nil,
			"docs/gates.tsv at the head of layup-setup: no kind"},
	} {
		standIn(t, c.g, c.answersErr, c.recordErr, nil)
		tbl, err := Run("/w", steps(new([]string)), io.Discard)
		var in *InputError
		if !errors.As(err, &in) || !strings.HasPrefix(err.Error(), c.want) || len(tbl.Rows) != 0 {
			t.Errorf("%s: %v, %d rows; want an input error that starts with %q, and no row", c.name, err, len(tbl.Rows), c.want)
		}
		if strings.Contains(strings.Join(c.g.calls, "\n"), "worktree") {
			t.Errorf("%s: the calls %q; want no scratch tree", c.name, c.g.calls)
		}
		if err != nil && strings.Contains(err.Error(), "more text") {
			t.Errorf("%s: %v; want only the first line of the error of git", c.name, err)
		}
	}
	// The one-check call of a step reads no manifest, so none is needed.
	standIn(t, &standInGit{version: "2.54.0", head: commitA}, nil, nil, nil)
	if _, err := Check("/w", []string{"pin"}, steps(new([]string)), io.Discard); err != nil {
		t.Errorf("Check(pin) with no manifest: %v; want no error", err)
	}
}

// A scratch tree that cannot be added makes each built check not-active with
// a fixed reason, and the error goes to out; one that cannot be removed gives
// the whole table and a *CleanupError (exit 2), as layup gate does.
func TestTheScratchTree(t *testing.T) {
	g := goodGit()
	g.addErr = errors.New("fatal: '/stand-in/scratch/tree' already exists")
	removed, _ := standIn(t, g, nil, nil, nil)
	var out bytes.Buffer
	tbl, err := Run("/w", steps(new([]string)), &out)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range tbl.Rows {
		want := "not built yet"
		if r.Check == "pin" || r.Check == "kit-history" || r.Check == "identity" {
			want = "scratch tree: add failed"
		}
		if r.Result != "not-active" || r.Reason != want {
			t.Errorf("the row %q; want not-active, %s", r, want)
		}
	}
	if !strings.Contains(out.String(), "already exists") || strings.Contains(strings.Join(g.calls, "\n"), "worktree remove") ||
		!reflect.DeepEqual(*removed, []string{"/stand-in/scratch"}) {
		t.Errorf("out %q, calls %q, removed %q; want the error on out, no removal of a tree, and the scratch directory removed", out.String(), g.calls, *removed)
	}

	g = goodGit()
	g.removeErr = errors.New("fatal: cannot remove")
	standIn(t, g, nil, nil, nil)
	tbl, err = Run("/w", steps(new([]string)), io.Discard)
	var cleanup *CleanupError
	if !errors.As(err, &cleanup) || cleanup.Path != "/stand-in/scratch" || len(tbl.Rows) != 15 {
		t.Errorf("a tree that cannot be removed: %v, %d rows; want a *CleanupError for /stand-in/scratch and the whole table", err, len(tbl.Rows))
	}

	// Round 1, finding 3: a failed add whose directory cannot be removed
	// leaves a scratch directory, so it is a *CleanupError too.
	g = goodGit()
	g.addErr = errors.New("fatal: cannot add")
	standIn(t, g, nil, nil, nil)
	removeAll = func(string) error { return errors.New("permission denied") }
	tbl, err = Run("/w", steps(new([]string)), io.Discard)
	if !errors.As(err, &cleanup) || cleanup.Path != "/stand-in/scratch" || len(tbl.Rows) != 15 {
		t.Errorf("a failed add that leaves its directory: %v, %d rows; want a *CleanupError and the whole table", err, len(tbl.Rows))
	}

	// Round 1, finding 2: a temporary directory in the work area is an input
	// error, before any scratch tree.
	g = goodGit()
	standIn(t, g, nil, nil, nil)
	tempRoot = func() string { return "/w/out" }
	var in *InputError
	if _, err := Run("/w", steps(new([]string)), io.Discard); !errors.As(err, &in) || !strings.Contains(err.Error(), "/w/out") ||
		strings.Contains(strings.Join(g.calls, "\n"), "worktree") {
		t.Errorf("TMPDIR in the work area: %v, calls %q; want an input error before any scratch tree", err, g.calls)
	}
	if _, err := Check("/w", []string{"jobs", "gates"}, steps(new([]string)), io.Discard); err != nil {
		t.Errorf("TMPDIR in the work area, no built check: %v; want no error, because no scratch tree is needed", err)
	}
}

// Round 1, finding 5: the scratch tree is added in the step of the first
// built check and removed in the step of the last row, so the progress lines
// cover both; a call with no built check makes no scratch tree.
func TestTheProgressLinesCoverTheScratchTree(t *testing.T) {
	g := goodGit()
	standIn(t, g, nil, nil, nil)
	if _, err := Run("/w", steps(&g.calls), io.Discard); err != nil {
		t.Fatal(err)
	}
	at := func(prefix string) int {
		return slices.IndexFunc(g.calls, func(c string) bool { return strings.HasPrefix(c, prefix) })
	}
	if add, pin := at("worktree add"), at("2/15 pin"); add < pin || add > at("3/15 kit-history") {
		t.Errorf("the calls %q; want the worktree add inside the step of pin", g.calls)
	}
	if remove, last := at("worktree remove"), at("15/15 gate:layout"); remove < last || g.calls[len(g.calls)-1] != "end" {
		t.Errorf("the calls %q; want the worktree remove inside the step of the last row", g.calls)
	}
	g = goodGit()
	standIn(t, g, nil, nil, nil)
	if _, err := Check("/w", []string{"jobs", "gates"}, steps(new([]string)), io.Discard); err != nil || strings.Contains(strings.Join(g.calls, "\n"), "worktree") {
		t.Errorf("Check(jobs, gates): %v, calls %q; want no scratch tree", err, g.calls)
	}
}

// The table is a record of the block setup-verify: — for the reason of a
// pass.
func TestTheTable(t *testing.T) {
	tbl := Table{Rows: []Row{{"pin", "pass", ""}, {"kit-history", "fail", "orphan: docs/tasks/T-1.md has no line"}, notBuiltRow("jobs")}}
	var b bytes.Buffer
	if err := tbl.Write(&b); err != nil {
		t.Fatal(err)
	}
	want := "check\tresult\treason\npin\tpass\t—\nkit-history\tfail\torphan: docs/tasks/T-1.md has no line\njobs\tnot-active\tnot built yet\n"
	if b.String() != want {
		t.Errorf("the table\n%s\nwant\n%s", b.String(), want)
	}
	if got := tbl.Results(); !reflect.DeepEqual(got, []string{"pass", "fail", "not-active"}) {
		t.Errorf("Results() = %q", got)
	}
}
