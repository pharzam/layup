package run

import (
	"reflect"
	"testing"

	workpkg "github.com/pharzam/layup/internal/work"
)

// The answers of the setup (docs/spec/scaffold.md, The answers of the setup;
// The work area; the acceptance test "The answers of the setup"; row 54 of the
// plan, task T-eahr, #277): the rows of inputs/answers.tsv from the rows of
// questions.tsv, intake/answers.tsv and intake/markers.tsv.

const issueURL = "https://github.com/acme/target/issues/1#issuecomment-"

func commentURL(comment string) string { return issueURL + comment }

// sq is a row of questions.tsv, of the form that records.ReadQuestions gives.
func sq(id, source, ref, block, text, state string) []string {
	return []string{id, source, ref, block, "before-delivery", "x", "", "", text, "2026-10-08T13:00:00Z", "layup run", "", state}
}

// sa is a row of intake/answers.tsv: the question, the answer, the comment, the author, the time.
func sa(question, answer, comment string) []string {
	return []string{question, answer, comment, "101", "2026-10-08T13:10:00Z"}
}

// setupWorld is the Intake of a target with three markers: the owner and the
// port with a proposal, the host with none; Q-001 is a question of the idea
// owner, and Q-007 a follow-up of it. Q-006 asks the owner marker again.
func setupWorld(answers ...[]string) setupData {
	return setupData{
		Questions: [][]string{
			sq("Q-001", "rule", "Q-001", "idea-owner", "Which database?", "answered"),
			sq("Q-002", "setup", "S01-stack", "operator", "Which technology stack?", "answered"),
			sq("Q-003", "setup", "M-aaaaaaaa", "operator", "the value of <owner name> at NOTES.md:2", "answered"),
			sq("Q-004", "setup", "M-bbbbbbbb", "operator", "the value of <port> at NOTES.md:3", "answered"),
			sq("Q-005", "setup", "M-cccccccc", "operator", "the value of <host> at docs/deploy.md:1", "asked"),
			sq("Q-006", "follow-up", "Q-003", "operator", "the value of <owner name> at NOTES.md:2", "asked"),
			sq("Q-007", "follow-up", "Q-001", "idea-owner", "Which database?", "asked"),
		},
		Answers: answers,
		Markers: [][]string{
			{"M-aaaaaaaa", "NOTES.md", "2", "<owner name>", "Carol Jones", "10", "20", "S-11112222"},
			{"M-bbbbbbbb", "NOTES.md", "3", "<port>", "8080", "30", "40", "S-11112222"},
			{"M-cccccccc", "docs/deploy.md", "1", "<host>", "", "", "", "S-11112222"},
		},
		URL:   commentURL,
		Batch: commentURL("601"),
	}
}

// with sets the state of a question of the world.
func (d setupData) with(question, state string) setupData {
	d.Questions = append([][]string(nil), d.Questions...)
	for i, q := range d.Questions {
		if q[0] == question {
			q = append([]string(nil), q...)
			q[12] = state
			d.Questions[i] = q
		}
	}
	return d
}

// The rows of the table of "The answers of the setup", for a marker with a
// proposal (M-aaaaaaaa, "Carol Jones"), one with none (M-cccccccc), and a gap
// with no answer, which gets the comment of the batch.
func TestTheAnswerOfAMarkerIsARowByTheTable(t *testing.T) {
	gapOf := func(id, text, source string) []string { return []string{id, "gap", "operator", source, text} }
	owner := "the value of <owner name> at NOTES.md:2"
	host := "the value of <host> at docs/deploy.md:1"
	for _, c := range []struct {
		name string
		data setupData
		id   string
		want []string
	}{
		{"accept, to a marker with a proposal: the proposed value", setupWorld(sa("Q-003", "accept", "804")), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", "Carol Jones", "operator", commentURL("804"), ""}},
		{"accept, to a marker with no proposal: no answer, so a gap at the batch's comment", setupWorld(sa("Q-005", "accept", "804")), "M-cccccccc",
			gapOf("M-cccccccc", host, commentURL("601"))},
		{"gap: a gap at the comment of the answer, with the question of the batch", setupWorld(sa("Q-003", "gap", "805")), "M-aaaaaaaa",
			gapOf("M-aaaaaaaa", owner, commentURL("805"))},
		{"a value, byte for byte", setupWorld(sa("Q-004", "9090 ", "806")), "M-bbbbbbbb",
			[]string{"M-bbbbbbbb", "9090 ", "operator", commentURL("806"), ""}},
		{"no answer: a gap at the batch's comment", setupWorld(), "M-cccccccc", gapOf("M-cccccccc", host, commentURL("601"))},
		{"open after the follow-up, with an answer: no answer", setupWorld(sa("Q-003", "Carol", "807")).with("Q-003", "open"), "M-aaaaaaaa",
			gapOf("M-aaaaaaaa", owner, commentURL("601"))},
		{"the answer to the follow-up replaces the answer, whatever the state of the question", setupWorld(sa("Q-003", "Carol", "807"), sa("Q-006", "Jane Roe", "808")).with("Q-003", "open"), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", "Jane Roe", "operator", commentURL("808"), ""}},
		{"the follow-up asks, nobody answers, and the question is not open: its answer stands", setupWorld(sa("Q-003", "Carol", "807")), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", "Carol", "operator", commentURL("807"), ""}},
		{"accept to the follow-up takes the proposal", setupWorld(sa("Q-006", "accept", "809")), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", "Carol Jones", "operator", commentURL("809"), ""}},
		{"only the word, in lower case, is an accept", setupWorld(sa("Q-003", "Accept", "810")), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", "Accept", "operator", commentURL("810"), ""}},
		{"only the word, in lower case, is a gap", setupWorld(sa("Q-004", "Gap", "812")), "M-bbbbbbbb",
			[]string{"M-bbbbbbbb", "Gap", "operator", commentURL("812"), ""}},
		{"a word with a space round it is a value", setupWorld(sa("Q-003", " gap", "811")), "M-aaaaaaaa",
			[]string{"M-aaaaaaaa", " gap", "operator", commentURL("811"), ""}},
	} {
		_, _, markers := setupAnswers(c.data)
		if got := markers[c.id]; !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: the row of %s is %q; want %q", c.name, c.id, got, c.want)
		}
	}
}

// Only the markers that the batch asked have a row to give, in each case a row.
func TestEachMarkerThatTheBatchAskedHasARow(t *testing.T) {
	_, _, markers := setupAnswers(setupWorld())
	if len(markers) != 3 || markers["M-aaaaaaaa"] == nil || markers["M-bbbbbbbb"] == nil || markers["M-cccccccc"] == nil {
		t.Errorf("the markers with a row: %q; want the three that the batch asked", markers)
	}
	d := setupWorld()
	d.Questions = d.Questions[:3] // the batch asked the stack and the owner only
	if _, _, markers := setupAnswers(d); len(markers) != 1 || markers["M-aaaaaaaa"] == nil || markers["M-bbbbbbbb"] != nil {
		t.Errorf("the markers with a row after a batch of one: %q; want M-aaaaaaaa alone, as a marker that the batch did not ask has no row", markers)
	}
}

// S01-stack gets a row for a value only: accept, no answer and a question
// that is open give none, so the runner stops on it and the step is fail.
func TestTheAnswerOfTheStackIsARowForAValueOnly(t *testing.T) {
	stack := func(answer, comment string) []string {
		return []string{"S01-stack", answer, "operator", commentURL(comment), ""}
	}
	for _, c := range []struct {
		name string
		data setupData
		want [][]string // the rows before Q-001's, which has no answer in these cases
	}{
		{"a value", setupWorld(sa("Q-002", "go", "804")), [][]string{stack("go", "804")}},
		{"accept: there is no value to accept", setupWorld(sa("Q-002", "accept", "804")), nil},
		{"no answer", setupWorld(), nil},
		{"open, with an answer", setupWorld(sa("Q-002", "go", "804")).with("Q-002", "open"), nil},
		{"the answer to the follow-up", func() setupData {
			d := setupWorld(sa("Q-002", "rust", "804"), sa("Q-008", "go", "805")).with("Q-002", "open")
			d.Questions = append(d.Questions, sq("Q-008", "follow-up", "Q-002", "operator", "Which technology stack?", "asked"))
			return d
		}(), [][]string{stack("go", "805")}},
	} {
		rows, _, _ := setupAnswers(c.data)
		if !reflect.DeepEqual(rows, c.want) {
			t.Errorf("%s: the rows %q; want %q", c.name, rows, c.want)
		}
	}
}

// The Q- rows: one for each answered question that is not a setup question and
// not the follow-up of one, in the order of the question IDs, by the role of its
// block, at the URL of its comment, with the text of the question for S04's raw
// fact; the row of S01-stack comes first.
func TestTheQRowsAreTheAnswersThatAreNotTheSetup(t *testing.T) {
	d := setupWorld(sa("Q-007", "MySQL", "813"), sa("Q-001", "PostgreSQL", "812"), sa("Q-002", "go", "804"), sa("Q-003", "accept", "805"), sa("Q-006", "Jane Roe", "806"))
	rows, asked, _ := setupAnswers(d)
	want := [][]string{
		{"S01-stack", "go", "operator", commentURL("804"), ""},
		{"Q-001", "PostgreSQL", "idea-owner", commentURL("812"), ""},
		{"Q-007", "MySQL", "idea-owner", commentURL("813"), ""},
	}
	if !reflect.DeepEqual(rows, want) {
		t.Errorf("the rows %q; want S01-stack, then Q-001 and Q-007 (the follow-up of a question that is not a setup question is a Q- row of its own), and no row of a marker or of its follow-up", rows)
	}
	if want := []workpkg.Question{{ID: "Q-001", Text: "Which database?"}, {ID: "Q-007", Text: "Which database?"}}; !reflect.DeepEqual(asked, want) {
		t.Errorf("the questions of the Q- rows %q; want %q", asked, want)
	}
	// A question with no answer has no row.
	if rows, asked, _ := setupAnswers(setupWorld(sa("Q-002", "go", "804"))); len(rows) != 1 || len(asked) != 0 {
		t.Errorf("with no answer to a question of the idea owner: the rows %q and the questions %q; want S01-stack alone", rows, asked)
	}
}

// The Q- rows are in the order of the numbers of the questions, not of their
// text: Q-999 comes before Q-1000, whatever the order of intake/answers.tsv.
func TestTheQRowsAreInTheOrderOfTheNumbersOfTheQuestions(t *testing.T) {
	d := setupWorld(sa("Q-1000", "b", "902"), sa("Q-999", "a", "901"), sa("Q-100", "c", "903"))
	d.Questions = append(d.Questions, sq("Q-999", "rule", "Q-999", "idea-owner", "Nine?", "answered"), sq("Q-1000", "rule", "Q-1000", "operator", "Thousand?", "answered"), sq("Q-100", "rule", "Q-100", "idea-owner", "Hundred?", "answered"))
	rows, asked, _ := setupAnswers(d)
	var ids []string
	for _, r := range rows {
		ids = append(ids, r[0])
	}
	if want := []string{"Q-100", "Q-999", "Q-1000"}; !reflect.DeepEqual(ids, want) {
		t.Errorf("the order of the rows %q; want %q", ids, want)
	}
	if len(asked) != 3 || asked[2].ID != "Q-1000" || asked[2].Text != "Thousand?" || rows[2][2] != "operator" {
		t.Errorf("the questions %q and the last row %q; want Q-1000 last, with its text and the role of its block", asked, rows[2])
	}
}
