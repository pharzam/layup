package run

import (
	"strings"
	"testing"
)

// fuQuestions is a batch with the intent form (Q-002 to Q-004) and a
// follow-up row for each field asked again.
func fuQuestions(asked ...string) [][]string {
	q := func(id, source, ref, kind, text, asker, state string) []string {
		return []string{id, source, ref, "idea-owner", "before-delivery", kind, "", "", text, "2026-10-11T10:00:00Z", asker, "", state}
	}
	rows := [][]string{
		q("Q-001", "rule", "Q-001", "G1", "Which stack?", "layup psb check", "closed"),
		q("Q-002", "intent", "band.B", "intent", "B?", "layup run", "open"),
		q("Q-003", "intent", "band.U", "intent", "U?", "layup run", "open"),
		q("Q-004", "intent", "appetite.hours", "intent", "Hours?", "layup run", "open"),
	}
	for i, ref := range asked {
		rows = append(rows, q("Q-00"+string(rune('5'+i)), "follow-up", ref, "follow-up", "again", "layup run", "asked"))
	}
	return rows
}

func fuAnswers(pairs ...string) [][]string {
	var out [][]string
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, []string{pairs[i], pairs[i+1], "802", "303", "2026-10-11T10:05:00Z"})
	}
	return out
}

func joined(rows [][]string) string {
	var out []string
	for _, r := range rows {
		out = append(out, strings.Join(r, " "))
	}
	return strings.Join(out, "; ")
}

// The budget after the follow-up (docs/spec/intake.md, The answers to the
// follow-up): not yet written, each field from the answer to its follow-up,
// else from the batch; written already, a parsed answer to a field asked again
// replaces its value; a field still missing or unparsed leaves its question
// open.
func TestBudgetAfterFollowUp(t *testing.T) {
	written := [][]string{{"band.B", "5.0", "Q-002"}, {"band.U", "10.0", "Q-003"}, {"appetite.hours", "40.0", "Q-004"}}
	for _, c := range []struct {
		name      string
		questions [][]string
		answers   [][]string
		budget    [][]string
		want      string
		open      string
	}{
		{"not yet written, B from the follow-up", fuQuestions("Q-002", "Q-003"),
			fuAnswers("Q-002", "10.0", "Q-003", "9.5", "Q-004", "40.0", "Q-005", "5.0", "Q-006", "10.0"), nil,
			"band.B 5.0 Q-005; band.U 10.0 Q-006; appetite.hours 40.0 Q-004", ""},
		{"not yet written, the hours still unparsed", fuQuestions("Q-004"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0", "Q-004", "forty", "Q-005", "forty again"), nil, "", "Q-004"},
		{"not yet written, B not below U after the follow-up", fuQuestions("Q-002", "Q-003"),
			fuAnswers("Q-002", "10.0", "Q-003", "9.5", "Q-004", "40.0", "Q-005", "12.0", "Q-006", "11.0"), nil, "", "Q-002 Q-003"},
		// The pair rule holds whatever the state of the hours (round 1 of #225).
		{"not yet written, B not below U and the hours unparsed", fuQuestions("Q-002", "Q-003", "Q-004"),
			fuAnswers("Q-002", "10.0", "Q-003", "9.5", "Q-004", "forty", "Q-005", "12.0", "Q-006", "11.0", "Q-007", "forty again"), nil, "", "Q-002 Q-003 Q-004"},
		{"not yet written, a field asked again with no answer", fuQuestions("Q-004"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0"), nil, "", "Q-004"},
		{"written, the hours asked again and replaced", fuQuestions("Q-004"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0", "Q-004", "40.0", "Q-005", "50.0"), written,
			"band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 50.0 Q-005", ""},
		{"written, an answer that does not parse", fuQuestions("Q-004"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0", "Q-004", "40.0", "Q-005", "fifty"), written,
			"band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", "Q-004"},
		{"written, no field of the intent form asked again", fuQuestions("Q-001"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0", "Q-004", "40.0", "Q-005", "Go"), written,
			"band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", ""},
		// A budget.tsv with no B or U (no writer of LAYUP gives one) is kept as
		// it was, and the field asked again is open.
		{"written, a budget of one row", fuQuestions("Q-004"),
			fuAnswers("Q-005", "50.0"), written[2:], "appetite.hours 40.0 Q-004", "Q-004"},
		{"written, a replacement that leaves B not below U", fuQuestions("Q-003"),
			fuAnswers("Q-002", "5.0", "Q-003", "10.0", "Q-004", "40.0", "Q-005", "4.0"), written,
			"band.B 5.0 Q-002; band.U 10.0 Q-003; appetite.hours 40.0 Q-004", "Q-003"},
	} {
		rows, open := budgetAfterFollowUp(c.questions, c.answers, c.budget)
		if got := joined(rows); got != c.want || strings.Join(open, " ") != c.open {
			t.Errorf("%s: %q, open %q; want %q, open %q", c.name, got, open, c.want, c.open)
		}
	}
}
