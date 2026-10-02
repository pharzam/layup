// Package cli of the fixture module reads its input from an environment
// variable: a breach of the input rule of docs/spec/README.md.
package cli

import "os"

// Run gives 0 when the input is set.
func Run(args []string) int {
	if os.Getenv("LAYUP_WORK") == "" {
		return 2
	}
	return 0
}
