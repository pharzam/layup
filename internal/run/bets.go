package run

// The record bets.tsv (docs/spec/records.md, the block bets; task T-2vva, #366, row 66j): each bet with its brief, its answer and its verdict.

import (
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// BetsSchema is the Go value of the block bets; its owner is this package.
var BetsSchema = tsv.Schema{Name: "bets", Location: "records:bets.tsv"}

// The columns of a row of bets.tsv.
const (
	betMilestone = iota
	betKind
	betBrief
	betBriefSHA
	betUSD
	betHours
	betAnswerID
	betOptions
	betVerdict
	betApprover
	betDecided
)

// ReadBets reads bets.tsv by BetsSchema and then makes the checks of its owner (docs/spec/records.md, bets).
func ReadBets(data []byte) ([][]string, error) { return nil, nil }

// betRow gives the row of bets.tsv of an answered brief: the posted copy of the brief, the appetite of the kept brief (empty for a
// no-bet), the copy row of the answer, the answer and the approver's row of approvers.tsv.
func betRow(milestone string, posted Posted, appetite spec.Cap, answer spec.Answer, copyRow, approver []string) []string {
	return nil
}
