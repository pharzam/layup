package spec

import (
	"fmt"
	"regexp"

	"github.com/pharzam/layup/internal/tsv"
)

// The Shape filter (docs/spec/loop.md, REQ-012 — Shape; task T-0pt8, #349):
// code, no model call. An architect session writes at least two options for a
// decision, each with a row for every constraint that applies; the filter
// reads that table, never the session's own judgement of it, and drops an
// option with a failed row, a missing row or a pass with no evidence. What it
// leaves are the survivors: one is chosen, more than one are shown with no
// recommendation (phase 2 has no panel), none is a round with no survivor. The
// rewrite once and the stop are the step shape of internal/run (row 66h).

// OptionsSchema is the form of the table of the options of one decision: the
// block options (docs/spec/records.md, NFR-001 — The records of the task loop).
var OptionsSchema = tsv.Schema{Name: "options", Location: "records:spec/options/<decision>.tsv", Columns: []tsv.Column{
	{Name: "option", Type: "text", Key: true},
	{Name: "constraint", Type: "text", Key: true},
	{Name: "result", Type: "enum(pass|fail)"},
	{Name: "evidence", Type: "text"},
}}

// The columns of an options row.
const (
	optName = iota
	optConstraint
	optResult
	optEvidence
)

// optionName is the form of an option's name: lowercase letters, digits and -,
// so that <decision>=<option> holds no space (loop.md, Shape).
var optionName = regexp.MustCompile(`^[a-z0-9-]+$`)

// ReadOptions reads a table of options by OptionsSchema and checks the rules
// of its owner, which the block gives in words (records.md): an option's name
// is of the form above, the result is never — (the reader gives — as the empty
// value in any column that is not a key, and the block allows it in none of
// these), and a fail has — in evidence. A pass with — in evidence is valid: the
// filter drops its option. The rows keep the decoded empty string for —. An
// error names the line and the column.
func ReadOptions(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, OptionsSchema)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		line := i + 2 // the header is line 1
		if !optionName.MatchString(r[optName]) {
			return nil, fmt.Errorf("line %d, column option: %q is not a word of lowercase letters, digits and -", line, r[optName])
		}
		if r[optResult] == "" {
			return nil, fmt.Errorf("line %d, column result: — is not a result, which is pass or fail", line)
		}
		if r[optResult] == "fail" && r[optEvidence] != "" {
			return nil, fmt.Errorf("line %d, column evidence: a fail has — in evidence, not %q", line, r[optEvidence])
		}
	}
	return rows, nil
}

// Fact is a numbered fact of the problem statement with its class, as a span
// row of the block spans or a row of the confirmed inventory gives it.
type Fact struct{ Number, Class string }

// FactsOfSpans gives the fact rows of a span table (rows as Read gives them),
// in the order of the table. A span that is not a fact has none.
func FactsOfSpans(spans [][]string) []Fact {
	var out []Fact
	for _, r := range spans {
		if r[colFact] != "" {
			out = append(out, Fact{r[colFact], r[colClass]})
		}
	}
	return out
}

// Requirement is a requirement of the draft or of the PRD with its kind
// (non-functional, or another) and its acceptance criterion. The draft's rows
// give them by RequirementsOfDraft; the PRD's come from the caller, as the
// reader of the PRD is the step shape's (docs/plan/README.md, row 66h).
type Requirement struct{ ID, Kind, Criterion string }

// RequirementsOfDraft gives the requirements of a draft (rows as the block draft
// gives them: id, kind, statement, covers, criterion, session), in its order.
func RequirementsOfDraft(draft [][]string) []Requirement {
	const colID, colKind, colCriterion = 0, 1, 4
	var out []Requirement
	for _, r := range draft {
		out = append(out, Requirement{r[colID], r[colKind], r[colCriterion]})
	}
	return out
}

// ADR is a decision record of the target with its status; only an Accepted one
// applies.
type ADR struct{ ID, Status string }

// Constraint is a constraint that applies to a decision: its ID, and the text
// that the architect session reads (the criterion of a requirement; "" for the
// others, whose text is in the records they name).
type Constraint struct{ ID, Text string }

// Standing lists the constraints that apply to a decision, from the records and
// never from the session (loop.md, Shape): each requirement of kind
// non-functional, by its ID and its criterion; each fact of class constraint,
// by F-<number>; each invariant of the PSB (Inv-<n>); and each decision record
// whose status is Accepted. Before the first bet the facts are those of the
// span table (FactsOfSpans); after it the caller gives those of the confirmed
// inventory instead, which row 62b builds. The order is the order of the four
// kinds and, in each, of the input, so two runs on one input give one list.
func Standing(reqs []Requirement, facts []Fact, invariants []string, adrs []ADR) []Constraint {
	var out []Constraint
	for _, r := range reqs {
		if r.Kind == "non-functional" {
			out = append(out, Constraint{r.ID, r.Criterion})
		}
	}
	for _, f := range facts {
		if f.Class == "constraint" {
			out = append(out, Constraint{"F-" + f.Number, ""})
		}
	}
	for _, v := range invariants {
		out = append(out, Constraint{v, ""})
	}
	for _, a := range adrs {
		if a.Status == "Accepted" {
			out = append(out, Constraint{a.ID, ""})
		}
	}
	return out
}

// IDs gives the IDs of the constraints, the list that Filter takes.
func IDs(cs []Constraint) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.ID
	}
	return out
}

// Option is one option of a decision with the reasons that drop it: a
// constraint and what is wrong with it, in the order of the constraints (those
// of the list first, then any other constraint that the table names).
type Option struct {
	Name    string
	Dropped bool
	Reasons []string
}

// Result is what the filter gives: each option in the order of its first row,
// and the survivors among them.
type Result struct {
	Options   []Option
	Survivors []Option
}

// Chosen gives the survivor when there is exactly one (it is chosen). With none
// or with several it gives none: several are shown with no recommendation and
// the bet decides, none is a round with no survivor.
func (r Result) Chosen() (string, bool) {
	if len(r.Survivors) == 1 {
		return r.Survivors[0].Name, true
	}
	return "", false
}

// Filter applies the rule to the rows of an options table (as ReadOptions gives
// them) and the list of the constraints that apply (Standing). An option is
// dropped when a row is fail, when a constraint of the list has no row of it,
// or when a pass has no evidence. A constraint that is not on the list, which
// the architect names, is judged too (a fail or a pass with no evidence drops
// the option), but no row is demanded for it.
func Filter(rows [][]string, constraints []string) Result {
	var res Result
	idx := map[string]int{}
	byOpt := map[string]map[string][]string{}
	order := map[string][]string{}
	for _, r := range rows {
		name := r[optName]
		if _, ok := idx[name]; !ok {
			idx[name] = len(res.Options)
			res.Options = append(res.Options, Option{Name: name})
			byOpt[name] = map[string][]string{}
		}
		byOpt[name][r[optConstraint]] = r
		order[name] = append(order[name], r[optConstraint])
	}
	for i := range res.Options {
		o := &res.Options[i]
		seen := map[string]bool{}
		judge := func(c string, r []string, ok bool) {
			switch {
			case !ok:
				o.Reasons = append(o.Reasons, c+": missing")
			case r[optResult] == "fail":
				o.Reasons = append(o.Reasons, c+": fail")
			case r[optEvidence] == "":
				o.Reasons = append(o.Reasons, c+": pass with no evidence")
			}
		}
		for _, c := range constraints {
			seen[c] = true
			r, ok := byOpt[o.Name][c]
			judge(c, r, ok)
		}
		for _, c := range order[o.Name] {
			if !seen[c] {
				seen[c] = true
				judge(c, byOpt[o.Name][c], true)
			}
		}
		o.Dropped = len(o.Reasons) > 0
		if !o.Dropped {
			res.Survivors = append(res.Survivors, *o)
		}
	}
	return res
}
