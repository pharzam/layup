// Package handoff holds the transition table of the task loop and the check
// of a handoff (docs/spec/loop.md, "REQ-005 — The handoff"; task T-h17s, #348,
// row 61 of the plan of M2e). It is code and no more: no forge call, no
// session and no program. It may import internal/tsv and internal/records only
// (docs/spec/packages.md, "The table of M2e"). The caller, internal/run, makes
// the facts that a check needs, and hands them in as a Handoff.
package handoff

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// A Row is one allowed handoff: the step it leaves, the step it enters, the
// kinds of artifact it names and the items that code checks (the block
// transitions).
type Row struct {
	From, To         string
	Artifacts, Items []string
}

// Schema is the Go value of the block transitions of docs/spec/records.md;
// its owner is this package, and its test compares it with the block.
var Schema = tsv.Schema{Name: "transitions", Location: "records:transitions.tsv", Columns: []tsv.Column{
	{Name: "from", Type: "text", Key: true},
	{Name: "to", Type: "text", Key: true},
	{Name: "artifacts", Type: "list(text)"},
	{Name: "items", Type: "list(text)"},
}}

// The columns of a row of the table.
const (
	colFrom = iota
	colTo
	colArtifacts
	colItems
)

// The kinds of artifact and the names of the items of docs/spec/loop.md,
// "The transition table", in the order of the specification.
var (
	kinds = []string{"plan", "plan-review", "review-record", "task-file", "commits", "rows", "patch", "rendering", "question", "answer"}
	items = []string{"schema", "artifacts", "plan-fields", "verdict-field", "tests-red-green", "spec-sections", "patch-per-kind", "register-rows", "rendering-equal", "nothing-material", "findings-located", "heading-ids", "question-row", "answer-row"}
)

// Kinds returns the kinds of artifact that a row may name.
func Kinds() []string { return slices.Clone(kinds) }

// Items returns the names of the items that a row may name.
func Items() []string { return slices.Clone(items) }

// Default returns the thirteen rows that the binary embeds: the twelve of
// docs/architecture.md §8 and the row of the target's architecture task, in the
// order of docs/spec/loop.md, "The checks of a handoff". Each call gives its
// own copy.
func Default() []Row {
	row := func(from, to string, artifacts, items []string) Row { return Row{from, to, artifacts, items} }
	w := func(s ...string) []string { return s }
	return []Row{
		row("task plan", "plan review", w("plan"), w("schema", "artifacts", "plan-fields")),
		row("plan review", "developer", w("plan-review"), w("schema", "artifacts", "verdict-field")),
		row("developer (build task)", "verifier", w("commits", "task-file"), w("schema", "artifacts", "tests-red-green")),
		row("developer (specification)", "verifier", w("commits", "task-file"), w("schema", "artifacts", "spec-sections")),
		row("developer (architecture)", "verifier", w("commits", "task-file"), w("schema", "artifacts")),
		row("developer (rule batch)", "verifier", w("commits", "task-file", "patch"), w("schema", "artifacts", "patch-per-kind")),
		row("developer (milestone plan)", "verifier", w("rows"), w("schema", "artifacts", "register-rows")),
		row("code (rendered task)", "verifier", w("rendering"), w("artifacts", "rendering-equal")),
		row("verifier", "close-out", w("review-record"), w("schema", "artifacts", "nothing-material")),
		row("verifier", "developer", w("review-record"), w("schema", "artifacts", "findings-located")),
		row("specification", "milestone plan", w("task-file"), w("artifacts", "heading-ids")),
		row("asking role", "owner role", w("question"), w("schema", "question-row")),
		row("owner role", "asking role", w("answer"), w("schema", "answer-row")),
	}
}

// InForce gives the rows that a read of the table checks and applies.
func InForce(rows []Row) []Row { return slices.Clone(rows) }

// WriteTable writes the table by Schema (the block transitions).
func WriteTable(w io.Writer, rows []Row) error {
	out := make([][]string, len(rows))
	for i, r := range rows {
		a, err := tsv.JoinList(r.Artifacts)
		if err != nil {
			return fmt.Errorf("row %d: artifacts: %w", i+1, err)
		}
		it, err := tsv.JoinList(r.Items)
		if err != nil {
			return fmt.Errorf("row %d: items: %w", i+1, err)
		}
		out[i] = []string{r.From, r.To, a, it}
	}
	return tsv.Write(w, Schema, out)
}

// ReadTable reads a table by Schema and refuses, naming the line, a row whose
// artifacts or items are empty or hold a word outside the lists of the
// specification. It does not check that the table has each default pair: that
// is CheckTable.
func ReadTable(data []byte) ([]Row, error) {
	recs, err := tsv.Read(data, Schema)
	if err != nil {
		return nil, err
	}
	rows := make([]Row, len(recs))
	for i, f := range recs {
		r := Row{From: f[colFrom], To: f[colTo], Artifacts: strings.Fields(f[colArtifacts]), Items: strings.Fields(f[colItems])}
		if err := checkRow(r); err != nil {
			return nil, fmt.Errorf("line %d: %w", i+2, err)
		}
		rows[i] = r
	}
	return rows, nil
}

// checkRow checks the words of one row: neither list is empty, and each word is
// of its list.
func checkRow(r Row) error {
	if len(r.Artifacts) == 0 {
		return fmt.Errorf("column artifacts: a row names at least one kind of artifact (%s -> %s)", r.From, r.To)
	}
	if len(r.Items) == 0 {
		return fmt.Errorf("column items: a row names at least one item (%s -> %s)", r.From, r.To)
	}
	for _, k := range r.Artifacts {
		if !slices.Contains(kinds, k) {
			return fmt.Errorf("column artifacts: %q is not a kind of artifact (%s -> %s)", k, r.From, r.To)
		}
	}
	for _, it := range r.Items {
		if !slices.Contains(items, it) {
			return fmt.Errorf("column items: %q is not an item (%s -> %s)", it, r.From, r.To)
		}
	}
	return nil
}

// renderedFrom is the step that code itself writes the handoff of: the rendered
// task has no session, so no typed result, and a row of it cannot list "schema".
const renderedFrom = "code (rendered task)"

// CheckTable checks a table, the default or an Operator's replacement: each
// row has words of the lists, no pair is listed twice, and each default pair
// has its row. A pair that is not a default pair is allowed, with words of the
// lists. The one other rule is for the row of code's own handoff, which lists
// no "schema" (Check could never find one: the handoff would be invalid at each
// attempt); another row keeps its "schema". It is the check that a replacement
// by "set transitions" passes, and that the register applies to the table it
// reads (docs/spec/loop.md, The transition table; task T-dx1h, #358).
func CheckTable(rows []Row) error {
	seen := map[[2]string]bool{}
	for i, r := range rows {
		if err := checkRow(r); err != nil {
			return fmt.Errorf("row %d: %w", i+1, err)
		}
		if r.From == renderedFrom && slices.Contains(r.Items, "schema") {
			return fmt.Errorf("row %d: the row of %s (%s -> %s) lists schema, but code's own handoff has no typed result", i+1, renderedFrom, r.From, r.To)
		}
		p := [2]string{r.From, r.To}
		if seen[p] {
			return fmt.Errorf("row %d: the pair %s -> %s is listed twice", i+1, r.From, r.To)
		}
		seen[p] = true
	}
	for _, d := range Default() {
		if !seen[[2]string{d.From, d.To}] {
			return fmt.Errorf("the table has no row for the pair %s -> %s", d.From, d.To)
		}
	}
	return nil
}
