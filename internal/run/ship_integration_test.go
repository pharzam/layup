//go:build integration

package run

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// A ship is the world of the close-out and the merge (row 66g, task T-xvvd, #363): the target with the task files and the
// product requirement document before the bet in its first tree, the verified head (the base, reviewed by a verifier of devin,
// an independent harness), and the registers that decide the verifier's independence.
type ship struct {
	*loopRun
	V  string // the verified head
	in CarryInput
}

const (
	taskFile      = "docs/tasks/T-ab12.md"
	completedFile = "docs/tasks/completed.md"
)

var shipBaseline = map[string]string{taskFile: "task\n", completedFile: "log\n", treePRDPath: treePRDPre}

func newShip(t *testing.T) *ship {
	t.Helper()
	l := newLoop(t, shipBaseline)
	l.r.Bot = 9001
	harnesses, models := commitRegisters(t, l)
	id := l.comment(1, "verification", "devin", record(l.a.base, "nothing material in scope"))
	return &ship{loopRun: l, V: l.a.base, in: CarryInput{Task: "T-ab12", Verifier: id, Harnesses: harnesses, Models: models}}
}

// verifierOf runs a verifier of devin at the head, whose record is the verdict; its session is the one the carry-over names.
func (s *ship) verifierOf(head, verdict string, bullets ...string) string {
	s.t.Helper()
	return s.comment(1, "verification", "devin", record(head, verdict, bullets...))
}

// prodBot gives the Sessions from the constructor that a run uses, with the App's bot as the run learns it.
func (l *loopRun) prodBot() *Sessions {
	l.t.Helper()
	st := &state{cfg: l.w.config("0.1.0-dev"), store: l.r.Store.(*recordsStore), target: l.r.Clone, defaultBranch: l.r.Branch, botID: 9001,
		start: [][]string{{"issue.control", "2", "run"}}}
	s := st.sessions()
	s.Now, s.RunID = l.r.Now, l.r.RunID
	return s
}

// closeOut runs a close-out session on the base, whose push hook is the close-out check of in, and gives the session and its
// error. shell is the work of the session in repo/ before its result.
func (s *ship) closeOut(base string, in CarryInput, work string) (string, error) {
	s.t.Helper()
	spec := harness(s.t, s.a, work)
	spec.Role, spec.Base, spec.Pair.Harness = roleOfStep("close-out"), base, "claude"
	spec.Push = func(id string, d session.Dir, r session.Run) error {
		return s.r.PushCloseOut(context.Background(), id, in)
	}
	return s.r.TaskSession(context.Background(), spec)
}

// closing is the work of a valid close-out: the two files changed in one commit, and the result.
func closing() string {
	return fileWork(map[string]string{taskFile: "task\nverdict\n", completedFile: "log\nline\n"}, okResult)
}

// refusalOf gives the detail of the event refused of the session, "" for none.
func (s *ship) refusalOf(id string) string {
	for _, e := range s.events() {
		if e[1] == "refused" && e[3] == id {
			return e[6]
		}
	}
	return ""
}

// branchTip gives the commit of the task branch on the forge, which the verifier's own push left at the verified head.
func (s *ship) branchTip() string {
	s.t.Helper()
	return strings.TrimSpace(gitOut(s.t, s.w.bare, "rev-parse", "refs/heads/task/T-ab12/1"))
}

func TestTheCloseOutCommitCarriesTheVerdict(t *testing.T) {
	s := newShip(t)
	ctx := context.Background()
	id, err := s.closeOut(s.V, s.in, closing())
	if err != nil {
		t.Fatalf("the close-out session: %v", err)
	}
	head := s.bound(id)
	if head == "" || s.refusalOf(id) != "" || s.branchTip() != head {
		t.Fatalf("bound %q, refused %q, the branch at %s; want the close-out pushed and bound", head, s.refusalOf(id), s.branchTip())
	}
	got, err := s.r.CarryVerify(ctx, s.in, id)
	if err != nil || got != head {
		t.Fatalf("CarryVerify = %q, %v; want the bound head %s", got, err, head)
	}
	posted := statusesOf(s.w, head)
	if len(posted) != 1 || posted[0]["context"] != VerifyContext || posted[0]["state"] != "success" || posted[0]["description"] != "verified head "+s.V+" plus the close-out" {
		t.Errorf("the statuses of the new head: %v; want layup/verify success, verified head %s plus the close-out", posted, s.V)
	}
	if _, err := s.r.CarryVerify(ctx, s.in, id); err != nil || len(statusesOf(s.w, head)) != 1 {
		t.Errorf("the second CarryVerify: %v with %d statuses; want one status in all", err, len(statusesOf(s.w, head)))
	}
}

// Every refusal is written as the others are: the event refused with its reason and the SHA-256 of the diff after a space, the
// diff in payloads/, nothing pushed and nothing bound.
func TestACloseOutThatChangesMoreThanTheTwoFilesIsRefusedAsScope(t *testing.T) {
	two := "docs/tasks/T-ab12.md"
	for _, c := range []struct{ name, work string }{
		{"another file", fileWork(map[string]string{taskFile: "t\n", "x.go": "package x\n"}, okResult)},
		{"a task file of another task", fileWork(map[string]string{"docs/tasks/T-cd34.md": "t\n"}, okResult)},
		{"the backlog", fileWork(map[string]string{taskFile: "t\n", "docs/tasks/backlog.md": "b\n"}, okResult)},
		{"two commits", "printf a > docs/tasks/completed.md && git add -A && git commit -qm first && " + closing()},
		{"an executable file", "printf t > " + two + " && chmod +x " + two + " && git add -A && git commit -qm work && printf '" + octal(okResult) + "' > ../result/result.tsv"},
		{"a link", "rm " + completedFile + " && ln -s x " + completedFile + " && git add -A && git commit -qm work && printf '" + octal(okResult) + "' > ../result/result.tsv"},
		{"a deletion", "git rm -q " + completedFile + " && git commit -qm work && printf '" + octal(okResult) + "' > ../result/result.tsv"},
		{"a commit that changes nothing, whose diff is empty", "git commit -q --allow-empty -m work && printf '" + octal(okResult) + "' > ../result/result.tsv"},
	} {
		s := newShip(t)
		tip := s.branchTip()
		id, err := s.closeOut(s.V, s.in, c.work)
		if err != nil {
			t.Errorf("%s: the session: %v", c.name, err)
			continue
		}
		reason := s.refusalOf(id)
		sum, ok := strings.CutPrefix(reason, "scope ")
		if !ok || len(sum) != 64 {
			t.Errorf("%s: refused %q; want scope and the SHA-256 of the diff", c.name, reason)
			continue
		}
		empty := sum == "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // the SHA-256 of no diff at all
		if (keptRecord(t, s.w, "payloads/"+sum) == "") == !empty || s.branchTip() != tip || s.bound(id) != "" {
			t.Errorf("%s: payload %q, the branch at %s (was %s), bound %q; want the diff kept and nothing pushed or bound", c.name, keptRecord(t, s.w, "payloads/"+sum), s.branchTip(), tip, s.bound(id))
		}
	}
}

func TestACloseOutOnAParentThatDoesNotCarryTheVerdictIsRefusedAsVerdict(t *testing.T) {
	// An ordinary child of the verified head: the close-out sits on it, and only the verified head carries.
	s := newShip(t)
	child := commitIn(s.loopRun, map[string]string{"x.go": "package x\n"})
	tip := s.branchTip()
	id, err := s.closeOut(child, s.in, closing())
	if err != nil {
		t.Fatal(err)
	}
	if reason := s.refusalOf(id); !strings.HasPrefix(reason, "verdict ") || s.branchTip() != tip {
		t.Errorf("a child of the verified head: refused %q, the branch at %s (was %s); want verdict and nothing pushed", reason, s.branchTip(), tip)
	}
	// A verifier whose record has a finding: the verified head is none.
	s = newShip(t)
	s.in.Verifier = s.verifierOf(s.V, "material", located)
	tip = s.branchTip()
	id, err = s.closeOut(s.V, s.in, closing())
	if err != nil {
		t.Fatal(err)
	}
	if reason := s.refusalOf(id); !strings.HasPrefix(reason, "verdict ") || s.branchTip() != tip {
		t.Errorf("a verifier with a finding: refused %q, the branch at %s (was %s); want verdict and nothing pushed", reason, s.branchTip(), tip)
	}
}

// A session of another task is a fault of the caller: an error, nothing written, nothing pushed.
func TestAVerifierOfAnotherTaskIsAnErrorAndNothingIsPushed(t *testing.T) {
	s := newShip(t)
	// The other task is registered: a task session posts its summary to the issue of its task (row 66v).
	files, err := registerFiles(mustRead(t, s.r, "tasks.tsv"), nil, Task{ID: "T-cd34", Issue: 9, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Opened: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	delete(files, "transitions.tsv")
	if err := s.r.commit(context.Background(), files, "tasks: register T-cd34"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.r.StartAttempt(context.Background(), "T-cd34", s.a.base); err != nil {
		t.Fatal(err)
	}
	spec := harness(t, s.a, "printf '"+octal(okResult)+"' > ../result/result.tsv")
	spec.Task, spec.Role, spec.Pair.Harness = "T-cd34", roleOfStep("verification"), "devin"
	other, err := s.r.TaskSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("the session of another task: %v", err)
	}
	s.in.Verifier = other
	tip := s.branchTip()
	id, err := s.closeOut(s.V, s.in, closing())
	if err == nil || s.branchTip() != tip || s.bound(id) != "" {
		t.Errorf("the close-out with a verifier of another task: %v, the branch at %s (was %s), bound %q; want an error and nothing pushed", err, s.branchTip(), tip, s.bound(id))
	}
	if _, err := s.r.CheckRendered(context.Background(), s.in); err == nil {
		t.Error("CheckRendered with no rendered commit named: no error")
	}
}

// CarryVerify checks the edge again, whatever pushed the head: a session that was bound by the push of any other task session is
// not a close-out (the plan review of #363, condition 3).
func TestCarryVerifyChecksTheWholeEdgeAndPostsNothingForAHeadThatIsNoCloseOut(t *testing.T) {
	ctx := context.Background()
	s := newShip(t)
	// A session of the role of the close-out that the ordinary push bound: a child of the verified head that changes a product file.
	id := s.as(roleOfStep("close-out"), 1, "claude", map[string]string{"x.go": "package x\n"}, okResult)
	head := s.bound(id)
	if head == "" {
		t.Fatal("the session was not bound")
	}
	if got, err := s.r.CarryVerify(ctx, s.in, id); err == nil || got != "" || len(statusesOf(s.w, head)) != 0 {
		t.Errorf("a child of the verified head that changes x.go: %q, %v with %d statuses; want an error and no status", got, err, len(statusesOf(s.w, head)))
	}
	// A session that never bound a head, and one that is unknown.
	if _, err := s.r.CarryVerify(ctx, s.in, "S-ffffffff"); err == nil {
		t.Error("an unknown session: no error")
	}
	// A close-out of the two files that is bound by the ordinary push, on a verifier that has a finding: no verdict to carry.
	s2 := newShip(t)
	id2 := s2.as(roleOfStep("close-out"), 1, "claude", map[string]string{taskFile: "t\n", completedFile: "l\n"}, okResult)
	s2.in.Verifier = s2.verifierOf(s2.V, "material", located)
	if _, err := s2.r.CarryVerify(ctx, s2.in, id2); err == nil || len(statusesOf(s2.w, s2.bound(id2))) != 0 {
		t.Errorf("a close-out above a head with no verdict: %v with %d statuses; want an error and no status", err, len(statusesOf(s2.w, s2.bound(id2))))
	}
}

// The rendered commit of the first bet (docs/spec/loop.md, The rendered commit): code made it, so no session is bound to it.
func renderedShip(t *testing.T) (*ship, string) {
	t.Helper()
	s := newShip(t)
	rendered := commitIn(s.loopRun, map[string]string{treePRDPath: renderedPRD(t, ""), treeFacts: treeFactsAt})
	s.in.Rendered = &RenderedInput{Head: rendered, Render: treeRecs}
	return s, rendered
}

func TestTheRenderedCommitCarriesAndTheCloseOutAfterItIsAccepted(t *testing.T) {
	ctx := context.Background()
	s, rendered := renderedShip(t)
	if ok, err := s.r.CheckRendered(ctx, s.in); err != nil || !ok {
		t.Fatalf("CheckRendered = %v, %v; want the rendered commit accepted", ok, err)
	}
	if got, err := s.r.CarryVerify(ctx, s.in, ""); err != nil || got != rendered {
		t.Fatalf("CarryVerify of the rendered commit = %q, %v; want %s", got, err, rendered)
	}
	if got := statusesOf(s.w, rendered); len(got) != 1 || got[0]["description"] != "verified head "+s.V+" plus the rendered commit" || got[0]["state"] != "success" {
		t.Errorf("the statuses of the rendered commit: %v; want success, verified head %s plus the rendered commit", got, s.V)
	}
	id, err := s.closeOut(rendered, s.in, closing())
	if err != nil {
		t.Fatal(err)
	}
	head := s.bound(id)
	if head == "" || s.refusalOf(id) != "" {
		t.Fatalf("the close-out after the rendered commit: bound %q, refused %q; want it accepted", head, s.refusalOf(id))
	}
	if got, err := s.r.CarryVerify(ctx, s.in, id); err != nil || got != head {
		t.Errorf("CarryVerify of the close-out after the rendered commit = %q, %v", got, err)
	}
}

func TestARenderedCommitThatIsNotAPermittedChildOfTheVerifiedHeadDoesNotCarry(t *testing.T) {
	ctx := context.Background()
	// The verifier found something: the rendered commit has no verified parent.
	s, _ := renderedShip(t)
	s.in.Verifier = s.verifierOf(s.V, "material", located)
	if ok, err := s.r.CheckRendered(ctx, s.in); err != nil || ok {
		t.Errorf("an unverified parent: CheckRendered = %v, %v; want false", ok, err)
	}
	if got, err := s.r.CarryVerify(ctx, s.in, ""); err == nil || len(statusesOf(s.w, s.in.Rendered.Head)) != 0 {
		t.Errorf("an unverified parent: CarryVerify = %q, %v with %d statuses; want an error and no status", got, err, len(statusesOf(s.w, s.in.Rendered.Head)))
	}
	id, err := s.closeOut(s.in.Rendered.Head, s.in, closing())
	if err != nil {
		t.Fatal(err)
	}
	if reason := s.refusalOf(id); !strings.HasPrefix(reason, "verdict ") {
		t.Errorf("the close-out after an unverified rendered commit: refused %q; want verdict", reason)
	}
}

func TestARenderedCommitOffTheRecordsIsNotAccepted(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name string
		make func(*ship) *RenderedInput
	}{
		{"a commit above another commit than the verified head", func(s *ship) *RenderedInput {
			commitIn(s.loopRun, map[string]string{"x.go": "package x\n"})
			return &RenderedInput{Head: commitIn(s.loopRun, map[string]string{treePRDPath: renderedPRD(t, ""), treeFacts: treeFactsAt}), Render: treeRecs}
		}},
		{"a path beyond the records", func(s *ship) *RenderedInput {
			return &RenderedInput{Head: commitIn(s.loopRun, map[string]string{treePRDPath: renderedPRD(t, ""), treeFacts: treeFactsAt, "x.go": "package x\n"}), Render: treeRecs}
		}},
		{"bytes that are not the rendering", func(s *ship) *RenderedInput {
			return &RenderedInput{Head: commitIn(s.loopRun, map[string]string{treePRDPath: renderedPRD(t, "does two things"), treeFacts: treeFactsAt}), Render: treeRecs}
		}},
		{"a record for one path of two", func(s *ship) *RenderedInput {
			return &RenderedInput{Head: commitIn(s.loopRun, map[string]string{treePRDPath: renderedPRD(t, ""), treeFacts: treeFactsAt}), Render: treeRecs[:1]}
		}},
		{"the verified head itself", func(s *ship) *RenderedInput { return &RenderedInput{Head: s.V, Render: treeRecs} }},
	} {
		s := newShip(t)
		s.in.Rendered = c.make(s)
		if ok, err := s.r.CheckRendered(ctx, s.in); err != nil || ok {
			t.Errorf("%s: CheckRendered = %v, %v; want false", c.name, ok, err)
		}
		if got, err := s.r.CarryVerify(ctx, s.in, ""); err == nil || len(statusesOf(s.w, s.in.Rendered.Head)) != 0 {
			t.Errorf("%s: CarryVerify = %q, %v with %d statuses; want an error and no status", c.name, got, err, len(statusesOf(s.w, s.in.Rendered.Head)))
		}
	}
}

// renderedPRD gives the rendering of the PRD of the fixture, with the text forged when given.
func renderedPRD(t *testing.T, forged string) string {
	t.Helper()
	prd, finds := spec.RenderPRD([]byte(treePRDPre), treeAnswer)
	if len(finds) != 0 {
		t.Fatalf("the fixture does not render: %v", finds)
	}
	if forged != "" {
		return strings.Replace(string(prd), "does one thing", forged, 1)
	}
	return string(prd)
}

// A verified head that is itself a merge commit carries the verdict as any other: what counts is how it got the verdict (the plan
// review of #363, condition 4). The attempt starts at the merge, so that the verifier's own bound event is outside the change.
func TestAVerifiedHeadThatIsAMergeCommitCarriesTheVerdict(t *testing.T) {
	ctx := context.Background()
	s := newShip(t)
	clone := s.r.Clone
	gitOut(t, clone, "switch", "-q", "-c", "side", s.V)
	commitIn(s.loopRun, map[string]string{"side.txt": "side\n"})
	gitOut(t, clone, "switch", "-q", "main")
	commitIn(s.loopRun, map[string]string{"main.txt": "main\n"})
	gitOut(t, clone, "merge", "-q", "--no-ff", "-m", "merge side", "side")
	merge := strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))
	if parents := strings.Fields(gitOut(t, clone, "log", "-1", "--format=%P", merge)); len(parents) != 2 {
		t.Fatalf("the verified head %s has the parents %v; want a merge commit", merge, parents)
	}
	n, err := s.r.StartAttempt(ctx, "T-ab12", merge)
	if err != nil || n != 2 {
		t.Fatalf("StartAttempt = %d, %v", n, err)
	}
	text := record(merge, "nothing material in scope")
	vspec := harness(t, s.a, planWork("plan.md", text, "comment\t1\tresult/plan.md\t"+commentSum(text)+"\t—"))
	vspec.Role, vspec.Attempt, vspec.Pair.Harness, vspec.Base = roleOfStep("verification"), 2, "devin", merge
	verifier, err := s.r.TaskSession(ctx, vspec)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.r.PostComment(ctx, "T-ab12", verifier, "plan.md"); err != nil {
		t.Fatal(err)
	}
	s.in.Verifier = verifier
	spec := harness(t, s.a, closing())
	spec.Role, spec.Attempt, spec.Base, spec.Pair.Harness = roleOfStep("close-out"), 2, merge, "claude"
	spec.Push = func(id string, d session.Dir, r session.Run) error { return s.r.PushCloseOut(ctx, id, s.in) }
	id, err := s.r.TaskSession(ctx, spec)
	if err != nil {
		t.Fatal(err)
	}
	head := s.bound(id)
	if head == "" || s.refusalOf(id) != "" {
		t.Fatalf("the close-out on a merge commit: bound %q, refused %q; want it accepted", head, s.refusalOf(id))
	}
	if got, err := s.r.CarryVerify(ctx, s.in, id); err != nil || got != head {
		t.Errorf("CarryVerify = %q, %v; want the close-out head %s", got, err, head)
	}
	if posted := statusesOf(s.w, head); len(posted) != 1 || posted[0]["description"] != "verified head "+merge+" plus the close-out" {
		t.Errorf("the statuses: %v; want verified head %s plus the close-out", posted, merge)
	}
}
