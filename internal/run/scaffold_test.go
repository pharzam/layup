package run

import (
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The blocks setup-tree and forge-check of the Scaffold (task T-tetz, #270):
// one valid file per reader, one case per rule that its block gives in words,
// one case per column whose rule has no clause for the empty value, and a wrong
// header. Each refusal must name the line, the column and the reason of its
// rule (refusedBy), because the form of a column refuses its empty value too,
// so a case that another rule refuses at the same column fails no mutation of
// its own.

const (
	setupTreeHeader  = "order\tcommit\ttree\tstep"
	forgeCheckHeader = "check\tresult\tdetail\trun\tat"
	scaffoldRunID    = "0123456789abcdef"
	scaffoldNever    = "never holds the empty value"
)

// scaffoldTable gives a record of the header and the rows.
func scaffoldTable(header string, rows ...string) []byte {
	return []byte(strings.Join(append([]string{header}, rows...), "\n") + "\n")
}

// scaffoldField gives row with its field of column (by the header) set to value.
func scaffoldField(header, row, column, value string) string {
	f := strings.Split(row, "\t")
	for i, n := range strings.Split(header, "\t") {
		if n == column {
			f[i] = value
		}
	}
	return strings.Join(f, "\t")
}

// refusedBy reports a test error unless read fails on a *tsv.Error of line and
// column whose reason holds the text reason.
func refusedBy(t *testing.T, name string, read func([]byte) ([][]string, error), data []byte, line int, column, reason string) {
	t.Helper()
	_, err := read(data)
	var te *tsv.Error
	switch {
	case err == nil:
		t.Errorf("%s: read, want an error of line %d, column %q", name, line, column)
	case !errors.As(err, &te):
		t.Errorf("%s: %v, want a *tsv.Error of line %d, column %q", name, err, line, column)
	case te.Line != line || te.Column != column || !strings.Contains(te.Reason, reason):
		t.Errorf("%s: %v, want line %d, column %q and a reason with %q", name, err, line, column, reason)
	}
}

// treeRow is a row of setup-tree: its commit is the character c, 40 times, and its tree the next.
func treeRow(order string, c byte, step string) string {
	return order + "\t" + strings.Repeat(string(c), 40) + "\t" + strings.Repeat(string(c+1), 40) + "\t" + step
}

func TestAValidSetupTreeIsRead(t *testing.T) {
	for name, data := range map[string][]byte{
		"the root and four commits, S14 before S11": scaffoldTable(setupTreeHeader, treeRow("0", 'a', "root"), treeRow("1", 'b', "S04"),
			treeRow("2", 'c', "S05"), treeRow("3", 'd', "S14"), treeRow("4", 'e', "S11")),
		"the root alone": scaffoldTable(setupTreeHeader, treeRow("0", 'a', "root")),
	} {
		if _, err := ReadSetupTree(data); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestSetupTreeRefusesEachBrokenRule(t *testing.T) {
	root, one := treeRow("0", 'a', "root"), treeRow("1", 'b', "S04")
	read := func(rows ...string) []byte { return scaffoldTable(setupTreeHeader, rows...) }
	with := func(row, column, value string) string { return scaffoldField(setupTreeHeader, row, column, value) }
	refusedBy(t, "the empty value in order", ReadSetupTree, read(with(root, "order", "—")), 2, "order", "key")
	for _, col := range []string{"commit", "tree", "step"} {
		refusedBy(t, "the empty value in "+col, ReadSetupTree, read(with(root, col, "—")), 2, col, scaffoldNever)
	}
	refusedBy(t, "the root with the step of a commit", ReadSetupTree, read(with(root, "step", "S04")), 2, "step", "root exactly when the order is 0")
	for _, step := range []string{"S03", "S15", "S4", "S004", "s05", "setup"} {
		refusedBy(t, "a commit of the step "+step, ReadSetupTree, read(root, with(one, "step", step)), 3, "step", "S04 to S14")
	}
	refusedBy(t, "a first row that is not the root", ReadSetupTree, read(one), 2, "order", "the order of this row is 0")
	refusedBy(t, "an order that skips a number", ReadSetupTree, read(root, with(one, "order", "2")), 3, "order", "the order of this row is 1")
	refusedBy(t, "a table with no row", ReadSetupTree, read(), 1, "order", "no row")
}

func TestSetupTreeRefusesAWrongHeader(t *testing.T) {
	if _, err := ReadSetupTree(scaffoldTable(strings.Replace(setupTreeHeader, "\t", "\tx", 1), treeRow("0", 'a', "root"))); err == nil {
		t.Error("setup-tree: a wrong header: read, want an error")
	}
}

// forgeRows are the eight checks of the table, in their order, each a pass,
// with the default branch named branch.
func forgeRows(branch string) []string {
	var rows []string
	for _, c := range []string{"tree", "rules:" + branch, "rules:layup-probe", "rules:layup-records", "bypass:default", "bypass:records", "probe", "activity"} {
		rows = append(rows, c+"\tpass\t—\t"+scaffoldRunID+"\t2026-10-13T04:00:00Z")
	}
	return rows
}

// forgeWith gives the rows with the field of column of row i set to value.
func forgeWith(rows []string, i int, column, value string) []byte {
	out := append([]string{}, rows...)
	out[i] = scaffoldField(forgeCheckHeader, out[i], column, value)
	return scaffoldTable(forgeCheckHeader, out...)
}

func TestAValidForgeCheckIsRead(t *testing.T) {
	rows, f40 := forgeRows("main"), strings.Repeat("f", 40)
	extra := forgeWith(rows, 1, "detail", "required_signatures")
	failed := append([]string{}, rows...)
	failed[0] = scaffoldField(forgeCheckHeader, scaffoldField(forgeCheckHeader, rows[0], "result", "fail"), "detail", "the head of main is "+f40)
	failed[6] = scaffoldField(forgeCheckHeader, scaffoldField(forgeCheckHeader, rows[6], "result", "fail"), "detail", "probe accepted: "+f40)
	for name, data := range map[string][]byte{
		"eight passes":                         scaffoldTable(forgeCheckHeader, rows...),
		"a default branch of another name":     scaffoldTable(forgeCheckHeader, forgeRows("trunk")...),
		"a pass of rules with its extra types": extra,
		"two fails, each with its reason":      scaffoldTable(forgeCheckHeader, failed...),
	} {
		if _, err := ReadForgeCheck(data); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
}

func TestForgeCheckRefusesEachBrokenRule(t *testing.T) {
	rows := forgeRows("main")
	table := func(rows ...string) []byte { return scaffoldTable(forgeCheckHeader, rows...) }
	refusedBy(t, "the empty value in check", ReadForgeCheck, forgeWith(rows, 0, "check", "—"), 2, "check", "key")
	for _, col := range []string{"result", "run", "at"} {
		refusedBy(t, "the empty value in "+col, ReadForgeCheck, forgeWith(rows, 0, col, "—"), 2, col, scaffoldNever)
	}
	refusedBy(t, "a fail with no reason", ReadForgeCheck, forgeWith(rows, 3, "result", "fail"), 5, "detail", "a fail names")
	refusedBy(t, "a pass of tree with a detail", ReadForgeCheck, forgeWith(rows, 0, "detail", "the head"), 2, "detail", "a pass has no detail")
	refusedBy(t, "a pass of bypass with a detail", ReadForgeCheck, forgeWith(rows, 4, "detail", "the bypass"), 6, "detail", "a pass has no detail")
	refusedBy(t, "a pass of probe with a detail", ReadForgeCheck, forgeWith(rows, 6, "detail", "refused"), 8, "detail", "a pass has no detail")
	for _, id := range []string{"0123456789abcde", "0123456789abcdef0", "0123456789ABCDEF", "0123456789abcdeg"} {
		refusedBy(t, "a run "+id, ReadForgeCheck, forgeWith(rows, 2, "run", id), 4, "run", "16 lowercase hexadecimal")
	}
	refusedBy(t, "a table of seven checks", ReadForgeCheck, table(rows[:7]...), 9, "check", "the next check is activity")
	refusedBy(t, "a table of nine checks", ReadForgeCheck, table(append(append([]string{}, rows...), strings.Replace(rows[7], "activity", "audit", 1))...), 10, "check", "eight rows")
	refusedBy(t, "a table with no row", ReadForgeCheck, table(), 2, "check", "the next check is tree")
	refusedBy(t, "a check that is not the first", ReadForgeCheck, forgeWith(rows, 0, "check", "probe2"), 2, "check", "the next check is tree")
	refusedBy(t, "two checks that swap", ReadForgeCheck, table(rows[0], rows[1], rows[2], rows[3], rows[5], rows[4], rows[6], rows[7]), 6, "check", "the next check is bypass:default")
	refusedBy(t, "the second check with no branch", ReadForgeCheck, forgeWith(rows, 1, "check", "rules:"), 3, "check", "the next check is rules:")
	refusedBy(t, "the second check that is not a rules check", ReadForgeCheck, forgeWith(rows, 1, "check", "bypass:main"), 3, "check", "the next check is rules:")
	refusedBy(t, "a default branch named like the next check: its key repeats", ReadForgeCheck, forgeWith(rows, 1, "check", "rules:layup-probe"), 4, "", "repeats the key")
}

func TestForgeCheckRefusesAWrongHeader(t *testing.T) {
	if _, err := ReadForgeCheck(scaffoldTable(strings.Replace(forgeCheckHeader, "\t", "\tx", 1), forgeRows("main")...)); err == nil {
		t.Error("forge-check: a wrong header: read, want an error")
	}
}
