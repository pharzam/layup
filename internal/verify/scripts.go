package verify

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"syscall"
)

// The baseline's own check scripts that the checks discipline-tests and
// link-lint run (D5 of #87).
const (
	discTests = "docs/tests/run-discipline-tests.sh"
	linkLint  = "docs/links/link-lint.sh"
)

// A scriptOutcome is how a script ended: its exit code, or the name of the
// signal that killed it (the form of internal/gate).
type scriptOutcome struct {
	code   int
	signal string
}

// notActive is the result of a check that could not run.
const notActive = "not-active"

// lookPath finds sh; the unit tests replace it.
var lookPath = exec.LookPath

// decideScript gives the row of a baseline script by the rule of a gate
// command (gate.md, D5 of #87): a script that is not a file of the tree is
// not-active, as is an sh that is not found; exit 0 is pass; another exit or
// a signal is fail. No reason holds a path of the scratch tree.
func decideScript(name, script string, fsys fs.FS, lookPath func(string) (string, error), run func(string) scriptOutcome) Row {
	if !isRegular(fsys, script) {
		return Row{name, notActive, "missing: " + script}
	}
	if _, err := lookPath("sh"); err != nil {
		return Row{name, notActive, "tool not found: sh"}
	}
	switch o := run(script); {
	case o.signal != "":
		return Row{name, "fail", "signal " + o.signal}
	case o.code != 0:
		return Row{name, "fail", fmt.Sprintf("exit %d", o.code)}
	}
	return Row{name, "pass", ""}
}

// runScript runs sh script in dir, with the environment of layup, and gives
// how it ended; its standard output goes to stdout and its standard error to
// stderr.
func runScript(dir, script string, stdout, stderr io.Writer) scriptOutcome {
	cmd := exec.Command("sh", script)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, stdout, stderr
	err := cmd.Run()
	var exit *exec.ExitError
	switch {
	case err == nil:
		return scriptOutcome{}
	case errors.As(err, &exit):
		if ws, ok := exit.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return scriptOutcome{signal: ws.Signal().String()}
		}
		return scriptOutcome{code: exit.ExitCode()}
	}
	fmt.Fprintln(stderr, err)
	return scriptOutcome{code: -1}
}

// scriptRow runs the baseline script in the checked-out tree at dir: the row
// of its check, and its output on out (standard error, never the table).
func scriptRow(dir, name, script string, out io.Writer) Row {
	return decideScript(name, script, os.DirFS(dir), lookPath, func(s string) scriptOutcome { return runScript(dir, s, out, out) })
}

// BrokenLinks gives the files whose links link-lint.sh of the checked-out
// tree at dir refuses, once each, in byte order: the call of S05, which gives
// WORK/target after its deletion and before its commit (K10, D6 of #87 with
// condition 1 of its plan review).
func BrokenLinks(dir string) ([]string, error) {
	if !isRegular(os.DirFS(dir), linkLint) {
		return nil, errors.New(linkLint + " is not a file of the tree")
	}
	var stderr bytes.Buffer
	o := runScript(dir, linkLint, io.Discard, &stderr)
	if o.signal != "" {
		return nil, errors.New(linkLint + ": signal " + o.signal)
	}
	return brokenLinks(o.code, stderr.Bytes())
}

var brokenLink = regexp.MustCompile(`^FAIL  L[0-9]+: (.+):[0-9]+ (?:links|uses) `)

// brokenLinks reads the standard error of link-lint.sh: each line FAIL
// L<n>: <path>:<line> names a file; a line FAIL that names none, or a failed
// exit with no file, is an error.
func brokenLinks(code int, stderr []byte) ([]string, error) {
	var files []string
	for _, line := range strings.Split(string(stderr), "\n") {
		if !strings.HasPrefix(line, "FAIL") {
			continue
		}
		m := brokenLink.FindStringSubmatch(line)
		if m == nil {
			return nil, errors.New("link-lint.sh: " + line)
		}
		files = append(files, m[1])
	}
	if code != 0 && len(files) == 0 {
		return nil, fmt.Errorf("link-lint.sh: exit %d with no file", code)
	}
	slices.Sort(files)
	return slices.Compact(files), nil
}
