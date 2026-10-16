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
	who, ok := AnswerAccepted(taskEvents(strings.Split(stdEvents, "\n")...), answered(), answerID)
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
		who, ok := AnswerAccepted(taskEvents(strings.Split(c.events, "\n")...), c.results, c.id)
		if ok != (c.want != "") || who != c.want {
			t.Errorf("%s: %q, %v; want %q", c.name, who, ok, c.want)
		}
	}
}
