package run

// This file holds the status layup/verify (docs/spec/loop.md, layup/verify and
// The independence of the reviewers; task T-q0kr, #362, row 66f of the plan of M2e).

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// A VerifyInput is what the caller knows of the verification of a head: the task,
// the base of its change, the head, the verifier's session (none before it
// starts), whether it waits for its turn (no producer in M2e), and the harness
// and model registers.
type VerifyInput struct {
	Task, Base, Head, Session string
	Waiting                   bool
	Harnesses, Models         [][]string
}

// VerifyAt gives the SHA that layup/verify is posted on and its status, by the
// order of loop.md, layup/verify: no harness outside the authors is `not-active`;
// a verifier that waits is pending; the named session must be a session of the
// task and of the role of the verification, and independent (its harness is no
// author and its pair is admitted); then its events and its kept copy decide. The
// status of a verifier is posted on the head it reviewed, the base of its start row.
func (s *Sessions) VerifyAt(ctx context.Context, in VerifyInput) (string, Status, error) {
	if !forge.FullSHA(in.Base) || !forge.FullSHA(in.Head) {
		return "", Status{}, fmt.Errorf("the base and the head of layup/verify are full SHAs")
	}
	status := func(state, d string) Status { return Status{VerifyContext, state, d} }
	if _, _, err := s.ReviewerPair(ctx, in.Task, Verification, in.Base, in.Head, in.Harnesses, in.Models); errors.Is(err, ErrNoHarnessOutsideAuthors) {
		return in.Head, status("failure", "not-active: "+ErrNoHarnessOutsideAuthors.Error()), nil
	} else if err != nil {
		return "", Status{}, err
	}
	if in.Waiting {
		return in.Head, status("pending", "waits for its turn"), nil
	}
	if in.Session == "" {
		return "", Status{}, fmt.Errorf("no verifier session is named, and none waits")
	}
	start, err := s.startRowOf(ctx, in.Session)
	if err != nil {
		return "", Status{}, err
	}
	if start[sessionTask] != in.Task || start[sessionRole] != stepRole("verification") {
		return "", Status{}, fmt.Errorf("the session %s is a %s session of %s, not a verifier of %s", in.Session, start[sessionRole], start[sessionTask], in.Task)
	}
	reviewed := start[sessionBase]
	authors, err := s.Authors(ctx, in.Task, Verification, in.Base, reviewed)
	if err != nil {
		return "", Status{}, err
	}
	admitted, err := s.admission(ctx, in.Harnesses, in.Models)
	if err != nil {
		return "", Status{}, err
	}
	if slices.Contains(authors, start[sessionHarness]) || !admitted(start[sessionHarness], start[sessionModel]) {
		return reviewed, status("failure", "not-active: the verifier is not independent"), nil
	}
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return "", Status{}, err
	}
	ended := slices.ContainsFunc(events, func(e []string) bool { return e[3] == in.Session && (e[1] == "result" || e[1] == "refused") })
	if !ended {
		return reviewed, status("pending", "verifier running"), nil
	}
	noRecord := status("failure", "no review record")
	if !passed(events, in.Session) {
		return reviewed, noRecord, nil
	}
	data, err := s.Store.ReadFile(ctx, "tasks/"+in.Task+"/results/"+in.Session+".tsv")
	if err != nil || data == nil {
		return reviewed, noRecord, err
	}
	result, err := records.ReadResult(data)
	if err != nil {
		return "", Status{}, fmt.Errorf("tasks/%s/results/%s.tsv: %w", in.Task, in.Session, err)
	}
	file := ""
	for _, r := range result {
		if r[0] == "comment" {
			if file != "" {
				return reviewed, noRecord, nil
			}
			file = strings.TrimPrefix(r[2], "result/")
		}
	}
	copy, ok, err := s.keptComment(ctx, in.Task, in.Session, file, result)
	if err != nil || !ok {
		return reviewed, noRecord, err
	}
	rec, ok := ParseReviewRecord(copy)
	switch {
	case !ok:
		return reviewed, noRecord, nil
	case NothingMaterialOK(copy, reviewed):
		return reviewed, status("success", "nothing material in scope"), nil
	case rec.Head != reviewed:
		return reviewed, status("failure", "review record of another head"), nil
	}
	return reviewed, status("failure", fmt.Sprintf("material: %d findings", len(rec.Findings))), nil
}
