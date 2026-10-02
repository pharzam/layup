//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/catalog"
)

// The fixtures of each entry of the binary's catalog, through the built
// binary (D10 of #91, with note 12 of its plan review): on a target rendered
// from the entry, layup gate gives exit 0 on the clean commit, and exit 1 with
// the row of the kind fail on the clean commit with the kind's fixture
// applied. The integration test of internal/catalog checks each row; this test
// checks the exit codes and the row of the fixture's kind.
func TestGateOnEachEntryOfTheCatalog(t *testing.T) {
	stacks, err := catalog.Stacks(catalog.Embedded())
	if err != nil || len(stacks) == 0 {
		t.Fatalf("the binary embeds no entry: %q, %v", stacks, err)
	}
	for _, s := range stacks {
		e, err := catalog.Read(catalog.Embedded(), s)
		if err != nil {
			t.Fatal(err)
		}
		files, err := e.Files("example.com/target")
		if err != nil {
			t.Fatal(err)
		}
		m, err := e.Manifest()
		if err != nil {
			t.Fatal(err)
		}
		all := map[string]string{"docs/gates.tsv": string(m)}
		for _, f := range files {
			all[f.Path] = string(f.Data)
		}
		dir := t.TempDir()
		fixtureGit(t, dir, "init", "-q", "-b", "main")
		clean := commitAll(t, dir, all)
		if r := layupWith(t, goEnv(t), "gate", dir, "--base", clean, "--head", clean); r.code != 0 {
			t.Errorf("%s: the clean commit: exit %d; want 0\n%s%s", s, r.code, r.stdout, r.stderr)
		}
		for _, k := range e.Kinds() {
			if k.State != "active" {
				continue
			}
			patch, err := e.Fixture(k.Name)
			if err != nil {
				t.Fatal(err)
			}
			p := filepath.Join(t.TempDir(), k.Name+".patch")
			if err := os.WriteFile(p, patch, 0o644); err != nil {
				t.Fatal(err)
			}
			fixtureGit(t, dir, "switch", "-q", "-c", "fixture-"+k.Name, clean)
			fixtureGit(t, dir, "apply", p)
			head := commitAll(t, dir, nil)
			r := layupWith(t, goEnv(t), "gate", dir, "--base", clean, "--head", head)
			if r.code != 1 || !strings.Contains(r.stdout, "\t"+k.Name+"\tactive\tfail\t") {
				t.Errorf("%s: the fixture of %s: exit %d; want 1 and the row %s fail\n%s%s", s, k.Name, r.code, k.Name, r.stdout, r.stderr)
			}
		}
	}
}
