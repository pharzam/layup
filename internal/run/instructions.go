package run

// This file holds the fixed instruction of each step of a task (docs/spec/loop.md,
// The instruction of a task step; task T-c6sx, #416, row 66s of the plan of M2e):
// the work, the outputs, the result file and the comment form that the parsers of
// the task loop accept. The text of record is the block of loop.md, and a test
// compares each constant with its block byte for byte, as it does for the probe.

// The names of the comment files of a task session, which row 66r's driver passes
// as the file of the handoff (loop.md, The instruction of a task step).
const (
	PlanFile         = "plan.md"
	PlanReviewFile   = "plan-review.md"
	ReviewRecordFile = "review-record.md"
)

// stepInstructions is the instruction of each task step that has one, by step.
var stepInstructions = map[string]string{}

// resultInstruction is the text shared by all the steps: the form of the result file.
const resultInstruction = ""

// A survivor is the options that the filter kept for one decision, in the
// order of its table.
type survivor struct {
	Decision string
	Options  []string
}

// taskExtras are the inputs of a prompt file that come with the task and not
// with the step table: the frozen list, the Shape records, the options with
// their survivors and the ledger.
type taskExtras struct {
	Frozen                 *kept
	Shape, Options, Ledger []kept
	Survivors              []survivor
	From                   string
}
