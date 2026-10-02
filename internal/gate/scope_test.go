package gate

import "testing"

func TestScopePatterns(t *testing.T) {
	for _, c := range []struct {
		pattern, path string
		match         bool
	}{
		{"./*.go", "x.go", true},
		{"./*.go", "a/b/c.go", true},
		{"./*.go", "go.mod", false},
		{"./*.go", "a/b/c.go.txt", false},
		{"./*.go", ".go", false},
		{"internal/*.go", "internal/a/b.go", true},
		{"internal/*.go", "internal/x.go", true},
		{"internal/*.go", "cmd/x.go", false},
		{"internal/*.go", "internal2/x.go", false},
		{"internal", "internal/x.go", true},
		{"internal", "internal", true},
		{"internal", "internal2/x.go", false},
		{"docs/gates.tsv", "docs/gates.tsv", true},
		{"docs/gates.tsv", "docs/gates.tsv.bak", false},
	} {
		p, err := parsePattern(c.pattern)
		if err != nil {
			t.Fatalf("parsePattern(%q): %v", c.pattern, err)
		}
		if got := p.match(c.path); got != c.match {
			t.Errorf("%q matches %q: %v, want %v", c.pattern, c.path, got, c.match)
		}
	}
}

func TestOtherFormsOfAPatternAreErrors(t *testing.T) {
	for _, s := range []string{"*.go", "P/*", "P/*/x.go", "./*.", "./**.go", "a/*.go/b", "/abs/*.go", "a/../*.go", "", ".", "a/", "./*.g*"} {
		if _, err := parsePattern(s); err == nil {
			t.Errorf("parsePattern(%q): no error", s)
		}
	}
}
