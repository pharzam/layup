//go:build integration

package run

import (
	"context"
	"errors"
	"os"
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
		if after := len(d.sessionsOf(spec)); after != 4 {
			t.Errorf("%s: %d sessions in all; want 4, one for each stage", c.name, after)
		}
		if n := countPosts(d, "findings:"); n != 1 {
			t.Errorf("%s: %d review records on the issue; want one", c.name, n)
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
