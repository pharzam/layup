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

// commitFiles writes each file (a path and its text; "" removes the path) into
// the repository at repo, and commits them; it gives the commit.
func commitFiles(t *testing.T, repo string, files map[string]string, message string) string {
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
