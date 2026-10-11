package run

// This file holds the step answers of Intake (docs/spec/intake.md, The
// answers; The raw fact of the answers; The steps of Intake; Input states;
// row 47a of the plan, task T-hxb9, #223): the comments of the Intake issue,
// copied, give the answers of the batch, and the intent form gives budget.tsv.

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/records"
)

// answerCopy is a row of copies.tsv and the body of its file.
type answerCopy struct {
	Row  []string
	Body string
}

// answerLine is the rule of an answer line (decided here, task T-hxb9): after
// a CR at its end is trimmed, the line starts with Q- and the question's
// number, then ": " exactly; the rest is the answer, byte for byte.
var answerLine = regexp.MustCompile(`^(Q-[0-9]{3,}): (.*)$`)

// parseAnswer gives the question ID and the answer of a line, or false: an
// empty answer, and one of exactly —, which the writer would write as the
// empty value, answer nothing.
func parseAnswer(line string) (string, string, bool) {
	m := answerLine.FindStringSubmatch(strings.TrimSuffix(line, "\r"))
	if m == nil || m[2] == "" || m[2] == "—" {
		return "", "", false
	}
	return m[1], m[2], true
}

// collectAnswers gives the rows of intake/answers.tsv of the questions given
// (the batch's, or the follow-up's, as records.ReadQuestions gives them), in
// the order of the questions, and the blocks that still wait. A copy counts when it is of the
// issue, made after the batch's asked (strictly, by its created), and a human
// decision of an approver in the role of the question's block (Decision: a
// first copy only, so an edit answers nothing; decided here, task T-hxb9, with
// run.md, A human decision). The last line of an ID, in the order of the
// copies, is its answer. A block that asks is in hand when a counted copy by
// its role holds an answer line of its questions.
func collectAnswers(questions, approvers [][]string, copies []answerCopy, issue string) ([][]string, []string) {
	if len(questions) == 0 {
		return nil, nil
	}
	asked := questions[0][9]
	block := map[string]string{} // the block of each question
	var blocks []string
	for _, q := range questions {
		block[q[0]] = q[3]
		if !slices.Contains(blocks, q[3]) {
			blocks = append(blocks, q[3])
		}
	}
	got := map[string][]string{}
	inHand := map[string]bool{}
	for _, c := range copies {
		if c.Row[2] != issue || c.Row[6] <= asked {
			continue
		}
		for _, line := range strings.Split(c.Body, "\n") {
			id, answer, ok := parseAnswer(line)
			b, asks := block[id]
			if !ok || !asks || !Decision(c.Row, approvers, b) {
				continue
			}
			got[id] = []string{id, answer, c.Row[0], c.Row[3], c.Row[6]}
			inHand[b] = true
		}
	}
	var waiting []string
	for _, b := range blocks {
		if !inHand[b] {
			waiting = append(waiting, b)
		}
	}
	var answers [][]string
	for _, q := range questions {
		if a, ok := got[q[0]]; ok {
			answers = append(answers, a)
		}
	}
	return answers, waiting
}

// decimal is the form of the type decimal of internal/tsv.
var decimal = regexp.MustCompile(`^(0|[1-9][0-9]*)\.[0-9]+$`)

// budgetOf gives the rows of budget.tsv when the three fields of the intent
// form parse: each answer, trimmed of the spaces around it (decided here, task
// T-hxb9), a decimal, and B below U as numbers. Else none and false.
func budgetOf(questions, answers [][]string) ([][]string, bool) {
	answerOf := map[string]string{}
	for _, a := range answers {
		answerOf[a[0]] = a[1]
	}
	value := map[string][2]string{} // the value and the question of each field
	for _, q := range questions {
		if q[1] != "intent" {
			continue
		}
		a, ok := answerOf[q[0]]
		if a = strings.TrimSpace(a); !ok || !decimal.MatchString(a) {
			return nil, false
		}
		value[q[2]] = [2]string{a, q[0]}
	}
	var rows [][]string
	for _, f := range intentFields {
		v, ok := value[f]
		if !ok {
			return nil, false
		}
		rows = append(rows, []string{f, v[0], v[1]})
	}
	b, _ := new(big.Rat).SetString(value["band.B"][0])
	u, _ := new(big.Rat).SetString(value["band.U"][0])
	if b.Cmp(u) >= 0 {
		return nil, false
	}
	return rows, true
}

// answersStep is the step answers of Intake: "done before" when
// intake/answers.tsv exists; else the comments of the Intake issue are copied
// (but the App's own, run.md, Copy before read), the answers collected, and
// while a block that asks waits, the step is done "waiting", the last row of
// the restart; else one records commit of intake/answers.tsv, questions.tsv
// with the answered states, and budget.tsv when the intent form parses.
func (r *state) answersStep(ctx context.Context) Step {
	const name = "answers"
	if data, err := r.store.ReadFile(ctx, "intake/answers.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return failf(name, "start.tsv has no number of the Intake issue")
	}
	questions, err := r.readRecord(ctx, "questions.tsv", records.ReadQuestions)
	if err != nil {
		return fail(name, err)
	}
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return fail(name, err)
	}
	for _, c := range comments {
		if c.App == r.cfg.Register.AppSlug {
			continue
		}
		if err := r.copy(ctx, c, issue); err != nil {
			return fail(name, err)
		}
	}
	var copies []answerCopy
	for _, row := range r.copies {
		if row[2] != strconv.Itoa(issue) {
			continue
		}
		body, err := r.store.ReadFile(ctx, row[9])
		if err != nil {
			return fail(name, err)
		}
		copies = append(copies, answerCopy{Row: row, Body: string(body)})
	}
	answers, waiting := collectAnswers(questions, r.approvers, copies, strconv.Itoa(issue))
	if len(waiting) > 0 {
		return done(name, "waiting: no answer yet in the block "+strings.Join(waiting, " and the block "))
	}
	when := map[string]string{}
	for _, a := range answers {
		when[a[0]] = a[4]
	}
	updated := make([][]string, len(questions))
	for i, q := range questions {
		updated[i] = slices.Clone(q)
		if t, ok := when[q[0]]; ok {
			updated[i][11], updated[i][12] = t, "answered"
		}
	}
	files := map[string][]byte{}
	for _, f := range []struct {
		path  string
		write func() ([]byte, error)
	}{
		{"intake/answers.tsv", func() ([]byte, error) { return table(records.AnswersSchema, answers) }},
		{"questions.tsv", func() ([]byte, error) { return table(records.QuestionsSchema, updated) }},
	} {
		if files[f.path], err = f.write(); err != nil {
			return fail(name, err)
		}
	}
	if _, err := records.ReadAnswers(files["intake/answers.tsv"]); err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadQuestions(files["questions.tsv"]); err != nil {
		return fail(name, err)
	}
	written := "budget.tsv not written (a field of the intent form does not parse)"
	if rows, ok := budgetOf(questions, answers); ok {
		if files["budget.tsv"], err = table(records.BudgetSchema, rows); err != nil {
			return fail(name, err)
		}
		written = "budget.tsv written"
	}
	if err := r.write(ctx, files, "layup run: the answers of the batch"); err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d answers, in intake/answers.tsv; %s", len(answers), written))
}
