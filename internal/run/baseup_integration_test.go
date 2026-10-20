//go:build integration

package run

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
)

// The world of a task whose base moved (task T-89xf, #424, row 66x; loop.md, A task whose base moved): the ship world, a developer's
// commit X of the task on the task branch, verified by a verifier of devin that pushes nothing, and the means to move the default branch of
// the forge. The Sessions are those of the production constructor.
type moved struct {
	*ship
	B, X string
	s    *Sessions
}

const movedBinding = "option\tconstraint\tresult\tevidence\ndb\tREQ-001\tpass\tdocs/spec/db.md\n"

// newMoved is newMovedAs for an architecture task of a decision that a bet opened, with its binding.
func newMoved(t *testing.T, taskFiles map[string]string) *moved {
	return newMovedAs(t, taskFiles, "architecture", true)
}

// newMovedAs makes the world of a task of the kind: with the binding of a decision, or with the draft of Intake in the records its first
// session read (a Shape task).
func newMovedAs(t *testing.T, taskFiles map[string]string, kind string, binding bool) *moved {
	t.Helper()
	ctx := context.Background()
	l := newLoop(t, shipBaseline)
	l.r.Bot = 9001
	harnesses, models := commitRegisters(t, l)
	tasks, err := ReadRegister([]byte(show(t, l.w.bare, "layup-records", "tasks.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].Kind = kind
	var reg strings.Builder
	if err := WriteRegister(&reg, tasks); err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{tasksPath: []byte(reg.String())}
	if binding {
		files[bindingPath("T-ab12", "db")] = []byte(movedBinding)
	} else {
		files["spec/draft.tsv"] = []byte("the draft of Intake\n")
	}
	if err := l.r.commit(ctx, files, "the kind and its inputs"); err != nil {
		t.Fatal(err)
	}
	l.a.records = l.r.Store.(*recordsStore).pushed
	s := &ship{loopRun: l, V: l.a.base, in: CarryInput{Task: "T-ab12", Harnesses: harnesses, Models: models}}
	id := s.dev(1, "claude", taskFiles, okResult)
	x := s.bound(id)
	if x == "" {
		t.Fatal("the developer's session has no bound head")
	}
	gitOut(t, s.r.Clone, "branch", "-f", "layup/base/T-ab12/1", x)
	text := record(x, "nothing material in scope")
	v := harness(t, s.a, planWork("plan.md", text, "comment\t1\tresult/plan.md\t"+commentSum(text)+"\t—"))
	v.Role, v.Attempt, v.Pair.Harness, v.Base, v.Branch = roleOfStep("verification"), 1, "devin", x, "layup/base/T-ab12/1"
	v.Push = func(string, session.Dir, session.Run) error { return nil }
	vid, err := s.r.TaskSession(ctx, v)
	if err != nil {
		t.Fatalf("the verifier: %q, %v", vid, err)
	}
	if _, err := s.r.PostComment(ctx, "T-ab12", vid, "plan.md"); err != nil {
		t.Fatal(err)
	}
	s.in.Verifier = vid
	// The register names the issue 3 for the task; the forge has opened two, so a new issue would take that number.
	if n, err := s.r.Forge.OpenIssue(ctx, "LAYUP task T-ab12", "the issue of the task"); err != nil || n != taskIssue {
		t.Fatalf("the issue of the task = %d, %v; want %d", n, err, taskIssue)
	}
	return &moved{ship: s, B: s.V, X: x, s: s.prodBot()}
}

// moveBase commits files on the default branch of the forge, as another task's merge would, and gives the commit.
func (m *moved) moveBase(files map[string]string) string {
	m.t.Helper()
	dir := m.t.TempDir()
	gitOut(m.t, "/", "clone", "-q", "-b", "main", "file://"+m.w.bare, dir)
	for p, d := range files {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, p)), 0o755); err != nil {
			m.t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, p), []byte(d), 0o644); err != nil {
			m.t.Fatal(err)
		}
	}
	gitOut(m.t, dir, "add", "-A")
	gitOut(m.t, dir, "commit", "-q", "-m", "the base moves")
	gitOut(m.t, dir, "push", "-q", "origin", "main")
	return strings.TrimSpace(gitOut(m.t, dir, "rev-parse", "HEAD"))
}

func (m *moved) bring(s *Sessions) (Advance, error) {
	m.t.Helper()
	return s.BringUpToDate(context.Background(), m.in)
}

// writes gives the requests of the forge that write, after the first n requests: a comment, a status and a token are not among them.
func (m *moved) writes(n int) []string {
	m.w.fake.mu.Lock()
	defer m.w.fake.mu.Unlock()
	var out []string
	for _, r := range m.w.fake.reqs[n:] {
		if !strings.HasPrefix(r, "GET ") && !strings.Contains(r, "/comments") && !strings.Contains(r, "/statuses/") && !strings.Contains(r, "/access_tokens") {
			out = append(out, r)
		}
	}
	return out
}

func (m *moved) nreqs() int {
	m.w.fake.mu.Lock()
	defer m.w.fake.mu.Unlock()
	return len(m.w.fake.reqs)
}

func (m *moved) eventsOfTask(task string) [][]string {
	m.t.Helper()
	rows, err := records.ReadEvents([]byte(keptRecord(m.t, m.w, "tasks/"+task+"/events.tsv")))
	if err != nil {
		m.t.Fatal(err)
	}
	return rows
}

func (m *moved) kindsOf(task, kind string) [][]string {
	return slices.DeleteFunc(m.eventsOfTask(task), func(e []string) bool { return e[1] != kind })
}

func (m *moved) tasks() []Task {
	m.t.Helper()
	tasks, err := ReadRegister([]byte(keptRecord(m.t, m.w, "tasks.tsv")))
	if err != nil {
		m.t.Fatal(err)
	}
	return tasks
}

func (m *moved) tip() string {
	return strings.TrimSpace(gitOut(m.t, m.w.bare, "rev-parse", "refs/heads/task/T-ab12/1"))
}

// pr opens the draft pull request of the attempt on the forge and registers it.
func (m *moved) pr() forge.PullRequest {
	m.t.Helper()
	pr, err := m.s.OpenPullRequest(context.Background(), PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "the work of the task"})
	if err != nil {
		m.t.Fatal(err)
	}
	return pr
}

var (
	taskPaths = map[string]string{"t.txt": "task\n"}
	basePaths = map[string]string{"u.txt": "base\n"}
)

// ----- the faults -----

var errStop = errors.New("stopped by the test")

// baseStoreFault is a records store that fails as a lost reply does: after the commit of a file it names, or with a lease that is lost from
// the n-th read, or with a fetch that fails.
type baseStoreFault struct {
	SessionStore
	after     string // a commit that writes this path is made and then fails with errStop; "" for none
	failIf    func(files map[string][]byte) bool
	fetchErr  error
	fetched   string // FetchBranch gives this commit of the clone, whatever the forge holds
	pushErr   error  // PushHead fails with this, and pushes nothing
	lostFrom  int    // the lease is another run's from the n-th read; 0 for never
	leaseRead int
}

func (f *baseStoreFault) Commit(ctx context.Context, files map[string][]byte, message string) error {
	if f.failIf != nil && f.failIf(files) {
		return errStop
	}
	if err := f.SessionStore.Commit(ctx, files, message); err != nil {
		return err
	}
	if _, ok := files[f.after]; ok && f.after != "" {
		return errStop
	}
	return nil
}

func (f *baseStoreFault) FetchBranch(ctx context.Context, branch string) (string, error) {
	if f.fetchErr != nil {
		return "", f.fetchErr
	}
	if f.fetched != "" {
		return f.fetched, nil
	}
	return f.SessionStore.FetchBranch(ctx, branch)
}

func (f *baseStoreFault) ReadLease(ctx context.Context) (LeaseRow, error) {
	f.leaseRead++
	row, err := f.SessionStore.ReadLease(ctx)
	if f.lostFrom > 0 && f.leaseRead >= f.lostFrom {
		row.Run = "another-run-0000"
	}
	return row, err
}

func (f *baseStoreFault) PushHead(ctx context.Context, sha, branch string) error {
	if f.pushErr != nil {
		return f.pushErr
	}
	return f.SessionStore.PushHead(ctx, sha, branch)
}

// baseForgeFault is the forge whose named call takes effect and then fails, and whose close of a pull request may answer merged.
type baseForgeFault struct {
	forge.Forge
	after       string // ClosePullRequest, CloseIssue or OpenIssue
	closeMerged bool
	prErr       error // the read of a pull request fails with this
}

func (f *baseForgeFault) PullRequest(ctx context.Context, q forge.PullRequestQuery) (forge.PullRequest, bool, error) {
	if err := f.prErr; err != nil {
		f.prErr = nil // the first read fails, the others do not
		return forge.PullRequest{}, false, err
	}
	return f.Forge.PullRequest(ctx, q)
}

func (f *baseForgeFault) ClosePullRequest(ctx context.Context, n int) (forge.PullRequest, error) {
	pr, err := f.Forge.ClosePullRequest(ctx, n)
	if f.closeMerged {
		pr.Merged = true
	}
	if err == nil && f.after == "ClosePullRequest" {
		err = errStop
	}
	return pr, err
}

func (f *baseForgeFault) CloseIssue(ctx context.Context, n int) error {
	err := f.Forge.CloseIssue(ctx, n)
	if err == nil && f.after == "CloseIssue" {
		err = errStop
	}
	return err
}

func (f *baseForgeFault) OpenIssue(ctx context.Context, title, body string) (int, error) {
	n, err := f.Forge.OpenIssue(ctx, title, body)
	if err == nil && f.after == "OpenIssue" {
		err = errStop
	}
	return n, err
}

// ----- the up-to-date task, and the carry -----

// An up-to-date task is left as it is: no request, no records commit (a green negative at the stub beside the carry's red).
func TestAnUpToDateVerifiedTaskIsLeftAsItIs(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	head, n := recordsHead(t, m.w), m.nreqs()
	got, err := m.bring(m.s)
	if err != nil || got.Task != "T-ab12" || !got.Verified || got.Head != m.X {
		t.Fatalf("BringUpToDate = %+v, %v; want {T-ab12 verified at %s}", got, err, m.X)
	}
	if current, def, err := m.s.UpToDate(context.Background(), m.X); err != nil || !current || def != m.B {
		t.Errorf("UpToDate(X) = %v, %s, %v; want current at the default head %s", current, def, err, m.B)
	}
	if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 || m.tip() != m.X {
		t.Errorf("an up-to-date task changed the records (%v), the forge (%v) or the branch (%s)", recordsHead(t, m.w) != head, m.writes(n), m.tip())
	}
}

// I3, the carry demo: the default branch moves on files the task does not touch; the merge is one commit of the two heads, pushed as
// code, with the events push and bound of session —, and layup/verify is carried over it; a close-out on it is accepted; a second call
// writes nothing.
func TestADisjointMovedBaseIsMergedAndCarriesTheVerdict(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	d := m.moveBase(basePaths)
	ctx := context.Background()
	m.s.Now = func() time.Time { return time.Date(2031, 1, 1, 0, 0, 0, 0, time.UTC) }
	got, err := m.bring(m.s)
	if err != nil || !got.Verified || got.Replaced != "" || got.Head == "" || got.Head == m.X {
		t.Fatalf("BringUpToDate = %+v, %v; want a verified merge commit", got, err)
	}
	merge := got.Head
	if parents, err := git.Parents(m.r.Clone, merge); err != nil || !slices.Equal(parents, []string{m.X, d}) {
		t.Errorf("the parents of %s = %v, %v; want [%s %s]", merge, parents, err, m.X, d)
	}
	if m.tip() != merge {
		t.Errorf("the branch holds %s; want the merge %s", m.tip(), merge)
	}
	var coded [][]string
	for _, e := range m.eventsOfTask("T-ab12") {
		if (e[1] == "push" || e[1] == "bound") && e[5] == merge {
			coded = append(coded, e)
		}
	}
	if len(coded) != 2 || coded[0][1] != "push" || coded[1][1] != "bound" || coded[0][3] != "" || coded[1][3] != "" || coded[1][6] != "task/T-ab12/1" {
		t.Errorf("the events of the merge: %v; want push then bound, session —, branch task/T-ab12/1", coded)
	}
	if obj := gitOut(t, m.r.Clone, "cat-file", "-p", merge); !strings.Contains(obj, "committer layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>") {
		t.Errorf("the merge commit is not the App's:\n%s", obj)
	}
	ct := func(rev string) int {
		n, _ := strconv.Atoi(strings.TrimSpace(gitOut(t, m.r.Clone, "log", "-1", "--format=%ct", rev)))
		return n
	}
	if got, want := ct(merge), max(ct(m.X), ct(d)); got != want {
		t.Errorf("the committer time of the merge is %d; want the later of the two heads, %d (never the clock of the run)", got, want)
	}
	if posted := statusesOf(m.w, merge); len(posted) != 1 || posted[0]["context"] != VerifyContext || posted[0]["state"] != "success" ||
		posted[0]["description"] != "verified head "+m.X+" plus a clean merge of the base" {
		t.Errorf("the statuses of the merge: %v; want layup/verify success, verified head %s plus a clean merge of the base", posted, m.X)
	}
	again := m.in
	again.Merge = merge
	if h, err := m.s.CarryVerify(ctx, again, ""); err != nil || h != merge {
		t.Errorf("CarryVerify on the merge = %q, %v", h, err)
	}
	head, n := recordsHead(t, m.w), m.nreqs()
	if second, err := m.bring(m.s); err != nil || second != got || recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("the second call = %+v, %v; writes %v; want the same answer and no write", second, err, m.writes(n))
	}
	// A close-out on the merge is accepted: the push hook calls carries, which accepts the merge.
	gitOut(t, m.r.Clone, "branch", "-f", "layup/base/T-ab12/m", merge)
	spec := harness(t, m.a, closing())
	spec.Role, spec.Base, spec.Branch, spec.Pair.Harness = roleOfStep("close-out"), merge, "layup/base/T-ab12/m", "claude"
	spec.Push = func(id string, d session.Dir, r session.Run) error { return m.r.PushCloseOut(ctx, id, m.in) }
	id, err := m.r.TaskSession(ctx, spec)
	if err != nil || m.bound(id) == "" || m.refusalOf(id) != "" {
		t.Errorf("a close-out on the merge: %q, %v, bound %q, refused %q; want it accepted", id, err, m.bound(id), m.refusalOf(id))
	}
}

// A merge that is not the production merge is not carried: each case below is a commit of the clone that CarryVerify must refuse, with no
// status posted. The first is the counterexample of condition 1 of the plan review: both parents change a file, and the merge takes the
// base's version, so the result diffs are {a.txt} and {}. A real merge of the same shape as the cases is the control that carries.
func TestAMergeThatTouchesATaskPathOrLeavesTheDefaultIsNotCarried(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	commitTree := func(m *moved, tree string, parents ...string) string {
		args := []string{"commit-tree", tree, "-m", "fabricated"}
		for _, p := range parents {
			args = append(args, "-p", p)
		}
		return strings.TrimSpace(gitOut(t, m.r.Clone, args...))
	}
	treeOf := func(m *moved, rev string) string {
		return strings.TrimSpace(gitOut(t, m.r.Clone, "rev-parse", rev+"^{tree}"))
	}
	refused := func(m *moved, name, bad string) {
		in := m.in
		in.Merge = bad
		if h, err := m.s.CarryVerify(ctx, in, ""); err == nil || len(statusesOf(m.w, bad)) != 0 {
			t.Errorf("%s: CarryVerify = %q, %v with %d statuses; want an error and no status", name, h, err, len(statusesOf(m.w, bad)))
		}
	}
	// the task and the base both change a.txt
	m := newMoved(t, map[string]string{"a.txt": "task\n"})
	d := m.moveBase(map[string]string{"a.txt": "base\n", "u.txt": "base\n"})
	m.fetchMain()
	side := m.moveBaseFree(map[string]string{"s.txt": "side\n"}) // a commit that no branch of the forge holds
	refused(m, "both parents changed a file, and the merge takes the base's version", commitTree(m, treeOf(m, d), m.X, d))
	refused(m, "a commit of one parent", commitTree(m, treeOf(m, m.X), m.X))
	gitOut(t, m.r.Clone, "switch", "-q", "--detach", m.X)
	gitOut(t, m.r.Clone, "merge", "-q", "--no-ff", "-m", "a clean merge of a side commit", side)
	offDefault := strings.TrimSpace(gitOut(t, m.r.Clone, "rev-parse", "HEAD"))
	gitOut(t, m.r.Clone, "switch", "-q", "main")
	refused(m, "a clean merge whose second parent is no ancestor of the default head", offDefault)
	refused(m, "a first parent that does not carry the verdict", commitTree(m, treeOf(m, d), side, d))
	// the paths are disjoint
	m = newMoved(t, taskPaths)
	d = m.moveBase(basePaths)
	m.fetchMain()
	real, clean, err := git.Merge(m.r.Clone, m.X, d, "the merge", git.Identity{Name: "layup-agent[bot]", Email: "9001+layup-agent[bot]@users.noreply.github.com", Time: timeOfCommit(t, m, d)})
	if err != nil || !clean {
		t.Fatalf("the control merge: %v, %v", clean, err)
	}
	gitOut(t, m.r.Clone, "switch", "-q", "--detach", m.X)
	gitOut(t, m.r.Clone, "merge", "-q", "--no-ff", "--no-commit", d)
	os.WriteFile(filepath.Join(m.r.Clone, "evil.txt"), []byte("evil\n"), 0o644)
	gitOut(t, m.r.Clone, "add", "-A")
	gitOut(t, m.r.Clone, "commit", "-q", "-m", "a merge with a file of its own")
	extra := strings.TrimSpace(gitOut(t, m.r.Clone, "rev-parse", "HEAD"))
	gitOut(t, m.r.Clone, "switch", "-q", "main")
	refused(m, "a merge that adds a file of its own", extra)
	refused(m, "a merge that drops a file of the task", commitTree(m, treeOf(m, d), m.X, d))
	refused(m, "a merge that drops the change of the base", commitTree(m, treeOf(m, m.X), m.X, d))
	in := m.in
	in.Merge = real
	if h, err := m.s.CarryVerify(ctx, in, ""); err != nil || h != real {
		t.Errorf("the control merge: CarryVerify = %q, %v; want it carried", h, err)
	}
}

func timeOfCommit(t *testing.T, m *moved, rev string) time.Time {
	t.Helper()
	at, err := git.CommitTime(m.r.Clone, rev)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

// fetchMain brings the default branch of the forge into the clone, as the first read of BringUpToDate does.
func (m *moved) fetchMain() {
	m.t.Helper()
	gitOut(m.t, m.r.Clone, "fetch", "-q", "--update-head-ok", "file://"+m.w.bare, "refs/heads/main:refs/heads/main")
}

// moveBaseFree makes a commit on a side of the clone that no branch of the forge holds, and gives it.
func (m *moved) moveBaseFree(files map[string]string) string {
	m.t.Helper()
	gitOut(m.t, m.r.Clone, "switch", "-q", "--detach", m.B)
	defer gitOut(m.t, m.r.Clone, "switch", "-q", "main")
	for p, d := range files {
		if err := os.WriteFile(filepath.Join(m.r.Clone, p), []byte(d), 0o644); err != nil {
			m.t.Fatal(err)
		}
	}
	gitOut(m.t, m.r.Clone, "add", "-A")
	gitOut(m.t, m.r.Clone, "commit", "-q", "-m", "a side commit")
	return strings.TrimSpace(gitOut(m.t, m.r.Clone, "rev-parse", "HEAD"))
}

// A failed git read is an error and never a redo: no request, no records commit, no PATCH and no POST.
func TestAFailedGitReadIsAnErrorNeverARedo(t *testing.T) {
	t.Parallel()
	for name, tweak := range map[string]func(m *moved, s *Sessions){
		"a fetch that fails": func(m *moved, s *Sessions) {
			s.Store = &baseStoreFault{SessionStore: s.Store, fetchErr: errors.New("the fetch failed")}
		},
		"a default branch that is not a descendant of the base": func(m *moved, s *Sessions) {
			root := strings.TrimSpace(gitOut(t, m.r.Clone, "commit-tree", "-m", "an unrelated root", strings.TrimSpace(gitOut(t, m.r.Clone, "rev-parse", m.B+"^{tree}"))))
			s.Store = &baseStoreFault{SessionStore: s.Store, fetched: root}
		},
		"a read of the pull request that fails": func(m *moved, s *Sessions) {
			s.Forge = &baseForgeFault{Forge: s.Forge, prErr: errors.New("the read failed")}
		},
		"a clone that lacks the objects": func(m *moved, s *Sessions) {
			empty := t.TempDir()
			gitOut(t, empty, "init", "-q")
			s.Clone = empty
		},
	} {
		m := newMoved(t, map[string]string{"a.txt": "task\n"})
		m.moveBase(map[string]string{"a.txt": "base\n"}) // an overlap: a redo, were the read not an error
		s := m.prodBot()
		tweak(m, s)
		head, n := recordsHead(t, m.w), m.nreqs()
		if got, err := m.bring(s); err == nil {
			t.Errorf("%s: BringUpToDate = %+v, nil; want an error", name, got)
		} else if strings.HasPrefix(name, "a default branch") && !strings.Contains(err.Error(), "is not an ancestor of the default branch") {
			t.Errorf("%s: %v; want the refusal that names the base", name, err)
		}
		if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
			t.Errorf("%s: the records moved (%v) or the forge was written to (%v); want neither", name, recordsHead(t, m.w) != head, m.writes(n))
		}
	}
}

// ----- the redo -----

// redoWorld is a moved world whose base overlaps the task, with a successor in the records, and the pull request of the attempt open.
func redoWorld(t *testing.T) *moved {
	t.Helper()
	m := newMoved(t, map[string]string{"a.txt": "task\n"})
	m.moveBase(map[string]string{"a.txt": "base\n"})
	tasks := m.tasks()
	next, err := AddTask(tasks, Task{ID: "T-sub1", Issue: 40, Kind: "build", Requirements: []string{"REQ-006"}, Size: "small", Predecessors: []string{"T-ab12"}, Opened: tasks[0].Opened})
	if err != nil {
		t.Fatal(err)
	}
	var b strings.Builder
	if err := WriteRegister(&b, next); err != nil {
		t.Fatal(err)
	}
	if err := m.r.commit(context.Background(), map[string][]byte{tasksPath: []byte(b.String())}, "a successor"); err != nil {
		t.Fatal(err)
	}
	m.pr()
	return m
}

// I6, the redo demo: the order of the forge's writes is the pull request, the issue, the new issue; one records commit holds the register,
// both events and the binding; the old row is replaced in place and its successor relinked.
func TestAnOverlappingMovedBaseRedoesTheTask(t *testing.T) {
	t.Parallel()
	m := redoWorld(t)
	n := m.nreqs()
	got, err := m.bring(m.s)
	if err != nil || got.Replaced == "" || got.Task != "T-ab12" || got.Verified {
		t.Fatalf("BringUpToDate = %+v, %v; want the task replaced", got, err)
	}
	fresh := got.Replaced
	if w := m.writes(n); !slices.Equal(w, []string{"PATCH /repos/acme/target/pulls/1", "PATCH /repos/acme/target/issues/3", "POST /repos/acme/target/issues"}) {
		t.Errorf("the writes of the forge = %v; want the pull request, the issue, the new issue, in that order", w)
	}
	tasks := m.tasks()
	if len(tasks) != 2 || tasks[0].ID != fresh || tasks[0].PR != 0 || tasks[0].Issue == 3 || tasks[0].Kind != "architecture" || tasks[1].ID != "T-sub1" || !slices.Equal(tasks[1].Predecessors, []string{fresh}) {
		t.Errorf("the register = %+v; want %s in place of T-ab12 with no pull request, and T-sub1 naming it", tasks, fresh)
	}
	if closed := m.kindsOf("T-ab12", "closed"); len(closed) != 1 || closed[0][2] != "1" || closed[0][6] != fresh {
		t.Errorf("the closed events = %v; want one, attempt 1, naming %s", closed, fresh)
	}
	base := strings.TrimSpace(gitOut(t, m.w.bare, "rev-parse", "refs/heads/main"))
	if rebased := m.kindsOf(fresh, "rebased"); len(rebased) != 1 || rebased[0][2] != "1" || rebased[0][4] != base || rebased[0][6] != "T-ab12" {
		t.Errorf("the rebased events = %v; want one: attempt 1, base %s, detail T-ab12", rebased, base)
	}
	// One records commit: the register, both events and the binding.
	commit := strings.TrimSpace(gitOut(t, m.w.bare, "log", "-1", "--format=%H", "layup-records", "--", "tasks/"+fresh+"/events.tsv"))
	changed := strings.Fields(gitOut(t, m.w.bare, "show", "--format=", "--name-only", commit))
	slices.Sort(changed)
	want := []string{"tasks.tsv", "tasks/T-ab12/events.tsv", "tasks/" + fresh + "/events.tsv", "tasks/" + fresh + "/options/db.tsv"}
	slices.Sort(want) // the new ID is drawn at random
	if !slices.Equal(changed, want) {
		t.Errorf("the commit that holds the rebased event changed %v; want %v", changed, want)
	}
	if intent := keptRecord(t, m.w, "tasks/T-ab12/redo.txt"); intent != fresh+" "+base+"\n" {
		t.Errorf("redo.txt = %q; want the new ID and its base", intent)
	}
	if st := m.w.fake.pulls.issueState[3]; st != "closed not_planned" || !m.w.fake.pullOf(1).closed {
		t.Errorf("the old issue %q and the old pull request closed %v; want closed as not planned, and closed", st, m.w.fake.pullOf(1).closed)
	}
	if title := m.w.fake.issueTexts[len(m.w.fake.issueTexts)-1][0]; !strings.Contains(title, fresh) {
		t.Errorf("the new issue is titled %q; want the new task's ID", title)
	}
	// A second call is a read: it names the replacement and writes nothing.
	head, n := recordsHead(t, m.w), m.nreqs()
	if again, err := m.bring(m.s); err != nil || again != got || recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("the second call = %+v, %v; writes %v; want the same answer and no write", again, err, m.writes(n))
	}
}

// Two paths that do not overlap can still conflict: the task adds a/b, the base adds the file a.
func TestAConflictingMergeOfDisjointPathsRedoesTheTask(t *testing.T) {
	t.Parallel()
	m := newMoved(t, map[string]string{"a/b": "task\n"})
	m.moveBase(map[string]string{"a": "base\n"})
	got, err := m.bring(m.s)
	if err != nil || got.Replaced == "" || len(m.kindsOf("T-ab12", "closed")) != 1 {
		t.Fatalf("BringUpToDate = %+v, %v; closed events %d; want a redo", got, err, len(m.kindsOf("T-ab12", "closed")))
	}
	if m.tip() != m.X {
		t.Errorf("the branch holds %s; want it untouched at %s", m.tip(), m.X)
	}
}

// A restart after each act of the redo completes it once, and a second call after that writes nothing. The window after OpenIssue
// opens a second issue; the records hold one redo.
func TestARestartAfterEachActOfTheRedoCompletesItOnce(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		stop   string
		issues int
	}{{"intent", 1}, {"ClosePullRequest", 1}, {"CloseIssue", 1}, {"OpenIssue", 2}} {
		m := redoWorld(t)
		s := m.prodBot()
		s.Forge = &baseForgeFault{Forge: s.Forge, after: c.stop}
		if c.stop == "intent" {
			s.Store = &baseStoreFault{SessionStore: s.Store, after: "tasks/T-ab12/redo.txt"}
		}
		before := len(m.w.fake.issueTexts)
		if _, err := m.bring(s); !errors.Is(err, errStop) {
			t.Fatalf("stop %s: %v; want the stop", c.stop, err)
		}
		if len(m.kindsOf("T-ab12", "closed")) != 0 {
			t.Fatalf("stop %s: a closed event before the combined commit", c.stop)
		}
		if c.stop == "intent" { // a decision that would now be a wait: the pull request is closed by a person; the intent goes on
			m.w.fake.pullOf(1).closed = true
		}
		var got Advance
		var err error
		for i := 0; i < 2; i++ { // the restart, then a second call
			if got, err = m.bring(m.prodBot()); err != nil || got.Replaced == "" {
				t.Fatalf("stop %s, restart %d: %+v, %v; want the task replaced", c.stop, i+1, got, err)
			}
		}
		// Cleanup is done before the first attempt of the new task; the driver then starts exactly one, at the rebased base.
		if st := m.w.fake.pulls.issueState[3]; st != "closed not_planned" || !m.w.fake.pullOf(1).closed || len(m.kindsOf(got.Replaced, "attempt")) != 0 {
			t.Fatalf("stop %s: before the first attempt the old issue is %q, the old pull request closed %v, attempts %d; want cleanup done and no attempt", c.stop, st, m.w.fake.pullOf(1).closed, len(m.kindsOf(got.Replaced, "attempt")))
		}
		for i := 0; i < 2; i++ {
			if adv, err := m.s2().AdvanceShapeTask(context.Background(), m.driverIn(got.Replaced)); err != nil || !strings.Contains(adv.Stop, "no admitted pair") {
				t.Fatalf("stop %s, driver entry %d on the new task: %+v, %v; want the stop for the pair after the attempt", c.stop, i+1, adv, err)
			}
		}
		if att := m.kindsOf(got.Replaced, "attempt"); len(att) != 1 || att[0][4] != strings.TrimSpace(gitOut(t, m.w.bare, "rev-parse", "refs/heads/main")) || len(m.kindsOf(got.Replaced, "rebased")) != 1 {
			t.Errorf("stop %s: attempts %v, rebased %d; want one attempt at the rebased base and one rebased", c.stop, att, len(m.kindsOf(got.Replaced, "rebased")))
		}
		tasks := m.tasks()
		if len(tasks) != 2 || tasks[0].ID != got.Replaced || len(m.kindsOf("T-ab12", "closed")) != 1 || len(m.kindsOf(got.Replaced, "rebased")) != 1 {
			t.Errorf("stop %s: the register %+v, closed %d, rebased %d; want one replacement", c.stop, tasks, len(m.kindsOf("T-ab12", "closed")), len(m.kindsOf(got.Replaced, "rebased")))
		}
		if opened := len(m.w.fake.issueTexts) - before; opened != c.issues {
			t.Errorf("stop %s: %d issues opened; want %d", c.stop, opened, c.issues)
		}
		if st := m.w.fake.pulls.issueState[3]; st != "closed not_planned" || !m.w.fake.pullOf(1).closed {
			t.Errorf("stop %s: the old issue %q, the old pull request closed %v", c.stop, st, m.w.fake.pullOf(1).closed)
		}
		if tasks[0].Issue == 3 {
			t.Errorf("stop %s: the new row names the old issue 3", c.stop)
		}
	}
}

// A pull request merged before the redo is a stop (the merge is recorded, not redone), and one merged as the redo closes it is an error;
// in neither is the old task closed.
func TestAPullRequestMergedAfterItsRedoStops(t *testing.T) {
	t.Parallel()
	m := redoWorld(t)
	m.w.fake.pullOf(1).merged = true
	n := m.nreqs()
	if got, err := m.bring(m.s); err != nil || !strings.Contains(got.Stop, "the pull request 1 is merged") || got.Replaced != "" {
		t.Errorf("a merged pull request: %+v, %v; want a stop that says merged", got, err)
	}
	if len(m.writes(n)) != 0 || len(m.kindsOf("T-ab12", "closed")) != 0 {
		t.Errorf("a merged pull request: the writes %v, closed events %d; want none", m.writes(n), len(m.kindsOf("T-ab12", "closed")))
	}
	m = redoWorld(t)
	s := m.prodBot()
	s.Forge = &baseForgeFault{Forge: s.Forge, closeMerged: true}
	n = m.nreqs()
	if got, err := m.bring(s); err == nil || !strings.Contains(err.Error(), "merged after its redo") {
		t.Errorf("a pull request that merged as it was closed: %+v, %v; want \"merged after its redo\"", got, err)
	}
	for _, w := range m.writes(n) {
		if strings.Contains(w, "/issues") {
			t.Errorf("the write %s after a merged pull request; want the issue left open", w)
		}
	}
	if len(m.kindsOf("T-ab12", "closed")) != 0 {
		t.Error("a closed event after a pull request that merged")
	}
}

// A pull request that a person closed, with no redo decided, is a wait and writes nothing.
func TestAClosedUnmergedPullRequestWithNoRedoIsAWait(t *testing.T) {
	t.Parallel()
	m := redoWorld(t)
	m.w.fake.pullOf(1).closed = true
	head, n := recordsHead(t, m.w), m.nreqs()
	got, err := m.bring(m.s)
	if err != nil || got.Wait != "waiting for the Operator: T-ab12: the pull request 1 is closed and not merged" || got.Replaced != "" {
		t.Errorf("BringUpToDate = %+v, %v; want the wait for the Operator", got, err)
	}
	if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("a wait wrote to the records (%v) or the forge (%v)", recordsHead(t, m.w) != head, m.writes(n))
	}
}

// A lost lease stops the redo before any act of the forge, and one lost in the middle stops before the next.
func TestALostLeaseOpensNoIssueAndWritesNothingAfterIt(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name     string
		lostFrom int
		patches  int
	}{{"lost at the start", 1, 0}, {"lost after the pull request", 2, 1}, {"lost after the issue", 3, 2}} {
		m := redoWorld(t)
		s := m.prodBot()
		s.Store = &baseStoreFault{SessionStore: s.Store, lostFrom: c.lostFrom}
		n := m.nreqs()
		_, err := m.bring(s)
		var lost *LostError
		if !errors.As(err, &lost) {
			t.Errorf("%s: %v; want a *LostError", c.name, err)
		}
		w := m.writes(n)
		patches := 0
		for _, r := range w {
			if strings.HasPrefix(r, "PATCH") {
				patches++
			}
			if strings.HasPrefix(r, "POST") {
				t.Errorf("%s: %s; want no new issue", c.name, r)
			}
		}
		if patches != c.patches {
			t.Errorf("%s: %d PATCH requests (%v); want %d", c.name, patches, w, c.patches)
		}
	}
	// The carry: a lost lease pushes nothing.
	m := newMoved(t, taskPaths)
	m.moveBase(basePaths)
	s := m.prodBot()
	s.Store = &baseStoreFault{SessionStore: s.Store, lostFrom: 1}
	var lost *LostError
	if _, err := m.bring(s); !errors.As(err, &lost) || m.tip() != m.X {
		t.Errorf("a lost lease before the push: %v, the branch at %s; want a *LostError and the branch at %s", err, m.tip(), m.X)
	}
}

// A task that holds the rendered commit is a stop in M2e (the Operator's decision on #424): nothing is closed or written.
func TestARedoThatWouldDropTheRenderedCommitIsAStop(t *testing.T) {
	t.Parallel()
	m := redoWorld(t)
	clone := m.r.Clone
	gitOut(t, clone, "switch", "-q", "--detach", m.X)
	os.WriteFile(filepath.Join(clone, "r.txt"), []byte("rendered\n"), 0o644)
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "the rendered commit")
	rendered := strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))
	gitOut(t, clone, "switch", "-q", "main")
	gitOut(t, clone, "push", "-q", "file://"+m.w.bare, rendered+":refs/heads/task/T-ab12/1")
	for _, kind := range []string{"push", "bound"} { // the events of code: session —
		data, err := m.r.appendEvents(context.Background(), "T-ab12", 1, "", []string{kind + " " + rendered + " task/T-ab12/1"})
		if err != nil {
			t.Fatal(err)
		}
		if err := m.r.commit(context.Background(), map[string][]byte{eventsPath("T-ab12"): data}, "the rendered commit "+kind); err != nil {
			t.Fatal(err)
		}
	}
	head, n := recordsHead(t, m.w), m.nreqs()
	got, err := m.bring(m.prodBot())
	if err != nil || got.Stop != "waiting for the Operator: T-ab12: a redo would drop the rendered commit" {
		t.Errorf("BringUpToDate = %+v, %v; want the stop for the rendered commit", got, err)
	}
	if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("a stop wrote to the records (%v) or the forge (%v)", recordsHead(t, m.w) != head, m.writes(n))
	}
}

// ----- the older attempt's pull request -----

// The pull request of an attempt that a newer attempt replaced is closed once, and only when the register names the newer one.
func TestTheOlderAttemptsPullRequestIsClosedOnce(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	m.pr()
	ctx := context.Background()
	if n, err := m.r.StartAttempt(ctx, "T-ab12", m.B); err != nil || n != 2 {
		t.Fatalf("StartAttempt = %d, %v", n, err)
	}
	// The control: attempt 2 exists and has no pull request; attempt 1's stays open, with no write.
	n := m.nreqs()
	if _, err := m.s.RepairTask(ctx, "T-ab12"); err != nil || len(m.writes(n)) != 0 || m.w.fake.pullOf(1).closed {
		t.Fatalf("RepairTask before the new pull request: %v, writes %v, closed %v; want no close", err, m.writes(n), m.w.fake.pullOf(1).closed)
	}
	id := m.as(roleOfStep("implement"), 2, "claude", map[string]string{"t2.txt": "again\n"}, okResult)
	if m.bound(id) == "" {
		t.Fatal("the developer of attempt 2 is not bound")
	}
	// An interrupted registration: the pull request of attempt 2 is open on the forge and the register still names attempt 1's.
	if _, err := m.s.Forge.OpenDraft(ctx, "task/T-ab12/2", "main", "feat: the second try", "the body"); err != nil {
		t.Fatal(err)
	}
	n = m.nreqs()
	if _, err := m.s.RepairTask(ctx, "T-ab12"); err != nil || len(m.writes(n)) != 0 || m.w.fake.pullOf(1).closed {
		t.Fatalf("RepairTask with the register still naming attempt 1: %v, writes %v, closed %v; want no close", err, m.writes(n), m.w.fake.pullOf(1).closed)
	}
	if _, err := m.s.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 2, Type: "feat", Description: "the second try"}); err != nil {
		t.Fatal(err)
	}
	n = m.nreqs()
	if _, err := m.s.RepairTask(ctx, "T-ab12"); err != nil || !slices.Equal(m.writes(n), []string{"PATCH /repos/acme/target/pulls/1"}) || !m.w.fake.pullOf(1).closed {
		t.Fatalf("RepairTask: %v, writes %v; want one PATCH of pull request 1", err, m.writes(n))
	}
	n = m.nreqs()
	if _, err := m.s.RepairTask(ctx, "T-ab12"); err != nil || len(m.writes(n)) != 0 {
		t.Errorf("the second RepairTask wrote %v; want no second PATCH", m.writes(n))
	}
}

// ----- the push of the merge -----

// A merge whose push was made and whose bound was lost is adopted at the restart: the same merge, no second push event, the status once.
func TestAPushedMergeCommitIsAdoptedOnARestart(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	m.moveBase(basePaths)
	s := m.prodBot()
	s.Store = &baseStoreFault{SessionStore: s.Store, failIf: func(files map[string][]byte) bool {
		e, ok := files[eventsPath("T-ab12")]
		return ok && strings.Contains(string(e), "\tbound\t") && strings.Count(string(e), "\tbound\t") > 1
	}}
	if _, err := m.bring(s); !errors.Is(err, errStop) {
		t.Fatalf("BringUpToDate with the bound lost: %v; want the stop", err)
	}
	merge := m.tip()
	if merge == m.X || len(statusesOf(m.w, merge)) != 0 {
		t.Fatalf("the branch holds %s and the merge has %d statuses; want the push made and no status before the bound", merge, len(statusesOf(m.w, merge)))
	}
	got, err := m.bring(m.prodBot())
	if err != nil || got.Head != merge || !got.Verified {
		t.Fatalf("the restart = %+v, %v; want the same merge %s", got, err, merge)
	}
	pushes := 0
	for _, e := range m.eventsOfTask("T-ab12") {
		if e[1] == "push" && e[3] == "" && e[5] == merge {
			pushes++
		}
	}
	if pushes != 1 || len(statusesOf(m.w, merge)) != 1 || m.tip() != merge {
		t.Errorf("push events of the merge %d, statuses %d, the branch %s; want one, one and %s", pushes, len(statusesOf(m.w, merge)), m.tip(), merge)
	}
}

// A merge that was announced and not pushed, when the default branch moves before the restart: a new merge with its own announcement.
func TestAMergeAnnouncedAndNotPushedIsMadeAgainWhenTheBaseMovesAgain(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	m.moveBase(basePaths)
	s := m.prodBot()
	s.Store = &baseStoreFault{SessionStore: s.Store, pushErr: errStop}
	if _, err := m.bring(s); !errors.Is(err, errStop) || m.tip() != m.X {
		t.Fatalf("BringUpToDate with the push lost: %v, the branch at %s; want the stop and the branch at %s", err, m.tip(), m.X)
	}
	m.moveBase(map[string]string{"v.txt": "again\n"})
	got, err := m.bring(m.prodBot())
	if err != nil || !got.Verified || m.tip() != got.Head {
		t.Fatalf("the restart = %+v, %v, the branch %s; want the new merge pushed", got, err, m.tip())
	}
	var announced []string
	for _, e := range m.eventsOfTask("T-ab12") {
		if e[1] == "push" && e[3] == "" {
			announced = append(announced, e[5])
		}
	}
	if len(announced) != 2 || announced[1] != got.Head || announced[0] == announced[1] {
		t.Errorf("the announced pushes %v; want two merges, the second being %s", announced, got.Head)
	}
}

// A task whose close-out is its head: the merge is made when the base leaves the close-out's two files, and the task is redone when it
// does not (the completed log is one of them).
func TestTheCloseOutHeadIsMergedWhenTheBaseLeavesItsFilesAndRedoneWhenItDoesNot(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name  string
		files map[string]string
		redo  bool
	}{{"the base leaves the two files", basePaths, false}, {"the base appends to the completed log", map[string]string{completedFile: "log\nanother\n"}, true}} {
		m := newMoved(t, taskPaths)
		ctx := context.Background()
		gitOut(t, m.r.Clone, "branch", "-f", "layup/base/T-ab12/c", m.X)
		spec := harness(t, m.a, closing())
		spec.Role, spec.Base, spec.Branch, spec.Pair.Harness = roleOfStep("close-out"), m.X, "layup/base/T-ab12/c", "claude"
		spec.Push = func(id string, d session.Dir, r session.Run) error { return m.r.PushCloseOut(ctx, id, m.in) }
		id, err := m.r.TaskSession(ctx, spec)
		if err != nil || m.bound(id) == "" {
			t.Fatalf("%s: the close-out: %q, %v, bound %q", c.name, id, err, m.bound(id))
		}
		closeOut := m.bound(id)
		m.moveBase(c.files)
		got, err := m.bring(m.prodBot())
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if c.redo != (got.Replaced != "") || (!c.redo && (!got.Verified || got.Head == closeOut)) {
			t.Errorf("%s: %+v; want redo %v", c.name, got, c.redo)
		}
		if !c.redo {
			if parents, err := git.Parents(m.r.Clone, got.Head); err != nil || parents[0] != closeOut {
				t.Errorf("%s: the parents of %s = %v, %v; want the close-out %s first", c.name, got.Head, parents, err, closeOut)
			}
		}
	}
}

// A merge whose status could not be posted is posted at the next call: the head is current then, and the carry is still undone.
func TestACarryWhoseStatusFailedIsPostedAtTheRestart(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	m.moveBase(basePaths)
	m.w.fake.failStatus = 1
	if got, err := m.bring(m.s); err == nil {
		t.Fatalf("BringUpToDate with the status refused = %+v, nil; want an error", got)
	}
	merge := m.tip()
	if merge == m.X || len(statusesOf(m.w, merge)) != 0 {
		t.Fatalf("the branch holds %s and the merge has %d statuses; want the merge pushed and no status", merge, len(statusesOf(m.w, merge)))
	}
	if got, err := m.bring(m.prodBot()); err != nil || got.Head != merge || len(statusesOf(m.w, merge)) != 1 {
		t.Errorf("the restart = %+v, %v with %d statuses; want the merge %s carried once", got, err, len(statusesOf(m.w, merge)), merge)
	}
}

// The new issue is composed from the inputs the records keep: the decision and the option of the binding, or the draft as the first session
// of the old task read it, never a later one; a kind with no composer is refused before any write.
func TestTheNewIssueIsComposedFromTheInputsTheRecordsKeep(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		name    string
		kind    string
		binding bool
		want    []string
	}{
		{"an architecture task of a decision", "architecture", true, []string{"the architecture of the decision db", "Chosen option: db, recorded in tasks/", "Requirement IDs: REQ-005"}},
		{"a Shape architecture task", "architecture", false, []string{"the target's architecture", "SHA-256 " + digest([]byte("the draft of Intake\n")), "Requirement IDs: REQ-005"}},
		{"a specification task", "specification", false, []string{"the specification", "Pinned draft: spec/draft.tsv, SHA-256 " + digest([]byte("the draft of Intake\n"))}},
	} {
		m := newMovedAs(t, map[string]string{"a.txt": "task\n"}, c.kind, c.binding)
		m.moveBase(map[string]string{"a.txt": "base\n"})
		if !c.binding { // a later draft changes nothing of the issue
			if err := m.r.commit(context.Background(), map[string][]byte{"spec/draft.tsv": []byte("a later draft\n")}, "a later draft"); err != nil {
				t.Fatal(err)
			}
		}
		got, err := m.bring(m.s)
		if err != nil || got.Replaced == "" {
			t.Fatalf("%s: %+v, %v; want a redo", c.name, got, err)
		}
		text := m.w.fake.issueTexts[len(m.w.fake.issueTexts)-1]
		for _, w := range append(c.want, "LAYUP task "+got.Replaced) {
			if !strings.Contains(text[0]+"\n"+text[1], w) {
				t.Errorf("%s: the new issue %q lacks %q", c.name, text, w)
			}
		}
		if strings.Contains(text[1], "T-ab12") {
			t.Errorf("%s: the new issue names the old task:\n%s", c.name, text[1])
		}
	}
	m := newMovedAs(t, map[string]string{"a.txt": "task\n"}, "build", true)
	m.moveBase(map[string]string{"a.txt": "base\n"})
	head, n := recordsHead(t, m.w), m.nreqs()
	if got, err := m.bring(m.s); err == nil || !strings.Contains(err.Error(), "no composer of an issue") || recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("a build task: %+v, %v; want the refusal before any write", got, err)
	}
}

// The first attempt of a redone task starts at the base of its event rebased (the driver of a Shape task).
func TestTheDriverStartsTheFirstAttemptOfARedoneTaskAtTheRebasedBase(t *testing.T) {
	t.Parallel()
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	in := d.input(spec)
	in.Harnesses = nil // the first act is the attempt; the next is the stop for the pair
	rebased := strings.Repeat("c", 40)
	s := d.r.sessions()
	data, err := s.appendEvent(context.Background(), spec, "rebased", 1, "", rebased, "T-zzzz")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.commit(context.Background(), map[string][]byte{eventsPath(spec): data}, "a redo"); err != nil {
		t.Fatal(err)
	}
	if got, err := s.AdvanceShapeTask(context.Background(), in); err != nil || !strings.Contains(got.Stop, "no admitted pair") {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the stop for the pair after the attempt", got, err)
	}
	var base string
	for _, e := range d.eventsOf(spec) {
		if e[1] == "attempt" {
			base = e[4]
		}
	}
	if base != rebased {
		t.Errorf("the base of the first attempt is %q; want the base of the event rebased, %s", base, rebased)
	}
}

// A head that is not the verified one cannot carry the verdict: the merge is refused before anything is written or pushed.
func TestAHeadThatDoesNotCarryTheVerdictIsNotMerged(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	clone := m.r.Clone
	gitOut(t, clone, "switch", "-q", "--detach", m.X)
	os.WriteFile(filepath.Join(clone, "t2.txt"), []byte("after the verdict\n"), 0o644)
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "after the verdict")
	later := strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))
	gitOut(t, clone, "switch", "-q", "main")
	gitOut(t, clone, "push", "-q", "file://"+m.w.bare, later+":refs/heads/task/T-ab12/1")
	for _, kind := range []string{"push", "bound"} {
		data, err := m.r.appendEvents(context.Background(), "T-ab12", 1, "", []string{kind + " " + later + " task/T-ab12/1"})
		if err != nil {
			t.Fatal(err)
		}
		if err := m.r.commit(context.Background(), map[string][]byte{eventsPath("T-ab12"): data}, "a later head "+kind); err != nil {
			t.Fatal(err)
		}
	}
	m.moveBase(basePaths)
	head, n := recordsHead(t, m.w), m.nreqs()
	if got, err := m.bring(m.s); err == nil {
		t.Errorf("BringUpToDate = %+v, nil; want an error", got)
	}
	if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 || m.tip() != later {
		t.Errorf("the records moved (%v), the forge was written (%v) or the branch moved to %s; want none", recordsHead(t, m.w) != head, m.writes(n), m.tip())
	}
}

// s2 gives fresh Sessions of the production constructor; driverIn is what the shared driver is given for the task.
func (m *moved) s2() *Sessions { return m.prodBot() }

func (m *moved) driverIn(task string) ShapeTaskInput {
	st := m.s.Store.(*recordsStore)
	return ShapeTaskInput{Task: task, Base: m.B, Records: st.Base}
}

func lostBound(files map[string][]byte) bool {
	e, ok := files[eventsPath("T-ab12")]
	return ok && strings.Count(string(e), "\tbound\t") > 1
}

// I1 of round 1 (late test): the push of M1 was made and its bound lost; the default moves again before the restart. The restart adopts
// M1 (its bound, its status), then merges D2 into it, twice over with no second write.
func TestAnInterruptedCarryIsAdoptedWhenTheDefaultMovesAgain(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	m.moveBase(basePaths)
	s := m.prodBot()
	s.Store = &baseStoreFault{SessionStore: s.Store, failIf: lostBound}
	if _, err := m.bring(s); !errors.Is(err, errStop) {
		t.Fatalf("BringUpToDate with the bound lost: %v; want the stop", err)
	}
	m1 := m.tip()
	d2 := m.moveBase(map[string]string{"v.txt": "again\n"})
	var got Advance
	for i := 0; i < 2; i++ {
		var err error
		if got, err = m.bring(m.prodBot()); err != nil || !got.Verified || got.Head == m1 {
			t.Fatalf("restart %d = %+v, %v; want a merge of D2 over the adopted %s", i+1, got, err, m1)
		}
	}
	if parents, err := git.Parents(m.r.Clone, got.Head); err != nil || !slices.Equal(parents, []string{m1, d2}) || m.tip() != got.Head {
		t.Errorf("the parents of %s = %v, %v; the branch %s; want [%s %s] and the branch at the merge", got.Head, parents, err, m.tip(), m1, d2)
	}
	for _, sha := range []string{m1, got.Head} {
		n := 0
		for _, e := range m.eventsOfTask("T-ab12") {
			if e[1] == "bound" && e[3] == "" && e[5] == sha {
				n++
			}
		}
		if n != 1 || len(statusesOf(m.w, sha)) != 1 {
			t.Errorf("%s: %d bound events and %d statuses; want one and one", sha, n, len(statusesOf(m.w, sha)))
		}
	}
	head, n := recordsHead(t, m.w), m.nreqs()
	if _, err := m.bring(m.prodBot()); err != nil || recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("a third call: %v; writes %v; want none", err, m.writes(n))
	}
}

// I2 of round 1 (late test): the pull request opened before the carry is adopted at the merge, again and again, with no draft and no
// register change; a head that moved with no explanation is still refused.
func TestOpenPullRequestAdoptsTheSamePullRequestAfterACarry(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	pr := m.pr()
	m.moveBase(basePaths)
	got, err := m.bring(m.s)
	if err != nil || !got.Verified || got.Head == m.X {
		t.Fatalf("BringUpToDate = %+v, %v; want the carry", got, err)
	}
	m.w.fake.movePull(pr.Number, got.Head) // the push of the merge moved the head of the pull request
	head, n := recordsHead(t, m.w), m.nreqs()
	for i := 0; i < 3; i++ {
		again, err := m.s.OpenPullRequest(context.Background(), PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "the work of the task"})
		if err != nil || again.Number != pr.Number || again.HeadSHA != got.Head {
			t.Fatalf("OpenPullRequest %d = %+v, %v; want pull request %d at %s", i+1, again, err, pr.Number, got.Head)
		}
	}
	if recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("the register or the forge changed: records moved %v, writes %v; want neither", recordsHead(t, m.w) != head, m.writes(n))
	}
	m.w.fake.movePull(pr.Number, strings.Repeat("e", 40))
	if bad, err := m.s.OpenPullRequest(context.Background(), PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "the work of the task"}); err == nil {
		t.Errorf("a head that moved with no explanation: %+v, nil; want an error", bad)
	}
}

// I3 of round 1 (late test): the pin of the draft of a Shape task is the one of the original issue, through any chain of redos.
func TestASecondRedoOfAShapeTaskKeepsTheOriginalPin(t *testing.T) {
	t.Parallel()
	m := newMovedAs(t, map[string]string{"a.txt": "task\n"}, "architecture", false)
	ctx := context.Background()
	pin := ""
	{
		rows, err := records.ReadSessions([]byte(keptRecord(t, m.w, "sessions.tsv")))
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range rows {
			if r[sessionTask] == "T-ab12" {
				pin = r[sessionRecords]
				break
			}
		}
	}
	original := "SHA-256 " + digest([]byte("the draft of Intake\n")) + ", read at the records commit " + pin
	m.moveBase(map[string]string{"a.txt": "base\n"})
	first, err := m.bring(m.s)
	if err != nil || first.Replaced == "" {
		t.Fatalf("the first redo: %+v, %v", first, err)
	}
	if text := m.w.fake.issueTexts[len(m.w.fake.issueTexts)-1][1]; !strings.Contains(text, original) {
		t.Errorf("the first replacement issue lacks the original pin %q:\n%s", original, text)
	}
	// A later draft, then the first session of the replacement reads it.
	if err := m.r.commit(ctx, map[string][]byte{"spec/draft.tsv": []byte("a later draft\n")}, "a later draft"); err != nil {
		t.Fatal(err)
	}
	rows, err := records.ReadSessions([]byte(keptRecord(t, m.w, "sessions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	row := sessionRow("S-d0000009", first.Replaced, "plan-of-a-task", "claude")
	row[sessionRecords] = m.r.Store.(*recordsStore).pushed
	data, err := table(records.SessionsSchema, append(rows, row))
	if err != nil {
		t.Fatal(err)
	}
	if err := m.r.commit(ctx, map[string][]byte{"sessions.tsv": data}, "a session of the replacement"); err != nil {
		t.Fatal(err)
	}
	rowB, err := m.s.taskRow(ctx, first.Replaced)
	if err != nil {
		t.Fatal(err)
	}
	_, body, err := m.s.redoIssue(ctx, rowB, "T-zzzz", nil)
	if err != nil || !strings.Contains(body, original) {
		t.Errorf("the second replacement issue = %q, %v; want the original pin %q", body, err, original)
	}
}

// I4a of round 1 (a characterization on the code of the head, with its mutant): the carry over a rendered first parent.
func TestTheCarryGoesOverTheRenderedFirstParent(t *testing.T) {
	t.Parallel()
	m := newMoved(t, taskPaths)
	ctx := context.Background()
	clone := m.r.Clone
	gitOut(t, clone, "switch", "-q", "--detach", m.X)
	for p, d := range map[string]string{treePRDPath: renderedPRD(t, ""), treeFacts: treeFactsAt} {
		os.MkdirAll(filepath.Dir(filepath.Join(clone, p)), 0o755)
		os.WriteFile(filepath.Join(clone, p), []byte(d), 0o644)
	}
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "the rendered commit")
	rendered := strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))
	gitOut(t, clone, "switch", "-q", "main")
	gitOut(t, clone, "push", "-q", "file://"+m.w.bare, rendered+":refs/heads/task/T-ab12/1")
	for _, kind := range []string{"push", "bound"} {
		data, err := m.r.appendEvents(ctx, "T-ab12", 1, "", []string{kind + " " + rendered + " task/T-ab12/1"})
		if err != nil {
			t.Fatal(err)
		}
		if err := m.r.commit(ctx, map[string][]byte{eventsPath("T-ab12"): data}, "the rendered commit "+kind); err != nil {
			t.Fatal(err)
		}
	}
	m.in.Rendered = &RenderedInput{Head: rendered, Render: treeRecs}
	d := m.moveBase(basePaths)
	got, err := m.bring(m.prodBot())
	if err != nil || !got.Verified || got.Head == rendered {
		t.Fatalf("BringUpToDate = %+v, %v; want a merge over the rendered commit", got, err)
	}
	if parents, err := git.Parents(clone, got.Head); err != nil || !slices.Equal(parents, []string{rendered, d}) || m.tip() != got.Head {
		t.Errorf("the parents of %s = %v, %v, the branch %s; want [%s %s] and the branch at the merge", got.Head, parents, err, m.tip(), rendered, d)
	}
	pushes, bounds := 0, 0
	for _, e := range m.eventsOfTask("T-ab12") {
		if e[3] == "" && e[5] == got.Head {
			if e[1] == "push" {
				pushes++
			} else if e[1] == "bound" {
				bounds++
			}
		}
	}
	if st := statusesOf(m.w, got.Head); pushes != 1 || bounds != 1 || len(st) != 1 || st[0]["description"] != "verified head "+m.X+" plus a clean merge of the base" {
		t.Errorf("push events %d, bound events %d, statuses %v; want one, one and the carry of %s", pushes, bounds, st, m.X)
	}
	head, n := recordsHead(t, m.w), m.nreqs()
	if again, err := m.bring(m.prodBot()); err != nil || again != got || recordsHead(t, m.w) != head || len(m.writes(n)) != 0 {
		t.Errorf("the repeat call = %+v, %v; writes %v; want the same answer and no write", again, err, m.writes(n))
	}
}
