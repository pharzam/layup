package run

// This file holds the step completeness of Intake (docs/spec/intake.md, The
// completeness review; The steps of Intake; Input states; row 45b of the plan,
// task T-0mey, #221): a session of a phase of qa-engineer, on a harness other
// than the draft's, reviews the draft's coverage of the facts.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/tsv"
)

// coverageFile is the form of result/coverage.tsv: the block coverage without
// id and session, which code adds.
var coverageFile = sessionSchema(records.CoverageSchema, "id", "session")

// completenessKeep gives the Keep of the step completeness on the facts of
// spans and the rows of draft (as records.ReadSpans and records.ReadDraft give
// them): the session's file read by its form; each fact a fact of the span
// table and each requirement an ID of the draft (decided here, task T-0mey:
// each row becomes a question of the batch, which names what exists); the
// rows with C-NNN in their order and the session ID, as intake/coverage.tsv,
// by the block's reader.
func completenessKeep(spans, draft [][]string) func(id string, data []byte) (map[string][]byte, error) {
	var facts, ids []string
	for _, s := range spans {
		if s[2] != "" {
			facts = append(facts, s[2])
		}
	}
	for _, d := range draft {
		ids = append(ids, d[0])
	}
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, coverageFile)
		if err != nil {
			return nil, err
		}
		out := make([][]string, len(rows))
		for i, r := range rows {
			if fact := r[1]; fact != "" && !slices.Contains(facts, fact) {
				return nil, fmt.Errorf("row %d: the fact %s is no fact of spec/spans.tsv", i+1, fact)
			}
			if req := r[2]; req != "" && !slices.Contains(ids, req) {
				return nil, fmt.Errorf("row %d: the requirement %s is no row of spec/draft.tsv", i+1, req)
			}
			out[i] = append(append([]string{fmt.Sprintf("C-%03d", i+1)}, r...), id)
		}
		table, err := table(records.CoverageSchema, out)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadCoverage(table); err != nil {
			return nil, err
		}
		return map[string][]byte{"intake/coverage.tsv": table}, nil
	}
}

// completenessPrompt is the fixed text of the prompt of the completeness
// review, with {input}, {spans}, {draft} and {output}, each path on a line of
// its own.
const completenessPrompt = `You review a draft of requirements for its coverage of the facts of a problem
statement, for LAYUP. Change no file in this repository and make no commit.

The problem statement: {input}
The span table: {spans}
The draft: {draft}
The output file: {output}

The span table holds the facts (the column fact) and their classes; the draft
holds the requirements (the column id) and the facts that each covers. Write
the output file as tab-separated values, with the header line
"kind<TAB>fact<TAB>requirement<TAB>question" and one row per item: uncovered,
a need, constraint or measure fact that no requirement or success criterion
covers (its fact, the requirement "—"); unsupported, a requirement that no fact
supports (the fact "—", its requirement); class, a fact whose class you doubt
(its fact, the requirement "—"). Each with a question for the idea owner, on
one line. With no item, write the header line alone. Then end.
`

func completenessText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{spans}", input["spans.tsv"],
		"{draft}", input["draft.tsv"], "{output}", output).Replace(completenessPrompt))
}

// draftHarness gives the harness of the draft's session: the session column
// of spec/draft.tsv, its row of sessions.tsv (decided here, task T-0mey: the
// numbering's harness is the draft's by the routing rule of The draft, and the
// records name the draft's session; the span table names none).
func (r *state) draftHarness(ctx context.Context, draft [][]string) (string, error) {
	if len(draft) == 0 {
		return "", errors.New("spec/draft.tsv has no row")
	}
	s := r.sessions()
	start, err := s.startRowOf(ctx, draft[0][5])
	if err != nil {
		return "", err
	}
	return start[sessionHarness], nil
}

// completenessStep is the step completeness of Intake: "done before" when
// intake/coverage.tsv exists; else one session of qa-engineer, on the first
// admitted pair of the role whose harness is not the draft's; with none (no
// admitted pair at all, one admitted harness, or none other), fail "not run:
// no admitted pair on another harness", and no session starts.
func (r *state) completenessStep(ctx context.Context) Step {
	const name = "completeness"
	if data, err := r.store.ReadFile(ctx, "intake/coverage.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	var inputs [3][]byte
	for i, path := range []string{"start/problem-statement.md", "spec/spans.tsv", "spec/draft.tsv"} {
		data, err := r.store.ReadFile(ctx, path)
		if err != nil {
			return fail(name, err)
		}
		if data == nil {
			return failf(name, "the records have no %s", path)
		}
		inputs[i] = data
	}
	spans, err := records.ReadSpans(inputs[1])
	if err != nil {
		return fail(name, fmt.Errorf("spec/spans.tsv: %w", err))
	}
	draft, err := records.ReadDraft(inputs[2])
	if err != nil {
		return fail(name, fmt.Errorf("spec/draft.tsv: %w", err))
	}
	except, err := r.draftHarness(ctx, draft)
	if err != nil {
		return fail(name, err)
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "qa-engineer", r.cfg.Harnesses, r.cfg.Models, except)
	if errors.Is(err, route.ErrNoPair) {
		return failf(name, "not run: no admitted pair on another harness")
	} else if err != nil {
		return fail(name, err)
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	if _, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "qa-engineer", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(),
		Input:  map[string][]byte{"problem-statement.md": inputs[0], "spans.tsv": inputs[1], "draft.tsv": inputs[2]},
		Output: "coverage.tsv", Prompt: completenessText, Keep: completenessKeep(spans, draft)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "intake/coverage.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadCoverage(data)
	if err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d items, in intake/coverage.tsv", len(rows)))
}
