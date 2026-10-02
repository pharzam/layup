//go:build integration

package verify

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
)

// The fixture harness (D6 of #84): the core of each built check runs on the
// cases of docs/setup/tests/, which run.sh runs through setup-check.sh, and
// gives the same lines. Each group of docs/setup/tests/ is in exactly one
// list; a later row of the plan moves a group from notYetBuilt to built.
var (
	built = map[string]func(fsys fs.FS, h history) []string{
		"pin":         pinFindings,
		"kit-history": func(fsys fs.FS, _ history) []string { return kitHistoryFindings(fsys, "github.com/pharzam/armature") },
		"identity":    func(fsys fs.FS, _ history) []string { return identityFindings(fsys) },
	}
	notYetBuilt    = []string{"facts", "onboarding", "glossary", "guardrails", "markers", "adapted"}
	notForATarget  = []string{"ci", "procedure", "protection"}
	fixtureCommits = git.Identity{Name: "fixture", Email: "fixture@invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
)

func TestTheFixturesOfSetupCheck(t *testing.T) {
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1"} {
		t.Setenv(k, v)
	}
	root := filepath.Join("..", "..", "docs", "setup", "tests")
	entries, err := os.ReadDir(root)
	if err != nil {
		t.Fatal(err)
	}
	var groups []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "kit" {
			groups = append(groups, e.Name())
		}
	}
	listed := append(append(append([]string{"frame"}, notYetBuilt...), notForATarget...), "pin", "kit-history", "identity")
	for _, g := range groups {
		if n := len(slices.DeleteFunc(slices.Clone(listed), func(l string) bool { return l != g })); n != 1 {
			t.Errorf("the group %s is in %d lists; want 1", g, n)
		}
	}
	for _, l := range listed {
		if !slices.Contains(groups, l) {
			t.Errorf("the listed group %s has no directory under docs/setup/tests/", l)
		}
	}
	cases := 0
	for _, g := range append([]string{"frame"}, "pin", "kit-history", "identity") {
		dirs, err := filepath.Glob(filepath.Join(root, g, "*", "EXPECT"))
		if err != nil {
			t.Fatal(err)
		}
		for _, expectPath := range dirs {
			cases++
			dir := filepath.Dir(expectPath)
			name := g + "/" + filepath.Base(dir)
			expect := readExpect(t, expectPath)
			repo := fixtureRepo(t, filepath.Join(root, "kit"), filepath.Join(dir, "overlay"), expect["mode"])
			run := []string{g}
			if o := expect["only"]; len(o) > 0 {
				run = strings.Split(o[0], ",")
			} else if g == "frame" {
				run = []string{"pin", "kit-history", "identity"}
			}
			found := false
			for _, check := range run {
				core, ok := built[check]
				if !ok {
					continue
				}
				var want []string
				for _, l := range expect["line"] {
					if strings.HasPrefix(l, "setup-check: "+check+" ") {
						want = append(want, l)
					}
				}
				findings := core(os.DirFS(repo), gitHistory{repo})
				found = found || len(findings) > 0
				if len(want) == 0 && g != "frame" {
					continue
				}
				if len(want) == 0 { // a frame overlay holds a good setup for each check that it runs (round 1, finding 4)
					want = []string{"setup-check: " + check + " OK"}
				}
				got := []string{"setup-check: " + check + " OK"}
				if len(findings) > 0 {
					got = nil
					for _, f := range findings {
						got = append(got, "setup-check: "+check+" FAIL "+f)
					}
				}
				slices.Sort(got)
				slices.Sort(want)
				if !slices.Equal(got, want) {
					t.Errorf("%s: the lines of %s\n got %q\nwant %q", name, check, got, want)
				}
			}
			if exit := map[bool]string{false: "0", true: "1"}[found]; g != "frame" && (len(expect["exit"]) != 1 || expect["exit"][0] != exit) {
				t.Errorf("%s: exit %s; want %q", name, exit, expect["exit"])
			}
		}
	}
	if cases == 0 {
		t.Fatal("no case found: a harness that tested nothing is not a pass")
	}
}

// readExpect reads the key=value lines of an EXPECT file.
func readExpect(t *testing.T, path string) map[string][]string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string][]string{}
	for _, l := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		if k, v, ok := strings.Cut(l, "="); ok {
			m[k] = append(m[k], v)
		}
	}
	return m
}

// fixtureRepo builds the repository of a case as run.sh does: the tree kit/
// as the root commit, overlay/ as a second commit, and for the mode shallow a
// clone of depth 1, which the test starts with git itself (K8).
func fixtureRepo(t *testing.T, kit, overlay string, mode []string) string {
	t.Helper()
	repo := filepath.Join(t.TempDir(), "repo")
	for _, err := range []error{git.Init(repo), copyTree(repo, kit), git.Add(repo), git.Commit(repo, "kit", fixtureCommits)} {
		if err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(overlay); err == nil {
		for _, err := range []error{copyTree(repo, overlay), git.Add(repo), git.Commit(repo, "overlay", fixtureCommits)} {
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	if len(mode) == 1 && mode[0] == "shallow" {
		shallow := filepath.Join(t.TempDir(), "shallow")
		if out, err := exec.Command("git", "-c", "core.hooksPath=/dev/null", "clone", "-q", "--depth", "1", "file://"+repo, shallow).CombinedOutput(); err != nil {
			t.Fatalf("git clone --depth 1: %v\n%s", err, out)
		}
		return shallow
	}
	return repo
}

// copyTree copies each file under src into dst, over a file of the same name,
// as cp -R does in run.sh.
func copyTree(dst, src string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(dst, filepath.Dir(rel)), 0o755); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), data, 0o644)
	})
}
