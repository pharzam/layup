package run

// This file holds the checks before a push of a task session
// (docs/spec/session.md, REQ-003 — Before a push: the check of the base, and
// a workflow or rule-path change with its refused diff as a payload; row 37a
// of the plan, task T-z027, #166). The push and the bind are row 37b's.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/rules"
)

// sessionRecords is the column records of a row of sessions.tsv: the records
// commit that the prompt was built from.
const sessionRecords = sessionBase + 1

// passed reports whether the end of the session id passed its result: its
// event result is done, and it has no event refused (decided here, task
// T-z027: a done session refused at its end, artifact among them, is never
// checked or pushed).
func passed(events [][]string, id string) bool {
	done := false
	for _, e := range events {
		switch {
		case e[3] != id:
		case e[1] == "refused":
			return false
		case e[1] == "result":
			done = e[6] == "done"
		}
	}
	return done
}

// pushTask is the push hook of a task session id: for a session whose end
// passed, the checks before a push on the head that its end fetched into the
// run's clone, so no command of layup runs in repo/. A head that does not
// descend from the base is refused (base); a records commit with no register
// that reads is refused (rule-paths, decided here: fail closed, ADR-0017
// decision 2); a change of .github/workflows/ or of a rule path is refused
// (workflow, rule-path), its git diff --binary from the base to the head
// written as payloads/<sha256>, which the event names after the reason. Each
// refusal is one records commit through Fenced, and nothing is pushed.
func (s *Sessions) pushTask(ctx context.Context, id string) error {
	start, err := s.startRowOf(ctx, id)
	if err != nil {
		return err
	}
	task := start[sessionTask]
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil || !passed(events, id) {
		return err
	}
	head, err := git.RevParse(s.Clone, "refs/layup/sessions/"+id)
	if err != nil {
		return err
	}
	base := start[sessionBase]
	if ok, err := git.IsAncestor(s.Clone, base, head); err != nil {
		return err
	} else if !ok {
		return s.refuse(ctx, start, id, "base", nil)
	}
	// No register in the records commit refuses; a commit or a git that
	// cannot be read is the run's own error.
	entries, err := git.LsTree(s.Clone, start[sessionRecords], "rule-paths.tsv")
	if err != nil {
		return err
	}
	if len(entries) != 1 || entries[0].Path != "rule-paths.tsv" {
		return s.refuse(ctx, start, id, "rule-paths", nil)
	}
	data, err := git.Show(s.Clone, start[sessionRecords], "rule-paths.tsv")
	if err != nil {
		return err
	}
	reg, err := rules.ReadRegister(data)
	if err != nil {
		return s.refuse(ctx, start, id, "rule-paths", nil)
	}
	names, err := git.DiffNames(s.Clone, base, head)
	if err != nil {
		return err
	}
	var change rules.Change
	if slices.Contains(names, rules.GuardrailsPath) {
		if change.Diff, err = git.DiffFile(s.Clone, base, head, rules.GuardrailsPath); err != nil {
			return err
		}
		// A head that lacks the file keeps an empty Head, which no added
		// line passes: a deletion is a rule-path change.
		if entries, err := git.LsTree(s.Clone, head, rules.GuardrailsPath); err != nil {
			return err
		} else if len(entries) == 1 && entries[0].Path == rules.GuardrailsPath {
			if change.Head, err = git.Show(s.Clone, head, rules.GuardrailsPath); err != nil {
				return err
			}
		}
	}
	reason, _ := rules.Check(names, reg, change)
	if reason == "" {
		return nil // the push and the bind are row 37b's
	}
	diff, err := git.DiffBinary(s.Clone, base, head)
	if err != nil {
		return err
	}
	return s.refuse(ctx, start, id, reason, diff)
}

// refuse commits the event refused of the session with its reason, and with a
// payload its SHA-256 after one space and the file payloads/<sha256>.
func (s *Sessions) refuse(ctx context.Context, start []string, id, reason string, payload []byte) error {
	task := start[sessionTask]
	attempt, _ := strconv.Atoi(start[sessionAttempt])
	files := map[string][]byte{}
	detail := reason
	if payload != nil {
		sum := sha256.Sum256(payload)
		name := hex.EncodeToString(sum[:])
		files["payloads/"+name] = payload
		detail += " " + name
	}
	events, err := s.appendEvents(ctx, task, attempt, id, []string{"refused " + detail})
	if err != nil {
		return err
	}
	files[eventsPath(task)] = events
	return s.commit(ctx, files, fmt.Sprintf("layup run: %s refused before a push (%s)", id, reason))
}
