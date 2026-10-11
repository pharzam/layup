package run

// This file holds the step draft of Intake (docs/spec/intake.md, The draft;
// The steps of Intake; Input states; row 45a of the plan, task T-2et0, #220):
// a session of a phase of product-owner drafts the requirements of the facts,
// and code checks the draft's coverage of the facts.

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/tsv"
)

// draftFile is the form of result/draft.tsv: the block draft without id and
// session, which code adds.
var draftFile = sessionSchema(records.DraftSchema, "id", "session")

// checkDraft is the check of the draft (code): each row's covers names at
// least one number, each a fact of spans (rows as records.ReadSpans gives
// them); each need and constraint fact is covered by a row; each measure fact
// is covered by a row of kind success-criterion. rows are the session's, as
// draftFile reads them: kind, statement, covers, criterion.
func checkDraft(rows, spans [][]string) error {
	class := map[string]string{} // the class of each fact, by its number
	var facts []string
	for _, s := range spans {
		if s[2] != "" {
			class[s[2]] = s[3]
			facts = append(facts, s[2])
		}
	}
	covered, criterion := map[string]bool{}, map[string]bool{}
	for i, r := range rows {
		covers := strings.Fields(r[2])
		if len(covers) == 0 {
			return fmt.Errorf("row %d: covers names no fact", i+1)
		}
		for _, n := range covers {
			if _, ok := class[n]; !ok {
				return fmt.Errorf("row %d: covers names %s, which is no fact", i+1, n)
			}
			covered[n] = true
			if r[0] == "success-criterion" {
				criterion[n] = true
			}
		}
	}
	for _, n := range facts {
		switch c := class[n]; {
		case (c == "need" || c == "constraint") && !covered[n]:
			return fmt.Errorf("the %s fact %s is covered by no row", c, n)
		case c == "measure" && !criterion[n]:
			return fmt.Errorf("the measure fact %s has no success criterion", n)
		}
	}
	return nil
}

// draftKeep gives the Keep of the step draft on the facts of spans: the
// session's file read by its form and checked; the rows with R-NNN in their
// order and the session ID, as spec/draft.tsv, by the block's reader.
func draftKeep(spans [][]string) func(id string, data []byte) (map[string][]byte, error) {
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, draftFile)
		if err != nil {
			return nil, err
		}
		if err := checkDraft(rows, spans); err != nil {
			return nil, err
		}
		out := make([][]string, len(rows))
		for i, r := range rows {
			out[i] = append(append([]string{fmt.Sprintf("R-%03d", i+1)}, r...), id)
		}
		table, err := table(records.DraftSchema, out)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadDraft(table); err != nil {
			return nil, err
		}
		return map[string][]byte{"spec/draft.tsv": table}, nil
	}
}

// draftPrompt is the fixed text of the prompt of the draft, with {input},
// {spans} and {output}; each path is on a line of its own.
const draftPrompt = `You draft the requirements of the facts of a problem statement, for LAYUP.
Change no file in this repository and make no commit.

The problem statement: {input}
The span table: {spans}
The output file: {output}

The span table holds the facts of the problem statement: each row a span by its
byte offsets, and for a fact its number (the column fact) and its class (need,
constraint, measure, context). Write the output file as tab-separated values,
with the header line "kind<TAB>statement<TAB>covers<TAB>criterion" and these
rows: one or more of kind requirement for each need fact; one of kind
non-functional for each constraint fact; one of kind success-criterion for each
measure fact. The statement is the requirement; covers is the numbers of the
facts that it covers, separated by one space; the criterion is its acceptance
criterion, or the measure of a success criterion. Give no priority. Then end.
`

func draftText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{spans}", input["spans.tsv"], "{output}", output).Replace(draftPrompt))
}

// draftStep is the step draft of Intake: "done before" when spec/draft.tsv
// exists; else one session of product-owner (the first admitted pair of that
// role, as for numbering, so both run on one harness) on the problem statement
// and spec/spans.tsv of the records (decided here, task T-2et0: the session
// names fact numbers, so its input holds the span table).
func (r *state) draftStep(ctx context.Context) Step {
	const name = "draft"
	if data, err := r.store.ReadFile(ctx, "spec/draft.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	psb, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	spansData, err := r.store.ReadFile(ctx, "spec/spans.tsv")
	if err != nil {
		return fail(name, err)
	}
	if psb == nil || spansData == nil {
		return failf(name, "the records have no start/problem-statement.md or spec/spans.tsv")
	}
	spans, err := records.ReadSpans(spansData)
	if err != nil {
		return fail(name, fmt.Errorf("spec/spans.tsv: %w", err))
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "product-owner", r.cfg.Harnesses, r.cfg.Models)
	if errors.Is(err, route.ErrNoPair) {
		return failf(name, "pair: %v", err)
	} else if err != nil {
		return fail(name, err)
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	if _, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "product-owner", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(), Input: map[string][]byte{"problem-statement.md": psb, "spans.tsv": spansData},
		Output: "draft.tsv", Prompt: draftText, Keep: draftKeep(spans)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "spec/draft.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadDraft(data)
	if err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d rows, in spec/draft.tsv", len(rows)))
}
