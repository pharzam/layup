package setup

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/catalog"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// goEntry gives the Go entry of the catalog of the binary.
func goEntry(t *testing.T) *catalog.Entry {
	t.Helper()
	e, err := catalog.Read(catalog.Embedded(), "go")
	if err != nil {
		t.Fatal(err)
	}
	return e
}

// lastRecord gives a record whose S01 has the stack go and the name
// acme/widget, with the pin rows of S02.
func lastRecord() work.Record {
	return append(work.Record{{"S01", "stack", "go", "answer", "S01-stack"}}, pinRecord()...)
}

// S12 writes the files of the entry with the module path of S01-name, the
// manifest, and after the rows that the head has, the row of open-gaps.tsv of
// each gap; its value rows name each file and gap (D1 of #92).
func TestS12(t *testing.T) {
	e := goEntry(t)
	gap := e.Gaps()[0]
	files, err := e.Files("github.com/acme/widget")
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := e.Manifest()
	if err != nil {
		t.Fatal(err)
	}
	gapRow := gap.Path + "\t" + gap.Marker + "\t" + gap.Question + "\n"
	f := &fakeRepo{trees: map[string][]git.TreeEntry{".": {{Path: "README.md"}, {Path: work.OpenGapsPath}}},
		shows: map[string]string{work.OpenGapsPath: "docs/a.md\tM\tWhich a?"}}
	f.install(t)
	o := runS12(Input{Dir: "w", Record: lastRecord()})
	want := map[string]string{"w/target/docs/gates.tsv": string(manifest), "w/target/" + work.OpenGapsPath: "docs/a.md\tM\tWhich a?\n" + gapRow}
	values := [][]string{{"module", "github.com/acme/widget", "answer", "S01-name"}}
	for _, file := range files {
		want["w/target/"+file.Path] = string(file.Data)
		values = append(values, []string{"catalog:" + file.Path, sha(file.Data), "catalog", "go/files/" + file.Path + ".tmpl"})
	}
	values = append(values, []string{"catalog:docs/gates.tsv", sha(manifest), "catalog", "go/kinds.tsv"},
		[]string{fmt.Sprintf("marker:%s:%d", gap.Path, gap.Line), gap.Marker, "gap", work.OpenGapsPath})
	if o.Kind != Done || o.Evidence != "checks jobs and gates" || !o.Commit || !reflect.DeepEqual(o.Values, values) {
		t.Errorf("S12: %s %q, values %q; want done, and the values %q", o.Kind, o.Evidence, o.Values, values)
	}
	for _, p := range slices.Sorted(mapKeys(want, f.files)) {
		if f.files[p] != want[p] {
			t.Errorf("S12 wrote %s:\n%q\nwant\n%q", p, f.files[p], want[p])
		}
	}
	if f.files["w/target/go.mod"] != "module github.com/acme/widget\n\ngo 1.26\n" || !strings.Contains(f.files["w/target/docs/gates.tsv"], "\ntest\tactive\tgo\t") {
		t.Errorf("go.mod %q and the manifest %q; want the module path of S01-name, and the kinds of the entry", f.files["w/target/go.mod"], f.files["w/target/docs/gates.tsv"])
	}

	// The head of a run that stopped after the commit of S12 holds the same
	// bytes and the gap row: no second row.
	head := []git.TreeEntry{{Path: "go.mod"}, {Path: work.OpenGapsPath}}
	f = &fakeRepo{trees: map[string][]git.TreeEntry{".": head}, shows: map[string]string{"go.mod": want["w/target/go.mod"], work.OpenGapsPath: gapRow}}
	f.install(t)
	if o := runS12(Input{Dir: "w", Record: lastRecord()}); o.Kind != Done || f.files["w/target/"+work.OpenGapsPath] != "" {
		t.Errorf("the head of S12: %s %q, open-gaps.tsv %q; want done, and no row written", o.Kind, o.Evidence, f.files["w/target/"+work.OpenGapsPath])
	}
	for _, c := range []struct {
		name     string
		f        *fakeRepo
		record   work.Record
		evidence string
	}{
		{"a baseline go.mod", &fakeRepo{trees: map[string][]git.TreeEntry{".": head}, shows: map[string]string{"go.mod": "module other\n"}}, lastRecord(),
			"go.mod: the head holds this path of the catalog entry with other bytes, and the baseline's own files stay byte for byte (REQ-018)"},
		{"a directory at a path of the entry", &fakeRepo{trees: map[string][]git.TreeEntry{".": {{Path: "docs/gates.tsv/x"}}}}, lastRecord(),
			"docs/gates.tsv: the head holds this path of the catalog entry with other bytes, and the baseline's own files stay byte for byte (REQ-018)"},
		{"no entry", &fakeRepo{}, append(work.Record{{"S01", "stack", "cobol", "answer", "S01-stack"}}, pinRecord()...), "the catalog entry cobol: catalog: "},
		{"no tree", &fakeRepo{fail: "ls-tree"}, lastRecord(), "the tree of the head: ls-tree HEAD .: failed"},
	} {
		c.f.install(t)
		if o := runS12(Input{Dir: "w", Record: c.record}); o.Kind != Fail || !strings.HasPrefix(o.Evidence, c.evidence) || len(c.f.files) != 0 {
			t.Errorf("%s: %s %q, %d files; want fail %q and no file", c.name, o.Kind, o.Evidence, len(c.f.files), c.evidence)
		}
	}
}

const protectionText = `{
  "required_status_checks": {
    "strict": true,
    "checks": [
      {
        "context": "static",
        "app_id": 15368
      },
      {
        "context": "layout",
        "app_id": 15368
      }
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "dismiss_stale_reviews": false,
    "require_code_owner_reviews": false,
    "require_last_push_approval": false,
    "required_approving_review_count": 0
  },
  "restrictions": null,
  "required_linear_history": false,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "block_creations": false,
  "required_conversation_resolution": true,
  "lock_branch": false,
  "allow_fork_syncing": false
}
`

const rulesetText = `{
  "name": "layup: the default branch",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": [
        "~DEFAULT_BRANCH",
        "refs/heads/layup-probe"
      ],
      "exclude": []
    }
  },
  "rules": [
    {
      "type": "pull_request",
      "parameters": {
        "dismiss_stale_reviews_on_push": false,
        "require_code_owner_review": false,
        "require_last_push_approval": false,
        "required_approving_review_count": 0,
        "required_review_thread_resolution": true
      }
    },
    {
      "type": "required_status_checks",
      "parameters": {
        "required_status_checks": [
          {
            "context": "static",
            "integration_id": 15368
          },
          {
            "context": "layout",
            "integration_id": 15368
          }
        ],
        "strict_required_status_checks_policy": true
      }
    },
    {
      "type": "non_fast_forward"
    },
    {
      "type": "deletion"
    }
  ],
  "bypass_actors": []
}
`

// twoKinds is a manifest with an active kind and a pending one, whose config
// path is docs/gates/x.txt.
const twoKinds = "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
	"static\tactive\tgo\tgo vet ./...\t./*.go\t—\n" +
	"layout\tpending\t—\t—\t./*.go\tdocs/gates/x.txt\n"

// S13 writes the protection file in the form of LAYUP's own and the ruleset,
// with one required check per kind of the manifest at the head, pinned to
// GitHub Actions, and a record row of each file; its commands are the push of
// layup-setup and the apply of the ruleset (D4 of #92).
func TestS13(t *testing.T) {
	f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	f.install(t)
	o := runS13(Input{Dir: "w", Record: lastRecord()})
	values := [][]string{{"branch-protection.sha256", sha([]byte(protectionText)), "computed", "sha256 docs/setup/branch-protection.json"},
		{"ruleset.sha256", sha([]byte(rulesetText)), "computed", "sha256 out/ruleset-default.json"}}
	if o.Kind != Done || o.Evidence != "the ruleset file, and its commands in commands.sh" || !o.Commit || !reflect.DeepEqual(o.Values, values) {
		t.Errorf("S13: %s %q %v, values %q; want done with the rows %q", o.Kind, o.Evidence, o.Commit, o.Values, values)
	}
	if got := f.files["w/target/docs/setup/branch-protection.json"]; got != protectionText {
		t.Errorf("branch-protection.json:\n%s\nwant\n%s", got, protectionText)
	}
	if got := f.files["w/out/ruleset-default.json"]; got != rulesetText {
		t.Errorf("ruleset-default.json:\n%s\nwant\n%s", got, rulesetText)
	}
	want := []Command{
		{2, "the push of the setup commits onto the default branch, a fast-forward from the root commit", "git -C target push origin layup-setup:main"},
		{4, "the apply of the ruleset of the default branch, with the Operator's login", "gh api --method POST 'repos/acme/widget/rulesets' --input out/ruleset-default.json"},
	}
	if got := commandsS13(lastRecord()); !reflect.DeepEqual(got, want) {
		t.Errorf("the commands %+v; want %+v", got, want)
	}
	f = &fakeRepo{}
	f.install(t)
	if o := runS13(Input{Dir: "w", Record: lastRecord()}); o.Kind != Fail || o.Evidence != "the head: docs/gates.tsv: no file docs/gates.tsv" || len(f.files) != 0 {
		t.Errorf("no manifest: %s %q; want fail and no file", o.Kind, o.Evidence)
	}
}

// The register holds the entries of setup.md in their order, each .sh file of
// the tree and each config path, with their sources; a path listed twice has
// one row, with the source catalog (D5 of #92, note 9 of its plan review).
func TestTheRulePathRegister(t *testing.T) {
	tree := []string{".github/gates.sh", ".githooks/install.sh", "README.md", "docs/adr/adr-lint.sh", "docs/gates/x.sh"}
	got := rulePaths(tree, []string{"docs/gates/floor.txt", "docs/gates/x.sh"}, []string{".github/gates.sh", "go.mod"})
	want := [][]string{{".github/", "", "baseline"}, {".githooks/", "", "baseline"}, {".gitattributes", "", "baseline"}, {"AGENTS.md", "", "baseline"},
		{"CLAUDE.md", "", "baseline"}, {"docs/engineering-discipline.md", "", "baseline"}, {"docs/issue-workflow.md", "", "baseline"},
		{"docs/ci/", "", "baseline"}, {"docs/tests/", "", "baseline"}, {".github/gates.sh", "", "catalog"}, {".githooks/install.sh", "", "baseline"},
		{"docs/adr/adr-lint.sh", "", "baseline"}, {"docs/gates/x.sh", "", "catalog"}, {"docs/gates.tsv", "", "catalog"},
		{"docs/gates/floor.txt", "", "catalog"}, {"docs/setup/", "", "architecture"}, {"docs/facts/", "", "architecture"},
		{"docs/guardrails.md", "added lines in section 2", "baseline"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the register\n got %q\nwant %q", got, want)
	}
}

const setupHead = "5555555555555555555555555555555555555555"

// S15 stops with O-verify when out/verify.tsv is missing, refuses a table of
// another form (exit 2) and a row that is not pass or clear (fail), and else
// writes the rule-path register from the head of layup-setup, with its value
// rows (D5, D6 of #92).
func TestS15(t *testing.T) {
	verify := "check\tresult\treason\npin\tpass\t—\ngate:layout\tclear\tpending: fixture not run\n"
	record := append(lastRecord(), []string{"S12", "catalog:.github/gates.sh", "x", "catalog", "go/files/.github/gates.sh.tmpl"})
	repo := func(disk map[string]string) *fakeRepo {
		return &fakeRepo{revs: map[string]string{"refs/heads/layup-setup^{commit}": setupHead}, disk: disk,
			trees: map[string][]git.TreeEntry{".": {{Path: ".github/gates.sh"}, {Path: "docs/x.sh"}, {Path: "go.mod"}}},
			shows: map[string]string{"docs/gates.tsv": twoKinds}}
	}
	f := repo(map[string]string{"w/out/verify.tsv": verify})
	f.install(t)
	o := runS15(Input{Dir: "w", Record: record})
	reg := "pattern\texception\tsource\n.github/\t—\tbaseline\n.githooks/\t—\tbaseline\n.gitattributes\t—\tbaseline\nAGENTS.md\t—\tbaseline\n" +
		"CLAUDE.md\t—\tbaseline\ndocs/engineering-discipline.md\t—\tbaseline\ndocs/issue-workflow.md\t—\tbaseline\ndocs/ci/\t—\tbaseline\n" +
		"docs/tests/\t—\tbaseline\n.github/gates.sh\t—\tcatalog\ndocs/x.sh\t—\tbaseline\ndocs/gates.tsv\t—\tcatalog\ndocs/gates/x.txt\t—\tcatalog\n" +
		"docs/setup/\t—\tarchitecture\ndocs/facts/\t—\tarchitecture\ndocs/guardrails.md\tadded lines in section 2\tbaseline\n"
	values := [][]string{{"setup.head", setupHead, "computed", "git rev-parse layup-setup"}, {"verify.sha256", sha([]byte(verify)), "computed", "sha256 out/verify.tsv"},
		{"rule-paths.sha256", sha([]byte(reg)), "computed", "sha256 out/rule-paths.tsv"}}
	if o.Kind != Done || o.Evidence != "every row of verify.tsv is pass or clear" || o.Commit || !reflect.DeepEqual(o.Values, values) {
		t.Errorf("S15: %s %q %v, values %q; want done, no commit of the runner, and %q", o.Kind, o.Evidence, o.Commit, o.Values, values)
	}
	if got := f.files["w/out/rule-paths.tsv"]; got != reg {
		t.Errorf("rule-paths.tsv:\n%s\nwant\n%s", got, reg)
	}
	f = repo(nil)
	f.install(t)
	stop := []StopRow{{"S15", "O-verify", "Run layup setup verify 'w' > 'w/out/verify.tsv', then run layup setup 'w' again.", ""}}
	if o := runS15(Input{Dir: "w", Record: record}); o.Kind != Stop || !reflect.DeepEqual(o.Stops, stop) || f.calls != nil {
		t.Errorf("no verify.tsv: %s %q, calls %q; want the stop %q and no call", o.Kind, o.Stops, f.calls, stop)
	}
	for _, c := range []struct{ text, kind, evidence string }{
		{"check\tresult\n", Invalid, "out/verify.tsv: line 1"},
		{"check\tresult\treason\njobs\tfail\tno CI job for the kind static\n", Fail, "out/verify.tsv: the row jobs is fail: no CI job for the kind static"},
		{"check\tresult\treason\ngate:static\tnot-active\tthe clean run: tool not found: go\n", Fail,
			"out/verify.tsv: the row gate:static is not-active: the clean run: tool not found: go"},
	} {
		f = repo(map[string]string{"w/out/verify.tsv": c.text})
		f.install(t)
		if o := runS15(Input{Dir: "w", Record: record}); o.Kind != c.kind || !strings.HasPrefix(o.Evidence, c.evidence) || len(f.files) != 0 {
			t.Errorf("verify.tsv %q: %s %q; want %s %q and no file", c.text, o.Kind, o.Evidence, c.kind, c.evidence)
		}
	}
}

// The hook of S15 makes the first commit of the orphan branch layup-records in
// a scratch work tree, with the four files, by who at pin.time; it takes the
// commit of a run that stopped, refuses another commit, and refuses a file
// that changed after S15 read it (D6 of #92).
func TestTheRecordsCommit(t *testing.T) {
	const records = "7777777777777777777777777777777777777777"
	verify, reg := "check\tresult\treason\npin\tpass\t—\n", "pattern\texception\tsource\nAGENTS.md\t—\tbaseline\n"
	record := append(lastRecord(), []string{"S15", "verify.sha256", sha([]byte(verify)), "computed", "sha256 out/verify.tsv"},
		[]string{"S15", "rule-paths.sha256", sha([]byte(reg)), "computed", "sha256 out/rule-paths.tsv"},
		[]string{"S15", "done", "every row of verify.tsv is pass or clear", "step", ""})
	var rec bytes.Buffer
	if err := tsv.Write(&rec, work.RecordSchema, record); err != nil {
		t.Fatal(err)
	}
	disk := map[string]string{"w/out/verify.tsv": verify, "w/out/rule-paths.tsv": reg}
	f := &fakeRepo{revs: map[string]string{"refs/heads/layup-setup^{commit}": setupHead}, disk: disk}
	f.install(t)
	if err := recordsS15("w", record, testWho); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{"README.md": recordsReadme, "setup/record.tsv": rec.String(), "setup/verify.tsv": verify, "rule-paths.tsv": reg}
	want := map[string]string{}
	for p, text := range files {
		want["/tmp/s/tree/"+p] = text
	}
	calls := []string{"rev-parse refs/heads/layup-records^{commit}", "rev-parse refs/heads/layup-setup^{commit}", "tempdir", "worktree add /tmp/s/tree " + setupHead,
		"switch --orphan layup-records", "write /tmp/s/tree/README.md", "write /tmp/s/tree/rule-paths.tsv", "write /tmp/s/tree/setup/record.tsv",
		"write /tmp/s/tree/setup/verify.tsv", `commit /tmp/s/tree "chore: the records of the setup" by LAYUP test at 2026-09-30T08:00:00Z`,
		"worktree remove /tmp/s/tree", "remove /tmp/s"}
	if !reflect.DeepEqual(f.files, want) || !slices.Equal(f.calls, calls) {
		t.Errorf("the records commit: files %q, calls %q; want %q and %q", f.files, f.calls, want, calls)
	}
	for _, fail := range []string{"worktree add", "switch", "commit", "worktree remove"} {
		f := &fakeRepo{revs: map[string]string{"refs/heads/layup-setup^{commit}": setupHead}, disk: disk, fail: fail}
		f.install(t)
		err := recordsS15("w", record, testWho)
		var in *InputError
		if err == nil || errors.As(err, &in) || f.calls[len(f.calls)-1] != "remove /tmp/s" {
			t.Errorf("%s fails: %v, calls %q; want an error that is not an input error, and the scratch directory removed", fail, err, f.calls)
		}
	}

	// The branch of a run that stopped after the commit.
	var entries []git.TreeEntry
	for p := range files {
		entries = append(entries, git.TreeEntry{Mode: "100644", Type: "blob", Path: p})
	}
	taken := func() *fakeRepo {
		return &fakeRepo{revs: map[string]string{"refs/heads/layup-records^{commit}": records}, roots: []string{records}, msg: recordsMessage,
			trees: map[string][]git.TreeEntry{".": entries}, shows: copyOf(files), disk: disk}
	}
	f = taken()
	f.install(t)
	if err := recordsS15("w", record, testWho); err != nil || slices.ContainsFunc(f.calls, func(c string) bool { return strings.HasPrefix(c, "commit") || strings.HasPrefix(c, "write") }) {
		t.Errorf("the commit of a run that stopped: %v, calls %q; want it taken, with no write and no commit", err, f.calls)
	}
	for name, change := range map[string]func(f *fakeRepo){
		"another file":    func(f *fakeRepo) { f.shows["README.md"] = "# other\n" },
		"a file more":     func(f *fakeRepo) { f.trees["."] = append(f.trees["."], git.TreeEntry{Mode: "100644", Path: "x"}) },
		"a parent":        func(f *fakeRepo) { f.roots = []string{setupHead} },
		"another message": func(f *fakeRepo) { f.msg = "chore: other" },
		"an executable": func(f *fakeRepo) {
			f.trees["."] = append([]git.TreeEntry{{Mode: "100755", Path: "README.md"}}, f.trees["."][1:]...)
		},
	} {
		f = taken()
		change(f)
		f.install(t)
		var in *InputError
		if err := recordsS15("w", record, testWho); !errors.As(err, &in) || err.Error() != "the branch layup-records of target has a commit that S15 did not make: start again in a new work area" {
			t.Errorf("%s: %v; want an input error", name, err)
		}
	}
	changed := map[string]string{"w/out/verify.tsv": verify + "x\tpass\t—\n", "w/out/rule-paths.tsv": reg}
	f = &fakeRepo{revs: map[string]string{"refs/heads/layup-setup^{commit}": setupHead}, disk: changed}
	f.install(t)
	if err := recordsS15("w", record, testWho); err == nil || err.Error() != "out/verify.tsv changed after S15 read it" || f.calls != nil {
		t.Errorf("a verify.tsv that changed: %v, calls %q; want an error before any call", err, f.calls)
	}
	if got := commandsS15(record); !reflect.DeepEqual(got, []Command{{3, "the push of the first records commit, the branch layup-records", "git -C target push origin layup-records"}}) {
		t.Errorf("the command of S15: %+v", got)
	}
}

// mapKeys gives the keys of a and b.
func mapKeys(a, b map[string]string) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range a {
			if !yield(k) {
				return
			}
		}
		for k := range b {
			if _, ok := a[k]; !ok && !yield(k) {
				return
			}
		}
	}
}

// copyOf gives a copy of m.
func copyOf(m map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

// The runner gives the hook Records the record that it then writes, with the
// step's value rows and its done row, and the identity of the run; an error of
// the hook is a fail with no done row, and an *InputError is exit 2 (note 3 of
// the plan review of #92).
func TestTheRecordsHook(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
		want string // the result and the evidence of S15, or the input error
	}{
		{"a commit", nil, "done records"},
		{"a commit that failed", errors.New("git commit: failed\nmore"), "fail the records commit failed: git commit: failed"},
		{"another commit", &InputError{errors.New("another commit")}, "input another commit"},
	} {
		f := &fakeSys{record: doneRows(14)}
		f.install(t)
		m := steps(f, map[string]Outcome{"S15": {Kind: Done, Evidence: "records", Values: [][]string{{"setup.head", "h", "computed", "git rev-parse layup-setup"}}}})
		var got work.Record
		var who git.Identity
		s := m["S15"]
		s.Records = func(dir string, r work.Record, w git.Identity) error { got, who = r, w; return c.err }
		m["S15"] = s
		before := len(f.record)
		res, err := Run("/w", m, testWho, noStep)
		row := ""
		if err != nil {
			row = "input " + err.Error()
		} else if len(res.Steps) == 15 {
			row = res.Steps[14].Result + " " + res.Steps[14].Evidence
		}
		var last []string
		if len(got) > 0 {
			last = got[len(got)-1]
		}
		if row != c.want || who != testWho || len(got) != before+2 || !slices.Equal(last, []string{"S15", "done", "records", "step", ""}) {
			t.Errorf("%s: %q, the hook got %d rows ending %q, by %v; want %q, the record with the value row and the done row", c.name, row, len(got), last, who, c.want)
		}
		if _, done := f.record.Value("S15", "done"); done != (c.err == nil) || c.err == nil && !reflect.DeepEqual(f.record, got) {
			t.Errorf("%s: the written record has the done row of S15: %v; want %v, and the record of the hook", c.name, done, c.err == nil)
		}
	}
}

// Steps gives S12 its evidence, the checks jobs and gates, S13 its commands,
// and S15 its command and the hook of the records commit.
func TestTheLastSteps(t *testing.T) {
	var named []string
	m := Steps(Brief{}, Calls{Checks: func(_ string, names []string) string { named = names; return "" }})
	if e := m["S12"].Evidence; e != nil {
		e("w")
	}
	if !slices.Equal(named, []string{"jobs", "gates"}) || m["S13"].Commands == nil || !m["S13"].HandOff || m["S13"].Evidence != nil ||
		m["S15"].Commands == nil || m["S15"].Records == nil || m["S15"].Evidence != nil {
		t.Errorf("the evidence of S12 %q; want jobs and gates, S13 a hand-off with commands, and S15 a command and the hook", named)
	}
}
