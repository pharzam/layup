package cli

import (
	"bytes"
	"testing"
	"time"
)

// A stoppedClock gives each step a ticker that ticks only when the test sends
// on it, so no test waits ten seconds.
type stoppedClock struct {
	ticks   []chan time.Time
	stopped int
}

func (c *stoppedClock) ticker(time.Duration) (<-chan time.Time, func()) {
	tick := make(chan time.Time)
	c.ticks = append(c.ticks, tick)
	return tick, func() { c.stopped++ }
}

// tick sends one tick on the ticker of step i, and fails when step i started
// no ticker.
func tick(t *testing.T, clock *stoppedClock, i int) {
	t.Helper()
	if len(clock.ticks) <= i {
		t.Fatalf("%d tickers started, want a ticker for step %d", len(clock.ticks), i+1)
	}
	clock.ticks[i] <- time.Time{}
}

func TestProgressPrintsALinePerStepAndABeatEveryTenSeconds(t *testing.T) {
	var out bytes.Buffer
	clock := &stoppedClock{}
	p := &progress{w: &out, command: "gate", every: 10 * time.Second, ticker: clock.ticker}

	p.step(1, 2, "go-build")
	tick(t, clock, 0)
	tick(t, clock, 0)
	p.step(2, 2, "go-test")
	tick(t, clock, 1)
	p.end()

	want := "layup gate: [1/2] go-build\n" +
		"layup gate: [1/2] go-build: 10 s\n" +
		"layup gate: [1/2] go-build: 20 s\n" +
		"layup gate: [2/2] go-test\n" +
		"layup gate: [2/2] go-test: 10 s\n"
	if out.String() != want {
		t.Fatalf("the lines:\n%s\nwant:\n%s", out.String(), want)
	}
	if clock.stopped != 2 {
		t.Fatalf("%d tickers stopped, want 2", clock.stopped)
	}
}

func TestProgressEndWithNoStepPrintsNothing(t *testing.T) {
	var out bytes.Buffer
	p := &progress{w: &out, command: "gate", every: 10 * time.Second, ticker: (&stoppedClock{}).ticker}
	p.end()
	p.end()
	if out.Len() != 0 {
		t.Fatalf("the lines %q, want none", out.String())
	}
}

func TestNewProgressBeatsEveryTenSeconds(t *testing.T) {
	if p := newProgress(&bytes.Buffer{}, "setup"); p.every != 10*time.Second || p.command != "setup" {
		t.Fatalf("every %v, command %q; want 10s and setup", p.every, p.command)
	}
}

// A line of a wait (docs/spec/run.md, The command) is printed with the step,
// and the beat starts again, so a wait shows one line every ten seconds.
func TestANoteOfAWaitStartsTheBeatAgain(t *testing.T) {
	var out bytes.Buffer
	clock := &stoppedClock{}
	p := &progress{w: &out, command: "run", every: 10 * time.Second, ticker: clock.ticker}
	p.step(4, 9, "root-push")
	tick(t, clock, 0)
	p.note("waiting for the push of the root commit")
	tick(t, clock, 1)
	p.end()
	p.note("after the end: printed with no step")
	want := "layup run: [4/9] root-push\n" +
		"layup run: [4/9] root-push: 10 s\n" +
		"layup run: [4/9] root-push: waiting for the push of the root commit\n" +
		"layup run: [4/9] root-push: 20 s\n" +
		"layup run: after the end: printed with no step\n"
	if out.String() != want {
		t.Fatalf("the lines:\n%s\nwant:\n%s", out.String(), want)
	}
	if clock.stopped != 2 {
		t.Fatalf("%d tickers stopped, want 2", clock.stopped)
	}
}
