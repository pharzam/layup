package forge

import (
	"context"
	"fmt"
	"time"
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
// capability of issues and comments, comes with M2c, which specifies its call.
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
	HasCommit     bool // the list of branches is not empty
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
func CheckPermissions(got map[string]string) []string {
	rank := map[string]int{"read": 1, "write": 2}
	var missing []string
	for _, p := range Permissions {
		if rank[got[p.Name]] < rank[p.Level] {
			missing = append(missing, p.Name+": "+p.Level)
		}
	}
	return missing
}
