//go:build integration

package run

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/tsv"
)

// The specification task that a changed class reopens on the production path (task T-4naq, #428, row 66ab; loop.md, The specification task that
// a changed class reopens): the world of the first bet whose answer re-classes four facts (1 need to constraint, 2 measure to need, Q-002 context to
// measure, Q-003 context to need), state.sessions(), the httptest forge, real git. A restart is a fresh state on the same records. The tests of the
// stops replace the package variable betStopAt and so stay serial (main_integration_test.go).

const reclassAnswer = "REQ-001: Must, phase 1\nfact 1: constraint\nfact 2: need\nfact Q-002: measure\nfact Q-003: need\noption db: pg\noption web: go\nbet\n"

type reopening struct {
	*finishing
	rendered string
}

func newReopening(t *testing.T) *reopening {
	t.Helper()
	f, _, r := newBetFinishingWith(t, reclassAnswer)
	return &reopening{f, r}
}

// merge merges the pull request of the specification task outside, as a person does, and records its event merged.
func (w *reopening) merge() string {
	w.t.Helper()
	_, branch := w.specBranch()
	w.follow(branch)
	w.mergePull(w.pr(w.spec))
	if _, _, err := w.r.sessions().MergeTask(context.Background(), w.spec, w.rendered); err != nil {
		w.t.Fatalf("MergeTask: %v", err)
	}
	return strings.TrimSpace(gitOut(w.t, w.w.bare, "rev-parse", "refs/heads/main"))
}

func (w *reopening) open() ([]Task, error) { return w.r.openReopenings(context.Background(), w.spec) }

func (w *reopening) issues() int { return w.requests("POST /repos/acme/target/issues") }

func (w *reopening) rows(t *testing.T) [][]string {
	t.Helper()
	if !w.has(reopeningsPath) {
		return nil
	}
	rows, err := ReadReopenings([]byte(keptRecord(t, w.w, reopeningsPath)))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// openWatch calls on at each OpenIssue, before it is made.
type openWatch struct {
	forge.Forge
	on func()
}

func (o openWatch) OpenIssue(ctx context.Context, title, body string) (int, error) {
	o.on()
	return o.Forge.OpenIssue(ctx, title, body)
}

// I1 (the demo) with I4 and the observation of M9: after the event merged of the specification, one specification task is opened for each re-classed
// fact, with the identities of the record, which was committed before each issue; a restart opens nothing and writes nothing.
func TestAReclassedFactOpensOneSpecificationTaskAfterTheSpecificationMerges(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	w.merge()
	var recorded []bool
	w.r.cfg.Forge = openWatch{w.r.cfg.Forge, func() { recorded = append(recorded, w.has(reopeningsPath)) }}
	issues := w.issues()
	got, err := w.open()
	if err != nil || len(got) != 4 {
		t.Fatalf("openReopenings = %v, %v; want four tasks", got, err)
	}
	rows := w.rows(t)
	for i, want := range [][]string{{"1", "1", "constraint", "F-0004#1", "NFR-002"}, {"1", "2", "need", "F-0004#2", "REQ-003"}, {"1", "Q-002", "measure", "Q-002", ""},
		{"1", "Q-003", "need", "Q-003", "REQ-004"}} {
		if !slices.Equal(rows[i][:5], want) || rows[i][reoCopy] != "7004" || got[i].ID != rows[i][reoTask] {
			t.Errorf("row %d = %v, task %s; want %v with the copy 7004 and the task of the row", i, rows[i], got[i].ID, want)
		}
		reqs := []string{rows[i][reoRequirement]}
		if rows[i][reoRequirement] == "" {
			reqs = []string{"REQ-001"} // a measure names the requirements of the specification
		}
		if got[i].Kind != "specification" || !slices.Equal(got[i].Requirements, reqs) || !slices.Equal(got[i].Predecessors, []string{w.spec}) || got[i].Issue == 0 {
			t.Errorf("task %d = %+v; want a specification task for %v after %s", i, got[i], reqs, w.spec)
		}
		if w.has(pinPath(got[i].ID)) {
			t.Errorf("%s has an opening pin: a redo would compose the issue of the original specification", got[i].ID)
		}
	}
	if reg := w.register(); len(reg) != 7 || !slices.EqualFunc(reg[3:], got, func(a, b Task) bool { return a.ID == b.ID && a.Issue == b.Issue }) {
		t.Errorf("the register = %v; want the three tasks of the bet and the four reopened, in the order of the rows", reg)
	}
	if n := w.issues() - issues; n != 4 || !slices.Equal(recorded, []bool{true, true, true, true}) {
		t.Errorf("%d issues opened, with the record committed before each: %v; want 4 and true each time", n, recorded)
	}
	for i, issue := range [][2]string{w.w.fake.issueTexts[got[0].Issue-1], w.w.fake.issueTexts[got[2].Issue-1]} {
		id, token, what := got[i*2].ID, rows[i*2][reoToken], []string{"requirement", "success criterion"}[i]
		if !strings.HasPrefix(issue[0], "LAYUP task "+id+": the "+what) || !strings.Contains(issue[1], "the token "+token+".") || !strings.Contains(issue[1], "specification task "+w.spec+":") {
			t.Errorf("the issue of %s = %q; want the block reopen-issue for the %s %s", id, issue, what, token)
		}
	}
	head, posts := recordsHead(t, w.w), w.issues()
	w.restart()
	if again, err := w.open(); err != nil || len(again) != 4 || again[2].ID != got[2].ID || recordsHead(t, w.w) != head || w.issues() != posts {
		t.Errorf("a second call = %v, %v; head moved %v, issues %d (was %d); want the same four tasks and no write", again, err, recordsHead(t, w.w) != head, w.issues(), posts)
	}
}

// I2, a green negative beside I1: nothing is opened, written or asked of the forge before the event merged.
func TestNothingOpensBeforeTheMergedEvent(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	head, posts := recordsHead(t, w.w), w.issues()
	if got, err := w.open(); err != nil || len(got) != 0 || recordsHead(t, w.w) != head || w.issues() != posts || w.has(reopeningsPath) {
		t.Errorf("openReopenings before merged = %v, %v; want none and no write", got, err)
	}
	w.merge()
	if got, err := w.open(); err != nil || len(got) != 4 {
		t.Errorf("the control, after merged: %v, %v; want four tasks", got, err)
	}
}

// I3: a stop after the allocation and before any issue leaves the record; the restart opens one task for each row, under the IDs of the record.
func TestARestartAfterTheAllocationOpensOneTaskWithTheSameIDs(t *testing.T) {
	w := newReopening(t)
	w.merge()
	issues := w.issues()
	restore := stopAt("reopenings")
	_, err := w.open()
	restore()
	if !errors.Is(err, errStopped) || len(w.rows(t)) != 4 || w.issues() != issues || len(w.register()) != 3 {
		t.Fatalf("the stop: %v, %d rows, %d issues, %d tasks; want the stop with the allocation and nothing else", err, len(w.rows(t)), w.issues()-issues, len(w.register()))
	}
	before := keptRecord(t, w.w, reopeningsPath)
	w.restart()
	got, err := w.open()
	if err != nil || len(got) != 4 || keptRecord(t, w.w, reopeningsPath) != before || w.issues() != issues+4 || len(w.register()) != 7 {
		t.Fatalf("the restart = %v, %v; want four tasks under the allocated IDs, 4 issues and the record unchanged", got, err)
	}
	for i, r := range w.rows(t) {
		if got[i].ID != r[reoTask] {
			t.Errorf("task %d is %s; want %s", i, got[i].ID, r[reoTask])
		}
	}
}

// I5, the known limit: a stop between the issue and the row of the first task opens a second issue under the same ID and writes one row.
func TestAStopBetweenTheIssueAndTheRowOpensASecondIssueAndOneRow(t *testing.T) {
	w := newReopening(t)
	w.merge()
	issues := w.issues()
	restore := stopAt("issue")
	_, err := w.open()
	restore()
	if !errors.Is(err, errStopped) || w.issues() != issues+1 || len(w.register()) != 3 || len(w.rows(t)) != 4 {
		t.Fatalf("the stop: %v, %d issues, %d tasks; want the stop after one issue, no row and the allocation", err, w.issues()-issues, len(w.register()))
	}
	w.restart()
	got, err := w.open()
	ids := map[string]int{}
	for _, tk := range w.register() {
		ids[tk.ID]++
	}
	if err != nil || len(got) != 4 || w.issues() != issues+5 || len(w.register()) != 7 || ids[got[0].ID] != 1 {
		t.Errorf("the restart = %v, %v; %d issues, %d tasks; want 5 issues (the second of the first task is the known limit) and one row for each task", got, err, w.issues()-issues, len(w.register()))
	}
}

// I6: a reopened task is not the specification task returned by its kind, and it is not redone from the issue of the original.
func TestTheReopeningIsNotTheSpecificationReusedByKind(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	w.merge()
	// A task of the register that names a higher requirement than the draft does: the numbers go on from it (one number space for the register too).
	if err := w.r.registerTask(context.Background(), Task{ID: "T-bld1", Issue: 90, Kind: "build", Requirements: []string{"REQ-009"}, Opened: time0}); err != nil {
		t.Fatal(err)
	}
	got, err := w.open()
	if err != nil || len(got) != 4 {
		t.Fatalf("openReopenings = %v, %v; want four tasks", got, err)
	}
	if rows := w.rows(t); rows[0][reoRequirement] != "NFR-010" || rows[1][reoRequirement] != "REQ-011" || rows[3][reoRequirement] != "REQ-012" {
		t.Errorf("the requirements %v; want NFR-010, REQ-011 and REQ-012, after the REQ-009 of the register", rows)
	}
	n := 0
	for _, tk := range w.register() {
		if tk.Kind == "specification" {
			n++
		}
	}
	if n != 5 || slices.ContainsFunc(got, func(tk Task) bool { return tk.ID == w.spec }) {
		t.Errorf("%d specification tasks; want the original and four of their own", n)
	}
	s, ctx := w.r.sessions(), context.Background()
	if _, _, err := s.redoIssue(ctx, got[0], "T-newx", nil); err == nil || !strings.Contains(err.Error(), "no opening pin") {
		t.Errorf("the redo of a reopened task: %v; want the refusal for its missing pin", err)
	}
	orig, _ := s.taskRow(ctx, w.spec)
	if title, _, err := s.redoIssue(ctx, orig, "T-newx", nil); err != nil || !strings.Contains(title, "the specification") {
		t.Errorf("the control, the redo of the original: %q, %v", title, err)
	}
}

// I7: the rendered head of a bet that re-classes facts passes the four statuses, so the specification merges by FinishTask (K4); the full check
// there exempts the facts that no task reopens yet.
func TestTheRenderedHeadOfAReclassingBetPassesLayupSpec(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	got, err := w.finishTask(w.spec, w.harnesses, nil)
	if err != nil || !got.Merged {
		t.Fatalf("FinishTask of the specification = %+v, %v; want Merged", got, err)
	}
	if st := headStates(w.drvWorld, w.rendered); !slices.Contains(st, SpecContext+" success") {
		t.Errorf("the statuses at the rendered commit: %v; want layup/spec success", st)
	}
}

// I10: a lease lost after the allocation opens no issue.
func TestALostLeaseOpensNoIssue(t *testing.T) {
	w := newReopening(t)
	w.merge()
	issues := w.issues()
	restore := stopAt("reopenings")
	w.open()
	restore()
	if err := w.r.store.WriteLease(context.Background(), LeaseRow{Run: "fedcba9876543210", Host: "host-2", Version: "0.1.0-dev", Started: time.Now().UTC().Truncate(time.Second), Heartbeat: 1, State: "held"}, "takeover"); err != nil {
		t.Fatal(err)
	}
	w.restart()
	var lost *LostError
	if _, err := w.open(); !errors.As(err, &lost) || w.issues() != issues || len(w.register()) != 3 {
		t.Errorf("openReopenings = %v; %d issues, %d tasks; want a LostError, no issue and no task", err, w.issues()-issues, len(w.register()))
	}
}

// Condition 1 of the plan review: two reopened tasks reach green checks in order. A fact is exempt at every head but its own task's until its task merged;
// the head of a task without its line fails, and the head of a task that removes a merged line fails.
func TestTheReopenedTasksPassTheFullCheckOneAfterTheOther(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	main := w.merge()
	tasks, err := w.open()
	if err != nil || len(tasks) != 4 {
		t.Fatalf("openReopenings = %v, %v; want four tasks", tasks, err)
	}
	b, d := tasks[1], tasks[3] // the facts 2 and Q-003, both needs
	ctx, s := context.Background(), w.r.sessions()
	prd := gitOut(t, w.w.bare, "show", main+":docs/prd/PRD-0001-x.md")
	line := func(id, token string) string {
		return id + ": a statement. Facts: " + token + " F-0004. Acceptance criterion: a check.\n"
	}
	headWith := func(lines ...string) string {
		if len(lines) == 0 {
			return main // the head that adds nothing
		}
		text := prd + "\n## Waits for a bet\n\n" + strings.Join(lines, "")
		h, err := w.r.commitFiles(main, map[string][]byte{"docs/prd/PRD-0001-x.md": []byte(text)}, "a head", w.r.who(time.Unix(1, 0)))
		if err != nil {
			t.Fatal(err)
		}
		return h
	}
	copyText, err := w.r.betCopy(ctx)
	if err != nil {
		t.Fatal(err)
	}
	needs := func(task, head string) string {
		t.Helper()
		if _, err := s.SpecAt(ctx, SpecInput{Task: task, Base: main, Head: head, Copy: copyText}); err != nil {
			t.Fatal(err)
		}
		rows, _, err := s.SpecTable(ctx, task, head)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range rows {
			if r[0] == "needs" {
				return r[1] + " " + r[2]
			}
		}
		return "no row"
	}
	lb, ld := line(b.Requirements[0], "F-0004#2"), line(d.Requirements[0], "Q-003")
	for _, c := range []struct{ name, task, head, want string }{
		{"the head of the first task, with its line", b.ID, headWith(lb), "pass "},
		{"the head of the first task, without its line", b.ID, headWith(), "fail the need fact 2 has no requirement"},
		{"the head of the second task, with the line of the first only", d.ID, headWith(lb), "fail the need fact Q-003 has no requirement"},
		{"the head of the second task, with its line", d.ID, headWith(ld), "pass "},
	} {
		if got := needs(c.task, c.head); got != c.want {
			t.Errorf("%s: needs = %q; want %q", c.name, got, c.want)
		}
	}
	// The first task merges: its fact is no longer exempt anywhere.
	data, err := s.appendEvents(ctx, b.ID, 1, "", []string{"merged " + headWith(lb) + " 99"})
	if err == nil {
		err = s.commit(ctx, map[string][]byte{eventsPath(b.ID): data}, "the first reopened task merged")
	}
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ name, head, want string }{
		{"the second task, on a base with the first line and its own", headWith(lb, ld), "pass "},
		{"the second task, on a base that lost the first line", headWith(strings.Replace(ld, "a statement", "another statement", 1)), "fail the need fact 2 has no requirement"},
	} {
		if got := needs(d.ID, c.head); got != c.want {
			t.Errorf("%s: needs = %q; want %q", c.name, got, c.want)
		}
	}
}

// I9: the block reopenings equals its Go schema, and the blocks reopen-issue and instruction-specification-reopened equal the code.
func TestTheReopeningsBlockAndTheNewTextBlocksEqualTheCode(t *testing.T) {
	t.Parallel()
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tsv.Compare(blocks["reopenings"], ReopeningsSchema); err != nil {
		t.Errorf("the block reopenings and ReopeningsSchema: %v", err)
	}
	doc, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	for name, text := range map[string]string{"reopen-issue": reopenIssueBlock, "instruction-specification-reopened": reopenedInstruction} {
		_, block, ok := strings.Cut(string(doc), "\n```text "+name+"\n")
		block, _, _ = strings.Cut(block, "```\n")
		if !ok || block != text || text == "" {
			t.Errorf("loop.md has the block %s: %v; it is equal to the code: %v", name, ok, block == text && text != "")
		}
	}
}

// I8: a harness that follows the reopened instruction (it writes the line that the prompt shows, and the section of the requirement) gets a valid
// handoff and a head whose layup/spec is green; the pinned prd-lint accepts the head, and the line stands under the heading.
func TestAHarnessThatFollowsTheReopenedInstructionGivesAValidHandoffAndAGreenSpecStatus(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	main := w.merge()
	tasks, err := w.open()
	if err != nil || len(tasks) != 4 {
		t.Fatalf("openReopenings = %v, %v; want four tasks", tasks, err)
	}
	b, ctx := tasks[1], context.Background() // the fact 2, a need: REQ-003
	in := w.input(b.ID)
	in.Base = main // the base of the driver is the default branch after the merge
	pc, err := w.r.sessions().PromptFile(ctx, PromptIn{Task: b.ID, Step: "specification", Base: in.Base})
	if err != nil {
		t.Fatal(err)
	}
	if items := itemsOf(pc, session.KindRequirements); !slices.Contains(items, "spec/reopenings.tsv") || !slices.Contains(items, "docs/prd/PRD-0001-x.md") {
		t.Errorf("the requirements of the prompt = %v; want the record of the reopenings and the PRD, which holds no row of %s", items, b.Requirements[0])
	}
	line := shownLine(string(pc.Text("/x")), "need")
	w.ctlWrite("reopen-"+b.ID, line)
	w.ctlWrite("reopen-section-"+b.ID, b.Requirements[0])
	if in.Copy, err = w.r.betCopy(ctx); err != nil {
		t.Fatal(err)
	}
	got, err := w.r.sessions().AdvanceShapeTask(ctx, in)
	if err != nil || !got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the verified head (handoff = %v)", got, err, w.kinds(b.ID))
	}
	for _, e := range w.eventsOf(b.ID) {
		if e[1] == "handoff" && !strings.HasPrefix(e[6], "valid") || e[1] == "refused" {
			t.Errorf("the event %v; want only valid handoffs", e)
		}
	}
	rows, found, err := w.r.sessions().SpecTable(ctx, b.ID, got.Head)
	if err != nil || !found {
		t.Fatalf("no spec table at the head: %v", err)
	}
	for _, r := range rows {
		if r[1] == "fail" {
			t.Errorf("the rule %s fails at the head: %s", r[0], r[2])
		}
	}
	if prd := gitOut(t, w.r.target, "show", got.Head+":docs/prd/PRD-0001-x.md"); !strings.Contains(prd, "\n## Waits for a bet\n\n"+line) {
		t.Errorf("the PRD at the head lacks the line under the heading %q:\n%s", line, prd)
	}
	if ok, lines, err := w.r.sessions().waitingSpecAccepted(in.Base, got.Head); err != nil || !ok {
		t.Errorf("the pinned prd-lint at the head: %v, %q, %v; want it accepted", ok, lines, err)
	}
}

// Condition 4 of the plan review: a task of the register that holds a reserved ID must agree with the reservation, and the rows of the record must
// agree with the version; else nothing is opened.
func TestARegisteredTaskThatDisagreesWithItsReservationIsRefused(t *testing.T) {
	t.Parallel()
	w := newReopening(t)
	w.merge()
	rows := [][]string{{"1", "1", "constraint", "F-0004#1", "NFR-002", "T-rsv1", "7004"}, {"1", "2", "need", "F-0004#2", "REQ-003", "T-rsv2", "7004"},
		{"1", "Q-002", "measure", "Q-002", "", "T-rsv3", "7004"}, {"1", "Q-003", "need", "Q-003", "REQ-004", "T-rsv4", "7004"}}
	data, err := table(ReopeningsSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	w.push(map[string][]byte{reopeningsPath: data}, "a reservation")
	w.restart()
	if err := w.r.registerTask(context.Background(), Task{ID: "T-rsv1", Issue: 91, Kind: "specification", Requirements: []string{"REQ-009"}, Predecessors: []string{w.spec}, Opened: time0}); err != nil {
		t.Fatal(err)
	}
	issues := w.issues()
	if got, err := w.open(); err == nil || !strings.Contains(err.Error(), "T-rsv1") || len(got) != 0 || w.issues() != issues {
		t.Errorf("openReopenings = %v, %v; %d issues; want an error that names T-rsv1 and no issue", got, err, w.issues()-issues)
	}
	rows[1][reoClass], rows[1][reoRequirement] = "constraint", "NFR-003"
	if data, err = table(ReopeningsSchema, rows); err != nil {
		t.Fatal(err)
	}
	w.push(map[string][]byte{reopeningsPath: data}, "another reservation")
	w.restart()
	if got, err := w.open(); err == nil || !strings.Contains(err.Error(), "class") || len(got) != 0 || w.issues() != issues {
		t.Errorf("openReopenings with a class the version does not give = %v, %v; want the refusal of the class and no issue", got, err)
	}
}
