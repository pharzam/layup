package run

// This file holds the step answers of Intake (docs/spec/intake.md, The
// answers; The raw fact of the answers; The steps of Intake; Input states;
// row 47a of the plan, task T-hxb9, #223): the comments of the Intake issue,
// copied, give the answers of the batch, and the intent form gives budget.tsv.

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// answerCopy is a row of copies.tsv and the body of its file.
type answerCopy struct {
	Row  []string
	Body string
}

// answerLine is the rule of an answer line (decided here, task T-hxb9): after
// a CR at its end is trimmed, the line starts with Q- and the question's
// number, then ": " exactly; the rest is the answer, byte for byte.
var answerLine = regexp.MustCompile(`^(Q-[0-9]{3,}): (.*)$`)

// parseAnswer gives the question ID and the answer of a line, or false: an
// empty answer, and one of exactly —, which the writer would write as the
// empty value, answer nothing.
func parseAnswer(line string) (string, string, bool) {
	m := answerLine.FindStringSubmatch(strings.TrimSuffix(line, "\r"))
	if m == nil || m[2] == "" || m[2] == "—" {
		return "", "", false
	}
	return m[1], m[2], true
}

// collectAnswers gives the rows of intake/answers.tsv of the questions given
// (the batch's, or the follow-up's, as records.ReadQuestions gives them), in
// the order of the questions, and the blocks that still wait. A copy counts when it is of the
// issue, made after the batch's asked (strictly, by its created), and a human
// decision of an approver in the role of the question's block (Decision: a
// first copy only, so an edit answers nothing; decided here, task T-hxb9, with
// run.md, A human decision). The last line of an ID, in the order of the
// copies, is its answer. A block that asks is in hand when a counted copy by
// its role holds an answer line of its questions.
func collectAnswers(questions, approvers [][]string, copies []answerCopy, issue string) ([][]string, []string) {
	if len(questions) == 0 {
		return nil, nil
	}
	asked := questions[0][9]
	block := map[string]string{} // the block of each question
	var blocks []string
	for _, q := range questions {
		block[q[0]] = q[3]
		if !slices.Contains(blocks, q[3]) {
			blocks = append(blocks, q[3])
		}
	}
	got := map[string][]string{}
	inHand := map[string]bool{}
	for _, c := range copies {
		if c.Row[2] != issue || c.Row[6] <= asked {
			continue
		}
		for _, line := range strings.Split(c.Body, "\n") {
			id, answer, ok := parseAnswer(line)
			b, asks := block[id]
			if !ok || !asks || !Decision(c.Row, approvers, b) {
				continue
			}
			got[id] = []string{id, answer, c.Row[0], c.Row[3], c.Row[6]}
			inHand[b] = true
		}
	}
	var waiting []string
	for _, b := range blocks {
		if !inHand[b] {
			waiting = append(waiting, b)
		}
	}
	var answers [][]string
	for _, q := range questions {
		if a, ok := got[q[0]]; ok {
			answers = append(answers, a)
		}
	}
	return answers, waiting
}

// pointLine is the rule of a line that names the approver of a planned point
// in a comment of the batch's owner (docs/spec/loop.md, The approvers and the
// success criteria, Decided here, task T-8rm6): the point's name, then ": "
// exactly, then the login.
var pointLine = regexp.MustCompile(`^([a-z][a-z-]*): ([A-Za-z0-9][-A-Za-z0-9]*)$`)

// pointBindings gives the login that each planned point's last counted line
// names, and for each login a counted line names the earliest comment that
// named it (its ID and its created — the source and the since of the row the
// step writes). A copy counts by the rule of an answer: of the issue, made
// after the batch's asked, and a decision of the idea owner (first copy, no
// App). A line whose name is none of the set's, and the rest of a counted
// comment, are input only and name nothing.
func pointBindings(copies []answerCopy, approvers [][]string, asked, issue string) (bound map[string]string, named map[string][2]string) {
	bound = map[string]string{}
	named = map[string][2]string{}
	for _, c := range copies {
		if c.Row[2] != issue || c.Row[6] <= asked || !Decision(c.Row, approvers, "idea-owner") {
			continue
		}
		for _, line := range strings.Split(c.Body, "\n") {
			m := pointLine.FindStringSubmatch(strings.TrimSuffix(line, "\r"))
			if m == nil || !slices.Contains(approvalPoints, m[1]) {
				continue
			}
			bound[m[1]] = m[2]
			if _, ok := named[m[2]]; !ok {
				named[m[2]] = [2]string{c.Row[0], c.Row[6]}
			}
		}
	}
	return bound, named
}

// approverRows gives the rows the step appends to approvers.tsv for each login
// a counted line names that holds no row of any role: one row each, in the
// order of the points, its ID from the forge, its since the naming comment's
// time and its source the earliest comment that named it. A login the forge
// does not know is recorded input and gives no row; a forge error is the
// step's (docs/spec/loop.md, The approvers and the success criteria; the
// Operator's and the idea owner's logins get none, architecture.md §5 and the
// decision of #365).
func approverRows(ctx context.Context, idOf func(context.Context, string) (int64, error), bound map[string]string, named map[string][2]string, approvers [][]string) ([][]string, error) {
	rowed := map[string]bool{}
	for _, a := range approvers {
		rowed[a[0]] = true // the account's ID: the login's spelling varies, the account does not
	}
	seen := map[string]bool{}
	var rows [][]string
	for _, p := range approvalPoints {
		login, ok := bound[p]
		if !ok || seen[login] {
			continue
		}
		seen[login] = true
		id, err := idOf(ctx, login)
		if err != nil {
			var fe *forge.Error
			if errors.As(err, &fe) && fe.Status == 404 {
				continue
			}
			return nil, err
		}
		key := strconv.FormatInt(id, 10)
		if rowed[key] {
			continue
		}
		rowed[key] = true
		rows = append(rows, []string{key, "approver", login, named[login][1], named[login][0]})
	}
	return rows, nil
}

// approverFor gives the approvers.tsv row whose account may answer the
// approval point (the answer #366's caller computes): the row of the login the
// point's last counted line names, whatever its role — the naming comment is
// enough for an account that already has a row, and a login is one account
// whatever its spelling's case — and the idea owner's row when the point names
// no account that holds one. nil when approvers holds no idea owner.
func approverFor(point string, bound map[string]string, approvers [][]string) []string {
	if login, ok := bound[point]; ok {
		for _, a := range approvers {
			if strings.EqualFold(a[2], login) {
				return a
			}
		}
	}
	for _, a := range approvers {
		if a[1] == "idea-owner" {
			return a
		}
	}
	return nil
}

// pointsServed gives the planned points the row serves, in the order of the
// set: each whose binding resolves to it.
func pointsServed(row []string, bound map[string]string, approvers [][]string) []string {
	var out []string
	for _, p := range approvalPoints {
		if a := approverFor(p, bound, approvers); a != nil && a[0] == row[0] {
			out = append(out, p)
		}
	}
	return out
}

// criteriaOf gives the criteria of the intent answer's answer: split on ";",
// each segment trimmed, the empty dropped (docs/spec/loop.md, Decided here).
func criteriaOf(answer string) []string {
	var out []string
	for _, s := range strings.Split(answer, ";") {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// criteriaAnswer gives the ID of the success-criteria question of the intent
// form and the criteria of its answer, none when it has none.
func criteriaAnswer(questions, answers [][]string) (string, []string) {
	var id string
	for _, q := range questions {
		if q[1] == "intent" && q[2] == "success-criteria" {
			id = q[0]
		}
	}
	if id == "" {
		return "", nil
	}
	for _, a := range answers {
		if a[0] == id {
			return id, criteriaOf(a[1])
		}
	}
	return id, nil
}

// decimal is the form of the type decimal of internal/tsv.
var decimal = regexp.MustCompile(`^(0|[1-9][0-9]*)\.[0-9]+$`)

// budgetOf gives the rows of budget.tsv when the three fields of the intent
// form parse: each answer, trimmed of the spaces around it (decided here, task
// T-hxb9), a decimal, and B below U as numbers. Else none and false.
func budgetOf(questions, answers [][]string) ([][]string, bool) {
	answerOf := map[string]string{}
	for _, a := range answers {
		answerOf[a[0]] = a[1]
	}
	value := map[string][2]string{} // the value and the question of each field
	for _, q := range questions {
		if q[1] != "intent" || !slices.Contains(budgetFields, q[2]) {
			continue
		}
		a, ok := answerOf[q[0]]
		if a = strings.TrimSpace(a); !ok || !decimal.MatchString(a) {
			return nil, false
		}
		value[q[2]] = [2]string{a, q[0]}
	}
	var rows [][]string
	for _, f := range budgetFields {
		v, ok := value[f]
		if !ok {
			return nil, false
		}
		rows = append(rows, []string{f, v[0], v[1]})
	}
	b, _ := new(big.Rat).SetString(value["band.B"][0])
	u, _ := new(big.Rat).SetString(value["band.U"][0])
	if b.Cmp(u) >= 0 {
		return nil, false
	}
	return rows, true
}

// answersStep is the step answers of Intake: "done before" when
// intake/answers.tsv exists; else the comments of the Intake issue are copied
// (but the App's own, run.md, Copy before read), the answers collected, and
// while a block that asks waits, the step is done "waiting", the last row of
// the restart; else one records commit of intake/answers.tsv, questions.tsv
// with the answered states, and budget.tsv when the intent form parses.
func (r *state) answersStep(ctx context.Context) Step {
	const name = "answers"
	if data, err := r.store.ReadFile(ctx, "intake/answers.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return failf(name, "start.tsv has no number of the Intake issue")
	}
	questions, err := r.readRecord(ctx, "questions.tsv", records.ReadQuestions)
	if err != nil {
		return fail(name, err)
	}
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return fail(name, err)
	}
	for _, c := range comments {
		if c.App == r.cfg.Register.AppSlug {
			continue
		}
		if err := r.copy(ctx, c, issue); err != nil {
			return fail(name, err)
		}
	}
	var copies []answerCopy
	for _, row := range r.copies {
		if row[2] != strconv.Itoa(issue) {
			continue
		}
		body, err := r.store.ReadFile(ctx, row[9])
		if err != nil {
			return fail(name, err)
		}
		copies = append(copies, answerCopy{Row: row, Body: string(body)})
	}
	answers, waiting := collectAnswers(questions, r.approvers, copies, strconv.Itoa(issue))
	if len(waiting) > 0 {
		return done(name, "waiting: no answer yet in the block "+strings.Join(waiting, " and the block "))
	}
	var newApprovers [][]string
	if len(questions) > 0 {
		bound, named := pointBindings(copies, r.approvers, questions[0][9], strconv.Itoa(issue))
		if newApprovers, err = approverRows(ctx, r.cfg.Forge.UserID, bound, named, r.approvers); err != nil {
			return fail(name, err)
		}
	}
	when := map[string]string{}
	for _, a := range answers {
		when[a[0]] = a[4]
	}
	updated := make([][]string, len(questions))
	for i, q := range questions {
		updated[i] = slices.Clone(q)
		if t, ok := when[q[0]]; ok {
			updated[i][11], updated[i][12] = t, "answered"
		}
	}
	files := map[string][]byte{}
	for _, f := range []struct {
		path  string
		write func() ([]byte, error)
	}{
		{"intake/answers.tsv", func() ([]byte, error) { return table(records.AnswersSchema, answers) }},
		{"questions.tsv", func() ([]byte, error) { return table(records.QuestionsSchema, updated) }},
	} {
		if files[f.path], err = f.write(); err != nil {
			return fail(name, err)
		}
	}
	if _, err := records.ReadAnswers(files["intake/answers.tsv"]); err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadQuestions(files["questions.tsv"]); err != nil {
		return fail(name, err)
	}
	written := "budget.tsv not written (a field of the intent form does not parse)"
	if rows, ok := budgetOf(questions, answers); ok {
		if files["budget.tsv"], err = table(records.BudgetSchema, rows); err != nil {
			return fail(name, err)
		}
		written = "budget.tsv written"
	}
	if len(newApprovers) > 0 {
		held := append(slices.Clone(r.approvers), newApprovers...)
		if files["approvers.tsv"], err = table(records.ApproversSchema, held); err != nil {
			return fail(name, err)
		}
		if _, err := records.ReadApprovers(files["approvers.tsv"]); err != nil {
			return fail(name, err)
		}
	}
	if err := r.write(ctx, files, "layup run: the answers of the batch"); err != nil {
		return fail(name, err)
	}
	r.approvers = append(r.approvers, newApprovers...)
	named := ""
	if len(newApprovers) > 0 {
		named = fmt.Sprintf("; approvers.tsv: %d approver row(s)", len(newApprovers))
	}
	return done(name, fmt.Sprintf("%d answers, in intake/answers.tsv; %s%s", len(answers), written, named))
}
