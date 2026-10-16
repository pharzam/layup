package spec

// The checks of the bet (docs/spec/loop.md, REQ-012 — The bet; task T-3q0b,
// #350): the brief and the appetite, pure functions that the bet step of
// internal/run calls (row 66j): the caller resolves the links and reads the
// records, so no interface takes a forge client, a path or a callback.

import (
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strings"
)

// Finding is one broken rule of a check of this package: the rule's name and
// the detail that names the place (a heading, a line, a token or a number),
// which the session gets in a new attempt.
type Finding struct{ Rule, Detail string }

// briefHeadings are the five headings of the brief's five parts, each
// required once (loop.md: "Code checks the five headings"). The content rule
// of part 5 binds at the first bet (architecture.md §7 step 7); the bet's
// kind is the caller's, not the check's.
var briefHeadings = []string{
	"## The problem",
	"## The appetite",
	"## The solution",
	"## The rabbit holes and the no-gos",
	"## The facts",
}

var (
	// inlineLink is a link of the form [text](target); refLink is the
	// reference form [text][ref], which the check does not support.
	inlineLink = regexp.MustCompile(`\[[^\]]*\]\(([^)]*)\)`)
	refLink    = regexp.MustCompile(`\[[^\]]*\]\[[^\]]*\]`)

	appetiteLine = regexp.MustCompile(`^appetite: USD ([0-9]+(\.[0-9]+)?), ([0-9]+(\.[0-9]+)?) hours$`)
	archLine     = regexp.MustCompile(`^architecture change: (yes|no)$`)
	batchLine    = regexp.MustCompile(`^approve batch ([^:]+): ([0-9a-f]{40}) ([0-9a-f]{64})$`)
	batchStart   = regexp.MustCompile(`^approve batch `)
)

// VerifiedBatch is a rule batch that the brief asks to approve: its point,
// the Git head SHA and the SHA-256 of its rule file, as the caller verified
// them. The check compares, never approves (row 66j owns that).
type VerifiedBatch struct{ Point, Head, Rule string }

// Brief is the parsed bet brief: the numbers of its appetite line, whether it
// changes the architecture, and the batches its approve lines name.
type Brief struct {
	USD, Hours string
	ArchChange bool
	Batches    []VerifiedBatch
}

// CheckBrief checks a bet brief: the five headings each once, the line limit
// (the parameter brief.lines, whose unset is "brief.lines is not set"), at
// least one link with each target resolved by the caller's map, the fixed
// appetite and architecture-change lines, and each approve-batch line against
// the verified batches. A last line without a line feed counts too.
func CheckBrief(brief []byte, lines *int, links map[string]bool, verified []VerifiedBatch) (Brief, []Finding) {
	var b Brief
	var findings []Finding
	text := string(brief)
	n := strings.Count(text, "\n")
	if text != "" && !strings.HasSuffix(text, "\n") {
		n++
	}
	switch {
	case lines == nil:
		findings = append(findings, Finding{"limit", "brief.lines is not set"})
	case n > *lines:
		findings = append(findings, Finding{"limit", fmt.Sprintf("the brief is %d lines, over the limit %d", n, *lines)})
	}

	section := ""
	count := map[string]int{}
	var appetite, arch []string
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(l, "## ") {
			section = l
			if !slices.Contains(briefHeadings, l) {
				findings = append(findings, Finding{"heading", fmt.Sprintf("the heading %q is not one of the five", l)})
			}
			count[l]++
		}
		if appetiteLine.MatchString(l) && section == "## The appetite" {
			appetite = append(appetite, l)
		}
		if archLine.MatchString(l) && section == "## The rabbit holes and the no-gos" {
			arch = append(arch, l)
		}
		// A malformed fixed line fails closed: a line that starts like one
		// but does not match is a finding of that line's rule.
		switch {
		case strings.HasPrefix(l, "appetite:") && !appetiteLine.MatchString(l):
			findings = append(findings, Finding{"appetite", fmt.Sprintf("the line %q is not `appetite: USD <n>, <n> hours`", l)})
		case strings.HasPrefix(l, "architecture change:") && !archLine.MatchString(l):
			findings = append(findings, Finding{"architecture", fmt.Sprintf("the line %q is not `architecture change: yes|no`", l)})
		case batchStart.MatchString(l) && !batchLine.MatchString(l):
			findings = append(findings, Finding{"batch", fmt.Sprintf("the line %q is not `approve batch <point>: <head-sha> <rule-sha256>`", l)})
		case batchLine.MatchString(l):
			m := batchLine.FindStringSubmatch(l)
			vb := VerifiedBatch{m[1], m[2], m[3]}
			if !slices.Contains(verified, vb) {
				findings = append(findings, Finding{"batch", fmt.Sprintf("the batch %q is not one the caller verified", m[1])})
			}
			b.Batches = append(b.Batches, vb)
		}
	}
	for _, h := range briefHeadings {
		switch count[h] {
		case 1:
		case 0:
			findings = append(findings, Finding{"heading", fmt.Sprintf("the heading %q is missing", h)})
		default:
			findings = append(findings, Finding{"heading", fmt.Sprintf("the heading %q is there %d times", h, count[h])})
		}
	}

	found := 0
	for _, m := range inlineLink.FindAllStringSubmatch(text, -1) {
		found++
		if !links[m[1]] {
			findings = append(findings, Finding{"link", fmt.Sprintf("the link to %s does not resolve", m[1])})
		}
	}
	if found == 0 {
		findings = append(findings, Finding{"link", "the brief holds no link"})
	}
	for _, m := range refLink.FindAllStringSubmatch(text, -1) {
		findings = append(findings, Finding{"link", fmt.Sprintf("the link %s is not of the inline form [text](target)", m[0])})
	}

	one := func(rule, what string, xs []string) string {
		if len(xs) != 1 {
			findings = append(findings, Finding{rule, fmt.Sprintf("the %s is there %d times, wants one", what, len(xs))})
			return ""
		}
		return xs[0]
	}
	if l := one("appetite", "line `appetite: USD <n>, <n> hours`", appetite); l != "" {
		m := appetiteLine.FindStringSubmatch(l)
		b.USD, b.Hours = m[1], m[3]
	}
	if l := one("architecture", "line `architecture change: yes|no`", arch); l != "" {
		b.ArchChange = strings.HasSuffix(l, "yes")
	}
	return b, findings
}

// Cap is a bet's cap of money and of wall-clock, as decimal text.
type Cap struct{ USD, Hours string }

// CheckAppetite checks a proposed cap against the band and the appetite
// (loop.md, REQ-012 — The bet: the money caps of the bets so far, this one
// included, at most band.U; the wall-clock caps at most appetite.hours; the
// caps plus the spend of milestone 0 within band.U; a cap below band.B is
// allowed). Each amount is exact decimal; an absent or unusable one is a
// finding, never a zero.
func CheckAppetite(prior []Cap, proposed Cap, bandU, hours, milestoneZeroUSD string) []Finding {
	var findings []Finding
	amount := func(what, s string) *big.Rat {
		if s == "" {
			findings = append(findings, Finding{"amount", "the " + what + " is not set"})
			return nil
		}
		r, ok := new(big.Rat).SetString(s)
		if !ok || r.Sign() < 0 {
			findings = append(findings, Finding{"amount", fmt.Sprintf("the %s %q is no usable amount", what, s)})
			return nil
		}
		return r
	}
	money := amount("proposed money cap", proposed.USD)
	clock := amount("proposed hour cap", proposed.Hours)
	u := amount("band.U", bandU)
	h := amount("appetite.hours", hours)
	s0 := amount("milestone-0 spend", milestoneZeroUSD)
	if findings != nil {
		return findings
	}
	m, hh := new(big.Rat), new(big.Rat)
	for _, c := range prior {
		if a := amount("a prior bet's money cap", c.USD); a != nil {
			m.Add(m, a)
		}
		if a := amount("a prior bet's hour cap", c.Hours); a != nil {
			hh.Add(hh, a)
		}
	}
	m.Add(m, money)
	hh.Add(hh, clock)
	switch {
	case m.Cmp(u) > 0:
		findings = append(findings, Finding{"money", fmt.Sprintf("the money caps of the bets are %s, over band.U %s", m.RatString(), u.RatString())})
	case hh.Cmp(h) > 0:
		findings = append(findings, Finding{"hours", fmt.Sprintf("the hour caps of the bets are %s, over the appetite %s", hh.RatString(), h.RatString())})
	case new(big.Rat).Add(m, s0).Cmp(u) > 0:
		findings = append(findings, Finding{"money", fmt.Sprintf("the money caps %s plus the milestone-0 spend %s are over band.U %s", m.RatString(), s0.RatString(), u.RatString())})
	}
	return findings
}
