package cli

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/gate"
)

// standIn puts a stand-in for gate.Run in place until the test ends: it calls
// step for each row and writes a block line, then gives the rows and err.
func standIn(t *testing.T, rows []gate.Row, err error) *[]string {
	t.Helper()
	var args []string
	saved := gateRun
	gateRun = func(repo, base, head string, step func(i, n int, kind string) func(), out io.Writer) (gate.Table, error) {
		args = []string{repo, base, head}
		for i, r := range rows {
			done := step(i+1, len(rows), r.Kind)
			done()
			io.WriteString(out, "the output of "+r.Kind+"\n")
		}
		if errors.As(err, new(*gate.InputError)) {
			return gate.Table{}, err
		}
		return gate.Table{Base: strings.Repeat("1", 40), Head: strings.Repeat("2", 40), Rows: rows}, err
	}
	t.Cleanup(func() { gateRun = saved })
	return &args
}

func TestGateUsageErrors(t *testing.T) {
	standIn(t, nil, nil)
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{[]string{"gate"}, "missing flag --base"},
		{[]string{"gate", "r", "--base", "b"}, "missing flag --head"},
		{[]string{"gate", "--base", "b", "--head", "h"}, "missing argument REPO"},
		{[]string{"gate", "r", "--base", "b", "--head", "h", "x"}, `extra argument "x"`},
		{[]string{"gate", "r", "--base", "b", "--head", "h", "--x", "y"}, `unknown flag "--x"`},
	} {
		code, out, errOut := run(c.args...)
		if code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: "+c.reason+"\n\nusage: layup ") {
			t.Errorf("%q: exit %d, stdout %q, stderr %q; want 2, nothing and %q with the usage", c.args, code, out, errOut, c.reason)
		}
	}
}

func TestGatePrintsTheTableAndGivesItsExitCode(t *testing.T) {
	for _, c := range []struct {
		results []string
		want    int
	}{
		{[]string{"pass", "clear"}, 0},
		{[]string{"pass", "fail"}, 1},
		{[]string{"clear", "not-active"}, 1},
	} {
		var rows []gate.Row
		for i, r := range c.results {
			rows = append(rows, gate.Row{Kind: []string{"static", "layout"}[i], State: "active", Result: r, Reason: "x"})
		}
		args := standIn(t, rows, nil)
		code, out, errOut := run("gate", "repo", "--head", "HEAD", "--base=main")
		if code != c.want || !strings.HasPrefix(out, "base\thead\tkind\tstate\tresult\treason\n") || strings.Count(out, "\n") != 3 {
			t.Errorf("%q: exit %d, stdout %q; want %d and the table", c.results, code, out, c.want)
		}
		if strings.Join(*args, " ") != "repo main HEAD" {
			t.Errorf("gate.Run got %q; want repo main HEAD", *args)
		}
		if want := "layup gate: [1/2] static\nthe output of static\nlayup gate: [2/2] layout\nthe output of layout\n"; errOut != want {
			t.Errorf("stderr %q; want the step line and the block of each kind:\n%q", errOut, want)
		}
	}
}

func TestGateGivesTwoOnAnInputErrorOrALeftoverScratchTree(t *testing.T) {
	standIn(t, nil, &gate.InputError{Err: errors.New("the revision \"x\" is not a commit")})
	if code, out, errOut := run("gate", "repo", "--base", "x", "--head", "HEAD"); code != 2 || out != "" || errOut != "layup: the revision \"x\" is not a commit\n" {
		t.Errorf("an input error: exit %d, stdout %q, stderr %q; want 2, nothing and the reason", code, out, errOut)
	}
	standIn(t, []gate.Row{{Kind: "static", State: "active", Result: "pass"}}, &gate.CleanupError{Path: "/tmp/layup-gate-1", Err: errors.New("busy")})
	if code, out, errOut := run("gate", "repo", "--base", "main", "--head", "HEAD"); code != 2 || strings.Count(out, "\n") != 2 || !strings.Contains(errOut, "layup: the scratch tree /tmp/layup-gate-1 is not removed") {
		t.Errorf("a leftover scratch tree: exit %d, stdout %q, stderr %q; want 2, the table and the path", code, out, errOut)
	}
}

func TestTheUsageListsGate(t *testing.T) {
	if _, _, errOut := run(); !strings.Contains(errOut, "\n  gate REPO --base REV --head REV ") {
		t.Errorf("the usage has no gate:\n%s", errOut)
	}
}
