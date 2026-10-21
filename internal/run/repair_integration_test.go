//go:build integration

package run

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// The repairs of a task session at a restart (task T-syxm, row 66v), against the
// httptest forge through the adapter and real git.

// pushFault stops the run at the push of a session's head: land is whether the
// forge took the push before the stop.
type pushFault struct {
	SessionStore
	land bool
}

func (p *pushFault) PushHead(ctx context.Context, sha, branch string) error {
	if p.land {
		if err := p.SessionStore.PushHead(ctx, sha, branch); err != nil {
			return err
		}
	}
	return errors.New("the run stopped at the push")
}

// attemptFault fails the first records commit that adds an event attempt.
type attemptFault struct {
	SessionStore
	failed bool
}

func (a *attemptFault) Commit(ctx context.Context, files map[string][]byte, message string) error {
	if data := files["tasks/T-ab12/events.tsv"]; !a.failed && strings.Contains(string(data[strings.LastIndex(strings.TrimSpace(string(data)), "\n"):]), "\tattempt\t") {
		a.failed = true
		return errors.New("the run stopped before the attempt")
	}
	return a.SessionStore.Commit(ctx, files, message)
}

// readFault fails the read of one file of the records.
type readFault struct {
	SessionStore
	path string
}

func (r *readFault) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if path == r.path {
		return nil, errors.New("the read of " + path + " failed")
	}
	return r.SessionStore.ReadFile(ctx, path)
}

// thief gives the lease to another run once stolen is set.
type thief struct {
	SessionStore
	stolen bool
}

func (th *thief) ReadLease(ctx context.Context) (LeaseRow, error) {
	if th.stolen {
		return LeaseRow{Run: "fedcba9876543210", State: "held"}, nil
	}
	return th.SessionStore.ReadLease(ctx)
}

// forgeFault fails the read of the comments of an issue, or runs a hook just before it.
type forgeFault struct {
	forge.Forge
	err  error
	hook func()
}

func (f *forgeFault) Comments(ctx context.Context, issue int) ([]forge.Comment, error) {
	if f.hook != nil {
		f.hook()
	}
	if f.err != nil {
		return nil, f.err
	}
	return f.Forge.Comments(ctx, issue)
}

// stopped runs a task session whose push the run's stop interrupts, and gives its ID.
func stopped(t *testing.T, r *Sessions, base worldAt, task string, land bool) string {
	t.Helper()
	s := afterRestart(r)
	s.Store = &pushFault{SessionStore: r.Store, land: land}
	spec := harness(t, base, work(helloSum))
	spec.Task = task
	id, err := s.TaskSession(context.Background(), spec)
	if err == nil || !strings.Contains(err.Error(), "stopped at the push") {
		t.Fatalf("TaskSession: %q, %v; want the stop at the push", id, err)
	}
	return id
}

// kindsOf gives the kinds of the events of task in the records branch.
func kindsOf(t *testing.T, w *world, task string) []string {
	t.Helper()
	rows, err := records.ReadEvents([]byte(keptRecord(t, w, "tasks/"+task+"/events.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	var out []string
	for _, e := range rows {
		out = append(out, e[1])
	}
	return out
}

// summaries gives the bodies of the comments of issue that are the comments of a task session of T-ab12.
func summaries(w *world, issue int) []string {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	var out []string
	for _, c := range w.fake.comments[issue] {
		if b, _ := c["body"].(string); strings.Contains(b, " session of T-ab12, attempt ") {
			out = append(out, b)
		}
	}
	return out
}

func line(id, tail string) string {
	return id + ": developer session of T-ab12, attempt 1: " + tail + "\n"
}

func TestARestartBindsAPushThatTheForgeAccepted(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	id := stopped(t, r, base, "T-ab12", true)
	got, err := afterRestart(r).RepairTask(context.Background(), "T-ab12")
	if want := []string{"attempt", "session", "result", "push", "bound"}; err != nil || !slices.Equal(kindsOf(t, w, "T-ab12"), want) || !slices.Equal(got.Bound, []string{id}) {
		t.Fatalf("events %v, Repair %+v, %v; want %v bound", kindsOf(t, w, "T-ab12"), got, err, want)
	}
	if s := summaries(w, taskIssue); !slices.Equal(s, []string{line(id, "done")}) {
		t.Errorf("the comments of the task's issue %q; want the one line, not posted again", s)
	}
}

func TestARestartRefusesALostPushAndStartsOneAttempt(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	id := stopped(t, r, base, "T-ab12", false)
	got, err := afterRestart(r).RepairTask(context.Background(), "T-ab12")
	if want := []string{"attempt", "session", "result", "push", "refused", "attempt"}; err != nil || !slices.Equal(kindsOf(t, w, "T-ab12"), want) || got.Attempt != 2 {
		t.Fatalf("events %v, Repair %+v, %v; want %v", kindsOf(t, w, "T-ab12"), got, err, want)
	}
	if s, want := summaries(w, taskIssue), []string{line(id, "done"), line(id, "done, refused push-lost")}; !slices.Equal(s, want) {
		t.Errorf("the comments of the task's issue %q; want %q (the refusal reaches the comment)", s, want)
	}
}

// Shape runs before any bet: the register holds a task of kind specification and the records no bets.tsv.
func TestARestartRepairsAShapePushBeforeAnyBet(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	ctx := context.Background()
	files, err := registerFiles(mustRead(t, r, "tasks.tsv"), nil, Task{ID: "T-sp01", Issue: 4, Kind: "specification", Requirements: []string{"REQ-005"}, Size: "small", Opened: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	delete(files, "transitions.tsv")
	if err := r.commit(ctx, files, "tasks: register T-sp01"); err != nil {
		t.Fatal(err)
	}
	if _, err := r.StartAttempt(ctx, "T-sp01", base.base); err != nil {
		t.Fatal(err)
	}
	stopped(t, r, base, "T-sp01", false)
	if data, err := r.Store.ReadFile(ctx, "bets.tsv"); err != nil || data != nil {
		t.Fatalf("bets.tsv: %q, %v; want none", data, err)
	}
	if _, err := afterRestart(r).RepairTask(ctx, "T-sp01"); err != nil {
		t.Fatal(err)
	}
	if want := []string{"attempt", "session", "result", "push", "refused", "attempt"}; !slices.Equal(kindsOf(t, w, "T-sp01"), want) {
		t.Errorf("events %v; want %v", kindsOf(t, w, "T-sp01"), want)
	}
}

func mustRead(t *testing.T, r *Sessions, path string) []byte {
	t.Helper()
	data, err := r.Store.ReadFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// A stop between the commit of refused and the commit of attempt: the next runs start one attempt, and no more.
func TestARestartBetweenPushLostAndTheAttemptStartsExactlyOne(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	ctx := context.Background()
	stopped(t, r, base, "T-ab12", false)
	first := afterRestart(r)
	first.Store = &attemptFault{SessionStore: r.Store}
	if _, err := first.RepairTask(ctx, "T-ab12"); err == nil || !strings.Contains(err.Error(), "before the attempt") {
		t.Fatalf("the stop: %v; want its error", err)
	}
	if want := []string{"attempt", "session", "result", "push", "refused"}; !slices.Equal(kindsOf(t, w, "T-ab12"), want) {
		t.Fatalf("events after the stop %v; want %v", kindsOf(t, w, "T-ab12"), want)
	}
	if _, err := afterRestart(r).RepairTask(ctx, "T-ab12"); err != nil {
		t.Fatal(err)
	}
	head := gitOut(t, w.bare, "rev-parse", "layup-records")
	if _, err := afterRestart(r).RepairTask(ctx, "T-ab12"); err != nil {
		t.Fatal(err)
	}
	if want := []string{"attempt", "session", "result", "push", "refused", "attempt"}; !slices.Equal(kindsOf(t, w, "T-ab12"), want) || gitOut(t, w.bare, "rev-parse", "layup-records") != head {
		t.Errorf("events %v; want %v and no commit by the third run", kindsOf(t, w, "T-ab12"), want)
	}
}

// A commit pushed onto the branch after the push: the head is a descendant of the SHA, not the SHA.
func TestAMovedBranchHeadIsPushLost(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	stopped(t, r, base, "T-ab12", true)
	pushedSHA := strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "refs/heads/task/T-ab12/1"))
	next := strings.TrimSpace(gitOut(t, w.bare, "commit-tree", pushedSHA+"^{tree}", "-p", pushedSHA, "-m", "moved"))
	gitOut(t, w.bare, "update-ref", "refs/heads/task/T-ab12/1", next)
	if _, err := afterRestart(r).RepairTask(context.Background(), "T-ab12"); err != nil {
		t.Fatal(err)
	}
	if want := []string{"attempt", "session", "result", "push", "refused", "attempt"}; !slices.Equal(kindsOf(t, w, "T-ab12"), want) {
		t.Errorf("events %v; want %v", kindsOf(t, w, "T-ab12"), want)
	}
}

func TestATaskSessionsCommentIsOnItsIssue(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	id := commentSession(t, r, base, "developer", "a note\n")
	if got, want := summaries(w, taskIssue), []string{line(id, "done")}; !slices.Equal(got, want) {
		t.Errorf("the comments of the task's issue %q; want %q", got, want)
	}
	if got := summaries(w, 2); len(got) != 0 {
		t.Errorf("the control issue holds %q; want no comment of a task session", got)
	}
}

func TestASessionCommentThatTheAppPostedIsNotPostedAgain(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	commentSession(t, r, base, "developer", "a note\n")
	got, err := afterRestart(r).RepairTask(context.Background(), "T-ab12")
	if err != nil || len(got.Posted) != 0 || len(summaries(w, taskIssue)) != 1 {
		t.Errorf("Repair %+v, %v, %d comments; want none posted", got, err, len(summaries(w, taskIssue)))
	}
}

func TestAMissingSessionCommentIsPostedAgainFromTheRecords(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	w.fake.mu.Lock()
	w.fake.failNote = " session of T-ab12, attempt 1: done"
	w.fake.mu.Unlock()
	spec := harness(t, base, work(helloSum))
	id, err := r.TaskSession(context.Background(), spec)
	if err == nil || len(summaries(w, taskIssue)) != 0 {
		t.Fatalf("a post that the forge refuses: %v, %d comments; want an error and none", err, len(summaries(w, taskIssue)))
	}
	for call := 1; call <= 2; call++ {
		got, err := afterRestart(r).RepairTask(context.Background(), "T-ab12")
		if want := []string{line(id, "done")}; err != nil || !slices.Equal(summaries(w, taskIssue), want) || (call == 1) != slices.Equal(got.Posted, []string{id}) {
			t.Errorf("call %d: Repair %+v, %v, comments %q; want %q posted by the first call only", call, got, err, summaries(w, taskIssue), want)
		}
	}
}

// With no usable row of the task in the register, no comment is posted, on any issue.
func TestNoRegisterRowOrABrokenRegisterPostsNowhere(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct{ name, want string }{{"no row", "has no row of T-ab12"}, {"refused", "tasks.tsv:"}, {"absent", "no tasks.tsv"}} {
		t.Run(tt.name, func(t *testing.T) {
			w, r, base := commentWorld(t)
			ctx := context.Background()
			switch tt.name {
			case "no row":
				files, err := registerFiles(nil, nil, Task{ID: "T-zz99", Issue: 9, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Opened: time.Now().UTC().Truncate(time.Second)})
				if err != nil {
					t.Fatal(err)
				}
				if err := r.commit(ctx, map[string][]byte{"tasks.tsv": files["tasks.tsv"]}, "tasks: another"); err != nil {
					t.Fatal(err)
				}
			case "refused":
				if err := r.commit(ctx, map[string][]byte{"tasks.tsv": []byte("broken\n")}, "tasks: broken"); err != nil {
					t.Fatal(err)
				}
			default:
				if err := r.Store.(*recordsStore).Remove(ctx, []string{"tasks.tsv"}, "tasks: none"); err != nil {
					t.Fatal(err)
				}
			}
			id, err := r.TaskSession(ctx, harness(t, base, work(helloSum)))
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("TaskSession: %q, %v; want an error with %q", id, err, tt.want)
			}
			if _, err := afterRestart(r).RepairTask(ctx, "T-ab12"); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Errorf("RepairTask: %v; want an error with %q", err, tt.want)
			}
			for _, issue := range []int{2, taskIssue} {
				if s := summaries(w, issue); len(s) != 0 {
					t.Errorf("issue %d holds %q; want no comment of the session", issue, s)
				}
			}
		})
	}
}

// Each guard of the post holds on both paths: the session's own and the hook's repost. The comment
// is missing in the second path (the forge refused it once), then the guard is in place.
func TestTheGuardsOfThePostHoldOnBothPaths(t *testing.T) {
	t.Parallel()
	guards := []struct {
		name string
		arm  func(s *Sessions) // puts the guard in s
		want func(err error) bool
	}{
		{"a failed read of the comments", func(s *Sessions) { s.Forge = &forgeFault{Forge: s.Forge, err: errors.New("comments: Bad Gateway")} },
			func(err error) bool { return err != nil && strings.Contains(err.Error(), "Bad Gateway") }},
		{"a failed read of the register", func(s *Sessions) { s.Store = &readFault{SessionStore: s.Store, path: "tasks.tsv"} },
			func(err error) bool { return err != nil && strings.Contains(err.Error(), "read of tasks.tsv failed") }},
		{"a lease lost before the post", func(s *Sessions) {
			th := &thief{SessionStore: s.Store}
			s.Store = th
			s.Forge = &forgeFault{Forge: s.Forge, hook: func() { th.stolen = true }}
		}, func(err error) bool {
			var lost *LostError
			return errors.As(err, &lost) && lost.Holder == "fedcba9876543210"
		}},
	}
	for _, g := range guards {
		t.Run(g.name+", the session's path", func(t *testing.T) {
			w, r, base := commentWorld(t)
			s := afterRestart(r)
			g.arm(s)
			id, err := s.TaskSession(context.Background(), harness(t, base, work(helloSum)))
			if !g.want(err) || len(summaries(w, taskIssue)) != 0 {
				t.Errorf("TaskSession %q: %v, %d comments; want the guard's error and no post", id, err, len(summaries(w, taskIssue)))
			}
		})
		t.Run(g.name+", the repost", func(t *testing.T) {
			w, r, base := commentWorld(t)
			w.fake.mu.Lock()
			w.fake.failNote = " session of T-ab12, attempt 1: done"
			w.fake.mu.Unlock()
			if _, err := r.TaskSession(context.Background(), harness(t, base, work(helloSum))); err == nil {
				t.Fatal("the first post: no error")
			}
			s := afterRestart(r)
			g.arm(s)
			if _, err := s.RepairTask(context.Background(), "T-ab12"); !g.want(err) || len(summaries(w, taskIssue)) != 0 {
				t.Errorf("RepairTask: %v, %d comments; want the guard's error and no post", err, len(summaries(w, taskIssue)))
			}
		})
	}
}

// Only the App's own comment of the same bytes counts as posted: the same text by a person, or by another App, does not suppress it.
func TestTheSameTextFromAPersonOrAnotherAppDoesNotSuppressTheSummary(t *testing.T) {
	t.Parallel()
	w, r, base := commentWorld(t)
	w.fake.mu.Lock()
	w.fake.failNote = " session of T-ab12, attempt 1: done"
	w.fake.mu.Unlock()
	id, err := r.TaskSession(context.Background(), harness(t, base, work(helloSum)))
	if err == nil {
		t.Fatal("the first post: no error")
	}
	text := line(id, "done")
	w.fake.mu.Lock()
	w.fake.comments[taskIssue] = append(w.fake.comments[taskIssue],
		map[string]any{"id": 801, "user": map[string]any{"id": 101, "login": "pharzam"}, "created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": text},
		map[string]any{"id": 802, "user": map[string]any{"id": 777, "login": "other-app[bot]"}, "performed_via_github_app": map[string]any{"slug": "other-app"},
			"created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": text})
	w.fake.mu.Unlock()
	for call := 1; call <= 2; call++ {
		if _, err := afterRestart(r).RepairTask(context.Background(), "T-ab12"); err != nil {
			t.Fatal(err)
		}
		if n := len(summaries(w, taskIssue)); n != 3 {
			t.Errorf("after call %d: %d comments of this line; want the two foreign ones and ours, once (session %s)", call, n, id)
		}
	}
}
