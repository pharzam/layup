package setup

// The last steps of a setup (task T-d6q5, #92): S12 writes the files of the
// catalog entry of the stack, S13 the protection of the default branch and its
// commands, and S15 the rule-path register and the first records commit.

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/catalog"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// catalogEntry gives the entry of a stack from the catalog of the binary; the
// unit tests replace it.
var catalogEntry = func(stack string) (*catalog.Entry, error) { return catalog.Read(catalog.Embedded(), stack) }

// The paths of the last steps: the gate manifest and the protection file of a
// target, and the files of the work area that S13 and S15 write.
const (
	manifestPath   = "docs/gates.tsv"
	protectionPath = "docs/setup/branch-protection.json"
	rulesetPath    = "out/ruleset-default.json"
	rulePathsPath  = "out/rule-paths.tsv"
)

// runS12 is S12 (D1 of #92): the files of the catalog entry of the stack, with
// the module path github.com/OWNER/NAME of the answer S01-name, the manifest
// docs/gates.tsv, and for each gap of the entry its row of open-gaps.tsv and
// its record row. A path of the entry that the head holds with other bytes is
// a fail: the baseline's own files stay byte for byte (REQ-018), and the head
// that holds the same bytes is the commit of a run that stopped.
func runS12(in Input) Outcome {
	target := filepath.Join(in.Dir, work.TargetPath)
	stack, _ := in.Record.Value("S01", "stack")
	name, _ := in.Record.Value("S01", "name")
	entry, err := catalogEntry(stack)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the catalog entry " + stack + ": " + firstLine(err)}
	}
	module := "github.com/" + name
	files, err := entry.Files(module)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: firstLine(err)}
	}
	manifest, err := entry.Manifest()
	if err != nil {
		return Outcome{Kind: Fail, Evidence: firstLine(err)}
	}
	files = append(files, catalog.File{Path: manifestPath, Data: manifest})
	head, err := sys.lsTree(target, "HEAD", ".")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the tree of the head: " + firstLine(err)}
	}
	values := [][]string{{"module", module, "answer", "S01-name"}}
	out := map[string][]byte{}
	for _, f := range files {
		for _, e := range head {
			if e.Path != f.Path && !strings.HasPrefix(e.Path, f.Path+"/") {
				continue
			}
			if old, err := sys.show(target, "HEAD", f.Path); e.Path != f.Path || err != nil || !bytes.Equal(old, f.Data) {
				return Outcome{Kind: Fail, Evidence: f.Path + ": the head holds this path of the catalog entry with other bytes, and the baseline's own files stay byte for byte (REQ-018)"}
			}
		}
		ref := stack + "/files/" + f.Path + ".tmpl"
		if f.Path == manifestPath {
			ref = stack + "/kinds.tsv"
		}
		values = append(values, []string{"catalog:" + f.Path, sha(f.Data), "catalog", ref})
		out[f.Path] = f.Data
	}
	list, _ := sys.show(target, "HEAD", work.OpenGapsPath)
	var gaps [][]string
	for _, g := range entry.Gaps() {
		values = append(values, []string{fmt.Sprintf("marker:%s:%d", g.Path, g.Line), g.Marker, "gap", work.OpenGapsPath})
		if !hasGapRow(list, g.Path, g.Marker) && !slices.ContainsFunc(gaps, func(r []string) bool { return r[0] == g.Path && r[1] == g.Marker }) {
			gaps = append(gaps, []string{g.Path, g.Marker, g.Question})
		}
	}
	if len(gaps) > 0 {
		if len(list) > 0 && list[len(list)-1] != '\n' {
			list = append(list, '\n')
		}
		var b bytes.Buffer
		if err := tsv.Write(&b, work.OpenGapsSchema, gaps); err != nil {
			return Outcome{Kind: Fail, Evidence: work.OpenGapsPath + ": " + err.Error()}
		}
		out[work.OpenGapsPath] = append(list, b.Bytes()...)
	}
	for _, p := range slices.Sorted(maps.Keys(out)) {
		if err := sys.write(targetFile(in.Dir, p), out[p]); err != nil {
			return Outcome{Kind: Fail, Evidence: "S12: " + firstLine(err)}
		}
	}
	return Outcome{Kind: Done, Evidence: "checks jobs and gates", Commit: true, Values: values}
}

// hasGapRow reports whether the text of open-gaps.tsv has the row of the file
// and the marker, read by tabs as check markers reads it.
func hasGapRow(list []byte, file, marker string) bool {
	for _, l := range strings.Split(string(list), "\n") {
		f := strings.Split(strings.TrimSuffix(l, "\r"), "\t")
		if len(f) > 1 && f[0] == file && f[1] == marker {
			return true
		}
	}
	return false
}

// manifestKinds gives the kinds of docs/gates.tsv at rev of the target, and the
// config paths of the kinds, each in their order.
func manifestKinds(target, rev string) (kinds, configs []string, err error) {
	data, err := sys.show(target, rev, manifestPath)
	if err != nil {
		return nil, nil, errors.New(manifestPath + ": " + firstLine(err))
	}
	rows, err := tsv.Read(data, catalog.ManifestSchema)
	if err != nil {
		return nil, nil, errors.New(manifestPath + ": " + err.Error())
	}
	for _, r := range rows {
		kinds = append(kinds, r[0])
		configs = append(configs, strings.Fields(r[5])...)
	}
	return kinds, configs, nil
}

// actionsApp is the ID of the GitHub Actions app, which runs the gate jobs:
// each required check is pinned to it (D4 of #92; its evidence is V-21 of
// docs/setup/record-T-n1hp.md, the body that GitHub read back).
const actionsApp = 15368

// jsonText gives v as JSON with an indent of two spaces and a line feed at
// its end, with no HTML escape.
func jsonText(v any) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	err := e.Encode(v)
	return b.Bytes(), err
}

// protectionBody gives docs/setup/branch-protection.json (D4 of #92): the
// form of LAYUP's own file, the body of the protection of a branch, with one
// required check per kind.
func protectionBody(kinds []string) ([]byte, error) {
	type check struct {
		Context string `json:"context"`
		AppID   int    `json:"app_id"`
	}
	type checks struct {
		Strict bool    `json:"strict"`
		Checks []check `json:"checks"`
	}
	type reviews struct {
		DismissStale bool `json:"dismiss_stale_reviews"`
		CodeOwner    bool `json:"require_code_owner_reviews"`
		LastPush     bool `json:"require_last_push_approval"`
		Approvals    int  `json:"required_approving_review_count"`
	}
	body := struct {
		Checks         checks    `json:"required_status_checks"`
		EnforceAdmins  bool      `json:"enforce_admins"`
		Reviews        reviews   `json:"required_pull_request_reviews"`
		Restrictions   *struct{} `json:"restrictions"`
		LinearHistory  bool      `json:"required_linear_history"`
		ForcePushes    bool      `json:"allow_force_pushes"`
		Deletions      bool      `json:"allow_deletions"`
		BlockCreations bool      `json:"block_creations"`
		Resolution     bool      `json:"required_conversation_resolution"`
		LockBranch     bool      `json:"lock_branch"`
		ForkSyncing    bool      `json:"allow_fork_syncing"`
	}{Checks: checks{Strict: true, Checks: []check{}}, EnforceAdmins: true, Resolution: true}
	for _, k := range kinds {
		body.Checks.Checks = append(body.Checks.Checks, check{k, actionsApp})
	}
	return jsonText(body)
}

// rulesetBody gives the ruleset of the default branch and of the ref
// layup-probe (D4 of #92): a pull request with no required review (an
// approval is an issue comment, and the App merges), one required check per
// kind pinned to GitHub Actions, no force push, no deletion, and no bypass.
func rulesetBody(kinds []string) ([]byte, error) {
	type check struct {
		Context       string `json:"context"`
		IntegrationID int    `json:"integration_id"`
	}
	type rule struct {
		Type       string `json:"type"`
		Parameters any    `json:"parameters,omitempty"`
	}
	type refs struct {
		Include []string `json:"include"`
		Exclude []string `json:"exclude"`
	}
	checks := []check{}
	for _, k := range kinds {
		checks = append(checks, check{k, actionsApp})
	}
	return jsonText(struct {
		Name        string `json:"name"`
		Target      string `json:"target"`
		Enforcement string `json:"enforcement"`
		Conditions  struct {
			RefName refs `json:"ref_name"`
		} `json:"conditions"`
		Rules  []rule `json:"rules"`
		Bypass []any  `json:"bypass_actors"`
	}{
		Name: rulesetName, Target: "branch", Enforcement: "active",
		Conditions: struct {
			RefName refs `json:"ref_name"`
		}{refs{[]string{"~DEFAULT_BRANCH", "refs/heads/layup-probe"}, []string{}}},
		Rules: []rule{
			{"pull_request", struct {
				DismissStale bool `json:"dismiss_stale_reviews_on_push"`
				CodeOwner    bool `json:"require_code_owner_review"`
				LastPush     bool `json:"require_last_push_approval"`
				Approvals    int  `json:"required_approving_review_count"`
				Resolution   bool `json:"required_review_thread_resolution"`
			}{Resolution: true}},
			{"required_status_checks", struct {
				Checks []check `json:"required_status_checks"`
				Strict bool    `json:"strict_required_status_checks_policy"`
			}{checks, true}},
			{"non_fast_forward", nil},
			{"deletion", nil},
		},
		Bypass: []any{},
	})
}

// rulesetName names the tool that wrote the ruleset and the branch, for the
// Operator who reads the rulesets of the repository (note 1 of the plan review
// of #92).
const rulesetName = "layup: the default branch"

// runS13 is S13 (D4 of #92): docs/setup/branch-protection.json in the tree and
// out/ruleset-default.json in the work area, each with its record row, from
// the kinds of the manifest at the head, which are the names of the gate jobs
// (K30); its commands are the hand-off.
func runS13(in Input) Outcome {
	kinds, _, err := manifestKinds(filepath.Join(in.Dir, work.TargetPath), "HEAD")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the head: " + err.Error()}
	}
	protection, err := protectionBody(kinds)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "S13: " + err.Error()}
	}
	ruleset, err := rulesetBody(kinds)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "S13: " + err.Error()}
	}
	if err := sys.write(targetFile(in.Dir, protectionPath), protection); err != nil {
		return Outcome{Kind: Fail, Evidence: "S13: " + firstLine(err)}
	}
	if err := sys.write(filepath.Join(in.Dir, filepath.FromSlash(rulesetPath)), ruleset); err != nil {
		return Outcome{Kind: Fail, Evidence: "S13: " + firstLine(err)}
	}
	return Outcome{Kind: Done, Evidence: "the ruleset file, and its commands in commands.sh", Commit: true, Values: [][]string{
		{"branch-protection.sha256", sha(protection), "computed", "sha256 " + protectionPath},
		{"ruleset.sha256", sha(ruleset), "computed", "sha256 " + rulesetPath},
	}}
}

// commandsS13 gives the commands of S13 (D4 of #92), run from the work area:
// the push of layup-setup onto the default branch, and the apply of the
// ruleset with the Operator's login.
func commandsS13(r work.Record) []Command {
	name, _ := r.Value("S01", "name")
	return []Command{
		{Order: 2, Comment: "the push of the setup commits onto the default branch, a fast-forward from the root commit", Text: "git -C " + work.TargetPath + " push origin layup-setup:main"},
		{Order: 4, Comment: "the apply of the ruleset of the default branch, with the Operator's login", Text: "gh api --method POST " + quote("repos/"+name+"/rulesets") + " --input " + rulesetPath},
	}
}

// RulePathsSchema is the form of the rule-path register: the block rule-paths
// of docs/spec/setup.md.
var RulePathsSchema = tsv.Schema{Name: "rule-paths", Location: "records:rule-paths.tsv", Columns: []tsv.Column{
	{Name: "pattern", Type: "text", Key: true},
	{Name: "exception", Type: "text"},
	{Name: "source", Type: "enum(baseline|catalog|architecture)"},
}}

// rulePaths gives the rows of the rule-path register (D5 of #92) in the order
// of the list of setup.md: tree is the paths of the tree at the head of
// layup-setup, configs the config paths of its manifest, and written the paths
// that S12 wrote from the catalog entry. A path listed twice has one row, with
// the source catalog (note 9 of the plan review).
func rulePaths(tree, configs, written []string) [][]string {
	var rows [][]string
	add := func(p, exception, source string) {
		if i := slices.IndexFunc(rows, func(r []string) bool { return r[0] == p }); i >= 0 {
			if source == "catalog" {
				rows[i][2] = source
			}
			return
		}
		rows = append(rows, []string{p, exception, source})
	}
	for _, p := range []string{".github/", ".githooks/", ".gitattributes", "AGENTS.md", "CLAUDE.md", "docs/engineering-discipline.md",
		"docs/issue-workflow.md", "docs/ci/", "docs/tests/"} {
		add(p, "", "baseline")
	}
	for _, p := range tree {
		if strings.HasSuffix(p, ".sh") {
			source := "baseline"
			if slices.Contains(written, p) {
				source = "catalog"
			}
			add(p, "", source)
		}
	}
	add(manifestPath, "", "catalog")
	for _, p := range configs {
		add(p, "", "catalog")
	}
	add("docs/setup/", "", "architecture")
	add("docs/facts/", "", "architecture")
	add("docs/guardrails.md", "added lines in section 2", "baseline")
	return rows
}

// register gives the bytes of the rule-path register of the target at the
// head of layup-setup; the record's rows catalog:<path> of S12 name the files
// of the catalog entry.
func register(target string, r work.Record) ([]byte, error) {
	const head = "refs/heads/layup-setup"
	entries, err := sys.lsTree(target, head, ".")
	if err != nil {
		return nil, errors.New("the tree of layup-setup: " + firstLine(err))
	}
	_, configs, err := manifestKinds(target, head)
	if err != nil {
		return nil, err
	}
	var tree, written []string
	for _, e := range entries {
		tree = append(tree, e.Path)
	}
	slices.Sort(tree)
	for _, row := range r {
		if p, ok := strings.CutPrefix(row[1], "catalog:"); ok && row[0] == "S12" {
			written = append(written, p)
		}
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, RulePathsSchema, rulePaths(tree, configs, written)); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// The records branch and its first commit (D6 of #92).
const (
	recordsBranch  = "layup-records"
	recordsMessage = "chore: the records of the setup"
)

// recordsReadme is the fixed text of README.md of the records branch (D6 of
// #92, note 2 of its plan review); docs/spec/setup.md holds the same text.
const recordsReadme = "# The records of this repository\n\n" +
	"This branch, `layup-records`, holds the records that LAYUP keeps for this\n" +
	"repository. It is an orphan branch: it shares no commit with the default\n" +
	"branch.\n\n" +
	"Only `layup run` writes this branch from Start on. Its first commit holds the\n" +
	"records of the setup.\n\n" +
	"Each file is a table of tab-separated values with a header row, or Markdown,\n" +
	"so a person reads it with no tool. A plain `git clone` carries the branch as\n" +
	"`origin/layup-records`.\n\n" +
	"- `setup/record.tsv`: each value of the setup, with its source, and one done\n" +
	"  row per step.\n" +
	"- `setup/verify.tsv`: the table of `layup setup verify` before the records\n" +
	"  commit.\n" +
	"- `rule-paths.tsv`: the rule-path register: the paths whose change is a change\n" +
	"  of the rules.\n"

// verifyAsk is the ask of the stop O-verify: the command that writes
// out/verify.tsv of the work area at dir.
func verifyAsk(dir string) string {
	return "Run layup setup verify " + quote(dir) + " > " + quote(filepath.Join(dir, filepath.FromSlash(work.VerifyPath))) +
		", then run layup setup " + quote(dir) + " again."
}

// runS15 is S15 (D5, D6 of #92): with no out/verify.tsv, it stops with the
// question O-verify; a verify.tsv that its schema refuses is an input error,
// and one with a row that is not pass or clear is a fail. Else it writes the
// rule-path register and its value rows; its hook Records makes the first
// records commit.
func runS15(in Input) Outcome {
	data, err := sys.read(filepath.Join(in.Dir, filepath.FromSlash(work.VerifyPath)))
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return Outcome{Kind: Stop, Stops: []StopRow{{"S15", "O-verify", verifyAsk(in.Dir), ""}}}
	case err != nil:
		return Outcome{Kind: Fail, Evidence: work.VerifyPath + ": " + firstLine(err)}
	}
	rows, err := tsv.Read(data, work.VerifySchema)
	if err != nil {
		return Outcome{Kind: Invalid, Evidence: work.VerifyPath + ": " + err.Error()}
	}
	for _, r := range rows {
		if r[1] != "pass" && r[1] != "clear" {
			return Outcome{Kind: Fail, Evidence: fmt.Sprintf("%s: the row %s is %s: %s", work.VerifyPath, r[0], r[1], r[2])}
		}
	}
	target := filepath.Join(in.Dir, work.TargetPath)
	head, err := sys.revParse(target, "refs/heads/layup-setup^{commit}")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the head of layup-setup: " + firstLine(err)}
	}
	reg, err := register(target, in.Record)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the rule-path register: " + firstLine(err)}
	}
	if err := sys.write(filepath.Join(in.Dir, filepath.FromSlash(rulePathsPath)), reg); err != nil {
		return Outcome{Kind: Fail, Evidence: "S15: " + firstLine(err)}
	}
	return Outcome{Kind: Done, Evidence: "every row of verify.tsv is pass or clear", Values: [][]string{
		{"setup.head", head, "computed", "git rev-parse layup-setup"},
		{"verify.sha256", sha(data), "computed", "sha256 " + work.VerifyPath},
		{"rule-paths.sha256", sha(reg), "computed", "sha256 " + rulePathsPath},
	}}
}

// recordFiles gives the four files of the records commit: the README, the
// record r, and the two files of the work area whose hashes r holds.
func recordFiles(dir string, r work.Record) (map[string][]byte, error) {
	var rec bytes.Buffer
	if err := tsv.Write(&rec, work.RecordSchema, r); err != nil {
		return nil, err
	}
	files := map[string][]byte{"README.md": []byte(recordsReadme), "setup/record.tsv": rec.Bytes()}
	for _, f := range []struct{ row, from, to string }{
		{"verify.sha256", work.VerifyPath, "setup/verify.tsv"},
		{"rule-paths.sha256", rulePathsPath, "rule-paths.tsv"},
	} {
		data, err := sys.read(filepath.Join(dir, filepath.FromSlash(f.from)))
		if err != nil {
			return nil, err
		}
		if want, _ := r.Value("S15", f.row); sha(data) != want {
			return nil, errors.New(f.from + " changed after S15 read it")
		}
		files[f.to] = data
	}
	return files, nil
}

// recordsS15 is the hook Records of S15 (note 3 of the plan review of #92):
// with the record that the runner then writes, it makes the first commit of
// the orphan branch layup-records in a scratch work tree outside the work
// area, by who at pin.time, so main and layup-setup do not move. A branch
// layup-records that is the commit of these files, with no parent and the
// message of S15, is the commit of a run that stopped, and S15 takes it, as
// S03 takes its root commit; another is an input error.
func recordsS15(dir string, r work.Record, who git.Identity) error {
	target := filepath.Join(dir, work.TargetPath)
	files, err := recordFiles(dir, r)
	if err != nil {
		return err
	}
	if commit, err := sys.revParse(target, "refs/heads/"+recordsBranch+"^{commit}"); err == nil {
		if !sameRecords(target, commit, files) {
			return &InputError{errors.New("the branch " + recordsBranch + " of " + work.TargetPath + " has a commit that S15 did not make: start again in a new work area")}
		}
		return nil
	}
	at, _ := r.Value("S02", "pin.time")
	when, err := time.Parse(timeForm, at)
	if err != nil {
		return errors.New("the record has no pin.time of the form " + timeForm)
	}
	who.Time = when
	head, err := sys.revParse(target, "refs/heads/layup-setup^{commit}")
	if err != nil {
		return err
	}
	scratch, err := sys.tempDir()
	if err != nil {
		return err
	}
	tree := filepath.Join(scratch, "tree")
	if err = sys.worktreeAdd(target, tree, head); err == nil {
		err = errors.Join(commitRecords(tree, files, who), sys.worktreeRemove(target, tree))
	}
	return errors.Join(err, sys.removeAll(scratch))
}

// commitRecords puts the scratch tree on the new branch layup-records, with no
// commit and no file, and commits the files there.
func commitRecords(tree string, files map[string][]byte, who git.Identity) error {
	if err := sys.switchOrphan(tree, recordsBranch); err != nil {
		return err
	}
	for _, p := range slices.Sorted(maps.Keys(files)) {
		if err := sys.write(filepath.Join(tree, filepath.FromSlash(p)), files[p]); err != nil {
			return err
		}
	}
	return sys.commit(tree, recordsMessage, who)
}

// sameRecords reports whether commit is the records commit of files: no
// parent, the message of S15, and a tree of the files, each a regular file
// with its bytes.
func sameRecords(target, commit string, files map[string][]byte) bool {
	roots, err := sys.rootCommits(target, commit)
	msg, merr := sys.message(target, commit)
	entries, lerr := sys.lsTree(target, commit, ".")
	if err != nil || merr != nil || lerr != nil || !slices.Equal(roots, []string{commit}) || msg != recordsMessage || len(entries) != len(files) {
		return false
	}
	for _, e := range entries {
		want, ok := files[e.Path]
		got, err := sys.show(target, commit, e.Path)
		if !ok || e.Mode != "100644" || err != nil || !bytes.Equal(got, want) {
			return false
		}
	}
	return true
}

// commandsS15 gives the command of S15: the push of the records branch.
func commandsS15(work.Record) []Command {
	return []Command{{Order: 3, Comment: "the push of the first records commit, the branch " + recordsBranch, Text: "git -C " + work.TargetPath + " push origin " + recordsBranch}}
}

// tempDir gives a new scratch directory in the temporary directory of the
// host.
func tempDir() (string, error) { return os.MkdirTemp("", "layup-setup-") }
