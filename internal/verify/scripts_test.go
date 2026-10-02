package verify

import (
	"errors"
	"slices"
	"testing"
	"testing/fstest"
)

// The checks discipline-tests and link-lint run the baseline's own scripts,
// with the rule of a gate command (D5 of #87, notes 5 and 6 of its plan
// review).
func TestTheBaselineScripts(t *testing.T) {
	scripts := fstest.MapFS{"docs/tests/run-discipline-tests.sh": {Data: []byte("x")}, "docs/links/link-lint.sh": {Data: []byte("x")}}
	found := func(string) (string, error) { return "/bin/sh", nil }
	for _, c := range []struct {
		name     string
		fsys     fstest.MapFS
		lookPath func(string) (string, error)
		o        scriptOutcome
		want     Row
	}{
		{"exit 0", scripts, found, scriptOutcome{}, Row{"link-lint", "pass", ""}},
		{"exit 1", scripts, found, scriptOutcome{code: 1}, Row{"link-lint", "fail", "exit 1"}},
		{"a signal", scripts, found, scriptOutcome{signal: "terminated"}, Row{"link-lint", "fail", "signal terminated"}},
		{"sh found and not started", scripts, found, scriptOutcome{code: -1}, Row{"link-lint", "fail", "exit -1"}},
		{"no script", fstest.MapFS{}, found, scriptOutcome{}, Row{"link-lint", "not-active", "missing: docs/links/link-lint.sh"}},
		{"no sh", scripts, func(string) (string, error) { return "", errors.New("not found") }, scriptOutcome{}, Row{"link-lint", "not-active", "tool not found: sh"}},
	} {
		ran := false
		got := decideScript("link-lint", "docs/links/link-lint.sh", c.fsys, c.lookPath, func(script string) scriptOutcome {
			ran = true
			return c.o
		})
		if got != c.want || ran != (c.want.Result != "not-active") {
			t.Errorf("%s: %q, ran %v; want %q", c.name, got, ran, c.want)
		}
	}
}

// The call of K10 for S05: the files that the lines of link-lint.sh on its
// standard error name (D6 of #87, condition 1 of its plan review).
func TestTheBrokenLinks(t *testing.T) {
	for _, c := range []struct {
		name   string
		code   int
		stderr string
		want   []string
		err    bool
	}{
		{"no broken link", 0, "", nil, false},
		{"the files of the lines, once each, in byte order", 1,
			"FAIL  L2: docs/b.md:3 links docs/decisions/x.md, but that file has no heading with that anchor\n" +
				"FAIL  L6: docs/a.md:9 uses reference label [x], but this file defines no [x]: target\n" +
				"FAIL  L4: docs/b.md:7 links ../x.md, which escapes the repository root\n", []string{"docs/a.md", "docs/b.md"}, false},
		{"a line that names no file", 1, "FAIL  L5: no in-tree link was resolved \u2014 this run checked nothing\n", nil, true},
		{"no root", 1, "FAIL  link-lint: root not found: x\n", nil, true},
		{"a line that names no file beside one that does", 1,
			"FAIL  L2: docs/a.md:3 links x.md, but that file has no heading with that anchor\nFAIL  L5: no in-tree link was resolved \u2014 this run checked nothing\n", nil, true},
		{"a failure with no line", 2, "sh: syntax error\n", nil, true},
	} {
		got, err := brokenLinks(c.code, []byte(c.stderr))
		if !slices.Equal(got, c.want) || (err != nil) != c.err {
			t.Errorf("%s: %q, %v; want %q, an error %v", c.name, got, err, c.want, c.err)
		}
	}
}
