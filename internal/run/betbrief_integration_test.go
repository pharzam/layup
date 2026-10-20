//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
func shapedWorld(t *testing.T) *drvWorld {
	t.Helper()
	d := newDrvWorldOpts(t, loopPerms, drvOpts{noOption: true})
	// The clock of the world runs 600 times faster than the real one, so the 8 hours of the cap of milestone 0 are 48 seconds, and a host under
	// load would end a test at that cap by the clock. The hours are widened; the money cap stays 50.
	start := keptRecord(t, d.w, "start/start.tsv")
	if !strings.Contains(start, "intake.cap\t50.0,8.0\t") {
		t.Fatalf("start.tsv holds no intake.cap of 50.0,8.0:\n%s", start)
	}
	d.push(map[string][]byte{"start/start.tsv": []byte(strings.Replace(start, "intake.cap\t50.0,8.0\t", "intake.cap\t50.0,4000.0\t", 1))}, "the hours of the cap of milestone 0 widened")
	d.settled()
	d.ctlOptions(1, survivorDB)
	d.ctlWrite("commit", "")
	if got := d.r.shapeStep(context.Background()); got.Result != "done" || !strings.Contains(got.Detail, " verified") {
		t.Fatalf("the step shape = %+v; want both tasks verified", got)
	}
	d.setParameter("brief.lines", "60")
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
	d.setParameter("brief.lines", "")
	head := recordsHead(t, d.w)
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); !errors.Is(err, ErrBriefLinesUnset) {
		t.Errorf("err = %v; want ErrBriefLinesUnset", err)
	}
	if d.briefSessions() != 0 || recordsHead(t, d.w) != head {
		t.Errorf("%d sessions, head moved: %v; want none (a green negative beside the error)", d.briefSessions(), recordsHead(t, d.w) != head)
	}
}

// The parts of code over brief.lines wait before any session, and a fresh build refuses a bets.tsv that exists.
func TestAFreshBuildRefusesWhatItCannotBuildBefore(t *testing.T) {
	t.Parallel()
	d := shapedWorld(t)
	d.setParameter("brief.lines", "10")
	if _, _, err := d.r.betBrief(context.Background(), "1", "first"); !errors.Is(err, ErrBriefOverLimit) || d.briefSessions() != 0 {
		t.Errorf("brief.lines 10: err = %v after %d sessions; want ErrBriefOverLimit and none", err, d.briefSessions())
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
