package run

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
)

// The step forge-check (docs/spec/scaffold.md, REQ-003, Scaffold 7, with its paragraph "Decided here"
// of task T-wphp; the five input states of the forge check; the acceptance rows "The read-back" and
// "The probe"; row 57a of the plan, #280): the read-back of the rules, the bypass lists, the verdict of
// the probe, and the step with its seams as fakes. The table has eight rows in every case.

func rule(typ string, params map[string]any) forge.Rule {
	return forge.Rule{Type: typ, Parameters: params}
}

func pair(context string, id float64) map[string]any {
	return map[string]any{"context": context, "integration_id": id}
}

func gates(strict bool, pairs ...map[string]any) forge.Rule {
	list := make([]any, len(pairs))
	for i, p := range pairs {
		list[i] = p
	}
	return rule("required_status_checks", map[string]any{"strict_required_status_checks_policy": strict, "required_status_checks": list})
}

func TestTheReadBackComparesEachRuleOfTheFileWithOneReturnedRule(t *testing.T) {
	pr := func(m ...string) forge.Rule {
		list := make([]any, len(m))
		for i, s := range m {
			list[i] = s
		}
		return rule("pull_request", map[string]any{"allowed_merge_methods": list, "required_review_thread_resolution": true, "required_approving_review_count": float64(0)})
	}
	creation, deletion := rule("creation", nil), rule("deletion", nil)
	for _, c := range []struct {
		name      string
		want, got []forge.Rule
		ok        bool
		detail    string
	}{
		{"equal", []forge.Rule{creation, pr("merge")}, []forge.Rule{pr("merge"), creation}, true, ""},
		{"a rule that the file holds is missing", []forge.Rule{creation, deletion}, []forge.Rule{creation}, false, "deletion is missing"},
		{"the call returns none", []forge.Rule{creation}, nil, false, "creation is missing"},
		{"a parameter that differs", []forge.Rule{pr("merge")}, []forge.Rule{pr("squash")}, false, "pull_request: allowed_merge_methods differs"},
		{"a parameter of the file that the call lacks, false", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false})}, []forge.Rule{rule("update", nil)}, false, "update: update_allows_fetch_and_merge is missing"},
		{"a parameter of the file that the call lacks, zero", []forge.Rule{rule("x", map[string]any{"n": float64(0)})}, []forge.Rule{rule("x", map[string]any{})}, false, "x: n is missing"},
		{"a parameter of the call that the file lacks", []forge.Rule{creation, rule("update", map[string]any{})}, []forge.Rule{creation, rule("update", map[string]any{"update_allows_fetch_and_merge": false})}, true, ""},
		{"absent parameters and an empty object are alike", []forge.Rule{rule("a", nil), rule("b", map[string]any{})}, []forge.Rule{rule("a", map[string]any{}), rule("b", nil)}, true, ""},
		{"the same rule twice", []forge.Rule{creation}, []forge.Rule{creation, creation}, true, ""},
		{"two entries of a type that each hold a part of the parameters, one order",
			[]forge.Rule{rule("x", map[string]any{"a": float64(1), "b": float64(2)})}, []forge.Rule{rule("x", map[string]any{"a": float64(1)}), rule("x", map[string]any{"b": float64(2)})}, false, "x: b is missing"},
		{"two entries of a type that each hold a part of the parameters, the other order",
			[]forge.Rule{rule("x", map[string]any{"a": float64(1), "b": float64(2)})}, []forge.Rule{rule("x", map[string]any{"b": float64(2)}), rule("x", map[string]any{"a": float64(1)})}, false, "x: a is missing"},
		{"two parameters that differ: the first by name", []forge.Rule{rule("x", map[string]any{"b": float64(2), "a": float64(1)})}, []forge.Rule{rule("x", map[string]any{"a": float64(9), "b": float64(9)})}, false, "x: a differs"},
		{"two entries of a type, the second equal", []forge.Rule{pr("merge")}, []forge.Rule{pr("squash"), pr("merge")}, true, ""},
		{"the required checks in another order", []forge.Rule{gates(true, pair("a", 1), pair("b", 2))}, []forge.Rule{gates(true, pair("b", 2), pair("a", 1))}, true, ""},
		{"an integration_id that differs", []forge.Rule{gates(true, pair("layup/gates", 42))}, []forge.Rule{gates(true, pair("layup/gates", 15368))}, false, "required_status_checks: required_status_checks differs"},
		{"an integration_id of 0 that the call lacks", []forge.Rule{gates(true, pair("a", 0))}, []forge.Rule{gates(true, map[string]any{"context": "a"})}, false, "required_status_checks: required_status_checks differs"},
		{"an integration_id that the file lacks and the call holds as 0", []forge.Rule{gates(true, map[string]any{"context": "a"})}, []forge.Rule{gates(true, pair("a", 0))}, false, "required_status_checks: required_status_checks differs"},
		{"a context that differs", []forge.Rule{gates(true, pair("a", 1))}, []forge.Rule{gates(true, pair("b", 1))}, false, "required_status_checks: required_status_checks differs"},
		{"an extra required check", []forge.Rule{gates(true, pair("a", 1))}, []forge.Rule{gates(true, pair("a", 1), pair("b", 2))}, false, "required_status_checks"},
		{"a required check missing", []forge.Rule{gates(true, pair("a", 1), pair("b", 2))}, []forge.Rule{gates(true, pair("a", 1))}, false, "required_status_checks"},
		{"strict that differs", []forge.Rule{gates(true, pair("a", 1))}, []forge.Rule{gates(false, pair("a", 1))}, false, "required_status_checks: strict_required_status_checks_policy differs"},
		{"an extra rule is named and passes", []forge.Rule{creation}, []forge.Rule{rule("required_signatures", nil), creation}, true, "extra: required_signatures"},
		{"extra rules are sorted", []forge.Rule{creation}, []forge.Rule{rule("workflows", nil), creation, rule("required_signatures", nil), rule("workflows", nil)}, true, "extra: required_signatures, workflows"},
	} {
		detail, ok := readBack(c.want, c.got)
		if ok != c.ok || !strings.Contains(detail, c.detail) || (ok && detail != c.detail) {
			t.Errorf("%s: %v, %q; want %v and %q", c.name, ok, detail, c.ok, c.detail)
		}
	}
}

func TestARulesetFileIsReadAsItsRulesInOrder(t *testing.T) {
	got, err := rulesetRules([]byte(`{"name":"x","rules":[{"type":"creation"},{"type":"update","parameters":{"update_allows_fetch_and_merge":false}}],"bypass_actors":[]}`))
	if err != nil || len(got) != 2 || got[0].Type != "creation" || got[1].Type != "update" || got[1].Parameters["update_allows_fetch_and_merge"] != false {
		t.Errorf("the rules of a ruleset file: %v, %v; want creation, then update with its parameter", got, err)
	}
	for _, c := range []struct{ name, text, reason string }{{"no rules", `{"name":"x"}`, "it holds no rule"}, {"an empty list", `{"rules":[]}`, "it holds no rule"},
		{"not JSON", `{"rules":`, "it is not a JSON object"}, {"an array", `[]`, "it is not a JSON object"}, {"empty", ``, "it is not a JSON object"}} {
		if got, err := rulesetRules([]byte(c.text)); err == nil || err.Error() != c.reason || got != nil {
			t.Errorf("%s: %v, %v; want the error %q and no rule", c.name, got, err, c.reason)
		}
	}
}

func TestTheOutputsOfTheAppliesAreReadBackByTheirFraming(t *testing.T) {
	const d, r = `{"source":"acme/target","bypass_actors":[]}`, `{"source":"acme/target","bypass_actors":[1]}`
	for _, c := range []struct{ name, d, r string }{
		{"neither ends with a line feed", d, r},
		{"each ends with a line feed", d + "\n", r + "\n"},
		{"the first does", d + "\n", r},
		{"the second does", d, r + "\n"},
	} {
		record, err := rulesetsRecord("acme/target", map[string][]byte{"default.json": []byte(c.d), "records.json": []byte(c.r)})
		if err != nil {
			t.Fatal(err)
		}
		got, err := readApplied(record)
		if err != nil || string(got["default.json"]) != c.d || string(got["records.json"]) != c.r || len(got) != 2 {
			t.Errorf("%s: %q, %v; want each output byte for byte", c.name, got, err)
		}
	}
	for _, c := range []struct{ name, text string }{
		{"empty", ""}, {"the first marker missing", "{}\n== records.json\n{}\n\n"}, {"the second marker missing", "== default.json\n{}\n\n"},
		{"no separator at the end", "== default.json\n{}\n== records.json\n{}"}, {"the markers in the other order", "== records.json\n{}\n== default.json\n{}\n"},
		{"a first marker that is another text of its length", "== defaultXjson\n{}\n== records.json\n{}\n\n"}, {"nothing after the second marker", "== default.json\n{}\n== records.json\n"},
		{"the first output and its separator missing", "== default.json\n== records.json\n{}\n\n"},
	} {
		if got, err := readApplied([]byte(c.text)); err == nil || got != nil {
			t.Errorf("%s: %q, %v; want an error and no output", c.name, got, err)
		}
	}
}

func TestTheBypassListsAreEmptyForTheDefaultAndTheAppAloneForTheRecords(t *testing.T) {
	app := `{"actor_id":42,"actor_type":"Integration","bypass_mode":"always"}`
	for _, c := range []struct {
		name, text string
		ok         bool
	}{
		{"an empty list", `{"bypass_actors":[]}`, true},
		{"a list that is not empty", `{"bypass_actors":[` + app + `]}`, false},
		{"no list", `{"id":1}`, false},
		{"not JSON", `{`, false},
	} {
		if err := bypassDefault([]byte(c.text)); (err == nil) != c.ok {
			t.Errorf("default, %s: %v; want ok %v", c.name, err, c.ok)
		}
	}
	for _, c := range []struct {
		name, text string
		ok         bool
	}{
		{"the App alone", `{"bypass_actors":[` + app + `]}`, true},
		{"no actor", `{"bypass_actors":[]}`, false},
		{"another actor", `{"bypass_actors":[{"actor_id":7,"actor_type":"Integration","bypass_mode":"always"}]}`, false},
		{"another type", `{"bypass_actors":[{"actor_id":42,"actor_type":"RepositoryRole","bypass_mode":"always"}]}`, false},
		{"another mode", `{"bypass_actors":[{"actor_id":42,"actor_type":"Integration","bypass_mode":"pull_request"}]}`, false},
		{"the App and another", `{"bypass_actors":[` + app + `,{"actor_id":1,"actor_type":"OrganizationAdmin","bypass_mode":"always"}]}`, false},
		{"no list", `{}`, false},
	} {
		if err := bypassRecords([]byte(c.text), 42); (err == nil) != c.ok {
			t.Errorf("records, %s: %v; want ok %v", c.name, err, c.ok)
		}
	}
}

func TestTheVerdictOfTheProbeIsAPassOnlyForTheRefusalOfTheRulesetByName(t *testing.T) {
	const sha = "0123456789012345678901234567890123456789"
	refused := func(stderr string) error {
		return &git.FailedError{Args: []string{"push"}, Code: 1, Stderr: stderr, Err: errors.New("exit status 1")}
	}
	for _, c := range []struct {
		name, result, detail string
		err                  error
	}{
		{"a refusal that names the ruleset", "pass", "", refused("remote: error: GH013: Repository rule violations found for refs/heads/layup-probe.\nremote: - Rulesets: layup: the default branch\n")},
		{"a push that the forge accepts", "fail", "probe accepted: " + sha, nil},
		{"a refusal by the ruleset of the records branch", "fail", "probe not run: remote: error: GH013", refused("remote: error: GH013: Repository rule violations found\nremote: - Rulesets: layup: the records branch\n")},
		{"a refusal by another ruleset", "fail", "probe not run: remote: error: GH013: Repository rule violations found", refused("remote: error: GH013: Repository rule violations found\nremote: - Rulesets: another\n")},
		{"a refusal with another text", "fail", "probe not run: remote: error: denied", refused("remote: error: denied\nlayup: the default branch\n")},
		{"a remote that cannot be reached", "fail", "probe not run: fatal: unable to access", &git.FailedError{Args: []string{"push"}, Code: 128, Stderr: "fatal: unable to access 'x'\n", Err: errors.New("exit status 128")}},
		{"an error that is no refusal of git", "fail", "probe not run: the token expired", errors.New("the token expired\nmore")},
		{"a call that did not start, whose text holds both", "fail", "probe not run: remote: error: GH013", &git.FailedError{Args: []string{"push"}, Code: -1, Stderr: "remote: error: GH013\nlayup: the default branch\n", Err: errors.New("not started")}},
		{"a failure with no standard error", "fail", "probe not run: git push: exit status 1", &git.FailedError{Args: []string{"push"}, Code: 1, Err: errors.New("exit status 1")}},
	} {
		result, detail := probeVerdict(sha, c.err)
		if result != c.result || !strings.HasPrefix(detail, c.detail) || (result == "pass" && detail != "") {
			t.Errorf("%s: %s, %q; want %s, %q", c.name, result, detail, c.result, c.detail)
		}
	}
}

// checkFake is the world of a job: all eight checks pass until a case breaks one.
type checkFake struct {
	treeText, applied string
	recordsErr        error
	areaErr           error
	files             map[string][]byte
	repo              forge.Repository
	repoErr           error
	rules             map[string][]forge.Rule
	rulesErr          map[string]error
	probeErr          error
	activityErr       error
	commitErr         error
	committed         [][]byte
	clock             time.Time
	areaFor           string
	probedHead        string
	activityRef       string
	areaCalls         int
}

const lastCommit = "0123456789012345678901234567890123456789"

func newCheckFake() *checkFake {
	var t bytes.Buffer
	tsv.Write(&t, SetupTreeSchema, [][]string{{"0", "1111111111111111111111111111111111111111", "2222222222222222222222222222222222222222", "root"}, {"1", lastCommit, "3333333333333333333333333333333333333333", "S13"}})
	defaultFile := `{"rules":[{"type":"creation"},{"type":"deletion"},{"type":"required_status_checks","parameters":{"strict_required_status_checks_policy":true,"required_status_checks":[{"context":"layup/gates","integration_id":42}]}}],"bypass_actors":[]}`
	recordsFile := `{"rules":[{"type":"update","parameters":{"update_allows_fetch_and_merge":false}},{"type":"deletion"}],"bypass_actors":[]}`
	rec, _ := rulesetsRecord("acme/target", map[string][]byte{"default.json": []byte(`{"source":"acme/target","bypass_actors":[]}`), "records.json": []byte(`{"source":"acme/target","bypass_actors":[{"actor_id":42,"actor_type":"Integration","bypass_mode":"always"}]}`)})
	def, _ := rulesetRules([]byte(defaultFile))
	recs, _ := rulesetRules([]byte(recordsFile))
	return &checkFake{treeText: t.String(), applied: string(rec), files: map[string][]byte{"ruleset-default.json": []byte(defaultFile), "ruleset-records.json": []byte(recordsFile)},
		repo:  forge.Repository{DefaultBranch: "main", Heads: map[string]string{"main": lastCommit}},
		rules: map[string][]forge.Rule{"main": def, "layup-probe": def, "layup-records": recs}, rulesErr: map[string]error{},
		probeErr: &git.FailedError{Args: []string{"push"}, Code: 1, Stderr: "remote: error: GH013: Repository rule violations found\nremote: - Rulesets: layup: the default branch\n", Err: errors.New("exit status 1")},
		clock:    time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC)}
}

func (f *checkFake) job() forgeCheckJob {
	return forgeCheckJob{
		branch: "main", app: 42, run: "0123456789abcdef",
		now: func() time.Time { f.clock = f.clock.Add(time.Second); return f.clock },
		records: func() ([]byte, []byte, error) {
			var tree, applied []byte
			if f.treeText != "" {
				tree = []byte(f.treeText)
			}
			if f.applied != "" {
				applied = []byte(f.applied)
			}
			return tree, applied, f.recordsErr
		},
		area: func(last string) error { f.areaFor, f.areaCalls = last, f.areaCalls+1; return f.areaErr },
		read: func(name string) ([]byte, error) {
			if data, ok := f.files[name]; ok {
				return data, nil
			}
			return nil, fmt.Errorf("%s: no such file", name)
		},
		repo:     func() (forge.Repository, error) { return f.repo, f.repoErr },
		rules:    func(branch string) ([]forge.Rule, error) { return f.rules[branch], f.rulesErr[branch] },
		probe:    func(head string) (string, error) { f.probedHead = head; return lastCommit, f.probeErr },
		activity: func(ref string) error { f.activityRef = ref; return f.activityErr },
		clean:    func(s string) string { return s },
		commit:   func(table []byte) error { f.committed = append(f.committed, table); return f.commitErr },
	}
}

// rowsOf reads the table that the job committed, by the reader of the block.
func (f *checkFake) rowsOf(t *testing.T) [][]string {
	t.Helper()
	if len(f.committed) != 1 {
		t.Fatalf("%d commits of the table; want one", len(f.committed))
	}
	rows, err := ReadForgeCheck(f.committed[0])
	if err != nil {
		t.Fatalf("the table does not read back by its schema: %v\n%s", err, f.committed[0])
	}
	return rows
}

func TestTheStepForgeCheckCommitsEightRowsAndIsDoneWhenEachIsAPass(t *testing.T) {
	f := newCheckFake()
	f.rules["main"] = append(f.rules["main"], rule("required_signatures", nil))
	got := f.job().do()
	if got != (Step{"forge-check", "done", "the forge check passed"}) {
		t.Fatalf("the step %+v; want done", got)
	}
	rows := f.rowsOf(t)
	names := []string{"tree", "rules:main", "rules:layup-probe", "rules:layup-records", "bypass:default", "bypass:records", "probe", "activity"}
	if len(rows) != len(names) || f.areaFor != lastCommit || f.probedHead != lastCommit || f.activityRef != "main" {
		t.Fatalf("%d rows, the area for %q, the probe on %q, the activity of %q; want eight rows, the last commit for both, and main", len(rows), f.areaFor, f.probedHead, f.activityRef)
	}
	for i, r := range rows {
		wantDetail := ""
		if names[i] == "rules:main" {
			wantDetail = "extra: required_signatures"
		}
		if want := []string{names[i], "pass", wantDetail, "0123456789abcdef", fmt.Sprintf("2026-10-13T12:00:%02dZ", i+1)}; strings.Join(r, "|") != strings.Join(want, "|") {
			t.Errorf("row %d: %q; want %q", i, r, want)
		}
	}
}

func TestEachCheckThatFailsIsItsRowAndTheTableStillHasEightRows(t *testing.T) {
	boom := errors.New("boom\nthe second line")
	for _, c := range []struct {
		name     string
		change   func(f *checkFake)
		fail     map[string]string // the rows that fail, by check, with a part of the detail; the others pass
		stepWord string
	}{
		{"a head that is not the last commit", func(f *checkFake) { f.repo.Heads["main"] = "abcdefabcdefabcdefabcdefabcdefabcdefabcd" },
			map[string]string{"tree": "abcdefabcdefabcdefabcdefabcdefabcdefabcd"}, "tree: "},
		{"a default branch that the forge lacks", func(f *checkFake) { delete(f.repo.Heads, "main") },
			map[string]string{"tree": "none", "probe": "probe not run: the forge gives no head"}, "tree: "},
		{"a read of the repository that fails", func(f *checkFake) { f.repoErr = boom },
			map[string]string{"tree": "the forge: boom", "probe": "probe not run: the head of main was not read"}, "tree: "},
		{"a rule that is missing on the default branch", func(f *checkFake) { f.rules["main"] = f.rules["main"][1:] },
			map[string]string{"rules:main": "creation is missing"}, "rules:main: "},
		{"a rule that is missing on layup-probe", func(f *checkFake) { f.rules["layup-probe"] = nil },
			map[string]string{"rules:layup-probe": "creation is missing"}, "rules:layup-probe: "},
		{"a rule that is missing on layup-records", func(f *checkFake) { f.rules["layup-records"] = f.rules["layup-records"][:1] },
			map[string]string{"rules:layup-records": "deletion is missing"}, "rules:layup-records: "},
		{"a call of rules that fails", func(f *checkFake) { f.rulesErr["layup-probe"] = boom },
			map[string]string{"rules:layup-probe": "boom"}, "rules:layup-probe: "},
		{"a ruleset file that is missing", func(f *checkFake) { delete(f.files, "ruleset-records.json") },
			map[string]string{"rules:layup-records": "ruleset-records.json: no such file"}, "rules:layup-records: "},
		{"a ruleset file that holds no rule", func(f *checkFake) { f.files["ruleset-default.json"] = []byte(`{"rules":[]}`) },
			map[string]string{"rules:main": "ruleset-default.json: it holds no rule", "rules:layup-probe": "ruleset-default.json: it holds no rule"}, "rules:main: "},
		{"a work area that cannot be made sure of", func(f *checkFake) { f.areaErr = errors.New("layup-setup is at x, and setup/tree.tsv ends at y") },
			map[string]string{"rules:main": "the work area: layup-setup is at x", "rules:layup-probe": "the work area", "rules:layup-records": "the work area"}, "rules:main: "},
		{"a setup/tree.tsv that is not in the records", func(f *checkFake) { f.treeText = "" },
			map[string]string{"tree": "setup/tree.tsv is not in the records", "rules:main": "setup/tree.tsv is not in the records", "rules:layup-probe": "setup/tree.tsv is not in the records", "rules:layup-records": "setup/tree.tsv is not in the records"}, "tree: "},
		{"a setup/tree.tsv that its reader refuses", func(f *checkFake) { f.treeText = "order\tcommit\n" },
			map[string]string{"tree": "setup/tree.tsv: line 1", "rules:main": "setup/tree.tsv: line 1", "rules:layup-probe": "setup/tree.tsv: line 1", "rules:layup-records": "setup/tree.tsv: line 1"}, "tree: "},
		{"a setup/rulesets.txt that is not in the records", func(f *checkFake) { f.applied = "" },
			map[string]string{"bypass:default": "setup/rulesets.txt is not in the records", "bypass:records": "setup/rulesets.txt is not in the records"}, "bypass:default: "},
		{"a setup/rulesets.txt that does not read back", func(f *checkFake) { f.applied = "not framed" },
			map[string]string{"bypass:default": "setup/rulesets.txt", "bypass:records": "setup/rulesets.txt"}, "bypass:default: "},
		{"a bypass list of the default that is not empty", func(f *checkFake) {
			f.applied = strings.Replace(f.applied, `"source":"acme/target","bypass_actors":[]`, `"source":"acme/target","bypass_actors":[{"actor_id":1}]`, 1)
		}, map[string]string{"bypass:default": "not empty"}, "bypass:default: "},
		{"a bypass list of the records that is not the App", func(f *checkFake) { f.app(7) }, map[string]string{"bypass:records": "App"}, "bypass:records: "},
		{"a probe that the forge accepts", func(f *checkFake) { f.probeErr = nil }, map[string]string{"probe": "probe accepted: " + lastCommit}, "probe: "},
		{"a probe that fails for another reason", func(f *checkFake) { f.probeErr = boom }, map[string]string{"probe": "probe not run: boom"}, "probe: "},
		{"an activity read that fails", func(f *checkFake) { f.activityErr = boom }, map[string]string{"activity": "boom"}, "activity: "},
		{"a ruleset of the forge that is in evaluate (no rule is returned)", func(f *checkFake) { f.rules = map[string][]forge.Rule{} },
			map[string]string{"rules:main": "missing", "rules:layup-probe": "missing", "rules:layup-records": "missing"}, "rules:main: "},
	} {
		f := newCheckFake()
		c.change(f)
		got := f.job().do()
		rows := f.rowsOf(t)
		if len(rows) != 8 {
			t.Fatalf("%s: %d rows; want eight", c.name, len(rows))
		}
		for _, r := range rows {
			want, failing := c.fail[r[0]]
			switch {
			case failing && (r[1] != "fail" || !strings.Contains(r[2], want)):
				t.Errorf("%s: the row %s is %s, %q; want fail with %q", c.name, r[0], r[1], r[2], want)
			case !failing && r[1] != "pass":
				t.Errorf("%s: the row %s is %s, %q; want pass", c.name, r[0], r[1], r[2])
			}
		}
		// The probe is made on the head that the forge gives, and the work area is made sure of only for a last commit.
		wantHead := f.repo.Heads["main"]
		if f.repoErr != nil {
			wantHead = ""
		}
		noLast := f.treeText == "" || strings.HasPrefix(f.treeText, "order\tcommit\n") // no last commit to make the area sure of
		if f.probedHead != wantHead || (f.areaCalls == 0) != noLast {
			t.Errorf("%s: the probe on %q, %d calls of the area; want the head %q, and the area for a last commit only", c.name, f.probedHead, f.areaCalls, wantHead)
		}
		if got.Name != "forge-check" || got.Result != "fail" || !strings.HasPrefix(got.Detail, c.stepWord) {
			t.Errorf("%s: the step %+v; want fail, naming %q first", c.name, got, c.stepWord)
		}
	}
}

// app changes the App's ID in the apply record of the records branch.
func (f *checkFake) app(id int) {
	f.applied = strings.Replace(f.applied, `"actor_id":42`, fmt.Sprintf(`"actor_id":%d`, id), 1)
}

func TestOnlyAnUnavailableRecordsBaseOrAFailedCommitEndsTheStepWithNoTable(t *testing.T) {
	boom := errors.New("boom")
	f := newCheckFake()
	f.recordsErr = boom
	if got := f.job().do(); got.Result != "fail" || !strings.Contains(got.Detail, "no table was committed") || len(f.committed) != 0 {
		t.Errorf("an unavailable records base: %+v, %d commits; want fail, no table was committed, and no commit", got, len(f.committed))
	}
	f = newCheckFake()
	f.commitErr = boom
	got := f.job().do()
	if got.Result != "fail" || !strings.Contains(got.Detail, "the records commit: boom") || !strings.Contains(got.Detail, "no table was committed") || len(f.committed) != 1 {
		t.Errorf("a commit that fails: %+v, %d commits; want fail, the commit that was attempted, and no table was committed", got, len(f.committed))
	}
	f = newCheckFake()
	f.repo.Heads["main"] = "abcdefabcdefabcdefabcdefabcdefabcdefabcd"
	f.probeErr = boom
	f.commitErr = boom
	if got := f.job().do(); !strings.Contains(got.Detail, "the records commit: boom") {
		t.Errorf("a check that fails and a commit that fails: %+v; want the commit named, as the table was not committed", got)
	}
}

func TestADetailIsTheFirstLineOfItsReasonInValidUTF8WithDIRForTheHost(t *testing.T) {
	for _, c := range []struct{ reason, want string }{{"boom\nthe second line", "boom"}, {"bad \xff byte", "bad ? byte"}, {"open /srv/layup-host/x: denied", "open DIR/x: denied"}} {
		f := newCheckFake()
		f.activityErr = errors.New(c.reason)
		job := f.job()
		job.clean = strings.NewReplacer("/srv/layup-host", "DIR").Replace
		if got := job.do(); got != (Step{"forge-check", "fail", "activity: " + c.want}) || f.rowsOf(t)[7][2] != c.want {
			t.Errorf("a reason %q: the step %+v and the row %q; want the detail %q", c.reason, got, f.rowsOf(t)[7], c.want)
		}
	}
}

func TestATableThatItsReaderWouldRefuseIsNotCommitted(t *testing.T) {
	f := newCheckFake()
	job := f.job()
	job.run = "not a run ID"
	if got := job.do(); got.Result != "fail" || !strings.HasPrefix(got.Detail, "the table: ") || !strings.Contains(got.Detail, "no table was committed") || len(f.committed) != 0 {
		t.Errorf("a run that is no run ID: %+v, %d commits; want the table refused by its reader, and no commit", got, len(f.committed))
	}
}
