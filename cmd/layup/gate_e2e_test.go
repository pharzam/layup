//go:build e2e

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// goEnv is the environment of a scenario that runs go: the cache of the host,
// so the standard library is not compiled again, and no download.
func goEnv(t *testing.T) []string {
	t.Helper()
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		t.Fatalf("go env GOCACHE: %v", err)
	}
	return []string{"GOCACHE=" + strings.TrimSpace(string(cache)), "GOTOOLCHAIN=local", "GOPROXY=off", "GOFLAGS="}
}

// fixtureGit runs git in dir as a tool of the test: no configuration of the
// host, a fixed identity and fixed dates, so a commit has one ID.
func fixtureGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	date := "2026-01-01T00:00:00Z"
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=LAYUP test", "GIT_AUTHOR_EMAIL=test@layup.invalid", "GIT_AUTHOR_DATE=" + date,
		"GIT_COMMITTER_NAME=LAYUP test", "GIT_COMMITTER_EMAIL=test@layup.invalid", "GIT_COMMITTER_DATE=" + date}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// commitAll writes the files into dir, commits them, and gives the commit.
func commitAll(t *testing.T, dir string, files map[string]string) string {
	t.Helper()
	for name, text := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, name)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	fixtureGit(t, dir, "add", "--all")
	fixtureGit(t, dir, "commit", "-q", "-m", "fixture")
	return fixtureGit(t, dir, "rev-parse", "HEAD")
}

// goRepo makes a repository of one small Go package with the manifest, on
// the go line of the running toolchain, and gives its path and its commit.
func goRepo(t *testing.T, manifest string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	fixtureGit(t, dir, "init", "-q", "-b", "main")
	base := commitAll(t, dir, map[string]string{
		"go.mod":         "module example.com/fixture\n\ngo " + strings.TrimPrefix(runtime.Version(), "go") + "\n",
		"x.go":           "package fixture\n\n// X is one.\nfunc X() int { return 1 }\n",
		"docs/gates.tsv": manifest,
	})
	return dir, base
}

const goManifest = "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
	"static\tactive\tgofmt\ttest -z \"$(gofmt -l .)\"\t./*.go\t—\n" +
	"vet\tactive\tgo\tgo vet ./...\t./*.go\t—\n" +
	"docs\tactive\tsh\ttrue\tdocs/*.md\t—\n" +
	"layout\tpending\tgo\tgo test ./layout/\tlayout\t—\n" +
	"contract\tpending\tgo\tgo test ./contract/\tcontract\t—\n"

// The demo of #82: on a Go repository, one verdict per kind, with exit code 0
// only when each kind passes or is clear.
func TestGateOnAGoRepository(t *testing.T) {
	dir, base := goRepo(t, goManifest)
	good := commitAll(t, dir, map[string]string{"y.go": "package fixture\n\n// Y is two.\nfunc Y() int { return 2 }\n"})
	r := layupWith(t, goEnv(t), "gate", dir, "--base", base, "--head", good)
	want := "base\thead\tkind\tstate\tresult\treason\n" +
		base + "\t" + good + "\tstatic\tactive\tpass\t—\n" +
		base + "\t" + good + "\tvet\tactive\tpass\t—\n" +
		base + "\t" + good + "\tdocs\tactive\tclear\tno product path\n" +
		base + "\t" + good + "\tlayout\tpending\tclear\tpending: no product path\n" +
		base + "\t" + good + "\tcontract\tpending\tclear\tpending: no product path\n"
	if r.code != 0 || r.stdout != want {
		t.Fatalf("each kind passes or is clear: exit %d, stdout:\n%s\nwant 0 and:\n%s\nstderr:\n%s", r.code, r.stdout, want, r.stderr)
	}
	bad := commitAll(t, dir, map[string]string{"y.go": "package fixture\nfunc  Y()  int { return 2 }\n", "layout/l.go": "package layout\n"})
	r = layupWith(t, goEnv(t), "gate", dir, "--base", base, "--head", bad)
	want = "base\thead\tkind\tstate\tresult\treason\n" +
		base + "\t" + bad + "\tstatic\tactive\tfail\texit 1\n" +
		base + "\t" + bad + "\tvet\tactive\tpass\t—\n" +
		base + "\t" + bad + "\tdocs\tactive\tclear\tno product path\n" +
		base + "\t" + bad + "\tlayout\tpending\tfail\tpending: product path changed: layout/l.go\n" +
		base + "\t" + bad + "\tcontract\tpending\tclear\tpending: no product path\n"
	if r.code != 1 || r.stdout != want {
		t.Fatalf("a kind fails: exit %d, stdout:\n%s\nwant 1 and:\n%s\nstderr:\n%s", r.code, r.stdout, want, r.stderr)
	}
	if !strings.Contains(r.stderr, "layup gate: [1/5] static\n") || !strings.Contains(r.stderr, "layup gate: [5/5] contract\n") {
		t.Fatalf("stderr has no step line of each kind:\n%s", r.stderr)
	}
	again := repeatWith(t, goEnv(t), "gate", dir, "--base", base, "--head", bad)
	if again.stdout != r.stdout || strings.Contains(r.stdout, "layup-gate-") || strings.Contains(r.stdout, dir) {
		t.Fatalf("the repeat rule: a second run gives other bytes, or the table holds a path:\n%s", r.stdout)
	}
}

// NFR-004: a kind whose tool is not found is not-active, exit 1; a base with
// no manifest and a revision that does not resolve give exit 2 and no table,
// never an empty table with exit 0.
func TestGateNeverPassesACheckThatDidNotRun(t *testing.T) {
	dir, base := goRepo(t, "kind\tstate\ttool\tcommand\tscope\tconfig\nimports\tactive\tlayup-no-such-tool-xyz\tlayup-no-such-tool-xyz ./...\t./*.go\t—\n")
	r := layupWith(t, goEnv(t), "gate", dir, "--base", base, "--head", base)
	if r.code != 1 || !strings.HasSuffix(r.stdout, "\timports\tactive\tnot-active\ttool not found: layup-no-such-tool-xyz\n") {
		t.Fatalf("a missing tool: exit %d, stdout %q; want 1 and not-active", r.code, r.stdout)
	}
	empty := t.TempDir()
	fixtureGit(t, empty, "init", "-q", "-b", "main")
	commitAll(t, empty, map[string]string{"x.go": "package x\n"})
	for name, args := range map[string][]string{
		"a base with no manifest":          {"gate", empty, "--base", "HEAD", "--head", "HEAD"},
		"a revision that does not resolve": {"gate", dir, "--base", "no-such-revision", "--head", "HEAD"},
		"a REPO that is not a repository":  {"gate", t.TempDir(), "--base", "HEAD", "--head", "HEAD"},
	} {
		r := layupWith(t, goEnv(t), args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: ") || strings.Contains(r.stderr, "usage:") {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 2, no table, and the reason with no usage", name, r.code, r.stdout, r.stderr)
		}
	}
}
