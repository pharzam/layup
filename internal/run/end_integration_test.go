//go:build integration

package run

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/session"
)

// worldAt is the base commit of a session and the records commit of its prompt,
// which holds the rule-path register.
type worldAt struct{ base, records string }

// register is the rule-path register of the test world.
const register = "pattern\texception\tsource\nAGENTS.md\t—\tbaseline\ndocs/guardrails.md\tadded lines in section 2\tbaseline\ndocs/rules/\t—\tarchitecture\n"

// sessionWorld gives the world of Start, the Sessions of its run on a clone
// of the target, with the stand-in forge, the base and a records commit with
// the register: attempt 1 of T-ab12 is started.
func sessionWorld(t *testing.T) (*world, *Sessions, *standInForge, worldAt) {
	return sessionWorldWith(t, nil)
}

// sessionWorldWith is sessionWorld whose baseline, the target's first tree, also holds files by path.
func sessionWorldWith(t *testing.T, files map[string]string) (*world, *Sessions, *standInForge, worldAt) {
	t.Helper()
	w := newWorldWith(t, files)
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	clone := filepath.Join(t.TempDir(), "clone")
	gitOut(t, "/", "clone", "-q", "-b", "main", "file://"+w.bare, clone)
	read := strings.TrimSpace(gitOut(t, clone, "rev-parse", "origin/layup-records"))
	base := strings.TrimSpace(gitOut(t, clone, "rev-parse", "main"))
	store := &recordsStore{dir: clone, url: "file://" + w.bare, read: read,
		auth: func(context.Context) (git.Auth, error) { return git.Auth{}, nil },
		who: func() git.Identity {
			return git.Identity{Name: "layup-agent[bot]", Email: "9001+layup-agent[bot]@users.noreply.github.com", Time: time.Now()}
		}}
	f := &standInForge{}
	r := &Sessions{Store: store, RunID: "0123456789abcdef", Host: t.TempDir(), Target: "acme/target", Clone: clone, Branch: "main",
		Now: time.Now, Forge: f, Control: 2}
	if n, err := r.StartAttempt(context.Background(), "T-ab12", base); err != nil || n != 1 {
		t.Fatalf("StartAttempt: %d, %v", n, err)
	}
	if err := r.commit(context.Background(), map[string][]byte{"rule-paths.tsv": []byte(register)}, "the register"); err != nil {
		t.Fatal(err)
	}
	// The task register: the comment of a task session is on the task's issue (task T-syxm).
	rows, err := registerFiles(nil, nil, Task{ID: "T-ab12", Issue: taskIssue, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", Opened: time.Now().UTC().Truncate(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if err := r.commit(context.Background(), rows, "tasks: register T-ab12"); err != nil {
		t.Fatal(err)
	}
	return w, r, f, worldAt{base, store.pushed}
}

const resultHead = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n"

// harness gives a TaskSpec whose fake harness runs body in repo/.
func harness(t *testing.T, a worldAt, body string) TaskSpec {
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte("#!/bin/sh\n"+body+"\n"), 0o755)
	return TaskSpec{Task: "T-ab12", Attempt: 1, Role: "developer", Base: a.base, Records: a.records, Prompt: []byte("the task\n"),
		Pair: Pair{Harness: "fake", Model: "m1", VersionCommand: filepath.Join(bin, "version"), Command: filepath.Join(bin, "harness"),
			PromptMode: "stdin", Wall: 1, Rules: []string{"LAYUP-TEST-RULES.md"}, Context: 1000,
			Session: session.Harness{CredentialTo: "—"}, Billing: "subscription", Usage: "none"},
		Admit: func(string) error { return nil }}
}

// work commits out.txt in repo/ on the task branch, and writes the result
// file with the artifact's SHA-256 sum.
func work(sum string) string {
	return `printf 'hello\n' > out.txt && git add out.txt && git commit -qm work &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `artifact\t1\tout.txt\t` + sum + `\t—\n' > ../result/result.tsv`
}

const helloSum = "5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03"

// The demo of row 36b: the result file of a task session is committed byte
// for byte as tasks/<task>/results/<session>.tsv, in the records commit that
// holds its telemetry row; the comment is posted and the directory removed.
func TestTheResultIsCommittedWithItsTelemetryRow(t *testing.T) {
	w, r, f, base := sessionWorld(t)
	id, err := r.TaskSession(context.Background(), harness(t, base, work(helloSum)))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	want := resultHead + "artifact\t1\tout.txt\t" + helloSum + "\t—\n"
	if got := show(t, w.bare, "layup-records", "tasks/T-ab12/results/"+id+".tsv"); got != want {
		t.Errorf("the result committed:\n%q\nwant byte for byte:\n%q", got, want)
	}
	// The commit of the result; the push and the bind come after it.
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "layup-records", "--", "tasks/T-ab12/results/"+id+".tsv"))
	changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", commit))
	if !slices.Equal(changed, []string{"tasks/T-ab12/events.tsv", "tasks/T-ab12/results/" + id + ".tsv", "telemetry.tsv"}) {
		t.Errorf("the commit of the result changed %q; want the result, the events and the telemetry row in one", changed)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\tT-ab12\t—\tdeveloper\tfake\tm1\tsubscription\t") {
		t.Errorf("telemetry.tsv has no row of %s:\n%s", id, tel)
	}
	events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
	if !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tdone\t") || strings.Contains(events, "refused") {
		t.Errorf("the events:\n%s", events)
	}
	if got := f.posted[2]; len(got) != 1 || got[0] != id+": developer session of T-ab12, attempt 1: done\n" {
		t.Errorf("the comments on the control issue %q", got)
	}
	if _, err := os.Stat(filepath.Join(r.Host, "sessions", id)); !os.IsNotExist(err) {
		t.Errorf("the session directory after the call: %v; want it removed", err)
	}
}

func TestAnArtifactThatDiffersIsRefusedAndKept(t *testing.T) {
	w, r, f, base := sessionWorld(t)
	id, err := r.TaskSession(context.Background(), harness(t, base, work(strings.Repeat("0", 64))))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
	if !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tdone\t") || !strings.Contains(events, "\trefused\t1\t"+id+"\t—\t—\tartifact\t") {
		t.Errorf("the events:\n%s", events)
	}
	if gitOut(t, w.bare, "ls-tree", "--name-only", "layup-records", "tasks/T-ab12/results/"+id+".tsv") == "" {
		t.Error("the result of a differing artifact is not committed; want it kept")
	}
	if got := f.posted[2]; len(got) != 1 || got[0] != id+": developer session of T-ab12, attempt 1: done, refused artifact\n" {
		t.Errorf("the comment %q", got)
	}
}

func TestNoValidResultIsNoResultFile(t *testing.T) {
	w, r, _, base := sessionWorld(t)
	id, err := r.TaskSession(context.Background(), harness(t, base, "exit 0"))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
	if !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tno-result\t") {
		t.Errorf("the events:\n%s", events)
	}
	if gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records", "tasks/T-ab12/results") != "" {
		t.Error("a results/ file with no valid result file")
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\t") {
		t.Errorf("no telemetry row of %s", id)
	}
}

func TestAHeadThatCannotBeReadIsRefusedBranch(t *testing.T) {
	w, r, f, base := sessionWorld(t)
	// The session removes the ref of its own branch: no head to fetch.
	body := work(helloSum) + " && git update-ref -d refs/heads/task/T-ab12/1"
	id, err := r.TaskSession(context.Background(), harness(t, base, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
	if !strings.Contains(events, "\trefused\t1\t"+id+"\t—\t—\tbranch\t") {
		t.Errorf("the events:\n%s", events)
	}
	if gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records", "tasks/T-ab12/results") != "" {
		t.Error("a results/ file of a head that cannot be read")
	}
	if got := f.posted[2]; len(got) != 1 || !strings.HasSuffix(got[0], ": done, refused branch\n") {
		t.Errorf("the comment %q", got)
	}
}

// A closed attempt seen at the end: the end hook commits an event closed of
// the attempt first, as the task loop of M2e would, then ends the session.
func TestAResultOfAClosedAttemptIsRefused(t *testing.T) {
	w, r, f, base := sessionWorld(t)
	spec := harness(t, base, work(helloSum))
	ctx := context.Background()
	spec.End = func(id string, d session.Dir, run session.Run, err error) error {
		data, aErr := r.appendEvents(ctx, "T-ab12", 1, "", []string{"closed —"})
		if aErr != nil {
			return aErr
		}
		data = []byte(strings.Replace(string(data), "\tclosed\t1\t—\t—\t—\t—\t", "\tclosed\t1\t—\t—\t—\tby the test\t", 1))
		if cErr := r.commit(ctx, map[string][]byte{"tasks/T-ab12/events.tsv": data}, "closed by the test"); cErr != nil {
			return cErr
		}
		return r.endTask(ctx, spec, id, d, run, err)
	}
	id, err := r.TaskSession(ctx, spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
	if !strings.Contains(events, "\trefused\t1\t"+id+"\t—\t—\tclosed-attempt\t") {
		t.Errorf("the events:\n%s", events)
	}
	if gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records", "tasks/T-ab12/results") != "" {
		t.Error("a results/ file of a closed attempt")
	}
	if got := f.posted[2]; len(got) != 1 || !strings.HasSuffix(got[0], ": done, refused closed-attempt\n") {
		t.Errorf("the comment %q", got)
	}
}

func TestAnArtifactThatTheHeadLacksIsRefused(t *testing.T) {
	w, r, _, base := sessionWorld(t)
	body := strings.Replace(work(helloSum), `artifact\t1\tout.txt`, `artifact\t1\tnone.txt`, 1)
	id, err := r.TaskSession(context.Background(), harness(t, base, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv"); !strings.Contains(events, "\trefused\t1\t"+id+"\t—\t—\tartifact\t") {
		t.Errorf("the events:\n%s", events)
	}
}

// A program that cannot start is the class start, with its telemetry row:
// its end is its start, as it has no end of its own.
func TestAProgramThatCannotStartHasItsRow(t *testing.T) {
	w, r, _, base := sessionWorld(t)
	spec := harness(t, base, "exit 0")
	spec.Pair.Command = filepath.Join(t.TempDir(), "none")
	id, err := r.TaskSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv"); !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tstart\t") {
		t.Errorf("the events:\n%s", events)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\n"+id+"\t") {
		t.Errorf("no telemetry row of %s:\n%s", id, tel)
	}
}

// The comment comes after every record of the session, the push hook's too:
// a refusal that the push hook records (row 37a's base, here a stand-in) is
// named in its first line.
func TestTheCommentNamesARefusalOfThePushHook(t *testing.T) {
	_, r, f, base := sessionWorld(t)
	spec := harness(t, base, work(helloSum))
	ctx := context.Background()
	spec.Push = func(id string, d session.Dir, run session.Run) error {
		data, err := r.appendEvents(ctx, "T-ab12", 1, id, []string{"refused base"})
		if err != nil {
			return err
		}
		return r.commit(ctx, map[string][]byte{"tasks/T-ab12/events.tsv": data}, "refused by the test")
	}
	id, err := r.TaskSession(ctx, spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := f.posted[2]; len(got) != 1 || got[0] != id+": developer session of T-ab12, attempt 1: done, refused base\n" {
		t.Errorf("the comment %q; want it after the push hook's refusal", got)
	}
}

// A second session of the attempt: its comment names its own result only,
// not the refusal of the first.
func TestTheCommentOfASecondSessionIsItsOwn(t *testing.T) {
	_, r, f, base := sessionWorld(t)
	ctx := context.Background()
	if _, err := r.TaskSession(ctx, harness(t, base, work(strings.Repeat("0", 64)))); err != nil {
		t.Fatal(err)
	}
	id, err := r.TaskSession(ctx, harness(t, base, work(helloSum)))
	if err != nil {
		t.Fatalf("the second session: %v", err)
	}
	if got := f.posted[2]; len(got) != 2 || got[1] != id+": developer session of T-ab12, attempt 1: done\n" {
		t.Errorf("the comments %q; want the second one its own", got)
	}
}

// A usage format that is not of the list is the run's own error, not a row
// with a false reason.
func TestAUsageFormatOfNoListIsAnError(t *testing.T) {
	_, r, _, base := sessionWorld(t)
	spec := harness(t, base, work(helloSum))
	spec.Pair.Usage = "other"
	if _, err := r.TaskSession(context.Background(), spec); err == nil || !strings.Contains(err.Error(), "not one of the list") {
		t.Errorf("%v; want the error of the format", err)
	}
}

// A stdout that cannot be read gives tokens unavailable, with its reason.
func TestAStdoutThatCannotBeReadIsUnavailable(t *testing.T) {
	w, r, _, base := sessionWorld(t)
	spec := harness(t, base, work(helloSum)+" && rm ../stdout")
	spec.Pair.Usage = "claude-result"
	id, err := r.TaskSession(context.Background(), spec)
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if tel := show(t, w.bare, "layup-records", "telemetry.tsv"); !strings.Contains(tel, "\tunavailable\tstdout cannot be read\t") {
		t.Errorf("telemetry.tsv:\n%s", tel)
	}
}

// octal writes text as the escapes of printf, byte for byte.
func octal(text string) string {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		fmt.Fprintf(&b, "\\%03o", text[i])
	}
	return b.String()
}

// planWork writes result/<name> with text, and a result file whose status row
// is followed by rows (the lines of the table, with real tabs), in the shell of
// a harness; it commits nothing.
func planWork(name, text string, rows ...string) string {
	out := resultHead
	for _, r := range rows {
		out += r + "\n"
	}
	return "printf '" + octal(text) + "' > ../result/" + name + " && printf '" + octal(out) + "' > ../result/result.tsv"
}

// planText is a valid plan whose bytes keep a \r, a blank line and a last line with no end.
const planText = "size: small\r\nopen questions: —\nnew interfaces: —\nnew dependencies: —\ntests:\n# The plan\r\n\r\nno final newline, trailing  \t"

func keptFiles(t *testing.T, w *world, id string) string {
	t.Helper()
	return gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records", "tasks/T-ab12/results/"+id)
}

// The end of a session that names a comment: the file of result/ is checked
// against the row's hash and kept byte for byte, in the one records commit of
// the result; the session posts nothing (loop.md, The artifact that is a comment).
func TestTheCommentRowOfASessionIsKeptWithTheResultInOneCommit(t *testing.T) {
	w, r, f, base := sessionWorld(t)
	row := "comment\t1\tresult/plan.md\t" + commentSum(planText) + "\t—"
	id, err := r.TaskSession(context.Background(), harness(t, base, planWork("plan.md", planText, row)))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := show(t, w.bare, "layup-records", "tasks/T-ab12/results/"+id+"/plan.md"); got != planText {
		t.Errorf("the kept file %q; want the bytes %q", got, planText)
	}
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "layup-records", "--", "tasks/T-ab12/results/"+id+".tsv"))
	changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", commit))
	want := []string{"tasks/T-ab12/events.tsv", "tasks/T-ab12/results/" + id + ".tsv", "tasks/T-ab12/results/" + id + "/plan.md", "telemetry.tsv"}
	if !slices.Equal(changed, want) { // git lists a path by byte order: ".tsv" before "/"
		t.Errorf("the commit of the result changed %q; want %q in one", changed, want)
	}
	if events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv"); strings.Contains(events, "refused") {
		t.Errorf("the events:\n%s", events)
	}
	if len(f.posted) != 1 || len(f.posted[2]) != 1 {
		t.Errorf("the session's end posted %v; want only the comment of the session on the control issue", f.posted)
	}
}

// A comment row whose file differs, is missing or is a link adds the event
// refused, artifact (session.md, The result of a session): the result is still
// committed, as for an artifact, and that row's file is not kept.
func TestACommentFileThatFailsItsRowIsRefusedArtifactAndNotKept(t *testing.T) {
	for _, c := range []struct{ name, body string }{
		{"differs", planWork("plan.md", planText, "comment\t1\tresult/plan.md\t"+commentSum("another")+"\t—")},
		{"missing", planWork("plan.md", planText, "comment\t1\tresult/none.md\t"+commentSum(planText)+"\t—")},
		{"a link", planWork("real.md", planText, "comment\t1\tresult/plan.md\t"+commentSum(planText)+"\t—") + " && ln -s real.md ../result/plan.md"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, _, base := sessionWorld(t)
			id, err := r.TaskSession(context.Background(), harness(t, base, c.body))
			if err != nil {
				t.Fatalf("TaskSession: %q, %v", id, err)
			}
			events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv")
			if !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tdone\t") || !strings.Contains(events, "\trefused\t1\t"+id+"\t—\t—\tartifact\t") {
				t.Errorf("the events:\n%s", events)
			}
			if gitOut(t, w.bare, "ls-tree", "--name-only", "layup-records", "tasks/T-ab12/results/"+id+".tsv") == "" {
				t.Error("the result is not committed; want it kept, as for an artifact")
			}
			if got := keptFiles(t, w, id); got != "" {
				t.Errorf("kept %q; want no file of a row that fails", got)
			}
		})
	}
}

// A comment row that holds a comment ID, in its value or its reason, is not
// a valid result: the class is no-result and nothing is kept (loop.md, "The
// comment ID is never in the session's result").
func TestAResultThatHoldsACommentIDIsNoResult(t *testing.T) {
	for _, c := range []struct{ name, row string }{
		{"in the reason", "comment\t1\tresult/plan.md\t" + commentSum(planText) + "\t6213845433-1"},
		{"as the value", "comment\t1\t6213845433\t" + commentSum(planText) + "\t—"},
	} {
		t.Run(c.name, func(t *testing.T) {
			w, r, _, base := sessionWorld(t)
			id, err := r.TaskSession(context.Background(), harness(t, base, planWork("plan.md", planText, c.row)))
			if err != nil {
				t.Fatalf("TaskSession: %q, %v", id, err)
			}
			if events := show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv"); !strings.Contains(events, "\tresult\t1\t"+id+"\t—\t—\tno-result\t") {
				t.Errorf("the events:\n%s", events)
			}
			if got := keptFiles(t, w, id); got != "" {
				t.Errorf("kept %q; want nothing", got)
			}
		})
	}
}
