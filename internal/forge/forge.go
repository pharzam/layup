package forge

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

// Capability is one of the six capabilities of a forge (docs/spec/forge.md,
// The six capabilities; architecture.md §1).
type Capability string

// The six capabilities, in the order of the table of forge.md.
const (
	IssuesAndComments Capability = "issues and comments"
	PullRequests      Capability = "pull requests with a draft state"
	CommitStatuses    Capability = "commit statuses"
	BranchRules       Capability = "branch rules"
	Activity          Capability = "repository activity"
	AppIdentity       Capability = "App identity"
)

// Capabilities is the six, in the order of the table of forge.md.
var Capabilities = []Capability{IssuesAndComments, PullRequests, CommitStatuses, BranchRules, Activity, AppIdentity}

// Forge is the interface of the calls of M2a (forge.md, The calls of M2a). An
// adapter declares its capabilities with Capabilities, so the step forge of
// layup run checks a stand-in adapter the same way. Comment, a method of the
// capability of issues and comments, comes with M2b (forge.md, The calls of
// M2b; task T-fsjp, O-189).
type Forge interface {
	Capabilities() []Capability
	// Installation gives the installation of the App on the repository.
	Installation(ctx context.Context) (Installation, error)
	// Token gives an installation token with five minutes or more left: the
	// one it holds, or a new one.
	Token(ctx context.Context) (Token, error)
	Repository(ctx context.Context) (Repository, error)
	// UserID gives the numeric ID of a login.
	UserID(ctx context.Context, login string) (int64, error)
	// OpenIssue opens an issue and gives its number.
	OpenIssue(ctx context.Context, title, body string) (int, error)
	// Comments gives each comment of an issue, every page.
	Comments(ctx context.Context, issue int) ([]Comment, error)
	// Comment posts a comment on an issue and gives its ID (M2b: the comment
	// of a session, task T-fsjp, O-189).
	Comment(ctx context.Context, issue int, body string) (int64, error)
	// EffectiveRules gives each active rule that applies to a branch, also to one
	// that does not exist yet (M2d: forge.md, The calls of M2d; task T-jpkc).
	EffectiveRules(ctx context.Context, branch string) ([]Rule, error)
	// Activity gives the activities of a ref, one page (M2d: forge.md, The calls
	// of M2d; task T-jpkc).
	Activity(ctx context.Context, ref string) ([]ActivityEntry, error)
	// The calls of M2e (forge.md, The calls of M2e; task T-cr8j).
	//
	// OpenDraft opens a draft pull request of head onto base and requests no
	// review.
	OpenDraft(ctx context.Context, head, base, title, body string) (PullRequest, error)
	// MarkReady takes a pull request out of the draft state, by its node ID.
	MarkReady(ctx context.Context, nodeID string) (PullRequest, error)
	// SetStatus posts a commit status on a SHA: one of the four states, and a
	// description of 140 characters or fewer.
	SetStatus(ctx context.Context, sha, context, state, description string) error
	// Statuses gives each status and each check run of a head, every page.
	Statuses(ctx context.Context, sha string) (HeadChecks, error)
	// PullRequest gives a pull request by its number or by its branch; for a
	// branch with none, found is false and there is no error.
	PullRequest(ctx context.Context, q PullRequestQuery) (pr PullRequest, found bool, err error)
	// Merge merges a pull request at its head SHA, by a merge commit, and gives
	// the merge commit's SHA; a head that moved is a forge error of status 409.
	Merge(ctx context.Context, number int, sha string) (string, error)
	// EditBody replaces the body of a pull request and gives the pull request; a
	// number below 1 is refused before any request (forge.md, The calls of M2e;
	// task T-08rj).
	EditBody(ctx context.Context, number int, body string) (PullRequest, error)
	// ClosePullRequest closes a pull request by its number and gives it as the
	// forge answers; one already closed, or merged, is no error (forge.md, The
	// calls of M2e; task T-xyws).
	ClosePullRequest(ctx context.Context, number int) (PullRequest, error)
	// CloseIssue closes an issue as not planned; one already closed is no error.
	CloseIssue(ctx context.Context, number int) error
}

// Installation is the installation of the App on the repository.
type Installation struct {
	ID          int64
	Permissions map[string]string // the name of a permission to its level
}

// Token is an installation token and its end time. It is never written to a
// file, a record or a log.
type Token struct {
	Token   string
	Expires time.Time
}

// Repository is what the step forge reads of the repository.
type Repository struct {
	DefaultBranch string
	Visibility    string
	HasCommit     bool              // the list of branches is not empty
	Branches      []string          // the names of the branches, each page
	Heads         map[string]string // the head of each branch, by its name (M2d, task T-jpkc)
}

// Rule is a rule that applies to a branch: its type and its parameters, the
// object that the forge gives, decoded; nil when the forge gives none or an
// empty object (forge.md, The calls of M2d). It also names the ruleset that it is
// of, as the forge does (task T-gkvi, #323): the type of the ruleset's source
// (Repository, Organization), the source and the ID of the ruleset, each the zero
// value when the answer lacks it. The read-back reads by branch and does not use
// them; the attribution of the probe's refusal does.
type Rule struct {
	Type              string
	Parameters        map[string]any
	RulesetSourceType string
	RulesetSource     string
	RulesetID         int64
}

// ActivityEntry is an activity of a ref: a push, a force push, the creation or
// the deletion of a branch, a merge. An activity whose actor the forge gives as
// null has the actor ID 0 and the login ghost, as a comment of a deleted account
// has the author 0 (forge.md, The calls of M2d).
type ActivityEntry struct {
	ActorID    int64
	ActorLogin string
	Ref        string
	Before     string
	After      string
	Timestamp  time.Time
	Type       string // the activity_type of the forge
}

// Comment is a comment of an issue.
type Comment struct {
	ID          int64
	AuthorID    int64 // 0 for a deleted account
	AuthorLogin string
	App         string // the slug of the App that made it, or empty
	Created     time.Time
	Updated     time.Time // the time of its last edit
	Body        string
}

// PullRequest is a pull request, as the calls of M2e give it (forge.md, The calls
// of M2e).
type PullRequest struct {
	Number   int
	NodeID   string
	State    string // open or closed
	Draft    bool
	Merged   bool
	HeadSHA  string
	MergeSHA string // the merge commit's SHA, when it is merged
	Body     string // the body; empty when the forge gives none (null)
}

// PullRequestQuery names a pull request by its number, or by the branch it
// comes from; exactly one of the two is set.
type PullRequestQuery struct {
	Number int
	Branch string
}

// Status is a commit status of a head: the creator is the account that posted it,
// which the caller compares with the App's bot (forge.md, The calls of M2e).
type Status struct {
	Context      string
	State        string
	Description  string
	CreatorID    int64
	CreatorLogin string
}

// CheckRun is a check run of a head, the newest of its name: a job of GitHub
// Actions is one of the App with the ID 15368.
type CheckRun struct {
	Name       string
	Status     string
	Conclusion string
	AppID      int64
	AppSlug    string
}

// HeadChecks is each status of a head, newest first as the forge gives them, and
// each latest check run.
type HeadChecks struct {
	Statuses  []Status
	CheckRuns []CheckRun
}

// The four states of a commit status.
var StatusStates = []string{"success", "failure", "pending", "error"}

// MaxDescription is the longest description of a commit status, in characters.
const MaxDescription = 140

// FullSHA tells a full object ID: 40 lowercase hexadecimal characters. The calls
// of M2e that take a SHA refuse any other text before any request.
func FullSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

// ValidStatus gives an error when state is not one of the four, or description is
// longer than MaxDescription characters (the check is made before any request).
func ValidStatus(state, description string) error {
	if !slices.Contains(StatusStates, state) {
		return fmt.Errorf("the state %q is none of %s", state, strings.Join(StatusStates, ", "))
	}
	if n := utf8.RuneCountInString(description); n > MaxDescription {
		return fmt.Errorf("the description has %d characters, over %d", n, MaxDescription)
	}
	return nil
}

// Error is a forge call that failed (forge.md, Forge errors): the call, the
// status (0 when no response came) and the first line of the message.
type Error struct {
	Call    string
	Status  int
	Message string
}

func (e *Error) Error() string {
	if e.Status == 0 {
		return fmt.Sprintf("%s: %s", e.Call, e.Message)
	}
	return fmt.Sprintf("%s: status %d: %s", e.Call, e.Status, e.Message)
}

// Missing gives each of the six that declared lacks, in the order of the table.
func Missing(declared []Capability) []Capability {
	have := map[Capability]bool{}
	for _, c := range declared {
		have[c] = true
	}
	var missing []Capability
	for _, c := range Capabilities {
		if !have[c] {
			missing = append(missing, c)
		}
	}
	return missing
}

// The permissions of M2a (forge.md, The six capabilities; run.md, step 1).
var Permissions = []struct{ Name, Level string }{
	{"contents", "write"}, {"issues", "write"}, {"metadata", "read"},
}

// CheckPermissions gives each permission of M2a that got lacks, as
// "name: level"; write holds read.
func CheckPermissions(got map[string]string) []string { return missingOf(got, Permissions) }

func missingOf(got map[string]string, list []struct{ Name, Level string }) []string {
	rank := map[string]int{"read": 1, "write": 2}
	var missing []string
	for _, p := range list {
		if rank[got[p.Name]] < rank[p.Level] {
			missing = append(missing, p.Name+": "+p.Level)
		}
	}
	return missing
}

// LoopPermissions are the permissions that the loop of M2e adds to those of M2a
// (forge.md, The calls of M2e): checks read amends ADR-0014 decision 3 (#336).
var LoopPermissions = []struct{ Name, Level string }{
	{"pull_requests", "write"}, {"statuses", "write"}, {"checks", "read"},
}

// CheckLoopPermissions gives each permission of the loop that got lacks, as
// "name: level"; write holds read.
func CheckLoopPermissions(got map[string]string) []string { return missingOf(got, LoopPermissions) }
