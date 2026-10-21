//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/tsv"
)

// The task register on the real adapter, the httptest forge and real git (task T-cak3, #357, row 66a of the plan of M2e).
// The step `loop` (row 66k, task T-9nbf) calls what these tests call, on the state of a restart.

// The Go value of the block tasks equals its block of docs/spec/records.md.
func TestTheSchemaOfTasksEqualsItsBlock(t *testing.T) {
	t.Parallel()
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["tasks"]
	if !ok {
		t.Fatal("docs/spec/ has no block tasks")
	}
	if err := tsv.Compare(block, TasksSchema); err != nil {
		t.Error(err)
	}
}

// loopPerms are the permissions of the installation when the Operator has granted those of the loop.
var loopPerms = map[string]string{"contents": "write", "issues": "write", "metadata": "read", "pull_requests": "write", "statuses": "write", "checks": "read"}

// registerWorld makes a world whose Start is done, and the forge answers perms from now on.
func registerWorld(t *testing.T, perms map[string]string) *world {
	t.Helper()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("1.0.0")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	w.fake.mu.Lock()
	w.fake.perms = perms
	w.fake.mu.Unlock()
	return w
}

// restarted gives the state of a restart as far as its clone: the forge step, which reads the installation, and the clone of the records.
func restarted(t *testing.T, w *world) *state {
	t.Helper()
	ctx := context.Background()
	r := &state{cfg: w.config("1.0.0")}
	if s := r.forgeStep(ctx, false); s.Result != "done" {
		t.Fatalf("the forge step of the restart: %+v", s)
	}
	if s := r.cloneStep(ctx); s.Result != "done" {
		t.Fatalf("the clone step of the restart: %+v", s)
	}
	return r
}

func recordsHead(t *testing.T, w *world) string {
	t.Helper()
	return strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records"))
}

func firstTask() Task {
	return Task{ID: "T-aaaa", Issue: 7, Kind: "specification", Requirements: []string{"REQ-005", "NFR-001"}, Predecessors: nil}
}

func TestTheLoopWritesNothingWhenThePermissionsAreMissing(t *testing.T) {
	t.Parallel()
	w := registerWorld(t, nil) // the three permissions of M2a: Start and the restart go on
	r := restarted(t, w)
	head := recordsHead(t, w)
	s := r.loopPermissionsStep()
	if s.Name != "forge" || s.Result != "fail" {
		t.Errorf("the permissions step: %+v; want forge: fail", s)
	}
	err := r.registerTask(context.Background(), firstTask())
	if err == nil {
		t.Fatal("a task was registered by an App that lacks the permissions of the loop")
	}
	for _, want := range []string{"pull_requests: write", "statuses: write", "checks: read"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the refusal %q does not name %q", err, want)
		}
	}
	if now := recordsHead(t, w); now != head {
		t.Errorf("the records head moved from %s to %s", head, now)
	}
	// The permissions granted, the same task is registered by the restart that reads them.
	w.fake.mu.Lock()
	w.fake.perms = loopPerms
	w.fake.mu.Unlock()
	r = restarted(t, w)
	if s := r.loopPermissionsStep(); s.Result != "done" {
		t.Errorf("the permissions step with the loop's permissions: %+v", s)
	}
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Errorf("with the permissions granted: %v", err)
	}
}

func TestTheFirstTaskIsOneRecordsCommitOfTheRegisterAndTheDefaultTable(t *testing.T) {
	t.Parallel()
	w := registerWorld(t, loopPerms)
	r := restarted(t, w)
	before := recordsHead(t, w)
	from := w.clock.Now().Truncate(time.Second)
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Fatal(err)
	}
	to := w.clock.Now()
	after := recordsHead(t, w)
	if n := strings.TrimSpace(gitOut(t, w.bare, "rev-list", "--count", before+".."+after)); n != "1" {
		t.Errorf("the records branch gained %s commits; want 1", n)
	}
	if parent := strings.TrimSpace(gitOut(t, w.bare, "rev-parse", after+"^")); parent != before {
		t.Errorf("the parent of the commit is %s; want the head %s", parent, before)
	}
	// Row 66b (task T-dx1h): loop.md writes parameters.tsv in the commit that writes transitions.tsv.
	if names := gitOut(t, w.bare, "diff", "--name-only", before, after); names != "parameters.tsv\ntasks.tsv\ntransitions.tsv\n" {
		t.Errorf("the commit changes %q; want parameters.tsv, tasks.tsv and transitions.tsv", names)
	}
	if who := gitOut(t, w.bare, "log", "-1", "--format=%an <%ae>|%cn <%ce>|%s", after); who != "layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>|layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>|tasks: register T-aaaa\n" {
		t.Errorf("the author, committer and subject %q; want the App's bot and %q", who, "tasks: register T-aaaa")
	}
	rows, err := ReadRegister([]byte(show(t, w.bare, after, "tasks.tsv")))
	if err != nil || len(rows) != 1 {
		t.Fatalf("tasks.tsv holds %+v, %v", rows, err)
	}
	if rows[0].Opened.Before(from) || rows[0].Opened.After(to) {
		t.Errorf("opened %s is outside the run's clock, %s to %s", rows[0].Opened, from, to)
	}
	want := firstTask()
	want.Opened = rows[0].Opened
	if !reflect.DeepEqual(rows[0], want) {
		t.Errorf("tasks.tsv holds %+v; want %+v", rows[0], want)
	}
	var def bytes.Buffer
	if err := handoff.WriteTable(&def, handoff.Default()); err != nil {
		t.Fatal(err)
	}
	if got := show(t, w.bare, after, "transitions.tsv"); got != def.String() {
		t.Errorf("transitions.tsv:\n%s\nwant the default table:\n%s", got, def.String())
	}
}

func TestASecondTaskWritesOnlyTheRegisterAndKeepsTheOperatorsTable(t *testing.T) {
	t.Parallel()
	w := registerWorld(t, loopPerms)
	r := restarted(t, w)
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Fatal(err)
	}
	// The Operator replaces the table (row 66b makes the replacement; here it is a file of the records).
	custom := append(handoff.Default(), handoff.Row{From: "verifier", To: "asking role", Artifacts: []string{"question"}, Items: []string{"schema"}})
	var b bytes.Buffer
	if err := handoff.WriteTable(&b, custom); err != nil {
		t.Fatal(err)
	}
	pushRecords(t, w.bare, map[string]string{"transitions.tsv": b.String()}, "the Operator's table")
	r = restarted(t, w)
	before := recordsHead(t, w)
	second := Task{ID: "T-bbbb", Issue: 8, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Predecessors: []string{"T-aaaa"}}
	if err := r.registerTask(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	after := recordsHead(t, w)
	if names := gitOut(t, w.bare, "diff", "--name-only", before, after); names != "tasks.tsv\n" {
		t.Errorf("the commit of the second task changes %q; want tasks.tsv alone", names)
	}
	if got := show(t, w.bare, after, "transitions.tsv"); got != b.String() {
		t.Errorf("transitions.tsv after the second task:\n%s\nwant the Operator's table:\n%s", got, b.String())
	}
	rows, err := ReadRegister([]byte(show(t, w.bare, after, "tasks.tsv")))
	if err != nil || len(rows) != 2 || rows[0].ID != "T-aaaa" || rows[1].ID != "T-bbbb" {
		t.Errorf("tasks.tsv holds %+v, %v; want the two tasks in order", rows, err)
	}
}

func TestARefusedTaskWritesNothing(t *testing.T) {
	t.Parallel()
	noReq := firstTask()
	noReq.Requirements = nil
	self := firstTask()
	self.Predecessors = []string{"T-aaaa"}
	absent := firstTask()
	absent.Predecessors = []string{"T-zzzz"}
	var def bytes.Buffer
	handoff.WriteTable(&def, handoff.Default()[1:])
	cycle := "task\tissue\tkind\trequirements\tsize\tpredecessors\tpr\topened\n" +
		"T-cccc\t5\tbuild\tREQ-005\t—\tT-dddd\t—\t2026-10-16T07:00:00Z\nT-dddd\t6\tbuild\tREQ-005\t—\tT-cccc\t—\t2026-10-16T07:00:00Z\n"
	for _, c := range []struct {
		name  string
		files map[string]string // the records before the write
		task  Task
		want  string
	}{
		{"no requirement", nil, noReq, "requirement"},
		{"a task that is its own predecessor", nil, self, "cycle"},
		{"a predecessor that is no task", nil, absent, "T-zzzz"},
		{"a table that lacks a default pair", map[string]string{"transitions.tsv": def.String()}, firstTask(), "transitions.tsv"},
		{"a register that read refuses", map[string]string{"tasks.tsv": cycle}, Task{ID: "T-eeee", Issue: 9, Kind: "build", Requirements: []string{"REQ-005"}}, "tasks.tsv"},
	} {
		w := registerWorld(t, loopPerms)
		if c.files != nil {
			pushRecords(t, w.bare, c.files, "records of the test")
		}
		r := restarted(t, w)
		head := recordsHead(t, w)
		err := r.registerTask(context.Background(), c.task)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want a refusal that names %q", c.name, err, c.want)
		}
		if now := recordsHead(t, w); now != head {
			t.Errorf("%s: the records head moved from %s to %s", c.name, head, now)
		}
		for _, f := range []string{"tasks.tsv", "transitions.tsv"} {
			_, had := c.files[f]
			if out := gitOut(t, w.bare, "ls-tree", "--name-only", head, f); (out != "") != had {
				t.Errorf("%s: %s is in the records head: %v; was before: %v", c.name, f, out != "", had)
			}
		}
	}
}

func TestTheRegisterWriteIsFencedByTheLease(t *testing.T) {
	t.Parallel()
	// refuseOnce makes the bare repository refuse the next push, once, as a head that moved is refused.
	refuseOnce := func(t *testing.T, w *world) {
		hook := "#!/bin/sh\nif [ -f refuse-once ]; then rm refuse-once; echo refused >&2; exit 1; fi\n"
		if err := os.WriteFile(filepath.Join(w.bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(w.bare, "refuse-once"), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	held := func(r *state, run string) {
		t.Helper()
		row := LeaseRow{Run: run, Host: "host-1", Version: "1.0.0", Started: opened, Heartbeat: 1, State: "held"}
		if err := r.store.WriteLease(context.Background(), row, "lease: held by the test"); err != nil {
			t.Fatal(err)
		}
	}
	t.Run("a push that is refused once is tried again while the lease is held", func(t *testing.T) {
		w := registerWorld(t, loopPerms)
		r := restarted(t, w)
		held(r, r.cfg.RunID)
		before := recordsHead(t, w)
		refuseOnce(t, w)
		if err := r.registerTask(context.Background(), firstTask()); err != nil {
			t.Fatalf("the write: %v", err)
		}
		if n := strings.TrimSpace(gitOut(t, w.bare, "rev-list", "--count", before+".."+recordsHead(t, w))); n != "1" {
			t.Errorf("the records branch gained %s commits; want 1", n)
		}
	})
	t.Run("a lease held by another run is lost, and nothing is written", func(t *testing.T) {
		w := registerWorld(t, loopPerms)
		r := restarted(t, w)
		held(r, "fedcba9876543210")
		before := recordsHead(t, w)
		refuseOnce(t, w)
		err := r.registerTask(context.Background(), firstTask())
		var lost *LostError
		if !errors.As(err, &lost) || lost.Holder != "fedcba9876543210" {
			t.Errorf("the write: %v; want a LostError whose holder is another-run", err)
		}
		if now := recordsHead(t, w); now != before {
			t.Errorf("the records head moved from %s to %s", before, now)
		}
	})
}
