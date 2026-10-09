//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/session"
)

// sessionWorld gives the world of Start, the Sessions of its run on a clone
// of the target, with the stand-in forge, and the base: attempt 1 of T-ab12 is
// started.
func sessionWorld(t *testing.T) (*world, *Sessions, *standInForge, string) {
	t.Helper()
	w := newWorld(t)
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
	return w, r, f, base
}

const resultHead = "kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n"

// harness gives a TaskSpec whose fake harness runs body in repo/.
func harness(t *testing.T, base, body string) TaskSpec {
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte("#!/bin/sh\n"+body+"\n"), 0o755)
	return TaskSpec{Task: "T-ab12", Attempt: 1, Role: "developer", Base: base, Records: base, Prompt: []byte("the task\n"),
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
	changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", "layup-records"))
	if !slices.Equal(changed, []string{"tasks/T-ab12/events.tsv", "tasks/T-ab12/results/" + id + ".tsv", "telemetry.tsv"}) {
		t.Errorf("the last records commit changed %q; want the result, the events and the telemetry row in one", changed)
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
