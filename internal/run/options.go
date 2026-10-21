package run

// This file holds the options of the step shape (docs/spec/loop.md, REQ-012 — Shape; task T-3f94, #364, row 66h of the plan of
// M2e): the constraints that apply, the options session and its output, and the filter of every decision.

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
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/verify"
)

// adrFile is the name of a decision record of the baseline: four digits, a dash, anything, .md.
var adrFile = regexp.MustCompile(`^[0-9]{4}-.*\.md$`)

// statusHeading is the heading that the baseline's linter finds a record's status after.
var statusHeading = regexp.MustCompile(`^## Status[[:space:]]*$`)

// adrsAt gives the decision records of the files of docs/adr/ (by base name) with their status lines, in the order of the names: the
// first non-empty line after the heading "## Status", as the baseline's linter reads it. A record with no such line is an error.
func adrsAt(files map[string][]byte) ([]spec.ADR, error) {
	var out []spec.ADR
	for _, name := range slices.Sorted(maps.Keys(files)) {
		if !adrFile.MatchString(name) {
			continue
		}
		status, in := "", false
		for _, line := range strings.Split(string(files[name]), "\n") {
			line = strings.TrimSuffix(line, "\r")
			if in && strings.TrimSpace(line) != "" {
				status = strings.TrimSpace(line)
				break
			}
			in = in || statusHeading.MatchString(line)
		}
		if status == "" {
			return nil, fmt.Errorf("docs/adr/%s has no status line", name)
		}
		out = append(out, spec.ADR{ID: "ADR-" + name[:4], Status: status})
	}
	return out, nil
}

// constraintsOf gives the constraints that apply from the rows of the draft and of the spans, the text of the target's guardrails.md
// and the files of its docs/adr/ by base name. A requirement is named by its final ID (NFR-NNN), as the Shape tasks name it.
func constraintsOf(draft, spans [][]string, guardrails []byte, adrFiles map[string][]byte) ([]spec.Constraint, error) {
	reqs := spec.RequirementsOfDraft(draft)
	for i, q := range reqs {
		if id, ok := finalID(q.ID, q.Kind); ok {
			reqs[i].ID = id
		}
	}
	adrs, err := adrsAt(adrFiles)
	if err != nil {
		return nil, err
	}
	return spec.Standing(reqs, spec.FactsOfSpans(spans), verify.InvariantIDs(guardrails), adrs), nil
}

// constraintsAt gives the constraints that apply, from the records and from the target at base. M2e runs Shape once, before the first
// bet, so the requirements are those of the draft and the facts those of the spans.
func (r *state) constraintsAt(ctx context.Context, base string) ([]spec.Constraint, error) {
	draft, err := r.readRecord(ctx, "spec/draft.tsv", records.ReadDraft)
	if err != nil {
		return nil, err
	}
	spans, err := r.readRecord(ctx, "spec/spans.tsv", records.ReadSpans)
	if err != nil {
		return nil, err
	}
	guardrails, err := git.Show(r.target, base, "docs/guardrails.md")
	if err != nil {
		return nil, fmt.Errorf("the target has no docs/guardrails.md at %s: %w", base, err)
	}
	adrs, err := filesAt(r.target, base, "docs/adr", func(name string) bool { return adrFile.MatchString(name) })
	if err != nil {
		return nil, err
	}
	return constraintsOf(draft, spans, guardrails, adrs)
}

// optionsHeader is the header line of the output of an options session and of a table of the block options.
const (
	optionsHeader = "decision\toption\tconstraint\tresult\tevidence\n"
	tableHeader   = "option\tconstraint\tresult\tevidence\n"
)

// reservedDecisions are the names that the prompts of the task steps use for an input file; no decision has one.
var reservedDecisions = []string{"draft", "spans", "answers", "tests", "tasks"}

// splitOptions checks the output of an options session and splits it into the table of each decision, in the form of the block
// options. It refuses an output with another header, no decision, a name that is no word of lowercase letters, digits and -, a table
// that the reader of the block refuses, and a decision with fewer than two options.
func splitOptions(data []byte) (map[string][]byte, error) {
	lines := strings.SplitAfter(string(data), "\n")
	if lines[0] != optionsHeader {
		return nil, fmt.Errorf("the header is not %q (the output starts %.60q)", strings.TrimSuffix(optionsHeader, "\n"), data)
	}
	tables := map[string]*strings.Builder{}
	for i, line := range lines[1:] {
		if line == "" {
			continue
		}
		decision, row, ok := strings.Cut(line, "\t")
		switch {
		case !ok:
			return nil, fmt.Errorf("line %d holds no tab", i+2)
		case !decisionName.MatchString(decision):
			return nil, fmt.Errorf("line %d: the decision %q is not a word of lowercase letters, digits and -", i+2, decision)
		case slices.Contains(reservedDecisions, decision):
			return nil, fmt.Errorf("line %d: the decision %q is reserved (the name of an input file of the task prompts)", i+2, decision)
		}
		if tables[decision] == nil {
			tables[decision] = &strings.Builder{}
			tables[decision].WriteString(tableHeader)
		}
		tables[decision].WriteString(row)
		if !strings.HasSuffix(row, "\n") {
			tables[decision].WriteString("\n")
		}
	}
	if len(tables) == 0 {
		return nil, errors.New("the output holds no decision")
	}
	out := map[string][]byte{}
	for decision, b := range tables {
		rows, err := spec.ReadOptions([]byte(b.String()))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", decision, err)
		}
		var names []string
		for _, r := range rows {
			if !slices.Contains(names, r[0]) {
				names = append(names, r[0])
			}
		}
		if len(names) < 2 {
			return nil, fmt.Errorf("decision %s has fewer than two options", decision)
		}
		out[decision] = []byte(b.String())
	}
	return out, nil
}

// optionsKeep is the Keep of the options session: the table of each decision as spec/options/<decision>.tsv.
func optionsKeep(id string, data []byte) (map[string][]byte, error) {
	tables, err := splitOptions(data)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for decision, t := range tables {
		files["spec/options/"+decision+".tsv"] = t
	}
	return files, nil
}

// optionTables splits the files of spec/options/ (by path below it) into the tables of the current round and those of round 1, by
// decision. A file below another directory is neither.
func optionTables(files map[string][]byte) (current, round1 map[string][]byte) {
	current, round1 = map[string][]byte{}, map[string][]byte{}
	decision := func(name string) (string, bool) {
		d, ok := strings.CutSuffix(name, ".tsv")
		return d, ok && !strings.Contains(d, "/")
	}
	for name, data := range files {
		if d, ok := decision(name); ok {
			current[d] = data
		} else if rest, ok := strings.CutPrefix(name, "round-1/"); ok {
			if d, ok := decision(rest); ok {
				round1[d] = data
			}
		}
	}
	return current, round1
}

// filterAll applies the filter to the table of each decision.
func filterAll(tables map[string][]byte, ids []string) (map[string]spec.Result, error) {
	out := map[string]spec.Result{}
	for _, d := range slices.Sorted(maps.Keys(tables)) {
		rows, err := spec.ReadOptions(tables[d])
		if err != nil {
			return nil, fmt.Errorf("spec/options/%s.tsv: %w", d, err)
		}
		out[d] = spec.Filter(rows, ids)
	}
	return out, nil
}

// noSurvivors gives, from the results of the filter, the decisions with no survivor that have no table of round 1 (the one rewrite is
// due) and those that have one (the round made no progress: the stop), each sorted.
func noSurvivors(round1 map[string][]byte, res map[string]spec.Result) (rewrite, stop []string) {
	for _, d := range slices.Sorted(maps.Keys(res)) {
		if len(res[d].Survivors) > 0 {
			continue
		}
		if _, has := round1[d]; has {
			stop = append(stop, d)
		} else {
			rewrite = append(rewrite, d)
		}
	}
	return rewrite, stop
}

// rewriteFiles checks the output of the rewrite session against the tables it replaces (by decision) and gives the records to keep:
// each new table, and the table it replaces as the table of round 1. The output holds exactly the decisions of old.
func rewriteFiles(old map[string][]byte, data []byte) (map[string][]byte, error) {
	tables, err := splitOptions(data)
	if err != nil {
		return nil, err
	}
	for d := range tables {
		if _, ok := old[d]; !ok {
			return nil, fmt.Errorf("the rewrite holds the decision %s, which it does not replace", d)
		}
	}
	files := map[string][]byte{}
	for d, was := range old {
		if _, ok := tables[d]; !ok {
			return nil, fmt.Errorf("the rewrite lacks the decision %s", d)
		}
		files["spec/options/"+d+".tsv"], files["spec/options/round-1/"+d+".tsv"] = tables[d], was
	}
	return files, nil
}

// droppedLines names each dropped option of a result with its reasons.
func droppedLines(res spec.Result) string {
	var b strings.Builder
	for _, o := range res.Options {
		if o.Dropped {
			fmt.Fprintf(&b, "- %s: %s\n", o.Name, strings.Join(o.Reasons, "; "))
		}
	}
	return b.String()
}

// noSurvivorComment gives the comment that names, for each decision, each dropped option of both rounds with its reasons.
func noSurvivorComment(decisions []string, round1, round2 map[string]spec.Result) string {
	var b strings.Builder
	b.WriteString(noSurvivorHeader + "\n")
	for _, d := range decisions {
		fmt.Fprintf(&b, "\n## %s\n\nRound 1:\n%s\nRound 2, the rewrite:\n%s", d, droppedLines(round1[d]), droppedLines(round2[d]))
	}
	return b.String()
}

const optionsPromptText = `You write the options of the architecture decisions of a target, for LAYUP. Change no file in this repository and make no commit.

The draft: {draft}
The spans: {spans}
The problem statement: {problem-statement}
The answers: {answers}
The output file: {output}

Find each architecture decision that the draft needs. For each decision write at least two options. Write the output file as
tab-separated values, with the header line
"decision<TAB>option<TAB>constraint<TAB>result<TAB>evidence"
and one row for each option and each constraint of the list below.
- The decision and the option are each a word of lowercase letters, digits and -. The decision is none of the names draft, spans, answers, tests and tasks.
- The constraint is an ID of the list below. An option has a row for every constraint of the list, and may have rows for others.
- The result is pass or fail. For a pass the evidence is the path of a file or a link that shows it, never a dash. For a fail the
  evidence is a dash, and nothing else.
{constraints}
The constraints that apply, by their IDs:
{list}
Then end.
`

const rewritePromptText = `You rewrite the options of architecture decisions of a target, for LAYUP. Change no file in this repository and make no commit.

The draft: {draft}
The problem statement: {problem-statement}
The reasons: {reasons}
The output file: {output}

Every option of the decisions below was dropped by the filter, for the reasons in the file above. For each decision write new
options, at least two, that meet the constraints of the list below. The tables that were dropped:
{tables}
Write the output file in the form of the tables, as tab-separated values with the header line
"decision<TAB>option<TAB>constraint<TAB>result<TAB>evidence"
and one row for each option and each constraint of the list, for these decisions and for no other. A pass has a path or a link as its
evidence, a fail has a dash. The decision and the option are each a word of lowercase letters, digits and -. The decision is none of the names draft, spans, answers, tests and tasks.

The constraints that apply, by their IDs:
{list}
Then end.
`

// constraintList is the list of the constraints, one on each line: the ID, and its criterion when it has one.
func constraintList(cs []spec.Constraint) string {
	var b strings.Builder
	for _, c := range cs {
		if c.Text != "" {
			fmt.Fprintf(&b, "- %s: %s\n", c.ID, c.Text)
		} else {
			fmt.Fprintf(&b, "- %s\n", c.ID)
		}
	}
	return b.String()
}

// optionsPrompt is the prompt of the options session, with the constraints that apply written into it.
func optionsPrompt(cs []spec.Constraint) func(input map[string]string, output string) []byte {
	return func(in map[string]string, output string) []byte {
		return []byte(strings.NewReplacer("{draft}", in["draft.tsv"], "{spans}", in["spans.tsv"], "{problem-statement}", in["problem-statement.md"],
			"{answers}", in["answers.tsv"], "{output}", output, "{constraints}", "", "{list}", constraintList(cs)).Replace(optionsPromptText))
	}
}

// rewritePrompt is the prompt of the rewrite session: the files of the dropped tables are the inputs named dropped-<decision>.tsv.
func rewritePrompt(cs []spec.Constraint) func(input map[string]string, output string) []byte {
	return func(in map[string]string, output string) []byte {
		var tables strings.Builder
		for _, name := range slices.Sorted(maps.Keys(in)) {
			if d, ok := strings.CutPrefix(strings.TrimSuffix(name, ".tsv"), "dropped-"); ok && strings.HasSuffix(name, ".tsv") {
				fmt.Fprintf(&tables, "%s: %s\n", d, in[name])
			}
		}
		return []byte(strings.NewReplacer("{draft}", in["draft.tsv"], "{problem-statement}", in["problem-statement.md"], "{reasons}", in["reasons.txt"],
			"{output}", output, "{tables}", tables.String(), "{list}", constraintList(cs)).Replace(rewritePromptText))
	}
}

// defaultConstraints gives the head of the default branch and the constraints that apply at it.
func (r *state) defaultConstraints(ctx context.Context) (string, []spec.Constraint, error) {
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return "", nil, err
	}
	cs, err := r.constraintsAt(ctx, base)
	return base, cs, err
}

// optionResults gives the results of the filter of each decision's current table in the records, as the step shape makes them and no session runs: the options of
// the architecture task of Shape and of the task that redoes it (loop.md, The step loop, item 9).
func (r *state) optionResults(ctx context.Context) (map[string]spec.Result, error) {
	_, cs, err := r.defaultConstraints(ctx)
	if err != nil {
		return nil, err
	}
	files, err := r.store.ReadFiles(ctx, "spec/options")
	if err != nil {
		return nil, err
	}
	current, _ := optionTables(files)
	return filterAll(current, spec.IDs(cs))
}

// optionsStage is the second stage of the step: it has the options written, filtered and, when a decision has no survivor, rewritten
// once; it gives the results of the filter, and, for a decision with no survivor after the rewrite, the detail of the stop (the
// comment is posted once).
func (r *state) optionsStage(ctx context.Context) (map[string]spec.Result, string, error) {
	base, cs, err := r.defaultConstraints(ctx)
	if err != nil {
		return nil, "", err
	}
	ids := spec.IDs(cs)
	tables := func() (current, round1 map[string][]byte, err error) {
		files, err := r.store.ReadFiles(ctx, "spec/options")
		current, round1 = optionTables(files)
		return current, round1, err
	}
	current, round1, err := tables()
	if err != nil {
		return nil, "", err
	}
	inputs := map[string][]byte{}
	for name, path := range map[string]string{"draft.tsv": "spec/draft.tsv", "spans.tsv": "spec/spans.tsv", "problem-statement.md": "start/problem-statement.md", "answers.tsv": "intake/answers.tsv"} {
		if inputs[name], err = r.store.ReadFile(ctx, path); err != nil {
			return nil, "", err
		} else if inputs[name] == nil && name != "spans.tsv" {
			return nil, "", fmt.Errorf("the records have no %s", path)
		}
	}
	session := func(step string, input map[string][]byte, prompt func(map[string]string, string) []byte, keep func(string, []byte) (map[string][]byte, error)) error {
		s := r.sessions()
		h, m, err := s.PhasePair(ctx, "systems-architect", r.cfg.Harnesses, r.cfg.Models, "")
		if errors.Is(err, route.ErrNoPair) {
			return fmt.Errorf("pair: %v", err)
		} else if err != nil {
			return err
		}
		_, err = s.PhaseSession(ctx, PhaseSpec{Step: step, Role: "systems-architect", Harness: h, Model: m, Models: r.cfg.Models, Base: base, Records: r.store.Base(),
			Input: input, Output: "options.tsv", Prompt: prompt, Keep: keep})
		return err
	}
	if len(current) == 0 {
		if err := session("options", inputs, optionsPrompt(cs), optionsKeep); err != nil {
			return nil, "", err
		}
		if current, round1, err = tables(); err != nil {
			return nil, "", err
		}
	}
	results, err := filterAll(current, ids)
	if err != nil {
		return nil, "", err
	}
	if rewrite, _ := noSurvivors(round1, results); len(rewrite) > 0 {
		old, input, reasons := map[string][]byte{}, map[string][]byte{"draft.tsv": inputs["draft.tsv"], "problem-statement.md": inputs["problem-statement.md"]}, strings.Builder{}
		for _, d := range rewrite {
			old[d], input["dropped-"+d+".tsv"] = current[d], current[d]
			fmt.Fprintf(&reasons, "decision %s\n%s\n", d, droppedLines(results[d]))
		}
		input["reasons.txt"] = []byte(reasons.String())
		if err := session("options-rewrite", input, rewritePrompt(cs), func(id string, data []byte) (map[string][]byte, error) { return rewriteFiles(old, data) }); err != nil {
			return nil, "", err
		}
		if current, round1, err = tables(); err != nil {
			return nil, "", err
		}
		if results, err = filterAll(current, ids); err != nil {
			return nil, "", err
		}
	}
	_, stop := noSurvivors(round1, results)
	if len(stop) == 0 {
		return results, "", nil
	}
	first, err := filterAll(round1, ids)
	if err != nil {
		return nil, "", err
	}
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return nil, "", errors.New("start.tsv has no number of the Intake issue")
	}
	body := noSurvivorComment(stop, first, results)
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return nil, "", err
	}
	if !slices.ContainsFunc(comments, func(c forge.Comment) bool { return c.App == r.cfg.Register.AppSlug && c.Body == body }) {
		if err := r.sessions().holdsLease(ctx); err != nil {
			return nil, "", err
		}
		if _, err := r.cfg.Forge.Comment(ctx, issue, body); err != nil {
			return nil, "", err
		}
	}
	return results, "waiting for the Operator: shape: no survivor", nil
}
