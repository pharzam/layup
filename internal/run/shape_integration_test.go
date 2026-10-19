//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// The step shape on the production path (task T-3f94, #364, row 66h): 66r's drvWorld restarted from its records, the production
// constructor state.sessions(), real git, the httptest forge and scripted harnesses.

// qrow is a row of questions.tsv of the batch: an intent field, or the question of the stack.
func qrow(id, ref, state string) []string {
	source, kind, block := "intent", "intent", "idea-owner"
	if ref == "S01-stack" {
		source, kind, block = "setup", "setup", "operator"
	}
	answered := ""
	if state == "answered" {
		answered = "2026-10-08T12:40:00Z"
	}
	return []string{id, source, ref, block, "before-delivery", kind, "", "", "The question " + id, "2026-10-08T12:30:00Z", "layup run", answered, state}
}

// push puts the records in the world and restarts it.
func (d *drvWorld) push(files map[string][]byte, message string) {
	d.t.Helper()
	in := map[string]string{}
	for p, data := range files {
		in[p] = string(data)
	}
	pushRecords(d.t, d.w.bare, in, message)
	d.restart()
}

// restart restarts the world as far as its clone, with the harnesses and models of the world's scripts in the configuration.
func (d *drvWorld) restart() {
	d.t.Helper()
	d.r = restarted(d.t, d.w)
	d.r.cfg.Harnesses, d.r.cfg.Models = d.harnesses, d.models
}

// questions gives the rows of questions.tsv of the records.
func (d *drvWorld) questions() [][]string {
	d.t.Helper()
	rows, err := records.ReadQuestions([]byte(keptRecord(d.t, d.w, "questions.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	return rows
}

func stateOf(rows [][]string) map[string]string {
	out := map[string]string{}
	for _, r := range rows {
		out[r[0]] = r[12]
	}
	return out
}

// has reports whether the records hold the path.
func (d *drvWorld) has(path string) bool {
	_, err := gitRunOK(d.w.bare, "cat-file", "-e", "layup-records:"+path)
	return err == nil
}

// appComments gives the bodies of the comments of the App on the Intake issue that start with the header.
func (d *drvWorld) appComments(header string) []string {
	d.w.fake.mu.Lock()
	defer d.w.fake.mu.Unlock()
	var out []string
	for _, c := range d.w.fake.comments[1] {
		if body, _ := c["body"].(string); strings.HasPrefix(body, header+"\n") {
			out = append(out, body)
		}
	}
	return out
}

// settled makes the records of an Intake with nothing left open and the budget written.
func (d *drvWorld) settled() {
	d.t.Helper()
	d.push(map[string][]byte{
		"questions.tsv":      cavTable(d.t, records.QuestionsSchema, [][]string{qrow("Q-001", "S01-stack", "answered"), qrow("Q-002", "band.B", "answered"), qrow("Q-003", "band.U", "answered"), qrow("Q-004", "appetite.hours", "answered")}),
		"budget.tsv":         cavTable(d.t, records.BudgetSchema, [][]string{{"band.B", "5.0", "Q-002"}, {"band.U", "10.0", "Q-003"}, {"appetite.hours", "40.0", "Q-004"}}),
		"intake/markers.tsv": cavTable(d.t, records.MarkersSchema, nil),
		"intake/answers.tsv": cavTable(d.t, records.AnswersSchema, [][]string{{"Q-002", "5.0", "802", "303", "2026-10-08T12:40:00Z"}, {"Q-003", "10.0", "802", "303", "2026-10-08T12:40:00Z"}, {"Q-004", "40.0", "802", "303", "2026-10-08T12:40:00Z"}}),
	}, "the records of an Intake with nothing left open")
}

// flakyForge fails the post of a comment, or the first read of the comments after a post.
type flakyForge struct {
	forge.Forge
	failPost, failAfter, armed bool
}

func (f *flakyForge) Comment(ctx context.Context, issue int, body string) (int64, error) {
	if f.failPost {
		return 0, context.DeadlineExceeded
	}
	id, err := f.Forge.Comment(ctx, issue, body)
	f.armed = err == nil && f.failAfter
	return id, err
}

func (f *flakyForge) Comments(ctx context.Context, issue int) ([]forge.Comment, error) {
	if f.armed {
		f.armed = false
		return nil, context.DeadlineExceeded
	}
	return f.Forge.Comments(ctx, issue)
}

// I1 (plan review conditions 6 and 7; loop.md, The questions left open): the questions left open are asked in one comment that is
// copied; a post that fails after the rows, and a read that fails after the post, are recovered with one comment; the answers set the
// states, an unparsed intent value waits and is asked again, and the correction writes budget.tsv with Shape's question as the source.
func TestTheQuestionsLeftOpenAreAskedAndTheIntentAnswerWritesBudget(t *testing.T) {
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	d.push(map[string][]byte{
		"questions.tsv": cavTable(t, records.QuestionsSchema, [][]string{qrow("Q-001", "S01-stack", "open"), qrow("Q-002", "band.B", "open"), qrow("Q-003", "band.U", "open"),
			qrow("Q-004", "appetite.hours", "open"), qrow("Q-005", "success-criteria", "open")}),
		"intake/markers.tsv": cavTable(t, records.MarkersSchema, nil), "intake/answers.tsv": cavTable(t, records.AnswersSchema, nil),
	}, "an Intake that left five questions open")
	ctx := context.Background()
	flaky := &flakyForge{Forge: d.r.cfg.Forge, failPost: true}
	d.r.cfg.Forge = flaky
	if _, err := d.r.shapeQuestions(ctx); err == nil || len(d.questions()) != 10 || len(d.appComments(shapeQuestionsHeader)) != 0 {
		t.Fatalf("a post that fails: %v, %d rows, %d comments; want the error, the five Shape rows and no comment", err, len(d.questions()), len(d.appComments(shapeQuestionsHeader)))
	}
	flaky.failPost, flaky.failAfter = false, true
	if _, err := d.r.shapeQuestions(ctx); err == nil || len(d.appComments(shapeQuestionsHeader)) != 1 || strings.Contains(keptRecord(t, d.w, "copies.tsv"), "601\t") {
		t.Fatalf("a read that fails after the post: %v, %d comments, copied: %v; want the error, the comment posted and not copied", err, len(d.appComments(shapeQuestionsHeader)), strings.Contains(keptRecord(t, d.w, "copies.tsv"), "601\t"))
	}
	d.restart()
	wait, err := d.r.shapeQuestions(ctx)
	if want := []string{"Q-002", "Q-003", "Q-004"}; err != nil || !slices.Equal(wait, want) || len(d.appComments(shapeQuestionsHeader)) != 1 || !strings.Contains(keptRecord(t, d.w, "copies.tsv"), "601\t") {
		t.Fatalf("the restart: wait %v, %v, %d comments, copied: %v; want %v, one comment, copied", wait, err, len(d.appComments(shapeQuestionsHeader)), strings.Contains(keptRecord(t, d.w, "copies.tsv"), "601\t"), want)
	}
	if body := d.appComments(shapeQuestionsHeader)[0]; !strings.Contains(body, "Q-007: The question Q-002\n") || !strings.Contains(body, "## For the Operator\n\nQ-006: The question Q-001\n") {
		t.Errorf("the comment of the round:\n%s\nwant the five questions, each in the block of its question", body)
	}
	// A second run writes nothing.
	head := recordsHead(t, d.w)
	d.restart()
	if _, err := d.r.shapeQuestions(ctx); err != nil || recordsHead(t, d.w) != head || len(d.appComments(shapeQuestionsHeader)) != 1 {
		t.Errorf("a second run: %v, the records head moved: %v, %d comments; want none", err, recordsHead(t, d.w) != head, len(d.appComments(shapeQuestionsHeader)))
	}
	// The answers: B and U parse, the hours do not.
	addComment(d.w, 1, 901, 303, "carol", "2030-01-01T00:00:00Z", "Q-007: 5.0\nQ-008: 10.0\nQ-009: soon\nQ-010: it works\n")
	addComment(d.w, 1, 902, 101, "pharzam", "2030-01-01T00:00:01Z", "Q-006: go\n")
	d.restart()
	wait, err = d.r.shapeQuestions(ctx)
	if want := []string{"Q-004"}; err != nil || !slices.Equal(wait, want) || d.has("budget.tsv") {
		t.Fatalf("the unparsed hours: wait %v, %v, budget.tsv: %v; want %v and no budget.tsv", wait, err, d.has("budget.tsv"), want)
	}
	want := map[string]string{"Q-001": "answered", "Q-002": "answered", "Q-003": "answered", "Q-004": "open", "Q-005": "answered", "Q-006": "answered", "Q-007": "answered", "Q-008": "answered",
		"Q-009": "answered", "Q-010": "answered", "Q-011": "asked"}
	if got := stateOf(d.questions()); !equalMaps(got, want) || len(d.appComments(shapeQuestionsHeader)) != 2 {
		t.Errorf("the states after the first answers = %v, %d comments; want %v and a second comment that asks the hours again", got, len(d.appComments(shapeQuestionsHeader)), want)
	}
	// The correction is the answer to the new row.
	addComment(d.w, 1, 903, 303, "carol", "2030-01-02T00:00:00Z", "Q-011: 40.0\n")
	d.restart()
	wait, err = d.r.shapeQuestions(ctx)
	rows, rerr := records.ReadBudget([]byte(keptRecord(t, d.w, "budget.tsv")))
	if wantRows := [][]string{{"band.B", "5.0", "Q-007"}, {"band.U", "10.0", "Q-008"}, {"appetite.hours", "40.0", "Q-011"}}; err != nil || rerr != nil || len(wait) != 0 || !slices.EqualFunc(rows, wantRows, slices.Equal) {
		t.Errorf("the correction: wait %v, %v, budget %v, %v; want no wait and %v", wait, err, rows, rerr, wantRows)
	}
	if got := stateOf(d.questions()); got["Q-004"] != "answered" || got["Q-011"] != "answered" || len(d.appComments(shapeQuestionsHeader)) != 2 {
		t.Errorf("after the correction Q-004 is %s, Q-011 is %s, %d comments; want answered, answered, 2", got["Q-004"], got["Q-011"], len(d.appComments(shapeQuestionsHeader)))
	}
}

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// ctlOptions leaves the output of the n-th options session for the scripted harness.
func (d *drvWorld) ctlOptions(n int, rows string) {
	d.ctlWrite("options-"+string(rune('0'+n)), "decision\toption\tconstraint\tresult\tevidence\n"+rows)
}

const (
	droppedDB  = "db\tpg\tInv-1\tfail\t—\ndb\tsqlite\tInv-1\tfail\t—\n"
	survivorDB = "db\tpg\tInv-1\tpass\tdocs/a.md\ndb\tsqlite\tInv-1\tpass\thttps://example.invalid/b\n"
)

func (d *drvWorld) sessionsStarted() int {
	b, _ := os.ReadFile(filepath.Join(d.ctl, "n-options"))
	n, _ := strconv.Atoi(strings.TrimSpace(string(b)))
	return n
}

// I2 (loop.md, The round, the rewrite and the stop): round 1 all dropped, one rewrite, all dropped again, then one comment, the row
// waiting for the Operator, and a second run that writes and starts nothing.
func TestTheRewriteAndTheRoundWithNoProgress(t *testing.T) {
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	d.settled()
	d.ctlOptions(1, droppedDB)
	d.ctlOptions(2, droppedDB)
	ctx := context.Background()
	got := d.r.shapeStep(ctx)
	if want := (Step{"shape", "done", "waiting for the Operator: shape: no survivor"}); got != want {
		t.Fatalf("the step = %+v; want %+v", got, want)
	}
	comments := d.appComments(noSurvivorHeader)
	if len(comments) != 1 || !strings.Contains(comments[0], "## db\n\nRound 1:\n- pg: Inv-1: fail\n- sqlite: Inv-1: fail\n\nRound 2, the rewrite:\n- pg: Inv-1: fail\n") {
		t.Errorf("the comments of the stop: %q; want one that names both rounds", comments)
	}
	if n := d.sessionsStarted(); n != 2 || !d.has("spec/options/round-1/db.tsv") || !d.has("spec/options/db.tsv") {
		t.Errorf("%d sessions, round-1 copy: %v; want two sessions and the copy of round 1", n, d.has("spec/options/round-1/db.tsv"))
	}
	if log := gitOut(t, d.w.bare, "log", "--format=%H", "layup-records", "--", "spec/options/round-1/db.tsv"); len(strings.Fields(log)) != 1 {
		t.Errorf("the copy of round 1 is in %d commits; want one, the rewrite's", len(strings.Fields(log)))
	}
	if d.w.fake.issues != 2 {
		t.Errorf("%d issues; want the two of Start, and no task opened on a stop", d.w.fake.issues)
	}
	// A restart gives the same row, writes nothing and starts no session.
	head := recordsHead(t, d.w)
	d.restart()
	if again := d.r.shapeStep(ctx); again != got || recordsHead(t, d.w) != head || len(d.appComments(noSurvivorHeader)) != 1 || d.sessionsStarted() != 2 {
		t.Errorf("a restart: %+v, head moved: %v, %d comments, %d sessions; want the same row, no write, one comment, two sessions", again, recordsHead(t, d.w) != head, len(d.appComments(noSurvivorHeader)), d.sessionsStarted())
	}
}

// I3: a stop during the rewrite keeps nothing; the restart runs the rewrite again, once, and the copy of round 1 is made once.
func TestAStopDuringTheRewriteRunsItAgainOnce(t *testing.T) {
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	d.settled()
	d.ctlOptions(1, droppedDB)
	d.ctlWrite("options-fail-2", "")
	d.ctlOptions(3, survivorDB)
	ctx := context.Background()
	if _, wait, err := d.r.optionsStage(ctx); err == nil || wait != "" || d.has("spec/options/round-1/db.tsv") {
		t.Fatalf("the rewrite that stops: wait %q, %v, round-1 copy: %v; want the error and no copy", wait, err, d.has("spec/options/round-1/db.tsv"))
	}
	d.restart()
	results, wait, err := d.r.optionsStage(ctx)
	if chosen := len(results["db"].Survivors); err != nil || wait != "" || chosen != 2 || d.sessionsStarted() != 3 {
		t.Fatalf("the restart: %d survivors, wait %q, %v, %d sessions; want two survivors, no wait, three sessions in all", chosen, wait, err, d.sessionsStarted())
	}
	if log := gitOut(t, d.w.bare, "log", "--format=%H", "layup-records", "--", "spec/options/round-1/db.tsv"); len(strings.Fields(log)) != 1 {
		t.Errorf("the copy of round 1 is in %d commits; want one", len(strings.Fields(log)))
	}
	if first := keptRecord(t, d.w, "spec/options/round-1/db.tsv"); !strings.Contains(first, "pg\tInv-1\tfail") {
		t.Errorf("the copy of round 1 = %q; want the dropped table", first)
	}
}

// I7, the demo (R11): after the Scaffold the step shape asks no question (none is left), writes and filters the options, and takes the
// specification task and then the architecture task to their verified heads, on a target with no PRD (#432) and on one with a PRD; a
// second run writes nothing.
func TestTheStepShape(t *testing.T) {
	for _, noPRD := range []bool{true, false} {
		d := newDrvWorldOpts(t, loopPerms, drvOpts{noPRD: noPRD, noOption: true})
		d.settled()
		d.ctlOptions(1, survivorDB)
		d.ctlWrite("commit", "")
		ctx := context.Background()
		got := d.r.shapeStep(ctx)
		if got.Result != "done" || !strings.Contains(got.Detail, " verified") {
			t.Fatalf("noPRD %v: the step = %+v; want done, with both tasks verified", noPRD, got)
		}
		tasks := d.register()
		if len(tasks) != 2 {
			t.Fatalf("noPRD %v: the step %+v; the register holds %d tasks; want 2", noPRD, got, len(tasks))
		}
		if want := (Step{"shape", "done", "1 decisions, each with a survivor; " + tasks[0].ID + " and " + tasks[1].ID + " verified"}); got != want {
			t.Fatalf("noPRD %v: the step = %+v; want %+v", noPRD, got, want)
		}
		spec, arch := tasks[0].ID, tasks[1].ID
		heads := d.boundOf(spec)
		if len(heads) != 1 {
			t.Fatalf("the specification is bound %d times; want once", len(heads))
		}
		var base string
		for _, e := range d.eventsOf(arch) {
			if e[1] == "attempt" {
				base = e[4]
			}
		}
		if base != heads[0] {
			t.Errorf("noPRD %v: the architecture task starts at %s; want the verified head %s of the specification", noPRD, base, heads[0])
		}
		for _, task := range []string{spec, arch} {
			var got []string
			for _, s := range statusesOf(d.w, d.boundOf(task)[0]) {
				got = append(got, s["context"].(string)+" "+s["state"].(string))
			}
			for _, w := range []string{GatesContext + " success", SpecContext + " success", RulesContext + " success", VerifyContext + " success"} {
				if !slices.Contains(got, w) {
					t.Errorf("noPRD %v: the statuses of %s = %v; want %q", noPRD, task, got, w)
				}
			}
		}
		if files := gitOut(t, d.r.target, "diff", "--name-only", heads[0], d.boundOf(arch)[0]); strings.Contains(files, "docs/prd/") {
			t.Errorf("noPRD %v: the files of the architecture task = %q; want none of docs/prd/", noPRD, files)
		}
		for _, p := range d.promptsOf(arch) {
			if strings.Contains(p, "The verifier") && !strings.Contains(p, heads[0]) {
				t.Errorf("noPRD %v: the prompt of a session of the architecture task does not name the base %s", noPRD, heads[0])
			}
		}
		// A second run writes nothing.
		head, posts, patches := recordsHead(t, d.w), d.requests("POST /repos/"), d.requests("PATCH ")
		d.restart()
		if again := d.r.shapeStep(ctx); again.Result != "done" || again.Detail != "done before" || recordsHead(t, d.w) != head || d.requests("POST /repos/") != posts || d.requests("PATCH ") != patches {
			t.Errorf("noPRD %v: a second run = %+v; head moved: %v, POST %d (was %d), PATCH %d (was %d); want done before and no write", noPRD, again, recordsHead(t, d.w) != head, d.requests("POST /repos/"), posts, d.requests("PATCH "), patches)
		}
	}
}

// The wait of the step is the row waiting: the intent form, the last row; with no budget.tsv and no intent question left to wait for,
// the step fails closed (loop.md, The questions left open: the wait).
func TestTheStepWaitsOnTheIntentFormAndFailsClosedWhenNothingIsOpen(t *testing.T) {
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	d.push(map[string][]byte{
		"questions.tsv":      cavTable(t, records.QuestionsSchema, [][]string{qrow("Q-001", "band.B", "open"), qrow("Q-002", "band.U", "open"), qrow("Q-003", "appetite.hours", "open")}),
		"intake/markers.tsv": cavTable(t, records.MarkersSchema, nil), "intake/answers.tsv": cavTable(t, records.AnswersSchema, nil),
	}, "an Intake that left the intent form open")
	ctx := context.Background()
	if got, want := d.r.shapeStep(ctx), (Step{"shape", "done", "waiting: the intent form: Q-001 Q-002 Q-003"}); got != want {
		t.Errorf("the step = %+v; want %+v", got, want)
	}
	if d.sessionsStarted() != 0 {
		t.Error("a session started before the intent form was answered; want none")
	}
	d.push(map[string][]byte{"questions.tsv": cavTable(t, records.QuestionsSchema, [][]string{qrow("Q-001", "band.B", "answered"), qrow("Q-002", "band.U", "answered"), qrow("Q-003", "appetite.hours", "answered")})}, "the questions answered, no budget")
	if got := d.r.shapeStep(ctx); got.Result != "fail" || got.Detail != "budget.tsv is not written and no question of the intent form is open" {
		t.Errorf("the step with no budget.tsv and nothing open = %+v; want the fail that names it", got)
	}
}
