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
	setupSteps = func(setup.Brief, setup.Checks) map[string]setup.Step {
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
	code, out, errOut := run("setup", w)
	if code != 1 || !strings.HasPrefix(out, "step\tactor\tresult\tevidence\n"+
		"S01\tlayup-setup\tdone\tevery answer present; the stack has a catalog entry\n"+
		"S02\tlayup-setup\tdone\tthe commit and the tree\n"+
		"S03\tlayup-setup\tdone\troot tree = pin.tree\n"+
		"S04\tlayup-setup\tdone\tchecks pin and facts\n"+
		"S05\tlayup-setup\tnot-active\tnot built yet\n") {
		t.Fatalf("the next run: exit %d, stderr %q, stdout\n%s\nwant 1, S01 to S04 done and S05 not built yet", code, errOut, out)
	}
	target := filepath.Join(w, "target")
	if got := gitIn(t, target, "diff", "--name-only", "main", "layup-setup"); got != "docs/adr/0009-pin-the-baseline.md\ndocs/adr/README.md\n"+
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
	setupSteps = func(br setup.Brief, c setup.Checks) map[string]setup.Step { // S01 to S03 of this version, then a stop
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
	if code, out, _ := run("setup", w); code != 1 || !strings.Contains(out, "\nS04\tlayup-setup\tdone\tchecks pin and facts\n") ||
		gitIn(t, target, "log", "--format=%s", "main..layup-setup") != "chore: setup S04" {
		t.Errorf("the run with the record put right: exit %d\n%s\nwant S04 done with one commit", code, out)
	}
}
