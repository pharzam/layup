//go:build integration

package run

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
)

// The driver of a Shape task through the Sessions of state.sessions(), the httptest forge, real git and scripted harnesses
// (task T-jy6f, #415; loop.md, The driver of a Shape task).

// advance calls the driver on Sessions that the production constructor made; tweak injects a fault into them.
func (d *drvWorld) advance(task string, tweak func(*Sessions)) (Advance, error) {
	d.t.Helper()
	s := d.r.sessions()
	if tweak != nil {
		tweak(s)
	}
	return s.AdvanceShapeTask(context.Background(), d.input(task))
}

func (d *drvWorld) mustVerify(task string) {
	d.t.Helper()
	got, err := d.advance(task, nil)
	if err != nil || !got.Verified || got.Task != task {
		d.t.Fatalf("AdvanceShapeTask(%s) = %+v, %v; want the verified head", task, got, err)
	}
}

// sessionsOf gives the rows of sessions.tsv of the task.
func (d *drvWorld) sessionsOf(task string) [][]string {
	d.t.Helper()
	rows, err := records.ReadSessions([]byte(keptRecord(d.t, d.w, "sessions.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	return slices.DeleteFunc(rows, func(r []string) bool { return r[sessionTask] != task })
}

func (d *drvWorld) boundOf(task string) []string {
	var out []string
	for _, e := range d.eventsOf(task) {
		if e[1] == "bound" {
			out = append(out, e[5])
		}
	}
	return out
}

func (d *drvWorld) promptsOf(task string) []string {
	files, _ := filepath.Glob(filepath.Join(d.ctl, "prompt-"+task+"-*"))
	slices.Sort(files)
	var out []string
	for _, f := range files {
		b, _ := os.ReadFile(f)
		out = append(out, string(b))
	}
	return out
}

var oneRoundEvents = []string{"attempt", "session", "result", "handoff", "session", "result", "handoff", "session", "result", "push", "bound", "handoff", "session", "result", "handoff"}

// I4, the demo: each Shape task goes from its opened issue and row through the plan, the plan review, the developer's push, the
// pull request, the four statuses and the verifier to a verified head. The sessions that only read commit, and are bound to
// nothing; a draft that changed after the opening changes no registered ID; a second call writes nothing.
func TestTheShapeTasksGoThroughTheLoopToTheirVerifiedHeads(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("commit", "")
	spec, arch := d.open()
	pushRecords(t, d.w.bare, map[string]string{"spec/draft.tsv": draftTSV + "R-002\tnon-functional\tFast\t1\tc\tS-aaaaaaaa\n"}, "a later draft")
	d.r = restarted(t, d.w)
	main := strings.TrimSpace(gitOut(t, d.r.target, "rev-parse", "refs/heads/main"))
	for _, task := range []string{spec, arch} {
		d.mustVerify(task)
	}
	for _, task := range []string{spec, arch} {
		if got := d.kinds(task); !slices.Equal(got, oneRoundEvents) {
			t.Errorf("the events of %s = %v; want %v", task, got, oneRoundEvents)
		}
		heads := d.boundOf(task)
		if len(heads) != 1 {
			t.Fatalf("%s is bound %d times; want once, the developer's", task, len(heads))
		}
		var got []string
		for _, s := range statusesOf(d.w, heads[0]) {
			got = append(got, s["context"].(string)+" "+s["state"].(string))
		}
		for _, want := range []string{GatesContext + " success", SpecContext + " success", RulesContext + " success", VerifyContext + " success"} {
			if !slices.Contains(got, want) {
				t.Errorf("the statuses of the head of %s = %v; want %q", task, got, want)
			}
		}
		if got[0] != VerifyContext+" success" {
			t.Errorf("the newest status of the head of %s is %q; want the verifier's success after its pending one", task, got[0])
		}
		roles := map[string]int{}
		for _, r := range d.sessionsOf(task) {
			roles[r[sessionRole]]++
		}
		if roles[stepRole("plan-of-a-task")] != 1 || roles[stepRole("verification")] != 2 {
			t.Errorf("the sessions of %s by role = %v; want one plan session and the plan review and the verifier of one role", task, roles)
		}
	}
	if len(d.w.fake.pulls.list) != 2 || !strings.Contains(d.w.fake.pulls.list[0].body, "The review record: ") || !strings.Contains(d.w.fake.pulls.list[1].body, "The review record: ") {
		t.Errorf("the pull requests: %d; want two, each with the link of its review record", len(d.w.fake.pulls.list))
	}
	if got := strings.TrimSpace(gitOut(t, d.r.target, "rev-parse", "refs/heads/main")); got != main {
		t.Errorf("the default branch of the run's clone moved from %s to %s; want it never moved (plan amendment 1)", main, got)
	}
	// A second call writes nothing.
	posts, patches, n := d.requests("POST "), d.requests("PATCH "), len(d.eventsOf(spec))+len(d.eventsOf(arch))
	for _, task := range []string{spec, arch} {
		d.mustVerify(task)
	}
	if d.requests("POST ") != posts || d.requests("PATCH ") != patches || len(d.eventsOf(spec))+len(d.eventsOf(arch)) != n {
		t.Errorf("a second call wrote: %d POST, %d PATCH, %d events; want %d, %d, %d", d.requests("POST "), d.requests("PATCH "), len(d.eventsOf(spec))+len(d.eventsOf(arch)), posts, patches, n)
	}
}

// A gate failure starts a new attempt from the same base with the table's findings, counts no review round, and is
// followed to a verified head.
func TestAGateFailureOfAShapeTaskStartsANewAttempt(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("gofile", "")
	spec, _ := d.open()
	d.mustVerify(spec)
	got := d.kinds(spec)
	if n := slices.Index(got, "attempt"); n != 0 || strings.Count(strings.Join(got, " "), "attempt") != 2 {
		t.Errorf("the events %v; want two attempts", got)
	}
	prompts := d.promptsOf(spec)
	// The prompts are those of the plan, the plan review, the developer of attempt 1, the developer of attempt 2 and the verifier.
	if len(prompts) != 5 || !strings.Contains(prompts[3], "gate:fine") || strings.Contains(prompts[2], "gate:fine") {
		t.Errorf("%d prompts, the developer's prompt of attempt 2 holds the finding: %v; want the finding of the gate in it and not in attempt 1's", len(prompts), len(prompts) == 5 && strings.Contains(prompts[3], "gate:fine"))
	}
	bound := d.boundOf(spec)
	if len(bound) != 2 || bound[0] == bound[1] {
		t.Fatalf("the bound heads %v; want two", bound)
	}
	if st := statusesOf(d.w, bound[0]); len(st) == 0 || !slices.ContainsFunc(st, func(m map[string]any) bool { return m["context"] == GatesContext && m["state"] == "failure" }) {
		t.Errorf("the statuses of the first head: %v; want a failure of layup/gates", st)
	}
	if n := len(slices.DeleteFunc(d.sessionsOf(spec), func(r []string) bool { return r[sessionAttempt] != "2" })); n != 2 {
		t.Errorf("%d sessions in attempt 2; want its developer and its verifier", n)
	}
}

// One material round starts the next attempt with the finding in the prompt, the second waits at the cap of 1; each failed
// verdict is posted at its reviewed head before the attempt or the wait; a second call waits again, with no session.
func TestAMaterialRoundStartsOneAttemptAndTheSecondWaitsAtTheCap(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	d.ctlWrite("verdicts-"+spec, "material\nmaterial\n")
	got, err := d.advance(spec, nil)
	want := "waiting for the Operator: " + spec + ": the review-cycle cap 1 is reached with findings open"
	if err != nil || got.Wait != want || got.Verified || got.Stop != "" {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the wait %q", got, err, want)
	}
	heads := d.boundOf(spec)
	if len(heads) != 2 {
		t.Fatalf("the bound heads %v; want two attempts", heads)
	}
	for _, h := range heads {
		if st := statusesOf(d.w, h); len(st) == 0 || st[0]["context"] != VerifyContext || st[0]["state"] != "failure" || !strings.Contains(st[0]["description"].(string), "material: 1 findings") {
			t.Errorf("the newest status of %s = %v; want the failure of layup/verify (material: 1 findings)", h[:7], st)
		}
	}
	// The prompts: plan, plan review, developer 1, verifier 1, developer 2, verifier 2.
	if p := d.promptsOf(spec); len(p) != 6 || !strings.Contains(p[4], "spec docs/spec/alpha.md:1: the section is thin") || strings.Contains(p[2], "the section is thin") {
		t.Errorf("the prompt of attempt 2 does not hold the finding of the record")
	}
	n := len(d.eventsOf(spec))
	if again, err := d.advance(spec, nil); err != nil || again.Wait != want || len(d.eventsOf(spec)) != n {
		t.Errorf("a second call = %+v, %v with %d events; want the same wait and no event", again, err, len(d.eventsOf(spec))-n)
	}
}

func TestNotMergeableStartsNoAttempt(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	d.ctlWrite("verdicts-"+spec, "not mergeable, findings recorded\n")
	got, err := d.advance(spec, nil)
	if err != nil || got.Wait != "waiting for the Operator: "+spec+": not mergeable, findings recorded" {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the wait for the Operator", got, err)
	}
	if n := strings.Count(strings.Join(d.kinds(spec), " "), "attempt"); n != 1 {
		t.Errorf("%d attempts; want one: not mergeable is never another cycle", n)
	}
	if st := statusesOf(d.w, d.boundOf(spec)[0]); st[0]["state"] != "failure" {
		t.Errorf("the newest status = %v; want the failure of layup/verify before the wait", st[0])
	}
}

// loop.md, layup/verify: with no harness outside the authors, the status not-active is posted before the stop, once, and no
// verifier starts.
func TestAVerifierWithNoIndependentHarnessPostsNotActiveBeforeTheStop(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	pushRecords(t, d.w.bare, map[string]string{"routing.tsv": "role\ttier\tposition\tharness\tmodel\nsoftware-engineer\treasoning\t1\tclaude\tm1\nproduct-owner\treasoning\t1\tdevin\tm1\nsystems-architect\treasoning\t1\tdevin\tm1\nqa-engineer\treasoning\t1\tdevin\tm1\n"}, "routing")
	d.r = restarted(t, d.w)
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || got.Stop != spec+": not-active: no harness outside the authors" {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop for not-active", got, err)
	}
	h := d.boundOf(spec)[0]
	failures := 0
	for _, st := range statusesOf(d.w, h) {
		if st["context"] == VerifyContext && st["state"] == "failure" && st["description"] == "not-active: no harness outside the authors" {
			failures++
		}
	}
	if failures != 1 {
		t.Errorf("%d posts of the not-active failure at the head; want one", failures)
	}
	if again, err := d.advance(spec, nil); err != nil || again.Stop == "" {
		t.Errorf("a second call = %+v, %v; want the same stop", again, err)
	}
	if n := len(slices.DeleteFunc(d.sessionsOf(spec), func(r []string) bool { return r[sessionRole] != stepRole("verification") })); n != 1 {
		t.Errorf("%d qa sessions; want only the plan review, no verifier", n)
	}
	if n := len(slices.DeleteFunc(statusesOf(d.w, h), func(m map[string]any) bool { return m["description"] != "not-active: no harness outside the authors" })); n != 1 {
		t.Errorf("%d not-active posts after a second call; want one", n)
	}
}

// Review conditions 3 and 6: a stage whose session ended and whose post, handoff or body edit was interrupted is finished from
// the kept files, and never run again.
func TestAnInterruptedStageIsFinishedFromItsKeptFilesAndRunsNoSecondSession(t *testing.T) {
	for _, c := range []struct {
		name  string
		setup func(d *drvWorld) func(*Sessions)
	}{
		{"the post of the plan", func(d *drvWorld) func(*Sessions) {
			d.w.fake.failNote = "size: small"
			return nil
		}},
		{"the post of the plan review", func(d *drvWorld) func(*Sessions) {
			d.w.fake.failNote = "verdict: approve"
			return nil
		}},
		{"the post of the review record", func(d *drvWorld) func(*Sessions) {
			d.w.fake.failNote = "findings:"
			return nil
		}},
		{"the handoff of the verifier", func(d *drvWorld) func(*Sessions) {
			return func(s *Sessions) {
				s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
					if strings.Contains(m, "the handoff verifier to close-out") {
						return errors.New("the run stopped")
					}
					return nil
				}}
			}
		}},
		{"the edit of the body", func(d *drvWorld) func(*Sessions) {
			return func(s *Sessions) { s.Forge = &editFault{Forge: s.Forge} }
		}},
	} {
		d := newDrvWorld(t, loopPerms)
		spec, _ := d.open()
		tweak := c.setup(d)
		if _, err := d.advance(spec, tweak); err == nil {
			t.Errorf("%s: the first call = nil; want the interruption", c.name)
			continue
		}
		d.mustVerify(spec)
		roles := map[string]int{}
		for _, r := range d.sessionsOf(spec) {
			roles[r[sessionRole]]++
		}
		if roles[stepRole("plan-of-a-task")] != 1 || roles[stepRole("verification")] != 2 || roles[stepRole("specification")] != 1 || len(d.sessionsOf(spec)) != 4 {
			t.Errorf("%s: the sessions by role %v; want one plan, one plan review and one verifier (qa-engineer 2), one developer and no second of a stage", c.name, roles)
		}
		for _, sub := range []string{"size: small", "verdict: approve", "findings:"} {
			if n := countPosts(d, sub); n != 1 {
				t.Errorf("%s: %d comments hold %q; want one (the plan, the plan review and the review record are each posted once)", c.name, n, sub)
			}
		}
	}
}

func countPosts(d *drvWorld, sub string) int {
	d.w.fake.mu.Lock()
	defer d.w.fake.mu.Unlock()
	n := 0
	for _, p := range d.w.fake.posted {
		if strings.Contains(p, sub) {
			n++
		}
	}
	return n
}

// faultStore fails a records commit whose message fail names.
type faultStore struct {
	SessionStore
	fail    func(message string) error
	started time.Time
}

func (f *faultStore) Commit(ctx context.Context, files map[string][]byte, message string) error {
	if err := f.fail(message); err != nil {
		f.fail = func(string) error { return nil } // once
		return err
	}
	return f.SessionStore.Commit(ctx, files, message)
}

func (f *faultStore) Started(ctx context.Context) (time.Time, error) {
	if !f.started.IsZero() {
		return f.started, nil
	}
	return f.SessionStore.Started(ctx)
}

// editFault fails the first edit of a body, after nothing is sent.
type editFault struct {
	forge.Forge
	failed bool
}

func (e *editFault) EditBody(ctx context.Context, n int, body string) (forge.PullRequest, error) {
	if !e.failed {
		e.failed = true
		return forge.PullRequest{}, errors.New("the forge did not answer")
	}
	return e.Forge.EditBody(ctx, n, body)
}

// A specification that lacks the section of a registered ID has an invalid handoff, opens no pull request and stops.
func TestAMissingSectionMakesTheHandoffInvalidAndStops(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("no-section", "")
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || !strings.Contains(got.Stop, "developer (specification) → verifier is invalid spec-sections") {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop for the invalid spec-sections", got, err)
	}
	if n := len(d.w.fake.pulls.list); n != 0 {
		t.Errorf("%d pull requests; want none", n)
	}
}

// A push that the developer never made (no commit) stops; a head that prd-lint does not accept as a waiting specification stops
// before the verifier, naming the line.
func TestASpecificationThatPrdLintDoesNotAcceptStopsBeforeTheVerifier(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("bad-prd", "")
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || !strings.Contains(got.Stop, spec+": the pinned prd-lint does not accept the head as a specification that waits for its bet: ") {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop that names prd-lint", got, err)
	}
	if n := len(slices.DeleteFunc(d.sessionsOf(spec), func(r []string) bool { return r[sessionRole] != stepRole("verification") })); n != 1 {
		t.Errorf("%d qa sessions; want the plan review only, no verifier", n)
	}
}

func TestADeveloperThatCommitsNothingMakesAnInvalidHandoffAndStops(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("no-commit", "")
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || !strings.Contains(got.Stop, "developer (specification) → verifier is invalid") {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop for the invalid handoff", got, err)
	}
}

// A restart after the row starts the first attempt once: a call that stops for want of a pair leaves the attempt, and the next
// call goes on from it.
func TestARestartAfterTheRowStartsTheFirstAttemptOnce(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	in := d.input(spec)
	in.Harnesses = nil
	got, err := d.r.sessions().AdvanceShapeTask(context.Background(), in)
	if err != nil || !strings.Contains(got.Stop, spec+": no admitted pair for the step plan-of-a-task") {
		t.Fatalf("AdvanceShapeTask with no harness = %+v, %v; want the stop for the pair", got, err)
	}
	d.r = restarted(t, d.w)
	d.mustVerify(spec)
	if n := strings.Count(strings.Join(d.kinds(spec), " "), "attempt"); n != 1 {
		t.Errorf("%d attempts; want one", n)
	}
}

// The cap of milestone 0 is asked before each session starts.
func TestTheCapOfMilestoneZeroIsAskedBeforeEachSession(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	_, err := d.advance(spec, func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(string) error { return nil }, started: s.Now().Add(-100 * time.Hour)}
	})
	if !errors.Is(err, ledger.ErrMilestoneCap) {
		t.Fatalf("AdvanceShapeTask with milestone 0 over its time = %v; want ledger.ErrMilestoneCap", err)
	}
	if gitOut(t, d.w.bare, "ls-tree", "--name-only", "layup-records", "sessions.tsv") != "" {
		t.Error("a session started; want none")
	}
}

// statusFault fails the post of the failure status of layup/verify that is the failOn-th, once.
type statusFault struct {
	forge.Forge
	failOn, n int
}

func (f *statusFault) SetStatus(ctx context.Context, sha, name, state, description string) error {
	if name == VerifyContext && state == "failure" {
		if f.n++; f.n == f.failOn {
			return errors.New("the forge did not take the status")
		}
	}
	return f.Forge.SetStatus(ctx, sha, name, state, description)
}

// Review condition 4: when the post of the failure of a material review fails, at the end of an attempt or before the terminal wait, the
// next run posts it before the next attempt starts or the wait is given, and adds no session and no comment.
func TestAFailedPostOfTheVerdictStatusIsPostedBeforeTheAttemptOrTheWait(t *testing.T) {
	for _, c := range []struct {
		name     string
		failOn   int
		attempts int // the attempts when the first call stops
	}{{"the end of the first attempt", 1, 1}, {"the terminal wait", 2, 2}} {
		d := newDrvWorld(t, loopPerms)
		spec, _ := d.open()
		d.ctlWrite("verdicts-"+spec, "material\nmaterial\n")
		if _, err := d.advance(spec, func(s *Sessions) { s.Forge = &statusFault{Forge: s.Forge, failOn: c.failOn} }); err == nil {
			t.Fatalf("%s: the first call = nil; want the failed post", c.name)
		}
		if n := strings.Count(strings.Join(d.kinds(spec), " "), "attempt"); n != c.attempts {
			t.Errorf("%s: %d attempts after the failed post; want %d (no attempt starts before the failure is posted)", c.name, n, c.attempts)
		}
		heads := d.boundOf(spec)
		if st := statusesOf(d.w, heads[len(heads)-1]); st[0]["state"] != "pending" {
			t.Errorf("%s: the newest status of the reviewed head is %v; want it still pending, the failure not posted", c.name, st[0])
		}
		sessions, comments := len(d.sessionsOf(spec)), countPosts(d, "findings:")
		got, err := d.advance(spec, nil)
		if err != nil || !strings.Contains(got.Wait, "the review-cycle cap 1 is reached") {
			t.Fatalf("%s: the second call = %+v, %v; want the wait at the cap", c.name, got, err)
		}
		heads = d.boundOf(spec)
		for _, h := range heads {
			if st := statusesOf(d.w, h); st[0]["state"] != "failure" || st[0]["context"] != VerifyContext {
				t.Errorf("%s: the newest status of %s is %v; want the failure of layup/verify", c.name, h[:7], st[0])
			}
		}
		if c.failOn == 2 && (len(d.sessionsOf(spec)) != sessions || countPosts(d, "findings:") != comments) {
			t.Errorf("%s: a session or a comment was added by the call that recovered", c.name)
		}
		if c.failOn == 1 && (len(d.sessionsOf(spec)) != 6 || countPosts(d, "findings:") != 2 || len(heads) != 2) {
			t.Errorf("%s: %d sessions, %d review records, %d heads; want 6, 2 and 2", c.name, len(d.sessionsOf(spec)), countPosts(d, "findings:"), len(heads))
		}
	}
}

// loop.md, layup/verify: with a gate failure and no harness outside the authors, the status not-active is posted and the driver stops; it
// starts no second attempt, as the failure of a gate alone would.
func TestAGateFailureWithNoIndependentVerifierStopsAndStartsNoAttempt(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	pushRecords(t, d.w.bare, map[string]string{"routing.tsv": "role\ttier\tposition\tharness\tmodel\nsoftware-engineer\treasoning\t1\tclaude\tm1\nproduct-owner\treasoning\t1\tdevin\tm1\nsystems-architect\treasoning\t1\tdevin\tm1\nqa-engineer\treasoning\t1\tdevin\tm1\n"}, "routing")
	d.r = restarted(t, d.w)
	d.ctlWrite("gofile", "")
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || got.Stop != spec+": not-active: no harness outside the authors" {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop for not-active", got, err)
	}
	if n := strings.Count(strings.Join(d.kinds(spec), " "), "attempt"); n != 1 {
		t.Errorf("%d attempts; want one: the gate failure starts none when the verifier cannot run", n)
	}
	var got2 []string
	for _, st := range statusesOf(d.w, d.boundOf(spec)[0]) {
		got2 = append(got2, st["context"].(string)+" "+st["state"].(string))
	}
	for _, want := range []string{GatesContext + " failure", VerifyContext + " failure"} {
		if !slices.Contains(got2, want) {
			t.Errorf("the statuses of the head = %v; want %q", got2, want)
		}
	}
}

// registerAfterBound answers a read of sessions.tsv, once the task has a bound, with the records commit of the developer's sessions set to rev: the
// head is then judged by the register of that commit, which a push check made at the start of the session could not have seen (a rule-path change
// cannot be pushed, so layup/rules cannot fail on a head of the loop's own making; this is the seeded state of the review of #415, finding 2).
type registerAfterBound struct {
	SessionStore
	task, rev string
}

func (r *registerAfterBound) ReadFile(ctx context.Context, path string) ([]byte, error) {
	data, err := r.SessionStore.ReadFile(ctx, path)
	if path != "sessions.tsv" || err != nil || data == nil {
		return data, err
	}
	if ev, _ := r.SessionStore.ReadFile(ctx, eventsPath(r.task)); !strings.Contains(string(ev), "\tbound\t") {
		return data, nil
	}
	rows, err := records.ReadSessions(data)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row[sessionRole] == stepRole("specification") {
			row[sessionRecords] = r.rev
		}
	}
	return table(records.SessionsSchema, rows)
}

// sideCommit makes a commit in the clone, on no ref, of a tree that holds the files, and gives its SHA.
func sideCommit(t *testing.T, clone string, files map[string]string) string {
	t.Helper()
	run := func(stdin string, env []string, args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir, cmd.Stdin = clone, strings.NewReader(stdin)
		cmd.Env = append(os.Environ(), append(env, "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@x.invalid", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@x.invalid")...)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("git %v: %v", args, err)
		}
		return strings.TrimSpace(string(out))
	}
	index := filepath.Join(t.TempDir(), "index")
	for p, text := range files {
		blob := run(text, nil, "hash-object", "-w", "--stdin")
		run("", []string{"GIT_INDEX_FILE=" + index}, "update-index", "--add", "--cacheinfo", "100644,"+blob+","+p)
	}
	tree := run("", []string{"GIT_INDEX_FILE=" + index}, "write-tree")
	return run("a side commit\n", nil, "commit-tree", tree)
}

// plan.md:147 and loop.md: a failure of layup/spec or layup/rules, and the status error of the gates, each stop the step; no verifier starts and
// no attempt follows. The control is the same world with no fault, which reaches the verified head.
func TestAFailedStatusOfTheHeadStopsTheStepWithNoVerifierAndNoNewAttempt(t *testing.T) {
	for _, c := range []struct {
		name, want string
		gates      string
		prep       func(d *drvWorld)
		tweak      func(d *drvWorld, task string) func(*Sessions)
	}{
		{"no fault (control)", "", gatesTSV, nil, nil},
		{"a failed layup/spec", ": layup/spec failed: criteria: ", gatesTSV, func(d *drvWorld) { d.ctlWrite("spec-fail", "") }, nil},
		{"a failed layup/rules", ": layup/rules failed: ", gatesTSV, nil, func(d *drvWorld, task string) func(*Sessions) {
			rev := sideCommit(t, d.r.target, map[string]string{"rule-paths.tsv": register + "docs/tasks/\t—\tarchitecture\n"})
			return func(s *Sessions) { s.Store = &registerAfterBound{SessionStore: s.Store, task: task, rev: rev} }
		}},
		{"a gate status of error", ": layup/gates is error: ", "kind\tstate\ttool\tcommand\tscope\tconfig\n", nil, nil},
	} {
		d := newDrvWorldWith(t, loopPerms, c.gates)
		if c.prep != nil {
			c.prep(d)
		}
		spec, _ := d.open()
		var tweak func(*Sessions)
		if c.tweak != nil {
			tweak = c.tweak(d, spec)
		}
		got, err := d.advance(spec, tweak)
		if c.want == "" {
			if err != nil || !got.Verified {
				t.Errorf("%s: %+v, %v; want the verified head", c.name, got, err)
			}
			continue
		}
		if err != nil || !strings.HasPrefix(got.Stop, spec) || !strings.Contains(got.Stop, c.want) {
			t.Errorf("%s: %+v, %v; want a stop holding %q", c.name, got, err, c.want)
		}
		if n := strings.Count(strings.Join(d.kinds(spec), " "), "attempt"); n != 1 {
			t.Errorf("%s: %d attempts; want one: a stop starts no attempt", c.name, n)
		}
		if n := len(slices.DeleteFunc(d.sessionsOf(spec), func(r []string) bool { return r[sessionRole] != stepRole("verification") })); n != 1 {
			t.Errorf("%s: %d qa sessions; want the plan review only, no verifier", c.name, n)
		}
	}
}

// I8 (plan amendment 1 of #364, #432): the architecture task starts at the specification task's verified head. On a target with no
// PRD its head passes the full check of layup/spec (the PRD is in its base); with a PRD in the first tree the same run passes from
// either base, as before. The driver gives the verified head, and every session of the architecture task, the plan and the
// developer included, clones the scratch branch of that base.
func TestTheArchitectureHeadOfATargetWithNoPRDPassesLayupSpec(t *testing.T) {
	for _, noPRD := range []bool{false, true} {
		d := newDrvWorldOpts(t, loopPerms, drvOpts{noPRD: noPRD})
		d.ctlWrite("commit", "")
		spec, arch := d.open()
		first, err := d.advance(spec, nil)
		if err != nil || !first.Verified || (noPRD && first.Head != d.boundOf(spec)[0]) {
			t.Fatalf("noPRD %v: the specification: %+v, %v; want a verified head that is its bound head", noPRD, first, err)
		}
		in := d.input(arch)
		if noPRD {
			in.Base = first.Head
		}
		got, err := d.r.sessions().AdvanceShapeTask(context.Background(), in)
		if err != nil || !got.Verified || got.Stop != "" {
			t.Fatalf("noPRD %v: the architecture task: %+v, %v; want the verified head", noPRD, got, err)
		}
		if !noPRD {
			continue
		}
		for _, e := range d.eventsOf(arch) {
			if e[1] == "attempt" && e[4] != first.Head {
				t.Errorf("the attempt %s of the architecture task starts at %s; want the head %s of the specification", e[2], e[4], first.Head)
			}
		}
		head := d.boundOf(arch)[0]
		table := keptRecord(t, d.w, "tasks/"+arch+"/spec-check.tsv")
		for _, rule := range []string{"covers", "ran"} {
			if !strings.Contains(table, head+"\t"+rule+"\tpass\t") {
				t.Errorf("the rule %s of the head of the architecture task is not a pass:\n%s", rule, table)
			}
		}
		var states []string
		for _, s := range statusesOf(d.w, head) {
			states = append(states, s["context"].(string)+" "+s["state"].(string))
		}
		for _, want := range []string{SpecContext + " success", RulesContext + " success"} {
			if !slices.Contains(states, want) {
				t.Errorf("the statuses of the head of the architecture task = %v; want %q", states, want)
			}
		}
		if files := gitOut(t, d.r.target, "diff", "--name-only", first.Head, head); strings.Contains(files, "docs/prd/") || !strings.Contains(files, "docs/spec/architecture.md") {
			t.Errorf("the files of the architecture task from the specification's head = %q; want its own changes, none of docs/prd/", files)
		}
	}
}
