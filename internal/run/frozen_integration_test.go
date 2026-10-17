//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/tsv"
)

// A loopRun is the world of #361: the target with its first tree, the register
// that opens T-ab12 on issue 3, the lease held, the App's forge, and the
// sessions of the run on a real clone and a real bare repository.
type loopRun struct {
	t *testing.T
	w *world
	r *Sessions
	a worldAt
}

func newLoop(t *testing.T, baseline map[string]string) *loopRun {
	t.Helper()
	w, r, a := commentWorldWith(t, baseline)
	return &loopRun{t: t, w: w, r: r, a: a}
}

// prod gives the Sessions from the constructor that a run uses, state.sessions(),
// with the clock of l.r, not a literal (guardrails §2, #360).
func (l *loopRun) prod() *Sessions {
	st := &state{cfg: l.w.config("0.1.0-dev"), store: l.r.Store.(*recordsStore), target: l.r.Clone, defaultBranch: l.r.Branch,
		start: [][]string{{"issue.control", "2", "run"}}}
	s := st.sessions()
	s.Now = l.r.Now
	return s
}

// fixed makes every time of the run the one second, so that nothing is ordered by the clock.
func (l *loopRun) fixed() {
	at := time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)
	l.r.Now = func() time.Time { return at }
}

// fileWork is the shell of a harness that writes the files in repo/, commits them and writes the result.
func fileWork(files map[string]string, result string) string {
	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	var b strings.Builder
	for _, p := range paths {
		b.WriteString("mkdir -p \"$(dirname '" + p + "')\" && printf '" + octal(files[p]) + "' > '" + p + "' && ")
	}
	b.WriteString("git add -A && git commit -qm work && printf '" + octal(result) + "' > ../result/result.tsv")
	return b.String()
}

func (l *loopRun) attempt(n int) {
	l.t.Helper()
	if got, err := l.r.StartAttempt(context.Background(), "T-ab12", l.a.base); err != nil || got != n {
		l.t.Fatalf("StartAttempt = %d, %v; want %d", got, err, n)
	}
}

// dev runs a developer session of the attempt on the harness, with the frozen
// files that the records hold, as the loop will call it, and gives its ID.
func (l *loopRun) dev(attempt int, harnessID string, files map[string]string, result string) string {
	l.t.Helper()
	ctx := context.Background()
	frozen, err := l.r.FrozenFiles(ctx, "T-ab12")
	if err != nil {
		l.t.Fatalf("FrozenFiles: %v", err)
	}
	spec := harness(l.t, l.a, fileWork(files, result))
	spec.Role, spec.Attempt, spec.Pair.Harness, spec.Frozen = roleOfStep("implement"), attempt, harnessID, frozen
	id, err := l.r.TaskSession(ctx, spec)
	if err != nil {
		l.t.Fatalf("TaskSession: %q, %v", id, err)
	}
	return id
}

const okResult = resultHead

// bound gives the SHA that the event bound of the session holds, from the records.
func (l *loopRun) bound(id string) string {
	l.t.Helper()
	for _, e := range l.events() {
		if e[1] == "bound" && e[3] == id {
			return e[5]
		}
	}
	return ""
}

func (l *loopRun) events() [][]string {
	l.t.Helper()
	rows, err := records.ReadEvents([]byte(keptRecord(l.t, l.w, "tasks/T-ab12/events.tsv")))
	if err != nil {
		l.t.Fatal(err)
	}
	return rows
}

func (l *loopRun) testsFile() string {
	l.t.Helper()
	return gitOut(l.t, l.w.bare, "ls-tree", "--name-only", "layup-records", "tasks/T-ab12/tests.tsv")
}

func (l *loopRun) tests() [][]string {
	l.t.Helper()
	rows, err := ReadTests([]byte(keptRecord(l.t, l.w, "tasks/T-ab12/tests.tsv")))
	if err != nil {
		l.t.Fatalf("tests.tsv: %v", err)
	}
	return rows
}

func green(id, src string) TestResult {
	return TestResult{ID: id, Source: src, FailsAtBase: true, PassesAtHead: true}
}

// toVerifier records the handoff of a developer session to the verifier.
func (l *loopRun) toVerifier(s *Sessions, id string, freeze *Freeze, items map[string]bool) (handoff.Verdict, error) {
	l.t.Helper()
	head := l.bound(id)
	if head == "" {
		head = strings.Repeat("e", 40)
	}
	return s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: id, From: "developer (build task)", To: "verifier",
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: head, Found: head}, {Kind: "task-file", Path: "docs/tasks/T-ab12.md", Named: hashA, Found: hashA}},
		Items:     items, Freeze: freeze})
}

// comment runs a session of the step that writes a comment, and posts it.
func (l *loopRun) comment(attempt int, step, harnessID, text string) string {
	l.t.Helper()
	ctx := context.Background()
	spec := harness(l.t, l.a, planWork("plan.md", text, "comment\t1\tresult/plan.md\t"+commentSum(text)+"\t—"))
	spec.Role, spec.Attempt, spec.Pair.Harness = roleOfStep(step), attempt, harnessID
	id, err := l.r.TaskSession(ctx, spec)
	if err != nil {
		l.t.Fatalf("TaskSession %s: %q, %v", step, id, err)
	}
	if _, err := l.r.PostComment(ctx, "T-ab12", id, "plan.md"); err != nil {
		l.t.Fatalf("PostComment: %v", err)
	}
	return id
}

func (l *loopRun) handoffOf(id, from, to, kind, item string) handoff.Verdict {
	l.t.Helper()
	v, err := l.r.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: id, From: from, To: to, File: "plan.md", Kind: kind, Items: map[string]bool{item: true}})
	if err != nil {
		l.t.Fatalf("RecordHandoff %s: %v", kind, err)
	}
	return v
}

func (l *loopRun) plan(attempt int, harnessID, text string) string {
	l.t.Helper()
	id := l.comment(attempt, "plan-of-a-task", harnessID, text)
	if v := l.handoffOf(id, "task plan", "plan review", "plan", "plan-fields"); !v.Valid {
		l.t.Fatalf("the plan handoff: %+v", v)
	}
	return id
}

func (l *loopRun) review(attempt int, harnessID, text string) string {
	l.t.Helper()
	id := l.comment(attempt, "plan-review", harnessID, text)
	if v := l.handoffOf(id, "plan review", "developer", "plan-review", "verdict-field"); !v.Valid {
		l.t.Fatalf("the plan review handoff: %+v", v)
	}
	return id
}

const (
	aOne = "package x\n\n// TestA one\n"
	aTwo = "package x\n\n// TestA two\n"
	bOne = "package x\n\n// TestB one\n"
)

var listAB = &Freeze{Tests: []TestResult{green("TestA", "a_test.go"), green("TestB", "b_test.go")}}

// frozenAt freezes the list of a_test.go and b_test.go by a first attempt of the harness.
func (l *loopRun) frozenAt(harnessID string) string {
	l.t.Helper()
	id := l.dev(1, harnessID, map[string]string{"a_test.go": aOne, "b_test.go": bOne}, okResult)
	if v, err := l.toVerifier(l.r, id, listAB, nil); err != nil || !v.Valid {
		l.t.Fatalf("the first handoff to the verifier: %+v, %v", v, err)
	}
	return id
}

func TestTheTestsBlockEqualsItsGoSchema(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["task-tests"]
	if !ok {
		t.Fatal("docs/spec/ has no block task-tests")
	}
	if err := tsv.Compare(block, TestsSchema); err != nil {
		t.Error(err)
	}
}

// The demo of loop.md, the frozen test list: the first valid handoff to the
// verifier freezes the list, in the commit of its event, with the hash of each
// source at the head that the session was bound to.
func TestTheFrozenListHoldsTheHashesOfTheFilesAtTheBoundHead(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne, "b_test.go": bOne}, okResult)
	if l.testsFile() != "" {
		t.Fatal("a tests.tsv before the handoff")
	}
	v, err := l.toVerifier(l.prod(), id, listAB, map[string]bool{"tests-red-green": false})
	if err != nil || !v.Valid || v.Detail != "valid" {
		t.Fatalf("the handoff = %+v, %v; want valid (the caller's own claim for tests-red-green is not read)", v, err)
	}
	rows := l.tests()
	want := [][]string{{"TestA", "a_test.go", sumOf([]byte(aOne)), "claude"}, {"TestB", "b_test.go", sumOf([]byte(bOne)), "claude"}}
	if len(rows) != 2 {
		t.Fatalf("tests.tsv = %v", rows)
	}
	for i := range want {
		if !slices.Equal(rows[i][:4], want[i]) || rows[i][4] != rows[0][4] {
			t.Errorf("row %d = %v, want %v and one time", i, rows[i], want[i])
		}
	}
	commit := strings.TrimSpace(gitOut(t, l.w.bare, "log", "-1", "--format=%H", "layup-records", "--", "tasks/T-ab12/tests.tsv"))
	changed := strings.Fields(gitOut(t, l.w.bare, "show", "--format=", "--name-only", commit))
	if !slices.Equal(changed, []string{"tasks/T-ab12/events.tsv", "tasks/T-ab12/tests.tsv"}) {
		t.Errorf("the commit of the list changed %v; want the event and the list in one commit", changed)
	}
	if ev := l.events(); ev[len(ev)-1][1] != "handoff" || ev[len(ev)-1][6] != "valid" {
		t.Errorf("the last event %v; want the handoff valid", ev[len(ev)-1])
	}
}

// A session that was never bound (its push was refused) has no head: the run's
// own error, and neither the event nor the list is written.
func TestASessionWithNoBoundHeadFreezesNothing(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne, ".github/workflows/ci.yml": "on: push\n"}, okResult)
	if l.bound(id) != "" {
		t.Fatal("the session was bound")
	}
	n := len(l.events())
	if _, err := l.toVerifier(l.r, id, &Freeze{Tests: []TestResult{green("TestA", "a_test.go")}}, nil); !errors.Is(err, ErrNoBoundHead) {
		t.Errorf("a handoff of a session with no bound head: %v; want ErrNoBoundHead", err)
	}
	if l.testsFile() != "" || len(l.events()) != n {
		t.Errorf("a list or an event was written for a session that has no head (%d events, was %d)", len(l.events()), n)
	}
}

// loop.md: a test that fails at the head, or passes at the base, makes the
// handoff invalid, so the list is never frozen on a test that proves nothing.
func TestAListedTestThatIsNotRedThenGreenMakesTheHandoffInvalidAndFreezesNothing(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne}, okResult)
	for _, c := range []struct {
		name  string
		tests []TestResult
	}{
		{"it passes at the base", []TestResult{{ID: "TestA", Source: "a_test.go", FailsAtBase: false, PassesAtHead: true}}},
		{"it fails at the head", []TestResult{{ID: "TestA", Source: "a_test.go", FailsAtBase: true, PassesAtHead: false}}},
		{"it was not run", []TestResult{{ID: "TestA", Source: "a_test.go"}}},
		{"the list is empty", nil},
	} {
		v, err := l.toVerifier(l.r, id, &Freeze{Tests: c.tests}, map[string]bool{"tests-red-green": true})
		if err != nil || v.Valid || v.Detail != "invalid tests-red-green" {
			t.Errorf("%s: %+v, %v; want invalid tests-red-green (the caller's true is not read)", c.name, v, err)
		}
		if l.testsFile() != "" {
			t.Errorf("%s: a list was frozen", c.name)
		}
	}
}

func TestAnInvalidHandoffFreezesNothing(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne}, okResult)
	head := l.bound(id)
	in := HandoffInput{Task: "T-ab12", Session: id, From: "developer (build task)", To: "verifier", Freeze: &Freeze{Tests: []TestResult{green("TestA", "a_test.go")}},
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: head, Found: strings.Repeat("e", 40)}, {Kind: "task-file", Named: hashA, Found: hashA}}}
	v, err := l.r.RecordHandoff(context.Background(), in)
	if err != nil || v.Valid || v.Detail != "invalid artifacts" {
		t.Fatalf("the handoff = %+v, %v; want invalid artifacts", v, err)
	}
	if l.testsFile() != "" {
		t.Error("a list was frozen by an invalid handoff")
	}
}

// A table of the Operator's that leaves out tests-red-green is the Operator's:
// the handoff is as the table says, and the list is still never frozen on a
// test that proves nothing.
func TestATableWithoutTestsRedGreenNeverFreezesOnAFailingTest(t *testing.T) {
	l := newLoop(t, nil)
	rows := handoff.Default()
	for i := range rows {
		if rows[i].From == "developer (build task)" {
			rows[i].Items = []string{"schema", "artifacts"}
		}
	}
	var buf bytes.Buffer
	if err := handoff.WriteTable(&buf, rows); err != nil {
		t.Fatal(err)
	}
	if err := l.r.commit(context.Background(), map[string][]byte{transitionsPath: buf.Bytes()}, "the Operator's table"); err != nil {
		t.Fatal(err)
	}
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne}, okResult)
	v, err := l.toVerifier(l.r, id, &Freeze{Tests: []TestResult{{ID: "TestA", Source: "a_test.go", FailsAtBase: false, PassesAtHead: true}}}, nil)
	if err != nil || !v.Valid {
		t.Fatalf("the handoff = %+v, %v; want valid under the Operator's table", v, err)
	}
	if l.testsFile() != "" {
		t.Error("the list was frozen on a test that passes at the base")
	}
}

// Another valid handoff, with nothing changed, writes the list no more.
func TestASecondValidHandoffDoesNotWriteTheListAgain(t *testing.T) {
	l := newLoop(t, nil)
	l.frozenAt("claude")
	before := keptRecord(t, l.w, "tasks/T-ab12/tests.tsv")
	l.attempt(2)
	id := l.dev(2, "devin", map[string]string{"c.txt": "other\n"}, okResult)
	v, err := l.toVerifier(l.r, id, listAB, nil)
	if err != nil || !v.Valid {
		t.Fatalf("the second handoff = %+v, %v", v, err)
	}
	if after := keptRecord(t, l.w, "tasks/T-ab12/tests.tsv"); after != before {
		t.Errorf("tests.tsv changed:\n%s\nwas:\n%s", after, before)
	}
	if n := len(strings.Fields(gitOut(t, l.w.bare, "log", "--format=%H", "layup-records", "--", "tasks/T-ab12/tests.tsv"))); n != 1 {
		t.Errorf("%d commits wrote tests.tsv; want one", n)
	}
}

// loop.md: a later head where the source of a listed test differs is a finding
// "test changed", and its handoff to the verifier is invalid: the test is not
// counted as passed.
func TestAChangedFrozenSourceIsAFindingAndTheHandoffIsInvalid(t *testing.T) {
	l := newLoop(t, nil)
	first := l.frozenAt("claude")
	before := keptRecord(t, l.w, "tasks/T-ab12/tests.tsv")
	l.attempt(2)
	id := l.dev(2, "devin", map[string]string{"a_test.go": aTwo}, okResult)
	v, err := l.toVerifier(l.prod(), id, listAB, nil)
	if err != nil || v.Valid || v.Detail != "invalid tests-red-green" {
		t.Fatalf("the handoff of a changed test = %+v, %v; want invalid tests-red-green (every test is green by the caller)", v, err)
	}
	got, err := l.prod().FrozenFindings(context.Background(), "T-ab12", l.bound(id))
	if err != nil || len(got) != 1 || got[0].Rule != "test changed" || got[0].File != "a_test.go" || got[0].Text != "TestA" {
		t.Errorf("FrozenFindings = %+v, %v; want the one finding test changed of a_test.go", got, err)
	}
	if got, err := l.prod().FrozenFindings(context.Background(), "T-ab12", l.bound(first)); err != nil || len(got) != 0 {
		t.Errorf("FrozenFindings at the freezing head = %+v, %v; want none", got, err)
	}
	if after := keptRecord(t, l.w, "tasks/T-ab12/tests.tsv"); after != before {
		t.Error("a changed test refroze the list with no amendment")
	}
}

// amended runs the amendment of a changed test: the plan, then its review, each
// by a harness of its own; and the handoff of the changed head that comes next.
func (l *loopRun) changedAttempt() string {
	l.attempt(2)
	id := l.dev(2, "devin", map[string]string{"a_test.go": aTwo}, okResult)
	if v, err := l.toVerifier(l.r, id, listAB, nil); err != nil || v.Valid {
		l.t.Fatalf("the handoff of a changed test: %+v, %v; want invalid", v, err)
	}
	return id
}

// The next attempt of the developer after an amendment writes the new version of the test, and its handoff freezes again.
func (l *loopRun) afterAmendment() (handoff.Verdict, error, string) {
	l.attempt(3)
	id := l.dev(3, "codex", map[string]string{"a_test.go": aTwo}, okResult)
	v, err := l.toVerifier(l.r, id, listAB, nil)
	return v, err, id
}

func amendmentFlow(t *testing.T, fixed bool) {
	l := newLoop(t, nil)
	if fixed {
		l.fixed()
	}
	l.frozenAt("claude")
	first := l.tests()
	l.changedAttempt()
	l.plan(2, "claude", "plan v2\n")
	l.review(2, "devin", "review of plan v2\n")
	v, err, id := l.afterAmendment()
	if err != nil || !v.Valid {
		t.Fatalf("the handoff after an amendment = %+v, %v; want valid", v, err)
	}
	now := l.tests()
	if now[0][2] != sumOf([]byte(aTwo)) || now[0][3] != "codex" || now[1][2] != first[1][2] || now[1][3] != "claude" {
		t.Errorf("the list after the amendment %v\nwas %v; want TestA at its new hash by codex and TestB as it was, by claude", now, first)
	}
	if got, err := l.r.FrozenFindings(context.Background(), "T-ab12", l.bound(id)); err != nil || len(got) != 0 {
		t.Errorf("the finding is not closed: %+v, %v", got, err)
	}
	// The amendment is consumed: a next changed head, with no new plan, is invalid again.
	l.attempt(4)
	id4 := l.dev(4, "devin", map[string]string{"a_test.go": aOne}, okResult)
	if v, err := l.toVerifier(l.r, id4, listAB, nil); err != nil || v.Valid {
		t.Errorf("a second change after the amendment was spent: %+v, %v; want invalid", v, err)
	}
}

func TestAnAmendmentFreezesTheListAgainAndTheFindingCloses(t *testing.T) { amendmentFlow(t, false) }

// The two freezings, the plan and its review all fall in one second: the order
// is the events', never the clock's (guardrails §2).
func TestTwoFreezesInTheSameSecondAreOrderedByTheEvents(t *testing.T) { amendmentFlow(t, true) }

func TestWithoutAnAmendmentAChangedTestIsNotFrozenAgain(t *testing.T) {
	for _, c := range []struct {
		name string
		run  func(l *loopRun)
	}{
		{"a plan with no review", func(l *loopRun) { l.plan(2, "claude", "plan v2\n") }},
		{"a review with no plan", func(l *loopRun) { l.review(2, "devin", "review\n") }},
		{"a review before the plan", func(l *loopRun) {
			l.review(2, "devin", "review\n")
			l.plan(2, "claude", "plan v2\n")
		}},
		{"a verifier's record that started before the plan", func(l *loopRun) {
			ctx := context.Background()
			spec := harness(l.t, l.a, planWork("plan.md", "record\n", "comment\t1\tresult/plan.md\t"+commentSum("record\n")+"\t—"))
			spec.Role, spec.Attempt, spec.Pair.Harness = roleOfStep("verification"), 2, "devin"
			vid, err := l.r.TaskSession(ctx, spec)
			if err != nil {
				l.t.Fatal(err)
			}
			if _, err := l.r.PostComment(ctx, "T-ab12", vid, "plan.md"); err != nil {
				l.t.Fatal(err)
			}
			l.plan(2, "claude", "plan v2\n")
			v, err := l.r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: vid, From: "verifier", To: "close-out", File: "plan.md", Kind: "review-record", Items: map[string]bool{"nothing-material": true}})
			if err != nil || !v.Valid {
				l.t.Fatalf("the stale record's handoff: %+v, %v", v, err)
			}
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			l := newLoop(t, nil)
			l.frozenAt("claude")
			l.changedAttempt()
			c.run(l)
			l.attempt(3)
			id := l.dev(3, "codex", map[string]string{"a_test.go": aTwo}, okResult)
			if v, err := l.toVerifier(l.r, id, listAB, nil); err != nil || v.Valid || v.Detail != "invalid tests-red-green" {
				t.Errorf("the handoff = %+v, %v; want invalid tests-red-green: this is no amendment", v, err)
			}
		})
	}
}

// lostStore is a records store whose push is refused and whose lease row is
// another run's: the state in which Fenced stops a run that lost the lease.
type lostStore struct{ SessionStore }

func (lostStore) Commit(context.Context, map[string][]byte, string) error { return ErrRefused }
func (lostStore) ReadLease(context.Context) (LeaseRow, error) {
	return LeaseRow{Run: "fedcba9876543210", State: "held"}, nil
}

// A lost lease writes nothing: neither the event nor the list.
func TestALostLeaseWritesNoFrozenList(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne}, okResult)
	n := len(l.events())
	lost := afterRestart(l.r)
	lost.Store = lostStore{l.r.Store}
	_, err := l.toVerifier(lost, id, &Freeze{Tests: []TestResult{green("TestA", "a_test.go")}}, nil)
	var le *LostError
	if !errors.As(err, &le) || le.Holder != "fedcba9876543210" {
		t.Errorf("the handoff of a run that lost the lease: %v; want a LostError naming the other run", err)
	}
	if l.testsFile() != "" || len(l.events()) != n {
		t.Error("a record was written by a run that lost the lease")
	}
}

// A restart that records the same handoff again writes neither a second event
// nor a second list.
func TestARestartRecordsTheSameHandoffOnceAndTheListStays(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"a_test.go": aOne}, okResult)
	freeze := &Freeze{Tests: []TestResult{green("TestA", "a_test.go")}}
	if v, err := l.toVerifier(l.prod(), id, freeze, nil); err != nil || !v.Valid {
		t.Fatalf("first: %+v, %v", v, err)
	}
	list := keptRecord(t, l.w, "tasks/T-ab12/tests.tsv")
	n := len(l.events())
	if v, err := l.toVerifier(l.prod(), id, freeze, nil); err != nil || !v.Valid {
		t.Fatalf("restart: %+v, %v", v, err)
	}
	if len(l.events()) != n || keptRecord(t, l.w, "tasks/T-ab12/tests.tsv") != list {
		t.Errorf("a restart wrote again (%d events, was %d)", len(l.events()), n)
	}
}

// loop.md, the attempts: code puts the frozen test files at the start of each
// new attempt and commits them, so that their hashes stay in the head.
func TestTheFrozenFilesAreInTheNextAttemptAtTheirHashes(t *testing.T) {
	l := newLoop(t, nil)
	l.fixed()
	l.frozenAt("claude")
	l.attempt(2)
	id := l.dev(2, "devin", map[string]string{"c.txt": "other\n"}, okResult)
	head := l.bound(id)
	for path, want := range map[string]string{"a_test.go": aOne, "b_test.go": bOne} {
		if got := gitOut(t, l.w.bare, "show", head+":"+path); got != want {
			t.Errorf("%s at the head = %q, want the frozen bytes %q", path, got, want)
		}
	}
	log := strings.Split(strings.TrimSpace(gitOut(t, l.w.bare, "log", "--reverse", "--format=%an|%at|%s", l.a.base+".."+head)), "\n")
	if len(log) != 2 {
		t.Fatalf("the commits of the attempt: %q; want the frozen tests, then the work", log)
	}
	wantAt := strconv.FormatInt(time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC).Unix(), 10)
	f := strings.Split(log[0], "|")
	if !strings.HasPrefix(f[0], "layup session S-") || f[1] != wantAt || f[2] != "layup: the frozen tests" {
		t.Errorf("the first commit %q; want the session's identity, the time of the run and the message of the frozen tests", log[0])
	}
}

func TestNoListPutsNoFileInAnAttempt(t *testing.T) {
	l := newLoop(t, nil)
	id := l.dev(1, "claude", map[string]string{"c.txt": "work\n"}, okResult)
	log := strings.Fields(gitOut(t, l.w.bare, "log", "--format=%H", l.a.base+".."+l.bound(id)))
	if len(log) != 1 {
		t.Errorf("%d commits in an attempt with no frozen list; want the work only", len(log))
	}
}

// The frozen bytes are found by their hash among the heads that were bound: a
// list whose hash no head holds is no list that can be applied.
func TestAFrozenListWhoseHashNoHeadHoldsIsRefused(t *testing.T) {
	l := newLoop(t, nil)
	l.frozenAt("claude")
	rows := l.tests()
	rows[0][2] = hashC
	data, err := WriteTests(rows)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.r.commit(context.Background(), map[string][]byte{"tasks/T-ab12/tests.tsv": data}, "a list that no head holds"); err != nil {
		t.Fatal(err)
	}
	if got, err := l.r.FrozenFiles(context.Background(), "T-ab12"); err == nil {
		t.Errorf("FrozenFiles = %v; want an error: no bound head holds the bytes of TestA", got)
	}
}

// checkout makes a clone-like repository with one commit, and a directory
// outside it, for the tests of the write of the frozen files.
func checkout(t *testing.T) (session.Dir, string, *Sessions) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	repo, outside := filepath.Join(t.TempDir(), "repo"), t.TempDir()
	if err := git.Init(repo); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repo, "README.md"), []byte("r\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := git.Add(repo); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(repo, "base", git.Identity{Name: "t", Email: "t@t.invalid", Time: time.Unix(1, 0)}); err != nil {
		t.Fatal(err)
	}
	return session.Dir{Repo: repo}, outside, &Sessions{Now: func() time.Time { return time.Unix(2, 0) }}
}

func untouched(t *testing.T, path, want string) {
	t.Helper()
	if got, err := os.ReadFile(path); err != nil || string(got) != want {
		t.Errorf("%s = %q, %v; want %q: the write of a frozen file left the checkout", path, got, err, want)
	}
}

// guardrails §2, a path that is clean by its spelling still leaves through a
// link: a parent of the frozen path, or the file itself, that is a link in the
// checkout is not followed; the write is refused, naming the path, and nothing
// outside the checkout is read or written.
func TestAFrozenFileThroughAParentLinkIsRefusedAndWritesNothingOutside(t *testing.T) {
	d, outside, s := checkout(t)
	if err := os.WriteFile(filepath.Join(outside, "a_test.go"), []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(d.Repo, "tests")); err != nil {
		t.Fatal(err)
	}
	err := s.applyFrozen(d, "S-aaaaaaaa", []FrozenFile{{Path: "tests/a_test.go", Data: []byte("frozen")}})
	if err == nil || !strings.Contains(err.Error(), "tests/a_test.go") {
		t.Errorf("applyFrozen through a link to a directory: %v; want an error naming tests/a_test.go", err)
	}
	untouched(t, filepath.Join(outside, "a_test.go"), "original")
}

func TestAFrozenFileThatIsALinkIsRefusedAndWritesNothingOutside(t *testing.T) {
	d, outside, s := checkout(t)
	target := filepath.Join(outside, "victim.go")
	if err := os.WriteFile(target, []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, filepath.Join(d.Repo, "a_test.go")); err != nil {
		t.Fatal(err)
	}
	err := s.applyFrozen(d, "S-aaaaaaaa", []FrozenFile{{Path: "a_test.go", Data: []byte("frozen")}})
	if err == nil || !strings.Contains(err.Error(), "a_test.go") {
		t.Errorf("applyFrozen over a link to a file: %v; want an error naming a_test.go", err)
	}
	untouched(t, target, "original")
}

// The rule is on the link, not on its target: a link that points inside the
// checkout is refused too, and the file it points to is not written.
func TestAFrozenFileThroughALinkInsideTheCheckoutIsRefused(t *testing.T) {
	d, _, s := checkout(t)
	if err := os.WriteFile(filepath.Join(d.Repo, "real.go"), []byte("original"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("real.go", filepath.Join(d.Repo, "a_test.go")); err != nil {
		t.Fatal(err)
	}
	if err := s.applyFrozen(d, "S-aaaaaaaa", []FrozenFile{{Path: "a_test.go", Data: []byte("frozen")}}); err == nil {
		t.Error("applyFrozen over a link that points inside the checkout: no error")
	}
	untouched(t, filepath.Join(d.Repo, "real.go"), "original")
}

// A refusal writes nothing for the files before it either would be wrong, so the
// check is made for the whole list first: a plain file earlier in the list is not
// written when a later one is refused.
func TestARefusedFrozenListWritesNoFile(t *testing.T) {
	d, outside, s := checkout(t)
	if err := os.Symlink(outside, filepath.Join(d.Repo, "tests")); err != nil {
		t.Fatal(err)
	}
	err := s.applyFrozen(d, "S-aaaaaaaa", []FrozenFile{{Path: "a_test.go", Data: []byte("one")}, {Path: "tests/b_test.go", Data: []byte("two")}})
	if err == nil {
		t.Fatal("no error")
	}
	if _, err := os.Lstat(filepath.Join(d.Repo, "a_test.go")); err == nil {
		t.Error("a_test.go was written before the later file was refused")
	}
	if entries, _ := os.ReadDir(outside); len(entries) != 0 {
		t.Errorf("%d entries outside the checkout", len(entries))
	}
}

// The control: a plain checkout gets the frozen bytes, a missing parent is made,
// an existing regular file is replaced, and the commit has the session's identity.
func TestAPlainCheckoutGetsTheFrozenFilesInOneCommit(t *testing.T) {
	d, _, s := checkout(t)
	if err := os.WriteFile(filepath.Join(d.Repo, "a_test.go"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	files := []FrozenFile{{Path: "a_test.go", Data: []byte("frozen a")}, {Path: "pkg/deep/b_test.go", Data: []byte("frozen b")}}
	if err := s.applyFrozen(d, "S-aaaaaaaa", files); err != nil {
		t.Fatal(err)
	}
	untouched(t, filepath.Join(d.Repo, "a_test.go"), "frozen a")
	untouched(t, filepath.Join(d.Repo, "pkg/deep/b_test.go"), "frozen b")
	head, err := git.Message(d.Repo, "HEAD")
	if err != nil || head != frozenMessage {
		t.Errorf("the message of HEAD = %q, %v", head, err)
	}
}
