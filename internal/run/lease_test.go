package run

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// clock is a stand-in clock: Sleep moves it, and each turn of it may change
// the stand-in store (another run that beats).
type clock struct {
	t      time.Time
	onTurn func(time.Time)
}

func (c *clock) Now() time.Time { return c.t }
func (c *clock) Sleep(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.t = c.t.Add(d)
	if c.onTurn != nil {
		c.onTurn(c.t)
	}
	return ctx.Err() // a context that ended during the sleep ends it
}

// store is a stand-in records store: the lease row of layup-records, the
// messages of its records commits, and the pushes that it refuses.
type store struct {
	row      LeaseRow
	messages []string
	refuse   int    // refuse this many writes
	takenBy  string // after a refusal, another run holds the lease
}

func (s *store) ReadLease(context.Context) (LeaseRow, error) { return s.row, nil }
func (s *store) WriteLease(_ context.Context, r LeaseRow, message string) error {
	if s.refuse > 0 {
		s.refuse--
		if s.takenBy != "" {
			s.row = LeaseRow{Run: s.takenBy, Host: "h2", Version: "v0.1.0", Heartbeat: 0, State: "held"}
		}
		return ErrRefused
	}
	s.row = r
	s.messages = append(s.messages, message)
	return nil
}

const (
	h     = 5 * time.Minute
	runA  = "aaaaaaaaaaaaaaaa"
	runB  = "bbbbbbbbbbbbbbbb"
	runC  = "cccccccccccccccc"
	ver   = "v0.1.0"
	start = "2026-10-08T12:00:00Z"
)

func at() time.Time { t, _ := time.Parse(time.RFC3339, start); return t }

func held(run string, beat int) LeaseRow {
	return LeaseRow{Run: run, Host: "h1", Version: ver, Started: at().Add(-time.Hour), Heartbeat: beat, State: "held"}
}

func me() LeaseRow { return LeaseRow{Run: runB, Host: "h2", Version: ver} }

func TestAReleasedLeaseIsTakenAtOnce(t *testing.T) {
	s := &store{row: LeaseRow{Run: runA, Host: "h1", Version: ver, State: "released"}}
	c := &clock{t: at()}
	got, err := Take(context.Background(), s, c, me(), h, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	if want := (LeaseRow{Run: runB, Host: "h2", Version: ver, Started: at(), Heartbeat: 0, State: "held"}); got != want || s.row != want {
		t.Errorf("Take = %+v, the store %+v; want %+v", got, s.row, want)
	}
	if c.t != at() {
		t.Errorf("Take waited %v", c.t.Sub(at()))
	}
}

// The demo of row 25a: a second run takes the lease only after the heartbeat
// has not moved for 3 × lease.H by its own clock.
func TestAStoppedHeartbeatIsTakenOverAfterThreeTimesH(t *testing.T) {
	s := &store{row: held(runA, 7)}
	c := &clock{t: at()}
	var lines []string
	got, err := Take(context.Background(), s, c, me(), h, func(l string) { lines = append(lines, l) })
	if err != nil {
		t.Fatal(err)
	}
	if waited := c.t.Sub(at()); waited < 3*h || waited > 3*h+10*time.Second {
		t.Errorf("taken over after %v, want 3 × H (%v), to the next read", waited, 3*h)
	}
	if got.Run != runB || s.row.Run != runB || s.row.State != "held" {
		t.Errorf("the lease is %+v after the takeover", s.row)
	}
	if len(s.messages) != 1 || !strings.Contains(s.messages[0], runA) {
		t.Errorf("the messages %q; want one that names the old run %s", s.messages, runA)
	}
	if len(lines) == 0 || !strings.Contains(lines[0], runA) {
		t.Errorf("progress lines %q; want a line at each read that names the holder", lines)
	}
}

func TestAMovingHeartbeatIsNotTakenOver(t *testing.T) {
	s := &store{row: held(runA, 7)}
	c := &clock{t: at()}
	c.onTurn = func(now time.Time) { // the holder beats every H
		if now.Sub(at())%h == 0 {
			s.row.Heartbeat++
		}
	}
	_, err := Take(context.Background(), s, c, me(), h, func(string) {})
	var heldErr *HeldError
	if !errors.As(err, &heldErr) || heldErr.Holder != runA {
		t.Fatalf("Take with a moving counter = %v; want a HeldError naming %s", err, runA)
	}
	if waited := c.t.Sub(at()); waited < 3*h || waited > 3*h+10*time.Second {
		t.Errorf("gave up after %v, want 3 × H from the start of the wait", waited)
	}
	if s.row.Run != runA || len(s.messages) != 0 {
		t.Errorf("the lease changed: %+v, %q", s.row, s.messages)
	}
}

// Condition 1 of the plan review: the takeover rule comes first, so a counter
// that moved once and then stopped gives exit 1 at 3 × H from the start.
func TestACounterThatMovedOnceThenStoppedGivesHeld(t *testing.T) {
	s := &store{row: held(runA, 7)}
	c := &clock{t: at()}
	c.onTurn = func(now time.Time) {
		if now.Sub(at()) == time.Minute {
			s.row.Heartbeat++
		}
	}
	_, err := Take(context.Background(), s, c, me(), h, func(string) {})
	var heldErr *HeldError
	if !errors.As(err, &heldErr) {
		t.Fatalf("Take = %v; want a HeldError", err)
	}
	if waited := c.t.Sub(at()); waited < 3*h || waited > 3*h+10*time.Second {
		t.Errorf("gave up after %v, want 3 × H from the start of the wait", waited)
	}
}

func TestARefusedTakeoverIsLost(t *testing.T) {
	s := &store{row: held(runA, 7), refuse: 1, takenBy: runC}
	_, err := Take(context.Background(), s, &clock{t: at()}, me(), h, func(string) {})
	var lost *LostError
	if !errors.As(err, &lost) || lost.Holder != runC {
		t.Fatalf("a refused takeover = %v; want a LostError naming %s", err, runC)
	}
}

func TestTheWaitEndsWithTheContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := Take(ctx, &store{row: held(runA, 7)}, &clock{t: at()}, me(), h, func(string) {}); !errors.Is(err, context.Canceled) {
		t.Errorf("Take with a cancelled context = %v", err)
	}
}

func TestBeatAndReleaseGoThroughFencing(t *testing.T) {
	ctx := context.Background()
	mine := held(runB, 3)
	s := &store{row: mine}
	if err := Beat(ctx, s, &mine); err != nil || s.row.Heartbeat != 4 || mine.Heartbeat != 4 {
		t.Fatalf("Beat = %v; the store %d, the run %d; want 4", err, s.row.Heartbeat, mine.Heartbeat)
	}
	s.refuse = 1 // a refusal while this run still holds the lease: one more try
	if err := Beat(ctx, s, &mine); err != nil || s.row.Heartbeat != 5 {
		t.Fatalf("Beat after one refusal = %v, heartbeat %d; want the second try to pass with 5", err, s.row.Heartbeat)
	}
	s.refuse = 2
	var lost *LostError
	if err := Beat(ctx, s, &mine); !errors.As(err, &lost) || lost.Holder != runB {
		t.Fatalf("Beat after two refusals = %v; want a LostError", err)
	}
	s.refuse, s.takenBy = 1, runC // after a takeover the old run's beat is refused
	if err := Beat(ctx, s, &mine); !errors.As(err, &lost) || lost.Holder != runC {
		t.Fatalf("Beat after a takeover = %v; want a LostError naming %s", err, runC)
	}
	mine, s = held(runB, 9), &store{row: held(runB, 9)}
	if err := Release(ctx, s, &mine); err != nil || s.row.State != "released" || s.row.Heartbeat != 9 {
		t.Fatalf("Release = %v, the store %+v", err, s.row)
	}
}

func TestHeartbeatBeatsAtEachH(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mine := held(runB, 0)
	s := &store{row: mine}
	c := &clock{t: at()}
	c.onTurn = func(now time.Time) {
		if now.Sub(at()) > 3*h {
			cancel()
		}
	}
	if err := Heartbeat(ctx, s, c, &mine, h); err != nil {
		t.Fatalf("Heartbeat ended with %v; want nil when its context ends", err)
	}
	if s.row.Heartbeat != 3 || len(s.messages) != 3 {
		t.Errorf("heartbeat %d with %d commits after 3 × H; want 3, a commit each", s.row.Heartbeat, len(s.messages))
	}
	s.refuse, s.takenBy = 1, runC
	if err := Heartbeat(context.Background(), s, c, &mine, h); err == nil {
		t.Error("Heartbeat after a takeover: nil; want a LostError")
	}
}
