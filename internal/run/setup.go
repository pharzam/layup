package run

// This file holds the loop of the step setup (docs/spec/scaffold.md, The stops of
// the runner; The prose rows; The steps of the Scaffold; row 54 of the plan, task
// T-eahr, #277): the step runner of internal/setup is run on the work area until
// it does not stop, and each stop is answered from the records, by a prose session,
// or not at all.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
	workpkg "github.com/pharzam/layup/internal/work"
)

// proseDone is a prose session that ended done: its ID, and its files by their
// paths below inputs/files, as they were committed to the records.
type proseDone struct {
	Session string
	Files   map[string][]byte
}

// setupLoop is the loop of the step with its seams: the runner on the work
// area; the row that a marker gets, for a marker that the batch asked; the
// rows added to inputs/answers.tsv; one session for the F- rows of a stop, which
// gives the files that the session may have committed with any error; the
// removal of a session's files from the records, with the reason; and the bound
// of the runs (decided here, task T-eahr: the loop has no other).
type setupLoop struct {
	run     func() (setup.Result, error)
	markers func(id string) ([]string, bool)
	answer  func(rows [][]string) error
	prose   func(ctx context.Context, rows []setup.StopRow) (proseDone, error)
	remove  func(ctx context.Context, p proseDone, reason string) error
	max     int
}

// setupDoneDetail is the detail of the step when the runner has done S01 to S14.
const setupDoneDetail = "S01 to S14 done, in the work area"

// do runs the loop and gives the row of the step setup. The runner stops on the
// one row O-verify of S15 when it has done S01 to S14 (decided here, task T-eahr):
// the step is then done, and the step verify, which writes out/verify.tsv, runs
// the runner again. A stop with M- rows is answered with the row of each marker
// that the batch asked (a marker no stop lists gets none, and an M- row that the
// batch did not ask is a fail); a stop with F- rows with one prose session; any
// other row is a fail. A fail after a prose session of the run, whichever run
// ends in it, removes the files of the last of them from the records, unless the
// runner has passed its rows since (a later stop with F- rows), and also when
// the session itself ends in an error after its commit; the same stop twice in a
// row is a fail.
func (l setupLoop) do(ctx context.Context) Step {
	const name = "setup"
	var last *proseDone
	var prev []setup.StopRow
	failed := func(detail string) Step {
		if last != nil && len(last.Files) > 0 {
			if err := l.remove(ctx, *last, detail); err != nil {
				detail += "; its files could not be removed from the records: " + strings.SplitN(err.Error(), "\n", 2)[0]
			}
		}
		return Step{name, "fail", detail}
	}
	for i := 0; i < l.max; i++ {
		if err := ctx.Err(); err != nil {
			return failed(err.Error())
		}
		res, err := l.run()
		if err != nil {
			return failed(strings.SplitN(err.Error(), "\n", 2)[0])
		}
		if len(res.Stops) == 0 {
			for _, row := range res.Steps {
				if row.Result != setup.Done && row.Result != setup.Operator {
					return failed(row.Step + ": " + row.Evidence)
				}
			}
			return Step{name, "done", setupDoneDetail}
		}
		rows := res.Stops
		if len(rows) == 1 && rows[0].Step == "S15" && rows[0].Question == "O-verify" {
			return Step{name, "done", setupDoneDetail}
		}
		if reflect.DeepEqual(rows, prev) {
			return failed("the runner stopped on the same rows again: " + questions(rows))
		}
		prev = rows
		var markers, files, other []setup.StopRow
		for _, r := range rows {
			switch {
			case strings.HasPrefix(r.Question, "M-"):
				markers = append(markers, r)
			case strings.HasPrefix(r.Question, "F-"):
				files = append(files, r)
			default:
				other = append(other, r)
			}
		}
		if len(other) > 0 {
			return failed("the runner stopped for " + questions(other) + ", which nothing in the records answers")
		}
		if len(files) > 0 {
			last = nil // the runner passed the rows of an earlier session, whose files stay
		}
		if len(markers) > 0 {
			var added [][]string
			var unasked []setup.StopRow
			for _, r := range markers {
				if row, ok := l.markers(r.Question); ok {
					added = append(added, row)
				} else {
					unasked = append(unasked, r)
				}
			}
			if len(unasked) > 0 {
				return failed("the stop lists " + questions(unasked) + ", which the batch did not ask")
			}
			if err := l.answer(added); err != nil {
				return failed(strings.SplitN(err.Error(), "\n", 2)[0])
			}
		}
		if len(files) > 0 {
			p, err := l.prose(ctx, files)
			last = &p // its files are in the records once it has ended, whatever its end
			if err != nil {
				return failed(strings.SplitN(err.Error(), "\n", 2)[0])
			}
		}
	}
	return failed(fmt.Sprintf("the runner did not end after %d runs", l.max))
}

// questions names the rows of a stop, in its order.
func questions(rows []setup.StopRow) string {
	names := make([]string, len(rows))
	for i, r := range rows {
		names[i] = r.Question
	}
	return strings.Join(names, ", ")
}

// setupCalls gives the calls of internal/verify that the steps of the setup read,
// as internal/cli gives them to layup setup (D9 of #90): the one-check call, the
// markers of a tree with their columns, the files whose links break, the files
// that check adapted flags, the link rule of the baseline, and the markers that an
// input of S14 loses. The run starts no program for the setup: the programs that
// internal/verify starts are its own (docs/spec/packages.md, The table of M2d).
func setupCalls() setup.Calls {
	return setup.Calls{
		Checks: func(dir string, names []string) string {
			t, err := verify.Check(dir, names, func(int, int, string) func() { return func() {} }, io.Discard)
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
		},
		Markers: func(tree string) ([]setup.Marker, error) {
			found, err := verify.Markers(tree)
			marks := make([]setup.Marker, len(found))
			for i, m := range found {
				marks[i] = setup.Marker{File: m.File, Line: m.Line, Col: m.Col, Text: m.Text}
			}
			return marks, err
		},
		BrokenLinks:   verify.BrokenLinks,
		Flagged:       verify.Flagged,
		LinksBaseline: verify.LinksBaseline,
		LostMarkers:   verify.LostMarkers,
	}
}

// setupRunner gives the call of the step runner on the work area at dir, with the
// steps of phase 2 for p; progress, when it is not nil, gets a line for each step
// that runs.
func setupRunner(dir string, p setup.Phase2, progress func(string)) func() (setup.Result, error) {
	steps := setup.StepsPhase2(p, setupCalls())
	return func() (setup.Result, error) {
		return setup.Run(dir, steps, setup.Who, func(i, n int, name string) func() {
			if progress != nil {
				progress(fmt.Sprintf("setup: [%d/%d] %s", i, n, name))
			}
			return func() {}
		})
	}
}

// setupInputs reads what the answers of the setup are made from: the questions,
// the answers and the markers of Intake, and the address of each comment of the
// Intake issue, with that of the comment of the batch.
func (r *state) setupInputs(ctx context.Context) (setupData, error) {
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return setupData{}, errors.New("start.tsv has no number of the Intake issue")
	}
	var d setupData
	if d.Questions, err = r.readRecord(ctx, "questions.tsv", records.ReadQuestions); err != nil {
		return d, err
	}
	if d.Answers, err = r.readRecord(ctx, "intake/answers.tsv", records.ReadAnswers); err != nil {
		return d, err
	}
	if d.Markers, err = r.readRecord(ctx, "intake/markers.tsv", records.ReadMarkers); err != nil {
		return d, err
	}
	d.URL = func(comment string) string {
		return fmt.Sprintf("%s/%s/%s/issues/%d#issuecomment-%s", r.cfg.Register.Web, r.cfg.Owner, r.cfg.Name, issue, comment)
	}
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return d, err
	}
	for _, c := range comments {
		if c.App == r.cfg.Register.AppSlug && strings.HasPrefix(c.Body, batchHeader+"\n") {
			d.Batch = d.URL(strconv.FormatInt(c.ID, 10))
			return d, nil
		}
	}
	return d, errors.New("the Intake issue has no comment of the batch")
}

// workSource reads what the work area is built from, with the rows of
// inputs/answers.tsv given: the briefs of start/, the prose texts of
// setup/inputs/files/ and setup/tree.tsv when the step verify has written it.
func (r *state) workSource(ctx context.Context, rows [][]string) (workSource, error) {
	var src workSource
	var err error
	if src.Brief, err = r.store.ReadFile(ctx, "start/problem-statement.md"); err != nil {
		return src, err
	}
	if src.Brief == nil {
		return src, errors.New("the records have no start/problem-statement.md")
	}
	if src.Vision, err = r.store.ReadFile(ctx, "start/vision.md"); err != nil {
		return src, err
	}
	if src.Tree, err = r.store.ReadFile(ctx, "setup/tree.tsv"); err != nil {
		return src, err
	}
	if src.Files, err = r.store.ReadFiles(ctx, "setup/inputs/files"); err != nil {
		return src, err
	}
	src.Answers = rows
	return src, nil
}

// appendAnswers adds rows to inputs/answers.tsv of the work area at dir.
func appendAnswers(dir string, rows [][]string) error {
	have, err := workpkg.ReadAnswers(dir)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, workpkg.AnswersSchema, append(slices.Clone(have), rows...)); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, filepath.FromSlash(workpkg.AnswersPath)), b.Bytes(), 0o644)
}

// proseSession runs one session of the role product-owner for the F- rows of a
// stop, on the target of the work area at dir at the head of layup-setup (the
// Sessions has that clone and that branch, so the session's repo/ is on its task
// branch at that head); its files are committed to the records by the end of the
// session, and written to inputs/files/ of the work area, where the runner reads
// them before the next run (docs/spec/scaffold.md, The prose rows). No admitted
// pair is an error that starts with "pair:". The files that the session's check
// of its output kept come with any error that follows, as they may be committed.
func (r *state) proseSession(ctx context.Context, dir string, rows []setup.StopRow) (proseDone, error) {
	names, err := proseNames(rows)
	if err != nil {
		return proseDone{}, err
	}
	table, err := stopTable(rows)
	if err != nil {
		return proseDone{}, err
	}
	input := map[string][]byte{"stop.tsv": table}
	if input["problem-statement.md"], err = r.store.ReadFile(ctx, "start/problem-statement.md"); err != nil {
		return proseDone{}, err
	}
	if vision, err := r.store.ReadFile(ctx, "start/vision.md"); err != nil {
		return proseDone{}, err
	} else if vision != nil {
		input["vision.md"] = vision
	}
	s := r.sessions()
	s.Clone, s.Branch = filepath.Join(dir, workpkg.TargetPath), "layup-setup"
	h, m, err := s.PhasePair(ctx, "product-owner", r.cfg.Harnesses, r.cfg.Models, "")
	if errors.Is(err, route.ErrNoPair) {
		return proseDone{}, fmt.Errorf("pair: %v", err)
	} else if err != nil {
		return proseDone{}, err
	}
	base, err := git.RevParse(s.Clone, "refs/heads/layup-setup")
	if err != nil {
		return proseDone{}, err
	}
	var kept map[string][]byte
	keep := proseKeep(names)
	target := r.cfg.Owner + "/" + r.cfg.Name // the name that the prompt gives to README.md (task T-92lm)
	id, err := s.PhaseSession(ctx, PhaseSpec{Step: "setup", Role: "product-owner", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(), Input: input, Output: proseOutput, Files: names,
		Prompt: func(input map[string]string, output string) []byte { return proseText(target, input, output) },
		KeepFiles: func(id string, files map[string][]byte) (map[string][]byte, error) {
			records, err := keep(id, files)
			if err == nil {
				kept = files
			}
			return records, err
		}})
	if err == nil {
		err = copyProse(dir, kept)
	}
	// Whatever the error, once the session's files were kept they may be in the records: the caller removes them.
	return proseDone{Session: id, Files: kept}, err
}

// copyProse writes the files of a prose session below inputs/files/ of the work area at dir.
func copyProse(dir string, files map[string][]byte) error {
	for p, data := range files {
		dst := filepath.Join(dir, "inputs", "files", filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

// removeProse removes the files of a prose session from the records in one
// records commit, through fencing, with its message naming the session and the
// reason (docs/spec/scaffold.md, The stops of the runner).
func (r *state) removeProse(ctx context.Context, p proseDone, reason string) error {
	paths := make([]string, 0, len(p.Files))
	for f := range p.Files {
		paths = append(paths, "setup/inputs/files/"+f)
	}
	slices.Sort(paths)
	reason = strings.SplitN(reason, "\n", 2)[0]
	if len(reason) > 200 {
		reason = reason[:200]
	}
	message := "layup run: the files of the session " + p.Session + ", removed: " + reason
	return Fenced(ctx, r.store, r.cfg.RunID, func(ctx context.Context) error { return r.store.Remove(ctx, paths, message) })
}

// setupStep is the step setup of the restart, after follow-up-answers: "done before"
// when setup/record.tsv exists (the step verify writes it); else the work area at
// DIR/work/OWNER/NAME is built again from the records, with the answers of the setup
// written from the answers of Intake, and the step runner is run on it until it does
// not stop, each stop answered by the table of "The stops of the runner"
// (docs/spec/scaffold.md; row 54 of the plan, task T-eahr, #277).
func (r *state) setupStep(ctx context.Context) Step {
	const name = "setup"
	if recorded, err := r.setupRecorded(ctx); err != nil {
		return fail(name, err)
	} else if recorded {
		return done(name, "done before")
	}
	_, l, err := r.setupArea(ctx)
	if err != nil {
		return fail(name, err)
	}
	return l.do(ctx)
}

// setupRecorded reports whether the records hold setup/record.tsv, which the commit
// of the step verify writes: the steps setup and verify are then done before.
func (r *state) setupRecorded(ctx context.Context) (bool, error) {
	data, err := r.store.ReadFile(ctx, "setup/record.tsv")
	return data != nil, err
}

// setupArea builds the work area at DIR/work/OWNER/NAME again from the records,
// with the answers of the setup written from the answers of Intake, and gives its
// path with the loop of the step runner on it: the beginning of the steps setup and
// verify, which each remove the area and build it again (docs/spec/scaffold.md, The
// work area).
func (r *state) setupArea(ctx context.Context) (string, setupLoop, error) {
	d, err := r.setupInputs(ctx)
	if err != nil {
		return "", setupLoop{}, err
	}
	rows, asked, markers := setupAnswers(d)
	src, err := r.workSource(ctx, rows)
	if err != nil {
		return "", setupLoop{}, err
	}
	host, err := filepath.Abs(r.cfg.Dir)
	if err != nil {
		return "", setupLoop{}, err
	}
	dir := filepath.Join(host, "work", r.cfg.Owner, r.cfg.Name)
	owner, migrate, err := r.codeOwnerInputs(ctx, dir)
	if err != nil {
		return "", setupLoop{}, err
	}
	if err := buildWorkArea(dir, src, r.target, r.defaultBranch, r.url()); err != nil {
		return "", setupLoop{}, err
	}
	if migrate != nil && r.cfg.Progress != nil {
		for _, note := range migrateNotes(dir) {
			r.cfg.Progress(note)
		}
	}
	p := setup.Phase2{Name: r.cfg.Owner + "/" + r.cfg.Name, Visibility: r.value("forge.visibility"), PinSource: r.value("pin.source"),
		PinCommit: r.value("pin.commit"), PinTree: r.value("pin.tree"), PinTime: r.value("pin.time"), BriefSum: r.value("psb.sha256"),
		AppID: r.cfg.Register.AppID, Host: host, DefaultBranch: r.defaultBranch, CodeOwner: owner, Migrate: migrate, Questions: asked}
	l := setupLoop{max: 12, run: setupRunner(dir, p, r.cfg.Progress),
		markers: func(id string) ([]string, bool) { row, ok := markers[id]; return row, ok },
		answer:  func(rows [][]string) error { return appendAnswers(dir, rows) },
		prose: func(ctx context.Context, rows []setup.StopRow) (proseDone, error) {
			return r.proseSession(ctx, dir, rows)
		},
		remove: r.removeProse}
	return dir, l, nil
}

// codeOwnerInputs reads what S13 needs for the code-owner rule, before the work area is touched: the Operator's login from
// approvers.tsv, and, from the two records of the setup, whether the target was set up before the rule (codeOwnerMode).
func (r *state) codeOwnerInputs(ctx context.Context, area string) (string, *setup.Migration, error) {
	login, err := operatorLogin(r.approvers)
	if err != nil {
		return "", nil, err
	}
	record, err := r.store.ReadFile(ctx, "setup/record.tsv")
	if err != nil {
		return "", nil, err
	}
	rulesets, err := r.store.ReadFile(ctx, "setup/rulesets.txt")
	if err != nil {
		return "", nil, err
	}
	migrate, err := codeOwnerMode(r.cfg.Owner+"/"+r.cfg.Name, area, record, rulesets)
	return login, migrate, err
}
