package run

import (
	"strconv"
	"strings"
	"testing"
)

// taskEvents gives the rows of an events.tsv from the lines "kind attempt
// session detail", with — for none.
func taskEvents(lines ...string) [][]string {
	var out [][]string
	for i, l := range lines {
		f := strings.Fields(l)
		for len(f) < 4 {
			f = append(f, "—")
		}
		for j := range f {
			if f[j] == "—" {
				f[j] = ""
			}
		}
		out = append(out, []string{strconv.Itoa(i + 1), f[0], f[1], f[2], "", "", f[3], "2026-10-09T12:00:00Z"})
	}
	return out
}

const (
	askSess    = "S-aaaaaaaa" // attempt 1 ends with a question
	ownerSess  = "S-bbbbbbbb" // the owner role's answer
	nextSess   = "S-cccccccc" // the next attempt
	answerID   = "S-bbbbbbbb.1"
	askedRow   = "question\t1\tMay it?\t—\tinterface-contract no docs/spec/forge.md"
	cited      = "completed\t—\tdone with " + answerID
	stdEvents  = "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done\nsession 1 " + ownerSess + "\nresult 1 " + ownerSess + " done\nattempt 2\nsession 2 " + nextSess + "\nresult 2 " + nextSess + " done"
	needsHuman = "needs_context\t—\tasks"
)

// stdRoles are the roles of the sessions of stdEvents: the question is of the
// interface-contract kind, whose owner is the Software Architect.
func stdRoles() Roles {
	return Roles{askSess: "software-developer", ownerSess: "software-architect", nextSess: "software-developer"}
}

func resultOfRows(status string, extra ...string) [][]string {
	var out [][]string
	f := strings.Split(status, "\t")
	out = append(out, []string{"status", "1", f[0], "", f[2]})
	for _, e := range extra {
		r := strings.Split(e, "\t")
		for i := range r {
			if r[i] == "—" {
				r[i] = ""
			}
		}
		out = append(out, r)
	}
	return out
}

func answered() map[string][][]string {
	return map[string][][]string{
		askSess:   resultOfRows(needsHuman, askedRow),
		ownerSess: resultOfRows("completed\t—\tanswered", "answer\t1\tYes\t—\tdocs/spec/forge.md"),
		nextSess:  resultOfRows(cited),
	}
}

func TestAnAnswerIsAcceptedByTheNextAttemptThatCitesIt(t *testing.T) {
	who, ok := AnswerAccepted(taskEvents(strings.Split(stdEvents, "\n")...), answered(), stdRoles(), Owners(DefaultOwners()), answerID, nil)
	if !ok || who != nextSess {
		t.Errorf("AnswerAccepted = %q, %v; want accepted by %s, the session of attempt 2", who, ok, nextSess)
	}
}

func TestAnAnswerIsNotAcceptedUnlessEachConditionHolds(t *testing.T) {
	with := func(session, rows string) map[string][][]string {
		m := answered()
		if rows == "" {
			delete(m, session)
		} else {
			m[session] = resultOfRows(strings.Split(rows, "|")[0], strings.Split(rows, "|")[1:]...)
		}
		return m
	}
	for _, c := range []struct {
		name    string
		events  string
		results map[string][][]string
		id      string
		want    string // the accepting session, "" for none
	}{
		{"the next attempt failed", stdEvents, with(nextSess, "failed\t—\tstopped "+answerID), answerID, ""},
		{"the next attempt blocked", stdEvents, with(nextSess, "blocked\t—\twaits "+answerID), answerID, ""},
		{"the next attempt completed and does not cite", stdEvents, with(nextSess, "completed\t—\tdone"), answerID, ""},
		{"the cite is part of a longer word", stdEvents, with(nextSess, "completed\t—\tdone x"+answerID+"x"), answerID, ""},
		{"the cite is of another answer", stdEvents, with(nextSess, "completed\t—\tdone S-bbbbbbbb.2"), answerID, ""},
		{"a question of the next attempt cites the answer", stdEvents, with(nextSess, cited+"|question\t1\tAgain?\t—\tinterface-contract no "+answerID), answerID, ""},
		{"a question of the next attempt cites another record", stdEvents, with(nextSess, cited+"|question\t1\tMore?\t—\tinterface-contract no docs/spec/forge.md"), answerID, nextSess},
		{"the next attempt has no result record", stdEvents, with(nextSess, ""), answerID, ""},
		{"the result of the next attempt did not end done", strings.Replace(stdEvents, "result 2 "+nextSess+" done", "result 2 "+nextSess+" crash", 1), answered(), answerID, ""},
		{"the attempt after the next one cites it", strings.Replace(strings.Replace(stdEvents, "attempt 2", "attempt 3", 1), "session 2", "session 3", 1), answered(), answerID, ""},
		{"the attempt started before the answer was given", "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done\nattempt 2\nsession 1 " + ownerSess + "\nresult 1 " + ownerSess + " done\nsession 2 " + nextSess + "\nresult 2 " + nextSess + " done", answered(), answerID, ""},
		{"the next attempt was closed", stdEvents + "\nclosed 2", answered(), answerID, ""},
		{"the next attempt was rebased", stdEvents + "\nrebased 2", answered(), answerID, ""},
		{"no attempt asked", "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done\nsession 1 " + ownerSess + "\nresult 1 " + ownerSess + " done\nattempt 2\nsession 2 " + nextSess + "\nresult 2 " + nextSess + " done",
			with(askSess, "completed\t—\tno question"), answerID, ""},
		{"the answer's session did not end done", strings.Replace(stdEvents, "result 1 "+ownerSess+" done", "result 1 "+ownerSess+" crash", 1), answered(), answerID, ""},
		{"the answer row is not in the result", stdEvents, answered(), "S-bbbbbbbb.2", ""},
		{"the answer ID is no ID", stdEvents, answered(), "bbbbbbbb", ""},
		{"the answer ID is empty", stdEvents, answered(), "", ""},
		{"a result with no answer row", stdEvents, with(ownerSess, "completed\t—\tanswered"), answerID, ""},
		{"a second session of the next attempt that cites it is the accepting one", stdEvents + "\nsession 2 S-dddddddd\nresult 2 S-dddddddd done",
			func() map[string][][]string {
				m := with(nextSess, "failed\t—\tstopped")
				m["S-dddddddd"] = resultOfRows(cited)
				return m
			}(), answerID, "S-dddddddd"},
	} {
		who, ok := AnswerAccepted(taskEvents(strings.Split(c.events, "\n")...), c.results, stdRoles(), Owners(DefaultOwners()), c.id, nil)
		if ok != (c.want != "") || who != c.want {
			t.Errorf("%s: %q, %v; want %q", c.name, who, ok, c.want)
		}
	}
}

// The ledger of stdEvents with the roles and the question rows changed: the
// session of the answer is accepted or not by its role against the owner of
// every question of the attempt it answers (loop.md, A question during the work).
func TestAnAnswerOfASessionOutsideTheOwnerRoleIsNotAccepted(t *testing.T) {
	const second = "S-a2a2a2a2" // a second session of attempt 1 that ends needs_context
	two := "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done\nsession 1 " + second + "\nresult 1 " + second + " done\nsession 1 " + ownerSess + "\nresult 1 " + ownerSess + " done\nattempt 2\nsession 2 " + nextSess + "\nresult 2 " + nextSess + " done"
	other := "question\t1\tWhere?\t—\tenvironment-or-infrastructure no docs/spec/forge.md"
	same := "question\t1\tHow?\t—\tinterface-contract no docs/spec/forge.md"
	withRows := func(session string, rows ...string) map[string][][]string {
		m := answered()
		m[session] = resultOfRows(needsHuman, rows...)
		return m
	}
	roles := func(kv ...string) Roles {
		r := stdRoles()
		for i := 0; i < len(kv); i += 2 {
			r[kv[i]] = kv[i+1]
		}
		return r
	}
	for _, c := range []struct {
		name    string
		events  string
		results map[string][][]string
		roles   Roles
		owners  Owners
		want    string // the accepting session, "" for none
	}{
		{"the owner role", stdEvents, answered(), stdRoles(), nil, nextSess},
		{"a developer wrote the answer", stdEvents, answered(), roles(ownerSess, "software-developer"), nil, ""},
		{"a QA engineer wrote the answer", stdEvents, answered(), roles(ownerSess, "qa-engineer"), nil, ""},
		{"the owner of another kind wrote the answer", stdEvents, answered(), roles(ownerSess, "software-engineer"), nil, ""},
		{"the name of the owner, not its role ID", stdEvents, answered(), roles(ownerSess, "Software Architect"), nil, ""},
		{"the role of the answer is empty", stdEvents, answered(), roles(ownerSess, ""), nil, ""},
		{"the session of the answer has no start row", stdEvents, answered(), Roles{askSess: "software-developer", nextSess: "software-developer"}, nil, ""},
		{"the accepting session is the owner role, not the answering one", stdEvents, answered(), roles(ownerSess, "software-developer", nextSess, "software-architect"), nil, ""},
		{"the question needs a human", stdEvents, withRows(askSess, "question\t1\tMay it?\t—\tinterface-contract yes docs/spec/forge.md"), stdRoles(), nil, ""},
		{"the question needs a human and a session of the operator role answers", stdEvents, withRows(askSess, "question\t1\tMay it?\t—\tinterface-contract yes docs/spec/forge.md"), roles(ownerSess, "operator"), nil, ""},
		{"the label is of no owner", stdEvents, withRows(askSess, "question\t1\tMay it?\t—\tother-kind no docs/spec/forge.md"), stdRoles(), nil, ""},
		{"the attempt has no question row", stdEvents, withRows(askSess), stdRoles(), nil, ""},
		{"a map of the operator as the owner", stdEvents, answered(), roles(ownerSess, "operator"), Owners{"interface-contract": "operator"}, ""},
		{"one session asks two owners, the first answers", stdEvents, withRows(askSess, askedRow, other), stdRoles(), nil, ""},
		{"one session asks two owners, the second answers", stdEvents, withRows(askSess, askedRow, other), roles(ownerSess, "software-engineer"), nil, ""},
		{"two sessions ask two owners, the first answers", two, withRows(second, other), stdRoles(), nil, ""},
		{"two sessions ask two owners, the second answers", two, withRows(second, other), roles(ownerSess, "software-engineer"), nil, ""},
		{"one session asks two questions of one owner", stdEvents, withRows(askSess, askedRow, same), stdRoles(), nil, nextSess},
		{"two sessions ask one owner", two, withRows(second, same), stdRoles(), nil, nextSess},
		{"a replaced owner map names the role of the answer", stdEvents, answered(), roles(ownerSess, "product-owner"), Owners{"interface-contract": "Product Owner"}, nextSess},
		{"a replaced owner map names another role", stdEvents, answered(), stdRoles(), Owners{"interface-contract": "Product Owner"}, ""},
	} {
		owners := c.owners
		if owners == nil {
			owners = DefaultOwners()
		}
		who, ok := AnswerAccepted(taskEvents(strings.Split(c.events, "\n")...), c.results, c.roles, owners, answerID, nil)
		if ok != (c.want != "") || who != c.want {
			t.Errorf("%s: %q, %v; want %q", c.name, who, ok, c.want)
		}
	}
}

// The demo at the level of the helper: two ledgers that differ in one value, the
// role of the session that wrote the answer, and in nothing else.
func TestTwoLedgersThatDifferOnlyInTheRoleOfTheAnsweringSessionAreToldApart(t *testing.T) {
	ev := taskEvents(strings.Split(stdEvents, "\n")...)
	for role, want := range map[string]string{"software-architect": nextSess, "software-developer": ""} {
		roles := stdRoles()
		roles[ownerSess] = role
		if who, ok := AnswerAccepted(ev, answered(), roles, DefaultOwners(), answerID, nil); who != want || ok != (want != "") {
			t.Errorf("an answer of %s: %q, %v; want %q", role, who, ok, want)
		}
	}
}

// The roles come from the start rows of sessions.tsv by session and role, and
// the owner role from the owner map given, not from one the helper holds.
func TestTheOwnerRoleIsReadFromTheQuestionsAndTheOwnerMapGiven(t *testing.T) {
	rows := [][]string{
		{askSess, "T-ab12", "1", "software-developer", "claude"},
		{ownerSess, "T-ab12", "1", "software-architect", "codex"},
	}
	if got := sessionRoles(rows); len(got) != 2 || got[askSess] != "software-developer" || got[ownerSess] != "software-architect" {
		t.Errorf("sessionRoles = %v; want the role of each session by its ID", got)
	}
	ev := taskEvents(strings.Split(stdEvents, "\n")...)
	if who, ok := answerOwner(ev, answered(), "1", Owners(DefaultOwners())); !ok || who != "software-architect" {
		t.Errorf("answerOwner = %q, %v; want software-architect", who, ok)
	}
	if who, ok := answerOwner(ev, answered(), "1", Owners{"interface-contract": "Systems Architect"}); !ok || who != "systems-architect" {
		t.Errorf("answerOwner with a replaced map = %q, %v; want systems-architect", who, ok)
	}
	if who, ok := answerOwner(ev, answered(), "2", Owners(DefaultOwners())); ok {
		t.Errorf("answerOwner of an attempt that asked nothing = %q, true", who)
	}
}

// loop.md, An uncertain end gets one new attempt (task T-qjjd, #475): "the next attempt" of an answer is the first one after the asking attempt that
// a retry did not replace; an attempt started for another reason is the next one, and is the only one.
func TestAnAnswerIsAcceptedByTheAttemptThatReplacesTwoCrashedOnes(t *testing.T) {
	const crashes = "attempt 1\nsession 1 " + askSess + "\nresult 1 " + askSess + " done\nsession 1 " + ownerSess + "\nresult 1 " + ownerSess + " done\n" +
		"attempt 2\nsession 2 S-dddddddd\nresult 2 S-dddddddd crash\nattempt 3\nsession 3 S-eeeeeeee\nresult 3 S-eeeeeeee crash\n"
	events := func(tail string) [][]string { return taskEvents(strings.Split(crashes+tail, "\n")...) }
	both, one := Replaced{"2": "S-dddddddd", "3": "S-eeeeeeee"}, Replaced{"2": "S-dddddddd"}
	atFour := "attempt 4\nsession 4 " + nextSess + "\nresult 4 " + nextSess + " done"
	if who, ok := AnswerAccepted(events(atFour), answered(), stdRoles(), Owners(DefaultOwners()), answerID, both); !ok || who != nextSess {
		t.Errorf("after two replaced attempts: %q, %v; want accepted by %s of attempt 4", who, ok, nextSess)
	}
	if n := OpenQuestions(events(atFour), answered(), stdRoles(), Owners(DefaultOwners()), both); n != 0 {
		t.Errorf("OpenQuestions after two replaced attempts = %d; want 0, the answer accepted", n)
	}
	// Guards: with one of the two replaced, attempt 3 is the next one and it has no completed session; an attempt that no retry replaced is the next.
	if who, ok := AnswerAccepted(events(atFour), answered(), stdRoles(), Owners(DefaultOwners()), answerID, one); ok {
		t.Errorf("with attempt 3 not replaced: accepted by %q; want none, attempt 3 is the next and it did not end done", who)
	}
	if who, ok := AnswerAccepted(events(atFour), answered(), stdRoles(), Owners(DefaultOwners()), answerID, nil); ok {
		t.Errorf("with no attempt replaced: accepted by %q; want none, attempt 2 is the next", who)
	}
}
