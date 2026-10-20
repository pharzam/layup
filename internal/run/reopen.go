package run

// This file holds the specification task that a changed class reopens (docs/spec/loop.md, The specification task that a changed class
// reopens; task T-4naq, #428, row 66ab of the plan of M2e).

import (
	"context"

	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// reopeningsPath is the record of the facts that the bet re-classed.
const reopeningsPath = "spec/reopenings.tsv"

// ReopeningsSchema is the Go value of the block reopenings; its owner is this package.
var ReopeningsSchema = tsv.Schema{Name: "reopenings", Location: "records:spec/reopenings.tsv"}

// The columns of a row of spec/reopenings.tsv.
const (
	reoN = iota
	reoFact
	reoClass
	reoToken
	reoRequirement
	reoTask
	reoCopy
)

// reopenIssueBlock is the body of the issue of a reopened task: the block reopen-issue of loop.md.
const reopenIssueBlock = ""

// reopenedInstruction is the instruction of the step specification of a reopened task: the block instruction-specification-reopened.
const reopenedInstruction = ""

// ReadReopenings reads spec/reopenings.tsv by ReopeningsSchema and then makes the checks of its owner.
func ReadReopenings(data []byte) ([][]string, error) { return nil, nil }

// reopenedTasks gives the tasks that the rows name.
func reopenedTasks(rows [][]string) map[string]bool { return map[string]bool{} }

// reclassed selects the rows of a version of the inventory whose class differs from the class of the fact at Intake.
func reclassed(version [][]string, intake []spec.InventoryFact) [][]string { return nil }

// tokenOf is the citation of a fact in the numbered facts record.
func tokenOf(fact, record string) string { return "" }

// nextNumber is one more than the highest number of the REQ- and NFR- IDs given, 1 for none.
func nextNumber(ids ...[]string) int { return 0 }

// allocate makes the rows of spec/reopenings.tsv of the selected rows of version n.
func allocate(sel [][]string, n int, record string, next int, tasks []string) [][]string { return nil }

// bindReopenings compares the rows of the record with the selection of the version.
func bindReopenings(rows, sel [][]string, record string) error { return nil }

// reopeningOf gives the row of the task, nil for none.
func reopeningOf(rows [][]string, task string) []string { return nil }

// reopeningFacts gives the facts that the full check exempts at the head of the task.
func reopeningFacts(version [][]string, intake []spec.InventoryFact, rows [][]string, merged []string, task string) []string {
	return nil
}

// openReopenings opens the specification tasks that the classes of the first bet reopen, after the event merged of the specification task.
func (r *state) openReopenings(ctx context.Context, specTask string) ([]Task, error) { return nil, nil }

// taskIDs draws n task IDs, none of a row of sessions.tsv, tasks.tsv or spec/reopenings.tsv, and no two alike.
func (s *Sessions) taskIDs(ctx context.Context, n int) ([]string, error) { return nil, nil }
