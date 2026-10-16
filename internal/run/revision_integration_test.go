//go:build integration

package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
)

// The integration tests of the step revision and of the approvers' naming
// (docs/spec/loop.md, The approvers and the success criteria and The revision
// of the draft; the acceptance row of M2e; row 66i of the plan, task T-8rm6,
// #365). The world's real Restart and PhaseSession run them; its scripted
// harness answers the revision with bin/revision.tsv, or, with none, keeps
// the draft byte for byte, and the fake forge resolves dave alone of the
// logins named (pharzam and carol hold their rows of the Start).

// jumpClock is a world's clock with an added jump the test sets once: the
// bound of the retry of the revision is the cap of milestone 0, read before
// each attempt (docs/spec/intake.md, The budget of Intake).
type jumpClock struct {
	Clock
	mu   sync.Mutex
	plus time.Duration
}

func (c *jumpClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Clock.Now().Add(c.plus)
}

func (c *jumpClock) jump(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.plus += d
}

// intakeAnswers posts the answers of a whole batch on the Intake issue of a
// world that waits at the step answers: the idea owner answers each question
// of her block — the intent fields with their decimals and the success
// criteria with criteria — and names the approvers in lines `<point>:
// <login>` of the same comment (points may be ""); the Operator answers each
// question of his block: S01-stack with "go", a marker with a proposal
// "accept", one without "gap". It gives the rows of questions.tsv.
func intakeAnswers(t *testing.T, w *world, criteria, points string) [][]string {
	t.Helper()
	questions, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	markers, err := records.ReadMarkers([]byte(show(t, w.bare, "refs/heads/layup-records", "intake/markers.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	proposed := map[string]bool{}
	for _, m := range markers {
		proposed[m[0]] = m[4] != ""
	}
	asked, err := time.Parse(time.RFC3339, questions[0][9])
	if err != nil {
		t.Fatal(err)
	}
	var owner, operator strings.Builder
	for _, q := range questions {
		answer := "the answer to " + q[0]
		switch {
		case q[3] != "operator":
			switch q[2] {
			case "band.B":
				answer = "5.0"
			case "band.U":
				answer = "10.0"
			case "appetite.hours":
				answer = "40.0"
			case "success-criteria":
				answer = criteria
			}
		case q[2] == "S01-stack":
			answer = "go"
		case proposed[q[2]]:
			answer = "accept"
		default:
			answer = "gap"
		}
		if q[3] == "operator" {
			fmt.Fprintf(&operator, "%s: %s\n", q[0], answer)
		} else {
			fmt.Fprintf(&owner, "%s: %s\n", q[0], answer)
		}
	}
	owner.WriteString(points)
	at := func(d time.Duration) string { return asked.Add(d).UTC().Format(time.RFC3339) }
	addComment(w, 1, 802, 303, "carol", at(2*time.Minute), owner.String())
	addComment(w, 1, 803, 101, "pharzam", at(3*time.Minute), operator.String())
	return questions
}

// throughFollowUp gives the steps a restart with the answers done runs in
// this package's worlds, revision last before the step setup fails on the
// file the baseline lacks.
var throughRevision = []string{"forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft",
	"completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision"}

// sessionsOf counts the rows of sessions.tsv of the role product-owner, the
// revision's alone of the sessions these worlds run; telemetryOf the rows of
// telemetry.tsv of the same role.
func sessionsOf(t *testing.T, w *world) int {
	t.Helper()
	return strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tproduct-owner\t")
}

func telemetryOf(t *testing.T, w *world) int {
	t.Helper()
	return strings.Count(show(t, w.bare, "refs/heads/layup-records", "telemetry.tsv"), "\tproduct-owner\t")
}

// draftBefore gives the rows of spec/draft.tsv of the records and the ID of
// the criteria question.
func draftOf(t *testing.T, w *world) [][]string {
	t.Helper()
	rows, err := records.ReadDraft([]byte(show(t, w.bare, "refs/heads/layup-records", "spec/draft.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// The demo of the acceptance row of M2e: the idea owner's answer to the
// intent field success-criteria lands in spec/draft.tsv as one
// success-criterion row per criterion, appended by the code after the
// session's kept file, in the order of the answer and covering the question;
// the approver the idea owner names for a point gets its row of approvers.tsv
// with the naming comment as its source, while a login the forge does not
// know names nothing; and a row of the session may cover an answer's Q-NNN
// together with the fact of its quote. A second restart is "done before", and
// a restart after a stop between the commit and the post reposts the
// session's comment from the record, byte for byte.
func TestTheRevisionOfTheDraft(t *testing.T) {
	w, rows := markersToBatch(t, proposals())
	allDone(t, "the restart", rows, intakeSteps...)
	questions := intakeAnswers(t, w, "fast; cheap", "bet: dave\naudit: pharzam\nretrospective: mallory\n")
	var criteriaQ, portQ string
	var port []string
	spans, err := records.ReadSpans([]byte(show(t, w.bare, "refs/heads/layup-records", "spec/spans.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range questions {
		switch q[2] {
		case "success-criteria":
			criteriaQ = q[0]
		case markerKeys()[1]: // the marker of the port
			portQ, port = q[0], q
		}
	}
	portFact := factOfQuote(port[6], port[7], spans)
	if criteriaQ == "" || portQ == "" || portFact == "" {
		t.Fatalf("the questions: criteria %q, port %q, the fact of its quote %q", criteriaQ, portQ, portFact)
	}
	// The session's file keeps the draft's three rows and adds a row that
	// covers the answer to the question of the port's marker, with the fact of
	// its quote (docs/spec/loop.md, The revision of the draft).
	os.WriteFile(filepath.Join(w.dir, "bin", "revision.tsv"), []byte(
		"kind\tstatement\tcovers\tcriterion\n"+
			"requirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.\n"+
			"requirement\tThe statement of fact 2.\t2\tThe criterion of fact 2.\n"+
			"requirement\tThe statement of fact 3.\t3\tThe criterion of fact 3.\n"+
			"requirement\tThe port as the answer names it.\t"+portFact+" "+portQ+"\tThe port is set.\n"), 0o644)
	rows = beforeSetup(t, Restart(context.Background(), w.config("0.1.0-dev")))
	allDone(t, "the restart with the answers", rows, throughRevision...)
	d := rowOf(rows, "revision").Detail
	if d != "the draft revised, 6 rows of which 2 the success criteria, in spec/draft.tsv" {
		t.Errorf("the detail of revision %q", d)
	}
	draft := draftOf(t, w)
	session := draft[0][5]
	if len(draft) != 6 || session == "" {
		t.Fatalf("the draft of the revision:\n%s", show(t, w.bare, "refs/heads/layup-records", "spec/draft.tsv"))
	}
	for i, want := range []string{
		"requirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.",
		"requirement\tThe statement of fact 2.\t2\tThe criterion of fact 2.",
		"requirement\tThe statement of fact 3.\t3\tThe criterion of fact 3.",
		"requirement\tThe port as the answer names it.\t" + portFact + " " + portQ + "\tThe port is set.",
		"success-criterion\tThe delivery's success criteria.\t" + criteriaQ + "\tfast",
		"success-criterion\tThe delivery's success criteria.\t" + criteriaQ + "\tcheap",
	} {
		if got := strings.Join(draft[i][1:5], "\t"); got != want || draft[i][0] != fmt.Sprintf("R-%03d", i+1) || draft[i][5] != session {
			t.Errorf("row %d of the kept draft %q; want %q, R-%03d and the session %s", i+1, draft[i], want, i+1, session)
		}
	}
	// dave gets the one new row, its since the naming comment's time and its
	// source the comment; pharzam names audit on his own row, which stays
	// single, and the unknown mallory names nothing.
	approvers := show(t, w.bare, "refs/heads/layup-records", "approvers.tsv")
	if !strings.Contains(approvers, "\n707\tapprover\tdave\t") || !strings.HasSuffix(approvers, "\t802\n") ||
		strings.Count(approvers, "\n101\t") != 1 || strings.Count(approvers, "\n") != 4 || strings.Contains(approvers, "mallory") {
		t.Errorf("approvers.tsv:\n%s", approvers)
	}
	sessions := show(t, w.bare, "refs/heads/layup-records", "sessions.tsv")
	if !strings.Contains(sessions, "\n"+session+"\t") || !strings.Contains(sessions, session+"\tT-") {
		t.Errorf("sessions.tsv has no row of the session of the revision %s:\n%s", session, sessions)
	}
	line := session + ": product-owner session of revision: "
	w.fake.mu.Lock()
	var posted []string
	for _, c := range w.fake.comments[2] {
		posted = append(posted, c["body"].(string))
	}
	w.fake.mu.Unlock()
	var comment string
	for _, p := range posted {
		if strings.HasPrefix(p, line) {
			comment = p
		}
	}
	if comment != line+"done\n" {
		t.Errorf("the comments of the control issue %q; want %q", posted, line+"done\\n")
	}
	// A second restart is done before and changes no record.
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:spec/draft.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if r := rowOf(again, "revision"); r.Result != "done" || r.Detail != "done before" {
		t.Errorf("the second restart's row %+v; want revision done before", r)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:spec/draft.tsv") != blob {
		t.Error("the second restart changed spec/draft.tsv")
	}
	// A run that stopped between the commit and the post: the comment is
	// posted from the record, byte for byte, with no second commit.
	w.fake.mu.Lock()
	var kept []map[string]any
	for _, c := range w.fake.comments[2] {
		if body, _ := c["body"].(string); !strings.HasPrefix(body, line) {
			kept = append(kept, c)
		}
	}
	w.fake.comments[2] = kept
	w.fake.mu.Unlock()
	third := Restart(context.Background(), w.config("0.1.0-dev"))
	if r := rowOf(third, "revision"); r.Result != "done" || r.Detail != "the comment of the session posted from the record" {
		t.Errorf("the third restart's row %+v; want revision done, posted from the record", r)
	}
	w.fake.mu.Lock()
	back := ""
	for _, c := range w.fake.comments[2] {
		if strings.HasPrefix(c["body"].(string), line) {
			back = c["body"].(string)
		}
	}
	w.fake.mu.Unlock()
	if back != comment {
		t.Errorf("the comment posted again %q; want the first's bytes %q", back, comment)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:spec/draft.tsv") != blob {
		t.Error("the repost restart changed spec/draft.tsv")
	}
}

// A file that the check refuses ends the session done but refused, with its
// telemetry row committed; the step runs the session again with the check's
// finding in the prompt, each attempt printed, until a kept file (the bound
// is the cap of milestone 0, asked before each attempt: docs/spec/intake.md,
// The budget of Intake).
func TestTheRevisionRefusedThenKept(t *testing.T) {
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), intakeSteps...)
	intakeAnswers(t, w, "fast; cheap", "")
	// Attempt 1 drops a fact of the class need; attempt 2, written by the hook
	// when its attempt line prints, keeps the draft.
	os.WriteFile(filepath.Join(w.dir, "bin", "revision.tsv"), []byte(
		"kind\tstatement\tcovers\tcriterion\nrequirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.\n"), 0o644)
	bad := "the need fact 2 is covered by no row"
	telemetry := telemetryOf(t, w)
	w.hook = func(line string) {
		if line == "revision: attempt 2" {
			os.WriteFile(filepath.Join(w.dir, "bin", "revision.tsv"), []byte(
				"kind\tstatement\tcovers\tcriterion\n"+
					"requirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.\n"+
					"requirement\tThe statement of fact 2.\t2\tThe criterion of fact 2.\n"), 0o644)
		}
	}
	rows := beforeSetup(t, Restart(context.Background(), w.config("0.1.0-dev")))
	allDone(t, "the restart with the answers", rows, throughRevision...)
	if d := rowOf(rows, "revision").Detail; d != "the draft revised, 4 rows of which 2 the success criteria, in spec/draft.tsv" {
		t.Errorf("the detail of revision %q", d)
	}
	w.mu.Lock()
	lines := strings.Join(w.lines, "\n")
	w.mu.Unlock()
	if !strings.Contains(lines, "revision: attempt 1") || !strings.Contains(lines, "revision: attempt 2") || strings.Contains(lines, "revision: attempt 3") {
		t.Errorf("the attempts printed: %q", lines)
	}
	prompt, err := os.ReadFile(filepath.Join(w.dir, "bin", "revision-prompt.txt"))
	if err != nil || !strings.Contains(string(prompt), "The check refused the last revision: "+bad+"\n") {
		t.Errorf("the prompt of the second attempt has no finding %q: %v\n%s", bad, err, prompt)
	}
	// The refused attempt committed its telemetry row and its comment, and
	// the kept one the same once more.
	if n := telemetryOf(t, w); n != telemetry+2 {
		t.Errorf("telemetry.tsv gained %d rows of product-owner; want the two attempts", n-telemetry)
	}
	w.fake.mu.Lock()
	var ends []string
	for _, c := range w.fake.comments[2] {
		if strings.Contains(c["body"].(string), " session of revision: ") {
			ends = append(ends, strings.TrimSpace(c["body"].(string)))
		}
	}
	w.fake.mu.Unlock()
	if len(ends) != 2 || !strings.HasSuffix(ends[0], "done, refused file") || !strings.HasSuffix(ends[1], "done") {
		t.Errorf("the comments of the revision's sessions %q; want refused file, then done", ends)
	}
}

// The cap of milestone 0, asked before each attempt, ends the retry: the
// refused attempt's telemetry row is the record, no second session starts and
// the step is fail with the cap's refusal (docs/spec/intake.md, The budget of
// Intake; row 66i).
func TestTheRevisionStopsAtTheCap(t *testing.T) {
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	allDone(t, "the restart", Restart(context.Background(), w.config("0.1.0-dev")), intakeSteps...)
	intakeAnswers(t, w, "fast; cheap", "")
	os.WriteFile(filepath.Join(w.dir, "bin", "revision.tsv"), []byte(
		"kind\tstatement\tcovers\tcriterion\nrequirement\tThe statement of fact 1.\t1\tThe criterion of fact 1.\n"), 0o644)
	sessions := sessionsOf(t, w)
	cfg := w.config("0.1.0-dev")
	clock := &jumpClock{Clock: cfg.Clock}
	cfg.Clock = clock
	w.hook = func(line string) {
		if line == "revision: attempt 2" {
			clock.jump(9 * time.Hour) // past intake.cap's 8.0 hours
		}
	}
	rows := Restart(context.Background(), cfg)
	last := rows[len(rows)-1]
	if last.Name != "revision" || last.Result != "fail" || last.Detail != "the cap of milestone 0: the time since Start passes 8.0 hours" {
		t.Fatalf("the last row %+v; want revision fail, the cap\n%s", last, steps(rows))
	}
	if n := sessionsOf(t, w); n != sessions+1 {
		t.Errorf("%d rows of sessions.tsv of product-owner after the run; want the one refused session", n-sessions)
	}
	w.mu.Lock()
	lines := strings.Join(w.lines, "\n")
	w.mu.Unlock()
	if strings.Contains(lines, "revision: attempt 3") || !strings.Contains(lines, "revision: attempt 2") {
		t.Errorf("the attempts printed: %q; want 1 and 2, no third", lines)
	}
	draft := show(t, w.bare, "refs/heads/layup-records", "spec/draft.tsv")
	if strings.Contains(draft, "success-criterion\tThe delivery's success criteria.") {
		t.Errorf("a criteria row of a refused revision is a record:\n%s", draft)
	}
}
