package records

import (
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// The records of a session of milestone M2b (docs/spec/records.md, NFR-001 —
// The records of a session; task T-3py1, #153): the schemas of sessions.tsv,
// harnesses.tsv, routing.tsv, tasks/<task>/events.tsv and
// tasks/<task>/results/<session>.tsv, and the rules that their blocks give in
// words. internal/run writes and reads them (docs/spec/session.md).
//
// Each "— when" and "— for" of these blocks is read as "exactly when" (the
// plan review of #153, condition 3), where the row can see the condition. A
// rule that needs another file (a model of models.tsv; the admission by the
// last probe of a harness) is the writer's or its reader's, not a rule of a
// row.

// SessionsSchema is the form of sessions.tsv: the block sessions.
var SessionsSchema = tsv.Schema{Name: "sessions", Location: "records:sessions.tsv", Columns: []tsv.Column{
	{Name: "session", Type: "id(S-xxxxxxxx)", Key: true},
	{Name: "task", Type: "id(T-xxxx)"},
	{Name: "attempt", Type: "int"},
	{Name: "role", Type: "text"},
	{Name: "harness", Type: "id(<word>)"},
	{Name: "version", Type: "text"},
	{Name: "model", Type: "text"},
	{Name: "base", Type: "sha1"},
	{Name: "records", Type: "sha1"},
	{Name: "prompt_bytes", Type: "int"},
	{Name: "prompt_tokens", Type: "int"},
	{Name: "context", Type: "int"},
	{Name: "cap", Type: "decimal"},
	{Name: "wall", Type: "int"},
	{Name: "vars", Type: "list(text)"},
	{Name: "policy", Type: "list(text)"},
}}

// HarnessesSchema is the form of harnesses.tsv: the block harnesses.
var HarnessesSchema = tsv.Schema{Name: "harnesses", Location: "records:harnesses.tsv", Columns: []tsv.Column{
	{Name: "session", Type: "id(S-xxxxxxxx)", Key: true},
	{Name: "harness", Type: "id(<word>)"},
	{Name: "version", Type: "text"},
	{Name: "model", Type: "text"},
	{Name: "result", Type: "enum(passed|failed)"},
	{Name: "reason", Type: "text"},
	{Name: "files", Type: "list(text)"},
	{Name: "models", Type: "list(text)"},
	{Name: "end", Type: "time"},
}}

// RoutingSchema is the form of routing.tsv: the block routing.
var RoutingSchema = tsv.Schema{Name: "routing", Location: "records:routing.tsv", Columns: []tsv.Column{
	{Name: "role", Type: "text", Key: true},
	{Name: "tier", Type: "enum(reasoning|execution)", Key: true},
	{Name: "position", Type: "int", Key: true},
	{Name: "harness", Type: "id(<word>)"},
	{Name: "model", Type: "text"},
}}

// EventsSchema is the form of tasks/<task>/events.tsv: the block events.
var EventsSchema = tsv.Schema{Name: "events", Location: "records:tasks/<task>/events.tsv", Columns: []tsv.Column{
	{Name: "n", Type: "int", Key: true},
	{Name: "kind", Type: "enum(attempt|session|refused|result|push|bound|closed|rebased)"},
	{Name: "attempt", Type: "int"},
	{Name: "session", Type: "id(S-xxxxxxxx)"},
	{Name: "base", Type: "sha1"},
	{Name: "sha", Type: "sha1"},
	{Name: "detail", Type: "text"},
	{Name: "time", Type: "time"},
}}

// ResultSchema is the form of tasks/<task>/results/<session>.tsv: the block
// result.
var ResultSchema = tsv.Schema{Name: "result", Location: "records:tasks/<task>/results/<session>.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "enum(status|artifact)", Key: true},
	{Name: "n", Type: "int", Key: true},
	{Name: "value", Type: "text"},
	{Name: "sha256", Type: "sha256"},
	{Name: "reason", Type: "text"},
}}

// The columns whose block rule has no clause for the empty value.
var (
	sessionsRequired  = []string{"task", "attempt", "role", "harness", "version", "model", "base", "records", "prompt_bytes", "prompt_tokens", "context", "wall"}
	harnessesRequired = []string{"harness", "model", "result", "end"}
	routingRequired   = []string{"harness", "model"}
	eventsRequired    = []string{"kind", "attempt", "time"}
)

// The forms that the rules of the blocks give in words.
var (
	varForm    = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*=`)                    // NAME=VALUE of vars
	branchForm = regexp.MustCompile(`^task/T-[a-z0-9]{4}/(0|[1-9][0-9]*)$`)        // the branch of push and bound
	ends       = []string{"done", "start", "crash", "no-result", "wall", "output"} // the classes of the end (session.md)
	// The first words of the reason of a failed probe: the block's words and
	// the classes of the end other than done.
	probeReasons = []string{"version", "rules", "context", "prompt", "outside", "not-used", "token", "files", "start", "crash", "no-result", "wall", "output"}
	statuses     = []string{"completed", "blocked", "needs_context", "decision_needed", "failed"}
)

// isPath reports whether v is of the type path of docs/spec/README.md, as
// internal/tsv checks it: no empty, . or .. part, no / at the start or the end.
func isPath(v string) bool {
	return v != "." && fs.ValidPath(v) && !strings.ContainsAny(v, "\t\n\r")
}

// atLeastOne reports whether v is an int of 1 or more.
func atLeastOne(v string) bool { n, err := strconv.Atoi(v); return err == nil && n >= 1 }

// CheckSession checks the rules of the block sessions on one row: the session
// ID, the attempt (1 for a probe), the estimate of the prompt, the context,
// the wall-clock limit, and the forms of vars and policy.
func CheckSession(r []string) error {
	f := fields(SessionsSchema, r)
	if c := missing(f, sessionsRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	bytes, _ := strconv.Atoi(f["prompt_bytes"])
	tokens, _ := strconv.Atoi(f["prompt_tokens"])
	context, _ := strconv.Atoi(f["context"])
	switch {
	case !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	case !atLeastOne(f["attempt"]):
		return &RowError{"attempt", "the attempt is 1 or more"}
	case f["role"] == "probe" && f["attempt"] != "1":
		return &RowError{"attempt", "the attempt of a probe is 1"}
	case tokens != (bytes+3)/4:
		return &RowError{"prompt_tokens", fmt.Sprintf("prompt_tokens is prompt_bytes over four, rounded up: %d", (bytes+3)/4)}
	case context < tokens:
		return &RowError{"context", "the context is never below prompt_tokens"}
	case !atLeastOne(f["wall"]):
		return &RowError{"wall", "the wall-clock limit is 1 minute or more"}
	}
	for _, v := range strings.Fields(f["vars"]) {
		if !varForm.MatchString(v) {
			return &RowError{"vars", "each fixed variable is NAME=VALUE"}
		}
	}
	for _, p := range strings.Fields(f["policy"]) {
		if !strings.HasPrefix(p, "/") {
			return &RowError{"policy", "each policy path is absolute"}
		}
	}
	return nil
}

// ReadSessions reads sessions.tsv by its schema and its rules.
func ReadSessions(data []byte) ([][]string, error) { return read(data, SessionsSchema, CheckSession) }

// CheckHarness checks the rules of the block harnesses on one row: the
// session ID; reason exactly when the probe failed, its first word one of the
// block's; version empty exactly when the version check refused the start.
func CheckHarness(r []string) error {
	f := fields(HarnessesSchema, r)
	if c := missing(f, harnessesRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	word, _, _ := strings.Cut(f["reason"], " ")
	switch {
	case !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	case (f["result"] == "passed") != (f["reason"] == ""):
		return &RowError{"reason", "reason is the empty value exactly when the probe passed"}
	case f["result"] == "failed" && !slices.Contains(probeReasons, word):
		return &RowError{"reason", "a failed probe's reason starts with a word of the block or a class of the end: " + strings.Join(probeReasons, ", ")}
	case (word == "version") != (f["version"] == ""):
		return &RowError{"version", "version is the empty value exactly when the version check refused the start"}
	}
	return nil
}

// ReadHarnesses reads harnesses.tsv by its schema and its rules.
func ReadHarnesses(data []byte) ([][]string, error) {
	return read(data, HarnessesSchema, CheckHarness)
}

// CheckRouting checks the rules of the block routing: no column holds the
// empty value, and the positions of each role and tier are 1 to k, each once,
// in any order of the file (the key refuses a repeat). A broken order is a
// *tsv.Error.
func CheckRouting(rows [][]string) error {
	count, most, line := map[string]int{}, map[string]int{}, map[string]int{}
	var lists []string
	for i, r := range rows {
		f := fields(RoutingSchema, r)
		if c := missing(f, routingRequired); c != "" {
			return &tsv.Error{Line: i + 2, Column: c, Reason: "this column never holds the empty value"}
		}
		k := f["role"] + " " + f["tier"]
		if _, ok := count[k]; !ok {
			lists = append(lists, k)
		}
		p, _ := strconv.Atoi(f["position"])
		count[k]++
		if p > most[k] {
			most[k], line[k] = p, i+2
		}
	}
	for _, k := range lists {
		if most[k] != count[k] {
			return &tsv.Error{Line: line[k], Column: "position", Reason: fmt.Sprintf("the positions of %s run 1 to %d with no gap", k, count[k])}
		}
	}
	return nil
}

// ReadRouting reads routing.tsv by its schema and its rules.
func ReadRouting(data []byte) ([][]string, error) {
	rows, err := tsv.Read(data, RoutingSchema)
	if err != nil {
		return nil, err
	}
	if err := CheckRouting(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckEvent checks the rules of the block events on one row: which columns
// each kind fills, the session ID, the class of a result, and the branch of a
// push or a bind, task/<task>/<attempt> with the row's attempt.
func CheckEvent(r []string) error {
	f := fields(EventsSchema, r)
	if c := missing(f, eventsRequired); c != "" {
		return &RowError{c, "this column never holds the empty value"}
	}
	kind, detail := f["kind"], f["detail"]
	reason, _, _ := strings.Cut(detail, " ") // the reason of a refused, before its value
	noSession := kind == "attempt" || kind == "closed" || kind == "rebased" || (kind == "refused" && reason == "pair")
	switch {
	case !atLeastOne(f["attempt"]):
		return &RowError{"attempt", "the attempt is 1 or more"}
	case f["session"] != "" && !sessionForm.MatchString(f["session"]):
		return &RowError{"session", "a session ID is S- and 8 lowercase hexadecimal characters"}
	case noSession != (f["session"] == ""):
		return &RowError{"session", "session is the empty value exactly for attempt, closed, rebased and a refused whose reason is pair"}
	case (kind == "attempt" || kind == "rebased") != (f["base"] != ""):
		return &RowError{"base", "base is given exactly for attempt and rebased"}
	case (kind == "push" || kind == "bound") != (f["sha"] != ""):
		return &RowError{"sha", "sha is given exactly for push and bound"}
	case (kind == "attempt" || kind == "session") && detail != "":
		return &RowError{"detail", "detail is the empty value for attempt and session"}
	case kind == "result" && !slices.Contains(ends, detail):
		return &RowError{"detail", "the detail of a result is a class of the end: " + strings.Join(ends, ", ")}
	case kind == "refused" && detail == "":
		return &RowError{"detail", "the detail of a refused is its reason"}
	case kind == "push" || kind == "bound":
		if m := branchForm.FindStringSubmatch(detail); m == nil || m[1] != f["attempt"] {
			return &RowError{"detail", "the detail of a push or a bind is the branch task/<task>/" + f["attempt"]}
		}
	}
	return nil
}

// CheckEvents checks the order of events.tsv: n runs 1, 2, ... with no gap,
// in the order of the file. A broken order is a *tsv.Error.
func CheckEvents(rows [][]string) error {
	for i, r := range rows {
		if f := fields(EventsSchema, r); f["n"] != strconv.Itoa(i+1) {
			return &tsv.Error{Line: i + 2, Column: "n", Reason: fmt.Sprintf("the next event is n %d", i+1)}
		}
	}
	return nil
}

// ReadEvents reads events.tsv by its schema, the rules of each row and the
// rule of the order.
func ReadEvents(data []byte) ([][]string, error) {
	rows, err := read(data, EventsSchema, CheckEvent)
	if err != nil {
		return nil, err
	}
	if err := CheckEvents(rows); err != nil {
		return nil, err
	}
	return rows, nil
}

// CheckEventsAppend checks that after holds the rows of before, unchanged and
// first: each row is one event, and no row is changed (§3).
func CheckEventsAppend(before, after [][]string) error {
	if len(after) < len(before) {
		return fmt.Errorf("an event was removed: %d events before, %d after", len(before), len(after))
	}
	for i := range before {
		if !slices.Equal(before[i], after[i]) {
			return fmt.Errorf("event %d was changed; an event is only appended", i+1)
		}
	}
	return nil
}

// CheckResultRow checks the rules of the block result on one row: the words
// of a status, the path of an artifact, and which of sha256 and reason each
// kind fills.
func CheckResultRow(r []string) error {
	f := fields(ResultSchema, r)
	if f["value"] == "" {
		return &RowError{"value", "this column never holds the empty value"}
	}
	status := f["kind"] == "status"
	switch {
	case status && !slices.Contains(statuses, f["value"]):
		return &RowError{"value", "a status is one of " + strings.Join(statuses, ", ")}
	case !status && !isPath(f["value"]):
		return &RowError{"value", "an artifact is a path of the target, of the type path"}
	case status == (f["sha256"] != ""):
		return &RowError{"sha256", "sha256 is given exactly for an artifact"}
	case status != (f["reason"] != ""):
		return &RowError{"reason", "reason is given exactly for a status"}
	case status && f["n"] != "1":
		return &RowError{"n", "a result holds one status row, with n 1"}
	}
	return nil
}

// CheckResult checks the numbers of a result: within each kind, n is 1 to k,
// each once, in any order of the file (the block gives a number, not an
// order; the key refuses a repeat), and one status row. A broken rule is a
// *tsv.Error.
func CheckResult(rows [][]string) error {
	count, most, line := map[string]int{}, map[string]int{}, map[string]int{}
	for i, r := range rows {
		f := fields(ResultSchema, r)
		n, _ := strconv.Atoi(f["n"])
		count[f["kind"]]++
		if n > most[f["kind"]] {
			most[f["kind"]], line[f["kind"]] = n, i+2
		}
	}
	for _, k := range []string{"status", "artifact"} {
		if most[k] != count[k] {
			return &tsv.Error{Line: line[k], Column: "n", Reason: fmt.Sprintf("the numbers of the %s rows run 1 to %d with no gap", k, count[k])}
		}
	}
	if count["status"] != 1 {
		return &tsv.Error{Line: 1, Column: "kind", Reason: "a result holds one status row"}
	}
	return nil
}

// ReadResult reads a result by its schema, the rules of each row and the
// rules of the table.
func ReadResult(data []byte) ([][]string, error) {
	rows, err := read(data, ResultSchema, CheckResultRow)
	if err != nil {
		return nil, err
	}
	if err := CheckResult(rows); err != nil {
		return nil, err
	}
	return rows, nil
}
