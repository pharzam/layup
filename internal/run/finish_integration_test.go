//go:build integration

package run

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

// The close-out and the merge of a verified task on the production path (task T-ehb4, #425, row 66y; loop.md, The merge at the head SHA and
// the close-out commit, D1 to D10): the Shape world of the driver, state.sessions(), the httptest forge, real git and the scripted harnesses,
// with a route for the close-out and a required job of GitHub Actions. A second call and a restart are new Sessions on the same records.

// A finishing is the world of FinishTask: the architecture task verified by the driver, and the route of the close-out.
type finishing struct {
	*drvWorld
	spec, arch string
}

// newFinishing makes the world; reasoning and execution are the harnesses of the two routes of the close-out.
func newFinishing(t *testing.T, reasoning, execution string) *finishing {
	t.Helper()
	d := newDrvWorld(t, loopPerms)
	routing := keptRecord(t, d.w, "routing.tsv") + "software-developer\treasoning\t1\t" + reasoning + "\tm1\nsoftware-developer\texecution\t1\t" + execution + "\tm1\n"
	pushRecords(t, d.w.bare, map[string]string{"routing.tsv": routing}, "the route of the close-out")
	d.restart()
	spec, arch := d.open()
	d.mustVerify(arch)
	f := &finishing{d, spec, arch}
	f.job("completed", "success")
	return f
}

// job sets the run of the required job build at every head, and the rules that require it.
func (f *finishing) job(status, conclusion string) {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	f.w.fake.rules = map[string]string{"main": requiredBuild}
	f.w.fake.checkRuns = map[string][]map[string]any{"*": {actionsRun(status, conclusion, 15368)}}
	f.w.fake.pulls.follow = true
}

func (f *finishing) finishTask(task string, harnesses [][]string, tweak func(*Sessions)) (Advance, error) {
	f.t.Helper()
	s := f.r.sessions()
	if tweak != nil {
		tweak(s)
	}
	in := FinishInput{Task: task, Records: f.r.store.Base, Harnesses: harnesses, Models: f.models}
	in.Copy = f.r.betCopy
	if task == f.spec {
		in.Bet = f.r.renderedBet
	}
	return s.FinishTask(context.Background(), in)
}

func (f *finishing) finish(tweak func(*Sessions)) (Advance, error) {
	return f.finishTask(f.arch, f.harnesses, tweak)
}

// closers gives the rows of sessions.tsv of the close-out sessions of the task.
func (f *finishing) closers(task string) [][]string {
	return slices.DeleteFunc(f.sessionsOf(task), func(r []string) bool { return r[sessionRole] != stepRole("close-out") })
}

// closedHead gives the bound head of the last close-out session of the task, "" for none.
func (f *finishing) closedHead(task string) string {
	head := ""
	for _, c := range f.closers(task) {
		for _, e := range f.eventsOf(task) {
			if e[1] == "bound" && e[3] == c[sessionID] {
				head = e[5]
			}
		}
	}
	return head
}

// writes counts the requests of the forge that change something.
func (f *finishing) writes() int {
	return f.requests("POST ") + f.requests("PATCH ") + f.requests("PUT ")
}

// order gives the place of the first request that holds each word, -1 for none.
func (f *finishing) order(words ...string) []int {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	out := make([]int, len(words))
	for i, w := range words {
		out[i] = slices.IndexFunc(f.w.fake.reqs, func(r string) bool { return strings.Contains(r, w) })
	}
	return out
}

func (f *finishing) pr(task string) int {
	for _, t := range f.register() {
		if t.ID == task {
			return t.PR
		}
	}
	return 0
}

func (f *finishing) mergedEvents(task string) [][]string {
	return slices.DeleteFunc(f.eventsOf(task), func(e []string) bool { return e[1] != "merged" })
}

// I1, the demo, with I6: a verified task goes through the close-out session at the tier of TierOf (the position of a Shape task is uphill, so the reasoning
// route, not the one of the plan alone), its checked push, layup/verify carried to the new head, the three other statuses, the mark and the merge at the head
// SHA, in that order, and the event merged; a second call is Merged and writes nothing.
func TestAVerifiedTaskIsClosedOutCarriedAndMergedAtItsHead(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	v := f.boundOf(f.arch)[0]
	got, err := f.finish(nil)
	h := f.closedHead(f.arch)
	if err != nil || !got.Merged || got.Task != f.arch || got.Head != h || h == "" {
		t.Fatalf("FinishTask = %+v, %v; want Merged at the bound head of the close-out %q", got, err, h)
	}
	if c := f.closers(f.arch); len(c) != 1 || c[0][sessionHarness] != "claude" || c[0][sessionBase] != v {
		t.Errorf("the close-out sessions %v; want one on the reasoning route (claude: the plan is small, the position uphill) at the verified head %s", c, v)
	}
	var got4 []string
	for _, s := range slices.Backward(statusesOf(f.w, h)) {
		got4 = append(got4, s["context"].(string)+" "+s["state"].(string))
		if s["context"] == VerifyContext && s["description"] != "verified head "+v+" plus the close-out" {
			t.Errorf("layup/verify at the new head: %q; want verified head %s plus the close-out", s["description"], v)
		}
	}
	if !slices.Equal(got4, []string{VerifyContext + " success", GatesContext + " success", SpecContext + " success", RulesContext + " success"}) {
		t.Errorf("the statuses at the new head, oldest first: %v; want verify carried, then gates, spec and rules, each success", got4)
	}
	n := strconv.Itoa(f.pr(f.arch))
	if o := f.order("/statuses/"+h, "POST /graphql", "PUT /repos/acme/target/pulls/"+n+"/merge"); o[0] < 0 || o[0] >= o[1] || o[1] >= o[2] {
		t.Errorf("the order of the writes %v; want the statuses, the mark, the merge at %s", o, h)
	}
	if m := f.mergedEvents(f.arch); len(m) != 1 || m[0][2] != "1" || m[0][5] != strings.TrimSpace(gitOut(t, f.w.bare, "rev-parse", "refs/heads/main")) {
		t.Errorf("the events merged: %v; want one, of attempt 1, at the merge commit", m)
	}
	if parents := strings.Fields(gitOut(t, f.w.bare, "rev-list", "--parents", "-n1", "refs/heads/main")); !slices.Contains(parents, h) {
		t.Errorf("the merge commit has the parents %v; want the head %s among them", parents[1:], h)
	}
	before, sessions := f.writes(), len(f.sessionsOf(f.arch))
	again, err := f.finish(nil)
	if err != nil || !again.Merged || again.Head != got.Head || f.writes() != before || len(f.sessionsOf(f.arch)) != sessions {
		t.Errorf("a second call = %+v, %v with %d writes and %d sessions more; want Merged and nothing written", again, err, f.writes()-before, len(f.sessionsOf(f.arch))-sessions)
	}
}

// failing wraps the forge of Sessions: each call is made, and the answer of the n-th SetStatus, of the first MarkReady and of the first Merge is lost.
type failing struct {
	forge.Forge
	set, n      int
	mark, merge bool
}

func (f *failing) SetStatus(ctx context.Context, sha, c, state, d string) error {
	err := f.Forge.SetStatus(ctx, sha, c, state, d)
	if f.n++; err == nil && f.n == f.set {
		return errors.New("the answer of the status was lost")
	}
	return err
}

func (f *failing) MarkReady(ctx context.Context, node string) (forge.PullRequest, error) {
	pr, err := f.Forge.MarkReady(ctx, node)
	if err == nil && f.mark {
		return pr, errors.New("the answer of the mark was lost")
	}
	return pr, err
}

func (f *failing) Merge(ctx context.Context, n int, sha string) (string, error) {
	merged, err := f.Forge.Merge(ctx, n, sha)
	if err == nil && f.merge {
		return merged, errors.New("the answer of the merge was lost")
	}
	return merged, err
}

// refusing refuses the records commits whose message holds the word, as a lost lease does, n times.
func refusing(word string, n int) func(*Sessions) {
	return func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if n > 0 && strings.Contains(m, word) {
				n--
				return &LostError{Holder: "0000000000000000"}
			}
			return nil
		}}
	}
}

// I3 (conditions 5 and 6 of the plan review): a restart after the bound, after an accepted status whose answer was lost, after the mark, after the merge
// and after two refused commits of the merged event adopts what was written. Each call is one stop; at the end there is one session, each status was posted
// once, the pull request was marked once and merged once, and the event merged is written once.
func TestARestartAfterEachWriteAdoptsIt(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	stops := []func(*Sessions){
		refusing(" bound to task/"+f.arch, 1),
		func(s *Sessions) { s.Forge = &failing{Forge: s.Forge, set: 3} },
		func(s *Sessions) { s.Forge = &failing{Forge: s.Forge, mark: true} },
		func(s *Sessions) { s.Forge = &failing{Forge: s.Forge, merge: true} },
		refusing(" merged as ", 1),
		refusing(" merged as ", 1),
	}
	for i, stop := range stops {
		if got, err := f.finish(stop); err == nil {
			t.Fatalf("call %d = %+v; want its stop", i+1, got)
		}
	}
	got, err := f.finish(nil)
	h := f.closedHead(f.arch)
	if err != nil || !got.Merged || got.Head != h {
		t.Fatalf("the last call = %+v, %v; want Merged at %s", got, err, h)
	}
	if n := len(f.closers(f.arch)); n != 1 {
		t.Errorf("%d close-out sessions; want one", n)
	}
	if st := statusesOf(f.w, h); len(st) != 4 {
		t.Errorf("%d statuses at the new head: %v; want the four, each posted once", len(st), headStates(f.drvWorld, h))
	}
	n := strconv.Itoa(f.pr(f.arch))
	if marks, merges := f.requests("POST /graphql"), f.requests("PUT /repos/acme/target/pulls/"+n+"/merge"); marks != 1 || merges != 1 {
		t.Errorf("%d marks and %d merges; want one of each", marks, merges)
	}
	if m := f.mergedEvents(f.arch); len(m) != 1 {
		t.Errorf("the events merged: %v; want one", m)
	}
}

// noPush makes the push of a head fail without pushing it, as a network that went down does.
type noPush struct{ SessionStore }

func (noPush) PushHead(context.Context, string, string) error {
	return errors.New("the network went down")
}

// noComment makes the post of a comment fail.
type noComment struct{ forge.Forge }

func (noComment) Comment(context.Context, int, string) (int64, error) {
	return 0, errors.New("the forge did not answer")
}

// I4, with I10: a pull request that a person closed without a merge is a wait that writes nothing, though a close-out's push is unresolved and its comment
// is missing (a repair would write both); once the pull request is open again, the repair finds the lost push and starts the next attempt, which is the
// handback in Advance.Attempt and not a Go error.
func TestAClosedUnmergedPullRequestIsAWaitWithNoWriteAndARepairThatStartsAnAttemptIsAHandback(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	if _, err := f.finish(func(s *Sessions) { s.Store, s.Forge = noPush{s.Store}, noComment{s.Forge} }); err == nil {
		t.Fatal("the first call gave no error; want the push that went down")
	}
	pr := f.pr(f.arch)
	f.w.fake.closePull(pr)
	head, writes, kinds := recordsHead(t, f.w), f.writes(), f.kinds(f.arch)
	want := "waiting for the Operator: " + f.arch + ": the pull request #" + strconv.Itoa(pr) + " of task/" + f.arch + "/1 is closed and not merged"
	for i := 0; i < 2; i++ {
		if got, err := f.finish(nil); err != nil || got.Wait != want || got.Merged || got.Attempt != 0 {
			t.Fatalf("call %d = %+v, %v; want the wait %q", i+1, got, err, want)
		}
	}
	if recordsHead(t, f.w) != head || f.writes() != writes || !slices.Equal(f.kinds(f.arch), kinds) {
		t.Errorf("the records head moved: %v, %d writes more, %d events more; want nothing written", recordsHead(t, f.w) != head, f.writes()-writes, len(f.kinds(f.arch))-len(kinds))
	}
	f.w.fake.mu.Lock()
	f.w.fake.pulls.list[pr-1].closed = false
	f.w.fake.mu.Unlock()
	got, err := f.finish(nil)
	if err != nil || got.Attempt != 2 || got.Merged || got.Wait != "" || got.Stop != "" {
		t.Fatalf("FinishTask after the repair = %+v, %v; want the handback of attempt 2", got, err)
	}
	if a := f.attempts(f.arch); len(a) != 2 {
		t.Errorf("the attempts %v; want the replacement started by the repair", a)
	}
}

// I5: a close-out outside its scope is refused by the push hook, is a stop that names the session and the reason, and is not tried again: a second call
// starts no session, and the new head is nowhere.
func TestACloseOutOutsideItsScopeIsRefusedAndNotRetried(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	f.ctlWrite("closeout-extra", "")
	got, err := f.finish(nil)
	if err != nil || !strings.HasPrefix(got.Stop, f.arch+": the session ") || !strings.Contains(got.Stop, " was refused (scope ") || got.Merged {
		t.Fatalf("FinishTask = %+v, %v; want the stop that names the refused session (scope)", got, err)
	}
	if again, err := f.finish(nil); err != nil || again != got || len(f.closers(f.arch)) != 1 {
		t.Errorf("a second call = %+v, %v with %d close-out sessions; want the same stop and one session", again, err, len(f.closers(f.arch)))
	}
	if h := f.closedHead(f.arch); h != "" || len(f.mergedEvents(f.arch)) != 0 {
		t.Errorf("a head %q was bound or a merge was recorded; want neither", h)
	}
}

// I7: a head that is not ready is a wait that names the reason, and a restart merges it once the job is green, with no second session.
func TestNotReadyWaitsAndARestartMerges(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	f.job("in_progress", "")
	got, err := f.finish(nil)
	h := f.closedHead(f.arch)
	want := "waiting: " + f.arch + ": not ready at " + h + ": job build: in_progress"
	if err != nil || got.Wait != want || got.Merged {
		t.Fatalf("FinishTask = %+v, %v; want the wait %q", got, err, want)
	}
	f.job("completed", "success")
	if got, err = f.finish(nil); err != nil || !got.Merged || got.Head != h || len(f.closers(f.arch)) != 1 {
		t.Errorf("the second call = %+v, %v with %d close-out sessions; want Merged at %s with one session", got, err, len(f.closers(f.arch)), h)
	}
}

// I8 (condition 5): when the close-out's harness leaves no harness outside the authors, VerifyAt at the new head would be not-active; the carried
// success stands, as no layup/verify is computed there.
func TestNoVerifyIsComputedAtTheCloseOutHead(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "devin", "devin")
	two := f.harnesses[:2] // claude and devin: the developer's and the close-out's
	got, err := f.finishTask(f.arch, two, nil)
	h := f.closedHead(f.arch)
	if err != nil || !got.Merged || got.Head != h {
		t.Fatalf("FinishTask = %+v, %v; want Merged at %s", got, err, h)
	}
	var verify []string
	for _, s := range statusesOf(f.w, h) {
		if s["context"] == VerifyContext {
			verify = append(verify, s["state"].(string)+" "+s["description"].(string))
		}
	}
	if len(verify) != 1 || !strings.HasPrefix(verify[0], "success verified head ") || !strings.HasSuffix(verify[0], " plus the close-out") {
		t.Errorf("the statuses layup/verify at the new head: %v; want the one success of the carry", verify)
	}
	base := f.eventsOf(f.arch)[0][4]
	if _, st, err := f.r.sessions().VerifyAt(context.Background(), VerifyInput{Task: f.arch, Base: base, Head: h, Waiting: true, Harnesses: two, Models: f.models}); err != nil || st.State != "failure" || !strings.HasPrefix(st.Description, "not-active") {
		t.Errorf("VerifyAt at the new head = %+v, %v; want not-active (the control: this world is the one where a computed verify would overwrite the carry)", st, err)
	}
}

// moveMain pushes a commit to the default branch of the target, as another merge does, and fetches it into the run's clone.
func (f *finishing) moveMain() {
	f.t.Helper()
	work := f.t.TempDir()
	gitOut(f.t, work, "clone", "-q", "-b", "main", f.w.bare, ".")
	if err := os.WriteFile(filepath.Join(work, "moved.txt"), []byte("moved\n"), 0o644); err != nil {
		f.t.Fatal(err)
	}
	gitOut(f.t, work, "add", "-A")
	gitOut(f.t, work, "commit", "-q", "-m", "another merge")
	gitOut(f.t, work, "push", "-q", "origin", "main")
	gitOut(f.t, f.r.target, "fetch", "-q", f.w.bare, "main")
}

// I9 (D9): a default branch that moved past the base of the attempt is a Go error and writes nothing; the base that is the tip passes.
func TestADefaultBranchThatMovedStopsTheCloseOutBeforeAnyWrite(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	f.moveMain()
	head, writes := recordsHead(t, f.w), f.writes()
	got, err := f.finish(nil)
	if want := f.arch + ": the base moved (row 66x's)"; err == nil || err.Error() != want {
		t.Fatalf("FinishTask = %+v, %v; want the error %q", got, err, want)
	}
	if recordsHead(t, f.w) != head || f.writes() != writes || len(f.closers(f.arch)) != 0 {
		t.Errorf("the records or the forge were written, or a session started; want nothing")
	}
}

// I11 (condition 3): a close-out session that ended completed, whose push hook did not run (the records commit of its push failed), resumes its checked
// push from the head it kept in the run's clone, with no second session.
func TestACloseOutThatEndedWithNoPushResumesItsCheckedPush(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	stop := func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, "layup run: the push of ") {
				return errors.New("the records branch did not answer")
			}
			return nil
		}}
	}
	if _, err := f.finish(stop); err == nil {
		t.Fatal("the first call gave no error; want the stop at the push")
	}
	closers := f.closers(f.arch)
	if len(closers) != 1 || slices.ContainsFunc(f.eventsOf(f.arch), func(e []string) bool { return e[3] == closers[0][sessionID] && (e[1] == "push" || e[1] == "bound") }) {
		t.Fatalf("the close-out sessions %v; want one that ended with no push or bound", closers)
	}
	got, err := f.finish(nil)
	if h := f.closedHead(f.arch); err != nil || !got.Merged || got.Head != h || h == "" || len(f.closers(f.arch)) != 1 {
		t.Errorf("the second call = %+v, %v with %d close-out sessions; want Merged at the head of the one session", got, err, len(f.closers(f.arch)))
	}
}

// newBetFinishing makes the world after the first bet: the specification task at its rendered commit r, the architecture task verified before the bet,
// the route of the close-out and the required job. v is the specification's verified head.
func newBetFinishing(t *testing.T) (f *finishing, v, r string) {
	t.Helper()
	d := betWorld(t, betOptionsTable)
	if got := d.bet(); got.Result != "done" || !strings.HasPrefix(got.Detail, "waiting: the answer of the bet") {
		t.Fatalf("the step bet = %+v; want the wait for the answer", got)
	}
	v = d.boundOf(d.register()[0].ID)[0]
	if _, err := d.r.renderedBet(context.Background(), v); err == nil {
		t.Error("renderedBet before the bet gave no error")
	}
	d.reply(7004, 303, "carol", "2026-10-08T14:03:00Z", betAnswerText)
	routing := keptRecord(t, d.w, "routing.tsv") + "software-developer\treasoning\t1\tclaude\tm1\nsoftware-developer\texecution\t1\tcodex\tm1\n"
	pushRecords(t, d.w.bare, map[string]string{"routing.tsv": routing}, "the route of the close-out")
	d.restart()
	if got := d.bet(); !strings.HasPrefix(got.Detail, "the bet 1: ") {
		t.Fatalf("the step bet = %+v; want the bet made", got)
	}
	spec, branch := d.specBranch()
	f = &finishing{drvWorld: d, spec: spec, arch: d.register()[1].ID}
	f.job("completed", "success")
	return f, v, d.branchHead(branch)
}

// I2 (conditions 5 and 6): the specification task's rendered commit gets the three statuses with its own inputs (the copy, the rendered records, the
// verified head as the pre-image, the register of the developer's session), the close-out goes on it, and the new head gets the same inputs with the
// register of the close-out session; the task merges at the head SHA.
func TestTheRenderedHeadGetsItsStatusesAndTheCloseOutGoesOnIt(t *testing.T) {
	t.Parallel()
	f, v, r := newBetFinishing(t)
	d, spec := f.drvWorld, f.spec
	rb, err := d.r.renderedBet(context.Background(), v)
	if err != nil || len(rb.Records) != 2 || !strings.HasPrefix(string(rb.Copy), "REQ-001: Must, phase 1\n") {
		t.Fatalf("renderedBet = %d records, copy %q, %v; want the two rendered records and the answer's copy", len(rb.Records), rb.Copy, err)
	}
	writes := f.writes()
	in := FinishInput{Task: spec, Records: d.r.store.Base, Harnesses: d.harnesses, Models: d.models}
	if _, err := d.r.sessions().FinishTask(context.Background(), in); err == nil || !strings.Contains(err.Error(), "no inputs of it") {
		t.Errorf("FinishTask with no inputs of the rendered head = %v; want the error that names them", err)
	}
	in.Bet = func(context.Context, string) (RenderedBet, error) {
		return RenderedBet{Copy: rb.Copy, Records: rb.Records[:1]}, nil
	}
	if _, err := d.r.sessions().FinishTask(context.Background(), in); err == nil || !strings.Contains(err.Error(), "is not a permitted rendered commit above the verified head "+v) || f.writes() != writes {
		t.Errorf("FinishTask with the records of another rendering = %v with %d writes more; want the refusal of the rendered commit and no write", err, f.writes()-writes)
	}
	got, err := f.finishTask(spec, d.harnesses, nil)
	h := f.closedHead(spec)
	if err != nil || !got.Merged || got.Head != h || h == "" {
		t.Fatalf("FinishTask = %+v, %v; want Merged at the head of the close-out %q", got, err, h)
	}
	for _, c := range []struct{ head, what string }{{r, "the rendered commit"}, {h, "the close-out"}} {
		states := headStates(d, c.head)
		for _, ctx := range []string{GatesContext, SpecContext, RulesContext, VerifyContext} {
			if !slices.Contains(states, ctx+" success") {
				t.Errorf("the statuses at %s: %v; want %s success", c.what, states, ctx)
			}
		}
	}
	if st := statusesOf(d.w, r); len(st) == 0 || st[len(st)-1]["description"] != "verified head "+v+" plus the rendered commit" {
		t.Errorf("the oldest status at the rendered commit: %v; want the carry of the step bet", st)
	}
	if c := f.closers(spec); len(c) != 1 || c[0][sessionBase] != r {
		t.Errorf("the close-out sessions %v; want one on the rendered commit %s", c, r)
	}
	if parents := strings.Fields(gitOut(t, d.w.bare, "rev-list", "--parents", "-n1", h)); len(parents) != 2 || parents[1] != r {
		t.Errorf("the close-out head has the parents %v; want the rendered commit %s", parents, r)
	}
	if m := f.mergedEvents(spec); len(m) != 1 || m[0][5] != strings.TrimSpace(gitOut(t, d.w.bare, "rev-parse", "refs/heads/main")) {
		t.Errorf("the events merged: %v; want one at the merge commit", m)
	}
}

// I12 (D5, finding 1 of round 1): a close-out session that ended completed, whose push hook did not run, is finished by its own checked push after the
// run's clone is made anew, as a restart makes it: the head it kept survives the new clone, and there is one producing session, however often it is called.
func TestACloseOutThatEndedWithNoPushResumesAfterTheCloneIsMadeAnew(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	stop := func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, "layup run: the push of ") {
				return errors.New("the records branch did not answer")
			}
			return nil
		}}
	}
	if _, err := f.finish(stop); err == nil {
		t.Fatal("the first call gave no error; want the stop at the push")
	}
	id := f.closers(f.arch)[0][sessionID]
	f.restart()
	got, err := f.finish(nil)
	h := f.closedHead(f.arch)
	if err != nil || !got.Merged || got.Head != h || h == "" {
		t.Fatalf("the call after the new clone = %+v, %v; want Merged at the head of the first session", got, err)
	}
	again, err := f.finish(nil)
	if c := f.closers(f.arch); err != nil || !again.Merged || len(c) != 1 || c[0][sessionID] != id {
		t.Errorf("the second call = %+v, %v with the close-out sessions %v; want Merged and the one session %s", again, err, c, id)
	}
}

// I12b: a kept head that is in no clone, and that no restart brought, is an error that names the session, not a second session.
func TestACloseOutWhoseKeptHeadIsGoneIsAnErrorAndNoSecondSession(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	stop := func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, "layup run: the push of ") {
				return errors.New("the records branch did not answer")
			}
			return nil
		}}
	}
	if _, err := f.finish(stop); err == nil {
		t.Fatal("the first call gave no error; want the stop at the push")
	}
	id := f.closers(f.arch)[0][sessionID]
	gitOut(t, f.r.target, "update-ref", "-d", "refs/layup/sessions/"+id)
	if _, err := f.finish(nil); err == nil || !strings.Contains(err.Error(), id) || len(f.closers(f.arch)) != 1 {
		t.Errorf("FinishTask = %v with %d sessions; want the error that names %s and no second session", err, len(f.closers(f.arch)), id)
	}
}

// I12c: a read of the refs of the clone that fails is the error of the read, not the absence of the kept head, and starts no second session.
func TestAFailedReadOfTheKeptHeadsIsTheErrorOfTheReadAndNotAnAbsence(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	stop := func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, "layup run: the push of ") {
				return errors.New("the records branch did not answer")
			}
			return nil
		}}
	}
	if _, err := f.finish(stop); err == nil {
		t.Fatal("the first call gave no error; want the stop at the push")
	}
	_, err := f.finish(func(s *Sessions) { s.Clone = filepath.Join(t.TempDir(), "gone") })
	if err == nil || strings.Contains(err.Error(), "is not in the run's clone") || len(f.closers(f.arch)) != 1 {
		t.Errorf("FinishTask = %v with %d sessions; want the error of the failed read and no second session", err, len(f.closers(f.arch)))
	}
}

// I14 (finding 2 of round 1): the architecture task that the bet opened, based on the rendered commit and verified with the answer's copy, is closed out: the
// full check of its new head reads the same copy, so that classes and ran pass, and it merges; without the copy the head is not ready.
func TestAnArchitectureCloseOutAfterTheBetCarriesTheAnswerCopy(t *testing.T) {
	t.Parallel()
	f, _, r := newBetFinishing(t)
	tasks := f.register()
	if len(tasks) != 3 {
		t.Fatalf("the register holds %d tasks; want the two of Shape and the one the bet opened", len(tasks))
	}
	arch := tasks[2].ID
	copy, err := f.r.betCopy(context.Background())
	if err != nil || !strings.HasPrefix(string(copy), "REQ-001: Must, phase 1\n") {
		t.Fatalf("betCopy = %q, %v; want the answer's copy", copy, err)
	}
	in := f.input(arch)
	in.Base, in.Copy = r, copy
	if got, err := f.r.sessions().AdvanceShapeTask(context.Background(), in); err != nil || !got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the verified head on the rendered commit", got, err)
	}
	got, err := f.finishTask(arch, f.harnesses, nil)
	h := f.closedHead(arch)
	if err != nil || !got.Merged || got.Head != h || h == "" {
		t.Fatalf("FinishTask = %+v, %v; want Merged at the head of the close-out (%v)", got, err, headStates(f.drvWorld, h))
	}
	if st := headStates(f.drvWorld, h); !slices.Contains(st, SpecContext+" success") {
		t.Errorf("the statuses at the new head: %v; want layup/spec success", st)
	}
}

// refusing409 closes the pull request when the merge is sent and answers 409, as a forge does for a head that changed in the same request.
type refusing409 struct {
	forge.Forge
	close func()
}

func (c refusing409) Merge(context.Context, int, string) (string, error) {
	c.close()
	return "", &forge.Error{Call: "merge", Status: 409, Message: "Head branch was modified"}
}

// I15 (finding 4 of round 1): a pull request that is open at every read of the call, closed when the merge is sent, and found closed at the same head by the
// read after the 409, is Q11's wait, with no error and no new attempt.
func TestAPullRequestClosedAtTheRefusedMergeIsTheWaitOfTheOperator(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	pr := f.pr(f.arch)
	got, err := f.finish(func(s *Sessions) { s.Forge = refusing409{s.Forge, func() { f.w.fake.closePull(pr) }} })
	want := "waiting for the Operator: " + f.arch + ": the pull request #" + strconv.Itoa(pr) + " of task/" + f.arch + "/1 is closed and not merged"
	if err != nil || got.Wait != want || got.Merged || got.Attempt != 0 || len(f.attempts(f.arch)) != 1 {
		t.Errorf("FinishTask = %+v, %v with %d attempts; want the wait %q and no new attempt", got, err, len(f.attempts(f.arch)), want)
	}
}

// closingPull closes the pull request of the task on the first read of it by number, which is the merge's.
type closingPull struct {
	forge.Forge
	close func()
}

func (c closingPull) PullRequest(ctx context.Context, q forge.PullRequestQuery) (forge.PullRequest, bool, error) {
	if q.Number != 0 {
		c.close()
	}
	return c.Forge.PullRequest(ctx, q)
}

// I13 (Q7 of the plan, condition 2): a pull request that a fresh read finds closed during the call, at the merge, is the wait of the Operator and not a wait for
// the readiness.
func TestAPullRequestClosedDuringTheCallIsTheWaitOfTheOperator(t *testing.T) {
	t.Parallel()
	f := newFinishing(t, "claude", "codex")
	pr := f.pr(f.arch)
	got, err := f.finish(func(s *Sessions) { s.Forge = closingPull{s.Forge, func() { f.w.fake.closePull(pr) }} })
	want := "waiting for the Operator: " + f.arch + ": the pull request #" + strconv.Itoa(pr) + " of task/" + f.arch + "/1 is closed and not merged"
	if err != nil || got.Wait != want || got.Merged || len(f.mergedEvents(f.arch)) != 0 {
		t.Errorf("FinishTask = %+v, %v; want the wait %q and no merge", got, err, want)
	}
}
