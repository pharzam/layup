//go:build integration

package run

// The question in the loop on the production path (task T-r7b9, #426, row 66z; loop.md, A question during the work, "The question in the
// loop"): the Sessions of state.sessions(), the httptest forge, real git and a scripted harness that asks and answers, with no bet.

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

const (
	askArch  = "question\t1\tMay the section name the forge?\t—\tinterface-contract no docs/spec/alpha.md\n"
	askHuman = "question\t1\tMay the section name the forge?\t—\tinterface-contract yes docs/spec/alpha.md\n"
	askSys   = "question\t2\tWhere is the boundary?\t—\tarchitecture-boundary no docs/spec/alpha.md\n"
	answer1  = "answer\t1\tYes, it may.\t—\tdocs/spec/alpha.md\n"
)

// askHarness is the scripted harness of the Shape world that also asks and answers: a file ask-<step> of ctl/ makes the next session of that step
// end needs_context with its rows (once), the step answer-to-a-question writes the rows of ctl/answers, and a session whose prompt lists the
// IDs of answers cites them in its status row, as the instruction tells it to.
func askHarness(ctl string) string {
	return `#!/bin/sh
` + harnessPrelude + `ctl=` + ctl + `
case "$line" in "Prompt file"*) keep ;; esac
keep() { :; }
ids=$(sed -n 's/.*Their IDs are \(.*\)\. Write each of these.*/\1/p' ../prompt.md)
res_head="kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone${ids:+ $ids}\n"
if [ -e $ctl/ask-$step ]; then
  printf 'kind\tn\tvalue\tsha256\treason\nstatus\t1\tneeds_context\t—\tasks\n' > ../result/result.tsv
  cat $ctl/ask-$step >> ../result/result.tsv; rm $ctl/ask-$step; exit 0
fi
if [ "$step" = answer-to-a-question ]; then
  { printf "$res_head"; cat $ctl/answers; } > ../result/result.tsv; exit 0
fi
` + harnessSteps(ctl)
}

// newAskWorld is the world of the Shape tasks with the owner role software-architect routed and the harness that asks.
func newAskWorld(t *testing.T) *drvWorld {
	t.Helper()
	d := newDrvWorld(t, loopPerms)
	pushRecords(t, d.w.bare, map[string]string{"routing.tsv": keptRecord(t, d.w, "routing.tsv") + "software-architect\treasoning\t1\tclaude\tm1\n"}, "the owner role")
	d.r = restarted(t, d.w)
	if err := os.WriteFile(d.harnesses[0][3], []byte(askHarness(d.ctl)), 0o755); err != nil {
		t.Fatal(err)
	}
	d.ctlWrite("answers", answer1)
	return d
}

// postedWith gives the bodies of the comments the run posted that hold the text.
func postedWith(d *drvWorld, sub string) []string {
	d.w.fake.mu.Lock()
	defer d.w.fake.mu.Unlock()
	var out []string
	for _, p := range d.w.fake.posted {
		if strings.Contains(p, sub) {
			out = append(out, p)
		}
	}
	return out
}

// byRole gives the IDs of the sessions of the task of the role.
func (d *drvWorld) byRole(task, role string) []string {
	var out []string
	for _, r := range d.sessionsOf(task) {
		if r[sessionRole] == role {
			out = append(out, r[sessionID])
		}
	}
	return out
}

// handoffsOf gives the details of the events handoff of a session.
func (d *drvWorld) handoffsOf(task, session string) []string {
	var out []string
	for _, e := range d.eventsOf(task) {
		if e[1] == "handoff" && e[3] == session {
			out = append(out, e[6])
		}
	}
	return out
}

// askedAndAnswered asks the plan question of the specification task, runs the driver to the verified head and gives the world, the task, the
// asking session and the answer session.
func askedAndAnswered(t *testing.T) (d *drvWorld, task, asker, owner string) {
	d = newAskWorld(t)
	task, _ = d.open()
	d.ctlWrite("ask-plan-of-a-task", askArch)
	if got, err := d.advance(task, nil); err != nil || !got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the verified head after the answer", got, err)
	}
	asker = d.byRole(task, stepRole("plan-of-a-task"))[0]
	owner = d.byRole(task, "software-architect")[0]
	return d, task, asker, owner
}

// I1, the demo (condition 2 of #367, row 66z): a plan that ends needs_context during Shape, with no bet, has its question posted and copied on
// the task's issue, its handoff recorded, the answer session of the owner role run read-only at the base of the asking session, its answer
// posted, the handoff back recorded and the next attempt started, whose plan cites the answer and is accepted by it.
func TestANeedsContextEndIsAnsweredByTheOwnerRoleAndTheNextAttemptStarts(t *testing.T) {
	d, task, asker, owner := askedAndAnswered(t)
	if _, err := os.Stat(filepath.Join(d.r.target, "bets.tsv")); err == nil {
		t.Error("a bets.tsv exists; the question must be reachable with no bet")
	}
	wantQ := asker + ": software-engineer session of " + task + ", attempt 1: question " + asker + ".1\n\nMay the section name the forge?"
	if got := postedWith(d, ": question "); len(got) != 1 || got[0] != wantQ {
		t.Errorf("the question comments %q; want the one %q", got, wantQ)
	}
	wantA := owner + ": software-architect session of " + task + ", attempt 1: answer " + owner + ".1 of " + asker + ".1\n\nYes, it may."
	if got := postedWith(d, ": answer "); len(got) != 1 || got[0] != wantA {
		t.Errorf("the answer comments %q; want the one %q", got, wantA)
	}
	copies := keptRecord(t, d.w, "copies.tsv")
	for _, body := range []string{wantQ, wantA} {
		if !strings.Contains(copies, "\t"+digest([]byte(body))+"\t") {
			t.Errorf("copies.tsv has no row of the comment %q", body[:40])
		}
	}
	if got := d.handoffsOf(task, asker); !slices.Equal(got, []string{"valid"}) {
		t.Errorf("the handoffs of the asking session %v; want one, valid", got)
	}
	if got := d.handoffsOf(task, owner); !slices.Equal(got, []string{"valid"}) {
		t.Errorf("the handoffs of the answer session %v; want one, valid", got)
	}
	kinds := d.kinds(task)
	if !slices.Equal(kinds[:8], []string{"attempt", "session", "result", "handoff", "session", "result", "handoff", "attempt"}) {
		t.Errorf("the first events %v; want the question, the answer session and attempt 2", kinds[:8])
	}
	for _, r := range d.sessionsOf(task) {
		if r[sessionID] == owner && (r[sessionAttempt] != "1" || r[sessionBase] != d.sessionBase(task, asker)) {
			t.Errorf("the start row of the answer session %v; want attempt 1 at the base of the asking session", r)
		}
	}
	for _, e := range d.eventsOf(task) {
		if e[3] == owner && e[1] == "push" {
			t.Errorf("the answer session pushed (%v); it only reads", e)
		}
	}
	// The next attempt cites the answer and the answer is accepted by it.
	events, results, roles := d.ledgerOf(task)
	if by, ok := AnswerAccepted(events, results, roles, Owners(DefaultOwners()), owner+".1"); !ok || roles[by] != stepRole("plan-of-a-task") {
		t.Errorf("AnswerAccepted(%s.1) = %q, %v; want the plan session of attempt 2", owner, by, ok)
	}
}

// ledgerOf gives what AnswerAccepted reads: the events, the typed rows of each result and the roles.
func (d *drvWorld) ledgerOf(task string) ([][]string, map[string][][]string, Roles) {
	d.t.Helper()
	events, results := d.eventsOf(task), map[string][][]string{}
	for _, e := range events {
		if e[1] == "result" {
			rows, err := records.ReadResult([]byte(keptRecord(d.t, d.w, "tasks/"+task+"/results/"+e[3]+".tsv")))
			if err != nil {
				d.t.Fatal(err)
			}
			results[e[3]] = rows
		}
	}
	return events, results, sessionRoles(d.sessionsOf(task))
}

func (d *drvWorld) sessionBase(task, session string) string {
	for _, r := range d.sessionsOf(task) {
		if r[sessionID] == session {
			return r[sessionBase]
		}
	}
	return ""
}

// I6 (Q6, Q7): the prompt of the answer session holds the question IDs and the result of the asking session, and the prompts of attempt 2 hold the
// ID of the answer; the prompts of attempt 1 (the plan that asked) hold no answer block.
func TestTheAnswerPromptHoldsTheQuestionsAndTheNextPromptsTheAnswerIDs(t *testing.T) {
	d, task, asker, owner := askedAndAnswered(t)
	prompts := d.promptsOf(task) // the plan that asked, the answer session, then the plan, the plan review, the developer, the verifier of attempt 2
	if len(prompts) != 6 {
		t.Fatalf("%d prompts; want 6", len(prompts))
	}
	if !strings.Contains(prompts[1], "step answer-to-a-question") || !strings.Contains(prompts[1], "asked: "+asker+".1.") || !strings.Contains(prompts[1], "/input/"+asker+".tsv ") {
		t.Errorf("the prompt of the answer session lacks the question ID or the result of the asking session:\n%s", prompts[1])
	}
	if strings.Contains(prompts[0], "Their IDs are") || strings.Contains(prompts[1], "Their IDs are") {
		t.Error("a prompt of attempt 1 holds the answers block")
	}
	for i, p := range prompts[2:] {
		if !strings.Contains(p, "Their IDs are "+owner+".1.") {
			t.Errorf("the prompt %d of attempt 2 lacks the ID of the answer %s.1", i+3, owner)
		}
	}
}

// I7: the result that a harness writes by the answer instruction crosses the real session boundary and its handoff owner role to asking role is
// valid, by the item answer-row that code computes from the result.
func TestAHarnessThatFollowsTheAnswerInstructionGivesAValidAnswerHandoff(t *testing.T) {
	d, task, _, owner := askedAndAnswered(t)
	kept := keptRecord(t, d.w, "tasks/"+task+"/results/"+owner+".tsv")
	rows, err := records.ReadResult([]byte(kept))
	if err != nil || !AnswerRowOK(rows) || statusOf(rows) != "completed" {
		t.Fatalf("the committed result of the answer session (%v):\n%s", err, kept)
	}
	if got := d.handoffsOf(task, owner); !slices.Equal(got, []string{"valid"}) {
		t.Errorf("the handoff of the answer session %v; want valid", got)
	}
}

// I2 and I3 (Q10, Q13): a question that goes to a human, and questions of more than one owner role, each end in the wait with the ID of the first
// question and start no answer session; the comments of the questions are posted once and a second call waits again, writing nothing.
func TestAHumanRoutedQuestionAndAMixedOwnerQuestionWaitForTheOperator(t *testing.T) {
	for _, c := range []struct {
		name, rows string
		comments   int
	}{{"a human-routed question", askHuman, 1}, {"questions of two owner roles", askArch + askSys, 2}} {
		d := newAskWorld(t)
		task, _ := d.open()
		d.ctlWrite("ask-plan-of-a-task", c.rows)
		asker := ""
		want := func() string { return "waiting for the Operator: " + task + ": question " + asker + ".1" }
		got, err := d.advance(task, nil)
		asker = d.byRole(task, stepRole("plan-of-a-task"))[0]
		if err != nil || got.Wait != want() || got.Verified || got.Stop != "" {
			t.Fatalf("%s: AdvanceShapeTask = %+v, %v; want the wait %q", c.name, got, err, want())
		}
		if n := len(d.byRole(task, "software-architect")) + len(d.byRole(task, "systems-architect")); n != 0 {
			t.Errorf("%s: %d answer sessions; want none", c.name, n)
		}
		if n := len(postedWith(d, ": question ")); n != c.comments {
			t.Errorf("%s: %d question comments; want %d", c.name, n, c.comments)
		}
		if h := d.handoffsOf(task, asker); !slices.Equal(h, []string{"valid"}) {
			t.Errorf("%s: the handoffs of the asking session %v; want valid", c.name, h)
		}
		events, posts := len(d.eventsOf(task)), len(postedWith(d, ""))
		again, err := d.advance(task, nil)
		if err != nil || again.Wait != want() || len(d.eventsOf(task)) != events || len(postedWith(d, "")) != posts {
			t.Errorf("%s: a second call = %+v, %v with %d new events and %d new comments; want the same wait and nothing written", c.name, again, err,
				len(d.eventsOf(task))-events, len(postedWith(d, ""))-posts)
		}
	}
}

// I4 (Q3, Q9, condition 3): an interruption after each act of the question, then a restart, writes each comment, handoff, answer session and attempt
// once; a session that started and has no result is run again (the stale start row stays).
func TestARestartAtEachStageOfTheQuestionWritesEachThingOnce(t *testing.T) {
	for _, c := range []struct {
		name   string
		fault  func(d *drvWorld) func(*Sessions)
		owners int // the sessions of the owner role after the restart
	}{
		{"the post of the question", func(d *drvWorld) func(*Sessions) { d.w.fake.failNote = ": question "; return nil }, 1},
		{"the copy of the question (adoption)", func(d *drvWorld) func(*Sessions) { return failCommit("records: the copy of the comment", 1) }, 1},
		{"the handoff to the owner role", func(d *drvWorld) func(*Sessions) { return failCommit("the handoff asking role to owner role", 1) }, 1},
		{"the end of the answer session", func(d *drvWorld) func(*Sessions) { return failCommit("layup run: the end of", 2) }, 2},
		{"the post of the answer", func(d *drvWorld) func(*Sessions) { d.w.fake.failNote = ": answer "; return nil }, 1},
		{"the copy of the answer (adoption)", func(d *drvWorld) func(*Sessions) { return failCommit("records: the copy of the comment", 2) }, 1},
		{"the handoff back to the asking role", func(d *drvWorld) func(*Sessions) { return failCommit("the handoff owner role to asking role", 1) }, 1},
		{"the start of the next attempt", func(d *drvWorld) func(*Sessions) { return failCommit("layup run: attempt 2 of", 1) }, 1},
	} {
		d := newAskWorld(t)
		task, _ := d.open()
		d.ctlWrite("ask-plan-of-a-task", askArch)
		if _, err := d.advance(task, c.fault(d)); err == nil {
			t.Errorf("%s: the first call = nil; want the interruption", c.name)
			continue
		}
		d.r = restarted(t, d.w)
		d.mustVerify(task)
		if n := len(postedWith(d, ": question ")); n != 1 {
			t.Errorf("%s: %d question comments; want 1", c.name, n)
		}
		if n := len(postedWith(d, ": answer ")); n != 1 {
			t.Errorf("%s: %d answer comments; want 1", c.name, n)
		}
		if n := len(d.byRole(task, "software-architect")); n != c.owners {
			t.Errorf("%s: %d sessions of the owner role; want %d", c.name, n, c.owners)
		}
		if n := strings.Count(strings.Join(d.kinds(task), " "), "attempt"); n != 2 {
			t.Errorf("%s: %d attempts; want 2", c.name, n)
		}
		for _, s := range d.sessionsOf(task) {
			if h := d.handoffsOf(task, s[sessionID]); (s[sessionRole] == stepRole("plan-of-a-task") || s[sessionRole] == "software-architect") && len(h) > 1 {
				t.Errorf("%s: the session %s has %d handoffs %v; want at most one each", c.name, s[sessionID], len(h), h)
			}
		}
	}
}

// failCommit fails the n-th records commit whose message holds the text, once.
func failCommit(text string, n int) func(*Sessions) {
	return func(s *Sessions) {
		seen := 0
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, text) {
				if seen++; seen == n {
					return errors.New("the run stopped")
				}
			}
			return nil
		}}
	}
}

// I5 (O11, condition 5): an open material finding of the last qualified review record keeps a task whose list is frozen uphill, through the
// Sessions of the production constructor; a later record of nothing material clears it, and the tier follows the position.
func TestAnOpenMaterialFindingKeepsTheTaskUphill(t *testing.T) {
	l := newLoop(t, nil)
	ctx := context.Background()
	s := l.prod()
	position := func() Position {
		t.Helper()
		p, err := s.Position(ctx, "T-ab12")
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	large := Plan{Size: "large"}
	l.frozenAt("claude")
	if p := position(); p != Downhill || TierOf(large, p) != "execution" {
		t.Fatalf("a frozen list and no review: %s, %s; want downhill, execution (the control)", p, TierOf(large, p))
	}
	l.record(1, "head: "+l.a.base+"\nverdict: material\nfindings:\n- docs/spec/forge.md:3 spec: the section is thin\n", "developer")
	if p := position(); p != Uphill || TierOf(large, p) != "reasoning" {
		t.Errorf("an open material finding: %s, %s; want uphill, reasoning", p, TierOf(large, p))
	}
	l.record(2, "head: "+l.a.base+"\nverdict: nothing material in scope\nfindings:\n", "close-out")
	if p := position(); p != Downhill || TierOf(large, p) != "execution" {
		t.Errorf("a later record of nothing material: %s, %s; want downhill, execution", p, TierOf(large, p))
	}
}

// record runs a verifier session that writes the text as its review record, posts it and records its handoff to the role.
func (l *loopRun) record(attempt int, text, to string) string {
	l.t.Helper()
	ctx := context.Background()
	if attempt > 1 {
		l.attempt(attempt)
	}
	spec := harness(l.t, l.a, planWork("review-record.md", text, "comment\t1\tresult/review-record.md\t"+commentSum(text)+"\t—"))
	spec.Role, spec.Attempt, spec.Pair.Harness = stepRole("verification"), attempt, "devin"
	id, err := l.r.TaskSession(ctx, spec)
	if err != nil {
		l.t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if _, err := l.r.PostComment(ctx, "T-ab12", id, "review-record.md"); err != nil {
		l.t.Fatalf("PostComment: %v", err)
	}
	if v, err := l.r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: id, From: "verifier", To: to, File: "review-record.md", Kind: "review-record"}); err != nil || !v.Valid {
		l.t.Fatalf("the handoff of the review record: %+v, %v", v, err)
	}
	return id
}
