package records

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/tsv"
)

// The records of Intake of milestone M2c (docs/spec/records.md, NFR-001 — The
// records of Intake; task T-529i, #212): the schemas of intake/meaning.tsv,
// spec/spans.tsv, spec/draft.tsv, intake/coverage.tsv, questions.tsv,
// intake/answers.tsv, intake/answers-review.tsv and budget.tsv, and the rules
// that their blocks give in words, read in the committed form (a session's own
// file, without the columns that code adds, is the session of a phase's, row
// 42a). Each "— when" and "— for" is read as "exactly when", as for the
// records of a session.
//
// A rule that reads another file is its writer's, not a row's: the quote in
// the problem statement (row 43b), a covers number that is a fact of
// spec/spans.tsv (row 45a), band.B below band.U (rows 47a and 47c). The
// cover, the order and the numbering of the spans are rules of the span
// check, each with its row of spec-check (row 41a): a span table that breaks
// one is committed at the step numbering and reported by the step spans, and a
// restart runs numbering again, where a refusal by this reader would commit
// nothing (intake.md, The span check); so the row keeps only start below end
// and a fact of 1 or more.

// MeaningSchema is the form of intake/meaning.tsv: the block meaning.
var MeaningSchema = tsv.Schema{Name: "meaning", Location: "records:intake/meaning.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(M-NNN)", Key: true},
	{Name: "kind", Type: "enum(two-readings|unmeasurable|missing-intent|new-work|fit|no-design|decision)"},
	{Name: "start", Type: "int"},
	{Name: "end", Type: "int"},
	{Name: "quote", Type: "text"},
	{Name: "question", Type: "text"},
	{Name: "session", Type: "id(S-xxxxxxxx)"},
}}

// SpansSchema is the form of spec/spans.tsv: the block spans.
var SpansSchema = tsv.Schema{Name: "spans", Location: "records:spec/spans.tsv", Columns: []tsv.Column{
	{Name: "start", Type: "int", Key: true},
	{Name: "end", Type: "int"},
	{Name: "fact", Type: "int"},
	{Name: "class", Type: "enum(need|constraint|measure|context)"},
	{Name: "reason", Type: "text"},
	{Name: "text", Type: "text"},
}}

// DraftSchema is the form of spec/draft.tsv: the block draft.
var DraftSchema = tsv.Schema{Name: "draft", Location: "records:spec/draft.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(R-NNN)", Key: true},
	{Name: "kind", Type: "enum(requirement|non-functional|success-criterion)"},
	{Name: "statement", Type: "text"},
	{Name: "covers", Type: "list(int)"},
	{Name: "criterion", Type: "text"},
	{Name: "session", Type: "id(S-xxxxxxxx)"},
}}

// CoverageSchema is the form of intake/coverage.tsv: the block coverage.
var CoverageSchema = tsv.Schema{Name: "coverage", Location: "records:intake/coverage.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(C-NNN)", Key: true},
	{Name: "kind", Type: "enum(uncovered|unsupported|class)"},
	{Name: "fact", Type: "int"},
	{Name: "requirement", Type: "id(R-NNN)"},
	{Name: "question", Type: "text"},
	{Name: "session", Type: "id(S-xxxxxxxx)"},
}}

// QuestionsSchema is the form of questions.tsv: the block questions.
var QuestionsSchema = tsv.Schema{Name: "questions", Location: "records:questions.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "id(Q-NNN)", Key: true},
	{Name: "source", Type: "enum(rule|meaning|coverage|intent|follow-up)"},
	{Name: "ref", Type: "text"},
	{Name: "block", Type: "enum(idea-owner|operator)"},
	{Name: "phase", Type: "enum(before-delivery)"},
	{Name: "kind", Type: "text"},
	{Name: "start", Type: "int"},
	{Name: "end", Type: "int"},
	{Name: "question", Type: "text"},
	{Name: "asked", Type: "time"},
	{Name: "asker", Type: "text"},
	{Name: "answered", Type: "time"},
	{Name: "state", Type: "enum(asked|answered|closed|open)"},
}}

// AnswersSchema is the form of intake/answers.tsv: the block answers.
var AnswersSchema = tsv.Schema{Name: "answers", Location: "records:intake/answers.tsv", Columns: []tsv.Column{
	{Name: "question", Type: "id(Q-NNN)", Key: true},
	{Name: "answer", Type: "text"},
	{Name: "comment", Type: "int"},
	{Name: "author", Type: "int"},
	{Name: "time", Type: "time"},
}}

// AnswersReviewSchema is the form of intake/answers-review.tsv: the block
// answers-review.
var AnswersReviewSchema = tsv.Schema{Name: "answers-review", Location: "records:intake/answers-review.tsv", Columns: []tsv.Column{
	{Name: "question", Type: "id(Q-NNN)", Key: true},
	{Name: "state", Type: "enum(closed|open)"},
	{Name: "reason", Type: "text"},
	{Name: "session", Type: "id(S-xxxxxxxx)"},
}}

// BudgetSchema is the form of budget.tsv: the block budget.
var BudgetSchema = tsv.Schema{Name: "budget", Location: "records:budget.tsv", Columns: []tsv.Column{
	{Name: "name", Type: "text", Key: true},
	{Name: "value", Type: "decimal"},
	{Name: "source", Type: "id(Q-NNN)"},
}}

// The columns whose rule has no clause for the empty value, which never hold
// it (docs/spec/README.md, The types).
var (
	meaningRequired       = []string{"id", "kind", "start", "end", "quote", "question", "session"}
	spansRequired         = []string{"start", "end", "text"}
	draftRequired         = []string{"id", "kind", "statement", "criterion", "session"}
	coverageRequired      = []string{"id", "kind", "question", "session"}
	questionsRequired     = []string{"id", "source", "ref", "block", "phase", "kind", "question", "asked", "asker", "state"}
	answersRequired       = []string{"question", "answer", "comment", "author", "time"}
	answersReviewRequired = []string{"question", "state", "reason", "session"}
	budgetRequired        = []string{"name", "value", "source"}
)

var (
	meaningKinds  = []string{"two-readings", "unmeasurable", "missing-intent", "new-work", "fit", "no-design", "decision"}
	coverageKinds = []string{"uncovered", "unsupported", "class"}
	ruleKinds     = []string{"G1", "G2", "G3", "G4", "G5"}
	intentFields  = []string{"band.B", "band.U", "appetite.hours"} // a later milestone adds its names here first
	qForm         = regexp.MustCompile(`^Q-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2,})$`)
	mForm         = regexp.MustCompile(`^M-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2,})$`)
	cForm         = regexp.MustCompile(`^C-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2,})$`)
)

// below reports whether a and b are ints with a below b.
func below(a, b string) bool {
	x, err1 := strconv.Atoi(a)
	y, err2 := strconv.Atoi(b)
	return err1 == nil && err2 == nil && x < y
}

// inOrder checks that the column id of rows runs prefix-001, prefix-002, ...
// in the order of the file. A broken order is a *tsv.Error.
func inOrder(s tsv.Schema, rows [][]string, column, prefix string) error {
	for i, r := range rows {
		if want := fmt.Sprintf("%s-%03d", prefix, i+1); fields(s, r)[column] != want {
			return &tsv.Error{Line: i + 2, Column: column, Reason: "the next ID is " + want}
		}
	}
	return nil
}

// readInOrder reads data by its schema and its row rules, then the order of
// its IDs.
func readInOrder(data []byte, s tsv.Schema, check func([]string) error, prefix string) ([][]string, error) {
	rows, err := read(data, s, check)
	if err != nil {
		return nil, err
	}
	if err := inOrder(s, rows, "id", prefix); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckMeaning checks the rules of the block meaning on one row: start below
// end, the quote the file's bytes from start to end (its length end - start),
// and the session's form.
func CheckMeaning(r []string) error {
	f := fields(MeaningSchema, r)
	if c := missing(f, meaningRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	start, _ := strconv.Atoi(f["start"])
	end, _ := strconv.Atoi(f["end"])
	switch {
	case start >= end:
		return &RowError{"end", "end is above start"}
	case len(f["quote"]) != end-start || !utf8.ValidString(f["quote"]):
		return &RowError{"quote", fmt.Sprintf("the quote is the file's bytes from start to end: %d bytes", end-start)}
	case !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	}
	return nil
}

// ReadMeaning reads intake/meaning.tsv by its schema, its row rules and the
// order of its IDs.
func ReadMeaning(data []byte) ([][]string, error) {
	return readInOrder(data, MeaningSchema, CheckMeaning, "M")
}

// CheckSpan checks the rules of the block spans on one row: start below end,
// a fact of 1 or more with its class and no reason, or no fact, with a reason
// and no class.
func CheckSpan(r []string) error {
	f := fields(SpansSchema, r)
	if c := missing(f, spansRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	isFact := f["fact"] != ""
	switch {
	case !below(f["start"], f["end"]):
		return &RowError{"end", "end is above start"}
	case isFact && !atLeastOne(f["fact"]):
		return &RowError{"fact", "a fact's number is 1 or more"}
	case isFact != (f["class"] != ""):
		return &RowError{"class", "class is the empty value exactly when fact is"}
	case isFact == (f["reason"] != ""):
		return &RowError{"reason", "reason is the empty value exactly when fact is not"}
	}
	return nil
}

// ReadSpans reads spec/spans.tsv by its schema and its row rules.
func ReadSpans(data []byte) ([][]string, error) { return read(data, SpansSchema, CheckSpan) }

// CheckDraft checks the rules of the block draft on one row: each covered
// fact 1 or more, and the session's form.
func CheckDraft(r []string) error {
	f := fields(DraftSchema, r)
	if c := missing(f, draftRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	for _, n := range strings.Fields(f["covers"]) {
		if !atLeastOne(n) {
			return &RowError{"covers", "each covered fact's number is 1 or more"}
		}
	}
	if !sessionForm.MatchString(f["session"]) {
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	}
	return nil
}

// ReadDraft reads spec/draft.tsv by its schema, its row rules and the order of
// its IDs.
func ReadDraft(data []byte) ([][]string, error) {
	return readInOrder(data, DraftSchema, CheckDraft, "R")
}

// CheckCoverage checks the rules of the block coverage on one row: a fact of 1
// or more exactly when the item is not unsupported, a requirement exactly when
// it is, and the session's form.
func CheckCoverage(r []string) error {
	f := fields(CoverageSchema, r)
	if c := missing(f, coverageRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	unsupported := f["kind"] == "unsupported"
	switch {
	case unsupported == (f["fact"] != ""):
		return &RowError{"fact", "fact is the empty value exactly for unsupported"}
	case !unsupported && !atLeastOne(f["fact"]):
		return &RowError{"fact", "a fact's number is 1 or more"}
	case unsupported != (f["requirement"] != ""):
		return &RowError{"requirement", "requirement is given exactly for unsupported"}
	case !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	}
	return nil
}

// ReadCoverage reads intake/coverage.tsv by its schema, its row rules and the
// order of its IDs.
func ReadCoverage(data []byte) ([][]string, error) {
	return readInOrder(data, CoverageSchema, CheckCoverage, "C")
}

// CheckQuestion checks the rules of the block questions on one row: the forms
// of ref, kind and asker by the source; start and end both empty or both given,
// start below end, empty for the intent form and an unsupported item; and
// answered empty exactly when the state is asked, either way for open.
func CheckQuestion(r []string) error {
	f := fields(QuestionsSchema, r)
	if c := missing(f, questionsRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	var refOK, kindOK, askerOK bool
	switch f["source"] {
	case "rule":
		refOK, kindOK, askerOK = qForm.MatchString(f["ref"]), slices.Contains(ruleKinds, f["kind"]), f["asker"] == "layup psb check"
	case "meaning":
		refOK, kindOK, askerOK = mForm.MatchString(f["ref"]), slices.Contains(meaningKinds, f["kind"]), sessionForm.MatchString(f["asker"])
	case "coverage":
		refOK, kindOK, askerOK = cForm.MatchString(f["ref"]), slices.Contains(coverageKinds, f["kind"]), sessionForm.MatchString(f["asker"])
	case "intent":
		refOK, kindOK, askerOK = slices.Contains(intentFields, f["ref"]), f["kind"] == "intent", f["asker"] == "layup run"
	case "follow-up":
		refOK, kindOK, askerOK = qForm.MatchString(f["ref"]), f["kind"] == "follow-up", f["asker"] == "layup run"
	}
	noQuote := f["source"] == "intent" || f["kind"] == "unsupported"
	switch {
	case !refOK:
		return &RowError{"ref", "ref names the source row of the question in the form of its source"}
	case !kindOK:
		return &RowError{"kind", "kind is a kind of the question's source"}
	case !askerOK:
		return &RowError{"asker", "the asker is layup psb check for a rule gap, the session for a gap of meaning or coverage, layup run for the intent form and a follow-up"}
	case f["start"] == "" && f["end"] != "":
		return &RowError{"start", "start and end are both the empty value or both given"}
	case noQuote && f["start"] != "":
		return &RowError{"start", "the intent form and an unsupported item have no quote"}
	case f["start"] != "" && !below(f["start"], f["end"]):
		return &RowError{"end", "end is given with start, and above it"}
	case f["state"] == "asked" && f["answered"] != "":
		return &RowError{"answered", "an asked question has no answer"}
	case (f["state"] == "answered" || f["state"] == "closed") && f["answered"] == "":
		return &RowError{"answered", "an answered or closed question has the time of its answer"}
	}
	return nil
}

// CheckQuestions checks the rules of questions.tsv that read more than one
// row: the IDs Q-001, ... in order, and a follow-up's ref names a question
// before it.
func CheckQuestions(rows [][]string) error {
	if err := inOrder(QuestionsSchema, rows, "id", "Q"); err != nil {
		return err
	}
	seen := map[string]bool{}
	for i, r := range rows {
		f := fields(QuestionsSchema, r)
		if f["source"] == "follow-up" && !seen[f["ref"]] {
			return &tsv.Error{Line: i + 2, Column: "ref", Reason: "a follow-up asks again a question before it"}
		}
		seen[f["id"]] = true
	}
	return nil
}

// ReadQuestions reads questions.tsv by its schema, its row rules and the rules
// of the table.
func ReadQuestions(data []byte) ([][]string, error) {
	rows, err := read(data, QuestionsSchema, CheckQuestion)
	if err != nil {
		return nil, err
	}
	if err := CheckQuestions(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckAnswer checks the block answers on one row: no column holds the empty
// value. A bound on the forge's IDs is none of the block's (as for copies).
func CheckAnswer(r []string) error {
	if c := missing(fields(AnswersSchema, r), answersRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	return nil
}

// ReadAnswers reads intake/answers.tsv by its schema and its row rules.
func ReadAnswers(data []byte) ([][]string, error) { return read(data, AnswersSchema, CheckAnswer) }

// CheckAnswerReview checks the block answers-review on one row: no column
// holds the empty value, and the session's form.
func CheckAnswerReview(r []string) error {
	f := fields(AnswersReviewSchema, r)
	if c := missing(f, answersReviewRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	if !sessionForm.MatchString(f["session"]) {
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	}
	return nil
}

// ReadAnswersReview reads intake/answers-review.tsv by its schema and its row
// rules.
func ReadAnswersReview(data []byte) ([][]string, error) {
	return read(data, AnswersReviewSchema, CheckAnswerReview)
}

// CheckBudget checks the block budget on one row: no column holds the empty
// value, and the name is one of the block's. band.B below band.U is its
// writers' (rows 47a and 47c), not the block's.
func CheckBudget(r []string) error {
	f := fields(BudgetSchema, r)
	if c := missing(f, budgetRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	if !slices.Contains(intentFields, f["name"]) {
		return &RowError{"name", "the name is one of " + strings.Join(intentFields, ", ")}
	}
	return nil
}

// ReadBudget reads budget.tsv by its schema and its row rules.
func ReadBudget(data []byte) ([][]string, error) { return read(data, BudgetSchema, CheckBudget) }
