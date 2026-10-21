package run

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
)

func TestTheReasonOfAProbe(t *testing.T) {
	d := session.Dir{Root: "/h/sessions/S-1a2b3c4d", Repo: "/h/sessions/S-1a2b3c4d/repo"}
	models := [][]string{{"claude", "m1", "1000", "https://x", "2026-10-09T12:00:00Z", "yes", ""}, {"claude", "m2", "1000", "https://x", "2026-10-09T12:00:00Z", "no", "too dear"}}
	const want = "0123456789abcdef" // the token of the prompt
	rows := func(r ...[]string) [][]string { return r }
	good := rows([]string{"token", want}, []string{"file", "AGENTS.md"})
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
		{"no AGENTS.md", "done", nil, rows([]string{"token", want}, []string{"file", "CLAUDE.md"}), nil, nil, "files"},
		{"a file outside", "done", nil, append(good, []string{"file", "/home/op/.claude/CLAUDE.md"}), nil, nil, "outside"},
		{"a relative path that leaves the session directory", "done", nil, append(good, []string{"file", "../../../x.md"}), nil, nil, "outside"},
		// A relative value is read from repo/: ../home/ is inside.
		{"a relative path into home/", "done", nil, append(good, []string{"file", "../home/.claude/CLAUDE.md"}), nil, []string{"m1"}, ""},
		{"an absolute path whose .. leaves the session directory", "done", nil, append(good, []string{"file", "/h/sessions/S-1a2b3c4d/../../etc/x.md"}), nil, nil, "outside"},
		{"a model of use no", "done", nil, good, nil, []string{"m1", "m2"}, "not-used"},
	} {
		if got := probeReason(c.class, c.resultErr, c.rows, want, d, c.policy, c.used, models, "claude"); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	// A probe.tsv that its block refuses makes the class no-result: the class
	// is the reason.
	if got := probeReason("no-result", errors.New("malformed"), nil, want, d, nil, nil, models, "claude"); got != "no-result" {
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

// drawFrom is the random seam that gives each call the next bytes: the byte b is the character
// number b of the alphabet of a task ID, so {0,0,0,0} is T-0000 and {1,1,1,1} is T-1111.
func drawFrom(batches ...[4]byte) func([]byte) (int, error) {
	i := 0
	return func(b []byte) (int, error) {
		copy(b, batches[i][:])
		i++
		return len(b), nil
	}
}

// sessionRowOf is a row of sessions.tsv of the task.
func sessionRowOf(id, task string) []string {
	return []string{id, task, "1", "developer", "claude", "1.0", "m1", strings.Repeat("a", 40), strings.Repeat("b", 40), "12", "3", "1000", "", "1", "", ""}
}

func storeWithSessions(t *testing.T, rows ...[]string) *standinStore {
	t.Helper()
	store := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	data, err := table(records.SessionsSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	store.files["sessions.tsv"] = data
	return store
}

// The characterization of the draw as it was extracted: an ID that a row of sessions.tsv holds is
// drawn again (the draw of a probe and of a phase session; the register's rows join it in U2).
func TestADrawnIDIsDrawnAgainWhileASessionHoldsIt(t *testing.T) {
	s := newSessions(storeWithSessions(t, sessionRowOf("S-aaaaaaaa", "T-0000")))
	s.random = drawFrom([4]byte{0, 0, 0, 0}, [4]byte{1, 1, 1, 1})
	if id, err := s.taskID(context.Background()); err != nil || id != "T-1111" {
		t.Errorf("taskID = %q, %v; want T-1111, the second draw", id, err)
	}
}

// loop.md, The opening of a Shape task: the draw avoids the register's tasks as well, so a probe or a phase session never takes
// the ID of a task of the register.
func TestADrawnIDIsDrawnAgainWhileTheRegisterHoldsIt(t *testing.T) {
	store := storeWithSessions(t)
	var reg strings.Builder
	if err := WriteRegister(&reg, []Task{{ID: "T-0000", Issue: 4, Kind: "specification", Requirements: []string{"REQ-001"}, Opened: time0}}); err != nil {
		t.Fatal(err)
	}
	store.files["tasks.tsv"] = []byte(reg.String())
	s := newSessions(store)
	s.random = drawFrom([4]byte{0, 0, 0, 0}, [4]byte{1, 1, 1, 1})
	if id, err := s.taskID(context.Background()); err != nil || id != "T-1111" {
		t.Errorf("taskID = %q, %v; want T-1111, the second draw", id, err)
	}
}
