package route

import (
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

const harnessesRecordHeader = "session\tharness\tversion\tmodel\tresult\treason\tfiles\tmodels\tend"

// probes reads rows of records:harnesses.tsv by records.ReadHarnesses.
func probes(t *testing.T, rows ...string) [][]string {
	t.Helper()
	r, err := records.ReadHarnesses(table(harnessesRecordHeader, rows...))
	if err != nil {
		t.Fatal(err)
	}
	return r
}

const (
	passed21  = "S-00000001\tclaude\t2.1.295\tclaude-opus-5-5\tpassed\t—\tAGENTS.md\tclaude-opus-5-5\t2026-10-09T12:00:00Z"
	failed21  = "S-00000002\tclaude\t2.1.295\tclaude-opus-5-5\tfailed\ttoken\tAGENTS.md\t—\t2026-10-09T12:05:00Z"
	passed22  = "S-00000003\tclaude\t2.1.296\tclaude-opus-5-5\tpassed\t—\tAGENTS.md\tclaude-opus-5-5\t2026-10-09T12:10:00Z"
	devinPass = "S-00000004\tdevin\t3000.11.3\tswe-2-high\tpassed\t—\tAGENTS.md\t—\t2026-10-09T12:15:00Z"
)

func TestProbed(t *testing.T) {
	for name, c := range map[string]struct {
		rows    []string
		version string
		want    bool
	}{
		"a passed probe at the version":                    {[]string{passed21}, "2.1.295", true},
		"a later failed probe at the same version":         {[]string{passed21, failed21}, "2.1.295", false},
		"a passed probe after a failed one":                {[]string{failed21, passed21}, "2.1.295", true},
		"a passed probe at another version only":           {[]string{passed22}, "2.1.295", false},
		"a later probe at another version is not this one": {[]string{passed21, strings.Replace(failed21, "2.1.295", "2.1.296", 1)}, "2.1.295", true},
		"a later failed probe of another harness at the same version": {[]string{passed21,
			"S-00000005\tdevin\t2.1.295\tswe-2-high\tfailed\ttoken\t—\t—\t2026-10-09T12:20:00Z"}, "2.1.295", true},
		"a harness with no row": {[]string{devinPass}, "2.1.295", false},
		"no row at all":         {nil, "2.1.295", false},
	} {
		if got := Probed(probes(t, c.rows...), "claude", c.version); got != c.want {
			t.Errorf("%s: %v, want %v", name, got, c.want)
		}
	}
}

func TestAdmitted(t *testing.T) {
	models := rowsOf(t, ReadModels, table(modelsHeader, opusRow, haikuRow, sweRow))
	record := probes(t, passed21, devinPass)
	for name, c := range map[string]struct {
		harness, version, model string
		want                    bool
	}{
		"a probed harness and a model in use":       {"claude", "2.1.295", "claude-opus-5-5", true},
		"a model of use no":                         {"claude", "2.1.295", "claude-haiku-5-5", false},
		"a model with no row of models.tsv":         {"claude", "2.1.295", "claude-fable-5-1", false},
		"a model of another harness":                {"claude", "2.1.295", "swe-2-high", false},
		"a version with no passed probe":            {"claude", "2.1.296", "claude-opus-5-5", false},
		"another harness, probed, its model in use": {"devin", "3000.11.3", "swe-2-high", true},
	} {
		if got := Admitted(record, models, c.harness, c.version, c.model); got != c.want {
			t.Errorf("%s: %v, want %v", name, got, c.want)
		}
	}
}

const routingRecordHeader = "role\ttier\tposition\tharness\tmodel"

func TestPair(t *testing.T) {
	// The file's order is not the order of position.
	routing, err := records.ReadRouting(table(routingRecordHeader,
		"developer\texecution\t2\tdevin\tswe-2-high",
		"developer\texecution\t1\tclaude\tclaude-opus-5-5",
		"developer\texecution\t3\tcodex\tgpt-6",
		"developer\treasoning\t1\tclaude\tclaude-fable-5-1",
		"verifier\texecution\t1\tcodex\tgpt-6"))
	if err != nil {
		t.Fatal(err)
	}
	all := func(string, string) bool { return true }
	not := func(h string) func(string, string) bool { return func(harness, _ string) bool { return harness != h } }
	only := func(h string) func(string, string) bool { return func(harness, _ string) bool { return harness == h } }
	for name, c := range map[string]struct {
		role, tier     string
		admitted       func(string, string) bool
		harness, model string
	}{
		"the first by position":                     {"developer", "execution", all, "claude", "claude-opus-5-5"},
		"the second when the first is not admitted": {"developer", "execution", not("claude"), "devin", "swe-2-high"},
		"the third when the first two are not":      {"developer", "execution", only("codex"), "codex", "gpt-6"},
		"the tier's own list":                       {"developer", "reasoning", all, "claude", "claude-fable-5-1"},
		"the role's own list":                       {"verifier", "execution", all, "codex", "gpt-6"},
	} {
		h, m, err := Pair(routing, c.role, c.tier, c.admitted)
		if err != nil || h != c.harness || m != c.model {
			t.Errorf("%s: %s %s %v, want %s %s", name, h, m, err, c.harness, c.model)
		}
	}
	for name, c := range map[string]struct {
		role, tier string
		admitted   func(string, string) bool
	}{
		"no admitted pair":  {"developer", "execution", func(string, string) bool { return false }},
		"a role of no list": {"tester", "execution", all},
		"a tier of no list": {"verifier", "reasoning", all},
	} {
		if _, _, err := Pair(routing, c.role, c.tier, c.admitted); !errors.Is(err, ErrNoPair) {
			t.Errorf("%s: %v, want ErrNoPair", name, err)
		}
	}
}
