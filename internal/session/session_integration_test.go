//go:build integration

package session

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func run(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@layup.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@layup.invalid"}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func entries(t *testing.T, dir string) []string {
	t.Helper()
	var names []string
	filepath.WalkDir(dir, func(p string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(dir, p)
			names = append(names, rel)
		}
		return nil
	})
	return names
}

// runClone makes the run's clone: main with two files, a records branch and a
// tag; it gives the clone and the base commit.
func runClone(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	run(t, dir, "init", "-q", "-b", "main")
	os.WriteFile(filepath.Join(dir, "AGENTS.md"), []byte("# rules\n"), 0o644)
	os.MkdirAll(filepath.Join(dir, "internal"), 0o755)
	os.WriteFile(filepath.Join(dir, "internal", "x.go"), []byte("package x\n"), 0o644)
	run(t, dir, "add", "-A")
	run(t, dir, "commit", "-q", "-m", "chore: base")
	run(t, dir, "branch", "layup-records")
	run(t, dir, "tag", "v1")
	base := run(t, dir, "rev-parse", "HEAD")
	// A commit after the base, with a file the base lacks, so "at the base"
	// can fail on its own rule (round 1 of #159, note 2).
	os.WriteFile(filepath.Join(dir, "later.txt"), []byte("later\n"), 0o644)
	run(t, dir, "add", "-A")
	run(t, dir, "commit", "-q", "-m", "chore: after the base")
	return dir, base
}

func TestMakeASessionDirectory(t *testing.T) {
	host := t.TempDir()
	src, base := runClone(t)
	key := filepath.Join(t.TempDir(), "devin.toml")
	os.WriteFile(key, []byte("token = \"x\"\n"), 0o600)
	h := Harness{Credential: key, CredentialTo: "file:.config/devin/credentials.toml"}
	d, err := Make(host, "acme/target", "S-1a2b3c4d", src, "main", "T-ab12", 1, base, []byte("Do the task.\n"), h)
	if err != nil {
		t.Fatal(err)
	}
	if d.Root != filepath.Join(host, "sessions", "S-1a2b3c4d") || !filepath.IsAbs(d.Home) || !filepath.IsAbs(d.Tmp) {
		t.Fatalf("the directory: %+v", d)
	}
	if refs := run(t, d.Repo, "for-each-ref", "--format=%(refname)"); refs != "refs/heads/main\nrefs/heads/task/T-ab12/1" {
		t.Errorf("the refs of repo/: %q; want main and the task branch only, no records branch and no tag", refs)
	}
	if b := run(t, d.Repo, "symbolic-ref", "HEAD"); b != "refs/heads/task/T-ab12/1" || run(t, d.Repo, "rev-parse", "HEAD") != base {
		t.Errorf("repo/ is on %q, not task/T-ab12/1 at the base", b)
	}
	if remotes := run(t, d.Repo, "remote"); remotes != "" {
		t.Errorf("repo/ has a remote: %q", remotes)
	}
	if st := run(t, d.Repo, "status", "--porcelain"); st != "" || !slices.Contains(entries(t, d.Repo), "internal/x.go") ||
		slices.Contains(entries(t, d.Repo), "later.txt") {
		t.Errorf("the work tree and the index of repo/ are not the base: %q", st)
	}
	if got := entries(t, d.Home); !slices.Equal(got, []string{".config/devin/credentials.toml", ".gitconfig"}) {
		t.Errorf("home/ holds %q; want .gitconfig and the credential of file: only", got)
	}
	if info, err := os.Stat(filepath.Join(d.Home, ".config/devin/credentials.toml")); err != nil || info.Mode().Perm() != 0o600 {
		t.Errorf("the copied credential: %v, %v; want mode 0600", info, err)
	}
	cfg, _ := os.ReadFile(filepath.Join(d.Home, ".gitconfig"))
	if string(cfg) != "[user]\n\tname = layup session S-1a2b3c4d\n\temail = S-1a2b3c4d@sessions.layup.invalid\n" {
		t.Errorf("home/.gitconfig: %q", cfg)
	}
	if got := entries(t, d.Tmp); len(got) != 0 {
		t.Errorf("tmp/ holds %q; want nothing", got)
	}
	if got := entries(t, d.Result); len(got) != 0 {
		t.Errorf("result/ holds %q; want nothing", got)
	}
	if p, _ := os.ReadFile(d.Prompt); string(p) != "Do the task.\n" {
		t.Errorf("prompt.md: %q", p)
	}
	if _, err := Make(host, "acme/target", "S-1a2b3c4d", src, "main", "T-ab12", 1, base, nil, Harness{}); !errors.Is(err, fs.ErrExist) {
		t.Errorf("a second Make of the same ID: %v; want the directory refused as one that exists (a session ID is never reused)", err)
	}
}

// Sweep removes the directories that a stopped run of the target left, and
// keeps those of another target, whose own lease guards them.
func TestSweep(t *testing.T) {
	host := t.TempDir()
	src, base := runClone(t)
	for _, c := range []struct{ target, id string }{{"acme/target", "S-00000001"}, {"acme/other", "S-00000002"}} {
		if _, err := Make(host, c.target, c.id, src, "main", "T-ab12", 1, base, nil, Harness{}); err != nil {
			t.Fatal(err)
		}
	}
	os.MkdirAll(filepath.Join(host, "sessions", "S-00000003", "repo"), 0o700)        // a Make that stopped before its first write
	os.WriteFile(filepath.Join(host, "sessions", "notes.txt"), []byte("x\n"), 0o600) // not a session directory: skipped
	if err := Sweep(host, "acme/target"); err != nil {
		t.Fatal(err)
	}
	left, _ := os.ReadDir(filepath.Join(host, "sessions"))
	var names []string
	for _, e := range left {
		names = append(names, e.Name())
	}
	if !slices.Equal(names, []string{"S-00000002", "notes.txt"}) {
		t.Errorf("after Sweep: %q; want the session of the other target only", names)
	}
	if err := Sweep(t.TempDir(), "acme/target"); err != nil {
		t.Errorf("a host with no sessions: %v", err)
	}
}

func TestEnvironReadsTheCredentialFile(t *testing.T) {
	key := filepath.Join(t.TempDir(), "claude.key")
	os.WriteFile(key, []byte("sk-test\n"), 0o600)
	env, err := Environ(Dir{Home: "/h/home", Tmp: "/h/tmp"}, Harness{Credential: key, CredentialTo: "var:ANTHROPIC_API_KEY"})
	if err != nil || env[len(env)-1] != "ANTHROPIC_API_KEY=sk-test" {
		t.Errorf("Environ: %q, %v", env, err)
	}
}
