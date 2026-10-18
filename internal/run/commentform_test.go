package run

import (
	"reflect"
	"strings"
	"testing"
)

// The unit tests of the comment forms (docs/spec/loop.md, The fields and the
// blocks of a comment; task T-4pje, #396): every pinned rule has a case that
// breaks it and a control that differs from it by that rule alone.

const goodPlan = "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n- t1 internal/a_test.go\n- t2 internal/b_test.go\n"

var hexA, hexB = strings.Repeat("a", 40), strings.Repeat("b", 40)

// swap gives s with its one old replaced by new; a miss is a bug of the test.
func swap(t *testing.T, s, old, new string) string {
	t.Helper()
	if strings.Count(s, old) != 1 {
		t.Fatalf("%q is not once in %q", old, s)
	}
	return strings.Replace(s, old, new, 1)
}

// record is a review record with the head, the verdict and the bullets.
func record(head, verdict string, bullets ...string) string {
	out := "head: " + head + "\nverdict: " + verdict + "\nfindings:\n"
	for _, b := range bullets {
		out += "- " + b + "\n"
	}
	return out
}

const located = "internal/a.go:12 r-x: a text: with a colon"

type formCase struct {
	name, body string
	want       bool
}

func runCases(t *testing.T, ok func([]byte) bool, cases []formCase) {
	t.Helper()
	for _, c := range cases {
		if got := ok([]byte(c.body)); got != c.want {
			t.Errorf("%s: %v for %q; want %v", c.name, got, c.body, c.want)
		}
	}
}

func TestThePlanCommentParsesInThePinnedForm(t *testing.T) {
	p, tests, ok := ParsePlanComment([]byte("note: ignored\nthe tests: prose\n- stray bullet\n" + goodPlan + "\n- t3 after.go\n"))
	wantTests := []PlanTest{{"t1", "internal/a_test.go"}, {"t2", "internal/b_test.go"}}
	if !ok || p.Size != "small" || len(p.OpenQuestions)+len(p.NewInterfaces)+len(p.NewDependencies) != 0 || !reflect.DeepEqual(tests, wantTests) {
		t.Errorf("ParsePlanComment = %+v, %v, %v; want small, no lists, %v", p, tests, ok, wantTests)
	}
	reordered := "new dependencies: d\ntests:\n- t1 a.go\nsize: large\nnew interfaces: i\nopen questions: q\n"
	p, tests, ok = ParsePlanComment([]byte(reordered))
	if !ok || p.Size != "large" || !reflect.DeepEqual(p.NewDependencies, []string{"d"}) || !reflect.DeepEqual(p.NewInterfaces, []string{"i"}) || !reflect.DeepEqual(p.OpenQuestions, []string{"q"}) || len(tests) != 1 {
		t.Errorf("a plan in another order = %+v, %v, %v", p, tests, ok)
	}
	runCases(t, PlanFieldsOK, []formCase{
		{"the plan", goodPlan, true},
		{"a specification task, an empty block", swap(t, goodPlan, "- t1 internal/a_test.go\n- t2 internal/b_test.go\n", ""), true},
		{"a size that is neither", swap(t, goodPlan, "size: small", "size: medium"), false},
		{"a size in capitals", swap(t, goodPlan, "size: small", "size: Small"), false},
		{"a size of —", swap(t, goodPlan, "size: small", "size: —"), false},
	})
}

func TestAPlanFieldTwiceOrAbsentIsNoPlan(t *testing.T) {
	cases := []formCase{{"the control", goodPlan, true}}
	for _, line := range []string{"size: small\n", "open questions: —\n", "new interfaces: —\n", "new dependencies: —\n", "tests:\n"} {
		cases = append(cases,
			formCase{"absent " + line, swap(t, goodPlan, line, ""), false},
			formCase{"twice " + line, goodPlan + line, false})
	}
	runCases(t, PlanFieldsOK, cases)
}

func TestTheTestsMarkerOnceWithBulletsOfTwoWords(t *testing.T) {
	head := "size: small\nopen questions: —\nnew interfaces: —\nnew dependencies: —\n"
	runCases(t, PlanFieldsOK, []formCase{
		{"two words", head + "tests:\n- t1 a.go\n", true},
		{"a bullet of one word", head + "tests:\n- t1\n", false},
		{"a bullet of three words", head + "tests:\n- t1 a.go extra\n", false},
		{"a double space", head + "tests:\n- t1  a.go\n", false},
		{"a path with a space", head + "tests:\n- t1 a b.go\n", false},
		{"the same ID twice", head + "tests:\n- t1 a.go\n- t1 b.go\n", false},
		{"the same path twice", head + "tests:\n- t1 a.go\n- t2 a.go\n", true},
		{"an absolute path", head + "tests:\n- t1 /a.go\n", false},
		{"a path that leaves the tree", head + "tests:\n- t1 ../a.go\n", false},
		{"an unclean path", head + "tests:\n- t1 a/./b.go\n", false},
		{"a path into .git", head + "tests:\n- t1 .git/x\n", false},
		{"a bad bullet after the block ends is prose", head + "tests:\n- t1 a.go\nprose\n- bad bullet here\n", true},
		{"a bad bullet inside the block", head + "tests:\n- t1 a.go\n- bad\n- t2 b.go\n", false},
		{"a bullet that is not under the marker is prose", head + "- bad bullet here\ntests:\n- t1 a.go\n", true},
		{"tests: — is no marker", head + "tests: —\n", false},
		{"tests: x is no marker", head + "tests: x\n", false},
		{"a marker after a line of prose", head + "prose\ntests:\n- t1 a.go\n", true},
	})
}

func TestAPlanListFieldIsOneLineOfText(t *testing.T) {
	p, _, ok := ParsePlanComment([]byte(swap(t, goodPlan, "open questions: —", "open questions: a, b")))
	if !ok || !reflect.DeepEqual(p.OpenQuestions, []string{"a, b"}) {
		t.Errorf("a text with a comma = %v, %v; want the one element %q", p.OpenQuestions, ok, "a, b")
	}
	if p, _, ok = ParsePlanComment([]byte(goodPlan)); !ok || p.OpenQuestions != nil || PlanTier(p) != "execution" {
		t.Errorf("— = %v, %v; want the empty list and the execution tier", p.OpenQuestions, ok)
	}
	if p, _, ok = ParsePlanComment([]byte(swap(t, goodPlan, "new interfaces: —", "new interfaces: x"))); !ok || PlanTier(p) != "reasoning" {
		t.Errorf("a text in a list = %v; want the reasoning tier", ok)
	}
}

func TestAFieldValueIgnoresOneCarriageReturn(t *testing.T) {
	runCases(t, PlanFieldsOK, []formCase{
		{"one \\r after the size", swap(t, goodPlan, "size: small", "size: small\r"), true},
		{"trailing spaces after the size", swap(t, goodPlan, "size: small", "size: small  "), true},
		{"two \\r after the size", swap(t, goodPlan, "size: small", "size: small\r\r"), false},
		{"a letter after the size", swap(t, goodPlan, "size: small", "size: smallX"), false},
		{"a \\r after —", swap(t, goodPlan, "open questions: —", "open questions: —\r"), true},
	})
	if p, _, _ := ParsePlanComment([]byte(swap(t, goodPlan, "open questions: —", "open questions: —\r"))); p.OpenQuestions != nil {
		t.Errorf("— and a \\r = %v; want the empty list", p.OpenQuestions)
	}
	runCases(t, VerdictFieldOK, []formCase{{"a verdict and a \\r", "verdict: approve\r\n", true}})
}

func TestAMarkerOrBulletWithACarriageReturnIsNotRead(t *testing.T) {
	runCases(t, PlanFieldsOK, []formCase{
		{"the marker with a \\r is no marker", swap(t, goodPlan, "tests:\n", "tests:\r\n"), false},
		{"a test bullet with a \\r", swap(t, goodPlan, "- t1 internal/a_test.go\n", "- t1 internal/a_test.go\r\n"), false},
		{"a test bullet, the control", goodPlan, true},
	})
	runCases(t, FindingsLocatedOK, []formCase{
		{"a finding with a \\r in its text", record(hexA, "material", "internal/a.go:1 r: text\r"), false},
		{"a finding with a \\r in its path", record(hexA, "material", "internal/a.go\r:1 r: text"), false},
		{"a finding, the control", record(hexA, "material", "internal/a.go:1 r: text"), true},
		{"the marker with a \\r is no marker", swap(t, record(hexA, "material", located), "findings:\n", "findings:\r\n"), false},
	})
}

func TestAFieldWithNothingAfterTheColonIsMalformed(t *testing.T) {
	runCases(t, PlanFieldsOK, []formCase{
		{"nothing after the space", swap(t, goodPlan, "size: small", "size: "), false},
		{"nothing after the colon", swap(t, goodPlan, "size: small", "size:"), false},
		{"spaces only", swap(t, goodPlan, "open questions: —", "open questions:   "), false},
		{"the control", goodPlan, true},
	})
	runCases(t, VerdictFieldOK, []formCase{{"a verdict with nothing", "verdict: \n", false}, {"a verdict, the control", "verdict: approve\n", true}})
}

func TestTheLineReaderIgnoresAFencedBlock(t *testing.T) {
	runCases(t, PlanFieldsOK, []formCase{
		{"no fence, a second size", goodPlan + "size: large\n", false},
		{"a second size in a backtick fence", goodPlan + "```\nsize: large\n```\n", true},
		{"a second size in a tilde fence", goodPlan + "~~~\nsize: large\n~~~\n", true},
		{"an info string on the opener", goodPlan + "```text a b\nsize: large\n```\n", true},
		{"a longer closer closes", swap(t, goodPlan, "size: small\n", "") + "```\nx\n`````\nsize: small\n", true},
		{"a tests: marker in a fence", goodPlan + "```\ntests:\n```\n", true},
		{"the size is in a fence", "```\nsize: small\n```\n" + swap(t, goodPlan, "size: small\n", ""), false},
		{"a closer of another character does not close", goodPlan + "```\nx\n~~~\nsize: large\n", true},
		{"a shorter closer does not close", goodPlan + "````\nx\n```\nsize: large\n", true},
		{"an opener with a \\r opens", "```\r\n" + goodPlan, false},
		{"a closer with a \\r does not close, a field after it", goodPlan + "```\nx\n```\r\nsize: large\n", true},
		{"an unclosed fence runs to the end", goodPlan + "```\nsize: large\n", true},
		{"an unclosed fence hides what follows", swap(t, goodPlan, "size: small\n", "") + "```\nsize: small\n", false},
		{"a fence indented is no fence", goodPlan + " ```\nsize: large\n", false},
		{"two backticks are no fence", goodPlan + "``\nsize: large\n", false},
		{"a fence ends the block", swap(t, goodPlan, "- t2 internal/b_test.go\n", "```\n```\n- bad\n"), true},
	})
}

func TestACommentThatFailsItsFormIsFalseEverywhere(t *testing.T) {
	good := record(hexA, "nothing material in scope")
	if !NothingMaterialOK([]byte(good), hexA) {
		t.Fatal("the control record is not valid")
	}
	bad := good + "head: " + hexB + "\n"
	if NothingMaterialOK([]byte(bad), hexA) || NothingMaterialOK([]byte(bad), hexB) || FindingsLocatedOK([]byte(record(hexA, "material", located)+"head: "+hexB+"\n")) {
		t.Error("a head twice is read for an item")
	}
	if FindingsLocatedOK([]byte(record("A"+hexA[1:], "material", located))) || FindingsLocatedOK([]byte(record(hexA, "bogus", located))) {
		t.Error("a bad head or verdict does not fail findings-located")
	}
	if FindingsLocatedOK([]byte(record(hexA, "material", located, "bad"))) || !FindingsLocatedOK([]byte(record(hexA, "material", located))) {
		t.Error("a malformed bullet does not fail the whole record")
	}
	if _, ok := ParseReviewRecord([]byte(record(hexA, "bogus"))); ok {
		t.Error("ParseReviewRecord reads a bad verdict")
	}
}

func TestVerdictFieldIsApproveOrApproveWithConditions(t *testing.T) {
	runCases(t, VerdictFieldOK, []formCase{
		{"approve", "verdict: approve\n", true},
		{"approve-with-conditions", "x\nverdict: approve-with-conditions\n", true},
		{"reject", "verdict: reject\n", false},
		{"absent", "no verdict here\n", false},
		{"twice", "verdict: approve\nverdict: approve\n", false},
		{"—", "verdict: —\n", false},
		{"a word of nobody's", "verdict: approved\n", false},
		{"in a fence", "```\nverdict: approve\n```\n", false},
		{"a fenced second verdict is not read", "verdict: approve\n```\nverdict: reject\n```\n", true},
		{"a value after a table row", "| Verdict | approve |\nverdict: approve\n", true},
	})
	if v, ok := ParsePlanReview([]byte("verdict: reject\n")); !ok || v != "reject" {
		t.Errorf("ParsePlanReview = %q, %v; want reject, true (a parse, not an approval)", v, ok)
	}
}

func TestFindingsLocatedNeedsALocatedFinding(t *testing.T) {
	one := func(b string) string { return record(hexA, "material", b) }
	runCases(t, FindingsLocatedOK, []formCase{
		{"one finding", one(located), true},
		{"line 1", one("p.go:1 r: t"), true},
		{"a 64-hex head", record(strings.Repeat("c", 64), "material", located), true},
		{"not mergeable with a finding", record(hexA, "not mergeable, findings recorded", located), true},
		{"line 0", one("p.go:0 r: t"), false},
		{"a negative line", one("p.go:-1 r: t"), false},
		{"a line that is no integer", one("p.go:x r: t"), false},
		{"a signed line", one("p.go:+1 r: t"), false},
		{"no line", one("p.go r: t"), false},
		{"a rule of two words", one("p.go:1 my rule: t"), false},
		{"a colon inside the rule", one("p.go:1 r:x: t"), false},
		{"a rule with a hyphen and a letter, the control", one("p.go:1 rx: t"), true},
		{"an empty rule", one("p.go:1 : t"), false},
		{"no colon after the rule", one("p.go:1 r t"), false},
		{"an empty text", one("p.go:1 r: "), false},
		{"a path with a space", one("a b.go:1 r: t"), false},
		{"an absolute path", one("/a.go:1 r: t"), false},
		{"a path that leaves the tree", one("../a.go:1 r: t"), false},
		{"the empty block", record(hexA, "material"), false},
		{"not mergeable with an empty block", record(hexA, "not mergeable, findings recorded"), false},
		{"the marker absent", "head: " + hexA + "\nverdict: material\n- " + located + "\n", false},
		{"the marker twice", record(hexA, "material", located) + "findings:\n- " + located + "\n", false},
		{"a head of 39 characters", record(hexA[1:], "material", located), false},
		{"a head in capitals", record(strings.ToUpper(hexA), "material", located), false},
		{"the head absent", "verdict: material\nfindings:\n- " + located + "\n", false},
		{"the verdict absent", "head: " + hexA + "\nfindings:\n- " + located + "\n", false},
	})
	r, ok := ParseReviewRecord([]byte(record(hexA, "material", located, "q/b.go:3 r2: t2")))
	want := []Finding{{Rule: "r-x", File: "internal/a.go", Line: 12, Text: "a text: with a colon"}, {Rule: "r2", File: "q/b.go", Line: 3, Text: "t2"}}
	if !ok || r.Verdict != "material" || r.Head != hexA || !reflect.DeepEqual(r.Findings, want) {
		t.Errorf("ParseReviewRecord = %+v, %v; want the verdict, the head and %v", r, ok, want)
	}
}

func TestNothingMaterialNeedsTheVerdictTheHeadAndAnEmptyBlock(t *testing.T) {
	nm := record(hexA, "nothing material in scope")
	for _, c := range []struct {
		name, body, head string
		want             bool
	}{
		{"the control", nm, hexA, true},
		{"64 hexadecimal characters", record(strings.Repeat("d", 64), "nothing material in scope"), strings.Repeat("d", 64), true},
		{"another head", nm, hexB, false},
		{"no head held", nm, "", false},
		{"the head in capitals", record(strings.ToUpper(hexA), "nothing material in scope"), strings.ToUpper(hexA), false},
		{"a finding in the block", record(hexA, "nothing material in scope", located), hexA, false},
		{"material", record(hexA, "material"), hexA, false},
		{"not mergeable with an empty block", record(hexA, "not mergeable, findings recorded"), hexA, false},
		{"the marker absent", "head: " + hexA + "\nverdict: nothing material in scope\n", hexA, false},
		{"the marker with a value", "head: " + hexA + "\nverdict: nothing material in scope\nfindings: none\n", hexA, false},
		{"a verdict of —", record(hexA, "—"), hexA, false},
	} {
		if got := NothingMaterialOK([]byte(c.body), c.head); got != c.want {
			t.Errorf("%s: %v; want %v", c.name, got, c.want)
		}
	}
}

func TestOpenMaterialFailsUphill(t *testing.T) {
	parse := func(body string) (ReviewRecord, bool) { return ParseReviewRecord([]byte(body)) }
	for _, c := range []struct {
		name, body string
		want       bool
	}{
		{"nothing material and an empty block", record(hexA, "nothing material in scope"), false},
		{"nothing material and a finding", record(hexA, "nothing material in scope", located), true},
		{"material and an empty block", record(hexA, "material"), true},
		{"not mergeable and an empty block", record(hexA, "not mergeable, findings recorded"), true},
		{"material and a finding", record(hexA, "material", located), true},
		{"a record that does not parse", "nothing material in scope\n", true},
		{"a record with a bad bullet", record(hexA, "nothing material in scope", "bad"), true},
		{"an empty comment", "", true},
	} {
		if got := OpenMaterial(parse(c.body)); got != c.want {
			t.Errorf("%s: %v; want %v", c.name, got, c.want)
		}
	}
	if !OpenMaterial(ReviewRecord{Verdict: "nothing material in scope"}, false) {
		t.Error("a record that did not parse reads as none open")
	}
}

func TestAnEmptyOrInvalidCommentParsesToNothing(t *testing.T) {
	for _, body := range []string{"", "\xff\xfe", goodPlan + "\xff", "\n\n"} {
		b := []byte(body)
		if PlanFieldsOK(b) || VerdictFieldOK(b) || FindingsLocatedOK(b) || NothingMaterialOK(b, hexA) {
			t.Errorf("%q is read for an item", body)
		}
		if _, tests, ok := ParsePlanComment(b); ok || tests != nil {
			t.Errorf("ParsePlanComment(%q) = %v, %v", body, tests, ok)
		}
		if _, ok := ParsePlanReview(b); ok {
			t.Errorf("ParsePlanReview(%q) is ok", body)
		}
		if _, ok := ParseReviewRecord(b); ok {
			t.Errorf("ParseReviewRecord(%q) is ok", body)
		}
	}
	if !PlanFieldsOK([]byte(goodPlan)) {
		t.Error("the control plan is not read")
	}
}
