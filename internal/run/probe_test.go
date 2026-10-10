package run

import (
	"errors"
	"testing"

	"github.com/pharzam/layup/internal/session"
)

func TestTheReasonOfAProbe(t *testing.T) {
	d := session.Dir{Root: "/h/sessions/S-1a2b3c4d", Repo: "/h/sessions/S-1a2b3c4d/repo"}
	models := [][]string{{"claude", "m1", "1000", "https://x", "2026-10-09T12:00:00Z", "yes", ""}, {"claude", "m2", "1000", "https://x", "2026-10-09T12:00:00Z", "no", "too dear"}}
	const token = "0123456789abcdef"
	rows := func(r ...[]string) [][]string { return r }
	good := rows([]string{"token", token}, []string{"file", "AGENTS.md"})
	for _, c := range []struct {
		name      string
		class     string
		resultErr error
		rows      [][]string
		policy    []string
		used      []string
		want      string
	}{
		{"a probe that passes", "done", nil, good, nil, []string{"m1"}, ""},
		{"a policy path outside, of the row", "done", nil, append(good, []string{"file", "/etc/claude/policy.json"}), []string{"/etc/claude/policy.json"}, []string{"m1"}, ""},
		{"a file inside the session directory, given absolute", "done", nil, append(good, []string{"file", "/h/sessions/S-1a2b3c4d/home/.claude/CLAUDE.md"}), nil, []string{"m1"}, ""},
		{"a file inside repo/ by a relative path with ..", "done", nil, append(good, []string{"file", "docs/../AGENTS.md"}), nil, []string{"m1"}, ""},
		{"an end that is not done", "crash", nil, good, nil, nil, "crash"},
		{"another token", "done", nil, rows([]string{"token", "ffffffffffffffff"}, []string{"file", "AGENTS.md"}), nil, nil, "token"},
		{"no token", "done", nil, rows([]string{"file", "AGENTS.md"}), nil, nil, "token"},
		{"no AGENTS.md", "done", nil, rows([]string{"token", token}, []string{"file", "CLAUDE.md"}), nil, nil, "files"},
		{"a file outside", "done", nil, append(good, []string{"file", "/home/op/.claude/CLAUDE.md"}), nil, nil, "outside"},
		{"a relative path that leaves the session directory", "done", nil, append(good, []string{"file", "../../../x.md"}), nil, nil, "outside"},
		// A relative value is read from repo/: ../home/ is inside.
		{"a relative path into home/", "done", nil, append(good, []string{"file", "../home/.claude/CLAUDE.md"}), nil, []string{"m1"}, ""},
		{"an absolute path whose .. leaves the session directory", "done", nil, append(good, []string{"file", "/h/sessions/S-1a2b3c4d/../../etc/x.md"}), nil, nil, "outside"},
		{"a model of use no", "done", nil, good, nil, []string{"m1", "m2"}, "not-used"},
	} {
		if got := probeReason(c.class, c.resultErr, c.rows, token, d, c.policy, c.used, models, "claude"); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	// A probe.tsv that its block refuses makes the class no-result: the class
	// is the reason.
	if got := probeReason("no-result", errors.New("malformed"), nil, token, d, nil, nil, models, "claude"); got != "no-result" {
		t.Errorf("no valid probe.tsv: %q, want no-result", got)
	}
}

func TestTheRoutingRegistersDiffer(t *testing.T) {
	a := [][]string{{"developer", "execution", "1", "claude", "m1"}, {"reviewer", "reasoning", "1", "claude", "m2"}}
	b := [][]string{{"reviewer", "reasoning", "1", "claude", "m2"}, {"developer", "execution", "1", "claude", "m1"}}
	if routingDiffers(a, b) {
		t.Error("the same rows in another order: they differ; want the same table")
	}
	if !routingDiffers(a, a[:1]) || !routingDiffers(nil, a) {
		t.Error("another table: the same; want them to differ")
	}
	c := [][]string{{"developer", "execution", "1", "claude", "m3"}, {"reviewer", "reasoning", "1", "claude", "m2"}}
	if !routingDiffers(a, c) {
		t.Error("another model in one row: the same; want them to differ")
	}
}
