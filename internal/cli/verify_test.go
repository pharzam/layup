package cli

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/verify"
)

// standInVerify puts a stand-in for verify.Run in place until the test ends:
// it calls step for each row, then gives the rows and err.
func standInVerify(t *testing.T, rows []verify.Row, err error) *string {
	t.Helper()
	var work string
	saved := verifyRun
	verifyRun = func(dir string, step func(i, n int, check string) func(), out io.Writer) (verify.Table, error) {
		work = dir
		for i, r := range rows {
			step(i+1, len(rows), r.Check)()
		}
		if errors.As(err, new(*verify.InputError)) {
			return verify.Table{}, err
		}
		return verify.Table{Rows: rows}, err
	}
	t.Cleanup(func() { verifyRun = saved })
	return &work
}

func TestSetupVerify(t *testing.T) {
	standInVerify(t, nil, nil)
	for _, c := range []struct {
		args   []string
		reason string
	}{
		{[]string{"setup", "verify"}, "missing argument WORK"},
		{[]string{"setup", "verify", "w", "x"}, `extra argument "x"`},
		{[]string{"setup", "verify", "w", "--x", "y"}, `unknown flag "--x"`},
	} {
		if code, out, errOut := run(c.args...); code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: "+c.reason+"\n\nusage: layup ") {
			t.Errorf("%q: exit %d, stdout %q, stderr %q; want 2, nothing and %q with the usage", c.args, code, out, errOut, c.reason)
		}
	}
	for _, c := range []struct {
		results []string
		want    int
	}{{[]string{"pass", "clear"}, 0}, {[]string{"pass", "not-active"}, 1}, {[]string{"fail", "pass"}, 1}} {
		rows := []verify.Row{{Check: "pin", Result: c.results[0]}, {Check: "kit-history", Result: c.results[1], Reason: "x"}}
		work := standInVerify(t, rows, nil)
		code, out, errOut := run("setup", "verify", "w")
		if code != c.want || !strings.HasPrefix(out, "check\tresult\treason\npin\t") || strings.Count(out, "\n") != 3 || *work != "w" {
			t.Errorf("%q: exit %d, stdout %q, WORK %q; want %d and the table of w", c.results, code, out, *work, c.want)
		}
		if want := "layup setup verify: [1/2] pin\nlayup setup verify: [2/2] kit-history\n"; errOut != want {
			t.Errorf("stderr %q; want the step line of each row %q", errOut, want)
		}
	}
	standInVerify(t, nil, &verify.InputError{Err: errors.New("target: no commit at the branch layup-setup")})
	if code, out, errOut := run("setup", "verify", "w"); code != 2 || out != "" || errOut != "layup: target: no commit at the branch layup-setup\n" {
		t.Errorf("an input error: exit %d, stdout %q, stderr %q; want 2, nothing and the reason", code, out, errOut)
	}
	standInVerify(t, []verify.Row{{Check: "pin", Result: "pass"}}, &verify.CleanupError{Path: "/tmp/layup-verify-1", Err: errors.New("busy")})
	if code, out, errOut := run("setup", "verify", "w"); code != 2 || strings.Count(out, "\n") != 2 || !strings.Contains(errOut, "layup: the scratch tree /tmp/layup-verify-1 is not removed") {
		t.Errorf("a leftover scratch tree: exit %d, stdout %q, stderr %q; want 2, the table and the path", code, out, errOut)
	}
	if _, _, errOut := run(); !strings.Contains(errOut, "\n  setup verify WORK ") {
		t.Errorf("the usage has no setup verify:\n%s", errOut)
	}
}
