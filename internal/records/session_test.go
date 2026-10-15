package records

import (
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/tsv"
)

// The records of a session (task T-3py1, #153): one valid file per record,
// one case per rule that its block gives in words, and one case per column
// whose rule has no clause for the empty value. Each refusal must name the
// column of the rule, so a case that fails for another reason fails the test.

// refusedAt reads data with read and reports a test error unless the read
// fails on a *tsv.Error that names column.
func refusedAt(t *testing.T, name string, read func([]byte) ([][]string, error), data []byte, column string) {
	t.Helper()
	_, err := read(data)
	var te *tsv.Error
	switch {
	case err == nil:
		t.Errorf("%s: read, want an error in column %s", name, column)
	case !errors.As(err, &te):
		t.Errorf("%s: %v, want a *tsv.Error in column %s", name, err, column)
	case te.Column != column:
		t.Errorf("%s: %v, want an error in column %s", name, err, column)
	}
}

// field gives row with its field of column (by the header) set to value.
func field(header, row, column, value string) string {
	names := strings.Split(header, "\t")
	f := strings.Split(row, "\t")
	for i, n := range names {
		if n == column {
			f[i] = value
		}
	}
	return strings.Join(f, "\t")
}

const sessionsHeader = "session\ttask\tattempt\trole\tharness\tversion\tmodel\tbase\trecords\tprompt_bytes\tprompt_tokens\tcontext\tcap\twall\tvars\tpolicy"

const sessionRow = "S-1a2b3c4d\tT-ab12\t1\tdeveloper\tclaude\t2.1.295\tclaude-opus-5-5\t" + sha1A + "\t" + sha1A + "\t1001\t251\t200000\t10.0\t60\tANTHROPIC_DEFAULT_HAIKU_MODEL=claude-opus-5-5\t/etc/claude-code/managed-settings.json"

func readSessions(data []byte) ([][]string, error) { return ReadSessions(data) }

func TestAValidSessionsIsRead(t *testing.T) {
	probe := field(sessionsHeader, field(sessionsHeader, field(sessionsHeader, field(sessionsHeader, sessionRow, "session", "S-00000000"), "role", "probe"), "cap", "—"), "vars", "—")
	probe = field(sessionsHeader, probe, "policy", "—")
	if _, err := ReadSessions(table(sessionsHeader, sessionRow, probe)); err != nil {
		t.Fatal(err)
	}
}

func TestSessionsRefusesEachBrokenRule(t *testing.T) {
	if _, err := ReadSessions(table(strings.Replace(sessionsHeader, "wall", "minutes", 1), sessionRow)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"task", "attempt", "role", "harness", "version", "model", "base", "records", "prompt_bytes", "prompt_tokens", "context", "wall"} {
		refusedAt(t, "the empty value in "+c, readSessions, table(sessionsHeader, field(sessionsHeader, sessionRow, c, "—")), c)
	}
	for _, c := range []struct{ name, column, value string }{
		{"a session ID that is not S- and 8 hexadecimal characters", "session", "S-1a2b3c4z"},
		{"attempt 0", "attempt", "0"},
		{"wall 0", "wall", "0"},
		{"prompt_tokens that is not prompt_bytes over four, rounded up", "prompt_tokens", "250"},
		{"a context below prompt_tokens", "context", "250"},
		{"a fixed variable that is not NAME=VALUE", "vars", "ANTHROPIC_DEFAULT_HAIKU_MODEL"},
		{"a policy path that is not absolute", "policy", "etc/claude-code/managed-settings.json"},
	} {
		refusedAt(t, c.name, readSessions, table(sessionsHeader, field(sessionsHeader, sessionRow, c.column, c.value)), c.column)
	}
	probe := field(sessionsHeader, field(sessionsHeader, sessionRow, "role", "probe"), "attempt", "2")
	refusedAt(t, "a probe at attempt 2", readSessions, table(sessionsHeader, probe), "attempt")
}

const harnessesHeader = "session\tharness\tversion\tmodel\tresult\treason\tfiles\tmodels\tend"

const passedRow = "S-1a2b3c4d\tclaude\t2.1.295\tclaude-opus-5-5\tpassed\t—\tAGENTS.md CLAUDE.md\tclaude-opus-5-5\t2026-10-09T12:00:00Z"

func readHarnesses(data []byte) ([][]string, error) { return ReadHarnesses(data) }

func TestAValidHarnessesIsRead(t *testing.T) {
	refused := "S-0000000a\tdevin\t—\tswe-2-high\tfailed\tversion\t—\t—\t2026-10-09T12:01:00Z"
	rule := "S-0000000b\tcodex\t0.9.1\tgpt-6\tfailed\trules /home/layup/AGENTS.md\t—\t—\t2026-10-09T12:02:00Z"
	ended := "S-0000000c\tclaude\t2.1.296\tclaude-opus-5-5\tfailed\twall\t—\t—\t2026-10-09T12:03:00Z"
	if _, err := ReadHarnesses(table(harnessesHeader, passedRow, refused, rule, ended)); err != nil {
		t.Fatal(err)
	}
}

func TestHarnessesRefusesEachBrokenRule(t *testing.T) {
	if _, err := ReadHarnesses(table(strings.Replace(harnessesHeader, "end", "ended", 1), passedRow)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"harness", "model", "result", "end"} {
		refusedAt(t, "the empty value in "+c, readHarnesses, table(harnessesHeader, field(harnessesHeader, passedRow, c, "—")), c)
	}
	failed := field(harnessesHeader, passedRow, "result", "failed")
	for _, c := range []struct{ name, row, column string }{
		{"a session ID that is not S- and 8 hexadecimal characters", field(harnessesHeader, passedRow, "session", "S-1a2b3c4z"), "session"},
		{"a passed probe with a reason", field(harnessesHeader, passedRow, "reason", "token"), "reason"},
		{"a failed probe with no reason", failed, "reason"},
		{"a failed probe whose reason is not a word of the block", field(harnessesHeader, failed, "reason", "slow"), "reason"},
		{"a failed probe whose reason is done", field(harnessesHeader, failed, "reason", "done"), "reason"},
		{"a refused version with a version", field(harnessesHeader, failed, "reason", "version"), "version"},
		{"no version with another reason", field(harnessesHeader, field(harnessesHeader, failed, "reason", "token"), "version", "—"), "version"},
		{"a passed probe with no version", field(harnessesHeader, passedRow, "version", "—"), "version"},
	} {
		refusedAt(t, c.name, readHarnesses, table(harnessesHeader, c.row), c.column)
	}
}

const routingHeader = "role\ttier\tposition\tharness\tmodel"

func readRouting(data []byte) ([][]string, error) { return ReadRouting(data) }

func TestAValidRoutingIsRead(t *testing.T) {
	// The positions of a role and tier run 1 to k in any order of the file.
	rows := []string{
		"developer\texecution\t2\tdevin\tswe-2-high",
		"developer\texecution\t1\tclaude\tclaude-opus-5-5",
		"developer\treasoning\t1\tclaude\tclaude-fable-5-1",
	}
	if _, err := ReadRouting(table(routingHeader, rows...)); err != nil {
		t.Fatal(err)
	}
}

func TestRoutingRefusesEachBrokenRule(t *testing.T) {
	row := "developer\texecution\t1\tclaude\tclaude-opus-5-5"
	if _, err := ReadRouting(table(strings.Replace(routingHeader, "position", "place", 1), row)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"harness", "model"} {
		refusedAt(t, "the empty value in "+c, readRouting, table(routingHeader, field(routingHeader, row, c, "—")), c)
	}
	refusedAt(t, "a gap in the positions of a role and tier", readRouting,
		table(routingHeader, row, "developer\texecution\t3\tdevin\tswe-2-high"), "position")
	refusedAt(t, "a list that starts at 2", readRouting,
		table(routingHeader, "developer\texecution\t2\tdevin\tswe-2-high"), "position")
}

const eventsHeader = "n\tkind\tattempt\tsession\tbase\tsha\tdetail\ttime"

// eventRows is a valid events.tsv of task T-ab12: an attempt, its session,
// the result, the push and the bind, a second attempt refused before a
// session (pair), the task loop's closed and rebased, the handoffs of the task
// loop (with a session, with none for a rendered task, valid, invalid, with the
// copy of a comment) and a merged.
func eventRows() []string {
	return []string{
		"1\tattempt\t1\t—\t" + sha1A + "\t—\t—\t2026-10-09T12:00:00Z",
		"2\tsession\t1\tS-1a2b3c4d\t—\t—\t—\t2026-10-09T12:00:01Z",
		"3\tresult\t1\tS-1a2b3c4d\t—\t—\tdone\t2026-10-09T12:10:00Z",
		"4\tpush\t1\tS-1a2b3c4d\t—\t" + sha1A + "\ttask/T-ab12/1\t2026-10-09T12:10:01Z",
		"5\tbound\t1\tS-1a2b3c4d\t—\t" + sha1A + "\ttask/T-ab12/1\t2026-10-09T12:10:02Z",
		"6\tattempt\t2\t—\t" + sha1A + "\t—\t—\t2026-10-09T12:11:00Z",
		"7\trefused\t2\t—\t—\t—\tpair\t2026-10-09T12:11:01Z",
		"8\trefused\t2\tS-0000000a\t—\t—\tcontext 300000 200000\t2026-10-09T12:11:02Z",
		"9\tclosed\t2\t—\t—\t—\t—\t2026-10-09T12:12:00Z",
		"10\trebased\t3\t—\t" + sha1A + "\t—\tmain moved\t2026-10-09T12:13:00Z",
		"11\thandoff\t1\tS-1a2b3c4d\t—\t—\tvalid comment 6191163525-1\t2026-10-09T12:14:00Z",
		"12\thandoff\t1\t—\t—\t—\tinvalid spec-sections\t2026-10-09T12:15:00Z",
		"13\tmerged\t1\t—\t—\t" + sha1A + "\t334\t2026-10-09T12:16:00Z",
		"14\thandoff\t1\tS-1a2b3c4d\t—\t—\tvalid\t2026-10-09T12:17:00Z",
	}
}

func readEvents(data []byte) ([][]string, error) { return ReadEvents(data) }

// event gives eventRows with row n (from 1) set by column and value.
func event(n int, column, value string) []string {
	rows := eventRows()
	rows[n-1] = field(eventsHeader, rows[n-1], column, value)
	return rows
}

func TestAValidEventsIsRead(t *testing.T) {
	if _, err := ReadEvents(table(eventsHeader, eventRows()...)); err != nil {
		t.Fatal(err)
	}
}

func TestEventsRefusesEachBrokenRule(t *testing.T) {
	if _, err := ReadEvents(table(strings.Replace(eventsHeader, "detail", "details", 1), eventRows()...)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	for _, c := range []string{"kind", "attempt", "time"} {
		refusedAt(t, "the empty value in "+c, readEvents, table(eventsHeader, event(2, c, "—")...), c)
	}
	gap := eventRows()
	gap[1] = field(eventsHeader, gap[1], "n", "3")
	gap[2] = field(eventsHeader, gap[2], "n", "2")
	for _, c := range []struct {
		name   string
		rows   []string
		column string
	}{
		{"attempt 0", event(1, "attempt", "0"), "attempt"},
		{"the events out of their order", gap, "n"},
		{"a session ID that is not S- and 8 hexadecimal characters", event(2, "session", "S-1a2b3c4g"), "session"},
		{"an attempt with a session", event(1, "session", "S-1a2b3c4d"), "session"},
		{"a session event with no session", event(2, "session", "—"), "session"},
		{"a refused with no session whose reason is not pair", event(7, "detail", "context 3 2"), "session"},
		{"a refused pair with a session", event(7, "session", "S-1a2b3c4d"), "session"},
		{"a closed with a session", event(9, "session", "S-1a2b3c4d"), "session"},
		{"an attempt with no base", event(1, "base", "—"), "base"},
		{"a result with a base", event(3, "base", sha1A), "base"},
		{"a push with no sha", event(4, "sha", "—"), "sha"},
		{"a session event with a sha", event(2, "sha", sha1A), "sha"},
		{"an attempt with a detail", event(1, "detail", "first"), "detail"},
		{"a session event with a detail", event(2, "detail", "started"), "detail"},
		{"a result whose detail is not a class of the end", event(3, "detail", "fine"), "detail"},
		{"a refused with no detail", event(8, "detail", "—"), "detail"},
		{"a push to a branch of another form", event(4, "detail", "main"), "detail"},
		{"a bound to the branch of another attempt", event(5, "detail", "task/T-ab12/2"), "detail"},
		{"a push to a task that is not an ID", event(4, "detail", "task/feature/1"), "detail"},
		{"a handoff with a base", event(11, "base", sha1A), "base"},
		{"a handoff with a sha", event(11, "sha", sha1A), "sha"},
		{"a handoff with a session that is not an ID", event(11, "session", "S-1a2b3c"), "session"},
		{"a handoff with no detail", event(11, "detail", "—"), "detail"},
		{"a handoff whose detail is neither valid nor invalid", event(11, "detail", "ok"), "detail"},
		{"an invalid handoff with no item", event(12, "detail", "invalid"), "detail"},
		{"an invalid handoff whose item is not a word of the table", event(12, "detail", "invalid Spec_Sections"), "detail"},
		{"a valid handoff with an item", event(14, "detail", "valid spec-sections"), "detail"},
		{"an invalid handoff with a well-formed word that is not an item", event(12, "detail", "invalid nonsense"), "detail"},
		{"an invalid handoff with an item that has a double hyphen", event(12, "detail", "invalid spec--sections"), "detail"},
		{"a handoff naming a comment with ID 0", event(11, "detail", "valid comment 0-1"), "detail"},
		{"a handoff naming a comment with no copy number", event(11, "detail", "valid comment 6191163525"), "detail"},
		{"a handoff naming a comment with copy 0", event(11, "detail", "valid comment 6191163525-0"), "detail"},
		{"a handoff naming a comment and more", event(11, "detail", "valid comment 6191163525-1 again"), "detail"},
		{"a merged with a session", event(13, "session", "S-1a2b3c4d"), "session"},
		{"a merged with a base", event(13, "base", sha1A), "base"},
		{"a merged with no sha", event(13, "sha", "—"), "sha"},
		{"a merged with no pull request number", event(13, "detail", "—"), "detail"},
		{"a merged whose detail is not a number", event(13, "detail", "the PR"), "detail"},
		{"a merged whose pull request is 0", event(13, "detail", "0"), "detail"},
		{"a merged whose pull request has a leading zero", event(13, "detail", "0334"), "detail"},
	} {
		refusedAt(t, c.name, readEvents, table(eventsHeader, c.rows...), c.column)
	}
}

func TestEventsAreOnlyAppended(t *testing.T) {
	before := eventRows()[:5]
	if err := CheckEventsAppend(rows(before), rows(eventRows())); err != nil {
		t.Fatalf("an appended event: %v", err)
	}
	changed := eventRows()
	changed[2] = field(eventsHeader, changed[2], "detail", "crash")
	if err := CheckEventsAppend(rows(before), rows(changed)); err == nil {
		t.Error("a changed event: want an error")
	}
	if err := CheckEventsAppend(rows(eventRows()), rows(before)); err == nil {
		t.Error("a removed event: want an error")
	}
}

// rows splits each line into its fields, as tsv.Read gives them ("" for —).
func rows(lines []string) [][]string {
	var out [][]string
	for _, l := range lines {
		f := strings.Split(l, "\t")
		for i := range f {
			if f[i] == "—" {
				f[i] = ""
			}
		}
		out = append(out, f)
	}
	return out
}

const resultHeader = "kind\tn\tvalue\tsha256\treason"

func resultRows() []string {
	return []string{
		"status\t1\tcompleted\t—\tthe test passes",
		"artifact\t1\tinternal/records/session.go\t" + sha256A + "\t—",
		"artifact\t2\tinternal/records/session_test.go\t" + sha256A + "\t—",
		"question\t1\tMay a draft pull request be marked ready?\t—\tinterface-contract yes docs/spec/forge.md docs/spec/loop.md",
		"decision\t1\tRead the check runs of the head\t—\tthe merge state of a draft is draft",
		"lesson\t1\tRead the forge before a write\t—\t—",
		"comment\t1\tresult/plan.md\t" + sha256A + "\t—",
		"question\t2\tWhich label is this?\t—\tsome-unknown-label maybe",
		"lesson\t2\tCheck the neighbours of a fix\t—\tdocs/spec/loop.md T-xs4y",
	}
}

func readResult(data []byte) ([][]string, error) { return ReadResult(data) }

func TestAValidResultIsRead(t *testing.T) {
	if _, err := ReadResult(table(resultHeader, resultRows()...)); err != nil {
		t.Fatal(err)
	}
	r := resultRows()
	if _, err := ReadResult(table(resultHeader, r[2], r[0], r[1])); err != nil {
		t.Fatalf("the rows in another order of the file: %v", err)
	}
	if _, err := ReadResult(table(resultHeader, "status\t1\tfailed\t—\tthe harness stopped")); err != nil {
		t.Fatalf("a status alone: %v", err)
	}
}

func TestResultRefusesEachBrokenRule(t *testing.T) {
	if _, err := ReadResult(table(strings.Replace(resultHeader, "reason", "why", 1), resultRows()...)); err == nil {
		t.Error("a wrong header: read, want an error")
	}
	set := func(i int, column, value string) []string {
		r := resultRows()
		r[i] = field(resultHeader, r[i], column, value)
		return r
	}
	for _, c := range []struct {
		name   string
		rows   []string
		column string
	}{
		{"the empty value in value", set(1, "value", "—"), "value"},
		{"a status that is not a word of the block", set(0, "value", "done"), "value"},
		{"an artifact path that is not a path", set(1, "value", "/etc/passwd"), "value"},
		{"an artifact with no sha256", set(1, "sha256", "—"), "sha256"},
		{"a status with a sha256", set(0, "sha256", sha256A), "sha256"},
		{"a status with no reason", set(0, "reason", "—"), "reason"},
		{"an artifact with a reason", set(1, "reason", "it is new"), "reason"},
		{"a gap in the artifacts", set(2, "n", "3"), "n"},
		{"no status row", resultRows()[1:], "kind"},
		{"a status row with n 2", set(0, "n", "2"), "n"},
		{"a question with no text", set(3, "value", "—"), "value"},
		{"a question with a sha256", set(3, "sha256", sha256A), "sha256"},
		{"a question with no reason", set(3, "reason", "—"), "reason"},
		{"a question with a label and no needs_human word", set(3, "reason", "interface-contract"), "reason"},
		{"a decision with a sha256", set(4, "sha256", sha256A), "sha256"},
		{"a decision with no reason", set(4, "reason", "—"), "reason"},
		{"a lesson with a sha256", set(5, "sha256", sha256A), "sha256"},
		{"a comment whose path is not under result/", set(6, "value", "plan.md"), "value"},
		{"a comment whose path is in a directory of result/", set(6, "value", "result/a/plan.md"), "value"},
		{"a comment whose path is the directory", set(6, "value", "result/"), "value"},
		{"a comment with no sha256", set(6, "sha256", "—"), "sha256"},
		{"a comment with a reason", set(6, "reason", "the plan"), "reason"},
		{"a gap in the questions", append(resultRows(), "question\t4\tAnd another?\t—\tlabel no"), "n"},
		{"a gap in the comments", set(6, "n", "2"), "n"},
		{"a gap in the decisions", append(resultRows(), "decision\t3\tAnother\t—\tbecause"), "n"},
		{"a gap in the lessons", append(resultRows(), "lesson\t4\tAnother\t—\t—"), "n"},
		{"a comment whose path is result/..", set(6, "value", "result/.."), "value"},
		{"a comment whose path is result/.", set(6, "value", "result/."), "value"},
		{"a question reason with two spaces", set(3, "reason", "interface-contract  yes"), "reason"},
		{"a question reason with a space at its end", set(3, "reason", "interface-contract yes "), "reason"},
		{"a lesson reason with two spaces", set(5, "reason", "docs/spec/loop.md  T-xs4y"), "reason"},
	} {
		refusedAt(t, c.name, readResult, table(resultHeader, c.rows...), c.column)
	}
	two := append(resultRows(), "status\t2\tfailed\t—\tagain")
	refusedAt(t, "two status rows", readResult, table(resultHeader, two...), "n")
}
