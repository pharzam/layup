package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// The Shape filter (docs/spec/loop.md, REQ-012 — Shape; task T-0pt8, #349):
// code, no model call. An architect session writes at least two options for a
// decision, each with a row for every constraint that applies; the filter
// reads that table, never the session's own judgement of it, and drops an
// option with a failed row, a missing row or a pass with no evidence. One
// survivor is chosen; several are shown with no recommendation (phase 2 has no
// panel); none is a round with no survivor. The rewrite once and the stop are
// the step shape of internal/run (row 66h).

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

// ReadOptions reads a table of options by OptionsSchema and checks the rules of
// its owner (records.md): an option's name is of the form above, the result is
// never — (the reader gives — as the empty value in a column that is not a key,
// and the block allows it in none of these), and a fail has — in evidence. A
// pass with — in evidence is valid: the filter drops its option. The rows keep
// the empty string for —. An error names the line and the column.
func ReadOptions(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, OptionsSchema)
	if err != nil {
		return nil, err
	}
	for i, r := range rows {
		line := i + 2 // the header is line 1
		switch {
		case !optionName.MatchString(r[optName]):
			return nil, fmt.Errorf("line %d, column option: %q is not a word of lowercase letters, digits and -", line, r[optName])
		case r[optResult] == "":
			return nil, fmt.Errorf("line %d, column result: — is not a result, which is pass or fail", line)
		case r[optResult] == "fail" && r[optEvidence] != "":
			return nil, fmt.Errorf("line %d, column evidence: a fail has — in evidence, not %q", line, r[optEvidence])
		}
	}
	return rows, nil
}

// Fact is a fact with its class: a span row of the block spans (Number is its
// fact number) or a row of the confirmed inventory (Number is a number of the
// spans, or the ID Q-NNN of a fact of the answers). The caller chooses the
// source (docs/spec/loop.md: the span facts before the first bet, the
// inventory's after it); row 66h does.
type Fact struct{ Number, Class string }

// FactsOfSpans gives the fact rows of a span table (rows as Read gives them), in
// the order of the table. A span that is not a fact has none.
func FactsOfSpans(spans [][]string) []Fact {
	var out []Fact
	for _, r := range spans {
		if r[colFact] != "" {
			out = append(out, Fact{r[colFact], r[colClass]})
		}
	}
	return out
}

// Requirement is a requirement of the draft or of the PRD with its kind and its
// acceptance criterion. The draft's rows give them by RequirementsOfDraft; the
// PRD's come from the caller (the reader of the PRD, like those of the
// invariants and of the decision records, is the step shape's, row 66h).
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

// ADR is a decision record of the target with its status line. A record applies
// when it is accepted: the status is Accepted, or Accepted. Amended by
// ADR-NNNN, as the baseline's record rules write an extended record, whose old
// decision still holds. A proposed, rejected or superseded record does not.
type ADR struct{ ID, Status string }

// Constraint is a constraint that applies: its ID, and the text that the
// architect session reads (the criterion of a requirement, "" for the others,
// whose text is in the records they name).
type Constraint struct{ ID, Text string }

// Standing lists the constraints that apply to a decision, from the records and
// never from the session (loop.md, Shape): each requirement of kind
// non-functional, by its ID and its criterion; each fact of class constraint,
// by F-<number>, and a fact of the answers by its own ID Q-NNN; each invariant
// of the PSB; and each accepted decision record. The order is that of the four
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
			id := "F-" + f.Number
			if strings.HasPrefix(f.Number, "Q-") {
				id = f.Number
			}
			out = append(out, Constraint{id, ""})
		}
	}
	for _, v := range invariants {
		out = append(out, Constraint{v, ""})
	}
	for _, a := range adrs {
		if a.Status == "Accepted" || strings.HasPrefix(a.Status, "Accepted. Amended by ") {
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

// Option is one option of a decision with the reasons that drop it, each a
// constraint and what is wrong with it: those of the list first, in its order,
// then any other constraint that the table names.
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

// Chosen gives the survivor when there is exactly one. With none or with several
// it gives none: several are shown with no recommendation and the bet decides,
// none is a round with no survivor.
func (r Result) Chosen() (string, bool) {
	if len(r.Survivors) == 1 {
		return r.Survivors[0].Name, true
	}
	return "", false
}

// Filter applies the rule to the rows of an options table (as ReadOptions gives
// them) and the list of the constraints that apply (IDs of Standing). An option
// is dropped when a row is fail, when a constraint of the list has no row of it,
// or when a pass has no evidence. A constraint that is not on the list, which
// the architect names, is judged too (a fail, or a pass with no evidence, drops
// the option), but no row is demanded for it.
func Filter(rows [][]string, constraints []string) Result {
	var res Result
	idx := map[string]int{}
	byOpt := map[string]map[string][]string{}
	named := map[string][]string{}
	for _, r := range rows {
		name := r[optName]
		if _, ok := idx[name]; !ok {
			idx[name] = len(res.Options)
			res.Options = append(res.Options, Option{Name: name})
			byOpt[name] = map[string][]string{}
		}
		byOpt[name][r[optConstraint]] = r
		named[name] = append(named[name], r[optConstraint])
	}
	for i := range res.Options {
		o := &res.Options[i]
		judge := func(c string, r []string) {
			switch {
			case r == nil:
				o.Reasons = append(o.Reasons, c+": missing")
			case r[optResult] == "fail":
				o.Reasons = append(o.Reasons, c+": fail")
			case r[optEvidence] == "":
				o.Reasons = append(o.Reasons, c+": pass with no evidence")
			}
		}
		listed := map[string]bool{}
		for _, c := range constraints {
			listed[c] = true
			judge(c, byOpt[o.Name][c])
		}
		for _, c := range named[o.Name] {
			if !listed[c] {
				listed[c] = true
				judge(c, byOpt[o.Name][c])
			}
		}
		o.Dropped = len(o.Reasons) > 0
		if !o.Dropped {
			res.Survivors = append(res.Survivors, *o)
		}
	}
	return res
}
