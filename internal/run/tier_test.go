package run

import (
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/session"
)

// The tier of a plan (loop.md, The tier and the context of a task): the
// reasoning tier when size is large or any of the three lists is not empty,
// the execution tier otherwise. One fixture for each way to reach reasoning,
// so that a rule dropped alone fails its own case.
func TestTheTierOfEachPlanClass(t *testing.T) {
	for _, c := range []struct {
		name string
		plan Plan
		want string
	}{
		{"a small plan with three empty lists", Plan{Size: "small"}, "execution"},
		{"large alone", Plan{Size: "large"}, "reasoning"},
		{"an open question alone", Plan{Size: "small", OpenQuestions: []string{"Q?"}}, "reasoning"},
		{"a new interface alone", Plan{Size: "small", NewInterfaces: []string{"internal/x"}}, "reasoning"},
		{"a new dependency alone", Plan{Size: "small", NewDependencies: []string{"a/b"}}, "reasoning"},
		{"a size that is neither small nor large", Plan{Size: "medium"}, "reasoning"},
		{"no size", Plan{}, "reasoning"},
	} {
		if got := PlanTier(c.plan); got != c.want {
			t.Errorf("%s: tier %q, want %q", c.name, got, c.want)
		}
	}
}

// The progress position moves the tier at an attempt: uphill to reasoning,
// downhill to execution (architecture.md §9); with no position the plan alone
// decides, as at its acceptance.
func TestTheProgressPositionMovesTheTier(t *testing.T) {
	small, large := Plan{Size: "small"}, Plan{Size: "large"}
	for _, c := range []struct {
		name string
		plan Plan
		pos  Position
		want string
	}{
		{"uphill moves a small plan to reasoning", small, Uphill, "reasoning"},
		{"downhill moves a large plan to execution", large, Downhill, "execution"},
		{"no position leaves a small plan at execution", small, "", "execution"},
		{"no position leaves a large plan at reasoning", large, "", "reasoning"},
	} {
		if got := TierOf(c.plan, c.pos); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

func eventsOf(lines string) [][]string { return taskEvents(strings.Split(lines, "\n")...) }

// The open unknowns that the records give in M2e are the questions of the last
// attempt that ended needs_context, until an answer to it is accepted.
func TestOpenQuestionsAreTheQuestionsNoAcceptedAnswerClosed(t *testing.T) {
	asked := map[string][][]string{askSess: resultOfRows(needsHuman, askedRow, askedRow)}
	pending := "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done"
	if n := OpenQuestions(eventsOf(pending), asked, stdRoles(), DefaultOwners()); n != 2 {
		t.Errorf("two questions of the last attempt: %d open, want 2", n)
	}
	if n := OpenQuestions(eventsOf("attempt 1\nsession 1 "+askSess+"\nresult 1 "+askSess+" done"), map[string][][]string{askSess: resultOfRows("completed\t—\tdone")}, stdRoles(), DefaultOwners()); n != 0 {
		t.Errorf("an attempt that ended completed: %d open, want 0", n)
	}
	if n := OpenQuestions(eventsOf(stdEvents), answered(), stdRoles(), DefaultOwners()); n != 0 {
		t.Errorf("an answer accepted by the next attempt: %d open, want 0", n)
	}
	// The next attempt cites the answer and asks it again: not accepted, so the question stays open.
	again := answered()
	again[nextSess] = resultOfRows(cited, "question\t1\tAgain?\t—\tinterface-contract no docs/spec/forge.md "+answerID)
	if n := OpenQuestions(eventsOf(stdEvents), again, stdRoles(), DefaultOwners()); n != 1 {
		t.Errorf("an answer that the next attempt asks again: %d open, want 1 (the questions of the last attempt, the answer not accepted)", n)
	}
	if n := OpenQuestions(nil, nil, nil, nil); n != 0 {
		t.Errorf("no event: %d open, want 0", n)
	}
}

// An answer of a session outside the owner role closes no question, and the
// owner role's does (loop.md, A question during the work, T-b3s3).
func TestAWrongRoleAnswerClosesNoOpenQuestion(t *testing.T) {
	roles := stdRoles()
	roles[ownerSess] = "software-developer"
	if n := OpenQuestions(eventsOf(stdEvents), answered(), roles, DefaultOwners()); n != 1 {
		t.Errorf("an answer of a developer: %d open, want 1", n)
	}
	if n := OpenQuestions(eventsOf(stdEvents), answered(), nil, DefaultOwners()); n != 1 {
		t.Errorf("an answer with no role known: %d open, want 1", n)
	}
	if n := OpenQuestions(eventsOf(stdEvents), answered(), stdRoles(), DefaultOwners()); n != 0 {
		t.Errorf("an answer of the owner role: %d open, want 0", n)
	}
}

// No row of the step table keeps the table's own word by-class once it is
// resolved with a tier, and none keeps owner-of-its-kind once it has an owner.
func TestEveryRowOfTheStepTableResolvesToATier(t *testing.T) {
	for _, row := range session.StepTable() {
		got, err := ResolveStep(row.Step, "execution", "Systems Architect")
		if err != nil {
			t.Errorf("%s: %v", row.Step, err)
			continue
		}
		if got.Tier != "reasoning" && got.Tier != "execution" {
			t.Errorf("%s: tier %q is neither reasoning nor execution", row.Step, got.Tier)
		}
		if row.Tier == "by-class" && got.Tier != "execution" {
			t.Errorf("%s: by-class with the tier execution gave %q", row.Step, got.Tier)
		}
		if row.Tier == "reasoning" && got.Tier != "reasoning" {
			t.Errorf("%s: a reasoning row gave %q", row.Step, got.Tier)
		}
		if got.Role == "owner-of-its-kind" || got.Role == "" {
			t.Errorf("%s: role %q is not resolved", row.Step, got.Role)
		}
	}
	if _, err := ResolveStep("implement", "", "Systems Architect"); err == nil {
		t.Error("by-class with no tier: no error")
	}
	if _, err := ResolveStep("implement", "by-class", ""); err == nil {
		t.Error("by-class with the tier by-class: no error")
	}
	if _, err := ResolveStep("no-such-step", "reasoning", ""); err == nil {
		t.Error("an unknown step: no error")
	}
}

// The owner map's role becomes a routing key by the one rule, and each default
// owner that the step table has as a role comes out as exactly that ID.
func TestOwnerOfItsKindIsResolvedByTheOwnerMap(t *testing.T) {
	var ids []string
	for _, row := range session.StepTable() {
		ids = append(ids, row.Role)
	}
	want := map[string]string{"Domain Expert": "domain-expert", "Systems Architect": "systems-architect", "Software Architect": "software-architect", "Software Engineer": "software-engineer"}
	for _, owner := range DefaultOwners() {
		got, err := ResolveStep("answer-to-a-question", "", owner)
		if err != nil || got.Role != want[owner] {
			t.Errorf("%s: role %q, %v; want %q", owner, got.Role, err, want[owner])
		}
		if owner != "Domain Expert" && !slices.Contains(ids, got.Role) {
			t.Errorf("%s: %q is no role of the step table %v", owner, got.Role, ids)
		}
	}
	if got, err := ResolveStep("answer-to-a-question", "", "Quality Lead"); err != nil || got.Role != "quality-lead" {
		t.Errorf("a replaced map's role: %q, %v", got.Role, err)
	}
	if _, err := ResolveStep("answer-to-a-question", "", "  "); err == nil {
		t.Error("an owner with no name: no error")
	}
	if _, err := ResolveStep("answer-to-a-question", "", "operator"); err == nil {
		t.Error("the Operator is no owner role: no error")
	}
}

func TestTheRoutingRoleOfAnOwnerIsTheRoleIDOfTheStepTable(t *testing.T) {
	for in, want := range map[string]string{"Systems Architect": "systems-architect", "  Software   Engineer ": "software-engineer", "Domain Expert": "domain-expert", "": ""} {
		if got := RoutingRole(in); got != want {
			t.Errorf("RoutingRole(%q) = %q, want %q", in, got, want)
		}
	}
}
