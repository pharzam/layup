package run

// This file holds the step meaning of Intake (docs/spec/intake.md, The review
// of meaning; The steps of Intake; Input states; row 43b of the plan, task
// T-1u2v, #218): a session of a phase of systems-architect writes the gaps of
// meaning of the problem statement, and code finds each quote in the file.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/tsv"
)

// sessionSchema gives the schema of a session's file: the block s without the
// columns drop, which code adds, in the block's order; a key column that is
// kept stays a key (docs/spec/intake.md, The sessions of Intake; decided
// here, task T-1u2v).
func sessionSchema(s tsv.Schema, drop ...string) tsv.Schema {
	out := tsv.Schema{Name: s.Name, Location: "host:sessions/<session>/result/" + s.Name + ".tsv"}
	for _, c := range s.Columns {
		if !slices.Contains(drop, c.Name) {
			out.Columns = append(out.Columns, c)
		}
	}
	return out
}

// quoteAt gives the byte offsets of quote in the problem statement psb: a
// quote is text of one line of the file, byte for byte, found exactly once.
// One that is empty, holds a tab or a line break, occurs nowhere or occurs
// more than once is refused, with its reason.
func quoteAt(psb []byte, quote string) (int, int, error) {
	switch {
	case quote == "":
		return 0, 0, errors.New("the quote is empty")
	case strings.Contains(quote, "\t"):
		return 0, 0, fmt.Errorf("the quote %q holds a tab", quote)
	case strings.ContainsAny(quote, "\r\n"):
		return 0, 0, fmt.Errorf("the quote %q holds a line break", quote)
	}
	q := []byte(quote)
	// Each place counts, those that overlap too ("--" is twice in "---"),
	// as a quote must name one place (round 1 of #218); bytes.Count does not.
	n, first := 0, -1
	for i := 0; i <= len(psb)-len(q); {
		j := bytes.Index(psb[i:], q)
		if j < 0 {
			break
		}
		if n++; n == 1 {
			first = i + j
		}
		i += j + 1
	}
	switch n {
	case 0:
		return 0, 0, fmt.Errorf("the quote %q occurs nowhere", quote)
	case 1:
		return first, first + len(q), nil
	default:
		return 0, 0, fmt.Errorf("the quote %q occurs %d times", quote, n)
	}
}

// meaningFile is the form of result/meaning.tsv.
var meaningFile = sessionSchema(records.MeaningSchema, "id", "start", "end", "session")

// meaningKeep gives the Keep of the step meaning on the problem statement psb:
// the session's file read by its form; each quote found by quoteAt; the rows
// with M-NNN in their order, the offsets and the session ID, as
// intake/meaning.tsv. A refused quote or row refuses the file.
func meaningKeep(psb []byte) func(id string, data []byte) (map[string][]byte, error) {
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, meaningFile)
		if err != nil {
			return nil, err
		}
		var out [][]string
		for i, r := range rows {
			kind, quote, question := r[0], r[1], r[2]
			start, end, err := quoteAt(psb, quote)
			if err != nil {
				return nil, fmt.Errorf("row %d: %w", i+1, err)
			}
			row := []string{fmt.Sprintf("M-%03d", i+1), kind, strconv.Itoa(start), strconv.Itoa(end), quote, question, id}
			if err := records.CheckMeaning(row); err != nil {
				return nil, fmt.Errorf("row %d: %w", i+1, err)
			}
			out = append(out, row)
		}
		table, err := table(records.MeaningSchema, out)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadMeaning(table); err != nil {
			return nil, err
		}
		return map[string][]byte{"intake/meaning.tsv": table}, nil
	}
}

// meaningPrompt is the fixed text of the prompt of the review of meaning, with
// {input} and {output}; each path is on a line of its own.
const meaningPrompt = `You review a problem statement for its gaps of meaning, for LAYUP. Change no
file in this repository and make no commit.

The problem statement: {input}
The output file: {output}

Write the output file as tab-separated values, with the header line
"kind<TAB>quote<TAB>question" and one row per gap of meaning that you find.
The kind is one of: two-readings (a term with two readings), unmeasurable (a
metric that cannot be measured as written), missing-intent (a missing intent
value), new-work (new technical work), fit (an assumption about how parts fit),
no-design (a design that may not exist), decision (a hard decision to settle
before the work). The quote is text of one line of the problem statement, byte
for byte, with no tab, that occurs once in the file: quote enough to be unique.
The question is the question for the idea owner, on one line. With no gap,
write the header line alone. Then end.
`

func meaningText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{output}", output).Replace(meaningPrompt))
}

// sessions gives the Sessions of the run's steps of a session.
func (r *state) sessions() *Sessions {
	control, _ := strconv.Atoi(r.value("issue.control"))
	return &Sessions{Store: r.store, RunID: r.cfg.RunID, Host: r.cfg.Dir, Target: r.cfg.Owner + "/" + r.cfg.Name, Clone: r.target,
		Branch: r.defaultBranch, Now: r.cfg.Clock.Now, Forge: r.cfg.Forge, Control: control, Prices: r.cfg.Prices}
}

// meaningStep is the step meaning of Intake: "done before" when
// intake/meaning.tsv exists; else one session of systems-architect on the
// problem statement of the records. No admitted pair is fail (pair); a session
// that ends other than done, a refused file and the cap of milestone 0 are
// fail with their text.
func (r *state) meaningStep(ctx context.Context) Step {
	const name = "meaning"
	if data, err := r.store.ReadFile(ctx, "intake/meaning.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil {
		return done(name, "done before")
	}
	psb, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	if psb == nil {
		return failf(name, "the records have no start/problem-statement.md")
	}
	s := r.sessions()
	h, m, err := s.PhasePair(ctx, "systems-architect", r.cfg.Harnesses, r.cfg.Models)
	if errors.Is(err, route.ErrNoPair) {
		return failf(name, "pair: %v", err)
	} else if err != nil {
		return fail(name, err)
	}
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		return fail(name, err)
	}
	if _, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "systems-architect", Harness: h, Model: m, Models: r.cfg.Models,
		Base: base, Records: r.store.Base(), Input: map[string][]byte{"problem-statement.md": psb}, Output: "meaning.tsv",
		Prompt: meaningText, Keep: meaningKeep(psb)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "intake/meaning.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadMeaning(data)
	if err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d gaps of meaning, in intake/meaning.tsv", len(rows)))
}
