package run

// The question in the loop (docs/spec/loop.md, A question during the work, "The question in the loop"; task T-r7b9, #426, row 66z): the stage
// of the question is told from a view of the records alone, and the open material finding of Position from the kept record of the last review.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

const (
	qa1, qa2   = "S-0a1b2c3d", "S-1b2c3d4e" // the asking sessions
	qow, qow2  = "S-2c3d4e5f", "S-3d4e5f60" // the answer sessions
	roleOwner  = "software-architect"       // the owner role of the label interface-contract
	qArch      = "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n"
	qHuman     = "question\t1\tMay it?\t—\tinterface-contract yes docs/spec/forge.md\n"
	qSys       = "question\t2\tWhere?\t—\tarchitecture-boundary no docs/spec/forge.md\n"
	qArch2     = "question\t2\tAnd it?\t—\tinterface-contract no docs/spec/forge.md\n"
	answerRow1 = "answer\t1\tYes.\t—\tdocs/spec/forge.md\n"
	answerRow2 = "answer\t2\tThere.\t—\tdocs/spec/forge.md\n"
	answerRow3 = "answer\t3\tAlso.\t—\tdocs/spec/forge.md\n"
)

// qStart gives the start row of a session as the reader of sessions.tsv gives it.
func qStart(t *testing.T, id, role, attempt, base string) []string {
	t.Helper()
	data := "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\n" +
		fmt.Sprintf("%s\tT-sp01\t%s\t%s\tclaude\t2.1.295\tclaude-opus-5-5\t%s\t%s\t1001\t251\t200000\t10.0\t60\t—\t—\n", id, attempt, role, base, base)
	rows, err := records.ReadSessions([]byte(data))
	if err != nil || len(rows) != 1 {
		t.Fatalf("the start row of %s: %v", id, err)
	}
	return rows[0]
}

// qResult gives the rows of a result file (the status row is the first).
func qResult(t *testing.T, status, rows string) [][]string {
	t.Helper()
	got, err := records.ReadResult([]byte("kind\tn\tvalue\tsha256\treason\nstatus\t1\t" + status + "\t—\tr\n" + rows))
	if err != nil {
		t.Fatalf("the result: %v", err)
	}
	return got
}

// qBuilder builds the view of the question stages.
type qBuilder struct {
	*viewBuilder
	t *testing.T
}

func newQ(t *testing.T) *qBuilder {
	b := &qBuilder{newView("specification"), t}
	b.v.Results, b.v.Issue = map[string][][]string{}, 3
	b.ev("attempt", "1", "", shB, "", "")
	return b
}

// session adds a started session with its full start row.
func (b *qBuilder) session(id, role, attempt, base string) *qBuilder {
	b.v.Sessions = append(b.v.Sessions, qStart(b.t, id, role, attempt, base))
	b.ev("session", attempt, id, "", "", "")
	return b
}

// finish adds the result of a session that ended done with the status and the rows.
func (b *qBuilder) finish(id, attempt, status, rows string) *qBuilder {
	b.v.Results[id] = qResult(b.t, status, rows)
	b.end(id, attempt, "done", status)
	return b
}

// asks adds an asking session that ended needs_context with the question rows.
func (b *qBuilder) asks(id, role, rows string) *qBuilder {
	return b.session(id, role, "1", shB).finish(id, "1", "needs_context", rows)
}

func (b *qBuilder) handoff(id, detail string) *qBuilder {
	b.ev("handoff", "1", id, "", "", detail)
	return b
}

// copied adds the copy of the comment of a question row of the session or of an answer row (with its question ID).
func (b *qBuilder) copied(id, kind string, n int, questionID string) *qBuilder {
	var row []string
	for _, r := range b.v.Results[id] {
		if r[0] == kind && r[1] == strconv.Itoa(n) {
			row = r
		}
	}
	var start []string
	for _, s := range b.v.Sessions {
		if s[sessionID] == id {
			start = s
		}
	}
	body, err := QuestionComment(start, row)
	if kind == "answer" {
		body, err = AnswerComment(start, row, questionID)
	}
	if err != nil {
		b.t.Fatalf("the comment of %s %s %d: %v", id, kind, n, err)
	}
	c := make([]string, copySHA+2)
	c[copyComment], c[copySeen], c[copyIssue], c[copySHA], c[copyBody] = strconv.Itoa(700+len(b.v.Copies)), "1", "3", digest(body), "copies/x.md"
	b.v.Copies = append(b.v.Copies, c)
	return b
}

func (b *qBuilder) stage() qStep { return questionPosition(b.v, 1) }

// flow builds the records of one question of the plan session up to the n-th fact: 0 the asking session, 1 its handoff, 2 the copy of its
// question, 3 the answer session started, 4 its result with one answer row, 5 the copy of the answer, 6 its handoff.
func flow(t *testing.T, n int, rows string) *qBuilder {
	b := newQ(t).asks(qa1, rolePlan, rows)
	if n >= 1 {
		b.handoff(qa1, "valid")
	}
	if n >= 2 {
		b.copied(qa1, "question", 1, "")
	}
	if n >= 3 {
		b.session(qow, roleOwner, "1", shB)
	}
	if n >= 4 {
		b.finish(qow, "1", "completed", answerRow1)
	}
	if n >= 5 {
		b.copied(qow, "answer", 1, qa1+".1")
	}
	if n >= 6 {
		b.handoff(qow, "valid")
	}
	return b
}

// U1: the stage is the first one the records show not done, P1 to P8; each act is followed by a fresh read, and a stage never does an
// earlier one again. An attempt with no asking session has none (the green negative beside the red controls).
func TestTheQuestionPositionIsTheFirstStageTheRecordsShowNotDone(t *testing.T) {
	want := []qStep{
		{Stage: qHandoffAsk, Session: qa1},
		{Stage: qPostQuestion, Session: qa1},
		{Stage: qRunOwner, Owner: roleOwner, Base: shB},
		{Stage: qRunOwner, Owner: roleOwner, Base: shB},
		{Stage: qPostAnswer, Session: qow, Question: qa1 + ".1"},
		{Stage: qHandoffAnswer, Session: qow},
		{Stage: qNextAttempt},
	}
	for n, w := range want {
		got := flow(t, n, qArch).stage()
		if got.Stage != w.Stage || got.Session != w.Session || got.Owner != w.Owner || got.Base != w.Base || got.Question != w.Question {
			t.Errorf("after %d facts: stage %d, session %q, owner %q, base %.7q, question %q; want stage %d, %q, %q, %.7q, %q",
				n, got.Stage, got.Session, got.Owner, got.Base, got.Question, w.Stage, w.Session, w.Owner, w.Base, w.Question)
		}
	}
	if got := newQ(t).session(qa1, rolePlan, "1", shB).finish(qa1, "1", "completed", "").stage(); got.Stage != qNone {
		t.Errorf("an attempt with no asking session: stage %d; want none", got.Stage)
	}
	if got := flow(t, 1, qArch).stage(); len(got.Row) == 0 || got.Row[0] != "question" || got.Row[1] != "1" {
		t.Errorf("the row of the stage P2 = %v; want the question row", got.Row)
	}
}

// Review condition 3: a recorded handoff that is not valid is a stop on every restart (the presence of its event never lets the stage
// go on), for the question and for the answer.
func TestAnInvalidRecordedHandoffStaysAStop(t *testing.T) {
	for _, c := range []struct{ name, want string }{
		{"the question's", "T-sp01: the handoff asking role → owner role is invalid question-row"},
		{"the answer's", "T-sp01: the handoff owner role → asking role is invalid answer-row"},
	} {
		b := flow(t, 0, qArch).handoff(qa1, "invalid question-row")
		if c.name != "the question's" {
			b = flow(t, 5, qArch).handoff(qow, "invalid answer-row")
		}
		for i := 0; i < 2; i++ { // the same on a second read: nothing advances
			if got := b.stage(); got.Stage != qStop || got.Msg != c.want {
				t.Errorf("%s handoff, read %d: stage %d, %q; want the stop %q", c.name, i+1, got.Stage, got.Msg, c.want)
			}
		}
	}
}

// U2 (Q10, Q13): a question that goes to a human is a wait with the ID of the first question and no answer session; the comment of the
// question is posted before it. Control: the same question that needs no human has its answer session (U1).
func TestAHumanQuestionIsAWaitWithNoAnswerSession(t *testing.T) {
	const want = "waiting for the Operator: T-sp01: question " + qa1 + ".1"
	if got := flow(t, 1, qHuman).stage(); got.Stage != qPostQuestion {
		t.Errorf("a human question before its post: stage %d; want the post (P2)", got.Stage)
	}
	got := flow(t, 2, qHuman).stage()
	if got.Stage != qWait || got.Msg != want || got.Owner != "" {
		t.Errorf("a human question: stage %d, %q, owner %q; want the wait %q", got.Stage, got.Msg, got.Owner, want)
	}
	if got := flow(t, 2, qArch).stage(); got.Stage != qRunOwner {
		t.Errorf("the control, a question of an owner role: stage %d; want the answer session", got.Stage)
	}
}

// U3 (Q10): questions of two owner roles are a wait whose ID is the first question; two questions of one owner role are not (the control).
func TestQuestionsOfTwoOwnerRolesAreAWait(t *testing.T) {
	const want = "waiting for the Operator: T-sp01: question " + qa1 + ".1"
	two := func() *qBuilder {
		return flow(t, 1, qArch+qSys).copied(qa1, "question", 1, "").copied(qa1, "question", 2, "")
	}
	if got := two().stage(); got.Stage != qWait || got.Msg != want {
		t.Errorf("two owner roles: stage %d, %q; want the wait %q", got.Stage, got.Msg, want)
	}
	one := flow(t, 1, qArch+qArch2).copied(qa1, "question", 1, "").copied(qa1, "question", 2, "")
	if got := one.stage(); got.Stage != qRunOwner || got.Owner != roleOwner {
		t.Errorf("the control, two questions of one owner role: stage %d, owner %q; want the answer session", got.Stage, got.Owner)
	}
}

// qTwoAsking is a view of two asking sessions, both with a question of n 1, that started before either result, and an answer session
// with the answer rows.
func qTwoAsking(t *testing.T, answers string) *qBuilder {
	b := newQ(t)
	b.session(qa1, rolePlan, "1", shB).session(qa2, rolePlan, "1", shB)
	b.finish(qa1, "1", "needs_context", qArch).finish(qa2, "1", "needs_context", qArch)
	b.handoff(qa1, "valid").handoff(qa2, "valid").copied(qa1, "question", 1, "").copied(qa2, "question", 1, "")
	return b.session(qow, roleOwner, "1", shB).finish(qow, "1", "completed", answers)
}

// U4 (condition 4): the list of questions runs over the asking sessions and their rows, a session and the row's n each, so two sessions whose
// rows both have n 1 are two questions; the n-th answer row answers the n-th question; fewer answer rows and more are each a stop.
func TestTheNthAnswerNamesTheNthQuestionAndAMissingOrExtraAnswerStops(t *testing.T) {
	first := qTwoAsking(t, answerRow1+answerRow2)
	if got := first.stage(); got.Stage != qPostAnswer || got.Question != qa1+".1" || got.Session != qow || got.Row[1] != "1" {
		t.Errorf("two answers, none posted: stage %d, question %q, row %v; want the post of the answer 1 for %s.1", got.Stage, got.Question, got.Row, qa1)
	}
	second := qTwoAsking(t, answerRow1+answerRow2).copied(qow, "answer", 1, qa1+".1")
	if got := second.stage(); got.Stage != qPostAnswer || got.Question != qa2+".1" || got.Row[1] != "2" {
		t.Errorf("the first answer posted: stage %d, question %q; want the post of the answer 2 for %s.1 (two questions of n 1)", got.Stage, got.Question, qa2)
	}
	for _, c := range []struct{ name, answers, want string }{
		{"one answer for two questions", answerRow1, "T-sp01: the question " + qa2 + ".1 has no answer"},
		{"no answer row", "", "T-sp01: the question " + qa1 + ".1 has no answer"},
		{"three answers for two questions", answerRow1 + answerRow2 + answerRow3, "T-sp01: the answer " + qow + ".3 answers no question"},
	} {
		if got := qTwoAsking(t, c.answers).stage(); got.Stage != qStop || got.Msg != c.want {
			t.Errorf("%s: stage %d, %q; want the stop %q", c.name, got.Stage, got.Msg, c.want)
		}
	}
}

// U7 (Q9): an answer session that ends other than done with completed, asks a question back or is not of the owner role is a stop that names
// it, and its own question is no new asking session (the stop repeats; no handoff is asked for).
func TestAnOwnerSessionThatDoesNotAnswerStops(t *testing.T) {
	for _, c := range []struct {
		name string
		make func(b *qBuilder) *qBuilder
		want string
	}{
		{"failed", func(b *qBuilder) *qBuilder { b.session(qow, roleOwner, "1", shB).end(qow, "1", "wall", ""); return b }, "T-sp01: the session " + qow + " of answer-to-a-question ended wall"},
		{"status failed", func(b *qBuilder) *qBuilder { return b.session(qow, roleOwner, "1", shB).finish(qow, "1", "failed", "") }, "T-sp01: the session " + qow + " of answer-to-a-question ended failed"},
		{"asks again", func(b *qBuilder) *qBuilder {
			return b.session(qow, roleOwner, "1", shB).finish(qow, "1", "needs_context", qArch)
		}, "T-sp01: the session " + qow + " of answer-to-a-question ended needs_context"},
		{"a wrong role", func(b *qBuilder) *qBuilder {
			return b.session(qow, "software-developer", "1", shB).finish(qow, "1", "completed", answerRow1)
		}, "T-sp01: the session " + qow + " started after the question is of the role software-developer, not of the owner role " + roleOwner},
	} {
		b := c.make(flow(t, 2, qArch))
		if got := b.stage(); got.Stage != qStop || got.Msg != c.want {
			t.Errorf("%s: stage %d, %q; want the stop %q", c.name, got.Stage, got.Msg, c.want)
		}
	}
}

// loopPosition: a session of the latest attempt that ended needs_context is the position question, not the stop (the expected change of
// the case of #415); the other statuses keep the stop.
func TestAnAskingSessionIsThePositionQuestionAndOtherStatusesStayStops(t *testing.T) {
	if got := flow(t, 0, qArch).at(); got.Pos != posQuestion || got.Ask.Stage != qHandoffAsk {
		t.Errorf("a plan that asked: position %d (%q), stage %d; want the question at the handoff", got.Pos, got.Msg, got.Ask.Stage)
	}
	for _, status := range []string{"blocked", "failed", "decision_needed"} {
		b := upTo("specification", 2).end("S-p1", "1", "done", status)
		if got := b.at(); got.Pos != posStop || !strings.Contains(got.Msg, "ended "+status) {
			t.Errorf("a plan with the status %s: position %d, %q; want the stop", status, got.Pos, got.Msg)
		}
	}
}

// Review condition 3: the stage that asked is run again in the next attempt, and the plan review is still told from a verifier when an owner
// session of the plan role handed back (its handoff is no plan handoff).
func TestTheStageThatAskedRunsAgainInTheNextAttempt(t *testing.T) {
	again := func(b *viewBuilder) loopStep { return b.ev("attempt", "2", "", shB, "", "").at() }
	if got := again(upTo("specification", 2).end("S-p1", "1", "done", "needs_context")); got.Pos != posPlan {
		t.Errorf("a plan that asked, attempt 2: position %d (%q); want the plan again", got.Pos, got.Msg)
	}
	if got := again(upTo("specification", 5).end("S-r1", "1", "done", "needs_context")); got.Pos != posReview {
		t.Errorf("a plan review that asked, attempt 2: position %d (%q); want the plan review again", got.Pos, got.Msg)
	}
	if got := again(upTo("specification", 8).end("S-d1", "1", "done", "needs_context")); got.Pos != posDeveloper {
		t.Errorf("a developer that asked, attempt 2: position %d (%q); want the developer again", got.Pos, got.Msg)
	}
	// The developer asked, and an owner session of the plan role answered and handed back, in attempt 1: the plan review of attempt 1 stays the plan review.
	b := upTo("specification", 8).end("S-d1", "1", "done", "needs_context")
	b.start(qow, rolePlan, "1", shB).ev("handoff", "1", qow, "", "", "valid")
	if got := again(b); got.Pos != posDeveloper {
		t.Errorf("an owner handoff of the plan role: position %d (%q); want the developer, the plan review unchanged", got.Pos, got.Msg)
	}
}

// U6: the session that answers only reads, so it has no push.
func TestTheAnswerStepIsReadOnly(t *testing.T) {
	if !readOnlyStep("answer-to-a-question") {
		t.Error("readOnlyStep(answer-to-a-question) = false; want true")
	}
	if readOnlyStep("implement") {
		t.Error("readOnlyStep(implement) = true; want false (the control)")
	}
}

// The open material finding (condition 5; loop.md, the open unknowns): the record of the last review that ended, qualified.
type materialWorld struct {
	*handoffWorld
	n int
}

const reviewRecord = "head: " + aSHA + "\nverdict: %s\nfindings:\n%s"

// review adds a session of the role qa-engineer that ended done with a comment row of review-record.md and the text kept as its file; named is the
// hash that the result names (empty: the hash of the text); handoff is the detail of its handoff, "" for none.
func (m *materialWorld) review(t *testing.T, text, named, handoff string) string {
	t.Helper()
	m.n++
	id := fmt.Sprintf("S-4e5f60%02x", m.n)
	m.session(t, id, 1, "qa-engineer", "")
	if named == "" {
		named = commentSum(text)
	}
	m.store.files["tasks/T-ab12/results/"+id+".tsv"] = []byte(resultTop + "status\t1\tcompleted\t—\tr\ncomment\t1\tresult/review-record.md\t" + named + "\t—\n")
	if text != "" {
		m.store.files[KeptPath("T-ab12", id, "review-record.md")] = []byte(text)
	}
	if handoff != "" {
		m.events = append(m.events, []string{strconv.Itoa(len(m.events) + 1), "handoff", "1", id, "", "", handoff, "2026-10-09T11:00:00Z"})
		m.store.files["tasks/T-ab12/events.tsv"], _ = table(records.EventsSchema, m.events)
	}
	return id
}

func (m *materialWorld) open(t *testing.T) (bool, error) {
	t.Helper()
	events := m.eventsNow(t)
	var starts [][]string
	if data := m.store.files["sessions.tsv"]; data != nil {
		var err error
		if starts, err = records.ReadSessions(data); err != nil {
			t.Fatal(err)
		}
	}
	results := map[string][][]string{}
	for _, e := range events {
		if e[1] == "result" {
			rows, err := records.ReadResult(m.store.files["tasks/T-ab12/results/"+e[3]+".tsv"])
			if err != nil {
				t.Fatal(err)
			}
			results[e[3]] = rows
		}
	}
	return m.s.openMaterial(context.Background(), "T-ab12", events, starts, results)
}

func clean() string { return fmt.Sprintf(reviewRecord, "nothing material in scope", "") }
func material() string {
	return fmt.Sprintf(reviewRecord, "material", "- docs/spec/forge.md:3 spec: the section is thin\n")
}

// keptReadFault fails the read of one path.
type keptReadFault struct {
	SessionStore
	path string
}

func (r keptReadFault) ReadFile(ctx context.Context, p string) ([]byte, error) {
	if p == r.path {
		return nil, errors.New("the records did not answer")
	}
	return r.SessionStore.ReadFile(ctx, p)
}

// U5 (O11, condition 5): the last qualified review record decides. Nothing is open when no review ended; an unqualified record (no handoff yet,
// an invalid handoff, a missing or other kept file, another head, a form that does not parse) is open, and never lets an older clean record
// stand in its place; a plan review is never read as a review record; a failed read is an error.
func TestTheLastQualifiedReviewRecordIsTheSourceOfAnOpenMaterialFinding(t *testing.T) {
	const valid = "valid comment 801-1"
	for _, c := range []struct {
		name  string
		build func(t *testing.T, m *materialWorld)
		want  bool
	}{
		{"no review (green negative)", func(t *testing.T, m *materialWorld) {}, false},
		{"a clean record, valid handoff (green negative)", func(t *testing.T, m *materialWorld) { m.review(t, clean(), "", valid) }, false},
		{"a material record", func(t *testing.T, m *materialWorld) { m.review(t, material(), "", valid) }, true},
		{"a later clean record clears an earlier material one", func(t *testing.T, m *materialWorld) {
			m.review(t, material(), "", valid)
			m.review(t, clean(), "", valid)
		}, false},
		{"a later material record opens an earlier clean one (the last, not the first)", func(t *testing.T, m *materialWorld) {
			m.review(t, clean(), "", valid)
			m.review(t, material(), "", valid)
		}, true},
		{"a completed review awaiting its handoff", func(t *testing.T, m *materialWorld) { m.review(t, clean(), "", "") }, true},
		{"an invalid handoff", func(t *testing.T, m *materialWorld) { m.review(t, clean(), "", "invalid nothing-material") }, true},
		{"a missing kept file", func(t *testing.T, m *materialWorld) { m.review(t, "", commentSum(clean()), valid) }, true},
		{"a hash that is not the result's", func(t *testing.T, m *materialWorld) { m.review(t, clean(), commentSum("other"), valid) }, true},
		{"another head", func(t *testing.T, m *materialWorld) {
			m.review(t, strings.Replace(clean(), aSHA, strings.Repeat("e", 40), 1), "", valid)
		}, true},
		{"a form that does not parse", func(t *testing.T, m *materialWorld) { m.review(t, "verdict: nothing material in scope\n", "", valid) }, true},
		{"an unqualified later record never lets the older clean one stand", func(t *testing.T, m *materialWorld) {
			m.review(t, clean(), "", valid)
			m.review(t, clean(), commentSum("other"), valid)
		}, true},
		{"a plan review is no review record", func(t *testing.T, m *materialWorld) {
			id := m.review(t, material(), "", valid)
			m.store.files["tasks/T-ab12/results/"+id+".tsv"] = []byte(resultTop + "status\t1\tcompleted\t—\tr\ncomment\t1\tresult/plan-review.md\t" + commentSum(material()) + "\t—\n")
		}, false},
	} {
		m := &materialWorld{handoffWorld: newHandoffWorld(t)}
		c.build(t, m)
		if got, err := m.open(t); err != nil || got != c.want {
			t.Errorf("%s: open = %v, %v; want %v", c.name, got, err, c.want)
		}
	}
	m := &materialWorld{handoffWorld: newHandoffWorld(t)}
	id := m.review(t, clean(), "", valid)
	m.s.Store = keptReadFault{m.s.Store, KeptPath("T-ab12", id, "review-record.md")}
	if got, err := m.open(t); err == nil {
		t.Errorf("a failed read of the kept file: open = %v, %v; want an error", got, err)
	}
}

// Review condition 3: the answer session and its prompt are bound to the head the asking session was bound to: the base of the start row of the
// last asking session (a verifier's is its bound head), not the base of the attempt.
func TestTheAnswerSessionStartsAtTheBaseOfTheLastAskingSession(t *testing.T) {
	b := newQ(t)
	b.session(qa1, rolePlan, "1", shB).session(qa2, roleQA, "1", shH)
	b.finish(qa1, "1", "needs_context", qArch).finish(qa2, "1", "needs_context", qArch)
	b.handoff(qa1, "valid").handoff(qa2, "valid").copied(qa1, "question", 1, "").copied(qa2, "question", 1, "")
	if got := b.stage(); got.Stage != qRunOwner || got.Base != shH {
		t.Errorf("stage %d at base %.7q; want the answer session at the base %.7q of the last asking session", got.Stage, got.Base, shH)
	}
}
