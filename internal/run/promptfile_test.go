package run

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/session"
)

func sumOf(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }

// tree is a target at its base in memory: the files by path.
func memTree(files map[string]string) targetTree {
	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	return targetTree{Paths: paths, Read: func(p string) ([]byte, error) {
		if s, ok := files[p]; ok {
			return []byte(s), nil
		}
		return nil, fmt.Errorf("no file %s", p)
	}}
}

const (
	prdText  = "# PRD\n\n| ID | Text |\n| -- | ---- |\n| REQ-005 | A handoff is a record. |\n| REQ-0050 | Another. |\nREQ-005 in prose.\n"
	specText = "# Spec\n\n## REQ-005 — The task loop\n\nSee ADR-0019 and ADR-0020.\n\n### Inside\n\nADR-0021 here.\n\n## REQ-007 — Statuses\n\nSee ADR-0016.\n"
)

func baseSrc(t *testing.T, step string) promptSrc {
	t.Helper()
	var row session.StepRow
	for _, r := range session.StepTable() {
		if r.Step == step {
			row = r
		}
	}
	return promptSrc{
		Step: row, Task: "T-ab12", Register: kept{"tasks.tsv", []byte("task\tissue\n")},
		Plan:         &kept{"copies/701-1.md", []byte("the plan\n")},
		Handoffs:     []kept{{"copies/702-1.md", []byte("the review\n")}},
		Answers:      []kept{{"tasks/T-ab12/results/S-bbbbbbbb.tsv", []byte("answer\n")}},
		Requirements: []string{"REQ-005"},
		Rules:        []rules.Entry{{Pattern: "AGENTS.md"}, {Pattern: "docs/rules/"}},
		Tree: memTree(map[string]string{
			"docs/prd/PRD-0001.md":        prdText,
			"docs/spec/loop.md":           specText,
			"docs/adr/0019-phase-loop.md": "adr 19\n",
			"docs/adr/0020-route.md":      "adr 20\n",
			"docs/adr/0021-smart-if.md":   "adr 21\n",
			"docs/adr/0016-gates.md":      "adr 16\n",
			"AGENTS.md":                   "rules\n",
			"docs/rules/a.md":             "rule a\n",
			"docs/other.md":               "other\n",
		}),
		Findings: []Finding{{Rule: "gate:static", Text: "command failed"}},
	}
}

func itemsOf(c *PromptContext, kind session.Kind) []string {
	var out []string
	for _, it := range c.Items {
		if it.Kind == kind {
			out = append(out, it.Path)
		}
	}
	return out
}

// loop.md, The tier and the context of a task: the prompt file holds, in this
// order, the issue and plan, the requirements, the decision records, the
// handoffs and answers, and the rule files.
func TestThePromptFileListsTheFiveKindsInTheirOrder(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	var kinds []session.Kind
	for _, it := range c.Items {
		if len(kinds) == 0 || kinds[len(kinds)-1] != it.Kind {
			kinds = append(kinds, it.Kind)
		}
	}
	if !slices.Equal(kinds, session.ContextKinds()) {
		t.Errorf("the kinds of the items in order: %v, want %v", kinds, session.ContextKinds())
	}
	text := string(c.Text("/s/S-aaaaaaaa"))
	last := -1
	for _, k := range session.ContextKinds() {
		i := strings.Index(text, string(k))
		if i <= last {
			t.Errorf("the heading %s is at %d, not after the one before it (%d):\n%s", k, i, last, text)
		}
		last = i
	}
}

func TestEachItemHasItsPathAndTheSHA256OfItsBytes(t *testing.T) {
	src := baseSrc(t, "implement")
	c, err := buildPrompt(src)
	if err != nil {
		t.Fatal(err)
	}
	text := string(c.Text("/s/S-aaaaaaaa"))
	for _, it := range c.Items {
		var data []byte
		if it.Where == "records" {
			for _, k := range append(append([]kept{src.Register, *src.Plan}, src.Handoffs...), src.Answers...) {
				if k.Path == it.Path {
					data = k.Data
				}
			}
		} else {
			data, _ = src.Tree.Read(it.Path)
		}
		if data == nil || it.SHA != sumOf(data) {
			t.Errorf("the item %s %s has the sha256 %s, want that of its bytes", it.Where, it.Path, it.SHA)
		}
		if !strings.Contains(text, it.SHA) {
			t.Errorf("the text does not hold the sha256 of %s", it.Path)
		}
	}
	// A records item is read under input/ of the session, a target item by its path in the clone.
	if !strings.Contains(text, "/s/S-aaaaaaaa/input/701-1.md") || strings.Contains(text, "/s/S-aaaaaaaa/input/loop.md") {
		t.Errorf("the paths of the text:\n%s", text)
	}
	if string(c.Input["701-1.md"]) != "the plan\n" || c.Input["loop.md"] != nil {
		t.Errorf("the input files %v; want the records items by their file name, and no target file", c.Input)
	}
}

func TestTheIssueAndPlanArePartOneAndAMissingPlanIsAnEmptyPartAndNotAnError(t *testing.T) {
	src := baseSrc(t, "plan-of-a-task")
	src.Plan, src.Rules = nil, nil // no plan, and no rule file: that part has no item
	c, err := buildPrompt(src)
	if err != nil {
		t.Fatalf("a task with no plan: %v", err)
	}
	if got := itemsOf(c, session.KindIssueAndPlan); !slices.Equal(got, []string{"tasks.tsv"}) {
		t.Errorf("issue-and-plan = %v, want only the register", got)
	}
	if text := string(c.Text("/s/x")); !strings.Contains(text, "rule-files:\n- none\n") {
		t.Errorf("an empty part is not said to be empty:\n%s", text)
	}
	c2, _ := buildPrompt(baseSrc(t, "implement"))
	if got := itemsOf(c2, session.KindIssueAndPlan); !slices.Equal(got, []string{"tasks.tsv", "copies/701-1.md"}) {
		t.Errorf("issue-and-plan = %v, want the register then the plan's copy", got)
	}
}

// A requirement's row is a table row of the PRD, its section a heading of the
// spec; a neighbour ID is not the ID (REQ-0050 is not REQ-005), and a word in
// prose is neither.
func TestTheRequirementItemsAreTheFilesWhoseRowOrHeadingHoldTheID(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	if got := itemsOf(c, session.KindRequirements); !slices.Equal(got, []string{"docs/prd/PRD-0001.md", "docs/spec/loop.md"}) {
		t.Errorf("requirements = %v", got)
	}
	src := baseSrc(t, "implement")
	src.Requirements = []string{"REQ-0050"}
	c, _ = buildPrompt(src)
	if got := itemsOf(c, session.KindRequirements); !slices.Equal(got, []string{"docs/prd/PRD-0001.md"}) {
		t.Errorf("REQ-0050 = %v; want only the PRD, whose row holds it (the spec's REQ-005 is not it)", got)
	}
	// A longer ID that starts with the ID, and one that ends with it, are not the ID.
	src = baseSrc(t, "implement")
	src.Tree = memTree(map[string]string{"docs/prd/PRD-0001.md": "| REQ-0050 | a |\n| XREQ-005 | b |\n"})
	src.Rules = nil
	c, _ = buildPrompt(src)
	if got := itemsOf(c, session.KindRequirements); len(got) != 0 {
		t.Errorf("REQ-0050 and XREQ-005 hold REQ-005: %v", got)
	}
	// A file that holds the ID in prose only is no item.
	src = baseSrc(t, "implement")
	src.Tree = memTree(map[string]string{"docs/prd/PRD-0001.md": "REQ-005 in prose.\n", "docs/spec/loop.md": "# Spec\n\nREQ-005 in prose.\n"})
	src.Rules = nil
	c, _ = buildPrompt(src)
	if got := itemsOf(c, session.KindRequirements); len(got) != 0 {
		t.Errorf("prose only: %v", got)
	}
	// Two requirements, in the order of the task's row, each file once.
	src = baseSrc(t, "implement")
	src.Requirements = []string{"REQ-007", "REQ-005"}
	c, _ = buildPrompt(src)
	if got := itemsOf(c, session.KindRequirements); !slices.Equal(got, []string{"docs/spec/loop.md", "docs/prd/PRD-0001.md"}) {
		t.Errorf("two requirements = %v; want the files of REQ-007 first, each once", got)
	}
}

// The decision records are those that the section of the heading names: its
// own text down to the next heading of its level or above.
func TestADecisionRecordIsOneThatASectionNames(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"docs/adr/0019-phase-loop.md", "docs/adr/0020-route.md", "docs/adr/0021-smart-if.md"}
	if got := itemsOf(c, session.KindDecisionRecords); !slices.Equal(got, want) {
		t.Errorf("decision records = %v, want %v (ADR-0021 is in a sub-section; ADR-0016 is of REQ-007's section)", got, want)
	}
	src := baseSrc(t, "implement")
	src.Tree = memTree(map[string]string{"docs/spec/loop.md": specText})
	if _, err := buildPrompt(src); err == nil || !strings.Contains(err.Error(), "ADR-0019") {
		t.Errorf("a decision record that a section names and the target lacks: %v; want an error naming it", err)
	}
}

func TestTheRuleFilesAreTheFilesOfTheRulePathRegister(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	if got := itemsOf(c, session.KindRuleFiles); !slices.Equal(got, []string{"AGENTS.md", "docs/rules/a.md"}) {
		t.Errorf("rule files = %v; want the files that the register's two patterns match", got)
	}
}

func TestTheHandoffsAreTheCopiesOfValidHandoffsAndTheResultsThatHoldAnAnswer(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	if got := itemsOf(c, session.KindHandoffsAndAnswers); !slices.Equal(got, []string{"copies/702-1.md", "tasks/T-ab12/results/S-bbbbbbbb.tsv"}) {
		t.Errorf("handoffs and answers = %v", got)
	}
}

// loop.md, The attempts: the findings of the last attempt are in the prompt
// file, and its diff is not.
func TestTheFindingsAreInThePromptAndThePromptHoldsNoDiff(t *testing.T) {
	c, err := buildPrompt(baseSrc(t, "implement"))
	if err != nil {
		t.Fatal(err)
	}
	text := string(c.Text("/s/S-aaaaaaaa"))
	if !strings.Contains(text, "gate:static") || !strings.Contains(text, "command failed") {
		t.Errorf("the findings are not in the text:\n%s", text)
	}
	if strings.Contains(text, "diff --git") || strings.Contains(text, "\n+++ ") {
		t.Errorf("the text holds a diff:\n%s", text)
	}
	src := baseSrc(t, "implement")
	src.Findings = nil
	c, _ = buildPrompt(src)
	if strings.Contains(string(c.Text("/s/x")), "gate:") {
		t.Error("a prompt with no finding names one")
	}
}

func TestTwoRecordsOfOneFileNameAreRefused(t *testing.T) {
	src := baseSrc(t, "implement")
	src.Handoffs = append(src.Handoffs, kept{"other/701-1.md", []byte("not the plan\n")})
	if _, err := buildPrompt(src); err == nil || !strings.Contains(err.Error(), "701-1.md") {
		t.Errorf("two records named 701-1.md: %v; want an error naming the file", err)
	}
	// One record named twice is one item.
	src = baseSrc(t, "implement")
	src.Handoffs = append(src.Handoffs, *src.Plan)
	c, err := buildPrompt(src)
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, it := range c.Items {
		if it.Path == "copies/701-1.md" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("a record named twice is %d items, want 1", n)
	}
}

func TestAStepOfNoTaskHasNoPromptFile(t *testing.T) {
	src := baseSrc(t, "implement")
	src.Step = session.StepRow{Step: "bet-brief", Role: "product-owner", Tier: "reasoning"}
	if _, err := buildPrompt(src); err == nil {
		t.Error("a step that has no record kinds has a prompt file")
	}
}
