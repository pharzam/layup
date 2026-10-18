//go:build integration

package rules

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
)

// The basic layup/rules over the real git (docs/spec/loop.md, layup/rules;
// row 63a of the plan, #352): each case builds a repository in a temporary
// directory, with a base commit that holds the true register and a head that
// changes it, and judges the two by Basic over Git(dir).

// itWorld is a repository with a home of its own and a fixed environment, so
// that no configuration of the host reaches the git processes of the test
// (guardrails §2, a tool under test reads the host's home configuration).
type itWorld struct {
	t    *testing.T
	dir  string
	env  []string
	keep bool // the index holds an entry that add -A would undo
}

func newWorld(t *testing.T) *itWorld {
	t.Helper()
	home := t.TempDir()
	w := &itWorld{t: t, dir: t.TempDir(), env: []string{"PATH=" + os.Getenv("PATH"), "HOME=" + home, "LC_ALL=C", "GIT_CONFIG_NOSYSTEM=1",
		"GIT_CONFIG_GLOBAL=/dev/null", "GIT_ATTR_NOSYSTEM=1", "GIT_TERMINAL_PROMPT=0", "GIT_AUTHOR_NAME=LAYUP test",
		"GIT_AUTHOR_EMAIL=test@layup.invalid", "GIT_COMMITTER_NAME=LAYUP test", "GIT_COMMITTER_EMAIL=test@layup.invalid"}}
	w.git("init", "-q", "-b", "main")
	return w
}

func (w *itWorld) git(args ...string) string {
	w.t.Helper()
	full := append([]string{"-c", "maintenance.auto=false", "-c", "core.autocrlf=false", "-c", "commit.gpgsign=false",
		"-c", "core.hooksPath=/dev/null", "-c", "core.attributesFile=/dev/null"}, args...)
	cmd := exec.Command("git", full...)
	cmd.Dir, cmd.Env = w.dir, w.env
	out, err := cmd.CombinedOutput()
	if err != nil {
		w.t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

func (w *itWorld) write(path, data string) {
	w.t.Helper()
	full := filepath.Join(w.dir, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		w.t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(data), 0o644); err != nil {
		w.t.Fatal(err)
	}
}

func (w *itWorld) remove(path string) {
	w.t.Helper()
	if err := os.RemoveAll(filepath.Join(w.dir, path)); err != nil {
		w.t.Fatal(err)
	}
}

// commit stages everything and commits it, and gives the commit's ID; after
// cache it commits the index as it is, which add -A would undo.
func (w *itWorld) commit(message string) string {
	w.t.Helper()
	if !w.keep {
		w.git("add", "-A")
	}
	w.keep = false
	w.git("commit", "-q", "--allow-empty", "-m", message)
	return w.git("rev-parse", "HEAD")
}

// branch puts the work tree at from, on a new branch.
func (w *itWorld) branch(name, from string) {
	w.t.Helper()
	w.git("checkout", "-q", "-f", "-b", name, from)
}

// cache puts an entry of the given mode into the index with the bytes data.
func (w *itWorld) cache(mode, path, data string) {
	w.t.Helper()
	cmd := exec.Command("git", "hash-object", "-w", "--stdin")
	cmd.Dir, cmd.Env, cmd.Stdin = w.dir, w.env, strings.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		w.t.Fatal(err)
	}
	w.cacheinfo(mode, strings.TrimSpace(string(out)), path)
}

func (w *itWorld) cacheinfo(mode, object, path string) {
	w.t.Helper()
	w.git("update-index", "--add", "--cacheinfo", mode+","+object+","+path)
	w.keep = true
}

// baseOf builds the base commit: the true register at its root (the stand-in
// of the records commit), a rule file, a raw fact, the guardrails and a
// product file.
func baseOf(t *testing.T) (*itWorld, string) {
	t.Helper()
	w := newWorld(t)
	w.write("rule-paths.tsv", prdRegister)
	w.write("AGENTS.md", "rules\n")
	w.write("docs/facts/F-0001.md", "raw\n")
	w.write(GuardrailsPath, head)
	w.write("README.md", "readme\n")
	return w, w.commit("base")
}

func (w *itWorld) judge(in Input) Result {
	w.t.Helper()
	return Basic(Git(w.dir), in)
}

func (w *itWorld) in(base, head string, records ...Record) Input {
	return Input{Base: base, Head: head, RegisterRev: base, PreRev: base, Records: records}
}

func TestTheRegisterIsReadAtTheBaseInARealRepo(t *testing.T) {
	w, b := baseOf(t)
	w.branch("forged", b)
	w.write("AGENTS.md", "forged\n")
	forged := header + "\nnothing/\t—\tbaseline\n"
	w.write("rule-paths.tsv", forged)
	w.write("docs/setup/rule-paths.tsv", forged)
	h := w.commit("a head that forges the register")
	want(t, "the true register", w.judge(w.in(b, h)), Failure, ReasonRulePath, "AGENTS.md")
	// The control: the register of the head, read by hand, lets the same
	// change through, so the case above can fail.
	forgedBytes, err := git.Show(w.dir, h, "rule-paths.tsv")
	if err != nil {
		t.Fatal(err)
	}
	if r := Judge(read(t, string(forgedBytes)), []Touched{file("AGENTS.md", "forged\n")}, nil, nil); r.State != Success {
		t.Errorf("the control, a forged register of the head: %+v, want a success", r)
	}
	in := w.in(b, h)
	in.RegisterRev = h
	want(t, "the head as the register's commit", w.judge(in), Error, ReasonRegister, "")
}

func TestARegisterThatCannotBeReadIsAnErrorInARealRepo(t *testing.T) {
	w, b := baseOf(t)
	variants := map[string]func(){
		"no file":       func() { w.remove("rule-paths.tsv") },
		"a header only": func() { w.write("rule-paths.tsv", header+"\n") },
		"a directory":   func() { w.remove("rule-paths.tsv"); w.write("rule-paths.tsv/x", "x\n") },
		"a bad exception": func() {
			w.write("rule-paths.tsv", strings.Replace(prdRegister, "AGENTS.md\t—", "AGENTS.md\tadded lines in section 2", 1))
		},
		"a link":            func() { w.remove("rule-paths.tsv"); w.cache("120000", "rule-paths.tsv", "AGENTS.md") },
		"not a register":    func() { w.write("rule-paths.tsv", "hello\n") },
		"a carriage return": func() { w.write("rule-paths.tsv", strings.ReplaceAll(prdRegister, "\n", "\r\n")) },
	}
	w.branch("head", b)
	w.write("AGENTS.md", "changed\n")
	h := w.commit("head")
	for name, edit := range variants {
		w.branch("v-"+strings.ReplaceAll(name, " ", "-"), b)
		edit()
		r := w.commit(name)
		in := w.in(b, h)
		in.RegisterRev = r
		want(t, name, w.judge(in), Error, ReasonRegister, "")
	}
	in := w.in(b, h)
	in.RegisterRev = "no-such-branch"
	want(t, "an unknown revision", w.judge(in), Error, ReasonRegister, "")
	want(t, "the control, the true register", w.judge(w.in(b, h)), Failure, ReasonRulePath, "AGENTS.md")
}

func TestARenderedFactsRecordInARealRepo(t *testing.T) {
	w, b := baseOf(t)
	w.branch("h", b)
	w.write(factsPath, factsWant)
	h := w.commit("the rendered record")
	want(t, "equal bytes", w.judge(w.in(b, h, factsRecord)), Success, ReasonRenderedRecord, "")
	want(t, "the record not named", w.judge(w.in(b, h)), Failure, ReasonRulePath, factsPath)
	w.branch("h2", b)
	w.write(factsPath, strings.TrimSuffix(factsWant, "\n"))
	h2 := w.commit("one byte short")
	want(t, "one byte short", w.judge(w.in(b, h2, factsRecord)), Failure, ReasonNotEqual, factsPath)
}

// The bytes that decide are those of the blob, whatever the attributes would
// make of the checkout: a text file whose checkout has carriage returns is
// stored with line feeds, and a file marked binary keeps its carriage returns.
func TestTheBlobIsComparedNotTheCheckout(t *testing.T) {
	crlf := strings.ReplaceAll(factsWant, "\n", "\r\n")
	for name, c := range map[string]struct {
		attr  string
		state State
		why   string
	}{
		"a text file with eol=crlf": {factsPath + " text eol=crlf\n", Success, ReasonRenderedRecord},
		"a file marked binary":      {factsPath + " -text\n", Failure, ReasonNotEqual},
	} {
		w := newWorld(t)
		w.write("rule-paths.tsv", prdRegister)
		w.write(".gitattributes", c.attr)
		b := w.commit("base")
		w.branch("h", b)
		w.write(factsPath, crlf)
		h := w.commit("the record with carriage returns")
		path := ""
		if c.state == Failure {
			path = factsPath
		}
		want(t, name, w.judge(w.in(b, h, factsRecord)), c.state, c.why, path)
	}
}

func TestARenderedPRDFromItsPreImage(t *testing.T) {
	w, b := baseOf(t)
	w.branch("spec", b)
	w.write(prdPath, prdPre)
	s := w.commit("the specification")
	w.write(prdPath, string(prdRendered(t)))
	r := w.commit("the rendered commit")
	in := w.in(b, r, prdRecord)
	in.PreRev = s
	want(t, "the first bet, from the specification head", w.judge(in), Success, ReasonRenderedRecord, "")

	w.branch("more", s)
	w.write(prdPath, string(prdRendered(t))+"one more line\n")
	r2 := w.commit("a rendered commit with one more line")
	in = w.in(b, r2, prdRecord)
	in.PreRev = s
	want(t, "one more line", w.judge(in), Failure, ReasonNotEqual, prdPath)

	w.branch("side", b)
	w.write("README.md", "side\n")
	side := w.commit("a side branch")
	for name, pre := range map[string]string{"the base, with no PRD": b, "the head": r, "a side branch": side} {
		in = w.in(b, r, prdRecord)
		in.PreRev = pre
		want(t, name, w.judge(in), Error, ReasonInput, "")
	}

	// A later bet: the base holds the PRD and is the pre-image.
	w2 := newWorld(t)
	w2.write("rule-paths.tsv", prdRegister)
	w2.write(prdPath, prdPre)
	b2 := w2.commit("base with a PRD")
	w2.branch("later", b2)
	w2.write(prdPath, string(prdRendered(t)))
	r3 := w2.commit("the rendering")
	want(t, "a later bet, from the base", w2.judge(w2.in(b2, r3, prdRecord)), Success, ReasonRenderedRecord, "")
}

func TestALinkASubmoduleAndAnExecutableAreNoRecords(t *testing.T) {
	for name, c := range map[string]struct {
		mode, data string
		state      State
		reason     string
	}{
		"a regular file (the control)":         {"100644", factsWant, Success, ReasonRenderedRecord},
		"a link whose target is the rendering": {"120000", factsWant, Failure, ReasonNotARegularFile},
		"an executable with the bytes":         {"100755", factsWant, Failure, ReasonNotARegularFile},
		"a submodule":                          {"160000", "", Failure, ReasonNotARegularFile},
	} {
		w, b := baseOf(t)
		w.branch("h", b)
		if c.mode == "160000" {
			w.cacheinfo("160000", "1111111111111111111111111111111111111111", factsPath)
		} else {
			w.cache(c.mode, factsPath, c.data)
		}
		h := w.commit(name)
		path := ""
		if c.state == Failure {
			path = factsPath
		}
		want(t, name, w.judge(w.in(b, h, factsRecord)), c.state, c.reason, path)
	}
}

func TestADeletionARenameAndADirectoryAtTheRecordPath(t *testing.T) {
	w := newWorld(t)
	w.write("rule-paths.tsv", prdRegister)
	w.write(factsPath, factsWant)
	b := w.commit("base with the record")

	w.branch("deleted", b)
	w.remove(factsPath)
	h := w.commit("deleted")
	want(t, "the record deleted", w.judge(w.in(b, h, factsRecord)), Failure, ReasonNotEqual, factsPath)
	want(t, "the record deleted and not named", w.judge(w.in(b, h)), Failure, ReasonRulePath, factsPath)

	w.branch("renamed", b)
	w.git("mv", factsPath, "docs/facts/other.md")
	h = w.commit("renamed inside docs/facts")
	other := factsRecord
	other.Path = "docs/facts/other.md"
	want(t, "renamed inside docs/facts, the new name named", w.judge(w.in(b, h, other)), Failure, ReasonRulePath, factsPath)

	w.branch("moved", b)
	w.git("mv", factsPath, "docs/moved.md")
	h = w.commit("moved out of docs/facts")
	want(t, "moved out of the register's directory", w.judge(w.in(b, h, factsRecord)), Failure, ReasonNotEqual, factsPath)

	w.branch("dir", b)
	w.remove(factsPath)
	w.write(factsPath+"/x", factsWant)
	h = w.commit("a directory at the path")
	want(t, "a directory at the record's path", w.judge(w.in(b, h, factsRecord)), Failure, ReasonNotEqual, factsPath)
}

func TestTheNamesOfARealDiffAreExact(t *testing.T) {
	for _, name := range []string{"docs/tests/a b\"c é.md", "docs/tests/line\nfeed.md", "docs/tests/tab\there.md", "docs/tests/\\back.md"} {
		w, b := baseOf(t)
		w.branch("h", b)
		w.write(name, "x\n")
		h := w.commit("a rule path with an awkward name")
		want(t, name, w.judge(w.in(b, h)), Failure, ReasonRulePath, name)
	}
}

func TestTheBaseMustBeAnAncestorOfTheHeadInARealRepo(t *testing.T) {
	w, b := baseOf(t)
	w.git("checkout", "-q", "--orphan", "unrelated")
	w.git("rm", "-rq", "--cached", ".")
	w.write("AGENTS.md", "other history\n")
	h := w.commit("an unrelated head")
	want(t, "a head that does not hold the base", w.judge(w.in(b, h)), Error, ReasonBase, "")
	// The base as the head changes nothing; the register lives at a commit of
	// its own, as it does on the records branch.
	w.branch("records", b)
	r := w.commit("a records commit")
	in := w.in(b, b)
	in.RegisterRev = r
	want(t, "the base as the head", w.judge(in), Success, ReasonNoRulePath, "")
}

func TestTheGuardrailsExceptionInARealRepo(t *testing.T) {
	edit := func(name, text string, extra ...func(*itWorld)) (*itWorld, string, string) {
		w, b := baseOf(t)
		w.branch("h", b)
		w.write(GuardrailsPath, text)
		for _, f := range extra {
			f(w)
		}
		return w, b, w.commit(name)
	}
	added := strings.Replace(head, "- c\n", "- c\n- d\n", 1)
	w, b, h := edit("a line added in §2", added)
	want(t, "a line added in §2", w.judge(w.in(b, h)), Success, ReasonNoRulePath, "")
	w, b, h = edit("two lines added in §2 and a rendered record", strings.Replace(head, "- c\n", "- c\n- d\n- e\n", 1), func(w *itWorld) { w.write(factsPath, factsWant) })
	want(t, "with a rendered record", w.judge(w.in(b, h, factsRecord)), Success, ReasonRenderedRecord, "")
	w, b, h = edit("a line added in §3", strings.Replace(head, "three\n", "three\nfour\n", 1))
	want(t, "a line added in §3", w.judge(w.in(b, h)), Failure, ReasonRulePath, GuardrailsPath)
	w, b, h = edit("a line removed", strings.Replace(head, "- a\n", "", 1))
	want(t, "a line removed", w.judge(w.in(b, h)), Failure, ReasonRulePath, GuardrailsPath)
	w, b, h = edit("a rule path beside it", added, func(w *itWorld) { w.write("AGENTS.md", "changed\n") })
	want(t, "a rule path beside the exception", w.judge(w.in(b, h)), Failure, ReasonRulePath, "AGENTS.md")
	w, b, h = edit("the file deleted", head, func(w *itWorld) { w.remove(GuardrailsPath) })
	want(t, "the guardrails deleted", w.judge(w.in(b, h)), Failure, ReasonRulePath, GuardrailsPath)
	w, b, _ = edit("a link", added)
	w.branch("link", b)
	w.remove(GuardrailsPath)
	w.cache("120000", GuardrailsPath, added)
	h = w.commit("the guardrails as a link")
	want(t, "the guardrails as a link", w.judge(w.in(b, h)), Failure, ReasonRulePath, GuardrailsPath)
}

// The item rendering-equal reads a real tree (task T-yeav, #404): the mode of
// each entry, and a pre-image that is a real earlier commit.
func TestRenderingEqualReadsARealTreeWithLinksAndModes(t *testing.T) {
	for name, c := range map[string]struct {
		mode, data string
		want       bool
	}{
		"a regular file (the control)": {"100644", factsWant, true},
		"a link to the rendering":      {"120000", factsWant, false},
		"an executable":                {"100755", factsWant, false},
		"a submodule":                  {"160000", "", false},
	} {
		w, b := baseOf(t)
		w.branch("h", b)
		if c.mode == "160000" {
			w.cacheinfo("160000", "1111111111111111111111111111111111111111", factsPath)
		} else {
			w.cache(c.mode, factsPath, c.data)
		}
		h := w.commit(name)
		if ok, err := RenderingEqual(Git(w.dir), h, "", []string{factsPath}, []Record{factsRecord}); ok != c.want || err != nil {
			t.Errorf("%s: %v, %v; want %v", name, ok, err, c.want)
		}
	}
	w, b := baseOf(t)
	w.write(prdPath, prdPre)
	pre := w.commit("the PRD before the bet")
	w.write(prdPath, string(prdRendered(t)))
	h := w.commit("the PRD rendered")
	call := func(head, preRev string) (bool, error) {
		return RenderingEqual(Git(w.dir), head, preRev, []string{prdPath}, []Record{prdRecord})
	}
	if ok, err := call(h, pre); !ok || err != nil {
		t.Errorf("a PRD of its earlier revision: %v, %v; want true", ok, err)
	}
	for name, c := range map[string][2]string{"the head as its own pre-image": {h, h}, "the base, which holds no PRD": {h, b}} {
		if ok, err := call(c[0], c[1]); ok || err != nil {
			t.Errorf("%s: %v, %v; want false and no error", name, ok, err)
		}
	}
	for name, c := range map[string][2]string{"a PreRev that is no revision": {h, "nope"}, "a head that is no revision": {"nope", pre}} {
		if ok, err := call(c[0], c[1]); ok || err == nil {
			t.Errorf("%s: %v, %v; want false and an error", name, ok, err)
		}
	}
}
