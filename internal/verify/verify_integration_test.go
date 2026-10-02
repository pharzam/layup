//go:build integration

package verify

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

func isolate(t *testing.T) {
	t.Helper()
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1"} {
		t.Setenv(k, v)
	}
}

func noSteps(i, n int, check string) func() { return func() {} }

// state gives the refs of WORK/target and the hashes of the files of WORK/out
// and WORK/inputs: a run of layup setup verify changes none of them.
func state(t *testing.T, w string) string {
	t.Helper()
	out, err := exec.Command("git", "-C", filepath.Join(w, work.TargetPath), "for-each-ref", "--format=%(refname) %(objectname)").Output()
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, dir := range []string{"out", "inputs"} {
		filepath.WalkDir(filepath.Join(w, dir), func(p string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				data, _ := os.ReadFile(p)
				s += fmt.Sprintf("%x %s\n", sha256.Sum256(data), p)
			}
			return err
		})
	}
	return s
}

// The demo of #84: on a stand-in work area the three built checks pass, each
// other check is not-active, and each kind of the manifest has its row. The
// run changes no ref and no file of the work area, and removes its scratch
// tree.
func TestRunOnAStandInWorkArea(t *testing.T) {
	isolate(t)
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	before := state(t, w.Dir)
	var scratch string
	saved := tempDir
	tempDir = func() (string, error) { d, err := saved(); scratch = d; return d, err }
	t.Cleanup(func() { tempDir = saved })
	tbl, err := Run(w.Dir, noSteps, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, r := range tbl.Rows {
		if r.Result != "not-active" || r.Reason != "not built yet" {
			got = append(got, r.Check+" "+r.Result+" "+r.Reason)
		}
	}
	if want := []string{"discipline-tests pass ", "pin pass ", "kit-history pass ", "facts pass ", "onboarding pass ", "glossary pass ", "guardrails pass ",
		"markers pass ", "adapted pass ", "identity pass ", "link-lint pass ", "sources pass "}; strings.Join(got, "|") != strings.Join(want, "|") || len(tbl.Rows) != 15 ||
		tbl.Rows[13].Check != "gate:static" || tbl.Rows[14].Check != "gate:layout" {
		t.Errorf("the rows %q; want %q, the others not built yet, and gate:static and gate:layout last", tbl.Rows, want)
	}
	if after := state(t, w.Dir); after != before {
		t.Errorf("the run changed the work area:\n%s\nto\n%s", before, after)
	}
	if _, err := os.Stat(scratch); scratch == "" || !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("the scratch directory %q is still there (%v)", scratch, err)
	}
	if out, _ := exec.Command("git", "-C", filepath.Join(w.Dir, work.TargetPath), "worktree", "list").Output(); strings.Count(string(out), "\n") != 1 {
		t.Errorf("git worktree list:\n%s\nwant the main work tree only", out)
	}
}

// Each finding of the target's part of a check (D3 of #84) fails its row, on
// a real work area.
func TestEachFindingOfATarget(t *testing.T) {
	isolate(t)
	dir := t.TempDir()
	link := filepath.ToSlash(filepath.Join(dir, "baseline")) + "/"
	for _, c := range []struct {
		name    string
		o       standin.Options
		check   string
		reason  string
		answers string // rows added to inputs/answers.tsv; "=" and a text replaces the file
	}{
		{"a decision of the kit", standin.Options{Files: map[string]string{"docs/decisions/D-0002.md": "x\n"}}, "kit-history",
			"decisions: docs/decisions/ exists (kit step 4 deletes it)", ""},
		{"a link to the baseline of the target", standin.Options{Files: map[string]string{"docs/tasks/backlog.md": "- [#1](file://" + link + "issues/1)\n"}}, "kit-history",
			"kit-link: docs/tasks/backlog.md links " + link + " (a kit task or note)", ""},
		{"a link to the repository of the baseline (round 1, finding 1)", standin.Options{Files: map[string]string{"docs/tasks/backlog.md": "- [the baseline](file://" + strings.TrimSuffix(link, "/") + ")\n"}},
			"kit-history", "kit-link: docs/tasks/backlog.md links " + link + " (a kit task or note)", ""},
		{"no pin", standin.Options{Files: map[string]string{"docs/setup/armature.pin": ""}}, "pin", "missing: docs/setup/armature.pin is absent", ""},
		{"a record of another commit", standin.Options{Record: map[string]string{"S02 pin.commit": strings.Repeat("1", 40)}}, "pin", "commit: the pin names ", ""},
		{"a README.md with no name", standin.Options{Files: map[string]string{"README.md": "# x\n\n[pin](docs/setup/armature.pin)\n"}}, "identity",
			"name: README.md does not hold the name of the target, " + standin.Name, ""},
		{"a README.md that names no records branch", standin.Options{Files: map[string]string{"README.md": "# " + standin.Name + "\n\n[pin](docs/setup/armature.pin)\n"}}, "identity",
			"branch: README.md does not name the branch layup-records", ""},
		{"the phrase of the kit", standin.Options{Files: map[string]string{"AGENTS.md": "Agent context for **Armature**\n"}}, "identity",
			`kit: AGENTS.md says the repository is the Armature kit ("Agent context for **Armature**")`, ""},
		// The four checks in a target's form (#89), on a record whose steps
		// are done through S14 (D10).
		{"a brief that does not match its hash", standin.Options{Done: true, Files: map[string]string{standin.BriefPath: "another text\n"}}, "facts",
			"hash: docs/facts/problem-statement-brief.md does not match docs/setup/facts.sha256", ""},
		{"an answer with no fact in the record of S04", standin.Options{Done: true}, "facts",
			"answers: Q-001 of inputs/answers.tsv is not a fact of F-0001", "Q-001\tten\tidea-owner\tu\t\u2014\n"},
		{"an M- answer after S11 and no record of S11", standin.Options{Done: true}, "facts",
			"record: expected one docs/facts/F-NNNN-marker-answers.md", "M-0123abcd\t8080\toperator\tu\t\u2014\n"},
		{"no onboarding file", standin.Options{Done: true, Files: map[string]string{"docs/onboarding-for-engineers.md": ""}}, "onboarding",
			"missing: docs/onboarding-for-engineers.md is absent", ""},
		{"an onboarding file with no link", standin.Options{Done: true, Files: map[string]string{"docs/onboarding-for-engineers.md": "# Onboarding\n"}}, "onboarding",
			"link: docs/onboarding-for-engineers.md has no link to facts/problem-statement-brief.md", ""},
		{"a glossary citation that does not resolve", standin.Options{Done: true, Files: map[string]string{"docs/glossary.md": "| A | a (`F-0001#9`) |\n"}}, "glossary",
			"fact: F-0001#9 is not a fact of the F-0001 record", ""},
		// The four checks of #87.
		{"an unlisted marker", standin.Options{Files: map[string]string{"docs/m.md": "a \u2039x\u203a marker\n"}}, "markers",
			"unlisted: docs/m.md \u2039x\u203a", ""},
		{"an answer row that the record names and the answers do not have", standin.Options{}, "sources",
			"source: S01 stack: the answer S01-stack is not a row of inputs/answers.tsv",
			"=question\tanswer\tby\tsource\tquestion_text\nS01-name\tstand-in-target\toperator\tu\t\u2014\nS01-visibility\tpublic\toperator\tu\t\u2014\nS01-baseline\tx\toperator\tu\t\u2014\n"},
		{"discipline tests that fail", standin.Options{Files: map[string]string{"docs/tests/run-discipline-tests.sh": "exit 1\n"}}, "discipline-tests",
			"exit 1", ""},
		{"a guardrails entry with no valid check", standin.Options{Done: true, Files: map[string]string{"docs/guardrails.md": "- **Inv-1** a rule. Check: later\n"}}, "guardrails",
			`check: Inv-1 has no valid Check: value ("later")`, ""},
	} {
		if err := os.RemoveAll(dir); err != nil {
			t.Fatal(err)
		}
		w, err := standin.Make(dir, c.o)
		if err != nil {
			t.Fatal(err)
		}
		if text, ok := strings.CutPrefix(c.answers, "="); ok {
			if err := os.WriteFile(filepath.Join(w.Dir, filepath.FromSlash(work.AnswersPath)), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		} else if c.answers != "" {
			f, err := os.OpenFile(filepath.Join(w.Dir, filepath.FromSlash(work.AnswersPath)), os.O_APPEND|os.O_WRONLY, 0)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := f.WriteString(c.answers); err != nil || f.Close() != nil {
				t.Fatal(err)
			}
		}
		tbl, err := Check(w.Dir, []string{c.check}, noSteps, io.Discard)
		if err != nil || len(tbl.Rows) != 1 || tbl.Rows[0].Result != "fail" || !strings.HasPrefix(tbl.Rows[0].Reason, c.reason) {
			t.Errorf("%s: %q, %v; want %s fail, %s", c.name, tbl.Rows, err, c.check, c.reason)
		}
	}
	// A baseline script that is not there is not-active, never pass (D5 of
	// #87, NFR-004).
	if err := os.RemoveAll(dir); err != nil {
		t.Fatal(err)
	}
	w, err := standin.Make(dir, standin.Options{Files: map[string]string{"docs/links/link-lint.sh": ""}})
	if err != nil {
		t.Fatal(err)
	}
	if tbl, err := Check(w.Dir, []string{"link-lint"}, noSteps, io.Discard); err != nil || len(tbl.Rows) != 1 || tbl.Rows[0] != (Row{"link-lint", "not-active", "missing: docs/links/link-lint.sh"}) {
		t.Errorf("no link-lint.sh: %q, %v; want link-lint not-active, missing: docs/links/link-lint.sh", tbl.Rows, err)
	}
}

// The input errors of a real work area: no answers, no target, no manifest.
func TestTheInputErrorsOfAWorkArea(t *testing.T) {
	isolate(t)
	for _, c := range []struct {
		name   string
		change func(w string) error
		want   string
	}{
		{"no answers", func(w string) error { return os.Remove(filepath.Join(w, work.AnswersPath)) }, "inputs/answers.tsv: "},
		{"no target", func(w string) error { return os.RemoveAll(filepath.Join(w, work.TargetPath)) }, "target: no commit at the branch layup-setup: "},
	} {
		w, err := standin.Make(t.TempDir(), standin.Options{})
		if err != nil {
			t.Fatal(err)
		}
		if err := c.change(w.Dir); err != nil {
			t.Fatal(err)
		}
		var in *InputError
		if _, err := Run(w.Dir, noSteps, io.Discard); !errors.As(err, &in) || !strings.HasPrefix(err.Error(), c.want) {
			t.Errorf("%s: %v; want an input error that starts with %q", c.name, err, c.want)
		}
	}
	w, err := standin.Make(t.TempDir(), standin.Options{Files: map[string]string{"docs/gates.tsv": ""}})
	if err != nil {
		t.Fatal(err)
	}
	var in *InputError
	if _, err := Run(w.Dir, noSteps, io.Discard); !errors.As(err, &in) || !strings.HasPrefix(err.Error(), "docs/gates.tsv at the head of layup-setup: ") {
		t.Errorf("no manifest: %v; want an input error", err)
	}

	// Round 1, finding 2: a TMPDIR in the work area is an input error, and the
	// run leaves the work area as it was.
	w, err = standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("TMPDIR", filepath.Join(w.Dir, "out"))
	before := state(t, w.Dir)
	if _, err := Run(w.Dir, noSteps, io.Discard); !errors.As(err, &in) || !strings.Contains(err.Error(), "set TMPDIR to a directory outside it") {
		t.Errorf("TMPDIR in WORK/out: %v; want an input error", err)
	}
	if after := state(t, w.Dir); after != before {
		t.Errorf("the run changed the work area:\n%s\nto\n%s", before, after)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if rel, err := filepath.Rel(wd, w.Dir); err != nil {
		t.Fatal(err)
	} else if _, err := Run(rel, noSteps, io.Discard); !errors.As(err, &in) {
		t.Errorf("TMPDIR in WORK/out, WORK given as the relative path %s: %v; want an input error", rel, err)
	}
}

// The Go schema of the table equals its block.
func TestTheTableSchemaEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["setup-verify"]
	if !ok {
		t.Fatal("docs/spec/ has no block setup-verify")
	}
	if err := tsv.Compare(block, TableSchema); err != nil {
		t.Error(err)
	}
}
