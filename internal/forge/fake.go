package forge

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

// Fake is the stand-in forge of the tests of other packages (docs/spec/forge.md,
// The calls of M2e; task T-cr8j, which settles the home of the stand-in that
// issue #373 asked about): the whole interface in memory, with no network and no
// clock. Production code never makes one. A test sets the fields before it starts
// and reads them after; the calls and the helpers hold a lock, so a Fake serves
// goroutines. A call that names Fail gives a forge error of status 502.
type Fake struct {
	mu sync.Mutex

	Fail        string            // the name of the call that gives a forge error
	FailAfter   string            // the name of the call that takes effect and then gives a forge error (EditBody only)
	Declared    []Capability      // nil declares the six
	Permissions map[string]string // of the installation
	Repo        Repository
	IDs         map[string]int64 // the numeric ID of a login
	Rules       map[string][]Rule
	Activities  map[string][]ActivityEntry
	AppBot      Status // the creator (ID and login) that SetStatus gives a status

	issues   int
	closed   map[int]string // an issue closed by CloseIssue, to its reason (task T-xyws)
	comments map[int][]Comment
	prs      []*fakePR
	heads    map[string]string
	statuses map[string][]Status // by SHA, newest first
	checks   map[string][]CheckRun
}

type fakePR struct {
	PullRequest
	branch string
}

// NewFake gives an empty stand-in that declares the six capabilities.
func NewFake() *Fake {
	return &Fake{Permissions: map[string]string{}, heads: map[string]string{}, comments: map[int][]Comment{},
		statuses: map[string][]Status{}, checks: map[string][]CheckRun{}}
}

func (f *Fake) err(call string) error {
	if f.Fail == call {
		return &Error{Call: call, Status: 502, Message: "Bad Gateway"}
	}
	return nil
}

// SetHead sets the head of a branch: a pull request opened from it has that head.
func (f *Fake) SetHead(branch, sha string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.heads[branch] = sha
}

// MoveHead moves the head of a pull request, as a push to its branch would.
func (f *Fake) MoveHead(number int, sha string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if pr := f.pr(number); pr != nil {
		pr.HeadSHA = sha
	}
}

// AddStatus adds a status of another source to a head, as the newest.
func (f *Fake) AddStatus(sha string, s Status) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.statuses[sha] = slices.Insert(f.statuses[sha], 0, s)
}

// SetCheckRuns sets the check runs of a head.
func (f *Fake) SetCheckRuns(sha string, runs []CheckRun) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.checks[sha] = slices.Clone(runs)
}

func (f *Fake) pr(number int) *fakePR {
	if number < 1 || number > len(f.prs) {
		return nil
	}
	return f.prs[number-1]
}

func (f *Fake) Capabilities() []Capability {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Declared != nil {
		return f.Declared
	}
	return Capabilities
}

func (f *Fake) Installation(context.Context) (Installation, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return Installation{ID: 1, Permissions: f.Permissions}, f.err("Installation")
}

func (f *Fake) Token(context.Context) (Token, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return Token{Token: "ghs_fake"}, f.err("Token")
}

func (f *Fake) Repository(context.Context) (Repository, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Repo, f.err("Repository")
}

func (f *Fake) UserID(_ context.Context, login string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.err("UserID"); err != nil {
		return 0, err
	}
	id, ok := f.IDs[login]
	if !ok {
		return 0, &Error{Call: "UserID", Status: 404, Message: "Not Found"}
	}
	return id, nil
}

func (f *Fake) OpenIssue(context.Context, string, string) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.err("OpenIssue"); err != nil {
		return 0, err
	}
	f.issues++
	return f.issues, nil
}

func (f *Fake) Comments(_ context.Context, issue int) ([]Comment, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.comments[issue]), f.err("Comments")
}

func (f *Fake) Comment(_ context.Context, issue int, body string) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.err("Comment"); err != nil {
		return 0, err
	}
	id := int64(len(f.comments[issue]) + 1)
	f.comments[issue] = append(f.comments[issue], Comment{ID: id, Body: body})
	return id, nil
}

func (f *Fake) EffectiveRules(_ context.Context, branch string) ([]Rule, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Rules[branch], f.err("EffectiveRules")
}

func (f *Fake) Activity(_ context.Context, ref string) ([]ActivityEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.Activities[ref], f.err("Activity")
}

func (f *Fake) OpenDraft(_ context.Context, head, base, title, body string) (PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.err("OpenDraft"); err != nil {
		return PullRequest{}, err
	}
	n := len(f.prs) + 1
	pr := &fakePR{branch: head, PullRequest: PullRequest{Number: n, NodeID: fmt.Sprintf("PR_fake%d", n), State: "open", Draft: true, HeadSHA: f.heads[head], Body: body}}
	f.prs = append(f.prs, pr)
	return pr.PullRequest, nil
}

func (f *Fake) MarkReady(_ context.Context, nodeID string) (PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if err := f.err("MarkReady"); err != nil {
		return PullRequest{}, err
	}
	for _, pr := range f.prs {
		if pr.NodeID == nodeID {
			pr.Draft = false
			return pr.PullRequest, nil
		}
	}
	return PullRequest{}, &Error{Call: "MarkReady", Status: 200, Message: "Could not resolve to a node with the global id " + nodeID}
}

func (f *Fake) SetStatus(_ context.Context, sha, context, state, description string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !FullSHA(sha) {
		return &Error{Call: "SetStatus", Message: "the SHA is not a full object ID"}
	}
	if err := ValidStatus(state, description); err != nil {
		return &Error{Call: "SetStatus", Message: err.Error()}
	}
	if err := f.err("SetStatus"); err != nil {
		return err
	}
	s := Status{Context: context, State: state, Description: description, CreatorID: f.AppBot.CreatorID, CreatorLogin: f.AppBot.CreatorLogin}
	f.statuses[sha] = slices.Insert(f.statuses[sha], 0, s)
	return nil
}

func (f *Fake) Statuses(_ context.Context, sha string) (HeadChecks, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if !FullSHA(sha) {
		return HeadChecks{}, &Error{Call: "Statuses", Message: "the SHA is not a full object ID"}
	}
	return HeadChecks{Statuses: slices.Clone(f.statuses[sha]), CheckRuns: slices.Clone(f.checks[sha])}, f.err("Statuses")
}

func (f *Fake) PullRequest(_ context.Context, q PullRequestQuery) (PullRequest, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if (q.Number != 0) == (q.Branch != "") || q.Number < 0 {
		return PullRequest{}, false, &Error{Call: "PullRequest", Message: "the query names a number or a branch, one of the two"}
	}
	if err := f.err("PullRequest"); err != nil {
		return PullRequest{}, false, err
	}
	switch {
	case q.Number != 0:
		pr := f.pr(q.Number)
		if pr == nil {
			return PullRequest{}, false, &Error{Call: "PullRequest", Status: 404, Message: "Not Found"}
		}
		return pr.PullRequest, true, nil
	}
	for i := len(f.prs) - 1; i >= 0; i-- { // the newest first, as the list of the forge
		if f.prs[i].branch == q.Branch {
			return f.prs[i].PullRequest, true, nil
		}
	}
	return PullRequest{}, false, nil
}

func (f *Fake) Merge(_ context.Context, number int, sha string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 || !FullSHA(sha) {
		return "", &Error{Call: "Merge", Message: "the pull request number or the head SHA is not valid"}
	}
	if err := f.err("Merge"); err != nil {
		return "", err
	}
	pr := f.pr(number)
	switch {
	case pr == nil:
		return "", &Error{Call: "Merge", Status: 404, Message: "Not Found"}
	case pr.Merged || pr.State != "open":
		return "", &Error{Call: "Merge", Status: 405, Message: "Pull Request is not mergeable"}
	case pr.Draft:
		return "", &Error{Call: "Merge", Status: 405, Message: "Pull Request is still a draft"}
	case sha != pr.HeadSHA:
		return "", &Error{Call: "Merge", Status: 409, Message: "Head branch was modified. Review and try the merge again."}
	}
	pr.Merged, pr.State = true, "closed"
	pr.MergeSHA = fmt.Sprintf("%040x", 0xabc000+number)
	return pr.MergeSHA, nil
}

// EditBody replaces the body of a pull request, also of a closed or merged one, as
// the forge does: a caller that must refuse such an edit reads the state first. A
// number below 1 is refused before Fail is read, as the adapter refuses it before any
// request. FailAfter names the call that takes effect and then gives the forge error,
// which is how a caller meets an edit whose answer was lost.
func (f *Fake) EditBody(_ context.Context, number int, body string) (PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 {
		return PullRequest{}, &Error{Call: "EditBody", Message: "the pull request number is not valid"}
	}
	if err := f.err("EditBody"); err != nil {
		return PullRequest{}, err
	}
	pr := f.pr(number)
	if pr == nil {
		return PullRequest{}, &Error{Call: "EditBody", Status: 404, Message: "Not Found"}
	}
	pr.Body = body
	if f.FailAfter == "EditBody" {
		return PullRequest{}, &Error{Call: "EditBody", Status: 502, Message: "Bad Gateway"}
	}
	return pr.PullRequest, nil
}

// ClosePullRequest closes a pull request. A closed one stays closed, and a merged
// one changes nothing: it is given with Merged true, as the forge answers (forge.md,
// The calls of M2e, The closes; task T-xyws). Number below 1 is refused before the
// failure of the call.
func (f *Fake) ClosePullRequest(_ context.Context, number int) (PullRequest, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 {
		return PullRequest{}, &Error{Call: "ClosePullRequest", Message: "the pull request number is not valid"}
	}
	if err := f.err("ClosePullRequest"); err != nil {
		return PullRequest{}, err
	}
	pr := f.pr(number)
	if pr == nil {
		return PullRequest{}, &Error{Call: "ClosePullRequest", Status: 404, Message: "Not Found"}
	}
	pr.State = "closed"
	return pr.PullRequest, nil
}

// CloseIssue closes an issue as not planned; a repeat changes nothing.
func (f *Fake) CloseIssue(_ context.Context, number int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 {
		return &Error{Call: "CloseIssue", Message: "the issue number is not valid"}
	}
	if err := f.err("CloseIssue"); err != nil {
		return err
	}
	if number > f.issues {
		return &Error{Call: "CloseIssue", Status: 404, Message: "Not Found"}
	}
	if f.closed == nil {
		f.closed = map[int]string{}
	}
	if _, done := f.closed[number]; !done {
		f.closed[number] = "not_planned"
	}
	return nil
}

// IssueState gives the state of an issue and, when closed, its reason: open with an
// empty reason, or closed with the reason that CloseIssue gave; ok is false for a
// number that is none. A test reads what a close changed with it.
func (f *Fake) IssueState(number int) (state, reason string, ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 || number > f.issues {
		return "", "", false
	}
	if r, done := f.closed[number]; done {
		return "closed", r, true
	}
	return "open", "", true
}

var _ Forge = (*Fake)(nil)
