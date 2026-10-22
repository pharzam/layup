package run

// This file holds the tier of a task and the resolution of the words of the
// step table (docs/spec/loop.md, The tier and the context of a task; task
// T-6zj3, #361, row 66e of the plan of M2e). The plan's four fields are a typed
// value that the caller reads from the plan: the text form of a plan is #387's.

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
)

// A Plan is the four fixed fields of a task's plan (architecture.md §9):
// size, small or large, and the three lists.
type Plan struct {
	Size                                          string
	OpenQuestions, NewInterfaces, NewDependencies []string
}

// A Position is the progress position of a task (architecture.md §11).
type Position string

const (
	// Uphill is a task whose test list is not frozen, or that has an open unknown.
	Uphill Position = "uphill"
	// Downhill is any other task.
	Downhill Position = "downhill"
)

// PlanTier is the tier that a plan alone gives: the reasoning tier when its size
// is large or any of the three lists is not empty, the execution tier otherwise.
// A size that is neither small nor large counts as large (decided here: the
// field has no other value, and a typo must not choose the cheaper tier).
func PlanTier(p Plan) string {
	if p.Size != "small" || len(p.OpenQuestions) > 0 || len(p.NewInterfaces) > 0 || len(p.NewDependencies) > 0 {
		return "reasoning"
	}
	return "execution"
}

// TierOf is the tier of a task at an attempt: the progress position moves it,
// uphill to the reasoning tier and downhill to the execution tier
// (architecture.md §9); with no position the plan alone decides, as at its
// acceptance (loop.md).
func TierOf(p Plan, pos Position) string {
	switch pos {
	case Uphill:
		return "reasoning"
	case Downhill:
		return "execution"
	}
	return PlanTier(p)
}

// OpenQuestions counts the open questions of a task that its records give in
// M2e: the question rows of the last attempt that ended needs_context, unless an
// answer of an owner role's session that started after it is accepted by the
// next attempt (AnswerAccepted, which also asks the session of the owner role of
// roles and owners). The open material findings are counted by Position
// (openMaterial, row 66z) and the open escalations (M3) are added when they exist.
func OpenQuestions(events [][]string, results map[string][][]string, roles Roles, owners Owners, rep Replaced) int {
	last := -1
	for i, e := range events {
		if e[1] == "result" && e[6] == "done" && statusOf(results[e[3]]) == "needs_context" {
			last = i
		}
	}
	if last < 0 {
		return 0
	}
	for id, rows := range results {
		start := slices.IndexFunc(events, func(e []string) bool { return e[1] == "session" && e[3] == id })
		if start <= last {
			continue
		}
		for _, r := range rows {
			if r[0] != "answer" {
				continue
			}
			if _, ok := AnswerAccepted(events, results, roles, owners, id+"."+r[1], rep); ok {
				return 0
			}
		}
	}
	n := 0
	for _, e := range events {
		if e[1] == "result" && e[2] == events[last][2] && e[6] == "done" && statusOf(results[e[3]]) == "needs_context" {
			for _, r := range results[e[3]] {
				if r[0] == "question" {
					n++
				}
			}
		}
	}
	return n
}

// Position gives the progress position of a task from its records: downhill
// when its test list is frozen, it has no open question and no open material
// finding (openMaterial), else uphill.
func (s *Sessions) Position(ctx context.Context, task string) (Position, error) {
	frozen, err := s.readTests(ctx, task)
	if err != nil {
		return "", err
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return "", err
	}
	results := map[string][][]string{}
	for _, e := range events {
		if e[1] != "result" {
			continue
		}
		rows, err := s.readTable(ctx, "tasks/"+task+"/results/"+e[3]+".tsv", records.ReadResult)
		if err != nil {
			return "", err
		}
		if rows != nil {
			results[e[3]] = rows
		}
	}
	starts, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return "", err
	}
	row, err := s.taskRow(ctx, task)
	if err != nil {
		return "", err
	}
	view := loopView{Task: task, Kind: row.Kind, Events: events, Status: map[string]string{}}
	for _, r := range starts {
		if r[sessionTask] == task {
			view.Sessions = append(view.Sessions, r)
		}
	}
	for _, e := range events {
		if e[1] == "result" && e[6] == "done" {
			view.Status[e[3]] = statusOf(results[e[3]])
		}
	}
	if frozen != nil && OpenQuestions(events, results, sessionRoles(starts), Owners(DefaultOwners()), replacedOf(view)) == 0 {
		open, err := s.openMaterial(ctx, task, events, starts, results)
		if err != nil {
			return "", err
		}
		if !open {
			return Downhill, nil
		}
	}
	return Uphill, nil
}

// RoutingRole turns the role of an owner map ("Systems Architect") into the key
// of the routing register, the role ID of the step table ("systems-architect"):
// lower case, each run of spaces one hyphen.
func RoutingRole(owner string) string {
	return strings.ToLower(strings.Join(strings.Fields(owner), "-"))
}

// stepRole is the role of a step of the production step table, "" for none.
func stepRole(step string) string {
	for _, r := range session.StepTable() {
		if r.Step == step {
			return r.Role
		}
	}
	return ""
}

// Resolved is a step of the step table with its two words resolved: the role
// and the tier to route by.
type Resolved struct{ Step, Role, Tier string }

// ResolveStep resolves the words of the step table that are not routing keys
// (session.md): by-class becomes tier, the tier of the task (reasoning or
// execution); owner-of-its-kind becomes the routing role of ownerRole, the role
// of the owner map. A human has no routing role, and a step that is not in the
// table is an error.
func ResolveStep(step, tier, ownerRole string) (Resolved, error) {
	for _, row := range session.StepTable() {
		if row.Step != step {
			continue
		}
		out := Resolved{Step: step, Role: row.Role, Tier: row.Tier}
		if row.Tier == "by-class" {
			if tier != "reasoning" && tier != "execution" {
				return Resolved{}, fmt.Errorf("the tier of %s is the task's class, reasoning or execution; %q is neither", step, tier)
			}
			out.Tier = tier
		}
		if row.Role == "owner-of-its-kind" {
			if out.Role = RoutingRole(ownerRole); out.Role == "" || out.Role == "operator" {
				return Resolved{}, fmt.Errorf("the owner role %q of %s is no role of the routing register", ownerRole, step)
			}
		}
		return out, nil
	}
	return Resolved{}, fmt.Errorf("the step table has no step %q", step)
}

// openMaterial says whether a material finding is open (loop.md, the open unknowns): the typed record of the last verifier is read after it is
// selected and qualified. The last verifier is the last session by the order of the events of its start, of the role qa-engineer, that is a
// verifier by its position and head binding (verifierBound, the classification of loopView.qaKind): it started after an event bound of a developer's
// session of its own attempt, at the SHA that is the base of its start row. A plan review has no such bound, so it is never selected, and a session
// of the role that is bound to no head is no verifier. With none, nothing is open. The output of the selected verifier is qualified only after
// that: it must have ended done with the status completed, hold a comment row of result/review-record.md and have a valid handoff that names a
// comment; its kept copy must exist and be the file the result names by its hash (keptComment), parse, and its head must be the base of the start row.
// Otherwise the finding is open: a verifier that has not ended or asked, whose record is missing, whose handoff is not recorded yet or not valid,
// never clears a finding of an older record, and no older record stands in its place. Otherwise OpenMaterial of the record decides. A failed read is
// an error.
func (s *Sessions) openMaterial(ctx context.Context, task string, events, starts [][]string, results map[string][][]string) (bool, error) {
	roles := sessionRoles(starts)
	last := ""
	for _, e := range events {
		if e[1] != "session" || roles[e[3]] != stepRole("verification") {
			continue
		}
		if i := slices.IndexFunc(starts, func(r []string) bool { return r[sessionID] == e[3] }); i >= 0 && verifierBound(events, roles, starts[i]) {
			last = e[3]
		}
	}
	if last == "" {
		return false, nil
	}
	if _, done := resultEnd(events, last); !done || statusOf(results[last]) != "completed" || !slices.ContainsFunc(results[last], func(r []string) bool { return r[0] == "comment" && r[2] == "result/"+ReviewRecordFile }) {
		return true, nil
	}
	detail, found := "", false
	for _, e := range events {
		if e[1] == "handoff" && e[3] == last {
			detail, found = e[6], true
		}
	}
	if valid, comment := validHandoff(detail); !found || !valid || !comment {
		return true, nil
	}
	data, ok, err := s.keptComment(ctx, task, last, ReviewRecordFile, results[last])
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}
	rec, parsed := ParseReviewRecord(data)
	if i := slices.IndexFunc(starts, func(r []string) bool { return r[sessionID] == last }); i < 0 || rec.Head != starts[i][sessionBase] {
		return true, nil
	}
	return OpenMaterial(rec, parsed), nil
}

// verifierBound says whether a session of the role of the verification is a verifier by its position and head binding: an event bound of a session of
// a developer's role, of its attempt and at the base of its start row, comes before the event of its start. It is the classification of loopView.qaKind
// told from the events and the roles alone, for the task of any kind.
func verifierBound(events [][]string, roles Roles, start []string) bool {
	developers := []string{stepRole("implement"), stepRole("specification"), stepRole("shape-architecture")}
	for _, e := range events {
		if e[1] == "session" && e[3] == start[sessionID] {
			return false
		}
		if e[1] == "bound" && e[2] == start[sessionAttempt] && e[5] == start[sessionBase] && slices.Contains(developers, roles[e[3]]) {
			return true
		}
	}
	return false
}
