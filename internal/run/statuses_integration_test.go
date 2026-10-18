//go:build integration

package run

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The block task-gates is the Go value GatesSchema (the owner's comparison that
// moves it from notYetBuilt to built).
func TestTheGatesBlockEqualsItsGoSchema(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["task-gates"]
	if !ok {
		t.Fatal("docs/spec/ has no block task-gates")
	}
	if err := tsv.Compare(block, GatesSchema); err != nil {
		t.Error(err)
	}
}

// statusWorldOfARun is the world of a run with real git, the records branch, and
// the adapter to the stand-in forge, where the App's bot is the ID 9001.
func statusWorldOfARun(t *testing.T) *loopRun {
	l := newLoop(t, nil)
	l.r.Bot = 9001
	return l
}

// childCommit makes a commit in the clone, a child of parent, with the files, and gives its SHA.
func childCommit(t *testing.T, clone, parent string, files map[string]string) string {
	t.Helper()
	gitOut(t, clone, "checkout", "-q", "--detach", parent)
	for p, text := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(clone, p)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(clone, p), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "test")
	return strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))
}

const gatesManifest = "kind\tstate\ttool\tcommand\tscope\tconfig\n"

// gateWorld gives a base with the manifest and a head that changes x.go.
func gateWorld(t *testing.T, kinds string) (*loopRun, string, string) {
	t.Helper()
	l := statusWorldOfARun(t)
	base := childCommit(t, l.r.Clone, l.a.base, map[string]string{"docs/gates.tsv": gatesManifest + kinds, "x.go": "package x\n"})
	head := childCommit(t, l.r.Clone, base, map[string]string{"x.go": "package x // changed\n"})
	return l, base, head
}

func statusesOf(w *world, sha string) []map[string]any {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	return append([]map[string]any(nil), w.fake.statuses[sha]...)
}

func TestTheGatesStatusIsPostedOnceFromTheRecordedTableAndItsOutputIsKept(t *testing.T) {
	count := filepath.Join(t.TempDir(), "count")
	l, base, head := gateWorld(t, "fine\tactive\tsh\techo ran >> "+count+"; echo held-output-of-fine\t./*.go\t—\ndocs\tpending\tsh\ttrue\t./*.md\t—\n")
	ctx := context.Background()
	var out bytes.Buffer
	var steps []string
	in := GatesInput{Task: "T-ab12", Base: base, Head: head, Out: &out, Step: func(i, n int, kind string) func() {
		steps = append(steps, kind)
		return func() {}
	}}
	st, err := l.r.GatesAt(ctx, in)
	want := Status{GatesContext, "success", "clear: pending: no product path"}
	if err != nil || st != want {
		t.Fatalf("GatesAt = %+v, %v; want %+v", st, err, want)
	}
	if !strings.Contains(out.String(), "held-output-of-fine") || strings.Join(steps, " ") != "fine docs" {
		t.Errorf("output %q, steps %v; want the held output of the command and the step of each kind", out.String(), steps)
	}
	if ok, err := l.r.PostStatus(ctx, head, st); err != nil || !ok {
		t.Fatalf("PostStatus: %v, %v", ok, err)
	}
	// A second call reads the table of the head from the records: no second gate run, no second post.
	out.Reset()
	if st2, err := l.r.GatesAt(ctx, in); err != nil || st2 != st || out.Len() != 0 {
		t.Fatalf("the second GatesAt = %+v, %v, output %q; want the same status from the table and no run", st2, err, out.String())
	}
	if ok, err := l.r.PostStatus(ctx, head, st); err != nil || ok {
		t.Errorf("the second PostStatus: %v, %v; want not posted", ok, err)
	}
	if data, _ := os.ReadFile(count); strings.Count(string(data), "ran") != 1 {
		t.Errorf("the command of the kind ran %q times; want once", string(data))
	}
	if got := statusesOf(l.w, head); len(got) != 1 || got[0]["context"] != GatesContext || got[0]["state"] != "success" {
		t.Errorf("the forge holds %v; want the one status", got)
	}
	tbl, found, err := l.r.GatesTable(ctx, "T-ab12", head)
	if err != nil || !found || len(tbl.Rows) != 2 || tbl.Base != base || tbl.Rows[0].Result != "pass" {
		t.Errorf("GatesTable = %+v, %v, %v; want the two rows of the head", tbl, found, err)
	}
}

func TestAFailedPostIsRecoveredFromTheSavedTableWithoutRunningTheGateAgain(t *testing.T) {
	count := filepath.Join(t.TempDir(), "count")
	l, base, head := gateWorld(t, "bad\tactive\tsh\techo ran >> "+count+"; exit 3\t./*.go\t—\n")
	ctx := context.Background()
	st, err := l.r.GatesAt(ctx, GatesInput{Task: "T-ab12", Base: base, Head: head})
	if want := (Status{GatesContext, "failure", "bad: fail: exit 3"}); err != nil || st != want {
		t.Fatalf("GatesAt = %+v, %v; want %+v", st, err, want)
	}
	l.w.fake.failStatus = 1 // the forge refuses the post, after the table is committed
	if ok, err := l.r.PostStatus(ctx, head, st); err == nil || ok {
		t.Fatalf("the failed post: %v, %v; want an error", ok, err)
	}
	if n := len(statusesOf(l.w, head)); n != 0 {
		t.Fatalf("%d statuses after the failed post; want none", n)
	}
	r2 := afterRestart(l.r)
	if r2.Bot != 9001 {
		t.Fatalf("the Sessions of a restart have the bot %d; want 9001", r2.Bot)
	}
	st2, err := r2.GatesAt(ctx, GatesInput{Task: "T-ab12", Base: base, Head: head})
	if err != nil || st2 != st {
		t.Fatalf("GatesAt after the restart = %+v, %v; want %+v from the saved table", st2, err, st)
	}
	if ok, err := r2.PostStatus(ctx, head, st2); err != nil || !ok {
		t.Fatalf("the post after the restart: %v, %v", ok, err)
	}
	if data, _ := os.ReadFile(count); strings.Count(string(data), "ran") != 1 {
		t.Errorf("the gate ran %q times; want once, the saved table is the record", string(data))
	}
	rows, err := ReadGates([]byte(keptRecord(t, l.w, "tasks/T-ab12/gates.tsv")))
	if err != nil || len(rows) != 1 || rows[0][0] != head || rows[0][4] != "fail" {
		t.Errorf("gates.tsv = %v, %v; want the one row of the head", rows, err)
	}
}

func TestAPendingKindWithAProductPathAndANotActiveKindAreFailures(t *testing.T) {
	for name, c := range map[string]struct{ kinds, want string }{
		"a pending kind and a product path that changed": {"boundary\tpending\tsh\ttrue\t./*.go\t—\n", "boundary: fail: pending: product path changed: x.go"},
		"a tool that is not found":                       {"lint\tactive\tno-such-tool-for-layup\ttrue\t./*.go\t—\n", "lint: not-active: tool not found: no-such-tool-for-layup"},
	} {
		l, base, head := gateWorld(t, c.kinds)
		st, err := l.r.GatesAt(context.Background(), GatesInput{Task: "T-ab12", Base: base, Head: head})
		if want := (Status{GatesContext, "failure", c.want}); err != nil || st != want {
			t.Errorf("%s: %+v, %v; want %+v", name, st, err, want)
		}
		if ok, err := l.r.PostStatus(context.Background(), head, st); err != nil || !ok || statusesOf(l.w, head)[0]["state"] != "failure" {
			t.Errorf("%s: posted %v, %v, read back %v; want a failure", name, ok, err, statusesOf(l.w, head))
		}
	}
}

func TestABaseWithNoManifestIsAPostedErrorAndWritesNoRow(t *testing.T) {
	l := statusWorldOfARun(t)
	head := childCommit(t, l.r.Clone, l.a.base, map[string]string{"x.go": "package x\n"})
	st, err := l.r.GatesAt(context.Background(), GatesInput{Task: "T-ab12", Base: l.a.base, Head: head})
	if err != nil || st.State != "error" || !strings.HasPrefix(st.Description, "error: docs/gates.tsv at the base") {
		t.Errorf("GatesAt = %+v, %v; want error, error: docs/gates.tsv at the base …", st, err)
	}
	if _, found, err := l.r.GatesTable(context.Background(), "T-ab12", head); err != nil || found {
		t.Errorf("GatesTable: found=%v, %v; want no row for an input error", found, err)
	}
	if st, err := l.r.GatesAt(context.Background(), GatesInput{Task: "T-ab12", Base: "main", Head: head}); err == nil {
		t.Errorf("a base that is no full SHA: %+v, no error", st)
	}
}

func TestTheRulesStatusReadsTheRegisterOfTheRecordsCommit(t *testing.T) {
	l := statusWorldOfARun(t)
	ctx := context.Background()
	// A head that code produced edits a rule path and the register at the root; the caller names the records commit.
	head := childCommit(t, l.r.Clone, l.a.base, map[string]string{"AGENTS.md": "changed\n", "rule-paths.tsv": "pattern\texception\tsource\nnothing\t—\tbaseline\n"})
	st, err := l.r.RulesAt(ctx, RulesInput{Task: "T-ab12", Base: l.a.base, Head: head, RegisterRev: l.a.records})
	if want := (Status{RulesContext, "failure", "rule path: AGENTS.md"}); err != nil || st != want {
		t.Errorf("RulesAt = %+v, %v; want %+v (the register of the records commit, not the head's)", st, err, want)
	}
	if _, err := l.r.RulesAt(ctx, RulesInput{Task: "T-ab12", Base: l.a.base, Head: head}); err == nil {
		t.Error("a head with no bound session and no commit of the caller: no error")
	}
	// A head that a session pushed: the register commit is the records column of its start row.
	id := l.dev(1, "claude", map[string]string{"x.go": "package x\n"}, okResult)
	bound := l.bound(id)
	if st, err := l.r.RulesAt(ctx, RulesInput{Task: "T-ab12", Base: l.a.base, Head: bound}); err != nil || st != (Status{RulesContext, "success", "no rule path"}) {
		t.Errorf("RulesAt of the pushed head = %+v, %v; want success, no rule path", st, err)
	}
	if _, err := l.r.RulesAt(ctx, RulesInput{Task: "T-ab12", Base: l.a.base, Head: bound, RegisterRev: l.a.records}); err == nil {
		t.Error("a commit of the caller for a head that has a bound session: no error")
	}
}

func TestTheVerifyStatusFollowsTheVerifierFromTheRecords(t *testing.T) {
	l := statusWorldOfARun(t)
	ctx := context.Background()
	s, harnesses, models := pairWorld(t, "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\n")
	files := map[string][]byte{}
	for _, p := range []string{"harnesses.tsv", "routing.tsv"} {
		files[p] = s.Store.(*standinStore).files[p]
	}
	if err := l.r.commit(ctx, files, "the registers"); err != nil {
		t.Fatal(err)
	}
	// The verifier is a session of devin that reviews the base; its record is nothing material at that head.
	text := record(l.a.base, "nothing material in scope")
	id := l.comment(1, "verification", "devin", text)
	in := VerifyInput{Task: "T-ab12", Base: l.a.base, Head: l.a.base, Session: id, Harnesses: harnesses, Models: models}
	sha, st, err := l.r.VerifyAt(ctx, in)
	if want := (Status{VerifyContext, "success", "nothing material in scope"}); err != nil || sha != l.a.base || st != want {
		t.Fatalf("VerifyAt = %s, %+v, %v; want %+v on the reviewed head", sha, st, err, want)
	}
	if ok, err := l.r.PostStatus(ctx, sha, st); err != nil || !ok {
		t.Fatalf("PostStatus: %v, %v", ok, err)
	}
	if ok, err := l.r.PostStatus(ctx, sha, st); err != nil || ok {
		t.Errorf("the same verdict again: %v, %v; want not posted", ok, err)
	}
}
