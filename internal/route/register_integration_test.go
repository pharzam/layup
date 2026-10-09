//go:build integration

package route

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The Go schemas of the blocks harness-register, models and routing-register
// equal their blocks of docs/spec/records.md (the owner's comparison;
// docs/spec/README.md, The schema block).
func TestTheSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"harness-register": HarnessRegisterSchema, "models": ModelsSchema,
		"routing-register": RoutingRegisterSchema} {
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

// CheckCredential reads a real file: mode 0600 of the user of the run passes,
// mode 0644 is refused, naming the file.
func TestCheckCredentialOnARealFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "claude.key")
	if err := os.WriteFile(path, []byte("sk-test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := CheckCredential(path); err != nil {
		t.Fatalf("a file of mode 0600: %v", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CheckCredential(path); err == nil {
		t.Error("a file of mode 0644: no error")
	}
}
