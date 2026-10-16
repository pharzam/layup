//go:build integration

package run

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/tsv"
)

// The Retrospective step on the real adapter, the httptest forge and real git
// (task T-v2sv, #359, row 66c of the plan of M2e). The helper is unwired — no
// caller exists in M2e — so a test calls what a later step loop may call, on
// the state of a restart.

// The Go value of the block retrospectives equals its block of
// docs/spec/records.md.
func TestTheSchemaOfRetrospectivesEqualsItsBlock(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS("../../docs/spec"))
	if err != nil {
		t.Fatal(err)
	}
	block, ok := blocks["retrospectives"]
	if !ok {
		t.Fatal("docs/spec/ has no block retrospectives")
	}
	if err := tsv.Compare(block, RetrospectivesSchema); err != nil {
		t.Error(err)
	}
}

// One records commit holds the event of the milestone and nothing else; on a
// restart the row exists and the step writes nothing.
func TestTheRetrospectiveIsOneCommitAndDoneBeforeOnARestart(t *testing.T) {
	w := registerWorld(t, loopPerms)
	r := restarted(t, w)
	ctx := context.Background()
	head := recordsHead(t, w)
	next, err := r.retrospective(ctx, "M2e")
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" {
		t.Fatalf("the step the loop enters next: got %q, want %q", next, "bet")
	}
	head2 := recordsHead(t, w)
	if head2 == head {
		t.Fatal("the event's commit did not move the records head")
	}
	// The commit's tree holds retrospectives.tsv alone: one file, one row.
	got := gitOut(t, w.bare, "show", "--name-only", "--format=", head2)
	if files := strings.Fields(got); len(files) != 1 || files[0] != "retrospectives.tsv" {
		t.Fatalf("the commit's files: got %v, want [retrospectives.tsv]", files)
	}
	text := show(t, w.bare, head2, "retrospectives.tsv")
	row := strings.Split(strings.TrimSpace(text), "\n")
	if len(row) != 2 || row[0] != "milestone\ttime" {
		t.Fatalf("the table:\n%s", text)
	}
	fields := strings.Split(row[1], "\t")
	if fields[0] != "M2e" {
		t.Fatalf("the milestone: got %q, want M2e", fields[0])
	}
	// The world's clock has sub-second parts; the row is second-precision UTC.
	if _, err := time.Parse(timeForm, fields[1]); err != nil {
		t.Fatalf("the row's time %q: %v", fields[1], err)
	}
	// A restart: the row exists, so the step is done before — no second write.
	r2 := restarted(t, w)
	next, err = r2.retrospective(ctx, "M2e")
	if err != nil {
		t.Fatal(err)
	}
	if next != "bet" {
		t.Fatalf("on the restart: got %q, want %q", next, "bet")
	}
	if h := recordsHead(t, w); h != head2 {
		t.Fatalf("a done-before step wrote: the head moved to %s", h)
	}
}
