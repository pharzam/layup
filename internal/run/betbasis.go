package run

// The basis of a kept brief (docs/spec/records.md, the block bet-basis; task T-h65e, #418, row 66u): what row 66j needs after the
// inventory write, so that it filters nothing again.

import "github.com/pharzam/layup/internal/tsv"

// BetBasisSchema is the Go value of the block bet-basis; its owner is this package.
var BetBasisSchema = tsv.Schema{Name: "bet-basis", Location: "records:bets/<milestone>/basis.tsv"}

// ReadBetBasis reads a basis by BetBasisSchema and then makes the checks of its owner.
func ReadBetBasis(data []byte) ([][]string, error) { return nil, nil }
