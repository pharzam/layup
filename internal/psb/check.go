// Package psb checks a problem statement for gaps before delivery starts, so the
// idea owner answers them in one batch (PSB §6 "Problem Statement Quality",
// Decision Point 2). Every rule is deterministic: text in, gaps out (ADR-0011).
package psb

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

// Gap is one question for the idea owner.
type Gap struct {
	Rule     string // G1 … G5
	Line     int    // 1-based line; 0 when the gap is an absence
	Excerpt  string // the trimmed line, TABs as spaces, at most 80 runes
	Question string
}

var (
	// G1: "technology stack", then only spaces or '*', then ':' and a value.
	namedStack = regexp.MustCompile(`(?i)technology stack[ *]*:[ *]*\S`)
	// G3: 2 to 6 capitals (digits after the first), with ASCII word boundaries.
	abbrev = regexp.MustCompile(`(^|[^A-Za-z0-9_])([A-Z][A-Z0-9]{1,5})($|[^A-Za-z0-9_])`)
	bold   = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	// G4: the kit's vague words ("One reading, not two"), whole word, any case.
	vague = regexp.MustCompile(`(?i)(^|[^A-Za-z0-9_])(fast|robust|soon|clean|better|handle)($|[^A-Za-z0-9_])`)
	digit = regexp.MustCompile(`[0-9]`)
	// G5: a target that a pilot must still calibrate.
	startValue = regexp.MustCompile(`(?i)\(start value\)`)
)

var noMeasure = map[string]bool{"": true, "—": true, "-": true, "tbd": true, "not measured": true}

// Check returns the gaps of a Markdown problem statement, sorted by line, then
// by rule, then by their order on the line.
func Check(text string) []Gap {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var gaps []Gap
	add := func(rule string, n int, q string) {
		ex := "—"
		if n > 0 {
			ex = excerpt(lines[n-1])
		}
		gaps = append(gaps, Gap{Rule: rule, Line: n, Excerpt: ex, Question: q})
	}

	stack := false
	for _, l := range lines {
		if namedStack.MatchString(l) {
			stack = true
		}
	}
	if !stack {
		add("G1", 0, "Which technology stack does the product use (languages, frameworks, tools)?")
	}

	defined := definedTerms(lines)
	seen := map[string]bool{}
	col := -1 // index of the Measurement or Verification column in the current table
	for i, l := range lines {
		n := i + 1
		cells, isRow := tableCells(l)
		switch {
		case !isRow:
			col = -1
		case isSeparator(cells):
		case col < 0:
			col = measureColumn(cells) // a header row; -2 when it has no such column
			if col == -1 {
				col = -2
			}
		case col >= 0 && col < len(cells) && noMeasure[strings.ToLower(strings.TrimSuffix(strings.TrimSpace(cells[col]), "."))]:
			add("G2", n, "How is this metric measured? The row gives no measurement method.")
		}
		for _, m := range allAbbrevs(l) {
			if !defined[m] && !seen[m] {
				seen[m] = true
				add("G3", n, fmt.Sprintf("What does %q mean? The terms table does not define it.", m))
			}
		}
		if m := vague.FindStringSubmatch(l); m != nil && !digit.MatchString(l) {
			add("G4", n, fmt.Sprintf("Which number or threshold does %q stand for here?", strings.ToLower(m[2])))
		}
		if startValue.MatchString(l) {
			add("G5", n, "Which start value does the pilot use for this target, and who sets it?")
		}
	}
	sort.SliceStable(gaps, func(a, b int) bool {
		if gaps[a].Line != gaps[b].Line {
			return gaps[a].Line < gaps[b].Line
		}
		return gaps[a].Rule < gaps[b].Rule
	})
	return gaps
}

// WriteTSV writes the batch: a header row, then one row per gap with ids
// Q-001, Q-002, … in the order of gaps.
func WriteTSV(w io.Writer, gaps []Gap) {
	fmt.Fprint(w, "id\trule\tline\texcerpt\tquestion\n")
	for i, g := range gaps {
		fmt.Fprintf(w, "Q-%03d\t%s\t%d\t%s\t%s\n", i+1, g.Rule, g.Line, g.Excerpt, g.Question)
	}
}

func excerpt(l string) string {
	s := strings.TrimSpace(strings.ReplaceAll(l, "\t", " "))
	if utf8.RuneCountInString(s) > 80 {
		s = string([]rune(s)[:80])
	}
	return s
}

func tableCells(l string) ([]string, bool) {
	t := strings.TrimSpace(l)
	if !strings.HasPrefix(t, "|") {
		return nil, false
	}
	t = strings.TrimSuffix(strings.TrimPrefix(t, "|"), "|")
	return strings.Split(t, "|"), true
}

func isSeparator(cells []string) bool {
	for _, c := range cells {
		if strings.Trim(strings.TrimSpace(c), ":-") != "" {
			return false
		}
	}
	return true
}

func measureColumn(header []string) int {
	for i, c := range header {
		switch strings.ToLower(strings.TrimSpace(c)) {
		case "measurement", "verification":
			return i
		}
	}
	return -1
}

func allAbbrevs(l string) []string {
	var out []string
	for rest := l; ; {
		loc := abbrev.FindStringSubmatchIndex(rest)
		if loc == nil {
			return out
		}
		out = append(out, rest[loc[4]:loc[5]])
		rest = rest[loc[5]:] // the boundary character may start the next match
	}
}

// definedTerms reads the first table under a heading that contains "Terms":
// each capital word inside a bold span of a row's first cell is defined.
func definedTerms(lines []string) map[string]bool {
	def := map[string]bool{}
	under, inTable := false, false
	for _, l := range lines {
		t := strings.TrimSpace(l)
		if strings.HasPrefix(t, "#") {
			if inTable {
				return def
			}
			under = strings.Contains(t, "Terms")
			continue
		}
		cells, isRow := tableCells(l)
		if !under || !isRow {
			if inTable {
				return def
			}
			continue
		}
		inTable = true
		for _, b := range bold.FindAllStringSubmatch(cells[0], -1) {
			for _, m := range allAbbrevs(b[1]) {
				def[m] = true
			}
		}
	}
	return def
}
