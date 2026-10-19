package run

// This file holds the steps answers-review and follow-up of Intake
// (docs/spec/intake.md, The review of the answers and the follow-up; The steps
// of Intake; Input states; row 47b of the plan, task T-vwu3, #224): a session
// of a phase of systems-architect closes or opens each answered question, and
// one follow-up comment asks again the questions still asked or open and the
// fields of the intent form that do not parse.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/tsv"
)

// followUpHeader is the first line of the follow-up comment (decided here,
// task T-vwu3: it is not the batch's, so neither is read as the other).
const followUpHeader = "LAYUP Intake: the follow-up"

// reviewFile is the form of result/answers-review.tsv: the block
// answers-review without session, which code adds.
var reviewFile = sessionSchema(records.AnswersReviewSchema, "session")

// emptyReview is the record of a review with no answered question: the header
// alone, so a restart knows that the step ran (decided here, task T-vwu3).
func emptyReview() []byte {
	data, _ := table(records.AnswersReviewSchema, nil)
	return data
}

// reviewKeep gives the Keep of the step answers-review on the rows of
// questions.tsv: the session's file read by its form; its questions are the
// answered questions, each once; the record in the order of the questions with
// the session, and questions.tsv with each reviewed question closed or open.
func reviewKeep(questions [][]string) func(id string, data []byte) (map[string][]byte, error) {
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, reviewFile)
		if err != nil {
			return nil, err
		}
		review := map[string][]string{}
		for _, r := range rows {
			i := slices.IndexFunc(questions, func(q []string) bool { return q[0] == r[0] })
			switch {
			case i < 0 || questions[i][12] != "answered":
				return nil, fmt.Errorf("%s is not an answered question", r[0])
			case review[r[0]] != nil:
				return nil, fmt.Errorf("%s twice", r[0])
			}
			review[r[0]] = r
		}
		var out [][]string
		updated := make([][]string, len(questions))
		for i, q := range questions {
			updated[i] = slices.Clone(q)
			if q[12] != "answered" {
				continue
			}
			r, ok := review[q[0]]
			if !ok {
				return nil, fmt.Errorf("the answered question %s has no row", q[0])
			}
			out = append(out, []string{r[0], r[1], r[2], id})
			updated[i][12] = r[1]
		}
		files := map[string][]byte{}
		if files["intake/answers-review.tsv"], err = table(records.AnswersReviewSchema, out); err != nil {
			return nil, err
		}
		if _, err := records.ReadAnswersReview(files["intake/answers-review.tsv"]); err != nil {
			return nil, err
		}
		if files["questions.tsv"], err = table(records.QuestionsSchema, updated); err != nil {
			return nil, err
		}
		if _, err := records.ReadQuestions(files["questions.tsv"]); err != nil {
			return nil, err
		}
		return files, nil
	}
}

// reviewPrompt is the fixed text of the prompt of the review of the answers,
// with {questions}, {answers} and {output}, each path on a line of its own.
const reviewPrompt = `You review the answers of the idea owner to the questions of a problem
statement, for LAYUP. Change no file in this repository and make no commit.

The problem statement: {input}
The questions: {questions}
The answers: {answers}
The output file: {output}

The questions are tab-separated values, one row per question (its ID in the
column id, the question in the column question, its state in the column
state, and the byte offsets of its quote in the problem statement in the
columns start and end); the answers name the question they answer. Write the output file as
tab-separated values, with the header line "question<TAB>state<TAB>reason" and
one row for each question whose state is answered: the state closed when the
answer closes the gap, or open when it does not, and the reason on one line.
Then end.
`

func reviewText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{questions}", input["questions.tsv"],
		"{answers}", input["answers.tsv"], "{output}", output).Replace(reviewPrompt))
}

// answersReviewStep is the step answers-review of Intake: "done before" when
// intake/answers-review.tsv exists; with no answered question, the record with
// no row and no session; else one session of systems-architect on the problem
// statement, questions.tsv and answers.tsv of the records (intake.md, The
// sessions of Intake: the problem statement for each step).
func (r *state) answersReviewStep(ctx context.Context) Step {
	const name = "answers-review"
	if data, err := r.store.ReadFile(ctx, "intake/answers-review.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	var inputs [3][]byte
	for i, path := range []string{"questions.tsv", "intake/answers.tsv", "start/problem-statement.md"} {
		data, err := r.store.ReadFile(ctx, path)
		if err != nil {
			return fail(name, err)
		}
		if data == nil {
			return failf(name, "the records have no %s", path)
		}
		inputs[i] = data
	}
	questions, err := records.ReadQuestions(inputs[0])
	if err != nil {
		return fail(name, fmt.Errorf("questions.tsv: %w", err))
	}
	if !slices.ContainsFunc(questions, func(q []string) bool { return q[12] == "answered" }) {
		if err := r.write(ctx, map[string][]byte{"intake/answers-review.tsv": emptyReview()}, "layup run: no answer to review"); err != nil {
			return fail(name, err)
		}
		return done(name, "no answered question; no session")
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "systems-architect", r.cfg.Harnesses, r.cfg.Models, "")
	if errors.Is(err, route.ErrNoPair) {
		return failf(name, "pair: %v", err)
	} else if err != nil {
		return fail(name, err)
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	if _, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "systems-architect", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(), Input: map[string][]byte{"questions.tsv": inputs[0], "answers.tsv": inputs[1], "problem-statement.md": inputs[2]},
		Output: "answers-review.tsv", Prompt: reviewText, Keep: reviewKeep(questions)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "intake/answers-review.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadAnswersReview(data)
	if err != nil {
		return fail(name, err)
	}
	open := 0
	for _, row := range rows {
		if row[1] == "open" {
			open++
		}
	}
	return done(name, fmt.Sprintf("%d answers reviewed, %d open, in intake/answers-review.tsv", len(rows), open))
}

// intentParsed gives, for each question of the intent form, whether its answer
// parses by the rule of row 47a: a decimal once the spaces around it are
// trimmed, and B below U as numbers (a B not below U is an unparsed answer of
// both).
func intentParsed(questions, answers [][]string) map[string]bool {
	parsed := map[string]bool{}
	answerOf := map[string]string{}
	for _, a := range answers {
		answerOf[a[0]] = a[1]
	}
	field := map[string]string{} // the question of each field
	for _, q := range questions {
		if q[1] == "intent" {
			a, ok := answerOf[q[0]]
			parsed[q[0]] = ok && decimal.MatchString(strings.TrimSpace(a))
			field[q[2]] = q[0]
		}
	}
	b, u := field["band.B"], field["band.U"]
	if parsed[b] && parsed[u] {
		if _, ok := budgetOf(questions, [][]string{{b, answerOf[b]}, {u, answerOf[u]}, {field["appetite.hours"], "1.0"}}); !ok {
			parsed[b], parsed[u] = false, false
		}
	}
	return parsed
}

// followUpRows gives the rows of questions.tsv of the follow-up: one for each
// question of the batch still asked, each open one, and each field of the
// intent form that does not parse (budget.tsv is written exactly when the
// three parse, so no field of a written budget.tsv is named by this rule); one
// row per question, whatever the number of rules that name it; IDs after the
// last.
func followUpRows(questions, answers [][]string, asked string) [][]string {
	parsed := intentParsed(questions, answers)
	var rows [][]string
	for _, q := range questions {
		if q[1] == "follow-up" {
			continue
		}
		unparsed := q[1] == "intent" && slices.Contains(budgetFields, q[2]) && !parsed[q[0]]
		if q[12] != "asked" && q[12] != "open" && !unparsed {
			continue
		}
		rows = append(rows, []string{fmt.Sprintf("Q-%03d", len(questions)+len(rows)+1), "follow-up", q[0], q[3], q[4], "follow-up",
			"", "", q[8], asked, "layup run", "", "asked"})
	}
	return rows
}

// followUpComment gives the comment of the follow-up rows, in the form of the
// batch, each question in the block of the question it asks again, with its
// quote and, for a marker, the value that the session proposed (the rows of
// intake/markers.tsv; task T-4c5c, #276).
func followUpComment(questions, rows [][]string, psb []byte, markers [][]string) string {
	return askedAgainComment(followUpHeader, questions, rows, psb, markers)
}

// askedAgainComment is the comment of rows that ask questions again, under header: the follow-up's, and the step shape's
// (task T-3f94, #364).
func askedAgainComment(header string, questions, rows [][]string, psb []byte, markers [][]string) string {
	asked := map[string][]string{}
	for _, q := range questions {
		asked[q[0]] = q
	}
	quote := func(r []string) (string, string) {
		if q := asked[r[2]]; q != nil {
			return q[6], q[7]
		}
		return "", ""
	}
	about := func(r []string) (bool, string) {
		if q := asked[r[2]]; q != nil {
			return aboutAMarker(q, markers)
		}
		return false, ""
	}
	return questionsComment(header, rows, quote, about, psb)
}

// followUpStep is the step follow-up of Intake: with follow-up rows in
// questions.tsv and the App's comment of the follow-up, "done before"; with
// the rows and no comment, the comment posted from the record; else the rows
// to ask again, in one records commit, then one comment on the Intake issue;
// with none, Intake is done and nothing is written.
func (r *state) followUpStep(ctx context.Context) Step {
	const name = "follow-up"
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return failf(name, "start.tsv has no number of the Intake issue")
	}
	psbData, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	if psbData == nil {
		return failf(name, "the records have no start/problem-statement.md")
	}
	questions, err := r.readRecord(ctx, "questions.tsv", records.ReadQuestions)
	if err != nil {
		return fail(name, err)
	}
	markers, err := r.readRecord(ctx, "intake/markers.tsv", records.ReadMarkers)
	if err != nil {
		return fail(name, err)
	}
	held := intakeFollowUps(questions)
	if held != nil {
		comments, err := r.cfg.Forge.Comments(ctx, issue)
		if err != nil {
			return fail(name, err)
		}
		for _, c := range comments {
			if c.App == r.cfg.Register.AppSlug && strings.HasPrefix(c.Body, followUpHeader+"\n") {
				return done(name, "done before")
			}
		}
		if _, err := r.cfg.Forge.Comment(ctx, issue, followUpComment(questions, held, psbData, markers)); err != nil {
			return fail(name, err)
		}
		return done(name, "the comment of the follow-up posted from questions.tsv")
	}
	answers, err := r.readRecord(ctx, "intake/answers.tsv", records.ReadAnswers)
	if err != nil {
		return fail(name, err)
	}
	rows := followUpRows(questions, answers, stamp(r.cfg.Clock.Now()))
	if rows == nil {
		return done(name, "nothing to ask again; Intake is done")
	}
	data, err := table(records.QuestionsSchema, append(slices.Clone(questions), rows...))
	if err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadQuestions(data); err != nil {
		return fail(name, err)
	}
	if err := r.write(ctx, map[string][]byte{"questions.tsv": data}, "layup run: the follow-up"); err != nil {
		return fail(name, err)
	}
	if _, err := r.cfg.Forge.Comment(ctx, issue, followUpComment(questions, rows, psbData, markers)); err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d questions asked again, in questions.tsv and the comment of the follow-up", len(rows)))
}

// intakeFollowUps gives the rows of the follow-up of Intake: the rows of source follow-up that are not asked by the step
// shape (kind shape, task T-3f94).
func intakeFollowUps(questions [][]string) [][]string {
	var out [][]string
	for _, q := range questions {
		if q[1] == "follow-up" && !isShapeRow(q) {
			out = append(out, q)
		}
	}
	return out
}
