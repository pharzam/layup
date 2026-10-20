package run

// This file holds the close-out and the merge of a verified task (docs/spec/loop.md, The merge at the head SHA and the close-out
// commit; task T-ehb4, #425, row 66y of the plan of M2e): FinishTask takes a verified, up-to-date task to a recorded merge.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/session"
)

// A RenderedBet is the typed input of the three checks on a head that holds the rendered commit of the first bet: the copy of the answer and
// the rendered records.
type RenderedBet struct {
	Copy    []byte
	Records []rules.Record
}

// A FinishInput is what FinishTask is given for one task: the task, the run's records commit as it is now, the harness and model registers,
// and, for the specification task, the inputs of its rendered head (D10).
type FinishInput struct {
	Task              string
	Records           func() string
	Harnesses, Models [][]string
	Bet               func(ctx context.Context, verified string) (RenderedBet, error)
}

// maxCloseActs bounds the acts of the close-out in one call: the start of a session, the resume of its push, and the read that finds it bound.
const maxCloseActs = 3

// FinishTask takes a verified task to its recorded merge (loop.md, D1 to D10): it reads the pull request of the open attempt first and writes
// nothing for a pull request that was closed and not merged, repairs the task's sessions, checks the verified head and the base, posts the
// statuses of the rendered commit, runs the close-out session at the carried head, carries the verdict, posts the other statuses and merges at the
// head SHA. A merge, a wait, a stop and a handback (Attempt) are the results; a stop and a wait repeat on the next call.
func (s *Sessions) FinishTask(ctx context.Context, in FinishInput) (Advance, error) {
	s.fns()
	task := in.Task
	pr, attempt, events, err := s.attemptPull(ctx, task)
	if err != nil {
		return Advance{}, err
	}
	if adv, done, err := s.terminalPull(ctx, task, attempt, pr, events, false); done || err != nil {
		return adv, err
	}
	rep, err := s.RepairTask(ctx, task)
	if err != nil {
		return Advance{}, err
	}
	if rep.Attempt != 0 {
		return Advance{Task: task, Attempt: rep.Attempt}, nil
	}
	if pr, attempt, events, err = s.attemptPull(ctx, task); err != nil {
		return Advance{}, err
	}
	if adv, done, err := s.terminalPull(ctx, task, attempt, pr, events, true); done || err != nil {
		return adv, err
	}
	v, err := s.shapeView(ctx, task)
	if err != nil {
		return Advance{}, err
	}
	cut, cv := finishCut(v), v
	if cut >= 0 {
		cv.Events = v.Events[:cut+1]
	}
	st := loopPosition(cv)
	if cut < 0 || st.Pos != posVerified || st.Attempt != attempt {
		return Advance{}, fmt.Errorf("%s: attempt %d has no verified head", task, attempt)
	}
	cin := CarryInput{Task: task, Verifier: st.Session, Harnesses: in.Harnesses, Models: in.Models}
	hin, carried := HeadInput{Task: task, Base: st.Base}, st.Head
	if v.Kind == "specification" {
		rendered, err := s.renderedHead(ctx, in, v, cin, attempt, st.Head)
		if err != nil {
			return Advance{}, err
		}
		cin.Rendered, carried = &RenderedInput{Head: rendered.Head, Render: rendered.Records}, rendered.Head
		hin.Copy, hin.Records, hin.PreRev = rendered.Copy, rendered.Records, st.Head
	}
	if ok, err := s.upToDate(ctx, st.Base, carried); err != nil {
		return Advance{}, err
	} else if !ok {
		return Advance{}, fmt.Errorf("%s: the base moved (row 66x's)", task)
	}
	if cin.Rendered != nil {
		if _, err := s.CarryVerify(ctx, cin, ""); err != nil {
			return Advance{}, err
		}
		reg, err := registerRevOf(task, st.Head, v.Events, v.Sessions, "")
		if err != nil {
			return Advance{}, err
		}
		at := hin
		at.Head, at.RegisterRev = carried, reg
		if _, err := s.PostChecks(ctx, at); err != nil {
			return Advance{}, err
		}
	}
	return s.closeOut(ctx, in, cin, hin, attempt, carried)
}

// A renderedHead is the rendered commit of the first bet, with its typed inputs.
type renderedHead struct {
	Head string
	RenderedBet
}

// renderedHead gives the event bound of code in the attempt, which CheckRendered accepts above the verified head, with the inputs of D10.
func (s *Sessions) renderedHead(ctx context.Context, in FinishInput, v loopView, cin CarryInput, attempt int, verified string) (renderedHead, error) {
	var out renderedHead
	for _, e := range v.Events {
		if n, _ := strconv.Atoi(e[2]); e[1] == "bound" && e[3] == "" && n == attempt {
			out.Head = e[5]
		}
	}
	if out.Head == "" || in.Bet == nil {
		return out, fmt.Errorf("%s: the specification task has no rendered commit bound in attempt %d, or no inputs of it: the step bet finishes it", v.Task, attempt)
	}
	var err error
	if out.RenderedBet, err = in.Bet(ctx, verified); err != nil {
		return out, err
	}
	cin.Rendered = &RenderedInput{Head: out.Head, Render: out.Records}
	ok, err := s.CheckRendered(ctx, cin)
	if err == nil && !ok {
		err = fmt.Errorf("%s: %s is not a permitted rendered commit above the verified head %s", v.Task, out.Head, verified)
	}
	return out, err
}

// attemptPull gives the pull request of the branch of the task's latest attempt, the attempt and the events: its number must be the register's.
func (s *Sessions) attemptPull(ctx context.Context, task string) (forge.PullRequest, int, [][]string, error) {
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return forge.PullRequest{}, 0, nil, err
	}
	attempt := 0
	for _, e := range events {
		if n, _ := strconv.Atoi(e[2]); e[1] == "attempt" && n > attempt {
			attempt = n
		}
	}
	if attempt == 0 {
		return forge.PullRequest{}, 0, nil, fmt.Errorf("%s has no attempt", task)
	}
	row, err := s.taskRow(ctx, task)
	if err != nil {
		return forge.PullRequest{}, 0, nil, err
	}
	branch := fmt.Sprintf("task/%s/%d", task, attempt)
	pr, found, err := s.Forge.PullRequest(ctx, forge.PullRequestQuery{Branch: branch})
	switch {
	case err != nil:
		return pr, attempt, events, err
	case !found:
		return pr, attempt, events, fmt.Errorf("%s: the pull request of %s is not on the forge", task, branch)
	case pr.Number != row.PR:
		return pr, attempt, events, fmt.Errorf("%s: the pull request %d of %s is not the register's %d", task, pr.Number, branch, row.PR)
	}
	return pr, attempt, events, nil
}

// terminalPull gives the end of a pull request that is closed or merged: the wait of D2 for a closed one that was not merged, and, for a merged one,
// Merged when its event is written, or, after the repair, the event.
func (s *Sessions) terminalPull(ctx context.Context, task string, attempt int, pr forge.PullRequest, events [][]string, repaired bool) (Advance, bool, error) {
	switch {
	case pr.State == "closed" && !pr.Merged:
		return Advance{Task: task, Wait: closedWait(task, fmt.Sprintf("task/%s/%d", task, attempt), pr.Number)}, true, nil
	case !pr.Merged:
		return Advance{}, false, nil
	case slices.ContainsFunc(events, func(e []string) bool { return e[1] == "merged" && e[5] == pr.MergeSHA }):
		return Advance{Task: task, Merged: true, Head: pr.HeadSHA}, true, nil
	case !repaired:
		return Advance{}, false, nil
	}
	if err := s.RecordMerged(ctx, task, pr); err != nil {
		return Advance{}, true, err
	}
	return Advance{Task: task, Merged: true, Head: pr.HeadSHA}, true, nil
}

// closeOut runs the positions of the close-out session at the carried head and, once it is bound, carries the verdict, posts the checks and merges.
func (s *Sessions) closeOut(ctx context.Context, in FinishInput, cin CarryInput, hin HeadInput, attempt int, carried string) (Advance, error) {
	task := in.Task
	for i := 0; i < maxCloseActs; i++ {
		v, err := s.shapeView(ctx, task)
		if err != nil {
			return Advance{}, err
		}
		step := closeOutPosition(v, finishCut(v), attempt, carried)
		if step.Pos == closeResume {
			if _, err := git.RevParse(s.Clone, "refs/layup/sessions/"+step.Session); err != nil { // the run's clone is made anew at a restart (D5)
				step.Pos = closeStart
			}
		}
		stop := ""
		switch step.Pos {
		case closeStop:
			return Advance{Task: task, Stop: step.Msg}, nil
		case closeBound:
			return s.mergeAt(ctx, cin, hin, step)
		case closeResume:
			err = s.PushCloseOut(ctx, step.Session, cin)
		default:
			stop, err = s.startCloseOut(ctx, in, v, attempt, carried, cin)
		}
		var refusal session.Refusal
		switch {
		case stop != "":
			return Advance{Task: task, Stop: stop}, nil
		case errors.As(err, &refusal): // recorded: the next position names it
		case err != nil:
			return Advance{}, err
		}
	}
	return Advance{}, fmt.Errorf("%s: the close-out did not come to a bound head in %d acts", task, maxCloseActs)
}

// startCloseOut starts a close-out session at the carried head, at the tier of the plan and the position, whose push is the close-out check (D4).
func (s *Sessions) startCloseOut(ctx context.Context, in FinishInput, v loopView, attempt int, carried string, cin CarryInput) (string, error) {
	plan, err := s.planOf(ctx, v)
	if err != nil {
		return "", err
	}
	pos, err := s.Position(ctx, v.Task)
	if err != nil {
		return "", err
	}
	res, err := closeOutStep(plan, pos)
	if err != nil {
		return "", err
	}
	h, m, err := s.PairFor(ctx, res.Role, res.Tier, in.Harnesses, in.Models, nil)
	if err != nil {
		return pairStop(v.Task, "close-out", err)
	}
	push := func(id string, _ session.Dir, _ session.Run) error { return s.PushCloseOut(ctx, id, cin) }
	return "", s.startAs(ctx, ShapeTaskInput{Task: v.Task, Records: in.Records, Harnesses: in.Harnesses, Models: in.Models}, v.Kind, "close-out", res.Role, attempt, carried, h, m, nil, true, push)
}

// planOf gives the plan of the task: the kept copy of the last valid comment handoff of a session of the plan role.
func (s *Sessions) planOf(ctx context.Context, v loopView) (Plan, error) {
	role, path := map[string]string{}, ""
	for _, r := range v.Sessions {
		role[r[sessionID]] = r[sessionRole]
	}
	for _, e := range v.Events {
		if planHandoff(e, role) {
			path = copyPathOf(e)
		}
	}
	if path == "" {
		return Plan{}, fmt.Errorf("%s has no valid handoff of a plan", v.Task)
	}
	data, err := s.Store.ReadFile(ctx, path)
	if err != nil {
		return Plan{}, err
	}
	plan, _, ok := ParsePlanComment(data)
	if !ok {
		return Plan{}, fmt.Errorf("%s: the plan in %s does not parse", v.Task, path)
	}
	return plan, nil
}

// mergeAt carries the verdict to the bound head of the close-out, posts the other three statuses there, and merges at the head SHA (D6).
func (s *Sessions) mergeAt(ctx context.Context, cin CarryInput, hin HeadInput, step closeStep) (Advance, error) {
	if _, err := s.CarryVerify(ctx, cin, step.Session); err != nil {
		return Advance{}, err
	}
	hin.Head, hin.RegisterRev = step.Head, ""
	if _, err := s.PostChecks(ctx, hin); err != nil {
		return Advance{}, err
	}
	_, r, err := s.MergeTask(ctx, cin.Task, step.Head)
	if err != nil || r.Ready {
		return notReady(cin.Task, step.Head, r), err
	}
	pr, attempt, events, err := s.attemptPull(ctx, cin.Task)
	if err != nil {
		return Advance{}, err
	}
	if adv, done, err := s.terminalPull(ctx, cin.Task, attempt, pr, events, true); done || err != nil {
		return adv, err
	}
	return notReady(cin.Task, step.Head, r), nil
}

// closePos is a position of the close-out session.
type closePos int

const (
	closeStart  closePos = iota // no session of the position, or one with no result: start one
	closeResume                 // a session ended with no push: its checked push
	closeStop                   // refused, or ended not completed
	closeBound                  // bound at Head
)

// A closeStep is the position of the close-out with the session it concerns, the message of a stop and the bound head.
type closeStep struct {
	Pos           closePos
	Session, Head string
	Msg           string
}

// finishCut is the index of the last handoff of a session of the verification role in the events, -1 for none: what follows is the close-out's.
func finishCut(v loopView) int {
	cut := -1
	for i, e := range v.Events {
		if e[1] == "handoff" && v.role(e[3]) == stepRole("verification") {
			cut = i
		}
	}
	return cut
}

// closeOutPosition reads the close-out session of the attempt at the carried head from the events after the cut (D5): the last session of the close-out
// role of the attempt, started after the cut with the carried head as its base, decides; a start refused before its row is a stop.
func closeOutPosition(v loopView, cut, attempt int, carried string) closeStep {
	var ids []string
	for _, e := range v.Events[min(cut+1, len(v.Events)):] {
		i := slices.IndexFunc(v.Sessions, func(r []string) bool { return r[sessionID] == e[3] })
		switch {
		case e[1] == "refused" && e[3] != "" && i < 0:
			return closeStep{Pos: closeStop, Session: e[3], Msg: fmt.Sprintf("%s: the start of the session %s was refused (%s)", v.Task, e[3], e[6])}
		case e[1] == "session" && i >= 0 && v.Sessions[i][sessionRole] == stepRole("close-out") && v.Sessions[i][sessionAttempt] == strconv.Itoa(attempt) && v.Sessions[i][sessionBase] == carried:
			ids = append(ids, e[3])
		}
	}
	if len(ids) == 0 {
		return closeStep{Pos: closeStart}
	}
	last := ids[len(ids)-1]
	switch o, msg := v.outcome(last, "close-out"); o {
	case running:
		return closeStep{Pos: closeStart}
	case failed:
		return closeStep{Pos: closeStop, Session: last, Msg: msg}
	}
	pushed, head := false, ""
	for _, e := range v.Events {
		switch {
		case e[3] != last:
		case e[1] == "push":
			pushed = true
		case e[1] == "bound":
			head = e[5]
		}
	}
	switch {
	case head != "":
		return closeStep{Pos: closeBound, Session: last, Head: head}
	case pushed:
		return closeStep{Pos: closeStop, Session: last, Msg: fmt.Sprintf("%s: the push of the session %s has no result (its repair is RepairTask's)", v.Task, last)}
	}
	return closeStep{Pos: closeResume, Session: last}
}

// closeOutStep resolves the close-out at the tier of the plan and the position (D4).
func closeOutStep(plan Plan, pos Position) (Resolved, error) {
	return ResolveStep("close-out", TierOf(plan, pos), "")
}

// closedWait is the wait of a pull request closed and not merged (D2).
func closedWait(task, branch string, number int) string {
	return fmt.Sprintf("waiting for the Operator: %s: the pull request #%d of %s is closed and not merged", task, number, branch)
}

// notReady is the result of a merge call (D6): the merge, a stop for a head that moved, or a wait with the reason.
func notReady(task, head string, r Readiness) Advance {
	switch {
	case r.Ready:
		return Advance{Task: task, Merged: true, Head: head}
	case strings.HasPrefix(r.Reason, "the head moved to "):
		return Advance{Task: task, Stop: fmt.Sprintf("%s: not ready at %s: %s", task, head, r.Reason)}
	}
	return Advance{Task: task, Wait: fmt.Sprintf("waiting: %s: not ready at %s: %s", task, head, r.Reason)}
}

// upToDate tells whether the default branch's head, read from the forge, is the attempt's base or an ancestor of head, the head that carries the
// verdict (D9; row 66x consumes it with the head after its clean merge of the base). A failed read is an error, never a moved base.
func (s *Sessions) upToDate(ctx context.Context, base, head string) (bool, error) {
	tip, found, err := s.Store.BranchHead(ctx, s.Branch)
	switch {
	case err != nil:
		return false, err
	case !found:
		return false, fmt.Errorf("the default branch %s is not on the forge", s.Branch)
	case tip == base:
		return true, nil
	}
	ok, err := git.IsAncestor(s.Clone, tip, head)
	if err != nil {
		return false, fmt.Errorf("the default branch %s holds %s, not the base %s, and the run's clone cannot tell whether it is an ancestor of %s: %w", s.Branch, tip, base, head, err)
	}
	return ok, nil
}

// renderedBet rebuilds the typed inputs of the rendered head of the first bet from the records and the verified head (D10), as the step bet built
// them: the answer that the latest version of the inventory names, its copy, and the rendering of the PRD and of the numbered facts record.
func (r *state) renderedBet(ctx context.Context, verified string) (RenderedBet, error) {
	data, err := r.store.ReadFile(ctx, "bets.tsv")
	if err != nil {
		return RenderedBet{}, err
	}
	bets, err := ReadBets(data)
	if err != nil {
		return RenderedBet{}, fmt.Errorf("bets.tsv: %w", err)
	}
	i := slices.IndexFunc(bets, func(b []string) bool { return b[betVerdict] == "bet" })
	if i < 0 {
		return RenderedBet{}, errors.New("bets.tsv holds no bet: the step bet comes first")
	}
	milestone := bets[i][betMilestone]
	prds, err := filesAt(r.target, verified, "docs/prd", isPRDName)
	if err != nil {
		return RenderedBet{}, err
	}
	if len(prds) != 1 {
		return RenderedBet{}, fmt.Errorf("the verified head %s holds %d files docs/prd/PRD-*.md; want one", verified, len(prds))
	}
	var prd []byte
	for _, b := range prds {
		prd = b
	}
	basis, err := r.briefBasis(ctx, milestone)
	if err != nil {
		return RenderedBet{}, err
	}
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return RenderedBet{}, errors.New("start.tsv has no number of the Intake issue")
	}
	if r.copies, err = r.readRecord(ctx, "copies.tsv", records.ReadCopies); err != nil {
		return RenderedBet{}, err
	}
	copies, err := r.answerCopies(ctx, issue)
	if err != nil {
		return RenderedBet{}, err
	}
	applied, err := r.appliedAnswer(ctx, milestone)
	if err != nil || applied == "" {
		return RenderedBet{}, errors.Join(err, fmt.Errorf("the inventory names no answer of the bet %s", milestone))
	}
	answer, row, _, err := r.restoreAnswer(applied, copies, basis.input(phaseTagsOf(prd)))
	if err != nil {
		return RenderedBet{}, err
	}
	answer.Milestone = milestone
	facts, texts, err := r.factsOf(ctx)
	if err != nil {
		return RenderedBet{}, err
	}
	index, versions, err := r.inventoryOf(ctx)
	if err != nil {
		return RenderedBet{}, err
	}
	copyID, _ := strconv.Atoi(row[copyComment])
	files, version, err := nextVersion(index, versions, facts, answer, copyID)
	if err != nil {
		return RenderedBet{}, err
	}
	if files != nil {
		return RenderedBet{}, fmt.Errorf("the latest version of the inventory does not name the copy %d", copyID)
	}
	_, recs, err := r.renderedAt(verified, answer, version, texts)
	if err != nil {
		return RenderedBet{}, err
	}
	i = slices.IndexFunc(copies, func(c answerCopy) bool { return c.Row[copyComment] == applied })
	return RenderedBet{Copy: []byte(copies[i].Body), Records: recs}, nil
}
