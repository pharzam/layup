//go:build integration

package tsv

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// built lists each block of docs/spec/ whose owner package compares it with
// its Go schema in the owner's own test (tsv.Compare). The owner moves the name
// here from notYetBuilt in the same change as that test.
var built = []string{"catalog-kinds", "gate-manifest", "gate-result", "open-gaps", "psb-gaps", "setup-answers", "setup-record", "setup-steps", "setup-stop", "setup-verify"}

// notYetBuilt lists each block that no owner compares with a Go schema yet.
// It only becomes shorter: a rule for the reviewer of each owner task; a test
// cannot read the list of its base (docs/spec/README.md, The schema block).
var notYetBuilt = []string{
	"prices", "rule-paths", "stalls", "telemetry",
}

// Each block of docs/spec/ has the form of the README, and is in exactly one
// of the two lists; each listed name is a block. A block that a later section
// adds fails this test until it is listed.
func TestEverySchemaBlockIsBuiltOrNotYetBuilt(t *testing.T) {
	blocks, err := ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	listed := map[string]int{}
	for _, name := range append(slices.Clone(built), notYetBuilt...) {
		listed[name]++
		if _, ok := blocks[name]; !ok {
			t.Errorf("%s is listed, but no block of docs/spec/ has that name", name)
		}
	}
	for _, name := range slices.Sorted(maps.Keys(blocks)) {
		if listed[name] != 1 {
			t.Errorf("the block %s is listed %d times in built and notYetBuilt; want 1", name, listed[name])
		}
	}
}
