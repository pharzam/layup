package run

// This file holds the step follow-up-answers of Intake (docs/spec/intake.md,
// The review of the answers and the follow-up: the answers to the follow-up;
// The steps of Intake; Input states; row 47c of the plan, task T-tmm2, #225):
// the answers to the follow-up are recorded, budget.tsv is written or a value
// of it replaced, and Intake ends.

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/records"
)

// budgetAfterFollowUp gives the rows of budget.tsv after the follow-up (nil
// when it stays unwritten) and the questions of the intent form to set open,
// from the rows of questions.tsv, the answers (the batch's and the
// follow-up's) and the rows of budget.tsv before (nil for none). Not yet
// written: each field from the answer to its follow-up when the follow-up
// asked it and it has an answer, else from the batch's; written when the
// three parse with B below U, each source the question whose answer gives it;
// else each field still missing or unparsed is open (a B not below U: both,
// decided here, task T-tmm2). Written already: each field asked again whose
// answer parses replaces its value, with the follow-up as its source, when
// the three still parse with B below U; else budget.tsv stays as it was and
// the fields asked again are open (decided here, task T-tmm2).
func budgetAfterFollowUp(questions, answers, budget [][]string) ([][]string, []string) {
	answerOf := map[string]string{}
	for _, a := range answers {
		answerOf[a[0]] = a[1]
	}
	batch := map[string]string{}  // the batch's question of each field
	follow := map[string]string{} // the follow-up of each question asked again
	for _, q := range questions {
		switch q[1] {
		case "intent":
			batch[q[2]] = q[0]
		case "follow-up":
			follow[q[2]] = q[0]
		}
	}
	parse := func(id string) (string, bool) {
		a, ok := answerOf[id]
		a = strings.TrimSpace(a)
		return a, ok && decimal.MatchString(a)
	}
	type value struct{ v, source string }
	pick := map[string]value{}
	var open []string
	below := func() bool {
		b, _ := new(big.Rat).SetString(pick["band.B"].v)
		u, _ := new(big.Rat).SetString(pick["band.U"].v)
		return b.Cmp(u) < 0
	}
	rows := func() [][]string {
		var out [][]string
		for _, f := range intentFields {
			out = append(out, []string{f, pick[f].v, pick[f].source})
		}
		return out
	}
	if budget == nil {
		for _, f := range intentFields {
			id := batch[f]
			if fu, asked := follow[id]; asked {
				if _, answered := answerOf[fu]; answered {
					id = fu
				}
			}
			if v, ok := parse(id); ok {
				pick[f] = value{v, id}
			} else {
				open = append(open, batch[f])
			}
		}
		if open != nil {
			return nil, open
		}
		if !below() {
			return nil, []string{batch["band.B"], batch["band.U"]}
		}
		return rows(), nil
	}
	for _, r := range budget {
		pick[r[0]] = value{r[1], r[2]}
	}
	before := rows()
	var asked []string
	for _, f := range intentFields {
		fu, ok := follow[batch[f]]
		if !ok {
			continue
		}
		asked = append(asked, batch[f])
		if v, ok := parse(fu); ok {
			pick[f] = value{v, fu}
		} else {
			open = append(open, batch[f])
		}
	}
	if !below() {
		return before, asked
	}
	return rows(), open
}

// followUpAnswersStep is the step follow-up-answers of Intake: with no
// follow-up row, nothing was asked again and Intake is done; "done before"
// when an answer to a follow-up is in intake/answers.tsv; else the comments of
// the Intake issue are copied (but the App's own), the answers to the
// follow-up collected by the rule of the step answers, and while a block that
// the follow-up asks has none, the step waits; else one records commit of the
// answers, their states, the intent questions still open and budget.tsv, and
// Intake ends.
func (r *state) followUpAnswersStep(ctx context.Context) Step {
	const name = "follow-up-answers"
	questions, err := r.readRecord(ctx, "questions.tsv", records.ReadQuestions)
	if err != nil {
		return fail(name, err)
	}
	var follow [][]string
	for _, q := range questions {
		if q[1] == "follow-up" {
			follow = append(follow, q)
		}
	}
	if follow == nil {
		return done(name, "no follow-up; Intake is done")
	}
	answers, err := r.readRecord(ctx, "intake/answers.tsv", records.ReadAnswers)
	if err != nil {
		return fail(name, err)
	}
	for _, a := range answers {
		if slices.ContainsFunc(follow, func(q []string) bool { return q[0] == a[0] }) {
			return done(name, "done before")
		}
	}
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return failf(name, "start.tsv has no number of the Intake issue")
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
	got, waiting := collectAnswers(follow, r.approvers, copies, strconv.Itoa(issue))
	if len(waiting) > 0 {
		return done(name, "waiting: no answer yet to the follow-up in the block "+strings.Join(waiting, " and the block "))
	}
	all := append(slices.Clone(answers), got...)
	when := map[string]string{}
	for _, a := range got {
		when[a[0]] = a[4]
	}
	updated := make([][]string, len(questions))
	for i, q := range questions {
		updated[i] = slices.Clone(q)
		if t, ok := when[q[0]]; ok {
			updated[i][11], updated[i][12] = t, "answered"
		}
	}
	before, err := r.store.ReadFile(ctx, "budget.tsv")
	if err != nil {
		return fail(name, err)
	}
	var held [][]string
	if before != nil {
		if held, err = records.ReadBudget(before); err != nil {
			return fail(name, fmt.Errorf("budget.tsv: %w", err))
		}
	}
	budget, open := budgetAfterFollowUp(questions, all, held)
	for i, q := range updated {
		if slices.Contains(open, q[0]) {
			updated[i][12] = "open"
		}
	}
	files := map[string][]byte{}
	if files["intake/answers.tsv"], err = table(records.AnswersSchema, all); err != nil {
		return fail(name, err)
	}
	if files["questions.tsv"], err = table(records.QuestionsSchema, updated); err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadAnswers(files["intake/answers.tsv"]); err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadQuestions(files["questions.tsv"]); err != nil {
		return fail(name, err)
	}
	said := "budget.tsv not written"
	if budget != nil {
		data, err := table(records.BudgetSchema, budget)
		if err != nil {
			return fail(name, err)
		}
		switch {
		case before == nil:
			files["budget.tsv"], said = data, "budget.tsv written"
		case string(data) != string(before):
			files["budget.tsv"], said = data, "budget.tsv replaced"
		default:
			said = "budget.tsv kept"
		}
	}
	if err := r.write(ctx, files, "layup run: the answers to the follow-up"); err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d answers to the follow-up, in intake/answers.tsv; %s; Intake is done", len(got), said))
}
