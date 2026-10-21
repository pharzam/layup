//go:build integration

package run

import (
	"context"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
)

// qaWorld ends a developer session whose result holds two questions and an
// answer, and gives what the end committed: the start row, the rows of the
// result, the session ID, and the comment that End posted.
func qaWorld(t *testing.T) (w *world, f *standInForge, id string, start []string, rows [][]string) {
	t.Helper()
	w, r, f, base := sessionWorld(t)
	table := resultHead + "question\t1\tFirst?\t—\tlabel needs_human\nquestion\t2\tWhich owner?\t—\tlabel needs_human\nanswer\t1\tUse the plan.\t—\tdocs/spec/loop.md\n"
	id, err := r.TaskSession(context.Background(), harness(t, base, "printf '"+octal(table)+"' > ../result/result.tsv"))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	for _, s := range startRows(t, w) {
		if s[sessionID] == id {
			start = s
		}
	}
	return w, f, id, start, resultsOf(t, w, id)[id]
}

func qaRowOf(rows [][]string, kind, n string) []string {
	for _, r := range rows {
		if r[0] == kind && r[1] == n {
			return r
		}
	}
	return nil
}

// The composers read the rows the real end committed; the expected bytes are
// written here with the session ID, never taken from the composer, and the
// first line starts like the comment that End posted.
func TestTheQuestionAndAnswerCommentsAreComposedFromTheRowsOfARealEnd(t *testing.T) {
	t.Parallel()
	w, f, id, start, rows := qaWorld(t)
	prefix := id + ": developer session of T-ab12, attempt 1: "
	if got := f.posted[taskIssue]; len(got) != 1 || got[0] != prefix+"done\n" {
		t.Fatalf("the end posted %q; want %q", got, prefix+"done\n")
	}
	head := gitOut(t, w.bare, "rev-parse", "layup-records")
	for _, c := range []struct{ n, want string }{
		{"1", prefix + "question " + id + ".1\n\nFirst?"},
		{"2", prefix + "question " + id + ".2\n\nWhich owner?"},
	} {
		got, err := QuestionComment(start, qaRowOf(rows, "question", c.n))
		if err != nil || string(got) != c.want {
			t.Errorf("question %s: %q, %v; want %q", c.n, got, err, c.want)
		}
	}
	want := prefix + "answer " + id + ".1 of S-0a1b2c3d.2\n\nUse the plan."
	got, err := AnswerComment(start, qaRowOf(rows, "answer", "1"), "S-0a1b2c3d.2")
	if err != nil || string(got) != want {
		t.Errorf("answer: %q, %v; want %q", got, err, want)
	}
	again, err := records.ReadResult([]byte(keptRecord(t, w, "tasks/T-ab12/results/"+id+".tsv"))) // a restart reads it again
	if err != nil {
		t.Fatal(err)
	}
	if second, err := AnswerComment(start, qaRowOf(again, "answer", "1"), "S-0a1b2c3d.2"); err != nil || string(second) != want {
		t.Errorf("after a second read: %q, %v; want %q", second, err, want)
	}
	if len(f.posted[taskIssue]) != 1 || gitOut(t, w.bare, "rev-parse", "layup-records") != head {
		t.Errorf("composing posted %q or moved the records; want neither", f.posted[taskIssue])
	}
}

func TestARowOfARealResultThatIsNoQuestionOrAnswerIsRefused(t *testing.T) {
	t.Parallel()
	_, _, _, start, rows := qaWorld(t)
	if got, err := QuestionComment(start, qaRowOf(rows, "question", "1")); err != nil || !strings.HasSuffix(string(got), "\n\nFirst?") {
		t.Fatalf("the control, a question row: %q, %v", got, err)
	}
	for name, got := range map[string]func() ([]byte, error){
		"status row":                 func() ([]byte, error) { return QuestionComment(start, qaRowOf(rows, "status", "1")) },
		"answer row to the question": func() ([]byte, error) { return QuestionComment(start, qaRowOf(rows, "answer", "1")) },
		"question row to the answer": func() ([]byte, error) { return AnswerComment(start, qaRowOf(rows, "question", "1"), "S-0a1b2c3d.2") },
		"status row to the answer":   func() ([]byte, error) { return AnswerComment(start, qaRowOf(rows, "status", "1"), "S-0a1b2c3d.2") },
	} {
		if b, err := got(); err == nil || b != nil {
			t.Errorf("%s: %q, %v; want an error and no bytes", name, b, err)
		}
	}
}
