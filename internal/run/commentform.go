package run

// This file holds the parse of the three comment forms of the task loop and the
// four handoff items read from them (docs/spec/loop.md, The items a caller
// computes: The fields and the blocks of a comment; task T-4pje, #396, row 66l of
// the plan of M2e). The step that wires the parse into the loop is row 66k's.

import (
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/rules"
)

// A PlanTest is a line of the plan's tests: block, a test's ID and its path.
type PlanTest struct{ ID, Path string }

// A ReviewRecord is a parsed review record: its verdict, the head it names and
// its located findings.
type ReviewRecord struct {
	Verdict, Head string
	Findings      []Finding
}

// A commentForm is the lines of a comment that count: its fields by key, in the
// order of the comment, and its marker words with the payload of each bullet of
// the block under them, one entry for each time the marker appears.
type commentForm struct {
	fields  map[string][]string
	markers map[string][][]string
}

var (
	fieldLine = regexp.MustCompile(`^[a-z]+( [a-z]+)*: `)
	headRule  = regexp.MustCompile(`^[0-9a-f]{40}([0-9a-f]{24})?$`)
)

// emptyValue is the field value that means an empty list or no text.
const emptyValue = "—"

// readComment reads the lines of a comment: a field, a marker (exactly "tests:" or
// "findings:") with the run of "- " lines under it, and nothing else; a line
// inside a fenced block is not read. A body that is empty or not valid UTF-8 is no
// form (ok is false).
func readComment(data []byte) (f commentForm, ok bool) {
	if len(data) == 0 || !utf8.Valid(data) {
		return f, false
	}
	f = commentForm{fields: map[string][]string{}, markers: map[string][][]string{}}
	var fence string // the opener's character repeated, when inside a fence
	block := -1      // the index of the marker's entry that takes bullets, or -1
	var word string
	for _, line := range strings.Split(string(data), "\n") {
		if fence != "" {
			if strings.Trim(line, fence[:1]) == "" && len(line) >= len(fence) {
				fence = ""
			}
			continue
		}
		if c := fenceOf(line); c != "" {
			fence, block = c, -1
			continue
		}
		if payload, bullet := strings.CutPrefix(line, "- "); bullet {
			if block >= 0 {
				f.markers[word][block] = append(f.markers[word][block], payload)
			}
			continue
		}
		block = -1
		if line == "tests:" || line == "findings:" {
			word = strings.TrimSuffix(line, ":")
			f.markers[word] = append(f.markers[word], []string{})
			block = len(f.markers[word]) - 1
			continue
		}
		if loc := fieldLine.FindString(line); loc != "" {
			value := strings.TrimRight(line[len(loc):], " ")
			value = strings.TrimRight(strings.TrimSuffix(value, "\r"), " ")
			key := loc[:len(loc)-2]
			f.fields[key] = append(f.fields[key], value)
		}
	}
	return f, true
}

// fenceOf gives the run of backticks or tildes that opens a fence at the start of
// a line (three or more), or "".
func fenceOf(line string) string {
	for _, c := range []string{"`", "~"} {
		if n := len(line) - len(strings.TrimLeft(line, c)); n >= 3 {
			return strings.Repeat(c, n)
		}
	}
	return ""
}

// field gives the value of a field that appears once and is not empty.
func (f commentForm) field(key string) (string, bool) {
	v := f.fields[key]
	return strings.Join(v, ""), len(v) == 1 && v[0] != ""
}

// block gives the bullets of a marker that appears once.
func (f commentForm) block(word string) ([]string, bool) {
	b := f.markers[word]
	if len(b) != 1 {
		return nil, false
	}
	return b[0], true
}

// listField gives a list field of the plan: — is the empty list, any other text
// one element.
func (f commentForm) listField(key string) ([]string, bool) {
	v, ok := f.field(key)
	if v == emptyValue {
		return nil, ok
	}
	return []string{v}, ok
}

// ParsePlanComment parses a plan comment: the four fields once each and the tests:
// marker once, its bullets "- <id> <path>" of two words, the ID unique and the
// path the register's own path rule. ok is false for a comment that fails any part.
func ParsePlanComment(data []byte) (Plan, []PlanTest, bool) {
	f, ok := readComment(data)
	if !ok {
		return Plan{}, nil, false
	}
	var p Plan
	size, ok := f.field("size")
	if !ok || (size != "small" && size != "large") {
		return Plan{}, nil, false
	}
	p.Size = size
	for _, l := range []struct {
		key string
		to  *[]string
	}{{"open questions", &p.OpenQuestions}, {"new interfaces", &p.NewInterfaces}, {"new dependencies", &p.NewDependencies}} {
		if *l.to, ok = f.listField(l.key); !ok {
			return Plan{}, nil, false
		}
	}
	bullets, ok := f.block("tests")
	if !ok {
		return Plan{}, nil, false
	}
	var tests []PlanTest
	seen := map[string]bool{}
	for _, b := range bullets {
		w := strings.Split(b, " ")
		if len(w) != 2 || w[0] == "" || seen[w[0]] || strings.Contains(b, "\r") || !rules.ValidPath(w[1]) {
			return Plan{}, nil, false
		}
		seen[w[0]] = true
		tests = append(tests, PlanTest{w[0], w[1]})
	}
	return p, tests, true
}

// ParsePlanReview parses a plan review comment and gives its verdict word, one of
// approve, approve-with-conditions and reject.
func ParsePlanReview(data []byte) (string, bool) {
	f, ok := readComment(data)
	if !ok {
		return "", false
	}
	v, ok := f.field("verdict")
	if !ok || (v != "approve" && v != "approve-with-conditions" && v != "reject") {
		return "", false
	}
	return v, true
}

// findingOf parses a bullet "<path>:<line> <rule>: <text>" of a review record.
func findingOf(b string) (Finding, bool) {
	where, rest, ok := strings.Cut(b, " ")
	colon := strings.LastIndex(where, ":")
	if !ok || colon < 0 || strings.Contains(b, "\r") {
		return Finding{}, false
	}
	file, num := where[:colon], where[colon+1:]
	line, err := strconv.Atoi(num)
	rule, text, ok := strings.Cut(rest, ": ")
	if err != nil || strings.Trim(num, "0123456789") != "" || line < 1 || !rules.ValidPath(file) ||
		!ok || rule == "" || strings.ContainsAny(rule, " :") || text == "" {
		return Finding{}, false
	}
	return Finding{Rule: rule, File: file, Line: line, Text: text}, true
}

// ParseReviewRecord parses a review record comment: head: once (40 or 64
// lowercase hexadecimal characters), verdict: once, and the findings: marker once
// with a located finding on each bullet of its block.
func ParseReviewRecord(data []byte) (ReviewRecord, bool) {
	f, ok := readComment(data)
	if !ok {
		return ReviewRecord{}, false
	}
	var r ReviewRecord
	if r.Head, ok = f.field("head"); !ok || !headRule.MatchString(r.Head) {
		return ReviewRecord{}, false
	}
	switch r.Verdict, ok = f.field("verdict"); r.Verdict {
	case "material", "nothing material in scope", "not mergeable, findings recorded":
	default:
		return ReviewRecord{}, false
	}
	bullets, ok2 := f.block("findings")
	if !ok || !ok2 {
		return ReviewRecord{}, false
	}
	for _, b := range bullets {
		fd, ok := findingOf(b)
		if !ok {
			return ReviewRecord{}, false
		}
		r.Findings = append(r.Findings, fd)
	}
	return r, true
}

// PlanFieldsOK is the item plan-fields of a plan comment's kept copy.
func PlanFieldsOK(data []byte) bool {
	_, _, ok := ParsePlanComment(data)
	return ok
}

// VerdictFieldOK is the item verdict-field of a plan review's kept copy: the
// verdict is approve or approve-with-conditions.
func VerdictFieldOK(data []byte) bool {
	v, ok := ParsePlanReview(data)
	return ok && v != "reject"
}

// PlanReviewReasons gives the reasons of a plan review's kept copy: each of its lines that is not blank and is not the line verdict:, in order, with a
// trailing carriage return and trailing spaces left off. A copy that is no UTF-8 text has none.
func PlanReviewReasons(data []byte) []string {
	if !utf8.Valid(data) {
		return nil
	}
	var out []string
	for _, line := range strings.Split(string(data), "\n") {
		if line = strings.TrimRight(line, " \r"); line != "" && !strings.HasPrefix(line, "verdict: ") {
			out = append(out, line)
		}
	}
	return out
}

// RejectReasonsOK is the item reject-reasons of a plan review's kept copy: it parses with the verdict reject and gives at least one reason, so a
// bare reject, which gives the planner nothing to act on, is false (task T-t94t).
func RejectReasonsOK(data []byte) bool {
	v, ok := ParsePlanReview(data)
	return ok && v == "reject" && len(PlanReviewReasons(data)) > 0
}

// FindingsLocatedOK is the item findings-located of a review record's kept copy:
// the whole record parses and the block holds at least one finding.
func FindingsLocatedOK(data []byte) bool {
	r, ok := ParseReviewRecord(data)
	return ok && len(r.Findings) > 0
}

// NothingMaterialOK is the item nothing-material of a review record's kept copy
// and the head that code holds for the verifier's session: the verdict is nothing
// material in scope, the record names that head and its block is empty.
func NothingMaterialOK(data []byte, head string) bool {
	r, ok := ParseReviewRecord(data)
	return ok && r.Verdict == "nothing material in scope" && r.Head == head && len(r.Findings) == 0
}

// OpenMaterial is the predicate that fails uphill: true when the record did not
// parse, its verdict is not "nothing material in scope", or its block holds a
// finding, so no invalid or empty record reads as no open material finding.
func OpenMaterial(r ReviewRecord, ok bool) bool {
	return !ok || r.Verdict != "nothing material in scope" || len(r.Findings) > 0
}
