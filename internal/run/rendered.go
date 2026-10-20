package run

// This file holds the inventory version and the rendered commit of the first bet (docs/spec/loop.md, The step bet; task T-2vva).

import (
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/spec"
)

// nextVersion gives the files of the next version of the confirmed inventory, spec/inventory-<n>.tsv and spec/inventory.tsv, and the rows
// of the version: the rows of the latest version, with no file, when that version already names this copy. The sha256 of each version
// the index names is checked against its bytes first.
func nextVersion(index []byte, versions map[int][]byte, facts []spec.InventoryFact, answer spec.Answer, copyID int) (map[string][]byte, [][]string, error) {
	return nil, nil, nil
}

// renderedFiles gives the files of the rendered commit and their records: the one PRD at the verified head, rendered from the answer, and
// the numbered facts record under id, rendered from the rows of the version. Two PRDs, no PRD or a numbered facts record that exists is an error.
func renderedFiles(prds, facts map[string][]byte, id string, answer spec.Answer, version [][]string, texts map[string]string) (map[string][]byte, []rules.Record, error) {
	return nil, nil, nil
}
