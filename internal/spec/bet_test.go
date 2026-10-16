package spec

import (
	"strings"
	"testing"
)

// The checks of the bet (docs/spec/loop.md, REQ-012 — The bet; task T-3q0b,
// #350): the five headings of the brief each once, the line limit, at least
// one link and every link resolved, the fixed appetite and
// architecture-change lines, and the verified-batch lines; and the appetite
// of the bets against the band.

// goodBrief is a brief within a limit of 20 lines: the five headings, the two
// fixed lines, one link and one verified batch.
var goodBrief = `## The problem
The milestone's requirements.

## The appetite
appetite: USD 12.50, 40 hours

## The solution
The architecture stands.

## The rabbit holes and the no-gos
architecture change: no

## The facts
1. **need:** a fact.
See [the loop](docs/spec/loop.md).
approve batch gate: 0123456789012345678901234567890123456789 0123456789012345678901234567890123456789012345678901234567890123
`

var goodLinks = map[string]bool{"docs/spec/loop.md": true}

var goodVerified = []VerifiedBatch{{"gate", "0123456789012345678901234567890123456789",
	"0123456789012345678901234567890123456789012345678901234567890123"}}

// checked runs CheckBrief on the good brief with the given changes.
func checked(t *testing.T, brief string, lines *int, links map[string]bool, verified []VerifiedBatch) (Brief, []Finding) {
	t.Helper()
	return CheckBrief([]byte(brief), lines, links, verified)
}

func limit(n int) *int { return &n }

func TestABriefThatFitsPasses(t *testing.T) {
	b, f := checked(t, goodBrief, limit(20), goodLinks, goodVerified)
	if len(f) != 0 {
		t.Fatalf("findings %v, want none", f)
	}
	if b.USD != "12.50" || b.Hours != "40" || b.ArchChange || len(b.Batches) != 1 || b.Batches[0].Point != "gate" {
		t.Errorf("the parsed brief %+v", b)
	}
}

func TestEachHeadingIsRequiredOnce(t *testing.T) {
	for _, h := range []string{"## The problem", "## The appetite", "## The solution", "## The rabbit holes and the no-gos", "## The facts"} {
		missing := strings.Replace(goodBrief, h, h+"s", 1)
		_, f := checked(t, missing, limit(20), goodLinks, goodVerified)
		if !hasFinding(f, "heading", h+`" is missing`) || !hasFinding(f, "heading", "is not one of the five") {
			t.Errorf("the heading %s missing: findings %v", h, f)
		}
		twice := strings.Replace(goodBrief, h, h+"\n"+h, 1)
		if _, f := checked(t, twice, limit(21), goodLinks, goodVerified); !hasFinding(f, "heading", "is there 2 times") {
			t.Errorf("the heading %s twice: findings %v", h, f)
		}
	}
}

func TestTheLineLimit(t *testing.T) {
	if _, f := checked(t, goodBrief, nil, goodLinks, goodVerified); !hasFinding(f, "limit", "brief.lines is not set") {
		t.Errorf("no limit: findings %v", f)
	}
	if _, f := checked(t, goodBrief, limit(16), goodLinks, goodVerified); len(f) != 0 {
		t.Errorf("at the limit: findings %v", f)
	}
	if _, f := checked(t, goodBrief, limit(15), goodLinks, goodVerified); !hasFinding(f, "limit", "16") || !hasFinding(f, "limit", "15") {
		t.Errorf("over the limit names both numbers: findings %v", f)
	}
	// A last line without a line feed counts too.
	if _, f := checked(t, strings.TrimSuffix(goodBrief, "\n"), limit(15), goodLinks, goodVerified); !hasFinding(f, "limit", "16") {
		t.Errorf("the unterminated last line: findings %v", f)
	}
}

func TestTheLinks(t *testing.T) {
	for _, c := range []struct {
		name  string
		brief string
		links map[string]bool
		want  string
	}{
		{"no link at all", strings.Replace(goodBrief, "See [the loop](docs/spec/loop.md).", "See the loop.", 1), goodLinks, "link"},
		{"a link absent from the map", goodBrief, map[string]bool{}, "docs/spec/loop.md"},
		{"a link whose target does not resolve", goodBrief, map[string]bool{"docs/spec/loop.md": false}, "docs/spec/loop.md"},
		{"a reference link, an unsupported form", strings.Replace(goodBrief, "[the loop](docs/spec/loop.md)", "[the loop][1]", 1), goodLinks, "not of the inline form"},
	} {
		if _, f := checked(t, c.brief, limit(20), c.links, goodVerified); !hasFinding(f, "link", c.want) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
	if _, f := checked(t, strings.Replace(goodBrief, "See [the loop](docs/spec/loop.md).", "See [a](docs/a.md) and [b](docs/b.md).", 1),
		limit(20), map[string]bool{"docs/a.md": true, "docs/b.md": true}, goodVerified); len(f) != 0 {
		t.Errorf("two resolved links: findings %v", f)
	}
}

func TestTheFixedLines(t *testing.T) {
	for _, c := range []struct{ name, brief, rule, detail string }{
		{"no appetite line", strings.Replace(goodBrief, "appetite: USD 12.50, 40 hours\n", "", 1), "appetite", ""},
		{"two appetite lines", strings.Replace(goodBrief, "appetite: USD 12.50, 40 hours", "appetite: USD 12.50, 40 hours\nappetite: USD 1.0, 2 hours", 1), "appetite", ""},
		{"a bad appetite number", strings.Replace(goodBrief, "USD 12.50", "USD x", 1), "appetite", "is not `appetite: USD"},
		{"no architecture line", strings.Replace(goodBrief, "architecture change: no\n", "", 1), "architecture", ""},
		{"two architecture lines", strings.Replace(goodBrief, "architecture change: no", "architecture change: no\narchitecture change: yes", 1), "architecture", ""},
		{"a bad architecture value", strings.Replace(goodBrief, "architecture change: no", "architecture change: maybe", 1), "architecture", "is not `architecture change:"},
		{"the appetite line in the wrong section", strings.Replace(goodBrief, "## The appetite\nappetite: USD 12.50, 40 hours", "## The appetite\n\n## The solution\nappetite: USD 12.50, 40 hours\n\n## The appetite", 1), "appetite", ""},
		{"the architecture line in the wrong section", strings.Replace(goodBrief, "architecture change: no", "## The solution\narchitecture change: no\n\n## The rabbit holes and the no-gos", 1), "architecture", ""},
	} {
		if _, f := checked(t, c.brief, limit(20), goodLinks, goodVerified); !hasFinding(f, c.rule, c.detail) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
	if b, f := checked(t, strings.Replace(goodBrief, "architecture change: no", "architecture change: yes", 1), limit(20), goodLinks, goodVerified); len(f) != 0 || !b.ArchChange {
		t.Errorf("architecture change: yes: %+v, findings %v", b, f)
	}
}

func TestTheApproveBatchLines(t *testing.T) {
	for _, c := range []struct{ name, brief string }{
		{"a batch the caller did not verify", strings.Replace(goodBrief, "approve batch gate:", "approve batch other:", 1)},
		{"a batch whose head differs", strings.Replace(goodBrief, "approve batch gate: 0123", "approve batch gate: 1123", 1)},
		{"a batch whose rule hash differs", goodBrief + "approve batch x: " + strings.Repeat("1", 40) + " " + strings.Repeat("2", 64)},
		{"a batch line of a bad form", strings.Replace(goodBrief, "approve batch gate: 0123456789012345678901234567890123456789", "approve batch gate: sha", 1)},
	} {
		if _, f := checked(t, c.brief, limit(21), goodLinks, goodVerified); !hasFinding(f, "batch", "") {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
}

func TestTheAppetiteAgainstTheBand(t *testing.T) {
	for _, c := range []struct {
		name                string
		prior               []Cap
		proposed            Cap
		bandU, hours, spend string
		refused             bool
		detail              string
	}{
		{"two caps of 0.6 U", []Cap{{"0.6", "10"}}, Cap{"0.6", "10"}, "1.0", "40.0", "0", true, "the money caps of the bets are"},
		{"two caps that fit, each alone and together", []Cap{{"0.4", "10"}}, Cap{"0.5", "20"}, "1.0", "40.0", "0", false, ""},
		{"two hour caps whose sum exceeds", []Cap{{"25", "10"}}, Cap{"20", "35"}, "100.0", "40.0", "0", true, "the hour caps"},
		{"the spend of milestone 0 over the band", []Cap{{"0.4", "10"}}, Cap{"0.5", "20"}, "1.0", "40.0", "0.2", true, "plus the milestone-0 spend"},
		{"the money equal to the band", []Cap{{"0.4", "10"}}, Cap{"0.6", "20"}, "1.0", "40.0", "0", false, ""},
		{"the hours equal to the appetite", []Cap{{"0.4", "10"}}, Cap{"0.5", "30"}, "1.0", "40.0", "0", false, ""},
		{"the money with the spend equal to the band", []Cap{{"0.4", "10"}}, Cap{"0.5", "20"}, "1.0", "40.0", "0.1", false, ""},
		{"a cap below band.B is no minimum", nil, Cap{"0.01", "1"}, "100.0", "40.0", "0", false, ""},
		{"no proposed money", nil, Cap{"", "1"}, "100.0", "40.0", "0", true, ""},
		{"a proposed money that is no number", nil, Cap{"x", "1"}, "100.0", "40.0", "0", true, ""},
		{"no proposed hours", nil, Cap{"0.5", ""}, "100.0", "40.0", "0", true, ""},
		{"a band that is no number", nil, Cap{"0.5", "1"}, "x", "40.0", "0", true, ""},
		{"no hour appetite", nil, Cap{"0.5", "1"}, "100.0", "", "0", true, ""},
		{"no milestone-0 spend", nil, Cap{"0.5", "1"}, "100.0", "40.0", "", true, ""},
		{"exact decimal arithmetic: 0.1 + 0.2", []Cap{{"0.1", "0"}}, Cap{"0.2", "1"}, "0.3", "40.0", "0", false, ""},
	} {
		f := CheckAppetite(c.prior, c.proposed, c.bandU, c.hours, c.spend)
		if (len(f) != 0) != c.refused || (c.detail != "" && !hasFinding(f, "money", c.detail) && !hasFinding(f, "hours", c.detail)) {
			t.Errorf("%s: findings %v, refused %v", c.name, f, c.refused)
		}
	}
}

func hasFinding(fs []Finding, rule, detail string) bool {
	for _, f := range fs {
		if (rule == "" || f.Rule == rule) && strings.Contains(f.Detail, detail) {
			return true
		}
	}
	return false
}

// The answer tests (the same contract).

var goodAnswer = `REQ-001: Must, phase 2
NFR-001: Won't, phase —
fact 3: out of scope
fact Q-002: constraint
option stack: go
bet`

var (
	wantIDs    = []string{"REQ-001", "NFR-001"}
	wantPhases = []string{"1", "2", "alpha"}
	wantFacts  = []string{"3", "Q-002"}
	wantSurviv = map[string][]string{"stack": {"go", "rust"}}
)

func answered(t *testing.T, body string, ids, phases, facts []string, surviv map[string][]string) (Answer, []Finding) {
	t.Helper()
	return ParseAnswer([]byte(body), ids, phases, facts, surviv)
}

func TestAnAnswerThatFollowsTheFormParses(t *testing.T) {
	a, f := answered(t, goodAnswer, wantIDs, wantPhases, wantFacts, wantSurviv)
	if len(f) != 0 {
		t.Fatalf("findings %v, want none", f)
	}
	if a.Reqs["REQ-001"] != (ReqAnswer{"Must", "2"}) || a.Reqs["NFR-001"] != (ReqAnswer{"Won't", ""}) ||
		!a.Facts["3"].Out || a.Facts["Q-002"].Class != "constraint" ||
		a.Options["stack"] != "go" || a.Verdict != "bet" {
		t.Errorf("the parsed answer %+v", a)
	}
}

func TestTheRequirementLines(t *testing.T) {
	for _, c := range []struct{ name, body, detail string }{
		{"a requirement with no line", "REQ-001: Must, phase 2\nfact 3: out of scope\nbet", ""},
		{"a requirement the brief did not name", strings.Replace(goodAnswer, "\nbet", "\nREQ-002: Must, phase 2\nbet", 1), "which the brief did not name"},
		{"a requirement twice", strings.Replace(goodAnswer, "NFR-001:", "REQ-001:", 1), "has a line already"},
		{"a class that is none of the four", strings.Replace(goodAnswer, "Must,", "Shall,", 1), "of no form"},
		{"a phase — on a requirement that is not Won't", strings.Replace(goodAnswer, "phase 2", "phase —", 1), "for Won't alone"},
		{"a Won't with a phase", strings.Replace(goodAnswer, "Won't, phase —", "Won't, phase 2", 1), "a Won't has `phase —`"},
		{"a phase that is no tag of the target", strings.Replace(goodAnswer, "phase 2", "phase 9", 1), "no tag of the target"},
		{"a line with no phase", strings.Replace(goodAnswer, ", phase 2", "", 1), "of no form"},
	} {
		if _, f := answered(t, c.body, wantIDs, wantPhases, wantFacts, wantSurviv); len(f) == 0 || (c.detail != "" && !hasFinding(f, "line", c.detail)) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
	// A phase tag of the project that is no number is valid.
	if _, f := answered(t, strings.Replace(goodAnswer, "phase 2", "phase alpha", 1), wantIDs, wantPhases, wantFacts, wantSurviv); len(f) != 0 {
		t.Errorf("the phase tag alpha: findings %v", f)
	}
}

func TestTheFactLines(t *testing.T) {
	for _, c := range []struct{ name, body, detail string }{
		{"a fact the records do not know", strings.Replace(goodAnswer, "\nbet", "\nfact 9: out of scope\nbet", 1), "no fact the records know"},
		{"a fact twice", strings.Replace(goodAnswer, "\nbet", "\nfact 3: constraint\nbet", 1), "has a line already"},
		{"a class that is none of the four", strings.Replace(goodAnswer, "Q-002: constraint", "Q-002: context-ish", 1), "a class of the four"},
	} {
		if _, f := answered(t, c.body, wantIDs, wantPhases, wantFacts, wantSurviv); len(f) == 0 || (c.detail != "" && !hasFinding(f, "line", c.detail)) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
	for _, c := range []struct{ name, body string }{
		{"a reclassed fact", strings.Replace(goodAnswer, "3: out of scope", "3: measure", 1)},
		{"an answer fact out of scope", strings.Replace(goodAnswer, "Q-002: constraint", "Q-002: out of scope", 1)},
	} {
		if _, f := answered(t, c.body, wantIDs, wantPhases, wantFacts, wantSurviv); len(f) != 0 {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
}

func TestTheOptionLinesAndTheVerdict(t *testing.T) {
	for _, c := range []struct{ name, body, detail string }{
		{"a bet with no option line for a decision of two survivors", strings.Replace(goodAnswer, "option stack: go\n", "", 1), "no option line"},
		{"an option that is no survivor", strings.Replace(goodAnswer, "stack: go", "stack: python", 1), ""},
		{"a decision the Shape left no choice on", strings.Replace(goodAnswer, "stack: go", "ui: go", 1), ""},
		{"the same decision twice", goodAnswer + "\noption stack: rust", ""},
		{"no verdict", strings.Replace(goodAnswer, "\nbet", "", 1), "not `bet` or `no bet`"},
		{"a line after the verdict", goodAnswer + "\nREQ-002: Could, phase 1", "after the verdict"},
		{"a malformed line", "REQ-001: Must, phase 2\nNFR-001: Won't, phase —\nfact 3: out of scope\nfact Q-002: constraint\noption stack: go\nhello world\nbet", ""},
	} {
		if _, f := answered(t, c.body, wantIDs, wantPhases, wantFacts, wantSurviv); len(f) == 0 || (c.detail != "" && !hasFinding(f, "", c.detail)) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
	// A no-bet needs the requirement lines but no option line; an option line
	// that comes still names a survivor.
	a, f := answered(t, "REQ-001: Could, phase 1\nNFR-001: Won't, phase —\nno bet", wantIDs, wantPhases, nil, nil)
	if len(f) != 0 || a.Verdict != "no-bet" {
		t.Errorf("no bet: verdict %q, findings %v", a.Verdict, f)
	}
	// A no-bet asks no option of a decision of two survivors either.
	if _, f = answered(t, "REQ-001: Could, phase 1\nNFR-001: Won't, phase —\nno bet", wantIDs, wantPhases, nil, wantSurviv); len(f) != 0 {
		t.Errorf("no bet with survivors: findings %v", f)
	}
}
