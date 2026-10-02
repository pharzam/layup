package main

import (
	"fmt"
	"testing"
)

// sameBytes gives nil when a and b are equal, and else an error that names
// the first byte where they differ.
func sameBytes(a, b []byte) error {
	for i := range min(len(a), len(b)) {
		if a[i] != b[i] {
			return fmt.Errorf("byte %d differs: %q and %q", i, a[i], b[i])
		}
	}
	if len(a) != len(b) {
		return fmt.Errorf("the lengths differ: %d and %d bytes", len(a), len(b))
	}
	return nil
}

func TestSameBytes(t *testing.T) {
	for _, c := range []struct {
		a, b string
		same bool
	}{
		{"layup 0.1.0-dev\n", "layup 0.1.0-dev\n", true},
		{"", "", true},
		{"layup 0.1.0-dev\n", "layup 0.1.1-dev\n", false},
		{"layup\n", "layup", false},
	} {
		if err := sameBytes([]byte(c.a), []byte(c.b)); (err == nil) != c.same {
			t.Errorf("sameBytes(%q, %q) = %v, want same %v", c.a, c.b, err, c.same)
		}
	}
}
