package run

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// The records of a restart (task T-syxm): the first attempt of T-ab12 at base, a
// session that was pushed and has no bound and no refused.
const (
	pushSess  = "S-1a2b3c4d"
	laterSess = "S-5e6f7a8b"
	baseSHA   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	pushedSHA = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	movedSHA  = "cccccccccccccccccccccccccccccccccccccccc"
	taskOne   = "task/T-ab12/1"
)

// ev gives an event row; its number is set by evs.
func evRow(kind string, attempt int, session, base, sha, detail string) []string {
	return []string{"", kind, strconv.Itoa(attempt), session, base, sha, detail, "2026-10-09T12:00:00Z"}
}

func attemptEv(n int, base string) []string { return evRow("attempt", n, "", base, "", "") }

// pushed is the events of a session that was pushed: session, result done, push.
func pushed(id string, attempt int) [][]string {
	return [][]string{evRow("session", attempt, id, "", "", ""), evRow("result", attempt, id, "", "", "done"),
		evRow("push", attempt, id, "", pushedSHA, "task/T-ab12/"+strconv.Itoa(attempt))}
}

func evs(groups ...[][]string) [][]string {
	var all [][]string
	for _, g := range groups {
		all = append(all, g...)
	}
	for i := range all {
		all[i][0] = strconv.Itoa(i + 1)
	}
	return all
}

// repairWorld gives the store and the Sessions of a restart whose events.tsv of T-ab12 is rows and
// whose target holds heads.
func repairWorld(t *testing.T, rows [][]string, heads map[string]string) (*standinStore, *Sessions) {
	t.Helper()
	data, err := table(records.EventsSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	s := &standinStore{files: map[string][]byte{"tasks/T-ab12/events.tsv": data}, holder: "0123456789abcdef", heads: heads}
	return s, newSessions(s)
}

// kinds gives "kind[:detail]" of each event of T-ab12.
func kinds(t *testing.T, s *standinStore) []string {
	t.Helper()
	var out []string
	for _, e := range events(t, s, "T-ab12") {
		out = append(out, e[1])
	}
	return out
}

func reads(s *standinStore) int {
	n := 0
	for _, op := range s.ops {
		if strings.HasPrefix(op, "read ") {
			n++
		}
	}
	return n
}

func count(list []string, kind string) int {
	n := 0
	for _, k := range list {
		if k == kind {
			n++
		}
	}
	return n
}

func TestAPushWhoseBranchHoldsItsSHAIsBound(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1)), map[string]string{taskOne: pushedSHA})
	got, err := r.RepairPushes(context.Background(), "T-ab12")
	want := []string{"attempt", "session", "result", "push", "bound"}
	if k := kinds(t, s); err != nil || !slices.Equal(k, want) {
		t.Fatalf("events = %v, %v; want %v", k, err, want)
	}
	last := events(t, s, "T-ab12")[4]
	if last[3] != pushSess || last[5] != pushedSHA || last[6] != taskOne || !slices.Equal(got.Bound, []string{pushSess}) || got.Attempt != 0 {
		t.Errorf("the bound event %q, Repair %+v; want session %s, the SHA and the branch, and no attempt", last, got, pushSess)
	}
}

func TestAPushWhoseBranchIsAbsentIsPushLostAndStartsOneAttempt(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1)), nil)
	got, err := r.RepairPushes(context.Background(), "T-ab12")
	want := []string{"attempt", "session", "result", "push", "refused", "attempt"}
	if k := kinds(t, s); err != nil || !slices.Equal(k, want) {
		t.Fatalf("events = %v, %v; want %v", k, err, want)
	}
	if got.Attempt != 2 || !slices.Equal(got.Lost, []string{pushSess}) {
		t.Errorf("Repair.Attempt = %d, Lost %v; want 2 and [%s]", got.Attempt, got.Lost, pushSess)
	}
	if e := events(t, s, "T-ab12"); e[4][3] != pushSess || e[4][6] != "push-lost" || e[5][2] != "2" {
		t.Errorf("the events %q; want refused push-lost of %s, then attempt 2", e[4:], pushSess)
	}
}

// Being the only writer proves no equality (Q7): another head, a descendant of the
// pushed SHA included (the integration test I5 pushes a real one), is a lost push.
func TestABranchHeadThatDiffersFromThePushIsPushLost(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1)), map[string]string{taskOne: movedSHA})
	if _, err := r.RepairPushes(context.Background(), "T-ab12"); err != nil {
		t.Fatal(err)
	}
	if e := events(t, s, "T-ab12"); len(e) < 5 || e[4][1] != "refused" || e[4][6] != "push-lost" {
		t.Errorf("events %q; want the 5th event refused push-lost", e)
	}
}

func TestAFailedBranchReadIsAnErrorAndWritesNothing(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1)), nil)
	s.headErr = errors.New("the remote cannot be reached")
	if _, err := r.RepairPushes(context.Background(), "T-ab12"); !errors.Is(err, s.headErr) {
		t.Errorf("err = %v; want the read's error", err)
	}
	if k := kinds(t, s); len(s.commits) != 0 || len(k) != 4 {
		t.Errorf("%d commits, events %v; want none written", len(s.commits), k)
	}
}

// A push of an attempt that is no longer the latest is resolved and starts no
// attempt: attempt 2 is already the open one.
func TestAPushLostInAnOlderAttemptStartsNoAttempt(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{attemptEv(2, baseSHA)}), nil)
	if _, err := r.RepairPushes(context.Background(), "T-ab12"); err != nil {
		t.Fatal(err)
	}
	k := kinds(t, s)
	e := events(t, s, "T-ab12")
	if len(e) != 6 || e[5][1] != "refused" || e[5][2] != "1" || e[5][6] != "push-lost" || count(k, "attempt") != 2 {
		t.Errorf("events %q; want refused of attempt 1 as the 6th event and still two attempts", e)
	}
}

// A restart between the commit of refused and the commit of attempt: the records
// hold the refusal and no attempt after it; the repair starts exactly one, and a
// second call writes nothing.
func TestAnOpenAttemptHoldingPushLostGetsExactlyOneAttempt(t *testing.T) {
	rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{evRow("refused", 1, pushSess, "", "", "push-lost")})
	s, r := repairWorld(t, rows, nil)
	for call := 1; call <= 2; call++ {
		if _, err := r.RepairPushes(context.Background(), "T-ab12"); err != nil {
			t.Fatal(err)
		}
		if n := count(kinds(t, s), "attempt"); n != 2 {
			t.Errorf("after call %d: %d attempts; want 2", call, n)
		}
	}
	if len(s.commits) != 1 || reads(s) != 0 {
		t.Errorf("%d commits and %d reads; want one commit (the attempt) and no read", len(s.commits), reads(s))
	}
}

// Several pushes lost in the open attempt give one replacement.
func TestSeveralLostPushesOfTheOpenAttemptGiveOneAttempt(t *testing.T) {
	rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), pushed(laterSess, 1))
	s, r := repairWorld(t, rows, nil)
	got, err := r.RepairPushes(context.Background(), "T-ab12")
	k := kinds(t, s)
	if err != nil || count(k, "refused") != 2 || count(k, "attempt") != 2 || got.Attempt != 2 || len(got.Lost) != 2 {
		t.Errorf("events %v, Repair %+v, %v; want two refused and attempt 2 once", k, got, err)
	}
}

// The latest attempt may be closed or rebased without another attempt event; the
// refusal is written and no attempt starts (the existing meaning of an open attempt).
func TestAClosedOrRebasedLatestAttemptGetsNoReplacement(t *testing.T) {
	for _, closing := range [][]string{evRow("closed", 1, "", "", "", "T-cd34"), evRow("rebased", 1, "", pushedSHA, "", "T-cd34")} {
		t.Run(closing[1], func(t *testing.T) {
			rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{closing})
			s, r := repairWorld(t, rows, nil)
			got, err := r.RepairPushes(context.Background(), "T-ab12")
			k := kinds(t, s)
			if err != nil || count(k, "refused") != 1 || count(k, "attempt") != 1 || got.Attempt != 0 {
				t.Errorf("events %v, Repair %+v, %v; want the refusal and no new attempt", k, got, err)
			}
		})
	}
}

// A push that a bound or a refused follows is not read again.
func TestAPushThatIsAlreadyBoundOrRefusedIsNotReadAgain(t *testing.T) {
	for name, after := range map[string]string{"bound": "bound", "refused": "refused"} {
		t.Run(name, func(t *testing.T) {
			tail := evRow("bound", 1, pushSess, "", pushedSHA, taskOne)
			if after == "refused" {
				tail = evRow("refused", 1, pushSess, "", "", "push-refused")
			}
			s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{tail}), nil)
			got, err := r.RepairPushes(context.Background(), "T-ab12")
			if err != nil || reads(s) != 0 || len(s.commits) != 0 || got.Attempt != 0 {
				t.Errorf("%d reads, %d commits, Repair %+v, %v; want none", reads(s), len(s.commits), got, err)
			}
		})
	}
}

func TestTheNewAttemptStartsAtTheLostAttemptsBase(t *testing.T) {
	s, r := repairWorld(t, evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1)), nil)
	if _, err := r.RepairPushes(context.Background(), "T-ab12"); err != nil {
		t.Fatal(err)
	}
	if e := events(t, s, "T-ab12"); len(e) != 6 || e[5][1] != "attempt" || e[5][4] != baseSHA {
		t.Errorf("events %q; want attempt 2 at the base %s of attempt 1", e, baseSHA)
	}
}

// A replacement attempt that was not committed is not reported as written (round 1 of
// #422, finding 1): with both fenced commits refused, Repair.Attempt is 0 beside the
// error, and the retry commits attempt 2 once.
func TestAReplacementThatWasNotCommittedIsNotReportedAsWritten(t *testing.T) {
	rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{evRow("refused", 1, pushSess, "", "", "push-lost")})
	store, sessions := repairWorld(t, rows, nil)
	store.refuse = 2
	got, err := sessions.RepairPushes(context.Background(), "T-ab12")
	if err == nil || len(store.commits) != 0 || count(kinds(t, store), "attempt") != 1 {
		t.Fatalf("both commits refused: Repair %+v, %v, %d commits, events %v; want the error and no attempt", got, err, len(store.commits), kinds(t, store))
	}
	if got.Attempt != 0 {
		t.Errorf("Repair.Attempt = %d with no replacement committed; want 0", got.Attempt)
	}
	got, err = sessions.RepairPushes(context.Background(), "T-ab12")
	if err != nil || got.Attempt != 2 || count(kinds(t, store), "attempt") != 2 || len(store.commits) != 1 {
		t.Fatalf("the retry: Repair %+v, %v, events %v, %d commits; want attempt 2 committed once", got, err, kinds(t, store), len(store.commits))
	}
	if got, err = sessions.RepairPushes(context.Background(), "T-ab12"); err != nil || got.Attempt != 0 || len(store.commits) != 1 {
		t.Errorf("a third call: Repair %+v, %v, %d commits; want nothing more", got, err, len(store.commits))
	}
}

// A sessions.tsv that the reader refuses stops the repost and the repair of a task:
// the error names the path, and nothing is committed or posted (round 1 of #422,
// finding 3). The Sessions has no forge, so a post would panic. A green negative on
// the code as it stands: its red mutant is the read's error returned as nil, nil
// (runs/T-syxm/mutants.md, M23).
func TestARefusedSessionsTableStopsTheRepostAndTheRepair(t *testing.T) {
	rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1), [][]string{evRow("bound", 1, pushSess, "", pushedSHA, "task/T-ab12/1")})
	for _, c := range []struct {
		name string
		call func(*Sessions) ([]string, error)
	}{
		{"RepostComments", func(s *Sessions) ([]string, error) { return s.RepostComments(context.Background(), "T-ab12") }},
		{"RepairTask", func(s *Sessions) ([]string, error) {
			r, err := s.RepairTask(context.Background(), "T-ab12")
			return r.Posted, err
		}},
	} {
		t.Run(c.name, func(t *testing.T) {
			store, sessions := repairWorld(t, rows, nil)
			store.files["sessions.tsv"] = []byte("not a sessions table\n")
			posted, err := c.call(sessions)
			if err == nil || !strings.Contains(err.Error(), "sessions.tsv:") || len(posted) != 0 || len(store.commits) != 0 {
				t.Errorf("%s: posted %v, %v, %d commits; want an error that names sessions.tsv, nothing posted, no commit", c.name, posted, err, len(store.commits))
			}
		})
	}
}

// An events.tsv that the reader refuses stops RepairPushes and RepostComments: each
// returns the read's error, which names the path, and writes and reads nothing else
// (round 2 of #422, finding 1). A green negative on the code as it stands: its red
// mutants are the two events reads taken as empty (runs/T-syxm/mutants.md, M26 and
// M27).
func TestARefusedEventsTableStopsTheRepairAndTheRepost(t *testing.T) {
	rows := evs([][]string{attemptEv(1, baseSHA)}, pushed(pushSess, 1))
	for _, c := range []struct {
		name string
		call func(*Sessions) error
	}{
		{"RepairPushes", func(s *Sessions) error { _, err := s.RepairPushes(context.Background(), "T-ab12"); return err }},
		{"RepostComments", func(s *Sessions) error { _, err := s.RepostComments(context.Background(), "T-ab12"); return err }},
	} {
		t.Run(c.name, func(t *testing.T) {
			store, sessions := repairWorld(t, rows, map[string]string{taskOne: pushedSHA})
			store.files["tasks/T-ab12/events.tsv"] = []byte("not an events table\n")
			err := c.call(sessions)
			if err == nil || !strings.Contains(err.Error(), "events.tsv:") || len(store.commits) != 0 || len(store.ops) != 0 {
				t.Errorf("%s: %v, %d commits, %d branch reads; want an error that names events.tsv, no commit, no read", c.name, err, len(store.commits), len(store.ops))
			}
		})
	}
}
