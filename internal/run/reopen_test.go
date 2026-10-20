package run

import (
	"context"
	"crypto/sha256"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// The unit level of the specification task that a changed class reopens (loop.md, The specification task that a changed class reopens;
// task T-4naq, #428, row 66ab): the selection, the identities, the record and its binding, the version and the exemption.

func ivr(fact, class, scope string) []string { return []string{fact, class, scope, "9"} }
func inf(id, class string) spec.InventoryFact {
	return spec.InventoryFact{ID: id, Class: class, Text: "t"}
}

// factsOf gives the facts of the rows (their first column).
func factsOf(rows [][]string) []string {
	var out []string
	for _, r := range rows {
		out = append(out, r[0])
	}
	return out
}

// U1: the facts of the selection are the version's rows in scope whose class is one of the three and differs from Intake's.
func TestTheReclassedFactsAreTheVersionRowsWhoseClassDiffersFromIntake(t *testing.T) {
	version := [][]string{ivr("3", "need", "in"), ivr("5", "need", "in"), ivr("7", "constraint", "in"), ivr("8", "measure", "out"), ivr("9", "context", "in"),
		ivr("12", "measure", "in"), ivr("Q-002", "need", "in"), ivr("Q-003", "context", "in"), ivr("13", "measure", "in"), ivr("14", "need", "out")}
	intake := []spec.InventoryFact{inf("3", "context"), inf("5", "need"), inf("7", "need"), inf("8", "context"), inf("9", "need"), inf("12", "context"),
		inf("Q-002", "context"), inf("Q-003", "context"), inf("13", "measure"), inf("14", "context")}
	got, err := reclassed(version, intake)
	if want := []string{"3", "7", "12", "Q-002"}; err != nil || !slices.Equal(factsOf(got), want) {
		t.Errorf("facts = %v, %v; want %v (a need from context, a constraint from need, a measure from context, a fact of the answers; not an equal class, an out-of-scope fact or a context)", factsOf(got), err, want)
	}
	if _, err := reclassed(append(version, ivr("20", "need", "in")), intake); err == nil || !strings.Contains(err.Error(), "20") {
		t.Errorf("a fact of the version that Intake does not know: %v; want an error that names it", err)
	}
}

// U2: one number space for REQ and NFR; the next number is one more than every ID known, from each of the three sources, and a measure takes none.
func TestTheNextRequirementIDIsOneMoreThanEveryKnownID(t *testing.T) {
	for _, c := range []struct {
		name  string
		ids   [][]string
		first int
	}{
		{"all three sources", [][]string{{"REQ-001", "NFR-002", "REQ-003"}, {"REQ-004"}, {"REQ-005"}}, 6},
		{"the draft alone", [][]string{{"REQ-008", "NFR-002"}, nil, nil}, 9},
		{"the register alone", [][]string{nil, {"NFR-009"}, nil}, 10},
		{"the reopenings alone", [][]string{nil, nil, {"REQ-010"}}, 11},
		{"none", [][]string{nil, nil, nil}, 1},
	} {
		if got := nextNumber(c.ids...); got != c.first {
			t.Errorf("%s: nextNumber = %d; want %d", c.name, got, c.first)
		}
	}
	sel := [][]string{ivr("3", "need", "in"), ivr("7", "constraint", "in"), ivr("12", "measure", "in"), ivr("Q-002", "need", "in")}
	got := allocate(sel, 1, "F-0002", 6, []string{"T-aaaa", "T-bbbb", "T-cccc", "T-dddd"})
	want := [][]string{{"1", "3", "need", "F-0002#3", "REQ-006", "T-aaaa", "9"}, {"1", "7", "constraint", "F-0002#7", "NFR-007", "T-bbbb", "9"},
		{"1", "12", "measure", "F-0002#12", "", "T-cccc", "9"}, {"1", "Q-002", "need", "Q-002", "REQ-008", "T-dddd", "9"}}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("ids = %v; want %v (REQ for a need, NFR for a constraint, none for a measure, one counter in the order of the rows)", got, want)
	}
	if got := allocate(sel[:1], 1, "F-0002", 1000, []string{"T-aaaa"}); len(got) != 1 || got[0][reoRequirement] != "REQ-1000" {
		t.Errorf("the 1000th requirement = %v; want REQ-1000", got)
	}
}

const reoHeader = "n\tfact\tclass\ttoken\trequirement\ttask\tcopy\n"

// U3: each rule of the owner of spec/reopenings.tsv refuses with its line, its column and its reason; the control reads.
func TestAReopeningsTableRefusesEachBrokenRule(t *testing.T) {
	const r1, r2 = "1\t3\tneed\tF-0002#3\tREQ-006\tT-aaaa\t9\n", "1\t12\tmeasure\tF-0002#12\t—\tT-bbbb\t9\n"
	if rows, err := ReadReopenings([]byte(reoHeader + r1 + r2)); err != nil || len(rows) != 2 || rows[1][reoRequirement] != "" {
		t.Fatalf("the control: %d rows, %v; want 2, a measure with the empty requirement", len(rows), err)
	}
	for _, c := range []struct{ name, body, want string }{
		{"a key twice", r1 + r1, `line 3: the key (n "1", fact "3") repeats the key of line 2`},
		{"a class that is context", "1\t3\tcontext\tF-0002#3\tREQ-006\tT-aaaa\t9\n", `line 2, column "class"`},
		{"a measure with an ID", "1\t12\tmeasure\tF-0002#12\tREQ-006\tT-bbbb\t9\n", `line 2, column "requirement": a measure has no requirement ID`},
		{"a need with no ID", "1\t3\tneed\tF-0002#3\t—\tT-aaaa\t9\n", `line 2, column "requirement": a need`},
		{"a need with an NFR", "1\t3\tneed\tF-0002#3\tNFR-006\tT-aaaa\t9\n", `line 2, column "requirement": a need`},
		{"a constraint with a REQ", "1\t3\tconstraint\tF-0002#3\tREQ-006\tT-aaaa\t9\n", `line 2, column "requirement": a constraint`},
		{"an ID of two digits", "1\t3\tneed\tF-0002#3\tREQ-06\tT-aaaa\t9\n", `line 2, column "requirement"`},
		{"a requirement twice", r1 + "1\t4\tneed\tF-0002#4\tREQ-006\tT-bbbb\t9\n", `line 3, column "requirement": the requirement REQ-006 is held by line 2`},
		{"a task twice", r1 + "1\t4\tneed\tF-0002#4\tREQ-007\tT-aaaa\t9\n", `line 3, column "task": the task T-aaaa is held by line 2`},
		{"a token of another fact", "1\t3\tneed\tF-0002#4\tREQ-006\tT-aaaa\t9\n", `line 2, column "token": a token is F-NNNN#3`},
		{"a bare Q token for a number", "1\t3\tneed\tQ-003\tREQ-006\tT-aaaa\t9\n", `line 2, column "token"`},
		{"a token of another Q", "1\tQ-002\tneed\tQ-003\tREQ-006\tT-aaaa\t9\n", `line 2, column "token": a token is Q-002`},
		{"no token", "1\t3\tneed\t—\tREQ-006\tT-aaaa\t9\n", `line 2, column "token": this column never holds the empty value`},
		{"no task", "1\t3\tneed\tF-0002#3\tREQ-006\t—\t9\n", `line 2, column "task": this column never holds the empty value`},
		{"no copy", "1\t3\tneed\tF-0002#3\tREQ-006\tT-aaaa\t—\n", `line 2, column "copy": this column never holds the empty value`},
		{"a version 0", "0\t3\tneed\tF-0002#3\tREQ-006\tT-aaaa\t9\n", `line 2, column "n": a version counts from 1`},
		{"a version that is split", r1 + "2\t4\tneed\tF-0002#4\tREQ-007\tT-bbbb\t9\n" + "1\t5\tneed\tF-0002#5\tREQ-008\tT-cccc\t9\n", `line 4, column "n": the rows of the version 1 are not together`},
	} {
		if _, err := ReadReopenings([]byte(reoHeader + c.body)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
	}
}

// U3b: the rows are bound to the selection of the version at each call: its facts and their order, their classes, the copy and the token.
func TestTheRowsOfTheRecordAreBoundToTheSelection(t *testing.T) {
	sel := [][]string{ivr("3", "need", "in"), ivr("12", "measure", "in")}
	rows := func() [][]string { return allocate(sel, 1, "F-0002", 6, []string{"T-aaaa", "T-bbbb"}) }
	if len(rows()) != 2 {
		t.Fatalf("allocate gave %d rows; want 2", len(rows()))
	}
	if err := bindReopenings(rows(), sel, "F-0002"); err != nil {
		t.Fatalf("the control: %v", err)
	}
	for _, c := range []struct {
		name, want string
		change     func([][]string) [][]string
	}{
		{"a fact less", "2 facts", func(r [][]string) [][]string { return r[:1] }},
		{"a fact more", "2 facts", func(r [][]string) [][]string {
			return append(r, []string{"1", "20", "need", "F-0002#20", "REQ-007", "T-cccc", "9"})
		}},
		{"another order", "fact 12", func(r [][]string) [][]string { r[0], r[1] = r[1], r[0]; return r }},
		{"another class", "class", func(r [][]string) [][]string { r[0][reoClass] = "constraint"; return r }},
		{"another copy", "copy", func(r [][]string) [][]string { r[1][reoCopy] = "10"; return r }},
		{"another record in the token", "token", func(r [][]string) [][]string { r[0][reoToken] = "F-0005#3"; return r }},
	} {
		if err := bindReopenings(c.change(rows()), sel, "F-0002"); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
	}
}

func versionBytes(rows [][]string) []byte {
	b, err := table(spec.InventorySchema, rows)
	if err != nil {
		panic(err)
	}
	return b
}

func indexBytes(t *testing.T, rows ...[]string) []byte {
	t.Helper()
	b, err := table(spec.InventoryIndexSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// U4: the one version is read after its hash is checked; a second version (a later bet) fails closed, and so do an absent index, a version that
// is not there and one whose bytes differ from the index.
func TestALaterVersionIsRefused(t *testing.T) {
	v1 := versionBytes([][]string{ivr("3", "need", "in")})
	v2 := versionBytes([][]string{ivr("3", "need", "in"), ivr("4", "need", "in")})
	sum := func(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
	one := indexBytes(t, []string{"1", sum(v1), "M2e"})
	if rows, err := firstVersion(one, map[int][]byte{1: v1}); err != nil || len(rows) != 1 || rows[0][0] != "3" {
		t.Fatalf("the control: %v, %v; want the row of version 1", rows, err)
	}
	two := indexBytes(t, []string{"1", sum(v1), "M2e"}, []string{"2", sum(v2), "M2g"})
	for _, c := range []struct {
		name  string
		index []byte
		vs    map[int][]byte
		want  string
	}{
		{"a second version", two, map[int][]byte{1: v1, 2: v2}, "the reopening of the classes of a later bet is not built in M2e"},
		{"no index", nil, nil, "no confirmed inventory"},
		{"no version", one, map[int][]byte{}, "spec/inventory-1.tsv"},
		{"another hash", one, map[int][]byte{1: v2}, "sha256"},
	} {
		if _, err := firstVersion(c.index, c.vs); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
	}
}

// U8a: the facts the full check exempts, by the states of loop.md, The full spec check: two re-classed facts, 7 and 12, with the tasks A and B.
func TestTheExemptFactsFollowTheStatesOfTheTwoTasks(t *testing.T) {
	version := [][]string{ivr("7", "need", "in"), ivr("12", "measure", "in"), ivr("5", "need", "in"), ivr("8", "constraint", "out")}
	intake := []spec.InventoryFact{inf("7", "context"), inf("12", "context"), inf("5", "need"), inf("8", "context")}
	rows := [][]string{{"1", "7", "need", "F-0002#7", "REQ-006", "T-aaaa", "9"}, {"1", "12", "measure", "F-0002#12", "", "T-bbbb", "9"}}
	for _, c := range []struct {
		name   string
		rows   [][]string
		merged []string
		task   string
		want   []string
	}{
		{"before the allocation, any head", nil, nil, "T-sp01", []string{"7", "12"}},
		{"the head of A", rows, nil, "T-aaaa", []string{"12"}},
		{"the head of B", rows, nil, "T-bbbb", []string{"7"}},
		{"the head of another task", rows, nil, "T-arch", []string{"7", "12"}},
		{"after A has merged, the head of B", rows, []string{"T-aaaa"}, "T-bbbb", nil},
		{"after A has merged, another head", rows, []string{"T-aaaa"}, "T-arch", []string{"12"}},
		{"after both have merged", rows, []string{"T-aaaa", "T-bbbb"}, "T-arch", nil},
		{"a merged task of no fact", rows, []string{"T-sp01"}, "T-arch", []string{"7", "12"}},
	} {
		got, err := reopeningFacts(version, intake, c.rows, c.merged, c.task)
		slices.Sort(got)
		slices.Sort(c.want)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
	if got, err := reopeningFacts(nil, intake, rows, nil, "T-sp01"); err != nil || len(got) != 0 {
		t.Errorf("before the first bet: %v, %v; want none", got, err)
	}
}

// The draws of a task ID avoid a reservation that no task holds yet, and the draws of one batch avoid each other (condition 4 of the plan review).
func TestADrawAvoidsAnUnopenedReservationAndItsOwnBatch(t *testing.T) {
	store := storeWithSessions(t)
	store.files[reopeningsPath] = []byte(reoHeader + "1\t3\tneed\tF-0002#3\tREQ-006\tT-0000\t9\n")
	s := newSessions(store)
	s.random = drawFrom([4]byte{0, 0, 0, 0}, [4]byte{1, 1, 1, 1})
	if id, err := s.taskID(context.Background()); err != nil || id != "T-1111" {
		t.Errorf("taskID = %q, %v; want T-1111, the draw after the one that the reservation holds", id, err)
	}
	s.random = drawFrom([4]byte{1, 1, 1, 1}, [4]byte{1, 1, 1, 1}, [4]byte{2, 2, 2, 2}, [4]byte{0, 0, 0, 0}, [4]byte{3, 3, 3, 3})
	if ids, err := s.taskIDs(context.Background(), 3); err != nil || !slices.Equal(ids, []string{"T-1111", "T-2222", "T-3333"}) {
		t.Errorf("taskIDs = %v, %v; want T-1111, T-2222 and T-3333: a draw again for the one of the batch and for the reservation", ids, err)
	}
}

// reopenSrc is the source of the prompt file of a step of a task that a class reopens: the row of spec/reopenings.tsv of a need (REQ-006 of the
// fact 3) or of a measure (the fact 12, which names the specification's REQ-005), and a target that holds the numbered facts record F-0004.
func reopenSrc(t *testing.T, step, class string) promptSrc {
	t.Helper()
	src := stepSrc(t, step, "specification")
	row := []string{"1", "3", "need", "F-0004#3", "REQ-006", "T-ab12", "9"}
	src.Requirements = []string{"REQ-006"}
	if class == "measure" {
		row, src.Requirements = []string{"1", "12", "measure", "F-0004#12", "", "T-ab12", "9"}, []string{"REQ-005"}
	}
	src.Reopen = row
	if step != "close-out" {
		src.Extras.Shape = append(src.Extras.Shape, kept{reopeningsPath, []byte("reopenings\n")})
	}
	files := map[string]string{"docs/prd/PRD-0001.md": prdText, "docs/spec/loop.md": specText, "AGENTS.md": "rules\n", "docs/rules/a.md": "rule a\n",
		"docs/facts/F-0003-a.md": "facts\n", "docs/facts/F-0004-numbered-facts.md": "numbered\n",
		"docs/adr/0019-phase-loop.md": "adr 19\n", "docs/adr/0020-route.md": "adr 20\n", "docs/adr/0021-smart-if.md": "adr 21\n"}
	src.Tree = memTree(files)
	return src
}

// U5: the prompt of the step specification of a reopened task is its golden file, for a need and for a measure; the golden files differ from
// the standard one in the instruction and in the requirements (the PRD with no row of the ID, and the row of the record).
func TestTheReopenedSpecificationPromptIsItsGoldenFile(t *testing.T) {
	for _, class := range []string{"need", "measure"} {
		c, err := buildPrompt(reopenSrc(t, "specification", class))
		if err != nil {
			t.Fatalf("%s: %v", class, err)
		}
		want, err := promptGoldens.ReadFile("testdata/prompts/specification-reopened-" + class + ".txt")
		if err != nil {
			t.Fatal(err)
		}
		if got := string(c.Text(rootDir)); got != string(want) {
			i := 0
			for i < len(got) && i < len(want) && got[i] == want[i] {
				i++
			}
			t.Errorf("%s: Text differs from its golden file: got %d bytes, want %d; first difference at byte %d: got %q, want %q", class, len(got), len(want), i,
				got[i:min(len(got), i+40)], string(want[i:min(len(want), i+40)]))
		}
		if m := leftover.FindString(string(c.Text(rootDir))); m != "" {
			t.Errorf("%s: the prompt keeps the placeholder %s", class, m)
		}
	}
}

// Condition 3 of the plan review: every step of a reopened task sees the PRD although no row holds the ID, and the row of the record; the step
// that is not the specification keeps its standard instruction.
func TestEveryStepOfAReopenedTaskSeesThePRDAndItsRow(t *testing.T) {
	for _, step := range []string{"plan-of-a-task", "plan-review", "specification", "verification"} {
		c, err := buildPrompt(reopenSrc(t, step, "need"))
		if err != nil {
			t.Fatalf("%s: %v", step, err)
		}
		got := itemsOf(c, session.KindRequirements)
		if !slices.Contains(got, "docs/prd/PRD-0001.md") || !slices.Contains(got, "spec/reopenings.tsv") || slices.Contains(got, "docs/spec/loop.md") {
			t.Errorf("%s: the requirements = %v; want the PRD with no row of REQ-006, the record of the reopenings and no section", step, got)
		}
		if reopened := strings.Contains(string(c.Text(rootDir)), "Write the requirement that waits for a bet"); reopened != (step == "specification") {
			t.Errorf("%s: holds the reopened instruction: %v", step, reopened)
		}
	}
	// A task that no class reopens gets neither the PRD with no row nor the instruction.
	src := stepSrc(t, "specification", "specification")
	src.Requirements = []string{"REQ-006"}
	c, err := buildPrompt(src)
	if err != nil || slices.Contains(itemsOf(c, session.KindRequirements), "docs/prd/PRD-0001.md") || strings.Contains(string(c.Text(rootDir)), "waits for a bet") {
		t.Errorf("a task with no row of reopenings: %v, %v; want the standard prompt", itemsOf(c, session.KindRequirements), err)
	}
}

// The token's record is the numbered facts record found at the base, and it is the only one: the prompt names no predicted record.
func TestTheReopenedPromptPinsTheNumberedRecord(t *testing.T) {
	for _, c := range []struct{ name, token, want string }{
		{"another record in the token", "F-0005#3", "F-0005"},
		{"no numbered record", "F-0004#3", "numbered facts record"},
		{"two numbered records", "F-0004#3", "2 numbered facts records"},
	} {
		src := reopenSrc(t, "specification", "need")
		src.Reopen[reoToken] = c.token
		files := map[string]string{"docs/prd/PRD-0001.md": prdText, "docs/facts/F-0003-a.md": "facts\n", "docs/facts/F-0004-numbered-facts.md": "n\n"}
		switch c.name {
		case "no numbered record":
			delete(files, "docs/facts/F-0004-numbered-facts.md")
		case "two numbered records":
			files["docs/facts/F-0005-numbered-facts.md"] = "n\n"
		}
		src.Tree = memTree(files)
		if _, err := buildPrompt(src); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
	}
	if c, err := buildPrompt(reopenSrc(t, "specification", "need")); err != nil || !strings.Contains(string(c.Text(rootDir)), "Facts: F-0004#3 F-0004.") {
		t.Errorf("the control: %v; want the line with the token and the numbered record", err)
	}
}

// shownLine gives the line that a filled instruction shows for the class, with its placeholders filled.
func shownLine(text, class string) string {
	prefix := map[string]string{"need": "- need or constraint: \"", "measure": "- measure: \""}[class]
	for _, l := range strings.Split(text, "\n") {
		if rest, ok := strings.CutPrefix(l, prefix); ok {
			line, _, _ := strings.Cut(rest, "\"")
			return strings.NewReplacer("<the requirement, one sentence>", "The system works.", "<the acceptance criterion>", "a check of it", "<the success criterion>", "it is fast").Replace(line)
		}
	}
	return ""
}

func reopenedLine(t *testing.T, class string) string {
	t.Helper()
	c, err := buildPrompt(reopenSrc(t, "specification", class))
	if err != nil {
		t.Fatal(err)
	}
	if line := shownLine(string(c.Text(rootDir)), class); line != "" {
		return line
	}
	t.Fatalf("the instruction shows no line for a %s", class)
	return ""
}

// U6: the lines that the instruction shows, written under the heading, pass the full check at a post-bet head: the fact 3 (a need) and the fact 4
// (a measure) have no requirement but the line; the control without the line, or with a line that cites another fact, fails.
func TestTheReopenedFormPassesTheFullCheck(t *testing.T) {
	in := specimenInput(t, "")
	in.Spans = [][]string{{"0", "10", "1", "need", "", "s1"}, {"10", "20", "2", "measure", "", "s2"}, {"20", "30", "3", "context", "", "s3"}, {"30", "40", "4", "context", "", "s4"}}
	ver := versionBytes([][]string{{"1", "need", "in", "9"}, {"2", "measure", "in", "9"}, {"3", "need", "in", "9"}, {"4", "measure", "in", "9"}})
	in.Index = [][]string{{"1", fmt.Sprintf("%x", sha256.Sum256(ver)), "M2e"}}
	in.Versions, in.PSBRecords, in.Records = map[int][]byte{1: ver}, []string{"F-0004"}, []string{"F-0003", "F-0004"}
	with := func(lines ...string) spec.FullInput {
		out := in
		out.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(specimen(t, -1) + "\n## Waits for a bet\n\n" + strings.Join(lines, "\n") + "\n")}
		return out
	}
	need := reopenedLine(t, "need")
	measure := strings.ReplaceAll(reopenedLine(t, "measure"), "F-0004#12", "F-0004#4")
	if !strings.HasPrefix(need, "REQ-006: ") || !strings.HasPrefix(measure, "Success criterion: ") {
		t.Fatalf("the lines the instruction shows: %q and %q", need, measure)
	}
	rows := spec.CheckFull(with(need, measure))
	for _, r := range rows {
		if r[1] == "fail" && r[0] != "classes" && r[0] != "ran" { // the copy of the bet is not in this input
			t.Errorf("the rule %s fails on the lines of the instruction: %s", r[0], r[2])
		}
	}
	for _, c := range []struct{ name, rule string }{{"no line", ""}, {"no line for the need", "needs"}, {"no line for the measure", "measures"}} {
		var lines []string
		switch c.rule {
		case "needs":
			lines = []string{measure}
		case "measures":
			lines = []string{need}
		}
		failed := false
		for _, r := range spec.CheckFull(with(lines...)) {
			failed = failed || (r[1] == "fail" && (c.rule == "" && (r[0] == "needs" || r[0] == "measures") || r[0] == c.rule))
		}
		if !failed {
			t.Errorf("%s: no rule fails; want the rule needs or measures", c.name)
		}
	}
}

// U8: the Reopening of the full input is read from the records: the facts of the selection of the latest version against Intake, less those whose
// row names a task that merged or the task whose head is judged; no index is none, an absent record is no task, a record the reader refuses is an error.
func TestTheReopeningSetOfTheFullInputIsReadFromTheRecords(t *testing.T) {
	ctx := context.Background()
	build := func(reopenings string) (*Sessions, spec.FullInput) {
		files := cavFiles(t)
		spans, _ := records.ReadSpans(files["spec/spans.tsv"])
		qs, _ := records.ReadQuestions(files["questions.tsv"])
		as, _ := records.ReadAnswers(files["intake/answers.tsv"])
		var version [][]string
		for _, f := range factRows(spans, qs, as) {
			class := f[1]
			if f[0] == "1" {
				class = "constraint" // Intake has fact 1 as a need
			}
			version = append(version, []string{f[0], class, "in", "9"})
		}
		v := versionBytes(version)
		files["spec/inventory-1.tsv"], files["spec/inventory.tsv"] = v, indexBytes(t, []string{"1", fmt.Sprintf("%x", sha256.Sum256(v)), "M2e"})
		if reopenings != "" {
			files[reopeningsPath] = []byte(reoHeader + reopenings)
		}
		s := cavSessions(files)
		in, err := cavRecords(t, s)
		if err != nil {
			t.Fatal(err)
		}
		return s, in
	}
	const row = "1\t1\tconstraint\tF-0003#1\tNFR-002\t%s\t9\n"
	for _, c := range []struct {
		name, reopenings, task string
		want                   []string
	}{
		{"no record yet", "", "T-cccc", []string{"1"}},
		{"its own task, not merged", fmt.Sprintf(row, "T-bbbb"), "T-bbbb", nil},
		{"another task, the task of the fact not merged", fmt.Sprintf(row, "T-bbbb"), "T-cccc", []string{"1"}},
		{"another task, the task of the fact merged", fmt.Sprintf(row, "T-aaaa"), "T-cccc", nil},
	} {
		s, in := build(c.reopenings)
		if err := s.specReopening(ctx, &in, c.task); err != nil || !slices.Equal(in.Reopening, c.want) {
			t.Errorf("%s: Reopening = %v, %v; want %v", c.name, in.Reopening, err, c.want)
		}
	}
	s, in := build(fmt.Sprintf(row, "T-bbbb"))
	s.Store.(*standinStore).files[reopeningsPath] = []byte(reoHeader + "1\t1\tconstraint\tF-0003#1\tREQ-002\tT-bbbb\t9\n")
	if err := s.specReopening(ctx, &in, "T-cccc"); err == nil || !strings.Contains(err.Error(), reopeningsPath) {
		t.Errorf("a record the reader refuses: %v; want an error that names %s", err, reopeningsPath)
	}
	in = spec.FullInput{}
	if err := s.specReopening(ctx, &in, "T-cccc"); err != nil {
		t.Errorf("before the first bet: %v; want none", err)
	}
}
