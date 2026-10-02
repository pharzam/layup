package setup

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// table is a gap table as internal/psb writes it: one gap with no line (G1)
// and one at line 3 (G4).
const table = "id\trule\tline\texcerpt\tquestion\n" +
	"Q-001\tG1\t0\t\u2014\tWhich technology stack does the product use?\n" +
	"Q-002\tG4\t3\tWe need it fast.\tWhat does fast mean here, as a number?\n"

// The gap table that internal/cli hands to S01 reads by the block psb-gaps
// (D5 of #86); a table of another form is an error.
func TestTheGapTable(t *testing.T) {
	gaps, err := readGaps([]byte(table))
	want := []Gap{{"Q-001", 0, "Which technology stack does the product use?"}, {"Q-002", 3, "What does fast mean here, as a number?"}}
	if err != nil || !slices.Equal(gaps, want) {
		t.Errorf("readGaps = %+v, %v; want %+v", gaps, err, want)
	}
	if gaps, err := readGaps([]byte("id\trule\tline\texcerpt\tquestion\n")); err != nil || len(gaps) != 0 {
		t.Errorf("a table with no gap: %+v, %v", gaps, err)
	}
	for _, bad := range []string{"", "id\trule\n", "id\trule\tline\texcerpt\tquestion\nQ-1\tG1\t0\t\u2014\tq\n"} {
		if _, err := readGaps([]byte(bad)); err == nil {
			t.Errorf("readGaps(%q): no error", bad)
		}
	}
}

// ans gives the answers of S01, with the changes: a question and its answer,
// "" removes the row.
func ans(change map[string]string) work.Answers {
	rows := [][]string{
		{"S01-stack", "go"}, {"S01-name", "acme/widget"}, {"S01-visibility", "public"},
		{"S01-baseline", "https://github.com/pharzam/armature"}, {"Q-001", "Go"}, {"Q-002", "under 200 ms"},
	}
	var a work.Answers
	for _, r := range rows {
		v, ok := change[r[0]]
		if ok && v == "" {
			continue
		}
		if ok {
			r[1] = v
		}
		a = append(a, []string{r[0], r[1], "operator", "https://github.invalid/acme/widget/issues/1#issuecomment-1", ""})
	}
	for q, v := range change {
		if !slices.ContainsFunc(rows, func(r []string) bool { return r[0] == q }) && v != "" {
			a = append(a, []string{q, v, "idea-owner", "https://github.invalid/acme/widget/issues/1#issuecomment-2", ""})
		}
	}
	return a
}

func runS01(t *testing.T, a work.Answers) Outcome {
	t.Helper()
	return Steps(Brief{Gaps: []byte(table), Sum: "abc"}, nil)["S01"].Run(Input{Dir: t.TempDir(), Answers: a})
}

// S01 stops once for each missing answer, of both kinds, in the order of the
// stop table, with the question texts of internal/work and of the gap table,
// and the line of a gap (K33); an empty answer is a missing one (D4 of #86).
func TestS01StopsForEachMissingAnswer(t *testing.T) {
	o := runS01(t, nil)
	var want []StopRow
	for _, q := range work.S01Questions {
		want = append(want, StopRow{"S01", q.ID, q.Text, ""})
	}
	want = append(want, StopRow{"S01", "Q-001", "Which technology stack does the product use?", ""},
		StopRow{"S01", "Q-002", "What does fast mean here, as a number?", work.BriefPath + ":3"})
	if o.Kind != Stop || !slices.Equal(o.Stops, want) {
		t.Errorf("no answer: %s %+v; want a stop with\n%+v", o.Kind, o.Stops, want)
	}
	a := ans(map[string]string{"S01-name": ""})
	for _, r := range a {
		if r[0] == "Q-002" {
			r[1] = "" // the reader of answers.tsv gives the empty value for \u2014
		}
	}
	o = runS01(t, a)
	if o.Kind != Stop || len(o.Stops) != 2 || o.Stops[0].Question != "S01-name" || o.Stops[1].Question != "Q-002" {
		t.Errorf("two missing answers: %s %+v; want a stop for S01-name and Q-002", o.Kind, o.Stops)
	}
}

// An answer to a Q- question that the gap table does not hold is an input
// error (D4, D11 of #86).
func TestS01RefusesAnAnswerToNoGap(t *testing.T) {
	if o := runS01(t, ans(map[string]string{"Q-009": "x"})); o.Kind != Invalid || !strings.Contains(o.Evidence, "Q-009") {
		t.Errorf("an answer to Q-009: %s %q; want an input error that names Q-009", o.Kind, o.Evidence)
	}
}

// S01 with each answer: the checks of the stack, the name, the visibility and
// the baseline (D1, D2, D4 of #86), and the record rows with their sources
// (D3).
func TestS01ChecksTheAnswers(t *testing.T) {
	o := runS01(t, ans(nil))
	want := [][]string{
		{"stack", "go", "answer", "S01-stack"}, {"name", "acme/widget", "answer", "S01-name"},
		{"visibility", "public", "answer", "S01-visibility"}, {"baseline", "https://github.com/pharzam/armature", "answer", "S01-baseline"},
		{"brief.sha256", "abc", "computed", "sha256 " + work.BriefPath},
	}
	if o.Kind != Done || o.Evidence != "every answer present; the stack has a catalog entry" || !slices.EqualFunc(o.Values, want, slices.Equal) {
		t.Errorf("each answer: %s %q %q; want done with %q", o.Kind, o.Evidence, o.Values, want)
	}
	for _, c := range []struct{ q, a, reason string }{
		{"S01-stack", "rust", "the stack rust has no entry in the catalog"},
		{"S01-stack", "test", "the stack test has no entry in the catalog"},
		{"S01-name", "widget", "S01-name: widget is not OWNER/NAME"},
		{"S01-name", "acme/widget/x", "is not OWNER/NAME"},
		{"S01-name", "-acme/widget", "is not OWNER/NAME"},
		{"S01-name", "acme-/widget", "is not OWNER/NAME"},
		{"S01-name", "ac--me/widget", "is not OWNER/NAME"},
		{"S01-name", "acme/..", "is not OWNER/NAME"},
		{"S01-name", strings.Repeat("a", 40) + "/w", "is not OWNER/NAME"},
		{"S01-name", "acme/" + strings.Repeat("w", 101), "is not OWNER/NAME"},
		{"S01-visibility", "internal", "S01-visibility: internal is not public or private"},
		{"S01-baseline", "ssh://github.com/pharzam/armature", "S01-baseline: ssh://github.com/pharzam/armature is not a URL"},
		{"S01-baseline", "git@github.com:pharzam/armature.git", "is not a URL"},
		{"S01-baseline", "github.com/pharzam/armature", "is not a URL"},
	} {
		if o := runS01(t, ans(map[string]string{c.q: c.a})); o.Kind != Fail || !strings.Contains(o.Evidence, c.reason) {
			t.Errorf("%s %s: %s %q; want fail with %q", c.q, c.a, o.Kind, o.Evidence, c.reason)
		}
	}
	for _, good := range []map[string]string{
		{"S01-name": "a/b"}, {"S01-name": "acme-co/widget.v2_x-y"}, {"S01-name": strings.Repeat("a", 39) + "/" + strings.Repeat("w", 100)},
		{"S01-visibility": "private"}, {"S01-baseline": "file:///tmp/armature"}, {"S01-baseline": "http://host/x.git"}, {"S01-baseline": "git://host/x"},
	} {
		if o := runS01(t, ans(good)); o.Kind != Done {
			t.Errorf("%v: %s %q; want done", good, o.Kind, o.Evidence)
		}
	}
}

// An answer whose source is a fact citation gives a row of source fact (D3).
func TestS01RowOfAFact(t *testing.T) {
	a := ans(nil)
	for _, r := range a {
		if r[0] == "S01-stack" {
			r[3] = "F-0003#44"
		}
	}
	o := runS01(t, a)
	if o.Kind != Done || !slices.Equal(o.Values[0], []string{"stack", "go", "fact", "F-0003#44"}) {
		t.Errorf("a fact source: %s %q; want the row stack go fact F-0003#44", o.Kind, o.Values)
	}
}

// Once S01 is done, a problem statement whose SHA-256 differs from the row
// brief.sha256 is an input error of the run (D6 of #86).
func TestS01UnchangedBrief(t *testing.T) {
	u := Steps(Brief{Gaps: []byte(table), Sum: "abc"}, nil)["S01"].Unchanged
	if u == nil {
		t.Fatal("S01 has no check of its inputs")
	}
	rec := func(sum string) work.Record {
		r := work.Record{{"S01", "stack", "go", "answer", "S01-stack"}}
		if sum != "" {
			r = append(r, []string{"S01", "brief.sha256", sum, "computed", "sha256 " + work.BriefPath})
		}
		return r
	}
	if err := u(rec("abc")); err != nil {
		t.Errorf("the same brief: %v", err)
	}
	if err := u(rec("def")); err == nil || !strings.Contains(err.Error(), "an input that changed after a step read it") || !strings.HasPrefix(err.Error(), work.BriefPath) {
		t.Errorf("another brief: %v; want an error that names the brief", err)
	}
	if err := u(rec("")); err == nil || !strings.Contains(err.Error(), "brief.sha256") {
		t.Errorf("no row brief.sha256: %v", err)
	}
}

// A step that ends with an input error ends the run with that error and no
// table; the runner writes no row and no commands.sh (D11 of #86).
func TestTheRunnerGivesTheInputErrorOfAStep(t *testing.T) {
	f := &fakeSys{}
	f.install(t)
	res, err := Run("w", steps(f, map[string]Outcome{"S01": {Kind: Invalid, Evidence: "inputs/answers.tsv: line 6, Q-009: no step of this run asked this question"}}), testWho, noStep)
	var input *InputError
	if !errors.As(err, &input) || !strings.Contains(err.Error(), "Q-009") || res.Steps != nil || res.Stops != nil || f.record != nil || f.commands != "" {
		t.Errorf("an input error of S01: %+v, %v, record %q, commands %q; want the input error and nothing written", res, err, f.record, f.commands)
	}
}

// The check of a done step's inputs runs before any step; an error is an input
// error (D6 of #86).
func TestTheRunnerChecksTheInputsOfADoneStep(t *testing.T) {
	f := &fakeSys{record: work.Record{{"S01", "done", "x", "step", ""}}}
	f.install(t)
	m := steps(f, nil)
	s := m["S01"]
	s.Reads = nil
	s.Unchanged = func(work.Record) error {
		return errors.New(work.BriefPath + ": an input that changed after a step read it")
	}
	m["S01"] = s
	_, err := Run("w", m, testWho, noStep)
	var input *InputError
	if !errors.As(err, &input) || f.calls != nil {
		t.Errorf("a changed input of a done step: %v, calls %q; want an input error before any step", err, f.calls)
	}
}

// The evidence of a step runs after its commit: each check passes, and the
// step is done; a check fails, and the step fails with its reason, with no
// done row, and the runner moves layup-setup back to the step's parent, but
// only when the head is the commit that the step made (D12 of #86, condition 1
// of its plan review).
func TestTheEvidenceOfAStep(t *testing.T) {
	pin := work.Record{{"S02", "pin.time", "2026-10-02T09:30:00Z", "computed", "the clock of the LAYUP host"}}
	for _, c := range []struct {
		name, reason string
		heads        []string
		done         bool
		resets       []string
	}{
		{"passes", "", []string{"p", "m"}, true, nil},
		{"fails", "pin: the pin file does not match", []string{"p", "m", "m"}, false, []string{"p"}},
		{"fails, and the head moved", "facts: x", []string{"p", "m", "other"}, false, nil},
	} {
		f := &fakeSys{record: pin, heads: c.heads}
		f.install(t)
		m := steps(f, map[string]Outcome{"S04": {Kind: Done, Evidence: "the pin", Commit: true}})
		s := m["S04"]
		s.Evidence = func(string) string { return c.reason }
		m["S04"] = s
		res, err := Run("w", m, testWho, noStep)
		_, done := f.record.Value("S04", "done")
		row := res.Steps[3]
		if err != nil || done != c.done || !slices.Equal(f.resets, c.resets) || len(f.commits) != 1 ||
			(!c.done && (row.Result != Fail || row.Evidence != c.reason)) {
			t.Errorf("%s: %v, done %v, resets %q, commits %q, row %+v; want done %v, resets %q", c.name, err, done, f.resets, f.commits, row, c.done, c.resets)
		}
	}
}

// The fixed text of the pin file, in the form of NFR-006 (D10 of #86).
func TestThePinText(t *testing.T) {
	got := pinText("https://github.com/pharzam/armature", "a95965534b14b0bf14ad74da0c9a45b5f4aedf88", "8ffb250afd584da8b418bc220fb6d72e802924ce", "2026-10-02")
	want := "source=https://github.com/pharzam/armature\ncommit=a95965534b14b0bf14ad74da0c9a45b5f4aedf88\ntree=8ffb250afd584da8b418bc220fb6d72e802924ce\n" +
		"method=git clone https://github.com/pharzam/armature, checkout a95965534b14b0bf14ad74da0c9a45b5f4aedf88, .git removed\ndate=2026-10-02\n"
	if got != want {
		t.Errorf("pinText =\n%s\nwant\n%s", got, want)
	}
}

// The decision record of the pin, in the baseline's form (D10 of #86): its
// file name, its title line, its date, its status, its three sections with the
// pin values, and no reference to an issue.
func TestTheDecisionRecordOfThePin(t *testing.T) {
	if f := adrFile(9); f != "0009-pin-the-baseline.md" {
		t.Errorf("adrFile(9) = %q", f)
	}
	text := adrText(9, "2026-10-02", "https://github.com/pharzam/armature", "a95965534b14b0bf14ad74da0c9a45b5f4aedf88", "8ffb250afd584da8b418bc220fb6d72e802924ce")
	for _, want := range []string{"# 0009. Pin the baseline\n\nDate: 2026-10-02\n\n## Status\n\nAccepted\n\n## Context\n", "\n## Decision\n", "\n## Consequences\n",
		"`https://github.com/pharzam/armature`", "`a95965534b14b0bf14ad74da0c9a45b5f4aedf88`", "`8ffb250afd584da8b418bc220fb6d72e802924ce`", "`docs/setup/armature.pin`"} {
		if !strings.Contains(text, want) {
			t.Errorf("the record holds no %q:\n%s", want, text)
		}
	}
	if strings.Contains(text, "#1") || strings.Contains(text, "issue") || strings.Contains(text, "pull request") {
		t.Errorf("the record names an issue:\n%s", text)
	}
}

// A row of an index table goes after the last row of the table under the
// heading Index; a table whose one row is the placeholder _none yet_ gets the
// row in its place (D10 of #86); a file with no such table is an error.
func TestAddIndexRow(t *testing.T) {
	adr := "# ADRs\n\n## Index\n\n| ADR | Title | Status |\n| --- | ----- | ------ |\n| [0001](0001-a.md) | A | Accepted |\n\n<!-- Add one row per ADR. -->\n"
	got, err := addIndexRow(adr, "| [0002](0002-b.md) | B | Accepted |")
	if err != nil || got != strings.Replace(adr, "| A | Accepted |\n", "| A | Accepted |\n| [0002](0002-b.md) | B | Accepted |\n", 1) {
		t.Errorf("the ADR index:\n%s\n%v", got, err)
	}
	facts := "# Facts\n\n## Index\n\n| Fact doc | Source | Collected | Status |\n| -------- | ------ | --------- | ------ |\n| _none yet_ | | | |\n"
	got, err = addIndexRow(facts, "| [F-0001](F-0001-setup-answers.md) | The answers | 2026-10-02 | Raw |")
	if err != nil || got != strings.Replace(facts, "| _none yet_ | | | |\n", "| [F-0001](F-0001-setup-answers.md) | The answers | 2026-10-02 | Raw |\n", 1) {
		t.Errorf("the facts index:\n%s\n%v", got, err)
	}
	end := "# ADRs\n\n## Index\n\n| ADR | Title | Status |\n| --- | ----- | ------ |\n| [0001](0001-a.md) | A | Accepted |"
	if got, err := addIndexRow(end, "| [0002](0002-b.md) | B | Accepted |"); err != nil || got != end+"\n| [0002](0002-b.md) | B | Accepted |\n" {
		t.Errorf("a table at the end of a file with no line feed:\n%q\n%v", got, err)
	}
	for _, bad := range []string{"# Facts\n\nno index\n", "# Facts\n\n## Index\n\nno table\n"} {
		if _, err := addIndexRow(bad, "| x |"); err == nil {
			t.Errorf("addIndexRow(%q): no error", bad)
		}
	}
}

// The answers record of S04, in the form of K15 and K17 (D10 of #86, with
// condition 2 of its plan review): the header table with each value, the date
// of pin.time, one fact per answer in the order of the stop table, and each
// angle quote of a recorded text written as an entity: of the answer, of its
// source and of the question (finding 1 of review round 1).
func TestTheAnswersRecord(t *testing.T) {
	asked := append(slices.Clone(work.S01Questions), work.Question{ID: "Q-001", Text: "Which stack? \u2039x\u203a"})
	a := work.Answers{
		{"Q-001", "Go \u20391.26\u203a", "idea-owner", "said \u2039here\u203a", ""},
		{"S01-stack", "go", "operator", "https://github.invalid/c/1", ""},
		{"S01-name", "acme/widget", "operator", "https://github.invalid/c/1", ""},
		{"S01-visibility", "public", "operator", "https://github.invalid/c/1", ""},
		{"S01-baseline", "https://github.com/pharzam/armature", "operator", "F-0003#5", ""},
	}
	text := answersRecord("F-0001", "2026-10-02", asked, a)
	want := "# F-0001. The answers to the questions of the setup\n\n| Field | Value |\n| ------------ | ----- |\n" +
		"| Fact ID | `F-0001` |\n" +
		"| Source | The answers of `inputs/answers.tsv` to the questions of S01 and to the gaps of the problem statement |\n" +
		"| Collected by | `layup setup`, step S04 |\n| Date collected | 2026-10-02 |\n" +
		"| Origin | `inputs/answers.tsv` of the work area |\n| Status | `Raw` |\n\n## Facts as collected\n\n" +
		"1. `S01-stack` go \u2014 by operator; source https://github.invalid/c/1; the question: " + work.S01Questions[0].Text + "\n" +
		"2. `S01-name` acme/widget \u2014 by operator; source https://github.invalid/c/1; the question: " + work.S01Questions[1].Text + "\n" +
		"3. `S01-visibility` public \u2014 by operator; source https://github.invalid/c/1; the question: " + work.S01Questions[2].Text + "\n" +
		"4. `S01-baseline` https://github.com/pharzam/armature \u2014 by operator; source F-0003#5; the question: " + work.S01Questions[3].Text + "\n" +
		"5. `Q-001` Go &lsaquo;1.26&rsaquo; \u2014 by idea-owner; source said &lsaquo;here&rsaquo;; the question: Which stack? &lsaquo;x&rsaquo;\n" +
		"\n## Notes on capture\n\nEach angle quote of a recorded text is written as `&lsaquo;` or `&rsaquo;`.\n"
	if text != want {
		t.Errorf("answersRecord =\n%s\nwant\n%s", text, want)
	}
}

// fakeRepo is a fake of internal/git and of the files of the work area w, for
// the unit tests of S02 to S04 (step 3 of the plan of #86): it logs each call,
// keeps each file written, and fails the first call that starts with fail.
type fakeRepo struct {
	calls  []string
	files  map[string]string
	exists map[string]bool
	revs   map[string]string          // a revision and its object name
	roots  []string                   // the root commits of main
	branch string                     // the branch of the target
	msg    string                     // the message of each commit
	trees  map[string][]git.TreeEntry // a directory of the root commit and its files
	shows  map[string]string          // a file of the root commit and its text
	fail   string
}

func (f *fakeRepo) install(t *testing.T) {
	t.Helper()
	saved := sys
	t.Cleanup(func() { sys = saved })
	f.files = map[string]string{}
	call := func(s string) error {
		f.calls = append(f.calls, s)
		if f.fail != "" && strings.HasPrefix(s, f.fail) {
			return errors.New(s + ": failed\nthe second line")
		}
		return nil
	}
	rev := func(r string) (string, error) {
		if err := call("rev-parse " + r); err != nil {
			return "", err
		}
		if v, ok := f.revs[r]; ok {
			return v, nil
		}
		return "", errors.New("unknown revision " + r)
	}
	sys = system{
		now:       func() time.Time { return time.Date(2026, 10, 2, 9, 30, 0, 0, time.UTC) },
		exists:    func(p string) bool { return f.exists[p] },
		removeAll: func(p string) error { return call("remove " + p) },
		lsRemote: func(url, ref string) (string, error) {
			return f.revs["ls-remote"], call("ls-remote " + url + " " + ref)
		},
		clone:          func(url, dir string) error { return call("clone " + url + " " + dir) },
		checkoutDetach: func(dir, c string) error { return call("checkout " + dir + " " + c) },
		revParse:       func(_, r string) (string, error) { return rev(r) },
		head:           func(string) (string, error) { return rev("HEAD") },
		rootCommits:    func(string, string) ([]string, error) { return f.roots, call("roots") },
		message:        func(_, r string) (string, error) { return f.msg, call("message " + r) },
		initRepo:       func(dir string) error { return call("init " + dir) },
		commit: func(dir, msg string, who git.Identity) error {
			return call(fmt.Sprintf("commit %s %q by %s at %s", dir, msg, who.Name, who.Time.UTC().Format(time.RFC3339)))
		},
		rename:       func(a, b string) error { return call("rename " + a + " " + b) },
		branch:       func(string) (string, error) { return f.branch, call("branch") },
		switchCreate: func(_, b, c string) error { return call("switch -c " + b + " " + c) },
		lsTree:       func(_, r, p string) ([]git.TreeEntry, error) { return f.trees[p], call("ls-tree " + r + " " + p) },
		show: func(_, r, p string) ([]byte, error) {
			if err := call("show " + r + ":" + p); err != nil {
				return nil, err
			}
			if v, ok := f.shows[p]; ok {
				return []byte(v), nil
			}
			return nil, errors.New("no file " + p)
		},
		write: func(p string, data []byte) error { f.files[p] = string(data); return call("write " + p) },
	}
}

const (
	baseURL = "https://github.com/pharzam/armature"
	commitA = "a95965534b14b0bf14ad74da0c9a45b5f4aedf88"
	treeA   = "8ffb250afd584da8b418bc220fb6d72e802924ce"
	rootA   = "1111111111111111111111111111111111111111"
)

// pinRecord gives a record whose S01 and S02 are done, with the pin rows.
func pinRecord() work.Record {
	return work.Record{
		{"S01", "name", "acme/widget", "answer", "S01-name"},
		{"S01", "baseline", baseURL, "fact", "F-0003#5"},
		{"S01", "done", "x", "step", ""},
		{"S02", "pin.source", baseURL, "fact", "F-0003#5"},
		{"S02", "pin.commit", commitA, "computed", "git ls-remote " + baseURL + " HEAD"},
		{"S02", "pin.tree", treeA, "computed", "git rev-parse " + commitA + "^{tree}"},
		{"S02", "pin.time", "2026-09-30T08:00:00Z", "computed", "the clock of the LAYUP host"},
		{"S02", "done", "the commit and the tree", "step", ""},
	}
}

// S02 resolves the commit once, at the time of the clock, clones into
// target.part, checks the commit out, removes .git and renames the directory
// at its end (D7, D8 of #86); a target that exists is an input error, and a
// call that fails is a fail with its first line.
func TestS02(t *testing.T) {
	f := &fakeRepo{revs: map[string]string{"ls-remote": commitA, commitA + "^{tree}": treeA}}
	f.install(t)
	o := runS02(Input{Dir: "w", Record: pinRecord()[:3]})
	want := [][]string{
		{"pin.source", baseURL, "fact", "F-0003#5"},
		{"pin.commit", commitA, "computed", "git ls-remote " + baseURL + " HEAD"},
		{"pin.tree", treeA, "computed", "git rev-parse " + commitA + "^{tree}"},
		{"pin.time", "2026-10-02T09:30:00Z", "computed", "the clock of the LAYUP host"},
	}
	calls := []string{"remove w/target.part", "ls-remote " + baseURL + " HEAD", "clone " + baseURL + " w/target.part", "checkout w/target.part " + commitA,
		"rev-parse " + commitA + "^{tree}", "remove w/target.part/.git", "rename w/target.part w/target"}
	if o.Kind != Done || o.Evidence != "the commit and the tree" || !reflect.DeepEqual(o.Values, want) || !slices.Equal(f.calls, calls) {
		t.Errorf("S02: %s %q %q, calls %q; want done, %q and %q", o.Kind, o.Evidence, o.Values, f.calls, want, calls)
	}
	for _, c := range []struct{ fail, kind, evidence string }{
		{"ls-remote", Fail, "git ls-remote " + baseURL + " HEAD: ls-remote " + baseURL + " HEAD: failed"},
		{"clone", Fail, "git clone " + baseURL + ": clone " + baseURL + " w/target.part: failed"},
		{"checkout", Fail, "git checkout " + commitA + ": checkout w/target.part " + commitA + ": failed"},
		{"rename", Fail, "the clone: rename w/target.part w/target: failed"},
	} {
		f := &fakeRepo{revs: map[string]string{"ls-remote": commitA, commitA + "^{tree}": treeA}, fail: c.fail}
		f.install(t)
		if o := runS02(Input{Dir: "w", Record: pinRecord()[:3]}); o.Kind != c.kind || o.Evidence != c.evidence || o.Values != nil {
			t.Errorf("%s fails: %s %q; want %s %q", c.fail, o.Kind, o.Evidence, c.kind, c.evidence)
		}
	}
	f = &fakeRepo{exists: map[string]bool{"w/target": true}}
	f.install(t)
	if o := runS02(Input{Dir: "w", Record: pinRecord()[:3]}); o.Kind != Invalid || !strings.HasPrefix(o.Evidence, "target of the work area exists, and S02 is not done") || f.calls != nil {
		t.Errorf("a target that exists: %s %q, calls %q; want an input error and no call", o.Kind, o.Evidence, f.calls)
	}
}

// S03 makes the root commit of the copy on main, by the identity of the run at
// pin.time, and checks its tree; a target whose main is one commit with the
// tree pin.tree is the commit of a run that stopped; another history is an
// input error (D8 of #86).
func TestS03(t *testing.T) {
	f := &fakeRepo{revs: map[string]string{"HEAD^{tree}": treeA}}
	f.install(t)
	o := runS03(Input{Dir: "w", Record: pinRecord(), Who: Who})
	calls := []string{"init w/target", `commit w/target "chore: the unmodified baseline at ` + commitA + `" by layup-agent[bot] at 2026-09-30T08:00:00Z`, "rev-parse HEAD^{tree}"}
	if o.Kind != Done || o.Evidence != "root tree = pin.tree" || o.Values != nil || !slices.Equal(f.calls, calls) {
		t.Errorf("S03: %s %q, calls %q; want done and %q", o.Kind, o.Evidence, f.calls, calls)
	}
	f = &fakeRepo{revs: map[string]string{"HEAD^{tree}": "2222222222222222222222222222222222222222"}}
	f.install(t)
	if o := runS03(Input{Dir: "w", Record: pinRecord(), Who: Who}); o.Kind != Fail || o.Evidence != "the root tree 2222222222222222222222222222222222222222 is not pin.tree "+treeA {
		t.Errorf("another tree: %s %q; want fail", o.Kind, o.Evidence)
	}
	own := "chore: the unmodified baseline at " + commitA
	for _, c := range []struct {
		name      string
		roots     []string
		tree, msg string
		kind      string
	}{
		{"its own commit", []string{rootA}, treeA, own, Done},
		{"another tree", []string{rootA}, "2222222222222222222222222222222222222222", own, Invalid},
		{"two root commits", []string{rootA, "3333333333333333333333333333333333333333"}, treeA, own, Invalid},
		{"another message", []string{rootA}, treeA, own + "\n\nby hand", Invalid},
	} {
		f := &fakeRepo{exists: map[string]bool{"w/target/.git": true}, roots: c.roots, msg: c.msg, revs: map[string]string{"refs/heads/main^{commit}": rootA, rootA + "^{tree}": c.tree}}
		f.install(t)
		o := runS03(Input{Dir: "w", Record: pinRecord(), Who: Who})
		if o.Kind != c.kind || slices.ContainsFunc(f.calls, func(s string) bool { return strings.HasPrefix(s, "init") || strings.HasPrefix(s, "commit") }) {
			t.Errorf("a target with a history, %s: %s %q, calls %q; want %s and no new commit", c.name, o.Kind, o.Evidence, f.calls, c.kind)
		}
	}
	bad := pinRecord()
	bad[6][2] = "2026-10-02"
	if o := runS03(Input{Dir: "w", Record: bad, Who: Who}); o.Kind != Fail || !strings.Contains(o.Evidence, "pin.time") {
		t.Errorf("a pin.time of another form: %s %q; want fail", o.Kind, o.Evidence)
	}
}

// The commands of S03 (D9 of #86): the remote origin of the target on GitHub
// and the push of main, each with its comment; the comment says whether the
// commit is LAYUP's own pin.
func TestTheCommandsOfS03(t *testing.T) {
	r := pinRecord()
	want := []Command{
		{1, "the remote of the target: its empty repository on GitHub", "git -C target remote add origin 'https://github.com/acme/widget.git'"},
		{1, "the push of the root commit: the unmodified baseline at " + commitA + ", the commit of LAYUP's own pin a959655", "git -C target push origin main"},
	}
	if got := commandsS03(r); !reflect.DeepEqual(got, want) {
		t.Fatalf("commandsS03 = %+v; want %+v", got, want)
	}
	r[4][2] = "4444444444444444444444444444444444444444"
	if got := commandsS03(r); !strings.HasSuffix(got[1].Comment, ", not the commit of LAYUP's own pin a959655") {
		t.Errorf("another commit: %q", got[1].Comment)
	}
}

// adrIndex is the index of docs/adr/ at LAYUP's pin, cut to two rows.
const adrIndex = "# ADRs\n\n## Index\n\n| ADR | Title | Status |\n| --- | ----- | ------ |\n| [0001](0001-a.md) | A | Accepted |\n| [0008](0008-h.md) | H | Accepted |\n\n<!-- Add one row per ADR. -->\n"

// factsIndex is the index of docs/facts/ at LAYUP's pin.
const factsIndex = "# Facts\n\n## Index\n\n| Fact doc | Source | Collected | Status |\n| -------- | ------ | --------- | ------ |\n| _none yet_ | | | |\n"

// s04Repo gives a fake of a target on main at its root commit, with the files
// of LAYUP's pin that S04 reads.
func s04Repo(branch, head string) *fakeRepo {
	return &fakeRepo{branch: branch, revs: map[string]string{"refs/heads/main^{commit}": rootA, "HEAD": head},
		trees: map[string][]git.TreeEntry{"docs/adr": {{Path: "docs/adr/0001-a.md"}, {Path: "docs/adr/0008-h.md"}, {Path: "docs/adr/README.md"}, {Path: "docs/adr/tests/0009-x.md"}},
			"docs/facts": {{Path: "docs/facts/README.md"}, {Path: "docs/facts/template.md"}}},
		shows: map[string]string{"docs/adr/README.md": adrIndex, "docs/facts/README.md": factsIndex}}
}

// S04 writes, on the branch layup-setup from the root commit, the pin file,
// the decision record of the pin with its index row, and the answers record
// with its index row and its line of facts.sha256, each from the files of the
// root commit and the date of pin.time (D10 of #86, conditions 1 and 2 of its
// plan review).
func TestS04(t *testing.T) {
	f := s04Repo("refs/heads/main", rootA)
	f.install(t)
	gaps := Brief{Gaps: []byte(table)}
	o := runS04(gaps, Input{Dir: "w", Record: pinRecord(), Answers: ans(nil)})
	record := answersRecord("F-0001", "2026-09-30", append(slices.Clone(work.S01Questions),
		work.Question{ID: "Q-001", Text: "Which technology stack does the product use?"}, work.Question{ID: "Q-002", Text: "What does fast mean here, as a number?"}), ans(nil))
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
	if o.Kind != Done || o.Evidence != "checks pin and facts" || !o.Commit || !reflect.DeepEqual(o.Values, values) || !slices.Contains(f.calls, "switch -c layup-setup "+rootA) {
		t.Errorf("S04: %s %q %v %q, calls %q; want done, the values %q and the branch made", o.Kind, o.Evidence, o.Commit, o.Values, f.calls, values)
	}
	if !reflect.DeepEqual(f.files, files) {
		for p, text := range files {
			if f.files[p] != text {
				t.Errorf("S04 wrote %s:\n%s\nwant\n%s", p, f.files[p], text)
			}
		}
		t.Errorf("S04 wrote %d files; want %d", len(f.files), len(files))
	}
	// A facts.sha256 of the root commit keeps its lines, also a last line with
	// no line feed.
	for _, list := range []string{"abc  docs/facts/x.md\n", "abc  docs/facts/x.md"} {
		f = s04Repo("refs/heads/main", rootA)
		f.shows["docs/setup/facts.sha256"] = list
		f.install(t)
		if o := runS04(gaps, Input{Dir: "w", Record: pinRecord(), Answers: ans(nil)}); o.Kind != Done || f.files["w/target/docs/setup/facts.sha256"] != "abc  docs/facts/x.md\n"+sum+"  docs/facts/F-0001-setup-answers.md\n" {
			t.Errorf("a list of hashes in the root commit, %q: %s %q", list, o.Kind, f.files["w/target/docs/setup/facts.sha256"])
		}
	}
}

// S04 takes a branch layup-setup at the root commit (a stop before its commit,
// or the undo of its evidence); a branch at another commit, or a target on
// another branch, is an input error; a missing pin row is a fail (condition 1
// of the plan review of #86).
func TestS04OnABranchThatExists(t *testing.T) {
	for _, c := range []struct {
		name, branch, head string
		exists             bool
		kind               string
	}{
		{"layup-setup at the root commit", "refs/heads/layup-setup", rootA, true, Done},
		{"layup-setup at another commit", "refs/heads/layup-setup", "5555555555555555555555555555555555555555", true, Invalid},
		{"main, and layup-setup exists", "refs/heads/main", rootA, true, Invalid},
		{"another branch", "refs/heads/x", rootA, false, Invalid},
	} {
		f := s04Repo(c.branch, c.head)
		if c.exists {
			f.revs["refs/heads/layup-setup^{commit}"] = c.head
		}
		f.install(t)
		o := runS04(Brief{Gaps: []byte(table)}, Input{Dir: "w", Record: pinRecord(), Answers: ans(nil)})
		if o.Kind != c.kind || slices.ContainsFunc(f.calls, func(s string) bool { return strings.HasPrefix(s, "switch") }) || (c.kind != Done && len(f.files) != 0) ||
			(c.branch == "refs/heads/main" && o.Evidence != "target of the work area is on main, and its branch layup-setup exists: put it on layup-setup, or start again in a new work area") {
			t.Errorf("%s: %s %q, calls %q, %d files; want %s, no new branch", c.name, o.Kind, o.Evidence, f.calls, len(f.files), c.kind)
		}
	}
	f := s04Repo("refs/heads/main", rootA)
	f.install(t)
	if o := runS04(Brief{Gaps: []byte(table)}, Input{Dir: "w", Record: pinRecord()[:6], Answers: ans(nil)}); o.Kind != Fail || o.Evidence != "the record has no row S02 pin.time" {
		t.Errorf("no row pin.time: %s %q", o.Kind, o.Evidence)
	}
}
