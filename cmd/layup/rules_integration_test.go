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

// The package rules of docs/spec/packages.md hold on the real module, and the
// same checker finds the seeded breaches of rule 5 in a fixture module: an
// import of net/http, and one of net/smtp behind a build constraint (NFR-005,
// NFR-007).
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
	for _, want := range []string{"rule 5: cmd/layup depends on net/http", "rule 5: internal/psb depends on net"} {
		if !slices.Contains(f, want) {
			t.Errorf("the fixture gives\n%s\nwant the finding %s", strings.Join(f, "\n"), want)
		}
	}
	for _, s := range f {
		if !strings.HasPrefix(s, "rule 5: ") {
			t.Errorf("the fixture breaks another rule too: %s", s)
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

// #48: a checkout with core.autocrlf=true (git's default on Windows) keeps the
// bytes of the raw facts files and of setup-check.sh, so check facts passes on
// it. The test clones the commit HEAD of this repository; a change of
// .gitattributes counts when it is committed.
func TestAnAutocrlfCheckoutKeepsTheFacts(t *testing.T) {
	home := t.TempDir()
	env := append(os.Environ(), "HOME="+home, "XDG_CONFIG_HOME="+home, "GIT_CONFIG_NOSYSTEM=1")
	run := func(dir string, name string, args ...string) string {
		t.Helper()
		cmd := exec.Command(name, args...)
		cmd.Dir, cmd.Env = dir, env
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, out)
		}
		return strings.TrimSpace(string(out))
	}
	root, _ := filepath.Abs(filepath.Join("..", ".."))
	head := run(root, "git", "rev-parse", "HEAD")
	clone := filepath.Join(t.TempDir(), "clone")
	run("", "git", "-c", "core.hooksPath=/dev/null", "clone", "-q", "--no-hardlinks", "--no-checkout", "--config", "core.autocrlf=true", root, clone)
	run(clone, "git", "-c", "core.hooksPath=/dev/null", "checkout", "-q", "--detach", head)
	cmd := exec.Command("sh", filepath.Join(clone, "docs", "setup", "setup-check.sh"), "--only", "facts", clone)
	cmd.Env = env
	if out, err := cmd.CombinedOutput(); err != nil || !strings.Contains(string(out), "setup-check: facts OK") {
		t.Errorf("check facts on an autocrlf checkout: %v\n%s", err, out)
	}
}
