//go:build integration

package run

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// The link of the review record through the production Sessions, the adapter, the HTTP fixture and real git (row 66q, task
// T-08rj, #413): a developer's head bound, the draft opened by OpenPullRequest, a verifier's review record posted and its
// handoff recorded, then the edit of the body.
type linkShip struct {
	*loopRun
	prod     *Sessions
	head     string
	dev      string
	opening  string // the body of the draft, as the forge holds it
	verifier string
	posted   Posted
}

func newLinkShip(t *testing.T) *linkShip {
	t.Helper()
	l := newLoop(t, nil)
	l.r.Bot = 9001
	ctx := context.Background()
	dev := l.dev(1, "claude", map[string]string{"x.go": "package x\n"}, okResult)
	head := l.bound(dev)
	if head == "" {
		t.Fatal("the developer's session was not bound")
	}
	prod := l.prodBot()
	if _, err := prod.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"}); err != nil {
		t.Fatalf("OpenPullRequest: %v", err)
	}
	// A session starts from a clone of the run's main alone, so the reviewed head is moved onto it, as the loop would before a verifier.
	gitOut(t, l.r.Clone, "update-ref", "refs/heads/main", head)
	ls := &linkShip{loopRun: l, prod: prod, head: head, dev: dev, opening: l.w.fake.pulls.list[0].body}
	ls.verifier = ls.session("verification", record(head, "nothing material in scope"))
	var err error
	if ls.posted, err = l.r.PostComment(ctx, "T-ab12", ls.verifier, "plan.md"); err != nil {
		t.Fatalf("PostComment: %v", err)
	}
	v, err := l.r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: ls.verifier, From: "verifier", To: "close-out", File: "plan.md", Kind: "review-record", Items: map[string]bool{"nothing-material": true}})
	if err != nil || !v.Valid {
		t.Fatalf("RecordHandoff = %+v, %v", v, err)
	}
	return ls
}

// session runs a session of the step on the reviewed head, which writes the text as its comment and changes nothing, so it
// binds the head again at the snapshot of its own start.
func (l *linkShip) session(step, text string) string {
	l.t.Helper()
	spec := harness(l.t, l.a, planWork("plan.md", text, "comment\t1\tresult/plan.md\t"+commentSum(text)+"\t—"))
	spec.Role, spec.Base, spec.Pair.Harness = roleOfStep(step), l.head, "devin"
	spec.Records = strings.TrimSpace(gitOut(l.t, l.w.bare, "rev-parse", "layup-records")) // the snapshot at this session's start, later than the developer's
	id, err := l.r.TaskSession(context.Background(), spec)
	if err != nil {
		l.t.Fatalf("TaskSession %s: %q, %v", step, id, err)
	}
	return id
}

func (l *linkShip) records(session string) string {
	l.t.Helper()
	row, err := l.r.startRowOf(context.Background(), session)
	if err != nil {
		l.t.Fatal(err)
	}
	return row[sessionRecords]
}

func (l *linkShip) body() string { return l.w.fake.pulls.list[0].body }

func (l *linkShip) want() string {
	return PullBody(3, l.prod.Web, l.prod.Target, l.records(l.dev), "", reviewLine(l.prod.Web, l.prod.Target, 3, l.posted.ID))
}

// I1, the demo: one PATCH after the comment's POST, the exact bytes, and no second edit on a second call.
func TestTheReviewRecordIsLinkedFromTheBodyAfterItsPost(t *testing.T) {
	t.Parallel()
	l := newLinkShip(t)
	if l.records(l.dev) == l.records(l.verifier) {
		t.Fatal("the developer and the verifier share a snapshot: the case of two snapshots at one head is not made")
	}
	if l.opening != PullBody(3, l.prod.Web, l.prod.Target, l.records(l.dev)) {
		t.Fatalf("the draft's body = %q; want the two lines of the developer's snapshot", l.opening)
	}
	edited, err := l.prod.LinkReviewRecord(context.Background(), "T-ab12", l.verifier)
	if err != nil || !edited {
		t.Fatalf("LinkReviewRecord = %v, %v; want true, nil", edited, err)
	}
	if got, want := l.body(), l.want(); got != want || !strings.HasSuffix(got, fmt.Sprintf("#issuecomment-%d", l.posted.ID)) {
		t.Errorf("the body = %q; want %q", got, want)
	}
	if !strings.Contains(l.body(), "file://"+l.w.web+"/acme/target/issues/3#issuecomment-") { // the web address is the forge register's, set by state.sessions()
		t.Errorf("the body %q does not link the comment at the web address of the forge register", l.body())
	}
	reqs := l.w.fake.reqs
	post, patch := -1, -1
	for i, r := range reqs {
		switch {
		case strings.HasPrefix(r, "POST /repos/acme/target/issues/3/comments"):
			post = i
		case r == "PATCH /repos/acme/target/pulls/1":
			patch = i
		}
	}
	if countOf(reqs, "PATCH ") != 1 || patch < post || post < 0 {
		t.Errorf("the requests %v; want one PATCH of the pull request after the POST of the comment", reqs)
	}
	if edited, err := l.prod.LinkReviewRecord(context.Background(), "T-ab12", l.verifier); err != nil || edited || countOf(l.w.fake.reqs, "PATCH ") != 1 {
		t.Errorf("a second call = %v, %v with %d PATCH; want false, nil and the one", edited, err, countOf(l.w.fake.reqs, "PATCH "))
	}
}

// I2. The link survives the close-out and a restart: a close-out session binds the head at a third snapshot, a new production
// Sessions edits nothing, and the records line of the body is still the opening's.
func TestTheLinkSurvivesTheCloseOutAndARestart(t *testing.T) {
	t.Parallel()
	l := newLinkShip(t)
	if _, err := l.prod.LinkReviewRecord(context.Background(), "T-ab12", l.verifier); err != nil {
		t.Fatal(err)
	}
	closeOut := l.session("close-out", "the close-out\n")
	if l.records(closeOut) == l.records(l.dev) {
		t.Fatal("the close-out shares the developer's snapshot")
	}
	restarted := l.prodBot()
	if edited, err := restarted.LinkReviewRecord(context.Background(), "T-ab12", l.verifier); err != nil || edited || countOf(l.w.fake.reqs, "PATCH ") != 1 {
		t.Errorf("after the close-out and a restart: %v, %v with %d PATCH; want false, nil and the one", edited, err, countOf(l.w.fake.reqs, "PATCH "))
	}
	first := strings.SplitN(l.opening, "\n", 3)[2]
	if !strings.Contains(l.body(), first) {
		t.Errorf("the body %q lost the records line of the opening %q", l.body(), first)
	}
}

// I3. An edit that the forge refuses is a forge error of EditBody and the body stays.
func TestAnEditRefusedByTheForgeIsAForgeErrorOfEditBody(t *testing.T) {
	t.Parallel()
	l := newLinkShip(t)
	l.w.fake.pulls.denyPatch = 403
	edited, err := l.prod.LinkReviewRecord(context.Background(), "T-ab12", l.verifier)
	if err == nil || edited || !strings.Contains(err.Error(), "EditBody: status 403") {
		t.Errorf("LinkReviewRecord = %v, %v; want EditBody: status 403", edited, err)
	}
	if l.body() != l.opening {
		t.Errorf("the body changed to %q", l.body())
	}
}
