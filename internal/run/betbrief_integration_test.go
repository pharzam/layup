//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// The brief of the first bet on the production path (task T-h65e, #418, row 66u): the world of a run after Shape (the two Shape tasks
// at their verified heads), the production constructor state.sessions(), real git, the httptest forge and a scripted harness for the
// step bet-brief. A restart is a fresh state on the same records.

const briefProblem, briefHoles = "The problem in prose.\n", "No rabbit holes but one.\n"

func briefAppetite(usd string) string { return "appetite: USD " + usd + ", 12 hours\nWhy: small.\n" }

// ctlBrief leaves the three files of the n-th session of the step bet-brief for the scripted harness, and, for a cost, the money that
// its usage report gives.
func (d *drvWorld) ctlBrief(n int, problem, appetite, holes, cost string) {
	d.t.Helper()
	dir := filepath.Join(d.ctl, fmt.Sprintf("brief-%d", n))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		d.t.Fatal(err)
	}
	for name, text := range map[string]string{"problem.md": problem, "appetite.md": appetite, "rabbit-holes.md": holes} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
			d.t.Fatal(err)
		}
	}
	if cost != "" {
		d.ctlWrite(fmt.Sprintf("cost-%d", n), cost)
	}
}

// briefSessions gives the number of sessions of the step bet-brief that the harness ran, and the prompt of the n-th.
func (d *drvWorld) briefSessions() int {
	b, _ := os.ReadFile(filepath.Join(d.ctl, "n-brief"))
	var n int
	fmt.Sscan(strings.TrimSpace(string(b)), &n)
	return n
}

func (d *drvWorld) briefPrompt(n int) string {
	b, _ := os.ReadFile(filepath.Join(d.ctl, fmt.Sprintf("brief-prompt-%d", n)))
	return string(b)
}

// setParameter sets a row of parameters.tsv as a set of the Operator would.
func (d *drvWorld) setParameter(name, value string) {
	d.t.Helper()
	rows, err := ReadParameters([]byte(keptRecord(d.t, d.w, "parameters.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	for i := range rows {
		if rows[i].Name == name {
			rows[i].Value, rows[i].Comment, rows[i].Evidence = value, 6200000001, "the idea owner reads that many lines"
			if value == "" { // unset: the row of the table, with no set
				rows[i].Comment, rows[i].Evidence = 0, ""
			}
		}
	}
	var b bytes.Buffer
	if err := WriteParameters(&b, rows); err != nil {
		d.t.Fatal(err)
	}
	d.push(map[string][]byte{"parameters.tsv": b.Bytes()}, "a set of "+name)
}

// shapedWorld is a run after Shape: both Shape tasks at their verified heads, db with two survivors, brief.lines 60, and the harness of
// the claude row reporting its cost.
func shapedWorld(t *testing.T) *drvWorld { return shapedWorldWith(t, survivorDB, "", "60") }

// shapedWorldWith is shapedWorld with the options table of the architect, the PRD that the specification session writes (the specimen
// of the instruction when empty) and the value of brief.lines.
func shapedWorldWith(t *testing.T, options, prd, lines string) *drvWorld {
	t.Helper()
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	if prd != "" {
		d.ctlWrite("prd.md", prd)
	}
	// The clock of the world runs 600 times faster than the real one, so the 8 hours of the cap of milestone 0 are 48 seconds, and a host under
	// load would end a test at that cap by the clock. The hours are widened; the money cap stays 50.
	start := keptRecord(t, d.w, "start/start.tsv")
	if !strings.Contains(start, "intake.cap\t50.0,8.0\t") {
		t.Fatalf("start.tsv holds no intake.cap of 50.0,8.0:\n%s", start)
	}
	d.push(map[string][]byte{"start/start.tsv": []byte(strings.Replace(start, "intake.cap\t50.0,8.0\t", "intake.cap\t50.0,4000.0\t", 1))}, "the hours of the cap of milestone 0 widened")
	d.settled()
	d.ctlOptions(1, options)
	d.ctlWrite("commit", "")
	if got := d.r.shapeStep(context.Background()); got.Result != "done" || !strings.Contains(got.Detail, " verified") {
		t.Fatalf("the step shape = %+v; want both tasks verified", got)
	}
	d.setParameter("brief.lines", lines)
	d.harnesses[0][regBilling], d.harnesses[0][regUsage] = "api", "claude-result"
	d.restart()
	return d
}

func (d *drvWorld) sessionRows() [][]string {
	d.t.Helper()
	rows, err := records.ReadSessions([]byte(keptRecord(d.t, d.w, "sessions.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	return rows
}

func (d *drvWorld) candidates() []string {
	d.t.Helper()
	out := gitOut(d.t, d.w.bare, "ls-tree", "--name-only", "layup-records", "bets/1/attempts/")
	return strings.Fields(out)
}

// I1, the demo (R11): one session of the step bet-brief, its input files equal to the records, the brief and its basis kept in one
// commit, the roles and the candidate tied by the session ID.
func TestTheBriefIsKeptFromOneSessionOfTheRun(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	before := len(d.sessionRows())
	brief, done, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || done || len(brief) == 0 {
		t.Fatalf("betBrief = %d bytes, done %v, %v; want the bytes of bets/1/brief.md", len(brief), done, err)
	}
	if kept := keptRecord(t, d.w, "bets/1/brief.md"); kept != string(brief) {
		t.Errorf("bets/1/brief.md =\n%s\nwant the bytes that betBrief gave\n%s", kept, brief)
	}
	rows := d.sessionRows()
	if len(rows) != before+1 || d.briefSessions() != 1 {
		t.Fatalf("%d sessions after the call (was %d), the harness ran %d; want one more and one run", len(rows), before, d.briefSessions())
	}
	id := rows[len(rows)-1][sessionID]
	if rows[len(rows)-1][sessionRole] != "product-owner" || !strings.HasPrefix(d.briefPrompt(1), "You write three parts of the brief of the bet of milestone 1, for LAYUP.") {
		t.Errorf("the new session %v with the prompt %q; want the role product-owner and the prompt of bet-brief", rows[len(rows)-1], d.briefPrompt(1))
	}
	if c := d.candidates(); !slices.Equal(c, []string{"bets/1/attempts/" + id + ".md"}) || keptRecord(t, d.w, c[0]) != string(brief) {
		t.Errorf("the candidates %v; want the one of the session %s with the bytes of the brief", c, id)
	}
	for name, path := range map[string]string{"problem-statement.md": "start/problem-statement.md", "spans.tsv": "spec/spans.tsv", "draft.tsv": "spec/draft.tsv", "questions.tsv": "questions.tsv", "answers.tsv": "intake/answers.tsv", "budget.tsv": "budget.tsv"} {
		got, err := os.ReadFile(filepath.Join(d.ctl, "brief-input-1", name))
		if want := keptRecord(t, d.w, path); err != nil || string(got) != want {
			t.Errorf("input/%s = %q, %v; want the bytes of %s", name, got, err, path)
		}
	}
	reqs, err := os.ReadFile(filepath.Join(d.ctl, "brief-input-1", "requirements.md"))
	if want := "### The requirements\n- REQ-001: Alpha Criterion: a check\n### The success criteria\n- faster\n"; err != nil || string(reqs) != want {
		t.Errorf("input/requirements.md = %q, %v; want the words of the PRD at the verified head %q (the draft says The loop)", reqs, err, want)
	}
	tasks := d.register()
	for _, want := range []string{"## The problem\n" + string(reqs) + briefProblem + "## The appetite\n" + briefAppetite("5") + "## The solution\n### The verified heads\n",
		fmt.Sprintf("- %s (specification): [pull request #%d](https://github.com/acme/target/pull/%d), head %s\n", tasks[0].ID, tasks[0].PR, tasks[0].PR, d.boundOf(tasks[0].ID)[0]),
		"### Decision db\n- pg (a survivor)\n- sqlite (a survivor)\n", "### The questions still open\n- none\narchitecture change: no\n" + briefHoles + "## The facts\n- 1 (need): s1\n- 2 (measure): s2\n- Q-002 (context): 5.0\n"} {
		if !strings.Contains(string(brief), want) {
			t.Errorf("the brief lacks %q:\n%s", want, brief)
		}
	}
	basis, err := ReadBetBasis([]byte(keptRecord(t, d.w, "bets/1/basis.tsv")))
	if err != nil || len(basis) != 10 {
		t.Errorf("bets/1/basis.tsv = %d rows, %v; want 10", len(basis), err)
	}
	if n := strings.Count(keptRecord(t, d.w, "telemetry.tsv"), id+"\t"); n != 1 {
		t.Errorf("telemetry.tsv holds %d rows of the session %s; want 1", n, id)
	}
}

// I2: a refused brief is tried again with its finding, and the refused session's telemetry stands.
func TestARefusedBriefIsTriedAgainWithItsFinding(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, "Text\n## The solution\n", briefAppetite("5"), briefHoles, "")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || d.briefSessions() != 2 {
		t.Fatalf("betBrief = %d bytes, %v after %d sessions; want a brief after two", len(brief), err, d.briefSessions())
	}
	if p := d.briefPrompt(2); !strings.Contains(p, "The check refused the last brief: file: problem.md has a line that starts with \"## \"\nWrite the three files again.\nThen end.\n") || strings.Contains(d.briefPrompt(1), "refused the last") {
		t.Errorf("the second prompt:\n%s\nwant the finding, and none in the first", p)
	}
	rows := d.sessionRows()
	if c := d.candidates(); len(c) != 1 || !strings.HasSuffix(c[0], rows[len(rows)-1][sessionID]+".md") || strings.Count(keptRecord(t, d.w, "telemetry.tsv"), "\n"+rows[len(rows)-2][sessionID]+"\t") != 1 {
		t.Errorf("candidates %v, telemetry of the first session: %v; want the second session's candidate alone, and the first one's telemetry row", c, keptRecord(t, d.w, "telemetry.tsv"))
	}
}

// Condition 3 of the plan review: a missing output file ends the session no-result, not refused file; its telemetry stands, the call
// returns the error, and the next run starts a new session. A semantic refusal is the case of the test above.
func TestAMissingOutputFileEndsNoResultAndTheNextRunStartsAgain(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	d.ctlWrite("brief-drop-1", "appetite.md")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
	_, _, err := d.r.betBrief(context.Background(), "1", "first")
	var pe *PhaseError
	if !errors.As(err, &pe) || !strings.Contains(pe.Reason, "no-result: the output file appetite.md is missing") || d.briefSessions() != 1 {
		t.Fatalf("betBrief = %v after %d sessions; want a *PhaseError no-result for the missing file after one", err, d.briefSessions())
	}
	rows := d.sessionRows()
	if n := strings.Count(keptRecord(t, d.w, "telemetry.tsv"), rows[len(rows)-1][sessionID]+"\t"); n != 1 || len(d.candidates()) != 0 {
		t.Errorf("the telemetry rows of the session: %d, candidates %v; want 1 and none", n, d.candidates())
	}
	d.restart()
	if brief, _, err := d.r.betBrief(context.Background(), "1", "first"); err != nil || len(brief) == 0 || d.briefSessions() != 2 {
		t.Errorf("the next run = %d bytes, %v after %d sessions; want a brief from a second session", len(brief), err, d.briefSessions())
	}
}

// I3 (condition 4 of the plan review of #366): the appetite is judged after the telemetry of the brief session. The money of the world
// before the call is unknown, so attempt 1 (caps 8 plus its own 3.0 against U 10) is refused only through its own reported spend.
func TestTheAppetiteIsJudgedAfterTheBriefSessionsTelemetry(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	if known, err := ledger.KnownMoney(mustTelemetry(t, d)); err != nil || known != "0.0" {
		t.Fatalf("the known money before the call = %q, %v; want 0.0 (a green negative: the 8 alone fits U 10)", known, err)
	}
	d.ctlBrief(1, briefProblem, briefAppetite("8"), briefHoles, "3.0")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "1.0")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || d.briefSessions() != 2 || !strings.Contains(string(brief), "appetite: USD 5, 12 hours") {
		t.Fatalf("betBrief = %q, %v after %d sessions; want the brief of attempt 2 after two", brief, err, d.briefSessions())
	}
	if p := d.briefPrompt(2); !strings.Contains(p, "The check refused the last brief: money: the money caps 8 plus the milestone-0 spend 3 are over band.U 10\n") {
		t.Errorf("the second prompt:\n%s\nwant the finding of the appetite with the session's own 3", p)
	}
	if c := d.candidates(); len(c) != 2 {
		t.Errorf("candidates %v; want two, the refused one standing as a record", c)
	}
}

func mustTelemetry(t *testing.T, d *drvWorld) [][]string {
	t.Helper()
	rows, err := records.ReadTelemetry([]byte(keptRecord(t, d.w, "telemetry.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// I4: a stop between the candidate and the promotion (the control comment of the session's end fails, as endPhase returns it after the
// candidate's commit) leaves the candidate; a fresh state promotes it with no new session and writes the brief and its basis in one
// commit. The promotion write is r.write, through Fenced: a lost lease is a *LostError and promotes nothing.
func TestAStopBeforeThePromotionPromotesWithNoSession(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	d.r.cfg.Forge = &flakyForge{Forge: d.r.cfg.Forge, failPost: true}
	if brief, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || brief != nil || d.has("bets/1/brief.md") || len(d.candidates()) != 1 {
		t.Fatalf("a stop at the comment: %d bytes, %v, brief.md %v, candidates %v; want the error, no brief and the candidate", len(brief), err, d.has("bets/1/brief.md"), d.candidates())
	}
	d.restart()
	candidate := keptRecord(t, d.w, d.candidates()[0])
	brief, done, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || done || string(brief) != candidate || d.briefSessions() != 1 {
		t.Errorf("the restart = %d bytes, done %v, %v after %d sessions; want the candidate's bytes and no new session", len(brief), done, err, d.briefSessions())
	}
	if keptRecord(t, d.w, "bets/1/brief.md") != candidate || !d.has("bets/1/basis.tsv") || strings.Count(gitOut(t, d.w.bare, "log", "--format=%H", "layup-records", "--", "bets/1/brief.md", "bets/1/basis.tsv"), "\n") != 1 {
		t.Error("the brief and its basis are not one commit of the candidate's bytes")
	}
}

// I5: a kept brief is done before: no write, no forge write, no session.
func TestAKeptBriefIsDoneBeforeWithNoWrite(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	first, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil {
		t.Fatal(err)
	}
	d.setParameter("brief.lines", "")
	d.setPrecondition()
	head, posts, patches := recordsHead(t, d.w), d.requests("POST /repos/"), d.requests("PATCH ")
	brief, done, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || !done || string(brief) != string(first) || recordsHead(t, d.w) != head || d.requests("POST /repos/") != posts || d.requests("PATCH ") != patches || d.briefSessions() != 1 {
		t.Errorf("a kept brief: %d bytes, done %v, %v; head moved %v, POST %d (was %d), PATCH %d (was %d), %d sessions; want done, the same bytes and no write",
			len(brief), done, err, recordsHead(t, d.w) != head, d.requests("POST /repos/"), posts, d.requests("PATCH "), patches, d.briefSessions())
	}
}

// setPrecondition makes every precondition of a fresh build fail, so that a kept brief shows it comes first: brief.lines unset and a
// bets.tsv present.
func (d *drvWorld) setPrecondition() {
	d.t.Helper()
	d.push(map[string][]byte{"bets.tsv": []byte("x\n")}, "a bets.tsv, as a later row would write it")
}

// I6: brief.lines unset starts no session, and nothing is written.
func TestBriefLinesUnsetStartsNoSession(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.operatorRows("Use Go.", "The question Q-006") // an Intake that leaves something out of the brief: the annex is prepared only after the waits
	d.setParameter("brief.lines", "")
	head := recordsHead(t, d.w)
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); !errors.Is(err, ErrBriefLinesUnset) {
		t.Errorf("err = %v; want ErrBriefLinesUnset", err)
	}
	if d.briefSessions() != 0 || recordsHead(t, d.w) != head || d.has("bets/1/facts.md") || d.has("bets/1/open.md") {
		t.Errorf("%d sessions, head moved: %v, facts.md %v, open.md %v; want none (a green negative beside the error)", d.briefSessions(), recordsHead(t, d.w) != head, d.has("bets/1/facts.md"), d.has("bets/1/open.md"))
	}
}

// The parts of code over brief.lines wait before any session, and a fresh build refuses a bets.tsv that exists.
func TestAFreshBuildRefusesWhatItCannotBuildBefore(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.operatorRows("Use Go.", "The question Q-006") // both annex files would be needed: the wait writes neither
	d.setParameter("brief.lines", "10")
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); !errors.Is(err, ErrBriefOverLimit) || d.briefSessions() != 0 || d.has("bets/1/facts.md") || d.has("bets/1/open.md") {
		t.Errorf("brief.lines 10: err = %v after %d sessions, facts.md %v, open.md %v; want ErrBriefOverLimit, none and no annex", err, d.briefSessions(), d.has("bets/1/facts.md"), d.has("bets/1/open.md"))
	}
	d.setParameter("brief.lines", "60")
	if err := d.r.store.Remove(context.Background(), []string{"start/problem-statement.md"}, "no problem statement"); err != nil {
		t.Fatal(err)
	}
	head := recordsHead(t, d.w)
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || !strings.Contains(err.Error(), "the records have no start/problem-statement.md") || d.briefSessions() != 0 || recordsHead(t, d.w) != head {
		t.Errorf("no problem-statement.md: err = %v after %d sessions, head moved %v; want the missing path refused, no session and no write", err, d.briefSessions(), recordsHead(t, d.w) != head)
	}
	d.setPrecondition()
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || !strings.Contains(err.Error(), "bets.tsv") || d.briefSessions() != 0 {
		t.Errorf("a bets.tsv: err = %v after %d sessions; want a refusal that names bets.tsv and none", err, d.briefSessions())
	}
	if _, _, err := d.r.betBrief(context.Background(), "1", "later"); err == nil || d.briefSessions() != 0 {
		t.Errorf("kind later: err = %v; want a refusal", err)
	}
}

// I7: the cap of milestone 0 ends the retry: attempt 1 is refused and reports 60 USD, above the cap of 50, so attempt 2 does not start.
func TestTheCapOfMilestoneZeroEndsTheRetry(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, "Text\n## The solution\n", briefAppetite("5"), briefHoles, "60.0")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
	_, _, err := d.r.betBrief(context.Background(), "1", "first")
	rows := d.sessionRows()
	if !errors.Is(err, ledger.ErrMilestoneCap) || d.briefSessions() != 1 || strings.Count(keptRecord(t, d.w, "telemetry.tsv"), rows[len(rows)-1][sessionID]+"\t") != 1 {
		t.Errorf("err = %v after %d sessions; want ErrMilestoneCap after one, with its one telemetry row", err, d.briefSessions())
	}
}

// I8: the links are the pull requests and the changed files of the verified heads; a pull request whose head moved is an error before
// any session.
func TestTheLinkSetIsTheChangedFilesAndPullRequestsOfTheVerifiedHeads(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	tasks := d.register()
	spHead := d.boundOf(tasks[0].ID)[0]
	blobOf := func(head, path string) string { return "https://github.com/acme/target/blob/" + head + "/" + path }
	d.ctlBrief(1, "See [the spec]("+blobOf(spHead, "AGENTS.md")+").\n", briefAppetite("5"), briefHoles, "")
	d.ctlBrief(2, "See [the spec]("+blobOf(spHead, "docs/spec/alpha.md")+").\n", briefAppetite("5"), briefHoles, "")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || d.briefSessions() != 2 || !strings.Contains(string(brief), "[the spec]("+blobOf(spHead, "docs/spec/alpha.md")+")") {
		t.Fatalf("betBrief = %v after %d sessions; want attempt 2 kept, attempt 1 refused for the link to a file the head did not change", err, d.briefSessions())
	}
	if p := d.briefPrompt(2); !strings.Contains(p, "link: the link to "+blobOf(spHead, "AGENTS.md")+" does not resolve") {
		t.Errorf("the second prompt lacks the finding of the link:\n%s", p)
	}
	// A pull request whose head moved is an error before any session.
	m := shapedWorld(t)
	m.w.fake.movePull(m.register()[0].PR, strings.Repeat("e", 40))
	m.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	if _, _, err := m.r.betBrief(context.Background(), "1", "first"); err == nil || !strings.Contains(err.Error(), "head") || m.briefSessions() != 0 {
		t.Errorf("a pull request whose head moved: err = %v after %d sessions; want an error naming the head, before any session", err, m.briefSessions())
	}
}

// A Shape task whose position is not verified is an error that names it, before any session: the brief reads the records and drives nothing.
func TestAShapeTaskWithNoVerifiedHeadIsAnErrorBeforeAnySession(t *testing.T) {
	t.Parallel()
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	d.settled()
	specTask, _ := d.open()
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	head := recordsHead(t, d.w)
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || !strings.Contains(err.Error(), specTask+": Shape has no verified head") || d.briefSessions() != 0 || recordsHead(t, d.w) != head {
		t.Errorf("betBrief on an unverified task = %v after %d sessions; want an error that names %s, no session and no write", err, d.briefSessions(), specTask)
	}
}

// I9 (loop.md, The bet: two survivors give a brief with no recommendation): code presents both survivors with no mark, and a session
// that names an option is refused with its finding.
func TestTwoSurvivorsGiveABriefWithNoRecommendation(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), "I would use pg here.\n", "")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || d.briefSessions() != 2 || !strings.Contains(string(brief), "- pg (a survivor)\n- sqlite (a survivor)\n") || strings.Contains(string(brief), "I would use pg") {
		t.Fatalf("betBrief = %v after %d sessions; want the second attempt kept, with both survivors unmarked", err, d.briefSessions())
	}
	if p := d.briefPrompt(2); !strings.Contains(p, "option: the brief names the option pg of db outside code's part") {
		t.Errorf("the second prompt lacks the finding of the option word:\n%s", p)
	}
}

// I10 (loop.md, The bet): a question still open is shown in the rabbit holes, with the text of the question.
func TestAQuestionStillOpenIsShownInTheRabbitHoles(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	qs := append(d.questions(), []string{"Q-005", "intent", "success-criteria", "idea-owner", "before-delivery", "intent", "", "", "The question Q-005", "2026-10-08T12:30:00Z", "layup run", "", "open"})
	data, err := table(records.QuestionsSchema, qs)
	if err != nil {
		t.Fatal(err)
	}
	d.push(map[string][]byte{"questions.tsv": data}, "a question left open")
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if want := "## The rabbit holes and the no-gos\n### The questions still open\n- Q-005: The question Q-005\narchitecture change: no\n"; err != nil || !strings.Contains(string(brief), want) {
		t.Errorf("betBrief = %v:\n%s\nwant the section %q", err, brief, want)
	}
}

// Conditions 1 and 4: the draft changed after the tasks were opened is an error naming the difference, and the basis is read back in a
// fresh state with the original requirements, survivors, facts, heads and appetite, after the facts' classes changed and another
// architecture task appeared, with nothing filtered again.
func TestTheBasisRoundTripsThroughBriefBasisInAFreshState(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err != nil {
		t.Fatal(err)
	}
	tasks := d.register()
	// The inventory class changes, another architecture task appears, and the options are filtered differently now.
	spans := strings.Replace(keptRecord(t, d.w, "spec/spans.tsv"), "\tneed\t", "\tconstraint\t", 1)
	d.push(map[string][]byte{"spec/spans.tsv": []byte(spans), "spec/options/db.tsv": []byte("option\tconstraint\tresult\tevidence\npg\tInv-1\tfail\t—\nsqlite\tInv-1\tfail\t—\n")}, "a class changed, the options refiltered")
	extra := Task{ID: "T-zz99", Issue: 99, Kind: "architecture", Requirements: tasks[0].Requirements, Opened: tasks[0].Opened}
	if err := d.r.registerTask(context.Background(), extra); err != nil {
		t.Fatal(err)
	}
	d.restart()
	got, err := d.r.briefBasis(context.Background(), "1")
	want := BriefBasis{Requirements: []string{"REQ-001"}, Survivors: map[string][]string{"db": {"pg", "sqlite"}},
		Facts: []spec.Fact{{Number: "1", Class: "need"}, {Number: "2", Class: "measure"}, {Number: "Q-002", Class: "context"}, {Number: "Q-003", Class: "context"}, {Number: "Q-004", Class: "context"}},
		Heads: map[string]string{tasks[0].ID: d.boundOf(tasks[0].ID)[0], tasks[1].ID: d.boundOf(tasks[1].ID)[0]}, Appetite: spec.Cap{USD: "5", Hours: "12"}}
	if err != nil || fmt.Sprint(got.Requirements, got.Survivors, got.Facts, got.Heads, got.Appetite, got.Open) != fmt.Sprint(want.Requirements, want.Survivors, want.Facts, want.Heads, want.Appetite, []string(nil)) {
		t.Errorf("briefBasis = %+v, %v\nwant %+v", got, err, want)
	}
	// A draft that lists another requirement than the register row is refused before any session.
	d2 := shapedWorld(t)
	draft := keptRecord(t, d2.w, "spec/draft.tsv") + "R-002\trequirement\tBeta\t1\tanother check\tS-aaaaaaaa\n"
	d2.push(map[string][]byte{"spec/draft.tsv": []byte(draft)}, "a draft with another requirement")
	if _, _, err := d2.r.betBrief(context.Background(), "1", "first"); err == nil || !strings.Contains(err.Error(), "REQ-002") || d2.briefSessions() != 0 {
		t.Errorf("a draft with REQ-002 over the register row: err = %v after %d sessions; want REQ-002 named, before any session", err, d2.briefSessions())
	}
}

// I11: the block bet-basis equals its Go schema.
func TestTheBetBasisBlockEqualsItsGoSchema(t *testing.T) {
	t.Parallel()
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tsv.Compare(blocks["bet-basis"], BetBasisSchema); err != nil {
		t.Error(err)
	}
}

// I12: the two text blocks of loop.md are the constants of the code, byte for byte (beside TestTheTextBlocksOfRunMd).
func TestTheBetBriefBlocksOfLoopMd(t *testing.T) {
	t.Parallel()
	loop, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{"bet-brief": briefPrompt, "bet-brief-refused": briefRefused} {
		_, text, ok := strings.Cut(string(loop), "\n```text "+name+"\n")
		text, _, _ = strings.Cut(text, "```\n")
		if !ok || text != want {
			t.Errorf("the block %s of loop.md:\n%s\nwant\n%s", name, text, want)
		}
	}
}

// F1 of round 1: a run stopped between the candidate and the promotion, and the fixed input changes before the recovery: the old candidate
// is not promoted over the new basis; a new session writes the brief again with the finding. (The control is the test of the stop above.)
func TestAChangedFixedInputBeforeTheRecoveryStartsANewAttempt(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
	d.r.cfg.Forge = &flakyForge{Forge: d.r.cfg.Forge, failPost: true}
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || len(d.candidates()) != 1 {
		t.Fatalf("a stop at the comment: %v, candidates %v; want the error and the candidate", err, d.candidates())
	}
	old := keptRecord(t, d.w, d.candidates()[0])
	lines := strings.Split(keptRecord(t, d.w, "spec/spans.tsv"), "\n")
	i := slices.IndexFunc(lines, func(l string) bool { return strings.Contains(l, "\tneed\t") })
	lines[i] += " (revised)"
	d.push(map[string][]byte{"spec/spans.tsv": []byte(strings.Join(lines, "\n"))}, "a fact's text changed before the recovery")
	brief, _, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || string(brief) == old || d.briefSessions() != 2 || !strings.Contains(d.briefPrompt(2), "fixed:") || !strings.Contains(string(brief), " (revised)") {
		t.Errorf("the recovery = %d bytes, %v after %d sessions; want a new brief of the new text and the finding fixed: in the second prompt", len(brief), err, d.briefSessions())
	}
}

// F2 of round 1: the reader of the basis refuses a kept brief that is malformed (the first call is the valid control).
func TestBriefBasisRefusesAMalformedKeptBrief(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err != nil {
		t.Fatal(err)
	}
	good := keptRecord(t, d.w, "bets/1/brief.md")
	d.restart()
	if b, err := d.r.briefBasis(context.Background(), "1"); err != nil || b.Appetite.USD != "5" {
		t.Fatalf("control: briefBasis = %+v, %v; want the appetite and no error", b, err)
	}
	for name, brief := range map[string]string{
		"only the appetite":        "## The appetite\nappetite: USD 5, 12 hours\n",
		"no facts heading":         strings.Replace(good, "## The facts\n", "", 1),
		"architecture change yes?": strings.Replace(good, "architecture change: no", "architecture change: perhaps", 1),
	} {
		d.push(map[string][]byte{"bets/1/brief.md": []byte(brief)}, name)
		if _, err := d.r.briefBasis(context.Background(), "1"); err == nil {
			t.Errorf("%s: briefBasis = no error; want the malformed brief refused", name)
		}
	}
}

// operatorRows leaves Intake with the Operator's rows of the block operator beside the idea owner's: Q-005, a setup question answered
// with answer, and Q-006, a setup question still open with the text openText (the shaped world holds Q-001..Q-004 and the answers Q-002 to Q-004).
func (d *drvWorld) operatorRows(answer, openText string) {
	d.t.Helper()
	qs := d.questions()[:4]
	open := qrow("Q-006", "S01-stack", "open")
	open[8] = openText
	qdata, err := table(records.QuestionsSchema, append(qs, qrow("Q-005", "S01-stack", "answered"), open))
	if err != nil {
		d.t.Fatal(err)
	}
	rows, err := records.ReadAnswers([]byte(keptRecord(d.t, d.w, "intake/answers.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	adata, err := table(records.AnswersSchema, append(rows[:3:3], []string{"Q-005", answer, "802", "303", "2026-10-08T12:40:00Z"}))
	if err != nil {
		d.t.Fatal(err)
	}
	d.push(map[string][]byte{"questions.tsv": qdata, "intake/answers.tsv": adata}, "the Operator's answer and question")
}

// annexSHA gives the last commit of the records that changed the linked files of the bet 1, from the bare repository.
func (d *drvWorld) annexSHA() string {
	d.t.Helper()
	return strings.TrimSpace(gitOut(d.t, d.w.bare, "log", "-1", "--format=%H", "layup-records", "--", "bets/1/facts.md", "bets/1/open.md"))
}

// annexCommits gives the number of records commits that changed a linked file.
func (d *drvWorld) annexCommits() int {
	d.t.Helper()
	return len(strings.Fields(gitOut(d.t, d.w.bare, "log", "--format=%H", "layup-records", "--", "bets/1/facts.md", "bets/1/open.md")))
}

// stopAtTheComment runs betBrief with the post of the control comment failing, so that a candidate stays and nothing is promoted.
func (d *drvWorld) stopAtTheComment() {
	d.t.Helper()
	d.r.cfg.Forge = &flakyForge{Forge: d.r.cfg.Forge, failPost: true}
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); err == nil || len(d.candidates()) == 0 {
		d.t.Fatalf("a stop at the comment: %v, candidates %v; want the error and a candidate", err, d.candidates())
	}
}

// I13, the demo (R11; task T-36xw, #471): the brief of an Intake with the Operator's rows links the facts and the questions left out at a
// pinned records commit that is written before the session; the files hold every text; a restart after the candidate reuses it at the
// same SHA, also after an unrelated records commit; the basis keeps every fact so that the answer can change an omitted one.
func TestTheLinkedFactsAreKeptAtAPinnedRecordsCommitBeforeTheSession(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.operatorRows("Use Go.", "The question Q-006")
	d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
	d.stopAtTheComment()
	sha := d.annexSHA()
	if !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(sha) || d.annexCommits() != 1 || d.briefSessions() != 1 {
		t.Fatalf("the annex commit %q, %d annex commits, %d sessions; want one commit and one session", sha, d.annexCommits(), d.briefSessions())
	}
	spans, _ := records.ReadSpans([]byte(keptRecord(t, d.w, "spec/spans.tsv")))
	questions := d.questions()
	answers, _ := records.ReadAnswers([]byte(keptRecord(t, d.w, "intake/answers.tsv")))
	wantFacts, wantOpen := factsPart(spans, questions, answers), openPart(questions)
	if got := keptRecord(t, d.w, "bets/1/facts.md"); got != wantFacts || !strings.Contains(got, "- Q-005 (context): Use Go.\n") {
		t.Errorf("bets/1/facts.md = %q; want %q", got, wantFacts)
	}
	if got := keptRecord(t, d.w, "bets/1/open.md"); got != wantOpen || !strings.Contains(got, "- Q-006: The question Q-006\n") {
		t.Errorf("bets/1/open.md = %q; want %q", got, wantOpen)
	}
	if got := gitOut(t, d.w.bare, "show", sha+":bets/1/facts.md"); got != wantFacts {
		t.Errorf("git show %s:bets/1/facts.md = %q; want the bytes of the records", sha, got)
	}
	if got := gitOut(t, d.w.bare, "show", "--format=%s", "--name-only", sha); got != "layup run: the facts and questions of the bet 1\n\nbets/1/facts.md\nbets/1/open.md\n" {
		t.Errorf("the annex commit shows %q; want its subject and the two files only", got)
	}
	// The annex commit comes before the session's commit of the candidate.
	if order := gitOut(t, d.w.bare, "log", "--format=%H", "layup-records", "--", d.candidates()[0]); strings.Index(gitOut(t, d.w.bare, "log", "--format=%H", "layup-records"), strings.TrimSpace(order)) > strings.Index(gitOut(t, d.w.bare, "log", "--format=%H", "layup-records"), sha) {
		t.Error("the candidate's commit is older than the annex commit; want the annex first")
	}
	factsURL := "https://github.com/acme/target/blob/" + sha + "/bets/1/facts.md"
	openURL := "https://github.com/acme/target/blob/" + sha + "/bets/1/open.md"
	candidate := keptRecord(t, d.w, d.candidates()[0])
	links, _ := os.ReadFile(filepath.Join(d.ctl, "brief-input-1", "links.md"))
	for _, want := range []string{"[the facts of the bet 1](" + factsURL + ")", "[the questions still open of the bet 1](" + openURL + ")",
		"- Q-005 (context): the answers of the Operator's questions;", "- Q-006: the Operator's questions still open,"} {
		if !strings.Contains(candidate, want) {
			t.Errorf("the candidate lacks %q:\n%s", want, candidate)
		}
	}
	if !strings.Contains(string(links), factsURL+"\n") || !strings.Contains(string(links), openURL+"\n") {
		t.Errorf("input/links.md = %q; want both annex URLs", links)
	}
	if strings.Contains(candidate, "Use Go.") || strings.Contains(candidate, "- Q-006: The question Q-006") {
		t.Errorf("the candidate holds the text of an answer or a question left out:\n%s", candidate)
	}

	// A restart after an unrelated records commit: the candidate is reused with no session and no annex commit, at the same SHA.
	d.push(map[string][]byte{"unrelated.txt": []byte("a heartbeat\n")}, "an unrelated records commit")
	brief, done, err := d.r.betBrief(context.Background(), "1", "first")
	if err != nil || done || string(brief) != candidate || d.briefSessions() != 1 || d.annexCommits() != 1 || d.annexSHA() != sha {
		t.Errorf("the restart = %d bytes, done %v, %v after %d sessions, %d annex commits, annex SHA %s; want the candidate, one session, one annex commit at %s",
			len(brief), done, err, d.briefSessions(), d.annexCommits(), d.annexSHA(), sha)
	}

	// The basis keeps every fact: a fact left out of the brief is a fact of the basis, and the answer can change it.
	basis, err := ReadBetBasis([]byte(keptRecord(t, d.w, "bets/1/basis.tsv")))
	if err != nil || !slices.ContainsFunc(basis, func(r []string) bool { return r[0] == "fact" && r[1] == "Q-005" && r[2] == "context" }) ||
		!slices.ContainsFunc(basis, func(r []string) bool { return r[0] == "open" && r[1] == "Q-006" }) {
		t.Errorf("bets/1/basis.tsv = %v, %v; want a fact row for Q-005 and an open row for Q-006", basis, err)
	}
	b, err := d.r.briefBasis(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	in := b.input([]string{"1"})
	answer, findings := spec.ParseAnswer([]byte("REQ-001: Must, phase 1\nfact Q-005: need\noption db: pg\nbet\n"), in.Required, in.Tags, in.Facts, in.Survivors)
	if len(findings) != 0 || answer.Facts["Q-005"] != (spec.FactAnswer{Class: "need"}) {
		t.Errorf("ParseAnswer on `fact Q-005: need` = %+v, findings %v; want the change of the omitted fact", answer.Facts, findings)
	}
	// The rendered facts record holds every fact with its text, the Operator's answer too: the rendering of the inventory of the basis.
	facts, texts, err := d.r.factsOf(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var version [][]string
	for _, f := range facts {
		version = append(version, []string{f.ID, f.Class, "in", "7004"})
	}
	rendered, fs := spec.RenderFacts(version, texts)
	want := "1. **need:** s1\n2. **measure:** s2\nQ-002. **context:** 5.0\nQ-003. **context:** 10.0\nQ-004. **context:** 40.0\nQ-005. **context:** Use Go.\n"
	if len(fs) != 0 || string(rendered) != want {
		t.Errorf("the rendered facts = %q, findings %v; want %q", rendered, fs, want)
	}
}

// Condition 4 of the plan review of #471: a text left out of the brief that changes while a candidate stands changes the pinned SHA, and the
// old candidate is refused with the rule fixed. The existing recovery test changes an inline span fact; neither case here can be seen there.
func TestAChangedTextLeftOutOfTheBriefRefusesTheOldCandidate(t *testing.T) {
	t.Parallel()
	for name, c := range map[string]struct {
		change func(*drvWorld)
		file   string
	}{
		"an omitted answer":        {func(d *drvWorld) { d.operatorRows("Use Rust.", "The question Q-006") }, "bets/1/facts.md"},
		"an omitted open question": {func(d *drvWorld) { d.operatorRows("Use Go.", "The question Q-006, reworded") }, "bets/1/open.md"},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			d := shapedWorld(t)
			d.operatorRows("Use Go.", "The question Q-006")
			d.ctlBrief(1, briefProblem, briefAppetite("5"), briefHoles, "")
			d.ctlBrief(2, briefProblem, briefAppetite("5"), briefHoles, "")
			d.stopAtTheComment()
			sha, old := d.annexSHA(), keptRecord(t, d.w, d.candidates()[0])
			c.change(d)
			brief, _, err := d.r.betBrief(context.Background(), "1", "first")
			if err != nil || string(brief) == old || d.briefSessions() != 2 || !strings.Contains(d.briefPrompt(2), "fixed:") {
				t.Fatalf("the recovery = %d bytes, %v after %d sessions; want a new brief, and the finding fixed: in the second prompt", len(brief), err, d.briefSessions())
			}
			if now := d.annexSHA(); now == sha || d.annexCommits() != 2 || !strings.Contains(string(brief), "/blob/"+now+"/bets/1/") || strings.Contains(string(brief), sha) {
				t.Errorf("the annex SHA %s (was %s) after %d annex commits; want a new commit, linked by the new brief and not the old SHA", now, sha, d.annexCommits())
			}
			if strings.Contains(keptRecord(t, d.w, c.file), "Use Go.\n") && c.file == "bets/1/facts.md" {
				t.Errorf("%s still holds the old text", c.file)
			}
		})
	}
}
