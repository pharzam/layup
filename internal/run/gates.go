package run

// This file holds the table of layup gate of a head as a record, and the status
// layup/gates (docs/spec/loop.md, layup/gates; docs/spec/records.md, the block
// task-gates; task T-q0kr, #362, row 66f of the plan of M2e).

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/tsv"
)

// GatesSchema is the Go value of the block task-gates of docs/spec/records.md;
// its owner is this package, and its test compares it with the block. It is not
// gate.ResultSchema, whose columns and keys differ (docs/spec/gate.md).
var GatesSchema = tsv.Schema{Name: "task-gates", Location: "records:tasks/<task>/gates.tsv", Columns: []tsv.Column{
	{Name: "head", Type: "sha1", Key: true},
	{Name: "kind", Type: "id(<word>)", Key: true},
	{Name: "base", Type: "sha1"},
	{Name: "state", Type: "enum(active|pending)"},
	{Name: "result", Type: "enum(pass|fail|not-active|clear)"},
	{Name: "reason", Type: "text"},
}}

func gatesPath(task string) string { return "tasks/" + task + "/gates.tsv" }

// checkHeads refuses a table whose rows of one head are not together, and, when
// baseCol is not negative, the rows of a head that have two bases. The head is
// the first column.
func checkHeads(rows [][]string, baseCol int) error {
	seen := map[string]bool{}
	for i, r := range rows {
		if i > 0 && rows[i-1][0] == r[0] {
			if baseCol >= 0 && rows[i-1][baseCol] != r[baseCol] {
				return &tsv.Error{Line: i + 2, Column: "base", Reason: "the base is the same on each row of a head"}
			}
			continue
		}
		if seen[r[0]] {
			return &tsv.Error{Line: i + 2, Column: "head", Reason: "the rows of a head are together"}
		}
		seen[r[0]] = true
	}
	return nil
}

// ReadGates reads gates.tsv by GatesSchema and the owner checks of its block
// (records.md): the rows of a head are together and have one base, and the reason
// is — exactly for a pass. An error names its line and column. A table that holds
// every kind of its base is a check of GatesTable, which reads the manifest.
func ReadGates(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, GatesSchema)
	if err != nil {
		return nil, err
	}
	if err := checkHeads(rows, 2); err != nil {
		return nil, err
	}
	for i, r := range rows {
		if (r[4] == "pass") != (r[5] == "") {
			return nil, &tsv.Error{Line: i + 2, Column: "reason", Reason: "the reason is — exactly for a pass"}
		}
	}
	return rows, nil
}

// WriteGates writes rows by GatesSchema. It writes nothing for rows that ReadGates refuses.
func WriteGates(rows [][]string) ([]byte, error) {
	data, err := table(GatesSchema, rows)
	if err != nil {
		return nil, err
	}
	if _, err := ReadGates(data); err != nil {
		return nil, err
	}
	return data, nil
}

// gateRows gives the rows of gates.tsv of a table of layup gate.
func gateRows(t gate.Table) [][]string {
	var rows [][]string
	for _, r := range t.Rows {
		reason := r.Reason
		if r.Result == "pass" {
			reason = ""
		}
		rows = append(rows, []string{t.Head, r.Kind, t.Base, r.State, r.Result, reason})
	}
	return rows
}

// gatesComplete says whether the rows of one head hold exactly the kinds and the
// states of the manifest, in its order: a table that omits a kind never gives a
// success (NFR-004).
func gatesComplete(rows [][]string, manifest []byte) error {
	kinds, err := gate.ReadManifest(manifest)
	if err != nil {
		return fmt.Errorf("docs/gates.tsv at the base: %w", err)
	}
	if len(rows) != len(kinds) {
		return fmt.Errorf("the table has %d rows and the manifest %d kinds", len(rows), len(kinds))
	}
	for i, k := range kinds {
		if rows[i][1] != k.Name || rows[i][3] != k.State {
			return fmt.Errorf("row %d is %s (%s); the manifest has %s (%s)", i+1, rows[i][1], rows[i][3], k.Name, k.State)
		}
	}
	return nil
}

// GatesStatus gives the status layup/gates of the table of a head (loop.md,
// layup/gates): success when every row passes or is clear, with the reason of
// the first clear row; else failure with the kind, the result and the reason of
// the first row, in the order of the manifest, that is neither. A table with no
// row is an error, never a pass.
func GatesStatus(t gate.Table) Status {
	if len(t.Rows) == 0 {
		return Status{GatesContext, "error", "error: the table has no row"}
	}
	clear := ""
	for _, r := range t.Rows {
		switch r.Result {
		case "pass":
		case "clear":
			if clear == "" {
				clear = r.Reason
			}
		default:
			return Status{GatesContext, "failure", fit(r.Kind + ": " + r.Result + ": " + r.Reason)}
		}
	}
	if clear != "" {
		return Status{GatesContext, "success", fit("clear: " + clear)}
	}
	return Status{GatesContext, "success", "pass"}
}

// A GatesInput names the head whose table is asked, and where the progress of the
// run and the held output of its commands go (the step and the output of
// gate.Run; os.Stderr and a line for each kind when nil).
type GatesInput struct {
	Task, Base, Head string
	Step             func(i, n int, kind string) func()
	Out              io.Writer
}

// manifest gives docs/gates.tsv at the base, nil for none.
func (s *Sessions) manifest(base string) ([]byte, error) {
	if s.manifestAt != nil {
		return s.manifestAt(base)
	}
	return s.blobAt(base, "docs/gates.tsv")
}

// GatesTable gives the recorded table of a head (rebuilt from gates.tsv, with the
// check that it holds every kind of the manifest at its base), and whether the
// head has one: what GateFindings and NextAttempt take after a restart. A table
// that its reader refuses, or that is not complete, is an error.
func (s *Sessions) GatesTable(ctx context.Context, task, head string) (gate.Table, bool, error) {
	data, err := s.Store.ReadFile(ctx, gatesPath(task))
	if err != nil || data == nil {
		return gate.Table{}, false, err
	}
	rows, err := ReadGates(data)
	if err != nil {
		return gate.Table{}, false, fmt.Errorf("%s: %w", gatesPath(task), err)
	}
	var mine [][]string
	for _, r := range rows {
		if r[0] == head {
			mine = append(mine, r)
		}
	}
	if len(mine) == 0 {
		return gate.Table{}, false, nil
	}
	m, err := s.manifest(mine[0][2])
	if err != nil {
		return gate.Table{}, false, err
	}
	if m == nil {
		return gate.Table{}, false, fmt.Errorf("%s has the table of %s, and the base %s has no docs/gates.tsv", gatesPath(task), head, mine[0][2])
	}
	if err := gatesComplete(mine, m); err != nil {
		return gate.Table{}, false, fmt.Errorf("%s, the table of %s: %w", gatesPath(task), head, err)
	}
	t := gate.Table{Base: mine[0][2], Head: head}
	for _, r := range mine {
		t.Rows = append(t.Rows, gate.Row{Kind: r[1], State: r[3], Result: r[4], Reason: r[5]})
	}
	return t, true, nil
}

// gateProgress is the step of a run with no step of the caller: one line when a
// kind starts and one when it ends, with the count and the time it took.
func gateProgress(out io.Writer) func(i, n int, kind string) func() {
	return func(i, n int, kind string) func() {
		start := time.Now()
		fmt.Fprintf(out, "layup run: gate [%d/%d] %s…\n", i, n, kind)
		return func() {
			fmt.Fprintf(out, "layup run: gate [%d/%d] %s done (%s)\n", i, n, kind, time.Since(start).Round(time.Millisecond))
		}
	}
}

// GatesAt gives the status layup/gates of a head. A head that has a table in
// gates.tsv is never run again: the status is rebuilt from its rows (and a table
// of another base is an error). Otherwise layup gate runs once, in a scratch tree
// with the base's gate files (gate.Run), and its rows are committed in one fenced
// records commit before the status is given, so a stop between the commit and the
// post posts from the record. An input error of the gate (a revision, the manifest
// at the base, git, sh) is the status error with no row written. A scratch tree
// that cannot be removed comes with the whole table: the rows are written, and the
// status is given with the error.
func (s *Sessions) GatesAt(ctx context.Context, in GatesInput) (Status, error) {
	if !forge.FullSHA(in.Base) || !forge.FullSHA(in.Head) {
		return Status{}, fmt.Errorf("the base and the head of layup/gates are full SHAs")
	}
	tbl, found, err := s.GatesTable(ctx, in.Task, in.Head)
	if err != nil {
		return Status{}, err
	}
	if found {
		if tbl.Base != in.Base {
			return Status{}, fmt.Errorf("%s has the table of %s at the base %s, not %s", gatesPath(in.Task), in.Head, tbl.Base, in.Base)
		}
		return GatesStatus(tbl), nil
	}
	out := in.Out
	if out == nil {
		out = os.Stderr
	}
	step := in.Step
	if step == nil {
		step = gateProgress(out)
	}
	t, runErr := gate.Run(s.Clone, in.Base, in.Head, step, out)
	var input *gate.InputError
	var cleanup *gate.CleanupError
	switch {
	case errors.As(runErr, &input):
		first, _, _ := strings.Cut(strings.TrimSpace(input.Error()), "\n")
		return Status{GatesContext, "error", fit("error: " + first)}, nil
	case runErr != nil && !errors.As(runErr, &cleanup):
		return Status{}, runErr
	}
	rows, err := s.readTable(ctx, gatesPath(in.Task), ReadGates)
	if err != nil {
		return Status{}, fmt.Errorf("%s: %w", gatesPath(in.Task), err)
	}
	data, err := WriteGates(append(rows, gateRows(t)...))
	if err != nil {
		return Status{}, err
	}
	if err := s.commit(ctx, map[string][]byte{gatesPath(in.Task): data}, "layup run: the gate table of "+in.Head+" of "+in.Task); err != nil {
		return Status{}, err
	}
	return GatesStatus(t), runErr
}
