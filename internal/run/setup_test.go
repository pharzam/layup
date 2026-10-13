package run

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
)

// The loop of the step setup (docs/spec/scaffold.md, The stops of the runner;
// the acceptance test "The stops of the runner under layup run (inventory)";
// row 54 of the plan, task T-eahr, #277): the runner is run again until it does
// not stop, each stop answered from the records, by a prose session, or not at
// all. The runner, the answers, the session and the records are fakes here.

func mRow(id string) setup.StopRow {
	return setup.StopRow{Step: "S10", Question: id, Ask: "What is the value?", Where: "NOTES.md:2 <x>"}
}
func fRow(path string) setup.StopRow {
	return setup.StopRow{Step: "S07", Question: "F-" + path, Ask: "Write " + path, Where: path}
}
func stopOf(rows ...setup.StopRow) setup.Result { return setup.Result{Stops: rows} }
func verifyStop() setup.Result {
	return stopOf(setup.StopRow{Step: "S15", Question: "O-verify", Ask: "Run layup setup verify", Where: ""})
}

// failedAt is the table of a run whose step failed: the steps before it are
// done, the step is fail with its evidence, and each later step is not-active.
func failedAt(step, evidence string) setup.Result {
	var rows []setup.StepRow
	state := setup.Done
	for i := 1; i <= 15; i++ {
		id := fmt.Sprintf("S%02d", i)
		switch {
		case id == step:
			rows, state = append(rows, setup.StepRow{Step: id, Actor: "layup-setup", Result: setup.Fail, Evidence: evidence}), setup.NotActive
		case state == setup.NotActive:
			rows = append(rows, setup.StepRow{Step: id, Actor: "layup-setup", Result: state, Evidence: "not run: " + step + " did not pass"})
		default:
			rows = append(rows, setup.StepRow{Step: id, Actor: "layup-setup", Result: state, Evidence: "fine"})
		}
	}
	return setup.Result{Steps: rows}
}

// loopFake is the world of a loop: the results of the runner in order, and what
// the loop asked of the answers, the session and the records.
type loopFake struct {
	results  []any // a setup.Result or an error
	ran      int
	answered [][]string
	prosed   [][]setup.StopRow
	proseErr map[int]error // by the number of the call of the session, from 0
	removed  []string
	removeIn error
}

func (f *loopFake) loop(markers map[string][]string) setupLoop {
	return setupLoop{
		max: 6,
		run: func() (setup.Result, error) {
			if f.ran >= len(f.results) {
				return setup.Result{}, errors.New("the fake runner has no result left")
			}
			r := f.results[f.ran]
			f.ran++
			if err, ok := r.(error); ok {
				return setup.Result{}, err
			}
			return r.(setup.Result), nil
		},
		markers: func(id string) ([]string, bool) { row, ok := markers[id]; return row, ok },
		answer: func(rows [][]string) error {
			f.answered = append(f.answered, rows...)
			return nil
		},
		prose: func(_ context.Context, rows []setup.StopRow) (proseDone, error) {
			n := len(f.prosed)
			f.prosed = append(f.prosed, rows)
			if err := f.proseErr[n]; err != nil {
				return proseDone{}, err
			}
			return proseDone{Session: fmt.Sprintf("S-0000000%d", n+1), Files: map[string][]byte{"a.md": []byte("a")}}, nil
		},
		remove: func(_ context.Context, p proseDone, reason string) error {
			f.removed = append(f.removed, p.Session+": "+reason)
			return f.removeIn
		},
	}
}

var asked = map[string][]string{
	"M-aaaaaaaa": {"M-aaaaaaaa", "Carol", "operator", "https://x/1#issuecomment-1", ""},
	"M-bbbbbbbb": {"M-bbbbbbbb", "gap", "operator", "https://x/1#issuecomment-2", "the value of <b>"},
	"M-cccccccc": {"M-cccccccc", "9", "operator", "https://x/1#issuecomment-3", ""},
}

const wantSetupDone = "S01 to S14 done, in the work area"

func TestARunThatDoesNotStopGoesOn(t *testing.T) {
	f := &loopFake{results: []any{setup.Result{Steps: failedAt("", "").Steps}}}
	if got := f.loop(asked).do(context.Background()); got != (Step{"setup", "done", wantSetupDone}) {
		t.Errorf("a run with every step done: %+v; want setup done", got)
	}
}

// S13's done row is a hand-off, whose result is operator: it is done for the step.
func TestAHandOffIsDone(t *testing.T) {
	r := failedAt("", "")
	r.Steps[12].Result = setup.Operator
	f := &loopFake{results: []any{r}}
	if got := f.loop(asked).do(context.Background()); got != (Step{"setup", "done", wantSetupDone}) {
		t.Errorf("a run whose S13 is handed off: %+v; want setup done", got)
	}
}

// The one row O-verify of S15 is the end: the runner has done S01 to S14.
func TestTheStopOfOVerifyIsTheEnd(t *testing.T) {
	f := &loopFake{results: []any{verifyStop()}}
	if got := f.loop(asked).do(context.Background()); got != (Step{"setup", "done", wantSetupDone}) || f.ran != 1 || len(f.prosed) != 0 || len(f.answered) != 0 {
		t.Errorf("a stop of O-verify: %+v after %d runs; want done, one run, no session and no answer", got, f.ran)
	}
}

// A stop with M- rows adds the row of each marker that the batch asked, in the
// order of the stop, and the runner is run again.
func TestAStopWithMRowsAddsTheRowsOfTheMarkersTheBatchAsked(t *testing.T) {
	f := &loopFake{results: []any{stopOf(mRow("M-cccccccc"), mRow("M-aaaaaaaa")), verifyStop()}}
	got := f.loop(asked).do(context.Background())
	if want := [][]string{asked["M-cccccccc"], asked["M-aaaaaaaa"]}; got.Result != "done" || !reflect.DeepEqual(f.answered, want) || f.ran != 2 {
		t.Errorf("the step %+v, the rows added %q after %d runs; want done, the rows of M-cccccccc and M-aaaaaaaa in the order of the stop, and two runs", got, f.answered, f.ran)
	}
}

// A marker that the batch asked and that S05's deletion or the prose step took
// out of the tree is in no stop, so it gets no row and fails nothing; one that
// no question of the batch names is the fail.
func TestAMarkerNoStopListsGetsNoRow(t *testing.T) {
	f := &loopFake{results: []any{stopOf(mRow("M-bbbbbbbb")), verifyStop()}}
	if got := f.loop(asked).do(context.Background()); got.Result != "done" || len(f.answered) != 1 || f.answered[0][0] != "M-bbbbbbbb" {
		t.Errorf("the step %+v, the rows added %q; want done and the row of the one marker that the stop lists", got, f.answered)
	}
}

func TestAnMRowTheBatchDidNotAskIsFail(t *testing.T) {
	f := &loopFake{results: []any{stopOf(mRow("M-aaaaaaaa"), mRow("M-zzzzzzzz"), mRow("M-yyyyyyyy")), verifyStop()}}
	got := f.loop(asked).do(context.Background())
	if got.Result != "fail" || !strings.Contains(got.Detail, "M-zzzzzzzz, M-yyyyyyyy") || !strings.Contains(got.Detail, "did not ask") || len(f.answered) != 0 || f.ran != 1 {
		t.Errorf("the step %+v, the rows added %q after %d runs; want fail, naming M-zzzzzzzz and M-yyyyyyyy, and no row added", got, f.answered, f.ran)
	}
}

// A stop with F- rows starts one prose session for the rows of the stop.
func TestAStopWithFRowsStartsAProseSession(t *testing.T) {
	f := &loopFake{results: []any{stopOf(fRow("docs/a.md"), fRow("README.md")), verifyStop()}}
	got := f.loop(asked).do(context.Background())
	if want := [][]setup.StopRow{{fRow("docs/a.md"), fRow("README.md")}}; got.Result != "done" || !reflect.DeepEqual(f.prosed, want) || f.ran != 2 || len(f.removed) != 0 {
		t.Errorf("the step %+v, the sessions for %q after %d runs, removed %q; want done after one session for the two rows", got, f.prosed, f.ran, f.removed)
	}
}

// Any other row is a fail that names the rows, whatever else the table holds:
// S01-stack with no answer, an M- row of no marker is not here (it is the case
// above), any other O- row, and O-verify of another step or with another row.
func TestAnyOtherRowIsFail(t *testing.T) {
	for _, c := range []struct {
		name string
		stop setup.Result
		want string
	}{
		{"S01-stack with no answer", stopOf(setup.StopRow{Step: "S01", Question: "S01-stack", Ask: "Which stack?"}), "S01-stack"},
		{"another O- row", stopOf(setup.StopRow{Step: "S15", Question: "O-other", Ask: "x"}), "O-other"},
		{"O-verify of another step", stopOf(setup.StopRow{Step: "S13", Question: "O-verify", Ask: "x"}), "O-verify"},
		{"O-verify and another row", stopOf(setup.StopRow{Step: "S15", Question: "O-verify", Ask: "x"}, mRow("M-aaaaaaaa")), "O-verify"},
		{"O-verify and an F- row", stopOf(fRow("README.md"), setup.StopRow{Step: "S15", Question: "O-verify", Ask: "x"}), "O-verify"},
		{"an S01- row that is not the stack", stopOf(setup.StopRow{Step: "S01", Question: "S01-name", Ask: "x"}), "S01-name"},
	} {
		f := &loopFake{results: []any{c.stop}}
		got := f.loop(asked).do(context.Background())
		if got.Result != "fail" || !strings.Contains(got.Detail, c.want) || len(f.answered) != 0 || len(f.prosed) != 0 || f.ran != 1 {
			t.Errorf("%s: %+v, answered %q, sessions %d; want fail naming %s, with no answer and no session", c.name, got, f.answered, len(f.prosed), c.want)
		}
	}
}

// An input error of the runner, a fail and a not-active step make the step fail.
func TestAnInputErrorAFailAndANotActiveAreFail(t *testing.T) {
	f := &loopFake{results: []any{&setup.InputError{Err: errors.New("inputs/answers.tsv: an input that changed after a step read it")}}}
	if got := f.loop(asked).do(context.Background()); got.Result != "fail" || !strings.Contains(got.Detail, "an input that changed after a step read it") {
		t.Errorf("an input error: %+v; want fail with its text", got)
	}
	f = &loopFake{results: []any{failedAt("S12", "the stack cobol has no entry in the catalog")}}
	if got := f.loop(asked).do(context.Background()); got != (Step{"setup", "fail", "S12: the stack cobol has no entry in the catalog"}) {
		t.Errorf("a failed step: %+v; want fail naming the step and its evidence, not the not-active steps after it", got)
	}
	notActive := setup.Result{Steps: []setup.StepRow{{Step: "S01", Result: setup.Done, Evidence: "x"}, {Step: "S02", Result: setup.NotActive, Evidence: "its check could not run"}}}
	f = &loopFake{results: []any{notActive}}
	if got := f.loop(asked).do(context.Background()); got != (Step{"setup", "fail", "S02: its check could not run"}) {
		t.Errorf("a not-active step: %+v; want fail", got)
	}
}

// A fail after a prose session of the run removes that session's files from the
// records, with the session and the reason; one before it removes none; and a
// session whose rows the runner passed stays.
func TestAFailAfterAProseSessionRemovesItsFilesAndOneBeforeRemovesNone(t *testing.T) {
	for _, c := range []struct {
		name    string
		results []any
		prose   map[int]error
		want    []string // what was removed
		detail  string
	}{
		{"a check that fails right after the session", []any{stopOf(fRow("a.md")), failedAt("S09", "guardrails: fail: no rule")}, nil,
			[]string{"S-00000001: S09: guardrails: fail: no rule"}, "S09: guardrails: fail: no rule"},
		{"after the answers of the markers, which do not end the session's turn", []any{stopOf(fRow("a.md")), stopOf(mRow("M-aaaaaaaa")), failedAt("S11", "markers: fail")}, nil,
			[]string{"S-00000001: S11: markers: fail"}, "S11: markers: fail"},
		{"a marker that the prose text brought and the batch did not ask", []any{stopOf(fRow("a.md")), stopOf(mRow("M-newnewne"))}, nil,
			[]string{"S-00000001: the stop lists M-newnewne, which the batch did not ask"}, "M-newnewne"},
		{"two sessions, the second's files are removed and the first's stay", []any{stopOf(fRow("a.md")), stopOf(fRow("b.md")), failedAt("S08", "glossary: fail")}, nil,
			[]string{"S-00000002: S08: glossary: fail"}, "S08: glossary: fail"},
		{"a second session that ends badly leaves the first, whose rows passed", []any{stopOf(fRow("a.md")), stopOf(fRow("b.md"))}, map[int]error{1: errors.New("the session S-bad: no-result: missing")},
			nil, "the session S-bad: no-result: missing"},
		{"a fail before any session", []any{stopOf(mRow("M-aaaaaaaa")), failedAt("S12", "the stack x has no entry")}, nil, nil, "S12: the stack x has no entry"},
		{"an input error after a session", []any{stopOf(fRow("a.md")), &setup.InputError{Err: errors.New("a bad row")}}, nil,
			[]string{"S-00000001: a bad row"}, "a bad row"},
	} {
		f := &loopFake{results: c.results, proseErr: c.prose}
		got := f.loop(asked).do(context.Background())
		if got.Result != "fail" || !strings.Contains(got.Detail, c.detail) || !slices.Equal(f.removed, c.want) {
			t.Errorf("%s: %+v, removed %q; want fail with %q, removed %q", c.name, got, f.removed, c.detail, c.want)
		}
	}
	// A removal that fails is said, and the step is fail all the same.
	f := &loopFake{results: []any{stopOf(fRow("a.md")), failedAt("S09", "x")}, removeIn: errors.New("the push was refused")}
	if got := f.loop(asked).do(context.Background()); got.Result != "fail" || !strings.Contains(got.Detail, "S09: x") || !strings.Contains(got.Detail, "the push was refused") {
		t.Errorf("a removal that fails: %+v; want fail with the reason and the error of the removal", got)
	}
}

// The same stop twice in a row is a fail, and removes the files of the session
// that did not settle it; the loop has no other bound, so it ends after max runs.
func TestAStopThatDoesNotChangeIsFailAndTheLoopHasABound(t *testing.T) {
	f := &loopFake{results: []any{stopOf(fRow("a.md")), stopOf(fRow("a.md"))}}
	got := f.loop(asked).do(context.Background())
	if got.Result != "fail" || !strings.Contains(got.Detail, "the same rows again") || len(f.prosed) != 1 || len(f.removed) != 1 {
		t.Errorf("the same stop twice: %+v, %d sessions, removed %q; want fail after one session, and its files removed", got, len(f.prosed), f.removed)
	}
	var results []any
	for i := 0; i < 20; i++ {
		results = append(results, stopOf(fRow(fmt.Sprintf("f%d.md", i%2))), stopOf(mRow(fmt.Sprintf("M-%08d", i))))
	}
	many := map[string][]string{}
	for i := 0; i < 20; i++ {
		many[fmt.Sprintf("M-%08d", i)] = []string{fmt.Sprintf("M-%08d", i), "v", "operator", "u", ""}
	}
	f = &loopFake{results: results}
	got = f.loop(many).do(context.Background())
	if got.Result != "fail" || !strings.Contains(got.Detail, "did not end") || f.ran > 6 {
		t.Errorf("a runner that never ends: %+v after %d runs; want fail within the bound of 6", got, f.ran)
	}
}

// A session that is refused or ends badly is the fail of the step, with its own
// text, which names the session; no session before it means none removed.
func TestAProseSessionThatEndsBadlyIsFail(t *testing.T) {
	f := &loopFake{results: []any{stopOf(fRow("a.md"))}, proseErr: map[int]error{0: &PhaseError{Session: "S-1a2b3c4d", Reason: "no-result: the output file a.md is missing"}}}
	got := f.loop(asked).do(context.Background())
	if got != (Step{"setup", "fail", "the session S-1a2b3c4d: no-result: the output file a.md is missing"}) || len(f.removed) != 0 {
		t.Errorf("a refused output: %+v, removed %q; want fail with the session's text and nothing removed", got, f.removed)
	}
}

// A cancelled context ends the loop between two runs, as a lost lease does.
func TestACancelledContextEndsTheLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := &loopFake{results: []any{stopOf(mRow("M-aaaaaaaa")), verifyStop()}}
	l := f.loop(asked)
	l.answer = func(rows [][]string) error { cancel(); return nil }
	if got := l.do(ctx); got.Result != "fail" || !strings.Contains(got.Detail, "context canceled") || f.ran != 1 {
		t.Errorf("a cancelled context: %+v after %d runs; want fail before the second run", got, f.ran)
	}
}
