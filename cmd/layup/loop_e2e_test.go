//go:build e2e

package main

// The e2e of the specification task (docs/spec/loop.md, "The acceptance tests of M2e", the row "The specification task, end to end", and its
// paragraph "Decided here" of task T-3xvc, row 67 of the plan of M2e, #368): the built binary, a fake forge on loopback, local bare repositories and
// scripted harness programs, with no secret and no network, from layup run --new to the merge of the two tasks of Shape.
//
// The routes of M2e below are the e2e copy of the integration level's fakeGitHub and pullRoutes (internal/run, files of package run, which this
// package cannot import): the two are kept in step by hand (D1). The fake stands for the target's CI (D3): each required job of GitHub Actions has a
// completed success check run at every head, and the four layup/ statuses come from the binary alone.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	layuprun "github.com/pharzam/layup/internal/run"
	"github.com/pharzam/layup/internal/setup"
)

// A pull is a pull request of the fake forge; its head is the tip of its branch at each request, as the forge's is.
type pull struct {
	number                   int
	node, branch, base, body string
	head, mergeSHA           string
	draft, merged            bool
}

// m2eRoutes are the calls of M2e that the fake forge plays: the comments of the issues of the tasks, the commit statuses, the check runs, and the pull requests.
type m2eRoutes struct {
	comments map[int][]map[string]any // by issue, the App's comments
	statuses map[string][]map[string]any
	atMerge  map[string][]map[string]any // the statuses of a SHA as they stood at the merge request for it, newest first
	pulls    []*pull
}

func (p *pull) wire() map[string]any {
	w := map[string]any{"number": p.number, "node_id": p.node, "state": "open", "draft": p.draft, "merged": p.merged, "merged_at": nil,
		"merge_commit_sha": strings.Repeat("9", 40), "head": map[string]any{"sha": p.head}, "body": p.body}
	if p.merged {
		w["state"], w["merged_at"], w["merge_commit_sha"] = "closed", "2026-10-18T12:00:00Z", p.mergeSHA
	}
	return w
}

// serve plays one call under the lock of the forge; it gives false for a call it does not know. A call that writes adds a line to the log of the forge.
func (m *m2eRoutes) serve(f *fakeForge, w http.ResponseWriter, r *http.Request) bool {
	const issues, pulls = "/repos/acme/target/issues/", "/repos/acme/target/pulls"
	var in struct {
		Title, Head, Base, Body, Context, State, Description, SHA string
		MergeMethod                                               string `json:"merge_method"`
		Draft                                                     bool
		Variables                                                 struct{ ID string }
	}
	reply := func(code int, v any) {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(v)
	}
	app := map[string]any{"id": 9001, "login": "layup-agent[bot]"}
	for _, p := range m.pulls {
		if tip, err := gitRun(f.bare, "rev-parse", "refs/heads/"+p.branch); err == nil && !p.merged {
			p.head = strings.TrimSpace(tip)
		}
	}
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&in)
	}
	switch path, n := r.URL.Path, 0; {
	case r.Method == "POST" && strings.HasPrefix(path, issues) && strings.HasSuffix(path, "/comments"):
		n, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(path, issues), "/comments"))
		if m.comments == nil {
			m.comments = map[int][]map[string]any{}
		}
		id := 1001
		for _, list := range m.comments {
			id += len(list)
		}
		m.comments[n] = append(m.comments[n], map[string]any{"id": id, "user": app, "performed_via_github_app": map[string]any{"slug": "layup-agent"},
			"created_at": "2026-10-08T13:00:00Z", "updated_at": "2026-10-08T13:00:00Z", "body": in.Body})
		reply(201, map[string]any{"id": id})
	case r.Method == "GET" && strings.HasPrefix(path, issues) && strings.HasSuffix(path, "/comments"):
		n, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(path, issues), "/comments"))
		reply(200, append([]map[string]any{}, m.comments[n]...))
	case r.Method == "POST" && strings.HasPrefix(path, "/repos/acme/target/statuses/"):
		sha := strings.TrimPrefix(path, "/repos/acme/target/statuses/")
		if m.statuses == nil {
			m.statuses = map[string][]map[string]any{}
		}
		m.statuses[sha] = append([]map[string]any{{"context": in.Context, "state": in.State, "description": in.Description, "creator": app}}, m.statuses[sha]...)
		reply(201, map[string]any{"id": len(m.statuses[sha])})
	case r.Method == "GET" && strings.HasPrefix(path, "/repos/acme/target/commits/") && strings.HasSuffix(path, "/statuses"):
		reply(200, append([]map[string]any{}, m.statuses[strings.TrimSuffix(strings.TrimPrefix(path, "/repos/acme/target/commits/"), "/statuses")]...))
	case r.Method == "GET" && strings.HasPrefix(path, "/repos/acme/target/commits/") && strings.HasSuffix(path, "/check-runs"):
		// Each required job of GitHub Actions that the served ruleset names (the integration_id of the Actions app, 15368), completed and success.
		var set struct {
			Rules []struct {
				Type       string
				Parameters struct {
					Checks []struct {
						Context       string `json:"context"`
						IntegrationID int    `json:"integration_id"`
					} `json:"required_status_checks"`
				}
			}
		}
		data, _ := os.ReadFile(filepath.Join(f.host, "work", "acme", "target", "out", "ruleset-default.json"))
		json.Unmarshal(data, &set)
		runs := []map[string]any{}
		for _, rule := range set.Rules {
			for _, c := range rule.Parameters.Checks {
				if rule.Type == "required_status_checks" && c.IntegrationID == 15368 {
					runs = append(runs, map[string]any{"name": c.Context, "status": "completed", "conclusion": jobConclusion(c.Context), "app": map[string]any{"id": 15368, "slug": "github-actions"}})
				}
			}
		}
		reply(200, map[string]any{"total_count": len(runs), "check_runs": runs})
	case r.Method == "POST" && path == pulls:
		head, err := gitRun(f.bare, "rev-parse", "refs/heads/"+in.Head)
		if err != nil {
			reply(422, map[string]string{"message": "head does not exist"})
			return true
		}
		p := &pull{number: len(m.pulls) + 1, branch: in.Head, base: in.Base, body: in.Body, head: strings.TrimSpace(head), draft: in.Draft}
		p.node = "PR_" + strconv.Itoa(p.number)
		m.pulls = append(m.pulls, p)
		reply(201, p.wire())
	case r.Method == "GET" && path == pulls:
		list := []map[string]any{}
		for i := len(m.pulls) - 1; i >= 0; i-- {
			if m.pulls[i].branch == strings.TrimPrefix(r.URL.Query().Get("head"), "acme:") {
				list = append(list, m.pulls[i].wire())
			}
		}
		reply(200, list)
	case strings.HasPrefix(path, pulls+"/") && !strings.HasSuffix(path, "/merge") && (r.Method == "GET" || r.Method == "PATCH"):
		n, _ = strconv.Atoi(strings.TrimPrefix(path, pulls+"/"))
		if n < 1 || n > len(m.pulls) {
			reply(404, map[string]string{"message": "Not Found"})
			return true
		}
		if r.Method == "PATCH" {
			m.pulls[n-1].body = in.Body
		}
		reply(200, m.pulls[n-1].wire())
	case r.Method == "POST" && path == "/graphql":
		for _, p := range m.pulls {
			if p.node == in.Variables.ID {
				p.draft = false
				f.reqs = append(f.reqs, "mark "+p.node)
				reply(200, map[string]any{"data": map[string]any{"markPullRequestReadyForReview": map[string]any{"pullRequest": map[string]any{
					"id": p.node, "number": p.number, "isDraft": false, "state": "OPEN", "headRefOid": p.head}}}})
				return true
			}
		}
		reply(200, map[string]any{"errors": []map[string]string{{"message": "Could not resolve to a node"}}})
	case r.Method == "PUT" && strings.HasPrefix(path, pulls+"/") && strings.HasSuffix(path, "/merge"):
		n, _ = strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(path, pulls+"/"), "/merge"))
		f.reqs = append(f.reqs, fmt.Sprintf("merge %d %s", n, in.SHA))
		if m.atMerge == nil {
			m.atMerge = map[string][]map[string]any{}
		}
		m.atMerge[in.SHA] = slices.Clone(m.statuses[in.SHA])
		switch {
		case n < 1 || n > len(m.pulls):
			reply(404, map[string]string{"message": "Not Found"})
		case in.MergeMethod != "merge" || m.pulls[n-1].draft || m.pulls[n-1].merged:
			reply(405, map[string]string{"message": "Pull Request is not mergeable"})
		case in.SHA != m.pulls[n-1].head:
			reply(409, map[string]string{"message": "Head branch was modified. Review and try the merge again."})
		default:
			p := m.pulls[n-1]
			tree, e1 := gitRun(f.bare, "rev-parse", p.head+"^{tree}")
			main, e2 := gitRun(f.bare, "rev-parse", "refs/heads/main")
			merge, e3 := gitRun(f.bare, "commit-tree", strings.TrimSpace(tree), "-p", strings.TrimSpace(main), "-p", p.head, "-m", fmt.Sprintf("Merge pull request #%d", n))
			if _, e4 := gitRun(f.bare, "update-ref", "refs/heads/main", strings.TrimSpace(merge)); e1 != nil || e2 != nil || e3 != nil || e4 != nil {
				reply(500, map[string]string{"message": "the merge commit could not be made"})
				return true
			}
			p.merged, p.mergeSHA = true, strings.TrimSpace(merge)
			reply(200, map[string]any{"sha": p.mergeSHA, "merged": true, "message": "Pull Request successfully merged"})
		}
	default:
		return false
	}
	return true
}

// The sessions of the tasks and of the brief, which the world's harness answers before its other cases (the task steps of the integration level's harnessSteps).
// The PRD is made from the forms of the instruction in the prompt file, filled from the draft that the prompt lists (condition 3 of the plan review): the
// draft has the requirement R-001 and the success criterion R-002 that the setup answers give; the harness reads both from the draft listed in the prompt.
const loopSessions = `env | grep -q 'ghs_\|GITHUB_TOKEN' && exit 9
first=$(printf '%s\n' "$p" | head -n 1)
case "$first" in
"You write three parts of the brief"*)
	out=$(printf '%s\n' "$p" | sed -n 's/^Write these three files of Markdown in the directory \(.*\), and no other file there:$/\1/p')
	mkdir -p "$out" && printf 'The problem in prose.\n' > "$out/problem.md" && printf 'appetite: USD 5, 12 hours\nWhy: small.\n' > "$out/appetite.md" &&
		printf 'No rabbit hole.\n' > "$out/rabbit-holes.md"
	exit $? ;;
"Prompt file of the step "*)
	step=${first#Prompt file of the step }; step=${step%% of the task *}
	task=${first##* of the task }; task=${task%.}
	res='kind\tn\tvalue\tsha256\treason\nstatus\t1\tcompleted\t—\tdone\n'
	comment() { printf '%s' "$2" > ../result/$1; printf "$res"'comment\t1\tresult/%s\t%s\t—\n' "$1" "$(sha256sum ../result/$1 | cut -d' ' -f1)" > ../result/result.tsv; }
	artifacts() {
		git add -A && git commit -qm work && printf "$res" > ../result/result.tsv; n=0
		for f in "$@"; do n=$((n+1)); printf 'artifact\t%s\t%s\t%s\t—\n' $n "$f" "$(sha256sum "$f" | cut -d' ' -f1)" >> ../result/result.tsv; done
	}
	mkdir -p docs/prd docs/spec docs/adr docs/tasks
	case $step in
	plan-of-a-task) comment plan.md 'size: small
open questions: —
new interfaces: —
new dependencies: —
tests:
' ;;
	plan-review) comment plan-review.md 'verdict: approve
' ;;
	specification)
		draft=$(sed -n 's/^- \(.*draft\.tsv\) [0-9a-f]*$/\1/p' ../prompt.md)
		col() { awk -F '\t' -v k="$1" -v c="$2" '$2 == k { print $c; exit }' "$draft"; }
		{ printf '# PRD-0001. The target\n\n'
		sed -n '/^## 6\. Functional requirements$/,/^Put the row of a requirement/p' ../prompt.md | sed '$d' |
			sed -e "/<the success criterion>/{s|<the success criterion>|$(col success-criterion 5)|;s|F-[0-9]*#<n>|$(col success-criterion 4)|}" \
			-e "s|<NNN>|001|g;s|<the requirement>|$(col requirement 3)|;s|<the acceptance criterion>|$(col requirement 5)|;s|<n>|$(col requirement 4)|g" \
			-e 's|<the phase tag>|1|;s|<the name of the phase>|The core|;s|<what ships in it>|Everything.|'
		printf '**Non-goals.** none.\n'; } > docs/prd/PRD-0001-x.md
		printf '## REQ-001 — Alpha\n\nThe text.\n' > docs/spec/alpha.md
		printf 'the task file\n' > docs/tasks/$task.md
		artifacts docs/tasks/$task.md docs/prd/PRD-0001-x.md docs/spec/alpha.md ;;
	shape-architecture)
		printf '## The database: pg\n\nThe decision db is decided: pg.\n' > docs/spec/architecture.md
		printf 'adr\n' > docs/adr/0001-db.md
		printf 'the task file\n' > docs/tasks/$task.md
		artifacts docs/tasks/$task.md docs/spec/architecture.md docs/adr/0001-db.md ;;
	verification) comment review-record.md "head: $(git rev-parse HEAD)
verdict: nothing material in scope
findings:
" ;;
	close-out)
		printf 'the task file, closed out\n' > docs/tasks/$task.md; printf 'the log of %s\n' $task >> docs/tasks/completed.md
		artifacts docs/tasks/$task.md docs/tasks/completed.md ;;
	esac
	exit $? ;;
esac
`

// loopWorldFiles changes the files of the scaffold world: a third harness, the routing of the task steps (the authors on fake, the verifier on spare or
// third, and the close-out on the resolved tier of software-developer, both tiers), and the harness' cases of loopSessions.
func loopWorldFiles(t *testing.T) func(files map[string]string) {
	return func(files map[string]string) {
		for _, p := range []string{"registers/harnesses.tsv", "registers/models.tsv"} {
			lines := strings.Split(strings.TrimSuffix(files[p], "\n"), "\n")
			files[p] = strings.Join(append(lines, strings.Replace(lines[len(lines)-1], "spare", "third", 1)), "\n") + "\n"
		}
		files["registers/routing.tsv"] = "role\ttier\tposition\tharness\tmodel\nsoftware-engineer\treasoning\t1\tfake\tm1\nsoftware-developer\texecution\t1\tfake\tm1\n" +
			"software-developer\treasoning\t1\tfake\tm1\nproduct-owner\treasoning\t1\tfake\tm1\nsystems-architect\treasoning\t1\tfake\tm1\n" +
			"qa-engineer\treasoning\t1\tfake\tm1\nqa-engineer\treasoning\t2\tspare\tm1\nqa-engineer\treasoning\t3\tthird\tm1\n"
		if strings.Count(files["bin/harness"], "p=$(cat)\n") != 1 {
			t.Fatal("the harness of the world has no one place after p=$(cat)")
		}
		files["bin/harness"] = strings.Replace(files["bin/harness"], "p=$(cat)\n", "p=$(cat)\n"+loopSessions, 1)
	}
}

// eventsOfTask gives the events of a task from the records branch.
func eventsOfTask(t *testing.T, w runWorld, task string) [][]string {
	t.Helper()
	rows, err := records.ReadEvents([]byte(record(t, w, "tasks/"+task+"/events.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	return rows
}

// failedJob names the required job that the fake answers with a failure: none in the scenario, and the control C1 (NFR-004) sets it in a copy of the tree.
const failedJob = ""

func jobConclusion(job string) string {
	if failedJob != "" && job == failedJob {
		return "failure"
	}
	return "success"
}

// wantMerges is the number of merge requests the scenario expects: one for each task of Shape, and none when a required job is a failure.
func wantMerges() int {
	if failedJob != "" {
		return 0
	}
	return 2
}

// mergeCountFinding gives the finding of a count of merge requests against the expected one, "" for none.
func mergeCountFinding(got, want int) string {
	if got == want {
		return ""
	}
	return fmt.Sprintf("%d merge requests in the log; want exactly %d", got, want)
}

// notGreen gives the contexts of the four layup/ statuses whose newest status in the snapshot (newest first) is not the App's (creator 9001) and success.
func notGreen(snapshot []map[string]any) (bad []string) {
	for _, c := range []string{"layup/gates", "layup/spec", "layup/rules", "layup/verify"} {
		i := slices.IndexFunc(snapshot, func(s map[string]any) bool { return s["context"] == c })
		if i < 0 || snapshot[i]["state"] != "success" || snapshot[i]["creator"].(map[string]any)["id"] != 9001 {
			bad = append(bad, c)
		}
	}
	return bad
}

// loopMerged checks the merge of a task of Shape after the run, from the log of the forge, the bare repository and the records (condition 4 of the plan
// review): H is the head of the branch that the merge request named, M the merge commit that it made.
func loopMerged(t *testing.T, w runWorld, task layuprun.Task) (head string) {
	t.Helper()
	w.forge.mu.Lock()
	defer w.forge.mu.Unlock()
	var at []int
	for i, l := range w.forge.reqs {
		if strings.HasPrefix(l, fmt.Sprintf("merge %d ", task.PR)) {
			at = append(at, i)
			head = strings.Fields(l)[2]
		}
	}
	if len(at) != 1 || !regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(head) {
		t.Fatalf("%s: %d merge requests of the pull request %d, at the head %q; want one", task.ID, len(at), task.PR, head)
	}
	if tip, _ := gitRun(w.bare, "rev-parse", "refs/heads/task/"+task.ID+"/1"); strings.TrimSpace(tip) != head {
		t.Errorf("%s: the merge request's sha %s is not the head of the branch %s", task.ID, head, strings.TrimSpace(tip))
	}
	pre := w.forge.reqs[:at[0]]
	if !slices.Contains(pre, "mark PR_"+strconv.Itoa(task.PR)) || slices.Index(pre, "POST /repos/acme/target/statuses/"+head) < 0 {
		t.Errorf("%s: no ready mark or no status at %s before the merge in the log", task.ID, head)
	}
	if bad := notGreen(w.forge.m2e.atMerge[head]); len(bad) > 0 {
		t.Errorf("%s: at the merge request the newest statuses at %s of %q were not the App's success: %v", task.ID, head, bad, w.forge.m2e.atMerge[head])
	}
	merge := ""
	log, _ := gitRun(w.bare, "log", "--first-parent", "--merges", "--format=%H %P", "refs/heads/main")
	for _, l := range strings.Split(strings.TrimSpace(log), "\n") {
		if f := strings.Fields(l); len(f) == 3 && f[2] == head {
			merge = f[0]
		}
	}
	var merged []string
	for _, e := range eventsOfTask(t, w, task.ID) {
		switch {
		case e[1] == "merged":
			merged = e
		case e[1] == "handoff" && !strings.HasPrefix(e[6], "valid"):
			t.Errorf("%s: the handoff event %v is not valid", task.ID, e)
		}
	}
	if merge == "" || merged == nil || merged[5] != merge || merged[6] != strconv.Itoa(task.PR) {
		t.Errorf("%s: the merge commit with the second parent %s is %q; the merged event is %v; want M and the number %d", task.ID, head, merge, merged, task.PR)
	}
	return head
}

func TestTheSpecificationTaskEndToEnd(t *testing.T) {
	lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "prd", "prd-lint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	// The second pinned baseline holds the pinned prd-lint that the driver reads at the base (D2); one fact, and one decision whose option pg
	// passes the invariant Inv-1 with the evidence docs/guardrails.md and whose option sqlite fails it: one survivor and one drop (Q2).
	bin := buildBaselineBinary(t, &loopBuild, map[string]string{"docs/prd/prd-lint.sh": string(lint)})
	w, first := startScaffoldWorldWith(t, bin, "db\\tpg\\tInv-1\\tpass\\tdocs/guardrails.md\\ndb\\tpg\\tADR-0001\\tpass\\tdocs/adr/0001-pin-the-baseline.md\\ndb\\tsqlite\\tInv-1\\tfail\\t\\342\\200\\224\\n", loopWorldFiles(t))
	if first.code != 0 {
		t.Fatalf("layup run --new: exit %d\n%s\n%s", first.code, first.stdout, first.stderr)
	}
	leaks(t, "Start", w, first)
	env := goEnv(t)
	runs := 1 // the binary is run seven times: Start and six restarts
	restart := func(stage string, code int, names []string) (result, map[string][2]string) {
		t.Helper()
		runs++
		r := w.runWith(t, bin, env, "run", "acme/target", "--host", ".")
		leaks(t, stage, w, r)
		got, rows := table(r.stdout)
		if r.code != code || !slices.Equal(got, names) {
			t.Fatalf("%s: exit %d, the rows %q; want exit %d, %q\n%s\n%s", stage, r.code, got, code, names, r.stdout, r.stderr)
		}
		return r, rows
	}
	front := []string{"forge", "clone", "version", "lease"}
	intake := []string{"probe", "rule-gaps", "meaning", "numbering", "spans", "draft", "completeness", "markers", "batch", "answers", "answers-review", "follow-up", "follow-up-answers", "revision"}
	scaffold := []string{"setup", "verify", "scaffold-push"}
	later := slices.Concat(front, []string{"forge-check"}, intake, scaffold) // once the Scaffold is done, forge-check follows lease

	// Stage 1 to 3: the restart to the batch, the answers (and the Operator's set of brief.lines), the restart to the wait at scaffold-push.
	_, rows := restart("stage 1", 0, slices.Concat(front, intake[:10]))
	if rows["probe"][1] != "probed and passed 3, probed and failed 0, skipped 0" {
		t.Errorf("stage 1: the probe row %q; want the three harnesses passed", rows["probe"])
	}
	port, owner := setup.MarkerID("docs/ops.md", angled("port")), setup.MarkerID("docs/ops.md", angled("owner"))
	var idea, operator strings.Builder
	questions := questionsOf(t, w)
	for _, q := range questions {
		text := map[string]string{"band.B": "5.0", "band.U": "10.0", "appetite.hours": "40.0", "success-criteria": "the setup completes", "S01-stack": "go", port: "8080", owner: "gap"}[q[2]]
		if text == "" {
			text = "the answer to " + q[0]
		}
		if q[3] == "operator" {
			operator.WriteString(q[0] + ": " + text + "\n")
		} else {
			idea.WriteString(q[0] + ": " + text + "\n")
		}
	}
	answer(t, w, 901, questions[0], idea.String())
	answerBy(t, w, 903, 101, "pharzam", questions[0], operator.String())
	w.forge.mu.Lock()
	w.forge.people2 = append(w.forge.people2, map[string]any{"id": 904, "user": map[string]any{"id": 101, "login": "pharzam"},
		"created_at": "2026-10-08T13:30:00Z", "updated_at": "2026-10-08T13:30:00Z", "body": "set brief.lines 80 because the e2e brief\n"})
	w.forge.mu.Unlock()
	restart("stage 3", 0, slices.Concat(front, intake, scaffold))

	// Stage 4: the test is the Operator and runs the printed file; the restart after it ends at forge-check, as the tail is chosen before the check runs.
	sh := exec.Command("sh", filepath.Join(w.dir, "work", "acme", "target", "out", "commands.sh"))
	sh.Env = []string{"PATH=" + filepath.Join(w.dir, "ops") + ":" + os.Getenv("PATH"), "HOME=" + t.TempDir(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "LC_ALL=C"}
	if out, err := sh.CombinedOutput(); err != nil {
		t.Fatalf("out/commands.sh: %v\n%s", err, out)
	}
	restart("stage 4", 0, slices.Concat(front, intake, scaffold, []string{"forge-check"}))

	// Stage 5: the restart that runs Shape to the verified heads of its two tasks, and the bet to its wait for the answer on the brief.
	tail := slices.Concat(later, []string{"shape", "bet"})
	r, rows := restart("stage 5", 0, tail)
	reg, err := layuprun.ReadRegister([]byte(record(t, w, "tasks.tsv")))
	if err != nil || len(reg) != 2 || reg[0].Kind != "specification" || reg[1].Kind != "architecture" {
		t.Fatalf("tasks.tsv: %+v, %v; want the specification then the architecture task of Shape", reg, err)
	}
	spec, arch := reg[0], reg[1]
	briefs := func() (n int) {
		w.forge.mu.Lock()
		defer w.forge.mu.Unlock()
		for _, c := range w.forge.intake {
			if strings.HasPrefix(c, "## The problem\n") {
				n++
			}
		}
		return n
	}
	if rows["shape"] != [2]string{"done", "1 decisions, each with a survivor; " + spec.ID + " and " + arch.ID + " verified"} || !strings.HasPrefix(rows["bet"][1], "waiting: the answer of the bet on the brief ") ||
		briefs() != 1 || !strings.Contains(record(t, w, "parameters.tsv"), "\nbrief.lines\t80\t") {
		t.Fatalf("stage 5: the rows %q, %d briefs\n%s\n%s", rows, briefs(), r.stdout, r.stderr)
	}

	// Stage 6: the idea owner approves the brief; the restart goes on from the records to the merge of both tasks.
	w.forge.mu.Lock()
	w.forge.people = append(w.forge.people, map[string]any{"id": 905, "user": map[string]any{"id": 303, "login": "carol"},
		"created_at": "2026-10-08T14:00:00Z", "updated_at": "2026-10-08T14:00:00Z", "body": "REQ-001: Must, phase 1\nbet\n"})
	w.forge.mu.Unlock()
	r, rows = restart("stage 6", 0, slices.Concat(later, []string{"shape", "bet", "loop"}))
	// Two merge requests, one for each task, or exactly none under the control C1 (a required job is a failure, NFR-004): never one.
	w.forge.mu.Lock()
	merges := len(slices.DeleteFunc(slices.Clone(w.forge.reqs), func(l string) bool { return !strings.HasPrefix(l, "merge ") }))
	w.forge.mu.Unlock()
	if f := mergeCountFinding(merges, wantMerges()); f != "" {
		t.Errorf("stage 6: %s", f)
	}
	if rows["shape"][1] != "done before" || !strings.HasPrefix(rows["bet"][1], "the bet 1: the rendered commit ") || rows["loop"] != [2]string{"done", "the loop merged " + spec.ID + ", " + arch.ID} {
		t.Fatalf("stage 6: the rows %q\n%s\n%s", rows, r.stdout, r.stderr)
	}
	if lease := record(t, w, "lease.tsv"); !strings.HasSuffix(lease, "\treleased\n") {
		t.Errorf("lease.tsv: %q; want released", lease)
	}

	// The merges: each task merged once, at its head H, with the four statuses of the App at H, and M with the second parent H.
	reg, _ = layuprun.ReadRegister([]byte(record(t, w, "tasks.tsv")))
	heads := map[string]string{}
	for _, task := range reg {
		heads[task.ID] = loopMerged(t, w, task)
	}
	order, _ := gitRun(w.bare, "log", "--first-parent", "--merges", "--reverse", "--format=%P", "refs/heads/main")
	if lines := strings.Split(strings.TrimSpace(order), "\n"); len(lines) != 2 || !strings.HasSuffix(lines[0], " "+heads[spec.ID]) || !strings.HasSuffix(lines[1], " "+heads[arch.ID]) {
		t.Errorf("the merges of main in order: %q; want the specification %s, then the architecture %s", lines, heads[spec.ID], heads[arch.ID])
	}

	// The first bet: one row of bets.tsv with the verdict bet, the inventory, and the rendered commit R on the specification's branch, whose parent V is
	// the verified head: R changes the two cells of the requirement's row and adds the facts record, and nothing else (the cell boundary, on the bytes).
	if bets := strings.Split(strings.TrimSpace(record(t, w, "bets.tsv")), "\n"); len(bets) != 2 || !strings.Contains(bets[1], "\tbet\t") {
		t.Errorf("bets.tsv:\n%s\nwant one row with the verdict bet", bets)
	}
	record(t, w, "spec/inventory.tsv")
	rendered, _ := gitRun(w.bare, "log", "-1", "--format=%H", "--grep=^layup run: the rendered commit of the bet", "refs/heads/task/"+spec.ID+"/1")
	rendered = strings.TrimSpace(rendered)
	diff, _ := gitRun(w.bare, "diff", "--name-status", rendered+"^", rendered)
	if !regexp.MustCompile(`^A\tdocs/facts/F-\d{4}-numbered-facts\.md\nM\tdocs/prd/PRD-0001-x\.md\n$`).MatchString(diff) {
		t.Fatalf("the rendered commit %q changes\n%s\nwant the PRD and a new numbered facts record", rendered, diff)
	}
	before, _ := gitRun(w.bare, "show", rendered+"^:docs/prd/PRD-0001-x.md")
	after, _ := gitRun(w.bare, "show", rendered+":docs/prd/PRD-0001-x.md")
	changed, was, now := 0, strings.Split(before, "\n"), strings.Split(after, "\n")
	if len(was) != len(now) {
		t.Fatalf("the rendered commit changes the number of lines of the PRD: %d to %d", len(was), len(now))
	}
	for i := range was {
		if was[i] == now[i] {
			continue
		}
		changed++
		b, a := strings.Split(was[i], "|"), strings.Split(now[i], "|")
		if len(b) != len(a) || len(b) < 6 || strings.TrimSpace(b[1]) != "REQ-001" || strings.TrimSpace(a[3]) != "Must" || strings.TrimSpace(a[4]) != "1" || !slices.Equal(slices.Concat(b[:3], b[5:]), slices.Concat(a[:3], a[5:])) {
			t.Errorf("line %d of the PRD changed outside the cells MoSCoW and Phase of REQ-001:\n%s\n%s", i+1, was[i], now[i])
		}
	}
	if main, _ := gitRun(w.bare, "show", "refs/heads/main:docs/prd/PRD-0001-x.md"); changed != 1 || main != after {
		t.Errorf("%d lines changed by the rendered commit; the PRD of main is the rendered one: %v", changed, main == after)
	}

	// Stage 7: a repeat restart records "done before" and makes no write to the forge (a request other than a read and the issue of a token).
	w.forge.mu.Lock()
	seen, comments := len(w.forge.reqs), briefs
	w.forge.mu.Unlock()
	_, rows = restart("stage 7", 0, slices.Concat(later, []string{"shape", "bet", "loop"}))
	w.forge.mu.Lock()
	writes := slices.DeleteFunc(slices.Clone(w.forge.reqs[seen:]), func(l string) bool {
		return strings.HasPrefix(l, "GET ") || l == "POST /app/installations/7/access_tokens"
	})
	w.forge.mu.Unlock()
	if rows["loop"] != [2]string{"done", "done before"} || len(writes) != 0 || comments() != 1 {
		t.Errorf("stage 7: the loop row %q, the writes %q, %d briefs; want done before and none", rows["loop"], writes, comments())
	}
	if runs != 7 {
		t.Errorf("the binary was run %d times; want 7", runs)
	}

	// No secret: no session exited 9 (the harness refuses a token in its environment), and no token is in the records.
	if out, err := gitRun(w.bare, "grep", "-q", "ghs_", "refs/heads/layup-records"); err == nil {
		t.Errorf("a token of the forge is in the records branch:\n%s", out)
	}
}

// The oracle of the merge refuses a trace in which only layup/verify was posted before the merge request and the other three contexts after it: the snapshot
// at the request holds one status (review round 1 of #368, finding 1).
func TestTheMergeOracleRefusesStatusesPostedAfterTheMergeRequest(t *testing.T) {
	app := map[string]any{"id": 9001, "login": "layup-agent[bot]"}
	st := func(c, state string, creator map[string]any) map[string]any {
		return map[string]any{"context": c, "state": state, "creator": creator}
	}
	if bad := notGreen([]map[string]any{st("layup/verify", "success", app)}); !slices.Equal(bad, []string{"layup/gates", "layup/spec", "layup/rules"}) {
		t.Errorf("a snapshot with layup/verify alone: %q; want the other three contexts", bad)
	}
	all := []map[string]any{st("layup/verify", "success", app), st("layup/rules", "success", app), st("layup/spec", "success", app), st("layup/gates", "success", app)}
	if bad := notGreen(all); len(bad) != 0 {
		t.Errorf("a snapshot of the four contexts at success: %q; want none", bad)
	}
	other := map[string]any{"id": 5, "login": "x"}
	if bad := notGreen(append([]map[string]any{st("layup/gates", "failure", app), st("layup/spec", "success", other)}, all...)); !slices.Equal(bad, []string{"layup/gates", "layup/spec"}) {
		t.Errorf("a newer failure and a newer status of another creator: %q; want those two contexts", bad)
	}
}

// The count of merge requests is exact: the control C1 accepts zero and refuses a single request, the scenario accepts two and refuses one (round 2 of #368).
func TestTheCountOfMergeRequestsIsExact(t *testing.T) {
	for _, c := range []struct {
		got, want int
		bad       bool
	}{{0, 0, false}, {1, 0, true}, {2, 0, true}, {2, 2, false}, {1, 2, true}, {0, 2, true}} {
		if f := mergeCountFinding(c.got, c.want); (f != "") != c.bad {
			t.Errorf("%d merge requests, %d expected: finding %q; want a finding: %v", c.got, c.want, f, c.bad)
		}
	}
	if failedJob == "" && wantMerges() != 2 {
		t.Errorf("the scenario expects %d merge requests; want 2", wantMerges())
	}
}
