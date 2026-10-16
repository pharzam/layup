package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// parametersPath is the file of the parameters in the records.
const parametersPath = "parameters.tsv"

// A Parameter is one row of parameters.tsv (docs/spec/records.md, the block
// parameters; task T-dx1h, #358, row 66b of the plan of M2e).
type Parameter struct {
	Name, Value, Default, Bound, Role, Evidence string // Value, Default and Evidence are "" for —
	Comment                                     int64  // 0 is —
	Since                                       time.Time
}

// ParametersSchema is the Go value of the block parameters of docs/spec/records.md;
// its owner is this package, and its test compares it with the block.
var ParametersSchema = tsv.Schema{Name: "parameters", Location: "records:parameters.tsv", Columns: []tsv.Column{
	{Name: "name", Type: "text", Key: true},
	{Name: "value", Type: "text"},
	{Name: "default", Type: "text"},
	{Name: "bound", Type: "text"},
	{Name: "role", Type: "enum(operator|idea-owner|approver)"},
	{Name: "evidence", Type: "text"},
	{Name: "comment", Type: "int"},
	{Name: "since", Type: "time"},
}}

// A paramSpec is one row of the table of the parameters in code: the table of
// docs/spec/loop.md, The parameters of M2e. Type and Bound are the words of that
// table; ok is the check of the type, within the check of the bound (nil for
// none beyond the type); Start marks a row whose first value is start.tsv's.
type paramSpec struct {
	Name, Type, Bound, Role, Default string
	Start                            bool
	ok, within                       func(string) bool
}

// parameterTable is the table of docs/spec/loop.md, REQ-012 — The parameter change, "The parameters of M2e": the six
// rows that a source gives completely, in its order (the plan review of #358, C3). The idea owner has no row: no source
// names the allowed dependencies' row.
var parameterTable = []paramSpec{
	{Name: "brief.lines", Type: "int", Bound: "a whole number of lines, 1 or more", Role: "operator", ok: wholeNumber, within: atLeastOne},
	{Name: "build.parallel", Type: "int", Bound: "a whole number of tasks, 1 or more", Role: "operator", ok: wholeNumber, within: atLeastOne},
	{Name: "transitions", Type: "sha256", Role: "operator", ok: isSHA256,
		Bound: "the SHA-256 of a copied comment of the control issue whose whole body is a table that passes the check of a table: a row for each default pair, only the kinds and items of the lists, and no `schema` in the row of code's own handoff"},
	{Name: "lease.H", Type: "int", Bound: "whole minutes, 1 or more, as the flag `--lease-h`", Role: "operator", Start: true, ok: startForm("lease.H"), within: atLeastOne},
	{Name: "watch.T", Type: "int", Bound: "whole minutes, 1 or more, as the flag `--watch-t`", Role: "operator", Start: true, ok: startForm("watch.T"), within: atLeastOne},
	{Name: "intake.cap", Type: "MONEY,HOURS", Bound: "MONEY,HOURS, each a decimal with a fraction, as the flag `--intake-cap`", Role: "operator", Start: true, ok: startForm("intake.cap")},
}

var (
	wholeNumber = regexp.MustCompile(`^(0|[1-9][0-9]*)$`).MatchString
	isSHA256    = regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString
)

// atLeastOne is the bound of a whole number: it has no leading zero, so it is 1 or more unless it is 0.
func atLeastOne(v string) bool { return v != "0" }

// startForm is the form that the block start gives a name (the flag's own check).
func startForm(name string) func(string) bool {
	return func(v string) bool { return records.CheckValue(name, v) == nil }
}

func (p Parameter) fields() []string {
	comment := ""
	if p.Comment != 0 {
		comment = strconv.FormatInt(p.Comment, 10)
	}
	since := ""
	if !p.Since.IsZero() {
		since = stamp(p.Since)
	}
	return []string{p.Name, p.Value, p.Default, p.Bound, p.Role, p.Evidence, comment, since}
}

// DefaultParameters gives the first write of parameters.tsv: a row for each of
// the table, with no value and no default but for the rows of Start, which take
// their value from start (the values of start.tsv) and say command as their
// evidence (docs/spec/loop.md, The parameter change). A value of start that is
// missing or out of its row's form is an error: the first write has no guess.
func DefaultParameters(table []paramSpec, start func(string) string, now time.Time) ([]Parameter, error) {
	rows := make([]Parameter, len(table))
	for i, spec := range table {
		rows[i] = Parameter{Name: spec.Name, Default: spec.Default, Bound: spec.Bound, Role: spec.Role, Since: now}
		if !spec.Start {
			continue
		}
		v := start(spec.Name)
		if v == "" {
			return nil, fmt.Errorf("start.tsv has no value of %s", spec.Name)
		}
		if err := checkValue(spec, v); err != nil {
			return nil, fmt.Errorf("start.tsv, %s: %v", spec.Name, err)
		}
		rows[i].Value, rows[i].Evidence = v, startEvidence
	}
	return rows, nil
}

// startEvidence is the evidence of a value that the command gave at Start.
const startEvidence = "command"

// checkValue checks the value of a row: its type, then its bound.
func checkValue(spec paramSpec, v string) error {
	if !spec.ok(v) {
		return fmt.Errorf("%q is not of the type %s of %s", v, spec.Type, spec.Name)
	}
	if spec.within != nil && !spec.within(v) {
		return fmt.Errorf("%q is out of the bound of %s: %s", v, spec.Name, spec.Bound)
	}
	return nil
}

// WriteParameters writes the table by ParametersSchema. It writes nothing for a
// table that ReadParameters would refuse.
func WriteParameters(w io.Writer, rows []Parameter) error {
	out := make([][]string, len(rows))
	for i, p := range rows {
		out[i] = p.fields()
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, ParametersSchema, out); err != nil {
		return err
	}
	if _, err := ReadParameters(b.Bytes()); err != nil {
		return err
	}
	_, err := w.Write(b.Bytes())
	return err
}

// ReadParameters reads a table by ParametersSchema and then makes the checks of
// its owner (CheckParameters) against the table of the code.
func ReadParameters(data []byte) ([]Parameter, error) {
	recs, err := tsv.Read(data, ParametersSchema)
	if err != nil {
		return nil, err
	}
	rows := make([]Parameter, len(recs))
	for i, f := range recs {
		p := Parameter{Name: f[0], Value: f[1], Default: f[2], Bound: f[3], Role: f[4], Evidence: f[5]}
		if f[6] != "" {
			if p.Comment, err = strconv.ParseInt(f[6], 10, 64); err != nil {
				return nil, &tsv.Error{Line: i + 2, Column: "comment", Reason: err.Error()}
			}
		}
		if p.Since, err = time.Parse(timeForm, f[7]); err != nil {
			return nil, &tsv.Error{Line: i + 2, Column: "since", Reason: err.Error()}
		}
		rows[i] = p
	}
	if err := CheckParameters(parameterTable, rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckParameters makes the checks that docs/spec/records.md gives the owner of
// parameters, after the types: a row for each name of the table, in its order,
// and no other; its bound words, role and default those of the table, and since
// never —. A row with no comment is its own origin: a Start row has the value
// of start.tsv in the form of its type and the evidence command, and any other
// row has neither a value nor evidence. A row with a comment is a set: a value
// of its type and within its bound, and its reason as evidence (so the word
// command is a reason there). The refusal names the line and the column.
func CheckParameters(table []paramSpec, rows []Parameter) error {
	for i, spec := range table {
		if i >= len(rows) {
			return &tsv.Error{Line: i + 2, Column: "name", Reason: "the next name is " + spec.Name + " (the names of the table, in its order)"}
		}
		p, line := rows[i], i+2
		at := func(column, reason string) error { return &tsv.Error{Line: line, Column: column, Reason: reason} }
		switch {
		case p.Name != spec.Name:
			return at("name", "the next name is "+spec.Name+" (the names of the table, in its order)")
		case p.Bound != spec.Bound:
			return at("bound", "the bound of "+spec.Name+" is the words of the table: "+spec.Bound)
		case p.Role != spec.Role:
			return at("role", "the role of "+spec.Name+" is "+spec.Role)
		case p.Default != spec.Default:
			return at("default", "the default of "+spec.Name+" is "+orDash(spec.Default))
		case p.Since.IsZero():
			return at("since", "never —")
		case p.Comment != 0 && p.Evidence == "":
			return at("evidence", "a set has its reason as evidence, never —")
		case p.Comment != 0 && p.Value == "":
			return at("value", "a set has a value, never —")
		case p.Comment != 0:
			if err := checkValue(spec, p.Value); err != nil {
				return at("value", err.Error())
			}
		case spec.Start:
			if err := checkValue(spec, p.Value); err != nil {
				return at("value", "a value of Start is start.tsv's, "+err.Error())
			}
			if p.Evidence != startEvidence {
				return at("evidence", "the evidence of a value of Start with no comment is "+startEvidence)
			}
		case p.Value != spec.Default:
			return at("value", spec.Name+" has no value until a set gives one (a value with no comment)")
		case p.Evidence != "":
			return at("evidence", "a row with no value and no comment has no evidence")
		}
	}
	if len(rows) > len(table) {
		return &tsv.Error{Line: len(table) + 2, Column: "name", Reason: rows[len(table)].Name + " is not a name of the table"}
	}
	return nil
}

func orDash(s string) string {
	if s == "" {
		return "—"
	}
	return s
}

// A setLine is the three parts of a set: its name, its value and its reason.
type setLine struct{ Name, Value, Reason string }

// parseSet reads the body of a comment as `set <name> <value> because <reason>`
// (docs/spec/loop.md, The parameter change, the form): one line, once its
// trailing white space is trimmed, with one space between the parts, a name and
// a value of no space, and a reason that is not empty and not —. A body whose
// first word is not set is no set (false, no error); one that starts with the
// word and breaks the form is a set that is refused (true, the error).
func parseSet(body string) (setLine, bool, error) {
	text := strings.TrimRight(body, " \t\r\n")
	first := text
	if i := strings.IndexAny(text, " \t\r\n"); i >= 0 {
		first = text[:i]
	}
	if first != "set" {
		return setLine{}, false, nil
	}
	const form = "the form is `set <name> <value> because <reason>`, on one line"
	if strings.ContainsAny(text, "\r\n") {
		return setLine{}, true, fmt.Errorf("%s: this comment has more than one line", form)
	}
	f := strings.Split(text, " ")
	if len(f) < 5 || f[1] == "" || f[2] == "" || f[3] != "because" {
		return setLine{}, true, fmt.Errorf("%s", form)
	}
	reason := strings.Join(f[4:], " ")
	if reason == "" || reason == "—" {
		return setLine{}, true, fmt.Errorf("%s: the reason is not empty and not —", form)
	}
	return setLine{Name: f[1], Value: f[2], Reason: reason}, true, nil
}

// A setInput is the first copy of a comment: its row of copies.tsv and its body.
type setInput struct {
	Row  []string
	Body string
}

// A Refusal is a set that was refused: its comment, the item that refused it
// (form, author, name, value or bound) and the reason.
type Refusal struct {
	Comment      int64
	Item, Reason string
}

// A Change is a set that passed. Table holds the bytes of the new transitions.tsv
// for the name transitions, the bytes of the copy that its value names.
type Change struct {
	Comment             int64
	Name, Value, Reason string
	Table               []byte
}

// A judgement is the result of judgeSet: a change, a refusal, or neither: the
// comment is no set, or it is superseded (its row records a comment that is the
// same or later), which is no change, no refusal and no report.
type judgement struct {
	Change     *Change
	Refusal    *Refusal
	Superseded bool
}

// judgeSet judges the first copy c of a comment against the table, the approvers
// and the rows now in force, in the order of docs/spec/loop.md: the form; the
// author (a first copy, from no App, by an account of approvers.tsv, and, once the
// name is known, in the role of its row); the name; the value; the bound. A set
// whose comment is not later than the comment that its row records is
// superseded, and judged no more (the plan review of #358, C2). For the name
// transitions, the bound is that lookup gives the body of a copy with the value as
// its SHA-256 and that the body reads as a table that passes the check of a table.
func judgeSet(table []paramSpec, approvers [][]string, c setInput, current []Parameter, lookup func(sha string) ([]byte, bool)) judgement {
	line, isSet, err := parseSet(c.Body)
	if !isSet {
		return judgement{}
	}
	id, _ := strconv.ParseInt(c.Row[0], 10, 64)
	refuse := func(item, format string, a ...any) judgement {
		return judgement{Refusal: &Refusal{Comment: id, Item: item, Reason: fmt.Sprintf(format, a...)}}
	}
	if err != nil {
		return refuse("form", "%v", err)
	}
	known := slices.ContainsFunc(approvers, func(a []string) bool { return a[0] == c.Row[3] })
	if c.Row[1] != "1" || c.Row[5] != "" || !known {
		return refuse("author", "a decision is the first copy of a comment by no App, by an account of approvers.tsv")
	}
	i := slices.IndexFunc(table, func(s paramSpec) bool { return s.Name == line.Name })
	if i < 0 {
		return refuse("name", "%s is not a row of the table of parameters", line.Name)
	}
	spec := table[i]
	if !Decision(c.Row, approvers, spec.Role) {
		return refuse("author", "%s is changed by the %s, and this account has not that role", spec.Name, spec.Role)
	}
	if j := slices.IndexFunc(current, func(p Parameter) bool { return p.Name == spec.Name }); j >= 0 && current[j].Comment >= id {
		return judgement{Superseded: true}
	}
	if !spec.ok(line.Value) {
		return refuse("value", "%q is not of the type %s of %s", line.Value, spec.Type, spec.Name)
	}
	if spec.within != nil && !spec.within(line.Value) {
		return refuse("bound", "%q is out of the bound of %s: %s", line.Value, spec.Name, spec.Bound)
	}
	change := Change{Comment: id, Name: spec.Name, Value: line.Value, Reason: line.Reason}
	if spec.Name == "transitions" {
		body, found := lookup(line.Value)
		if !found {
			return refuse("bound", "no copy of a comment of the control issue has the SHA-256 %s", line.Value)
		}
		rows, err := handoff.ReadTable(body)
		if err == nil {
			err = handoff.CheckTable(rows)
		}
		if err != nil {
			return refuse("bound", "the copy is not a table that passes the check of a table: %v", err)
		}
		change.Table = body
	}
	return judgement{Change: &change}
}

// applyChange gives the table with the row of the change replaced: its value, its
// reason as evidence, its comment and the time; the rows are a copy, and the
// default, bound and role are kept.
func applyChange(rows []Parameter, ch Change, now time.Time) ([]Parameter, error) {
	i := slices.IndexFunc(rows, func(p Parameter) bool { return p.Name == ch.Name })
	if i < 0 {
		return nil, fmt.Errorf("%s is not a row of parameters.tsv", ch.Name)
	}
	next := slices.Clone(rows)
	next[i].Value, next[i].Evidence, next[i].Comment, next[i].Since = ch.Value, ch.Reason, ch.Comment, now
	return next, nil
}

// parameterFiles gives the file of the first write of the parameters: parameters.tsv
// from the defaults when the records have none (existing is nil); one that exists
// is read and checked, and kept (no file); one that its reader refuses is an error
// and no file (docs/spec/loop.md, The parameter change: written in the commit that
// writes transitions.tsv, and never from the default again).
func parameterFiles(existing []byte, start func(string) string, now time.Time) (map[string][]byte, error) {
	if existing != nil {
		if _, err := ReadParameters(existing); err != nil {
			return nil, fmt.Errorf("%s: %w", parametersPath, err)
		}
		return map[string][]byte{}, nil
	}
	rows, err := DefaultParameters(parameterTable, start, now)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", parametersPath, err)
	}
	var b bytes.Buffer
	if err := WriteParameters(&b, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", parametersPath, err)
	}
	return map[string][]byte{parametersPath: b.Bytes()}, nil
}

// An Applied is a change that the boundary wrote, with the note for the reader:
// empty, or that a value of Start is recorded and not read.
type Applied struct {
	Change Change
	Note   string
}

// startNote is said of a set of a row of Start (docs/spec/loop.md: nothing of M2e
// reads those rows from parameters.tsv).
const startNote = "recorded; this run reads start.tsv, so M2e still acts on the old value"

// parametersBoundary is what the loop calls between two steps (the call is row
// 66k's, #367): the next step boundary of docs/spec/loop.md, The parameter
// change. It copies the comments of the control issue that are no App's (copy
// before read), judges the first copy of each in the order of the comment IDs,
// and writes each change as one records commit, through the fence of the lease,
// holding parameters.tsv and, for transitions, transitions.tsv too. It gives the
// changes and the refusals, each of which it also reports as a progress line; a
// refused set changes and records nothing, and is judged again at the next
// boundary with the same result.
func (r *state) parametersBoundary(ctx context.Context) ([]Applied, []Refusal, error) {
	issue, err := strconv.Atoi(r.value("issue.control"))
	if err != nil {
		return nil, nil, fmt.Errorf("no control issue")
	}
	data, err := r.store.ReadFile(ctx, parametersPath)
	if err != nil {
		return nil, nil, err
	}
	if data == nil {
		return nil, nil, fmt.Errorf("%s is not in the records: the first task of the loop writes it", parametersPath)
	}
	rows, err := ReadParameters(data)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", parametersPath, err)
	}
	comments, err := r.cfg.Forge.Comments(ctx, issue)
	if err != nil {
		return nil, nil, err
	}
	for _, c := range comments {
		if c.App != "" {
			continue
		}
		if err := r.copy(ctx, c, issue); err != nil {
			return nil, nil, err
		}
	}
	var firsts [][]string // the first copies of the control issue by no App, in the order of the comment IDs
	for _, row := range r.copies {
		if row[1] == "1" && row[2] == strconv.Itoa(issue) && row[5] == "" {
			firsts = append(firsts, row)
		}
	}
	id := func(row []string) int64 { n, _ := strconv.ParseInt(row[0], 10, 64); return n }
	sort.Slice(firsts, func(i, j int) bool { return id(firsts[i]) < id(firsts[j]) })
	var lookupErr error
	lookup := func(sha string) ([]byte, bool) {
		for _, row := range r.copies {
			if row[2] == strconv.Itoa(issue) && row[8] == sha {
				body, err := r.store.ReadFile(ctx, row[9])
				if err != nil || body == nil {
					lookupErr = fmt.Errorf("the copy %s: %v", row[9], err)
					return nil, false
				}
				return body, true
			}
		}
		return nil, false
	}
	var applied []Applied
	var refused []Refusal
	for _, row := range firsts {
		body, err := r.store.ReadFile(ctx, row[9])
		if err != nil || body == nil {
			return applied, refused, fmt.Errorf("the copy %s: %v", row[9], err)
		}
		j := judgeSet(parameterTable, r.approvers, setInput{Row: row, Body: string(body)}, rows, lookup)
		if lookupErr != nil {
			return applied, refused, lookupErr
		}
		switch {
		case j.Refusal != nil:
			refused = append(refused, *j.Refusal)
			r.cfg.Progress(fmt.Sprintf("parameters: refused comment %d: %s: %s", j.Refusal.Comment, j.Refusal.Item, j.Refusal.Reason))
		case j.Change != nil:
			next, err := applyChange(rows, *j.Change, r.cfg.Clock.Now().UTC().Truncate(time.Second))
			if err != nil {
				return applied, refused, err
			}
			var b bytes.Buffer
			if err := WriteParameters(&b, next); err != nil {
				return applied, refused, err
			}
			files := map[string][]byte{parametersPath: b.Bytes()}
			if j.Change.Table != nil {
				files[transitionsPath] = j.Change.Table
			}
			if err := r.write(ctx, files, "parameters: set "+j.Change.Name); err != nil {
				return applied, refused, err
			}
			rows = next
			a := Applied{Change: *j.Change}
			line := fmt.Sprintf("parameters: applied set %s %s (comment %d)", j.Change.Name, j.Change.Value, j.Change.Comment)
			if i := slices.IndexFunc(parameterTable, func(s paramSpec) bool { return s.Name == j.Change.Name }); parameterTable[i].Start {
				a.Note = startNote
				line += ": " + startNote
			}
			applied = append(applied, a)
			r.cfg.Progress(line)
		}
	}
	return applied, refused, nil
}
