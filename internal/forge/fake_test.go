package forge

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

const (
	sha1a = "1111111111111111111111111111111111111111"
	sha2b = "2222222222222222222222222222222222222222"
)

func TestFakeDeclaresTheSixCapabilities(t *testing.T) {
	if got := Missing(NewFake().Capabilities()); len(got) != 0 {
		t.Errorf("the stand-in does not declare %v", got)
	}
}

func TestFakeOpensADraftMarksItReadyAndMergesItAtItsHead(t *testing.T) {
	f := NewFake()
	f.SetHead("task/T-aaaa/1", sha1a)
	ctx := context.Background()
	pr, err := f.OpenDraft(ctx, "task/T-aaaa/1", "main", "feat: T-aaaa x", "Closes #1")
	if err != nil || !pr.Draft || pr.State != "open" || pr.Number != 1 || pr.NodeID == "" || pr.HeadSHA != sha1a {
		t.Fatalf("OpenDraft = %+v, %v", pr, err)
	}
	if _, err := f.Merge(ctx, pr.Number, sha1a); err == nil {
		t.Errorf("a draft merged")
	} else if e := new(Error); !errors.As(err, &e) || e.Status != 405 {
		t.Errorf("Merge of a draft = %v, want status 405", err)
	}
	rdy, err := f.MarkReady(ctx, pr.NodeID)
	if err != nil || rdy.Draft || rdy.Number != pr.Number {
		t.Fatalf("MarkReady = %+v, %v", rdy, err)
	}
	sha, err := f.Merge(ctx, pr.Number, sha1a)
	if err != nil || sha == "" || sha == sha1a {
		t.Fatalf("Merge = %q, %v", sha, err)
	}
	got, found, err := f.PullRequest(ctx, PullRequestQuery{Number: pr.Number})
	if err != nil || !found || !got.Merged || got.State != "closed" || got.MergeSHA != sha {
		t.Errorf("PullRequest after the merge = %+v, %v, %v", got, found, err)
	}
	if _, err := f.Merge(ctx, pr.Number, sha1a); err == nil {
		t.Errorf("a merged pull request merged again")
	}
}

func TestFakeRefusesAMergeOfAHeadThatMoved(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	pr, _ := f.OpenDraft(ctx, "b", "main", "t", "")
	f.MarkReady(ctx, pr.NodeID)
	f.MoveHead(pr.Number, sha2b)
	_, err := f.Merge(ctx, pr.Number, sha1a)
	if e := new(Error); !errors.As(err, &e) || e.Status != 409 {
		t.Fatalf("Merge at the old head = %v, want status 409", err)
	}
	if got, _, _ := f.PullRequest(ctx, PullRequestQuery{Number: pr.Number}); got.Merged || got.HeadSHA != sha2b {
		t.Errorf("after the refused merge: %+v", got)
	}
	if _, err := f.Merge(ctx, pr.Number, sha2b); err != nil {
		t.Errorf("Merge at the new head = %v", err)
	}
}

func TestFakeFindsAPullRequestByItsBranchAndGivesNoneWithoutAnError(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	if _, found, err := f.PullRequest(ctx, PullRequestQuery{Branch: "task/T-aaaa/1"}); found || err != nil {
		t.Errorf("no pull request: found = %v, err = %v", found, err)
	}
	pr, _ := f.OpenDraft(ctx, "task/T-aaaa/1", "main", "t", "")
	got, found, err := f.PullRequest(ctx, PullRequestQuery{Branch: "task/T-aaaa/1"})
	if err != nil || !found || got.Number != pr.Number {
		t.Errorf("by branch = %+v, %v, %v", got, found, err)
	}
	if _, _, err := f.PullRequest(ctx, PullRequestQuery{Number: 99}); err == nil {
		t.Errorf("a number that is no pull request gave no error")
	}
	if _, _, err := f.PullRequest(ctx, PullRequestQuery{}); err == nil {
		t.Errorf("an empty query gave no error")
	}
	if _, _, err := f.PullRequest(ctx, PullRequestQuery{Number: 1, Branch: "b"}); err == nil {
		t.Errorf("a query with both fields gave no error")
	}
}

func TestFakeGivesTheStatusesNewestFirstWithTheirCreatorAndTheCheckRuns(t *testing.T) {
	f := NewFake()
	f.AppBot = Status{CreatorID: 77, CreatorLogin: "layup-agent[bot]"}
	ctx := context.Background()
	if err := f.SetStatus(ctx, sha1a, "layup/gates", "pending", "running"); err != nil {
		t.Fatal(err)
	}
	if err := f.SetStatus(ctx, sha1a, "layup/gates", "success", "pass"); err != nil {
		t.Fatal(err)
	}
	f.AddStatus(sha1a, Status{Context: "layup/gates", State: "success", CreatorID: 5, CreatorLogin: "someone"})
	f.SetCheckRuns(sha1a, []CheckRun{{Name: "layout", Status: "completed", Conclusion: "success", AppID: 15368, AppSlug: "github-actions"}})
	got, err := f.Statuses(ctx, sha1a)
	if err != nil {
		t.Fatal(err)
	}
	want := []Status{
		{Context: "layup/gates", State: "success", CreatorID: 5, CreatorLogin: "someone"},
		{Context: "layup/gates", State: "success", Description: "pass", CreatorID: 77, CreatorLogin: "layup-agent[bot]"},
		{Context: "layup/gates", State: "pending", Description: "running", CreatorID: 77, CreatorLogin: "layup-agent[bot]"},
	}
	if !reflect.DeepEqual(got.Statuses, want) {
		t.Errorf("Statuses = %+v, want %+v", got.Statuses, want)
	}
	if len(got.CheckRuns) != 1 || got.CheckRuns[0].Name != "layout" {
		t.Errorf("CheckRuns = %+v", got.CheckRuns)
	}
	if other, _ := f.Statuses(ctx, sha2b); len(other.Statuses) != 0 || len(other.CheckRuns) != 0 {
		t.Errorf("another head has %+v: a status is bound to its SHA", other)
	}
}

func TestFakeRefusesAStatusThatTheCallRefuses(t *testing.T) {
	f := NewFake()
	for _, c := range []struct{ sha, state string }{{sha1a, "ok"}, {"abc", "success"}, {"", "success"}} {
		if err := f.SetStatus(context.Background(), c.sha, "layup/gates", c.state, "x"); err == nil {
			t.Errorf("SetStatus(%q, %q) gave no error", c.sha, c.state)
		}
	}
}

func TestFakeFailsTheNamedCall(t *testing.T) {
	f := NewFake()
	f.Fail = "Merge"
	ctx := context.Background()
	pr, _ := f.OpenDraft(ctx, "b", "main", "t", "")
	f.MarkReady(ctx, pr.NodeID)
	_, err := f.Merge(ctx, pr.Number, "")
	if e := new(Error); !errors.As(err, &e) || e.Call != "Merge" {
		t.Errorf("Merge with Fail = %v", err)
	}
}
