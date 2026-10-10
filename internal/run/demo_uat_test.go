//go:build uat

package run

// The uat test of M2b (docs/spec/session.md, The acceptance tests of M2b, The
// demo; O-180 of #147; task T-x7cs, #170). It runs on the host with the
// Operator's real harnesses and credentials and writes to a real target, so no
// CI job and no documented command of the ladder runs it; a person runs it:
//
//	go test -tags=uat -run TestTheDemoOfM2b -timeout 60m -v ./internal/run -args -host DIR -target OWNER/NAME
//
// It runs the steps of the restart on one run, and between probe and phase,
// with the lease held, the developer session of a small fixed task.

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/forge/github"
	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
	"github.com/pharzam/layup/internal/rules"
)

var (
	uatHost    = flag.String("host", "", "the host directory of layup run")
	uatTarget  = flag.String("target", "", "the target OWNER/NAME, started")
	uatVersion = flag.String("version", "0.1.0-dev", "the LAYUP version that started the target")
)

type uatClock struct{}

func (uatClock) Now() time.Time { return time.Now() }
func (uatClock) Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// uatConfig reads the host's registers with the readers that internal/cli
// uses (internal/cli/run.go, hostInputs), which this package cannot import.
func uatConfig(t *testing.T) Config {
	t.Helper()
	owner, name, ok := strings.Cut(*uatTarget, "/")
	if *uatHost == "" || !ok {
		t.Fatal("give -host DIR -target OWNER/NAME")
	}
	dir, err := filepath.Abs(*uatHost)
	if err != nil {
		t.Fatal(err)
	}
	file := func(parts ...string) []byte {
		data, err := os.ReadFile(filepath.Join(append([]string{dir}, parts...)...))
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	reg, err := forge.ReadForgeRegister(file("registers", "forge.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	harnesses, ids, err := route.ReadHarnesses(file("registers", "harnesses.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	models, err := route.ReadModels(file("registers", "models.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	routing, err := route.ReadRoutingRegister(file("registers", "routing.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	if err := route.CheckRegisters(ids, models, routing); err != nil {
		t.Fatal(err)
	}
	for _, h := range harnesses {
		if h[regCredential] != "" {
			if err := route.CheckCredential(h[regCredential]); err != nil {
				t.Fatal(err)
			}
		}
	}
	prices, err := ReadPrices(dir)
	if err != nil {
		t.Fatal(err)
	}
	key, err := forge.CheckKeyFile(reg.KeyFile)
	if err != nil {
		t.Fatal(err)
	}
	host, err := os.Hostname()
	if err != nil {
		t.Fatal(err)
	}
	progress := func(line string) { t.Log(line) }
	return Config{Owner: owner, Name: name, Dir: dir, Register: reg, Harnesses: harnesses, Models: models, Routing: routing,
		Prices: prices, Clock: uatClock{}, Progress: progress, RunID: fmt.Sprintf("%016x", time.Now().UnixNano()),
		HostName: host, Version: *uatVersion,
		Forge: github.New(github.Config{API: reg.API, AppID: reg.AppID, Key: key, Owner: owner, Name: name, Progress: progress})}
}

// demoRulePaths are the entries of phase 1 that are fixed paths
// (docs/spec/setup.md, The rule-path register), the register of a target that
// layup run --new made, which has none until layup setup writes it (task
// T-x7cs, D2 of the plan of #170).
var demoRulePaths = [][]string{
	{".github/", "—", "baseline"}, {".githooks/", "—", "baseline"}, {".gitattributes", "—", "baseline"},
	{"AGENTS.md", "—", "baseline"}, {"CLAUDE.md", "—", "baseline"}, {"docs/engineering-discipline.md", "—", "baseline"},
	{"docs/issue-workflow.md", "—", "baseline"}, {"docs/ci/", "—", "baseline"}, {"docs/tests/", "—", "baseline"},
	{"docs/gates.tsv", "—", "catalog"}, {"docs/setup/", "—", "architecture"}, {"docs/facts/", "—", "architecture"},
	{"docs/guardrails.md", "added lines in section 2", "baseline"},
}

// demoPrompt is the fixed task of the developer session; {task} is its ID.
const demoPrompt = `This is the demo of LAYUP M2b, task {task}. Do only this:

1. Create the file demo/{task}.txt in this repository, with the one line "The demo of LAYUP M2b, task {task}." and a final line feed.
2. Commit it on the current branch with the message "docs: {task} the file of the demo". Do not push.
3. Write the file ../result/result.tsv (the directory result/ beside this repository), of tab-separated values with a final line feed, with the header line
"kind<TAB>n<TAB>value<TAB>sha256<TAB>reason" and two rows:
"status<TAB>1<TAB>completed<TAB>—<TAB>the file is written and committed", and
"artifact<TAB>1<TAB>demo/{task}.txt<TAB><sha><TAB>—", where <sha> is the SHA-256 of demo/{task}.txt in lowercase hexadecimal (sha256sum prints it).

Then end.
`

func TestTheDemoOfM2b(t *testing.T) {
	cfg := uatConfig(t)
	r := &state{cfg: cfg}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r.cancel = cancel
	defer r.endBeat()
	step := func(s Step) {
		t.Helper()
		t.Logf("step %s: %s, %s", s.Name, s.Result, s.Detail)
		if s.Result != "done" {
			t.Fatalf("the step %s: %s", s.Name, s.Detail)
		}
	}
	step(r.forgeStep(ctx, false))
	step(r.cloneStep(ctx))
	step(r.versionStep(ctx))
	step(r.takeStep(ctx))
	// A test that stops before phase still releases the lease, so the target
	// is not left held until the takeover rule.
	phased := false
	defer func() {
		if !phased {
			s := r.leaseStep(context.Background())
			t.Logf("the lease, released after a stop: %s, %s", s.Result, s.Detail)
		}
	}()
	step(r.probeStep(ctx))

	// The probe of each registered harness passed at its version.
	control, _ := strconv.Atoi(r.value("issue.control"))
	s := &Sessions{Store: r.store, RunID: cfg.RunID, Host: cfg.Dir, Target: cfg.Owner + "/" + cfg.Name, Clone: r.target,
		Branch: r.defaultBranch, Now: cfg.Clock.Now, Forge: cfg.Forge, Control: control, Prices: cfg.Prices}
	held, err := s.readTable(ctx, "harnesses.tsv", records.ReadHarnesses)
	if err != nil {
		t.Fatal(err)
	}
	versions := map[string]string{}
	for _, h := range held {
		versions[h[1]] = h[2] // the last row of each harness, in the order of the file
	}
	for _, h := range cfg.Harnesses {
		if v := versions[h[regHarness]]; v == "" || !route.Probed(held, h[regHarness], v) {
			t.Fatalf("the harness %s: no probe that passed at its version (%q); harnesses.tsv %q", h[regHarness], v, held)
		}
	}

	// The rule-path register, when the records branch holds none (D2).
	if data, err := s.Store.ReadFile(ctx, "rule-paths.tsv"); err != nil {
		t.Fatal(err)
	} else if data == nil {
		data, err := table(rules.RulePathsSchema, demoRulePaths)
		if err != nil {
			t.Fatal(err)
		}
		if err := r.write(ctx, map[string][]byte{"rule-paths.tsv": data}, "layup run: the rule-path register of the demo of M2b"); err != nil {
			t.Fatal(err)
		}
	}

	// The developer session of the fixed task (D3).
	base, err := git.RevParse(r.target, "refs/remotes/origin/"+r.defaultBranch)
	if err != nil {
		t.Fatal(err)
	}
	task, err := s.probeTaskID(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if n, err := s.StartAttempt(ctx, task, base); err != nil || n != 1 {
		t.Fatalf("StartAttempt: %d, %v", n, err)
	}
	harness, model, err := route.Pair(cfg.Routing, "developer", "execution", func(h, m string) bool {
		return route.Admitted(held, cfg.Models, h, versions[h], m)
	})
	if err != nil {
		t.Fatal(err)
	}
	row := cfg.Harnesses[slices.IndexFunc(cfg.Harnesses, func(h []string) bool { return h[regHarness] == harness })]
	modelRow := cfg.Models[slices.IndexFunc(cfg.Models, func(m []string) bool { return m[0] == harness && m[1] == model })]
	recs := r.store.Base()
	spec := TaskSpec{Task: task, Attempt: 1, Role: "developer", Base: base, Records: recs,
		Prompt: []byte(strings.ReplaceAll(demoPrompt, "{task}", task)), Pair: pairOf(row, modelRow),
		Admit: s.Admit(ctx, ProbeSpec{Harness: row, Models: cfg.Models, Base: base, Records: recs})}
	t.Logf("the task %s, attempt 1, the pair %s %s, the base %s", task, harness, model, base)
	id, err := s.TaskSession(ctx, spec)
	t.Logf("the session %s: %v", id, err)
	if err != nil {
		t.Fatal(err)
	}

	// Its events, its branch on the target and its telemetry row (D4).
	events, err := s.readTable(ctx, eventsPath(task), records.ReadEvents)
	if err != nil {
		t.Fatal(err)
	}
	var kinds []string
	for _, e := range events {
		kinds = append(kinds, e[1])
		t.Logf("event %q", e)
	}
	if want := []string{"attempt", "session", "result", "push", "bound"}; !slices.Equal(kinds, want) {
		t.Errorf("the events %q, want %q", kinds, want)
	}
	branch := "refs/heads/task/" + task + "/1"
	head, err := git.LsRemote(r.url(), branch)
	t.Logf("%s on the target: %s, %v", branch, head, err)
	if err != nil || head == "" {
		t.Errorf("no %s on the target", branch)
	}
	telemetry, err := s.readTable(ctx, "telemetry.tsv", records.ReadTelemetry)
	if err != nil {
		t.Fatal(err)
	}
	if i := slices.IndexFunc(telemetry, func(row []string) bool { return row[0] == id }); i < 0 {
		t.Errorf("no telemetry row of %s", id)
	} else {
		t.Logf("telemetry %q", telemetry[i])
	}

	phased = true
	step(r.phaseStep(ctx))
}
