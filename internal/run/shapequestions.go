package run

// This file holds the questions left open of the step shape (docs/spec/loop.md, The questions left open; task T-3f94, #364,
// row 66h of the plan of M2e).

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// shapeQuestionsHeader is the first line of the comment that asks the questions left open.
const shapeQuestionsHeader = "LAYUP Shape: the questions left open"

// isShapeRow reports whether a row of questions.tsv is asked by the step shape.
func isShapeRow(q []string) bool { return q[1] == "follow-up" && q[5] == "shape" }

// leftOpen gives the IDs of the questions to ask, in the order of questions.tsv: a question that is not a follow-up, that no Shape row
// still asks, and whose state is open, or asked with no answered follow-up that asks it again.
func leftOpen(questions [][]string) []string {
	asking, again := map[string]bool{}, map[string]bool{}
	for _, q := range questions {
		if isShapeRow(q) && q[12] == "asked" {
			asking[q[2]] = true
		}
		if q[1] == "follow-up" && q[12] == "answered" {
			again[q[2]] = true
		}
	}
	var out []string
	for _, q := range questions {
		if q[1] != "follow-up" && !asking[q[0]] && (q[12] == "open" || (q[12] == "asked" && !again[q[0]])) {
			out = append(out, q[0])
		}
	}
	return out
}

// shapeRows gives the rows of questions.tsv that ask the questions ids again, in the block of each, with the IDs after the last.
func shapeRows(questions [][]string, ids []string, asked string) [][]string {
	byID := map[string][]string{}
	for _, q := range questions {
		byID[q[0]] = q
	}
	var rows [][]string
	for i, id := range ids {
		q := byID[id]
		rows = append(rows, []string{fmt.Sprintf("Q-%03d", len(questions)+i+1), "follow-up", id, q[3], q[4], "shape", "", "", q[8], asked, "layup run", "", "asked"})
	}
	return rows
}

// shapeComment gives the comment of the rows of a round: the form of the follow-up under the header of Shape.
func shapeComment(questions, rows [][]string, psb []byte, markers [][]string) string {
	return askedAgainComment(shapeQuestionsHeader, questions, rows, psb, markers)
}

// applyShapeAnswers gives questions.tsv and budget.tsv after the answers got (rows of answers.tsv to the Shape questions). A Shape
// row and the question it asks again become answered at the time of the answer, but a field of the intent form whose value does not
// parse stays open; a Shape row still asked is left out of the reading of budget.tsv. The budget is that of budgetAfterFollowUp: the
// last Shape row of a field is its source, and a value is replaced only while the three parse with B below U. The new budget is nil
// when budget.tsv stays unwritten.
func applyShapeAnswers(questions, answers, budget, got [][]string) (updated, newBudget [][]string) {
	when := map[string]string{}
	for _, a := range got {
		when[a[0]] = a[4]
	}
	var read [][]string
	for _, q := range questions {
		if _, ok := when[q[0]]; !isShapeRow(q) || q[12] != "asked" || ok {
			read = append(read, q)
		}
	}
	newBudget, open := budgetAfterFollowUp(read, append(slices.Clone(answers), got...), budget)
	updated = make([][]string, len(questions))
	for i, q := range questions {
		updated[i] = slices.Clone(q)
	}
	set := func(id, time, state string) {
		for _, q := range updated {
			if q[0] == id {
				q[11], q[12] = time, state
			}
		}
	}
	for _, q := range questions {
		if t, ok := when[q[0]]; ok && isShapeRow(q) {
			set(q[0], t, "answered")
			if slices.Contains(open, q[2]) {
				set(q[2], "", "open")
			} else {
				set(q[2], t, "answered")
			}
		}
	}
	return updated, newBudget
}

// intentWait gives the questions of the intent form that block Shape: those of band.B, band.U and appetite.hours that are open or
// that a Shape row still asks.
func intentWait(questions [][]string) []string {
	asking := map[string]bool{}
	for _, q := range questions {
		if isShapeRow(q) && q[12] == "asked" {
			asking[q[2]] = true
		}
	}
	var out []string
	for _, q := range questions {
		if q[1] == "intent" && slices.Contains(budgetFields, q[2]) && (q[12] == "open" || asking[q[0]]) {
			out = append(out, q[0])
		}
	}
	return out
}

// shapeQuestions is the first stage of the step: it posts and copies the comment of the rows that are asked, records the answers, asks
// the questions that are left open, and gives the questions of the intent form that the run waits for. It reads the records and the
// forge at each run: the rows, the comment and its copy are each made once, whichever of them an earlier run stopped before.
func (r *state) shapeQuestions(ctx context.Context) ([]string, error) {
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return nil, errors.New("start.tsv has no number of the Intake issue")
	}
	read := func(path string, f func([]byte) ([][]string, error)) ([][]string, error) {
		return r.readRecord(ctx, path, f)
	}
	questions, err := read("questions.tsv", records.ReadQuestions)
	if err != nil {
		return nil, err
	}
	markers, err := read("intake/markers.tsv", records.ReadMarkers)
	if err != nil {
		return nil, err
	}
	answers, err := read("intake/answers.tsv", records.ReadAnswers)
	if err != nil {
		return nil, err
	}
	psb, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return nil, err
	} else if psb == nil {
		return nil, errors.New("the records have no start/problem-statement.md")
	}
	before, err := r.store.ReadFile(ctx, "budget.tsv")
	if err != nil {
		return nil, err
	}
	var held [][]string
	if before != nil {
		if held, err = records.ReadBudget(before); err != nil {
			return nil, fmt.Errorf("budget.tsv: %w", err)
		}
	}
	// The rows that are asked have their comment, posted and copied.
	if err := r.shapePosted(ctx, issue, questions, psb, markers); err != nil {
		return nil, err
	}
	if pending := shapeAsked(questions); len(pending) > 0 {
		comments, err := r.cfg.Forge.Comments(ctx, issue)
		if err != nil {
			return nil, err
		}
		for _, c := range comments {
			if c.App == r.cfg.Register.AppSlug {
				continue
			}
			if err := r.copy(ctx, c, issue); err != nil {
				return nil, err
			}
		}
		var copies []answerCopy
		for _, row := range r.copies {
			if row[2] != strconv.Itoa(issue) {
				continue
			}
			body, err := r.store.ReadFile(ctx, row[9])
			if err != nil {
				return nil, err
			}
			copies = append(copies, answerCopy{Row: row, Body: string(body)})
		}
		got, _ := collectAnswers(pending, r.approvers, copies, strconv.Itoa(issue))
		if len(got) > 0 {
			bound, named := pointBindings(copies, r.approvers, pending[0][9], strconv.Itoa(issue))
			newApprovers, err := approverRows(ctx, r.cfg.Forge.UserID, bound, named, r.approvers)
			if err != nil {
				return nil, err
			}
			var budget [][]string
			if questions, budget = applyShapeAnswers(questions, answers, held, got); budget != nil {
				held = budget
			}
			answers = append(slices.Clone(answers), got...)
			files := map[string][]byte{}
			if files["intake/answers.tsv"], err = table(records.AnswersSchema, answers); err != nil {
				return nil, err
			}
			if files["questions.tsv"], err = table(records.QuestionsSchema, questions); err != nil {
				return nil, err
			}
			if budget != nil {
				data, err := table(records.BudgetSchema, budget)
				if err != nil {
					return nil, err
				}
				if string(data) != string(before) {
					files["budget.tsv"], before = data, data
				}
			}
			if len(newApprovers) > 0 {
				if files["approvers.tsv"], err = table(records.ApproversSchema, append(slices.Clone(r.approvers), newApprovers...)); err != nil {
					return nil, err
				}
			}
			for p, read := range map[string]func([]byte) ([][]string, error){"intake/answers.tsv": records.ReadAnswers, "questions.tsv": records.ReadQuestions, "budget.tsv": records.ReadBudget, "approvers.tsv": records.ReadApprovers} {
				if data, ok := files[p]; ok {
					if _, err := read(data); err != nil {
						return nil, fmt.Errorf("%s: %w", p, err)
					}
				}
			}
			if err := r.write(ctx, files, "layup run: the answers to the questions left open"); err != nil {
				return nil, err
			}
			r.approvers = append(r.approvers, newApprovers...)
		}
	}
	// The questions that are left open are asked.
	if ids := leftOpen(questions); len(ids) > 0 {
		rows := shapeRows(questions, ids, stamp(r.cfg.Clock.Now()))
		questions = append(slices.Clone(questions), rows...)
		data, err := table(records.QuestionsSchema, questions)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadQuestions(data); err != nil {
			return nil, err
		}
		if err := r.write(ctx, map[string][]byte{"questions.tsv": data}, "layup run: the questions left open"); err != nil {
			return nil, err
		}
		if err := r.shapePosted(ctx, issue, questions, psb, markers); err != nil {
			return nil, err
		}
	}
	wait := intentWait(questions)
	if before == nil && len(wait) == 0 {
		return nil, errors.New("budget.tsv is not written and no question of the intent form is open")
	}
	return wait, nil
}

// shapeAsked gives the rows of the questions that the step shape asks and that are still asked.
func shapeAsked(questions [][]string) [][]string {
	var out [][]string
	for _, q := range questions {
		if isShapeRow(q) && q[12] == "asked" {
			out = append(out, q)
		}
	}
	return out
}

// shapePosted makes sure that each round of the rows that are asked has its comment on the Intake issue and its copy in the records:
// the comment is posted when no comment of the App has its body, under the lease, and copied whenever it is found.
func (r *state) shapePosted(ctx context.Context, issue int, questions [][]string, psb []byte, markers [][]string) error {
	rounds := map[string][][]string{} // the rows of each round by their asked
	var order []string
	for _, q := range shapeAsked(questions) {
		if _, ok := rounds[q[9]]; !ok {
			order = append(order, q[9])
		}
		rounds[q[9]] = append(rounds[q[9]], q)
	}
	for _, asked := range order {
		body := shapeComment(questions, rounds[asked], psb, markers)
		find := func() (forge.Comment, bool, error) {
			comments, err := r.cfg.Forge.Comments(ctx, issue)
			for _, c := range comments {
				if c.App == r.cfg.Register.AppSlug && c.Body == body {
					return c, true, err
				}
			}
			return forge.Comment{}, false, err
		}
		c, found, err := find()
		if err != nil {
			return err
		}
		if !found {
			if err := r.sessions().holdsLease(ctx); err != nil {
				return err
			}
			if _, err := r.cfg.Forge.Comment(ctx, issue, body); err != nil {
				return err
			}
			if c, found, err = find(); err != nil {
				return err
			} else if !found {
				return errors.New("the comment that asks the questions left open is not on the Intake issue after its post")
			}
		}
		if err := r.copy(ctx, c, issue); err != nil {
			return err
		}
	}
	return nil
}
