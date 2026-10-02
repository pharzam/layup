package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/setup"
)

// standInSetup puts a stand-in for setup.Run in place until the test ends: it
// calls step for each row of res, then gives res and err.
func standInSetup(t *testing.T, res setup.Result, err error) *string {
	t.Helper()
	var work string
	saved := setupRun
	setupRun = func(dir string, _ map[string]setup.Step, _ git.Identity, step func(i, n int, name string) func()) (setup.Result, error) {
		work = dir
		for i, r := range res.Steps {
			step(i+1, len(res.Steps), r.Step)()
		}
		return res, err
	}
	t.Cleanup(func() { setupRun = saved })
	return &work
}

func TestSetupCommand(t *testing.T) {
	standInSetup(t, setup.Result{}, nil)
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{[]string{"setup"}, "missing argument WORK"},
		{[]string{"setup", "w", "x"}, `extra argument "x"`},
		{[]string{"setup", "w", "--x", "y"}, `unknown flag "--x"`},
	} {
		if code, out, errOut := run(c.args...); code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: "+c.reason+"\n\nusage: layup ") {
			t.Errorf("%q: exit %d, stdout %q, stderr %q; want 2, nothing and %q with the usage", c.args, code, out, errOut, c.reason)
		}
	}
	steps := func(results ...string) setup.Result {
		var r setup.Result
		for i, res := range results {
			r.Steps = append(r.Steps, setup.StepRow{Step: []string{"S01", "S02"}[i], Actor: "layup-setup", Result: res, Evidence: "x"})
		}
		return r
	}
	for _, c := range []struct {
		name string
		res  setup.Result
		code int
		head string
	}{
		{"done and a hand-off", steps("done", "operator"), 0, "step\tactor\tresult\tevidence\n"},
		{"a fail", steps("done", "fail"), 1, "step\tactor\tresult\tevidence\n"},
		{"a stop", setup.Result{Stops: []setup.StopRow{{Step: "S01", Question: "S01-stack", Ask: "Which stack?"}}}, 3, "step\tquestion\task\twhere\n"},
	} {
		work := standInSetup(t, c.res, nil)
		code, out, errOut := run("setup", "w")
		if code != c.code || !strings.HasPrefix(out, c.head) || *work != "w" {
			t.Errorf("%s: exit %d, stdout %q, WORK %q; want %d and the table of w", c.name, code, out, *work, c.code)
		}
		if c.res.Steps != nil && !strings.HasPrefix(errOut, "layup setup: [1/2] S01\n") {
			t.Errorf("%s: stderr %q; want the progress line of each step", c.name, errOut)
		}
	}
	standInSetup(t, setup.Result{}, &setup.InputError{Err: errors.New("the work area w is not a directory")})
	if code, out, errOut := run("setup", "w"); code != 2 || out != "" || errOut != "layup: the work area w is not a directory\n" {
		t.Errorf("an input error: exit %d, stdout %q, stderr %q; want 2, nothing and the reason", code, out, errOut)
	}
	if _, _, errOut := run(); !strings.Contains(errOut, "\n  setup WORK ") || !strings.Contains(errOut, "\n  setup verify WORK ") {
		t.Errorf("the usage has no setup WORK beside setup verify WORK:\n%s", errOut)
	}
}
