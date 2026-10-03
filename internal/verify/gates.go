package verify

// The rows gate:<kind> (task T-d6q5, D3 of #92, with condition 1 of its plan
// review): an active kind passes when the gate gives pass or clear on the
// setup head, and fail on a commit of the kind's known-bad fixture.

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/catalog"
	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// The calls of the rows, so the unit tests can replace them.
var (
	gateRun   = gate.Run
	fixtureOf = func(stack, kind string) ([]byte, error) {
		e, err := catalog.Read(catalog.Embedded(), stack)
		if err != nil {
			return nil, err
		}
		return e.Fixture(kind)
	}
)

// fixtureWho is the author and the committer of a fixture commit, at the date
// pin.time, so one input gives one commit ID; the commit is on no ref, and no
// push sends it.
var fixtureWho = git.Identity{Name: "layup setup verify", Email: "verify@layup.invalid"}

// gateRows is the state of the rows gate:<kind> of one run: the work area, the
// setup head, and the clean run, which the first active row makes.
type gateRows struct {
	target, head string
	record       work.Record
	out          io.Writer
	clean        *gate.Table
	cleanReason  string // why the clean run gave no table
	left         error  // a scratch tree or directory of a row that stays
}

// row gives the row of the kind k. The first rule that matches decides, as
// the table of setup.md says: the clean run fail gives fail; either run
// not-active gives not-active; a fixture run that is not fail gives fail,
// reason fixture not detected; then pass. A pending kind is clear, with no
// run; a fixture that does not apply is not-active (condition 1).
func (g *gateRows) row(k gate.Kind) Row {
	name := "gate:" + k.Name
	if k.State != "active" {
		return Row{name, "clear", "pending: fixture not run"}
	}
	if g.clean == nil && g.cleanReason == "" {
		fmt.Fprintf(g.out, "layup setup verify: the clean run of the gate on %s\n", g.head)
		g.clean, g.cleanReason = g.run(g.head)
	}
	if g.clean == nil {
		return Row{name, "not-active", "the clean run: " + g.cleanReason}
	}
	if r := kindRow(g.clean, k.Name); r.Result == "fail" || r.Result == "not-active" {
		return Row{name, r.Result, "the clean run: " + r.Reason}
	}
	commit, reason := g.fixture(k.Name)
	if reason != "" {
		return Row{name, "not-active", reason}
	}
	fmt.Fprintf(g.out, "layup setup verify: the fixture run of %s on %s\n", k.Name, commit)
	t, reason := g.run(commit)
	if t == nil {
		return Row{name, "not-active", "the fixture run: " + reason}
	}
	switch r := kindRow(t, k.Name); r.Result {
	case "not-active":
		return Row{name, r.Result, "the fixture run: " + r.Reason}
	case "fail":
		return Row{Check: name, Result: "pass"}
	}
	return Row{name, "fail", "fixture not detected"}
}

// run runs the gate on head with the setup head as its base, and gives its
// table, or why there is none. A scratch tree of the gate that stays is kept
// for the end of the run.
func (g *gateRows) run(head string) (*gate.Table, string) {
	t, err := gateRun(g.target, g.head, head, func(int, int, string) func() { return func() {} }, g.out)
	if cleanup := new(gate.CleanupError); errors.As(err, &cleanup) {
		g.left = errors.Join(g.left, err)
	} else if err != nil {
		return nil, firstLine(err)
	}
	return &t, ""
}

// kindRow gives the row of the kind in t.
func kindRow(t *gate.Table, kind string) gate.Row {
	for _, r := range t.Rows {
		if r.Kind == kind {
			return r
		}
	}
	return gate.Row{Kind: kind, Result: "not-active", Reason: "the gate gave no row of the kind"}
}

// fixture makes the commit of the known-bad fixture of kind on the setup head
// (D3 of #92): in a scratch work tree outside the work area, detached at the
// head, git apply of the patch of the catalog entry of the stack, then a
// commit by fixtureWho at the date of pin.time, on no ref. It gives the
// commit, or the reason of a row not-active.
func (g *gateRows) fixture(kind string) (commit, reason string) {
	stack, missing := value(g.record, "S01", "stack")
	if missing != "" {
		return "", missing
	}
	at, missing := value(g.record, "S02", "pin.time")
	if missing != "" {
		return "", missing
	}
	when, err := time.Parse(timeForm, at)
	if err != nil {
		return "", "record: S02 pin.time is not of the form " + timeForm
	}
	patch, err := fixtureOf(stack, kind)
	if err != nil {
		return "", "fixture: " + firstLine(err)
	}
	scratch, err := tempDir()
	if err != nil {
		return "", "fixture: " + scratchFail
	}
	defer func() {
		if err := removeAll(scratch); err != nil {
			g.left = errors.Join(g.left, &CleanupError{Path: scratch, Err: err})
		}
	}()
	file, tree := filepath.Join(scratch, "fixture.patch"), filepath.Join(scratch, "tree")
	if err := writeFile(file, patch, 0o644); err != nil {
		return "", "fixture: " + scratchFail
	}
	if err := repoAPI.WorktreeAdd(g.target, tree, g.head); err != nil {
		return "", "fixture: " + scratchFail
	}
	defer func() {
		if err := repoAPI.WorktreeRemove(g.target, tree); err != nil {
			g.left = errors.Join(g.left, &CleanupError{Path: tree, Err: err})
		}
	}()
	if err := repoAPI.Apply(tree, file); err != nil {
		return "", "fixture does not apply: " + gitReason(err, "git apply failed")
	}
	who := fixtureWho
	who.Time = when
	if err := repoAPI.CommitAll(tree, "test: the known-bad fixture of the kind "+kind, who); err != nil {
		return "", "fixture: the commit failed: " + gitReason(err, "git commit failed")
	}
	if commit, err = repoAPI.RevParse(tree, "HEAD^{commit}"); err != nil {
		return "", "fixture: the commit failed: " + gitReason(err, "git rev-parse failed")
	}
	return commit, ""
}

// gitReason gives the first line that git wrote on standard error, which names
// the paths of the patch and no path of the scratch tree, so two runs give
// one table; else the text other.
func gitReason(err error, other string) string {
	var failed *git.FailedError
	if errors.As(err, &failed) {
		if s := strings.TrimSpace(failed.Stderr); s != "" {
			line, _, _ := strings.Cut(s, "\n")
			return line
		}
	}
	return other
}
