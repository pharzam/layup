// Package cli is the command-line surface of layup. It holds no state: each
// command reads and writes plain files in a repository (ADR-0011).
package cli

import (
	"fmt"
	"io"
)

// Version is the version that `layup version` prints.
const Version = "0.1.0-dev"

const usage = `usage: layup <command>

commands:
  version   print the version of layup
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
	default:
		fmt.Fprintf(stderr, "layup: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}
