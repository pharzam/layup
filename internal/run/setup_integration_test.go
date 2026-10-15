//go:build integration

package run

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/verify"
	workpkg "github.com/pharzam/layup/internal/work"
)

// The step setup on the real step runner (docs/spec/scaffold.md, The answers of
// the setup; The work area; The stops of the runner; The prose rows; the
// acceptance tests of M2d; row 54 of the plan, task T-eahr, #277). The world's
// baseline is the stand-in baseline of internal/standin with the files that S05
// to S14 act on (as internal/cli's test of S05 to S15 has them), and the world's
// scripted harness writes the prose files from bin/files/ for the F- rows of the
// stop table that its prompt names.

// The two markers of the baseline, with the angle quotes that angled gives (markers_test.go).
var (
	portMark  = angled("port")
	ownerMark = angled("owner")
)

// scaffoldBaseline makes the baseline of the world at a new directory: the
// stand-in baseline, the link-lint.sh of this repository (the baseline's own, as
// LAYUP keeps it), a note of the kit in a blockquote of the backlog, a guide
// whose link the deletion of the history breaks, two markers (one twice in a
// file), and two files that check adapted flags. It gives the directory and the
// commit.
func scaffoldBaseline(t *testing.T) (string, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "baseline")
	url, _, err := standin.Baseline(dir)
	if err != nil {
		t.Fatal(err)
	}
	lint, err := os.ReadFile(filepath.Join("..", "..", "docs", "links", "link-lint.sh"))
	if err != nil {
		t.Fatal(err)
	}
	for p, text := range map[string]string{
		"docs/links/link-lint.sh": string(lint),
		"docs/tasks/backlog.md":   "# Backlog\n\n- **T-0001**: a task of the baseline ([#1](" + url + "/issues/1))\n\n> a note of the kit\n> with a link ([#2](" + url + "/issues/2))\n\nKeep this line.\n",
		"docs/guide.md":           "# Guide\n\nSee [the decision](decisions/D-0001-stand-in.md).\n",
		"docs/ops.md":             "# Ops\n\nThe port is " + portMark + ".\nAgain " + portMark + ".\nThe owner is " + ownerMark + ".\n",
		"docs/how-to.md":          "# How to\n\nAdapt this to your project.\n",
		"docs/glossary.md":        "# Glossary\n\nThe words of the kit.\n",
	} {
		writeFile(t, filepath.Join(dir, filepath.FromSlash(p)), text)
	}
	// One commit: the world's test of the Operator pushes the HEAD of the root clone as soon as
	// it exists, which can be the baseline's own commit, and the default branch must be one root commit.
	if err := os.RemoveAll(filepath.Join(dir, ".git")); err != nil {
		t.Fatal(err)
	}
	if err := git.Init(dir); err != nil {
		t.Fatal(err)
	}
	if err := git.Add(dir); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(dir, "the stand-in baseline with the scaffold", git.Identity{Name: "t", Email: "t@layup.invalid", Time: time.Unix(0, 0)}); err != nil {
		t.Fatal(err)
	}
	commit, err := git.RevParse(dir, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	return dir, commit
}

// proseFiles are the files that the world's session writes, by path, for the
// stops of the setup: the guide that S05 asks for, and the six of the prose step.
func proseFiles() map[string]string {
	const name = "acme/target"
	return map[string]string{
		"docs/guide.md":                    "# Guide\n\nThe decisions of the target are in its records.\n",
		"docs/onboarding-for-engineers.md": "# Onboarding\n\nThe work starts from [the problem statement](facts/problem-statement-brief.md).\n",
		"docs/guardrails.md":               "# Guardrails\n\n## 1. Rules\n\n- **Inv-1** — the setup keeps its record. Check: no check yet\n",
		"docs/glossary.md":                 "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository |\n",
		"docs/how-to.md":                   "# How to\n\nRun the tests of the product.\n",
		"README.md":                        "# " + name + "\n\nThe product repository, set up from its pinned baseline ([the pin](docs/setup/armature.pin)). Its records are on the branch `layup-records`.\n",
		"AGENTS.md":                        "# AGENTS.md\n\nAgent context for **" + name + "**.\n",
	}
}

// scaffoldWorld is a world on the scaffold baseline, with the prose files in bin/files/. The
// step verify of row 55 has its own test: in this world layup setup verify is not run (its call
// is an error), so the step verify is fail and no record of the setup is written, as before it.
func scaffoldWorld(t *testing.T, files map[string]string) *world {
	t.Helper()
	saved := verifyRun
	t.Cleanup(func() { verifyRun = saved })
	verifyRun = func(string, func(int, int, string) func(), io.Writer) (verify.Table, error) {
		return verify.Table{}, errors.New("not run in this world")
	}
	w := newWorld(t)
	w.baseline, w.pin = scaffoldBaseline(t)
	for p, text := range files {
		writeFile(t, filepath.Join(w.dir, "bin", "files", filepath.FromSlash(p)), text)
	}
	return w
}

// toSetup runs Start, the restart to the batch, the answers of the idea owner and of
// the Operator (the stack as given, the port 8080 and the owner a gap), and the
// restart after them: Intake ends, and the step setup runs. It gives the rows of
// that restart.
func toSetup(t *testing.T, w *world, stack string) []Step {
	t.Helper()
	allDone(t, "Start", w.start(w.config("0.1.0-dev")), "forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	if d := rowOf(rows, "answers").Detail; !strings.HasPrefix(d, "waiting") {
		t.Fatalf("the restart to the batch: %q; want the step answers waiting\n%s", d, steps(rows))
	}
	questions, err := records.ReadQuestions([]byte(show(t, w.bare, "refs/heads/layup-records", "questions.tsv")))
	if err != nil {
		t.Fatal(err)
	}
	asked, err := time.Parse(time.RFC3339, questions[0][9])
	if err != nil {
		t.Fatal(err)
	}
	at := func(d time.Duration) string { return asked.Add(d).UTC().Format(time.RFC3339) }
	idPort, idOwner := setup.MarkerID("docs/ops.md", portMark), setup.MarkerID("docs/ops.md", ownerMark)
	var owner, operator strings.Builder
	for _, q := range questions {
		answer := "the answer to " + q[0]
		switch {
		case q[2] == "band.B":
			answer = "9.5"
		case q[2] == "band.U":
			answer = "10.0"
		case q[2] == "appetite.hours":
			answer = "40.0"
		case q[2] == "S01-stack":
			answer = stack
		case q[2] == idPort:
			answer = "8080"
		case q[2] == idOwner:
			answer = "gap"
		}
		if q[3] == "operator" {
			operator.WriteString(q[0] + ": " + answer + "\n")
		} else {
			owner.WriteString(q[0] + ": " + answer + "\n")
		}
	}
	addComment(w, 1, 802, 303, "carol", at(2*time.Minute), owner.String())
	addComment(w, 1, 803, 101, "pharzam", at(3*time.Minute), operator.String())
	return Restart(context.Background(), w.config("0.1.0-dev"))
}

// workTarget is the target/ of the work area of the world.
func workTarget(w *world) string { return filepath.Join(w.dir, "work", "acme", "target", "target") }

// setupSessions counts the comments on the control issue that say a session of the step setup was done.
func setupSessions(w *world) int {
	w.fake.mu.Lock()
	defer w.fake.mu.Unlock()
	n := 0
	for _, c := range w.fake.posted {
		n += strings.Count(c, " session of setup: ")
	}
	return n
}

// The demo of row 54: after Intake, a restart builds the work area from the
// records, answers each stop of the runner from them (the markers of S10) or with a
// prose session (S05's guide, then the six files of the prose step), and runs S01
// to S14: the step setup is done, target/ holds layup-setup with the commits of S04
// to S13, origin of target/ is the forge, the texts are in the records, and a second
// restart builds the work area again to the same head and starts no session, and
// a restart with setup/record.tsv in the records finds the step done before.
func TestTheStepSetupRunsS01ToS14(t *testing.T) {
	w := scaffoldWorld(t, proseFiles())
	rows := toSetup(t, w, "go")
	if last := rowOf(rows, "setup"); last != (Step{"setup", "done", "S01 to S14 done, in the work area"}) {
		t.Fatalf("the row of the step setup: %+v; want setup done\n%s", last, steps(rows))
	}
	target := workTarget(w)
	if got := gitOut(t, target, "log", "--reverse", "--format=%s", "main..layup-setup"); got != "chore: setup S04\nchore: setup S05\nchore: setup S06\nchore: setup S07\nchore: setup S08\n"+
		"chore: setup S09\nchore: setup S14\nchore: setup S11\nchore: setup S12\nchore: setup S13\n" {
		t.Errorf("the setup commits: %q", got)
	}
	if got, want := strings.TrimSpace(gitOut(t, target, "remote", "get-url", "origin")), w.config("0.1.0-dev").Register.Web+"/acme/target.git"; got != want {
		t.Errorf("origin of target/: %q; want %q", got, want)
	}
	for path, want := range proseFiles() {
		if got := show(t, w.bare, "refs/heads/layup-records", "setup/inputs/files/"+path); got != want {
			t.Errorf("setup/inputs/files/%s in the records: %q; want %q", path, got, want)
		}
		if got, err := os.ReadFile(filepath.Join(w.dir, "work", "acme", "target", "inputs", "files", filepath.FromSlash(path))); err != nil || string(got) != want {
			t.Errorf("inputs/files/%s of the work area: %q, %v; want %q", path, got, err, want)
		}
	}
	answers, err := workpkg.ReadAnswers(filepath.Join(w.dir, "work", "acme", "target"))
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, a := range answers {
		got[a[0]] = a[1]
	}
	if got["S01-stack"] != "go" || got[setup.MarkerID("docs/ops.md", portMark)] != "8080" || got[setup.MarkerID("docs/ops.md", ownerMark)] != "gap" {
		t.Errorf("inputs/answers.tsv: %q; want go for S01-stack, 8080 for the port and a gap for the owner", got)
	}
	if n := setupSessions(w); n != 2 {
		t.Errorf("%d sessions of the step setup; want 2 (the guide of S05, the six files of the prose step)", n)
	}
	// The prompt that the last prose session read has the rules of the checks and the name of this world's target (task T-92lm, #326).
	prompt, err := os.ReadFile(filepath.Join(w.dir, "bin", "prose-prompt.txt"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"the name of this repository, acme/target, as text", "must name the branch layup-records", `"kit", "kits"`, "Before you end, search each file that you wrote"} {
		if !strings.Contains(strings.Join(strings.Fields(string(prompt)), " "), want) {
			t.Errorf("the prompt of the prose session does not say %q:\n%s", want, prompt)
		}
	}
	head := gitOut(t, target, "rev-parse", "layup-setup")
	sessions := setupSessions(w)
	t.Run("a second restart before verify rebuilds the same head and starts no session", func(t *testing.T) {
		again := Restart(context.Background(), w.config("0.1.0-dev"))
		if last := rowOf(again, "setup"); last != (Step{"setup", "done", "S01 to S14 done, in the work area"}) {
			t.Fatalf("the row of the step setup of a second restart: %+v; want setup done\n%s", last, steps(again))
		}
		if now := gitOut(t, target, "rev-parse", "layup-setup"); now != head || setupSessions(w) != sessions {
			t.Errorf("a second restart: the head %s (was %s), %d sessions (were %d); want the same head and no session", strings.TrimSpace(now), strings.TrimSpace(head), setupSessions(w), sessions)
		}
	})
	// With setup/record.tsv in the records the step is "done before": it builds no work area,
	// runs no runner and starts no session (the step verify writes the file).
	t.Run("setup is done before when the record exists", func(t *testing.T) {
		clone := filepath.Join(t.TempDir(), "records")
		gitOut(t, "/", "clone", "-q", "--branch", "layup-records", "--single-branch", "file://"+w.bare, clone)
		var b bytes.Buffer
		if err := tsv.Write(&b, workpkg.RecordSchema, [][]string{{"S01", "done", "S01-stack answered", "step", ""}}); err != nil {
			t.Fatal(err)
		}
		writeFile(t, filepath.Join(clone, "setup", "record.tsv"), b.String())
		gitOut(t, clone, "add", "-A")
		gitOut(t, clone, "commit", "-q", "-m", "the record of the setup, as the step verify writes it")
		gitOut(t, clone, "push", "-q", "origin", "layup-records")
		if err := os.RemoveAll(filepath.Join(w.dir, "work")); err != nil {
			t.Fatal(err)
		}
		again := Restart(context.Background(), w.config("0.1.0-dev"))
		if last := rowOf(again, "setup"); last != (Step{"setup", "done", "done before"}) {
			t.Fatalf("the row of the step setup of the restart: %+v; want setup done before\n%s", last, steps(again))
		}
		if _, err := os.Stat(filepath.Join(w.dir, "work")); err == nil || setupSessions(w) != sessions {
			t.Errorf("a step done before: the work area exists (%v) or a session started (%d, was %d); want neither", err, setupSessions(w), sessions)
		}
	})
}

// A prose text whose check fails (its citation resolves to no fact): the step is fail, naming the check; the files of that
// session are removed from the records and the files of the earlier session, whose
// row the runner passed, stay; with no admitted pair of product-owner the next
// restart is fail (pair) and starts no session; and a restart with the text put right
// starts a new session and ends done.
func TestAProseTextWhoseCheckFailsIsRemovedAndANewSessionStarts(t *testing.T) {
	files := proseFiles()
	files["docs/glossary.md"] = "# Glossary\n\n| Term | Meaning |\n| ---- | ------- |\n| Target | the product repository (`F-0099#1`) |\n" // a citation that resolves to no fact
	w := scaffoldWorld(t, files)
	rows := toSetup(t, w, "go")
	last := rows[len(rows)-1]
	if last.Name != "setup" || last.Result != "fail" || !strings.Contains(last.Detail, "glossary") {
		t.Fatalf("the last row of the restart: %+v; want setup fail, naming the check glossary\n%s", last, steps(rows))
	}
	if got := show(t, w.bare, "refs/heads/layup-records", "setup/inputs/files/docs/guide.md"); got != proseFiles()["docs/guide.md"] {
		t.Errorf("the guide of the first session: %q; want it kept, as the runner passed its row", got)
	}
	for path := range proseFiles() {
		if path == "docs/guide.md" {
			continue
		}
		if out, err := gitRunOK(w.bare, "show", "refs/heads/layup-records:setup/inputs/files/"+path); err == nil {
			t.Errorf("setup/inputs/files/%s is in the records after the fail: %q; want the files of the failed session removed", path, out)
		}
	}
	if log := gitOut(t, w.bare, "log", "--format=%s", "refs/heads/layup-records"); !strings.Contains(log, "layup run: the files of the session S-") {
		t.Errorf("the log of the records has no commit that removes the files of the session:\n%s", log)
	}
	if n := setupSessions(w); n != 2 {
		t.Errorf("%d sessions after the fail; want 2", n)
	}
	// No admitted pair of product-owner: the host's routing register has none (the step probe copies it into the records).
	cfg := w.config("0.1.0-dev")
	var routing [][]string
	for _, r := range cfg.Routing {
		if r[0] != "product-owner" {
			routing = append(routing, r)
		}
	}
	cfg.Routing = routing
	rows = Restart(context.Background(), cfg)
	if last := rows[len(rows)-1]; last.Name != "setup" || last.Result != "fail" || !strings.HasPrefix(last.Detail, "pair: ") || setupSessions(w) != 2 {
		t.Errorf("a restart with no pair of product-owner: %+v, %d sessions; want setup fail with pair, and no session", last, setupSessions(w))
	}
	// The pair back and the text put right: a new session, and done.
	writeFile(t, filepath.Join(w.dir, "bin", "files", "docs", "glossary.md"), proseFiles()["docs/glossary.md"])
	rows = Restart(context.Background(), w.config("0.1.0-dev"))
	if last := rowOf(rows, "setup"); last != (Step{"setup", "done", "S01 to S14 done, in the work area"}) || setupSessions(w) != 3 {
		t.Errorf("the restart with the text put right: %+v, %d sessions; want setup done after a third session\n%s", last, setupSessions(w), steps(rows))
	}
}

// A prose session that ends after its records commit (its comment on the control issue is refused): the step is fail,
// and the files of that session are removed from the records, as for any fail that follows a session.
func TestAProseSessionWhoseCommentFailsLeavesNoFileInTheRecords(t *testing.T) {
	w := scaffoldWorld(t, proseFiles())
	w.fake.failNote = " session of setup: "
	rows := toSetup(t, w, "go")
	if last := rows[len(rows)-1]; last.Name != "setup" || last.Result != "fail" {
		t.Fatalf("the last row of the restart: %+v; want setup fail\n%s", last, steps(rows))
	}
	if out, err := gitRunOK(w.bare, "show", "refs/heads/layup-records:setup/inputs/files/docs/guide.md"); err == nil {
		t.Errorf("docs/guide.md is in the records after the fail: %q; want the files of the session removed", out)
	}
	if log := gitOut(t, w.bare, "log", "--format=%s", "refs/heads/layup-records"); !strings.Contains(log, "layup run: the files of the session S-") {
		t.Errorf("the log of the records has no commit that removes the files of the session:\n%s", log)
	}
}

// The real runner on a work area that layup run built, with no prose session: the
// two input states of S01. With no row of S01-stack the runner stops on it, the step
// is fail and names it, and no setup commit is kept: layup-setup does not exist. With
// a stack that the catalog does not have S01 is fail with its own reason, and the same.
func TestTheRealRunnerOnTheStackAnswersOfTheStepSetup(t *testing.T) {
	for _, c := range []struct {
		name    string
		answers [][]string
		want    string
	}{
		{"no answer to S01-stack", nil, "S01-stack"},
		{"a stack that the catalog does not have", [][]string{{"S01-stack", "cobol", "operator", "https://github.invalid/acme/target/issues/1#issuecomment-5", ""}}, "S01: the stack cobol has no entry in the catalog"},
	} {
		bare, ids := forgeRepo(t, "main", map[string]string{"README.md": "# Root\n"})
		dir := filepath.Join(t.TempDir(), "work", "acme", "target")
		if err := buildWorkArea(dir, workSource{Brief: []byte("The problem.\n"), Answers: c.answers}, runClone(t, bare), "main", originURL); err != nil {
			t.Fatal(err)
		}
		p := setup.Phase2{Name: "acme/target", Visibility: "public", PinSource: "file:///no/such/baseline", PinCommit: ids[0],
			PinTree: strings.TrimSpace(gitOut(t, bare, "rev-parse", ids[0]+"^{tree}")), PinTime: "2026-10-13T05:00:00Z", BriefSum: "0123", AppID: 42, Host: t.TempDir(), DefaultBranch: "main"}
		l := setupLoop{run: setupRunner(dir, p, nil), max: 4,
			markers: func(string) ([]string, bool) { return nil, false },
			answer:  func([][]string) error { return nil },
			prose: func(context.Context, []setup.StopRow) (proseDone, error) {
				t.Error("a session started")
				return proseDone{}, nil
			},
			remove: func(context.Context, proseDone, string) error { t.Error("a removal"); return nil }}
		if got := l.do(context.Background()); got.Result != "fail" || !strings.Contains(got.Detail, c.want) {
			t.Errorf("%s: %+v; want fail with %q", c.name, got, c.want)
		}
		if out, err := gitRunOK(filepath.Join(dir, "target"), "rev-parse", "--verify", "refs/heads/layup-setup"); err == nil {
			t.Errorf("%s: layup-setup exists (%s); want no setup commit kept", c.name, strings.TrimSpace(out))
		}
	}
}

// gitRunOK runs git in dir, as gitOut does, and gives its output and its error, for a call that may fail.
func gitRunOK(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	return string(out), err
}
