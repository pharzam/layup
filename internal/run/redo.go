package run

// This file holds the redo of a task whose base moved (docs/spec/loop.md, A task whose base moved: the carry and the redo; task T-89xf).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
)

// intentPath is the file of the decision to redo a task (records.md: tasks/<task>/redo.txt).
func intentPath(task string) string { return "tasks/" + task + "/redo.txt" }

// redoRegister gives the register with the row of old replaced in place by the row of the task id on issue: the kind, the requirements,
// the size and the predecessors kept, the pr cleared; each task that named old names id. The table given is not changed.
func redoRegister(tasks []Task, old, id string, issue int, opened time.Time) ([]Task, error) {
	at := slices.IndexFunc(tasks, func(t Task) bool { return t.ID == old })
	if at < 0 {
		return nil, fmt.Errorf("%s has no row of %s", tasksPath, old)
	}
	if slices.ContainsFunc(tasks, func(t Task) bool { return t.ID == id }) {
		return nil, fmt.Errorf("%s is already a task of %s", id, tasksPath)
	}
	out := make([]Task, len(tasks))
	for i, t := range tasks {
		t.Predecessors = slices.Clone(t.Predecessors)
		for j, p := range t.Predecessors {
			if p == old {
				t.Predecessors[j] = id
			}
		}
		out[i] = t
	}
	out[at] = Task{ID: id, Issue: issue, Kind: tasks[at].Kind, Requirements: slices.Clone(tasks[at].Requirements), Size: tasks[at].Size,
		Predecessors: slices.Clone(tasks[at].Predecessors), Opened: opened.UTC().Truncate(time.Second)}
	if err := CheckRegister(out); err != nil {
		return nil, err
	}
	return out, nil
}

// A redoPlan is what the combined records commit of a redo is made of.
type redoPlan struct {
	Tasks     []Task            // the register before
	Old, New  string            // the task IDs
	Issue     int               // the new task's issue
	Attempt   int               // the open attempt of the old task
	Base      string            // the base of the new task
	At        time.Time         // the time of both events
	OldEvents [][]string        // the events of the old task
	Bindings  map[string][]byte // tasks/<old>/options/<decision>.tsv, by the name below tasks/<old>/options
	Pin       []byte            // tasks/<old>/pin.txt of a Shape task, nil for another
}

// redoFiles gives the files of the combined records commit of a redo, each read back by its own reader: the register with the row replaced,
// the events of the old task with closed, the events of the new task with rebased, and the bindings.
func redoFiles(p redoPlan) (map[string][]byte, error) {
	tasks, err := redoRegister(p.Tasks, p.Old, p.New, p.Issue, p.At)
	if err != nil {
		return nil, err
	}
	var reg bytes.Buffer
	if err := WriteRegister(&reg, tasks); err != nil {
		return nil, err
	}
	if !slices.ContainsFunc(p.OldEvents, func(e []string) bool { return e[1] == "attempt" && e[2] == strconv.Itoa(p.Attempt) }) {
		return nil, fmt.Errorf("the events of %s have no attempt %d", p.Old, p.Attempt)
	}
	at := p.At.UTC().Format(timeForm)
	closed := append(slices.Clone(p.OldEvents), []string{strconv.Itoa(len(p.OldEvents) + 1), "closed", strconv.Itoa(p.Attempt), "", "", "", p.New, at})
	if err := records.CheckEventsAppend(p.OldEvents, closed); err != nil {
		return nil, err
	}
	files := map[string][]byte{tasksPath: reg.Bytes()}
	for path, rows := range map[string][][]string{eventsPath(p.Old): closed, eventsPath(p.New): {{"1", "rebased", "1", "", p.Base, "", p.Old, at}}} {
		data, err := table(records.EventsSchema, rows)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadEvents(data); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		files[path] = data
	}
	if p.Pin != nil {
		files[pinPath(p.New)] = p.Pin
	}
	for name, data := range p.Bindings {
		if _, err := spec.ReadOptions(data); err != nil {
			return nil, fmt.Errorf("the binding %s: %w", name, err)
		}
		files["tasks/"+p.New+"/options/"+name] = data
	}
	return files, nil
}

var intentForm = regexp.MustCompile(`^(T-[a-z0-9]{4}) ([0-9a-f]{40})\n$`)

// intentBytes and parseIntent are the file tasks/<old>/redo.txt: the new task's ID, a space, its base, and a line feed.
func intentBytes(id, base string) []byte { return []byte(id + " " + base + "\n") }

func parseIntent(data []byte) (id, base string, err error) {
	m := intentForm.FindSubmatch(data)
	if m == nil {
		return "", "", fmt.Errorf("the intent of a redo is a task ID, a space, a base SHA and a line feed; got %q", data)
	}
	return string(m[1]), string(m[2]), nil
}

// firstBase gives the base of the first attempt of a task: the base of its event rebased, else def.
func firstBase(events [][]string, def string) string {
	for _, e := range events {
		if e[1] == "rebased" {
			return e[4]
		}
	}
	return def
}

// redoIssue gives the title and the body of the new issue from what the records keep of the old task: the kind's composer with the new ID,
// the requirement IDs of the row, the decision and the option of the binding of an architecture task that a bet opened, and, for a Shape task,
// the opening pin of a Shape task (the digest of the draft and the records commit its issue was opened with, kept in tasks/<old>/pin.txt, and inherited by the replacement). A kind with no composer, or a Shape task with no pin, is an error.
func (s *Sessions) redoIssue(ctx context.Context, row Task, id string, bindings map[string][]byte) (title, body string, err error) {
	switch {
	case row.Kind == "architecture" && len(bindings) > 1:
		return "", "", fmt.Errorf("%s has %d bindings; a task of one decision has one", row.ID, len(bindings))
	case row.Kind == "architecture" && len(bindings) == 1:
		for name, data := range bindings {
			decision, _ := strings.CutSuffix(name, ".tsv")
			rows, err := spec.ReadOptions(data)
			if err != nil {
				return "", "", err
			}
			if len(rows) == 0 {
				return "", "", fmt.Errorf("the binding %s has no row", name)
			}
			return archIssueTitle(id, decision), archIssueBody(id, decision, rows[0][0], row.Requirements), nil
		}
	case row.Kind == "specification" || row.Kind == "architecture":
		pin, err := s.Store.ReadFile(ctx, pinPath(row.ID))
		if err != nil {
			return "", "", err
		}
		if pin == nil {
			return "", "", fmt.Errorf("%s has no opening pin (%s): its issue cannot be composed again", row.ID, pinPath(row.ID))
		}
		sum, at, err := parsePin(pin)
		if err != nil {
			return "", "", fmt.Errorf("%s: %w", pinPath(row.ID), err)
		}
		return shapeIssueTitle(id, row.Kind), shapeIssueText(id, row.Kind, row.Requirements, sum, at), nil
	}
	return "", "", fmt.Errorf("the kind %s of %s has no composer of an issue: a task of that kind is not redone in M2e", row.Kind, row.ID)
}

// redo redoes the task as a new task, in the order of loop.md, A task whose base moved: the pull requests of the attempts are read (a merged
// one is an error), the combined commit is built and read back by its readers with a placeholder issue, the intent is committed, and then
// each pull request that is open is closed, the old issue is closed and the new one opened, each after the lease is read, and the one
// records commit is made last. With an intent given the decision was made: its new ID and base are used, whatever the comparison would say now.
func (s *Sessions) redo(ctx context.Context, in CarryInput, events [][]string, attempt int, d string, intent []byte) (Advance, error) {
	task := in.Task
	row, err := s.taskRow(ctx, task)
	if err != nil {
		return Advance{}, err
	}
	var open []forge.PullRequest
	for a := 1; a <= attempt; a++ {
		pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: fmt.Sprintf("task/%s/%d", task, a)})
		if err != nil {
			return Advance{}, err
		}
		if found && pr.Merged {
			return Advance{}, fmt.Errorf("%s: the pull request %d merged: it is recorded, not redone", task, pr.Number)
		}
		if found && pr.State == "open" {
			open = append(open, pr)
		}
	}
	id := ""
	if intent != nil {
		if id, d, err = parseIntent(intent); err != nil {
			return Advance{}, fmt.Errorf("%s: %w", intentPath(task), err)
		}
	} else if id, err = s.taskID(ctx); err != nil {
		return Advance{}, err
	}
	data, err := s.Store.ReadFile(ctx, tasksPath)
	if err != nil {
		return Advance{}, err
	}
	register, err := ReadRegister(data)
	if err != nil {
		return Advance{}, fmt.Errorf("%s: %w", tasksPath, err)
	}
	bindings, err := s.Store.ReadFiles(ctx, "tasks/"+task+"/options")
	if err != nil {
		return Advance{}, err
	}
	var pin []byte
	if row.Kind != "architecture" || len(bindings) == 0 {
		if pin, err = s.Store.ReadFile(ctx, pinPath(task)); err != nil {
			return Advance{}, err
		}
	}
	plan := redoPlan{Tasks: register, Old: task, New: id, Issue: 1, Attempt: attempt, Base: d, At: s.Now(), OldEvents: events, Bindings: bindings, Pin: pin}
	if _, err := redoFiles(plan); err != nil {
		return Advance{}, err
	}
	title, body, err := s.redoIssue(ctx, row, id, bindings)
	if err != nil {
		return Advance{}, err
	}
	if intent == nil {
		if err := s.commit(ctx, map[string][]byte{intentPath(task): intentBytes(id, d)}, fmt.Sprintf("layup run: the redo of %s as %s decided", task, id)); err != nil {
			return Advance{}, err
		}
	}
	for _, pr := range open {
		if err := s.holdsLease(ctx); err != nil {
			return Advance{}, err
		}
		if got, err := s.Forge.ClosePullRequest(ctx, pr.Number); err != nil {
			return Advance{}, err
		} else if got.Merged {
			return Advance{}, fmt.Errorf("%s: the pull request %d merged after its redo", task, pr.Number)
		}
	}
	if err := s.holdsLease(ctx); err != nil {
		return Advance{}, err
	}
	if err := s.Forge.CloseIssue(ctx, row.Issue); err != nil {
		return Advance{}, err
	}
	if err := s.holdsLease(ctx); err != nil {
		return Advance{}, err
	}
	if plan.Issue, err = s.Forge.OpenIssue(ctx, title, body); err != nil {
		return Advance{}, err
	}
	files, err := redoFiles(plan)
	if err != nil {
		return Advance{}, err
	}
	if err := s.commit(ctx, files, fmt.Sprintf("layup run: %s redone as %s", task, id)); err != nil {
		return Advance{}, err
	}
	return Advance{Task: task, Replaced: id}, nil
}

// closeOlderPulls closes the open pull request of each attempt older than the open one (O22), only when the register names the pull request of
// the open attempt's branch: a new attempt with no pull request leaves the older one open. A task with one attempt, or no row or no pull
// request in the register, needs no forge.
func (s *Sessions) closeOlderPulls(ctx context.Context, task string) error {
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return err
	}
	attempt, _ := latestAttempt(events)
	if attempt < 2 {
		return nil
	}
	data, err := s.Store.ReadFile(ctx, tasksPath)
	if err != nil || data == nil {
		return err
	}
	register, err := ReadRegister(data)
	if err != nil {
		return fmt.Errorf("%s: %w", tasksPath, err)
	}
	i := slices.IndexFunc(register, func(t Task) bool { return t.ID == task })
	if i < 0 || register[i].PR == 0 {
		return nil
	}
	if cur, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: fmt.Sprintf("task/%s/%d", task, attempt)}); err != nil || !found || cur.Number != register[i].PR {
		return err
	}
	for a := 1; a < attempt; a++ {
		pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: fmt.Sprintf("task/%s/%d", task, a)})
		if err != nil {
			return err
		}
		if !found || pr.Merged || pr.State != "open" {
			continue
		}
		if err := s.holdsLease(ctx); err != nil {
			return err
		}
		if got, err := s.Forge.ClosePullRequest(ctx, pr.Number); err != nil {
			return err
		} else if got.Merged {
			return errors.New(task + ": the pull request " + strconv.Itoa(pr.Number) + " merged as it was closed")
		}
	}
	return nil
}

// pinPath is the file of the opening pin of a Shape task (records.md: tasks/<task>/pin.txt); pinBytes and parsePin are its bytes:
// the SHA-256 of the draft and the records commit that the issue was opened with, a space between and a line feed.
func pinPath(task string) string { return "tasks/" + task + "/pin.txt" }

func pinBytes(sum, records string) []byte { return []byte(sum + " " + records + "\n") }

var pinForm = regexp.MustCompile(`^([^ \n]+) ([^ \n]+)\n$`)

func parsePin(data []byte) (sum, records string, err error) {
	m := pinForm.FindSubmatch(data)
	if m == nil {
		return "", "", fmt.Errorf("a pin is the SHA-256 of the draft, a space, the records commit and a line feed; got %q", data)
	}
	return string(m[1]), string(m[2]), nil
}
