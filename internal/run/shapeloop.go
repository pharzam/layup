package run

// This file holds the driver of a Shape task (docs/spec/loop.md, The driver of a Shape task; task
// T-jy6f, #415, row 66r of the plan of M2e): the position of a task from its records, and the acts that
// take it, one at a time, to a verified head.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// ShapeTaskInput is what the driver is given for one task: the task, the base of the run (the head of the
// default branch, where the first attempt starts), the run's records commit as it is now, the harness and
// model registers, and, for an architecture task, the results of the filter of each decision.
type ShapeTaskInput struct {
	Task, Base        string
	Records           func() string
	Harnesses, Models [][]string
	Options           map[string]spec.Result
	Copy              []byte // the copy of the answer of the confirmed bet, which the full check of every head after it reads (row 66y); nil before the first bet
}

// Advance is where the driver stopped: a verified head, a wait for the Operator, or a stop. A stop and a
// wait repeat on the next call, which reads the same records. FinishTask (row 66y) gives Merged, with Head the merged head, or
// Attempt, the attempt that its repair started.
type Advance struct {
	Task       string
	Verified   bool
	Merged     bool
	Head       string // the bound head that the verifier found nothing material in, or the merged head
	Attempt    int    // the attempt that a repair started, 0 for none
	Wait, Stop string
	Replaced   string // the ID of the task that redoes this one (task T-89xf, row 66x)
}

// shapeCycleCap is the review-cycle cap of a Shape task: the target's rule is 1, and 2 for a change of a
// gate, which a Shape task cannot make (a gate file is a rule path, refused before any push).
const shapeCycleCap = 1

// maxActs bounds the acts of one call: each act moves the records on, so a call that does not end is a defect.
const maxActs = 100

type position int

const (
	posNone           position = iota
	posAttempt                 // no attempt yet
	posPlan                    // run the plan session
	posPlanFinish              // post the plan and record its handoff
	posReview                  // run the plan review
	posReviewFinish            // post the plan review and record its handoff
	posDeveloper               // run the developer's session
	posToVerifier              // record the developer's handoff
	posShip                    // the pull request, the statuses, the acceptance, then the verifier
	posVerifierFinish          // post the review record and record its handoff
	posVerified                // the verdict is nothing material: link it, post the statuses, done
	posNextAttempt             // the verdict is material and the cap allows one more attempt
	posWait                    // wait for the Operator
	posStop                    // a case the driver does not go on from
	posQuestion                // a session of the attempt asked a question: the stages of questionloop.go
	posRetry                   // the last session of the attempt ended crash or no-result: one new attempt at its base
	posReplan                  // the plan review rejected the plan: one new attempt at its base, the plan stage reopened
)

// loopView is what the position is told from: the events and the rows of sessions.tsv of the task, the
// status of the kept result of each session that ended done, and the review record of the kept copy of
// each valid handoff that has one.
type loopView struct {
	Task, Kind string
	Events     [][]string
	Sessions   [][]string
	Status     map[string]string
	Records    map[string]ReviewRecord
	// The question in the loop (row 66z) is told from the typed rows of each result that ended done, the rows of copies.tsv
	// and the issue of the task.
	Results map[string][][]string
	Copies  [][]string
	Issue   int
	// Rejects holds the reason lines of each plan review whose kept copy is the file its result names by hash and
	// rejects with at least one reason.
	Rejects map[string][]string
	// Unread holds the plan reviews that ended completed whose verdict cannot be established from the kept copy (absent, other bytes than the result
	// names by hash, or no plan review form); Verdicts holds the verdict word of each plan review whose kept copy is read.
	Unread   map[string]bool
	Verdicts map[string]string
}

// loopStep is a position with what its act needs: the session it concerns, the message of a stop or a wait,
// the latest attempt with its base, and the developer's bound head.
type loopStep struct {
	Pos          position
	Session, Msg string
	Attempt      int
	Base, Head   string
	Ask          qStep // the stage of the question, for posQuestion
}

const nothingMaterial = "nothing material in scope"

// devStep and devFrom are the step of the developer of a Shape task and the name of its side of the table of transitions.
func devStep(kind string) string {
	if kind == "architecture" {
		return "shape-architecture"
	}
	return "specification"
}

func devFrom(kind string) string {
	if kind == "architecture" {
		return "developer (architecture)"
	}
	return "developer (specification)"
}

func (v loopView) role(id string) string {
	for _, s := range v.Sessions {
		if s[sessionID] == id {
			return s[sessionRole]
		}
	}
	return ""
}

// startIndex is the place of the event session of a session in the events, -1 for none.
func (v loopView) startIndex(id string) int {
	return slices.IndexFunc(v.Events, func(e []string) bool { return e[1] == "session" && e[3] == id })
}

// lastHandoff gives the detail of the latest event handoff of a session.
func (v loopView) lastHandoff(id string) (string, bool) {
	detail, found := "", false
	for _, e := range v.Events {
		if e[1] == "handoff" && e[3] == id {
			detail, found = e[6], true
		}
	}
	return detail, found
}

// qaKind tells the two sessions of the role qa-engineer apart by the order of the events, never the clock: the
// verifier starts after a bound of the developer's session of its own attempt, at that SHA; the plan review starts
// after the first valid handoff of the plan and before the first push. The first, never the last: a plan handed
// off after the review (the plan of a replanned attempt) must not turn an earlier review into neither (task T-t94t).
func (v loopView) qaKind(s []string) (review, verifier bool) {
	start := v.startIndex(s[sessionID])
	if start < 0 {
		return false, false
	}
	plan, push := -1, len(v.Events)
	for i, e := range v.Events {
		if ok, comment := validHandoff(e[6]); e[1] == "handoff" && ok && comment && v.role(e[3]) == stepRole("plan-of-a-task") {
			if plan < 0 {
				plan = i // a handoff of a comment: the question and the answer handoffs of a session of this role are none
			}
		}
		if e[1] == "push" && i < push {
			push = i
		}
	}
	for _, e := range v.Events[:start] {
		if e[1] == "bound" && e[2] == s[sessionAttempt] && e[5] == s[sessionBase] && v.role(e[3]) == stepRole(devStep(v.Kind)) {
			return false, true
		}
	}
	return plan >= 0 && start > plan && start < push, false
}

// outcome says how a session ended: running (no result yet, a stopped run), ended done with the status
// completed, or a stop that names the session.
type outcome int

const (
	running outcome = iota
	finished
	failed
)

func (v loopView) outcome(id, step string) (outcome, string) {
	result := ""
	for _, e := range v.Events {
		switch {
		case e[3] != id:
		case e[1] == "refused":
			return failed, fmt.Sprintf("%s: the session %s was refused (%s)", v.Task, id, e[6])
		case e[1] == "result":
			result = e[6]
		}
	}
	switch status := v.Status[id]; {
	case result == "":
		return running, ""
	case result != "done":
		return failed, fmt.Sprintf("%s: the session %s of %s ended %s", v.Task, id, step, result)
	case status != "completed":
		if status == "" {
			status = "with no kept result"
		}
		return failed, fmt.Sprintf("%s: the session %s of %s ended %s", v.Task, id, step, status)
	}
	return finished, ""
}

type stage int

const (
	stRun stage = iota
	stFinish
	stDone
	stStop
)

// stage tells how far the last of the sessions of a stage has come: to run (none, or one that stopped), to
// finish (it ended and has no handoff), done (a valid handoff), or a stop (it failed, or its handoff is invalid).
func (v loopView) stage(ids []string, step, pair string) (stage, string, string) {
	if len(ids) == 0 {
		return stRun, "", ""
	}
	last := ids[len(ids)-1]
	switch o, msg := v.outcome(last, step); o {
	case running:
		return stRun, "", ""
	case failed:
		return stStop, last, msg
	}
	detail, found := v.lastHandoff(last)
	if !found {
		return stFinish, last, ""
	}
	if ok, _ := validHandoff(detail); !ok {
		return stStop, last, fmt.Sprintf("%s: the handoff %s is %s", v.Task, pair, detail)
	}
	return stDone, last, ""
}

// reviewRounds counts the rounds with findings: the verifiers whose valid handoff went to the developer.
func reviewRounds(v loopView) int {
	n := 0
	for _, s := range v.Sessions {
		rec, ok := v.Records[s[sessionID]]
		if _, verifier := v.qaKind(s); ok && verifier && rec.Verdict != nothingMaterial {
			n++
		}
	}
	return n
}

// findingsOf gives the findings that the developer of an attempt is given: those of the review record of the
// attempt before it that no retry replaced, else none and the bound head of that attempt, whose table of the gates gives them. A record
// of an older attempt is never used.
func findingsOf(v loopView, attempt int) (found []Finding, gateHead string) {
	prev := strconv.Itoa(priorAttempt(replacedOf(v), attempt))
	for _, e := range v.Events {
		if e[1] == "bound" && e[2] == prev {
			gateHead = e[5]
		}
	}
	for _, s := range v.Sessions {
		rec, ok := v.Records[s[sessionID]]
		if _, verifier := v.qaKind(s); ok && verifier && s[sessionAttempt] == prev && s[sessionBase] == gateHead && rec.Verdict != nothingMaterial {
			return rec.Findings, ""
		}
	}
	return nil, gateHead
}

// Replaced names the attempts that a retry replaced: the attempt number to the session whose eligible end caused it.
type Replaced map[string]string

// The stages of a task session that a retry bounds: the plan, the plan review, the developer and the verifier.
const (
	stagePlan, stageReview, stageDeveloper, stageVerifier = "plan", "review", "developer", "verifier"
)

// The classes of an eligible end: a crash or no-result (uncertain), and a done refused for the artifact alone.
const (
	classUncertain, classArtifact, classRulePath = "uncertain", "artifact", "rule-path"
)

// rulePathFindings gives one finding for each path of a refused payload that the register matches (rules.Match, so the unconditional .sh clause counts),
// in the order of the payload; an ordinary path gives none. The register is the one given, the current one (loop.md, The findings of a refused rule path;
// task T-dtw6, #491). docs/guardrails.md is named with the exception of its own register row, whether or not the refused change passed it.
func rulePathFindings(paths []string, reg []rules.Entry) []Finding {
	var out []Finding
	for _, p := range paths {
		e, ok := rules.Match(reg, p)
		if !ok {
			continue
		}
		text := "a rule path that the last attempt changed, so its session was refused and nothing of it was pushed: change nothing of it"
		if e.Exception != "" {
			text = "a rule path that the last attempt changed, so its session was refused and nothing of it was pushed: only " + e.Exception + " may change it"
		}
		out = append(out, Finding{Rule: classRulePath, File: p, Text: text})
	}
	return out
}

// stageOf gives the stage of a session by its start row and its role, "" for none: an answer session has none, and neither has the
// close-out or an owner role.
func (v loopView) stageOf(id string) string {
	s := v.startRow(id)
	if s == nil || v.answerSession(s) {
		return ""
	}
	switch s[sessionRole] {
	case stepRole("plan-of-a-task"):
		return stagePlan
	case stepRole(devStep(v.Kind)):
		return stageDeveloper
	case stepRole("verification"):
		if review, verifier := v.qaKind(s); review {
			return stageReview
		} else if verifier {
			return stageVerifier
		}
	}
	return ""
}

// eligibleEnd gives the last session of the attempt, its stage and the class of its end, when the end is one that a retry replaces
// (loop.md, The driver of a Shape task, "An uncertain end gets one new attempt"): "uncertain" for a crash or no-result with no refusal, and "artifact"
// for a done with the status completed whose refusals are exactly one event, artifact. A session that has no stage has none. "" for none.
func (v loopView) eligibleEnd(attempt string) (id, stage, class string) {
	for _, e := range v.Events {
		if e[1] == "session" && e[2] == attempt {
			id = e[3]
		}
	}
	end, refusals, detail := "", 0, ""
	for _, e := range v.Events {
		switch {
		case e[3] != id:
		case e[1] == "refused":
			refusals, detail = refusals+1, e[6]
		case e[1] == "result":
			end = e[6]
		}
	}
	switch {
	case id == "":
		return "", "", ""
	case refusals == 0 && (end == "crash" || end == "no-result"):
		class = classUncertain
	case refusals == 1 && detail == "artifact" && end == "done" && v.Status[id] == "completed":
		class = classArtifact
	default:
		return "", "", ""
	}
	return id, v.stageOf(id), class
}

// retries walks the attempts in order. An attempt is replaced by a retry when its last session has an eligible end, no earlier replacement
// was caused by an eligible end of the same stage and the same class (the bound is once per stage and class: crash and no-result share one retry,
// an artifact refusal has its own), it was not closed or rebased and holds no lost push, and the
// next attempt was recorded after that end at the same base. It gives the replaced attempts, and for the latest attempt, when it meets the
// rule and has no next attempt, the session whose end the next attempt would replace.
func retries(v loopView) (rep Replaced, pending string) {
	rep, spent, last := Replaced{}, map[[2]string]bool{}, 0
	for _, e := range v.Events {
		if n, _ := strconv.Atoi(e[2]); e[1] == "attempt" && n > last {
			last = n
		}
	}
	for n := 1; n <= last; n++ {
		a := strconv.Itoa(n)
		id, stage, class := v.eligibleEnd(a)
		if id == "" || stage == "" || spent[[2]string{stage, class}] {
			continue
		}
		base, ended, next := "", -1, -1
		for i, e := range v.Events {
			switch {
			case (e[1] == "closed" || e[1] == "rebased" || (e[1] == "refused" && e[6] == "push-lost")) && e[2] == a:
				stage = ""
			case e[1] == "attempt" && e[2] == a:
				base = e[4]
			case e[1] == "result" && e[3] == id:
				ended = i
			case e[1] == "attempt" && e[2] == strconv.Itoa(n+1) && next < 0:
				next = i
			}
		}
		switch {
		case stage == "":
		case next >= 0 && next > ended && v.Events[next][4] == base:
			rep[a], spent[[2]string{stage, class}] = id, true
		case n == last:
			pending = id
		}
	}
	return rep, pending
}

// replacedOf is the one rule of the replacement (loop.md, The driver of a Shape task, "An uncertain end gets one new attempt"), for the readers of "the attempt before".
func replacedOf(v loopView) Replaced { rep, _ := retries(v); return rep }

// priorAttempt is the nearest attempt before b that a retry did not replace, 0 for none: a retry stands for the attempt it replaces.
func priorAttempt(rep Replaced, b int) int {
	a := b - 1
	for a > 0 && rep[strconv.Itoa(a)] != "" {
		a--
	}
	return a
}

// rejectOf is the one rule of a rejection (loop.md, The driver of a Shape task, "A rejected plan goes back to the planner"), for counting, for
// reopening the plan stage and for the reasons the planner is given: a plan review that no retry replaced, whose result is completed, whose kept
// copy is the file the result names by hash and rejects with a reason (the view's Rejects), and whose latest handoff is valid. It gives the reasons.
func (v loopView) rejectOf(id string, rep Replaced) ([]string, bool) {
	reasons, found := v.Rejects[id]
	if !found || v.stageOf(id) != stageReview || rep[v.startRow(id)[sessionAttempt]] == id {
		return nil, false
	}
	if o, _ := v.outcome(id, "the plan review"); o != finished {
		return nil, false
	}
	detail, handedOff := v.lastHandoff(id)
	if ok, _ := validHandoff(detail); !handedOff || !ok {
		return nil, false
	}
	return reasons, true
}

// legacyRejectOf reports whether a plan review is a reject whose handoff an earlier binary recorded as plan review to developer, invalid
// verdict-field: it meets every condition of rejectOf but the handoff, which is to be recorded again, as the pair plan review to task plan.
func (v loopView) legacyRejectOf(id string, rep Replaced) bool {
	if _, found := v.Rejects[id]; !found || v.stageOf(id) != stageReview || rep[v.startRow(id)[sessionAttempt]] == id {
		return false
	}
	detail, _ := v.lastHandoff(id)
	o, _ := v.outcome(id, "the plan review")
	return o == finished && strings.HasPrefix(detail, "invalid verdict-field")
}

// planRejects counts the plan reviews that rejected the plan by the rule of rejectOf.
func planRejects(v loopView) int {
	rep, n := replacedOf(v), 0
	for _, s := range v.Sessions {
		if _, ok := v.rejectOf(s[sessionID], rep); ok {
			n++
		}
	}
	return n
}

// reopenedBy gives the highest attempt before this one in which a plan review rejected the plan (rejectOf), 0 for none, and its reasons: the plan
// stage runs again in the attempts after it, and the sessions of the plan and of the plan review of the attempts up to it are passed over.
func reopenedBy(v loopView, rep Replaced, attempt int) (int, []string) {
	best, reasons := 0, []string(nil)
	for _, s := range v.Sessions {
		if r, ok := v.rejectOf(s[sessionID], rep); ok {
			if a, _ := strconv.Atoi(s[sessionAttempt]); a < attempt && a >= best {
				best, reasons = a, r
			}
		}
	}
	return best, reasons
}

// rejectionOf gives the reasons of the plan review that rejected the plan in the latest attempt before this one, as the findings the planner of this
// attempt is given; none for the first plan.
func rejectionOf(v loopView, attempt int) []Finding {
	_, reasons := reopenedBy(v, replacedOf(v), attempt)
	var out []Finding
	for _, r := range reasons {
		out = append(out, Finding{Rule: "plan-review:reject", Text: r})
	}
	return out
}

// loopPosition gives the first step of the table of loop.md, The task loop, that the records show not done.
func loopPosition(v loopView) loopStep {
	attempt, base := 0, ""
	for _, e := range v.Events {
		if n, _ := strconv.Atoi(e[2]); e[1] == "attempt" && n >= attempt {
			attempt, base = n, e[4]
		}
	}
	if attempt == 0 {
		return loopStep{Pos: posAttempt}
	}
	step := loopStep{Attempt: attempt, Base: base}
	stop := func(msg string) loopStep { step.Pos, step.Msg = posStop, msg; return step }
	// Only the developer's session is ever bound (a session that reads has no push): a bound of any other stops the step.
	for _, e := range v.Events {
		if e[1] == "bound" && e[3] != "" && v.role(e[3]) != stepRole(devStep(v.Kind)) { // a bound of code has no session
			return stop(fmt.Sprintf("%s: the session %s has a bound and is not of the developer role", v.Task, e[3]))
		}
	}
	// A start that was refused has an event refused and no start row: no stage would ever see it.
	since := 0
	for i, e := range v.Events {
		if n, _ := strconv.Atoi(e[2]); e[1] == "attempt" && n == attempt {
			since = i
		}
	}
	for _, e := range v.Events[since:] {
		if e[1] == "refused" && !slices.ContainsFunc(v.Sessions, func(r []string) bool { return r[sessionID] == e[3] }) {
			return stop(fmt.Sprintf("%s: the start of the session %s was refused (%s)", v.Task, e[3], e[6]))
		}
	}
	if q := questionPosition(v, attempt); q.Stage != qNone { // a session of this attempt asked: the stages of questionloop.go
		step.Pos, step.Ask = posQuestion, q
		return step
	}
	rep, pending := retries(v)
	if pending != "" { // the last session of the attempt has an eligible end (crash, no-result or a refused artifact): one new attempt at the same base
		step.Pos, step.Session = posRetry, pending
		return step
	}
	reopened, _ := reopenedBy(v, rep, attempt) // the plan and the plan review of the attempts up to it are passed over
	var plans, reviews, devs, verifiers []string
	for _, s := range v.Sessions {
		id, a := s[sessionID], s[sessionAttempt]
		n, _ := strconv.Atoi(a)
		if rep[a] == id || v.answerSession(s) || (a != strconv.Itoa(attempt) && v.Status[id] == "needs_context") {
			continue // the answer session of an attempt, and a session that asked in an earlier attempt: its stage runs again in this one
		}
		switch s[sessionRole] {
		case stepRole("plan-of-a-task"):
			if n > reopened {
				plans = append(plans, id)
			}
		case stepRole("verification"):
			review, verifier := v.qaKind(s)
			switch {
			case review && n <= reopened:
			case review:
				reviews = append(reviews, id)
			case verifier && a == strconv.Itoa(attempt):
				verifiers = append(verifiers, id)
			case !verifier:
				return stop(fmt.Sprintf("%s: the session %s is neither the plan review nor a verifier", v.Task, id))
			}
		case stepRole(devStep(v.Kind)):
			if a == strconv.Itoa(attempt) {
				devs = append(devs, id)
			}
		}
	}
	for _, c := range []struct {
		ids        []string
		step, pair string
		run, end   position
	}{{plans, "the plan", "task plan → plan review", posPlan, posPlanFinish}, {reviews, "the plan review", "plan review → developer", posReview, posReviewFinish}} {
		pair := c.pair
		if c.end == posReviewFinish && len(c.ids) > 0 {
			last := c.ids[len(c.ids)-1]
			if detail, _ := v.lastHandoff(last); v.Verdicts[last] == "reject" && !strings.HasPrefix(detail, "invalid verdict-field") {
				pair = "plan review → task plan" // a reject is checked on the pair to the planner; an earlier binary's failure keeps the developer pair
			}
		}
		st, id, msg := v.stage(c.ids, c.step, pair)
		step.Session = id
		if c.end == posReviewFinish && st == stDone && v.Unread[id] {
			return stop(fmt.Sprintf("%s: the plan review %s cannot be read from its kept copy", v.Task, id))
		}
		if c.end == posReviewFinish {
			switch _, rejected := v.rejectOf(id, rep); {
			case st == stDone && !rejected && v.Verdicts[id] == "reject": // a reject with no reason, whose handoff the Operator's table admitted: never an approval
				return stop(fmt.Sprintf("%s: the plan review %s rejected the plan and gave no reason (its handoff is valid on this table, and a reject with no reason stops)", v.Task, id))
			case st == stStop && v.legacyRejectOf(id, rep): // an earlier binary's invalid verdict-field of a reject: its handoff is recorded again
				step.Pos = posReviewFinish
				return step
			case st == stDone && rejected && planRejects(v) <= shapeCycleCap:
				step.Pos = posReplan
				return step
			case st == stDone && rejected:
				return stop(fmt.Sprintf("%s: the plan review rejected the plan %d times", v.Task, planRejects(v)))
			}
		}
		switch st {
		case stRun:
			step.Pos = c.run
			return step
		case stFinish:
			step.Pos = c.end
			return step
		case stStop:
			return stop(msg)
		}
	}
	if len(devs) == 0 {
		step.Pos = posDeveloper
		return step
	}
	last := devs[len(devs)-1]
	step.Session = last
	switch o, msg := v.outcome(last, devStep(v.Kind)); o {
	case running:
		step.Pos = posDeveloper
		return step
	case failed:
		return stop(msg)
	}
	pushed := false
	for _, e := range v.Events {
		if e[3] == last && e[1] == "push" {
			pushed = true
		}
		if e[3] == last && e[1] == "bound" {
			step.Head = e[5]
		}
	}
	switch {
	case step.Head == "" && pushed:
		return stop(fmt.Sprintf("%s: a push with no result (its repair is the step loop's)", v.Task))
	case step.Head == "":
		return stop(fmt.Sprintf("%s: the session %s ended with no push (its resume is the step loop's)", v.Task, last))
	}
	pair := devFrom(v.Kind) + " → verifier"
	detail, found := v.lastHandoff(last)
	if ok, _ := validHandoff(detail); !found {
		step.Pos = posToVerifier
		return step
	} else if !ok {
		return stop(fmt.Sprintf("%s: the handoff %s is %s", v.Task, pair, detail))
	}
	verifiers = slices.DeleteFunc(verifiers, func(id string) bool {
		i := slices.IndexFunc(v.Sessions, func(s []string) bool { return s[sessionID] == id })
		return v.Sessions[i][sessionBase] != step.Head
	})
	var rec ReviewRecord
	var haveRec bool
	if n := len(verifiers); n > 0 {
		rec, haveRec = v.Records[verifiers[n-1]]
	}
	st, id, msg := v.stage(verifiers, "the verification", "verifier")
	step.Session = id
	switch st {
	case stRun:
		step.Pos = posShip
		return step
	case stFinish:
		step.Pos = posVerifierFinish
		return step
	case stStop:
		return stop(msg)
	}
	switch {
	case !haveRec:
		return stop(fmt.Sprintf("%s: the review record of the session %s has no kept copy", v.Task, id))
	case rec.Verdict == nothingMaterial && len(rec.Findings) == 0:
		step.Pos = posVerified
	case rec.Verdict == "not mergeable, findings recorded":
		step.Pos, step.Msg = posWait, fmt.Sprintf("waiting for the Operator: %s: not mergeable, findings recorded", v.Task)
	case reviewRounds(v) <= shapeCycleCap:
		step.Pos = posNextAttempt
	default:
		step.Pos, step.Msg = posWait, fmt.Sprintf("waiting for the Operator: %s: the review-cycle cap %d is reached with findings open", v.Task, shapeCycleCap)
	}
	return step
}

// readPlanReview reads the plan review of a session from its result: review is whether the result names a comment file plan-review.md (a verifier's
// does not), and verdict is its verdict word when the kept copy is the file the result names by hash and has the plan review form, else "" (the
// review is unread: the copy is absent, of other bytes, or no review). data are the bytes read.
func (s *Sessions) readPlanReview(ctx context.Context, task, id string, result [][]string) (review bool, verdict string, data []byte, err error) {
	if !slices.ContainsFunc(result, func(r []string) bool { return r[0] == "comment" && r[2] == "result/"+PlanReviewFile }) {
		return false, "", nil, nil
	}
	data, ok, err := s.keptComment(ctx, task, id, PlanReviewFile, result)
	if err != nil || !ok {
		return true, "", nil, err
	}
	verdict, _ = ParsePlanReview(data)
	return true, verdict, data, nil
}

// shapeView reads the records that the position is told from.
func (s *Sessions) shapeView(ctx context.Context, task string) (loopView, error) {
	t, err := s.taskRow(ctx, task)
	if err != nil {
		return loopView{}, err
	}
	v := loopView{Task: task, Kind: t.Kind, Status: map[string]string{}, Records: map[string]ReviewRecord{}, Results: map[string][][]string{}, Issue: t.Issue, Rejects: map[string][]string{}, Unread: map[string]bool{}, Verdicts: map[string]string{}}
	if v.Events, err = s.readTable(ctx, eventsPath(task), records.ReadEvents); err != nil {
		return v, err
	}
	all, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return v, err
	}
	for _, r := range all {
		if r[sessionTask] == task {
			v.Sessions = append(v.Sessions, r)
		}
	}
	for _, e := range v.Events {
		if e[1] != "result" || e[6] != "done" {
			continue
		}
		rows, err := s.readTable(ctx, "tasks/"+task+"/results/"+e[3]+".tsv", records.ReadResult)
		if err != nil {
			return v, err
		}
		v.Results[e[3]] = rows
		for _, r := range rows {
			if r[0] == "status" {
				v.Status[e[3]] = r[2]
			}
		}
	}
	for _, sess := range v.Sessions { // the reasons of each plan review that rejected: its kept copy is the file its result names by hash
		if id := sess[sessionID]; sess[sessionRole] == stepRole("verification") && v.Status[id] == "completed" {
			review, verdict, data, err := s.readPlanReview(ctx, task, id, v.Results[id])
			switch {
			case err != nil:
				return v, err
			case !review:
			case verdict == "":
				v.Unread[id] = true // a valid handoff over a copy that cannot be read is no approval
			default:
				v.Verdicts[id] = verdict
				if RejectReasonsOK(data) {
					v.Rejects[id] = PlanReviewReasons(data)
				}
			}
		}
	}
	copies, err := s.readTable(ctx, "copies.tsv", records.ReadCopies)
	if err != nil {
		return v, err
	}
	v.Copies = copies
	for _, sess := range v.Sessions {
		id := sess[sessionID]
		detail, found := v.lastHandoff(id)
		var cid, seen int
		if _, err := fmt.Sscanf(detail, "valid comment %d-%d", &cid, &seen); !found || err != nil || sess[sessionRole] != stepRole("verification") {
			continue
		}
		for _, c := range copies {
			if c[copyComment] != strconv.Itoa(cid) || c[copySeen] != strconv.Itoa(seen) || c[copyIssue] != strconv.Itoa(t.Issue) {
				continue
			}
			data, err := s.Store.ReadFile(ctx, c[copyBody])
			if err != nil {
				return v, err
			}
			if rec, ok := ParseReviewRecord(data); ok {
				v.Records[id] = rec
			}
		}
	}
	return v, nil
}

// AdvanceShapeTask takes one Shape task to its verified head (loop.md, The driver of a Shape task): it repairs
// the task's sessions, then acts on the position of its records and reads them again, until the head is verified,
// the Operator is to answer, or a case stops it.
func (s *Sessions) AdvanceShapeTask(ctx context.Context, in ShapeTaskInput) (Advance, error) {
	s.fns()
	if _, err := s.RepairTask(ctx, in.Task); err != nil {
		return Advance{}, err
	}
	for i := 0; i < maxActs; i++ {
		v, err := s.shapeView(ctx, in.Task)
		if err != nil {
			return Advance{}, err
		}
		st := loopPosition(v)
		var stop, wait string
		switch st.Pos {
		case posStop:
			return Advance{Task: in.Task, Stop: st.Msg}, nil
		case posWait:
			if err := s.settle(ctx, in, v, st); err != nil {
				return Advance{}, err
			}
			return Advance{Task: in.Task, Wait: st.Msg}, nil
		case posQuestion:
			stop, wait, err = s.askStage(ctx, in, v, st)
		case posAttempt:
			_, err = s.StartAttempt(ctx, in.Task, firstBase(v.Events, in.Base))
		case posRetry, posReplan:
			_, err = s.StartAttempt(ctx, in.Task, st.Base)
		case posPlan, posReview, posDeveloper:
			stop, err = s.runStage(ctx, in, v, st)
		case posPlanFinish, posReviewFinish, posVerifierFinish:
			err = s.finishStage(ctx, v, st)
		case posToVerifier:
			err = s.toVerifier(ctx, v, st)
		case posShip:
			stop, err = s.ship(ctx, in, v, st)
		case posVerified, posNextAttempt:
			if st.Pos == posVerified && renderedPushed(v.Events, st.Attempt) { // the rendered commit is the pull request's head: the verified head was settled before
				return Advance{Task: in.Task, Verified: true, Head: st.Head}, nil
			}
			if stop, err = s.settleVerdict(ctx, in, v, st); err == nil && stop == "" && st.Pos == posNextAttempt {
				_, err = s.StartAttempt(ctx, in.Task, st.Base)
			}
			if err == nil && stop == "" && st.Pos == posVerified {
				return Advance{Task: in.Task, Verified: true, Head: st.Head}, nil
			}
		default:
			return Advance{}, fmt.Errorf("%s: the position %d is no step of the driver", in.Task, st.Pos)
		}
		var refusal session.Refusal
		var posting HandoffRefusal
		switch {
		case stop != "":
			return Advance{Task: in.Task, Stop: stop}, nil
		case wait != "":
			return Advance{Task: in.Task, Wait: wait}, nil
		case errors.As(err, &refusal), errors.As(err, &posting):
			// The refusal is recorded; the next position names it.
		case err != nil:
			return Advance{}, err
		}
	}
	return Advance{}, fmt.Errorf("%s: the driver did not come to a verified head, a wait or a stop in %d acts", in.Task, maxActs)
}

// startSession starts one session of the step through the one call of TaskSession: the cap of milestone 0, the
// prompt file at the base, the admission of the pair. A base that is not on the default branch is made reachable
// first: the scratch branch layup/base/<task>/<attempt> of the run's clone is set to it, and the session clones that
// branch (plan amendment 1 of #415).
func (s *Sessions) startSession(ctx context.Context, in ShapeTaskInput, kind, step string, attempt int, base string, pair, model []string, findings []Finding, scratch bool) error {
	return s.startAs(ctx, in, kind, step, stepRole(step), attempt, base, pair, model, findings, scratch, nil)
}

// startAs is startSession for a session whose role is not the step table's: the answer session is of the owner role of the question
// (row 66z), the table's word for it being owner-of-its-kind. push is the session's push hook, nil for the checks of any task session
// (the close-out of row 66y has its own).
func (s *Sessions) startAs(ctx context.Context, in ShapeTaskInput, kind, step, role string, attempt int, base string, pair, model []string, findings []Finding, scratch bool, push func(id string, d session.Dir, r session.Run) error) error {
	p := pairOf(pair, model)
	if err := s.milestoneCap(ctx, p.Cap); err != nil {
		return err
	}
	rec := in.Records()
	pin := PromptIn{Task: in.Task, Step: step, Base: base, Findings: findings}
	if kind == "architecture" {
		pin.Options = in.Options
	}
	pc, err := s.PromptFile(ctx, pin)
	if err != nil {
		return err
	}
	branch := ""
	if scratch {
		branch = fmt.Sprintf("layup/base/%s/%d", in.Task, attempt)
		if err := git.UpdateRef(s.Clone, "refs/heads/"+branch, base); err != nil {
			return err
		}
	}
	// The admission may run a probe at the same base, so it clones the same branch.
	ts := TaskSpec{Task: in.Task, Attempt: attempt, Role: role, Base: base, Branch: branch, Records: rec, Pair: p, Context: pc,
		Admit: s.Admit(ctx, ProbeSpec{Harness: pair, Models: in.Models, Base: base, Records: rec, Branch: branch}), Push: push}
	_, err = s.TaskSession(ctx, ts)
	return err
}

// pairStop is the stop of a step that has no admitted pair.
func pairStop(task, step string, err error) (string, error) {
	if errors.Is(err, route.ErrNoPair) {
		return fmt.Sprintf("%s: no admitted pair for the step %s", task, step), nil
	}
	return "", err
}

// runStage starts the session of the plan, the plan review or the developer.
func (s *Sessions) runStage(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (string, error) {
	step := map[position]string{posPlan: "plan-of-a-task", posReview: "plan-review", posDeveloper: devStep(v.Kind)}[st.Pos]
	res, err := ResolveStep(step, "", "")
	if err != nil {
		return "", err
	}
	var h, m []string
	var findings []Finding
	switch st.Pos {
	case posPlan:
		findings = rejectionOf(v, st.Attempt) // the reasons of a plan review that rejected the plan of an attempt before this one
		h, m, err = s.PairFor(ctx, res.Role, res.Tier, in.Harnesses, in.Models, nil)
	case posReview:
		if h, m, err = s.ReviewerPair(ctx, in.Task, PlanReview, st.Base, "", in.Harnesses, in.Models); errors.Is(err, ErrNoHarnessOutsideAuthors) {
			return fmt.Sprintf("%s: no harness outside the authors of the plan", in.Task), nil
		}
	case posDeveloper:
		var gateHead string
		if findings, gateHead = findingsOf(v, st.Attempt); gateHead != "" {
			var stop string
			if findings, stop, err = s.headFindings(ctx, v.Task, v.Kind, st.Base, gateHead); err != nil || stop != "" {
				return stop, err
			}
		}
		fallthrough
	default:
		h, m, err = s.PairFor(ctx, res.Role, res.Tier, in.Harnesses, in.Models, nil)
	}
	if err != nil {
		return pairStop(in.Task, step, err)
	}
	return "", s.startSession(ctx, in, v.Kind, step, st.Attempt, st.Base, h, m, findings, true)
}

// headFindings gives the findings that the developer of the next attempt is given from the head that failed
// its checks and was never judged by a review: the rows of the head's gate table that did not pass, then the findings
// of the other checks. A stop of those checks comes before the findings.
func (s *Sessions) headFindings(ctx context.Context, task, kind, base, head string) ([]Finding, string, error) {
	table, found, err := s.GatesTable(ctx, task, head)
	if err != nil {
		return nil, "", err
	}
	var out []Finding
	if found {
		out = RetryGateFindings(kind, table)
	}
	check, err := s.CheckFindings(ctx, HeadInput{Task: task, Base: base, Head: head})
	if err != nil || check.Stop != "" {
		return nil, taskStop(task, check.Stop), err
	}
	return append(out, check.Findings...), "", nil
}

// taskStop prefixes the stop of a check with the task, and keeps "" for none.
func taskStop(task, stop string) string {
	if stop == "" {
		return ""
	}
	return task + ": " + stop
}

// finishStage posts the comment of a session that ended and records its handoff; both read before they write, so
// a stage that an interruption left half done is finished and its session is not run again.
func (s *Sessions) finishStage(ctx context.Context, v loopView, st loopStep) error {
	file, from, to, kind := PlanFile, "task plan", "plan review", "plan"
	switch st.Pos {
	case posReviewFinish:
		file, from, to, kind = PlanReviewFile, "plan review", "developer", "plan-review"
		data, err := s.Store.ReadFile(ctx, KeptPath(v.Task, st.Session, file))
		if err != nil {
			return err
		}
		if verdict, ok := ParsePlanReview(data); ok && verdict == "reject" {
			to = "task plan" // a reject goes back to the planner; the item reject-reasons says whether it gave a reason
		}
	case posVerifierFinish:
		file, from, kind = ReviewRecordFile, "verifier", "review-record"
		data, err := s.Store.ReadFile(ctx, KeptPath(v.Task, st.Session, file))
		if err != nil {
			return err
		}
		rec, ok := ParseReviewRecord(data)
		to = "developer"
		if ok && rec.Verdict == nothingMaterial {
			to = "close-out"
		}
	}
	if _, err := s.PostComment(ctx, v.Task, st.Session, file); err != nil {
		return err
	}
	_, err := s.RecordHandoff(ctx, HandoffInput{Task: v.Task, Session: st.Session, From: from, To: to, File: file, Kind: kind})
	return err
}

// toVerifier records the handoff of the developer's bound head: the commit, and the task file as the result named it
// and as the head holds it.
func (s *Sessions) toVerifier(ctx context.Context, v loopView, st loopStep) error {
	path := "docs/tasks/" + v.Task + ".md"
	rows, err := s.readTable(ctx, "tasks/"+v.Task+"/results/"+st.Session+".tsv", records.ReadResult)
	if err != nil {
		return err
	}
	named := ""
	for _, r := range rows {
		if r[0] == "artifact" && r[2] == path {
			named = r[3]
		}
	}
	found, err := s.blobHash(st.Head, path)
	if err != nil {
		return err
	}
	_, err = s.RecordHandoff(ctx, HandoffInput{Task: v.Task, Session: st.Session, From: devFrom(v.Kind), To: "verifier",
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: st.Head, Found: st.Head}, {Kind: "task-file", Path: path, Named: named, Found: found}}})
	return err
}

// head gives the four statuses of the bound head with the verifier of the session (none: it waits for its turn).
func (s *Sessions) head(ctx context.Context, in ShapeTaskInput, st loopStep, session string) ([]StatusPost, error) {
	return s.PostStatuses(ctx, HeadInput{Task: in.Task, Base: st.Base, Head: st.Head, Copy: in.Copy, Session: session, Waiting: session == "", Harnesses: in.Harnesses, Models: in.Models})
}

// shapeDescription is the description of the title of the pull request of each kind.
var shapeDescription = map[string]string{"specification": "the specification of the target", "architecture": "the architecture of the target"}

// ship opens the draft pull request of the bound head, posts its four statuses and judges them, accepts a
// specification that waits for its bet, and starts the verifier. An error or the status not-active is a stop, and no
// attempt starts; a failure of layup/gates, layup/spec or layup/rules starts the next attempt with the findings of all three.
// It reads the same records again at a restart, so a stop between the statuses and the attempt starts exactly one.
func (s *Sessions) ship(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (string, error) {
	if _, err := s.OpenPullRequest(ctx, PullInput{Task: in.Task, Attempt: st.Attempt, Type: "docs", Description: shapeDescription[v.Kind]}); err != nil {
		return "", err
	}
	posts, err := s.head(ctx, in, st, "")
	if err != nil {
		return "", err
	}
	stop, gates, checks := judgeStatuses(in.Task, posts)
	if stop != "" {
		return stop, nil
	}
	var check CheckResult
	if checks {
		if check, err = s.CheckFindings(ctx, HeadInput{Task: in.Task, Base: st.Base, Head: st.Head}); err != nil || check.Stop != "" {
			return taskStop(in.Task, check.Stop), err
		}
	}
	if gates || checks {
		var table gate.Table
		found := true
		if gates {
			if table, found, err = s.GatesTable(ctx, in.Task, st.Head); err != nil {
				return "", err
			}
		}
		if _, nErr := s.NextAttempt(ctx, NextAttemptInput{Task: in.Task, Kind: v.Kind, Base: st.Base, Table: table, Findings: check.Findings}); !found || errors.Is(nErr, ErrNoFailure) {
			return fmt.Sprintf("%s: a layup/ status failed and the records hold no failed row of it", in.Task), nil
		} else if nErr != nil {
			return "", nErr
		}
		return "", nil
	}
	reopened := false
	if v.Kind == "specification" {
		var err error
		if reopened, err = s.reopenedTask(ctx, in.Task); err != nil {
			return "", err
		}
	}
	if v.Kind == "specification" && !reopened { // a reopened task is judged by layup/spec and by the target's prd-lint at the ready rule (D9)
		ok, lines, err := s.waitingSpecAccepted(st.Base, st.Head)
		if err != nil {
			return "", err
		}
		if !ok {
			return fmt.Sprintf("%s: the pinned prd-lint does not accept the head as a specification that waits for its bet: %q", in.Task, lines), nil
		}
	}
	h, m, err := s.ReviewerPair(ctx, in.Task, Verification, st.Base, st.Head, in.Harnesses, in.Models)
	if errors.Is(err, ErrNoHarnessOutsideAuthors) {
		return fmt.Sprintf("%s: not-active: %s", in.Task, ErrNoHarnessOutsideAuthors), nil
	} else if err != nil {
		return pairStop(in.Task, "verification", err)
	}
	return "", s.startSession(ctx, in, v.Kind, "verification", st.Attempt, st.Head, h, m, nil, true)
}

// judgeStatuses judges the four statuses posted at a bound head: the first stop (an error, or the status not-active of
// layup/verify) comes before everything, else whether layup/gates failed and whether layup/spec or layup/rules did.
func judgeStatuses(task string, posts []StatusPost) (stop string, gates, checks bool) {
	for _, p := range posts {
		switch name, state, d := p.Status.Context, p.Status.State, p.Status.Description; {
		case state == "error":
			return fmt.Sprintf("%s: %s is error: %s", task, name, d), false, false
		case state == "failure" && name == VerifyContext && len(d) > len("not-active: ") && d[:len("not-active: ")] == "not-active: ":
			return fmt.Sprintf("%s: %s", task, d), false, false
		case state == "failure" && name == GatesContext:
			gates = true
		case state == "failure" && (name == SpecContext || name == RulesContext):
			checks = true
		}
	}
	return "", gates, checks
}

// waitingSpecAccepted runs the pinned prd-lint of the target (its script at the base) on the docs/prd/ and docs/facts/
// of the head, in a scratch directory, and gives spec.WaitingSpecAccepted's answer.
func (s *Sessions) waitingSpecAccepted(base, head string) (bool, []string, error) {
	script, err := git.Show(s.Clone, base, "docs/prd/prd-lint.sh")
	if err != nil {
		return false, nil, err
	}
	root, err := os.MkdirTemp("", "layup-prdlint-")
	if err != nil {
		return false, nil, err
	}
	defer os.RemoveAll(root)
	for _, dir := range []string{"docs/prd", "docs/facts"} {
		files, err := filesAt(s.Clone, head, dir, func(string) bool { return true })
		if err != nil {
			return false, nil, err
		}
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			return false, nil, err
		}
		for name, data := range files {
			if err := os.WriteFile(filepath.Join(root, dir, name), data, 0o644); err != nil {
				return false, nil, err
			}
		}
	}
	path := filepath.Join(root, "prd-lint.sh")
	if err := os.WriteFile(path, script, 0o644); err != nil {
		return false, nil, err
	}
	return spec.WaitingSpecAccepted(path, filepath.Join(root, "docs", "prd"))
}

// settle links the review record in the body of the pull request and posts the four statuses with the verifier's verdict;
// every one reads the forge before it writes.
func (s *Sessions) settle(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) error {
	_, err := s.settleStatuses(ctx, in, v, st)
	return err
}

func (s *Sessions) settleStatuses(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) ([]StatusPost, error) {
	if _, err := s.LinkReviewRecord(ctx, in.Task, st.Session); err != nil {
		return nil, err
	}
	return s.head(ctx, in, st, st.Session)
}

// settleVerdict settles a verdict that ends the attempt, and for the verified head requires layup/verify to be a success.
func (s *Sessions) settleVerdict(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (string, error) {
	posts, err := s.settleStatuses(ctx, in, v, st)
	if err != nil || st.Pos != posVerified {
		return "", err
	}
	for _, p := range posts {
		if p.Status.Context == VerifyContext && p.Status.State != "success" {
			return fmt.Sprintf("%s: %s is %s: %s", in.Task, VerifyContext, p.Status.State, p.Status.Description), nil
		}
	}
	return "", nil
}

// readOnlyStep reports whether the session of a step only reads: it has no push.
func readOnlyStep(step string) bool {
	return step == "plan-of-a-task" || step == "plan-review" || step == "verification" || step == "answer-to-a-question"
}
