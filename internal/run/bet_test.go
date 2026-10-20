package run

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/spec"
)

// The step bet (docs/spec/loop.md, The step bet; task T-2vva, #366, row 66j): the pure parts. A test marked "control" or "negative" is
// green at the inert stub by design (its red twin is beside it).

const (
	btSHA     = "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1"
	btCreated = "2026-10-08T13:00:00Z"
)

func btRow(m, kind, usd, hours, verdict string) string {
	return strings.Join([]string{m, kind, "901", btSHA[:64], usd, hours, "905", "db=pg", verdict, "77", btCreated}, "\t") + "\n"
}

const btHeader = "milestone\tkind\tbrief\tbrief_sha256\tappetite_usd\tappetite_hours\tanswer\toptions\tverdict\tapprover\tdecided\n"

// U1: the rules of the owner of bets, each with its line, its column and its reason; the control is a valid table.
func TestABetsTableRefusesEachBrokenRule(t *testing.T) {
	for _, c := range []struct{ name, table, want string }{
		{"a second first bet", btRow("1", "first", "40.0", "30.0", "bet") + btRow("2", "first", "40.0", "30.0", "bet"), `line 3, column "kind"`},
		{"a later bet before the first", btRow("2", "later", "40.0", "30.0", "bet"), `line 2, column "kind"`},
		{"no brief hash", strings.Replace(btRow("1", "first", "40.0", "30.0", "bet"), btSHA[:64], "—", 1), `line 2, column "brief_sha256"`},
		{"a bet with no appetite", strings.Replace(btRow("1", "first", "—", "30.0", "bet"), "—", "—", 1), `line 2, column "appetite_usd"`},
		{"a no-bet with an appetite", btRow("1", "first", "40.0", "30.0", "no-bet"), `line 2, column "appetite_usd"`},
	} {
		_, err := ReadBets([]byte(btHeader + c.table))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %s", c.name, err, c.want)
		}
	}
	rows, err := ReadBets([]byte(btHeader + btRow("1", "first", "40.0", "30.0", "bet") + btRow("2", "later", "10.0", "5.0", "bet") + strings.Replace(strings.Replace(btRow("3", "later", "40.0", "30.0", "no-bet"), "40.0", "—", 1), "30.0", "—", 1)))
	if err != nil || len(rows) != 3 {
		t.Errorf("control: %d rows, %v; want 3 and no error", len(rows), err)
	}
}

// U2: the row of a bet is its answer and its brief; a no-bet has no appetite; an integer cap is written as a decimal.
func TestTheBetRowIsItsAnswerAndItsBrief(t *testing.T) {
	copyRow := []string{"905", "1", "1", "77", "carol", "", btCreated, "2026-10-08T14:00:00Z", "abc", "copies/905-1.md"}
	approver := []string{"77", "approver", "carol", btCreated, "905"}
	posted := Posted{ID: 901, Seen: 1, SHA: btSHA[:64]}
	a := spec.Answer{Verdict: "bet", Options: map[string]string{"web": "go", "db": "pg"}}
	got := betRow("1", posted, spec.Cap{USD: "40", Hours: "30.5"}, a, copyRow, approver)
	want := []string{"1", "first", "901", btSHA[:64], "40.0", "30.5", "905", "db=pg web=go", "bet", "77", "2026-10-08T14:00:00Z"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("betRow = %q; want %q", got, want)
	}
	got = betRow("1", posted, spec.Cap{USD: "40", Hours: "30"}, spec.Answer{Verdict: "no-bet"}, copyRow, approver)
	want = []string{"1", "first", "901", btSHA[:64], "", "", "905", "", "no-bet", "77", "2026-10-08T14:00:00Z"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("betRow of a no-bet = %q; want %q", got, want)
	}
}

// U5: the options column is sorted by decision.
func TestTheOptionLinesAreTheColumnOptions(t *testing.T) {
	if got := optionsColumn(map[string]string{"web": "go", "db": "pg"}); !slices.Equal(got, []string{"db=pg", "web=go"}) {
		t.Errorf("optionsColumn = %q; want [db=pg web=go]", got)
	}
	if got := optionsColumn(nil); len(got) != 0 {
		t.Errorf("optionsColumn of none = %q; want none", got)
	}
}

// U6: the phase tags are the first column of the table of section 9, whatever the tag is.
func TestThePhaseTagsAreTheFirstColumnOfSection9(t *testing.T) {
	prd := "## 8. Success metrics\n\n| 9 | not here |\n\n## 9. Rollout & phases\n\nThe tags.\n\n| Phase | Name | Ships |\n| ----- | ---- | ----- |\n| 1 | A | x |\n| alpha | B | y |\n| 3 | C | z |\n\n## 10. Risks\n\n| 4 | D | w |\n"
	if got := phaseTagsOf([]byte(prd)); !slices.Equal(got, []string{"1", "alpha", "3"}) {
		t.Errorf("tags = %q; want [1 alpha 3]", got)
	}
	for name, in := range map[string]string{"no section": "## 8. x\n| 1 | a |\n", "no table": "## 9. Rollout\n\nprose\n", "a header only": "## 9. Rollout\n\n| Phase | Name |\n| --- | --- |\n"} {
		if got := phaseTagsOf([]byte(in)); len(got) != 0 {
			t.Errorf("negative %s: tags = %q; want none", name, got)
		}
	}
}

func btCopy(id, seen, author, app, created, body string) answerCopy {
	return answerCopy{Row: []string{id, seen, "1", author, "login", app, created, "2026-10-08T15:00:00Z", "sha", "copies/" + id + ".md"}, Body: body}
}

// U3: the answer is the earliest first copy of the approver, after the brief, that parses; the others answer nothing.
func TestTheAnswerIsTheEarliestFirstCopyOfTheApproverAfterTheBrief(t *testing.T) {
	good := "REQ-001: Must, phase 1\nfact 1: out of scope\noption db: pg\nbet\n"
	in := betInput{Required: []string{"REQ-001"}, Tags: []string{"1", "2"}, Facts: []string{"1", "Q-001"}, Survivors: map[string][]string{"db": {"pg", "sqlite"}}}
	brief := []string{"90", "1", "1", "999", "layup", "layup", btCreated, "2026-10-08T13:00:01Z", "sha", "copies/90-1.md"}
	approver := []string{"303", "idea-owner", "carol", "2026-10-08T12:00:00Z", "802"}
	copies := []answerCopy{
		btCopy("100", "1", "303", "", "2026-10-08T14:01:00Z", good), // valid, later
		btCopy("99", "1", "303", "", "2026-10-08T14:00:00Z", good),  // the answer: 99 before 100 by number, not by text
		btCopy("80", "1", "303", "", "2026-10-08T12:30:00Z", good),  // before the brief
		btCopy("95", "1", "303", "", btCreated, good),               // at the brief's time, not after it
		btCopy("91", "1", "303", "layup", "2026-10-08T14:00:00Z", good),
		btCopy("92", "1", "202", "", "2026-10-08T14:00:00Z", good), // the Operator
		btCopy("93", "2", "303", "", "2026-10-08T14:00:00Z", good), // an edit
		btCopy("94", "1", "303", "", "2026-10-08T14:00:00Z", "REQ-001: Must, phase 1\n"),
	}
	a, row, ok := betAnswer(copies, brief, approver, in)
	if !ok || row[0] != "99" || a.Verdict != "bet" || a.Options["db"] != "pg" {
		t.Fatalf("betAnswer = %+v, %v, %v; want the copy 99 with the option db=pg", a, row, ok)
	}
	other := []string{"404", "approver", "bob", "2026-10-08T12:00:00Z", "803"}
	if _, row, ok := betAnswer(append(copies, btCopy("101", "1", "404", "", "2026-10-08T14:02:00Z", good)), brief, other, in); !ok || row[0] != "101" {
		t.Errorf("with the approver bob: %v, %v; want the copy 101, since the idea owner's 99 answers nothing", row, ok)
	}
	if _, _, ok := betAnswer(copies[2:], brief, approver, in); ok {
		t.Errorf("negative: only the copies that answer nothing gave an answer")
	}
}

// U13: a settled decision is an option line for a decision that had more than one survivor.
func TestTheSettledDecisionsAreTheOptionLines(t *testing.T) {
	survivors := map[string][]string{"db": {"pg", "sqlite"}, "web": {"go"}}
	if got := settledDecisions(map[string]string{"db": "pg", "web": "go"}, survivors); !slices.Equal(got, []string{"db"}) {
		t.Errorf("settled = %q; want [db] (web had one survivor)", got)
	}
	if got := settledDecisions(nil, survivors); len(got) != 0 {
		t.Errorf("negative: settled with no option lines = %q; want none", got)
	}
}

// U14: the requirements of an option are its mapped rows.
func TestTheRequirementsOfAnOptionAreItsMappedRows(t *testing.T) {
	draft := [][]string{{"R-002", "non-functional"}, {"R-004", "requirement"}, {"R-005", "success-criterion"}}
	row := func(c string) []string { return []string{"pg", c, "pass", "docs/x.md"} }
	got, err := archRequirements("db", "pg", [][]string{row("R-004"), row("R-002"), row("NFR-002"), row("F-3"), row("Inv-1"), row("NFR-009")}, draft)
	if err != nil || !slices.Equal(got, []string{"NFR-002", "REQ-004"}) {
		t.Errorf("ids = %q, %v; want [NFR-002 REQ-004] once each (NFR-009 is no requirement of the draft)", got, err)
	}
	for _, none := range [][]string{{"F-3", "Inv-1", "ADR-0001", "Q-002"}, {"R-005"}, nil} {
		var rows [][]string
		for _, c := range none {
			rows = append(rows, row(c))
		}
		if _, err := archRequirements("db", "pg", rows, draft); err == nil || err.Error() != "the decision db: its option pg names no requirement" {
			t.Errorf("rows %q: err = %v; want the decision db: its option pg names no requirement", none, err)
		}
	}
}

const btPRD = "# PRD\n\n## 6. Functional requirements\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| --- | --- | --- | --- | --- |\n| REQ-001 | s | | | F-0001#1 |\n\n**Non-goals.** None.\n\n## 7. Non-functional requirements\n"

// U11: the rendered files are the one PRD and the numbered facts record, by code's rendering; the other cases are refused.
func TestTheRenderedRecordsAreOnePRDAndTheFactsFile(t *testing.T) {
	ans := spec.Answer{Reqs: map[string]spec.ReqAnswer{"REQ-001": {Class: "Must", Phase: "1"}}, Facts: map[string]spec.FactAnswer{"2": {Out: true}}, Verdict: "bet"}
	version := [][]string{{"1", "need", "in", "905"}, {"2", "constraint", "out", "905"}}
	texts := map[string]string{"1": "a", "2": "b"}
	files, recs, err := renderedFiles(map[string][]byte{"PRD-0001-x.md": []byte(btPRD)}, nil, "F-0005", ans, version, texts)
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	if !slices.Equal(paths, []string{"docs/facts/F-0005-numbered-facts.md", "docs/prd/PRD-0001-x.md"}) {
		t.Fatalf("paths = %q; want the facts record F-0005 and the PRD", paths)
	}
	if got := string(files["docs/facts/F-0005-numbered-facts.md"]); got != "1. **need:** a\n2. **constraint (out of scope):** b\n" {
		t.Errorf("the facts record = %q", got)
	}
	if got := string(files["docs/prd/PRD-0001-x.md"]); !strings.Contains(got, "| REQ-001 | s | Must | 1 | F-0001#1 |") || !strings.Contains(got, "Out of scope by the bet: 2.") {
		t.Errorf("the PRD = %q; want the priorities and the out-of-scope fact", got)
	}
	for _, r := range recs {
		if want, f := files[r.Path], (func() []byte { b, _ := rules.Render(r, []byte(btPRD)); return b })(); string(want) != string(f) {
			t.Errorf("record %s is not code's rendering of its file", r.Path)
		}
	}
	two := map[string][]byte{"PRD-0001-x.md": []byte(btPRD), "PRD-0002-y.md": []byte(btPRD)}
	for _, c := range []struct {
		name     string
		prds     map[string][]byte
		existing map[string][]byte
		want     string
	}{{"two PRDs", two, nil, "one docs/prd"}, {"no PRD", nil, nil, "one docs/prd"}, {"a facts record that exists", map[string][]byte{"PRD-0001-x.md": []byte(btPRD)}, map[string][]byte{"F-0004-numbered-facts.md": {}}, "numbered facts record"}} {
		if _, _, err := renderedFiles(c.prds, c.existing, "F-0005", ans, version, texts); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %q", c.name, err, c.want)
		}
	}
}

// U-version: the next version of the inventory, its skip when it names this copy already, and the check of the hash of the one before.
func TestTheNextVersionOfTheInventory(t *testing.T) {
	facts := []spec.InventoryFact{{ID: "1", Class: "need", Text: "a"}, {ID: "2", Class: "constraint", Text: "b"}}
	ans := spec.Answer{Milestone: "1", Facts: map[string]spec.FactAnswer{"2": {Out: true}}}
	files, rows, err := nextVersion(nil, nil, facts, ans, 905)
	if err != nil || len(files) != 2 || len(rows) != 2 || rows[1][2] != "out" {
		t.Fatalf("the first version: %d files, %d rows, %v", len(files), len(rows), err)
	}
	versions := map[int][]byte{1: files["spec/inventory-1.tsv"]}
	if again, again2, err := nextVersion(files["spec/inventory.tsv"], versions, facts, ans, 905); err != nil || again != nil || len(again2) != 2 {
		t.Errorf("a version that names this copy: %d files, %d rows, %v; want no file and its rows", len(again), len(again2), err)
	}
	next, _, err := nextVersion(files["spec/inventory.tsv"], versions, facts, ans, 906)
	if err != nil || next["spec/inventory-2.tsv"] == nil {
		t.Errorf("the version of another copy: %v, %v; want spec/inventory-2.tsv", next, err)
	}
	versions[1] = append(slices.Clone(versions[1]), ' ')
	if _, _, err := nextVersion(files["spec/inventory.tsv"], versions, facts, ans, 906); err == nil || !strings.Contains(err.Error(), "sha256") {
		t.Errorf("a version whose bytes differ from the index: err = %v; want the sha256", err)
	}
}

// U15: bet follows shape, and only once the Scaffold is done.
func TestBetFollowsShapeOnlyOnceTheScaffoldIsDone(t *testing.T) {
	if got := restartNames(true, true); fmt.Sprint(got[len(got)-2:]) != "[shape bet]" {
		t.Errorf("the last rows = %v; want shape, bet", got[len(got)-2:])
	}
	for _, scaffold := range []bool{false} {
		if got := restartNames(true, scaffold); slices.Contains(got, "bet") {
			t.Errorf("negative: scaffold %v has the row bet: %v", scaffold, got)
		}
	}
}

// The readers of an event bound skip the bound of code (Q9, plan amendment 1): the rendered commit has no session.
func TestACodeBoundWithNoSessionIsSkippedByItsReaders(t *testing.T) {
	shaA, shaR := strings.Repeat("a", 40), strings.Repeat("b", 40)
	at := "2026-10-08T14:00:00Z"
	events := [][]string{{"1", "attempt", "1", "", shaA, "", "", at}, {"2", "bound", "1", "S-aaaaaaaa", "", shaA, "task/T-aaaa/1", at}, {"3", "bound", "1", "", "", shaR, "task/T-aaaa/1", at}}
	if sha, session := lastBound(events, 1); sha != shaA || session != "S-aaaaaaaa" {
		t.Errorf("lastBound = %s, %q; want the last bound of a session", sha, session)
	}
	start := make([]string, sessionRecords+1)
	start[sessionID], start[sessionTask], start[sessionRecords] = "S-aaaaaaaa", "T-aaaa", strings.Repeat("c", 40)
	if rev, err := registerRevOf("T-aaaa", shaR, events, [][]string{start}, strings.Repeat("d", 40)); err != nil || rev != strings.Repeat("d", 40) {
		t.Errorf("registerRevOf of the rendered commit = %q, %v; want the commit that the caller gives", rev, err)
	}
	if _, err := authorsFrom(Verification, "T-aaaa", [][]string{start}, events, nil, func(string) (bool, error) { return true, nil }); err != nil {
		t.Errorf("authorsFrom with a bound of code: %v; want none", err)
	}
	if st := loopPosition(loopView{Task: "T-aaaa", Kind: "specification", Events: events[:1:1]}); st.Pos == posStop {
		t.Fatalf("control: %+v", st)
	}
	if st := loopPosition(loopView{Task: "T-aaaa", Kind: "specification", Events: [][]string{events[0], events[2]}}); st.Pos == posStop {
		t.Errorf("loopPosition with a bound of code = %+v; want no stop", st)
	}
	if !renderedBound(events, 1) || renderedBound(events[:2], 1) {
		t.Errorf("renderedBound = %v, %v; want true with the bound of code and false without", renderedBound(events, 1), renderedBound(events[:2], 1))
	}
}
