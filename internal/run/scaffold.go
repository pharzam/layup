package run

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// The blocks setup-tree and forge-check of milestone M2d (docs/spec/records.md,
// NFR-001 — The records of the Scaffold; task T-tetz, #270): the schemas of
// setup/tree.tsv and setup/forge-check.tsv, which this package alone writes and
// reads, and the rules that their blocks give in words, as the records of
// Intake do in internal/records ("— when" is "exactly when"). What reads
// another file or the forge is the writer's: for setup-tree the step verify,
// row 55 (that its last row is the head of layup-setup and that each commit has
// its tree), and for forge-check the step forge-check, row 57a (that the time
// rises and that each row has the run of the table).

// SetupTreeSchema is the form of setup/tree.tsv: the block setup-tree.
var SetupTreeSchema = tsv.Schema{Name: "setup-tree", Location: "records:setup/tree.tsv", Columns: []tsv.Column{
	{Name: "order", Type: "int", Key: true},
	{Name: "commit", Type: "sha1"},
	{Name: "tree", Type: "sha1"},
	{Name: "step", Type: "text"},
}}

// ForgeCheckSchema is the form of setup/forge-check.tsv: the block forge-check.
var ForgeCheckSchema = tsv.Schema{Name: "forge-check", Location: "records:setup/forge-check.tsv", Columns: []tsv.Column{
	{Name: "check", Type: "text", Key: true},
	{Name: "result", Type: "enum(pass|fail)"},
	{Name: "detail", Type: "text"},
	{Name: "run", Type: "text"},
	{Name: "at", Type: "time"},
}}

// rowFields gives the value of each column of a row of s by its name, as
// tsv.Read gives it ("" for the empty value).
func rowFields(s tsv.Schema, r []string) map[string]string {
	m := map[string]string{}
	for i, c := range s.Columns {
		if i < len(r) {
			m[c.Name] = r[i]
		}
	}
	return m
}

// emptyColumn gives the first of names whose value is the empty value, or "".
func emptyColumn(f map[string]string, names ...string) string {
	for _, n := range names {
		if f[n] == "" {
			return n
		}
	}
	return ""
}

// readTable reads data by the schema s, checks each row with row, then the
// table with table; a broken rule is a *tsv.Error with its line and column.
func readTable(data []byte, s tsv.Schema, row func([]string) error, table func([][]string) error) ([][]string, error) {
	rows, err := tsv.Read(data, s)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		if err := row(r); err != nil {
			var re *records.RowError
			if errors.As(err, &re) {
				return nil, &tsv.Error{Line: i + 2, Column: re.Column, Reason: re.Reason}
			}
			return nil, err
		}
	}
	if err := table(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// setupStepForm is the step of a setup commit: S04 to S14, from its message
// "chore: setup <step>". The steps are not in order in the table: S14 commits
// before S11 to S13.
var setupStepForm = regexp.MustCompile(`^S(0[4-9]|1[0-4])$`)

// CheckSetupCommit checks the rules of the block setup-tree on one row: no
// column holds the empty value, and the step is root exactly for the order 0,
// else S04 to S14. A broken rule is a *records.RowError.
func CheckSetupCommit(r []string) error {
	f := rowFields(SetupTreeSchema, r)
	if c := emptyColumn(f, "commit", "tree", "step"); c != "" {
		return &records.RowError{Column: c, Reason: "this column never holds the empty value"}
	}
	switch {
	case (f["order"] == "0") != (f["step"] == "root"):
		return &records.RowError{Column: "step", Reason: "the step is root exactly when the order is 0"}
	case f["order"] != "0" && !setupStepForm.MatchString(f["step"]):
		return &records.RowError{Column: "step", Reason: "the step of a setup commit is S04 to S14"}
	}
	return nil
}

// CheckSetupTree checks the order of setup/tree.tsv, which reads more than one
// row: the order runs 0, 1, 2, ... in the order of the rows, so the first row
// is the root commit, and a table with no row is refused. A broken rule is a
// *tsv.Error.
func CheckSetupTree(rows [][]string) error {
	if len(rows) == 0 {
		return &tsv.Error{Line: 1, Column: "order", Reason: "no row: the first row is the root commit, of order 0"}
	}
	for i, r := range rows {
		if want := strconv.Itoa(i); rowFields(SetupTreeSchema, r)["order"] != want {
			return &tsv.Error{Line: i + 2, Column: "order", Reason: "the order of this row is " + want + ": 0 for the root commit, then 1, 2, ... in the order of the rows"}
		}
	}
	return nil
}

// ReadSetupTree reads setup/tree.tsv by its schema, the rules of each row and
// the rule of the order.
func ReadSetupTree(data []byte) ([][]string, error) {
	return readTable(data, SetupTreeSchema, CheckSetupCommit, CheckSetupTree)
}

// forgeChecks are the checks of the table in the order of the block, one row
// each, whatever the result (scaffold.md, REQ-003, Scaffold 7): the second is
// rules: and the name of the default branch, which is not empty.
var forgeChecks = []string{"tree", "rules:<default branch>", "rules:layup-probe", "rules:layup-records", "bypass:default", "bypass:records", "probe", "activity"}

// CheckForgeCheck checks the rules of the block forge-check on one row: result,
// run and at never hold the empty value; a fail has a detail; a pass has none,
// but a pass of a rules check, which may name the extra rule types; and the run
// has the form of a run ID. The name of the check is the rule of the table.
// A broken rule is a *records.RowError.
func CheckForgeCheck(r []string) error {
	f := rowFields(ForgeCheckSchema, r)
	if c := emptyColumn(f, "result", "run", "at"); c != "" {
		return &records.RowError{Column: c, Reason: "this column never holds the empty value"}
	}
	switch {
	case f["result"] == "fail" && f["detail"] == "":
		return &records.RowError{Column: "detail", Reason: "a fail names the first missing or different rule, or the reason"}
	case f["result"] == "pass" && f["detail"] != "" && !strings.HasPrefix(f["check"], "rules:"):
		return &records.RowError{Column: "detail", Reason: "a pass has no detail, but a pass of a rules check, which may name its extra rule types"}
	case !records.ValidRunID(f["run"]):
		return &records.RowError{Column: "run", Reason: "a run ID is 16 lowercase hexadecimal characters"}
	}
	return nil
}

// CheckForgeChecks checks the checks of setup/forge-check.tsv, which reads more
// than one row: exactly the eight of the block, in its order. The default
// branch of the second row is not empty, and a name of either of the next two
// rows repeats a key, which tsv.Read refuses before this check. A broken rule
// is a *tsv.Error.
func CheckForgeChecks(rows [][]string) error {
	for i, r := range rows {
		got := rowFields(ForgeCheckSchema, r)["check"]
		switch {
		case i >= len(forgeChecks):
			return &tsv.Error{Line: i + 2, Column: "check", Reason: "the table has eight rows, one for each check"}
		case i == 1 && (!strings.HasPrefix(got, "rules:") || got == "rules:"):
			return &tsv.Error{Line: i + 2, Column: "check", Reason: "the next check is " + forgeChecks[i]}
		case i != 1 && got != forgeChecks[i]:
			return &tsv.Error{Line: i + 2, Column: "check", Reason: "the next check is " + forgeChecks[i]}
		}
	}
	if len(rows) < len(forgeChecks) {
		return &tsv.Error{Line: len(rows) + 2, Column: "check", Reason: "the next check is " + forgeChecks[len(rows)]}
	}
	return nil
}

// ReadForgeCheck reads setup/forge-check.tsv by its schema, the rules of each
// row and the rule of its checks.
func ReadForgeCheck(data []byte) ([][]string, error) {
	return readTable(data, ForgeCheckSchema, CheckForgeCheck, CheckForgeChecks)
}
