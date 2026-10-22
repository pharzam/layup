package run

import (
	"maps"
	"slices"
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
		{"a plan review whose copy is no reject (task T-t94t: a reject with a reason is a handoff to the planner)", upTo("specification", 6).ev("handoff", "1", "S-r1", "", "", "invalid verdict-field"), posStop, "plan review → developer is invalid verdict-field"},
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
		{"a done with a refusal of the branch", upTo("specification", 5).end("S-r1", "1", "done", "completed").ev("refused", "1", "S-r1", "", "", "branch"), "was refused (branch)"},
		{"a done with a closed attempt", upTo("specification", 5).end("S-r1", "1", "done", "completed").ev("refused", "1", "S-r1", "", "", "closed-attempt"), "was refused (closed-attempt)"},
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

// retryPending is the session whose eligible end the next attempt would replace, "" for none.
func retryPending(b *viewBuilder) string { _, pending := retries(b.v); return pending }

// artifactEnd ends the session done, with the status given and the refusals given, each its own event refused.
func artifactEnd(b *viewBuilder, id, attempt, status string, reasons ...string) *viewBuilder {
	b.end(id, attempt, "done", status)
	for _, r := range reasons {
		b.ev("refused", attempt, id, "", "", r)
	}
	return b
}

// loop.md, An uncertain end gets one new attempt, the class artifact (task T-0sj8, #478): a session that ended done with the status completed
// and one refusal, artifact, is an eligible end of each of the four stages, as a crash is.
func TestARefusedArtifactIsAnEligibleEndOfEachStage(t *testing.T) {
	for _, st := range []struct {
		name string
		n    int
		id   string
		next position
	}{{"the plan", 2, "S-p1", posPlan}, {"the plan review", 5, "S-r1", posReview}, {"the developer", 8, "S-d1", posDeveloper}, {"the verifier", 13, "S-v1", posDeveloper}} {
		b := artifactEnd(upTo("specification", st.n), st.id, "1", "completed", "artifact")
		if got := b.at(); got.Pos != posRetry || got.Attempt != 1 || got.Base != shB || got.Session != st.id {
			t.Errorf("%s refused for an artifact: %+v; want the retry at attempt 1, base %s, of %s", st.name, got, shB[:7], st.id)
		}
		if got := b.ev("attempt", "2", "", shB, "", "").at(); got.Pos != st.next || got.Attempt != 2 {
			t.Errorf("%s refused for an artifact, attempt 2 recorded: %+v; want position %d in attempt 2", st.name, got, st.next)
		}
	}
}

// The demo's records (T-ut9g): the plan review of attempt 1 crashed, the one of attempt 2 at the same base ended done and refused artifact:
// the next run gives attempt 3 at the base, with the session start rows and the valid plan handoff that make both sessions the plan review.
func TestTheDemoRecordsOfTUt9gGiveAttempt3(t *testing.T) {
	b := upTo("specification", 4).start("S-38157561", roleQA, "1", shB).end("S-38157561", "1", "crash", "").ev("attempt", "2", "", shB, "", "").
		start("S-141cc16c", roleQA, "2", shB)
	artifactEnd(b, "S-141cc16c", "2", "completed", "artifact")
	if rep, pending := retries(b.v); !maps.Equal(rep, Replaced{"1": "S-38157561"}) || pending != "S-141cc16c" {
		t.Fatalf("retries = %v, %q; want attempt 1 replaced and S-141cc16c pending", rep, pending)
	}
	if got := b.at(); got.Pos != posRetry || got.Attempt != 2 || got.Base != shB || got.Session != "S-141cc16c" {
		t.Errorf("the demo's records: %+v; want the retry of S-141cc16c in attempt 2 at the base", got)
	}
	if got := b.ev("attempt", "3", "", shB, "", "").at(); got.Pos != posReview || got.Attempt != 3 || got.Base != shB {
		t.Errorf("the demo's records and attempt 3: %+v; want the plan review of attempt 3 at the base", got)
	}
}

// Each negative differs from the retry in one fact only; none is an eligible end, so none is pending and none is replaced.
func TestARefusalOfAnotherKindOrStatusIsNoEligibleEnd(t *testing.T) {
	asked := func(role string, reasons ...string) *viewBuilder {
		return artifactEnd(upTo("specification", 2).end("S-p1", "1", "done", "needs_context").start("S-a1", role, "1", shB), "S-a1", "1", "completed", reasons...)
	}
	closeOut := func() *viewBuilder {
		return artifactEnd(upTo("specification", 15).start("S-c1", stepRole("close-out"), "1", shH), "S-c1", "1", "completed", "artifact")
	}
	for _, c := range []struct {
		name string
		b    *viewBuilder
	}{
		{"closed-attempt", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "closed-attempt")},
		{"branch", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "branch")},
		{"push-lost", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "push-lost")},
		{"two refusal events, artifact first", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "artifact", "branch")},
		{"two refusal events, artifact last", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "closed-attempt", "artifact")},
		{"two refusal events of artifact", artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "artifact", "artifact")},
		{"a crash", upTo("specification", 5).end("S-r1", "1", "crash", "").ev("refused", "1", "S-r1", "", "", "artifact")},
		{"a no-result", upTo("specification", 5).end("S-r1", "1", "no-result", "").ev("refused", "1", "S-r1", "", "", "artifact")},
		{"a wall", upTo("specification", 5).end("S-r1", "1", "wall", "").ev("refused", "1", "S-r1", "", "", "artifact")},
		{"the status blocked", artifactEnd(upTo("specification", 5), "S-r1", "1", "blocked", "artifact")},
		{"the status with no kept result", artifactEnd(upTo("specification", 5), "S-r1", "1", "", "artifact")},
		{"the close-out session", closeOut()},
		{"an answer session of the plan role", asked(rolePlan, "artifact")},
		{"an answer session of the developer role", asked(roleSpec, "artifact")},
	} {
		if p := retryPending(c.b); p != "" {
			t.Errorf("%s: %q is pending; want no eligible end", c.name, p)
		}
		if got := replacedOf(c.b.ev("attempt", "2", "", shB, "", "").v); got != nil && len(got) != 0 {
			t.Errorf("%s: replaced %v; want none", c.name, got)
		}
	}
}

// Precedence (condition 2 of the plan review): the artifact class needs the status completed. A refused result of the status needs_context
// keeps the question path, with no retry; the answer sessions of its attempt stay no stage, refused or not.
func TestARefusedResultThatAsksAQuestionKeepsTheQuestionPath(t *testing.T) {
	b := artifactEnd(upTo("specification", 2), "S-p1", "1", "needs_context", "artifact")
	if got := b.at(); got.Pos != posQuestion {
		t.Errorf("a plan refused for an artifact that asks: %d %q; want the question", got.Pos, got.Msg)
	}
	if p := retryPending(b); p != "" {
		t.Errorf("a plan refused for an artifact that asks: %q is pending; want no retry", p)
	}
}

// The separate bound (the Operator, 6338065517): crash and no-result share one retry per stage and an artifact refusal has its own, in either
// order; a second end of the same class stops.
func TestAnArtifactRefusalHasItsOwnRetryPerStageInEitherOrder(t *testing.T) {
	end := func(b *viewBuilder, id, n, class string) *viewBuilder {
		if class == "artifact" {
			return artifactEnd(b, id, n, "completed", "artifact")
		}
		return b.end(id, n, class, "")
	}
	for _, c := range [][3]string{{"crash", "artifact", "crash"}, {"artifact", "crash", "artifact"}, {"no-result", "artifact", "no-result"}, {"artifact", "no-result", "artifact"}} {
		b := end(upTo("specification", 5), "S-r1", "1", c[0]).ev("attempt", "2", "", shB, "", "").start("S-r2", roleQA, "2", shB)
		end(b, "S-r2", "2", c[1])
		if got := b.at(); got.Pos != posRetry || got.Attempt != 2 || got.Session != "S-r2" {
			t.Fatalf("%s, then %s: %+v; want the retry of S-r2 in attempt 2", c[0], c[1], got)
		}
		b.ev("attempt", "3", "", shB, "", "")
		if got := b.at(); got.Pos != posReview || got.Attempt != 3 {
			t.Errorf("%s, then %s, attempt 3 recorded: %+v; want the plan review of attempt 3", c[0], c[1], got)
		}
		// The third end: the class of the first end is spent, and so is the class of the second.
		for _, third := range []string{c[0], c[1]} {
			x := &viewBuilder{loopView{Task: b.v.Task, Kind: b.v.Kind, Events: slices.Clone(b.v.Events), Sessions: slices.Clone(b.v.Sessions), Status: maps.Clone(b.v.Status), Records: b.v.Records}}
			x.start("S-r3", roleQA, "3", shB)
			end(x, "S-r3", "3", third)
			if got := x.at(); got.Pos != posStop {
				t.Errorf("%s, %s, then %s: %d %q; want the stop, each class has spent its retry", c[0], c[1], third, got.Pos, got.Msg)
			}
		}
	}
	// The same class twice stops at the second end.
	b := artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "artifact").ev("attempt", "2", "", shB, "", "").start("S-r2", roleQA, "2", shB)
	if got := artifactEnd(b, "S-r2", "2", "completed", "artifact").at(); got.Pos != posStop || !strings.Contains(got.Msg, "S-r2 was refused (artifact)") {
		t.Errorf("two artifact refusals of one stage: %d %q; want the stop", got.Pos, got.Msg)
	}
	// Another stage has its own artifact retry.
	d := artifactEnd(upTo("specification", 5), "S-r1", "1", "completed", "artifact").ev("attempt", "2", "", shB, "", "").start("S-r2", roleQA, "2", shB).end("S-r2", "2", "done", "completed").
		ev("handoff", "2", "S-r2", "", "", "valid comment 15-1").start("S-d2", roleSpec, "2", shB)
	if got := artifactEnd(d, "S-d2", "2", "completed", "artifact").at(); got.Pos != posRetry || got.Session != "S-d2" {
		t.Errorf("a developer refused for an artifact after a retried plan review: %+v; want the retry", got)
	}
}

// The verifier refused for an artifact: the retry is at the task base and runs a developer; it counts no review round, as a crash counts none.
func TestAnArtifactRefusalOfTheVerifierCountsNoReviewRound(t *testing.T) {
	r := artifactEnd(upTo("specification", 13), "S-v1", "1", "completed", "artifact")
	if got := r.at(); got.Pos != posRetry || got.Base != shB || got.Session != "S-v1" {
		t.Errorf("the verifier refused for an artifact: %+v; want the retry at the task base %s", got, shB[:7])
	}
	c := r.ev("attempt", "2", "", shB, "", "").devRound("2", shG).verdict("2", shG, "material", 1)
	if got := c.at(); got.Pos != posNextAttempt {
		t.Errorf("one material round after a verifier refused for an artifact: %d %q; want the next attempt", got.Pos, got.Msg)
	}
	if n := reviewRounds(c.v); n != 1 {
		t.Errorf("reviewRounds = %d; want 1", n)
	}
	d := c.ev("attempt", "3", "", shB, "", "").devRound("3", shG).verdict("3", shG, "material", 1)
	if got := d.at(); got.Pos != posWait || !strings.Contains(got.Msg, "the review-cycle cap 1 is reached") {
		t.Errorf("two material rounds around the retry: %d %q; want the wait at the cap", got.Pos, got.Msg)
	}
}

// Task T-t94t (#483): a plan review that rejects the plan sends the task back to its planner. The rule of the one admitted reject: a review no retry
// replaced, whose result is completed, whose kept copy is the file its result names by hash and rejects with a reason (the view's Rejects), and whose
// latest handoff is valid. Counting, reopening and the reasons all read it.

// reject adds the reasons of the kept copy of a plan review, as shapeView reads them.
func (b *viewBuilder) reject(id string, reasons ...string) *viewBuilder {
	if b.v.Rejects == nil {
		b.v.Rejects = map[string][]string{}
	}
	b.v.Rejects[id] = reasons
	return b
}

// planRound adds a plan and its handoff, then a plan review that ended and the handoff of the detail, in the attempt.
func (b *viewBuilder) planRound(attempt, plan, review string, planComment, reviewComment int, detail string) *viewBuilder {
	b.start(plan, rolePlan, attempt, shB).end(plan, attempt, "done", "completed").ev("handoff", attempt, plan, "", "", "valid comment "+strconv.Itoa(planComment)+"-1")
	return b.start(review, roleQA, attempt, shB).end(review, attempt, "done", "completed").ev("handoff", attempt, review, "", "", detail+" comment "+strconv.Itoa(reviewComment)+"-1")
}

func TestARejectedPlanStartsANewAttemptWithThePlanReopened(t *testing.T) {
	b := upTo("specification", 7).reject("S-r1", "no ordered steps")
	if got := b.at(); got.Pos != posReplan || got.Attempt != 1 || got.Base != shB {
		t.Fatalf("a valid reject: %d %q; want the replan of attempt 1 at the base", got.Pos, got.Msg)
	}
	b.ev("attempt", "2", "", shB, "", "")
	if got := b.at(); got.Pos != posPlan || got.Attempt != 2 {
		t.Fatalf("after the new attempt: %+v; want the plan of attempt 2, the plan and the review of attempt 1 passed over", got)
	}
	b.start("S-p2", rolePlan, "2", shB).end("S-p2", "2", "done", "completed")
	if got := b.at(); got.Pos != posPlanFinish || got.Session != "S-p2" {
		t.Fatalf("the new plan ended: %+v; want to finish S-p2", got)
	}
	b.ev("handoff", "2", "S-p2", "", "", "valid comment 14-1")
	if got := b.at(); got.Pos != posReview || got.Session != "" {
		t.Fatalf("the new plan handed off: %+v; want the plan review to run (not the review of attempt 1)", got)
	}
	b.start("S-r2", roleQA, "2", shB).end("S-r2", "2", "done", "completed").ev("handoff", "2", "S-r2", "", "", "valid comment 15-1")
	if got := b.at(); got.Pos != posDeveloper || got.Attempt != 2 {
		t.Errorf("the approving review of attempt 2: %+v; want the developer", got)
	}
	// The same facts with a review whose copy rejects with no reason, or an approve, are no reject.
	if got := upTo("specification", 7).at(); got.Pos != posDeveloper {
		t.Errorf("a valid handoff of a review that is no reject: %+v; want the developer", got)
	}
}

func TestTheRejectsAreBoundedByAnOwnCountOfTheCycleCap(t *testing.T) {
	b := upTo("specification", 7).reject("S-r1", "no ordered steps").ev("attempt", "2", "", shB, "", "").planRound("2", "S-p2", "S-r2", 14, 15, "valid").reject("S-r2", "still none")
	if n := planRejects(b.v); n != 2 {
		t.Errorf("planRejects = %d; want 2", n)
	}
	if got := b.at(); got.Pos != posStop || got.Msg != "T-sp01: the plan review rejected the plan 2 times" {
		t.Errorf("the second reject: %d %q; want the stop \"T-sp01: the plan review rejected the plan 2 times\"", got.Pos, got.Msg)
	}
	// A material round of the verifier counts none of them.
	m := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1)
	if n := planRejects(m.v); n != 0 {
		t.Errorf("planRejects with a material verifier = %d; want 0", n)
	}
	if got := m.at(); got.Pos != posNextAttempt {
		t.Errorf("a material round: %+v; want the next attempt (the plan count is its own)", got)
	}
}

func TestARejectAfterARetryReopensThePlan(t *testing.T) {
	b := upTo("specification", 4).start("S-r1", roleQA, "1", shB).end("S-r1", "1", "crash", "").ev("attempt", "2", "", shB, "", "").
		start("S-r2", roleQA, "2", shB).end("S-r2", "2", "done", "completed").ev("handoff", "2", "S-r2", "", "", "valid comment 15-1").reject("S-r2", "no steps")
	if got := b.at(); got.Pos != posReplan || got.Attempt != 2 {
		t.Fatalf("a reject in attempt 2: %+v; want the replan", got)
	}
	if rep := replacedOf(b.v); !maps.Equal(rep, Replaced{"1": "S-r1"}) {
		t.Errorf("replacedOf = %v; want attempt 1 replaced by the retry", rep)
	}
	if got := b.ev("attempt", "3", "", shB, "", "").at(); got.Pos != posPlan || got.Attempt != 3 {
		t.Errorf("attempt 3: %+v; want the plan", got)
	}
}

// demoRecords is the records of T-ut9g up to the handoff of S-0ef59a25: the plan review of attempt 1 crashed, the one of attempt 2 ended done and was
// refused for an artifact, the one of attempt 3 rejected and its handoff is the invalid verdict-field of the earlier binary. Both kept copies reject,
// with different reasons.
func demoRecords() *viewBuilder {
	b := newView("specification").ev("attempt", "1", "", shB, "", "").start("S-fb051a4c", rolePlan, "1", shB).end("S-fb051a4c", "1", "done", "completed").
		ev("handoff", "1", "S-fb051a4c", "", "", "valid comment 11-1").start("S-38157561", roleQA, "1", shB).end("S-38157561", "1", "crash", "").
		ev("attempt", "2", "", shB, "", "").start("S-141cc16c", roleQA, "2", shB)
	artifactEnd(b, "S-141cc16c", "2", "completed", "artifact").reject("S-141cc16c", "an older reason of attempt 2")
	return b.ev("attempt", "3", "", shB, "", "").start("S-0ef59a25", roleQA, "3", shB).end("S-0ef59a25", "3", "done", "completed").
		ev("handoff", "3", "S-0ef59a25", "", "", "invalid verdict-field comment 6339776504-1").reject("S-0ef59a25", "no ordered steps", "no plan for the Won't rows", "no size base", "no demo sentence")
}

func TestTheDemoRecordsOfTUt9gGiveAttempt4(t *testing.T) {
	b := demoRecords()
	if got := b.at(); got.Pos != posReviewFinish || got.Session != "S-0ef59a25" || got.Attempt != 3 {
		t.Fatalf("the demo's records: %+v; want the handoff of S-0ef59a25 recorded again (the invalid verdict-field of a reject with reasons)", got)
	}
	b.ev("handoff", "3", "S-0ef59a25", "", "", "valid comment 6339776504-1")
	if got := b.at(); got.Pos != posReplan || got.Attempt != 3 || got.Base != shB {
		t.Fatalf("the demo's records, handoff recorded: %+v; want the replan of attempt 3 at the base", got)
	}
	if n := planRejects(b.v); n != 1 {
		t.Errorf("planRejects = %d; want 1, S-141cc16c (refused, replaced) not counted", n)
	}
	if got := b.ev("attempt", "4", "", shB, "", "").at(); got.Pos != posPlan || got.Attempt != 4 {
		t.Fatalf("attempt 4: %+v; want the plan", got)
	}
	want := []Finding{{Rule: "plan-review:reject", Text: "no ordered steps"}, {Rule: "plan-review:reject", Text: "no plan for the Won't rows"}, {Rule: "plan-review:reject", Text: "no size base"}, {Rule: "plan-review:reject", Text: "no demo sentence"}}
	if got := rejectionOf(b.v, 4); !slices.Equal(got, want) {
		t.Errorf("rejectionOf(4) = %+v; want the four reasons of S-0ef59a25 and none of S-141cc16c", got)
	}
}

// Condition 2 of the plan review: a review that was refused or replaced is no admitted reject, whatever its copy says.
func TestARefusedOrReplacedReviewWhoseCopyRejectsReopensNothingAndSuppliesNothing(t *testing.T) {
	// S-141cc16c, refused for an artifact and not yet replaced (attempt 2 is the latest): the retry, not the replan.
	b := newView("specification").ev("attempt", "1", "", shB, "", "").start("S-fb051a4c", rolePlan, "1", shB).end("S-fb051a4c", "1", "done", "completed").
		ev("handoff", "1", "S-fb051a4c", "", "", "valid comment 11-1").start("S-141cc16c", roleQA, "1", shB)
	artifactEnd(b, "S-141cc16c", "1", "completed", "artifact").reject("S-141cc16c", "other reasons")
	if got := b.at(); got.Pos != posRetry || got.Session != "S-141cc16c" {
		t.Errorf("a refused review whose copy rejects: %+v; want the retry", got)
	}
	if n := planRejects(b.v); n != 0 {
		t.Errorf("planRejects = %d; want 0", n)
	}
	if got := b.ev("attempt", "2", "", shB, "", "").at(); got.Pos != posReview || got.Attempt != 2 {
		t.Errorf("after the retry's attempt: %+v; want the plan review (the plan is not reopened)", got)
	}
	if got := rejectionOf(b.v, 2); got != nil {
		t.Errorf("rejectionOf(2) = %+v; want none from a refused review", got)
	}
	// A done review with a copy that rejects and no handoff yet is finished first, not a reject.
	c := upTo("specification", 6).reject("S-r1", "x")
	if got := c.at(); got.Pos != posReviewFinish {
		t.Errorf("a review with no handoff: %+v; want to finish it", got)
	}
	if n := planRejects(c.v); n != 0 {
		t.Errorf("planRejects with no handoff = %d; want 0", n)
	}
}

// A bare reject is an invalid reject-reasons: the stop, and no attempt (condition 3).
func TestABareRejectIsAStopAndStartsNoAttempt(t *testing.T) {
	b := upTo("specification", 6).ev("handoff", "1", "S-r1", "", "", "invalid reject-reasons comment 12-1")
	if got := b.at(); got.Pos != posStop || !strings.Contains(got.Msg, "is invalid reject-reasons comment 12-1") {
		t.Errorf("a bare reject: %+v; want the stop with the invalid handoff", got)
	}
	if n := planRejects(b.v); n != 0 {
		t.Errorf("planRejects = %d; want 0", n)
	}
}

// qaKind classifies by the first valid plan handoff, which no later plan moves: a review stays a review and no retry is unwound (the red control).
func TestAPlanReviewKeepsItsStageWhenALaterPlanIsHandedOff(t *testing.T) {
	b := demoRecords()
	before := replacedOf(b.v)
	b.ev("handoff", "3", "S-0ef59a25", "", "", "valid comment 6339776504-1").ev("attempt", "4", "", shB, "", "").start("S-p4", rolePlan, "4", shB).end("S-p4", "4", "done", "completed").
		ev("handoff", "4", "S-p4", "", "", "valid comment 7000000001-1")
	for _, id := range []string{"S-38157561", "S-141cc16c", "S-0ef59a25"} {
		if got := b.v.stageOf(id); got != stageReview {
			t.Errorf("stageOf(%s) = %q after the plan of attempt 4; want %q", id, got, stageReview)
		}
	}
	if after := replacedOf(b.v); !maps.Equal(before, after) || len(after) != 2 {
		t.Errorf("replacedOf after the new plan = %v, before %v; want the same two replaced attempts", after, before)
	}
	if got := b.at(); got.Pos != posReview || got.Attempt != 4 {
		t.Errorf("attempt 4 with its plan handed off: %+v; want the plan review", got)
	}
}

// The reasons go to the planner of every attempt after the reject, across a question and a crashed plan's retry, and to no other step.
func TestThePlanOfTheAttemptAfterARejectIsGivenTheReasons(t *testing.T) {
	b := upTo("specification", 7).reject("S-r1", "no ordered steps", "no demo sentence")
	want := []Finding{{Rule: "plan-review:reject", Text: "no ordered steps"}, {Rule: "plan-review:reject", Text: "no demo sentence"}}
	if got := rejectionOf(b.v, 1); got != nil {
		t.Errorf("rejectionOf(1) = %+v; want none for the first plan", got)
	}
	b.ev("attempt", "2", "", shB, "", "").start("S-p2", rolePlan, "2", shB).end("S-p2", "2", "crash", "").ev("attempt", "3", "", shB, "", "")
	for _, attempt := range []int{2, 3} {
		if got := rejectionOf(b.v, attempt); !slices.Equal(got, want) {
			t.Errorf("rejectionOf(%d) = %+v; want the two reasons", attempt, got)
		}
	}
	if got := b.at(); got.Pos != posPlan || got.Attempt != 3 {
		t.Errorf("after the crashed plan's retry: %+v; want the plan of attempt 3", got)
	}
}

// Round 1 of T-t94t (F1): a plan review whose verdict cannot be established from its hash-checked kept copy is no approval. A valid handoff over a
// lost copy was read as an approval (the reject left no entry), so the developer ran on a rejected or an unreviewed plan.
func (b *viewBuilder) unread(id string) *viewBuilder {
	if b.v.Unread == nil {
		b.v.Unread = map[string]bool{}
	}
	b.v.Unread[id] = true
	return b
}

func TestReviewerUnreadReviewIsNoApprovalInThePosition(t *testing.T) {
	const stop = "the plan review S-r1 cannot be read from its kept copy"
	// Before the new plan: the review is the last of the stage. The control has its copy.
	before := func() *viewBuilder { return upTo("specification", 7) }
	// After the new plan of the replan attempt is handed off, and with the earlier reject lost, the old review would stand for the new plan's.
	after := func() *viewBuilder {
		return upTo("specification", 7).ev("attempt", "2", "", shB, "", "").start("S-p2", rolePlan, "2", shB).end("S-p2", "2", "done", "completed").
			ev("handoff", "2", "S-p2", "", "", "valid comment 14-1")
	}
	for _, c := range []struct {
		name  string
		build func() *viewBuilder
		ctl   position
	}{{"before the new plan", before, posReplan}, {"after the new plan", after, posReview}} {
		if got := c.build().reject("S-r1", "no steps").at(); got.Pos != c.ctl {
			t.Errorf("%s, the control with its copy: %d %q; want position %d", c.name, got.Pos, got.Msg, c.ctl)
		}
		for _, lost := range []string{"missing", "mismatched"} { // both leave the review unread: the view cannot tell them apart, the reader refuses both
			got := c.build().unread("S-r1").at()
			if got.Pos != posStop || !strings.Contains(got.Msg, stop) {
				t.Errorf("%s, %s kept copy: %d %q; want the stop %q and never the developer", c.name, lost, got.Pos, got.Msg, stop)
			}
		}
	}
}

// Round 1 of T-t94t (F2): a bare reject is checked on the pair plan review to task plan, and the stop names that pair.
func TestReviewerBareRejectNamesActualPair(t *testing.T) {
	b := upTo("specification", 6).ev("handoff", "1", "S-r1", "", "", "invalid reject-reasons comment 12-1")
	b.v.Verdicts = map[string]string{"S-r1": "reject"}
	if got := b.at(); got.Pos != posStop || !strings.Contains(got.Msg, "plan review → task plan is invalid reject-reasons comment 12-1") {
		t.Errorf("a bare reject: %d %q; want the stop naming plan review → task plan", got.Pos, got.Msg)
	}
	// The developer pair keeps its label for its own failures.
	c := upTo("specification", 6).ev("handoff", "1", "S-r1", "", "", "invalid verdict-field comment 12-1")
	c.v.Verdicts = map[string]string{"S-r1": "approve"}
	if got := c.at(); got.Pos != posStop || !strings.Contains(got.Msg, "plan review → developer is invalid verdict-field") {
		t.Errorf("an invalid approve: %d %q; want the developer pair", got.Pos, got.Msg)
	}
}
