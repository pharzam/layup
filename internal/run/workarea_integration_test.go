//go:build integration

package run

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	workpkg "github.com/pharzam/layup/internal/work"
)

// The work area of a target under layup run (docs/spec/scaffold.md, The work
// area; row 54 of the plan, task T-eahr, #277) on a real repository: removed
// and built again from the records, with target/ a clone of the default branch
// from the run's clone, main at the root commit, and origin the forge's URL.

// forgeRepo makes a bare repository on branch with one commit for each set of
// files, and gives its path and the IDs of its commits; the run's clone is made
// of it as the step clone makes it.
func forgeRepo(t *testing.T, branch string, sets ...map[string]string) (bare string, ids []string) {
	t.Helper()
	src := t.TempDir()
	gitOut(t, src, "init", "-q", "-b", branch)
	for i, files := range sets {
		for path, text := range files {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(src, path)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(src, path), []byte(text), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		gitOut(t, src, "add", "-A")
		gitOut(t, src, "commit", "-q", "-m", fmt.Sprintf("commit %d", i+1))
		ids = append(ids, strings.TrimSpace(gitOut(t, src, "rev-parse", "HEAD")))
	}
	bare = filepath.Join(t.TempDir(), "forge.git")
	gitOut(t, "/", "clone", "-q", "--bare", src, bare)
	return bare, ids
}

// runClone is the run's clone of the bare repository: no checkout.
func runClone(t *testing.T, bare string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "clone")
	if err := git.Clone(bare, dir, git.Auth{}); err != nil {
		t.Fatal(err)
	}
	return dir
}

// treeRows is setup/tree.tsv for the commits of the repository, the first the
// root: order, commit, tree and step.
func treeRows(t *testing.T, bare string, ids []string) []byte {
	t.Helper()
	var rows [][]string
	for i, id := range ids {
		step := "root"
		if i > 0 {
			step = fmt.Sprintf("S%02d", i+3)
		}
		rows = append(rows, []string{fmt.Sprint(i), id, strings.TrimSpace(gitOut(t, bare, "rev-parse", id+"^{tree}")), step})
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, SetupTreeSchema, rows); err != nil {
		t.Fatal(err)
	}
	return b.Bytes()
}

const originURL = "https://github.invalid/acme/target.git"

var answersRows = [][]string{{"S01-stack", "go", "operator", "https://github.invalid/acme/target/issues/1#issuecomment-5", ""}, {"Q-001", "PostgreSQL", "idea-owner", "https://github.invalid/acme/target/issues/1#issuecomment-6", ""}}

func TestTheWorkAreaIsBuiltFromTheRecords(t *testing.T) {
	bare, ids := forgeRepo(t, "main", map[string]string{"README.md": "# Target\n", "docs/a.md": "a\n"})
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	if err := os.MkdirAll(dir, 0o755); err != nil { // a work area of an earlier run, which is removed
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "stale.txt"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := workSource{Brief: []byte("The problem.\n"), Vision: []byte("The vision.\n"), Answers: answersRows,
		Files: map[string][]byte{"docs/glossary.md": []byte("# Glossary\n"), "README.md": []byte("# Acme\n")}}
	if err := buildWorkArea(dir, src, runClone(t, bare), "main", originURL); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{"inputs/briefs/problem-statement.md": "The problem.\n", "inputs/briefs/vision.md": "The vision.\n",
		"inputs/files/docs/glossary.md": "# Glossary\n", "inputs/files/README.md": "# Acme\n", "target/README.md": "# Target\n", "target/docs/a.md": "a\n"} {
		if got, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(path))); err != nil || string(got) != want {
			t.Errorf("%s: %q, %v; want %q", path, got, err, want)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "stale.txt")); err == nil {
		t.Error("stale.txt of the earlier work area is still there; want the area removed and built again")
	}
	if got, err := workpkg.ReadAnswers(dir); err != nil || !reflect.DeepEqual([][]string(got), answersRows) {
		t.Errorf("inputs/answers.tsv: %q, %v; want the rows given", got, err)
	}
	target := filepath.Join(dir, "target")
	for _, c := range []struct{ name, want string }{
		{"rev-parse HEAD", ids[0]}, {"symbolic-ref --short HEAD", "main"}, {"remote", "origin"}, {"remote get-url origin", originURL}, {"status --porcelain", ""},
		{"for-each-ref --format=%(refname)", "refs/heads/main"},
	} {
		if got := strings.TrimSpace(gitOut(t, target, strings.Fields(c.name)...)); got != c.want {
			t.Errorf("git %s in target/: %q; want %q", c.name, got, c.want)
		}
	}
}

// With no vision brief the area has none; a second build gives the same head and a clean tree.
func TestTheWorkAreaIsBuiltAgainFromScratch(t *testing.T) {
	bare, ids := forgeRepo(t, "main", map[string]string{"README.md": "# Target\n"})
	clone := runClone(t, bare)
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	src := workSource{Brief: []byte("The problem.\n"), Answers: answersRows}
	for i := 0; i < 2; i++ {
		if err := buildWorkArea(dir, src, clone, "main", originURL); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(dir, "inputs", "briefs", "vision.md")); err == nil {
			t.Errorf("build %d: vision.md exists; want none with no vision brief", i+1)
		}
		if got := strings.TrimSpace(gitOut(t, filepath.Join(dir, "target"), "rev-parse", "HEAD")); got != ids[0] {
			t.Errorf("build %d: HEAD %s; want %s", i+1, got, ids[0])
		}
	}
	if err := os.WriteFile(filepath.Join(dir, "target", "dirty.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := buildWorkArea(dir, src, clone, "main", originURL); err != nil {
		t.Fatal(err)
	}
	if got := gitOut(t, filepath.Join(dir, "target"), "status", "--porcelain"); got != "" {
		t.Errorf("the target after a third build: %q; want a clean tree, as a build removes what was there", got)
	}
}

// After the Operator's push the default branch holds the setup commits: target/
// has main at the root commit, the row of order 0 of setup/tree.tsv, and its
// files, not the files of the branch's tip.
func TestTheWorkAreaTargetIsAtTheRootCommitOfTheTree(t *testing.T) {
	bare, ids := forgeRepo(t, "main", map[string]string{"README.md": "# Root\n"}, map[string]string{"docs/pin": "pin\n"}, map[string]string{"docs/glossary.md": "g\n"})
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	src := workSource{Brief: []byte("b\n"), Tree: treeRows(t, bare, ids)}
	if err := buildWorkArea(dir, src, runClone(t, bare), "main", originURL); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "target")
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "main")); got != ids[0] {
		t.Errorf("main of target/ is %s; want the root commit %s", got, ids[0])
	}
	if _, err := os.Stat(filepath.Join(target, "docs", "glossary.md")); err == nil {
		t.Error("docs/glossary.md of the tip is in target/; want the files of the root commit")
	}
	if got := gitOut(t, target, "status", "--porcelain"); got != "" {
		t.Errorf("the tree of target/: %q; want it clean at the root commit", got)
	}
}

// With no setup/tree.tsv the default branch must hold the one commit that Start read back.
func TestTheWorkAreaRefusesMoreThanTheRootCommitWithNoTree(t *testing.T) {
	bare, _ := forgeRepo(t, "main", map[string]string{"README.md": "# Root\n"}, map[string]string{"docs/pin": "pin\n"})
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	err := buildWorkArea(dir, workSource{Brief: []byte("b\n")}, runClone(t, bare), "main", originURL)
	if err == nil || !strings.Contains(err.Error(), "setup/tree.tsv") {
		t.Errorf("a default branch of two commits and no tree: %v; want an error that names setup/tree.tsv", err)
	}
}

// The root commit of the tree must be in the clone.
func TestTheWorkAreaRefusesARootThatTheCloneLacks(t *testing.T) {
	bare, ids := forgeRepo(t, "main", map[string]string{"README.md": "# Root\n"})
	other, otherIDs := forgeRepo(t, "main", map[string]string{"README.md": "# Another\n"})
	_ = ids
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	err := buildWorkArea(dir, workSource{Brief: []byte("b\n"), Tree: treeRows(t, other, otherIDs)}, runClone(t, bare), "main", originURL)
	if err == nil || !strings.Contains(err.Error(), otherIDs[0]) || !strings.Contains(err.Error(), "setup/tree.tsv") {
		t.Errorf("a tree whose root commit is not in the clone: %v; want an error that names %s and setup/tree.tsv", err, otherIDs[0])
	}
}

// The runner reads refs/heads/main (S03, S04): the branch of target/ is main
// whatever the forge's default branch is called.
func TestTheTargetIsOnMainWhateverTheDefaultBranchIs(t *testing.T) {
	bare, ids := forgeRepo(t, "trunk", map[string]string{"README.md": "# Root\n"})
	dir := filepath.Join(t.TempDir(), "work", "acme", "target")
	if err := buildWorkArea(dir, workSource{Brief: []byte("b\n")}, runClone(t, bare), "trunk", originURL); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "target")
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "refs/heads/main")); got != ids[0] {
		t.Errorf("refs/heads/main of target/ is %q; want the root commit %s", got, ids[0])
	}
	if got := strings.TrimSpace(gitOut(t, target, "symbolic-ref", "--short", "HEAD")); got != "main" {
		t.Errorf("HEAD of target/ is on %q; want main", got)
	}
	if got := gitOut(t, target, "status", "--porcelain"); got != "" {
		t.Errorf("the tree of target/: %q; want it clean", got)
	}
}

// The work area refuses, before it touches the disk, a file of the records
// whose path leaves inputs/files (docs/spec/scaffold.md, The work area): the
// records are LAYUP's own, and a path with .. is no path of the target.
func TestTheWorkAreaRefusesAFilePathThatLeavesItsDirectory(t *testing.T) {
	for _, path := range []string{"..", "../escape.md", "docs/../../escape.md", "/etc/passwd", "docs//glossary.md", "", ".", "docs/./glossary.md"} {
		dir := filepath.Join(t.TempDir(), "work", "acme", "target")
		err := buildWorkArea(dir, workSource{Brief: []byte("b"), Files: map[string][]byte{path: []byte("x")}}, "/no/such/clone", "main", originURL)
		if err == nil || !strings.Contains(err.Error(), "setup/inputs/files/") {
			t.Errorf("the path %q: %v; want an error that names the path in the records", path, err)
		}
		if _, statErr := os.Stat(dir); statErr == nil {
			t.Errorf("the path %q: the work area exists after the refusal; want no change of the disk", path)
		}
	}
}

// A link at DIR/work or at DIR/work/OWNER would take the removal of the area outside DIR
// (docs/spec/scaffold.md, The work area): the build is refused, and the files outside stay.
func TestTheWorkAreaIsNotBuiltThroughALink(t *testing.T) {
	bare, _ := forgeRepo(t, "main", map[string]string{"README.md": "# Target\n"})
	clone := runClone(t, bare)
	for _, link := range []string{"work/acme", "work"} {
		host, outside := t.TempDir(), t.TempDir()
		keep := []string{"target/keep.txt", "acme/target/keep.txt"}
		for _, k := range keep {
			if err := os.MkdirAll(filepath.Join(outside, filepath.Dir(k)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(outside, k), []byte("keep\n"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		if err := os.MkdirAll(filepath.Dir(filepath.Join(host, link)), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(host, link)); err != nil {
			t.Fatal(err)
		}
		err := buildWorkArea(filepath.Join(host, "work", "acme", "target"), workSource{Brief: []byte("b\n")}, clone, "main", originURL)
		if err == nil || !strings.Contains(err.Error(), "is a link") {
			t.Errorf("a link at %s: %v; want an error that names the link", link, err)
		}
		for _, k := range keep {
			if data, err := os.ReadFile(filepath.Join(outside, k)); err != nil || string(data) != "keep\n" {
				t.Errorf("a link at %s: %s outside DIR is %q, %v; want it kept", link, k, data, err)
			}
		}
	}
}
