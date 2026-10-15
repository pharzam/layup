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
		{"a parameter that GitHub omits at its default: update_allows_fetch_and_merge false", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false})}, []forge.Rule{rule("update", nil)}, true, ""},
		{"the same parameter true in the file, omitted by the call", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": true})}, []forge.Rule{rule("update", nil)}, false, "update: update_allows_fetch_and_merge is missing"},
		{"the same parameter false in the file, true returned", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false})}, []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": true})}, false, "update: update_allows_fetch_and_merge differs"},
		{"the same parameter true in the file, false returned", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": true})}, []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false})}, false, "update: update_allows_fetch_and_merge differs"},
		{"the same name, false and omitted, on a rule of another type", []forge.Rule{rule("x", map[string]any{"update_allows_fetch_and_merge": false})}, []forge.Rule{rule("x", nil)}, false, "x: update_allows_fetch_and_merge is missing"},
		{"another parameter of update, false and omitted", []forge.Rule{rule("update", map[string]any{"other": false})}, []forge.Rule{rule("update", nil)}, false, "update: other is missing"},
		{"the omitted default and another missing parameter, named after it", []forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false, "z_other": true})}, []forge.Rule{rule("update", nil)}, false, "update: z_other is missing"},
		{"update with no parameters on both sides", []forge.Rule{rule("update", nil)}, []forge.Rule{rule("update", nil)}, true, ""},
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

// GitHub's refusal for a rule violation is the line `GH013: Repository rule violations found for refs/heads/<ref>.`, the
// line of the rules URL and one line for each violated rule; it names no ruleset (ten public reports, issue #323). The
// verdict passes it for the probe's ref alone, as the whole ref, and only when the attribution has no fault.
func TestTheVerdictOfTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRefThatARulesetExplains(t *testing.T) {
	const sha = "0123456789012345678901234567890123456789"
	const notAttributed = "probe refused, but not attributed to the default branch's ruleset: "
	refused := func(stderr string) error {
		return &git.FailedError{Args: []string{"push"}, Code: 1, Stderr: stderr, Err: errors.New("exit status 1")}
	}
	text := func(ref string) string {
		return "remote: error: GH013: Repository rule violations found for " + ref + ".\nremote: Review all repository rules at https://github.com/acme/target/rules?ref=" + ref +
			"\nremote: \nremote: - Cannot create ref due to creations being restricted.\nremote: \n"
	}
	fault := errors.New("the creation of layup-probe is also restricted by ruleset 9 (Organization acme)")
	for _, c := range []struct {
		name, result, detail string
		err, attributed      error
	}{
		{"GitHub's refusal of the probe ref, attributed", "pass", "", refused(text("refs/heads/layup-probe")), nil},
		{"the same with a line that names a ruleset, as the tests used to invent", "pass", "", refused(text("refs/heads/layup-probe") + "remote: - Rulesets: layup: the default branch\n"), nil},
		{"the line without its final period", "pass", "", refused("remote: error: GH013: Repository rule violations found for refs/heads/layup-probe\n"), nil},
		{"the line that ends with a carriage return", "pass", "", refused("remote: error: GH013: Repository rule violations found for refs/heads/layup-probe.\r\n"), nil},
		{"a refusal for a ref that starts with the probe's name", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/layup-probe-other.", refused(text("refs/heads/layup-probe-other")), nil},
		{"a refusal for a ref below the probe's name", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/layup-probe/child.", refused(text("refs/heads/layup-probe/child")), nil},
		{"a refusal for a ref with an extension", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/layup-probe.x.", refused(text("refs/heads/layup-probe.x")), nil},
		{"a refusal for the default branch", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/main.", refused(text("refs/heads/main")), nil},
		{"the line with a word after the ref", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/layup-probe for", refused("remote: error: GH013: Repository rule violations found for refs/heads/layup-probe for the user.\n"), nil},
		{"a refusal with no GH013", "fail", "probe not run: remote: error: denied", refused("remote: error: denied\nremote: Repository rule violations found for refs/heads/layup-probe.\n"), nil},
		{"a refusal that the attribution does not explain", "fail", notAttributed + fault.Error(), refused(text("refs/heads/layup-probe")), fault},
		{"a fault of the attribution and a refusal for another ref: the text decides first", "fail", "probe not run: remote: error: GH013: Repository rule violations found for refs/heads/main.", refused(text("refs/heads/main")), fault},
		{"a push that the forge accepts, with a fault", "fail", "probe accepted: " + sha, nil, fault},
		{"a push that the forge accepts", "fail", "probe accepted: " + sha, nil, nil},
		{"a remote that cannot be reached", "fail", "probe not run: fatal: unable to access", &git.FailedError{Args: []string{"push"}, Code: 128, Stderr: "fatal: unable to access 'x'\n", Err: errors.New("exit status 128")}, fault},
		{"an error that is no refusal of git", "fail", "probe not run: the token expired", errors.New("the token expired\nmore"), nil},
		{"a call that did not start, whose text matches", "fail", "probe not run: remote: error: GH013", &git.FailedError{Args: []string{"push"}, Code: -1, Stderr: text("refs/heads/layup-probe"), Err: errors.New("not started")}, nil},
		{"a failure with no standard error", "fail", "probe not run: git push: exit status 1", &git.FailedError{Args: []string{"push"}, Code: 1, Err: errors.New("exit status 1")}, nil},
	} {
		result, detail := probeVerdict(sha, c.err, c.attributed)
		if result != c.result || !strings.HasPrefix(detail, c.detail) || (result == "pass" && detail != "") {
			t.Errorf("%s: %s, %q; want %s, %q", c.name, result, detail, c.result, c.detail)
		}
	}
}

// The attribution of a refusal: every creation rule that the effective rules of layup-probe list is the default branch's
// ruleset's, by its source type and its ID, and there is one at least; a rule of another type of another ruleset is no fault,
// as only a creation rule refuses the creation of a ref that does not exist.
func TestAProbeIsAttributedOnlyWhenEachCreationRuleOfTheRefIsTheDefaultRulesets(t *testing.T) {
	of := func(typ, source string, id int64) forge.Rule {
		return forge.Rule{Type: typ, RulesetSourceType: "Repository", RulesetSource: source, RulesetID: id}
	}
	other := func(typ, sourceType string, id int64) forge.Rule {
		return forge.Rule{Type: typ, RulesetSourceType: sourceType, RulesetSource: "acme", RulesetID: id}
	}
	for _, c := range []struct {
		name  string
		rules []forge.Rule
		id    int64
		fault string // a part of the error, or "" for none
	}{
		{"the default's creation rule alone", []forge.Rule{of("creation", "acme/target", 77)}, 77, ""},
		{"the default's rules", []forge.Rule{of("creation", "acme/target", 77), of("deletion", "acme/target", 77), of("non_fast_forward", "acme/target", 77)}, 77, ""},
		{"a rule of another type of another ruleset", []forge.Rule{of("creation", "acme/target", 77), other("pull_request", "Organization", 9)}, 77, ""},
		{"the creation rule twice, as two entries of the default", []forge.Rule{of("creation", "acme/target", 77), of("creation", "acme/target", 77)}, 77, ""},
		{"a second ruleset's creation rule", []forge.Rule{of("creation", "acme/target", 77), other("creation", "Organization", 9)}, 77, "the creation of layup-probe is also restricted by ruleset 9 (Organization acme)"},
		{"a second ruleset's creation rule first", []forge.Rule{other("creation", "Repository", 8), of("creation", "acme/target", 77)}, 77, "ruleset 8 (Repository acme)"},
		{"the default's ID of another source type", []forge.Rule{other("creation", "Organization", 77)}, 77, "ruleset 77 (Organization acme)"},
		{"a repository ruleset of another ID", []forge.Rule{of("creation", "acme/target", 78)}, 77, "ruleset 78 (Repository acme/target)"},
		{"no creation rule", []forge.Rule{of("deletion", "acme/target", 77), of("non_fast_forward", "acme/target", 77)}, 77, "no creation rule applies to layup-probe"},
		{"no rule", nil, 77, "no creation rule applies to layup-probe"},
		{"a creation rule of no ruleset, for an ID of zero", []forge.Rule{{Type: "creation"}}, 0, "no ID"},
		{"an ID that is not positive", []forge.Rule{of("creation", "acme/target", -3)}, -3, "no ID"},
	} {
		err := probeAttribution(c.rules, c.id)
		if (err == nil) != (c.fault == "") || (err != nil && !strings.Contains(err.Error(), c.fault)) {
			t.Errorf("%s: %v; want a fault with %q (none when empty)", c.name, err, c.fault)
		}
	}
}

// The ID of the default branch's ruleset is the `id` of the output of its apply, a positive integer of GitHub's object.
func TestTheIDOfTheDefaultRulesetIsThePositiveIntegerIDOfItsApplyOutput(t *testing.T) {
	for _, c := range []struct {
		name, output string
		id           int64
		fault        string
	}{
		{"an output with its id", `{"id":24509051,"name":"layup: the default branch","source":"acme/target"}`, 24509051, ""},
		{"an output with a line feed after it", "{\"id\":77}\n", 77, ""},
		{"an output with spaces round it", " \t{\"id\":77}  \n\n", 77, ""},
		{"an id beyond the exact range of a float", `{"id":9007199254740993}`, 9007199254740993, ""},
		{"two objects, the first with the id", `{"id":77}{"id":99}`, 0, "not one JSON object"},
		{"two objects with a line feed between", "{\"id\":77}\n{\"id\":99}", 0, "not one JSON object"},
		{"text after the object", `{"id":77} x`, 0, "not one JSON object"},
		{"the start of a second object", `{"id":77}{`, 0, "not one JSON object"},
		{"a value after the object", `{"id":77} 5`, 0, "not one JSON object"},
		{"an output with no id", `{"source":"acme/target"}`, 0, "no id"},
		{"an id that is a string", `{"id":"77"}`, 0, "no id"},
		{"an id of zero", `{"id":0}`, 0, "no id"},
		{"an id below zero", `{"id":-5}`, 0, "no id"},
		{"an id that is a fraction", `{"id":1.5}`, 0, "no id"},
		{"an id of null", `{"id":null}`, 0, "no id"},
		{"a bare null", `null`, 0, "not one JSON object"},
		{"a list", `[{"id":77}]`, 0, "not one JSON object"},
		{"a text that is no JSON", `not json`, 0, "not one JSON object"},
		{"an empty output", ``, 0, "not one JSON object"},
	} {
		id, err := rulesetID([]byte(c.output))
		if id != c.id || (err == nil) != (c.fault == "") || (err != nil && !strings.Contains(err.Error(), c.fault)) {
			t.Errorf("%s: %d, %v; want %d and a fault with %q (none when empty)", c.name, id, err, c.id, c.fault)
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
	rulesRead         map[string]int // the reads of the effective rules of each branch
}

// The IDs that the apply outputs of the fake give the two rulesets, and the rules as the forge returns them for a ruleset.
const (
	defaultRulesetID = 77
	recordsRulesetID = 78
)

// asGitHubGives gives the rules of a ruleset file as the real forge returns them: the rule update with no parameters, though the file wrote
// update_allows_fetch_and_merge false (first hand, run 6 of the demo of M2d, 2026-10-15; task T-j7c5, #330; runs/T-j7c5/test-runs.md).
func asGitHubGives(rules []forge.Rule) []forge.Rule {
	out := make([]forge.Rule, len(rules))
	for i, r := range rules {
		if r.Type == "update" && r.Parameters["update_allows_fetch_and_merge"] == false {
			r.Parameters = nil
		}
		out[i] = r
	}
	return out
}

func ofRuleset(rules []forge.Rule, id int64) []forge.Rule {
	out := make([]forge.Rule, len(rules))
	for i, r := range rules {
		r.RulesetSourceType, r.RulesetSource, r.RulesetID = "Repository", "acme/target", id
		out[i] = r
	}
	return out
}

// GitHub's refusal of the probe: the line of GH013 with the ref, the line of the rules URL, one line for the rule.
const githubRefusal = "remote: error: GH013: Repository rule violations found for refs/heads/layup-probe.\nremote: Review all repository rules at https://github.com/acme/target/rules?ref=refs/heads/layup-probe\nremote: \nremote: - Cannot create ref due to creations being restricted.\nremote: \n"

const lastCommit = "0123456789012345678901234567890123456789"

func newCheckFake() *checkFake {
	var t bytes.Buffer
	tsv.Write(&t, SetupTreeSchema, [][]string{{"0", "1111111111111111111111111111111111111111", "2222222222222222222222222222222222222222", "root"}, {"1", lastCommit, "3333333333333333333333333333333333333333", "S13"}})
	defaultFile := `{"rules":[{"type":"creation"},{"type":"deletion"},{"type":"required_status_checks","parameters":{"strict_required_status_checks_policy":true,"required_status_checks":[{"context":"layup/gates","integration_id":42}]}}],"bypass_actors":[]}`
	recordsFile := `{"rules":[{"type":"update","parameters":{"update_allows_fetch_and_merge":false}},{"type":"deletion"}],"bypass_actors":[]}`
	rec, _ := rulesetsRecord("acme/target", map[string][]byte{"default.json": []byte(`{"id":77,"source":"acme/target","bypass_actors":[]}`), "records.json": []byte(`{"id":78,"source":"acme/target","bypass_actors":[{"actor_id":42,"actor_type":"Integration","bypass_mode":"always"}]}`)})
	def, _ := rulesetRules([]byte(defaultFile))
	recs, _ := rulesetRules([]byte(recordsFile))
	return &checkFake{treeText: t.String(), applied: string(rec), files: map[string][]byte{"ruleset-default.json": []byte(defaultFile), "ruleset-records.json": []byte(recordsFile)},
		repo:  forge.Repository{DefaultBranch: "main", Heads: map[string]string{"main": lastCommit}},
		rules: map[string][]forge.Rule{"main": ofRuleset(def, defaultRulesetID), "layup-probe": ofRuleset(def, defaultRulesetID), "layup-records": ofRuleset(asGitHubGives(recs), recordsRulesetID)}, rulesErr: map[string]error{}, rulesRead: map[string]int{},
		probeErr: &git.FailedError{Args: []string{"push"}, Code: 1, Stderr: githubRefusal, Err: errors.New("exit status 1")},
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
		repo: func() (forge.Repository, error) { return f.repo, f.repoErr },
		rules: func(branch string) ([]forge.Rule, error) {
			f.rulesRead[branch]++
			return f.rules[branch], f.rulesErr[branch]
		},
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

// The fake of the job gives the rules of layup-records as the real forge does (#330): the rule update with no parameters, while the ruleset
// file that the apply sent holds update_allows_fetch_and_merge written false.
func TestTheFakeOfTheJobGivesTheRecordsRulesAsGitHubDoes(t *testing.T) {
	f := newCheckFake()
	var update []forge.Rule
	for _, r := range f.rules["layup-records"] {
		if r.Type == "update" {
			update = append(update, r)
		}
	}
	if len(update) != 1 || len(update[0].Parameters) != 0 {
		t.Errorf("the rules of layup-records in the fake: %+v; want one update rule with no parameters", f.rules["layup-records"])
	}
	if !strings.Contains(string(f.files["ruleset-records.json"]), `"update_allows_fetch_and_merge":false`) {
		t.Errorf("the ruleset file of the records in the fake: %s; want the parameter written false", f.files["ruleset-records.json"])
	}
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
	const notAttributed = "probe refused, but not attributed to the default branch's ruleset: "
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
			map[string]string{"rules:layup-probe": "creation is missing", "probe": notAttributed + "no creation rule applies to layup-probe"}, "rules:layup-probe: "},
		{"a rule that is missing on layup-records", func(f *checkFake) { f.rules["layup-records"] = f.rules["layup-records"][:1] },
			map[string]string{"rules:layup-records": "deletion is missing"}, "rules:layup-records: "},
		{"a records file that wrote true, the forge that returns false", func(f *checkFake) {
			f.files["ruleset-records.json"] = []byte(strings.Replace(string(f.files["ruleset-records.json"]), `"update_allows_fetch_and_merge":false`, `"update_allows_fetch_and_merge":true`, 1))
			f.rules["layup-records"] = ofRuleset([]forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": false}), rule("deletion", nil)}, recordsRulesetID)
		}, map[string]string{"rules:layup-records": "update: update_allows_fetch_and_merge differs"}, "rules:layup-records: "},
		{"the rule update that the forge lacks on layup-records", func(f *checkFake) { f.rules["layup-records"] = f.rules["layup-records"][1:] },
			map[string]string{"rules:layup-records": "update is missing"}, "rules:layup-records: "},
		{"a records file that wrote true, the forge that omits it", func(f *checkFake) {
			f.files["ruleset-records.json"] = []byte(strings.Replace(string(f.files["ruleset-records.json"]), `"update_allows_fetch_and_merge":false`, `"update_allows_fetch_and_merge":true`, 1))
		}, map[string]string{"rules:layup-records": "update: update_allows_fetch_and_merge is missing"}, "rules:layup-records: "},
		{"a records file that wrote false, the forge that returns true", func(f *checkFake) {
			f.rules["layup-records"] = ofRuleset([]forge.Rule{rule("update", map[string]any{"update_allows_fetch_and_merge": true}), rule("deletion", nil)}, recordsRulesetID)
		}, map[string]string{"rules:layup-records": "update: update_allows_fetch_and_merge differs"}, "rules:layup-records: "},
		{"a call of rules that fails", func(f *checkFake) { f.rulesErr["layup-probe"] = boom },
			map[string]string{"rules:layup-probe": "boom", "probe": notAttributed + "the rules of layup-probe were not read: boom"}, "rules:layup-probe: "},
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
			map[string]string{"bypass:default": "setup/rulesets.txt is not in the records", "bypass:records": "setup/rulesets.txt is not in the records", "probe": notAttributed + "setup/rulesets.txt is not in the records"}, "bypass:default: "},
		{"a setup/rulesets.txt that does not read back", func(f *checkFake) { f.applied = "not framed" },
			map[string]string{"bypass:default": "setup/rulesets.txt", "bypass:records": "setup/rulesets.txt", "probe": notAttributed + "setup/rulesets.txt"}, "bypass:default: "},
		{"an apply output of the default ruleset with no id", func(f *checkFake) { f.applied = strings.Replace(f.applied, `"id":77,`, "", 1) },
			map[string]string{"probe": notAttributed + "the output of the default branch's ruleset has no id"}, "probe: "},
		{"an apply output of the default ruleset that is two objects", func(f *checkFake) {
			f.applied = strings.Replace(f.applied, `"source":"acme/target","bypass_actors":[]}`, `"source":"acme/target","bypass_actors":[]}{"id":99}`, 1)
		}, map[string]string{"bypass:default": "the output has no list of bypass actors", "probe": notAttributed + "the output of the default branch's ruleset is not one JSON object"}, "bypass:default: "},
		{"an apply output of the default ruleset with the id of another ruleset", func(f *checkFake) { f.applied = strings.Replace(f.applied, `"id":77,`, `"id":99,`, 1) },
			map[string]string{"probe": notAttributed + "the creation of layup-probe is also restricted by ruleset 77 (Repository acme/target)"}, "probe: "},
		{"a creation rule of another ruleset on layup-probe", func(f *checkFake) {
			f.rules["layup-probe"] = append(f.rules["layup-probe"], forge.Rule{Type: "creation", RulesetSourceType: "Organization", RulesetSource: "acme", RulesetID: 9})
		}, map[string]string{"probe": notAttributed + "the creation of layup-probe is also restricted by ruleset 9 (Organization acme)"}, "probe: "},
		{"a refusal for another ref", func(f *checkFake) {
			f.probeErr = &git.FailedError{Args: []string{"push"}, Code: 1, Stderr: strings.ReplaceAll(githubRefusal, "refs/heads/layup-probe", "refs/heads/layup-probe-other"), Err: errors.New("exit status 1")}
		}, map[string]string{"probe": "probe not run: remote: error: GH013"}, "probe: "},
		{"a bypass list of the default that is not empty", func(f *checkFake) {
			f.applied = strings.Replace(f.applied, `"source":"acme/target","bypass_actors":[]`, `"source":"acme/target","bypass_actors":[{"actor_id":1}]`, 1)
		}, map[string]string{"bypass:default": "not empty"}, "bypass:default: "},
		{"a bypass list of the records that is not the App", func(f *checkFake) { f.app(7) }, map[string]string{"bypass:records": "App"}, "bypass:records: "},
		{"a probe that the forge accepts", func(f *checkFake) { f.probeErr = nil }, map[string]string{"probe": "probe accepted: " + lastCommit}, "probe: "},
		{"a probe that fails for another reason", func(f *checkFake) { f.probeErr = boom }, map[string]string{"probe": "probe not run: boom"}, "probe: "},
		{"an activity read that fails", func(f *checkFake) { f.activityErr = boom }, map[string]string{"activity": "boom"}, "activity: "},
		{"a ruleset of the forge that is in evaluate (no rule is returned)", func(f *checkFake) { f.rules = map[string][]forge.Rule{} },
			map[string]string{"rules:main": "missing", "rules:layup-probe": "missing", "rules:layup-records": "missing", "probe": notAttributed + "no creation rule applies to layup-probe"}, "rules:main: "},
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

// The rules of layup-probe are read once, with the answer of the same read for its rules row and for the attribution of the probe,
// and independently of the work area and of the ruleset file; the probe is still made on the forge's head in each case.
func TestTheRulesOfLayupProbeAreReadOnceWhateverTheWorkAreaAndTheRulesetFile(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(f *checkFake)
	}{
		{"nothing broken", func(f *checkFake) {}},
		{"a work area that cannot be made sure of", func(f *checkFake) { f.areaErr = errors.New("layup-setup is at x") }},
		{"a ruleset file that is missing", func(f *checkFake) { delete(f.files, "ruleset-default.json") }},
		{"a setup/tree.tsv that is not in the records", func(f *checkFake) { f.treeText = "" }},
		{"a call of rules that fails", func(f *checkFake) { f.rulesErr["layup-probe"] = errors.New("boom") }},
	} {
		f := newCheckFake()
		c.change(f)
		f.job().do()
		if f.rulesRead["layup-probe"] != 1 || f.probedHead != lastCommit {
			t.Errorf("%s: %d reads of the rules of layup-probe, the probe on %q; want one read, and the probe on the forge's head", c.name, f.rulesRead["layup-probe"], f.probedHead)
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
