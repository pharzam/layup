package run

import "testing"

// The step shape (task T-3f94, #364; docs/spec/loop.md, REQ-012 — Shape, its base).

func attempt(n, base string) []string { return []string{n, "attempt", n, "—", base, "", "—"} }

func TestTheArchitectureTaskStartsAtTheVerifiedHeadOfTheSpecification(t *testing.T) {
	const h, x = "1111111111111111111111111111111111111111", "2222222222222222222222222222222222222222"
	verified := Advance{Task: "T-sp01", Verified: true, Head: h}
	// The control: verified at H, no attempt of the architecture task: it starts at H.
	if base, stop, err := archStart(verified, "T-ar01", nil); base != h || stop != "" || err != nil {
		t.Errorf("archStart(verified, no events) = %q, %q, %v; want %s", base, stop, err, h)
	}
	if base, stop, err := archStart(verified, "T-ar01", [][]string{attempt("1", h), attempt("2", h)}); base != h || stop != "" || err != nil {
		t.Errorf("archStart(two attempts at H) = %q, %q, %v; want %s, no stop", base, stop, err, h)
	}
	want := "T-ar01: its base " + x + " is not the verified head " + h + " of T-sp01"
	if _, stop, err := archStart(verified, "T-ar01", [][]string{attempt("1", h), attempt("2", x)}); stop != want || err != nil {
		t.Errorf("archStart(an attempt at X) stop = %q, %v; want %q", stop, err, want)
	}
	if _, _, err := archStart(Advance{Task: "T-sp01", Verified: true}, "T-ar01", nil); err == nil {
		t.Error("archStart(verified, no head) = nil; want an error")
	}
	// The specification in a wait or a stop: no start, no stop of its own.
	if base, stop, err := archStart(Advance{Task: "T-sp01", Wait: "waiting for the Operator: x"}, "T-ar01", nil); base != "" || stop != "" || err != nil {
		t.Errorf("archStart(the specification in a wait) = %q, %q, %v; want none", base, stop, err)
	}
}
