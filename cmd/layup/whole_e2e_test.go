//go:build e2e

package main

// The whole setup (task T-dep6, #93): one shared run of the binary on a
// stand-in baseline with no network, from S01 to S15 (K27), and the scenarios
// that read it: the stops and the end of the run, the table of layup setup
// verify with its seeded defects (NFR-004), the records in the target's Git
// (NFR-001), the pushes of commands.sh, and the target with LAYUP absent
// (NFR-002).

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// The two angle quotes of a marker and the empty value of a table, from their
// code points, so that this file holds no marker (K11).
var (
	mOpen, mClose = string(rune(0x2039)), string(rune(0x203a))
	empty         = string(rune(0x2014))
)

// at is the comment that holds each answer of the shared run.
const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"

// wholeBrief is the problem statement of the shared run: one gap, Q-001 (G4).
const wholeBrief = "# The problem statement of " + standin.Name + "\n\nTechnology stack: Go\n\nThe stand-in target needs one product, and its search must be fast.\n"

// testWho is the author and the committer of each commit that a scenario makes
// itself, at a fixed time.
var testWho = git.Identity{Name: "LAYUP test", Email: "test@layup.invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

// A whole is the shared whole-setup run (D1 of #93). Its directory is in the
// directory of TestMain, so that no scenario removes it when the scenario
// ends. It holds the baseline and its URL, each run of the binary in order,
// and two snapshots of the work area: at the stop O-verify, and at the end.
type whole struct {
	dir, baseline, url string
	env                []string // the go environment of the runs
	runs               []run    // each run of layup setup, in order
	pinCommit          string   // pin.commit after the run that did S02
	records            [2]string
	stopped, done      string // the two snapshots
	verify             result // layup setup verify at the stop O-verify
}

// A run is one stop or end of the shared run: the result of the binary, and
// of a second run with the same input (NFR-005).
type run struct {
	name          string
	first, second result
}

var (
	wholeOnce sync.Once
	wholeRun  *whole
	wholeErr  error
)

// shared gives the shared run; the first scenario that calls it makes it. An
// error ends each scenario that calls it.
func shared(t *testing.T) *whole {
	t.Helper()
	wholeOnce.Do(func() {
		start := time.Now()
		wholeRun, wholeErr = makeWhole()
		fmt.Fprintf(os.Stderr, "the shared whole-setup run: %s\n", time.Since(start).Round(time.Second))
	})
	if wholeErr != nil {
		t.Fatal(wholeErr)
	}
	return wholeRun
}

// hostGoEnv is the environment of a run that starts go: the build cache of the
// host, and no download.
func hostGoEnv() ([]string, error) {
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		return nil, err
	}
	return []string{"GOCACHE=" + strings.TrimSpace(string(cache)), "GOTOOLCHAIN=local", "GOPROXY=off", "GOFLAGS="}, nil
}

// layupIn runs the binary as layupWith does, with its home and its temporary
// directory under dir, outside each work area, and the entries of env.
func layupIn(dir string, env []string, args ...string) (result, error) {
	home, tmp := filepath.Join(dir, "home"), filepath.Join(dir, "tmp")
	for _, d := range []string{home, tmp} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return result{}, err
		}
	}
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	cmd.Env = append([]string{"HOME=" + home, "TMPDIR=" + tmp, "PATH=" + os.Getenv("PATH"), "LC_ALL=C"}, env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	code := 0
	if err := cmd.Run(); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			return result{}, err
		}
		code = exit.ExitCode()
	}
	return result{stdout.String(), stderr.String(), code}, nil
}

// writeFiles writes each file (a path from root, and its text).
func writeFiles(root string, files map[string]string) error {
	for p, text := range files {
		path := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// commitIn stages each change of the repository at dir and commits it.
func commitIn(dir, message string) error {
	if err := git.Add(dir); err != nil {
		return err
	}
	return git.Commit(dir, message, testWho)
}

// opsText is the file of the markers: one marker twice, one that does not
// close on its line, and one that the Operator keeps as a gap.
func opsText() string {
	return "# Ops\n\nThe port is " + mOpen + "port" + mClose + ".\nAgain " + mOpen + "port" + mClose + ".\nThe owner is " +
		mOpen + "owner" + mClose + ".\nThe zone " + mOpen + "zone\n"
}

// scaffold adds the cases of a whole run to the baseline at dir (D2 of #93):
// LAYUP's own link-lint.sh, a note of the history in a blockquote of the
// backlog, a guide whose link the deletion of S05 breaks, the markers, a file
// that check adapted flags, a glossary that it flags, and the baseline's own
// workflow.
func scaffold(dir, url string) error {
	lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "links", "link-lint.sh"))
	if err != nil {
		return err
	}
	if err := writeFiles(dir, map[string]string{
		"docs/links/link-lint.sh":            string(lint),
		"docs/tasks/backlog.md":              "# Backlog\n\n- **T-0001**: a task of the baseline ([#1](" + url + "/issues/1))\n\n> a note of the kit\n> with a link ([#2](" + url + "/issues/2))\n\nKeep this line.\n",
		"docs/guide.md":                      "# Guide\n\nSee [the decision](decisions/D-0001-stand-in.md).\n",
		"docs/ops.md":                        opsText(),
		"docs/how-to.md":                     "# How to\n\nAdapt this to your project.\n",
		"docs/glossary.md":                   "# Glossary\n\nThe words of the kit.\n",
		".github/workflows/ci.yml":           baselineWorkflow,
		"docs/tests/run-discipline-tests.sh": "#!/bin/sh\necho 'run-discipline-tests: the stand-in: 0 failed'\n",
	}); err != nil {
		return err
	}
	return commitIn(dir, "the scaffold of the baseline")
}

// baselineWorkflow is the baseline's own workflow, which stays byte for byte
// after S12.
const baselineWorkflow = "name: ci\n\non:\n  pull_request:\n\njobs:\n  links:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n      - run: sh docs/links/link-lint.sh\n"

// proseInputs are the input files of the prose step, and the fix of the link
// that S05 breaks.
var proseInputs = map[string]string{
	"docs/onboarding-for-engineers.md": "# Onboarding\n\nThe work starts from [the problem statement](facts/problem-statement-brief.md).\n",
	"docs/glossary.md":                 "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository |\n",
	"docs/guardrails.md":               "# Guardrails\n\n## 1. Rules\n\n- **Inv-1** " + empty + " the setup keeps its record. Check: no check yet\n",
	"README.md":                        "# " + standin.Name + "\n\nThe product repository, set up from its pinned baseline ([the pin](docs/setup/armature.pin)). Its records are on the branch `layup-records`.\n",
	"AGENTS.md":                        "# AGENTS.md\n\nAgent context for **" + standin.Name + "**.\n",
	"docs/how-to.md":                   "# How to\n\nRun the tests of the product.\n",
}

// answersText gives the text of answers.tsv for the rows.
func answersText(rows [][]string) (string, error) {
	var b bytes.Buffer
	err := tsv.Write(&b, work.AnswersSchema, rows)
	return b.String(), err
}

// makeWhole makes the shared run: each run of layup setup until its end, with
// the answers and the input files that each stop asks for, and layup setup
// verify at the stop O-verify (D1 to D3 of #93).
func makeWhole() (*whole, error) {
	env, err := hostGoEnv()
	if err != nil {
		return nil, err
	}
	w := &whole{dir: filepath.Join(filepath.Dir(binary), "whole"), env: env}
	w.baseline = filepath.Join(w.dir, "baseline")
	if w.url, _, err = standin.Baseline(w.baseline); err != nil {
		return nil, err
	}
	if err := scaffold(w.baseline, w.url); err != nil {
		return nil, err
	}
	area := filepath.Join(w.dir, "work")
	if err := writeFiles(area, map[string]string{work.BriefPath: wholeBrief}); err != nil {
		return nil, err
	}
	// stop runs the binary on the work area twice; between the two runs, between
	// does what the test needs there (D3: a change of the baseline).
	stop := func(name string, between func() error) error {
		first, err := layupIn(w.dir, w.env, "setup", area)
		if err != nil {
			return err
		}
		if between != nil {
			if err := between(); err != nil {
				return err
			}
		}
		second, err := layupIn(w.dir, w.env, "setup", area)
		if err != nil {
			return err
		}
		w.runs = append(w.runs, run{name, first, second})
		return nil
	}
	answers := [][]string{{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", w.url, "operator", at, ""}}
	give := func(rows ...[]string) error {
		answers = append(answers, rows...)
		text, err := answersText(answers)
		if err != nil {
			return err
		}
		return writeFiles(area, map[string]string{work.AnswersPath: text})
	}
	if err := stop("S01", nil); err != nil {
		return nil, err
	}
	if err := give([]string{"Q-001", "below 200 ms for each search", "idea-owner", at, ""}); err != nil {
		return nil, err
	}
	// The run that does S02 resolves the pin once; a new commit of the baseline
	// after it changes nothing (D3).
	if err := stop("S05", func() error {
		record, err := work.ReadRecord(area)
		if err != nil {
			return err
		}
		w.pinCommit, _ = record.Value("S02", "pin.commit")
		if err := writeFiles(w.baseline, map[string]string{"docs/later.md": "# A later file of the baseline\n"}); err != nil {
			return err
		}
		return commitIn(w.baseline, "a later commit of the baseline")
	}); err != nil {
		return nil, err
	}
	if err := writeFiles(filepath.Join(area, "inputs", "files"), map[string]string{"docs/guide.md": "# Guide\n\nThe decisions of the target are in its records.\n"}); err != nil {
		return nil, err
	}
	if err := stop("the prose step", nil); err != nil {
		return nil, err
	}
	if err := writeFiles(filepath.Join(area, "inputs", "files"), proseInputs); err != nil {
		return nil, err
	}
	if err := stop("S10", nil); err != nil {
		return nil, err
	}
	for _, m := range []struct{ text, answer, question string }{
		{mOpen + "port" + mClose, "8080", ""}, {mOpen + "owner" + mClose, "gap", "Who owns the operations?"}, {mOpen + "zone", "eu-west", ""},
	} {
		if err := give([]string{setup.MarkerID("docs/ops.md", m.text), m.answer, "operator", at, m.question}); err != nil {
			return nil, err
		}
	}
	if err := stop("S15", nil); err != nil {
		return nil, err
	}
	w.stopped = filepath.Join(w.dir, "stopped")
	if err := copyTree(area, w.stopped); err != nil {
		return nil, err
	}
	if w.verify, err = layupIn(w.dir, w.env, "setup", "verify", area); err != nil {
		return nil, err
	}
	if err := writeFiles(area, map[string]string{work.VerifyPath: w.verify.stdout}); err != nil {
		return nil, err
	}
	// The records commit before and after the rerun; "" when there is none,
	// which the scenarios see.
	target := filepath.Join(area, work.TargetPath)
	if err := stop("the end", func() error {
		w.records[0], _ = git.RevParse(target, "refs/heads/layup-records")
		return nil
	}); err != nil {
		return nil, err
	}
	w.records[1], _ = git.RevParse(target, "refs/heads/layup-records")
	w.done = area
	return w, nil
}

// copyTree copies the directory from to the new directory to, byte for byte,
// with the mode of each file.
func copyTree(from, to string) error {
	return filepath.WalkDir(from, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(from, p)
		if err != nil {
			return err
		}
		dst := filepath.Join(to, rel)
		info, err := d.Info()
		if err != nil {
			return err
		}
		switch {
		case d.IsDir():
			return os.MkdirAll(dst, 0o755)
		case info.Mode().Type() == fs.ModeSymlink:
			link, err := os.Readlink(p)
			if err != nil {
				return err
			}
			return os.Symlink(link, dst)
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, data, info.Mode().Perm())
	})
}

// copyOf gives a copy of the work area from, which a scenario may change.
func copyOf(t *testing.T, from string) string {
	t.Helper()
	to := filepath.Join(t.TempDir(), "work")
	if err := copyTree(from, to); err != nil {
		t.Fatal(err)
	}
	return to
}

// gitOut runs a plain git in dir with no setting of the host and no
// background maintenance, and gives its output; an error ends the test.
func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "maintenance.auto=false"}, args...)...)
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "LC_ALL=C"}
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

// stopTable gives a stop table of the rows.
func stopTable(rows ...[]string) string {
	out := "step\tquestion\task\twhere\n"
	for _, r := range rows {
		out += strings.Join(r, "\t") + "\n"
	}
	return out
}

// The demo of #93: on a stand-in baseline with no network, a whole setup
// stops at S01, S05, the prose step, S10 and S15, each stop the same bytes on
// two runs; with the answers, the input files and the table of layup setup
// verify, the last run does S15 and exits 0; a rerun gives the same table and
// no second records commit (D1 to D3).
func TestAWholeSetupOnAStandInBaseline(t *testing.T) {
	w := shared(t)
	var s01 [][]string
	for _, q := range work.S01Questions {
		s01 = append(s01, []string{"S01", q.ID, q.Text, empty})
	}
	ask := func(file, m string) string {
		return "What is the value of " + m + " in " + file + "? Answer gap to keep it as an open gap, with its question as question_text."
	}
	write := func(p string) []string {
		return []string{stepOfInput(p), "F-" + p, "Write the text of " + p + " for the target, and give it as inputs/files/" + p + ".", p}
	}
	port, owner, zone := mOpen+"port"+mClose, mOpen+"owner"+mClose, mOpen+"zone"
	area := w.done
	verifyAsk := "Run layup setup verify '" + area + "' > '" + filepath.Join(area, "out", "verify.tsv") + "', then run layup setup '" + area + "' again."
	want := []struct{ name, stdout string }{
		{"S01", stopTable(append(s01, []string{"S01", "Q-001", `Which number or threshold does "fast" stand for here?`, "inputs/briefs/problem-statement.md:5"})...)},
		{"S05", stopTable([]string{"S05", "F-docs/guide.md", "Fix the links of docs/guide.md that the deletion of the baseline's history breaks, and give the file as inputs/files/docs/guide.md.", "docs/guide.md"})},
		{"the prose step", stopTable(write("AGENTS.md"), write("README.md"), write("docs/glossary.md"), write("docs/guardrails.md"),
			[]string{"S14", "F-docs/how-to.md", "Adapt docs/how-to.md, which check adapted flags, to the target, and give it as inputs/files/docs/how-to.md.", "docs/how-to.md"},
			write("docs/onboarding-for-engineers.md"))},
		{"S10", stopTable([]string{"S10", setup.MarkerID("docs/ops.md", port), ask("docs/ops.md", port), "docs/ops.md:3 " + port},
			[]string{"S10", setup.MarkerID("docs/ops.md", owner), ask("docs/ops.md", owner), "docs/ops.md:5 " + owner},
			[]string{"S10", setup.MarkerID("docs/ops.md", zone), ask("docs/ops.md", zone), "docs/ops.md:6 " + zone})},
		{"S15", stopTable([]string{"S15", "O-verify", verifyAsk, empty})},
	}
	if len(w.runs) != 6 {
		t.Fatalf("the shared run has %d runs; want the five stops and the end", len(w.runs))
	}
	for i, s := range want {
		r := w.runs[i]
		if r.name != s.name || r.first.code != 3 || r.first.stdout != s.stdout {
			t.Errorf("the stop %s: exit %d, stderr %q\n%s\nwant 3 and\n%s", s.name, r.first.code, r.first.stderr, r.first.stdout, s.stdout)
		}
		if r.second.code != r.first.code || r.second.stdout != r.first.stdout {
			t.Errorf("the stop %s on a second run: exit %d\n%s\nwant the same bytes (NFR-005)", s.name, r.second.code, r.second.stdout)
		}
	}
	steps := "step\tactor\tresult\tevidence\n"
	for _, s := range [][]string{{"S01", "done", "every answer present; the stack has a catalog entry"}, {"S02", "done", "the commit and the tree"},
		{"S03", "done", "root tree = pin.tree"}, {"S04", "done", "checks pin and facts"}, {"S05", "done", "checks kit-history and link-lint"},
		{"S06", "done", "check facts"}, {"S07", "done", "check onboarding"}, {"S08", "done", "check glossary"}, {"S09", "done", "check guardrails"},
		{"S10", "done", "every marker has an answer row"}, {"S11", "done", "checks markers, sources and facts"}, {"S12", "done", "checks jobs and gates"},
		{"S13", "operator", "handed to the Operator: the ruleset file, and its commands in commands.sh"}, {"S14", "done", "checks identity and adapted"},
		{"S15", "done", "every row of verify.tsv is pass or clear"}} {
		steps += s[0] + "\tlayup-setup\t" + s[1] + "\t" + s[2] + "\n"
	}
	end := w.runs[5]
	if end.first.code != 0 || end.first.stdout != steps || end.second.code != 0 || end.second.stdout != steps {
		t.Errorf("the end: exit %d and %d, stderr %q\n%s\nwant 0, the same table twice, and\n%s", end.first.code, end.second.code, end.first.stderr, end.first.stdout, steps)
	}
	if w.records[0] == "" || w.records[0] != w.records[1] {
		t.Errorf("layup-records %q before the rerun and %q after it; want one records commit", w.records[0], w.records[1])
	}
	// D3: the pin is resolved once, and the root tree is pin.tree (NFR-006).
	record, err := work.ReadRecord(area)
	if err != nil {
		t.Fatal(err)
	}
	commit, _ := record.Value("S02", "pin.commit")
	tree, _ := record.Value("S02", "pin.tree")
	target := filepath.Join(area, work.TargetPath)
	if head := gitOut(t, w.baseline, "rev-parse", "HEAD"); commit == "" || commit != w.pinCommit || head == commit {
		t.Errorf("pin.commit %q after the run that did S02 %q, the baseline now at %q; want the pin kept while the baseline moved", commit, w.pinCommit, head)
	}
	if got := gitOut(t, target, "rev-parse", "main^{tree}"); got != tree {
		t.Errorf("the root tree %s; want pin.tree %s", got, tree)
	}
	// S12 keeps the baseline's own workflow byte for byte.
	if got := gitOut(t, target, "show", "layup-setup:.github/workflows/ci.yml") + "\n"; got != baselineWorkflow {
		t.Errorf("the baseline's workflow on layup-setup:\n%s\nwant it byte for byte", got)
	}
}

// stepOfInput gives the step that asks for the input file p.
func stepOfInput(p string) string {
	switch p {
	case "docs/onboarding-for-engineers.md":
		return "S07"
	case "docs/glossary.md":
		return "S08"
	case "docs/guardrails.md":
		return "S09"
	}
	return "S14"
}

// verifyRows is the table of layup setup verify on the shared run: each check
// passes, and each kind of the Go entry has its row.
func verifyRows() [][]string {
	var rows [][]string
	for _, c := range []string{"discipline-tests", "pin", "kit-history", "facts", "onboarding", "glossary", "guardrails", "markers", "adapted", "identity",
		"link-lint", "sources", "jobs", "gate:static"} {
		rows = append(rows, []string{c, "pass", empty})
	}
	for _, k := range []string{"layout", "boundary", "contract"} {
		rows = append(rows, []string{"gate:" + k, "clear", "pending: fixture not run"})
	}
	return append(rows, []string{"gate:test", "pass", empty})
}

// tableOf gives the table of layup setup verify of the rows.
func tableOf(rows [][]string) string {
	out := "check\tresult\treason\n"
	for _, r := range rows {
		out += strings.Join(r, "\t") + "\n"
	}
	return out
}

// The table of layup setup verify on a work area that layup setup made: each
// row passes or is clear, and the command exits 0; each seeded defect gives
// fail on its own row, the other rows unchanged, and exit 1; a missing
// baseline script gives not-active (NFR-004), and S15 then refuses that table
// (D4 of #93).
func TestLayupSetupVerifyOnAWholeSetup(t *testing.T) {
	w := shared(t)
	clean := verifyRows()
	if w.verify.code != 0 || w.verify.stdout != tableOf(clean) {
		t.Fatalf("layup setup verify at the stop O-verify: exit %d, stderr %q\n%s\nwant 0 and\n%s", w.verify.code, w.verify.stderr, w.verify.stdout, tableOf(clean))
	}
	if r, err := layupIn(w.dir, w.env, "setup", "verify", w.done); err != nil || r.code != 0 || r.stdout != tableOf(clean) {
		t.Errorf("layup setup verify after S15: %v, exit %d\n%s\nwant 0 and the same table", err, r.code, r.stdout)
	}
	for _, c := range []struct {
		name   string
		seed   func(t *testing.T, area string)
		row    string
		result string
		reason string
	}{
		{"a value row whose answer is not a row of answers.tsv", func(t *testing.T, area string) {
			p := filepath.Join(area, filepath.FromSlash(work.RecordPath))
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			text := regexp.MustCompile(`(?m)^(S12\tmodule\t[^\t]*\tanswer\t)S01-name$`).ReplaceAllString(string(data), "${1}S01-nothing")
			if text == string(data) {
				t.Fatal("no row S12 module to change")
			}
			if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "sources", "fail", "source: S12 module: the answer S01-nothing is not a row of inputs/answers.tsv"},
		{"a marker that a commit adds", func(t *testing.T, area string) {
			target := filepath.Join(area, work.TargetPath)
			if err := writeFiles(target, map[string]string{"docs/new.md": "a " + mOpen + "x" + mClose + " marker\n"}); err != nil {
				t.Fatal(err)
			}
			if err := commitIn(target, "a new marker"); err != nil {
				t.Fatal(err)
			}
		}, "markers", "fail", "unlisted: docs/new.md " + mOpen + "x" + mClose},
		{"a pin tree that differs", func(t *testing.T, area string) {
			p := filepath.Join(area, filepath.FromSlash(work.RecordPath))
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			text := regexp.MustCompile(`(?m)^(S02\tpin\.tree\t)[0-9a-f]+`).ReplaceAllString(string(data), "${1}"+strings.Repeat("1", 40))
			if text == string(data) {
				t.Fatal("no row S02 pin.tree to change")
			}
			if err := os.WriteFile(p, []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		}, "pin", "fail", "tree: the pin names "},
		{"a removed gate job", func(t *testing.T, area string) {
			target := filepath.Join(area, work.TargetPath)
			p := filepath.Join(target, ".github", "workflows", "gates.yml")
			data, err := os.ReadFile(p)
			if err != nil {
				t.Fatal(err)
			}
			i := bytes.Index(data, []byte("\n  test:\n"))
			if i < 0 {
				t.Fatal("no job test in the workflow")
			}
			if err := os.WriteFile(p, data[:i+1], 0o644); err != nil {
				t.Fatal(err)
			}
			if err := commitIn(target, "no job test"); err != nil {
				t.Fatal(err)
			}
		}, "jobs", "fail", "no CI job for the kind test"},
	} {
		area := copyOf(t, w.done)
		c.seed(t, area)
		r, err := layupIn(w.dir, w.env, "setup", "verify", area)
		if err != nil {
			t.Fatal(err)
		}
		got, err := tsv.Read([]byte(r.stdout), work.VerifySchema)
		if err != nil || r.code != 1 || len(got) != len(clean) {
			t.Errorf("%s: %v, exit %d\n%s\nwant 1 and a whole table", c.name, err, r.code, r.stdout)
			continue
		}
		for i, row := range got {
			want := clean[i]
			if row[0] == c.row {
				if row[1] != c.result || !strings.HasPrefix(row[2], c.reason) {
					t.Errorf("%s: the row %q; want %s %s, %s", c.name, row, c.row, c.result, c.reason)
				}
				continue
			}
			if want[1] == "pass" {
				want = []string{want[0], "pass", ""}
			}
			if !slices.Equal(row, want) {
				t.Errorf("%s: the row %q; want %q, unchanged", c.name, row, want)
			}
		}
	}

	// A missing baseline script: not-active, and S15 refuses the table.
	area := copyOf(t, w.stopped)
	target := filepath.Join(area, work.TargetPath)
	if err := os.Remove(filepath.Join(target, "docs", "tests", "run-discipline-tests.sh")); err != nil {
		t.Fatal(err)
	}
	if err := commitIn(target, "no discipline tests"); err != nil {
		t.Fatal(err)
	}
	r, err := layupIn(w.dir, w.env, "setup", "verify", area)
	if err != nil {
		t.Fatal(err)
	}
	if r.code != 1 || !strings.HasPrefix(r.stdout, "check\tresult\treason\ndiscipline-tests\tnot-active\tmissing: docs/tests/run-discipline-tests.sh\n") {
		t.Fatalf("no discipline tests: exit %d\n%s\nwant 1 and discipline-tests not-active", r.code, r.stdout)
	}
	if err := writeFiles(area, map[string]string{work.VerifyPath: r.stdout}); err != nil {
		t.Fatal(err)
	}
	s, err := layupIn(w.dir, w.env, "setup", area)
	if err != nil {
		t.Fatal(err)
	}
	if s.code != 1 || !strings.HasSuffix(s.stdout, "\nS15\tlayup-setup\tfail\tout/verify.tsv: the row discipline-tests is not-active: missing: docs/tests/run-discipline-tests.sh\n") {
		t.Errorf("S15 with a not-active row: exit %d\n%s\nwant 1 and S15 fail", s.code, s.stdout)
	}
}

// The records are in the target's Git (NFR-001, D5 of #93): the orphan branch
// layup-records holds the fixed README and the three tables, each valid and
// equal to its file of the work area; each answer is a fact of a raw answers
// record of layup-setup, with its source; commands.sh holds the four commands
// in their order.
func TestTheRecordsAreInTheTargetsGit(t *testing.T) {
	w := shared(t)
	area := w.done
	target := filepath.Join(area, work.TargetPath)
	cmd := exec.Command("git", "-C", target, "merge-base", "main", "layup-records")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null"}
	if out, err := cmd.Output(); err == nil || len(out) != 0 {
		t.Errorf("git merge-base main layup-records: %q, %v; want no common commit", out, err)
	}
	spec, err := os.ReadFile(filepath.Join("..", "..", "docs", "spec", "setup.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, readme, _ := strings.Cut(string(spec), "\n```text records-readme\n")
	readme, _, _ = strings.Cut(readme, "\n```\n")
	if got := gitOut(t, target, "show", "layup-records:README.md"); readme == "" || got != readme {
		t.Errorf("README.md of layup-records:\n%s\nwant the fixed text of setup.md:\n%s", got, readme)
	}
	for _, f := range []struct {
		branch, out string
		schema      tsv.Schema
	}{
		{"setup/record.tsv", work.RecordPath, work.RecordSchema},
		{"setup/verify.tsv", work.VerifyPath, work.VerifySchema},
		{"rule-paths.tsv", "out/rule-paths.tsv", setup.RulePathsSchema},
	} {
		got := gitOut(t, target, "show", "layup-records:"+f.branch) + "\n"
		data, err := os.ReadFile(filepath.Join(area, filepath.FromSlash(f.out)))
		if _, rerr := tsv.Read([]byte(got), f.schema); err != nil || rerr != nil || got != string(data) {
			t.Errorf("%s of layup-records: %v, %v; want it valid by its schema, and equal to %s", f.branch, err, rerr, f.out)
		}
	}
	if got := gitOut(t, target, "ls-tree", "-r", "--name-only", "layup-records"); got != "README.md\nrule-paths.tsv\nsetup/record.tsv\nsetup/verify.tsv" {
		t.Errorf("the files of layup-records: %q; want the four of S15", got)
	}
	// Each answer is a fact of a raw answers record of layup-setup (S04 the
	// S01- and Q- rows, S11 the M- rows), with its by and its source (K15).
	answers, err := work.ReadAnswers(area)
	if err != nil {
		t.Fatal(err)
	}
	var facts string
	for _, p := range strings.Split(gitOut(t, target, "ls-tree", "--name-only", "layup-setup", "docs/facts/"), "\n") {
		if strings.HasSuffix(p, "-answers.md") {
			facts += gitOut(t, target, "show", "layup-setup:"+p) + "\n"
		}
	}
	for _, a := range answers {
		if !strings.Contains(facts, "`"+a[0]+"` "+a[1]+" "+empty+" by "+a[2]+"; source "+a[3]+"; ") {
			t.Errorf("the answer %q is not a fact of a raw answers record of layup-setup", a)
		}
	}
	if len(answers) != 8 {
		t.Errorf("%d answers; want the four of S01, Q-001 and three markers", len(answers))
	}
	cmds, err := os.ReadFile(filepath.Join(area, filepath.FromSlash(work.CommandsPath)))
	if err != nil {
		t.Fatal(err)
	}
	var lines []string
	for _, l := range strings.Split(strings.TrimSuffix(string(cmds), "\n"), "\n") {
		if !strings.HasPrefix(l, "#") {
			lines = append(lines, l)
		}
	}
	want := []string{"git -C target remote add origin 'https://github.com/" + standin.Name + ".git'", "git -C target push origin main",
		"git -C target push origin layup-setup:main", "git -C target push origin layup-records",
		"gh api --method POST 'repos/" + standin.Name + "/rulesets' --input out/ruleset-default.json"}
	if !slices.Equal(lines, want) {
		t.Errorf("the commands of commands.sh:\n%q\nwant, in this order (S03, S13, S15, S13):\n%q", lines, want)
	}
}

// pushed runs commands.sh of a copy of the work area w (D7 of #93): the remote
// origin is a local bare repository through url.<bare>.insteadOf, which the
// test gives in the environment of the run, and a stub of gh, first on PATH,
// records its arguments with no network. It gives the copy, the bare
// repository and the arguments that gh got.
func pushed(t *testing.T, w *whole) (area, bare, gh string) {
	t.Helper()
	area = copyOf(t, w.done)
	tmp := t.TempDir()
	bare, stub := filepath.Join(tmp, "target.git"), filepath.Join(tmp, "bin")
	gitOut(t, tmp, "init", "-q", "--bare", "-b", "main", bare) // the default branch of the repository on GitHub
	args := filepath.Join(tmp, "gh-args")
	if err := writeFiles(stub, map[string]string{"gh": "#!/bin/sh\nprintf '%s\\n' \"$*\" >> '" + args + "'\n"}); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(stub, "gh"), 0o755); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", filepath.Join("out", "commands.sh"))
	cmd.Dir = area
	cmd.Env = []string{"PATH=" + stub + string(os.PathListSeparator) + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null", "LC_ALL=C", "GIT_CONFIG_COUNT=2",
		"GIT_CONFIG_KEY_0=url.file://" + filepath.ToSlash(bare) + ".insteadOf", "GIT_CONFIG_VALUE_0=https://github.com/" + standin.Name + ".git",
		"GIT_CONFIG_KEY_1=maintenance.auto", "GIT_CONFIG_VALUE_1=false"}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("sh out/commands.sh: %v\n%s", err, out)
	}
	data, err := os.ReadFile(args)
	if err != nil {
		t.Fatal(err)
	}
	return area, bare, string(data)
}

// The pushes of commands.sh (D7 of #93): on a bare repository in the place of
// GitHub, main is the setup head, a fast-forward from the root commit whose
// tree is pin.tree, and layup-records is the records commit with no parent;
// the apply of the ruleset goes to gh; the target's own origin stays the URL of
// GitHub.
func TestThePushesOfCommandsSh(t *testing.T) {
	w := shared(t)
	area, bare, gh := pushed(t, w)
	target := filepath.Join(area, work.TargetPath)
	record, err := work.ReadRecord(area)
	if err != nil {
		t.Fatal(err)
	}
	tree, _ := record.Value("S02", "pin.tree")
	head := gitOut(t, target, "rev-parse", "layup-setup")
	root := gitOut(t, target, "rev-parse", "main")
	if got := gitOut(t, bare, "rev-parse", "main"); got != head {
		t.Errorf("main of the bare repository %s; want the setup head %s", got, head)
	}
	if got := gitOut(t, bare, "rev-list", "--max-parents=0", "main"); got != root || gitOut(t, bare, "rev-parse", root+"^{tree}") != tree {
		t.Errorf("the root of main %s; want the root commit %s, whose tree is pin.tree %s", got, root, tree)
	}
	records := gitOut(t, target, "rev-parse", "layup-records")
	if got := gitOut(t, bare, "rev-parse", "layup-records"); got != records || gitOut(t, bare, "rev-list", "--count", "layup-records") != "1" {
		t.Errorf("layup-records of the bare repository %s; want the records commit %s, with no parent", got, records)
	}
	if want := "api --method POST repos/" + standin.Name + "/rulesets --input out/ruleset-default.json\n"; gh != want {
		t.Errorf("gh got %q; want %q", gh, want)
	}
	if got := gitOut(t, target, "config", "remote.origin.url"); got != "https://github.com/"+standin.Name+".git" {
		t.Errorf("the origin of the target %q; want the URL of GitHub, so the rewrite is the test's own", got)
	}
}

// commandWord is a word that a line under .github/ may not hold (condition 1
// of the plan review of #93): layup or setup-check, delimited by the start or
// the end of the line or a character that is not a letter, a digit, - or _.
var commandWord = regexp.MustCompile(`(^|[^A-Za-z0-9_-])(layup|setup-check)([^A-Za-z0-9_-]|$)`)

// The target is independent of LAYUP (NFR-002, D6 of #93): no file of
// layup-setup is a LAYUP program or setup-check.sh, no line under .github/
// names layup or setup-check, and no required check is a layup/ check; a
// plain clone carries origin/layup-records; and in that clone, with no layup
// on PATH, the gate job of each kind passes or is clear on the setup head, and
// on a commit of a Go file each active kind passes and each pending kind waits
// for the first bet.
func TestTheTargetPassesItsGateWithLAYUPAbsent(t *testing.T) {
	w := shared(t)
	area, bare, _ := pushed(t, w)
	target := filepath.Join(area, work.TargetPath)
	magic := [][]byte{{0x7f, 'E', 'L', 'F'}, {0xfe, 0xed, 0xfa, 0xce}, {0xfe, 0xed, 0xfa, 0xcf}, {0xce, 0xfa, 0xed, 0xfe}, {0xcf, 0xfa, 0xed, 0xfe},
		{0xca, 0xfe, 0xba, 0xbe}, {'M', 'Z'}}
	files := strings.Split(gitOut(t, target, "ls-tree", "-r", "--name-only", "layup-setup"), "\n")
	for _, p := range files {
		if base := filepath.Base(p); base == "layup" || base == "setup-check.sh" {
			t.Errorf("%s: a file of layup-setup is a LAYUP program or LAYUP's setup check", p)
		}
		data, err := os.ReadFile(filepath.Join(target, filepath.FromSlash(p)))
		if err != nil {
			t.Fatal(err)
		}
		if slices.ContainsFunc(magic, func(m []byte) bool { return bytes.HasPrefix(data, m) }) {
			t.Errorf("%s: a file of layup-setup is an executable program", p)
		}
		if !strings.HasPrefix(p, ".github/") {
			continue
		}
		for i, l := range strings.Split(string(data), "\n") {
			l, _, _ = strings.Cut(l, "#")
			if commandWord.MatchString(l) {
				t.Errorf("%s:%d: %q names layup or setup-check", p, i+1, l)
			}
		}
	}
	var contexts []string
	var ruleset struct {
		Rules []struct {
			Parameters struct {
				Checks []struct{ Context string } `json:"required_status_checks"`
			} `json:"parameters"`
		} `json:"rules"`
	}
	var protection struct {
		Checks struct {
			Checks []struct{ Context string } `json:"checks"`
		} `json:"required_status_checks"`
	}
	for p, v := range map[string]any{filepath.Join(area, "out", "ruleset-default.json"): &ruleset, filepath.Join(target, "docs", "setup", "branch-protection.json"): &protection} {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, v); err != nil {
			t.Fatal(err)
		}
	}
	for _, r := range ruleset.Rules {
		for _, c := range r.Parameters.Checks {
			contexts = append(contexts, c.Context)
		}
	}
	for _, c := range protection.Checks.Checks {
		contexts = append(contexts, c.Context)
	}
	if len(contexts) != 10 || slices.ContainsFunc(contexts, func(c string) bool { return strings.HasPrefix(c, "layup/") }) {
		t.Errorf("the required checks %q; want the five kinds twice, and no layup/ check", contexts)
	}

	clone := filepath.Join(t.TempDir(), "clone")
	gitOut(t, filepath.Dir(clone), "clone", "-q", bare, clone)
	if got := gitOut(t, clone, "branch", "-r", "--format=%(refname:short)"); !slices.Contains(strings.Split(got, "\n"), "origin/layup-records") {
		t.Errorf("the branches of a plain clone: %q; want origin/layup-records", got)
	}
	// A PATH with no layup: no directory that holds a file layup.
	var dirs []string
	for _, d := range filepath.SplitList(os.Getenv("PATH")) {
		if _, err := os.Stat(filepath.Join(d, "layup")); err != nil {
			dirs = append(dirs, d)
		}
	}
	path := strings.Join(dirs, string(os.PathListSeparator))
	env := append([]string{"PATH=" + path, "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "LC_ALL=C"}, w.env...)
	if lp, err := exec.LookPath("layup"); err == nil && slices.Contains(dirs, filepath.Dir(lp)) {
		t.Fatalf("layup is on the PATH of the test: %s", lp)
	}
	manifest := gitOut(t, clone, "show", "HEAD:docs/gates.tsv")
	workflow := gitOut(t, clone, "show", "HEAD:.github/workflows/gates.yml")
	var kinds []string
	for _, l := range strings.Split(manifest, "\n")[1:] {
		kind, _, _ := strings.Cut(l, "\t")
		kinds = append(kinds, kind)
		if !strings.Contains(workflow, "\n      - run: sh .github/gates.sh "+kind+"\n") {
			t.Errorf("the workflow has no run line of the gate job of %s", kind)
		}
	}
	job := func(kind, base, head string) (string, int) {
		cmd := exec.Command("sh", ".github/gates.sh", kind)
		cmd.Dir = clone
		cmd.Env = append(slices.Clone(env), "GATE_BASE="+base, "GATE_HEAD="+head)
		out, err := cmd.Output()
		code := 0
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			code = exit.ExitCode()
		} else if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
		return lines[len(lines)-1], code
	}
	root := gitOut(t, clone, "rev-list", "--max-parents=0", "HEAD")
	head := gitOut(t, clone, "rev-parse", "HEAD")
	want := map[string]string{"static": "clear\tno product path", "layout": "clear\tpending: no product path", "boundary": "clear\tpending: no product path",
		"contract": "clear\tpending: no product path", "test": "clear\tno product path"}
	for _, k := range kinds {
		if row, code := job(k, root, head); code != 0 || row != k+"\t"+want[k] {
			t.Errorf("the job of %s on the setup head: %q, exit %d; want %q, exit 0", k, row, code, k+"\t"+want[k])
		}
	}
	if err := writeFiles(clone, map[string]string{
		"widget.go":      "// Package widget is the first product file of the target.\npackage widget\n\n// Double gives twice n.\nfunc Double(n int) int { return 2 * n }\n",
		"widget_test.go": "package widget\n\nimport \"testing\"\n\nfunc TestDouble(t *testing.T) {\n\tif Double(2) != 4 {\n\t\tt.Fatal(\"Double(2) is not 4\")\n\t}\n}\n",
	}); err != nil {
		t.Fatal(err)
	}
	gitOut(t, clone, "add", "--all")
	gitOut(t, clone, "-c", "user.name=LAYUP test", "-c", "user.email=test@layup.invalid", "commit", "-q", "-m", "the first product file")
	change := gitOut(t, clone, "rev-parse", "HEAD")
	want = map[string]string{"static": "pass\t" + empty, "test": "pass\t" + empty}
	for _, k := range kinds {
		row, code := job(k, head, change)
		switch active, ok := want[k]; {
		case ok && (code != 0 || row != k+"\t"+active):
			t.Errorf("the job of the active kind %s on a Go change: %q, exit %d; want %q, exit 0", k, row, code, k+"\t"+active)
		case !ok && (code != 1 || row != k+"\tfail\tpending: product path changed: widget.go"):
			t.Errorf("the job of the pending kind %s on a Go change: %q, exit %d; want fail, pending: product path changed, exit 1", k, row, code)
		}
	}
}
