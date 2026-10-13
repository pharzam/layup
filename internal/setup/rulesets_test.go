package setup

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// The two rulesets of S13 (task T-9v92, #274): the file of the default branch in both phases and the
// file of the records branch in phase 2 (docs/spec/scaffold.md, "The ruleset of the default branch"
// and "The ruleset of the records branch"), by their bytes and by their rules, and the step of phase 2
// that writes them with the row without-layup (NFR-002). Each case breaks one rule only, and each
// refusal is compared with its whole reason.

// rulesetPhase2Text is the file of the default branch in phase 2 for the kinds of twoKinds and the App 42:
// the gate jobs pinned to GitHub Actions, then the four layup/ checks pinned to the App.
const rulesetPhase2Text = `{
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
      "type": "creation"
    },
    {
      "type": "deletion"
    },
    {
      "type": "non_fast_forward"
    },
    {
      "type": "pull_request",
      "parameters": {
        "dismiss_stale_reviews_on_push": false,
        "require_code_owner_review": false,
        "require_last_push_approval": false,
        "required_approving_review_count": 0,
        "required_review_thread_resolution": true,
        "allowed_merge_methods": [
          "merge"
        ],
        "require_extra_approval_for_unattributed_changes": false
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
          },
          {
            "context": "layup/gates",
            "integration_id": 42
          },
          {
            "context": "layup/spec",
            "integration_id": 42
          },
          {
            "context": "layup/verify",
            "integration_id": 42
          },
          {
            "context": "layup/rules",
            "integration_id": 42
          }
        ],
        "strict_required_status_checks_policy": true
      }
    }
  ],
  "bypass_actors": []
}
`

// rulesetRecordsText is the file of the records branch for the App 42.
const rulesetRecordsText = `{
  "name": "layup: the records branch",
  "target": "branch",
  "enforcement": "active",
  "conditions": {
    "ref_name": {
      "include": [
        "refs/heads/layup-records"
      ],
      "exclude": []
    }
  },
  "rules": [
    {
      "type": "update",
      "parameters": {
        "update_allows_fetch_and_merge": false
      }
    },
    {
      "type": "deletion"
    },
    {
      "type": "non_fast_forward"
    }
  ],
  "bypass_actors": [
    {
      "actor_id": 42,
      "actor_type": "Integration",
      "bypass_mode": "always"
    }
  ]
}
`

// withoutLayup is the line of the row without-layup (scaffold.md, NFR-002).
const withoutLayup = `To merge with LAYUP absent, remove the required checks layup/gates, layup/spec, layup/verify and layup/rules from the ruleset "layup: the default branch".`

// The file of the default branch in phase 2, byte for byte.
func TestTheRulesetOfTheDefaultBranchInPhase2(t *testing.T) {
	got, err := rulesetBody([]string{"static", "layout"}, 42)
	if err != nil || string(got) != rulesetPhase2Text {
		t.Errorf("the file of the default branch (err %v):\n%s\nwant\n%s", err, got, rulesetPhase2Text)
	}
}

// The file of the records branch, byte for byte.
func TestTheRulesetOfTheRecordsBranch(t *testing.T) {
	got, err := rulesetRecordsBody(42)
	if err != nil || string(got) != rulesetRecordsText {
		t.Errorf("the file of the records branch (err %v):\n%s\nwant\n%s", err, got, rulesetRecordsText)
	}
}

// decode reads a ruleset file as the forge does: the rules keep their order, an object has no order.
func decode(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("the file is not JSON: %v\n%s", err, body)
	}
	return m
}

func rulesOf(m map[string]any) []map[string]any {
	var out []map[string]any
	for _, r := range m["rules"].([]any) {
		out = append(out, r.(map[string]any))
	}
	return out
}

// Each rule of each table of scaffold.md is in the file, in the order of the table, with its
// parameters and no other: a rule that has none has no "parameters" key; the layup/ checks are
// those of phase 2 only, each pinned to the App given, and the gate jobs to GitHub Actions.
func TestTheRulesAreThoseOfTheTables(t *testing.T) {
	check := func(context string, id int) map[string]any {
		return map[string]any{"context": context, "integration_id": float64(id)}
	}
	for _, c := range []struct {
		name   string
		app    int64
		checks []any
	}{
		{"phase 1", 0, []any{check("static", 15368), check("layout", 15368)}},
		{"phase 2", 5118100, []any{check("static", 15368), check("layout", 15368), check("layup/gates", 5118100), check("layup/spec", 5118100),
			check("layup/verify", 5118100), check("layup/rules", 5118100)}},
	} {
		body, err := rulesetBody([]string{"static", "layout"}, c.app)
		if err != nil {
			t.Fatal(err)
		}
		m := decode(t, body)
		rules := rulesOf(m)
		var types []string
		for _, r := range rules {
			types = append(types, r["type"].(string))
		}
		if want := []string{"creation", "deletion", "non_fast_forward", "pull_request", "required_status_checks"}; !reflect.DeepEqual(types, want) {
			t.Errorf("%s: the rules are %q; want %q, the order of the table", c.name, types, want)
		}
		for _, r := range rules[:3] {
			if _, ok := r["parameters"]; ok {
				t.Errorf("%s: the rule %v has parameters; want none", c.name, r["type"])
			}
		}
		pr := map[string]any{"dismiss_stale_reviews_on_push": false, "require_code_owner_review": false, "require_last_push_approval": false,
			"required_approving_review_count": float64(0), "required_review_thread_resolution": true, "allowed_merge_methods": []any{"merge"},
			"require_extra_approval_for_unattributed_changes": false}
		if got := rules[3]["parameters"]; !reflect.DeepEqual(got, any(pr)) {
			t.Errorf("%s: pull_request has %v; want %v", c.name, got, pr)
		}
		rs := map[string]any{"required_status_checks": c.checks, "strict_required_status_checks_policy": true}
		if got := rules[4]["parameters"]; !reflect.DeepEqual(got, any(rs)) {
			t.Errorf("%s: required_status_checks has %v; want %v", c.name, got, rs)
		}
		refs := map[string]any{"include": []any{"~DEFAULT_BRANCH", "refs/heads/layup-probe"}, "exclude": []any{}}
		if got := m["conditions"].(map[string]any)["ref_name"]; !reflect.DeepEqual(got, any(refs)) {
			t.Errorf("%s: the refs are %v; want %v", c.name, got, refs)
		}
		if got := m["bypass_actors"]; !reflect.DeepEqual(got, any([]any{})) {
			t.Errorf("%s: bypass_actors is %v; want the empty list, so that no one creates layup-probe", c.name, got)
		}
		if m["name"] != "layup: the default branch" || m["target"] != "branch" || m["enforcement"] != "active" {
			t.Errorf("%s: name %v, target %v, enforcement %v", c.name, m["name"], m["target"], m["enforcement"])
		}
	}
}

// The App's ID is the pin of each layup/ check and the bypass actor of the records branch, and
// nothing else: a gate job stays pinned to GitHub Actions.
func TestTheAppIDIsThePinAndTheBypassActor(t *testing.T) {
	body, err := rulesetRecordsBody(5118100)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"actor_id": float64(5118100), "actor_type": "Integration", "bypass_mode": "always"}
	if got := decode(t, body)["bypass_actors"]; !reflect.DeepEqual(got, any([]any{want})) {
		t.Errorf("the bypass actors are %v; want only %v", got, want)
	}
	rules := rulesOf(decode(t, body))
	var types []string
	for _, r := range rules {
		types = append(types, r["type"].(string))
	}
	if want := []string{"update", "deletion", "non_fast_forward"}; !reflect.DeepEqual(types, want) {
		t.Errorf("the rules of the records branch are %q; want %q", types, want)
	}
	if got := rules[0]["parameters"]; !reflect.DeepEqual(got, any(map[string]any{"update_allows_fetch_and_merge": false})) {
		t.Errorf("the rule update has %v; want update_allows_fetch_and_merge false, so that only the bypass actor updates the branch", got)
	}
	for _, r := range rules[1:] {
		if _, ok := r["parameters"]; ok {
			t.Errorf("the rule %v has parameters; want none", r["type"])
		}
	}
	m := decode(t, body)
	refs := map[string]any{"include": []any{"refs/heads/layup-records"}, "exclude": []any{}}
	if got := m["conditions"].(map[string]any)["ref_name"]; !reflect.DeepEqual(got, any(refs)) {
		t.Errorf("the refs are %v; want %v", got, refs)
	}
	if m["name"] != "layup: the records branch" || m["target"] != "branch" || m["enforcement"] != "active" {
		t.Errorf("name %v, target %v, enforcement %v", m["name"], m["target"], m["enforcement"])
	}
	if !strings.Contains(string(body), `"actor_id": 5118100`) {
		t.Errorf("the file does not hold the App's ID:\n%s", body)
	}
}

// phase2App gives the inputs of phase 2 with the App 42, which layup run reads from the forge register.
func phase2App(id int64) Phase2 {
	p := phase2()
	p.AppID = id
	return p
}

// S13 in phase 2 writes the protection file as phase 1 does, the two ruleset files, and the rows
// branch-protection.sha256, ruleset.sha256 (the default branch's file) and without-layup (NFR-002,
// source computed, ref out/ruleset-default.json); it commits as phase 1's S13 does.
func TestPhase2S13WritesBothRulesetsAndTheRow(t *testing.T) {
	f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	f.install(t)
	o := StepsPhase2(phase2App(42), Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()})
	values := [][]string{{"branch-protection.sha256", sha([]byte(protectionText)), "computed", "sha256 docs/setup/branch-protection.json"},
		{"ruleset.sha256", sha([]byte(rulesetPhase2Text)), "computed", "sha256 out/ruleset-default.json"},
		{"without-layup", withoutLayup, "computed", "out/ruleset-default.json"}}
	if o.Kind != Done || o.Evidence != "the ruleset files, and their commands in commands.sh" || !o.Commit || !reflect.DeepEqual(o.Values, values) {
		t.Errorf("S13 in phase 2: %s %q %v, values %q; want done with the rows %q", o.Kind, o.Evidence, o.Commit, o.Values, values)
	}
	for path, want := range map[string]string{"w/target/docs/setup/branch-protection.json": protectionText,
		"w/out/ruleset-default.json": rulesetPhase2Text, "w/out/ruleset-records.json": rulesetRecordsText} {
		if got := f.files[path]; got != want {
			t.Errorf("%s:\n%s\nwant\n%s", path, got, want)
		}
	}
	if len(f.files) != 3 {
		t.Errorf("S13 in phase 2 wrote %d files; want the three above", len(f.files))
	}
}

// S13 in phase 1 writes no records file and no row without-layup, and its default file has no layup/ check.
func TestS13OfPhase1WritesNoRecordsRulesetAndNoRow(t *testing.T) {
	f := &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}
	f.install(t)
	o := runS13(Input{Dir: "w", Record: lastRecord()})
	for _, v := range o.Values {
		if v[0] == "without-layup" {
			t.Errorf("S13 in phase 1 wrote the row %q; want none", v)
		}
	}
	if _, ok := f.files["w/out/ruleset-records.json"]; ok || len(f.files) != 2 || strings.Contains(f.files["w/out/ruleset-default.json"], "layup/") {
		t.Errorf("S13 in phase 1 wrote %d files, the records file %v; want the two files of phase 1, with no layup/ check", len(f.files), ok)
	}
}

// An App's ID below 1 is a fail of S13 in phase 2, with no file: a pinned check or a bypass actor
// with that ID is a ruleset that the forge refuses or reads as no one. A missing manifest fails as
// in phase 1.
func TestPhase2S13Refuses(t *testing.T) {
	for _, c := range []struct {
		name     string
		app      int64
		repo     *fakeRepo
		evidence string
	}{
		{"an App of ID 0", 0, &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}, "S13: the app_id of the forge register is 0, not a positive number"},
		{"an App of a negative ID", -7, &fakeRepo{shows: map[string]string{"docs/gates.tsv": twoKinds}}, "S13: the app_id of the forge register is -7, not a positive number"},
		{"no manifest", 42, &fakeRepo{}, "the head: docs/gates.tsv: no file docs/gates.tsv"},
	} {
		c.repo.install(t)
		o := StepsPhase2(phase2App(c.app), Calls{})["S13"].Run(Input{Dir: "w", Record: lastRecord()})
		if o.Kind != Fail || o.Evidence != c.evidence || len(c.repo.files) != 0 {
			t.Errorf("%s: %s %q, %d files; want fail %q and no file", c.name, o.Kind, o.Evidence, len(c.repo.files), c.evidence)
		}
	}
}

// S13 in phase 2 keeps what phase 1 gives it but its code and its commands: the actor and the
// hand-off (the commands of phase 2 are those of row 52c, task T-56dj).
func TestPhase2S13KeepsItsActorAndHandOff(t *testing.T) {
	one, two := Steps(Brief{}, Calls{})["S13"], StepsPhase2(phase2App(42), Calls{})["S13"]
	if reflect.ValueOf(one.Run).Pointer() == reflect.ValueOf(two.Run).Pointer() {
		t.Errorf("S13 in phase 2 has the code of phase 1")
	}
	if reflect.ValueOf(one.Commands).Pointer() == reflect.ValueOf(two.Commands).Pointer() || two.Commands == nil {
		t.Errorf("S13 in phase 2 has the commands of phase 1, or none")
	}
	if two.Actor != one.Actor || !two.HandOff || !reflect.DeepEqual(one.Reads, two.Reads) {
		t.Errorf("S13 in phase 2: actor %q, hand-off %v, reads %q; want the actor %q, the hand-off and the reads of phase 1", two.Actor, two.HandOff, two.Reads, one.Actor)
	}
}
