//go:build integration

package run

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
)

// failingTable runs the real gate over a repository whose manifest has a kind
// that passes and a kind that exits 3, and gives its table.
func failingTable(t *testing.T) gate.Table {
	t.Helper()
	dir := t.TempDir()
	if err := git.Init(dir); err != nil {
		t.Fatal(err)
	}
	who := git.Identity{Name: "t", Email: "t@t.invalid", Time: time.Unix(1, 0)}
	write := func(files map[string]string) {
		for name, text := range files {
			if err := os.WriteFile(filepath.Join(dir, name), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if err := git.Add(dir); err != nil {
			t.Fatal(err)
		}
		if err := git.Commit(dir, "test", who); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(dir, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	write(map[string]string{"docs/gates.tsv": "kind\tstate\ttool\tcommand\tscope\tconfig\nfine\tactive\tsh\ttrue\t./*.go\t—\nexits\tactive\tsh\texit 3\t./*.go\t—\n", "x.go": "package x\n"})
	base, err := git.RevParse(dir, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	write(map[string]string{"x.go": "package x // changed\n"})
	table, err := gate.Run(dir, base, "HEAD", func(int, int, string) func() { return func() {} }, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	return table
}

// The demo of loop.md, the attempts: a failure of step 5 starts a new attempt
// at the base, and the prompt of its session holds the findings and not the
// diff of the attempt before.
func TestAGateFailureStartsTheNextAttemptFromTheBaseWithItsFindings(t *testing.T) {
	l := newLoop(t, nil)
	ctx := context.Background()
	l.dev(1, "claude", map[string]string{"x.go": "package x // marker-of-attempt-one\n"}, okResult)
	table := failingTable(t)
	s := l.prod()
	got, err := s.NextAttempt(ctx, NextAttemptInput{Task: "T-ab12", Base: l.a.base, Table: table})
	if err != nil || got.N != 2 || len(got.Findings) != 1 || got.Findings[0].Rule != "gate:exits" || got.Findings[0].Text != "exit 3" {
		t.Fatalf("NextAttempt = %+v, %v; want attempt 2 with the finding gate:exits, exit 3 (the kind that passed is no finding)", got, err)
	}
	ev := l.events()
	if last := ev[len(ev)-1]; last[1] != "attempt" || last[2] != "2" || last[4] != l.a.base {
		t.Errorf("the last event %v; want attempt 2 at the base %s", last, l.a.base)
	}
	pc, err := s.PromptFile(ctx, PromptIn{Task: "T-ab12", Step: "implement", Base: l.a.base, Findings: got.Findings})
	if err != nil {
		t.Fatal(err)
	}
	spec := harness(t, l.a, `cp ../prompt.md seen-prompt.txt && git add -A && git commit -qm seen && printf '`+octal(okResult)+`' > ../result/result.tsv`)
	spec.Role, spec.Attempt, spec.Context = roleOfStep("implement"), 2, pc
	id, err := l.r.TaskSession(ctx, spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	head := l.bound(id)
	prompt := gitOut(t, l.w.bare, "show", head+":seen-prompt.txt")
	if !strings.Contains(prompt, "gate:exits") || !strings.Contains(prompt, "exit 3") {
		t.Errorf("the prompt of attempt 2 holds no finding:\n%s", prompt)
	}
	if strings.Contains(prompt, "marker-of-attempt-one") || strings.Contains(prompt, "diff --git") {
		t.Errorf("the prompt of attempt 2 holds the diff of attempt 1:\n%s", prompt)
	}
	// Attempt 2 starts from the base, not from the head of attempt 1: its one commit is the prompt's copy.
	if n := len(strings.Fields(gitOut(t, l.w.bare, "rev-list", l.a.base+".."+head))); n != 1 {
		t.Errorf("%d commits from the base to the head of attempt 2; want 1", n)
	}
}

func TestAnAttemptThatAQuestionEndedStillCountsInTheNumbering(t *testing.T) {
	l := newLoop(t, nil)
	asking := resultHead[:strings.Index(resultHead, "status")] + "status\t1\tneeds_context\t—\tasks\nquestion\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md\n"
	l.dev(1, "claude", map[string]string{"q.txt": "q\n"}, asking)
	got, err := l.prod().NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: l.a.base, Table: failingTable(t)})
	if err != nil || got.N != 2 {
		t.Errorf("NextAttempt = %+v, %v; want attempt 2", got, err)
	}
}

// Nothing is written for a table that passed.
func TestNoFailedRowStartsNoAttemptOnTheRealRecords(t *testing.T) {
	l := newLoop(t, nil)
	n := len(l.events())
	pass := gate.Table{Base: l.a.base, Head: l.a.base, Rows: []gate.Row{{Kind: "fine", State: "active", Result: "pass", Reason: "—"}}}
	if _, err := l.prod().NextAttempt(context.Background(), NextAttemptInput{Task: "T-ab12", Base: l.a.base, Table: pass}); err == nil || len(l.events()) != n {
		t.Errorf("a passed table: %v, %d events more", err, len(l.events())-n)
	}
}
