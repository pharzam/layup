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
		"adapted":     adaptedFindings,
		"identity":    func(fsys fs.FS, _ history) []string { return identityFindings(fsys) },
		"facts":       func(fsys fs.FS, _ history) []string { return factsHashFindings(fsys) },
		"onboarding":  func(fsys fs.FS, _ history) []string { return onboardingFindings(fsys) },
		"glossary":    func(fsys fs.FS, _ history) []string { return glossaryFindings(fsys) },
		"guardrails":  func(fsys fs.FS, _ history) []string { return guardrailsFindings(fsys) },
	}
	builtGroups = []string{"pin", "kit-history", "facts", "onboarding", "glossary", "guardrails", "adapted", "identity"} // the keys of built, in the order of the table
	notYetBuilt = []string{"markers"}
	// targetKinds names, for each check in a target's form, the kinds of the
	// lines of its sh function that the target keeps with the same text (D9
	// of #89); the other kinds are LAYUP's form, which the engine does not run.
	targetKinds = map[string][]string{"facts": {"hash"}, "onboarding": {"missing", "marker", "link", "fact"},
		"glossary": {}, "guardrails": {"check"}}
	// layupOnly names each case whose EXPECT has no line of a kind that the
	// target keeps: the harness does not compare it (condition 1 of the plan
	// review of #89). The facts cases fail in the engine only because their
	// overlays have no docs/setup/facts.sha256; facts/good-autocrlf is a test
	// of LAYUP's own .gitattributes (#48), not of a rule.
	layupOnly = []string{"facts/bad-answers-blank", "facts/bad-answers-missing", "facts/bad-answers-repeat",
		"facts/bad-batch-absent", "facts/bad-batch-rows", "facts/bad-blank-tab", "facts/good-autocrlf",
		"glossary/bad-no-section", "glossary/bad-rows", "guardrails/bad-no-section"}
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
	listed := append(append(append([]string{"frame"}, notYetBuilt...), notForATarget...), builtGroups...)
	var skipped []string
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
	for _, g := range append([]string{"frame"}, builtGroups...) {
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
				run = builtGroups
			}
			found := false
			for _, check := range run {
				core, ok := built[check]
				if !ok {
					continue
				}
				kinds, target := targetKinds[check]
				keep := func(l string) bool { // a line of a kind that the target keeps, or the OK line
					kind, _, _ := strings.Cut(strings.TrimPrefix(l, "setup-check: "+check+" FAIL "), ":")
					return !target || l == "setup-check: "+check+" OK" || slices.Contains(kinds, kind)
				}
				var want []string
				for _, l := range expect["line"] {
					if strings.HasPrefix(l, "setup-check: "+check+" ") && keep(l) {
						want = append(want, l)
					}
				}
				if target && g != "frame" && (len(want) == 0 || expect["mode"] != nil && expect["mode"][0] == "autocrlf") {
					skipped = append(skipped, name)
					continue
				}
				var findings []string
				for _, f := range core(os.DirFS(repo), gitHistory{repo}) {
					if keep("setup-check: " + check + " FAIL " + f) {
						findings = append(findings, f)
					}
				}
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
			if exit := map[bool]string{false: "0", true: "1"}[found]; g != "frame" && !slices.Contains(skipped, name) && (len(expect["exit"]) != 1 || expect["exit"][0] != exit) {
				t.Errorf("%s: exit %s; want %q", name, exit, expect["exit"])
			}
		}
	}
	if cases == 0 {
		t.Fatal("no case found: a harness that tested nothing is not a pass")
	}
	slices.Sort(skipped)
	if want := slices.Sorted(slices.Values(layupOnly)); !slices.Equal(skipped, want) {
		t.Errorf("the cases that the harness did not compare\n got %q\nwant %q (layupOnly)", skipped, want)
	}
}

// The lists of check adapted equal AD_EXCLUDE and ad_allowed of
// setup-check.sh, read at test time (D1 of #88: the rule stays as it is).
func TestTheListsOfAdaptedEqualTheSh(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "docs", "setup", "setup-check.sh"))
	if err != nil {
		t.Fatal(err)
	}
	var exclude string
	var allowed [][2]string
	in := false
	for _, l := range strings.Split(string(data), "\n") {
		switch {
		case strings.HasPrefix(l, "AD_EXCLUDE='") && strings.HasSuffix(l, "'"):
			exclude = strings.TrimSuffix(strings.TrimPrefix(l, "AD_EXCLUDE='"), "'")
		case strings.TrimSpace(l) == "cat <<'AD_ALLOW'":
			in = true
		case l == "AD_ALLOW":
			in = false
		case in:
			p, reason, _ := strings.Cut(l, "\t")
			allowed = append(allowed, [2]string{p, reason})
		}
	}
	if exclude != adExclude {
		t.Errorf("AD_EXCLUDE of setup-check.sh\n%q\nthe Go copy\n%q", exclude, adExclude)
	}
	if len(allowed) == 0 || !slices.Equal(allowed, adAllowed) {
		t.Errorf("ad_allowed of setup-check.sh\n%q\nthe Go copy\n%q", allowed, adAllowed)
	}
}

// Flagged gives the files of a work tree that check adapted flags, for the
// prose step (O-123, K12): on the fixture repositories of the group adapted,
// and an error where git cannot list the files.
func TestFlagged(t *testing.T) {
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1"} {
		t.Setenv(k, v)
	}
	root := filepath.Join("..", "..", "docs", "setup", "tests")
	for c, want := range map[string][]string{
		"bad-rule-1":   {"docs/a.md", "docs/adr/0013-new.md", "docs/tests/t.md"},
		"bad-rule-2":   {"docs/b.md"},
		"good-allowed": nil,
	} {
		repo := fixtureRepo(t, filepath.Join(root, "kit"), filepath.Join(root, "adapted", c, "overlay"), nil)
		if got, err := Flagged(repo); err != nil || !slices.Equal(got, want) {
			t.Errorf("%s: %q, %v; want %q", c, got, err, want)
		}
	}
	if got, err := Flagged(t.TempDir()); err == nil {
		t.Errorf("a directory that is not a repository: %q, no error", got)
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
