package run

import (
	"strconv"
	"strings"
	"testing"
)

// The command `cycle <task> because <reason>` (task T-0xr7, #493; docs/spec/loop.md, One more cycle at the cap), on the pure judge and the views of
// shapeloop_test.go.

const (
	opID, ownerID, otherID = "101", "303", "999"
	ctlIssue               = 2
)

var cycleApprovers = [][]string{
	{opID, "operator", "pharzam", "2026-10-08T12:00:00Z", "start"},
	{ownerID, "idea-owner", "carol", "2026-10-08T12:00:00Z", "start"},
}

// aCopy is a candidate: the first copy of the comment id on issue by the account author (app "" for none).
func aCopy(id int, seen string, issue int, author, app, body string) setInput {
	return setInput{Row: []string{strconv.Itoa(id), seen, strconv.Itoa(issue), author, "login", app, "2026-10-08T13:00:00Z", "2026-10-08T13:00:01Z", "sha", "copies/x"}, Body: body}
}

func operatorSays(id int, body string) setInput {
	return aCopy(id, "1", ctlIssue, opID, "", body)
}

// atTheCap is the view of T-sp01 with two rounds with findings: the records of the cap round are the comment 14 and its verdict is the second argument.
func atTheCap(verdict string, cands ...setInput) *viewBuilder {
	b := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1).ev("attempt", "2", "", shB, "", "")
	b.start("S-d2", roleSpec, "2", shB).end("S-d2", "2", "done", "completed").ev("push", "2", "S-d2", "", shG, "task/T-sp01/2").ev("bound", "2", "S-d2", "", shG, "task/T-sp01/2").
		ev("handoff", "2", "S-d2", "", "", "valid").start("S-v2", roleQA, "2", shG).end("S-v2", "2", "done", "completed").
		ev("handoff", "2", "S-v2", "", "", "valid comment 14-1").copyOf("S-v2", verdict, 2)
	b.v.Control, b.v.Approvers, b.v.Cycles = ctlIssue, cycleApprovers, cands
	return b
}

func TestParseCycleReadsTheFormAndRefusesABrokenOne(t *testing.T) {
	for _, c := range []struct {
		body    string
		task    string
		isCycle bool
		bad     string // a word of the error, "" for none
	}{
		{"cycle T-ab12 because the finding is real", "T-ab12", true, ""},
		{"cycle T-ab12 because the first\nsecond line", "", true, "more than one line"},
		{"cycle T-ab12", "", true, "the form is"},
		{"cycle T-ab12 because —", "", true, "not empty and not"},
		{"cycle T-AB12 because x", "", true, "T-xxxx"},
		{"cycle because x y", "", true, "T-xxxx"},
		{"cycles T-ab12 because x", "", false, ""},
		{"set brief.lines 8 because x", "", false, ""},
	} {
		line, isCycle, err := parseCycle(c.body)
		if isCycle != c.isCycle || line.Task != c.task || (c.bad == "") != (err == nil) || (err != nil && !strings.Contains(err.Error(), c.bad)) {
			t.Errorf("parseCycle(%q) = %+v, %v, %v; want task %q, isCycle %v, error with %q", c.body, line, isCycle, err, c.task, c.isCycle, c.bad)
		}
	}
}

// The first three items of the order of judgement, each reached: form, author, issue; and the order itself: a comment that fails two is refused for the first.
func TestACommandIsJudgedByTheFormTheAuthorAndTheIssueInThatOrder(t *testing.T) {
	good := "cycle T-sp01 because the finding is real"
	for _, c := range []struct {
		name string
		in   setInput
		item string // "" for a command that passes
	}{
		{"the Operator on the control issue", operatorSays(20, good), ""},
		{"a broken form", operatorSays(20, "cycle T-sp01"), "form"},
		{"the idea owner only", aCopy(20, "1", ctlIssue, ownerID, "", good), "author"},
		{"a stranger", aCopy(20, "1", ctlIssue, otherID, "", good), "author"},
		{"the App", aCopy(20, "1", ctlIssue, opID, "layup-agent", good), "author"},
		{"an edit of the Operator's comment", aCopy(20, "2", ctlIssue, opID, "", good), "author"},
		{"a stranger on another issue", aCopy(20, "1", 3, otherID, "", good), "author"},
		{"the Operator on the issue of the task", aCopy(20, "1", 3, opID, "", good), "issue"},
	} {
		_, isCycle, ref := judgeCommand(c.in, ctlIssue, cycleApprovers)
		switch {
		case !isCycle:
			t.Errorf("%s: judged no command", c.name)
		case c.item == "" && ref != nil, c.item != "" && (ref == nil || ref.Item != c.item || ref.Comment != 20 || ref.Reason == ""):
			t.Errorf("%s: the refusal %+v; want the item %q", c.name, ref, c.item)
		}
		if b := atTheCap("material", c.in); (c.item == "") != (grantOf(b.v) == 1) || (c.item != "" && b.at().Pos != posWait) {
			t.Errorf("%s: the grant %d, the position %+v; want a grant exactly when the command passes, and the wait otherwise", c.name, grantOf(b.v), b.at())
		}
	}
	if _, isCycle, ref := judgeCommand(operatorSays(20, "thanks"), ctlIssue, cycleApprovers); isCycle || ref != nil {
		t.Errorf("a comment that is no command: %v, %+v; want neither", isCycle, ref)
	}
}

func TestAGrantOfTheOperatorAtTheCapStartsTheNextAttempt(t *testing.T) {
	b := atTheCap("material")
	got := b.at()
	want := "waiting for the Operator: T-sp01: the review-cycle cap 1 is reached with findings open; the Operator may grant one more cycle with `cycle T-sp01 because <reason>` on the control issue"
	if got.Pos != posWait || got.Msg != want {
		t.Fatalf("the wait with no grant: %+v; want %q", got, want)
	}
	b.v.Cycles = []setInput{operatorSays(20, "cycle T-sp01 because the finding of round 2 is real")}
	if g, refused := judgeGrants(b.v, b.v.Cycles); g != 20 || len(refused) != 0 {
		t.Errorf("judgeGrants = %d, %v; want the grant of the comment 20", g, refused)
	}
	if got := b.at(); got.Pos != posNextAttempt || got.Session != "S-v2" || got.Base != shB || got.Attempt != 2 {
		t.Errorf("the position with the grant: %+v; want the next attempt at the base, after the verdict of S-v2", got)
	}
	// A grant that names another task grants this one nothing.
	b.v.Cycles = []setInput{operatorSays(20, "cycle T-zz99 because x")}
	if g, refused := judgeGrants(b.v, b.v.Cycles); g != 0 || len(refused) != 0 {
		t.Errorf("a command of another task: %d, %v; want no grant and no refusal here", g, refused)
	}
}

// Condition 3 and 4 of the plan review: IDs below, equal to and above the comment of the cap round; and a command before the cap round, at round 1.
func TestAGrantBeforeOrAtTheCapRoundIsRefusedForGoodAndOneAfterItIsTheGrant(t *testing.T) {
	for _, c := range []struct {
		id   int
		item string // "" for the grant
	}{{5, "cap"}, {13, "cap"}, {14, "cap"}, {15, ""}, {6347893426, ""}} {
		b := atTheCap("material", operatorSays(c.id, "cycle T-sp01 because x"))
		g, refused := judgeGrants(b.v, b.v.Cycles)
		switch {
		case c.item == "" && (g != int64(c.id) || len(refused) != 0), c.item != "" && (g != 0 || len(refused) != 1 || refused[0].Item != c.item || refused[0].Comment != int64(c.id)):
			t.Errorf("the comment %d: grant %d, refusals %+v; want the item %q", c.id, g, refused, c.item)
		}
	}
	// At round 1 there is no cap round: the command is refused with the item cap, and the position is the next attempt as before.
	b := upTo("specification", 14).ev("handoff", "1", "S-v1", "", "", "valid comment 13-1").copyOf("S-v1", "material", 1)
	b.v.Control, b.v.Approvers, b.v.Cycles = ctlIssue, cycleApprovers, []setInput{operatorSays(20, "cycle T-sp01 because early")}
	if g, refused := judgeGrants(b.v, b.v.Cycles); g != 0 || len(refused) != 1 || refused[0].Item != "cap" {
		t.Errorf("a command at round 1: %d, %+v; want the item cap", g, refused)
	}
	if got := b.at(); got.Pos != posNextAttempt {
		t.Errorf("round 1 with a command: %+v; want the next attempt, the cap is not reached", got)
	}
	// A command posted before the cap stays refused after the cap round exists and after a later command is the grant.
	later := atTheCap("material", operatorSays(10, "cycle T-sp01 because early"), operatorSays(20, "cycle T-sp01 because now"))
	if g, refused := judgeGrants(later.v, later.v.Cycles); g != 20 || len(refused) != 1 || refused[0].Comment != 10 || refused[0].Item != "cap" {
		t.Errorf("an early command and a later one: %d, %+v; want the grant 20 and the comment 10 refused", g, refused)
	}
}

func TestOneGrantPerTaskAndTheCapAfterItIsTwo(t *testing.T) {
	b := atTheCap("material", operatorSays(21, "cycle T-sp01 because again"), operatorSays(20, "cycle T-sp01 because the first")) // not in the order of the IDs
	g, refused := judgeGrants(b.v, b.v.Cycles)
	if g != 20 || len(refused) != 1 || refused[0].Item != "bound" || refused[0].Comment != 21 {
		t.Fatalf("two commands: %d, %+v; want the grant 20 and the comment 21 refused with the item bound", g, refused)
	}
	// Attempt 3 starts; its developer and its verifier find a third round with findings: the cap is 2, the wait names no command, and the grant stays.
	b.ev("attempt", "3", "", shB, "", "")
	if got := b.at(); got.Pos != posDeveloper || got.Attempt != 3 {
		t.Fatalf("after the new attempt: %+v; want the developer of attempt 3", got)
	}
	if grantOf(b.v) != 1 {
		t.Errorf("grantOf after the attempt started = %d; want 1, an accepted grant stays accepted", grantOf(b.v))
	}
	b.start("S-d3", roleSpec, "3", shB).end("S-d3", "3", "done", "completed").ev("push", "3", "S-d3", "", shH, "task/T-sp01/3").ev("bound", "3", "S-d3", "", shH, "task/T-sp01/3").
		ev("handoff", "3", "S-d3", "", "", "valid").start("S-v3", roleQA, "3", shH).end("S-v3", "3", "done", "completed").
		ev("handoff", "3", "S-v3", "", "", "valid comment 30-1").copyOf("S-v3", "material", 1)
	want := "waiting for the Operator: T-sp01: the review-cycle cap 2 is reached with findings open"
	if got := b.at(); got.Pos != posWait || got.Msg != want {
		t.Errorf("the third round: %+v; want the wait %q", got, want)
	}
	// A command after the third round is the second grant: refused with bound, and the wait stays.
	b.v.Cycles = append(b.v.Cycles, operatorSays(40, "cycle T-sp01 because a third cycle"))
	if _, refused := judgeGrants(b.v, b.v.Cycles); len(refused) != 2 || refused[1].Item != "bound" || refused[1].Comment != 40 {
		t.Errorf("a command after the third round: %+v; want it refused with the item bound", refused)
	}
	if got := b.at(); got.Pos != posWait {
		t.Errorf("the third round with a second command: %+v; want the wait", got)
	}
}

// Plan review, condition 3: the wait of a record `not mergeable, findings recorded` grants no cycle and consumes no grant, whatever the commands; and a
// task that is past the cap with the cap round material and no command is the wait that names the command.
func TestANotMergeableCapRoundGrantsNoCycleAndConsumesNoGrant(t *testing.T) {
	b := atTheCap("not mergeable, findings recorded", operatorSays(20, "cycle T-sp01 because x"), operatorSays(21, "cycle T-sp01 because y"))
	g, refused := judgeGrants(b.v, b.v.Cycles)
	if g != 0 || len(refused) != 2 || refused[0].Item != "cap" || refused[1].Item != "cap" {
		t.Errorf("commands on a not mergeable cap round: %d, %+v; want both refused with the item cap and no grant consumed", g, refused)
	}
	if got := b.at(); got.Pos != posWait || got.Msg != "waiting for the Operator: T-sp01: not mergeable, findings recorded" {
		t.Errorf("the position: %+v; want the wait for the not mergeable record", got)
	}
}
