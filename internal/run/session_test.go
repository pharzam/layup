package run

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
)

// standinStore is a records branch in memory: its files at the run's base,
// each commit of the run, and the pushes that git refuses before one passes.
type standinStore struct {
	files   map[string][]byte
	commits []map[string][]byte
	refuse  int    // the next pushes that are refused
	holder  string // the run that the lease names
}

func (s *standinStore) ReadLease(context.Context) (LeaseRow, error) {
	return LeaseRow{Run: s.holder, State: "held"}, nil
}
func (s *standinStore) WriteLease(context.Context, LeaseRow, string) error { return nil }
func (s *standinStore) ReadFile(_ context.Context, path string) ([]byte, error) {
	return s.files[path], nil
}
func (s *standinStore) Commit(_ context.Context, files map[string][]byte, _ string) error {
	if s.refuse > 0 {
		s.refuse--
		return ErrRefused
	}
	s.commits = append(s.commits, files)
	for p, d := range files {
		s.files[p] = d
	}
	return nil
}

const aSHA = "0123456789abcdef0123456789abcdef01234567"

var sessionAt = time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)

func newSessions(store *standinStore) *Sessions {
	return &Sessions{Store: store, RunID: "0123456789abcdef", Host: "/h", Target: "acme/target", Clone: "/c", Branch: "main",
		Now: func() time.Time { return sessionAt }}
}

func events(t *testing.T, s *standinStore, task string) [][]string {
	t.Helper()
	data := s.files["tasks/"+task+"/events.tsv"]
	if data == nil {
		return nil
	}
	rows, err := records.ReadEvents(data)
	if err != nil {
		t.Fatalf("events.tsv: %v", err)
	}
	return rows
}

func TestStartAttempt(t *testing.T) {
	s := &standinStore{files: map[string][]byte{}, holder: "0123456789abcdef"}
	r := newSessions(s)
	for want := 1; want <= 2; want++ {
		n, err := r.StartAttempt(context.Background(), "T-ab12", aSHA)
		if err != nil || n != want {
			t.Fatalf("attempt %d: %d, %v", want, n, err)
		}
	}
	got := events(t, s, "T-ab12")
	want := [][]string{{"1", "attempt", "1", "", aSHA, "", "", "2026-10-09T12:00:00Z"}, {"2", "attempt", "2", "", aSHA, "", "", "2026-10-09T12:00:00Z"}}
	if !slices.EqualFunc(got, want, slices.Equal) || len(s.commits) != 2 {
		t.Errorf("events %q in %d commits; want %q in 2", got, len(s.commits), want)
	}
	// An event of another kind does not raise the next attempt, whatever its
	// attempt: only an event attempt does.
	s.files["tasks/T-ab12/events.tsv"] = append(s.files["tasks/T-ab12/events.tsv"], "3\trefused\t7\tS-1a2b3c4d\t—\t—\tversion\t2026-10-09T12:00:00Z\n"...)
	if n, err := r.StartAttempt(context.Background(), "T-ab12", aSHA); err != nil || n != 3 {
		t.Errorf("after a refusal of attempt 7: %d, %v; want 3", n, err)
	}
	s.files["tasks/T-cd34/events.tsv"] = []byte("broken\n")
	if _, err := r.StartAttempt(context.Background(), "T-cd34", aSHA); err == nil || len(s.commits) != 3 {
		t.Errorf("an events.tsv that its reader refuses: %v, %d commits; want an error and no commit", err, len(s.commits))
	}
}

// task gives a TaskSpec of attempt 1 of T-ab12 and Sessions whose session
// functions record their calls and touch no file.
func task(store *standinStore, calls *[]string) (*Sessions, TaskSpec) {
	r := newSessions(store)
	r.sweep = func(host, target string) error { *calls = append(*calls, "sweep"); return nil }
	r.make = func(id string, spec TaskSpec) (session.Dir, error) {
		*calls = append(*calls, "make")
		return session.Dir{Root: "/h/sessions/" + id, Repo: "/h/sessions/" + id + "/repo", Home: "/h/sessions/" + id + "/home",
			Tmp: "/h/sessions/" + id + "/tmp", Prompt: "/h/sessions/" + id + "/prompt.md"}, nil
	}
	r.remove = func(path string) error { *calls = append(*calls, "remove "+path); return nil }
	// The environment of a stand-in: the credential of var: as CRED, which
	// the version check must not get.
	r.environ = func(d session.Dir, h session.Harness) ([]string, error) {
		env := []string{"PATH=/bin", "HOME=" + d.Home, "TMPDIR=" + d.Tmp}
		if strings.HasPrefix(h.CredentialTo, "var:") {
			env = append(env, "CRED=secret")
		}
		return append(env, h.Vars...), nil
	}
	r.version = func(_ context.Context, command string, env []string, dir string) (string, error) {
		*calls = append(*calls, "version in "+dir+" with "+lookup(env, "HOME")+lookup(env, "CRED"))
		return "2.1.295 (Claude Code)", nil
	}
	r.rules = func(d session.Dir, rules, policy []string) ([]string, error) {
		*calls = append(*calls, "rules")
		return []string{"/etc/claude/policy.json"}, nil
	}
	r.process = func(s session.Spec) (session.Run, error) {
		*calls = append(*calls, "process with "+lookup(s.Env, "CRED"))
		return session.Run{}, nil
	}
	r.newID = func() (string, error) { return "S-1a2b3c4d", nil }
	r.exists = func(string) bool { return false } // a unit test touches no file
	return r, TaskSpec{Task: "T-ab12", Attempt: 1, Role: "developer", Base: aSHA, Records: aSHA, Prompt: []byte("the task\n"),
		Pair: Pair{Harness: "claude", Model: "claude-fable-5-1", VersionCommand: "claude --version", Command: "claude -p {prompt}",
			PromptMode: "arg", Cap: "5.0", Wall: 30, Rules: []string{"CLAUDE.md"}, Policy: []string{"/etc/claude/policy.json"},
			Context: 1000, Session: session.Harness{CredentialTo: "var:TOKEN", Credential: "/k", Vars: []string{"A=1"}}},
		Admit: func(v string) error { *calls = append(*calls, "admit "+v); return nil },
		End: func(id string, d session.Dir, r session.Run, err error) error {
			*calls = append(*calls, "end")
			return nil
		},
	}
}

func lookup(env []string, name string) string {
	for _, e := range env {
		if v, ok := strings.CutPrefix(e, name+"="); ok {
			return v
		}
	}
	return ""
}

// withAttempt gives a store whose T-ab12 has the event attempt 1.
func withAttempt() *standinStore {
	return &standinStore{holder: "0123456789abcdef", files: map[string][]byte{
		"tasks/T-ab12/events.tsv": []byte("n\tkind\tattempt\tsession\tbase\tsha\tdetail\ttime\n1\tattempt\t1\t—\t" + aSHA + "\t—\t—\t2026-10-09T11:00:00Z\n")}}
}

func TestTheStartRowBeforeTheProcess(t *testing.T) {
	store := withAttempt()
	var calls []string
	r, spec := task(store, &calls)
	store.refuse = 1 // one refused push: Fenced reads the lease and tries once more
	id, err := r.TaskSession(context.Background(), spec)
	want := []string{"sweep", "make", "version in /h/sessions/S-1a2b3c4d/repo with /h/sessions/S-1a2b3c4d/tmp", "admit 2.1.295 (Claude Code)", "rules", "process with secret", "end"}
	if err != nil || id != "S-1a2b3c4d" || !slices.Equal(calls, want) {
		t.Fatalf("%q, %v, the calls %q; want %q", id, err, calls, want)
	}
	// One commit holds the start row and the event session.
	if len(store.commits) != 1 || store.commits[0]["sessions.tsv"] == nil || store.commits[0]["tasks/T-ab12/events.tsv"] == nil {
		t.Fatalf("the commits %d; want one with sessions.tsv and events.tsv", len(store.commits))
	}
	rows, err := records.ReadSessions(store.files["sessions.tsv"])
	wantRow := []string{"S-1a2b3c4d", "T-ab12", "1", "developer", "claude", "2.1.295 (Claude Code)", "claude-fable-5-1", aSHA, aSHA, "9", "3", "1000", "5.0", "30", "A=1", "/etc/claude/policy.json"}
	if err != nil || len(rows) != 1 || !slices.Equal(rows[0], wantRow) {
		t.Errorf("sessions.tsv %q, %v; want %q", rows, err, wantRow)
	}
	ev := events(t, store, "T-ab12")
	if got := ev[len(ev)-1]; !slices.Equal(got, []string{"2", "session", "1", "S-1a2b3c4d", "", "", "", "2026-10-09T12:00:00Z"}) {
		t.Errorf("the event session %q", got)
	}
}

func TestARefusedStart(t *testing.T) {
	for _, c := range []struct {
		reason, detail string
		set            func(r *Sessions, spec *TaskSpec, store *standinStore)
	}{
		{"attempt", "attempt", func(r *Sessions, spec *TaskSpec, _ *standinStore) { spec.Attempt = 2 }},
		{"no events", "attempt", func(_ *Sessions, _ *TaskSpec, s *standinStore) { delete(s.files, "tasks/T-ab12/events.tsv") }},
		{"version", "version", func(r *Sessions, _ *TaskSpec, _ *standinStore) {
			r.version = func(context.Context, string, []string, string) (string, error) {
				return "", session.Refusal{Reason: "version"}
			}
		}},
		{"probe", "probe", func(_ *Sessions, spec *TaskSpec, _ *standinStore) {
			spec.Admit = func(string) error { return session.Refusal{Reason: "probe"} }
		}},
		{"context", "context 3 2", func(_ *Sessions, spec *TaskSpec, _ *standinStore) { spec.Pair.Context = 2 }},
		{"prompt", "prompt 131072", func(_ *Sessions, spec *TaskSpec, _ *standinStore) { spec.Prompt = []byte(strings.Repeat("x", 131072)) }},
		{"rules", "rules /h/CLAUDE.md", func(r *Sessions, _ *TaskSpec, _ *standinStore) {
			r.rules = func(session.Dir, []string, []string) ([]string, error) {
				return nil, session.Refusal{Reason: "rules", Value: "/h/CLAUDE.md"}
			}
		}},
	} {
		store := withAttempt()
		var calls []string
		r, spec := task(store, &calls)
		if c.reason == "prompt" {
			spec.Pair.Context = 1 << 20
		}
		c.set(r, &spec, store)
		id, err := r.TaskSession(context.Background(), spec)
		var refusal session.Refusal
		if !errors.As(err, &refusal) || id != "S-1a2b3c4d" {
			t.Errorf("%s: %q, %v; want a refusal", c.reason, id, err)
		}
		if slices.ContainsFunc(calls, isProcess) || !slices.Contains(calls, "remove /h/sessions/S-1a2b3c4d") {
			t.Errorf("%s: the calls %q; want no process and the directory removed", c.reason, calls)
		}
		if store.files["sessions.tsv"] != nil || len(store.commits) != 1 {
			t.Errorf("%s: %d commits, sessions.tsv %q; want one commit and no start row", c.reason, len(store.commits), store.files["sessions.tsv"])
		}
		ev := events(t, store, "T-ab12")
		want := []string{strconv.Itoa(len(ev)), "refused", strconv.Itoa(spec.Attempt), "S-1a2b3c4d", "", "", c.detail, "2026-10-09T12:00:00Z"}
		if got := ev[len(ev)-1]; !slices.Equal(got, want) {
			t.Errorf("%s: the event %q; want %q", c.reason, got, want)
		}
	}
}

func TestALostLeaseStopsTheStart(t *testing.T) {
	store := withAttempt()
	var calls []string
	r, spec := task(store, &calls)
	store.refuse = 2 // refused twice: Fenced's LostError
	_, err := r.TaskSession(context.Background(), spec)
	var lost *LostError
	if !errors.As(err, &lost) || slices.ContainsFunc(calls, isProcess) {
		t.Errorf("%v, the calls %q; want a LostError and no process", err, calls)
	}
	if !slices.Contains(calls, "remove /h/sessions/S-1a2b3c4d") {
		t.Errorf("the calls %q; want the directory removed", calls)
	}
	// A lease of another run: the first refusal is a LostError at once.
	store, calls = withAttempt(), nil
	r, spec = task(store, &calls)
	store.refuse, store.holder = 1, "fedcba9876543210"
	if _, err := r.TaskSession(context.Background(), spec); !errors.As(err, &lost) || lost.Holder != "fedcba9876543210" {
		t.Errorf("a lease of another run: %v", err)
	}
}

func isProcess(call string) bool { return strings.HasPrefix(call, "process") }

// A start whose directory-making fails part-way removes what it made: a
// copied credential does not wait for the next sweep. A directory that was
// there before is another session's, and is kept.
func TestAFailedMakeRemovesItsDirectory(t *testing.T) {
	for _, existed := range []bool{false, true} {
		store := withAttempt()
		var calls []string
		r, spec := task(store, &calls)
		r.make = func(id string, spec TaskSpec) (session.Dir, error) {
			calls = append(calls, "make")
			return session.Dir{}, errors.New("the clone failed")
		}
		// The root exists before Make when another session holds it, and
		// after a Make that made it and then failed.
		r.exists = func(path string) bool { return existed || slices.Contains(calls, "make") }
		_, err := r.TaskSession(context.Background(), spec)
		removed := slices.Contains(calls, "remove /h/sessions/S-1a2b3c4d")
		if err == nil || removed == existed || len(store.commits) != 0 {
			t.Errorf("a directory there before %v: %v, the calls %q, %d commits; want an error, removed %v, no commit", existed, err, calls, len(store.commits), !existed)
		}
	}
}
