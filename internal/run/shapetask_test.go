package run

import (
	"slices"
	"strings"
	"testing"
)

// loop.md, The requirement IDs of a Shape task: REQ for a requirement, NFR for a non-functional row, the same number, no ID
// for a success criterion, in the order of the draft.
func TestTheFinalIDsAreTheKindsOfThePinnedDraft(t *testing.T) {
	draft := [][]string{
		{"R-001", "requirement", "a", "1", "c", "S-aaaaaaaa"}, {"R-002", "non-functional", "b", "1", "c", "S-aaaaaaaa"},
		{"R-003", "success-criterion", "c", "1", "c", "S-aaaaaaaa"}, {"R-004", "requirement", "d", "1", "c", "S-aaaaaaaa"}}
	if got, want := finalIDs(draft), []string{"REQ-001", "NFR-002", "REQ-004"}; !slices.Equal(got, want) {
		t.Errorf("finalIDs = %v; want %v", got, want)
	}
	for _, c := range []struct {
		id, kind, want string
		ok             bool
	}{{"R-007", "requirement", "REQ-007", true}, {"R-012", "non-functional", "NFR-012", true}, {"R-003", "success-criterion", "", false}} {
		if got, ok := finalID(c.id, c.kind); got != c.want || ok != c.ok {
			t.Errorf("finalID(%s, %s) = %q, %v; want %q, %v", c.id, c.kind, got, ok, c.want, c.ok)
		}
	}
}

func TestTheIssueOfAShapeTaskHasItsTitleAndTheFilledBlockOfItsKind(t *testing.T) {
	if got := shapeIssueTitle("T-ab12", "specification"); got != "LAYUP task T-ab12: the specification" {
		t.Errorf("the title of a specification = %q", got)
	}
	if got := shapeIssueTitle("T-ab12", "architecture"); got != "LAYUP task T-ab12: the target's architecture" {
		t.Errorf("the title of an architecture task = %q", got)
	}
	sum := strings.Repeat("e", 64)
	for _, kind := range []string{"specification", "architecture"} {
		text := shapeIssueText("T-ab12", kind, []string{"REQ-001", "NFR-002"}, sum, "f00d")
		for _, want := range []string{"LAYUP task T-ab12: ", "Requirement IDs: REQ-001 NFR-002\n", "SHA-256 " + sum + ", read at the records commit f00d.\n", "docs/tasks/T-ab12.md"} {
			if !strings.Contains(text, want) {
				t.Errorf("%s: the text lacks %q:\n%s", kind, want, text)
			}
		}
		if strings.ContainsAny(text, "{}") {
			t.Errorf("%s: a placeholder is left:\n%s", kind, text)
		}
	}
}
