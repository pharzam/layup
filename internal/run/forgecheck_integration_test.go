//go:build integration

package run

import (
	"context"
	"encoding/json"
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
// GH013 and the name of the ruleset, and exit 1; and, once, the push of a table of the forge check
// while the file refuse-once is there; any other push is accepted.
func refuseProbe(t *testing.T, bare, text string) {
	t.Helper()
	hook := "#!/bin/sh\nwhile read old new ref; do\n\tif [ \"$ref\" = refs/heads/layup-probe ]; then\n\t\tprintf '%s\\n' '" + text + "' >&2\n\t\texit 1\n\tfi\n\tif [ -f refuse-once ] && [ \"$(git log -1 --format=%s $new)\" = 'layup run: the forge check' ]; then\n\t\trm refuse-once\n\t\techo refused once >&2\n\t\texit 1\n\tfi\ndone\nexit 0\n"
	if err := os.WriteFile(filepath.Join(bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
}

const refusal = "error: GH013: Repository rule violations found for refs/heads/layup-probe.\n- Rulesets: layup: the default branch"

// The probe is pushed with the real internal/git: only the refusal that holds GH013 and the name of the
// ruleset of the default branch is a pass; a refusal by another ruleset, one with another text, a remote
// that cannot be reached and a push that is accepted are not; the commit is empty on top of the head, by
// the identity that is given, and on no branch of the clone.
func TestTheProbeIsAPassOnlyForTheRefusalOfTheRulesetOfTheDefaultBranch(t *testing.T) {
	who := git.Identity{Name: "layup-agent[bot]", Email: "9001+layup-agent[bot]@users.noreply.github.com", Time: time.Date(2026, 10, 13, 12, 0, 0, 0, time.UTC)}
	for _, c := range []struct {
		name, hook, result, detail string
		gone                       bool
	}{
		{"the refusal that names the ruleset", refusal, "pass", "", false},
		{"a refusal by another ruleset", strings.Replace(refusal, "layup: the default branch", "another", 1), "fail", "probe not run: remote: error: GH013", false},
		{"a refusal with another text", "error: denied\n- Rulesets: layup: the default branch", "fail", "probe not run: remote: error: denied", false},
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
		result, detail := probeVerdict(sha, err)
		if result != c.result || !strings.HasPrefix(detail, c.detail) || sha == "" {
			t.Errorf("%s: %s, %q, %v; want %s, %q, and the commit", c.name, result, detail, err, c.result, c.detail)
		}
		// No branch of the clone holds the commit; git makes the remote-tracking ref of a push that the forge accepts, to its own URL.
		accepted := c.detail == "probe accepted: "
		wantRefs := map[bool]string{true: "refs/remotes/origin/layup-probe\n", false: ""}[accepted]
		_, onForge := gitRunOK(bare, "rev-parse", "--verify", "-q", "refs/heads/layup-probe")
		bot := "layup-agent[bot]|9001+layup-agent[bot]@users.noreply.github.com"
		wantCommit := head + "|" + strings.TrimSpace(gitOut(t, clone, "rev-parse", head+"^{tree}")) + "|" + bot + "|" + bot + "|1791892800|layup: the probe of the forge check\n"
		if got := gitOut(t, clone, "for-each-ref", "--format=%(refname)", "--contains", sha); got != wantRefs || (onForge == nil) != accepted ||
			gitOut(t, clone, "show", "-s", "--format=%P|%T|%an|%ae|%cn|%ce|%at|%s", sha) != wantCommit {
			t.Errorf("%s: the commit is on %q of the clone, layup-probe on the forge: %v; want %q, and only where accepted, an empty commit on the head by the bot at the run's time", c.name, got, onForge, wantRefs)
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
	writeFile(t, filepath.Join(applies, "default.json"), "{\"source\":\"acme/target\",\"bypass_actors\":[]}\n")
	writeFile(t, filepath.Join(applies, "records.json"), "{\"source\":\"acme/target\",\"bypass_actors\":[{\"actor_id\":42,\"actor_type\":\"Integration\",\"bypass_mode\":\"always\"}]}")
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
		w.fake.rules = map[string]string{"main": rulesOf("ruleset-default.json"), "layup-probe": rulesOf("ruleset-default.json"), "layup-records": rulesOf("ruleset-records.json")}
	})
	tables := func() int {
		return len(strings.Fields(gitOut(t, w.bare, "log", "--format=%H", "--grep=^layup run: the forge check$", "refs/heads/layup-records")))
	}
	restart := func(want Step) []Step {
		t.Helper()
		rows := Restart(context.Background(), cfg)
		if r := rowOf(rows, "forge-check"); r.Name != want.Name || r.Result != want.Result || !strings.Contains(r.Detail, want.Detail) || rows[len(rows)-1].Name != "forge-check" {
			t.Fatalf("the step forge-check: %+v; want %+v, the last row\n%s", r, want, steps(rows))
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
	rows := restart(Step{"forge-check", "done", "the forge check passed"})
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
	rows = restart(Step{"forge-check", "done", "the forge check passed"})
	if got := strings.TrimSpace(gitOut(t, target, "rev-parse", "layup-setup")); rowOf(rows, "scaffold-push").Detail != "done before" || tables() != 2 || got != head || setupSessions(w) != sessions {
		t.Errorf("the restart: %+v, %d tables, layup-setup %s (was %s), %d sessions (were %d); want done before, two tables, the same head and no session", rowOf(rows, "scaffold-push"), tables(), got, head, setupSessions(w), sessions)
	}
	// Three checks fail, the rules of the records branch, the probe (accepted) and the activity, and the table is committed; a work area that is there stays.
	writeFile(t, filepath.Join(area, "stays.txt"), "x")
	fake(func() { w.fake.rules["layup-records"], w.fake.noActivity = "", true })
	must(os.Remove(filepath.Join(w.bare, "hooks", "pre-receive")))
	restart(Step{"forge-check", "fail", "rules:layup-records: update is missing"})
	passes("rules:layup-records", "probe", "activity")
	accepted := strings.Fields(table()[6][2])
	out, err := probeRef()
	by := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%an|%ae|%s", "refs/heads/layup-probe"))
	if _, stays := os.Stat(filepath.Join(area, "stays.txt")); stays != nil || tables() != 3 || len(accepted) != 3 || accepted[1] != "accepted:" || err != nil || strings.TrimSpace(out) != accepted[2] ||
		by != "layup-agent[bot]|9001+layup-agent[bot]@users.noreply.github.com|"+probeMessage {
		t.Errorf("the probe row %q, %d tables, the work area kept: %v, layup-probe %q, %v by %q; want the commit of the App's bot, accepted", table()[6], tables(), stays, out, err, by)
	}
	// A default branch one commit ahead of the last commit of setup/tree.tsv, and a ruleset file gone from the work area: the probe is made on that head.
	gitOut(t, w.bare, "update-ref", "-d", "refs/heads/layup-probe")
	refuseProbe(t, w.bare, refusal)
	fake(func() { w.fake.rules["layup-records"], w.fake.noActivity = rulesOf("ruleset-records.json"), false })
	gitOut(t, target, "switch", "-q", "-c", "ahead", "layup-setup")
	gitOut(t, target, "commit", "-q", "--allow-empty", "-m", "one commit ahead")
	ahead := strings.TrimSpace(gitOut(t, target, "rev-parse", "HEAD"))
	gitOut(t, target, "push", "-q", w.bare, "ahead:main")
	must(os.Remove(filepath.Join(area, "out", "ruleset-records.json")))
	restart(Step{"forge-check", "fail", "tree: the head of main is " + ahead + ", and setup/tree.tsv ends at " + head})
	passes("tree", "rules:layup-records")
	if got := table()[3][2]; !strings.HasPrefix(got, "open DIR/work/") || !strings.HasSuffix(got, "/out/ruleset-records.json: no such file or directory") || strings.Contains(got, w.dir) {
		t.Errorf("the rules row of a ruleset file that is gone: %q; want its path and the reason, with DIR for the host", got)
	}
	// A work area that cannot be built again (a file in the place of its directory): the three rules rows fail with that reason, DIR for the host.
	must(os.RemoveAll(filepath.Join(w.dir, "work")))
	writeFile(t, filepath.Join(w.dir, "work"), "x")
	restart(Step{"forge-check", "fail", "tree: "})
	for _, r := range table()[1:4] {
		if r[1] != "fail" || !strings.HasPrefix(r[2], "the work area: ") || strings.Contains(r[2], w.dir) || !strings.Contains(r[2], "DIR") {
			t.Errorf("the row %s: %s, %q; want a fail for the work area, with DIR for the host", r[0], r[1], r[2])
		}
	}
	// A work area built again whose head is not the last commit of setup/tree.tsv: the three rules rows fail, naming both.
	must(os.Remove(filepath.Join(w.dir, "work")))
	lines := strings.Split(strings.TrimRight(recordsFile(t, w, "setup/tree.tsv"), "\n"), "\n")
	cells := strings.Split(lines[len(lines)-1], "\t")
	cells[1] = strings.Repeat("1", 40)
	lines[len(lines)-1] = strings.Join(cells, "\t")
	pushRecords(t, w.bare, map[string]string{"setup/tree.tsv": strings.Join(lines, "\n") + "\n"}, "a tree whose last commit is not the head")
	restart(Step{"forge-check", "fail", "tree: "})
	for _, r := range table()[1:4] {
		if r[1] != "fail" || !strings.Contains(r[2], head) || !strings.Contains(r[2], strings.Repeat("1", 40)) {
			t.Errorf("the row %s: %s, %q; want a fail for the work area that names its head %s and the last commit of the tree", r[0], r[1], r[2], head)
		}
	}
	// A push of the table that the forge refuses once is tried once more, as the lease is still the run's: the table is committed.
	before := tables()
	must(os.WriteFile(filepath.Join(w.bare, "refuse-once"), nil, 0o644))
	restart(Step{"forge-check", "fail", "tree: "})
	if tables() != before+1 {
		t.Errorf("%d tables after a push that was refused once; want %d", tables(), before+1)
	}
}
