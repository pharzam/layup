// Package rules holds, in milestone M2b, the reader of the rule-path register
// and the check of a session's diff before a push (docs/spec/session.md, A
// workflow or rule-path change; row 32 of the plan, task T-m1dx, #158). The
// basic rule of layup/rules is in basic.go (M2e, row 63a, #352); the rest of
// layup/rules and rule batches come in M2f. The functions of this file read
// bytes and names, and start no program: the caller (internal/run) gives the
// register of the session's records commit, the names of DiffNames, and the
// diff and the head of docs/guardrails.md.
package rules

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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
// no entry of its own, it gives the entry {*.sh}. The exception that Check reads
// is the one of the row that Match gives first: with the register of S15 no
// directory row covers docs/guardrails.md, so its own row is that row.
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
		default: // a removed line, a context line, a mode change (old mode, new mode), a binary body, or other text
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
	return start, len(strings.Split(strings.TrimSuffix(string(head), "\n"), "\n")) + 1
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

// PayloadPaths gives the paths that a refused payload (git diff --binary --no-renames) changes, in the order of its headers (docs/spec/loop.md, The findings
// of a refused rule path; task T-dtw6, #491). A line that starts at the first column with "diff --git" is a candidate header and must be well formed: the
// sides a/P and b/P, each raw or quoted by git, that both decode to one non-empty P; the first malformed candidate is an error. Every other line (metadata,
// hunk lines, the base85 body of a binary patch, blank lines) is ignored, and a payload with no header is an error. The caller adds the payload's name.
func PayloadPaths(diff []byte) ([]string, error) {
	var paths []string
	for _, line := range strings.Split(string(diff), "\n") {
		if !strings.HasPrefix(line, "diff --git") {
			continue
		}
		p, err := headerPath(strings.TrimPrefix(line, "diff --git"))
		if err != nil {
			return nil, fmt.Errorf("%q: %w", line, err)
		}
		paths = append(paths, p)
	}
	if len(paths) == 0 {
		return nil, errors.New("no diff --git header")
	}
	return paths, nil
}

// headerPath reads the sides of a header after "diff --git". A path may hold spaces, so each space is tried as the separator and the first split whose two
// decoded sides are a/P and b/P of one P wins.
func headerPath(rest string) (string, error) {
	rest, ok := strings.CutPrefix(rest, " ")
	if !ok {
		return "", errors.New("no space after diff --git")
	}
	why := errors.New("not two sides")
	for i := 0; i < len(rest); i++ {
		if rest[i] != ' ' {
			continue
		}
		a, errA := side(rest[:i], "a/")
		b, errB := side(rest[i+1:], "b/")
		switch {
		case errA != nil:
			why = errA
		case errB != nil:
			why = errB
		case a != b:
			why = fmt.Errorf("the sides differ: %q and %q", a, b)
		default:
			return a, nil
		}
	}
	return "", why
}

// side decodes one side of a header: raw or quoted as git quotes a path, then the prefix; the path left must be non-empty.
func side(s, prefix string) (string, error) {
	if strings.HasPrefix(s, `"`) {
		u, err := strconv.Unquote(s)
		if err != nil {
			return "", fmt.Errorf("a malformed quoted side %s", s)
		}
		s = u
	}
	p, ok := strings.CutPrefix(s, prefix)
	if !ok || p == "" {
		return "", fmt.Errorf("the side %q is not %sP", s, prefix)
	}
	return p, nil
}
