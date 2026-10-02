// Package cli is the command-line surface of layup. It holds no state: each
// command reads and writes plain files in a repository (ADR-0011). It is the
// frame of every command (docs/spec/README.md, Commands): one command table,
// the argument rules, the usage, the exit codes and the progress lines.
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/pharzam/layup/internal/psb"
)

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

// psbCheck keeps the rule of its own section: 0 no gap, 1 a gap.
func psbCheck(in call) int {
	src, err := os.ReadFile(in.args[0])
	if err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	gaps := psb.Check(string(src))
	psb.WriteTSV(in.stdout, gaps)
	if len(gaps) > 0 {
		return exitFail
	}
	return exitPass
}
