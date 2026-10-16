package run

// The mode of S13 for the code-owner rule (docs/spec/scaffold.md, The code-owner rule; task T-jew3, #356): a work
// area is built for a new target, or for one that was set up before the rule, from records that are already there.

import (
	"fmt"
	"path/filepath"

	"github.com/pharzam/layup/internal/setup"
	workpkg "github.com/pharzam/layup/internal/work"
)

// codeOwnerMode decides, from setup/record.tsv and setup/rulesets.txt of the records (nil when absent), whether the work
// area at area is built for a target set up before the rule. With no setup record the target is new; a record with the row
// codeowners.sha256 of S13 is a new target built again, whose commits are the recorded ones; a record without it is a target
// set up before the rule, and its migration is the update of the ruleset that the Operator applied, which only the record of
// the apply names. A target set up before the rule whose rulesets were never applied, or whose record gives no positive ID, is
// an error that names the recovery: the setup commits of such a target hold no .github/CODEOWNERS, so a new area cannot
// reproduce them, and no ID is guessed and nothing is applied or created.
func codeOwnerMode(name, area string, record, rulesets []byte) (*setup.Migration, error) {
	if record == nil {
		return nil, nil
	}
	rows, err := workpkg.ParseRecord(record)
	if err != nil {
		return nil, fmt.Errorf("setup/record.tsv: %w", err)
	}
	if _, ok := rows.Value("S13", "codeowners.sha256"); ok {
		return nil, nil
	}
	host := filepath.Dir(filepath.Dir(filepath.Dir(area)))
	if rulesets == nil {
		return nil, fmt.Errorf("the rulesets of %s were never applied (setup/rulesets.txt is not in the records), and its work area %s is gone; "+
			"this code cannot build that area again, because the setup commits recorded before the code-owner rule hold no .github/CODEOWNERS. "+
			"Restore %s from a copy of the host and run: layup run %s --host %s (its out/commands.sh applies both rulesets, as printed before), or start a new target with layup run --new", name, area, area, name, host)
	}
	noID := func(reason string) error {
		return fmt.Errorf("setup/rulesets.txt of %s does not give a positive id for default.json (%s); no id is guessed and no ruleset is applied or created. "+
			"Restore %s from a copy of the host, or start a new target", name, reason, area)
	}
	outs, err := readApplied(rulesets)
	if err != nil {
		return nil, noID(err.Error())
	}
	id, err := rulesetID(outs[applyFiles[0]])
	if err != nil {
		return nil, noID(err.Error())
	}
	return &setup.Migration{RulesetID: id}, nil
}

// operatorLogin gives the login of the Operator's row of approvers.tsv: the owner of the gate paths. No such row is an
// error, never an empty owner, and the repository's owner is not read.
func operatorLogin(approvers [][]string) (string, error) {
	for _, a := range approvers {
		if a[1] == "operator" {
			return a[2], nil
		}
	}
	return "", fmt.Errorf("approvers.tsv has no row of the role operator")
}

// migrateNotes gives the lines that name, for the Operator, the files of the migration of the area at area: forge-check
// prints no commands, so the build of the area is where the path is given.
func migrateNotes(area string) []string {
	return []string{
		"this target was set up before the code-owner rule: apply the update with your own login: sh " + filepath.Join(area, "out", "commands.sh"),
		filepath.Join(area, "out", "land-code-owner.sh") + " is not applicable until the tree check admits its commit (Q1 of #356, open)",
	}
}
