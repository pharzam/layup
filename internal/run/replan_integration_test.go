//go:build integration

package run

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// A plan review that rejects the plan sends the task back to its planner (task T-t94t, #483; loop.md, The driver of a Shape task, "A rejected
// plan goes back to the planner"): the production path through state.sessions(), the httptest forge, real git and the scripted harnesses, whose
// knob plan-verdicts holds the verdict of the n-th plan review.

const rejectReason = "the plan gives no ordered steps"

// planPrompts gives the prompts of the plan sessions of the task, in order.
func (d *drvWorld) planPrompts(task string) []string {
	var out []string
	for _, p := range d.promptsOf(task) {
		if strings.HasPrefix(p, "Prompt file of the step plan-of-a-task ") {
			out = append(out, p)
		}
	}
	return out
}

// invalidHandoffs gives the session and the detail of each event handoff of the task that is not valid.
func (d *drvWorld) invalidHandoffs(task string) []string {
	var out []string
	for _, e := range d.eventsOf(task) {
		if e[1] == "handoff" && !strings.HasPrefix(e[6], "valid") {
			out = append(out, e[3]+" "+e[6])
		}
	}
	return out
}

func TestARejectedPlanGoesBackToThePlannerWithTheReasons(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("plan-verdicts", "reject\napprove\n")
	spec, _ := d.open()
	d.mustVerify(spec)
	if got := d.attempts(spec); len(got) != 2 || got[0] != got[1] {
		t.Errorf("the bases of the attempts %v; want two attempts at one base", got)
	}
	plans := d.byRole(spec, stepRole("plan-of-a-task"))
	if len(plans) != 2 {
		t.Fatalf("%d plan sessions; want 2", len(plans))
	}
	prompts := d.planPrompts(spec)
	if len(prompts) != 2 || strings.Contains(prompts[0], "plan-review:reject") || !strings.Contains(prompts[1], "findings of the last attempt:\n- plan-review:reject: "+rejectReason+"\n") {
		t.Errorf("the plan prompts %q; want the reason in the second only", prompts)
	}
	if bad := d.invalidHandoffs(spec); len(bad) != 0 {
		t.Errorf("invalid handoffs %v; want every handoff valid, the reject too", bad)
	}
	authors := map[string]bool{}
	harness := map[string]string{}
	for _, r := range d.sessionsOf(spec) {
		harness[r[sessionID]] = r[sessionHarness]
		if r[sessionRole] == stepRole("plan-of-a-task") {
			authors[r[sessionHarness]] = true
		}
	}
	qa := d.byRole(spec, stepRole("verification"))
	if len(qa) != 3 || authors[harness[qa[1]]] {
		t.Errorf("%d qa sessions, the new plan review on %s; want 3 and a harness outside the plan authors %v", len(qa), harness[qa[1]], authors)
	}
}

func TestASecondRejectStopsAtTheCapOfItsOwn(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("plan-verdicts", "reject\nreject\n")
	spec, _ := d.open()
	want := spec + ": the plan review rejected the plan 2 times"
	got, err := d.advance(spec, nil)
	if err != nil || got.Stop != want || got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop %q", got, err, want)
	}
	events, sessions := len(d.eventsOf(spec)), len(d.sessionsOf(spec))
	if again, err := d.advance(spec, nil); err != nil || again.Stop != want || len(d.eventsOf(spec)) != events || len(d.sessionsOf(spec)) != sessions {
		t.Errorf("the second call = %+v, %v; %d events, %d sessions; want the same stop and no event or session", again, err, len(d.eventsOf(spec))-events, len(d.sessionsOf(spec))-sessions)
	}
	if n := len(d.attempts(spec)); n != 2 {
		t.Errorf("%d attempts; want 2", n)
	}
}

// Condition 3 of the plan review: a reject with no reason is an invalid reject-reasons, a stop, and starts no attempt.
func TestABareRejectIsAnInvalidHandoffAndStartsNoAttempt(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("plan-verdicts", "bare-reject\n")
	spec, _ := d.open()
	got, err := d.advance(spec, nil)
	if err != nil || !strings.Contains(got.Stop, "plan review → task plan is invalid reject-reasons comment ") || got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop with an invalid reject-reasons", got, err)
	}
	if bad := d.invalidHandoffs(spec); len(bad) != 1 || !strings.Contains(bad[0], "invalid reject-reasons comment ") {
		t.Errorf("invalid handoffs %v; want the one of the review", bad)
	}
	if n := len(d.attempts(spec)); n != 1 {
		t.Errorf("%d attempts; want 1, no new attempt", n)
	}
}

func TestARejectAfterACrashedPlanReviewReopensThePlan(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("crash-plan-review", "1")
	d.ctlWrite("plan-verdicts", "approve\nreject\napprove\n")
	spec, _ := d.open()
	d.mustVerify(spec)
	if got := d.attempts(spec); len(got) != 3 || got[0] != got[1] || got[1] != got[2] {
		t.Errorf("the bases of the attempts %v; want three attempts at one base", got)
	}
	if n := len(d.byRole(spec, stepRole("plan-of-a-task"))); n != 2 {
		t.Errorf("%d plan sessions; want 2", n)
	}
}

func TestARestartBeforeOrAfterTheReplanAttemptStartsExactlyOne(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ name, fail string }{{"before the attempt", "layup run: attempt 2 of"}, {"after the attempt", ", attempt 2 of "}} {
		d := newDrvWorld(t, loopPerms)
		d.ctlWrite("plan-verdicts", "reject\napprove\n")
		spec, _ := d.open()
		_, err := d.advance(spec, func(s *Sessions) {
			s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
				if strings.Contains(m, c.fail) {
					return errors.New("the commit failed")
				}
				return nil
			}}
		})
		if err == nil {
			t.Fatalf("%s: AdvanceShapeTask did not fail", c.name)
		}
		d.restart()
		d.mustVerify(spec)
		if got := d.attempts(spec); len(got) != 2 {
			t.Errorf("%s: %d attempts; want exactly two", c.name, len(got))
		}
		if p := len(d.byRole(spec, stepRole("plan-of-a-task"))); p != 2 {
			t.Errorf("%s: %d plan sessions; want 2", c.name, p)
		}
	}
}

// keptFault makes the read of one kept file answer data (nil: the file is absent).
type keptFault struct {
	SessionStore
	path string
	data []byte
}

func (k *keptFault) ReadFile(ctx context.Context, p string) ([]byte, error) {
	if p == k.path {
		return k.data, nil
	}
	return k.SessionStore.ReadFile(ctx, p)
}

// The demo's records (T-ut9g): a target whose table of transitions holds the thirteen rows of its first task, a plan review that crashed, one that was
// refused for an artifact, one that rejected, and the handoff that the earlier binary recorded, plan review to developer, invalid verdict-field. The next
// run gives attempt 4 with a new plan session whose prompt holds the reasons, with no repair: the table keeps its bytes and the old event stays.
func TestTheDemoRecordsOfTUt9gGiveAttempt4WithNoRepair(t *testing.T) {
	t.Parallel()
	// The setup: a rejection whose handoff the new binary could not record (its commit fails), then the earlier binary's own handoff.
	setup := func(t *testing.T) (*drvWorld, string, string, []byte) {
		d := newDrvWorld(t, loopPerms)
		old := transitionBytes(t, original13())
		pushRecords(t, d.w.bare, map[string]string{"transitions.tsv": string(old)}, "the table of the first task")
		d.restart()
		d.ctlWrite("crash-plan-review", "1")
		d.ctlWrite("artifact-plan-review", "2")
		d.ctlWrite("plan-verdicts", "approve\napprove\nreject\napprove\n")
		spec, _ := d.open()
		got, err := d.advance(spec, func(s *Sessions) {
			s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
				if strings.Contains(m, "the handoff plan review to task plan") {
					return errors.New("the commit failed")
				}
				return nil
			}}
		})
		if err == nil && !strings.Contains(got.Stop, "invalid verdict-field") { // the earlier binary's stop is already the record
			t.Fatalf("AdvanceShapeTask = %+v, %v; want a failed commit of the handoff of the reject", got, err)
		}
		d.restart()
		qa := d.byRole(spec, stepRole("verification"))
		if len(qa) != 3 {
			t.Fatalf("%d qa sessions; want the crashed, the refused and the rejecting review", len(qa))
		}
		if len(d.handoffsOf(spec, qa[2])) == 0 {
			if _, err := d.r.sessions().RecordHandoff(context.Background(), HandoffInput{Task: spec, Session: qa[2], From: "plan review", To: "developer", File: PlanReviewFile, Kind: "plan-review"}); err != nil {
				t.Fatal(err)
			}
		}
		if got := d.handoffsOf(spec, qa[2]); len(got) != 1 || !strings.HasPrefix(got[0], "invalid verdict-field comment ") {
			t.Fatalf("the handoff of the earlier binary: %v; want invalid verdict-field", got)
		}
		return d, spec, qa[2], old
	}
	d, spec, review, old := setup(t)
	d.restart()
	d.mustVerify(spec)
	if got := d.attempts(spec); len(got) != 4 {
		t.Errorf("%d attempts; want 4: the crash, the refusal, the reject and the replan", len(got))
	}
	if got := d.handoffsOf(spec, review); len(got) != 2 || !strings.HasPrefix(got[0], "invalid verdict-field") || !strings.HasPrefix(got[1], "valid comment ") {
		t.Errorf("the handoffs of the review %v; want the old event kept and a valid one after it", got)
	}
	if prompts := d.planPrompts(spec); len(prompts) != 2 || !strings.Contains(prompts[1], "- plan-review:reject: "+rejectReason+"\n") {
		t.Errorf("the plan prompts %d; want a second one that holds the reason", len(prompts))
	}
	if got := keptRecord(t, d.w, "transitions.tsv"); got != string(old) {
		t.Error("transitions.tsv changed; want its bytes kept")
	}
	// Condition 5: kept bytes that are absent, or are not the bytes the result names by hash, trigger neither the recovery nor a replan.
	for _, c := range []struct {
		name string
		data []byte
	}{{"absent", nil}, {"not the named bytes", []byte("verdict: reject\n" + rejectReason + "\nand more\n")}} {
		d, spec, review, _ := setup(t)
		d.restart()
		got, err := d.advance(spec, func(s *Sessions) {
			s.Store = &keptFault{SessionStore: s.Store, path: KeptPath(spec, review, PlanReviewFile), data: c.data}
		})
		if err != nil || !strings.Contains(got.Stop, "is invalid verdict-field comment ") || got.Verified {
			t.Errorf("%s: AdvanceShapeTask = %+v, %v; want the stop with the invalid verdict-field", c.name, got, err)
		}
		if n := len(d.attempts(spec)); n != 3 {
			t.Errorf("%s: %d attempts; want 3, no replan", c.name, n)
		}
		if got := d.handoffsOf(spec, review); len(got) != 1 {
			t.Errorf("%s: the handoffs of the review %v; want only the old one", c.name, got)
		}
	}
}

// Round 1 (F1): the reject handoff is valid and its kept copy is then lost (absent, or other bytes than the result names). The driver stops with the
// review named, before and after the new plan is handed off, and runs no developer; the controls are the same runs with the copy intact.
func TestReviewerValidRejectLostCopyCannotAuthorizeDeveloper(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name, fail string
		nth        int // the commit of the matching message to fail, counted from 1
		attempts   int
	}{{"before the new plan", "layup run: attempt 2 of", 1, 1}, {"after the new plan", "the handoff task plan to plan review of", 2, 2}} {
		for _, lost := range []string{"intact", "missing", "mismatched"} {
			d := newDrvWorld(t, loopPerms)
			d.ctlWrite("plan-verdicts", "reject\napprove\n")
			spec, _ := d.open()
			seen := 0
			_, err := d.advance(spec, func(s *Sessions) {
				s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
					if strings.Contains(m, c.fail) {
						if seen++; seen == c.nth {
							return errors.New("the commit failed")
						}
					}
					return nil
				}}
			})
			if err == nil {
				t.Fatalf("%s: AdvanceShapeTask did not fail", c.name)
			}
			d.restart()
			review := d.byRole(spec, stepRole("verification"))[0]
			tweak := func(s *Sessions) {}
			switch lost {
			case "missing":
				tweak = func(s *Sessions) {
					s.Store = &keptFault{SessionStore: s.Store, path: KeptPath(spec, review, PlanReviewFile)}
				}
			case "mismatched":
				tweak = func(s *Sessions) {
					s.Store = &keptFault{SessionStore: s.Store, path: KeptPath(spec, review, PlanReviewFile), data: []byte("verdict: approve\n")}
				}
			}
			got, err := d.advance(spec, tweak)
			if lost == "intact" {
				if err != nil || !got.Verified {
					t.Errorf("%s, the control: %+v, %v; want the verified head", c.name, got, err)
				}
				continue
			}
			if err != nil || got.Verified || !strings.Contains(got.Stop, "the plan review "+review+" cannot be read from its kept copy") {
				t.Errorf("%s, %s copy: %+v, %v; want the stop that names the review", c.name, lost, got, err)
			}
			if n := len(d.byRole(spec, stepRole("specification"))); n != 0 {
				t.Errorf("%s, %s copy: %d developer sessions; want none", c.name, lost, n)
			}
			if n := len(d.attempts(spec)); n != c.attempts {
				t.Errorf("%s, %s copy: %d attempts; want %d", c.name, lost, n, c.attempts)
			}
		}
	}
}
