package setup

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
)

// The code-owner rule of S13 in phase 2 (docs/spec/scaffold.md, The code-owner rule; docs/spec/loop.md, The
// gate jobs read the base, part 2; task T-jew3, #356, plan review 6214007155). Each case breaks one rule only,
// and each refusal is compared with its whole reason.

// owner is the login of the Operator's row of approvers.tsv in these tests.
const owner = "octo-operator"

// codeOwnersTwoKinds is .github/CODEOWNERS for the manifest twoKinds, whose one config path lies under /docs/gates/.
const codeOwnersTwoKinds = "# Written by layup setup (S13): a change to a gate path waits for the Operator's review (docs/spec/loop.md).\n" +
	"/.github/ @octo-operator\n/docs/gates.tsv @octo-operator\n/docs/gates/ @octo-operator\n"

// The four classes of paths, root-anchored, and a config path that a line already covers (by the directory prefix
// at a slash, never by a shared name) or that is listed twice has no line of its own; the order is the manifest's.
func TestCodeOwnersText(t *testing.T) {
	for _, c := range []struct {
		name    string
		configs []string
		extra   string
	}{
		{"no config path", nil, ""},
		{"the config path of the Go entry, under /docs/gates/", []string{"docs/gates/coverage-floor.txt"}, ""},
		{"a config path that is the manifest itself", []string{"docs/gates.tsv"}, ""},
		{"an uncovered path, a duplicate, one under it, and a name that only shares a prefix", []string{"tools/cfg.yml", "tools/cfg.yml", "tools/cfg.yml/sub", "docs/gatesx/a"},
			"/tools/cfg.yml @octo-operator\n/docs/gatesx/a @octo-operator\n"},
	} {
		got, err := CodeOwners(owner, c.configs)
		if want := codeOwnersTwoKinds + c.extra; err != nil || string(got) != want {
			t.Errorf("%s (err %v):\n%s\nwant\n%s", c.name, err, got, want)
		}
	}
}

// A login or a config path that CODEOWNERS cannot hold as one pattern is a refusal that names it, with no text:
// no guess, no quoting. A login is at most 39 characters, as the owner of a repository is (ownerName).
func TestCodeOwnersRefusals(t *testing.T) {
	login := func(s string) string {
		return "CODEOWNERS: the login " + strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll("\""+s+"\"", "\n", `\n`), "\t", `\t`), "\r", `\r`) + " is not a GitHub login"
	}
	path := func(s string) string {
		return "CODEOWNERS: the config path \"" + s + "\" cannot be written as a root-anchored pattern"
	}
	long := strings.Repeat("a", 40)
	for _, c := range []struct {
		name, login string
		configs     []string
		want        string
	}{
		{"an empty login", "", nil, login("")},
		{"a login with a space", "a b", nil, login("a b")},
		{"a login with @", "@x", nil, login("@x")},
		{"a login with a line feed", "x\ny", nil, login("x\ny")},
		{"a login that begins with a hyphen", "-x", nil, login("-x")},
		{"a login of 40 characters", long, nil, login(long)},
		{"a config path with a space", owner, []string{"a b"}, path("a b")},
		{"a config path with #", owner, []string{"a#b"}, path("a#b")},
		{"a config path with !", owner, []string{"!x"}, path("!x")},
		{"a config path with [", owner, []string{"[x]"}, path("[x]")},
		{"a config path with *", owner, []string{"a/*.yml"}, path("a/*.yml")},
		{"a config path with ..", owner, []string{"a/../b"}, path("a/../b")},
		{"a config path with an empty segment", owner, []string{"a//b"}, path("a//b")},
		{"an absolute config path", owner, []string{"/abs"}, path("/abs")},
		{"a config path that ends in a slash", owner, []string{"a/"}, path("a/")},
		{"a config path that is a dot", owner, []string{"."}, path(".")},
	} {
		if got, err := CodeOwners(c.login, c.configs); err == nil || err.Error() != c.want || got != nil {
			t.Errorf("%s: %q, %v; want the error %q and no text", c.name, got, err, c.want)
		}
	}
	if !ownerLogin(strings.Repeat("a", 39)) || ownerLogin(long) || !ownerLogin("a-b") || ownerLogin("a--b") {
		t.Errorf("ownerLogin: 39 characters %v, 40 %v, a-b %v, a--b %v; want true, false, true, false", ownerLogin(strings.Repeat("a", 39)), ownerLogin(long), ownerLogin("a-b"), ownerLogin("a--b"))
	}
}

// phase2Owner gives the inputs of phase 2 with the Operator's login and the App 42.
func phase2Owner() Phase2 {
	p := phase2App(42)
	p.CodeOwner = owner
	return p
}

// A new target carries .github/CODEOWNERS in its tree, in the commit of S13, with a row of its hash; the
// protection file, the ruleset files and the commands are as before. A login that is not one is a fail of S13
// with no file written.
func TestS13WritesCodeOwnersInTheTree(t *testing.T) {
	f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	f.install(t)
	o := StepsPhase2(phase2Owner(), Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()})
	row := []string{"codeowners.sha256", sha([]byte(codeOwnersTwoKinds)), "computed", "sha256 .github/CODEOWNERS"}
	if o.Kind != Done || !o.Commit || len(o.Values) != 4 || !reflect.DeepEqual(o.Values[2], row) {
		t.Errorf("S13: %s %q %v, values %q; want done, committed, and the row %q third", o.Kind, o.Evidence, o.Commit, o.Values, row)
	}
	if got := f.files["w/target/.github/CODEOWNERS"]; got != codeOwnersTwoKinds {
		t.Errorf("CODEOWNERS in the tree:\n%s\nwant\n%s", got, codeOwnersTwoKinds)
	}
	if got := f.files["w/target/docs/setup/branch-protection.json"]; got != protectionText {
		t.Errorf("branch-protection.json changed:\n%s", got)
	}
	if _, ok := f.files["w/out/CODEOWNERS"]; ok || len(f.files) != 4 {
		t.Errorf("S13 wrote %d files, out/CODEOWNERS %v; want the four of a new target", len(f.files), ok)
	}
	bad := phase2Owner()
	bad.CodeOwner = "a b"
	g := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	g.install(t)
	want := "S13: " + `CODEOWNERS: the login "a b" is not a GitHub login`
	if o := StepsPhase2(bad, Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()}); o.Kind != Fail || o.Evidence != want || len(g.files) != 0 {
		t.Errorf("a bad login: %s %q, %d files; want fail %q and no file", o.Kind, o.Evidence, len(g.files), want)
	}
}

// A target set up before the rule is rebuilt with the files of M2d in its tree (the commits of the rebuilt area
// must be the recorded ones, scaffold.md, The work area): the protection file alone, no CODEOWNERS in the tree and
// no row for it. The file for the landing is in out/, with the new ruleset file, the landing script (not
// applicable) and the update in place of the push and the POST of commands.sh.
func TestS13InMigrateModeWritesTheFilesOfM2d(t *testing.T) {
	p := phase2Owner()
	p.Migrate = &Migration{RulesetID: 25145510}
	f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	f.install(t)
	steps := StepsPhase2(p, Calls{})
	o := steps["S13"].Run(Input{Dir: "w", Record: lastRecord()})
	if o.Kind != Done || !o.Commit || len(o.Values) != 3 || o.Values[1][0] != "ruleset.sha256" || o.Values[2][0] != "without-layup" {
		t.Errorf("S13: %s %q %v, values %q; want the three rows of M2d", o.Kind, o.Evidence, o.Commit, o.Values)
	}
	for path, want := range map[string]string{"w/target/docs/setup/branch-protection.json": protectionText, "w/out/ruleset-default.json": rulesetPhase2Text,
		"w/out/CODEOWNERS": codeOwnersTwoKinds, "w/out/land-code-owner.sh": landScript(p, false), "w/out/ruleset-records.json": rulesetRecordsText} {
		if got, ok := f.files[path]; !ok || got != want {
			t.Errorf("%s:\n%s\nwant\n%s", path, got, want)
		}
	}
	if _, ok := f.files["w/target/.github/CODEOWNERS"]; ok || len(f.files) != 5 {
		t.Errorf("S13 wrote %d files, CODEOWNERS in the tree %v; want the five above and none in the tree", len(f.files), ok)
	}
	got := steps["S13"].Commands(s13Record())
	if want := updateCommands(p); !reflect.DeepEqual(got, want) || len(got) == 0 {
		t.Errorf("the commands %+v; want the update %+v", got, want)
	}
	var all string
	for _, c := range got {
		all += c.Text + "\n"
	}
	if strings.Contains(all, "POST") || strings.Contains(all, "git push") || !strings.Contains(all, "rulesets/25145510") {
		t.Errorf("the update has a POST, a push or no recorded ID:\n%s", all)
	}
	for _, bad := range []int64{0, -3} {
		p.Migrate = &Migration{RulesetID: bad}
		if o := StepsPhase2(p, Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()}); o.Kind != Fail || !strings.Contains(o.Evidence, "ruleset ID") {
			t.Errorf("a ruleset ID of %d: %s %q; want a fail that names the ID", bad, o.Kind, o.Evidence)
		}
	}
}

// liveRuleset gives the ruleset that the apply of the M2d demo returned (runs/T-sdtb/rulesets.txt, byte for byte).
func liveRuleset(t *testing.T) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", "..", "runs", "T-sdtb", "rulesets.txt"))
	if err != nil {
		t.Fatal(err)
	}
	_, rest, _ := strings.Cut(string(data), "== default.json\n")
	text, _, _ := strings.Cut(rest, "\n== records.json")
	var m map[string]any
	if err := json.Unmarshal([]byte(text), &m); err != nil {
		t.Fatal(err)
	}
	return m
}

// The fake of the forge takes its form from a first-hand answer (guardrails.md §2, a fake that echoes the file): the
// object that the apply of the M2d demo returned holds fields that the file of S13 never wrote, so a body made from
// the file is not the live ruleset, and the tests of the update serve this object, not the file.
func TestTheStubFormIsTheLiveForm(t *testing.T) {
	live := liveRuleset(t)
	for _, k := range []string{"id", "source_type", "source", "node_id", "created_at", "updated_at", "current_user_can_bypass", "_links"} {
		if _, ok := live[k]; !ok {
			t.Errorf("the live object has no %s", k)
		}
	}
	rules := rulesOf(live)
	if _, ok := rules[3]["parameters"].(map[string]any)["required_reviewers"]; !ok {
		t.Errorf("the pull_request rule of the live object has no required_reviewers: %v", rules[3])
	}
	if _, ok := rules[4]["parameters"].(map[string]any)["do_not_enforce_on_create"]; !ok {
		t.Errorf("the required_status_checks rule of the live object has no do_not_enforce_on_create: %v", rules[4])
	}
	file, err := rulesetBody([]string{"static", "layout", "boundary", "contract", "test"}, 5118100)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := rulesOf(decode(t, file))[3]["parameters"].(map[string]any)["required_reviewers"]; ok {
		t.Errorf("the file of S13 holds required_reviewers: the difference that the update must keep has gone")
	}
}

// The commit that lands CODEOWNERS on a target set up before the rule (docs/spec/scaffold.md, The tree check; task
// T-4cep, #384; the Operator's ruling 6217133467). classifyLanding judges the parents and the two listings of whole
// trees, and gives the first condition that fails, as a whole sentence, or "" for the one commit that is admitted.
const (
	lastSetup   = "1111111111111111111111111111111111111111"
	strayParent = "2222222222222222222222222222222222222222"
	ownersBlob  = "7efc45d58203810bb12f9115bc722ce5fec65a58" // git hash-object of codeOwnersTwoKinds
	strayBlob   = "0000000000000000000000000000000000000001"
)

func entry(mode, typ, object, path string) git.TreeEntry {
	return git.TreeEntry{Mode: mode, Type: typ, Object: object, Path: path}
}

// lastTree is the tree of the last setup commit of these cases: files of three depths, none of them CODEOWNERS.
func lastTree() []git.TreeEntry {
	return []git.TreeEntry{
		entry("100644", "blob", "aaaa", ".github/workflows/ci.yml"),
		entry("100644", "blob", "bbbb", "README.md"),
		entry("100644", "blob", "cccc", "docs/gates.tsv"),
	}
}

// withEntries gives the tree of lastTree with the entries of change applied: a path with a nil entry is removed.
func withEntries(change map[string]*git.TreeEntry) []git.TreeEntry {
	var out []git.TreeEntry
	for _, e := range lastTree() {
		if c, ok := change[e.Path]; ok {
			if c == nil {
				continue
			}
			e = *c
		}
		out = append(out, e)
	}
	for p, c := range change {
		if c != nil && !slices.ContainsFunc(lastTree(), func(e git.TreeEntry) bool { return e.Path == p }) {
			out = append(out, *c)
		}
	}
	return out
}

func ptr(e git.TreeEntry) *git.TreeEntry { return &e }

func TestClassifyLanding(t *testing.T) {
	good := ptr(entry("100644", "blob", ownersBlob, ".github/CODEOWNERS"))
	for _, c := range []struct {
		name    string
		parents []string
		last    []git.TreeEntry
		head    []git.TreeEntry
		want    string
	}{
		{"the one commit that is admitted", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}), ""},
		{"the same, the listing in another order", []string{lastSetup}, lastTree(), slices.Concat([]git.TreeEntry{*good}, lastTree()), ""},
		{"no parent: a root commit", nil, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}),
			"the head has 0 parents, and must have exactly one, the last setup commit " + lastSetup},
		{"two parents: a merge", []string{lastSetup, strayParent}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}),
			"the head has 2 parents, and must have exactly one, the last setup commit " + lastSetup},
		{"another parent: a grandchild, an ancestor, an unrelated commit", []string{strayParent}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}),
			"the parent of the head is " + strayParent + ", not the last setup commit " + lastSetup},
		{"the last commit already has CODEOWNERS", []string{lastSetup},
			withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good}),
			"the last setup commit already has .github/CODEOWNERS"},
		{"nothing added: an empty commit", []string{lastSetup}, lastTree(), lastTree(), "the head adds no .github/CODEOWNERS"},
		{"another path: CODEOWNERS at the root", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{"CODEOWNERS": ptr(entry("100644", "blob", ownersBlob, "CODEOWNERS"))}),
			`the tree of the head differs from that of the last setup commit at "CODEOWNERS"`},
		{"another path: docs/CODEOWNERS", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{"docs/CODEOWNERS": ptr(entry("100644", "blob", ownersBlob, "docs/CODEOWNERS"))}),
			`the tree of the head differs from that of the last setup commit at "docs/CODEOWNERS"`},
		{"another path: the name in lower case", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/codeowners": ptr(entry("100644", "blob", ownersBlob, ".github/codeowners"))}),
			`the tree of the head differs from that of the last setup commit at ".github/codeowners"`},
		{"a second file added", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good, "docs/x.md": ptr(entry("100644", "blob", "dddd", "docs/x.md"))}),
			`the tree of the head differs from that of the last setup commit at "docs/x.md"`},
		{"a file changed", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good, "README.md": ptr(entry("100644", "blob", "eeee", "README.md"))}),
			`the tree of the head differs from that of the last setup commit at "README.md"`},
		{"a file removed", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good, "README.md": nil}),
			`the tree of the head differs from that of the last setup commit at "README.md"`},
		{"a mode changed elsewhere", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": good, "README.md": ptr(entry("100755", "blob", "bbbb", "README.md"))}),
			`the tree of the head differs from that of the last setup commit at "README.md"`},
		{"a symbolic link", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": ptr(entry("120000", "blob", ownersBlob, ".github/CODEOWNERS"))}),
			"the head adds .github/CODEOWNERS as mode 120000, type blob, not a regular file of mode 100644"},
		{"a submodule", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": ptr(entry("160000", "commit", ownersBlob, ".github/CODEOWNERS"))}),
			"the head adds .github/CODEOWNERS as mode 160000, type commit, not a regular file of mode 100644"},
		{"an executable file", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": ptr(entry("100755", "blob", ownersBlob, ".github/CODEOWNERS"))}),
			"the head adds .github/CODEOWNERS as mode 100755, type blob, not a regular file of mode 100644"},
		{"other bytes", []string{lastSetup}, lastTree(), withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": ptr(entry("100644", "blob", strayBlob, ".github/CODEOWNERS"))}),
			"the bytes of .github/CODEOWNERS are not those that S13 writes for the login " + owner},
	} {
		if got := classifyLanding(c.parents, lastSetup, c.last, c.head, []byte(codeOwnersTwoKinds), owner); got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// CheckCodeOwnerLanding reads the parents of the head, the manifest and the whole tree of the last commit, and the
// whole tree of the head, and nothing else: no author, committer, message, date, signature or actor is read, so no
// identity is inferred (condition 3 of the plan review 6217115601). A head whose parent is not the last commit
// is judged before any tree is read.
func TestCheckCodeOwnerLandingReadsOnlyTheObjectsOfTheConditions(t *testing.T) {
	const head = "3333333333333333333333333333333333333333"
	saved := sys
	t.Cleanup(func() { sys = saved })
	var calls []string
	parents := []string{lastSetup}
	sys.parents = func(dir, rev string) ([]string, error) { calls = append(calls, "parents "+rev); return parents, nil }
	sys.lsTree = func(dir, rev, path string) ([]git.TreeEntry, error) {
		calls = append(calls, "ls-tree "+rev+" "+path)
		if rev == head {
			return withEntries(map[string]*git.TreeEntry{".github/CODEOWNERS": ptr(entry("100644", "blob", ownersBlob, ".github/CODEOWNERS"))}), nil
		}
		return lastTree(), nil
	}
	sys.show = func(dir, rev, path string) ([]byte, error) {
		calls = append(calls, "show "+rev+" "+path)
		return []byte(twoKinds), nil
	}
	if err := CheckCodeOwnerLanding("d", lastSetup, head, owner); err != nil {
		t.Errorf("the canonical head: %v; want nil", err)
	}
	want := []string{"parents " + head, "show " + lastSetup + " docs/gates.tsv", "ls-tree " + lastSetup + " .", "ls-tree " + head + " ."}
	if !reflect.DeepEqual(calls, want) {
		t.Errorf("the reads %q; want %q (the manifest at the last commit, never at the head)", calls, want)
	}
	calls, parents = nil, []string{strayParent}
	err := CheckCodeOwnerLanding("d", lastSetup, head, owner)
	if wantErr := "the parent of the head is " + strayParent + ", not the last setup commit " + lastSetup; err == nil || err.Error() != wantErr {
		t.Errorf("a wrong parent: %v; want %q", err, wantErr)
	}
	if want := []string{"parents " + head}; !reflect.DeepEqual(calls, want) {
		t.Errorf("a wrong parent read %q; want %q, no tree", calls, want)
	}
	// a login that CODEOWNERS cannot hold is an error that names it, not a pass
	parents = []string{lastSetup}
	if err := CheckCodeOwnerLanding("d", lastSetup, head, "a b"); err == nil || !strings.Contains(err.Error(), `the login "a b"`) {
		t.Errorf("a login with a space: %v; want an error that names it", err)
	}
	// a read that fails is an error of the call, never an admission
	sys.lsTree = func(dir, rev, path string) ([]git.TreeEntry, error) { return nil, errors.New("boom") }
	if err := CheckCodeOwnerLanding("d", lastSetup, head, owner); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Errorf("a tree that cannot be read: %v; want an error with the cause", err)
	}
}
