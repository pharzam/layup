//go:build e2e

package main

// The e2e scenarios of layup run (docs/spec/run.md, the acceptance row
// "layup run --new, then layup run TARGET"): the built binary against a fake
// GitHub on loopback (the forge register's api, plain http, so TLS is not
// exercised) and local bare repositories (its web), with no secret and no
// network. Step 4 waits for the Operator's push with the real clock, so each
// Start takes about ten seconds.

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	runOnce    sync.Once
	runDir     string // the directory of the binary and the baseline; TestMain removes it
	runBinary  string // the binary whose pin is the local baseline
	runBase    string // the bare baseline, shared by each world
	runCommit  string // its commit
	runOnceErr error
)

// gitRun runs git as a tool of the test, with no configuration of the host.
func gitRun(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + dir, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t",
		"GIT_AUTHOR_DATE=@1767225600 +0000", "GIT_COMMITTER_DATE=@1767225600 +0000"}
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// buildRunBinary makes the bare baseline once (a commit with fixed dates, so
// its ID is the same in each run of the tests) and builds the binary with
// LAYUP's pin set to it (decision 2 of #131; the release binary keeps the
// values of docs/setup/armature.pin).
func buildRunBinary(t *testing.T) (string, string) {
	t.Helper()
	runOnce.Do(func() {
		dir, err := os.MkdirTemp("", "layup-run-e2e-")
		if err != nil {
			runOnceErr = err
			return
		}
		runDir = dir
		src := filepath.Join(dir, "src")
		os.MkdirAll(src, 0o755)
		os.WriteFile(filepath.Join(src, "README.md"), []byte("# Armature\n"), 0o644)
		for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"commit", "-q", "-m", "baseline"}} {
			if out, err := gitRun(src, args...); err != nil {
				runOnceErr = errors.New(out)
				return
			}
		}
		commit, _ := gitRun(src, "rev-parse", "HEAD")
		runCommit = strings.TrimSpace(commit)
		runBase = filepath.Join(dir, "armature.git")
		if out, err := gitRun(dir, "clone", "-q", "--bare", src, runBase); err != nil {
			runOnceErr = errors.New(out)
			return
		}
		runBinary = filepath.Join(dir, "layup")
		ld := "-X github.com/pharzam/layup/internal/cli.pinSource=file://" + runBase +
			" -X github.com/pharzam/layup/internal/cli.pinCommit=" + strings.TrimSpace(commit)
		build := exec.Command("go", "build", "-ldflags", ld, "-o", runBinary, ".")
		build.Env = append(os.Environ(), "GOFLAGS=", "GOENV=off", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
		if out, err := build.CombinedOutput(); err != nil {
			runOnceErr = errors.New(string(out))
		}
	})
	if runOnceErr != nil {
		t.Fatal(runOnceErr)
	}
	return runBinary, runBase
}

// fakeForge plays the calls of M2a; its branches are those of the bare target.
type fakeForge struct {
	bare   string
	mu     sync.Mutex
	issues int
}

func (f *fakeForge) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	call := r.Method + " " + r.URL.Path
	enc := json.NewEncoder(w)
	switch {
	case call == "GET /repos/acme/target/installation":
		enc.Encode(map[string]any{"id": 7, "permissions": map[string]string{"contents": "write", "issues": "write", "metadata": "read"}})
	case call == "POST /app/installations/7/access_tokens":
		w.WriteHeader(201)
		enc.Encode(map[string]any{"token": "ghs_e2e", "expires_at": "2099-01-01T00:00:00Z"})
	case call == "GET /repos/acme/target":
		enc.Encode(map[string]any{"default_branch": "main", "visibility": "public"})
	case call == "GET /repos/acme/target/branches":
		out, _ := gitRun(f.bare, "for-each-ref", "--format=%(refname:short)", "refs/heads")
		list := []map[string]string{}
		for _, n := range strings.Fields(out) {
			list = append(list, map[string]string{"name": n})
		}
		enc.Encode(list)
	case strings.HasPrefix(call, "GET /users/"):
		ids := map[string]int{"pharzam": 101, "carol": 303, "layup-agent[bot]": 9001}
		if id, ok := ids[strings.TrimPrefix(r.URL.Path, "/users/")]; ok {
			enc.Encode(map[string]any{"id": id})
			return
		}
		http.Error(w, `{"message":"Not Found"}`, 404)
	case call == "POST /repos/acme/target/issues":
		f.issues++
		w.WriteHeader(201)
		enc.Encode(map[string]any{"number": f.issues})
	case call == "GET /repos/acme/target/issues/2/comments":
		enc.Encode([]map[string]any{{"id": 501, "user": map[string]any{"id": 777, "login": "layup-watch[bot]"},
			"performed_via_github_app": map[string]any{"slug": "layup-watch"},
			"created_at":               "2026-10-08T12:00:00Z", "updated_at": "2026-10-08T12:00:00Z", "body": "watching\n"}})
	default:
		http.Error(w, `{"message":"Not Found"}`, 404)
	}
}

// runWorld is one host: its directory with the registers, the key and the
// briefs, a bare target under its web, and a fake forge.
type runWorld struct {
	dir, bare string
}

func newRunWorld(t *testing.T) runWorld {
	t.Helper()
	w := runWorld{dir: t.TempDir()}
	web := filepath.Join(t.TempDir(), "web")
	w.bare = filepath.Join(web, "acme", "target.git")
	os.MkdirAll(w.bare, 0o755)
	if out, err := gitRun(w.bare, "init", "-q", "--bare"); err != nil {
		t.Fatal(out)
	}
	srv := httptest.NewServer(http.HandlerFunc((&fakeForge{bare: w.bare}).serve))
	t.Cleanup(srv.Close)
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(w.dir, "app.pem") // made at run time: no key enters the tree
	files := map[string]string{
		"registers/forge.tsv": "forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb\n" +
			"github\t42\tlayup-agent\t" + keyFile + "\tlayup-watch\t" + srv.URL + "\tfile://" + web + "\n",
		"registers/harnesses.tsv": "harness\tcap\twall\tcommand\tprompt\tversion\tcredential\tcredential_to\trules\tpolicy\tusage\tbilling\tvars\n" +
			"claude\t10.0\t60\tclaude -p --model {model} --max-budget-usd {cap} {prompt}\targ\tclaude --version\t—\t—\tCLAUDE.md\t—\tclaude-result\tapi\t—\n",
		"registers/models.tsv": "harness\tmodel\tcontext\tsource\tdate\tuse\treason\n" +
			"claude\tclaude-opus-5-5\t1000000\thttps://docs.claude.com/models\t2026-10-09T12:00:00Z\tyes\t—\n",
		"registers/routing.tsv": "role\ttier\tposition\tharness\tmodel\ndeveloper\texecution\t1\tclaude\tclaude-opus-5-5\n",
		"psb.md":                "# The problem\n",
	}
	for p, text := range files {
		os.MkdirAll(filepath.Dir(filepath.Join(w.dir, p)), 0o755)
		os.WriteFile(filepath.Join(w.dir, p), []byte(text), 0o644)
	}
	os.WriteFile(keyFile, pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}), 0o600)
	return w
}

// operator plays the Operator's push of the root commit, as the printed
// command does, once the root commit exists; stop ends it.
func (w runWorld) operator(stop <-chan struct{}) {
	root := filepath.Join(w.dir, "roots", "acme", "target")
	for {
		select {
		case <-stop:
			return
		default:
		}
		if id, err := gitRun(root, "rev-parse", "--verify", "HEAD"); err == nil {
			if _, err := gitRun(root, "push", "-q", "file://"+w.bare, strings.TrimSpace(id)+":refs/heads/main"); err == nil {
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// run runs the binary in the world's directory, so the arguments are the same
// in each world.
func (w runWorld) run(t *testing.T, binary string, args ...string) result {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = w.dir
	cmd.Env = []string{"HOME=" + t.TempDir(), "TMPDIR=" + t.TempDir(), "PATH=" + os.Getenv("PATH"), "LC_ALL=C"}
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("layup %q: %v", args, err)
		}
		code = exit.ExitCode()
	}
	return result{stdout.String(), stderr.String(), code}
}

var startArgs = []string{"run", "--new", "acme/target", "--host", ".", "--psb", "psb.md", "--operator", "pharzam",
	"--idea-owner", "carol", "--plan", "free", "--intake-cap", "50.0,8.0", "--lease-h", "5", "--watch-t", "1"}

// startWorld runs layup run --new in a new world, with the stand-in Operator.
func startWorld(t *testing.T, binary string) (runWorld, result) {
	t.Helper()
	w := newRunWorld(t)
	stop := make(chan struct{})
	go w.operator(stop)
	defer close(stop)
	return w, w.run(t, binary, startArgs...)
}

// The demo of row 26: layup run --new, then layup run TARGET, give their
// tables, the same bytes on a repeat. Start's first run changes its forge (a
// second --new on the same target is forge: fail, O-163), so its repeat is a
// new world with the same arguments; the restart repeats in its world.
func TestRunNewThenRestart(t *testing.T) {
	binary, _ := buildRunBinary(t)
	w, first := startWorld(t, binary)
	if first.code != 0 {
		t.Fatalf("layup run --new: exit %d\nstdout:\n%s\nstderr:\n%s", first.code, first.stdout, first.stderr)
	}
	rows := strings.Split(strings.TrimSuffix(first.stdout, "\n"), "\n")
	if len(rows) != 10 || rows[0] != "step\tresult\tdetail" {
		t.Fatalf("the table of Start:\n%s", first.stdout)
	}
	for i, name := range []string{"forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease"} {
		if !strings.HasPrefix(rows[i+1], name+"\tdone\t") {
			t.Errorf("row %d of Start is %q; want %s done", i+1, rows[i+1], name)
		}
	}
	for _, line := range []string{"layup run: [1/9] forge\n", "layup run: [9/9] lease\n", "push -- file://"} {
		if !strings.Contains(first.stderr, line) {
			t.Errorf("standard error has no %q:\n%s", line, first.stderr)
		}
	}
	if strings.Contains(first.stdout, w.dir) || strings.Contains(first.stdout+first.stderr, "ghs_") {
		t.Errorf("the table holds a path of the host, or an output holds a token:\n%s", first.stdout)
	}
	// The pin of the binary is the local baseline (-ldflags -X): a name that
	// the linker did not find would leave the real baseline of armature.pin.
	if out, err := gitRun(w.bare, "show", "refs/heads/layup-records:start/start.tsv"); err != nil || !strings.Contains(out, "watch\tconfirmed\trun\n") ||
		!strings.Contains(out, "pin.source\tfile://"+runBase+"\trun\n") || !strings.Contains(out, "pin.commit\t"+runCommit+"\trun\n") {
		t.Errorf("start.tsv of the records branch: %v\n%s", err, out)
	}
	_, second := startWorld(t, binary)
	if err := sameBytes([]byte(first.stdout), []byte(second.stdout)); err != nil || first.code != second.code {
		t.Errorf("layup run --new in two worlds: exit %d and %d; standard output: %v", first.code, second.code, err)
	}

	again := w.run(t, binary, "run", "acme/target", "--host", ".")
	if again.code != 0 || !strings.HasPrefix(again.stdout, "step\tresult\tdetail\nforge\tdone\t") || !strings.Contains(again.stdout, "\nphase\tdone\t") {
		t.Fatalf("layup run TARGET: exit %d\n%s\n%s", again.code, again.stdout, again.stderr)
	}
	if third := w.run(t, binary, "run", "acme/target", "--host", "."); third.code != again.code || sameBytes([]byte(third.stdout), []byte(again.stdout)) != nil {
		t.Errorf("layup run TARGET twice: exit %d and %d\n%s\n%s", again.code, third.code, again.stdout, third.stdout)
	}
}

func TestRunUsageAndInputErrors(t *testing.T) {
	binary, _ := buildRunBinary(t)
	w := newRunWorld(t)
	if r := w.run(t, binary, "run", "--new", "acme/target", "--host", "."); r.code != 2 || r.stdout != "" || !strings.Contains(r.stderr, "missing flag --psb") || !strings.Contains(r.stderr, "usage:") {
		t.Errorf("a missing flag: exit %d, stdout %q, stderr %q", r.code, r.stdout, r.stderr)
	}
	os.WriteFile(filepath.Join(w.dir, "registers", "forge.tsv"), []byte("forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb\n"), 0o644)
	if r := w.run(t, binary, startArgs...); r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: ") || strings.Contains(r.stderr, "usage:") {
		t.Errorf("a forge register with no row: exit %d, stdout %q, stderr %q", r.code, r.stdout, r.stderr)
	}
}
