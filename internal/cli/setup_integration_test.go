//go:build integration

package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
	"github.com/pharzam/layup/internal/work"
)

// stubSteps gives the steps of a setup: each done, unless out names another
// outcome (condition 4 of the plan review of #85).
func stubSteps(t *testing.T, out map[string]setup.Outcome) {
	t.Helper()
	saved := setupSteps
	setupSteps = func(setup.Brief, setup.Calls) map[string]setup.Step {
		m := setup.Stubs()
		for id, s := range m {
			o, ok := out[id]
			if !ok {
				o = setup.Outcome{Kind: setup.Done, Evidence: "ok " + id}
			}
			s.Run = func(setup.Input) setup.Outcome { return o }
			if id == "S13" {
				s.Commands = func(work.Record) []setup.Command {
					return []setup.Command{{Order: 2, Comment: "the push of layup-setup", Text: "git -C target push origin layup-setup:main"}}
				}
			}
			m[id] = s
		}
		return m
	}
	t.Cleanup(func() { setupSteps = saved })
}

// On a real work area, a run reads out/record.tsv, resumes after its done rows,
// and gives each exit code: 3 with only the stop table, 0 on the resumed run,
// 1 for a step that did not pass, 2 for a record of another form.
func TestSetupExitCodesOnAWorkArea(t *testing.T) {
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	stubSteps(t, map[string]setup.Outcome{"S03": {Kind: setup.Stop, Stops: []setup.StopRow{{Step: "S03", Question: "O-push", Ask: "push the root commit"}}}})
	if code, out, _ := run("setup", w.Dir); code != 3 || out != "step\tquestion\task\twhere\nS03\tO-push\tpush the root commit\t—\n" {
		t.Errorf("a stop: exit %d, stdout %q; want 3 and only the stop table", code, out)
	}
	stubSteps(t, nil)
	if code, out, errOut := run("setup", w.Dir); code != 0 || !strings.Contains(out, "\nS01\tlayup-setup\tdone\tevery answer present; the stack has a catalog entry\n") ||
		!strings.Contains(out, "\nS13\tlayup-setup\toperator\thanded to the Operator: ok S13\n") || strings.Contains(errOut, "] S01\n") {
		t.Errorf("the resumed run: exit %d, stdout %q, stderr %q; want 0, S01 and S02 from the record and not run again", code, out, errOut)
	}
	w2, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	stubSteps(t, map[string]setup.Outcome{"S04": {Kind: setup.Fail, Evidence: "the pin differs"}})
	if code, out, _ := run("setup", w2.Dir); code != 1 || !strings.Contains(out, "\nS05\tlayup-setup\tnot-active\tnot run: S04 did not pass\n") {
		t.Errorf("a fail: exit %d, stdout %q; want 1", code, out)
	}
	if err := os.WriteFile(filepath.Join(w2.Dir, filepath.FromSlash(work.RecordPath)), []byte("step\tname\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := run("setup", w2.Dir); code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: out/record.tsv: line 1") {
		t.Errorf("a record of another form: exit %d, stdout %q, stderr %q; want 2, nothing and the line", code, out, errOut)
	}
	// A done S01 with no row answers.sha256 is exit 2 (finding 1 of round 1).
	w3, err := standin.Make(t.TempDir(), standin.Options{Record: map[string]string{"S01 answers.sha256": ""}})
	if err != nil {
		t.Fatal(err)
	}
	if code, out, errOut := run("setup", w3.Dir); code != 2 || out != "" || errOut != "layup: out/record.tsv: S01 is done and has no row answers.sha256\n" {
		t.Errorf("a done S01 with no hash: exit %d, stdout %q, stderr %q; want 2 and the reason", code, out, errOut)
	}
}

// adrBaseline makes a stand-in baseline at dir with the files of docs/adr/
// that S04 reads at LAYUP's pin: the records 0001 to 0008, their index, and
// adr-lint.sh, which LAYUP keeps as its root commit has it. It gives the
// file:// URL.
func adrBaseline(t *testing.T, dir string) string {
	t.Helper()
	url, _, err := standin.Baseline(dir)
	if err != nil {
		t.Fatal(err)
	}
	lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "adr", "adr-lint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "docs", "adr", "adr-lint.sh"), string(lint))
	index := "# Architecture decision records\n\n## Index\n\n| ADR | Title | Status |\n| --- | ----- | ------ |\n"
	for n := 1; n <= 8; n++ {
		name := fmt.Sprintf("%04d-record-%d.md", n, n)
		writeFile(t, filepath.Join(dir, "docs", "adr", name), fmt.Sprintf("# %04d. Record %d\n\nDate: 2026-01-01\n\n## Status\n\nAccepted\n\n"+
			"## Context\n\nA context.\n\n## Decision\n\nA decision.\n\n## Consequences\n\nA consequence.\n", n, n))
		index += fmt.Sprintf("| [%04d](%s) | Record %d | Accepted |\n", n, name, n)
	}
	writeFile(t, filepath.Join(dir, "docs", "adr", "README.md"), index+"\n<!-- Add one row per ADR as you write them. Keep the newest at the bottom. -->\n")
	if err := git.Add(dir); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(dir, "the records of the baseline", git.Identity{Name: "t", Email: "t@layup.invalid", Time: time.Unix(0, 0)}); err != nil {
		t.Fatal(err)
	}
	return url
}

// gitIn runs git in dir and gives its output.
func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %q: %v", args, err)
	}
	return strings.TrimSpace(string(out))
}

// The demo of #86: on a problem statement with gaps, the first run stops at
// S01 with one table, the four questions of S01 and then the Q- rows of layup
// psb check on the same file, with their texts (condition 3 of the plan
// review). With the answers, the next run does S01 to S04: the evidence of S04
// is checks pin and facts, the baseline's own adr-lint.sh passes on the
// target, and commands.sh holds the commands of S03.
func TestSetupRunsS01ToS04(t *testing.T) {
	tmp := t.TempDir()
	url := adrBaseline(t, filepath.Join(tmp, "baseline"))
	w := filepath.Join(tmp, "work")
	brief := filepath.Join(w, "inputs", "briefs", "problem-statement.md")
	if err := os.MkdirAll(filepath.Dir(brief), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, brief, "# The brief\n\nThe product must be fast.\n\n| Goal | Measure |\n| ---- | ------- |\n| Speed | TBD |\n")
	_, gaps, _ := run("psb", "check", brief)
	want := "step\tquestion\task\twhere\n"
	for _, q := range work.S01Questions {
		want += "S01\t" + q.ID + "\t" + q.Text + "\t—\n"
	}
	var answers [][]string
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	for _, q := range [][]string{{"S01-stack", "go"}, {"S01-name", standin.Name}, {"S01-visibility", "private"}, {"S01-baseline", url}} {
		answers = append(answers, []string{q[0], q[1], "operator", at, ""})
	}
	for _, line := range strings.Split(strings.TrimSpace(gaps), "\n")[1:] {
		f := strings.Split(line, "\t")
		where := "—"
		if f[2] != "0" {
			where = "inputs/briefs/problem-statement.md:" + f[2]
		}
		want += "S01\t" + f[0] + "\t" + f[4] + "\t" + where + "\n"
		answers = append(answers, []string{f[0], "the answer to " + f[0], "idea-owner", "F-0003#5", ""})
	}
	if code, out, errOut := run("setup", w); code != 3 || out != want || len(answers) < 6 {
		t.Fatalf("the first run: exit %d, stderr %q, stdout\n%s\nwant 3 and\n%s(%d answers)", code, errOut, out, want, len(answers))
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, work.AnswersSchema, answers); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(w, "inputs", "answers.tsv"), b.String())
	code, out, errOut := run("setup", w) // S01 to S06, then the prose step stops for its inputs (task T-b3r1)
	record, _ := work.ReadRecord(w)
	for _, id := range []string{"S01", "S02", "S03", "S04", "S05", "S06"} {
		if v, _ := record.Value(id, "done"); code != 3 || v == "" || !strings.HasPrefix(out, "step\tquestion\task\twhere\nS14\tF-AGENTS.md\t") {
			t.Fatalf("the next run: exit %d, stderr %q, stdout\n%s\nwant 3, %s done and the stop of the prose step", code, errOut, out, id)
		}
	}
	target := filepath.Join(w, "target")
	s04 := gitIn(t, target, "rev-list", "--grep=^chore: setup S04$", "layup-setup")
	if got := gitIn(t, target, "diff", "--name-only", "main", s04); got != "docs/adr/0009-pin-the-baseline.md\ndocs/adr/README.md\n"+
		"docs/facts/F-0001-setup-answers.md\ndocs/facts/README.md\ndocs/setup/armature.pin\ndocs/setup/facts.sha256" {
		t.Errorf("the files of S04: %q", got)
	}
	lint := exec.Command("sh", "docs/adr/adr-lint.sh")
	lint.Dir = target
	if out, err := lint.Output(); err != nil || string(out) != "adr-lint: OK\n" {
		t.Errorf("sh docs/adr/adr-lint.sh on the target: %v, %q; want exit 0", err, out)
	}
	tab, err := verify.Check(w, []string{"pin", "facts"}, func(int, int, string) func() { return func() {} }, io.Discard)
	if err != nil || len(tab.Rows) != 2 || tab.Rows[0].Result != "pass" || tab.Rows[1].Result != "pass" {
		t.Errorf("checks pin and facts after S04: %v, %+v; want pass", err, tab.Rows)
	}
	cmds, err := os.ReadFile(filepath.Join(w, "out", "commands.sh"))
	if err != nil || !strings.Contains(string(cmds), "\ngit -C target remote add origin 'https://github.com/"+standin.Name+".git'\n") ||
		!strings.HasSuffix(string(cmds), "\ngit -C target push origin main\n") {
		t.Errorf("commands.sh: %v\n%s", err, cmds)
	}
	if out, err := exec.Command("sh", "-n", filepath.Join(w, "out", "commands.sh")).CombinedOutput(); err != nil {
		t.Errorf("sh -n commands.sh: %v\n%s", err, out)
	}
}

// A broken pin fails the evidence of S04 (D12 of #86): the run gives S04 fail
// with the reason of check pin, layup-setup goes back to the root commit, and
// a run with the record put right does S04 again, with one commit.
func TestABrokenPinFailsTheEvidenceOfS04(t *testing.T) {
	tmp := t.TempDir()
	url, _, err := standin.Baseline(filepath.Join(tmp, "baseline"))
	if err != nil {
		t.Fatal(err)
	}
	w := filepath.Join(tmp, "work")
	if err := os.MkdirAll(filepath.Join(w, "inputs", "briefs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(w, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	var b bytes.Buffer
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	if err := tsv.Write(&b, work.AnswersSchema, [][]string{{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", url, "operator", at, ""}}); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(w, "inputs", "answers.tsv"), b.String())
	stubSteps(t, map[string]setup.Outcome{"S04": {Kind: setup.Fail, Evidence: "a stop before S04"}})
	saved := setupSteps
	setupSteps = func(br setup.Brief, c setup.Calls) map[string]setup.Step { // S01 to S03 of this version, then a stop
		m, real := saved(br, c), setup.Steps(br, c)
		for _, id := range []string{"S01", "S02", "S03"} {
			m[id] = real[id]
		}
		return m
	}
	if code, out, _ := run("setup", w); code != 1 || !strings.Contains(out, "\nS04\tlayup-setup\tfail\ta stop before S04\n") {
		t.Fatalf("the run to S03: exit %d\n%s", code, out)
	}
	setupSteps = setup.Steps
	record := filepath.Join(w, "out", "record.tsv")
	good, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	tree := gitIn(t, filepath.Join(w, "target"), "rev-parse", "main^{tree}")
	writeFile(t, record, strings.Replace(string(good), "\tpin.tree\t"+tree+"\t", "\tpin.tree\t"+strings.Repeat("0", 40)+"\t", 1))
	code, out, _ := run("setup", w)
	if code != 1 || !strings.Contains(out, "\nS04\tlayup-setup\tfail\tpin: fail: tree: pin names "+strings.Repeat("0", 40)+", root commit has "+tree+"\n") {
		t.Errorf("a broken pin: exit %d\n%s\nwant S04 fail with the reason of check pin", code, out)
	}
	target := filepath.Join(w, "target")
	if gitIn(t, target, "rev-parse", "layup-setup") != gitIn(t, target, "rev-parse", "main") {
		t.Errorf("layup-setup is not at the root commit after the undo")
	}
	writeFile(t, record, string(good))
	code, out, _ = run("setup", w) // S04 again, then S05 and S06, and the prose step stops (task T-b3r1)
	if done, _ := os.ReadFile(record); code != 3 || !strings.Contains(string(done), "\nS04\tdone\tchecks pin and facts\t") ||
		strings.SplitN(gitIn(t, target, "log", "--reverse", "--format=%s", "main..layup-setup"), "\n", 2)[0] != "chore: setup S04" ||
		gitIn(t, target, "rev-list", "--count", "--grep=^chore: setup S04$", "layup-setup") != "1" {
		t.Errorf("the run with the record put right: exit %d\n%s\nwant S04 done with one commit", code, out)
	}
}

// scaffoldBaseline makes a stand-in baseline at dir for S05 to S14: the
// baseline's own link-lint.sh (LAYUP keeps it as its root commit has it), a
// note of the kit in a blockquote of the backlog, a guide whose link the
// deletion of the history breaks, two markers (one on two lines), and a file
// that check adapted flags. It gives the file:// URL.
func scaffoldBaseline(t *testing.T, dir string) string {
	t.Helper()
	url, _, err := standin.Baseline(dir)
	if err != nil {
		t.Fatal(err)
	}
	lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "links", "link-lint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	for p, text := range map[string]string{
		"docs/links/link-lint.sh": string(lint),
		"docs/tasks/backlog.md":   "# Backlog\n\n- **T-0001**: a task of the baseline ([#1](" + url + "/issues/1))\n\n> a note of the kit\n> with a link ([#2](" + url + "/issues/2))\n\nKeep this line.\n",
		"docs/guide.md":           "# Guide\n\nSee [the decision](decisions/D-0001-stand-in.md).\n",
		"docs/ops.md":             "# Ops\n\nThe port is \u2039port\u203a.\nAgain \u2039port\u203a.\nThe owner is \u2039owner\u203a.\n",
		"docs/how-to.md":          "# How to\n\nAdapt this to your project.\n",
	} {
		writeFile(t, filepath.Join(dir, filepath.FromSlash(p)), text)
	}
	if err := git.Add(dir); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(dir, "the scaffold of the baseline", git.Identity{Name: "t", Email: "t@layup.invalid", Time: time.Unix(0, 0)}); err != nil {
		t.Fatal(err)
	}
	return url
}

// writeInputs writes input files of the work area at w, under inputs/files/.
func writeInputs(t *testing.T, w string, files map[string]string) {
	t.Helper()
	for p, text := range files {
		path := filepath.Join(w, "inputs", "files", filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, path, text)
	}
}

// The demo of #90: after S04, S05 stops for the file whose links the deletion
// of the history breaks; the prose step stops once, with one row for each of
// its inputs, the flagged file too; S10 stops with one row per marker; with
// the inputs and the answers, each step is done with its checks as the
// evidence, and S12 is not built yet. A README.md that names no records
// branch fails the evidence of S14, and the run after the fix gives the same
// rows (condition 1 of the plan review).
func TestSetupRunsS05ToS14(t *testing.T) {
	tmp := t.TempDir()
	url := scaffoldBaseline(t, filepath.Join(tmp, "baseline"))
	w := filepath.Join(tmp, "work")
	if err := os.MkdirAll(filepath.Join(w, "inputs", "briefs"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(w, "inputs", "briefs", "problem-statement.md"), standin.Brief)
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	answers := [][]string{{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", url, "operator", at, ""}}
	writeAnswers := func() {
		var b bytes.Buffer
		if err := tsv.Write(&b, work.AnswersSchema, answers); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(w, "inputs", "answers.tsv"), b.String())
	}
	writeAnswers()
	if code, out, _ := run("setup", w); code != 3 || out != "step\tquestion\task\twhere\n"+
		"S05\tF-docs/guide.md\tFix the links of docs/guide.md that the deletion of the baseline's history breaks, and give the file as inputs/files/docs/guide.md.\tdocs/guide.md\n" {
		t.Fatalf("the stop of S05: exit %d\n%s", code, out)
	}
	writeInputs(t, w, map[string]string{"docs/guide.md": "# Guide\n\nThe decisions of the target are in its records.\n"})
	code, out, errOut := run("setup", w)
	want := "step\tquestion\task\twhere\n"
	for _, r := range [][]string{{"S14", "AGENTS.md", "Write the text of"}, {"S14", "README.md", "Write the text of"}, {"S08", "docs/glossary.md", "Write the text of"},
		{"S09", "docs/guardrails.md", "Write the text of"}, {"S14", "docs/how-to.md", "Adapt"}, {"S07", "docs/onboarding-for-engineers.md", "Write the text of"}} { // by path
		ask := r[2] + " " + r[1] + " for the target, and give it as inputs/files/" + r[1] + "."
		if r[2] == "Adapt" {
			ask = "Adapt " + r[1] + ", which check adapted flags, to the target, and give it as inputs/files/" + r[1] + "."
		}
		want += r[0] + "\tF-" + r[1] + "\t" + ask + "\t" + r[1] + "\n"
	}
	if code != 3 || out != want {
		t.Fatalf("the stop of the prose step: exit %d, stderr %q\n%s\nwant 3 and\n%s", code, errOut, out, want)
	}
	inputs := map[string]string{
		"docs/onboarding-for-engineers.md": "# Onboarding\n\nThe work starts from [the problem statement](facts/problem-statement-brief.md).\n",
		"docs/glossary.md":                 "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository |\n",
		"docs/guardrails.md":               "# Guardrails\n\n## 1. Rules\n\n- **Inv-1** — the setup keeps its record. Check: no check yet\n",
		"README.md":                        "# " + standin.Name + "\n\nThe product repository, set up from its pinned baseline ([the pin](docs/setup/armature.pin)).\n",
		"AGENTS.md":                        "# AGENTS.md\n\nAgent context for **" + standin.Name + "**.\n",
		"docs/how-to.md":                   "# How to\n\nRun the tests of the product.\n",
	}
	writeInputs(t, w, inputs)
	if code, out, _ := run("setup", w); code != 1 || !strings.Contains(out, "\nS14\tlayup-setup\tfail\tidentity: fail: branch: README.md does not name the branch layup-records\n") {
		t.Fatalf("a README.md that names no records branch: exit %d\n%s", code, out)
	}
	target := filepath.Join(w, "target")
	head := gitIn(t, target, "rev-parse", "layup-setup")
	inputs["README.md"] = strings.TrimSuffix(inputs["README.md"], "\n") + " Its records are on the branch `layup-records`.\n"
	writeInputs(t, w, inputs)
	code, out, _ = run("setup", w)
	ask := func(file, m string) string {
		return "What is the value of " + m + " in " + file + "? Answer gap to keep it as an open gap, with its question as question_text."
	}
	port, owner := "\u2039port\u203a", "\u2039owner\u203a"
	idPort, idOwner := setup.MarkerID("docs/ops.md", port), setup.MarkerID("docs/ops.md", owner)
	if code != 3 || out != "step\tquestion\task\twhere\n"+"S10\t"+idPort+"\t"+ask("docs/ops.md", port)+"\tdocs/ops.md:3 "+port+"\n"+
		"S10\t"+idOwner+"\t"+ask("docs/ops.md", owner)+"\tdocs/ops.md:5 "+owner+"\n" {
		t.Fatalf("the stop of S10: exit %d\n%s", code, out)
	}
	answers = append(answers, []string{idPort, "8080", "operator", at, ""}, []string{idOwner, "gap", "operator", at, "Who owns the operations?"})
	writeAnswers()
	if code, out, errOut := run("setup", w); code != 1 || !strings.Contains(out, "\nS11\tlayup-setup\tdone\tchecks markers, sources and facts\n") ||
		!strings.Contains(out, "\nS12\tlayup-setup\tnot-active\tnot built yet\n") {
		t.Fatalf("the run with the answers: exit %d, stderr %q\n%s", code, errOut, out)
	}
	if got := gitIn(t, target, "log", "--format=%s", head+"..layup-setup"); got != "chore: setup S11\nchore: setup S14" {
		t.Errorf("the commits after the undo of S14: %q; want S14 once, then S11", got)
	}
	if got := gitIn(t, target, "log", "--reverse", "--format=%s", "main..layup-setup"); got != "chore: setup S04\nchore: setup S05\nchore: setup S06\n"+
		"chore: setup S07\nchore: setup S08\nchore: setup S09\nchore: setup S14\nchore: setup S11" {
		t.Errorf("the setup commits: %q", got)
	}
	for p, want := range map[string]string{
		"docs/ops.md":              "# Ops\n\nThe port is 8080.\nAgain 8080.\nThe owner is " + owner + ".\n",
		"docs/setup/open-gaps.tsv": "docs/ops.md\t" + owner + "\tWho owns the operations?\n",
		"docs/guide.md":            "# Guide\n\nThe decisions of the target are in its records.\n",
		"docs/tasks/backlog.md":    "# Backlog\n\n\n\nKeep this line.\n",
	} {
		if got := gitIn(t, target, "show", "layup-setup:"+p) + "\n"; got != want {
			t.Errorf("%s on layup-setup:\n%q\nwant\n%q", p, got, want)
		}
	}
	for _, p := range []string{"docs/decisions", "docs/audit", "docs/tasks/T-0001.md"} {
		if out, err := exec.Command("git", "-C", target, "cat-file", "-e", "layup-setup:"+p).CombinedOutput(); err == nil {
			t.Errorf("%s is on layup-setup: %s", p, out)
		}
	}
	record, _ := work.ReadRecord(w)
	for _, r := range [][]string{
		{"S05", "file:docs/guide.md"}, {"S06", "brief.copy"}, {"S06", "brief.copy.sha256"}, {"S07", "file:docs/onboarding-for-engineers.md"},
		{"S14", "file:README.md"}, {"S14", "file:docs/how-to.md"}, {"S11", "marker:docs/ops.md:3"}, {"S11", "marker:docs/ops.md:4"},
		{"S11", "marker:docs/ops.md:5"}, {"S11", "marker.record"}, {"S10", "answers.sha256"},
	} {
		if v, ok := record.Value(r[0], r[1]); !ok || v == "" {
			t.Errorf("the record has no row %s %s", r[0], r[1])
		}
	}
	tab, err := verify.Check(w, []string{"kit-history", "link-lint", "facts", "onboarding", "glossary", "guardrails", "markers", "adapted", "identity", "sources"},
		func(int, int, string) func() { return func() {} }, io.Discard)
	for _, r := range tab.Rows {
		if err != nil || r.Result != "pass" {
			t.Errorf("check %s after S11: %s %q, %v; want pass", r.Check, r.Result, r.Reason, err)
		}
	}
}
