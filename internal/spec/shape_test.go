package spec

import (
	"slices"
	"strings"
	"testing"
)

// The Shape filter (docs/spec/loop.md, REQ-012 — Shape; task T-0pt8, #349): an
// option is dropped by a failed row, a missing row or a pass with no
// evidence; the survivors are what is left, and code lists the constraints
// that apply from the records.

// orow is a row of the block options: option, constraint, result, evidence
// ("" for —).
func orow(option, constraint, result, evidence string) []string {
	return []string{option, constraint, result, evidence}
}

// standing is the list of constraints that apply in these cases.
var standing = []string{"R-003", "F-2", "Inv-1", "ADR-0010"}

// goodOptions gives the rows of the option name: a pass with evidence for each
// constraint of standing.
func goodOptions(option string) [][]string {
	var rows [][]string
	for _, c := range standing {
		rows = append(rows, orow(option, c, "pass", "docs/adr/0010.md"))
	}
	return rows
}

func names(os []Option) []string {
	var out []string
	for _, o := range os {
		out = append(out, o.Name)
	}
	return out
}

func TestAnOptionWithAFailedRowIsDropped(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")...)
	rows[5] = orow("b", "F-2", "fail", "")
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("survivors %q, want a", got)
	}
	b := r.Options[1]
	if !b.Dropped || len(b.Reasons) != 1 || b.Reasons[0] != "F-2: fail" {
		t.Errorf("option b: dropped %v, reasons %q, want dropped with F-2: fail", b.Dropped, b.Reasons)
	}
}

func TestAnOptionWithAMissingRowIsDropped(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")[:3]...) // b has no row for ADR-0010
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("survivors %q, want a", got)
	}
	if b := r.Options[1]; len(b.Reasons) != 1 || b.Reasons[0] != "ADR-0010: missing" {
		t.Errorf("option b reasons %q, want ADR-0010: missing", b.Reasons)
	}
}

func TestAPassWithNoEvidenceDropsTheOption(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")...)
	rows[6] = orow("b", "Inv-1", "pass", "")
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("survivors %q, want a", got)
	}
	if b := r.Options[1]; len(b.Reasons) != 1 || b.Reasons[0] != "Inv-1: pass with no evidence" {
		t.Errorf("option b reasons %q, want Inv-1: pass with no evidence", b.Reasons)
	}
}

func TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints(t *testing.T) {
	rows := append(goodOptions("a"), orow("b", "ADR-0010", "fail", ""), orow("b", "R-003", "pass", ""))
	r := Filter(rows, standing)
	want := []string{"R-003: pass with no evidence", "F-2: missing", "Inv-1: missing", "ADR-0010: fail"}
	if got := r.Options[1].Reasons; !slices.Equal(got, want) {
		t.Errorf("reasons %q, want %q", got, want)
	}
}

func TestOneSurvivorIsChosen(t *testing.T) {
	rows := append(goodOptions("a"), orow("b", "R-003", "fail", ""))
	r := Filter(rows, standing)
	name, ok := r.Chosen()
	if !ok || name != "a" {
		t.Errorf("chosen %q %v, want a true", name, ok)
	}
}

func TestMoreThanOneSurvivorChoosesNone(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")...)
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a", "b"}) {
		t.Fatalf("survivors %q, want a b", got)
	}
	if name, ok := r.Chosen(); ok {
		t.Errorf("chosen %q, want none: the bet decides, with no recommendation", name)
	}
}

func TestEveryOptionDroppedLeavesNoSurvivor(t *testing.T) {
	rows := [][]string{orow("a", "R-003", "fail", ""), orow("b", "F-2", "fail", "")}
	r := Filter(rows, standing)
	if len(r.Survivors) != 0 {
		t.Errorf("survivors %q, want none", names(r.Survivors))
	}
	if name, ok := r.Chosen(); ok {
		t.Errorf("chosen %q, want none", name)
	}
	for _, o := range r.Options {
		if !o.Dropped {
			t.Errorf("option %s survives", o.Name)
		}
	}
}

func TestAMissingRowDropsOnlyTheOptionThatLacksIt(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")[:2]...) // a has all four rows, b the first two
	r := Filter(rows, standing)
	if a, b := r.Options[0], r.Options[1]; a.Dropped || !b.Dropped || len(b.Reasons) != 2 {
		t.Errorf("a dropped %v (%q), b dropped %v (%q): want a kept, b dropped for Inv-1 and ADR-0010", a.Dropped, a.Reasons, b.Dropped, b.Reasons)
	}
}

func TestAnExtraRowWithAPassAndNoEvidenceDropsItsOption(t *testing.T) {
	rows := append(goodOptions("a"), goodOptions("b")...)
	rows = append(rows, orow("b", "R-010", "pass", ""))
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("survivors %q, want a", got)
	}
	if b := r.Options[1]; len(b.Reasons) != 1 || b.Reasons[0] != "R-010: pass with no evidence" {
		t.Errorf("b reasons %q, want R-010: pass with no evidence", b.Reasons)
	}
}

func TestAnExtraConstraintIsJudgedButNotDemanded(t *testing.T) {
	// a requirement that the architect names is judged: a fail drops the option;
	// an option with no row for it is not dropped.
	rows := append(goodOptions("a"), goodOptions("b")...)
	rows = append(rows, orow("b", "R-010", "fail", ""))
	r := Filter(rows, standing)
	if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
		t.Fatalf("survivors %q, want a (b fails R-010, which a has no row for)", got)
	}
	if a := r.Options[0]; a.Dropped {
		t.Errorf("a is dropped for %q: no row is demanded for R-010", a.Reasons)
	}
}

// draft is a requirement draft (the block draft: id, kind, statement, covers,
// criterion, session).
func draft() [][]string {
	return [][]string{
		{"R-001", "requirement", "a need", "1", "the need is met", "S-00000001"},
		{"R-003", "non-functional", "a limit", "2", "stays under 5 s", "S-00000001"},
		{"R-002", "success-criterion", "a measure", "3", "7 of 10", "S-00000001"},
		{"R-004", "non-functional", "another limit", "2", "no network", "S-00000001"},
	}
}

func TestTheNonFunctionalRequirementsOfTheDraftApplyWithTheirCriterion(t *testing.T) {
	got := Standing(RequirementsOfDraft(draft()), nil, nil, nil)
	want := []Constraint{{"R-003", "stays under 5 s"}, {"R-004", "no network"}}
	if !slices.Equal(got, want) {
		t.Errorf("constraints %v, want %v: the kinds requirement and success-criterion do not apply", got, want)
	}
}

func TestTheNonFunctionalRequirementsOfThePRDApplyToo(t *testing.T) {
	prd := []Requirement{{"REQ-1", "functional", "a need"}, {"NFR-2", "non-functional", "git is the record"}}
	got := Standing(prd, nil, nil, nil)
	if want := []Constraint{{"NFR-2", "git is the record"}}; !slices.Equal(got, want) {
		t.Errorf("constraints %v, want %v", got, want)
	}
}

func TestTheConstraintFactsOfTheSpansApplyBeforeTheFirstBet(t *testing.T) {
	facts := []Fact{{"1", "need"}, {"2", "constraint"}, {"3", "measure"}, {"4", "constraint"}, {"5", "context"}}
	if got := IDs(Standing(nil, facts, nil, nil)); !slices.Equal(got, []string{"F-2", "F-4"}) {
		t.Errorf("constraints %q, want F-2 F-4: only the facts of class constraint", got)
	}
}

func TestTheConfirmedInventoryReplacesTheFactsOfTheSpans(t *testing.T) {
	spans := []Fact{{"2", "constraint"}, {"4", "constraint"}}
	inventory := []Fact{{"2", "need"}, {"4", "constraint"}, {"6", "constraint"}} // a bet changed fact 2 and an answer gave fact 6
	if got := IDs(Standing(nil, spans, nil, nil)); !slices.Equal(got, []string{"F-2", "F-4"}) {
		t.Fatalf("before the bet %q, want F-2 F-4", got)
	}
	if got := IDs(Standing(nil, inventory, nil, nil)); !slices.Equal(got, []string{"F-4", "F-6"}) {
		t.Errorf("after the bet %q, want F-4 F-6: the caller passes the confirmed inventory's facts, not the spans'", got)
	}
}

func TestTheInvariantsApply(t *testing.T) {
	if got := IDs(Standing(nil, nil, []string{"Inv-2", "Inv-1"}, nil)); !slices.Equal(got, []string{"Inv-2", "Inv-1"}) {
		t.Errorf("constraints %q, want Inv-2 Inv-1, in the order given", got)
	}
}

func TestOnlyTheAcceptedDecisionRecordsApply(t *testing.T) {
	adrs := []ADR{{"ADR-0011", "Accepted"}, {"ADR-0012", "Proposed"}, {"ADR-0013", "Superseded by ADR-0014"}, {"ADR-0010", "Accepted"}}
	if got := IDs(Standing(nil, nil, nil, adrs)); !slices.Equal(got, []string{"ADR-0011", "ADR-0010"}) {
		t.Errorf("constraints %q, want ADR-0011 ADR-0010: a record that is not Accepted does not apply", got)
	}
}

func TestTheFourKindsComeInOneOrder(t *testing.T) {
	got := IDs(Standing(RequirementsOfDraft(draft()), []Fact{{"2", "constraint"}}, []string{"Inv-1"}, []ADR{{"ADR-0010", "Accepted"}}))
	want := []string{"R-003", "R-004", "F-2", "Inv-1", "ADR-0010"}
	if !slices.Equal(got, want) {
		t.Errorf("constraints %q, want %q: requirements, facts, invariants, records", got, want)
	}
}

func TestTheFactsOfASpanTableAreItsFactRows(t *testing.T) {
	spans := [][]string{
		row("0", "7", "", "", "a heading", "# Brief"),
		row("9", "36", "1", "need", "", "x"),
		row("40", "50", "2", "constraint", "", "y"),
	}
	got := FactsOfSpans(spans)
	want := []Fact{{"1", "need"}, {"2", "constraint"}}
	if !slices.Equal(got, want) {
		t.Errorf("facts %v, want %v", got, want)
	}
}

func TestAnOptionsTableThatBreaksAnOwnerRuleIsRefused(t *testing.T) {
	head := "option\tconstraint\tresult\tevidence\n"
	for _, c := range []struct{ name, body, want string }{
		{"a fail with evidence", "a\tR-003\tfail\tdocs/x.md\n", "line 2, column evidence"},
		{"an option name with a capital", "A\tR-003\tpass\tdocs/x.md\n", "line 2, column option"},
		{"an option name with a space", "a b\tR-003\tpass\tdocs/x.md\n", "option"},
		{"an option name that is —", "—\tR-003\tpass\tdocs/x.md\n", "option"},
		{"a result outside the two", "a\tR-003\tmaybe\tdocs/x.md\n", "result"},
		{"an empty result", "a\tR-003\t—\t—\n", "result"},
		{"a row twice", "a\tR-003\tpass\tdocs/x.md\na\tR-003\tfail\t—\n", "key"},
	} {
		if _, err := ReadOptions([]byte(head + c.body)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v, want one that names %q", c.name, err, c.want)
		}
	}
	rows, err := ReadOptions([]byte(head + "a-1\tR-003\tpass\t—\na-1\tF-2\tfail\t—\n"))
	if err != nil {
		t.Fatalf("a table with a pass with no evidence and a fail with none: %v", err)
	}
	if rows[0][3] != "" || rows[1][3] != "" {
		t.Errorf("evidence %q %q, want the empty value", rows[0][3], rows[1][3])
	}
}
