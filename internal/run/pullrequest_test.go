package run

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

// The title is the form of the checked-in pr-title check (docs/spec/loop.md, The pull request; the plan review of #363,
// condition 2): a type of the fixed list, a description of one line that begins with a lowercase letter, has three
// characters or more, ends with no full stop and holds no control character. The task ID is not in it.
func TestTheTitleIsTheFormOfThePRTitleCheck(t *testing.T) {
	for _, c := range []struct {
		name, typ, description string
		ok                     bool
	}{
		{"the control, a title of the form", "feat", "add the draft", true},
		{"each type of the list", "revert", "undo the merge", true},
		{"three characters", "fix", "abc", true},
		{"a type outside the list", "wip", "add the draft", false},
		{"a type in capitals", "Feat", "add the draft", false},
		{"no type", "", "add the draft", false},
		{"a description that begins with a capital", "feat", "Add the draft", false},
		{"a task ID in front, as the earlier text of the specification wrote it", "feat", "T-ab12 add the draft", false},
		{"a description that begins with a space", "feat", " add the draft", false},
		{"two characters", "fix", "ab", false},
		{"a full stop at the end", "feat", "add the draft.", false},
		{"two lines", "feat", "add\nthe draft", false},
		{"a tab", "feat", "add\tthe draft", false},
		{"an escape", "feat", "add\x1b the draft", false},
		{"empty", "feat", "", false},
	} {
		if err := checkTitle(c.typ, c.description); (err == nil) != c.ok {
			t.Errorf("%s: checkTitle(%q, %q) = %v; want ok = %v", c.name, c.typ, c.description, err, c.ok)
		}
	}
}

func TestThePullBodyIsTheTwoLinesAndTakesFurtherLines(t *testing.T) {
	want := "Closes #3\n\nThe records commit this task started from: " + prWeb + "/acme/target/commit/" + recA
	if got := PullBody(3, prWeb, "acme/target", recA); got != want {
		t.Errorf("PullBody = %q; want %q", got, want)
	}
	more := PullBody(3, prWeb, "acme/target", recA, "The review record: https://x/issues/3#issuecomment-9")
	if !strings.HasPrefix(more, want+"\n") || !strings.HasSuffix(more, "issuecomment-9") {
		t.Errorf("PullBody with a line = %q; want the two lines and the line after them", more)
	}
}

func TestTheDraftIsOpenedOnceWithItsTitleAndBody(t *testing.T) {
	w := newPRWorld(t)
	ctx := context.Background()
	pr, err := w.s.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"})
	if err != nil || pr.Number != 1 || !pr.Draft {
		t.Fatalf("OpenPullRequest = %+v, %v; want the draft number 1", pr, err)
	}
	if w.f.title != "feat: add the draft" || w.f.body != "Closes #3\n\nThe records commit this task started from: "+prWeb+"/acme/target/commit/"+recA {
		t.Errorf("the title %q and the body %q", w.f.title, w.f.body)
	}
	if n := w.prOf(t); n != 1 {
		t.Errorf("the register holds pr %d; want 1", n)
	}
	commits := len(w.st.commits)
	again, err := w.s.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"})
	if err != nil || again.Number != 1 || w.f.count("OpenDraft") != 1 || len(w.st.commits) != commits {
		t.Errorf("the second call = %+v, %v with %d drafts and %d new commits; want the same pull request, one draft, no commit", again, err, w.f.count("OpenDraft"), len(w.st.commits)-commits)
	}
}

// A draft that opened and a register that could not be written (a failure that Fenced does not retry) leaves the draft: the next call finds
// it by its branch and writes the number, and opens no second (the plan review of #363, condition 8).
func TestAStopAfterTheDraftOpensNoSecondDraft(t *testing.T) {
	w := newPRWorld(t)
	ctx := context.Background()
	in := PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"}
	w.st.failure = errors.New("the disk is full")
	if _, err := w.s.OpenPullRequest(ctx, in); err == nil || !strings.Contains(err.Error(), "the disk is full") {
		t.Fatalf("the call with a register that cannot be written: %v; want that error", err)
	}
	if w.f.count("OpenDraft") != 1 || w.prOf(t) != 0 {
		t.Fatalf("%d drafts and pr %d after the stop; want one draft and no number yet", w.f.count("OpenDraft"), w.prOf(t))
	}
	w.st.failure = nil
	pr, err := w.s.OpenPullRequest(ctx, in)
	if err != nil || pr.Number != 1 || w.f.count("OpenDraft") != 1 || w.prOf(t) != 1 {
		t.Errorf("the next call = %+v, %v with %d drafts and pr %d; want the first draft found and its number written", pr, err, w.f.count("OpenDraft"), w.prOf(t))
	}
}

// The records commit of the body is the start row of the session bound to the attempt's head: the snapshot of attempt 1 for
// attempt 1 after attempt 2 started, and never the snapshot of a later session.
func TestTheRecordsCommitIsTheBoundSessionsAndNotALaterOnes(t *testing.T) {
	for _, c := range []struct {
		attempt int
		want    string
	}{{1, recA}, {2, recB}} {
		w := newPRWorld(t)
		in := PullInput{Task: "T-ab12", Attempt: c.attempt, Type: "fix", Description: "mend the draft"}
		if _, err := w.s.OpenPullRequest(context.Background(), in); err != nil {
			t.Fatalf("attempt %d: %v", c.attempt, err)
		}
		if !strings.HasSuffix(w.f.body, "/commit/"+c.want) {
			t.Errorf("attempt %d: the body %q; want the link to the records commit %s of its own session", c.attempt, w.f.body, c.want)
		}
	}
}

func TestATitleOffItsFormIsRefusedBeforeAnyCall(t *testing.T) {
	w := newPRWorld(t)
	_, err := w.s.OpenPullRequest(context.Background(), PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "T-ab12 add the draft"})
	if err == nil || !strings.Contains(err.Error(), "title") {
		t.Errorf("a title with the task ID: %v; want an error that names the title", err)
	}
	if len(w.f.log) != 0 || len(w.st.commits) != 0 {
		t.Errorf("the calls %v and %d commits; want none before the refusal", w.f.log, len(w.st.commits))
	}
}

func TestAnAttemptWithNoBoundHeadAndATaskOutsideTheRegisterAreErrorsBeforeAnyCall(t *testing.T) {
	for name, in := range map[string]PullInput{
		"an attempt with no bound event": {Task: "T-ab12", Attempt: 3, Type: "feat", Description: "add the draft"},
		"a task outside the register":    {Task: "T-zz99", Attempt: 1, Type: "feat", Description: "add the draft"},
	} {
		w := newPRWorld(t)
		if pr, err := w.s.OpenPullRequest(context.Background(), in); err == nil {
			t.Errorf("%s: %+v, no error", name, pr)
		}
		if len(w.f.log) != 0 {
			t.Errorf("%s: the calls %v; want none", name, w.f.log)
		}
	}
}

func TestAPullRequestOfTheBranchIsReadBeforeAnythingIsOpened(t *testing.T) {
	ctx := context.Background()
	in := PullInput{Task: "T-ab12", Attempt: 2, Type: "feat", Description: "add the draft"}
	// An open pull request that the register does not know yet: used, its number written, no draft.
	w := newPRWorld(t)
	open, _ := w.f.Fake.OpenDraft(ctx, "task/T-ab12/2", "main", "feat: x", "b")
	if pr, err := w.s.OpenPullRequest(ctx, in); err != nil || pr.Number != open.Number || w.f.count("OpenDraft") != 0 || w.prOf(t) != open.Number {
		t.Errorf("an open pull request: %+v, %v with %d drafts and pr %d; want it used and its number written", pr, err, w.f.count("OpenDraft"), w.prOf(t))
	}
	// A merged one is returned untouched: no draft, no write.
	w = newPRWorld(t)
	merged := forge.PullRequest{Number: 4, State: "closed", Merged: true, HeadSHA: head2, MergeSHA: head3}
	w.f.found = &merged
	if pr, err := w.s.OpenPullRequest(ctx, in); err != nil || pr != merged || w.f.count("OpenDraft") != 0 || len(w.st.commits) != 0 {
		t.Errorf("a merged pull request: %+v, %v with %d drafts and %d commits; want it returned untouched", pr, err, w.f.count("OpenDraft"), len(w.st.commits))
	}
	// A closed one that is not merged is the redo of row 66k: an error, no new draft.
	w = newPRWorld(t)
	closed := forge.PullRequest{Number: 4, State: "closed", HeadSHA: head2}
	w.f.found = &closed
	if pr, err := w.s.OpenPullRequest(ctx, in); err == nil || w.f.count("OpenDraft") != 0 {
		t.Errorf("a closed pull request: %+v, %v with %d drafts; want an error and no draft", pr, err, w.f.count("OpenDraft"))
	}
}

func TestWithPRSetsTheNumberOfOneTaskAndKeepsTheOthers(t *testing.T) {
	a := Task{ID: "T-aa11", Issue: 1, Kind: "build", Requirements: []string{"REQ-005"}}
	b := Task{ID: "T-bb22", Issue: 2, Kind: "build", Requirements: []string{"REQ-005"}}
	got, err := withPR([]Task{a, b}, "T-bb22", 7)
	if err != nil || len(got) != 2 || !reflect.DeepEqual(got[0], a) || got[1].PR != 7 || got[1].ID != "T-bb22" {
		t.Errorf("withPR = %+v, %v; want the second task on pull request 7 and the first as it was", got, err)
	}
	if a.PR != 0 {
		t.Error("withPR changed the table it was given")
	}
	if got, err := withPR([]Task{a, b}, "T-zz99", 7); err == nil {
		t.Errorf("a task outside the register: %+v, no error", got)
	}
}

// The head of the pull request is the attempt's last bound head: a branch that moved after the event bound is refused for a new
// draft and for an open one, and nothing is registered (round 1, finding 2).
func TestABranchHeadThatIsNotTheBoundHeadIsRefusedAndNothingIsRegistered(t *testing.T) {
	ctx := context.Background()
	in := PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"}
	// A new draft, the branch at head3 while the last bound event names head1.
	w := newPRWorld(t)
	w.f.SetHead("task/T-ab12/1", head3)
	if pr, err := w.s.OpenPullRequest(ctx, in); err == nil || !strings.Contains(err.Error(), "moved") || w.prOf(t) != 0 {
		t.Errorf("a new draft at a moved head: %+v, %v with pr %d; want a moved-head error and no number", pr, err, w.prOf(t))
	}
	// An open pull request whose head moved after the bound event.
	w = newPRWorld(t)
	open, _ := w.f.Fake.OpenDraft(ctx, "task/T-ab12/1", "main", "feat: x", "b")
	w.f.MoveHead(open.Number, head3)
	if pr, err := w.s.OpenPullRequest(ctx, in); err == nil || !strings.Contains(err.Error(), "moved") || w.prOf(t) != 0 || len(w.st.commits) != 0 {
		t.Errorf("an open pull request at a moved head: %+v, %v with pr %d; want a moved-head error and no write", pr, err, w.prOf(t))
	}
	// The control: the heads agree, the pull request is used and registered.
	w = newPRWorld(t)
	open, _ = w.f.Fake.OpenDraft(ctx, "task/T-ab12/1", "main", "feat: x", "b")
	if pr, err := w.s.OpenPullRequest(ctx, in); err != nil || pr.Number != open.Number || w.prOf(t) != open.Number {
		t.Errorf("an open pull request at the bound head: %+v, %v with pr %d; want it used and registered", pr, err, w.prOf(t))
	}
}
