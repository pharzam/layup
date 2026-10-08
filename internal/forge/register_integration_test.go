//go:build integration

package forge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go schema of the block forge-register equals its block of docs/spec/records.md
// (the owner's comparison; docs/spec/README.md, The schema block).
func TestTheSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["forge-register"]
	if !ok {
		t.Fatal("docs/spec/ has no block forge-register")
	}
	if err := tsv.Compare(block, ForgeRegisterSchema); err != nil {
		t.Error(err)
	}
}
