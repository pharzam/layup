//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The problem statement of the tests of layup spec check, and its span tables.
const (
	specPSB    = "Ab cd\n"
	spansHead  = "start\tend\tfact\tclass\treason\ttext\n"
	spansCover = spansHead + "0\t2\t1\tneed\t—\tAb\n3\t5\t2\tconstraint\t—\tcd\n"
	spansTwice = spansHead + "0\t4\t1\tneed\t—\tAb c\n3\t5\t2\tconstraint\t—\tcd\n"
	allPass    = "rule\tresult\tdetail\ncover\tpass\t—\nbounds\tpass\t—\ntext\tpass\t—\nids\tpass\t—\nclasses\tpass\t—\n"
)

// specFiles writes the problem statement and a span table and gives their paths.
func specFiles(t *testing.T, psb, spans string) (string, string) {
	t.Helper()
	dir := t.TempDir()
	p, s := filepath.Join(dir, "psb.md"), filepath.Join(dir, "spans.tsv")
	if err := os.WriteFile(p, []byte(psb), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(s, []byte(spans), 0o644); err != nil {
		t.Fatal(err)
	}
	return p, s
}

// The demo of row 41b: the built binary prints the table of the span check
// and exits 0 on a span table that covers the file, the same bytes on a
// repeat (NFR-005), and 1 on a span table with a byte in two spans.
func TestSpecCheckPassesAndFails(t *testing.T) {
	p, s := specFiles(t, specPSB, spansCover)
	if r := repeat(t, "spec", "check", "--facts", s, "--psb", p); r.code != 0 || r.stdout != allPass || r.stderr != "" {
		t.Errorf("a span table that covers the file: exit %d, stdout %q, stderr %q; want 0, five pass rows and nothing", r.code, r.stdout, r.stderr)
	}
	p, s = specFiles(t, specPSB, spansTwice)
	r := repeat(t, "spec", "check", "--facts", s, "--psb", p)
	if r.code != 1 || !strings.Contains(r.stdout, "\ncover\tfail\tbyte 3 is in 2 spans\n") || r.stderr != "" {
		t.Errorf("a byte in two spans: exit %d, stdout %q, stderr %q; want 1 and the row cover fail", r.code, r.stdout, r.stderr)
	}
}

// A file that is missing or not valid UTF-8, and a span table that its block
// refuses, give exit 2, nothing on standard output, and the reason, which
// names the file, with no usage.
func TestSpecCheckInputErrors(t *testing.T) {
	p, s := specFiles(t, specPSB, spansCover)
	dir := filepath.Dir(p)
	bad := filepath.Join(dir, "bad.md")
	if err := os.WriteFile(bad, []byte("Ab\n\xff\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	refused := filepath.Join(dir, "refused.tsv")
	if err := os.WriteFile(refused, []byte("start\tend\n0\t2\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	missing := filepath.Join(dir, "missing")
	for _, c := range []struct {
		name, facts, psb, reason string
	}{
		{"a missing problem statement", s, missing, missing + ": no such file or directory"},
		{"a missing span table", missing, p, missing + ": no such file or directory"},
		{"a problem statement not valid UTF-8", s, bad, "layup: " + bad + ": line 2 is not valid UTF-8\n"},
		{"a span table not valid UTF-8", bad, p, "layup: " + bad + ": line 2 is not valid UTF-8\n"},
		{"a span table that its block refuses", refused, p, "layup: " + refused + ": "},
	} {
		r := layup(t, "spec", "check", "--facts", c.facts, "--psb", c.psb)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: ") || !strings.Contains(r.stderr, c.reason) || strings.Contains(r.stderr, "usage:") {
			t.Errorf("%s: exit %d, stdout %q, stderr %q; want 2, nothing, and %q with no usage", c.name, r.code, r.stdout, r.stderr, c.reason)
		}
	}
}

// The usage errors of layup spec check, each with its reason and the usage.
func TestSpecCheckUsageErrors(t *testing.T) {
	p, s := specFiles(t, specPSB, spansCover)
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{[]string{"spec", "check"}, "missing flag --facts"},
		{[]string{"spec", "check", "--facts", s}, "missing flag --psb"},
		{[]string{"spec", "check", "--facts", s, "--facts", s, "--psb", p}, "flag --facts given twice"},
		{[]string{"spec", "check", "--facts=", "--psb", p}, "flag --facts has no value"},
		{[]string{"spec", "check", "--facts", s, "--psb", p, "extra"}, "extra argument \"extra\""},
	} {
		r := layup(t, c.args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: "+c.reason+"\n") || !strings.Contains(r.stderr, "\nusage: layup ") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing, and %q with the usage", c.args, r.code, r.stdout, r.stderr, c.reason)
		}
	}
}
