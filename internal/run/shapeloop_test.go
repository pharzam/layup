package run

import (
	"maps"
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
		// Expected change (row 66z, #426): a session that ended needs_context is the position question, not a stop; the other statuses stay stops
		// (TestAnAskingSessionIsThePositionQuestionAndOtherStatusesStayStops).
		{"a plan whose status is needs_context", upTo("specification", 2).end("S-p1", "1", "done", "needs_context"), posQuestion, ""},
		{"a plan whose status is blocked", upTo("specification", 2).end("S-p1", "1", "done", "blocked"), posStop, "ended blocked"},
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

// U8: the judgement of the four statuses of a bound head (task T-tnxn, #427): an error or a not-active verifier stops before any retry, whatever else
// failed; a failure of layup/spec or layup/rules is a check to retry with, and one of layup/gates a gate failure.
func TestAFailedSpecOrRulesStatusIsACheckToRetryAndAnErrorStopsFirst(t *testing.T) {
	post := func(context, state, d string) StatusPost { return StatusPost{Status: Status{context, state, d}} }
	ok := []StatusPost{post(GatesContext, "success", "pass"), post(SpecContext, "success", "pass"), post(RulesContext, "success", "no rule path"), post(VerifyContext, "pending", "waiting")}
	with := func(i int, p StatusPost) []StatusPost {
		out := append([]StatusPost(nil), ok...)
		out[i] = p
		return out
	}
	for _, c := range []struct {
		name          string
		posts         []StatusPost
		stop          string
		gates, checks bool
	}{
		{"all pass", ok, "", false, false},
		{"layup/spec failed", with(1, post(SpecContext, "failure", "criteria: x")), "", false, true},
		{"layup/rules failed", with(2, post(RulesContext, "failure", "rule path: a")), "", false, true},
		{"layup/gates failed", with(0, post(GatesContext, "failure", "exits")), "", true, false},
		{"gates and spec failed", append(with(0, post(GatesContext, "failure", "exits"))[:1:1], with(1, post(SpecContext, "failure", "criteria: x"))[1:]...), "", true, true},
		{"layup/gates error beside a spec failure", append(with(0, post(GatesContext, "error", "no manifest"))[:1:1], with(1, post(SpecContext, "failure", "criteria: x"))[1:]...), "T-sp01: layup/gates is error: no manifest", false, false},
		{"layup/rules error after a spec failure", append(with(1, post(SpecContext, "failure", "criteria: x"))[:2:2], with(2, post(RulesContext, "error", "register"))[2:]...), "T-sp01: layup/rules is error: register", false, false},
		{"not-active after a rules failure", append(with(2, post(RulesContext, "failure", "rule path: a"))[:3:3], post(VerifyContext, "failure", "not-active: no harness")), "T-sp01: not-active: no harness", false, false},
	} {
		stop, gates, checks := judgeStatuses("T-sp01", c.posts)
		if stop != c.stop || gates != c.gates || checks != c.checks {
			t.Errorf("%s: judgeStatuses = %q, gates %v, checks %v; want %q, %v, %v", c.name, stop, gates, checks, c.stop, c.gates, c.checks)
		}
	}
}

// A devRound adds the developer's session of the attempt n, pushed and bound at head; a verdict adds the verifier's session of the attempt n at base.
func (b *viewBuilder) devRound(n, head string) *viewBuilder {
	d := "S-d" + n
	return b.start(d, roleSpec, n, shB).end(d, n, "done", "completed").ev("push", n, d, "", head, "task/T-sp01/"+n).ev("bound", n, d, "", head, "task/T-sp01/"+n).ev("handoff", n, d, "", "", "valid")
}

func (b *viewBuilder) verdict(n, base, verdict string, findings int) *viewBuilder {
	v := "S-v" + n
	return b.start(v, roleQA, n, base).end(v, n, "done", "completed").ev("handoff", n, v, "", "", "valid comment 20-1").copyOf(v, verdict, findings)
}

// loop.md, An uncertain end gets one new attempt (task T-qjjd, #475): the last session of the latest attempt, of one of the four
// stages, ended crash or no-result and was not refused: the position is posRetry at the base of the attempt; with the new attempt
// recorded, the stage runs again in it (the verifier's retry runs a developer, the verifier's own base being the developer's head).
func TestAnEligibleEndOfAStageIsTheRetryPosition(t *testing.T) {
	for _, st := range []struct {
		name string
		n    int
		id   string
		next position
	}{{"the plan", 2, "S-p1", posPlan}, {"the plan review", 5, "S-r1", posReview}, {"the developer", 8, "S-d1", posDeveloper}, {"the verifier", 13, "S-v1", posDeveloper}} {
		for _, class := range []string{"crash", "no-result"} {
			b := upTo("specification", st.n).end(st.id, "1", class, "")
			if got := b.at(); got.Pos != posRetry || got.Attempt != 1 || got.Base != shB || got.Session != st.id {
				t.Errorf("%s ended %s: %+v; want the retry at attempt 1, base %s, of %s", st.name, class, got, shB[:7], st.id)
			}
			if got := b.ev("attempt", "2", "", shB, "", "").at(); got.Pos != st.next || got.Attempt != 2 {
				t.Errorf("%s ended %s, attempt 2 recorded: %+v; want position %d in attempt 2", st.name, class, got, st.next)
			}
		}
	}
}

// The replacement is told from the eligible end that caused it: each negative differs from the retry in one fact only.
func TestAnAttemptIsReplacedOnlyByTheRetryOfItsEligibleEnd(t *testing.T) {
	crashed := func() *viewBuilder { return upTo("specification", 5) }
	again := func(b *viewBuilder) *viewBuilder { return b.ev("attempt", "2", "", shB, "", "") }
	closeOut := func() *viewBuilder {
		return upTo("specification", 15).start("S-c1", stepRole("close-out"), "1", shH).end("S-c1", "1", "crash", "")
	}
	asked := func(role string) *viewBuilder {
		return upTo("specification", 2).end("S-p1", "1", "done", "needs_context").start("S-a1", role, "1", shB).end("S-a1", "1", "crash", "")
	}
	for _, c := range []struct {
		name string
		b    *viewBuilder
		want Replaced
	}{
		{"a crash and its retry", again(crashed().end("S-r1", "1", "crash", "")), Replaced{"1": "S-r1"}},
		{"a no-result and its retry", again(crashed().end("S-r1", "1", "no-result", "")), Replaced{"1": "S-r1"}},
		{"a refused session", again(crashed().ev("refused", "1", "S-r1", "", "", "artifact").end("S-r1", "1", "crash", "")), nil},
		{"a closed attempt", again(crashed().end("S-r1", "1", "crash", "").ev("closed", "1", "", "", "", "")), nil},
		{"a rebased attempt", again(crashed().end("S-r1", "1", "crash", "").ev("rebased", "1", "", "", "", "")), nil},
		{"a lost push", again(crashed().end("S-r1", "1", "crash", "").ev("refused", "1", "S-d0", "", "", "push-lost")), nil},
		{"a successor at another base", crashed().end("S-r1", "1", "crash", "").ev("attempt", "2", "", shG, "", ""), nil},
		{"a successor before the end", again(crashed()).end("S-r1", "1", "crash", ""), nil},
		{"an end that is not the last session", again(crashed().end("S-r1", "1", "crash", "").start("S-r2", roleQA, "1", shB)), nil},
		{"the close-out session", again(closeOut()), nil},
		{"a crashed answer session of the plan role", again(asked(rolePlan)), nil},
		{"a crashed answer session of the developer role", again(asked(roleSpec)), nil},
		{"a wall", again(crashed().end("S-r1", "1", "wall", "")), nil},
		{"an output", again(crashed().end("S-r1", "1", "output", "")), nil},
		{"a start", again(crashed().end("S-r1", "1", "start", "")), nil},
		{"a gate failure", again(upTo("specification", 12)), nil},
		{"a material round", again(upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1)), nil},
		{"the handback of an answer", again(upTo("specification", 2).end("S-p1", "1", "done", "needs_context")), nil},
	} {
		if got := replacedOf(c.b.v); !maps.Equal(got, c.want) {
			t.Errorf("%s: replaced %v; want %v", c.name, got, c.want)
		}
	}
}

// A guard (green before and after): the ends that are not eligible stay stops, and a refusal is never skipped.
func TestEndsThatAreNotEligibleStayStops(t *testing.T) {
	for _, c := range []struct {
		name string
		b    *viewBuilder
		msg  string
	}{
		{"a crash with a refusal", upTo("specification", 5).end("S-r1", "1", "crash", "").ev("refused", "1", "S-r1", "", "", "closed-attempt"), "was refused (closed-attempt)"},
		{"a done with a refusal", upTo("specification", 5).end("S-r1", "1", "done", "completed").ev("refused", "1", "S-r1", "", "", "artifact"), "was refused (artifact)"},
		{"a wall", upTo("specification", 5).end("S-r1", "1", "wall", ""), "S-r1 of the plan review ended wall"},
		{"an output", upTo("specification", 5).end("S-r1", "1", "output", ""), "S-r1 of the plan review ended output"},
		{"a start", upTo("specification", 5).end("S-r1", "1", "start", ""), "S-r1 of the plan review ended start"},
	} {
		if got := c.b.at(); got.Pos != posStop || !strings.Contains(got.Msg, c.msg) {
			t.Errorf("%s: %d %q; want a stop holding %q", c.name, got.Pos, got.Msg, c.msg)
		}
	}
}

// The bound is once per stage, whatever the classes of the two ends; the end that spends a stage's retry is the eligible one that
// caused the replacement, so a verifier's crash spends the verifier's retry and not the developer's.
func TestTheBoundIsOncePerStageAcrossClasses(t *testing.T) {
	for _, c := range [][2]string{{"crash", "no-result"}, {"no-result", "crash"}, {"crash", "crash"}} {
		b := upTo("specification", 5).end("S-r1", "1", c[0], "").ev("attempt", "2", "", shB, "", "").start("S-r2", roleQA, "2", shB).end("S-r2", "2", c[1], "")
		if got := b.at(); got.Pos != posStop || got.Msg != "T-sp01: the session S-r2 of the plan review ended "+c[1] {
			t.Errorf("a plan review that ended %s, then %s: %d %q; want the stop \"ended %s\"", c[0], c[1], got.Pos, got.Msg, c[1])
		}
	}
	// Another stage has its own retry.
	b := upTo("specification", 5).end("S-r1", "1", "crash", "").ev("attempt", "2", "", shB, "", "").start("S-r2", roleQA, "2", shB).end("S-r2", "2", "done", "completed").
		ev("handoff", "2", "S-r2", "", "", "valid comment 15-1").start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "crash", "")
	if got := b.at(); got.Pos != posRetry || got.Attempt != 2 {
		t.Errorf("a developer crash after a retried plan review: %+v; want the retry in attempt 2", got)
	}
	// The reverse stage: attempt 1 holds a developer that ended done, which spends nothing.
	v := upTo("specification", 13).end("S-v1", "1", "crash", "").ev("attempt", "2", "", shB, "", "").start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "crash", "")
	if got := v.at(); got.Pos != posRetry || got.Attempt != 2 || got.Base != shB {
		t.Fatalf("a verifier crash, its retry, then the developer's first crash: %+v; want the retry in attempt 2", got)
	}
	if got := v.ev("attempt", "3", "", shB, "", "").start("S-d3", roleSpec, "3", shB).end("S-d3", "3", "crash", "").at(); got.Pos != posStop || !strings.Contains(got.Msg, "S-d3 of specification ended crash") {
		t.Errorf("the developer's second crash: %d %q; want the stop", got.Pos, got.Msg)
	}
}

// A retry counts no review round: a crash does not spend the cap, and a material round after one is still a round.
func TestTheReviewCycleCapCountsNoRetry(t *testing.T) {
	v1 := func() *viewBuilder {
		return upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1)
	}
	b := v1().ev("attempt", "2", "", shB, "", "").start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "crash", "").ev("attempt", "3", "", shB, "", "").devRound("3", shG).verdict("3", shG, "material", 1)
	if got := b.at(); got.Pos != posWait || !strings.Contains(got.Msg, "the review-cycle cap 1 is reached") {
		t.Errorf("two material rounds around a developer retry: %d %q; want the wait at the cap", got.Pos, got.Msg)
	}
	c := upTo("specification", 13).end("S-v1", "1", "crash", "").ev("attempt", "2", "", shB, "", "").devRound("2", shG).verdict("2", shG, "material", 1)
	if got := c.at(); got.Pos != posNextAttempt {
		t.Errorf("one material round after a verifier retry: %d %q; want the next attempt", got.Pos, got.Msg)
	}
}

// Two consecutive replacements of different stages: the findings of the round before them reach the developer of the attempt after them.
func TestTheFindingsPassOverTwoReplacedAttempts(t *testing.T) {
	b := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1).
		ev("attempt", "2", "", shB, "", "").start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "crash", "").
		ev("attempt", "3", "", shB, "", "").devRound("3", shG).start("S-v3", roleQA, "3", shG).end("S-v3", "3", "crash", "").ev("attempt", "4", "", shB, "", "")
	if found, head := findingsOf(b.v, 4); len(found) != 1 || head != "" {
		t.Errorf("findingsOf for attempt 4 = %+v, %q; want the finding of S-v1 and no gate head", found, head)
	}
	if got := priorAttempt(replacedOf(b.v), 4); got != 1 {
		t.Errorf("priorAttempt(4) = %d; want 1, passing over the attempts 3 and 2", got)
	}
}
