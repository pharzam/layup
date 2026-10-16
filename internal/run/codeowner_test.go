package run

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/tsv"
	workpkg "github.com/pharzam/layup/internal/work"
)

// The mode of S13 for the code-owner rule (docs/spec/scaffold.md, The code-owner rule; task T-jew3, #356, plan
// review 6214007155): the work area is built for a new target, or for one set up before the rule, from the two
// records that are already there. Each case breaks one rule only, and a refusal is compared with its whole reason.

// recordText gives setup/record.tsv with the rows of S13 that name give.
func recordText(t *testing.T, names ...string) []byte {
	t.Helper()
	rows := [][]string{{"S13", "branch-protection.sha256", "a", "computed", "sha256 docs/setup/branch-protection.json"}}
	for _, n := range names {
		rows = append(rows, []string{"S13", n, "b", "computed", "sha256 x"})
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, workpkg.RecordSchema, rows); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

// appliedText gives setup/rulesets.txt with the output of the default branch's apply.
func appliedText(def string) []byte {
	return []byte("== default.json\n" + def + "\n== records.json\n{\"id\":2}\n")
}

func TestSetupAreaChoosesTheMode(t *testing.T) {
	const name, area = "acme/widget", "DIR/work/acme/widget"
	good := appliedText(`{"id":25145510,"name":"layup: the default branch"}`)
	never := "the rulesets of acme/widget were never applied (setup/rulesets.txt is not in the records), and its work area DIR/work/acme/widget is gone; " +
		"this code cannot build that area again, because the setup commits recorded before the code-owner rule hold no .github/CODEOWNERS. " +
		"Restore DIR/work/acme/widget from a copy of the host and run: layup run acme/widget --host DIR (its out/commands.sh applies both rulesets, as printed before), or start a new target with layup run --new"
	noID := func(reason string) string {
		return "setup/rulesets.txt of acme/widget does not give a positive id for default.json (" + reason + "); no id is guessed and no ruleset is applied or created. " +
			"Restore DIR/work/acme/widget from a copy of the host, or start a new target"
	}
	for _, c := range []struct {
		label    string
		record   []byte
		rulesets []byte
		want     *setup.Migration
		err      string
	}{
		{"no setup record: a new target", nil, nil, nil, ""},
		{"no setup record, rulesets there: still new", nil, good, nil, ""},
		{"recorded with the row of CODEOWNERS: new, rebuilt", recordText(t, "codeowners.sha256"), good, nil, ""},
		{"recorded with the row, rulesets never applied: new, rebuilt", recordText(t, "codeowners.sha256"), nil, nil, ""},
		{"recorded before, applied: migrate with the ID", recordText(t, "ruleset.sha256"), good, &setup.Migration{RulesetID: 25145510}, ""},
		{"recorded before, rulesets never applied", recordText(t, "ruleset.sha256"), nil, nil, never},
		{"recorded before, rulesets.txt not the record of two applies", recordText(t, "ruleset.sha256"), []byte("garbage\n"), nil, noID("it is not the record of the two applies")},
		{"recorded before, default.json not an object", recordText(t, "ruleset.sha256"), appliedText("[1]"), nil, noID("the output of the default branch's ruleset is not one JSON object")},
		{"recorded before, default.json with no id", recordText(t, "ruleset.sha256"), appliedText(`{"name":"x"}`), nil, noID("the output of the default branch's ruleset has no id")},
		{"recorded before, an id of 0", recordText(t, "ruleset.sha256"), appliedText(`{"id":0}`), nil, noID("the output of the default branch's ruleset has no id")},
		{"recorded before, an id below 0", recordText(t, "ruleset.sha256"), appliedText(`{"id":-4}`), nil, noID("the output of the default branch's ruleset has no id")},
		{"recorded before, two objects", recordText(t, "ruleset.sha256"), appliedText(`{"id":1}{"id":2}`), nil, noID("the output of the default branch's ruleset is not one JSON object")},
		{"a record that is not setup/record.tsv", []byte("not a table\n"), good, nil, "setup/record.tsv: "},
	} {
		got, err := codeOwnerMode(name, area, c.record, c.rulesets)
		if c.err == "" && (err != nil || !reflect.DeepEqual(got, c.want)) || c.err != "" && (err == nil || !strings.HasPrefix(err.Error(), c.err) || got != nil) {
			t.Errorf("%s: %+v, %v; want %+v, error %q", c.label, got, err, c.want, c.err)
		}
	}
}

// The Operator's login is the login of the row of the role operator of approvers.tsv; no such row is an error,
// never an empty owner, and the repository's owner is not read.
func TestOperatorLoginIsTheRowOfTheRole(t *testing.T) {
	rows := [][]string{{"1", "idea-owner", "someone", "t", "start"}, {"2", "operator", "pharzam", "t", "start"}}
	if got, err := operatorLogin(rows); err != nil || got != "pharzam" {
		t.Errorf("operatorLogin: %q, %v; want pharzam", got, err)
	}
	if got, err := operatorLogin(rows[:1]); err == nil || got != "" || err.Error() != "approvers.tsv has no row of the role operator" {
		t.Errorf("no row: %q, %v; want an error and no login", got, err)
	}
}

// The area built for a target set up before the rule names the files of the migration once, in the order that the Operator
// runs them, and an area for a new target names none: forge-check prints no commands, so this is the path that the Operator
// is given (task T-4cep: the tree check admits the commit of the landing).
func TestTheMigrationNamesItsFiles(t *testing.T) {
	want := []string{
		"this target was set up before the code-owner rule: apply the update with your own login: sh DIR/work/acme/widget/out/commands.sh",
		"then land the code-owner file with your own login, once the update is read back: sh DIR/work/acme/widget/out/land-code-owner.sh",
	}
	if got := migrateNotes("DIR/work/acme/widget"); !reflect.DeepEqual(got, want) {
		t.Errorf("migrateNotes: %q; want %q", got, want)
	}
}

// The inputs of the check of the code-owner commit are the record of the setup and the Operator's row of approvers.tsv, and
// nothing else (the ruling 6217133467: a target whose setup record has the row S13 codeowners.sha256 takes no such commit).
// Each case breaks one rule only, and a refusal is compared with its whole reason.
func TestLandingInputs(t *testing.T) {
	approvers := [][]string{{"1", "idea-owner", "someone", "t", "start"}, {"2", "operator", "pharzam", "t", "start"}}
	for _, c := range []struct {
		name      string
		record    []byte
		approvers [][]string
		login     string
		err       string
	}{
		{"a record without the row, the Operator's row", recordText(t, "ruleset.sha256"), approvers, "pharzam", ""},
		{"no setup record", nil, approvers, "", "setup/record.tsv is not in the records"},
		{"a record that is not setup/record.tsv", []byte("not a table\n"), approvers, "", "setup/record.tsv: "},
		{"the row S13 codeowners.sha256", recordText(t, "codeowners.sha256"), approvers, "",
			"the setup record has the row S13 codeowners.sha256: .github/CODEOWNERS is in the tree of the last setup commit, and a target set up with the rule takes no code-owner commit"},
		{"no row of the role operator", recordText(t, "ruleset.sha256"), approvers[:1], "", "approvers.tsv has no row of the role operator"},
	} {
		login, err := landingInputs(c.record, c.approvers)
		if c.err == "" && (err != nil || login != c.login) || c.err != "" && (err == nil || !strings.HasPrefix(err.Error(), c.err) || login != "") {
			t.Errorf("%s: %q, %v; want %q, error %q", c.name, login, err, c.login, c.err)
		}
	}
}
