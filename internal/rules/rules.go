// Package rules holds, in milestone M2b, the reader of the rule-path register
// and the check of a session's diff before a push (docs/spec/session.md, A
// workflow or rule-path change; row 32 of the plan, task T-m1dx, #158). The
// check layup/rules and rule batches come in M2f. Its functions read bytes and
// names, and start no program: the caller (internal/run) gives the register of
// the session's records commit, the names of DiffNames, and the diff and the
// head of docs/guardrails.md.
package rules

import (
	"bufio"
	"bytes"

	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// RulePathsSchema is the form of the rule-path register, the block rule-paths
// of docs/spec/setup.md. Its owner is internal/setup, which writes the
// register; this package may not import it, so it holds a second Go value of
// the block, which its own integration test compares with the block.
var RulePathsSchema = tsv.Schema{Name: "rule-paths", Location: "records:rule-paths.tsv", Columns: []tsv.Column{
	{Name: "pattern", Type: "text", Key: true},
	{Name: "exception", Type: "text"},
	{Name: "source", Type: "enum(baseline|catalog|architecture)"},
}}

const (
	// GuardrailsPath is the one path with an exception.
	GuardrailsPath = "docs/guardrails.md"
	// ExceptionGuardrails is the text of its exception in the register.
	ExceptionGuardrails = "added lines in section 2"
	workflows           = ".github/workflows/"
)

// An Entry is one row of the register: a pattern, and its exception ("" for
// none).
type Entry struct{ Pattern, Exception string }

// ReadRegister reads rule-paths.tsv by its schema, and refuses an exception on
// any row but docs/guardrails.md, and any other exception there (the plan
// review of #158, condition 2). An error names its line.
func ReadRegister(data []byte) ([]Entry, error) {
	rows, err := tsv.Read(data, RulePathsSchema)
	if err != nil {
		return nil, err
	}
	reg := make([]Entry, 0, len(rows))
	for i, r := range rows {
		e := Entry{Pattern: r[0], Exception: r[1]}
		if e.Exception != "" && (e.Pattern != GuardrailsPath || e.Exception != ExceptionGuardrails) {
			return nil, &tsv.Error{Line: i + 2, Column: "exception", Reason: "the one exception is " + ExceptionGuardrails + ", of " + GuardrailsPath}
		}
		reg = append(reg, e)
	}
	return reg, nil
}

// Match gives the first entry of the register that matches path: a pattern
// that ends with / matches each path under it, another pattern the path
// itself. Any path whose name ends with .sh matches, as the register's entry
// "each file of the tree whose name ends with .sh" names a kind of file; with
// no entry of its own, it gives the entry {*.sh}.
func Match(reg []Entry, path string) (Entry, bool) {
	for _, e := range reg {
		if (strings.HasSuffix(e.Pattern, "/") && strings.HasPrefix(path, e.Pattern)) || e.Pattern == path {
			return e, true
		}
	}
	if strings.HasSuffix(path, ".sh") {
		return Entry{Pattern: "*.sh"}, true
	}
	return Entry{}, false
}

// GuardrailsAdditions reports whether a diff of docs/guardrails.md (git diff
// -U0) only adds lines inside its section 2: the diff is read as hunks, holds
// at least one added line and no removed one, and each added line lies at the
// head strictly between the line that starts "## 2." and the next line that
// starts "## ". A diff with no hunk, a binary diff, or other text is not the
// exception (the plan review of #158, condition 1).
func GuardrailsAdditions(diff, head []byte) bool {
	start, end := section2(head)
	if start == 0 {
		return false
	}
	added, inHunk, next := 0, false, 0
	sc := bufio.NewScanner(bytes.NewReader(diff))
	sc.Buffer(make([]byte, 0, 64*1024), 16*1024*1024)
	for sc.Scan() {
		line := sc.Text()
		switch {
		case strings.HasPrefix(line, "@@ "):
			first, ok := hunkStart(line)
			if !ok {
				return false
			}
			inHunk, next = true, first
		case !inHunk && (strings.HasPrefix(line, "diff --git ") || strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "--- ") || strings.HasPrefix(line, "+++ ")):
		case inHunk && strings.HasPrefix(line, "+"):
			if next <= start || next >= end {
				return false
			}
			added, next = added+1, next+1
		case inHunk && strings.HasPrefix(line, `\ `): // "\ No newline at end of file"
		default: // a removed line, a context line, a mode change, a binary body, or other text
			return false
		}
	}
	return sc.Err() == nil && added > 0
}

// section2 gives the line numbers (from 1) of the line that starts "## 2."
// and of the next line that starts "## " (one past the last line if none);
// start 0 when the head has no such line.
func section2(head []byte) (start, end int) {
	for i, l := range strings.Split(string(head), "\n") {
		switch {
		case start == 0 && strings.HasPrefix(l, "## 2."):
			start = i + 1
		case start != 0 && strings.HasPrefix(l, "## "):
			return start, i + 1
		}
	}
	return start, strings.Count(string(head), "\n") + 2
}

// hunkStart reads the first line of the new side of a hunk header
// "@@ -a[,b] +c[,d] @@".
func hunkStart(header string) (int, bool) {
	f := strings.Fields(header)
	if len(f) < 4 || f[0] != "@@" || f[3] != "@@" || !strings.HasPrefix(f[2], "+") {
		return 0, false
	}
	c, _, _ := strings.Cut(strings.TrimPrefix(f[2], "+"), ",")
	n, err := strconv.Atoi(c)
	return n, err == nil
}

// A Change is the diff (git diff -U0) and the head of docs/guardrails.md; the
// zero Change is no diff given.
type Change struct{ Diff, Head []byte }

// Check gives the first refusal of a session's diff, by the names of
// DiffNames: first "workflow" for a path under .github/workflows/ over every
// path, then "rule-path" for a path that the register matches over every
// path, in the order of the names within each, as the specification gives the
// two sentences; docs/guardrails.md passes when its entry has the exception
// and guardrails only adds lines in its section 2. No refusal gives "" and "".
func Check(names []string, reg []Entry, guardrails Change) (reason, path string) {
	for _, p := range names {
		if strings.HasPrefix(p, workflows) {
			return "workflow", p
		}
	}
	for _, p := range names {
		e, ok := Match(reg, p)
		if !ok {
			continue
		}
		if p == GuardrailsPath && e.Exception == ExceptionGuardrails && GuardrailsAdditions(guardrails.Diff, guardrails.Head) {
			continue
		}
		return "rule-path", p
	}
	return "", ""
}
