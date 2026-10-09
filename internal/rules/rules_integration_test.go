//go:build integration

package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The second Go value of the block rule-paths equals its block of
// docs/spec/setup.md, as the owner's value does (internal/setup).
func TestTheSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["rule-paths"]
	if !ok {
		t.Fatal("docs/spec/ has no block rule-paths")
	}
	if err := tsv.Compare(block, RulePathsSchema); err != nil {
		t.Error(err)
	}
}
