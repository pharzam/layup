// Package git is the one caller of the git program (docs/spec/packages.md,
// rule 3 and "The calls of internal/git"). Each call runs git with a fixed
// environment and fixed -c values, so no file and no variable of the host
// changes a branch, a tree, a file list, a hook, an author or a signature,
// and no call asks a question.
package git

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// MinVersion is the oldest git that this package supports: the first with
// GIT_CONFIG_GLOBAL.
const MinVersion = "2.32.0"

// Identity is the author and the committer of a commit, and its time.
type Identity struct {
	Name, Email string
	Time        time.Time
}

// NotFoundError says that the git program is not on the PATH.
type NotFoundError struct{ Err error }

func (e *NotFoundError) Error() string { return "git not found: " + e.Err.Error() }
func (e *NotFoundError) Unwrap() error { return e.Err }

// FailedError says that a call of git failed: git exited with a code that is
// not 0, did not start or was stopped (Code -1), or exited 0 with an output
// that the call cannot read (Code 0).
type FailedError struct {
	Args   []string // the arguments after the -c values
	Code   int
	Stderr string
	Err    error
}

func (e *FailedError) Error() string {
	msg := fmt.Sprintf("git %s: %v", strings.Join(e.Args, " "), e.Err)
	if s := strings.TrimSpace(e.Stderr); s != "" {
		msg += ": " + s
	}
	return msg
}

func (e *FailedError) Unwrap() error { return e.Err }

// config is the -c values that every call starts with.
var config = strings.Fields("-c core.hooksPath=/dev/null -c core.attributesFile=/dev/null " +
	"-c core.excludesFile=/dev/null -c core.autocrlf=false -c commit.gpgsign=false -c http.emptyAuth=false")

// environ is the environment of a call: fixed values, PATH and TMPDIR of the
// host when they are set, and extra. No other variable of the host reaches
// git. HOME is a path under which no file can be, so libcurl reads no .netrc;
// git starts no ssh and no remote helper, which find the user's home without
// HOME.
func environ(extra ...string) []string {
	env := []string{"LC_ALL=C", "HOME=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0", "GIT_ALLOW_PROTOCOL=file:git:http:https"}
	for _, k := range []string{"PATH", "TMPDIR"} {
		if v, ok := os.LookupEnv(k); ok {
			env = append(env, k+"="+v)
		}
	}
	return append(env, extra...)
}

// run starts git with no standard input, and gives its standard output and
// standard error. The unit tests replace it.
var run = func(dir string, args, env []string) ([]byte, []byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir, cmd.Env = dir, env
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}

// call runs git in dir and gives its standard output, or an error of one of
// the two kinds.
func call(dir string, env []string, args ...string) ([]byte, error) {
	stdout, stderr, err := run(dir, append(append([]string{}, config...), args...), env)
	if err == nil {
		return stdout, nil
	}
	if errors.Is(err, exec.ErrNotFound) {
		return nil, &NotFoundError{Err: err}
	}
	code := -1
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		code = exit.ExitCode()
	}
	return nil, &FailedError{Args: args, Code: code, Stderr: string(stderr), Err: err}
}

// do runs git in dir for a call that reads no output.
func do(dir string, args ...string) error {
	_, err := call(dir, environ(), args...)
	return err
}

// Version gives the version of git, for example "2.54.0 (Apple Git-157)".
func Version() (string, error) {
	out, err := call("", environ(), "--version")
	return strings.TrimPrefix(strings.TrimSpace(string(out)), "git version "), err
}

// Supported reports whether a version from Version is 2.32 or newer
// (MinVersion). A version that it cannot read is not supported.
func Supported(version string) bool {
	var major, minor int
	n, _ := fmt.Sscanf(version, "%d.%d", &major, &minor)
	return n == 2 && (major > 2 || major == 2 && minor >= 32)
}

// LsRemote gives the commit of ref in the repository at url.
func LsRemote(url, ref string) (string, error) {
	args := []string{"ls-remote", "--exit-code", "--", url, ref}
	out, err := call("", environ(), args...)
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if id, name, ok := strings.Cut(line, "\t"); ok && name == ref {
			return id, nil
		}
	}
	return "", &FailedError{Args: args, Err: fmt.Errorf("the output has no line for %s", ref)}
}

// Clone clones the repository at url into dir, and checks out no file.
func Clone(url, dir string) error { return do("", "clone", "--no-checkout", "--", url, dir) }

// CheckoutDetach checks out rev, on no branch.
func CheckoutDetach(dir, rev string) error {
	return do(dir, "checkout", "--detach", "--end-of-options", rev)
}

// Init makes an empty repository at dir, on the branch main.
func Init(dir string) error { return do("", "init", "-b", "main", "--", dir) }

// Add stages each change under the paths, deletions too; no path is the
// whole tree.
func Add(dir string, paths ...string) error {
	return do(dir, append([]string{"add", "--all", "--"}, paths...)...)
}

// Commit commits the staged changes. who is the author and the committer, and
// its time is both dates.
func Commit(dir, message string, who Identity) error {
	date := fmt.Sprintf("@%d +0000", who.Time.Unix())
	env := environ("GIT_AUTHOR_NAME="+who.Name, "GIT_AUTHOR_EMAIL="+who.Email, "GIT_AUTHOR_DATE="+date,
		"GIT_COMMITTER_NAME="+who.Name, "GIT_COMMITTER_EMAIL="+who.Email, "GIT_COMMITTER_DATE="+date)
	_, err := call(dir, env, "commit", "-m", message)
	return err
}

// SwitchCreate makes the branch at start and puts the work tree on it.
func SwitchCreate(dir, branch, start string) error {
	return do(dir, "switch", "-c", branch, "--end-of-options", start)
}

// Branch makes the branch name at start, with no switch.
func Branch(dir, name, start string) error {
	return do(dir, "branch", "--end-of-options", name, start)
}

// SwitchOrphan puts the work tree on a new branch with no commit; git removes
// the tracked files from the work tree.
func SwitchOrphan(dir, branch string) error { return do(dir, "switch", "--orphan", branch) }

// RevParse gives the object name of rev, for example of "<commit>^{tree}".
func RevParse(dir, rev string) (string, error) {
	out, err := call(dir, environ(), "rev-parse", "--verify", "--end-of-options", rev)
	return strings.TrimSpace(string(out)), err
}

// RootCommits gives the commits with no parent that rev reaches.
func RootCommits(dir, rev string) ([]string, error) {
	out, err := call(dir, environ(), "rev-list", "--max-parents=0", "--end-of-options", rev, "--")
	return strings.Fields(string(out)), err
}

// LsFiles gives the paths of the tracked files.
func LsFiles(dir string) ([]string, error) {
	out, err := call(dir, environ(), "ls-files", "-z")
	return names(out), err
}

// WorktreeAdd makes a work tree at path with rev checked out, on no branch.
func WorktreeAdd(dir, path, rev string) error {
	return do(dir, "worktree", "add", "--detach", "--", path, rev)
}

// WorktreeRemove removes the work tree at path, with its changes.
func WorktreeRemove(dir, path string) error {
	return do(dir, "worktree", "remove", "--force", "--", path)
}

// Show gives the bytes of the file at path in rev.
func Show(dir, rev, path string) ([]byte, error) {
	return call(dir, environ(), "show", "--end-of-options", rev+":"+path, "--")
}

// DiffNames gives the paths that differ between base and head; a renamed path
// at both ends.
func DiffNames(dir, base, head string) ([]string, error) {
	out, err := call(dir, environ(), "diff", "--name-only", "--no-renames", "-z", "--end-of-options", base, head, "--")
	return names(out), err
}

// Apply applies the patch file to the work tree at dir.
func Apply(dir, patch string) error { return do(dir, "apply", "--", patch) }

// names reads a list of paths that git wrote with -z.
func names(out []byte) []string {
	if s := strings.TrimSuffix(string(out), "\x00"); s != "" {
		return strings.Split(s, "\x00")
	}
	return nil
}
