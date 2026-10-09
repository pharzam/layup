//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The package rules of docs/spec/packages.md hold on the real module, with the
// tables of phase 1, of M2a and of M2b, the line of rule 5 and the line of the
// engine checks. The same checker finds the seeded breaches of two fixture
// modules: in netimport, an import of net/http outside the adapter, and one of
// net/smtp behind a build constraint, which also breaks the rule of the engine
// checks in internal/psb; in enginedep, a package of the engine checks that
// depends on internal/session (NFR-005, NFR-007).
func TestPackageRules(t *testing.T) {
	text, err := os.ReadFile(filepath.Join("..", "..", "docs", "spec", "packages.md"))
	if err != nil {
		t.Fatal(err)
	}
	rows, err := readTable(string(text))
	if err != nil || len(rows) == 0 {
		t.Fatalf("the tables of phase 1 and of M2a: %d rows, %v", len(rows), err)
	}
	adapter, err := readAdapter(string(text))
	if err != nil {
		t.Fatal(err)
	}
	engine, err := readEngine(string(text))
	if err != nil {
		t.Fatal(err)
	}
	real := load(t, filepath.Join("..", ".."))
	if len(real.sources["internal/git"]) == 0 || len(real.sources["cmd/layup"]) == 0 {
		t.Fatal("no source of internal/git or cmd/layup was read, so the scan checks nothing")
	}
	for _, p := range engine {
		if _, ok := real.sources[p]; !ok {
			t.Errorf("the line of the engine checks names %s, which is no package of the module", p)
		}
	}
	if f := checkRules(rows, adapter, engine, real); len(f) != 0 {
		t.Errorf("the module breaks the package rules:\n%s", strings.Join(f, "\n"))
	}
	for fixture, want := range map[string][]string{
		"netimport": {"engine checks: internal/psb depends on net",
			"rule 5: cmd/layup imports net/http; only " + adapter + " imports them", "rule 5: internal/psb depends on net"},
		"enginedep": {"May import: internal/gate imports internal/session, which its row does not allow",
			"engine checks: internal/gate depends on internal/session"},
	} {
		f := checkRules(rows, adapter, engine, load(t, filepath.Join("testdata", fixture)))
		for _, w := range want {
			if !slices.Contains(f, w) {
				t.Errorf("%s gives\n%s\nwant the finding %s", fixture, strings.Join(f, "\n"), w)
			}
		}
		for _, s := range f {
			psb := fixture == "netimport" && (strings.HasPrefix(s, "rule 5: internal/psb depends on ") || strings.HasPrefix(s, "engine checks: internal/psb depends on "))
			if !slices.Contains(want, s) && !psb {
				t.Errorf("%s breaks another rule too: %s", fixture, s)
			}
		}
	}
}

// load reads the module at root: go.mod by go mod edit -json, its packages
// and their dependencies by go list -deps -json, and the non-test Go files of
// its packages, also the ones that a build constraint leaves out. An import
// that only such a file has is listed too, with its dependencies.
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
	listed := map[string]bool{"C": true} // C is the cgo pseudo-package
	list := func(args ...string) {
		for dec := json.NewDecoder(bytes.NewReader(goCmd(t, root, append([]string{"list", "-deps", "-json"}, args...)...))); dec.More(); {
			var p goPackage
			if err := dec.Decode(&p); err != nil {
				t.Fatal(err)
			}
			if !listed[p.ImportPath] {
				listed[p.ImportPath] = true
				m.packages = append(m.packages, p)
			}
		}
	}
	list("./...")
	for _, p := range m.packages {
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
	var more []string
	for _, files := range m.sources {
		for _, imp := range fileImports(files) {
			if !listed[imp] && !slices.Contains(more, imp) {
				more = append(more, imp)
			}
		}
	}
	if len(more) > 0 {
		list(append([]string{"-e"}, more...)...) // -e: a package of no module is a finding of rule 2
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

// No non-test Go file of the module reads an environment variable or the
// standard input, except environ of internal/git; the same scan finds the
// seeded read of a fixture module (docs/spec/README.md, Commands: Arguments;
// NFR-005).
func TestInputRule(t *testing.T) {
	real := load(t, filepath.Join("..", ".."))
	if len(real.sources["internal/git"]) == 0 || len(real.sources["internal/cli"]) == 0 {
		t.Fatal("no source of internal/git or internal/cli was read, so the scan checks nothing")
	}
	if f := checkInputs(real); len(f) != 0 {
		t.Errorf("the module breaks the input rule:\n%s", strings.Join(f, "\n"))
	}
	want := "input rule: internal/cli/cli.go:9 reads os.Getenv"
	if f := checkInputs(load(t, filepath.Join("testdata", "envread"))); !slices.Equal(f, []string{want}) {
		t.Errorf("the fixture gives\n%s\nwant only the finding %s", strings.Join(f, "\n"), want)
	}
}
