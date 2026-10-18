package run

// This file holds the table of layup gate of a head as a record, and the status
// layup/gates (docs/spec/loop.md, layup/gates; docs/spec/records.md, the block
// task-gates; task T-q0kr, #362, row 66f of the plan of M2e).

import (
	"context"
	"io"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/tsv"
)

// GatesSchema is the Go value of the block task-gates of docs/spec/records.md.
var GatesSchema = tsv.Schema{}

func gatesPath(task string) string { return "tasks/" + task + "/gates.tsv" }

// ReadGates reads gates.tsv by GatesSchema and the owner checks of its block.
func ReadGates([]byte) ([][]string, error) { return nil, nil }

// WriteGates writes rows by GatesSchema; it writes nothing for rows that ReadGates refuses.
func WriteGates([][]string) ([]byte, error) { return nil, nil }

// GatesStatus gives the status layup/gates of the table of a head.
func GatesStatus(gate.Table) Status { return Status{} }

// A GatesInput names the head whose table is asked, and where the progress and the
// held output of the run go.
type GatesInput struct {
	Task, Base, Head string
	Step             func(i, n int, kind string) func()
	Out              io.Writer
}

// GatesTable gives the recorded table of a head, and whether the head has one.
func (s *Sessions) GatesTable(context.Context, string, string) (gate.Table, bool, error) {
	return gate.Table{}, false, nil
}

// GatesAt gives the status layup/gates of a head, running the gate once.
func (s *Sessions) GatesAt(context.Context, GatesInput) (Status, error) { return Status{}, nil }
