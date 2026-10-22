//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// A restart that holds the lease releases it whatever its last row is (docs/spec/run.md, The lease and fencing; task T-yngf, #479), on the restart world.

const nextRun = "fedcba9876543210"

// readClock counts the waits of the lease (Take reads every readEvery) and nothing else.
type readClock struct {
	Clock
	reads *atomic.Int32
}

func (c readClock) Sleep(ctx context.Context, d time.Duration) error {
	if d == readEvery {
		c.reads.Add(1)
	}
	return c.Clock.Sleep(ctx, d)
}

func recordsLog(t *testing.T, w *world) string {
	return gitOut(t, w.bare, "log", "--format=%s", "refs/heads/layup-records")
}

func leaseText(t *testing.T, w *world) string {
	return show(t, w.bare, "refs/heads/layup-records", "lease.tsv")
}

// nextTakesAtOnce runs a restart of another run and checks that it took the released lease with no wait of its own.
func nextTakesAtOnce(t *testing.T, w *world, from, as string) []Step {
	t.Helper()
	cfg := w.config("0.1.0-dev")
	cfg.RunID = as
	var reads atomic.Int32
	cfg.Clock = readClock{cfg.Clock, &reads}
	rows := Restart(context.Background(), cfg)
	if l := rowOf(rows, "lease"); l.Result != "done" || reads.Load() != 0 {
		t.Errorf("the next restart: lease %+v after %d reads of the lease; want done with none", l, reads.Load())
	}
	if !strings.Contains(recordsLog(t, w), "takes the released lease from the run "+from) {
		t.Errorf("no commit takes the released lease:\n%s", recordsLog(t, w))
	}
	return rows
}

// midRestart runs the steps of a restart up to the lease and then tail, and gives the rows before the release at its end.
func midRestart(t *testing.T, w *world, tail ...func(*state, context.Context) Step) (*state, context.Context, []Step) {
	t.Helper()
	cfg := w.config("0.1.0-dev")
	cfg.RunID = nextRun
	r := &state{cfg: cfg}
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	t.Cleanup(func() { r.endBeat(); cancel() })
	chain := []func(context.Context) Step{func(ctx context.Context) Step { return r.forgeStep(ctx, false) }, r.cloneStep, r.versionStep, r.takeStep}
	for _, step := range tail {
		chain = append(chain, func(ctx context.Context) Step { return step(r, ctx) })
	}
	return r, ctx, r.steps(ctx, make([]string, len(chain)), 0, chain)
}

func TestAFailedRestartReleasesTheLeaseAndTheNextTakesItAtOnce(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name, step string
		startFails bool
		arm        func(*world, *Config)
	}{
		{"a step of Intake", "meaning", false, func(_ *world, cfg *Config) { cfg.Routing = nil }},
		{"a step before phase", "probe", true, func(w *world, _ *Config) { w.fake.defBranch = "gone" }},
		{"a step inside phase", "phase", true, func(w *world, _ *Config) { w.fake.failOpen = 2 }},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := newWorld(t)
			if c.startFails {
				w.fake.failOpen = 2
				w.start(w.config("0.1.0-dev"))
			} else {
				allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
			}
			cfg := w.config("0.1.0-dev")
			c.arm(w, &cfg)
			rows := Restart(context.Background(), cfg)
			if last := rows[len(rows)-1]; last.Name != c.step || last.Result != "fail" || strings.Contains(last.Detail, "the lease:") {
				t.Fatalf("the last row %+v; want %s fail with its own detail", last, c.step)
			}
			if !strings.Contains(leaseText(t, w), "\treleased\n") || gitOut(t, w.bare, "log", "-1", "--format=%s", "refs/heads/layup-records") != "lease: released\n" {
				t.Errorf("the lease after the fail is not released by the last commit:\n%s%s", leaseText(t, w), recordsLog(t, w))
			}
			w.fake.defBranch = ""
			nextTakesAtOnce(t, w, "0123456789abcdef", nextRun)
		})
	}
}

// A control, green on the base: a restart whose phase releases the lease itself releases it once.
func TestARestartThatRunsPhaseReleasesOnce(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	w.fake.failOpen = 2
	w.start(w.config("0.1.0-dev"))
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), "forge", "clone", "version", "lease", "probe", "phase")
	if n := strings.Count(recordsLog(t, w), "lease: released\n"); n != 1 {
		t.Errorf("%d commits release the lease; want 1", n)
	}
}

// Another run took the lease during a step: the release is refused, and the other run's row and records are untouched.
func TestARestartThatLostTheLeaseToAnotherRunPublishesNoRelease(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	r, ctx, rows := midRestart(t, w, func(*state, context.Context) Step {
		pushRecords(t, w.bare, map[string]string{"lease.tsv": "run\thost\tversion\tstarted\theartbeat\tstate\n" +
			"cccccccccccccccc\thost-2\t0.1.0-dev\t2026-10-08T13:00:00Z\t0\theld\n"}, "lease: the run cccccccccccccccc takes the lease over")
		return failf("meaning", "own detail")
	})
	head := recordsHead(t, w)
	rows = r.endLease(ctx, rows)
	last := rows[len(rows)-1]
	if !strings.Contains(last.Detail, "the lease names the run cccccccccccccccc") || last.Result != "fail" || recordsHead(t, w) != head || !strings.Contains(leaseText(t, w), "cccccccccccccccc") {
		t.Errorf("the last row %+v, the lease:\n%s\nwant the loss named, the head %s unchanged", last, leaseText(t, w), head)
	}
}

// Two refusals of the release while the row names this run are a loss too, and no third push is tried, though it would succeed (the hook lets the third pass).
func TestTwoRefusedReleasesMakeNoThirdAttempt(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	w.fake.failOpen = 2
	w.start(w.config("0.1.0-dev"))
	count := filepath.Join(t.TempDir(), "count")
	hook := "#!/bin/sh\nwhile read old new ref; do\n[ \"$(git log -1 --format=%s \"$new\")\" = \"lease: released\" ] || continue\n" +
		"echo x >> " + count + "\n[ \"$(wc -l < " + count + ")\" -le 2 ] && exit 1\ndone\nexit 0\n"
	os.MkdirAll(filepath.Join(w.bare, "hooks"), 0o755)
	if err := os.WriteFile(filepath.Join(w.bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := w.config("0.1.0-dev")
	cfg.RunID = nextRun
	rows := Restart(context.Background(), cfg)
	last := rows[len(rows)-1]
	data, _ := os.ReadFile(count)
	if last.Name != "phase" || last.Result != "fail" || !strings.Contains(last.Detail, "the lease names the run "+nextRun) || strings.Count(string(data), "x") != 2 ||
		!strings.Contains(leaseText(t, w), "\theld\n") || strings.Contains(recordsLog(t, w), "lease: released") {
		t.Errorf("the last row %+v, %d refused pushes, the lease:\n%s\nwant phase fail naming the loss, 2 pushes, held", last, strings.Count(string(data), "x"), leaseText(t, w))
	}
}

// A records prefix is committed, and the next write fails: the release keeps the prefix, publishes nothing of the failed payload, and a fresh restart resumes.
func TestAPartialRecordsStepReleasesAndKeepsItsPrefix(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	r, ctx, rows := midRestart(t, w,
		func(r *state, ctx context.Context) Step { return r.ruleGapsStep(ctx) },
		func(r *state, ctx context.Context) Step {
			// A path with a NUL cannot be written: the commit fails before a push, and not as a refusal.
			err := r.write(ctx, map[string][]byte{"intake/a\x00b": []byte("x")}, "records: the failed payload")
			return Step{"meaning", "fail", "the write: " + strings.SplitN(err.Error(), ":", 2)[0]}
		})
	if last := rows[len(rows)-1]; last.Name != "meaning" || last.Result != "fail" {
		t.Fatalf("the rows before the release: %+v", rows)
	}
	rows = r.endLease(ctx, rows)
	if last := rows[len(rows)-1]; last.Result != "fail" || strings.Contains(last.Detail, "the lease:") {
		t.Errorf("the last row %+v; want fail with its own detail", last)
	}
	log := recordsLog(t, w)
	if !strings.HasPrefix(log, "lease: released\nlayup run: the rule gaps of the problem statement\n") || strings.Contains(log, "failed payload") ||
		!strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "intake/psb.tsv") {
		t.Errorf("the records after the release (the prefix is the rule gaps, the failed payload is no commit):\n%s", log)
	}
	again := nextTakesAtOnce(t, w, nextRun, "aaaaaaaaaaaaaaaa")
	if g := rowOf(again, "rule-gaps"); g.Result != "done" || g.Detail != "done before" {
		t.Errorf("the fresh restart's row %+v; want rule-gaps done before, the prefix resumed", g)
	}
}
