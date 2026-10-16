package session

import (
	"slices"
	"testing"
)

// The expected step table, written apart from the code: docs/architecture.md
// §9 (the default step table, its grouped rows split) and docs/spec/loop.md
// (The tier and the context of a task). A task step has the five kinds, in the
// order of loop.md; a phase step with no task has none.
var wantKinds = []Kind{"issue-and-plan", "requirements", "decision-records", "handoffs-and-answers", "rule-files"}

type wantRow struct {
	step, role, tier string
	task             bool
}

var wantRows = []wantRow{
	{"gap-review", "systems-architect", "reasoning", false},
	{"completeness-review", "qa-engineer", "reasoning", false},
	{"numbering-and-draft", "product-owner", "reasoning", false},
	{"specification", "product-owner", "reasoning", true},
	{"bet-brief", "product-owner", "reasoning", false},
	{"shape-architecture", "systems-architect", "reasoning", true},
	{"shape-activation-batch", "software-architect", "reasoning", true},
	{"plan-of-a-milestone", "software-engineer", "reasoning", true},
	{"plan-of-a-task", "software-engineer", "reasoning", true},
	{"implement", "software-developer", "by-class", true},
	{"close-out", "software-developer", "by-class", true},
	{"plan-review", "qa-engineer", "reasoning", true},
	{"verification", "qa-engineer", "reasoning", true},
	{"answer-to-a-question", "owner-of-its-kind", "reasoning", true},
	{"retrospective-lessons", "systems-architect", "reasoning", false},
}

func TestContextKindsAreInTheOrderOfTheSpecification(t *testing.T) {
	if got := ContextKinds(); !slices.Equal(got, wantKinds) {
		t.Errorf("ContextKinds() = %v, want %v", got, wantKinds)
	}
}

func TestTheStepTableHasEachStepWithItsRoleTierAndKinds(t *testing.T) {
	table := StepTable()
	if len(table) != len(wantRows) {
		t.Fatalf("%d rows, want %d", len(table), len(wantRows))
	}
	for i, w := range wantRows {
		r := table[i]
		if r.Step != w.step || r.Role != w.role || r.Tier != w.tier {
			t.Errorf("row %d = %q %q %q, want %q %q %q", i, r.Step, r.Role, r.Tier, w.step, w.role, w.tier)
		}
		want := []Kind(nil)
		if w.task {
			want = wantKinds
		}
		if !slices.Equal(r.Kinds, want) {
			t.Errorf("step %q: kinds %v, want %v", w.step, r.Kinds, want)
		}
	}
}

func TestNoTwoStepsShareAName(t *testing.T) {
	seen := map[string]bool{}
	for _, r := range StepTable() {
		if seen[r.Step] {
			t.Errorf("the step %q is twice in the table", r.Step)
		}
		seen[r.Step] = true
	}
}

func TestACallersChangeDoesNotReachTheTableNorTheKinds(t *testing.T) {
	table := StepTable()
	table[0].Role = "changed"
	table[3].Kinds[0] = "changed"
	table = append(table[:1], table[2:]...)
	kinds := ContextKinds()
	kinds[0] = "changed"
	again := StepTable()
	if len(again) != len(wantRows) || again[0].Role != "systems-architect" || again[3].Kinds[0] != wantKinds[0] || ContextKinds()[0] != wantKinds[0] {
		t.Errorf("a change of a returned value reached the table: %v %v", again[0], again[3])
	}
	_ = table
}
