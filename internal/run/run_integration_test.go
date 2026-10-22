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
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/forge/github"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/psb"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
	workpkg "github.com/pharzam/layup/internal/work"
)

// fakeGitHub plays the calls of M2a for the steps of Start and the restart
// (docs/spec/run.md, the acceptance row Start): its branches are those of the
// bare target repository, so the stand-in Operator's push shows.
type fakeGitHub struct {
	t          *testing.T
	bare       string
	mu         sync.Mutex
	issues     int
	failOpen   int    // the number of the issue whose opening fails, or 0
	failNote   string // a comment whose body holds this fails once, or ""
	comments   map[int][]map[string]any
	posted     []string                    // the bodies of the comments that the run posted, on any issue (row 43b)
	rules      map[string]string           // the JSON list that the rules of a branch answer, by the branch (row 57a); [] when none
	noActivity bool                        // the read of the activity answers 404 (row 57a)
	perms      map[string]string           // the permissions of the installation; the three of M2a when nil (task T-cak3)
	statuses   map[string][]map[string]any // the commit statuses of each SHA, newest first (row 66f)
	failStatus int                         // the next posts of a status that answer 500, after okStatus have passed
	okStatus   int
	reqs       []string                    // every request, "METHOD /path", in order (row 66g)
	checkRuns  map[string][]map[string]any // the latest check runs of each SHA (row 66g)
	pulls      pullRoutes                  // the pull requests (row 66g)
	issueTexts [][2]string                 // the title and the body of each issue opened (task T-jy6f)
	defBranch  string                      // the default branch that the repository answers; "main" when empty (task T-yngf)
}

func (f *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	call := r.Method + " " + r.URL.Path
	f.reqs = append(f.reqs, call)
	switch {
	case call == "GET /repos/acme/target/installation":
		perms := map[string]string{"contents": "write", "issues": "write", "metadata": "read"}
		if f.perms != nil {
			perms = f.perms
		}
		json.NewEncoder(w).Encode(map[string]any{"id": 7, "permissions": perms})
	case call == "POST /app/installations/7/access_tokens":
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"token": "ghs_runtest", "expires_at": "2099-01-01T00:00:00Z"})
	case call == "GET /repos/acme/target":
		def := "main"
		if f.defBranch != "" {
			def = f.defBranch
		}
		json.NewEncoder(w).Encode(map[string]any{"default_branch": def, "visibility": "public"})
	case call == "GET /repos/acme/target/branches":
		out, err := exec.Command("git", "-C", f.bare, "for-each-ref", "--format=%(refname:short) %(objectname)", "refs/heads").Output()
		if err != nil {
			http.Error(w, `{"message":"no bare repository"}`, 500)
			return
		}
		list := []map[string]any{}
		for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
			if name, sha, ok := strings.Cut(line, " "); ok {
				list = append(list, map[string]any{"name": name, "commit": map[string]string{"sha": sha}})
			}
		}
		json.NewEncoder(w).Encode(list)
	case strings.HasPrefix(call, "GET /repos/acme/target/rules/branches/"):
		text := f.rules[strings.TrimPrefix(r.URL.Path, "/repos/acme/target/rules/branches/")]
		if text == "" {
			text = "[]"
		}
		w.Write([]byte(text))
	case call == "GET /repos/acme/target/activity":
		if f.noActivity {
			http.Error(w, `{"message":"Not Found"}`, 404)
			return
		}
		w.Write([]byte("[]"))
	case strings.HasPrefix(call, "GET /users/"):
		ids := map[string]int{"pharzam": 101, "carol": 303, "layup-agent[bot]": 9001, "dave": 707}
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
		var in struct{ Title, Body string }
		json.NewDecoder(r.Body).Decode(&in)
		f.issueTexts = append(f.issueTexts, [2]string{in.Title, in.Body})
		f.issues++
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"number": f.issues})
	case strings.HasPrefix(call, "POST /repos/acme/target/issues/") && strings.HasSuffix(call, "/comments"):
		var in struct{ Body string }
		json.NewDecoder(r.Body).Decode(&in)
		if f.failNote != "" && strings.Contains(in.Body, f.failNote) {
			f.failNote = ""
			http.Error(w, `{"message":"Server Error"}`, 500)
			return
		}
		f.posted = append(f.posted, in.Body)
		// The comment is the App's, and a later read of its issue sees it
		// (row 46: the step batch reads the Intake issue).
		n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/repos/acme/target/issues/"), "/comments"))
		id := 600 + len(f.posted)
		f.comments[n] = append(f.comments[n], map[string]any{"id": id, "user": map[string]any{"id": 9001, "login": "layup-agent[bot]"},
			"performed_via_github_app": map[string]any{"slug": "layup-agent"},
			"created_at":               "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": in.Body})
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": id})
	case strings.HasPrefix(call, "GET /repos/acme/target/issues/") && strings.HasSuffix(call, "/comments"):
		n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/repos/acme/target/issues/"), "/comments"))
		list := f.comments[n]
		if list == nil {
			list = []map[string]any{}
		}
		json.NewEncoder(w).Encode(list)
	case strings.HasPrefix(call, "POST /repos/acme/target/statuses/"):
		if f.okStatus > 0 {
			f.okStatus--
		} else if f.failStatus > 0 {
			f.failStatus--
			http.Error(w, `{"message":"Server Error"}`, 500)
			return
		}
		var in struct{ State, Context, Description string }
		json.NewDecoder(r.Body).Decode(&in)
		sha := strings.TrimPrefix(r.URL.Path, "/repos/acme/target/statuses/")
		if f.statuses == nil {
			f.statuses = map[string][]map[string]any{}
		}
		f.statuses[sha] = append([]map[string]any{{"context": in.Context, "state": in.State, "description": in.Description,
			"creator": map[string]any{"id": 9001, "login": "layup-agent[bot]"}}}, f.statuses[sha]...)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": len(f.statuses[sha])})
	case strings.HasPrefix(call, "GET /repos/acme/target/commits/") && strings.HasSuffix(call, "/statuses"):
		sha := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/repos/acme/target/commits/"), "/statuses")
		list := f.statuses[sha]
		if list == nil {
			list = []map[string]any{}
		}
		json.NewEncoder(w).Encode(list)
	case strings.HasPrefix(call, "GET /repos/acme/target/commits/") && strings.HasSuffix(call, "/check-runs"):
		sha := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/repos/acme/target/commits/"), "/check-runs")
		runs := f.checkRuns[sha]
		if runs == nil {
			runs = f.checkRuns["*"] // the runs of every SHA that has none of its own (task T-ehb4)
		}
		if runs == nil {
			runs = []map[string]any{}
		}
		json.NewEncoder(w).Encode(map[string]any{"check_runs": runs})
	case f.pulls.serve(f, w, r):
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
	clock    fastClock         // the world's one clock: no config of it goes back in time
}

func newWorld(t *testing.T) *world { return newWorldWith(t, nil) }

// newWorldWith makes a world whose baseline also holds files, by path (task
// T-4c5c, #276: the step markers scans the pinned tree, which is the baseline's).
func newWorldWith(t *testing.T, files map[string]string) *world {
	w := &world{t: t, dir: t.TempDir(), web: t.TempDir(),
		clock: fastClock{begin: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC), real: time.Now()}}
	src := t.TempDir()
	os.WriteFile(filepath.Join(src, "README.md"), []byte("# Armature\n"), 0o644)
	for path, text := range files {
		os.MkdirAll(filepath.Dir(filepath.Join(src, path)), 0o755)
		os.WriteFile(filepath.Join(src, path), []byte(text), 0o644)
	}
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
	// The target receives pushes, and the git that receives one starts its own
	// maintenance in the background unless the repository says no: the -c of the
	// pushing call does not reach it (T-2hm3).
	gitOut(t, w.bare, "config", "receive.autogc", "false")
	gitOut(t, w.bare, "config", "maintenance.auto", "false")
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
	// probe as probeHarness does; and a session of a phase with the file of
	// bin/ of its output's name, which a test may write (bin/meaning.tsv: no
	// gap), else, for the numbering, a span table that covers the problem
	// statement, one fact per line that is not empty (row 44), for the draft
	// one requirement per need fact of the span table (row 45a), for the
	// completeness review no item (row 45b), for the step markers no proposal
	// (row 53), and for the review of the answers each answered question closed
	// (row 47b).
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
	in=$(printf '%s\n' "$p" | sed -n 's/^The problem statement: //p')
	f="$(dirname "$0")/$(basename "$out")"
	sp=$(printf '%s\n' "$p" | sed -n 's/^The span table: //p')
	# Each session of a phase reads the problem statement (intake.md, The
	# sessions of Intake; round 1 of #224).
	[ -f "$in" ] || exit 3
	# The prose session of the step setup (row 54, task T-eahr): repo/ is the target on its task
	# branch at the head of layup-setup, and the files of the F- rows of the stop table, in
	# bin/files/, are written below the output directory.
	if printf '%s\n' "$p" | grep -q '^The output directory: '; then
		printf '%s\n' "$p" > "$(dirname "$0")/prose-prompt.txt"
		dir=$(printf '%s\n' "$p" | sed -n 's/^The output directory: //p')
		stop=$(printf '%s\n' "$p" | sed -n 's/^The stop table: //p')
		[ -f docs/setup/armature.pin ] || exit 5
		[ "$(git rev-parse --abbrev-ref HEAD | cut -c1-5)" = task/ ] || exit 6
		for want in $(awk -F '\t' 'NR > 1 { print $4 }' "$stop"); do
			src="$(dirname "$0")/files/$want"
			[ -f "$src" ] || exit 7
			mkdir -p "$dir/$(dirname "$want")" && cp "$src" "$dir/$want" || exit 8
		done
		exit 0
	fi
	# The session of the step markers reads the list of the markers too (row 53).
	if printf '%s\n' "$p" | grep -q '^The list of markers: '; then
		ls=$(printf '%s\n' "$p" | sed -n 's/^The list of markers: //p')
		[ "$(head -n 1 "$ls")" = "$(printf 'id\tfile\tline\tmarker\tvalue\tstart\tend\tsession')" ] || exit 4
	fi
	# The revision (row 66i, task T-8rm6): each attempt's prompt is kept for the
	# test (the finding of a refused attempt is in the next one).
	if [ "$(basename "$out")" = revision.tsv ]; then
		printf '%s\n' "$p" > "$(dirname "$0")/revision-prompt.txt"
	fi
	if [ -f "$f" ]; then
		cat "$f" > "$out"
	elif [ "$(basename "$out")" = revision.tsv ]; then
		dr=$(printf '%s\n' "$p" | sed -n 's/^The draft before the revision: //p')
		# The record's row is the session's with id first and session last.
		LC_ALL=C awk -F '\t' 'BEGIN { OFS = "\t" } NR == 1 { print "kind", "statement", "covers", "criterion" } NR > 1 { print $2, $3, $4, $5 }' "$dr" > "$out"
	elif [ "$(basename "$out")" = answers-review.tsv ]; then
		qs=$(printf '%s\n' "$p" | sed -n 's/^The questions: //p')
		LC_ALL=C awk -F '\t' 'BEGIN { printf "question\tstate\treason\n" } NR > 1 && $13 == "answered" { printf "%s\tclosed\tthe answer is clear\n", $1 }' "$qs" > "$out"
	elif [ "$(basename "$out")" = coverage.tsv ]; then
		printf 'kind\tfact\trequirement\tquestion\n' > "$out"
	elif [ "$(basename "$out")" = markers.tsv ]; then
		printf 'id\tvalue\tquote\n' > "$out"
	elif [ "$(basename "$out")" = draft.tsv ]; then
		LC_ALL=C awk -F '\t' 'BEGIN { printf "kind\tstatement\tcovers\tcriterion\n" }
			NR > 1 && $4 == "need" { printf "requirement\tThe statement of fact %s.\t%s\tThe criterion of fact %s.\n", $3, $3, $3 }' "$sp" > "$out"
	else
		LC_ALL=C awk 'BEGIN { printf "start\tend\tfact\tclass\treason\n"; o = 0; n = 0 }
			{ l = length($0); if (l > 0) { n++; printf "%d\t%d\t%d\tneed\t\342\200\224\n", o, o + l, n }; o += l + 1 }' "$in" > "$out"
	fi ;;
esac
`), 0o755)
	return w
}

func (w *world) config(version string) Config {
	clock := w.clock
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
			"", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""},
			// devin runs the same scripts, a second harness for the
			// completeness review (row 45b).
			{"devin", "", "30", filepath.Join(w.dir, "bin", "harness"), "stdin", filepath.Join(w.dir, "bin", "version"),
				"", "", "LAYUP-TEST-RULES.md", "", "none", "subscription", ""}},
		Models: [][]string{{"claude", "m1", "1000000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""},
			{"devin", "m2", "1000000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""}},
		Routing: [][]string{{"systems-architect", "reasoning", "1", "claude", "m1"}, {"product-owner", "reasoning", "1", "claude", "m1"},
			// The role's first pair is on the draft's harness, which the
			// review passes over.
			{"qa-engineer", "reasoning", "1", "claude", "m1"}, {"qa-engineer", "reasoning", "2", "devin", "m2"}},
		Forge: adapter, Clock: clock,
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

// beforeSetup gives the rows of a restart that ran to the end of Intake without
// the last, the step setup, which these worlds cannot do: their baseline holds no
// file of a setup and their scripted harness writes no prose, so the step is fail
// at S04, whose file docs/adr/README.md the baseline lacks (row 54, task T-eahr).
// The rows before it are asserted as they were.
func beforeSetup(t *testing.T, rows []Step) []Step {
	t.Helper()
	last := rows[len(rows)-1]
	if last.Name != "setup" || last.Result != "fail" || !strings.HasPrefix(last.Detail, "S04: docs/adr/README.md of the baseline") {
		t.Errorf("the last row of a restart that ends Intake: %+v; want the step setup, fail at S04 for the file that its baseline lacks", last)
	}
	return rows[:len(rows)-1]
}

// The clock of a world does not go back from one config of it to the next. A
// restart gets a new config, and a clock that began at 12:00:00 again let two
// restarts of one run write the same lease row, second for second (T-0p5p).
func TestTheClockOfAWorldDoesNotGoBack(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	first := w.config("0.1.0-dev")
	time.Sleep(20 * time.Millisecond) // 12 seconds of the clock
	seen := first.Clock.Now()
	if got := w.config("0.1.0-dev").Clock.Now(); got.Before(seen) {
		t.Errorf("the clock of a second config reads %s, before the %s that the first read", got.Format(time.RFC3339Nano), seen.Format(time.RFC3339Nano))
	}
}

// A push into the target of a world starts no maintenance of git in it. The
// git that receives a push runs `git maintenance run --auto --quiet --detach`
// in the target, and the -c maintenance.auto=false of the pushing call does
// not reach it: a local transport clears the values of -c before it starts the
// receiving git. From git 2.54 (measured on 2.55, CI's) that maintenance repacks
// the target in the background, and the removal of the test's directory races
// with it (T-2hm3).
func TestAPushIntoTheTargetOfAWorldStartsNoMaintenance(t *testing.T) {
	w := newWorld(t)
	src := t.TempDir()
	writeFile(t, filepath.Join(src, "a.txt"), "a\n")
	gitOut(t, src, "init", "-q", "-b", "main")
	gitOut(t, src, "add", "-A")
	gitOut(t, src, "commit", "-q", "-m", "one")
	trace := filepath.Join(t.TempDir(), "trace2.json")
	t.Setenv("GIT_TRACE2_EVENT", trace) // gitOut starts from this environment; internal/git does not
	gitOut(t, src, "push", "-q", "file://"+w.bare, "main:refs/heads/probe")
	b, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	// The receiving git wrote an event of its own, so the trace is that of the
	// process that would start the maintenance.
	if !strings.Contains(got, `"name":"receive-pack"`) {
		t.Fatalf("the trace of the push holds no event of receive-pack:\n%s", got)
	}
	for _, started := range []string{`"maintenance","run"`, `"name":"maintenance"`} {
		if strings.Contains(got, started) {
			t.Errorf("a push into the target started the maintenance of git: %s is in the trace", started)
		}
	}
}

// The demo of row 25b: against the httptest forge and a local bare
// repository, Start makes the first records commit with the pin, the briefs,
// approvers.tsv, start.tsv and the lease row.
func TestStartMakesTheFirstRecordsCommit(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.2.0"))
	if got := steps(rows); !strings.HasPrefix(got, "step\tresult\tdetail\nforge\tdone\t") || !strings.Contains(got, "\nclone\tdone\t") || !strings.Contains(got, "\nversion\tfail\t") {
		t.Errorf("the restart with another version:\n%s", got)
	}
	// The same version: the lease is released, so it is taken at once; each
	// step of Start is done.
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
}

// A run that stopped at opening: the restart takes the lease over from the
// stopped run and runs step 7 again for that issue.
func TestARestartAfterARunThatStoppedAtOpening(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	scaffold, err := os.ReadFile("../../docs/spec/scaffold.md")
	if err != nil {
		t.Fatal(err)
	}
	for name, c := range map[string]struct{ doc, want string }{"start-readme": {string(spec), startReadme}, "intake-issue": {string(spec), intakeBody},
		"control-issue": {string(spec), controlBody}, "scaffold-readme": {string(scaffold), scaffoldReadme}} {
		_, text, ok := strings.Cut(c.doc, "\n```text "+name+"\n")
		text, _, _ = strings.Cut(text, "```\n")
		if !ok || text != c.want {
			t.Errorf("the block %s:\n%s\nwant\n%s", name, text, c.want)
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
			os.MkdirAll(filepath.Dir(filepath.Join(work, p)), 0o755) // a record of a directory that the branch lacks (task T-4c5c)
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	// With one admitted harness left, the completeness review has no pair on
	// another harness (row 45b): the restart stops there.
	lost := Restart(context.Background(), cfg)
	allDone(t, "the restart with a harness lost", lost[:len(lost)-1], "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft")
	if last := lost[len(lost)-1]; last.Name != "completeness" || last.Result != "fail" || last.Detail != "not run: no admitted pair on another harness" {
		t.Errorf("the last row of the restart with a harness lost %+v; want completeness fail, not run", last)
	}
}

// The step rule-gaps (row 43a, task T-qk66, #217): after a complete Start, the
// restart commits the table of layup psb check on start/problem-statement.md
// as intake/psb.tsv, in place of the row phase, and releases the lease; a
// second restart finds the record and runs nothing, so the blob is the same.
func TestTheStepRuleGaps(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	cfg := w.config("0.1.0-dev")
	cfg.PSB = []byte("# The problem\n\nThe service is fast.\n")
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	var progress []string
	cfg = w.config("0.1.0-dev")
	cfg.Step = func(i, n int, name string) { progress = append(progress, fmt.Sprintf("%d/%d %s", i, n, name)) }
	rows := Restart(context.Background(), cfg)
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if want := []string{"1/6 forge", "2/6 clone", "3/22 version", "4/22 lease", "5/22 probe", "6/22 rule-gaps", "7/22 meaning", "8/22 numbering", "9/22 spans", "10/22 draft", "11/22 completeness", "12/22 markers", "13/22 batch", "14/22 answers"}; strings.Join(progress, ", ") != strings.Join(want, ", ") {
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
	if gaps := rowOf(again, "rule-gaps"); gaps.Name != "rule-gaps" || gaps.Result != "done" || gaps.Detail != "done before" {
		t.Errorf("the second restart's row %+v, want rule-gaps done before", gaps)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/psb.tsv") != blob {
		t.Error("the second restart changed intake/psb.tsv")
	}
}

// Start is done only when the watch is written too: a start.tsv whose watch is
// still — gives the row phase, which runs the watch again (row 43a).
func TestARestartWhoseWatchIsNotWrittenRunsPhase(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	w := newWorld(t)
	os.WriteFile(filepath.Join(w.dir, "bin", "meaning.tsv"), []byte("kind\tquote\tquestion\nunmeasurable\tThe service is fast.\tHow fast?\n"), 0o644)
	cfg := w.config("0.1.0-dev")
	cfg.PSB = []byte("# The problem\n\nThe service is fast.\n")
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if d := rowOf(rows, "meaning").Detail; d != "1 gaps of meaning, in intake/meaning.tsv" {
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
	if m := rowOf(again, "meaning"); m.Result != "done" || m.Detail != "done before" {
		t.Errorf("the second restart's row %+v, want meaning done before", m)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/meaning.tsv") != blob {
		t.Error("the second restart changed intake/meaning.tsv")
	}
}

// A quote that occurs twice refuses the session's file: meaning is fail,
// naming the session, and nothing of the file is committed; with no route for
// systems-architect, meaning is fail (pair) and no session starts. Each fail
// releases the lease (run.md, The lease and fencing; task T-yngf, #479: it was held, the known limit).
func TestTheStepMeaningFails(t *testing.T) {
	t.Parallel()
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
			if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
				t.Error("the lease is not released after a fail") // T-yngf: was held, the known limit
			}
			// With no pair, no session of systems-architect starts (round 1 of #218).
			if c.meaning == "" && strings.Contains(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tsystems-architect\t") {
				t.Error("a session of systems-architect started with no pair")
			}
		})
	}
}

// rowOf gives the row of the step name, or the zero Step.
func rowOf(rows []Step, name string) Step {
	for _, r := range rows {
		if r.Name == name {
			return r
		}
	}
	return Step{}
}

// The steps numbering and spans (row 44, task T-hiv9, #219): a session of
// product-owner writes the span table, code adds each span's text, and the
// span check of the records commits its table; the lease is released, and a
// second restart finds both records and runs nothing.
func TestTheStepsNumberingAndSpans(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if n, s := rowOf(rows, "numbering").Detail, rowOf(rows, "spans").Detail; n != "2 spans, 2 facts, in spec/spans.tsv" || s != "the span check passed, in spec/check.tsv" {
		t.Errorf("the details %q and %q", n, s)
	}
	// The world's brief is "# The problem\n\nbyte for byte\n".
	if got := show(t, w.bare, "refs/heads/layup-records", "spec/spans.tsv"); got != "start\tend\tfact\tclass\treason\ttext\n0\t13\t1\tneed\t—\t# The problem\n15\t28\t2\tneed\t—\tbyte for byte\n" {
		t.Errorf("spec/spans.tsv:\n%q", got)
	}
	// spec/check.tsv is the table of spec.Check on the records, by its schema.
	spans, err := spec.Read([]byte(show(t, w.bare, "refs/heads/layup-records", "spec/spans.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	var want bytes.Buffer
	tsv.Write(&want, spec.CheckSchema, spec.Check(spans, []byte(show(t, w.bare, "refs/heads/layup-records", "start/problem-statement.md"))))
	if got := show(t, w.bare, "refs/heads/layup-records", "spec/check.tsv"); got != want.String() || strings.Contains(got, "fail") {
		t.Errorf("spec/check.tsv:\n%s\nwant:\n%s", got, want.String())
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tproduct-owner\tclaude\t1.0.0\tm1\t") ||
		!strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("no session of product-owner, or the lease not released")
	}
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if n, s := rowOf(again, "numbering"), rowOf(again, "spans"); n.Detail != "done before" || s.Detail != "done before" || s.Result != "done" {
		t.Errorf("the second restart's rows %+v and %+v; want both done before", n, s)
	}
}

// The demo of row 44: after a failed span check the run stops at spans, so no
// row after it runs and no batch is posted, with spec/check.tsv committed and
// the lease released (T-yngf; it was held); the next restart runs numbering again
// with a new session, and spans passes, both records replaced.
func TestAFailedSpanCheckRunsNumberingAgain(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	// Two spans over byte 10 of "# The problem\n\nbyte for byte\n".
	bad := filepath.Join(w.dir, "bin", "spans.tsv")
	os.WriteFile(bad, []byte("start\tend\tfact\tclass\treason\n0\t13\t1\tneed\t—\n10\t28\t2\tneed\t—\n"), 0o644)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	last := rows[len(rows)-1]
	if last.Name != "spans" || last.Result != "fail" || last.Detail != "the span check: cover fail: byte 10 is in 2 spans" {
		t.Fatalf("the last row %+v; want spans fail, the first fail row of the check", last)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "spec/check.tsv"); !strings.Contains(got, "\ncover\tfail\tbyte 10 is in 2 spans\n") {
		t.Errorf("spec/check.tsv:\n%s", got)
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("the lease is not released after the fail") // T-yngf: was held, the known limit
	}
	// A rerun whose session writes the same table fails the same way: its
	// table equals the last spec/check.tsv, which is kept as it is.
	same := Restart(context.Background(), w.config("0.1.0-dev"))
	if last := same[len(same)-1]; last.Name != "spans" || last.Result != "fail" || last.Detail != "the span check: cover fail: byte 10 is in 2 spans" {
		t.Fatalf("a rerun with the same table: the last row %+v; want spans fail as before", last)
	}
	os.Remove(bad) // the next session covers the file
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart after the fail", again, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if n := rowOf(again, "numbering").Detail; n != "2 spans, 2 facts, in spec/spans.tsv" {
		t.Errorf("numbering after a failed check: %q; want a new session's spans", n)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "spec/check.tsv"); strings.Contains(got, "fail") {
		t.Errorf("spec/check.tsv after the rerun:\n%s", got)
	}
	// Three sessions of the numbering, each rerun's a new one, and the
	// draft's after the check passed (row 45a).
	if n := strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tproduct-owner\t"); n != 4 {
		t.Errorf("%d sessions of product-owner; want 4: three of the numbering and the draft's", n)
	}
}

// With no route for product-owner, numbering is fail (pair) and no session of
// product-owner starts.
func TestNumberingWithNoPair(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	cfg := w.config("0.1.0-dev")
	cfg.Routing = cfg.Routing[:1] // systems-architect only
	rows := Restart(context.Background(), cfg)
	if last := rows[len(rows)-1]; last.Name != "numbering" || last.Result != "fail" || !strings.HasPrefix(last.Detail, "pair: no admitted pair") {
		t.Errorf("the last row %+v; want numbering fail (pair)", last)
	}
	if strings.Contains(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tproduct-owner\t") {
		t.Error("a session of product-owner started with no pair")
	}
}

// The step draft (row 45a, task T-2et0, #220): one session of product-owner
// on the problem statement and the span table; its draft covers each need, so
// spec/draft.tsv holds its rows with their IDs, and the lease is released; a
// second restart finds the record and runs nothing.
func TestTheStepDraft(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if d := rowOf(rows, "draft").Detail; d != "2 rows, in spec/draft.tsv" {
		t.Errorf("the detail of draft %q", d)
	}
	got := show(t, w.bare, "refs/heads/layup-records", "spec/draft.tsv")
	if !strings.HasPrefix(got, "id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.\tS-") ||
		!strings.Contains(got, "\nR-002\trequirement\tThe statement of fact 2.\t2\t") {
		t.Errorf("spec/draft.tsv:\n%s", got)
	}
	if strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tproduct-owner\t") != 2 ||
		!strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("not two sessions of product-owner (the numbering and the draft), or the lease not released")
	}
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if d := rowOf(again, "draft"); d.Result != "done" || d.Detail != "done before" {
		t.Errorf("the second restart's row %+v; want draft done before", d)
	}
}

// The demo of row 45a: a draft whose covers names no fact is refused: draft is
// fail, naming the session, and nothing of the file is committed.
func TestADraftThatCoversNoFactIsRefused(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	os.WriteFile(filepath.Join(w.dir, "bin", "draft.tsv"), []byte("kind\tstatement\tcovers\tcriterion\nrequirement\tS.\t1\tC.\nrequirement\tS.\t2\tC.\nrequirement\tS.\t9\tC.\n"), 0o644)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	last := rows[len(rows)-1]
	if last.Name != "draft" || last.Result != "fail" || !strings.HasPrefix(last.Detail, "the session S-") || !strings.HasSuffix(last.Detail, ": refused file: row 3: covers names 9, which is no fact") {
		t.Fatalf("the last row %+v; want draft fail, naming the session and the row", last)
	}
	if strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "spec/draft.tsv") {
		t.Error("a refused draft is a record")
	}
}

// The completeness review (row 45b, task T-0mey, #221): its pair is the first
// admitted pair of qa-engineer whose harness is not the draft's, even when the
// role's first pair is on that harness; its file is intake/coverage.tsv, and a
// second restart finds the record and runs nothing.
func TestTheStepCompleteness(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	os.WriteFile(filepath.Join(w.dir, "bin", "coverage.tsv"), []byte("kind\tfact\trequirement\tquestion\nclass\t1\t—\tIs fact 1 a need?\n"), 0o644)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if d := rowOf(rows, "completeness").Detail; d != "1 items, in intake/coverage.tsv" {
		t.Errorf("the detail of completeness %q", d)
	}
	sessions := show(t, w.bare, "refs/heads/layup-records", "sessions.tsv")
	if !strings.Contains(sessions, "\tqa-engineer\tdevin\t1.0.0\tm2\t") || strings.Contains(sessions, "\tqa-engineer\tclaude\t") ||
		!strings.Contains(sessions, "\tproduct-owner\tclaude\t") {
		t.Errorf("the review is not on devin, the harness that did not number or draft:\n%s", sessions)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "intake/coverage.tsv"); !strings.HasPrefix(got, "id\tkind\tfact\trequirement\tquestion\tsession\nC-001\tclass\t1\t—\tIs fact 1 a need?\tS-") {
		t.Errorf("intake/coverage.tsv:\n%s", got)
	}
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if c := rowOf(again, "completeness"); c.Result != "done" || c.Detail != "done before" {
		t.Errorf("the second restart's row %+v; want completeness done before", c)
	}
}

// The demo of row 45b: with no admitted pair of qa-engineer on another
// harness, completeness is fail, the last row (no batch), no session of
// qa-engineer starts, and the lease is released (T-yngf; it was held).
func TestTheCompletenessReviewWithNoOtherHarness(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name   string
		routes int
	}{
		{"qa-engineer on the draft's harness only", 3},
		{"no route of qa-engineer at all", 2}, // round 1 of #221, note 3
	} {
		t.Run(c.name, func(t *testing.T) {
			w := newWorld(t)
			allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
			cfg := w.config("0.1.0-dev")
			cfg.Routing = cfg.Routing[:c.routes]
			rows := Restart(context.Background(), cfg)
			if last := rows[len(rows)-1]; last.Name != "completeness" || last.Result != "fail" || last.Detail != "not run: no admitted pair on another harness" {
				t.Fatalf("the last row %+v; want completeness fail, not run", last)
			}
			if strings.Contains(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tqa-engineer\t") {
				t.Error("a session of qa-engineer started")
			}
			if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
				t.Error("the lease is not released after the fail") // T-yngf: was held, the known limit
			}
		})
	}
}

// A review's file that names a fact of no span is refused: completeness is
// fail, naming the session, and no intake/coverage.tsv is committed.
func TestACoverageFileThatNamesNoFactIsRefused(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	os.WriteFile(filepath.Join(w.dir, "bin", "coverage.tsv"), []byte("kind\tfact\trequirement\tquestion\nuncovered\t9\t—\tWhat of 9?\n"), 0o644)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	last := rows[len(rows)-1]
	if last.Name != "completeness" || last.Result != "fail" || !strings.HasPrefix(last.Detail, "the session S-") ||
		!strings.HasSuffix(last.Detail, ": refused file: row 1: the fact 9 is no fact of spec/spans.tsv") {
		t.Fatalf("the last row %+v; want completeness fail, naming the session and the row", last)
	}
	if strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "intake/coverage.tsv") {
		t.Error("a refused file is a record")
	}
}

// intakeComments gives the bodies of the App's comments on the Intake issue,
// issue 1 of the world.
func intakeComments(w *world) []string {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	var out []string
	for _, c := range w.fake.comments[1] {
		if app, _ := c["performed_via_github_app"].(map[string]any); app != nil && app["slug"] == "layup-agent" {
			out = append(out, c["body"].(string))
		}
	}
	return out
}

// The demo of row 46 (task T-8r3u, #222): after the completeness review, one
// comment of the App on the Intake issue holds the merged questions in two
// blocks, after the records commit of questions.tsv; a second restart finds
// both and posts nothing.
func TestTheStepBatch(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	questions := show(t, w.bare, "refs/heads/layup-records", "questions.tsv")
	n := strings.Count(questions, "\n") - 1
	if d := rowOf(rows, "batch").Detail; d != fmt.Sprintf("%d questions, in questions.tsv and the comment of the batch", n) || n < 4 {
		t.Errorf("the detail of batch %q, %d questions", d, n)
	}
	for _, f := range []string{"band.B", "band.U", "appetite.hours"} {
		if !strings.Contains(questions, "\tintent\t"+f+"\tidea-owner\tbefore-delivery\tintent\t—\t—\t") {
			t.Errorf("questions.tsv has no question of the intent field %s:\n%s", f, questions)
		}
	}
	// With no marker in the baseline the Operator's block asks S01-stack alone, whatever rule G1
	// found (row 53, task T-4c5c).
	stack := fmt.Sprintf("Q-%03d\tsetup\tS01-stack\toperator\tbefore-delivery\tsetup\t—\t—\t%s\t", n, workpkg.S01Questions[0].Text)
	if !strings.Contains(questions, "\n"+stack) || strings.Count(questions, "\tsetup\t") != 2 { // the source and the kind of the one row
		t.Errorf("questions.tsv has no question of S01-stack as its last row:\n%s", questions)
	}
	c := intakeComments(w)
	if len(c) != 1 || !strings.HasPrefix(c[0], batchHeader+"\n") || !strings.Contains(c[0], "\n## For the idea owner\n") ||
		!strings.Contains(c[0], "\n## For the Operator\n\n"+fmt.Sprintf("Q-%03d: ", n)+workpkg.S01Questions[0].Text+"\n\n") || !strings.Contains(c[0], fmt.Sprintf("\nQ-%03d: ", n)) ||
		strings.Contains(c[0], "`Q-NNN: accept`") {
		t.Fatalf("the App's comments on the Intake issue %q; want one batch of %d questions in two blocks, S01-stack in the Operator's", c, n)
	}
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:questions.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if b := rowOf(again, "batch"); b.Result != "done" || b.Detail != "done before" || len(intakeComments(w)) != 1 {
		t.Errorf("the second restart's row %+v, %d comments; want batch done before and one comment", b, len(intakeComments(w)))
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:questions.tsv") != blob {
		t.Error("the second restart changed questions.tsv")
	}
}

// A run that stopped between the records commit of questions.tsv and the post
// leaves the record and no comment: the next restart posts the comment from
// the record, with no new commit of questions.tsv.
func TestABatchWithNoCommentIsPostedFromTheRecord(t *testing.T) {
	t.Parallel()
	// A world with markers, so the comment has the proposals that intake/markers.tsv holds (row 53).
	w, rows := markersToBatch(t, proposals())
	allDone(t, "the restart", rows, intakeSteps...)
	first := intakeComments(w)
	if len(first) != 1 || !strings.Contains(first[0], "\nproposed: Carol Jones\n") {
		t.Fatalf("the batch %q; want one comment with the proposal of the owner", first)
	}
	// The post never happened. The issue holds a comment of a person that
	// starts with the header and another comment of the App: neither is the
	// App's batch.
	w.fake.mu.Lock()
	w.fake.comments[1] = []map[string]any{
		{"id": 701, "user": map[string]any{"id": 303, "login": "carol"}, "created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z",
			"body": batchHeader + "\na copy by hand, with no answer line\n"},
		{"id": 702, "user": map[string]any{"id": 9001, "login": "layup-agent[bot]"}, "performed_via_github_app": map[string]any{"slug": "layup-agent"},
			"created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": "another comment of the App\n"},
	}
	w.fake.mu.Unlock()
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:questions.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if b := rowOf(again, "batch"); b.Result != "done" || b.Detail != "the comment of the batch posted from questions.tsv" {
		t.Errorf("the row %+v; want batch done, posted from the record", b)
	}
	if c := intakeComments(w); len(c) != 2 || len(first) != 1 || c[1] != first[0] {
		t.Errorf("the comment posted again %q; want the first's bytes %q", c, first)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:questions.tsv") != blob {
		t.Error("the restart changed questions.tsv")
	}
}

// addComment adds a comment of a person on issue n of the world's forge, at
// the time created.
func addComment(w *world, n int, id, author int, login, created, body string) {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	w.fake.comments[n] = append(w.fake.comments[n], map[string]any{"id": id, "user": map[string]any{"id": author, "login": login},
		"created_at": created, "updated_at": created, "body": body})
}

// The demo of row 47a (task T-hxb9, #223): after the batch, the step answers
// waits, the last row of a restart that releases the lease; after a comment of
// the idea owner with an answer line of each question (one corrected by a
// later comment) and one of a person who is no approver, the next restart
// copies the comments that are not the App's and writes intake/answers.tsv,
// the answered states and budget.tsv; a third restart finds the record.
func TestTheStepAnswers(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	if d := rowOf(rows, "answers").Detail; d != "waiting: no answer yet in the block idea-owner and the block operator" {
		t.Errorf("the detail of answers %q; want waiting in the two blocks (row 53: the Operator's block asks S01-stack)", d)
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("the lease is not released after a wait")
	}
	questions, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	asked, err := time.Parse(time.RFC3339, questions[0][9])
	if err != nil {
		t.Fatal(err)
	}
	at := func(d time.Duration) string { return asked.Add(d).UTC().Format(time.RFC3339) }
	var first, last, operator strings.Builder
	for _, q := range questions {
		answer := "the answer to " + q[0]
		switch q[2] {
		case "band.B":
			answer = "9.5"
		case "band.U":
			answer = "10.0"
		case "appetite.hours":
			answer = "40.0"
		}
		if q[3] == "operator" { // the Operator's block, S01-stack, answered by the Operator
			fmt.Fprintf(&operator, "%s: go\n", q[0])
			continue
		}
		fmt.Fprintf(&first, "%s: %s\n", q[0], answer)
	}
	fmt.Fprintf(&last, "%s: the answer corrected\n", questions[0][0])
	addComment(w, 1, 801, 999, "mallory", at(time.Minute), questions[0][0]+": not an approver\n")
	addComment(w, 1, 802, 303, "carol", at(2*time.Minute), first.String())
	addComment(w, 1, 803, 303, "carol", at(3*time.Minute), last.String())
	addComment(w, 1, 804, 101, "pharzam", at(4*time.Minute), operator.String())
	rows = Restart(context.Background(), w.config("0.1.0-dev"))
	// The review (row 47b) closes each answer, and the follow-up asks nothing.
	rows = beforeSetup(t, rows)
	allDone(t, "the restart with the answers", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision")
	if d := rowOf(rows, "follow-up").Detail; d != "nothing to ask again; Intake is done" {
		t.Errorf("the detail of follow-up %q", d)
	}
	if d := rowOf(rows, "follow-up-answers").Detail; d != "no follow-up; Intake is done" {
		t.Errorf("the detail of follow-up-answers %q (row 47c)", d)
	}
	if d := rowOf(rows, "answers").Detail; d != fmt.Sprintf("%d answers, in intake/answers.tsv; budget.tsv written", len(questions)) {
		t.Errorf("the detail of answers %q", d)
	}
	answers := show(t, w.bare, "refs/heads/layup-records", "intake/answers.tsv")
	if !strings.Contains(answers, "\n"+questions[0][0]+"\tthe answer corrected\t803\t303\t") || strings.Contains(answers, "not an approver") ||
		strings.Count(answers, "\n") != len(questions)+1 {
		t.Errorf("intake/answers.tsv:\n%s", answers)
	}
	// Each question has the time of its answer; the review then closed it.
	after, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil || len(after) != len(questions) {
		t.Fatalf("questions.tsv: %d rows, %v", len(after), err)
	}
	for _, q := range after {
		if q[11] == "" || q[12] != "closed" {
			t.Errorf("the question %q: want the time of its answer, and closed by the review", q)
		}
	}
	budget := show(t, w.bare, "refs/heads/layup-records", "budget.tsv")
	if !strings.Contains(budget, "\nband.B\t9.5\t") || !strings.Contains(budget, "\nband.U\t10.0\t") || !strings.Contains(budget, "\nappetite.hours\t40.0\t") {
		t.Errorf("budget.tsv:\n%s", budget)
	}
	copies := show(t, w.bare, "refs/heads/layup-records", "copies.tsv")
	for _, id := range []string{"801", "802", "803", "804"} {
		if !strings.Contains(copies, "\n"+id+"\t1\t1\t") {
			t.Errorf("copies.tsv has no copy of the comment %s:\n%s", id, copies)
		}
	}
	if strings.Contains(copies, "\tlayup-agent\t") {
		t.Errorf("copies.tsv holds the App's own comment:\n%s", copies)
	}
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if a := rowOf(again, "answers"); a.Result != "done" || a.Detail != "done before" {
		t.Errorf("the third restart's row %+v; want answers done before", a)
	}
}

// The demo of row 47b (task T-vwu3, #224): the idea owner answers all but the
// first question, with a B not below U, and the review leaves the hours open;
// one follow-up comment of the App on the Intake issue asks again exactly the
// missing question, the two unparsed fields and the open one; the next
// restart finds the records and posts nothing.
func TestTheReviewOfTheAnswersAndTheFollowUp(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers")
	questions, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	asked, _ := time.Parse(time.RFC3339, questions[0][9])
	var body, operator, review strings.Builder
	review.WriteString("question\tstate\treason\n")
	ref := map[string]string{} // the ID of each field of the intent form
	for _, q := range questions[1:] {
		answer := "the answer to " + q[0]
		switch q[2] {
		case "band.B":
			answer = "10.0" // not below U: both unparsed
		case "band.U":
			answer = "9.5"
		case "appetite.hours":
			answer = "40.0"
		}
		ref[q[2]] = q[0]
		if q[3] == "operator" { // S01-stack, answered by the Operator
			fmt.Fprintf(&operator, "%s: go\n", q[0])
		} else {
			fmt.Fprintf(&body, "%s: %s\n", q[0], answer)
		}
		state := "closed"
		if q[2] == "appetite.hours" {
			state = "open"
		}
		fmt.Fprintf(&review, "%s\t%s\tthe review's reason\n", q[0], state)
	}
	os.WriteFile(filepath.Join(w.dir, "bin", "answers-review.tsv"), []byte(review.String()), 0o644)
	addComment(w, 1, 802, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), body.String())
	addComment(w, 1, 803, 101, "pharzam", asked.Add(2*time.Minute).UTC().Format(time.RFC3339), operator.String())
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart with the answers", rows, "forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers")
	if d := rowOf(rows, "follow-up").Detail; d != "4 questions asked again, in questions.tsv and the comment of the follow-up" {
		t.Errorf("the detail of follow-up %q", d)
	}
	after, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	var again []string
	for _, q := range after[len(questions):] {
		if q[1] != "follow-up" {
			t.Errorf("a row after the batch %q is no follow-up", q)
		}
		again = append(again, q[2])
	}
	if want := []string{questions[0][0], ref["band.B"], ref["band.U"], ref["appetite.hours"]}; strings.Join(again, " ") != strings.Join(want, " ") {
		t.Errorf("the follow-up asks again %q; want the missing, the two unparsed and the open %q", again, want)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "intake/answers-review.tsv"); !strings.Contains(got, "\n"+ref["appetite.hours"]+"\topen\tthe review's reason\tS-") {
		t.Errorf("intake/answers-review.tsv:\n%s", got)
	}
	c := intakeComments(w)
	if len(c) != 2 || !strings.HasPrefix(c[1], followUpHeader+"\n") || !strings.Contains(c[1], "\n"+after[len(questions)][0]+": "+questions[0][8]+"\n") {
		t.Fatalf("the App's comments on the Intake issue %q; want the batch and one follow-up", c)
	}
	third := Restart(context.Background(), w.config("0.1.0-dev"))
	for _, name := range []string{"answers", "answers-review", "follow-up"} {
		if s := rowOf(third, name); s.Result != "done" || s.Detail != "done before" {
			t.Errorf("the third restart's row %+v; want %s done before", s, name)
		}
	}
	if s := rowOf(third, "follow-up-answers"); s.Detail != "waiting: no answer yet to the follow-up in the block idea-owner" {
		t.Errorf("the third restart's row %+v; want follow-up-answers waiting (row 47c)", s)
	}
	if len(intakeComments(w)) != 2 {
		t.Errorf("the third restart posted a comment: %d", len(intakeComments(w)))
	}
}

// intakeToFollowUp runs a world to the follow-up: the idea owner answers each
// question of the batch, the intent form with b, u and h, and the review
// leaves open the fields of the intent form named in open. It gives the world,
// the rows of questions.tsv after the follow-up, and the restart's rows.
func intakeToFollowUp(t *testing.T, b, u, h string, open ...string) (*world, [][]string, []Step) {
	t.Helper()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	Restart(context.Background(), w.config("0.1.0-dev"))
	questions, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	asked, _ := time.Parse(time.RFC3339, questions[0][9])
	var body, operator, review strings.Builder
	review.WriteString("question\tstate\treason\n")
	for _, q := range questions {
		answer, state := "the answer to "+q[0], "closed"
		switch q[2] {
		case "band.B":
			answer = b
		case "band.U":
			answer = u
		case "appetite.hours":
			answer = h
		}
		if slices.Contains(open, q[2]) {
			state = "open"
		}
		if q[3] == "operator" { // S01-stack, answered by the Operator
			fmt.Fprintf(&operator, "%s: go\n", q[0])
		} else {
			fmt.Fprintf(&body, "%s: %s\n", q[0], answer)
		}
		fmt.Fprintf(&review, "%s\t%s\tthe review's reason\n", q[0], state)
	}
	os.WriteFile(filepath.Join(w.dir, "bin", "answers-review.tsv"), []byte(review.String()), 0o644)
	addComment(w, 1, 802, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), body.String())
	addComment(w, 1, 801, 101, "pharzam", asked.Add(2*time.Minute).UTC().Format(time.RFC3339), operator.String())
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	after, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return w, after, rows
}

// The demo of row 47c (task T-tmm2, #225): budget.tsv is written at answers,
// the review leaves the hours open, the follow-up asks them again; the step
// follow-up-answers waits, then records the idea owner's answer, which
// replaces the hours in budget.tsv with the follow-up as its source, and
// Intake ends; a third restart finds the record.
func TestTheAnswersToTheFollowUpReplaceAValueOfTheBudget(t *testing.T) {
	t.Parallel()
	w, questions, rows := intakeToFollowUp(t, "5.0", "10.0", "40.0", "appetite.hours")
	if d := rowOf(rows, "follow-up-answers").Detail; d != "waiting: no answer yet to the follow-up in the block idea-owner" {
		t.Fatalf("the detail of follow-up-answers %q; want waiting", d)
	}
	follow := questions[len(questions)-1]
	if follow[1] != "follow-up" {
		t.Fatalf("the last question %q is no follow-up", follow)
	}
	asked, _ := time.Parse(time.RFC3339, follow[9])
	addComment(w, 1, 803, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), follow[0]+": 50.0\n")
	rows = Restart(context.Background(), w.config("0.1.0-dev"))
	if d := rowOf(rows, "follow-up-answers").Detail; d != "1 answers to the follow-up, in intake/answers.tsv; budget.tsv replaced; Intake is done" {
		t.Errorf("the detail of follow-up-answers %q", d)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "budget.tsv"); !strings.Contains(got, "\nappetite.hours\t50.0\t"+follow[0]+"\n") || !strings.Contains(got, "\nband.B\t5.0\t") {
		t.Errorf("budget.tsv:\n%s", got)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "intake/answers.tsv"); !strings.Contains(got, "\n"+follow[0]+"\t50.0\t803\t303\t") {
		t.Errorf("intake/answers.tsv:\n%s", got)
	}
	after, _ := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	for _, q := range after {
		switch {
		case q[0] == follow[0] && q[12] != "answered":
			t.Errorf("the follow-up %q: want answered", q)
		case q[0] == follow[2] && q[12] != "open":
			t.Errorf("the question asked again %q: want its own state, open", q)
		}
	}
	if strings.Contains(show(t, w.bare, "refs/heads/layup-records", "copies.tsv"), "\tlayup-agent\t") {
		t.Error("copies.tsv holds a comment of the App")
	}
	third := Restart(context.Background(), w.config("0.1.0-dev"))
	if s := rowOf(third, "follow-up-answers"); s.Result != "done" || s.Detail != "done before" {
		t.Errorf("the third restart's row %+v; want follow-up-answers done before", s)
	}
}

// With the batch's B not below U, budget.tsv is first written from the
// answers to the follow-up, each value with the question that gives it.
func TestTheAnswersToTheFollowUpWriteTheBudget(t *testing.T) {
	t.Parallel()
	w, questions, _ := intakeToFollowUp(t, "10.0", "9.5", "40.0")
	var follow [][]string
	for _, q := range questions {
		if q[1] == "follow-up" {
			follow = append(follow, q)
		}
	}
	if len(follow) != 2 {
		t.Fatalf("the follow-up rows %q; want B and U asked again", follow)
	}
	asked, _ := time.Parse(time.RFC3339, follow[0][9])
	addComment(w, 1, 803, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), follow[0][0]+": 5.0\n"+follow[1][0]+": 10.0\n")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	if d := rowOf(rows, "follow-up-answers").Detail; d != "2 answers to the follow-up, in intake/answers.tsv; budget.tsv written; Intake is done" {
		t.Errorf("the detail of follow-up-answers %q", d)
	}
	got := show(t, w.bare, "refs/heads/layup-records", "budget.tsv")
	if !strings.Contains(got, "\nband.B\t5.0\t"+follow[0][0]+"\n") || !strings.Contains(got, "\nband.U\t10.0\t"+follow[1][0]+"\n") || !strings.Contains(got, "\nappetite.hours\t40.0\t") {
		t.Errorf("budget.tsv:\n%s", got)
	}
}

// Input states: an answer to the follow-up that still does not parse leaves
// budget.tsv unwritten and its question of the intent form open; Intake is
// done.
func TestAnUnparsedAnswerToTheFollowUpLeavesItsQuestionOpen(t *testing.T) {
	t.Parallel()
	w, questions, _ := intakeToFollowUp(t, "10.0", "9.5", "40.0")
	var follow [][]string
	for _, q := range questions {
		if q[1] == "follow-up" {
			follow = append(follow, q)
		}
	}
	asked, _ := time.Parse(time.RFC3339, follow[0][9])
	addComment(w, 1, 803, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), follow[0][0]+": five\n"+follow[1][0]+": 10.0\n")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	if d := rowOf(rows, "follow-up-answers").Detail; d != "2 answers to the follow-up, in intake/answers.tsv; budget.tsv not written; Intake is done" {
		t.Errorf("the detail of follow-up-answers %q", d)
	}
	if strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "budget.tsv") {
		t.Error("budget.tsv is written")
	}
	after, _ := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	for _, q := range after {
		if q[0] == follow[0][2] && q[12] != "open" {
			t.Errorf("the question of B %q: want open", q)
		}
	}
}
