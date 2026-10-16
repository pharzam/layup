package run

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/session"
)

func sessionRow(id, task, step, harness string) []string {
	return []string{id, task, "1", roleOfStep(step), harness, "2.1", "m1", aSHA, aSHA, "9", "3", "1000", "5.0", "30", "", ""}
}

// roleOfStep reads the role of a step of the production step table.
func roleOfStep(step string) string {
	for _, r := range session.StepTable() {
		if r.Step == step {
			return r.Role
		}
	}
	panic("the step table has no step " + step)
}

func boundEvent(n int, session, sha string) []string {
	return []string{strconv.Itoa(n), "bound", "1", session, "", sha, "task/T-ab12/1", "2026-10-16T12:00:00Z"}
}

func inDiffOf(shas ...string) func(string) (bool, error) {
	return func(sha string) (bool, error) { return slices.Contains(shas, sha), nil }
}

const (
	shaOne = "1111111111111111111111111111111111111111"
	shaTwo = "2222222222222222222222222222222222222222"
	shaOld = "3333333333333333333333333333333333333333"
)

// loop.md, The independence of the reviewers: the authors for a plan review are
// the harnesses of the task's plan sessions (sessions.tsv).
func TestTheAuthorsOfAPlanReviewAreTheHarnessesOfThePlanSessions(t *testing.T) {
	sessions := [][]string{
		sessionRow("S-a0000001", "T-ab12", "plan-of-a-task", "claude"),
		sessionRow("S-a0000002", "T-ab12", "plan-of-a-task", "devin"),
		sessionRow("S-a0000003", "T-ab12", "plan-of-a-task", "claude"),
		sessionRow("S-a0000004", "T-ab12", "implement", "codex"),
		sessionRow("S-a0000005", "T-ab12", "plan-review", "opencode"),
		sessionRow("S-a0000006", "T-zz99", "plan-of-a-task", "grok"),
	}
	got, err := authorsFrom(PlanReview, "T-ab12", sessions, nil, nil, nil)
	if err != nil || !slices.Equal(got, []string{"claude", "devin"}) {
		t.Errorf("the authors of a plan review = %v, %v; want [claude devin]: the plan sessions of the task, each once, in order", got, err)
	}
	// A session of the same role that answered an environment question is counted too: the limit of the record, on the safe side.
	sessions = append(sessions, sessionRow("S-a0000007", "T-ab12", "plan-of-a-task", "codex"))
	if got, _ := authorsFrom(PlanReview, "T-ab12", sessions, nil, nil, nil); !slices.Equal(got, []string{"claude", "codex", "devin"}) {
		t.Errorf("a session of the plan's role that answered: %v", got)
	}
}

func TestTheAuthorsOfAVerificationAreTheHarnessesBoundToACommitOfTheDiff(t *testing.T) {
	sessions := [][]string{
		sessionRow("S-d0000001", "T-ab12", "implement", "claude"),
		sessionRow("S-d0000002", "T-ab12", "implement", "devin"),
		sessionRow("S-d0000003", "T-ab12", "implement", "codex"),
		sessionRow("S-d0000004", "T-ab12", "plan-of-a-task", "opencode"),
	}
	events := [][]string{boundEvent(1, "S-d0000001", shaOne), boundEvent(2, "S-d0000002", shaTwo), boundEvent(3, "S-d0000003", shaOld)}
	got, err := authorsFrom(Verification, "T-ab12", sessions, events, nil, inDiffOf(shaOne, shaTwo))
	if err != nil || !slices.Equal(got, []string{"claude", "devin"}) {
		t.Errorf("the authors of a verification = %v, %v; want [claude devin]", got, err)
	}
	// The plan session is no author of a verification: only a commit of the diff makes one.
	if slices.Contains(got, "opencode") {
		t.Error("a plan session's harness is an author of a verification")
	}
}

// A commit that the diff from the base does not hold (an earlier attempt's, or
// a commit of the base) makes no author, and a push with no bound is no commit
// bound to a session.
func TestASessionBoundToACommitOutsideTheDiffIsNoAuthor(t *testing.T) {
	sessions := [][]string{sessionRow("S-d0000001", "T-ab12", "implement", "claude"), sessionRow("S-d0000002", "T-ab12", "implement", "devin")}
	events := [][]string{boundEvent(1, "S-d0000001", shaOld), {"2", "push", "1", "S-d0000002", "", shaOne, "task/T-ab12/1", "2026-10-16T12:00:00Z"}}
	got, err := authorsFrom(Verification, "T-ab12", sessions, events, nil, inDiffOf(shaOne))
	if err != nil || len(got) != 0 {
		t.Errorf("authors = %v, %v; want none: one session is bound outside the diff, the other pushed and was never bound", got, err)
	}
}

// Of another task's session, a bound commit makes no author of this task's change.
func TestASessionOfAnotherTaskIsNoAuthor(t *testing.T) {
	sessions := [][]string{sessionRow("S-d0000001", "T-zz99", "implement", "claude")}
	got, _ := authorsFrom(Verification, "T-ab12", sessions, [][]string{boundEvent(1, "S-d0000001", shaOne)}, nil, inDiffOf(shaOne))
	if len(got) != 0 {
		t.Errorf("authors = %v; want none", got)
	}
}

func TestTheWriterOfTheFrozenTestsIsAnAuthor(t *testing.T) {
	tests := [][]string{{"TestA", "a_test.go", hashA, "codex", "2026-10-16T12:00:00Z"}, {"TestB", "b_test.go", hashB, "claude", "2026-10-16T12:00:00Z"}}
	got, err := authorsFrom(Verification, "T-ab12", nil, nil, tests, inDiffOf())
	if err != nil || !slices.Equal(got, []string{"claude", "codex"}) {
		t.Errorf("the authors = %v, %v; want the writers of the frozen tests, [claude codex]", got, err)
	}
}

func TestAnAuthorsReadThatFailsIsNoAnswer(t *testing.T) {
	sessions := [][]string{sessionRow("S-d0000001", "T-ab12", "implement", "claude")}
	boom := errors.New("git failed")
	if got, err := authorsFrom(Verification, "T-ab12", sessions, [][]string{boundEvent(1, "S-d0000001", shaOne)}, nil, func(string) (bool, error) { return false, boom }); !errors.Is(err, boom) || got != nil {
		t.Errorf("authors = %v, %v; want the error", got, err)
	}
	// A bound event of a session that the records do not have is an error, never a skipped author.
	if _, err := authorsFrom(Verification, "T-ab12", nil, [][]string{boundEvent(1, "S-d0000009", shaOne)}, nil, inDiffOf(shaOne)); err == nil {
		t.Error("a bound event of a session with no start row: no error")
	}
}

// pairWorld is the records of a routing register and the probes of three
// harnesses, on the stand-in store.
func pairWorld(t *testing.T, routing string) (*Sessions, [][]string, [][]string) {
	t.Helper()
	store := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	var probes [][]string
	var harnesses, models [][]string
	for i, h := range []string{"claude", "devin", "codex"} {
		probes = append(probes, []string{"S-b000000" + strconv.Itoa(i+1), h, "1.0", "m1", "passed", "—", "AGENTS.md", "m1", "2026-10-16T12:00:00Z"})
		harnesses = append(harnesses, []string{h, "", "1", "/bin/true", "stdin", "/bin/true", "", "", "AGENTS.md", "", "none", "subscription", ""})
		models = append(models, []string{h, "m1", "1000", "https://example.invalid", "2026-10-16T12:00:00Z", "yes", "—"})
	}
	data, err := table(records.HarnessesSchema, probes)
	if err != nil {
		t.Fatal(err)
	}
	store.files["harnesses.tsv"] = data
	store.files["routing.tsv"] = []byte("role\ttier\tposition\tharness\tmodel\n" + routing)
	return newSessions(store), harnesses, models
}

func TestThePairIsTheFirstAdmittedOneOutsideTheAuthors(t *testing.T) {
	r, harnesses, models := pairWorld(t, "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\nqa-engineer\treasoning\t3\tcodex\tm1\nqa-engineer\texecution\t1\tcodex\tm1\n")
	h, m, err := r.PairFor(context.Background(), "qa-engineer", "reasoning", harnesses, models, []string{"claude", "devin"})
	if err != nil || idOf(h) != "codex" || len(m) < 2 || m[1] != "m1" {
		t.Errorf("PairFor = %v, %v, %v; want codex (the first outside claude and devin)", h, m, err)
	}
	h, _, err = r.PairFor(context.Background(), "qa-engineer", "reasoning", harnesses, models, []string{"claude"})
	if err != nil || idOf(h) != "devin" {
		t.Errorf("PairFor outside claude = %v, %v; want devin, the first outside it, not the last", h, err)
	}
	h, _, err = r.PairFor(context.Background(), "qa-engineer", "reasoning", harnesses, models, nil)
	if err != nil || idOf(h) != "claude" {
		t.Errorf("PairFor with no author = %v, %v; want claude", h, err)
	}
	// The tier is the call's, not a constant: the execution list has its own pair.
	h, _, err = r.PairFor(context.Background(), "qa-engineer", "execution", harnesses, models, []string{"claude"})
	if err != nil || idOf(h) != "codex" {
		t.Errorf("PairFor of the execution tier = %v, %v; want codex", h, err)
	}
}

func TestNoHarnessOutsideTheAuthorsIsTheErrorOfThePair(t *testing.T) {
	r, harnesses, models := pairWorld(t, "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\n")
	_, _, err := r.PairFor(context.Background(), "qa-engineer", "reasoning", harnesses, models, []string{"claude", "devin"})
	if !errors.Is(err, route.ErrNoPair) {
		t.Errorf("PairFor with every pair an author: %v; want route.ErrNoPair (the one refusal of the start)", err)
	}
}

func TestAHarnessThatFailedItsProbeAndAModelOfUseNoAreNotAdmitted(t *testing.T) {
	r, harnesses, models := pairWorld(t, "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\nqa-engineer\treasoning\t3\tcodex\tm1\n")
	probes := r.Store.(*standinStore).files["harnesses.tsv"]
	rows, rerr := records.ReadHarnesses(probes)
	if rerr != nil {
		t.Fatal(rerr)
	}
	rows[0][4] = "failed" // claude failed its probe
	rows[0][5] = "token"
	data, err := table(records.HarnessesSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	r.Store.(*standinStore).files["harnesses.tsv"] = data
	models[1][5] = "no" // devin's model is on the "not used" list
	h, _, err := r.PairFor(context.Background(), "qa-engineer", "reasoning", harnesses, models, nil)
	if err != nil || idOf(h) != "codex" {
		t.Errorf("PairFor = %v, %v; want codex: claude failed its probe and devin's model is not used", h, err)
	}
}

// idOf is the harness ID of a harness register row, "" for none.
func idOf(row []string) string {
	if len(row) == 0 {
		return ""
	}
	return row[0]
}
