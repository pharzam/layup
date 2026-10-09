//go:build integration

package git

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// isolate gives the test a home of its own, so that no git process of the
// test reads or writes the real one.
func isolate(t *testing.T) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	return home
}

// plain is the environment of a git process that the test starts itself,
// outside the package: no configuration file of the host, and who.
func plain(home string, extra ...string) []string {
	return append([]string{"PATH=" + os.Getenv("PATH"), "HOME=" + home, "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0", "GIT_AUTHOR_NAME=LAYUP test",
		"GIT_AUTHOR_EMAIL=test@layup.invalid", "GIT_COMMITTER_NAME=LAYUP test", "GIT_COMMITTER_EMAIL=test@layup.invalid"}, extra...)
}

// gitIn runs git as a tool of the test, and gives its output.
func gitIn(dir string, env []string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir, cmd.Env = dir, env
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// gitOK is gitIn, and an error ends the test.
func gitOK(t *testing.T, dir string, env []string, args ...string) string {
	t.Helper()
	out, err := gitIn(dir, env, args...)
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// write writes each file under dir; a file whose text starts with "#!" is
// executable.
func write(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for name, text := range files {
		mode := os.FileMode(0o644)
		if strings.HasPrefix(text, "#!") {
			mode = 0o755
		}
		must(t, os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755))
		must(t, os.WriteFile(filepath.Join(dir, name), []byte(text), mode))
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// commitTree makes a repository in dir that holds files and one commit by
// who, and gives the commit.
func commitTree(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	must(t, Init(dir))
	write(t, dir, files)
	must(t, Add(dir))
	must(t, Commit(dir, "chore: one", who))
	id, err := RevParse(dir, "HEAD")
	if err != nil || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(id) {
		t.Fatalf("commit %q (%v): want 40 hexadecimal characters", id, err)
	}
	return id
}

// The root commit is on main, one input gives one commit ID (NFR-005), and
// the tree keeps each file name byte for byte: decomposed (NFD) or composed.
func TestInitAddAndCommit(t *testing.T) {
	home := isolate(t)
	tree := map[string]string{"a.txt": "a\n", "docs/b.md": "b\n", "run.sh": "#!/bin/sh\n", "cafe\u0301.txt": "nfd\n", "\u00e9t\u00e9.md": "nfc\n"}
	a := t.TempDir()
	if idA, idB := commitTree(t, a, tree), commitTree(t, t.TempDir(), tree); idA != idB {
		t.Fatalf("two commits of one input: %s and %s", idA, idB)
	}
	if head := gitOK(t, a, plain(home), "symbolic-ref", "HEAD"); head != "refs/heads/main\n" {
		t.Errorf("HEAD %q, want refs/heads/main", head)
	}
	obj := gitOK(t, a, plain(home), "cat-file", "-p", "HEAD")
	for _, role := range []string{"author", "committer"} {
		if line := role + " LAYUP test <test@layup.invalid> 1767225600 +0000\n"; !strings.Contains(obj, line) {
			t.Errorf("the commit object has no line %q:\n%s", line, obj)
		}
	}
	if names, err := LsFiles(a); err != nil || !reflect.DeepEqual(names, []string{"a.txt", "cafe\u0301.txt", "docs/b.md", "run.sh", "\u00e9t\u00e9.md"}) {
		t.Errorf("LsFiles: %+q, %v", names, err)
	}
}

// S02: ls-remote and clone of a bare repository by a file URL, the checkout
// of the commit, and the tree of the commit.
func TestCloneAndCheckout(t *testing.T) {
	home := isolate(t)
	src, bare, dst := t.TempDir(), filepath.Join(t.TempDir(), "baseline.git"), filepath.Join(t.TempDir(), "target")
	id := commitTree(t, src, map[string]string{"a.txt": "a\n"})
	gitOK(t, "", plain(home), "clone", "--bare", "--", src, bare)
	if got, err := LsRemote("file://"+bare, "HEAD"); err != nil || got != id {
		t.Fatalf("LsRemote: %q, %v; want %s", got, err, id)
	}
	must(t, Clone("file://"+bare, dst, Auth{}))
	if exists(filepath.Join(dst, "a.txt")) {
		t.Fatal("the clone checked out a.txt; want no file before the checkout")
	}
	must(t, CheckoutDetach(dst, id))
	if b, err := os.ReadFile(filepath.Join(dst, "a.txt")); err != nil || string(b) != "a\n" {
		t.Fatalf("a.txt after the checkout: %q, %v", b, err)
	}
	want, _ := RevParse(src, "HEAD^{tree}")
	if got, err := RevParse(dst, id+"^{tree}"); err != nil || got != want {
		t.Errorf("the tree of the copy: %q, %v; want %s", got, err, want)
	}
}

// S04 and S15: the setup branch and an orphan branch.
func TestBranchesAndAnOrphan(t *testing.T) {
	home := isolate(t)
	dir := t.TempDir()
	root := commitTree(t, dir, map[string]string{"a.txt": "a\n"})
	must(t, SwitchCreate(dir, "layup-setup", root))
	write(t, dir, map[string]string{"b.txt": "b\n"})
	must(t, Add(dir))
	must(t, Commit(dir, "chore: setup S04", who))
	if head := gitOK(t, dir, plain(home), "symbolic-ref", "HEAD"); head != "refs/heads/layup-setup\n" {
		t.Errorf("HEAD %q, want refs/heads/layup-setup", head)
	}
	must(t, SwitchOrphan(dir, "layup-records"))
	if names, err := LsFiles(dir); err != nil || len(names) != 0 {
		t.Fatalf("LsFiles on the orphan branch: %q, %v; want no file", names, err)
	}
	write(t, dir, map[string]string{"README.md": "records\n"})
	must(t, Add(dir, "README.md"))
	must(t, Commit(dir, "chore: setup S15", who))
	records, _ := RevParse(dir, "layup-records")
	for branch, want := range map[string]string{"layup-records": records, "layup-setup": root} {
		if roots, err := RootCommits(dir, branch); err != nil || !reflect.DeepEqual(roots, []string{want}) {
			t.Errorf("the roots of %s: %q, %v; want only %s", branch, roots, err, want)
		}
	}
}

// The step runner (task T-79y7): the branch of the work tree, and an error
// for a detached HEAD, so a setup commit lands only on layup-setup.
func TestBranch(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	root := commitTree(t, dir, map[string]string{"a.txt": "a\n"})
	for _, c := range []struct {
		name, want string
		move       func() error
	}{
		{"after init", "refs/heads/main", func() error { return nil }},
		{"after switch -c", "refs/heads/layup-setup", func() error { return SwitchCreate(dir, "layup-setup", root) }},
		{"detached", "", func() error { return CheckoutDetach(dir, root) }},
	} {
		must(t, c.move())
		if got, err := Branch(dir); got != c.want || (err == nil) != (c.want != "") {
			t.Errorf("%s: Branch %q, %v; want %q", c.name, got, err, c.want)
		}
	}
}

// layup gate and the fixture run of layup setup verify: a scratch work tree,
// a patch, a commit on no ref, the changed paths, a file at a revision, and
// the scratch tree removed.
func TestAScratchTree(t *testing.T) {
	home := isolate(t)
	dir, scratch, patch := t.TempDir(), filepath.Join(t.TempDir(), "scratch"), filepath.Join(t.TempDir(), "bad.patch")
	base := commitTree(t, dir, map[string]string{"a.txt": "a\n", "b.txt": "b\n"})
	must(t, WorktreeAdd(dir, scratch, base))
	must(t, os.WriteFile(patch, []byte("--- a/a.txt\n+++ b/a.txt\n@@ -1 +1 @@\n-a\n+A\n"), 0o644))
	must(t, Apply(scratch, patch))
	must(t, Add(scratch))
	must(t, Commit(scratch, "fixture", who))
	fixture, _ := RevParse(scratch, "HEAD")
	if refs := gitOK(t, dir, plain(home), "for-each-ref", "--contains", fixture); refs != "" {
		t.Errorf("refs that hold the fixture commit: %q; want none", refs)
	}
	if names, err := DiffNames(dir, base, fixture); err != nil || !reflect.DeepEqual(names, []string{"a.txt"}) {
		t.Errorf("DiffNames: %q, %v; want a.txt", names, err)
	}
	must(t, os.Rename(filepath.Join(dir, "b.txt"), filepath.Join(dir, "c.txt")))
	must(t, Add(dir))
	must(t, Commit(dir, "a rename", who))
	if names, err := DiffNames(dir, base, "HEAD"); err != nil || !reflect.DeepEqual(names, []string{"b.txt", "c.txt"}) {
		t.Errorf("DiffNames of a rename: %q, %v; want both ends, b.txt and c.txt", names, err)
	}
	if b, err := Show(dir, base, "b.txt"); err != nil || string(b) != "b\n" {
		t.Errorf("Show of a present path: %q, %v", b, err)
	}
	var failed *FailedError
	if _, err := Show(dir, base, "missing.txt"); !errors.As(err, &failed) || failed.Code != 128 {
		t.Errorf("Show of a missing path: %v; want a *FailedError with code 128", err)
	}
	must(t, WorktreeRemove(dir, scratch))
	if list := gitOK(t, dir, plain(home), "worktree", "list", "--porcelain"); exists(scratch) || strings.Count(list, "worktree ") != 1 {
		t.Errorf("after the removal: the scratch tree is there (%v), or another work tree is listed:\n%s", exists(scratch), list)
	}
}

// No call starts the automatic maintenance of git: a commit runs git
// maintenance run --auto (git 2.29), whose tasks go to the background (2.47);
// in CI of the pull request #115, on git 2.55.0, a background repack wrote
// into .git/objects while a test removed the directory. A repository
// whose own configuration asks for the loose-objects task at once, in the
// foreground, keeps its loose objects after a Commit; a control shows that a
// plain git commit packs them, so the test can fail.
func TestNoCallStartsTheMaintenance(t *testing.T) {
	home := isolate(t)
	setUp := func() string {
		dir := filepath.Join(t.TempDir(), "r")
		gitOK(t, "", plain(home), "init", "-q", "-b", "main", dir)
		for _, kv := range [][2]string{{"maintenance.auto", "true"}, {"maintenance.autoDetach", "false"}, {"gc.autoDetach", "false"},
			{"maintenance.loose-objects.enabled", "true"}, {"maintenance.loose-objects.auto", "1"}} {
			gitOK(t, dir, plain(home), "config", kv[0], kv[1])
		}
		write(t, dir, map[string]string{"a.txt": "a\n"})
		gitOK(t, dir, plain(home), "add", "-A")
		return dir
	}
	packs := func(dir string) []string {
		m, _ := filepath.Glob(filepath.Join(dir, ".git", "objects", "pack", "*.pack"))
		return m
	}
	ctl := setUp()
	gitOK(t, ctl, plain(home, "GIT_AUTHOR_DATE=@0 +0000", "GIT_COMMITTER_DATE=@0 +0000"), "commit", "-q", "-m", "control")
	if len(packs(ctl)) == 0 {
		t.Fatal("control: a plain git commit started no maintenance, so this test cannot fail")
	}
	dir := setUp()
	must(t, Commit(dir, "chore: one", who))
	if p := packs(dir); len(p) != 0 {
		t.Errorf("a Commit started the maintenance of git: the packs %q", p)
	}
}

// Condition 1 of the plan review of #79: the inputs (a) to (e) of the host, a
// global configuration and a hook change no byte or mode of a tree, no file,
// no author, no branch and no repository. A control first shows that each of
// (a) to (e) changes a plain git run, so the test can fail. Commit sets its
// identity after the fixed list, so (e) cannot change a commit here: the proof
// of (e) is TestTheEnvironmentIsAFixedList.
func TestAHostileHostChangesNothing(t *testing.T) {
	home := isolate(t)
	tree := map[string]string{"crlf.txt": "a\r\nb\r\n", "doc.md": "# doc\n", "run.sh": "#!/bin/sh\n"}
	want := commitTree(t, t.TempDir(), tree)
	for _, dir := range []string{os.Getenv("XDG_CONFIG_HOME"), filepath.Join(home, ".config")} {
		write(t, dir, map[string]string{"git/attributes": "* text=auto\n", "git/ignore": "*.md\n"}) // (a), (b)
	}
	write(t, home, map[string]string{".gitconfig": "[init]\n\tdefaultBranch = master\n[core]\n\tautocrlf = true\n[commit]\n\tgpgsign = true\n"})
	hostile := []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=core.fileMode", "GIT_CONFIG_VALUE_0=false", // (c)
		"GIT_AUTHOR_NAME=Hostile"} // (e)
	hook := func(dir string) {
		must(t, os.WriteFile(filepath.Join(dir, ".git", "hooks", "pre-commit"), []byte("#!/bin/sh\ntouch hook-ran\nexit 1\n"), 0o755))
	}

	t.Run("control: each of (a) to (e) changes a plain git run", func(t *testing.T) {
		first := strings.Fields("-c core.hooksPath=/dev/null -c core.autocrlf=false -c commit.gpgsign=false -c maintenance.auto=false") // D3 before the review, and no background repack (#92)
		ctl, decoy, other := t.TempDir(), t.TempDir(), t.TempDir()
		gitOK(t, "", plain(home), "init", "-q", "-b", "main", ctl)
		write(t, ctl, tree)
		gitOK(t, ctl, plain(home, hostile...), append(first, "add", "-A")...)
		gitOK(t, ctl, plain(home, hostile...), append(first, "commit", "-q", "-m", "chore: one")...)
		files := gitOK(t, ctl, plain(home), "ls-files", "-s")
		if blob := gitOK(t, ctl, plain(home), "cat-file", "-p", "HEAD:crlf.txt"); blob != "a\nb\n" {
			t.Errorf("(a) the attributes file left %q", blob)
		}
		if strings.Contains(files, "doc.md") || !regexp.MustCompile(`(?m)^100644 .*\trun\.sh$`).MatchString(files) {
			t.Errorf("(b) the ignore file kept doc.md, or (c) GIT_CONFIG_COUNT kept the mode of run.sh:\n%s", files)
		}
		if author := gitOK(t, ctl, plain(home), "log", "-1", "--format=%an"); author != "Hostile\n" {
			t.Errorf("(e) the author is %q", author)
		}
		gitOK(t, "", plain(home), "init", "-q", "-b", "main", decoy)
		write(t, other, tree)
		gitOK(t, other, plain(home, "GIT_DIR="+filepath.Join(decoy, ".git")), append(first, "add", "-A")...)
		gitOK(t, other, plain(home, "GIT_DIR="+filepath.Join(decoy, ".git")), append(first, "commit", "-q", "-m", "x")...)
		if _, err := gitIn(decoy, plain(home), "rev-parse", "--verify", "-q", "HEAD"); err != nil {
			t.Error("(d) GIT_DIR did not send the commit to the other repository")
		}
	})

	for _, kv := range append(hostile, "GIT_DIR="+filepath.Join(t.TempDir(), "decoy.git")) { // and (d)
		k, v, _ := strings.Cut(kv, "=")
		t.Setenv(k, v)
	}
	dir := t.TempDir()
	must(t, Init(dir))
	if !exists(filepath.Join(dir, ".git")) {
		t.Fatalf("Init made no repository in %s: GIT_DIR of the host sent it to %s", dir, os.Getenv("GIT_DIR"))
	}
	write(t, dir, tree)
	hook(dir)
	must(t, Add(dir))
	must(t, Commit(dir, "chore: one", who))
	if got, err := RevParse(dir, "HEAD"); err != nil || got != want {
		show, _ := gitIn(dir, plain(home), "show", "--stat", "--format=fuller", "HEAD")
		t.Errorf("commit %q (%v), want %s: the host changed the tree or the author\n%s", got, err, want, show)
	}
	head, _ := gitIn(dir, plain(home), "symbolic-ref", "HEAD")
	if head != "refs/heads/main\n" || exists(filepath.Join(dir, "hook-ran")) || exists(os.Getenv("GIT_DIR")) {
		t.Errorf("HEAD %q, want refs/heads/main; or the hook ran (%v); or a call wrote to GIT_DIR (%v)",
			head, exists(filepath.Join(dir, "hook-ran")), exists(os.Getenv("GIT_DIR")))
	}
}

// No call uses a credential of the host or asks a question (K31). A loopback
// server that wants a password gets no Authorization header, although the
// host's .netrc has one for it; an ssh URL and a remote helper's URL start no
// program of the PATH; and each clone fails at once.
func TestNoCallUsesACredentialOfTheHost(t *testing.T) {
	home, bin, ran := isolate(t), t.TempDir(), t.TempDir()
	write(t, home, map[string]string{".netrc": "machine 127.0.0.1 login host password secret\n"})
	for _, name := range []string{"ssh", "git-remote-layuptest"} {
		write(t, bin, map[string]string{name: "#!/bin/sh\ntouch '" + filepath.Join(ran, name) + "'\nexit 1\n"})
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	var asked, sent atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if asked.Add(1); r.Header.Get("Authorization") != "" {
			sent.Add(1)
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="layup-test"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	for url, want := range map[string]string{srv.URL + "/baseline.git": "terminal prompts disabled",
		"ssh://git@127.0.0.1/baseline.git": "transport 'ssh' not allowed", "layuptest::baseline": "transport 'layuptest' not allowed"} {
		done := make(chan error, 1)
		go func() { done <- Clone(url, filepath.Join(t.TempDir(), "target"), Auth{}) }()
		select {
		case err := <-done:
			var failed *FailedError
			if !errors.As(err, &failed) || !strings.Contains(failed.Stderr, want) {
				t.Errorf("clone of %s: %v; want a *FailedError that says %s", url, err, want)
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("the clone of %s did not end in 60 s: it waits for an answer", url)
		}
	}
	if asked.Load() == 0 || sent.Load() != 0 {
		t.Errorf("the server got %d requests, %d with an Authorization header; want one or more, none with it", asked.Load(), sent.Load())
	}
	if started, _ := os.ReadDir(ran); len(started) != 0 {
		t.Errorf("%d programs of the PATH started, first %s; want none", len(started), started[0].Name())
	}
}

// The two error kinds with real processes, and the version of the host's git.
func TestTheErrorKindsAndTheVersion(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	must(t, Init(dir))
	var failed *FailedError
	if _, err := RevParse(dir, "no-such-revision"); !errors.As(err, &failed) || failed.Code != 128 || !strings.Contains(failed.Stderr, "fatal") {
		t.Errorf("git failed: %v; want a *FailedError with code 128 and the standard error", err)
	}
	v, err := Version()
	must(t, err)
	t.Logf("git %s", v)
	if !Supported(v) {
		t.Errorf("git %s is older than %s", v, MinVersion)
	}
	t.Setenv("PATH", "")
	var notFound *NotFoundError
	if _, err := Version(); !errors.As(err, &notFound) {
		t.Errorf("an empty PATH: %v; want a *NotFoundError", err)
	}
}

// LsTree gives each file at or under a path, with its mode: a file, an
// executable, a directory with its files, a symbolic link; and no entry for a
// path that the commit does not have.
func TestLsTree(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	write(t, dir, map[string]string{"layout/a_test.go": "package a\n", "layout/b/c.txt": "c\n", "run.sh": "#!/bin/sh\n"})
	must(t, os.Symlink("run.sh", filepath.Join(dir, "link.sh")))
	head := commitTree(t, dir, nil)
	for path, want := range map[string][]string{
		"layout":           {"100644 blob layout/a_test.go", "100644 blob layout/b/c.txt"},
		"run.sh":           {"100755 blob run.sh"},
		"link.sh":          {"120000 blob link.sh"},
		"layout/a_test.go": {"100644 blob layout/a_test.go"},
		"missing":          nil,
	} {
		entries, err := LsTree(dir, head, path)
		var got []string
		for _, e := range entries {
			got = append(got, e.Mode+" "+e.Type+" "+e.Path)
			if len(e.Object) != 40 {
				t.Errorf("%s: the object %q is not 40 characters", path, e.Object)
			}
		}
		if err != nil || !reflect.DeepEqual(got, want) {
			t.Errorf("LsTree(%q) = %q, %v; want %q", path, got, err, want)
		}
	}
}

// Check pin: a repository with its whole history is not shallow, and a clone
// of depth 1 is, also in a work tree of it (layup setup verify reads a
// scratch work tree).
func TestIsShallow(t *testing.T) {
	home := isolate(t)
	full := t.TempDir()
	commitTree(t, full, map[string]string{"a.txt": "a\n"})
	shallow := filepath.Join(t.TempDir(), "shallow")
	gitOK(t, "", plain(home), "clone", "-q", "--depth", "1", "file://"+full, shallow)
	tree := filepath.Join(t.TempDir(), "tree")
	must(t, WorktreeAdd(shallow, tree, "HEAD"))
	for _, c := range []struct {
		dir  string
		want bool
	}{{full, false}, {shallow, true}, {tree, true}} {
		if got, err := IsShallow(c.dir); err != nil || got != c.want {
			t.Errorf("IsShallow(%s): %v, %v; want %v", c.dir, got, err, c.want)
		}
	}
}

// The step runner (task T-b3r1): Staged sees a change only when the index
// holds it, and ResetHard puts the index and the work tree back to HEAD,
// with a file that only the index holds removed (D11, D12 of #90).
func TestStagedAndResetHard(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	commitTree(t, dir, map[string]string{"a.txt": "a\n"})
	if staged, err := Staged(dir); err != nil || staged {
		t.Fatalf("a clean tree: %v, %v; want no staged change", staged, err)
	}
	write(t, dir, map[string]string{"a.txt": "changed\n", "new.txt": "new\n"})
	if staged, err := Staged(dir); err != nil || staged {
		t.Errorf("a change that is not staged: %v, %v; want false", staged, err)
	}
	must(t, Add(dir))
	if staged, err := Staged(dir); err != nil || !staged {
		t.Errorf("a staged change: %v, %v; want true", staged, err)
	}
	must(t, ResetHard(dir))
	data, err := os.ReadFile(filepath.Join(dir, "a.txt"))
	if staged, serr := Staged(dir); err != nil || string(data) != "a\n" || exists(filepath.Join(dir, "new.txt")) || serr != nil || staged {
		t.Errorf("after ResetHard: a.txt %q, new.txt there %v, staged %v; want the tree of HEAD", data, exists(filepath.Join(dir, "new.txt")), staged)
	}
}

// The demo of row 23 (#128): a push to a local bare repository passes, a
// fast-forward passes, and a push that is not a fast-forward is refused with
// code 1, the branch of the bare repository unchanged (fencing; packages.md,
// The calls of M2a). Then Init and Fetch read the branch back, as step 5 of
// layup run does.
func TestPushAndFetchOfABareRepository(t *testing.T) {
	home := isolate(t)
	bare := filepath.Join(t.TempDir(), "target.git")
	gitOK(t, "", plain(home), "init", "--bare", "-q", bare)
	url := "file://" + bare

	work := t.TempDir()
	first := commitTree(t, work, map[string]string{"a.txt": "a\n"})
	must(t, Push(work, url, first, "main", Auth{}))
	write(t, work, map[string]string{"b.txt": "b\n"})
	must(t, Add(work))
	must(t, Commit(work, "chore: two", who))
	second, err := RevParse(work, "HEAD")
	must(t, err)
	must(t, Push(work, url, second, "main", Auth{}))

	// A commit that does not descend from the branch: an orphan of work.
	write(t, work, map[string]string{"c.txt": "c\n"})
	gitOK(t, work, plain(home), "checkout", "-q", "--orphan", "side")
	must(t, Add(work))
	must(t, Commit(work, "chore: side", who))
	side, err := RevParse(work, "HEAD")
	must(t, err)
	err = Push(work, url, side, "main", Auth{})
	var failed *FailedError
	if !errors.As(err, &failed) || failed.Code != 1 {
		t.Fatalf("a push that is not a fast-forward: %v; want a *FailedError of code 1", err)
	}
	if head := strings.TrimSpace(gitOK(t, bare, plain(home), "rev-parse", "refs/heads/main")); head != second {
		t.Errorf("the branch of the bare repository is %s after the refused push, want %s", head, second)
	}

	clone := filepath.Join(t.TempDir(), "clone")
	must(t, Init(clone))
	must(t, Fetch(clone, url, "refs/heads/main", Auth{}))
	if got, err := RevParse(clone, "refs/heads/main"); err != nil || got != second {
		t.Errorf("Init, then Fetch of refs/heads/main: %s, %v; want %s", got, err, second)
	}

	err = Push(work, "file://"+filepath.Join(t.TempDir(), "none.git"), side, "main", Auth{})
	if !errors.As(err, &failed) || failed.Code == 1 {
		t.Errorf("a push to a remote that does not exist: %v; want a *FailedError of a code other than 1", err)
	}
}

// With a token, the request of git carries the header of forge.md (The App
// identity) to the web of the register; with none, it carries no header.
func TestTheTokenReachesTheServerAsAHeaderOnly(t *testing.T) {
	isolate(t)
	var got atomic.Value
	got.Store("")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header.Get("Authorization"); h != "" {
			got.Store(h)
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	dir := filepath.Join(t.TempDir(), "clone")
	must(t, Init(dir))
	err := Fetch(dir, srv.URL+"/acme/target.git", "refs/heads/main", Auth{Web: srv.URL, Token: "ghs_testtoken"})
	if err == nil || strings.Contains(err.Error(), "ghs_testtoken") {
		t.Errorf("the fetch from a server that has no repository: %v; want an error with no token", err)
	}
	if h := got.Load().(string); h != "Basic eC1hY2Nlc3MtdG9rZW46Z2hzX3Rlc3R0b2tlbg==" {
		t.Errorf("the server got Authorization %q, want Basic of x-access-token and the token", h)
	}
	got.Store("")
	Fetch(dir, srv.URL+"/acme/target.git", "refs/heads/main", Auth{})
	if h := got.Load().(string); h != "" {
		t.Errorf("with no token the server got Authorization %q, want none", h)
	}
}

// CloneLocal makes a --no-local clone of one branch: no other branch, no tag,
// no remote, and objects of its own, so the log reads after the source is
// gone (packages.md, The calls of M2b; session.md, The session directory).
func TestCloneLocal(t *testing.T) {
	home := isolate(t)
	src := t.TempDir()
	base := commitTree(t, src, map[string]string{"a.txt": "a\n"})
	gitOK(t, src, plain(home), "branch", "layup-records")
	gitOK(t, src, plain(home), "tag", "v1")
	dir := filepath.Join(t.TempDir(), "repo")
	must(t, CloneLocal(src, dir, "main"))
	if refs := gitOK(t, dir, plain(home), "for-each-ref", "--format=%(refname)"); refs != "refs/heads/main\n" {
		t.Errorf("the refs of the clone: %q; want refs/heads/main alone", refs)
	}
	if cfg, err := os.ReadFile(filepath.Join(dir, ".git", "config")); err != nil || strings.Contains(string(cfg), "[remote") {
		t.Errorf("the configuration of the clone holds a remote (%v):\n%s", err, cfg)
	}
	if present := exists(filepath.Join(dir, ".git", "objects", "info", "alternates")); present {
		t.Error("the clone shares the objects of its source")
	}
	must(t, os.RemoveAll(src))
	if got := gitOK(t, dir, plain(home), "log", "-1", "--format=%H", "main"); got != base+"\n" {
		t.Errorf("the log after the source is gone: %q; want %s", got, base)
	}
	must(t, SwitchCreate(dir, "task/T-ab12/1", base))
	if head := gitOK(t, dir, plain(home), "symbolic-ref", "HEAD"); head != "refs/heads/task/T-ab12/1\n" {
		t.Errorf("HEAD %q after SwitchCreate", head)
	}
}

// session makes a session's clone of the run's clone, on its task branch at
// the base, with one commit of the session; it gives the clone and the head.
func session(t *testing.T, home, runClone, base string) (string, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "repo")
	must(t, CloneLocal(runClone, dir, "main"))
	must(t, SwitchCreate(dir, "task/T-ab12/1", base))
	write(t, dir, map[string]string{"b.txt": "b\n"})
	gitOK(t, dir, plain(home), "add", "b.txt")
	gitOK(t, dir, plain(home), "commit", "-q", "-m", "feat: the session's change")
	return dir, strings.TrimSpace(gitOK(t, dir, plain(home), "rev-parse", "HEAD"))
}

// FetchSession brings the session's head into the run's clone, and leaves no
// scratch directory; a SHA of no commit is refused and sets no ref.
func TestFetchSession(t *testing.T) {
	home := isolate(t)
	t.Setenv("TMPDIR", t.TempDir())
	runClone := t.TempDir()
	base := commitTree(t, runClone, map[string]string{"a.txt": "a\n"})
	dir, head := session(t, home, runClone, base)
	dst := "refs/heads/task/T-ab12/1"
	must(t, FetchSession(runClone, filepath.Join(dir, ".git"), head, dst))
	if got := gitOK(t, runClone, plain(home), "rev-parse", dst); got != head+"\n" {
		t.Errorf("%s in the run's clone: %q; want %s", dst, got, head)
	}
	blob := strings.TrimSpace(gitOK(t, dir, plain(home), "rev-parse", "HEAD:b.txt"))
	for name, sha := range map[string]string{"a blob": blob, "an absent object": strings.Repeat("e", 40)} {
		err := FetchSession(runClone, filepath.Join(dir, ".git"), sha, "refs/heads/task/T-ab12/2")
		if !errors.Is(err, ErrNotACommit) {
			t.Errorf("%s: %v; want ErrNotACommit", name, err)
		}
		if _, err := gitIn(runClone, plain(home), "rev-parse", "--verify", "-q", "refs/heads/task/T-ab12/2"); err == nil {
			t.Errorf("%s: a ref was set", name)
		}
	}
	if left, _ := os.ReadDir(os.Getenv("TMPDIR")); len(left) != 0 {
		t.Errorf("scratch directories left in TMPDIR: %v", left)
	}
}

// IsAncestor, DiffFile and DiffBinary read two commits of a real repository;
// the binary diff applies back to the base.
func TestTheReadsOfM2b(t *testing.T) {
	home := isolate(t)
	runClone := t.TempDir()
	base := commitTree(t, runClone, map[string]string{"a.txt": "a\n", "docs/guardrails.md": "# G\n\n## 2. Pitfalls\n\n- one\n\n## 3. Validation\n"})
	write(t, runClone, map[string]string{"docs/guardrails.md": "# G\n\n## 2. Pitfalls\n\n- one\n- two\n\n## 3. Validation\n"})
	must(t, os.WriteFile(filepath.Join(runClone, "bin.dat"), []byte{0, 1, 2, 0, 255}, 0o644))
	must(t, Add(runClone))
	must(t, Commit(runClone, "chore: two", who))
	head, err := RevParse(runClone, "HEAD")
	must(t, err)
	if yes, err := IsAncestor(runClone, base, head); err != nil || !yes {
		t.Errorf("IsAncestor(base, head) = %v, %v; want true", yes, err)
	}
	if yes, err := IsAncestor(runClone, head, base); err != nil || yes {
		t.Errorf("IsAncestor(head, base) = %v, %v; want false", yes, err)
	}
	if _, err := IsAncestor(runClone, base, strings.Repeat("e", 40)); err == nil {
		t.Error("IsAncestor of an absent commit: no error")
	}
	diff, err := DiffFile(runClone, base, head, "docs/guardrails.md")
	if err != nil || !strings.Contains(string(diff), "\n+- two\n") || regexp.MustCompile(`(?m)^-[^-]`).Match(diff) {
		t.Errorf("DiffFile: %v\n%s", err, diff)
	}
	patch, err := DiffBinary(runClone, base, head)
	if err != nil || !strings.Contains(string(patch), "GIT binary patch") {
		t.Fatalf("DiffBinary: %v\n%s", err, patch)
	}
	back := t.TempDir()
	gitOK(t, "", plain(home), "clone", "-q", runClone, back)
	gitOK(t, back, plain(home), "checkout", "-q", base)
	p := filepath.Join(t.TempDir(), "p.patch")
	must(t, os.WriteFile(p, patch, 0o644))
	gitOK(t, back, plain(home), "apply", p)
	if got, _ := os.ReadFile(filepath.Join(back, "bin.dat")); string(got) != string([]byte{0, 1, 2, 0, 255}) {
		t.Errorf("the binary file after git apply: %v", got)
	}
}

// hostileKeys are the keys that `git help --config` of the host's git lists
// (2.47.3 on the host of task T-z5dj; the Operator's host has 2.54.0) and that
// name a program or a shell command, each with a marker program of its own
// name; each driver key has a driver of its own, so one key does not hide
// another. A key that only picks a tool (help.browser, web.browser,
// instaweb.browser, diff.tool, merge.tool) reaches a program only through a key
// of the list. live are the ones that a plain read in the clone fires on this
// host (the control shows each); notLive are the ones that no plain read fires,
// so their absence after FetchSession is asserted, not shown live.
// remote.v.vcs names a helper git-remote-<vcs> of the PATH, not a path, so it
// has no marker.
var (
	hostileKeys = []string{"core.fsmonitor", "core.hooksPath", "core.sshCommand", "core.askPass", "core.editor", "core.pager",
		"credential.helper", "diff.external", "diff.x.command", "diff.t.textconv", "filter.y.clean", "filter.y.smudge",
		"filter.p.process", "merge.x.driver", "uploadpack.packObjectsHook", "sequence.editor", "gpg.program",
		"remote.origin.uploadpack", "remote.origin.receivepack", "remote.ext.url", "core.alternateRefsCommand", "include.path",
		"includeIf.path", "alias.y", "core.gitProxy", "browser.x.cmd", "browser.x.path", "difftool.x.cmd", "mergetool.x.cmd",
		"man.x.cmd", "man.x.path", "guitool.x.cmd", "gc.recentObjectsHook", "gpg.ssh.program", "gpg.ssh.defaultKeyCommand",
		"imap.tunnel", "instaweb.httpd", "interactive.diffFilter", "pager.status", "sendemail.ccCmd", "sendemail.headerCmd",
		"sendemail.toCmd", "sendemail.smtpServer", "submodule.x.update", "remote.v.vcs"}
	notLive = []string{"core.sshCommand", "core.askPass", "core.editor", "core.pager", "credential.helper",
		"uploadpack.packObjectsHook", "sequence.editor", "gpg.program", "remote.origin.receivepack", "merge.x.driver",
		"browser.x.cmd", "browser.x.path", "difftool.x.cmd", "mergetool.x.cmd", "man.x.cmd", "man.x.path", "guitool.x.cmd",
		"gc.recentObjectsHook", "gpg.ssh.program", "gpg.ssh.defaultKeyCommand", "imap.tunnel", "instaweb.httpd",
		"interactive.diffFilter", "pager.status", "sendemail.ccCmd", "sendemail.headerCmd", "sendemail.toCmd",
		"sendemail.smtpServer", "submodule.x.update", "remote.v.vcs"}
)

// arm writes into the clone at dir a configuration that holds each hostile
// key, each program touching a marker of its name in marks; the hooks
// directory holds every hook of githooks(5). up is a repository with a commit
// that the clone lacks, for the fetch that reads the clone's alternates.
func arm(t *testing.T, dir, marks, up string) {
	t.Helper()
	bin := t.TempDir()
	prog := func(name string) string {
		p := filepath.Join(bin, strings.NewReplacer("/", "_", ".", "_").Replace(name))
		write(t, "/", map[string]string{strings.TrimPrefix(p, "/"): "#!/bin/sh\ntouch '" + filepath.Join(marks, name) + "'\nexit 1\n"})
		return p
	}
	hooks := filepath.Join(t.TempDir(), "hooks")
	for _, h := range strings.Fields("applypatch-msg pre-applypatch post-applypatch pre-commit pre-merge-commit prepare-commit-msg " +
		"commit-msg post-commit pre-rebase post-checkout post-merge pre-push pre-receive update proc-receive post-receive " +
		"post-update reference-transaction push-to-checkout pre-auto-gc post-rewrite sendemail-validate fsmonitor-watchman " +
		"p4-changelist p4-prepare-changelist p4-post-changelist p4-pre-submit post-index-change") {
		write(t, "/", map[string]string{strings.TrimPrefix(filepath.Join(hooks, h), "/"): "#!/bin/sh\ntouch '" + filepath.Join(marks, "core.hooksPath") + "'\nexit 0\n"})
	}
	include := filepath.Join(t.TempDir(), "included")
	write(t, "/", map[string]string{strings.TrimPrefix(include, "/"): "[alias]\n\tx = !" + prog("include.path") + "\n"})
	includeIf := filepath.Join(t.TempDir(), "includedIf")
	write(t, "/", map[string]string{strings.TrimPrefix(includeIf, "/"): "[alias]\n\tz = !" + prog("includeIf.path") + "\n"})
	cfg := "[core]\n\tfsmonitor = " + prog("core.fsmonitor") + "\n\thooksPath = " + hooks + "\n\tsshCommand = " + prog("core.sshCommand") +
		"\n\taskPass = " + prog("core.askPass") + "\n\teditor = " + prog("core.editor") + "\n\tpager = " + prog("core.pager") +
		"\n\talternateRefsCommand = " + prog("core.alternateRefsCommand") +
		"\n[credential]\n\thelper = " + prog("credential.helper") +
		"\n[diff]\n\texternal = " + prog("diff.external") +
		"\n[diff \"x\"]\n\tcommand = " + prog("diff.x.command") + "\n[diff \"t\"]\n\ttextconv = " + prog("diff.t.textconv") +
		"\n[filter \"y\"]\n\tclean = " + prog("filter.y.clean") + "\n\tsmudge = " + prog("filter.y.smudge") +
		"\n[filter \"p\"]\n\tprocess = " + prog("filter.p.process") +
		"\n[merge \"x\"]\n\tdriver = " + prog("merge.x.driver") +
		"\n[uploadpack]\n\tpackObjectsHook = " + prog("uploadpack.packObjectsHook") +
		"\n[sequence]\n\teditor = " + prog("sequence.editor") +
		"\n[gpg]\n\tprogram = " + prog("gpg.program") +
		"\n[remote \"origin\"]\n\turl = " + up + "\n\tuploadpack = " + prog("remote.origin.uploadpack") + "\n\treceivepack = " + prog("remote.origin.receivepack") +
		"\n[remote \"up\"]\n\turl = " + up +
		"\n[remote \"ext\"]\n\turl = ext::" + prog("remote.ext.url") +
		"\n[protocol \"ext\"]\n\tallow = always\n[include]\n\tpath = " + include +
		"\n[includeIf \"gitdir:" + dir + "/\"]\n\tpath = " + includeIf +
		"\n[alias]\n\ty = !" + prog("alias.y") +
		"\n[core]\n\tgitProxy = " + prog("core.gitProxy") + "\n[remote \"gp\"]\n\turl = git://layup.invalid/x" +
		"\n[browser \"x\"]\n\tcmd = " + prog("browser.x.cmd") + "\n\tpath = " + prog("browser.x.path") +
		"\n[difftool \"x\"]\n\tcmd = " + prog("difftool.x.cmd") + "\n[mergetool \"x\"]\n\tcmd = " + prog("mergetool.x.cmd") +
		"\n[man \"x\"]\n\tcmd = " + prog("man.x.cmd") + "\n\tpath = " + prog("man.x.path") +
		"\n[guitool \"x\"]\n\tcmd = " + prog("guitool.x.cmd") + "\n[gc]\n\trecentObjectsHook = " + prog("gc.recentObjectsHook") +
		"\n[gpg \"ssh\"]\n\tprogram = " + prog("gpg.ssh.program") + "\n\tdefaultKeyCommand = " + prog("gpg.ssh.defaultKeyCommand") +
		"\n[imap]\n\ttunnel = " + prog("imap.tunnel") + "\n[instaweb]\n\thttpd = " + prog("instaweb.httpd") +
		"\n[interactive]\n\tdiffFilter = " + prog("interactive.diffFilter") + "\n[pager]\n\tstatus = " + prog("pager.status") +
		"\n[sendemail]\n\tccCmd = " + prog("sendemail.ccCmd") + "\n\theaderCmd = " + prog("sendemail.headerCmd") +
		"\n\ttoCmd = " + prog("sendemail.toCmd") + "\n\tsmtpServer = " + prog("sendemail.smtpServer") +
		"\n[submodule \"x\"]\n\tupdate = !" + prog("submodule.x.update") + "\n[remote \"v\"]\n\tvcs = layup-hostile\n"
	f, err := os.OpenFile(filepath.Join(dir, ".git", "config"), os.O_APPEND|os.O_WRONLY, 0)
	must(t, err)
	_, err = f.WriteString(cfg)
	must(t, err)
	must(t, f.Close())
	write(t, dir, map[string]string{".git/info/attributes": "a.txt diff=x merge=x\nt.txt diff=t\n*.y filter=y\n*.p filter=p\n"})
}

// marked gives the names of the markers in marks.
func marked(t *testing.T, marks string) []string {
	t.Helper()
	entries, err := os.ReadDir(marks)
	must(t, err)
	var names []string
	for _, e := range entries {
		names = append(names, e.Name())
	}
	return names
}

// A session configuration that holds each key of git-config(1) that starts a
// program runs none of them in FetchSession (session.md, The fetch by SHA). The
// control shows, one plain read at a time in a copy of the same clone, which
// keys are live on this host.
func TestAHostileSessionRunsNothing(t *testing.T) {
	home := isolate(t)
	runClone := t.TempDir()
	base := commitTree(t, runClone, map[string]string{"a.txt": "a\n", "b.y": "y\n"})
	up := t.TempDir()
	gitOK(t, "", plain(home), "clone", "-q", runClone, up)
	write(t, up, map[string]string{"u.txt": "u\n"})
	gitOK(t, up, plain(home), "add", "u.txt")
	gitOK(t, up, plain(home), "commit", "-q", "-m", "chore: up")

	t.Run("control: plain reads in the hostile clone fire each live key", func(t *testing.T) {
		dir, _ := session(t, home, runClone, base)
		marks := t.TempDir()
		arm(t, dir, marks, up)
		env := plain(home)
		// core.alternateRefsCommand is read in a repository that has an
		// alternate, so this copy gets one before its fetch of up.
		write(t, dir, map[string]string{".git/objects/info/alternates": filepath.Join(runClone, ".git", "objects") + "\n",
			"a.txt": "changed\n", "t.txt": "t\n", "z.txt": "z\n", "b.y": "y\n", "c.p": "p\n"})
		gitIn(dir, env, "add", "-N", "t.txt", "z.txt")
		for _, read := range [][]string{{"status"}, {"diff", "--", "a.txt"}, {"diff", "--no-ext-diff", "--", "t.txt"}, {"diff", "--", "z.txt"}, {"add", "b.y"}, {"add", "c.p"}, {"commit", "-q", "-m", "x", "--", "a.txt"},
			{"x"}, {"z"}, {"y"}, {"fetch", "origin"}, {"fetch", "ext"}, {"fetch", "up"}, {"fetch", "gp"}} {
			gitIn(dir, env, read...) // a read may fail: the marker is the evidence
		}
		must(t, os.Remove(filepath.Join(dir, "b.y")))
		gitIn(dir, env, "checkout", "--", "b.y")
		got := marked(t, marks)
		for _, k := range hostileKeys {
			if !slices.Contains(notLive, k) && !slices.Contains(got, k) {
				t.Errorf("the key %s did not fire in the control; the list is not live on this host (fired: %v)", k, got)
			}
		}
	})

	dir, head := session(t, home, runClone, base)
	marks := t.TempDir()
	arm(t, dir, marks, up)
	must(t, FetchSession(runClone, filepath.Join(dir, ".git"), head, "refs/heads/task/T-ab12/1"))
	if got := marked(t, marks); len(got) != 0 {
		t.Errorf("FetchSession ran the session's programs: %v", got)
	}
}
