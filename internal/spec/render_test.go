package spec

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The rendering of the bet (docs/spec/loop.md, REQ-012 — The rendered commit
// and the confirmed inventory; task T-3q0b, #350): the MoSCoW and Phase cells
// of the requirement tables and the non-goals portion change, and nothing
// else; the numbered facts record is one deterministic line per fact.

// prdFixture is a PRD in small: the tables of §6 and §7, the acceptance
// criteria with their quoted IDs, the non-goals paragraph and a matrix of
// §12 that must keep each byte.
var prdFixture = `# PRD

## 3. Goals & non-goals

**Non-goals.** The target ships no plugin system and no hosted service.

## 6. Functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| REQ-001 | The app does one thing. | Must | 1 | F-0001#1 |
| REQ-015 | The app paints itself. | Won't | — | F-0001#2 |

## 7. Non-functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| NFR-001 | Git is the record. | Must | 1 | F-0001#3 |

### 7.1 Acceptance criteria

| Requirement | Criterion |
| ` + "`REQ-001`" + ` | It does one thing. |

## 12. Requirements traceability matrix

| REQ | Task | Test |
| REQ-001 | T-3q0b | this test |
| REQ-001 | T-3q0b | a wide row | is | no requirement row |
`

func renderAnswer() Answer {
	return Answer{Reqs: map[string]ReqAnswer{
		"REQ-001": {"Should", "3"},
		"REQ-015": {"Won't", ""},
		"NFR-001": {"Could", "alpha"},
	}, Facts: map[string]FactAnswer{
		"2":     {Out: true},
		"10":    {Out: true},
		"Q-001": {Out: true},
		"Q-002": {Out: true},
	}}
}

func TestRenderPRDChangesOnlyThePermittedCells(t *testing.T) {
	out, f := RenderPRD([]byte(prdFixture), renderAnswer())
	if len(f) != 0 {
		t.Fatalf("findings %v", f)
	}
	got := string(out)
	for _, want := range []string{
		"| REQ-001 | The app does one thing. | Should | 3 | F-0001#1 |",
		"| REQ-015 | The app paints itself. | Won't | — | F-0001#2 |",
		"| NFR-001 | Git is the record. | Could | alpha | F-0001#3 |",
		"| REQ-001 | T-3q0b | this test |",
		"| REQ-001 | T-3q0b | a wide row | is | no requirement row |",
		"Out of scope by the bet: 2, 10, Q-001, Q-002.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the rendering lacks %q", want)
		}
	}
	// The whole-PRD diff: mask the MoSCoW and Phase cells of the requirement
	// rows of §6 and §7 and the generated non-goals portion; each other byte
	// is unchanged.
	if maskPRD(got) != maskPRD(prdFixture) {
		t.Errorf("a byte outside the permitted cells changed:\n%s", maskPRD(got))
	}
	// A second render gives the same bytes.
	out2, f2 := RenderPRD(out, renderAnswer())
	if len(f2) != 0 || string(out2) != got {
		t.Errorf("the second render: findings %v, bytes equal %v", f2, string(out2) == got)
	}
	// A later bet with other out-of-scope marks replaces the portion.
	out3, _ := RenderPRD(out, Answer{Reqs: map[string]ReqAnswer{"REQ-001": {"Should", "3"}, "REQ-015": {"Won't", ""}, "NFR-001": {"Could", "alpha"}},
		Facts: map[string]FactAnswer{"5": {Out: true}}})
	if !strings.Contains(string(out3), "Out of scope by the bet: 5.") || strings.Count(string(out3), "Out of scope by the bet:") != 1 {
		t.Errorf("the generated portion was not replaced:\n%s", out3)
	}
	// A PRD with no non-goals paragraph keeps every byte but the cells.
	bare := "## 6. Functional requirements\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| REQ-001 | x | Must | 1 | F-0001#1 |\n"
	out4, _ := RenderPRD([]byte(bare), Answer{Reqs: map[string]ReqAnswer{"REQ-001": {"Could", "2"}}, Facts: map[string]FactAnswer{"5": {Out: true}}})
	if want := strings.Replace(bare, "| Must | 1 |", "| Could | 2 |", 1); string(out4) != want {
		t.Errorf("the rendering of a PRD with no non-goals:\n%s", out4)
	}
}

// maskPRD replaces the MoSCoW and Phase cells of the requirement rows of §6
// and §7 with a mark and drops the generated non-goals portion, so a byte
// diff masks precisely the permitted locations.
func maskPRD(prd string) string {
	var out []string
	section := 0
	for _, l := range strings.Split(prd, "\n") {
		if strings.HasPrefix(l, "## ") {
			section = map[bool]int{true: 6, false: 7}[strings.HasPrefix(l, "## 6.")]
		}
		if section != 0 && strings.HasPrefix(l, "|") && strings.Count(l, "|") >= 6 {
			cells := strings.Split(l, "|")
			if id := strings.TrimSpace(cells[1]); strings.HasPrefix(id, "REQ-") || strings.HasPrefix(id, "NFR-") {
				cells[3], cells[4] = " X ", " X "
				l = strings.Join(cells, "|")
			}
		}
		if i := strings.Index(l, " Out of scope by the bet:"); i >= 0 {
			l = l[:i]
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}

func TestRenderPRDRefusals(t *testing.T) {
	for _, c := range []struct {
		name   string
		prd    string
		answer Answer
		detail string
	}{
		{"a requirement with no row", prdFixture, Answer{Reqs: map[string]ReqAnswer{"REQ-099": {"Must", "1"}}}, "has no row"},
		{"a requirement twice", strings.Replace(prdFixture, "## 7. Non-functional requirements",
			"| REQ-001 | A repeat. | Must | 2 | F-0001#4 |\n\n## 7. Non-functional requirements", 1), renderAnswer(), "has 2 rows"},
	} {
		if _, f := RenderPRD([]byte(c.prd), c.answer); len(f) == 0 || !strings.Contains(f[0].Detail, c.detail) {
			t.Errorf("%s: findings %v", c.name, f)
		}
	}
}

func TestRenderFactsOfTheVersion(t *testing.T) {
	version := [][]string{{"1", "need", "in", "42"}, {"2", "constraint", "out", "42"}, {"Q-001", "context", "in", "42"}}
	text := map[string]string{"1": "one thing", "2": "a limit", "Q-001": "the answer"}
	out, f := RenderFacts(version, text)
	want := "1. **need:** one thing\n2. **constraint (out of scope):** a limit\nQ-001. **context:** the answer\n"
	if len(f) != 0 || string(out) != want {
		t.Fatalf("findings %v, bytes %q, want %q", f, out, want)
	}
	// The render is deterministic and reports no text that was not supplied.
	for _, c := range []struct {
		name    string
		version [][]string
		text    map[string]string
	}{
		{"a fact with no text", version, map[string]string{"1": "one thing", "Q-001": "the answer"}},
		{"a fact with no class", [][]string{{"1", "", "in", "42"}}, text},
	} {
		if _, f := RenderFacts(c.version, c.text); len(f) == 0 {
			t.Errorf("%s: no finding", c.name)
		}
	}
}

// The inventory tests (the same contract).

func TestTheInventoryReaders(t *testing.T) {
	good := "fact\tclass\tscope\tcopy\n1\tneed\tin\t42\n2\tconstraint\tout\t42\nQ-001\tcontext\tin\t42\n"
	if rows, err := ReadInventory([]byte(good)); err != nil || len(rows) != 3 {
		t.Fatalf("a good inventory: %d rows, %v", len(rows), err)
	}
	index := "n\tsha256\tbet\n1\t" + strings.Repeat("a", 64) + "\tM2e\n"
	if rows, err := ReadInventoryIndex([]byte(index)); err != nil || len(rows) != 1 {
		t.Fatalf("a good index: %d rows, %v", len(rows), err)
	}
	for _, c := range []struct{ name, data, want string }{
		{"class —", "fact\tclass\tscope\tcopy\n1\t—\tin\t42\n", "column \"class\""},
		{"scope —", "fact\tclass\tscope\tcopy\n1\tneed\t—\t42\n", "column \"scope\""},
		{"copy —", "fact\tclass\tscope\tcopy\n1\tneed\tin\t—\n", "column \"copy\""},
		{"copy 0", "fact\tclass\tscope\tcopy\n1\tneed\tin\t0\n", "column \"copy\""},
		{"a fact twice", "fact\tclass\tscope\tcopy\n1\tneed\tin\t42\n1\tneed\tin\t42\n", "key"},
		{"index sha256 —", "n\tsha256\tbet\n1\t—\tM2e\n", "column \"sha256\""},
		{"index bet —", "n\tsha256\tbet\n1\t" + strings.Repeat("a", 64) + "\t—\n", "column \"bet\""},
	} {
		read := ReadInventory
		if strings.HasPrefix(c.name, "index") {
			read = ReadInventoryIndex
		}
		if _, err := read([]byte(c.data)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v, want one that names %s", c.name, err, c.want)
		}
	}
}

// versionRows parses the bytes of a version or an index by its schema.
func versionRows(t *testing.T, schema tsv.Schema, data []byte) [][]string {
	t.Helper()
	rows, err := tsv.Read(data, schema)
	if err != nil {
		t.Fatalf("%s: %v", schema.Name, err)
	}
	return rows
}

func TestTheFirstInventory(t *testing.T) {
	facts := []InventoryFact{{"1", "need", "text one"}, {"2", "constraint", "text two"}, {"Q-001", "context", "text q"}}
	a := Answer{Facts: map[string]FactAnswer{"2": {Out: true}}, Milestone: "M2e"}
	v, idx, f := NextInventory(nil, nil, facts, a, 42)
	if len(f) != 0 {
		t.Fatalf("findings %v, want none", f)
	}
	rows := versionRows(t, InventorySchema, v)
	want := [][]string{{"1", "need", "in", "42"}, {"2", "constraint", "out", "42"}, {"Q-001", "context", "in", "42"}}
	if fmt.Sprint(rows) != fmt.Sprint(want) {
		t.Errorf("version 1 rows %v, want %v", rows, want)
	}
	irows := versionRows(t, InventoryIndexSchema, idx)
	hash := fmt.Sprintf("%x", sha256.Sum256(v))
	if len(irows) != 1 || irows[0][0] != "1" || irows[0][1] != hash || irows[0][2] != "M2e" {
		t.Errorf("index rows %v, want {{1 %s M2e}}", irows, hash)
	}
}

func TestTheNextInventory(t *testing.T) {
	facts := []InventoryFact{{"1", "need", "t1"}, {"2", "constraint", "t2"}, {"Q-001", "context", "tq"}}
	v1, i1, f := NextInventory(nil, nil, facts, Answer{Facts: map[string]FactAnswer{"2": {Out: true}}, Milestone: "M2e"}, 42)
	if len(f) != 0 {
		t.Fatal(f)
	}
	prev, _ := ReadInventory(v1)
	index, _ := ReadInventoryIndex(i1)
	// The second bet reclasses fact 1 and leaves fact 2's out mark.
	a2 := Answer{Facts: map[string]FactAnswer{"1": {Class: "measure"}}, Milestone: "M2f"}
	v2, i2, f2 := NextInventory(prev, index, facts, a2, 77)
	if len(f2) != 0 {
		t.Fatalf("findings %v", f2)
	}
	rows := versionRows(t, InventorySchema, v2)
	want := [][]string{{"1", "measure", "in", "77"}, {"2", "constraint", "out", "77"}, {"Q-001", "context", "in", "77"}}
	if fmt.Sprint(rows) != fmt.Sprint(want) {
		t.Errorf("version 2 rows %v, want %v", rows, want)
	}
	irows := versionRows(t, InventoryIndexSchema, i2)
	if len(irows) != 2 || irows[0][0] != "1" || irows[1][0] != "2" {
		t.Errorf("index rows %v, want 1 and 2", irows)
	}
	// The earlier version's bytes and the earlier index row are retained.
	if !bytes.HasPrefix(i2, i1) {
		t.Errorf("the index lost its first row:\n%s", i2)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(v2))
	if irows[1][1] != hash {
		t.Errorf("the hash of version 2 is %s, want %s", irows[1][1], hash)
	}
}

func TestInventoryRefusals(t *testing.T) {
	facts := []InventoryFact{{"1", "need", "t"}}
	for _, c := range []struct {
		name     string
		previous [][]string
		index    [][]string
		facts    []InventoryFact
		answer   Answer
		copyID   int
	}{
		{"an unknown fact in the answer", nil, nil, facts, Answer{Facts: map[string]FactAnswer{"9": {Out: true}}, Milestone: "M2e"}, 42},
		{"a copy ID that is no positive comment ID", nil, nil, facts, Answer{Milestone: "M2e"}, 0},
		{"a fact with no class", nil, nil, []InventoryFact{{"1", "", "t"}}, Answer{Milestone: "M2e"}, 42},
		{"no milestone for the index", nil, nil, facts, Answer{}, 42},
		{"a second version with no first", nil, [][]string{{"1", strings.Repeat("a", 64), "M2e"}}, facts, Answer{Milestone: "M2f"}, 77},
		{"an index out of order", [][]string{{"1", "need", "in", "42"}}, [][]string{{"3", strings.Repeat("a", 64), "M2e"}}, facts, Answer{Milestone: "M2f"}, 77},
	} {
		if _, _, f := NextInventory(c.previous, c.index, c.facts, c.answer, c.copyID); len(f) == 0 {
			t.Errorf("%s: no finding", c.name)
		}
	}
	// A version 2 needs version 1 — with an index in order so no other rule fires.
	if _, _, f := NextInventory(nil, [][]string{{"1", strings.Repeat("a", 64), "M2e"}}, facts, Answer{Milestone: "M2f"}, 77); len(f) == 0 || f[0].Rule != "version" {
		t.Errorf("a second version with no first: findings %v", f)
	}
	// The answer's class is one of the four or a finding, as for the fact's own —
	// the owner rule, not the writer's enum check, names the fact.
	if _, _, f := NextInventory(nil, nil, facts, Answer{Facts: map[string]FactAnswer{"1": {Class: "bogus"}}, Milestone: "M2e"}, 42); len(f) == 0 || f[0].Rule != "fact" {
		t.Errorf("a class that is none of the four: findings %v", f)
	}
}

// The quote tests (the same contract).

// qpsb is the problem statement, with two fact spans.
var qpsb = []byte("The app does one thing. It sends reminders.")

var qspans = [][]string{
	{"0", "22", "1", "need", "", "The app does one thing."},
	{"23", "44", "2", "constraint", "", "It sends reminders."},
}

// qrow builds a questions.tsv row: id, source, ref, block, phase, kind,
// start, end, question, asked, asker, answered, state.
func qrow(id, source, ref, start, end string) []string {
	return []string{id, source, ref, "idea-owner", "before-delivery", "two-readings", start, end,
		"the question", "2026-10-10T12:00:00Z", "S-0000000a", "2026-10-10T13:00:00Z", "closed"}
}

// arow builds an answers.tsv row: question, answer, comment, author, time.
func arow(q string) []string {
	return []string{q, "the answer", "4001", "1675602", "2026-10-10T13:00:00Z"}
}

// drow builds a draft row: id, kind, statement, covers, criterion, session.
func drow(id, covers string) []string {
	return []string{id, "requirement", "the statement", covers, "the criterion", "S-0000000a"}
}

func TestQuoteFactOfFindsTheOneContainingFact(t *testing.T) {
	for _, c := range []struct {
		name       string
		start, end int
		want       string
	}{
		{"a quote inside fact 1", 4, 7, "1"},
		{"a quote inside fact 2", 30, 40, "2"},
		{"the whole span of fact 1", 0, 22, "1"},
		{"a quote in no fact", 0, 0, "err"}, // invalid: start == end
		{"a quote before the file", -1, 5, "err:bounds"},
		{"a quote just past the file", 40, 50, "err:bounds"},
		{"a quote past the file", 0, 90, "err"},
		{"a quote on two facts", 20, 25, "err"},
	} {
		got, err := QuoteFactOf(c.start, c.end, qspans, qpsb)
		if c.want == "err" || c.want == "err:bounds" {
			if err == nil || (c.want == "err:bounds" && !strings.Contains(err.Error(), "no quote of a problem statement")) {
				t.Errorf("%s: fact %q, error %v, want a bounds error", c.name, got, err)
			}
			continue
		}
		if err != nil || got != c.want {
			t.Errorf("%s: fact %q, error %v, want %q", c.name, got, err, c.want)
		}
	}
	// A span table whose spans overlap gives two candidates.
	two := append(append([][]string{}, qspans...), []string{"10", "30", "3", "need", "", "x"})
	if _, err := QuoteFactOf(25, 28, two, qpsb); err == nil {
		t.Error("two containing facts: no error")
	}
}

func TestCheckQuoteFacts(t *testing.T) {
	questions := [][]string{
		qrow("Q-001", "meaning", "M-001", "4", "7"), // quotes "app", inside fact 1
		qrow("Q-002", "intent", "band.U", "", ""),   // no quote
		qrow("Q-003", "follow-up", "Q-001", "", ""), // asks Q-001 again
	}
	answers := [][]string{arow("Q-001"), arow("Q-002"), arow("Q-003")}
	with := func(rows [][]string, extra ...[]string) [][]string {
		return append(append([][]string{}, rows...), extra...)
	}
	for _, c := range []struct {
		name   string
		covers string
		qs, as [][]string
		want   bool // a finding
	}{
		{"a covered answer and its quote's fact", "1 Q-001", questions, answers, false},
		{"the answer's fact alone", "Q-001", questions, answers, true},
		{"the answer's fact and another fact", "2 Q-001", questions, answers, true},
		{"an answer that has no quote", "Q-002", questions, answers, false},
		{"a follow-up resolves to the quote's fact", "1 Q-003", questions, answers, false},
		{"a follow-up without the quote's fact", "Q-003", questions, answers, true},
		{"a Q that is no answer", "Q-009", questions, answers, true},
		{"a question the answers lack", "Q-014", with(questions, qrow("Q-014", "meaning", "M-011", "", "")), answers, true},
		{"a Q that is no question", "Q-002", questions[:1], answers, true},
		{"a token that is no fact", "abc", questions, answers, true},
		{"a fact number that is no fact", "9", questions, answers, true},
		{"a cyclic chain", "Q-010",
			with(questions, qrow("Q-010", "follow-up", "Q-010", "", "")), with(answers, arow("Q-010")), true},
		{"a follow-up of a question that is not there", "Q-011",
			with(questions, qrow("Q-011", "follow-up", "Q-099", "", "")), with(answers, arow("Q-011")), true},
		{"an unpaired quote offset", "Q-012",
			with(questions, qrow("Q-012", "meaning", "M-009", "7", "")), with(answers, arow("Q-012")), true},
		{"a quote outside the problem statement", "Q-013",
			with(questions, qrow("Q-013", "meaning", "M-010", "4", "90")), with(answers, arow("Q-013")), true},
	} {
		if f := CheckQuoteFacts([][]string{drow("R-001", c.covers)}, c.qs, c.as, qspans, qpsb); (len(f) != 0) != c.want {
			t.Errorf("%s: findings %v, want a finding %v", c.name, f, c.want)
		}
	}
	// An unpaired offset names the pairing, not the parse — the rules are distinct.
	f1 := CheckQuoteFacts([][]string{drow("R-002", "Q-012")},
		with(questions, qrow("Q-012", "meaning", "M-009", "7", "")), with(answers, arow("Q-012")), qspans, qpsb)
	if len(f1) == 0 || !strings.Contains(f1[0].Detail, "one quote offset") {
		t.Errorf("the unpaired offset: findings %v", f1)
	}
	// The good form: each row's tokens resolve; a finding names the row.
	f := CheckQuoteFacts([][]string{drow("R-007", "Q-001")}, questions, answers, qspans, qpsb)
	if len(f) == 0 || !strings.Contains(f[0].Detail, "R-007") || !strings.Contains(f[0].Detail, "1") {
		t.Errorf("the finding names the row and the missing fact: %v", f)
	}
}
