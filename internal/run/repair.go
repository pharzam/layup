package run

// This file holds the repairs of a task session at a restart (docs/spec/loop.md,
// REQ-005 — The steps of the loop after the Scaffold; row 66v of the plan, task
// T-syxm, #422): a push with no result, and the comment of a session on the
// task's issue. They read the events, the sessions and the register of the
// records and the forge, never bets.tsv, so they work during Shape.

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/pharzam/layup/internal/records"
)

// Repair is what RepairTask wrote: the sessions it bound, the sessions whose push
// it refused as push-lost, the sessions whose comment it posted, and the new
// attempt it started (0 for none).
type Repair struct {
	Bound, Lost, Posted []string
	Attempt             int
}

// RepairTask repairs the sessions of task before a driver reads its position:
// the pushes first, so that a refusal reaches the comment, then the pull request of an older attempt (O22, row 66x), then the comments.
// Each stage keeps what an earlier stage committed; a stage that fails stops the
// repair with what it wrote so far, and the next run repeats it from the records.
func (s *Sessions) RepairTask(ctx context.Context, task string) (Repair, error) {
	out, err := s.RepairPushes(ctx, task)
	if err != nil {
		return out, err
	}
	if err := s.closeOlderPulls(ctx, task); err != nil {
		return out, err
	}
	out.Posted, err = s.RepostComments(ctx, task)
	return out, err
}

// RepairPushes resolves each push of task that has no bound and no refused after
// it: the branch of the push event, read from the forge, holds the pushed SHA
// exactly (bound), or it does not (refused, push-lost: absent, another SHA, a
// descendant of the pushed one included); a read that fails is an error and
// writes nothing. A push of an attempt that is not the open one gets the same
// event and no new attempt. Then, when an open attempt holds a push-lost
// refusal, one new attempt starts at its base: the new attempt holds none, so
// the rule starts one however many pushes were lost and however often it runs,
// and finds a refusal that an interrupted run left.
func (s *Sessions) RepairPushes(ctx context.Context, task string) (Repair, error) {
	s.fns()
	var out Repair
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return out, err
	}
	for i, e := range events {
		if e[1] != "push" || e[3] == "" || resolved(events[i+1:], e[3]) { // a push of code has no session: its step resolves it
			continue
		}
		id, branch, pushed := e[3], e[6], e[5]
		attempt, _ := strconv.Atoi(e[2])
		head, found, err := s.Store.BranchHead(ctx, branch)
		if err != nil {
			return out, err
		}
		kind, message := "refused push-lost", fmt.Sprintf("layup run: the push of %s was lost", id)
		if found && head == pushed {
			kind, message = "bound "+pushed+" "+branch, fmt.Sprintf("layup run: %s bound to %s", id, branch)
		}
		data, err := s.appendEvents(ctx, task, attempt, id, []string{kind})
		if err != nil {
			return out, err
		}
		if err := s.commit(ctx, map[string][]byte{eventsPath(task): data}, message); err != nil {
			return out, err
		}
		if found && head == pushed {
			out.Bound = append(out.Bound, id)
		} else {
			out.Lost = append(out.Lost, id)
		}
	}
	events, err = s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return out, err
	}
	for _, e := range events {
		attempt, _ := strconv.Atoi(e[2])
		if e[1] != "refused" || e[6] != "push-lost" || !openAttempt(events, e[3], attempt) {
			continue
		}
		base := ""
		for _, a := range events {
			if a[1] == "attempt" && a[2] == e[2] {
				base = a[4]
			}
		}
		if base == "" {
			return out, errors.New("the events of " + task + " have no attempt " + e[2])
		}
		attempt, err = s.StartAttempt(ctx, task, base)
		if err == nil {
			out.Attempt = attempt
		}
		return out, err
	}
	return out, nil
}

// resolved reports whether the events after a push hold a bound or a refused of its session.
func resolved(after [][]string, id string) bool {
	for _, e := range after {
		if e[3] == id && (e[1] == "bound" || e[1] == "refused") {
			return true
		}
	}
	return false
}

// RepostComments posts the comment of each session of task that has a result and
// whose line the App has not posted on the task's issue, from the events, in the
// order of sessions.tsv; it gives the sessions that it posted. A line that a
// repair changed is a new comment, and the old one stays. A task with no such
// session reads neither the register nor the forge. It needs the App's slug, the
// only way to tell the App's own comment.
func (s *Sessions) RepostComments(ctx context.Context, task string) ([]string, error) {
	s.fns()
	sessions, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return nil, err
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	var lines, ids []string
	for _, r := range sessions {
		if line := sessionLine(r, events); r[sessionTask] == task && line != "" {
			lines, ids = append(lines, line), append(ids, r[sessionID])
		}
	}
	if len(lines) == 0 {
		return nil, nil
	}
	if s.App == "" {
		return nil, errors.New("the App's slug is unknown, so its own comments cannot be told")
	}
	posted, err := s.postSessionComments(ctx, task, lines)
	var out []string
	for _, p := range posted {
		for i, l := range lines {
			if l == p {
				out = append(out, ids[i])
			}
		}
	}
	return out, err
}
