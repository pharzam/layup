package spec

// The pinned prd-lint of the target (docs/spec/loop.md, The rendered commit
// and the confirmed inventory; docs/spec/packages.md, the table of M2e; task
// T-4wyg, #351): the one process seam of the package. The target's own script
// is run as `sh <script> <prd dir>`; a specification that waits for its bet
// leaves the MoSCoW and Phase cells empty, so its head is accepted when every
// failure line of the script names an empty cell — and only then. The script
// is never modified, relaxed or replaced.

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// emptyCellFailure is the full-line form of the one failure the waiting
// specification accepts, exactly as the pinned producer writes it (a `FAIL  `
// line of the awk checks, stdout): a well-formed PRD file name, a requirement
// ID, and the empty MoSCoW or the absent phase.
var emptyCellFailure = regexp.MustCompile(`^FAIL  PRD-[0-9]{4}-[a-z0-9][a-z0-9-]*\.md: (REQ|NFR)-[0-9]{3} (MoSCoW "" not in Must\|Should\|Could\|Won't|has no phase value)$`)

// RunPRDLint runs the pinned script as `sh <script> <prdDir>` and gives the
// exit code and every non-empty line of each stream, separately: the awk
// checks print to stdout and the script's err() to stderr, and the split is
// part of the contract. The directory is an argument, not a word of a shell
// line: no path of it is interpreted.
func RunPRDLint(script, prdDir string) (code int, stdout, stderr []string, err error) {
	var out, errOut bytes.Buffer
	cmd := exec.Command("sh", script, prdDir)
	cmd.Stdout, cmd.Stderr = &out, &errOut
	runErr := cmd.Run()
	if runErr != nil {
		exit, ok := runErr.(*exec.ExitError)
		if !ok {
			return 0, nil, nil, fmt.Errorf("the pinned prd-lint could not run: %w", runErr)
		}
		code = exit.ExitCode()
	}
	return code, nonEmpty(out.String()), nonEmpty(errOut.String()), nil
}

// nonEmpty gives the non-empty lines of a stream's text.
func nonEmpty(s string) []string {
	var lines []string
	for _, l := range strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n") {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	return lines
}

// WaitingSpecAccepted runs the pinned linter of the target on prdDir and
// reports whether the head is accepted as a specification that waits for its
// bet: the script exits 0, or exits 1 with at least one failure line and with
// every non-empty line of both streams an empty-cell failure. Any other
// outcome refuses: lines carries the non-conforming output lines for the
// report. A script that cannot run or exits otherwise is an error, never an
// acceptance.
func WaitingSpecAccepted(script, prdDir string) (ok bool, lines []string, err error) {
	code, stdout, stderr, err := RunPRDLint(script, prdDir)
	if err != nil {
		return false, nil, err
	}
	if code != 0 && code != 1 {
		return false, nil, fmt.Errorf("the pinned prd-lint exited %d, which is not 0 or 1 (stderr: %s)", code, strings.Join(stderr, "; "))
	}
	ok, lines = waitingVerdict(code, stdout, stderr)
	return ok, lines, nil
}

// waitingVerdict is the pure decision over the producer's exit code and the
// non-empty lines of both streams: exit 0 accepts; exit 1 accepts only when
// every line is a classified empty-cell failure and at least one exists —
// the vacuous exit and any foreign line refuse (review condition C1).
func waitingVerdict(code int, stdout, stderr []string) (bool, []string) {
	if code == 0 {
		return true, nil
	}
	if code != 1 {
		return false, append(append([]string{}, stdout...), stderr...)
	}
	var offend []string
	fails := 0
	for _, l := range append(append([]string{}, stdout...), stderr...) {
		if !emptyCellFailure.MatchString(l) {
			offend = append(offend, l)
		} else {
			fails++
		}
	}
	if fails == 0 {
		return false, append(append([]string{}, stdout...), stderr...)
	}
	return len(offend) == 0, offend
}
