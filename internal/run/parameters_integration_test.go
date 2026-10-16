//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/tsv"
)

// The parameters on the real adapter, the httptest forge and real git (task T-dx1h, #358, row 66b of the plan of M2e).
// No step `loop` exists yet (row 66k, #367, which owns the call between steps); a test calls what that step will call,
// on the state of a restart.

func TestTheSchemaOfParametersEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["parameters"]
	if !ok {
		t.Fatal("docs/spec/ has no block parameters")
	}
	if err := tsv.Compare(block, ParametersSchema); err != nil {
		t.Error(err)
	}
}

// The table of the code equals the table of loop.md ("The parameters of M2e"): the owner's test that the spec asks for,
// and the six rows of the plan review's C3, each with a source for each of its other cells.
func TestTheParameterTableEqualsTheSpecification(t *testing.T) {
	data, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	_, section, found := strings.Cut(text, "\n### The parameters of M2e\n")
	if !found {
		t.Fatal("loop.md has no section 'The parameters of M2e'")
	}
	if i := strings.Index(section, "\n#"); i >= 0 {
		section = section[:i]
	}
	var rows [][]string
	for _, l := range strings.Split(section, "\n") {
		if !strings.HasPrefix(l, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(l, "| "), " | ")
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		rows = append(rows, cells)
	}
	if len(rows) < 3 || strings.Join(rows[0], "|") != "Name|Type|Bound|Role|Default|Source" {
		t.Fatalf("the table of loop.md: %d rows, header %v; want Name | Type | Bound | Role | Default | Source", len(rows), rows)
	}
	rows = rows[2:] // the header and its rule
	if len(rows) != len(parameterTable) {
		t.Fatalf("loop.md has %d rows, the code %d; want the same six", len(rows), len(parameterTable))
	}
	unquote := func(s string) string { return strings.Trim(s, "`") }
	for i, spec := range parameterTable {
		r := rows[i]
		if len(r) != 6 {
			t.Errorf("row %d of loop.md has %d cells; want 6", i+1, len(r))
			continue
		}
		if unquote(r[0]) != spec.Name || unquote(r[1]) != spec.Type || r[2] != spec.Bound || unquote(r[3]) != spec.Role || unquote(r[4]) != "—" {
			t.Errorf("row %d of loop.md %q; the code has name %q, type %q, bound %q, role %q, default —", i+1, r, spec.Name, spec.Type, spec.Bound, spec.Role)
		}
		for _, cell := range []string{"name:", "type:", "bound:", "role:", "default:"} {
			if !strings.Contains(r[5], cell) {
				t.Errorf("the source of %s has no %s cell: %q", spec.Name, cell, r[5])
			}
		}
	}
	// The open gap of The bet stays: its marker is written as escapes here, as a marker in a test file would be unlisted.
	flat := strings.Join(strings.Fields(text), " ") // the sentences wrap
	for _, want := range []string{"the others are read by their milestones", "\u2039the line limit of the bet brief\u203a"} {
		if !strings.Contains(flat, want) {
			t.Errorf("loop.md lost %q", want)
		}
	}
}

// seededTable is the table of defaults as text, made here and not by the code under test, so that a test of the boundary
// does not rest on the first writer.
func seededTable(t *testing.T) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("name\tvalue\tdefault\tbound\trole\tevidence\tcomment\tsince\n")
	for _, r := range defaultParams(t) {
		dash := func(s string) string {
			if s == "" {
				return "—"
			}
			return s
		}
		b.WriteString(strings.Join([]string{r.Name, dash(r.Value), "—", r.Bound, r.Role, dash(r.Evidence), "—", "2026-10-08T12:30:00Z"}, "\t") + "\n")
	}
	return b.String()
}

// loopWorld is a world whose first task is registered, with a table of parameters that the test made.
func loopWorld(t *testing.T) (*world, *state) {
	t.Helper()
	w := registerWorld(t, loopPerms)
	pushRecords(t, w.bare, map[string]string{"parameters.tsv": seededTable(t)}, "the table of defaults, made by the test")
	r := restarted(t, w)
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Fatal(err)
	}
	return w, restarted(t, w)
}

// say adds a comment of the account author on the control issue (2 in these worlds).
func say(w *world, id, author int, login, body string) {
	addComment(w, 2, id, author, login, "2026-10-08T13:00:00Z", body)
}

const (
	operatorID, ideaOwnerID, strangerID = 101, 303, 999
)

// A commitSummary is one commit of the records branch: its subject and the files it changes, in a line.
type commitSummary struct{ subject, files string }

func commitsSince(t *testing.T, w *world, before string) []commitSummary {
	t.Helper()
	var out []commitSummary
	for _, id := range strings.Fields(gitOut(t, w.bare, "rev-list", "--reverse", before+".."+recordsHead(t, w))) {
		files := strings.Fields(gitOut(t, w.bare, "diff-tree", "--no-commit-id", "--name-only", "-r", id))
		out = append(out, commitSummary{strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%s", id)), strings.Join(files, ",")})
	}
	return out
}

func parametersAt(t *testing.T, w *world, rev string) []Parameter {
	t.Helper()
	rows, err := ReadParameters([]byte(show(t, w.bare, rev, "parameters.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

func rowNamed(rows []Parameter, name string) Parameter {
	for _, r := range rows {
		if r.Name == name {
			return r
		}
	}
	return Parameter{}
}

// progressLines gives the progress lines of the world that hold word.
func progressLines(w *world, word string) []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	var out []string
	for _, l := range w.lines {
		if strings.Contains(l, word) {
			out = append(out, l)
		}
	}
	return out
}

func TestTheFirstTaskWritesTheDefaultParametersFromStartTsv(t *testing.T) {
	w := newWorld(t)
	cfg := w.config("1.0.0")
	cfg.LeaseH, cfg.WatchT, cfg.IntakeCap = 7, 3, "20.0,4.0"
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	w.fake.mu.Lock()
	w.fake.perms = loopPerms
	w.fake.mu.Unlock()
	r := restarted(t, w) // the restart's own flags are not Start's: start.tsv is the evidence
	from := w.clock.Now().Truncate(time.Second)
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Fatal(err)
	}
	to := w.clock.Now()
	rows := parametersAt(t, w, "refs/heads/layup-records")
	if len(rows) != 6 {
		t.Fatalf("parameters.tsv has %d rows; want six", len(rows))
	}
	for name, want := range map[string]string{"lease.H": "7", "watch.T": "3", "intake.cap": "20.0,4.0"} {
		if r := rowNamed(rows, name); r.Value != want || r.Evidence != "command" || r.Comment != 0 {
			t.Errorf("%s is %+v; want %q from start.tsv, evidence command, no comment", name, r, want)
		}
	}
	for _, r := range rows {
		if r.Since.Before(from) || r.Since.After(to) {
			t.Errorf("%s: since %s is outside the clock, %s to %s", r.Name, r.Since, from, to)
		}
	}
	for _, name := range []string{"brief.lines", "build.parallel", "transitions"} {
		if r := rowNamed(rows, name); r.Value != "" || r.Evidence != "" || r.Comment != 0 {
			t.Errorf("%s is %+v; want unset (the 100 of the demo target is the Operator's input at row 68b, not a default)", name, r)
		}
	}
}

func TestAnExistingParametersFileThatItsReaderRefusesRefusesTheTaskAndWritesNothing(t *testing.T) {
	w := registerWorld(t, loopPerms)
	pushRecords(t, w.bare, map[string]string{"parameters.tsv": "name\tvalue\nbrief.lines\t7\n"}, "a table by hand")
	r := restarted(t, w)
	before := recordsHead(t, w)
	err := r.registerTask(context.Background(), firstTask())
	if err == nil || !strings.Contains(err.Error(), "parameters.tsv") {
		t.Errorf("registerTask: %v; want a refusal that names parameters.tsv", err)
	}
	if after := recordsHead(t, w); after != before {
		t.Errorf("the records head moved from %s to %s", before, after)
	}
}

func TestAnExistingParametersFileIsKeptByTheFirstTask(t *testing.T) {
	w := registerWorld(t, loopPerms)
	r := restarted(t, w)
	rows, err := DefaultParameters(parameterTable, r.value, w.clock.Now().Truncate(time.Second))
	if err != nil || len(rows) != 6 {
		t.Fatalf("the default rows: %d, %v", len(rows), err)
	}
	rows[1].Value, rows[1].Evidence, rows[1].Comment = "3", "the host has three cores", 9000
	var b bytes.Buffer
	if err := WriteParameters(&b, rows); err != nil {
		t.Fatal(err)
	}
	pushRecords(t, w.bare, map[string]string{"parameters.tsv": b.String()}, "a table of an earlier run")
	r = restarted(t, w)
	before := recordsHead(t, w)
	if err := r.registerTask(context.Background(), firstTask()); err != nil {
		t.Fatal(err)
	}
	if names := gitOut(t, w.bare, "diff", "--name-only", before, recordsHead(t, w)); names != "tasks.tsv\ntransitions.tsv\n" {
		t.Errorf("the commit changes %q; want tasks.tsv and transitions.tsv, the table that exists kept", names)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "parameters.tsv"); got != b.String() {
		t.Errorf("parameters.tsv after the first task:\n%s\nwant the table that existed:\n%s", got, b.String())
	}
}

func TestASetIsCopiedAndAppliedAtTheBoundaryAsOneRecordsCommit(t *testing.T) {
	w, r := loopWorld(t)
	say(w, 9101, operatorID, "pharzam", "set brief.lines 7 because the idea owner reads seven lines")
	r = restarted(t, w) // a restart: no state of the first run
	before := recordsHead(t, w)
	prior := parametersAt(t, w, before)
	if got := rowNamed(prior, "brief.lines"); got.Value != "" {
		t.Fatalf("before the boundary brief.lines is %+v", got)
	}
	from := w.clock.Now().Truncate(time.Second)
	applied, refused, err := r.parametersBoundary(context.Background())
	if err != nil || len(refused) != 0 || len(applied) != 1 {
		t.Fatalf("the boundary: applied %+v, refused %+v, %v; want one change", applied, refused, err)
	}
	to := w.clock.Now()
	if c := applied[0].Change; c.Comment != 9101 || c.Name != "brief.lines" || c.Value != "7" || applied[0].Note != "" {
		t.Errorf("the applied change %+v", applied[0])
	}
	got := commitsSince(t, w, before)
	want := []commitSummary{{"records: the copy of the comment 9101", "copies.tsv,copies/9101-1.md"}, {"parameters: set brief.lines", "parameters.tsv"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the commits %+v; want the copy and then the change, %+v", got, want)
	}
	if who := gitOut(t, w.bare, "log", "-1", "--format=%an <%ae>|%cn <%ce>", recordsHead(t, w)); who != "layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>|layup-agent[bot] <9001+layup-agent[bot]@users.noreply.github.com>\n" {
		t.Errorf("the author and committer of the change %q; want the App's bot", who)
	}
	rows := parametersAt(t, w, recordsHead(t, w))
	row := rowNamed(rows, "brief.lines")
	wantRow := prior[0]
	wantRow.Value, wantRow.Evidence, wantRow.Comment, wantRow.Since = "7", "the idea owner reads seven lines", 9101, row.Since
	if !reflect.DeepEqual(row, wantRow) || row.Since.Before(from) || row.Since.After(to) {
		t.Errorf("the row %+v; want %+v, since in %s to %s", row, wantRow, from, to)
	}
	for i := range rows {
		if rows[i].Name != "brief.lines" && !reflect.DeepEqual(rows[i], prior[i]) {
			t.Errorf("the row %s changed: %+v", rows[i].Name, rows[i])
		}
	}
	if l := progressLines(w, "parameters: applied"); len(l) != 1 || !strings.Contains(l[0], "brief.lines") {
		t.Errorf("the progress lines of the change %q; want one that names brief.lines", l)
	}
}

func TestASetOfAStartParameterIsRecordedAndSaysThatThisRunDoesNotReadIt(t *testing.T) {
	w, _ := loopWorld(t)
	say(w, 9101, operatorID, "pharzam", "set lease.H 9 because the Operator wants a longer lease")
	r := restarted(t, w)
	startBlob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:start/start.tsv")
	applied, refused, err := r.parametersBoundary(context.Background())
	if err != nil || len(refused) != 0 || len(applied) != 1 {
		t.Fatalf("the boundary: applied %+v, refused %+v, %v", applied, refused, err)
	}
	if !strings.Contains(applied[0].Note, "recorded") || !strings.Contains(applied[0].Note, "start.tsv") {
		t.Errorf("the note %q; want it to say that the change is recorded and that this run reads start.tsv", applied[0].Note)
	}
	if row := rowNamed(parametersAt(t, w, recordsHead(t, w)), "lease.H"); row.Value != "9" || row.Comment != 9101 {
		t.Errorf("lease.H is %+v; want the set recorded", row)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:start/start.tsv") != startBlob {
		t.Error("the change rewrote start.tsv")
	}
	if l := progressLines(w, "start.tsv"); len(l) == 0 {
		t.Error("no progress line says that the run reads start.tsv")
	}
}

func TestARefusedSetChangesNothingAndIsReportedWithItsItem(t *testing.T) {
	w, _ := loopWorld(t)
	for _, c := range []struct {
		id, author int
		body       string
	}{
		{9101, strangerID, "set brief.lines 7 because a stranger asks"},
		{9102, operatorID, "set stall.T 5 because the name is not a row"},
		{9103, operatorID, "set brief.lines 0 because zero lines"},
		{9104, operatorID, "set brief.lines abc because not a number"},
		{9105, operatorID, "set brief.lines"},
		{9106, ideaOwnerID, "set brief.lines 7 because the idea owner has no right to this row"},
		{9107, operatorID, "thanks, that reads well"},
	} {
		say(w, c.id, c.author, "login", c.body)
	}
	// A comment of an App is no human input: it is not copied and not judged.
	w.fake.mu.Lock()
	w.fake.comments[2] = append(w.fake.comments[2], map[string]any{"id": 9108, "user": map[string]any{"id": 777, "login": "layup-watch[bot]"}, "performed_via_github_app": map[string]any{"slug": "layup-watch"},
		"created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": "set brief.lines 7 because an App asks"})
	w.fake.mu.Unlock()
	r := restarted(t, w)
	before := recordsHead(t, w)
	blob := gitOut(t, w.bare, "rev-parse", before+":parameters.tsv")
	items := func(refused []Refusal) []string {
		var out []string
		for _, x := range refused {
			out = append(out, x.Item)
		}
		return out
	}
	wantItems := []string{"author", "name", "bound", "value", "form", "author"}
	applied, refused, err := r.parametersBoundary(context.Background())
	if err != nil || len(applied) != 0 || !reflect.DeepEqual(items(refused), wantItems) {
		t.Fatalf("the boundary: applied %+v, items %v, %v; want none applied and the items %v", applied, items(refused), err, wantItems)
	}
	for i, id := range []int64{9101, 9102, 9103, 9104, 9105, 9106} {
		if refused[i].Comment != id || refused[i].Reason == "" {
			t.Errorf("refusal %d is %+v; want the comment %d and a reason", i, refused[i], id)
		}
	}
	for _, c := range commitsSince(t, w, before) {
		if !strings.HasPrefix(c.subject, "records: the copy of the comment") {
			t.Errorf("a commit that is no copy: %+v", c)
		}
	}
	if gitOut(t, w.bare, "rev-parse", recordsHead(t, w)+":parameters.tsv") != blob {
		t.Error("a refused set changed parameters.tsv")
	}
	if strings.Contains(show(t, w.bare, "refs/heads/layup-records", "copies.tsv"), "9108\t") {
		t.Error("a comment of an App was copied")
	}
	if n := len(progressLines(w, "parameters: refused")); n != 6 {
		t.Errorf("%d progress lines of a refusal; want 6", n)
	}
	// The copy is the record: a second boundary judges again, with the same result, and writes nothing.
	again := restarted(t, w)
	head := recordsHead(t, w)
	applied, refused, err = again.parametersBoundary(context.Background())
	if err != nil || len(applied) != 0 || !reflect.DeepEqual(items(refused), wantItems) {
		t.Errorf("the second boundary: applied %+v, items %v, %v", applied, items(refused), err)
	}
	if recordsHead(t, w) != head {
		t.Error("the second boundary wrote a commit")
	}
}

// C2 of the plan review: a set is applied at most once, a superseded one never again, in the order of the comments'
// IDs, not of the forge's list.
func TestASetThatIsAppliedIsNotAppliedAgainAndTheLaterOfTwoWins(t *testing.T) {
	w, _ := loopWorld(t)
	say(w, 9102, operatorID, "pharzam", "set build.parallel 4 because the host has four cores")
	say(w, 9101, operatorID, "pharzam", "set build.parallel 2 because the host has two cores") // listed after, but the earlier comment
	r := restarted(t, w)
	before := recordsHead(t, w)
	applied, refused, err := r.parametersBoundary(context.Background())
	if err != nil || len(refused) != 0 || len(applied) != 2 || applied[0].Change.Comment != 9101 || applied[1].Change.Comment != 9102 {
		t.Fatalf("the boundary: applied %+v, refused %+v, %v; want 9101 and then 9102", applied, refused, err)
	}
	changes := func(since string) int {
		n := 0
		for _, c := range commitsSince(t, w, since) {
			if c.subject == "parameters: set build.parallel" {
				n++
			}
		}
		return n
	}
	if n := changes(before); n != 2 {
		t.Errorf("%d change commits; want one for each set", n)
	}
	if row := rowNamed(parametersAt(t, w, recordsHead(t, w)), "build.parallel"); row.Value != "4" || row.Comment != 9102 {
		t.Errorf("build.parallel is %+v; want the later set's value", row)
	}
	head := recordsHead(t, w)
	applied, refused, err = restarted(t, w).parametersBoundary(context.Background())
	if err != nil || len(applied) != 0 || len(refused) != 0 || recordsHead(t, w) != head {
		t.Errorf("the second boundary: applied %+v, refused %+v, %v, head moved %v; want nothing", applied, refused, err, recordsHead(t, w) != head)
	}
	// A later set is applied once, and the earlier two are not applied again over it.
	say(w, 9103, operatorID, "pharzam", "set build.parallel 8 because eight cores")
	applied, _, err = restarted(t, w).parametersBoundary(context.Background())
	if err != nil || len(applied) != 1 || applied[0].Change.Comment != 9103 || changes(head) != 1 {
		t.Errorf("the third boundary: applied %+v, %v, %d change commits; want only 9103", applied, err, changes(head))
	}
	head = recordsHead(t, w)
	if applied, _, err = restarted(t, w).parametersBoundary(context.Background()); err != nil || len(applied) != 0 || recordsHead(t, w) != head {
		t.Errorf("the fourth boundary: applied %+v, %v, head moved %v", applied, err, recordsHead(t, w) != head)
	}
	if row := rowNamed(parametersAt(t, w, recordsHead(t, w)), "build.parallel"); row.Value != "8" || row.Comment != 9103 {
		t.Errorf("build.parallel is %+v; want 8 of the last set", row)
	}
}

func TestAnEditOfACommentIsNoSet(t *testing.T) {
	w, _ := loopWorld(t)
	say(w, 9110, operatorID, "pharzam", "thanks, that reads well")
	if applied, refused, err := restarted(t, w).parametersBoundary(context.Background()); err != nil || len(applied)+len(refused) != 0 {
		t.Fatalf("the first boundary: %+v %+v %v", applied, refused, err)
	}
	w.fake.mu.Lock()
	for _, c := range w.fake.comments[2] {
		if c["id"] == 9110 {
			c["body"], c["updated_at"] = "set build.parallel 9 because an edit", "2026-10-08T14:00:00Z"
		}
	}
	w.fake.mu.Unlock()
	applied, refused, err := restarted(t, w).parametersBoundary(context.Background())
	if err != nil || len(applied)+len(refused) != 0 {
		t.Errorf("the boundary after the edit: applied %+v, refused %+v, %v; want none (the first copy is the decision)", applied, refused, err)
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "copies.tsv"), "9110\t2\t") {
		t.Error("the edit has no copy of its own")
	}
	if row := rowNamed(parametersAt(t, w, recordsHead(t, w)), "build.parallel"); row.Value != "" {
		t.Errorf("build.parallel is %+v after an edit", row)
	}
}

// customTable is the default table and a pair of the Operator's.
func customTable(t *testing.T) []byte {
	t.Helper()
	return handoffTable(t, append(handoff.Default(), handoff.Row{From: "verifier", To: "asking role", Artifacts: []string{"question"}, Items: []string{"schema"}}))
}

func TestSetTransitionsReplacesTheTableInOneCommitAndTheOldTableStaysInHistory(t *testing.T) {
	w, _ := loopWorld(t)
	table := customTable(t)
	say(w, 9201, ideaOwnerID, "carol", string(table)) // the author of the comment that holds the table is not checked
	say(w, 9202, operatorID, "pharzam", "set transitions "+sum(table)+" because the Operator adds a handoff")
	r := restarted(t, w)
	before := recordsHead(t, w)
	old := show(t, w.bare, before, "transitions.tsv")
	applied, refused, err := r.parametersBoundary(context.Background())
	if err != nil || len(refused) != 0 || len(applied) != 1 || applied[0].Change.Name != "transitions" {
		t.Fatalf("the boundary: applied %+v, refused %+v, %v", applied, refused, err)
	}
	got := commitsSince(t, w, before)
	if n := len(got); n != 3 || got[2] != (commitSummary{"parameters: set transitions", "parameters.tsv,transitions.tsv"}) {
		t.Errorf("the commits %+v; want two copies and then one change of both files", got)
	}
	head := recordsHead(t, w)
	if file := show(t, w.bare, head, "transitions.tsv"); file != string(table) {
		t.Errorf("transitions.tsv:\n%s\nwant the bytes of the copy:\n%s", file, table)
	}
	if row := rowNamed(parametersAt(t, w, head), "transitions"); row.Value != sum(table) || row.Comment != 9202 || row.Value != sum([]byte(show(t, w.bare, head, "transitions.tsv"))) {
		t.Errorf("the row transitions %+v; want the SHA-256 of the file in force", row)
	}
	if prior := show(t, w.bare, recordsHead(t, w)+"^", "transitions.tsv"); prior != old {
		t.Errorf("the table before the change in Git history:\n%s\nwant\n%s", prior, old)
	}
	if rows, err := handoff.ReadTable([]byte(show(t, w.bare, head, "transitions.tsv"))); err != nil || handoff.CheckTable(rows) != nil {
		t.Errorf("the new table does not read and pass: %v", err)
	}
	// The next task keeps the Operator's table.
	second := Task{ID: "T-bbbb", Issue: 8, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Predecessors: []string{"T-aaaa"}}
	if err := restarted(t, w).registerTask(context.Background(), second); err != nil {
		t.Fatal(err)
	}
	if show(t, w.bare, recordsHead(t, w), "transitions.tsv") != string(table) {
		t.Error("the next task rewrote the table from the default")
	}
}

func TestASetOfTransitionsThatFailsItsCheckWritesNothing(t *testing.T) {
	w, r := loopWorld(t)
	without := handoff.Default()
	without = append(without[:3:3], without[4:]...)
	rendered := handoff.Default()
	for i := range rendered {
		if rendered[i].From == "code (rendered task)" {
			rendered[i].Items = []string{"schema", "artifacts", "rendering-equal"}
		}
	}
	good := customTable(t)
	crlf := []byte(strings.ReplaceAll(string(good), "\n", "\r\n"))
	elsewhere := handoffTable(t, append(handoff.Default(), handoff.Row{From: "verifier", To: "asking role", Artifacts: []string{"answer"}, Items: []string{"schema"}}))
	// A copy of a comment of the Intake issue: its hash names a table that no comment of the control issue holds.
	addComment(w, 1, 9300, ideaOwnerID, "carol", "2026-10-08T13:00:00Z", string(elsewhere))
	for i, body := range [][]byte{handoffTable(t, without), handoffTable(t, rendered), crlf} {
		say(w, 9301+i, ideaOwnerID, "carol", string(body))
		say(w, 9311+i, operatorID, "pharzam", "set transitions "+sum(body)+" because the table is bad")
	}
	say(w, 9314, operatorID, "pharzam", "set transitions "+sum(elsewhere)+" because it is a copy of the Intake issue")
	say(w, 9315, operatorID, "pharzam", "set transitions "+strings.Repeat("0", 64)+" because no copy has this hash")
	say(w, 9316, operatorID, "pharzam", "set transitions abc because not a SHA-256")
	ctx := context.Background()
	comments, err := r.cfg.Forge.Comments(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	for _, cm := range comments {
		if cm.ID == 9300 {
			if err := r.copy(ctx, cm, 1); err != nil {
				t.Fatal(err)
			}
		}
	}
	r = restarted(t, w)
	before := recordsHead(t, w)
	blobs := gitOut(t, w.bare, "rev-parse", before+":parameters.tsv", before+":transitions.tsv")
	applied, refused, err := r.parametersBoundary(ctx)
	var got []string
	for _, x := range refused {
		got = append(got, x.Item)
	}
	want := []string{"bound", "bound", "bound", "bound", "bound", "value"}
	if err != nil || len(applied) != 0 || !reflect.DeepEqual(got, want) {
		t.Fatalf("the boundary: applied %+v, items %v, %v; want none applied and %v", applied, got, err, want)
	}
	if !strings.Contains(refused[1].Reason, "code (rendered task)") {
		t.Errorf("the reason of the rendered row %q; want it to name the row", refused[1].Reason)
	}
	if now := gitOut(t, w.bare, "rev-parse", recordsHead(t, w)+":parameters.tsv", recordsHead(t, w)+":transitions.tsv"); now != blobs {
		t.Errorf("a refused replacement changed a file: %s then %s", blobs, now)
	}
}

// The boundary writes through the fence of the lease, as the register does (row 66a).
func TestTheBoundaryIsFencedByTheLease(t *testing.T) {
	w, _ := loopWorld(t)
	say(w, 9101, operatorID, "pharzam", "set build.parallel 2 because two cores")
	r := restarted(t, w)
	// The comment is copied before the push is refused, so that the refusal falls on the change's own commit.
	ctx := context.Background()
	comments, err := r.cfg.Forge.Comments(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range comments {
		if c.App == "" {
			if err := r.copy(ctx, c, 2); err != nil {
				t.Fatal(err)
			}
		}
	}
	held := LeaseRow{Run: "fedcba9876543210", Host: "host-1", Version: "1.0.0", Started: opened, Heartbeat: 1, State: "held"}
	if err := r.store.WriteLease(context.Background(), held, "lease: held by the test"); err != nil {
		t.Fatal(err)
	}
	hook := "#!/bin/sh\nif [ -f refuse-once ]; then rm refuse-once; echo refused >&2; exit 1; fi\n"
	if err := os.WriteFile(filepath.Join(w.bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(w.bare, "refuse-once"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	_, _, err = r.parametersBoundary(ctx)
	var lost *LostError
	if !errors.As(err, &lost) || lost.Holder != "fedcba9876543210" {
		t.Errorf("the boundary: %v; want a LostError whose holder is the other run", err)
	}
	if row := rowNamed(parametersAt(t, w, recordsHead(t, w)), "build.parallel"); row.Value != "" {
		t.Errorf("build.parallel is %+v after a lost lease", row)
	}
}
