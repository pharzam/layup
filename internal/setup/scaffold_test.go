package setup

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// sum gives the SHA-256 of text in hexadecimal.
func sum(text string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(text))) }

// The history of a stand-in baseline at the head of layup-setup, for S05.
const (
	kitBacklog = "# Backlog\n\n- **T-0001** \u2014 a task ([detail](T-0001.md))\n- **T-9999** \u2014 a task of the target\n- **T-00010** \u2014 another\n" +
		"\n> a note of the kit\n> on two lines ([#2](https://github.com/pharzam/armature/issues/2))\n\nkeep this line\n"
	kitCompleted = "# Completed\n\n- **2026-01-01** \u2014 **T-0002** \u2014 done ([#1](https://github.com/pharzam/armature/issues/1))\n" +
		"- **2026-01-02** \u2014 **T-0003** \u2014 linked ([#3](https://github.com/pharzam/armature/issues/3))\n"
)

// s05Repo gives a fake of a target whose head holds the history of the kit.
func s05Repo(inputs map[string]string) *fakeRepo {
	return &fakeRepo{
		trees: map[string][]git.TreeEntry{
			"docs/tasks":     {{Path: "docs/tasks/T-0001.md"}, {Path: "docs/tasks/T-0002.md"}, {Path: "docs/tasks/backlog.md"}, {Path: "docs/tasks/completed.md"}},
			"docs/decisions": {{Path: "docs/decisions/D-0001.md"}},
			"docs/audit":     {{Path: "docs/audit/README.md"}},
		},
		shows: map[string]string{"docs/tasks/backlog.md": kitBacklog, "docs/tasks/completed.md": kitCompleted},
		disk:  inputs,
	}
}

// s05Calls gives the calls of S05: the files whose links break, one list per
// call, and the link rule of the baseline.
func s05Calls(broken ...[]string) Calls {
	return Calls{
		BrokenLinks: func(string) ([]string, error) {
			if len(broken) == 0 {
				return nil, nil
			}
			r := broken[0]
			broken = broken[1:]
			return r, nil
		},
		LinksBaseline: func(source, line string) bool {
			return source == baseURL && strings.Contains(line, "github.com/pharzam/armature")
		},
	}
}

// S05 deletes the history of the baseline from its head: the two directories,
// each task file, each line of the two indexes that names a deleted task or
// links the baseline's repository, and each blockquote that holds such a
// line; a file whose links break gets its input file, with a record row
// (D1 to D3 of #90).
func TestS05(t *testing.T) {
	f := s05Repo(map[string]string{"w/inputs/files/docs/x.md": "fixed\n"})
	f.install(t)
	o := runS05(s05Calls([]string{"docs/x.md"}, nil), Input{Dir: "w", Record: pinRecord()})
	removed := []string{"remove w/target/docs/decisions", "remove w/target/docs/audit", "remove w/target/docs/tasks/T-0001.md", "remove w/target/docs/tasks/T-0002.md"}
	files := map[string]string{
		"w/target/docs/tasks/backlog.md":   "# Backlog\n\n- **T-9999** \u2014 a task of the target\n- **T-00010** \u2014 another\n\n\nkeep this line\n",
		"w/target/docs/tasks/completed.md": "# Completed\n\n",
		"w/target/docs/x.md":               "fixed\n",
	}
	values := [][]string{{"file:docs/x.md", sum("fixed\n"), "computed", "sha256 inputs/files/docs/x.md"}}
	var got []string
	for _, c := range f.calls {
		if strings.HasPrefix(c, "remove ") {
			got = append(got, c)
		}
	}
	if o.Kind != Done || o.Evidence != "checks kit-history and link-lint" || !o.Commit || !reflect.DeepEqual(o.Values, values) || !slices.Equal(got, removed) || !reflect.DeepEqual(f.files, files) {
		t.Errorf("S05: %s %q %v %q, removed %q, files %q; want done, %q, %q and %q", o.Kind, o.Evidence, o.Commit, o.Values, got, f.files, values, removed, files)
	}
}

// A file whose links break and that has no input file is one stop row F-<path>
// of S05; an input that still breaks, or a failed list, is a fail (D2 of #90).
func TestS05StopsForTheLinksThatBreak(t *testing.T) {
	f := s05Repo(nil)
	f.install(t)
	o := runS05(s05Calls([]string{"docs/x.md", "docs/y.md"}), Input{Dir: "w", Record: pinRecord()})
	want := []StopRow{
		{"S05", "F-docs/x.md", "Fix the links of docs/x.md that the deletion of the baseline's history breaks, and give the file as inputs/files/docs/x.md.", "docs/x.md"},
		{"S05", "F-docs/y.md", "Fix the links of docs/y.md that the deletion of the baseline's history breaks, and give the file as inputs/files/docs/y.md.", "docs/y.md"},
	}
	if o.Kind != Stop || !reflect.DeepEqual(o.Stops, want) || f.files["w/target/docs/x.md"] != "" {
		t.Errorf("no input: %s %+v; want a stop with %+v and no copy", o.Kind, o.Stops, want)
	}
	f = s05Repo(map[string]string{"w/inputs/files/docs/x.md": "still broken\n"})
	f.install(t)
	if o := runS05(s05Calls([]string{"docs/x.md"}, []string{"docs/x.md"}), Input{Dir: "w", Record: pinRecord()}); o.Kind != Fail ||
		o.Evidence != "the links of these files still break after their input files: docs/x.md" {
		t.Errorf("an input that still breaks: %s %q", o.Kind, o.Evidence)
	}
	c := s05Calls()
	c.BrokenLinks = func(string) ([]string, error) { return nil, errors.New("link-lint.sh: exit 2 with no file\nmore") }
	f = s05Repo(nil)
	f.install(t)
	if o := runS05(c, Input{Dir: "w", Record: pinRecord()}); o.Kind != Fail || o.Evidence != "the links of the tree: link-lint.sh: exit 2 with no file" {
		t.Errorf("a failed list: %s %q", o.Kind, o.Evidence)
	}
}

// factsIndexS04 is the index of docs/facts/ after S04.
const factsIndexS04 = "# Facts\n\n## Index\n\n| Fact doc | Source | Collected | Status |\n| -------- | ------ | --------- | ------ |\n" +
	"| [F-0001](F-0001-setup-answers.md) | The answers to the questions of the setup | 2026-09-30 | Raw |\n"

// S06 copies each brief byte for byte, with its line in facts.sha256, its
// index row and its two record rows (D4 of #90, note 3 of its plan review).
func TestS06(t *testing.T) {
	const list = "abc  docs/facts/F-0001-setup-answers.md\n"
	for _, vision := range []bool{false, true} {
		f := &fakeRepo{shows: map[string]string{factsSumsPath: list, factsDir + "/README.md": factsIndexS04},
			disk: map[string]string{"w/inputs/briefs/problem-statement.md": "brief\n"}}
		files := map[string]string{
			"w/target/docs/facts/problem-statement-brief.md": "brief\n",
			"w/target/" + factsSumsPath:                      list + sum("brief\n") + "  docs/facts/problem-statement-brief.md\n",
			"w/target/docs/facts/README.md": factsIndexS04 +
				"| [problem-statement-brief.md](problem-statement-brief.md) | The problem statement of the idea owner, from the work area | 2026-09-30 | Raw |\n",
		}
		values := [][]string{
			{"brief.copy", "docs/facts/problem-statement-brief.md", "computed", "the raw file name of a brief"},
			{"brief.copy.sha256", sum("brief\n"), "computed", "sha256 docs/facts/problem-statement-brief.md"},
		}
		if vision {
			f.disk["w/inputs/briefs/vision.md"] = "vision\n"
			files["w/target/docs/facts/architectural-vision-brief.md"] = "vision\n"
			files["w/target/"+factsSumsPath] += sum("vision\n") + "  docs/facts/architectural-vision-brief.md\n"
			files["w/target/docs/facts/README.md"] += "| [architectural-vision-brief.md](architectural-vision-brief.md) | The vision brief of the idea owner, from the work area | 2026-09-30 | Raw |\n"
			values = append(values, []string{"vision.copy", "docs/facts/architectural-vision-brief.md", "computed", "the raw file name of a brief"},
				[]string{"vision.copy.sha256", sum("vision\n"), "computed", "sha256 docs/facts/architectural-vision-brief.md"})
		}
		f.install(t)
		o := runS06(Input{Dir: "w", Record: pinRecord()})
		if o.Kind != Done || o.Evidence != "check facts" || !o.Commit || !reflect.DeepEqual(o.Values, values) || !reflect.DeepEqual(f.files, files) {
			t.Errorf("S06 (vision %v): %s %q %q, files %q; want done, %q and %q", vision, o.Kind, o.Evidence, o.Values, f.files, values, files)
		}
	}
}

// Each step of the prose step copies its named files from inputs/files/, with
// a record row for each (K42); a missing input is one stop row F-<path> of the
// step (D5 of #90).
func TestTheProseStep(t *testing.T) {
	f := &fakeRepo{disk: map[string]string{"w/inputs/files/docs/glossary.md": "# Glossary\n"}}
	f.install(t)
	o := runProse("S08", Input{Dir: "w"}, []string{glossaryPath})
	want := [][]string{{"file:docs/glossary.md", sum("# Glossary\n"), "computed", "sha256 inputs/files/docs/glossary.md"}}
	if o.Kind != Done || o.Evidence != "check glossary" || !o.Commit || !reflect.DeepEqual(o.Values, want) || f.files["w/target/docs/glossary.md"] != "# Glossary\n" {
		t.Errorf("S08: %s %q %q, files %q", o.Kind, o.Evidence, o.Values, f.files)
	}
	f = &fakeRepo{}
	f.install(t)
	o = runProse("S07", Input{Dir: "w"}, []string{onboardingPath})
	stop := []StopRow{{"S07", "F-docs/onboarding-for-engineers.md", "Write the text of docs/onboarding-for-engineers.md for the target, and give it as inputs/files/docs/onboarding-for-engineers.md.", "docs/onboarding-for-engineers.md"}}
	if o.Kind != Stop || !reflect.DeepEqual(o.Stops, stop) || len(f.files) != 0 {
		t.Errorf("S07 with no input: %s %+v, files %q; want %+v", o.Kind, o.Stops, f.files, stop)
	}
}

// S14 copies README.md, AGENTS.md and each other file that check adapted
// flags, less the files of S07 to S09, so no question is in the group's table
// twice (D5 of #90, condition 2 of its plan review).
func TestS14(t *testing.T) {
	flagged := []string{"AGENTS.md", "README.md", "docs/engineering-discipline.md", glossaryPath, guardrailsPath, onboardingPath}
	c := Calls{Flagged: func(string) ([]string, error) { return flagged, nil }}
	f := &fakeRepo{disk: map[string]string{"w/inputs/files/README.md": "r\n", "w/inputs/files/AGENTS.md": "a\n"}}
	f.install(t)
	o := runS14(c, Input{Dir: "w"})
	stop := []StopRow{{"S14", "F-docs/engineering-discipline.md", "Adapt docs/engineering-discipline.md, which check adapted flags, to the target, and give it as inputs/files/docs/engineering-discipline.md.", "docs/engineering-discipline.md"}}
	if o.Kind != Stop || !reflect.DeepEqual(o.Stops, stop) || len(f.files) != 0 {
		t.Errorf("S14 with no input for a flagged file: %s %+v, files %q; want %+v", o.Kind, o.Stops, f.files, stop)
	}
	f.disk["w/inputs/files/docs/engineering-discipline.md"] = "e\n"
	f.install(t)
	o = runS14(c, Input{Dir: "w"})
	values := [][]string{
		{"file:README.md", sum("r\n"), "computed", "sha256 inputs/files/README.md"},
		{"file:AGENTS.md", sum("a\n"), "computed", "sha256 inputs/files/AGENTS.md"},
		{"file:docs/engineering-discipline.md", sum("e\n"), "computed", "sha256 inputs/files/docs/engineering-discipline.md"},
	}
	if o.Kind != Done || o.Evidence != "checks identity and adapted" || !reflect.DeepEqual(o.Values, values) || len(f.files) != 3 {
		t.Errorf("S14: %s %q %q, files %q; want done and %q", o.Kind, o.Evidence, o.Values, f.files, values)
	}
	c.Flagged = func(string) ([]string, error) { return nil, errors.New("git: cannot list the tracked files") }
	if o := runS14(c, Input{Dir: "w"}); o.Kind != Fail || o.Evidence != "the files that check adapted flags: git: cannot list the tracked files" {
		t.Errorf("a failed list: %s %q", o.Kind, o.Evidence)
	}
}

var (
	mx = "\u2039x\u203a"
	my = "\u2039y\u203a"
)

// markersOf gives the call of the markers of a tree, with the markers marks.
func markersOf(marks ...Marker) Calls {
	return Calls{Markers: func(string) ([]Marker, error) { return marks, nil }}
}

// mAnswers gives answers.tsv with the M- rows: an ID, its answer and its
// question_text.
func mAnswers(rows ...[]string) work.Answers {
	var a work.Answers
	for _, r := range rows {
		a = append(a, []string{r[0], r[1], "operator", r[2], r[3]})
	}
	return a
}

// S10 asks one question per file and marker, at the first line of the marker,
// all in one table; an M- answer to a marker that the tree does not hold is an
// input error (D6 of #90).
func TestS10(t *testing.T) {
	c := markersOf(Marker{"docs/a.md", 3, 0, mx}, Marker{"docs/a.md", 7, 2, mx}, Marker{"docs/b.md", 1, 4, my})
	idx, idy := MarkerID("docs/a.md", mx), MarkerID("docs/b.md", my)
	ask := func(file, m string) string {
		return "What is the value of " + m + " in " + file + "? Answer gap to keep it as an open gap, with its question as question_text."
	}
	o := runS10(c, Input{Dir: "w"})
	want := []StopRow{{"S10", idx, ask("docs/a.md", mx), "docs/a.md:3 " + mx}, {"S10", idy, ask("docs/b.md", my), "docs/b.md:1 " + my}}
	if o.Kind != Stop || !reflect.DeepEqual(o.Stops, want) {
		t.Errorf("no answer: %s %+v; want %+v", o.Kind, o.Stops, want)
	}
	if o := runS10(c, Input{Dir: "w", Answers: mAnswers([]string{idx, "8080", "u", ""})}); o.Kind != Stop || !reflect.DeepEqual(o.Stops, want[1:]) {
		t.Errorf("one answer: %s %+v; want %+v", o.Kind, o.Stops, want[1:])
	}
	a := mAnswers([]string{idx, "gap", "u", "Which x?"}, []string{idy, "y", "u", ""})
	if o := runS10(c, Input{Dir: "w", Answers: a}); o.Kind != Done || o.Evidence != "every marker has an answer row" || o.Commit {
		t.Errorf("each answer: %s %q %v; want done and no commit", o.Kind, o.Evidence, o.Commit)
	}
	if o := runS10(c, Input{Dir: "w", Answers: append(a, mAnswers([]string{"M-00000000", "z", "u", ""})...)}); o.Kind != Invalid || !strings.Contains(o.Evidence, "M-00000000") {
		t.Errorf("an answer to no marker: %s %q; want an input error", o.Kind, o.Evidence)
	}
	if o := runS10(markersOf(), Input{Dir: "w"}); o.Kind != Done {
		t.Errorf("no marker: %s %q; want done", o.Kind, o.Evidence)
	}
	c.Markers = func(string) ([]Marker, error) { return nil, errors.New("git: cannot list") }
	if o := runS10(c, Input{Dir: "w"}); o.Kind != Fail || o.Evidence != "the markers of the tree: git: cannot list" {
		t.Errorf("a failed list: %s %q", o.Kind, o.Evidence)
	}
}

// s11Repo gives a fake of a target with two files that hold markers, after
// S04 and S06.
func s11Repo() *fakeRepo {
	return &fakeRepo{
		trees: map[string][]git.TreeEntry{factsDir: {{Path: "docs/facts/F-0001-setup-answers.md"}, {Path: "docs/facts/README.md"}, {Path: "docs/facts/problem-statement-brief.md"}}},
		shows: map[string]string{factsSumsPath: "abc  docs/facts/F-0001-setup-answers.md\n", factsDir + "/README.md": factsIndexS04},
		disk: map[string]string{
			"w/target/docs/a.md": "one\ntwo\n" + mx + " three\nfour\nfive\nsix\n  " + mx + " seven and the `\u2039` mention\n",
			"w/target/docs/b.md": "see " + my + " here\nand " + my + " there\n",
		},
	}
}

// S11 fills each marker with a value at the places that the scanner found,
// keeps a gap with its row of open-gaps.tsv, writes a record row per place,
// and writes the second answers record (D8 of #90, condition 3 of its plan
// review).
func TestS11(t *testing.T) {
	c := markersOf(Marker{"docs/a.md", 3, 0, mx}, Marker{"docs/a.md", 7, 2, mx}, Marker{"docs/b.md", 1, 4, my}, Marker{"docs/b.md", 2, 4, my})
	idx, idy := MarkerID("docs/a.md", mx), MarkerID("docs/b.md", my)
	a := mAnswers([]string{idx, "8080", "F-0003#5", ""}, []string{idy, "gap", "https://github.invalid/c/1", "Which y?"})
	f := s11Repo()
	f.install(t)
	o := runS11(c, Input{Dir: "w", Record: pinRecord(), Answers: a})
	asked := []work.Question{{ID: idx, Text: markerAsk("docs/a.md", mx)}, {ID: idy, Text: markerAsk("docs/b.md", my)}}
	record := answersRecord("F-0002", "2026-09-30", markerRecord, asked, a)
	files := map[string]string{
		"w/target/docs/a.md":                           "one\ntwo\n8080 three\nfour\nfive\nsix\n  8080 seven and the `\u2039` mention\n",
		"w/target/" + work.OpenGapsPath:                "docs/b.md\t" + my + "\tWhich y?\n",
		"w/target/docs/facts/F-0002-marker-answers.md": record,
		"w/target/docs/facts/README.md":                factsIndexS04 + "| [F-0002](F-0002-marker-answers.md) | The answers to the markers of the setup | 2026-09-30 | Raw |\n",
		"w/target/" + factsSumsPath:                    "abc  docs/facts/F-0001-setup-answers.md\n" + sum(record) + "  docs/facts/F-0002-marker-answers.md\n",
	}
	values := [][]string{
		{"marker:docs/a.md:3", "8080", "fact", "F-0003#5"},
		{"marker:docs/a.md:7", "8080", "fact", "F-0003#5"},
		{"marker:docs/b.md:1", my, "gap", work.OpenGapsPath},
		{"marker:docs/b.md:2", my, "gap", work.OpenGapsPath},
		{"marker.record", "docs/facts/F-0002-marker-answers.md", "computed", "the next free ID of docs/facts/"},
		{"marker.record.sha256", sum(record), "computed", "sha256 docs/facts/F-0002-marker-answers.md"},
	}
	if o.Kind != Done || o.Evidence != "checks markers, sources and facts" || !o.Commit || !reflect.DeepEqual(o.Values, values) {
		t.Errorf("S11: %s %q %v\n%q\nwant done and\n%q", o.Kind, o.Evidence, o.Commit, o.Values, values)
	}
	for p, text := range files {
		if f.files[p] != text {
			t.Errorf("S11 wrote %s:\n%q\nwant\n%q", p, f.files[p], text)
		}
	}
	if len(f.files) != len(files) {
		t.Errorf("S11 wrote %d files; want %d", len(f.files), len(files))
	}
	if !strings.Contains(record, "| Collected by | `layup setup`, step S11 |") || !strings.Contains(record, "# F-0002. The answers to the markers of the setup") ||
		!strings.Contains(record, "| Source | The answers of `inputs/answers.tsv` to the markers that S10 listed |") || strings.Contains(record, "\u2039") {
		t.Errorf("the second answers record:\n%s", record)
	}
}

// S11 checks each answer and each marker before it changes a file: a value
// with an angle quote, a gap marker in a file whose name holds a tab, or a
// marker with no answer is a fail with no file written; with no marker, S11
// writes nothing (D8 of #90, condition 1 of its plan review).
func TestS11ChecksBeforeItWrites(t *testing.T) {
	ida := MarkerID("docs/a.md", mx)
	tab := "docs/t\tb.md"
	for _, c := range []struct {
		name   string
		marks  []Marker
		a      work.Answers
		reason string
	}{
		{"a value with an angle quote", []Marker{{"docs/a.md", 3, 0, mx}}, mAnswers([]string{ida, "\u2039v\u203a", "u", ""}), ida + ": the value holds an angle quote, so it would be a marker"},
		{"a gap in a file whose name holds a tab", []Marker{{"docs/a.md", 3, 0, mx}, {tab, 1, 0, mx}},
			mAnswers([]string{ida, "8080", "u", ""}, []string{MarkerID(tab, mx), "gap", "u", "Which?"}), "docs/t\tb.md: a file whose name holds a tab cannot have a row in docs/setup/open-gaps.tsv"},
		{"a marker with no answer", []Marker{{"docs/a.md", 3, 0, mx}}, nil, ida + ": no answer"},
	} {
		f := s11Repo()
		f.install(t)
		if o := runS11(markersOf(c.marks...), Input{Dir: "w", Record: pinRecord(), Answers: c.a}); o.Kind != Fail || o.Evidence != c.reason || len(f.files) != 0 {
			t.Errorf("%s: %s %q, %d files; want fail %q and no file", c.name, o.Kind, o.Evidence, len(f.files), c.reason)
		}
	}
	f := s11Repo()
	f.install(t)
	if o := runS11(markersOf(), Input{Dir: "w", Record: pinRecord()}); o.Kind != Done || o.Values != nil || len(f.files) != 0 {
		t.Errorf("no marker: %s %q, files %q; want done with nothing written", o.Kind, o.Values, f.files)
	}
}

// Before each step from S04 to S14 the runner puts the target back to its
// head, so a leftover of a stop or of an undo never enters a commit; a reset
// that fails is a fail of the step (D11 of #90).
func TestTheRunnerResetsTheTargetBeforeAStep(t *testing.T) {
	f := &fakeSys{record: doneRows(3)}
	f.install(t)
	res, err := Run("w", steps(f, nil), testWho, noStep)
	want := slices.Repeat([]string{"w/target"}, 11) // S04 to S13 and S14: not S15
	if err != nil || !slices.Equal(f.hards, want) || res.Steps[3].Result != Done {
		t.Errorf("the resets: %v, %q; want one per step from S04 to S14", err, f.hards)
	}
	f = &fakeSys{record: doneRows(3), hardFails: true}
	f.install(t)
	res, err = Run("w", steps(f, nil), testWho, noStep)
	if err != nil || res.Steps[3] != (StepRow{"S04", "layup-setup", Fail, "the reset of the target failed: fatal: reset failed"}) || slices.Contains(f.calls, "S04") {
		t.Errorf("a reset that fails: %v, %+v, calls %q; want S04 fail and not run", err, res.Steps[3], f.calls)
	}
}
