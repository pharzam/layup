package run

// This file holds the step bet (docs/spec/loop.md, The step bet; task T-2vva, #366, row 66j of the plan of M2e): the milestone, the
// brief and its post, the answer of the approver, and the order of the acts that follow it, the row of bets.tsv last.

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
)

// betInput is what the answer to a brief is parsed against: the requirement IDs, the phase tags of the target, the fact IDs and the
// survivors of each decision, all of the kept basis of the brief.
type betInput struct {
	Required, Tags, Facts []string
	Survivors             map[string][]string
}

// separatorCell is a cell of the separator row of a table: hyphens, with a colon at either end or both.
var separatorCell = regexp.MustCompile(`^:?-+:?$`)

// phaseTagsOf gives the first column of the table under the heading "## 9." of a PRD: the phase tags of the target (the section
// "Rollout & phases" of the baseline's PRD form; the specification instruction asks for it). The table is the first one of the section and
// has the header Phase | Name | Ships, a separator row and rows of three cells with a tag; any other form gives no tag.
func phaseTagsOf(prd []byte) []string {
	var rows [][]string
	in, started := false, false
scan:
	for _, raw := range strings.Split(string(prd), "\n") {
		l := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(l, "## "):
			if in {
				break scan
			}
			in = strings.HasPrefix(l, "## 9.")
		case !in:
		case strings.HasPrefix(l, "|") && strings.HasSuffix(l, "|") && len(l) > 1:
			started = true
			cells := strings.Split(l[1:len(l)-1], "|")
			for i := range cells {
				cells[i] = strings.TrimSpace(cells[i])
			}
			rows = append(rows, cells)
		case started:
			break scan
		}
	}
	if len(rows) < 3 || !slices.Equal(rows[0], []string{"Phase", "Name", "Ships"}) {
		return nil
	}
	if len(rows[1]) != 3 {
		return nil
	}
	for _, c := range rows[1] {
		if !separatorCell.MatchString(c) {
			return nil
		}
	}
	var tags []string
	for _, r := range rows[2:] {
		if len(r) != 3 || r[0] == "" {
			return nil
		}
		tags = append(tags, r[0])
	}
	return tags
}

// optionsColumn gives the column options of bets.tsv: the chosen options as <decision>=<option>, sorted by decision, none for none.
func optionsColumn(options map[string]string) []string {
	var out []string
	for _, d := range slices.Sorted(maps.Keys(options)) {
		out = append(out, d+"="+options[d])
	}
	return out
}

// betAnswer gives the answer to the brief: the earliest copy, by numeric comment ID, that is of the issue, the first copy of its comment,
// no App's, by the approver's account, created strictly after the brief's copy, and that ParseAnswer takes without a finding. It gives
// the parsed answer and the row of the copy.
func betAnswer(copies []answerCopy, brief, approver []string, in betInput) (spec.Answer, []string, bool) {
	sorted := slices.Clone(copies)
	id := func(c answerCopy) int { n, _ := strconv.Atoi(c.Row[copyComment]); return n }
	sort.SliceStable(sorted, func(i, j int) bool { return id(sorted[i]) < id(sorted[j]) })
	after, err := time.Parse(timeForm, brief[copyCreated])
	if err != nil {
		return spec.Answer{}, nil, false
	}
	for _, c := range sorted {
		made, err := time.Parse(timeForm, c.Row[copyCreated])
		if err != nil || c.Row[copyIssue] != brief[copyIssue] || c.Row[copySeen] != "1" || c.Row[copyApp] != "" || c.Row[copyAuthorID] != approver[0] || !made.After(after) {
			continue
		}
		if a, findings := spec.ParseAnswer([]byte(c.Body), in.Required, in.Tags, in.Facts, in.Survivors); len(findings) == 0 {
			return a, c.Row, true
		}
	}
	return spec.Answer{}, nil, false
}

// betStopAt is the narrow seam of the tests of the atomicity of the step: it is called after each act of the step that a stop can
// follow ("version", "pushevent", "pushed", "bound", "carry", "issue", "task"), and an error it gives stops the step there.
var betStopAt = func(stage string) error { return nil }

// renderedPushed reports whether the attempt has an event push or bound of code, with no session: the rendered commit of the first bet is
// pushed or its push is pending, and the step bet, not the driver of the Shape task, finishes it.
func renderedPushed(events [][]string, attempt int) bool {
	return slices.ContainsFunc(events, func(e []string) bool {
		a, _ := strconv.Atoi(e[2])
		return (e[1] == "push" || e[1] == "bound") && e[3] == "" && a == attempt
	})
}

// betStep is the step bet. It reads the records and the forge at each run and holds nothing: the stages of loop.md, The step bet, each
// skipped when its record exists, and bets.tsv written last.
func (r *state) betStep(ctx context.Context) Step {
	const name = "bet"
	res, err := r.bet(ctx)
	switch {
	case errors.Is(err, ErrBriefLinesUnset), errors.Is(err, ErrBriefOverLimit):
		return done(name, "waiting for the Operator: "+err.Error())
	case err != nil:
		return fail(name, err)
	}
	return done(name, res)
}

func (r *state) bet(ctx context.Context) (string, error) {
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return "", errors.New("start.tsv has no number of the Intake issue")
	}
	s := r.sessions()
	// 0. Read: the Shape tasks at their verified heads, the milestone, the bets so far.
	register, err := r.readRegisterRows(ctx)
	if err != nil {
		return "", err
	}
	var specTask, archTask Task
	for _, t := range register {
		if t.Kind == "specification" && specTask.ID == "" {
			specTask = t
		} else if t.Kind == "architecture" && archTask.ID == "" {
			archTask = t
		}
	}
	if specTask.ID == "" || archTask.ID == "" {
		return "", errors.New("the register holds no Shape task: the step shape comes first")
	}
	st, err := s.verifiedStep(ctx, specTask.ID)
	if err != nil {
		return "", err
	}
	prds, err := filesAt(r.target, st.Head, "docs/prd", isPRDName)
	if err != nil {
		return "", err
	}
	if len(prds) != 1 {
		return "", fmt.Errorf("the verified head %s holds %d files docs/prd/PRD-*.md; want one", st.Head, len(prds))
	}
	var prd []byte
	for _, b := range prds {
		prd = b
	}
	tags := phaseTagsOf(prd)
	if len(tags) == 0 {
		return "", fmt.Errorf("the PRD at the verified head %s has no phase tag in its section 9: the milestone is not known", st.Head)
	}
	milestone := tags[0]
	if err := checkBriefRequest(milestone, "first"); err != nil {
		return "", err
	}
	betsData, err := r.store.ReadFile(ctx, "bets.tsv")
	if err != nil {
		return "", err
	}
	var bets [][]string
	if betsData != nil {
		if bets, err = ReadBets(betsData); err != nil {
			return "", fmt.Errorf("bets.tsv: %w", err)
		}
	}
	for _, b := range bets {
		if b[betMilestone] == milestone && b[betVerdict] == "bet" {
			return "done before", nil
		} else if b[betMilestone] == milestone {
			return "waiting for the Operator: no bet", nil
		}
	}
	if len(bets) > 0 {
		return "", errors.New("a later bet's rendered task is row 66t's (moved to M2g): M2e builds the first bet only")
	}
	// 1 and 2. The brief and its post.
	brief, _, err := r.betBrief(ctx, milestone, "first")
	if err != nil {
		return "", err
	}
	if err := CheckCommentBody(brief); err != nil {
		return "", err
	}
	posted, err := s.postOnIssue(ctx, issue, brief)
	if err != nil {
		return "", err
	}
	// 3. The answer: every comment that is not the App's is copied first.
	if r.copies, err = r.readRecord(ctx, "copies.tsv", records.ReadCopies); err != nil {
		return "", err
	}
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return "", err
	}
	for _, c := range comments {
		if c.App == r.cfg.Register.AppSlug {
			continue
		}
		if err := r.copy(ctx, c, issue); err != nil {
			return "", err
		}
	}
	var copies []answerCopy
	for _, row := range r.copies {
		if row[copyIssue] != strconv.Itoa(issue) {
			continue
		}
		body, err := r.store.ReadFile(ctx, row[copyBody])
		if err != nil {
			return "", err
		}
		copies = append(copies, answerCopy{Row: row, Body: string(body)})
	}
	briefCopy, ok := copyNamed(r.copies, issue, posted.SHA)
	if !ok {
		return "", errors.New("the copy of the brief is not in copies.tsv")
	}
	questions, err := r.readRecord(ctx, "questions.tsv", records.ReadQuestions)
	if err != nil {
		return "", err
	}
	asked := ""
	if len(questions) > 0 {
		asked = questions[0][9]
	}
	basis, err := r.briefBasis(ctx, milestone)
	if err != nil {
		return "", err
	}
	in := betInput{Required: basis.Requirements, Tags: tags, Survivors: basis.Survivors}
	for _, f := range basis.Facts {
		in.Facts = append(in.Facts, f.Number)
	}
	var answer spec.Answer
	var copyRow, approver []string
	if applied, err := r.appliedAnswer(ctx, milestone); err != nil {
		return "", err
	} else if applied != "" {
		// The inventory names the answer of this bet: a later comment or a later binding of the approver does not replace it.
		if answer, copyRow, approver, err = r.restoreAnswer(applied, copies, in); err != nil {
			return "", err
		}
	} else {
		bound, _ := pointBindings(copies, r.approvers, asked, strconv.Itoa(issue))
		if approver, err = approverFor(ctx, r.cfg.Forge.UserID, "bet", bound, r.approvers); err != nil {
			return "", err
		} else if approver == nil {
			return "", errors.New("approvers.tsv holds no idea owner to answer the bet")
		}
		var ok bool
		if answer, copyRow, ok = betAnswer(copies, briefCopy, approver, in); !ok {
			return fmt.Sprintf("waiting: the answer of the bet on the brief %d", posted.ID), nil
		}
	}
	answer.Milestone = milestone
	write := func(bet []string, message string) error {
		data, err := table(BetsSchema, [][]string{bet})
		if err == nil {
			_, err = ReadBets(data)
		}
		if err != nil {
			return err
		}
		return r.write(ctx, map[string][]byte{"bets.tsv": data}, message)
	}
	row := betRow(milestone, posted, basis.Appetite, answer, copyRow, approver)
	if answer.Verdict == "no-bet" {
		if err := write(row, "layup run: the bet "+milestone+" (no bet)"); err != nil {
			return "", err
		}
		return "waiting for the Operator: no bet", nil
	}
	// 4. Validate before any write of the inventory, the rendering and the tasks.
	settled := settledDecisions(answer.Options, basis.Survivors)
	draft, err := r.readRecord(ctx, "spec/draft.tsv", records.ReadDraft)
	if err != nil {
		return "", err
	}
	requirements := map[string][]string{}
	tables := map[string][][]string{}
	for _, d := range settled {
		data, err := r.store.ReadFile(ctx, "spec/options/"+d+".tsv")
		if err != nil {
			return "", err
		} else if data == nil {
			return "", fmt.Errorf("the records have no spec/options/%s.tsv", d)
		}
		if tables[d], err = spec.ReadOptions(data); err != nil {
			return "", fmt.Errorf("spec/options/%s.tsv: %w", d, err)
		}
		if requirements[d], err = archRequirements(d, answer.Options[d], tables[d], draft); err != nil {
			return "", err
		}
	}
	// 5. The version of the inventory.
	spans, err := r.readRecord(ctx, "spec/spans.tsv", records.ReadSpans)
	if err != nil {
		return "", err
	}
	answers, err := r.readRecord(ctx, "intake/answers.tsv", records.ReadAnswers)
	if err != nil {
		return "", err
	}
	var facts []spec.InventoryFact
	texts := map[string]string{}
	for _, f := range factRows(spans, questions, answers) {
		facts, texts[f[0]] = append(facts, spec.InventoryFact{ID: f[0], Class: f[1], Text: f[2]}), f[2]
	}
	index, err := r.store.ReadFile(ctx, "spec/inventory.tsv")
	if err != nil {
		return "", err
	}
	versions := map[int][]byte{}
	if index != nil {
		rows, err := spec.ReadInventoryIndex(index)
		if err != nil {
			return "", fmt.Errorf("spec/inventory.tsv: %w", err)
		}
		for n := range rows {
			if versions[n+1], err = r.store.ReadFile(ctx, fmt.Sprintf("spec/inventory-%d.tsv", n+1)); err != nil {
				return "", err
			}
		}
	}
	copyID, _ := strconv.Atoi(copyRow[copyComment])
	files, version, err := nextVersion(index, versions, facts, answer, copyID)
	if err != nil {
		return "", err
	}
	if files != nil {
		if err := r.write(ctx, files, "layup run: the inventory of the bet "+milestone+" (copy "+copyRow[copyComment]+")"); err != nil {
			return "", err
		}
		if err := betStopAt("version"); err != nil {
			return "", err
		}
	}
	// 6. The rendered commit and its carry.
	rendered, err := r.renderedStage(ctx, s, specTask.ID, st, milestone, answer, copyRow, version, texts)
	if err != nil {
		return "", err
	}
	// 7. The architecture tasks.
	if err := r.archStage(ctx, settled, answer, tables, requirements, []string{specTask.ID, archTask.ID}); err != nil {
		return "", err
	}
	// 8. bets.tsv last.
	if err := write(row, "layup run: the bet "+milestone); err != nil {
		return "", err
	}
	return fmt.Sprintf("the bet %s: the rendered commit %s, %d architecture tasks settled", milestone, rendered[:12], len(settled)), nil
}

// readRegisterRows gives the register, none for a records branch with no tasks.tsv.
func (r *state) readRegisterRows(ctx context.Context) ([]Task, error) {
	data, err := r.store.ReadFile(ctx, tasksPath)
	if err != nil || data == nil {
		return nil, err
	}
	tasks, err := ReadRegister(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tasksPath, err)
	}
	return tasks, nil
}

// renderedStage makes the rendered commit of the first bet on the verified head, pushes it as code, and carries layup/verify to it;
// each act is skipped when its event exists. It gives the commit.
func (r *state) renderedStage(ctx context.Context, s *Sessions, task string, st loopStep, milestone string, answer spec.Answer, copyRow []string, version [][]string, texts map[string]string) (string, error) {
	prds, err := filesAt(r.target, st.Head, "docs/prd", isPRDName)
	if err != nil {
		return "", err
	}
	existing, err := filesAt(r.target, st.Head, "docs/facts", isNumberedFacts)
	if err != nil {
		return "", err
	}
	entries, err := git.LsTree(r.target, st.Head, "docs/facts")
	if err != nil {
		return "", err
	}
	id := nextFactsID(factRecordIDs(factsNames(entries)))
	files, recs, err := renderedFiles(prds, existing, id, answer, version, texts)
	if err != nil {
		return "", err
	}
	copied, err := time.Parse(timeForm, copyRow[copyCopied])
	if err != nil {
		return "", err
	}
	commit, err := r.commitFiles(st.Head, files, fmt.Sprintf("layup run: the rendered commit of the bet %s (copy %s)", milestone, copyRow[copyComment]), r.who(copied))
	if err != nil {
		return "", err
	}
	branch := fmt.Sprintf("task/%s/%d", task, st.Attempt)
	in := CarryInput{Task: task, Verifier: st.Session, Harnesses: r.cfg.Harnesses, Models: r.cfg.Models, Rendered: &RenderedInput{Head: commit, Render: recs}}
	if ok, err := s.CheckRendered(ctx, in); err != nil {
		return "", err
	} else if !ok {
		return "", fmt.Errorf("%s is not a permitted rendered commit above the verified head %s", commit, st.Head)
	}
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		return "", err
	}
	has := func(kind string) bool {
		return slices.ContainsFunc(events, func(e []string) bool { return e[1] == kind && e[3] == "" && e[5] == commit && e[6] == branch })
	}
	event := func(kind, message string) error {
		data, err := s.appendEvents(ctx, task, st.Attempt, "", []string{kind + " " + commit + " " + branch})
		if err != nil {
			return err
		}
		return s.commit(ctx, map[string][]byte{eventsPath(task): data}, message)
	}
	if !has("bound") {
		if !has("push") {
			if err := event("push", "layup run: the push of the rendered commit to "+branch); err != nil {
				return "", err
			}
			if err := betStopAt("pushevent"); err != nil {
				return "", err
			}
		}
		// Only the verified head (the push is to come) and the rendered commit (the push was made) are the branch's tips that this step goes on from.
		switch tip, found, err := s.Store.BranchHead(ctx, branch); {
		case err != nil:
			return "", err
		case !found || (tip != st.Head && tip != commit):
			return "", fmt.Errorf("%w: %s holds %q, not the verified head %s or the rendered commit %s", errPushRefused, branch, tip, st.Head, commit)
		}
		if err := s.Store.PushHead(ctx, commit, branch); err != nil {
			return "", err
		}
		if err := betStopAt("pushed"); err != nil {
			return "", err
		}
		if err := event("bound", "layup run: the rendered commit bound to "+branch); err != nil {
			return "", err
		}
		if err := betStopAt("bound"); err != nil {
			return "", err
		}
	}
	if _, err := s.CarryVerify(ctx, in, ""); err != nil {
		return "", err
	}
	return commit, betStopAt("carry")
}

// commitFiles makes the commit of files on a work tree at parent, by who at the time of who, and gives its ID; the work tree is removed.
func (r *state) commitFiles(parent string, files map[string][]byte, message string, who git.Identity) (string, error) {
	tmp, err := tempDir()
	if err != nil {
		return "", err
	}
	defer removeAll(tmp)
	tree := tmp + "/tree"
	if err := git.WorktreeAdd(r.target, tree, parent); err != nil {
		return "", err
	}
	defer git.WorktreeRemove(r.target, tree)
	if err := writeFiles(tree, files); err != nil {
		return "", err
	}
	if err := git.Add(tree); err != nil {
		return "", err
	}
	if err := git.Commit(tree, message, who); err != nil {
		return "", err
	}
	return git.RevParse(tree, "HEAD")
}

// appliedAnswer gives the comment ID of the answer that the latest version of the inventory names, when that version is the one of this
// bet's milestone: the answer is applied, and a restart goes on with it. It gives "" when the bet has no version yet.
func (r *state) appliedAnswer(ctx context.Context, milestone string) (string, error) {
	data, err := r.store.ReadFile(ctx, "spec/inventory.tsv")
	if err != nil || data == nil {
		return "", err
	}
	index, err := spec.ReadInventoryIndex(data)
	if err != nil {
		return "", fmt.Errorf("spec/inventory.tsv: %w", err)
	}
	if len(index) == 0 || index[len(index)-1][2] != milestone {
		return "", nil
	}
	path := fmt.Sprintf("spec/inventory-%s.tsv", index[len(index)-1][0])
	version, err := r.store.ReadFile(ctx, path)
	if err != nil {
		return "", err
	} else if version == nil {
		return "", fmt.Errorf("the records have no %s", path)
	}
	rows, err := spec.ReadInventory(version)
	if err != nil || len(rows) == 0 {
		return "", fmt.Errorf("%s: %v", path, errors.Join(err, errors.New("no row")))
	}
	return rows[0][3], nil
}

// restoreAnswer gives the answer of the copy that the inventory names (its first copy), the copy's row and the row of approvers.tsv of the
// account that wrote it; the answer must still parse against the kept basis.
func (r *state) restoreAnswer(applied string, copies []answerCopy, in betInput) (spec.Answer, []string, []string, error) {
	for _, c := range copies {
		if c.Row[copyComment] != applied || c.Row[copySeen] != "1" {
			continue
		}
		a, fs := spec.ParseAnswer([]byte(c.Body), in.Required, in.Tags, in.Facts, in.Survivors)
		if len(fs) > 0 {
			return a, nil, nil, fmt.Errorf("the answer %s that the inventory names does not parse: %s", applied, findings(fs))
		}
		for _, ap := range r.approvers {
			if ap[0] == c.Row[copyAuthorID] {
				return a, c.Row, ap, nil
			}
		}
		return a, nil, nil, fmt.Errorf("the author %s of the answer %s is no row of approvers.tsv", c.Row[copyAuthorID], applied)
	}
	return spec.Answer{}, nil, nil, fmt.Errorf("the copies hold no first copy of the answer %s that the inventory names", applied)
}
