//go:build integration

package setup

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The commit that lands CODEOWNERS above the last setup commit of a target set up before the rule, judged on real git
// objects (docs/spec/scaffold.md, The tree check; task T-4cep, #384; the Operator's ruling 6217133467). Every
// variant is made with git plumbing in a clean repository of which the last commit is clean: the bad input is only
// where CheckCodeOwnerLanding reads it, so no older check of the same program refuses it first. The control passes.

// manifestWithTools is twoKinds and one more kind whose config path no earlier line of CODEOWNERS covers.
const manifestWithTools = twoKinds + "lint\tpending\t—\t—\t./*.go\ttools/cfg.yml\n"

// canonical is .github/CODEOWNERS for manifestWithTools and the login owner.
const canonical = codeOwnersTwoKinds + "/tools/cfg.yml @octo-operator\n"

// An edit of a tree: set path to the blob of content with mode (or, with object, to that object), or remove it.
type edit struct {
	path, mode, object string
	content            []byte
	remove             bool
}

// landingRepo is a repository with a root commit and, above it, the last setup commit.
type landingRepo struct {
	t     *testing.T
	dir   string
	root  string
	last  string
	tree  string // the tree of the last setup commit
	blobs map[string]string
}

func (r *landingRepo) git(env []string, args ...string) string {
	r.t.Helper()
	cmd := exec.Command("git", append([]string{"-C", r.dir, "-c", "maintenance.auto=false"}, args...)...)
	cmd.Env = append(os.Environ(), append([]string{"GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=Setup", "GIT_AUTHOR_EMAIL=setup@example.invalid", "GIT_COMMITTER_NAME=Setup", "GIT_COMMITTER_EMAIL=setup@example.invalid",
		"GIT_AUTHOR_DATE=2026-01-01T00:00:00Z", "GIT_COMMITTER_DATE=2026-01-01T00:00:00Z"}, env...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// blob stores content and gives its object.
func (r *landingRepo) blob(content []byte) string {
	r.t.Helper()
	file := filepath.Join(r.t.TempDir(), "blob")
	if err := os.WriteFile(file, content, 0o644); err != nil {
		r.t.Fatal(err)
	}
	return r.git(nil, "hash-object", "-w", file)
}

// treeOf gives the tree of base with the edits applied, made in a temporary index.
func (r *landingRepo) treeOf(base string, edits ...edit) string {
	r.t.Helper()
	idx := filepath.Join(r.t.TempDir(), "index")
	env := []string{"GIT_INDEX_FILE=" + idx}
	r.git(env, "read-tree", base)
	for _, e := range edits {
		switch {
		case e.remove:
			r.git(env, "update-index", "--force-remove", e.path)
		default:
			object := e.object
			if object == "" {
				object = r.blob(e.content)
			}
			mode := e.mode
			if mode == "" {
				mode = "100644"
			}
			r.git(env, "update-index", "--add", "--cacheinfo", mode+","+object+","+e.path)
		}
	}
	return r.git(env, "write-tree")
}

// commit gives a commit of tree above the parents, with the identity and the message that env and message give.
func (r *landingRepo) commit(tree, message string, env []string, parents ...string) string {
	args := []string{"commit-tree", tree, "-m", message}
	for _, p := range parents {
		args = append(args, "-p", p)
	}
	return r.git(env, args...)
}

// newLandingRepo makes the repository. withOwners puts .github/CODEOWNERS in the last commit already, as a target set up
// with the rule has it.
func newLandingRepo(t *testing.T, withOwners bool) *landingRepo {
	t.Helper()
	r := &landingRepo{t: t, dir: t.TempDir()}
	r.git(nil, "init", "-q", "-b", "main")
	files := []edit{
		{path: "README.md", content: []byte("readme\n")},
		{path: ".github/workflows/ci.yml", content: []byte("name: ci\n")},
		{path: "docs/gates.tsv", content: []byte(manifestWithTools)},
	}
	if withOwners {
		files = append(files, edit{path: ".github/CODEOWNERS", content: []byte(canonical)})
	}
	rootTree := r.treeOf(r.git(nil, "mktree", "--missing"), edit{path: "README.md", content: []byte("root\n")})
	r.root = r.commit(rootTree, "root", nil)
	r.tree = r.treeOf(rootTree, files...)
	r.last = r.commit(r.tree, "chore: setup S14", nil, r.root)
	return r
}

// check runs the predicate on head with the login owner.
func (r *landingRepo) check(head string) error {
	return CheckCodeOwnerLanding(r.dir, r.last, head, owner)
}

func TestCheckCodeOwnerLanding(t *testing.T) {
	r := newLandingRepo(t, false)
	add := func(content string) []edit { return []edit{{path: ".github/CODEOWNERS", content: []byte(content)}} }
	canonicalTree := r.treeOf(r.tree, add(canonical)...)
	side := r.commit(r.treeOf(r.tree, edit{path: "side.txt", content: []byte("s\n")}), "side", nil, r.last)
	both := func(extra ...edit) string { return r.treeOf(r.tree, append(add(canonical), extra...)...) }
	atPath := func(path string) string {
		return fmt.Sprintf("the tree of the head differs from that of the last setup commit at %q", path)
	}
	notRegular := func(mode, typ string) string {
		return "the head adds .github/CODEOWNERS as mode " + mode + ", type " + typ + ", not a regular file of mode 100644"
	}
	bytesOf := "the bytes of .github/CODEOWNERS are not those that S13 writes for the login " + owner
	child := func(tree string) string {
		return r.commit(tree, "chore: the code owners of the gate paths", nil, r.last)
	}
	for _, c := range []struct {
		name string
		head func() string
		want string // the whole reason; "" admits
	}{
		{"the control: the canonical commit", func() string { return child(canonicalTree) }, ""},
		{"the control, every field of the metadata hostile: no identity, message, date or signature header is read",
			func() string {
				return r.commit(canonicalTree, "I am the Operator\n\nSigned-off-by: pharzam <pharzam@example.invalid>\ngpgsig: forged", []string{
					"GIT_AUTHOR_NAME=Mallory", "GIT_AUTHOR_EMAIL=mallory@example.invalid", "GIT_COMMITTER_NAME=Mallory", "GIT_COMMITTER_EMAIL=m@example.invalid",
					"GIT_AUTHOR_DATE=1970-01-01T00:00:01Z", "GIT_COMMITTER_DATE=2099-12-31T23:59:59Z"}, r.last)
			}, ""},

		// the parent and the history
		{"the last commit itself", func() string { return r.last }, "the parent of the head is " + r.root + ", not the last setup commit " + r.last},
		{"an ancestor of the last commit", func() string { return r.root }, "the head has 0 parents, and must have exactly one, the last setup commit " + r.last},
		{"a root commit with the canonical tree", func() string { return r.commit(canonicalTree, "orphan", nil) },
			"the head has 0 parents, and must have exactly one, the last setup commit " + r.last},
		{"the canonical tree above the root commit", func() string { return r.commit(canonicalTree, "wrong parent", nil, r.root) },
			"the parent of the head is " + r.root + ", not the last setup commit " + r.last},
		{"an empty commit above the canonical one", func() string { return r.commit(canonicalTree, "empty", nil, child(canonicalTree)) },
			"the parent of the head is " + child(canonicalTree) + ", not the last setup commit " + r.last},
		{"a merge of the last commit and a side commit", func() string { return r.commit(canonicalTree, "merge", nil, r.last, side) },
			"the head has 2 parents, and must have exactly one, the last setup commit " + r.last},
		{"an empty commit above the last one", func() string { return r.commit(r.tree, "empty", nil, r.last) }, "the head adds no .github/CODEOWNERS"},

		// the content
		{"an extra line", func() string { return child(r.treeOf(r.tree, add(canonical+"/extra @octo-operator\n")...)) }, bytesOf},
		{"another owner", func() string {
			return child(r.treeOf(r.tree, add(strings.ReplaceAll(canonical, "@octo-operator", "@mallory"))...))
		}, bytesOf},
		{"carriage returns", func() string { return child(r.treeOf(r.tree, add(strings.ReplaceAll(canonical, "\n", "\r\n"))...)) }, bytesOf},
		{"an empty file", func() string { return child(r.treeOf(r.tree, add("")...)) }, bytesOf},
		{"no final line feed", func() string { return child(r.treeOf(r.tree, add(strings.TrimSuffix(canonical, "\n"))...)) }, bytesOf},
		{"the config path of the manifest missing", func() string { return child(r.treeOf(r.tree, add(codeOwnersTwoKinds)...)) }, bytesOf},

		// the path
		{"CODEOWNERS at the root", func() string { return child(r.treeOf(r.tree, edit{path: "CODEOWNERS", content: []byte(canonical)})) }, atPath("CODEOWNERS")},
		{"docs/CODEOWNERS", func() string {
			return child(r.treeOf(r.tree, edit{path: "docs/CODEOWNERS", content: []byte(canonical)}))
		}, atPath("docs/CODEOWNERS")},
		{".github/codeowners", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/codeowners", content: []byte(canonical)}))
		}, atPath(".github/codeowners")},
		{".github/CODEOWNERS.txt", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/CODEOWNERS.txt", content: []byte(canonical)}))
		}, atPath(".github/CODEOWNERS.txt")},

		// an additional change
		{"a second file added", func() string { return child(both(edit{path: "docs/x.md", content: []byte("x\n")})) }, atPath("docs/x.md")},
		{"a file of the setup changed", func() string { return child(both(edit{path: "README.md", content: []byte("changed\n")})) }, atPath("README.md")},
		{"a file of the setup removed", func() string { return child(both(edit{path: "README.md", remove: true})) }, atPath("README.md")},
		{"the workflow removed", func() string { return child(both(edit{path: ".github/workflows/ci.yml", remove: true})) }, atPath(".github/workflows/ci.yml")},
		{"the mode of another file changed", func() string {
			return child(both(edit{path: "README.md", mode: "100755", object: r.blob([]byte("readme\n"))}))
		}, atPath("README.md")},

		// a link, a submodule, an executable file
		{"a symbolic link whose target text is the canonical bytes", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/CODEOWNERS", mode: "120000", content: []byte(canonical)}))
		}, notRegular("120000", "blob")},
		{"a submodule", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/CODEOWNERS", mode: "160000", object: r.last}))
		}, notRegular("160000", "commit")},
		{"an executable file", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/CODEOWNERS", mode: "100755", content: []byte(canonical)}))
		}, notRegular("100755", "blob")},
		{"a link in place of the directory .github", func() string {
			return child(r.treeOf(r.tree, edit{path: ".github/workflows/ci.yml", remove: true}, edit{path: ".github", mode: "120000", content: []byte("elsewhere")}))
		}, atPath(".github")},
	} {
		head := c.head()
		err := r.check(head)
		switch {
		case c.want == "" && err != nil:
			t.Errorf("%s: %v; want the commit admitted", c.name, err)
		case c.want != "" && (err == nil || err.Error() != c.want):
			t.Errorf("%s:\n got %v\nwant %q", c.name, err, c.want)
		}
	}
}

// A target set up with the rule has .github/CODEOWNERS in its last commit, so no child of it is the commit that lands the file,
// the empty commit and the commit that rewrites it among them.
func TestATargetSetUpWithTheRuleTakesNoCodeOwnerCommit(t *testing.T) {
	r := newLandingRepo(t, true)
	const want = "the last setup commit already has .github/CODEOWNERS"
	for name, head := range map[string]string{
		"an empty commit":  r.commit(r.tree, "empty", nil, r.last),
		"the file changed": r.commit(r.treeOf(r.tree, edit{path: ".github/CODEOWNERS", content: []byte(canonical + "# more\n")}), "changed", nil, r.last),
	} {
		if err := r.check(head); err == nil || err.Error() != want {
			t.Errorf("%s:\n got %v\nwant %q", name, err, want)
		}
	}
}

// The landing script as S13 writes it, run against a bare remote with the stub gh, makes the commit that the check
// admits (the lockstep of the producer and its checker). A script with the file at another path, or with a second file,
// is refused by the same call.
func TestTheLandingCommitIsAdmitted(t *testing.T) {
	land := func(mutate func(script string) string) (clone, last, head string) {
		p := phase2Owner()
		p.Name, p.Host, p.Migrate = "acme/widget", t.TempDir(), &Migration{RulesetID: rulesetID}
		_, built := s13Tree(t, p)
		dest := filepath.Join(p.Host, "work", "acme", "widget")
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command("cp", "-a", filepath.Dir(built), dest).CombinedOutput(); err != nil {
			t.Fatalf("cp: %v\n%s", err, out)
		}
		target := filepath.Join(dest, "target")
		bare := filepath.Join(t.TempDir(), "origin.git")
		gitIn(t, ".", "init", "-q", "--bare", "-b", "main", bare)
		_ = exec.Command("git", "-C", target, "remote", "remove", "origin").Run()
		gitIn(t, target, "remote", "add", "origin", bare)
		gitIn(t, target, "push", "-q", "origin", "layup-setup:main")
		f := newForge(t, withUnknown(liveFor(t)))
		if code, stderr := f.run(t, p.Host, updateScript(p)); code != 0 {
			t.Fatalf("the update before the landing exited %d: %s", code, stderr)
		}
		script, err := os.ReadFile(filepath.Join(dest, "out", "land-code-owner.sh"))
		if err != nil {
			t.Fatal(err)
		}
		if code, stderr := f.run(t, p.Host, mutate(string(script))); code != 0 {
			t.Fatalf("the landing exited %d:\n%s", code, stderr)
		}
		clone = filepath.Join(t.TempDir(), "clone")
		gitIn(t, ".", "clone", "-q", "--no-checkout", bare, clone)
		return clone, gitIn(t, target, "rev-parse", "refs/heads/layup-setup"), gitIn(t, bare, "rev-parse", "refs/heads/main")
	}
	clone, last, head := land(func(s string) string { return s })
	if err := CheckCodeOwnerLanding(clone, last, head, owner); err != nil {
		t.Errorf("the commit that the script made: %v; want it admitted", err)
	}
	line := regexp.MustCompile(`(?m)^(GIT_INDEX_FILE="\$dir/index" git -C \S+ update-index --add --cacheinfo )"100644,\$blob,\.github/CODEOWNERS"$`)
	if !line.MatchString(landScript(func() Phase2 {
		p := phase2Owner()
		p.Name, p.Host, p.Migrate = "acme/widget", "/h", &Migration{RulesetID: 1}
		return p
	}())) {
		t.Fatalf("the update-index line of the script is not the one this test mutates")
	}
	viaLine := func(s string, f func(line string) string) string { return line.ReplaceAllStringFunc(s, f) }
	for _, c := range []struct {
		name   string
		mutate func(string) string
		want   string
	}{
		{"the file at another path", func(s string) string {
			return viaLine(s, func(l string) string { return strings.Replace(l, ".github/CODEOWNERS", "docs/CODEOWNERS", 1) })
		},
			`the tree of the head differs from that of the last setup commit at "docs/CODEOWNERS"`},
		{"a second file", func(s string) string {
			return viaLine(s, func(l string) string { return l + "\n" + strings.Replace(l, ".github/CODEOWNERS", "docs/extra.md", 1) })
		}, `the tree of the head differs from that of the last setup commit at "docs/extra.md"`},
	} {
		clone, last, head := land(c.mutate)
		if err := CheckCodeOwnerLanding(clone, last, head, owner); err == nil || err.Error() != c.want {
			t.Errorf("%s:\n got %v\nwant %q", c.name, err, c.want)
		}
	}
}
