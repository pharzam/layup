package cli

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/verify"
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

// standInBrief puts stand-ins in place until the test ends: WORK is a
// directory, readFile gives the problem statement src (nil: absent), and
// setupSteps keeps the brief that the command hands to the steps.
func standInBrief(t *testing.T, src []byte) *setup.Brief {
	t.Helper()
	got := &setup.Brief{Sum: "not called"}
	savedDir, savedRead, savedSteps := isDir, readFile, setupSteps
	isDir = func(string) bool { return true }
	readFile = func(name string) ([]byte, error) {
		if name != filepath.Join("w", "inputs", "briefs", "problem-statement.md") && name != "brief.md" {
			t.Errorf("readFile(%q); want the problem statement of w", name)
		}
		if src == nil {
			return nil, fs.ErrNotExist
		}
		return src, nil
	}
	setupSteps = func(b setup.Brief, _ setup.Checks) map[string]setup.Step { *got = b; return setup.Stubs() }
	t.Cleanup(func() { isDir, readFile, setupSteps = savedDir, savedRead, savedSteps })
	return got
}

// The command reads the problem statement of WORK once and hands its gap
// table, as layup psb check writes it, and its SHA-256 to the steps (D5 of
// #86); a problem statement that is absent or not valid UTF-8 is exit 2.
func TestSetupReadsTheProblemStatement(t *testing.T) {
	work := standInSetup(t, setup.Result{Steps: []setup.StepRow{{Step: "S01", Actor: "layup-setup", Result: "done", Evidence: "x"}}}, nil)
	src := []byte("# A brief\n\nThe product must be fast.\n")
	got := standInBrief(t, src)
	if code, _, errOut := run("setup", "w"); code != 0 || *work != "w" {
		t.Fatalf("a problem statement: exit %d, stderr %q; want 0 and a run", code, errOut)
	}
	_, table, _ := run("psb", "check", "brief.md")
	if string(got.Gaps) != table || !strings.Contains(table, "\nQ-") || got.Sum != fmt.Sprintf("%x", sha256.Sum256(src)) {
		t.Errorf("the brief: gaps %q, sum %q; want the table of layup psb check, %q, and the SHA-256 of the file", got.Gaps, got.Sum, table)
	}
	for _, c := range []struct {
		src    []byte
		reason string
	}{
		{nil, "inputs/briefs/problem-statement.md: the problem statement is absent"},
		{[]byte("# A brief\n\xff\n"), "inputs/briefs/problem-statement.md: line 2 is not valid UTF-8"},
	} {
		*work = ""
		got := standInBrief(t, c.src)
		if code, out, errOut := run("setup", "w"); code != 2 || out != "" || errOut != "layup: "+c.reason+"\n" || *work != "" || got.Sum != "not called" {
			t.Errorf("%q: exit %d, stdout %q, stderr %q, run %q; want 2, nothing, %q and no run", c.src, code, out, errOut, *work, c.reason)
		}
	}
}

// The evidence call gives the reason of the first row that is not pass or
// clear, or "" (D12 of #86).
func TestTheEvidenceCall(t *testing.T) {
	saved := verifyCheck
	t.Cleanup(func() { verifyCheck = saved })
	var dir string
	var names []string
	for _, c := range []struct {
		rows []verify.Row
		err  error
		want string
	}{
		{[]verify.Row{{Check: "pin", Result: "pass", Reason: ""}, {Check: "facts", Result: "clear", Reason: ""}}, nil, ""},
		{[]verify.Row{{Check: "pin", Result: "pass", Reason: ""}, {Check: "facts", Result: "fail", Reason: "hash: x does not match"}, {Check: "x", Result: "not-active", Reason: "y"}}, nil, "facts: fail: hash: x does not match"},
		{[]verify.Row{{Check: "pin", Result: "not-active", Reason: "no layup-setup"}}, nil, "pin: not-active: no layup-setup"},
		{nil, &verify.InputError{Err: errors.New("the work area w has no target\nmore")}, "the evidence call: the work area w has no target"},
	} {
		verifyCheck = func(d string, n []string, _ func(i, n int, check string) func(), _ io.Writer) (verify.Table, error) {
			dir, names = d, n
			return verify.Table{Rows: c.rows}, c.err
		}
		if got := evidence(io.Discard)("w", []string{"pin", "facts"}); got != c.want || dir != "w" || !slices.Equal(names, []string{"pin", "facts"}) {
			t.Errorf("%v %v: %q on %q %q; want %q on w, pin and facts", c.rows, c.err, got, dir, names, c.want)
		}
	}
}
