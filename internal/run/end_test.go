package run

import (
	"strconv"
	"strings"
	"testing"
)

// ev gives an events table of rows "kind attempt session".
func ev(rows ...string) [][]string {
	var t [][]string
	for i, r := range rows {
		f := strings.Fields(r)
		t = append(t, []string{itoa(i + 1), f[0], f[1], strings.ReplaceAll(f[2], "—", ""), "", "", "", "2026-10-09T12:00:00Z"})
	}
	return t
}

func itoa(n int) string { return strconv.Itoa(n) }

func TestTheOpenAttempt(t *testing.T) {
	const id = "S-1a2b3c4d"
	for _, c := range []struct {
		name   string
		events [][]string
		open   bool
	}{
		{"no later event", ev("attempt 1 —", "session 1 "+id), true},
		{"a closed of its attempt after", ev("attempt 1 —", "session 1 "+id, "closed 1 —"), false},
		{"a rebased of its attempt after", ev("attempt 1 —", "session 1 "+id, "rebased 1 —"), false},
		{"an attempt of another after", ev("attempt 1 —", "session 1 "+id, "attempt 2 —"), false},
		{"a closed of another attempt after", ev("attempt 1 —", "attempt 2 —", "session 1 "+id, "closed 2 —"), true},
		{"a closed of its attempt before", ev("attempt 1 —", "closed 1 —", "session 1 "+id), true},
		{"a session of another after", ev("attempt 1 —", "session 1 "+id, "session 1 S-99999999"), true},
		{"no event session", ev("attempt 1 —"), false},
	} {
		if got := openAttempt(c.events, id, 1); got != c.open {
			t.Errorf("%s: %v, want %v", c.name, got, c.open)
		}
	}
}
