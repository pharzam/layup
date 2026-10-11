package run

import (
	"strings"
	"testing"
)

// answersFixture is a batch of four questions of the idea owner's block, asked
// at 10:00:00: a rule gap and the three fields of the intent form.
var answersQuestions = [][]string{
	{"Q-001", "rule", "Q-001", "idea-owner", "before-delivery", "G1", "", "", "Which stack?", "2026-10-11T10:00:00Z", "layup psb check", "", "asked"},
	{"Q-002", "intent", "band.B", "idea-owner", "before-delivery", "intent", "", "", "B?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
	{"Q-003", "intent", "band.U", "idea-owner", "before-delivery", "intent", "", "", "U?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
	{"Q-004", "intent", "appetite.hours", "idea-owner", "before-delivery", "intent", "", "", "Hours?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
}

var answersApprovers = [][]string{
	{"303", "idea-owner", "carol", "2026-10-08T12:00:00Z", "start"},
	{"101", "operator", "pharzam", "2026-10-08T12:00:00Z", "start"},
}

// copyOf gives a copy of a comment of issue 1: a row of copies.tsv and its body.
func copyOf(id, seen, author, app, created, body string) answerCopy {
	return answerCopy{Row: []string{id, seen, "1", author, "x", app, created, created, strings.Repeat("a", 64), "copies/" + id + "-" + seen + ".md"}, Body: body}
}

// The answers (docs/spec/intake.md, The answers): a line Q-NNN: <answer> of a
// first copy after the batch by an approver in the role of the question's
// block answers it; the last such line of an ID is its answer; a block waits
// until such a comment holds at least one answer line.
func TestCollectAnswers(t *testing.T) {
	ignored := []answerCopy{
		copyOf("701", "1", "303", "", "2026-10-11T09:59:59Z", "Q-001: before the batch\n"),
		copyOf("702", "1", "303", "", "2026-10-11T10:00:00Z", "Q-001: at the time of the batch\n"),
		copyOf("703", "1", "101", "", "2026-10-11T10:01:00Z", "Q-001: the Operator, in the other block's role\n"),
		copyOf("704", "1", "999", "", "2026-10-11T10:02:00Z", "Q-001: not an approver\n"),
		copyOf("705", "1", "303", "layup-agent", "2026-10-11T10:03:00Z", "Q-001: by an App\n"),
		copyOf("706", "1", "303", "", "2026-10-11T10:04:00Z", "Thank you, I will answer soon.\n"),
	}
	if answers, waiting := collectAnswers(answersQuestions, answersApprovers, ignored, "1"); len(answers) != 0 || strings.Join(waiting, " ") != "idea-owner" {
		t.Fatalf("no answer line of an approver in the role: %q, waiting %q; want none, the idea owner's block waiting", answers, waiting)
	}
	copies := append(ignored,
		copyOf("707", "1", "303", "", "2026-10-11T10:05:00Z",
			"Q-001: first\nQ-002: 9.5\r\nQ-003: 10.0 \nQ-004: 40.0\nQ-009: no such question\n xQ-001: not at the start\nQ-001:no space\n"),
		// A later line corrects an earlier one; an empty answer and one of
		// exactly — after it answer nothing.
		copyOf("708", "1", "303", "", "2026-10-11T10:06:00Z", "Q-001: second\nQ-001: \nQ-001: —\n"),
		copyOf("707", "2", "303", "", "2026-10-11T10:07:00Z", "Q-001: an edit, never an answer\n"),
	)
	answers, waiting := collectAnswers(answersQuestions, answersApprovers, copies, "1")
	want := [][]string{
		{"Q-001", "second", "708", "303", "2026-10-11T10:06:00Z"},
		{"Q-002", "9.5", "707", "303", "2026-10-11T10:05:00Z"},   // the CR before the line feed trimmed
		{"Q-003", "10.0 ", "707", "303", "2026-10-11T10:05:00Z"}, // trailing spaces kept
		{"Q-004", "40.0", "707", "303", "2026-10-11T10:05:00Z"},
	}
	if len(waiting) != 0 || len(answers) != len(want) {
		t.Fatalf("the answers %q, waiting %q; want %q", answers, waiting, want)
	}
	for i := range want {
		if strings.Join(answers[i], "|") != strings.Join(want[i], "|") {
			t.Errorf("answer %d: %q, want %q", i, answers[i], want[i])
		}
	}
}

// The intent form parses when each field is a decimal and B is below U, as
// numbers; then budget.tsv has its three rows, each with its question.
func TestBudgetOf(t *testing.T) {
	answer := func(b, u, h string) [][]string {
		return [][]string{{"Q-002", b, "707", "303", "t"}, {"Q-003", u, "707", "303", "t"}, {"Q-004", h, "707", "303", "t"}}
	}
	rows, ok := budgetOf(answersQuestions, answer("9.5", "10.0", "40.0"))
	if !ok || len(rows) != 3 || strings.Join(rows[0], " ") != "band.B 9.5 Q-002" || strings.Join(rows[1], " ") != "band.U 10.0 Q-003" ||
		strings.Join(rows[2], " ") != "appetite.hours 40.0 Q-004" {
		t.Errorf("9.5, 10.0, 40.0: %q, %v; want the three rows (9.5 below 10.0 as numbers)", rows, ok)
	}
	// Spaces around a value are trimmed by the parse; the answer itself is kept
	// byte for byte in answers.tsv.
	if rows, ok := budgetOf(answersQuestions, answer(" 9.5", "10.0 ", "40.0")); !ok || rows[0][1] != "9.5" || rows[1][1] != "10.0" {
		t.Errorf("spaces around the values: %q, %v; want them parsed", rows, ok)
	}
	for _, c := range []struct{ name, b, u, h string }{
		{"B not below U", "10.0", "9.5", "40.0"},
		{"B equal to U", "10.0", "10.0", "40.0"},
		{"a value that is no decimal", "9.5", "10.0", "40"},
		{"a value with a word", "9.5 USD", "10.0", "40.0"},
	} {
		if rows, ok := budgetOf(answersQuestions, answer(c.b, c.u, c.h)); ok || rows != nil {
			t.Errorf("%s: %q, %v; want budget.tsv not written", c.name, rows, ok)
		}
	}
	if rows, ok := budgetOf(answersQuestions, answer("9.5", "10.0", "40.0")[:2]); ok || rows != nil {
		t.Errorf("a field with no answer: %q, %v; want budget.tsv not written", rows, ok)
	}
}
