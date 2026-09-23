// Command layup is the LAYUP core engine (ADR-0011).
package main

import (
	"os"

	"github.com/pharzam/layup/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
