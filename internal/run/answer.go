package run

// This file holds the acceptance of an owner role's answer by the next attempt
// of the asking role (docs/spec/loop.md, The checks of a handoff, "A question
// during the work"; docs/spec/records.md, the block result; task T-q5k4, #360,
// row 66d).

import (
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// answerForm is the ID of an answer: the session of the owner role and the n of
// its answer row.
var answerForm = regexp.MustCompile(`^(S-[0-9a-f]{8})\.([1-9][0-9]*)$`)

// citesID reports whether id is a word of the reason of a row.
func citesID(row []string, id string) bool { return slices.Contains(strings.Fields(row[4]), id) }

// resultEnd gives the index of the event result of the session in events, and
// whether its class is done.
func resultEnd(events [][]string, session string) (int, bool) {
	for i, e := range events {
		if e[1] == "result" && e[3] == session {
			return i, e[6] == "done"
		}
	}
	return -1, false
}

// statusOf gives the value of the status row of a result.
func statusOf(result [][]string) string {
	for _, r := range result {
		if r[0] == "status" {
			return r[2]
		}
	}
	return ""
}

// AnswerAccepted says whether the answer with the ID answerID (<session>.<n>)
// is accepted, and by which session, from the events of the task (events.tsv)
// and the typed result of each session (by its ID). The answer is a row of the
// kind answer of a session of the owner role that ended done; it answers the
// attempt that last ended needs_context before that session started; it is
// accepted by the next attempt (the number after it, started after the answer
// and neither closed nor rebased) when one of its sessions ended done with the
// status completed, its status row cites the ID, and none of its question rows
// does. The accepting actor is that session. Nothing is recorded: M3c adds the
// columns (records.md).
func AnswerAccepted(events [][]string, results map[string][][]string, answerID string) (string, bool) {
	m := answerForm.FindStringSubmatch(answerID)
	if m == nil {
		return "", false
	}
	owner := m[1]
	if !slices.ContainsFunc(results[owner], func(r []string) bool { return r[0] == "answer" && r[1] == m[2] }) {
		return "", false
	}
	end, done := resultEnd(events, owner)
	if !done {
		return "", false
	}
	start := slices.IndexFunc(events, func(e []string) bool { return e[1] == "session" && e[3] == owner })
	asking := 0
	for _, e := range events[:max(start, 0)] {
		if e[1] == "result" && e[6] == "done" && statusOf(results[e[3]]) == "needs_context" {
			asking, _ = strconv.Atoi(e[2])
		}
	}
	if asking == 0 {
		return "", false
	}
	next := strconv.Itoa(asking + 1)
	began := slices.IndexFunc(events[end+1:], func(e []string) bool { return e[1] == "attempt" && e[2] == next })
	if began < 0 {
		return "", false
	}
	after := events[end+1+began:]
	if slices.ContainsFunc(after, func(e []string) bool { return (e[1] == "closed" || e[1] == "rebased") && e[2] == next }) {
		return "", false
	}
	for _, e := range after {
		if e[1] != "session" || e[2] != next {
			continue
		}
		result := results[e[3]]
		if _, done := resultEnd(events, e[3]); !done || statusOf(result) != "completed" {
			continue
		}
		cites, asks := false, false
		for _, r := range result {
			switch r[0] {
			case "status":
				cites = citesID(r, answerID)
			case "question":
				asks = asks || citesID(r, answerID)
			}
		}
		if cites && !asks {
			return e[3], true
		}
	}
	return "", false
}
