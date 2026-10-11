package run

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/ledger"
)

// startTSV is a valid start.tsv of one harness, fake, with intake.cap.
func startTSV(intakeCap string) []byte {
	sha1, sha256 := strings.Repeat("a", 40), strings.Repeat("b", 64)
	return []byte(strings.Join([]string{"name\tvalue\tsource",
		"layup.version\tv0.1.0\trun", "psb.sha256\t" + sha256 + "\tcommand", "vision.sha256\t—\tcommand",
		"forge.plan\tfree\tcommand", "forge.visibility\tpublic\tforge", "app.permissions\tcontents:write\tforge",
		"operator.id\t123\tforge", "idea-owner.id\t124\tforge", "intake.cap\t" + intakeCap + "\tcommand",
		"lease.H\t5\tcommand", "watch.T\t10\tcommand", "harness.fake.cap\t—\tregister", "harness.fake.wall\t1\tregister",
		"pin.source\thttps://example.invalid/armature\trun", "pin.commit\t" + sha1 + "\trun", "pin.tree\t" + sha1 + "\trun",
		"pin.time\t2026-10-07T12:00:00Z\trun", "issue.intake\t1\tforge", "issue.control\t2\tforge", "watch\tconfirmed\trun",
	}, "\n") + "\n")
}

// The demo of row 42b at the unit level: with the cap of milestone 0 passed,
// PhaseSession gives ErrMilestoneCap, draws no session ID and makes no
// commit; the detail is fixed by the records, with no clock reading (round 1
// of #216).
func TestTheCapOfMilestoneZeroDrawsNoSession(t *testing.T) {
	for _, c := range []struct {
		name, sessionCap string
		later            time.Duration
		want             string
	}{
		{"the money", "1.01", 0, "the cap of milestone 0: the known money 0.0 plus the session's cap 1.01 passes 1.0 USD"},
		{"the hours", "", 9 * time.Hour, "the cap of milestone 0: the time since Start passes 8.0 hours"},
	} {
		store := &standinStore{files: map[string][]byte{"start/start.tsv": startTSV("1.0,8.0")}, holder: "0123456789abcdef"}
		r := newSessions(store)
		r.Now = func() time.Time { return sessionAt.Add(c.later) }
		drawn := 0
		r.newID = func() (string, error) { drawn++; return "S-1a2b3c4d", nil }
		h := []string{"fake", c.sessionCap, "1", "harness", "stdin", "version", "", "", "R.md", "", "none", "subscription", ""}
		m := []string{"fake", "m1", "1000", "https://example.invalid", "2026-10-09T12:00:00Z", "yes", ""}
		id, err := r.PhaseSession(context.Background(), PhaseSpec{Step: "meaning", Role: "systems-architect", Harness: h, Model: m,
			Models: [][]string{m}, Base: aSHA, Records: aSHA, Output: "meaning.tsv"})
		if !errors.Is(err, ledger.ErrMilestoneCap) || err.Error() != c.want || id != "" {
			t.Errorf("%s: %q, %v; want ErrMilestoneCap, %q", c.name, id, err, c.want)
		}
		if drawn != 0 || len(store.commits) != 0 {
			t.Errorf("%s: %d session IDs drawn, %d commits; want none", c.name, drawn, len(store.commits))
		}
	}
}
