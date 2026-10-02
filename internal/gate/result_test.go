package gate

import (
	"errors"
	"testing"
)

func kind(t *testing.T, name, state, scope string) Kind {
	t.Helper()
	p, err := parsePattern(scope)
	if err != nil {
		t.Fatal(err)
	}
	return Kind{Name: name, State: state, Tool: "go", Command: "go vet ./...", Scope: []pattern{p}}
}

func found(string) (string, error)   { return "/usr/bin/go", nil }
func missing(string) (string, error) { return "", errors.New("not found") }

func TestDecideActiveByTheFirstLineThatMatches(t *testing.T) {
	k := kind(t, "static", "active", "./*.go")
	for _, c := range []struct {
		name           string
		lookPath       func(string) (string, error)
		files          []string
		o              outcome
		result, reason string
		runs           int
	}{
		{"a missing tool", missing, []string{"x.go"}, outcome{}, notActive, "tool not found: go", 0},
		{"a missing tool and no product path", missing, []string{"README.md"}, outcome{}, notActive, "tool not found: go", 0},
		{"no product path", found, []string{"README.md"}, outcome{}, clear, "no product path", 0},
		{"exit 0", found, []string{"x.go"}, outcome{}, pass, "", 1},
		{"exit 3", found, []string{"a/x.go"}, outcome{code: 3}, fail, "exit 3", 1},
		{"a signal", found, []string{"x.go"}, outcome{signal: "terminated"}, fail, "signal terminated", 1},
	} {
		runs := 0
		result, reason := decideActive(k, c.files, c.lookPath, func(string) outcome { runs++; return c.o })
		if result != c.result || reason != c.reason || runs != c.runs {
			t.Errorf("%s: %s %q with %d runs; want %s %q with %d", c.name, result, reason, runs, c.result, c.reason, c.runs)
		}
	}
}

func TestDecidePendingNamesTheFirstChangedProductPath(t *testing.T) {
	k := kind(t, "layout", "pending", "./*.go")
	if r, why := decidePending(k, []string{"README.md", "a/b.go", "c.go"}); r != fail || why != "pending: product path changed: a/b.go" {
		t.Errorf("a changed product path: %s %q", r, why)
	}
	if r, why := decidePending(k, []string{"README.md", "go.mod"}); r != clear || why != "pending: no product path" {
		t.Errorf("no changed product path: %s %q", r, why)
	}
	if r, why := decidePending(k, nil); r != clear || why != "pending: no product path" {
		t.Errorf("no change: %s %q", r, why)
	}
}

// NFR-004: over every input of the rules, a kind is pass only when it is
// active, its tool is found, it has a product path and its command exits 0.
// A pending kind with a missing tool gives its pending result (K19).
func TestACheckThatDidNotRunNeverPasses(t *testing.T) {
	for _, state := range []string{"active", "pending"} {
		for _, look := range []func(string) (string, error){found, missing} {
			for _, files := range [][]string{{"x.go"}, {"README.md"}} {
				for _, o := range []outcome{{}, {code: 1}, {signal: "killed"}} {
					k := kind(t, "test", state, "./*.go")
					var result string
					if state == "pending" {
						result, _ = decidePending(k, files)
					} else {
						result, _ = decideActive(k, files, look, func(string) outcome { return o })
					}
					_, err := look("go")
					ran := state == "active" && err == nil && files[0] == "x.go"
					if result == pass && !(ran && o == (outcome{})) {
						t.Errorf("%s, tool found %v, files %q, %+v: pass for a check that did not pass", state, err == nil, files, o)
					}
					if ran && o == (outcome{}) && result != pass {
						t.Errorf("%s, files %q: %s; want pass", state, files, result)
					}
				}
			}
		}
	}
}
