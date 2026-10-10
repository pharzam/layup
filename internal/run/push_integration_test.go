//go:build integration

package run

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// change commits the file path with text in repo/, and writes a valid result
// with no artifact.
func change(path, text string) string {
	return `mkdir -p "$(dirname '` + path + `')" && printf '` + text + `' > '` + path + `' && git add -A && git commit -qm change &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `' > ../result/result.tsv`
}

func refusals(t *testing.T, w *world, id string) []string {
	t.Helper()
	var out []string
	for _, l := range strings.Split(show(t, w.bare, "layup-records", "tasks/T-ab12/events.tsv"), "\n") {
		if f := strings.Split(l, "\t"); len(f) == 8 && f[1] == "refused" && f[3] == id {
			out = append(out, f[6])
		}
	}
	return out
}

// The demo of row 37a: a session head that changes a rule path is refused
// before any push, with its diff as a payload.
func TestARulePathChangeIsRefusedWithItsDiff(t *testing.T) {
	w, r, f, a := sessionWorld(t)
	// Two files change, so the payload is the diff of the head, not of the
	// rule path alone.
	id, err := r.TaskSession(context.Background(), harness(t, a, "printf 'more\\n' >> README.md && "+change("docs/rules/new.md", "a rule\\n")))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	got := refusals(t, w, id)
	if len(got) != 1 || !strings.HasPrefix(got[0], "rule-path ") {
		t.Fatalf("the refusals %q; want one rule-path with its payload", got)
	}
	sum := strings.TrimPrefix(got[0], "rule-path ")
	payload := show(t, w.bare, "layup-records", "payloads/"+sum)
	head := strings.TrimSpace(gitOut(t, r.Clone, "rev-parse", "refs/layup/sessions/"+id))
	// gitOut reads no configuration of the host, as internal/git does not.
	if diff := gitOut(t, r.Clone, "diff", "--binary", "--no-renames", a.base, head); payload != diff {
		t.Errorf("the payload:\n%q\nwant git diff --binary from the base to the head:\n%q", payload, diff)
	}
	if s := sha256.Sum256([]byte(payload)); hex.EncodeToString(s[:]) != sum {
		t.Errorf("the payload's name %s is not the SHA-256 of its bytes", sum)
	}
	// The payload and its event are one records commit.
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "layup-records", "--", "payloads/"+sum))
	if changed := strings.Fields(gitOut(t, w.bare, "show", "--format=", "--name-only", commit)); len(changed) != 2 || changed[0] != "payloads/"+sum || changed[1] != "tasks/T-ab12/events.tsv" {
		t.Errorf("the commit of the payload changed %q; want the payload and the events", changed)
	}
	if c := f.posted[2]; len(c) != 1 || !strings.HasSuffix(c[0], ": done, refused rule-path\n") {
		t.Errorf("the comment %q", c)
	}
	// Nothing is pushed: the task branch is not on the target.
	if out, _ := exec.Command("git", "-C", w.bare, "rev-parse", "--verify", "-q", "refs/heads/task/T-ab12/1").Output(); len(out) > 0 {
		t.Errorf("the task branch is on the target: %s", out)
	}
}

func TestAWorkflowChangeIsRefused(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	id, err := r.TaskSession(context.Background(), harness(t, a, change(".github/workflows/ci.yml", "on: push\\n")))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 1 || !strings.HasPrefix(got[0], "workflow ") {
		t.Errorf("the refusals %q; want one workflow with its payload", got)
	}
}

func TestAHeadThatDoesNotDescendFromTheBaseIsRefused(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	body := `git checkout -q --orphan other && printf 'x\n' > x.md && git add -A && git commit -qm other && git branch -qf task/T-ab12/1 other &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `' > ../result/result.tsv`
	id, err := r.TaskSession(context.Background(), harness(t, a, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 1 || got[0] != "base" {
		t.Errorf("the refusals %q; want base", got)
	}
}

func TestAddedLinesInSection2OfTheGuardrailsPass(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	g := "# G\n\n## 2. Known pitfalls\n\n- one\n\n## 3. Next\n"
	gitOut(t, r.Clone, "checkout", "-q", "main")
	writeFile(t, filepath.Join(r.Clone, "docs", "guardrails.md"), g)
	gitOut(t, r.Clone, "add", "-A")
	gitOut(t, r.Clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "guardrails")
	a.base = strings.TrimSpace(gitOut(t, r.Clone, "rev-parse", "HEAD"))
	body := `printf '# G\n\n## 2. Known pitfalls\n\n- one\n- two\n\n## 3. Next\n' > docs/guardrails.md && git add -A && git commit -qm pitfall &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `' > ../result/result.tsv`
	id, err := r.TaskSession(context.Background(), harness(t, a, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 0 {
		t.Errorf("the refusals %q; want none", got)
	}
}

// A head that deletes docs/guardrails.md changes a rule path: the exception
// holds only for added lines, so the result is refused (rule-path).
func TestAHeadThatDeletesTheGuardrailsIsRefused(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	gitOut(t, r.Clone, "checkout", "-q", "main")
	writeFile(t, filepath.Join(r.Clone, "docs", "guardrails.md"), "# G\n\n## 2. Known pitfalls\n\n- one\n")
	gitOut(t, r.Clone, "add", "-A")
	gitOut(t, r.Clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "guardrails")
	a.base = strings.TrimSpace(gitOut(t, r.Clone, "rev-parse", "HEAD"))
	body := `git rm -q docs/guardrails.md && git commit -qm gone &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `' > ../result/result.tsv`
	id, err := r.TaskSession(context.Background(), harness(t, a, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 1 || !strings.HasPrefix(got[0], "rule-path ") {
		t.Errorf("the refusals %q; want one rule-path with its payload", got)
	}
}

// A head whose docs/guardrails.md is a gitlink, to a commit that the clone
// lacks, changes a rule path: an entry that is no file keeps an empty head.
func TestAHeadWhoseGuardrailsIsAGitlinkIsRefused(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	gitOut(t, r.Clone, "checkout", "-q", "main")
	writeFile(t, filepath.Join(r.Clone, "docs", "guardrails.md"), "# G\n\n## 2. Known pitfalls\n\n- one\n")
	gitOut(t, r.Clone, "add", "-A")
	gitOut(t, r.Clone, "-c", "user.name=t", "-c", "user.email=t@t", "commit", "-qm", "guardrails")
	a.base = strings.TrimSpace(gitOut(t, r.Clone, "rev-parse", "HEAD"))
	body := `git rm -q --cached docs/guardrails.md && rm docs/guardrails.md &&
git update-index --add --cacheinfo 160000,1111111111111111111111111111111111111111,docs/guardrails.md && git commit -qm gitlink &&
printf '` + strings.ReplaceAll(resultHead, "\t", `\t`) + `' > ../result/result.tsv`
	id, err := r.TaskSession(context.Background(), harness(t, a, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 1 || !strings.HasPrefix(got[0], "rule-path ") {
		t.Errorf("the refusals %q; want one rule-path with its payload", got)
	}
}

func TestARecordsCommitWithNoRegisterRefuses(t *testing.T) {
	for _, broken := range []bool{false, true} {
		w, r, _, a := sessionWorld(t)
		a.records = a.base // a commit of main: no rule-paths.tsv
		if broken {
			if err := r.commit(context.Background(), map[string][]byte{"rule-paths.tsv": []byte("broken\n")}, "a broken register"); err != nil {
				t.Fatal(err)
			}
			a.records = r.Store.(*recordsStore).pushed
		}
		id, err := r.TaskSession(context.Background(), harness(t, a, change("README.md", "changed\\n")))
		if err != nil {
			t.Fatalf("TaskSession: %q, %v", id, err)
		}
		if got := refusals(t, w, id); len(got) != 1 || got[0] != "rule-paths" {
			t.Errorf("a register that is broken %v: the refusals %q; want rule-paths", broken, got)
		}
	}
}

// A session whose result is refused by its end is not checked: a refused
// artifact and a rule-path change give the one refusal of the end.
func TestASessionRefusedAtItsEndIsNotChecked(t *testing.T) {
	w, r, _, a := sessionWorld(t)
	body := `mkdir -p docs/rules && printf 'r\n' > docs/rules/x.md && git add -A && git commit -qm rule && ` + work(strings.Repeat("0", 64))
	id, err := r.TaskSession(context.Background(), harness(t, a, body))
	if err != nil {
		t.Fatalf("TaskSession: %q, %v", id, err)
	}
	if got := refusals(t, w, id); len(got) != 1 || got[0] != "artifact" {
		t.Errorf("the refusals %q; want artifact only", got)
	}
}

func writeFile(t *testing.T, path, text string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
}
