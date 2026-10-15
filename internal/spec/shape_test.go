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

// orow is a row of the block options: option, constraint, result, evidence ("" for —).
func orow(option, constraint, result, evidence string) []string {
	return []string{option, constraint, result, evidence}
}

var standing = []string{"R-003", "F-2", "Inv-1", "ADR-0010"} // the constraints that apply in these cases

// goodOptions gives a pass with evidence of the option for each constraint of standing.
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

// twoOptions gives a (good) and b (good, but for the rows that change makes).
func twoOptions(change func(b [][]string) [][]string) [][]string {
	return append(goodOptions("a"), change(goodOptions("b"))...)
}

func TestEachBreakDropsOnlyItsOption(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func([][]string) [][]string
		reason string
	}{
		{"a failed row", func(b [][]string) [][]string { b[1] = orow("b", "F-2", "fail", ""); return b }, "F-2: fail"},
		{"a missing row", func(b [][]string) [][]string { return b[:3] }, "ADR-0010: missing"},
		{"a pass with no evidence", func(b [][]string) [][]string { b[2] = orow("b", "Inv-1", "pass", ""); return b }, "Inv-1: pass with no evidence"},
		{"an extra row that fails", func(b [][]string) [][]string { return append(b, orow("b", "R-010", "fail", "")) }, "R-010: fail"},
		{"an extra row, a pass with no evidence", func(b [][]string) [][]string { return append(b, orow("b", "R-010", "pass", "")) }, "R-010: pass with no evidence"},
	} {
		r := Filter(twoOptions(c.change), standing)
		if got := names(r.Survivors); !slices.Equal(got, []string{"a"}) {
			t.Errorf("%s: survivors %q, want a only", c.name, got)
		}
		if b := r.Options[1]; !b.Dropped || !slices.Equal(b.Reasons, []string{c.reason}) {
			t.Errorf("%s: b dropped %v, reasons %q, want %q", c.name, b.Dropped, b.Reasons, c.reason)
		}
	}
}

func TestAnExtraConstraintIsJudgedButNotDemanded(t *testing.T) {
	// b names R-010 and fails it; a has no row for it and is not dropped.
	r := Filter(twoOptions(func(b [][]string) [][]string { return append(b, orow("b", "R-010", "fail", "")) }), standing)
	if a := r.Options[0]; a.Dropped {
		t.Errorf("a is dropped for %q: no row is demanded for R-010", a.Reasons)
	}
}

func TestAMissingRowDropsOnlyTheOptionThatLacksIt(t *testing.T) {
	r := Filter(twoOptions(func(b [][]string) [][]string { return b[:2] }), standing) // b: R-003 and F-2 only
	if a, b := r.Options[0], r.Options[1]; a.Dropped || !slices.Equal(b.Reasons, []string{"Inv-1: missing", "ADR-0010: missing"}) {
		t.Errorf("a dropped %v, b reasons %q: want a kept, b dropped for Inv-1 and ADR-0010", a.Dropped, b.Reasons)
	}
}

func TestAnOptionWithSeveralBreaksNamesEachInTheOrderOfTheConstraints(t *testing.T) {
	r := Filter(append(goodOptions("a"), orow("b", "ADR-0010", "fail", ""), orow("b", "R-003", "pass", "")), standing)
	want := []string{"R-003: pass with no evidence", "F-2: missing", "Inv-1: missing", "ADR-0010: fail"}
	if got := r.Options[1].Reasons; !slices.Equal(got, want) {
		t.Errorf("reasons %q, want %q", got, want)
	}
}

func TestTheSurvivorsAreChosenOnlyWhenThereIsOne(t *testing.T) {
	fail := func(name string) [][]string { return [][]string{orow(name, "R-003", "fail", "")} }
	for _, c := range []struct {
		name      string
		rows      [][]string
		survivors []string
		chosen    string
	}{
		{"one survivor", append(goodOptions("a"), fail("b")...), []string{"a"}, "a"},
		{"more than one: no recommendation", append(goodOptions("a"), goodOptions("b")...), []string{"a", "b"}, ""},
		{"every option dropped", append(fail("a"), fail("b")...), nil, ""},
	} {
		r := Filter(c.rows, standing)
		name, ok := r.Chosen()
		if got := names(r.Survivors); !slices.Equal(got, c.survivors) || name != c.chosen || ok != (c.chosen != "") {
			t.Errorf("%s: survivors %q, chosen %q %v, want %q and %q", c.name, got, name, ok, c.survivors, c.chosen)
		}
	}
}

func draft() [][]string {
	return [][]string{
		{"R-001", "requirement", "a need", "1", "the need is met", "S-00000001"},
		{"R-003", "non-functional", "a limit", "2", "stays under 5 s", "S-00000001"},
		{"R-002", "success-criterion", "a measure", "3", "7 of 10", "S-00000001"},
		{"R-004", "non-functional", "another limit", "2", "no network", "S-00000001"},
	}
}

func TestTheConstraintsThatApplyComeFromTheRecordsInOneOrder(t *testing.T) {
	facts := []Fact{{"1", "need"}, {"2", "constraint"}, {"3", "measure"}, {"4", "constraint"}, {"5", "context"}}
	adrs := []ADR{{"ADR-0011", "Accepted"}, {"ADR-0012", "Proposed"}, {"ADR-0013", "Superseded by ADR-0014"}, {"ADR-0007", "Accepted. Amended by ADR-0015"}}
	prd := []Requirement{{"REQ-1", "functional", "a need"}, {"NFR-2", "non-functional", "git is the record"}}
	for _, c := range []struct {
		name string
		got  []Constraint
		want []string
	}{
		{"the draft: only the non-functional requirements", Standing(RequirementsOfDraft(draft()), nil, nil, nil), []string{"R-003", "R-004"}},
		{"the PRD: only the non-functional requirements", Standing(prd, nil, nil, nil), []string{"NFR-2"}},
		{"the facts: only those of class constraint", Standing(nil, facts, nil, nil), []string{"F-2", "F-4"}},
		{"the invariants, as given", Standing(nil, nil, []string{"Inv-2", "Inv-1"}, nil), []string{"Inv-2", "Inv-1"}},
		{"the records: the accepted and the amended, not the others", Standing(nil, nil, nil, adrs), []string{"ADR-0011", "ADR-0007"}},
		{"a fact of the answers by its own ID", Standing(nil, []Fact{{"Q-004", "constraint"}, {"6", "constraint"}}, nil, nil), []string{"Q-004", "F-6"}},
		{"the four kinds in one order", Standing(RequirementsOfDraft(draft()), facts[1:2], []string{"Inv-1"}, adrs[:1]), []string{"R-003", "R-004", "F-2", "Inv-1", "ADR-0011"}},
	} {
		if got := IDs(c.got); !slices.Equal(got, c.want) {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
	}
	if got := Standing(RequirementsOfDraft(draft()), nil, nil, nil); got[0].Text != "stays under 5 s" {
		t.Errorf("the criterion of R-003 is %q, want its criterion", got[0].Text)
	}
}

// The caller chooses the source of the facts: the span facts before the first
// bet, the confirmed inventory's after it (row 66h); Standing lists what it is given.
func TestTheCallerChoosesTheSourceOfTheFacts(t *testing.T) {
	spans := []Fact{{"2", "constraint"}, {"4", "constraint"}}
	inventory := []Fact{{"2", "need"}, {"4", "constraint"}, {"6", "constraint"}}
	if got := IDs(Standing(nil, spans, nil, nil)); !slices.Equal(got, []string{"F-2", "F-4"}) {
		t.Errorf("the span facts %q, want F-2 F-4", got)
	}
	if got := IDs(Standing(nil, inventory, nil, nil)); !slices.Equal(got, []string{"F-4", "F-6"}) {
		t.Errorf("the inventory's facts %q, want F-4 F-6", got)
	}
}

func TestTheFactsOfASpanTableAreItsFactRows(t *testing.T) {
	spans := [][]string{row("0", "7", "", "", "a heading", "# Brief"), row("9", "36", "1", "need", "", "x"), row("40", "50", "2", "constraint", "", "y")}
	if got, want := FactsOfSpans(spans), []Fact{{"1", "need"}, {"2", "constraint"}}; !slices.Equal(got, want) {
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
		{"an empty result", "a\tR-003\t—\t—\n", "line 2, column result"},
		{"a row twice", "a\tR-003\tpass\tdocs/x.md\na\tR-003\tfail\t—\n", "key"},
	} {
		if _, err := ReadOptions([]byte(head + c.body)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v, want one that names %q", c.name, err, c.want)
		}
	}
	rows, err := ReadOptions([]byte(head + "a-1\tR-003\tpass\t—\na-1\tF-2\tfail\t—\n"))
	if err != nil || rows[0][3] != "" || rows[1][3] != "" {
		t.Errorf("a pass with no evidence and a fail with none: rows %q, error %v, want the empty value and no error", rows, err)
	}
}
