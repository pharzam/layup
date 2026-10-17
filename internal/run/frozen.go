package run

// This file holds the frozen test list of a task (docs/spec/loop.md, The frozen
// test list; docs/spec/records.md, the block task-tests; task T-6zj3, #361, row
// 66e of the plan of M2e): the Go schema of tasks/<task>/tests.tsv, the freezing
// at the first valid handoff to the verifier, the finding "test changed", the
// amendment that freezes the list again, and the frozen files that an attempt
// starts with. The text form of a plan's list of tests, the file that holds a
// test and the run that gives a test's results are #387's: the caller passes
// them as typed values.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// TestsSchema is the Go value of the block task-tests of docs/spec/records.md;
// its owner is this package, and its test compares it with the block.
var TestsSchema = tsv.Schema{Name: "task-tests", Location: "records:tasks/<task>/tests.tsv", Columns: []tsv.Column{
	{Name: "test", Type: "text", Key: true},
	{Name: "source", Type: "path"},
	{Name: "sha256", Type: "sha256"},
	{Name: "harness", Type: "id(<word>)"},
	{Name: "frozen", Type: "time"},
}}

// testsPath is the path in the records of the frozen list of a task.
func testsPath(task string) string { return "tasks/" + task + "/tests.tsv" }

// ReadTests reads a frozen list by TestsSchema and the rules of its block: at
// least one test, never — in source, sha256, harness or frozen, and one time on
// every row. An error names its line and column.
func ReadTests(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, TestsSchema)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, &tsv.Error{Line: 1, Reason: "a frozen list has at least one test"}
	}
	for i, r := range rows {
		for c := 1; c < len(TestsSchema.Columns); c++ {
			if r[c] == "" {
				return nil, &tsv.Error{Line: i + 2, Column: TestsSchema.Columns[c].Name, Reason: "this column never holds the empty value"}
			}
		}
		if r[4] != rows[0][4] {
			return nil, &tsv.Error{Line: i + 2, Column: "frozen", Reason: "the time is the same on each row of a freezing"}
		}
	}
	return rows, nil
}

// WriteTests writes a frozen list by TestsSchema. It writes nothing for a list
// that ReadTests would refuse.
func WriteTests(rows [][]string) ([]byte, error) {
	data, err := table(TestsSchema, rows)
	if err != nil {
		return nil, err
	}
	if _, err := ReadTests(data); err != nil {
		return nil, err
	}
	return data, nil
}

// A TestResult is what the caller knows of one listed test: its ID in the plan,
// the file that holds it, and whether it fails at the base and passes at the
// head. The producer of the results (a gate or test run) is not this package's.
type TestResult struct {
	ID, Source                string
	FailsAtBase, PassesAtHead bool
}

// A Freeze is the list of tests of a handoff to the verifier, with their results.
type Freeze struct{ Tests []TestResult }

// A FrozenFile is a frozen test file: its path in the target and its bytes.
type FrozenFile struct {
	Path string
	Data []byte
}

// redGreen is the item tests-red-green computed from typed results (loop.md,
// The frozen test list): a list of at least one test, each with an ID and a
// source, none repeated, each failing at the base and passing at the head. An
// empty list proves nothing and is not red-green (guardrails §2).
func redGreen(tests []TestResult) bool {
	if len(tests) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, t := range tests {
		if t.ID == "" || t.Source == "" || seen[t.ID] || !t.FailsAtBase || !t.PassesAtHead {
			return false
		}
		seen[t.ID] = true
	}
	return true
}

// freezeRows gives the rows of the list that a freezing writes: one for each
// test in the order given, with the hash of its source and the one time. A test
// that was frozen before at the same source with the same hash keeps its harness;
// any other takes the harness of the freezing session.
func freezeRows(prev [][]string, tests []TestResult, hashes map[string]string, harness string, at time.Time) [][]string {
	rows := make([][]string, 0, len(tests))
	for _, t := range tests {
		h, sum := harness, hashes[t.Source]
		for _, p := range prev {
			if p[0] == t.ID && p[1] == t.Source && p[2] == sum {
				h = p[3]
			}
		}
		rows = append(rows, []string{t.ID, t.Source, sum, h, stamp(at)})
	}
	return rows
}

// sameAsFrozen reports whether the tests are the frozen list: the same IDs at
// the same sources, each source at its frozen hash, and no other test.
func sameAsFrozen(rows [][]string, tests []TestResult, hashes map[string]string) bool {
	if len(rows) != len(tests) {
		return false
	}
	for _, r := range rows {
		if !slices.ContainsFunc(tests, func(t TestResult) bool { return t.ID == r[0] && t.Source == r[1] && hashes[t.Source] == r[2] }) {
			return false
		}
	}
	return true
}

// changedTests gives the finding "test changed" of each frozen test whose
// source, read by hashAt ("" for a file that the head lacks), has another hash.
func changedTests(rows [][]string, hashAt func(path string) (string, error)) ([]Finding, error) {
	var out []Finding
	for _, r := range rows {
		sum, err := hashAt(r[1])
		if err != nil {
			return nil, err
		}
		if sum != r[2] {
			out = append(out, Finding{Rule: "test changed", File: r[1], Text: r[0]})
		}
	}
	return out, nil
}

// validHandoff reports whether the detail of an event handoff is valid, and
// whether it names a comment.
func validHandoff(detail string) (valid, comment bool) {
	return detail == "valid" || strings.HasPrefix(detail, "valid "), strings.HasPrefix(detail, "valid comment ")
}

// amendmentPending reports whether the events hold an amendment that no freezing
// has used (loop.md, The frozen test list): after the last valid handoff of a
// developer's session, a valid handoff with a comment from a session of the plan
// role, then one from a session of the plan-review role whose session started
// after that handoff. The order is the events', never the clock's; a verifier's
// record shares the plan reviewer's role, so only a session that started after
// the new plan is its review.
func amendmentPending(events, sessions [][]string) bool {
	role := map[string]string{}
	for _, s := range sessions {
		role[s[sessionID]] = s[sessionRole]
	}
	started := func(id string) int {
		return slices.IndexFunc(events, func(e []string) bool { return e[1] == "session" && e[3] == id })
	}
	last := -1
	for i, e := range events {
		if v, _ := validHandoff(e[6]); e[1] == "handoff" && v && role[e[3]] == stepRole("implement") {
			last = i
		}
	}
	for i := last + 1; i < len(events); i++ {
		e := events[i]
		if _, c := validHandoff(e[6]); e[1] != "handoff" || !c || role[e[3]] != stepRole("plan-of-a-task") {
			continue
		}
		for j := i + 1; j < len(events); j++ {
			f := events[j]
			if _, c := validHandoff(f[6]); f[1] == "handoff" && c && role[f[3]] == stepRole("plan-review") && started(f[3]) > i {
				return true
			}
		}
	}
	return false
}

// shouldFreeze reports whether a handoff writes the list: a valid handoff to
// the verifier that was given one.
func shouldFreeze(in HandoffInput, v handoff.Verdict) bool {
	return in.Freeze != nil && !in.CodeWritten && in.To == "verifier" && v.Valid
}

// ErrNoBoundHead is the refusal to freeze a session that has no bound head.
var ErrNoBoundHead = errors.New("the session has no bound head")

// boundHead gives the SHA that the event bound of the session holds, "" for none.
func boundHead(events [][]string, session string) string {
	for _, e := range events {
		if e[1] == "bound" && e[3] == session {
			return e[5]
		}
	}
	return ""
}

// blobHash gives the SHA-256 of the file at path in the commit head of the
// run's clone, "" when the head has no such file.
func (s *Sessions) blobHash(head, path string) (string, error) {
	data, err := s.blobAt(head, path)
	if err != nil || data == nil {
		return "", err
	}
	return digest(data), nil
}

// blobAt gives the bytes of the file at path in head, nil when it is none.
func (s *Sessions) blobAt(head, path string) ([]byte, error) {
	entries, err := git.LsTree(s.Clone, head, path)
	if err != nil {
		return nil, err
	}
	if len(entries) != 1 || entries[0].Path != path || entries[0].Type != "blob" {
		return nil, nil
	}
	return git.Show(s.Clone, head, path)
}

// readTests reads the frozen list of a task, nil when it has none.
func (s *Sessions) readTests(ctx context.Context, task string) ([][]string, error) {
	data, err := s.Store.ReadFile(ctx, testsPath(task))
	if err != nil || data == nil {
		return nil, err
	}
	rows, err := ReadTests(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", testsPath(task), err)
	}
	return rows, nil
}

// freezing is what a handoff to the verifier that was given a list decides:
// the item tests-red-green, and the rows to write (nil for none).
type freezing struct {
	redGreen bool
	rows     [][]string
}

// freeze decides the list of a handoff (loop.md, The frozen test list): the
// hash of each source is read at the head the session was bound to; with no
// list the results must be red-green and the list is the first; with a list and
// no amendment pending the tests must be the frozen ones at their frozen hashes
// (a changed source is not counted as passed); with an amendment pending the
// results alone decide and the list is written again whole, when it differs.
func (s *Sessions) freeze(ctx context.Context, in HandoffInput) (freezing, error) {
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return freezing{}, err
	}
	head := boundHead(events, in.Session)
	if head == "" {
		return freezing{}, fmt.Errorf("%w: %s", ErrNoBoundHead, in.Session)
	}
	start, err := s.startRowOf(ctx, in.Session)
	if err != nil {
		return freezing{}, err
	}
	sessions, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return freezing{}, err
	}
	cur, err := s.readTests(ctx, in.Task)
	if err != nil {
		return freezing{}, err
	}
	hashes := map[string]string{}
	for _, t := range in.Freeze.Tests {
		if t.Source == "" {
			continue
		}
		if hashes[t.Source], err = s.blobHash(head, t.Source); err != nil {
			return freezing{}, err
		}
	}
	ok := redGreen(in.Freeze.Tests)
	for _, t := range in.Freeze.Tests {
		ok = ok && hashes[t.Source] != ""
	}
	if !ok {
		return freezing{}, nil
	}
	switch {
	case cur == nil:
		return freezing{true, freezeRows(nil, in.Freeze.Tests, hashes, start[sessionHarness], s.Now())}, nil
	case amendmentPending(events, sessions):
		rows := freezeRows(cur, in.Freeze.Tests, hashes, start[sessionHarness], s.Now())
		if sameList(cur, rows) {
			return freezing{true, nil}, nil
		}
		return freezing{true, rows}, nil
	}
	return freezing{redGreen: sameAsFrozen(cur, in.Freeze.Tests, hashes)}, nil
}

// sameList reports whether two lists hold the same tests, sources, hashes and
// harnesses, whatever their time.
func sameList(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !slices.Equal(a[i][:4], b[i][:4]) {
			return false
		}
	}
	return true
}

// FrozenFindings gives the findings of the frozen list of the task at head: the
// finding "test changed" of each listed test whose source at head has another
// hash than the frozen one (a source the head lacks included). No list, no finding.
func (s *Sessions) FrozenFindings(ctx context.Context, task, head string) ([]Finding, error) {
	rows, err := s.readTests(ctx, task)
	if err != nil || rows == nil {
		return nil, err
	}
	return changedTests(rows, func(path string) (string, error) { return s.blobHash(head, path) })
}

// FrozenFiles gives the files of the frozen list, each the bytes whose SHA-256
// is the frozen one, found among the heads that were bound to the task's
// sessions, the newest first; a test whose bytes no bound head holds is an
// error, never a guess. A task with no list has no file.
func (s *Sessions) FrozenFiles(ctx context.Context, task string) ([]FrozenFile, error) {
	rows, err := s.readTests(ctx, task)
	if err != nil || rows == nil {
		return nil, err
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	var heads []string
	for _, e := range events {
		if e[1] == "bound" {
			heads = append([]string{e[5]}, heads...)
		}
	}
	var out []FrozenFile
	for _, r := range rows {
		var found []byte
		for _, h := range heads {
			data, err := s.blobAt(h, r[1])
			if err != nil {
				return nil, err
			}
			if data != nil && digest(data) == r[2] {
				found = data
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("no bound head of %s holds %s at the frozen hash %s", task, r[1], r[2])
		}
		if !slices.ContainsFunc(out, func(f FrozenFile) bool { return f.Path == r[1] }) {
			out = append(out, FrozenFile{Path: r[1], Data: found})
		}
	}
	return out, nil
}
