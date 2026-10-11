package run

// This file holds the step batch of Intake (docs/spec/intake.md, The merge and
// the IDs; The comment of the batch; The steps of Intake; row 46 of the plan,
// task T-8r3u, #222): the questions of the rule gaps, of meaning, of the
// completeness review and the intent form, merged into questions.tsv and
// posted as one comment on the Intake issue.

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/psb"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// psbGapsSchema is the form of intake/psb.tsv, the table of layup psb check.
var psbGapsSchema = psb.GapsSchema

// batchHeader is the first line of the comment of the batch, by which a
// restart knows that the batch was posted (decided here, task T-8r3u).
const batchHeader = "LAYUP Intake: the batch of questions"

// The fields of the intent form, in their order, and their fixed questions
// (architecture.md §12).
var (
	intentFields    = []string{"band.B", "band.U", "appetite.hours"}
	intentQuestions = map[string]string{
		"band.B":         "B, in USD: below it the delivery goes on without asking you. A decimal, below U.",
		"band.U":         "U, in USD: the upper edge of the band, the money appetite of the delivery. A decimal, above B.",
		"appetite.hours": "The wall-clock appetite of the delivery, in hours. A decimal.",
	}
)

// batchRecords are the records that the batch merges: the problem statement;
// intake/psb.tsv by psbGapsSchema (id, rule, line, excerpt, question); and
// intake/meaning.tsv, intake/coverage.tsv and spec/spans.tsv, as their
// readers give them.
type batchRecords struct {
	PSB                            []byte
	Gaps, Meaning, Coverage, Spans [][]string
}

// lineSpan gives the offsets of line n (from 1) of psb, trimmed of its line
// break and of a CR before it; ok is false for a line that the file lacks.
func lineSpan(psb []byte, n int) (int, int, bool) {
	start := 0
	for i := 1; i < n; i++ {
		j := bytes.IndexByte(psb[start:], '\n')
		if j < 0 {
			return 0, 0, false
		}
		start += j + 1
	}
	if start > len(psb) || (start == len(psb) && n > 1) {
		return 0, 0, false
	}
	end := len(psb)
	if j := bytes.IndexByte(psb[start:], '\n'); j >= 0 {
		end = start + j
	}
	if end > start && psb[end-1] == '\r' {
		end--
	}
	return start, end, true
}

// mergeQuestions gives the rows of questions.tsv of the batch: the rule gaps,
// the gaps of meaning, the items of the completeness review and the intent
// form, in that order. A question whose quote (its offsets) and kind equal
// those of a question before it is dropped; a question with no quote never is
// (decided here, task T-8r3u: the three fields of the intent form share one
// kind and have no quote). Each kept question gets Q-NNN in the order of the
// merge; asked is the time the step read from the run's clock before it wrote
// the rows (decided here, task T-8r3u).
func mergeQuestions(b batchRecords, asked string) ([][]string, error) {
	type question struct{ source, ref, kind, start, end, text, asker string }
	var qs []question
	for _, g := range b.Gaps {
		q := question{"rule", g[0], g[1], "", "", g[4], "layup psb check"}
		if line, _ := strconv.Atoi(g[2]); line > 0 {
			start, end, ok := lineSpan(b.PSB, line)
			if !ok {
				return nil, fmt.Errorf("intake/psb.tsv: %s names line %d, which the problem statement lacks", g[0], line)
			}
			q.start, q.end = strconv.Itoa(start), strconv.Itoa(end)
		}
		qs = append(qs, q)
	}
	for _, m := range b.Meaning {
		qs = append(qs, question{"meaning", m[0], m[1], m[2], m[3], m[5], m[6]})
	}
	span := map[string][2]string{} // the offsets of each fact
	for _, s := range b.Spans {
		if s[2] != "" {
			span[s[2]] = [2]string{s[0], s[1]}
		}
	}
	for _, c := range b.Coverage {
		q := question{"coverage", c[0], c[1], "", "", c[4], c[5]}
		if c[1] != "unsupported" {
			sp, ok := span[c[2]]
			if !ok {
				return nil, fmt.Errorf("intake/coverage.tsv: %s names the fact %s, which spec/spans.tsv lacks", c[0], c[2])
			}
			q.start, q.end = sp[0], sp[1]
		}
		qs = append(qs, q)
	}
	for _, f := range intentFields {
		qs = append(qs, question{"intent", f, "intent", "", "", intentQuestions[f], "layup run"})
	}
	var rows [][]string
	seen := map[string]bool{}
	for _, q := range qs {
		if q.start != "" {
			key := q.start + " " + q.end + " " + q.kind
			if seen[key] {
				continue
			}
			seen[key] = true
		}
		rows = append(rows, []string{fmt.Sprintf("Q-%03d", len(rows)+1), q.source, q.ref, "idea-owner", "before-delivery", q.kind,
			q.start, q.end, q.text, asked, q.asker, "", "asked"})
	}
	return rows, nil
}

// backquote writes s as Markdown code: in a run of backquotes one longer than
// the longest run in s, with a space inside when s begins or ends with a
// backquote (decided here, task T-8r3u, as CommonMark reads it).
func backquote(s string) string {
	longest, run := 0, 0
	for _, c := range s {
		if c == '`' {
			run++
			longest = max(longest, run)
		} else {
			run = 0
		}
	}
	fence := strings.Repeat("`", longest+1)
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		return fence + " " + s + " " + fence
	}
	return fence + s + fence
}

// batchComment gives the comment of the batch of the rows of questions.tsv
// (as records.ReadQuestions gives them) on the problem statement psb: the
// header line; the block for the idea owner, each question as "Q-NNN:
// <question>" and its quote on the next line; the block for the Operator,
// which has no question in M2c; the rule of the answers.
func batchComment(rows [][]string, psb []byte) string {
	return questionsComment(batchHeader, rows, func(r []string) (string, string) { return r[6], r[7] }, psb)
}

// questionsComment gives a comment in the form of the batch: header, the block
// for the idea owner, each row's question and the quote at the offsets that
// quote gives for it (none for ""), the block for the Operator and the rule of
// the answers; the batch and the follow-up share it (row 47b).
func questionsComment(header string, rows [][]string, quote func([]string) (string, string), psb []byte) string {
	var b strings.Builder
	b.WriteString(header + "\n\n## For the idea owner\n\n")
	for _, r := range rows {
		fmt.Fprintf(&b, "%s: %s\n", r[0], r[8])
		if s, e := quote(r); s != "" {
			start, _ := strconv.Atoi(s)
			end, _ := strconv.Atoi(e)
			b.WriteString(backquote(string(psb[start:end])) + "\n")
		}
		b.WriteString("\n")
	}
	b.WriteString("## For the Operator\n\nNo setup question in M2c.\n\n")
	b.WriteString("Answer each block in one comment, one line per question ID: `Q-NNN: <answer>`.\n")
	return b.String()
}

// readRecord reads the record path by read; a missing record is an error.
func (r *state) readRecord(ctx context.Context, path string, read func([]byte) ([][]string, error)) ([][]string, error) {
	data, err := r.store.ReadFile(ctx, path)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("the records have no %s", path)
	}
	rows, err := read(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return rows, nil
}

// batchPosted reports whether the Intake issue holds a comment of the App that
// starts with the batch's header line. The App's own comment is read with no
// copy: it is no human input, and no decision rests on it (decided here, task
// T-8r3u; run.md, Copy before read).
func (r *state) batchPosted(ctx context.Context, issue int) (bool, error) {
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return false, err
	}
	for _, c := range comments {
		if c.App == r.cfg.Register.AppSlug && strings.HasPrefix(c.Body, batchHeader+"\n") {
			return true, nil
		}
	}
	return false, nil
}

// batchStep is the step batch of Intake: with no questions.tsv, the merge, one
// records commit of questions.tsv, then the comment on the Intake issue (the
// commit announces it, by fencing); with questions.tsv and the comment, "done
// before"; with questions.tsv and no comment (a run that stopped between the
// commit and the post), the comment posted from the record, with no second
// commit (decided here, task T-8r3u).
func (r *state) batchStep(ctx context.Context) Step {
	const name = "batch"
	issue, err := strconv.Atoi(r.value("issue.intake"))
	if err != nil {
		return failf(name, "start.tsv has no number of the Intake issue")
	}
	psbData, err := r.store.ReadFile(ctx, "start/problem-statement.md")
	if err != nil {
		return fail(name, err)
	}
	if psbData == nil {
		return failf(name, "the records have no start/problem-statement.md")
	}
	held, err := r.store.ReadFile(ctx, "questions.tsv")
	if err != nil {
		return fail(name, err)
	}
	if held != nil {
		posted, err := r.batchPosted(ctx, issue)
		if err != nil {
			return fail(name, err)
		}
		if posted {
			return done(name, "done before")
		}
		rows, err := records.ReadQuestions(held)
		if err != nil {
			return fail(name, fmt.Errorf("questions.tsv: %w", err))
		}
		if _, err := r.cfg.Forge.Comment(ctx, issue, batchComment(rows, psbData)); err != nil {
			return fail(name, err)
		}
		return done(name, "the comment of the batch posted from questions.tsv")
	}
	b := batchRecords{PSB: psbData}
	for _, rec := range []struct {
		path string
		read func([]byte) ([][]string, error)
		into *[][]string
	}{
		{"intake/psb.tsv", func(d []byte) ([][]string, error) { return tsv.Read(d, psbGapsSchema) }, &b.Gaps},
		{"intake/meaning.tsv", records.ReadMeaning, &b.Meaning},
		{"intake/coverage.tsv", records.ReadCoverage, &b.Coverage},
		{"spec/spans.tsv", records.ReadSpans, &b.Spans},
	} {
		if *rec.into, err = r.readRecord(ctx, rec.path, rec.read); err != nil {
			return fail(name, err)
		}
	}
	rows, err := mergeQuestions(b, stamp(r.cfg.Clock.Now()))
	if err != nil {
		return fail(name, err)
	}
	data, err := table(records.QuestionsSchema, rows)
	if err != nil {
		return fail(name, err)
	}
	if _, err := records.ReadQuestions(data); err != nil {
		return fail(name, err)
	}
	if err := r.write(ctx, map[string][]byte{"questions.tsv": data}, "layup run: the batch of questions"); err != nil {
		return fail(name, err)
	}
	if _, err := r.cfg.Forge.Comment(ctx, issue, batchComment(rows, psbData)); err != nil {
		return fail(name, err)
	}
	return done(name, fmt.Sprintf("%d questions, in questions.tsv and the comment of the batch", len(rows)))
}
