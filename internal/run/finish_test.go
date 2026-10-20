package run

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// The close-out and the merge of a verified task, at the unit level (task T-ehb4, #425, row 66y; loop.md, The merge at the head SHA and
// the close-out commit, D1 to D10): the positions on the stand-in store and the stand-in forge, and the pure functions.

const (
	carriedC   = "cccccccccccccccccccccccccccccccccccccccc"
	closeID    = "S-c0000001"
	verifierID = "S-v0000001"
)

// finishWorld is the world of the pull request, with the pull request of attempt 2 open and registered.
func finishWorld(t *testing.T) (*prWorld, forge.PullRequest) {
	t.Helper()
	w := newPRWorld(t)
	return w, w.open(t)
}

func finish(w *prWorld) (Advance, error) {
	return w.s.FinishTask(context.Background(), FinishInput{Task: "T-ab12", Records: func() string { return recB }})
}

// U1 (D2, D7): a pull request closed and not merged is a wait whose text names the pull request and the branch, and it writes nothing: no
// records commit, no read of a branch for a repair, even with an unresolved push of a session that has no comment; the forge is only read.
func TestAClosedUnmergedPullRequestOfTheOpenAttemptIsAWaitWithNoWrite(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	w.events = append(w.events, []string{"7", "push", "2", "S-d0000002", "", head2, "task/T-ab12/2", "2026-10-18T10:05:00Z"})
	w.write(t)
	if _, err := w.f.Fake.ClosePullRequest(context.Background(), pr.Number); err != nil {
		t.Fatal(err)
	}
	commits, ops := len(w.st.commits), len(w.st.ops)
	want := "waiting for the Operator: T-ab12: the pull request #" + strconv.Itoa(pr.Number) + " of task/T-ab12/2 is closed and not merged"
	got, err := finish(w)
	if err != nil || got != (Advance{Task: "T-ab12", Wait: want}) {
		t.Fatalf("FinishTask = %+v, %v; want the wait %q", got, err, want)
	}
	if len(w.st.commits) != commits || len(w.st.ops) != ops {
		t.Errorf("%d commits and %d operations of the store; want none beyond the %d and %d before: no repair, no write", len(w.st.commits), len(w.st.ops), commits, ops)
	}
	for _, l := range w.f.log {
		if l != "PullRequest" {
			t.Errorf("the forge was asked %q; want reads of the pull request only", l)
		}
	}
}

// U2 (O16): a merged pull request whose merged event was lost gets the event, once, for the head it was merged at.
func TestAMergedPullRequestWithNoEventGetsItsMergedEvent(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	if _, err := w.f.Fake.MarkReady(context.Background(), pr.NodeID); err != nil {
		t.Fatal(err)
	}
	merge, err := w.f.Fake.Merge(context.Background(), pr.Number, head2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := finish(w)
	if err != nil || got != (Advance{Task: "T-ab12", Merged: true, Head: head2}) {
		t.Fatalf("FinishTask = %+v, %v; want Merged at %s", got, err, head2)
	}
	last := events(t, w.st, "T-ab12")
	if e := last[len(last)-1]; e[1] != "merged" || e[2] != "2" || e[5] != merge {
		t.Errorf("the last event %v; want merged %s in attempt 2", e, merge)
	}
}

// U3: a merged pull request with its event writes nothing and reads nothing for a repair, though a push of a session is unresolved; U2 is its control.
func TestAMergedPullRequestWithItsEventWritesNothing(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	w.green(head2)
	if _, r, err := w.s.MergeTask(context.Background(), "T-ab12", head2); err != nil || !r.Ready {
		t.Fatalf("MergeTask = %+v, %v", r, err)
	}
	w.events = append(events(t, w.st, "T-ab12"), []string{"8", "push", "2", "S-d0000002", "", head2, "task/T-ab12/2", "2026-10-18T10:05:00Z"})
	w.write(t)
	commits, ops := len(w.st.commits), len(w.st.ops)
	got, err := finish(w)
	if err != nil || got != (Advance{Task: "T-ab12", Merged: true, Head: head2}) || pr.Number != 1 {
		t.Fatalf("FinishTask = %+v, %v; want Merged at %s", got, err, head2)
	}
	if len(w.st.commits) != commits || len(w.st.ops) != ops {
		t.Errorf("%d commits and %d operations after; want none beyond the %d and %d before", len(w.st.commits), len(w.st.ops), commits, ops)
	}
}

// U7 (D8): the pull request of the branch of the open attempt must be the register's.
func TestThePullRequestOfTheOpenAttemptMustBeTheRegisters(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	w.register(t, pr.Number+4)
	want := fmt.Sprintf("T-ab12: the pull request %d of task/T-ab12/2 is not the register's %d", pr.Number, pr.Number+4)
	if got, err := finish(w); err == nil || err.Error() != want {
		t.Errorf("FinishTask = %+v, %v; want the error %q", got, err, want)
	}
	w.register(t, 0)
	if _, err := finish(w); err == nil {
		t.Error("a register with no pull request: FinishTask gave no error")
	}
	if _, err := finish(newPRWorld(t)); err == nil || !strings.Contains(err.Error(), "task/T-ab12/2 is not on the forge") {
		t.Errorf("FinishTask with no pull request on the branch = %v; want the error that it is not on the forge", err)
	}
}

// U10 (D1, position 5): an attempt with no handoff of a verifier has no verified head, and nothing is written.
func TestAnAttemptWithNoVerifierHandoffHasNoVerifiedHead(t *testing.T) {
	t.Parallel()
	w, _ := finishWorld(t)
	commits := len(w.st.commits)
	if _, err := finish(w); err == nil || err.Error() != "T-ab12: attempt 2 has no verified head" || len(w.st.commits) != commits {
		t.Errorf("FinishTask = %v with %d commits more; want the error that attempt 2 has no verified head and no write", err, len(w.st.commits)-commits)
	}
}

// U4 (D4, condition 8 of #367): the close-out is resolved at the tier of the plan moved by the position.
func TestTheCloseOutStepIsResolvedAtTierOfThePlanAndThePosition(t *testing.T) {
	t.Parallel()
	large, small := Plan{Size: "large"}, Plan{Size: "small"}
	for _, c := range []struct {
		name string
		plan Plan
		pos  Position
		tier string
	}{
		{"large and downhill", large, Downhill, "execution"},
		{"small and uphill", small, Uphill, "reasoning"},
		{"large and no position", large, "", "reasoning"},
		{"small and no position", small, "", "execution"},
	} {
		got, err := closeOutStep(c.plan, c.pos)
		if want := (Resolved{Step: "close-out", Role: "software-developer", Tier: c.tier}); err != nil || got != want {
			t.Errorf("%s: closeOutStep = %+v, %v; want %+v", c.name, got, err, want)
		}
	}
}

// view builds the view of a task whose events are given as kind, attempt, session, sha and detail.
func view(rows [][]string, sessions [][]string, status map[string]string) loopView {
	v := loopView{Task: "T-ab12", Kind: "architecture", Sessions: sessions, Status: status}
	for i, r := range rows {
		v.Events = append(v.Events, []string{strconv.Itoa(i + 1), r[0], r[1], r[2], "", r[3], r[4], "2026-10-20T10:00:00Z"})
	}
	return v
}

func closer(id, base string) []string {
	r := sessionRow(id, "T-ab12", "close-out", "claude")
	r[sessionAttempt], r[sessionBase] = "1", base
	return r
}

// U9: the cut is the last handoff of a session of the verification role.
func TestTheCutIsTheLastHandoffOfASessionOfTheVerificationRole(t *testing.T) {
	t.Parallel()
	qa, dev := sessionRow(verifierID, "T-ab12", "verification", "devin"), sessionRow("S-d0000001", "T-ab12", "implement", "claude")
	v := view([][]string{{"attempt", "1", "", "", ""}, {"handoff", "1", "S-d0000001", "", "valid"}, {"handoff", "1", verifierID, "", "valid comment 5-1"}, {"handoff", "1", "S-d0000001", "", "valid"}},
		[][]string{qa, dev}, nil)
	if got := finishCut(v); got != 2 {
		t.Errorf("finishCut = %d; want 2, the handoff of the verifier", got)
	}
	if got := finishCut(view(nil, nil, nil)); got != -1 {
		t.Errorf("finishCut of no event = %d; want -1", got)
	}
}

// U5 (D5, condition 3): the close-out session of the position is the last one after the cut at the carried base, told by the events.
func TestTheCloseOutSessionOfThePositionIsTheLastOneAfterTheCutAtTheCarriedBase(t *testing.T) {
	t.Parallel()
	head := []string{"attempt", "1", "", "", ""}
	cut := []string{"handoff", "1", verifierID, "", "valid comment 5-1"}
	qa := sessionRow(verifierID, "T-ab12", "verification", "devin")
	start := []string{"session", "1", closeID, "", ""}
	done := []string{"result", "1", closeID, "", "done"}
	pushed := []string{"push", "1", closeID, aSHA, "task/T-ab12/1"}
	bound := []string{"bound", "1", closeID, aSHA, "task/T-ab12/1"}
	completed := map[string]string{closeID: "completed"}
	for _, c := range []struct {
		name   string
		events [][]string
		rows   [][]string
		status map[string]string
		want   closeStep
	}{
		{"none", [][]string{head, cut}, nil, nil, closeStep{Pos: closeStart}},
		{"running", [][]string{head, cut, start}, [][]string{closer(closeID, carriedC)}, nil, closeStep{Pos: closeStart}},
		{"ended with no push", [][]string{head, cut, start, done}, [][]string{closer(closeID, carriedC)}, completed, closeStep{Pos: closeResume, Session: closeID}},
		{"pushed and bound", [][]string{head, cut, start, done, pushed, bound}, [][]string{closer(closeID, carriedC)}, completed, closeStep{Pos: closeBound, Session: closeID, Head: aSHA}},
		{"refused", [][]string{head, cut, start, done, {"refused", "1", closeID, "", "scope abc"}}, [][]string{closer(closeID, carriedC)}, completed,
			closeStep{Pos: closeStop, Session: closeID, Msg: "T-ab12: the session " + closeID + " was refused (scope abc)"}},
		{"pushed with no bound", [][]string{head, cut, start, done, pushed}, [][]string{closer(closeID, carriedC)}, completed,
			closeStep{Pos: closeStop, Session: closeID, Msg: "T-ab12: the push of the session " + closeID + " has no result (its repair is RepairTask's)"}},
		{"ended failed", [][]string{head, cut, start, {"result", "1", closeID, "", "failed"}}, [][]string{closer(closeID, carriedC)}, nil,
			closeStep{Pos: closeStop, Session: closeID, Msg: "T-ab12: the session " + closeID + " of close-out ended failed"}},
		{"another base", [][]string{head, cut, start, done, pushed, bound}, [][]string{closer(closeID, aSHA)}, completed, closeStep{Pos: closeStart}},
		{"started before the cut", [][]string{head, start, cut, done, pushed, bound}, [][]string{closer(closeID, carriedC)}, completed, closeStep{Pos: closeStart}},
		{"the last decides", [][]string{head, cut, start, done, pushed, bound, {"session", "1", "S-c0000002", "", ""}}, [][]string{closer(closeID, carriedC), closer("S-c0000002", carriedC)}, completed,
			closeStep{Pos: closeStart}},
		{"a start refused before its row", [][]string{head, cut, {"refused", "1", "S-c0000003", "", "budget"}}, nil, nil,
			closeStep{Pos: closeStop, Session: "S-c0000003", Msg: "T-ab12: the start of the session S-c0000003 was refused (budget)"}},
	} {
		v := view(c.events, append([][]string{qa}, c.rows...), c.status)
		if got := closeOutPosition(v, indexOf(c.events, cut), 1, carriedC); got != c.want {
			t.Errorf("%s: closeOutPosition = %+v; want %+v", c.name, got, c.want)
		}
	}
}

func indexOf(rows [][]string, row []string) int {
	for i, r := range rows {
		if fmt.Sprint(r) == fmt.Sprint(row) {
			return i
		}
	}
	return -1
}

// U6 (D6): not ready is a wait that names the reason, a moved head is a stop, and a ready head is the merge.
func TestNotReadyIsAWaitAndAMovedHeadIsAStop(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		r    Readiness
		want Advance
	}{
		{Readiness{Reason: "layup/gates: pending"}, Advance{Task: "T-ab12", Wait: "waiting: T-ab12: not ready at " + head2 + ": layup/gates: pending"}},
		{Readiness{Reason: "job build: in_progress"}, Advance{Task: "T-ab12", Wait: "waiting: T-ab12: not ready at " + head2 + ": job build: in_progress"}},
		{Readiness{Reason: "the head moved to " + head3}, Advance{Task: "T-ab12", Stop: "T-ab12: not ready at " + head2 + ": the head moved to " + head3}},
		{Readiness{Ready: true}, Advance{Task: "T-ab12", Merged: true, Head: head2}},
	} {
		if got := notReady("T-ab12", head2, c.r); got != c.want {
			t.Errorf("notReady(%+v) = %+v; want %+v", c.r, got, c.want)
		}
	}
	if got, want := closedWait("T-ab12", "task/T-ab12/2", 7), "waiting for the Operator: T-ab12: the pull request #7 of task/T-ab12/2 is closed and not merged"; got != want {
		t.Errorf("closedWait = %q; want %q", got, want)
	}
}

// U8 (D9): the default branch's head equal to the base is up to date with no git read; a head that is not on the forge and a failed read are errors,
// never a moved base.
func TestTheDefaultHeadMustBeTheBaseOrAnAncestorOfTheCarriedHead(t *testing.T) {
	t.Parallel()
	w, _ := finishWorld(t)
	w.st.heads = map[string]string{"main": shaOld}
	if ok, err := w.s.upToDate(context.Background(), shaOld, head2); err != nil || !ok {
		t.Errorf("upToDate with the tip at the base = %v, %v; want true", ok, err)
	}
	w.st.heads = nil
	if ok, err := w.s.upToDate(context.Background(), shaOld, head2); err == nil || ok {
		t.Errorf("upToDate with no default branch = %v, %v; want an error", ok, err)
	}
	w.st.heads, w.st.headErr = map[string]string{"main": shaOld}, errors.New("the forge did not answer")
	if ok, err := w.s.upToDate(context.Background(), shaOld, head2); err == nil || ok {
		t.Errorf("upToDate with a failed read = %v, %v; want an error", ok, err)
	}
}

// U11 (finding 3 of round 1): a pull request merged outside, at the developer's head, while a close-out session's push of another head is unresolved, gets its
// event merged for the attempt of that head, and the repair that resolves the push starts no attempt; a second call writes nothing.
func TestAMergedPullRequestIsRecordedAndStartsNoReplacementAttempt(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	rows, err := records.ReadSessions(w.st.files["sessions.tsv"])
	if err != nil {
		t.Fatal(err)
	}
	closer := sessionRow("S-c0000001", "T-ab12", "close-out", "claude")
	closer[sessionAttempt], closer[sessionBase], closer[sessionRecords] = "2", head2, recB
	if w.st.files["sessions.tsv"], err = table(records.SessionsSchema, append(rows, closer)); err != nil {
		t.Fatal(err)
	}
	at := "2026-10-18T10:05:00Z"
	w.events = append(w.events,
		[]string{"7", "session", "2", "S-c0000001", "", "", "", at},
		[]string{"8", "result", "2", "S-c0000001", "", "", "done", at},
		[]string{"9", "push", "2", "S-c0000001", "", head3, "task/T-ab12/2", at})
	w.write(t)
	w.st.heads, w.s.App = map[string]string{"task/T-ab12/2": head2}, "layup-agent"
	if _, err := w.f.Fake.MarkReady(context.Background(), pr.NodeID); err != nil {
		t.Fatal(err)
	}
	merge, err := w.f.Fake.Merge(context.Background(), pr.Number, head2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := finish(w)
	if err != nil || got != (Advance{Task: "T-ab12", Merged: true, Head: head2}) {
		t.Fatalf("FinishTask = %+v, %v; want Merged at %s with no new attempt", got, err, head2)
	}
	n, merged := 0, 0
	for _, e := range events(t, w.st, "T-ab12") {
		switch {
		case e[1] == "attempt":
			n++
		case e[1] == "merged" && e[2] == "2" && e[5] == merge:
			merged++
		}
	}
	if n != 2 || merged != 1 {
		t.Errorf("%d attempts and %d events merged of attempt 2; want 2 and 1", n, merged)
	}
	commits := len(w.st.commits)
	if again, err := finish(w); err != nil || !again.Merged || len(w.st.commits) != commits {
		t.Errorf("a second call = %+v, %v with %d commits more; want Merged and nothing written", again, err, len(w.st.commits)-commits)
	}
}

// U12 (finding 3): a merged pull request whose head was pushed by a session with no bound yet gets that bound, and then its event merged, still with no
// replacement attempt: the resolution of the pushes establishes the binding of the merged head.
func TestAMergedPullRequestWhosePushWasNeverBoundGetsItsBoundAndItsMergedEvent(t *testing.T) {
	t.Parallel()
	w, pr := finishWorld(t)
	w.events[5] = []string{"6", "push", "2", "S-d0000002", "", head2, "task/T-ab12/2", "2026-10-18T10:00:00Z"}
	w.write(t)
	w.st.heads = map[string]string{"task/T-ab12/2": head2}
	if _, err := w.f.Fake.MarkReady(context.Background(), pr.NodeID); err != nil {
		t.Fatal(err)
	}
	merge, err := w.f.Fake.Merge(context.Background(), pr.Number, head2)
	if err != nil {
		t.Fatal(err)
	}
	got, err := finish(w)
	if err != nil || got != (Advance{Task: "T-ab12", Merged: true, Head: head2}) {
		t.Fatalf("FinishTask = %+v, %v; want Merged at %s", got, err, head2)
	}
	kinds := ""
	for _, e := range events(t, w.st, "T-ab12")[5:] {
		kinds += e[1] + ":" + e[5] + " "
	}
	if want := "push:" + head2 + " bound:" + head2 + " merged:" + merge + " "; kinds != want {
		t.Errorf("the events from the sixth: %q; want %q: the bound, then the merged, and no attempt", kinds, want)
	}
}
