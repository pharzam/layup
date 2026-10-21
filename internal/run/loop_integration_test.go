//go:build integration

package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

// The step loop on the production path (task T-9nbf, #367, row 66k; loop.md, The step loop): the world of the first bet (newBetFinishing), the tail
// of a restart after the Scaffold (afterScaffold), state.sessions(), the httptest forge, real git and the scripted harnesses. A restart is a fresh state on
// the same records; a fault is the hook of the Sessions that the step makes (state.sessionsHook) or a wrapper of the forge.

var ctx0 = context.Background()

func (d *drvWorld) loop() Step { d.t.Helper(); return d.r.loopStep(ctx0) }

// tail restarts the world and runs the steps after the Scaffold as Restart does, to the first fail or wait.
func (d *drvWorld) tail() []Step {
	d.t.Helper()
	d.restart()
	return d.r.steps(ctx0, []string{"shape", "bet", "loop"}, 0, d.r.afterScaffold())
}

func (f *finishing) taskOf(id string) Task {
	f.t.Helper()
	for _, t := range f.register() {
		if t.ID == id {
			return t
		}
	}
	f.t.Fatalf("the register has no %s", id)
	return Task{}
}

// opened counts the issues that the forge opened.
func (f *finishing) opened() int {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	return len(f.w.fake.issueTexts)
}

// requestsFrom gives the place of the first request, from the place start, that is the line given or holds the word for a line that ends in a slash, -1 for none.
func (f *finishing) requestsFrom(start int, lines ...string) []int {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	out := make([]int, len(lines))
	for i, l := range lines {
		out[i] = -1
		for j := start; j < len(f.w.fake.reqs); j++ {
			if r := f.w.fake.reqs[j]; r == l || (strings.HasSuffix(l, "/") && strings.HasPrefix(r, l) && !strings.Contains(strings.TrimPrefix(r, l), "/")) {
				out[i] = j
				break
			}
		}
	}
	return out
}

func (f *finishing) nreqs() int {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	return len(f.w.fake.reqs)
}

// changes counts the requests of the forge that change a task: the writes but the tokens, which each restart asks for.
func (f *finishing) changes() int { return f.writes() - f.requests("POST /app/installations") }

// mergeSpecOutside merges the pull request of the specification at the rendered commit, as a person does, so that no close-out of its own is made.
func (f *finishing) mergeSpecOutside(rendered string) {
	f.t.Helper()
	_, branch := f.specBranch()
	f.follow(branch)
	f.mergePull(f.pr(f.spec))
	if _, _, err := f.r.sessions().MergeTask(ctx0, f.spec, rendered); err != nil {
		f.t.Fatal(err)
	}
}

// mergeAt is the place of the request that merges the pull request of the task.
func (f *finishing) mergeAt(task string) int {
	return f.order("PUT /repos/acme/target/pulls/" + strconv.Itoa(f.pr(task)) + "/merge")[0]
}

// moveMainWith pushes a commit that writes the files to the default branch of the target and fetches it into the run's clone.
func (f *finishing) moveMainWith(files map[string]string) {
	f.t.Helper()
	work := f.t.TempDir()
	gitOut(f.t, work, "clone", "-q", "-b", "main", f.w.bare, ".")
	for p, data := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(work, p)), 0o755); err != nil {
			f.t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(work, p), []byte(data), 0o644); err != nil {
			f.t.Fatal(err)
		}
	}
	gitOut(f.t, work, "add", "-A")
	gitOut(f.t, work, "commit", "-q", "-m", "another merge")
	gitOut(f.t, work, "push", "-q", "origin", "main")
	gitOut(f.t, f.r.target, "fetch", "-q", f.w.bare, "main")
}

// verifierOf gives the last session of the role of the verification: the plan review comes before the verifier.
func (f *finishing) verifierOf(task string) string {
	id := ""
	for _, r := range f.sessionsOf(task) {
		if r[sessionRole] == stepRole("verification") {
			id = r[sessionID]
		}
	}
	return id
}

// taskComments counts the comments of an issue that are the lines of the sessions of the task.
func (f *finishing) taskComments(issue int, task string) int {
	f.w.fake.mu.Lock()
	defer f.w.fake.mu.Unlock()
	n := 0
	for _, c := range f.w.fake.comments[issue] {
		if b, _ := c["body"].(string); strings.Contains(b, " session of "+task+", attempt ") {
			n++
		}
	}
	return n
}

// noPushOf makes the push of the heads of one task fail without pushing them, as a network that went down does.
type noPushOf struct {
	SessionStore
	task string
}

func (n noPushOf) PushHead(ctx context.Context, sha, branch string) error {
	if strings.HasPrefix(branch, "task/"+n.task+"/") {
		return errors.New("the network went down")
	}
	return n.SessionStore.PushHead(ctx, sha, branch)
}

// lostPost makes the call of the forge and loses its answer, once: the post of a comment, or the opening of the draft pull request of the task.
type lostPost struct {
	forge.Forge
	kind, task string
	fired      bool
}

func (l *lostPost) Comment(ctx context.Context, issue int, body string) (int64, error) {
	id, err := l.Forge.Comment(ctx, issue, body)
	if err == nil && l.kind == "comment" && !l.fired {
		l.fired = true
		return 0, errors.New("the answer of the comment was lost")
	}
	return id, err
}

func (l *lostPost) OpenDraft(ctx context.Context, head, base, title, body string) (forge.PullRequest, error) {
	pr, err := l.Forge.OpenDraft(ctx, head, base, title, body)
	if err == nil && l.kind == "open" && !l.fired && strings.HasPrefix(head, "task/"+l.task+"/") {
		l.fired = true
		return pr, errors.New("the answer of the opening was lost")
	}
	return pr, err
}

// With every task registered and verified or not, loop merges them in the order of the register after the bet. I1, the demo of goal a and
// condition 3 of the plan review: the merge of the specification, the architecture task of Shape carried over the merged base and merged, nothing to
// reopen, and the bet's architecture task run with the option it is bound to as the only survivor of its prompt; a second loop writes nothing.
func TestAfterTheBetLoopMergesTheShapeTasksAndTheBetsTaskInOrder(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	b := f.register()[2].ID
	issues, tasks := f.opened(), len(f.register())
	rows := f.tail()
	want := []Step{{"shape", "done", "done before"}, {"bet", "done", "done before"}, {"loop", "done", "the loop merged " + f.spec + ", " + f.arch + ", " + b}}
	if !slices.Equal(rows, want) {
		t.Fatalf("the rows = %+v; want %+v", rows, want)
	}
	if at := []int{f.mergeAt(f.spec), f.mergeAt(f.arch), f.mergeAt(b)}; at[0] < 0 || at[0] >= at[1] || at[1] >= at[2] {
		t.Errorf("the merges are at %v of the request log; want the specification, the architecture task of Shape, then the bet's", at)
	}
	for _, task := range []string{f.spec, f.arch, b} {
		if m := f.mergedEvents(task); len(m) != 1 {
			t.Errorf("%s: the events merged %v; want one", task, m)
		}
	}
	if f.opened() != issues || len(f.register()) != tasks {
		t.Errorf("an issue or a task was opened; want nothing to reopen")
	}
	prompts := strings.Join(f.promptsOf(b), "\n")
	if !strings.Contains(prompts, "the survivors of the filter:\n- db: pg\n") || strings.Contains(prompts, "pg sqlite") || strings.Contains(prompts, "- web:") {
		t.Errorf("the prompt files of %s do not name the option it is bound to as the only survivor:\n%s", b, prompts)
	}
	head, writes := recordsHead(t, f.w), f.changes()
	if again := f.tail(); again[2] != (Step{"loop", "done", "done before"}) || recordsHead(t, f.w) != head || f.changes() != writes {
		t.Errorf("a second loop = %+v with %d writes more; want done before and no write", again, f.changes()-writes)
	}
}

// I2, F1: a restart after the first close-out session is bound reaches loop: shape and bet are done before, and the cut of "verified" is the same
// in the step bet, in FinishTask and in the dispatch of loop.
func TestARestartAfterTheFirstCloseOutReachesLoop(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	f.job("in_progress", "") // the close-out is bound and the merge waits
	if got, err := f.finishTask(f.spec, f.harnesses, nil); err != nil || !strings.HasPrefix(got.Wait, "waiting: "+f.spec+": not ready at ") {
		t.Fatalf("FinishTask = %+v, %v; want the wait of a head that is not ready", got, err)
	}
	rows := f.tail()
	if len(rows) != 3 || rows[0] != (Step{"shape", "done", "done before"}) || rows[1] != (Step{"bet", "done", "done before"}) || rows[2].Result != "done" || !strings.HasPrefix(rows[2].Detail, "waiting: "+f.spec+": not ready at ") {
		t.Errorf("the rows = %+v; want shape and bet done before, then the wait of loop", rows)
	}
	if st, err := f.r.sessions().verifiedStep(ctx0, f.spec); err != nil || st.Head == "" {
		t.Errorf("verifiedStep after the close-out = %+v, %v; want the verified head (the cut of FinishTask)", st, err)
	}
}

// I2b, F2: a restart after the architecture task of Shape was redone (its new task starts at the default branch, not at the specification's verified head)
// reaches loop too.
func TestARestartAfterTheArchitectureTaskIsReplacedReachesLoop(t *testing.T) {
	t.Parallel()
	f, _, rendered := newBetFinishing(t)
	f.mergeSpecOutside(rendered) // no close-out of its own masks F1 here
	f.moveMainWith(map[string]string{"docs/spec/architecture.md": "## Moved\n\nAnother merge wrote this.\n"})
	adv, err := f.r.sessions().BringUpToDate(ctx0, CarryInput{Task: f.arch, Verifier: f.verifierOf(f.arch), Harnesses: f.harnesses, Models: f.models})
	if err != nil || adv.Replaced == "" {
		t.Fatalf("BringUpToDate = %+v, %v; want the redo of the task", adv, err)
	}
	in := f.input(adv.Replaced)
	in.Copy, _ = f.r.betCopy(ctx0) // every task after the bet is given the first copy of the answer (D10)
	if got, err := f.r.sessions().AdvanceShapeTask(ctx0, in); err != nil || !got.Verified {
		t.Fatalf("AdvanceShapeTask of the new task = %+v, %v; want its verified head", got, err)
	}
	f.restart()
	if got := f.r.atBoundary("shape", f.r.shapeStep)(ctx0); got != (Step{"shape", "done", "done before"}) {
		t.Errorf("shape after the redo = %+v; want done before", got)
	}
}

// I3, the characterization of row 66v with its composition control: a close-out push that was lost after the bet is refused push-lost, one attempt
// starts, the driver takes it to a verified head and the task merges (D7).
func TestALostCloseOutPushAfterTheBetIsRefusedRetriedAndMerged(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	b := f.register()[2].ID
	f.r.sessionsHook = func(s *Sessions) { s.Store = noPushOf{s.Store, b} }
	if got := f.loop(); got.Result != "fail" || !strings.Contains(got.Detail, "the network went down") {
		t.Fatalf("the first loop = %+v; want the fail of the push that went down", got)
	}
	f.restart()
	if got := f.loop(); got.Result != "done" || got.Detail != "the loop merged "+b {
		t.Fatalf("the second loop = %+v; want %s merged", got, b)
	}
	kinds := strings.Join(f.kinds(b), " ")
	if !strings.Contains(kinds, "push refused attempt") || len(f.attempts(b)) != 2 || len(f.mergedEvents(b)) != 1 {
		t.Errorf("the events of %s: %s, %d attempts; want push, refused, a second attempt, and one merge", b, kinds, len(f.attempts(b)))
	}
}

// I4, the characterization of row 66v with its control: the comment of a session of a task after the bet is on its issue and not on the control issue.
func TestASessionCommentOfATaskAfterTheBetIsOnItsIssue(t *testing.T) {
	t.Parallel()
	f, _, rendered := newBetFinishing(t)
	f.mergeSpecOutside(rendered)
	arch := f.taskOf(f.arch)
	control, _ := strconv.Atoi(f.r.value("issue.control"))
	before := f.taskComments(arch.Issue, arch.ID)
	if adv, err := f.r.actOn(ctx0, arch, false); err != nil || !adv.Merged {
		t.Fatalf("the act on %s = %+v, %v; want it merged", arch.ID, adv, err)
	}
	if on, off := f.taskComments(arch.Issue, arch.ID)-before, f.taskComments(control, arch.ID); on == 0 || off != 0 {
		t.Errorf("%d new comments of the close-out of %s on its issue and %d on the control issue; want some and none", on, arch.ID, off)
	}
}

// I5, condition 2 of the plan review: the merge was made and the commit of its event was refused; loop records the event of a merged pull request before
// any base movement, and sends no second merge.
func TestALostMergedEventIsRecordedByLoopWithNoSecondMerge(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	f.r.sessionsHook = refusing(" merged as ", 2)
	if got := f.loop(); got.Result != "fail" || len(f.mergedEvents(f.spec)) != 0 || f.requests("PUT /repos/acme/target/pulls/"+strconv.Itoa(f.pr(f.spec))+"/merge") != 1 {
		t.Fatalf("the first loop = %+v with the events %v; want a fail after the merge and before its event", got, f.mergedEvents(f.spec))
	}
	f.restart()
	if got := f.loop(); got.Result != "done" || !strings.HasPrefix(got.Detail, "the loop merged "+f.spec+", ") {
		t.Fatalf("the second loop = %+v; want the specification recorded and the others merged", got)
	}
	if m := f.mergedEvents(f.spec); len(m) != 1 || f.requests("PUT /repos/acme/target/pulls/"+strconv.Itoa(f.pr(f.spec))+"/merge") != 1 {
		t.Errorf("the events %v and %d merge requests of the specification; want one event and one request", m, f.requests("PUT /repos/acme/target/pulls/"+strconv.Itoa(f.pr(f.spec))+"/merge"))
	}
}

// I6, condition 7: a restart after each write adopts it: the comment of the close-out, a status, the mark, the merge, its event and the opening of the
// pull request of the last task. At the end, one close-out session and one merge of each pull request, one mark each, and no comment twice.
func TestARestartAfterEachWriteOfTheLoopAdoptsIt(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	b := f.register()[2].ID
	pulls := f.requests("POST /repos/acme/target/pulls")
	stops := []func(){
		func() { f.r.cfg.Forge = &lostPost{Forge: f.r.cfg.Forge, kind: "comment"} },
		func() { f.r.cfg.Forge = &failing{Forge: f.r.cfg.Forge, set: 3} },
		func() { f.r.cfg.Forge = &failing{Forge: f.r.cfg.Forge, mark: true} },
		func() { f.r.cfg.Forge = &failing{Forge: f.r.cfg.Forge, merge: true} },
		func() { f.r.sessionsHook = refusing(" merged as ", 1) },
		func() { f.r.cfg.Forge = &lostPost{Forge: f.r.cfg.Forge, kind: "open", task: b} },
	}
	for i, stop := range stops {
		f.restart()
		stop()
		if got := f.loop(); got.Result != "fail" {
			t.Fatalf("stop %d: the loop = %+v; want the fail of its fault", i+1, got)
		}
	}
	if got := f.tail(); got[2].Result != "done" || got[2].Detail != "the loop merged "+b {
		t.Fatalf("the last restart = %+v; want the last task merged (the others did in the restarts before)", got)
	}
	for _, task := range []string{f.spec, f.arch, b} {
		n := strconv.Itoa(f.pr(task))
		if c, m := len(f.closers(task)), f.requests("PUT /repos/acme/target/pulls/"+n+"/merge"); c != 1 || m != 1 || len(f.mergedEvents(task)) != 1 {
			t.Errorf("%s: %d close-out sessions, %d merge requests, %d events merged; want one of each", task, c, m, len(f.mergedEvents(task)))
		}
		if n := f.taskComments(f.taskOf(task).Issue, task); n != len(f.sessionsOf(task)) {
			t.Errorf("%s: %d session comments for %d sessions; want one comment of each", task, n, len(f.sessionsOf(task)))
		}
	}
	if marks, opened := f.requests("POST /graphql"), f.requests("POST /repos/acme/target/pulls")-pulls; marks != 3 || opened != 1 {
		t.Errorf("%d marks and %d pull requests opened; want 3 and the one of the last task", marks, opened)
	}
}

// afterMerge runs on after the first merge of a pull request, as another merge to the default branch would come.
type afterMerge struct {
	forge.Forge
	on func()
	n  int
}

func (a *afterMerge) Merge(ctx context.Context, number int, sha string) (string, error) {
	merged, err := a.Forge.Merge(ctx, number, sha)
	if err == nil {
		if a.n++; a.n == 1 {
			a.on()
		}
	}
	return merged, err
}

// watchControl wraps the forge: after the n-th read of the comments of the control issue it runs on(n), and it notes how many reads came before the first merge.
type watchControl struct {
	forge.Forge
	issue, reads, atMerge int
	on                    func(n int)
}

func (w *watchControl) Comments(ctx context.Context, issue int) ([]forge.Comment, error) {
	got, err := w.Forge.Comments(ctx, issue)
	if issue == w.issue {
		w.reads++
		w.on(w.reads)
	}
	return got, err
}

func (w *watchControl) Merge(ctx context.Context, n int, sha string) (string, error) {
	if w.atMerge == 0 {
		w.atMerge = w.reads
	}
	return w.Forge.Merge(ctx, n, sha)
}

// I7, condition 4: a set is applied at each boundary: before shape, before bet, before loop and between two acts of loop. The comments of the control issue
// are read three times before the first merge, and a set that comes after the first act is applied before the second.
func TestASetIsAppliedAtEachBoundaryOfTheTail(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	control, _ := strconv.Atoi(f.r.value("issue.control"))
	set := func(id int, text string) {
		addComment(f.w, control, id, 101, "pharzam", "2026-10-08T15:00:00Z", text)
	}
	set(6300000001, "set brief.lines 70 because the Operator wants longer briefs")
	w := &watchControl{issue: control, on: func(n int) {
		switch n {
		case 1:
			set(6300000002, "set lease.H 9 because the Operator wants a longer lease")
		case 3:
			set(6300000003, "set build.parallel 4 because the host has four cores")
		}
	}}
	f.restart()
	w.Forge, f.r.cfg.Forge = f.r.cfg.Forge, w
	rows := f.r.steps(ctx0, []string{"shape", "bet", "loop"}, 0, f.r.afterScaffold())
	if len(rows) != 3 || rows[2].Result != "done" {
		t.Fatalf("the rows = %+v; want the three steps done", rows)
	}
	var log []string
	for _, s := range strings.Split(gitOut(t, f.w.bare, "log", "--reverse", "--format=%s", "layup-records"), "\n") {
		if strings.HasPrefix(s, "parameters: set ") || strings.Contains(s, " merged as ") {
			log = append(log, strings.TrimPrefix(s, "parameters: set "))
		}
	}
	order := strings.Join(log, "|")
	if !strings.HasPrefix(order, "brief.lines|lease.H|") || !strings.Contains(order, "|build.parallel|") || strings.Index(order, "build.parallel") < strings.Index(order, " merged as ") ||
		strings.Count(order[:strings.Index(order, "build.parallel")], " merged as ") != 1 {
		t.Errorf("the records commits %q; want brief.lines, lease.H, the merge of the specification, build.parallel between the first and the second act", order)
	}
	if w.atMerge != 3 {
		t.Errorf("%d reads of the control issue before the first merge; want 3, one for each of shape, bet and loop", w.atMerge)
	}
}

// The condition 4 skip: an absent parameters.tsv is skipped only before the first task, and a table that is absent after it, refused by its reader or unreadable
// is an error that names it; the F3 probe is the first assertion at the stub.
func TestAnAbsentParametersTableIsSkippedOnlyBeforeTheFirstTask(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	head := recordsHead(t, d.w)
	if err := d.r.boundary(ctx0); err != nil || recordsHead(t, d.w) != head {
		t.Fatalf("the boundary before the first task = %v (records moved: %v); want a skip", err, recordsHead(t, d.w) != head)
	}
	d.open()
	d.restart()
	if err := d.r.boundary(ctx0); err != nil {
		t.Fatalf("the boundary after the first task = %v; want it called", err)
	}
	if err := d.r.store.Commit(ctx0, map[string][]byte{parametersPath: nil}, "the table removed"); err != nil {
		t.Fatal(err)
	}
	if err := d.r.boundary(ctx0); err == nil || !strings.Contains(err.Error(), parametersPath) {
		t.Errorf("the boundary with the table absent after the first task = %v; want an error that names it", err)
	}
	pushRecords(t, d.w.bare, map[string]string{parametersPath: "broken\n"}, "a table the reader refuses")
	d.restart()
	if got := d.r.atBoundary("loop", func(context.Context) Step { return done("loop", "ran") })(ctx0); got.Result != "fail" || !strings.HasPrefix(got.Detail, "parameters: ") {
		t.Errorf("a refused table: %+v; want parameters: …", got)
	}
}

// I8, condition 4: a missing permission of the loop is the row of the step, forge: …, before any boundary write or act, for each of shape, bet and loop.
func TestAMissingLoopPermissionFailsTheRowOfEachStepBeforeAnyWrite(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	control, _ := strconv.Atoi(f.r.value("issue.control"))
	addComment(f.w, control, 6300000001, 101, "pharzam", "2026-10-08T15:00:00Z", "set brief.lines 70 because the Operator wants longer briefs")
	f.w.fake.mu.Lock()
	f.w.fake.perms = map[string]string{"contents": "write", "issues": "write", "metadata": "read", "pull_requests": "write", "checks": "read"}
	f.w.fake.mu.Unlock()
	f.restart()
	head, writes := recordsHead(t, f.w), f.changes()
	for _, c := range []struct {
		name string
		step func(context.Context) Step
	}{{"shape", f.r.shapeStep}, {"bet", f.r.betStep}, {"loop", f.r.loopStep}} {
		got := f.r.atBoundary(c.name, c.step)(ctx0)
		if c.name == "loop" {
			got = f.loop()
		}
		if got.Name != c.name || got.Result != "fail" || got.Detail != "forge: the installation lacks the permissions of the task loop: statuses: write" {
			t.Errorf("%s = %+v; want its row, fail, forge: the installation lacks the permissions of the task loop: statuses: write", c.name, got)
		}
	}
	if recordsHead(t, f.w) != head || f.changes() != writes {
		t.Errorf("the records head moved: %v, %d writes more; want no boundary write and no act", recordsHead(t, f.w) != head, f.changes()-writes)
	}
}

// I9, the inventory row of the stand-in forge through loop: a gate failure goes back to the developer as a new attempt, the pull request body links the records
// commit, and the merge is at the head SHA that was read (the stand-in refuses a merge at any other).
func TestTheTaskLoopWithTheStandInForgeThroughLoop(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	id := f.register()[2].ID
	f.ctlWrite("gofile-arch", "")
	if got := f.loop(); got.Result != "done" || !strings.HasSuffix(got.Detail, id) {
		t.Fatalf("loop = %+v; want the tasks merged", got)
	}
	b := f.taskOf(id)
	if len(f.attempts(b.ID)) != 2 {
		t.Errorf("the attempts of %s: %v; want the gate failure to start a second", b.ID, f.attempts(b.ID))
	}
	f.w.fake.mu.Lock()
	body := f.w.fake.pulls.list[b.PR-1].body
	f.w.fake.mu.Unlock()
	if !strings.HasPrefix(body, fmt.Sprintf("Closes #%d\n\n", b.Issue)) || !strings.Contains(body, "/commit/") {
		t.Errorf("the body of the pull request of %s: %q; want Closes #%d and the records commit", b.ID, body, b.Issue)
	}
	if parents := strings.Fields(gitOut(t, f.w.bare, "rev-list", "--parents", "-n1", "refs/heads/main")); !slices.Contains(parents[1:], f.closedHead(b.ID)) {
		t.Errorf("the merge commit has the parents %v; want the head %s of the close-out among them", parents[1:], f.closedHead(b.ID))
	}
}

// I9b: a close-out that touches README.md is refused for its scope, and loop stops there with the reason.
func TestACloseOutThatTouchesTheReadmeIsRefusedThroughLoop(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	f.ctlWrite("closeout-extra", "")
	if got := f.loop(); got.Result != "fail" || !strings.Contains(got.Detail, " was refused (scope ") || len(f.mergedEvents(f.spec)) != 0 {
		t.Errorf("loop = %+v; want a fail that names the refused session (scope) and no merge", got)
	}
}

// I10, goal b: the architecture task of Shape overlaps what another merge wrote and is redone through loop: the pull request is closed, then the issue, a new
// issue is opened, one records commit holds the events closed and rebased, and the new task is driven to its merge, with the bet's task after it. The
// control: a base that moved by a file the task does not touch makes no redo.
func TestAnOverlappingMovedBaseIsRedoneThroughLoop(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	b := f.register()[2].ID
	pr, start := f.pr(f.arch), f.nreqs()
	f.r.cfg.Forge = &afterMerge{Forge: f.r.cfg.Forge, on: func() {
		f.moveMainWith(map[string]string{"docs/spec/architecture.md": "## Moved\n\nAnother merge wrote this.\n"})
	}}
	got := f.loop()
	if got.Result != "done" || !strings.HasPrefix(got.Detail, "the loop merged "+f.spec+", ") || strings.Contains(got.Detail, f.arch) || !strings.HasSuffix(got.Detail, ", "+b) {
		t.Fatalf("loop = %+v; want the specification, the task that redoes %s, and %s merged", got, f.arch, b)
	}
	kinds := f.kinds(f.arch)
	if !slices.Contains(kinds, "closed") || len(f.mergedEvents(f.arch)) != 0 {
		t.Errorf("the events of %s: %v; want it closed and not merged", f.arch, kinds)
	}
	o := f.requestsFrom(start, "PATCH /repos/acme/target/pulls/"+strconv.Itoa(pr), "PATCH /repos/acme/target/issues/", "POST /repos/acme/target/issues")
	if o[0] < 0 || o[0] >= o[1] || o[1] >= o[2] {
		t.Errorf("the order of the redo %v; want the pull request closed, the issue closed, then a new issue", o)
	}
	if reg := f.register(); len(reg) != 3 || reg[1].ID == f.arch || !slices.Equal(f.taskOf(b).Predecessors, []string{f.spec, reg[1].ID}) {
		t.Errorf("the register %+v; want the new row in place of %s, and %s after it", reg, f.arch, b)
	}
	c, _, _ := newBetFinishing(t)
	opened := c.opened()
	c.r.cfg.Forge = &afterMerge{Forge: c.r.cfg.Forge, on: func() { c.moveMainWith(map[string]string{"elsewhere.txt": "x\n"}) }}
	if got := c.loop(); got.Result != "done" || c.opened() != opened || slices.Contains(c.kinds(c.arch), "closed") || len(c.mergedEvents(c.arch)) != 1 {
		t.Errorf("the control: loop = %+v with %d issues opened; want no redo for a base that moved elsewhere, and %s merged", got, c.opened()-opened, c.arch)
	}
}

// I11 (6304631026, D9): a reopened task is driven to a verified head and merged by FinishTask with no rendered commit asked, and its head is not judged as a
// specification that waits for its bet.
func TestAReopenedTaskIsMergedByLoopWithNoRenderedCommit(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	main := w.merge()
	tasks, err := w.open()
	if err != nil || len(tasks) == 0 {
		t.Fatalf("openReopenings = %v, %v", tasks, err)
	}
	task := tasks[0]
	in := w.input(task.ID)
	in.Base = main
	pc, err := w.r.sessions().PromptFile(ctx0, PromptIn{Task: task.ID, Step: "specification", Base: main})
	if err != nil {
		t.Fatal(err)
	}
	w.ctlWrite("reopen-"+task.ID, shownLine(string(pc.Text("/x")), "need"))
	w.ctlWrite("reopen-section-"+task.ID, task.Requirements[0])
	if in.Copy, err = w.r.betCopy(ctx0); err != nil {
		t.Fatal(err)
	}
	if adv, err := w.r.sessions().AdvanceShapeTask(ctx0, in); err != nil || !adv.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the verified head", adv, err)
	}
	adv, err := w.r.actOn(ctx0, w.taskOf(task.ID), false)
	if err != nil || !adv.Merged || len(w.mergedEvents(task.ID)) != 1 {
		t.Errorf("the act on the reopened task = %+v, %v; want it merged with no rendered commit asked (%v)", adv, err, w.kinds(task.ID))
	}
}

// I13 (F4, Q2): a lost close-out push of the specification task after the bet starts an attempt that has no rendered commit; loop stops with the
// Operator's wait before any driver, at every restart.
func TestAnAttemptAfterTheBetWithNoRenderedCommitIsTheStopOfTheOperator(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	if _, err := f.finishTask(f.spec, f.harnesses, func(s *Sessions) { s.Store = noPushOf{s.Store, f.spec} }); err == nil || !strings.Contains(err.Error(), "the network went down") {
		t.Fatalf("FinishTask with the push down = %v; want its error", err)
	}
	f.restart()
	if got, err := f.finishTask(f.spec, f.harnesses, nil); err != nil || got.Attempt != 2 {
		t.Fatalf("FinishTask after the loss = %+v, %v; want the handback of attempt 2", got, err)
	}
	// F4: the attempt has no rendered commit and no pull request, and FinishTask has nothing to finish: an error, while the step bet is done.
	if _, err := f.finishTask(f.spec, f.harnesses, nil); err == nil {
		t.Fatal("FinishTask of the attempt with no rendered commit gave no error; the premise of F4 is gone")
	}
	want := Step{"loop", "fail", "waiting for the Operator: " + f.spec + ": an attempt after the bet has no rendered commit"}
	for i := 0; i < 2; i++ {
		f.restart()
		sessions := len(f.sessionsOf(f.spec))
		if got := f.loop(); got != want || len(f.sessionsOf(f.spec)) != sessions {
			t.Fatalf("loop %d after the loss = %+v with %d sessions more; want %+v and no session of a driver", i+1, got, len(f.sessionsOf(f.spec))-sessions, want)
		}
	}
	if len(f.attempts(f.spec)) != 2 || len(f.mergedEvents(f.spec)) != 0 {
		t.Errorf("%d attempts, %d merges; want the replacement attempt and no merge", len(f.attempts(f.spec)), len(f.mergedEvents(f.spec)))
	}
}

// I12 (Q1): the close-out of a task that loop acts on is routed from the progress position, uphill for every task of M2e, and not from the size of its
// plan: the plan here is small, which alone would give the execution route.
func TestTheCloseOutOfALoopTaskIsRoutedFromItsPositionNotItsPlan(t *testing.T) {
	t.Parallel()
	f, _, rendered := newBetFinishing(t)
	f.mergeSpecOutside(rendered) // the base of the architecture task of Shape is the specification's verified head, an ancestor of the default branch once it has merged
	arch := f.taskOf(f.arch)
	if adv, err := f.r.actOn(ctx0, arch, false); err != nil || !adv.Merged {
		t.Fatalf("the act on %s = %+v, %v; want it merged", arch.ID, adv, err)
	}
	if c := f.closers(f.arch); len(c) != 1 || c[0][sessionHarness] != "claude" {
		t.Errorf("the close-out sessions %v; want one on the reasoning route (claude), the execution route being codex", c)
	}
}

// Condition 2 of the plan review, M18: a pull request that a person closed without a merge is the wait of loop, and loop writes nothing, though a close-out's
// push is unresolved and its comment is missing (a repair would write both); it repeats at every restart.
func TestAClosedUnmergedPullRequestIsTheWaitOfLoopWithNoWrite(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	if _, err := f.finishTask(f.spec, f.harnesses, func(s *Sessions) { s.Store, s.Forge = noPush{s.Store}, noComment{s.Forge} }); err == nil {
		t.Fatal("FinishTask with the push down gave no error")
	}
	pr := f.pr(f.spec)
	f.w.fake.closePull(pr)
	head, writes, kinds := recordsHead(t, f.w), f.changes(), f.kinds(f.spec)
	want := Step{"loop", "done", "waiting for the Operator: " + f.spec + ": the pull request #" + strconv.Itoa(pr) + " of task/" + f.spec + "/1 is closed and not merged"}
	for i := 0; i < 2; i++ {
		f.restart()
		if got := f.loop(); got != want {
			t.Fatalf("loop %d = %+v; want %+v", i+1, got, want)
		}
	}
	if recordsHead(t, f.w) != head || f.changes() != writes || !slices.Equal(f.kinds(f.spec), kinds) {
		t.Errorf("the records head moved: %v, %d writes more; want nothing written", recordsHead(t, f.w) != head, f.changes()-writes)
	}
}

// Condition 2 of the plan review: the event merged of the specification was recorded before loop starts and nothing was opened after it; loop opens the tasks
// that the bet reopens at its start (the reading is whenever the event exists, not only when this invocation writes it).
func TestLoopOpensTheReopeningsOfAMergeThatWasRecordedBefore(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	w.merge()
	tasks := len(w.register())
	w.restart()
	w.loop() // the drive of the reopened tasks is I11's: the row is not the question here
	if reopened := len(w.register()) - tasks; reopened != 4 || len(w.rows(t)) != 4 {
		t.Errorf("%d tasks and %d reservations opened by loop; want the 4 that the bet re-classed", reopened, len(w.rows(t)))
	}
}

// A base that moved by a file the specification does not touch: loop carries the verdict over the merge above the rendered commit (the primitive is row 66x's;
// the loop gives it the rendered input) and merges the three tasks.
func TestAMovedBaseOfTheSpecificationWithItsRenderedCommitIsCarriedThroughLoop(t *testing.T) {
	t.Parallel()
	f, _, _ := newBetFinishing(t)
	b := f.register()[2].ID
	f.moveMainWith(map[string]string{"elsewhere.txt": "x\n"})
	if got := f.loop(); got.Result != "done" || got.Detail != "the loop merged "+f.spec+", "+f.arch+", "+b {
		t.Fatalf("loop = %+v; want the three tasks merged", got)
	}
	if merges := slices.DeleteFunc(f.eventsOf(f.spec), func(e []string) bool { return e[1] != "bound" || e[3] != "" }); len(merges) < 2 {
		t.Errorf("the bounds of code of %s: %v; want the rendered commit and the merge of the default branch above it", f.spec, merges)
	}
}

// Before the bet there is nothing for loop to do: it fails and writes nothing.
func TestLoopBeforeTheBetFailsAndWritesNothing(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	head := recordsHead(t, d.w)
	if got := d.loop(); got.Result != "fail" || got.Detail != "bets.tsv holds no bet: the step bet comes first" || recordsHead(t, d.w) != head {
		t.Errorf("loop before the bet = %+v; want the fail, and no write", got)
	}
}
