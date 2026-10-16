package run

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
)

const (
	planSess   = "S-1a2b3c4d"
	reviewSess = "S-5e6f7a8b"
	taskIssue  = 3
	resultTop  = "kind\tn\tvalue\tsha256\treason\n"
)

// commentForge is a stand-in forge for the comments of one run: Comments and
// Comment over a list of its own, the calls it got, and the switches that make a
// call fail. A comment it stores is the App's (layup-agent), as the adapter's
// answer is. Its answer is the one the adapter decodes from fakeGitHub's JSON
// (the fields of forge.Comment), not the body it was given echoed back.
type commentForge struct {
	forge.Forge
	comments    map[int][]forge.Comment
	calls       []string
	next        int64
	failComment bool // Comment gives an error and stores nothing
	lostAnswer  bool // Comment stores the comment and gives an error
	failList    bool // Comments gives an error
}

func newCommentForge() *commentForge {
	return &commentForge{comments: map[int][]forge.Comment{}, next: 700}
}

func (f *commentForge) Comments(_ context.Context, issue int) ([]forge.Comment, error) {
	f.calls = append(f.calls, "Comments")
	if f.failList {
		return nil, &forge.Error{Call: "Comments", Status: 502, Message: "Bad Gateway"}
	}
	return slices.Clone(f.comments[issue]), nil
}

func (f *commentForge) Comment(_ context.Context, issue int, body string) (int64, error) {
	f.calls = append(f.calls, "Comment")
	if f.failComment {
		return 0, &forge.Error{Call: "Comment", Status: 502, Message: "Bad Gateway"}
	}
	f.next++
	f.comments[issue] = append(f.comments[issue], forge.Comment{ID: f.next, AuthorID: 9001, AuthorLogin: "layup-agent[bot]", App: "layup-agent",
		Created: sessionAt, Updated: sessionAt, Body: body})
	if f.lostAnswer {
		return 0, &forge.Error{Call: "Comment", Status: 504, Message: "Gateway Timeout"}
	}
	return f.next, nil
}

// hold adds a comment that the forge already has.
func (f *commentForge) hold(issue int, c forge.Comment) {
	f.comments[issue] = append(f.comments[issue], c)
}

// edit changes the body of a comment.
func (f *commentForge) edit(issue int, id int64, body string) {
	for i := range f.comments[issue] {
		if f.comments[issue][i].ID == id {
			f.comments[issue][i].Body = body
		}
	}
}

func (f *commentForge) count(call string) int {
	n := 0
	for _, c := range f.calls {
		if c == call {
			n++
		}
	}
	return n
}

// A handoffWorld is the records of task T-ab12 (issue 3) with the default
// transition table, and a Sessions on a stand-in store and forge.
type handoffWorld struct {
	s      *Sessions
	store  *standinStore
	f      *commentForge
	events [][]string
	sess   [][]string
}

func newHandoffWorld(t *testing.T) *handoffWorld {
	t.Helper()
	store := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	files, err := registerFiles(nil, nil, Task{ID: "T-ab12", Issue: taskIssue, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Opened: sessionAt})
	if err != nil {
		t.Fatal(err)
	}
	for p, d := range files {
		store.files[p] = d
	}
	w := &handoffWorld{store: store, f: newCommentForge()}
	w.events = [][]string{{"1", "attempt", "1", "", aSHA, "", "", "2026-10-09T11:00:00Z"}}
	w.s = newSessions(store)
	w.s.Forge, w.s.App = w.f, "layup-agent"
	return w
}

// session records the end of a done session of the attempt with a comment row
// for plan.md and the text of the file kept in the records, as the end of
// a session leaves them; the status row's reason names a comment ID that the
// handoff must not read.
func (w *handoffWorld) session(t *testing.T, id string, attempt int, role, text string, extra ...string) {
	t.Helper()
	w.sess = append(w.sess, []string{id, "T-ab12", strconv.Itoa(attempt), role, "claude", "2.1", "m1", aSHA, aSHA, "9", "3", "1000", "5.0", "30", "", ""})
	n := len(w.events)
	a := strconv.Itoa(attempt)
	w.events = append(w.events,
		[]string{strconv.Itoa(n + 1), "session", a, id, "", "", "", "2026-10-09T11:00:00Z"},
		[]string{strconv.Itoa(n + 2), "result", a, id, "", "", "done", "2026-10-09T11:00:00Z"})
	var err error
	if w.store.files["sessions.tsv"], err = table(records.SessionsSchema, w.sess); err != nil {
		t.Fatal(err)
	}
	if w.store.files["tasks/T-ab12/events.tsv"], err = table(records.EventsSchema, w.events); err != nil {
		t.Fatal(err)
	}
	rows := resultTop + "status\t1\tcompleted\t—\tposted as comment 999-1\n" + strings.Join(extra, "")
	if text != "" {
		rows += "comment\t1\tresult/plan.md\t" + sum(text) + "\t—\n"
		w.store.files[KeptPath("T-ab12", id, "plan.md")] = []byte(text)
	}
	w.store.files["tasks/T-ab12/results/"+id+".tsv"] = []byte(rows)
}

func (w *handoffWorld) eventsNow(t *testing.T) [][]string {
	t.Helper()
	return events(t, w.store, "T-ab12")
}

func (w *handoffWorld) copies(t *testing.T) [][]string {
	t.Helper()
	data := w.store.files["copies.tsv"]
	if data == nil {
		return nil
	}
	rows, err := records.ReadCopies(data)
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

const planBody = "# The plan\r\n\r\nno final newline, trailing  \t"

func TestTheCommentIsPostedFromTheKeptBytes(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	if err != nil {
		t.Fatal(err)
	}
	if bodies := w.f.comments[taskIssue]; len(bodies) != 1 || bodies[0].Body != planBody {
		t.Fatalf("the comments of the task's issue %+v; want one with the bytes of the kept file", bodies)
	}
	if got != (Posted{ID: 701, Seen: 1, SHA: sum(planBody)}) {
		t.Errorf("Posted = %+v; want the ID that the forge gave (not 999, which the result's text names)", got)
	}
	rows := w.copies(t)
	if len(rows) != 1 || rows[0][0] != "701" || rows[0][1] != "1" || rows[0][2] != "3" || rows[0][3] != "9001" || rows[0][5] != "layup-agent" || rows[0][8] != sum(planBody) || rows[0][9] != "copies/701-1.md" {
		t.Errorf("copies.tsv %q; want comment 701, seen 1, issue 3, the App's author, the SHA-256 of the file", rows)
	}
	if string(w.store.files["copies/701-1.md"]) != planBody {
		t.Errorf("the copy %q; want the bytes of the file", w.store.files["copies/701-1.md"])
	}
	if len(w.store.commits) != 1 || len(w.store.commits[0]) != 2 {
		t.Errorf("%d commits; want one holding copies.tsv and the body", len(w.store.commits))
	}
	if !slices.Equal(w.f.calls, []string{"Comments", "Comment", "Comments"}) {
		t.Errorf("the forge calls %q; want a read, the one post, a read back", w.f.calls)
	}
}

func TestACommentThatHasACopyIsNotPostedAgain(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	first, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	if err != nil {
		t.Fatal(err)
	}
	second, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	if err != nil || second != first || w.f.count("Comment") != 1 || len(w.store.commits) != 1 {
		t.Errorf("the second call: %+v, %v; %d posts, %d commits; want the first copy, one post and one commit", second, err, w.f.count("Comment"), len(w.store.commits))
	}
	// The same bytes on another issue are no copy of this issue's comment.
	w2 := newHandoffWorld(t)
	w2.session(t, planSess, 1, "developer", planBody)
	other := sum(planBody)
	w2.store.files["copies.tsv"], _ = table(records.CopiesSchema, [][]string{{"650", "1", "9", "9001", "layup-agent[bot]", "layup-agent", "2026-10-09T11:00:00Z", "2026-10-09T11:00:00Z", other, "copies/650-1.md"}})
	if _, err := w2.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err != nil || w2.f.count("Comment") != 1 {
		t.Errorf("a copy of the same bytes on issue 9: %v, %d posts; want the post on issue 3", err, w2.f.count("Comment"))
	}
}

func TestACommentThatTheForgeHoldsIsCopiedAndNotPostedAgain(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	// A person's comment of the same text, and the App's of another text, are no copy.
	w.f.hold(taskIssue, forge.Comment{ID: 650, AuthorID: 101, AuthorLogin: "pharzam", Body: planBody, Created: sessionAt, Updated: sessionAt})
	w.f.hold(taskIssue, forge.Comment{ID: 651, AuthorID: 9001, AuthorLogin: "layup-agent[bot]", App: "layup-agent", Body: planBody + "!", Created: sessionAt, Updated: sessionAt})
	w.f.next = 700
	if _, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err != nil || w.f.count("Comment") != 1 {
		t.Fatalf("a person's comment and a different App comment: %v, %d posts; want one post", err, w.f.count("Comment"))
	}
	// The stop between the post and the copy: the App's comment is there, copies.tsv has no row.
	w = newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.f.hold(taskIssue, forge.Comment{ID: 655, AuthorID: 9001, AuthorLogin: "layup-agent[bot]", App: "layup-agent", Body: planBody, Created: sessionAt, Updated: sessionAt})
	got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	if err != nil || got != (Posted{ID: 655, Seen: 1, SHA: sum(planBody)}) || w.f.count("Comment") != 0 {
		t.Errorf("the comment the forge holds: %+v, %v, %d posts; want it copied and none posted", got, err, w.f.count("Comment"))
	}
	if rows := w.copies(t); len(rows) != 1 || rows[0][0] != "655" {
		t.Errorf("copies.tsv %q; want the row of comment 655", rows)
	}
}

func TestAPostThatFailsLeavesNoCopyAndTheNextCallPosts(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.f.failComment = true
	if _, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err == nil || len(w.store.commits) != 0 {
		t.Fatalf("a post that fails: %v, %d commits; want an error and no copy", err, len(w.store.commits))
	}
	w.f.failComment = false
	if got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err != nil || got.ID != 701 || w.f.count("Comment") != 2 {
		t.Errorf("the next call: %+v, %v, %d post calls; want the post from the kept copy", got, err, w.f.count("Comment"))
	}
	w = newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.f.failList = true
	if _, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err == nil || w.f.count("Comment") != 0 {
		t.Errorf("a read of the comments that fails: %v, %d posts; want an error before any post", err, w.f.count("Comment"))
	}
}

func TestAPostWhoseAnswerWasLostIsFoundAndCopiedNotPostedTwice(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.f.lostAnswer = true
	if _, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err == nil || len(w.store.commits) != 0 {
		t.Fatalf("a post whose answer is lost: %v, %d commits; want an error and no copy", err, len(w.store.commits))
	}
	w.f.lostAnswer = false
	got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	if err != nil || got.ID != 701 || w.f.count("Comment") != 1 || len(w.copies(t)) != 1 {
		t.Errorf("the next call: %+v, %v, %d posts; want the stored comment copied and no second post", got, err, w.f.count("Comment"))
	}
}

func TestALostLeaseStopsTheComment(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.store.holder = "ffffffffffffffff"
	_, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	var lost *LostError
	if !errors.As(err, &lost) || w.f.count("Comment") != 0 || len(w.store.commits) != 0 {
		t.Errorf("a lease of another run: %v, %d posts, %d commits; want a LostError and no write", err, w.f.count("Comment"), len(w.store.commits))
	}
}

func TestThePostIsFencedByTheLease(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.store.refuse = 1 // one refused push: Fenced reads the lease, held by this run, and tries once more
	if got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err != nil || got.ID != 701 || len(w.store.commits) != 1 {
		t.Errorf("one refused push: %+v, %v, %d commits; want the copy written at the second try", got, err, len(w.store.commits))
	}
	w = newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	w.store.refuse = 2
	_, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
	var lost *LostError
	if !errors.As(err, &lost) {
		t.Fatalf("two refused pushes: %v; want a LostError", err)
	}
	// The comment was posted before the copy was refused: the next run adopts it.
	w.store.refuse = 0
	if got, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err != nil || got.ID != 701 || w.f.count("Comment") != 1 {
		t.Errorf("the next call: %+v, %v, %d posts; want the comment adopted", got, err, w.f.count("Comment"))
	}
}

func TestAFileThatCannotBeACommentIsRefusedBeforeThePost(t *testing.T) {
	for _, c := range []struct{ name, text string }{
		{"over the limit", strings.Repeat("a", MaxCommentChars+1)},
		{"not valid UTF-8", "a\xffb"},
	} {
		w := newHandoffWorld(t)
		w.session(t, planSess, 1, "developer", c.text)
		_, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
		var refusal HandoffRefusal
		if !errors.As(err, &refusal) || refusal.Reason != "handoff" || len(w.f.calls) != 0 {
			t.Errorf("%s: %v, forge calls %q; want a HandoffRefusal and no forge call", c.name, err, w.f.calls)
		}
		ev := w.eventsNow(t)
		if last := ev[len(ev)-1]; last[1] != "refused" || last[3] != planSess || last[6] != "handoff" {
			t.Errorf("%s: the last event %q; want refused, handoff, of the session", c.name, last)
		}
		n := len(ev)
		w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md")
		if again := w.eventsNow(t); len(again) != n {
			t.Errorf("%s: a second call added %d events; want none", c.name, len(again)-n)
		}
	}
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	if _, err := w.s.PostComment(context.Background(), "T-ab12", planSess, "plan.md"); err == nil || len(w.f.calls) != 0 {
		t.Errorf("no kept file: %v; want an error and no forge call", err)
	}
}

// planInput is the handoff of a plan to its review.
func planInput() HandoffInput {
	return HandoffInput{Task: "T-ab12", Session: planSess, From: "task plan", To: "plan review", File: "plan.md", Kind: "plan",
		Items: map[string]bool{"plan-fields": true}}
}

func posted(t *testing.T, w *handoffWorld, session, text string) {
	t.Helper()
	w.session(t, session, 1, "developer", text)
	if _, err := w.s.PostComment(context.Background(), "T-ab12", session, "plan.md"); err != nil {
		t.Fatal(err)
	}
}

func TestTheHandoffEventNamesTheCopyThatCodeRead(t *testing.T) {
	w := newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	commits := len(w.store.commits)
	v, err := w.s.RecordHandoff(context.Background(), planInput())
	if err != nil || !v.Valid || v.Detail != "valid comment 701-1" {
		t.Fatalf("RecordHandoff = %+v, %v; want valid comment 701-1 (not 999-1 of the result's text)", v, err)
	}
	ev := w.eventsNow(t)
	want := []string{strconv.Itoa(len(ev)), "handoff", "1", planSess, "", "", "valid comment 701-1", "2026-10-09T12:00:00Z"}
	if last := ev[len(ev)-1]; !slices.Equal(last, want) {
		t.Errorf("the event %q; want %q", last, want)
	}
	if len(w.store.commits) != commits+1 || len(w.store.commits[commits]) != 1 {
		t.Errorf("%d commits after; want one more, holding events.tsv alone", len(w.store.commits)-commits)
	}
}

func TestACommentThatWasNeverPostedIsInvalidArtifactsWithNoCopy(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", planBody)
	v, err := w.s.RecordHandoff(context.Background(), planInput())
	if err != nil || v.Valid || v.Detail != "invalid artifacts" {
		t.Errorf("RecordHandoff = %+v, %v; want invalid artifacts with no comment named", v, err)
	}
	if len(w.copies(t)) != 0 {
		t.Error("a copy exists for a comment that was never posted")
	}
	// No comment row for the file: the artifact is absent.
	w = newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	if v, err := w.s.RecordHandoff(context.Background(), planInput()); err != nil || v.Detail != "invalid artifacts" {
		t.Errorf("a file the result does not name: %+v, %v; want invalid artifacts", v, err)
	}
}

func TestAnEditedCommentIsANewCopyAndTheFirstHandoffIsInvalid(t *testing.T) {
	w := newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	w.f.edit(taskIssue, 701, "edited after the copy\n")
	v, err := w.s.RecordHandoff(context.Background(), planInput())
	if err != nil || v.Valid || v.Detail != "invalid artifacts comment 701-1" {
		t.Fatalf("RecordHandoff = %+v, %v; want invalid artifacts naming the first copy", v, err)
	}
	rows := w.copies(t)
	if len(rows) != 2 || rows[0][1] != "1" || rows[0][8] != sum(planBody) || rows[1][1] != "2" || rows[1][8] != sum("edited after the copy\n") {
		t.Fatalf("copies.tsv %q; want the first copy kept and the edit as seen 2", rows)
	}
	if string(w.store.files["copies/701-1.md"]) != planBody || string(w.store.files["copies/701-2.md"]) != "edited after the copy\n" {
		t.Error("the copy files are not the two bodies")
	}
	last := w.store.commits[len(w.store.commits)-1]
	if last["copies.tsv"] == nil || last["copies/701-2.md"] == nil || last["tasks/T-ab12/events.tsv"] == nil {
		t.Errorf("the commit holds %d files; want the new copy and the event in one", len(last))
	}
	// A new handoff names the new copy: a session whose file is the edited text.
	w.session(t, reviewSess, 1, "developer", "edited after the copy\n")
	if _, err := w.s.PostComment(context.Background(), "T-ab12", reviewSess, "plan.md"); err != nil || w.f.count("Comment") != 1 {
		t.Fatalf("the post of the edited text: %v, %d posts; want the copy seen 2 adopted, no post", err, w.f.count("Comment"))
	}
	in := planInput()
	in.Session = reviewSess
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || !v.Valid || v.Detail != "valid comment 701-2" {
		t.Errorf("the new handoff = %+v, %v; want valid comment 701-2", v, err)
	}
	// An edit back to the first text is a third copy, not the first again.
	w.f.edit(taskIssue, 701, planBody)
	in.Session = planSess
	in.Items = map[string]bool{"plan-fields": true}
	v, err = w.s.RecordHandoff(context.Background(), in)
	if err != nil || !v.Valid || v.Detail != "valid comment 701-3" {
		t.Errorf("after the edit back: %+v, %v; want valid, naming the last copy of these bytes (seen 3), which the read of the edit added", v, err)
	}
	if rows := w.copies(t); len(rows) != 3 || rows[2][1] != "3" {
		t.Errorf("copies.tsv %q; want seen 3 for the text that came back", rows)
	}
}

func TestAnItemThatTheCallerDoesNotPassIsFalse(t *testing.T) {
	for _, c := range []struct {
		name  string
		items map[string]bool
		want  string
	}{
		{"none", nil, "invalid plan-fields comment 701-1"},
		{"false", map[string]bool{"plan-fields": false}, "invalid plan-fields comment 701-1"},
		{"another item only", map[string]bool{"verdict-field": true}, "invalid plan-fields comment 701-1"},
		{"true", map[string]bool{"plan-fields": true}, "valid comment 701-1"},
	} {
		w := newHandoffWorld(t)
		posted(t, w, planSess, planBody)
		in := planInput()
		in.Items = c.items
		if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != c.want {
			t.Errorf("%s: %+v, %v; want %s", c.name, v, err, c.want)
		}
	}
}

func TestTheProgramItemsOfTheCallerAreHandedToTheCheck(t *testing.T) {
	h64 := strings.Repeat("c", 64)
	h40 := strings.Repeat("d", 40)
	// code (rendered task): no session, no schema; rendering-equal is the caller's.
	for _, c := range []struct {
		name  string
		items map[string]bool
		want  string
	}{
		{"rendering-equal true", map[string]bool{"rendering-equal": true}, "valid"},
		{"rendering-equal false", map[string]bool{"rendering-equal": false}, "invalid rendering-equal"},
		{"rendering-equal not passed", nil, "invalid rendering-equal"},
	} {
		w := newHandoffWorld(t)
		in := HandoffInput{Task: "T-ab12", Attempt: 1, From: "code (rendered task)", To: "verifier", CodeWritten: true, Items: c.items,
			Artifacts: []handoff.Artifact{{Kind: "rendering", Named: h64, Found: h64}}}
		v, err := w.s.RecordHandoff(context.Background(), in)
		if err != nil || v.Detail != c.want {
			t.Errorf("%s: %+v, %v; want %s", c.name, v, err, c.want)
			continue
		}
		ev := w.eventsNow(t)
		if last := ev[len(ev)-1]; last[1] != "handoff" || last[3] != "" || last[2] != "1" || last[6] != c.want {
			t.Errorf("%s: the event %q; want a handoff of attempt 1 with no session", c.name, last)
		}
	}
	// developer (specification) -> verifier: spec-sections is the caller's.
	for _, c := range []struct {
		name  string
		items map[string]bool
		want  string
	}{
		{"spec-sections true", map[string]bool{"spec-sections": true}, "valid"},
		{"spec-sections false", map[string]bool{"spec-sections": false}, "invalid spec-sections"},
		{"not passed", nil, "invalid spec-sections"},
	} {
		w := newHandoffWorld(t)
		w.session(t, planSess, 1, "developer", "")
		in := HandoffInput{Task: "T-ab12", Session: planSess, From: "developer (specification)", To: "verifier", Items: c.items,
			Artifacts: []handoff.Artifact{{Kind: "commits", Named: h40, Found: h40}, {Kind: "task-file", Named: h64, Found: h64}}}
		if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != c.want {
			t.Errorf("%s: %+v, %v; want %s", c.name, v, err, c.want)
		}
	}
	// An artifact the caller names whose hash differs from the one code found is invalid.
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	in := HandoffInput{Task: "T-ab12", Session: planSess, From: "developer (specification)", To: "verifier", Items: map[string]bool{"spec-sections": true},
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: h40, Found: strings.Repeat("e", 40)}, {Kind: "task-file", Named: h64, Found: h64}}}
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid artifacts" {
		t.Errorf("a hash that differs: %+v, %v; want invalid artifacts", v, err)
	}
}

func TestAClaimOfTheCallerForWhatCodeComputesIsNotRead(t *testing.T) {
	all := map[string]bool{"schema": true, "artifacts": true, "question-row": true, "answer-row": true, "plan-fields": true}
	// A result file that the block refuses.
	w := newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	w.store.files["tasks/T-ab12/results/"+planSess+".tsv"] = []byte("not a result\n")
	in := planInput()
	in.Items = all
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid schema" {
		t.Errorf("a result the block refuses: %+v, %v; want invalid schema", v, err)
	}
	// No result file at all.
	w = newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	delete(w.store.files, "tasks/T-ab12/results/"+planSess+".tsv")
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid schema" {
		t.Errorf("no result file: %+v, %v; want invalid schema", v, err)
	}
	// question-row and answer-row are computed from the result, whatever the map says.
	w = newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	if v, err := w.s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: planSess, From: "asking role", To: "owner role", Items: all}); err != nil || v.Detail != "invalid question-row" {
		t.Errorf("a result with no question: %+v, %v; want invalid question-row", v, err)
	}
	if v, err := w.s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: planSess, From: "owner role", To: "asking role", Items: all}); err != nil || v.Detail != "invalid answer-row" {
		t.Errorf("a result with no answer: %+v, %v; want invalid answer-row", v, err)
	}
}

func TestTheQuestionAndAnswerHandoffsAreComputedFromTheRecords(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n")
	w.session(t, reviewSess, 2, "software-architect", "", "answer\t1\tYes\t—\tdocs/spec/forge.md\n")
	if v, err := w.s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: planSess, From: "asking role", To: "owner role"}); err != nil || v.Detail != "valid" {
		t.Errorf("a question with a record: %+v, %v; want valid", v, err)
	}
	if v, err := w.s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: reviewSess, From: "owner role", To: "asking role"}); err != nil || v.Detail != "valid" {
		t.Errorf("an answer with a record: %+v, %v; want valid", v, err)
	}
	w.session(t, "S-9a8b7c6d", 3, "developer", "", "question\t1\tMay it?\t—\tinterface-contract no\n")
	if v, err := w.s.RecordHandoff(context.Background(), HandoffInput{Task: "T-ab12", Session: "S-9a8b7c6d", From: "asking role", To: "owner role"}); err != nil || v.Detail != "invalid question-row" {
		t.Errorf("a question with no record: %+v, %v; want invalid question-row", v, err)
	}
}

func TestAPairWithNoRowIsInvalidPairAndAReplacedTableIsRead(t *testing.T) {
	w := newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	in := planInput()
	in.From = "task plan"
	in.To = "nobody"
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid pair comment 701-1" {
		t.Errorf("a pair with no row: %+v, %v; want invalid pair", v, err)
	}
	// The Operator's table: the plan row lists no plan-fields, so no caller result is needed.
	rows, err := handoff.ReadTable(w.store.files["transitions.tsv"])
	if err != nil {
		t.Fatal(err)
	}
	for i := range rows {
		if rows[i].From == "task plan" {
			rows[i].Items = []string{"schema", "artifacts"}
		}
	}
	var b strings.Builder
	if err := handoff.WriteTable(&b, rows); err != nil {
		t.Fatal(err)
	}
	w.store.files["transitions.tsv"] = []byte(b.String())
	in = planInput()
	in.Items = nil
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "valid comment 701-1" {
		t.Errorf("the replaced table: %+v, %v; want valid with no caller item", v, err)
	}
	// A table that fails its check is an error, and nothing is written.
	w.store.files["transitions.tsv"] = []byte("not a table\n")
	n := len(w.store.commits)
	in.Session = planSess
	if _, err := w.s.RecordHandoff(context.Background(), in); err == nil || len(w.store.commits) != n {
		t.Errorf("a table that fails its check: %v, %d new commits; want an error and none", err, len(w.store.commits)-n)
	}
}

func TestAHandoffIsNotRecordedTwice(t *testing.T) {
	w := newHandoffWorld(t)
	posted(t, w, planSess, planBody)
	first, err := w.s.RecordHandoff(context.Background(), planInput())
	if err != nil {
		t.Fatal(err)
	}
	n, commits := len(w.eventsNow(t)), len(w.store.commits)
	second, err := w.s.RecordHandoff(context.Background(), planInput())
	if err != nil || second != first || len(w.eventsNow(t)) != n || len(w.store.commits) != commits {
		t.Errorf("the second call: %+v, %v; %d events, %d commits; want the recorded verdict and no write", second, err, len(w.eventsNow(t))-n, len(w.store.commits)-commits)
	}
	// A different outcome of the same session is a new event.
	in := planInput()
	in.Items = nil
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid plan-fields comment 701-1" || len(w.eventsNow(t)) != n+1 {
		t.Errorf("another outcome: %+v, %v; want a new event", v, err)
	}
}

func TestAHandoffOfASessionWithoutAStartRowIsAnError(t *testing.T) {
	w := newHandoffWorld(t)
	if _, err := w.s.RecordHandoff(context.Background(), planInput()); err == nil || len(w.store.commits) != 0 {
		t.Errorf("a session with no start row: %v, %d commits; want an error and no write", err, len(w.store.commits))
	}
}
