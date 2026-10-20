//go:build integration

package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// targetFiles is the first tree of the target of the prompt tests: a PRD whose
// table holds REQ-005, a spec section with a heading that holds it and names
// two decision records, a decision record that no section names, and the rule
// files that the register of the world matches.
var targetFiles = map[string]string{
	"docs/prd/PRD-0001.md":        "# PRD\n\n| ID | Text |\n| -- | ---- |\n| REQ-005 | A handoff is a record. |\n",
	"docs/spec/loop.md":           "# Spec\n\n## REQ-005 — The loop\n\nSee ADR-0019 and ADR-0020.\n\n## Other\n\nADR-0021.\n",
	"docs/adr/0019-phase-loop.md": "adr 19\n",
	"docs/adr/0020-route.md":      "adr 20\n",
	"docs/adr/0021-smart-if.md":   "adr 21\n",
	"docs/rules/a.md":             "a rule\n",
	"AGENTS.md":                   "the agent rules\n",
	"docs/other.md":               "not a rule file\n",
}

// copyAll is the shell of a harness that commits its prompt file and its input files, and the result.
func copyAll() string {
	return `cp ../prompt.md seen-prompt.txt && mkdir seen-input && cp ../input/* seen-input/ && git add -A && git commit -qm seen && printf '` + octal(okResult) + `' > ../result/result.tsv`
}

// The demo of loop.md, the tier and the context of a task: the prompt file is
// built from the records of a real task (a plan posted, copied and recorded by
// the calls of #360, the register, the rule-path register) and the target at
// the base, and the session reads the files that it names, each at its hash.
func TestAPromptFileIsBuiltFromTheRecordsOfARealTaskAndTheSessionReadsTheFiles(t *testing.T) {
	t.Parallel()
	l := newLoop(t, targetFiles)
	ctx := context.Background()
	planID := l.plan(1, "claude", "the plan\n")
	_ = planID
	l.attempt(2)
	s := l.prod()
	pc, err := s.PromptFile(ctx, PromptIn{Task: "T-ab12", Step: "implement", Base: l.a.base, Findings: []Finding{{Rule: "gate:static", Text: "command failed"}}})
	if err != nil {
		t.Fatal(err)
	}
	want := map[session.Kind][]string{
		session.KindIssueAndPlan:       {"tasks.tsv", "copies/"},
		session.KindRequirements:       {"docs/prd/PRD-0001.md", "docs/spec/loop.md"},
		session.KindDecisionRecords:    {"docs/adr/0019-phase-loop.md", "docs/adr/0020-route.md"},
		session.KindHandoffsAndAnswers: nil,
		session.KindRuleFiles:          {"AGENTS.md", "docs/rules/a.md"},
	}
	for kind, paths := range want {
		got := itemsOf(pc, kind)
		if len(got) != len(paths) {
			t.Errorf("%s = %v, want %v", kind, got, paths)
			continue
		}
		for i, p := range paths {
			if !strings.HasPrefix(got[i], p) {
				t.Errorf("%s[%d] = %s, want %s", kind, i, got[i], p)
			}
		}
	}
	// A decision record that no section of the task names is no item; a file that is no rule path is none.
	for _, it := range pc.Items {
		if it.Path == "docs/adr/0021-smart-if.md" || it.Path == "docs/other.md" {
			t.Errorf("the item %s is listed", it.Path)
		}
	}
	spec := harness(t, l.a, copyAll())
	spec.Role, spec.Attempt, spec.Context = roleOfStep("implement"), 2, pc
	id, err := l.r.TaskSession(ctx, spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	head := l.bound(id)
	prompt := gitOut(t, l.r.Clone, "show", head+":seen-prompt.txt")
	root := filepath.Join(l.r.Host, "sessions", id)
	if abs, err := filepath.Abs(root); err == nil {
		root = abs
	}
	if prompt != string(pc.Text(root)) {
		t.Errorf("the prompt file that the session read differs from Text(root):\n%s\nwant:\n%s", prompt, pc.Text(root))
	}
	for _, it := range pc.Items {
		named := it.Path
		if it.Where == "records" {
			named = "/input/" + filepath.Base(it.Path) // a records item is read under input/
		}
		if !strings.Contains(prompt, it.SHA) || !strings.Contains(prompt, named) {
			t.Errorf("the prompt does not name %s with its sha256 %s", it.Path, it.SHA)
		}
		var data string
		if it.Where == "records" {
			data = gitOut(t, l.w.bare, "show", head+":seen-input/"+filepath.Base(it.Path))
		} else {
			data = gitOut(t, l.w.bare, "show", l.a.base+":"+it.Path)
		}
		if got := sumOf([]byte(data)); got != it.SHA {
			t.Errorf("%s %s: the bytes the session reads have the sha256 %s, the item says %s", it.Where, it.Path, got, it.SHA)
		}
	}
	if !strings.Contains(prompt, "gate:static") {
		t.Errorf("the findings are not in the prompt:\n%s", prompt)
	}
}

// The start row holds the size of that prompt; a context that does not fit
// refuses the start with both numbers (architecture.md §9).
func TestAPromptThatDoesNotFitRefusesTheStartWithBothNumbers(t *testing.T) {
	t.Parallel()
	l := newLoop(t, targetFiles)
	ctx := context.Background()
	l.attempt(2)
	pc, err := l.prod().PromptFile(ctx, PromptIn{Task: "T-ab12", Step: "implement", Base: l.a.base})
	if err != nil {
		t.Fatal(err)
	}
	spec := harness(t, l.a, copyAll())
	spec.Role, spec.Attempt, spec.Context = roleOfStep("implement"), 2, pc
	spec.Pair.Context = 10
	id, err := l.r.TaskSession(ctx, spec)
	if err == nil {
		t.Fatalf("a prompt over the context: no error (%s)", id)
	}
	root := filepath.Join(l.r.Host, "sessions", "S-00000000")
	if abs, e := filepath.Abs(root); e == nil {
		root = abs
	}
	est := (len(pc.Text(root)) + 3) / 4
	if !strings.Contains(keptRecord(t, l.w, "tasks/T-ab12/events.tsv"), "\trefused\t2\t"+id+"\t") || !strings.Contains(keptRecord(t, l.w, "tasks/T-ab12/events.tsv"), "context "+strconv.Itoa(est)+" 10") {
		t.Errorf("the events:\n%s\nwant a refusal context %d 10", keptRecord(t, l.w, "tasks/T-ab12/events.tsv"), est)
	}
}

// The tier and the position of a retry come from the plan and the records: no
// list is uphill; a frozen list with no open question is downhill; a question
// of the last attempt makes it uphill again until an answer is accepted.
func TestThePlanAndTheRecordsGiveTheTierOfARetry(t *testing.T) {
	t.Parallel()
	l := newLoop(t, nil)
	ctx := context.Background()
	s := l.prod()
	pos := func() Position {
		t.Helper()
		p, err := s.Position(ctx, "T-ab12")
		if err != nil {
			t.Fatal(err)
		}
		return p
	}
	large := Plan{Size: "large"}
	if p := pos(); p != Uphill || TierOf(large, p) != "reasoning" {
		t.Errorf("before the list: %s, %s; want uphill, reasoning", p, TierOf(large, p))
	}
	l.frozenAt("claude")
	if p := pos(); p != Downhill || TierOf(large, p) != "execution" || PlanTier(large) != "reasoning" {
		t.Errorf("a frozen list: %s, %s; want downhill, execution (the plan alone is reasoning)", p, TierOf(large, p))
	}
	got, err := ResolveStep("implement", TierOf(large, pos()), "")
	if err != nil || got.Tier != "execution" || got.Role != "software-developer" {
		t.Errorf("ResolveStep = %+v, %v", got, err)
	}
	// Attempt 2 ends with a question: uphill again.
	l.attempt(2)
	asking := "kind\tn\tvalue\tsha256\treason\nstatus\t1\tneeds_context\t—\tasks\nquestion\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n"
	l.dev(2, "devin", map[string]string{"q.txt": "q\n"}, asking)
	if p := pos(); p != Uphill {
		t.Errorf("an open question: %s; want uphill", p)
	}
	// The owner's answer, and attempt 3 that uses it: accepted, so downhill.
	owner := l.as("software-architect", 2, "codex", map[string]string{"o.txt": "o\n"}, "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tanswered\nanswer\t1\tYes\t—\tdocs/spec/forge.md\n")
	l.attempt(3)
	l.dev(3, "devin", map[string]string{"d.txt": "d\n"}, "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tuses the answer "+owner+".1\n")
	if p := pos(); p != Downhill {
		t.Errorf("an answer accepted by the next attempt: %s; want downhill", p)
	}
}

// The same ledger twice, differing only in the role of the session that wrote
// the answer row: the owner role of the question moves the task downhill, any
// other role leaves it uphill (loop.md, A question during the work, T-b3s3).
func TestAnAnswerOfAWrongRoleSessionLeavesThePositionUphill(t *testing.T) {
	t.Parallel()
	const ask = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tneeds_context\t—\tasks\n"
	const answer = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tanswered\nanswer\t1\tYes\t—\tdocs/spec/forge.md\n"
	for _, c := range []struct {
		name, role, questions string
		want                  Position
	}{
		{"the owner role of the question", "software-architect", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n", Downhill},
		{"a developer", "software-developer", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n", Uphill},
		{"the owner of another kind", "software-engineer", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n", Uphill},
		{"one of two owners of one session's questions", "software-architect", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\nquestion\t2\tWhere?\t—\tenvironment-or-infrastructure no docs/spec/forge.md\n", Uphill},
		{"the owner of two questions of one kind", "software-architect", "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\nquestion\t2\tHow?\t—\tinterface-contract no docs/spec/forge.md\n", Downhill},
	} {
		l := newLoop(t, nil)
		ctx := context.Background()
		s := l.prod()
		l.frozenAt("claude")
		l.attempt(2)
		l.dev(2, "devin", map[string]string{"q.txt": "q\n"}, ask+c.questions)
		owner := l.as(c.role, 2, "codex", map[string]string{"o.txt": "o\n"}, answer)
		l.attempt(3)
		l.dev(3, "devin", map[string]string{"d.txt": "d\n"}, "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tuses the answer "+owner+".1\n")
		if p, err := s.Position(ctx, "T-ab12"); err != nil || p != c.want {
			t.Errorf("%s: %s, %v; want %s", c.name, p, err, c.want)
		}
	}
}

// The records of a specification task and of an architecture task: the revised
// draft and its facts, and the options table of one decision.
const (
	draftTSV   = "id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\trequirement\tThe loop\t1\ta check\tS-aaaaaaaa\n"
	spansTSV   = "start\tend\tfact\tclass\treason\ttext\n0\t10\t1\tneed\t—\tthe need\n"
	answersTSV = "question\tanswer\tcomment\tauthor\ttime\nQ-001\tyes\t11\t1\t2026-01-01T00:00:00Z\n"
	testsTSV   = "test\tsource\tsha256\tharness\tfrozen\nt1\tinternal/a.go\t" + sha64 + "\tclaude\t2026-10-18T00:00:00Z\n"
)

var shapeFiles = map[string][]byte{"spec/draft.tsv": []byte(draftTSV), "spec/spans.tsv": []byte(spansTSV),
	"intake/answers.tsv": []byte(answersTSV), "spec/options/db.tsv": optionsTable}

var dbResult = map[string]spec.Result{"db": resultOfOptions([]string{"pg", "sqlite"}, "pg", "sqlite")}

// shapeWorld is newLoop with the target's facts record, a specification task
// T-sp01 and an architecture task T-ar01 in the register; the records of Shape
// are the caller's to add.
func shapeWorld(t *testing.T) *loopRun {
	t.Helper()
	files := map[string]string{"docs/facts/F-0003-answers.md": "the answers\n"}
	for p, c := range targetFiles {
		files[p] = c
	}
	l := newLoop(t, files)
	for _, task := range []Task{{ID: "T-sp01", Issue: 4, Kind: "specification", Requirements: []string{"REQ-005"}}, {ID: "T-ar01", Issue: 5, Kind: "architecture", Requirements: []string{"REQ-005"}}} {
		task.Opened = time.Now().UTC().Truncate(time.Second)
		reg, trans := []byte(keptRecord(t, l.w, "tasks.tsv")), []byte(keptRecord(t, l.w, "transitions.tsv"))
		out, err := registerFiles(reg, trans, task)
		if err == nil {
			err = l.r.commit(context.Background(), out, "tasks: register "+task.ID)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	return l
}

func (l *loopRun) records(files map[string][]byte) {
	l.t.Helper()
	if err := l.r.commit(context.Background(), files, "records"); err != nil {
		l.t.Fatal(err)
	}
}

// The demo of row 66s: for each of the seven task steps, a session started with
// the prompt file of a real records tree reads that file and the input files,
// each at its hash, and the file ends with the instruction of the step.
func TestEachStepsPromptFileFromARealRecordsTree(t *testing.T) {
	t.Parallel()
	l := shapeWorld(t)
	ctx := context.Background()
	l.records(shapeFiles)
	for _, c := range []struct{ step, task string }{{"plan-of-a-task", "T-sp01"}, {"plan-review", "T-ar01"}, {"specification", "T-sp01"}, {"shape-architecture", "T-ar01"},
		{"implement", "T-ab12"}, {"verification", "T-ab12"}, {"close-out", "T-ab12"}} {
		n, err := l.r.StartAttempt(ctx, c.task, l.a.base)
		if err != nil && !strings.Contains(err.Error(), "attempt") {
			t.Fatal(err)
		}
		in := PromptIn{Task: c.task, Step: c.step, Base: l.a.base}
		if c.task == "T-ar01" {
			in.Options = dbResult
		}
		pc, err := l.prod().PromptFile(ctx, in)
		if err != nil {
			t.Fatalf("%s: %v", c.step, err)
		}
		spec := harness(t, l.a, copyAll())
		spec.Pair.Context = 100000
		spec.Task, spec.Role, spec.Attempt, spec.Context = c.task, roleOfStep(c.step), n, pc
		id, err := l.r.TaskSession(ctx, spec)
		if err != nil {
			t.Fatalf("%s: TaskSession: %q, %v", c.step, id, err)
		}
		// The head is the one that the end of the session fetched into the run's clone: a step that only reads has no push, so no bound
		// (task T-jy6f).
		head := strings.TrimSpace(gitOut(t, l.r.Clone, "rev-parse", "refs/layup/sessions/"+id))
		root, _ := filepath.Abs(filepath.Join(l.r.Host, "sessions", id))
		prompt := gitOut(t, l.r.Clone, "show", head+":seen-prompt.txt")
		if prompt != string(pc.Text(root)) {
			t.Errorf("%s: the prompt that the session read differs from Text(root):\n%s", c.step, prompt)
		}
		if first := strings.SplitN(stepInstructions[c.step], "{", 2)[0]; !strings.Contains(prompt, "\nThe work of this step:\n"+first) {
			t.Errorf("the prompt of %s holds no line \"The work of this step:\" followed by its instruction", c.step)
		}
		for _, it := range pc.Items {
			if it.Where == "records" && sumOf([]byte(gitOut(t, l.r.Clone, "show", head+":seen-input/"+filepath.Base(it.Path)))) != it.SHA {
				t.Errorf("%s: the input %s is not at its hash", c.step, it.Path)
			}
		}
		if shape := c.task != "T-ab12" && c.step != "close-out"; shape != strings.Contains(prompt, "/input/draft.tsv ") {
			t.Errorf("%s of %s: the draft is in the prompt: %v", c.step, c.task, !shape)
		}
	}
}

// Each refusal of PromptFile has its control: the same call with the input made
// right goes through, so the refusal is the case's and no other read's.
func TestPromptFileRefusalsHaveTheirControls(t *testing.T) {
	t.Parallel()
	l := shapeWorld(t)
	ctx := context.Background()
	call := func(task, step string, opts map[string]spec.Result) error {
		_, err := l.prod().PromptFile(ctx, PromptIn{Task: task, Step: step, Base: l.a.base, Options: opts})
		return err
	}
	refused := func(name string, err error, want string) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: %v; want an error holding %q", name, err, want)
		}
	}
	passes := func(name string, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("%s (the control): %v", name, err)
		}
	}
	// Stage 0: a step with no instruction is refused before every read, even for a base that no read could resolve.
	// Expected change (row 66z, #426): answer-to-a-question has its instruction and leaves the refused steps.
	for _, step := range []string{"shape-activation-batch", "plan-of-a-milestone"} {
		_, err := l.prod().PromptFile(ctx, PromptIn{Task: "T-ab12", Step: step, Base: "0000000000000000000000000000000000000000"})
		if err == nil || err.Error() != "the step "+step+" has no instruction in M2e" {
			t.Errorf("the step %s: %v", step, err)
		}
	}
	_, err := l.prod().PromptFile(ctx, PromptIn{Task: "T-ab12", Step: "implement", Base: "0000000000000000000000000000000000000000"})
	if err == nil || strings.Contains(err.Error(), "has no instruction") {
		t.Errorf("the control of stage 0: %v; want the failure of a later read", err)
	}
	// Stage 1: the step and the kind, and the options.
	refused("specification of a build task", call("T-ab12", "specification", nil), "the step specification is not a step of a task of kind build")
	refused("implement of a specification task", call("T-sp01", "implement", nil), "the step implement is not a step of a task of kind specification")
	passes("plan-review of a build task", call("T-ab12", "plan-review", nil))
	refused("no options for the architecture", call("T-ar01", "shape-architecture", nil), "the options of the task T-ar01")
	refused("options for a build task", call("T-ab12", "implement", dbResult), "the options of the task T-ab12")
	refused("options for the close-out", call("T-ar01", "close-out", dbResult), "the options of the task T-ar01")
	// Stage 3: the Shape records, one by one, then a draft that its reader refuses.
	for _, p := range []string{"spec/draft.tsv", "spec/spans.tsv", "intake/answers.tsv"} { // the brief is in the records from Start
		refused("the missing "+p, call("T-sp01", "specification", nil), "the records have no "+p)
		l.records(map[string][]byte{p: shapeFiles[p]})
	}
	passes("the specification with its records", call("T-sp01", "specification", nil))
	// Each required record, removed from a records tree that has all of them, is a refusal that names it; restored, the same call passes;
	// a record that its reader refuses is an error that names it (the brief has no reader).
	store := l.r.Store.(*recordsStore)
	for _, p := range []string{"spec/draft.tsv", "spec/spans.tsv", "start/problem-statement.md", "intake/answers.tsv"} {
		orig := []byte(keptRecord(t, l.w, p))
		if err := store.Remove(ctx, []string{p}, "remove "+p); err != nil {
			t.Fatal(err)
		}
		refused("the removed "+p, call("T-sp01", "specification", nil), "the records have no "+p)
		l.records(map[string][]byte{p: orig})
		passes("the restored "+p, call("T-sp01", "specification", nil))
		if !strings.HasSuffix(p, ".md") {
			l.records(map[string][]byte{p: []byte("not a table\n")})
			refused("a refused "+p, call("T-sp01", "specification", nil), p+": ")
			l.records(map[string][]byte{p: orig})
		}
	}
	l.records(map[string][]byte{"spec/draft.tsv": []byte(strings.Replace(draftTSV, "requirement", "wish", 1))})
	refused("a draft with a bad kind", call("T-sp01", "specification", nil), "spec/draft.tsv: ")
	l.records(map[string][]byte{"spec/draft.tsv": []byte(draftTSV)})
	// The options: a table that the records lack, a result that disagrees, and two records of one name.
	refused("a missing table", call("T-ar01", "shape-architecture", dbResult), "the records have no spec/options/db.tsv")
	l.records(map[string][]byte{"spec/options/db.tsv": optionsTable, "spec/options/draft.tsv": optionsTable})
	passes("the architecture with its table", call("T-ar01", "shape-architecture", dbResult))
	refused("a renamed option", call("T-ar01", "shape-architecture", map[string]spec.Result{"db": resultOfOptions([]string{"pg", "lite"}, "pg")}), "decision db: the options of the result are")
	refused("a decision named like a record", call("T-ar01", "shape-architecture", map[string]spec.Result{"draft": dbResult["db"]}), "two records are named draft.tsv")
	// The frozen list: a malformed one is an error naming the path; a good one is an item.
	l.records(map[string][]byte{"tasks/T-sp01/tests.tsv": []byte("not a frozen list\n")})
	refused("a malformed frozen list", call("T-sp01", "specification", nil), "tasks/T-sp01/tests.tsv: ")
	l.records(map[string][]byte{"tasks/T-sp01/tests.tsv": []byte(testsTSV)})
	pc, err := l.prod().PromptFile(ctx, PromptIn{Task: "T-sp01", Step: "specification", Base: l.a.base})
	if err != nil || !strings.HasSuffix(strings.Join(itemsOf(pc, session.KindIssueAndPlan), " "), "tasks/T-sp01/tests.tsv") {
		t.Errorf("a good frozen list: %v, %v; want it after the plan", err, pc)
	}
	// The verifier: the base of the last attempt, and no attempt at all.
	refused("verification with no attempt", call("T-sp01", "verification", nil), "the task T-sp01 has no attempt")
	for i := 0; i < 2; i++ {
		if _, err := l.r.StartAttempt(ctx, "T-sp01", l.a.base); err != nil {
			t.Fatal(err)
		}
	}
	passes("verification after an attempt", call("T-sp01", "verification", nil))
	// The ledger: no sessions.tsv before a session, both after one, and a telemetry table its reader refuses.
	refused("close-out with no ledger", call("T-ab12", "close-out", nil), "the records have no sessions.tsv")
	l.attempt(2)
	spec := harness(t, l.a, copyAll())
	spec.Pair.Context = 100000
	if id, err := l.r.TaskSession(ctx, spec); err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	passes("close-out with the ledger", call("T-ab12", "close-out", nil))
	for _, p := range []string{"telemetry.tsv", "sessions.tsv"} { // the other table of the ledger is there
		orig := []byte(keptRecord(t, l.w, p))
		if err := store.Remove(ctx, []string{p}, "remove "+p); err != nil {
			t.Fatal(err)
		}
		refused("close-out without "+p, call("T-ab12", "close-out", nil), "the records have no "+p)
		l.records(map[string][]byte{p: orig})
		passes("close-out with "+p+" restored", call("T-ab12", "close-out", nil))
	}
	l.records(map[string][]byte{"telemetry.tsv": []byte("not telemetry\n")})
	refused("a refused telemetry table", call("T-ab12", "close-out", nil), "telemetry.tsv: ")
	l.records(map[string][]byte{"sessions.tsv": []byte("not sessions\n")}) // refused by the read of stage 2, before the ledger is read
	refused("a refused sessions table", call("T-ab12", "close-out", nil), "sessions.tsv: ")
}

// The harness reads only prompt.md: it writes the comment file in the form that
// the text shows, with the placeholders filled, and the result file from the rows
// that the text quotes; the handoff is then valid, the verifier's by the head that
// code holds for its start row.
func TestAHarnessThatFollowsTheInstructionGivesValidHandoffs(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		step, file, kind, from, to, item, verdict string
	}{
		{"plan-of-a-task", PlanFile, "plan", "task plan", "plan review", "plan-fields", ""},
		{"plan-review", PlanReviewFile, "plan-review", "plan review", "developer", "verdict-field", ""},
		{"verification", ReviewRecordFile, "review-record", "verifier", "close-out", "nothing-material", "nothing material in scope"},
	} {
		t.Run(c.step, func(t *testing.T) {
			l := newLoop(t, targetFiles)
			ctx := context.Background()
			l.attempt(2)
			pc, err := l.prod().PromptFile(ctx, PromptIn{Task: "T-ab12", Step: c.step, Base: l.a.base})
			if err != nil {
				t.Fatal(err)
			}
			text := string(pc.Text("/x"))
			if !strings.Contains(text, "result/"+c.file+"<TAB>") {
				t.Errorf("the prompt of %s does not name result/%s", c.step, c.file)
			}
			comment := string(shown(formLines(text), c.verdict, c.step == "plan-of-a-task")) // a build task lists a test; a clean record lists no finding
			result := followedResult(t, text, "completed", commentSum(comment), false)
			spec := harness(t, l.a, "printf '"+octal(comment)+"' > ../result/"+c.file+" && printf '"+octal(string(result))+"' > ../result/result.tsv")
			spec.Pair.Context = 100000
			spec.Role, spec.Attempt, spec.Context = roleOfStep(c.step), 2, pc
			id, err := l.r.TaskSession(ctx, spec)
			if err != nil {
				t.Fatalf("TaskSession: %q, %v", id, err)
			}
			r2 := afterRestart(l.r)
			p, err := r2.PostComment(ctx, "T-ab12", id, c.file)
			if err != nil {
				t.Fatalf("PostComment: %v", err)
			}
			v, err := r2.RecordHandoff(ctx, HandoffInput{Task: "T-ab12", Session: id, From: c.from, To: c.to, File: c.file, Kind: c.kind, Items: map[string]bool{c.item: true}})
			if want := fmt.Sprintf("valid comment %d-1", p.ID); err != nil || !v.Valid || v.Detail != want {
				t.Errorf("the handoff %s -> %s of a harness that followed the instruction = %+v, %v; want %s", c.from, c.to, v, err, want)
			}
			if _, err := records.ReadResult(result); err != nil {
				t.Error(err)
			}
		})
	}
}

// The blocks of loop.md are the text of record of the instructions (as probe-prompt
// is of the probe): each constant is its block, byte for byte.
func TestTheInstructionBlocksOfLoopMd(t *testing.T) {
	t.Parallel()
	doc, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]string{"instruction-result": resultInstruction, "instruction-answers": answersInstruction}
	for step, text := range stepInstructions {
		want["instruction-"+step] = text
	}
	if len(want) != 10 {
		t.Errorf("%d texts in the code; want the eight steps, the shared result text and the conditional answers text", len(want))
	}
	// Review condition 6 of #426: a block of the document that the code does not know is as wrong as a constant that differs.
	for _, line := range strings.Split(string(doc), "\n") {
		if name, ok := strings.CutPrefix(line, "```text instruction-"); ok {
			if _, known := want["instruction-"+name]; !known {
				t.Errorf("loop.md has the block instruction-%s and the code has no text for it", name)
			}
		}
	}
	for name, text := range want {
		_, block, ok := strings.Cut(string(doc), "\n```text "+name+"\n")
		block, _, _ = strings.Cut(block, "```\n")
		if !ok || block != text {
			t.Errorf("loop.md has the block %s: %v; it is equal to the code: %v", name, ok, block == text)
		}
	}
}

// loop.md, The specification output: the specimen that follows the forms of the
// instruction is a specification that waits for its bet, by the pinned prd-lint of the target.
func TestTheSpecimenOfTheSpecificationInstructionIsAcceptedByThePinnedPrdLint(t *testing.T) {
	t.Parallel()
	script, err := filepath.Abs("../../docs/prd/prd-lint.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		drop int
		want bool
	}{{"the specimen", -1, true}, {"the specimen without its matrix", 3, false}} {
		root := t.TempDir()
		for p, text := range map[string]string{"docs/prd/PRD-0001-x.md": specimen(t, c.drop), "docs/facts/F-0003-answers.md": "# the record\n"} {
			os.MkdirAll(filepath.Dir(filepath.Join(root, p)), 0o755)
			os.WriteFile(filepath.Join(root, p), []byte(text), 0o644)
		}
		ok, lines, err := spec.WaitingSpecAccepted(script, filepath.Join(root, "docs", "prd"))
		if err != nil || ok != c.want {
			t.Errorf("%s: accepted %v, %v, %v; want %v", c.name, ok, lines, err, c.want)
		}
	}
}

// A harness that follows the text of the steps that write no comment, with the
// outputs that the text requires, committed and hashed: the result crosses the real
// session boundary, for a finished step and for a step that a question stops.
func TestAHarnessThatFollowsTheInstructionGivesValidResultsForTheOtherStepsAndAStoppingQuestion(t *testing.T) {
	t.Parallel()
	l := shapeWorld(t)
	ctx := context.Background()
	l.records(shapeFiles)
	for _, c := range []struct{ step, task string }{{"specification", "T-sp01"}, {"shape-architecture", "T-ar01"}, {"implement", "T-ab12"}, {"close-out", "T-ab12"}} {
		for _, status := range []string{"completed", "needs_context"} {
			t.Run(c.step+"/"+status, func(t *testing.T) {
				n, err := l.r.StartAttempt(ctx, c.task, l.a.base)
				if err != nil {
					t.Fatal(err)
				}
				in := PromptIn{Task: c.task, Step: c.step, Base: l.a.base}
				if c.task == "T-ar01" {
					in.Options = dbResult
				}
				pc, err := l.prod().PromptFile(ctx, in)
				if err != nil {
					t.Fatal(err)
				}
				text := string(pc.Text("/x"))
				outputs := outputsOf(text)
				if c.step == "close-out" && len(outputs) != 2 {
					t.Fatalf("the close-out text names the outputs %v; want its task file and the completed log", outputs)
				}
				result := followedResult(t, text, status, sha64, status == "needs_context")
				var sh strings.Builder
				for _, o := range outputs {
					sh.WriteString("mkdir -p '" + filepath.Dir(o) + "' && printf '" + octal(outputContent(o)) + "' > '" + o + "' && ")
				}
				sh.WriteString("git add -A && git commit -qm work && printf '" + octal(string(result)) + "' > ../result/result.tsv")
				spec := harness(t, l.a, sh.String())
				spec.Pair.Context = 100000
				spec.Task, spec.Role, spec.Attempt, spec.Context = c.task, roleOfStep(c.step), n, pc
				id, err := l.r.TaskSession(ctx, spec)
				if err != nil {
					t.Fatalf("TaskSession: %q, %v", id, err)
				}
				kept := keptRecord(t, l.w, "tasks/"+c.task+"/results/"+id+".tsv")
				rows, err := records.ReadResult([]byte(kept))
				if kept != string(result) || err != nil {
					t.Fatalf("the committed result differs from the followed one (%v):\n%s", err, kept)
				}
				if err := checkOutputs(rows, outputs); err != nil {
					t.Error(err)
				}
				if qs := QuestionsOf(rows); (len(qs) == 1) != (status == "needs_context") {
					t.Errorf("questions %+v for the status %s", qs, status)
				}
				if events := keptRecord(t, l.w, "tasks/"+c.task+"/events.tsv"); strings.Contains(events, "\trefused\t"+strconv.Itoa(n)+"\t"+id+"\t") {
					t.Errorf("the session %s was refused:\n%s", id, events)
				}
			})
		}
	}
}
