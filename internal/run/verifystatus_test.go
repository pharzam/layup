package run

import (
	"context"
	"strconv"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

const verifier = "S-e0000001"

// A verifyWorld is the records of T-ab12 on the stand-in store with three admitted
// harnesses (claude, devin, codex), the routing of the step verification, and a
// frozen list whose tests were written by claude, an author.
type verifyWorld struct {
	s                 *Sessions
	store             *standinStore
	harnesses, models [][]string
	in                VerifyInput
	events            [][]string
}

func newVerifyWorld(t *testing.T, reviewer string) *verifyWorld {
	t.Helper()
	s, harnesses, models := pairWorld(t, "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\nqa-engineer\treasoning\t3\tcodex\tm1\n")
	store := s.Store.(*standinStore)
	frozen, err := WriteTests([][]string{{"T-1", "x_test.go", hashA, "claude", "2026-10-16T12:00:00Z"}})
	if err != nil {
		t.Fatal(err)
	}
	store.files[testsPath("T-ab12")] = frozen
	w := &verifyWorld{s: s, store: store, harnesses: harnesses, models: models,
		in: VerifyInput{Task: "T-ab12", Base: shaOld, Head: headA, Session: verifier, Harnesses: harnesses, Models: models}}
	w.events = [][]string{{"1", "attempt", "1", "", shaOld, "", "", "2026-10-16T11:00:00Z"}}
	w.session(t, reviewer, headA)
	return w
}

// session records the start of the verifier on the harness, reviewing head.
func (w *verifyWorld) session(t *testing.T, harness, reviewed string) {
	t.Helper()
	row := sessionRow(verifier, "T-ab12", "verification", harness)
	row[sessionBase] = reviewed
	data, err := table(records.SessionsSchema, [][]string{row})
	if err != nil {
		t.Fatal(err)
	}
	w.store.files["sessions.tsv"] = data
	w.events = append(w.events, []string{strconv.Itoa(len(w.events) + 1), "session", "1", verifier, "", "", "", "2026-10-16T11:00:00Z"})
	w.write(t)
}

func (w *verifyWorld) write(t *testing.T) {
	t.Helper()
	data, err := table(records.EventsSchema, w.events)
	if err != nil {
		t.Fatal(err)
	}
	w.store.files[eventsPath("T-ab12")] = data
}

// end records the result of the verifier: its event, its result file and its kept copy.
func (w *verifyWorld) end(t *testing.T, result, text string) {
	t.Helper()
	w.events = append(w.events, []string{strconv.Itoa(len(w.events) + 1), "result", "1", verifier, "", "", result, "2026-10-16T11:30:00Z"})
	w.write(t)
	rows := resultTop + "status\t1\tcompleted\t—\tdone\n"
	if text != "" {
		rows += "comment\t1\tresult/review.md\t" + commentSum(text) + "\t—\n"
		w.store.files[KeptPath("T-ab12", verifier, "review.md")] = []byte(text)
	}
	w.store.files["tasks/T-ab12/results/"+verifier+".tsv"] = []byte(rows)
}

func (w *verifyWorld) verify(t *testing.T) (string, Status) {
	t.Helper()
	sha, st, err := w.s.VerifyAt(context.Background(), w.in)
	if err != nil {
		t.Fatalf("VerifyAt: %v", err)
	}
	return sha, st
}

func TestNoHarnessOutsideTheAuthorsIsFailureNotActive(t *testing.T) {
	w := newVerifyWorld(t, "devin")
	w.store.files["routing.tsv"] = []byte("role\ttier\tposition\tharness\tmodel\nqa-engineer\treasoning\t1\tclaude\tm1\n")
	w.in.Session = ""
	sha, st := w.verify(t)
	if sha != headA || st != (Status{VerifyContext, "failure", "not-active: no harness outside the authors"}) {
		t.Errorf("VerifyAt = %s, %+v; want failure, not-active: no harness outside the authors, at the head", sha, st)
	}
	w.in.Waiting = true // the control: a verifier that waits is pending only when a harness is outside the authors
	if _, st = w.verify(t); st.State != "failure" {
		t.Errorf("waiting with no harness outside the authors: %+v; want not-active before any wait", st)
	}
	w.store.files["routing.tsv"] = []byte("role\ttier\tposition\tharness\tmodel\nqa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\n")
	if _, st = w.verify(t); st.State != "pending" {
		t.Errorf("the control, a harness outside the authors: %+v; want pending", st)
	}
}

func TestAWaitingVerifierIsPendingAndARunningOneIsToo(t *testing.T) {
	w := newVerifyWorld(t, "devin")
	w.in.Waiting = true
	if sha, st := w.verify(t); sha != headA || st != (Status{VerifyContext, "pending", "waits for its turn"}) {
		t.Errorf("a waiting verifier: %s, %+v", sha, st)
	}
	w.in.Waiting = false
	w.in.Head = headB // the head moved while the verifier reviews headA: the status stays on the reviewed head
	if sha, st := w.verify(t); sha != headA || st != (Status{VerifyContext, "pending", "verifier running"}) {
		t.Errorf("a running verifier: %s, %+v; want pending, verifier running, on the head it reviews", sha, st)
	}
}

func TestVerifyIsSuccessOnlyForNothingMaterialAtTheReviewedHeadWithAResultThatIsDone(t *testing.T) {
	clean := record(headA, "nothing material in scope")
	for _, c := range []struct {
		name, result, text string
		want               Status
		tamper             func(*verifyWorld)
	}{
		{"nothing material at the reviewed head", "done", clean, Status{VerifyContext, "success", "nothing material in scope"}, nil},
		{"a record of another head", "done", record(headB, "nothing material in scope"), Status{VerifyContext, "failure", "review record of another head"}, nil},
		{"a finding", "done", record(headA, "material", located), Status{VerifyContext, "failure", "material: 1 findings"}, nil},
		{"material with an empty block", "done", record(headA, "material"), Status{VerifyContext, "failure", "material: 0 findings"}, nil},
		{"nothing material with a finding in the block", "done", record(headA, "nothing material in scope", located), Status{VerifyContext, "failure", "material: 1 findings"}, nil},
		{"a copy that is no review record", "done", "all good\n", Status{VerifyContext, "failure", "no review record"}, nil},
		{"a result that is not done", "crash", clean, Status{VerifyContext, "failure", "no review record"}, nil},
		{"a refused result", "done", clean, Status{VerifyContext, "failure", "no review record"}, func(w *verifyWorld) {
			w.events = append(w.events, []string{strconv.Itoa(len(w.events) + 1), "refused", "1", verifier, "", "", "handoff", "2026-10-16T11:31:00Z"})
		}},
		{"a kept copy that the result does not name", "done", clean, Status{VerifyContext, "failure", "no review record"}, func(w *verifyWorld) {
			w.store.files[KeptPath("T-ab12", verifier, "review.md")] = []byte(clean + "- edited\n")
		}},
		{"no kept copy", "done", clean, Status{VerifyContext, "failure", "no review record"}, func(w *verifyWorld) {
			delete(w.store.files, KeptPath("T-ab12", verifier, "review.md"))
		}},
		{"a result with no comment row", "done", "", Status{VerifyContext, "failure", "no review record"}, nil},
	} {
		w := newVerifyWorld(t, "devin")
		w.end(t, c.result, c.text)
		if c.tamper != nil {
			c.tamper(w)
			w.write(t)
		}
		if sha, st := w.verify(t); sha != headA || st != c.want {
			t.Errorf("%s: %s, %+v; want %+v on the reviewed head", c.name, sha, st, c.want)
		}
	}
}

func TestAVerifierThatIsAnAuthorOrNotAdmittedIsNotIndependent(t *testing.T) {
	want := Status{VerifyContext, "failure", "not-active: the verifier is not independent"}
	w := newVerifyWorld(t, "claude") // an author of the change: it wrote the frozen tests; devin and codex are available
	w.end(t, "done", record(headA, "nothing material in scope"))
	if sha, st := w.verify(t); sha != headA || st != want {
		t.Errorf("a clean record of an author: %s, %+v; want %+v though another harness is admitted", sha, st, want)
	}
	for name, change := range map[string]func(*verifyWorld){
		"a harness that failed its probe": func(w *verifyWorld) {
			rows, _ := records.ReadHarnesses(w.store.files["harnesses.tsv"])
			rows[1][4], rows[1][5] = "failed", "token"
			w.store.files["harnesses.tsv"], _ = table(records.HarnessesSchema, rows)
		},
		"a model of use no":            func(w *verifyWorld) { w.models[1][5] = "no"; w.in.Models = w.models },
		"a harness of no register row": func(w *verifyWorld) { w.in.Harnesses = [][]string{w.harnesses[0], w.harnesses[2]} },
	} {
		w := newVerifyWorld(t, "devin")
		w.end(t, "done", record(headA, "nothing material in scope"))
		if _, st := w.verify(t); st.State != "success" {
			t.Fatalf("%s, the control (devin is independent): %+v", name, st)
		}
		change(w)
		if _, st := w.verify(t); st != want {
			t.Errorf("%s: %+v; want %+v", name, st, want)
		}
	}
}

func TestAVerifierOfAnotherTaskOrRoleOrNoneIsAnError(t *testing.T) {
	for name, change := range map[string]func(*verifyWorld, *testing.T){
		"a session of another task": func(w *verifyWorld, t *testing.T) {
			row := sessionRow(verifier, "T-zz99", "verification", "devin")
			w.store.files["sessions.tsv"], _ = table(records.SessionsSchema, [][]string{row})
		},
		"a session of the developer role": func(w *verifyWorld, t *testing.T) {
			row := sessionRow(verifier, "T-ab12", "implement", "devin")
			w.store.files["sessions.tsv"], _ = table(records.SessionsSchema, [][]string{row})
		},
		"a session with no start row": func(w *verifyWorld, t *testing.T) { w.in.Session = "S-e0000009" },
		"no session and no wait":      func(w *verifyWorld, t *testing.T) { w.in.Session = "" },
	} {
		w := newVerifyWorld(t, "devin")
		change(w, t)
		if sha, st, err := w.s.VerifyAt(context.Background(), w.in); err == nil {
			t.Errorf("%s: %s, %+v, no error", name, sha, st)
		}
	}
	w := newVerifyWorld(t, "devin")
	if _, _, err := w.s.VerifyAt(context.Background(), w.in); err != nil {
		t.Errorf("the control, a verifier of the task and the role: %v", err)
	}
}
