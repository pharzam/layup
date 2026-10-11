// Package spec holds the span check of the facts of a problem statement
// (docs/spec/intake.md, The span check; task T-4ywt, #213): code, no model
// call, an engine check of NFR-005. It reads a span table by the block spans
// and writes the table of the check by the block spec-check
// (docs/spec/records.md, NFR-001 — The records of Intake). It may import
// internal/tsv only (docs/spec/packages.md, The table of M2c), so it holds a
// second Go value of the block spans, which its own test compares with the
// block; the owner of spans, and of its row rules, is internal/records.
package spec

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// CheckSchema is the form of the table of the span check: the block
// spec-check.
var CheckSchema = tsv.Schema{Name: "spec-check", Location: "stdout", Columns: []tsv.Column{
	{Name: "rule", Type: "enum(cover|bounds|text|ids|classes)", Key: true},
	{Name: "result", Type: "enum(pass|fail)"},
	{Name: "detail", Type: "text"},
}}

// SpansSchema is a second Go value of the block spans, whose owner is
// internal/records; this package reads the span table with it.
var SpansSchema = tsv.Schema{Name: "spans", Location: "records:spec/spans.tsv", Columns: []tsv.Column{
	{Name: "start", Type: "int", Key: true},
	{Name: "end", Type: "int"},
	{Name: "fact", Type: "int"},
	{Name: "class", Type: "enum(need|constraint|measure|context)"},
	{Name: "reason", Type: "text"},
	{Name: "text", Type: "text"},
}}

// The columns of a span row.
const (
	colStart = iota
	colEnd
	colFact
	colClass
	colReason
	colText
)

var classes = []string{"need", "constraint", "measure", "context"}

// Read reads a span table by SpansSchema: the header, the types and the key.
// A table that its types accept is always checked: a broken pairing of class
// and reason is a fail of the rule classes, not an error of the read.
func Read(data []byte) ([][]string, error) { return tsv.Read(data, SpansSchema) }

// WriteTSV writes the table of the span check, rows as Check gives them, by
// CheckSchema: the writer of layup spec check (task T-hpid), as internal/cli
// may not import internal/tsv.
func WriteTSV(w io.Writer, rows [][]string) error { return tsv.Write(w, CheckSchema, rows) }

// isSpace reports whether b is white space by the rule cover: ASCII space,
// tab, CR or LF.
func isSpace(b byte) bool { return b == ' ' || b == '\t' || b == '\r' || b == '\n' }

// span is a span row's offsets, and whether they are within the file and in
// order (the rule bounds).
type span struct {
	start, end int
	ok         bool
}

// Check gives the table of the span check of spans (rows as Read gives them,
// "" for the empty value) against the problem statement psb: one row per rule,
// in the order of the block, pass or fail, with the first broken place of
// that rule in detail ("" for pass). Each rule reports only its own break.
func Check(spans [][]string, psb []byte) [][]string {
	sp := make([]span, len(spans))
	for i, r := range spans {
		sp[i].start, _ = strconv.Atoi(r[colStart])
		sp[i].end, _ = strconv.Atoi(r[colEnd])
	}
	return [][]string{
		result("cover", cover(sp, psb)),
		result("bounds", bounds(sp, len(psb))),
		result("text", text(spans, sp, psb)),
		result("ids", ids(spans)),
		result("classes", classesOf(spans)),
	}
}

func result(rule, detail string) []string {
	if detail == "" {
		return []string{rule, "pass", ""}
	}
	return []string{rule, "fail", detail}
}

// cover: every byte that is not white space lies in exactly one span, and no
// span holds only white space. A span out of the file counts the bytes that
// are in it.
func cover(sp []span, psb []byte) string {
	count := make([]int, len(psb))
	for _, s := range sp {
		for i := max(s.start, 0); i < min(s.end, len(psb)); i++ {
			count[i]++
		}
	}
	for i, b := range psb {
		if !isSpace(b) && count[i] != 1 {
			if count[i] == 0 {
				return fmt.Sprintf("byte %d is in no span", i)
			}
			return fmt.Sprintf("byte %d is in %d spans", i, count[i])
		}
	}
	for _, s := range sp {
		blank := true
		for i := max(s.start, 0); i < min(s.end, len(psb)); i++ {
			blank = blank && isSpace(psb[i])
		}
		if blank && s.start < s.end && s.start < len(psb) {
			return fmt.Sprintf("the span at %d is white space only", s.start)
		}
	}
	return ""
}

// bounds: each span has 0 <= start < end <= size, and the spans are in the
// order of start with no overlap. It marks the spans that hold.
func bounds(sp []span, size int) string {
	first, prev := "", 0
	for i := range sp {
		s := &sp[i]
		var why string
		switch {
		case s.start >= s.end:
			why = fmt.Sprintf("the span at %d: start not below end %d", s.start, s.end)
		case s.end > size:
			why = fmt.Sprintf("the span at %d: end %d past the file of %d bytes", s.start, s.end, size)
		case s.start < prev:
			why = fmt.Sprintf("the span at %d: out of order, or over the span before it, which ends at %d", s.start, prev)
		}
		s.ok = why == ""
		if why == "" {
			prev = s.end
		} else if first == "" {
			first = why
		}
	}
	return first
}

// text: each span's text is the file's text at its offsets, with each tab, CR
// and LF written as one space (the field rule of the column text). A field of
// "" is the text —, which the writer writes as — and the reader gives back as
// the empty value. A span that breaks bounds has no text to compare, and the
// rule skips it (decided here, task T-4ywt): bounds reports it, one rule, one
// place, so two runs print the same bytes.
func text(spans [][]string, sp []span, psb []byte) string {
	norm := strings.NewReplacer("\t", " ", "\r", " ", "\n", " ")
	for i, r := range spans {
		if !sp[i].ok {
			continue
		}
		want := norm.Replace(string(psb[sp[i].start:sp[i].end]))
		got := r[colText]
		if got == "" {
			got = "—"
		}
		if got != want {
			return fmt.Sprintf("the span at %d: its text differs from the file's %q", sp[i].start, want)
		}
	}
	return ""
}

// ids: the fact numbers are 1 to n, each once, in the order of the spans.
func ids(spans [][]string) string {
	n := 0
	for _, r := range spans {
		if r[colFact] == "" {
			continue
		}
		n++
		if r[colFact] != strconv.Itoa(n) {
			return fmt.Sprintf("the span at %s: fact %s, want %d", r[colStart], r[colFact], n)
		}
	}
	return ""
}

// classesOf: a fact has a class of the four and no reason; a span that is not
// a fact has a reason and no class (the pairings of the block spans).
func classesOf(spans [][]string) string {
	for _, r := range spans {
		at := "the span at " + r[colStart] + ": "
		if r[colFact] != "" {
			switch {
			case !slices.Contains(classes, r[colClass]):
				return at + "a fact with no class of the four"
			case r[colReason] != "":
				return at + "a fact with a reason"
			}
			continue
		}
		switch {
		case r[colClass] != "":
			return at + "a span that is not a fact, with a class"
		case r[colReason] == "":
			return at + "a span that is not a fact, with no reason"
		}
	}
	return ""
}
