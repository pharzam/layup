//go:build integration

package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/setup"
	workpkg "github.com/pharzam/layup/internal/work"
)

// The integration tests of the step markers and of the setup questions of the
// batch (docs/spec/scaffold.md, The step markers; The answers of the setup; the
// acceptance tests "The step markers" and "The setup questions in the batch and
// the comment"; row 53 of the plan, task T-4c5c, #276). The world's baseline
// holds the markers, and its harness writes bin/markers.tsv.

// markersFiles are the files of the baseline: two markers in one file, one in
// another, so the pinned tree has three places, in the order of S10's stop table.
func markersFiles() map[string]string {
	return map[string]string{
		"NOTES.md":       "# Notes\nThe owner is " + angled("owner name") + ".\nThe port is " + angled("port") + ".\n",
		"docs/deploy.md": "Deploy to " + angled("host") + ".\n",
	}
}

// markersBrief is the problem statement of the world: it states the owner and the
// port, and not the host.
const markersBrief = "# The problem\n\nThe owner is Carol Jones.\nThe port is 8080.\n"

// markerKeys are the IDs of the three markers, in the order of the places.
func markerKeys() [3]string {
	return [3]string{setup.MarkerID("NOTES.md", angled("owner name")), setup.MarkerID("NOTES.md", angled("port")), setup.MarkerID("docs/deploy.md", angled("host"))}
}

// markersToBatch runs a world with markers through Start and the restart that
// ends at the wait of the step answers, with the session's file in bin/.
func markersToBatch(t *testing.T, file string) (*world, []Step) {
	t.Helper()
	w := newWorldWith(t, markersFiles())
	if file != "" {
		os.WriteFile(filepath.Join(w.dir, "bin", "markers.tsv"), []byte(file), 0o644)
	}
	cfg := w.config("0.1.0-dev")
	cfg.PSB = []byte(markersBrief)
	allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	return w, Restart(context.Background(), w.config("0.1.0-dev"))
}

// proposals is the session's file: the port first, then the owner (the file's
// order is no matter), and none for the host.
func proposals() string {
	k := markerKeys()
	return "id\tvalue\tquote\n" + k[1] + "\t8080\tThe port is 8080.\n" + k[0] + "\tCarol Jones\tThe owner is Carol Jones.\n"
}

var intakeSteps = []string{"forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers"}

// The demo of row 53 (the scan, the session and the record): after the
// completeness review, one session of systems-architect proposes values for the
// markers of the pinned tree with quotes that code finds once; intake/markers.tsv
// holds a row for each marker, in the order of S10's stop table, with code's
// offsets; a second restart finds the record and starts no session.
func TestTheStepMarkers(t *testing.T) {
	t.Parallel()
	w, rows := markersToBatch(t, proposals())
	allDone(t, "the restart", rows, intakeSteps...)
	if d := rowOf(rows, "markers").Detail; d != "3 markers, 2 with a proposal, in intake/markers.tsv" {
		t.Errorf("the detail of markers %q", d)
	}
	got := show(t, w.bare, "refs/heads/layup-records", "intake/markers.tsv")
	table, err := records.ReadMarkers([]byte(got))
	if err != nil || len(table) != 3 {
		t.Fatalf("intake/markers.tsv: %d rows, %v:\n%s", len(table), err, got)
	}
	k := markerKeys()
	owner, port := strings.Index(markersBrief, "The owner is Carol Jones."), strings.Index(markersBrief, "The port is 8080.")
	want := [][]string{
		{k[0], "NOTES.md", "2", angled("owner name"), "Carol Jones", fmt.Sprint(owner), fmt.Sprint(owner + len("The owner is Carol Jones."))},
		{k[1], "NOTES.md", "3", angled("port"), "8080", fmt.Sprint(port), fmt.Sprint(port + len("The port is 8080."))},
		{k[2], "docs/deploy.md", "1", angled("host"), "", "", ""},
	}
	session := table[0][7]
	for i, r := range table {
		if strings.Join(r[:7], "|") != strings.Join(want[i], "|") || r[7] != session {
			t.Errorf("row %d of intake/markers.tsv %q; want %q and the session %s", i+1, r, want[i], session)
		}
	}
	sessions := show(t, w.bare, "refs/heads/layup-records", "sessions.tsv")
	if !strings.Contains(sessions, "\n"+session+"\t") || strings.Count(sessions, "\tsystems-architect\t") != 2 { // the review of meaning and this one
		t.Errorf("sessions.tsv has no start row of the session %s, or not two of systems-architect:\n%s", session, sessions)
	}
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
		t.Error("the lease is not released after the last row")
	}
	blob := gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/markers.tsv")
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if m := rowOf(again, "markers"); m.Result != "done" || m.Detail != "done before" {
		t.Errorf("the second restart's row %+v; want markers done before", m)
	}
	if gitOut(t, w.bare, "rev-parse", "refs/heads/layup-records:intake/markers.tsv") != blob ||
		strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tsystems-architect\t") != 2 {
		t.Error("the second restart changed intake/markers.tsv or started a session")
	}
}

// With no marker in the pinned tree the step writes a table with no row and
// starts no session; the batch then asks S01-stack alone in the Operator's block.
func TestTheStepMarkersWithNoMarker(t *testing.T) {
	t.Parallel()
	w := newWorld(t)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart", rows, intakeSteps...)
	if d := rowOf(rows, "markers").Detail; d != "no marker in the pinned tree; intake/markers.tsv has no row" {
		t.Errorf("the detail of markers %q", d)
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "intake/markers.tsv"); got != "id\tfile\tline\tmarker\tvalue\tstart\tend\tsession\n" {
		t.Errorf("intake/markers.tsv %q; want the header alone", got)
	}
	if n := strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tsystems-architect\t"); n != 1 { // the review of meaning alone
		t.Errorf("%d sessions of systems-architect; want the review of meaning alone", n)
	}
}

// Input states: a file of the session that breaks a rule is refused, and the step
// is fail, naming the session; nothing of the file is a record, and the lease
// is released (T-yngf; it stayed held). With no pair for systems-architect the step is fail (pair) and no
// session starts for it.
func TestTheStepMarkersFails(t *testing.T) {
	t.Parallel()
	k := markerKeys()
	head := "id\tvalue\tquote\n"
	for _, c := range []struct{ name, file, detail string }{
		{"a quote that occurs more than once", head + k[0] + "\tCarol\tThe\n", `: refused file: row 1: the quote "The" occurs 3 times`},
		{"a quote that occurs nowhere", head + k[0] + "\tCarol\tThe owner is Dave.\n", `: refused file: row 1: the quote "The owner is Dave." occurs nowhere`},
		{"a quote with a tab", head + k[0] + "\tCarol\tThe\towner\n", ": refused file: "},
		{"a key that the list lacks", head + "M-00000000\tx\tThe port is 8080.\n", ": refused file: row 1: the key M-00000000 is not in the list"},
		{"an empty value", head + k[0] + "\t—\tThe owner is Carol Jones.\n", ": refused file: row 1: the value is empty"},
		{"no route", "", "pair: no admitted pair"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := newWorldWith(t, markersFiles())
			if c.file != "" {
				os.WriteFile(filepath.Join(w.dir, "bin", "markers.tsv"), []byte(c.file), 0o644)
			}
			cfg := w.config("0.1.0-dev")
			cfg.PSB = []byte(markersBrief)
			allDone(t, "Start", w.start(cfg), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
			cfg = w.config("0.1.0-dev")
			if c.file == "" {
				// The review of meaning runs first, on the same pair: it is a record already, and the pair is gone.
				pushRecords(t, w.bare, map[string]string{"intake/meaning.tsv": "id\tkind\tstart\tend\tquote\tquestion\tsession\n"}, "records: no gap of meaning")
				cfg.Routing = [][]string{{"product-owner", "reasoning", "1", "claude", "m1"}, {"qa-engineer", "reasoning", "1", "claude", "m1"}, {"qa-engineer", "reasoning", "2", "devin", "m2"}}
			}
			rows := Restart(context.Background(), cfg)
			last := rows[len(rows)-1]
			if last.Name != "markers" || last.Result != "fail" || !strings.Contains(last.Detail, c.detail) {
				t.Fatalf("the last row %+v; want markers fail with %q", last, c.detail)
			}
			if c.file != "" && !strings.HasPrefix(last.Detail, "the session S-") {
				t.Errorf("the detail %q does not name the session", last.Detail)
			}
			if strings.Contains(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "refs/heads/layup-records"), "intake/markers.tsv") {
				t.Error("a refused file is a record")
			}
			if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "lease.tsv"), "\treleased\n") {
				t.Error("the lease is not released after a fail") // T-yngf: was held, the known limit
			}
			// The review of meaning and the refused session, or, with no pair, none: the review of meaning is a record already.
			want := 2
			if c.file == "" {
				want = 0
			}
			if sessions := strings.Count(show(t, w.bare, "refs/heads/layup-records", "sessions.tsv"), "\tsystems-architect\t"); sessions != want {
				t.Errorf("%d sessions of systems-architect; want %d", sessions, want)
			}
		})
	}
}

// The demo of the setup questions (the merge and the comment): with the
// markers of the world, the step batch merges, after the questions of the idea
// owner's block, S01-stack and one question for each marker, in the order of
// intake/markers.tsv, each a row of questions.tsv with the source setup, in the
// Operator's block; the comment writes a marker's question with its proposal
// and its quote, and the rule of the answers has the forms of a marker's
// answer. S01-stack is asked as well as the rule gap of G1.
func TestTheStepBatchAsksTheSetupQuestions(t *testing.T) {
	t.Parallel()
	w, rows := markersToBatch(t, proposals())
	allDone(t, "the restart", rows, intakeSteps...)
	table, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	var setups [][]string
	for _, q := range table {
		if q[1] == "setup" {
			setups = append(setups, q)
		}
	}
	k := markerKeys()
	owner, port := strings.Index(markersBrief, "The owner is Carol Jones."), strings.Index(markersBrief, "The port is 8080.")
	n := len(table)
	want := [][]string{
		{fmt.Sprintf("Q-%03d", n-3), "S01-stack", "", "", workpkg.S01Questions[0].Text, "layup run"},
		{fmt.Sprintf("Q-%03d", n-2), k[0], fmt.Sprint(owner), fmt.Sprint(owner + len("The owner is Carol Jones.")), "the value of " + angled("owner name") + " at NOTES.md:2", ""},
		{fmt.Sprintf("Q-%03d", n-1), k[1], fmt.Sprint(port), fmt.Sprint(port + len("The port is 8080.")), "the value of " + angled("port") + " at NOTES.md:3", ""},
		{fmt.Sprintf("Q-%03d", n), k[2], "", "", "the value of " + angled("host") + " at docs/deploy.md:1", ""},
	}
	if len(setups) != 4 {
		t.Fatalf("%d setup questions %q; want S01-stack and the three markers", len(setups), setups)
	}
	markersTable, _ := records.ReadMarkers([]byte(show(t, w.bare, "refs/heads/layup-records", "intake/markers.tsv")))
	for i, q := range setups {
		asker := want[i][5]
		if asker == "" {
			asker = markersTable[i-1][7] // the session of its row of intake/markers.tsv
		}
		if q[0] != want[i][0] || q[2] != want[i][1] || q[3] != "operator" || q[4] != "before-delivery" || q[5] != "setup" ||
			q[6] != want[i][2] || q[7] != want[i][3] || q[8] != want[i][4] || q[10] != asker || q[11] != "" || q[12] != "asked" {
			t.Errorf("setup question %d %q; want %q, asked by %q", i+1, q, want[i], asker)
		}
	}
	// S01-stack is asked with the rule gap of G1 (the world's problem statement names no stack).
	if !strings.Contains(show(t, w.bare, "refs/heads/layup-records", "questions.tsv"), "\trule\tQ-001\tidea-owner\tbefore-delivery\tG1\t") {
		t.Error("questions.tsv has no rule gap G1: the test does not show that S01-stack is asked whatever G1 found")
	}
	c := intakeComments(w)
	if len(c) != 1 {
		t.Fatalf("%d comments of the App; want the batch", len(c))
	}
	for _, line := range []string{
		fmt.Sprintf("\n## For the Operator\n\nQ-%03d: %s\n\n", n-3, workpkg.S01Questions[0].Text),
		fmt.Sprintf("\nQ-%03d: the value of %s at NOTES.md:2\nproposed: Carol Jones\n`The owner is Carol Jones.`\n\n", n-2, angled("owner name")),
		fmt.Sprintf("\nQ-%03d: the value of %s at NOTES.md:3\nproposed: 8080\n`The port is 8080.`\n\n", n-1, angled("port")),
		fmt.Sprintf("\nQ-%03d: the value of %s at docs/deploy.md:1\n\n", n, angled("host")),
		"`Q-NNN: accept`", "`Q-NNN: <value>`", "`Q-NNN: gap`",
	} {
		if !strings.Contains(c[0], line) {
			t.Errorf("the comment of the batch has no %q:\n%s", line, c[0])
		}
	}
}

// The follow-up of a setup question (row 53 and docs/spec/intake.md, The review
// of the answers and the follow-up): the idea owner answers every question, the
// Operator only S01-stack, so the follow-up asks the three markers again, in the
// Operator's block, each with the quote of its first question and, for a marker
// with a proposal, the line proposed:; the idea owner's block says it asks none.
func TestTheFollowUpAsksTheMarkersAgain(t *testing.T) {
	t.Parallel()
	w, _ := markersToBatch(t, proposals())
	table, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	asked, _ := time.Parse(time.RFC3339, table[0][9])
	var idea, operator strings.Builder
	for _, q := range table {
		switch {
		case q[1] == "setup" && q[2] == "S01-stack":
			fmt.Fprintf(&operator, "%s: go\n", q[0])
		case q[1] == "intent":
			fmt.Fprintf(&idea, "%s: %s\n", q[0], map[string]string{"band.B": "5.0", "band.U": "10.0", "appetite.hours": "40.0", "success-criteria": "it does its work"}[q[2]])
		case q[1] != "setup":
			fmt.Fprintf(&idea, "%s: the answer to %s\n", q[0], q[0])
		}
	}
	addComment(w, 1, 802, 303, "carol", asked.Add(time.Minute).UTC().Format(time.RFC3339), idea.String())
	addComment(w, 1, 803, 101, "pharzam", asked.Add(2*time.Minute).UTC().Format(time.RFC3339), operator.String())
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	allDone(t, "the restart with the answers", rows, append(append([]string{}, intakeSteps...), "answers-review", "follow-up", "follow-up-answers")...)
	if d := rowOf(rows, "follow-up").Detail; d != "3 questions asked again, in questions.tsv and the comment of the follow-up" {
		t.Errorf("the detail of follow-up %q", d)
	}
	if d := rowOf(rows, "follow-up-answers").Detail; d != "waiting: no answer yet to the follow-up in the block operator" {
		t.Errorf("the detail of follow-up-answers %q; want the Operator's block to wait", d)
	}
	n := len(table)
	c := intakeComments(w)
	if len(c) != 2 || !strings.HasPrefix(c[1], followUpHeader+"\n") {
		t.Fatalf("the App's comments on the Intake issue %q; want the batch and one follow-up", c)
	}
	for _, want := range []string{
		"\n## For the idea owner\n\nNo question for the idea owner.\n\n## For the Operator\n\n",
		fmt.Sprintf("\nQ-%03d: the value of %s at NOTES.md:2\nproposed: Carol Jones\n`The owner is Carol Jones.`\n\n", n+1, angled("owner name")),
		fmt.Sprintf("\nQ-%03d: the value of %s at NOTES.md:3\nproposed: 8080\n`The port is 8080.`\n\n", n+2, angled("port")),
		fmt.Sprintf("\nQ-%03d: the value of %s at docs/deploy.md:1\n\n", n+3, angled("host")),
		"`Q-NNN: accept`",
	} {
		if !strings.Contains(c[1], want) {
			t.Errorf("the follow-up has no %q:\n%s", want, c[1])
		}
	}
	// A run that stopped between the rows and the post leaves the rows: the next restart posts
	// the comment again from the records, byte for byte, with the proposals.
	w.fake.mu.Lock()
	var kept []map[string]any
	for _, cm := range w.fake.comments[1] {
		if body, _ := cm["body"].(string); !strings.HasPrefix(body, followUpHeader+"\n") {
			kept = append(kept, cm)
		}
	}
	w.fake.comments[1] = kept
	w.fake.mu.Unlock()
	again := Restart(context.Background(), w.config("0.1.0-dev"))
	if f := rowOf(again, "follow-up"); f.Result != "done" || f.Detail != "the comment of the follow-up posted from questions.tsv" {
		t.Errorf("the row %+v; want follow-up done, posted from the record", f)
	}
	if c2 := intakeComments(w); len(c2) != 2 || c2[1] != c[1] {
		t.Errorf("the comment posted again %q; want the first's bytes", c2)
	}
}

// The scan reads the tree of the one root commit of the default branch, which is
// pin.tree, in a scratch work tree that it removes: not the head of the branch,
// which after Scaffold 6 holds the setup commits; and a root tree that is not
// pin.tree, or a default branch of two root commits, is an error (the step is
// fail, scaffold.md, The step markers).
func TestThePinnedTreeOfTheMarkers(t *testing.T) {
	t.Setenv("TMPDIR", t.TempDir())
	dir := t.TempDir()
	gitOut(t, dir, "init", "-q", "-b", "main")
	gitOut(t, dir, "config", "user.name", "t")
	gitOut(t, dir, "config", "user.email", "t@example.invalid")
	os.WriteFile(filepath.Join(dir, "NOTES.md"), []byte("The port is "+angled("port")+".\n"), 0o644)
	gitOut(t, dir, "add", "-A")
	gitOut(t, dir, "commit", "-q", "-m", "root")
	tree := strings.TrimSpace(gitOut(t, dir, "rev-parse", "HEAD^{tree}"))
	// The head holds another marker, as the setup commits may.
	os.WriteFile(filepath.Join(dir, "LATER.md"), []byte(angled("later")+"\n"), 0o644)
	gitOut(t, dir, "add", "-A")
	gitOut(t, dir, "commit", "-q", "-m", "head")
	gitOut(t, dir, "update-ref", "refs/remotes/origin/main", "HEAD")
	r := &state{target: dir, defaultBranch: "main", start: [][]string{{"pin.tree", tree, "run"}}}
	places, err := r.pinnedPlaces()
	if err != nil || len(places) != 1 || places[0].File != "NOTES.md" || places[0].Line != 1 || places[0].Marker != angled("port") {
		t.Fatalf("the places %+v, %v; want the one marker of the root commit's tree", places, err)
	}
	if list := strings.TrimSpace(gitOut(t, dir, "worktree", "list")); strings.Contains(list, "\n") {
		t.Errorf("the scratch work tree is not removed:\n%s", list)
	}
	if left, _ := filepath.Glob(filepath.Join(os.TempDir(), "layup-markers-*")); len(left) != 0 {
		t.Errorf("the scratch directory is not removed: %q", left)
	}
	r.start = [][]string{{"pin.tree", strings.Repeat("0", 40), "run"}}
	if _, err := r.pinnedPlaces(); err == nil || !strings.Contains(err.Error(), "is not pin.tree") {
		t.Errorf("a root tree that is not pin.tree: %v", err)
	}
	// Two root commits: an unrelated history merged into the branch.
	r.start = [][]string{{"pin.tree", tree, "run"}}
	gitOut(t, dir, "checkout", "-q", "--orphan", "other")
	gitOut(t, dir, "rm", "-rfq", ".")
	os.WriteFile(filepath.Join(dir, "OTHER.md"), []byte("other\n"), 0o644)
	gitOut(t, dir, "add", "-A")
	gitOut(t, dir, "commit", "-q", "-m", "other root")
	gitOut(t, dir, "checkout", "-q", "main")
	gitOut(t, dir, "merge", "-q", "--allow-unrelated-histories", "-m", "merge", "other")
	gitOut(t, dir, "update-ref", "refs/remotes/origin/main", "HEAD")
	if _, err := r.pinnedPlaces(); err == nil || !strings.Contains(err.Error(), "has 2 root commits") {
		t.Errorf("two root commits: %v", err)
	}
}
