package run

import (
	"bytes"
	"slices"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
)

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
