package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// The fixtures of the revision's check, each through its production reader: a
// span table of a heading and four facts (need, constraint, measure, context),
// three questions (Q-001's quote holds fact 3's span, Q-002 is the intent
// field success-criteria with no quote, Q-003's quote holds fact 1's span),
// the answers to Q-001 and Q-002, and the draft before the revision.
func revisionFixture(t *testing.T) (spans, questions, answers, before [][]string) {
	t.Helper()
	read := func(data string, s tsv.Schema, r func([]byte) ([][]string, error)) [][]string {
		rows, err := r([]byte(data))
		if err != nil {
			t.Fatal(err)
		}
		return rows
	}
	spans = read("start\tend\tfact\tclass\treason\ttext\n"+
		"0\t3\t—\t—\ta heading\t# A\n"+
		"4\t10\t1\tneed\t—\tneed 1\n"+
		"11\t17\t2\tconstraint\t—\tcons 2\n"+
		"18\t24\t3\tmeasure\t—\tmeas 3\n"+
		"25\t31\t4\tcontext\t—\tcont 4\n", records.SpansSchema, records.ReadSpans)
	questions = read("id\tsource\tref\tblock\tphase\tkind\tstart\tend\tquestion\tasked\tasker\tanswered\tstate\n"+
		"Q-001\tcoverage\tC-001\tidea-owner\tbefore-delivery\tuncovered\t18\t24\tWhat covers fact 3?\t2026-10-11T10:00:00Z\tS-22222222\t2026-10-11T10:05:00Z\tanswered\n"+
		"Q-002\tintent\tsuccess-criteria\tidea-owner\tbefore-delivery\tintent\t—\t—\tThe criteria?\t2026-10-11T10:00:00Z\tlayup run\t2026-10-11T10:05:00Z\tanswered\n"+
		"Q-003\tmeaning\tM-001\tidea-owner\tbefore-delivery\tunmeasurable\t5\t10\tMeasurable?\t2026-10-11T10:00:00Z\tS-11111111\t—\tasked\n", records.QuestionsSchema, records.ReadQuestions)
	answers = read("question\tanswer\tcomment\tauthor\ttime\n"+
		"Q-001\tR-004 covers it.\t801\t303\t2026-10-11T10:05:00Z\n"+
		"Q-002\tfast; cheap\t801\t303\t2026-10-11T10:05:00Z\n", records.AnswersSchema, records.ReadAnswers)
	before = read("id\tkind\tstatement\tcovers\tcriterion\tsession\n"+
		"R-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS-0a1b2c3d\n"+
		"R-002\tnon-functional\tThe tool keeps 2.\t2\tIt keeps 2.\tS-0a1b2c3d\n"+
		"R-003\tsuccess-criterion\tThe tool meets 3.\t3\tIt meets 3.\tS-0a1b2c3d\n", records.DraftSchema, records.ReadDraft)
	return spans, questions, answers, before
}

// The check of the revision (docs/spec/loop.md, The revision of the draft):
// the draft check's rules, widened so a covers token may be the ID of an
// answer's fact — a Q-NNN with a row of intake/answers.tsv; a row covering an
// answer's fact names the fact of the answer's quote; a fact whose covering
// kinds changed is covered by a row naming an answer's fact; and no row of the
// session covers the success-criteria answer, whose rows the code appends.
func TestCheckRevision(t *testing.T) {
	spans, questions, answers, before := revisionFixture(t)
	file := func(rows ...string) [][]string {
		var out [][]string
		for _, r := range rows {
			out = append(out, strings.Split(r, "\t"))
		}
		return out
	}
	good := file(
		"requirement\tThe tool does 1.\t1\tIt does 1.",
		"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
		"success-criterion\tThe tool meets 3.\t3\tIt meets 3.",
	)
	covered4, err := records.ReadDraft([]byte("id\tkind\tstatement\tcovers\tcriterion\tsession\n" +
		"R-001\trequirement\tThe tool does 1.\t1\tIt does 1.\tS-0a1b2c3d\n" +
		"R-002\tnon-functional\tThe tool keeps 2.\t2\tIt keeps 2.\tS-0a1b2c3d\n" +
		"R-003\tsuccess-criterion\tThe tool meets 3.\t3\tIt meets 3.\tS-0a1b2c3d\n" +
		"R-004\trequirement\tThe context of 4.\t4\tIt rests on 4.\tS-0a1b2c3d\n"))
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name   string
		rows   [][]string
		before [][]string // nil for the fixture's draft
		want   string     // "" for a pass, else a word of the finding
	}{
		{"the draft kept, no citation needed", good, nil, ""},
		{"an answer's fact with the fact of its quote", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3 Q-001\tIt meets 3."), nil, ""},
		{"a row covering an answer's fact, the quote's fact missing", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3.",
			"requirement\tThe answer of Q-001.\tQ-001\t—"), nil, "Q-001"},
		{"a Q-NNN that is no fact of the answers", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3 Q-003\tIt meets 3."), nil, "Q-003"},
		{"a Q-NNN of no answer, its quote's fact covered", file(
			"requirement\tThe tool does 1.\t1 Q-003\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3."), nil, "Q-003"},
		{"the success-criteria answer covered by the session", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3.",
			"success-criterion\tThe answer of Q-002.\tQ-002\tfast."), nil, "Q-002"},
		{"a need fact covered by no row", file(
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3."), nil, "need"},
		{"a measure fact with no success criterion", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"requirement\tThe tool does 3.\t3\tIt does 3."), nil, "measure"},
		{"a row covering no fact", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3.",
			"requirement\tNothing.\t—\t—"), nil, "covers"},
		{"a changed class set with no answer's fact", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"requirement\tThe tool does 3.\t3\tIt does 3.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3."), nil, "3"},
		{"a changed class set named by an answer's fact", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"requirement\tThe tool does 3.\t3 Q-001\tIt does 3.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3."), nil, ""},
		{"a context fact newly covered, no citation needed", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3.",
			"requirement\tThe context of 4.\t4\tIt rests on 4."), nil, ""},
		{"a context fact dropped, no citation needed", file(
			"requirement\tThe tool does 1.\t1\tIt does 1.",
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.",
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3."), covered4, ""},
	} {
		old := before
		if c.before != nil {
			old = c.before
		}
		err := checkRevision(c.rows, spans, questions, answers, old, "Q-002")
		if c.want == "" && err != nil {
			t.Errorf("%s: %v; want a pass", c.name, err)
		}
		if c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)) {
			t.Errorf("%s: %v; want a refusal naming %q", c.name, err, c.want)
		}
	}
}

// The fact of an answer's quote (docs/spec/loop.md, Decided here): the fact of
// the span that holds the quote's first byte; a quote in no fact's span has
// none.
func TestTheFactOfAQuote(t *testing.T) {
	spans, _, _, _ := revisionFixture(t)
	for _, c := range []struct{ start, end, want string }{
		{"18", "24", "3"}, {"5", "10", "1"}, {"25", "31", "4"}, {"0", "3", ""}, {"32", "35", ""},
	} {
		if got := factOfQuote(c.start, c.end, spans); got != c.want {
			t.Errorf("the fact of the quote %s..%s: %q; want %q", c.start, c.end, got, c.want)
		}
	}
}

// The Keep of the revision: the session's file, checked, R-NNN in order with
// the session, and the success criteria of the intent answer appended as
// success-criterion rows covering its Q-NNN, each criterion the idea owner's
// words (docs/spec/loop.md, The approvers and the success criteria). A refused
// file is the error.
func TestKeepOfRevision(t *testing.T) {
	spans, questions, answers, before := revisionFixture(t)
	_, criteria := criteriaAnswer(questions, answers)
	keep := revisionKeep(spans, questions, answers, before, "Q-002", criteria)
	files, err := keep("S-aaaabbbb", []byte(
		"kind\tstatement\tcovers\tcriterion\n"+
			"requirement\tThe tool does 1.\t1\tIt does 1.\n"+
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.\n"+
			"success-criterion\tThe tool meets 3.\t3 Q-001\tIt meets 3.\n"))
	if err != nil {
		t.Fatal(err)
	}
	rows, err := records.ReadDraft(files["spec/draft.tsv"])
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"R-001 requirement The tool does 1. 1 It does 1. S-aaaabbbb",
		"R-002 non-functional The tool keeps 2. 2 It keeps 2. S-aaaabbbb",
		"R-003 success-criterion The tool meets 3. 3 Q-001 It meets 3. S-aaaabbbb",
		"R-004 success-criterion The delivery's success criteria. Q-002 fast S-aaaabbbb",
		"R-005 success-criterion The delivery's success criteria. Q-002 cheap S-aaaabbbb",
	}
	if len(rows) != len(want) {
		t.Fatalf("the kept draft %v; want %v", rows, want)
	}
	for i := range want {
		if strings.Join(rows[i], " ") != want[i] {
			t.Errorf("row %d: %q; want %q", i+1, rows[i], want[i])
		}
	}
	if _, err := keep("S-aaaabbbb", []byte(
		"kind\tstatement\tcovers\tcriterion\n"+
			"requirement\tThe tool does 1.\t1\tIt does 1.\n"+
			"non-functional\tThe tool keeps 2.\t2\tIt keeps 2.\n"+
			"success-criterion\tThe tool meets 3.\t3\tIt meets 3.\n"+
			"requirement\tThe answer of Q-001.\tQ-001\t—\n")); err == nil || !strings.Contains(err.Error(), "Q-001") {
		t.Errorf("a file that drops the quote's fact: %v; want a refusal naming Q-001", err)
	}
}
