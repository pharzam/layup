package run

// The record bets.tsv (docs/spec/records.md, the block bets; task T-2vva, #366, row 66j): each bet with its brief, its answer and its verdict.

import (
	"fmt"
	"strings"

	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// BetsSchema is the Go value of the block bets; its owner is this package.
var BetsSchema = tsv.Schema{Name: "bets", Location: "records:bets.tsv", Columns: []tsv.Column{
	{Name: "milestone", Type: "text", Key: true},
	{Name: "kind", Type: "enum(first|later)"},
	{Name: "brief", Type: "int"},
	{Name: "brief_sha256", Type: "sha256"},
	{Name: "appetite_usd", Type: "decimal"},
	{Name: "appetite_hours", Type: "decimal"},
	{Name: "answer", Type: "int"},
	{Name: "options", Type: "list(text)"},
	{Name: "verdict", Type: "enum(bet|no-bet)"},
	{Name: "approver", Type: "int"},
	{Name: "decided", Type: "time"},
}}

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

// ReadBets reads bets.tsv by BetsSchema and then makes the checks of its owner (docs/spec/records.md, bets): brief_sha256 and verdict
// are never —, the appetite is — exactly for a no-bet, the kind first is there once, and a later bet comes after it.
func ReadBets(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, BetsSchema)
	if err != nil {
		return nil, err
	}
	first := false
	for i, r := range rows {
		bad := func(column, reason string) error { return &tsv.Error{Line: i + 2, Column: column, Reason: reason} }
		switch {
		case r[betBriefSHA] == "":
			return nil, bad("brief_sha256", "this column never holds the empty value")
		case r[betVerdict] == "":
			return nil, bad("verdict", "this column never holds the empty value")
		case (r[betVerdict] == "no-bet") != (r[betUSD] == ""):
			return nil, bad("appetite_usd", "the appetite is the empty value exactly for a no-bet")
		case (r[betVerdict] == "no-bet") != (r[betHours] == ""):
			return nil, bad("appetite_hours", "the appetite is the empty value exactly for a no-bet")
		case r[betKind] == "first" && first:
			return nil, bad("kind", "a second row of the kind first: the first bet is made once")
		case r[betKind] == "later" && !first:
			return nil, bad("kind", "a later bet comes after the first")
		}
		first = first || r[betKind] == "first"
	}
	return rows, nil
}

// decimalText writes a cap of the brief as a decimal of the record: 5 is 5.0.
func decimalText(s string) string {
	if s == "" || strings.Contains(s, ".") {
		return s
	}
	return s + ".0"
}

// betRow gives the row of bets.tsv of an answered brief: the posted copy of the brief, the appetite of the kept brief (empty for a
// no-bet), the copy row of the answer, the answer and the approver's row of approvers.tsv.
func betRow(milestone string, posted Posted, appetite spec.Cap, answer spec.Answer, copyRow, approver []string) []string {
	usd, hours := decimalText(appetite.USD), decimalText(appetite.Hours)
	if answer.Verdict == "no-bet" {
		usd, hours = "", ""
	}
	return []string{milestone, "first", fmt.Sprint(posted.ID), posted.SHA, usd, hours, copyRow[copyComment], strings.Join(optionsColumn(answer.Options), " "), answer.Verdict, approver[0], copyRow[copyCopied]}
}
