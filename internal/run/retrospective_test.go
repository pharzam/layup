package run

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/tsv"
)

// The Retrospective step on a fake records store (task T-v2sv, #359, row 66c
// of the plan of M2e). The fake is an in-memory eventStore: ReadFile answers
// from files, Commit lands its files and records the call; lease and
// commitErr script the fencing cases.

type fakeEventStore struct {
	files     map[string][]byte
	commits   [][]string // the sorted paths of each commit's files
	messages  []string   // the commit messages
	commitErr []error    // given to the next commits, in order
	lease     LeaseRow
	readErr   error
}

func (f *fakeEventStore) ReadFile(_ context.Context, path string) ([]byte, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}
	return f.files[path], nil
}

func (f *fakeEventStore) Commit(_ context.Context, files map[string][]byte, message string) error {
	if len(f.commitErr) > 0 {
		err := f.commitErr[0]
		f.commitErr = f.commitErr[1:]
		if err != nil {
			return err
		}
	}
	paths := make([]string, 0, len(files))
	for path, data := range files {
		f.files[path] = data
		paths = append(paths, path)
	}
	slices.Sort(paths)
	f.commits = append(f.commits, paths)
	f.messages = append(f.messages, message)
	return nil
}

func (f *fakeEventStore) ReadLease(context.Context) (LeaseRow, error) { return f.lease, nil }

func (f *fakeEventStore) WriteLease(_ context.Context, row LeaseRow, _ string) error {
	f.lease = row
	return nil
}

func newFakeStore() *fakeEventStore { return &fakeEventStore{files: map[string][]byte{}} }

var retroAt = time.Date(2026, 10, 16, 11, 30, 0, 0, time.UTC)

func TestTheStepWritesItsEventAndGivesTheNextBet(t *testing.T) {
	f := newFakeStore()
	next, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt)
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" {
		t.Fatalf("the step the loop enters next: got %q, want %q", next, "bet")
	}
	if len(f.commits) != 1 {
		t.Fatalf("the records commits: got %d, want 1", len(f.commits))
	}
	want := "milestone\ttime\nM2e\t2026-10-16T11:30:00Z\n"
	if got := string(f.files["retrospectives.tsv"]); got != want {
		t.Fatalf("retrospectives.tsv:\n%swant:\n%s", got, want)
	}
}

// The proof that no session starts is the signature, not the package's
// imports (internal/run imports session in other files): eventStore holds
// only the lease, ReadFile and Commit, so no call of it can start a session;
// the assertion then shows the one commit touched retrospectives.tsv alone.
func TestTheStepStartsNoSession(t *testing.T) {
	f := newFakeStore()
	if _, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt); err != nil {
		t.Fatal(err)
	}
	for _, paths := range f.commits {
		for _, p := range paths {
			if p != "retrospectives.tsv" {
				t.Errorf("the commit touched %s; only %s may change", p, "retrospectives.tsv")
			}
		}
	}
}

func TestAnExistingTableKeepsItsRowsAndAppends(t *testing.T) {
	first := "M2d\t2026-10-01T00:00:00Z"
	table := "milestone\ttime\n" + first + "\n"
	f := &fakeEventStore{files: map[string][]byte{"retrospectives.tsv": []byte(table)}}
	at := time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)
	next, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", at)
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" {
		t.Fatalf("the step the loop enters next: got %q, want %q", next, "bet")
	}
	want := table + "M2e\t2026-10-16T12:00:00Z\n"
	got := string(f.files["retrospectives.tsv"])
	if got != want {
		t.Fatalf("retrospectives.tsv:\n%swant:\n%s", got, want)
	}
	if line := strings.Split(got, "\n")[1]; line != first {
		t.Fatalf("the first row changed: got %q, want %q", line, first)
	}
}

func TestARetrospectiveDoneBeforeWritesNothing(t *testing.T) {
	table := "milestone\ttime\nM2e\t2026-10-10T00:00:00Z\n"
	f := &fakeEventStore{files: map[string][]byte{"retrospectives.tsv": []byte(table)}}
	next, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt)
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" {
		t.Fatalf("the step the loop enters next: got %q, want %q", next, "bet")
	}
	if len(f.commits) != 0 {
		t.Fatalf("a row of the milestone exists, yet %d commits ran", len(f.commits))
	}
}

func TestAnEmptyOrDashedMilestoneIsRefusedAndNothingCommits(t *testing.T) {
	for _, m := range []string{"", "—"} {
		f := newFakeStore()
		next, err := retrospectiveStep(context.Background(), f, "run-1", m, retroAt)
		if err == nil {
			t.Errorf("milestone %q: no refusal", m)
			continue
		}
		var e *tsv.Error
		if !errors.As(err, &e) || e.Column != "milestone" {
			t.Errorf("milestone %q: the refusal %v names no column milestone", m, err)
		}
		if next != "" {
			t.Errorf("milestone %q: a refused step returned %q", m, next)
		}
		if len(f.commits) != 0 {
			t.Errorf("milestone %q: %d commits ran on a refusal", m, len(f.commits))
		}
	}
}

func TestTheEventTimeIsTheRunsClockInUTC(t *testing.T) {
	// A clock with sub-second parts and a zone that is not UTC: the row holds
	// the run's instant truncated to the second, in UTC.
	at := time.Date(2026, 10, 16, 13, 30, 0, 123456789, time.FixedZone("X", 2*3600))
	f := newFakeStore()
	if _, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", at); err != nil {
		t.Fatal(err)
	}
	row := strings.Split(string(f.files["retrospectives.tsv"]), "\n")[1]
	when, err := time.Parse(timeForm, strings.Split(row, "\t")[1])
	if err != nil {
		t.Fatalf("the row's time %q: %v", row, err)
	}
	if !when.Equal(at.UTC().Truncate(time.Second)) {
		t.Fatalf("the event's time is %v, want %v", when, at.UTC().Truncate(time.Second))
	}
}

func TestATableThatReadRefusesStopsTheStep(t *testing.T) {
	bad := "milestone\ttime\nM2d\tnot-a-time\n"
	f := &fakeEventStore{files: map[string][]byte{"retrospectives.tsv": []byte(bad)}}
	if _, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt); err == nil {
		t.Fatal("a table that read refuses was not refused")
	} else if !strings.Contains(err.Error(), "retrospectives.tsv") {
		t.Fatalf("the error names no file: %v", err)
	}
	if len(f.commits) != 0 {
		t.Fatalf("%d commits ran on a refused table", len(f.commits))
	}
}

func TestAReadErrorIsReturned(t *testing.T) {
	f := &fakeEventStore{files: map[string][]byte{}, readErr: errors.New("the store is down")}
	if _, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt); err == nil {
		t.Fatal("the store's error was not returned")
	}
	if len(f.commits) != 0 {
		t.Fatalf("%d commits ran on a failed read", len(f.commits))
	}
}

func TestTheWriteIsFencedByTheLease(t *testing.T) {
	// The lease still held by this run: a refused push is tried once more.
	f := &fakeEventStore{files: map[string][]byte{}, commitErr: []error{ErrRefused},
		lease: LeaseRow{Run: "run-1", State: "held"}}
	next, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt)
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" || len(f.commits) != 1 {
		t.Fatalf("after one refused push: next %q, commits %d", next, len(f.commits))
	}
	// The lease held by another run: the write is a LostError and lands nothing.
	g := &fakeEventStore{files: map[string][]byte{}, commitErr: []error{ErrRefused},
		lease: LeaseRow{Run: "run-2", State: "held"}}
	next, err = retrospectiveStep(context.Background(), g, "run-1", "M2e", retroAt)
	var lost *LostError
	if !errors.As(err, &lost) {
		t.Fatalf("a push refused with the lease lost: got %v, want a LostError", err)
	}
	if lost.Holder != "run-2" {
		t.Fatalf("the LostError's holder: got %q, want %q", lost.Holder, "run-2")
	}
	if next != "" || len(g.commits) != 0 || g.files["retrospectives.tsv"] != nil {
		t.Fatalf("the event landed on a lost lease: next %q, commits %d", next, len(g.commits))
	}
}

// A non-key column holds — unless its owner checks it, so ReadRetrospectives
// refuses a — in time itself (docs/spec/records.md, the owner checks of the
// task loop's records).
func TestADashInTimeIsRefusedByTheOwner(t *testing.T) {
	table := "milestone\ttime\nM2d\t—\n"
	_, err := ReadRetrospectives([]byte(table))
	if err == nil {
		t.Fatal("a — in the column time was not refused")
	}
	var e *tsv.Error
	if !errors.As(err, &e) || e.Column != "time" {
		t.Fatalf("the refusal %v names no column time", err)
	}
	f := &fakeEventStore{files: map[string][]byte{"retrospectives.tsv": []byte(table)}}
	if _, err := retrospectiveStep(context.Background(), f, "run-1", "M2e", retroAt); err == nil {
		t.Fatal("the step wrote on a table that the owner refuses")
	}
	if len(f.commits) != 0 {
		t.Fatalf("%d commits ran on a refused table", len(f.commits))
	}
}

func TestATableWithADuplicateMilestoneIsRefused(t *testing.T) {
	table := "milestone\ttime\nM2e\t2026-10-10T00:00:00Z\nM2e\t2026-10-11T00:00:00Z\n"
	if _, err := ReadRetrospectives([]byte(table)); err == nil {
		t.Fatal("a second row of a milestone was not refused")
	}
}
