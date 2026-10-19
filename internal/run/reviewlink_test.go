package run

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// The link of the review record in the body (docs/spec/loop.md, The link of the review record in the body; task T-08rj, #413,
// row 66q). The world is T-ab12 on issue 3 in attempt 1: the developer S-d0000001 bound head1 at the snapshot recA, the verifier
// S-e0000001 reviewed head1 and bound it again at the snapshot recB, and its latest handoff names the copy of comment 601.
const (
	lkDev     = "S-d0000001"
	lkVer     = "S-e0000001"
	lkClose   = "S-c0000001"
	lkComment = 601
	lkURL     = prWeb + "/acme/target/issues/3#issuecomment-601"
)

var lkRecord = record(head1, "nothing material in scope")

// A bodyForge logs the calls of the link and lets a test change the answer of the edit.
type bodyForge struct {
	*forge.Fake
	calls  []string
	answer func(forge.PullRequest) forge.PullRequest
}

func (f *bodyForge) PullRequest(ctx context.Context, q forge.PullRequestQuery) (forge.PullRequest, bool, error) {
	f.calls = append(f.calls, "PullRequest")
	return f.Fake.PullRequest(ctx, q)
}

func (f *bodyForge) EditBody(ctx context.Context, number int, body string) (forge.PullRequest, error) {
	f.calls = append(f.calls, "EditBody")
	pr, err := f.Fake.EditBody(ctx, number, body)
	if err == nil && f.answer != nil {
		pr = f.answer(pr)
	}
	return pr, err
}

func (f *bodyForge) count(call string) int {
	n := 0
	for _, c := range f.calls {
		if c == call {
			n++
		}
	}
	return n
}

// A failingStore fails the read of one path, "lease" for the lease.
type failingStore struct {
	*standinStore
	path string
	err  error
}

func (s *failingStore) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if path == s.path {
		return nil, s.err
	}
	return s.standinStore.ReadFile(ctx, path)
}

func (s *failingStore) ReadLease(ctx context.Context) (LeaseRow, error) {
	if s.path == "lease" {
		return LeaseRow{}, s.err
	}
	return s.standinStore.ReadLease(ctx)
}

type linkWorld struct {
	s        *Sessions
	st       *standinStore
	f        *bodyForge
	sessions [][]string
	events   [][]string
	copies   [][]string
	copyText string
	pr       int
	pr1      forge.PullRequest
}

func newLinkWorld(t *testing.T) *linkWorld {
	t.Helper()
	fake := forge.NewFake()
	fake.SetHead("task/T-ab12/1", head1)
	st := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	w := &linkWorld{st: st, f: &bodyForge{Fake: fake}, copyText: lkRecord}
	w.s = newSessions(st)
	w.s.Forge, w.s.Web = w.f, prWeb
	pr, err := fake.OpenDraft(context.Background(), "task/T-ab12/1", "main", "feat: x", PullBody(3, prWeb, "acme/target", recA))
	if err != nil {
		t.Fatal(err)
	}
	w.pr, w.pr1 = pr.Number, pr
	dev := sessionRow(lkDev, "T-ab12", "implement", "claude")
	dev[sessionBase], dev[sessionRecords] = shaOld, recA
	ver := sessionRow(lkVer, "T-ab12", "verification", "devin")
	ver[sessionBase], ver[sessionRecords] = head1, recB
	w.sessions = [][]string{dev, ver}
	at := "2026-10-18T10:00:00Z"
	w.events = [][]string{
		{"1", "attempt", "1", "", shaOld, "", "", at},
		{"2", "session", "1", lkDev, "", "", "", at},
		{"3", "bound", "1", lkDev, "", head1, "task/T-ab12/1", at},
		{"4", "session", "1", lkVer, "", "", "", at},
		{"5", "bound", "1", lkVer, "", head1, "task/T-ab12/1", at},
		{"6", "handoff", "1", lkVer, "", "", "valid comment 601-1", at},
	}
	w.copies = [][]string{{"601", "1", "3", "9001", "layup-agent[bot]", "layup-agent", at, at, commentSum(lkRecord), "copies/601-1.md"}}
	w.write(t)
	return w
}

// write rewrites the records from the rows of the world.
func (w *linkWorld) write(t *testing.T) {
	t.Helper()
	put := func(path string, data []byte, err error) {
		if err != nil {
			t.Fatal(err)
		}
		w.st.files[path] = data
	}
	data, err := table(records.SessionsSchema, w.sessions)
	put("sessions.tsv", data, err)
	data, err = table(records.EventsSchema, w.events)
	put(eventsPath("T-ab12"), data, err)
	data, err = table(records.CopiesSchema, w.copies)
	put("copies.tsv", data, err)
	w.st.files["copies/601-1.md"] = []byte(w.copyText)
	data, err = registerBytes(w.pr)
	put(tasksPath, data, err)
}

// registerBytes is tasks.tsv with T-ab12 on issue 3 and the pull request pr.
func registerBytes(pr int) ([]byte, error) {
	var b bytes.Buffer
	err := WriteRegister(&b, []Task{{ID: "T-ab12", Issue: 3, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", PR: pr, Opened: time.Date(2026, 10, 18, 9, 0, 0, 0, time.UTC)}})
	return b.Bytes(), err
}

func (w *linkWorld) link() (bool, error) {
	return w.s.LinkReviewRecord(context.Background(), "T-ab12", lkVer)
}

// bodyOf gives the body of the pull request n in the stand-in.
func (w *linkWorld) bodyOf(t *testing.T, n int) string {
	t.Helper()
	pr, found, err := w.f.Fake.PullRequest(context.Background(), forge.PullRequestQuery{Number: n})
	if err != nil || !found {
		t.Fatalf("PullRequest %d: %v, %v", n, found, err)
	}
	return pr.Body
}

func wantBody(records string) string {
	return PullBody(3, prWeb, "acme/target", records, "", reviewLine(prWeb, "acme/target", 3, lkComment))
}

func wantErr(t *testing.T, name string, err error, text string) {
	t.Helper()
	if err == nil || !strings.Contains(err.Error(), text) {
		t.Errorf("%s: err = %v; want one that says %q", name, err, text)
	}
}

// U1. The bytes of the body: the draft's lines unchanged, a blank line, the link of the comment, no final line feed.
func TestTheBodyWithAReviewRecordEndsInTheLinkOfItsComment(t *testing.T) {
	got := wantBody(recA)
	want := "Closes #3\n\nThe records commit this task started from: " + prWeb + "/acme/target/commit/" + recA + "\n\nThe review record: " + lkURL
	if got != want {
		t.Errorf("the body = %q; want %q", got, want)
	}
	if line := reviewLine(prWeb+"/", "acme/target", 3, lkComment); line != "The review record: "+lkURL {
		t.Errorf("reviewLine with a web address that ends in a slash = %q", line)
	}
}

// U2. The first call edits once, the second finds the body equal and makes no call (U3 is its green negative).
func TestTheBodyIsEditedOnceToLinkTheReviewRecord(t *testing.T) {
	w := newLinkWorld(t)
	edited, err := w.link()
	if err != nil || !edited {
		t.Fatalf("LinkReviewRecord = %v, %v; want true, nil", edited, err)
	}
	if got, want := w.bodyOf(t, w.pr), wantBody(recA); got != want {
		t.Errorf("the body = %q; want %q (the snapshot of the developer that produced the head, not the verifier's %s)", got, want, recB)
	}
	if w.f.count("EditBody") != 1 {
		t.Errorf("calls %v; want one EditBody", w.f.calls)
	}
	edited, err = w.link()
	if err != nil || edited || w.f.count("EditBody") != 1 {
		t.Errorf("a second call = %v, %v with the calls %v; want false, nil and no second EditBody", edited, err, w.f.calls)
	}
}

// U3. A body that is already the composed one is read and not edited.
func TestAnEqualBodyIsNotEdited(t *testing.T) {
	w := newLinkWorld(t)
	if _, err := w.f.Fake.EditBody(context.Background(), w.pr, wantBody(recA)); err != nil {
		t.Fatal(err)
	}
	edited, err := w.link()
	if err != nil || edited || strings.Join(w.f.calls, " ") != "PullRequest" {
		t.Errorf("LinkReviewRecord = %v, %v with the calls %v; want false, nil and one PullRequest", edited, err, w.f.calls)
	}
}

// U4. A close-out that binds the same head later, at the snapshot recC, does not change the link; a restart (new Sessions on
// the same records) makes no edit; a body that a person edited is replaced whole.
func TestTheRecordsLinkIsTheProducingSessionsNotTheCloseOutsNorTheVerifiers(t *testing.T) {
	w := newLinkWorld(t)
	if _, err := w.link(); err != nil {
		t.Fatal(err)
	}
	closeOut := sessionRow(lkClose, "T-ab12", "close-out", "claude")
	closeOut[sessionBase], closeOut[sessionRecords] = head1, strings.Repeat("c", 40)
	w.sessions = append(w.sessions, closeOut)
	w.events = append(w.events,
		[]string{"7", "session", "1", lkClose, "", "", "", "2026-10-18T11:00:00Z"},
		[]string{"8", "bound", "1", lkClose, "", head1, "task/T-ab12/1", "2026-10-18T11:00:00Z"})
	w.write(t)
	w.s = newSessions(w.st)
	w.s.Forge, w.s.Web = w.f, prWeb
	if edited, err := w.link(); err != nil || edited {
		t.Errorf("after the close-out and a restart: %v, %v; want false, nil (the link kept the snapshot of the producing session)", edited, err)
	}
	if _, err := w.f.Fake.EditBody(context.Background(), w.pr, "a person's text"); err != nil {
		t.Fatal(err)
	}
	if edited, err := w.link(); err != nil || !edited || w.bodyOf(t, w.pr) != wantBody(recA) {
		t.Errorf("after a person's edit: %v, %v, body %q; want the whole body replaced", edited, err, w.bodyOf(t, w.pr))
	}
}

// A linkRefusal breaks one token of the world; fix repairs it, and the call then edits.
type linkRefusal struct {
	name     string
	brk, fix func(w *linkWorld)
	want     string
}

// U5. Each refusal of the records is a Go error before any forge call, and its one-token control passes.
func TestEachRecordThatFailsRefusesTheLinkBeforeAnyCall(t *testing.T) {
	var rows []linkRefusal
	add := func(name, want string, brk, fix func(w *linkWorld)) {
		rows = append(rows, linkRefusal{name, brk, fix, want})
	}
	add("the session is of another task", "not of T-ab12",
		func(w *linkWorld) { w.sessions[1][sessionTask] = "T-zz99" }, func(w *linkWorld) { w.sessions[1][sessionTask] = "T-ab12" })
	add("the session is a developer's", "not of the verification",
		func(w *linkWorld) { w.sessions[1][sessionRole] = roleOfStep("implement") },
		func(w *linkWorld) { w.sessions[1][sessionRole] = roleOfStep("verification") })
	add("a newer attempt exists", "newest attempt",
		func(w *linkWorld) {
			w.events = append(w.events, []string{"7", "attempt", "2", "", head1, "", "", "2026-10-18T10:30:00Z"})
		},
		func(w *linkWorld) { w.events = w.events[:6] })
	add("the session has no handoff", "has no handoff",
		func(w *linkWorld) { w.events = w.events[:5] },
		func(w *linkWorld) {
			w.events = append(w.events, []string{"6", "handoff", "1", lkVer, "", "", "valid comment 601-1", "2026-10-18T10:00:00Z"})
		})
	add("the latest handoff is invalid after an earlier valid one", "latest handoff",
		func(w *linkWorld) {
			w.events = append(w.events, []string{"7", "handoff", "1", lkVer, "", "", "invalid nothing-material comment 601-1", "2026-10-18T10:20:00Z"})
		},
		func(w *linkWorld) { w.events = w.events[:6] })
	add("the latest handoff names no comment", "latest handoff",
		func(w *linkWorld) { w.events[5][6] = "valid" },
		func(w *linkWorld) { w.events[5][6] = "valid comment 601-1" })
	add("the copy is not in copies.tsv", "no copy of the comment",
		func(w *linkWorld) { w.events[5][6] = "valid comment 777-1" },
		func(w *linkWorld) { w.events[5][6] = "valid comment 601-1" })
	add("the copy is of another issue", "no copy of the comment",
		func(w *linkWorld) { w.copies[0][copyIssue] = "9" },
		func(w *linkWorld) { w.copies[0][copyIssue] = "3" })
	add("the copy is a plan review", "not a review record",
		func(w *linkWorld) { w.copyText = "verdict: approve\n" },
		func(w *linkWorld) { w.copyText = lkRecord })
	add("the record is of another head", "is of the head",
		func(w *linkWorld) { w.copyText = record(shaOld, "nothing material in scope") },
		func(w *linkWorld) { w.copyText = lkRecord })
	add("the register names no pull request", "no pull request in the register",
		func(w *linkWorld) { w.pr = 0 },
		func(w *linkWorld) { w.pr = w.pr1.Number })
	add("no session bound the reviewed head", "no event bound at the head",
		func(w *linkWorld) { w.events[2][5], w.events[4][5] = shaOld, shaOld },
		func(w *linkWorld) { w.events[2][5], w.events[4][5] = head1, head1 })
	for _, r := range rows {
		w := newLinkWorld(t)
		r.brk(w)
		w.write(t)
		edited, err := w.link()
		wantErr(t, r.name, err, r.want)
		if edited || len(w.f.calls) != 0 {
			t.Errorf("%s: edited %v, calls %v; want no edit and no call of the forge", r.name, edited, w.f.calls)
		}
		r.fix(w)
		w.write(t)
		if edited, err := w.link(); err != nil || !edited {
			t.Errorf("%s: the control with the one token repaired = %v, %v; want true, nil", r.name, edited, err)
		}
	}
}

// U6. A lease that is not held edits nothing.
func TestALostLeaseEditsNothing(t *testing.T) {
	w := newLinkWorld(t)
	w.st.holder = "ffffffffffffffff"
	edited, err := w.link()
	var lost *LostError
	if !errors.As(err, &lost) || edited || w.f.count("EditBody") != 0 {
		t.Errorf("LinkReviewRecord = %v, %v with the calls %v; want *LostError and no EditBody", edited, err, w.f.calls)
	}
}

// U7. A closed or a merged pull request is not edited.
func TestAClosedOrMergedPullRequestIsNotEdited(t *testing.T) {
	w := newLinkWorld(t)
	if _, err := w.f.Fake.ClosePullRequest(context.Background(), w.pr); err != nil {
		t.Fatal(err)
	}
	edited, err := w.link()
	wantErr(t, "a closed pull request", err, "is closed")
	if edited || w.f.count("EditBody") != 0 {
		t.Errorf("a closed pull request: edited %v, calls %v", edited, w.f.calls)
	}
	w = newLinkWorld(t)
	w.f.Fake.MarkReady(context.Background(), w.pr1.NodeID)
	if _, err := w.f.Fake.Merge(context.Background(), w.pr, head1); err != nil {
		t.Fatal(err)
	}
	edited, err = w.link()
	wantErr(t, "a merged pull request", err, "is closed")
	if edited || w.f.count("EditBody") != 0 {
		t.Errorf("a merged pull request: edited %v, calls %v", edited, w.f.calls)
	}
}

// U8. A failed edit is an error and the next call edits once; an edit that took effect and gave an error is read by the next
// call, which finds the body equal and makes no second edit (the review's condition 4).
func TestAFailedEditIsAnErrorAndAnEditThatTookEffectIsNotRepeated(t *testing.T) {
	w := newLinkWorld(t)
	w.f.Fake.Fail = "EditBody"
	edited, err := w.link()
	wantErr(t, "a failed edit", err, "EditBody: status 502")
	if edited || w.bodyOf(t, w.pr) == wantBody(recA) {
		t.Errorf("a failed edit: edited %v, body %q", edited, w.bodyOf(t, w.pr))
	}
	w.f.Fake.Fail = ""
	if edited, err := w.link(); err != nil || !edited || w.f.count("EditBody") != 2 {
		t.Errorf("the next call = %v, %v with the calls %v; want one edit", edited, err, w.f.calls)
	}

	w = newLinkWorld(t)
	w.f.Fake.FailAfter = "EditBody"
	edited, err = w.link()
	wantErr(t, "an edit that took effect", err, "EditBody: status 502")
	if edited || w.bodyOf(t, w.pr) != wantBody(recA) {
		t.Fatalf("an edit that took effect: edited %v, body %q; want the body edited", edited, w.bodyOf(t, w.pr))
	}
	w.f.Fake.FailAfter = ""
	if edited, err := w.link(); err != nil || edited || w.f.count("EditBody") != 1 {
		t.Errorf("the next call = %v, %v with the calls %v; want false, nil and no second EditBody", edited, err, w.f.calls)
	}
}

// U9. A failed read of the pull request edits nothing.
func TestAFailedReadOfThePullRequestEditsNothing(t *testing.T) {
	w := newLinkWorld(t)
	w.f.Fake.Fail = "PullRequest"
	edited, err := w.link()
	wantErr(t, "a failed read", err, "PullRequest: status 502")
	if edited || w.f.count("EditBody") != 0 {
		t.Errorf("edited %v, calls %v", edited, w.f.calls)
	}
}

// U10. Each read of the store that fails is an error, with the earlier reads succeeding and no call of the forge.
func TestAFailedStoreReadIsAnErrorWithNoForgeCall(t *testing.T) {
	boom := errors.New("the store failed")
	for _, path := range []string{tasksPath, "sessions.tsv", eventsPath("T-ab12"), "copies.tsv", "copies/601-1.md"} {
		w := newLinkWorld(t)
		w.s.Store = &failingStore{standinStore: w.st, path: path, err: boom}
		if edited, err := w.link(); !errors.Is(err, boom) || edited || len(w.f.calls) != 0 {
			t.Errorf("a failed read of %s: %v, %v with the calls %v; want the store's error and no call", path, edited, err, w.f.calls)
		}
	}
	w := newLinkWorld(t)
	w.s.Store = &failingStore{standinStore: w.st, path: "lease", err: boom}
	if edited, err := w.link(); !errors.Is(err, boom) || edited || w.f.count("EditBody") != 0 {
		t.Errorf("a failed read of the lease: %v, %v with the calls %v; want the store's error and no edit", edited, err, w.f.calls)
	}
	w = newLinkWorld(t)
	delete(w.st.files, "copies/601-1.md")
	if _, err := w.link(); err == nil || len(w.f.calls) != 0 {
		t.Errorf("a copy whose file is missing: err %v, calls %v; want an error and no call", err, w.f.calls)
	}
}

// U11. An answer with another body is an error.
func TestAnAnswerWithAnotherBodyIsAnError(t *testing.T) {
	w := newLinkWorld(t)
	w.f.answer = func(pr forge.PullRequest) forge.PullRequest { pr.Body = "x"; return pr }
	edited, err := w.link()
	wantErr(t, "another body", err, "gave another body")
	if edited {
		t.Errorf("edited true with another body in the answer")
	}
}

// U12. The window of a newer attempt: the verifier is of attempt 2 while the register still names the pull request of attempt 1.
// The link edits neither pull request, whether the forge has a pull request of the attempt's branch or not; with the register on
// the pull request of the branch it edits that one alone (the review's condition 3).
func TestANewerAttemptWithTheRegisterOnTheOlderPullRequestEditsNothing(t *testing.T) {
	w := newLinkWorld(t)
	dev := sessionRow("S-d0000002", "T-ab12", "implement", "claude")
	dev[sessionAttempt], dev[sessionBase], dev[sessionRecords] = "2", head1, recB
	w.sessions[1][sessionAttempt], w.sessions[1][sessionBase], w.sessions[1][sessionRecords] = "2", head2, recB
	w.sessions = append(w.sessions, dev)
	at := "2026-10-18T12:00:00Z"
	w.events = [][]string{
		{"1", "attempt", "1", "", shaOld, "", "", at},
		{"2", "session", "1", lkDev, "", "", "", at},
		{"3", "bound", "1", lkDev, "", head1, "task/T-ab12/1", at},
		{"4", "attempt", "2", "", head1, "", "", at},
		{"5", "session", "2", "S-d0000002", "", "", "", at},
		{"6", "bound", "2", "S-d0000002", "", head2, "task/T-ab12/2", at},
		{"7", "session", "2", lkVer, "", "", "", at},
		{"8", "bound", "2", lkVer, "", head2, "task/T-ab12/2", at},
		{"9", "handoff", "2", lkVer, "", "", "valid comment 601-1", at},
	}
	w.copyText = record(head2, "nothing material in scope")
	w.copies[0][copySHA] = commentSum(w.copyText)
	w.write(t)
	// No pull request of the branch of attempt 2 yet: the register names the older one.
	edited, err := w.link()
	wantErr(t, "no pull request of the branch", err, "no pull request")
	if edited || w.f.count("EditBody") != 0 {
		t.Errorf("no pull request of attempt 2: edited %v, calls %v", edited, w.f.calls)
	}
	// The pull request of attempt 2 exists, and the register still names attempt 1's.
	w.f.Fake.SetHead("task/T-ab12/2", head2)
	pr2, err := w.f.Fake.OpenDraft(context.Background(), "task/T-ab12/2", "main", "feat: y", "Closes #3")
	if err != nil {
		t.Fatal(err)
	}
	edited, err = w.link()
	wantErr(t, "the register names the older pull request", err, "not the register's")
	if edited || w.f.count("EditBody") != 0 || w.bodyOf(t, w.pr1.Number) != w.pr1.Body {
		t.Errorf("the older pull request: edited %v, calls %v", edited, w.f.calls)
	}
	// The register names the pull request of the branch: that one is edited, the older one is not.
	w.pr = pr2.Number
	w.write(t)
	if edited, err := w.link(); err != nil || !edited {
		t.Fatalf("with the register on %d: %v, %v; want true, nil", pr2.Number, edited, err)
	}
	if got := w.bodyOf(t, pr2.Number); !strings.HasPrefix(got, "Closes #3\n\nThe records commit this task started from: "+prWeb+"/acme/target/commit/"+recB) || !strings.HasSuffix(got, "\n\nThe review record: "+lkURL) {
		t.Errorf("the body of the pull request %d = %q", pr2.Number, got)
	}
	if w.bodyOf(t, w.pr1.Number) != w.pr1.Body {
		t.Errorf("the pull request of attempt 1 was edited: %q", w.bodyOf(t, w.pr1.Number))
	}
}

// U13. A head of the pull request other than the reviewed one is an error.
func TestAPullRequestWhoseHeadMovedIsNotEdited(t *testing.T) {
	w := newLinkWorld(t)
	w.f.Fake.MoveHead(w.pr, head3)
	edited, err := w.link()
	wantErr(t, "a moved head", err, "not the reviewed")
	if edited || w.f.count("EditBody") != 0 {
		t.Errorf("edited %v, calls %v", edited, w.f.calls)
	}
}
