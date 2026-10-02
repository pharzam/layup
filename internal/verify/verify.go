// Package verify checks the setup of a target from outside the target: the
// checks of layup setup verify (docs/spec/setup.md, The checks of layup setup
// verify). It reads the work area of the target, runs each check on a scratch
// work tree of the head of the branch layup-setup, and gives the table
// setup-verify.
package verify

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// TableSchema is the form of the table of layup setup verify: the block
// setup-verify of docs/spec/setup.md.
var TableSchema = tsv.Schema{Name: "setup-verify", Location: "stdout", Columns: []tsv.Column{
	{Name: "check", Type: "text", Key: true},
	{Name: "result", Type: "enum(pass|fail|not-active|clear)"},
	{Name: "reason", Type: "text"},
}}

// The fixed reasons of a row that did not run.
const (
	notBuilt    = "not built yet"
	scratchFail = "scratch tree: add failed"
)

// A Row is one row of the table.
type Row struct{ Check, Result, Reason string }

// A Table is the rows of a run, in the order of the checks.
type Table struct{ Rows []Row }

// Results gives the result column of the table.
func (t Table) Results() []string {
	var out []string
	for _, r := range t.Rows {
		out = append(out, r.Result)
	}
	return out
}

// Write writes the table by TableSchema.
func (t Table) Write(w io.Writer) error {
	var rows [][]string
	for _, r := range t.Rows {
		rows = append(rows, []string{r.Check, r.Result, r.Reason})
	}
	return tsv.Write(w, TableSchema, rows)
}

// An InputError is an error that the user fixes in the work area or on the
// host: the command gives exit code 2 and no table.
type InputError struct{ Err error }

func (e *InputError) Error() string { return e.Err.Error() }
func (e *InputError) Unwrap() error { return e.Err }

// A CleanupError says that the run could not remove its scratch tree. The
// table is complete; the command gives exit code 2.
type CleanupError struct {
	Path string
	Err  error
}

func (e *CleanupError) Error() string {
	return fmt.Sprintf("the scratch tree %s is not removed: %v", e.Path, e.Err)
}
func (e *CleanupError) Unwrap() error { return e.Err }

// An input is what a check reads: the scratch tree and the two records.
type input struct {
	fsys    fs.FS   // the files of the scratch tree
	area    fs.FS   // the files of the work area
	repo    history // its git history
	record  work.Record
	answers work.Answers
}

// A check is one row of the table of setup.md: run gives its findings, or
// script names the baseline's own script that it runs (D5 of #87); a check
// with neither is one that this version of layup does not have yet.
type check struct {
	name   string
	run    func(in input) []string
	script string
}

// built reports whether this version of layup has the check.
func (c check) built() bool { return c.run != nil || c.script != "" }

// checks is the table of setup.md, in its order (D4 of #84). Rows 10 to 15 of
// the plan add the other checks; row 11 added adapted, row 12 facts,
// onboarding, glossary and guardrails in a target's form, row 10 markers,
// sources and the two scripts.
var checks = []check{
	{name: "discipline-tests", script: discTests},
	{name: "pin", run: checkPin},
	{name: "kit-history", run: checkKitHistory},
	{name: "facts", run: checkFacts},
	{name: "onboarding", run: checkOnboarding},
	{name: "glossary", run: checkGlossary},
	{name: "guardrails", run: checkGuardrails},
	{name: "markers", run: checkMarkers},
	{name: "adapted", run: checkAdapted},
	{name: "identity", run: checkIdentity},
	{name: "link-lint", script: linkLint},
	{name: "sources", run: checkSources},
	{name: "jobs"},
}

// gates is the name, in the one-check call, of each row gate:<kind> (K13).
const gates = "gates"

// gitAPI is the calls of internal/git that a run makes.
type gitAPI interface {
	Version() (string, error)
	RevParse(dir, rev string) (string, error)
	Show(dir, rev, path string) ([]byte, error)
	WorktreeAdd(dir, path, rev string) error
	WorktreeRemove(dir, path string) error
	history(dir string) history
}

type realGit struct{}

func (realGit) Version() (string, error)                 { return git.Version() }
func (realGit) RevParse(dir, rev string) (string, error) { return git.RevParse(dir, rev) }
func (realGit) Show(dir, rev, p string) ([]byte, error)  { return git.Show(dir, rev, p) }
func (realGit) WorktreeAdd(dir, p, rev string) error     { return git.WorktreeAdd(dir, p, rev) }
func (realGit) WorktreeRemove(dir, p string) error       { return git.WorktreeRemove(dir, p) }
func (realGit) history(dir string) history               { return gitHistory{dir} }

// gitHistory is the history of the repository at dir, through internal/git.
type gitHistory struct{ dir string }

func (h gitHistory) IsShallow() (bool, error)       { return git.IsShallow(h.dir) }
func (h gitHistory) RootCommits() ([]string, error) { return git.RootCommits(h.dir, "HEAD") }
func (h gitHistory) Files() ([]string, error)       { return git.LsFiles(h.dir) }
func (h gitHistory) Tree(commit string) (string, error) {
	return git.RevParse(h.dir, commit+"^{tree}")
}

// The things that a run uses, so the unit tests can replace them.
var (
	repoAPI     gitAPI = realGit{}
	readAnswers        = work.ReadAnswers
	readRecord         = work.ReadRecord
	tempRoot           = os.TempDir
	tempDir            = func() (string, error) { return os.MkdirTemp(tempRoot(), "layup-verify-") }
	removeAll          = os.RemoveAll
)

// Run runs each check of the table, and each row gate:<kind>, on the work
// area at dir (layup setup verify WORK). step is called at the start of each
// row and gives a function that the run calls when the row ends; out gets the
// diagnostics. An *InputError comes with no table; a *CleanupError comes with
// the whole table.
func Run(dir string, step func(i, n int, check string) func(), out io.Writer) (Table, error) {
	names := []string{gates}
	for _, c := range checks {
		names = append(names, c.name)
	}
	return Check(dir, names, step, out)
}

// Check runs the named checks on the work area at dir, in the order of the
// table: the one-check call that internal/cli makes as the evidence of a step
// of layup setup (D5 of #84). The name gates gives each row gate:<kind>, from
// the manifest at the head of layup-setup, which no other name reads. A name
// that is not a check is an error, never a row.
func Check(dir string, names []string, step func(i, n int, check string) func(), out io.Writer) (Table, error) {
	if len(names) == 0 {
		return Table{}, errors.New("no check is named")
	}
	for _, n := range names {
		if n != gates && !slices.ContainsFunc(checks, func(c check) bool { return c.name == n }) {
			return Table{}, fmt.Errorf("%q is not a check of layup setup verify", n)
		}
	}
	if v, err := repoAPI.Version(); err != nil || !git.Supported(v) {
		if err == nil {
			err = errors.New("version " + v)
		}
		return Table{}, &InputError{fmt.Errorf("git %s or newer is needed: %v", git.MinVersion, err)}
	}
	answers, err := readAnswers(dir)
	if err != nil {
		return Table{}, &InputError{err}
	}
	record, err := readRecord(dir)
	if err != nil {
		return Table{}, &InputError{err}
	}
	target := filepath.Join(dir, work.TargetPath)
	head, err := repoAPI.RevParse(target, "refs/heads/layup-setup^{commit}")
	if err != nil {
		return Table{}, &InputError{fmt.Errorf("%s: no commit at the branch layup-setup: %s", work.TargetPath, firstLine(err))}
	}
	var rows []check
	for _, c := range checks {
		if slices.Contains(names, c.name) {
			rows = append(rows, c)
		}
	}
	if slices.Contains(names, gates) {
		kinds, err := manifest(target, head)
		if err != nil {
			return Table{}, err
		}
		for _, k := range kinds {
			rows = append(rows, check{name: "gate:" + k.Name}) // row 15 builds the check
		}
	}

	// The scratch tree is outside the work area, so the run changes no file
	// of it (round 1 of #84, finding 2).
	if slices.ContainsFunc(rows, check.built) && within(tempRoot(), dir) {
		return Table{}, &InputError{fmt.Errorf("the temporary directory %s is in the work area %s: set TMPDIR to a directory outside it", tempRoot(), dir)}
	}
	// The scratch tree is added in the step of the first built check, and
	// removed in the step of the last row, so the progress lines cover both
	// (finding 5).
	in := input{record: record, answers: answers, area: os.DirFS(dir)}
	var scratch, tree string
	var added bool
	var addErr, left error
	var t Table
	for i, c := range rows {
		done := step(i+1, len(rows), c.name)
		if c.built() && !added {
			added = true
			if scratch, tree, addErr, left = add(target, head); addErr != nil {
				fmt.Fprintf(out, "layup setup verify: the scratch tree: %v\n", addErr)
			} else {
				in.fsys, in.repo = os.DirFS(tree), repoAPI.history(tree)
			}
		}
		row := Row{Check: c.name, Result: "not-active", Reason: notBuilt}
		switch {
		case !c.built():
		case addErr != nil:
			row.Reason = scratchFail
		case c.script != "":
			row = scriptRow(tree, c.name, c.script, out)
		default:
			row = Row{Check: c.name, Result: "pass"}
			if f := c.run(in); len(f) > 0 {
				row.Result, row.Reason = "fail", f[0]
			}
		}
		if i == len(rows)-1 && added && addErr == nil {
			if left = repoAPI.WorktreeRemove(target, tree); left == nil {
				left = removeAll(scratch)
			}
		}
		done()
		t.Rows = append(t.Rows, row)
	}
	if left != nil {
		return t, &CleanupError{Path: scratch, Err: left}
	}
	return t, nil
}

// within reports whether path is dir or a path under it, each made absolute
// and with its symbolic links resolved where it exists.
func within(path, dir string) bool {
	rel, err := filepath.Rel(resolve(dir), resolve(path))
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// resolve makes p absolute first: filepath.EvalSymlinks keeps a relative path
// relative, which filepath.Rel cannot compare with an absolute one.
func resolve(p string) string {
	if a, err := filepath.Abs(p); err == nil {
		p = a
	}
	if r, err := filepath.EvalSymlinks(p); err == nil {
		return r
	}
	return p
}

// manifest reads docs/gates.tsv at head. A missing or malformed manifest, or
// one with no row, is an input error, as for layup gate (gate.md, NFR-004
// item 4): the run cannot name the rows gate:<kind>.
func manifest(target, head string) ([]gate.Kind, error) {
	data, err := repoAPI.Show(target, head, "docs/gates.tsv")
	if err != nil {
		return nil, &InputError{fmt.Errorf("docs/gates.tsv at the head of layup-setup: %s", firstLine(err))}
	}
	kinds, err := gate.ReadManifest(data)
	if err != nil {
		return nil, &InputError{fmt.Errorf("docs/gates.tsv at the head of layup-setup: %w", err)}
	}
	return kinds, nil
}

// add makes the scratch work tree of head in a new temporary directory, and
// gives the directory and the tree. When the add fails, it removes the
// directory, and gives the error of that removal too: a directory that stays
// is a scratch tree that the run could not remove (finding 3).
func add(target, head string) (scratch, tree string, err, left error) {
	if scratch, err = tempDir(); err != nil {
		return "", "", err, nil
	}
	tree = filepath.Join(scratch, "tree")
	if err = repoAPI.WorktreeAdd(target, tree, head); err != nil {
		return scratch, "", err, removeAll(scratch)
	}
	return scratch, tree, nil, nil
}

// firstLine gives the first line of the text of err: the rest of an error of
// git can be long.
func firstLine(err error) string {
	s, _, _ := strings.Cut(err.Error(), "\n")
	return s
}
