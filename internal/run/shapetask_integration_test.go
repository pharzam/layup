//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
)

// The Shape tasks on the production path (task T-jy6f, #415, row 66r): the world of a restart of a run, the opening of the two
// tasks by openShapeTask, and the driver AdvanceShapeTask of the Sessions that state.sessions() makes, with real git, the httptest
// forge through the adapter and scripted harnesses.

const gatesTSV = "kind\tstate\ttool\tcommand\tscope\tconfig\nfine\tactive\tsh\tfalse\t./*.go\t—\n"

// The prelude of every scripted harness: a comment file with its row, and the artifact rows of the files of a commit.
const harnessPrelude = `res_head='kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n'
extra() { # a session that only reads and still commits (the push rule is its test)
  if [ -e $ctl/commit ]; then echo x > extra.txt; git add -A; git commit -qm extra; fi
}
keep() { n=$(ls $ctl/prompt-$task-* 2>/dev/null | wc -l); cp ../prompt.md $ctl/prompt-$task-$n; }
comment() {
  printf '%s' "$2" > ../result/$1
  sum=$(sha256sum ../result/$1 | cut -d' ' -f1)
  printf "$res_head"'comment\t1\tresult/%s\t%s\t—\n' "$1" "$sum" > ../result/result.tsv
}
artifacts() {
  git add -A && git commit -qm work
  printf "$res_head" > ../result/result.tsv; n=0
  for p in "$@"; do [ -e "$p" ] || continue; n=$((n+1)); printf 'artifact\t%s\t%s\t%s\t—\n' $n "$p" "$(sha256sum "$p" | cut -d' ' -f1)" >> ../result/result.tsv; done
}
read -r line
step=${line#Prompt file of the step }; step=${step%% of the task *}
task=${line##* of the task }; task=${task%.}
`

type drvWorld struct {
	t                 *testing.T
	w                 *world
	r                 *state
	ctl               string // the directory where a test leaves its instructions to the scripts
	harnesses, models [][]string
}

// ctlWrite leaves a file for the scripts.
func (d *drvWorld) ctlWrite(name, text string) {
	d.t.Helper()
	if err := os.WriteFile(filepath.Join(d.ctl, name), []byte(text), 0o644); err != nil {
		d.t.Fatal(err)
	}
}

// harnessSteps is what the scripted harness does at each step; a test names a file of ctl/ to change it.
func harnessSteps(ctl string) string {
	return `ctl=` + ctl + `
case "$line" in "Prompt file"*) keep ;; esac
case $step in
"You write the options"*|"You rewrite the options"*) # the options sessions of Shape: the file options-<n> of ctl/ is the n-th output
  out=$(sed -n 's/^The output file: //p')
  n=$(cat $ctl/n-options 2>/dev/null || echo 0); n=$((n+1)); echo $n > $ctl/n-options
  [ -e $ctl/options-fail-$n ] && exit 3
  cat $ctl/options-$n > "$out" ;;
plan-of-a-task) comment plan.md 'size: small
open questions: —
new interfaces: —
new dependencies: —
tests:
'; extra ;;
plan-review) comment plan-review.md "verdict: $(cat $ctl/plan-verdict 2>/dev/null || echo approve)
"; extra ;;
specification)
  mkdir -p docs/prd docs/spec docs/tasks
  { cat $ctl/prd.md; printf '\n<!-- %s -->\n' $(date +%s%N); } > docs/prd/PRD-0001-x.md
  if [ -e $ctl/no-section ]; then rm -f docs/spec/alpha.md; else printf '## REQ-001 — Alpha\n\nThe text. %s\n' $(date +%s%N) > docs/spec/alpha.md; fi
  if [ -e $ctl/bad-prd ]; then sed -i 's/^| REQ-001 | Alpha |  |  |/| REQ-001 | Alpha | Maybe | 1 |/' docs/prd/PRD-0001-x.md; fi
  printf 'the task file %s\n' $(date +%s%N) > docs/tasks/$task.md
  if [ -e $ctl/gofile ]; then echo 'package x' > x.go; rm $ctl/gofile; fi
  if [ -e $ctl/spec-fail ]; then printf '\n| REQ-009 | x | Must | 1 | F-0003 |\n' >> docs/prd/PRD-0001-x.md; fi
  if [ -e $ctl/no-commit ]; then printf "$res_head" > ../result/result.tsv; exit 0; fi
  artifacts docs/tasks/$task.md docs/prd/PRD-0001-x.md $(ls docs/spec/*.md 2>/dev/null) x.go ;;
shape-architecture)
  mkdir -p docs/spec docs/adr docs/tasks
  printf '## Decision pending the bet\n\nA\n' > docs/spec/architecture.md
  printf 'adr\n' > docs/adr/0001-db.md
  printf 'the task file %s\n' $(date +%s%N) > docs/tasks/$task.md
  artifacts docs/tasks/$task.md docs/spec/architecture.md docs/adr/0001-db.md ;;
verification)
  n=$(cat $ctl/n-$task 2>/dev/null || echo 0); n=$((n+1)); echo $n > $ctl/n-$task
  verdict=$(sed -n ${n}p $ctl/verdicts-$task 2>/dev/null); [ -n "$verdict" ] || verdict='nothing material in scope'
  bullets=''; case $verdict in nothing*) ;; *) bullets='- docs/spec/alpha.md:1 spec: the section is thin
' ;; esac
  comment review-record.md "head: $(git rev-parse HEAD)
verdict: $verdict
findings:
$bullets"; extra ;;
esac
`
}

func newDrvWorld(t *testing.T, perms map[string]string) *drvWorld {
	return newDrvWorldWith(t, perms, gatesTSV)
}

// newDrvWorldWith is newDrvWorld whose target holds the manifest of the gates given.
func newDrvWorldWith(t *testing.T, perms map[string]string, gates string) *drvWorld {
	return newDrvWorldOpts(t, perms, drvOpts{gates: gates})
}

// drvOpts changes the world: noPRD makes a target with no PRD and no section in its first tree (the case of #432), noOptions leaves
// the records without the options of Shape.
type drvOpts struct {
	gates           string
	noPRD, noOption bool
}

func newDrvWorldOpts(t *testing.T, perms map[string]string, o drvOpts) *drvWorld {
	t.Helper()
	lint, err := os.ReadFile("../../docs/prd/prd-lint.sh")
	if err != nil {
		t.Fatal(err)
	}
	gates := o.gates
	if gates == "" {
		gates = gatesTSV
	}
	target := map[string]string{"AGENTS.md": "the agent rules\n", "docs/facts/F-0003-answers.md": "the answers\n", "docs/gates.tsv": gates,
		"docs/prd/prd-lint.sh": string(lint), "docs/guardrails.md": "# Guardrails\n\n- **Inv-1** — a rule. Check: no check yet\n"}
	if !o.noPRD {
		// A target that already has a PRD and its section: the full check of layup/spec judges a head that touches docs/spec/,
		// and the architecture head of a target with no PRD fails it with "no PRD" unless it starts at the specification's head (#432).
		target["docs/prd/PRD-0001-x.md"], target["docs/spec/alpha.md"] = specimen(t, -1), "## REQ-001 — Alpha\n\nThe text.\n"
	}
	w := newWorldWith(t, target)
	allDone(t, "Start", w.start(w.config("1.0.0")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	w.fake.mu.Lock()
	w.fake.perms = perms
	w.fake.mu.Unlock()
	d := &drvWorld{t: t, w: w, ctl: t.TempDir()}
	d.ctlWrite("prd.md", specimen(t, -1))
	bin := t.TempDir()
	os.WriteFile(filepath.Join(bin, "version"), []byte("#!/bin/sh\necho 1.0.0\n"), 0o755)
	os.WriteFile(filepath.Join(bin, "harness"), []byte("#!/bin/sh\n"+harnessPrelude+harnessSteps(d.ctl)), 0o755)
	var probes [][]string
	for i, h := range []string{"claude", "devin", "codex"} {
		probes = append(probes, []string{"S-b000000" + string(rune('1'+i)), h, "1.0.0", "m1", "passed", "—", "AGENTS.md", "m1", "2026-10-08T12:00:00Z"})
		d.harnesses = append(d.harnesses, []string{h, "", "1", filepath.Join(bin, "harness"), "stdin", filepath.Join(bin, "version"), "", "", "AGENTS.md", "", "none", "subscription", ""})
		d.models = append(d.models, []string{h, "m1", "1000000", "https://example.invalid", "2026-10-08T12:00:00Z", "yes", "—"})
	}
	harnessesTSV, err := table(records.HarnessesSchema, probes)
	if err != nil {
		t.Fatal(err)
	}
	routing := "role\ttier\tposition\tharness\tmodel\n"
	for _, role := range []string{"software-engineer", "product-owner", "systems-architect"} {
		routing += role + "\treasoning\t1\tclaude\tm1\n"
	}
	routing += "qa-engineer\treasoning\t1\tclaude\tm1\nqa-engineer\treasoning\t2\tdevin\tm1\nqa-engineer\treasoning\t3\tcodex\tm1\n"
	files := map[string]string{"harnesses.tsv": string(harnessesTSV), "routing.tsv": routing, "rule-paths.tsv": register}
	// The records of Intake that the full check of layup/spec reads, with the draft and the options of Shape.
	for p, data := range cavFiles(t) {
		if p != "tasks.tsv" && !strings.HasPrefix(p, "tasks/") {
			files[p] = string(data)
		}
	}
	for _, p := range []string{"spec/draft.tsv", "spec/options/db.tsv"} {
		if p != "spec/options/db.tsv" || !o.noOption {
			files[p] = string(shapeFiles[p])
		}
	}
	files["lease.tsv"] = "run\thost\tversion\tstarted\theartbeat\tstate\n0123456789abcdef\thost-1\t1.0.0\t2026-10-08T12:00:00Z\t1\theld\n"
	pushRecords(t, w.bare, files, "the records of Shape")
	d.restart()
	return d
}

// input is what the driver is given for a task; a test changes it before the call.
func (d *drvWorld) input(task string) ShapeTaskInput {
	base, err := git.RevParse(d.r.target, "refs/remotes/origin/main")
	if err != nil {
		d.t.Fatal(err)
	}
	in := ShapeTaskInput{Task: task, Base: base, Records: d.r.store.Base, Harnesses: d.harnesses, Models: d.models}
	if rows, err := d.r.sessions().taskRow(context.Background(), task); err == nil && rows.Kind == "architecture" {
		in.Options = dbResult
	}
	return in
}

// open opens the two Shape tasks and gives the IDs of the specification and of the architecture.
func (d *drvWorld) open() (spec, arch string) {
	d.t.Helper()
	for _, kind := range []string{"specification", "architecture"} {
		task, err := d.r.openShapeTask(context.Background(), kind)
		if err != nil {
			d.t.Fatalf("openShapeTask(%s): %v", kind, err)
		}
		if kind == "specification" {
			spec = task.ID
		} else {
			arch = task.ID
		}
	}
	return spec, arch
}

// requests counts the requests of the forge whose line starts with the prefix.
func (d *drvWorld) requests(prefix string) int {
	d.w.fake.mu.Lock()
	defer d.w.fake.mu.Unlock()
	n := 0
	for _, r := range d.w.fake.reqs {
		if strings.HasPrefix(r, prefix) {
			n++
		}
	}
	return n
}

func (d *drvWorld) register() []Task {
	d.t.Helper()
	tasks, err := ReadRegister([]byte(keptRecord(d.t, d.w, "tasks.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	return tasks
}

// eventsOf gives the events of a task from the records branch.
func (d *drvWorld) eventsOf(task string) [][]string {
	d.t.Helper()
	rows, err := records.ReadEvents([]byte(keptRecord(d.t, d.w, "tasks/"+task+"/events.tsv")))
	if err != nil {
		d.t.Fatal(err)
	}
	return rows
}

// kinds gives the kind of each event of a task, with the attempt, in order.
func (d *drvWorld) kinds(task string) []string {
	var out []string
	for _, e := range d.eventsOf(task) {
		out = append(out, e[1])
	}
	return out
}

// I1 (R10, plan review conditions 2 and 5): no issue is opened before the permissions of the loop; with them, two issues and
// two rows with the exact set, and the body of each issue is the fixed text of its kind.
func TestAShapeTaskIsOpenedOnlyWithTheLoopPermissions(t *testing.T) {
	d := newDrvWorld(t, nil)
	head, before := recordsHead(t, d.w), d.requests("POST /repos/acme/target/issues")
	if _, err := d.r.openShapeTask(context.Background(), "specification"); err == nil || !strings.Contains(err.Error(), "lacks the permissions of the task loop: ") {
		t.Fatalf("openShapeTask without the permissions = %v; want the permissions named", err)
	}
	if n := d.requests("POST /repos/acme/target/issues") - before; n != 0 || recordsHead(t, d.w) != head {
		t.Errorf("%d issues opened and the records head moved: %v; want none and no commit", n, recordsHead(t, d.w) != head)
	}
	d.w.fake.mu.Lock()
	d.w.fake.perms = loopPerms
	d.w.fake.mu.Unlock()
	d.r = restarted(t, d.w)
	spec, arch := d.open()
	tasks := d.register()
	if len(tasks) != 2 || tasks[0].ID != spec || tasks[1].ID != arch || tasks[0].Kind != "specification" || tasks[1].Kind != "architecture" {
		t.Fatalf("the register = %+v; want the specification then the architecture", tasks)
	}
	sum := digest([]byte(draftTSV))
	for i, task := range tasks {
		if strings.Join(task.Requirements, " ") != "REQ-001" || task.Issue != 3+i || task.Predecessors != nil {
			t.Errorf("the row of %s = %+v; want REQ-001, the issue %d and no predecessor", task.ID, task, 3+i)
		}
		want := [2]string{shapeIssueTitle(task.ID, task.Kind), shapeIssueText(task.ID, task.Kind, task.Requirements, sum, d.r.store.Base())}
		if got := d.w.fake.issueTexts[i+2]; got[0] != want[0] || got[1] == "" || !strings.Contains(got[1], "SHA-256 "+sum) {
			t.Errorf("the issue of %s = %q; want the title %q and a body that names the SHA-256 of the pinned draft", task.ID, got, want[0])
		}
	}
}

// Review condition 2: a draft with no requirement, a commit that the dry run refuses (a table of transitions, a parameter
// file), and a lost lease each open no issue.
func TestARefusedOpeningWritesNoIssue(t *testing.T) {
	for _, c := range []struct {
		name  string
		prep  func(d *drvWorld)
		want  string
		after func(d *drvWorld)
	}{
		{"a draft with no requirement", func(d *drvWorld) {
			pushRecords(t, d.w.bare, map[string]string{"spec/draft.tsv": "id\tkind\tstatement\tcovers\tcriterion\tsession\nR-001\tsuccess-criterion\tfast\t1\tc\tS-aaaaaaaa\n"}, "a draft")
		}, "spec/draft.tsv has no row of kind requirement or non-functional", nil},
		{"a table of transitions the dry run refuses", func(d *drvWorld) {
			pushRecords(t, d.w.bare, map[string]string{"transitions.tsv": "from\tto\tartifacts\titems\nx\ty\tplan\tschema\n"}, "a table")
		}, "transitions.tsv", nil},
		{"a file of parameters the dry run refuses", func(d *drvWorld) {
			pushRecords(t, d.w.bare, map[string]string{"parameters.tsv": "broken\n"}, "parameters")
		}, "parameters.tsv", nil},
		{"a lost lease", func(d *drvWorld) {
			pushRecords(t, d.w.bare, map[string]string{"lease.tsv": strings.Replace(keptRecord(t, d.w, "lease.tsv"), "0123456789abcdef", "fedcba9876543210", 1)}, "a lease")
		}, "lease", nil},
	} {
		d := newDrvWorld(t, loopPerms)
		c.prep(d)
		d.r = restarted(t, d.w)
		before := d.requests("POST /repos/acme/target/issues")
		_, err := d.r.openShapeTask(context.Background(), "specification")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
		if n := d.requests("POST /repos/acme/target/issues") - before; n != 0 {
			t.Errorf("%s: %d issues opened; want none", c.name, n)
		}
	}
}

// loop.md, the opening: an opened task is returned unchanged, so a second call opens no second issue and a restart that
// stopped between two tasks opens the other one only.
func TestAnOpenedShapeTaskIsReusedAndOpensNoSecondIssue(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	first, err := d.r.openShapeTask(context.Background(), "specification")
	if err != nil {
		t.Fatal(err)
	}
	d.r = restarted(t, d.w)
	spec, arch := d.open()
	if spec != first.ID || arch == first.ID || d.requests("POST /repos/acme/target/issues") != 4 {
		t.Errorf("the specification %s then %s, %d issues; want the first reused and 4 issues in all, the 2 of Start and one for each task", spec, first.ID, d.requests("POST /repos/acme/target/issues"))
	}
	if n := len(d.register()); n != 2 {
		t.Errorf("the register holds %d tasks; want 2", n)
	}
}

// Known limit (plan Q14): no call of the forge lists issues, so a stop between the issue and the row leaves an issue that the
// next run does not find; it opens a second one.
func TestAStopBetweenTheIssueAndTheRowOpensASecondIssue(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	held := keptRecord(t, d.w, "lease.tsv")
	d.r.cfg.Forge = &stealingForge{Forge: d.r.cfg.Forge, steal: func() {
		pushRecords(t, d.w.bare, map[string]string{"lease.tsv": strings.Replace(held, "0123456789abcdef", "fedcba9876543210", 1)}, "a lease")
	}}
	if _, err := d.r.openShapeTask(context.Background(), "specification"); err == nil {
		t.Fatal("openShapeTask with the lease lost after the issue = nil; want the error of the lost lease")
	}
	if gitOut(t, d.w.bare, "ls-tree", "--name-only", "layup-records", "tasks.tsv") != "" {
		t.Error("a row was registered by a run that lost its lease")
	}
	if n := d.requests("POST /repos/acme/target/issues"); n != 3 {
		t.Errorf("%d issues; want 3, the 2 of Start and the one that stands", n)
	}
	pushRecords(t, d.w.bare, map[string]string{"lease.tsv": held}, "the lease back")
	d.r = restarted(t, d.w)
	if _, err := d.r.openShapeTask(context.Background(), "specification"); err != nil {
		t.Fatal(err)
	}
	if n := d.requests("POST /repos/acme/target/issues"); n != 4 {
		t.Errorf("%d issues after the next run; want 4, a second issue for the task (the known limit)", n)
	}
}

// stealingForge acts after the issue of the forge is opened.
type stealingForge struct {
	forge.Forge
	steal func()
}

func (f *stealingForge) OpenIssue(ctx context.Context, title, body string) (int, error) {
	n, err := f.Forge.OpenIssue(ctx, title, body)
	if err == nil {
		f.steal()
	}
	return n, err
}

// loop.md, the issue text: each constant is its block of loop.md, byte for byte.
func TestTheShapeIssueBlocksOfLoopMd(t *testing.T) {
	doc, err := os.ReadFile("../../docs/spec/loop.md")
	if err != nil {
		t.Fatal(err)
	}
	if len(shapeIssueBlocks) != 2 {
		t.Errorf("%d texts in the code; want the two kinds of Shape task", len(shapeIssueBlocks))
	}
	for kind, text := range shapeIssueBlocks {
		_, block, ok := strings.Cut(string(doc), "\n```text shape-issue-"+kind+"\n")
		block, _, _ = strings.Cut(block, "```\n")
		if !ok || block != text {
			t.Errorf("loop.md has the block shape-issue-%s: %v; it is equal to the code: %v", kind, ok, block == text)
		}
	}
}
