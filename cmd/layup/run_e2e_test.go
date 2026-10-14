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
	"context"
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
	host   string // the host directory, whose work area holds the ruleset files that the rules are read from (row 58)
	mu     sync.Mutex
	issues int
	posted []string         // the bodies of the comments on the control issue
	intake []string         // the bodies of the comments on the Intake issue, the App's (row 46)
	people []map[string]any // the comments of persons on the Intake issue, scripted by a test (row 48)
	// noRules names the branches whose rules the forge gives as none (row 58).
	noRules map[string]bool
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
		out, _ := gitRun(f.bare, "for-each-ref", "--format=%(refname:short) %(objectname)", "refs/heads")
		list := []map[string]any{}
		for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
			if name, sha, ok := strings.Cut(line, " "); ok {
				list = append(list, map[string]any{"name": name, "commit": map[string]string{"sha": sha}})
			}
		}
		enc.Encode(list)
	case strings.HasPrefix(call, "GET /repos/acme/target/rules/branches/"):
		// The rules of a branch are those of its ruleset file in the work area, as an active ruleset of the forge gives them.
		branch := strings.TrimPrefix(r.URL.Path, "/repos/acme/target/rules/branches/")
		file := "ruleset-default.json"
		if branch == "layup-records" {
			file = "ruleset-records.json"
		}
		var set struct {
			Rules json.RawMessage `json:"rules"`
		}
		data, err := os.ReadFile(filepath.Join(f.host, "work", "acme", "target", "out", file))
		if err != nil || f.noRules[branch] || json.Unmarshal(data, &set) != nil {
			w.Write([]byte("[]"))
			return
		}
		w.Write(set.Rules)
	case call == "GET /repos/acme/target/activity":
		w.Write([]byte("[]"))
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
	case call == "POST /repos/acme/target/issues/2/comments":
		var in struct{ Body string }
		json.NewDecoder(r.Body).Decode(&in)
		f.posted = append(f.posted, in.Body)
		w.WriteHeader(201)
		enc.Encode(map[string]any{"id": 600 + len(f.posted)})
	case call == "POST /repos/acme/target/issues/1/comments":
		var in struct{ Body string }
		json.NewDecoder(r.Body).Decode(&in)
		f.intake = append(f.intake, in.Body)
		w.WriteHeader(201)
		enc.Encode(map[string]any{"id": 800 + len(f.intake)})
	case call == "GET /repos/acme/target/issues/1/comments":
		list := []map[string]any{}
		for i, body := range f.intake {
			list = append(list, map[string]any{"id": 801 + i, "user": map[string]any{"id": 9001, "login": "layup-agent[bot]"},
				"performed_via_github_app": map[string]any{"slug": "layup-agent"},
				"created_at":               "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": body})
		}
		enc.Encode(append(list, f.people...))
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
	forge     *fakeForge
}

func newRunWorld(t *testing.T) runWorld { return newRunWorldWith(t, nil) }

// newRunWorldWith is newRunWorld with a change of the files of the host's directory, by path, which a scenario
// makes to replace one, as the harness, or to add one (row 58).
func newRunWorldWith(t *testing.T, change func(files map[string]string)) runWorld {
	t.Helper()
	w := runWorld{dir: t.TempDir()}
	web := filepath.Join(t.TempDir(), "web")
	w.bare = filepath.Join(web, "acme", "target.git")
	os.MkdirAll(w.bare, 0o755)
	// -b main: the bare repository stands in for GitHub, whose HEAD is the
	// default branch, so the run's clone has a local main (guardrails §2).
	if out, err := gitRun(w.bare, "init", "-q", "--bare", "-b", "main"); err != nil {
		t.Fatal(out)
	}
	// The target receives pushes: the git that receives one starts no
	// maintenance of its own in the background (T-2hm3).
	for _, kv := range [][2]string{{"receive.autogc", "false"}, {"maintenance.auto", "false"}} {
		if out, err := gitRun(w.bare, "config", kv[0], kv[1]); err != nil {
			t.Fatal(out)
		}
	}
	w.forge = &fakeForge{bare: w.bare, host: w.dir}
	srv := httptest.NewServer(http.HandlerFunc(w.forge.serve))
	t.Cleanup(srv.Close)
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	keyFile := filepath.Join(w.dir, "app.pem") // made at run time: no key enters the tree
	files := map[string]string{
		"registers/forge.tsv": "forge\tapp_id\tapp_slug\tkey_file\twatch_slug\tapi\tweb\n" +
			"github\t42\tlayup-agent\t" + keyFile + "\tlayup-watch\t" + srv.URL + "\tfile://" + web + "\n",
		// The harness is the world's own scripts, so no table depends on the
		// host's PATH (condition 3 of the plan review of #168).
		"registers/harnesses.tsv": "harness\tcap\twall\tcommand\tprompt\tversion\tcredential\tcredential_to\trules\tpolicy\tusage\tbilling\tvars\n" +
			"fake\t—\t1\t" + filepath.Join(w.dir, "bin", "harness") + "\tstdin\t" + filepath.Join(w.dir, "bin", "version") + "\t—\t—\tLAYUP-TEST-RULES.md\t—\tnone\tsubscription\t—\n" +
			// spare, the same scripts, the second harness of the
			// completeness review (row 45b).
			"spare\t—\t1\t" + filepath.Join(w.dir, "bin", "harness") + "\tstdin\t" + filepath.Join(w.dir, "bin", "version") + "\t—\t—\tLAYUP-TEST-RULES.md\t—\tnone\tsubscription\t—\n",
		"registers/models.tsv": "harness\tmodel\tcontext\tsource\tdate\tuse\treason\n" +
			"fake\tm1\t1000000\thttps://example.invalid/models\t2026-10-09T12:00:00Z\tyes\t—\n" +
			"spare\tm1\t1000000\thttps://example.invalid/models\t2026-10-09T12:00:00Z\tyes\t—\n",
		"registers/routing.tsv": "role\ttier\tposition\tharness\tmodel\ndeveloper\texecution\t1\tfake\tm1\nsystems-architect\treasoning\t1\tfake\tm1\nproduct-owner\treasoning\t1\tfake\tm1\n" +
			"qa-engineer\treasoning\t1\tfake\tm1\nqa-engineer\treasoning\t2\tspare\tm1\n",
		"bin/version": "#!/bin/sh\necho 1.0.0\n",
		"bin/harness": `#!/bin/sh
p=$(cat)
case "$p" in
"This is a probe"*)
	result=$(printf '%s\n' "$p" | sed -n 's/^Write one file, \(.*\), of tab-separated values, with the header line$/\1/p')
	token=$(printf '%s\n' "$p" | sed -n 's/.*one row "token<TAB>\([0-9a-f]*\)".*/\1/p')
	printf 'kind\tvalue\ntoken\t%s\nfile\tAGENTS.md\n' "$token" > "$result" ;;
*)
	out=$(printf '%s\n' "$p" | sed -n 's/^The output file: //p')
	in=$(printf '%s\n' "$p" | sed -n 's/^The problem statement: //p')
	sp=$(printf '%s\n' "$p" | sed -n 's/^The span table: //p')
	case "$out" in
	*/meaning.tsv) printf 'kind\tquote\tquestion\n' > "$out" ;;
	*/coverage.tsv) printf 'kind\tfact\trequirement\tquestion\n' > "$out" ;;
	*/answers-review.tsv)
		qs=$(printf '%s\n' "$p" | sed -n 's/^The questions: //p')
		LC_ALL=C awk -F '\t' 'BEGIN { printf "question\tstate\treason\n" } NR > 1 && $13 == "answered" { printf "%s\tclosed\tthe answer is clear\n", $1 }' "$qs" > "$out" ;;
	*/draft.tsv) LC_ALL=C awk -F '\t' 'BEGIN { printf "kind\tstatement\tcovers\tcriterion\n" }
		NR > 1 && $4 == "need" { printf "requirement\tThe statement of fact %s.\t%s\tThe criterion of fact %s.\n", $3, $3, $3 }' "$sp" > "$out" ;;
	*) LC_ALL=C awk 'BEGIN { printf "start\tend\tfact\tclass\treason\n"; o = 0; n = 0 }
		{ l = length($0); if (l > 0) { n++; printf "%d\t%d\t%d\tneed\t\342\200\224\n", o, o + l, n }; o += l + 1 }' "$in" > "$out" ;;
	esac ;;
esac
`,
		"psb.md": "# The problem\n",
	}
	if change != nil {
		change(files)
	}
	for p, text := range files {
		os.MkdirAll(filepath.Dir(filepath.Join(w.dir, p)), 0o755)
		mode := os.FileMode(0o644)
		if strings.HasPrefix(p, "bin/") || strings.HasPrefix(p, "ops/") {
			mode = 0o755
		}
		os.WriteFile(filepath.Join(w.dir, p), []byte(text), mode)
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
	return w.runWith(t, binary, nil, args...)
}

// runWith is run with the variables of env added to the fixed environment: the settings of Go for
// the step that runs the real verification (row 58).
func (w runWorld) runWith(t *testing.T, binary string, env []string, args ...string) result {
	t.Helper()
	// A run that waits for a lease that no one releases ends the scenario with its reason, not at the limit of the package.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	cmd.Dir = w.dir
	cmd.Env = append([]string{"HOME=" + t.TempDir(), "TMPDIR=" + t.TempDir(), "PATH=" + os.Getenv("PATH"), "LC_ALL=C"}, env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if ctx.Err() != nil || !errors.As(err, &exit) {
			t.Fatalf("layup %q: %v (no end in three minutes: %v)\n%s", args, err, ctx.Err(), stderr.String())
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
	return w, w.start(t, binary)
}

// start runs layup run --new in the world, with the stand-in Operator.
func (w runWorld) start(t *testing.T, binary string) result {
	t.Helper()
	stop := make(chan struct{})
	go w.operator(stop)
	defer close(stop)
	return w.run(t, binary, startArgs...)
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
	if again.code != 0 || !strings.HasPrefix(again.stdout, "step\tresult\tdetail\nforge\tdone\t") || !strings.Contains(again.stdout, "\nrule-gaps\tdone\t") ||
		!strings.Contains(again.stdout, "\nmeaning\tdone\t0 gaps of meaning, in intake/meaning.tsv\n") ||
		!strings.Contains(again.stdout, "\nnumbering\tdone\t1 spans, 1 facts, in spec/spans.tsv\nspans\tdone\tthe span check passed, in spec/check.tsv\n"+
			"draft\tdone\t1 rows, in spec/draft.tsv\ncompleteness\tdone\t0 items, in intake/coverage.tsv\nmarkers\tdone\tno marker in the pinned tree; intake/markers.tsv has no row\nbatch\tdone\t") ||
		!strings.HasSuffix(again.stdout, " questions, in questions.tsv and the comment of the batch\nanswers\tdone\twaiting: no answer yet in the block idea-owner and the block operator\n") || len(w.forge.intake) != 1 || !strings.HasPrefix(w.forge.intake[0], "LAYUP Intake: the batch of questions\n") {
		t.Fatalf("layup run TARGET: exit %d\n%s\n%s", again.code, again.stdout, again.stderr)
	}
	// The demo of row 38: the restart's step probe probes the scripted
	// harnesses (two from row 45b), writes their rows and posts a comment
	// for each on the control issue.
	if !strings.Contains(again.stdout, "\nprobe\tdone\tprobed and passed 2, probed and failed 0, skipped 0\n") {
		t.Errorf("the step probe of the restart:\n%s", again.stdout)
	}
	if out, err := gitRun(w.bare, "show", "refs/heads/layup-records:harnesses.tsv"); err != nil || !strings.Contains(out, "\tfake\t1.0.0\tm1\tpassed\t—\tAGENTS.md\t") {
		t.Errorf("harnesses.tsv of the records branch: %v\n%s", err, out)
	}
	// The steps meaning (row 43b), numbering (row 44), draft (row 45a) and
	// completeness (row 45b) each run one session, which posts its own
	// comment after the probes'.
	if len(w.forge.posted) != 6 || !strings.HasPrefix(w.forge.posted[0], "S-") || !strings.Contains(w.forge.posted[0], ": probe of fake 1.0.0: passed\n") ||
		!strings.Contains(w.forge.posted[1], ": probe of spare 1.0.0: passed\n") ||
		!strings.HasSuffix(w.forge.posted[2], ": systems-architect session of meaning: done\n") ||
		!strings.HasSuffix(w.forge.posted[3], ": product-owner session of numbering: done\n") ||
		!strings.HasSuffix(w.forge.posted[4], ": product-owner session of draft: done\n") ||
		!strings.HasSuffix(w.forge.posted[5], ": qa-engineer session of completeness: done\n") {
		t.Errorf("the comments on the control issue %q", w.forge.posted)
	}
	// The same bytes on a repeat in a new world, as the first run changes
	// the records (README.md, Determinism).
	w2, _ := startWorld(t, binary)
	if repeat := w2.run(t, binary, "run", "acme/target", "--host", "."); repeat.code != again.code || sameBytes([]byte(repeat.stdout), []byte(again.stdout)) != nil {
		t.Errorf("layup run TARGET in two worlds: exit %d and %d\n%s\n%s", again.code, repeat.code, again.stdout, repeat.stdout)
	}
	// A second restart in the same world probes none: the last probe at the
	// version passed.
	third := w.run(t, binary, "run", "acme/target", "--host", ".")
	if third.code != 0 || !strings.Contains(third.stdout, "\nprobe\tdone\tprobed and passed 0, probed and failed 0, skipped 0\n") || len(w.forge.posted) != 6 ||
		!strings.Contains(third.stdout, "\nrule-gaps\tdone\tdone before\n") || !strings.Contains(third.stdout, "\nmeaning\tdone\tdone before\n") ||
		!strings.HasSuffix(third.stdout, "\nnumbering\tdone\tdone before\nspans\tdone\tdone before\ndraft\tdone\tdone before\ncompleteness\tdone\tdone before\n"+
			"markers\tdone\tdone before\nbatch\tdone\tdone before\nanswers\tdone\twaiting: no answer yet in the block idea-owner and the block operator\n") || len(w.forge.intake) != 1 {
		t.Errorf("layup run TARGET again: exit %d, %d comments\n%s", third.code, len(w.forge.posted), third.stdout)
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
