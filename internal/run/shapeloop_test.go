package run

import (
	"strconv"
	"strings"
	"testing"
)

// The position of a Shape task is told from the records alone (loop.md, The driver of a Shape task):
// a view of synthetic rows, one more fact at each step.
type viewBuilder struct{ v loopView }

var (
	shB = strings.Repeat("a", 40) // the base of the attempts
	shH = strings.Repeat("b", 40) // the bound head of the developer
	shG = strings.Repeat("c", 40) // the bound head of the next attempt
)

func newView(kind string) *viewBuilder {
	return &viewBuilder{loopView{Task: "T-sp01", Kind: kind, Status: map[string]string{}, Records: map[string]ReviewRecord{}}}
}

func (b *viewBuilder) ev(kind string, attempt string, id, base, sha, detail string) *viewBuilder {
	n := strconv.Itoa(len(b.v.Events) + 1)
	b.v.Events = append(b.v.Events, []string{n, kind, attempt, id, base, sha, detail, "2026-10-19T00:00:00Z"})
	return b
}

// start adds the start row and the event session of a session of the role.
func (b *viewBuilder) start(id, role, attempt, base string) *viewBuilder {
	b.v.Sessions = append(b.v.Sessions, []string{id, "T-sp01", attempt, role, "claude", "1.0", "m1", base})
	return b.ev("session", attempt, id, "", "", "")
}

// end adds the event result of a session and the status of its kept result.
func (b *viewBuilder) end(id, attempt, class, status string) *viewBuilder {
	b.ev("result", attempt, id, "", "", class)
	b.v.Status[id] = status
	return b
}

func (b *viewBuilder) at() loopStep { return loopPosition(b.v) }

const (
	rolePlan = "software-engineer"
	roleQA   = "qa-engineer"
	roleSpec = "product-owner"
)

func (b *viewBuilder) copyOf(id string, verdict string, findings int) *viewBuilder {
	r := ReviewRecord{Verdict: verdict, Head: shH}
	for i := 0; i < findings; i++ {
		r.Findings = append(r.Findings, Finding{Rule: "r", File: "docs/spec/a.md", Line: 3, Text: "t"})
	}
	b.v.Records[id] = r
	return b
}

// upTo builds the records of the first attempt up to the step: the numbers are those of the table below.
func upTo(kind string, n int) *viewBuilder {
	b := newView(kind)
	steps := []func(){
		func() { b.ev("attempt", "1", "", shB, "", "") },
		func() { b.start("S-p1", rolePlan, "1", shB) },
		func() { b.end("S-p1", "1", "done", "completed") },
		func() { b.ev("handoff", "1", "S-p1", "", "", "valid comment 11-1") },
		func() { b.start("S-r1", roleQA, "1", shB) },
		func() { b.end("S-r1", "1", "done", "completed") },
		func() { b.ev("handoff", "1", "S-r1", "", "", "valid comment 12-1") },
		func() { b.start("S-d1", stepRole(devStep(kind)), "1", shB) },
		func() { b.end("S-d1", "1", "done", "completed") },
		func() { b.ev("push", "1", "S-d1", "", shH, "task/T-sp01/1") },
		func() { b.ev("bound", "1", "S-d1", "", shH, "task/T-sp01/1") },
		func() { b.ev("handoff", "1", "S-d1", "", "", "valid") },
		func() { b.start("S-v1", roleQA, "1", shH) },
		func() { b.end("S-v1", "1", "done", "completed") },
		func() {
			b.ev("handoff", "1", "S-v1", "", "", "valid comment 13-1")
			b.copyOf("S-v1", "nothing material in scope", 0)
		},
	}
	for _, f := range steps[:n] {
		f()
	}
	return b
}

func TestThePositionIsTheFirstStepTheRecordsShowNotDone(t *testing.T) {
	for _, kind := range []string{"specification", "architecture"} {
		want := []position{posAttempt, posPlan, posPlan, posPlanFinish, posReview, posReview, posReviewFinish, posDeveloper, posDeveloper,
			posStop, posStop, posToVerifier, posShip, posShip, posVerifierFinish, posVerified}
		for n, w := range want {
			if got := upTo(kind, n).at(); got.Pos != w {
				t.Errorf("%s, %d facts: position %d (%q); want %d", kind, n, got.Pos, got.Msg, w)
			}
		}
	}
	// What an act needs comes with the position.
	if got := upTo("specification", 11).at(); got.Attempt != 1 || got.Base != shB || got.Head != shH || got.Session != "S-d1" {
		t.Errorf("the position to the verifier = %+v; want attempt 1, base, head and the developer's session", got)
	}
	if got := upTo("specification", 3).at(); got.Session != "S-p1" {
		t.Errorf("the position to finish the plan = %+v; want the plan session", got)
	}
}

// Condition 3 of the plan review: a stage whose session ended is finished from its kept files, not run again; a stopped run
// is run again; a failure of the session or an invalid handoff is a stop that names it.
func TestAStageThatEndedIsFinishedAndAStageThatStoppedIsRunAgain(t *testing.T) {
	for _, c := range []struct {
		name string
		b    *viewBuilder
		want position
		msg  string
	}{
		{"the plan session with no result", upTo("specification", 2), posPlan, ""},
		{"the verifier ended and posted nothing yet", upTo("specification", 14), posVerifierFinish, ""},
		{"a plan session that ended wall", upTo("specification", 2).end("S-p1", "1", "wall", ""), posStop, "ended wall"},
		{"a plan whose status is needs_context", upTo("specification", 2).end("S-p1", "1", "done", "needs_context"), posStop, "ended needs_context"},
		{"a refused session", upTo("specification", 2).ev("refused", "1", "S-p1", "", "", "artifact"), posStop, "was refused (artifact)"},
		{"an invalid plan handoff", upTo("specification", 3).ev("handoff", "1", "S-p1", "", "", "invalid plan-fields"), posStop, "task plan → plan review is invalid plan-fields"},
		{"a plan review that rejects", upTo("specification", 6).ev("handoff", "1", "S-r1", "", "", "invalid verdict-field"), posStop, "plan review → developer is invalid verdict-field"},
		{"an invalid developer handoff", upTo("specification", 11).ev("handoff", "1", "S-d1", "", "", "invalid spec-sections"), posStop, "is invalid spec-sections"},
		{"a developer that ended with no commit pushed", upTo("specification", 9), posStop, "ended with no push"},
		{"a push with no result", upTo("specification", 10), posStop, "a push with no result"},
		{"a refused push", upTo("specification", 9).ev("push", "1", "S-d1", "", shH, "task/T-sp01/1").ev("refused", "1", "S-d1", "", "", "push-refused"), posStop, "was refused (push-refused)"},
		{"a start that was refused", upTo("specification", 1).ev("refused", "1", "S-x1", "", "", "context 1175 1000"), posStop, "the start of the session S-x1 was refused (context 1175 1000)"},
		{"a start refused in an attempt before", upTo("specification", 1).ev("refused", "1", "S-x1", "", "", "version").ev("attempt", "2", "", shB, "", ""), posPlan, ""},
		{"an invalid verifier handoff", upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "invalid nothing-material"), posStop, "the handoff verifier is invalid nothing-material"},
	} {
		got := c.b.at()
		if got.Pos != c.want || !strings.Contains(got.Msg, c.msg) {
			t.Errorf("%s: %d %q; want %d holding %q", c.name, got.Pos, got.Msg, c.want, c.msg)
		}
	}
}

// loop.md, The driver of a Shape task: the plan review is the qa-engineer session between the last valid plan handoff and the
// first push; the verifier is the one that starts at the bound head of its attempt's developer; any other is a stop.
func TestThePlanReviewIsBeforeTheFirstPushAndTheVerifierAfterItsBound(t *testing.T) {
	// A verifier-like session before the push is the plan review: the position finishes it.
	if got := upTo("specification", 5).end("S-r1", "1", "done", "completed").at(); got.Pos != posReviewFinish {
		t.Errorf("a qa session before the first push: %+v; want the plan review to finish", got)
	}
	// A qa session at another base than the bound head, after the push, is neither.
	b := upTo("specification", 12).start("S-x1", roleQA, "1", shB)
	if got := b.at(); got.Pos != posStop || !strings.Contains(got.Msg, "S-x1 is neither the plan review nor a verifier") {
		t.Errorf("an unclassified qa session: %+v; want the stop", got)
	}
	// A qa session before any plan handoff is neither.
	b = newView("specification").ev("attempt", "1", "", shB, "", "").start("S-x1", roleQA, "1", shB)
	if got := b.at(); got.Pos != posStop {
		t.Errorf("a qa session before the plan: %+v; want the stop", got)
	}
}

// loop.md, The cycle cap of a Shape task is 1: one material round starts the next attempt, the second waits; a
// record that is not mergeable always waits; neither starts an attempt of its own.
func TestOneMaterialRoundStartsAnAttemptAndTheSecondWaitsAtTheCap(t *testing.T) {
	material := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1)
	if got := material.at(); got.Pos != posNextAttempt || got.Session != "S-v1" || got.Base != shB {
		t.Errorf("the first material round: %+v; want the next attempt at the base, after the verdict of S-v1", got)
	}
	// The next attempt, its developer, and a second material round.
	b := material.ev("attempt", "2", "", shB, "", "")
	if got := b.at(); got.Pos != posDeveloper || got.Attempt != 2 {
		t.Fatalf("after the new attempt: %+v; want the developer of attempt 2", got)
	}
	b.start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "done", "completed").ev("push", "2", "S-d2", "", shG, "task/T-sp01/2").ev("bound", "2", "S-d2", "", shG, "task/T-sp01/2").
		ev("handoff", "2", "S-d2", "", "", "valid").start("S-v2", roleQA, "2", shG).end("S-v2", "2", "done", "completed").
		ev("handoff", "2", "S-v2", "", "", "valid comment 14-1").copyOf("S-v2", "material", 2)
	if got := b.at(); got.Pos != posWait || !strings.Contains(got.Msg, "waiting for the Operator: T-sp01: the review-cycle cap 1 is reached with findings open") {
		t.Errorf("the second material round: %+v; want the wait at the cap", got)
	}
	// A gate failure counts no round: a second attempt that began without a verifier round is still the first round.
	c := upTo("specification", 12).ev("attempt", "2", "", shB, "", "").start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "done", "completed").
		ev("push", "2", "S-d2", "", shG, "task/T-sp01/2").ev("bound", "2", "S-d2", "", shG, "task/T-sp01/2").ev("handoff", "2", "S-d2", "", "", "valid").
		start("S-v2", roleQA, "2", shG).end("S-v2", "2", "done", "completed").ev("handoff", "2", "S-v2", "", "", "valid comment 14-1").copyOf("S-v2", "material", 1)
	if got := c.at(); got.Pos != posNextAttempt {
		t.Errorf("a first round after an attempt of a gate failure: %+v; want the next attempt", got)
	}
}

func TestNotMergeableFindingsRecordedIsAWaitNeverAnAttempt(t *testing.T) {
	b := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "not mergeable, findings recorded", 1)
	if got := b.at(); got.Pos != posWait || !strings.Contains(got.Msg, "waiting for the Operator: T-sp01: not mergeable, findings recorded") {
		t.Errorf("not mergeable: %+v; want the wait", got)
	}
}

// The verified head needs nothing material at the bound head of the latest attempt: a record of an earlier attempt is not it.
func TestAVerifiedHeadNeedsNothingMaterialAtTheLastBoundHead(t *testing.T) {
	if got := upTo("specification", 15).at(); got.Pos != posVerified {
		t.Fatalf("the control: %+v; want the verified head", got)
	}
	b := upTo("specification", 15).ev("attempt", "2", "", shB, "", "")
	if got := b.at(); got.Pos != posDeveloper {
		t.Errorf("a new attempt after a verified head of the first: %+v; want its developer, not the verified head", got)
	}
}

func TestTheFindingsOfTheNextAttemptAreTheLastRoundsOrTheGatesTable(t *testing.T) {
	b := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1).ev("attempt", "2", "", shB, "", "")
	found, head := findingsOf(b.v, 2)
	if len(found) != 1 || found[0] != (Finding{Rule: "r", File: "docs/spec/a.md", Line: 3, Text: "t"}) || head != "" {
		t.Errorf("findingsOf after a material round = %+v, %q; want its finding and no gate head", found, head)
	}
	// An attempt that began after a head that the review never judged: the gate table of that head.
	g := upTo("specification", 12).ev("attempt", "2", "", shB, "", "")
	found, head = findingsOf(g.v, 2)
	if len(found) != 0 || head != shH {
		t.Errorf("findingsOf after a gate failure = %+v, %q; want no record and the bound head %s", found, head, shH)
	}
	// An older record is not used for a later attempt that had none of its own.
	o := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1).ev("attempt", "2", "", shB, "", "").
		start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "done", "completed").ev("push", "2", "S-d2", "", shG, "task/T-sp01/2").ev("bound", "2", "S-d2", "", shG, "task/T-sp01/2").
		ev("attempt", "3", "", shB, "", "")
	if found, head = findingsOf(o.v, 3); len(found) != 0 || head != shG {
		t.Errorf("findingsOf for attempt 3 = %+v, %q; want no finding of the older record and the gate head %s", found, head, shG)
	}
}

// loop.md, The driver of a Shape task: only the developer's session is ever bound; a bound of a session of another role stops the step,
// and the same ledger without it is the verified head.
func TestABoundOfASessionOutsideTheDeveloperRoleStopsTheStep(t *testing.T) {
	if got := upTo("specification", 15).at(); got.Pos != posVerified {
		t.Fatalf("the control: %+v; want the verified head", got)
	}
	for _, c := range []struct{ name, session string }{{"the verifier", "S-v1"}, {"the plan", "S-p1"}, {"the plan review", "S-r1"}} {
		b := upTo("specification", 15).ev("push", "1", c.session, "", shH, "task/T-sp01/1").ev("bound", "1", c.session, "", shH, "task/T-sp01/1")
		if got := b.at(); got.Pos != posStop || !strings.Contains(got.Msg, c.session+" has a bound and is not of the developer role") {
			t.Errorf("a bound of %s: %+v; want the stop that names %s", c.name, got, c.session)
		}
	}
}
