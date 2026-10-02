//go:build integration

package standin

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

func isolate(t *testing.T) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
}

// The stand-in baseline builds with no network, and git ls-remote of its
// file:// URL gives its one commit (verify-test-baseline). It holds the
// history and the phrases of the kit, which a setup removes.
func TestAStandInBaseline(t *testing.T) {
	isolate(t)
	dir := filepath.Join(t.TempDir(), "baseline")
	url, commit, err := Baseline(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := git.LsRemote(url, "HEAD"); err != nil || got != commit || !strings.HasPrefix(url, "file://") {
		t.Fatalf("ls-remote %s HEAD: %q, %v; want the commit %s", url, got, err, commit)
	}
	if roots, err := git.RootCommits(dir, "HEAD"); err != nil || len(roots) != 1 || roots[0] != commit {
		t.Errorf("the root commits %q, %v; want the one commit", roots, err)
	}
	files, err := git.LsFiles(dir)
	for _, f := range []string{"README.md", "AGENTS.md", "docs/decisions/D-0001-stand-in.md", "docs/audit/README.md",
		"docs/tasks/T-0001.md", "docs/tasks/backlog.md", "docs/tasks/completed.md"} {
		if err != nil || !slices.Contains(files, f) {
			t.Errorf("the baseline has no file %s (%v)", f, err)
		}
	}
	backlog, _ := os.ReadFile(filepath.Join(dir, "docs/tasks/backlog.md"))
	if !strings.Contains(string(backlog), strings.TrimPrefix(url, "file://")+"/issues/1") {
		t.Errorf("backlog.md does not link the baseline:\n%s", backlog)
	}
}

// A work area: the root commit of the unchanged baseline on main, the setup
// by hand on layup-setup, and the two records by their schemas; each option
// changes its part.
func TestAStandInWorkArea(t *testing.T) {
	isolate(t)
	w, err := Make(t.TempDir(), Options{})
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(w.Dir, work.TargetPath)
	roots, err := git.RootCommits(target, "layup-setup")
	if err != nil || len(roots) != 1 {
		t.Fatalf("the root commits of layup-setup: %q, %v; want one", roots, err)
	}
	if tree, err := git.RevParse(target, roots[0]+"^{tree}"); err != nil || tree != w.Tree {
		t.Errorf("the tree of the root commit %s, %v; want the pinned tree %s", tree, err, w.Tree)
	}
	if main, err := git.RevParse(target, "refs/heads/main"); err != nil || main != roots[0] {
		t.Errorf("main is %s, %v; want the root commit %s", main, err, roots[0])
	}
	show := func(path string) string {
		b, err := git.Show(target, "layup-setup", path)
		if err != nil {
			return ""
		}
		return string(b)
	}
	if pin := show("docs/setup/armature.pin"); pin != pinText(w.URL, w.Commit, w.Tree) {
		t.Errorf("the pin at the head of layup-setup:\n%s", pin)
	}
	if show("docs/decisions/D-0001-stand-in.md") != "" || show("docs/tasks/T-0001.md") != "" || strings.Contains(show("docs/tasks/backlog.md"), "/issues/") {
		t.Error("the history of the baseline is still at the head of layup-setup")
	}
	if readme := show("README.md"); !strings.Contains(readme, Name) || !strings.Contains(readme, "](docs/setup/armature.pin)") {
		t.Errorf("README.md at the head of layup-setup does not name the target and link the pin:\n%s", readme)
	}
	if show("docs/gates.tsv") == "" {
		t.Error("no docs/gates.tsv at the head of layup-setup")
	}
	record, err := work.ReadRecord(w.Dir)
	if v, _ := record.Value("S02", "pin.commit"); err != nil || v != w.Commit {
		t.Errorf("the record: pin.commit %q, %v; want %s", v, err, w.Commit)
	}
	if answers, err := work.ReadAnswers(w.Dir); err != nil || len(answers) == 0 {
		t.Errorf("the answers: %d rows, %v", len(answers), err)
	}

	w, err = Make(t.TempDir(), Options{Files: map[string]string{"docs/decisions/x.md": "x\n", "README.md": ""},
		Record: map[string]string{"S02 pin.commit": strings.Repeat("1", 40), "S01 name": ""}})
	if err != nil {
		t.Fatal(err)
	}
	target = filepath.Join(w.Dir, work.TargetPath)
	if show("docs/decisions/x.md") != "x\n" || show("README.md") != "" {
		t.Error("the option Files did not change the head of layup-setup")
	}
	record, _ = work.ReadRecord(w.Dir)
	if v, _ := record.Value("S02", "pin.commit"); v != strings.Repeat("1", 40) {
		t.Errorf("the option Record did not change pin.commit: %q", v)
	}
	if _, ok := record.Value("S01", "name"); ok {
		t.Error("the option Record did not remove the row S01 name")
	}
}
