//go:build integration

// The run with a fake of internal/git. The fake writes the head into a real
// temporary directory, so these tests are at the integration level (review
// round 1, note 4; docs/tests/test-levels.md: a unit test touches no file).

package gate

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
)

const (
	baseID = "1111111111111111111111111111111111111111"
	headID = "2222222222222222222222222222222222222222"
)

// fakeGit stands in for internal/git: revisions, files and trees by
// "<rev>:<path>", the changed paths, and a WorktreeAdd that writes the files
// of the head into the tree.
type fakeGit struct {
	version          string
	commits          map[string]string
	files            map[string]string
	trees            map[string][]git.TreeEntry
	changed          []string
	head             map[string]string
	link, linkTarget string // a path of the head that is a symbolic link out of the tree, and its target
	diffErr, addErr  error
	removeErr        error
	added, removed   []string
}

func (f *fakeGit) Version() (string, error) {
	if f.version == "" {
		return "", &git.NotFoundError{Err: errors.New("executable file not found in $PATH")}
	}
	return f.version, nil
}

func (f *fakeGit) RevParse(dir, rev string) (string, error) {
	if id, ok := f.commits[rev]; ok {
		return id, nil
	}
	return "", &git.FailedError{Args: []string{"rev-parse"}, Code: 128, Stderr: "fatal: Needed a single revision", Err: errors.New("exit status 128")}
}

func (f *fakeGit) Show(dir, rev, p string) ([]byte, error) {
	if b, ok := f.files[rev+":"+p]; ok {
		return []byte(b), nil
	}
	return nil, &git.FailedError{Args: []string{"show"}, Code: 128, Stderr: "fatal: path '" + p + "' does not exist", Err: errors.New("exit status 128")}
}

func (f *fakeGit) LsTree(dir, rev, p string) ([]git.TreeEntry, error) { return f.trees[rev+":"+p], nil }

func (f *fakeGit) DiffNames(dir, base, head string) ([]string, error) { return f.changed, f.diffErr }

func (f *fakeGit) WorktreeAdd(dir, p, rev string) error {
	f.added = append(f.added, p)
	if f.addErr != nil {
		return f.addErr
	}
	for name, text := range f.head {
		mode := os.FileMode(0o644)
		if strings.HasPrefix(text, "#!") {
			mode = 0o755
		}
		os.MkdirAll(filepath.Dir(filepath.Join(p, name)), 0o755)
		os.WriteFile(filepath.Join(p, name), []byte(text), mode)
	}
	if f.link != "" {
		os.Symlink(f.linkTarget, filepath.Join(p, f.link))
	}
	return os.WriteFile(filepath.Join(p, ".git"), []byte("gitdir: elsewhere\n"), 0o644)
}

func (f *fakeGit) WorktreeRemove(dir, p string) error {
	f.removed = append(f.removed, p)
	if f.removeErr != nil {
		return f.removeErr
	}
	return os.RemoveAll(p)
}

// manifest is a manifest of the base with one kind of each kind of result.
const manifest = manifestHeader +
	"pass\tactive\tgo\tok\t./*.go\t—\n" +
	"fail\tactive\tgo\texit 3\t./*.go\t—\n" +
	"nocode\tactive\tgo\tok\tdocs/*.md\t—\n" +
	"notool\tactive\tlayup-no-such-tool-xyz\tok\t./*.go\t—\n" +
	"waiting\tpending\tgo\tok\tlayout\t—\n" +
	"changed\tpending\tgo\tok\t./*.go\t—\n"

// newFake gives a fake whose base holds the manifest and whose head holds a
// Go file; the change touches x.go.
func newFake() *fakeGit {
	return &fakeGit{
		version: "2.54.0",
		commits: map[string]string{"main^{commit}": baseID, "HEAD^{commit}": headID},
		files:   map[string]string{baseID + ":docs/gates.tsv": manifest},
		changed: []string{"README.md", "x.go"},
		head:    map[string]string{"x.go": "package x\n", "README.md": "x\n"},
	}
}

// use puts f, a lookPath that finds each tool but layup-no-such-tool-xyz, and
// a shell that gives the outcome of each command ("ok", "exit 3", "signal")
// in place of the real ones until the test ends. seen gets the tree of each
// run.
func use(t *testing.T, f *fakeGit, seen func(tree string)) {
	t.Helper()
	t.Setenv("TMPDIR", t.TempDir())
	savedGit, savedLook, savedShell := repoAPI, lookPath, runShell
	repoAPI = f
	lookPath = func(tool string) (string, error) {
		if tool == "layup-no-such-tool-xyz" {
			return "", errors.New("not found")
		}
		return "/bin/" + tool, nil
	}
	runShell = func(dir, command string) (outcome, []byte) {
		if seen != nil {
			seen(dir)
		}
		switch command {
		case "exit 3":
			return outcome{code: 3}, []byte("vet: x.go:1: a finding\n")
		case "signal":
			return outcome{signal: "terminated"}, nil
		}
		return outcome{}, []byte("ok, no newline")
	}
	t.Cleanup(func() { repoAPI, lookPath, runShell = savedGit, savedLook, savedShell })
}

// events records the calls of the step function and the writes of the run.
type events struct{ log []string }

func (e *events) step(i, n int, kind string) func() {
	e.log = append(e.log, fmt.Sprintf("step %d/%d %s", i, n, kind))
	return func() { e.log = append(e.log, "done "+kind) }
}

func (e *events) Write(p []byte) (int, error) {
	e.log = append(e.log, "out "+string(p))
	return len(p), nil
}

func TestRunGivesOneRowPerKindInTheOrderOfTheManifest(t *testing.T) {
	f := newFake()
	use(t, f, nil)
	var e events
	table, err := Run("repo", "main", "HEAD", e.step, &e)
	if err != nil {
		t.Fatal(err)
	}
	want := []Row{
		{"pass", "active", pass, ""},
		{"fail", "active", fail, "exit 3"},
		{"nocode", "active", clear, "no product path"},
		{"notool", "active", notActive, "tool not found: layup-no-such-tool-xyz"},
		{"waiting", "pending", clear, "pending: no product path"},
		{"changed", "pending", fail, "pending: product path changed: x.go"},
	}
	if table.Base != baseID || table.Head != headID || !slices.Equal(table.Rows, want) {
		t.Fatalf("the table %+v\nwant base %s, head %s, rows %+v", table, baseID, headID, want)
	}
	wantLog := []string{"step 1/6 pass", "done pass", "out ok, no newline", "out \n", "step 2/6 fail", "done fail",
		"out vet: x.go:1: a finding\n", "step 3/6 nocode", "done nocode", "step 4/6 notool", "done notool",
		"step 5/6 waiting", "done waiting", "step 6/6 changed", "done changed"}
	if !slices.Equal(e.log, wantLog) {
		t.Fatalf("the steps and the output:\n%q\nwant:\n%q", e.log, wantLog)
	}
	if len(f.added) != 1 || !slices.Equal(f.removed, f.added) {
		t.Fatalf("added %q, removed %q; want one scratch tree, removed", f.added, f.removed)
	}
	var b bytes.Buffer
	if err := table.Write(&b); err != nil || !strings.HasPrefix(b.String(), "base\thead\tkind\tstate\tresult\treason\n"+baseID+"\t"+headID+"\tpass\tactive\tpass\t—\n") {
		t.Fatalf("the table as written: %q, %v", b.String(), err)
	}
	if got := table.Results(); !slices.Equal(got, []string{pass, fail, clear, notActive, clear, fail}) {
		t.Fatalf("Results = %q", got)
	}
}

func TestRunGivesAnInputErrorAndNoTable(t *testing.T) {
	for name, change := range map[string]func(f *fakeGit){
		"git not found":               func(f *fakeGit) { f.version = "" },
		"git older than 2.32":         func(f *fakeGit) { f.version = "2.31.8" },
		"a base that is not a commit": func(f *fakeGit) { delete(f.commits, "main^{commit}") },
		"a head that is not a commit": func(f *fakeGit) { delete(f.commits, "HEAD^{commit}") },
		"no manifest at the base":     func(f *fakeGit) { delete(f.files, baseID+":docs/gates.tsv") },
		"a malformed manifest":        func(f *fakeGit) { f.files[baseID+":docs/gates.tsv"] = "kind\n" },
		"a manifest with no row":      func(f *fakeGit) { f.files[baseID+":docs/gates.tsv"] = manifestHeader },
		"a config path that is a link at the base": func(f *fakeGit) {
			f.files[baseID+":docs/gates.tsv"] = manifestHeader + "pass\tactive\tgo\tok\t./*.go\trun.sh\n"
			f.trees = map[string][]git.TreeEntry{baseID + ":run.sh": {{Mode: "120000", Type: "blob", Path: "run.sh"}}}
		},
		"a config path that is a submodule at the base": func(f *fakeGit) {
			f.files[baseID+":docs/gates.tsv"] = manifestHeader + "pass\tactive\tgo\tok\t./*.go\tvendor\n"
			f.trees = map[string][]git.TreeEntry{baseID + ":vendor": {{Mode: "160000", Type: "commit", Path: "vendor"}}}
		},
	} {
		f := newFake()
		change(f)
		use(t, f, nil)
		var e events
		table, err := Run("repo", "main", "HEAD", e.step, &e)
		var input *InputError
		if !errors.As(err, &input) || len(table.Rows) != 0 || len(f.added) != 0 {
			t.Errorf("%s: %d rows, error %v, %d scratch trees; want an *InputError, no row and no scratch tree", name, len(table.Rows), err, len(f.added))
		}
	}
	f := newFake()
	use(t, f, nil)
	lookPath = func(tool string) (string, error) { return "", errors.New("not found") }
	if _, err := Run("repo", "main", "HEAD", (&events{}).step, &events{}); !errors.As(err, new(*InputError)) {
		t.Errorf("sh not found: %v; want an *InputError", err)
	}
}

// Condition 1 of the plan review: a failure gives a fixed reason, with no part
// of the error; the error goes to standard error. A pending row needs no
// scratch tree, so it keeps its result (note 4).
func TestAFailureOfTheScratchTreeOrTheDiffGivesNotActiveRows(t *testing.T) {
	f := newFake()
	f.addErr = errors.New("fatal: '/tmp/layup-gate-123/tree' already exists")
	use(t, f, nil)
	var e events
	table, err := Run("repo", "main", "HEAD", e.step, &e)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range table.Rows {
		want := map[string]string{"active": "scratch tree: add failed", "pending": r.Reason}[r.State]
		if r.State == "active" && (r.Result != notActive || r.Reason != want) || r.State == "pending" && r.Result == notActive {
			t.Errorf("after a failed add: %+v", r)
		}
		if strings.Contains(r.Reason, "/") {
			t.Errorf("the reason %q holds a path", r.Reason)
		}
	}
	if !strings.Contains(strings.Join(e.log, ""), "already exists") {
		t.Errorf("the error of git is not on standard error: %q", e.log)
	}
	f = newFake()
	f.diffErr = errors.New("fatal: bad object")
	use(t, f, nil)
	table, _ = Run("repo", "main", "HEAD", (&events{}).step, &events{})
	for _, r := range table.Rows {
		if r.State == "pending" && (r.Result != notActive || r.Reason != "diff failed") || r.State == "active" && r.Result == notActive && r.Kind != "notool" {
			t.Errorf("after a failed diff: %+v", r)
		}
	}
}

func TestACleanupFailureComesWithTheWholeTable(t *testing.T) {
	f := newFake()
	f.removeErr = errors.New("fatal: cannot remove")
	use(t, f, nil)
	table, err := Run("repo", "main", "HEAD", (&events{}).step, &events{})
	var cleanup *CleanupError
	if !errors.As(err, &cleanup) || len(table.Rows) != 6 {
		t.Fatalf("%d rows, %v; want 6 rows and a *CleanupError", len(table.Rows), err)
	}
}

// The scratch tree holds the base's manifest and the base's files of each
// config path, with their modes, and no file of the head under a config path.
func TestTheOverlayPutsTheGateFilesOfTheBase(t *testing.T) {
	f := newFake()
	f.files[baseID+":docs/gates.tsv"] = manifestHeader + "pass\tactive\tgo\tok\t./*.go\tlayout run.sh gone.txt\n"
	f.files[baseID+":layout/a_test.go"] = "package layout // base\n"
	f.files[baseID+":run.sh"] = "#!/bin/sh\necho base\n"
	f.trees = map[string][]git.TreeEntry{
		baseID + ":layout": {{Mode: "100644", Type: "blob", Path: "layout/a_test.go"}},
		baseID + ":run.sh": {{Mode: "100755", Type: "blob", Path: "run.sh"}},
	}
	f.head = map[string]string{"x.go": "package x\n", "docs/gates.tsv": "the head's own manifest\n",
		"layout/a_test.go": "package layout // head\n", "layout/new_test.go": "package layout // head only\n",
		"run.sh": "echo head\n", "gone.txt": "the head's own\n"}
	got := map[string]string{}
	use(t, f, func(tree string) {
		for _, p := range []string{"docs/gates.tsv", "layout/a_test.go", "layout/new_test.go", "run.sh", "gone.txt", "x.go"} {
			b, err := os.ReadFile(filepath.Join(tree, p))
			got[p] = string(b)
			if err != nil {
				got[p] = "absent"
			}
		}
		if info, err := os.Stat(filepath.Join(tree, "run.sh")); err != nil || info.Mode().Perm() != 0o755 {
			got["run.sh mode"] = fmt.Sprint(info.Mode().Perm())
		}
	})
	if _, err := Run("repo", "main", "HEAD", (&events{}).step, &events{}); err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"docs/gates.tsv": f.files[baseID+":docs/gates.tsv"], "layout/a_test.go": "package layout // base\n",
		"layout/new_test.go": "absent", "run.sh": "#!/bin/sh\necho base\n", "gone.txt": "absent", "x.go": "package x\n"}
	for p, w := range want {
		if got[p] != w {
			t.Errorf("%s in the scratch tree: %q; want %q", p, got[p], w)
		}
	}
	if m, ok := got["run.sh mode"]; ok {
		t.Errorf("the mode of run.sh: %s; want 755", m)
	}
}

// A symbolic link of the head cannot send a write of the overlay out of the
// tree: the write fails, and the active rows are not-active.
func TestTheOverlayDoesNotWriteOutOfTheTree(t *testing.T) {
	f := newFake()
	f.link, f.linkTarget = "docs", t.TempDir()
	use(t, f, nil)
	table, err := Run("repo", "main", "HEAD", (&events{}).step, &events{})
	if err != nil || len(table.Rows) != 6 {
		t.Fatalf("%d rows, %v; want 6 rows", len(table.Rows), err)
	}
	for _, r := range table.Rows {
		if r.State == "active" && (r.Result != notActive || r.Reason != "scratch tree: overlay failed") {
			t.Errorf("%+v; want not-active, scratch tree: overlay failed", r)
		}
	}
	if _, err := os.Stat(filepath.Join(f.linkTarget, "gates.tsv")); err == nil {
		t.Error("the overlay wrote gates.tsv out of the tree, through the link")
	}
}

// Review round 1, finding 1: two config paths that overlap give one result on
// each run, whatever the order of a map: the base's file a, and no a/b.
func TestOverlappingConfigPathsGiveOneResult(t *testing.T) {
	for run := 0; run < 20; run++ {
		f := newFake()
		f.files[baseID+":docs/gates.tsv"] = manifestHeader + "pass\tactive\tgo\tok\t./*.go\ta a/b\n"
		f.files[baseID+":a"] = "the base's a\n"
		f.trees = map[string][]git.TreeEntry{baseID + ":a": {{Mode: "100644", Type: "blob", Path: "a"}}}
		f.head = map[string]string{"x.go": "package x\n", "a/b": "the head's a/b\n"}
		var a string
		use(t, f, func(tree string) {
			b, err := os.ReadFile(filepath.Join(tree, "a"))
			a = fmt.Sprint(string(b), err)
		})
		table, err := Run("repo", "main", "HEAD", (&events{}).step, &events{})
		if err != nil || len(table.Rows) != 1 || table.Rows[0].Result != pass || a != "the base's a\n<nil>" {
			t.Fatalf("run %d: %+v, %v, a = %q; want pass and the base's file a", run, table.Rows, err, a)
		}
	}
}
