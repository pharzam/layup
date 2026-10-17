//go:build integration

package spec

// The prd-lint seam against the real producer: `sh` runs the pinned script of
// this repository on throwaway trees in t.TempDir(), both streams are read
// separately, and the real captured lines go through the classifier — the
// tests never write into the tree (task T-4wyg, #351; review conditions
// C1–C3).

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// prdLintScript is this repository's pinned producer, the same file the
// quality gate runs.
func prdLintScript(t *testing.T) string {
	t.Helper()
	p, err := filepath.Abs(filepath.Join("..", "..", "docs", "prd", "prd-lint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// fixtureTree writes t.TempDir()/docs/{prd,facts}: the fact files so that
// prd-lint resolves F-NNNN tokens, and the PRD files given.
func fixtureTree(t *testing.T, prds map[string]string, facts []string) string {
	t.Helper()
	root := t.TempDir()
	for _, dir := range []string{"prd", "facts"} {
		if err := os.MkdirAll(filepath.Join(root, "docs", dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for name, text := range prds {
		if err := os.WriteFile(filepath.Join(root, "docs", "prd", name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range facts {
		if err := os.WriteFile(filepath.Join(root, "docs", "facts", f+".md"), []byte("# "+f+"\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return filepath.Join(root, "docs", "prd")
}

// waitingPRD is a specification that waits for its bet: every MoSCoW and
// Phase cell empty, but fact cites and the matrix complete — the only
// failure lines the producer may emit are the empty-cell ones.
func waitingPRD(moscow, phase string) string {
	return `## 6. Functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| --- | --------- | ------ | ----- | ----- |
| REQ-001 | One. | ` + moscow + ` | ` + phase + ` | F-0001#1 |

| REQ | Facts | Guardrail | ADR | Task | Test |
| --- | ----- | --------- | --- | ---- | ---- |
| REQ-001 | F-0001#1 | — | — | — | — |
`
}

func TestRunPRDLintAcceptsAWaitingSpecification(t *testing.T) {
	dir := fixtureTree(t, map[string]string{"PRD-0001-x.md": waitingPRD("", "")}, []string{"F-0001"})
	code, stdout, stderr, err := RunPRDLint(prdLintScript(t), dir)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Fatalf("the waiting specification: exit %d, want 1 (lines %v / %v)", code, stdout, stderr)
	}
	if len(stderr) != 0 {
		t.Fatalf("the waiting specification prints to stderr: %v", stderr)
	}
	if len(stdout) != 2 || !emptyCellFailure.MatchString(stdout[0]) || !emptyCellFailure.MatchString(stdout[1]) {
		t.Fatalf("the real lines are not the two empty-cell failures: %v", stdout)
	}
	ok, offend, err := WaitingSpecAccepted(prdLintScript(t), dir)
	if err != nil || !ok {
		t.Fatalf("a waiting specification is refused: %v %v", ok, offend)
	}
}

func TestRunPRDLintRefusesARealViolation(t *testing.T) {
	bad := strings.Replace(waitingPRD("", ""), "F-0001#1 |", "F-9999#1 |", 1)
	dir := fixtureTree(t, map[string]string{"PRD-0001-x.md": bad}, []string{"F-0001"})
	ok, offend, err := WaitingSpecAccepted(prdLintScript(t), dir)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("a cite to no fact record is accepted")
	}
	if len(offend) != 1 || !strings.Contains(offend[0], "no resolvable F-NNNN fact") {
		t.Fatalf("the offending lines %v, want the resolvable-fact failure", offend)
	}
}

// Both streams through the real sh (C2): a well-named waiting PRD plus a
// badly named one — the filename failure arrives on stderr while the
// empty-cell failures arrive on stdout, and the pair refuses.
func TestRunPRDLintReadsBothStreams(t *testing.T) {
	dir := fixtureTree(t, map[string]string{
		"PRD-0001-x.md": waitingPRD("", ""),
		"PRD-bad.md":    waitingPRD("", ""),
	}, []string{"F-0001"})
	code, stdout, stderr, err := RunPRDLint(prdLintScript(t), dir)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Fatalf("exit %d, want 1", code)
	}
	if len(stderr) != 1 || !strings.Contains(stderr[0], "filename must be PRD-NNNN-kebab-case.md") {
		t.Fatalf("the stderr split: %v", stderr)
	}
	// The bad name is never linted — the script's positional list holds only
	// well-formed names — so stdout is the good file's two cell lines, and the
	// refusal is the one stderr line.
	if len(stdout) != 2 || !emptyCellFailure.MatchString(stdout[0]) || !emptyCellFailure.MatchString(stdout[1]) {
		t.Fatalf("the stdout split: %v, want the good file's two cell lines", stdout)
	}
	ok, offend, err := WaitingSpecAccepted(prdLintScript(t), dir)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("a bad filename on stderr is accepted")
	}
	if len(offend) != 1 || !strings.Contains(offend[0], "filename") {
		t.Fatalf("the offending lines %v, want the filename failure", offend)
	}
}

// A directory with no PRD-*.md exits 0 with `prd-lint: OK` (C3's exit-0
// case); a — MoSCoW is no empty cell and refuses (C3); a — Phase on a
// non-Won't row is the producer's `has no phase value` and is the empty form.
func TestRunPRDLintEdgeCases(t *testing.T) {
	t.Run("an empty directory", func(t *testing.T) {
		dir := fixtureTree(t, nil, nil)
		code, stdout, _, err := RunPRDLint(prdLintScript(t), dir)
		if err != nil {
			t.Fatal(err)
		}
		if code != 0 || len(stdout) != 1 || stdout[0] != "prd-lint: OK" {
			t.Fatalf("exit %d with %v, want `prd-lint: OK`", code, stdout)
		}
		if ok, _, err := WaitingSpecAccepted(prdLintScript(t), dir); err != nil || !ok {
			t.Fatalf("exit 0 is not accepted: %v %v", ok, err)
		}
	})
	t.Run("an em-dash MoSCoW refuses", func(t *testing.T) {
		dir := fixtureTree(t, map[string]string{"PRD-0001-x.md": waitingPRD("—", "—")}, []string{"F-0001"})
		ok, offend, err := WaitingSpecAccepted(prdLintScript(t), dir)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			t.Fatal(`MoSCoW "—" is accepted as an empty cell`)
		}
		if len(offend) != 1 || !strings.Contains(offend[0], `MoSCoW "—"`) {
			t.Fatalf("the offending lines %v", offend)
		}
	})
	t.Run("a dash phase is the empty form", func(t *testing.T) {
		dir := fixtureTree(t, map[string]string{"PRD-0001-x.md": waitingPRD("", "—")}, []string{"F-0001"})
		ok, offend, err := WaitingSpecAccepted(prdLintScript(t), dir)
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("a — phase of a non-Won't row is the empty form, refused: %v", offend)
		}
	})
	t.Run("a missing script is an error", func(t *testing.T) {
		// sh exits 127 for a file it cannot open: not 0, not 1, so the verdict
		// is an error and never an acceptance.
		ok, _, err := WaitingSpecAccepted(filepath.Join(t.TempDir(), "none.sh"), t.TempDir())
		if err == nil || ok {
			t.Fatalf("a script that does not exist gives %v, %v; want an error", ok, err)
		}
	})
}
