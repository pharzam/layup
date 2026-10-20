//go:build e2e

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	layuprun "github.com/pharzam/layup/internal/run"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/work"
)

// The e2e of the Scaffold (docs/spec/scaffold.md, "The acceptance tests of M2d", the row "The steps of the
// Scaffold, end to end" and its paragraph "Decided here" of task T-eevt; row 58 of the plan, #282): the built
// binary, a fake forge on loopback, local bare repositories and the world's scripted harness programs, with no
// secret and no network, in five stages: each restart's exit code, rows (the whole sequence of its steps), details and
// standard error are asserted, the printed file is run by the test, and the records are read from the bare repository.

// table gives the rows of a table of run-steps: the steps in order, and for each its result and detail.
func table(stdout string) ([]string, map[string][2]string) {
	var names []string
	rows := map[string][2]string{}
	for i, l := range strings.Split(strings.TrimSuffix(stdout, "\n"), "\n") {
		if f := strings.SplitN(l, "\t", 3); i > 0 && len(f) == 3 {
			names, rows[f[0]] = append(names, f[0]), [2]string{f[1], f[2]}
		}
	}
	return names, rows
}

// record gives a file of layup-records.
func record(t *testing.T, w runWorld, path string) string {
	t.Helper()
	out, err := gitRun(w.bare, "show", "refs/heads/layup-records:"+path)
	if err != nil {
		t.Fatalf("%s of layup-records: %v\n%s", path, err, out)
	}
	return out
}

// leaks fails the scenario when a table of a run holds a path of the host or a token of the forge in a detail, or an
// output holds a token; the printed commands hold absolute paths, by their rule, and are not a detail of a table.
func leaks(t *testing.T, stage string, w runWorld, r result) {
	t.Helper()
	names, rows := table(r.stdout)
	for _, name := range names {
		if d := rows[name][1]; strings.Contains(d, w.dir) || strings.Contains(d, "ghs_") {
			t.Errorf("%s: the detail of %s holds a path of the host or a token: %q", stage, name, d)
		}
	}
	if strings.Contains(r.stdout+r.stderr, "ghs_") {
		t.Errorf("%s: an output holds a token of the forge:\n%s\n%s", stage, r.stdout, r.stderr)
	}
}

func TestTheStepsOfTheScaffold(t *testing.T) {
	bin := buildScaffoldBinary(t)
	w, first := startScaffoldWorld(t, bin)
	if first.code != 0 {
		t.Fatalf("layup run --new: exit %d\n%s\n%s", first.code, first.stdout, first.stderr)
	}
	leaks(t, "Start", w, first)
	env := goEnv(t)
	restart := func(stage string) result {
		t.Helper()
		r := w.runWith(t, bin, env, "run", "acme/target", "--host", ".")
		leaks(t, stage, w, r)
		return r
	}
	// The rows that each restart makes, in order: the four steps before the steps of Intake, those of Intake, the three of the
	// Scaffold, and forge-check: last at the first start, and after lease once the Scaffold is done.
	front := []string{"forge", "clone", "version", "lease"}
	intake := []string{"probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision"}
	scaffold := []string{"setup", "verify", "scaffold-push"}
	stage1, stage2 := slices.Concat(front, intake[:10]), slices.Concat(front, intake, scaffold)
	atFirst, atLater := slices.Concat(stage2, []string{"forge-check"}), slices.Concat(front, []string{"forge-check"}, intake, scaffold, []string{"shape"})
	area := filepath.Join(w.dir, "work", "acme", "target")
	leaseIs := func(state string) {
		t.Helper()
		if lease := record(t, w, "lease.tsv"); !strings.HasSuffix(lease, "\t"+state+"\n") {
			t.Errorf("lease.tsv: %q; want its state %s", lease, state)
		}
	}

	// Stage 1: the restart to the batch, on a baseline that holds two markers, asks the setup questions of the Operator's block.
	r := restart("stage 1")
	names, rows := table(r.stdout)
	if r.code != 0 || !slices.Equal(names, stage1) || !strings.HasSuffix(r.stdout, "\nanswers\tdone\twaiting: no answer yet in the block idea-owner and the block operator\n") {
		t.Fatalf("stage 1, the restart to the batch: exit %d, the rows %q\n%s\n%s", r.code, names, r.stdout, r.stderr)
	}
	port, owner := setup.MarkerID("docs/ops.md", angled("port")), setup.MarkerID("docs/ops.md", angled("owner"))
	questions := questionsOf(t, w)
	if len(questions) != 8 || !strings.HasPrefix(rows["markers"][1], "2 markers") || !strings.HasPrefix(rows["batch"][1], "8 questions") {
		t.Fatalf("stage 1: the markers row %q, the batch row %q, %d questions; want two markers and eight questions: G1, the four fields, S01-stack and the two markers\n%q", rows["markers"], rows["batch"], len(questions), questions)
	}

	// Stage 2: the answers, then the restart that ends Intake, runs setup and verify, and waits at scaffold-push.
	var idea, operator strings.Builder
	for _, q := range questions {
		text := "the answer to " + q[0]
		switch q[2] {
		case "band.B":
			text = "5.0"
		case "band.U":
			text = "10.0"
		case "appetite.hours":
			text = "40.0"
		case "success-criteria":
			text = "the setup completes"
		case "S01-stack":
			text = "go"
		case port:
			text = "8080"
		case owner:
			text = "gap"
		}
		if q[3] == "operator" {
			operator.WriteString(q[0] + ": " + text + "\n")
		} else {
			idea.WriteString(q[0] + ": " + text + "\n")
		}
	}
	answer(t, w, 901, questions[0], idea.String())
	answerBy(t, w, 903, 101, "pharzam", questions[0], operator.String())
	r = restart("stage 2")
	names, rows = table(r.stdout)
	waiting := "waiting for the Operator's pushes and ruleset apply"
	if r.code != 0 || !slices.Equal(names, stage2) || rows["scaffold-push"] != [2]string{"done", waiting} ||
		rows["follow-up"][1] != "nothing to ask again; Intake is done" || rows["follow-up-answers"][1] != "no follow-up; Intake is done" ||
		rows["setup"][0] != "done" || rows["verify"] != [2]string{"done", "the setup verified and its records committed"} {
		t.Fatalf("stage 2, the restart that runs the setup: exit %d, the rows %q\n%s\n%s", r.code, names, r.stdout, r.stderr)
	}
	commands, err := os.ReadFile(filepath.Join(area, "out", "commands.sh"))
	if err != nil || !strings.Contains(string(commands), w.dir) {
		t.Fatalf("out/commands.sh of the work area: %v; want its paths absolute\n%s", err, commands)
	}
	for _, line := range append(strings.Split(strings.TrimRight(string(commands), "\n"), "\n"), "scaffold-push: run "+filepath.Join(area, "out", "commands.sh")+" with your own login, from a clone with no hooks installed") {
		if !strings.Contains(r.stderr, line+"\n") {
			t.Errorf("standard error has no line %q:\n%s", line, r.stderr)
		}
	}
	// The records of the setup: the real verification's table, the thirteen checks of its block in order, each a pass, then a row for
	// each kind of the Go entry's manifest, in its order (the active kinds pass, the pending kinds are clear); the three files that the
	// step publishes are those of out/ of the work area; the README is the block of the specification; the tree.
	lines := strings.Split(strings.TrimSuffix(record(t, w, "setup/verify.tsv"), "\n"), "\n")
	var got []string
	for _, l := range lines {
		f := strings.Split(l, "\t")
		got = append(got, f[0]+" "+f[1])
	}
	want := []string{"check result", "discipline-tests pass", "pin pass", "kit-history pass", "facts pass", "onboarding pass", "glossary pass", "guardrails pass", "markers pass", "adapted pass",
		"identity pass", "link-lint pass", "sources pass", "jobs pass", "gate:static pass", "gate:layout clear", "gate:boundary clear", "gate:contract clear", "gate:test pass"}
	if !slices.Equal(got, want) {
		t.Errorf("setup/verify.tsv: %q; want %q", got, want)
	}
	for path, from := range map[string]string{"setup/record.tsv": work.RecordPath, "setup/verify.tsv": work.VerifyPath, "rule-paths.tsv": "out/rule-paths.tsv"} {
		if out, err := os.ReadFile(filepath.Join(area, filepath.FromSlash(from))); err != nil || string(out) != record(t, w, path) {
			t.Errorf("%s of layup-records is not %s of the work area: %v\n%s", path, from, err, out)
		}
	}
	if paths := record(t, w, "rule-paths.tsv"); !strings.HasPrefix(paths, "pattern\texception\tsource\n.github/\t\u2014\tbaseline\n") || !strings.Contains(paths, "\ndocs/guardrails.md\tadded lines in section 2\tbaseline\n") {
		t.Errorf("rule-paths.tsv:\n%s", paths)
	}
	spec, err := os.ReadFile(filepath.Join("..", "..", "docs", "spec", "scaffold.md"))
	if err != nil {
		t.Fatal(err)
	}
	_, block, found := strings.Cut(string(spec), "\n```text scaffold-readme\n")
	block, _, _ = strings.Cut(block, "```\n")
	if readme := record(t, w, "README.md"); !found || readme != block {
		t.Errorf("README.md of layup-records is not the block scaffold-readme of scaffold.md:\n%s\nwant\n%s", readme, block)
	}
	if out, err := gitRun(w.bare, "show", "refs/heads/layup-records:setup/rulesets.txt"); err == nil {
		t.Errorf("setup/rulesets.txt is in the records while the step waits:\n%s", out)
	}
	leaseIs("released")

	// Stage 3: the test is the Operator: it runs the printed file, as it is, with a stand-in gh and no configuration of git.
	sh := exec.Command("sh", filepath.Join(area, "out", "commands.sh"))
	sh.Env = []string{"PATH=" + filepath.Join(w.dir, "ops") + ":" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "LC_ALL=C"}
	if out, err := sh.CombinedOutput(); err != nil {
		t.Fatalf("out/commands.sh: %v\n%s", err, out)
	}
	head, _ := gitRun(filepath.Join(area, "target"), "rev-parse", "layup-setup")
	if main, _ := gitRun(w.bare, "rev-parse", "refs/heads/main"); main != head {
		t.Errorf("the default branch of the forge is %q; want the head of layup-setup %q", strings.TrimSpace(main), strings.TrimSpace(head))
	}
	applies := map[string]string{}
	for _, f := range []string{"default.json", "records.json"} {
		out, err := os.ReadFile(filepath.Join(w.dir, "applies", "acme", "target", f))
		if err != nil || !strings.Contains(string(out), `"source":"acme/target"`) {
			t.Fatalf("the apply output %s: %v\n%s", f, err, out)
		}
		applies[f] = string(out)
	}

	// Stage 4: the restart that records the applies and runs the forge check; the bare repository refuses the probe.
	r = restart("stage 4")
	names, rows = table(r.stdout)
	if r.code != 0 || !slices.Equal(names, atFirst) || rows["scaffold-push"] != [2]string{"done", "the apply of the rulesets recorded"} || rows["forge-check"] != [2]string{"done", "the forge check passed"} {
		t.Fatalf("stage 4, the restart that records the applies: exit %d, the rows %q\n%s\n%s", r.code, names, r.stdout, r.stderr)
	}
	if got, want := record(t, w, "setup/rulesets.txt"), "== default.json\n"+applies["default.json"]+"\n== records.json\n"+applies["records.json"]+"\n"; got != want {
		t.Errorf("setup/rulesets.txt:\n%q\nwant\n%q", got, want)
	}
	checks := func() [][]string {
		t.Helper()
		got, err := layuprun.ReadForgeCheck([]byte(record(t, w, "setup/forge-check.tsv")))
		if err != nil {
			t.Fatal(err)
		}
		return got
	}
	table1 := checks()
	for i, c := range table1 {
		if c[1] != "pass" || (i == 1 && c[0] != "rules:main") {
			t.Errorf("setup/forge-check.tsv row %d: %q; want eight checks, each pass, rules:main second", i, c)
		}
	}
	if out, err := gitRun(w.bare, "rev-parse", "--verify", "-q", "refs/heads/layup-probe"); err == nil {
		t.Errorf("layup-probe is on the forge after a refused probe: %s", out)
	}
	// The fake forge answers as GitHub does (#330): the rule update of layup-records has no parameters, while the file that the writer made holds
	// update_allows_fetch_and_merge written false.
	rec := httptest.NewRecorder()
	w.forge.serve(rec, httptest.NewRequest(http.MethodGet, "/repos/acme/target/rules/branches/layup-records", nil))
	var served []map[string]any
	var file struct {
		Rules []map[string]any `json:"rules"`
	}
	data, err := os.ReadFile(filepath.Join(area, "out", "ruleset-records.json"))
	if err != nil || json.Unmarshal(data, &file) != nil || json.Unmarshal(rec.Body.Bytes(), &served) != nil || len(served) == 0 || len(file.Rules) == 0 {
		t.Fatalf("the rules of layup-records from the fake: %s; the ruleset file: %v", rec.Body.String(), err)
	}
	fileParams, _ := file.Rules[0]["parameters"].(map[string]any)
	if _, has := served[0]["parameters"]; served[0]["type"] != "update" || has || file.Rules[0]["type"] != "update" || fileParams["update_allows_fetch_and_merge"] != false {
		t.Errorf("the fake gives %v for the rule update of layup-records, the file %v; want no parameters in the answer and update_allows_fetch_and_merge false in the file", served[0], file.Rules[0])
	}
	if tree := record(t, w, "setup/tree.tsv"); !strings.Contains(tree, "\t"+strings.TrimSpace(head)+"\t") {
		t.Errorf("setup/tree.tsv does not end at the head of the default branch %s:\n%s", strings.TrimSpace(head), tree)
	}
	leaseIs("released")

	// Stage 5: a later start runs the forge check after the lease, writes the table again, and ends the rows at a row that is not a pass.
	time.Sleep(1100 * time.Millisecond) // the clock of the rows has a second: the next table's is later than the first's
	r = restart("stage 5")
	names, rows = table(r.stdout)
	if r.code != 0 || !slices.Equal(names, atLater) || rows["forge-check"] != [2]string{"done", "the forge check passed"} {
		t.Fatalf("stage 5, a later start: exit %d, the rows %q\n%s\n%s", r.code, names, r.stdout, r.stderr)
	}
	own := map[string]string{"follow-up": "nothing to ask again; Intake is done", "follow-up-answers": "no follow-up; Intake is done", "shape": "waiting for the Operator: shape: no survivor"}
	for _, name := range names[5:] {
		want, ok := own[name]
		if !ok && name != "probe" {
			want = "done before"
		}
		if rows[name][0] != "done" || (name != "probe" && rows[name][1] != want) {
			t.Errorf("at a later start the step %s is %q; want done, %q", name, rows[name], want)
		}
	}
	commits, _ := gitRun(w.bare, "log", "--format=%H", "--grep=^layup run: the forge check$", "refs/heads/layup-records")
	if table2 := checks(); table2[0][4] <= table1[0][4] || len(strings.Fields(commits)) != 2 {
		t.Errorf("the table of the later start: %q after %q, in %d commits; want a new records commit, later", table2[0], table1[0], len(strings.Fields(commits)))
	}
	leaseIs("released")
	w.forge.mu.Lock()
	w.forge.noRules = map[string]bool{"layup-records": true}
	w.forge.mu.Unlock()
	r = restart("stage 5, a forge that gives no rule")
	names, rows = table(r.stdout)
	if r.code != 1 || !slices.Equal(names, atLater[:5]) || rows["forge-check"] != [2]string{"fail", "rules:layup-records: update is missing"} {
		t.Fatalf("stage 5, a later start whose forge gives no rule: exit %d, the rows %q\n%s\n%s", r.code, names, r.stdout, r.stderr)
	}
	if got := checks(); got[3][0] != "rules:layup-records" || got[3][1] != "fail" {
		t.Errorf("setup/forge-check.tsv after the failing later start: %q; want the row rules:layup-records fail", got)
	}
	leaseIs("held")
}

// angled gives the marker of a name, with the angle quotes U+2039 and U+203A.
func angled(name string) string { return string(rune(0x2039)) + name + string(rune(0x203A)) }

var (
	scaffoldOnce sync.Once
	scaffoldDir  string // the directory of the second binary and its baseline; TestMain removes it
	scaffoldBin  string
	scaffoldErr  error
)

// buildScaffoldBinary makes the scaffold baseline once, as the integration level's scaffoldBaseline does (the
// stand-in baseline of internal/standin, the link-lint.sh of this repository, a note of the kit in the backlog, a
// guide whose link the deletion of the history breaks, two markers, one of them twice in a file, and two files that
// check adapted flags, in one commit), and builds the binary with its pin, as buildRunBinary does for a README.
func buildScaffoldBinary(t *testing.T) string {
	t.Helper()
	scaffoldOnce.Do(func() {
		dir, err := os.MkdirTemp("", "layup-scaffold-e2e-")
		if err != nil {
			scaffoldErr = err
			return
		}
		scaffoldDir = dir
		base := filepath.Join(dir, "baseline")
		url, _, err := standin.Baseline(base)
		if err != nil {
			scaffoldErr = err
			return
		}
		lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "links", "link-lint.sh"))
		if err != nil {
			scaffoldErr = err
			return
		}
		for p, text := range map[string]string{
			"docs/links/link-lint.sh": string(lint),
			"docs/tasks/backlog.md":   "# Backlog\n\n- **T-0001**: a task of the baseline ([#1](" + url + "/issues/1))\n\n> a note of the kit\n> with a link ([#2](" + url + "/issues/2))\n\nKeep this line.\n",
			"docs/guide.md":           "# Guide\n\nSee [the decision](decisions/D-0001-stand-in.md).\n",
			"docs/ops.md":             "# Ops\n\nThe port is " + angled("port") + ".\nAgain " + angled("port") + ".\nThe owner is " + angled("owner") + ".\n",
			"docs/how-to.md":          "# How to\n\nAdapt this to your project.\n",
			"docs/glossary.md":        "# Glossary\n\nThe words of the kit.\n",
		} {
			if err := os.MkdirAll(filepath.Dir(filepath.Join(base, p)), 0o755); err == nil {
				err = os.WriteFile(filepath.Join(base, p), []byte(text), 0o644)
			}
			if err != nil {
				scaffoldErr = err
				return
			}
		}
		// One commit: the default branch of the target must be one root commit, which the Operator pushes.
		os.RemoveAll(filepath.Join(base, ".git"))
		for _, args := range [][]string{{"init", "-q", "-b", "main"}, {"add", "-A"}, {"commit", "-q", "-m", "the stand-in baseline with the scaffold"}} {
			if out, err := gitRun(base, args...); err != nil {
				scaffoldErr = errors.New(out)
				return
			}
		}
		commit, _ := gitRun(base, "rev-parse", "HEAD")
		scaffoldBin = filepath.Join(dir, "layup")
		ld := "-X github.com/pharzam/layup/internal/cli.pinSource=" + url + " -X github.com/pharzam/layup/internal/cli.pinCommit=" + strings.TrimSpace(commit)
		build := exec.Command("go", "build", "-ldflags", ld, "-o", scaffoldBin, ".")
		build.Env = append(os.Environ(), "GOFLAGS=", "GOENV=off", "GOWORK=off", "GOTOOLCHAIN=local", "GOPROXY=off")
		if out, err := build.CombinedOutput(); err != nil {
			scaffoldErr = errors.New(string(out))
		}
	})
	if scaffoldErr != nil {
		t.Fatal(scaffoldErr)
	}
	return scaffoldBin
}

// prose are the files that the world's harness writes for the F- rows of the stop table of the prose step of the setup
// (the guide that S05 asks for, and the six of the prose step), by path.
var prose = map[string]string{
	"docs/guide.md":                    "# Guide\n\nThe decisions of the target are in its records.\n",
	"docs/onboarding-for-engineers.md": "# Onboarding\n\nThe work starts from [the problem statement](facts/problem-statement-brief.md).\n",
	"docs/guardrails.md":               "# Guardrails\n\n## 1. Rules\n\n- **Inv-1** \u2014 the setup keeps its record. Check: no check yet\n",
	"docs/glossary.md":                 "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository |\n",
	"docs/how-to.md":                   "# How to\n\nRun the tests of the product.\n",
	"README.md":                        "# acme/target\n\nThe product repository, set up from its pinned baseline ([the pin](docs/setup/armature.pin)). Its records are on the branch `layup-records`.\n",
	"AGENTS.md":                        "# AGENTS.md\n\nAgent context for **acme/target**.\n",
}

// proseSession is the part of the world's harness that answers the prose session of the setup: the target is on its task
// branch at the head of layup-setup, and the files of bin/files/ that the F- rows of the stop table name are copied below
// the output directory.
const proseSession = `	if printf '%s\n' "$p" | grep -q '^The output directory: '; then
		dir=$(printf '%s\n' "$p" | sed -n 's/^The output directory: //p')
		stop=$(printf '%s\n' "$p" | sed -n 's/^The stop table: //p')
		[ -f docs/setup/armature.pin ] || exit 5
		[ "$(git rev-parse --abbrev-ref HEAD | cut -c1-5)" = task/ ] || exit 6
		for want in $(awk -F '\t' 'NR > 1 { print $4 }' "$stop"); do
			src="$(dirname "$0")/files/$want"
			[ -f "$src" ] || exit 7
			mkdir -p "$dir/$(dirname "$want")" && cp "$src" "$dir/$want" || exit 8
		done
		exit 0
	fi
`

// startScaffoldWorld makes the world of the scenario: the world of the other scenarios of layup run, whose harness
// answers a session of the step markers with no proposal (the baseline holds two markers) and the prose session of the
// setup, and starts it.
func startScaffoldWorld(t *testing.T, binary string) (runWorld, result) {
	t.Helper()
	w := newRunWorldWith(t, func(files map[string]string) {
		insert := func(anchor, text string) {
			if strings.Count(files["bin/harness"], anchor) != 1 {
				t.Fatalf("the harness of the world has no one place after %q", anchor)
			}
			files["bin/harness"] = strings.Replace(files["bin/harness"], anchor, anchor+text, 1)
		}
		insert("\tcase \"$out\" in\n", "\t*/markers.tsv) printf 'id\\tvalue\\tquote\\n' > \"$out\" ;;\n")
		// The options of Shape (task T-3f94, row 66h): two options of one decision that both fail the invariant Inv-1, in the first round
		// and in the rewrite, so the step shape ends in its wait, with no survivor.
		insert("\tcase \"$out\" in\n", "\t*/options.tsv) printf 'decision\\toption\\tconstraint\\tresult\\tevidence\\ndb\\tpg\\tInv-1\\tfail\\t\\342\\200\\224\\ndb\\tsqlite\\tInv-1\\tfail\\t\\342\\200\\224\\n' > \"$out\" ;;\n")
		insert("\tsp=$(printf '%s\\n' \"$p\" | sed -n 's/^The span table: //p')\n", proseSession)
		for path, text := range prose {
			files["bin/files/"+path] = text
		}
		// The Operator's gh, which the test puts first on the PATH of commands.sh: it takes the one call that the file makes, the
		// POST of a ruleset file of the work area to /repos/acme/target/rulesets, and answers it with the input and the source of
		// the repository, as the forge does; any other call is an error.
		files["ops/gh"] = "#!/bin/sh\n[ \"$#\" = 6 ] && [ \"$1 $2 $3 $4 $5\" = 'api --method POST /repos/acme/target/rulesets --input' ] && " +
			"case \"$6\" in */out/ruleset-default.json | */out/ruleset-records.json) [ -f \"$6\" ] ;; *) false ;; esac || { echo \"gh: not the apply of a ruleset of acme/target: $*\" >&2; exit 2; }\n" +
			"case \"$6\" in */out/ruleset-default.json) id=7001 ;; *) id=7002 ;; esac\n" +
			"sed \"1s/^{/{\\\"source\\\":\\\"acme\\\\/target\\\",\\\"id\\\":$id,/\" \"$6\"\n"
	})
	// The forge refuses the push of the probe, as the ruleset of the default branch does, in the form that GitHub prints: GH013 and the ref,
	// the line of the rules URL and the violated rule, and no name of a ruleset (task T-gkvi).
	hook := "#!/bin/sh\nwhile read old new ref; do\n\tif [ \"$ref\" = refs/heads/layup-probe ]; then\n" +
		"\t\tprintf '%s\\n' 'error: GH013: Repository rule violations found for refs/heads/layup-probe.' 'Review all repository rules at https://github.com/acme/target/rules?ref=refs/heads/layup-probe' '' '- Cannot create ref due to creations being restricted.' >&2\n\t\texit 1\n\tfi\ndone\nexit 0\n"
	if err := os.WriteFile(filepath.Join(w.bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	return w, w.start(t, binary)
}
