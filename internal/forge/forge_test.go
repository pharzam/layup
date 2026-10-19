package forge

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

// standIn is a stand-in adapter of a test: it declares the capabilities it is
// given and makes no call.
type standIn struct{ declared []Capability }

func (s standIn) Capabilities() []Capability { return s.declared }
func (standIn) Installation(context.Context) (Installation, error) {
	return Installation{}, nil
}
func (standIn) Token(context.Context) (Token, error)           { return Token{}, nil }
func (standIn) Repository(context.Context) (Repository, error) { return Repository{}, nil }
func (standIn) UserID(context.Context, string) (int64, error)  { return 0, nil }
func (standIn) OpenIssue(context.Context, string, string) (int, error) {
	return 0, nil
}
func (standIn) Comments(context.Context, int) ([]Comment, error)          { return nil, nil }
func (standIn) Comment(context.Context, int, string) (int64, error)       { return 0, nil }
func (standIn) EffectiveRules(context.Context, string) ([]Rule, error)    { return nil, nil }
func (standIn) Activity(context.Context, string) ([]ActivityEntry, error) { return nil, nil }

func (standIn) OpenDraft(context.Context, string, string, string, string) (PullRequest, error) {
	return PullRequest{}, nil
}
func (standIn) MarkReady(context.Context, string) (PullRequest, error) { return PullRequest{}, nil }
func (standIn) SetStatus(context.Context, string, string, string, string) error {
	return nil
}
func (standIn) Statuses(context.Context, string) (HeadChecks, error) { return HeadChecks{}, nil }
func (standIn) PullRequest(context.Context, PullRequestQuery) (PullRequest, bool, error) {
	return PullRequest{}, false, nil
}
func (standIn) Merge(context.Context, int, string) (string, error) { return "", nil }
func (standIn) EditBody(context.Context, int, string) (PullRequest, error) {
	return PullRequest{}, nil
}
func (standIn) ClosePullRequest(context.Context, int) (PullRequest, error) {
	return PullRequest{}, nil
}
func (standIn) CloseIssue(context.Context, int) error { return nil }

var _ Forge = standIn{}

func TestMissingNamesEachUndeclaredCapability(t *testing.T) {
	if len(Capabilities) != 6 {
		t.Fatalf("Capabilities has %d, want the six of forge.md", len(Capabilities))
	}
	five := append([]Capability{}, Capabilities[:2]...)
	five = append(five, Capabilities[3:]...)
	if got := Missing(standIn{five}.Capabilities()); !reflect.DeepEqual(got, []Capability{Capabilities[2]}) {
		t.Errorf("a stand-in that declares five: Missing = %v, want [%s]", got, Capabilities[2])
	}
	if got := Missing(standIn{Capabilities}.Capabilities()); len(got) != 0 {
		t.Errorf("a stand-in that declares six: Missing = %v, want none", got)
	}
	if got := Missing(nil); !reflect.DeepEqual(got, Capabilities) {
		t.Errorf("no declaration: Missing = %v, want the six in order", got)
	}
}

func TestCheckPermissionsNamesEachMissingPermission(t *testing.T) {
	for _, c := range []struct {
		name string
		got  map[string]string
		want []string
	}{
		{"the full set", map[string]string{"contents": "write", "issues": "write", "metadata": "read"}, nil},
		{"write holds read", map[string]string{"contents": "write", "issues": "write", "metadata": "write"}, nil},
		{"more than M2a uses", map[string]string{"contents": "write", "issues": "write", "metadata": "read", "statuses": "write"}, nil},
		{"no issues", map[string]string{"contents": "write", "metadata": "read"}, []string{"issues: write"}},
		{"contents read, not write", map[string]string{"contents": "read", "issues": "write", "metadata": "read"}, []string{"contents: write"}},
		{"none", nil, []string{"contents: write", "issues: write", "metadata: read"}},
		{"an unknown level", map[string]string{"contents": "admin", "issues": "write", "metadata": "read"}, []string{"contents: write"}},
	} {
		if got := CheckPermissions(c.got); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: CheckPermissions = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestCheckLoopPermissionsNamesEachMissingPermission(t *testing.T) {
	for _, c := range []struct {
		name string
		got  map[string]string
		want []string
	}{
		{"the full set", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "read"}, nil},
		{"write holds read", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "write"}, nil},
		{"more than the loop uses", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "read", "contents": "write"}, nil},
		{"no checks", map[string]string{"pull_requests": "write", "statuses": "write"}, []string{"checks: read"}},
		{"statuses read, not write", map[string]string{"pull_requests": "write", "statuses": "read", "checks": "read"}, []string{"statuses: write"}},
		{"pull requests read, not write", map[string]string{"pull_requests": "read", "statuses": "write", "checks": "read"}, []string{"pull_requests: write"}},
		{"none", nil, []string{"pull_requests: write", "statuses: write", "checks: read"}},
		{"an unknown level", map[string]string{"pull_requests": "admin", "statuses": "write", "checks": "read"}, []string{"pull_requests: write"}},
	} {
		if got := CheckLoopPermissions(c.got); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: CheckLoopPermissions = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestTheLoopPermissionsLeaveTheM2aOnesAlone(t *testing.T) {
	if got := CheckPermissions(map[string]string{"contents": "write", "issues": "write", "metadata": "read"}); got != nil {
		t.Errorf("CheckPermissions with the M2a set = %v, want none: the loop's three are not M2a's", got)
	}
}

func TestValidStatusRefusesAWrongStateOrATooLongDescription(t *testing.T) {
	long := strings.Repeat("é", MaxDescription) // 140 characters, 280 bytes
	for _, c := range []struct {
		state, desc string
		ok          bool
	}{
		{"success", "pass", true},
		{"failure", "", true},
		{"pending", "waits for its turn", true},
		{"error", long, true},
		{"success", long + "x", false},
		{"ok", "pass", false},
		{"", "pass", false},
		{"Success", "pass", false},
	} {
		err := ValidStatus(c.state, c.desc)
		if (err == nil) != c.ok {
			t.Errorf("ValidStatus(%q, %d characters) = %v, ok = %v", c.state, len([]rune(c.desc)), err, c.ok)
		}
	}
	if len(StatusStates) != 4 {
		t.Errorf("StatusStates = %v, want the four of forge.md", StatusStates)
	}
}

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

func TestFakeRefusesWhatTheAdapterRefusesBeforeAnyRequest(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	pr, _ := f.OpenDraft(ctx, "b", "main", "t", "")
	f.MarkReady(ctx, pr.NodeID)
	refused := func(call string, err error) {
		t.Helper()
		if e := new(Error); !errors.As(err, &e) || e.Call != call || e.Status != 0 {
			t.Errorf("%s = %v, want a refusal of that call with the status 0 (a 404 or a 502 is not one)", call, err)
		}
	}
	_, err := f.Statuses(ctx, "abc")
	refused("Statuses", err)
	_, err = f.Merge(ctx, pr.Number, "")
	refused("Merge", err)
	_, err = f.Merge(ctx, 0, sha1a)
	refused("Merge", err)
	_, _, err = f.PullRequest(ctx, PullRequestQuery{Number: -1})
	refused("PullRequest", err)
	f.Fail = "SetStatus"
	refused("SetStatus", f.SetStatus(ctx, "abc", "layup/gates", "success", "x")) // a refusal comes before a failure of the call
	if !FullSHA(sha1a) || FullSHA("") || FullSHA(sha1a[:39]) || FullSHA(sha1a+"0") || FullSHA("ABCDEF"+sha1a[6:]) {
		t.Errorf("FullSHA is not 40 lowercase hexadecimal characters")
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
	_, err := f.Merge(ctx, pr.Number, sha1a)
	if e := new(Error); !errors.As(err, &e) || e.Call != "Merge" || e.Status != 502 {
		t.Errorf("Merge with Fail = %v, want the forge error of status 502 that Fail gives, not a refusal", err)
	}
	f.Fail = ""
	if _, err := f.Merge(ctx, pr.Number, "main"); err == nil {
		t.Errorf("Merge at a SHA that is none gave no error")
	} else if e := new(Error); errors.As(err, &e) && e.Status != 0 {
		t.Errorf("a refusal has the status %d, want 0: it is made before any request", e.Status)
	}
}

func TestFakeKeepsTheBodyOfADraftAndEditsIt(t *testing.T) {
	f := NewFake()
	f.SetHead("task/T-aaaa/1", sha1a)
	ctx := context.Background()
	pr, err := f.OpenDraft(ctx, "task/T-aaaa/1", "main", "t", "Closes #1")
	if err != nil || pr.Body != "Closes #1" {
		t.Fatalf("OpenDraft = %+v, %v, want the body it was given", pr, err)
	}
	byNumber := func() PullRequest {
		t.Helper()
		got, found, err := f.PullRequest(ctx, PullRequestQuery{Number: pr.Number})
		if err != nil || !found {
			t.Fatalf("PullRequest by number = %+v, %v, %v", got, found, err)
		}
		return got
	}
	byBranch := func() PullRequest {
		t.Helper()
		got, found, err := f.PullRequest(ctx, PullRequestQuery{Branch: "task/T-aaaa/1"})
		if err != nil || !found {
			t.Fatalf("PullRequest by branch = %+v, %v, %v", got, found, err)
		}
		return got
	}
	if got := byNumber(); got.Body != "Closes #1" {
		t.Errorf("PullRequest by number: Body = %q; want %q", got.Body, "Closes #1")
	}
	if got := byBranch(); got.Body != "Closes #1" {
		t.Errorf("PullRequest by branch: Body = %q; want %q", got.Body, "Closes #1")
	}
	edited, err := f.EditBody(ctx, pr.Number, "x")
	if err != nil || edited.Number != pr.Number || edited.Body != "x" || edited.HeadSHA != sha1a || !edited.Draft {
		t.Fatalf("EditBody = %+v, %v; want Number %d, Body %q and the rest of the pull request unchanged", edited, err, pr.Number, "x")
	}
	if got := byNumber(); got.Body != "x" {
		t.Errorf("by number after the edit: Body = %q; want %q", got.Body, "x")
	}
	if got := byBranch(); got.Body != "x" {
		t.Errorf("by branch after the edit: Body = %q; want %q", got.Body, "x")
	}
	// the forge edits a closed or merged pull request too: the refusal is the caller's
	f.MarkReady(ctx, pr.NodeID)
	if _, err := f.Merge(ctx, pr.Number, sha1a); err != nil {
		t.Fatal(err)
	}
	if edited, err := f.EditBody(ctx, pr.Number, "after the merge"); err != nil || !edited.Merged || edited.Body != "after the merge" {
		t.Errorf("EditBody of a merged pull request = %+v, %v; want it edited and still merged", edited, err)
	}
}

func TestFakeRefusesAnEditBelow1BeforeItsFailureAndGivesTheForgeErrors(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	f.OpenDraft(ctx, "b", "main", "t", "")
	f.Fail = "EditBody"
	for _, n := range []int{0, -1} {
		_, err := f.EditBody(ctx, n, "x")
		if e := new(Error); !errors.As(err, &e) || e.Call != "EditBody" || e.Status != 0 {
			t.Errorf("EditBody(%d) with Fail = %v; want a refusal of the call EditBody with the status 0 (a 404 or a 502 is not one)", n, err)
		}
	}
	f.Fail = ""
	_, err := f.EditBody(ctx, 99, "x")
	if e := new(Error); !errors.As(err, &e) || e.Call != "EditBody" || e.Status != 404 {
		t.Errorf("EditBody(99) = %v; want a forge error of status 404", err)
	}
	f.Fail = "EditBody"
	_, err = f.EditBody(ctx, 1, "x")
	if e := new(Error); !errors.As(err, &e) || e.Call != "EditBody" || e.Status != 502 {
		t.Errorf("EditBody(1) with Fail = %v; want the forge error of status 502", err)
	}
	f.Fail = ""
	if got, _, _ := f.PullRequest(ctx, PullRequestQuery{Number: 1}); got.Body != "" {
		t.Errorf("a refused or failed edit changed the body to %q", got.Body)
	}
}

func TestFakeCanTakeAnEditAndStillGiveAnError(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	f.OpenDraft(ctx, "b", "main", "t", "old")
	f.FailAfter = "EditBody"
	_, err := f.EditBody(ctx, 1, "new")
	if e := new(Error); !errors.As(err, &e) || e.Call != "EditBody" || e.Status != 502 {
		t.Errorf("EditBody with FailAfter = %v; want the forge error of status 502", err)
	}
	if got, _, _ := f.PullRequest(ctx, PullRequestQuery{Number: 1}); got.Body != "new" {
		t.Errorf("after the error: Body = %q; want %q (the edit took effect)", got.Body, "new")
	}
	f.FailAfter = ""
	if _, err := f.EditBody(ctx, 0, "x"); err == nil {
		t.Errorf("a refusal below 1 gave no error")
	}
}

// The closes (task T-xyws, #423; forge.md, The calls of M2e).

func TestFakeClosesAPullRequestAndReadsItClosedByNumberAndByBranch(t *testing.T) {
	f := NewFake()
	f.SetHead("task/T-aaaa/1", sha1a)
	ctx := context.Background()
	opened, _ := f.OpenDraft(ctx, "task/T-aaaa/1", "main", "feat: T-aaaa x", "Closes #1")
	f.MarkReady(ctx, opened.NodeID)
	pr, err := f.ClosePullRequest(ctx, opened.Number)
	if err != nil || pr.Number != opened.Number || pr.State != "closed" || pr.Merged {
		t.Fatalf("ClosePullRequest(%d) = %+v, %v; want Number %d, State closed, Merged false", opened.Number, pr, err, opened.Number)
	}
	byNumber, found, err := f.PullRequest(ctx, PullRequestQuery{Number: opened.Number})
	if err != nil || !found || byNumber.State != "closed" || byNumber.Merged {
		t.Errorf("PullRequest by number after the close = %+v, %v, %v; want State closed, Merged false", byNumber, found, err)
	}
	byBranch, found, err := f.PullRequest(ctx, PullRequestQuery{Branch: "task/T-aaaa/1"})
	if err != nil || !found || byBranch.State != "closed" {
		t.Errorf("PullRequest by branch after the close = %+v, %v, %v; want State closed", byBranch, found, err)
	}
	again, err := f.ClosePullRequest(ctx, opened.Number) // a repeat is no error
	if err != nil || again.State != "closed" || again.Merged {
		t.Errorf("a second ClosePullRequest = %+v, %v; want State closed, Merged false, no error", again, err)
	}
	if _, err := f.Merge(ctx, opened.Number, sha1a); err == nil {
		t.Errorf("a closed pull request merged")
	}
}

func TestFakeCloseOfAMergedPullRequestChangesNothingAndGivesItMerged(t *testing.T) {
	f := NewFake()
	f.SetHead("b", sha1a)
	ctx := context.Background()
	opened, _ := f.OpenDraft(ctx, "b", "main", "t", "")
	f.MarkReady(ctx, opened.NodeID)
	sha, err := f.Merge(ctx, opened.Number, sha1a)
	if err != nil {
		t.Fatal(err)
	}
	pr, err := f.ClosePullRequest(ctx, opened.Number)
	if err != nil || pr.Number != opened.Number || !pr.Merged || pr.State != "closed" || pr.MergeSHA != sha {
		t.Errorf("ClosePullRequest of a merged pull request = %+v, %v; want Number %d, Merged true, State closed, MergeSHA %q, no error", pr, err, opened.Number, sha)
	}
	got, _, _ := f.PullRequest(ctx, PullRequestQuery{Number: opened.Number})
	if !got.Merged || got.MergeSHA != sha {
		t.Errorf("PullRequest after the close of a merged one = %+v; want Merged true and the merge commit kept", got)
	}
}

func TestFakeClosesAnIssueAsNotPlannedAndKeepsItAfterASecondClose(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	n, _ := f.OpenIssue(ctx, "t", "b")
	if st, reason, ok := f.IssueState(n); st != "open" || reason != "" || !ok {
		t.Errorf("IssueState of an open issue = %q, %q, %v; want open, an empty reason, true", st, reason, ok)
	}
	if err := f.CloseIssue(ctx, n); err != nil {
		t.Fatalf("CloseIssue(%d) = %v", n, err)
	}
	if st, reason, ok := f.IssueState(n); st != "closed" || reason != "not_planned" || !ok {
		t.Errorf("IssueState after the close = %q, %q, %v; want closed, not_planned, true", st, reason, ok)
	}
	if err := f.CloseIssue(ctx, n); err != nil { // a repeat is no error
		t.Errorf("a second CloseIssue = %v, want no error", err)
	}
	if st, reason, ok := f.IssueState(n); st != "closed" || reason != "not_planned" || !ok {
		t.Errorf("IssueState after a second close = %q, %q, %v; want closed, not_planned, true", st, reason, ok)
	}
	for _, bad := range []int{n + 1, 0, -1} {
		if st, reason, ok := f.IssueState(bad); st != "" || reason != "" || ok {
			t.Errorf("IssueState(%d) = %q, %q, %v; want an empty state and reason, false", bad, st, reason, ok)
		}
	}
}

func TestFakeRefusesACloseBelow1BeforeItsFailureAndGivesTheForgeErrors(t *testing.T) {
	f := NewFake()
	ctx := context.Background()
	f.OpenDraft(ctx, "b", "main", "t", "")
	f.OpenIssue(ctx, "t", "b")
	calls := []struct {
		call string
		do   func(n int) error
	}{
		{"ClosePullRequest", func(n int) error { _, err := f.ClosePullRequest(ctx, n); return err }},
		{"CloseIssue", func(n int) error { return f.CloseIssue(ctx, n) }},
	}
	for _, c := range calls {
		f.Fail = c.call
		for _, n := range []int{0, -1} { // a refusal comes before the failure of the call
			if e := new(Error); !errors.As(c.do(n), &e) || e.Call != c.call || e.Status != 0 {
				t.Errorf("%s(%d) with Fail = %v, want a refusal of that call with the status 0", c.call, n, c.do(n))
			}
		}
		if e := new(Error); !errors.As(c.do(1), &e) || e.Call != c.call || e.Status != 502 {
			t.Errorf("%s(1) with Fail = %v, want the forge error of status 502", c.call, c.do(1))
		}
		f.Fail = ""
		if e := new(Error); !errors.As(c.do(99), &e) || e.Call != c.call || e.Status != 404 {
			t.Errorf("%s(99) = %v, want a forge error of status 404", c.call, c.do(99))
		}
	}
}
