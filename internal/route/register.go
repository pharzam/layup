// Package route reads the host registers of the harnesses (docs/spec/records.md,
// the blocks harness-register, models and routing-register): the harness
// register of Start (row 22b of the plan, task T-1g1q, #137) with the columns
// that M2b adds, the models and the Operator's routing register, the checks
// across the three, and the check of a harness's credential file (row 30a,
// task T-ysph, #155). Admission and the routing order come in row 30b.
package route

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// HarnessRegisterSchema is the form of host:registers/harnesses.tsv: the
// columns that Start reads, and the ten that M2b adds after them.
var HarnessRegisterSchema = tsv.Schema{Name: "harness-register", Location: "host:registers/harnesses.tsv", Columns: []tsv.Column{
	{Name: "harness", Type: "id(<word>)", Key: true},
	{Name: "cap", Type: "decimal"},
	{Name: "wall", Type: "int"},
	{Name: "command", Type: "text"},
	{Name: "prompt", Type: "enum(file|arg|stdin)"},
	{Name: "version", Type: "text"},
	{Name: "credential", Type: "text"},
	{Name: "credential_to", Type: "text"},
	{Name: "rules", Type: "list(text)"},
	{Name: "policy", Type: "list(text)"},
	{Name: "usage", Type: "enum(claude-result|none)"},
	{Name: "billing", Type: "enum(api|subscription)"},
	{Name: "vars", Type: "list(text)"},
}}

// The columns of the harness register whose rule has no clause for the
// empty value (the plan review of #155, condition 1).
var harnessRequired = []string{"wall", "command", "prompt", "version", "rules", "usage", "billing"}

var (
	// nameForm is the form of NAME of var:NAME and of a fixed variable
	// (docs/spec/session.md, The environment and the harness credential).
	nameForm = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	// refused is each name that no session variable of a register may take:
	// the named list of a session's environment, and the forge credentials and
	// agent sockets that layup knows.
	refused = []string{"PATH", "LANG", "HOME", "TMPDIR", "GIT_CONFIG_NOSYSTEM",
		"GH_TOKEN", "GITHUB_TOKEN", "GH_ENTERPRISE_TOKEN", "GITHUB_ENTERPRISE_TOKEN", "SSH_AUTH_SOCK"}
)

// fields gives the value of each column of a row of s by its name.
func fields(s tsv.Schema, r []string) map[string]string {
	m := map[string]string{}
	for i, c := range s.Columns {
		if i < len(r) {
			m[c.Name] = r[i]
		}
	}
	return m
}

// words reports whether v is words separated by one space, with no space at
// either end.
func words(v string) bool {
	return v == strings.TrimSpace(v) && !strings.Contains(v, "  ")
}

// checkHarness checks the rules of the block harness-register on one row; a
// broken rule is the column and why.
func checkHarness(f map[string]string) (string, string) {
	for _, c := range harnessRequired {
		if f[c] == "" {
			return c, "this column never holds the empty value"
		}
	}
	if f["wall"] == "0" {
		return "wall", "the wall-clock limit is 1 minute or more"
	}
	switch {
	case strings.Contains(f["command"], "{cap}") != (f["cap"] != ""):
		return "command", "the command holds {cap} exactly when cap is not the empty value"
	case !words(f["command"]):
		return "command", "the words of the command are separated by one space"
	case !words(f["version"]):
		return "version", "the words of the version command are separated by one space"
	case f["credential"] != "" && !filepath.IsAbs(f["credential"]):
		return "credential", "the credential is an absolute path"
	case (f["credential"] == "") != (f["credential_to"] == ""):
		return "credential_to", "credential_to is the empty value exactly when credential is"
	}
	credVar := ""
	if to := f["credential_to"]; to != "" {
		if name, ok := strings.CutPrefix(to, "var:"); ok {
			if !nameForm.MatchString(name) || slices.Contains(refused, name) {
				return "credential_to", "var:NAME takes a name of the form of a variable, not of the named list or of a forge credential"
			}
			credVar = name
		} else if path, ok := strings.CutPrefix(to, "file:"); ok {
			if path == "" || filepath.IsAbs(path) || slices.Contains(strings.Split(path, "/"), "..") {
				return "credential_to", "file:PATH takes a relative path with no part .."
			}
		} else {
			return "credential_to", "credential_to is var:NAME or file:PATH"
		}
	}
	for _, p := range strings.Fields(f["policy"]) {
		if !filepath.IsAbs(p) {
			return "policy", "each policy path is absolute"
		}
	}
	for _, v := range strings.Fields(f["vars"]) {
		name, _, ok := strings.Cut(v, "=")
		switch {
		case !ok || !nameForm.MatchString(name):
			return "vars", "each fixed variable is NAME=VALUE, NAME of the form of a variable"
		case slices.Contains(refused, name) || name == credVar:
			return "vars", name + " is a name of the named list, of the credential, or of a forge credential or an agent socket"
		}
	}
	return "", ""
}

// ReadHarnesses reads the harness register by its schema and the rules of its
// block, and gives its rows and the harness IDs in their order; a register
// with no row is allowed (docs/spec/run.md, Input states). An error names its
// line and column.
func ReadHarnesses(data []byte) ([][]string, []string, error) {
	rows, err := tsv.Read(data, HarnessRegisterSchema)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]string, 0, len(rows))
	for i, r := range rows {
		if column, reason := checkHarness(fields(HarnessRegisterSchema, r)); column != "" {
			return nil, nil, &tsv.Error{Line: i + 2, Column: column, Reason: reason}
		}
		ids = append(ids, r[0])
	}
	return rows, ids, nil
}
