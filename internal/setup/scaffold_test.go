package setup

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
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
	f = s05Repo(map[string]string{"w/inputs/files/docs/x.md": "fixed\n"})
	f.install(t)
	o = runS05(s05Calls([]string{"docs/x.md", "docs/y.md"}), Input{Dir: "w", Record: pinRecord()})
	if _, copied := f.files["w/target/docs/x.md"]; o.Kind != Stop || !reflect.DeepEqual(o.Stops, want[1:]) || copied {
		t.Errorf("one input of two: %s %+v, x.md copied %v; want a stop with %+v and no copy", o.Kind, o.Stops, copied, want[1:])
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

// proseCalls gives the call of the files that check adapted flags, as on the
// baseline at LAYUP's pin: the named files of S07 to S09 too.
func proseCalls() Calls {
	flagged := []string{"AGENTS.md", "README.md", "docs/engineering-discipline.md", glossaryPath, onboardingPath}
	return Calls{Flagged: func(string) ([]string, error) { return flagged, nil },
		LostMarkers: func(string, []byte, []byte) []string { return nil }}
}

// proseInputs gives the input files of the prose step, less those of omit.
func proseInputs(omit ...string) map[string]string {
	in := map[string]string{}
	for _, p := range []string{onboardingPath, glossaryPath, guardrailsPath, "README.md", "AGENTS.md", "docs/engineering-discipline.md"} {
		if !slices.Contains(omit, p) {
			in["w/inputs/files/"+p] = "text of " + p + "\n"
		}
		in["w/target/"+p] = "the file of the head, " + p + "\n"
	}
	return in
}

func fileRows(paths ...string) [][]string {
	var rows [][]string
	for _, p := range paths {
		rows = append(rows, []string{"file:" + p, sum("text of " + p + "\n"), "computed", "sha256 inputs/files/" + p})
	}
	return rows
}

// The prose step is one unit of input (O-123, finding 1 of review round 1 of
// #90): while an input of the group is missing, no step of the group copies a
// file, and each step gives the stop rows of its own missing inputs, so the
// run gives one table; a step whose own inputs exist waits, with no row. S14
// asks for README.md, AGENTS.md and each other file that check adapted flags,
// less the named file of a step of S07 to S09 that is not done, so no
// question is twice in the table (condition 2 of the plan review).
func TestTheProseStepIsOneUnit(t *testing.T) {
	stop := func(step, p string) StopRow {
		return StopRow{step, "F-" + p, "Write the text of " + p + " for the target, and give it as inputs/files/" + p + ".", p}
	}
	adapt := StopRow{"S14", "F-docs/engineering-discipline.md", "Adapt docs/engineering-discipline.md, which check adapted flags, to the target, and give it as inputs/files/docs/engineering-discipline.md.", "docs/engineering-discipline.md"}
	for _, c := range []struct {
		name  string
		omit  []string
		stops map[string][]StopRow
	}{
		{"no input", []string{onboardingPath, glossaryPath, guardrailsPath, "README.md", "AGENTS.md", "docs/engineering-discipline.md"}, map[string][]StopRow{
			"S07": {stop("S07", onboardingPath)}, "S08": {stop("S08", glossaryPath)}, "S09": {stop("S09", guardrailsPath)},
			"S14": {stop("S14", "README.md"), stop("S14", "AGENTS.md"), adapt}}},
		{"each input but one of S07", []string{onboardingPath}, map[string][]StopRow{"S07": {stop("S07", onboardingPath)}, "S08": nil, "S09": nil, "S14": nil}},
		{"each input but one of S14", []string{"docs/engineering-discipline.md"}, map[string][]StopRow{"S07": nil, "S08": nil, "S09": nil, "S14": {adapt}}},
	} {
		for step, want := range c.stops {
			f := &fakeRepo{disk: proseInputs(c.omit...)}
			f.install(t)
			o := runProseStep(step, proseCalls(), Input{Dir: "w"})
			if o.Kind != Stop || !reflect.DeepEqual(o.Stops, want) || len(f.files) != 0 {
				t.Errorf("%s, %s: %s %+v, %d files; want a stop with %+v and no file", c.name, step, o.Kind, o.Stops, len(f.files), want)
			}
		}
	}
	f := &fakeRepo{disk: proseInputs()}
	f.install(t)
	for step, want := range map[string][][]string{"S07": fileRows(onboardingPath), "S08": fileRows(glossaryPath), "S09": fileRows(guardrailsPath),
		"S14": fileRows("README.md", "AGENTS.md", "docs/engineering-discipline.md")} {
		o := runProseStep(step, proseCalls(), Input{Dir: "w"})
		if o.Kind != Done || !reflect.DeepEqual(o.Values, want) || o.Evidence != map[string]string{"S07": "check onboarding", "S08": "check glossary", "S09": "check guardrails", "S14": "checks identity and adapted"}[step] {
			t.Errorf("each input, %s: %s %q %q; want done and %q", step, o.Kind, o.Evidence, o.Values, want)
		}
	}
	// After S07 to S09, a named file that check adapted still flags is S14's
	// too, so a new input for it reaches the tree.
	done := work.Record{{"S07", "done", "x", "step", ""}, {"S08", "done", "x", "step", ""}, {"S09", "done", "x", "step", ""}}
	if o := runProseStep("S14", proseCalls(), Input{Dir: "w", Record: done}); o.Kind != Done ||
		!reflect.DeepEqual(o.Values, fileRows("README.md", "AGENTS.md", "docs/engineering-discipline.md", glossaryPath, onboardingPath)) {
		t.Errorf("S14 after S07 to S09: %s %q", o.Kind, o.Values)
	}
	c := proseCalls()
	c.Flagged = func(string) ([]string, error) { return nil, errors.New("git: cannot list the tracked files") }
	if o := runProseStep("S08", c, Input{Dir: "w"}); o.Kind != Fail || o.Evidence != "the files that check adapted flags: git: cannot list the tracked files" {
		t.Errorf("a failed list: %s %q", o.Kind, o.Evidence)
	}
}

// S14 refuses an adapted input that loses a marker of a flagged line of the
// file before it, and copies no file (fix 2 of the first pilot, #97). It asks
// the call for each file that it adapts, not for a file that it writes; the
// evidence names the first lost marker and the count.
func TestS14RefusesAnInputThatLosesAMarker(t *testing.T) {
	f := &fakeRepo{disk: proseInputs()}
	f.install(t)
	c := proseCalls()
	var asked []string
	c.LostMarkers = func(path string, before, after []byte) []string {
		asked = append(asked, path+" | "+string(before)+" | "+string(after))
		return []string{path + ":3: the input loses the marker \u2039a\u203a", path + ":9: the input loses the marker \u2039b\u203a"}
	}
	o := runProseStep("S14", c, Input{Dir: "w"})
	if want := "docs/engineering-discipline.md:3: the input loses the marker \u2039a\u203a (2 lost markers in all)"; o.Kind != Fail || o.Evidence != want || len(f.files) != 0 {
		t.Errorf("S14: %s %q, %d files; want fail %q and no file", o.Kind, o.Evidence, len(f.files), want)
	}
	if want := []string{"docs/engineering-discipline.md | the file of the head, docs/engineering-discipline.md\n | text of docs/engineering-discipline.md\n"}; !reflect.DeepEqual(asked, want) {
		t.Errorf("the calls: %q; want %q", asked, want)
	}
	c.LostMarkers = nil
	if o := runProseStep("S14", c, Input{Dir: "w"}); o.Kind != Fail || o.Evidence != "the prose step has no call of the markers that an input loses" {
		t.Errorf("no call: %s %q", o.Kind, o.Evidence)
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
	// a marker over two lines of a CRLF file: its question, its where and its ask show its key, each line end as one space (fix 1 of the first pilot, #97)
	open := "\u2039open\r\nend\u203a"
	o = runS10(markersOf(Marker{"docs/c.md", 2, 0, open}), Input{Dir: "w"})
	if o.Kind != Stop || len(o.Stops) != 1 || o.Stops[0].Question != MarkerID("docs/c.md", "\u2039open end\u203a") ||
		o.Stops[0].Where != "docs/c.md:2 \u2039open end\u203a" || o.Stops[0].Ask != ask("docs/c.md", "\u2039open end\u203a") {
		t.Errorf("a marker over two lines of a CRLF file: %s %+v; want its key, each line end as one space", o.Kind, o.Stops)
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
	open := "\u2039open\r\nend\u203a"
	f2 := &fakeRepo{trees: f.trees, shows: f.shows, disk: map[string]string{"w/target/docs/c.md": "one\r\n" + open + "\r\n"}}
	f2.install(t)
	idOpen := MarkerID("docs/c.md", "\u2039open end\u203a")
	if o := runS11(markersOf(Marker{"docs/c.md", 2, 0, open}), Input{Dir: "w", Record: pinRecord(), Answers: mAnswers([]string{idOpen, "8080", "u", ""})}); o.Kind != Done ||
		f2.files["w/target/docs/c.md"] != "one\r\n8080\r\n" || o.Values[0][1] != "8080" {
		t.Errorf("a marker over two lines of a CRLF file: %s %q, %q; want the whole marker filled and the line end after it kept", o.Kind, f2.files["w/target/docs/c.md"], o.Values)
	}
	if !strings.Contains(record, "| Collected by | `layup setup`, step S11 |") || !strings.Contains(record, "# F-0002. The answers to the markers of the setup") ||
		!strings.Contains(record, "| Source | The answers of `inputs/answers.tsv` to the markers that S10 listed |") || strings.Contains(record, "\u2039") {
		t.Errorf("the second answers record:\n%s", record)
	}
}

// A line that holds more than one marker gives each place its own record row:
// the name of a row has the byte column of its marker after the line, so that
// the key (step, name) of the record is unique (the defect that the first
// pilot found, #97: one line of the baseline holds three markers). A line with
// one marker keeps the name marker:<file>:<line>.
func TestS11TwoMarkersOnOneLine(t *testing.T) {
	line := mx + " and " + my + " and " + mx
	colY, colX2 := len(mx)+len(" and "), len(mx)+len(" and ")+len(my)+len(" and ")
	c := markersOf(Marker{"docs/a.md", 2, 0, mx}, Marker{"docs/a.md", 2, colY, my}, Marker{"docs/a.md", 2, colX2, mx}, Marker{"docs/a.md", 3, 0, my})
	idx, idy := MarkerID("docs/a.md", mx), MarkerID("docs/a.md", my)
	a := mAnswers([]string{idx, "8080", "F-0003#5", ""}, []string{idy, "gap", "https://github.invalid/c/1", "Which y?"})
	f := &fakeRepo{
		trees: map[string][]git.TreeEntry{factsDir: {{Path: "docs/facts/F-0001-setup-answers.md"}, {Path: "docs/facts/README.md"}, {Path: "docs/facts/problem-statement-brief.md"}}},
		shows: map[string]string{factsSumsPath: "abc  docs/facts/F-0001-setup-answers.md\n", factsDir + "/README.md": factsIndexS04},
		disk:  map[string]string{"w/target/docs/a.md": "one\n" + line + "\n" + my + "\n"},
	}
	f.install(t)
	o := runS11(c, Input{Dir: "w", Record: pinRecord(), Answers: a})
	names := []string{"marker:docs/a.md:2:0", "marker:docs/a.md:2:" + strconv.Itoa(colY), "marker:docs/a.md:2:" + strconv.Itoa(colX2), "marker:docs/a.md:3"}
	if o.Kind != Done || len(o.Values) < 4 {
		t.Fatalf("S11: %s %q %q; want done with a row for each place", o.Kind, o.Evidence, o.Values)
	}
	for i, n := range names {
		if o.Values[i][0] != n {
			t.Errorf("row %d is named %q; want %q", i, o.Values[i][0], n)
		}
	}
	if got, want := f.files["w/target/docs/a.md"], "one\n8080 and "+my+" and 8080\n"+my+"\n"; got != want {
		t.Errorf("S11 wrote %q; want %q", got, want)
	}
	// the rows go into the record, whose key is (step, name): no key repeats
	var rows [][]string
	for _, v := range o.Values {
		rows = append(rows, append([]string{"S11"}, v...))
	}
	if err := tsv.Write(io.Discard, work.RecordSchema, rows); err != nil {
		t.Errorf("the record of the rows: %v; want no error", err)
	}
}

// S11 replaces the whole of a marker over more lines (fix 1 of the first
// pilot, #97: S11 filled the first line of a marker of docs/tasks/backlog.md
// only, and the rest of the placeholder stayed). The row of each place names
// its line in the tree that S11 writes, and a gap over more lines gets its key
// in its record row and in its row of open-gaps.tsv.
func TestS11AMarkerOverMoreLines(t *testing.T) {
	const file = "docs/tasks/backlog.md"
	scheme := "\u2039State your exact scheme\nhere \u2014 for example\nexist.\"\u203a"
	two := "\u2039State one\nthing\u203a"
	lead, col2 := "A suffix. `", 5+len(mx)+len(" and ")
	c := markersOf(Marker{file, 2, len(lead), scheme}, Marker{file, 5, 5, mx}, Marker{file, 5, col2, two})
	idS, idx, id2 := MarkerID(file, work.MarkerKey(scheme)), MarkerID(file, mx), MarkerID(file, work.MarkerKey(two))
	a := mAnswers([]string{idS, "T- plus four random lowercase letters or digits", "https://github.invalid/c/2", ""},
		[]string{idx, "gap", "https://github.invalid/c/2", "Which x?"}, []string{id2, "gap", "https://github.invalid/c/2", "What is the one thing?"})
	f := &fakeRepo{
		trees: map[string][]git.TreeEntry{factsDir: {{Path: "docs/facts/F-0001-setup-answers.md"}, {Path: "docs/facts/README.md"}, {Path: "docs/facts/problem-statement-brief.md"}}},
		shows: map[string]string{factsSumsPath: "abc  docs/facts/F-0001-setup-answers.md\n", factsDir + "/README.md": factsIndexS04},
		disk:  map[string]string{"w/target/" + file: "intro\n" + lead + scheme + "`\nnext " + mx + " and " + two + " end\n"},
	}
	f.install(t)
	o := runS11(c, Input{Dir: "w", Record: pinRecord(), Answers: a})
	if o.Kind != Done || len(o.Values) < 3 {
		t.Fatalf("S11: %s %q %q; want done with a row for each place", o.Kind, o.Evidence, o.Values)
	}
	want := [][]string{{"marker:" + file + ":2", "T- plus four random lowercase letters or digits", "answer", idS},
		{"marker:" + file + ":3:5", mx, "gap", work.OpenGapsPath}, {"marker:" + file + ":3:" + strconv.Itoa(col2), "\u2039State one thing\u203a", "gap", work.OpenGapsPath}}
	if !reflect.DeepEqual(o.Values[:3], want) {
		t.Errorf("the rows of the places:\n%q\nwant\n%q", o.Values[:3], want)
	}
	if got, w := f.files["w/target/"+file], "intro\n"+lead+"T- plus four random lowercase letters or digits`\nnext "+mx+" and "+two+" end\n"; got != w {
		t.Errorf("S11 wrote %q; want %q", got, w)
	}
	if got, w := f.files["w/target/"+work.OpenGapsPath], file+"\t"+mx+"\tWhich x?\n"+file+"\t\u2039State one thing\u203a\tWhat is the one thing?\n"; got != w {
		t.Errorf("S11 wrote the open gaps %q; want %q", got, w)
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
		{"a gap marker with a carriage return", []Marker{{"docs/a.md", 3, 0, "\u2039x\r\ny\u203a"}}, mAnswers([]string{MarkerID("docs/a.md", "\u2039x y\u203a"), "gap", "u", "Which?"}),
			"docs/a.md: the marker \u2039x y\u203a holds a tab or a carriage return, so it cannot have a row in docs/setup/open-gaps.tsv"},
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
