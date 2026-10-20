//go:build integration

package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// The step bet on the production path (task T-2vva, #366, row 66j): the world of a run after Shape (66u's shapedWorldWith), the
// production constructor state.sessions(), real git, the httptest forge and a scripted harness for the brief. A restart is a fresh
// state on the same records. The tests of the stops replace the package variable betStopAt and so stay serial (guardrails.md, the
// serial rule of T-05f4).

const (
	// betOptions: db has two survivors, web has one and one drop; both survivors of db name REQ-001 among their rows.
	betOptionsTable = "db\tpg\tInv-1\tpass\tdocs/a.md\ndb\tpg\tREQ-001\tpass\tdocs/a.md\ndb\tsqlite\tInv-1\tpass\thttps://example.invalid/b\ndb\tsqlite\tREQ-001\tpass\thttps://example.invalid/b\n" +
		"web\tgo\tInv-1\tpass\tdocs/a.md\nweb\tgo\tREQ-001\tpass\tdocs/a.md\nweb\trust\tInv-1\tfail\t—\nweb\trust\tREQ-001\tpass\tdocs/a.md\n"
	// the answer: a priority, a fact out of scope, an option line for each decision (the one for web, which had one survivor, settles nothing).
	betAnswerText = "REQ-001: Must, phase 1\nfact 2: out of scope\noption db: pg\noption web: go\nbet\n"
)

var errStopped = errors.New("stopped by the test")

func betWorld(t *testing.T, options string) *drvWorld {
	t.Helper()
	d := shapedWorldWith(t, options, specimen(t, -1)+"\n**Non-goals.** None.\n", "80")
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	return d
}

func (d *drvWorld) bet() Step { return d.r.betStep(context.Background()) }

func (d *drvWorld) reply(id int, author int, login, created, body string) {
	addComment(d.w, 1, id, author, login, created, body)
}

// specBranch gives the branch of the specification task's verified head.
func (d *drvWorld) specBranch() (spec, branch string) {
	tasks := d.register()
	return tasks[0].ID, "task/" + tasks[0].ID + "/1"
}

// branchHead gives the head of a branch of the target, "" for none.
func (d *drvWorld) branchHead(branch string) string {
	out, err := gitRunOK(d.w.bare, "rev-parse", "--verify", "-q", "refs/heads/"+branch)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

func stopAt(stage string) func() {
	old := betStopAt
	betStopAt = func(s string) error {
		if s == stage {
			return errStopped
		}
		return nil
	}
	return func() { betStopAt = old }
}

type statusWatch struct {
	forge.Forge
	on func(sha, context string)
}

func (s statusWatch) SetStatus(ctx context.Context, sha, c, state, description string) error {
	s.on(sha, c)
	return s.Forge.SetStatus(ctx, sha, c, state, description)
}

// I1: the block bets equals its Go schema, and the Go schema has the columns of the block.
func TestTheBetsBlockEqualsItsGoSchema(t *testing.T) {
	t.Parallel()
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tsv.Compare(blocks["bets"], BetsSchema); err != nil {
		t.Errorf("the block bets and BetsSchema: %v", err)
	}
}

// I14, the demo (R11), with I2, I4, I7 and I11: after Shape the brief is posted once and the run waits; the comments that answer nothing
// change nothing; the answer of the approver leads to the version, the rendered commit with its carry, the architecture task of the
// one decision that had two survivors and, last, the row; a second run writes nothing.
func TestTheStepBet(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	ctx := context.Background()
	d.reply(6500, 303, "carol", "2026-10-08T12:00:00Z", betAnswerText) // before the brief: it answers nothing
	got := d.bet()
	if got.Result != "done" || !strings.HasPrefix(got.Detail, "waiting: the answer of the bet on the brief ") {
		t.Fatalf("the step = %+v; want the wait for the answer", got)
	}
	brief := keptRecord(t, d.w, "bets/1/brief.md")
	posted := d.appComments("## The problem")
	if len(posted) != 1 || posted[0] != brief || !strings.Contains(brief, "- pg (a survivor)\n- sqlite (a survivor)\n") || d.briefSessions() != 1 {
		t.Fatalf("%d comments of the brief, %d sessions; want the kept brief posted once, with the two survivors unmarked", len(posted), d.briefSessions())
	}
	d.reply(7001, 303, "carol", "2026-10-08T14:00:00Z", "REQ-001: Must, phase 1\n") // a finding: no verdict
	d.reply(7002, 707, "dave", "2026-10-08T14:01:00Z", betAnswerText)               // no role
	d.reply(7003, 101, "pharzam", "2026-10-08T14:02:00Z", betAnswerText)            // the Operator
	d.restart()
	if again := d.bet(); again != got || d.briefSessions() != 1 || len(d.appComments("## The problem")) != 1 {
		t.Fatalf("a second run = %+v, %d sessions, %d posts; want the same wait, one session, one post", again, d.briefSessions(), len(d.appComments("## The problem")))
	}
	d.reply(7004, 303, "carol", "2026-10-08T14:03:00Z", betAnswerText)
	d.reply(7005, 303, "carol", "2026-10-08T14:04:00Z", betAnswerText) // a later valid answer is not the answer
	d.restart()
	specID, branch := d.specBranch()
	h := d.boundOf(specID)[0]
	var carry []string
	d.r.cfg.Forge = statusWatch{d.r.cfg.Forge, func(sha, c string) {
		if c != VerifyContext || sha != d.branchHead(branch) {
			return
		}
		for _, e := range d.eventsOf(specID) {
			if e[1] == "bound" && e[3] == "" && e[5] == sha {
				carry = append(carry, sha)
			}
		}
	}}
	got = d.bet()
	if got.Result != "done" || !strings.HasPrefix(got.Detail, "the bet 1: ") {
		t.Fatalf("the step = %+v; want the bet made", got)
	}
	// The row, last.
	bets, err := ReadBets([]byte(keptRecord(t, d.w, "bets.tsv")))
	if err != nil || len(bets) != 1 {
		t.Fatalf("bets.tsv: %v, %d rows", err, len(bets))
	}
	row, sum := bets[0], digest([]byte(brief))
	if row[betMilestone] != "1" || row[betKind] != "first" || row[betBriefSHA] != sum || row[betUSD] != "5.0" || row[betHours] != "12.0" || row[betAnswerID] != "7004" ||
		row[betOptions] != "db=pg web=go" || row[betVerdict] != "bet" || row[betApprover] != "303" {
		t.Errorf("the row = %q; want the brief %s, the cap 5.0 and 12.0, the answer 7004, the options db=pg web=go, bet, the approver 303", row, sum)
	}
	if subject := gitOut(t, d.w.bare, "log", "-1", "--format=%s", "layup-records"); !strings.Contains(subject, "the bet 1") || !strings.Contains(gitOut(t, d.w.bare, "show", "--name-only", "--format=", "layup-records"), "bets.tsv") {
		t.Errorf("the last records commit is %q; want the one that writes bets.tsv", subject)
	}
	// The version of the inventory: the class of a fact of the answers is that of its quote, the mark out is the answer's.
	idx, err := spec.ReadInventoryIndex([]byte(keptRecord(t, d.w, "spec/inventory.tsv")))
	if err != nil || len(idx) != 1 || idx[0][2] != "1" || idx[0][1] != digest([]byte(keptRecord(t, d.w, "spec/inventory-1.tsv"))) {
		t.Errorf("the index = %v, %v; want one version for the bet 1 with the hash of its bytes", idx, err)
	}
	inv, _ := spec.ReadInventory([]byte(keptRecord(t, d.w, "spec/inventory-1.tsv")))
	if len(inv) < 3 || inv[0][0] != "1" || inv[0][1] != "need" || inv[1][2] != "out" || inv[2][0] != "Q-002" || inv[2][3] != "7004" {
		t.Errorf("the version = %v; want fact 1 need in, fact 2 out, then the facts of the answers, each of the copy 7004", inv)
	}
	// The rendered commit: one child of the verified head, deterministic, with its carry after the bound.
	r := d.branchHead(branch)
	if parents := strings.Fields(gitOut(t, d.w.bare, "rev-list", "--parents", "-n1", r)); len(parents) != 2 || parents[1] != h {
		t.Fatalf("the branch %s is at %s with the parents %v; want a child of the verified head %s", branch, r, parents, h)
	}
	if names := strings.Fields(gitOut(t, d.w.bare, "diff", "--name-only", h, r)); !slices.Equal(names, []string{"docs/facts/F-0004-numbered-facts.md", "docs/prd/PRD-0001-x.md"}) {
		t.Errorf("the rendered commit changes %v; want the facts record F-0004 and the PRD", names)
	}
	if facts := gitOut(t, d.w.bare, "show", r+":docs/facts/F-0004-numbered-facts.md"); !strings.HasPrefix(facts, "1. **need:** s1\n2. **measure (out of scope):** s2\n") {
		t.Errorf("the facts record = %q", facts)
	}
	if prd := gitOut(t, d.w.bare, "show", r+":docs/prd/PRD-0001-x.md"); !strings.Contains(prd, "| REQ-001 | Alpha | Must | 1 |") || !strings.Contains(prd, "Out of scope by the bet: 2.") {
		t.Errorf("the PRD at the rendered commit lacks the priority or the mark out of scope")
	}
	var verify []string
	for _, s := range statusesOf(d.w, r) {
		if s["context"] == VerifyContext {
			verify = append(verify, fmt.Sprint(s["state"], " ", s["description"]))
		}
	}
	if !slices.Equal(verify, []string{"success verified head " + h + " plus the rendered commit"}) || !slices.Equal(carry, []string{r}) {
		t.Errorf("layup/verify on the rendered commit = %q, with the bound seen at the post %q; want one success, posted after the bound", verify, carry)
	}
	ev := d.eventsOf(specID)
	if n := len(ev); ev[n-2][1] != "push" || ev[n-1][1] != "bound" || ev[n-2][3] != "" || ev[n-1][3] != "" || ev[n-1][5] != r || ev[n-2][6] != branch {
		t.Errorf("the last events of %s = %v; want the push and the bound of %s with no session", specID, ev[max(0, n-2):], r)
	}
	// I11: the full check's record of the problem statement is the rendered one.
	in, err := d.r.sessions().specInput(ctx, strings.TrimSpace(gitOut(t, d.r.target, "rev-parse", "refs/remotes/origin/main")), r, nil)
	if err != nil || !slices.Equal(in.PSBRecords, []string{"F-0004"}) {
		t.Errorf("PSBRecords at the rendered commit = %v, %v; want [F-0004]", in.PSBRecords, err)
	}
	// The architecture task of the one decision that had two survivors, with its binding.
	tasks := d.register()
	if len(tasks) != 3 {
		t.Fatalf("the register holds %d tasks; want the two Shape tasks and one for the decision db", len(tasks))
	}
	a := tasks[2]
	if a.Kind != "architecture" || !slices.Equal(a.Requirements, []string{"REQ-001"}) || !slices.Equal(a.Predecessors, []string{tasks[0].ID, tasks[1].ID}) {
		t.Errorf("the new row = %+v; want an architecture task for REQ-001 after the two Shape tasks", a)
	}
	binding := keptRecord(t, d.w, bindingPath(a.ID, "db"))
	if binding != "option\tconstraint\tresult\tevidence\npg\tInv-1\tpass\tdocs/a.md\npg\tREQ-001\tpass\tdocs/a.md\n" {
		t.Errorf("the binding = %q; want the rows of the option pg", binding)
	}
	if issue := d.w.fake.issueTexts[a.Issue-1]; issue[0] != "LAYUP task "+a.ID+": the architecture of the decision db" || !strings.Contains(issue[1], "Chosen option: pg") {
		t.Errorf("the issue %d = %q", a.Issue, issue)
	}
	// A second run writes nothing.
	head, posts, patches := recordsHead(t, d.w), d.requests("POST /repos/"), d.requests("PATCH ")
	d.restart()
	if again := d.bet(); again.Result != "done" || again.Detail != "done before" || recordsHead(t, d.w) != head || d.requests("POST /repos/") != posts || d.requests("PATCH ") != patches {
		t.Errorf("a second run = %+v; head moved %v, POST %d (was %d); want done before and no write", again, recordsHead(t, d.w) != head, d.requests("POST /repos/"), posts)
	}
}

// I3: brief.lines unset is a wait for the Operator and starts no session.
func TestBriefLinesUnsetWaitsWithNoSession(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	d.setParameter("brief.lines", "")
	got := d.bet()
	if got.Result != "done" || got.Detail != "waiting for the Operator: brief.lines is not set" || d.briefSessions() != 0 || len(d.appComments("## The problem")) != 0 {
		t.Errorf("the step = %+v after %d sessions; want the wait for the Operator, no session and no post", got, d.briefSessions())
	}
}

// I4: an approver that the idea owner names is the only one who answers, and a lookup that fails chooses nobody.
func TestOnlyTheApproverOfTheBetAnswersAndALookupThatFailsChoosesNoOne(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	d.bet()
	var copies, approvers [][]string
	if d.has("copies.tsv") {
		copies, _ = records.ReadCopies([]byte(keptRecord(t, d.w, "copies.tsv")))
	}
	approvers, _ = records.ReadApprovers([]byte(keptRecord(t, d.w, "approvers.tsv")))
	body := "bet: dave\n"
	copies = append(copies, []string{"6400", "1", "1", "303", "carol", "", "2026-10-08T12:45:00Z", "2026-10-08T12:46:00Z", digest([]byte(body)), "copies/6400-1.md"})
	approvers = append(approvers, []string{"707", "approver", "dave", "2026-10-08T12:45:00Z", "6400"})
	pushRecords(t, d.w.bare, map[string]string{"copies.tsv": string(cavTable(t, records.CopiesSchema, copies)), "copies/6400-1.md": body,
		"approvers.tsv": string(cavTable(t, records.ApproversSchema, approvers))}, "the idea owner names the approver of the bet")
	d.reply(7010, 303, "carol", "2026-10-08T14:00:00Z", betAnswerText) // the idea owner, with an approver named
	d.restart()
	if got := d.bet(); got.Result != "done" || !strings.HasPrefix(got.Detail, "waiting: the answer") {
		t.Fatalf("with dave named: %+v; want the idea owner's answer to answer nothing", got)
	}
	d.reply(7011, 707, "dave", "2026-10-08T14:01:00Z", betAnswerText)
	d.restart()
	d.r.cfg.Forge = userFails{d.r.cfg.Forge}
	if got := d.bet(); got.Result != "fail" || !strings.Contains(got.Detail, "lookup down") || d.has("bets.tsv") {
		t.Fatalf("with the lookup down: %+v; want a fail, and nobody chosen", got)
	}
	d.restart()
	if got := d.bet(); got.Result != "done" || !strings.HasPrefix(got.Detail, "the bet 1: ") {
		t.Fatalf("with dave named: %+v; want the bet made on dave's answer", got)
	}
	if bets, _ := ReadBets([]byte(keptRecord(t, d.w, "bets.tsv"))); len(bets) != 1 || bets[0][betAnswerID] != "7011" || bets[0][betApprover] != "707" {
		t.Errorf("bets.tsv = %v; want the answer 7011 of the approver 707", bets)
	}
}

type userFails struct{ forge.Forge }

func (userFails) UserID(context.Context, string) (int64, error) { return 0, errors.New("lookup down") }

// I9: a decision whose chosen option names no requirement is refused before the inventory, the rendering and the tasks are written,
// at every run.
func TestADecisionWhoseOptionNamesNoRequirementIsRefusedBeforeAnyWrite(t *testing.T) {
	t.Parallel()
	d := betWorld(t, survivorDB)
	d.bet()
	d.reply(7020, 303, "carol", "2026-10-08T14:00:00Z", "REQ-001: Must, phase 1\nfact 2: out of scope\noption db: pg\nbet\n")
	d.restart()
	specID, branch := d.specBranch()
	for i := 0; i < 2; i++ {
		got := d.bet()
		if got.Result != "fail" || got.Detail != "the decision db: its option pg names no requirement" {
			t.Fatalf("run %d: %+v; want the refusal", i, got)
		}
		if i == 0 {
			d.restart()
		}
	}
	if d.has("spec/inventory.tsv") || d.has("bets.tsv") || d.branchHead(branch) != d.boundOf(specID)[0] || len(d.register()) != 2 || renderedBound(d.eventsOf(specID), 1) || slices.ContainsFunc(d.eventsOf(specID), func(e []string) bool { return e[1] == "push" && e[3] == "" }) {
		t.Errorf("the refusal left a write: inventory %v, bets %v, branch %q, %d tasks", d.has("spec/inventory.tsv"), d.has("bets.tsv"), d.branchHead(branch), len(d.register()))
	}
}

// I10: a no bet writes its row with no appetite and waits at every run; nothing else is written.
func TestANoBetWritesItsRowAndWaits(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	d.bet()
	d.reply(7030, 303, "carol", "2026-10-08T14:00:00Z", "REQ-001: Must, phase 1\nno bet\n")
	d.restart()
	want := Step{"bet", "done", "waiting for the Operator: no bet"}
	if got := d.bet(); got != want {
		t.Fatalf("the step = %+v; want %+v", got, want)
	}
	bets, err := ReadBets([]byte(keptRecord(t, d.w, "bets.tsv")))
	if err != nil || len(bets) != 1 || bets[0][betVerdict] != "no-bet" || bets[0][betUSD] != "" || bets[0][betHours] != "" || bets[0][betOptions] != "" {
		t.Fatalf("bets.tsv = %v, %v; want one no-bet row with no appetite", bets, err)
	}
	specID, branch := d.specBranch()
	head := recordsHead(t, d.w)
	d.restart()
	if got := d.bet(); got != want || recordsHead(t, d.w) != head || d.has("spec/inventory.tsv") || d.branchHead(branch) != d.boundOf(specID)[0] || len(d.register()) != 2 {
		t.Errorf("a second run = %+v; want the same wait and no inventory, no rendered commit, no task", got)
	}
}

// A later bet is built by no row of M2e: the step fails closed and writes nothing.
func TestALaterBetFailsClosed(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	row := strings.Join([]string{"9", "first", "901", digest([]byte("x")), "5.0", "12.0", "905", "—", "bet", "303", "2026-10-08T14:00:00Z"}, "\t")
	pushRecords(t, d.w.bare, map[string]string{"bets.tsv": "milestone\tkind\tbrief\tbrief_sha256\tappetite_usd\tappetite_hours\tanswer\toptions\tverdict\tapprover\tdecided\n" + row + "\n"}, "a bet of another milestone")
	d.restart()
	head := recordsHead(t, d.w)
	if got := d.bet(); got.Result != "fail" || !strings.Contains(got.Detail, "row 66t") || recordsHead(t, d.w) != head || d.briefSessions() != 0 {
		t.Errorf("the step = %+v; want a fail naming row 66t and no write", got)
	}
}

// I8, with the fact-class change and the stops of condition 5: a stop after each act that one can follow leaves a run that the
// next one finishes, with one push of the rendered commit, one post of layup/verify, one task for the decision, no new Shape session,
// and the row of bets.tsv last. The stop between the issue and its row opens a second issue, which is reported and not hidden.
func TestTheBetStepIsAtomicToItsLastAct(t *testing.T) {
	for _, stage := range []string{"version", "pushevent", "pushed", "bound", "carry", "issue", "task"} {
		t.Run(stage, func(t *testing.T) {
			d := betWorld(t, betOptionsTable)
			d.bet()
			answer := betAnswerText
			if stage == "version" { // a fact whose class the answer changes: the inventory now differs from the Shape basis
				answer = strings.Replace(answer, "fact 2: out of scope\n", "fact 1: constraint\nfact 2: out of scope\n", 1)
			}
			d.reply(7040, 303, "carol", "2026-10-08T14:00:00Z", answer)
			d.restart()
			sessions, issues := len(d.sessionRows()), d.requests("POST /repos/acme/target/issues")
			restore := stopAt(stage)
			got := d.bet()
			restore()
			if got.Result != "fail" || !strings.Contains(got.Detail, errStopped.Error()) {
				t.Fatalf("the stopped step = %+v; want the stop", got)
			}
			if d.has("bets.tsv") {
				t.Fatalf("bets.tsv is written after a stop at %s", stage)
			}
			// The event push announces the forge write: after the stop at it the branch is still at the verified head, and after the stop at the push it is at the commit with no bound.
			specID, branch := d.specBranch()
			code := func(kind string) bool {
				return slices.ContainsFunc(d.eventsOf(specID), func(e []string) bool { return e[1] == kind && e[3] == "" })
			}
			switch head := d.branchHead(branch); stage {
			case "pushevent":
				if !code("push") || head != d.boundOf(specID)[0] {
					t.Errorf("after the stop at the push event: the event %v, the branch at %s; want the event and the verified head", code("push"), head)
				}
			case "pushed":
				if !code("push") || code("bound") || head == d.boundOf(specID)[0] {
					t.Errorf("after the stop at the push: push %v, bound %v, the branch at %s; want the push, no bound and the new head", code("push"), code("bound"), head)
				}
			}
			d.restart()
			// The restart runs shape first: the code push without its bound must not start an attempt or a session.
			if shape := d.r.shapeStep(context.Background()); shape.Result != "done" || shape.Detail != "done before" {
				t.Fatalf("shape after the stop = %+v; want done before", shape)
			}
			if got = d.bet(); got.Result != "done" || !strings.HasPrefix(got.Detail, "the bet 1: ") {
				t.Fatalf("the step after the stop = %+v; want the bet made", got)
			}
			var pushes, bounds int
			for _, e := range d.eventsOf(specID) {
				if e[3] == "" && e[1] == "push" {
					pushes++
				} else if e[3] == "" && e[1] == "bound" {
					bounds++
				} else if e[1] == "attempt" && e[2] != "1" {
					t.Errorf("an attempt %s started", e[2])
				}
			}
			r, verifies := d.branchHead(branch), 0
			for _, s := range statusesOf(d.w, r) {
				if s["context"] == VerifyContext {
					verifies++
				}
			}
			wantIssues := 1
			if stage == "issue" {
				wantIssues = 2 // the known window between OpenIssue and the row: the first issue stays on the forge
			}
			if pushes != 1 || bounds != 1 || verifies != 1 || len(d.register()) != 3 || d.requests("POST /repos/acme/target/issues")-issues != wantIssues || len(d.sessionRows()) != sessions {
				t.Errorf("pushes %d, bounds %d, verifies %d, tasks %d, issues %d (want %d), new sessions %d; want one of each, one task and no session", pushes, bounds, verifies, len(d.register()), d.requests("POST /repos/acme/target/issues")-issues, wantIssues, len(d.sessionRows())-sessions)
			}
			if inv := keptRecord(t, d.w, "spec/inventory-1.tsv"); stage == "version" && !strings.Contains(inv, "1\tconstraint\tin\t7040") {
				t.Errorf("the version = %q; want fact 1 constraint", inv)
			}
			if subject := gitOut(t, d.w.bare, "log", "-1", "--format=%s", "layup-records"); !strings.Contains(subject, "the bet 1") {
				t.Errorf("the last records commit is %q; want the row of bets.tsv last", subject)
			}
		})
	}
}

// I12: after the rendered push the step shape is done before: no verifier session, no attempt, no write, though the pull request's
// head is the rendered commit.
func TestShapeStaysDoneAfterTheRenderedPush(t *testing.T) {
	t.Parallel()
	d := betWorld(t, betOptionsTable)
	d.bet()
	d.reply(7050, 303, "carol", "2026-10-08T14:00:00Z", betAnswerText)
	d.restart()
	if got := d.bet(); got.Result != "done" || !strings.HasPrefix(got.Detail, "the bet 1: ") {
		t.Fatalf("the step = %+v", got)
	}
	// The link of the review record, asked directly, accepts the pull request whose head is the rendered commit (the guard of LinkReviewRecord).
	specID, _ := d.specBranch()
	st, err := d.r.sessions().verifiedStep(context.Background(), specID)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.r.sessions().LinkReviewRecord(context.Background(), specID, st.Session); err != nil {
		t.Errorf("LinkReviewRecord with the rendered commit as the head of the pull request: %v; want none", err)
	}
	head, posts, patches, sessions := recordsHead(t, d.w), d.requests("POST /repos/"), d.requests("PATCH "), len(d.sessionRows())
	d.restart()
	if got := d.r.shapeStep(context.Background()); got.Result != "done" || got.Detail != "done before" || recordsHead(t, d.w) != head || d.requests("POST /repos/") != posts || d.requests("PATCH ") != patches || len(d.sessionRows()) != sessions {
		t.Errorf("shape after the rendered push = %+v; want done before and no write, no session", got)
	}
}

// loop.md, the issue text of an architecture task of a decision: the constant is its block, byte for byte.
func TestTheArchitectureIssueBlockOfLoopMd(t *testing.T) {
	t.Parallel()
	doc, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	_, block, ok := strings.Cut(string(doc), "\n```text bet-issue-architecture\n")
	block, _, _ = strings.Cut(block, "```\n")
	if !ok || block != archIssueBlock {
		t.Errorf("loop.md has the block bet-issue-architecture: %v; it is equal to the code: %v", ok, block == archIssueBlock)
	}
}
