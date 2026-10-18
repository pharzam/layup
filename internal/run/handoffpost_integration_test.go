//go:build integration

package run

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// commentWorld is sessionWorld with the adapter for the forge, the register
// that opens T-ab12 on issue 3, and the lease held by the run: the world of a
// restart that posts a comment and records a handoff.
func commentWorld(t *testing.T) (*world, *Sessions, worldAt) {
	return commentWorldWith(t, nil)
}

// commentWorldWith is commentWorld whose target also holds the baseline files by path.
func commentWorldWith(t *testing.T, baseline map[string]string) (*world, *Sessions, worldAt) {
	t.Helper()
	w, r, _, base := sessionWorldWith(t, baseline)
	r.Forge, r.App = w.config("0.1.0-dev").Forge, "layup-agent"
	ctx := context.Background()
	files, err := registerFiles(nil, nil, Task{ID: "T-ab12", Issue: taskIssue, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Opened: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.commit(ctx, files, "tasks: register T-ab12"); err != nil {
		t.Fatal(err)
	}
	if err := r.Store.WriteLease(ctx, LeaseRow{Run: r.RunID, Host: "host-1", Version: "0.1.0-dev", Started: time.Now().UTC().Truncate(time.Second), Heartbeat: 1, State: "held"}, "lease"); err != nil {
		t.Fatal(err)
	}
	return w, r, base
}

// afterRestart gives the Sessions of a restart: the same records, forge and run, no
// state of the first (the session directory is gone).
func afterRestart(r *Sessions) *Sessions {
	return &Sessions{Store: r.Store, RunID: r.RunID, Host: r.Host, Target: r.Target, Clone: r.Clone, Branch: r.Branch, Now: time.Now, Forge: r.Forge, Control: r.Control, App: r.App, Bot: r.Bot, Web: r.Web}
}

// commentSession runs a session of the role that writes result/plan.md with text
// and names it by a comment row.
func commentSession(t *testing.T, r *Sessions, base worldAt, role, text string) string {
	t.Helper()
	spec := harness(t, base, planWork("plan.md", text, "comment\t1\tresult/plan.md\t"+commentSum(text)+"\t—"))
	spec.Role = role
	id, err := r.TaskSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	return id
}

func commentsOf(w *world) []map[string]any {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	return append([]map[string]any(nil), w.fake.comments[taskIssue]...)
}

// copiesOfIssue counts the rows of copies.tsv that copy a comment of the task's issue
// (Start copied the dead-man notice of the control issue, whose rows are no concern).
func copiesOfIssue(t *testing.T, w *world) int {
	t.Helper()
	rows, err := records.ReadCopies([]byte(keptRecord(t, w, "copies.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	n := 0
	for _, r := range rows {
		if r[copyIssue] == fmt.Sprint(taskIssue) {
			n++
		}
	}
	return n
}

func editComment(w *world, i int, body string) {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	w.fake.comments[taskIssue][i]["body"] = body
}

func keptRecord(t *testing.T, w *world, path string) string {
	return show(t, w.bare, "layup-records", path)
}

// The demo of row 66d (loop.md, the artifact that is a comment): for each of the
// kinds plan, plan-review and review-record, the file of the session, the comment
// that code posts from its bytes after a restart, the copy and its row, the
// handoff event with the ID that code read, the three equal hashes.
func TestThePlanThePlanReviewAndTheReviewRecordArePostedFromTheKeptCopy(t *testing.T) {
	for _, c := range []struct {
		kind, role, from, to, item string
		text                       func(worldAt) string
	}{
		{"plan", "software-engineer", "task plan", "plan review", "plan-fields", func(worldAt) string { return planText }},
		{"plan-review", "qa-engineer", "plan review", "developer", "verdict-field", func(worldAt) string { return "verdict: approve-with-conditions\r\nno final newline" }},
		{"review-record", "qa-engineer", "verifier", "close-out", "nothing-material", func(b worldAt) string { return record(b.base, "nothing material in scope") }},
	} {
		t.Run(c.kind, func(t *testing.T) {
			w, r, base := commentWorld(t)
			body := c.text(base)
			id := commentSession(t, r, base, c.role, body)
			if n := len(commentsOf(w)); n != 0 {
				t.Fatalf("the session's end left %d comments on the task's issue; want none (code posts after the result)", n)
			}
			r2 := afterRestart(r)
			ctx := context.Background()
			p, err := r2.PostComment(ctx, "T-ab12", id, "plan.md")
			if err != nil {
				t.Fatalf("PostComment: %v", err)
			}
			cs := commentsOf(w)
			if len(cs) != 1 || cs[0]["body"] != body {
				t.Fatalf("the comments of the task's issue %v; want one with the bytes of the file", cs)
			}
			if p.Seen != 1 || p.SHA != commentSum(body) || cs[0]["id"] != int(p.ID) {
				t.Errorf("Posted %+v for the comment %v", p, cs[0]["id"])
			}
			copies := keptRecord(t, w, "copies.tsv")
			if !strings.Contains(copies, fmt.Sprintf("\n%d\t1\t%d\t9001\tlayup-agent[bot]\tlayup-agent\t", p.ID, taskIssue)) || !strings.Contains(copies, "\t"+commentSum(body)+"\t") {
				t.Errorf("copies.tsv:\n%s", copies)
			}
			if got := keptRecord(t, w, fmt.Sprintf("copies/%d-1.md", p.ID)); got != body {
				t.Errorf("the copy %q; want the bytes of the file", got)
			}
			v, err := r2.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: id, From: c.from, To: c.to, File: "plan.md", Kind: c.kind, Items: map[string]bool{c.item: true}})
			want := fmt.Sprintf("valid comment %d-1", p.ID)
			if err != nil || !v.Valid || v.Detail != want {
				t.Fatalf("RecordHandoff = %+v, %v; want %s", v, err, want)
			}
			if events := keptRecord(t, w, "tasks/T-ab12/events.tsv"); !strings.Contains(events, "\thandoff\t1\t"+id+"\t—\t—\t"+want+"\t") {
				t.Errorf("the events:\n%s", events)
			}
			if kept := keptRecord(t, w, "tasks/T-ab12/results/"+id+"/plan.md"); kept != body {
				t.Errorf("the kept file %q", kept)
			}
			// The records read through the real checks (guardrails §2: a record that no test read).
			if _, err := records.ReadEvents([]byte(keptRecord(t, w, "tasks/T-ab12/events.tsv"))); err != nil {
				t.Errorf("events.tsv: %v", err)
			}
			if _, err := records.ReadCopies([]byte(copies)); err != nil {
				t.Errorf("copies.tsv: %v", err)
			}
		})
	}
}

func TestAnEditedCommentMakesTheSecondCopyAndTheNewHandoffNamesIt(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	id := commentSession(t, r, base, "software-engineer", planText)
	p, err := r.PostComment(ctx, "T-ab12", id, "plan.md")
	if err != nil {
		t.Fatal(err)
	}
	editComment(w, 0, "edited by a person\n")
	in := HandoffInput{Task: "T-ab12", Session: id, From: "task plan", To: "plan review", File: "plan.md", Kind: "plan", Items: map[string]bool{"plan-fields": true}}
	v, err := r.RecordHandoff(ctx, in)
	if want := fmt.Sprintf("invalid artifacts comment %d-1", p.ID); err != nil || v.Valid || v.Detail != want {
		t.Fatalf("RecordHandoff = %+v, %v; want %s", v, err, want)
	}
	if got := keptRecord(t, w, fmt.Sprintf("copies/%d-2.md", p.ID)); got != "edited by a person\n" {
		t.Errorf("the second copy %q", got)
	}
	if got := keptRecord(t, w, fmt.Sprintf("copies/%d-1.md", p.ID)); got != planText {
		t.Errorf("the first copy %q; want it kept", got)
	}
}

// failOnCopy fails the first records commit that writes copies.tsv, as a stop
// between the post and the copy does.
type failOnCopy struct {
	SessionStore
	failed bool
}

func (s *failOnCopy) Commit(ctx context.Context, files map[string][]byte, message string) error {
	if _, ok := files["copies.tsv"]; ok && !s.failed {
		s.failed = true
		return errors.New("the run stopped before the copy")
	}
	return s.SessionStore.Commit(ctx, files, message)
}

func TestAStopBetweenThePostAndTheCopyDoesNotPostTwice(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	id := commentSession(t, r, base, "software-engineer", planText)
	stopped := afterRestart(r)
	stopped.Store = &failOnCopy{SessionStore: r.Store}
	if _, err := stopped.PostComment(ctx, "T-ab12", id, "plan.md"); err == nil {
		t.Fatal("the stop: no error")
	}
	if n := len(commentsOf(w)); n != 1 {
		t.Fatalf("%d comments after the stop; want the one that was posted", n)
	}
	p, err := afterRestart(r).PostComment(ctx, "T-ab12", id, "plan.md")
	if err != nil || len(commentsOf(w)) != 1 || p.Seen != 1 {
		t.Errorf("the restart: %+v, %v, %d comments; want the comment copied, none posted", p, err, len(commentsOf(w)))
	}
	if !strings.Contains(keptRecord(t, w, "copies.tsv"), fmt.Sprintf("\n%d\t1\t", p.ID)) {
		t.Error("copies.tsv has no row of the comment")
	}
}

func TestAStopBeforeThePostPostsFromTheKeptCopyAtTheRestart(t *testing.T) {
	w, r, base := commentWorld(t)
	id := commentSession(t, r, base, "software-engineer", planText)
	w.fake.mu.Lock()
	w.fake.failNote = "no final newline" // the first post of the plan fails once with a 500
	w.fake.mu.Unlock()
	if _, err := afterRestart(r).PostComment(context.Background(), "T-ab12", id, "plan.md"); err == nil || len(commentsOf(w)) != 0 {
		t.Fatalf("a post that the forge refuses: %v, %d comments; want an error and none", err, len(commentsOf(w)))
	}
	if n := copiesOfIssue(t, w); n != 0 {
		t.Errorf("%d copies of the task's issue exist for a post that failed", n)
	}
	if _, err := afterRestart(r).PostComment(context.Background(), "T-ab12", id, "plan.md"); err != nil || len(commentsOf(w)) != 1 || commentsOf(w)[0]["body"] != planText {
		t.Errorf("the restart: %v, %d comments; want the post of the kept bytes", err, len(commentsOf(w)))
	}
}

func TestALostLeaseDuringTheCommentPostsNothing(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	id := commentSession(t, r, base, "software-engineer", planText)
	// Another run takes the lease.
	if err := r.Store.WriteLease(ctx, LeaseRow{Run: "fedcba9876543210", Host: "host-2", Version: "0.1.0-dev", Started: time.Now().UTC().Truncate(time.Second), Heartbeat: 1, State: "held"}, "takeover"); err != nil {
		t.Fatal(err)
	}
	_, err := afterRestart(r).PostComment(ctx, "T-ab12", id, "plan.md")
	var lost *LostError
	if !errors.As(err, &lost) || lost.Holder != "fedcba9876543210" || len(commentsOf(w)) != 0 {
		t.Errorf("%v, %d comments; want a LostError naming the other run and no post", err, len(commentsOf(w)))
	}
}

func TestTheLimitOfACommentThroughTheAdapter(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	at := planText + strings.Repeat("€", MaxCommentChars-len([]rune(planText))) // 65,536 characters
	id := commentSession(t, r, base, "software-engineer", at)
	p, err := afterRestart(r).PostComment(ctx, "T-ab12", id, "plan.md")
	if cs := commentsOf(w); err != nil || len(cs) != 1 || cs[0]["body"] != at || p.SHA != commentSum(at) {
		t.Fatalf("the body at the limit: %v, %d comments; want it posted and read back equal", err, len(cs))
	}
	v, err := afterRestart(r).RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: id, From: "task plan", To: "plan review", File: "plan.md", Kind: "plan", Items: map[string]bool{"plan-fields": true}})
	if err != nil || !v.Valid {
		t.Errorf("the handoff of the body at the limit: %+v, %v; want valid", v, err)
	}
	// One character over: refused before the post, the event refused, handoff.
	w2, r2, base2 := commentWorld(t)
	over := at + "€"
	id2 := commentSession(t, r2, base2, "software-engineer", over)
	_, err = afterRestart(r2).PostComment(ctx, "T-ab12", id2, "plan.md")
	var refusal HandoffRefusal
	if !errors.As(err, &refusal) || len(commentsOf(w2)) != 0 {
		t.Fatalf("one character over: %v, %d comments; want a HandoffRefusal and no request", err, len(commentsOf(w2)))
	}
	if events := keptRecord(t, w2, "tasks/T-ab12/events.tsv"); !strings.Contains(events, "\trefused\t1\t"+id2+"\t—\t—\thandoff\t") {
		t.Errorf("the events:\n%s", events)
	}
}

// resultWork writes the result file of a session of the harness: status, then rows
// (with real tabs).
func resultWork(status string, rows ...string) string {
	out := "kind\tn\tvalue\tsha256\treason\nstatus\t1\t" + status + "\t—\twork\n"
	for _, r := range rows {
		out += r + "\n"
	}
	return "printf '" + octal(out) + "' > ../result/result.tsv"
}

func resultsOf(t *testing.T, w *world, ids ...string) map[string][][]string {
	t.Helper()
	out := map[string][][]string{}
	for _, id := range ids {
		rows, err := records.ReadResult([]byte(keptRecord(t, w, "tasks/T-ab12/results/"+id+".tsv")))
		if err != nil {
			t.Fatalf("the result of %s: %v", id, err)
		}
		out[id] = rows
	}
	return out
}

// A question goes to its owner role and its answer is accepted by the next
// attempt that cites it, from the records that the end of each session wrote.
func TestTheQuestionOfASessionIsRoutedAndItsAnswerAcceptedFromTheRecords(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	run := func(role string, attempt int, body string) string {
		spec := harness(t, base, body)
		spec.Role, spec.Attempt = role, attempt
		id, err := r.TaskSession(ctx, spec)
		if err != nil {
			t.Fatalf("TaskSession: %q, %v", id, err)
		}
		return id
	}
	asking := run("developer", 1, resultWork("needs_context", "question\t1\tMay a draft be marked ready?\t—\tinterface-contract no docs/spec/forge.md"))
	qs := QuestionsOf(resultsOf(t, w, asking)[asking])
	if len(qs) != 1 || RouteQuestion(qs[0], DefaultOwners()) != (Route{Role: "Software Architect"}) {
		t.Fatalf("the questions %+v; want one routed to the Software Architect", qs)
	}
	if v, err := r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: asking, From: "asking role", To: "owner role"}); err != nil || v.Detail != "valid" {
		t.Fatalf("the question handoff: %+v, %v; want valid", v, err)
	}
	owner := run("software-architect", 1, resultWork("completed", "answer\t1\tYes, by its node ID\t—\tdocs/spec/forge.md"))
	answer := owner + ".1"
	if v, err := r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: owner, From: "owner role", To: "asking role"}); err != nil || v.Detail != "valid" {
		t.Fatalf("the answer handoff: %+v, %v; want valid", v, err)
	}
	if _, err := r.StartAttempt(ctx, "T-ab12", base.base); err != nil {
		t.Fatal(err)
	}
	next := run("developer", 2, resultWork("completed"))
	// The next attempt that does not cite the answer does not accept it.
	roles := func() Roles { return sessionRoles(startRows(t, w)) }
	read := func() ([][]string, map[string][][]string) {
		ev, err := records.ReadEvents([]byte(keptRecord(t, w, "tasks/T-ab12/events.tsv")))
		if err != nil {
			t.Fatal(err)
		}
		return ev, resultsOf(t, w, asking, owner, next)
	}
	if ev, res := read(); true {
		if who, ok := AnswerAccepted(ev, res, roles(), Owners(DefaultOwners()), answer); ok {
			t.Errorf("accepted by %s with no cite", who)
		}
	}
	// A second session of attempt 2 that cites it is the accepting actor.
	citing := run("developer", 2, resultWork("completed")+" && "+"printf '"+octal("kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tuses the answer "+answer+"\n")+"' > ../result/result.tsv")
	ev, res := read()
	res[citing] = resultsOf(t, w, citing)[citing]
	if who, ok := AnswerAccepted(ev, res, roles(), Owners(DefaultOwners()), answer); !ok || who != citing {
		t.Errorf("AnswerAccepted = %q, %v; want accepted by %s", who, ok, citing)
	}
}

// The Sessions of a run come from the constructor of the run, state.sessions(),
// not from a literal: a stop between the post and the copy, and the restart
// that follows it, post one comment only when that constructor knows the App's
// slug (round 1 of #360, finding 1: it did not, so the adoption of the App's own
// comment was skipped and the restart posted twice).
func TestTheSessionsOfTheRunAdoptTheCommentOfALostCopyAndPostOnce(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	id := commentSession(t, r, base, "software-engineer", planText)
	st := &state{cfg: w.config("0.1.0-dev"), store: r.Store.(*recordsStore), target: r.Clone, defaultBranch: r.Branch,
		start: [][]string{{"issue.control", "2", "run"}}}
	stopped := st.sessions()
	stopped.Store = &failOnCopy{SessionStore: stopped.Store}
	if _, err := stopped.PostComment(ctx, "T-ab12", id, "plan.md"); err == nil || len(commentsOf(w)) != 1 {
		t.Fatalf("the stop after the post: %v, %d comments; want an error and the one posted comment", err, len(commentsOf(w)))
	}
	p, err := st.sessions().PostComment(ctx, "T-ab12", id, "plan.md")
	if err != nil || len(commentsOf(w)) != 1 || p.Seen != 1 {
		t.Errorf("the restart: %+v, %v, %d comments on the task's issue; want the comment adopted and copied, none posted again", p, err, len(commentsOf(w)))
	}
	if n := copiesOfIssue(t, w); n != 1 {
		t.Errorf("%d copies of the task's issue after the restart; want 1", n)
	}
}

// An answer row is a valid record whoever wrote it, so a developer session's
// answer to a question of the Software Architect passes every reader of the
// records; the acceptance alone reads the role of the answering session from
// sessions.tsv (loop.md, A question during the work, T-b3s3).
func TestAnAnswerOfADeveloperSessionIsNotAcceptedFromTheRecords(t *testing.T) {
	w, r, base := commentWorld(t)
	ctx := context.Background()
	run := func(role string, attempt int, body string) string {
		spec := harness(t, base, body)
		spec.Role, spec.Attempt = role, attempt
		id, err := r.TaskSession(ctx, spec)
		if err != nil {
			t.Fatalf("TaskSession: %q, %v", id, err)
		}
		return id
	}
	asking := run("developer", 1, resultWork("needs_context", "question\t1\tMay a draft be marked ready?\t—\tinterface-contract no docs/spec/forge.md"))
	wrong := run("developer", 1, resultWork("completed", "answer\t1\tYes\t—\tdocs/spec/forge.md"))
	if v, err := r.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: wrong, From: "owner role", To: "asking role"}); err != nil || v.Detail != "valid" {
		t.Fatalf("the handoff of the answer: %+v, %v; want valid (a known limit: it reads no role)", v, err)
	}
	if _, err := r.StartAttempt(ctx, "T-ab12", base.base); err != nil {
		t.Fatal(err)
	}
	id := wrong + ".1"
	next := run("developer", 2, resultWork("completed")+" && printf '"+octal("kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tuses the answer "+id+"\n")+"' > ../result/result.tsv")
	ev, err := records.ReadEvents([]byte(keptRecord(t, w, "tasks/T-ab12/events.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	roles := sessionRoles(startRows(t, w))
	if who, ok := AnswerAccepted(ev, resultsOf(t, w, asking, wrong, next), roles, Owners(DefaultOwners()), id); ok {
		t.Errorf("the answer of a developer session is accepted by %s", who)
	}
}

// startRows gives the rows of the sessions.tsv that the run kept.
func startRows(t *testing.T, w *world) [][]string {
	t.Helper()
	rows, err := records.ReadSessions([]byte(keptRecord(t, w, "sessions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// The demo of row 66l (group (a)): a handoff that lists a comment item is valid or
// invalid by what code reads from the kept copy and from the verifier's recorded
// base, whatever the caller says, through the records and the forge.
func TestTheCommentChecksOfTheTaskLoop(t *testing.T) {
	for _, c := range []struct {
		name, role, from, to, kind, item string
		text                             func(base worldAt) string
		claim                            bool
		want                             string
	}{
		{"nothing material at the verifier's base, the caller says false", "qa-engineer", "verifier", "close-out", "review-record", "nothing-material",
			func(b worldAt) string { return record(b.base, "nothing material in scope") }, false, "valid"},
		{"nothing material at another head, the caller says true", "qa-engineer", "verifier", "close-out", "review-record", "nothing-material",
			func(b worldAt) string { return record(hexB, "nothing material in scope") }, true, "invalid nothing-material"},
		{"nothing material with a finding, the caller says true", "qa-engineer", "verifier", "close-out", "review-record", "nothing-material",
			func(b worldAt) string { return record(b.base, "nothing material in scope", located) }, true, "invalid nothing-material"},
		{"findings located, the caller says false", "qa-engineer", "verifier", "developer", "review-record", "findings-located",
			func(b worldAt) string { return record(b.base, "material", located) }, false, "valid"},
		{"a plan of the form, the caller says false", "software-engineer", "task plan", "plan review", "plan", "plan-fields",
			func(worldAt) string { return planText }, false, "valid"},
		{"a plan that is not of the form, the caller says true", "software-engineer", "task plan", "plan review", "plan", "plan-fields",
			func(worldAt) string { return "size: small\n" }, true, "invalid plan-fields"},
		{"a plan review that rejects, the caller says true", "qa-engineer", "plan review", "developer", "plan-review", "verdict-field",
			func(worldAt) string { return "verdict: reject\n" }, true, "invalid verdict-field"},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, r, base := commentWorld(t)
			ctx := context.Background()
			id := commentSession(t, r, base, c.role, c.text(base))
			p, err := afterRestart(r).PostComment(ctx, "T-ab12", id, "plan.md")
			if err != nil {
				t.Fatal(err)
			}
			v, err := afterRestart(r).RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: id, From: c.from, To: c.to, File: "plan.md", Kind: c.kind, Items: map[string]bool{c.item: c.claim}})
			if want := fmt.Sprintf("%s comment %d-1", c.want, p.ID); err != nil || v.Detail != want {
				t.Errorf("RecordHandoff = %+v, %v; want %s", v, err, want)
			}
		})
	}
}
