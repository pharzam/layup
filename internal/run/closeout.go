package run

// This file holds the close-out commit and the carry-over of layup/verify
// (docs/spec/loop.md, The merge at the head SHA and the close-out commit; task
// T-xvvd, #363, row 66g). The step loop that calls it is row 66k's, and the push
// of the rendered commit row 66j's.

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/rules"
)

// A RenderedInput is the rendered commit above the verified head: the commit and
// the records that code rendered into it.
type RenderedInput struct {
	Head   string
	Render []rules.Record
}

// A CarryInput is what the caller knows of the verdict that a head carries: the
// task, the verifier's session, the registers that decide its independence, and
// the rendered commit, nil for none.
type CarryInput struct {
	Task, Verifier    string
	Harnesses, Models [][]string
	Rendered          *RenderedInput
	Merge             string // the merge commit of the default branch that BringUpToDate made, "" for none (task T-89xf, row 66x)
}

// closeOutPaths gives the two paths that the close-out commit may change: the task
// file and the completed log (decided here, task T-xvvd: no setup record holds them).
func closeOutPaths(task string) []string {
	return []string{"docs/tasks/" + task + ".md", "docs/tasks/completed.md"}
}

// verified gives the verified head of the verifier named, the head its review
// record passed at, and whether the verdict is success (nothing material in scope).
// The session must be a verifier of the task: VerifyAt refuses another, and the base
// of the task's change is the base of the event attempt of the verifier's attempt.
func (s *Sessions) verified(ctx context.Context, in CarryInput) (string, bool, error) {
	if in.Verifier == "" {
		return "", false, errors.New("no verifier session is named")
	}
	start, err := s.startRowOf(ctx, in.Verifier)
	if err != nil {
		return "", false, err
	}
	if start[sessionTask] != in.Task {
		return "", false, fmt.Errorf("the session %s is a session of %s, not of %s", in.Verifier, start[sessionTask], in.Task)
	}
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return "", false, err
	}
	base := ""
	for _, e := range events {
		if e[1] == "attempt" && e[2] == start[sessionAttempt] {
			base = e[4]
		}
	}
	if base == "" {
		return "", false, fmt.Errorf("%s has no event attempt %s: no base of the task's change", in.Task, start[sessionAttempt])
	}
	sha, st, err := s.VerifyAt(ctx, VerifyInput{Task: in.Task, Base: base, Head: start[sessionBase], Session: in.Verifier, Harnesses: in.Harnesses, Models: in.Models})
	if err != nil {
		return "", false, err
	}
	return sha, st.State == "success", nil
}

// regular tells whether the head holds the path as a regular file (mode 100644).
func (s *Sessions) regular(head, path string) (bool, error) {
	entries, err := git.LsTree(s.Clone, head, path)
	if err != nil {
		return false, err
	}
	return len(entries) == 1 && entries[0].Path == path && entries[0].Type == "blob" && entries[0].Mode == "100644", nil
}

// renderedOK tells whether the rendered commit is a permitted child of the verified head v: its only parent is v, the paths it
// changes are exactly the paths of the rendered records, each is a regular file, and the rendering is code's from the records with
// the pre-image of a PRD pinned to v (never a revision that the caller names).
func (s *Sessions) renderedOK(r *RenderedInput, v string) (bool, error) {
	if !forge.FullSHA(r.Head) || len(r.Render) == 0 {
		return false, nil
	}
	parents, err := git.Parents(s.Clone, r.Head)
	if err != nil || len(parents) != 1 || parents[0] != v {
		return false, err
	}
	var paths []string
	for _, rec := range r.Render {
		paths = append(paths, rec.Path)
	}
	names, err := git.DiffNames(s.Clone, v, r.Head)
	if err != nil {
		return false, err
	}
	slices.Sort(names)
	sorted := slices.Sorted(slices.Values(paths))
	if !slices.Equal(names, sorted) {
		return false, nil
	}
	for _, p := range paths {
		if ok, err := s.regular(r.Head, p); err != nil || !ok {
			return false, err
		}
	}
	return rules.RenderingEqual(rules.Git(s.Clone), r.Head, v, paths, r.Render)
}

// CheckRendered tells whether the rendered commit that the input names is a permitted child of the verified head. It checks no
// review: the verdict is the verified head's. It needs no session, as code made the commit.
func (s *Sessions) CheckRendered(ctx context.Context, in CarryInput) (bool, error) {
	if in.Rendered == nil {
		return false, errors.New("no rendered commit is named")
	}
	v, ok, err := s.verified(ctx, in)
	if err != nil || !ok {
		return false, err
	}
	return s.renderedOK(in.Rendered, v)
}

// scopeOK tells whether the head is one commit on the base that changes only the paths of the close-out, each a regular file.
func (s *Sessions) scopeOK(task, base, head string, names []string) (bool, error) {
	parents, err := git.Parents(s.Clone, head)
	if err != nil || len(parents) != 1 || parents[0] != base {
		return false, err
	}
	allowed := closeOutPaths(task)
	if len(names) == 0 {
		return false, nil
	}
	for _, n := range names {
		if !slices.Contains(allowed, n) {
			return false, nil
		}
		if ok, err := s.regular(head, n); err != nil || !ok {
			return false, err
		}
	}
	return true, nil
}

// carries tells whether the parent carries the verdict of the verified head v: it is v, the rendered commit that renderedOK accepts, or a
// merge commit that mergeEdge accepts whose first parent carries (an iterative walk down the first parents), the first parent being also
// the close-out head of such a head under a merge (task T-89xf, loop.md, A task whose base moved).
func (s *Sessions) carries(in CarryInput, v, parent string) (bool, error) {
	afterMerge := false
	for {
		if parent == v {
			return true, nil
		}
		if in.Rendered != nil && parent == in.Rendered.Head {
			return s.renderedOK(in.Rendered, v)
		}
		if !forge.FullSHA(parent) {
			return false, nil
		}
		parents, err := git.Parents(s.Clone, parent)
		if err != nil {
			return false, err
		}
		switch {
		case len(parents) == 2:
			if ok, err := s.mergeEdge(parents[0], parents[1], parent); err != nil || !ok {
				return false, err
			}
			afterMerge = true
		case len(parents) == 1 && afterMerge:
			names, err := git.DiffNames(s.Clone, parents[0], parent)
			if err != nil {
				return false, err
			}
			if ok, err := s.scopeOK(in.Task, parents[0], parent, names); err != nil || !ok {
				return false, err
			}
			afterMerge = false
		default:
			return false, nil
		}
		parent = parents[0]
	}
}

// PushCloseOut is the push hook of the close-out session: the checks of any task
// session, then three more, each refused as a developer's would be with the diff as
// a payload: scope when the head is not one commit on the session's base that changes
// only the two paths as regular files, and verdict when the base does not carry the
// verdict. A session of another task or of another role, a verifier of another task
// and a failed read are errors that write nothing.
func (s *Sessions) PushCloseOut(ctx context.Context, id string, in CarryInput) error {
	return s.pushChecked(ctx, id, func(start []string, base, head string, names []string) (string, error) {
		if start[sessionTask] != in.Task || start[sessionRole] != stepRole("close-out") {
			return "", fmt.Errorf("the session %s is a %s session of %s, not a close-out of %s", id, start[sessionRole], start[sessionTask], in.Task)
		}
		if ok, err := s.scopeOK(in.Task, base, head, names); err != nil || !ok {
			return "scope", err
		}
		v, ok, err := s.verified(ctx, in)
		if err != nil {
			return "", err
		}
		if ok, err = s.carriesFrom(in, v, ok, base); err != nil || !ok {
			return "verdict", err
		}
		return "", nil
	})
}

// carriesFrom is carries for a verdict that may not be success.
func (s *Sessions) carriesFrom(in CarryInput, v string, verdict bool, parent string) (bool, error) {
	if !verdict {
		return false, nil
	}
	return s.carries(in, v, parent)
}

// CarryVerify posts layup/verify success on the head that the verdict carries to, through
// PostStatus, and gives the head: the head of the close-out session named, its bound SHA, or,
// for the empty session, the rendered commit of the input. It checks the whole edge again, as
// the push did: the verdict, the session of the task and the role, the commit count, the
// paths and modes, and the carry of the parent. Anything less is an error and posts nothing.
func (s *Sessions) CarryVerify(ctx context.Context, in CarryInput, session string) (string, error) {
	v, ok, err := s.verified(ctx, in)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("the verifier %s did not find nothing material in scope: no verdict to carry", in.Verifier)
	}
	var head, what string
	if session == "" && in.Merge != "" {
		parents, err := git.Parents(s.Clone, in.Merge)
		if err != nil {
			return "", err
		}
		if ok, err := s.carries(in, v, in.Merge); err != nil || !ok || len(parents) != 2 {
			return "", errors.Join(err, fmt.Errorf("%s is not a clean merge of the base above the verified head %s", in.Merge, v))
		}
		head, what = in.Merge, "a clean merge of the base"
	} else if session == "" {
		if in.Rendered == nil {
			return "", errors.New("no session and no rendered commit are named")
		}
		if ok, err := s.renderedOK(in.Rendered, v); err != nil || !ok {
			return "", errors.Join(err, fmt.Errorf("%s is not a permitted rendered commit above the verified head %s", in.Rendered.Head, v))
		}
		head, what = in.Rendered.Head, "the rendered commit"
	} else {
		start, err := s.startRowOf(ctx, session)
		if err != nil {
			return "", err
		}
		if start[sessionTask] != in.Task || start[sessionRole] != stepRole("close-out") {
			return "", fmt.Errorf("the session %s is a %s session of %s, not a close-out of %s", session, start[sessionRole], start[sessionTask], in.Task)
		}
		events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
		if err != nil {
			return "", err
		}
		if head = boundHead(events, session); head == "" {
			return "", ErrNoBoundHead
		}
		names, err := git.DiffNames(s.Clone, start[sessionBase], head)
		if err != nil {
			return "", err
		}
		if ok, err := s.scopeOK(in.Task, start[sessionBase], head, names); err != nil || !ok {
			return "", errors.Join(err, fmt.Errorf("%s is not one commit on %s that changes only %v as regular files", head, start[sessionBase], closeOutPaths(in.Task)))
		}
		if ok, err := s.carries(in, v, start[sessionBase]); err != nil || !ok {
			return "", errors.Join(err, fmt.Errorf("%s does not carry the verdict of the verified head %s", start[sessionBase], v))
		}
		what = "the close-out"
	}
	if _, err := s.PostStatus(ctx, head, Status{VerifyContext, "success", fmt.Sprintf("verified head %s plus %s", v, what)}); err != nil {
		return "", err
	}
	return head, nil
}

// mergeEdge tells whether the commit m is a clean merge of the second parent p2, a commit of the default branch, into the first parent p1, by
// the rule of loop.md, A task whose base moved: with mb the merge base, T = DiffNames(mb, p1) the files the task touches and P = DiffNames(mb, p2),
// T and P are disjoint, m differs from p1 only in P and from p2 only in T. That proves the combined tree from the commit alone; whether p1 carries
// the verdict is the caller's.
func (s *Sessions) mergeEdge(p1, p2, m string) (bool, error) {
	def, err := git.RevParse(s.Clone, "refs/heads/"+s.Branch)
	if err != nil {
		return false, err
	}
	if ok, err := git.IsAncestor(s.Clone, p2, def); err != nil || !ok {
		return false, err
	}
	mb, err := git.MergeBase(s.Clone, p1, p2)
	if err != nil {
		return false, err
	}
	var t, p, a, b []string
	for _, r := range []struct {
		out        *[]string
		from, head string
	}{{&t, mb, p1}, {&p, mb, p2}, {&a, p1, m}, {&b, p2, m}} {
		if *r.out, err = git.DiffNames(s.Clone, r.from, r.head); err != nil {
			return false, err
		}
	}
	within := func(sub, set []string) bool {
		return !slices.ContainsFunc(sub, func(x string) bool { return !slices.Contains(set, x) })
	}
	return !slices.ContainsFunc(t, func(x string) bool { return slices.Contains(p, x) }) && within(a, p) && within(b, t), nil
}
