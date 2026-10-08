// Package run is layup run (docs/spec/run.md). In row 25a of the plan (task
// T-trej, #130, O-173) it holds the rules of a run: the lease and fencing, a
// human decision and copy before read, over two small interfaces, so row 25b
// builds the steps of Start and the restart on them with internal/git and the
// forge.
package run

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// LeaseRow is the one row of lease.tsv (docs/spec/records.md, the block lease).
type LeaseRow struct {
	Run, Host, Version string
	Started            time.Time
	Heartbeat          int
	State              string // held or released
}

// Store is the records branch as the lease sees it: ReadLease reads the row of
// lease.tsv as the forge holds it; WriteLease makes a records commit with row
// and message and pushes it, and gives ErrRefused when git refuses the push.
type Store interface {
	ReadLease(ctx context.Context) (LeaseRow, error)
	WriteLease(ctx context.Context, row LeaseRow, message string) error
}

// Clock is the run's own clock; Sleep ends early with the context.
type Clock interface {
	Now() time.Time
	Sleep(ctx context.Context, d time.Duration) error
}

// ErrRefused is a records push that git refused.
var ErrRefused = errors.New("the records push was refused")

// HeldError: the heartbeat of another run moved for 3 × lease.H of the wait.
type HeldError struct{ Holder string }

func (e *HeldError) Error() string { return "the lease is held by the run " + e.Holder }

// LostError: this run does not hold the lease, or a push stays refused.
type LostError struct{ Holder string }

func (e *LostError) Error() string {
	return "the records push was refused; the lease names the run " + e.Holder
}

// readEvery is the period of the reads of a run that waits for the lease.
const readEvery = 10 * time.Second

// Take takes the lease for me (its Run, Host and Version): at once when it is
// released; else it reads the row every ten seconds and, at each read, first
// takes it over when heartbeat has not changed for 3 × h since the last change
// it saw, by its own clock; else, 3 × h after the start of its wait, it gives a
// HeldError. progress gets one line at each read (docs/spec/run.md, The lease
// and fencing).
func Take(ctx context.Context, s Store, c Clock, me LeaseRow, h time.Duration, progress func(string)) (LeaseRow, error) {
	row, err := s.ReadLease(ctx)
	if err != nil {
		return LeaseRow{}, err
	}
	begin := c.Now()
	seen, changed := row, begin
	for row.State != "released" {
		if err := c.Sleep(ctx, readEvery); err != nil {
			return LeaseRow{}, err
		}
		if row, err = s.ReadLease(ctx); err != nil {
			return LeaseRow{}, err
		}
		now := c.Now()
		if row.Run != seen.Run || row.Heartbeat != seen.Heartbeat {
			seen, changed = row, now
		}
		switch {
		case row.State == "released":
		case now.Sub(changed) >= 3*h:
			return claim(ctx, s, c, me, fmt.Sprintf("lease: the run %s takes the lease over from the run %s", me.Run, row.Run))
		case now.Sub(begin) >= 3*h:
			return LeaseRow{}, &HeldError{Holder: row.Run}
		default:
			progress(fmt.Sprintf("lease: held by the run %s; its heartbeat has not changed for %d s",
				row.Run, int(now.Sub(changed).Seconds())))
		}
	}
	return claim(ctx, s, c, me, fmt.Sprintf("lease: the run %s takes the released lease from the run %s", me.Run, row.Run))
}

// claim writes me as the holder. The run does not hold the lease yet, so a
// refused push is lost: another run took it first.
func claim(ctx context.Context, s Store, c Clock, me LeaseRow, message string) (LeaseRow, error) {
	me.Started, me.Heartbeat, me.State = c.Now().UTC().Truncate(time.Second), 0, "held"
	if err := s.WriteLease(ctx, me, message); errors.Is(err, ErrRefused) {
		return LeaseRow{}, lost(ctx, s)
	} else if err != nil {
		return LeaseRow{}, err
	}
	return me, nil
}

// lost reads the lease again and names its holder.
func lost(ctx context.Context, s Store) error {
	row, err := s.ReadLease(ctx)
	if err != nil {
		return err
	}
	return &LostError{Holder: row.Run}
}

// Fenced runs a records write of the run that holds the lease mine. When git
// refuses its push, it reads the lease again: held by another run, or no
// longer held, is a LostError; still held by mine, the write is tried once
// more, and a second refusal is a LostError (docs/spec/run.md, Fencing).
func Fenced(ctx context.Context, s Store, mine string, write func(context.Context) error) error {
	err := write(ctx)
	if !errors.Is(err, ErrRefused) {
		return err
	}
	row, err := s.ReadLease(ctx)
	if err != nil {
		return err
	}
	if row.Run != mine || row.State != "held" {
		return &LostError{Holder: row.Run}
	}
	if err := write(ctx); errors.Is(err, ErrRefused) {
		return lost(ctx, s)
	} else if err != nil {
		return err
	}
	return nil
}

// Beat adds 1 to the heartbeat of mine, through Fenced.
func Beat(ctx context.Context, s Store, mine *LeaseRow) error {
	return Fenced(ctx, s, mine.Run, func(ctx context.Context) error {
		next := *mine
		next.Heartbeat++
		if err := s.WriteLease(ctx, next, fmt.Sprintf("lease: heartbeat %d", next.Heartbeat)); err != nil {
			return err
		}
		*mine = next
		return nil
	})
}

// Release gives the lease back, through Fenced.
func Release(ctx context.Context, s Store, mine *LeaseRow) error {
	return Fenced(ctx, s, mine.Run, func(ctx context.Context) error {
		next := *mine
		next.State = "released"
		if err := s.WriteLease(ctx, next, "lease: released"); err != nil {
			return err
		}
		*mine = next
		return nil
	})
}

// Heartbeat beats every h while the run holds the lease, until ctx ends (nil)
// or a beat is lost (its error).
func Heartbeat(ctx context.Context, s Store, c Clock, mine *LeaseRow, h time.Duration) error {
	for {
		if c.Sleep(ctx, h) != nil {
			return nil
		}
		if err := Beat(ctx, s, mine); err != nil {
			return err
		}
	}
}
