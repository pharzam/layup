package spec

import (
	"strings"
	"testing"
)

// The span check (docs/spec/intake.md, The span check; task T-4ywt, #213): a
// span table that covers a small file passes each rule, and each case that
// breaks a rule fails that rule's row, with the first broken place.

// The file: a heading, a blank line and one sentence. Its bytes:
// "# Brief\n" is 0 to 7 (with the line feed at 7), "\n" is 8, and
// "The app sends one reminder.\n" is 9 to 36 (the line feed at 36).
const psb = "# Brief\n\nThe app sends one reminder.\n"

// row is a span row: start, end, fact, class, reason, text ("" for —).
func row(start, end, fact, class, reason, text string) []string {
	return []string{start, end, fact, class, reason, text}
}

func good() [][]string {
	return [][]string{
		row("0", "7", "", "", "a heading", "# Brief"),
		row("9", "36", "1", "need", "", "The app sends one reminder."),
	}
}

// ruleOf gives the result and the detail of the rule's row of table.
func ruleOf(t *testing.T, table [][]string, rule string) (string, string) {
	t.Helper()
	for _, r := range table {
		if r[0] == rule {
			return r[1], r[2]
		}
	}
	t.Fatalf("no row %s in %q", rule, table)
	return "", ""
}

func TestASpanTableThatCoversTheFilePasses(t *testing.T) {
	table := Check(good(), []byte(psb))
	var rules []string
	for _, r := range table {
		rules = append(rules, r[0])
		if r[1] != "pass" || r[2] != "" {
			t.Errorf("%s: %s %q, want pass with no detail", r[0], r[1], r[2])
		}
	}
	if strings.Join(rules, " ") != "cover bounds text ids classes" {
		t.Errorf("the rules %q, want cover bounds text ids classes, in that order", rules)
	}
}

func TestEachRuleFailsOnItsOwnBreak(t *testing.T) {
	set := func(i, col int, v string) [][]string { s := good(); s[i][col] = v; return s }
	for _, c := range []struct {
		name, rule string
		spans      [][]string
		detail     string // a part of the detail
	}{
		{"a byte in no span", "cover", set(1, 1, "30"), "in no span"},
		{"a byte in two spans", "cover", [][]string{row("0", "7", "", "", "a heading", "# Brief"), row("6", "36", "1", "need", "", "f\n\nThe app sends one reminder.")}, "in 2 spans"},
		{"a span of white space only", "cover", [][]string{row("0", "7", "", "", "a heading", "# Brief"), row("7", "9", "", "", "a gap", "  "), row("9", "36", "1", "need", "", "The app sends one reminder.")}, "white space only"},
		{"an end past the file", "bounds", set(1, 1, "40"), "past the file"},
		{"a start not below the end", "bounds", set(1, 1, "9"), "not below"},
		{"two spans out of order", "bounds", [][]string{good()[1], good()[0]}, "out of order"},
		{"a text that differs from the file", "text", set(1, 5, "The app sends two reminders."), "differs"},
		{"a fact number that skips", "ids", set(1, 2, "2"), "want 1"},
		{"a fact with no class", "classes", set(1, 3, ""), "class"},
		{"a fact with a reason", "classes", set(1, 4, "a heading"), "reason"},
		{"a span that is not a fact, with a class", "classes", set(0, 3, "need"), "class"},
		{"a span that is not a fact, with no reason", "classes", set(0, 4, ""), "reason"},
	} {
		table := Check(c.spans, []byte(psb))
		got, detail := ruleOf(t, table, c.rule)
		if got != "fail" || !strings.Contains(detail, c.detail) {
			t.Errorf("%s: %s %s %q, want fail with %q", c.name, c.rule, got, detail, c.detail)
		}
	}
}

func TestAFactNumberThatRepeats(t *testing.T) {
	psb := "One.\nTwo.\n"
	spans := [][]string{row("0", "4", "1", "need", "", "One."), row("5", "9", "1", "need", "", "Two.")}
	if got, detail := ruleOf(t, Check(spans, []byte(psb)), "ids"); got != "fail" || !strings.Contains(detail, "want 2") {
		t.Errorf("a repeated fact number: %s %q, want fail with want 2", got, detail)
	}
}

func TestATextOfOneEmDash(t *testing.T) {
	// The text — of a span is written as —, which reads back as the empty
	// value (docs/spec/README.md), so the rule text compares an empty field
	// as the text —.
	psb := "—\n"
	spans := [][]string{row("0", "3", "", "", "a rule line", "")}
	if got, detail := ruleOf(t, Check(spans, []byte(psb)), "text"); got != "pass" {
		t.Errorf("a span whose text is one em dash: %s %q, want pass", got, detail)
	}
}

func TestReadTakesASpanTableByItsBlock(t *testing.T) {
	data := "start\tend\tfact\tclass\treason\ttext\n0\t7\t—\t—\ta heading\t# Brief\n9\t36\t1\tneed\t—\tThe app sends one reminder.\n"
	rows, err := Read([]byte(data))
	if err != nil || len(rows) != 2 {
		t.Fatalf("Read: %d rows, %v; want 2", len(rows), err)
	}
	if _, err := Read([]byte(strings.Replace(data, "reason", "why", 1))); err == nil {
		t.Error("a wrong header: read, want an error")
	}
}

func TestASpanAcrossALineFeed(t *testing.T) {
	// The field rule of the column text writes each tab, CR and LF as one
	// space, so the rule text compares after the same replacement.
	psb := "One fact\tacross\na line.\n"
	spans := [][]string{row("0", "23", "1", "need", "", "One fact across a line.")}
	if got, detail := ruleOf(t, Check(spans, []byte(psb)), "text"); got != "pass" {
		t.Errorf("a span across a tab and a line feed: %s %q, want pass", got, detail)
	}
}
