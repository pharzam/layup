package setup

import (
	"bytes"
	"cmp"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// fakeSys is an in-memory work area: the record, the answers, the commits and
// commands.sh.
type fakeSys struct {
	record             work.Record
	answers            work.Answers
	recordErr, ansErr  error
	commits, calls     []string
	commands           string
	head               string // the branch of the target: "" is layup-setup, "detached" none
	noDir, commitFails bool
}

func (f *fakeSys) install(t *testing.T) {
	t.Helper()
	saved := sys
	t.Cleanup(func() { sys = saved })
	sys = system{
		isDir:       func(string) bool { return !f.noDir },
		readRecord:  func(string) (work.Record, error) { return f.record, f.recordErr },
		readAnswers: func(string) (work.Answers, error) { return f.answers, f.ansErr },
		writeRecord: func(_ string, r work.Record) error { f.record = r; return nil },
		writeCommands: func(_ string, text string) error {
			f.commands = text
			return nil
		},
		branch: func(string) (string, error) {
			if f.head == "detached" {
				return "", errors.New("fatal: ref HEAD is not a symbolic ref")
			}
			return cmp.Or(f.head, "refs/heads/layup-setup"), nil
		},
		commit: func(_, message string, who git.Identity) error {
			f.commits = append(f.commits, fmt.Sprintf("%s by %s <%s> at %s", message, who.Name, who.Email, who.Time.UTC().Format(time.RFC3339)))
			if f.commitFails {
				return errors.New("commit failed")
			}
			return nil
		},
	}
}

var testWho = git.Identity{Name: "LAYUP test", Email: "test@layup.invalid"}

// steps gives a stub of each step: done, unless out names another outcome;
// each call is logged in f.
func steps(f *fakeSys, out map[string]Outcome) map[string]Step {
	m := map[string]Step{}
	for _, id := range ids() {
		o, ok := out[id]
		if !ok {
			o = Outcome{Kind: Done, Evidence: "ok " + id}
		}
		id := id
		m[id] = Step{ID: id, Actor: "layup-setup", HandOff: id == "S13", Reads: reads[id], Commands: commands[id], Run: func(Input) Outcome {
			f.calls = append(f.calls, id)
			return o
		}}
	}
	return m
}

var reads = map[string][]string{"S01": {"S01-", "Q-"}, "S10": {"M-"}}

// commands gives the hand-off, S13, its command for the Operator.
var commands = map[string]func(work.Record) []Command{"S13": func(work.Record) []Command {
	return []Command{{Order: 2, Comment: "the push of layup-setup", Text: "git push origin layup-setup:main"}}
}}

// doneRows gives a record whose steps are done through S<through>, with the
// row answers.sha256 of each done step that reads answers (no answer yet).
func doneRows(through int) work.Record {
	var r work.Record
	for i := 1; i <= through; i++ {
		id := fmt.Sprintf("S%02d", i)
		if reads[id] != nil {
			r = append(r, work.AnswersHash(id, nil, reads[id]))
		}
		r = append(r, []string{id, "done", "ok", "step", ""})
	}
	return append(r, []string{"S02", "pin.time", "2026-10-02T09:30:00Z", "computed", "the clock"})
}

// rowOf gives the row of the record with the step and the name, or nil.
func rowOf(r work.Record, step, name string) []string {
	for _, row := range r {
		if row[0] == step && row[1] == name {
			return row
		}
	}
	return nil
}

func noStep(i, n int, name string) func() { return func() {} }

// A run skips each step with a done row, and runs the others in the run order:
// the prose step S07, S08, S09 and S14 before S10 (D1, K16).
func TestARunResumesAfterItsDoneSteps(t *testing.T) {
	f := &fakeSys{record: doneRows(4)}
	f.install(t)
	var progress []string
	_, err := Run("/w", steps(f, nil), testWho, func(i, n int, name string) func() {
		progress = append(progress, fmt.Sprintf("%d/%d %s", i, n, name))
		return func() {}
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"S05", "S06", "S07", "S08", "S09", "S14", "S10", "S11", "S12", "S13", "S15"}
	if !reflect.DeepEqual(f.calls, want) {
		t.Errorf("the calls %q, want %q", f.calls, want)
	}
	if len(progress) != 11 || progress[0] != "1/11 S05" || progress[5] != "6/11 S14" || progress[10] != "11/11 S15" {
		t.Errorf("the progress %q; want one line per step that runs, in the run order", progress)
	}
}

// Each outcome ends the run as D2 says; after a step that did not pass, each
// later step is not-active; a hand-off writes its done row and goes on.
func TestTheOutcomes(t *testing.T) {
	stops := []StopRow{{"S03", "O-push", "push the root commit", "—"}}
	for _, c := range []struct {
		name  string
		out   map[string]Outcome
		rows  []string // the result and the evidence of S03, S04 and S13
		stops bool
	}{
		{"all done", nil, []string{"done ok S03", "done ok S04", "operator handed to the Operator: ok S13"}, false},
		{"a stop", map[string]Outcome{"S03": {Kind: Stop, Stops: stops}}, nil, true},
		{"a fail", map[string]Outcome{"S03": {Kind: Fail, Evidence: "root tree differs"}},
			[]string{"fail root tree differs", "not-active not run: S03 did not pass", "not-active not run: S03 did not pass"}, false},
		{"a not-active", map[string]Outcome{"S03": {Kind: NotActive, Evidence: "not built yet"}},
			[]string{"not-active not built yet", "not-active not run: S03 did not pass", "not-active not run: S03 did not pass"}, false},
		{"a hand-off", map[string]Outcome{"S13": {Kind: Operator, Evidence: "the push of layup-setup, the apply of the ruleset"}},
			[]string{"done ok S03", "done ok S04", "operator handed to the Operator: the push of layup-setup, the apply of the ruleset"}, false},
	} {
		f := &fakeSys{}
		f.install(t)
		res, err := Run("/w", steps(f, c.out), testWho, noStep)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if c.stops != (res.Stops != nil) || c.stops && !reflect.DeepEqual(res.Stops, stops) {
			t.Errorf("%s: the stops %q, want stops %v", c.name, res.Stops, c.stops)
			continue
		}
		if c.stops {
			continue
		}
		var got []string
		for _, r := range res.Steps {
			if r.Step == "S03" || r.Step == "S04" || r.Step == "S13" {
				got = append(got, r.Result+" "+r.Evidence)
			}
		}
		if !reflect.DeepEqual(got, c.rows) || len(res.Steps) != 15 {
			t.Errorf("%s: the rows %q, want %q (15 rows)", c.name, got, c.rows)
		}
	}
}

// A hand-off records what layup setup did: its done row says "handed to the
// Operator", never that the Operator ran the commands; a resumed run shows it
// as operator and does not run it again (condition 1 of the plan review).
func TestTheHandOff(t *testing.T) {
	f := &fakeSys{}
	f.install(t)
	if _, err := Run("/w", steps(f, map[string]Outcome{"S13": {Kind: Operator, Evidence: "the apply of the ruleset"}}), testWho, noStep); err != nil {
		t.Fatal(err)
	}
	if v, ok := f.record.Value("S13", "done"); !ok || v != "handed to the Operator: the apply of the ruleset" {
		t.Errorf("the done row of S13: %q, %v", v, ok)
	}
	f.calls = nil
	res, err := Run("/w", steps(f, nil), testWho, noStep)
	if err != nil || len(f.calls) != 0 || len(res.Steps) != 15 || res.Steps[12].Result != "operator" {
		t.Errorf("a resumed run: %v, calls %q, the rows %q; want no call and S13 operator", err, f.calls, res.Steps)
	}
	// A hand-off with no command for the Operator fails, with no done row and
	// no commit, so no row says that commands.sh holds a command that it does
	// not hold (finding 3 of round 1).
	for _, c := range []struct {
		name, step string
		out        Outcome
		cmds       func(work.Record) []Command
	}{
		{"S13 with no command", "S13", Outcome{Kind: Done, Evidence: "the ruleset file", Commit: true}, nil},
		{"an outcome operator with an empty list", "S12", Outcome{Kind: Operator, Evidence: "the gate files"}, func(work.Record) []Command { return nil }},
	} {
		f := &fakeSys{record: doneRows(4)}
		f.install(t)
		m := steps(f, map[string]Outcome{c.step: c.out})
		s := m[c.step]
		s.Commands = c.cmds
		m[c.step] = s
		res, err := Run("/w", m, testWho, noStep)
		n, _ := strconv.Atoi(c.step[1:])
		if _, done := f.record.Value(c.step, "done"); err != nil || done || len(f.commits) != 0 || len(res.Steps) != 15 ||
			res.Steps[n-1] != (StepRow{c.step, "layup-setup", "fail", "a hand-off with no command for the Operator"}) {
			t.Errorf("%s: %v, done %v, commits %q, the rows %q; want %s fail, no commit and no done row", c.name, err, done, f.commits, res.Steps, c.step)
		}
	}
}

// The prose step is one group: its missing inputs make one stop table, and a
// step of it that is done keeps its done row (D1, O-123).
func TestTheProseGroupStopsOnce(t *testing.T) {
	f := &fakeSys{record: doneRows(6)}
	f.install(t)
	res, err := Run("/w", steps(f, map[string]Outcome{
		"S07": {Kind: Stop, Stops: []StopRow{{"S07", "F-docs/onboarding-for-engineers.md", "write it", "docs/onboarding-for-engineers.md"}}},
		"S14": {Kind: Stop, Stops: []StopRow{{"S14", "F-README.md", "write it", "README.md"}}},
	}), testWho, noStep)
	if err != nil {
		t.Fatal(err)
	}
	var q []string
	for _, s := range res.Stops {
		q = append(q, s.Question)
	}
	if !reflect.DeepEqual(q, []string{"F-README.md", "F-docs/onboarding-for-engineers.md"}) || !reflect.DeepEqual(f.calls, []string{"S07", "S08", "S09", "S14"}) {
		t.Errorf("the stop rows %q, the calls %q; want one table for S07 and S14, and no S10", q, f.calls)
	}
	if _, ok := f.record.Value("S08", "done"); !ok {
		t.Error("S08 of the group is done, and has no done row")
	}
}

// The stop table is in the order of D6: S01- in the order of S01, Q- by the
// number of the ID, F- by path, M- by file and line, O- last (K33, note 4).
func TestTheStopTableOrder(t *testing.T) {
	rows := []StopRow{{"S01", "Q-1000", "a", "inputs/briefs/problem-statement.md:9"}, {"S01", "S01-baseline", "a", "—"},
		{"S10", "M-bbbbbbbb", "a", "b.md:2 x"}, {"S07", "F-docs/z.md", "a", "docs/z.md"}, {"S01", "Q-999", "a", "inputs/briefs/problem-statement.md:0"},
		{"S01", "S01-stack", "a", "—"}, {"S15", "O-verify", "a", "—"}, {"S10", "M-aaaaaaaa", "a", "a.md:10 x"},
		{"S14", "F-README.md", "a", "README.md"}, {"S01", "S01-name", "a", "—"}, {"S10", "M-cccccccc", "a", "a.md:9 x"}}
	sortStops(rows)
	var got []string
	for _, r := range rows {
		got = append(got, r.Question)
	}
	want := []string{"S01-stack", "S01-name", "S01-baseline", "Q-999", "Q-1000", "F-README.md", "F-docs/z.md", "M-cccccccc", "M-aaaaaaaa", "M-bbbbbbbb", "O-verify"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the order\n got %q\nwant %q", got, want)
	}
}

// Each table is a record of its block; two runs on one input give the same
// bytes (NFR-005).
func TestTheTables(t *testing.T) {
	var out [2]bytes.Buffer
	for i := range out {
		f := &fakeSys{}
		f.install(t)
		res, err := Run("/w", steps(f, map[string]Outcome{"S02": {Kind: Fail, Evidence: "no commit"}}), testWho, noStep)
		if err != nil || res.Write(&out[i]) != nil {
			t.Fatal(err)
		}
	}
	if out[0].String() != out[1].String() || !strings.HasPrefix(out[0].String(), "step\tactor\tresult\tevidence\nS01\tlayup-setup\tdone\tok S01\nS02\tlayup-setup\tfail\tno commit\n") {
		t.Errorf("the step table:\n%s", out[0].String())
	}
	var b bytes.Buffer
	if err := (Result{Stops: []StopRow{{"S01", "S01-stack", "Which stack?", ""}}}).Write(&b); err != nil || b.String() != "step\tquestion\task\twhere\nS01\tS01-stack\tWhich stack?\t—\n" {
		t.Errorf("the stop table %q, %v", b.String(), err)
	}
}

func answers(rows ...string) work.Answers {
	var a work.Answers
	for _, r := range rows {
		a = append(a, strings.Split(r, "|"))
	}
	return a
}

// The rows that no step can ask, and the rule of gap, are input errors before
// any step runs (D5).
func TestTheAnswersRule(t *testing.T) {
	good := []string{"S01-stack|go|operator|u|", "Q-001|ten|idea-owner|u|", "M-0123abcd|gap|operator|u|Which port?"}
	for _, c := range []struct {
		name, row, want string
	}{
		{"an F- row", "F-README.md|x|operator|u|", "F-README.md"},
		{"an O- row", "O-verify|x|operator|u|", "O-verify"},
		{"another prefix", "X-1|x|operator|u|", "X-1"},
		{"gap for a question that is not a marker", "Q-002|gap|idea-owner|u|Why?", "Q-002"},
		{"gap with no question text", "M-0000abcd|gap|operator|u|", "M-0000abcd"},
		{"an S01- question that S01 does not ask (finding 4 of round 1)", "S01-wrong|x|operator|u|", "S01-wrong: S01 does not ask this question"},
	} {
		f := &fakeSys{answers: answers(append(good, c.row)...)}
		f.install(t)
		var in *InputError
		if _, err := Run("/w", steps(f, nil), testWho, noStep); !errors.As(err, &in) || !strings.Contains(err.Error(), c.want) || len(f.calls) != 0 {
			t.Errorf("%s: %v, calls %q; want an input error that names %s, before any step", c.name, err, f.calls, c.want)
		}
	}
	f := &fakeSys{answers: answers(good...)}
	f.install(t)
	if _, err := Run("/w", steps(f, nil), testWho, noStep); err != nil {
		t.Errorf("good answers: %v", err)
	}
}

// A step checks the rows of its prefixes against the questions that it asked.
func TestCheckAsked(t *testing.T) {
	a := answers("S01-stack|go|operator|u|", "Q-001|x|idea-owner|u|", "Q-005|x|idea-owner|u|", "M-0123abcd|8080|operator|u|")
	if err := CheckAsked(a, []string{"S01-", "Q-"}, []string{"S01-stack", "Q-001", "Q-002"}); err == nil || !strings.Contains(err.Error(), "line 4") || !strings.Contains(err.Error(), "Q-005") {
		t.Errorf("an unasked Q- row: %v; want an error that names line 4 and Q-005", err)
	}
	if err := CheckAsked(a, []string{"S01-", "Q-"}, []string{"S01-stack", "Q-001", "Q-005"}); err != nil {
		t.Errorf("the asked rows: %v", err)
	}
}

// The answers that a done step read do not change, go or come: each run
// compares them with the step's answers.sha256 before any step runs
// (condition 2 of the plan review).
func TestTheAnswersOfADoneStep(t *testing.T) {
	base := []string{"S01-stack|go|operator|u|", "Q-001|ten|idea-owner|u|", "M-0123abcd|8080|operator|u|"}
	for _, c := range []struct {
		name      string
		rows      []string
		ok1, ok15 bool // the result when the steps are done through S01, and through S15
	}{
		{"unchanged", base, true, true},
		{"a changed Q- row", []string{base[0], "Q-001|eleven|idea-owner|u|", base[2]}, false, false},
		{"a removed Q- row", []string{base[0], base[2]}, false, false},
		{"a new Q- row", append(append([]string{}, base...), "Q-002|x|idea-owner|u|"), false, false},
		{"a changed M- row: S01 does not read it, S10 does", []string{base[0], base[1], "M-0123abcd|9090|operator|u|"}, true, false},
	} {
		for _, through := range []int{1, 15} {
			ok := map[int]bool{1: c.ok1, 15: c.ok15}[through]
			f := &fakeSys{answers: answers(base...)}
			f.install(t)
			out := map[string]Outcome{}
			for i := through + 1; i <= 15; i++ {
				out[fmt.Sprintf("S%02d", i)] = Outcome{Kind: NotActive, Evidence: "not built yet"}
			}
			if _, err := Run("/w", steps(f, out), testWho, noStep); err != nil {
				t.Fatal(err)
			}
			// The ref of the row follows the record's rule for a hash
			// (finding 5 of round 1).
			if row := rowOf(f.record, "S01", "answers.sha256"); row == nil || row[3] != "computed" || row[4] != "sha256 inputs/answers.tsv S01- Q-" {
				t.Fatalf("the row S01 answers.sha256 %q; want source computed and ref sha256 inputs/answers.tsv S01- Q-", row)
			}
			f.answers, f.calls = answers(c.rows...), nil
			var in *InputError
			_, err := Run("/w", steps(f, out), testWho, noStep)
			if got := err == nil; got != ok || !ok && (!errors.As(err, &in) || len(f.calls) != 0 || !strings.Contains(err.Error(), "an input that changed after a step read it")) {
				t.Errorf("%s, steps done through S%02d: %v, calls %q; want ok %v", c.name, through, err, f.calls, ok)
			}
		}
	}
}

// A done step that reads answers and has no row answers.sha256 is an input
// error: a change of its answers would pass with no check (finding 1 of round
// 1).
func TestADoneStepWithNoAnswersHash(t *testing.T) {
	for _, c := range []struct {
		through int
		step    string
	}{{1, "S01"}, {15, "S10"}} {
		var r work.Record
		for _, row := range doneRows(c.through) {
			if row[0] != c.step || row[1] != "answers.sha256" {
				r = append(r, row)
			}
		}
		f := &fakeSys{record: r}
		f.install(t)
		var in *InputError
		if _, err := Run("/w", steps(f, nil), testWho, noStep); !errors.As(err, &in) || err.Error() != "out/record.tsv: "+c.step+" is done and has no row answers.sha256" || len(f.calls) != 0 {
			t.Errorf("a done %s with no hash: %v, calls %q; want an input error before any step", c.step, err, f.calls)
		}
	}
}

// The marker ID is M- and the first 8 hexadecimal characters of the SHA-256 of
// the file, a tab and the marker (setup.md, The stop table).
func TestMarkerID(t *testing.T) {
	sum := sha256.Sum256([]byte("docs/a.md\t\u2039port\u203a"))
	if got, want := MarkerID("docs/a.md", "\u2039port\u203a"), fmt.Sprintf("M-%x", sum[:4]); got != want || len(got) != 10 {
		t.Errorf("MarkerID = %q, want %q", got, want)
	}
}

// commands.sh holds the commands of the done steps in the order of setup.md,
// each with its comment line; a step done in an earlier run keeps its command
// (D7).
func TestTheCommandsFile(t *testing.T) {
	f := &fakeSys{record: doneRows(3)}
	f.install(t)
	m := steps(f, nil)
	cmd := func(order int, text string) func(work.Record) []Command {
		return func(work.Record) []Command { return []Command{{Order: order, Comment: "the " + text, Text: text}} }
	}
	s := m["S03"]
	s.Commands = cmd(1, "git push root")
	m["S03"] = s
	s = m["S13"]
	s.Commands = func(work.Record) []Command {
		return []Command{{Order: 4, Comment: "the apply", Text: "gh api apply"}, {Order: 2, Comment: "the push of layup-setup", Text: "git push setup"}}
	}
	m["S13"] = s
	s = m["S15"]
	s.Commands = cmd(3, "git push records")
	m["S15"] = s
	if _, err := Run("/w", m, testWho, noStep); err != nil {
		t.Fatal(err)
	}
	want := "# the git push root\ngit push root\n# the push of layup-setup\ngit push setup\n# the git push records\ngit push records\n# the apply\ngh api apply\n"
	if f.commands != want {
		t.Errorf("commands.sh\n%s\nwant\n%s", f.commands, want)
	}
	// A run with no command writes commands.sh again, empty, so no command of
	// an earlier record stays (finding 3 of round 1).
	f = &fakeSys{commands: "echo stale\n"}
	f.install(t)
	if _, err := Run("/w", steps(f, map[string]Outcome{"S02": {Kind: Fail, Evidence: "no commit"}}), testWho, noStep); err != nil || f.commands != "" {
		t.Errorf("a run with no command: %v, commands.sh %q; want it written again, empty", err, f.commands)
	}
}

// A step from S04 to S14 that changed the tree makes one commit on
// layup-setup, by the identity of O-136 at the date pin.time; S03 and S15
// make their own commits (condition 5).
func TestTheCommitOfAStep(t *testing.T) {
	f := &fakeSys{record: doneRows(2)}
	f.install(t)
	if _, err := Run("/w", steps(f, map[string]Outcome{"S03": {Kind: Done, Evidence: "root", Commit: true},
		"S05": {Kind: Done, Evidence: "history", Commit: true}, "S15": {Kind: Done, Evidence: "records", Commit: true}}), testWho, noStep); err != nil {
		t.Fatal(err)
	}
	if want := []string{"chore: setup S05 by LAYUP test <test@layup.invalid> at 2026-10-02T09:30:00Z"}; !reflect.DeepEqual(f.commits, want) {
		t.Errorf("the commits %q, want %q", f.commits, want)
	}
	if Who.Name != "layup-agent[bot]" || Who.Email != "335371832+layup-agent[bot]@users.noreply.github.com" {
		t.Errorf("Who = %+v; want the bot of O-136", Who)
	}
	// A commit that fails makes its step fail, with no done row, so a rerun
	// does the step again.
	f = &fakeSys{record: doneRows(4), commitFails: true}
	f.install(t)
	res, err := Run("/w", steps(f, map[string]Outcome{"S05": {Kind: Done, Evidence: "history", Commit: true}}), testWho, noStep)
	if _, done := f.record.Value("S05", "done"); err != nil || done || len(res.Steps) != 15 || res.Steps[4].Result != "fail" || res.Steps[4].Evidence != "the commit of the step failed" {
		t.Errorf("a failed commit: %v, done %v, the rows %q; want S05 fail and no done row", err, done, res.Steps)
	}
	// A target that is not on the branch layup-setup gets no commit: the step
	// fails, with no done row (finding 2 of round 1).
	for _, head := range []string{"refs/heads/main", "detached"} {
		f = &fakeSys{record: doneRows(4), head: head}
		f.install(t)
		res, err := Run("/w", steps(f, map[string]Outcome{"S05": {Kind: Done, Evidence: "history", Commit: true}}), testWho, noStep)
		if _, done := f.record.Value("S05", "done"); err != nil || done || len(f.commits) != 0 || len(res.Steps) != 15 ||
			res.Steps[4] != (StepRow{"S05", "layup-setup", "fail", "the commit of the step failed: the target is not on the branch layup-setup"}) {
			t.Errorf("a target on %s: %v, done %v, commits %q, the rows %q; want S05 fail, no commit and no done row", head, err, done, f.commits, res.Steps)
		}
	}
}

// The input errors of a work area, and the files that a new work area does
// not have yet (D4, condition 3).
func TestTheInputs(t *testing.T) {
	missing := fmt.Errorf("out/record.tsv: %w", fs.ErrNotExist)
	for _, c := range []struct {
		name string
		f    *fakeSys
		ok   bool
	}{
		{"a work area that is not a directory", &fakeSys{noDir: true}, false},
		{"a record of another form", &fakeSys{recordErr: errors.New("out/record.tsv: line 2")}, false},
		{"answers of another form", &fakeSys{ansErr: errors.New("inputs/answers.tsv: line 1")}, false},
		{"no record and no answers", &fakeSys{recordErr: missing, ansErr: fmt.Errorf("inputs/answers.tsv: %w", fs.ErrNotExist)}, true},
	} {
		c.f.install(t)
		var in *InputError
		if _, err := Run("/w", steps(c.f, nil), testWho, noStep); (err == nil) != c.ok || !c.ok && !errors.As(err, &in) {
			t.Errorf("%s: %v; want ok %v", c.name, err, c.ok)
		}
	}
}

// The steps of this task are stubs: each is not-active, not built yet, so a
// new work area gives S01 not built yet and each later step not run (D9).
func TestTheStubs(t *testing.T) {
	f := &fakeSys{}
	f.install(t)
	res, err := Run("/w", Stubs(), testWho, noStep)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Steps) != 15 || res.Steps[0] != (StepRow{"S01", "layup-setup", "not-active", "not built yet"}) ||
		res.Steps[14] != (StepRow{"S15", "layup-setup", "not-active", "not run: S01 did not pass"}) {
		t.Errorf("the rows %q", res.Steps)
	}
	if s := Stubs(); !reflect.DeepEqual(s["S01"].Reads, []string{"S01-", "Q-"}) || !reflect.DeepEqual(s["S10"].Reads, []string{"M-"}) || !s["S13"].HandOff {
		t.Errorf("the stubs do not read their answers, or S13 is not a hand-off")
	}
}
