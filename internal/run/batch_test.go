package run

import (
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
	workpkg "github.com/pharzam/layup/internal/work"
)

// batchPSB is a problem statement of three lines: "# A" (bytes 0 to 3), a
// line with two undefined abbreviations (4 to 30), and "It is fast." (31 to 42).
const batchPSB = "# A\nWe use Go and TLA and XYZ.\nIt is fast.\n"

func batchFixture(t *testing.T) batchRecords {
	t.Helper()
	read := func(data string, s tsv.Schema) [][]string {
		rows, err := tsv.Read([]byte(data), s)
		if err != nil {
			t.Fatal(err)
		}
		return rows
	}
	return batchRecords{
		PSB: []byte(batchPSB),
		Gaps: read("id\trule\tline\texcerpt\tquestion\n"+
			"Q-001\tG1\t0\t—\tWhich stack?\n"+
			"Q-002\tG3\t2\tWe use Go and TLA and XYZ.\tWhat is TLA?\n"+
			"Q-003\tG3\t2\tWe use Go and TLA and XYZ.\tWhat is XYZ?\n", psbGapsSchema),
		Meaning: read("id\tkind\tstart\tend\tquote\tquestion\tsession\n"+
			"M-001\tunmeasurable\t31\t42\tIt is fast.\tHow fast?\tS-11111111\n"+
			"M-002\tunmeasurable\t31\t42\tIt is fast.\tFast as what?\tS-11111111\n"+
			"M-003\ttwo-readings\t31\t42\tIt is fast.\tFast to run or to build?\tS-11111111\n", records.MeaningSchema),
		Coverage: read("id\tkind\tfact\trequirement\tquestion\tsession\n"+
			"C-001\tuncovered\t2\t—\tWhat covers fact 2?\tS-22222222\n"+
			"C-002\tunsupported\t—\tR-001\tWhich fact supports R-001?\tS-22222222\n", records.CoverageSchema),
		Spans: read("start\tend\tfact\tclass\treason\ttext\n"+
			"0\t3\t—\t—\ta heading\t# A\n4\t30\t1\tneed\t—\tWe use Go and TLA and XYZ.\n31\t42\t2\tmeasure\t—\tIt is fast.\n", records.SpansSchema),
		// The markers of the pinned tree: one with a proposal, one with none, and one whose
		// proposal quotes the same line as the first's (a repeat of quote and kind, kept).
		Markers: read("id\tfile\tline\tmarker\tvalue\tstart\tend\tsession\n"+
			"M-0a1b2c3d\tdocs/a.md\t2\t"+angled("owner")+"\tCarol\t4\t30\tS-33333333\n"+
			"M-4e5f6a7b\tdocs/b.md\t1\t"+angled("port")+"\t—\t—\t—\tS-33333333\n"+
			"M-8c9d0e1f\tdocs/b.md\t5\t"+angled("host")+"\tGo\t4\t30\tS-33333333\n", records.MarkersSchema),
	}
}

// The merge and the IDs (docs/spec/intake.md): the rule gaps, the gaps of
// meaning, the items of the completeness review and the intent form, in that
// order, then the setup questions of the Operator's block, S01-stack and each
// marker; a question whose quote and kind repeat one before it is dropped (a
// repeat is within one source: the kinds of the sources differ), a question
// with no quote never, and a setup question never (task T-4c5c: two markers
// may share a quote); each kept question gets Q-NNN in the order of the merge.
func TestMergeQuestions(t *testing.T) {
	rows, err := mergeQuestions(batchFixture(t), "2026-10-11T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"Q-001 rule Q-001 G1 — — layup psb check",
		"Q-002 rule Q-002 G3 4 30 layup psb check", // the line, trimmed of its line break
		// Q-003 of the gap table: the same line and rule, dropped (a known limit)
		"Q-003 meaning M-001 unmeasurable 31 42 S-11111111",
		// M-002: the same quote and kind as M-001, dropped
		"Q-004 meaning M-003 two-readings 31 42 S-11111111",
		"Q-005 coverage C-001 uncovered 31 42 S-22222222", // its fact's span
		"Q-006 coverage C-002 unsupported — — S-22222222",
		"Q-007 intent band.B intent — — layup run",
		"Q-008 intent band.U intent — — layup run",
		"Q-009 intent appetite.hours intent — — layup run",
		"Q-010 setup S01-stack setup — — layup run",
		"Q-011 setup M-0a1b2c3d setup 4 30 S-33333333",
		"Q-012 setup M-4e5f6a7b setup — — S-33333333",
		"Q-013 setup M-8c9d0e1f setup 4 30 S-33333333", // the quote and the kind of Q-011: kept
	}
	var got []string
	for _, r := range rows {
		dash := func(s string) string {
			if s == "" {
				return "—"
			}
			return s
		}
		got = append(got, strings.Join([]string{r[0], r[1], r[2], r[5], dash(r[6]), dash(r[7]), r[10]}, " "))
		block := "idea-owner"
		if r[1] == "setup" {
			block = "operator"
		}
		if r[3] != block || r[4] != "before-delivery" || r[9] != "2026-10-11T10:00:00Z" || r[11] != "" || r[12] != "asked" {
			t.Errorf("the row %q: want %s, before-delivery, the time asked, no answer, asked", r, block)
		}
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("the merge:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	// The question of S01-stack is the text of S01 of phase 1; a marker's is "the value of
	// <marker> at <file>:<line>", as the specification writes it.
	if rows[9][8] != workpkg.S01Questions[0].Text {
		t.Errorf("the question of S01-stack %q; want %q", rows[9][8], workpkg.S01Questions[0].Text)
	}
	for i, want := range []string{"the value of " + angled("owner") + " at docs/a.md:2", "the value of " + angled("port") + " at docs/b.md:1", "the value of " + angled("host") + " at docs/b.md:5"} {
		if rows[10+i][8] != want {
			t.Errorf("the question of marker %d %q; want %q", i+1, rows[10+i][8], want)
		}
	}
	table, err := table(records.QuestionsSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := records.ReadQuestions(table); err != nil {
		t.Errorf("ReadQuestions of the merge: %v", err)
	}
}

// S01-stack is asked whatever the rule check found (F-21, decision 4 of #206: the stack is the
// Operator's), once, with no marker and with no rule gap at all; the batch of no marker has the one
// setup question.
func TestMergeQuestionsAsksTheStackWithNoMarker(t *testing.T) {
	b := batchFixture(t)
	b.Markers = nil
	b.Gaps = b.Gaps[1:] // no G1
	rows, err := mergeQuestions(b, "2026-10-11T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	last := rows[len(rows)-1]
	if last[1] != "setup" || last[2] != "S01-stack" || last[3] != "operator" || last[0] != "Q-009" {
		t.Errorf("the last row %q; want Q-009, setup, S01-stack, operator", last)
	}
	setups := 0
	for _, r := range rows {
		if r[1] == "setup" {
			setups++
		}
	}
	if setups != 1 {
		t.Errorf("%d setup questions; want S01-stack alone", setups)
	}
}

// The comment of the batch: the header line; the block for the idea owner,
// each question and its quote; the intent form; the block for the Operator,
// S01-stack and each marker, with its proposal and its quote; the rule of the
// answers, and the three forms of a marker's answer.
func TestTheCommentOfTheBatch(t *testing.T) {
	fix := batchFixture(t)
	rows, err := mergeQuestions(fix, "2026-10-11T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	c := batchComment(rows, []byte(batchPSB), fix.Markers)
	for _, want := range []string{
		batchHeader + "\n",
		"\n## For the idea owner\n",
		"\nQ-001: Which stack?\n\n",
		"\nQ-002: What is TLA?\n`We use Go and TLA and XYZ.`\n",
		"\nQ-003: How fast?\n`It is fast.`\n",
		"\nQ-006: Which fact supports R-001?\n\n",
		"\nQ-007: " + intentQuestions["band.B"] + "\n",
		"\nQ-009: " + intentQuestions["appetite.hours"] + "\n",
		"\nQ-009: " + intentQuestions["appetite.hours"] + "\n\n## For the Operator\n\nQ-010: " + workpkg.S01Questions[0].Text + "\n\n",
		"\nQ-011: the value of " + angled("owner") + " at docs/a.md:2\nproposed: Carol\n`We use Go and TLA and XYZ.`\n\n",
		"\nQ-012: the value of " + angled("port") + " at docs/b.md:1\n\n",
		"\nQ-013: the value of " + angled("host") + " at docs/b.md:5\nproposed: Go\n`We use Go and TLA and XYZ.`\n",
		"`Q-NNN: <answer>`",
		"`Q-NNN: accept`", "`Q-NNN: <value>`", "`Q-NNN: gap`",
	} {
		if !strings.Contains(c, want) {
			t.Errorf("the comment has no %q:\n%s", want, c)
		}
	}
	if !strings.HasPrefix(c, batchHeader+"\n") || strings.Index(c, "## For the idea owner") > strings.Index(c, "## For the Operator") {
		t.Errorf("the comment does not start with the header, or its blocks are out of order:\n%s", c)
	}
}

// Each question is in the block of its row, the idea owner's first; a block with
// no question says so; the three forms of a marker's answer are in the rule of
// the answers only when a marker is asked (task T-4c5c, #276).
func TestTheBlocksOfAComment(t *testing.T) {
	fix := batchFixture(t)
	rows, err := mergeQuestions(fix, "2026-10-11T10:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	var idea, operator [][]string
	for _, r := range rows {
		if r[3] == "operator" {
			operator = append(operator, r)
		} else {
			idea = append(idea, r)
		}
	}
	quote := func(r []string) (string, string) { return r[6], r[7] }
	about := func(r []string) (bool, string) { return aboutAMarker(r, fix.Markers) }
	c := questionsComment("H", idea, quote, about, []byte(batchPSB))
	if !strings.Contains(c, "\n## For the Operator\n\nNo question for the Operator.\n\n") || strings.Contains(c, "`Q-NNN: accept`") {
		t.Errorf("the comment of the idea owner's questions alone:\n%s", c)
	}
	c = questionsComment("H", operator[:1], quote, about, []byte(batchPSB))
	if !strings.Contains(c, "\n## For the idea owner\n\nNo question for the idea owner.\n\n## For the Operator\n\nQ-010: ") || strings.Contains(c, "`Q-NNN: accept`") {
		t.Errorf("the comment of S01-stack alone:\n%s", c)
	}
	c = questionsComment("H", operator, quote, about, []byte(batchPSB))
	if !strings.Contains(c, "`Q-NNN: accept`") || !strings.Contains(c, "`Q-NNN: gap`") {
		t.Errorf("the comment of the markers has no forms of a marker's answer:\n%s", c)
	}
	if strings.Index(c, "## For the idea owner") > strings.Index(c, "## For the Operator") || strings.Index(c, "Answer each block") < strings.Index(c, "## For the Operator") {
		t.Errorf("the blocks, or the rule of the answers, are out of order:\n%s", c)
	}
}

// A quote in a run of backquotes one longer than its longest run, with a space
// inside when it begins or ends with a backquote.
func TestBackquote(t *testing.T) {
	for _, c := range []struct{ in, want string }{
		{"plain", "`plain`"},
		{"a`b", "``a`b``"},
		{"`a", "`` `a ``"},
		{"a`", "`` a` ``"},
		{"a``b", "```a``b```"},
	} {
		if got := backquote(c.in); got != c.want {
			t.Errorf("backquote(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// The quote of a rule gap is its line, trimmed of its line break and of a CR
// before it; a line that the file lacks is none.
func TestLineSpan(t *testing.T) {
	psb := []byte("one\r\ntwo\nthree")
	for _, c := range []struct {
		n          int
		start, end int
		ok         bool
	}{{1, 0, 3, true}, {2, 5, 8, true}, {3, 9, 14, true}, {4, 0, 0, false}} {
		start, end, ok := lineSpan(psb, c.n)
		if start != c.start || end != c.end || ok != c.ok {
			t.Errorf("line %d: %d, %d, %v; want %d, %d, %v", c.n, start, end, ok, c.start, c.end, c.ok)
		}
	}
}
