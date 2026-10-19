package run

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// The questions left open of the step shape (task T-3f94, #364; docs/spec/loop.md, The questions left open).

// qr is a row of questions.tsv: the intent form, a follow-up and the others differ only in their source, ref, kind and state.
func qr(id, source, ref, kind, state string) []string {
	answered := ""
	if state == "answered" || state == "closed" {
		answered = "2026-10-10T13:00:00Z"
	}
	return []string{id, source, ref, "idea-owner", "before-delivery", kind, "", "", "text of " + id, "2026-10-10T12:00:00Z", "layup run", answered, state}
}

func TestTheQuestionsLeftOpenAreTheOpenAndTheUnansweredAsked(t *testing.T) {
	questions := [][]string{
		qr("Q-001", "intent", "band.B", "intent", "open"),          // asked
		qr("Q-002", "meaning", "M-001", "two-readings", "asked"),   // asked: no follow-up
		qr("Q-003", "meaning", "M-002", "two-readings", "asked"),   // its follow-up is answered
		qr("Q-004", "coverage", "C-001", "uncovered", "closed"),    // closed
		qr("Q-005", "rule", "Q-005", "G1", "asked"),                // a Shape row asks it
		qr("Q-006", "intent", "band.U", "intent", "open"),          // asked again: Shape's first row is answered
		qr("Q-007", "follow-up", "Q-003", "follow-up", "answered"), // Intake's follow-up
		qr("Q-008", "follow-up", "Q-005", "shape", "asked"),
		qr("Q-009", "follow-up", "Q-006", "shape", "answered"),
		qr("Q-010", "intent", "appetite.hours", "intent", "answered"),
	}
	if got, want := leftOpen(questions), []string{"Q-001", "Q-002", "Q-006"}; !slices.Equal(got, want) {
		t.Errorf("leftOpen = %v; want %v (the open and the unanswered asked; not the closed, the answered, the follow-ups, or one that Shape asks)", got, want)
	}
}

func TestAShapeRowAsksAgainInTheBlockOfItsQuestion(t *testing.T) {
	questions := [][]string{qr("Q-001", "intent", "band.B", "intent", "open"), qr("Q-002", "meaning", "M-001", "two-readings", "asked")}
	questions[1][3] = "operator"
	rows := shapeRows(questions, []string{"Q-002", "Q-001"}, "2026-10-12T09:00:00Z")
	want := [][]string{
		{"Q-003", "follow-up", "Q-002", "operator", "before-delivery", "shape", "", "", "text of Q-002", "2026-10-12T09:00:00Z", "layup run", "", "asked"},
		{"Q-004", "follow-up", "Q-001", "idea-owner", "before-delivery", "shape", "", "", "text of Q-001", "2026-10-12T09:00:00Z", "layup run", "", "asked"},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Fatalf("shapeRows = %v; want %v", rows, want)
	}
	for _, r := range rows {
		if err := records.CheckQuestion(r); err != nil {
			t.Errorf("the row %s: %v", r[0], err)
		}
	}
	if err := records.CheckQuestions(append(slices.Clone(questions), rows...)); err != nil {
		t.Errorf("the rows with the questions: %v", err)
	}
}

func TestTheCommentOfTheRoundIsTheFormOfTheBatch(t *testing.T) {
	questions := [][]string{qr("Q-001", "intent", "band.B", "intent", "open")}
	rows := shapeRows(questions, []string{"Q-001"}, "2026-10-12T09:00:00Z")
	got := shapeComment(questions, rows, []byte("the problem"), nil)
	want := strings.Replace(followUpComment(questions, rows, []byte("the problem"), nil), followUpHeader, shapeQuestionsHeader, 1)
	if got != want || !strings.HasPrefix(got, "LAYUP Shape: the questions left open\n\n## For the idea owner\n\nQ-002: text of Q-001\n") {
		t.Errorf("shapeComment = %q; want the comment of the follow-up under the header of Shape: %q", got, want)
	}
}

func TestOnlyTheRowsOfIntakeAreItsFollowUp(t *testing.T) {
	questions := [][]string{qr("Q-001", "intent", "band.B", "intent", "open"), qr("Q-002", "follow-up", "Q-001", "follow-up", "answered"), qr("Q-003", "follow-up", "Q-001", "shape", "asked")}
	if got := intakeFollowUps(questions); len(got) != 1 || got[0][0] != "Q-002" {
		t.Errorf("intakeFollowUps = %v; want the row Q-002 alone", got)
	}
}

// ans is a row of answers.tsv to a question, at a time.
func ans(id, text string) []string { return []string{id, text, "11", "1", "2026-10-12T10:00:00Z"} }

func states(rows [][]string) map[string]string {
	out := map[string]string{}
	for _, r := range rows {
		out[r[0]] = r[12]
	}
	return out
}

// intentRound are the four of the batch (open after the follow-up, answers that do not parse) and Shape's rows for the first three.
func intentRound() ([][]string, [][]string) {
	q := [][]string{qr("Q-001", "intent", "band.B", "intent", "open"), qr("Q-002", "intent", "band.U", "intent", "open"),
		qr("Q-003", "intent", "appetite.hours", "intent", "open"), qr("Q-004", "intent", "success-criteria", "intent", "answered"),
		qr("Q-005", "follow-up", "Q-001", "shape", "asked"), qr("Q-006", "follow-up", "Q-002", "shape", "asked"), qr("Q-007", "follow-up", "Q-003", "shape", "asked")}
	return q, [][]string{ans("Q-001", "soon"), ans("Q-002", "later"), ans("Q-003", "x"), ans("Q-004", "it works")}
}

// The answer to Shape's question sets the state of both questions, and writes budget.tsv with Shape's question as the source of each
// value, or replaces a value there only while the three still parse with B below U (plan review condition 6).
func TestAnIntentAnswerThatParsesWritesBudgetWithShapesQuestionAsSource(t *testing.T) {
	q, a := intentRound()
	got := [][]string{ans("Q-005", "5.0"), ans("Q-006", "10.0"), ans("Q-007", "40.0")}
	updated, budget := applyShapeAnswers(q, a, nil, got)
	if want := [][]string{{"band.B", "5.0", "Q-005"}, {"band.U", "10.0", "Q-006"}, {"appetite.hours", "40.0", "Q-007"}}; !reflect.DeepEqual(budget, want) {
		t.Errorf("budget missing: %v; want %v", budget, want)
	}
	for id, want := range map[string]string{"Q-001": "answered", "Q-002": "answered", "Q-003": "answered", "Q-005": "answered", "Q-006": "answered", "Q-007": "answered", "Q-004": "answered"} {
		if got := states(updated)[id]; got != want {
			t.Errorf("budget missing: the state of %s = %s; want %s", id, got, want)
		}
	}
	if updated[0][11] != "2026-10-12T10:00:00Z" || updated[4][11] != "2026-10-12T10:00:00Z" {
		t.Errorf("the times of Q-001 and Q-005 = %q, %q; want the time of the answer", updated[0][11], updated[4][11])
	}

	// A value that exists is replaced by the last Shape row of its field, and only with B below U.
	held := [][]string{{"band.B", "5.0", "Q-005"}, {"band.U", "10.0", "Q-006"}, {"appetite.hours", "40.0", "Q-007"}}
	q2 := append(slices.Clone(q), qr("Q-008", "follow-up", "Q-002", "shape", "asked"))
	for i := range q2[:7] {
		if i != 3 {
			q2[i][12] = "answered"
		}
	}
	a2 := append(slices.Clone(a), ans("Q-005", "5.0"), ans("Q-006", "10.0"), ans("Q-007", "40.0"))
	updated, budget = applyShapeAnswers(q2, a2, held, [][]string{ans("Q-008", "12.0")})
	if want := [][]string{{"band.B", "5.0", "Q-005"}, {"band.U", "12.0", "Q-008"}, {"appetite.hours", "40.0", "Q-007"}}; !reflect.DeepEqual(budget, want) {
		t.Errorf("a valid pair on a replace: %v; want %v", budget, want)
	}
	if states(updated)["Q-002"] != "answered" || states(updated)["Q-008"] != "answered" {
		t.Errorf("a valid pair on a replace: the states %v; want Q-002 and Q-008 answered", states(updated))
	}
	updated, budget = applyShapeAnswers(q2, a2, held, [][]string{ans("Q-008", "4.0")})
	if !reflect.DeepEqual(budget, held) || states(updated)["Q-002"] != "open" || states(updated)["Q-008"] != "answered" {
		t.Errorf("a B not below U on a replace: %v, %v; want budget as it was, Q-002 open, Q-008 answered", budget, states(updated))
	}
}

func TestAnUnparsedIntentAnswerLeavesTheQuestionOpenAndALaterCorrectionIsTheAnswerToANewRow(t *testing.T) {
	q, a := intentRound()
	got := [][]string{ans("Q-005", "soon"), ans("Q-006", "10.0"), ans("Q-007", "40.0")}
	updated, budget := applyShapeAnswers(q, a, nil, got)
	if budget != nil || states(updated)["Q-001"] != "open" || states(updated)["Q-005"] != "answered" || states(updated)["Q-002"] != "answered" {
		t.Fatalf("an unparsed B: %v, %v; want no budget, Q-001 open, Q-005 answered, Q-002 answered", budget, states(updated))
	}
	if want := []string{"Q-001"}; !slices.Equal(intentWait(updated), want) {
		t.Errorf("the wait = %v; want %v", intentWait(updated), want)
	}
	// The next round asks Q-001 again, in a new row, whose answer is the correction.
	rows := shapeRows(updated, leftOpen(updated), "2026-10-12T11:00:00Z")
	if len(rows) != 1 || rows[0][0] != "Q-008" || rows[0][2] != "Q-001" {
		t.Fatalf("the rows of the next round = %v; want Q-008 asking Q-001", rows)
	}
	all := append(slices.Clone(a), got...)
	updated, budget = applyShapeAnswers(append(updated, rows...), all, nil, [][]string{ans("Q-008", "5.0")})
	if want := [][]string{{"band.B", "5.0", "Q-008"}, {"band.U", "10.0", "Q-006"}, {"appetite.hours", "40.0", "Q-007"}}; !reflect.DeepEqual(budget, want) {
		t.Errorf("the correction: %v; want %v", budget, want)
	}
	if states(updated)["Q-001"] != "answered" || len(intentWait(updated)) != 0 {
		t.Errorf("the correction: the state of Q-001 is %s, the wait %v; want answered and no wait", states(updated)["Q-001"], intentWait(updated))
	}
}

// A partial answer records what was answered and leaves the rest asked.
func TestAPartialAnswerLeavesTheOtherQuestionsAsked(t *testing.T) {
	q, a := intentRound()
	updated, budget := applyShapeAnswers(q, a, nil, [][]string{ans("Q-005", "5.0")})
	want := map[string]string{"Q-001": "answered", "Q-002": "open", "Q-003": "open", "Q-005": "answered", "Q-006": "asked", "Q-007": "asked"}
	for id, w := range want {
		if states(updated)[id] != w {
			t.Errorf("the state of %s = %s; want %s", id, states(updated)[id], w)
		}
	}
	if budget != nil {
		t.Errorf("budget = %v; want none while two fields are unanswered", budget)
	}
	if wait := intentWait(updated); !slices.Equal(wait, []string{"Q-002", "Q-003"}) {
		t.Errorf("the wait = %v; want [Q-002 Q-003]", wait)
	}
}

func TestOnlyTheBudgetFieldsBlockAndAShapeRowStillAskedBlocks(t *testing.T) {
	q := [][]string{qr("Q-001", "intent", "band.B", "intent", "open"), qr("Q-002", "intent", "band.U", "intent", "answered"), qr("Q-003", "intent", "appetite.hours", "intent", "answered"),
		qr("Q-004", "intent", "success-criteria", "intent", "open"), qr("Q-005", "meaning", "M-001", "two-readings", "open"), qr("Q-006", "follow-up", "Q-002", "shape", "asked")}
	if got := intentWait(q); !slices.Equal(got, []string{"Q-001", "Q-002"}) {
		t.Errorf("intentWait = %v; want [Q-001 Q-002]: an open field, and one that a Shape row asks; not success-criteria, not a question of another source", got)
	}
}
