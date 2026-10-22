package run

// This file holds the step verify (docs/spec/scaffold.md, The step verify and the
// records of the setup; The README of the Scaffold; row 55 of the plan, task
// T-79zk, #278): layup setup verify on the work area, S15, and the one records
// commit of the setup.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
	workpkg "github.com/pharzam/layup/internal/work"
)

// verifyRun is layup setup verify in process; a variable, so that a test of a
// world puts a stand-in in its place.
var verifyRun = verify.Run

const (
	verifyDoneDetail = "the setup verified and its records committed"
	setupMessage     = "layup run: the records of the setup"
	rulePathsOut     = "out/rule-paths.tsv" // the register that S15 writes
)

// verifyJob is the step with its seams: the loop of the runner on the work area;
// layup setup verify; the write of its table to out/verify.tsv; the three files
// of out/ at their paths in the records, which are there when S15 completed; the
// rows of setup/tree.tsv; and the records commit, which writes the five paths
// and keeps the other records. A failure never reaches the commit (decided here,
// task T-79zk).
type verifyJob struct {
	loop   func() Step
	verify func() (verify.Table, error)
	write  func([]byte) error
	files  func() (map[string][]byte, error)
	tree   func() ([]byte, error)
	commit func(map[string][]byte) error
}

// do runs the job: the loop to the stop O-verify, layup setup verify, its table, the
// first row that is not pass or clear as the fail of the step, the loop again to the
// end of S15, and the commit. An error of layup setup verify is a fail whatever its
// table holds, as a table of passing rows with a scratch tree that is not removed.
func (j verifyJob) do() Step {
	const name = "verify"
	failed := func(detail string, err error) Step {
		return Step{name, "fail", detail + strings.SplitN(err.Error(), "\n", 2)[0]}
	}
	if s := j.loop(); s.Result != "done" {
		return Step{name, "fail", s.Detail}
	}
	t, err := j.verify()
	if err != nil {
		return failed("layup setup verify: ", err)
	}
	var table bytes.Buffer
	if err := t.Write(&table); err != nil {
		return failed("layup setup verify: ", err)
	}
	if err := j.write(table.Bytes()); err != nil {
		return failed("out/verify.tsv: ", err)
	}
	for _, r := range t.Rows {
		if r.Result != "pass" && r.Result != "clear" {
			return Step{name, "fail", fmt.Sprintf("layup setup verify: %s is %s: %s", r.Check, r.Result, r.Reason)}
		}
	}
	if s := j.loop(); s.Result != "done" {
		return Step{name, "fail", s.Detail}
	}
	files, err := j.files()
	if err != nil {
		return failed("the files of out/: ", err)
	}
	tree, err := j.tree()
	if err != nil {
		return failed("setup/tree.tsv: ", err)
	}
	files["setup/tree.tsv"], files["README.md"] = tree, []byte(scaffoldReadme)
	if err := j.commit(files); err != nil {
		return failed("the records commit: ", err)
	}
	return Step{name, "done", verifyDoneDetail}
}

// setupOutputs reads the three files of out/ that the records keep, at their
// paths there. S15 completed when out/record.tsv holds its rows setup.head,
// verify.sha256 and rule-paths.sha256, each hash that of its file, and the
// register holds a path.
func setupOutputs(dir string) (map[string][]byte, error) {
	files := map[string][]byte{}
	for _, f := range []struct{ from, to string }{{workpkg.RecordPath, "setup/record.tsv"}, {workpkg.VerifyPath, "setup/verify.tsv"}, {rulePathsOut, "rule-paths.tsv"}} {
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(f.from)))
		if err != nil {
			return nil, err
		}
		files[f.to] = data
	}
	record, err := workpkg.ReadRecord(dir)
	if err != nil {
		return nil, err
	}
	if v, ok := record.Value("S15", "setup.head"); !ok || v == "" {
		return nil, errors.New("out/record.tsv has no row setup.head of S15: S15 did not complete")
	}
	for _, f := range []struct{ row, path, from string }{{"verify.sha256", "setup/verify.tsv", workpkg.VerifyPath}, {"rule-paths.sha256", "rule-paths.tsv", rulePathsOut}} {
		if v, ok := record.Value("S15", f.row); !ok || v != fmt.Sprintf("%x", sha256.Sum256(files[f.path])) {
			return nil, fmt.Errorf("out/record.tsv: the row %s of S15 is not the hash of %s: S15 did not complete with this file", f.row, f.from)
		}
	}
	rows, err := tsv.Read(files["rule-paths.tsv"], setup.RulePathsSchema)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", rulePathsOut, err)
	} else if len(rows) == 0 {
		return nil, errors.New(rulePathsOut + " holds no path")
	}
	return files, nil
}

// setupTree gives setup/tree.tsv of the target at the path target: the commit of
// main with the step root, then each commit of layup-setup above it, oldest first,
// by first parents, each with its tree and the step of its message
// chore: setup <step>, so the last row is the head that the Operator pushes (the
// order is that of the commits: S14 comes before S11). A message of another form,
// a main that is not an ancestor of layup-setup, or no commit above main is an error.
func setupTree(target string) ([]byte, error) {
	root, err := git.RevParse(target, "refs/heads/main")
	if err != nil {
		return nil, err
	}
	head, err := git.RevParse(target, "refs/heads/layup-setup")
	if err != nil {
		return nil, err
	}
	commits := []string{root}
	for c := head; c != root; {
		commits = append(commits[:1], append([]string{c}, commits[1:]...)...)
		if c, err = git.RevParse(target, c+"^"); err != nil {
			return nil, fmt.Errorf("main is not an ancestor of layup-setup: %w", err)
		}
	}
	if len(commits) == 1 {
		return nil, errors.New("layup-setup holds no commit above main")
	}
	rows := make([][]string, len(commits))
	for i, c := range commits {
		tree, err := git.RevParse(target, c+"^{tree}")
		if err != nil {
			return nil, err
		}
		step := "root"
		if i > 0 {
			msg, err := git.Message(target, c)
			if err != nil {
				return nil, err
			}
			line, _, _ := strings.Cut(msg, "\n")
			if step = strings.TrimPrefix(line, "chore: setup "); step == line || !setupStepForm.MatchString(step) {
				return nil, fmt.Errorf("the commit %s has the message %q, not chore: setup <step>", c, line)
			}
		}
		rows[i] = []string{strconv.Itoa(i), c, tree, step}
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, SetupTreeSchema, rows); err != nil {
		return nil, err
	}
	if _, err := ReadSetupTree(b.Bytes()); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// verifyStep is the step verify of the restart, after setup: "done before" when
// setup/record.tsv exists; else the work area is built again from the records and
// the runner is run on it as setup runs it, layup setup verify and S15 follow, and
// one records commit through fencing updates setup/record.tsv, setup/verify.tsv,
// rule-paths.tsv, setup/tree.tsv and README.md (docs/spec/scaffold.md).
func (r *state) verifyStep(ctx context.Context) Step {
	const name = "verify"
	if recorded, err := r.setupRecorded(ctx); err != nil {
		return fail(name, err)
	} else if recorded {
		return done(name, "done before")
	}
	dir, l, err := r.setupArea(ctx)
	if err != nil {
		return fail(name, err)
	}
	return verifyJob{
		loop: func() Step { return l.do(ctx) },
		verify: func() (verify.Table, error) {
			return verifyRun(dir, func(i, n int, check string) func() {
				if r.cfg.Progress != nil {
					r.cfg.Progress(fmt.Sprintf("verify: [%d/%d] %s", i, n, check))
				}
				return func() {}
			}, io.Discard)
		},
		write: func(table []byte) error {
			p := filepath.Join(dir, filepath.FromSlash(workpkg.VerifyPath))
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				return err
			}
			return os.WriteFile(p, table, 0o644)
		},
		files: func() (map[string][]byte, error) { return setupOutputs(dir) },
		tree:  func() ([]byte, error) { return setupTree(filepath.Join(dir, workpkg.TargetPath)) },
		commit: func(files map[string][]byte) error {
			return r.fenced(ctx, func(ctx context.Context) error { return r.store.Commit(ctx, files, setupMessage) })
		},
	}.do()
}
