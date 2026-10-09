package route

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

const harnessHeader = "harness\tcap\twall\tcommand\tprompt\tversion\tcredential\tcredential_to\trules\tpolicy\tusage\tbilling\tvars"

// The rows of a valid register: a harness with a cap, a credential by var: and
// a fixed variable; and one with no cap, a credential by file: and a policy.
const (
	claudeRow = "claude\t10.0\t60\tclaude -p --model {model} --max-budget-usd {cap} {prompt}\targ\tclaude --version\t/home/u/.layup/claude.key\tvar:ANTHROPIC_API_KEY\tCLAUDE.md\t—\tclaude-result\tapi\tANTHROPIC_DEFAULT_HAIKU_MODEL=claude-opus-5-5"
	devinRow  = "devin\t—\t30\tdevin -p --model {model} {prompt}\tfile\tdevin --version\t/home/u/.layup/devin.toml\tfile:.config/devin/credentials.toml\tAGENTS.md .devin/rules.md\t/etc/devin/policy.json\tnone\tsubscription\t—"
	plainRow  = "codex\t—\t20\tcodex exec --model {model}\tstdin\tcodex --version\t—\t—\tAGENTS.md\t—\tnone\tapi\t—"
)

func harnesses(rows ...string) []byte {
	return []byte(strings.Join(append([]string{harnessHeader}, rows...), "\n") + "\n")
}

// field gives row with its field of column (by the header) set to value.
func field(header, row, column, value string) string {
	names := strings.Split(header, "\t")
	f := strings.Split(row, "\t")
	for i, n := range names {
		if n == column {
			f[i] = value
		}
	}
	return strings.Join(f, "\t")
}

// refusedAt reports a test error unless err is a *tsv.Error in column.
func refusedAt(t *testing.T, name string, err error, column string) {
	t.Helper()
	var te *tsv.Error
	switch {
	case err == nil:
		t.Errorf("%s: read, want an error in column %s", name, column)
	case !errors.As(err, &te):
		t.Errorf("%s: %v, want a *tsv.Error in column %s", name, err, column)
	case te.Column != column:
		t.Errorf("%s: %v, want an error in column %s", name, err, column)
	}
}

func TestAGoodHarnessRegisterIsRead(t *testing.T) {
	rows, ids, err := ReadHarnesses(harnesses(claudeRow, devinRow, plainRow, field(harnessHeader, claudeRow, "harness", "fable")))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ids, []string{"claude", "devin", "codex", "fable"}) || len(rows) != 4 || rows[1][1] != "" || rows[0][2] != "60" {
		t.Fatalf("ids %v, rows %v", ids, rows)
	}
	if rows, ids, err := ReadHarnesses(harnesses()); err != nil || len(rows) != 0 || len(ids) != 0 {
		t.Errorf("an empty register: %v, %v, %v; want it read, with no row", rows, ids, err)
	}
	// A PATH of file: whose name holds two dots is no part ".." (condition 2 of
	// the plan review of #155).
	if _, _, err := ReadHarnesses(harnesses(field(harnessHeader, devinRow, "credential_to", "file:a..b"))); err != nil {
		t.Errorf("file:a..b: %v; want it read", err)
	}
}

func TestHarnessRegisterRefusesEachBrokenRule(t *testing.T) {
	for name, data := range map[string][]byte{
		"two rows for one harness": harnesses(claudeRow, claudeRow),
		"an unknown column":        []byte(harnessHeader + "\tmodel\n" + claudeRow + "\tfable\n"),
		"the header of M2a":        []byte("harness\tcap\twall\nclaude\t10.0\t60\n"),
		"a harness ID in capitals": harnesses(field(harnessHeader, claudeRow, "harness", "Claude")),
		"a prompt of another word": harnesses(field(harnessHeader, claudeRow, "prompt", "pipe")),
		"a usage of another word":  harnesses(field(harnessHeader, claudeRow, "usage", "codex-json")),
	} {
		if _, _, err := ReadHarnesses(data); err == nil || !strings.Contains(err.Error(), "line ") {
			t.Errorf("%s: %v; want an error that names its line", name, err)
		}
	}
	// On plainRow, which has no cap, so the rule of {cap} cannot refuse an
	// empty command for it (round 1 of #155, finding 2).
	for _, c := range []string{"wall", "command", "prompt", "version", "rules", "usage", "billing"} {
		_, _, err := ReadHarnesses(harnesses(field(harnessHeader, plainRow, c, "—")))
		refusedAt(t, "the empty value in "+c, err, c)
	}
	for _, c := range []struct{ name, row, column string }{
		{"a wall of 0", field(harnessHeader, claudeRow, "wall", "0"), "wall"},
		{"a cap and no {cap}", field(harnessHeader, claudeRow, "command", "claude -p --model {model} {prompt}"), "command"},
		{"{cap} and no cap", field(harnessHeader, claudeRow, "cap", "—"), "command"},
		{"two spaces in command", field(harnessHeader, claudeRow, "command", "claude  -p {cap}"), "command"},
		{"a space at the end of version", field(harnessHeader, claudeRow, "version", "claude --version "), "version"},
		{"a credential that is not absolute", field(harnessHeader, claudeRow, "credential", "claude.key"), "credential"},
		{"a credential and no credential_to", field(harnessHeader, claudeRow, "credential_to", "—"), "credential_to"},
		{"credential_to and no credential", field(harnessHeader, plainRow, "credential_to", "var:OPENAI_API_KEY"), "credential_to"},
		{"credential_to of another form", field(harnessHeader, claudeRow, "credential_to", "env:ANTHROPIC_API_KEY"), "credential_to"},
		{"var: with a name of another form", field(harnessHeader, claudeRow, "credential_to", "var:ANTHROPIC-KEY"), "credential_to"},
		{"var: with a name of the named list", field(harnessHeader, claudeRow, "credential_to", "var:PATH"), "credential_to"},
		{"var: with a forge credential's name", field(harnessHeader, claudeRow, "credential_to", "var:GH_TOKEN"), "credential_to"},
		{"file: with an absolute PATH", field(harnessHeader, devinRow, "credential_to", "file:/etc/creds"), "credential_to"},
		{"file: with a part ..", field(harnessHeader, devinRow, "credential_to", "file:a/../b"), "credential_to"},
		{"file: with .. at its end", field(harnessHeader, devinRow, "credential_to", "file:a/.."), "credential_to"},
		{"file: with an empty PATH", field(harnessHeader, devinRow, "credential_to", "file:"), "credential_to"},
		{"a policy that is not absolute", field(harnessHeader, devinRow, "policy", "etc/devin/policy.json"), "policy"},
		{"a fixed variable that is not NAME=VALUE", field(harnessHeader, claudeRow, "vars", "ANTHROPIC_DEFAULT_HAIKU_MODEL"), "vars"},
		{"a fixed variable whose name is of another form", field(harnessHeader, claudeRow, "vars", "A-B=1"), "vars"},
		{"a fixed variable of the named list", field(harnessHeader, claudeRow, "vars", "HOME=/tmp"), "vars"},
		{"a fixed variable that is the credential's", field(harnessHeader, claudeRow, "vars", "ANTHROPIC_API_KEY=x"), "vars"},
		{"a fixed variable GH_TOKEN", field(harnessHeader, claudeRow, "vars", "GH_TOKEN=x"), "vars"},
		{"a fixed variable GITHUB_TOKEN", field(harnessHeader, claudeRow, "vars", "GITHUB_TOKEN=x"), "vars"},
		{"a fixed variable GH_ENTERPRISE_TOKEN", field(harnessHeader, claudeRow, "vars", "GH_ENTERPRISE_TOKEN=x"), "vars"},
		{"a fixed variable GITHUB_ENTERPRISE_TOKEN", field(harnessHeader, claudeRow, "vars", "GITHUB_ENTERPRISE_TOKEN=x"), "vars"},
		{"a fixed variable SSH_AUTH_SOCK", field(harnessHeader, claudeRow, "vars", "SSH_AUTH_SOCK=/tmp/a"), "vars"},
	} {
		_, _, err := ReadHarnesses(harnesses(c.row))
		refusedAt(t, c.name, err, c.column)
	}
}
