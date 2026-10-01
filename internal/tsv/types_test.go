package tsv

import (
	"strings"
	"testing"
)

func TestParseTypeTakesEachTypeOfTheClosedList(t *testing.T) {
	for _, expr := range []string{
		"text", "int", "decimal", "bool", "time", "sha1", "sha256", "path",
		"enum(active|pending)", "enum(G1|G2|G3|G4|G5)",
		"id(Q-NNN)", "id(SNN)", "id(ST-NNN)", "id(S-xxxxxxxx)", "id(T-xxxx)", "id(<word>)",
		"list(text)", "list(path)", "list(enum(in|out))",
	} {
		if _, err := parseType(expr); err != nil {
			t.Errorf("parseType(%q): %v; want no error", expr, err)
		}
	}
}

func TestParseTypeRefusesATypeOffTheList(t *testing.T) {
	for _, expr := range []string{
		"", "float", "Text", "int ", "text(1)",
		"enum()", "enum(a||b)", "enum(a|a)", "enum(a b)", "enum(a|b", "enum(a|b)c",
		"id()", "id(Q-<num>)", "id(Q NNN)", "id(Q-N)N)",
		"list()", "list(float)", "list(list(text))", "list(text))",
	} {
		if _, err := parseType(expr); err == nil {
			t.Errorf("parseType(%q): no error; want one", expr)
		}
	}
}

func TestTypeCheckTakesTheValuesOfItsType(t *testing.T) {
	for _, c := range []struct {
		typ       string
		good, bad []string
	}{
		{"text", []string{"a line", " ", "Q-001"}, []string{"", "a\tb", "a\nb", "a\rb", "\xff"}},
		{"int", []string{"0", "7", "1000", "12345678901234567890"}, []string{"", "-1", "+1", "07", "00", "1.0", "1e3", " 1"}},
		// D3: a decimal needs its point; the part before it is an int.
		{"decimal", []string{"0.42", "3.0", "10.50"}, []string{"3", ".5", "5.", "-0.5", "03.5", "1e3", "1.5e3", "1,5"}},
		{"bool", []string{"yes", "no"}, []string{"Yes", "true", "1"}},
		{"time", []string{"2026-10-01T08:09:21Z", "2024-02-29T23:59:59Z"}, []string{
			"2026-10-01T08:09:21", "2026-10-01T08:09:21.5Z", "2026-10-01T08:09:21+00:00",
			"2026-10-01T08:09:21z", "2026-13-01T08:09:21Z", "2026-02-29T00:00:00Z", "2026-10-01 08:09:21Z",
		}},
		{"sha1", []string{"7cdd34633e7be741e98e3013c7950109c8c7e7cf"}, []string{
			"7CDD34633E7BE741E98E3013C7950109C8C7E7CF", "7cdd346", "7cdd34633e7be741e98e3013c7950109c8c7e7cg",
		}},
		{"sha256", []string{strings.Repeat("ab", 32)}, []string{strings.Repeat("ab", 31) + "a", strings.Repeat("AB", 32)}},
		{"path", []string{"docs/gates.tsv", "a", ".golangci.yml", "a..b/c"}, []string{"", "/docs", "../x", "a/../b", "a//b", "a/", ".", "./a", "a\tb"}},
		{"enum(active|pending)", []string{"active", "pending"}, []string{"Active", "other", "active "}},
		// D2, note 2: a run of N is a minimum count of digits; x is one character.
		{"id(Q-NNN)", []string{"Q-001", "Q-999", "Q-1000"}, []string{"Q-01", "q-001", "Q001", "Q-00a", "Q-001 "}},
		{"id(SNN)", []string{"S01", "S15"}, []string{"S1", "s01", "S-01"}},
		{"id(T-xxxx)", []string{"T-18v6", "T-0drh"}, []string{"T-18V6", "T-18v", "T-18v66", "T-18-6"}},
		{"id(<word>)", []string{"static", "not-active"}, []string{"Static", "a1", "a b", "a_b"}},
		// D7: one space separates the values of a list.
		{"list(text)", []string{"a", "./*.go docs/"}, []string{"a  b", " a", "a ", "a\tb"}},
		{"list(path)", []string{"a b/c"}, []string{"a /x", "a ../b"}},
		{"list(enum(x|y))", []string{"x y x"}, []string{"x z"}},
	} {
		typ, err := parseType(c.typ)
		if err != nil {
			t.Fatalf("parseType(%q): %v", c.typ, err)
		}
		for _, v := range c.good {
			if err := typ.check(v); err != nil {
				t.Errorf("%s: check(%q): %v; want no error", c.typ, v, err)
			}
		}
		for _, v := range c.bad {
			if err := typ.check(v); err == nil {
				t.Errorf("%s: check(%q): no error; want one", c.typ, v)
			}
		}
	}
}
