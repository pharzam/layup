package run

// The composers of the question and the answer comment (task T-gyvo, #405).
// A fixture is read by the records' own readers, so it is the row a reader
// gives; a refusal that no reader would pass is an explicit slice that differs
// from a valid control in one token.

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/records"
)

const qaSHA = "0123456789abcdef0123456789abcdef01234567"

const qaStartTSV = "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy\n" +
	"S-0a1b2c3d\tT-ab12\t2\tdeveloper\tclaude\t2.1.295\tclaude-opus-5-5\t" + qaSHA + "\t" + qaSHA + "\t1001\t251\t200000\t10.0\t60\t—\t—\n"

const qaResultTSV = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n" +
	"question\t1\tFirst?\t—\tlabel needs_human\nquestion\t2\t2\t—\tlabel needs_human S-ffffffff.9\nquestion\t3\tWhich owner?\t—\tlabel needs_human\n" +
	"answer\t1\tUse the plan.\t—\tdocs/spec/loop.md\nanswer\t2\tThe architect owns it.\t—\tdocs/spec/loop.md S-ffffffff.9\n"

func qaStart(t *testing.T) []string {
	t.Helper()
	rows, err := records.ReadSessions([]byte(qaStartTSV))
	if err != nil || len(rows) != 1 {
		t.Fatalf("the fixture start row: %v", err)
	}
	return rows[0]
}

// qaRow gives the row of kind and n of the fixture result.
func qaRow(t *testing.T, kind, n string) []string {
	t.Helper()
	rows, err := records.ReadResult([]byte(qaResultTSV))
	if err != nil {
		t.Fatalf("the fixture result: %v", err)
	}
	for _, r := range rows {
		if r[0] == kind && r[1] == n {
			return r
		}
	}
	t.Fatalf("no row %s %s", kind, n)
	return nil
}

// qaWith gives a copy of row whose column i is v.
func qaWith(row []string, i int, v string) []string {
	c := slices.Clone(row)
	c[i] = v
	return c
}

const wantQuestion = "S-0a1b2c3d: developer session of T-ab12, attempt 2: question S-0a1b2c3d.3\n\nWhich owner?"

func TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText(t *testing.T) {
	got, err := QuestionComment(qaStart(t), qaRow(t, "question", "3"))
	if err != nil || string(got) != wantQuestion {
		t.Errorf("got %q, %v; want %q", got, err, wantQuestion)
	}
}

func TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers(t *testing.T) {
	start := qaWith(qaStart(t), 0, "S-1111aaaa")
	want := "S-1111aaaa: developer session of T-ab12, attempt 2: answer S-1111aaaa.2 of S-0a1b2c3d.3\n\nThe architect owns it."
	// The row cites another ID in its reason; the argument is the question.
	got, err := AnswerComment(start, qaRow(t, "answer", "2"), "S-0a1b2c3d.3")
	if err != nil || string(got) != want {
		t.Errorf("got %q, %v; want %q", got, err, want)
	}
	// A legal answer row that cites no ID composes.
	want = "S-1111aaaa: developer session of T-ab12, attempt 2: answer S-1111aaaa.1 of S-0a1b2c3d.3\n\nUse the plan."
	got, err = AnswerComment(start, qaRow(t, "answer", "1"), "S-0a1b2c3d.3")
	if err != nil || string(got) != want {
		t.Errorf("got %q, %v; want %q", got, err, want)
	}
}

func TestTheIDUsesTheColumnNNeverTheTextNorTheReason(t *testing.T) {
	start := qaStart(t)
	for _, c := range []struct{ n, text string }{{"7", "1"}, {"1", "2"}, {"2", "Which owner?"}} {
		q := []string{"question", c.n, c.text, "", "label needs_human S-ffffffff.9"}
		got, err := QuestionComment(start, q)
		if want := "question S-0a1b2c3d." + c.n + "\n\n" + c.text; err != nil || !strings.HasSuffix(string(got), want) {
			t.Errorf("question n %s text %q: got %q, %v; want it to end with %q", c.n, c.text, got, err, want)
		}
		a := []string{"answer", c.n, c.text, "", "docs/spec/loop.md S-ffffffff.9"}
		got, err = AnswerComment(start, a, "S-0a1b2c3d.3")
		if want := "answer S-0a1b2c3d." + c.n + " of S-0a1b2c3d.3\n\n" + c.text; err != nil || !strings.HasSuffix(string(got), want) {
			t.Errorf("answer n %s text %q: got %q, %v; want it to end with %q", c.n, c.text, got, err, want)
		}
	}
}

func TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers(t *testing.T) {
	start := qaStart(t)
	const want = "S-0a1b2c3d: developer session of T-ab12, attempt 2: "
	if got := firstLinePrefix(start); got != want {
		t.Fatalf("firstLinePrefix %q; want %q", got, want)
	}
	// Each field changed alone changes only its place in the line.
	line := func(id, task, attempt, role string) string {
		return id + ": " + role + " session of " + task + ", attempt " + attempt + ": "
	}
	for col, w := range map[int]string{
		sessionID:      line("S-0a1b2c3e", "T-ab12", "2", "developer"),
		sessionTask:    line("S-0a1b2c3d", "T-ab13", "2", "developer"),
		sessionAttempt: line("S-0a1b2c3d", "T-ab12", "3", "developer"),
		sessionRole:    line("S-0a1b2c3d", "T-ab12", "2", "reviewer"),
	} {
		v := map[int]string{sessionID: "S-0a1b2c3e", sessionTask: "T-ab13", sessionAttempt: "3", sessionRole: "reviewer"}[col]
		if got := firstLinePrefix(qaWith(start, col, v)); got != w {
			t.Errorf("column %d changed to %q gives the prefix %q; want %q", col, v, got, w)
		}
	}
	q, _ := QuestionComment(start, qaRow(t, "question", "3"))
	a, _ := AnswerComment(start, qaRow(t, "answer", "1"), "S-0a1b2c3d.3")
	if !strings.HasPrefix(string(q), want+"question ") || !strings.HasPrefix(string(a), want+"answer ") {
		t.Errorf("the first lines %q and %q do not start with the prefix", q, a)
	}
}

func TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved(t *testing.T) {
	start, row := qaStart(t), qaRow(t, "answer", "1")
	for _, ok := range []string{"S-0a1b2c3d.3", "S-00000000.1", "S-ffffffff.10"} {
		if got, err := AnswerComment(start, row, ok); err != nil || !strings.HasSuffix(strings.SplitN(string(got), "\n", 2)[0], " of "+ok) {
			t.Errorf("question ID %q: %q, %v; want it composed", ok, got, err)
		}
	}
	for _, bad := range []string{"", "S-0a1b2c3d", "0a1b2c3d.3", ".3", "S-0a1b2c3d.0", "S-0a1b2c3d.03", "S-0A1B2C3D.3", "S-0a1b2c3.3", "s-0a1b2c3d.3",
		" S-0a1b2c3d.3", "S-0a1b2c3d.3\n", "S-0a1b2c3d.3 ", "S-0a1b2c3d.3 x", "S-0a1b2c3d.3 S-0a1b2c3d.4"} {
		if got, err := AnswerComment(start, row, bad); err == nil || got != nil {
			t.Errorf("question ID %q: %q, %v; want an error and no bytes", bad, got, err)
		}
	}
}

// A refusal is one token away from its valid control, which composes. The
// boundary is the reader's: the schema and the types of the columns, then the
// rules of the row, with no repair of a byte.
func TestAMalformedRowGivesAnErrorAndNoBytes(t *testing.T) {
	start := qaStart(t)
	rows := map[string][]string{"q": qaRow(t, "question", "3"), "a": qaRow(t, "answer", "1")}
	compose := func(who string, s, r []string) ([]byte, error) {
		if who == "q" {
			return QuestionComment(s, r)
		}
		return AnswerComment(s, r, "S-0a1b2c3d.3")
	}
	type refusal struct {
		name, who string
		onStart   bool
		col       int
		bad, good string
	}
	var cases []refusal
	add := func(name, who string, onStart bool, col int, good string, bads ...string) {
		for _, b := range bads {
			cases = append(cases, refusal{name + " " + b, who, onStart, col, b, good})
		}
	}
	add("kind of question", "q", false, 0, "question", "answer", "decision", "lesson", "status", "comment")
	add("kind of answer", "a", false, 0, "answer", "question", "decision", "lesson")
	for _, who := range []string{"q", "a"} {
		add("n of "+who, who, false, 1, "3", "0", "01", "x", "", "-1", "1.5")
		add("text of "+who, who, false, 2, "Which owner?", "", "a\nb", "a\rb", "a\tb", "a\xffb", "—")
		add("session of "+who, who, true, 0, "S-0a1b2c3d", "S-0A1B2C3D", "S-0a1b2c3", "")
		add("task of "+who, who, true, 1, "T-ab12", "bad", "T-AB12", "")
		add("attempt of "+who, who, true, 2, "2", "02", "0", "x", "")
		add("role of "+who, who, true, 3, "developer", "", "de\nveloper", "de\tveloper", "de\rveloper")
		add("harness of "+who, who, true, 4, "claude", "two words", "")
	}
	add("sha of question", "q", false, 3, "", strings.Repeat("0", 64))
	add("reason of question", "q", false, 4, "label needs_human", "label", "label  needs_human", " label needs_human", "—")
	add("reason of answer", "a", false, 4, "docs/spec/loop.md", "", "a  b", "a ")
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, r, gs, gr := start, rows[c.who], start, rows[c.who]
			if c.onStart {
				s, gs = qaWith(start, c.col, c.bad), qaWith(start, c.col, c.good)
			} else {
				r, gr = qaWith(r, c.col, c.bad), qaWith(r, c.col, c.good)
			}
			if good, err := compose(c.who, gs, gr); err != nil || len(good) == 0 {
				t.Fatalf("the control %q does not compose: %v", c.good, err)
			}
			if got, err := compose(c.who, s, r); err == nil || got != nil {
				t.Errorf("got %q, %v; want an error and no bytes", got, err)
			}
		})
	}
	for who, r := range rows {
		for name, bad := range map[string][2][]string{
			"row of 4": {start, r[:4]}, "row of 6": {start, append(slices.Clone(r), "x")},
			"start of 15": {start[:15], r}, "start of 17": {append(slices.Clone(start), "x"), r}, "no start": {nil, r}, "no row": {start, nil},
		} {
			if got, err := compose(who, bad[0], bad[1]); err == nil || got != nil {
				t.Errorf("%s of %s: %q, %v; want an error and no bytes", name, who, got, err)
			}
		}
	}
}

func TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes(t *testing.T) {
	start := qaStart(t)
	for _, who := range []string{"q", "a"} {
		compose := func(text string) ([]byte, error) {
			if who == "q" {
				return QuestionComment(start, []string{"question", "3", text, "", "label needs_human"})
			}
			return AnswerComment(start, []string{"answer", "3", text, "", "docs/spec/loop.md"}, "S-0a1b2c3d.4")
		}
		probe, err := compose("x")
		if err != nil {
			t.Fatal(err)
		}
		fill := MaxCommentChars - (utf8.RuneCount(probe) - 1) // characters of text for a body of exactly the limit
		for name, text := range map[string]string{"limit": strings.Repeat("a", fill), "limit in three-byte letters": strings.Repeat("€", fill)} {
			got, err := compose(text)
			if err != nil || utf8.RuneCount(got) != MaxCommentChars {
				t.Errorf("%s of %s: %d characters, %v; want a body of %d", name, who, utf8.RuneCount(got), err, MaxCommentChars)
			}
		}
		got, err := compose(strings.Repeat("a", fill+1))
		var ref HandoffRefusal
		if !errors.As(err, &ref) || ref.Reason != "handoff" || got != nil {
			t.Errorf("one character over of %s: %q, %v; want a HandoffRefusal and no bytes", who, got, err)
		}
	}
}

func TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged(t *testing.T) {
	start, row := qaStart(t), qaRow(t, "question", "3")
	s0, r0 := slices.Clone(start), slices.Clone(row)
	first, err := QuestionComment(start, row)
	if err != nil || string(first) != wantQuestion {
		t.Fatalf("got %q, %v; want %q", first, err, wantQuestion)
	}
	first[0] = 'X'
	second, err := QuestionComment(start, row)
	if err != nil || string(second) != wantQuestion {
		t.Errorf("second call %q, %v; want %q", second, err, wantQuestion)
	}
	if !slices.Equal(start, s0) || !slices.Equal(row, r0) {
		t.Errorf("the inputs changed: %q %q", start, row)
	}
}
