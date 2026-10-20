package run

// This file holds the inventory version and the rendered commit of the first bet (docs/spec/loop.md, The step bet; task T-2vva).

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/spec"
)

func tempDir() (string, error) { return os.MkdirTemp("", "layup-rendered-") }
func removeAll(p string)       { os.RemoveAll(p) }

// writeFiles writes the files, by their paths below dir, as regular files.
func writeFiles(dir string, files map[string][]byte) error {
	for p, data := range files {
		path := filepath.Join(dir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func findings(fs []spec.Finding) string {
	var parts []string
	for _, f := range fs {
		parts = append(parts, f.Rule+": "+f.Detail)
	}
	return strings.Join(parts, "; ")
}

// nextVersion gives the files of the next version of the confirmed inventory, spec/inventory-<n>.tsv and spec/inventory.tsv, and the rows
// of the version: the rows of the latest version, with no file, when that version already names this copy. The sha256 of each version
// the index names is checked against its bytes first.
func nextVersion(index []byte, versions map[int][]byte, facts []spec.InventoryFact, answer spec.Answer, copyID int) (map[string][]byte, [][]string, error) {
	var idx [][]string
	if index != nil {
		var err error
		if idx, err = spec.ReadInventoryIndex(index); err != nil {
			return nil, nil, fmt.Errorf("spec/inventory.tsv: %w", err)
		}
	}
	for _, row := range idx {
		n, _ := strconv.Atoi(row[0])
		if got := digest(versions[n]); versions[n] == nil || got != row[1] {
			return nil, nil, fmt.Errorf("spec/inventory-%d.tsv has the sha256 %s; the index has %s", n, got, row[1])
		}
	}
	var prev [][]string
	if n := len(idx); n > 0 {
		var err error
		if prev, err = spec.ReadInventory(versions[n]); err != nil {
			return nil, nil, fmt.Errorf("spec/inventory-%d.tsv: %w", n, err)
		}
		if len(prev) > 0 && prev[0][3] == strconv.Itoa(copyID) {
			return nil, prev, nil
		}
	}
	version, next, fs := spec.NextInventory(prev, idx, facts, answer, copyID)
	if len(fs) > 0 {
		return nil, nil, fmt.Errorf("the inventory: %s", findings(fs))
	}
	rows, err := spec.ReadInventory(version)
	if err != nil {
		return nil, nil, err
	}
	return map[string][]byte{fmt.Sprintf("spec/inventory-%d.tsv", len(idx)+1): version, "spec/inventory.tsv": next}, rows, nil
}

// renderedFiles gives the files of the rendered commit and their records: the one PRD at the verified head, rendered from the answer, and
// the numbered facts record under id, rendered from the rows of the version. Two PRDs, no PRD or a numbered facts record that exists is an error.
func renderedFiles(prds, facts map[string][]byte, id string, answer spec.Answer, version [][]string, texts map[string]string) (map[string][]byte, []rules.Record, error) {
	if len(prds) != 1 {
		return nil, nil, fmt.Errorf("the rendered commit needs exactly one docs/prd/PRD-*.md at the verified head; it holds %d", len(prds))
	}
	if len(facts) > 0 {
		return nil, nil, fmt.Errorf("the verified head already holds a numbered facts record: %s", slices.Sorted(maps.Keys(facts))[0])
	}
	var name string
	var pre []byte
	for name, pre = range prds {
	}
	recs := []rules.Record{{Path: "docs/prd/" + name, Kind: rules.PRD, Answer: answer}, {Path: "docs/facts/" + id + "-numbered-facts.md", Kind: rules.Facts, Version: version, Texts: texts}}
	out := map[string][]byte{}
	for _, rec := range recs {
		data, fs := rules.Render(rec, pre)
		if len(fs) > 0 {
			return nil, nil, fmt.Errorf("the rendering of %s: %s", rec.Path, findings(fs))
		}
		out[rec.Path] = data
	}
	return out, recs, nil
}
