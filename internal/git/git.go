// Package git is the one caller of the git program (docs/spec/packages.md,
// rule 3 and "The calls of internal/git"). Each call runs git with a fixed
// environment and fixed -c values, so no file and no variable of the host
// changes a branch, a tree, a file list, a hook, an author or a signature,
// and no call asks a question or uses a credential of the host.
package git

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
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
// not 0, did not start (also when the call refused its input) or was stopped
// (Code -1), or exited 0 with an output that the call cannot read (Code 0).
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

// config is the -c values that every call starts with. maintenance.auto=false
// keeps a commit from starting the maintenance of git, whose tasks go on in
// the background after the call ends: in CI of the pull request #115, on git
// 2.55.0, one repacked (task T-d6q5).
var config = strings.Fields("-c core.hooksPath=/dev/null -c core.attributesFile=/dev/null " +
	"-c core.excludesFile=/dev/null -c core.autocrlf=false -c core.precomposeUnicode=false -c commit.gpgsign=false " +
	"-c http.emptyAuth=false -c maintenance.auto=false")

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
	if id, ok := refLine(out, ref); ok {
		return id, nil
	}
	return "", &FailedError{Args: args, Err: fmt.Errorf("the output has no line for %s", ref)}
}

// Clone clones the repository at url into dir, and checks out no file.
// With a token, auth gives git the header of the token, as Fetch and Push do.
func Clone(url, dir string, auth Auth) error {
	return doAuth("", auth, []string{"clone", "--no-checkout", "--", url, dir})
}

// CheckoutDetach checks out commit, a full object ID such as LsRemote gives,
// on no branch.
func CheckoutDetach(dir, commit string) error {
	return doAt(dir, commit, "checkout", "--detach", commit)
}

// Init makes an empty repository at dir, on the branch main.
func Init(dir string) error { return do("", "init", "-b", "main", "--", dir) }

// Add stages each change under the paths, deletions too; no path is the
// whole tree.
func Add(dir string, paths ...string) error {
	return do(dir, append([]string{"add", "--all", "--"}, paths...)...)
}

// commitEnv is the environment of a call that makes a commit: the fixed list,
// then the six variables of who, the author and the committer, and its time as
// both dates, so one input gives one commit ID (NFR-005). Commit and CommitTree
// take it from here.
func commitEnv(who Identity) []string {
	date := fmt.Sprintf("@%d +0000", who.Time.Unix())
	return environ("GIT_AUTHOR_NAME="+who.Name, "GIT_AUTHOR_EMAIL="+who.Email, "GIT_AUTHOR_DATE="+date,
		"GIT_COMMITTER_NAME="+who.Name, "GIT_COMMITTER_EMAIL="+who.Email, "GIT_COMMITTER_DATE="+date)
}

// Commit commits the staged changes. who is the author and the committer, and
// its time is both dates.
func Commit(dir, message string, who Identity) error {
	_, err := call(dir, commitEnv(who), "commit", "-m", message)
	return err
}

// CommitTree makes a commit of tree on parent, in dir, with who as its author
// and its committer, on no ref, and gives its ID: the probe of the forge check
// pushes it (task T-j4sd, #272; docs/spec/packages.md, The calls of M2d). tree
// and parent are full object IDs, and any other text, which could be an option,
// starts no git (so the call needs no --end-of-options); the message is the
// argument of -m. An empty commit is the call with the tree of its parent: the
// caller gives it. An output that is no object ID is a *FailedError with code 0.
func CommitTree(dir, tree, parent, message string, who Identity) (string, error) {
	args := []string{"commit-tree", tree, "-p", parent, "-m", message}
	if !fullObjectID(tree) || !fullObjectID(parent) {
		return "", &FailedError{Args: args, Code: -1, Err: errors.New("not a full object ID")}
	}
	out, err := call(dir, commitEnv(who), args...)
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(out))
	if !fullObjectID(id) {
		return "", &FailedError{Args: args, Code: 0, Err: fmt.Errorf("the output %q is no object ID", id)}
	}
	return id, nil
}

// SwitchCreate makes the branch at commit, a full object ID, and puts the
// work tree on it.
func SwitchCreate(dir, branch, commit string) error {
	return doAt(dir, commit, "switch", "-c", branch, commit)
}

// doAt is do for checkout and switch. They get no --end-of-options, which
// they may read as a revision before git 2.44, so commit must be a full
// object ID: any other text, which could be an option, starts no git.
func doAt(dir, commit string, args ...string) error {
	if !fullObjectID(commit) {
		return &FailedError{Args: args, Code: -1, Err: errors.New("not a full object ID")}
	}
	return do(dir, args...)
}

// Branch gives the branch of the work tree at dir, as refs/heads/<name>; a
// detached HEAD is an error.
func Branch(dir string) (string, error) {
	out, err := call(dir, environ(), "symbolic-ref", "--quiet", "HEAD")
	return strings.TrimSpace(string(out)), err
}

// ResetSoft moves the branch of HEAD in dir to commit, and keeps the index and
// the work tree: the runner of layup setup takes back the commit of a step
// whose evidence did not pass (D12 of #86).
func ResetSoft(dir, commit string) error { return doAt(dir, commit, "reset", "--soft", commit) }

// ResetHard puts the index and the work tree of dir back to HEAD: the runner
// of layup setup starts each step from the commit of the step before it (D11
// of #90).
func ResetHard(dir string) error { return do(dir, "reset", "--hard", "--quiet", "HEAD") }

// Staged reports whether the index of dir holds a change against HEAD: the
// runner of layup setup commits a step only then (D12 of #90).
func Staged(dir string) (bool, error) {
	_, err := call(dir, environ(), "diff", "--cached", "--quiet", "--exit-code")
	var failed *FailedError
	if errors.As(err, &failed) && failed.Code == 1 {
		return true, nil
	}
	return false, err
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

// Message gives the message of the commit rev, with no line feed at its end:
// S03 takes only its own root commit of a run that stopped (D8 of #86).
func Message(dir, rev string) (string, error) {
	out, err := call(dir, environ(), "log", "-1", "--format=%B", "--end-of-options", rev, "--")
	return strings.TrimRight(string(out), "\n"), err
}

// CommitTime gives the committer time of the commit rev, in UTC, read as
// seconds since the epoch (%ct), as the form of %cI differs between versions
// of git (docs/guardrails.md): the time of Start is the first commit of the
// records branch (task T-ybd5).
func CommitTime(dir, rev string) (time.Time, error) {
	out, err := call(dir, environ(), "log", "-1", "--format=%ct", "--end-of-options", rev, "--")
	if err != nil {
		return time.Time{}, err
	}
	n, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("git log: the committer time %q: %w", out, err)
	}
	return time.Unix(n, 0).UTC(), nil
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

// TreeEntry is one entry of LsTree: the mode (for example 100644, 100755,
// 120000 for a symbolic link), the type (blob, or commit for a submodule), the
// object and the path from the root of the tree.
type TreeEntry struct{ Mode, Type, Object, Path string }

// LsTree gives each entry at or under path in rev, recursively, in the order
// of git; a path that rev does not have gives none.
func LsTree(dir, rev, path string) ([]TreeEntry, error) {
	args := []string{"ls-tree", "-r", "-z", "--full-tree", "--end-of-options", rev, "--", path}
	out, err := call(dir, environ(), args...)
	if err != nil {
		return nil, err
	}
	var entries []TreeEntry
	for _, record := range names(out) {
		meta, p, ok := strings.Cut(record, "\t")
		f := strings.Fields(meta)
		if !ok || len(f) != 3 {
			return nil, &FailedError{Args: args, Err: fmt.Errorf("a record of another form: %q", record)}
		}
		entries = append(entries, TreeEntry{Mode: f[0], Type: f[1], Object: f[2], Path: p})
	}
	return entries, nil
}

// Parents gives the parents of the commit rev, in order: none for a root commit.
func Parents(dir, rev string) ([]string, error) {
	args := []string{"rev-list", "--parents", "-n", "1", "--end-of-options", rev, "--"}
	out, err := call(dir, environ(), args...)
	if err != nil {
		return nil, err
	}
	f := strings.Fields(string(out))
	if len(f) == 0 {
		return nil, &FailedError{Args: args, Err: errors.New("the output is empty")}
	}
	return f[1:], nil
}

// BlobID gives the object ID that git gives a blob of data, in a repository of SHA-1 (the object format of
// every target of GitHub): the hash of "blob", the length, a NUL and the data. It starts no process and
// writes no object.
func BlobID(data []byte) string {
	h := sha1.New()
	fmt.Fprintf(h, "blob %d\x00", len(data))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// IsShallow reports whether the repository at dir is a shallow clone.
func IsShallow(dir string) (bool, error) {
	args := []string{"rev-parse", "--is-shallow-repository"}
	out, err := call(dir, environ(), args...)
	if err != nil {
		return false, err
	}
	switch s := strings.TrimSpace(string(out)); s {
	case "true", "false":
		return s == "true", nil
	default:
		return false, &FailedError{Args: args, Err: fmt.Errorf("an output of another form: %q", s)}
	}
}

// Auth is the web of the forge register and the installation token of
// docs/spec/forge.md (The App identity). The zero Auth sends no header.
type Auth struct{ Web, Token string }

// env gives what a call with auth adds to the fixed list: the header of the
// token for http.<web>/, in the call's environment and never in its
// arguments, where another user of the host could read it.
func (a Auth) env() ([]string, error) {
	if a.Token == "" {
		return nil, nil
	}
	if a.Web == "" || strings.HasSuffix(a.Web, "/") {
		return nil, errors.New("a token needs the web of the forge register, with no final slash")
	}
	basic := base64.StdEncoding.EncodeToString([]byte("x-access-token:" + a.Token))
	return []string{"GIT_CONFIG_COUNT=1", "GIT_CONFIG_KEY_0=http." + a.Web + "/.extraHeader",
		"GIT_CONFIG_VALUE_0=Authorization: Basic " + basic}, nil
}

// Fetch fetches ref of url into the same ref of dir, with no tags. ref starts
// with refs/ and holds no ':', so no '+' forces it. --update-head-ok lets it
// set the branch of a repository that Init made, which git otherwise refuses
// as the branch that is checked out (packages.md, The calls of M2a).
func Fetch(dir, url, ref string, auth Auth) error {
	args := []string{"fetch", "--no-tags", "--update-head-ok", "--", url, ref + ":" + ref}
	if !strings.HasPrefix(ref, "refs/") || strings.Contains(ref, ":") {
		return &FailedError{Args: args, Code: -1, Err: errors.New("not a ref under refs/ with no ':'")}
	}
	return doAuth(dir, auth, args)
}

// Push pushes commit, a full object ID, to the branch of url, never with
// --force. A push that is not a fast-forward, or that the remote refuses, is a
// *FailedError of Code 1; a remote that cannot be reached gives another code.
func Push(dir, url, commit, branch string, auth Auth) error {
	args := []string{"push", "--porcelain", "--", url, commit + ":refs/heads/" + branch}
	if !fullObjectID(commit) || branch == "" || strings.Contains(branch, ":") {
		return &FailedError{Args: args, Code: -1, Err: errors.New("not a full object ID, or not a branch")}
	}
	return doAuth(dir, auth, args)
}

// doAuth runs git in dir with the environment of auth, and drops its output.
func doAuth(dir string, auth Auth, args []string) error {
	_, err := callAuth(dir, auth, args)
	return err
}

// callAuth runs git in dir with the environment of auth, and gives its
// output: a token is a header of the call's environment, never an argument.
func callAuth(dir string, auth Auth, args []string) ([]byte, error) {
	extra, err := auth.env()
	if err != nil {
		return nil, &FailedError{Args: args, Code: -1, Err: err}
	}
	return call(dir, environ(extra...), args...)
}

// RemoteHead gives the commit of the branch of the repository at url, with the
// token of auth when it has one: found is false for a branch that the remote
// does not have (exit 2 of ls-remote, or no line of exactly refs/heads/BRANCH,
// as ls-remote matches its pattern against the tail of a ref name). The ID on
// that line must be a full object ID, else the answer is a *FailedError of Code
// 0 and never an absence. An empty branch, one that holds ':' and one that starts
// with '-' are refused before git starts (task T-syxm, packages.md, The calls of
// M2e).
func RemoteHead(url, branch string, auth Auth) (sha string, found bool, err error) {
	ref := "refs/heads/" + branch
	args := []string{"ls-remote", "--exit-code", "--", url, ref}
	if branch == "" || strings.Contains(branch, ":") || strings.HasPrefix(branch, "-") {
		return "", false, &FailedError{Args: args, Code: -1, Err: errors.New("not a branch")}
	}
	out, err := callAuth("", auth, args)
	var failed *FailedError
	if errors.As(err, &failed) && failed.Code == 2 {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	id, ok := refLine(out, ref)
	switch {
	case !ok:
		return "", false, nil
	case !fullObjectID(id):
		return "", false, &FailedError{Args: args, Err: fmt.Errorf("the line for %s holds %q, which is no full object ID", ref, id)}
	}
	return id, true, nil
}

// refLine gives the ID on the line of exactly the ref name in the output of ls-remote.
func refLine(out []byte, ref string) (string, bool) {
	for _, line := range strings.Split(string(out), "\n") {
		if id, name, ok := strings.Cut(line, "\t"); ok && name == ref {
			return id, true
		}
	}
	return "", false
}

// ErrNotACommit is the refusal of FetchSession for a SHA that names no commit
// of the session's objects: the refusal branch of docs/spec/session.md. It is
// the Err of a *FailedError, so errors.As still finds the kind.
var ErrNotACommit = errors.New("the SHA names no commit of the session's objects")

// fullObjectID reports whether s is a full object ID: 40 or 64 lowercase
// hexadecimal characters.
func fullObjectID(s string) bool {
	return (len(s) == 40 || len(s) == 64) && strings.Trim(s, "0123456789abcdef") == ""
}

// CloneLocal clones the one branch of the local repository src into dir, with
// no checkout, no tag and its own objects, then removes its remote, so no
// command in dir reaches a remote through it (packages.md, The calls of M2b).
func CloneLocal(src, dir, branch string) error {
	args := []string{"clone", "--no-local", "--no-checkout", "--single-branch", "--no-tags", "--branch", branch, "--", src, dir}
	if branch == "" || strings.HasPrefix(branch, "-") {
		return &FailedError{Args: args, Code: -1, Err: errors.New("not a branch")}
	}
	if err := do("", args...); err != nil {
		return err
	}
	return do(dir, "remote", "remove", "origin")
}

// AddRemote adds the remote name, with url, to the repository at dir: layup run
// gives the clone of a work area its origin, the forge's URL with no credential
// (docs/spec/scaffold.md, The work area; task T-eahr). A name or a URL that is
// empty, that starts with a dash, or that holds a space (a name) or a control
// character, starts no git: git could read it as an option.
func AddRemote(dir, name, url string) error {
	args := []string{"remote", "add", "--", name, url}
	if !remoteName.MatchString(name) || url == "" || url[0] == '-' || strings.ContainsFunc(url, func(r rune) bool { return r < 0x20 || r == 0x7f }) {
		return &FailedError{Args: args, Code: -1, Err: errors.New("not a remote name and a URL that git cannot read as an option")}
	}
	return do(dir, args...)
}

// remoteName is the form of the name of a remote that AddRemote accepts.
var remoteName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// FetchSession fetches the commit sha of a session's clone, whose .git
// directory is session, into dst of the clone at dir. No git runs in the
// session's clone: a scratch bare repository whose one alternate is the
// session's object directory takes the ref, and git upload-pack runs there,
// with LAYUP's configuration; the fetch runs with an empty hooks directory.
// Both scratch directories are removed on each return. A sha that names no
// commit is ErrNotACommit (packages.md, The calls of M2b).
func FetchSession(dir, session, sha, dst string) error {
	if !fullObjectID(sha) || !strings.HasPrefix(dst, "refs/") || strings.Contains(dst, ":") {
		return &FailedError{Args: []string{"fetch", dst}, Code: -1, Err: errors.New("not a full object ID, or not a ref under refs/ with no ':'")}
	}
	tmp, err := os.MkdirTemp("", "layup-session-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	empty, err := os.MkdirTemp("", "layup-hooks-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(empty)
	if err := do("", "init", "--bare", "--", tmp); err != nil {
		return err
	}
	objects, err := filepath.Abs(filepath.Join(session, "objects"))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(tmp, "objects", "info"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, "objects", "info", "alternates"), []byte(objects+"\n"), 0o644); err != nil {
		return err
	}
	args := []string{"cat-file", "-t", sha}
	out, err := call(tmp, environ(), args...)
	var failed *FailedError
	switch {
	case errors.As(err, &failed):
		return &FailedError{Args: args, Code: failed.Code, Stderr: failed.Stderr, Err: ErrNotACommit}
	case err != nil:
		return err
	case strings.TrimSpace(string(out)) != "commit":
		return &FailedError{Args: args, Code: 0, Err: ErrNotACommit}
	}
	if err := do(tmp, "update-ref", "refs/heads/session", sha); err != nil {
		return err
	}
	return do(dir, "-c", "core.hooksPath="+empty, "fetch", "--no-tags", "--no-write-fetch-head", "--", tmp, "refs/heads/session:"+dst)
}

// IsAncestor reports whether base is an ancestor of head, both full object
// IDs; exit 1 of git is no, not an error.
func IsAncestor(dir, base, head string) (bool, error) {
	args := []string{"merge-base", "--is-ancestor", base, head}
	if !fullObjectID(base) || !fullObjectID(head) {
		return false, &FailedError{Args: args, Code: -1, Err: errors.New("not a full object ID")}
	}
	_, err := call(dir, environ(), args...)
	var failed *FailedError
	if errors.As(err, &failed) && failed.Code == 1 {
		return false, nil
	}
	return err == nil, err
}

// DiffFile gives the diff of path from base to head with no context line.
func DiffFile(dir, base, head, path string) ([]byte, error) {
	return call(dir, environ(), "diff", "-U0", "--no-color", "--no-renames", "--end-of-options", base, head, "--", path)
}

// DiffBinary gives the diff from base to head that git apply takes, binary
// files included: the payload of a refused diff.
func DiffBinary(dir, base, head string) ([]byte, error) {
	return call(dir, environ(), "diff", "--binary", "--no-renames", "--end-of-options", base, head, "--")
}

// updateRefName is the form of a ref that UpdateRef sets: below refs/heads/layup/, plain characters.
var updateRefName = regexp.MustCompile(`^refs/heads/layup/[A-Za-z0-9._/-]+$`)

// UpdateRef sets the ref to the commit sha of dir: a ref of layup's own namespace, refs/heads/layup/, and a
// full object ID, else it is refused before git starts (docs/spec/packages.md, The calls of M2e; task T-jy6f).
func UpdateRef(dir, ref, sha string) error {
	if !updateRefName.MatchString(ref) || strings.Contains(ref, "..") || strings.Contains(ref, "//") || strings.HasSuffix(ref, "/") || strings.HasSuffix(ref, ".lock") {
		return &FailedError{Args: []string{"update-ref", ref}, Code: -1, Err: errors.New("not a ref of refs/heads/layup/")}
	}
	return doAt(dir, sha, "update-ref", ref, sha)
}

// MergeBase gives the best common ancestor of the commits a and b, both full object IDs; commits with none are an error (task T-89xf,
// packages.md, The calls of M2e).
func MergeBase(dir, a, b string) (string, error) {
	args := []string{"merge-base", a, b}
	if !fullObjectID(a) || !fullObjectID(b) {
		return "", &FailedError{Args: args, Code: -1, Err: errors.New("not a full object ID")}
	}
	out, err := call(dir, environ(), args...)
	if err != nil {
		return "", err
	}
	id := strings.TrimSpace(string(out))
	if !fullObjectID(id) {
		return "", &FailedError{Args: args, Code: 0, Err: fmt.Errorf("the output %q is no object ID", id)}
	}
	return id, nil
}

// Merge merges the commit d into the commit x in a scratch work tree, with a merge commit by who at the time of who, and gives the
// commit: one input gives one commit (NFR-005). clean is false for a merge that conflicts (exit 1 of git with unmerged entries), which
// leaves no commit and no work tree; every other failure, including a failed read of the unmerged entries, is an error. x and d are full
// object IDs and the message is not empty (task T-89xf, packages.md, The calls of M2e).
func Merge(dir, x, d, message string, who Identity) (sha string, clean bool, err error) {
	args := []string{"merge", "--no-ff", "--no-edit", "-m", message, "--end-of-options", d}
	if !fullObjectID(x) || !fullObjectID(d) || message == "" {
		return "", false, &FailedError{Args: args, Code: -1, Err: errors.New("not full object IDs and a message")}
	}
	tmp, err := os.MkdirTemp("", "layup-merge-")
	if err != nil {
		return "", false, err
	}
	defer os.RemoveAll(tmp)
	tree := filepath.Join(tmp, "tree")
	if err := WorktreeAdd(dir, tree, x); err != nil {
		return "", false, err
	}
	defer WorktreeRemove(dir, tree)
	_, err = call(tree, commitEnv(who), args...)
	var failed *FailedError
	if errors.As(err, &failed) && failed.Code == 1 {
		out, lerr := call(tree, environ(), "ls-files", "-u", "-z")
		if lerr != nil {
			return "", false, lerr
		}
		if len(out) == 0 {
			return "", false, err
		}
		return "", false, nil
	} else if err != nil {
		return "", false, err
	}
	sha, err = RevParse(tree, "HEAD")
	return sha, err == nil, err
}
