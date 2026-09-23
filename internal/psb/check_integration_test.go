//go:build integration

package psb

import (
	"path/filepath"
	"testing"
)

// The real PSB (fact F-0001) must yield the golden batch, which holds the G1
// row: the gap that cost a question in the manual setup.
func TestGoldenRealPSB(t *testing.T) {
	golden(t, filepath.Join("..", "..", "docs", "facts", "problem-statement-brief.md"), filepath.Join("testdata", "psb.tsv"))
}
