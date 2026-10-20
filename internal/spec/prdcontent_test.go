package spec

import (
	"slices"
	"testing"
)

// ContentOfPRDs (task T-h65e, #418): the words of the rows of the PRD, by the parse of the full check.
func TestContentOfPRDsGivesTheWordsOfTheRows(t *testing.T) {
	prd := "# PRD-0001\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| --- | --------- | ------ | ----- | ----- |\n" +
		"| REQ-001 | Alpha |  |  | F-0003#1 |\n| NFR-002 | p99 \\| under 200 ms |  |  | F-0003#2 |\n\n" +
		"| Requirement | Acceptance criterion |\n| ----------- | -------------------- |\n| `REQ-001` | a check |\n| `NFR-002` | measured |\n\n" +
		"## 9. Success metrics\n| Metric | Fact |\n| ------ | ---- |\n| faster | F-0003#3 |\n| cheaper | F-0003#4 |\n"
	got := ContentOfPRDs(map[string][]byte{"PRD-0001-x.md": []byte(prd)})
	if got.Statements["REQ-001"] != "Alpha" || got.Statements["NFR-002"] != "p99 | under 200 ms" || len(got.Statements) != 2 {
		t.Errorf("Statements = %v; want REQ-001 Alpha and NFR-002 with its escaped bar", got.Statements)
	}
	if got.Criteria["REQ-001"] != "a check" || got.Criteria["NFR-002"] != "measured" || len(got.Criteria) != 2 {
		t.Errorf("Criteria = %v; want the two criteria", got.Criteria)
	}
	if !slices.Equal(got.Metrics, []string{"faster", "cheaper"}) {
		t.Errorf("Metrics = %v; want [faster cheaper]", got.Metrics)
	}
	if empty := ContentOfPRDs(nil); len(empty.Statements) != 0 || len(empty.Metrics) != 0 {
		t.Errorf("no PRD: %+v; want nothing", empty)
	}
}
