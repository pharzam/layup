package run

// The Retrospective of docs/spec/loop.md (The Retrospective), task T-v2sv,
// #359, row 66c of the plan of M2e; the record settles #372 under R10.
// architecture.md §8 puts the step between the Accept and the next Bet; in
// phases 2 and 3 it is a step with no session: code writes its event — a row
// of retrospectives.tsv, the block retrospectives — and the loop goes to the
// next Bet. No caller exists in M2e: loop.md gives the step loop only the
// Shape merges and the register's tasks, and Accept is M2g's. The helper is
// unwired, as registerTask is, and "the loop goes to the next Bet" is the
// value the call returns, not a transition made here.

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/pharzam/layup/internal/tsv"
)

// RetrospectivesSchema is the Go value of the block retrospectives of
// docs/spec/records.md; its owner is this package, and its test compares it
// with the block.
var RetrospectivesSchema = tsv.Schema{Name: "retrospectives", Location: "records:retrospectives.tsv", Columns: []tsv.Column{
	{Name: "milestone", Type: "text", Key: true},
	{Name: "time", Type: "time"},
}}

// retrospectivesPath is the record of the Retrospective's events, one row per
// milestone.
const retrospectivesPath = "retrospectives.tsv"

// nextBet is the step name the Retrospective gives: after it the loop goes to
// the next Bet. It is a value returned to the caller, not a transition made
// here.
const nextBet = "bet"

// eventStore is the records branch as the Retrospective step sees it: the
// lease, for Fenced; ReadFile gives a file at the run's records base, nil for
// none; Commit makes one records commit of files and pushes it, and gives
// ErrRefused when git refuses the push. The seam holds no call that starts a
// session or asks a budget.
type eventStore interface {
	Store
	ReadFile(ctx context.Context, path string) ([]byte, error)
	Commit(ctx context.Context, files map[string][]byte, message string) error
}

// ReadRetrospectives reads the table by RetrospectivesSchema and then makes
// the checks of its owner (CheckRetrospectives).
func ReadRetrospectives(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, RetrospectivesSchema)
	if err != nil {
		return nil, err
	}
	if err := CheckRetrospectives(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckRetrospectives makes the checks that docs/spec/records.md gives the
// owner of retrospectives, after the types: time is never — (a non-key column
// would hold it, so the owner checks it). The refusal names the line and the
// column.
func CheckRetrospectives(rows [][]string) error {
	for i, r := range rows {
		if r[1] == "" {
			return &tsv.Error{Line: i + 2, Column: "time", Reason: "the event's time is never —"}
		}
	}
	return nil
}

// retrospectiveFiles gives the files of the records commit that writes the
// event of milestone: retrospectives.tsv with its row appended — or no file
// when the row exists (a write is read before it is repeated, so a restart is
// done before). A table that ReadRetrospectives refuses, and a milestone that
// the key rule refuses, stop the write.
func retrospectiveFiles(table []byte, milestone string, at time.Time) (map[string][]byte, bool, error) {
	var rows [][]string
	if table != nil {
		var err error
		if rows, err = ReadRetrospectives(table); err != nil {
			return nil, false, fmt.Errorf("%s: %w", retrospectivesPath, err)
		}
	}
	for _, r := range rows {
		if r[0] == milestone {
			return nil, false, nil
		}
	}
	rows = append(rows, []string{milestone, stamp(at)})
	var b bytes.Buffer
	if err := tsv.Write(&b, RetrospectivesSchema, rows); err != nil {
		return nil, false, fmt.Errorf("%s: %w", retrospectivesPath, err)
	}
	if _, err := ReadRetrospectives(b.Bytes()); err != nil {
		return nil, false, fmt.Errorf("%s: %w", retrospectivesPath, err)
	}
	return map[string][]byte{retrospectivesPath: b.Bytes()}, true, nil
}

// retrospectiveStep is the Retrospective step on the records seam: the event
// of milestone in one fenced records commit, or none when the row exists, and
// the answer nextBet. It starts no session: the seam holds no call that does.
func retrospectiveStep(ctx context.Context, s eventStore, runID, milestone string, at time.Time) (string, error) {
	data, err := s.ReadFile(ctx, retrospectivesPath)
	if err != nil {
		return "", err
	}
	files, write, err := retrospectiveFiles(data, milestone, at)
	if err != nil {
		return "", err
	}
	if !write {
		return nextBet, nil
	}
	if err := Fenced(ctx, s, runID, func(ctx context.Context) error {
		return s.Commit(ctx, files, "retrospectives: the Retrospective of "+milestone)
	}); err != nil {
		return "", err
	}
	return nextBet, nil
}

// retrospective is the step on the run's own state: its store, its run ID and
// its clock. Nothing in M2e calls it.
func (r *state) retrospective(ctx context.Context, milestone string) (string, error) {
	return retrospectiveStep(ctx, r.store, r.cfg.RunID, milestone, r.cfg.Clock.Now().UTC().Truncate(time.Second))
}
