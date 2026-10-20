package run

// The basis of a kept brief (docs/spec/records.md, the block bet-basis; task T-h65e, #418, row 66u): what row 66j needs after the
// inventory write, so that it filters nothing again.

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// BetBasisSchema is the Go value of the block bet-basis; its owner is this package.
var BetBasisSchema = tsv.Schema{Name: "bet-basis", Location: "records:bets/<milestone>/basis.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "enum(requirement|survivor|dropped|fact|open|head)", Key: true},
	{Name: "id", Type: "text", Key: true},
	{Name: "value", Type: "text"},
}}

var (
	basisRequirement = regexp.MustCompile(`^(REQ|NFR)-([0-9]{3,})$`)
	basisFact        = regexp.MustCompile(`^([0-9]+|Q-[0-9]{3,})$`)
	basisQuestion    = regexp.MustCompile(`^Q-[0-9]{3,}$`)
	basisTask        = regexp.MustCompile(`^T-[a-z0-9]{4}$`)
	basisSHA         = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// ReadBetBasis reads a basis by BetBasisSchema and then makes the checks of its owner (docs/spec/records.md, bet-basis): the rule of
// each kind, a survivor for each decision, and the kinds that a basis needs.
func ReadBetBasis(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, BetBasisSchema)
	if err != nil {
		return nil, err
	}
	count, first, survives := map[string]int{}, map[string]int{}, map[string]bool{}
	for i, r := range rows {
		kind, id, value, line := r[0], r[1], r[2], i+2
		bad := func(column, reason string) error { return &tsv.Error{Line: line, Column: column, Reason: reason} }
		count[kind]++
		decision, option, split := strings.Cut(id, "=")
		switch kind {
		case "requirement":
			m := basisRequirement.FindStringSubmatch(id)
			if m == nil {
				return nil, bad("id", "a requirement is REQ-NNN or NFR-NNN")
			} else if value != "R-"+m[2] {
				return nil, bad("value", fmt.Sprintf("the row of the draft for %s is R-%s", id, m[2]))
			}
		case "survivor", "dropped":
			if !split || decision == "" || option == "" || strings.ContainsAny(id, " ") || strings.Contains(option, "=") {
				return nil, bad("id", "a survivor or a dropped option is <decision>=<option>")
			} else if kind == "survivor" && value != "" {
				return nil, bad("value", "a survivor holds —")
			} else if kind == "dropped" && value == "" {
				return nil, bad("value", "a dropped option holds its reasons, not —")
			}
			if _, seen := first[decision]; !seen {
				first[decision] = line
			}
			survives[decision] = survives[decision] || kind == "survivor"
		case "fact":
			if !basisFact.MatchString(id) {
				return nil, bad("id", "a fact is a number or a Q-NNN")
			} else if value != "need" && value != "constraint" && value != "measure" && value != "context" {
				return nil, bad("value", "a fact's class is need, constraint, measure or context")
			}
		case "open":
			if !basisQuestion.MatchString(id) {
				return nil, bad("id", "an open question is a Q-NNN")
			} else if value != "" {
				return nil, bad("value", "an open question holds —")
			}
		case "head":
			if !basisTask.MatchString(id) {
				return nil, bad("id", "a head is of a task T-xxxx")
			} else if !basisSHA.MatchString(value) {
				return nil, bad("value", "a head is 40 lowercase hexadecimal characters")
			}
		}
	}
	for _, d := range slices.Sorted(maps.Keys(first)) {
		if !survives[d] {
			return nil, &tsv.Error{Line: first[d], Column: "id", Reason: fmt.Sprintf("the decision %s has no survivor", d)}
		}
	}
	switch {
	case count["requirement"] == 0:
		return nil, fmt.Errorf("the basis has no requirement row")
	case count["fact"] == 0:
		return nil, fmt.Errorf("the basis has no fact row")
	case count["head"] < 2:
		return nil, fmt.Errorf("the basis has %d head rows, want two or more", count["head"])
	}
	return rows, nil
}
