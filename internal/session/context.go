package session

// This file holds the step table's column of record kinds (docs/spec/loop.md,
// The tier and the context of a task; docs/architecture.md §9, Roles and steps;
// row 63b of the plan, task T-svn3, #353): for each step, the kinds of record
// that go into its prompt file, in the order of the specification. The builder
// of the prompt file is row 66e's; this file only says what its list holds.

// Kind is a kind of record that goes into the prompt file of a step of a task.
type Kind string

const (
	// KindIssueAndPlan is the task's issue and plan.
	KindIssueAndPlan Kind = "issue-and-plan"
	// KindRequirements is the rows and the specification sections of the task's
	// requirement IDs.
	KindRequirements Kind = "requirements"
	// KindDecisionRecords is the decision records that those sections name.
	KindDecisionRecords Kind = "decision-records"
	// KindHandoffsAndAnswers is the task's handoffs and the answers to its questions.
	KindHandoffsAndAnswers Kind = "handoffs-and-answers"
	// KindRuleFiles is the target's rule files.
	KindRuleFiles Kind = "rule-files"
)

// contextKinds is the list of the specification, in its order.
var contextKinds = []Kind{KindIssueAndPlan, KindRequirements, KindDecisionRecords, KindHandoffsAndAnswers, KindRuleFiles}

// ContextKinds returns the kinds that the prompt file of a step of a task
// holds, in this order. The slice is a copy.
func ContextKinds() []Kind { return append([]Kind(nil), contextKinds...) }

// StepRow is a row of the step table: a step, its role and its tier as
// docs/architecture.md §9 gives them, and the record kinds of its prompt file.
// A step of a phase that has no task has no kinds: a session of Intake or of the
// Scaffold reads the input/ files of its own section (intake.md, scaffold.md),
// and the bet brief and the retrospective lessons have their own sections
// (loop.md, architecture.md §13). The values by-class and owner-of-its-kind are
// the table's own words, resolved before routing (session.md).
type StepRow struct {
	Step, Role, Tier string
	Kinds            []Kind
}

// defaultSteps is the default step table of §9, its grouped rows split into one row
// for each session step. A task step has the kinds; a phase step has none.
var defaultSteps = []struct {
	step, role, tier string
	task             bool
}{
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

// StepTable returns the step table, a copy: a caller's change reaches neither
// the rows nor their kinds.
func StepTable() []StepRow {
	out := make([]StepRow, len(defaultSteps))
	for i, s := range defaultSteps {
		out[i] = StepRow{Step: s.step, Role: s.role, Tier: s.tier}
		if s.task {
			out[i].Kinds = ContextKinds()
		}
	}
	return out
}
