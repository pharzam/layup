// Package session of the fixture module starts the command of a register row,
// which is not a string literal: its register row allows it.
package session

import "os/exec"

// Start starts the words of a command.
func Start(command []string) *exec.Cmd { return exec.Command(command[0], command[1:]...) }
