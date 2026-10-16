package run

// This file holds the route of a question to its owner role and the checks of
// the question and answer rows of a result (docs/spec/loop.md, The checks of a
// handoff, "A question during the work"; task T-q5k4, #360, row 66d).

import "strings"

// A Question is a question row of a typed result: its text, the label of its
// kind, its needs_human word and the records it concerns.
type Question struct {
	Text, Label, NeedsHuman string
	Records                 []string
}

// A Route is where a question goes: the owner role of its kind, or a human.
type Route struct {
	Role  string
	Human bool
}

// DefaultOwners gives the default owner map of architecture.md §9 by the label
// of a kind, a map of its own for each call; the confirmed map is phase 3's
// (owners.tsv).
func DefaultOwners() map[string]string {
	return map[string]string{
		"domain-or-business":            "Domain Expert",
		"architecture-boundary":         "Systems Architect",
		"interface-contract":            "Software Architect",
		"environment-or-infrastructure": "Software Engineer",
	}
}

// QuestionsOf gives the question rows of a result, in the order of the file:
// the reason of a row is the label, the needs_human word and the records it
// concerns, one space between (records.md, the block result).
func QuestionsOf(result [][]string) []Question {
	var out []Question
	for _, r := range result {
		if r[0] != "question" {
			continue
		}
		w := strings.Fields(r[4])
		q := Question{Text: r[2]}
		if len(w) > 0 {
			q.Label = w[0]
		}
		if len(w) > 1 {
			q.NeedsHuman = w[1]
		}
		if len(w) > 2 {
			q.Records = w[2:]
		}
		out = append(out, q)
	}
	return out
}

// RouteQuestion routes a question: to the owner role of its kind when it needs
// no human (needs_human is exactly "no") and its label is a kind of the owner
// map; otherwise to the Operator, a word other than yes or no and a label
// outside the map included (loop.md, "A question during the work": the
// conservative reading; architecture.md §10 allocates the Operator to a point
// that no clause gives to the idea owner).
func RouteQuestion(q Question, owners map[string]string) Route {
	if role, ok := owners[q.Label]; ok && q.NeedsHuman == "no" {
		return Route{Role: role}
	}
	return Route{Role: "operator", Human: true}
}

// QuestionRowOK is the item question-row: the result has at least one question
// row, and each names at least one record it concerns after its label and its
// needs_human word. It does not resolve a record: no rule defines resolution.
func QuestionRowOK(result [][]string) bool {
	qs := QuestionsOf(result)
	for _, q := range qs {
		if len(q.Records) == 0 {
			return false
		}
	}
	return len(qs) > 0
}

// AnswerRowOK is the item answer-row: the result has at least one answer row,
// and each cites at least one record it rests on (the block result refuses an
// answer with none; the check does not rest on that alone).
func AnswerRowOK(result [][]string) bool {
	n := 0
	for _, r := range result {
		if r[0] != "answer" {
			continue
		}
		if len(strings.Fields(r[4])) == 0 {
			return false
		}
		n++
	}
	return n > 0
}
