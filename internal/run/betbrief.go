package run

// This file holds the brief of the first bet (docs/spec/loop.md, The brief of the first bet; task T-h65e, #418, row 66u of the plan of
// M2e): its inputs, its session, its checks, the appetite after the session's telemetry, the promotion and the basis. Row 66j calls
// betBrief and posts the brief.

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/spec"
)

// ErrBriefLinesUnset and ErrBriefOverLimit are the two waits before any session; row 66j maps them to done waiting for the Operator.
var (
	ErrBriefLinesUnset = errors.New("brief.lines is not set")
	ErrBriefOverLimit  = errors.New("the parts of code are over brief.lines")
)

// BriefBasis is what a kept brief rests on, read from bets/<milestone>/basis.tsv and the brief.
type BriefBasis struct {
	Requirements []string            // the final IDs
	Survivors    map[string][]string // the survivors of each decision, in the order of the filter
	Dropped      map[string]string   // the reasons of each dropped option, by <decision>=<option>
	Facts        []spec.Fact         // the number or Q-NNN of each fact with its class
	Open         []string            // the Q-NNN of the questions still open
	Heads        map[string]string   // the verified head of each Shape task, by task ID
	Appetite     spec.Cap            // the numbers of the brief's appetite line
}

// headRef is a Shape task with its verified head and its pull request.
type headRef struct {
	Task, Kind, Head string
	PR               int
}

// briefParts are the four texts that code writes into the brief.
type briefParts struct{ requirements, solution, open, facts string }

// briefPrompt and briefRefused are the blocks bet-brief and bet-brief-refused of docs/spec/loop.md, byte for byte.
const briefPrompt = `You write three parts of the brief of the bet of milestone {milestone}, for LAYUP. Change no file in this repository and make no commit.

Code writes the other parts of the brief and puts yours between them. Read these files:
- The problem statement: {problem}
- The fact spans: {spans}
- The draft: {draft}
- The questions: {questions}
- The answers: {answers}
- The budget: {budget}
- The requirements and the success criteria, which code writes: {requirements}
- The solution, which code writes: {solution}
- The questions still open, which code writes: {open}
- The facts, which code writes: {facts}
- The links you may use, which code writes: {links}

Write these three files of Markdown in the directory {output}, and no other file there:
- problem.md: the problem of this milestone in plain prose for the idea owner, from the problem statement and the requirements.
- appetite.md: the money and the time you ask for, with the line "appetite: USD <n>, <n> hours" exactly once, each <n> a decimal such as 12.5, and the reasons in a few lines.
- rabbit-holes.md: the rabbit holes and the no-gos of this milestone in plain prose, one line each.

The whole brief is at most {lines} lines, and your three files together have {room} lines at most. Code checks the brief after you, refuses it and asks you again for each of these:
- A line that starts with "## " in any file: code writes the headings.
- A line that starts with "architecture change:" or "approve batch " in any file: code writes the first, and this bet has no batch.
- A link in any form but [text](target), or a target that is not a line of the links file.
- An option of a decision with more than one survivor, named anywhere in your files. Code lists the options in the solution part. You recommend none of them, not by name and not by a phrase such as "the first one": the idea owner chooses.
- No appetite line, or two. A money cap above U of the budget, hours above appetite.hours of the budget, or a money cap that with the spend of the milestone so far is above U.
- A brief of more than 65536 characters.

{refused}Then end.
`

const briefRefused = `The check refused the last brief: {finding}
Write the three files again.
`

// briefFiles are the files of the session's output directory.
var (
	briefFiles    = []string{"problem.md", "appetite.md", "rabbit-holes.md"}
	milestoneForm = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	safePath      = regexp.MustCompile(`^[A-Za-z0-9._/-]+$`)
	linkJoin      = strings.NewReplacer("](", "] (", "][", "] [")
)

func briefPath(milestone string) string   { return "bets/" + milestone + "/brief.md" }
func basisPath(milestone string) string   { return "bets/" + milestone + "/basis.tsv" }
func attemptsDir(milestone string) string { return "bets/" + milestone + "/attempts" }

// checkBriefRequest refuses a milestone that is not one path part and a kind that has no brief in M2e.
func checkBriefRequest(name, kind string) error {
	if !milestoneForm.MatchString(name) || name == "." || name == ".." {
		return fmt.Errorf("the milestone %q is not one path part", name)
	}
	switch kind {
	case "first":
		return nil
	case "later":
		return errors.New("the brief of a later bet is built by no row of M2e")
	}
	return fmt.Errorf("the kind %q is not first or later", kind)
}

// plain writes a text that code copies so that it is no link of the form [text](target) or [text][ref].
func plain(s string) string { return linkJoin.Replace(s) }

// requirementsPart is requirements.md: the words of the PRD at the verified head for exactly the requirement IDs of the register row.
func requirementsPart(ids []string, prd spec.PRDContent) (string, error) {
	var b strings.Builder
	b.WriteString("### The requirements\n")
	for _, id := range ids {
		statement, ok := prd.Statements[id]
		switch {
		case !ok:
			return "", fmt.Errorf("the requirement %s of the register row has no row in the PRD at the verified head", id)
		case strings.TrimSpace(statement) == "":
			return "", fmt.Errorf("the requirement %s has no statement in the PRD at the verified head", id)
		case prd.Criteria[id] == "":
			return "", fmt.Errorf("the requirement %s has no acceptance criterion in the PRD at the verified head", id)
		}
		fmt.Fprintf(&b, "- %s: %s Criterion: %s\n", id, plain(statement), plain(prd.Criteria[id]))
	}
	for _, id := range slices.Sorted(maps.Keys(prd.Statements)) {
		if !slices.Contains(ids, id) {
			return "", fmt.Errorf("the PRD at the verified head holds %s, which the register row does not", id)
		}
	}
	b.WriteString("### The success criteria\n")
	for _, m := range prd.Metrics {
		b.WriteString("- " + plain(m) + "\n")
	}
	if len(prd.Metrics) == 0 {
		b.WriteString("- none\n")
	}
	return b.String(), nil
}

func prURL(owner, name string, pr int) string {
	return "https://github.com/" + owner + "/" + name + "/pull/" + strconv.Itoa(pr)
}

// solutionPart is solution.md: the verified heads and, for each decision, each option as a survivor or dropped with its reasons. The
// presentation is code's alone and marks no option.
func solutionPart(owner, name string, heads []headRef, results map[string]spec.Result) string {
	var b strings.Builder
	b.WriteString("### The verified heads\n")
	for _, h := range heads {
		fmt.Fprintf(&b, "- %s (%s): [pull request #%d](%s), head %s\n", h.Task, h.Kind, h.PR, prURL(owner, name, h.PR), h.Head)
	}
	for _, d := range slices.Sorted(maps.Keys(results)) {
		fmt.Fprintf(&b, "### Decision %s\n", d)
		one, single := results[d].Chosen()
		for _, o := range results[d].Options {
			switch {
			case o.Dropped:
				fmt.Fprintf(&b, "- %s (dropped: %s)\n", o.Name, plain(strings.Join(o.Reasons, "; ")))
			case single && o.Name == one:
				fmt.Fprintf(&b, "- %s (the one survivor, decided at Shape)\n", o.Name)
			default:
				fmt.Fprintf(&b, "- %s (a survivor)\n", o.Name)
			}
		}
	}
	return b.String()
}

// openPart is open.md: each question still open, with its text.
func openPart(questions [][]string) string {
	var b strings.Builder
	b.WriteString("### The questions still open\n")
	text := map[string]string{}
	for _, q := range questions {
		text[q[0]] = q[8]
	}
	open := stillOpen(questions)
	for _, id := range open {
		fmt.Fprintf(&b, "- %s: %s\n", id, plain(text[id]))
	}
	if len(open) == 0 {
		b.WriteString("- none\n")
	}
	return b.String()
}

// answerFactClass is the class of the fact of an answer: that of the span fact that holds the first byte of the question's quote, the
// quote of the question it asks again for a follow-up (a Shape row too), and context for no quote or a quote outside every fact.
func answerFactClass(id string, questions, spans [][]string) string {
	byID := map[string][]string{}
	for _, q := range questions {
		byID[q[0]] = q
	}
	q := byID[id]
	for hops := 0; q != nil && q[1] == "follow-up" && q[6] == "" && hops < len(questions); hops++ {
		q = byID[q[2]]
	}
	if q == nil {
		return "context"
	}
	if fact := factOfQuote(q[6], q[7], spans); fact != "" {
		for _, r := range spans {
			if r[2] == fact {
				return r[3]
			}
		}
	}
	return "context"
}

// factRows gives the facts as number, class and text: the span facts in the order of the table, then one fact per answer.
func factRows(spans, questions, answers [][]string) [][3]string {
	var out [][3]string
	for _, r := range spans {
		if r[2] != "" {
			out = append(out, [3]string{r[2], r[3], r[5]})
		}
	}
	for _, a := range answers {
		out = append(out, [3]string{a[0], answerFactClass(a[0], questions, spans), a[1]})
	}
	return out
}

// factsPart is facts.md: each fact with its class and its text.
func factsPart(spans, questions, answers [][]string) string {
	var b strings.Builder
	for _, f := range factRows(spans, questions, answers) {
		fmt.Fprintf(&b, "- %s (%s): %s\n", f[0], f[1], plain(f[2]))
	}
	return b.String()
}

// linkSet is the set of the targets that a link of the brief may have: the pull request of each Shape task and the regular files that
// its head changed, by their blob URL at that head.
func linkSet(owner, name string, heads []headRef, blobs map[string][]string) map[string]bool {
	set := map[string]bool{}
	for _, h := range heads {
		set[prURL(owner, name, h.PR)] = true
		for _, p := range blobs[h.Task] {
			if safePath.MatchString(p) {
				set["https://github.com/"+owner+"/"+name+"/blob/"+h.Head+"/"+p] = true
			}
		}
	}
	return set
}

// assemble puts the three files of the session between the parts of code, in the fixed order of the brief.
func assemble(p briefParts, problem, appetite, holes string) string {
	return "## The problem\n" + p.requirements + problem + "## The appetite\n" + appetite + "## The solution\n" + p.solution +
		"## The rabbit holes and the no-gos\n" + p.open + "architecture change: no\n" + holes + "## The facts\n" + p.facts
}

// joinFindings is the text of the findings of a stage: rule: detail; rule: detail.
func joinFindings(fs []spec.Finding) string {
	parts := make([]string, len(fs))
	for i, f := range fs {
		parts[i] = f.Rule + ": " + f.Detail
	}
	return strings.Join(parts, "; ")
}

// briefKeeper judges the output of a session and a stored candidate.
type briefKeeper struct {
	parts briefParts
	multi map[string][]string // every option of each decision with more than one survivor
	lines *int
	links map[string]bool
}

// overLimitError is the wait of the parts of code that alone are over brief.lines.
type overLimitError struct{ n, limit int }

func (e overLimitError) Error() string {
	return fmt.Sprintf("the parts of code are %d lines, over brief.lines %d", e.n, e.limit)
}
func (e overLimitError) Is(target error) bool { return target == ErrBriefOverLimit }

// overLimit gives the wait before any session: brief.lines unset, or the parts of code alone above it.
func (k briefKeeper) overLimit() error {
	if k.lines == nil {
		return ErrBriefLinesUnset
	}
	if n := strings.Count(assemble(k.parts, "", "", ""), "\n"); n > *k.lines {
		return overLimitError{n, *k.lines}
	}
	return nil
}

// judge makes the checks that an assembled brief passes, from its bytes alone (the session's files are gone when a candidate is
// judged again): the checker of the brief, the option-word rule and the body of a comment. The first stage that finds gives its findings.
func (k briefKeeper) judge(assembled string) []spec.Finding {
	if fs := k.fixedParts(assembled); len(fs) > 0 {
		return fs
	}
	if _, fs := spec.CheckBrief([]byte(assembled), k.lines, k.links, nil); len(fs) > 0 {
		return fs
	}
	if fs := k.optionWords(assembled); len(fs) > 0 {
		return fs
	}
	var refusal HandoffRefusal
	if err := CheckCommentBody([]byte(assembled)); errors.As(err, &refusal) {
		return []spec.Finding{{Rule: refusal.Reason, Detail: refusal.Value}}
	}
	return nil
}

// fixedParts requires each part of code in its place, as the inputs stand now: a candidate of an earlier call, made on other inputs, is not
// promoted over a basis built from these.
func (k briefKeeper) fixedParts(assembled string) []spec.Finding {
	p := k.parts
	var fs []spec.Finding
	for _, c := range []struct {
		what string
		ok   bool
	}{
		{"requirements", strings.HasPrefix(assembled, "## The problem\n"+p.requirements)},
		{"solution", strings.Contains(assembled, "## The solution\n"+p.solution+"## The rabbit holes and the no-gos\n")},
		{"open questions", strings.Contains(assembled, "## The rabbit holes and the no-gos\n"+p.open+"architecture change: no\n")},
		{"facts", strings.HasSuffix(assembled, "## The facts\n"+p.facts)},
	} {
		if !c.ok {
			fs = append(fs, spec.Finding{Rule: "fixed", Detail: "the brief lacks the " + c.what + " part of the inputs as they stand"})
		}
	}
	return fs
}

// optionWords is the option-word rule (a partial evidence that the brief recommends no option, not the requirement itself): outside the
// parts of code, no word equals an option of a decision that has more than one survivor.
func (k briefKeeper) optionWords(assembled string) []spec.Finding {
	rest := assembled
	for _, part := range []string{k.parts.requirements, k.parts.solution, k.parts.open, k.parts.facts, "architecture change: no\n"} {
		rest = strings.Replace(rest, part, "", 1)
	}
	words := map[string]bool{}
	for _, w := range strings.FieldsFunc(strings.ToLower(rest), func(r rune) bool { return !(r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-') }) {
		words[w] = true
	}
	var fs []spec.Finding
	for _, d := range slices.Sorted(maps.Keys(k.multi)) {
		for _, o := range k.multi[d] {
			if words[strings.ToLower(o)] {
				fs = append(fs, spec.Finding{Rule: "option", Detail: fmt.Sprintf("the brief names the option %s of %s outside code's part", o, d)})
			}
		}
	}
	return fs
}

// keep is the Keep of the session: the files that the session may not write lines in, then the judge on the assembled brief; a pass
// keeps the candidate under the attempts of the milestone.
func (k briefKeeper) keep(name, id string, files map[string][]byte) (map[string][]byte, error) {
	text := map[string]string{}
	var fs []spec.Finding
	for _, n := range briefFiles {
		text[n] = string(files[n])
		if text[n] != "" && !strings.HasSuffix(text[n], "\n") {
			text[n] += "\n"
		}
		for _, prefix := range []string{"## ", "architecture change:", "approve batch "} {
			if slices.ContainsFunc(strings.Split(text[n], "\n"), func(l string) bool { return strings.HasPrefix(l, prefix) }) {
				fs = append(fs, spec.Finding{Rule: "file", Detail: fmt.Sprintf("%s has a line that starts with %q", n, prefix)})
			}
		}
	}
	assembled := assemble(k.parts, text["problem.md"], text["appetite.md"], text["rabbit-holes.md"])
	if len(fs) == 0 {
		fs = k.judge(assembled)
	}
	if len(fs) > 0 {
		return nil, errors.New(joinFindings(fs))
	}
	return map[string][]byte{attemptsDir(name) + "/" + id + ".md": []byte(assembled)}, nil
}

// appetiteFindings judges the appetite line of a brief against the band, the hours and the known money of milestone 0, which the
// telemetry holds (the session's own row too, so it is read after the session).
func appetiteFindings(brief string, telemetry [][]string, bandU, hours string) ([]spec.Finding, error) {
	b, _ := spec.CheckBrief([]byte(brief), nil, nil, nil)
	known, err := ledger.KnownMoney(telemetry)
	if err != nil {
		return nil, err
	}
	return spec.CheckAppetite(nil, spec.Cap{USD: b.USD, Hours: b.Hours}, bandU, hours, known), nil
}

// briefText is the prompt of the session: the block bet-brief with the paths of its files, and after a refusal the finding.
func briefText(name string, lines, room int, finding string) func(input map[string]string, output string) []byte {
	return func(input map[string]string, output string) []byte {
		note := ""
		if finding != "" {
			note = strings.NewReplacer("{finding}", finding).Replace(briefRefused)
		}
		return []byte(strings.NewReplacer("{milestone}", name, "{lines}", strconv.Itoa(lines), "{room}", strconv.Itoa(room),
			"{problem}", input["problem-statement.md"], "{spans}", input["spans.tsv"], "{draft}", input["draft.tsv"], "{questions}", input["questions.tsv"],
			"{answers}", input["answers.tsv"], "{budget}", input["budget.tsv"], "{requirements}", input["requirements.md"], "{solution}", input["solution.md"],
			"{open}", input["open.md"], "{facts}", input["facts.md"], "{links}", input["links.md"], "{output}", output, "{refused}", note).Replace(briefPrompt))
	}
}

// verifiedHead gives the bound head and the base of the open attempt of a Shape task whose position is verified; it reads the records
// and drives nothing.
func (s *Sessions) verifiedHead(ctx context.Context, task string) (head, base string, err error) {
	s.fns()
	v, err := s.shapeView(ctx, task)
	if err != nil {
		return "", "", err
	}
	if st := loopPosition(v); st.Pos == posVerified {
		return st.Head, st.Base, nil
	}
	return "", "", fmt.Errorf("%s: Shape has no verified head for it", task)
}

// briefInputs is everything a fresh build of the brief reads, before any session.
type briefInputs struct {
	inputs       map[string][]byte // the six records byte for byte and the five texts of code, by the name of their input file
	keeper       briefKeeper
	reqs, open   []string
	results      map[string]spec.Result
	heads        []headRef
	facts        [][3]string
	bandU, hours string
}

// readBrief reads a record that must exist, with its reader (none for a record of bytes alone).
func (r *state) readBrief(ctx context.Context, path string, read func([]byte) ([][]string, error)) ([][]string, []byte, error) {
	data, err := r.store.ReadFile(ctx, path)
	if err != nil {
		return nil, nil, err
	} else if data == nil {
		return nil, nil, fmt.Errorf("the records have no %s", path)
	}
	if read == nil {
		return nil, data, nil
	}
	rows, err := read(data)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", path, err)
	}
	return rows, data, nil
}

// briefInputs reads the inputs of the brief in the order of loop.md, The brief of the first bet, and writes nothing.
func (r *state) briefInputs(ctx context.Context) (*briefInputs, error) {
	if bets, err := r.store.ReadFile(ctx, "bets.tsv"); err != nil {
		return nil, err
	} else if bets != nil {
		return nil, errors.New("bets.tsv exists: the brief of a later bet is built by no row of M2e")
	}
	in := &briefInputs{inputs: map[string][]byte{}}
	rows := map[string][][]string{}
	for _, c := range []struct {
		name, path string
		read       func([]byte) ([][]string, error)
	}{{"problem-statement.md", "start/problem-statement.md", nil}, {"spans.tsv", "spec/spans.tsv", records.ReadSpans}, {"draft.tsv", "spec/draft.tsv", records.ReadDraft},
		{"questions.tsv", "questions.tsv", records.ReadQuestions}, {"answers.tsv", "intake/answers.tsv", records.ReadAnswers}, {"budget.tsv", "budget.tsv", records.ReadBudget}} {
		var err error
		if rows[c.name], in.inputs[c.name], err = r.readBrief(ctx, c.path, c.read); err != nil {
			return nil, err
		}
	}
	for _, b := range rows["budget.tsv"] {
		switch b[0] {
		case "band.U":
			in.bandU = b[1]
		case "appetite.hours":
			in.hours = b[1]
		}
	}
	if in.bandU == "" || in.hours == "" {
		return nil, errors.New("budget.tsv has no band.U or no appetite.hours")
	}
	_, params, err := r.readBrief(ctx, parametersPath, nil)
	if err != nil {
		return nil, err
	}
	table, err := ReadParameters(params)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parametersPath, err)
	}
	for _, p := range table {
		if n, err := strconv.Atoi(p.Value); p.Name == "brief.lines" && err == nil {
			in.keeper.lines = &n
		}
	}
	return in, r.briefShape(ctx, in, rows)
}

// briefShape reads the Shape tasks, the PRD at the verified head, the options and the links, and composes the texts of code.
func (r *state) briefShape(ctx context.Context, in *briefInputs, rows map[string][][]string) error {
	data, err := r.store.ReadFile(ctx, tasksPath)
	if err != nil {
		return err
	}
	register, err := ReadRegister(data)
	if err != nil {
		return fmt.Errorf("%s: %w", tasksPath, err)
	}
	s, blobs := r.sessions(), map[string][]string{}
	var specTask Task
	for _, kind := range []string{"specification", "architecture"} {
		i := slices.IndexFunc(register, func(t Task) bool { return t.Kind == kind })
		if i < 0 {
			return fmt.Errorf("the register has no task of the kind %s", kind)
		}
		t := register[i]
		head, base, err := s.verifiedHead(ctx, t.ID)
		if err != nil {
			return err
		}
		pr, found, err := r.cfg.Forge.PullRequest(ctx, forge.PullRequestQuery{Number: t.PR})
		if t.PR == 0 {
			return fmt.Errorf("%s has no pull request in the register", t.ID)
		} else if err != nil {
			return err
		} else if !found || pr.HeadSHA != head {
			return fmt.Errorf("%s: the pull request #%d is not found at the verified head %s", t.ID, t.PR, head)
		}
		names, err := git.DiffNames(r.target, base, head)
		if err != nil {
			return err
		}
		for _, p := range names {
			if es, err := git.LsTree(r.target, head, p); err != nil {
				return err
			} else if len(es) == 1 && es[0].Path == p && es[0].Type == "blob" && es[0].Mode == "100644" {
				blobs[t.ID] = append(blobs[t.ID], p)
			}
		}
		in.heads = append(in.heads, headRef{t.ID, t.Kind, head, t.PR})
		if kind == "specification" {
			specTask = t
		}
	}
	if ids := finalIDs(rows["draft.tsv"]); !slices.Equal(ids, specTask.Requirements) {
		return fmt.Errorf("the draft lists the requirements %v and the register row of %s lists %v", ids, specTask.ID, specTask.Requirements)
	}
	prds, err := filesAt(r.target, in.heads[0].Head, "docs/prd", isPRDName)
	if err != nil {
		return err
	}
	requirements, err := requirementsPart(specTask.Requirements, spec.ContentOfPRDs(prds))
	if err != nil {
		return err
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return err
	}
	cs, err := r.constraintsAt(ctx, base)
	if err != nil {
		return err
	}
	files, err := r.store.ReadFiles(ctx, "spec/options")
	if err != nil {
		return err
	}
	current, _ := optionTables(files)
	if in.results, err = filterAll(current, spec.IDs(cs)); err != nil {
		return err
	} else if len(in.results) == 0 {
		return errors.New("the records hold no options of a decision")
	}
	in.keeper.multi = map[string][]string{}
	for d, res := range in.results {
		if len(res.Survivors) == 0 {
			return fmt.Errorf("the decision %s has no survivor", d)
		} else if len(res.Survivors) > 1 {
			for _, o := range res.Options {
				in.keeper.multi[d] = append(in.keeper.multi[d], o.Name)
			}
		}
	}
	in.keeper.links = linkSet(r.cfg.Owner, r.cfg.Name, in.heads, blobs)
	in.keeper.parts = briefParts{requirements: requirements, solution: solutionPart(r.cfg.Owner, r.cfg.Name, in.heads, in.results),
		open: openPart(rows["questions.tsv"]), facts: factsPart(rows["spans.tsv"], rows["questions.tsv"], rows["answers.tsv"])}
	in.reqs, in.open, in.facts = specTask.Requirements, stillOpen(rows["questions.tsv"]), factRows(rows["spans.tsv"], rows["questions.tsv"], rows["answers.tsv"])
	for name, text := range map[string]string{"requirements.md": requirements, "solution.md": in.keeper.parts.solution, "open.md": in.keeper.parts.open,
		"facts.md": in.keeper.parts.facts, "links.md": strings.Join(slices.Sorted(maps.Keys(in.keeper.links)), "\n") + "\n"} {
		in.inputs[name] = []byte(text)
	}
	return nil
}

// basisRows gives the rows of bets/<milestone>/basis.tsv.
func (in *briefInputs) basisRows() [][]string {
	var rows [][]string
	for _, id := range in.reqs {
		rows = append(rows, []string{"requirement", id, "R-" + id[4:]})
	}
	for _, d := range slices.Sorted(maps.Keys(in.results)) {
		for _, o := range in.results[d].Options {
			if o.Dropped {
				rows = append(rows, []string{"dropped", d + "=" + o.Name, strings.Join(o.Reasons, "; ")})
			} else {
				rows = append(rows, []string{"survivor", d + "=" + o.Name, ""})
			}
		}
	}
	for _, f := range in.facts {
		rows = append(rows, []string{"fact", f[0], f[1]})
	}
	for _, id := range in.open {
		rows = append(rows, []string{"open", id, ""})
	}
	for _, h := range in.heads {
		rows = append(rows, []string{"head", h.Task, h.Head})
	}
	return rows
}

// candidate is a brief that passed the checks of a session and waits for the appetite and the promotion.
type candidate struct {
	id   string
	data []byte
}

// latestCandidate gives the candidate of the last session of sessions.tsv that has one, or none.
func (r *state) latestCandidate(ctx context.Context, name string) (candidate, error) {
	files, err := r.store.ReadFiles(ctx, attemptsDir(name))
	if err != nil {
		return candidate{}, err
	}
	sessions, err := r.sessions().readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return candidate{}, err
	}
	var c candidate
	for _, row := range sessions {
		if data, ok := files[row[sessionID]+".md"]; ok {
			c = candidate{row[sessionID], data}
		}
	}
	return c, nil
}

// judgeCandidate gives the finding of a candidate, or "": the checks of its bytes, then the appetite on the telemetry as it stands.
func (r *state) judgeCandidate(ctx context.Context, in *briefInputs, c candidate) (string, error) {
	if fs := in.keeper.judge(string(c.data)); len(fs) > 0 {
		return joinFindings(fs), nil
	}
	rows, err := r.sessions().readTable(ctx, "telemetry.tsv", records.ReadTelemetry)
	if err != nil {
		return "", err
	}
	fs, err := appetiteFindings(string(c.data), rows, in.bandU, in.hours)
	return joinFindings(fs), err
}

// betBrief keeps bets/<milestone>/brief.md for the first bet, or gives the bytes it kept before (done). It runs the session of the
// step bet-brief until a brief passes the checks and the appetite, and promotes it with its basis in one commit; a candidate that
// an earlier call left is judged again before any session starts.
func (r *state) betBrief(ctx context.Context, name, kind string) (brief []byte, done bool, err error) {
	if err := checkBriefRequest(name, kind); err != nil {
		return nil, false, err
	}
	if kept, err := r.store.ReadFile(ctx, briefPath(name)); err != nil || kept != nil {
		return kept, kept != nil, err
	}
	in, err := r.briefInputs(ctx)
	if err != nil {
		return nil, false, err
	}
	if err := in.keeper.overLimit(); err != nil {
		return nil, false, err
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "product-owner", r.cfg.Harnesses, r.cfg.Models, "")
	if errors.Is(err, route.ErrNoPair) {
		return nil, false, fmt.Errorf("pair: %v", err)
	} else if err != nil {
		return nil, false, err
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return nil, false, err
	}
	c, err := r.latestCandidate(ctx, name)
	if err != nil {
		return nil, false, err
	}
	room := *in.keeper.lines - strings.Count(assemble(in.keeper.parts, "", "", ""), "\n")
	for finding := ""; ; {
		if c.data != nil {
			if finding, err = r.judgeCandidate(ctx, in, c); err != nil {
				return nil, false, err
			} else if finding == "" {
				return c.data, false, r.promote(ctx, in, name, c)
			}
		}
		r.cfg.Progress("bet-brief: a session")
		id, err := s.PhaseSession(ctx, PhaseSpec{Step: "bet-brief", Role: "product-owner", Harness: h, Model: m, Models: r.cfg.Models, Base: base, Records: r.store.Base(),
			Input: in.inputs, Output: "brief", Files: briefFiles, Prompt: briefText(name, *in.keeper.lines, room, finding),
			KeepFiles: func(id string, files map[string][]byte) (map[string][]byte, error) {
				return in.keeper.keep(name, id, files)
			}})
		var pe *PhaseError
		if errors.As(err, &pe) && strings.HasPrefix(pe.Reason, "refused file: ") {
			finding, c = strings.TrimPrefix(pe.Reason, "refused file: "), candidate{}
			continue
		} else if err != nil {
			return nil, false, err
		}
		data, err := r.store.ReadFile(ctx, attemptsDir(name)+"/"+id+".md")
		if err != nil || data == nil {
			return nil, false, errors.Join(err, fmt.Errorf("the session %s kept no candidate", id))
		}
		c = candidate{id, data}
	}
}

// promote keeps the brief of the candidate and its basis in one fenced records commit.
func (r *state) promote(ctx context.Context, in *briefInputs, name string, c candidate) error {
	basis, err := table(BetBasisSchema, in.basisRows())
	if err != nil {
		return err
	}
	if _, err := ReadBetBasis(basis); err != nil {
		return err
	}
	return r.write(ctx, map[string][]byte{briefPath(name): c.data, basisPath(name): basis}, fmt.Sprintf("layup run: the brief of the bet %s (session %s)", name, c.id))
}

// briefBasis gives what the kept brief of the milestone rests on, from its basis and its appetite line, and filters nothing again.
func (r *state) briefBasis(ctx context.Context, name string) (BriefBasis, error) {
	var b BriefBasis
	if err := checkBriefRequest(name, "first"); err != nil {
		return b, err
	}
	_, brief, err := r.readBrief(ctx, briefPath(name), nil)
	if err != nil {
		return b, err
	}
	rows, _, err := r.readBrief(ctx, basisPath(name), ReadBetBasis)
	if err != nil {
		return b, err
	}
	b.Survivors, b.Dropped, b.Heads = map[string][]string{}, map[string]string{}, map[string]string{}
	for _, row := range rows {
		kind, id, value := row[0], row[1], row[2]
		decision, option, _ := strings.Cut(id, "=")
		switch kind {
		case "requirement":
			b.Requirements = append(b.Requirements, id)
		case "survivor":
			b.Survivors[decision] = append(b.Survivors[decision], option)
		case "dropped":
			b.Dropped[id] = value
		case "fact":
			b.Facts = append(b.Facts, spec.Fact{Number: id, Class: value})
		case "open":
			b.Open = append(b.Open, id)
		case "head":
			b.Heads[id] = value
		}
	}
	parsed, fs := spec.CheckBrief(brief, nil, nil, nil)
	// The lookups are absent here (no limit, no links, no verified batches), so only the structural rules count.
	if i := slices.IndexFunc(fs, func(f spec.Finding) bool {
		return f.Rule == "heading" || f.Rule == "appetite" || f.Rule == "architecture"
	}); i >= 0 {
		return b, fmt.Errorf("the kept brief %s is malformed: %s: %s", briefPath(name), fs[i].Rule, fs[i].Detail)
	}
	if parsed.USD == "" || parsed.Hours == "" {
		return b, fmt.Errorf("the kept brief %s has no appetite line", briefPath(name))
	}
	b.Appetite = spec.Cap{USD: parsed.USD, Hours: parsed.Hours}
	return b, nil
}
