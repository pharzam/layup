// Package cli is the command-line surface of layup. It holds no state: each
// command reads and writes plain files in a repository (ADR-0011).
package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/pharzam/layup/internal/psb"
)

// Version is the version that `layup version` prints.
const Version = "0.1.0-dev"

const usage = `usage: layup <command>

commands:
  version           print the version of layup
  psb check FILE    print the gap questions of a problem statement as TSV;
                    exit 0 with no gap, 1 with gaps, 2 on a usage error
`

// Run executes one layup command and returns its exit code: 0 for success,
// 2 for a usage error.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "layup %s\n", Version)
		return 0
	case "psb":
		return psbCommand(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "layup: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

func psbCommand(args []string, stdout, stderr io.Writer) int {
	if len(args) != 2 || args[0] != "check" {
		fmt.Fprintf(stderr, "layup: usage: layup psb check FILE\n")
		return 2
	}
	src, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Fprintf(stderr, "layup: %v\n", err)
		return 2
	}
	gaps := psb.Check(string(src))
	psb.WriteTSV(stdout, gaps)
	if len(gaps) > 0 {
		return 1
	}
	return 0
}
