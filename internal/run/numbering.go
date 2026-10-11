package run

// This file holds the steps numbering and spans of Intake (docs/spec/intake.md,
// The numbering; The span check; The steps of Intake; Input states; row 44 of
// the plan, task T-hiv9, #219): a session of a phase of product-owner splits
// the problem statement into spans, and code checks them.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// spansFile is the form of result/spans.tsv: the block spans without text,
// which code adds.
var spansFile = sessionSchema(records.SpansSchema, "text")

// numberingKeep gives the Keep of the step numbering on the problem statement
// psb: the session's file read by its form; each span's text read by code from
// the file at its offsets, which the writer of the table writes by the field
// rule of the column (each tab, CR and LF as one space); the rows as spec/spans.tsv, by
// the block's rules. A span partly out of the file gets the text of its part
// in the file, and the rule Bounds of the check reports it; one with no byte
// in the file has no text, and the block refuses it (decided here, task
// T-hiv9).
func numberingKeep(psb []byte) func(id string, data []byte) (map[string][]byte, error) {
	return func(_ string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, spansFile)
		if err != nil {
			return nil, err
		}
		out := make([][]string, len(rows))
		for i, r := range rows {
			start, _ := strconv.Atoi(r[0])
			end, _ := strconv.Atoi(r[1])
			text := ""
			if from, to := max(start, 0), min(end, len(psb)); from < to {
				text = string(psb[from:to]) // tsv.Write applies the field rule
			}
			out[i] = append(append([]string{}, r...), text)
		}
		table, err := table(records.SpansSchema, out)
		if err != nil {
			return nil, err
		}
		if _, err := records.ReadSpans(table); err != nil {
			return nil, err
		}
		return map[string][]byte{"spec/spans.tsv": table}, nil
	}
}

// numberingPrompt is the fixed text of the prompt of the numbering, with
// {input} and {output}; each path is on a line of its own.
const numberingPrompt = `You split a problem statement into its facts, for LAYUP. Change no file in
this repository and make no commit.

The problem statement: {input}
The output file: {output}

Write the output file as tab-separated values, with the header line
"start<TAB>end<TAB>fact<TAB>class<TAB>reason" and one row per span of the
problem statement, in the order of the file. A span is a part of the file by
its byte offsets: start from 0, end exclusive; a span may be part of a line.
Every byte that is not white space lies in exactly one span, and the spans do
not overlap. A fact has its number, 1 to n in the order of the spans, and its
class, one of need, constraint, measure, context, with the reason written as
"—"; a span that is not a fact (a heading, a rule line) has "—" as its fact and
its class, and its reason. Then end.
`

func numberingText(input map[string]string, output string) []byte {
	return []byte(strings.NewReplacer("{input}", input["problem-statement.md"], "{output}", output).Replace(numberingPrompt))
}

// checkFailed reports whether the records hold a spec/check.tsv with a fail
// row: then numbering and spans run again (intake.md, The span check).
func (r *state) checkFailed(ctx context.Context) (bool, error) {
	data, err := r.store.ReadFile(ctx, "spec/check.tsv")
	if err != nil || data == nil {
		return false, err
	}
	rows, err := tsv.Read(data, spec.CheckSchema)
	if err != nil {
		return false, fmt.Errorf("spec/check.tsv: %w", err)
	}
	for _, row := range rows {
		if row[1] == "fail" {
			return true, nil
		}
	}
	return false, nil
}

// numberingStep is the step numbering of Intake: "done before" when
// spec/spans.tsv exists and the last spec/check.tsv has no fail row; else one
// session of product-owner on the problem statement of the records, whose
// rows replace spec/spans.tsv.
func (r *state) numberingStep(ctx context.Context) Step {
	const name = "numbering"
	failed, err := r.checkFailed(ctx)
	if err != nil {
		return fail(name, err)
	}
	if data, err := r.store.ReadFile(ctx, "spec/spans.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil && !failed {
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
	h, m, err := s.PhasePair(ctx, "product-owner", r.cfg.Harnesses, r.cfg.Models, "")
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
		Base: base, Records: r.store.Base(), Input: map[string][]byte{"problem-statement.md": psb}, Output: "spans.tsv",
		Prompt: numberingText, Keep: numberingKeep(psb)}); err != nil {
		return fail(name, err)
	}
	data, err := r.store.ReadFile(ctx, "spec/spans.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadSpans(data)
	if err != nil {
		return fail(name, err)
	}
	facts := 0
	for _, row := range rows {
		if row[2] != "" {
			facts++
		}
	}
	return done(name, fmt.Sprintf("%d spans, %d facts, in spec/spans.tsv", len(rows), facts))
}

// spansStep is the step spans of Intake: "done before" when spec/check.tsv
// exists with no fail row; else the span check in process on spec/spans.tsv
// and the problem statement of the records, its table committed byte for byte
// as spec/check.tsv (tsv.Write by spec.CheckSchema); fail when a row fails,
// with the first fail row, so no step after it runs.
func (r *state) spansStep(ctx context.Context) Step {
	const name = "spans"
	failed, err := r.checkFailed(ctx)
	if err != nil {
		return fail(name, err)
	}
	if data, err := r.store.ReadFile(ctx, "spec/check.tsv"); err != nil {
		return fail(name, err)
	} else if data != nil && !failed {
		return done(name, "done before")
	}
	spansData, err := r.store.ReadFile(ctx, "spec/spans.tsv")
	if err != nil {
		return fail(name, err)
	}
	psb, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	spans, err := spec.Read(spansData)
	if err != nil {
		return fail(name, fmt.Errorf("spec/spans.tsv: %w", err))
	}
	rows := spec.Check(spans, psb)
	var buf bytes.Buffer
	if err := tsv.Write(&buf, spec.CheckSchema, rows); err != nil {
		return fail(name, err)
	}
	// A table equal to the last spec/check.tsv (a rerun whose session wrote
	// the same spans) is kept as it is: a records commit changes a file.
	if last, err := r.store.ReadFile(ctx, "spec/check.tsv"); err != nil {
		return fail(name, err)
	} else if !bytes.Equal(last, buf.Bytes()) {
		if err := r.write(ctx, map[string][]byte{"spec/check.tsv": buf.Bytes()}, "layup run: the span check of the facts"); err != nil {
			return fail(name, err)
		}
	}
	for _, row := range rows {
		if row[1] == "fail" {
			return failf(name, "the span check: %s fail: %s", row[0], row[2])
		}
	}
	return done(name, "the span check passed, in spec/check.tsv")
}
