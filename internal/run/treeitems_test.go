package run

// The tree items of a handoff, at the unit level (loop.md, The items a caller
// computes; task T-yeav, #404): no clone is read, so the head is none, or the
// row does not list the item.

import (
	"context"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/rules"
)

func TestSectionsOfNoRequirementIDAreFalse(t *testing.T) {
	files := map[string][]byte{"a.md": []byte("## REQ-005 — x\nbody\n## REQ-006 — y\nbody\n")}
	for _, c := range []struct {
		name string
		ids  []string
		want bool
	}{
		{"no ID", nil, false}, {"an empty list", []string{}, false},
		{"one ID, present (the control)", []string{"REQ-005"}, true},
		{"two IDs, both present", []string{"REQ-005", "REQ-006"}, true},
		{"two IDs, one missing", []string{"REQ-005", "REQ-007"}, false},
	} {
		if got := specSections(files, c.ids); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
	if specSections(nil, []string{"REQ-005"}) {
		t.Error("no file holds a section")
	}
}

func TestTheRequirementIDsAreTheRegisterRowOfTheTask(t *testing.T) {
	w := newHandoffWorld(t)
	if task, err := w.s.taskRow(context.Background(), "T-ab12"); err != nil || strings.Join(task.Requirements, " ") != "REQ-005" || task.Issue != taskIssue {
		t.Errorf("taskRow = %+v, %v; want the row of T-ab12 with REQ-005", task, err)
	}
	if _, err := w.s.taskRow(context.Background(), "T-zzzz"); err == nil {
		t.Error("a task with no row has one")
	}
	w.store.files["tasks.tsv"] = []byte("not a register\n")
	if _, err := w.s.taskRow(context.Background(), "T-ab12"); err == nil {
		t.Error("a corrupt register has a row")
	}
}

var h40, h64 = strings.Repeat("d", 40), strings.Repeat("c", 64)

func specHandoff(items map[string]bool) HandoffInput {
	return HandoffInput{Task: "T-ab12", Session: planSess, From: "developer (specification)", To: "verifier", Items: items,
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: h40, Found: h40}, {Kind: "task-file", Named: h64, Found: h64}}}
}

func renderedHandoff(items map[string]bool) HandoffInput {
	return HandoffInput{Task: "T-ab12", Attempt: 1, From: "code (rendered task)", To: "verifier", CodeWritten: true, Items: items,
		Artifacts: []handoff.Artifact{{Kind: "rendering", Path: "docs/facts/numbered.md", Named: h64, Found: h64}}}
}

func TestATreeItemWithNoHeadIsFalseAndNotAnError(t *testing.T) {
	both := map[string]bool{"spec-sections": true, "rendering-equal": true}
	// A session handoff has the head of its bound event, never the caller's Head,
	// whether or not it has one: this session has none.
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	in := specHandoff(both)
	in.Head = h40
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid spec-sections" {
		t.Errorf("a session with no bound head, the caller's Head and literal: %+v, %v; want invalid spec-sections", v, err)
	}
	if ev := w.eventsNow(t); ev[len(ev)-1][1] != "handoff" || ev[len(ev)-1][6] != "invalid spec-sections" {
		t.Errorf("the last event %q; want the handoff invalid spec-sections", ev[len(ev)-1])
	}
	// A handoff that code wrote has the Head that code sets; none is false.
	w = newHandoffWorld(t)
	if v, err := w.s.RecordHandoff(context.Background(), renderedHandoff(both)); err != nil || v.Detail != "invalid rendering-equal" {
		t.Errorf("code-written with no Head, the caller's literal true: %+v, %v; want invalid rendering-equal", v, err)
	}
	if v, err := w.s.RecordHandoff(context.Background(), renderedHandoff(nil)); err != nil || v.Detail != "invalid rendering-equal" {
		t.Errorf("code-written with no Head, no literal: %+v, %v; want invalid rendering-equal", v, err)
	}
}

func TestTheTreeItemsAreComputedOnlyWhereTheRowListsThem(t *testing.T) {
	both := map[string]bool{"spec-sections": true, "rendering-equal": true}
	// The default row of the architecture task lists neither item: nothing is
	// read (the clone of this world does not exist), and the caller's values are not read.
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	in := specHandoff(both)
	in.From = "developer (architecture)"
	in.Head = h40
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "valid" {
		t.Fatalf("a row without the item: %+v, %v; want valid with no read", v, err)
	}
	// A code-written row without rendering-equal reads no tree either: this
	// handoff has a head, and the records and artifact that would reach git.
	rows, err := handoff.ReadTable(w.store.files["transitions.tsv"])
	if err != nil {
		t.Fatal(err)
	}
	for i := range rows {
		if rows[i].From == "code (rendered task)" {
			rows[i].Items = []string{"artifacts"}
		}
	}
	var plain strings.Builder
	if err := handoff.WriteTable(&plain, rows); err != nil {
		t.Fatal(err)
	}
	w.store.files["transitions.tsv"] = []byte(plain.String())
	code := renderedHandoff(both)
	code.Head, code.Render = h40, []rules.Record{{Path: "docs/facts/numbered.md", Kind: rules.Facts}}
	if v, err := w.s.RecordHandoff(context.Background(), code); err != nil || v.Detail != "valid" {
		t.Fatalf("a code-written row without the item: %+v, %v; want valid with no read", v, err)
	}
	// The same row, replaced by one that lists spec-sections: computed there, and false with no head.
	if rows, err = handoff.ReadTable(w.store.files["transitions.tsv"]); err != nil {
		t.Fatal(err)
	}
	for i := range rows {
		if rows[i].From == "developer (architecture)" {
			rows[i].Items = append(rows[i].Items, "spec-sections")
		}
	}
	var b strings.Builder
	if err := handoff.WriteTable(&b, rows); err != nil {
		t.Fatal(err)
	}
	w.store.files["transitions.tsv"] = []byte(b.String())
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid spec-sections" {
		t.Errorf("a replaced row that lists the item: %+v, %v; want invalid spec-sections", v, err)
	}
}

func TestAnArtifactWhoseHashDiffersIsInvalidArtifacts(t *testing.T) {
	w := newHandoffWorld(t)
	w.session(t, planSess, 1, "developer", "")
	in := specHandoff(map[string]bool{"spec-sections": true})
	in.Artifacts[0].Found = strings.Repeat("e", 40)
	if v, err := w.s.RecordHandoff(context.Background(), in); err != nil || v.Detail != "invalid artifacts" {
		t.Errorf("a hash that differs: %+v, %v; want invalid artifacts", v, err)
	}
}
