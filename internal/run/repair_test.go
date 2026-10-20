package run

import (
	"bytes"
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
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

// The replacement of the row of a redone task (task T-89xf, row 66x; loop.md, A task whose base moved).
var (
	redoAt   = time.Date(2026, 10, 20, 9, 0, 0, 0, time.UTC)
	redoRows = func() []Task {
		return []Task{
			{ID: "T-spec", Issue: 2, Kind: "specification", Requirements: []string{"REQ-001"}, PR: 3, Opened: redoAt.Add(-time.Hour)},
			{ID: "T-old1", Issue: 5, Kind: "architecture", Requirements: []string{"REQ-001", "NFR-002"}, Size: "small", Predecessors: []string{"T-spec"}, PR: 9, Opened: redoAt.Add(-time.Hour)},
			{ID: "T-sub1", Issue: 6, Kind: "build", Requirements: []string{"REQ-002"}, Size: "large", Predecessors: []string{"T-old1", "T-spec"}, Opened: redoAt.Add(-time.Hour)},
			{ID: "T-sub2", Issue: 7, Kind: "build", Requirements: []string{"REQ-003"}, Size: "small", Predecessors: []string{"T-spec"}, Opened: redoAt.Add(-time.Hour)},
		}
	}
)

func TestTheReplacementTakesTheOldRowAndClearsItsPullRequest(t *testing.T) {
	before := redoRows()
	got, err := redoRegister(before, "T-old1", "T-new1", 12, redoAt)
	if err != nil || len(got) != 4 {
		t.Fatalf("redoRegister = %v, %v; want four rows", got, err)
	}
	want := Task{ID: "T-new1", Issue: 12, Kind: "architecture", Requirements: []string{"REQ-001", "NFR-002"}, Size: "small", Predecessors: []string{"T-spec"}, PR: 0, Opened: redoAt}
	if got[1].ID != want.ID || got[1].Issue != want.Issue || got[1].Kind != want.Kind || !slices.Equal(got[1].Requirements, want.Requirements) ||
		got[1].Size != want.Size || !slices.Equal(got[1].Predecessors, want.Predecessors) || got[1].PR != 0 || !got[1].Opened.Equal(want.Opened) {
		t.Errorf("the row in place of the old one = %+v; want %+v", got[1], want)
	}
	if got[0].ID != "T-spec" || got[2].ID != "T-sub1" || got[3].ID != "T-sub2" {
		t.Errorf("the order of the rows: %s %s %s %s; want the position kept", got[0].ID, got[1].ID, got[2].ID, got[3].ID)
	}
	if before[1].ID != "T-old1" || before[2].Predecessors[0] != "T-old1" {
		t.Error("the table given was changed")
	}
	for name, c := range map[string]struct {
		old, id string
	}{"an old task that is no row": {"T-zzzz", "T-new1"}, "a new ID that is a row": {"T-old1", "T-sub2"}, "a new ID that is the old": {"T-old1", "T-old1"}} {
		if rows, err := redoRegister(redoRows(), c.old, c.id, 12, redoAt); err == nil {
			t.Errorf("%s: %v, nil; want an error", name, rows)
		}
	}
}

func TestEachTaskThatNamesTheOldTaskNamesTheNewOne(t *testing.T) {
	got, err := redoRegister(redoRows(), "T-old1", "T-new1", 12, redoAt)
	if err != nil || len(got) != 4 {
		t.Fatalf("redoRegister = %v, %v; want four rows", got, err)
	}
	if !slices.Equal(got[2].Predecessors, []string{"T-new1", "T-spec"}) || !slices.Equal(got[3].Predecessors, []string{"T-spec"}) {
		t.Errorf("the predecessors of the successors = %v and %v; want [T-new1 T-spec] and [T-spec]", got[2].Predecessors, got[3].Predecessors)
	}
	var b bytes.Buffer
	if err := WriteRegister(&b, got); err != nil {
		t.Errorf("the register after the redo does not pass its own reader: %v", err)
	}
}

const redoBinding = "option\tconstraint\tresult\tevidence\ndb\tREQ-001\tpass\tdocs/spec/db.md\n"

func redoPlanOf() redoPlan {
	events := [][]string{
		{"1", "attempt", "1", "", redoOld, "", "", "2026-10-20T08:00:00Z"},
		{"2", "attempt", "2", "", redoOld, "", "", "2026-10-20T08:10:00Z"},
	}
	return redoPlan{Tasks: redoRows(), Old: "T-old1", New: "T-new1", Issue: 12, Attempt: 2, Base: redoNew, At: redoAt, OldEvents: events,
		Bindings: map[string][]byte{"db.tsv": []byte(redoBinding)}}
}

const (
	redoOld = "1111111111111111111111111111111111111111"
	redoNew = "2222222222222222222222222222222222222222"
)

func TestTheRedoIsOneCommitOfTheRegisterTheTwoEventsAndTheBindings(t *testing.T) {
	files, err := redoFiles(redoPlanOf())
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for p := range files {
		paths = append(paths, p)
	}
	slices.Sort(paths)
	if want := []string{"tasks.tsv", "tasks/T-new1/events.tsv", "tasks/T-new1/options/db.tsv", "tasks/T-old1/events.tsv"}; !slices.Equal(paths, want) {
		t.Fatalf("the files = %v; want %v", paths, want)
	}
	tasks, err := ReadRegister(files["tasks.tsv"])
	if err != nil || tasks[1].ID != "T-new1" || tasks[1].Issue != 12 || tasks[1].PR != 0 {
		t.Errorf("tasks.tsv: %v, %v; want T-new1 in place of the old row", tasks, err)
	}
	old, err := records.ReadEvents(files["tasks/T-old1/events.tsv"])
	if err != nil || len(old) != 3 || old[2][1] != "closed" || old[2][2] != "2" || old[2][3] != "" || old[2][6] != "T-new1" || old[2][7] != "2026-10-20T09:00:00Z" {
		t.Errorf("the events of the old task = %v, %v; want the two attempts and closed 2 naming T-new1", old, err)
	}
	renew, err := records.ReadEvents(files["tasks/T-new1/events.tsv"])
	if err != nil || len(renew) != 1 || renew[0][1] != "rebased" || renew[0][2] != "1" || renew[0][4] != redoNew || renew[0][6] != "T-old1" {
		t.Errorf("the events of the new task = %v, %v; want one rebased: attempt 1, base %s, detail T-old1", renew, err, redoNew)
	}
	if string(files["tasks/T-new1/options/db.tsv"]) != redoBinding {
		t.Errorf("the binding = %q; want the bytes of the old one", files["tasks/T-new1/options/db.tsv"])
	}
	if _, err := spec.ReadOptions(files["tasks/T-new1/options/db.tsv"]); err != nil {
		t.Errorf("the binding does not read: %v", err)
	}
	for name, mutate := range map[string]func(*redoPlan){
		"a binding that does not read":      func(p *redoPlan) { p.Bindings["db.tsv"] = []byte("not a table\n") },
		"a base that is no SHA":             func(p *redoPlan) { p.Base = "main" },
		"an old task that is no row":        func(p *redoPlan) { p.Old = "T-zzzz" },
		"an old task with no attempt event": func(p *redoPlan) { p.OldEvents = nil },
	} {
		p := redoPlanOf()
		mutate(&p)
		if files, err := redoFiles(p); err == nil {
			t.Errorf("%s: %d files, nil; want an error and no file", name, len(files))
		}
	}
}

func TestTheIntentOfARedoIsReadBackOrRefused(t *testing.T) {
	id, base, err := parseIntent(intentBytes("T-new1", redoNew))
	if err != nil || id != "T-new1" || base != redoNew {
		t.Errorf("parseIntent(intentBytes) = %q, %q, %v; want T-new1 and %s", id, base, err, redoNew)
	}
	if got := string(intentBytes("T-new1", redoNew)); got != "T-new1 "+redoNew+"\n" {
		t.Errorf("intentBytes = %q", got)
	}
	for _, bad := range []string{"", "T-new1\n", "T-new1 " + redoNew, "T-NEW1 " + redoNew + "\n", "T-new1 main\n", "T-new1 " + redoNew + " x\n", "T-new1 " + redoNew + "\n\n"} {
		if id, base, err := parseIntent([]byte(bad)); err == nil {
			t.Errorf("parseIntent(%q) = %q, %q, nil; want an error", bad, id, base)
		}
	}
}

func TestTheFirstAttemptOfARedoneTaskStartsAtTheRebasedBase(t *testing.T) {
	rebased := [][]string{{"1", "rebased", "1", "", redoNew, "", "T-old1", "2026-10-20T09:00:00Z"}}
	if got := firstBase(rebased, redoOld); got != redoNew {
		t.Errorf("firstBase(rebased) = %s; want the base of the event, %s", got, redoNew)
	}
	if got := firstBase(nil, redoOld); got != redoOld {
		t.Errorf("firstBase(none) = %s; want the base the caller gives (a green negative)", got)
	}
}

// F7 of round 2 (late test): the pin is inherited by value, so no length of the chain of redos changes it and none ends in a silent root.
func TestThePinIsInheritedThroughAnyChainOfRedos(t *testing.T) {
	pin := pinBytes(strings.Repeat("3e", 32), strings.Repeat("f6", 20))
	tasks := redoRows()
	old := "T-old1"
	for i := 0; i < 70; i++ {
		id := "T-" + strconv.FormatInt(int64(36*36*36+i), 36)
		p := redoPlan{Tasks: tasks, Old: old, New: id, Issue: 100 + i, Attempt: 1, Base: redoNew, At: redoAt,
			OldEvents: [][]string{{"1", "attempt", "1", "", redoOld, "", "", "2026-10-20T08:00:00Z"}}, Pin: pin}
		files, err := redoFiles(p)
		if err != nil {
			t.Fatalf("redo %d: %v", i+1, err)
		}
		if got := string(files[pinPath(id)]); got != string(pin) {
			t.Fatalf("redo %d: %s = %q; want the pin as it was", i+1, pinPath(id), got)
		}
		if tasks, err = ReadRegister(files["tasks.tsv"]); err != nil {
			t.Fatal(err)
		}
		old = id
	}
	if files, err := redoFiles(redoPlan{Tasks: redoRows(), Old: "T-old1", New: "T-new1", Issue: 12, Attempt: 2, Base: redoNew, At: redoAt, OldEvents: redoPlanOf().OldEvents}); err != nil || files[pinPath("T-new1")] != nil {
		t.Errorf("a plan with no pin: %v, %q; want no pin file", err, files[pinPath("T-new1")])
	}
}

func TestThePinIsReadBackOrRefused(t *testing.T) {
	sum, rec := strings.Repeat("3e", 32), strings.Repeat("f6", 20)
	gotSum, gotRec, err := parsePin(pinBytes(sum, rec))
	if err != nil || gotSum != sum || gotRec != rec || string(pinBytes(sum, rec)) != sum+" "+rec+"\n" {
		t.Errorf("parsePin(pinBytes) = %q, %q, %v", gotSum, gotRec, err)
	}
	for _, bad := range []string{"", "a\n", "a b", "a b c\n", "a  b\n", "a b\n\n", "not-a-digest not-a-commit\n", sum[:63] + " " + rec + "\n", sum + " " + rec[:39] + "\n",
		strings.Repeat("g", 64) + " " + rec + "\n", sum + " " + strings.ToUpper(rec) + "\n", sum + "\t" + rec + "\n", sum + " " + rec[:20] + "\t" + rec[20:] + "\n", sum + " " + rec + "\n" + sum + " " + rec + "\n"} {
		if gs, gr, err := parsePin([]byte(bad)); err == nil {
			t.Errorf("parsePin(%q) = %q, %q, nil; want an error", bad, gs, gr)
		}
	}
	p := redoPlanOf()
	p.Pin = []byte("not-a-digest not-a-commit\n")
	if files, err := redoFiles(p); err == nil {
		t.Errorf("redoFiles with a malformed pin: %d files, nil; want an error", len(files))
	}
}
