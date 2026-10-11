//go:build integration

package run

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/forge/github"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/psb"
	"github.com/pharzam/layup/internal/tsv"
)

// fakeGitHub plays the calls of M2a for the steps of Start and the restart
// (docs/spec/run.md, the acceptance row Start): its branches are those of the
// bare target repository, so the stand-in Operator's push shows.
type fakeGitHub struct {
	t        *testing.T
	bare     string
	mu       sync.Mutex
	issues   int
	failOpen int // the number of the issue whose opening fails, or 0
	comments map[int][]map[string]any
	posted   []string // the bodies of the comments that the run posted, on any issue (row 43b)
}

func (f *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	call := r.Method + " " + r.URL.Path
	switch {
	case call == "GET /repos/acme/target/installation":
		json.NewEncoder(w).Encode(map[string]any{"id": 7, "permissions": map[string]string{"contents": "write", "issues": "write", "metadata": "read"}})
	case call == "POST /app/installations/7/access_tokens":
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"token": "ghs_runtest", "expires_at": "2099-01-01T00:00:00Z"})
	case call == "GET /repos/acme/target":
		json.NewEncoder(w).Encode(map[string]any{"default_branch": "main", "visibility": "public"})
	case call == "GET /repos/acme/target/branches":
		out, err := exec.Command("git", "-C", f.bare, "for-each-ref", "--format=%(refname:short)", "refs/heads").Output()
		if err != nil {
			http.Error(w, `{"message":"no bare repository"}`, 500)
			return
		}
		list := []map[string]string{}
		for _, name := range strings.Fields(string(out)) {
			list = append(list, map[string]string{"name": name})
		}
		json.NewEncoder(w).Encode(list)
	case strings.HasPrefix(call, "GET /users/"):
		ids := map[string]int{"pharzam": 101, "carol": 303, "layup-agent[bot]": 9001}
		login := strings.TrimPrefix(r.URL.Path, "/users/")
		if id, ok := ids[login]; ok {
			json.NewEncoder(w).Encode(map[string]any{"id": id, "login": login})
			return
		}
		http.Error(w, `{"message":"Not Found"}`, 404)
	case call == "POST /repos/acme/target/issues":
		if f.issues+1 == f.failOpen {
			f.failOpen = 0
			http.Error(w, `{"message":"Server Error"}`, 500)
			return
		}
		f.issues++
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"number": f.issues})
	case strings.HasPrefix(call, "POST /repos/acme/target/issues/") && strings.HasSuffix(call, "/comments"):
		var in struct{ Body string }
		json.NewDecoder(r.Body).Decode(&in)
		f.posted = append(f.posted, in.Body)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": 600 + len(f.posted)})
	case strings.HasPrefix(call, "GET /repos/acme/target/issues/") && strings.HasSuffix(call, "/comments"):
		n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/repos/acme/target/issues/"), "/comments"))
		list := f.comments[n]
		if list == nil {
			list = []map[string]any{}
		}
		json.NewEncoder(w).Encode(list)
	default:
		http.Error(w, `{"message":"Not Found"}`, 404)
	}
}

// fastClock is the clock of the run and the adapter: its time goes 600 times
// faster than the real one, so a wait of ten seconds takes about 17 ms.
type fastClock struct{ begin, real time.Time }

func (c fastClock) Now() time.Time { return c.begin.Add(time.Since(c.real) * 600) }
func (c fastClock) Sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d / 600):
		return nil
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t",
		"GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}

// world is one test target: the host directory, the bare baseline, the bare
// target under web, and the fake forge.
type world struct {
	t        *testing.T
	dir      string
	web      string
	bare     string
	baseline string
	pin      string // the baseline's commit
	fake     *fakeGitHub
	srv      *httptest.Server
	key      *rsa.PrivateKey
	lines    []string
	mu       sync.Mutex
	hook     func(line string) // a test acts on a progress line
}

func newWorld(t *testing.T) *world {
	t.Setenv("HOME", t.TempDir())
	w := &world{t: t, dir: t.TempDir(), web: t.TempDir()}
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "README.md"), []byte("# Armature\n"), 0o644)
	gitOut(t, src, "init", "-q", "-b", "main")
	gitOut(t, src, "add", "-A")
	gitOut(t, src, "commit", "-q", "-m", "baseline")
	w.pin = strings.TrimSpace(gitOut(t, src, "rev-parse", "HEAD"))
	w.baseline = filepath.Join(t.TempDir(), "armature.git")
	gitOut(t, "/", "clone", "-q", "--bare", src, w.baseline)
	w.bare = filepath.Join(w.web, "acme", "target.git")
	// HEAD names main, as the forge's HEAD names its default branch, so a
	// clone of the target has its branch main (a session clones it, row 43b).
	gitOut(t, "/", "init", "-q", "--bare", "-b", "main", w.bare)
	w.fake = &fakeGitHub{t: t, bare: w.bare, comments: map[int][]map[string]any{
		2: {{"id": 501, "user": map[string]any{"id": 777, "login": "layup-watch[bot]"}, "performed_via_github_app": map[string]any{"slug": "layup-watch"},
			"created_at": "2026-10-08T12:00:00Z", "updated_at": "2026-10-08T12:00:00Z", "body": "the dead-man job is watching\n"}},
	}}
	w.srv = httptest.NewServer(http.HandlerFunc(w.fake.serve))
	t.Cleanup(w.srv.Close)
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	w.key = k
	// The world's harness, claude, is its own scripts (row 43b): it answers a
	// probe as probeHarness does, and the meaning prompt with bin/meaning.tsv,
	// which a test may replace (no gap by default).
	bin := filepath.Join(w.dir, "bin")
	os.MkdirAll(bin, 0o755)
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "meaning.tsv"), []byte("kind\tquote\tquestion\n"), 0o644)
	os.WriteFile(filepath.Join(bin, "harness"), []byte(`#!/bin/sh
p=$(cat)
case "$p" in
"This is a probe"*)
	result=$(printf '%s\n' "$p" | sed -n 's/^Write one file, \(.*\), of tab-separated values, with the header line$/\1/p')
	token=$(printf '%s\n' "$p" | sed -n 's/.*one row "token<TAB>\([0-9a-f]*\)".*/\1/p')
	printf 'kind\tvalue\ntoken\t%s\nfile\tAGENTS.md\n' "$token" > "$result" ;;
*)
	out=$(printf '%s\n' "$p" | sed -n 's/^The output file: //p')
	cat "$(dirname "$0")/meaning.tsv" > "$out" ;;
esac
`), 0o755)
	return w
}

func (w *world) config(version string) Config {
	clock := fastClock{begin: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC), real: time.Now()}
	reg := forge.Register{Forge: "github", AppID: 42, AppSlug: "layup-agent", KeyFile: "/k.pem", WatchSlug: "layup-watch", API: w.srv.URL, Web: "file://" + w.web}
	progress := func(l string) {
		w.mu.Lock()
		w.lines = append(w.lines, l)
		hook := w.hook
		w.mu.Unlock()
		if hook != nil {
			hook(l)
		}
	}
	adapter := github.New(github.Config{API: reg.API, AppID: reg.AppID, Key: w.key, Owner: "acme", Name: "target",
		Client: w.srv.Client(), Progress: progress, Now: clock.Now, Sleep: clock.Sleep})
	return Config{Owner: "acme", Name: "target", Dir: w.dir, PSB: []byte("# The problem\n\nbyte for byte\n"), Vision: []byte("# The vision\n"),
		Operator: "pharzam", IdeaOwner: "carol", Plan: "free", IntakeCap: "50.0,8.0", LeaseH: 5, WatchT: 1,
		Register: reg, Harnesses: [][]string{{"claude", "10.0", "60", filepath.Join(w.dir, "bin", "harness"), "stdin", filepath.Join(w.dir, "bin", "version"),
			"", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""}, {"devin", "", "30"}},
		Models:  [][]string{{"claude", "m1", "1000000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""}},
		Routing: [][]string{{"systems-architect", "reasoning", "1", "claude", "m1"}},
		Forge:   adapter, Clock: clock,
		Progress: progress, RunID: "0123456789abcdef", HostName: "host-1", Version: version,
		Baseline: "file://" + w.baseline, LayupPin: "a95965534b14b0bf14ad74da0c9a45b5f4aedf88"}
}

// operator plays the Operator's push of the root commit, as the printed
// command does, with no token, once the root commit exists.
func (w *world) operator(ctx context.Context) {
	root := filepath.Join(w.dir, "roots", "acme", "target")
	for ctx.Err() == nil {
		if id, err := git.RevParse(root, "HEAD"); err == nil {
			if git.Push(root, "file://"+w.bare, id, "main", git.Auth{}) == nil {
				return
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func (w *world) start(cfg Config) []Step {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	go w.operator(ctx)
	return Start(ctx, cfg)
}

func show(t *testing.T, bare, rev, path string) string {
	t.Helper()
	return gitOut(t, bare, "show", rev+":"+path)
}

func steps(rows []Step) string {
	var b strings.Builder
	WriteSteps(&b, rows)
	return b.String()
}

func allDone(t *testing.T, name string, rows []Step, want ...string) {
	t.Helper()
	var got []string
	for _, r := range rows {
		got = append(got, r.Name)
		if r.Result != "done" {
			t.Errorf("%s: the step %s is %s: %s", name, r.Name, r.Result, r.Detail)
		}
	}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("%s: the steps %v, want %v\n%s", name, got, want, steps(rows))
	}
}

// The demo of row 25b: against the httptest forge and a local bare
// repository, Start makes the first records commit with the pin, the briefs,
// approvers.tsv, start.tsv and the lease row.
func TestStartMakesTheFirstRecordsCommit(t *testing.T) {
	w := newWorld(t)
	rows := w.start(w.config("0.1.0-dev"))
	allDone(t, "Start", rows, "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")

	first := strings.Fields(gitOut(t, w.bare, "rev-list", "--reverse", "refs/heads/layup-records"))[0]
	if parents := strings.Fields(gitOut(t, w.bare, "rev-list", "--parents", "-n", "1", first)); len(parents) != 1 {
		t.Errorf("the first records commit has parents %v; want an orphan", parents[1:])
	}
	if who := gitOut(t, w.bare, "log", "-1", "--format=%an <%ae>|%cn <%ce>", first); who != "layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>|layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>\n" {
		t.Errorf("the author and committer %q; want the App's bot", who)
	}
	start := show(t, w.bare, first, "start/start.tsv")
	for _, line := range []string{"layup.version\t0.1.0-dev\trun", "pin.source\tfile://" + w.baseline + "\trun", "pin.commit\t" + w.pin + "\trun",
		"operator.id\t101\tforge", "idea-owner.id\t303\tforge", "app.permissions\tcontents:write issues:write metadata:read\tforge",
		"harness.devin.cap\t—\tregister", "issue.intake\t—\tforge", "watch\t—\trun"} {
		if !strings.Contains(start, line+"\n") {
			t.Errorf("the first start.tsv has no line %q:\n%s", line, start)
		}
	}
	if got := show(t, w.bare, first, "start/problem-statement.md"); got != "# The problem\n\nbyte for byte\n" {
		t.Errorf("the brief %q", got)
	}
	if got := show(t, w.bare, first, "start/vision.md"); got != "# The vision\n" {
		t.Errorf("the vision %q", got)
	}
	if got := show(t, w.bare, first, "approvers.tsv"); !strings.Contains(got, "101\toperator\tpharzam\t") || !strings.Contains(got, "303\tidea-owner\tcarol\t") {
		t.Errorf("approvers.tsv:\n%s", got)
	}
	if got := show(t, w.bare, first, "lease.tsv"); !strings.Contains(got, "0123456789abcdef\thost-1\t0.1.0-dev\t") || !strings.HasSuffix(got, "\theld\n") {
		t.Errorf("lease.tsv of the first commit:\n%s", got)
	}
	if got := show(t, w.bare, first, "README.md"); got != startReadme {
		t.Errorf("the README:\n%s", got)
	}
	// The last commit: the issues, the watch, the copy of the notice, the release.
	end := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv")
	for _, line := range []string{"issue.intake\t1\tforge", "issue.control\t2\tforge", "watch\tconfirmed\trun"} {
		if !strings.Contains(end, line+"\n") {
			t.Errorf("the last start.tsv has no line %q", line)
		}
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "lease.tsv"); !strings.HasSuffix(got, "\treleased\n") {
		t.Errorf("lease.tsv at the end:\n%s", got)
	}
	if last := gitOut(t, w.bare, "log", "-1", "--format=%s", "refs/heads/layup-records"); last != "lease: released\n" {
		t.Errorf("the last records commit is %q; want the release, after the heartbeat ended", last)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "copies/501-1.md"); got != "the dead-man job is watching\n" {
		t.Errorf("the copy of the notice %q", got)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "copies.tsv"); !strings.Contains(got, "501\t1\t2\t777\tlayup-watch[bot]\tlayup-watch\t") {
		t.Errorf("copies.tsv:\n%s", got)
	}
	for _, r := range rows {
		if strings.Contains(r.Detail, w.dir) || strings.Contains(r.Detail, "ghs_") {
			t.Errorf("the detail of %s holds a path of the host or a token: %q", r.Name, r.Detail)
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	pushed := false
	for _, l := range w.lines {
		pushed = pushed || strings.Contains(l, "push -- file://"+w.web+"/acme/target.git main")
		if strings.Contains(l, "ghs_") {
			t.Errorf("a progress line holds a token: %q", l)
		}
	}
	if !pushed {
		t.Errorf("no progress line prints the push of the root commit: %q", w.lines)
	}
}

func TestARestartWithAnotherVersionFailsAtVersion(t *testing.T) {
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.2.0"))
	if got := steps(rows); !strings.HasPrefix(got, "step\tresult\tdetail\nforge\tdone\t") || !strings.Contains(got, "\nclone\tdone\t") || !strings.Contains(got, "\nversion\tfail\t") {
		t.Errorf("the restart with another version:\n%s", got)
	}
	// The same version: the lease is released, so it is taken at once; each
	// step of Start is done.
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning")
}

// A run that stopped at opening: the restart takes the lease over from the
// stopped run and runs step 7 again for that issue.
func TestARestartAfterARunThatStoppedAtOpening(t *testing.T) {
	w := newWorld(t)
	w.fake.failOpen = 2
	rows := w.start(w.config("0.1.0-dev"))
	if last := rows[len(rows)-1]; last.Name != "issues" || last.Result != "fail" {
		t.Fatalf("the Start with a failed opening:\n%s", steps(rows))
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv"); !strings.Contains(got, "issue.control\topening\tforge\n") {
		t.Fatalf("the control issue is not opening:\n%s", got)
	}
	cfg := w.config("0.1.0-dev")
	cfg.RunID = "fedcba9876543210"
	cfg.WatchT, cfg.LeaseH = 0, 0 // a restart has no --watch-t and no --lease-h: start.tsv gives them
	notice := w.fake.comments[2]
	w.fake.mu.Lock()
	w.fake.comments = map[int][]map[string]any{} // the notice comes after the first read
	w.fake.mu.Unlock()
	var once sync.Once
	w.mu.Lock()
	w.hook = func(l string) {
		if strings.HasPrefix(l, "watch: waiting") {
			once.Do(func() { w.fake.mu.Lock(); w.fake.comments[2] = notice; w.fake.mu.Unlock() })
		}
	}
	w.mu.Unlock()
	rows = Restart(context.Background(), cfg)
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "phase")
	if got := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv"); !strings.Contains(got, "watch\tconfirmed\trun\n") {
		t.Errorf("the restart's watch, for up to watch.T of start.tsv, did not see the notice:\n%s", got)
	}
	if phase := rows[len(rows)-1].Detail; !strings.Contains(phase, "issues") || !strings.Contains(phase, "watch") {
		t.Errorf("the detail of phase %q; want the steps it ran again", phase)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv"); !strings.Contains(got, "issue.control\t2\tforge\n") {
		t.Errorf("the control issue after the restart:\n%s", got)
	}
	if got := gitOut(t, w.bare, "log", "--format=%s", "refs/heads/layup-records"); !strings.Contains(got, "takes the lease over from the run 0123456789abcdef") {
		t.Errorf("no commit takes the lease over:\n%s", got)
	}
}

// Rows 6 and 8 of the input states, with the real git: --new on a repository
// with a commit, and a restart whose start.tsv its reader refuses.
func TestTheInputStatesWithGit(t *testing.T) {
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Start(context.Background(), w.config("0.1.0-dev"))
	if len(rows) != 1 || rows[0].Name != "forge" || rows[0].Result != "fail" {
		t.Errorf("--new on a repository with a commit:\n%s", steps(rows))
	}
	work := t.TempDir()
	gitOut(t, work, "clone", "-q", "-b", "layup-records", w.bare, ".")
	os.WriteFile(filepath.Join(work, "start", "start.tsv"), []byte("name\tvalue\tsource\ncolour\tblue\trun\n"), 0o644)
	gitOut(t, work, "commit", "-q", "-am", "a broken record")
	gitOut(t, work, "push", "-q", "origin", "layup-records")
	rows = Restart(context.Background(), w.config("0.1.0-dev"))
	if got := steps(rows); !strings.Contains(got, "\nclone\tfail\t") || strings.Contains(got, "\nversion\t") {
		t.Errorf("a restart whose start.tsv its reader refuses:\n%s", got)
	}
}

func TestTheTextBlocksOfRunMd(t *testing.T) {
	spec, err := os.ReadFile("../../docs/spec/run.md")
	if err != nil {
		t.Fatal(err)
	}
	session, err := os.ReadFile("../../docs/spec/session.md")
	if err != nil {
		t.Fatal(err)
	}
	_, text, ok := strings.Cut(string(session), "\n```text probe-prompt\n")
	text, _, _ = strings.Cut(text, "```\n")
	if !ok || text != probePrompt {
		t.Errorf("the block probe-prompt of session.md:\n%s\nwant\n%s", text, probePrompt)
	}
	for name, want := range map[string]string{"start-readme": startReadme, "intake-issue": intakeBody, "control-issue": controlBody} {
		_, text, ok := strings.Cut(string(spec), "\n```text "+name+"\n")
		text, _, _ = strings.Cut(text, "```\n")
		if !ok || text != want {
			t.Errorf("the block %s of run.md:\n%s\nwant\n%s", name, text, want)
		}
	}
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tsv.Compare(blocks["run-steps"], RunStepsSchema); err != nil {
		t.Error(err)
	}
}

// pushRecords commits files on the head of layup-records of the bare target,
// as another run would, and pushes it.
// A beat of the run can come between its clone and its push, so it tries
// again on a refusal, as another run would.
func pushRecords(t *testing.T, bare string, files map[string]string, message string) {
	t.Helper()
	for try := 0; try < 5; try++ {
		work := t.TempDir()
		gitOut(t, work, "clone", "-q", "-b", "layup-records", bare, ".")
		for p, text := range files {
			os.WriteFile(filepath.Join(work, p), []byte(text), 0o644)
		}
		gitOut(t, work, "add", "-A")
		gitOut(t, work, "commit", "-q", "-m", message)
		if exec.Command("git", "-C", work, "push", "-q", "origin", "layup-records").Run() == nil {
			return
		}
	}
	t.Fatalf("five pushes of %q were refused", message)
}

// Condition 1 of the plan review: the old run made one more records commit,
// which does not change the lease row, after the restart's clone; the
// takeover is made on the commit of the last read, so it is not refused.
func TestATakeoverAfterACommitThatTheCloneDidNotSee(t *testing.T) {
	w := newWorld(t)
	w.fake.failOpen = 2
	w.start(w.config("0.1.0-dev"))
	var once sync.Once
	w.mu.Lock()
	w.hook = func(l string) {
		if strings.HasPrefix(l, "lease: held by") {
			once.Do(func() {
				pushRecords(t, w.bare, map[string]string{"note.md": "a copy of the old run\n"}, "records: a late commit of the old run")
			})
		}
	}
	w.mu.Unlock()
	cfg := w.config("0.1.0-dev")
	cfg.RunID = "fedcba9876543210"
	allDone(t, "the restart", Restart(context.Background(), cfg), "forge", "clone", "version", "lease", "probe", "phase")
	if got := show(t, w.bare, "refs/heads/layup-records", "note.md"); got != "a copy of the old run\n" {
		t.Errorf("the late commit is not in the history of the takeover: %q", got)
	}
}

// Condition 2 of the plan review: another run takes the lease during step 8;
// the next beat is refused, the lease names the other run, and the wait of
// step 8 ends with fail.
func TestALostBeatEndsTheWatch(t *testing.T) {
	w := newWorld(t)
	w.fake.comments = map[int][]map[string]any{} // no notice comes
	cfg := w.config("0.1.0-dev")
	cfg.WatchT = 60
	var once sync.Once
	w.mu.Lock()
	w.hook = func(l string) {
		if strings.HasPrefix(l, "watch: waiting") {
			once.Do(func() {
				pushRecords(t, w.bare, map[string]string{"lease.tsv": "run\thost\tversion\tstarted\theartbeat\tstate\n" +
					"cccccccccccccccc\thost-2\t0.1.0-dev\t2026-10-08T13:00:00Z\t0\theld\n"}, "lease: the run cccccccccccccccc takes the lease over")
			})
		}
	}
	w.mu.Unlock()
	rows := w.start(cfg)
	last := rows[len(rows)-1]
	if last.Name != "watch" || last.Result != "fail" || !strings.Contains(last.Detail, "cccccccccccccccc") {
		t.Errorf("a lost beat during step 8:\n%s", steps(rows))
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	waits := 0
	for _, l := range w.lines {
		if strings.HasPrefix(l, "watch: waiting") {
			waits++
		}
	}
	if waits >= 100 { // watch.T of 60 minutes is 360 reads
		t.Errorf("the watch read %d times after the beat was lost; want it ended at the loss", waits)
	}
}

// A restart whose harness register changed after Start (#199): a register
// that gained rows, with the steps of Start that the restart runs again, so
// that start.tsv is written again; and a register that lost one.
func TestARestartWithAHarnessRegisterThatChanged(t *testing.T) {
	w := newWorld(t)
	w.fake.failOpen = 2
	cfg := w.config("0.1.0-dev")
	cfg.Harnesses = nil
	rows := w.start(cfg)
	if last := rows[len(rows)-1]; last.Name != "issues" || last.Result != "fail" {
		t.Fatalf("the Start with no harness and a failed opening:\n%s", steps(rows))
	}
	cfg = w.config("0.1.0-dev") // two harnesses now
	cfg.RunID = "fedcba9876543210"
	cfg.WatchT, cfg.LeaseH = 0, 0
	allDone(t, "the restart with a harness gained", Restart(context.Background(), cfg), "forge", "clone", "version", "lease", "probe", "phase")
	got := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv")
	if strings.Contains(got, "harness.") || !strings.Contains(got, "\nissue.control\t2\tforge\n") {
		t.Errorf("start.tsv after the restart: the Start's rows, no harness row, the control issue 2:\n%s", got)
	}

	w = newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	cfg = w.config("0.1.0-dev")
	cfg.Harnesses = cfg.Harnesses[:1] // devin removed
	allDone(t, "the restart with a harness lost", Restart(context.Background(), cfg), "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning")
}

// The step rule-gaps (row 43a, task T-qk66, #217): after a complete Start, the
// restart commits the table of layup psb check on start/problem-statement.md
// as intake/psb.tsv, in place of the row phase, and releases the lease; a
// second restart finds the record and runs nothing, so the blob is the same.
func TestTheStepRuleGaps(t *testing.T) {
	w := newWorld(t)
	cfg := w.config("0.1.0-dev")
	cfg.PSB = []byte("# The problem\n\nThe service is fast.\n")
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	var progress []string
	cfg = w.config("0.1.0-dev")
	cfg.Step = func(i, n int, name string) { progress = append(progress, fmt.Sprintf("%d/%d %s", i, n, name)) }
	rows := Restart(context.Background(), cfg)
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning")
	if want := []string{"1/6 forge", "2/6 clone", "3/7 version", "4/7 lease", "5/7 probe", "6/7 rule-gaps", "7/7 meaning"}; strings.Join(progress, ", ") != strings.Join(want, ", ") {
		t.Errorf("the progress %q, want %q", progress, want)
	}
	// What layup psb check prints on the same file (internal/cli/cli.go):
	// psb.WriteTSV of psb.Check.
	var want bytes.Buffer
	psb.WriteTSV(&want, psb.Check(string(show(t, w.bare, "refs/heads/layup-records", "start/problem-statement.md"))))
	got := show(t, w.bare, "refs/heads/layup-records", "intake/psb.tsv")
	if got != want.String() || !strings.Contains(got, "\tG1\t") {
		t.Errorf("intake/psb.tsv:\n%s\nwant:\n%s", got, want.String())
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("the lease is not released after the last row")
	}
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/psb.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if gaps := again[len(again)-2]; gaps.Name != "rule-gaps" || gaps.Result != "done" || gaps.Detail != "done before" {
		t.Errorf("the second restart's row %+v, want rule-gaps done before", gaps)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/psb.tsv") != blob {
		t.Error("the second restart changed intake/psb.tsv")
	}
}

// Start is done only when the watch is written too: a start.tsv whose watch is
// still — gives the row phase, which runs the watch again (row 43a).
func TestARestartWhoseWatchIsNotWrittenRunsPhase(t *testing.T) {
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	start := show(t, w.bare, "refs/heads/layup-records", "start/start.tsv")
	i := strings.Index(start, "\nwatch\t")
	pushRecords(t, w.bare, map[string]string{"start/start.tsv": start[:i] + "\nwatch\t—\trun\n"}, "records: no watch yet")
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "phase")
}

// The step meaning (row 43b, task T-1u2v, #218): after rule-gaps, one session
// of systems-architect on the world's harness; its quote is found once, so
// intake/meaning.tsv holds its row with code's offsets, and the lease is
// released; a second restart finds the record and runs nothing.
func TestTheStepMeaning(t *testing.T) {
	w := newWorld(t)
	os.WriteFile(filepath.Join(w.dir, "bin", "meaning.tsv"), []byte("kind\tquote\tquestion\nunmeasurable\tThe service is fast.\tHow fast?\n"), 0o644)
	cfg := w.config("0.1.0-dev")
	cfg.PSB = []byte("# The problem\n\nThe service is fast.\n")
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning")
	if d := rows[len(rows)-1].Detail; d != "gaps of meaning: 1, in intake/meaning.tsv" {
		t.Errorf("the detail of meaning %q", d)
	}
	got := show(t, w.bare, "refs/heads/layup-records", "intake/meaning.tsv")
	if !strings.HasPrefix(got, "id\tkind\tstart\tend\tquote\tquestion\tsession\nM-001\tunmeasurable\t15\t35\tThe service is fast.\tHow fast?\tS-") || strings.Count(got, "\n") != 2 {
		t.Errorf("intake/meaning.tsv:\n%s", got)
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tsystems-architect\tclaude\t1.0.0\tm1\t") {
		t.Error("sessions.tsv has no start row of a session of systems-architect")
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("the lease is not released after the last row")
	}
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/meaning.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if last := again[len(again)-1]; last.Name != "meaning" || last.Result != "done" || last.Detail != "done before" {
		t.Errorf("the second restart's last row %+v, want meaning done before", last)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/meaning.tsv") != blob {
		t.Error("the second restart changed intake/meaning.tsv")
	}
}

// A quote that occurs twice refuses the session's file: meaning is fail,
// naming the session, and nothing of the file is committed; with no route for
// systems-architect, meaning is fail (pair) and no session starts. Each fail
// leaves the lease held (run.md, the known limit).
func TestTheStepMeaningFails(t *testing.T) {
	for _, c := range []struct{ name, meaning, detail string }{
		{"a quote that occurs twice", "kind\tquote\tquestion\nunmeasurable\tfast\tHow fast?\n", `: refused file: row 1: the quote "fast" occurs 2 times`},
		{"no route", "", "pair: no admitted pair"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := newWorld(t)
			if c.meaning != "" {
				os.WriteFile(filepath.Join(w.dir, "bin", "meaning.tsv"), []byte(c.meaning), 0o644)
			}
			cfg := w.config("0.1.0-dev")
			cfg.PSB = []byte("# The problem\n\nThe service is fast.\nIt is fast.\n")
			allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
			cfg = w.config("0.1.0-dev")
			if c.meaning == "" {
				cfg.Routing = nil
			}
			rows := Restart(context.Background(), cfg)
			last := rows[len(rows)-1]
			if last.Name != "meaning" || last.Result != "fail" || !strings.Contains(last.Detail, c.detail) {
				t.Fatalf("the last row %+v; want meaning fail with %q", last, c.detail)
			}
			if c.meaning != "" && !strings.HasPrefix(last.Detail, "the session S-") {
				t.Errorf("the detail %q does not name the session", last.Detail)
			}
			if strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "intake/meaning.tsv") {
				t.Error("a refused file is a record")
			}
			if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\theld\n") {
				t.Error("the lease is not held after a fail")
			}
		})
	}
}
