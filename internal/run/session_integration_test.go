//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/session"
)

// The demo of row 36a: against the httptest forge and a local bare
// repository, a task session's start row is pushed before its fake harness
// starts; the harness's first act is to read layup-records of the target.
func TestTheStartRowIsPushedBeforeTheProcess(t *testing.T) {
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
	r := &Sessions{Store: store, RunID: "0123456789abcdef", Host: t.TempDir(), Target: "acme/target", Clone: clone, Branch: "main", Now: time.Now}
	ctx := context.Background()
	n, err := r.StartAttempt(ctx, "T-ab12", base)
	if err != nil || n != 1 {
		t.Fatalf("StartAttempt: %d, %v", n, err)
	}

	seen := t.TempDir()
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte(`#!/bin/sh
git --git-dir="$BARE" show layup-records:sessions.tsv > "$SEEN/sessions"
git --git-dir="$BARE" show layup-records:tasks/T-ab12/events.tsv > "$SEEN/events"
`), 0o755)
	var ended session.Run
	id, err := r.TaskSession(ctx, TaskSpec{Task: "T-ab12", Attempt: 1, Role: "developer", Base: base, Records: read, Prompt: []byte("the task\n"),
		Pair: Pair{Harness: "fake", Model: "m1", VersionCommand: filepath.Join(bin, "version"), Command: filepath.Join(bin, "harness"),
			PromptMode: "stdin", Wall: 1, Rules: []string{"LAYUP-TEST-RULES.md"}, Context: 1000,
			Session: session.Harness{CredentialTo: "—", Vars: []string{"BARE=" + w.bare, "SEEN=" + seen}}},
		Admit: func(string) error { return nil },
		End: func(id string, d session.Dir, run session.Run, err error) error {
			ended = run
			return err
		},
	})
	if err != nil || ended.Exit != 0 || ended.StoppedBy != "" {
		t.Fatalf("TaskSession: %q, %v, the run %+v", id, err, ended)
	}
	sessions, _ := os.ReadFile(filepath.Join(seen, "sessions"))
	events, _ := os.ReadFile(filepath.Join(seen, "events"))
	if !strings.Contains(string(sessions), "\n"+id+"\tT-ab12\t1\tdeveloper\tfake\t1.0.0\tm1\t"+base+"\t") {
		t.Errorf("the harness did not see its start row:\n%s", sessions)
	}
	if !strings.Contains(string(events), "\n2\tsession\t1\t"+id+"\t") {
		t.Errorf("the harness did not see the event session:\n%s", events)
	}
	if got := gitOut(t, w.bare, "log", "-1", "--format=%s", "layup-records"); !strings.Contains(got, "the start of "+id) {
		t.Errorf("the last records commit is %q; want the start row's", got)
	}
}
