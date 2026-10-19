package run

// This file holds the questions left open of the step shape (docs/spec/loop.md, The questions left open; task T-3f94, #364,
// row 66h of the plan of M2e).

import "context"

// shapeQuestionsHeader is the first line of the comment that asks the questions left open.
const shapeQuestionsHeader = "LAYUP Shape: the questions left open"

// leftOpen gives the IDs of the questions to ask, in the order of questions.tsv.
func leftOpen(questions [][]string) []string { return nil }

// shapeRows gives the rows of questions.tsv that ask the questions ids again, with the IDs after the last.
func shapeRows(questions [][]string, ids []string, asked string) [][]string { return nil }

// shapeComment gives the comment of the rows of a round.
func shapeComment(questions, rows [][]string, psb []byte, markers [][]string) string { return "" }

// applyShapeAnswers gives questions.tsv and budget.tsv after the answers got (rows of answers.tsv to the Shape questions).
func applyShapeAnswers(questions, answers, budget, got [][]string) (updated, newBudget [][]string) {
	return questions, budget
}

// intentWait gives the questions of the intent form that block Shape: those of band.B, band.U and appetite.hours that are
// open or that a Shape row still asks.
func intentWait(questions [][]string) []string { return nil }

// shapeQuestions is the first stage of the step: it asks the questions left open, records the answers, and gives the questions
// that the run waits for.
func (r *state) shapeQuestions(ctx context.Context) ([]string, error) { return nil, nil }
