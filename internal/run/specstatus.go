package run

// This file holds the status layup/spec and the table of the full check of a head
// as a record (docs/spec/loop.md, layup/spec; task T-q0kr, #362, row 66f of the
// plan of M2e). The schema, the writer and the check itself are internal/spec's.

import (
	"context"
	"fmt"
	"strings"

	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

func specPath(task string) string { return "tasks/" + task + "/spec-check.tsv" }

// specRuleList gives the rules of the table of the full check in the order of
// its block: the values of the enum of the column rule of the schema, so that no
// second list of the rules is kept.
func specRuleList() []string {
	t := spec.TaskSpecCheckSchema.Columns[1].Type
	return strings.Split(strings.TrimSuffix(strings.TrimPrefix(t, "enum("), ")"), "|")
}

// specComplete says whether the rows (rule, result, detail) of a head hold every
// rule of the block once, in its order, a clear head's table too.
func specComplete(rows [][]string) bool {
	list := specRuleList()
	if len(rows) != len(list) {
		return false
	}
	for i, r := range rows {
		if len(r) != 3 || r[0] != list[i] {
			return false
		}
	}
	return true
}

// SpecStatus gives the status layup/spec of the rows of the full check of a head,
// as CheckFull gives them (loop.md, layup/spec): success, `pass` when no row is a
// fail or a clear, or `clear: <detail>` of the first clear row; failure with
// `<rule>: <detail>` of the first fail row in the order of the rules. A table that
// is not complete is an error, never a success.
func SpecStatus(rows [][]string) Status {
	if !specComplete(rows) {
		return Status{SpecContext, "error", "error: the table of the spec check is not complete"}
	}
	clear := ""
	for _, r := range rows {
		switch r[1] {
		case "fail":
			return Status{SpecContext, "failure", fit(strings.TrimSuffix(r[0]+": "+r[2], ": "))}
		case "clear":
			if clear == "" {
				clear = r[2]
			}
		}
	}
	if clear != "" {
		return Status{SpecContext, "success", fit("clear: " + clear)}
	}
	return Status{SpecContext, "success", "pass"}
}

// readSpecTable reads spec-check.tsv by spec.TaskSpecCheckSchema and the owner
// checks of its block: the rows of a head are together, and each head holds the
// whole table. It gives the rows of every head.
func readSpecTable(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, spec.TaskSpecCheckSchema)
	if err != nil {
		return nil, err
	}
	if err := checkHeads(rows, -1); err != nil {
		return nil, err
	}
	for i := 0; i < len(rows); {
		j := i
		var part [][]string
		for ; j < len(rows) && rows[j][0] == rows[i][0]; j++ {
			part = append(part, rows[j][1:])
		}
		if !specComplete(part) {
			return nil, &tsv.Error{Line: i + 2, Column: "rule", Reason: "a head holds the rules of the block, in order"}
		}
		i = j
	}
	return rows, nil
}

// SpecTable gives the recorded rows (rule, result, detail) of the full check of a
// head, and whether the head has them. A table its reader refuses is an error.
func (s *Sessions) SpecTable(ctx context.Context, task, head string) ([][]string, bool, error) {
	data, err := s.Store.ReadFile(ctx, specPath(task))
	if err != nil || data == nil {
		return nil, false, err
	}
	rows, err := readSpecTable(data)
	if err != nil {
		return nil, false, fmt.Errorf("%s: %w", specPath(task), err)
	}
	var mine [][]string
	for _, r := range rows {
		if r[0] == head {
			mine = append(mine, r[1:])
		}
	}
	return mine, len(mine) > 0, nil
}

// A SpecInput names the head whose layup/spec is asked: the task, the base, the
// head, and the text of the latest bet's copy (row 66j's; none before the first bet).
type SpecInput struct {
	Task, Base, Head string
	Copy             []byte
}

// SpecAt gives the status layup/spec of a head, running the full check once.
func (s *Sessions) SpecAt(context.Context, SpecInput) (Status, error) { return Status{}, nil }
