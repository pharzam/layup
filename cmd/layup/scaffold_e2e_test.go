//go:build e2e

package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	layuprun "github.com/pharzam/layup/internal/run"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
)

// The e2e of the Scaffold (docs/spec/scaffold.md, "The acceptance tests of M2d", the row "The steps of the
// Scaffold, end to end" and its paragraph "Decided here" of task T-eevt; row 58 of the plan, #282): the built
// binary, a fake forge on loopback, local bare repositories and the world's scripted harness programs, with no
// secret and no network, in five stages, each a restart whose exit code, table and standard error are asserted
// and whose records are read from the bare repository.

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

func TestTheStepsOfTheScaffold(t *testing.T) {
	bin := buildScaffoldBinary(t)
	w, first := startScaffoldWorld(t, bin)
	if first.code != 0 {
		t.Fatalf("layup run --new: exit %d\n%s\n%s", first.code, first.stdout, first.stderr)
	}
	env := goEnv(t)
	restart := func() result { return w.runWith(t, bin, env, "run", "acme/target", "--host", ".") }
	area := filepath.Join(w.dir, "work", "acme", "target")
	leaseIs := func(state string) {
		t.Helper()
		if lease := record(t, w, "lease.tsv"); !strings.HasSuffix(lease, "\t"+state+"\n") {
			t.Errorf("lease.tsv: %q; want its state %s", lease, state)
		}
	}

	// Stage 1: the restart to the batch, on a baseline that holds two markers, asks the setup questions of the Operator's block.
	r := restart()
	if r.code != 0 || !strings.HasSuffix(r.stdout, "\nanswers\tdone\twaiting: no answer yet in the block idea-owner and the block operator\n") {
		t.Fatalf("stage 1, the restart to the batch: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
	}
	port, owner := setup.MarkerID("docs/ops.md", angled("port")), setup.MarkerID("docs/ops.md", angled("owner"))
	questions := questionsOf(t, w)
	if _, rows := table(r.stdout); len(questions) != 7 || !strings.HasPrefix(rows["markers"][1], "2 markers") || !strings.HasPrefix(rows["batch"][1], "7 questions") {
		t.Fatalf("stage 1: the markers row %q, the batch row %q, %d questions; want two markers and seven questions: G1, the three fields, S01-stack and the two markers\n%q", rows["markers"], rows["batch"], len(questions), questions)
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
	r = restart()
	names, rows := table(r.stdout)
	waiting := "waiting for the Operator's pushes and ruleset apply"
	if r.code != 0 || names[len(names)-1] != "scaffold-push" || rows["scaffold-push"] != [2]string{"done", waiting} ||
		rows["follow-up"][1] != "nothing to ask again; Intake is done" || rows["follow-up-answers"][1] != "no follow-up; Intake is done" ||
		rows["setup"][0] != "done" || rows["verify"] != [2]string{"done", "the setup verified and its records committed"} {
		t.Fatalf("stage 2, the restart that runs the setup: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
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
	// The records of the setup: the real verification's table, the thirteen checks of its block in order, each a pass, then a row
	// for each kind of the manifest, a pass or a clear; the rule-paths register; the README of the Scaffold; the tree.
	lines := strings.Split(strings.TrimSuffix(record(t, w, "setup/verify.tsv"), "\n"), "\n")
	for i, check := range []string{"discipline-tests", "pin", "kit-history", "facts", "onboarding", "glossary", "guardrails", "markers", "adapted", "identity", "link-lint", "sources", "jobs"} {
		if lines[0] != "check\tresult\treason" || len(lines) < 15 || !strings.HasPrefix(lines[i+1], check+"\tpass\t") {
			t.Fatalf("setup/verify.tsv: row %d is %q; want %s, a pass\n%s", i+1, lines[i+1], check, strings.Join(lines, "\n"))
		}
	}
	for _, l := range lines[14:] {
		if !strings.HasPrefix(l, "gate:") || (!strings.Contains(l, "\tpass\t") && !strings.Contains(l, "\tclear\t")) {
			t.Errorf("setup/verify.tsv: the row %q of a kind; want gate:<kind>, a pass or a clear", l)
		}
	}
	if paths := record(t, w, "rule-paths.tsv"); !strings.HasPrefix(paths, "pattern\texception\tsource\n.github/\t\u2014\tbaseline\n") || !strings.Contains(paths, "\ndocs/guardrails.md\tadded lines in section 2\tbaseline\n") {
		t.Errorf("rule-paths.tsv:\n%s", paths)
	}
	readme := record(t, w, "README.md")
	for _, file := range []string{"setup/record.tsv", "setup/verify.tsv", "setup/tree.tsv", "setup/rulesets.txt", "setup/forge-check.tsv", "rule-paths.tsv"} {
		if !strings.Contains(readme, "`"+file+"`") {
			t.Errorf("the README of the records does not name %s:\n%s", file, readme)
		}
	}
	record(t, w, "setup/record.tsv")
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
	r = restart()
	names, rows = table(r.stdout)
	if r.code != 0 || names[len(names)-1] != "forge-check" || rows["scaffold-push"] != [2]string{"done", "the apply of the rulesets recorded"} || rows["forge-check"] != [2]string{"done", "the forge check passed"} {
		t.Fatalf("stage 4, the restart that records the applies: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
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
	if tree := record(t, w, "setup/tree.tsv"); !strings.Contains(tree, "\t"+strings.TrimSpace(head)+"\t") {
		t.Errorf("setup/tree.tsv does not end at the head of the default branch %s:\n%s", strings.TrimSpace(head), tree)
	}
	leaseIs("released")

	// Stage 5: a later start runs the forge check after the lease, writes the table again, and ends the rows at a row that is not a pass.
	time.Sleep(1100 * time.Millisecond) // the clock of the rows has a second: the next table's is later than the first's
	r = restart()
	names, rows = table(r.stdout)
	if r.code != 0 || strings.Join(names[:6], " ") != "forge clone version lease forge-check probe" || rows["forge-check"] != [2]string{"done", "the forge check passed"} || names[len(names)-1] != "scaffold-push" {
		t.Fatalf("stage 5, a later start: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
	}
	own := map[string]string{"follow-up": "nothing to ask again; Intake is done", "follow-up-answers": "no follow-up; Intake is done"}
	for _, name := range names[6:] {
		want, ok := own[name]
		if !ok {
			want = "done before"
		}
		if rows[name] != [2]string{"done", want} {
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
	r = restart()
	names, rows = table(r.stdout)
	if r.code != 1 || names[len(names)-1] != "forge-check" || rows["forge-check"] != [2]string{"fail", "rules:layup-records: update is missing"} || rows["probe"] != ([2]string{}) {
		t.Fatalf("stage 5, a later start whose forge gives no rule: exit %d\n%s\n%s", r.code, r.stdout, r.stderr)
	}
	if got := checks(); got[3][0] != "rules:layup-records" || got[3][1] != "fail" {
		t.Errorf("setup/forge-check.tsv after the failing later start: %q; want the row rules:layup-records fail", got)
	}
	leaseIs("held")
	if strings.Contains(r.stdout+r.stderr, "ghs_") {
		t.Errorf("an output holds a token of the forge")
	}
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
		insert("\tsp=$(printf '%s\\n' \"$p\" | sed -n 's/^The span table: //p')\n", proseSession)
		for path, text := range prose {
			files["bin/files/"+path] = text
		}
		// The Operator's gh, which the test puts first on the PATH of commands.sh: it answers the apply of a ruleset with its
		// input and the source of the repository, as the forge does.
		files["ops/gh"] = "#!/bin/sh\n[ \"$1\" = api ] || exit 2\nwhile [ $# -gt 0 ]; do [ \"$1\" = --input ] && in=$2; shift; done\n" +
			"sed '1s/^{/{\"source\":\"acme\\/target\",/' \"$in\"\n"
	})
	// The forge refuses the push of the probe, as the ruleset of the default branch does: GH013 and the name of the ruleset.
	hook := "#!/bin/sh\nwhile read old new ref; do\n\tif [ \"$ref\" = refs/heads/layup-probe ]; then\n" +
		"\t\tprintf '%s\\n' 'error: GH013: Repository rule violations found for refs/heads/layup-probe.' '- Rulesets: layup: the default branch' >&2\n\t\texit 1\n\tfi\ndone\nexit 0\n"
	if err := os.WriteFile(filepath.Join(w.bare, "hooks", "pre-receive"), []byte(hook), 0o755); err != nil {
		t.Fatal(err)
	}
	return w, w.start(t, binary)
}
