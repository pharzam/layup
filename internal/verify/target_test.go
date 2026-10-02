package verify

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pharzam/layup/internal/work"
)

// The tests of the checks facts, onboarding, glossary and guardrails in a
// target's form, and of the fact resolver (D2 to D9 of #89).

func sum(s string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(s))) }

const (
	briefPath  = "docs/facts/problem-statement-brief.md"
	visionPath = "docs/facts/architectural-vision-brief.md"
	setupPath  = "docs/facts/F-0001-setup-answers.md"
	markerPath = "docs/facts/F-0002-marker-answers.md"
	indexPath  = "docs/facts/README.md"
	sumsPath   = "docs/setup/facts.sha256"
)

// answersRecord gives an answers record in the form of D2, one fact per
// question ID.
func answersRecord(id string, qs ...string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s. The answers\n\n## Facts as collected\n\n", id)
	for i, q := range qs {
		fmt.Fprintf(&b, "%d. `%s` an answer \u2014 by operator; source https://x.invalid/1; the question: a question\n", i+1, q)
	}
	return b.String()
}

// index gives docs/facts/README.md with a row for each ID.
func index(ids ...string) string {
	s := "# Customer facts\n\n## Index\n\n| Fact doc | Source | Collected | Status |\n| -------- | ------ | --------- | ------ |\n"
	for _, id := range ids {
		s += "| [" + id + "](" + id + "-x.md) | the answers | 2026-10-02 | Raw |\n"
	}
	return s
}

// factsTree gives a target's tree after S04 and S06: the two briefs, the
// record of S04, the index and docs/setup/facts.sha256 with a line for each
// raw facts file. change adds or replaces a file, and "" removes it; a change
// of facts.sha256 is used as it is.
func factsTree(change map[string]string) fstest.MapFS {
	files := map[string]string{briefPath: "# The problem\n\nThe statement.\n", visionPath: "# The vision\n",
		setupPath: answersRecord("F-0001", "S01-stack", "Q-001"), indexPath: index("F-0001")}
	for p, s := range change {
		files[p] = s
	}
	if _, ok := change[sumsPath]; !ok {
		var b strings.Builder
		for _, p := range []string{briefPath, visionPath, setupPath, markerPath} {
			if s := files[p]; s != "" {
				fmt.Fprintf(&b, "%s  %s\n", sum(s), p)
			}
		}
		files[sumsPath] = b.String()
	}
	fsys := fstest.MapFS{}
	for p, s := range files {
		if s != "" {
			fsys[p] = &fstest.MapFile{Data: []byte(s)}
		}
	}
	return fsys
}

// doneRows gives a setup record with the done row of each step.
func doneRecord(steps ...string) work.Record {
	var r work.Record
	for _, s := range steps {
		r = append(r, []string{s, "done", "ok", "step", ""})
	}
	return r
}

var (
	allSteps  = []string{"S01", "S02", "S03", "S04", "S05", "S06", "S07", "S08", "S09", "S10", "S11", "S12", "S13", "S14"}
	setupRows = work.Answers{{"S01-stack", "go", "operator", "u", ""}, {"Q-001", "ten", "idea-owner", "u", ""}}
	withM     = append(append(work.Answers{}, setupRows...), []string{"M-0123abcd", "8080", "operator", "u", ""})
)

// The hash loop of check_facts, with its lines: each line of the list names a
// file whose SHA-256 it holds (D6).
func TestTheHashesOfFacts(t *testing.T) {
	same(t, "no list", factsHashFindings(fstest.MapFS{}), []string{"hash: docs/setup/facts.sha256 is absent"})
	same(t, "a good tree", factsHashFindings(factsTree(nil)), nil)
	list := sum("a") + "  docs/facts/a.md\nabc\n\n" + strings.Repeat("1", 64) + "  docs/facts/x.md\n" +
		sum("d") + "  docs/facts/d\n" + sum("b") + "\tdocs/facts/b.md"
	fsys := fstest.MapFS{sumsPath: {Data: []byte(list)}, "docs/facts/a.md": {Data: []byte("a")},
		"docs/facts/b.md": {Data: []byte("b")}, "docs/facts/d/x": {Data: []byte("x")}}
	same(t, "the lines", factsHashFindings(fsys), []string{
		"hash: line 2 of docs/setup/facts.sha256 has no path",
		"hash: docs/facts/x.md does not match docs/setup/facts.sha256",
		"hash: docs/facts/d does not match docs/setup/facts.sha256",
	})
}

// Check facts in a target's form reads each record and brief that exists in
// full, and the done rows of the setup record decide what must exist
// (D4 with condition 2 of the plan review, D6, K18).
func TestTheFactsOfATarget(t *testing.T) {
	twoSums := sum("# The problem\n\nThe statement.\n") + "  " + briefPath + "\n" + sum("# The vision\n") + "  " + visionPath + "\n"
	for _, c := range []struct {
		name    string
		change  map[string]string
		answers work.Answers
		done    []string
		want    []string
	}{
		{"a setup through S14", nil, setupRows, allSteps, nil},
		{"a new tree: nothing is required", map[string]string{setupPath: "", briefPath: "", visionPath: "", indexPath: ""}, setupRows, nil, nil},
		{"S04 is done and its record is not there", map[string]string{setupPath: ""}, setupRows, []string{"S04"},
			[]string{"record: expected one docs/facts/F-NNNN-setup-answers.md"}},
		{"two records of S04", map[string]string{"docs/facts/F-0003-setup-answers.md": answersRecord("F-0003", "S01-stack", "Q-001")}, setupRows, nil,
			[]string{"record: expected one docs/facts/F-NNNN-setup-answers.md"}},
		{"S06 is done and the problem statement is not there", map[string]string{briefPath: ""}, setupRows, []string{"S04", "S06"},
			[]string{"brief: docs/facts/problem-statement-brief.md is absent"}},
		{"no vision brief is no finding", map[string]string{visionPath: ""}, setupRows, allSteps, nil},
		{"a record not in the list", map[string]string{sumsPath: twoSums}, setupRows, nil,
			[]string{"listed: docs/facts/F-0001-setup-answers.md is not in docs/setup/facts.sha256"}},
		{"a brief not in the list", map[string]string{sumsPath: sum("# The vision\n") + "  " + visionPath + "\n" + sum(answersRecord("F-0001", "S01-stack", "Q-001")) + "  " + setupPath + "\n"}, setupRows, nil,
			[]string{"listed: docs/facts/problem-statement-brief.md is not in docs/setup/facts.sha256"}},
		{"no index row", map[string]string{indexPath: index()}, setupRows, nil,
			[]string{"index: docs/facts/README.md has no row for F-0001"}},
		{"an answer row with no fact", nil, append(append(work.Answers{}, setupRows...), []string{"Q-002", "x", "idea-owner", "u", ""}), nil,
			[]string{"answers: Q-002 of inputs/answers.tsv is not a fact of F-0001"}},
		{"a fact with no answer row", map[string]string{setupPath: answersRecord("F-0001", "S01-stack", "Q-001", "Q-009")}, setupRows, nil,
			[]string{"answers: F-0001 fact 3 names Q-009, which inputs/answers.tsv does not have"}},
		{"a question twice", map[string]string{setupPath: answersRecord("F-0001", "S01-stack", "Q-001", "Q-001")}, setupRows, nil,
			[]string{"answers: F-0001 holds Q-001 2 times"}},
		{"an M- question in the record of S04", map[string]string{setupPath: answersRecord("F-0001", "S01-stack", "Q-001", "M-0123abcd")}, withM, nil,
			[]string{"answers: F-0001 fact 3 names M-0123abcd, a question of the record of S11"}},
		{"a fact with no question ID", map[string]string{setupPath: answersRecord("F-0001", "S01-stack", "Q-001") + "3. an answer with no ID\n"}, setupRows, nil,
			[]string{"answers: F-0001 fact 3 names no question ID"}},
		{"an M- row before S11 is done is no finding", nil, withM, []string{"S04", "S06"}, nil},
		{"S11 is done, an M- row, and no record of S11", nil, withM, allSteps,
			[]string{"record: expected one docs/facts/F-NNNN-marker-answers.md"}},
		{"S11 is done and no M- row: no record of S11 is needed", nil, setupRows, allSteps, nil},
		{"a record of S11 is read in full before S11 is done", map[string]string{markerPath: answersRecord("F-0002", "M-0123abcd"), indexPath: index("F-0001", "F-0002")}, withM, []string{"S04", "S06"}, nil},
		{"an S01- question in the record of S11", map[string]string{markerPath: answersRecord("F-0002", "M-0123abcd", "S01-stack"), indexPath: index("F-0001", "F-0002")}, withM, allSteps,
			[]string{"answers: F-0002 fact 2 names S01-stack, a question of the record of S04"}},
	} {
		in := input{fsys: factsTree(c.change), repo: oneRoot, record: doneRecord(c.done...), answers: c.answers}
		same(t, c.name, checkFacts(in), c.want)
	}
}

// F-NNNN#n resolves when docs/facts/ holds exactly one F-NNNN-*.md with a
// numbered fact n (D7).
func TestTheFactResolver(t *testing.T) {
	one := fstest.MapFS{"docs/facts/F-0001-a.md": {Data: []byte("# F-0001\n\n1. one\n03. three\n")}}
	two := fstest.MapFS{"docs/facts/F-0001-a.md": {Data: []byte("1. one\n")}, "docs/facts/F-0001-b.md": {Data: []byte("1. one\n")}}
	for _, c := range []struct {
		name string
		fsys fstest.MapFS
		cite string
		want string
	}{
		{"a fact", one, "F-0001#1", ""},
		{"a fact with a leading zero", one, "F-0001#3", ""},
		{"no such fact", one, "F-0001#2", "fact: F-0001#2 is not a fact of the F-0001 record"},
		{"no record", one, "F-0002#1", "fact: F-0002#1 is not a fact of the F-0002 record"},
		{"two records", two, "F-0001#1", "fact: F-0001#1 is not a fact of the F-0001 record"},
	} {
		if got := resolveFact(c.fsys, c.cite); got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
}

const onboardingText = "# Onboarding\n\nRead [the brief](facts/problem-statement-brief.md) (`F-0001#1`).\n"

// Check onboarding in a target's form: the rules of check_onboarding, with
// the resolver for each citation (D8).
func TestTheOnboardingOfATarget(t *testing.T) {
	facts := "# F-0001\n\n1. one\n"
	for _, c := range []struct {
		name, text string
		want       []string
	}{
		{"a good file", onboardingText, nil},
		{"no file", "", []string{"missing: docs/onboarding-for-engineers.md is absent"}},
		{"a marker", onboardingText + "\u2039x\u203a\n", []string{"marker: docs/onboarding-for-engineers.md holds a \u2039 character"}},
		{"no link", "# Onboarding\n", []string{"link: docs/onboarding-for-engineers.md has no link to facts/problem-statement-brief.md"}},
		{"citations that do not resolve", onboardingText + "See `F-0001#9` and `F-0004#1`.\n",
			[]string{"fact: F-0001#9 is not a fact of the F-0001 record", "fact: F-0004#1 is not a fact of the F-0004 record"}},
	} {
		fsys := fstest.MapFS{"docs/facts/F-0001-setup-answers.md": {Data: []byte(facts)}}
		if c.text != "" {
			fsys["docs/onboarding-for-engineers.md"] = &fstest.MapFile{Data: []byte(c.text)}
		}
		same(t, c.name, onboardingFindings(fsys), c.want)
	}
}

// Check glossary in a target's form: each citation resolves; any heading and
// any number of rows (D8).
func TestTheGlossaryOfATarget(t *testing.T) {
	facts := fstest.MapFS{"docs/facts/F-0001-setup-answers.md": {Data: []byte("1. one\n")}}
	same(t, "no file", glossaryFindings(facts), []string{"missing: docs/glossary.md is absent"})
	facts["docs/glossary.md"] = &fstest.MapFile{Data: []byte("# Glossary\n\n## Our terms\n\n| Term | Meaning |\n| - | - |\n| A | an a (`F-0001#1`) |\n")}
	same(t, "any heading, one row", glossaryFindings(facts), nil)
	facts["docs/glossary.md"] = &fstest.MapFile{Data: []byte("| B | a b (`F-0001#2`) |\n")}
	same(t, "a citation that does not resolve", glossaryFindings(facts), []string{"fact: F-0001#2 is not a fact of the F-0001 record"})
}

// Check guardrails in a target's form: an entry is a bullet - **Inv-N** up to
// the next entry, a heading or the end; its Check: value as check_guardrails
// reads it; each citation resolves (D8, note 4 of the plan review).
func TestTheGuardrailsOfATarget(t *testing.T) {
	tree := func(text string) fstest.MapFS {
		return fstest.MapFS{"docs/guardrails.md": {Data: []byte(text)}, "README.md": {Data: []byte("x")},
			"docs/x.sh": {Data: []byte("x")}, "docs/dir/a": {Data: []byte("a")},
			"docs/facts/F-0001-setup-answers.md": {Data: []byte("1. one\n")}}
	}
	same(t, "no file", guardrailsFindings(fstest.MapFS{}), []string{"missing: docs/guardrails.md is absent"})
	same(t, "no entry", guardrailsFindings(tree("# Guardrails\n\nCheck: maybe later\n")), nil)
	text := "# Guardrails\n\n" +
		"- **Inv-1** \u2014 a rule (`F-0001#1`). Check: no check yet  \n" +
		"- **Inv-2** \u2014 a rule. Check: README.md (hook)\n" +
		"- **Inv-3** \u2014 a rule.\n  Check: docs/x.sh (ci:build)\n" +
		"- **Inv-4** \u2014 a rule. Check: docs/nothing.sh (ci:x)\n" +
		"- **Inv-5** \u2014 a rule. Check: README.md (cron)\n" +
		"- **Inv-6** \u2014 a rule. Check: maybe later\n" +
		"- **Inv-7** \u2014 a rule. Check: docs/dir (hook)\n" +
		"- **Inv-8** \u2014 a rule with no check value\n\n## Notes\n\nCheck: no check yet\n" +
		"- **Inv-9** \u2014 a rule. Check: a Check: no check yet\n" +
		"- **Inv-10** \u2014 a rule (`F-0001#2`). Check: README.md (ci:)\n"
	same(t, "the entries", guardrailsFindings(tree(text)), []string{
		"check: Inv-4 names \"docs/nothing.sh\", which is not a file",
		"check: Inv-5 gate \"cron\" is not hook or ci:<job>",
		"check: Inv-6 has no valid Check: value (\"maybe later\")",
		"check: Inv-7 names \"docs/dir\", which is not a file",
		"check: Inv-8 has no valid Check: value (\"\")",
		"check: Inv-10 gate \"ci:\" is not hook or ci:<job>",
		"fact: F-0001#2 is not a fact of the F-0001 record",
	})
}
