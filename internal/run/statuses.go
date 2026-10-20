package run

// This file holds the commit statuses of a head and the rule of their posting
// (docs/spec/loop.md, REQ-007 — The four layup/ statuses; task T-q0kr, #362, row
// 66f of the plan of M2e). The step loop that calls them is row 66k's.

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/rules"
)

// The four contexts of the statuses of a head.
const (
	GatesContext  = "layup/gates"
	SpecContext   = "layup/spec"
	RulesContext  = "layup/rules"
	VerifyContext = "layup/verify"
)

var statusContexts = []string{GatesContext, SpecContext, RulesContext, VerifyContext}

// A Status is a commit status to post: its context, one of the four states and
// a description of one line.
type Status struct{ Context, State, Description string }

// fit gives the description as one line of at most forge.MaxDescription
// characters: a CR or LF becomes a space, and a longer text is cut to one
// character less and ends in an ellipsis, so SetStatus never refuses it.
func fit(s string) string {
	s = strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
	if r := []rune(s); len(r) > forge.MaxDescription {
		return string(r[:forge.MaxDescription-1]) + "…"
	}
	return s
}

// appStatus gives the newest status of the context whose creator is the bot
// (HeadChecks holds the statuses newest first). A status of another creator is
// never the App's.
func appStatus(c forge.HeadChecks, context string, bot int64) (forge.Status, bool) {
	for _, st := range c.Statuses {
		if st.Context == context && st.CreatorID == bot {
			return st, true
		}
	}
	return forge.Status{}, false
}

// PostStatus posts st on the head sha unless the App's newest status of the
// context already has its state and description (loop.md, REQ-007). The SHA, the
// context, the state and the bot are checked before any request; the lease is
// read just before the post, and a lost lease posts nothing. It tells whether it
// posted.
func (s *Sessions) PostStatus(ctx context.Context, sha string, st Status) (bool, error) {
	if !forge.FullSHA(sha) {
		return false, fmt.Errorf("the status of %q: the SHA is not a full object ID", sha)
	}
	if !slices.Contains(statusContexts, st.Context) {
		return false, fmt.Errorf("the context %q is none of the four of layup run", st.Context)
	}
	st.Description = fit(st.Description)
	if err := forge.ValidStatus(st.State, st.Description); err != nil {
		return false, err
	}
	if s.Bot == 0 {
		return false, fmt.Errorf("the ID of the App's bot is not known: no status can be told to be the App's")
	}
	checks, err := s.Forge.Statuses(ctx, sha)
	if err != nil {
		return false, err
	}
	if cur, ok := appStatus(checks, st.Context, s.Bot); ok && cur.State == st.State && cur.Description == st.Description {
		return false, nil
	}
	if err := s.holdsLease(ctx); err != nil {
		return false, err
	}
	if err := s.Forge.SetStatus(ctx, sha, st.Context, st.State, st.Description); err != nil {
		return false, err
	}
	return true, nil
}

// RulesStatus gives the status of layup/rules for the result of the basic rule:
// its state, and the reason, with the path when the result has one. A result with
// no state is an error, never a pass.
func RulesStatus(r rules.Result) Status {
	switch r.State {
	case rules.Success, rules.Failure, rules.Error:
	default:
		return Status{RulesContext, "error", "error: no result"}
	}
	d := r.Reason
	if r.Path != "" {
		d += ": " + r.Path
	}
	return Status{RulesContext, string(r.State), fit(d)}
}

// A RulesInput is what the caller knows of the head whose layup/rules is asked:
// the task, the base, the head, the rendered records and the revision a PRD is
// rendered from, as rules.Input takes them, and RegisterRev, the commit of the
// register for a head that code produced.
type RulesInput struct {
	Task, Base, Head, PreRev, RegisterRev string
	Records                               []rules.Record
}

// registerRevOf gives the commit whose rule-path register judges the head: the
// records column of the start row of the session of the task that has the event
// bound of the head (sessions of one head must agree), or the commit that the
// caller gives for a head that no session pushed. Never the head's own register,
// and never a commit that the caller gives for a head that has a bound session.
func registerRevOf(task, head string, events, sessions [][]string, given string) (string, error) {
	found := ""
	for _, e := range events {
		if e[1] != "bound" || e[3] == "" || e[5] != head {
			continue
		}
		i := slices.IndexFunc(sessions, func(r []string) bool { return r[sessionID] == e[3] })
		if i < 0 {
			return "", fmt.Errorf("the event bound of %s has no start row in sessions.tsv", e[3])
		}
		if sessions[i][sessionTask] != task {
			continue
		}
		rec := sessions[i][sessionRecords]
		if found != "" && found != rec {
			return "", fmt.Errorf("the sessions bound to %s name two records commits", head)
		}
		found = rec
	}
	switch {
	case found != "" && given != "":
		return "", fmt.Errorf("%s has a bound session: the register commit is the one of its start row", head)
	case found != "":
		return found, nil
	case !forge.FullSHA(given):
		return "", fmt.Errorf("%s has no bound session of %s, and the register commit %q is no full SHA", head, task, given)
	}
	return given, nil
}

// RulesAt gives the status layup/rules of a head: RulesStatus of RulesResult.
func (s *Sessions) RulesAt(ctx context.Context, in RulesInput) (Status, error) {
	r, err := s.RulesResult(ctx, in)
	if err != nil {
		return Status{}, err
	}
	return RulesStatus(r), nil
}

// RulesResult gives the result of rules.Basic over the base and the head, by the register of
// the commit that registerRevOf names.
func (s *Sessions) RulesResult(ctx context.Context, in RulesInput) (rules.Result, error) {
	if !forge.FullSHA(in.Base) || !forge.FullSHA(in.Head) {
		return rules.Result{}, fmt.Errorf("the base and the head of layup/rules are full SHAs")
	}
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return rules.Result{}, err
	}
	sessions, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return rules.Result{}, err
	}
	rev, err := registerRevOf(in.Task, in.Head, events, sessions, in.RegisterRev)
	if err != nil {
		return rules.Result{}, err
	}
	return rules.Basic(rules.Git(s.Clone), rules.Input{Base: in.Base, Head: in.Head, RegisterRev: rev, PreRev: in.PreRev, Records: in.Records}), nil
}

// A HeadInput is what the caller knows of a head whose four statuses are asked:
// the inputs of GatesAt, SpecAt, RulesAt and VerifyAt.
type HeadInput struct {
	Task, Base, Head    string
	Copy                []byte
	Step                func(i, n int, kind string) func()
	Out                 io.Writer
	PreRev, RegisterRev string
	Records             []rules.Record
	Session             string
	Waiting             bool
	Harnesses, Models   [][]string
}

// A StatusPost is one status of PostStatuses: the SHA it is posted on, the status,
// and whether it was posted (false when the App's status was already at it).
type StatusPost struct {
	SHA    string
	Status Status
	Posted bool
}

// PostStatuses computes the four statuses of a head and posts each through
// PostStatus, in the order gates, spec, rules, verify, one after the other: a
// failure stops the call and keeps the earlier posts, and the next call posts the
// rest (loop.md, REQ-007). It gives each status with the SHA it is posted on. A
// scratch tree of the gates that cannot be removed does not stop the call: the
// error is returned after the four.
func (s *Sessions) PostStatuses(ctx context.Context, in HeadInput) ([]StatusPost, error) {
	var out []StatusPost
	var late error
	post := func(sha string, st Status) error {
		posted, err := s.PostStatus(ctx, sha, st)
		if err == nil {
			out = append(out, StatusPost{sha, st, posted})
		}
		return err
	}
	st, err := s.GatesAt(ctx, GatesInput{Task: in.Task, Base: in.Base, Head: in.Head, Step: in.Step, Out: in.Out})
	var cleanup *gate.CleanupError
	if errors.As(err, &cleanup) {
		late, err = err, nil
	}
	if err == nil {
		err = post(in.Head, st)
	}
	if err == nil {
		if st, err = s.SpecAt(ctx, SpecInput{Task: in.Task, Base: in.Base, Head: in.Head, Copy: in.Copy}); err == nil {
			err = post(in.Head, st)
		}
	}
	if err == nil {
		if st, err = s.RulesAt(ctx, RulesInput{Task: in.Task, Base: in.Base, Head: in.Head, PreRev: in.PreRev, RegisterRev: in.RegisterRev, Records: in.Records}); err == nil {
			err = post(in.Head, st)
		}
	}
	if err == nil {
		var sha string
		if sha, st, err = s.VerifyAt(ctx, VerifyInput{Task: in.Task, Base: in.Base, Head: in.Head, Session: in.Session, Waiting: in.Waiting, Harnesses: in.Harnesses, Models: in.Models}); err == nil {
			err = post(sha, st)
		}
	}
	return out, errors.Join(err, late)
}
