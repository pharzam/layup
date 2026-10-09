package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// wantConfig is the -c values that every call starts with.
var wantConfig = strings.Fields("-c core.hooksPath=/dev/null -c core.attributesFile=/dev/null " +
	"-c core.excludesFile=/dev/null -c core.autocrlf=false -c core.precomposeUnicode=false -c commit.gpgsign=false " +
	"-c http.emptyAuth=false -c maintenance.auto=false")

// fullID is an object ID of SHA-1, 40 hexadecimal characters.
const fullID = "0123456789abcdef0123456789abcdef01234567"

// who is the identity of the test commits, at 2026-01-01T00:00:00Z.
var who = Identity{Name: "LAYUP test", Email: "test@layup.invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

// A stubCall is one start of git that the stub saw.
type stubCall struct {
	dir       string
	args, env []string
}

// stub replaces the git process until the test ends: each start records its
// directory, arguments and environment, and gets stdout and err.
func stub(t *testing.T, stdout string, err error) *[]stubCall {
	t.Helper()
	calls, saved := &[]stubCall{}, run
	run = func(dir string, args, env []string) ([]byte, []byte, error) {
		*calls = append(*calls, stubCall{dir, args, env})
		return []byte(stdout), []byte("fatal: stub\n"), err
	}
	t.Cleanup(func() { run = saved })
	return calls
}

// only gives the one start of git that the stub saw.
func only(t *testing.T, calls *[]stubCall) stubCall {
	t.Helper()
	if len(*calls) != 1 {
		t.Fatalf("%d starts of git, want 1", len(*calls))
	}
	return (*calls)[0]
}

func TestEachCallRunsItsVerb(t *testing.T) {
	url := "https://example.invalid/b.git"
	tests := []struct {
		name, dir, args string // args: after the -c values, split at spaces
		call            func()
		env             []string // after the fixed list
	}{
		{"version", "", "--version", func() { Version() }, nil},
		{"ls-remote", "", "ls-remote --exit-code -- " + url + " HEAD", func() { LsRemote(url, "HEAD") }, nil},
		{"clone", "", "clone --no-checkout -- " + url + " w/target", func() { Clone(url, "w/target", Auth{}) }, nil},
		{"clone with a token", "", "clone --no-checkout -- " + url + " w/target", func() { Clone(url, "w/target", testAuth) }, tokenEnv},
		{"checkout --detach", "r", "checkout --detach " + fullID, func() { CheckoutDetach("r", fullID) }, nil},
		{"init -b main", "", "init -b main -- r", func() { Init("r") }, nil},
		{"add, the whole tree", "r", "add --all --", func() { Add("r") }, nil},
		{"add, two paths", "r", "add --all -- a.txt -b", func() { Add("r", "a.txt", "-b") }, nil},
		{"commit", "r", "commit -m chore:S04", func() { Commit("r", "chore:S04", who) }, []string{"GIT_AUTHOR_NAME=LAYUP test",
			"GIT_AUTHOR_EMAIL=test@layup.invalid", "GIT_AUTHOR_DATE=@1767225600 +0000", "GIT_COMMITTER_NAME=LAYUP test",
			"GIT_COMMITTER_EMAIL=test@layup.invalid", "GIT_COMMITTER_DATE=@1767225600 +0000"}},
		{"switch -c", "r", "switch -c layup-setup " + fullID, func() { SwitchCreate("r", "layup-setup", fullID) }, nil},
		{"switch --orphan", "r", "switch --orphan layup-records", func() { SwitchOrphan("r", "layup-records") }, nil},
		{"reset --soft", "r", "reset --soft " + fullID, func() { ResetSoft("r", fullID) }, nil},
		{"reset --hard", "r", "reset --hard --quiet HEAD", func() { ResetHard("r") }, nil},
		{"diff --cached --quiet", "r", "diff --cached --quiet --exit-code", func() { Staged("r") }, nil},
		{"symbolic-ref --quiet", "r", "symbolic-ref --quiet HEAD", func() { Branch("r") }, nil},
		{"rev-parse", "r", "rev-parse --verify --end-of-options abc^{tree}", func() { RevParse("r", "abc^{tree}") }, nil},
		{"rev-list --max-parents=0", "r", "rev-list --max-parents=0 --end-of-options main --", func() { RootCommits("r", "main") }, nil},
		{"log -1 --format=%B", "r", "log -1 --format=%B --end-of-options main --", func() { Message("r", "main") }, nil},
		{"rev-parse --is-shallow-repository", "r", "rev-parse --is-shallow-repository", func() { IsShallow("r") }, nil},
		{"ls-files", "r", "ls-files -z", func() { LsFiles("r") }, nil},
		{"worktree add --detach", "r", "worktree add --detach -- /s/scratch abc", func() { WorktreeAdd("r", "/s/scratch", "abc") }, nil},
		{"worktree remove", "r", "worktree remove --force -- /s/scratch", func() { WorktreeRemove("r", "/s/scratch") }, nil},
		{"show", "r", "show --end-of-options abc:docs/gates.tsv --", func() { Show("r", "abc", "docs/gates.tsv") }, nil},
		{"diff --name-only", "r", "diff --name-only --no-renames -z --end-of-options abc def --", func() { DiffNames("r", "abc", "def") }, nil},
		{"apply", "r", "apply -- /s/static.patch", func() { Apply("r", "/s/static.patch") }, nil},
		{"ls-tree", "r", "ls-tree -r -z --full-tree --end-of-options abc -- layout", func() { LsTree("r", "abc", "layout") }, nil},
		{"fetch", "r", "fetch --no-tags --update-head-ok -- " + url + " refs/heads/main:refs/heads/main",
			func() { Fetch("r", url, "refs/heads/main", Auth{}) }, nil},
		{"fetch with a token", "r", "fetch --no-tags --update-head-ok -- " + url + " refs/heads/main:refs/heads/main",
			func() { Fetch("r", url, "refs/heads/main", testAuth) }, tokenEnv},
		{"push", "r", "push --porcelain -- " + url + " " + fullID + ":refs/heads/layup-records",
			func() { Push("r", url, fullID, "layup-records", Auth{}) }, nil},
		{"push with a token", "r", "push --porcelain -- " + url + " " + fullID + ":refs/heads/layup-records",
			func() { Push("r", url, fullID, "layup-records", testAuth) }, tokenEnv},
		{"merge-base --is-ancestor", "r", "merge-base --is-ancestor " + fullID + " " + fullID, func() { IsAncestor("r", fullID, fullID) }, nil},
		{"diff -U0 of one file", "r", "diff -U0 --no-color --no-renames --end-of-options abc def -- docs/guardrails.md",
			func() { DiffFile("r", "abc", "def", "docs/guardrails.md") }, nil},
		{"diff --binary", "r", "diff --binary --no-renames --end-of-options abc def --", func() { DiffBinary("r", "abc", "def") }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := stub(t, "", nil)
			tt.call()
			c := only(t, calls)
			if want := append(append([]string{}, wantConfig...), strings.Fields(tt.args)...); c.dir != tt.dir || !reflect.DeepEqual(c.args, want) {
				t.Errorf("dir %q, args %q\nwant dir %q, args %q", c.dir, c.args, tt.dir, want)
			}
			if want := environ(tt.env...); !reflect.DeepEqual(c.env, want) {
				t.Errorf("env %q\nwant the fixed list, then %q", c.env, tt.env)
			}
		})
	}
}

// The environment is a fixed list: no other variable of the host, above all
// no GIT_* variable, no askpass program and not its HOME, reaches git. This
// test is the proof of input (e) of TestAHostileHostChangesNothing.
func TestTheEnvironmentIsAFixedList(t *testing.T) {
	for k, v := range map[string]string{"PATH": "/stub/bin", "HOME": "/stub/home", "TMPDIR": "/stub/tmp"} {
		t.Setenv(k, v)
	}
	for _, k := range strings.Fields("GIT_DIR GIT_WORK_TREE GIT_AUTHOR_NAME GIT_CONFIG_COUNT GIT_CONFIG_KEY_0 GIT_CONFIG_VALUE_0 " +
		"GIT_CONFIG_GLOBAL GIT_SSH_COMMAND GIT_ALLOW_PROTOCOL GIT_ASKPASS SSH_ASKPASS SSH_AUTH_SOCK XDG_CONFIG_HOME LANG") {
		t.Setenv(k, "hostile")
	}
	calls := stub(t, "", nil)
	RevParse("r", "HEAD")
	fixed := []string{"LC_ALL=C", "HOME=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0", "GIT_ALLOW_PROTOCOL=file:git:http:https", "PATH=/stub/bin", "TMPDIR=/stub/tmp"}
	if env := only(t, calls).env; !reflect.DeepEqual(env, fixed) {
		t.Fatalf("env\n got %q\nwant %q", env, fixed)
	}

	// A commit adds its identity, with both dates, so the same input gives
	// the same commit ID (NFR-005).
	calls = stub(t, "", nil)
	Commit("r", "chore: one", who)
	identity := []string{"GIT_AUTHOR_NAME=LAYUP test", "GIT_AUTHOR_EMAIL=test@layup.invalid", "GIT_AUTHOR_DATE=@1767225600 +0000",
		"GIT_COMMITTER_NAME=LAYUP test", "GIT_COMMITTER_EMAIL=test@layup.invalid", "GIT_COMMITTER_DATE=@1767225600 +0000"}
	if env := only(t, calls).env; !reflect.DeepEqual(env, append(fixed, identity...)) {
		t.Fatalf("commit env\n got %q\nwant the fixed list, then %q", env, identity)
	}

	// A variable of the list that the host does not set stays unset.
	os.Unsetenv("TMPDIR")
	calls = stub(t, "", nil)
	RevParse("r", "HEAD")
	if env := only(t, calls).env; !reflect.DeepEqual(env, fixed[:8]) {
		t.Fatalf("env with no TMPDIR\n got %q\nwant %q", env, fixed[:8])
	}
}

func TestTheOutputIsRead(t *testing.T) {
	for _, c := range []struct {
		name, stdout string
		call         func() (any, error)
		want         any
	}{
		{"ls-remote: the line of exactly HEAD", "1111\trefs/remotes/origin/HEAD\n2222\tHEAD\n",
			func() (any, error) { return LsRemote("u", "HEAD") }, "2222"},
		{"rev-parse", "abc\n", func() (any, error) { return RevParse("r", "HEAD") }, "abc"},
		{"rev-list", "r1\nr2\n", func() (any, error) { return RootCommits("r", "main") }, []string{"r1", "r2"}},
		{"log: the message with no line feed at its end", "chore: a\n\nb\n\n", func() (any, error) { return Message("r", "main") }, "chore: a\n\nb"},
		{"ls-files: the names as git wrote them", "a b.txt\x00c\nd.txt\x00",
			func() (any, error) { return LsFiles("r") }, []string{"a b.txt", "c\nd.txt"}},
		{"diff: no change", "", func() (any, error) { return DiffNames("r", "a", "b") }, []string(nil)},
		{"show: the bytes unchanged", "a\r\nb", func() (any, error) { return Show("r", "HEAD", "a.txt") }, []byte("a\r\nb")},
		{"version", "git version 2.54.0 (Apple Git-157)\n", func() (any, error) { return Version() }, "2.54.0 (Apple Git-157)"},
		{"is-shallow: a shallow clone", "true\n", func() (any, error) { return IsShallow("r") }, true},
		{"is-shallow: a full history", "false\n", func() (any, error) { return IsShallow("r") }, false},
	} {
		stub(t, c.stdout, nil)
		if got, err := c.call(); err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %q, %v; want %q", c.name, got, err, c.want)
		}
	}
	stub(t, "1111\trefs/remotes/origin/HEAD\n", nil)
	var failed *FailedError
	if _, err := LsRemote("u", "HEAD"); !errors.As(err, &failed) || failed.Code != 0 {
		t.Errorf("LsRemote with no line for HEAD: %v; want a *FailedError with code 0", err)
	}
	stub(t, "yes\n", nil)
	if shallow, err := IsShallow("r"); !errors.As(err, &failed) || failed.Code != 0 || shallow {
		t.Errorf("IsShallow with the output yes: %v, %v; want false and a *FailedError with code 0", shallow, err)
	}
}

// Staged reads the exit code of git diff --cached --quiet: 0 is no staged
// change, 1 a staged change, another code an error (D12 of #90).
func TestStagedReadsTheExitCode(t *testing.T) {
	for _, c := range []struct {
		code   int
		staged bool
		failed bool
	}{{0, false, false}, {1, true, false}, {128, false, true}} {
		var err error
		if c.code != 0 {
			err = exec.Command("sh", "-c", fmt.Sprintf("exit %d", c.code)).Run()
		}
		stub(t, "", err)
		staged, got := Staged("r")
		if staged != c.staged || (got != nil) != c.failed {
			t.Errorf("exit %d: %v, %v; want %v and an error %v", c.code, staged, got, c.staged, c.failed)
		}
	}
}

// CheckoutDetach and SwitchCreate take a full object ID and no
// --end-of-options; any other text, which could be an option, starts no git.
func TestCheckoutAndSwitchTakeAFullObjectID(t *testing.T) {
	for _, rev := range []string{"-x", "HEAD", fullID[:39], fullID + "0", strings.ToUpper(fullID)} {
		calls := stub(t, "", nil)
		var failed *FailedError
		for _, err := range []error{CheckoutDetach("r", rev), SwitchCreate("r", "b", rev)} {
			if !errors.As(err, &failed) || failed.Code != -1 {
				t.Errorf("%q: %v; want a *FailedError with code -1", rev, err)
			}
		}
		if len(*calls) != 0 {
			t.Errorf("%q: %d starts of git, want none", rev, len(*calls))
		}
	}
	calls := stub(t, "", nil)
	if err := CheckoutDetach("r", strings.Repeat("ab", 32)); err != nil || len(*calls) != 1 {
		t.Errorf("an object ID of SHA-256: %v, %d starts of git; want no error and 1 start", err, len(*calls))
	}
}

func TestSupported(t *testing.T) {
	for v, want := range map[string]bool{"2.32.0": true, "2.54.0 (Apple Git-157)": true, "2.39.5.windows.1": true,
		"2.32": true, "3.0.0": true, "2.31.9": false, "2.28.0": false, "1.99.0": false, "": false, "unknown": false} {
		if got := Supported(v); got != want {
			t.Errorf("Supported(%q) = %v, want %v", v, got, want)
		}
	}
}

func TestTheTwoErrorKinds(t *testing.T) {
	stub(t, "", &exec.Error{Name: "git", Err: exec.ErrNotFound})
	var notFound *NotFoundError
	if _, err := RevParse("r", "HEAD"); !errors.As(err, &notFound) {
		t.Errorf("git missing: %v; want a *NotFoundError", err)
	}
	stub(t, "", errors.New("signal: killed"))
	var failed *FailedError
	if _, err := RevParse("r", "HEAD"); !errors.As(err, &failed) {
		t.Fatalf("git failed: %v; want a *FailedError", err)
	}
	if failed.Code != -1 || failed.Stderr != "fatal: stub\n" || strings.Join(failed.Args, " ") != "rev-parse --verify --end-of-options HEAD" {
		t.Errorf("FailedError %+v: want code -1, the standard error, and the arguments after the -c values", failed)
	}
}

// LsTree reads each record of git ls-tree -z: the mode, the type, the object
// and the path, with a tab before the path; a record of another form is an
// error of the call, code 0.
func TestLsTreeReadsEachEntry(t *testing.T) {
	stub(t, "100644 blob "+fullID+"\tlayout/a test.go\x00100755 blob "+fullID+"\trun.sh\x00", nil)
	got, err := LsTree("r", "abc", "layout")
	want := []TreeEntry{{"100644", "blob", fullID, "layout/a test.go"}, {"100755", "blob", fullID, "run.sh"}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("LsTree = %q, %v; want %q", got, err, want)
	}
	stub(t, "", nil)
	if got, err := LsTree("r", "abc", "missing"); err != nil || len(got) != 0 {
		t.Fatalf("LsTree of no entry = %q, %v; want none", got, err)
	}
	stub(t, "100644 blob\tx\x00", nil)
	var failed *FailedError
	if _, err := LsTree("r", "abc", "x"); !errors.As(err, &failed) || failed.Code != 0 {
		t.Fatalf("LsTree of a record of another form: %v; want a *FailedError with code 0", err)
	}
}

// testAuth is the web of a forge register and a token of the tests; tokenEnv
// is what a call with it adds to the fixed list (docs/spec/forge.md, The App
// identity): the header of x-access-token and the token, in Basic.
var (
	testAuth = Auth{Web: "https://github.com", Token: "ghs_testtoken"}
	tokenEnv = []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=http.https://github.com/.extraHeader",
		"GIT_CONFIG_VALUE_0=Authorization: Basic eC1hY2Nlc3MtdG9rZW46Z2hzX3Rlc3R0b2tlbg=="}
)

func TestFetchAndPushRefuseTheirInputBeforeGitStarts(t *testing.T) {
	url := "https://example.invalid/b.git"
	for name, call := range map[string]func() error{
		"a ref that forces":              func() error { return Fetch("r", url, "+refs/heads/main", Auth{}) },
		"a ref with a colon":             func() error { return Fetch("r", url, "refs/heads/a:refs/heads/b", Auth{}) },
		"a ref not under refs/":          func() error { return Fetch("r", url, "main", Auth{}) },
		"a commit that is not a full ID": func() error { return Push("r", url, "+abc", "main", Auth{}) },
		"a branch with a colon":          func() error { return Push("r", url, fullID, "a:b", Auth{}) },
		"an empty branch":                func() error { return Push("r", url, fullID, "", Auth{}) },
		"a token with no web":            func() error { return Fetch("r", url, "refs/heads/main", Auth{Token: "ghs_x"}) },
		"a web that ends in a slash": func() error {
			return Fetch("r", url, "refs/heads/main", Auth{Web: "https://github.com/", Token: "ghs_x"})
		},
	} {
		calls := stub(t, "", nil)
		err := call()
		var failed *FailedError
		if !errors.As(err, &failed) || failed.Code != -1 || len(*calls) != 0 {
			t.Errorf("%s: %v, %d starts of git; want a *FailedError of code -1 and none", name, err, len(*calls))
		}
	}
}

func TestNoErrorHoldsTheToken(t *testing.T) {
	stub(t, "", errors.New("exit status 128"))
	for _, err := range []error{Fetch("r", "https://github.com/a/b.git", "refs/heads/main", testAuth),
		Push("r", "https://github.com/a/b.git", fullID, "main", testAuth)} {
		if err == nil || strings.Contains(err.Error(), "ghs_testtoken") || strings.Contains(err.Error(), "eC1hY2Nlc3MtdG9rZW46") {
			t.Errorf("the error %v is nil or holds the token", err)
		}
	}
}

// CloneLocal clones one branch with no remote: the clone, then the removal of
// its remote in the clone (packages.md, The calls of M2b).
func TestCloneLocalRemovesItsRemote(t *testing.T) {
	calls := stub(t, "", nil)
	if err := CloneLocal("/w/run", "/w/s/repo", "main"); err != nil {
		t.Fatal(err)
	}
	want := []stubCall{
		{"", append(append([]string{}, wantConfig...), strings.Fields("clone --no-local --no-checkout --single-branch --no-tags --branch main -- /w/run /w/s/repo")...), environ()},
		{"/w/s/repo", append(append([]string{}, wantConfig...), "remote", "remove", "origin"), environ()},
	}
	if !reflect.DeepEqual(*calls, want) {
		t.Fatalf("calls %q\nwant %q", *calls, want)
	}
}

// FetchSession runs no git in the session's clone: a scratch bare repository
// whose alternate is the session's object directory, the type of the SHA,
// the ref, and the fetch with an empty hooks directory; both directories are
// gone on each return (packages.md, The calls of M2b).
func TestFetchSessionRunsNoGitInTheSession(t *testing.T) {
	session := t.TempDir()
	calls := stub(t, "commit\n", nil)
	if err := FetchSession("r", session, fullID, "refs/heads/task/T-ab12/1"); err != nil {
		t.Fatal(err)
	}
	if len(*calls) != 4 {
		t.Fatalf("%d starts of git, want 4: %q", len(*calls), *calls)
	}
	tmp := (*calls)[0].args[len((*calls)[0].args)-1]
	var empty string
	for _, a := range (*calls)[3].args {
		if v, ok := strings.CutPrefix(a, "core.hooksPath="); ok && v != "/dev/null" {
			empty = v
		}
	}
	steps := []struct{ dir, args string }{
		{"", "init --bare -- " + tmp},
		{tmp, "cat-file -t " + fullID},
		{tmp, "update-ref refs/heads/session " + fullID},
		{"r", "-c core.hooksPath=" + empty + " fetch --no-tags --no-write-fetch-head -- " + tmp + " refs/heads/session:refs/heads/task/T-ab12/1"},
	}
	for i, st := range steps {
		if c := (*calls)[i]; c.dir != st.dir || !reflect.DeepEqual(c.args, append(append([]string{}, wantConfig...), strings.Fields(st.args)...)) {
			t.Errorf("start %d: dir %q, args %q\nwant dir %q, args %q", i+1, c.dir, c.args, st.dir, st.args)
		}
		if c := (*calls)[i]; c.dir == session || slices.Contains(c.args, session) {
			t.Errorf("start %d runs git in the session's clone: %q", i+1, c)
		}
	}
	if empty == "" || tmp == "" {
		t.Fatal("no scratch repository or no empty hooks directory")
	}
	if _, err := os.Stat(tmp); err == nil {
		t.Errorf("the scratch repository %s is left", tmp)
	}
	if _, err := os.Stat(empty); err == nil {
		t.Errorf("the empty hooks directory %s is left", empty)
	}
}

// A SHA that names no commit is refused with ErrNotACommit in a *FailedError,
// whether cat-file prints another type or fails, and no ref is set.
func TestFetchSessionRefusesASHAOfNoCommit(t *testing.T) {
	failed := exec.Command("sh", "-c", "exit 128").Run()
	for name, c := range map[string]struct {
		stdout string
		err    error
	}{"a blob": {"blob\n", nil}, "an absent object": {"", failed}} {
		calls := &[]stubCall{}
		saved := run
		run = func(dir string, args, env []string) ([]byte, []byte, error) {
			*calls = append(*calls, stubCall{dir, args, env})
			if slices.Contains(args, "cat-file") {
				return []byte(c.stdout), []byte("fatal: stub\n"), c.err
			}
			return nil, nil, nil
		}
		err := FetchSession("r", t.TempDir(), fullID, "refs/heads/task/T-ab12/1")
		run = saved
		var fe *FailedError
		if !errors.Is(err, ErrNotACommit) || !errors.As(err, &fe) {
			t.Errorf("%s: %v; want ErrNotACommit in a *FailedError", name, err)
		}
		if len(*calls) != 2 {
			t.Errorf("%s: %d starts of git, want 2 (no ref, no fetch)", name, len(*calls))
		}
		if tmp := (*calls)[0].args[len((*calls)[0].args)-1]; present(tmp) {
			t.Errorf("%s: the scratch repository %s is left", name, tmp)
		}
	}
}

func present(path string) bool { _, err := os.Stat(path); return err == nil }

func TestTheCallsOfM2bRefuseTheirInputBeforeGitStarts(t *testing.T) {
	for name, call := range map[string]func() error{
		"an empty branch to clone":            func() error { return CloneLocal("/w/run", "/w/s/repo", "") },
		"a branch to clone that is an option": func() error { return CloneLocal("/w/run", "/w/s/repo", "-u/bin/sh") },
		"a SHA that is not a full ID":         func() error { return FetchSession("r", "/w/s/repo/.git", "HEAD", "refs/heads/x") },
		"a destination not under refs/":       func() error { return FetchSession("r", "/w/s/repo/.git", fullID, "main") },
		"a destination with a colon":          func() error { return FetchSession("r", "/w/s/repo/.git", fullID, "refs/heads/a:b") },
		"a base that is not a full ID":        func() error { _, err := IsAncestor("r", "--all", fullID); return err },
		"a head that is not a full ID":        func() error { _, err := IsAncestor("r", fullID, "main"); return err },
	} {
		calls := stub(t, "", nil)
		err := call()
		var failed *FailedError
		if !errors.As(err, &failed) || failed.Code != -1 || len(*calls) != 0 {
			t.Errorf("%s: %v, %d starts of git; want a *FailedError of code -1 and none", name, err, len(*calls))
		}
	}
}

// IsAncestor reads the exit code: 0 is yes, 1 is no, any other is an error.
func TestIsAncestorReadsTheExitCode(t *testing.T) {
	for _, c := range []struct {
		code       int
		yes, isErr bool
	}{{0, true, false}, {1, false, false}, {128, false, true}} {
		var err error
		if c.code != 0 {
			err = exec.Command("sh", "-c", fmt.Sprintf("exit %d", c.code)).Run()
		}
		stub(t, "", err)
		yes, got := IsAncestor("r", fullID, fullID)
		if yes != c.yes || (got != nil) != c.isErr {
			t.Errorf("exit %d: %v, %v; want %v and an error %v", c.code, yes, got, c.yes, c.isErr)
		}
	}
}
