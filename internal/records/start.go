package records

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// The records of Start of milestone M2a (docs/spec/records.md, NFR-001 — The
// records of Start; task T-8kqn, #126): the schemas of start/start.tsv,
// approvers.tsv, lease.tsv and copies.tsv, and the rules that their blocks
// give in words. internal/run writes and reads them (docs/spec/run.md).

// StartSchema is the form of start/start.tsv: the block start.
var StartSchema = tsv.Schema{Name: "start", Location: "records:start/start.tsv", Columns: []tsv.Column{
	{Name: "name", Type: "text", Key: true},
	{Name: "value", Type: "text"},
	{Name: "source", Type: "enum(command|register|forge|run)"},
}}

// ApproversSchema is the form of approvers.tsv: the block approvers.
var ApproversSchema = tsv.Schema{Name: "approvers", Location: "records:approvers.tsv", Columns: []tsv.Column{
	{Name: "id", Type: "int", Key: true},
	{Name: "role", Type: "enum(operator|idea-owner|approver)", Key: true},
	{Name: "login", Type: "text"},
	{Name: "since", Type: "time"},
	{Name: "source", Type: "text"},
}}

// LeaseSchema is the form of lease.tsv: the block lease.
var LeaseSchema = tsv.Schema{Name: "lease", Location: "records:lease.tsv", Columns: []tsv.Column{
	{Name: "run", Type: "text", Key: true},
	{Name: "host", Type: "text"},
	{Name: "version", Type: "text"},
	{Name: "started", Type: "time"},
	{Name: "heartbeat", Type: "int"},
	{Name: "state", Type: "enum(held|released)"},
}}

// CopiesSchema is the form of copies.tsv: the block copies.
var CopiesSchema = tsv.Schema{Name: "copies", Location: "records:copies.tsv", Columns: []tsv.Column{
	{Name: "comment", Type: "int", Key: true},
	{Name: "seen", Type: "int", Key: true},
	{Name: "issue", Type: "int"},
	{Name: "author_id", Type: "int"},
	{Name: "author_login", Type: "text"},
	{Name: "app", Type: "text"},
	{Name: "created", Type: "time"},
	{Name: "copied", Type: "time"},
	{Name: "sha256", Type: "sha256"},
	{Name: "body", Type: "path"},
}}

// The forms of the types of docs/spec/README.md that the rule of value of
// start.tsv names; internal/tsv does not export its checks.
var (
	intForm     = regexp.MustCompile(`^(0|[1-9][0-9]*)$`)
	decimalForm = regexp.MustCompile(`^(0|[1-9][0-9]*)\.[0-9]+$`)
	sha1Form    = regexp.MustCompile(`^[0-9a-f]{40}$`)
	sha256Form  = regexp.MustCompile(`^[0-9a-f]{64}$`)
	wordForm    = regexp.MustCompile(`^[a-z-]+$`)      // the pattern <word> of the type id
	runForm     = regexp.MustCompile(`^[0-9a-f]{16}$`) // a run ID of lease.tsv
)

func isTime(v string) bool { _, ok := seconds(v, v); return ok }

// A startName is the rule of one name of start.tsv: the form of its value,
// whether the value may be the empty value, and its source.
type startName struct {
	value  func(string) bool
	empty  bool
	source string
}

func anyValue(string) bool { return true }

// startOrder gives the fixed names in the order of the block; "harness" stands
// for the two rows of each harness of the register.
var startOrder = []string{"layup.version", "psb.sha256", "vision.sha256", "forge.plan", "forge.visibility",
	"app.permissions", "operator.id", "idea-owner.id", "intake.cap", "lease.H", "watch.T", "harness",
	"pin.source", "pin.commit", "pin.tree", "pin.time", "issue.intake", "issue.control", "watch"}

// startNames gives the rule of each fixed name, and of harness.cap and
// harness.wall for the rows of each harness.
var startNames = map[string]startName{
	"layup.version":    {anyValue, false, "run"},
	"psb.sha256":       {sha256Form.MatchString, false, "command"},
	"vision.sha256":    {sha256Form.MatchString, true, "command"},
	"forge.plan":       {anyValue, false, "command"},
	"forge.visibility": {anyValue, false, "forge"},
	"app.permissions":  {permissions, false, "forge"},
	"operator.id":      {intForm.MatchString, false, "forge"},
	"idea-owner.id":    {intForm.MatchString, false, "forge"},
	"intake.cap":       {intakeCap, false, "command"},
	"lease.H":          {intForm.MatchString, false, "command"},
	"watch.T":          {intForm.MatchString, false, "command"},
	"harness.cap":      {decimalForm.MatchString, true, "register"},
	"harness.wall":     {intForm.MatchString, false, "register"},
	"pin.source":       {anyValue, false, "run"},
	"pin.commit":       {sha1Form.MatchString, false, "run"},
	"pin.tree":         {sha1Form.MatchString, false, "run"},
	"pin.time":         {isTime, false, "run"},
	"issue.intake":     {issue, true, "forge"},
	"issue.control":    {issue, true, "forge"},
	"watch":            {func(v string) bool { return v == "confirmed" || v == "not-confirmed" }, true, "run"},
}

// permissions: name:level pairs in the form of list(text).
func permissions(v string) bool {
	for _, p := range strings.Split(v, " ") {
		name, level, ok := strings.Cut(p, ":")
		if !ok || name == "" || level == "" {
			return false
		}
	}
	return true
}

// intakeCap: MONEY,HOURS, two values of the type decimal joined by a comma.
func intakeCap(v string) bool {
	money, hours, ok := strings.Cut(v, ",")
	return ok && decimalForm.MatchString(money) && decimalForm.MatchString(hours)
}

// issue: an issue number or opening.
func issue(v string) bool { return v == "opening" || intForm.MatchString(v) }

// CheckStart checks the rules of the block start on the rows of start.tsv, as
// tsv.Read gives them, for the harness register whose IDs are harnesses, in
// its order: each fixed name once, in the order of the block; the two rows of
// each harness together, cap first, in the order of the register; the form
// of each value, the empty value only where the block allows it, and the
// source of each name. A broken rule is a *tsv.Error.
func CheckStart(rows [][]string, harnesses []string) error {
	var want []string
	for _, n := range startOrder {
		if n != "harness" {
			want = append(want, n)
			continue
		}
		for _, h := range harnesses {
			want = append(want, "harness."+h+".cap", "harness."+h+".wall")
		}
	}
	at := func(i int, column, reason string) error {
		return &tsv.Error{Line: i + 2, Column: column, Reason: reason}
	}
	for i, r := range rows {
		f := fields(StartSchema, r)
		name := f["name"]
		if i >= len(want) || name != want[i] {
			next := "no more row"
			if i < len(want) {
				next = want[i]
			}
			return at(i, "name", fmt.Sprintf("the next name is %s (the names of the block, in its order, with the harnesses of the register)", next))
		}
		rule := startNames[name]
		if strings.HasPrefix(name, "harness.") {
			rule = startNames["harness."+name[strings.LastIndex(name, ".")+1:]]
		}
		switch {
		case f["value"] == "" && !rule.empty:
			return at(i, "value", name+" never holds the empty value")
		case f["value"] != "" && !rule.value(f["value"]):
			return at(i, "value", "not a value of the form that the block gives "+name)
		case f["source"] != rule.source:
			return at(i, "source", "the source of "+name+" is "+rule.source)
		}
	}
	if len(rows) < len(want) {
		return &tsv.Error{Line: len(rows) + 2, Column: "name", Reason: "the next name is " + want[len(rows)]}
	}
	return nil
}

// ReadStart reads start.tsv by its schema and its rules, for the harness
// register whose IDs are harnesses; each ID is of the form <word>.
func ReadStart(data []byte, harnesses []string) ([][]string, error) {
	for _, h := range harnesses {
		if !wordForm.MatchString(h) {
			return nil, fmt.Errorf("the harness ID %q is not of the form <word>", h)
		}
	}
	rows, err := tsv.Read(data, StartSchema)
	if err != nil {
		return nil, err
	}
	if err := CheckStart(rows, harnesses); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckApprover checks the rules of the block approvers on one row: login and
// since never hold the empty value, and source is start or a comment ID.
func CheckApprover(r []string) error {
	f := fields(ApproversSchema, r)
	if c := missing(f, []string{"login", "since", "source"}); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	if f["source"] != "start" && !intForm.MatchString(f["source"]) {
		return &RowError{"source", "the source is start or the ID of a comment"}
	}
	return nil
}

// ReadApprovers reads approvers.tsv by its schema and its rules.
func ReadApprovers(data []byte) ([][]string, error) {
	return read(data, ApproversSchema, CheckApprover)
}

// CheckLease checks the rules of the block lease on one row: no column holds
// the empty value, and run is 16 lowercase hexadecimal characters.
func CheckLease(r []string) error {
	f := fields(LeaseSchema, r)
	if c := missing(f, []string{"host", "version", "started", "heartbeat", "state"}); c != "" {
		return &RowError{c, "no column of lease.tsv holds the empty value"}
	}
	if !runForm.MatchString(f["run"]) {
		return &RowError{"run", "a run ID is 16 lowercase hexadecimal characters"}
	}
	return nil
}

// ReadLease reads lease.tsv by its schema and its rules: exactly one row.
func ReadLease(data []byte) ([][]string, error) {
	rows, err := read(data, LeaseSchema, CheckLease)
	if err != nil {
		return nil, err
	}
	if len(rows) != 1 {
		return nil, &tsv.Error{Line: 1, Reason: fmt.Sprintf("the lease table has exactly one row, not %d", len(rows))}
	}
	return rows, nil
}

// CheckCopy checks the rules of the block copies on one row: no column but
// app holds the empty value, and body is copies/<comment>-<seen>.md.
func CheckCopy(r []string) error {
	f := fields(CopiesSchema, r)
	if c := missing(f, []string{"issue", "author_id", "author_login", "created", "copied", "sha256", "body"}); c != "" {
		return &RowError{c, "no column of copies.tsv but app holds the empty value"}
	}
	if want := "copies/" + f["comment"] + "-" + f["seen"] + ".md"; f["body"] != want {
		return &RowError{"body", "the body of this copy is " + want}
	}
	return nil
}

// CheckCopies checks the order of copies.tsv: for each comment, seen runs 1,
// 2, ... with no gap, in the order of the file. A broken rule is a *tsv.Error.
func CheckCopies(rows [][]string) error {
	last := map[string]int{}
	for i, r := range rows {
		f := fields(CopiesSchema, r)
		seen, _ := strconv.Atoi(f["seen"])
		if want := last[f["comment"]] + 1; seen != want {
			return &tsv.Error{Line: i + 2, Column: "seen", Reason: fmt.Sprintf("the next copy of comment %s is seen %d", f["comment"], want)}
		}
		last[f["comment"]] = seen
	}
	return nil
}

// ReadCopies reads copies.tsv by its schema, the rules of each row and the
// rule of the order.
func ReadCopies(data []byte) ([][]string, error) {
	rows, err := read(data, CopiesSchema, CheckCopy)
	if err != nil {
		return nil, err
	}
	if err := CheckCopies(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckValue checks a value of a fixed name of start.tsv by the form that the
// block gives it, so a command can check a flag's value before a run makes
// the row (task T-mqty). The empty value is checked by CheckStart.
func CheckValue(name, value string) error {
	rule, ok := startNames[name]
	if !ok {
		return fmt.Errorf("%s is not a name of start.tsv", name)
	}
	if !rule.value(value) {
		return fmt.Errorf("%q is not a value of the form that the block start gives %s", value, name)
	}
	return nil
}
