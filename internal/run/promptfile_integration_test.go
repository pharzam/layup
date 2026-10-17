//go:build integration

package run

import (
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/session"
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
	prompt := gitOut(t, l.w.bare, "show", head+":seen-prompt.txt")
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
