package records

import (
	"strings"
	"testing"
)

// The records of Intake (task T-529i, #212): one valid file per record, one
// case per rule that its block gives in words, and one case per column whose
// rule has no clause for the empty value. Each refusal must name the column of
// the rule (refusedAt), so a case that fails for another reason fails the test.

// One question of each source, each the first row of its table.
var (
	rule1     = "Q-001\trule\tQ-001\tidea-owner\tbefore-delivery\tG1\t—\t—\tWhich technology stack?\t2026-10-10T12:00:00Z\tlayup psb check\t—\tasked"
	coverage1 = "Q-001\tcoverage\tC-001\tidea-owner\tbefore-delivery\tuncovered\t12\t40\tWhich requirement covers fact 1?\t2026-10-10T12:00:00Z\tS-1a2b3c4d\t—\tasked"
	intent1   = "Q-001\tintent\tband.U\tidea-owner\tbefore-delivery\tintent\t—\t—\tThe band U, in USD?\t2026-10-10T12:00:00Z\tlayup run\t—\tasked"
	followup1 = "Q-001\tfollow-up\tQ-001\tidea-owner\tbefore-delivery\tfollow-up\t—\t—\tWhich app, again?\t2026-10-10T14:00:00Z\tlayup run\t—\tasked"
	// The setup questions of M2d (docs/spec/scaffold.md, The answers of the setup; task T-4c5c, #276): the Operator's.
	stack1  = "Q-001\tsetup\tS01-stack\toperator\tbefore-delivery\tsetup\t—\t—\tWhich stack?\t2026-10-10T12:00:00Z\tlayup run\t—\tasked"
	marker1 = "Q-001\tsetup\tM-1a2b3c4d\toperator\tbefore-delivery\tsetup\t12\t40\tthe value of a marker\t2026-10-10T12:00:00Z\tS-1a2b3c4d\t—\tasked"
)

const (
	meaningHeader  = "id\tkind\tstart\tend\tquote\tquestion\tsession"
	meaningRow     = "M-001\ttwo-readings\t10\t17\tthe app\tWhich app is meant?\tS-1a2b3c4d"
	spansHeader    = "start\tend\tfact\tclass\treason\ttext"
	factRow        = "12\t40\t1\tneed\t—\tThe app sends one reminder a day."
	headingRow     = "0\t11\t—\t—\ta heading\t# The brief"
	draftHeader    = "id\tkind\tstatement\tcovers\tcriterion\tsession"
	draftRow       = "R-001\trequirement\tThe app sends one reminder a day.\t1 2\tOne reminder is sent per day.\tS-1a2b3c4d"
	coverageHeader = "id\tkind\tfact\trequirement\tquestion\tsession"
	uncoveredRow   = "C-001\tuncovered\t3\t—\tWhich requirement covers fact 3?\tS-1a2b3c4d"
	unsupportedRow = "C-002\tunsupported\t—\tR-002\tWhich fact supports R-002?\tS-1a2b3c4d"
	questionsHead  = "id\tsource\tref\tblock\tphase\tkind\tstart\tend\tquestion\tasked\tasker\tanswered\tstate"
	askedRow       = "Q-001\tmeaning\tM-001\tidea-owner\tbefore-delivery\ttwo-readings\t10\t17\tWhich app is meant?\t2026-10-10T12:00:00Z\tS-1a2b3c4d\t—\tasked"
	intentRow      = "Q-002\tintent\tband.B\tidea-owner\tbefore-delivery\tintent\t—\t—\tThe band B, in USD?\t2026-10-10T12:00:00Z\tlayup run\t2026-10-10T13:00:00Z\tanswered"
	answersHeader  = "question\tanswer\tcomment\tauthor\ttime"
	answerRow      = "Q-001\tThe mobile app.\t4001\t1675602\t2026-10-10T13:00:00Z"
	reviewHeader   = "question\tstate\treason\tsession"
	reviewRow      = "Q-001\tclosed\tThe answer names one app.\tS-1a2b3c4d"
	budgetHeader   = "name\tvalue\tsource"
)

func readMeaning(d []byte) ([][]string, error)       { return ReadMeaning(d) }
func readSpans(d []byte) ([][]string, error)         { return ReadSpans(d) }
func readDraft(d []byte) ([][]string, error)         { return ReadDraft(d) }
func readCoverage(d []byte) ([][]string, error)      { return ReadCoverage(d) }
func readQuestions(d []byte) ([][]string, error)     { return ReadQuestions(d) }
func readAnswers(d []byte) ([][]string, error)       { return ReadAnswers(d) }
func readAnswersReview(d []byte) ([][]string, error) { return ReadAnswersReview(d) }
func readBudget(d []byte) ([][]string, error)        { return ReadBudget(d) }

func TestValidIntakeRecordsAreRead(t *testing.T) {
	m2 := "M-002\tdecision\t50\t62\tone reminder\tOne reminder for whom?\tS-0000000a"
	open := strings.Replace(strings.Replace(intentRow, "Q-002", "Q-003", 1), "\tanswered", "\topen", 1)
	unanswered := strings.Replace(strings.Replace(open, "2026-10-10T13:00:00Z", "—", 1), "Q-003", "Q-004", 1)
	closed := strings.Replace(strings.Replace(askedRow, "Q-001", "Q-005", 1), "\t—\tasked", "\t2026-10-10T13:00:00Z\tclosed", 1)
	for name, c := range map[string]struct {
		read func([]byte) ([][]string, error)
		data []byte
	}{
		"meaning":  {readMeaning, table(meaningHeader, meaningRow, m2)},
		"spans":    {readSpans, table(spansHeader, headingRow, factRow)},
		"draft":    {readDraft, table(draftHeader, draftRow, "R-002\tsuccess-criterion\tNine in ten users act on it.\t3\tThe share of users who act, by the app's log.\tS-1a2b3c4d")},
		"coverage": {readCoverage, table(coverageHeader, uncoveredRow, unsupportedRow, "C-003\tclass\t4\t—\tIs fact 4 a need?\tS-1a2b3c4d")},
		"questions": {readQuestions, table(questionsHead, askedRow, intentRow, open, unanswered, closed,
			strings.Replace(rule1, "Q-001", "Q-006", 1), strings.Replace(coverage1, "Q-001", "Q-007", 1), strings.Replace(followup1, "Q-001	follow-up", "Q-008	follow-up", 1),
			strings.Replace(stack1, "Q-001", "Q-009", 1), strings.Replace(marker1, "Q-001", "Q-010", 1),
			field(questionsHead, field(questionsHead, strings.Replace(marker1, "Q-001", "Q-011", 1), "start", "—"), "end", "—"))},
		"answers":         {readAnswers, table(answersHeader, answerRow)},
		"answers-review":  {readAnswersReview, table(reviewHeader, reviewRow, "Q-002\topen\tThe band is not below U.\tS-1a2b3c4d")},
		"budget":          {readBudget, table(budgetHeader, "band.B\t50.0\tQ-002", "band.U\t200.0\tQ-003", "appetite.hours\t40.0\tQ-004")},
		"budget, in part": {readBudget, table(budgetHeader, "band.U\t200.0\tQ-003")},
	} {
		if _, err := c.read(c.data); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestIntakeRecordsRefuseEachBrokenRule(t *testing.T) {
	type c struct {
		name   string
		read   func([]byte) ([][]string, error)
		data   []byte
		column string
	}
	var cases []c
	empty := func(name string, read func([]byte) ([][]string, error), header, row string, columns ...string) {
		for _, col := range columns {
			cases = append(cases, c{name + ": the empty value in " + col, read, table(header, field(header, row, col, "—")), col})
		}
	}
	empty("meaning", readMeaning, meaningHeader, meaningRow, "id", "kind", "start", "end", "quote", "question", "session")
	empty("spans", readSpans, spansHeader, factRow, "start", "end", "text")
	empty("draft", readDraft, draftHeader, draftRow, "id", "kind", "statement", "criterion", "session")
	empty("coverage", readCoverage, coverageHeader, uncoveredRow, "id", "kind", "question", "session")
	empty("questions", readQuestions, questionsHead, askedRow, "id", "source", "ref", "block", "phase", "kind", "question", "asked", "asker", "state")
	empty("answers", readAnswers, answersHeader, answerRow, "question", "answer", "comment", "author", "time")
	empty("answers-review", readAnswersReview, reviewHeader, reviewRow, "question", "state", "reason", "session")
	empty("budget", readBudget, budgetHeader, "band.U\t200.0\tQ-003", "name", "value", "source")
	cases = append(cases,
		c{"meaning: an ID out of order", readMeaning, table(meaningHeader, field(meaningHeader, meaningRow, "id", "M-002")), "id"},
		c{"meaning: start not below end", readMeaning, table(meaningHeader, field(meaningHeader, meaningRow, "end", "10")), "end"},
		c{"meaning: offsets that differ from the quote's length", readMeaning, table(meaningHeader, field(meaningHeader, meaningRow, "end", "18")), "quote"},
		c{"meaning: a session that is not S- and 8 hexadecimal characters", readMeaning, table(meaningHeader, field(meaningHeader, meaningRow, "session", "S-1a2b3c4z")), "session"},
		c{"spans: start not below end", readSpans, table(spansHeader, field(spansHeader, factRow, "end", "12")), "end"},
		c{"spans: a fact with no class", readSpans, table(spansHeader, field(spansHeader, factRow, "class", "—")), "class"},
		c{"spans: a span that is not a fact, with a class", readSpans, table(spansHeader, field(spansHeader, headingRow, "class", "need")), "class"},
		c{"spans: a fact with a reason", readSpans, table(spansHeader, field(spansHeader, factRow, "reason", "a heading")), "reason"},
		c{"spans: a span that is not a fact, with no reason", readSpans, table(spansHeader, field(spansHeader, headingRow, "reason", "—")), "reason"},
		c{"spans: fact 0", readSpans, table(spansHeader, field(spansHeader, factRow, "fact", "0")), "fact"},
		c{"draft: an ID out of order", readDraft, table(draftHeader, field(draftHeader, draftRow, "id", "R-002")), "id"},
		c{"draft: a covered fact 0", readDraft, table(draftHeader, field(draftHeader, draftRow, "covers", "1 0")), "covers"},
		c{"draft: a session that is not S- and 8 hexadecimal characters", readDraft, table(draftHeader, field(draftHeader, draftRow, "session", "S-1a2b3c4z")), "session"},
		c{"coverage: an ID out of order", readCoverage, table(coverageHeader, unsupportedRow), "id"},
		c{"coverage: an uncovered item with no fact", readCoverage, table(coverageHeader, field(coverageHeader, uncoveredRow, "fact", "—")), "fact"},
		c{"coverage: an unsupported item with a fact", readCoverage, table(coverageHeader, field(coverageHeader, strings.Replace(unsupportedRow, "C-002", "C-001", 1), "fact", "2")), "fact"},
		c{"coverage: an uncovered item with a requirement", readCoverage, table(coverageHeader, field(coverageHeader, uncoveredRow, "requirement", "R-001")), "requirement"},
		c{"coverage: an unsupported item with no requirement", readCoverage, table(coverageHeader, field(coverageHeader, strings.Replace(unsupportedRow, "C-002", "C-001", 1), "requirement", "—")), "requirement"},
		c{"coverage: a session that is not S- and 8 hexadecimal characters", readCoverage, table(coverageHeader, field(coverageHeader, uncoveredRow, "session", "S-1a2b3c4z")), "session"},
		c{"questions: an ID out of order", readQuestions, table(questionsHead, intentRow), "id"},
		c{"questions: a start with no end", readQuestions, table(questionsHead, field(questionsHead, askedRow, "end", "—")), "end"},
		c{"questions: an end with no start", readQuestions, table(questionsHead, field(questionsHead, askedRow, "start", "—")), "start"},
		c{"questions: start not below end", readQuestions, table(questionsHead, field(questionsHead, askedRow, "end", "10")), "end"},
		c{"questions: an asked question with an answer", readQuestions, table(questionsHead, field(questionsHead, askedRow, "answered", "2026-10-10T13:00:00Z")), "answered"},
		c{"questions: an answered question with no answer", readQuestions, table(questionsHead, field(questionsHead, askedRow, "state", "answered")), "answered"},
		c{"questions: a closed question with no answer", readQuestions, table(questionsHead, field(questionsHead, askedRow, "state", "closed")), "answered"},
		c{"coverage: fact 0", readCoverage, table(coverageHeader, field(coverageHeader, uncoveredRow, "fact", "0")), "fact"},
		c{"questions: a meaning question whose ref is not M-NNN", readQuestions, table(questionsHead, field(questionsHead, askedRow, "ref", "C-001")), "ref"},
		c{"questions: an intent question whose ref is not a field", readQuestions, table(questionsHead, field(questionsHead, intent1, "ref", "band.C")), "ref"},
		c{"questions: a rule question whose ref is not Q-NNN", readQuestions, table(questionsHead, field(questionsHead, rule1, "ref", "G1:0")), "ref"},
		c{"questions: a coverage question whose ref is not C-NNN", readQuestions, table(questionsHead, field(questionsHead, coverage1, "ref", "M-001")), "ref"},
		c{"questions: a meaning question of a kind of no meaning", readQuestions, table(questionsHead, field(questionsHead, askedRow, "kind", "G2")), "kind"},
		c{"questions: a rule question of a kind that is no rule", readQuestions, table(questionsHead, field(questionsHead, rule1, "kind", "G6")), "kind"},
		c{"questions: an intent question whose kind is not intent", readQuestions, table(questionsHead, field(questionsHead, intent1, "kind", "decision")), "kind"},
		c{"questions: a coverage question of a kind of no coverage", readQuestions, table(questionsHead, field(questionsHead, coverage1, "kind", "fit")), "kind"},
		c{"questions: a rule question not asked by layup psb check", readQuestions, table(questionsHead, field(questionsHead, rule1, "asker", "layup run")), "asker"},
		c{"questions: a meaning question not asked by a session", readQuestions, table(questionsHead, field(questionsHead, askedRow, "asker", "layup run")), "asker"},
		c{"questions: an intent question not asked by layup run", readQuestions, table(questionsHead, field(questionsHead, intent1, "asker", "S-1a2b3c4d")), "asker"},
		c{"questions: an intent question with a quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, intent1, "start", "1"), "end", "5")), "start"},
		c{"questions: an unsupported item with a quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, coverage1, "kind", "unsupported"), "ref", "C-001")), "start"},
		c{"questions: a gap of meaning with no quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, askedRow, "start", "—"), "end", "—")), "start"},
		c{"questions: an uncovered item with no quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, coverage1, "start", "—"), "end", "—")), "start"},
		c{"questions: a rule gap G2 with no quote", readQuestions, table(questionsHead, field(questionsHead, rule1, "kind", "G2")), "start"},
		c{"questions: a rule gap G1 with a quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, rule1, "start", "0"), "end", "11")), "start"},
		c{"questions: a follow-up of no question before it", readQuestions, table(questionsHead, followup1), "ref"},
		c{"questions: a setup question whose ref is not S01-stack", readQuestions, table(questionsHead, field(questionsHead, stack1, "ref", "S01-name")), "ref"},
		c{"questions: a setup question whose ref is a gap of meaning's ID", readQuestions, table(questionsHead, field(questionsHead, marker1, "ref", "M-001")), "ref"},
		c{"questions: a setup question whose ref is a marker's key of 7 characters", readQuestions, table(questionsHead, field(questionsHead, marker1, "ref", "M-1a2b3c4")), "ref"},
		c{"questions: a setup question whose ref has an uppercase hexadecimal character", readQuestions, table(questionsHead, field(questionsHead, marker1, "ref", "M-1A2B3C4D")), "ref"},
		c{"questions: a setup question of a kind that is not setup", readQuestions, table(questionsHead, field(questionsHead, stack1, "kind", "intent")), "kind"},
		c{"questions: a meaning question of the kind setup", readQuestions, table(questionsHead, field(questionsHead, askedRow, "kind", "setup")), "kind"},
		c{"questions: S01-stack not asked by layup run", readQuestions, table(questionsHead, field(questionsHead, stack1, "asker", "S-1a2b3c4d")), "asker"},
		c{"questions: a marker's question not asked by a session", readQuestions, table(questionsHead, field(questionsHead, marker1, "asker", "layup run")), "asker"},
		c{"questions: S01-stack with a quote", readQuestions, table(questionsHead, field(questionsHead, field(questionsHead, stack1, "start", "1"), "end", "5")), "start"},
		c{"questions: a setup question in the idea owner's block", readQuestions, table(questionsHead, field(questionsHead, stack1, "block", "idea-owner")), "block"},
		c{"answers-review: a session that is not S- and 8 hexadecimal characters", readAnswersReview, table(reviewHeader, field(reviewHeader, reviewRow, "session", "S-1a2b3c4z")), "session"},
		c{"budget: a name that is not of the block", readBudget, table(budgetHeader, "appetite.money\t200.0\tQ-003"), "name"},
	)
	for _, k := range cases {
		refusedAt(t, k.name, k.read, k.data, k.column)
	}
}

func TestIntakeRecordsRefuseAWrongHeader(t *testing.T) {
	for name, c := range map[string]struct {
		read        func([]byte) ([][]string, error)
		header, row string
	}{
		"meaning": {readMeaning, meaningHeader, meaningRow}, "spans": {readSpans, spansHeader, factRow},
		"draft": {readDraft, draftHeader, draftRow}, "coverage": {readCoverage, coverageHeader, uncoveredRow},
		"questions": {readQuestions, questionsHead, askedRow}, "answers": {readAnswers, answersHeader, answerRow},
		"answers-review": {readAnswersReview, reviewHeader, reviewRow}, "budget": {readBudget, budgetHeader, "band.U\t200.0\tQ-003"},
	} {
		if _, err := c.read(table(strings.Replace(c.header, "\t", "\tx", 1), c.row)); err == nil {
			t.Errorf("%s: a wrong header: read, want an error", name)
		}
	}
}
