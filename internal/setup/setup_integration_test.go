//go:build integration

package setup

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// The Go schemas of the two tables equal their blocks.
func TestTheSchemasEqualTheirBlocks(t *testing.T) {
	blocks, err := tsv.ReadBlocks(os.DirFS(filepath.Join("..", "..", "docs", "spec")))
	if err != nil {
		t.Fatal(err)
	}
	for name, s := range map[string]tsv.Schema{"setup-steps": StepsSchema, "setup-stop": StopSchema} {
		block, ok := blocks[name]
		if !ok {
			t.Errorf("docs/spec/ has no block %s", name)
			continue
		}
		if err := tsv.Compare(block, s); err != nil {
			t.Error(err)
		}
	}
}

func gitOut(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	if err != nil {
		t.Fatalf("git %s: %v", strings.Join(args, " "), err)
	}
	return strings.TrimSpace(string(out))
}

// On a real work area, a step that changes the tree makes one commit
// "chore: setup <step>" on layup-setup, by the identity of the run at the
// date pin.time; a second run makes none. A host with a global hook that
// fails, commit.gpgsign=true and core.autocrlf=true changes nothing (K7).
// commands.sh is a shell file that sh -n accepts (D7).
func TestARunOnARealWorkArea(t *testing.T) {
	home := t.TempDir()
	hooks := filepath.Join(home, "hooks")
	if err := os.MkdirAll(hooks, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hooks, "pre-commit"), []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	gitconfig := "[core]\n\thooksPath = " + hooks + "\n\tautocrlf = true\n[commit]\n\tgpgsign = true\n"
	if err := os.WriteFile(filepath.Join(home, ".gitconfig"), []byte(gitconfig), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	target := filepath.Join(w.Dir, work.TargetPath)
	before := gitOut(t, target, "rev-parse", "layup-setup")
	m := doneSteps()
	res, err := Run(w.Dir, m, testWho, noStep)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(res.Results(), " "); got != strings.TrimSpace(strings.Repeat("done ", 12)+"operator done done") {
		t.Errorf("the results %q", got)
	}
	if n := gitOut(t, target, "rev-list", "--count", before+"..layup-setup"); n != "1" {
		t.Fatalf("%s new commits on layup-setup, want 1", n)
	}
	if got := gitOut(t, target, "log", "-1", "--format=%s|%an <%ae>|%aI|%cI", "layup-setup"); got != "chore: setup S05|LAYUP test <test@layup.invalid>|2026-10-02T09:30:00Z|2026-10-02T09:30:00Z" {
		t.Errorf("the commit %q", got)
	}
	if blob := gitOut(t, target, "cat-file", "-p", "layup-setup:x.md"); blob != "one\ntwo" {
		t.Errorf("the blob %q; want the bytes of the step, with no conversion", blob)
	}
	head := gitOut(t, target, "rev-parse", "layup-setup")
	if _, err := Run(w.Dir, m, testWho, noStep); err != nil || gitOut(t, target, "rev-parse", "layup-setup") != head {
		t.Errorf("a second run: %v; want no new commit", err)
	}
	cmds := filepath.Join(w.Dir, filepath.FromSlash(work.CommandsPath))
	if out, err := exec.Command("sh", "-n", cmds).CombinedOutput(); err != nil {
		t.Errorf("sh -n %s: %v\n%s", cmds, err, out)
	}
}

// doneSteps gives the stubs, each done: S05 writes x.md and asks for a commit,
// and S05 and S13 give a command for the Operator.
func doneSteps() map[string]Step {
	m := Stubs()
	for _, id := range ids() {
		s := m[id]
		s.Run = func(Input) Outcome { return Outcome{Kind: Done, Evidence: "ok " + id} }
		switch id {
		case "S05":
			s.Run = func(in Input) Outcome {
				if err := os.WriteFile(filepath.Join(in.Dir, work.TargetPath, "x.md"), []byte("one\ntwo\n"), 0o644); err != nil {
					return Outcome{Kind: Fail, Evidence: err.Error()}
				}
				return Outcome{Kind: Done, Evidence: "x.md", Commit: true}
			}
			s.Commands = func(work.Record) []Command {
				return []Command{{Order: 1, Comment: "the push of the root commit", Text: "git -C target push origin main"}}
			}
		case "S13":
			s.Commands = func(work.Record) []Command {
				return []Command{{Order: 2, Comment: "the push of layup-setup", Text: "git -C target push origin layup-setup:main"}}
			}
		}
		m[id] = s
	}
	return m
}

// A target that is not on the branch layup-setup gets no commit of a step: the
// step fails, and no branch moves (finding 2 of round 1).
func TestNoCommitOffTheSetupBranch(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(w.Dir, work.TargetPath)
	gitOut(t, target, "switch", "-q", "main")
	before := gitOut(t, target, "rev-parse", "main", "layup-setup")
	res, err := Run(w.Dir, doneSteps(), testWho, noStep)
	if err != nil || len(res.Steps) != 15 || res.Steps[4] != (StepRow{"S05", "layup-setup", "fail", "the commit of the step failed: the target is not on the branch layup-setup"}) {
		t.Errorf("a target on main: %v, the rows %q; want S05 fail", err, res.Steps)
	}
	if after := gitOut(t, target, "rev-parse", "main", "layup-setup"); after != before {
		t.Errorf("the branches moved: %q, then %q", before, after)
	}
}
