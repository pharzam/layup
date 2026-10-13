package run

// This file holds the answers of the setup (docs/spec/scaffold.md, The answers
// of the setup; The work area; row 54 of the plan, task T-eahr, #277): the
// rows of inputs/answers.tsv that layup run writes from the answers of Intake.

import (
	"slices"
	"strconv"
	"strings"

	workpkg "github.com/pharzam/layup/internal/work"
)

// setupData is what the answers of the setup are made from: the rows of
// questions.tsv, intake/answers.tsv and intake/markers.tsv as their readers
// give them; URL gives the address of a comment of the Intake issue by its ID,
// and Batch the address of the comment of the batch.
type setupData struct {
	Questions, Answers, Markers [][]string
	URL                         func(comment string) string
	Batch                       string
}

// Columns of the rows that setupAnswers reads.
const (
	qID, qSource, qRef, qBlock, qText, qState = 0, 1, 2, 3, 8, 12 // questions.tsv
	aQuestion, aAnswer, aComment              = 0, 1, 2           // intake/answers.tsv
	mID, mValue                               = 0, 4              // intake/markers.tsv
)

// setupAnswers gives the rows of inputs/answers.tsv before the first run of the
// step runner: the row of S01-stack when it has a value, then one row for each
// answer to a question that is not a setup question, in the order of the
// question IDs (decided here, task T-eahr: the follow-up of a setup question is
// no such question, and the follow-up of another is its own); the questions of
// those rows, each with its text, which S04's raw fact names; and, by the ID of
// the marker, the row that a stop of S10 gets for each marker that the batch
// asked, by the table of "The answers of the setup".
//
// The answer of a setup question is that of the follow-up that asks it again,
// when that has one; else its own, unless its state is open, which is no answer.
// accept and gap are the whole answer, in lower case; any other text is a value,
// byte for byte. No answer, and accept to a question with no proposal (a stack
// has none), is no row for S01-stack, and for a marker a gap, at the comment of
// the batch; a gap that was answered is at the comment of its answer.
func setupAnswers(d setupData) (rows [][]string, asked []workpkg.Question, markers map[string][]string) {
	question := map[string][]string{}
	for _, q := range d.Questions {
		question[q[qID]] = q
	}
	answer := map[string][]string{}
	for _, a := range d.Answers {
		answer[a[aQuestion]] = a
	}
	follow := map[string]string{} // a question asked again, to the follow-up that asks it
	for _, q := range d.Questions {
		if q[qSource] == "follow-up" {
			follow[q[qRef]] = q[qID]
		}
	}
	setup := func(q []string) bool { return q != nil && q[qSource] == "setup" }
	effective := func(id string) ([]string, bool) {
		if a, ok := answer[follow[id]]; ok {
			return a, true
		}
		if question[id][qState] == "open" {
			return nil, false
		}
		a, ok := answer[id]
		return a, ok
	}
	proposal := map[string]string{}
	for _, m := range d.Markers {
		proposal[m[mID]] = m[mValue]
	}
	markers = map[string][]string{}
	for _, q := range d.Questions {
		if !setup(q) {
			continue
		}
		a, ok := effective(q[qID])
		if q[qRef] == "S01-stack" {
			if ok && a[aAnswer] != "accept" {
				rows = append(rows, []string{"S01-stack", a[aAnswer], q[qBlock], d.URL(a[aComment]), ""})
			}
			continue
		}
		gap := func(source string) []string { return []string{q[qRef], "gap", q[qBlock], source, q[qText]} }
		switch {
		case !ok || (a[aAnswer] == "accept" && proposal[q[qRef]] == ""):
			markers[q[qRef]] = gap(d.Batch)
		case a[aAnswer] == "gap":
			markers[q[qRef]] = gap(d.URL(a[aComment]))
		case a[aAnswer] == "accept":
			markers[q[qRef]] = []string{q[qRef], proposal[q[qRef]], q[qBlock], d.URL(a[aComment]), ""}
		default:
			markers[q[qRef]] = []string{q[qRef], a[aAnswer], q[qBlock], d.URL(a[aComment]), ""}
		}
	}
	var qs [][]string // the answers that are not the setup's, by the number of the question
	for _, a := range d.Answers {
		q := question[a[aQuestion]]
		if q != nil && !setup(q) && !(q[qSource] == "follow-up" && setup(question[q[qRef]])) {
			qs = append(qs, a)
		}
	}
	slices.SortStableFunc(qs, func(x, y []string) int { return questionNumber(x[aQuestion]) - questionNumber(y[aQuestion]) })
	for _, a := range qs {
		q := question[a[aQuestion]]
		rows = append(rows, []string{a[aQuestion], a[aAnswer], q[qBlock], d.URL(a[aComment]), ""})
		asked = append(asked, workpkg.Question{ID: a[aQuestion], Text: q[qText]})
	}
	return rows, asked, markers
}

// questionNumber is the number of a question ID Q-NNN, so Q-999 comes before Q-1000.
func questionNumber(id string) int {
	n, _ := strconv.Atoi(strings.TrimPrefix(id, "Q-"))
	return n
}
