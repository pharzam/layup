//go:build integration

package spec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go value of the block spec-check, and the second Go value of the block
// spans, equal their blocks of docs/spec/records.md (task T-4ywt, #213); the
// blocks inventory and inventory-index, which this package owns, too (task
// T-3q0b, #350).
func TestTheSchemasEqualTheirBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"spec-check": CheckSchema, "spans": SpansSchema,
		"inventory": InventorySchema, "inventory-index": InventoryIndexSchema} {
		block, ok := blocks[name]
		if !ok {
			t.Errorf("docs/spec/ has no block %s", name)
			continue
		}
		if err := tsv.Compare(block, s); err != nil {
			t.Error(err)
		}
	}
}

// The Go value of the block options equals its block of docs/spec/records.md
// (task T-0pt8, #349).
func TestTheOptionsSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["options"]
	if !ok {
		t.Fatal("docs/spec/ has no block options")
	}
	if err := tsv.Compare(block, OptionsSchema); err != nil {
		t.Error(err)
	}
}
