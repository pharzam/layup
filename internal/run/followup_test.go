package run

import (
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	workpkg "github.com/pharzam/layup/internal/work"
)

// followQuestions is a batch after the answers: Q-001 a rule gap with a quote,
// answered; Q-002 a gap of meaning, answered; Q-003 a rule gap with no answer;
// and the intent form, answered.
func followQuestions(state1, state2, b, u, h string) [][]string {
	q := func(id, source, ref, kind, start, end, text, asker, answered, state string) []string {
		if state == "asked" {
			answered = "" // an asked question has no answer
		}
		return []string{id, source, ref, "idea-owner", "before-delivery", kind, start, end, text, "2026-10-11T10:00:00Z", asker, answered, state}
	}
	at := "2026-10-11T10:05:00Z"
	return [][]string{
		q("Q-001", "rule", "Q-001", "G3", "4", "30", "What is TLA?", "layup psb check", at, state1),
		q("Q-002", "meaning", "M-001", "unmeasurable", "31", "42", "How fast?", "S-11111111", at, state2),
		q("Q-003", "rule", "Q-002", "G1", "", "", "Which stack?", "layup psb check", "", "asked"),
		q("Q-004", "intent", "band.B", "intent", "", "", "B?", "layup run", at, b),
		q("Q-005", "intent", "band.U", "intent", "", "", "U?", "layup run", at, u),
		q("Q-006", "intent", "appetite.hours", "intent", "", "", "Hours?", "layup run", at, h),
	}
}

func followAnswers(b, u, h string) [][]string {
	a := func(id, answer string) []string { return []string{id, answer, "802", "303", "2026-10-11T10:05:00Z"} }
	return [][]string{a("Q-001", "Three Letters"), a("Q-002", "fast enough"), a("Q-004", b), a("Q-005", u), a("Q-006", h)}
}

// The Keep of the review of the answers: each answered question once, closed
// or open with a reason; the record in the order of the questions with the
// session, and questions.tsv with the new states.
func TestKeepOfTheReview(t *testing.T) {
	questions := followQuestions("answered", "answered", "answered", "answered", "answered")
	keep := reviewKeep(questions)
	head := "question\tstate\treason\n"
	files, err := keep("S-0a0b0c0d", []byte(head+"Q-002\topen\tfast is no measure\nQ-001\tclosed\tclear\nQ-004\tclosed\tok\nQ-005\tclosed\tok\nQ-006\tclosed\tok\n"))
	if err != nil {
		t.Fatalf("Keep: %v", err)
	}
	if got := string(files["intake/answers-review.tsv"]); !strings.HasPrefix(got, "question\tstate\treason\tsession\nQ-001\tclosed\tclear\tS-0a0b0c0d\nQ-002\topen\tfast is no measure\tS-0a0b0c0d\n") {
		t.Errorf("intake/answers-review.tsv:\n%s", got)
	}
	rows, err := records.ReadQuestions(files["questions.tsv"])
	if err != nil {
		t.Fatal(err)
	}
	if states := []string{rows[0][12], rows[1][12], rows[2][12]}; !slices.Equal(states, []string{"closed", "open", "asked"}) {
		t.Errorf("the states %q; want closed, open, and asked for the question with no answer", states)
	}
	for _, c := range []struct{ name, file, refused string }{
		{"a question with no answer", head + "Q-003\tclosed\tx\n", "Q-003 is not an answered question"},
		{"a question twice", head + "Q-001\tclosed\tx\nQ-001\topen\ty\nQ-002\tclosed\tx\nQ-004\tclosed\tx\nQ-005\tclosed\tx\nQ-006\tclosed\tx\n", "repeats the key"},
		{"an answered question with no row", head + "Q-001\tclosed\tx\n", "Q-002 has no row"},
		{"a file of the record's form", "question\tstate\treason\tsession\nQ-001\tclosed\tx\tS-0a0b0c0d\n", "header"},
	} {
		if files, err := keep("S-0a0b0c0d", []byte(c.file)); err == nil || files != nil || !strings.Contains(err.Error(), c.refused) {
			t.Errorf("%s: %q, %v; want it refused, %q", c.name, files, err, c.refused)
		}
	}
	// With no answered question the review writes the record with no row.
	if got := string(emptyReview()); got != "question\tstate\treason\tsession\n" {
		t.Errorf("the record of no review %q", got)
	}
}

// The follow-up asks again each question still asked, each open one, and each
// field of the intent form that does not parse (a B not below U is an
// unparsed answer of both); one row per question.
// idsJoin joins a follow-up's rows as "Q-NNN<ref" pairs.
func idsJoin(rows [][]string) string {
	var out []string
	for _, r := range rows {
		out = append(out, r[0]+"<"+r[2])
	}
	return strings.Join(out, " ")
}

func TestFollowUpRows(t *testing.T) {
	ids := idsJoin
	for _, c := range []struct {
		name      string
		questions [][]string
		answers   [][]string
		want      string
	}{
		{"missing, open, B not below U", followQuestions("closed", "open", "closed", "closed", "closed"), followAnswers("10.0", "9.5", "40.0"),
			"Q-007<Q-002 Q-008<Q-003 Q-009<Q-004 Q-010<Q-005"},
		{"an unparsed hours field", followQuestions("closed", "closed", "closed", "closed", "closed"), followAnswers("9.5", "10.0", "forty"),
			"Q-007<Q-003 Q-008<Q-006"},
		{"an open intent field, its value parsed (budget.tsv written)", followQuestions("closed", "closed", "closed", "closed", "open"), followAnswers("9.5", "10.0", "40.0"),
			"Q-007<Q-003 Q-008<Q-006"},
		{"an intent field with no answer, asked once", followQuestions("closed", "closed", "closed", "closed", "asked"), followAnswers("9.5", "10.0", "40.0")[:4],
			"Q-007<Q-003 Q-008<Q-006"},
	} {
		rows := followUpRows(c.questions, c.answers, "2026-10-11T11:00:00Z")
		if got := ids(rows); got != c.want {
			t.Errorf("%s: %s; want %s", c.name, got, c.want)
		}
		for _, r := range rows {
			if r[1] != "follow-up" || r[5] != "follow-up" || r[6] != "" || r[10] != "layup run" || r[12] != "asked" || r[9] != "2026-10-11T11:00:00Z" {
				t.Errorf("%s: the row %q is no follow-up row", c.name, r)
			}
		}
		table, _ := table(records.QuestionsSchema, append(slices.Clone(c.questions), rows...))
		if _, err := records.ReadQuestions(table); err != nil {
			t.Errorf("%s: ReadQuestions: %v", c.name, err)
		}
	}
	all := followQuestions("closed", "closed", "closed", "closed", "closed")
	all[2][11], all[2][12] = "2026-10-11T10:05:00Z", "closed"
	if rows := followUpRows(all, followAnswers("9.5", "10.0", "40.0"), "t"); rows != nil {
		t.Errorf("nothing to ask: %q; want no row", rows)
	}
}

// A field of the intent form that is no name of the budget — the success
// criteria — is not asked again for its answer parsing no decimal; an open or
// a still-asked one is asked like any question.
func TestTheCriteriaFieldIsAskedOnlyByState(t *testing.T) {
	criteria := func(state, answered string) []string {
		if state == "asked" {
			answered = ""
		}
		return []string{"Q-007", "intent", "success-criteria", "idea-owner", "before-delivery", "intent", "", "", "The criteria?", "2026-10-11T10:00:00Z", "layup run", answered, state}
	}
	answers := append(slices.Clone(followAnswers("9.5", "10.0", "40.0")),
		[]string{"Q-007", "fast; cheap", "802", "303", "2026-10-11T10:05:00Z"})
	for _, c := range []struct {
		name, state, answered, want string
	}{
		{"answered, its words no decimal", "closed", "2026-10-11T10:05:00Z", "Q-008<Q-003"},
		{"open, asked again by its state", "open", "2026-10-11T10:05:00Z", "Q-008<Q-003 Q-009<Q-007"},
		{"still asked, asked again", "asked", "", "Q-008<Q-003 Q-009<Q-007"},
	} {
		questions := append(slices.Clone(followQuestions("closed", "closed", "closed", "closed", "closed")), criteria(c.state, c.answered))
		if rows := followUpRows(questions, answers, "2026-10-11T11:00:00Z"); idsJoin(rows) != c.want {
			t.Errorf("%s: %v; want %s", c.name, rows, c.want)
		}
	}
}

// The comment of the follow-up: its own header, the two blocks, each question
// with the quote of the question it asks again.
func TestTheCommentOfTheFollowUp(t *testing.T) {
	questions := followQuestions("open", "closed", "closed", "closed", "closed")
	rows := followUpRows(questions, followAnswers("9.5", "10.0", "40.0"), "t")
	c := followUpComment(questions, rows, []byte(batchPSB), nil)
	for _, want := range []string{followUpHeader + "\n", "\n## For the idea owner\n", "\nQ-007: What is TLA?\n`We use Go and TLA and XYZ.`\n",
		"\nQ-008: Which stack?\n\n", "\n## For the Operator\n\nNo question for the Operator.\n", "`Q-NNN: <answer>`"} {
		if !strings.Contains(c, want) {
			t.Errorf("the comment has no %q:\n%s", want, c)
		}
	}
	if strings.HasPrefix(c, batchHeader) {
		t.Error("the follow-up starts with the header of the batch")
	}
}

// A follow-up of a setup question is in the Operator's block, and shows what the
// batch showed: S01-stack asked again, and a marker asked again with the line
// proposed: and the quote of its proposal; the rule gap is in the idea owner's
// block (task T-4c5c, #276; docs/spec/intake.md, The review of the answers and
// the follow-up).
func TestTheFollowUpOfASetupQuestion(t *testing.T) {
	questions := followQuestions("closed", "closed", "closed", "closed", "closed")
	setup := func(id, ref, start, end, text, asker, answered, state string) []string {
		return []string{id, "setup", ref, "operator", "before-delivery", "setup", start, end, text, "2026-10-11T10:00:00Z", asker, answered, state}
	}
	marker := "the value of " + angled("owner") + " at docs/a.md:2"
	questions = append(questions,
		setup("Q-007", "S01-stack", "", "", workpkg.S01Questions[0].Text, "layup run", "", "asked"),
		setup("Q-008", "M-0a1b2c3d", "4", "30", marker, "S-33333333", "2026-10-11T10:05:00Z", "open"))
	rows := followUpRows(questions, followAnswers("9.5", "10.0", "40.0"), "t")
	if len(rows) != 3 || rows[1][3] != "operator" || rows[2][3] != "operator" || rows[0][3] != "idea-owner" {
		t.Fatalf("the follow-up rows %q; want the rule gap, S01-stack and the marker, each in the block of its question", rows)
	}
	c := followUpComment(questions, rows, []byte(batchPSB), batchFixture(t).Markers)
	pointForm := "Name the approver of each planned point, one line `<point>: <login>` for each you name (a point you do not name is yours):\n" +
		"- `setup`\n- `bet`\n- `activation`\n- `retrospective`\n- `early-retrospective`\n- `audit`\n"
	want := "\n## For the idea owner\n\nQ-009: Which stack?\n\n" + pointForm + "\n## For the Operator\n\nQ-010: " + workpkg.S01Questions[0].Text + "\n\n" +
		"Q-011: " + marker + "\nproposed: Carol\n`We use Go and TLA and XYZ.`\n\n"
	if !strings.Contains(c, want) {
		t.Errorf("the comment has no %q:\n%s", want, c)
	}
	if !strings.Contains(c, "`Q-NNN: accept`") {
		t.Errorf("the comment asks a marker and has no forms of its answer:\n%s", c)
	}
}
