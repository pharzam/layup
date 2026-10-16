package run

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

// answersFixture is a batch of four questions of the idea owner's block, asked
// at 10:00:00: a rule gap and the three fields of the intent form.
var answersQuestions = [][]string{
	{"Q-001", "rule", "Q-001", "idea-owner", "before-delivery", "G1", "", "", "Which stack?", "2026-10-11T10:00:00Z", "layup psb check", "", "asked"},
	{"Q-002", "intent", "band.B", "idea-owner", "before-delivery", "intent", "", "", "B?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
	{"Q-003", "intent", "band.U", "idea-owner", "before-delivery", "intent", "", "", "U?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
	{"Q-004", "intent", "appetite.hours", "idea-owner", "before-delivery", "intent", "", "", "Hours?", "2026-10-11T10:00:00Z", "layup run", "", "asked"},
}

var answersApprovers = [][]string{
	{"303", "idea-owner", "carol", "2026-10-08T12:00:00Z", "start"},
	{"101", "operator", "pharzam", "2026-10-08T12:00:00Z", "start"},
}

// copyOf gives a copy of a comment of issue 1: a row of copies.tsv and its body.
func copyOf(id, seen, author, app, created, body string) answerCopy {
	return answerCopy{Row: []string{id, seen, "1", author, "x", app, created, created, strings.Repeat("a", 64), "copies/" + id + "-" + seen + ".md"}, Body: body}
}

// The answers (docs/spec/intake.md, The answers): a line Q-NNN: <answer> of a
// first copy after the batch by an approver in the role of the question's
// block answers it; the last such line of an ID is its answer; a block waits
// until such a comment holds at least one answer line.
func TestCollectAnswers(t *testing.T) {
	ignored := []answerCopy{
		copyOf("701", "1", "303", "", "2026-10-11T09:59:59Z", "Q-001: before the batch\n"),
		copyOf("702", "1", "303", "", "2026-10-11T10:00:00Z", "Q-001: at the time of the batch\n"),
		copyOf("703", "1", "101", "", "2026-10-11T10:01:00Z", "Q-001: the Operator, in the other block's role\n"),
		copyOf("704", "1", "999", "", "2026-10-11T10:02:00Z", "Q-001: not an approver\n"),
		copyOf("705", "1", "303", "layup-agent", "2026-10-11T10:03:00Z", "Q-001: by an App\n"),
		copyOf("706", "1", "303", "", "2026-10-11T10:04:00Z", "Thank you, I will answer soon.\n"),
	}
	if answers, waiting := collectAnswers(answersQuestions, answersApprovers, ignored, "1"); len(answers) != 0 || strings.Join(waiting, " ") != "idea-owner" {
		t.Fatalf("no answer line of an approver in the role: %q, waiting %q; want none, the idea owner's block waiting", answers, waiting)
	}
	copies := append(ignored,
		copyOf("707", "1", "303", "", "2026-10-11T10:05:00Z",
			"Q-001: first\nQ-002: 9.5\r\nQ-003: 10.0 \nQ-004: 40.0\nQ-009: no such question\n xQ-001: not at the start\nQ-001:no space\n"),
		// A later line corrects an earlier one; an empty answer and one of
		// exactly — after it answer nothing.
		copyOf("708", "1", "303", "", "2026-10-11T10:06:00Z", "Q-001: second\nQ-001: \nQ-001: —\n"),
		copyOf("707", "2", "303", "", "2026-10-11T10:07:00Z", "Q-001: an edit, never an answer\n"),
	)
	answers, waiting := collectAnswers(answersQuestions, answersApprovers, copies, "1")
	want := [][]string{
		{"Q-001", "second", "708", "303", "2026-10-11T10:06:00Z"},
		{"Q-002", "9.5", "707", "303", "2026-10-11T10:05:00Z"},   // the CR before the line feed trimmed
		{"Q-003", "10.0 ", "707", "303", "2026-10-11T10:05:00Z"}, // trailing spaces kept
		{"Q-004", "40.0", "707", "303", "2026-10-11T10:05:00Z"},
	}
	if len(waiting) != 0 || len(answers) != len(want) {
		t.Fatalf("the answers %q, waiting %q; want %q", answers, waiting, want)
	}
	for i := range want {
		if strings.Join(answers[i], "|") != strings.Join(want[i], "|") {
			t.Errorf("answer %d: %q, want %q", i, answers[i], want[i])
		}
	}
}

// The intent form parses when each field is a decimal and B is below U, as
// numbers; then budget.tsv has its three rows, each with its question.
func TestBudgetOf(t *testing.T) {
	answer := func(b, u, h string) [][]string {
		return [][]string{{"Q-002", b, "707", "303", "t"}, {"Q-003", u, "707", "303", "t"}, {"Q-004", h, "707", "303", "t"}}
	}
	rows, ok := budgetOf(answersQuestions, answer("9.5", "10.0", "40.0"))
	if !ok || len(rows) != 3 || strings.Join(rows[0], " ") != "band.B 9.5 Q-002" || strings.Join(rows[1], " ") != "band.U 10.0 Q-003" ||
		strings.Join(rows[2], " ") != "appetite.hours 40.0 Q-004" {
		t.Errorf("9.5, 10.0, 40.0: %q, %v; want the three rows (9.5 below 10.0 as numbers)", rows, ok)
	}
	// Spaces around a value are trimmed by the parse; the answer itself is kept
	// byte for byte in answers.tsv.
	if rows, ok := budgetOf(answersQuestions, answer(" 9.5", "10.0 ", "40.0")); !ok || rows[0][1] != "9.5" || rows[1][1] != "10.0" {
		t.Errorf("spaces around the values: %q, %v; want them parsed", rows, ok)
	}
	for _, c := range []struct{ name, b, u, h string }{
		{"B not below U", "10.0", "9.5", "40.0"},
		{"B equal to U", "10.0", "10.0", "40.0"},
		{"a value that is no decimal", "9.5", "10.0", "40"},
		{"a value with a word", "9.5 USD", "10.0", "40.0"},
	} {
		if rows, ok := budgetOf(answersQuestions, answer(c.b, c.u, c.h)); ok || rows != nil {
			t.Errorf("%s: %q, %v; want budget.tsv not written", c.name, rows, ok)
		}
	}
	if rows, ok := budgetOf(answersQuestions, answer("9.5", "10.0", "40.0")[:2]); ok || rows != nil {
		t.Errorf("a field with no answer: %q, %v; want budget.tsv not written", rows, ok)
	}
}

// The bindings of the planned points (docs/spec/loop.md, The approvers and the
// success criteria): a line <point>: <login> of a counted comment — a first
// copy, after the batch's asked, a decision of the idea owner — names the
// point's approver, and the last counted line wins; a name outside the set, a
// comment of another block or of an account with no approvers row, an edit and
// an App's comment count for nothing (each is input only, and answers
// nothing); and the earliest comment that named a login is the source of the
// row the step writes.
func TestThePointBindings(t *testing.T) {
	copies := []answerCopy{
		copyOf("701", "1", "303", "", "2026-10-11T10:05:00Z", "setup: pharzam\nbet: alice\n"),
		copyOf("702", "1", "101", "", "2026-10-11T10:06:00Z", "setup: mallory\n"),                          // the Operator's block: never
		copyOf("703", "1", "999", "", "2026-10-11T10:07:00Z", "bet: nobody\n"),                             // an account with no row: input only
		copyOf("701", "2", "303", "", "2026-10-11T10:08:00Z", "bet: edited\n"),                             // an edit is never a decision
		copyOf("704", "1", "303", "", "2026-10-11T10:09:00Z", "color: blue\nbet: mallory\naudit: alice\n"), // an unlisted name is input only; bet renames, alice serves twice
		copyOf("705", "1", "303", "layup-agent", "2026-10-11T10:10:00Z", "audit: ghost\n"),                 // an App's comment: never
		copyOf("706", "1", "303", "", "2026-10-11T09:59:00Z", "setup: early\n"),                            // before the batch
		{Row: []string{"707", "1", "77", "303", "x", "", "2026-10-11T10:11:00Z", "2026-10-11T10:11:00Z", "d", "copies/707-1.md"}, // another issue
			Body: "audit: ghost\n"},
	}
	bound, named := pointBindings(copies, answersApprovers, "2026-10-11T10:00:00Z", "1")
	want := map[string]string{"setup": "pharzam", "bet": "mallory", "audit": "alice"}
	if len(bound) != len(want) {
		t.Fatalf("the bindings %v; want %v", bound, want)
	}
	for p, l := range want {
		if bound[p] != l {
			t.Errorf("the point %q binds %q; want %q", p, bound[p], l)
		}
	}
	// named holds, for each login a counted line named, the earliest comment
	// that named it — the source and the time of the row the step writes.
	if named["pharzam"] != [2]string{"701", "2026-10-11T10:05:00Z"} ||
		named["alice"] != [2]string{"701", "2026-10-11T10:05:00Z"} ||
		named["mallory"] != [2]string{"704", "2026-10-11T10:09:00Z"} || len(named) != 3 {
		t.Errorf("the named logins %v; want pharzam/alice by 701 and mallory by 704", named)
	}
}

// The rows the step appends for the logins a counted line names and no row
// holds (docs/spec/loop.md, The approvers and the success criteria): one row
// each, the ID from the forge, since the naming comment's time and source the
// earliest comment that named it; the Operator's and the idea owner's logins
// get none (architecture.md §5: never a row of role approver for the idea
// owner, and the decision of #365 rules out the Operator's); a login the forge
// does not know is recorded input and gives no row; a forge error is the
// step's.
func TestTheApproverRows(t *testing.T) {
	ctx := context.Background()
	bound := map[string]string{
		"bet": "alice", "audit": "alice", "setup": "pharzam", "retrospective": "carol", "activation": "mallory",
	}
	named := map[string][2]string{
		"alice": {"701", "2026-10-11T10:05:00Z"}, "pharzam": {"701", "2026-10-11T10:05:00Z"},
		"carol": {"701", "2026-10-11T10:05:00Z"}, "mallory": {"704", "2026-10-11T10:09:00Z"},
	}
	ids := func(_ context.Context, login string) (int64, error) {
		switch login {
		case "alice":
			return 404, nil
		case "pharzam":
			return 101, nil
		case "carol":
			return 303, nil
		}
		return 0, &forge.Error{Call: "UserID", Status: 404, Message: "Not Found"}
	}
	rows, err := approverRows(ctx, ids, bound, named, answersApprovers)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"404", "approver", "alice", "2026-10-11T10:05:00Z", "701"}
	if len(rows) != 1 || strings.Join(rows[0], " ") != strings.Join(want, " ") {
		t.Errorf("the new rows %v; want %v", rows, want)
	}
	fail := func(context.Context, string) (int64, error) {
		return 0, &forge.Error{Call: "UserID", Status: 502, Message: "Bad Gateway"}
	}
	if _, err := approverRows(ctx, fail, bound, named, answersApprovers); err == nil || !strings.Contains(err.Error(), "502") {
		t.Errorf("a forge error: %v; want the error, not a row", err)
	}
}

// The account ID, not the login's spelling, is what an approver row is of: a
// counted line naming PHARZAM, which the forge resolves to the recorded
// Operator's ID, gives no second row of that account; and the point resolves
// to the Operator's row, never to the idea-owner fallback — the two accounts
// hold different IDs so the fallback cannot mask a wrong binding (decided
// here, task T-8rm6, after the audit of #365).
func TestAnAccountARowAlreadyServesIsNamedByAnotherSpelling(t *testing.T) {
	ctx := context.Background()
	bound := map[string]string{"bet": "dave", "activation": "PHARZAM", "audit": "DAVE"}
	named := map[string][2]string{
		"PHARZAM": {"801", "2026-10-11T10:05:00Z"},
		"dave":    {"802", "2026-10-11T10:06:00Z"},
		"DAVE":    {"803", "2026-10-11T10:07:00Z"},
	}
	ids := func(_ context.Context, login string) (int64, error) {
		switch {
		case strings.EqualFold(login, "pharzam"):
			return 101, nil // the recorded Operator, by another spelling
		case strings.EqualFold(login, "dave"):
			return 707, nil // one unrecorded account, two spellings
		}
		return 0, &forge.Error{Call: "UserID", Status: 404, Message: "Not Found"}
	}
	rows, err := approverRows(ctx, ids, bound, named, answersApprovers)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"707", "approver", "dave", "2026-10-11T10:06:00Z", "802"}
	if len(rows) != 1 || strings.Join(rows[0], " ") != strings.Join(want, " ") {
		t.Errorf("the new rows %v; want %v — PHARZAM resolves to the Operator's recorded ID 101, and DAVE to dave's ID 707, so one row of 707", rows, want)
	}
	got := approverFor("activation", bound, answersApprovers)
	if got == nil || got[0] != "101" || got[1] != "operator" {
		t.Errorf("the approver of activation %v; want the Operator's row 101, the account PHARZAM resolves to — not the idea owner's row 303 of the fallback", got)
	}
	if served := pointsServed(answersApprovers[1], bound, answersApprovers); !slices.Contains(served, "activation") {
		t.Errorf("the points the Operator's row serves %v; want activation among them", served)
	}
}

// The account that may answer a point, and the set of points a row serves (the
// answer #366's caller computes): the row of the login the point's last
// counted line names, whatever its role — the naming comment is enough for an
// account that already has a row — and the idea owner's row for a point that
// names none; a login named but holding no row binds nothing (its line was
// input only, and the point falls back to the idea owner).
func TestTheApproverOfEachPoint(t *testing.T) {
	rows := append(slices.Clone(answersApprovers), []string{"404", "approver", "alice", "2026-10-11T10:05:00Z", "701"})
	bound := map[string]string{"setup": "pharzam", "bet": "alice", "audit": "alice", "activation": "mallory"}
	for _, c := range []struct{ point, wantRole, wantLogin string }{
		{"bet", "approver", "alice"},             // a named login with its new row
		{"audit", "approver", "alice"},           // one login may serve two points
		{"setup", "operator", "pharzam"},         // bound to the Operator's row by the copy
		{"retrospective", "idea-owner", "carol"}, // no line names it: the idea owner's
		{"activation", "idea-owner", "carol"},    // named a login with no row: none
	} {
		got := approverFor(c.point, bound, rows)
		if got == nil || got[1] != c.wantRole || got[2] != c.wantLogin {
			t.Errorf("the approver of %s: %v; want the %s row of %s", c.point, got, c.wantRole, c.wantLogin)
		}
	}
	for _, c := range []struct {
		row  int
		want string
	}{
		{0, "activation retrospective early-retrospective"}, // the idea owner's: unbound, and the rowless naming
		{1, "setup"},     // the Operator's row serves setup
		{2, "bet audit"}, // alice's row serves both it names
	} {
		got := pointsServed(rows[c.row], bound, rows)
		if strings.Join(got, " ") != c.want {
			t.Errorf("the points of row %v: %v; want %q", rows[c.row], got, c.want)
		}
	}
}

// The criteria of the intent answer's answer (docs/spec/loop.md, The approvers
// and the success criteria, Decided here): split on ";", each segment trimmed,
// the empty dropped.
func TestTheSuccessCriteriaOfTheAnswer(t *testing.T) {
	for _, c := range []struct {
		in   string
		want string
	}{
		{"a;b", "a b"}, {"a; b", "a b"}, {" a ;b ;c ", "a b c"},
		{"a;;b", "a b"}, {"a;", "a"}, {"one", "one"}, {"  ", ""}, {"", ""},
	} {
		if got := criteriaOf(c.in); strings.Join(got, " ") != c.want {
			t.Errorf("the criteria of %q: %q; want %q", c.in, got, c.want)
		}
	}
	questions := append(slices.Clone(answersQuestions),
		[]string{"Q-006", "intent", "success-criteria", "idea-owner", "before-delivery", "intent", "", "", "The criteria?", "2026-10-11T10:00:00Z", "layup run", "", "asked"})
	other := [][]string{{"Q-001", "Go", "801", "303", "2026-10-11T10:07:00Z"}}
	answers := append(slices.Clone(other), []string{"Q-006", "fast; cheap", "801", "303", "2026-10-11T10:07:00Z"})
	id, got := criteriaAnswer(questions, answers)
	if id != "Q-006" || strings.Join(got, " ") != "fast cheap" {
		t.Errorf("the criteria answer: %q, %q; want Q-006 and [fast cheap]", id, got)
	}
	id, got = criteriaAnswer(questions, other)
	if id != "Q-006" || got != nil {
		t.Errorf("the criteria unanswered: %q, %v; want Q-006 and none", id, got)
	}
}
