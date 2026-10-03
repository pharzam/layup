//go:build integration

package records

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go schemas of telemetry.tsv and prices.tsv equal their blocks of
// docs/spec/records.md (the demo of #94).
func TestTheSchemasEqualTheirBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"telemetry": TelemetrySchema, "prices": PricesSchema} {
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
