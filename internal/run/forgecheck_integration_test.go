//go:build integration

package run

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/verify"
)

// The step forge-check on a real bare repository and on the real step runner (docs/spec/scaffold.md,
// REQ-003, Scaffold 7 and its paragraph "Decided here" of task T-wphp; the acceptance row "The probe";
// row 57a of the plan, #280).

// refuseProbe makes the bare repository refuse a push to layup-probe, as the forge does: the text of
// GH013 in the form that GitHub prints (it names the ref and the violated rule, and no ruleset; task
// T-gkvi, #323), and exit 1; and, once, the push of a table of the forge check
// while the file refuse-once is there; any other push is accepted.
func refuseProbe(t *testing.T, bare, text string) {
	t.Helper()
	hook := "#!/bin/sh\nwhile read old new ref; do\n\tif [ \"$ref\" = refs/heads/layup-probe ]; then\n\t\tprintf '%s\\n' '" + text + "' >&2\n\t\texit 1\n\tfi\n\tif [ -f refuse-once ] && [ \"$(git log -1 --format=%s $new)\" = 'layup run: the forge check' ]; then\n\t\trm refuse-once\n\t\techo refused once >&2\n\t\texit 1\n\tfi\ndone\nexit 0\n"
	if err := os.WriteFile(filepath.Join(bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
}

const refusal = "error: GH013: Repository rule violations found for refs/heads/layup-probe.\nReview all repository rules at https://github.com/acme/target/rules?ref=refs/heads/layup-probe\n\n- Cannot create ref due to creations being restricted.\n"

// inRuleset gives a JSON list of rules as the forge answers for a ruleset: each rule with the source type, the source and the ID of
// the ruleset, and the rule update with no parameters when the file wrote update_allows_fetch_and_merge false, as GitHub gives it (first
// hand, run 6 of the demo of M2d, 2026-10-15; task T-j7c5, #330; runs/T-j7c5/test-runs.md); withAnotherCreation adds a creation rule of an
// organization's ruleset, and withUpdateTrue gives the update parameter as true.
func inRuleset(t *testing.T, rules string, id int) string {
	t.Helper()
	var list []map[string]any
	if err := json.Unmarshal([]byte(rules), &list); err != nil {
		t.Fatal(err)
	}
	for _, r := range list {
		r["ruleset_source_type"], r["ruleset_source"], r["ruleset_id"] = "Repository", "acme/target", id
		if params, ok := r["parameters"].(map[string]any); ok && r["type"] == "update" && params["update_allows_fetch_and_merge"] == false {
			delete(r, "parameters")
		}
	}
	out, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func withUpdateTrue(t *testing.T, rules string) string {
	t.Helper()
	var list []map[string]any
	if err := json.Unmarshal([]byte(rules), &list); err != nil {
		t.Fatal(err)
	}
	for _, r := range list {
		if r["type"] == "update" {
			r["parameters"] = map[string]any{"update_allows_fetch_and_merge": true}
		}
	}
	out, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

func withAnotherCreation(t *testing.T, rules string) string {
	t.Helper()
	var list []map[string]any
	if err := json.Unmarshal([]byte(rules), &list); err != nil {
		t.Fatal(err)
	}
	list = append(list, map[string]any{"type": "creation", "ruleset_source_type": "Organization", "ruleset_source": "acme", "ruleset_id": 9001})
	out, err := json.Marshal(list)
	if err != nil {
		t.Fatal(err)
	}
	return string(out)
}

// The probe is pushed with the real internal/git: GitHub's refusal of the ref layup-probe is a pass, as far as the attribution goes (it
// is given here as none); a refusal of another ref, one with another text, a remote that cannot be reached and a push that is accepted
// are not; the commit is empty on top of the head, by the identity that is given, and on no branch of the clone.
func TestTheProbeIsAPassOnlyForGitHubsRefusalOfTheProbeRef(t *testing.T) {
	who := git.Identity{Name: "layup-agent[bot]", Email: "9001+layup-agent[bot]@users.noreply.github.com", Time: time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC)}
	for _, c := range []struct {
		name, hook, result, detail string
		gone                       bool
	}{
		{"GitHub's refusal of the probe ref", refusal, "pass", "", false},
		{"a refusal of another ref", strings.ReplaceAll(refusal, "layup-probe", "layup-probe-other"), "fail", "probe not run: remote: error: GH013", false},
		{"a refusal with another text", "error: denied\n- Cannot create ref due to creations being restricted.", "fail", "probe not run: remote: error: denied", false},
		{"a remote that cannot be reached", "", "fail", "probe not run: ", true},
		{"a push that is accepted", "", "fail", "probe accepted: ", false},
	} {
		src, bare, clone := t.TempDir(), filepath.Join(t.TempDir(), "target.git"), filepath.Join(t.TempDir(), "clone")
		gitOut(t, src, "init", "-q", "-b", "main")
		gitOut(t, src, "commit", "-q", "--allow-empty", "-m", "root")
		gitOut(t, "/", "clone", "-q", "--bare", src, bare)
		gitOut(t, bare, "config", "receive.autogc", "false")
		gitOut(t, "/", "clone", "-q", "file://"+bare, clone)
		head := strings.TrimSpace(gitOut(t, clone, "rev-parse", "origin/main"))
		if c.hook != "" {
			refuseProbe(t, bare, c.hook)
		}
		url := "file://" + bare
		if c.gone {
			url = "file://" + filepath.Join(bare, "nowhere.git")
		}
		sha, err := pushProbe(clone, url, git.Auth{}, head, who)
		result, detail := probeVerdict(sha, err, nil)
		if result != c.result || !strings.HasPrefix(detail, c.detail) || sha == "" {
			t.Errorf("%s: %s, %q, %v; want %s, %q, and the commit", c.name, result, detail, err, c.result, c.detail)
		}
		// No branch of the clone holds the commit; a push that the forge accepts, to the clone's own URL, makes the remote-tracking ref in
		// the git of this host and not in that of CI.
		accepted := c.detail == "probe accepted: "
		refs := strings.Fields(gitOut(t, clone, "for-each-ref", "--format=%(refname)", "--contains", sha))
		_, onForge := gitRunOK(bare, "rev-parse", "--verify", "-q", "refs/heads/layup-probe")
		bot := "layup-agent[bot]|9001+layup-agent[bot]@users.noreply.github.com"
		wantCommit := head + "|" + strings.TrimSpace(gitOut(t, clone, "rev-parse", head+"^{tree}")) + "|" + bot + "|" + bot + "|1791892800|layup: the probe of the forge check\n"
		if len(refs) > 1 || len(refs) == 1 && (!accepted || refs[0] != "refs/remotes/origin/layup-probe") || (onForge == nil) != accepted ||
			gitOut(t, clone, "show", "-s", "--format=%P|%T|%an|%ae|%cn|%ce|%at|%s", sha) != wantCommit {
			t.Errorf("%s: the commit is on %q of the clone, layup-probe on the forge: %v; want no ref, or only where accepted the remote-tracking ref, an empty commit on the head by the bot at the run's time", c.name, refs, onForge)
		}
	}
	if sha, err := pushProbe(t.TempDir(), "file:///nowhere", git.Auth{}, "1111", who); err == nil || sha != "" {
		t.Errorf("a head that is no object ID: %q, %v; want an error and no commit", sha, err)
	}
}

// The step on the real step runner: after the apply of the Operator it checks the forge and commits the table, whatever
// the result, at each restart; a work area that is missing is built again, one that is there is used as it is.
func TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart(t *testing.T) {
	saved := verifyRun
	t.Cleanup(func() { verifyRun = saved })
	w := scaffoldWorld(t, proseFiles())
	verifyRun = func(string, func(int, int, string) func(), io.Writer) (verify.Table, error) {
		return verify.Table{Rows: []verify.Row{{Check: "pin", Result: "pass", Reason: "—"}}}, nil
	}
	must := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	cfg := w.config("0.1.0-dev")
	toSetup(t, w, "go")
	area := filepath.Join(w.dir, "work", "acme", "target")
	target := filepath.Join(area, "target")
	head := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup"))
	// The Operator pushes the setup commits and applies the two rulesets; the forge answers their rules and refuses the probe.
	gitOut(t, target, "push", "-q", w.bare, "layup-setup:main")
	applies := filepath.Join(w.dir, "applies", "acme", "target")
	writeFile(t, filepath.Join(applies, "default.json"), "{\"id\":7001,\"source\":\"acme/target\",\"bypass_actors\":[]}\n")
	writeFile(t, filepath.Join(applies, "records.json"), "{\"id\":7002,\"source\":\"acme/target\",\"bypass_actors\":[{\"actor_id\":42,\"actor_type\":\"Integration\",\"bypass_mode\":\"always\"}]}")
	refuseProbe(t, w.bare, refusal)
	rulesOf := func(file string) string {
		var f struct {
			Rules json.RawMessage `json:"rules"`
		}
		data, err := os.ReadFile(filepath.Join(area, "out", file))
		must(err)
		must(json.Unmarshal(data, &f))
		return string(f.Rules)
	}
	fake := func(change func()) { w.fake.mu.Lock(); change(); w.fake.mu.Unlock() }
	fake(func() {
		w.fake.rules = map[string]string{"main": inRuleset(t, rulesOf("ruleset-default.json"), 7001), "layup-probe": inRuleset(t, rulesOf("ruleset-default.json"), 7001), "layup-records": inRuleset(t, rulesOf("ruleset-records.json"), 7002)}
	})
	// The fake answers as GitHub does: the rule update of layup-records has no parameters, while the file that the real writer made holds the
	// parameter written false (#330).
	var served, wrote []map[string]any
	must(json.Unmarshal([]byte(w.fake.rules["layup-records"]), &served))
	must(json.Unmarshal([]byte(rulesOf("ruleset-records.json")), &wrote))
	fileParams, _ := wrote[0]["parameters"].(map[string]any)
	if _, has := served[0]["parameters"]; served[0]["type"] != "update" || has || wrote[0]["type"] != "update" || fileParams["update_allows_fetch_and_merge"] != false {
		t.Errorf("the fake gives %v for the rule update of layup-records, the file %v; want no parameters in the answer and update_allows_fetch_and_merge false in the file", served[0], wrote[0])
	}
	tables := func() int {
		return len(strings.Fields(gitOut(t, w.bare, "log", "--format=%H", "--grep=^layup run: the forge check$", "refs/heads/layup-records")))
	}
	// restart runs a restart and checks which rows it makes, in which order, and the notes of progress: a row's note is its place of the 22 rows
	// that a restart with Start done plans, after the 1/6 forge and 2/6 clone of one with Start not done.
	first := []string{"forge", "clone", "version", "lease", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision", "setup", "verify", "scaffold-push", "forge-check"}
	later := []string{"forge", "clone", "version", "lease", "forge-check", "probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision", "setup", "verify", "scaffold-push"}
	var progress []string
	cfg.Step = func(i, n int, name string) { progress = append(progress, fmt.Sprintf("%d/%d %s", i, n, name)) }
	restart := func(want Step, names []string) []Step {
		t.Helper()
		progress = nil
		rows := Restart(context.Background(), cfg)
		if r := rowOf(rows, "forge-check"); r.Name != want.Name || r.Result != want.Result || !strings.Contains(r.Detail, want.Detail) {
			t.Fatalf("the step forge-check: %+v; want %+v\n%s", r, want, steps(rows))
		}
		var got, notes []string
		for _, row := range rows {
			got = append(got, row.Name)
		}
		for i, name := range names {
			total := 22
			if i < 2 {
				total = 6
			}
			notes = append(notes, fmt.Sprintf("%d/%d %s", i+1, total, name))
		}
		if !slices.Equal(got, names) || !slices.Equal(progress, notes) {
			t.Errorf("the rows %q and the notes of progress %q; want the rows %q and the notes %q\n%s", got, progress, names, notes, steps(rows))
		}
		return rows
	}
	table := func() [][]string {
		t.Helper()
		rows, err := ReadForgeCheck([]byte(recordsFile(t, w, "setup/forge-check.tsv")))
		must(err)
		return rows
	}
	passes := func(fails ...string) { // each row of the table but those of fails is a pass, with the run and the clock of the run
		t.Helper()
		var last time.Time
		for _, r := range table() {
			at, err := time.Parse(time.RFC3339, r[4])
			if err != nil || at.Before(last) || at.After(w.clock.Now()) || r[3] != "0123456789abcdef" {
				t.Errorf("the row %s: run %q at %q (after %v); want the run's ID and a time of the run's clock that does not go back", r[0], r[3], r[4], last)
			}
			last = at
			if want := map[bool]string{true: "fail", false: "pass"}[slices.Contains(fails, r[0])]; r[1] != want {
				t.Errorf("the row %s is %s, %q; want %s", r[0], r[1], r[2], want)
			}
		}
	}
	probeRef := func() (string, error) {
		return gitRunOK(w.bare, "rev-parse", "--verify", "-q", "refs/heads/layup-probe")
	}
	// The first restart after the apply: the record of the applies, then the table, in one commit of its path.
	rows := restart(Step{"forge-check", "done", "the forge check passed"}, first)
	passes()
	commit := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%H", "--grep=^layup run: the forge check$", "refs/heads/layup-records"))
	if changed := strings.Fields(gitOut(t, w.bare, "diff", "--name-only", commit+"~1", commit)); strings.Join(changed, " ") != "setup/forge-check.tsv" || len(table()) != 8 || table()[1][0] != "rules:main" {
		t.Errorf("the table of %d rows, rules:main second, in a commit that changes %q; want eight rows and setup/forge-check.tsv alone", len(table()), changed)
	}
	if _, err := probeRef(); rowOf(rows, "scaffold-push").Detail != "the apply of the rulesets recorded" || err == nil {
		t.Errorf("the step before: %+v, layup-probe on the forge: %v; want the apply recorded and no layup-probe after a refused probe", rowOf(rows, "scaffold-push"), err)
	}
	// A later restart, the work area gone: built again to the same commits, the table written again.
	sessions := setupSessions(w)
	must(os.RemoveAll(filepath.Join(w.dir, "work")))
	rows = restart(Step{"forge-check", "done", "the forge check passed"}, later)
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup")); rowOf(rows, "scaffold-push").Detail != "done before" || tables() != 2 || got != head || setupSessions(w) != sessions {
		t.Errorf("the restart: %+v, %d tables, layup-setup %s (was %s), %d sessions (were %d); want done before, two tables, the same head and no session", rowOf(rows, "scaffold-push"), tables(), got, head, setupSessions(w), sessions)
	}
	// At a later start each step of Intake and of the Scaffold is done before, but the two that give the detail of a finished Intake.
	doneBefore := map[string]string{"rule-gaps": "done before", "meaning": "done before", "numbering": "done before", "spans": "done before", "draft": "done before", "completeness": "done before",
		"markers": "done before", "batch": "done before", "answers": "done before", "answers-review": "done before", "follow-up": "nothing to ask again; Intake is done",
		"follow-up-answers": "no follow-up; Intake is done", "revision": "done before", "setup": "done before", "verify": "done before", "scaffold-push": "done before"}
	for name, detail := range doneBefore {
		if got := rowOf(rows, name); got.Result != "done" || got.Detail != detail {
			t.Errorf("at a later start the step %s is %+v; want done, %q", name, got, detail)
		}
	}
	// Another ruleset's creation rule on layup-probe: the push is refused as before, and the probe is a fail that names the ruleset; the table is committed.
	fake(func() { w.fake.rules["layup-probe"] = withAnotherCreation(t, w.fake.rules["layup-probe"]) })
	restart(Step{"forge-check", "fail", "probe: probe refused, but not attributed to the default branch's ruleset: the creation of layup-probe is also restricted by ruleset 9001 (Organization acme)"}, later[:5])
	passes("probe")
	// The Scaffold is not done while a row is not a pass, so the restart that follows is of the first-start chain; with the rule gone, it passes.
	fake(func() { w.fake.rules["layup-probe"] = inRuleset(t, rulesOf("ruleset-default.json"), 7001) })
	restart(Step{"forge-check", "done", "the forge check passed"}, first)
	passes()
	// The forge returns true for the update parameter that the file wrote false: the row of the records rules fails, naming it; with the answer back, it passes.
	fake(func() { w.fake.rules["layup-records"] = withUpdateTrue(t, w.fake.rules["layup-records"]) })
	restart(Step{"forge-check", "fail", "rules:layup-records: update: update_allows_fetch_and_merge differs"}, later[:5])
	passes("rules:layup-records")
	fake(func() { w.fake.rules["layup-records"] = inRuleset(t, rulesOf("ruleset-records.json"), 7002) })
	restart(Step{"forge-check", "done", "the forge check passed"}, first)
	passes()
	// Three checks fail, the rules of the records branch, the probe (accepted) and the activity, and the table is committed; a work area that is there stays.
	writeFile(t, filepath.Join(area, "stays.txt"), "x")
	fake(func() { w.fake.rules["layup-records"], w.fake.noActivity = "", true })
	must(os.Remove(filepath.Join(w.bare, "hooks", "pre-receive")))
	restart(Step{"forge-check", "fail", "rules:layup-records: update is missing"}, later[:5])
	passes("rules:layup-records", "probe", "activity")
	accepted := strings.Fields(table()[6][2])
	out, err := probeRef()
	by := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%an|%ae|%s", "refs/heads/layup-probe"))
	if _, stays := os.Stat(filepath.Join(area, "stays.txt")); stays != nil || tables() != 7 || len(accepted) != 3 || accepted[1] != "accepted:" || err != nil || strings.TrimSpace(out) != accepted[2] ||
		by != "layup-agent[bot]|9001+layup-agent[bot]@users.noreply.github.com|"+probeMessage {
		t.Errorf("the probe row %q, %d tables, the work area kept: %v, layup-probe %q, %v by %q; want the commit of the App's bot, accepted", table()[6], tables(), stays, out, err, by)
	}
	// A default branch one commit ahead of the last commit of setup/tree.tsv, and a ruleset file gone from the work area: the probe is made on that head.
	gitOut(t, w.bare, "update-ref", "-d", "refs/heads/layup-probe")
	refuseProbe(t, w.bare, refusal)
	fake(func() {
		w.fake.rules["layup-records"], w.fake.noActivity = inRuleset(t, rulesOf("ruleset-records.json"), 7002), false
	})
	gitOut(t, target, "switch", "-q", "-c", "ahead", "layup-setup")
	gitOut(t, target, "commit", "-q", "--allow-empty", "-m", "one commit ahead")
	ahead := strings.TrimSpace(gitOut(t, target, "rev-parse", "HEAD"))
	gitOut(t, target, "push", "-q", w.bare, "ahead:main")
	must(os.Remove(filepath.Join(area, "out", "ruleset-records.json")))
	restart(Step{"forge-check", "fail", "tree: the head of main is " + ahead + ", and setup/tree.tsv ends at " + head}, first)
	passes("tree", "rules:layup-records")
	if got := table()[3][2]; !strings.HasPrefix(got, "open DIR/work/") || !strings.HasSuffix(got, "/out/ruleset-records.json: no such file or directory") || strings.Contains(got, w.dir) {
		t.Errorf("the rules row of a ruleset file that is gone: %q; want its path and the reason, with DIR for the host", got)
	}
	// A work area that cannot be built again (a link where its directory belongs): the three rules rows fail with that reason, DIR for the host.
	must(os.RemoveAll(filepath.Join(w.dir, "work")))
	must(os.Symlink(t.TempDir(), filepath.Join(w.dir, "work")))
	restart(Step{"forge-check", "fail", "tree: "}, first)
	for _, r := range table()[1:4] {
		if r[1] != "fail" || !strings.HasPrefix(r[2], "the work area: DIR/work is a link") || strings.Contains(r[2], w.dir) {
			t.Errorf("the row %s: %s, %q; want a fail for the work area that is a link, with DIR for the host", r[0], r[1], r[2])
		}
	}
	// A work area built again whose loop fails (a glossary whose citation resolves to no fact): the three rules rows fail with its reason.
	must(os.Remove(filepath.Join(w.dir, "work")))
	glossary := "setup/inputs/files/docs/glossary.md"
	pushRecords(t, w.bare, map[string]string{glossary: "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository (`F-0099#1`) |\n"}, "a glossary whose citation resolves to no fact")
	restart(Step{"forge-check", "fail", "tree: "}, first)
	for _, r := range table()[1:4] {
		if r[1] != "fail" || !strings.HasPrefix(r[2], "the work area: ") || !strings.Contains(r[2], "glossary") {
			t.Errorf("the row %s: %s, %q; want a fail for the work area, naming the glossary", r[0], r[1], r[2])
		}
	}
	pushRecords(t, w.bare, map[string]string{glossary: proseFiles()["docs/glossary.md"]}, "the glossary again")
	// A work area built again whose head is not the last commit of setup/tree.tsv: the three rules rows fail, naming both (the half-built
	// area of the failed loop is removed first, as a work area that is there is used as it is).
	must(os.RemoveAll(filepath.Join(w.dir, "work")))
	lines := strings.Split(strings.TrimRight(recordsFile(t, w, "setup/tree.tsv"), "\n"), "\n")
	cells := strings.Split(lines[len(lines)-1], "\t")
	cells[1] = strings.Repeat("1", 40)
	lines[len(lines)-1] = strings.Join(cells, "\t")
	pushRecords(t, w.bare, map[string]string{"setup/tree.tsv": strings.Join(lines, "\n") + "\n"}, "a tree whose last commit is not the head")
	restart(Step{"forge-check", "fail", "tree: "}, first)
	for _, r := range table()[1:4] {
		if r[1] != "fail" || !strings.Contains(r[2], head) || !strings.Contains(r[2], strings.Repeat("1", 40)) {
			t.Errorf("the row %s: %s, %q; want a fail for the work area that names its head %s and the last commit of the tree", r[0], r[1], r[2], head)
		}
	}
	// A push of the table that the forge refuses once is tried once more, as the lease is still the run's: the table is committed.
	before := tables()
	must(os.WriteFile(filepath.Join(w.bare, "refuse-once"), nil, 0o644))
	restart(Step{"forge-check", "fail", "tree: "}, first)
	if tables() != before+1 {
		t.Errorf("%d tables after a push that was refused once; want %d", tables(), before+1)
	}
	// A table that its reader refuses, a header and no row (a reading of "each row is pass" would call it done): the first-start chain, and the table written again.
	before = tables()
	pushRecords(t, w.bare, map[string]string{"setup/forge-check.tsv": "check\tresult\tdetail\trun\tat\n"}, "a table that its reader refuses")
	restart(Step{"forge-check", "fail", "tree: "}, first)
	if _, err := ReadForgeCheck([]byte(recordsFile(t, w, "setup/forge-check.tsv"))); err != nil || tables() != before+1 {
		t.Errorf("the table after a refused one: %v, %d tables; want it readable again and %d tables", err, tables(), before+1)
	}
}
