// Package setup is the step runner of layup setup (docs/spec/setup.md, The
// command layup setup): it runs the steps of a target's setup in their order,
// resumes from the setup record, and gives the step table, or the stop table
// of a stop. The steps themselves are rows 9, 13 and 15 of the plan; until
// then they are stubs (Stubs).
package setup

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// StepsSchema is the form of the step table: the block setup-steps of
// docs/spec/setup.md.
var StepsSchema = tsv.Schema{Name: "setup-steps", Location: "stdout", Columns: []tsv.Column{
	{Name: "step", Type: "id(SNN)", Key: true},
	{Name: "actor", Type: "enum(layup-setup|operator|layup-run)"},
	{Name: "result", Type: "enum(done|fail|not-active|operator)"},
	{Name: "evidence", Type: "text"},
}}

// StopSchema is the form of the stop table: the block setup-stop.
var StopSchema = tsv.Schema{Name: "setup-stop", Location: "stdout", Columns: []tsv.Column{
	{Name: "step", Type: "id(SNN)"},
	{Name: "question", Type: "text", Key: true},
	{Name: "ask", Type: "text"},
	{Name: "where", Type: "text"},
}}

// The kinds of an outcome (D2 of #85).
const (
	Done      = "done"       // done and checked: its value rows and its done row
	Stop      = "stop"       // a missing input: the run ends with the stop table
	Fail      = "fail"       // its check failed: the run ends
	NotActive = "not-active" // its check could not run: the run ends
	Operator  = "operator"   // a hand-off: a command in commands.sh waits for the Operator
	Invalid   = "input"      // an input error: exit 2, no table and no row (D11 of #86)
)

// handedOff starts the evidence of a hand-off: what layup setup did, never
// that the Operator ran a command (condition 1 of the plan review of #85).
const handedOff = "handed to the Operator: "

// A StopRow is one missing input of a stop.
type StopRow struct{ Step, Question, Ask, Where string }

// A StepRow is one row of the step table.
type StepRow struct{ Step, Actor, Result, Evidence string }

// An Outcome is how a step ended.
type Outcome struct {
	Kind     string     // one of the kinds above
	Evidence string     // Done, Operator: the evidence line; Fail, NotActive: the reason
	Values   [][]string // Done: the value rows (name, value, source, ref)
	Stops    []StopRow  // Stop: each missing input
	Commit   bool       // Done: the step changed the tree of the target
}

// An Input is what a step reads: the work area and its two records, and the
// identity of the commits that the step makes (O-136).
type Input struct {
	Dir     string
	Record  work.Record
	Answers work.Answers
	Who     git.Identity
}

// A Command is a command for the Operator in commands.sh. Order is its place
// (setup.md, Where the records go in phase 1): 1 the push of the root commit
// (S03), 2 the push of layup-setup (S13), 3 the push of layup-records (S15),
// 4 the apply of the ruleset (S13).
type Command struct {
	Order         int
	Comment, Text string
}

// A Step is one step of a setup.
type Step struct {
	ID       string                        // S01 to S15
	Actor    string                        // who does it in phase 1
	Reads    []string                      // the prefixes of the questions of answers.tsv that it reads
	HandOff  bool                          // its done result is operator (S13)
	Run      func(in Input) Outcome        // the step
	Commands func(r work.Record) []Command // its commands for the Operator, once it is done
	// Unchanged checks, when the step is done, that the inputs that it read
	// did not change (S01: the problem statement, D6 of #86); an error is an
	// input error of the run.
	Unchanged func(r work.Record) error
	// Evidence runs the checks of the step's evidence on the work area at dir,
	// after its commit, and gives the reason of the first that does not pass,
	// or "" (S04 to S14, D12 of #86).
	Evidence func(dir string) string
}

// Order is the run order (D1 of #85, O-123): the prose step, S07, S08, S09 and
// S14, is one group before S10, whose missing inputs make one stop table.
var Order = [][]string{{"S01"}, {"S02"}, {"S03"}, {"S04"}, {"S05"}, {"S06"},
	{"S07", "S08", "S09", "S14"}, {"S10"}, {"S11"}, {"S12"}, {"S13"}, {"S15"}}

// Who is the author and committer of each setup commit and of the records
// commit (O-136); its time is pin.time of S02, so one input gives one commit
// ID.
var Who = git.Identity{Name: "layup-agent[bot]", Email: "335371832+layup-agent[bot]@users.noreply.github.com"}

// A Result is the table of a run: the stop table of a stop, or else the step
// table.
type Result struct {
	Steps []StepRow
	Stops []StopRow
}

// Write writes the stop table of a stop, or else the step table.
func (r Result) Write(w io.Writer) error {
	var rows [][]string
	if r.Stops != nil {
		for _, s := range r.Stops {
			rows = append(rows, []string{s.Step, s.Question, s.Ask, s.Where})
		}
		return tsv.Write(w, StopSchema, rows)
	}
	for _, s := range r.Steps {
		rows = append(rows, []string{s.Step, s.Actor, s.Result, s.Evidence})
	}
	return tsv.Write(w, StepsSchema, rows)
}

// Results gives the result column of the step table; a stop has none.
func (r Result) Results() []string {
	var out []string
	for _, s := range r.Steps {
		out = append(out, s.Result)
	}
	return out
}

// An InputError is an error that the user fixes in the work area: the
// command gives exit code 2 and no table.
type InputError struct{ Err error }

func (e *InputError) Error() string { return e.Err.Error() }
func (e *InputError) Unwrap() error { return e.Err }

// system is the files and the git calls of a run; the unit tests replace it.
type system struct {
	isDir         func(string) bool
	readRecord    func(string) (work.Record, error)
	readAnswers   func(string) (work.Answers, error)
	writeRecord   func(string, work.Record) error
	writeCommands func(string, string) error
	branch        func(dir string) (string, error)
	commit        func(dir, message string, who git.Identity) error
	head          func(dir string) (string, error) // the commit of HEAD
	resetSoft     func(dir, commit string) error   // moves the branch of HEAD to commit (D12 of #86)
	resetHard     func(dir string) error           // puts the index and the work tree back to HEAD (D11 of #90)
	read          func(path string) ([]byte, error)
	now           func() time.Time // the clock of the LAYUP host (pin.time, D7 of #86)
	// The calls of the steps S02 to S04 (steps.go).
	lsRemote       func(url, ref string) (string, error)
	clone          func(url, dir string) error
	checkoutDetach func(dir, commit string) error
	revParse       func(dir, rev string) (string, error)
	rootCommits    func(dir, rev string) ([]string, error)
	message        func(dir, rev string) (string, error)
	initRepo       func(dir string) error
	lsTree         func(dir, rev, path string) ([]git.TreeEntry, error)
	show           func(dir, rev, path string) ([]byte, error)
	switchCreate   func(dir, branch, commit string) error
	exists         func(path string) bool
	removeAll      func(path string) error
	rename         func(from, to string) error
	write          func(path string, data []byte) error
}

var sys = system{
	isDir:       func(p string) bool { fi, err := os.Stat(p); return err == nil && fi.IsDir() },
	readRecord:  work.ReadRecord,
	readAnswers: work.ReadAnswers,
	writeRecord: func(dir string, r work.Record) error {
		var b bytes.Buffer
		if err := tsv.Write(&b, work.RecordSchema, r); err != nil {
			return err
		}
		return writeFile(filepath.Join(dir, filepath.FromSlash(work.RecordPath)), b.Bytes())
	},
	writeCommands: func(dir, text string) error {
		return writeFile(filepath.Join(dir, filepath.FromSlash(work.CommandsPath)), []byte(text))
	},
	branch: git.Branch,
	commit: func(dir, message string, who git.Identity) error {
		if err := git.Add(dir); err != nil {
			return err
		}
		if staged, err := git.Staged(dir); err != nil || !staged { // a step that changes no file makes no commit (D12 of #90)
			return err
		}
		return git.Commit(dir, message, who)
	},
	head:           func(dir string) (string, error) { return git.RevParse(dir, "HEAD^{commit}") },
	resetSoft:      git.ResetSoft,
	resetHard:      git.ResetHard,
	read:           os.ReadFile,
	now:            time.Now,
	lsRemote:       git.LsRemote,
	clone:          git.Clone,
	checkoutDetach: git.CheckoutDetach,
	revParse:       git.RevParse,
	rootCommits:    git.RootCommits,
	message:        git.Message,
	initRepo:       git.Init,
	lsTree:         git.LsTree,
	show:           git.Show,
	switchCreate:   git.SwitchCreate,
	exists:         func(p string) bool { _, err := os.Lstat(p); return err == nil },
	removeAll:      os.RemoveAll,
	rename:         os.Rename,
	write:          writeFile,
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Run runs the steps that the setup record of the work area at dir does not
// have as done, in the run order (D2 of #85). step is called at the start of
// each step that runs and gives a function that the run calls when the step
// ends. who is the identity of the commits, at the time pin.time. An
// *InputError comes with no table.
func Run(dir string, steps map[string]Step, who git.Identity, step func(i, n int, name string) func()) (Result, error) {
	if !sys.isDir(dir) {
		return Result{}, &InputError{fmt.Errorf("the work area %s is not a directory", dir)}
	}
	record, err := sys.readRecord(dir)
	if errors.Is(err, fs.ErrNotExist) {
		record, err = nil, nil // a new work area: no step is done
	}
	if err != nil {
		return Result{}, &InputError{err}
	}
	answers, err := sys.readAnswers(dir)
	if errors.Is(err, fs.ErrNotExist) {
		answers, err = nil, nil // no answer yet; a step that needs one stops for it
	}
	if err != nil {
		return Result{}, &InputError{err}
	}
	if err := checkRows(answers); err != nil {
		return Result{}, &InputError{err}
	}
	done := map[string]bool{}
	for _, r := range record {
		if r[1] == "done" {
			done[r[0]] = true
		}
	}
	for _, id := range ids() { // the other inputs that a done step read do not change (D6 of #86)
		if done[id] && steps[id].Unchanged != nil {
			if err := steps[id].Unchanged(record); err != nil {
				return Result{}, &InputError{err}
			}
		}
	}
	for _, id := range ids() { // the answers that a done step read do not change
		reads := steps[id].Reads
		if !done[id] || len(reads) == 0 {
			continue
		}
		v, ok := record.Value(id, "answers.sha256")
		if !ok {
			return Result{}, &InputError{fmt.Errorf("%s: %s is done and has no row answers.sha256", work.RecordPath, id)}
		}
		if v != work.AnswersHash(id, answers, reads)[2] {
			return Result{}, &InputError{fmt.Errorf("%s: an input that changed after a step read it: the rows %s that %s read",
				work.AnswersPath, strings.Join(reads, " "), id)}
		}
	}
	var todo []string
	for _, g := range Order {
		for _, id := range g {
			if !done[id] {
				todo = append(todo, id)
			}
		}
	}
	ran := map[string]StepRow{} // the steps of this run that did not pass
	var stops []StopRow
	failed, i := "", 0
	for _, g := range Order {
		for _, id := range g {
			if done[id] {
				continue
			}
			i++
			end := step(i, len(todo), id)
			s := steps[id]
			if id >= "S04" && id <= "S14" { // each step starts from the commit of the step before it (D11 of #90)
				if err := sys.resetHard(filepath.Join(dir, work.TargetPath)); err != nil {
					ran[id], failed = StepRow{id, s.Actor, Fail, "the reset of the target failed: " + firstLine(err)}, id
					end()
					break
				}
			}
			o := s.Run(Input{Dir: dir, Record: record, Answers: answers, Who: who})
			switch o.Kind {
			case Invalid:
				end()
				return Result{}, &InputError{errors.New(o.Evidence)}
			case Done, Operator:
				next := slices.Clone(record)
				for _, v := range o.Values {
					next = append(next, append([]string{id}, v...))
				}
				handOff := o.Kind == Operator || s.HandOff
				if handOff && (s.Commands == nil || len(s.Commands(next)) == 0) { // the hand-off is its commands
					ran[id], failed = StepRow{id, s.Actor, Fail, "a hand-off with no command for the Operator"}, id
					break
				}
				target := filepath.Join(dir, work.TargetPath)
				var parent, made string
				if o.Commit && id >= "S04" && id <= "S14" { // S03 and S15 make their own commits
					parent, _ = sys.head(target)
					if reason := commit(dir, id, who, record); reason != "" {
						ran[id], failed = StepRow{id, s.Actor, Fail, reason}, id
						break
					}
					made, _ = sys.head(target)
				}
				if s.Evidence != nil { // the one-check call, after the commit (D12 of #86)
					if reason := s.Evidence(dir); reason != "" {
						// The undo moves layup-setup back only from the commit that
						// this run made (condition 1 of the plan review of #86).
						if cur, err := sys.head(target); err == nil && made != "" && cur == made && parent != "" && parent != made {
							_ = sys.resetSoft(target, parent)
						}
						ran[id], failed = StepRow{id, s.Actor, Fail, reason}, id
						break
					}
				}
				evidence := o.Evidence
				if handOff {
					evidence = handedOff + o.Evidence
				}
				record = next
				if len(s.Reads) > 0 {
					record = append(record, work.AnswersHash(id, answers, s.Reads))
				}
				record = append(record, []string{id, "done", evidence, "step", ""})
				if err := sys.writeRecord(dir, record); err != nil {
					end()
					return Result{}, err
				}
				done[id] = true
			case Stop:
				stops = append(stops, o.Stops...)
			default:
				ran[id], failed = StepRow{id, s.Actor, o.Kind, o.Evidence}, id
			}
			end()
			if failed != "" {
				break
			}
		}
		if failed != "" || stops != nil {
			break
		}
	}
	if err := writeCommands(dir, steps, record); err != nil {
		return Result{}, err
	}
	if failed == "" && stops != nil {
		sortStops(stops)
		return Result{Stops: stops}, nil
	}
	var rows []StepRow
	for _, id := range ids() {
		actor := steps[id].Actor
		row, ok := ran[id]
		switch v, isDone := record.Value(id, "done"); {
		case ok:
		case isDone && strings.HasPrefix(v, handedOff):
			row = StepRow{id, actor, Operator, v}
		case isDone:
			row = StepRow{id, actor, Done, v}
		default:
			row = StepRow{id, actor, NotActive, "not run: " + failed + " did not pass"}
		}
		rows = append(rows, row)
	}
	return Result{Steps: rows}, nil
}

// commit commits the change of step id on the branch layup-setup of the
// target, by who at the time pin.time of the record. It gives the evidence of
// a commit that failed, or "".
func commit(dir, id string, who git.Identity, record work.Record) string {
	const failed = "the commit of the step failed"
	v, _ := record.Value("S02", "pin.time")
	t, err := time.Parse("2006-01-02T15:04:05Z", v)
	if err != nil {
		return failed
	}
	target := filepath.Join(dir, work.TargetPath)
	if b, err := sys.branch(target); err != nil || b != "refs/heads/layup-setup" {
		return failed + ": the target is not on the branch layup-setup"
	}
	who.Time = t
	if sys.commit(target, "chore: setup "+id, who) != nil {
		return failed
	}
	return ""
}

// writeCommands writes commands.sh again from the commands of the done steps,
// in their fixed order, each with its comment line (D7 of #85); with no
// command, the file is empty, so no command of an earlier record stays.
func writeCommands(dir string, steps map[string]Step, record work.Record) error {
	var cmds []Command
	for _, id := range ids() {
		if _, ok := record.Value(id, "done"); ok && steps[id].Commands != nil {
			cmds = append(cmds, steps[id].Commands(record)...)
		}
	}
	sort.SliceStable(cmds, func(i, j int) bool { return cmds[i].Order < cmds[j].Order })
	var b strings.Builder
	for _, c := range cmds {
		fmt.Fprintf(&b, "# %s\n%s\n", c.Comment, c.Text)
	}
	return sys.writeCommands(dir, b.String())
}

// ids gives S01 to S15, the order of the step table.
func ids() []string {
	out := make([]string, 15)
	for i := range out {
		out[i] = fmt.Sprintf("S%02d", i+1)
	}
	return out
}

func hasPrefix(q string, prefixes []string) bool {
	return slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(q, p) })
}

// checkRows is the part of the answers rule that needs no step (D5 of #85): a
// question that no step asks, a question of S01 that is not one of its four,
// and the answer gap, which keeps a marker.
func checkRows(a work.Answers) error {
	for i, r := range a {
		at := fmt.Sprintf("%s: line %d, %s", work.AnswersPath, i+2, r[0])
		switch {
		case strings.HasPrefix(r[0], "F-") || strings.HasPrefix(r[0], "O-"):
			return fmt.Errorf("%s: an F- or O- question is answered by a file or a command, not by a row", at)
		case !hasPrefix(r[0], []string{"S01-", "Q-", "M-"}):
			return fmt.Errorf("%s: no step asks a question of this form", at)
		case strings.HasPrefix(r[0], "S01-") && !slices.Contains(s01, r[0]):
			return fmt.Errorf("%s: S01 does not ask this question", at)
		case r[1] == "gap" && !strings.HasPrefix(r[0], "M-"):
			return fmt.Errorf("%s: the answer gap keeps a marker, and this question is not a marker", at)
		case r[1] == "gap" && r[4] == "":
			return fmt.Errorf("%s: the answer gap needs its question_text", at)
		}
	}
	return nil
}

// CheckAsked gives an error for the first row of answers whose question has
// one of the prefixes and is not one of asked: S01 checks the Q- rows against
// its gap table, and S10 the M- rows against its markers (checkRows checks the
// S01- rows before any step).
func CheckAsked(a work.Answers, prefixes, asked []string) error {
	for i, r := range a {
		if hasPrefix(r[0], prefixes) && !slices.Contains(asked, r[0]) {
			return fmt.Errorf("%s: line %d, %s: no step of this run asked this question", work.AnswersPath, i+2, r[0])
		}
	}
	return nil
}

// MarkerID gives the question ID of a marker: M- and the first 8 hexadecimal
// characters of the SHA-256 of the file, a tab and the marker, so it stays
// while the marker stays (setup.md, The stop table).
func MarkerID(file, marker string) string {
	sum := sha256.Sum256([]byte(file + "\t" + marker))
	return fmt.Sprintf("M-%x", sum[:4])
}

// s01 is the order of the questions of S01 (setup.md, The steps).
var s01 = []string{"S01-stack", "S01-name", "S01-visibility", "S01-baseline"}

// sortStops puts the rows of a stop table in the order of D6 of #85: the S01-
// rows in the order of S01, the Q- rows by the number of their ID, the F- rows
// by path, the M- rows by file and line, and the other rows last.
func sortStops(rows []StopRow) {
	type key struct {
		rank, num int
		text      string
	}
	keyOf := func(r StopRow) key {
		q := r.Question
		switch {
		case strings.HasPrefix(q, "S01-"):
			if n := slices.Index(s01, q); n >= 0 {
				return key{0, n, ""}
			}
			return key{0, len(s01), q}
		case strings.HasPrefix(q, "Q-"):
			n, _ := strconv.Atoi(q[2:])
			return key{1, n, ""}
		case strings.HasPrefix(q, "F-"):
			return key{2, 0, q[2:]}
		case strings.HasPrefix(q, "M-"):
			at, _, _ := strings.Cut(r.Where, " ")
			if i := strings.LastIndex(at, ":"); i >= 0 {
				n, _ := strconv.Atoi(at[i+1:])
				return key{3, n, at[:i]}
			}
			return key{3, 0, at}
		}
		return key{4, 0, q}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		a, b := keyOf(rows[i]), keyOf(rows[j])
		if a.rank != b.rank {
			return a.rank < b.rank
		}
		if a.text != b.text {
			return a.text < b.text
		}
		return a.num < b.num
	})
}

// Stubs gives the steps of this version of layup: each is not-active, not
// built yet, until rows 9, 13 and 15 of the plan build it (D9 of #85).
func Stubs() map[string]Step {
	m := map[string]Step{}
	for _, id := range ids() {
		m[id] = Step{ID: id, Actor: "layup-setup", HandOff: id == "S13", Reads: map[string][]string{"S01": {"S01-", "Q-"}, "S10": {"M-"}}[id],
			Run: func(Input) Outcome { return Outcome{Kind: NotActive, Evidence: "not built yet"} }}
	}
	return m
}
