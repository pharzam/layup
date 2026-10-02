//go:build integration

package catalog

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// testCatalog is the test entry, embedded with the pattern form of an entry
// of the binary: the prefix all: keeps .github/ (docs/spec/setup.md, The stack
// catalog).
//
//go:embed all:testdata/test
var testCatalog embed.FS

// The embedded test entry is read by its rules, with its go.mod and its
// .github/ files (the demo of #81). The root of the fs.FS that Read and Stacks
// take is the catalog root, which holds one directory per stack: the binary's
// //go:embed all:<stack> gives that root directly, and this test takes the
// directory testdata/ of its own pattern as the root.
func TestTheEmbeddedTestEntry(t *testing.T) {
	root, err := fs.Sub(testCatalog, "testdata")
	if err != nil {
		t.Fatal(err)
	}
	if s, err := Stacks(root); err != nil || !slices.Equal(s, []string{"test"}) {
		t.Fatalf("Stacks = %q, %v; want test", s, err)
	}
	e, err := Read(root, "test")
	if err != nil {
		t.Fatal(err)
	}
	files, err := e.Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, f := range files {
		paths = append(paths, f.Path)
		if f.Path == "go.mod" && string(f.Data) != "module example.com/target\n\ngo 1.26\n" {
			t.Errorf("go.mod: %q", f.Data)
		}
	}
	if !slices.Equal(paths, []string{".github/workflows/gates.yml", "go.mod"}) {
		t.Fatalf("the files %q; want .github/workflows/gates.yml and go.mod", paths)
	}
	if p, err := e.Fixture("static"); err != nil || len(p) == 0 {
		t.Fatalf("the fixture of static: %q, %v", p, err)
	}
	for _, ref := range []string{"test/kinds.tsv", "test/files/go.mod.tmpl", "test/files/.github/workflows/gates.yml.tmpl", "test/fixtures/static.patch"} {
		if !e.Has(ref) {
			t.Errorf("Has(%q) = false", ref)
		}
	}
}

// The Go schemas of kinds.tsv and of the manifest equal their blocks in
// docs/spec/, and the manifest of the test entry reads back by the block of
// gate-manifest.
func TestTheSchemaBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"catalog-kinds": KindsSchema, "gate-manifest": ManifestSchema} {
		block, ok := blocks[name]
		if !ok {
			t.Fatalf("docs/spec/ has no block %s", name)
		}
		if err := tsv.Compare(block, s); err != nil {
			t.Error(err)
		}
	}
	root, _ := fs.Sub(testCatalog, "testdata")
	e, err := Read(root, "test")
	if err != nil {
		t.Fatal(err)
	}
	m, err := e.Manifest()
	if err != nil {
		t.Fatal(err)
	}
	rows, err := tsv.Read(m, blocks["gate-manifest"])
	if err != nil || len(rows) != 2 {
		t.Fatalf("the manifest by the block gate-manifest: %d rows, %v", len(rows), err)
	}
}
