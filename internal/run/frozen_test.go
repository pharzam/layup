package run

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/tsv"
)

const testsHead = "test\tsource\tsha256\tharness\tfrozen\n"

var (
	hashA = strings.Repeat("a", 64)
	hashB = strings.Repeat("b", 64)
	hashC = strings.Repeat("c", 64)
)

func frozenList(rows ...string) []byte { return []byte(testsHead + strings.Join(rows, "")) }

func goodRow(id, src, sum, harness string) string {
	return id + "\t" + src + "\t" + sum + "\t" + harness + "\t2026-10-16T12:00:00Z\n"
}

func TestAFrozenListIsReadByItsRules(t *testing.T) {
	data := frozenList(goodRow("TestA", "a_test.go", hashA, "claude"), goodRow("TestB", "pkg/b_test.go", hashB, "codex"))
	rows, err := ReadTests(data)
	if err != nil || len(rows) != 2 || rows[1][1] != "pkg/b_test.go" {
		t.Fatalf("ReadTests = %v, %v", rows, err)
	}
	back, err := WriteTests(rows)
	if err != nil || string(back) != string(data) {
		t.Errorf("WriteTests = %q, %v; want the bytes it read", back, err)
	}
}

// Each rule of the owner has a fixture of its own: the one broken rule is the
// only reason that can refuse the list (guardrails §2).
func TestAFrozenListRefusesEachBrokenRule(t *testing.T) {
	for _, c := range []struct {
		name, column string
		data         []byte
	}{
		{"a header and no row", "", []byte(testsHead)},
		{"a source of —", "source", frozenList("TestA\t—\t" + hashA + "\tclaude\t2026-10-16T12:00:00Z\n")},
		{"a sha256 of —", "sha256", frozenList("TestA\ta_test.go\t—\tclaude\t2026-10-16T12:00:00Z\n")},
		{"a harness of —", "harness", frozenList("TestA\ta_test.go\t" + hashA + "\t—\t2026-10-16T12:00:00Z\n")},
		{"a frozen time that differs between rows", "frozen", frozenList(goodRow("TestA", "a_test.go", hashA, "claude"), "TestB\tb_test.go\t"+hashB+"\tclaude\t2026-10-16T12:00:01Z\n")},
		{"a frozen time of —", "frozen", frozenList("TestA\ta_test.go\t" + hashA + "\tclaude\t—\n")},
		{"a repeated test ID", "", frozenList(goodRow("TestA", "a_test.go", hashA, "claude"), goodRow("TestA", "b_test.go", hashB, "claude"))},
		{"a bad hash", "sha256", frozenList("TestA\ta_test.go\tnothex\tclaude\t2026-10-16T12:00:00Z\n")},
		{"a source with a .. part", "source", frozenList("TestA\t../a_test.go\t" + hashA + "\tclaude\t2026-10-16T12:00:00Z\n")},
	} {
		_, err := ReadTests(c.data)
		var te *tsv.Error
		if err == nil || !errors.As(err, &te) {
			t.Errorf("%s: %v; want a *tsv.Error", c.name, err)
			continue
		}
		if c.column != "" && te.Column != c.column {
			t.Errorf("%s: refused at the column %q, want %q (%v)", c.name, te.Column, c.column, err)
		}
	}
	if _, err := WriteTests([][]string{{"TestA", "a_test.go", hashA, "claude", "2026-10-16T12:00:00Z"}, {"TestB", "b_test.go", hashB, "claude", "2026-10-16T12:00:09Z"}}); err == nil {
		t.Error("WriteTests wrote a list that ReadTests refuses")
	}
}

// loop.md, the frozen test list: each listed test fails at the base and passes
// at the head, and nothing proven is no proof.
func TestAnEmptyListFreezesNothingAndIsInvalid(t *testing.T) {
	if redGreen(nil) || redGreen([]TestResult{}) {
		t.Error("an empty list is red-green; \"each row is true\" of no row must not pass")
	}
}

func TestAListedTestThatDoesNotFailAtTheBaseMakesTheHandoffInvalid(t *testing.T) {
	ok := TestResult{ID: "A", Source: "a_test.go", FailsAtBase: true, PassesAtHead: true}
	if !redGreen([]TestResult{ok}) {
		t.Fatal("a test that fails at the base and passes at the head is not red-green")
	}
	for _, c := range []struct {
		name string
		bad  TestResult
	}{
		{"passes at the base", TestResult{ID: "B", Source: "b_test.go", FailsAtBase: false, PassesAtHead: true}},
		{"fails at the head", TestResult{ID: "B", Source: "b_test.go", FailsAtBase: true, PassesAtHead: false}},
		{"was not run", TestResult{ID: "B", Source: "b_test.go"}},
	} {
		if redGreen([]TestResult{ok, c.bad}) {
			t.Errorf("a list with a test that %s is red-green", c.name)
		}
	}
	// A result with no ID or no source cannot be listed.
	if redGreen([]TestResult{{Source: "a_test.go", FailsAtBase: true, PassesAtHead: true}}) || redGreen([]TestResult{{ID: "A", FailsAtBase: true, PassesAtHead: true}}) {
		t.Error("a test with no ID or no source is red-green")
	}
}

func frozenRows() [][]string {
	return [][]string{
		{"TestA", "a_test.go", hashA, "claude", "2026-10-16T12:00:00Z"},
		{"TestB", "b_test.go", hashB, "claude", "2026-10-16T12:00:00Z"},
	}
}

func TestAFreezingWritesTheRowsWithTheSessionsHarnessAndOneTime(t *testing.T) {
	at := time.Date(2026, 10, 16, 13, 0, 0, 0, time.UTC)
	tests := []TestResult{{ID: "TestA", Source: "a_test.go"}, {ID: "TestB", Source: "b_test.go"}}
	rows := freezeRows(nil, tests, map[string]string{"a_test.go": hashA, "b_test.go": hashB}, "codex", at)
	want := [][]string{{"TestA", "a_test.go", hashA, "codex", "2026-10-16T13:00:00Z"}, {"TestB", "b_test.go", hashB, "codex", "2026-10-16T13:00:00Z"}}
	if len(rows) != 2 || strings.Join(rows[0], "|") != strings.Join(want[0], "|") || strings.Join(rows[1], "|") != strings.Join(want[1], "|") {
		t.Errorf("freezeRows = %v, want %v", rows, want)
	}
}

// At a re-freeze the harness of a test whose source did not change is kept; a
// changed or a new test takes the harness of the freezing session; the list is
// written whole.
func TestTheHarnessOfAnUnchangedTestIsKept(t *testing.T) {
	at := time.Date(2026, 10, 16, 14, 0, 0, 0, time.UTC)
	tests := []TestResult{{ID: "TestA", Source: "a_test.go"}, {ID: "TestB", Source: "b_test.go"}, {ID: "TestN", Source: "n_test.go"}}
	rows := freezeRows(frozenRows(), tests, map[string]string{"a_test.go": hashA, "b_test.go": hashC, "n_test.go": hashA}, "codex", at)
	if len(rows) != 3 {
		t.Fatalf("freezeRows = %v; want the three tests of the new plan", rows)
	}
	for i, h := range []string{"claude", "codex", "codex"} {
		if rows[i][3] != h {
			t.Errorf("row %d (%s): harness %q, want %q", i, rows[i][0], rows[i][3], h)
		}
		if rows[i][4] != "2026-10-16T14:00:00Z" {
			t.Errorf("row %d: frozen %q, want the one time of the freezing", i, rows[i][4])
		}
	}
	if rows[1][2] != hashC {
		t.Errorf("the changed test keeps the hash %q, want the new %q", rows[1][2], hashC)
	}
	// A test that kept its hash at another path is a new test: its harness is the session's.
	moved := freezeRows(frozenRows(), []TestResult{{ID: "TestA", Source: "moved_test.go"}}, map[string]string{"moved_test.go": hashA}, "codex", at)
	if len(moved) != 1 || moved[0][3] != "codex" {
		t.Errorf("a test at another path: %v, want the harness of the freezing session", moved)
	}
}

func TestAListThatIsNotTheFrozenOneIsNotTheSame(t *testing.T) {
	same := []TestResult{{ID: "TestA", Source: "a_test.go"}, {ID: "TestB", Source: "b_test.go"}}
	hashes := map[string]string{"a_test.go": hashA, "b_test.go": hashB}
	if !sameAsFrozen(frozenRows(), same, hashes) {
		t.Fatal("the frozen tests at their hashes are not the frozen list")
	}
	for _, c := range []struct {
		name   string
		tests  []TestResult
		hashes map[string]string
	}{
		{"a source that changed", same, map[string]string{"a_test.go": hashA, "b_test.go": hashC}},
		{"a source that is missing", same, map[string]string{"a_test.go": hashA}},
		{"a test added", append(append([]TestResult{}, same...), TestResult{ID: "TestN", Source: "n_test.go"}), map[string]string{"a_test.go": hashA, "b_test.go": hashB, "n_test.go": hashA}},
		{"a test missing", same[:1], hashes},
		{"a test at another path", []TestResult{same[0], {ID: "TestB", Source: "c_test.go"}}, map[string]string{"a_test.go": hashA, "c_test.go": hashB}},
	} {
		if sameAsFrozen(frozenRows(), c.tests, c.hashes) {
			t.Errorf("%s: the list is the frozen one", c.name)
		}
	}
}

func TestAChangedSourceIsTheFindingTestChanged(t *testing.T) {
	cur := map[string]string{"a_test.go": hashA, "b_test.go": hashC}
	got, err := changedTests(frozenRows(), func(p string) (string, error) { return cur[p], nil })
	if err != nil || len(got) != 1 {
		t.Fatalf("changedTests = %v, %v; want one finding", got, err)
	}
	if f := got[0]; f.Rule != "test changed" || f.File != "b_test.go" || f.Text != "TestB" {
		t.Errorf("the finding %+v", f)
	}
	// A source that the head lacks is a change too.
	delete(cur, "a_test.go")
	if got, _ := changedTests(frozenRows(), func(p string) (string, error) { return cur[p], nil }); len(got) != 2 {
		t.Errorf("a missing source: %v, want a finding for each of the two", got)
	}
	if got, _ := changedTests(frozenRows(), func(p string) (string, error) {
		return map[string]string{"a_test.go": hashA, "b_test.go": hashB}[p], nil
	}); len(got) != 0 {
		t.Errorf("nothing changed: %v", got)
	}
	if _, err := changedTests(frozenRows(), func(string) (string, error) { return "", errors.New("git failed") }); err == nil {
		t.Error("a read that fails is no finding and no pass")
	}
}

// roleOf is the role of a step from the production step table.
func roleOf(t *testing.T, step string) string {
	t.Helper()
	for _, r := range session.StepTable() {
		if r.Step == step {
			return r.Role
		}
	}
	t.Fatalf("the step table has no step %s", step)
	return ""
}

// amend builds the sessions.tsv rows and the events of the amendment fixtures:
// each line is "kind session detail", the session of a line is the id.
type amend struct {
	t        *testing.T
	sessions [][]string
	events   [][]string
}

func (a *amend) session(id, step string) {
	a.sessions = append(a.sessions, []string{id, "T-ab12", "1", roleOf(a.t, step), "claude", "2.1", "m1", aSHA, aSHA, "9", "3", "1000", "5.0", "30", "", ""})
	a.add("session", id, "")
}

func (a *amend) add(kind, id, detail string) {
	a.events = append(a.events, []string{strconv.Itoa(len(a.events) + 1), kind, "1", id, "", "", detail, "2026-10-16T12:00:00Z"})
}

// freezeAt: the developer's valid handoff, after which an amendment is looked for.
func (a *amend) freezeAt(id string) {
	a.session(id, "implement")
	a.add("result", id, "done")
	a.add("handoff", id, "valid")
}

func (a *amend) pending() bool { return amendmentPending(a.events, a.sessions) }

func newAmend(t *testing.T) *amend {
	a := &amend{t: t}
	a.add("attempt", "", "")
	a.freezeAt("S-d0000001")
	return a
}

func (a *amend) plan(id, detail string) {
	a.session(id, "plan-of-a-task")
	a.add("result", id, "done")
	a.add("handoff", id, detail)
}

func (a *amend) review(id, detail string) {
	a.session(id, "plan-review")
	a.add("result", id, "done")
	a.add("handoff", id, detail)
}

func TestAnAmendmentIsANewPlanThenANewPlanReview(t *testing.T) {
	a := newAmend(t)
	if a.pending() {
		t.Fatal("a freezing with nothing after it is an amendment")
	}
	a.plan("S-a0000001", "valid comment 701-1")
	if a.pending() {
		t.Error("a new plan with no review is an amendment")
	}
	a.review("S-a0000002", "valid comment 702-1")
	if !a.pending() {
		t.Error("a new plan and its review after the freezing are no amendment")
	}
	// The next valid handoff of a developer consumes it.
	a.freezeAt("S-d0000002")
	if a.pending() {
		t.Error("the amendment is still pending after the freezing it allowed")
	}
}

func TestAnAmendmentIsNotAReviewWithoutAPlanOrInTheWrongOrder(t *testing.T) {
	only := newAmend(t)
	only.review("S-a0000002", "valid comment 702-1")
	if only.pending() {
		t.Error("a review with no plan is an amendment")
	}
	wrong := newAmend(t)
	wrong.review("S-a0000002", "valid comment 702-1")
	wrong.plan("S-a0000001", "valid comment 701-1")
	if wrong.pending() {
		t.Error("a review before the plan is an amendment")
	}
	// Two reviews and no plan, and two plans and no review, are no amendment: each leg needs its own role.
	twoReviews := newAmend(t)
	twoReviews.review("S-a0000002", "valid comment 702-1")
	twoReviews.review("S-a0000003", "valid comment 703-1")
	if twoReviews.pending() {
		t.Error("two reviews and no plan are an amendment")
	}
	twoPlans := newAmend(t)
	twoPlans.plan("S-a0000001", "valid comment 701-1")
	twoPlans.plan("S-a0000004", "valid comment 704-1")
	if twoPlans.pending() {
		t.Error("two plans and no review are an amendment")
	}
	// A plan or a review that is not valid, or that has no comment, is none.
	for name, plan := range map[string]string{"an invalid plan": "invalid plan-fields comment 701-1", "a plan with no comment": "valid"} {
		a := newAmend(t)
		a.plan("S-a0000001", plan)
		a.review("S-a0000002", "valid comment 702-1")
		if a.pending() {
			t.Errorf("%s: an amendment", name)
		}
	}
	for name, review := range map[string]string{"an invalid review": "invalid verdict-field comment 702-1", "a review with no comment": "valid"} {
		a := newAmend(t)
		a.plan("S-a0000001", "valid comment 701-1")
		a.review("S-a0000002", review)
		if a.pending() {
			t.Errorf("%s: an amendment", name)
		}
	}
}

// A verifier shares the role of the plan reviewer and gives a comment handoff
// too: its record, of a session that started before the new plan, is no review
// of it even when its handoff comes after the plan (the plan review of #361).
func TestAStaleVerifiersRecordIsNoReviewOfTheNewPlan(t *testing.T) {
	a := newAmend(t)
	a.session("S-v0000001", "verification") // starts before the plan
	a.plan("S-a0000001", "valid comment 701-1")
	a.add("result", "S-v0000001", "done")
	a.add("handoff", "S-v0000001", "valid comment 703-1") // its record comes after the plan
	if a.pending() {
		t.Error("a verifier's review record that started before the new plan is its review")
	}
	// A verifier that starts after the plan is a session of the plan-review role too: it reviews the new plan.
	a.review("S-a0000002", "valid comment 702-1")
	if !a.pending() {
		t.Error("the plan review that started after the new plan is not found behind the stale verifier")
	}
}

// A session of the plan's role that answers a question gives a handoff with no
// comment, so it is no plan.
func TestAnAnswerOfThePlanRoleIsNoNewPlan(t *testing.T) {
	a := newAmend(t)
	a.plan("S-a0000001", "valid")
	a.review("S-a0000002", "valid comment 702-1")
	if a.pending() {
		t.Error("an answer's handoff of the plan role is a plan")
	}
}

// loop.md, the frozen test list: the list freezes at a valid handoff to the
// verifier and at no other: not at a handoff to another role, not at an invalid
// one, and not when the caller gave no list.
func TestTheListFreezesOnlyAtAValidHandoffToTheVerifierThatGaveOne(t *testing.T) {
	list := &Freeze{Tests: []TestResult{{ID: "A", Source: "a_test.go", FailsAtBase: true, PassesAtHead: true}}}
	valid := handoff.Verdict{Valid: true, Detail: "valid"}
	if !shouldFreeze(HandoffInput{To: "verifier", Freeze: list}, valid) {
		t.Fatal("a valid handoff to the verifier with a list does not freeze")
	}
	for name, c := range map[string]struct {
		in HandoffInput
		v  handoff.Verdict
	}{
		"a handoff to the developer": {HandoffInput{To: "developer", Freeze: list}, valid},
		"a handoff to the close-out": {HandoffInput{To: "close-out", Freeze: list}, valid},
		"an invalid handoff":         {HandoffInput{To: "verifier", Freeze: list}, handoff.Verdict{Detail: "invalid tests-red-green"}},
		"no list":                    {HandoffInput{To: "verifier"}, valid},
	} {
		if shouldFreeze(c.in, c.v) {
			t.Errorf("%s freezes the list", name)
		}
	}
}

// A list given to a handoff that is not to the verifier is ignored whole: it
// asks for no head, so a session with none is no error there.
func TestAListGivenToAnotherHandoffIsIgnored(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, roleOfStep("implement"), "", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n")
	in := HandoffInput{Task: "T-ab12", Session: planSess, From: "asking role", To: "owner role",
		Freeze: &Freeze{Tests: []TestResult{{ID: "A", Source: "a_test.go"}}}}
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "valid" {
		t.Errorf("RecordHandoff = %+v, %v; want valid: the list of a handoff to the owner role is not read", v, err)
	}
	if _, ok := w.store.files[testsPath("T-ab12")]; ok {
		t.Error("a list was written")
	}
}
