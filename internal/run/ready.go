package run

// This file holds the ready and merge rule and the merge at the head SHA
// (docs/spec/loop.md, The ready and merge rule; task T-xvvd, #363, row 66g).
// The step loop that calls it is row 66k's.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/setup"
)

// A Readiness tells whether a head may be marked ready and merged, and the first
// reason when it may not.
type Readiness struct {
	Ready  bool
	Reason string
}

// jsonWhole reads a JSON number that is a whole number, as the forge's answer decodes it.
func jsonWhole(v any) (int64, bool) {
	switch n := v.(type) {
	case float64:
		return int64(n), n == math.Trunc(n) && math.Abs(n) < 1<<53
	case int:
		return int64(n), true
	case int64:
		return n, true
	case json.Number:
		i, err := n.Int64()
		return i, err == nil
	}
	return 0, false
}

// requiredJobs gives the names of the required checks of GitHub Actions that the
// rules hold, in the order of the rules. A rule of the type required_status_checks
// whose parameter is not a list of objects with a context that is a non-empty string
// and an integration_id that is a whole number of 1 or more is malformed, and is an
// error whatever its other entries are (docs/spec/loop.md, The ready and merge rule).
func requiredJobs(rules []forge.Rule) ([]string, error) {
	var jobs []string
	for _, r := range rules {
		if r.Type != "required_status_checks" {
			continue
		}
		list, ok := r.Parameters["required_status_checks"].([]any)
		if !ok {
			return nil, errors.New("a rule required_status_checks has no list required_status_checks")
		}
		for i, e := range list {
			m, ok := e.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("required_status_checks[%d] is not an object", i)
			}
			name, ok := m["context"].(string)
			if !ok || name == "" {
				return nil, fmt.Errorf("required_status_checks[%d] has no context that is a non-empty string", i)
			}
			id, ok := jsonWhole(m["integration_id"])
			if !ok || id < 1 {
				return nil, fmt.Errorf("required_status_checks[%d] (%s) has no integration_id that is a whole number of 1 or more", i, name)
			}
			if id == setup.ActionsApp && !slices.Contains(jobs, name) {
				jobs = append(jobs, name)
			}
		}
	}
	return jobs, nil
}

// Ready reads the forge and tells whether the head may be marked ready and merged
// (docs/spec/loop.md, The ready and merge rule): each of the four contexts is, by a
// read, a status of the App that is the newest of its context at the head, of any
// creator, and is success; and each required job of GitHub Actions has a check run
// of that App that is completed and success. It writes nothing. A read that fails or
// a rule that is malformed is an error, never a pass.
func (s *Sessions) Ready(ctx context.Context, head string) (Readiness, error) {
	if !forge.FullSHA(head) {
		return Readiness{}, fmt.Errorf("the head %q is not a full object ID", head)
	}
	if s.Bot == 0 {
		return Readiness{}, errors.New("the ID of the App's bot is not known: no status can be told to be the App's")
	}
	checks, err := s.Forge.Statuses(ctx, head)
	if err != nil {
		return Readiness{}, err
	}
	for _, c := range statusContexts {
		app, ok := appStatus(checks, c, s.Bot)
		if !ok {
			return Readiness{Reason: c + ": no status"}, nil
		}
		// The newest status of the context, of any creator, is the first one the forge gives.
		if i := slices.IndexFunc(checks.Statuses, func(st forge.Status) bool { return st.Context == c }); checks.Statuses[i].CreatorID != s.Bot {
			return Readiness{Reason: c + ": a status of another source is newer"}, nil
		}
		if app.State != "success" {
			return Readiness{Reason: c + ": " + app.State}, nil
		}
	}
	rules, err := s.Forge.EffectiveRules(ctx, s.Branch)
	if err != nil {
		return Readiness{}, err
	}
	jobs, err := requiredJobs(rules)
	if err != nil {
		return Readiness{}, fmt.Errorf("the rules of %s: %w", s.Branch, err)
	}
	if len(jobs) == 0 {
		return Readiness{Reason: "no required job of GitHub Actions"}, nil
	}
	for _, name := range jobs {
		found := false
		for _, run := range checks.CheckRuns {
			if run.Name != name || run.AppID != setup.ActionsApp {
				continue
			}
			found = true
			if run.Status != "completed" || run.Conclusion != "success" {
				return Readiness{Reason: "job " + name + ": " + strings.TrimSpace(run.Status+" "+run.Conclusion)}, nil
			}
		}
		if !found {
			return Readiness{Reason: "job " + name + ": no check run"}, nil
		}
	}
	return Readiness{Ready: true}, nil
}

// pullOf reads the pull request of the task by the number in the register.
func (s *Sessions) pullOf(ctx context.Context, task string) (forge.PullRequest, error) {
	row, err := s.taskRow(ctx, task)
	if err != nil {
		return forge.PullRequest{}, err
	}
	if row.PR == 0 {
		return forge.PullRequest{}, fmt.Errorf("%s has no pull request in %s", task, tasksPath)
	}
	pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Number: row.PR})
	if err != nil {
		return forge.PullRequest{}, err
	}
	if !found {
		return forge.PullRequest{}, fmt.Errorf("the pull request %d of %s is not on the forge", row.PR, task)
	}
	return pr, nil
}

// MarkReady marks the pull request of the task ready when its head is the head
// given and Ready holds. A pull request that is already ready needs no call.
func (s *Sessions) MarkReady(ctx context.Context, task, head string) (Readiness, error) {
	pr, err := s.pullOf(ctx, task)
	if err != nil {
		return Readiness{}, err
	}
	return s.markReady(ctx, pr, head)
}

func (s *Sessions) markReady(ctx context.Context, pr forge.PullRequest, head string) (Readiness, error) {
	switch {
	case pr.State != "open" || pr.Merged:
		return Readiness{Reason: "the pull request is not open"}, nil
	case head != pr.HeadSHA:
		return Readiness{Reason: "the head moved to " + pr.HeadSHA}, nil
	}
	r, err := s.Ready(ctx, head)
	if err != nil || !r.Ready || !pr.Draft {
		return r, err
	}
	if err := s.holdsLease(ctx); err != nil {
		return Readiness{}, err
	}
	if _, err := s.Forge.MarkReady(ctx, pr.NodeID); err != nil {
		return Readiness{}, err
	}
	return r, nil
}

// MergeTask merges the pull request of the task at the head SHA it read, and
// records the event merged; it gives the merge commit. A pull request that is
// merged only gets its event. It marks the draft ready, then reads the statuses
// and the jobs again just before the merge. A head that moved, before or by the
// forge's 409, gives the reason and no write; after a 409 the pull request and
// the statuses of the new head are read, and no second merge is sent from the old read.
func (s *Sessions) MergeTask(ctx context.Context, task, head string) (string, Readiness, error) {
	pr, err := s.pullOf(ctx, task)
	if err != nil {
		return "", Readiness{}, err
	}
	if pr.Merged {
		return pr.MergeSHA, Readiness{Ready: true}, s.RecordMerged(ctx, task, pr)
	}
	r, err := s.markReady(ctx, pr, head)
	if err != nil || !r.Ready {
		return "", r, err
	}
	if r, err = s.Ready(ctx, head); err != nil || !r.Ready {
		return "", r, err
	}
	if err := s.holdsLease(ctx); err != nil {
		return "", Readiness{}, err
	}
	merge, err := s.Forge.Merge(ctx, pr.Number, head)
	var refused *forge.Error
	if errors.As(err, &refused) && refused.Status == 409 {
		return s.afterRefusal(ctx, task, pr.Number, head, err)
	}
	if err != nil {
		return "", Readiness{}, err
	}
	pr.Merged, pr.State, pr.HeadSHA, pr.MergeSHA = true, "closed", head, merge
	return merge, Readiness{Ready: true}, s.RecordMerged(ctx, task, pr)
}

// afterRefusal is the forge's 409: the pull request is read again, and for a head that moved the statuses of the new head.
func (s *Sessions) afterRefusal(ctx context.Context, task string, number int, head string, refusal error) (string, Readiness, error) {
	pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Number: number})
	if err != nil {
		return "", Readiness{}, err
	}
	switch {
	case !found:
		return "", Readiness{}, fmt.Errorf("the pull request %d of %s is not on the forge after a refused merge", number, task)
	case pr.Merged:
		return pr.MergeSHA, Readiness{Ready: true}, s.RecordMerged(ctx, task, pr)
	case pr.HeadSHA == head:
		return "", Readiness{}, fmt.Errorf("the merge of the pull request %d at %s was refused and its head did not move: %w", number, head, refusal)
	}
	r, err := s.Ready(ctx, pr.HeadSHA)
	if err != nil {
		return "", Readiness{}, err
	}
	reason := "the head moved to " + pr.HeadSHA
	if !r.Ready {
		reason += ": " + r.Reason
	}
	return "", Readiness{Reason: reason}, nil
}

// RecordMerged writes the event merged of a merged pull request, once: the merge
// commit's SHA and the number of the pull request, in the attempt of the events
// bound of the head that was merged. A head that no session bound, or that
// sessions of more than one attempt bound, is an error and no event: the attempt
// is never guessed.
func (s *Sessions) RecordMerged(ctx context.Context, task string, pr forge.PullRequest) error {
	if !pr.Merged || !forge.FullSHA(pr.MergeSHA) || pr.Number < 1 {
		return fmt.Errorf("the pull request %d of %s has no merge commit to record", pr.Number, task)
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return err
	}
	attempts := map[int]bool{}
	for _, e := range events {
		switch {
		case e[1] == "merged" && e[5] == pr.MergeSHA:
			return nil
		case e[1] == "bound" && e[5] == pr.HeadSHA:
			a, _ := strconv.Atoi(e[2])
			attempts[a] = true
		}
	}
	if len(attempts) != 1 {
		return fmt.Errorf("the head %s of the pull request %d is the event bound of %d attempts of %s: the attempt of its event merged is not known", pr.HeadSHA, pr.Number, len(attempts), task)
	}
	var attempt int
	for a := range attempts {
		attempt = a
	}
	data, err := s.appendEvents(ctx, task, attempt, "", []string{fmt.Sprintf("merged %s %d", pr.MergeSHA, pr.Number)})
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{eventsPath(task): data}, fmt.Sprintf("layup run: %s merged as %s", task, pr.MergeSHA))
}
