//go:build e2e

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// write writes a file of a work area, with its directories.
func write(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}

// The usage and input errors of layup setup give exit 2, nothing on standard
// output, and the reason on standard error: a work area with no problem
// statement too (D5 of #86).
func TestSetup(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing")
	unasked := t.TempDir() // an answer to a question that S01 does not ask (finding 4 of round 1)
	write(t, filepath.Join(unasked, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	write(t, filepath.Join(unasked, "inputs", "answers.tsv"), "question\tanswer\tby\tsource\tquestion_text\nS01-wrong\tx\toperator\tu\t—\n")
	nobrief := t.TempDir()
	for _, c := range []struct {
		args          []string
		reason, usage string
	}{
		{[]string{"setup"}, "missing argument WORK", "\nusage: layup "},
		{[]string{"setup", "--x", "y", "w"}, `unknown flag "--x"`, "\nusage: layup "},
		{[]string{"setup", missing}, "the work area " + missing + " is not a directory", ""},
		{[]string{"setup", nobrief}, "inputs/briefs/problem-statement.md: the problem statement is absent", ""},
		{[]string{"setup", unasked}, "inputs/answers.tsv: line 2, S01-wrong: S01 does not ask this question", ""},
	} {
		r := layup(t, c.args...)
		if r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: "+c.reason+"\n") || (c.usage == "") == strings.Contains(r.stderr, "usage:") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing and %q", c.args, r.code, r.stdout, r.stderr, c.reason)
		}
	}
}

// The demos of #86, #90 and #92: on a new work area, S01 stops with one table
// of its four questions (exit 3), the same on each run; with the answers, the
// next run does S01 to S06 on the stand-in baseline by its file:// URL and the
// prose step stops with one table of its inputs; with the inputs, the next
// run does S07 to S14, S12 with the gates of the Go entry, and S15 stops for
// out/verify.tsv (exit 3); a last run gives the same table and starts at S15.
// The whole setup, with S15 done, is the e2e of row 16 of the plan.
func TestSetupRunsS01ToS15(t *testing.T) {
	env := goEnv(t)
	tmp := t.TempDir()
	url, _, err := standin.Baseline(filepath.Join(tmp, "baseline"))
	if err != nil {
		t.Fatal(err)
	}
	w := filepath.Join(tmp, "work")
	write(t, filepath.Join(w, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	stop := "step\tquestion\task\twhere\n"
	for _, q := range work.S01Questions {
		stop += "S01\t" + q.ID + "\t" + q.Text + "\t\u2014\n"
	}
	if r := repeat(t, "setup", w); r.code != 3 || r.stdout != stop || !strings.Contains(r.stderr, "layup setup: [1/15] S01\n") {
		t.Fatalf("a new work area: exit %d, stdout %q, stderr %q; want 3 and the stop table\n%s", r.code, r.stdout, r.stderr, stop)
	}
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	var b bytes.Buffer
	if err := tsv.Write(&b, work.AnswersSchema, [][]string{{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", url, "operator", at, ""}}); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(w, "inputs", "answers.tsv"), b.String())
	prose := "step\tquestion\task\twhere\n"
	for _, r := range [][]string{{"S14", "AGENTS.md"}, {"S14", "README.md"}, {"S08", "docs/glossary.md"}, {"S09", "docs/guardrails.md"}, {"S07", "docs/onboarding-for-engineers.md"}} {
		prose += r[0] + "\tF-" + r[1] + "\tWrite the text of " + r[1] + " for the target, and give it as inputs/files/" + r[1] + ".\t" + r[1] + "\n"
	}
	if r := layup(t, "setup", w); r.code != 3 || r.stdout != prose || !strings.Contains(r.stderr, "layup setup: [6/15] S06\n") {
		t.Fatalf("the run with the answers: exit %d, stdout %q, stderr %q; want 3 and the stop table of the prose step\n%s", r.code, r.stdout, r.stderr, prose)
	}
	for p, text := range map[string]string{
		"docs/onboarding-for-engineers.md": "# Onboarding\n\nThe work starts from [the problem statement](facts/problem-statement-brief.md).\n",
		"docs/glossary.md":                 "# Glossary\n",
		"docs/guardrails.md":               "# Guardrails\n",
		"README.md":                        "# " + standin.Name + "\n\nSet up from [the pin](docs/setup/armature.pin); the records are on the branch `layup-records`.\n",
		"AGENTS.md":                        "# AGENTS.md\n\nAgent context for **" + standin.Name + "**.\n",
	} {
		write(t, filepath.Join(w, "inputs", "files", filepath.FromSlash(p)), text)
	}
	r := layupWith(t, env, "setup", w)
	verify := "Run layup setup verify '" + w + "' > '" + filepath.Join(w, "out", "verify.tsv") + "', then run layup setup '" + w + "' again."
	if want := "step\tquestion\task\twhere\nS15\tO-verify\t" + verify + "\t\u2014\n"; r.code != 3 || r.stdout != want ||
		!strings.Contains(r.stderr, "layup setup: [7/9] S12\n") || !strings.Contains(r.stderr, "layup setup: [8/9] S13\n") {
		t.Fatalf("the run with the inputs: exit %d, stdout %q, stderr %q; want 3, S07 to S14 done and the stop of S15\n%s", r.code, r.stdout, r.stderr, want)
	}
	record, err := work.ReadRecord(w)
	for _, id := range []string{"S07", "S08", "S09", "S10", "S11", "S12", "S13", "S14"} {
		if _, ok := record.Value(id, "done"); err != nil || !ok {
			t.Errorf("the record has no done row of %s: %v", id, err)
		}
	}
	again := layupWith(t, env, "setup", w)
	if again.code != 3 || again.stdout != r.stdout || again.stderr != "layup setup: [1/1] S15\n" {
		t.Errorf("a last run: exit %d, stdout %q, stderr %q; want the same table, from S15", again.code, again.stdout, again.stderr)
	}
}
