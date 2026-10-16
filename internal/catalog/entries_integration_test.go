//go:build integration

package catalog

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/gate"
	"github.com/pharzam/layup/internal/git"
)

// who is the author and the committer of each commit of these tests, at a
// fixed time.
var who = git.Identity{Name: "fixture", Email: "fixture@layup.invalid", Time: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

// goEnv makes the commands of the kinds run on the toolchain of the host, with
// its build cache, with no download and with no setting of the host (L-A1;
// note 10 of the plan review of #91).
func goEnv(t *testing.T) {
	t.Helper()
	cache, err := exec.Command("go", "env", "GOCACHE").Output()
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1",
		"GOCACHE": strings.TrimSpace(string(cache)), "GOENV": "off", "GOFLAGS": "", "GOTOOLCHAIN": "local", "GOPROXY": "off"} {
		t.Setenv(k, v)
	}
}

// writeFiles writes each file (a path and its text; "" removes the path) into
// the repository at repo.
func writeFiles(t *testing.T, repo string, files map[string]string) {
	t.Helper()
	for p, text := range files {
		path := filepath.Join(repo, filepath.FromSlash(p))
		if text == "" {
			if err := os.Remove(path); err != nil {
				t.Fatal(err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

// commitFiles writes the files into the repository at repo, and commits them;
// it gives the commit.
func commitFiles(t *testing.T, repo string, files map[string]string, message string) string {
	t.Helper()
	writeFiles(t, repo, files)
	return commitAll(t, repo, message)
}

// commitAll commits the work tree of the repository at repo.
func commitAll(t *testing.T, repo, message string) string {
	t.Helper()
	if err := git.Add(repo); err != nil {
		t.Fatal(err)
	}
	if err := git.Commit(repo, message, who); err != nil {
		t.Fatal(err)
	}
	head, err := git.RevParse(repo, "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	return head
}

// render makes a new repository with the files and the manifest that the
// setup writes from the entry: the clean commit of a target.
func render(t *testing.T, e *Entry) (repo, clean string) {
	t.Helper()
	repo = filepath.Join(t.TempDir(), "target")
	if err := git.Init(repo); err != nil {
		t.Fatal(err)
	}
	files, err := e.Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	m, err := e.Manifest()
	if err != nil {
		t.Fatal(err)
	}
	all := map[string]string{"docs/gates.tsv": string(m)}
	for _, f := range files {
		all[f.Path] = string(f.Data)
	}
	return repo, commitFiles(t, repo, all, "the setup")
}

// withFixture commits the fixture of kind on the clean commit, on a branch of
// its own, and gives the commit. A fixture that does not apply fails the test.
func withFixture(t *testing.T, e *Entry, repo, clean, kind string) string {
	t.Helper()
	patch, err := e.Fixture(kind)
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), kind+".patch")
	if err := os.WriteFile(p, patch, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := git.SwitchCreate(repo, "fixture-"+kind, clean); err != nil {
		t.Fatal(err)
	}
	if err := git.Apply(repo, p); err != nil {
		t.Fatalf("the fixture of %s does not apply to the clean commit: %v", kind, err)
	}
	return commitFiles(t, repo, nil, "the fixture of "+kind)
}

func nop(int, int, string) func() { return func() {} }

// row gives the row of kind in the table.
func row(t *testing.T, tab gate.Table, kind string) gate.Row {
	t.Helper()
	for _, r := range tab.Rows {
		if r.Kind == kind {
			return r
		}
	}
	t.Fatalf("the table has no row for %s: %+v", kind, tab.Rows)
	return gate.Row{}
}

// Each entry of the binary, on a target rendered from it (D8 and D10 of #91;
// LAYUP's CI runs each fixture, §6): on the clean commit, which has no file in
// a scope of the Go entry, each active kind passes or is clear and each
// pending kind is clear; with its fixture applied, each active kind fails. The
// manifest reads back by the reader of layup gate.
func TestTheFixturesOfEachEntry(t *testing.T) {
	goEnv(t)
	stacks, err := Stacks(Embedded())
	if err != nil || len(stacks) == 0 {
		t.Fatalf("the binary embeds no entry: %q, %v", stacks, err)
	}
	for _, s := range stacks {
		e, err := Read(Embedded(), s)
		if err != nil {
			t.Fatal(err)
		}
		m, err := e.Manifest()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := gate.ReadManifest(m); err != nil {
			t.Errorf("%s: the manifest: %v", s, err)
		}
		repo, clean := render(t, e)
		tab, err := gate.Run(repo, clean, clean, nop, io.Discard)
		if err != nil {
			t.Fatalf("%s: the clean commit: %v", s, err)
		}
		for _, k := range e.Kinds() {
			r := row(t, tab, k.Name)
			if k.State == "pending" && (r.Result != "clear" || r.Reason != "pending: no product path") ||
				k.State == "active" && r.Result != "pass" && r.Result != "clear" {
				t.Errorf("%s: the clean commit: %s %s %q; want pass or clear", s, k.Name, r.Result, r.Reason)
			}
		}
		for _, k := range e.Kinds() {
			if k.State != "active" {
				continue
			}
			head := withFixture(t, e, repo, clean, k.Name)
			var out bytes.Buffer
			tab, err := gate.Run(repo, clean, head, nop, &out)
			if err != nil {
				t.Fatalf("%s: the fixture of %s: %v", s, k.Name, err)
			}
			if r := row(t, tab, k.Name); r.Result != "fail" {
				t.Errorf("%s: the fixture of %s: %s %q; want fail\n%s", s, k.Name, r.Result, r.Reason, out.String())
			}
		}
	}
}

// The active kinds of the Go entry pass on good code: the clean commit with a
// package that gofmt keeps, that go vet accepts, and whose test passes. So a
// command that fails on each tree cannot hide behind the clean commit, where
// the kinds are clear (note 11 of the plan review of #91).
func TestTheActiveKindsOfTheGoEntryPassOnGoodCode(t *testing.T) {
	goEnv(t)
	e, err := Read(Embedded(), "go")
	if err != nil {
		t.Fatal(err)
	}
	repo, clean := render(t, e)
	head := commitFiles(t, repo, map[string]string{
		"good/good.go":      "// Package good is good code for the gate kinds of the Go entry.\npackage good\n\n// Sum gives the sum of a and b.\nfunc Sum(a, b int) int { return a + b }\n",
		"good/good_test.go": "package good\n\nimport \"testing\"\n\nfunc TestSum(t *testing.T) {\n\tif Sum(1, 2) != 3 {\n\t\tt.Fatal(\"Sum(1, 2) != 3\")\n\t}\n}\n",
	}, "good code")
	var out bytes.Buffer
	tab, err := gate.Run(repo, clean, head, nop, &out)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range e.Kinds() {
		want := "pass"
		if k.State == "pending" {
			want = "fail" // a pending kind refuses each product path (gate.md)
		}
		if r := row(t, tab, k.Name); r.Result != want {
			t.Errorf("good code: %s %s %q; want %s\n%s", k.Name, r.Result, r.Reason, want, out.String())
		}
	}
}

// The job script of each entry of the binary (D5 of #91) and layup gate, on the
// same base and head, give the same pass or fail, with the same result and
// reason for each line of the table of the run; each input that layup gate
// refuses (exit 2) fails the job (note 4 of the plan review, finding 1 of round
// 1), and so does a kind with no row. Each sh of the host runs the script, with
// the same output.
func TestTheJobScriptOfEachEntry(t *testing.T) {
	goEnv(t)
	const header = "kind\tstate\ttool\tcommand\tscope\tconfig\n"
	const manifest = header +
		"ok\tactive\tsh\ttrue\t./*.txt\t—\n" +
		"bad\tactive\tsh\texit 3\t./*.txt\t—\n" +
		"notool\tactive\tno-such-tool-of-91\ttrue\t./*.txt\t—\n" +
		"none\tactive\tsh\ttrue\t./*.go\t—\n" +
		"pend\tpending\t—\t—\tsrc/*.txt docs/rules.md\t—\n"
	base := map[string]string{"docs/gates.tsv": manifest, "docs/rules.md": "the rules\n", "a.md": "a\n"}
	withRow := func(row string) map[string]string {
		return map[string]string{"docs/gates.tsv": strings.Replace(manifest, "ok\tactive\tsh\ttrue\t./*.txt\t—\n", row, 1)}
	}
	cases := []struct {
		name, kind   string
		base, head   map[string]string
		result       string // the result of the job; "" when layup gate refuses the manifest
		reason       string
		gateRefuses  bool
		scriptRefers string // a part of the reason of the job when layup gate refuses the manifest
		tree         bool   // the base is the tree of the base commit, not the commit
	}{
		{"active, the command exits 0", "ok", base, map[string]string{"x.txt": "x\n"}, "pass", "—", false, "", false},
		{"active, the command exits 3", "bad", base, map[string]string{"x.txt": "x\n"}, "fail", "exit 3", false, "", false},
		{"active, its tool is not found", "notool", base, map[string]string{"x.txt": "x\n"}, "not-active", "tool not found: no-such-tool-of-91", false, "", false},
		{"active, no product path", "none", base, map[string]string{"x.txt": "x\n"}, "clear", "no product path", false, "", false},
		{"pending, no product path changed", "pend", base, map[string]string{"x.txt": "x\n", "docs/rules.md.bak": "b\n", "src/.txt": "e\n"}, "clear", "pending: no product path", false, "", false},
		{"pending, a file under the directory of a pattern", "pend", base, map[string]string{"src/a/b.txt": "b\n"}, "fail", "pending: product path changed: src/a/b.txt", false, "", false},
		{"pending, the path of a pattern changed", "pend", base, map[string]string{"docs/rules.md": "new rules\n"}, "fail", "pending: product path changed: docs/rules.md", false, "", false},
		{"pending, the path of a pattern removed", "pend", base, map[string]string{"docs/rules.md": ""}, "fail", "pending: product path changed: docs/rules.md", false, "", false},
		{"no manifest", "ok", map[string]string{"a.md": "a\n"}, map[string]string{"x.txt": "x\n"}, "", "", true, "no manifest", false},
		{"a scope pattern of another form", "ok", map[string]string{"docs/gates.tsv": strings.Replace(manifest, "ok\tactive\tsh\ttrue\t./*.txt", "ok\tactive\tsh\ttrue\t*.txt", 1)},
			map[string]string{"x.txt": "x\n"}, "", "", true, "the scope pattern *.txt", false},
		{"a scope pattern of another form in another row", "ok", map[string]string{"docs/gates.tsv": strings.Replace(manifest, "src/*.txt", "src/*/x.txt", 1)},
			map[string]string{"x.txt": "x\n"}, "", "", true, "the scope pattern src/*/x.txt", false},
		{"a header of another form", "ok", map[string]string{"docs/gates.tsv": strings.Replace(manifest, "kind\tstate", "kind\tstatus", 1)},
			map[string]string{"x.txt": "x\n"}, "", "", true, "line 1 is not the header row", false},
		{"a kind with no row", "nosuch", base, map[string]string{"x.txt": "x\n"}, "", "", false, "no row for the kind nosuch", false},
		// Each form that the reader of the manifest refuses (round 1 of #91,
		// finding 1), a builtin as the tool (note 2), and a base that is no
		// commit (note 5).
		{"an empty field", "ok", withRow("ok\tactive\tsh\t\t./*.txt\t—\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "an empty field", false},
		{"an empty field of config", "ok", withRow("ok\tactive\tsh\ttrue\t./*.txt\t\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "an empty field", false},
		{"a config value that is not a path", "ok", withRow("ok\tactive\tsh\ttrue\t./*.txt\t../x\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "the config path ../x", false},
		{"a scope of two spaces", "ok", withRow("ok\tactive\tsh\ttrue\t./*.txt  ./*.md\t—\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "the scope", false},
		{"a carriage return in a row", "ok", withRow("ok\tactive\tsh\ttrue\t./*.txt\t—\r\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "a carriage return", false},
		{"an empty line", "ok", withRow("ok\tactive\tsh\ttrue\t./*.txt\t—\n\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "an empty line", false},
		{"no line feed at the end", "ok", map[string]string{"docs/gates.tsv": strings.TrimSuffix(manifest, "\n")}, map[string]string{"x.txt": "x\n"}, "", "", true, "no line feed", false},
		{"a byte that is not UTF-8", "ok", withRow("ok\tactive\tsh\ttrue \xff\t./*.txt\t—\n"), map[string]string{"x.txt": "x\n"}, "", "", true, "not UTF-8", false},
		{"a kind twice", "ok", map[string]string{"docs/gates.tsv": manifest + "ok\tactive\tsh\ttrue\t./*.txt\t—\n"}, map[string]string{"x.txt": "x\n"}, "", "", true, "the kind ok is there twice", false},
		{"a builtin as the tool", "ok", withRow("ok\tactive\t:\ttrue\t./*.txt\t—\n"), map[string]string{"x.txt": "x\n"}, "not-active", "tool not found: :", false, "", false},
		{"a base that is a tree", "ok", base, map[string]string{"x.txt": "x\n"}, "", "", true, "not a commit", true},
		// The cases of a run of round 2 that came after fifteen minutes and
		// is skipped (Bootstrap mode rule 4): a program that is also a
		// builtin of dash, a last byte NUL with no line feed, a header row
		// with no row.
		{"a program that is a builtin of dash", "ok", withRow("ok\tactive\ttrue\ttrue\t./*.txt\t—\n"), map[string]string{"x.txt": "x\n"}, "pass", "—", false, "", false},
		{"a last byte NUL", "ok", map[string]string{"docs/gates.tsv": manifest + "\x00"}, map[string]string{"x.txt": "x\n"}, "", "", true, "no line feed", false},
		{"a header row and no row", "ok", map[string]string{"docs/gates.tsv": header}, map[string]string{"x.txt": "x\n"}, "", "", true, "no row", false},
	}
	stacks, err := Stacks(Embedded())
	if err != nil || len(stacks) == 0 {
		t.Fatalf("the binary embeds no entry: %q, %v", stacks, err)
	}
	for _, s := range stacks {
		e, err := Read(Embedded(), s)
		if err != nil {
			t.Fatal(err)
		}
		files, err := e.Files("example.com/target")
		if err != nil {
			t.Fatal(err)
		}
		script := filepath.Join(t.TempDir(), "gates.sh")
		for _, f := range files {
			if f.Path == ".github/gates.sh" {
				if err := os.WriteFile(script, f.Data, 0o644); err != nil {
					t.Fatal(err)
				}
			}
		}
		// Each shell of the host that is a POSIX sh: CI runs the job with the
		// sh of ubuntu-latest, dash (round 1 of #91, note 9).
		var shells [][]string
		for _, sh := range [][]string{{"sh"}, {"dash"}, {"bash", "--posix"}} {
			if _, err := exec.LookPath(sh[0]); err == nil {
				shells = append(shells, sh)
			}
		}
		for _, sh := range shells {
			if out, err := exec.Command(sh[0], append(sh[1:], "-n", script)...).CombinedOutput(); err != nil {
				t.Fatalf("%s: %s -n .github/gates.sh: %v\n%s", s, sh, err, out)
			}
		}
		for _, c := range cases {
			repo := filepath.Join(t.TempDir(), "target")
			if err := git.Init(repo); err != nil {
				t.Fatal(err)
			}
			b := commitFiles(t, repo, c.base, "the base")
			h := commitFiles(t, repo, c.head, "the head")
			if c.tree {
				var err error
				if b, err = git.RevParse(repo, b+"^{tree}"); err != nil {
					t.Fatal(err)
				}
			}
			var out []byte
			var passed bool
			for i, sh := range shells {
				cmd := exec.Command(sh[0], append(sh[1:], script, c.kind)...)
				cmd.Dir, cmd.Env = repo, append(os.Environ(), "GATE_BASE="+b, "GATE_HEAD="+h)
				o, err := cmd.Output()
				var exit *exec.ExitError
				if err != nil && !errors.As(err, &exit) {
					t.Fatalf("%s: %s: the script with %s: %v", s, c.name, sh, err)
				}
				if i > 0 && (string(o) != string(out) || passed != (err == nil)) {
					t.Errorf("%s: %s: the job with %s gives %q (passed %v); with %s %q (passed %v)", s, c.name, sh, o, err == nil, shells[0], out, passed)
				}
				out, passed = o, err == nil
			}
			got := strings.Split(strings.TrimSuffix(string(out), "\n"), "\t")
			tab, gerr := gate.Run(repo, b, h, nop, io.Discard)
			switch {
			case c.result != "":
				r := row(t, tab, c.kind)
				if gerr != nil || len(got) != 3 || got[1] != r.Result || got[2] != reason(r) || r.Result != c.result || reason(r) != c.reason ||
					passed != (r.Result == "pass" || r.Result == "clear") {
					t.Errorf("%s: %s: the job %q (passed %v); layup gate %s %q, %v; want %s %q", s, c.name, out, passed, r.Result, reason(r), gerr, c.result, c.reason)
				}
			default:
				var input *gate.InputError
				if passed || len(got) != 3 || got[1] != "fail" || !strings.Contains(got[2], c.scriptRefers) || c.gateRefuses != errors.As(gerr, &input) {
					t.Errorf("%s: %s: the job %q (passed %v), layup gate %v; want a failed job that names %q, and layup gate refuses %v", s, c.name, out, passed, gerr, c.scriptRefers, c.gateRefuses)
				}
			}
		}
	}
}

// reason gives the reason of a row as the table writes it: — for an empty
// reason.
func reason(r gate.Row) string {
	if r.Reason == "" {
		return "—"
	}
	return r.Reason
}

// overlayStep gives the text of the step of the job of kind in the workflow that
// takes the gate files from the base (the lines between "- run: |" and "env:",
// without their indent), or fails the test when the job has none.
func overlayStep(t *testing.T, workflow, kind string) string {
	t.Helper()
	_, jobs, _ := strings.Cut(workflow, "\njobs:\n")
	_, job, ok := strings.Cut("\n"+jobs, "\n  "+kind+":\n")
	if !ok {
		t.Fatalf("the workflow has no job of the kind %s", kind)
	}
	_, step, ok := strings.Cut(job, "      - run: |\n")
	if !ok {
		t.Fatalf("the job of the kind %s has no step that takes the gate files from the base", kind)
	}
	body, _, _ := strings.Cut(step, "        env:\n")
	var out []string
	for _, l := range strings.Split(strings.TrimSuffix(body, "\n"), "\n") {
		out = append(out, strings.TrimPrefix(l, "          "))
	}
	return strings.Join(out, "\n") + "\n"
}

// The gate jobs read the base (F-22; loop.md, "The gate jobs read the base"):
// the step of each job runs the base's own script in its mode --overlay, which
// puts the base's manifest, script and config paths into the checkout of the
// head, and then the job runs. A head that weakens its own gate is judged by the
// base's files; the exception is a head ref of the form batch/<point>, with one
// segment after batch/. Each case also runs the job with no step (the control),
// which shows that the case can fail.
func TestTheJobsReadTheBase(t *testing.T) {
	goEnv(t)
	const header = "kind\tstate\ttool\tcommand\tscope\tconfig\n"
	const dirCmd = `test -x docs/gates/dir/run.sh && test "$(cat docs/gates/dir/a.txt)" = ok && test ! -e docs/gates/dir/extra.txt` +
		` && git diff --cached --quiet "$GATE_BASE" -- docs/gates/dir && git diff --quiet -- docs/gates/dir`
	manifestWith := func(cfg string) string {
		return header + "ok\tactive\tsh\ttrue\t./*.txt\t—\n" + "bad\tactive\tsh\texit 3\t./*.txt\t—\n" +
			"cfg\tactive\tsh\ttest \"$(cat docs/gates/cfg.txt)\" = ok\t./*.txt\t" + cfg + "\n" +
			"newcfg\tactive\tsh\ttest -f docs/gates/new.txt\t./*.txt\tdocs/gates/new.txt\n" +
			"dir\tactive\tsh\t" + dirCmd + "\t./*.txt\tdocs/gates/dir\n"
	}
	manifest := manifestWith("docs/gates/cfg.txt")
	weak := strings.Replace(manifest, "bad\tactive\tsh\texit 3", "bad\tactive\tsh\ttrue", 1)
	alwaysPass := "#!/bin/sh\nprintf 'bad\\tpass\\t—\\n'\nexit 0\n"
	base := map[string]string{"docs/gates.tsv": manifest, "docs/gates/cfg.txt": "no\n", "x.txt": "x\n"}
	with := func(over map[string]string) map[string]string {
		m := map[string]string{}
		for k, v := range base {
			m[k] = v
		}
		for k, v := range over {
			m[k] = v
		}
		return m
	}
	outside := t.TempDir()
	// link replaces the path of the repository by a symbolic link to target.
	link := func(p, target string) func(*testing.T, string) {
		return func(t *testing.T, repo string) {
			if err := os.RemoveAll(filepath.Join(repo, p)); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(target, filepath.Join(repo, p)); err != nil {
				t.Fatal(err)
			}
		}
	}
	chmod := func(p string, mode os.FileMode) func(*testing.T, string) {
		return func(t *testing.T, repo string) {
			if err := os.Chmod(filepath.Join(repo, p), mode); err != nil {
				t.Fatal(err)
			}
		}
	}
	dirBase, dirHead := with(map[string]string{"docs/gates/dir/a.txt": "ok\n", "docs/gates/dir/run.sh": "x\n"}),
		map[string]string{"docs/gates/dir/a.txt": "no\n", "docs/gates/dir/run.sh": "y\n", "docs/gates/dir/extra.txt": "e\n"}
	type scenario struct {
		name, kind, ref string
		base, head      map[string]string
		noScript        bool // the base has no .github/gates.sh
		baseIsTree      bool // GATE_BASE names the tree of the base
		baseEdit        func(t *testing.T, repo string)
		headEdit        func(t *testing.T, repo string) // both are made before their commit
		stepFails       string                          // a part of the output of the step when it must fail
		result, reason  string                          // the job after the step: the result and a part of its reason
		control         string                          // the result of the job with no step; "" for none
	}
	cases := []scenario{
		{name: "a head that makes a failing kind pass", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3", control: "pass"},
		{name: "a head whose script always passes", kind: "bad", ref: "feature/x", base: base, head: map[string]string{".github/gates.sh": alwaysPass}, result: "fail", reason: "exit 3", control: "pass"},
		{name: "a head that changes the file of a config path", kind: "cfg", ref: "feature/x", base: base, head: map[string]string{"docs/gates/cfg.txt": "ok\n"}, result: "fail", reason: "exit 1", control: "pass"},
		{name: "a head that adds a file of a config path that the base lacks", kind: "newcfg", ref: "feature/x", base: base, head: map[string]string{"docs/gates/new.txt": "n\n"}, result: "fail", reason: "exit 1", control: "pass"},
		{name: "a head that removes the file of a config path", kind: "cfg", ref: "feature/x", base: with(map[string]string{"docs/gates/cfg.txt": "ok\n"}), head: map[string]string{"docs/gates/cfg.txt": ""}, result: "pass", reason: "—", control: "fail"},
		{name: "a head that has no manifest of the base", kind: "ok", ref: "feature/x", base: with(map[string]string{"docs/gates.tsv": ""}), head: map[string]string{"docs/gates.tsv": manifest}, result: "fail", reason: "no manifest", control: "pass"},
		{name: "a directory config path: its files, their modes, no other", kind: "dir", ref: "feature/x", base: dirBase, head: dirHead, result: "pass", reason: "—", control: "fail",
			baseEdit: chmod("docs/gates/dir/run.sh", 0o755), headEdit: chmod("docs/gates/dir/run.sh", 0o644)},
		{name: "the head of a batch branch", kind: "bad", ref: "batch/activation", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "pass", reason: "—"},
		{name: "a ref that is batch/ alone", kind: "bad", ref: "batch/", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3"},
		{name: "a ref with two segments after batch/", kind: "bad", ref: "batch/a/b", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3"},
		{name: "a ref that is batch alone", kind: "bad", ref: "batch", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3"},
		{name: "a ref that holds batch/ after another word", kind: "bad", ref: "xbatch/a", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3"},
		{name: "a ref that holds batch/ after a segment", kind: "bad", ref: "feature/batch/x", base: base, head: map[string]string{"docs/gates.tsv": weak}, result: "fail", reason: "exit 3"},
		{name: "a base that is a tree", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, baseIsTree: true, stepFails: "not a commit"},
		{name: "a base with no script", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, noScript: true, stepFails: "gates.sh"},
		{name: "a config path in .git", kind: "cfg", ref: "feature/x", base: with(map[string]string{"docs/gates.tsv": manifestWith(".git/hooks/pre-commit")}), head: map[string]string{"docs/gates.tsv": manifest}, stepFails: "in .git"},
		{name: "a config path that is not a path", kind: "cfg", ref: "feature/x", base: with(map[string]string{"docs/gates.tsv": manifestWith("../x")}), head: map[string]string{"docs/gates.tsv": manifest}, stepFails: "is not a path"},
		{name: "a config path that is a link at the base", kind: "cfg", ref: "feature/x", base: with(map[string]string{"docs/gates/cfg.txt": "elsewhere\n"}), head: map[string]string{"docs/gates/cfg.txt": ""}, stepFails: "link",
			baseEdit: link("docs/gates/cfg.txt", "elsewhere")}, // the head has no link: only the base can refuse
		{name: "a config path under a link of the head to the outside", kind: "cfg", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link("docs/gates", outside)},
		{name: "a config path under a link of the head to a directory of the tree", kind: "cfg", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link("docs/gates", "../.github")},
		{name: "a config file that is a link of the head", kind: "cfg", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link("docs/gates/cfg.txt", "../../x.txt")},
		{name: "a manifest that is a link of the head to a file of the tree", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link("docs/gates.tsv", "../x.txt")},
		{name: "a manifest that is a link of the head to nothing", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link("docs/gates.tsv", filepath.Join(outside, "none"))},
		{name: "a script that is a link of the head to the outside", kind: "bad", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "link",
			headEdit: link(".github/gates.sh", filepath.Join(outside, "none"))},
		{name: "a file of the head above a config path", kind: "cfg", ref: "feature/x", base: base, head: map[string]string{"x.txt": "y\n"}, stepFails: "not a directory",
			headEdit: func(t *testing.T, repo string) {
				if err := os.RemoveAll(filepath.Join(repo, "docs", "gates")); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(repo, "docs", "gates"), []byte("f\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			}},
	}
	stacks, err := Stacks(Embedded())
	if err != nil || len(stacks) == 0 {
		t.Fatalf("the binary embeds no entry: %q, %v", stacks, err)
	}
	var shells [][]string
	for _, sh := range [][]string{{"sh"}, {"dash"}, {"bash", "--posix"}} {
		if _, err := exec.LookPath(sh[0]); err == nil {
			shells = append(shells, sh)
		}
	}
	for _, s := range stacks {
		e, err := Read(Embedded(), s)
		if err != nil {
			t.Fatal(err)
		}
		files, err := e.Files("example.com/target")
		if err != nil {
			t.Fatal(err)
		}
		entry := map[string]string{}
		for _, f := range files {
			entry[f.Path] = string(f.Data)
		}
		step := overlayStep(t, entry[".github/workflows/gates.yml"], "static")
		for _, c := range cases {
			all := map[string]string{}
			for k, v := range entry {
				all[k] = v
			}
			for k, v := range c.base {
				if v == "" {
					delete(all, k) // the base has no such file
				} else {
					all[k] = v
				}
			}
			if c.noScript {
				delete(all, ".github/gates.sh")
			}
			repo := filepath.Join(t.TempDir(), "target")
			if err := git.Init(repo); err != nil {
				t.Fatal(err)
			}
			writeFiles(t, repo, all)
			if c.baseEdit != nil {
				c.baseEdit(t, repo)
			}
			b := commitAll(t, repo, "the base")
			if c.baseIsTree {
				if b, err = git.RevParse(repo, b+"^{tree}"); err != nil {
					t.Fatal(err)
				}
			}
			writeFiles(t, repo, c.head)
			if c.headEdit != nil {
				c.headEdit(t, repo)
			}
			h := commitAll(t, repo, "the head")
			var firstStep, firstJob string
			for i, sh := range shells {
				env := append(os.Environ(), "GATE_BASE="+b, "GATE_HEAD="+h, "GATE_HEAD_REF="+c.ref, "RUNNER_TEMP="+t.TempDir())
				run := func(args ...string) (string, bool) {
					cmd := exec.Command(sh[0], append(sh[1:], args...)...)
					cmd.Dir, cmd.Env = repo, env
					o, err := cmd.CombinedOutput()
					var exit *exec.ExitError
					if err != nil && !errors.As(err, &exit) {
						t.Fatalf("%s: %s: %v", s, c.name, err)
					}
					return string(o), err == nil
				}
				reset := func() {
					if out, err := exec.Command("git", "-C", repo, "reset", "--hard", "-q", h).CombinedOutput(); err != nil {
						t.Fatalf("%v: %s", err, out)
					}
				}
				reset()
				if c.control != "" && i == 0 {
					// the control: the same checkout with no step
					o, _ := run(".github/gates.sh", c.kind)
					l := strings.Split(strings.TrimSuffix(o, "\n"), "\n")
					if g := strings.Split(l[len(l)-1], "\t"); len(g) != 3 || g[1] != c.control {
						t.Errorf("%s: %s: the control (no step) gave %q; want %s: the case cannot tell the overlay from none", s, c.name, l[len(l)-1], c.control)
					}
					reset()
				}
				stepOut, stepOK := run("-c", step)
				if c.stepFails != "" {
					if stepOK || !strings.Contains(stepOut, c.stepFails) {
						t.Errorf("%s: %s: with %s the step gave %q (passed %v); want a failure that names %q", s, c.name, sh, stepOut, stepOK, c.stepFails)
					}
					if ents, _ := os.ReadDir(outside); len(ents) != 0 {
						t.Errorf("%s: %s: the step wrote into a directory outside the checkout: %v", s, c.name, ents)
					}
					if st, err := exec.Command("git", "-C", repo, "status", "--porcelain").Output(); err != nil || len(st) != 0 {
						t.Errorf("%s: %s: with %s the failed step changed the checkout: %q (%v)", s, c.name, sh, st, err)
					}
					continue
				}
				if !stepOK {
					t.Errorf("%s: %s: with %s the step failed: %q", s, c.name, sh, stepOut)
					continue
				}
				jobOut, jobOK := run(".github/gates.sh", c.kind)
				lines := strings.Split(strings.TrimSuffix(jobOut, "\n"), "\n")
				last := lines[len(lines)-1]
				got := strings.Split(last, "\t")
				if len(got) != 3 || got[1] != c.result || !strings.Contains(got[2], c.reason) || jobOK != (c.result == "pass" || c.result == "clear") {
					t.Errorf("%s: %s: with %s the job gave %q (passed %v); want %s with a reason that holds %q", s, c.name, sh, last, jobOK, c.result, c.reason)
				}
				if i == 0 {
					firstStep, firstJob = stepOut, jobOut
				} else if stepOut != firstStep || jobOut != firstJob {
					t.Errorf("%s: %s: with %s the step and the job give %q and %q; with %s %q and %q", s, c.name, sh, stepOut, jobOut, shells[0], firstStep, firstJob)
				}
			}
		}
	}
}
