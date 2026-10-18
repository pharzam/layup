package run

// The step forge-check (docs/spec/scaffold.md, REQ-003, Scaffold 7 and its paragraph "Decided here" of
// task T-wphp; row 57a of the plan, #280): the eight checks of the forge, and their table.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
)

const (
	forgeCheckDone    = "the forge check passed"
	forgeCheckMessage = "layup run: the forge check"
	probeMessage      = "layup: the probe of the forge check"
	probeBranch       = "layup-probe"
)

// probeRefusal is the line of GitHub's refusal of a push for a rule violation, for the ref of the probe and for it alone, as the
// whole ref: `GH013: Repository rule violations found for refs/heads/layup-probe.` (the full stop may be left out). GitHub prints
// the violated rules after it and the URL of the rules of the ref, and names no ruleset (docs/spec/scaffold.md, The probe).
var probeRefusal = regexp.MustCompile(`(?m)GH013: Repository rule violations found for refs/heads/` + regexp.QuoteMeta(probeBranch) + `\.?\s*$`)

// forgeCheckJob is the step with its seams: the two records that the checks read (setup/tree.tsv and
// setup/rulesets.txt, nil when absent); the work area, made sure of for the last commit; a ruleset file
// of its out/; the forge's repository, rules of a branch and activity of a ref; the probe of a head, which
// gives its commit; the records commit of the table; the clean of a detail of host paths; the check of the
// code-owner commit, which a head other than the last commit of the tree record is given to (the tree check,
// docs/spec/scaffold.md); and a line of the step's own.
type forgeCheckJob struct {
	branch   string
	app      int64
	run      string
	now      func() time.Time
	records  func() (tree, applied []byte, err error)
	area     func(last string) error
	read     func(name string) ([]byte, error)
	repo     func() (forge.Repository, error)
	rules    func(branch string) ([]forge.Rule, error)
	probe    func(head string) (string, error)
	activity func(ref string) error
	commit   func(table []byte) error
	clean    func(string) string
	landing  func(head, last string) error // asked of a head that is another commit than the last: nil, it is the code-owner commit
	note     func(string)                  // a line of the step's own
	// merges gives the merge commits of the tasks: the SHA of each that a merged event records, and of each that the forge
	// shows merged and no event records, by the task; nil is none (the tree check below). parents gives the parents of a commit.
	merges  func() (recorded, unrecorded map[string]string, err error)
	parents func(sha string) ([]string, error)
}

// explainHead walks the first parents from head: each commit must be a merge of recorded or unrecorded, and the walk
// passes at a commit for which bottom holds, once a merge was accepted. It gives a line for each unrecorded merge.
func explainHead(head, last string, parents func(string) ([]string, error), recorded, unrecorded map[string]string, bottom func(string) bool) ([]string, error) {
	return nil, errors.New("not built")
}

// firstLine is the text up to its first line feed.
func firstLine(s string) string { return strings.SplitN(s, "\n", 2)[0] }

// do runs every check, whatever an earlier one gave, and commits the table in either outcome; the step is
// done when each row is a pass, else it names the first that is not. Records that cannot be read, a table
// that its reader refuses and a commit that fails end it with no table.
func (j forgeCheckJob) do() Step {
	const name = "forge-check"
	stop := func(what string, err error) Step {
		return Step{name, "fail", what + firstLine(err.Error()) + "; no table was committed"}
	}
	tree, applied, err := j.records()
	if err != nil {
		return stop("the records: ", err)
	}
	rows := j.checks(tree, applied)
	var table bytes.Buffer
	err = tsv.Write(&table, ForgeCheckSchema, rows)
	if err == nil {
		_, err = ReadForgeCheck(table.Bytes())
	}
	if err != nil {
		return stop("the table: ", err)
	}
	if err := j.commit(table.Bytes()); err != nil {
		return stop("the records commit: ", err)
	}
	for _, r := range rows {
		if r[1] != "pass" {
			return Step{name, "fail", r[0] + ": " + r[2]}
		}
	}
	return Step{name, "done", forgeCheckDone}
}

// checks gives the eight rows of the table, each with the run and the time of its check.
func (j forgeCheckJob) checks(tree, applied []byte) [][]string {
	var rows [][]string
	row := func(check, detail string, ok bool) {
		result := "pass"
		if !ok {
			result = "fail"
		}
		rows = append(rows, []string{check, result, j.clean(strings.ToValidUTF8(firstLine(detail), "?")), j.run, stamp(j.now())})
	}
	rowErr := func(check string, err error) {
		if err != nil {
			row(check, err.Error(), false)
		} else {
			row(check, "", true)
		}
	}
	last, lastErr := lastOfTree(tree)
	repo, repoErr := j.repo()
	head := repo.Heads[j.branch]
	switch {
	case lastErr != nil:
		row("tree", lastErr.Error(), false)
	case repoErr != nil:
		row("tree", "the forge: "+repoErr.Error(), false)
	default:
		detail, ok := j.treeHead(head, last)
		row("tree", detail, ok)
	}
	areaErr := lastErr
	if lastErr == nil {
		if err := j.area(last); err != nil {
			areaErr = fmt.Errorf("the work area: %w", err)
		}
	}
	// The effective rules of layup-probe are read once, whatever the work area and the ruleset file are: the row rules:layup-probe compares
	// them with the file, and the attribution of the probe's refusal reads the ruleset of each of their creation rules.
	probeRules, probeRulesErr := j.rules(probeBranch)
	for _, c := range []struct {
		check, file string
		read        func() ([]forge.Rule, error)
	}{
		{"rules:" + j.branch, "ruleset-default.json", func() ([]forge.Rule, error) { return j.rules(j.branch) }},
		{"rules:layup-probe", "ruleset-default.json", func() ([]forge.Rule, error) { return probeRules, probeRulesErr }},
		{"rules:layup-records", "ruleset-records.json", func() ([]forge.Rule, error) { return j.rules("layup-records") }},
	} {
		detail, ok := j.rulesRow(areaErr, c.file, c.read)
		row(c.check, detail, ok)
	}
	outs, outsErr := appliedOutputs(applied)
	for _, c := range []struct {
		check, file string
		test        func([]byte) error
	}{{"bypass:default", applyFiles[0], bypassDefault}, {"bypass:records", applyFiles[1], func(d []byte) error { return bypassRecords(d, j.app) }}} {
		err := outsErr
		if err == nil {
			err = c.test(outs[c.file])
		}
		rowErr(c.check, err)
	}
	switch {
	case repoErr != nil:
		row("probe", "probe not run: the head of "+j.branch+" was not read", false)
	case head == "":
		row("probe", "probe not run: the forge gives no head for "+j.branch, false)
	default:
		sha, err := j.probe(head)
		result, detail := probeVerdict(sha, err, attribution(probeRules, probeRulesErr, outs, outsErr))
		row("probe", detail, result == "pass")
	}
	rowErr("activity", j.activity(j.branch))
	return rows
}

// treeHead judges the head of the default branch (docs/spec/scaffold.md, The tree check): the last commit of
// setup/tree.tsv passes; so does the code-owner commit of a target set up before the code-owner rule, which only the
// check of the job can tell (never asked of no head), and then the step says so on its own line while the table keeps
// its empty detail. Any other head fails, with the reason of the check when it was asked.
func (j forgeCheckJob) treeHead(head, last string) (string, bool) {
	if head == last {
		return "", true
	}
	none := head
	if none == "" {
		none = "none"
	}
	plain := fmt.Sprintf("the head of %s is %s, and setup/tree.tsv ends at %s", j.branch, none, last)
	if head == "" {
		return plain, false
	}
	if err := j.landing(head, last); err != nil {
		return plain + "; it is not the code-owner commit: " + err.Error(), false
	}
	j.note("tree: " + head + " is the code-owner commit above " + last)
	return "", true
}

// lastOfTree gives the last commit of setup/tree.tsv, or the reason that it cannot.
func lastOfTree(data []byte) (string, error) {
	if data == nil {
		return "", errors.New("setup/tree.tsv is not in the records")
	}
	rows, err := ReadSetupTree(data)
	if err != nil {
		return "", fmt.Errorf("setup/tree.tsv: %w", err)
	}
	return rows[len(rows)-1][1], nil
}

// appliedOutputs gives the outputs of the applies from the text of setup/rulesets.txt, or the reason.
func appliedOutputs(data []byte) (map[string][]byte, error) {
	if data == nil {
		return nil, errors.New("setup/rulesets.txt is not in the records")
	}
	outs, err := readApplied(data)
	if err != nil {
		return nil, fmt.Errorf("setup/rulesets.txt: %w", err)
	}
	return outs, nil
}

// rulesRow reads back the rules of a branch, which read gives, against a ruleset file; a work area that failed fails it.
func (j forgeCheckJob) rulesRow(areaErr error, file string, read func() ([]forge.Rule, error)) (string, bool) {
	if areaErr != nil {
		return areaErr.Error(), false
	}
	data, err := j.read(file)
	if err != nil {
		return err.Error(), false
	}
	want, err := rulesetRules(data)
	if err != nil {
		return file + ": " + err.Error(), false
	}
	got, err := read()
	if err != nil {
		return err.Error(), false
	}
	return readBack(want, got)
}

// rulesetRules gives the rules of a ruleset file, in its order; a file with none is an error.
func rulesetRules(data []byte) ([]forge.Rule, error) {
	var f struct {
		Rules []struct {
			Type       string         `json:"type"`
			Parameters map[string]any `json:"parameters"`
		} `json:"rules"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, errors.New("it is not a JSON object")
	}
	if len(f.Rules) == 0 {
		return nil, errors.New("it holds no rule")
	}
	rules := make([]forge.Rule, len(f.Rules))
	for i, r := range f.Rules {
		rules[i] = forge.Rule{Type: r.Type, Parameters: r.Parameters}
	}
	return rules, nil
}

// readBack compares the rules that a ruleset file holds with those that the forge returns: each needs one
// returned rule of its type whose parameters equal each of its own, so parameters of two returned rules are
// never combined; a parameter that the forge returns and the file lacks is ignored, and one that the file
// holds at the default that GitHub omits (omittedAtDefault) and the forge lacks is equal. The detail of a pass
// is empty, or names the rule types that the forge returns and the file lacks; the detail of a fail names
// the first rule that is missing, or the first parameter that is missing or differs.
func readBack(want, got []forge.Rule) (string, bool) {
	wanted := map[string]bool{}
	for _, w := range want {
		wanted[w.Type] = true
		detail, found, first := w.Type+" is missing", false, true
		for _, g := range got {
			if g.Type != w.Type {
				continue
			}
			d := parameterDiff(w.Type, w.Parameters, g.Parameters)
			if d == "" {
				found = true
				break
			} else if first {
				detail, first = w.Type+": "+d, false
			}
		}
		if !found {
			return detail, false
		}
	}
	extra := map[string]bool{}
	for _, g := range got {
		if !wanted[g.Type] {
			extra[g.Type] = true
		}
	}
	if len(extra) == 0 {
		return "", true
	}
	types := make([]string, 0, len(extra))
	for t := range extra {
		types = append(types, t)
	}
	sort.Strings(types)
	return "extra: " + strings.Join(types, ", "), true
}

// omittedAtDefault names, by the type of the rule, the parameters that GitHub leaves out of the rule that it returns when
// the file wrote their default: on a repository that is not a fork the rule update comes back with no
// update_allows_fetch_and_merge, though the file wrote false (first hand, run 6 of the demo of M2d, 2026-10-15, in the
// answer of the apply, the read of the ruleset and the effective rules; task T-j7c5, #330; docs/spec/scaffold.md, The
// read-back). An entry needs first-hand evidence: GitHub echoes the other false parameters, those of pull_request.
var omittedAtDefault = map[string]map[string]any{"update": {"update_allows_fetch_and_merge": false}}

// parameterDiff names the first parameter of want, in the order of the names, that got lacks or holds
// with another value, or is "". A parameter that got lacks is no difference when it is one of
// omittedAtDefault for the type of the rule and want holds that default.
func parameterDiff(ruleType string, want, got map[string]any) string {
	names := make([]string, 0, len(want))
	for n := range want {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		g, ok := got[n]
		switch {
		case !ok:
			if d, omitted := omittedAtDefault[ruleType][n]; omitted && reflect.DeepEqual(want[n], d) {
				continue
			}
			return n + " is missing"
		case !sameParameter(n, want[n], g):
			return n + " differs"
		}
	}
	return ""
}

// sameParameter compares two parameters as the decoded JSON: the required checks as the sets of their
// pairs of a context and an integration_id, in any order, and every other value as it is.
func sameParameter(name string, want, got any) bool {
	if name == "required_status_checks" {
		w, wok := checkPairs(want)
		g, gok := checkPairs(got)
		if wok && gok {
			return reflect.DeepEqual(w, g)
		}
	}
	return reflect.DeepEqual(want, got)
}

// checkPairs gives the set of the pairs of a list of required checks; a context or an integration_id that is
// absent is no 0 and no empty string.
func checkPairs(v any) (map[string]bool, bool) {
	list, ok := v.([]any)
	if !ok {
		return nil, false
	}
	pairs := map[string]bool{}
	for _, e := range list {
		m, ok := e.(map[string]any)
		if !ok {
			return nil, false
		}
		pairs[fmt.Sprintf("%#v|%#v", m["context"], m["integration_id"])] = true
	}
	return pairs, true
}

// readApplied is the inverse of rulesetsRecord: the output of each apply, by its file, from the text of
// setup/rulesets.txt. An output is what lies between its marker line and the line feed before the next
// marker line, or before the end of the text; the final line feed is the separator.
func readApplied(record []byte) (map[string][]byte, error) {
	text, first, second := string(record), "== "+applyFiles[0]+"\n", "\n== "+applyFiles[1]+"\n"
	i := strings.Index(text, second)
	if !strings.HasPrefix(text, first) || i < len(first) || !strings.HasSuffix(text, "\n") || i+len(second) > len(text)-1 {
		return nil, errors.New("it is not the record of the two applies")
	}
	return map[string][]byte{applyFiles[0]: []byte(text[len(first):i]), applyFiles[1]: []byte(text[i+len(second) : len(text)-1])}, nil
}

// bypassActor is an entry of the bypass list of an apply output.
type bypassActor struct {
	ID   int64  `json:"actor_id"`
	Type string `json:"actor_type"`
	Mode string `json:"bypass_mode"`
}

// bypassActors gives the bypass_actors of an apply output; an output with no list is an error.
func bypassActors(output []byte) ([]bypassActor, error) {
	var out struct {
		Actors *[]bypassActor `json:"bypass_actors"`
	}
	if err := json.Unmarshal(output, &out); err != nil || out.Actors == nil {
		return nil, errors.New("the output has no list of bypass actors")
	}
	return *out.Actors, nil
}

// bypassDefault checks that the bypass list of the default branch's ruleset is empty.
func bypassDefault(output []byte) error {
	actors, err := bypassActors(output)
	if err == nil && len(actors) != 0 {
		err = errors.New("the bypass list of the default branch's ruleset is not empty")
	}
	return err
}

// bypassRecords checks that the bypass list of the records branch's ruleset is the App alone: one
// Integration of the App's ID, always.
func bypassRecords(output []byte, app int64) error {
	actors, err := bypassActors(output)
	if err == nil && (len(actors) != 1 || actors[0].ID != app || actors[0].Type != "Integration" || actors[0].Mode != "always") {
		err = fmt.Errorf("the bypass list of the records branch's ruleset is not the App alone (Integration %d, always)", app)
	}
	return err
}

// probeVerdict is the row of the probe from the push: a pass only for a refusal of the remote (git exits
// non-zero) whose standard error holds GitHub's line of GH013 for the ref of the probe (probeRefusal), when the
// attribution of the refusal to the default branch's ruleset has no fault (it is nil); a push that is accepted
// gives the commit, a refusal that the attribution does not explain is a fail that says so, and any other failure
// is "probe not run", never a pass.
func probeVerdict(sha string, err error, attributed error) (result, detail string) {
	if err == nil {
		return "fail", "probe accepted: " + sha
	}
	var failed *git.FailedError
	if errors.As(err, &failed) {
		if failed.Code > 0 && probeRefusal.MatchString(failed.Stderr) {
			if attributed != nil {
				return "fail", "probe refused, but not attributed to the default branch's ruleset: " + firstLine(attributed.Error())
			}
			return "pass", ""
		}
		if line := firstLine(strings.TrimSpace(failed.Stderr)); line != "" {
			return "fail", "probe not run: " + line
		}
	}
	return "fail", "probe not run: " + firstLine(err.Error())
}

// pushProbe makes the empty commit of the probe on top of head in dir (the tree of head, head as its parent,
// the identity who) and pushes it to layup-probe of url; it gives the commit, also when the push fails. The
// commit is on no branch of dir.
func pushProbe(dir, url string, auth git.Auth, head string, who git.Identity) (string, error) {
	tree, err := git.RevParse(dir, head+"^{tree}")
	if err != nil {
		return "", err
	}
	sha, err := git.CommitTree(dir, tree, head, probeMessage, who)
	if err != nil {
		return "", err
	}
	return sha, git.Push(dir, url, sha, probeBranch, auth)
}

// forgeCheckStep is the step forge-check of the restart: the job on the records, the forge, the work
// area and the run's clone. It follows scaffold-push when that step is done and not waiting, and, once the
// Scaffold is done, it comes between lease and probe (docs/spec/scaffold.md, At each later start).
func (r *state) forgeCheckStep(ctx context.Context) Step {
	const name = "forge-check"
	host, err := filepath.Abs(r.cfg.Dir)
	if err != nil {
		return fail(name, err)
	}
	area := filepath.Join(host, "work", r.cfg.Owner, r.cfg.Name)
	return forgeCheckJob{
		branch: r.defaultBranch, app: r.cfg.Register.AppID, run: r.cfg.RunID, now: r.cfg.Clock.Now, clean: r.clean,
		records: func() ([]byte, []byte, error) {
			tree, err := r.store.ReadFile(ctx, "setup/tree.tsv")
			if err != nil {
				return nil, nil, err
			}
			applied, err := r.store.ReadFile(ctx, "setup/rulesets.txt")
			return tree, applied, err
		},
		landing: func(head, last string) error {
			record, err := r.store.ReadFile(ctx, "setup/record.tsv")
			if err != nil {
				return err
			}
			return codeOwnerLanding(r.target, record, r.approvers)(head, last)
		},
		note: func(line string) {
			if r.cfg.Progress != nil {
				r.cfg.Progress(line)
			}
		},
		area:  func(last string) error { return r.workAreaAgain(ctx, area, last) },
		read:  func(file string) ([]byte, error) { return os.ReadFile(filepath.Join(area, "out", file)) },
		repo:  func() (forge.Repository, error) { return r.cfg.Forge.Repository(ctx) },
		rules: func(branch string) ([]forge.Rule, error) { return r.cfg.Forge.EffectiveRules(ctx, branch) },
		probe: func(head string) (string, error) {
			a, err := r.auth(ctx)
			if err != nil {
				return "", err
			}
			return pushProbe(r.target, r.url(), a, head, r.who(r.cfg.Clock.Now()))
		},
		activity: func(ref string) error { _, err := r.cfg.Forge.Activity(ctx, ref); return err },
		commit: func(table []byte) error {
			return Fenced(ctx, r.store, r.cfg.RunID, func(ctx context.Context) error {
				return r.store.Commit(ctx, map[string][]byte{"setup/forge-check.tsv": table}, forgeCheckMessage)
			})
		},
	}.do()
}

// attribution says whether a refusal of the probe can be attributed to the default branch's ruleset: the rules of layup-probe
// were read, the apply output of that ruleset is in the records and gives its ID, and each creation rule of the ref is of it.
func attribution(rules []forge.Rule, rulesErr error, outs map[string][]byte, outsErr error) error {
	if rulesErr != nil {
		return fmt.Errorf("the rules of %s were not read: %s", probeBranch, firstLine(rulesErr.Error()))
	}
	if outsErr != nil {
		return outsErr
	}
	id, err := rulesetID(outs[applyFiles[0]])
	if err != nil {
		return err
	}
	return probeAttribution(rules, id)
}

// probeAttribution gives nil when the creation of layup-probe is restricted by the ruleset of ID id and by no other: the rules of the
// ref, as the forge gives them for a ref that does not exist yet, hold at least one rule of type creation, and each of them is of a
// ruleset of the repository with that ID. Only a creation rule refuses the creation of a ref (docs/spec/scaffold.md, The ruleset of
// the default branch), so a rule of another type of another ruleset is no fault. GitHub's refusal names the ref and the violated rule and
// no ruleset, so this is how a refusal is attributed to the ruleset of LAYUP.
func probeAttribution(rules []forge.Rule, id int64) error {
	if id <= 0 {
		return errors.New("the default branch's ruleset has no ID")
	}
	found := false
	for _, r := range rules {
		if r.Type != "creation" {
			continue
		}
		if r.RulesetSourceType != "Repository" || r.RulesetID != id {
			return fmt.Errorf("the creation of %s is also restricted by ruleset %d (%s %s)", probeBranch, r.RulesetID, r.RulesetSourceType, r.RulesetSource)
		}
		found = true
	}
	if !found {
		return fmt.Errorf("no creation rule applies to %s", probeBranch)
	}
	return nil
}

// rulesetID gives the ID of a ruleset from the output of its apply: the `id` of the one JSON object that the output is, a positive
// integer. An output that is anything but one object, with a second value or other text after it, is refused: the ID of the first of
// two objects would be the ID of a ruleset that no apply made.
func rulesetID(output []byte) (int64, error) {
	notObject := errors.New("the output of the default branch's ruleset is not one JSON object")
	noID := errors.New("the output of the default branch's ruleset has no id")
	dec := json.NewDecoder(bytes.NewReader(output))
	dec.UseNumber()
	var object map[string]any
	if err := dec.Decode(&object); err != nil || object == nil {
		return 0, notObject
	}
	if _, err := dec.Token(); err != io.EOF {
		return 0, notObject
	}
	n, ok := object["id"].(json.Number)
	if !ok {
		return 0, noID
	}
	id, err := n.Int64()
	if err != nil || id <= 0 {
		return 0, noID
	}
	return id, nil
}

// mergeMaps reads the merges of the tasks of the register: the merge commit of each event merged, by its task, and the merge
// commit of each task with a pull request and no event merged whose pull request the forge shows merged.
func mergeMaps(ctx context.Context, store interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
}, f forge.Forge) (recorded, unrecorded map[string]string, err error) {
	return nil, nil, nil
}
