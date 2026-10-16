package rules

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/spec"
)

// The basic layup/rules (docs/spec/loop.md, layup/rules; row 63a of the plan,
// #352). Judge is pure: the paths that differ, what the head holds at each,
// the register and the rendered records in; a state, a reason and a path
// out. Basic reads a Repo, which these tests fake with different bytes at
// the base and at the head, so a read of the wrong one shows.

const (
	prdPath   = "docs/prd/PRD.md"
	factsPath = "docs/facts/numbered.md"
)

// prdRegister is the register of the entries of phase 1 with the PRD added,
// the register of a target that lists it (the PRD is not a rule path under
// the register of setup.md).
var prdRegister = register + prdPath + "\t—\tcatalog\n"

var (
	factsVersion = [][]string{{"1", "need", "in", "7"}, {"2", "constraint", "out", "7"}}
	factsTexts   = map[string]string{"1": "One.", "2": "Two."}
	factsWant    = "1. **need:** One.\n2. **constraint (out of scope):** Two.\n"
	factsRecord  = Record{Path: factsPath, Kind: Facts, Version: factsVersion, Texts: factsTexts}
)

const prdPre = `# PRD

**Non-goals.** No plugin system.

## 6. Functional requirements

| REQ | Statement | MoSCoW | Phase | Facts |
| REQ-001 | The app does one thing. |  |  | F-0001#1 |
| REQ-002 | The app paints itself. |  |  | F-0001#2 |

## 12. Requirements traceability matrix

| REQ | Task | Test |
| REQ-001 | T-xxxx | this test |
`

var prdAnswer = spec.Answer{
	Reqs:  map[string]spec.ReqAnswer{"REQ-001": {Class: "Must", Phase: "1"}, "REQ-002": {Class: "Won't"}},
	Facts: map[string]spec.FactAnswer{"2": {Out: true}},
}

var prdRecord = Record{Path: prdPath, Kind: PRD, Answer: prdAnswer}

func prdRendered(t *testing.T) []byte {
	t.Helper()
	out, f := spec.RenderPRD([]byte(prdPre), prdAnswer)
	if len(f) != 0 {
		t.Fatalf("the fixture does not render: %v", f)
	}
	return out
}

// file is a path of the head: a regular file with data.
func file(path string, data string) Touched {
	return Touched{Path: path, Present: true, Mode: "100644", Type: "blob", Data: []byte(data)}
}

func judge(t *testing.T, regText string, touched []Touched, records []Record, pre map[string][]byte) Result {
	t.Helper()
	return Judge(read(t, regText), touched, records, pre)
}

func want(t *testing.T, name string, got Result, state State, reason, path string) {
	t.Helper()
	if got != (Result{state, reason, path}) {
		t.Errorf("%s: %+v, want %s %q %q", name, got, state, reason, path)
	}
}

func TestJudgeANoRulePathChangeIsASuccess(t *testing.T) {
	in2 := Change{Diff: diff("@@ -5,0 +6 @@\n+- b\n"), Head: []byte(head)}
	for name, touched := range map[string][]Touched{
		"nothing changed":            nil,
		"paths the register lacks":   {file("internal/x/x.go", "x"), file("README.md", "r")},
		"only the guardrails in §2":  {{Path: GuardrailsPath, Present: true, Mode: "100644", Type: "blob", Data: []byte(head), Guardrails: in2}},
		"a path beside the register": {file("docs/cinema.md", "c"), file("docs/ci", "a file named like the directory")},
	} {
		want(t, name, judge(t, register, touched, nil, nil), Success, ReasonNoRulePath, "")
	}
}

func TestJudgeAChangedRulePathFails(t *testing.T) {
	for name, c := range map[string]struct {
		touched []Touched
		path    string
	}{
		"a directory pattern":       {[]Touched{file("docs/tests/test-levels.md", "x")}, "docs/tests/test-levels.md"},
		"an exact pattern":          {[]Touched{file("AGENTS.md", "x")}, "AGENTS.md"},
		"a script anywhere":         {[]Touched{file("tools/gen.sh", "x")}, "tools/gen.sh"},
		"a deleted rule path":       {[]Touched{{Path: "go.mod"}}, "go.mod"},
		"the first of two in order": {[]Touched{file("README.md", "x"), file("go.mod", "x"), file("AGENTS.md", "x")}, "go.mod"},
	} {
		want(t, name, judge(t, register, c.touched, nil, nil), Failure, ReasonRulePath, c.path)
	}
}

// The workflow refusal of Check is folded into the reason rule path
// (condition 3 of the plan review): the register here has no .github/ row, so
// only the workflow rule of Check can refuse the path.
func TestJudgeAWorkflowPathFailsAsARulePath(t *testing.T) {
	noGithub := strings.Replace(register, ".github/\t—\tbaseline\n", "", 1)
	if noGithub == register {
		t.Fatal("the fixture register keeps its .github/ row")
	}
	if reason, _ := Check([]string{".github/workflows/ci.yml"}, read(t, noGithub), Change{}); reason != "workflow" {
		t.Fatalf("Check: %q, want workflow", reason)
	}
	want(t, "a workflow", judge(t, noGithub, []Touched{file(".github/workflows/ci.yml", "on: push\n")}, nil, nil), Failure, ReasonRulePath, ".github/workflows/ci.yml")
}

func TestJudgeAPassingRecordDoesNotHideALaterRulePath(t *testing.T) {
	rec := []Record{factsRecord}
	for name, touched := range map[string][]Touched{
		"the record first": {file(factsPath, factsWant), file("AGENTS.md", "forged")},
		"the record last":  {file("AGENTS.md", "forged"), file(factsPath, factsWant)},
	} {
		want(t, name, judge(t, register, touched, rec, nil), Failure, ReasonRulePath, "AGENTS.md")
	}
}

func TestJudgeAnEqualRenderedFactsRecordPasses(t *testing.T) {
	want(t, "the facts record", judge(t, register, []Touched{file(factsPath, factsWant)}, []Record{factsRecord}, nil), Success, ReasonRenderedRecord, "")
}

func TestJudgeARecordThatDiffersByAByteFails(t *testing.T) {
	for name, got := range map[string]string{
		"the last byte":            strings.TrimSuffix(factsWant, ".\n") + "!\n",
		"no final line feed":       strings.TrimSuffix(factsWant, "\n"),
		"a second final line feed": factsWant + "\n",
		"a carriage return":        strings.Replace(factsWant, "One.\n", "One.\r\n", 1),
		"a byte-order mark":        "\ufeff" + factsWant,
		"a trailing space":         strings.Replace(factsWant, "One.\n", "One. \n", 1),
		"a swapped class":          strings.Replace(factsWant, "**need:**", "**measure:**", 1),
		"an extra line":            factsWant + "3. **need:** Three.\n",
		"a missing line":           strings.SplitAfter(factsWant, "\n")[0],
		"an empty file":            "",
	} {
		want(t, name, judge(t, register, []Touched{file(factsPath, got)}, []Record{factsRecord}, nil), Failure, ReasonNotEqual, factsPath)
	}
}

func TestJudgeAPRDEqualToTheRenderingOfItsPreImagePasses(t *testing.T) {
	good := prdRendered(t)
	for _, cell := range []string{"| REQ-001 | The app does one thing. | Must | 1 | F-0001#1 |", "Out of scope by the bet: 2."} {
		if !strings.Contains(string(good), cell) {
			t.Fatalf("the rendering lacks %q:\n%s", cell, good)
		}
	}
	pre := map[string][]byte{prdPath: []byte(prdPre)}
	rec := []Record{prdRecord}
	want(t, "the rendering", judge(t, prdRegister, []Touched{file(prdPath, string(good))}, rec, pre), Success, ReasonRenderedRecord, "")
	for name, got := range map[string]string{
		"a cell that differs":            strings.Replace(string(good), "| Must | 1 |", "| Should | 1 |", 1),
		"the non-goals that differ":      strings.Replace(string(good), "Out of scope by the bet: 2.", "Out of scope by the bet: 1.", 1),
		"a line outside the rendering":   strings.Replace(string(good), "does one thing", "does two things", 1),
		"the pre-image, not rendered":    prdPre,
		"the rendering with a final \\n": string(good) + "\n",
	} {
		want(t, name, judge(t, prdRegister, []Touched{file(prdPath, got)}, rec, pre), Failure, ReasonNotEqual, prdPath)
	}
}

func TestJudgeARenderingThatFailsIsNoPass(t *testing.T) {
	missingRow := Record{Path: prdPath, Kind: PRD, Answer: spec.Answer{Reqs: map[string]spec.ReqAnswer{"REQ-099": {Class: "Must", Phase: "1"}}}}
	noClass := Record{Path: factsPath, Kind: Facts, Version: [][]string{{"1", "", "in", "7"}}, Texts: factsTexts}
	noText := Record{Path: factsPath, Kind: Facts, Version: factsVersion, Texts: map[string]string{"1": "One."}}
	pre := map[string][]byte{prdPath: []byte(prdPre)}
	for name, c := range map[string]struct {
		rec  Record
		path string
	}{"a requirement with no row": {missingRow, prdPath}, "a fact with no class": {noClass, factsPath}, "a fact with no text": {noText, factsPath}} {
		// The head holds the bytes that a lenient rendering would give.
		got := judge(t, prdRegister, []Touched{file(c.path, "whatever")}, []Record{c.rec}, pre)
		want(t, name, got, Failure, ReasonNoRendering, c.path)
	}
}

func TestJudgeOneRecordDoesNotCoverAnotherPath(t *testing.T) {
	second := "docs/facts/second.md"
	want(t, "the same bytes at a second path", judge(t, register, []Touched{file(factsPath, factsWant), file(second, factsWant)}, []Record{factsRecord}, nil),
		Failure, ReasonRulePath, second)
	other := Record{Path: second, Kind: Facts, Version: [][]string{{"9", "need", "in", "7"}}, Texts: map[string]string{"9": "Nine."}}
	want(t, "the bytes of the other record", judge(t, register, []Touched{file(factsPath, "9. **need:** Nine.\n"), file(second, factsWant)}, []Record{factsRecord, other}, nil),
		Failure, ReasonNotEqual, factsPath)
}

func TestJudgeARecordEntryThatIsNoRegularFileFails(t *testing.T) {
	entry := func(mode, typ string) Touched {
		return Touched{Path: factsPath, Present: true, Mode: mode, Type: typ, Data: []byte(factsWant)}
	}
	for name, c := range map[string]struct {
		touched Touched
		reason  string
	}{
		"a link with the bytes":        {entry("120000", "blob"), ReasonNotARegularFile},
		"a submodule":                  {entry("160000", "commit"), ReasonNotARegularFile},
		"an executable with the bytes": {entry("100755", "blob"), ReasonNotARegularFile},
		"nothing at the path":          {Touched{Path: factsPath}, ReasonNotEqual},
	} {
		want(t, name, judge(t, register, []Touched{c.touched}, []Record{factsRecord}, nil), Failure, c.reason, factsPath)
	}
}

func TestJudgeTheGuardrailsException(t *testing.T) {
	touch := func(c Change) Touched {
		return Touched{Path: GuardrailsPath, Present: true, Mode: "100644", Type: "blob", Data: []byte(head), Guardrails: c}
	}
	in2 := Change{Diff: diff("@@ -5,0 +6 @@\n+- b\n"), Head: []byte(head)}
	in3 := Change{Diff: diff("@@ -8,0 +9 @@\n+three\n"), Head: []byte(head)}
	removed := Change{Diff: diff("@@ -4 +3,0 @@\n-- a\n"), Head: []byte(head)}
	rec := []Record{factsRecord}
	want(t, "additions in §2", judge(t, register, []Touched{touch(in2)}, rec, nil), Success, ReasonNoRulePath, "")
	want(t, "additions in §2 and a rendered record", judge(t, register, []Touched{touch(in2), file(factsPath, factsWant)}, rec, nil), Success, ReasonRenderedRecord, "")
	want(t, "an addition in §3", judge(t, register, []Touched{touch(in3)}, rec, nil), Failure, ReasonRulePath, GuardrailsPath)
	want(t, "a removed line", judge(t, register, []Touched{touch(removed)}, rec, nil), Failure, ReasonRulePath, GuardrailsPath)
	want(t, "no diff given", judge(t, register, []Touched{touch(Change{})}, rec, nil), Failure, ReasonRulePath, GuardrailsPath)
	noException := strings.Replace(register, "docs/guardrails.md\tadded lines in section 2", "docs/guardrails.md\t—", 1)
	want(t, "no exception in the register", judge(t, noException, []Touched{touch(in2)}, rec, nil), Failure, ReasonRulePath, GuardrailsPath)
	want(t, "the guardrails deleted", judge(t, register, []Touched{{Path: GuardrailsPath, Guardrails: Change{Diff: diff("@@ -1,9 +0,0 @@\n-x\n")}}}, rec, nil), Failure, ReasonRulePath, GuardrailsPath)
}

// The limit of the decision: a record whose path the register does not match
// is not judged, whatever the head holds there.
func TestJudgeARecordTheRegisterDoesNotMatchIsNotJudged(t *testing.T) {
	pre := map[string][]byte{prdPath: []byte(prdPre)}
	want(t, "a PRD the register lacks", judge(t, register, []Touched{file(prdPath, "not the rendering at all")}, []Record{prdRecord}, pre), Success, ReasonNoRulePath, "")
	want(t, "a PRD the register lists", judge(t, prdRegister, []Touched{file(prdPath, "not the rendering at all")}, []Record{prdRecord}, pre), Failure, ReasonNotEqual, prdPath)
}

func TestJudgeBadInputIsAnErrorNeverASuccess(t *testing.T) {
	// The head holds the clean, equal rendering: only the record is bad.
	touched := []Touched{file(factsPath, factsWant)}
	bad := func(path string) Record { r := factsRecord; r.Path = path; return r }
	for _, p := range []string{"", "/docs/facts/x.md", "../x", "docs/../x", "docs/facts/../../x", "./docs/facts/x.md", "docs//facts/x.md", "docs/facts/x.md/",
		"docs/facts/.", "docs\\facts\\x.md", "docs/facts/x\x00.md", ".git/config", "docs/.git/config", "docs/.GIT/config", ":x"} {
		want(t, fmt.Sprintf("the path %q", p), judge(t, register, touched, []Record{bad(p)}, nil), Error, ReasonInput, p)
	}
	for name, c := range map[string]struct {
		records []Record
		pre     map[string][]byte
		path    string
	}{
		"two records of one path":      {[]Record{factsRecord, factsRecord}, nil, factsPath},
		"an unknown kind":              {[]Record{{Path: factsPath, Kind: "pdf"}}, nil, factsPath},
		"no kind":                      {[]Record{{Path: factsPath}}, nil, factsPath},
		"a record of the guardrails":   {[]Record{{Path: GuardrailsPath, Kind: Facts}}, nil, GuardrailsPath},
		"a PRD with no pre-image":      {[]Record{prdRecord}, nil, prdPath},
		"a PRD with another pre-image": {[]Record{prdRecord}, map[string][]byte{"docs/other.md": nil}, prdPath},
	} {
		want(t, name, judge(t, register, touched, c.records, c.pre), Error, ReasonInput, c.path)
	}
	// A PRD pre-image that is empty is a pre-image.
	if r := judge(t, register, nil, []Record{prdRecord}, map[string][]byte{prdPath: {}}); r.State != Success {
		t.Errorf("an empty pre-image: %+v, want a success (nothing changed)", r)
	}
}

func TestRenderIsCodesRendering(t *testing.T) {
	gotF, ff := Render(factsRecord, nil)
	wantF, wf := spec.RenderFacts(factsVersion, factsTexts)
	if !bytes.Equal(gotF, wantF) || !reflect.DeepEqual(ff, wf) || string(gotF) != factsWant {
		t.Errorf("the facts record: %q %v, want %q %v", gotF, ff, wantF, wf)
	}
	gotP, pf := Render(prdRecord, []byte(prdPre))
	wantP, wp := spec.RenderPRD([]byte(prdPre), prdAnswer)
	if !bytes.Equal(gotP, wantP) || !reflect.DeepEqual(pf, wp) {
		t.Errorf("the PRD: %q %v, want %q %v", gotP, pf, wantP, wp)
	}
	if out, f := Render(Record{Path: factsPath, Kind: "pdf"}, nil); out != nil || len(f) == 0 {
		t.Errorf("an unknown kind: %q %v, want no bytes and a finding", out, f)
	}
}

func TestEachResultIsOneOfTheClosedSet(t *testing.T) {
	closed := map[State][]string{
		Success: {ReasonNoRulePath, ReasonRenderedRecord},
		Failure: {ReasonRulePath, ReasonNotEqual, ReasonNotARegularFile, ReasonNoRendering},
		Error:   {ReasonRegister, ReasonInput, ReasonBase, ReasonGit},
	}
	pre := map[string][]byte{prdPath: []byte(prdPre)}
	results := []Result{
		judge(t, register, nil, nil, nil),
		judge(t, register, []Touched{file(factsPath, factsWant)}, []Record{factsRecord}, nil),
		judge(t, register, []Touched{file("AGENTS.md", "x")}, nil, nil),
		judge(t, register, []Touched{file(factsPath, "x")}, []Record{factsRecord}, nil),
		judge(t, register, []Touched{{Path: factsPath, Present: true, Mode: "120000", Type: "blob"}}, []Record{factsRecord}, nil),
		judge(t, prdRegister, []Touched{file(prdPath, "x")}, []Record{{Path: prdPath, Kind: PRD, Answer: spec.Answer{Reqs: map[string]spec.ReqAnswer{"REQ-9": {Class: "Must", Phase: "1"}}}}}, pre),
		judge(t, register, nil, []Record{{Path: "../x", Kind: Facts}}, nil),
	}
	for _, r := range results {
		if !slices.Contains(closed[r.State], r.Reason) {
			t.Errorf("%+v is not in the closed set", r)
		}
	}
	if len(results) < 7 {
		t.Fatal("the table lost a case")
	}
	if r := results[len(results)-1]; r.State != Error {
		t.Errorf("a bad record: %+v, want an error", r)
	}
}

// Condition 5 of the plan review: the code reads no copy and checks no hash;
// the copy arrives as parsed values.
func TestBasicReadsNoCopyAndChecksNoHash(t *testing.T) {
	src, err := os.ReadFile("basic.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, word := range []string{"copies", "sha256", "crypto", "hash"} {
		if strings.Contains(strings.ToLower(string(src)), word) {
			t.Errorf("basic.go mentions %q: a copy and its hash are the caller's", word)
		}
	}
}

// fake is a Repo with a tree for each commit. A name maps to a commit ID;
// Show and LsTree answer from the tree of the commit they are given, and the
// calls of Show are kept so that a test can see which commit was read.
type fake struct {
	ids   map[string]string
	trees map[string]map[string]fakeFile
	names []string
	diffs map[string][]byte
	anc   map[string]bool
	errs  map[string]error
	shown []string
}

type fakeFile struct {
	mode, typ string
	data      string
}

func (f *fake) RevParse(rev string) (string, error) {
	if err := f.errs["RevParse"]; err != nil {
		return "", err
	}
	if id, ok := f.ids[rev]; ok {
		return id, nil
	}
	return "", errors.New("unknown revision " + rev)
}

func (f *fake) Show(rev, path string) ([]byte, error) {
	f.shown = append(f.shown, rev+":"+path)
	if err := f.errs["Show"]; err != nil {
		return nil, err
	}
	if v, ok := f.trees[rev][path]; ok {
		return []byte(v.data), nil
	}
	return nil, errors.New("no such path " + rev + ":" + path)
}

func (f *fake) LsTree(rev, path string) ([]git.TreeEntry, error) {
	if err := f.errs["LsTree"]; err != nil {
		return nil, err
	}
	var out []git.TreeEntry
	for p, v := range f.trees[rev] {
		if p == path || strings.HasPrefix(p, path+"/") {
			out = append(out, git.TreeEntry{Mode: v.mode, Type: v.typ, Object: "0", Path: p})
		}
	}
	return out, nil
}

func (f *fake) DiffNames(base, head string) ([]string, error) { return f.names, f.errs["DiffNames"] }
func (f *fake) DiffFile(base, head, path string) ([]byte, error) {
	return f.diffs[path], f.errs["DiffFile"]
}
func (f *fake) IsAncestor(base, head string) (bool, error) {
	if err := f.errs["IsAncestor"]; err != nil {
		return false, err
	}
	return base == head || f.anc[base+">"+head], nil
}

func reg(data string) fakeFile { return fakeFile{"100644", "blob", data} }

// scene is a repository of four commits: the base b, the records commit r
// (which holds the true register), the head h (which holds a forged one that
// lists nothing) and the pre-image commit p, in the chain b, p, h.
func scene() *fake {
	return &fake{
		ids: map[string]string{"base": "b", "head": "h", "records": "r", "pre": "p"},
		trees: map[string]map[string]fakeFile{
			"b": {},
			"r": {"rule-paths.tsv": reg(prdRegister)},
			"h": {"rule-paths.tsv": reg(header + "\nnothing/\t—\tbaseline\n"), "AGENTS.md": reg("forged"), factsPath: reg(factsWant)},
			"p": {prdPath: reg(prdPre)},
		},
		names: []string{"AGENTS.md"},
		anc:   map[string]bool{"b>h": true, "b>p": true, "p>h": true},
		errs:  map[string]error{},
	}
}

func input(records ...Record) Input {
	return Input{Base: "base", Head: "head", RegisterRev: "records", PreRev: "pre", Records: records}
}

func TestBasicReadsTheRegisterAtTheBase(t *testing.T) {
	f := scene()
	want(t, "a forged register at the head", Basic(f, input()), Failure, ReasonRulePath, "AGENTS.md")
	if !slices.Contains(f.shown, "r:rule-paths.tsv") {
		t.Errorf("the register was not read at the records commit: %v", f.shown)
	}
	for _, s := range f.shown {
		if strings.HasPrefix(s, "h:rule-paths.tsv") {
			t.Errorf("the register was read at the head: %v", f.shown)
		}
	}
}

func TestBasicRefusesARegisterItCannotUse(t *testing.T) {
	for name, mutate := range map[string]func(*fake){
		"no file at the records commit": func(f *fake) { f.trees["r"] = map[string]fakeFile{} },
		"a register of the header only": func(f *fake) { f.trees["r"]["rule-paths.tsv"] = reg(header + "\n") },
		"garbage":                       func(f *fake) { f.trees["r"]["rule-paths.tsv"] = reg("not a register\n") },
		"an exception of another path": func(f *fake) {
			f.trees["r"]["rule-paths.tsv"] = reg(strings.Replace(register, "AGENTS.md\t—", "AGENTS.md\tadded lines in section 2", 1))
		},
		"a link in place of the file":    func(f *fake) { f.trees["r"]["rule-paths.tsv"] = fakeFile{"120000", "blob", register} },
		"a directory in place of it":     func(f *fake) { f.trees["r"] = map[string]fakeFile{"rule-paths.tsv/x": reg("x")} },
		"the records commit is the head": func(f *fake) { f.ids["records"] = "h"; f.trees["h"]["rule-paths.tsv"] = reg(register) },
		"an unknown records commit":      func(f *fake) { delete(f.ids, "records") },
	} {
		f := scene()
		mutate(f)
		want(t, name, Basic(f, input()), Error, ReasonRegister, "")
	}
}

func TestBasicJudgesWhatTheHeadHolds(t *testing.T) {
	f := scene()
	f.names = []string{factsPath}
	want(t, "the equal record", Basic(f, input(factsRecord)), Success, ReasonRenderedRecord, "")
	f.trees["h"][factsPath] = reg(factsWant + "x")
	want(t, "one more byte", Basic(f, input(factsRecord)), Failure, ReasonNotEqual, factsPath)
	f.trees["h"][factsPath] = fakeFile{"120000", "blob", factsWant}
	want(t, "a link with the bytes", Basic(f, input(factsRecord)), Failure, ReasonNotARegularFile, factsPath)
	delete(f.trees["h"], factsPath)
	want(t, "the record deleted", Basic(f, input(factsRecord)), Failure, ReasonNotEqual, factsPath)
	// The head's copy of the rule-paths.tsv, forged to list nothing, is no
	// help: nothing of the head decides what a rule path is.
	f.names = []string{"AGENTS.md", GuardrailsPath}
	f.trees["h"][GuardrailsPath] = reg(head)
	f.diffs = map[string][]byte{GuardrailsPath: diff("@@ -5,0 +6 @@\n+- b\n")}
	want(t, "a rule path beside the exception", Basic(f, input()), Failure, ReasonRulePath, "AGENTS.md")
	f.names = []string{GuardrailsPath}
	want(t, "the exception alone", Basic(f, input()), Success, ReasonNoRulePath, "")
}

func TestBasicPreImageAndBaseChains(t *testing.T) {
	prdPost := func(f *fake) {
		f.names = []string{prdPath}
		f.trees["h"][prdPath] = reg(string(prdRendered(t)))
	}
	f := scene()
	prdPost(f)
	want(t, "the pre-image commit", Basic(f, input(prdRecord)), Success, ReasonRenderedRecord, "")
	for name, mutate := range map[string]func(*fake, *Input){
		"no PreRev":                      func(f *fake, in *Input) { in.PreRev = "" },
		"an unknown PreRev":              func(f *fake, in *Input) { in.PreRev = "nothing" },
		"PreRev the head":                func(f *fake, in *Input) { f.ids["pre"] = "h" },
		"PreRev not an ancestor of head": func(f *fake, in *Input) { delete(f.anc, "p>h") },
		"the base not an ancestor of it": func(f *fake, in *Input) { delete(f.anc, "b>p") },
		"no PRD at the pre-image":        func(f *fake, in *Input) { f.trees["p"] = map[string]fakeFile{} },
		"a link at the pre-image":        func(f *fake, in *Input) { f.trees["p"][prdPath] = fakeFile{"120000", "blob", prdPre} },
	} {
		f := scene()
		prdPost(f)
		in := input(prdRecord)
		mutate(f, &in)
		want(t, name, Basic(f, in), Error, ReasonInput, "")
	}
	f = scene()
	prdPost(f)
	f.trees["b"][prdPath] = reg(prdPre)
	f.ids["pre"] = "b"
	want(t, "the base as the pre-image (a later bet)", Basic(f, input(prdRecord)), Success, ReasonRenderedRecord, "")
	// A PRD with no record named needs no PreRev.
	f = scene()
	in := input()
	in.PreRev = ""
	want(t, "no PRD record, no PreRev", Basic(f, in), Failure, ReasonRulePath, "AGENTS.md")

	f = scene()
	delete(f.anc, "b>h")
	want(t, "a base that is not an ancestor of the head", Basic(f, input()), Error, ReasonBase, "")
	for _, call := range []string{"RevParse", "LsTree", "DiffNames", "IsAncestor"} {
		f = scene()
		f.errs[call] = errors.New("git failed")
		want(t, "a failing "+call, Basic(f, input()), Error, ReasonGit, "")
	}
	f = scene()
	f.names = []string{GuardrailsPath}
	f.errs["DiffFile"] = errors.New("git failed")
	want(t, "a failing DiffFile", Basic(f, input()), Error, ReasonGit, "")
	f = scene()
	f.names = []string{factsPath}
	f.errs["Show"] = errors.New("git failed")
	// Show fails for the register first: that is a register that cannot be read.
	if r := Basic(f, input(factsRecord)); r.State != Error {
		t.Errorf("a failing Show: %+v, want an error", r)
	}
}

func TestBasicRefusesABadRecordBeforeItReadsAnything(t *testing.T) {
	f := scene()
	bad := factsRecord
	bad.Path = "../x"
	want(t, "a path out of the tree", Basic(f, input(bad)), Error, ReasonInput, "../x")
	if len(f.shown) != 0 {
		t.Errorf("Show was called for a bad input: %v", f.shown)
	}
}
