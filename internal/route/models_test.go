package route

import (
	"io/fs"
	"strings"
	"testing"
)

const modelsHeader = "harness\tmodel\tcontext\tsource\tdate\tuse\treason"

const (
	opusRow  = "claude\tclaude-opus-5-5\t1000000\thttps://docs.claude.com/models\t2026-10-09T12:00:00Z\tyes\t—"
	haikuRow = "claude\tclaude-haiku-5-5\t200000\thttps://docs.claude.com/models\t2026-10-09T12:00:00Z\tno\tthe Operator's list of models not to use (engineering-discipline.md)"
	sweRow   = "devin\tswe-2-high\t200000\thttp://docs.devin.ai/models\t2026-10-09T12:00:00Z\tyes\t—"
)

func table(header string, rows ...string) []byte {
	return []byte(strings.Join(append([]string{header}, rows...), "\n") + "\n")
}

func TestAGoodModelsIsRead(t *testing.T) {
	if _, err := ReadModels(table(modelsHeader, opusRow, haikuRow, sweRow)); err != nil {
		t.Fatal(err)
	}
}

func TestModelsRefusesEachBrokenRule(t *testing.T) {
	if _, err := ReadModels(table(strings.Replace(modelsHeader, "use", "used", 1), opusRow)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"context", "source", "date", "use"} {
		_, err := ReadModels(table(modelsHeader, field(modelsHeader, opusRow, c, "—")))
		refusedAt(t, "the empty value in "+c, err, c)
	}
	for _, c := range []struct{ name, row, column string }{
		{"a model in use with a reason", field(modelsHeader, opusRow, "reason", "fast"), "reason"},
		{"a model not to use with no reason", field(modelsHeader, haikuRow, "reason", "—"), "reason"},
		{"a source that is not a URL", field(modelsHeader, opusRow, "source", "the docs"), "source"},
		{"a source of another scheme", field(modelsHeader, opusRow, "source", "ftp://docs.claude.com/models"), "source"},
	} {
		_, err := ReadModels(table(modelsHeader, c.row))
		refusedAt(t, c.name, err, c.column)
	}
}

const routingHeader = "role\ttier\tposition\tharness\tmodel"

func TestAGoodRoutingRegisterIsRead(t *testing.T) {
	// The positions of a role and tier run 1 to k in any order of the file, as
	// records:routing.tsv reads them (row 28, #153).
	rows := []string{"developer\texecution\t2\tdevin\tswe-2-high", "developer\texecution\t1\tclaude\tclaude-opus-5-5",
		"verifier\treasoning\t1\tclaude\tclaude-opus-5-5"}
	if _, err := ReadRoutingRegister(table(routingHeader, rows...)); err != nil {
		t.Fatal(err)
	}
}

func TestRoutingRegisterRefusesEachBrokenRule(t *testing.T) {
	row := "developer\texecution\t1\tclaude\tclaude-opus-5-5"
	if _, err := ReadRoutingRegister(table(strings.Replace(routingHeader, "tier", "level", 1), row)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"harness", "model"} {
		_, err := ReadRoutingRegister(table(routingHeader, field(routingHeader, row, c, "—")))
		refusedAt(t, "the empty value in "+c, err, c)
	}
	_, err := ReadRoutingRegister(table(routingHeader, row, "developer\texecution\t3\tdevin\tswe-2-high"))
	refusedAt(t, "a gap in the positions", err, "position")
}

// The checks across the three registers (session.md, Input states, row 7; the
// block routing-register: a model of that harness in models.tsv).
func TestCheckRegisters(t *testing.T) {
	ids := []string{"claude", "devin"}
	models := rowsOf(t, ReadModels, table(modelsHeader, opusRow, haikuRow, sweRow))
	routing := rowsOf(t, ReadRoutingRegister, table(routingHeader, "developer\texecution\t1\tclaude\tclaude-opus-5-5", "developer\texecution\t2\tdevin\tswe-2-high"))
	if err := CheckRegisters(ids, models, routing); err != nil {
		t.Fatal(err)
	}
	for name, c := range map[string]struct {
		ids             []string
		models, routing [][]string
		file, column    string
	}{
		"a models row of a harness that the register lacks":   {[]string{"claude"}, models, routing[:1], "models.tsv", "harness"},
		"a routing row of a harness that the register lacks":  {ids, models, append(routing, []string{"developer", "execution", "3", "codex", "gpt-6"}), "routing.tsv", "harness"},
		"a routing row whose model has no row of its harness": {ids, models, append(routing, []string{"developer", "execution", "3", "devin", "claude-opus-5-5"}), "routing.tsv", "model"},
	} {
		err := CheckRegisters(c.ids, c.models, c.routing)
		if err == nil || !strings.HasPrefix(err.Error(), c.file+": ") {
			t.Errorf("%s: %v; want an error of %s", name, err, c.file)
			continue
		}
		refusedAt(t, name, err, c.column)
	}
}

func rowsOf(t *testing.T, read func([]byte) ([][]string, error), data []byte) [][]string {
	t.Helper()
	rows, err := read(data)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// The checks of a credential file on its mode and owner, with a stand-in of
// both (session.md, Input states, row 3); the real file is the integration test's.
func TestCheckCredentialFile(t *testing.T) {
	const path = "/home/u/.layup/claude.key"
	if err := checkCredentialFile(path, 0o600, 1000, 1000); err != nil {
		t.Fatalf("mode 0600 of the runner: %v", err)
	}
	for name, c := range map[string]struct {
		mode          fs.FileMode
		owner, runner int
	}{
		"mode 0644":           {0o644, 1000, 1000},
		"mode 0400":           {0o400, 1000, 1000},
		"a directory":         {fs.ModeDir | 0o600, 1000, 1000},
		"another user's file": {0o600, 0, 1000},
	} {
		if err := checkCredentialFile(path, c.mode, c.owner, c.runner); err == nil || !strings.Contains(err.Error(), path) {
			t.Errorf("%s: %v; want an error that names the file", name, err)
		}
	}
	if err := CheckCredential("claude.key"); err == nil || !strings.Contains(err.Error(), "claude.key") {
		t.Errorf("a relative path: %v; want an error that names it", err)
	}
	if err := CheckCredential("/nonexistent/claude.key"); err == nil || !strings.Contains(err.Error(), "/nonexistent/claude.key") {
		t.Errorf("a missing file: %v; want an error that names it", err)
	}
}
