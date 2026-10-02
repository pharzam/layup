// Package standin builds, at test time, a stand-in of the pinned baseline and
// a work area of a target set up from it (K11 of the plan, D7 of #84). The
// tests of layup setup verify, and of layup setup in later rows, get a
// baseline with no network and with no history of LAYUP, and no file of it is
// in Git. Only test files import this package.
package standin

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// Name is the name of the target of a work area, and PinTime its pin.time.
const (
	Name    = "stand-in-target"
	PinTime = "2026-10-02T09:30:00Z"
)

// who is the author and the committer of each commit of a stand-in, at a
// fixed time, so one input gives one commit ID (NFR-005).
var who = git.Identity{Name: "stand-in", Email: "stand-in@layup.invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

// baselineFiles gives the files of the stand-in baseline at url: the history
// and the phrases of the kit, which a setup removes (S05, S14).
func baselineFiles(url string) map[string]string {
	return map[string]string{
		"README.md":                         "# The stand-in baseline\n\nThis repository is a generic **template** for a team that delivers a specified product.\n",
		"AGENTS.md":                         "# AGENTS.md\n\nAgent context for **Armature**, a stand-in.\n\nA domain-free **template**, not a product.\n",
		"docs/decisions/D-0001-stand-in.md": "# D-0001: a decision of the baseline\n",
		"docs/audit/README.md":              "# The audit of the baseline\n",
		"docs/tasks/T-0001.md":              "# T-0001: a task of the baseline\n",
		"docs/tasks/backlog.md":             "# Backlog\n\n- **T-0001**: a task of the baseline ([#1](" + url + "/issues/1))\n",
		"docs/tasks/completed.md":           "# Completed\n\n- **2026-01-01**: **T-0000**, the first task of the baseline ([#0](" + url + "/issues/0))\n",
	}
}

// The files that the setup by hand writes on layup-setup, by step.
const (
	gates = "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
		"static\tactive\tgo\tgo vet ./...\t./*.go\t—\n" +
		"layout\tpending\tgo\tgo test ./layout/\t./*.go\t—\n"
	readme = "# " + Name + "\n\nThe product repository of " + Name + ", set up from its pinned baseline ([the pin](docs/setup/armature.pin)).\n"
	agents = "# AGENTS.md\n\nAgent context for **" + Name + "**.\n"
)

// Baseline writes the stand-in baseline into a new repository at dir, with
// one commit on main, and gives its file:// URL and the commit.
func Baseline(dir string) (url, commit string, err error) {
	if dir, err = filepath.Abs(dir); err != nil {
		return "", "", err
	}
	url = "file://" + filepath.ToSlash(dir)
	if err := git.Init(dir); err != nil {
		return "", "", err
	}
	if err := commitFiles(dir, baselineFiles(url), "the stand-in baseline"); err != nil {
		return "", "", err
	}
	commit, err = git.RevParse(dir, "HEAD")
	return url, commit, err
}

// A Work is a work area that Make gives: its root, and the URL, the commit
// and the tree of its baseline, which its pin names.
type Work struct{ Dir, URL, Commit, Tree string }

// Options change one part of a work area, for the test of one finding.
type Options struct {
	// Files are the files of a last commit on layup-setup, after the setup
	// by hand: a path and its text; "" removes the path.
	Files map[string]string
	// Record changes the record: "<step> <name>" and the new value; ""
	// removes the row.
	Record map[string]string
}

// Make makes a new stand-in baseline at dir/baseline and, from it, a work
// area at dir/work: the unchanged baseline as the root commit on main (S03),
// a setup by hand on layup-setup, one commit per step (the pin of S04, the
// history removed as S05 does, the manifest of S12, README.md and AGENTS.md of
// S14), out/record.tsv with the rows of S01 and S02, and inputs/answers.tsv.
func Make(dir string, o Options) (Work, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return Work{}, err
	}
	base := filepath.Join(dir, "baseline")
	url, commit, err := Baseline(base)
	if err != nil {
		return Work{}, err
	}
	tree, err := git.RevParse(base, commit+"^{tree}")
	if err != nil {
		return Work{}, err
	}
	w := Work{Dir: filepath.Join(dir, "work"), URL: url, Commit: commit, Tree: tree}
	target := filepath.Join(w.Dir, work.TargetPath)
	if err := git.Init(target); err != nil {
		return Work{}, err
	}
	if err := commitFiles(target, baselineFiles(url), "chore: the unmodified baseline at "+commit); err != nil {
		return Work{}, err
	}
	root, err := git.RevParse(target, "HEAD")
	if err != nil {
		return Work{}, err
	}
	if rootTree, err := git.RevParse(target, root+"^{tree}"); err != nil || rootTree != tree {
		return Work{}, fmt.Errorf("the root commit has the tree %s (%v), not the pinned tree %s", rootTree, err, tree)
	}
	if err := git.SwitchCreate(target, "layup-setup", root); err != nil {
		return Work{}, err
	}
	steps := []struct {
		step  string
		files map[string]string
	}{
		{"S04", map[string]string{"docs/setup/armature.pin": pinText(url, commit, tree)}},
		{"S05", map[string]string{"docs/decisions": "", "docs/audit": "", "docs/tasks/T-0001.md": "",
			"docs/tasks/backlog.md": "# Backlog\n", "docs/tasks/completed.md": "# Completed\n"}},
		{"S12", map[string]string{"docs/gates.tsv": gates}},
		{"S14", map[string]string{"README.md": readme, "AGENTS.md": agents}},
	}
	for _, s := range steps {
		if err := commitFiles(target, s.files, "chore: setup "+s.step); err != nil {
			return Work{}, err
		}
	}
	if len(o.Files) > 0 {
		if err := commitFiles(target, o.Files, "chore: a change of the test"); err != nil {
			return Work{}, err
		}
	}
	var record, answers bytes.Buffer
	if err := tsv.Write(&record, work.RecordSchema, recordRows(url, commit, tree, o.Record)); err != nil {
		return Work{}, err
	}
	if err := tsv.Write(&answers, work.AnswersSchema, answerRows(url)); err != nil {
		return Work{}, err
	}
	return w, write(w.Dir, map[string]string{work.RecordPath: record.String(), work.AnswersPath: answers.String()})
}

// pinText gives the pin file that S04 writes from the pin rows, in the form
// of NFR-006.
func pinText(url, commit, tree string) string {
	return fmt.Sprintf("source=%s\ncommit=%s\ntree=%s\nmethod=git clone %s, checkout %s, .git removed\ndate=%s\n",
		url, commit, tree, url, commit, PinTime[:10])
}

// recordRows gives the rows of S01 and S02 of the record, with the changes;
// the row answers.sha256 of S01 is of the answers of answerRows.
func recordRows(url, commit, tree string, change map[string]string) [][]string {
	rows := [][]string{
		{"S01", "stack", "go", "answer", "S01-stack"},
		{"S01", "name", Name, "answer", "S01-name"},
		{"S01", "visibility", "public", "answer", "S01-visibility"},
		{"S01", "baseline", url, "answer", "S01-baseline"},
		work.AnswersHash("S01", answerRows(url), []string{"S01-", "Q-"}), // the rows that S01 reads (setup.Stubs)
		{"S01", "done", "every answer present; the stack has a catalog entry", "step", ""},
		{"S02", "pin.source", url, "answer", "S01-baseline"},
		{"S02", "pin.commit", commit, "computed", "git ls-remote " + url + " HEAD"},
		{"S02", "pin.tree", tree, "computed", "git rev-parse " + commit + "^{tree}"},
		{"S02", "pin.time", PinTime, "computed", "the clock of the LAYUP host"},
		{"S02", "done", "the commit and the tree", "step", ""},
	}
	var out [][]string
	for _, r := range rows {
		v, ok := change[r[0]+" "+r[1]]
		switch {
		case !ok:
			out = append(out, r)
		case v != "":
			out = append(out, []string{r[0], r[1], v, r[3], r[4]})
		}
	}
	return out
}

// answerRows gives the answers of the questions of S01.
func answerRows(url string) [][]string {
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	return [][]string{
		{"S01-stack", "go", "operator", at, ""},
		{"S01-name", Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""},
		{"S01-baseline", url, "operator", at, ""},
	}
}

// commitFiles writes files into the repository at dir and commits the
// change with message.
func commitFiles(dir string, files map[string]string, message string) error {
	if err := write(dir, files); err != nil {
		return err
	}
	if err := git.Add(dir); err != nil {
		return err
	}
	return git.Commit(dir, message, who)
}

// write writes each file under dir, with its directories; "" removes the
// path.
func write(dir string, files map[string]string) error {
	for name, text := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if text == "" {
			if err := os.RemoveAll(path); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			return err
		}
	}
	return nil
}
