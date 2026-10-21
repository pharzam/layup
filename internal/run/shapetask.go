package run

// This file holds the opening of a Shape task (docs/spec/loop.md, The task register; task T-jy6f,
// #415, row 66r of the plan of M2e).

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// shapeIssueBlocks is the body of the issue of each kind of Shape task: the blocks shape-issue-<kind> of
// loop.md, with the placeholders {task}, {ids}, {sha} and {records}.
var shapeIssueBlocks = map[string]string{
	"specification": `LAYUP task {task}: the specification of the target.

Requirement IDs: {ids}
Pinned draft: spec/draft.tsv, SHA-256 {sha}, read at the records commit {records}.

Definition of Done: docs/prd/ holds a row of the requirement table, an acceptance row and a matrix row for each requirement ID, with the MoSCoW and Phase cells empty; docs/spec/ holds a section for each requirement ID; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
`,
	"architecture": `LAYUP task {task}: the target's architecture.

Requirement IDs: {ids}
Pinned draft: spec/draft.tsv, SHA-256 {sha}, read at the records commit {records}.

Definition of Done: docs/spec/ and docs/adr/ hold each decision that has one survivor as decided, and a section "Decision pending the bet" for each decision that has more than one; docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
`,
}

// shapeIssueTitles is the part of the title that follows "LAYUP task <ID>: ".
var shapeIssueTitles = map[string]string{"specification": "the specification", "architecture": "the target's architecture"}

// finalID gives the stable requirement ID of a row R-NNN of the draft by its kind: REQ-NNN for a
// requirement, NFR-NNN for a non-functional row, the same number; a success criterion has none.
func finalID(id, kind string) (string, bool) {
	n, ok := strings.CutPrefix(id, "R-")
	switch {
	case !ok:
		return "", false
	case kind == "requirement":
		return "REQ-" + n, true
	case kind == "non-functional":
		return "NFR-" + n, true
	}
	return "", false
}

// finalIDs gives the IDs of every row of the draft that has one, in the order of the draft.
func finalIDs(draft [][]string) []string {
	var out []string
	for _, r := range draft {
		if id, ok := finalID(r[0], r[1]); ok {
			out = append(out, id)
		}
	}
	return out
}

// shapeIssueTitle is the title of the issue of a Shape task.
func shapeIssueTitle(id, kind string) string {
	return "LAYUP task " + id + ": " + shapeIssueTitles[kind]
}

// shapeIssueText is the body of the issue of a Shape task: the block of its kind in loop.md, filled.
func shapeIssueText(id, kind string, ids []string, sum, records string) string {
	return strings.NewReplacer("{task}", id, "{ids}", strings.Join(ids, " "), "{sha}", sum, "{records}", records).Replace(shapeIssueBlocks[kind])
}

// openShapeTask opens the Shape task of the kind once (loop.md, The opening of a Shape task): an
// existing row of the kind is returned as it is; else the permissions, the draft and the commit that
// will register the task are checked before the lease is read and the issue is opened.
func (r *state) openShapeTask(ctx context.Context, kind string) (Task, error) {
	return r.openTask(ctx, openSpec{kind: kind, reuse: true, title: func(id string) string { return shapeIssueTitle(id, kind) },
		body: func(id string, ids []string, sum, records string) string {
			return shapeIssueText(id, kind, ids, sum, records)
		}})
}

// openSpec parameterizes the opening of a task (task T-2vva: the architecture tasks that the bet opens): the kind, whether an existing
// row of the kind is returned as it is (the Shape tasks), the requirement IDs (the draft's when none), the predecessors, the title and the
// body of the issue, the files of the same commit that the task needs (its binding), and a call between the issue and the row.
type openSpec struct {
	id           string // the reserved ID of the task, drawn when empty
	kind         string
	reuse        bool
	ids          []string
	predecessors []string
	title        func(id string) string
	body         func(id string, ids []string, sum, records string) string
	extra        func(id string) map[string][]byte
	afterIssue   func() error
}

// openTask opens a task of the kind in the order of loop.md, The opening of a Shape task: the candidate row, the parameters and the
// extra files are checked in a dry run before the lease is read and the issue is opened.
func (r *state) openTask(ctx context.Context, o openSpec) (Task, error) {
	var files [3][]byte
	for i, p := range []string{tasksPath, transitionsPath, parametersPath} {
		var err error
		if files[i], err = r.store.ReadFile(ctx, p); err != nil {
			return Task{}, err
		}
	}
	tasks, transitions, parameters := files[0], files[1], files[2]
	if tasks != nil {
		register, err := ReadRegister(tasks)
		if err != nil {
			return Task{}, fmt.Errorf("%s: %w", tasksPath, err)
		}
		for _, t := range register {
			if o.reuse && t.Kind == o.kind {
				return t, nil
			}
		}
	}
	if s := r.loopPermissionsStep(); s.Result == "fail" {
		return Task{}, errors.New(s.Detail)
	}
	const draftPath = "spec/draft.tsv"
	draft, err := r.store.ReadFile(ctx, draftPath)
	if err != nil {
		return Task{}, err
	}
	if draft == nil {
		return Task{}, fmt.Errorf("the records have no %s", draftPath)
	}
	rows, err := records.ReadDraft(draft)
	if err != nil {
		return Task{}, fmt.Errorf("%s: %w", draftPath, err)
	}
	ids := finalIDs(rows)
	if o.ids != nil {
		ids = o.ids
	}
	if len(ids) == 0 {
		return Task{}, fmt.Errorf("%s has no row of kind requirement or non-functional", draftPath)
	}
	s := r.sessions()
	id := o.id
	if id == "" {
		var err error
		if id, err = s.taskID(ctx); err != nil {
			return Task{}, err
		}
	}
	t := Task{ID: id, Kind: o.kind, Requirements: ids, Predecessors: o.predecessors, Opened: r.cfg.Clock.Now().UTC().Truncate(time.Second)}
	// The commit that will register the task is built and refused here, with a placeholder for the issue,
	// before the forge is written; registerTask builds it again with the issue.
	dry := t
	dry.Issue = 1
	dryFiles, err := registerFiles(tasks, transitions, dry)
	if err != nil {
		return Task{}, err
	}
	var extra map[string][]byte
	if o.extra != nil {
		if extra = o.extra(id); extra == nil {
			return Task{}, fmt.Errorf("the task %s has no extra file", id)
		}
		for p := range extra {
			if _, clash := dryFiles[p]; clash {
				return Task{}, fmt.Errorf("the extra file %s is a file of the register", p)
			}
		}
	}
	if _, err := parameterFiles(parameters, r.value, t.Opened); err != nil {
		return Task{}, err
	}
	// The opening inputs are captured once: the issue and, for a Shape task, the kept pin (for the redo, task T-89xf) are made from these
	// values, whatever the records do before the issue is opened (a heartbeat moves the base).
	sum, opening := digest(draft), r.store.Base()
	if o.reuse {
		if extra == nil {
			extra = map[string][]byte{}
		}
		extra[pinPath(id)] = pinBytes(sum, opening)
	}
	if err := s.holdsLease(ctx); err != nil {
		return Task{}, err
	}
	if t.Issue, err = r.cfg.Forge.OpenIssue(ctx, o.title(id), o.body(id, ids, sum, opening)); err != nil {
		return Task{}, err
	}
	if o.afterIssue != nil {
		if err := o.afterIssue(); err != nil {
			return Task{}, err
		}
	}
	return t, r.registerTaskWith(ctx, t, extra)
}
