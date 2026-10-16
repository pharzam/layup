package run

import (
	"slices"
	"testing"
)

func question(label, needs string, records ...string) []string {
	reason := label + " " + needs
	for _, r := range records {
		reason += " " + r
	}
	return []string{"question", "1", "May a draft be marked ready?", "", reason}
}

func TestAQuestionRowIsRead(t *testing.T) {
	got := QuestionsOf([][]string{{"status", "1", "needs_context", "", "asks"}, question("interface-contract", "no", "docs/spec/forge.md", "T-cr8j"),
		{"question", "2", "Which?", "", "some-label maybe"}})
	want := []Question{{"May a draft be marked ready?", "interface-contract", "no", []string{"docs/spec/forge.md", "T-cr8j"}}, {"Which?", "some-label", "maybe", nil}}
	if !slices.EqualFunc(got, want, func(a, b Question) bool {
		return a.Text == b.Text && a.Label == b.Label && a.NeedsHuman == b.NeedsHuman && slices.Equal(a.Records, b.Records)
	}) {
		t.Errorf("QuestionsOf = %+v, want %+v", got, want)
	}
}

func TestTheRouteOfEachKindWithNoHumanNeeded(t *testing.T) {
	for label, role := range map[string]string{"domain-or-business": "Domain Expert", "architecture-boundary": "Systems Architect",
		"interface-contract": "Software Architect", "environment-or-infrastructure": "Software Engineer"} {
		got := RouteQuestion(Question{Label: label, NeedsHuman: "no", Records: []string{"R-001"}}, DefaultOwners())
		if got != (Route{Role: role}) {
			t.Errorf("%s: %+v; want the owner role %s", label, got, role)
		}
	}
}

func TestAQuestionThatNeedsAHumanGoesToTheOperator(t *testing.T) {
	for _, word := range []string{"yes", "maybe", "YES", "No", "true", "0"} {
		got := RouteQuestion(Question{Label: "interface-contract", NeedsHuman: word}, DefaultOwners())
		if got != (Route{Role: "operator", Human: true}) {
			t.Errorf("needs_human %q: %+v; want the Operator", word, got)
		}
	}
}

func TestALabelOutsideTheFourGoesToTheOperator(t *testing.T) {
	for _, label := range []string{"unknown", "Interface-Contract", "interface_contract", "domain", ""} {
		got := RouteQuestion(Question{Label: label, NeedsHuman: "no"}, DefaultOwners())
		if got != (Route{Role: "operator", Human: true}) {
			t.Errorf("label %q with no human needed: %+v; want the Operator", label, got)
		}
	}
}

func TestAReplacedOwnerMapRoutesByItsRoles(t *testing.T) {
	owners := map[string]string{"interface-contract": "Release Engineer"}
	if got := RouteQuestion(Question{Label: "interface-contract", NeedsHuman: "no"}, owners); got != (Route{Role: "Release Engineer"}) {
		t.Errorf("a replaced map: %+v; want its role", got)
	}
	if got := RouteQuestion(Question{Label: "domain-or-business", NeedsHuman: "no"}, owners); !got.Human {
		t.Errorf("a label the replaced map lacks: %+v; want the Operator", got)
	}
	d := DefaultOwners()
	d["interface-contract"] = "changed"
	if DefaultOwners()["interface-contract"] != "Software Architect" {
		t.Error("DefaultOwners gives a map that a caller changes for the next caller")
	}
}

func TestQuestionRowNeedsARecordAndAQuestion(t *testing.T) {
	status := []string{"status", "1", "needs_context", "", "asks"}
	for _, c := range []struct {
		name   string
		result [][]string
		want   bool
	}{
		{"one question with a record", [][]string{status, question("interface-contract", "no", "T-cr8j")}, true},
		{"two questions each with records", [][]string{status, question("a", "no", "T-1"), question("b", "yes", "T-2", "T-3")}, true},
		{"no question row", [][]string{status}, false},
		{"no row at all", nil, false},
		{"a question with a label and a word only", [][]string{status, question("interface-contract", "no")}, false},
		{"one question with a record and one without", [][]string{status, question("a", "no", "T-1"), question("b", "no")}, false},
		{"a decision row with records is no question", [][]string{status, {"decision", "1", "x", "", "because T-1 T-2 T-3"}}, false},
	} {
		if got := QuestionRowOK(c.result); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}

func TestAnAnswerRowNeedsAnAnswer(t *testing.T) {
	status := []string{"status", "1", "completed", "", "answered"}
	answer := []string{"answer", "1", "Yes", "", "docs/spec/forge.md"}
	for _, c := range []struct {
		name   string
		result [][]string
		want   bool
	}{
		{"one answer", [][]string{status, answer}, true},
		{"no answer row", [][]string{status}, false},
		{"no row at all", nil, false},
		{"a question row is no answer", [][]string{status, question("a", "no", "T-1")}, false},
		{"an answer with no record", [][]string{status, {"answer", "1", "Yes", "", ""}}, false},
	} {
		if got := AnswerRowOK(c.result); got != c.want {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
	}
}
