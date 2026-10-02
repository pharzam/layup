// Package cli is the command-line surface of layup. It holds no state: each
// command reads and writes plain files in a repository (ADR-0011). It is the
// frame of every command (docs/spec/README.md, Commands): one command table,
// the argument rules, the usage, the exit codes and the progress lines.
package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/psb"
	"github.com/pharzam/layup/internal/verify"
)

// gateRun runs layup gate; the unit tests replace it.
var gateRun = gate.Run

// readFile reads the FILE of layup psb check; the unit tests replace it.
var readFile = os.ReadFile

// verifyRun runs layup setup verify; the unit tests replace it.
var verifyRun = verify.Run

// Version is the version that `layup version` prints.
const Version = "0.1.0-dev"

// The exit codes of every command (docs/spec/README.md, Commands).
const (
	exitPass  = 0 // every row passed; for psb check, no gap
	exitFail  = 1 // a row failed or did not run; for psb check, a gap
	exitUsage = 2 // a usage or input error
	exitStop  = 3 // only layup setup: the run stopped for a human input
)

// exitCode maps the result column of a table to its exit code: exitPass only
// when the table has a row and each row is pass, clear, done or operator.
// Each other word, and a table with no row, give exitFail, so a check that
// did not run never gives 0 (NFR-004). A usage or input error and a stop of
// layup setup do not come from a table.
func exitCode(results []string) int {
	if len(results) == 0 {
		return exitFail
	}
	for _, r := range results {
		switch r {
		case "pass", "clear", "done", "operator":
		default:
			return exitFail
		}
	}
	return exitPass
}

// commands is the command table: the dispatch, the argument check and the
// usage all come from it, so a new command adds one row.
var commands = []command{
	{words: []string{"version"}, help: "print the version of layup", run: version},
	{words: []string{"psb", "check"}, args: []string{"FILE"},
		help: "print the gap questions of a problem statement as a table", run: psbCheck},
	{words: []string{"gate"}, args: []string{"REPO"}, flags: []flag{{"base", "REV"}, {"head", "REV"}},
		help: "run the gate kinds of the manifest at --base on --head", run: gateCommand},
	{words: []string{"setup", "verify"}, args: []string{"WORK"},
		help: "check the setup of the target of the work area WORK, from outside", run: setupVerify},
}

// Run executes one layup command and returns its exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	return dispatch(commands, args, stdout, stderr)
}

// dispatch executes one command of table. A usage error prints the reason
// and the usage on standard error, and nothing on standard output.
func dispatch(table []command, args []string, stdout, stderr io.Writer) int {
	c, in, err := parse(table, args)
	if err != nil {
		fmt.Fprintf(stderr, "layup: %v\n\n%s", err, usage(table))
		return exitUsage
	}
	in.stdout, in.stderr = stdout, stderr
	return c.run(in)
}

func version(in call) int {
	fmt.Fprintf(in.stdout, "layup %s\n", Version)
	return exitPass
}

// psbCheck keeps the rule of its own section: 0 no gap, 1 a gap; 2 for a FILE
// that it cannot read or that is not valid UTF-8 (K32), and for a table that it
// cannot write.
func psbCheck(in call) int {
	src, err := readFile(in.args[0])
	if err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	if n := invalidLine(src); n > 0 {
		fmt.Fprintf(in.stderr, "layup: %s: line %d is not valid UTF-8\n", in.args[0], n)
		return exitUsage
	}
	gaps := psb.Check(string(src))
	if err := psb.WriteTSV(in.stdout, gaps); err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	if len(gaps) > 0 {
		return exitFail
	}
	return exitPass
}

// invalidLine gives the line, from 1, of the first byte of src that is not
// valid UTF-8, or 0 when src is valid.
func invalidLine(src []byte) int {
	for i := 0; i < len(src); {
		r, size := utf8.DecodeRune(src[i:])
		if r == utf8.RuneError && size == 1 {
			return 1 + bytes.Count(src[:i], []byte{'\n'})
		}
		i += size
	}
	return 0
}

// gateCommand runs layup gate (docs/spec/gate.md): the table on standard
// output; the progress lines and the block of each kind on standard error.
// The beats of a kind stop before its block, and before the table.
func gateCommand(in call) int {
	p := newProgress(in.stderr, "gate")
	t, err := gateRun(in.args[0], in.flags["base"], in.flags["head"], func(i, n int, kind string) func() {
		p.step(i, n, kind)
		return p.end
	}, in.stderr)
	p.end()
	if errors.As(err, new(*gate.InputError)) {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	if werr := t.Write(in.stdout); werr != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", werr)
		return exitUsage
	}
	if err != nil { // the scratch tree is left in REPO; the table is complete
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	return exitCode(t.Results())
}

// setupVerify runs layup setup verify (docs/spec/setup.md): the table on
// standard output; the progress lines and the diagnostics on standard error.
func setupVerify(in call) int {
	p := newProgress(in.stderr, "setup verify")
	t, err := verifyRun(in.args[0], func(i, n int, check string) func() {
		p.step(i, n, check)
		return p.end
	}, in.stderr)
	p.end()
	if errors.As(err, new(*verify.InputError)) {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	if werr := t.Write(in.stdout); werr != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", werr)
		return exitUsage
	}
	if err != nil { // the scratch tree is left; the table is complete
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	return exitCode(t.Results())
}
