package run

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/handoff"
)

// The parameters (task T-dx1h, #358, row 66b of the plan of M2e; docs/spec/loop.md, REQ-012 — The parameter
// change and The transition table; docs/spec/records.md, the block parameters and its owner checks).

var (
	paramsAt = time.Date(2026, 10, 16, 12, 0, 0, 0, time.UTC)
	// The accounts of approvers.tsv: the Operator, the idea owner and an approver of a point.
	paramApprovers = [][]string{{"101", "operator", "pharzam", "2026-10-08T12:00:00Z", "start"}, {"303", "idea-owner", "carol", "2026-10-08T12:00:00Z", "start"},
		{"404", "approver", "dave", "2026-10-08T12:00:00Z", "5001"}}
	startValues = map[string]string{"lease.H": "5", "watch.T": "1", "intake.cap": "50.0,8.0"}
)

func startValue(name string) string { return startValues[name] }

// defaultParams is the table of defaults, made here from the table of the code and not by DefaultParameters, so that
// a test of another function does not rest on it.
func defaultParams(t *testing.T) []Parameter {
	t.Helper()
	var rows []Parameter
	for _, s := range parameterTable {
		r := Parameter{Name: s.Name, Bound: s.Bound, Role: s.Role, Since: paramsAt}
		if s.Start {
			r.Value, r.Evidence = startValue(s.Name), "command"
		}
		rows = append(rows, r)
	}
	if len(rows) != 6 {
		t.Fatalf("the table of the code has %d rows; want the six of the table of loop.md", len(rows))
	}
	return rows
}

// aSet is the first copy of the comment id by the account author, as the boundary reads it.
func aSet(id int64, author, body string) setInput {
	return setInput{Row: []string{strconv.FormatInt(id, 10), "1", "2", author, "login", "", "2026-10-16T12:00:00Z", "2026-10-16T12:00:01Z", strings.Repeat("a", 64), "copies/x.md"}, Body: body}
}

func noCopy(string) ([]byte, bool) { return nil, false }

func TestTheDefaultParametersAreTheSixRowsOfTheTable(t *testing.T) {
	rows, err := DefaultParameters(parameterTable, startValue, paramsAt)
	if err != nil || len(rows) != 6 {
		t.Fatalf("the default rows: %d rows, %v; want the six of the table", len(rows), err)
	}
	var names []string
	for _, r := range rows {
		names = append(names, r.Name)
		if r.Role != "operator" || r.Bound == "" || !r.Since.Equal(paramsAt) {
			t.Errorf("%s: role %q, bound %q, since %s; want operator, words, the time of the write", r.Name, r.Role, r.Bound, r.Since)
		}
		if r.Default != "" || r.Value == "100" {
			t.Errorf("%s: default %q, value %q; want no default (no source gives one)", r.Name, r.Default, r.Value)
		}
	}
	if want := []string{"brief.lines", "build.parallel", "transitions", "lease.H", "watch.T", "intake.cap"}; !reflect.DeepEqual(names, want) {
		t.Errorf("the rows %v; want the six of the table, in its order, %v", names, want)
	}
	for _, spec := range parameterTable {
		if spec.Role == "idea-owner" {
			t.Errorf("%s is a row of the idea owner: no source names that row (the plan review, C3)", spec.Name)
		}
	}
	for _, r := range rows[:3] {
		if r.Value != "" || r.Evidence != "" || r.Comment != 0 {
			t.Errorf("%s is %+v; want unset, with no evidence and no comment", r.Name, r)
		}
	}
}

func TestTheStartParametersTakeTheirValuesFromStartTsv(t *testing.T) {
	rows, err := DefaultParameters(parameterTable, func(n string) string {
		return map[string]string{"lease.H": "9", "watch.T": "3", "intake.cap": "20.0,4.0"}[n]
	}, paramsAt)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]Parameter{}
	for _, r := range rows {
		got[r.Name] = r
	}
	for name, want := range map[string]string{"lease.H": "9", "watch.T": "3", "intake.cap": "20.0,4.0"} {
		if r := got[name]; r.Value != want || r.Evidence != "command" || r.Comment != 0 || r.Default != "" {
			t.Errorf("%s is %+v; want value %q, evidence command, no comment, no default", name, r, want)
		}
	}
	if _, err := DefaultParameters(parameterTable, func(n string) string { return map[string]string{"lease.H": "5", "intake.cap": "50.0,8.0"}[n] }, paramsAt); err == nil || !strings.Contains(err.Error(), "watch.T") {
		t.Errorf("a start.tsv with no watch.T gave %v; want an error that names it", err)
	}
}

func TestParametersAreWrittenAndReadBackWithTheirDashes(t *testing.T) {
	rows := defaultParams(t)
	set := rows[0]
	set.Value, set.Evidence, set.Comment, set.Since = "7", "the idea owner reads seven lines", 6200000001, paramsAt.Add(time.Hour)
	rows[0] = set
	var b bytes.Buffer
	if err := WriteParameters(&b, rows); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSuffix(b.String(), "\n"), "\n")
	if lines[0] != "name\tvalue\tdefault\tbound\trole\tevidence\tcomment\tsince" || len(lines) != 7 {
		t.Fatalf("the table:\n%s", b.String())
	}
	if !strings.HasPrefix(lines[1], "brief.lines\t7\t—\t") || !strings.HasSuffix(lines[1], "\toperator\tthe idea owner reads seven lines\t6200000001\t2026-10-16T13:00:00Z") {
		t.Errorf("the row of a set: %q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "build.parallel\t—\t—\t") || !strings.HasSuffix(lines[2], "\toperator\t—\t—\t2026-10-16T12:00:00Z") {
		t.Errorf("the row of a default: %q", lines[2])
	}
	if !strings.HasPrefix(lines[4], "lease.H\t5\t—\t") || !strings.HasSuffix(lines[4], "\toperator\tcommand\t—\t2026-10-16T12:00:00Z") {
		t.Errorf("the row of a Start value: %q", lines[4])
	}
	got, err := ReadParameters(b.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, rows) {
		t.Errorf("the table read back\n%+v\nwant\n%+v", got, rows)
	}
	bad := slicesClone(rows)
	bad[1].Role = ""
	if err := WriteParameters(&bytes.Buffer{}, bad); err == nil {
		t.Error("a table that the reader refuses was written")
	}
}

func slicesClone(rows []Parameter) []Parameter { return append([]Parameter(nil), rows...) }

// The owner checks of the block parameters (records.md): after the types, bound, role and since are never —; a row
// with no comment is its own origin (a Start row says command, any other row is unset), and a row with a comment is a
// set, whose evidence is its reason (C1 of the plan review).
func TestTheOwnerChecksOfParametersRefuseADashWhereTheBlockAllowsNone(t *testing.T) {
	base := func() []Parameter {
		rows := defaultParams(t)
		rows[0].Value, rows[0].Evidence, rows[0].Comment = "7", "the idea owner reads seven lines", 55
		return rows
	}
	cases := []struct {
		name   string
		edit   func([]Parameter)
		column string // "" for a table that passes
		line   int
	}{
		{"the table as written", func([]Parameter) {}, "", 0},
		{"a set whose reason is the word command", func(r []Parameter) { r[0].Evidence = "command" }, "", 0},
		{"a bound with no words", func(r []Parameter) { r[1].Bound = "" }, "bound", 3},
		{"a role with no word", func(r []Parameter) { r[2].Role = "" }, "role", 4},
		{"a role that is not the row's", func(r []Parameter) { r[2].Role = "idea-owner" }, "role", 4},
		{"a row with no time", func(r []Parameter) { r[1].Since = time.Time{} }, "since", 3},
		{"a set with no value", func(r []Parameter) { r[0].Value = "" }, "value", 2},
		{"a set with no reason", func(r []Parameter) { r[0].Evidence = "" }, "evidence", 2},
		{"a value with no comment on a row that has no origin", func(r []Parameter) { r[1].Value = "3" }, "value", 3},
		{"evidence with no comment on such a row", func(r []Parameter) { r[1].Evidence = "because" }, "evidence", 3},
		{"a Start row with no evidence", func(r []Parameter) { r[3].Evidence = "" }, "evidence", 5},
		{"a Start row with another evidence than command and no comment", func(r []Parameter) { r[3].Evidence = "because" }, "evidence", 5},
		{"a Start row with no value", func(r []Parameter) { r[4].Value = "" }, "value", 6},
		{"a default that is not —", func(r []Parameter) { r[1].Default = "2" }, "default", 3},
		{"a name that is not the table's, in its place", func(r []Parameter) { r[1].Name = "stall.T" }, "name", 3},
	}
	for _, tc := range cases {
		rows := base()
		tc.edit(rows)
		err := CheckParameters(parameterTable, rows)
		if tc.column == "" {
			if err != nil {
				t.Errorf("%s: %v", tc.name, err)
			}
			continue
		}
		refusedAt(t, tc.name, err, tc.line, tc.column)
	}
	if err := CheckParameters(parameterTable, base()[:5]); err == nil {
		t.Error("a table with five rows of the six was accepted: a check that runs only when its row is there")
	}
	if err := CheckParameters(parameterTable, nil); err == nil {
		t.Error("a header and no row was accepted")
	}
}

func TestAParametersFileWithAnUnknownNameOrAValueOutOfItsRowIsRefused(t *testing.T) {
	var b bytes.Buffer
	if err := WriteParameters(&b, defaultParams(t)); err != nil {
		t.Fatal(err)
	}
	good := b.String()
	if strings.Count(good, "\n") != 7 {
		t.Fatalf("the table as written has %d lines; want its header and six rows:\n%s", strings.Count(good, "\n"), good)
	}
	if _, err := ReadParameters([]byte(good)); err != nil {
		t.Fatalf("the table as written: %v", err)
	}
	set := func(name, value string) string { // the row of name as a set of value, comment 9
		return replaceRow(good, name, func(f []string) { f[1], f[5], f[6] = value, "why", "9" })
	}
	for _, tc := range []struct {
		name, data, column string
	}{
		{"an unknown name in place of a row", strings.Replace(good, "build.parallel\t", "stall.T\t", 1), "name"},
		{"a value of another type", set("brief.lines", "abc"), "value"},
		{"a value out of the bound", set("brief.lines", "0"), "value"},
		{"a bound that is not the words of the code", replaceRow(good, "brief.lines", func(f []string) { f[3] = "any number" }), "bound"},
		{"a role that is not the code's", replaceRow(good, "brief.lines", func(f []string) { f[4] = "idea-owner" }), "role"},
		{"an intake.cap that is not MONEY,HOURS", set("intake.cap", "50.0"), "value"},
		{"a transitions value that is not a SHA-256", set("transitions", "abc"), "value"},
		{"a row of a name twice", good + strings.SplitN(good, "\n", 3)[1] + "\n", "name"},
		{"the header with a column renamed", strings.Replace(good, "evidence", "proof", 1), "evidence"},
	} {
		_, err := ReadParameters([]byte(tc.data))
		if err == nil || !strings.Contains(err.Error(), tc.column) {
			t.Errorf("%s: %v; want an error that names %q", tc.name, err, tc.column)
		}
	}
}

// replaceRow gives the table with the row of name edited by f, which gets its fields.
func replaceRow(table, name string, f func([]string)) string {
	lines := strings.Split(table, "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, name+"\t") {
			fields := strings.Split(l, "\t")
			f(fields)
			lines[i] = strings.Join(fields, "\t")
		}
	}
	return strings.Join(lines, "\n")
}

func TestTheFormOfASetIsOneLineAndAnythingElseIsNotASet(t *testing.T) {
	for _, tc := range []struct {
		body        string
		isSet, form bool // isSet: the comment is a set; form: its form is right
		want        setLine
	}{
		{"set brief.lines 7 because the idea owner reads seven lines", true, true, setLine{"brief.lines", "7", "the idea owner reads seven lines"}},
		{"set brief.lines 7 because x\n", true, true, setLine{"brief.lines", "7", "x"}},
		{"set brief.lines 7 because x  \r\n", true, true, setLine{"brief.lines", "7", "x"}},
		{"set brief.lines 7 because the reason  holds   spaces because", true, true, setLine{"brief.lines", "7", "the reason  holds   spaces because"}},
		{"set transitions " + strings.Repeat("a", 64) + " because the Operator's table", true, true, setLine{"transitions", strings.Repeat("a", 64), "the Operator's table"}},
		{"", false, false, setLine{}},
		{"thanks, that reads well", false, false, setLine{}},
		{"Set brief.lines 7 because x", false, false, setLine{}},
		{"settle this 7 because x", false, false, setLine{}},
		{"please set brief.lines 7 because x", false, false, setLine{}},
		{"set", true, false, setLine{}},
		{"set brief.lines", true, false, setLine{}},
		{"set brief.lines 7", true, false, setLine{}},
		{"set brief.lines 7 because", true, false, setLine{}},
		{"set brief.lines 7 because ", true, false, setLine{}},
		{"set brief.lines 7 because —", true, false, setLine{}},
		{"set brief.lines  7 because x", true, false, setLine{}},
		{"set  brief.lines 7 because x", true, false, setLine{}},
		{"set brief.lines 7 Because x", true, false, setLine{}},
		{"set brief.lines 7 becausex", true, false, setLine{}},
		{"set brief.lines\t7 because x", true, false, setLine{}},
		{"set brief.lines 7 because x\nand a second line", true, false, setLine{}},
		{"set brief.lines 7 because x\n\nset build.parallel 2 because y", true, false, setLine{}},
		{"set brief.lines a b because x", true, false, setLine{}},
	} {
		got, isSet, err := parseSet(tc.body)
		if isSet != tc.isSet || (tc.isSet && (err == nil) != tc.form) || (!tc.isSet && err != nil) || (tc.form && got != tc.want) {
			t.Errorf("%q: line %+v, set %v, err %v; want %+v, set %v, form right %v", tc.body, got, isSet, err, tc.want, tc.isSet, tc.form)
		}
	}
}

func TestASetIsRefusedByTheFirstFailingItemInTheOrderOfTheSpecification(t *testing.T) {
	deps := append(append([]paramSpec(nil), parameterTable...), paramSpec{Name: "x.deps", Type: "text", Bound: "words", Role: "idea-owner", ok: func(string) bool { return true }})
	both := append(append([][]string(nil), paramApprovers...), []string{"101", "idea-owner", "pharzam", "2026-10-08T12:00:00Z", "start"})
	current := defaultParams(t)
	current = append(current, Parameter{Name: "x.deps", Bound: "words", Role: "idea-owner", Since: paramsAt})
	appSet := aSet(10, "101", "set brief.lines 7 because x")
	appSet.Row[5] = "layup-agent"
	edited := aSet(10, "101", "set brief.lines 7 because x")
	edited.Row[1] = "2"
	// The author comes before the name: an edit or an App's comment of an unknown name is refused as author.
	editedUnknown := aSet(10, "101", "set stall.T 5 because x")
	editedUnknown.Row[1] = "2"
	appUnknown := aSet(10, "101", "set stall.T 5 because x")
	appUnknown.Row[5] = "layup-agent"
	for _, tc := range []struct {
		name      string
		table     []paramSpec
		approvers [][]string
		in        setInput
		item      string // "" for a set that passes
	}{
		{"the Operator sets a row of the Operator", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines 7 because x"), ""},
		{"a comment that is not of the form", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines 7"), "form"},
		{"a malformed set by no approver and of an unknown name", parameterTable, paramApprovers, aSet(10, "999", "set stall.T"), "form"},
		{"no approver, an unknown name", parameterTable, paramApprovers, aSet(10, "999", "set stall.T 5 because x"), "author"},
		{"no approver, a good set", parameterTable, paramApprovers, aSet(10, "999", "set brief.lines 7 because x"), "author"},
		{"the idea owner on a row of the Operator", parameterTable, paramApprovers, aSet(10, "303", "set brief.lines 7 because x"), "author"},
		{"an approver of a point on a row of the Operator", parameterTable, paramApprovers, aSet(10, "404", "set brief.lines 7 because x"), "author"},
		{"a comment of an App", parameterTable, paramApprovers, appSet, "author"},
		{"an edit (a second copy) is never a decision", parameterTable, paramApprovers, edited, "author"},
		{"an edit of a comment of an unknown name", parameterTable, paramApprovers, editedUnknown, "author"},
		{"an App's comment of an unknown name", parameterTable, paramApprovers, appUnknown, "author"},
		{"an unknown name by the Operator", parameterTable, paramApprovers, aSet(10, "101", "set stall.T 5 because x"), "name"},
		{"an unknown name and a value of no type", parameterTable, paramApprovers, aSet(10, "101", "set stall.T abc because x"), "name"},
		{"a value of no type", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines abc because x"), "value"},
		{"a value that is —", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines — because x"), "value"},
		{"a value with a leading zero", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines 07 because x"), "value"},
		{"a negative value", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines -1 because x"), "value"},
		{"a value out of its bound", parameterTable, paramApprovers, aSet(10, "101", "set brief.lines 0 because x"), "bound"},
		{"the idea owner sets a row of the idea owner (an injected table)", deps, paramApprovers, aSet(10, "303", "set x.deps github.com/a/b because x"), ""},
		{"the Operator does not set a row of the idea owner", deps, paramApprovers, aSet(10, "101", "set x.deps github.com/a/b because x"), "author"},
		{"one account of both roles sets a row of the idea owner", deps, both, aSet(10, "101", "set x.deps github.com/a/b because x"), ""},
	} {
		j := judgeSet(tc.table, tc.approvers, tc.in, current, noCopy)
		id := int64(10)
		switch {
		case j.Superseded:
			t.Errorf("%s: superseded", tc.name)
		case tc.item == "" && (j.Change == nil || j.Refusal != nil):
			t.Errorf("%s: %+v; want a change", tc.name, j)
		case tc.item == "" && (j.Change.Comment != id || j.Change.Reason != "x"):
			t.Errorf("%s: the change %+v", tc.name, *j.Change)
		case tc.item != "" && (j.Change != nil || j.Refusal == nil || j.Refusal.Item != tc.item || j.Refusal.Comment != id || j.Refusal.Reason == ""):
			t.Errorf("%s: %+v; want a refusal of the item %s with its reason and the comment's ID", tc.name, j, tc.item)
		}
	}
}

// C2 of the plan review: a set is applied only when its comment is later than the comment that the row records.
func TestASetThatIsAtOrBelowTheCommentOfItsRowIsSuperseded(t *testing.T) {
	current := defaultParams(t)
	current[0].Value, current[0].Evidence, current[0].Comment = "5", "later", 50
	for _, tc := range []struct {
		id         int64
		body       string
		superseded bool
	}{
		{40, "set brief.lines 7 because earlier", true},
		{50, "set brief.lines 7 because the same comment", true},
		{51, "set brief.lines 7 because later", false},
		{40, "set brief.lines 0 because earlier and out of bound", true}, // no report for a set that a later one superseded
		{40, "set build.parallel 2 because another row has no comment", false},
	} {
		j := judgeSet(parameterTable, paramApprovers, aSet(tc.id, "101", tc.body), current, noCopy)
		if j.Superseded != tc.superseded || (tc.superseded && (j.Change != nil || j.Refusal != nil)) || (!tc.superseded && j.Change == nil) {
			t.Errorf("comment %d %q: %+v; want superseded %v", tc.id, tc.body, j, tc.superseded)
		}
	}
}

func TestEachRowOfTheTableAcceptsJustInsideItsBoundAndRefusesJustOutside(t *testing.T) {
	table := handoffTable(t, handoff.Default())
	lookup := func(sha string) ([]byte, bool) { return table, sha == sum(table) }
	for _, tc := range []struct {
		name, value, item string // item "" for a value that passes
	}{
		{"brief.lines", "1", ""}, {"brief.lines", "2", ""}, {"brief.lines", "100000", ""}, {"brief.lines", "0", "bound"}, {"brief.lines", "-1", "value"}, {"brief.lines", "1.5", "value"}, {"brief.lines", "1e2", "value"},
		{"build.parallel", "1", ""}, {"build.parallel", "0", "bound"}, {"build.parallel", "00", "value"}, {"build.parallel", "four", "value"},
		{"lease.H", "1", ""}, {"lease.H", "0", "bound"}, {"lease.H", "x", "value"},
		{"watch.T", "1", ""}, {"watch.T", "0", "bound"}, {"watch.T", "x", "value"},
		{"intake.cap", "50.0,8.0", ""}, {"intake.cap", "0.0,0.5", ""}, {"intake.cap", "50,8", "value"}, {"intake.cap", "50.0", "value"}, {"intake.cap", "50.0,8.0,1.0", "value"}, {"intake.cap", "50.0;8.0", "value"},
		{"transitions", sum(table), ""}, {"transitions", strings.ToUpper(sum(table)), "value"}, {"transitions", sum(table)[1:], "value"}, {"transitions", strings.Repeat("0", 64), "bound"},
	} {
		j := judgeSet(parameterTable, paramApprovers, aSet(10, "101", "set "+tc.name+" "+tc.value+" because x"), defaultParams(t), lookup)
		switch {
		case tc.item == "" && j.Change == nil:
			t.Errorf("%s %q: %+v; want it accepted", tc.name, tc.value, j)
		case tc.item != "" && (j.Refusal == nil || j.Refusal.Item != tc.item):
			t.Errorf("%s %q: %+v; want the item %s", tc.name, tc.value, j, tc.item)
		}
	}
}

func TestApplyingASetReplacesOneRowAndLeavesTheCallersTable(t *testing.T) {
	rows := defaultParams(t)
	saved := slicesClone(rows)
	later := paramsAt.Add(time.Hour)
	got, err := applyChange(rows, Change{Comment: 60, Name: "build.parallel", Value: "4", Reason: "the host has four cores"}, later)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rows, saved) {
		t.Errorf("the caller's table changed: %+v", rows)
	}
	if len(got) != len(rows) {
		t.Fatalf("the result has %d rows; want the %d of the table", len(got), len(rows))
	}
	for i, r := range got {
		if r.Name != "build.parallel" {
			if !reflect.DeepEqual(r, saved[i]) {
				t.Errorf("row %s changed: %+v", r.Name, r)
			}
			continue
		}
		want := saved[i]
		want.Value, want.Evidence, want.Comment, want.Since = "4", "the host has four cores", 60, later
		if !reflect.DeepEqual(r, want) || r.Default != "" || r.Bound == "" || r.Role != "operator" {
			t.Errorf("the row %+v; want %+v with its default, bound and role kept", r, want)
		}
	}
	got[0].Value = "changed"
	if !reflect.DeepEqual(rows, saved) {
		t.Error("the result shares its rows with the caller's table")
	}
	if _, err := applyChange(rows, Change{Comment: 61, Name: "stall.T", Value: "5", Reason: "x"}, later); err == nil {
		t.Error("a change of a name that is no row was applied")
	}
}

func sum(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }

func handoffTable(t *testing.T, rows []handoff.Row) []byte {
	t.Helper()
	var b bytes.Buffer
	if err := handoff.WriteTable(&b, rows); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

func TestASetOfTransitionsIsJudgedByTheCopyItNames(t *testing.T) {
	withRendered := func(items ...string) []byte {
		rows := handoff.Default()
		for i := range rows {
			if rows[i].From == "code (rendered task)" {
				rows[i].Items = items
			}
		}
		return handoffTable(t, rows)
	}
	dropped := handoff.Default()
	dropped = append(dropped[:3:3], dropped[4:]...)
	unknownKind := append(handoff.Default(), handoff.Row{From: "a", To: "b", Artifacts: []string{"plan"}, Items: []string{"schema"}})
	good := handoffTable(t, handoff.Default())
	extra := handoffTable(t, unknownKind)
	for _, tc := range []struct {
		name   string
		body   []byte
		item   string // "" for a table that passes
		reason string
	}{
		{"the default table", good, "", ""},
		{"the default and a pair of the Operator's", extra, "", ""},
		{"a row of code's own handoff with no schema", withRendered("artifacts", "rendering-equal"), "", ""},
		{"a missing default pair", handoffTable(t, dropped), "bound", "developer (specification)"},
		{"a row of code's own handoff that lists schema", withRendered("schema", "artifacts", "rendering-equal"), "bound", "code (rendered task)"},
		{"a body that is no table", []byte("the table is below\n"), "bound", ""},
		{"a table in a code fence", []byte("```\n" + string(good) + "```\n"), "bound", ""},
		{"a table with CRLF line ends", []byte(strings.ReplaceAll(string(good), "\n", "\r\n")), "bound", ""},
		{"an unknown item", []byte(strings.Replace(string(good), "plan-fields", "vibes", 1)), "bound", "vibes"},
		{"an unknown kind", []byte(strings.Replace(string(good), "\tplan\t", "\tblueprint\t", 1)), "bound", "blueprint"},
		{"an empty body", nil, "bound", ""},
	} {
		sha := sum(tc.body)
		lookup := func(s string) ([]byte, bool) { return tc.body, s == sha }
		j := judgeSet(parameterTable, paramApprovers, aSet(10, "101", "set transitions "+sha+" because the Operator's table"), defaultParams(t), lookup)
		switch {
		case tc.item == "" && (j.Change == nil || !bytes.Equal(j.Change.Table, tc.body) || j.Change.Value != sha):
			t.Errorf("%s: %+v; want a change that carries the copy's bytes", tc.name, j)
		case tc.item != "" && (j.Refusal == nil || j.Refusal.Item != tc.item || !strings.Contains(j.Refusal.Reason, tc.reason)):
			t.Errorf("%s: %+v; want the item %s and a reason that holds %q", tc.name, j, tc.item, tc.reason)
		}
	}
	// no copy has the hash
	j := judgeSet(parameterTable, paramApprovers, aSet(10, "101", "set transitions "+sum(good)+" because x"), defaultParams(t), noCopy)
	if j.Refusal == nil || j.Refusal.Item != "bound" || !strings.Contains(j.Refusal.Reason, "copy") {
		t.Errorf("a hash of no copy: %+v", j)
	}
}

func TestParameterFilesAreWrittenOnlyWhenAbsent(t *testing.T) {
	files, err := parameterFiles(nil, startValue, paramsAt)
	if err != nil || len(files) != 1 {
		t.Fatalf("the files of a first write: %v, %v", files, err)
	}
	rows, err := ReadParameters(files["parameters.tsv"])
	if err != nil || !reflect.DeepEqual(rows, defaultParams(t)) {
		t.Errorf("parameters.tsv: %+v, %v; want the default rows", rows, err)
	}
	rows[1].Value, rows[1].Evidence, rows[1].Comment = "3", "why", 9
	var kept bytes.Buffer
	if err := WriteParameters(&kept, rows); err != nil {
		t.Fatal(err)
	}
	if again, err := parameterFiles(kept.Bytes(), startValue, paramsAt.Add(time.Hour)); err != nil || len(again) != 0 {
		t.Errorf("a table that exists: %v, %v; want it kept and nothing written", again, err)
	}
	for name, data := range map[string]string{"a file that is no table": "nothing\n", "a header alone": string(kept.Bytes()[:strings.Index(kept.String(), "\n")+1])} {
		if files, err := parameterFiles([]byte(data), startValue, paramsAt); err == nil || files != nil {
			t.Errorf("%s: %v, %v; want a refusal and no file", name, files, err)
		}
	}
}
