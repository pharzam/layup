package run

// This file holds the link of the review record in the body of the pull request
// (docs/spec/loop.md, The link of the review record in the body; task T-08rj,
// #413, row 66q of the plan of M2e). The step that calls it is row 66k's.

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// reviewLine is the last line of the body: the link of the comment of the review record.
func reviewLine(web, repo string, issue int, comment int64) string {
	return fmt.Sprintf("The review record: %s/%s/issues/%d#issuecomment-%d", strings.TrimRight(web, "/"), repo, issue, comment)
}

// LinkReviewRecord replaces the body of the task's pull request with the body that
// links the review record of the verifier's session; it gives whether it edited.
// Every input is from the records (loop.md, The link of the review record in the
// body): the comment is the one of the verifier's latest handoff, which must be
// valid, and the records commit is the one that the draft's body links, which a
// session that bound the reviewed head must have started at. The pull request of the attempt's branch must be the register's
// and at the reviewed head, and open; a body that is already the composed one is
// not edited; the lease is read just before the edit; the answer is read.
func (s *Sessions) LinkReviewRecord(ctx context.Context, task, verifier string) (bool, error) {
	start, err := s.startRowOf(ctx, verifier)
	if err != nil {
		return false, err
	}
	if start[sessionTask] != task {
		return false, fmt.Errorf("the session %s is of %s, not of %s", verifier, start[sessionTask], task)
	}
	if role := stepRole("verification"); start[sessionRole] != role {
		return false, fmt.Errorf("the session %s is of the role %s, not of the verification (%s)", verifier, start[sessionRole], role)
	}
	attempt, _ := strconv.Atoi(start[sessionAttempt])
	base := start[sessionBase]
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return false, err
	}
	newest, detail, bound := 0, "", []string{}
	for _, e := range events {
		a, _ := strconv.Atoi(e[2])
		switch {
		case e[1] == "attempt" && a > newest:
			newest = a
		case e[1] == "handoff" && e[3] == verifier:
			detail = e[6]
		case e[1] == "bound" && a == attempt && e[5] == base:
			bound = append(bound, e[3])
		}
	}
	if attempt != newest {
		return false, fmt.Errorf("the session %s is of attempt %d, not of the newest attempt %d of %s", verifier, attempt, newest, task)
	}
	if detail == "" {
		return false, fmt.Errorf("the session %s has no handoff", verifier)
	}
	var id, seen int
	if _, err := fmt.Sscanf(detail, "valid comment %d-%d", &id, &seen); err != nil || fmt.Sprintf("valid comment %d-%d", id, seen) != detail {
		return false, fmt.Errorf("the latest handoff of the session %s is %q, not \"valid comment <ID>-<seen>\"", verifier, detail)
	}
	t, err := s.taskRow(ctx, task)
	if err != nil {
		return false, err
	}
	if t.PR == 0 {
		return false, fmt.Errorf("%s has no pull request in the register", task)
	}
	copies, err := s.readTable(ctx, "copies.tsv", records.ReadCopies)
	if err != nil {
		return false, err
	}
	var copied []string
	for _, c := range copies {
		if c[copyComment] == strconv.Itoa(id) && c[copySeen] == strconv.Itoa(seen) && c[copyIssue] == strconv.Itoa(t.Issue) {
			copied = c
		}
	}
	if copied == nil {
		return false, fmt.Errorf("copies.tsv has no copy of the comment %d-%d of the issue %d", id, seen, t.Issue)
	}
	data, err := s.Store.ReadFile(ctx, copied[copyBody])
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, fmt.Errorf("the records have no %s", copied[copyBody])
	}
	rec, ok := ParseReviewRecord(data)
	if !ok {
		return false, fmt.Errorf("the copy %s is not a review record", copied[copyBody])
	}
	if rec.Head != base {
		return false, fmt.Errorf("the review record in %s is of the head %s, not of the verifier's base %s", copied[copyBody], rec.Head, base)
	}
	if len(bound) == 0 {
		return false, fmt.Errorf("attempt %d of %s has no event bound at the head %s", attempt, task, base)
	}

	branch := "task/" + task + "/" + start[sessionAttempt]
	pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: branch})
	switch {
	case err != nil:
		return false, err
	case !found:
		return false, fmt.Errorf("the branch %s has no pull request", branch)
	case pr.Number != t.PR:
		return false, fmt.Errorf("the pull request %d of the branch %s is not the register's %d", pr.Number, branch, t.PR)
	case pr.HeadSHA != base && !slices.ContainsFunc(events, func(e []string) bool { return e[1] == "bound" && e[3] == "" && e[5] == pr.HeadSHA }): // the rendered commit of the bet is the head of the pull request
		return false, fmt.Errorf("the head of the pull request %d is %s, not the reviewed %s", pr.Number, pr.HeadSHA, base)
	case pr.State != "open":
		return false, fmt.Errorf("the pull request %d is %s: only an open one is edited", pr.Number, pr.State)
	}
	// The snapshot of the opening is the one the draft's body links, if a session that bound the reviewed head started at it.
	commit, ok := bodyRecords(pr.Body, s.Web, s.Target)
	opened := false
	for _, id := range bound {
		row, err := s.startRowOf(ctx, id)
		if err != nil {
			return false, err
		}
		opened = opened || (ok && row[sessionRecords] == commit)
	}
	if !opened {
		return false, fmt.Errorf("the body of the pull request %d links no records commit of a session that bound the reviewed head %s", pr.Number, base)
	}
	want := PullBody(t.Issue, s.Web, s.Target, commit, "", reviewLine(s.Web, s.Target, t.Issue, int64(id)))
	if pr.Body == want {
		return false, nil
	}
	if err := s.holdsLease(ctx); err != nil {
		return false, err
	}
	got, err := s.Forge.EditBody(ctx, pr.Number, want)
	if err != nil {
		return false, err
	}
	if got.Body != want {
		return false, fmt.Errorf("the forge gave another body for the pull request %d than the one sent", pr.Number)
	}
	return true, nil
}
