package run

// This file holds the pull request of a task (docs/spec/loop.md, The pull
// request; task T-xvvd, #363, row 66g of the plan of M2e). The step loop that
// calls it is row 66k's; the edit of the body to link the review record is row
// 66q's.

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// A PullInput is what the caller knows of the draft to open: the task, its
// attempt, and the type and the description of the title.
type PullInput struct {
	Task        string
	Attempt     int
	Type        string
	Description string
}

// The title is the form of the checked-in pr-title check (.github/workflows/pr-title.yml):
// a type of its list and a subject that begins with a lowercase letter, has three characters
// or more and does not end with a full stop. The task ID is not in it (decided here, task T-xvvd).
var (
	titleTypes   = []string{"feat", "fix", "docs", "refactor", "perf", "test", "build", "ci", "chore", "revert"}
	titleSubject = regexp.MustCompile(`^[a-z].+[^.]$`)
)

// checkTitle refuses a title that the checked-in pr-title check would refuse, and one with a
// control character.
func checkTitle(typ, description string) error {
	if !slices.Contains(titleTypes, typ) {
		return fmt.Errorf("the type %q of the title is none of %s", typ, strings.Join(titleTypes, ", "))
	}
	if !titleSubject.MatchString(description) || strings.IndexFunc(description, unicode.IsControl) >= 0 {
		return fmt.Errorf("the description %q of the title is not one line of three characters or more that begins with a lowercase letter, ends with no full stop and holds no control character", description)
	}
	return nil
}

// PullBody gives the body of the draft: the line that closes the issue, a blank
// line, the line that links the records commit, and each further line, with no
// final line feed. The edit that links the review record (row 66q) is one more line.
func PullBody(issue int, web, repo, records string, more ...string) string {
	lines := []string{"Closes #" + strconv.Itoa(issue), "", recordsLine(web, repo) + records}
	return strings.Join(append(lines, more...), "\n")
}

// recordsLine is the line of the body that links the records commit, up to the commit itself.
func recordsLine(web, repo string) string {
	return "The records commit this task started from: " + strings.TrimRight(web, "/") + "/" + repo + "/commit/"
}

// bodyRecords gives the records commit that a body links, the inverse of PullBody: the snapshot the draft was opened with, which
// stays in the body of the pull request through later binds, a restart and the edit of row 66q. ok is false for a body with no such line.
func bodyRecords(body, web, repo string) (commit string, ok bool) {
	for _, l := range strings.Split(body, "\n") {
		if c, found := strings.CutPrefix(l, recordsLine(web, repo)); found && c != "" && strings.Trim(c, "0123456789abcdef") == "" {
			return c, true
		}
	}
	return "", false
}

// withPR gives the register with the pr of the task set; the table given is not changed.
func withPR(tasks []Task, id string, n int) ([]Task, error) {
	i := slices.IndexFunc(tasks, func(t Task) bool { return t.ID == id })
	if i < 0 {
		return nil, fmt.Errorf("%s has no row of %s", tasksPath, id)
	}
	next := slices.Clone(tasks)
	next[i].PR = n
	return next, nil
}

// lastBound gives the SHA and the session of the last event bound of the attempt, "" for none.
func lastBound(events [][]string, attempt int) (sha, session string) {
	for _, e := range events {
		if a, _ := strconv.Atoi(e[2]); e[1] == "bound" && e[3] != "" && a == attempt {
			sha, session = e[5], e[3]
		}
	}
	return sha, session
}

// OpenPullRequest opens the draft pull request of the attempt's last bound
// head, once, and records its number in the register. Every input is from the
// records: the head is the SHA of the last event bound of the attempt, and the
// records commit of the body is the records column of the start row of the
// session of that event. The pull request of the branch is read before anything
// is opened: an open one is used, a merged one is returned untouched, a closed
// one that is not merged is an error (the redo is row 66k's). A head of the open or new
// pull request other than the bound SHA is an error too, and nothing is registered.
func (s *Sessions) OpenPullRequest(ctx context.Context, in PullInput) (forge.PullRequest, error) {
	if err := checkTitle(in.Type, in.Description); err != nil {
		return forge.PullRequest{}, err
	}
	if in.Attempt < 1 {
		return forge.PullRequest{}, fmt.Errorf("the attempt %d of %s is not an attempt", in.Attempt, in.Task)
	}
	task, err := s.taskRow(ctx, in.Task)
	if err != nil {
		return forge.PullRequest{}, err
	}
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return forge.PullRequest{}, err
	}
	head, session := lastBound(events, in.Attempt)
	if head == "" {
		return forge.PullRequest{}, fmt.Errorf("attempt %d of %s has no event bound: no head to open a pull request for", in.Attempt, in.Task)
	}
	start, err := s.startRowOf(ctx, session)
	if err != nil {
		return forge.PullRequest{}, err
	}
	if start[sessionTask] != in.Task {
		return forge.PullRequest{}, fmt.Errorf("the session %s of the event bound is a session of %s, not of %s", session, start[sessionTask], in.Task)
	}
	branch := "task/" + in.Task + "/" + strconv.Itoa(in.Attempt)
	pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: branch})
	if err != nil {
		return forge.PullRequest{}, err
	}
	switch {
	case found && pr.Merged:
		return pr, nil
	case found && pr.State != "open":
		return forge.PullRequest{}, fmt.Errorf("the pull request %d of %s is closed and not merged: the redo of a task is not this call's", pr.Number, branch)
	case !found:
		if err := s.holdsLease(ctx); err != nil {
			return forge.PullRequest{}, err
		}
		body := PullBody(task.Issue, s.Web, s.Target, start[sessionRecords])
		if pr, err = s.Forge.OpenDraft(ctx, branch, s.Branch, in.Type+": "+in.Description, body); err != nil {
			return forge.PullRequest{}, err
		}
	}
	if pr.HeadSHA != head {
		return forge.PullRequest{}, fmt.Errorf("the head of the pull request %d of %s moved: %s, not the bound %s", pr.Number, branch, pr.HeadSHA, head)
	}
	if task.PR == pr.Number {
		return pr, nil
	}
	return pr, s.registerPR(ctx, in.Task, pr.Number)
}

// registerPR replaces the row of the task in the register with its pr set, in one fenced records commit.
func (s *Sessions) registerPR(ctx context.Context, task string, n int) error {
	data, err := s.Store.ReadFile(ctx, tasksPath)
	if err != nil {
		return err
	}
	tasks, err := ReadRegister(data)
	if err != nil {
		return fmt.Errorf("%s: %w", tasksPath, err)
	}
	if tasks, err = withPR(tasks, task, n); err != nil {
		return err
	}
	var b bytes.Buffer
	if err := WriteRegister(&b, tasks); err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{tasksPath: b.Bytes()}, fmt.Sprintf("layup run: the pull request %d of %s", n, task))
}
