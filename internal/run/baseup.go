package run

// This file holds a task whose base moved (docs/spec/loop.md, A task whose base moved: the carry and the redo; task T-89xf, #424, row 66x):
// the up-to-date predicate, the clean merge of the default branch with the carry of layup/verify over it, and the choice of the redo.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
)

// TipOf gives the SHA of the last event bound of the attempt, of any session (the bound of code has none, which lastBound skips), "" for none.
func TipOf(events [][]string, attempt int) string {
	tip := ""
	for _, e := range events {
		if a, _ := strconv.Atoi(e[2]); e[1] == "bound" && a == attempt {
			tip = e[5]
		}
	}
	return tip
}

// latestAttempt gives the largest attempt of an event attempt and its base, 0 for none.
func latestAttempt(events [][]string) (int, string) {
	n, base := 0, ""
	for _, e := range events {
		if a, _ := strconv.Atoi(e[2]); e[1] == "attempt" && a >= n {
			n, base = a, e[4]
		}
	}
	return n, base
}

// UpToDate fetches the default branch and gives whether its head is an ancestor of head, with that head: the one predicate of a head that is
// current, which the close-out (row 66y) asks too. A head brought current by the merge of BringUpToDate is current.
func (s *Sessions) UpToDate(ctx context.Context, head string) (current bool, def string, err error) {
	if def, err = s.Store.FetchBranch(ctx, s.Branch); err != nil {
		return false, "", err
	}
	current, err = git.IsAncestor(s.Clone, def, head)
	return current, def, err
}

// commitsOfCode gives the commits of the events bound of the attempt that have no session, in order.
func commitsOfCode(events [][]string, attempt int) []string {
	var out []string
	for _, e := range events {
		if a, _ := strconv.Atoi(e[2]); e[1] == "bound" && e[3] == "" && a == attempt {
			out = append(out, e[5])
		}
	}
	return out
}

// BringUpToDate leaves a verified task whose base moved in one recoverable end state: its head current with the default branch and carrying
// layup/verify, or the task redone. It is read-only for a head that is current (and posts the carry again for a merge of its own, which is
// idempotent). The loop-permission check is the caller's (row 66k).
func (s *Sessions) BringUpToDate(ctx context.Context, in CarryInput) (Advance, error) {
	s.fns()
	task := in.Task
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return Advance{}, err
	}
	for _, e := range events {
		if e[1] == "closed" { // the redo was made: its one records commit holds this event
			return Advance{Task: task, Replaced: e[6]}, nil
		}
	}
	attempt, base := latestAttempt(events)
	if attempt == 0 {
		return Advance{}, fmt.Errorf("%s has no event attempt", task)
	}
	if intent, err := s.Store.ReadFile(ctx, intentPath(task)); err != nil {
		return Advance{}, err
	} else if intent != nil { // a redo decided: it goes on before any comparison
		return s.redo(ctx, in, events, attempt, "", intent)
	}
	x := TipOf(events, attempt)
	if x == "" {
		return Advance{}, fmt.Errorf("attempt %d of %s has no event bound: no head to bring up to date", attempt, task)
	}
	if adopted, err := s.adoptPushed(ctx, in, events, attempt, x); err != nil {
		return Advance{}, err
	} else if adopted != "" { // a merge that was pushed and never bound: it is the tip now, whatever the default branch has become
		x = adopted
		if events, err = s.readTable(ctx, eventsPath(task), records.ReadEvents); err != nil {
			return Advance{}, err
		}
	}
	current, d, err := s.UpToDate(ctx, x)
	if err != nil {
		return Advance{}, err
	}
	if ok, err := git.IsAncestor(s.Clone, base, d); err != nil {
		return Advance{}, err
	} else if !ok {
		return Advance{}, fmt.Errorf("the base %s of attempt %d of %s is not an ancestor of the default branch %s", base, attempt, task, d)
	}
	if current {
		if parents, err := git.Parents(s.Clone, x); err != nil {
			return Advance{}, err
		} else if len(parents) == 2 && slices.Contains(commitsOfCode(events, attempt), x) { // the merge of this call: the status may be undone
			in.Merge = x
			if _, err := s.CarryVerify(ctx, in, ""); err != nil {
				return Advance{}, err
			}
		}
		return Advance{Task: task, Verified: true, Head: x}, nil
	}
	v, ok, err := s.verified(ctx, in)
	if err != nil {
		return Advance{}, err
	}
	if !ok {
		return Advance{}, fmt.Errorf("the verifier %s did not find nothing material in scope: no verdict to bring up to date", in.Verifier)
	}
	branch := fmt.Sprintf("task/%s/%d", task, attempt)
	if pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: branch}); err != nil {
		return Advance{}, err
	} else if found && pr.Merged {
		return Advance{Task: task, Stop: fmt.Sprintf("%s: the pull request %d is merged: the merge is recorded, not redone", task, pr.Number)}, nil
	} else if found && pr.State != "open" {
		return Advance{Task: task, Wait: fmt.Sprintf("waiting for the Operator: %s: the pull request %d is closed and not merged", task, pr.Number)}, nil
	}
	mb, err := git.MergeBase(s.Clone, x, d)
	if err != nil {
		return Advance{}, err
	}
	touched, err := git.DiffNames(s.Clone, mb, x)
	if err != nil {
		return Advance{}, err
	}
	moved, err := git.DiffNames(s.Clone, mb, d)
	if err != nil {
		return Advance{}, err
	}
	if !slices.ContainsFunc(touched, func(p string) bool { return slices.Contains(moved, p) }) {
		if m, clean, err := s.mergeInto(ctx, task, attempt, x, d); err != nil {
			return Advance{}, err
		} else if clean {
			return s.carryMerge(ctx, in, events, attempt, v, x, m)
		}
	}
	for _, c := range commitsOfCode(events, attempt) {
		if parents, err := git.Parents(s.Clone, c); err != nil {
			return Advance{}, err
		} else if len(parents) == 1 { // the rendered commit of the first bet
			return Advance{Task: task, Stop: fmt.Sprintf("waiting for the Operator: %s: a redo would drop the rendered commit", task)}, nil
		}
	}
	return s.redo(ctx, in, events, attempt, d, nil)
}

// mergeInto makes the merge commit of d into x, by the App at the later committer time of the two, and keeps it by a ref of layup's own
// namespace; clean is false for a merge that conflicts.
func (s *Sessions) mergeInto(ctx context.Context, task string, attempt int, x, d string) (m string, clean bool, err error) {
	if s.App == "" || s.Bot == 0 {
		return "", false, fmt.Errorf("the App's slug and bot ID are unknown, so the merge cannot be made as the App")
	}
	var at time.Time
	for _, c := range []string{x, d} {
		t, err := git.CommitTime(s.Clone, c)
		if err != nil {
			return "", false, err
		}
		if t.After(at) {
			at = t
		}
	}
	bot := s.App + "[bot]"
	who := git.Identity{Name: bot, Email: fmt.Sprintf("%d+%s@users.noreply.github.com", s.Bot, bot), Time: at}
	branch := fmt.Sprintf("task/%s/%d", task, attempt)
	if m, clean, err = git.Merge(s.Clone, x, d, fmt.Sprintf("layup run: merge %s into %s", short(d), branch), who); err != nil || !clean {
		return "", false, err
	}
	return m, true, git.UpdateRef(s.Clone, fmt.Sprintf("refs/heads/layup/merge/%s-%d", task, attempt), m)
}

// carryMerge pushes the merge commit m as code and carries layup/verify over it: m is checked by the carry rule before anything is written;
// the event push announces the push, the lease is read right before it, the event bound follows, and the status comes last. Each act is
// skipped when its event or the branch shows it done.
func (s *Sessions) carryMerge(ctx context.Context, in CarryInput, events [][]string, attempt int, v, x, m string) (Advance, error) {
	task := in.Task
	in.Merge = m
	if ok, err := s.carries(in, v, m); err != nil {
		return Advance{}, err
	} else if !ok {
		return Advance{}, fmt.Errorf("the merge %s of %s does not carry the verdict of the verified head %s", m, task, v)
	}
	branch := fmt.Sprintf("task/%s/%d", task, attempt)
	has := func(kind string) bool {
		return slices.ContainsFunc(events, func(e []string) bool { return e[1] == kind && e[3] == "" && e[5] == m && e[6] == branch })
	}
	event := func(kind string) error { return s.codeEvent(ctx, task, attempt, kind, m, branch) }
	if !has("bound") {
		if !has("push") {
			if err := event("push"); err != nil {
				return Advance{}, err
			}
		}
		switch tip, found, err := s.Store.BranchHead(ctx, branch); {
		case err != nil:
			return Advance{}, err
		case !found || (tip != x && tip != m):
			return Advance{}, fmt.Errorf("%w: %s holds %q, not the head %s or its merge %s", errPushRefused, branch, tip, x, m)
		case tip == x:
			if err := s.holdsLease(ctx); err != nil {
				return Advance{}, err
			}
			if err := s.Store.PushHead(ctx, m, branch); err != nil {
				return Advance{}, err
			}
		}
		if err := event("bound"); err != nil {
			return Advance{}, err
		}
	}
	if _, err := s.CarryVerify(ctx, in, ""); err != nil {
		return Advance{}, err
	}
	return Advance{Task: task, Verified: true, Head: m}, nil
}

// codeEvent commits the event push or bound of the commit sha of code (session —) to the branch, in one fenced records commit.
func (s *Sessions) codeEvent(ctx context.Context, task string, attempt int, kind, sha, branch string) error {
	data, err := s.appendEvents(ctx, task, attempt, "", []string{kind + " " + sha + " " + branch})
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{eventsPath(task): data}, "layup run: the "+kind+" of the merge of the base to "+branch)
}

// adoptPushed resolves an outstanding push of a merge before another merge is computed: when an event push of code has no bound and the
// branch of the forge holds that commit, the push was made and its bound was lost. The commit must be a merge whose first parent is the
// tip x and that the carry rule accepts; its bound is then written, its status posted, and it is given as the new tip. A branch that holds
// x (the push was not made) or no announced commit leaves everything as it is; any other head of the branch is an error.
func (s *Sessions) adoptPushed(ctx context.Context, in CarryInput, events [][]string, attempt int, x string) (string, error) {
	var announced []string
	for _, e := range events {
		if a, _ := strconv.Atoi(e[2]); e[1] == "push" && e[3] == "" && a == attempt && !slices.ContainsFunc(events, func(b []string) bool { return b[1] == "bound" && b[3] == "" && b[5] == e[5] }) {
			announced = append(announced, e[5])
		}
	}
	if len(announced) == 0 {
		return "", nil
	}
	branch := fmt.Sprintf("task/%s/%d", in.Task, attempt)
	tip, found, err := s.Store.BranchHead(ctx, branch)
	if err != nil || !found || tip == x || !slices.Contains(announced, tip) {
		return "", err
	}
	v, ok, err := s.verified(ctx, in)
	if err != nil {
		return "", err
	}
	parents, err := git.Parents(s.Clone, tip)
	if err != nil {
		return "", err
	}
	in.Merge = tip
	if carried, err := s.carries(in, v, tip); err != nil || !ok || !carried || len(parents) != 2 || parents[0] != x {
		return "", errors.Join(err, fmt.Errorf("%s holds the announced %s, which is not a merge of the base over %s that carries the verdict", branch, tip, x))
	}
	if err := s.codeEvent(ctx, in.Task, attempt, "bound", tip, branch); err != nil {
		return "", err
	}
	if _, err := s.CarryVerify(ctx, in, ""); err != nil {
		return "", err
	}
	return tip, nil
}

// producedFrom tells whether the tip is reached from the head by merges of code alone: each commit from the tip down is the commit of an
// event bound with no session and a merge, and its first parent is the next, down to the head.
func (s *Sessions) producedFrom(events [][]string, attempt int, tip, head string) (bool, error) {
	code := commitsOfCode(events, attempt)
	for cur := tip; cur != head; {
		if !slices.Contains(code, cur) {
			return false, nil
		}
		parents, err := git.Parents(s.Clone, cur)
		if err != nil || len(parents) != 2 {
			return false, err
		}
		cur = parents[0]
	}
	return true, nil
}
