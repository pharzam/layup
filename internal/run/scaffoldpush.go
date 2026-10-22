package run

// This file holds the step scaffold-push (docs/spec/scaffold.md, The commands for the
// Operator and its paragraph "Decided here" of task T-hfn1; row 56 of the plan, #279):
// the print of commands.sh, the wait for the Operator's pushes and applies, and the one
// records commit of setup/rulesets.txt.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pharzam/layup/internal/git"
	workpkg "github.com/pharzam/layup/internal/work"
)

const (
	scaffoldWaiting = "waiting for the Operator's pushes and ruleset apply"
	scaffoldDone    = "the apply of the rulesets recorded"
	appliesMessage  = "layup run: the apply of the rulesets"
)

// applyFiles are the outputs of the Operator's two applies, in the order of the record.
var applyFiles = []string{"default.json", "records.json"}

// scaffoldJob is the step with its seams: the last commit of setup/tree.tsv; the work area,
// made sure of for that commit; the head of the default branch on the forge; the apply outputs
// that exist, by name; the path and the text of commands.sh; the print, a note at a time; and
// the records commit. A wait or a failure never reaches the commit (decided here, task T-hfn1).
type scaffoldJob struct {
	source   string // OWNER/NAME
	last     func() (string, error)
	area     func(last string) error
	head     func() (string, error)
	applies  func() (map[string][]byte, error)
	commands func() (path, text string, err error)
	print    func(note string)
	commit   func(files map[string][]byte) error
}

// do runs the job: the work area is made sure of first, on either outcome; the step waits, and
// prints the rule of the Operator and commands.sh, while the head of the forge is not the last
// commit or an output is missing; else it validates the outputs and commits the record.
func (j scaffoldJob) do() Step {
	const name = "scaffold-push"
	failed := func(prefix string, err error) Step {
		return Step{name, "fail", prefix + strings.SplitN(err.Error(), "\n", 2)[0]}
	}
	last, err := j.last()
	if err != nil {
		return failed("setup/tree.tsv: ", err)
	}
	if err := j.area(last); err != nil {
		return failed("", err)
	}
	head, err := j.head()
	if err != nil {
		return failed("the forge: ", err)
	}
	files, err := j.applies()
	if err != nil {
		return failed("the apply outputs: ", err)
	}
	if head != last || len(files) < len(applyFiles) {
		path, text, err := j.commands()
		if err != nil {
			return failed("commands.sh: ", err)
		}
		j.print("scaffold-push: run " + path + " with your own login, from a clone with no hooks installed")
		for _, line := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
			j.print(line)
		}
		return Step{name, "done", scaffoldWaiting}
	}
	record, err := rulesetsRecord(j.source, files)
	if err != nil {
		return failed("", err)
	}
	if err := j.commit(map[string][]byte{"setup/rulesets.txt": record}); err != nil {
		return failed("the records commit: ", err)
	}
	return Step{name, "done", scaffoldDone}
}

// rulesetsRecord gives setup/rulesets.txt: for default.json and then records.json, the line
// == <file>, the output byte for byte, and one line feed that is the separator. Each output must be
// a JSON object whose source is OWNER/NAME, else the error names the file.
func rulesetsRecord(source string, outputs map[string][]byte) ([]byte, error) {
	var b bytes.Buffer
	for _, f := range applyFiles {
		var out map[string]any
		if err := json.Unmarshal(outputs[f], &out); err != nil || out == nil {
			return nil, fmt.Errorf("%s is not a JSON object", f)
		}
		if got, ok := out["source"].(string); !ok || got != source {
			return nil, fmt.Errorf("%s: its source is %v, not %s", f, out["source"], source)
		}
		b.WriteString("== " + f + "\n")
		b.Write(outputs[f])
		b.WriteString("\n")
	}
	return b.Bytes(), nil
}

// readApplies reads the outputs of the Operator's applies that exist in dir, by name; a name
// that cannot be read, a directory too, is an error that names it.
func readApplies(dir string) (map[string][]byte, error) {
	files := map[string][]byte{}
	for _, f := range applyFiles {
		data, err := os.ReadFile(filepath.Join(dir, f))
		if errors.Is(err, fs.ErrNotExist) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		files[f] = data
	}
	return files, nil
}

// workAreaAgain makes sure of the work area at area: one that is there is used as it is, and one that
// is missing is built again as setup does, with the runner run on it, and its layup-setup must be the
// last commit of setup/tree.tsv.
func (r *state) workAreaAgain(ctx context.Context, area, last string) error {
	if _, err := os.Stat(area); err == nil {
		return nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("the work area: %w", err)
	}
	dir, l, err := r.setupArea(ctx)
	if err != nil {
		return err
	}
	if s := l.do(ctx); s.Result != "done" {
		return errors.New(s.Detail)
	}
	head, err := git.RevParse(filepath.Join(dir, workpkg.TargetPath), "refs/heads/layup-setup")
	if err != nil {
		return err
	}
	if head != last {
		return fmt.Errorf("layup-setup of the work area built again is at %s, and setup/tree.tsv ends at %s", head, last)
	}
	return nil
}

// scaffoldPushStep is the step scaffold-push of the restart, after verify: "done before" when
// setup/rulesets.txt is in the records, else the job on the records, the forge and the host.
func (r *state) scaffoldPushStep(ctx context.Context) Step {
	const name = "scaffold-push"
	if data, err := r.store.ReadFile(ctx, "setup/rulesets.txt"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	host, err := filepath.Abs(r.cfg.Dir)
	if err != nil {
		return fail(name, err)
	}
	area := filepath.Join(host, "work", r.cfg.Owner, r.cfg.Name)
	return scaffoldJob{
		source: r.cfg.Owner + "/" + r.cfg.Name,
		last: func() (string, error) {
			data, err := r.store.ReadFile(ctx, "setup/tree.tsv")
			if err != nil {
				return "", err
			} else if data == nil {
				return "", errors.New("the records have no setup/tree.tsv")
			}
			rows, err := ReadSetupTree(data)
			if err != nil {
				return "", err
			}
			return rows[len(rows)-1][1], nil
		},
		area: func(last string) error { return r.workAreaAgain(ctx, area, last) },
		head: func() (string, error) {
			repo, err := r.cfg.Forge.Repository(ctx)
			return repo.Heads[r.defaultBranch], err
		},
		applies: func() (map[string][]byte, error) {
			return readApplies(filepath.Join(host, "applies", r.cfg.Owner, r.cfg.Name))
		},
		commands: func() (string, string, error) {
			p := filepath.Join(area, "out", "commands.sh")
			data, err := os.ReadFile(p)
			return p, string(data), err
		},
		print: func(note string) {
			if r.cfg.Progress != nil {
				r.cfg.Progress(note)
			}
		},
		commit: func(files map[string][]byte) error {
			return r.fenced(ctx, func(ctx context.Context) error { return r.store.Commit(ctx, files, appliesMessage) })
		},
	}.do()
}
