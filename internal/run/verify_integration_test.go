//go:build integration

package run

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
	workpkg "github.com/pharzam/layup/internal/work"
)

// The step verify on the real step runner (docs/spec/scaffold.md, The step verify
// and the records of the setup; The README of the Scaffold; The work area; the
// input state of a row that is not pass or clear; the acceptance row "layup run
// with the real step runner"; row 55 of the plan, task T-79zk, #278).

// setupRepo is a repository with main at one commit and layup-setup above it with
// one commit for each step, as the runner makes them.
func setupRepo(t *testing.T, steps ...string) string {
	t.Helper()
	dir := t.TempDir()
	gitOut(t, dir, "init", "-q", "-b", "main")
	writeFile(t, filepath.Join(dir, "baseline.txt"), "baseline\n")
	gitOut(t, dir, "add", "-A")
	gitOut(t, dir, "commit", "-q", "-m", "the baseline")
	gitOut(t, dir, "switch", "-q", "-c", "layup-setup")
	for _, step := range steps {
		writeFile(t, filepath.Join(dir, step+".txt"), step+"\n")
		gitOut(t, dir, "add", "-A")
		gitOut(t, dir, "commit", "-q", "-m", "chore: setup "+step)
	}
	return dir
}

// The tree of the setup is the commit of main with the step root, then each commit of
// layup-setup above it, oldest first, by its first parents, with its tree and the step
// of its message: the order of the commits (S14 comes before S11), not of the steps.
func TestTheTreeOfTheSetupIsTheCommitsOfLayupSetupAboveMain(t *testing.T) {
	t.Parallel()
	dir := setupRepo(t, "S04", "S05", "S14", "S11")
	data, err := setupTree(dir)
	if err != nil {
		t.Fatal(err)
	}
	rows, err := ReadSetupTree(data)
	if err != nil {
		t.Fatalf("setup/tree.tsv does not read back by its schema: %v\n%s", err, data)
	}
	commits := append([]string{strings.TrimSpace(gitOut(t, dir, "rev-parse", "main"))}, strings.Fields(gitOut(t, dir, "log", "--reverse", "--format=%H", "main..layup-setup"))...)
	steps := []string{"root", "S04", "S05", "S14", "S11"}
	if len(rows) != len(steps) || len(commits) != len(steps) {
		t.Fatalf("%d rows and %d commits; want %d of each:\n%s", len(rows), len(commits), len(steps), data)
	}
	for i, r := range rows {
		tree := strings.TrimSpace(gitOut(t, dir, "rev-parse", commits[i]+"^{tree}"))
		if want := []string{string(rune('0' + i)), commits[i], tree, steps[i]}; strings.Join(r, "|") != strings.Join(want, "|") {
			t.Errorf("row %d: %q; want %q", i, r, want)
		}
	}
	if head := strings.TrimSpace(gitOut(t, dir, "rev-parse", "layup-setup")); rows[len(rows)-1][1] != head {
		t.Errorf("the last row is %s; want the head of layup-setup %s", rows[len(rows)-1][1], head)
	}
}

func TestTheTreeOfTheSetupRefusesWhatTheRunnerNeverMakes(t *testing.T) {
	t.Parallel()
	bad := func(message string) string {
		dir := setupRepo(t, "S04")
		writeFile(t, filepath.Join(dir, "x.txt"), "x\n")
		gitOut(t, dir, "add", "-A")
		gitOut(t, dir, "commit", "-q", "-m", message)
		return dir
	}
	unrelated := setupRepo(t)
	gitOut(t, unrelated, "switch", "-q", "--orphan", "other")
	gitOut(t, unrelated, "commit", "-q", "--allow-empty", "-m", "chore: setup S04")
	gitOut(t, unrelated, "branch", "-q", "-M", "layup-setup")
	for _, c := range []struct{ name, dir, want string }{
		{"a message of another form", bad("wip"), `"wip"`},
		{"a step that does not begin with S", bad("chore: setup X04"), `"chore: setup X04"`},
		{"a step of one digit", bad("chore: setup S4"), `"chore: setup S4"`},
		{"a step with a letter for a digit", bad("chore: setup S0A"), `"chore: setup S0A"`},
		{"a main that is not an ancestor of layup-setup", unrelated, "main is not an ancestor of layup-setup"},
		{"no commit above main", setupRepo(t), "no commit above main"},
	} {
		if data, err := setupTree(c.dir); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %q, %v; want an error with %q", c.name, data, err, c.want)
		}
	}
}

// goEnv makes the go program of the gates of layup setup verify run on this host's
// build cache, with no download and no setting of the host, as internal/cli's test
// of the setup does; the HOME of the package is set once, by TestMain.
func goEnv(t *testing.T) {
	t.Helper()
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{"GOCACHE": strings.TrimSpace(string(cache)), "GOENV": "off", "GOFLAGS": "", "GOTOOLCHAIN": "local", "GOPROXY": "off"} {
		t.Setenv(k, v)
	}
}

func recordsFile(t *testing.T, w *world, path string) string {
	return show(t, w.bare, "refs/heads/layup-records", path)
}

// The demo of row 55: after S14 a restart runs layup setup verify and S15 and commits
// the records of the setup in one records commit, five paths and the earlier records
// kept; a row that is not pass or clear is the fail of the step with no such commit; a
// second restart finds setup and verify "done before" (the work area that is missing is
// built again by the step scaffold-push that follows, which waits); and, after the
// Operator's push of the setup commits and a records branch that lacks the record, the
// work area is built again to the same commits, its main the root commit by setup/tree.tsv.
func TestTheStepVerifyCommitsTheRecordsOfTheSetup(t *testing.T) {
	goEnv(t)
	saved := verifyRun
	t.Cleanup(func() { verifyRun = saved })
	w := scaffoldWorld(t, proseFiles())
	cfg := w.config("0.1.0-dev")
	// A row that is not pass or clear: the step is fail, naming the check and its reason, and no records commit is made.
	verifyRun = func(string, func(int, int, string) func(), io.Writer) (verify.Table, error) {
		return verify.Table{Rows: []verify.Row{{Check: "pin", Result: "pass", Reason: "—"}, {Check: "link-lint", Result: "fail", Reason: "docs/guide.md:3: a broken link"}}}, nil
	}
	rows := toSetup(t, w, "go")
	if r := rowOf(rows, "setup"); r != (Step{"setup", "done", setupDoneDetail}) {
		t.Fatalf("the step setup: %+v; want done\n%s", r, steps(rows))
	}
	if r := rowOf(rows, "verify"); r != (Step{"verify", "fail", "layup setup verify: link-lint is fail: docs/guide.md:3: a broken link"}) {
		t.Fatalf("the step verify with a row that is fail: %+v\n%s", r, steps(rows))
	}
	for _, p := range []string{"setup/record.tsv", "setup/verify.tsv", "rule-paths.tsv", "setup/tree.tsv"} {
		if out, err := gitRunOK(w.bare, "show", "refs/heads/layup-records:"+p); err == nil {
			t.Errorf("%s is in the records after the fail: %q; want no records commit of the setup", p, out)
		}
	}
	// The real call: one commit of five paths.
	verifyRun = saved
	rows = Restart(context.Background(), cfg)
	if r := rowOf(rows, "verify"); r != (Step{"verify", "done", verifyDoneDetail}) {
		t.Fatalf("the step verify: %+v; want done\n%s", r, steps(rows))
	}
	w.mu.Lock()
	lines := strings.Join(w.lines, "\n")
	w.mu.Unlock()
	if !strings.Contains(lines, "verify: [1/") {
		t.Errorf("no progress line of layup setup verify (verify: [i/n] check) among\n%s", lines)
	}
	area := filepath.Join(w.dir, "work", "acme", "target")
	target := filepath.Join(area, "target")
	for path, from := range map[string]string{"setup/record.tsv": "out/record.tsv", "setup/verify.tsv": "out/verify.tsv", "rule-paths.tsv": "out/rule-paths.tsv"} {
		want, err := os.ReadFile(filepath.Join(area, filepath.FromSlash(from)))
		if err != nil || recordsFile(t, w, path) != string(want) || len(want) == 0 {
			t.Errorf("%s in the records: %q, %v; want %s of the work area, %q", path, recordsFile(t, w, path), err, from, want)
		}
	}
	if got := recordsFile(t, w, "README.md"); got != scaffoldReadme {
		t.Errorf("README.md in the records:\n%s\nwant the README of the Scaffold", got)
	}
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "--grep=^layup run: the records of the setup$", "refs/heads/layup-records"))
	changed := strings.Fields(gitOut(t, w.bare, "diff", "--name-only", commit+"~1", commit))
	if want := "README.md rule-paths.tsv setup/record.tsv setup/tree.tsv setup/verify.tsv"; strings.Join(changed, " ") != want {
		t.Errorf("the paths of the records commit: %q; want exactly %q", changed, want)
	}
	for _, p := range []string{"start/start.tsv", "questions.tsv", "setup/inputs/files/docs/guide.md"} {
		if out, err := gitRunOK(w.bare, "show", "refs/heads/layup-records:"+p); err != nil {
			t.Errorf("%s is not in the records after the commit: %v %s", p, err, out)
		}
	}
	treeRows, err := ReadSetupTree([]byte(recordsFile(t, w, "setup/tree.tsv")))
	head := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup"))
	root := strings.TrimSpace(gitOut(t, target, "rev-parse", "main"))
	if err != nil || len(treeRows) != 11 || treeRows[0][1] != root || treeRows[10][1] != head {
		t.Fatalf("setup/tree.tsv: %d rows, %v; want 11, the root %s first and the head %s last\n%s", len(treeRows), err, root, head, recordsFile(t, w, "setup/tree.tsv"))
	}
	var got []string
	for _, r := range treeRows {
		got = append(got, r[3])
	}
	if want := "root S04 S05 S06 S07 S08 S09 S14 S11 S12 S13"; strings.Join(got, " ") != want {
		t.Errorf("the steps of setup/tree.tsv: %q; want %q", strings.Join(got, " "), want)
	}
	// The Operator's push, then a second restart with the work area gone: both steps are done before.
	gitOut(t, target, "push", "-q", w.bare, "layup-setup:main")
	if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
		t.Fatal(err)
	}
	rows = Restart(context.Background(), cfg)
	for _, name := range []string{"setup", "verify"} {
		if r := rowOf(rows, name); r != (Step{name, "done", "done before"}) {
			t.Errorf("the step %s of a second restart: %+v; want done before\n%s", name, r, steps(rows))
		}
	}
	if r := rowOf(rows, "scaffold-push"); r != (Step{"scaffold-push", "done", "waiting for the Operator's pushes and ruleset apply"}) {
		t.Errorf("the step scaffold-push of a second restart: %+v; want done, waiting\n%s", r, steps(rows))
	}
	if _, err := os.Stat(filepath.Join(w.dir, "work", "acme", "target")); err != nil {
		t.Errorf("the step scaffold-push, which waits, did not build the work area again: %v", err)
	}
	// A records branch without the record of the setup: the work area is built again from the records, with
	// main of target/ the root commit though the default branch holds the setup commits, to the same commits.
	clone := filepath.Join(t.TempDir(), "records")
	gitOut(t, "/", "clone", "-q", "--branch", "layup-records", "--single-branch", "file://"+w.bare, clone)
	gitOut(t, clone, "rm", "-q", "setup/record.tsv")
	gitOut(t, clone, "commit", "-q", "-m", "a records branch that lacks the record of the setup")
	gitOut(t, clone, "push", "-q", "origin", "layup-records")
	if got := strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "refs/heads/main")); got != head {
		t.Fatalf("main of the forge is %s; want the setup head %s after the push", got, head)
	}
	rows = Restart(context.Background(), cfg)
	if r := rowOf(rows, "verify"); r != (Step{"verify", "done", verifyDoneDetail}) {
		t.Fatalf("the step verify after the push: %+v; want done\n%s", r, steps(rows))
	}
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup")); got != head {
		t.Errorf("layup-setup of the work area built again: %s; want the same head %s", got, head)
	}
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "main")); got != root {
		t.Errorf("main of the work area built again: %s; want the root commit %s although the forge's default branch is at the head", got, root)
	}
	if got := recordsFile(t, w, "setup/tree.tsv"); got != gitOut(t, w.bare, "show", commit+":setup/tree.tsv") {
		t.Errorf("setup/tree.tsv of the second verify differs from the first")
	}
	again := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "--grep=^layup run: the records of the setup$", "refs/heads/layup-records"))
	if changed := strings.Fields(gitOut(t, w.bare, "diff", "--name-only", again+"~1", again)); again == commit || strings.Join(changed, " ") != "setup/record.tsv" {
		t.Errorf("the second records commit %s (the first was %s) changes %q; want a new commit that changes setup/record.tsv alone, as the other four paths are written with the bytes that the records hold", again, commit, changed)
	}
}

// The three files of out/ go to their paths in the records, byte for byte, once S15 completed
// (the agreed contract of the plan review): out/record.tsv holds its rows setup.head,
// verify.sha256 and rule-paths.sha256, each hash that of its file, and the register holds a
// path. A file that is missing, a row that is missing or a hash that is not its file's is an error.
func TestTheFilesOfOutAreTheRecordsOfTheSetup(t *testing.T) {
	t.Parallel()
	sum := func(text string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(text))) }
	var reg, none bytes.Buffer
	if err := tsv.Write(&reg, setup.RulePathsSchema, [][]string{{".github/", "—", "baseline"}}); err != nil {
		t.Fatal(err)
	}
	if err := tsv.Write(&none, setup.RulePathsSchema, nil); err != nil {
		t.Fatal(err)
	}
	const table = "table\n"
	// out writes the work area of a case: the rows of S15 that it lacks or holds wrong, and the files.
	out := func(drop string, wrong map[string]string, register string, missing string) string {
		rows := [][]string{{"S15", "setup.head", strings.Repeat("a", 40), "computed", "git rev-parse layup-setup"},
			{"S15", "verify.sha256", sum(table), "computed", "sha256 out/verify.tsv"}, {"S15", "rule-paths.sha256", sum(register), "computed", "sha256 out/rule-paths.tsv"}}
		var kept [][]string
		for _, r := range rows {
			if v, ok := wrong[r[1]]; ok {
				r[2] = v
			}
			if r[1] != drop {
				kept = append(kept, r)
			}
		}
		var rec bytes.Buffer
		if err := tsv.Write(&rec, workpkg.RecordSchema, kept); err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		for p, text := range map[string]string{"out/record.tsv": rec.String(), "out/verify.tsv": table, "out/rule-paths.tsv": register} {
			if p != missing {
				writeFile(t, filepath.Join(dir, filepath.FromSlash(p)), text)
			}
		}
		return dir
	}
	got, err := setupOutputs(out("", nil, reg.String(), ""))
	if err != nil || len(got) != 3 || string(got["setup/verify.tsv"]) != table || string(got["rule-paths.tsv"]) != reg.String() || !strings.Contains(string(got["setup/record.tsv"]), "setup.head") {
		t.Fatalf("the files of a work area whose S15 completed: %q, %v; want the three of out/, byte for byte", got, err)
	}
	for _, c := range []struct {
		name, dir, want string
	}{
		{"no row setup.head", out("setup.head", nil, reg.String(), ""), "setup.head"},
		{"no row verify.sha256", out("verify.sha256", nil, reg.String(), ""), "verify.sha256"},
		{"no row rule-paths.sha256", out("rule-paths.sha256", nil, reg.String(), ""), "rule-paths.sha256"},
		{"a hash of another table", out("", map[string]string{"verify.sha256": sum("other\n")}, reg.String(), ""), "verify.sha256"},
		{"a hash of another register", out("", map[string]string{"rule-paths.sha256": sum("other\n")}, reg.String(), ""), "rule-paths.sha256"},
		{"a register with no path", out("", nil, none.String(), ""), "holds no path"},
		{"a register that its schema refuses", out("", nil, "x\\n", ""), "rule-paths.tsv"},
		{"a missing table", out("", nil, reg.String(), "out/verify.tsv"), "verify.tsv"},
		{"a missing register", out("", nil, reg.String(), "out/rule-paths.tsv"), "rule-paths.tsv"},
		{"a missing record", out("", nil, reg.String(), "out/record.tsv"), "record.tsv"},
	} {
		if _, err := setupOutputs(c.dir); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error with %q", c.name, err, c.want)
		}
	}
}
