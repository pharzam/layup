//go:build integration

package run

import (
	"os"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go schemas of the blocks setup-tree and forge-check equal their blocks of
// docs/spec/records.md (the demo of #270; the block markers is compared in
// internal/records).
func TestTheSchemasOfTheScaffoldEqualTheirBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"setup-tree": SetupTreeSchema, "forge-check": ForgeCheckSchema} {
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
