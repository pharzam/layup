package setup

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// The steps of the setup in phase 2 (task T-ktgx, #273): the five steps that
// differ from phase 1, each with the cases of its words in scaffold.md, "The
// steps of the setup in phase 2". Each case breaks one rule only, and each
// refusal is compared with its whole reason, so a case that another rule
// refuses fails no mutation of its own.

// phase2 gives the inputs of phase 2 of these tests: the pin of pinRecord, the
// questions of the batch, and the sum of a brief.
func phase2() Phase2 {
	return Phase2{Name: "acme/widget", Visibility: "public", PinSource: baseURL, PinCommit: commitA, PinTree: treeA, PinTime: "2026-09-30T08:00:00Z", BriefSum: "abc", AppID: 42,
		Questions: []work.Question{{ID: "Q-001", Text: "Which technology stack does the product use?"}, {ID: "Q-002", Text: "What does fast mean here, as a number?"}}}
}

// answers2 gives the answers of phase 2, the stack and the two Q- rows of ans;
// change is as for ans, and a row of S01-name, S01-visibility or S01-baseline
// that a change adds stays.
func answers2(change map[string]string) work.Answers {
	c := map[string]string{"S01-name": "", "S01-visibility": "", "S01-baseline": ""}
	for k, v := range change {
		c[k] = v
	}
	return ans(c)
}

func runPhase2(t *testing.T, id string, p Phase2, in Input) Outcome {
	t.Helper()
	return StepsPhase2(p, Calls{})[id].Run(in)
}

// S01 asks the stack and nothing else: one stop row, with no answer and with
// an answer of the empty value, and no row for the other three questions of
// phase 1 or for a Q- question.
func TestPhase2S01AsksTheStackAlone(t *testing.T) {
	want := []StopRow{{"S01", "S01-stack", work.S01Questions[0].Text, ""}}
	for name, a := range map[string]work.Answers{
		"no answer":                    nil,
		"an answer of the empty value": {{"S01-stack", "", "operator", "u", ""}},
		"the Q- rows alone":            {{"Q-001", "Go", "idea-owner", "u", ""}},
	} {
		if o := runPhase2(t, "S01", phase2(), Input{Dir: t.TempDir(), Answers: a}); o.Kind != Stop || !slices.Equal(o.Stops, want) {
			t.Errorf("%s: %s %+v; want a stop with %+v", name, o.Kind, o.Stops, want)
		}
	}
}

// S01 takes the name, the visibility and the baseline from the records of
// Start, and the brief's sum from psb.sha256; it runs no layup psb check, so a
// Q- row of any ID is no input error, and a row of S01-name is not read.
func TestPhase2S01TakesTheValuesOfStart(t *testing.T) {
	const ref = "records:start/start.tsv"
	want := [][]string{{"stack", "go", "answer", "S01-stack"}, {"name", "acme/widget", "computed", ref}, {"visibility", "public", "computed", ref},
		{"baseline", baseURL, "computed", ref}, {"brief.sha256", "abc", "computed", "sha256 " + work.BriefPath}}
	for name, a := range map[string]work.Answers{
		"the stack and two Q- rows":                answers2(nil),
		"a Q- row of no question of the batch":     answers2(map[string]string{"Q-009": "x"}),
		"a row of S01-name that names another one": answers2(map[string]string{"S01-name": "other/name"}),
	} {
		o := runPhase2(t, "S01", phase2(), Input{Dir: t.TempDir(), Answers: a})
		if o.Kind != Done || o.Evidence != "S01-stack answered; the stack has a catalog entry" || !reflect.DeepEqual(o.Values, want) {
			t.Errorf("%s: %s %q %q; want done with %q", name, o.Kind, o.Evidence, o.Values, want)
		}
	}
	a := answers2(nil)
	a[0][3] = "F-0003#5" // the answer's source is a fact of the target
	want[0] = []string{"stack", "go", "fact", "F-0003#5"}
	if o := runPhase2(t, "S01", phase2(), Input{Dir: t.TempDir(), Answers: a}); o.Kind != Done || !reflect.DeepEqual(o.Values, want) {
		t.Errorf("an answer that cites a fact: %s %q; want the row %q", o.Kind, o.Values, want[0])
	}
}

// S01 makes the checks of phase 1 on the same values, and a value that fails
// is a fail, not an input error.
func TestPhase2S01ChecksTheValues(t *testing.T) {
	with := func(change func(*Phase2)) Phase2 { p := phase2(); change(&p); return p }
	for _, c := range []struct {
		name   string
		p      Phase2
		a      work.Answers
		reason string
	}{
		{"a stack of no entry", phase2(), answers2(map[string]string{"S01-stack": "rust"}), "the stack rust has no entry in the catalog"},
		{"a name that is not OWNER/NAME", with(func(p *Phase2) { p.Name = "widget" }), answers2(nil), "name: widget is not OWNER/NAME"},
		{"a visibility of another word", with(func(p *Phase2) { p.Visibility = "internal" }), answers2(nil), "visibility: internal is not public or private"},
		{"a baseline that is not a URL", with(func(p *Phase2) { p.PinSource = "git@github.com:pharzam/armature.git" }), answers2(nil),
			"baseline: git@github.com:pharzam/armature.git is not a URL with the scheme https, http, git or file"},
	} {
		if o := runPhase2(t, "S01", c.p, Input{Dir: t.TempDir(), Answers: c.a}); o.Kind != Fail || o.Evidence != c.reason || o.Values != nil {
			t.Errorf("%s: %s %q; want fail with %q", c.name, o.Kind, o.Evidence, c.reason)
		}
	}
}

// A done S01 compares the sum of the brief with psb.sha256 of start.tsv, as
// phase 1 compares the problem statement.
func TestPhase2S01ChecksItsBriefAgain(t *testing.T) {
	s := StepsPhase2(phase2(), Calls{})["S01"]
	if err := s.Unchanged(work.Record{{"S01", "brief.sha256", "abc", "computed", "x"}}); err != nil {
		t.Errorf("the same sum: %v", err)
	}
	if err := s.Unchanged(work.Record{{"S01", "brief.sha256", "zzz", "computed", "x"}}); err == nil || !strings.Contains(err.Error(), "zzz") {
		t.Errorf("another sum: %v; want an error that names it", err)
	}
	if err := s.Unchanged(nil); err == nil || !strings.Contains(err.Error(), "no row brief.sha256") {
		t.Errorf("no row: %v; want an error", err)
	}
}

// S02 resolves nothing: no call of git, no clone, no write; the four rows are
// those of start.tsv with source computed and the ref of the records of Start
// (NFR-006), and the actor of the step is layup-run.
func TestPhase2S02ResolvesNothing(t *testing.T) {
	f := &fakeRepo{}
	f.install(t)
	const ref = "records:start/start.tsv"
	want := [][]string{{"pin.source", baseURL, "computed", ref}, {"pin.commit", commitA, "computed", ref}, {"pin.tree", treeA, "computed", ref},
		{"pin.time", "2026-09-30T08:00:00Z", "computed", ref}}
	o := runPhase2(t, "S02", phase2(), Input{Dir: "w"})
	if o.Kind != Done || o.Evidence != "the commit and the tree of start.tsv" || !reflect.DeepEqual(o.Values, want) || f.calls != nil || len(f.files) != 0 {
		t.Errorf("S02: %s %q %q, calls %q, %d files; want done with %q and no call", o.Kind, o.Evidence, o.Values, f.calls, len(f.files), want)
	}
	if a := StepsPhase2(phase2(), Calls{})["S02"].Actor; a != "layup-run" {
		t.Errorf("the actor of S02 is %q; want layup-run", a)
	}
}

// S03 makes no commit and no command: it reads the tree of main in the clone
// that layup run made, and checks that it is pin.tree.
func TestPhase2S03MakesNoCommit(t *testing.T) {
	f := &fakeRepo{revs: map[string]string{"refs/heads/main^{tree}": treeA}}
	f.install(t)
	o := runPhase2(t, "S03", phase2(), Input{Dir: "w", Record: pinRecord(), Who: Who})
	if o.Kind != Done || o.Evidence != "root tree = pin.tree" || o.Commit || o.Values != nil || !slices.Equal(f.calls, []string{"rev-parse refs/heads/main^{tree}"}) {
		t.Errorf("S03: %s %q %v, calls %q; want done, no commit and the one read of the tree", o.Kind, o.Evidence, o.Commit, f.calls)
	}
	other := "2222222222222222222222222222222222222222"
	for _, c := range []struct {
		name   string
		f      *fakeRepo
		record work.Record
		reason string
	}{
		{"another tree", &fakeRepo{revs: map[string]string{"refs/heads/main^{tree}": other}}, pinRecord(), "the root tree " + other + " is not pin.tree " + treeA},
		{"an error of git", &fakeRepo{fail: "rev-parse"}, pinRecord(), "the tree of main: rev-parse refs/heads/main^{tree}: failed"},
		{"no row pin.tree", &fakeRepo{revs: map[string]string{"refs/heads/main^{tree}": treeA}}, pinRecord()[:5], "the record has no row S02 pin.tree"},
	} {
		c.f.install(t)
		if o := runPhase2(t, "S03", phase2(), Input{Dir: "w", Record: c.record, Who: Who}); o.Kind != Fail || o.Evidence != c.reason {
			t.Errorf("%s: %s %q; want fail with %q", c.name, o.Kind, o.Evidence, c.reason)
		}
	}
	if StepsPhase2(phase2(), Calls{})["S03"].Commands != nil || Steps(Brief{}, Calls{})["S03"].Commands == nil {
		t.Error("S03 gives no command in phase 2, and its two commands stay in phase 1")
	}
}

// S04 writes the pin file, the decision record and the raw fact, as in phase 1,
// from the pin rows of S02; its raw fact holds the answer to S01-stack and the
// Q- rows of the questions of the batch, each with the text of its question,
// and not a question of the batch that has no row.
func TestPhase2S04WritesTheRawFactOfTheBatch(t *testing.T) {
	f := s04Repo("refs/heads/main", rootA)
	f.install(t)
	a, p := answers2(nil), phase2()
	p.Questions = append(slices.Clone(p.Questions), work.Question{ID: "Q-003", Text: "A question with no answer"})
	o := runPhase2(t, "S04", p, Input{Dir: "w", Record: pinRecord(), Answers: a})
	record := answersRecord("F-0001", "2026-09-30", setupRecord, append([]work.Question{work.S01Questions[0]}, phase2().Questions...), a)
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(record)))
	files := map[string]string{
		"w/target/docs/setup/armature.pin":            pinText(baseURL, commitA, treeA, "2026-09-30"),
		"w/target/docs/adr/0009-pin-the-baseline.md":  adrText(9, "2026-09-30", baseURL, commitA, treeA),
		"w/target/docs/adr/README.md":                 strings.Replace(adrIndex, "| H | Accepted |\n", "| H | Accepted |\n| [0009](0009-pin-the-baseline.md) | Pin the baseline | Accepted |\n", 1),
		"w/target/docs/facts/F-0001-setup-answers.md": record,
		"w/target/docs/facts/README.md":               strings.Replace(factsIndex, "| _none yet_ | | | |", "| [F-0001](F-0001-setup-answers.md) | The answers to the questions of the setup | 2026-09-30 | Raw |", 1),
		"w/target/docs/setup/facts.sha256":            sum + "  docs/facts/F-0001-setup-answers.md\n",
	}
	values := [][]string{
		{"pin.adr", "docs/adr/0009-pin-the-baseline.md", "computed", "the next free number of docs/adr/"},
		{"answers.record", "docs/facts/F-0001-setup-answers.md", "computed", "the next free ID of docs/facts/"},
		{"answers.record.sha256", sum, "computed", "sha256 docs/facts/F-0001-setup-answers.md"},
	}
	if o.Kind != Done || o.Evidence != "checks pin and facts" || !o.Commit || !reflect.DeepEqual(o.Values, values) || !reflect.DeepEqual(f.files, files) {
		t.Errorf("S04: %s %q %v %q, %d files; want done, the values %q and the six files of phase 1", o.Kind, o.Evidence, o.Commit, o.Values, len(f.files), values)
	}
	for _, q := range []string{"`S01-stack` go", "`Q-001` Go", "the question: Which technology stack does the product use?", "`Q-002` under 200 ms"} {
		if !strings.Contains(record, q) {
			t.Errorf("the raw fact does not hold %q:\n%s", q, record)
		}
	}
	if strings.Contains(record, "Q-003") {
		t.Errorf("the raw fact holds a question with no row of the answers:\n%s", record)
	}
	// A Q- row that the batch does not hold: its text is not known.
	f = s04Repo("refs/heads/main", rootA)
	f.install(t)
	if o := runPhase2(t, "S04", phase2(), Input{Dir: "w", Record: pinRecord(), Answers: answers2(map[string]string{"Q-009": "x"})}); o.Kind != Fail ||
		o.Evidence != "the question Q-009 is not one of the questions of the batch, so its text is not known" || len(f.files) != 0 {
		t.Errorf("a Q- row of no question of the batch: %s %q, %d files; want fail and no file", o.Kind, o.Evidence, len(f.files))
	}
}

// S15 keeps phase 1's step, with no records commit and no command that pushes
// layup-records: layup run commits the records.
func TestPhase2S15MakesNoRecordsCommit(t *testing.T) {
	two, one := StepsPhase2(phase2(), Calls{})["S15"], Steps(Brief{}, Calls{})["S15"]
	if two.Records != nil || two.Commands != nil || one.Records == nil || one.Commands == nil {
		t.Errorf("S15 in phase 2 has Records %v and Commands %v; in phase 1 %v and %v; want nil in phase 2 and both in phase 1",
			two.Records != nil, two.Commands != nil, one.Records != nil, one.Commands != nil)
	}
	verify := "check\tresult\treason\npin\tpass\t\u2014\ngate:layout\tclear\tpending: fixture not run\n"
	record := append(lastRecord(), []string{"S12", "catalog:.github/gates.sh", "x", "catalog", "go/files/.github/gates.sh.tmpl"})
	var outcomes []Outcome
	var written []map[string]string
	for _, s := range []Step{one, two} {
		f := &fakeRepo{revs: map[string]string{"refs/heads/layup-setup^{commit}": setupHead}, disk: map[string]string{"w/out/verify.tsv": verify},
			trees: map[string][]git.TreeEntry{".": {{Path: ".github/gates.sh"}, {Path: "docs/x.sh"}, {Path: "go.mod"}}}, shows: map[string]string{"docs/gates.tsv": twoKinds}}
		f.install(t)
		outcomes, written = append(outcomes, s.Run(Input{Dir: "w", Record: record})), append(written, f.files)
	}
	if outcomes[0].Kind != Done || !reflect.DeepEqual(outcomes[0], outcomes[1]) || !reflect.DeepEqual(written[0], written[1]) {
		t.Errorf("S15 in phase 2: %+v, files %v; want the outcome and the files of phase 1: %+v, %v", outcomes[1], written[1], outcomes[0], written[0])
	}
}

// The other steps are phase 1's: the same code, the same hand-off and the same
// answers read. S01 to S04 have the code of phase 2, S15 keeps its code, and
// S02, S03 and S15 are layup run's in the step table.
func TestPhase2KeepsTheOtherStepsOfPhase1(t *testing.T) {
	one, two := Steps(Brief{}, Calls{}), StepsPhase2(phase2(), Calls{})
	changed := []string{"S01", "S02", "S03", "S04", "S13"}
	for _, id := range ids() {
		a, b := one[id], two[id]
		same := reflect.ValueOf(a.Run).Pointer() == reflect.ValueOf(b.Run).Pointer()
		actor := a.Actor
		if slices.Contains([]string{"S02", "S03", "S15"}, id) {
			actor = "layup-run"
		}
		if same == slices.Contains(changed, id) || b.Actor != actor || a.HandOff != b.HandOff || !slices.Equal(a.Reads, b.Reads) || (a.Evidence == nil) != (b.Evidence == nil) {
			t.Errorf("%s: same code %v, actor %q, hand-off %v, reads %q; want the code of phase 1 exactly when it is not one of %q, and the actor %q",
				id, same, b.Actor, b.HandOff, b.Reads, changed, actor)
		}
	}
}

// S04 keeps the evidence of phase 1: after its commit the runner runs the
// checks pin and facts (NFR-006: check pin compares the file with the record
// rows, as in phase 1).
func TestPhase2S04KeepsItsEvidence(t *testing.T) {
	var asked [][]string
	calls := Calls{Checks: func(dir string, names []string) string {
		asked = append(asked, append([]string{dir}, names...))
		return ""
	}}
	s := StepsPhase2(phase2(), calls)["S04"]
	if s.Evidence == nil || s.Evidence("w") != "" || !reflect.DeepEqual(asked, [][]string{{"w", "pin", "facts"}}) {
		t.Errorf("the evidence of S04 in phase 2 asked %q; want the checks pin and facts on w", asked)
	}
}

// S04 of phase 1 keeps its refusal of a gap table of another form as an input
// error, and writes no file: the body that both phases share (runS04With) gets
// its questions from the phase.
func TestS04OfPhase1RefusesABadGapTable(t *testing.T) {
	f := s04Repo("refs/heads/main", rootA)
	f.install(t)
	o := runS04(Brief{Gaps: []byte("id\trule\n")}, Input{Dir: "w", Record: pinRecord(), Answers: ans(nil)})
	if o.Kind != Invalid || !strings.Contains(o.Evidence, "the gap table of "+work.BriefPath) || len(f.files) != 0 {
		t.Errorf("a gap table of another form: %s %q, %d files; want an input error and no file", o.Kind, o.Evidence, len(f.files))
	}
}
