package rules

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/spec"
)

// RenderingEqual and RecordRendered (task T-yeav, #404): the handoff item
// rendering-equal and the per-record rule that it shares with Judge.

// renderScene is scene() whose head holds the rendered facts record and the
// rendered PRD, and whose pre-image commit p holds the PRD before the bet.
func renderScene(t *testing.T) *fake {
	t.Helper()
	f := scene()
	f.trees["h"][prdPath] = reg(string(prdRendered(t)))
	f.ids["alias"] = "h"
	return f
}

// An rcall is one call of RenderingEqual: the input of the good case.
type rcall struct {
	head, pre string
	paths     []string
	recs      []Record
}

func goodCall() rcall {
	return rcall{"head", "pre", []string{factsPath, prdPath}, []Record{factsRecord, prdRecord}}
}

func (c rcall) run(r Repo) (bool, error) { return RenderingEqual(r, c.head, c.pre, c.paths, c.recs) }

// counted counts the calls of git that a Repo gets.
type counted struct {
	Repo
	n int
}

func (c *counted) RevParse(rev string) (string, error) { c.n++; return c.Repo.RevParse(rev) }
func (c *counted) Show(rev, p string) ([]byte, error)  { c.n++; return c.Repo.Show(rev, p) }
func (c *counted) LsTree(rev, p string) ([]git.TreeEntry, error) {
	c.n++
	return c.Repo.LsTree(rev, p)
}
func (c *counted) IsAncestor(a, b string) (bool, error) { c.n++; return c.Repo.IsAncestor(a, b) }

func TestRecordRenderedIsTheOneRuleOfJudgeAndOfTheItem(t *testing.T) {
	good, pre := string(prdRendered(t)), []byte(prdPre)
	at := func(mode, typ, data string) Touched {
		return Touched{Path: factsPath, Present: true, Mode: mode, Type: typ, Data: []byte(data)}
	}
	noText := Record{Path: factsPath, Kind: Facts, Version: factsVersion, Texts: map[string]string{"1": "One."}}
	prdAt := Touched{Path: prdPath, Present: true, Mode: "100644", Type: "blob", Data: []byte(good)}
	for _, c := range []struct {
		name   string
		rec    Record
		t      Touched
		reason string
	}{
		{"an equal facts record", factsRecord, at("100644", "blob", factsWant), ""},
		{"a PRD of its pre-image", prdRecord, prdAt, ""},
		{"a byte off", factsRecord, at("100644", "blob", factsWant+"!"), ReasonNotEqual},
		{"absent", factsRecord, Touched{Path: factsPath}, ReasonNotEqual},
		{"a link", factsRecord, at("120000", "blob", factsWant), ReasonNotARegularFile},
		{"an executable", factsRecord, at("100755", "blob", factsWant), ReasonNotARegularFile},
		{"a gitlink", factsRecord, at("160000", "commit", ""), ReasonNotARegularFile},
		{"a rendering with a finding", noText, at("100644", "blob", factsWant), ReasonNoRendering},
	} {
		if got := RecordRendered(c.rec, c.t, pre); got != c.reason {
			t.Errorf("%s: RecordRendered = %q, want %q", c.name, got, c.reason)
		}
		want := Result{Success, ReasonRenderedRecord, ""}
		if c.reason != "" {
			want = Result{Failure, c.reason, c.rec.Path}
		}
		if got := Judge(read(t, prdRegister), []Touched{c.t}, []Record{c.rec}, map[string][]byte{c.rec.Path: pre}); got != want {
			t.Errorf("%s: Judge = %+v, want %+v", c.name, got, want)
		}
	}
}

func TestRenderingEqualIsTrueOnlyForTheRenderingOfARegularFile(t *testing.T) {
	if ok, err := goodCall().run(renderScene(t)); !ok || err != nil {
		t.Fatalf("the control, a facts record and a PRD of its pre-image: %v, %v", ok, err)
	}
	only := goodCall()
	only.paths, only.recs = []string{factsPath}, []Record{factsRecord}
	if ok, err := only.run(renderScene(t)); !ok || err != nil {
		t.Errorf("a facts record alone (it needs no PreRev): %v, %v", ok, err)
	}
	f := renderScene(t)
	f.trees["h"][factsPath] = reg(factsWant + "!")
	if ok, err := goodCall().run(f); ok || err != nil {
		t.Errorf("a byte off: %v, %v; want false and no error", ok, err)
	}
}

func TestRenderingEqualRefusesEachBadInputAndItsControlPasses(t *testing.T) {
	badPath := func(p string) func(*fake, *rcall) {
		return func(_ *fake, c *rcall) { c.paths, c.recs = []string{p}, []Record{{Path: p, Kind: Facts}} }
	}
	atHead := func(p string, v fakeFile) func(*fake, *rcall) {
		return func(f *fake, _ *rcall) { f.trees["h"][p] = v }
	}
	rows := []struct {
		name   string
		static bool // refused before any call of git
		mut    func(*fake, *rcall)
	}{
		{"no artifact", true, func(_ *fake, c *rcall) { c.paths = nil }},
		{"no artifact and no record", true, func(_ *fake, c *rcall) { c.paths, c.recs = nil, nil }},
		{"an artifact repeated", true, func(_ *fake, c *rcall) { c.paths = []string{factsPath, factsPath} }},
		{"an artifact with no record", true, func(_ *fake, c *rcall) { c.paths = []string{factsPath, "docs/x.md"} }},
		{"more artifacts than records", true, func(_ *fake, c *rcall) { c.paths = append(c.paths, "docs/x.md") }},
		{"a record with no artifact", true, func(_ *fake, c *rcall) { c.paths = c.paths[:1] }},
		{"a record repeated", true, func(_ *fake, c *rcall) { c.recs = []Record{factsRecord, factsRecord} }},
		{"a path that leaves the tree", true, badPath("../x")},
		{"an absolute path", true, badPath("/abs")},
		{"a path into .git", true, badPath(".git/x")},
		{"a path with a backslash", true, badPath(`a\b`)},
		{"a path that is a pathspec", true, badPath(":x")},
		{"the guardrails file", true, badPath(GuardrailsPath)},
		{"a kind that is none", true, func(_ *fake, c *rcall) { c.recs[0].Kind = "other" }},
		{"a PRD with no PreRev", false, func(_ *fake, c *rcall) { c.pre = "" }},
		{"a PreRev that is the head", false, func(_ *fake, c *rcall) { c.pre = "head" }},
		{"a PreRev that is another name of the head", false, func(_ *fake, c *rcall) { c.pre = "alias" }},
		{"a PreRev that is no ancestor", false, func(f *fake, _ *rcall) { f.anc["p>h"] = false }},
		{"a PRD absent at the PreRev", false, func(f *fake, _ *rcall) { delete(f.trees["p"], prdPath) }},
		{"a PRD at the PreRev that is a link", false, func(f *fake, _ *rcall) { f.trees["p"][prdPath] = fakeFile{"120000", "blob", prdPre} }},
		{"a file absent at the head", false, func(f *fake, _ *rcall) { delete(f.trees["h"], factsPath) }},
		{"a link at the head", false, atHead(factsPath, fakeFile{"120000", "blob", factsWant})},
		{"an executable at the head", false, atHead(factsPath, fakeFile{"100755", "blob", factsWant})},
		{"a gitlink at the head", false, atHead(factsPath, fakeFile{"160000", "commit", ""})},
		{"a byte off in the facts", false, atHead(factsPath, reg(factsWant+"!"))},
		{"a rendering with a finding", false, func(_ *fake, c *rcall) {
			c.recs[0].Texts = map[string]string{"1": "One."}
		}},
	}
	for _, r := range rows {
		ctl := &counted{Repo: renderScene(t)}
		if ok, err := goodCall().run(ctl); !ok || err != nil {
			t.Fatalf("%s: the control is %v, %v", r.name, ok, err)
		}
		f, c := renderScene(t), goodCall()
		c.paths, c.recs = slices.Clone(c.paths), slices.Clone(c.recs)
		r.mut(f, &c)
		cnt := &counted{Repo: f}
		if ok, err := c.run(cnt); ok || err != nil {
			t.Errorf("%s: %v, %v; want false and no error", r.name, ok, err)
		}
		if r.static && cnt.n != 0 {
			t.Errorf("%s: %d calls of git before the refusal; want none", r.name, cnt.n)
		}
	}
}

func TestRenderingEqualTakesThePreImageFromTheNamedRevisionNeverTheHead(t *testing.T) {
	good := string(prdRendered(t))
	// A byte outside the fields that RenderPRD writes, so the forged head is
	// the rendering of itself: only the earlier revision can show it.
	forged := strings.Replace(good, "does one thing", "does two things", 1)
	if again, _ := spec.RenderPRD([]byte(forged), prdAnswer); string(again) != forged || forged == good {
		t.Fatal("the premise: the forged head must render to itself and differ from the rendering")
	}
	run := func(head, pre string) (bool, error, []string) {
		f := renderScene(t)
		f.trees["h"][prdPath] = reg(head)
		ok, err := RenderingEqual(f, "head", pre, []string{prdPath}, []Record{prdRecord})
		return ok, err, f.shown
	}
	for _, pre := range []string{"pre", "head", "alias"} {
		if ok, err, _ := run(forged, pre); ok || err != nil {
			t.Errorf("the forged byte, PreRev %q: %v, %v; want false", pre, ok, err)
		}
	}
	ok, err, shown := run(good, "pre")
	if !ok || err != nil || !slices.Contains(shown, "p:"+prdPath) {
		t.Errorf("the restored byte: %v, %v, read %q; want true, reading the pre-image at p", ok, err, shown)
	}
}

func TestRenderingEqualGivesAnErrorForEachFailedCallOfGit(t *testing.T) {
	boom := errors.New("boom")
	for _, c := range []struct {
		name string
		mut  func(*fake, *rcall)
		is   error
	}{
		{"the head does not resolve", func(_ *fake, c *rcall) { c.head = "nope" }, nil},
		{"the PreRev does not resolve", func(_ *fake, c *rcall) { c.pre = "nope" }, nil},
		{"the rev-parse fails", func(f *fake, _ *rcall) { f.errs["RevParse"] = boom }, boom},
		{"the ancestry fails", func(f *fake, _ *rcall) { f.errs["IsAncestor"] = boom }, boom},
		{"the tree fails", func(f *fake, _ *rcall) { f.errs["LsTree"] = boom }, boom},
		{"the show fails", func(f *fake, _ *rcall) { f.errs["Show"] = boom }, boom},
		{"the tree of the head fails, with no PRD", func(f *fake, c *rcall) {
			c.paths, c.recs = []string{factsPath}, []Record{factsRecord}
			f.errs["LsTree"] = boom
		}, boom},
		{"the show of the head fails, with no PRD", func(f *fake, c *rcall) {
			c.paths, c.recs = []string{factsPath}, []Record{factsRecord}
			f.errs["Show"] = boom
		}, boom},
	} {
		f, call := renderScene(t), goodCall()
		c.mut(f, &call)
		ok, err := call.run(f)
		if ok || err == nil || (c.is != nil && !errors.Is(err, c.is)) {
			t.Errorf("%s: %v, %v; want false and the error", c.name, ok, err)
		}
	}
}
