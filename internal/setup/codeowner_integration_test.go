//go:build integration

package setup

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/work"
)

// The commands that update a ruleset for a target set up before the code-owner rule, and the script that lands
// its CODEOWNERS (docs/spec/scaffold.md, The code-owner rule of a target that M2d set up; task T-jew3, #356). They
// run with the real sh, git and gh (the Operator's own tools) against an httptest TLS server that serves the live
// object of the first-hand record runs/T-sdtb/rulesets.txt, as the forge does: the fields that the file never held,
// a write that bumps updated_at. The tests need gh on PATH and fail, never skip, when it is missing.

const rulesetID = 25145510

// liveFor gives the live object as the repository of these tests (acme/widget) would hold it: the source differs, nothing else.
func liveFor(t *testing.T) map[string]any {
	m := liveRuleset(t)
	m["source"] = "acme/widget"
	return m
}

// forgeFake is the ruleset endpoint of GitHub: GET gives the object, PUT replaces the fields of a write that it
// gets, as a real forge does, and both log "METHOD PATH". failPut makes the n-th PUT (from 1) answer 500; drop
// removes the bypass actors from the object after a PUT, as a forge that dropped a field would.
type forgeFake struct {
	mu      sync.Mutex
	obj     map[string]any
	calls   []string
	bodies  [][]byte
	puts    int
	failPut int
	failGet bool
	drop    bool
	env     []string
}

func newForge(t *testing.T, obj map[string]any) *forgeFake {
	t.Helper()
	if _, err := exec.LookPath("gh"); err != nil {
		t.Fatalf("the tests need gh on PATH (docs/tests/test-levels.md): %v", err)
	}
	f := &forgeFake{obj: obj}
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.calls = append(f.calls, r.Method+" "+strings.TrimPrefix(r.URL.Path, "/api/v3"))
		if !strings.HasSuffix(r.URL.Path, "/repos/acme/widget/rulesets/25145510") {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			if f.failGet {
				http.Error(w, `{"message":"boom"}`, 500)
				return
			}
		case http.MethodPut:
			f.puts++
			b, _ := io.ReadAll(r.Body)
			f.bodies = append(f.bodies, b)
			if f.puts == f.failPut {
				http.Error(w, `{"message":"boom"}`, 500)
				return
			}
			var in map[string]any
			if err := json.Unmarshal(b, &in); err != nil {
				http.Error(w, `{"message":"not JSON"}`, 400)
				return
			}
			for _, k := range []string{"name", "target", "enforcement", "conditions", "rules", "bypass_actors"} {
				if v, ok := in[k]; ok {
					f.obj[k] = v
				}
			}
			if f.drop {
				f.obj["bypass_actors"] = []any{}
			}
			f.obj["updated_at"] = "2026-10-16T12:00:00.000Z"
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(f.obj)
	}))
	t.Cleanup(srv.Close)
	dir := t.TempDir()
	pemFile := filepath.Join(dir, "ca.pem")
	if err := os.WriteFile(pemFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw}), 0o600); err != nil {
		t.Fatal(err)
	}
	f.env = []string{"GH_HOST=" + strings.TrimPrefix(srv.URL, "https://"), "GH_ENTERPRISE_TOKEN=x", "GH_CONFIG_DIR=" + dir, "SSL_CERT_FILE=" + pemFile,
		"GH_PROMPT_DISABLED=1", "NO_COLOR=1", "GH_NO_UPDATE_NOTIFIER=1", "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_AUTHOR_NAME=Op", "GIT_AUTHOR_EMAIL=op@example.invalid", "GIT_COMMITTER_NAME=Op", "GIT_COMMITTER_EMAIL=op@example.invalid"}
	return f
}

func (f *forgeFake) methods(m string) (n int) {
	for _, c := range f.calls {
		if strings.HasPrefix(c, m+" ") {
			n++
		}
	}
	return n
}

// ghRun runs script with sh in dir, with the environment of the fake forge and extra, and gives its status and standard error.
func (f *forgeFake) run(t *testing.T, dir, script string, extra ...string) (int, string) {
	t.Helper()
	file := filepath.Join(t.TempDir(), "script.sh")
	if err := os.WriteFile(file, []byte(script), 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("sh", file)
	cmd.Dir = dir
	cmd.Env = append(append(os.Environ(), f.env...), extra...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), stderr.String()
	} else if err != nil {
		t.Fatal(err)
	}
	return 0, stderr.String()
}

// migrationHost makes the host directory of a target set up by M2d: a bare origin whose main is the head of
// layup-setup (the pushed setup), the work area with its clone, out/CODEOWNERS and the applies directory. It
// gives the Phase2 of the migration, the bare repository and the work area.
func migrationHost(t *testing.T) (p Phase2, bare, area string) {
	t.Helper()
	host, bare := t.TempDir(), filepath.Join(t.TempDir(), "origin.git")
	area = filepath.Join(host, "work", "acme", "widget")
	git := func(dir string, args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=Op", "GIT_AUTHOR_EMAIL=op@example.invalid", "GIT_COMMITTER_NAME=Op", "GIT_COMMITTER_EMAIL=op@example.invalid")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	target := filepath.Join(area, "target")
	for _, d := range []string{filepath.Join(area, "out"), filepath.Dir(bare)} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	git(".", "init", "-q", "--bare", "-b", "main", bare)
	git(".", "init", "-q", "-b", "main", target)
	git(target, "remote", "add", "origin", bare)
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("root\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(target, "add", "-A")
	git(target, "commit", "-q", "-m", "root")
	git(target, "switch", "-q", "-c", "layup-setup")
	if err := os.WriteFile(filepath.Join(target, "docs.txt"), []byte("setup\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	git(target, "add", "-A")
	git(target, "commit", "-q", "-m", "chore: setup S13")
	git(target, "push", "-q", "origin", "layup-setup:main")
	if err := os.WriteFile(filepath.Join(area, "out", "CODEOWNERS"), []byte(codeOwnersTwoKinds), 0o644); err != nil {
		t.Fatal(err)
	}
	p = phase2Owner()
	p.Host, p.Migrate = host, &Migration{RulesetID: rulesetID}
	return p, bare, area
}

// updateScript gives commands.sh of S13 for the migration.
func updateScript(p Phase2) string {
	return commandsText(StepsPhase2(p, Calls{}), s13Record())
}

// reduced gives the object as the update writes it: the fields of a write only.
func reduced(m map[string]any) map[string]any {
	out := map[string]any{}
	for _, k := range []string{"name", "target", "enforcement", "conditions", "rules", "bypass_actors"} {
		out[k] = m[k]
	}
	return out
}

func decodeFile(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	return m
}

// withUnknown adds to the live object what the Operator may have added by hand and the file never held: a rule of
// a type LAYUP does not write, a bypass actor, and a parameter of pull_request that GitHub or the Operator set.
func withUnknown(m map[string]any) map[string]any {
	m["rules"] = append(m["rules"].([]any), map[string]any{"type": "required_signatures"})
	m["bypass_actors"] = []any{map[string]any{"actor_id": float64(5), "actor_type": "RepositoryRole", "bypass_mode": "pull_request"}}
	for _, r := range m["rules"].([]any) {
		if r := r.(map[string]any); r["type"] == "pull_request" {
			r["parameters"].(map[string]any)["unknown_parameter"] = float64(7)
		}
	}
	return m
}

// The update writes the live object back with one parameter changed and nothing else: the rule that LAYUP never
// wrote, the bypass actor and the parameter that the file never held stay; it never POSTs; it keeps the object
// before, the body, the answer of the PUT and the object after in DIR/applies; the read back equals the body; and a
// second run does nothing.
func TestTheUpdateKeepsWhatIsUnknown(t *testing.T) {
	p, _, _ := migrationHost(t)
	live := withUnknown(liveFor(t))
	f := newForge(t, live)
	want := reduced(liveFor(t))
	want["bypass_actors"] = live["bypass_actors"]
	want["rules"] = withUnknown(liveFor(t))["rules"]
	for _, r := range want["rules"].([]any) {
		if r := r.(map[string]any); r["type"] == "pull_request" {
			r["parameters"].(map[string]any)["require_code_owner_review"] = true
		}
	}
	if code, stderr := f.run(t, p.Host, updateScript(p)); code != 0 {
		t.Fatalf("the update exited %d:\n%s", code, stderr)
	}
	a := filepath.Join(p.Host, "applies", "acme", "widget")
	if f.methods("POST") != 0 || f.methods("PUT") != 1 {
		t.Errorf("the calls %v; want one PUT and no POST", f.calls)
	}
	var sent map[string]any
	if err := json.Unmarshal(f.bodies[0], &sent); err != nil || !jsonEqual(sent, want) {
		t.Errorf("the PUT body %s (err %v); want the live object with only require_code_owner_review changed:\n%v", f.bodies[0], err, want)
	}
	for _, name := range []string{"default.before.json", "default.body.json", "default.update.json", "default.after.json"} {
		if _, err := os.Stat(filepath.Join(a, name)); err != nil {
			t.Errorf("the output %s is not kept: %v", name, err)
		}
	}
	body, _ := os.ReadFile(filepath.Join(a, "default.body.json"))
	after, _ := os.ReadFile(filepath.Join(a, "default.after.json"))
	if !bytes.Equal(body, after) || !bytes.Equal(body, append(bytes.TrimRight(f.bodies[0], "\n"), '\n')) {
		t.Errorf("the body, the PUT body and the object read back differ:\n%s\n%s\n%s", body, f.bodies[0], after)
	}
	before := decodeFile(t, filepath.Join(a, "default.before.json"))
	if before["id"] != float64(rulesetID) || before["source"] != "acme/widget" {
		t.Errorf("the object before is %v; want the live object, with its id and source", before)
	}
	n := len(f.calls)
	if code, stderr := f.run(t, p.Host, updateScript(p)); code != 0 || len(f.calls) != n {
		t.Errorf("a second run: exit %d, %d new calls, %s; want exit 0 and no call", code, len(f.calls)-n, stderr)
	}
}

func jsonEqual(a, b any) bool {
	x, _ := json.Marshal(a)
	y, _ := json.Marshal(b)
	return bytes.Equal(x, y)
}

// The update refuses, with no write, a live ruleset that is not the one the records name, and stops where a call
// fails or the read back differs; a run again does what is not done.
func TestTheUpdateRefusals(t *testing.T) {
	for _, c := range []struct {
		name   string
		change func(m map[string]any)
		want   string
	}{
		{"another id", func(m map[string]any) { m["id"] = float64(7) }, "id 7"},
		{"another name", func(m map[string]any) { m["name"] = "other" }, "name other"},
		{"another source", func(m map[string]any) { m["source"] = "acme/other" }, "source acme/other"},
		{"a ruleset in evaluate mode", func(m map[string]any) { m["enforcement"] = "evaluate" }, "enforcement evaluate"},
		{"no pull_request rule", func(m map[string]any) { m["rules"] = []any{map[string]any{"type": "deletion"}} }, "0 pull_request rules"},
		{"two pull_request rules", func(m map[string]any) {
			m["rules"] = append(m["rules"].([]any), map[string]any{"type": "pull_request", "parameters": map[string]any{}})
		}, "2 pull_request rules"},
	} {
		p, _, _ := migrationHost(t)
		live := liveFor(t)
		c.change(live)
		f := newForge(t, live)
		code, stderr := f.run(t, p.Host, updateScript(p))
		if code == 0 || !strings.Contains(stderr, c.want) || f.methods("PUT") != 0 || fileExists(filepath.Join(p.Host, "applies", "acme", "widget", "default.body.json")) {
			t.Errorf("%s: exit %d, %d PUT, stderr %q; want a refusal that names %q and no write", c.name, code, f.methods("PUT"), stderr, c.want)
		}
	}

	// a PUT that fails: no answer file, no read back, exit 1; the run again does the PUT
	p, _, _ := migrationHost(t)
	f := newForge(t, liveFor(t))
	f.failPut = 1
	a := filepath.Join(p.Host, "applies", "acme", "widget")
	if code, _ := f.run(t, p.Host, updateScript(p)); code == 0 || fileExists(filepath.Join(a, "default.update.json")) || fileExists(filepath.Join(a, "default.after.json")) {
		t.Errorf("a PUT that fails: exit %d, or an output was kept", code)
	}
	if code, stderr := f.run(t, p.Host, updateScript(p)); code != 0 || !fileExists(filepath.Join(a, "default.after.json")) {
		t.Errorf("the run again: exit %d, %s; want it done", code, stderr)
	}

	// a forge that dropped a field: the read back differs from the body, exit 1, no after file
	p, _, _ = migrationHost(t)
	f = newForge(t, withUnknown(liveFor(t)))
	f.drop = true
	if code, stderr := f.run(t, p.Host, updateScript(p)); code == 0 || !strings.Contains(stderr, "differs") || fileExists(filepath.Join(p.Host, "applies", "acme", "widget", "default.after.json")) {
		t.Errorf("a read back that differs: exit %d, %q; want exit 1, the word differs, and no after file", code, stderr)
	}

	// a head that is not the head of layup-setup: no call to the forge at all
	p, bare, _ := migrationHost(t)
	gitIn(t, bare, "update-ref", "refs/heads/main", gitIn(t, bare, "rev-parse", "refs/heads/main~1"))
	f = newForge(t, liveFor(t))
	if code, stderr := f.run(t, p.Host, updateScript(p)); code == 0 || len(f.calls) != 0 || !strings.Contains(stderr, "not the head of layup-setup") {
		t.Errorf("a head that moved: exit %d, %d calls, %q; want a stop before any call", code, len(f.calls), stderr)
	}
}

func fileExists(p string) bool { _, err := os.Stat(p); return err == nil }

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// enforcements gives the enforcement that each PUT of the fake forge wrote, in order.
func (f *forgeFake) enforcements() (out []string) {
	for _, b := range f.bodies {
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		out = append(out, m["enforcement"].(string))
	}
	return out
}

// The landing script (docs/spec/scaffold.md): as produced for this task it refuses at its first command and runs
// no program; with its gate removed it makes one commit that adds .github/CODEOWNERS on the head, disables the ruleset,
// pushes, and puts back the AFTER object (the updated protection, never the one before), whatever the push does.
func TestTheLandingScript(t *testing.T) {
	// not applicable: the first command refuses and no gh or git runs
	q, _, _ := migrationHost(t)
	g := newForge(t, liveFor(t))
	if code, stderr := g.run(t, q.Host, landScript(q, false), "PATH="+t.TempDir()); code == 0 || !strings.Contains(stderr, "not applicable") || len(g.calls) != 0 {
		t.Errorf("the script as produced: exit %d, %q; want a refusal that says not applicable, and no call", code, stderr)
	}

	prepare := func(t *testing.T) (Phase2, string, *forgeFake) {
		p, bare, _ := migrationHost(t)
		f := newForge(t, withUnknown(liveFor(t)))
		if code, stderr := f.run(t, p.Host, updateScript(p)); code != 0 {
			t.Fatalf("the update before the landing exited %d: %s", code, stderr)
		}
		f.calls, f.bodies, f.puts = nil, nil, 0
		return p, bare, f
	}

	// the landing
	p, bare, f := prepare(t)
	before := gitIn(t, bare, "rev-parse", "refs/heads/main")
	if code, stderr := f.run(t, p.Host, landScript(p, true)); code != 0 {
		t.Fatalf("the landing exited %d:\n%s", code, stderr)
	}
	if got := gitIn(t, bare, "log", "--format=%H", "-n", "2", "refs/heads/main"); !strings.HasSuffix(got, before) || strings.Count(got, "\n") != 1 {
		t.Errorf("main: %q; want one new commit on %s", got, before)
	}
	if names := gitIn(t, bare, "diff", "--name-status", before, "refs/heads/main"); names != "A\t.github/CODEOWNERS" {
		t.Errorf("the commit changes %q; want only the addition of .github/CODEOWNERS", names)
	}
	if got := gitIn(t, bare, "show", "refs/heads/main:.github/CODEOWNERS"); got+"\n" != codeOwnersTwoKinds {
		t.Errorf("CODEOWNERS on main:\n%s", got)
	}
	if got := f.enforcements(); strings.Join(got, ",") != "disabled,active" {
		t.Errorf("the PUTs wrote %v; want disabled, then active", got)
	}
	after := decodeFile(t, filepath.Join(p.Host, "applies", "acme", "widget", "default.after.json"))
	if !jsonEqual(reduced(f.obj), after) || reduced(f.obj)["enforcement"] != "active" {
		t.Errorf("the ruleset left behind is not the AFTER object, active")
	}

	// the push fails: the AFTER object is put back, the exit is not 0
	p, bare, f = prepare(t)
	hook := filepath.Join(bare, "hooks", "pre-receive")
	if err := os.WriteFile(hook, []byte("#!/bin/sh\nexit 1\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	code, _ := f.run(t, p.Host, landScript(p, true))
	after = decodeFile(t, filepath.Join(p.Host, "applies", "acme", "widget", "default.after.json"))
	if code == 0 || strings.Join(f.enforcements(), ",") != "disabled,active" || !jsonEqual(reduced(f.obj), after) {
		t.Errorf("a push that fails: exit %d, PUTs %v; want exit 1, disabled then the AFTER object back", code, f.enforcements())
	}

	// the restore fails: loud, and it says enforcement may be off
	p, _, f = prepare(t)
	f.failPut = 2
	if code, stderr := f.run(t, p.Host, landScript(p, true)); code == 0 || !strings.Contains(stderr, "enforcement may be off") || !strings.Contains(stderr, "acme/widget") || !strings.Contains(stderr, "25145510") {
		t.Errorf("a restore that fails: exit %d, %q; want exit 1 and a message that names the repository, the ID and that enforcement may be off", code, stderr)
	}

	// the disable fails: no push, no restore
	p, bare, f = prepare(t)
	f.failPut = 1
	before = gitIn(t, bare, "rev-parse", "refs/heads/main")
	if code, _ := f.run(t, p.Host, landScript(p, true)); code == 0 || f.puts != 1 || gitIn(t, bare, "rev-parse", "refs/heads/main") != before {
		t.Errorf("a disable that fails: exit %d, %d PUTs; want exit 1, one PUT, no push", code, f.puts)
	}

	// the refusals before any write: no after file, an after file without the rule, a head that moved, a live ruleset that is not the AFTER object
	for _, c := range []struct {
		name   string
		change func(t *testing.T, p Phase2, bare string, f *forgeFake)
		want   string
	}{
		{"no after file", func(t *testing.T, p Phase2, _ string, _ *forgeFake) {
			_ = os.Remove(filepath.Join(p.Host, "applies", "acme", "widget", "default.after.json"))
		}, "default.after.json is missing"},
		{"an after file without the rule", func(t *testing.T, p Phase2, _ string, _ *forgeFake) {
			file := filepath.Join(p.Host, "applies", "acme", "widget", "default.after.json")
			data, _ := os.ReadFile(file)
			_ = os.WriteFile(file, bytes.ReplaceAll(data, []byte(`"require_code_owner_review":true`), []byte(`"require_code_owner_review":false`)), 0o644)
		}, "does not require the code owner"},
		{"a head that moved", func(t *testing.T, p Phase2, bare string, _ *forgeFake) {
			gitIn(t, bare, "update-ref", "refs/heads/main", gitIn(t, bare, "rev-parse", "refs/heads/main~1"))
		}, "not the head of layup-setup"},
		{"a live ruleset that changed since the update", func(t *testing.T, _ Phase2, _ string, f *forgeFake) {
			f.obj["rules"] = []any{map[string]any{"type": "deletion"}}
		}, "is not the object that the update left"},
	} {
		p, bare, f := prepare(t)
		c.change(t, p, bare, f)
		if code, stderr := f.run(t, p.Host, landScript(p, true)); code == 0 || !strings.Contains(stderr, c.want) || f.methods("PUT") != 0 {
			t.Errorf("%s: exit %d, %d PUTs, %q; want a refusal that says %q and no write", c.name, code, f.methods("PUT"), stderr, c.want)
		}
	}
}

// s13Tree runs the real S13 of phase 2 over a real work area of the stand-in (the other steps are done stubs, none
// of them commits), in the mode p gives, and gives the files that its commit changed on layup-setup, as
// "STATUS\tPATH" lines, and the tree of the area.
func s13Tree(t *testing.T, p Phase2) (changed, target string) {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	m := doneSteps()
	m["S05"] = Step{ID: "S05", Run: func(Input) Outcome { return Outcome{Kind: Done, Evidence: "ok"} }}
	m["S13"] = StepsPhase2(p, Calls{})["S13"]
	target = filepath.Join(w.Dir, work.TargetPath)
	before := gitOut(t, target, "rev-parse", "layup-setup")
	if res, err := Run(w.Dir, m, testWho, noStep); err != nil || slices.Contains(res.Results(), Fail) {
		t.Fatalf("the run: %v, %+v", err, res.Steps)
	}
	return gitOut(t, target, "diff", "--name-status", before, "layup-setup"), target
}

// A work area built again must reproduce the commits that were recorded (scaffold.md, The work area): for a target
// set up before the rule S13 writes into the tree the files of M2d alone, so the rebuilt head is the recorded one;
// a new target's commit adds .github/CODEOWNERS, with the owner of the gate paths.
func TestARebuiltAreaReproducesTheRecordedCommits(t *testing.T) {
	p := phase2Owner()
	p.Name, p.Host = standin.Name, t.TempDir()
	changed, target := s13Tree(t, p)
	if want := "A\t.github/CODEOWNERS\nA\tdocs/setup/branch-protection.json"; changed != want {
		t.Errorf("a new target: S13 changed %q; want %q", changed, want)
	}
	if got := gitOut(t, target, "cat-file", "-p", "layup-setup:.github/CODEOWNERS"); !strings.Contains(got, "/docs/gates/ @octo-operator") {
		t.Errorf("CODEOWNERS in the tree of a new target: %q", got)
	}
	p.Migrate = &Migration{RulesetID: rulesetID}
	changed, target = s13Tree(t, p)
	if want := "A\tdocs/setup/branch-protection.json"; changed != want {
		t.Errorf("a target set up before the rule: S13 changed %q; want %q, the file of M2d alone", changed, want)
	}
	if _, err := os.Stat(filepath.Join(target, ".github", "CODEOWNERS")); err == nil {
		t.Errorf("a rebuilt area of a target set up before the rule has .github/CODEOWNERS in its work tree")
	}
}
