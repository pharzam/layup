//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The package rules of docs/spec/packages.md hold on the real module, and the
// same checker finds the seeded breach of rule 5 in a fixture module that
// imports net/http (NFR-005, NFR-007).
func TestPackageRules(t *testing.T) {
	text, err := os.ReadFile(filepath.Join("..", "..", "docs", "spec", "packages.md"))
	if err != nil {
		t.Fatal(err)
	}
	rows, err := readTable(string(text))
	if err != nil || len(rows) == 0 {
		t.Fatalf("the table of phase 1: %d rows, %v", len(rows), err)
	}
	real := load(t, filepath.Join("..", ".."))
	if len(real.sources["internal/git"]) == 0 || len(real.sources["cmd/layup"]) == 0 {
		t.Fatal("no source of internal/git or cmd/layup was read, so the scan checks nothing")
	}
	if f := checkRules(rows, real); len(f) != 0 {
		t.Errorf("the module breaks the package rules:\n%s", strings.Join(f, "\n"))
	}
	f := checkRules(rows, load(t, filepath.Join("testdata", "netimport")))
	if !strings.Contains(strings.Join(f, "\n")+"\n", "rule 5: cmd/layup depends on net/http\n") {
		t.Errorf("the fixture that imports net/http gives\n%s\nwant the finding rule 5: cmd/layup depends on net/http", strings.Join(f, "\n"))
	}
	for _, s := range f {
		if !strings.HasPrefix(s, "rule 5: ") {
			t.Errorf("the fixture breaks another rule too: %s", s)
		}
	}
}

// load reads the module at root: go.mod by go mod edit -json, its packages
// and their dependencies by go list -deps -json, and the non-test Go files of
// its packages, also the ones that a build constraint leaves out.
func load(t *testing.T, root string) module {
	t.Helper()
	var mod struct {
		Module  struct{ Path string }
		Require []struct{ Path, Version string }
	}
	if err := json.Unmarshal(goCmd(t, root, "mod", "edit", "-json"), &mod); err != nil {
		t.Fatal(err)
	}
	m := module{path: mod.Module.Path, sources: map[string]map[string]string{}}
	for _, r := range mod.Require {
		m.requires = append(m.requires, r.Path+" "+r.Version)
	}
	for dec := json.NewDecoder(bytes.NewReader(goCmd(t, root, "list", "-deps", "-json", "./..."))); dec.More(); {
		var p goPackage
		if err := dec.Decode(&p); err != nil {
			t.Fatal(err)
		}
		m.packages = append(m.packages, p)
		rel, ok := inModule(m.path, p.ImportPath)
		if !ok {
			continue
		}
		m.sources[rel] = map[string]string{}
		for _, name := range append(append(p.GoFiles, p.CgoFiles...), p.IgnoredGoFiles...) {
			if !strings.HasSuffix(name, "_test.go") {
				src, err := os.ReadFile(filepath.Join(p.Dir, name))
				if err != nil {
					t.Fatal(err)
				}
				m.sources[rel][name] = string(src)
			}
		}
	}
	return m
}

// goCmd runs the go command at root, with no setting of the host that changes
// what it reads, and no download: GOFLAGS, the go env file, a workspace, a
// toolchain switch and the module proxy are off.
func goCmd(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("go", args...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "GOFLAGS=", "GOENV=off", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go %s in %s: %v\n%s", strings.Join(args, " "), root, err, stderr.String())
	}
	return out
}
