package gate

import (
	"fmt"
)

// The results of a kind (the table gate-result).
const (
	pass      = "pass"
	fail      = "fail"
	notActive = "not-active"
	clear     = "clear"
)

// An outcome is how the command of a kind ended: its exit code, or the name of
// the signal that killed it.
type outcome struct {
	code   int
	signal string
}

// decide gives the result and the reason of an active kind by the table of the
// run (docs/spec/gate.md, The run): the first line that matches decides. files
// are the files of the scratch tree; run runs the command and gives how it
// ended. The reason of pass is "" (— in the table).
func decideActive(k Kind, files []string, lookPath func(string) (string, error), run func(string) outcome) (string, string) {
	if _, err := lookPath(k.Tool); err != nil {
		return notActive, "tool not found: " + k.Tool
	}
	if _, ok := k.inScope(files); !ok {
		return clear, "no product path"
	}
	o := run(k.Command)
	switch {
	case o.signal != "":
		return fail, "signal " + o.signal
	case o.code != 0:
		return fail, fmt.Sprintf("exit %d", o.code)
	}
	return pass, ""
}

// decidePending gives the result and the reason of a pending kind from the
// paths that the change changes; its command never runs.
func decidePending(k Kind, changed []string) (string, string) {
	if p, ok := k.inScope(changed); ok {
		return fail, "pending: product path changed: " + p
	}
	return clear, "pending: no product path"
}
