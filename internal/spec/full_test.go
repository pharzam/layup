package spec

// The tests of the full spec check (task T-4wyg, #351): each of the ten rules,
// the before-first-bet mode and the whole-table clear, the record-bound fact
// tokens, and the waiting row of classes — each an assertion the plan review
// of #351 pins (conditions C4 to C13).

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// A small consistent world for the tests: the problem statement's facts are
// numbered 1 to 8 (the records F-0001 and F-0003 number the one file), the
// answers' record is F-0004, whose fact N answers Q-00N, and Q-001 is answered
// with a quote in the span of fact 7.
func testPSB() []byte {
	return []byte(strings.Repeat("x", 100))
}

// span: start, end, fact, class, reason, text — as spec.Read gives them.
func testSpans() [][]string {
	return [][]string{
		{"0", "10", "1", "need", "", "s1"},
		{"10", "20", "2", "constraint", "", "s2"},
		{"20", "30", "3", "measure", "", "s3"},
		{"30", "40", "4", "need", "", "s4"},
		{"40", "50", "5", "measure", "", "s5"},
		{"50", "60", "6", "context", "", "s6"},
		{"60", "70", "7", "need", "", "s7"}, // holds the quote of Q-001
		{"70", "80", "8", "need", "", "s8"}, // only a waits-for-a-bet line covers it
	}
}

func testQuestions() [][]string {
	return [][]string{
		{"Q-001", "rule", "g1", "", "", "", "60", "70", "the question", "2026-01-01T00:00:00Z", "layup psb check", "2026-01-01T01:00:00Z", "answered"},
	}
}

func testAnswers() [][]string {
	return [][]string{
		{"Q-001", "yes, it does", "11", "1", "2026-01-01T01:00:00Z"},
	}
}

// testPRD is the consistent PRD: every need and constraint of the spans is
// covered by a requirement row or by the waits-for-a-bet section, each measure
// has a success row, each requirement has a criterion, and every fact token
// resolves — F-0001#7 being the quote's fact of the answer Q-001 that REQ-002
// covers as F-0004#1.
const testPRD = `# PRD-0001. A stand-in

## 6. Functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| REQ-001 | Alpha and beta. | Must | 1 | F-0001#1, F-0001#2 |
| REQ-002 | The answer. | Should | 2 | F-0004#1, F-0001#7 |

## 7. Non-functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| NFR-001 | Delta. | Must | 1 | F-0001#4 |

### 7.1 Acceptance criteria

One criterion per requirement.

| Requirement | Acceptance criterion |
| ----------- | -------------------- |
| ` + "`REQ-001`" + ` | a check of it |
| ` + "`REQ-002`" + ` | a check of it |
| ` + "`NFR-001`" + ` | a check of it |

## 8. Success metrics

| Dimension | Start value | Fact |
| --------- | ----------- | ---- |
| one | x | F-0001#3 |
| two | y | F-0001#5 |

## waits for a bet

A requirement that waits for its next bet covers F-0001#8.
`

// testRegister is a tasks.tsv whose one task names every requirement and has
// merged.
func testRegister(t *testing.T) []byte {
	var b bytes.Buffer
	if err := tsv.Write(&b, tasksSchema, [][]string{
		{"T-aaaa", "1", "build", "REQ-001 REQ-002 NFR-001", "small", "", "1", "2026-01-01T02:00:00Z"},
	}); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// testInventory writes one version of the confirmed inventory and gives the
// index and the versions map, the pair the caller hands over.
func testInventory(t *testing.T, rows [][]string) (index [][]string, versions map[int][]byte) {
	var b bytes.Buffer
	if err := tsv.Write(&b, InventorySchema, rows); err != nil {
		t.Fatal(err)
	}
	return [][]string{{"1", fmt.Sprintf("%x", sha256.Sum256(b.Bytes())), "M2e"}}, map[int][]byte{1: b.Bytes()}
}

func inventoryOfSpans(t *testing.T, extra ...[]string) ([][]string, map[int][]byte) {
	var rows [][]string
	for _, s := range testSpans() {
		rows = append(rows, []string{s[colFact], s[colClass], "in", "9"})
	}
	rows = append(rows, []string{"Q-001", "context", "in", "9"})
	return testInventory(t, append(rows, extra...))
}

const testCopy = `REQ-001: Must, phase 1
REQ-002: Should, phase 2
NFR-001: Must, phase 1
fact 3: out of scope
option D-1: a
bet
`

// baseInput is the consistent post-bet world; tests mutate fields.
func baseInput(t *testing.T) FullInput {
	index, versions := inventoryOfSpans(t)
	return FullInput{
		Changed:       []string{"docs/spec/loop.md"},
		PSBRecords:    []string{"F-0001", "F-0003"},
		AnswersRecord: "F-0004",
		PRDs:          map[string][]byte{"PRD-0001-x.md": []byte(testPRD)},
		SpecFiles: map[string][]byte{"loop.md": []byte(`# s

## REQ-001 — One

text.

## REQ-002 — Two

text.

## NFR-001 — Three

text.
`)},
		PSB:       testPSB(),
		Spans:     testSpans(),
		Questions: testQuestions(),
		Answers:   testAnswers(),
		Index:     index,
		Versions:  versions,
		Copy:      []byte(testCopy),
		Register:  testRegister(t),
		Merged:    []string{"T-aaaa"},
	}
}

// preBetInput is the same world before the first bet: no index, no versions,
// no copy.
func preBetInput(t *testing.T) FullInput {
	in := baseInput(t)
	in.Index = nil
	in.Versions = nil
	in.Copy = nil
	return in
}

// The rows of CheckFull are rule, result, detail — as the block's columns
// after the head the writer adds — so ruleOf of the span check's tests
// applies: it gives (result, detail).
func wantResult(t *testing.T, rows [][]string, rule, result string) {
	t.Helper()
	if got, detail := ruleOf(t, rows, rule); got != result {
		t.Errorf("the rule %s: result %q, want %q (detail %q)", rule, got, result, detail)
	}
}

func wantFail(t *testing.T, rows [][]string, rule, detailPart string) {
	t.Helper()
	got, detail := ruleOf(t, rows, rule)
	if got != resFail || !strings.Contains(detail, detailPart) {
		t.Errorf("the rule %s: (%q, %q), want fail with %q in the detail", rule, got, detail, detailPart)
	}
}

// The table is the ten rules, each once, in the block's order. The want list
// is the block's literal order, not the code's own variable: a mutant that
// permutes fullRules must fail here.
func wantTenRows(t *testing.T, rows [][]string) {
	t.Helper()
	if len(rows) != 10 {
		t.Fatalf("the table has %d rows, want the ten rules", len(rows))
	}
	want := []string{"covers", "needs", "measures", "criteria", "classes",
		"inventory", "quotes", "sections", "tasks", "ran"}
	for i, rule := range want {
		if rows[i][0] != rule {
			t.Fatalf("row %d is %q, want the rule %s (the block's order)", i, rows[i][0], rule)
		}
	}
}

func TestCheckFullPassesOnAConsistentInput(t *testing.T) {
	rows := CheckFull(baseInput(t))
	wantTenRows(t, rows)
	for _, r := range rows {
		if r[1] != resPass {
			t.Errorf("the row %v: want pass on a consistent input", r)
		}
	}
}

func TestCheckFullBeforeTheFirstBet(t *testing.T) {
	rows := CheckFull(preBetInput(t))
	wantTenRows(t, rows)
	for _, rule := range []string{"classes", "inventory"} {
		got, detail := ruleOf(t, rows, rule)
		if got != resClear || detail != "before the first bet" {
			t.Errorf("the rule %s: (%q, %q), want clear with the reason `before the first bet`", rule, got, detail)
		}
	}
	for _, rule := range []string{"covers", "needs", "measures", "criteria", "quotes", "sections", "tasks", "ran"} {
		wantResult(t, rows, rule, resPass)
	}
}

// A docs/prd/ path takes the spec-touching mode as a docs/spec/ path does:
// the bet rules are clear and the numbered-facts rules judge (round-2 note
// N2 — the docs/prd/ half of touchesSpec).
func TestCheckFullBeforeTheFirstBetOnAPRDPath(t *testing.T) {
	in := preBetInput(t)
	in.Changed = []string{"docs/prd/PRD-0001-x.md"}
	rows := CheckFull(in)
	wantTenRows(t, rows)
	for _, rule := range []string{"classes", "inventory"} {
		got, detail := ruleOf(t, rows, rule)
		if got != resClear || detail != "before the first bet" {
			t.Errorf("the rule %s: (%q, %q), want clear with the reason `before the first bet`", rule, got, detail)
		}
	}
	for _, rule := range []string{"covers", "needs", "criteria", "tasks"} {
		wantResult(t, rows, rule, resPass)
	}
}

func TestCheckFullNoSpecificationYet(t *testing.T) {
	in := preBetInput(t)
	in.Changed = []string{"internal/spec/full.go"}
	rows := CheckFull(in)
	wantTenRows(t, rows)
	for _, r := range rows {
		if r[1] != resClear || r[2] != "no specification yet" {
			t.Errorf("the row %v: want clear with the reason `no specification yet`", r)
		}
	}
}

// After the first bet a head that touches neither the PRD nor docs/spec/ is
// judged on all ten rules — the whole-table clear is only before the first
// bet (review condition C4; loop.md scopes `no specification yet` to it).
func TestCheckFullJudgesAllTenPostBetWithoutSpecPath(t *testing.T) {
	in := baseInput(t)
	in.Changed = []string{"internal/spec/full.go"}
	rows := CheckFull(in)
	wantTenRows(t, rows)
	for _, r := range rows {
		if r[1] != resPass {
			t.Errorf("the row %v: a judged rule on a consistent input, want pass", r)
		}
	}
}

// covers: a fact token resolves only for the records that number the problem
// statement, or for the answers' record, whose fact N answers Q-00N (C7).
func TestCoversBindsEachTokenToItsRecord(t *testing.T) {
	prd := func(facts string) map[string][]byte {
		return map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
			"| REQ-001 | Alpha and beta. | Must | 1 | F-0001#1, F-0001#2 |",
			"| REQ-001 | Alpha and beta. | Must | 1 | "+facts+" |", 1))}
	}
	in := baseInput(t)
	in.PRDs = prd("F-0004#1, F-0001#2") // the answers' record: fact 1 is Q-001, answered
	rows := CheckFull(in)
	wantResult(t, rows, "covers", resPass)

	in.PRDs = prd("F-0004#9") // fact 9 of the answers is Q-009, never asked
	rows = CheckFull(in)
	wantFail(t, rows, "covers", "F-0004#9")

	in.PRDs = prd("F-0002#1") // a record that numbers no fact of the problem statement
	rows = CheckFull(in)
	wantFail(t, rows, "covers", "F-0002#1")

	in.PRDs = prd("Q-001") // the bare form is accepted
	rows = CheckFull(in)
	wantResult(t, rows, "covers", resPass)

	in.PRDs = prd("F-0001") // a bare record of the problem statement's facts
	rows = CheckFull(in)
	wantResult(t, rows, "covers", resPass)

	in.PRDs = prd("F-0004") // a bare cite of the answers' record resolves to no numbered fact
	rows = CheckFull(in)
	wantResult(t, rows, "covers", resPass)

	// A record that the head holds in docs/facts/ (the setup record, which prd-lint
	// reads) is accepted bare; its facts are no spans and no answers, and an unknown record is no record.
	in.Records = []string{"F-0005"}
	for _, c := range []struct {
		cite string
		pass bool
	}{{"F-0005", true}, {"F-0005#1", false}, {"F-0009", false}, {"F-0005, F-0001#1", true}} {
		in.PRDs = prd(c.cite)
		wantResult(t, CheckFull(in), "covers", map[bool]string{true: resPass, false: resFail}[c.pass])
	}
}

// needs: an in-scope need or constraint with no requirement and no waiting
// requirement is a fail; out-of-scope and waits-for-a-bet exempt (fact 8 of
// the base world is exempted by the waits section).
func TestNeedsJudgesTheFacts(t *testing.T) {
	in := baseInput(t)
	// Fact 4 loses its requirement.
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |",
		"| NFR-001 | Delta. | Must | 1 | F-0001#6 |", 1))}
	rows := CheckFull(in)
	wantFail(t, rows, "needs", "4")

	// A need marked out of scope in the inventory is exempt: the same PRD that
	// failed above passes once fact 4's scope is out.
	in = baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |",
		"| NFR-001 | Delta. | Must | 1 | F-0001#6 |", 1))}
	var inv [][]string
	for _, r := range testSpans() {
		scope := "in"
		if r[colFact] == "4" {
			scope = "out"
		}
		inv = append(inv, []string{r[colFact], r[colClass], scope, "9"})
	}
	inv = append(inv, []string{"Q-001", "context", "in", "9"})
	in.Index, in.Versions = testInventory(t, inv)
	rows = CheckFull(in)
	wantResult(t, rows, "needs", resPass)
}

// Fact 8, a need that only the waits-for-a-bet section names, is exempt.
func TestNeedsExemptsAWaitingRequirement(t *testing.T) {
	in := baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"covers F-0001#8.", "states no fact.", 1))}
	rows := CheckFull(in)
	wantFail(t, rows, "needs", "8")
}

// measures: only a success-criterion row's Fact cell covers a measure — a
// cite inside a requirement's acceptance cell does not (C8).
func TestMeasuresReadsOnlyTheSuccessRows(t *testing.T) {
	in := baseInput(t)
	// Move fact 3's cite from the success table into REQ-001's criterion.
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.NewReplacer(
		"| one | x | F-0001#3 |", "| one | x | F-0001#5 |",
		"`REQ-001` | a check of it |", "`REQ-001` | a check of it (F-0001#3) |",
	).Replace(testPRD))}
	rows := CheckFull(in)
	wantFail(t, rows, "measures", "3")
}

// measures' two exemptions (round-1 finding F3): a measure fact the confirmed
// version marks out of scope, and one a waits-for-a-bet section names, each
// pass with no success row — drop either guard and the rule fails.
func TestMeasuresExemptsAnOutOfScopeFact(t *testing.T) {
	in := baseInput(t)
	in.Index, in.Versions = inventoryOfSpans(t, []string{"9", "measure", "out", "9"})
	wantResult(t, CheckFull(in), "measures", resPass)
}

func TestMeasuresExemptsAWaitingFact(t *testing.T) {
	in := baseInput(t)
	in.Index, in.Versions = inventoryOfSpans(t, []string{"9", "measure", "in", "9"})
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"covers F-0001#8.", "covers F-0001#8 and F-0001#9.", 1))}
	wantResult(t, CheckFull(in), "measures", resPass)
}

// The waits-for-a-bet heading is the exact text, case-folded, at any heading
// level: a heading that only holds the phrase is prose (round-2 F1 — the
// stated contract needed these falsifiers).
func TestWaitsForABetIsTheExactFoldedHeading(t *testing.T) {
	prd := func(heading string) map[string][]byte {
		return map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
			"## waits for a bet", heading, 1))}
	}
	// The phrase inside a longer heading names no section: fact 8 loses its
	// exemption and needs fails, for trailing prose and for leading prose
	// alike — a suffix match would wrongly open the section.
	in := baseInput(t)
	in.PRDs = prd("## waits for a bet later")
	wantFail(t, CheckFull(in), "needs", "8")

	in = baseInput(t)
	in.PRDs = prd("## the team waits for a bet")
	wantFail(t, CheckFull(in), "needs", "8")

	// Any case and each of the six Markdown levels hold the section.
	for _, heading := range []string{
		"# waits for a bet", "## WAITS FOR A BET", "### Waits for a bet",
		"#### waits FOR a bet", "##### Waits For A Bet", "###### waits for a bet",
	} {
		in = baseInput(t)
		in.PRDs = prd(heading)
		wantResult(t, CheckFull(in), "needs", resPass)
	}
}

// The waits section ends at the next heading of any level: a cite in the
// prose after it exempts nothing.
func TestWaitsForABetEndsAtTheNextHeading(t *testing.T) {
	in := baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"## waits for a bet\n\nA requirement that waits for its next bet covers F-0001#8.\n",
		"## waits for a bet\n\nNothing waits.\n\n### Later notes\n\nA note covers F-0001#8.\n", 1))}
	wantFail(t, CheckFull(in), "needs", "8")
}

// The waits section names facts by the citation forms (C9): F-<record>#n
// (the base fixture's F-0001#8), F-<answers>#N, and the bare Q-NNN. A bare
// F-NNNN of a record that numbers nothing is a cite covers judges — it
// resolves to no fact of the universe.
func TestWaitsForABetNamesTheCitationForms(t *testing.T) {
	in := baseInput(t)
	in.Index, in.Versions = inventoryOfSpans(t,
		[]string{"Q-002", "need", "in", "9"}, []string{"Q-009", "need", "in", "9"})
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"covers F-0001#8.", "covers F-0001#8, F-0004#2 and Q-009.", 1))}
	wantResult(t, CheckFull(in), "needs", resPass)

	in = baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"covers F-0001#8.", "covers F-0001#8 and F-9999.", 1))}
	wantFail(t, CheckFull(in), "covers", "F-9999")
}

// The success table counts only under a `## … Success metrics` heading: the
// same table at another level covers no measure (round-2 F1).
func TestSuccessMetricsHeadingIsLevelTwo(t *testing.T) {
	for _, heading := range []string{"### 8. Success metrics", "# 8. Success metrics"} {
		in := baseInput(t)
		in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
			"## 8. Success metrics", heading, 1))}
		wantFail(t, CheckFull(in), "measures", "3")
	}
}

func TestCriteriaNeedsOneRowPerRequirement(t *testing.T) {
	in := baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| `NFR-001` | a check of it |\n", "", 1))}
	rows := CheckFull(in)
	wantFail(t, rows, "criteria", "NFR-001")
}

// A malformed row of the criteria table — one cell, or none — is no panic and
// no criterion: the requirement it truncates keeps no row and the rule fails
// (round-1 finding F1).
func TestAShortCriteriaRowIsNoPanicAndNoCriterion(t *testing.T) {
	for _, row := range []string{"| `NFR-001` |", "|"} {
		in := baseInput(t)
		in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
			"| `NFR-001` | a check of it |", row, 1))}
		rows := CheckFull(in)
		wantFail(t, rows, "criteria", "NFR-001")
	}
}

// classes: the PRD's cells agree with the latest copy, and a requirement
// waiting for its next bet — empty MoSCoW and empty or — Phase, no copy line —
// agrees too (C5). Only REQ/NFR-headed lines are copy lines; a copy with
// fact, option and bet lines parses clean (C6).
func TestClassesComparesWithTheLatestCopy(t *testing.T) {
	in := baseInput(t)
	in.Copy = []byte("REQ-001: Should, phase 1\nREQ-002: Should, phase 2\nNFR-001: Must, phase 1\nbet\n")
	rows := CheckFull(in)
	wantFail(t, rows, "classes", "REQ-001")

	// The copy names a requirement the PRD does not have.
	in = baseInput(t)
	in.Copy = []byte(testCopy + "REQ-009: Could, phase 3\n")
	rows = CheckFull(in)
	wantFail(t, rows, "classes", "REQ-009")

	// A PRD row that waits for its next bet: empty MoSCoW, a — Phase and no
	// copy line agrees (C5). Any other no-line row fails.
	in = baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |",
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |\n| REQ-009 | Later. |  | — | F-0001#6 |", 1))}
	rows = CheckFull(in)
	wantResult(t, rows, "classes", resPass)

	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |",
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |\n| REQ-009 | Later. | Should |  | F-0001#6 |", 1))}
	rows = CheckFull(in)
	wantFail(t, rows, "classes", "REQ-009") // a set MoSCoW without a copy line

	// Both cells must be empty to wait: an empty MoSCoW with a set Phase and
	// no copy line fails (round-2 note N1 — the exemption needs the whole row).
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |",
		"| NFR-001 | Delta. | Must | 1 | F-0001#4 |\n| REQ-009 | Later. |  | 1 | F-0001#6 |", 1))}
	rows = CheckFull(in)
	wantFail(t, rows, "classes", "REQ-009")

	// A REQ-headed line that breaks the grammar fails.
	in = baseInput(t)
	in.Copy = []byte(testCopy + "REQ-001: Later\n")
	rows = CheckFull(in)
	wantFail(t, rows, "classes", "REQ-001")
}

func TestInventoryVerifiesEachVersion(t *testing.T) {
	in := baseInput(t)
	in.Versions = map[int][]byte{1: []byte("fact\tclass\tscope\tcopy\n1\tneed\tin\t9\n")}
	rows := CheckFull(in)
	wantFail(t, rows, "inventory", "sha256")

	in = baseInput(t)
	in.Versions = map[int][]byte{}
	rows = CheckFull(in)
	wantFail(t, rows, "inventory", "did not run: version 1")
	wantFail(t, rows, "ran", "did not run") // an absent input is a did-not-run (note N2)

	// A version whose bytes match the index but are not a readable inventory:
	// inventory fails and the rules of the universe did not run (C13).
	bad := []byte("fact\tclass\tscope\tcopy\n1\tweird\tin\t9\n")
	in = baseInput(t)
	in.Index = [][]string{{"1", fmt.Sprintf("%x", sha256.Sum256(bad)), "M2e"}}
	in.Versions = map[int][]byte{1: bad}
	rows = CheckFull(in)
	wantFail(t, rows, "inventory", "version 1")
	for _, rule := range []string{"covers", "needs", "measures", "quotes"} {
		wantFail(t, rows, rule, notRun)
	}
	wantFail(t, rows, "ran", notRun)
}

// The post-bet universe is the latest confirmed version — the greatest n of
// the index, wherever its row sits (C13, round-2 F1). Version 1 carries a
// need fact no requirement covers; version 2, the index's latest, does not —
// a selection by the smallest n, or by the row's position, fails needs.
func TestTheUniverseIsTheLatestConfirmedVersion(t *testing.T) {
	var v1, v2 bytes.Buffer
	var rows [][]string
	for _, s := range testSpans() {
		rows = append(rows, []string{s[colFact], s[colClass], "in", "9"})
	}
	rows = append(rows, []string{"Q-001", "context", "in", "9"})
	if err := tsv.Write(&v2, InventorySchema, rows); err != nil {
		t.Fatal(err)
	}
	if err := tsv.Write(&v1, InventorySchema, append(rows, []string{"9", "need", "in", "9"})); err != nil {
		t.Fatal(err)
	}
	index := [][]string{
		{"1", fmt.Sprintf("%x", sha256.Sum256(v1.Bytes())), "M2e"},
		{"2", fmt.Sprintf("%x", sha256.Sum256(v2.Bytes())), "M2e"},
	}
	for _, idx := range [][][]string{index, {index[1], index[0]}} {
		in := baseInput(t)
		in.Index = idx
		in.Versions = map[int][]byte{1: v1.Bytes(), 2: v2.Bytes()}
		got := CheckFull(in)
		wantResult(t, got, "needs", resPass)
		wantResult(t, got, "inventory", resPass)
	}
}

// quotes: a requirement that covers an answer's fact also names the fact of
// that answer's quote — here F-0001#7, whose span holds Q-001's quote.
func TestQuotesFollowsTheChain(t *testing.T) {
	in := baseInput(t)
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD,
		"F-0004#1, F-0001#7", "F-0004#1", 1))}
	rows := CheckFull(in)
	wantFail(t, rows, "quotes", "7")

	// An answered question with no quote offsets asks nothing more.
	in = baseInput(t)
	in.Questions = [][]string{{"Q-001", "intent", "band.B", "", "", "", "", "", "q", "2026-01-01T00:00:00Z", "layup run", "2026-01-01T01:00:00Z", "answered"}}
	rows = CheckFull(in)
	wantResult(t, rows, "quotes", resPass)
}

// sections: a delivered requirement needs a non-empty docs/spec/ section whose
// heading holds its ID; one named only by unmerged tasks is not delivered.
func TestSectionsJudgesOnlyTheDelivered(t *testing.T) {
	in := baseInput(t)
	in.Merged = nil // no task merged: nothing is delivered
	rows := CheckFull(in)
	wantResult(t, rows, "sections", resPass)

	in = baseInput(t)
	delete(in.SpecFiles, "loop.md")
	in.SpecFiles["loop.md"] = []byte("## REQ-001 — One\n\ntext.\n## REQ-002 — Two\n\n## NFR-001 — Three\n\ntext.\n")
	rows = CheckFull(in)
	wantFail(t, rows, "sections", "REQ-002") // the section exists but is empty
}

func TestTasksJudgesTheRegister(t *testing.T) {
	var b bytes.Buffer
	if err := tsv.Write(&b, tasksSchema, [][]string{
		{"T-bbbb", "2", "build", "", "small", "", "", "2026-01-01T02:00:00Z"},
	}); err != nil {
		t.Fatal(err)
	}
	in := baseInput(t)
	in.Register = b.Bytes()
	rows := CheckFull(in)
	wantFail(t, rows, "tasks", "T-bbbb")

	in.Register = nil
	rows = CheckFull(in)
	wantFail(t, rows, "tasks", notRun)
	wantFail(t, rows, "ran", notRun)
}

// ran: a row that did not run is a fail, and a table that lacks a rule's row
// fails ran — did not run never counts as passed (NFR-004, C12).
func TestRanRowJudgesTheTable(t *testing.T) {
	rows := CheckFull(baseInput(t))
	wantResult(t, rows, "ran", resPass)

	missing := rows[:len(rows)-2] // drop `tasks` and `ran`
	if r := ranRow(missing); r[1] != resFail {
		t.Errorf("a table without a rule's row: ran gives %v, want fail", r)
	}
	withDidNotRun := slices.Clone(rows[:len(rows)-1])
	withDidNotRun[0][1], withDidNotRun[0][2] = resFail, notRun+"spec/spans.tsv"
	if r := ranRow(withDidNotRun); r[1] != resFail || !strings.Contains(r[2], "did not run") {
		t.Errorf("a did-not-run row: ran gives %v, want fail with the reason", r)
	}
}

// A rule whose input is absent did not run: no register, no PRD, an
// unreadable latest version — each is a fail row and fails ran (C13).
func TestCheckFullDidNotRunIsNeverAPass(t *testing.T) {
	in := baseInput(t)
	in.PRDs = nil
	rows := CheckFull(in)
	for _, rule := range []string{"covers", "needs", "measures", "criteria", "classes", "quotes", "sections"} {
		wantFail(t, rows, rule, notRun)
	}
	wantFail(t, rows, "ran", notRun)
}

// WriteSpecCheck writes the head in the first column; a Read of the block
// gives the rows back, the clear rows' reasons included (C10).
func TestWriteSpecCheckRoundTrip(t *testing.T) {
	var b bytes.Buffer
	head := strings.Repeat("a", 40)
	rows := CheckFull(preBetInput(t))
	if err := WriteSpecCheck(&b, head, rows); err != nil {
		t.Fatal(err)
	}
	back, err := tsv.Read(b.Bytes(), TaskSpecCheckSchema)
	if err != nil {
		t.Fatal(err)
	}
	if len(back) != 10 {
		t.Fatalf("the written table has %d rows, want 10", len(back))
	}
	for i, r := range back {
		if r[0] != head || r[1] != rows[i][0] || r[2] != rows[i][1] || r[3] != rows[i][2] {
			t.Errorf("row %d reads %v, want head plus %v", i, r, rows[i])
		}
	}
	if _, d := ruleOf(t, rows, "classes"); !strings.Contains(b.String(), d) {
		t.Errorf("the written table drops the clear reason %q", d)
	}
}

func TestParseReqLine(t *testing.T) {
	id, a, ok := parseReqLine("REQ-001: Must, phase 1")
	if !ok || id != "REQ-001" || a.Class != "Must" || a.Phase != "1" {
		t.Errorf("parseReqLine gives %s %v %v", id, a, ok)
	}
	if _, a, ok = parseReqLine("NFR-001: Won't, phase —"); !ok || a.Phase != "" {
		t.Errorf("a Won't's phase — gives %v %v", a, ok)
	}
	for _, l := range []string{"fact 3: out of scope", "option D-1: a", "bet", "REQ-001: Later", "prose"} {
		if _, _, ok := parseReqLine(l); ok {
			t.Errorf("parseReqLine takes the non-requirement line %q", l)
		}
	}
}

// HasSection (task T-yeav, #404): the one predicate that sections and the
// handoff item spec-sections both call.
func TestHasSectionNeedsANonEmptySectionWhoseHeadingHoldsTheID(t *testing.T) {
	files := func(text string) map[string][]byte { return map[string][]byte{"a.md": []byte(text)} }
	if !HasSection(files("## REQ-005 — x\nbody\n"), "REQ-005") {
		t.Fatal("the control, a heading and a body line, is no section")
	}
	for name, text := range map[string]string{
		"no body line":            "## REQ-005 — x\n",
		"only blank lines":        "## REQ-005 — x\n\n  \n",
		"the body is the next":    "## REQ-005 — x\n## Next\nbody\n",
		"the ID in the body":      "## Other\nREQ-005 is here\n",
		"plain text, no heading":  "REQ-005 — x\nbody\n",
		"an unrelated heading ID": "## REQ-006 — x\nbody\n",
	} {
		if HasSection(files(text), "REQ-005") {
			t.Errorf("%s: a section", name)
		}
	}
}

func TestHasSectionReadsTheIDAtItsBoundaries(t *testing.T) {
	has := func(heading, id string) bool {
		return HasSection(map[string][]byte{"a.md": []byte(heading + "\nbody\n")}, id)
	}
	for _, c := range []struct {
		heading, id string
		want        bool
	}{
		{"## REQ-005 — x", "REQ-00", false}, {"## REQ-0051 — x", "REQ-005", false}, {"## XREQ-005 — x", "REQ-005", false},
		{"## REQ-00 — x", "REQ-00", true}, {"## (REQ-005)", "REQ-005", true}, {"## REQ-005", "REQ-005", true},
	} {
		if got := has(c.heading, c.id); got != c.want {
			t.Errorf("%q holds %q: %v, want %v", c.heading, c.id, got, c.want)
		}
	}
}

func TestHasSectionOfNoIDOrNoFilesIsFalse(t *testing.T) {
	files := map[string][]byte{"a.md": []byte("## REQ-005 — x\nbody\n")}
	if HasSection(files, "") || HasSection(nil, "REQ-005") || HasSection(map[string][]byte{}, "REQ-005") {
		t.Error("no ID, or no file, is a section")
	}
	if !HasSection(files, "REQ-005") {
		t.Error("the control, the real ID, is no section")
	}
}

// U7 (task T-4naq, row 66ab; loop.md, The full spec check): the facts that FullInput.Reopening lists are exempt from needs and measures, those
// rules only, and only they: the same input without the list fails, and an uncovered fact that is not listed still fails.
func TestAReclassedFactIsExemptFromNeedsAndMeasuresWhileItIsListed(t *testing.T) {
	bare := func(in FullInput) FullInput { // fact 8 is a need and fact 9 a measure, and the waits section names neither
		in.Index, in.Versions = inventoryOfSpans(t, []string{"9", "measure", "in", "9"})
		in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(testPRD, "covers F-0001#8.", "states no fact.", 1))}
		return in
	}
	in := bare(baseInput(t))
	wantFail(t, CheckFull(in), "needs", "8")
	wantFail(t, CheckFull(in), "measures", "9")
	in.Reopening = []string{"8", "9"}
	rows := CheckFull(in)
	wantResult(t, rows, "needs", resPass)
	wantResult(t, rows, "measures", resPass)
	wantResult(t, rows, "ran", resPass)
	in.Reopening = []string{"8"}
	wantResult(t, CheckFull(in), "needs", resPass)
	wantFail(t, CheckFull(in), "measures", "9")
	in.Reopening = []string{"9"}
	wantFail(t, CheckFull(in), "needs", "8")
	wantResult(t, CheckFull(in), "measures", resPass)
	// The other rules judge it as before: a listed fact does not excuse a requirement without a criterion.
	in = bare(baseInput(t))
	in.Reopening = []string{"8", "9"}
	in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(strings.Replace(string(in.PRDs["PRD-0001-x.md"]), "| `REQ-002` | a check of it |\n", "", 1))}
	wantFail(t, CheckFull(in), "criteria", "REQ-002")
}
