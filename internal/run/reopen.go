package run

// This file holds the specification task that a changed class reopens (docs/spec/loop.md, The specification task that a changed class
// reopens; task T-4naq, #428, row 66ab of the plan of M2e): after the event merged of the specification, one task for each fact that the first bet
// re-classed, under the identities of spec/reopenings.tsv.

import (
	"context"
	"crypto/rand"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// reopeningsPath is the record of the facts that the bet re-classed.
const reopeningsPath = "spec/reopenings.tsv"

// ReopeningsSchema is the Go value of the block reopenings; its owner is this package.
var ReopeningsSchema = tsv.Schema{Name: "reopenings", Location: "records:spec/reopenings.tsv", Columns: []tsv.Column{
	{Name: "n", Type: "int", Key: true},
	{Name: "fact", Type: "text", Key: true},
	{Name: "class", Type: "enum(need|constraint|measure)"},
	{Name: "token", Type: "text"},
	{Name: "requirement", Type: "text"},
	{Name: "task", Type: "id(T-xxxx)"},
	{Name: "copy", Type: "int"},
}}

// The columns of a row of spec/reopenings.tsv.
const (
	reoN = iota
	reoFact
	reoClass
	reoToken
	reoRequirement
	reoTask
	reoCopy
)

// reopenIssueBlock is the body of the issue of a reopened task: the block reopen-issue of loop.md.
const reopenIssueBlock = `LAYUP task {task}: the {what} of the fact {fact}, re-classed by the bet as {class}.

Requirement IDs: {ids}
Reopened after the merge of the specification task {spec}: the row of the fact {fact} in spec/reopenings.tsv of the records, the token {token}.

Definition of Done: docs/prd/ holds, under the heading "Waits for a bet" and not in a table, {line}; {sections}docs/tasks/{task}.md is the task's file. The head passes the four layup/ statuses, and a verifier on a harness outside the authors finds nothing material in scope.
`

// reopenedInstruction is the instruction of the step specification of a reopened task: the block instruction-specification-reopened of loop.md.
const reopenedInstruction = `Write the requirement that waits for a bet, for the task {task} (kind {kind}, issue {issue}), from the files listed above: the row of this task in spec/reopenings.tsv, the facts spec/spans.tsv, start/problem-statement.md and intake/answers.tsv, the revised draft spec/draft.tsv, the product requirements document in docs/prd/ and the specification in docs/spec/. The bet changed the class of the fact {fact} to {class}, and no requirement names it yet; the token {token} cites it. Commit your work in this repository in one or more commits; push nothing.
Keep every line of the files that already exist, except the additions below, renumber nothing, and put no new row in any table. In the file of docs/prd/ write, under a heading "## Waits for a bet" (add it as the last heading of the file when the file has none), the line of your class, as plain text starting at the first column, never in a table:
- need or constraint: "{requirement}: <the requirement, one sentence> Facts: {token} {record}. Acceptance criterion: <the acceptance criterion>." In docs/spec/ write a section whose heading holds {requirement} as a whole word, for example "## {requirement} — <title>", with a non-empty text under it.
- measure: "Success criterion: <the success criterion> Fact: {token}." Write no section in docs/spec/.
The class of this task is {class}: write the line of that class only.
Write docs/tasks/{task}.md, the task's file, in the form of the target's docs/tasks/. Write nothing under docs/facts/.
The required outputs for the result file are docs/tasks/{task}.md, the file of docs/prd/ that you wrote and, for a need or a constraint, each file of docs/spec/ that you wrote.
`

var (
	tokenOfNumber = regexp.MustCompile(`^F-[0-9]{4}#([0-9]+)$`)
	factNumber    = regexp.MustCompile(`^[0-9]+$`)
	reopenReq     = map[string]*regexp.Regexp{"need": regexp.MustCompile(`^REQ-[0-9]{3,}$`), "constraint": regexp.MustCompile(`^NFR-[0-9]{3,}$`)}
	knownReq      = regexp.MustCompile(`^(?:REQ|NFR)-([0-9]{3,})$`)
)

// ReadReopenings reads spec/reopenings.tsv by ReopeningsSchema and then makes the checks of its owner (docs/spec/records.md, the block reopenings):
// the version counts from 1 and its rows are together, the token is that of the fact, the requirement is that of the class, and no requirement
// and no task is held by two rows. The comparison with the inventory is bindReopenings's.
func ReadReopenings(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, ReopeningsSchema)
	if err != nil {
		return nil, err
	}
	reqs, tasks, seenN, last := map[string]int{}, map[string]int{}, map[int]bool{}, 0
	for i, r := range rows {
		line := i + 2
		bad := func(column, reason string) error { return &tsv.Error{Line: line, Column: column, Reason: reason} }
		for _, c := range []int{reoToken, reoTask, reoCopy} {
			if r[c] == "" {
				return nil, bad(ReopeningsSchema.Columns[c].Name, "this column never holds the empty value")
			}
		}
		n, _ := strconv.Atoi(r[reoN])
		switch {
		case n < 1:
			return nil, bad("n", "a version counts from 1")
		case n != last && seenN[n]:
			return nil, bad("n", fmt.Sprintf("the rows of the version %d are not together", n))
		}
		seenN[n], last = true, n
		if factNumber.MatchString(r[reoFact]) {
			if m := tokenOfNumber.FindStringSubmatch(r[reoToken]); m == nil || m[1] != r[reoFact] {
				return nil, bad("token", fmt.Sprintf("a token is F-NNNN#%s for the fact %s", r[reoFact], r[reoFact]))
			}
		} else if r[reoToken] != r[reoFact] {
			return nil, bad("token", fmt.Sprintf("a token is %s, the fact itself", r[reoFact]))
		}
		switch re := reopenReq[r[reoClass]]; {
		case r[reoClass] == "measure" && r[reoRequirement] != "":
			return nil, bad("requirement", "a measure has no requirement ID")
		case r[reoClass] != "measure" && (r[reoRequirement] == "" || !re.MatchString(r[reoRequirement])):
			return nil, bad("requirement", fmt.Sprintf("a %s's requirement is %s-NNN, three digits or more", r[reoClass], map[string]string{"need": "REQ", "constraint": "NFR"}[r[reoClass]]))
		}
		if r[reoRequirement] != "" {
			if j, ok := reqs[r[reoRequirement]]; ok {
				return nil, bad("requirement", fmt.Sprintf("the requirement %s is held by line %d", r[reoRequirement], j))
			}
			reqs[r[reoRequirement]] = line
		}
		if j, ok := tasks[r[reoTask]]; ok {
			return nil, bad("task", fmt.Sprintf("the task %s is held by line %d", r[reoTask], j))
		}
		tasks[r[reoTask]] = line
	}
	return rows, nil
}

// reopenedTasks gives the tasks that the rows name.
func reopenedTasks(rows [][]string) map[string]bool {
	out := map[string]bool{}
	for _, r := range rows {
		out[r[reoTask]] = true
	}
	return out
}

// reopeningOf gives the row of the task, nil for none.
func reopeningOf(rows [][]string, task string) []string {
	for _, r := range rows {
		if r[reoTask] == task {
			return r
		}
	}
	return nil
}

// reclassed selects the rows of a version of the inventory in scope whose class is a need, a constraint or a measure and differs from the class of
// the fact at Intake; a fact of the version that Intake does not know is an error.
func reclassed(version [][]string, intake []spec.InventoryFact) ([][]string, error) {
	class := map[string]string{}
	for _, f := range intake {
		class[f.ID] = f.Class
	}
	var out [][]string
	for _, r := range version {
		was, ok := class[r[0]]
		switch {
		case !ok:
			return nil, fmt.Errorf("the fact %s of the inventory is no fact of Intake", r[0])
		case r[2] == "in" && slices.Contains([]string{"need", "constraint", "measure"}, r[1]) && r[1] != was:
			out = append(out, r)
		}
	}
	return out, nil
}

// tokenOf is the citation of a fact: F-NNNN#<n> of the numbered facts record for a number, the fact itself for a Q-NNN.
func tokenOf(fact, record string) string {
	if factNumber.MatchString(fact) {
		return record + "#" + fact
	}
	return fact
}

// nextNumber is one more than the highest number of the REQ- and NFR- IDs given, 1 for none: one number space for both prefixes.
func nextNumber(ids ...[]string) int {
	top := 0
	for _, list := range ids {
		for _, id := range list {
			if m := knownReq.FindStringSubmatch(id); m != nil {
				n, _ := strconv.Atoi(m[1])
				top = max(top, n)
			}
		}
	}
	return top + 1
}

// allocate makes the rows of spec/reopenings.tsv of the selected rows of version n, in their order: the token, the requirement ID of a need
// (REQ-) and of a constraint (NFR-) from next on, none for a measure, and the task of the same place of tasks.
func allocate(sel [][]string, n int, record string, next int, tasks []string) [][]string {
	var out [][]string
	for i, r := range sel {
		req := ""
		switch r[1] {
		case "need":
			req, next = fmt.Sprintf("REQ-%03d", next), next+1
		case "constraint":
			req, next = fmt.Sprintf("NFR-%03d", next), next+1
		}
		out = append(out, []string{strconv.Itoa(n), r[0], r[1], tokenOf(r[0], record), req, tasks[i], r[3]})
	}
	return out
}

// bindReopenings compares the rows of the record with the selection of the version: the same facts in the same order, each with the class and the
// copy of the version and the token of the numbered facts record.
func bindReopenings(rows, sel [][]string, record string) error {
	if len(rows) != len(sel) {
		return fmt.Errorf("%s has %d facts; the selection of the version has %d facts", reopeningsPath, len(rows), len(sel))
	}
	for i, r := range rows {
		switch s := sel[i]; {
		case r[reoFact] != s[0]:
			return fmt.Errorf("%s names the fact %s at the row %d; the selection has the fact %s", reopeningsPath, r[reoFact], i+1, s[0])
		case r[reoClass] != s[1]:
			return fmt.Errorf("%s gives the fact %s the class %s; the version has %s", reopeningsPath, s[0], r[reoClass], s[1])
		case r[reoCopy] != s[3]:
			return fmt.Errorf("%s gives the fact %s the copy %s; the version has %s", reopeningsPath, s[0], r[reoCopy], s[3])
		case r[reoToken] != tokenOf(s[0], record):
			return fmt.Errorf("%s gives the fact %s the token %s; the numbered facts record gives %s", reopeningsPath, s[0], r[reoToken], tokenOf(s[0], record))
		}
	}
	return nil
}

// firstVersion gives the rows of the one version of the inventory after its sha256 is checked against the index; a second version, a later
// bet, is refused.
func firstVersion(index []byte, versions map[int][]byte) ([][]string, error) {
	var rows [][]string
	if index != nil {
		var err error
		if rows, err = spec.ReadInventoryIndex(index); err != nil {
			return nil, fmt.Errorf("spec/inventory.tsv: %w", err)
		}
	}
	switch {
	case len(rows) == 0:
		return nil, fmt.Errorf("there is no confirmed inventory: spec/inventory.tsv has no version")
	case len(rows) > 1:
		return nil, fmt.Errorf("the reopening of the classes of a later bet is not built in M2e: the index names %d versions", len(rows))
	}
	data := versions[1]
	if data == nil {
		return nil, fmt.Errorf("the records have no spec/inventory-1.tsv")
	} else if got := digest(data); got != rows[0][1] {
		return nil, fmt.Errorf("spec/inventory-1.tsv has the sha256 %s; the index has %s", got, rows[0][1])
	}
	return spec.ReadInventory(data)
}

// reopeningFacts gives the facts that the full check exempts at the head of the task: those of the selection, less the facts whose row names the
// task or a task that has merged (loop.md, The full spec check).
func reopeningFacts(version [][]string, intake []spec.InventoryFact, rows [][]string, merged []string, task string) ([]string, error) {
	sel, err := reclassed(version, intake)
	if err != nil {
		return nil, err
	}
	owner := map[string]string{}
	for _, r := range rows {
		owner[r[reoFact]] = r[reoTask]
	}
	var out []string
	for _, r := range sel {
		if t, ok := owner[r[0]]; !ok || (t != task && !slices.Contains(merged, t)) {
			out = append(out, r[0])
		}
	}
	return out, nil
}

// specReopening sets Reopening of the input of the full check of a head of the task, from the latest version, Intake, spec/reopenings.tsv and the
// merged tasks that the input holds. Before the first bet it is empty; a file the reader refuses is an error.
func (s *Sessions) specReopening(ctx context.Context, in *spec.FullInput, task string) error {
	if len(in.Index) == 0 {
		return nil
	}
	fail := func(err error) error { return &SourceError{"Reopening", "records:" + reopeningsPath, err} }
	var rows [][]string
	if data, err := s.Store.ReadFile(ctx, reopeningsPath); err != nil {
		return fail(err)
	} else if data != nil {
		if rows, err = ReadReopenings(data); err != nil {
			return fail(fmt.Errorf("%s: %w", reopeningsPath, err))
		}
	}
	latest := 0
	for _, r := range in.Index {
		n, _ := strconv.Atoi(r[0])
		latest = max(latest, n)
	}
	version, err := spec.ReadInventory(in.Versions[latest]) // a version that cannot be read leaves the rules that need it not run
	if err != nil {
		return nil
	}
	var intake []spec.InventoryFact
	for _, f := range factRows(in.Spans, in.Questions, in.Answers) {
		intake = append(intake, spec.InventoryFact{ID: f[0], Class: f[1], Text: f[2]})
	}
	if in.Reopening, err = reopeningFacts(version, intake, rows, in.Merged, task); err != nil {
		return fail(err)
	}
	return nil
}

// taskIDs draws n task IDs, none of a row of sessions.tsv, tasks.tsv or spec/reopenings.tsv, and no two alike.
func (s *Sessions) taskIDs(ctx context.Context, n int) ([]string, error) {
	rows, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if register, err := s.Store.ReadFile(ctx, tasksPath); err != nil {
		return nil, err
	} else if register != nil {
		if tasks, err = ReadRegister(register); err != nil {
			return nil, fmt.Errorf("%s: %w", tasksPath, err)
		}
	}
	var reserved [][]string
	if data, err := s.Store.ReadFile(ctx, reopeningsPath); err != nil {
		return nil, err
	} else if data != nil {
		if reserved, err = ReadReopenings(data); err != nil {
			return nil, fmt.Errorf("%s: %w", reopeningsPath, err)
		}
	}
	random := s.random
	if random == nil {
		random = rand.Read
	}
	var ids []string
	for len(ids) < n {
		id, err := drawTaskID(random, func(id string) bool {
			return slices.Contains(ids, id) || slices.ContainsFunc(rows, func(r []string) bool { return r[sessionTask] == id }) ||
				slices.ContainsFunc(tasks, func(t Task) bool { return t.ID == id }) || slices.ContainsFunc(reserved, func(r []string) bool { return r[reoTask] == id })
		})
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// numberedRecordAt gives the one numbered facts record, F-NNNN, of the tree of the commit: the record the rendered commit of the bet wrote.
func (r *state) numberedRecordAt(commit string) (string, error) {
	files, err := filesAt(r.target, commit, "docs/facts", isNumberedFacts)
	if err != nil {
		return "", err
	}
	recs, err := psbRecordsOf(files)
	if err != nil {
		return "", fmt.Errorf("%s: %w", commit, err)
	} else if len(recs) == 0 {
		return "", fmt.Errorf("the commit %s holds no numbered facts record", commit)
	}
	return recs[0], nil
}

// reopenRecord gives the numbered facts record of the names of the files of docs/facts/ at the base, which the token of the row must name.
func reopenRecord(names []string, token string) (string, error) {
	files := map[string][]byte{}
	for _, n := range names {
		if isNumberedFacts(n) {
			files[n] = nil
		}
	}
	recs, err := psbRecordsOf(files)
	if err != nil {
		return "", err
	} else if len(recs) == 0 {
		return "", fmt.Errorf("the target has no numbered facts record at the base: the token %s cites none", token)
	}
	if strings.HasPrefix(token, "F-") && !strings.HasPrefix(token, recs[0]+"#") {
		return "", fmt.Errorf("the token %s names a record other than the numbered facts record %s of the base", token, recs[0])
	}
	return recs[0], nil
}

// reopenIssue fills the block reopen-issue for the row: the title (its first line, with no full stop) and the body.
func reopenIssue(row []string, ids []string, specTask string) (title, body string) {
	what, line, sections := "requirement", "one line that starts with "+row[reoRequirement]+": and cites "+row[reoToken], "docs/spec/ holds a section whose heading holds "+row[reoRequirement]+"; "
	if row[reoClass] == "measure" {
		what, line, sections = "success criterion", "one line that starts with Success criterion: and cites "+row[reoToken], ""
	}
	body = strings.NewReplacer("{task}", row[reoTask], "{what}", what, "{fact}", row[reoFact], "{class}", row[reoClass], "{ids}", strings.Join(ids, " "),
		"{spec}", specTask, "{token}", row[reoToken], "{line}", line, "{sections}", sections).Replace(reopenIssueBlock)
	first, _, _ := strings.Cut(body, "\n")
	return strings.TrimSuffix(first, "."), body
}

// openReopenings opens the specification tasks that the classes of the first bet reopen, after the event merged of the specification task (loop.md,
// The specification task that a changed class reopens): nothing before the event; then the rows of spec/reopenings.tsv are allocated once, in one
// records commit before any call of the forge, and a task is opened for each row that the register does not hold. It gives the tasks in the order
// of the rows.
func (r *state) openReopenings(ctx context.Context, specTask string) ([]Task, error) {
	s := r.sessions()
	register, err := r.readRegisterRows(ctx)
	if err != nil {
		return nil, err
	}
	si := slices.IndexFunc(register, func(t Task) bool { return t.ID == specTask && t.Kind == "specification" })
	if si < 0 {
		return nil, fmt.Errorf("the register has no specification task %s", specTask)
	}
	events, err := s.readTable(ctx, eventsPath(specTask), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	merge := ""
	for _, e := range events {
		if e[1] == "merged" {
			merge = e[5]
			break
		}
	}
	if merge == "" {
		return nil, nil
	}
	index, versions, err := r.inventoryOf(ctx)
	if err != nil {
		return nil, err
	}
	version, err := firstVersion(index, versions)
	if err != nil {
		return nil, err
	}
	intake, _, err := r.factsOf(ctx)
	if err != nil {
		return nil, err
	}
	sel, err := reclassed(version, intake)
	if err != nil || len(sel) == 0 {
		return nil, err
	}
	if _, err := s.Store.FetchBranch(ctx, s.Branch); err != nil {
		return nil, err
	}
	record, err := r.numberedRecordAt(merge)
	if err != nil {
		return nil, err
	}
	rows, err := r.reopeningRows(ctx, s, register, sel, record)
	if err != nil {
		return nil, err
	}
	var out []Task
	for _, row := range rows {
		reqs := []string{row[reoRequirement]}
		if row[reoRequirement] == "" {
			reqs = register[si].Requirements // a measure has no ID: it names the requirements of the specification
		}
		cur, err := r.readRegisterRows(ctx)
		if err != nil {
			return nil, err
		}
		if j := slices.IndexFunc(cur, func(t Task) bool { return t.ID == row[reoTask] }); j >= 0 {
			t := cur[j]
			if t.Kind != "specification" || !slices.Equal(t.Requirements, reqs) || !slices.Equal(t.Predecessors, []string{specTask}) {
				return nil, fmt.Errorf("the task %s of %s is registered as %s for %v after %v; the reservation is a specification task for %v after [%s]",
					t.ID, reopeningsPath, t.Kind, t.Requirements, t.Predecessors, reqs, specTask)
			}
			out = append(out, t)
			continue
		}
		title, body := reopenIssue(row, reqs, specTask)
		t, err := r.openTask(ctx, openSpec{id: row[reoTask], kind: "specification", ids: reqs, predecessors: []string{specTask},
			title:      func(string) string { return title },
			body:       func(string, []string, string, string) string { return body },
			afterIssue: func() error { return betStopAt("issue") }})
		if err != nil {
			return nil, err
		}
		if err := betStopAt("task"); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

// reopeningRows gives the rows of spec/reopenings.tsv: those of the record, bound to the selection, or, when the record is not there, the
// allocation of the selection, committed before it is returned.
func (r *state) reopeningRows(ctx context.Context, s *Sessions, register []Task, sel [][]string, record string) ([][]string, error) {
	data, err := r.store.ReadFile(ctx, reopeningsPath)
	if err != nil {
		return nil, err
	}
	if data != nil {
		rows, err := ReadReopenings(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", reopeningsPath, err)
		}
		return rows, bindReopenings(rows, sel, record)
	}
	draft, err := r.readRecord(ctx, "spec/draft.tsv", records.ReadDraft)
	if err != nil {
		return nil, err
	}
	var known []string
	for _, t := range register {
		known = append(known, t.Requirements...)
	}
	ids, err := s.taskIDs(ctx, len(sel))
	if err != nil {
		return nil, err
	}
	rows := allocate(sel, 1, record, nextNumber(finalIDs(draft), known), ids)
	if data, err = table(ReopeningsSchema, rows); err != nil {
		return nil, err
	}
	if _, err := ReadReopenings(data); err != nil {
		return nil, err
	}
	if err := r.write(ctx, map[string][]byte{reopeningsPath: data}, fmt.Sprintf("layup run: the reopening of the %d facts that the bet re-classed", len(rows))); err != nil {
		return nil, err
	}
	return rows, betStopAt("reopenings")
}
