package run

import (
	"context"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/gate"
)

const (
	headA = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	headB = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	baseX = "cccccccccccccccccccccccccccccccccccccccc"
	baseY = "dddddddddddddddddddddddddddddddddddddddd"
)

const gatesHeader = "head\tkind\tbase\tstate\tresult\treason\n"

// gateLine is one row of gates.tsv.
func gateLine(head, kind, base, state, result, reason string) string {
	return strings.Join([]string{head, kind, base, state, result, reason}, "\t") + "\n"
}

func gateRow(kind, state, result, reason string) gate.Row {
	return gate.Row{Kind: kind, State: state, Result: result, Reason: reason}
}

func TestEachRowCombinationMapsToTheGatesStatus(t *testing.T) {
	pass := gateRow("build", "active", "pass", "")
	for _, c := range []struct {
		name string
		rows []gate.Row
		want Status
	}{
		{"every kind passes", []gate.Row{pass, gateRow("static", "active", "pass", "—")}, Status{GatesContext, "success", "pass"}},
		{"a clear kind counts as a pass, with its reason", []gate.Row{pass, gateRow("docs", "pending", "clear", "pending: no product path")}, Status{GatesContext, "success", "clear: pending: no product path"}},
		{"the first clear reason", []gate.Row{gateRow("a", "active", "clear", "no product path"), gateRow("b", "pending", "clear", "pending: no product path")}, Status{GatesContext, "success", "clear: no product path"}},
		{"a failing kind", []gate.Row{pass, gateRow("go", "active", "fail", "exit 1")}, Status{GatesContext, "failure", "go: fail: exit 1"}},
		{"a not-active kind is a failure, never a success", []gate.Row{pass, gateRow("lint", "active", "not-active", "tool not found: golangci-lint")}, Status{GatesContext, "failure", "lint: not-active: tool not found: golangci-lint"}},
		{"a pending kind and a product path that changed", []gate.Row{pass, gateRow("boundary", "pending", "fail", "pending: product path changed: a.go")}, Status{GatesContext, "failure", "boundary: fail: pending: product path changed: a.go"}},
		{"the first row in the order of the manifest, past a clear one", []gate.Row{gateRow("c", "active", "clear", "no product path"), gateRow("n", "active", "not-active", "scratch tree: add failed"), gateRow("f", "active", "fail", "exit 2")}, Status{GatesContext, "failure", "n: not-active: scratch tree: add failed"}},
		{"a table with no row is no pass", nil, Status{GatesContext, "error", "error: the table has no row"}},
	} {
		if got := GatesStatus(gate.Table{Base: baseX, Head: headA, Rows: c.rows}); got != c.want {
			t.Errorf("%s: GatesStatus = %+v, want %+v", c.name, got, c.want)
		}
	}
}

func TestAGatesTableIsReadByItsRulesAndRefusesEachBrokenOne(t *testing.T) {
	good := gatesHeader + gateLine(headA, "build", baseX, "active", "pass", "—") + gateLine(headA, "docs", baseX, "pending", "clear", "pending: no product path") +
		gateLine(headB, "build", baseY, "active", "fail", "exit 1") + gateLine(headB, "docs", baseY, "pending", "clear", "pending: no product path")
	rows, err := ReadGates([]byte(good))
	if err != nil || len(rows) != 4 {
		t.Fatalf("the control, a good table of two heads: %d rows, %v", len(rows), err)
	}
	if again, err := WriteGates(rows); err != nil || string(again) != good {
		t.Errorf("WriteGates(ReadGates(good)) = %q, %v; want the same bytes", again, err)
	}
	for name, bad := range map[string]string{
		"heads interleaved":     gatesHeader + gateLine(headA, "build", baseX, "active", "pass", "—") + gateLine(headB, "build", baseX, "active", "pass", "—") + gateLine(headA, "docs", baseX, "pending", "clear", "x"),
		"two bases in a head":   gatesHeader + gateLine(headA, "build", baseX, "active", "pass", "—") + gateLine(headA, "docs", baseY, "pending", "clear", "x"),
		"a pass with a reason":  gatesHeader + gateLine(headA, "build", baseX, "active", "pass", "ok"),
		"a fail with no reason": gatesHeader + gateLine(headA, "build", baseX, "active", "fail", "—"),
		"a result of none":      gatesHeader + gateLine(headA, "build", baseX, "active", "skipped", "x"),
		"a head that is no SHA": gatesHeader + gateLine("abc", "build", baseX, "active", "pass", "—"),
		"a wrong header":        "head\tkind\tbase\tstate\tresult\twhy\n" + gateLine(headA, "build", baseX, "active", "pass", "—"),
	} {
		if _, err := ReadGates([]byte(bad)); err == nil {
			t.Errorf("%s: ReadGates accepted the table", name)
		}
	}
	if _, err := WriteGates([][]string{{headA, "build", baseX, "active", "pass", "ok"}}); err == nil {
		t.Error("WriteGates wrote a pass with a reason")
	}
	if rows, err := ReadGates([]byte(gatesHeader)); err != nil || len(rows) != 0 {
		t.Errorf("a header alone: %d rows, %v; want an empty table", len(rows), err)
	}
}

const manifest = "kind\tstate\ttool\tcommand\tscope\tconfig\nbuild\tactive\tsh\ttrue\t./*.go\t—\ndocs\tpending\tsh\ttrue\t./*.md\t—\n"

func TestASavedTableIsCompleteOnlyWhenItsKindsAreTheManifestOfItsBase(t *testing.T) {
	ok := [][]string{{headA, "build", baseX, "active", "pass", ""}, {headA, "docs", baseX, "pending", "clear", "pending: no product path"}}
	if err := gatesComplete(ok, []byte(manifest)); err != nil {
		t.Fatalf("the control, every kind in order: %v", err)
	}
	for name, rows := range map[string][][]string{
		"a header alone":                         nil,
		"an omitted last kind":                   ok[:1],
		"an omitted first kind":                  ok[1:],
		"the kinds in the wrong order":           {ok[1], ok[0]},
		"a kind that is not in the manifest":     {ok[0], ok[1], {headA, "extra", baseX, "active", "pass", ""}},
		"a state that differs from the manifest": {ok[0], {headA, "docs", baseX, "active", "clear", "x"}},
	} {
		if err := gatesComplete(rows, []byte(manifest)); err == nil {
			t.Errorf("%s: the table was complete", name)
		}
	}
	if err := gatesComplete(ok, []byte("kind\tstate\n")); err == nil {
		t.Error("a manifest that does not read: the table was complete")
	}
}

func TestTheGatesTableOfAHeadRebuildsTheFindings(t *testing.T) {
	store := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	s := newSessions(store)
	if _, found, err := s.GatesTable(context.Background(), "T-ab12", headA); err != nil || found {
		t.Fatalf("no gates.tsv: found=%v, %v; want not found and no error", found, err)
	}
	store.files[gatesPath("T-ab12")] = []byte(gatesHeader)
	if _, found, err := s.GatesTable(context.Background(), "T-ab12", headA); err != nil || found {
		t.Fatalf("a header alone: found=%v, %v; want not found", found, err)
	}
	store.files[gatesPath("T-ab12")] = []byte(gatesHeader + gateLine(headA, "build", baseX, "active", "pass", "—") + gateLine(headA, "docs", baseX, "pending", "fail", "pending: product path changed: a.md"))
	s.manifestAt = func(base string) ([]byte, error) { return []byte(manifest), nil }
	tbl, found, err := s.GatesTable(context.Background(), "T-ab12", headA)
	if err != nil || !found || tbl.Base != baseX || tbl.Head != headA || len(tbl.Rows) != 2 {
		t.Fatalf("GatesTable = %+v, %v, %v; want the two rows of the head", tbl, found, err)
	}
	if f := GateFindings(tbl); len(f) != 1 || f[0].Rule != "gate:docs" || f[0].Text != "pending: product path changed: a.md" {
		t.Errorf("GateFindings of the table = %+v; want the one finding gate:docs", f)
	}
	if _, found, err := s.GatesTable(context.Background(), "T-ab12", headB); err != nil || found {
		t.Errorf("another head: found=%v, %v; want not found", found, err)
	}
	s.manifestAt = func(string) ([]byte, error) { return nil, nil }
	if _, _, err := s.GatesTable(context.Background(), "T-ab12", headA); err == nil {
		t.Error("a base with no manifest: no error")
	}
}

// specRows is the ten rows of the full check as CheckFull gives them: every rule
// passes, then the changes of the case.
func specRows(change map[string][2]string) [][]string {
	var rows [][]string
	for _, rule := range []string{"covers", "needs", "measures", "criteria", "classes", "inventory", "quotes", "sections", "tasks", "ran"} {
		r := []string{rule, "pass", ""}
		if c, ok := change[rule]; ok {
			r = []string{rule, c[0], c[1]}
		}
		rows = append(rows, r)
	}
	return rows
}

func TestTheSpecStatusIsSuccessOnlyWithNoFailRowAndACompleteTable(t *testing.T) {
	whole := map[string][2]string{}
	for _, r := range specRows(nil) {
		whole[r[0]] = [2]string{"clear", "no specification yet"}
	}
	for _, c := range []struct {
		name string
		rows [][]string
		want Status
	}{
		{"every rule passes", specRows(nil), Status{SpecContext, "success", "pass"}},
		{"the whole table is clear", specRows(whole), Status{SpecContext, "success", "clear: no specification yet"}},
		{"classes and inventory clear before the first bet", specRows(map[string][2]string{"classes": {"clear", "before the first bet"}, "inventory": {"clear", "before the first bet"}}), Status{SpecContext, "success", "clear: before the first bet"}},
		{"the first fail in the order of the rules", specRows(map[string][2]string{"needs": {"fail", "REQ-001 has no need"}, "quotes": {"fail", "x"}}), Status{SpecContext, "failure", "needs: REQ-001 has no need"}},
		{"a rule that did not run", specRows(map[string][2]string{"sections": {"fail", "did not run: the spec files"}, "ran": {"fail", "did not run: sections"}}), Status{SpecContext, "failure", "sections: did not run: the spec files"}},
		{"a fail past a clear row", specRows(map[string][2]string{"covers": {"clear", "x"}, "tasks": {"fail", "T-aa11 has no predecessor"}}), Status{SpecContext, "failure", "tasks: T-aa11 has no predecessor"}},
		{"a header alone", nil, Status{SpecContext, "error", "error: the table of the spec check is not complete"}},
		{"an omitted rule", append(specRows(nil)[:3:3], specRows(nil)[4:]...), Status{SpecContext, "error", "error: the table of the spec check is not complete"}},
		{"a missing ran", specRows(nil)[:9], Status{SpecContext, "error", "error: the table of the spec check is not complete"}},
		{"the rules in the wrong order", append([][]string{specRows(nil)[1], specRows(nil)[0]}, specRows(nil)[2:]...), Status{SpecContext, "error", "error: the table of the spec check is not complete"}},
	} {
		if got := SpecStatus(c.rows); got != c.want {
			t.Errorf("%s: SpecStatus = %+v, want %+v", c.name, got, c.want)
		}
	}
}

func TestASpecTableIsReadByItsRules(t *testing.T) {
	store := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	s := newSessions(store)
	ctx := context.Background()
	put := func(parts ...[2]string) {
		var b strings.Builder
		b.WriteString("head\trule\tresult\tdetail\n")
		for _, p := range parts {
			for _, r := range specRows(nil) {
				b.WriteString(p[0] + "\t" + r[0] + "\t" + r[1] + "\t" + map[bool]string{true: "—", false: r[2]}[r[2] == ""] + "\n")
			}
		}
		store.files[specPath("T-ab12")] = []byte(b.String())
	}
	if _, found, err := s.SpecTable(ctx, "T-ab12", headA); err != nil || found {
		t.Fatalf("no spec-check.tsv: %v, %v", found, err)
	}
	put([2]string{headA}, [2]string{headB})
	rows, found, err := s.SpecTable(ctx, "T-ab12", headB)
	if err != nil || !found || len(rows) != 10 || rows[0][0] != "covers" || len(rows[0]) != 3 {
		t.Fatalf("SpecTable = %v, %v, %v; want the ten rows (rule, result, detail) of the head", rows, found, err)
	}
	if got := SpecStatus(rows); got.State != "success" {
		t.Errorf("the status of the rebuilt table: %+v", got)
	}
	whole := string(store.files[specPath("T-ab12")])
	lines := strings.Split(strings.TrimSuffix(whole, "\n"), "\n")
	for name, bad := range map[string][]string{
		"an omitted rule":   append(cloneLines(lines[:4]), lines[5:]...),
		"heads interleaved": append(append(cloneLines(lines[:6]), lines[11]), lines[6:11]...),
		"a head cut short":  lines[:15],
	} {
		store.files[specPath("T-ab12")] = []byte(strings.Join(bad, "\n") + "\n")
		if _, _, err := s.SpecTable(ctx, "T-ab12", headA); err == nil {
			t.Errorf("%s: SpecTable accepted the table", name)
		}
	}
}

func cloneLines(l []string) []string { return append([]string(nil), l...) }
