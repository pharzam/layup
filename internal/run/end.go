package run

// This file holds the end of a task session (docs/spec/session.md, REQ-005 —
// The result of a session, The open attempt, A comment for a session; row
// 36b of the plan, task T-fsjp, #165): the result file, the artifacts at the
// session's head, the open attempt, the telemetry row in the one records
// commit, and the comment on the control issue.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
)

// openAttempt reports whether the attempt of the session id is still its
// task's open attempt: after the session's event session, no event closed or
// rebased of that attempt, and no event attempt of another (§3).
func openAttempt(events [][]string, id string, attempt int) bool {
	a := strconv.Itoa(attempt)
	after := false
	for _, e := range events {
		switch {
		case !after:
			after = e[1] == "session" && e[3] == id
		case (e[1] == "closed" || e[1] == "rebased") && e[2] == a, e[1] == "attempt" && e[2] != a:
			return false
		}
	}
	return after
}

// The columns of a row of sessions.tsv, as records.ReadSessions gives it.
const (
	sessionID = iota
	sessionTask
	sessionAttempt
	sessionRole
	sessionHarness
	sessionVersion
	sessionModel
	sessionBase
)

// endTask is the end of a task session id in d, whose process gave r and
// err: it reads the result file (for the class done only, decided here), the
// start row, and the events at the run's last pushed records commit; it
// checks the open attempt, and for a kept result the head (read from the files
// of repo/.git, fetched into the run's clone with FetchSession, so no command
// of layup runs in repo/) and each artifact's SHA-256 there; it makes one
// records commit of the result file, the events result and refused, and the
// telemetry row, through Fenced; then it posts the comment on the control
// issue. An error of the process that is no class is the run's own.
func (s *Sessions) endTask(ctx context.Context, spec TaskSpec, id string, d session.Dir, r session.Run, err error) error {
	data, rows, resultErr := session.ResultFile(d, false)
	class := session.Class(r, err, resultErr)
	if class == "" {
		return err
	}
	start, err := s.startRowOf(ctx, id)
	if err != nil {
		return err
	}
	attempt, _ := strconv.Atoi(start[sessionAttempt])
	task := start[sessionTask]
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return err
	}
	var refused []string // the reasons, in order
	keep := class == "done"
	switch {
	case !openAttempt(events, id, attempt):
		refused, keep = append(refused, "closed-attempt"), false
	case keep:
		bad, err := s.artifacts(d, task, attempt, id, rows)
		switch {
		case errors.As(err, new(session.Refusal)):
			refused, keep = append(refused, "branch"), false
		case err != nil:
			return err
		case bad:
			refused = append(refused, "artifact")
		}
	}
	files := map[string][]byte{}
	if keep {
		files["tasks/"+task+"/results/"+id+".tsv"] = data
	}
	evData, err := s.appendEvents(ctx, task, attempt, id, append([]string{"result " + class}, prefix("refused ", refused)...))
	if err != nil {
		return err
	}
	files[eventsPath(task)] = evData
	tel, err := s.telemetry(ctx, spec, start, d, r)
	if err != nil {
		return err
	}
	files["telemetry.tsv"] = tel
	if err := s.commit(ctx, files, fmt.Sprintf("layup run: the end of %s, %s", id, class)); err != nil {
		return err
	}
	line := fmt.Sprintf("%s: %s session of %s, attempt %d: %s", id, start[sessionRole], task, attempt, class)
	for _, why := range refused {
		line += ", refused " + why
	}
	_, err = s.Forge.Comment(ctx, s.Control, line+"\n")
	return err
}

func prefix(p string, list []string) []string {
	out := make([]string, len(list))
	for i, v := range list {
		out[i] = p + v
	}
	return out
}

// startRowOf gives the row of sessions.tsv of the session id: the session
// ID, the task, the role, the attempt and the base come from it, never from
// the result file (§3).
func (s *Sessions) startRowOf(ctx context.Context, id string) ([]string, error) {
	rows, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		if r[sessionID] == id {
			return r, nil
		}
	}
	return nil, fmt.Errorf("sessions.tsv has no start row of %s", id)
}

// appendEvents gives events.tsv of task with each event "kind detail" added.
func (s *Sessions) appendEvents(ctx context.Context, task string, attempt int, id string, kinds []string) ([]byte, error) {
	before, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	after := append([][]string{}, before...)
	for _, k := range kinds {
		kind, detail, _ := strings.Cut(k, " ")
		after = append(after, []string{strconv.Itoa(len(after) + 1), kind, strconv.Itoa(attempt), id, "", "", detail, s.Now().UTC().Format(timeForm)})
	}
	if err := records.CheckEventsAppend(before, after); err != nil {
		return nil, err
	}
	if err := records.CheckEvents(after); err != nil {
		return nil, err
	}
	return table(records.EventsSchema, after)
}

// artifacts reads the session's head from the files of repo/.git, fetches it
// into the run's clone, and reports whether an artifact of the result differs
// from the file of its path at the head, or the head lacks it. A head that
// cannot be read or fetched is Refusal{branch}.
func (s *Sessions) artifacts(d session.Dir, task string, attempt int, id string, rows [][]string) (bool, error) {
	head, err := session.HeadOf(d.Repo, task, attempt)
	if err != nil {
		return false, err
	}
	if err := git.FetchSession(s.Clone, filepath.Join(d.Repo, ".git"), head, "refs/layup/sessions/"+id); err != nil {
		return false, session.Refusal{Reason: "branch", Value: err.Error()}
	}
	for _, r := range rows {
		if r[0] != "artifact" {
			continue
		}
		blob, err := git.Show(s.Clone, head, r[2])
		if err != nil {
			return true, nil // the head lacks it
		}
		sum := sha256.Sum256(blob)
		if hex.EncodeToString(sum[:]) != r[3] {
			return true, nil
		}
	}
	return false, nil
}

// telemetry gives telemetry.tsv with the session's row added: the start
// row's columns, the register row's billing, the process's own times, the
// usage report of the row's format in stdout, and host:prices.tsv.
func (s *Sessions) telemetry(ctx context.Context, spec TaskSpec, start []string, d session.Dir, r session.Run) ([]byte, error) {
	before, err := s.readTable(ctx, "telemetry.tsv", records.ReadTelemetry)
	if err != nil {
		return nil, err
	}
	u, err := session.UsageOf(spec.Pair.Usage, filepath.Join(d.Root, "stdout"))
	if err != nil {
		u = session.Usage{Status: "unavailable", Reason: "stdout cannot be read"}
	}
	end := r.End
	if end.IsZero() {
		end = r.Start // a program that could not start has no end of its own
	}
	row, err := ledger.Row(ledger.Session{ID: start[sessionID], Task: start[sessionTask], Role: start[sessionRole],
		Harness: start[sessionHarness], Model: start[sessionModel], Billing: spec.Pair.Billing,
		Start: r.Start, FirstOutput: r.First, End: end}, ledger.Usage(u), s.Prices)
	if err != nil {
		return nil, err
	}
	return table(records.TelemetrySchema, append(append([][]string{}, before...), row))
}
