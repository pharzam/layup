package run

// This file holds the start of a task session (docs/spec/session.md, The start
// of a session): the call that starts an attempt, and the hooks of the one call
// of internal/session for a task session, with its start row and its refused
// start (row 36a of the plan, task T-d8t9, #164).

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/tsv"
)

// SessionStore is the records branch as a session sees it: the lease, for
// Fenced; ReadFile gives a file at the run's records base, nil for none; Commit
// makes one records commit of files and pushes it, and gives ErrRefused when
// git refuses the push.
type SessionStore interface {
	Store
	ReadFile(ctx context.Context, path string) ([]byte, error)
	Commit(ctx context.Context, files map[string][]byte, message string) error
}

// Sessions starts the attempts and the task sessions of one run of one
// target: its records store, its run ID, the host directory, the target
// OWNER/NAME, the run's clone and its default branch, and its clock
// (time.Now when nil); for the end of a task session (row 36b), the forge,
// the control issue (issue.control of start.tsv) and the rows of
// host:prices.tsv as records.ReadPrices gives them (none for a missing file:
// an empty table, decided here, task T-fsjp). The function fields are the
// session's own, which a unit test replaces.
type Sessions struct {
	Store         SessionStore
	RunID         string
	Host, Target  string
	Clone, Branch string
	Now           func() time.Time
	Forge         forge.Forge
	Control       int
	Prices        [][]string
	newID         func() (string, error)
	sweep         func(host, target string) error
	make          func(id string, spec TaskSpec) (session.Dir, error)
	remove        func(path string) error
	version       func(ctx context.Context, command string, env []string, dir string) (string, error)
	rules         func(d session.Dir, rules, policy []string) ([]string, error)
	process       func(s session.Spec) (session.Run, error)
	environ       func(d session.Dir, h session.Harness) ([]string, error)
	exists        func(path string) bool
}

// Pair is the admitted pair of a session, as internal/route reads the harness
// register row and the models row: their IDs, the row's commands, how the
// prompt goes, the spend cap ("" for none), the wall-clock limit in minutes,
// the rule-file names and the policy paths, the model's context size, and the
// credential and fixed variables of the session's files.
type Pair struct {
	Harness, Model                      string
	VersionCommand, Command, PromptMode string
	Cap                                 string
	Wall                                int
	Rules, Policy                       []string
	Context                             int
	Session                             session.Harness
	Billing, Usage                      string // the row's billing and usage format (row 36b)
}

// TaskSpec is one task session to start: the task, its attempt and role, the
// base commit, the records commit that the prompt was built from, the prompt
// and the pair; Admit gives nil when the version read may start (row 38 runs
// the probe in it), else Refusal{probe}; End and Push are the end of the
// session and the checks before a push, the push and the bind (rows 36b, 37a
// and 37b), End nil for the end of a task session (endTask, row 36b), Push
// nil for the checks before a push (pushTask, row 37a).
type TaskSpec struct {
	Task          string
	Attempt       int
	Role          string
	Base, Records string
	Prompt        []byte
	Pair          Pair
	Admit         func(version string) error
	End           func(id string, d session.Dir, r session.Run, err error) error
	Push          func(id string, d session.Dir, r session.Run) error
}

func (s *Sessions) fns() {
	if s.newID == nil {
		s.newID, s.sweep, s.remove = session.NewID, session.Sweep, os.RemoveAll
		s.version, s.rules, s.process, s.environ = session.Version, session.CheckRuleFiles, session.Process, session.Environ
		s.make = func(id string, spec TaskSpec) (session.Dir, error) {
			return session.Make(s.Host, s.Target, id, s.Clone, s.Branch, spec.Task, spec.Attempt, spec.Base, spec.Prompt, spec.Pair.Session)
		}
	}
	if s.environ == nil {
		s.environ = session.Environ
	}
	if s.exists == nil {
		s.exists = func(path string) bool { _, err := os.Lstat(path); return err == nil }
	}
	if s.Now == nil {
		s.Now = time.Now
	}
}

func eventsPath(task string) string { return "tasks/" + task + "/events.tsv" }

// readTable reads the table at path by its reader; a missing file is an
// empty table (decided here, task T-d8t9).
func (s *Sessions) readTable(ctx context.Context, path string, read func([]byte) ([][]string, error)) ([][]string, error) {
	data, err := s.Store.ReadFile(ctx, path)
	if err != nil || data == nil {
		return nil, err
	}
	rows, err := read(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return rows, nil
}

// appendEvent gives the files of events.tsv of task with the event added, its
// number the next, after CheckEventsAppend.
func (s *Sessions) appendEvent(ctx context.Context, task string, kind string, attempt int, id, base, detail string) ([]byte, error) {
	before, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	row := []string{strconv.Itoa(len(before) + 1), kind, strconv.Itoa(attempt), id, base, "", detail, s.Now().UTC().Format(timeForm)}
	after := append(append([][]string{}, before...), row)
	if err := records.CheckEventsAppend(before, after); err != nil {
		return nil, err
	}
	if err := records.CheckEvents(after); err != nil {
		return nil, err
	}
	return table(records.EventsSchema, after)
}

// StartAttempt starts the next attempt of task at base: the event attempt,
// its number one more than the largest attempt of an event attempt, 1 for
// none (decided here, task T-d8t9), in one records commit with its push,
// through Fenced. It gives the attempt.
func (s *Sessions) StartAttempt(ctx context.Context, task, base string) (int, error) {
	s.fns()
	rows, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return 0, err
	}
	n := 1
	for _, r := range rows {
		if a, _ := strconv.Atoi(r[2]); r[1] == "attempt" && a >= n {
			n = a + 1
		}
	}
	data, err := s.appendEvent(ctx, task, "attempt", n, "", base, "")
	if err != nil {
		return 0, err
	}
	return n, s.commit(ctx, map[string][]byte{eventsPath(task): data}, fmt.Sprintf("layup run: attempt %d of %s", n, task))
}

func (s *Sessions) commit(ctx context.Context, files map[string][]byte, message string) error {
	return Fenced(ctx, s.Store, s.RunID, func(ctx context.Context) error { return s.Store.Commit(ctx, files, message) })
}

// TaskSession starts a task session through the one call of internal/session
// and gives its ID. Before the call it sweeps the directories that a stopped
// run of the target left; the directory is made right after the ID, so each
// check runs on the session's own files (decided here, task T-d8t9). The
// start row, a row of sessions.tsv and the event session, is one records
// commit, pushed before the process starts. A refused start removes the
// directory, then commits the event refused with the session ID and the
// reason and its value, and no start row; an error before the process starts,
// a lost lease among them, removes the directory too.
func (s *Sessions) TaskSession(ctx context.Context, spec TaskSpec) (string, error) {
	s.fns()
	if err := s.sweep(s.Host, s.Target); err != nil {
		return "", err
	}
	var (
		d        session.Dir
		made     bool
		started  bool
		version  string
		estimate int
		policy   []string
	)
	steps := session.Steps{
		ID: func() (string, error) {
			id, err := s.newID()
			if err != nil {
				return "", err
			}
			// A directory of this ID that was there before is another
			// session's; one that Make made before it failed is this start's.
			root := filepath.Join(s.Host, "sessions", id)
			if abs, err := filepath.Abs(root); err == nil {
				root = abs
			}
			existed := s.exists(root)
			if d, err = s.make(id, spec); err != nil {
				if !existed && s.exists(root) {
					d, made = session.Dir{Root: root}, true
				}
				return id, err
			}
			made = true
			return id, nil
		},
		Attempt: func(string) error {
			rows, err := s.readTable(ctx, eventsPath(spec.Task), records.ReadEvents)
			if err != nil {
				return err
			}
			for _, r := range rows {
				if r[1] == "attempt" && r[2] == strconv.Itoa(spec.Attempt) {
					return nil
				}
			}
			return session.Refusal{Reason: "attempt"}
		},
		Version: func(string) error {
			// No credential (step 2): no credential variable, and HOME the
			// session's tmp/, as home/ holds a credential of file: once
			// Make has run (decided here, task T-d8t9).
			env, err := s.environ(d, session.Harness{Vars: spec.Pair.Session.Vars})
			if err != nil {
				return err
			}
			for i, e := range env {
				if strings.HasPrefix(e, "HOME=") {
					env[i] = "HOME=" + d.Tmp
				}
			}
			if version, err = s.version(ctx, spec.Pair.VersionCommand, env, d.Repo); err != nil {
				return err
			}
			return spec.Admit(version)
		},
		Context: func(string) (err error) {
			estimate, err = session.CheckContext(spec.Prompt, spec.Pair.Context)
			return err
		},
		Prompt: func(string) error { return session.CheckPromptSize(spec.Prompt, spec.Pair.PromptMode) },
		Rules: func(string) (err error) {
			policy, err = s.rules(d, spec.Pair.Rules, spec.Pair.Policy)
			return err
		},
		StartRow: func(id string) error { return s.startRow(ctx, id, spec, version, estimate, policy) },
		Process: func(string) (session.Run, error) {
			started = true
			env, err := s.environ(d, spec.Pair.Session)
			if err != nil {
				return session.Run{}, err
			}
			words := session.Words(spec.Pair.Command, spec.Pair.Model, spec.Pair.Cap, session.Prompt(spec.Pair.PromptMode, d, spec.Prompt))
			return s.process(session.NewSpec(d, words, env, spec.Pair.PromptMode, time.Duration(spec.Pair.Wall)*time.Minute))
		},
		End: func(id string, r session.Run, err error) error {
			if spec.End == nil {
				return s.endTask(ctx, spec, id, d, r, err)
			}
			return spec.End(id, d, r, err)
		},
	}
	steps.Push = func(id string, r session.Run) error { return s.pushTask(ctx, id) }
	if spec.Push != nil {
		steps.Push = func(id string, r session.Run) error { return spec.Push(id, d, r) }
	}
	id, err := session.Call(steps)
	if started {
		// Once the process started, the directory is removed after the call,
		// whether the end and the push hooks succeed or fail, so a copied
		// credential never waits for the next sweep (decided here, task
		// T-fsjp).
		if rmErr := s.remove(d.Root); rmErr != nil {
			err = errors.Join(err, rmErr)
		}
		// The comment comes after every record of the session, the push
		// hook's too (decided here, task T-fsjp, round 1 of #165).
		// A lost lease stops the run before any other write (run.md).
		if !errors.As(err, new(*LostError)) {
			if cErr := s.comment(ctx, id); cErr != nil {
				err = errors.Join(err, cErr)
			}
		}
		return id, err
	}
	if err == nil {
		return id, nil
	}
	if made {
		if rmErr := s.remove(d.Root); rmErr != nil {
			return id, errors.Join(err, rmErr)
		}
	}
	var refusal session.Refusal
	if !errors.As(err, &refusal) {
		return id, err
	}
	detail := strings.TrimSpace(refusal.Reason + " " + refusal.Value)
	data, appendErr := s.appendEvent(ctx, spec.Task, "refused", spec.Attempt, id, "", detail)
	if appendErr != nil {
		return id, errors.Join(err, appendErr)
	}
	if cErr := s.commit(ctx, map[string][]byte{eventsPath(spec.Task): data}, fmt.Sprintf("layup run: %s refused (%s)", id, refusal.Reason)); cErr != nil {
		return id, errors.Join(err, cErr)
	}
	return id, err
}

// startRow commits the start row of a session: its row of sessions.tsv and
// the event session of its task, in one records commit (decided here, task
// T-d8t9: a crash leaves both or neither), through Fenced.
func (s *Sessions) startRow(ctx context.Context, id string, spec TaskSpec, version string, estimate int, policy []string) error {
	before, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return err
	}
	vars, err := tsv.JoinList(spec.Pair.Session.Vars)
	if err != nil {
		return err
	}
	pol, err := tsv.JoinList(policy)
	if err != nil {
		return err
	}
	row := []string{id, spec.Task, strconv.Itoa(spec.Attempt), spec.Role, spec.Pair.Harness, version, spec.Pair.Model, spec.Base, spec.Records,
		strconv.Itoa(len(spec.Prompt)), strconv.Itoa(estimate), strconv.Itoa(spec.Pair.Context), spec.Pair.Cap, strconv.Itoa(spec.Pair.Wall), vars, pol}
	if err := records.CheckSession(row); err != nil {
		return err
	}
	sessions, err := table(records.SessionsSchema, append(append([][]string{}, before...), row))
	if err != nil {
		return err
	}
	events, err := s.appendEvent(ctx, spec.Task, "session", spec.Attempt, id, "", "")
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{"sessions.tsv": sessions, eventsPath(spec.Task): events},
		fmt.Sprintf("layup run: the start of %s, attempt %d of %s", id, spec.Attempt, spec.Task))
}
