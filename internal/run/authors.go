package run

// This file holds the authors of a change and the reviewer's pair
// (docs/spec/loop.md, The independence of the reviewers; docs/architecture.md
// §9, Routing; task T-6zj3, #361, row 66e of the plan of M2e): the authors are
// computed from the records and the ancestry of the commits, never from a field
// that an agent writes. The status that posts `not-active` is #362's; the
// authors that an Operator's stall answer makes (a diff passed to an attempt)
// have no record before M3d, and are added there.

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
)

// A Review is the step whose authors are asked for.
type Review int

const (
	// PlanReview is the review of a task's plan: its authors are the harnesses of the plan sessions.
	PlanReview Review = iota
	// Verification is the review of a head: its authors are the harnesses of the
	// sessions bound to a commit of the diff from the base, and the harness that
	// wrote the frozen tests.
	Verification
)

// step is the step of the step table that the review is.
func (r Review) step() string {
	if r == PlanReview {
		return "plan-review"
	}
	return "verification"
}

// ErrNoHarnessOutsideAuthors is the refusal of a reviewer: no admitted harness is
// outside the authors. The text is the one that the status layup/verify posts
// after `not-active: ` (loop.md).
var ErrNoHarnessOutsideAuthors = errors.New("no harness outside the authors")

// authorsFrom computes the harnesses that are authors of a review of the task
// from its records: sessions.tsv (written by layup run), the events of the task,
// and the rows of its frozen list. inDiff says whether a commit is in the diff
// from the base to the head. The result is sorted, each harness once.
//
// The authors of a plan review are the harnesses of the sessions of the task
// whose role is the plan role of the step table; the record has a role and no
// step, so a session of that role that answered a question is counted too, on the
// safe side (known limit, loop.md).
func authorsFrom(r Review, task string, sessions, events, tests [][]string, inDiff func(sha string) (bool, error)) ([]string, error) {
	var out []string
	add := func(h string) {
		if !slices.Contains(out, h) {
			out = append(out, h)
		}
	}
	if r == PlanReview {
		for _, s := range sessions {
			if s[sessionTask] == task && s[sessionRole] == stepRole("plan-of-a-task") {
				add(s[sessionHarness])
			}
		}
		slices.Sort(out)
		return out, nil
	}
	for _, e := range events {
		if e[1] != "bound" {
			continue
		}
		i := slices.IndexFunc(sessions, func(s []string) bool { return s[sessionID] == e[3] })
		if i < 0 {
			return nil, fmt.Errorf("the event bound of %s has no start row in sessions.tsv", e[3])
		}
		if sessions[i][sessionTask] != task {
			continue
		}
		in, err := inDiff(e[5])
		if err != nil {
			return nil, err
		}
		if in {
			add(sessions[i][sessionHarness])
		}
	}
	for _, t := range tests {
		add(t[3])
	}
	slices.Sort(out)
	return out, nil
}

// Authors gives the harnesses that are authors of the review r of the task's
// change from base to head, from the records and the run's clone: a commit is in
// the diff when it is an ancestor of the head and not of the base.
func (s *Sessions) Authors(ctx context.Context, task string, r Review, base, head string) ([]string, error) {
	sessions, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return nil, err
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	tests, err := s.readTests(ctx, task)
	if err != nil {
		return nil, err
	}
	return authorsFrom(r, task, sessions, events, tests, func(sha string) (bool, error) {
		inHead, err := git.IsAncestor(s.Clone, sha, head)
		if err != nil || !inHead {
			return false, err
		}
		inBase, err := git.IsAncestor(s.Clone, sha, base)
		return !inBase, err
	})
}

// PairFor gives the first pair of records:routing.tsv for the role and the tier
// that is admitted and whose harness is not in except: the harness is in the
// register harnesses and its last probe, at the version of its last row of
// records:harnesses.tsv, passed, and the model's row of models has use yes
// (route.Admitted). It gives the harness register row and the row of models.tsv;
// with none, route.ErrNoPair.
func (s *Sessions) PairFor(ctx context.Context, role, tier string, harnesses, models [][]string, except []string) ([]string, []string, error) {
	routing, err := s.readTable(ctx, "routing.tsv", records.ReadRouting)
	if err != nil {
		return nil, nil, err
	}
	probes, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
	if err != nil {
		return nil, nil, err
	}
	last := map[string]string{} // the version of each harness's last row
	for _, r := range probes {
		last[r[1]] = r[2]
	}
	row := func(harness string) []string {
		for _, h := range harnesses {
			if h[regHarness] == harness {
				return h
			}
		}
		return nil
	}
	harness, model, err := route.Pair(routing, role, tier, func(h, m string) bool {
		v, probed := last[h]
		return !slices.Contains(except, h) && probed && row(h) != nil && route.Admitted(probes, models, h, v, m)
	})
	if err != nil {
		return nil, nil, err
	}
	for _, m := range models {
		if m[0] == harness && m[1] == model {
			return row(harness), m, nil
		}
	}
	return nil, nil, route.ErrNoPair // not reached: Admitted holds only for a row of models
}

// ReviewerPair gives the pair of a plan review or a verification of the task's
// change from base to head: the first admitted pair of the step's role and tier
// (the step table, never typed here) whose harness is not an author. With none
// the error holds ErrNoHarnessOutsideAuthors and route.ErrNoPair.
func (s *Sessions) ReviewerPair(ctx context.Context, task string, r Review, base, head string, harnesses, models [][]string) ([]string, []string, error) {
	res, err := ResolveStep(r.step(), "", "")
	if err != nil {
		return nil, nil, err
	}
	authors, err := s.Authors(ctx, task, r, base, head)
	if err != nil {
		return nil, nil, err
	}
	h, m, err := s.PairFor(ctx, res.Role, res.Tier, harnesses, models, authors)
	if errors.Is(err, route.ErrNoPair) {
		return nil, nil, fmt.Errorf("%w: %w", ErrNoHarnessOutsideAuthors, err)
	}
	return h, m, err
}
