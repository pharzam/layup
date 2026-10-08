// Package cli is the command-line surface of layup. It holds no state: each
// command reads and writes plain files in a repository (ADR-0011). It is the
// frame of every command (docs/spec/README.md, Commands): one command table,
// the argument rules, the usage, the exit codes and the progress lines.
package cli

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/psb"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/verify"
)

// gateRun runs layup gate; the unit tests replace it.
var gateRun = gate.Run

// readFile reads the FILE of layup psb check and the problem statement of
// layup setup, and isDir reports whether WORK is a directory; the unit tests
// replace them.
var (
	readFile = os.ReadFile
	isDir    = func(p string) bool { fi, err := os.Stat(p); return err == nil && fi.IsDir() }
)

// verifyRun runs layup setup verify, and verifyCheck the one-check call of a
// step of layup setup; the unit tests replace them.
var (
	verifyRun   = verify.Run
	verifyCheck = verify.Check
	// The lists of internal/verify that the steps read (D9 of #90); the unit
	// tests replace them.
	verifyMarkers       = verify.Markers
	verifyBrokenLinks   = verify.BrokenLinks
	verifyFlagged       = verify.Flagged
	verifyLinksBaseline = verify.LinksBaseline
)

// setupRun runs layup setup, and setupSteps gives its steps from the gap table
// of the problem statement and the evidence call; the tests replace them.
var (
	setupRun   = setup.Run
	setupSteps = setup.Steps
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
	{words: []string{"gate"}, args: []string{"REPO"}, flags: []flag{{"base", "REV"}, {"head", "REV"}},
		help: "run the gate kinds of the manifest at --base on --head", run: gateCommand},
	{words: []string{"setup"}, args: []string{"WORK"},
		help: "run the steps of the setup of the target of the work area WORK", run: setupCommand},
	{words: []string{"setup", "verify"}, args: []string{"WORK"},
		help: "check the setup of the target of the work area WORK, from outside", run: setupVerify},
	{words: []string{"run"}, selector: "new", optional: []string{"vision"},
		flags: []flag{{"new", "OWNER/NAME"}, {"host", "DIR"}, {"psb", "FILE"}, {"vision", "FILE"},
			{"operator", "LOGIN"}, {"idea-owner", "LOGIN"}, {"plan", "PLAN"}, {"intake-cap", "MONEY,HOURS"},
			{"lease-h", "MINUTES"}, {"watch-t", "MINUTES"}},
		help: "start a run on the empty repository OWNER/NAME of the forge", run: runNew},
	{words: []string{"run"}, args: []string{"TARGET"}, flags: []flag{{"host", "DIR"}},
		help: "restart the run on the repository TARGET (OWNER/NAME)", run: runAgain},
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

// setupCommand runs layup setup (docs/spec/setup.md): the step table, or the
// stop table of a stop, on standard output; the progress lines and the
// diagnostics on standard error.
func setupCommand(in call) int {
	var brief setup.Brief
	if isDir(in.args[0]) { // else setup.Run names the work area
		var err error
		if brief, err = readBrief(in.args[0]); err != nil {
			fmt.Fprintf(in.stderr, "layup: %v\n", err)
			return exitUsage
		}
	}
	p := newProgress(in.stderr, "setup")
	res, err := setupRun(in.args[0], setupSteps(brief, calls(in.stderr)), setup.Who, func(i, n int, step string) func() {
		p.step(i, n, step)
		return p.end
	})
	p.end()
	if err != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", err)
		return exitUsage
	}
	if werr := res.Write(in.stdout); werr != nil {
		fmt.Fprintf(in.stderr, "layup: %v\n", werr)
		return exitUsage
	}
	if res.Stops != nil {
		return exitStop
	}
	return exitCode(res.Results())
}

// readBrief reads the problem statement of the work area at dir once, runs the
// rules of layup psb check on it, and gives the gap table, as layup psb check
// writes it, and the SHA-256 of the bytes that the rules read (D5 of #86). A
// file that is absent or that is not valid UTF-8 is an input error (K32).
func readBrief(dir string) (setup.Brief, error) {
	src, err := readFile(filepath.Join(dir, filepath.FromSlash(setup.BriefPath)))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return setup.Brief{}, fmt.Errorf("%s: the problem statement is absent", setup.BriefPath)
	case err != nil:
		return setup.Brief{}, fmt.Errorf("%s: %v", setup.BriefPath, err)
	}
	if err := checkBrief(setup.BriefPath, src); err != nil {
		return setup.Brief{}, err
	}
	vision, err := readFile(filepath.Join(dir, filepath.FromSlash(setup.VisionPath)))
	switch {
	case errors.Is(err, fs.ErrNotExist): // the vision brief is optional
	case err != nil:
		return setup.Brief{}, fmt.Errorf("%s: %v", setup.VisionPath, err)
	default:
		if err := checkBrief(setup.VisionPath, vision); err != nil {
			return setup.Brief{}, err
		}
	}
	var table bytes.Buffer
	if err := psb.WriteTSV(&table, psb.Check(string(src))); err != nil {
		return setup.Brief{}, err
	}
	return setup.Brief{Gaps: table.Bytes(), Sum: fmt.Sprintf("%x", sha256.Sum256(src))}, nil
}

// checkBrief refuses a brief at path that is not valid UTF-8 (K32), or that
// holds a marker by the scanner of internal/verify: a raw fact never changes,
// and S11 would have to change it (D4 of #90); or that holds an angle quote
// with no pair, which check markers refuses (fix 3 of the first pilot, #97).
func checkBrief(path string, src []byte) error {
	if n := invalidLine(src); n > 0 {
		return fmt.Errorf("%s: line %d is not valid UTF-8", path, n)
	}
	if m := verify.TextMarkers(string(src)); len(m) > 0 {
		return fmt.Errorf("%s: line %d holds the marker %s, and a brief holds no marker: write the quote another way", path, m[0].Line, m[0].Text)
	}
	if n, q := verify.FirstUnpaired(string(src)); n > 0 {
		return fmt.Errorf("%s: line %d holds the angle quote %s with no pair, and a brief holds no marker: write the quote another way", path, n, q)
	}
	return nil
}

// calls gives the calls of internal/verify that the steps of layup setup read
// (D9 of #90): the one-check call, the markers of a tree with their columns,
// the files whose links break, the files that check adapted flags, the link
// rule of the baseline, and the markers that an input of S14 loses (fix 2 of
// the first pilot, #97).
func calls(out io.Writer) setup.Calls {
	return setup.Calls{
		Checks: evidence(out),
		Markers: func(tree string) ([]setup.Marker, error) {
			marks, err := verifyMarkers(tree)
			var out []setup.Marker
			for _, m := range marks {
				out = append(out, setup.Marker{File: m.File, Line: m.Line, Col: m.Col, Text: m.Text})
			}
			return out, err
		},
		BrokenLinks:   func(tree string) ([]string, error) { return verifyBrokenLinks(tree) },
		Flagged:       func(tree string) ([]string, error) { return verifyFlagged(tree) },
		LinksBaseline: func(source, line string) bool { return verifyLinksBaseline(source, line) },
		LostMarkers:   verify.LostMarkers,
	}
}

// evidence gives the one-check call of a step of layup setup (D12 of #86):
// verify.Check of the named checks on the work area, with its diagnostics on
// out; the reason of the first row that is not pass or clear, or "".
func evidence(out io.Writer) setup.Checks {
	return func(dir string, names []string) string {
		t, err := verifyCheck(dir, names, func(int, int, string) func() { return func() {} }, out)
		if err != nil {
			s, _, _ := strings.Cut(err.Error(), "\n")
			return "the evidence call: " + s
		}
		for _, r := range t.Rows {
			if r.Result != "pass" && r.Result != "clear" {
				return r.Check + ": " + r.Result + ": " + r.Reason
			}
		}
		return ""
	}
}
