//go:build integration

package run

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
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
	gitOut(t, "/", "init", "-q", "--bare", w.bare)
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
		Register: reg, Harnesses: [][]string{{"claude", "10.0", "60"}, {"devin", "", "30"}}, Forge: adapter, Clock: clock,
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
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "phase")
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
