package run

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"
)

// The release at the end of a restart (docs/spec/run.md, The lease and fencing; task T-yngf, #479), on the stand-in store of lease_test.go.

// endStore is that stand-in with a count of the writes that were tried and a fault: on runs first at each write, and its error is the write's.
type endStore struct {
	mu    sync.Mutex
	store // the row, the messages that landed, the refusals
	tried int
	land  bool // a heartbeat lands although its context has ended: a beat in flight that completes
	on    func(ctx context.Context, message string) error
}

func (s *endStore) ReadLease(ctx context.Context) (LeaseRow, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.ReadLease(ctx)
}

func (s *endStore) WriteLease(ctx context.Context, r LeaseRow, message string) error {
	s.mu.Lock()
	s.tried++
	on := s.on
	s.mu.Unlock()
	if on != nil {
		if err := on(ctx, message); err != nil {
			return err
		}
	}
	if err := ctx.Err(); err != nil && !(s.land && strings.HasPrefix(message, "lease: heartbeat")) {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store.WriteLease(ctx, r, message)
}

func (s *endStore) seen() (row LeaseRow, log []string, tried int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.row, append([]string(nil), s.messages...), s.tried
}

// tick is a clock whose sleep takes a millisecond, so a heartbeat does not spin.
type tick struct{}

func (tick) Now() time.Time { return at() }
func (tick) Sleep(ctx context.Context, _ time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(time.Millisecond):
		return nil
	}
}

func endState(s Store, mine LeaseRow) *state {
	r := &state{cfg: Config{RunID: runB, Clock: tick{}, LeaseH: 5}, viaLease: s, lease: mine}
	r.cancel = func() {}
	return r
}

func lastRow(result, detail string) []Step { return []Step{{"meaning", result, detail}} }

func TestAFailedRowKeepsItsDetailAndTheLeaseIsReleased(t *testing.T) {
	s := &endStore{store: store{row: held(runB, 3)}}
	got := endState(s, held(runB, 3)).endLease(context.Background(), lastRow("fail", "own"))
	row, log, _ := s.seen()
	if got[0] != (Step{"meaning", "fail", "own"}) || row.State != "released" || len(log) != 1 || log[0] != "lease: released" {
		t.Errorf("the row %+v, the lease %+v, the commits %q; want the row as it was, the lease released by one commit", got[0], row, log)
	}
}

func TestAFailedReleaseIsNamedAfterTheFailure(t *testing.T) {
	for _, c := range []struct{ name, result, detail, want string }{
		{"a failed row", "fail", "own", "own; the lease: disk full"},
		{"a done row", "done", "all", "the lease: disk full"}, // a control on the base: the rule of a row done is as it was
	} {
		t.Run(c.name, func(t *testing.T) {
			s := &endStore{store: store{row: held(runB, 3)}, on: func(context.Context, string) error { return errors.New("disk full") }}
			got := endState(s, held(runB, 3)).endLease(context.Background(), lastRow(c.result, c.detail))
			if got[0].Result != "fail" || got[0].Detail != c.want {
				t.Errorf("the row %+v; want fail with %q", got[0], c.want)
			}
		})
	}
}

// Another run took the lease: the release is refused and published nowhere, and the loss is named beside the failure.
func TestARunThatLostTheLeaseNamesTheLossAndPublishesNoRelease(t *testing.T) {
	s := &endStore{store: store{row: held(runB, 3), refuse: 1, takenBy: runC}}
	r := endState(s, held(runB, 3))
	got := r.endLease(context.Background(), lastRow("fail", "own"))
	row, log, _ := s.seen()
	if want := "own; the lease: the records push was refused; the lease names the run " + runC; got[0].Detail != want {
		t.Errorf("the detail %q; want %q", got[0].Detail, want)
	}
	if row.Run != runC || row.State != "held" || len(log) != 0 || r.lostLease == nil {
		t.Errorf("the lease %+v, the commits %q, the loss %v; want the other run's row, no commit, the loss remembered", row, log, r.lostLease)
	}
}

// Two refusals while the row still names this run are a LostError too (Fenced): the loss is terminal, so no third push is tried although it would succeed.
func TestTwoRefusalsAreTerminalAndNoReleaseIsTriedAfterThem(t *testing.T) {
	s := &endStore{store: store{row: held(runB, 3), refuse: 2}}
	r := endState(s, held(runB, 3))
	got := r.endLease(context.Background(), lastRow("fail", "own"))
	if want := "own; the lease: the records push was refused; the lease names the run " + runB; got[0].Detail != want {
		t.Errorf("the detail %q; want %q", got[0].Detail, want)
	}
	again := r.endLease(context.Background(), lastRow("fail", "own"))
	row, log, tried := s.seen()
	if tried != 2 || row.State != "held" || len(log) != 0 || again[0].Detail != "own" {
		t.Errorf("after the second end: %d writes tried, the lease %+v, the commits %q, the detail %q; want 2, held, none, own", tried, row, log, again[0].Detail)
	}
}

// Controls, green on the base too: a run that holds nothing, or whose lease is released, tries no write; a row done stays done.
func TestARunThatHoldsNothingReleasesNothing(t *testing.T) {
	for name, mine := range map[string]LeaseRow{"no lease": {}, "another run's": held(runC, 1), "released": {Run: runB, State: "released"}} {
		s := &endStore{store: store{row: held(runC, 1)}}
		got := endState(s, mine).endLease(context.Background(), lastRow("done", "all"))
		if _, _, tried := s.seen(); tried != 0 || got[0] != (Step{"meaning", "done", "all"}) {
			t.Errorf("%s: %d writes tried, the row %+v; want none, the row as it was", name, tried, got[0])
		}
	}
}

// phaseStep's own release failed (the step context ended): the lease is still held, one more try at the end, with no second line for the same error.
func TestTheEndTriesAgainAfterAReleaseOfThePhaseThatFailed(t *testing.T) {
	failed := false
	s := &endStore{store: store{row: held(runB, 3)}, on: func(context.Context, string) error {
		if !failed {
			failed = true
			return errors.New("disk full")
		}
		return nil
	}}
	r := endState(s, held(runB, 3))
	phase := r.leaseStep(context.Background())
	got := r.endLease(context.Background(), []Step{{"phase", "fail", phase.Detail}})
	row, _, tried := s.seen()
	if got[0].Detail != "disk full" || row.State != "released" || tried != 2 {
		t.Errorf("the row %+v, the lease %+v, %d writes; want disk full, released, 2", got[0], row, tried)
	}
	if again := r.endLease(context.Background(), got); func() int { _, _, n := s.seen(); return n }() != 2 || again[0].Detail != "disk full" {
		t.Error("a lease that is released was released again")
	}
}

// The heartbeat is joined before r.lease is read: with a beat in flight the release waits for it, is the last write, and no later beat restores held.
func TestTheHeartbeatEndsBeforeTheRelease(t *testing.T) {
	inFlight, proceed := make(chan struct{}), make(chan struct{})
	var once, unblock sync.Once
	s := &endStore{store: store{row: held(runB, 0)}, land: true, on: func(_ context.Context, m string) error {
		if strings.HasPrefix(m, "lease: heartbeat") {
			once.Do(func() { close(inFlight); <-proceed })
		}
		return nil
	}}
	t.Cleanup(func() { unblock.Do(func() { close(proceed) }) })
	r := endState(s, held(runB, 0))
	r.beat()
	<-inFlight
	finished := make(chan []Step, 1)
	go func() { finished <- r.endLease(context.Background(), lastRow("fail", "own")) }()
	select {
	case <-finished:
		t.Fatal("the end returned while a beat was in flight; it must join the heartbeat first")
	case <-time.After(50 * time.Millisecond):
	}
	unblock.Do(func() { close(proceed) })
	<-finished
	time.Sleep(20 * time.Millisecond)
	row, log, _ := s.seen()
	if row.State != "released" || len(log) < 2 || log[0] != "lease: heartbeat 1" || log[len(log)-1] != "lease: released" || r.lease.State != "released" {
		t.Errorf("the lease %+v, the last commits %q; want released, and the release the last write", row, log)
	}
}

// A heartbeat error that is not a loss proves nothing about the lease: the release is still tried. A row done is fail with the error.
func TestAHeartbeatErrorThatIsNotALossStillReleases(t *testing.T) {
	for _, c := range []struct{ result, detail, want string }{{"fail", "own", "own"}, {"done", "all", "the lease: net down"}} {
		s := &endStore{store: store{row: held(runB, 0)}, on: func(_ context.Context, m string) error {
			if strings.HasPrefix(m, "lease: heartbeat") {
				return errors.New("net down")
			}
			return nil
		}}
		r := endState(s, held(runB, 0))
		r.beat()
		for i := 0; i < 2000 && r.lost() == nil; i++ {
			time.Sleep(time.Millisecond)
		}
		got := r.endLease(context.Background(), lastRow(c.result, c.detail))
		row, _, _ := s.seen()
		if got[0].Detail != c.want || got[0].Result != "fail" && c.result == "done" || row.State != "released" || r.lostLease != nil {
			t.Errorf("%s row: %+v, the lease %+v, the loss %v; want %q, released, no loss", c.result, got[0], row, r.lostLease, c.want)
		}
	}
}

// The release does not use the context of the steps, which a lost beat or an interrupt has ended, and is bounded by its own timeout.
func TestTheReleaseOutlivesTheContextOfTheStepsAndIsBounded(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	s := &endStore{store: store{row: held(runB, 3)}}
	if endState(s, held(runB, 3)).endLease(cancelled, lastRow("fail", "own")); func() string { r, _, _ := s.seen(); return r.State }() != "released" {
		t.Error("a cancelled context of the steps stopped the release")
	}
	s = &endStore{store: store{row: held(runB, 3)}, on: func(ctx context.Context, _ string) error { <-ctx.Done(); return ctx.Err() }}
	r := endState(s, held(runB, 3))
	r.endTimeout = 30 * time.Millisecond
	begin := time.Now()
	got := r.endLease(cancelled, lastRow("fail", "own"))
	if got[0].Detail != "own; the lease: context deadline exceeded" || time.Since(begin) > 5*time.Second {
		t.Errorf("the row %+v after %v; want the deadline named, soon", got[0], time.Since(begin))
	}
}
