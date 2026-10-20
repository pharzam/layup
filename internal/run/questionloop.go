package run

// This file holds the question in the loop (docs/spec/loop.md, A question during the work, "The question in the loop"; task
// T-r7b9, #426, row 66z of the plan of M2e): the stages that take a session that ended needs_context to its owner role's answer and
// to the next attempt, as a position of the driver of a Shape task.

import (
	"context"
	"fmt"
	"strconv"
)

// A qStage is the first stage of the question that the records show not done.
type qStage int

const (
	qNone          qStage = iota // no asking session in the latest attempt
	qHandoffAsk                  // P1: record the handoff asking role to owner role of an asking session
	qPostQuestion                // P2: post and copy the comment of a question row
	qWait                        // P3: the questions are not of one owner role: wait for the Operator
	qRunOwner                    // P4: run the answer session
	qPostAnswer                  // P6: post and copy the comment of an answer row
	qHandoffAnswer               // P7: record the handoff owner role to asking role
	qNextAttempt                 // P8: start the next attempt
	qStop                        // a stop: P1, P5, P7, or a row that cannot be composed
)

// A qStep is a stage with what its act needs: the session it concerns (the asking session of P1 and P2, the answer session of P6 and P7),
// the row and, for an answer, the question it answers, the owner role and the base of the answer session, the stop or the wait.
type qStep struct {
	Stage    qStage
	Session  string
	Row      []string
	Question string
	Owner    string
	Base     string
	Msg      string
}

// An askedQuestion is a question row of an asking session with its ID, <session>.<n>.
type askedQuestion struct {
	ID, Session string
	Row         []string
}

// startRow gives the row of sessions.tsv of a session, nil for none.
func (v loopView) startRow(id string) []string {
	for _, s := range v.Sessions {
		if s[sessionID] == id {
			return s
		}
	}
	return nil
}

// baseOf gives the base of the start row of a session, "" for none.
func (v loopView) baseOf(id string) string {
	if r := v.startRow(id); r != nil {
		return r[sessionBase]
	}
	return ""
}

// asking gives the asking sessions of an attempt in the order of their start events, and the place of the first event result that
// ended one: the sessions of the attempt that ended done with the status needs_context and started before that place. A session of
// the attempt that started after it is an answer session (answers).
func (v loopView) asking(attempt string) (ids []string, first int) {
	first = -1
	for i, e := range v.Events {
		if e[1] == "result" && e[2] == attempt && e[6] == "done" && v.Status[e[3]] == "needs_context" {
			first = i
			break
		}
	}
	if first < 0 {
		return nil, -1
	}
	for _, e := range v.Events[:first] {
		if e[1] != "session" || e[2] != attempt {
			continue
		}
		if i, done := resultEnd(v.Events, e[3]); i >= 0 && done && v.Status[e[3]] == "needs_context" {
			ids = append(ids, e[3])
		}
	}
	return ids, first
}

// answerSession reports whether the session is the answer session of its attempt: it started after the first result that ended needs_context.
func (v loopView) answerSession(s []string) bool {
	_, first := v.asking(s[sessionAttempt])
	i := v.startIndex(s[sessionID])
	return first >= 0 && i > first
}

// answers gives the sessions of the attempt that started after the first asking result, in the order of their starts: the answer sessions.
func (v loopView) answers(attempt string, first int) []string {
	var ids []string
	if first < 0 {
		return nil
	}
	for _, e := range v.Events[first+1:] {
		if e[1] == "session" && e[2] == attempt {
			ids = append(ids, e[3])
		}
	}
	return ids
}

// questions gives the one list of the questions of the asking sessions: the sessions in order, the rows of each result in the order of its file.
func (v loopView) questions(asking []string) []askedQuestion {
	var out []askedQuestion
	for _, id := range asking {
		for _, r := range v.Results[id] {
			if r[0] == "question" {
				out = append(out, askedQuestion{ID: id + "." + r[1], Session: id, Row: r})
			}
		}
	}
	return out
}

// questionPosition tells the stage of the question of the attempt from the view alone: the first of P1 to P8 that the records show not done.
func questionPosition(v loopView, attempt int) qStep {
	a := strconv.Itoa(attempt)
	asking, first := v.asking(a)
	if len(asking) == 0 {
		return qStep{}
	}
	stop := func(format string, args ...any) qStep {
		return qStep{Stage: qStop, Msg: v.Task + ": " + fmt.Sprintf(format, args...)}
	}
	for _, id := range asking { // P1
		detail, found := v.lastHandoff(id)
		if !found {
			return qStep{Stage: qHandoffAsk, Session: id}
		}
		if ok, _ := validHandoff(detail); !ok {
			return stop("the handoff asking role → owner role is %s", detail)
		}
	}
	qs := v.questions(asking)
	if len(qs) == 0 {
		return stop("the session %s ended needs_context and asked no question", asking[0])
	}
	for _, q := range qs { // P2
		body, err := QuestionComment(v.startRow(q.Session), q.Row)
		if err != nil {
			return stop("the question %s cannot be composed: %v", q.ID, err)
		}
		if _, ok := copyNamed(v.Copies, v.Issue, digest(body)); !ok {
			return qStep{Stage: qPostQuestion, Session: q.Session, Row: q.Row}
		}
	}
	ans := v.answers(a, first)
	cut := len(v.Events) // P3: the owner role of the questions, from the results that asked and no other
	if len(ans) > 0 {
		cut = v.startIndex(ans[0])
	}
	owner, ok := answerOwner(v.Events[:cut], v.Results, a, Owners(DefaultOwners()))
	if !ok {
		return qStep{Stage: qWait, Msg: "waiting for the Operator: " + v.Task + ": question " + qs[0].ID}
	}
	last := ""
	if len(ans) > 0 {
		last = ans[len(ans)-1]
	}
	if last == "" {
		return qStep{Stage: qRunOwner, Owner: owner, Base: v.baseOf(asking[len(asking)-1])} // P4
	}
	if role := v.role(last); role != owner {
		return stop("the session %s started after the question is of the role %s, not of the owner role %s", last, role, owner)
	}
	switch o, msg := v.outcome(last, "answer-to-a-question"); o { // P4, P5
	case running:
		return qStep{Stage: qRunOwner, Owner: owner, Base: v.baseOf(asking[len(asking)-1])}
	case failed:
		return qStep{Stage: qStop, Msg: msg}
	}
	var rows [][]string
	for _, r := range v.Results[last] {
		if r[0] == "answer" {
			rows = append(rows, r)
		}
	}
	have := make([]bool, len(qs))
	for _, r := range rows {
		n, err := strconv.Atoi(r[1])
		if err != nil || n < 1 || n > len(qs) || have[n-1] {
			return stop("the answer %s.%s answers no question", last, r[1])
		}
		have[n-1] = true
	}
	for i, ok := range have {
		if !ok {
			return stop("the question %s has no answer", qs[i].ID)
		}
	}
	start := v.startRow(last)
	for _, r := range rows { // P6
		n, _ := strconv.Atoi(r[1])
		qid := qs[n-1].ID
		body, err := AnswerComment(start, r, qid)
		if err != nil {
			return stop("the answer %s.%s cannot be composed: %v", last, r[1], err)
		}
		if _, ok := copyNamed(v.Copies, v.Issue, digest(body)); !ok {
			return qStep{Stage: qPostAnswer, Session: last, Row: r, Question: qid}
		}
	}
	detail, found := v.lastHandoff(last) // P7
	if !found {
		return qStep{Stage: qHandoffAnswer, Session: last}
	}
	if ok, _ := validHandoff(detail); !ok {
		return stop("the handoff owner role → asking role is %s", detail)
	}
	return qStep{Stage: qNextAttempt} // P8
}

// askStage does the act of a stage of the question; the driver reads the records again after it. It gives the stop or the wait of the stage, if any.
func (s *Sessions) askStage(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (stop, wait string, err error) {
	q := st.Ask
	switch q.Stage {
	case qStop:
		return q.Msg, "", nil
	case qWait:
		return "", q.Msg, nil
	case qHandoffAsk:
		_, err = s.RecordHandoff(ctx, HandoffInput{Task: in.Task, Session: q.Session, From: "asking role", To: "owner role"})
	case qPostQuestion:
		var body []byte
		if body, err = QuestionComment(v.startRow(q.Session), q.Row); err == nil {
			_, err = s.postBody(ctx, in.Task, body)
		}
	case qRunOwner:
		stop, err = s.runOwner(ctx, in, v, st)
	case qPostAnswer:
		var body []byte
		if body, err = AnswerComment(v.startRow(q.Session), q.Row, q.Question); err == nil {
			_, err = s.postBody(ctx, in.Task, body)
		}
	case qHandoffAnswer:
		_, err = s.RecordHandoff(ctx, HandoffInput{Task: in.Task, Session: q.Session, From: "owner role", To: "asking role"})
	case qNextAttempt:
		_, err = s.StartAttempt(ctx, in.Task, st.Base)
	default:
		err = fmt.Errorf("%s: the stage %d of the question is no act of the driver", in.Task, q.Stage)
	}
	return stop, "", err
}

// runOwner starts the answer session: the step answer-to-a-question under the attempt of the asking sessions, as the owner role of their
// questions, at tier reasoning, read-only, at the base of the last asking session.
func (s *Sessions) runOwner(ctx context.Context, in ShapeTaskInput, v loopView, st loopStep) (string, error) {
	const step = "answer-to-a-question"
	res, err := ResolveStep(step, "", st.Ask.Owner)
	if err != nil {
		return "", err
	}
	h, m, err := s.PairFor(ctx, res.Role, res.Tier, in.Harnesses, in.Models, nil)
	if err != nil {
		return pairStop(in.Task, step, err)
	}
	return "", s.startAs(ctx, in, v.Kind, step, res.Role, st.Attempt, st.Ask.Base, h, m, nil, true)
}

// answersOf gives the IDs of the answers that the attempt after attempt has to cite: those of the answer session of the owner role of the
// questions of attempt, when it ended done with the status completed, in the order of its answer rows. It gives none otherwise.
func (v loopView) answersOf(attempt int) []string {
	a := strconv.Itoa(attempt)
	asking, first := v.asking(a)
	if len(asking) == 0 {
		return nil
	}
	ans := v.answers(a, first)
	if len(ans) == 0 {
		return nil
	}
	owner, ok := answerOwner(v.Events[:v.startIndex(ans[0])], v.Results, a, Owners(DefaultOwners()))
	last := ans[len(ans)-1]
	if _, done := resultEnd(v.Events, last); !ok || !done || v.role(last) != owner || v.Status[last] != "completed" {
		return nil
	}
	var ids []string
	for _, r := range v.Results[last] {
		if r[0] == "answer" {
			ids = append(ids, last+"."+r[1])
		}
	}
	return ids
}
