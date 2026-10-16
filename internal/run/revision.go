package run

// This file holds the step revision of Intake (docs/spec/loop.md, The
// revision of the draft; docs/spec/intake.md, The steps of Intake; row 66i of
// the plan, task T-8rm6, #365): after the answers, a session of a phase of
// product-owner revises the requirement draft on the answers of the batch;
// code's check binds each fact of an answer that a row covers to the fact of
// the answer's quote, refuses a class change that no answer grounds and a row
// that covers the success-criteria answer, and appends that answer's criteria
// as the measure rows.

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/tsv"
)

// coversQ is the form of a question's ID, as the records' qForm gives it.
var coversQ = regexp.MustCompile(`^Q-(00[1-9]|0[1-9][0-9]|[1-9][0-9]{2,})$`)

// factOfQuote gives the fact of an answer's quote: the fact of the span that
// holds the quote's first byte (docs/spec/loop.md, Decided here), or "" when
// the quote's start lies in no fact's span.
func factOfQuote(start, end string, spans [][]string) string {
	s, err := strconv.Atoi(start)
	if err != nil {
		return ""
	}
	for _, r := range spans {
		if r[2] == "" {
			continue
		}
		lo, _ := strconv.Atoi(r[0])
		hi, _ := strconv.Atoi(r[1])
		if s >= lo && s < hi {
			return r[2]
		}
	}
	return ""
}

// checkRevision is the check of the step revision on the session's rows
// (docs/spec/loop.md, The revision of the draft): the rules of the draft
// check, widened so a covers token may be the ID of a fact of the answers — a
// Q-NNN with a row of intake/answers.tsv; a row that covers an answer's fact
// also names the fact of the answer's quote; a fact whose set of covering
// kinds changed from the draft before it is covered by a row that names a
// fact of the answers; and no row covers the question of the intent field
// success-criteria, whose rows the code appends. rows are the session's, as
// draftFile reads them: kind, statement, covers, criterion; before are the
// draft's rows, as records.ReadDraft gives them.
func checkRevision(rows, spans, questions, answers, before [][]string, criteriaQ string) error {
	class := map[string]string{}
	var facts []string
	for _, s := range spans {
		if s[2] != "" {
			class[s[2]] = s[3]
			facts = append(facts, s[2])
		}
	}
	answered := map[string]bool{}
	for _, a := range answers {
		answered[a[0]] = true
	}
	quoteFact := map[string]string{} // the fact of each question's quote
	for _, q := range questions {
		if q[6] != "" {
			quoteFact[q[0]] = factOfQuote(q[6], q[7], spans)
		}
	}
	covered, criterion := map[string]bool{}, map[string]bool{}
	kinds := map[string]map[string]bool{} // the kinds of the rows that cover each fact
	for i, r := range rows {
		covers := strings.Fields(r[2])
		if len(covers) == 0 {
			return fmt.Errorf("row %d: covers names no fact", i+1)
		}
		for _, n := range covers {
			switch {
			case coversQ.MatchString(n):
				if !answered[n] {
					return fmt.Errorf("row %d: covers names %s, which is no fact of the answers", i+1, n)
				}
				if n == criteriaQ {
					return fmt.Errorf("row %d: covers names %s, the answer of the success criteria, whose rows the code appends", i+1, n)
				}
				if f := quoteFact[n]; f != "" && !slices.Contains(covers, f) {
					return fmt.Errorf("row %d: covers names %s, a fact of the answers, but not the fact %s of its quote", i+1, n, f)
				}
			case class[n] == "":
				return fmt.Errorf("row %d: covers names %s, which is no fact of spec/spans.tsv and no fact of the answers", i+1, n)
			default:
				covered[n] = true
				if kinds[n] == nil {
					kinds[n] = map[string]bool{}
				}
				kinds[n][r[0]] = true
				if r[0] == "success-criterion" {
					criterion[n] = true
				}
			}
		}
	}
	old := map[string]map[string]bool{}
	for _, r := range before {
		for _, n := range strings.Fields(r[3]) {
			if class[n] == "" {
				continue
			}
			if old[n] == nil {
				old[n] = map[string]bool{}
			}
			old[n][r[1]] = true
		}
	}
	for _, n := range facts {
		switch c := class[n]; {
		case (c == "need" || c == "constraint") && !covered[n]:
			return fmt.Errorf("the %s fact %s is covered by no row", c, n)
		case c == "measure" && !criterion[n]:
			return fmt.Errorf("the measure fact %s has no success criterion", n)
		case len(old[n]) > 0 && len(kinds[n]) > 0 && !maps.Equal(old[n], kinds[n]):
			cited := false
			for _, r := range rows {
				if !slices.Contains(strings.Fields(r[2]), n) {
					continue
				}
				if slices.ContainsFunc(strings.Fields(r[2]), func(t string) bool { return answered[t] }) {
					cited = true
				}
			}
			if !cited {
				return fmt.Errorf("the class of the fact %s changed, and no row that covers it names a fact of the answers", n)
			}
		}
	}
	return nil
}

// revisionKeep gives the Keep of the step revision on the span table, the
// questions, the answers and the draft before it, the ID of the
// success-criteria question and its criteria: the session's file read by its
// form and checked; the rows with the measure rows of the criteria appended —
// one success-criterion row each, covering the criteria's question, the
// criterion the idea owner's words — R-NNN in their order and the session, as
// spec/draft.tsv, by the block's reader.
func revisionKeep(spans, questions, answers, before [][]string, criteriaQ string, criteria []string) func(id string, data []byte) (map[string][]byte, error) {
	return func(id string, data []byte) (map[string][]byte, error) {
		rows, err := tsv.Read(data, draftFile)
		if err != nil {
			return nil, err
		}
		if err := checkRevision(rows, spans, questions, answers, before, criteriaQ); err != nil {
			return nil, err
		}
		var out [][]string
		for _, r := range rows {
			out = append(out, append(append([]string{fmt.Sprintf("R-%03d", len(out)+1)}, r...), id))
		}
		for _, c := range criteria {
			out = append(out, []string{fmt.Sprintf("R-%03d", len(out)+1), "success-criterion", "The delivery's success criteria.", criteriaQ, c, id})
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

// revisionPrompt is the fixed text of the prompt of the revision, with
// {input}, {spans}, {draft}, {questions}, {answers}, {output} and {finding},
// each path or note on a line of its own.
const revisionPrompt = `You revise the requirement draft of a problem statement, for LAYUP, as its
product owner, after the idea owner's answers to the questions of the batch.
Change no file in this repository and make no commit.

The problem statement: {input}
The span table: {spans}
The draft before the revision: {draft}
The questions: {questions}
The answers: {answers}
The output file: {output}
{finding}
Each table is tab-separated values. The span table is the facts of the problem
statement: a row's fact is its number, class its class. The draft's rows are
the requirement, the non-functional and the success-criterion kinds; a row's
covers names the numbers of the span table's facts, or the ID Q-NNN of an
answer's question when the row rests on the answer. A row that covers an
answer's Q-NNN also names the fact of the answer's quote (the fact of the
span that holds the question's quote). A fact whose covering kinds you change
rests on an answer: one row that covers it names its Q-NNN. The success
criteria of the intent answer are none of yours — code writes their rows.
Write the output file as tab-separated values, with the header line
"kind<TAB>statement<TAB>covers<TAB>criterion" and one row per row of the
revised draft. Then end.
`

// revisionText gives the prompt of the step revision on the inputs and the
// output path; finding is the check's text of the last refused revision (the
// session gets the finding in a new attempt), "" for the first.
func revisionText(finding string) func(input map[string]string, output string) []byte {
	return func(input map[string]string, output string) []byte {
		note := ""
		if finding != "" {
			note = "The check refused the last revision: " + finding + "\nRevise it again on the answers.\n"
		}
		return []byte(strings.NewReplacer("{input}", input["problem-statement.md"],
			"{spans}", input["spans.tsv"], "{draft}", input["draft.tsv"],
			"{questions}", input["questions.tsv"], "{answers}", input["answers.tsv"],
			"{output}", output, "{finding}", note).Replace(revisionPrompt))
	}
}

// revisionStep is the step revision of Intake: the draft's rows whose session
// started on a records base that holds intake/answers.tsv are the revised
// draft — "done before", with the session's comment reposted from the record
// when the issue lacks it (a stop between the commit and the post; the
// batch's precedent); else one session of product-owner on the problem
// statement, the span table, the draft, the questions and the answers of the
// records, each attempt after a refusal run with the finding in its prompt,
// until a kept draft or the step fails.
func (r *state) revisionStep(ctx context.Context) Step {
	const name = "revision"
	data, err := r.store.ReadFile(ctx, "spec/draft.tsv")
	if err != nil {
		return fail(name, err)
	}
	if data == nil {
		return failf(name, "the records have no spec/draft.tsv")
	}
	before, err := records.ReadDraft(data)
	if err != nil {
		return fail(name, fmt.Errorf("spec/draft.tsv: %w", err))
	}
	sessions, err := r.readRecord(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return fail(name, err)
	}
	baseOf := map[string]string{}
	for _, s := range sessions {
		baseOf[s[0]] = s[8]
	}
	revised := true
	for _, row := range before {
		base := baseOf[row[5]]
		if base == "" {
			return failf(name, "sessions.tsv has no session %s of the draft", row[5])
		}
		entries, err := git.LsTree(r.target, base, "intake/answers.tsv")
		if err != nil {
			return fail(name, err)
		}
		if len(entries) == 0 {
			revised = false
			break
		}
	}
	if revised {
		control, err := strconv.Atoi(r.value("issue.control"))
		if err != nil {
			return failf(name, "start.tsv has no number of the control issue")
		}
		comments, err := r.cfg.Forge.Comments(ctx, control)
		if err != nil {
			return fail(name, err)
		}
		line := before[0][5] + ": product-owner session of revision: "
		for _, c := range comments {
			if c.App == r.cfg.Register.AppSlug && strings.HasPrefix(c.Body, line) {
				return done(name, "done before")
			}
		}
		if _, err := r.cfg.Forge.Comment(ctx, control, line+"done\n"); err != nil {
			return fail(name, err)
		}
		return done(name, "the comment of the session posted from the record")
	}
	inputs := map[string][]byte{}
	for _, path := range []string{"start/problem-statement.md", "spec/spans.tsv", "spec/draft.tsv", "questions.tsv", "intake/answers.tsv"} {
		data, err := r.store.ReadFile(ctx, path)
		if err != nil {
			return fail(name, err)
		}
		if data == nil {
			return failf(name, "the records have no %s", path)
		}
		inputs[path[strings.LastIndex(path, "/")+1:]] = data
	}
	spans, err := records.ReadSpans(inputs["spans.tsv"])
	if err != nil {
		return fail(name, fmt.Errorf("spec/spans.tsv: %w", err))
	}
	questions, err := records.ReadQuestions(inputs["questions.tsv"])
	if err != nil {
		return fail(name, fmt.Errorf("questions.tsv: %w", err))
	}
	answers, err := records.ReadAnswers(inputs["answers.tsv"])
	if err != nil {
		return fail(name, fmt.Errorf("intake/answers.tsv: %w", err))
	}
	criteriaQ, criteria := criteriaAnswer(questions, answers)
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
	finding := ""
	for attempt := 1; ; attempt++ {
		r.cfg.Progress(fmt.Sprintf("revision: attempt %d", attempt))
		_, err := s.PhaseSession(ctx, PhaseSpec{Step: name, Role: "product-owner", Harness: h, Model: m, Models: r.cfg.Models,
			Base: base, Records: r.store.Base(), Input: inputs, Output: "revision.tsv",
			Prompt: revisionText(finding), Keep: revisionKeep(spans, questions, answers, before, criteriaQ, criteria)})
		var pe *PhaseError
		if errors.As(err, &pe) && strings.HasPrefix(pe.Reason, "refused file: ") {
			finding = strings.TrimPrefix(pe.Reason, "refused file: ")
			continue
		}
		if err != nil {
			return fail(name, err)
		}
		break
	}
	kept, err := r.store.ReadFile(ctx, "spec/draft.tsv")
	if err != nil {
		return fail(name, err)
	}
	rows, err := records.ReadDraft(kept)
	if err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("the draft revised, %d rows of which %d the success criteria, in spec/draft.tsv", len(rows), len(criteria)))
}
