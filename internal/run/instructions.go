package run

// This file holds the fixed instruction of each step of a task (docs/spec/loop.md,
// The instruction of a task step; task T-c6sx, #416, row 66s of the plan of M2e):
// the work, the outputs, the result file and the comment form that the parsers of
// the task loop accept. The text of record is the block of loop.md, and a test
// compares each constant with its block byte for byte, as it does for the probe.

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/spec"
)

// The names of the comment files of a task session, which row 66r's driver passes
// as the file of the handoff (loop.md, The instruction of a task step).
const (
	PlanFile         = "plan.md"
	PlanReviewFile   = "plan-review.md"
	ReviewRecordFile = "review-record.md"
)

// stepInstructions is the instruction of each task step that has one, by step: the
// blocks instruction-<step> of loop.md.
var stepInstructions = map[string]string{
	"plan-of-a-task": `Plan the task {task} (kind {kind}, issue {issue}) from the files listed above: the issue row, the requirements, the decision records, the handoffs and findings, and the rule files. Change no file of the repository and make no commit.
Write the plan as the file {comment}, in this form, each line starting at the first column and ended by a line feed:
size: <small or large>
open questions: <one line, or —>
new interfaces: <one line, or —>
new dependencies: <one line, or —>
tests:
- <id> <path>
The word after size: is exactly small or large. The character — is for the three list fields only, and a field with nothing after the colon and the space is wrong. The line tests: stands alone, once. Under it, write one line "- <id> <path>" for each test of the plan: exactly two words, each ID unique in the list, each path relative to the repository, clean and inside it (no space, no .., no leading /). A task of kind specification or architecture lists no test: write no line under tests:. A task of kind build lists at least one test.
In the result file write the row "comment<TAB>1<TAB>result/plan.md<TAB><the SHA-256 of {comment}><TAB>—".
`,
	"plan-review": `Review the plan of the task {task} (kind {kind}, issue {issue}), the newest plan among the handoffs listed above, against the issue row, the requirements, the decision records and the rule files: its coverage of the requirements, its list of tests, and its fit with the decision records and the rule files. Change no file of the repository and make no commit.
Write the review as the file {comment}, in this form, the first line starting at the first column and ended by a line feed:
verdict: <approve, approve-with-conditions or reject>
The line verdict: is written once, and its value is exactly one of the three words approve, approve-with-conditions and reject. After it, write the conditions or the reasons as plain lines that do not start with "verdict: ". This is the form the target reads; it is not the table "| Field | Value |" of LAYUP's own reviews, which you must not write.
In the result file write the row "comment<TAB>1<TAB>result/plan-review.md<TAB><the SHA-256 of {comment}><TAB>—".
`,
	"specification": `Write the specification of the task {task} (kind {kind}, issue {issue}) from the revised draft and its facts listed above: spec/draft.tsv, spec/spans.tsv, start/problem-statement.md and intake/answers.tsv. Commit your work in this repository in one or more commits; push nothing.
The IDs. The row R-NNN of kind requirement in spec/draft.tsv is the requirement REQ-NNN, and the row R-NNN of kind non-functional is NFR-NNN, with the same three digits. A row of kind success-criterion is no requirement and has no such ID. Keep the IDs as the draft listed above numbers them; never renumber one.
In docs/prd/ write the product requirements document of the target, with one row for each requirement ID, in these forms (every line starts at the first column):

| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| REQ-<NNN> | <the requirement> |  |  | F-{facts}#<n> F-{record} |

| Requirement | Acceptance criterion |
| ----------- | -------------------- |
| ` + "`" + `REQ-<NNN>` + "`" + ` | <the acceptance criterion> |

## 9. Success metrics
| Metric | Fact |
| ------ | ---- |
| <the success criterion> | F-{facts}#<n> |

| REQ | Facts | Guardrail | ADR | Task | Test |
| --- | ----- | --------- | --- | ---- | ---- |
| REQ-<NNN> | F-{facts}#<n> F-{record} | — | — | — | — |

Write NFR-<NNN> in place of REQ-<NNN> for a row of kind non-functional, and keep the backticks around the ID in the acceptance table. Leave the MoSCoW and Phase cells empty (not —): code fills them at the bet. In a Facts cell, cite each fact that the row covers: F-{facts}#<n> for the fact n of spec/spans.tsv, and Q-<NNN> for a fact of intake/answers.tsv; always cite F-{record}, the existing record, too. Write one acceptance row and one matrix row for each requirement ID, and one success row for each success-criterion row of the draft, citing the facts it covers.
In docs/spec/ write, for each requirement ID, a section whose heading holds the ID as a whole word, for example "## REQ-<NNN> — <title>", with a non-empty text under it.
Write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/. Write nothing under docs/facts/: code renders it at the first bet.
The required outputs for the result file are docs/tasks/{task}.md and each file of docs/prd/ and docs/spec/ that you wrote.
`,
	"shape-architecture": `Write the architecture of the target for the task {task} (kind {kind}, issue {issue}) from the options tables of its decisions and the survivors of the filter, both listed above under decision-records, with the revised draft and its facts. Commit your work in this repository in one or more commits; push nothing.
For each decision that has one survivor, write that option as decided. For each decision that has more than one survivor, write a section with the heading "Decision pending the bet" that lists its survivors and decides nothing. Never choose among survivors, and never write an option that is not a survivor.
Write the architecture under docs/spec/ and the decision records under docs/adr/, in the form of the rule files listed above, and write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/.
The required outputs for the result file are docs/tasks/{task}.md and each file of docs/spec/ and docs/adr/ that you wrote.
`,
	"implement": `Implement the task {task} (kind {kind}, issue {issue}) as its plan says: the plan and its review are among the handoffs listed above. For each test the plan lists, write the failing test first, watch it fail for the right reason, then write the change. Answer each finding listed under the findings of the last attempt. When tests.tsv is listed above, change none of the test files it names. Keep docs/tasks/{task}.md, the task's file, current. Commit your work in this repository in one or more commits; push nothing.
The required outputs for the result file are docs/tasks/{task}.md and each file you wrote or changed.
`,
	"verification": `Verify the change of the task {task} (kind {kind}, issue {issue}) from {from} to {head}, the commit your repository is at: review the changes from {from} to {head} against the plan, the requirements, the decision records, the rule files and the frozen list tests.tsv, as listed above. Change no file of the repository and make no commit.
Write the review record as the file {comment}, in this form, each line starting at the first column and ended by a line feed:
head: {head}
verdict: <material, nothing material in scope or not mergeable, findings recorded>
findings:
- <path>:<line> <rule>: <text>
The line head: is written as shown, once. The value of verdict: is exactly one of the three strings material, nothing material in scope and not mergeable, findings recorded. The line findings: stands alone, once. Under it write one line for each material finding: <path> is relative to the repository, clean and inside it (no space, no .., no leading /), <line> is a positive integer, <rule> is one word with no colon, and <text> is not empty. For nothing material in scope write no line under findings:. This is the form the target reads; it is not the table "| Field | Value |" of LAYUP's own reviews, which you must not write.
In the result file write the row "comment<TAB>1<TAB>result/review-record.md<TAB><the SHA-256 of {comment}><TAB>—".
`,
	"answer-to-a-question": answerInstruction,
	"close-out": `Close out the task {task} (kind {kind}, issue {issue}). From the verdict among the handoffs listed above and the ledger (sessions.tsv and telemetry.tsv, listed above), write the task's verdict, its line of the completed log and its resource record. Make one commit that changes only docs/tasks/{task}.md and docs/tasks/completed.md; push nothing.
The required outputs for the result file are docs/tasks/{task}.md and docs/tasks/completed.md.
`,
}

// answerInstruction is the instruction of the step answer-to-a-question, the block instruction-answer-to-a-question of loop.md (row 66z).
const answerInstruction = `Answer the questions that the sessions of the last attempt of the task {task} (kind {kind}, issue {issue}) asked: {questions}. The question rows are in the result files of those sessions, listed above under handoffs-and-answers; the IDs above are in the order of the questions, each a session and the n of its row. Read the files each question concerns, the requirements, the decision records and the rule files. Change no file of the repository and make no commit.
Write one row of the kind answer for each question, and no other row of that kind: "answer<TAB><n><TAB><the answer, one line><TAB>—<TAB><record>". The answer row of column n = 1 answers the first question of the list above, the row of column n = 2 the second, and so on. At least one <record> names a file the answer rests on, more are separated by one space, and there is no space at either end of the cell. Write the status row completed when every question is answered.
`

// answersInstruction is the block instruction-answers of loop.md: it follows the result text of a step when the attempt before has an
// answer of its owner role (row 66z).
const answersInstruction = `The sessions of the last attempt asked questions, and the owner role answered them; the answers are listed above under handoffs-and-answers. Their IDs are {answers}. Write each of these IDs as a word of the reason of the status row, and write none of them in a question row.
`

// resultInstruction is the text shared by all the steps: the form of the result file,
// the block instruction-result of loop.md.
const resultInstruction = `Write the result file {result}: UTF-8 text with no byte-order mark, one row per line, each line ended by a line feed, the last too, and no blank line. Separate the cells of a row with one tab character each, never with spaces. The first line is "kind<TAB>n<TAB>value<TAB>sha256<TAB>reason". The single character — (U+2014) is the empty cell; no cell holds a tab or a line break. The rows of each kind are numbered 1, 2, 3 and so on, with no gap.
Write exactly one status row, "status<TAB>1<TAB><status><TAB>—<TAB><reason>": <status> is completed when the work of this step is done, needs_context when a question stops it, and blocked or failed otherwise; <reason> is one line and not empty.
For each required output that the work above names, write the row "artifact<TAB><n><TAB><path><TAB><sha256><TAB>—": <path> is relative to the repository, and <sha256> is 64 lowercase hexadecimal characters, the SHA-256 of that file at your last commit. A file that this prompt only asks you to read is no output.
When a question stops the work, write one row for each question, "question<TAB><n><TAB><text><TAB>—<TAB><label> <needs_human> <record>": <label> is one of {labels}; <needs_human> is yes or no; at least one <record> names a file the question concerns, more are separated by one space, and there is no space at either end of the cell.
A comment file is plain Markdown with no code fence. Push nothing. Then end.
`

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

// commentFiles is the comment file of each step that writes one.
var commentFiles = map[string]string{"plan-of-a-task": PlanFile, "plan-review": PlanReviewFile, "verification": ReviewRecordFile}

// stepKinds is the kind of task that each of three steps belongs to.
var stepKinds = map[string]string{"specification": "specification", "shape-architecture": "architecture", "implement": "build"}

// checkStepKind refuses a step that is not a step of a task of the kind: the
// steps specification, shape-architecture and implement take one kind each, any
// other step takes any.
func checkStepKind(step, kind string) error {
	if want, ok := stepKinds[step]; ok && want != kind {
		return fmt.Errorf("the step %s is not a step of a task of kind %s", step, kind)
	}
	return nil
}

// inputsWanted says which inputs of the task a step receives: the revised draft
// and its facts for each step but the close-out of a specification or an
// architecture task, the options for the same steps of an architecture task, and
// the ledger for the close-out (loop.md, The records each step receives).
func inputsWanted(step, kind string) (shape, options, ledger bool) {
	ledger = step == "close-out"
	shape = !ledger && (kind == "specification" || kind == "architecture")
	return shape, shape && kind == "architecture", ledger
}

var decisionName = regexp.MustCompile(`^[a-z0-9-]+$`)

// optionsKept reads the table of each decision of opts and checks the result of
// the filter against it: the options of the result are the table's, in the order
// of their first row, there is a survivor, and each survivor is one of the options.
// It gives the tables as records, sorted by decision, and the survivors.
func optionsKept(opts map[string]spec.Result, read func(string) ([]byte, error)) ([]kept, []survivor, error) {
	var files []kept
	var surv []survivor
	for _, d := range slices.Sorted(maps.Keys(opts)) {
		res := opts[d]
		if !decisionName.MatchString(d) {
			return nil, nil, fmt.Errorf("decision %s: not a word of lowercase letters, digits and -", d)
		}
		p := "spec/options/" + d + ".tsv"
		data, err := read(p)
		if err != nil {
			return nil, nil, err
		}
		rows, err := spec.ReadOptions(data)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: %w", p, err)
		}
		var table, got []string
		for _, r := range rows {
			if !slices.Contains(table, r[0]) {
				table = append(table, r[0])
			}
		}
		for _, o := range res.Options {
			got = append(got, o.Name)
		}
		s := survivor{Decision: d}
		for _, o := range res.Survivors {
			s.Options = append(s.Options, o.Name)
		}
		switch i := slices.IndexFunc(s.Options, func(n string) bool { return !slices.Contains(table, n) }); {
		case !slices.Equal(got, table):
			return nil, nil, fmt.Errorf("decision %s: the options of the result are %v, those of %s are %v", d, got, p, table)
		case len(s.Options) == 0:
			return nil, nil, fmt.Errorf("decision %s: the result has no survivor", d)
		case i >= 0:
			return nil, nil, fmt.Errorf("decision %s: the survivor %s is not an option of the table", d, s.Options[i])
		}
		files, surv = append(files, kept{p, data}), append(surv, s)
	}
	return files, surv, nil
}

// labelList gives the labels of an owner map, sorted, joined by a comma and a space.
func labelList(owners map[string]string) string {
	return strings.Join(slices.Sorted(maps.Keys(owners)), ", ")
}

// fill replaces each {name} of text by the value of name.
func fill(text string, vars map[string]string) string {
	for _, k := range slices.Sorted(maps.Keys(vars)) {
		text = strings.ReplaceAll(text, "{"+k+"}", vars[k])
	}
	return text
}
