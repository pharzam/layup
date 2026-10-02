//go:build integration

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/work"
)

// stubSteps gives the steps of a setup: each done, unless out names another
// outcome (condition 4 of the plan review of #85).
func stubSteps(t *testing.T, out map[string]setup.Outcome) {
	t.Helper()
	saved := setupSteps
	setupSteps = func() map[string]setup.Step {
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
