//go:build integration

package gate

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
)

// The Go values of the two schemas equal their blocks in docs/spec/.
func TestTheSchemaBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"gate-manifest": ManifestSchema, "gate-result": ResultSchema} {
		block, ok := blocks[name]
		if !ok {
			t.Fatalf("docs/spec/ has no block %s", name)
		}
		if err := tsv.Compare(block, s); err != nil {
			t.Error(err)
		}
	}
}

// who is the identity of the commits of the tests.
var who = git.Identity{Name: "LAYUP test", Email: "test@layup.invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

// repo makes a repository with a commit of the files, and gives its path.
func repo(t *testing.T, files map[string]string) string {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	dir := t.TempDir()
	if err := git.Init(dir); err != nil {
		t.Fatal(err)
	}
	commitFiles(t, dir, files, nil)
	return dir
}

// commitFiles writes the files (a text that starts with "#!" is executable),
// removes the paths of gone, and commits.
func commitFiles(t *testing.T, dir string, files map[string]string, gone []string) {
	t.Helper()
	for name, text := range files {
		mode := os.FileMode(0o644)
		if strings.HasPrefix(text, "#!") {
			mode = 0o755
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), mode); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range gone {
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			t.Fatal(err)
		}
	}
	if err := git.Add(dir); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(dir, "test", who); err != nil {
		t.Fatal(err)
	}
}

// plainGit runs git as a tool of the test, with no configuration of the host.
func plainGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t.invalid"}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// recorder is a step function and a writer for Run.
type recorder struct{ out strings.Builder }

func (r *recorder) step(int, int, string) func() { return func() {} }
func (r *recorder) Write(p []byte) (int, error)  { return r.out.Write(p) }

const realManifest = "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
	"files\tactive\tsh\tpwd; cat layout/a_test.go docs/gates.tsv; test ! -e layout/new_test.go\t./*.go\tlayout\n" +
	"exits\tactive\tsh\texit 3\t./*.go\t—\n" +
	"killed\tactive\tsh\tkill -TERM $$\t./*.go\t—\n" +
	"waiting\tpending\tsh\ttrue\t./*.go\t—\n"

// On a real repository: the scratch tree holds the base's manifest and gate
// files; sh -c gives pass, exit 3 and a signal; the change of a Go file fails
// the pending kind; afterwards the tree is gone and REPO is as it was.
func TestRunOnARealRepository(t *testing.T) {
	dir := repo(t, map[string]string{"docs/gates.tsv": realManifest, "layout/a_test.go": "package layout // BASE\n", "x.go": "package x\n"})
	base := plainGit(t, dir, "rev-parse", "HEAD")
	commitFiles(t, dir, map[string]string{"docs/gates.tsv": realManifest + "extra\tactive\tsh\ttrue\t./*.go\t—\n",
		"layout/a_test.go": "package layout // HEAD\n", "layout/new_test.go": "package layout\n", "y.go": "package x\n"}, nil)
	head := plainGit(t, dir, "rev-parse", "HEAD")
	status := plainGit(t, dir, "status", "--porcelain")
	var r recorder
	table, err := Run(dir, base, "HEAD", r.step, &r)
	if err != nil {
		t.Fatal(err)
	}
	want := []Row{{"files", "active", pass, ""}, {"exits", "active", fail, "exit 3"},
		{"killed", "active", fail, "signal terminated"}, {"waiting", "pending", fail, "pending: product path changed: layout/a_test.go"}}
	if table.Base != base || table.Head != head || !slices.Equal(table.Rows, want) {
		t.Fatalf("the table %+v\nwant base %s, head %s, rows %+v", table, base, head, want)
	}
	block := r.out.String()
	if !strings.Contains(block, "// BASE") || strings.Contains(block, "// HEAD") || strings.Contains(block, "extra\tactive") {
		t.Fatalf("the block of the kind files does not show the base's gate files:\n%s", block)
	}
	tree := strings.SplitN(block, "\n", 2)[0]
	if _, err := os.Stat(tree); err == nil || !strings.Contains(tree, "layup-gate-") {
		t.Fatalf("the scratch tree %q: still there (%v), or not the run's tree", tree, err)
	}
	if list := plainGit(t, dir, "worktree", "list", "--porcelain"); strings.Count(list, "worktree ") != 1 {
		t.Fatalf("another work tree is listed:\n%s", list)
	}
	if s := plainGit(t, dir, "status", "--porcelain"); s != status || plainGit(t, dir, "rev-parse", "HEAD") != head {
		t.Fatalf("REPO changed: status %q, want %q", s, status)
	}
}

// A hook of REPO does not run, through a core.hooksPath of REPO either.
func TestAHookOfTheRepositoryDoesNotRun(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "hook-ran")
	dir := repo(t, map[string]string{"docs/gates.tsv": "kind\tstate\ttool\tcommand\tscope\tconfig\nok\tactive\tsh\ttrue\t./*.go\t—\n",
		"x.go": "package x\n", "hooks/post-checkout": "#!/bin/sh\ntouch " + marker + "\n"})
	plainGit(t, dir, "config", "core.hooksPath", "hooks")
	control := filepath.Join(t.TempDir(), "control")
	plainGit(t, dir, "worktree", "add", "--detach", control, "HEAD")
	if _, err := os.Stat(marker); err != nil {
		t.Fatal("the control: a plain git worktree add did not run the hook, so the test proves nothing")
	}
	plainGit(t, dir, "worktree", "remove", "--force", control)
	if err := os.Remove(marker); err != nil {
		t.Fatal(err)
	}
	table, err := Run(dir, "HEAD", "HEAD", (&recorder{}).step, &recorder{})
	if err != nil || len(table.Rows) != 1 || table.Rows[0].Result != pass {
		t.Fatalf("%+v, %v; want the row ok, pass", table.Rows, err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("the post-checkout hook of REPO ran")
	}
}

// A tag is peeled to its commit; a tree, a revision that does not resolve and
// a base with no manifest are input errors.
func TestTheRevisionsAndTheManifestOnARealRepository(t *testing.T) {
	dir := repo(t, map[string]string{"docs/gates.tsv": "kind\tstate\ttool\tcommand\tscope\tconfig\nok\tactive\tsh\ttrue\t./*.go\t—\n", "x.go": "package x\n"})
	plainGit(t, dir, "tag", "-a", "-m", "an annotated tag", "v1")
	commitID := plainGit(t, dir, "rev-parse", "HEAD")
	table, err := Run(dir, "HEAD", "v1", (&recorder{}).step, &recorder{})
	if err != nil || table.Head != commitID || plainGit(t, dir, "rev-parse", "v1") == commitID {
		t.Fatalf("the head of a tag: %s, %v; want the commit %s", table.Head, err, commitID)
	}
	for _, rev := range []string{"HEAD^{tree}", "no-such-revision", "HEAD:x.go"} {
		if _, err := Run(dir, rev, "HEAD", (&recorder{}).step, &recorder{}); !errors.As(err, new(*InputError)) {
			t.Errorf("the base %s: %v; want an *InputError", rev, err)
		}
	}
	empty := repo(t, map[string]string{"x.go": "package x\n"})
	if _, err := Run(empty, "HEAD", "HEAD", (&recorder{}).step, &recorder{}); !errors.As(err, new(*InputError)) {
		t.Errorf("a base with no manifest: %v; want an *InputError", err)
	}
}

// git diff names a renamed path at both ends and a deleted path, so both fail
// a pending kind.
func TestARenameAndADeleteChangeAProductPath(t *testing.T) {
	manifest := "kind\tstate\ttool\tcommand\tscope\tconfig\nwaiting\tpending\tsh\ttrue\tinternal/*.go\t—\n"
	dir := repo(t, map[string]string{"docs/gates.tsv": manifest, "internal/a.go": "package a\n", "internal/c.go": "package a\n"})
	base := plainGit(t, dir, "rev-parse", "HEAD")
	plainGit(t, dir, "mv", "internal/a.go", "internal/b.go")
	commitFiles(t, dir, nil, []string{"internal/c.go"})
	table, err := Run(dir, base, "HEAD", (&recorder{}).step, &recorder{})
	if err != nil || len(table.Rows) != 1 || table.Rows[0].Result != fail || table.Rows[0].Reason != "pending: product path changed: internal/a.go" {
		t.Fatalf("%+v, %v; want fail on internal/a.go, the first path that git names", table.Rows, err)
	}
}

// Condition 1 of the plan review: a scratch tree that cannot be made gives
// the active rows a reason with no path; the error names the path on
// standard error.
func TestAFailedScratchTreeGivesAReasonWithNoPath(t *testing.T) {
	dir := repo(t, map[string]string{"docs/gates.tsv": "kind\tstate\ttool\tcommand\tscope\tconfig\nok\tactive\tsh\ttrue\t./*.go\t—\n", "x.go": "package x\n"})
	t.Setenv("TMPDIR", filepath.Join(t.TempDir(), "no-such-directory"))
	var r recorder
	table, err := Run(dir, "HEAD", "HEAD", r.step, &r)
	if err != nil || len(table.Rows) != 1 || table.Rows[0] != (Row{"ok", "active", notActive, "scratch tree: add failed"}) {
		t.Fatalf("%+v, %v; want not-active, scratch tree: add failed", table.Rows, err)
	}
	if !regexp.MustCompile(`layup gate: the scratch tree: .*no-such-directory`).MatchString(r.out.String()) {
		t.Fatalf("standard error does not name the path:\n%s", r.out.String())
	}
}
