// Package gate runs the gate kinds of a target from outside the target
// (docs/spec/gate.md): it reads the manifest at the base, makes a scratch work
// tree of the head with the base's gate files, gives one result per kind, and
// removes the tree. It imports internal/tsv and internal/git, and starts only
// sh -c, for the gate commands.
package gate

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
)

// ResultSchema is the form of the table of layup gate: the block gate-result
// of docs/spec/gate.md.
var ResultSchema = tsv.Schema{Name: "gate-result", Location: "stdout", Columns: []tsv.Column{
	{Name: "base", Type: "sha1"},
	{Name: "head", Type: "sha1"},
	{Name: "kind", Type: "id(<word>)", Key: true},
	{Name: "state", Type: "enum(active|pending)"},
	{Name: "result", Type: "enum(pass|fail|not-active|clear)"},
	{Name: "reason", Type: "text"},
}}

// An InputError is an error that the user fixes in the input or on the host:
// a revision, the manifest at the base, git or sh. The command gives exit code
// 2, and there is no table.
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

// A Row is one row of the table: a kind of the manifest and its result.
type Row struct{ Kind, State, Result, Reason string }

// A Table is the result of a run: the base and head commits, and one row per
// kind in the order of the manifest.
type Table struct {
	Base, Head string
	Rows       []Row
}

// Results gives the result column of the table.
func (t Table) Results() []string {
	var out []string
	for _, r := range t.Rows {
		out = append(out, r.Result)
	}
	return out
}

// Write writes the table by the schema gate-result.
func (t Table) Write(w io.Writer) error {
	var rows [][]string
	for _, r := range t.Rows {
		rows = append(rows, []string{t.Base, t.Head, r.Kind, r.State, r.Result, r.Reason})
	}
	return tsv.Write(w, ResultSchema, rows)
}

// The things that a run uses, so the unit tests can replace them.
var (
	repoAPI  gitAPI = realGit{}
	lookPath        = exec.LookPath
	runShell        = shell
)

// Run runs the kinds of the manifest at base on head, in the repository at
// repo (docs/spec/gate.md, The run). step is called at the start of each kind
// and gives a function that the run calls when the kind ends; then the run
// writes the kind's held output to out. An *InputError comes with no table; a
// *CleanupError comes with the whole table.
func Run(repo, base, head string, step func(i, n int, kind string) func(), out io.Writer) (Table, error) {
	if v, err := repoAPI.Version(); err != nil || !git.Supported(v) {
		return Table{}, &InputError{fmt.Errorf("git %s or newer is needed: %v", git.MinVersion, cmpErr(err, "version "+v))}
	}
	if _, err := lookPath("sh"); err != nil {
		return Table{}, &InputError{fmt.Errorf("sh not found: %w", err)}
	}
	var t Table
	var err error
	if t.Base, err = commit(repo, base); err != nil {
		return Table{}, err
	}
	if t.Head, err = commit(repo, head); err != nil {
		return Table{}, err
	}
	data, err := repoAPI.Show(repo, t.Base, "docs/gates.tsv")
	if err != nil {
		return Table{}, &InputError{fmt.Errorf("docs/gates.tsv at the base %s: %v", t.Base, firstLine(err))}
	}
	kinds, err := readManifest(data)
	if err != nil {
		return Table{}, &InputError{fmt.Errorf("docs/gates.tsv at the base %s: %w", t.Base, err)}
	}
	files, err := baseFiles(repo, t.Base, kinds)
	if err != nil {
		return Table{}, err
	}
	// A failure gives each row that needs the failed part a fixed reason,
	// with no part of the error, which can name a scratch path; the error
	// itself goes to standard error. A pending row needs no scratch tree.
	changed, diffErr := repoAPI.DiffNames(repo, t.Base, t.Head)
	if diffErr != nil {
		fmt.Fprintf(out, "layup gate: git diff: %v\n", diffErr)
	}
	scratchReason := ""
	s, err := newScratch(repo, t.Head)
	if err != nil {
		scratchReason = "scratch tree: add failed"
	} else if err = s.overlay(t.Base, data, files); err != nil {
		scratchReason = "scratch tree: overlay failed"
	}
	var tree []string
	if err == nil {
		if tree, err = s.files(); err != nil {
			scratchReason = "scratch tree: overlay failed"
		}
	}
	if err != nil {
		fmt.Fprintf(out, "layup gate: the scratch tree: %v\n", err)
	}
	for i, k := range kinds {
		done := step(i+1, len(kinds), k.Name)
		row, output := Row{Kind: k.Name, State: k.State}, []byte(nil)
		switch {
		case k.State == "pending" && diffErr != nil:
			row.Result, row.Reason = notActive, "diff failed"
		case k.State == "pending":
			row.Result, row.Reason = decidePending(k, changed)
		case scratchReason != "":
			row.Result, row.Reason = notActive, scratchReason
		default:
			row.Result, row.Reason = decideActive(k, tree, lookPath, func(command string) outcome {
				o, b := runShell(s.tree, command)
				output = b
				return o
			})
		}
		done()
		if len(output) > 0 {
			out.Write(output)
			if output[len(output)-1] != '\n' {
				io.WriteString(out, "\n")
			}
		}
		t.Rows = append(t.Rows, row)
	}
	if s != nil {
		if err := s.remove(); err != nil {
			return t, &CleanupError{Path: s.dir, Err: err}
		}
	}
	return t, nil
}

// commit peels rev to a commit in repo; a revision that does not resolve to
// a commit is an input error.
func commit(repo, rev string) (string, error) {
	id, err := repoAPI.RevParse(repo, rev+"^{commit}")
	if err != nil {
		return "", &InputError{fmt.Errorf("the revision %q is not a commit of %s: %v", rev, repo, firstLine(err))}
	}
	return id, nil
}

// firstLine gives the first line of the text of err.
func firstLine(err error) string {
	s, _, _ := strings.Cut(strings.TrimSpace(err.Error()), "\n")
	return s
}

// cmpErr gives err, or the text when err is nil.
func cmpErr(err error, text string) error {
	if err != nil {
		return err
	}
	return errors.New(text)
}
