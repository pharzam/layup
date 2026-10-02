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
	must(t, Clone("file://"+bare, dst))
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
		first := strings.Fields("-c core.hooksPath=/dev/null -c core.autocrlf=false -c commit.gpgsign=false") // D3 before the review
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
		go func() { done <- Clone(url, filepath.Join(t.TempDir(), "target")) }()
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
