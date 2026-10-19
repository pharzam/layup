//go:build integration

package run

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

// The pull request and the merge through the adapter and the stand-in forge of the HTTP fixture, with real git (row 66g, task T-xvvd,
// #363): the draft of a bound head, the ready and merge rule, the merge at the head SHA, and the tree after it.

const requiredBuild = `[{"type":"required_status_checks","parameters":{"required_status_checks":[{"context":"build","integration_id":15368}]}}]`

// A pullShip is a developer's head bound on the task branch, the production Sessions with the App's bot, and the draft opened on it.
type pullShip struct {
	*loopRun
	s    *Sessions
	head string
	pr   forge.PullRequest
}

func newPullShip(t *testing.T) *pullShip {
	t.Helper()
	l := newLoop(t, nil)
	l.r.Bot = 9001
	head := l.bound(l.dev(1, "claude", map[string]string{"x.go": "package x\n"}, okResult))
	if head == "" {
		t.Fatal("the developer's session was not bound")
	}
	s := l.prodBot()
	// The draft is opened and registered directly, so that the tests of the ready and merge rule stand on the forge and the register alone.
	ctx := context.Background()
	pr, err := l.r.Forge.OpenDraft(ctx, "task/T-ab12/1", "main", "feat: add the draft", "Closes #3")
	if err != nil {
		t.Fatalf("OpenDraft: %v", err)
	}
	tasks, err := ReadRegister([]byte(keptRecord(t, l.w, "tasks.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].PR = pr.Number
	var b bytes.Buffer
	if err := WriteRegister(&b, tasks); err != nil {
		t.Fatal(err)
	}
	if err := l.r.commit(ctx, map[string][]byte{tasksPath: b.Bytes()}, "the pull request of T-ab12"); err != nil {
		t.Fatal(err)
	}
	l.w.fake.mu.Lock()
	l.w.fake.rules = map[string]string{"main": requiredBuild}
	l.w.fake.mu.Unlock()
	return &pullShip{loopRun: l, s: s, head: head, pr: pr}
}

func (p *pullShip) setRuns(sha string, runs ...map[string]any) {
	p.w.fake.mu.Lock()
	defer p.w.fake.mu.Unlock()
	if p.w.fake.checkRuns == nil {
		p.w.fake.checkRuns = map[string][]map[string]any{}
	}
	p.w.fake.checkRuns[sha] = runs
}

func actionsRun(status, conclusion string, app int) map[string]any {
	return map[string]any{"name": "build", "status": status, "conclusion": conclusion, "app": map[string]any{"id": app, "slug": "github-actions"}}
}

// green posts the four statuses of the App on the head and gives the job build a success of GitHub Actions.
func (p *pullShip) green(sha string) {
	p.t.Helper()
	for _, c := range statusContexts {
		if _, err := p.s.PostStatus(context.Background(), sha, Status{c, "success", "pass"}); err != nil {
			p.t.Fatal(err)
		}
	}
	p.setRuns(sha, actionsRun("completed", "success", 15368))
}

func (p *pullShip) foreign(sha, context, state string) {
	p.w.fake.mu.Lock()
	defer p.w.fake.mu.Unlock()
	p.w.fake.statuses[sha] = append([]map[string]any{{"context": context, "state": state, "description": "x", "creator": map[string]any{"id": 7, "login": "someone"}}}, p.w.fake.statuses[sha]...)
}

func (p *pullShip) requests() []string {
	p.w.fake.mu.Lock()
	defer p.w.fake.mu.Unlock()
	return slices.Clone(p.w.fake.reqs)
}

func countOf(reqs []string, prefix string) int {
	n := 0
	for _, r := range reqs {
		if strings.HasPrefix(r, prefix) {
			n++
		}
	}
	return n
}

func (p *pullShip) mainSHA() string {
	return strings.TrimSpace(gitOut(p.t, p.w.bare, "rev-parse", "refs/heads/main"))
}

func (p *pullShip) merges() []string {
	var out []string
	for _, e := range p.events() {
		if e[1] == "merged" {
			out = append(out, strings.Join(e[2:7], " "))
		}
	}
	return out
}

const pullsPath = "/repos/acme/target/pulls"

func TestTheDraftPullRequestOfABoundHead(t *testing.T) {
	l := newLoop(t, nil)
	l.r.Bot = 9001
	ctx := context.Background()
	head := l.bound(l.dev(1, "claude", map[string]string{"x.go": "package x\n"}, okResult))
	prod := l.prodBot()
	in := PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"}
	pr, err := prod.OpenPullRequest(ctx, in)
	if err != nil || pr.Number != 1 || !pr.Draft || pr.HeadSHA != head {
		t.Fatalf("OpenPullRequest = %+v, %v; want draft 1 at %s", pr, err, head)
	}
	got := l.w.fake.pulls.list[0]
	web := l.w.config("0.1.0-dev").Register.Web
	want := "Closes #3\n\nThe records commit this task started from: " + web + "/acme/target/commit/" + l.a.records
	if got.title != "feat: add the draft" || got.body != want || got.base != "main" || got.branch != "task/T-ab12/1" || !got.draft || l.w.fake.pulls.noDraft {
		t.Errorf("the request: title %q, body %q, base %q, branch %q, draft %v; want the title, the body %q, main, the task branch, a draft", got.title, got.body, got.base, got.branch, got.draft, want)
	}
	tasks, err := ReadRegister([]byte(keptRecord(t, l.w, "tasks.tsv")))
	if err != nil || tasks[0].PR != 1 {
		t.Errorf("tasks.tsv: %+v, %v; want pr 1", tasks, err)
	}
	// A restart: the same call finds the pull request by its branch and posts no second.
	if again, err := afterRestart(prod).OpenPullRequest(ctx, in); err != nil || again.Number != 1 || countOf(l.w.fake.reqs, "POST "+pullsPath) != 1 {
		t.Errorf("after a restart: %+v, %v with %d posts; want the same pull request and one post in all", again, err, countOf(l.w.fake.reqs, "POST "+pullsPath))
	}
	// The second attempt is opened by the Sessions of a restart that never saw the first: the address of the web comes from the run's own
	// constructor and survives the restart (the plan review of #363, condition 8).
	l.attempt(2)
	l.bound(l.dev(2, "claude", map[string]string{"y.go": "package x\n"}, okResult))
	if pr2, err := afterRestart(prod).OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 2, Type: "fix", Description: "mend the draft"}); err != nil || pr2.Number != 2 {
		t.Fatalf("the second attempt after a restart: %+v, %v", pr2, err)
	}
	if body := l.w.fake.pulls.list[1].body; !strings.Contains(body, web+"/acme/target/commit/") || strings.Contains(body, "/commit/\n") {
		t.Errorf("the body of the second draft %q; want the address of the web of the run", body)
	}
}

func TestTheReadyAndMergeRuleBlocksEachCaseAndNothingIsMarkedOrMerged(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name string
		do   func(*pullShip)
		want string
	}{
		{"no status", func(*pullShip) {}, "layup/gates: no status"},
		{"three of the four", func(p *pullShip) {
			p.green(p.head)
			p.w.fake.mu.Lock()
			p.w.fake.statuses[p.head] = p.w.fake.statuses[p.head][1:] // the newest is layup/verify
			p.w.fake.mu.Unlock()
		}, "layup/verify: no status"},
		{"a status of another source", func(p *pullShip) { p.green(p.head); p.foreign(p.head, SpecContext, "success") }, "layup/spec: a status of another source is newer"},
		{"an error of the App", func(p *pullShip) {
			p.green(p.head)
			if _, err := p.s.PostStatus(ctx, p.head, Status{GatesContext, "error", "x"}); err != nil {
				p.t.Fatal(err)
			}
		}, "layup/gates: error"},
		{"no check run", func(p *pullShip) { p.green(p.head); p.setRuns(p.head) }, "job build: no check run"},
		{"a neutral job", func(p *pullShip) { p.green(p.head); p.setRuns(p.head, actionsRun("completed", "neutral", 15368)) }, "job build: completed neutral"},
		{"a skipped job", func(p *pullShip) { p.green(p.head); p.setRuns(p.head, actionsRun("completed", "skipped", 15368)) }, "job build: completed skipped"},
		{"a success of another app", func(p *pullShip) { p.green(p.head); p.setRuns(p.head, actionsRun("completed", "success", 42)) }, "job build: no check run"},
		{"a ruleset with no job of GitHub Actions", func(p *pullShip) {
			p.green(p.head)
			p.w.fake.mu.Lock()
			p.w.fake.rules["main"] = "[]"
			p.w.fake.mu.Unlock()
		}, "no required job of GitHub Actions"},
	} {
		p := newPullShip(t)
		c.do(p)
		before := p.mainSHA()
		merge, r, err := p.s.MergeTask(ctx, "T-ab12", p.head)
		if err != nil || merge != "" || r.Ready || r.Reason != c.want {
			t.Errorf("%s: MergeTask = %q, %+v, %v; want not ready, %q", c.name, merge, r, err, c.want)
		}
		reqs := p.requests()
		if !p.w.fake.pullOf(1).draft || p.mainSHA() != before || countOf(reqs, "POST /graphql") != 0 || countOf(reqs, "PUT ") != 0 || len(p.merges()) != 0 {
			t.Errorf("%s: draft %v, main moved %v, marks %d, merges %d, events %v; want a draft, no mark, no merge, no event", c.name, p.w.fake.pullOf(1).draft, p.mainSHA() != before, countOf(reqs, "POST /graphql"), countOf(reqs, "PUT "), p.merges())
		}
		if mr, err := p.s.MarkReady(ctx, "T-ab12", p.head); err != nil || mr.Ready || !p.w.fake.pullOf(1).draft {
			t.Errorf("%s: MarkReady = %+v, %v with draft %v; want not ready and still a draft", c.name, mr, err, p.w.fake.pullOf(1).draft)
		}
	}
}

func TestTheFourSuccessesMarkTheDraftReadyAndMergeAtTheHeadSHA(t *testing.T) {
	ctx := context.Background()
	p := newPullShip(t)
	base := p.mainSHA()
	p.green(p.head)
	merge, r, err := p.s.MergeTask(ctx, "T-ab12", p.head)
	if err != nil || !r.Ready || merge == "" {
		t.Fatalf("MergeTask = %q, %+v, %v; want a merge", merge, r, err)
	}
	reqs := p.requests()
	mark, put := slices.Index(reqs, "POST /graphql"), slices.Index(reqs, "PUT "+pullsPath+"/1/merge")
	if mark < 0 || put < mark {
		t.Fatalf("the requests %v; want the mark and then the merge", reqs)
	}
	reread := false
	for _, r := range reqs[mark+1 : put] {
		reread = reread || r == "GET /repos/acme/target/commits/"+p.head+"/statuses"
	}
	if !reread {
		t.Errorf("the requests between the mark and the merge %v; want the statuses of the head read again", reqs[mark+1:put])
	}
	if p.mainSHA() != merge || strings.Fields(gitOut(t, p.w.bare, "log", "-1", "--format=%P", merge))[1] != p.head || strings.Fields(gitOut(t, p.w.bare, "log", "-1", "--format=%P", merge))[0] != base {
		t.Errorf("main is %s; want the merge commit %s with the parents %s and %s", p.mainSHA(), merge, base, p.head)
	}
	if len(p.merges()) != 1 || !strings.HasSuffix(p.merges()[0], merge+" 1") {
		t.Errorf("the events merged %v; want one, with the merge commit and pull request 1", p.merges())
	}
}

// A status that changes after the mark blocks the merge: the statuses are read again just before it.
func TestAStatusThatChangesAfterTheMarkBlocksTheMerge(t *testing.T) {
	p := newPullShip(t)
	p.green(p.head)
	p.w.fake.pulls.afterMark = func() {
		p.w.fake.statuses[p.head] = append([]map[string]any{{"context": VerifyContext, "state": "failure", "description": "x", "creator": map[string]any{"id": 9001, "login": "layup-agent[bot]"}}}, p.w.fake.statuses[p.head]...)
	}
	merge, r, err := p.s.MergeTask(context.Background(), "T-ab12", p.head)
	if err != nil || merge != "" || r.Ready || r.Reason != "layup/verify: failure" {
		t.Errorf("MergeTask = %q, %+v, %v; want not ready, layup/verify: failure", merge, r, err)
	}
	if countOf(p.requests(), "PUT ") != 0 || len(p.merges()) != 0 {
		t.Errorf("merges %d, events %v; want none", countOf(p.requests(), "PUT "), p.merges())
	}
}

// The head moves after the last read: the forge answers 409. The call reads the pull request, then the statuses of the new head,
// sends no second merge, and writes no event; the next call merges the new head (the plan review of #363, condition 6).
func TestAHeadThatMovesAfterTheLastReadIsRefusedAndTheNewHeadIsReadBeforeAnotherCall(t *testing.T) {
	ctx := context.Background()
	p := newPullShip(t)
	p.attempt(2)
	newHead := p.bound(p.dev(2, "claude", map[string]string{"y.go": "package x\n"}, okResult))
	p.green(p.head)
	p.green(newHead)
	p.w.fake.pulls.afterMark = func() { p.w.fake.pulls.list[0].head = newHead }
	merge, r, err := p.s.MergeTask(ctx, "T-ab12", p.head)
	if err != nil || merge != "" || r.Ready || !strings.HasPrefix(r.Reason, "the head moved to "+newHead) {
		t.Fatalf("MergeTask = %q, %+v, %v; want the head moved to %s", merge, r, err, newHead)
	}
	reqs := p.requests()
	put := slices.Index(reqs, "PUT "+pullsPath+"/1/merge")
	if put < 0 || countOf(reqs, "PUT ") != 1 {
		t.Fatalf("the requests %v; want exactly one merge", reqs)
	}
	after := reqs[put+1:]
	if len(after) < 2 || after[0] != "GET "+pullsPath+"/1" || !slices.Contains(after, "GET /repos/acme/target/commits/"+newHead+"/statuses") || slices.Contains(after, "GET /repos/acme/target/commits/"+p.head+"/statuses") {
		t.Errorf("the requests after the refused merge %v; want the pull request, then the statuses of %s, and nothing else", after, newHead)
	}
	if len(p.merges()) != 0 || p.w.fake.pullOf(1).merged {
		t.Errorf("events %v, merged %v; want none after a refusal", p.merges(), p.w.fake.pullOf(1).merged)
	}
	p.w.fake.pulls.afterMark = nil
	merge, r, err = p.s.MergeTask(ctx, "T-ab12", newHead)
	if err != nil || !r.Ready || merge == "" || len(p.merges()) != 1 || !strings.HasPrefix(p.merges()[0], "2 ") {
		t.Errorf("the next call = %q, %+v, %v with the events %v; want the new head merged and recorded in attempt 2", merge, r, err, p.merges())
	}
}

// The merge that was made and not recorded (the process was killed, or the lease was lost, between the call and the records commit):
// the tree check passes it with its line, and passes it with none once the event is written. A commit pushed to the branch in any
// other way, and the setup head after the merge, fail it.
func TestAMergeNotYetRecordedPassesTheTreeCheck(t *testing.T) {
	ctx := context.Background()
	p := newPullShip(t)
	base := p.mainSHA()
	p.green(p.head)
	// The merge is made through the forge, and the record of it is not.
	if r, err := p.s.MarkReady(ctx, "T-ab12", p.head); err != nil || !r.Ready {
		t.Fatalf("MarkReady = %+v, %v", r, err)
	}
	merge, err := p.s.Forge.Merge(ctx, 1, p.head)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	gitOut(t, p.r.Clone, "fetch", "-q", p.w.bare, "main")
	st := &state{cfg: p.w.config("0.1.0-dev"), store: p.r.Store.(*recordsStore), target: p.r.Clone, defaultBranch: "main", botID: 9001}
	tree := func(head string) (string, string, []string) {
		f := newCheckFake()
		f.lastForTest(base)
		f.repo = forge.Repository{DefaultBranch: "main", Heads: map[string]string{"main": head}}
		j := f.job()
		j.merges, j.parents = st.treeSeams(ctx)
		j.do()
		row := f.rowsOf(t)[0]
		return row[1], row[2], f.notes
	}
	if state, detail, notes := tree(merge); state != "pass" || detail != "" || !slices.Equal(notes, []string{"tree: " + merge + " merge of T-ab12, not yet recorded"}) {
		t.Errorf("a merge not yet recorded: %s %q %q; want a pass with no detail and the line", state, detail, notes)
	}
	pr, _, err := p.s.Forge.PullRequest(ctx, forge.PullRequestQuery{Number: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err := p.s.RecordMerged(ctx, "T-ab12", pr); err != nil {
		t.Fatalf("RecordMerged: %v", err)
	}
	if state, detail, notes := tree(merge); state != "pass" || detail != "" || len(notes) != 0 {
		t.Errorf("the merge once recorded: %s %q %q; want a pass and no line", state, detail, notes)
	}
	// A commit that reached the branch in another way, and the setup head after the first merge, fail.
	push := childCommit(t, p.r.Clone, merge, map[string]string{"z.go": "package x\n"})
	if state, detail, _ := tree(push); state != "fail" || detail != push+" is not a recorded merge" {
		t.Errorf("a direct push: %s %q; want the fail %q", state, detail, push+" is not a recorded merge")
	}
	if state, detail, _ := tree(base); state != "fail" || detail != base+" is not a recorded merge" {
		t.Errorf("the setup head after the merge: %s %q; want the fail", state, detail)
	}
}

// The demo of row 66g: the draft whose body links the issue and the records commit, the close-out that carries layup/verify, the
// four successes and the job, the merge at the head SHA with its event, and the tree check that passes the merge before and after
// its event is recorded.
func TestAPullRequestFromDraftToMergeAndTheTreeAfterIt(t *testing.T) {
	ctx := context.Background()
	s := newShip(t)
	id, err := s.closeOut(s.V, s.in, closing())
	if err != nil {
		t.Fatal(err)
	}
	head := s.bound(id)
	if head == "" {
		t.Fatalf("the close-out session %s was not bound", id)
	}
	prod := s.prodBot()
	pr, err := prod.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 1, Type: "docs", Description: "close out the task"})
	if err != nil || pr.HeadSHA != head {
		t.Fatalf("OpenPullRequest = %+v, %v; want a draft at the close-out head %s", pr, err, head)
	}
	if body := s.w.fake.pulls.list[0].body; !strings.HasPrefix(body, "Closes #3\n\n") || !strings.HasSuffix(body, s.a.records) {
		t.Errorf("the body %q; want Closes #3 and the records commit %s", body, s.a.records)
	}
	if _, err := prod.CarryVerify(ctx, s.in, id); err != nil {
		t.Fatalf("CarryVerify: %v", err)
	}
	for _, c := range []string{GatesContext, SpecContext, RulesContext} {
		if _, err := prod.PostStatus(ctx, head, Status{c, "success", "pass"}); err != nil {
			t.Fatal(err)
		}
	}
	s.w.fake.mu.Lock()
	s.w.fake.rules = map[string]string{"main": requiredBuild}
	s.w.fake.checkRuns = map[string][]map[string]any{head: {actionsRun("completed", "success", 15368)}}
	s.w.fake.mu.Unlock()
	base := strings.TrimSpace(gitOut(t, s.w.bare, "rev-parse", "refs/heads/main"))
	merge, r, err := prod.MergeTask(ctx, "T-ab12", head)
	if err != nil || !r.Ready || merge == "" {
		t.Fatalf("MergeTask = %q, %+v, %v", merge, r, err)
	}
	if got := strings.Fields(gitOut(t, s.w.bare, "log", "-1", "--format=%P", merge)); len(got) != 2 || got[0] != base || got[1] != head {
		t.Errorf("the merge commit has the parents %v; want %s and %s", got, base, head)
	}
	last := s.events()[len(s.events())-1]
	if last[1] != "merged" || last[5] != merge || last[6] != "1" {
		t.Errorf("the last event %q; want merged, %s, pull request 1", last, merge)
	}
	gitOut(t, s.r.Clone, "fetch", "-q", s.w.bare, "main")
	st := &state{cfg: s.w.config("0.1.0-dev"), store: s.r.Store.(*recordsStore), target: s.r.Clone, defaultBranch: "main", botID: 9001}
	f := newCheckFake()
	f.lastForTest(base)
	f.repo = forge.Repository{DefaultBranch: "main", Heads: map[string]string{"main": merge}}
	j := f.job()
	j.merges, j.parents = st.treeSeams(ctx)
	if got := j.do(); got.Result != "done" || len(f.notes) != 0 {
		t.Errorf("the tree check after the merge: %+v with the notes %q; want done with no line, the event is recorded", got, f.notes)
	}
}
