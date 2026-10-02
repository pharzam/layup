package setup

// The steps S05 to S11 and S14 (task T-b3r1, #90): the history of the
// baseline, the briefs, the prose step, and the markers. Each step reads the
// files that it changes from the head of layup-setup, where the runner puts
// the target before it (D11 of #90), so a rerun gives the same files and rows.

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// The named files of the prose step (D5 of #90): S07, S08 and S09 copy one
// each, and S14 copies README.md and AGENTS.md; the input files of a work
// area are under inputsDir.
const (
	onboardingPath = "docs/onboarding-for-engineers.md"
	glossaryPath   = "docs/glossary.md"
	guardrailsPath = "docs/guardrails.md"
	inputsDir      = "inputs/files"
)

var named = []string{onboardingPath, glossaryPath, guardrailsPath, "README.md", "AGENTS.md"}

// targetFile and inputFile give a file of the target and an input file of the
// work area at dir.
func targetFile(dir, p string) string {
	return filepath.Join(dir, work.TargetPath, filepath.FromSlash(p))
}
func inputFile(dir, p string) string {
	return filepath.Join(dir, filepath.FromSlash(inputsDir), filepath.FromSlash(p))
}

// sha gives the SHA-256 of data in hexadecimal.
func sha(data []byte) string { return fmt.Sprintf("%x", sha256.Sum256(data)) }

// fileRow gives the record row of a file that a step copied from its input
// file (K42, D3 of #90).
func fileRow(p string, data []byte) []string {
	return []string{"file:" + p, sha(data), "computed", "sha256 " + inputsDir + "/" + p}
}

// A wanted file is a file of the target that a step copies from its input
// file, with the question of its stop row.
type wanted struct{ path, ask string }

func writeAsk(p string) wanted {
	return wanted{p, "Write the text of " + p + " for the target, and give it as " + inputsDir + "/" + p + "."}
}
func adaptAsk(p string) wanted {
	return wanted{p, "Adapt " + p + ", which check adapted flags, to the target, and give it as " + inputsDir + "/" + p + "."}
}
func linkAsk(p string) wanted {
	return wanted{p, "Fix the links of " + p + " that the deletion of the baseline's history breaks, and give the file as " + inputsDir + "/" + p + "."}
}

// copyInputs copies each file from its input file into the target, with its
// record row, or gives one stop row of step per file with no input file; it
// writes nothing when a row stops.
func copyInputs(step, dir string, files []wanted) ([][]string, []StopRow, error) {
	data := map[string][]byte{}
	var stops []StopRow
	for _, f := range files {
		b, err := sys.read(inputFile(dir, f.path))
		switch {
		case errors.Is(err, fs.ErrNotExist):
			stops = append(stops, StopRow{step, "F-" + f.path, f.ask, f.path})
		case err != nil:
			return nil, nil, fmt.Errorf("%s/%s: %s", inputsDir, f.path, firstLine(err))
		default:
			data[f.path] = b
		}
	}
	if stops != nil {
		return nil, stops, nil
	}
	var rows [][]string
	for _, f := range files {
		if err := sys.write(targetFile(dir, f.path), data[f.path]); err != nil {
			return nil, nil, err
		}
		rows = append(rows, fileRow(f.path, data[f.path]))
	}
	return rows, nil, nil
}

var taskFile = regexp.MustCompile(`^docs/tasks/(T-[^/]*)\.md$`)

// runS05 is S05 (D1 to D3 of #90): from the head of layup-setup, it deletes
// docs/decisions/, docs/audit/ and each docs/tasks/T-*.md, removes each line
// of the two task indexes that names a deleted task or links the baseline's
// repository, with each blockquote that holds one, and copies the input file
// of each file whose links then break.
func runS05(c Calls, in Input) Outcome {
	if c.BrokenLinks == nil || c.LinksBaseline == nil {
		return Outcome{Kind: Fail, Evidence: "S05 has no call of the links of the tree"}
	}
	target := filepath.Join(in.Dir, work.TargetPath)
	source, _ := in.Record.Value("S02", "pin.source")
	var gone, ids []string
	for _, dir := range []string{"docs/decisions", "docs/audit", "docs/tasks"} {
		entries, err := sys.lsTree(target, "HEAD", dir)
		if err != nil {
			return Outcome{Kind: Fail, Evidence: dir + "/ of the head: " + firstLine(err)}
		}
		if dir != "docs/tasks" && len(entries) > 0 {
			gone = append(gone, dir)
		}
		for _, e := range entries {
			if m := taskFile.FindStringSubmatch(e.Path); dir == "docs/tasks" && m != nil {
				gone, ids = append(gone, e.Path), append(ids, m[1])
			}
		}
	}
	for _, p := range gone {
		if err := sys.removeAll(targetFile(in.Dir, p)); err != nil {
			return Outcome{Kind: Fail, Evidence: "S05: " + firstLine(err)}
		}
	}
	for _, index := range []string{"docs/tasks/backlog.md", "docs/tasks/completed.md"} {
		text, err := sys.show(target, "HEAD", index)
		if err != nil {
			continue // a baseline with no such index
		}
		kept := removeLines(string(text), func(line string) bool { return namesTask(line, ids) || c.LinksBaseline(source, line) })
		if kept == string(text) {
			continue
		}
		if err := sys.write(targetFile(in.Dir, index), []byte(kept)); err != nil {
			return Outcome{Kind: Fail, Evidence: "S05: " + firstLine(err)}
		}
	}
	files, err := c.BrokenLinks(target)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the links of the tree: " + firstLine(err)}
	}
	var want []wanted
	for _, f := range files {
		want = append(want, linkAsk(f))
	}
	rows, stops, err := copyInputs("S05", in.Dir, want)
	switch {
	case err != nil:
		return Outcome{Kind: Fail, Evidence: "S05: " + err.Error()}
	case stops != nil:
		return Outcome{Kind: Stop, Stops: stops}
	}
	if len(files) > 0 {
		still, err := c.BrokenLinks(target)
		if err != nil {
			return Outcome{Kind: Fail, Evidence: "the links of the tree: " + firstLine(err)}
		}
		if len(still) > 0 {
			return Outcome{Kind: Fail, Evidence: "the links of these files still break after their input files: " + strings.Join(still, ", ")}
		}
	}
	return Outcome{Kind: Done, Evidence: "checks kit-history and link-lint", Commit: true, Values: rows}
}

// removeLines gives text with each line that drop matches removed, and each
// blockquote (a run of lines that start with >) that holds such a line (D1 of
// #90).
func removeLines(text string, drop func(line string) bool) string {
	lines := strings.SplitAfter(text, "\n")
	match := func(l string) bool { return drop(strings.TrimRight(l, "\r\n")) }
	var out []string
	for i := 0; i < len(lines); {
		j := i + 1
		if strings.HasPrefix(lines[i], ">") {
			for j < len(lines) && strings.HasPrefix(lines[j], ">") {
				j++
			}
		}
		if !slices.ContainsFunc(lines[i:j], match) {
			out = append(out, lines[i:j]...)
		}
		i = j
	}
	return strings.Join(out, "")
}

// namesTask reports whether line holds one of ids as a word: with no letter,
// digit, - or _ just before or after it (D1 of #90).
func namesTask(line string, ids []string) bool {
	word := func(r rune) bool { return r == '-' || r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r) }
	for _, id := range ids {
		for i := 0; ; {
			j := strings.Index(line[i:], id)
			if j < 0 {
				break
			}
			at, end := i+j, i+j+len(id)
			before, _ := utf8.DecodeLastRuneInString(line[:at])
			after, _ := utf8.DecodeRuneInString(line[end:])
			if (at == 0 || !word(before)) && (end == len(line) || !word(after)) {
				return true
			}
			i = at + 1
		}
	}
	return false
}

// A brief is one brief of the work area that S06 copies (D4 of #90).
type brief struct {
	input, file, name, words string
	required                 bool
}

var briefs = []brief{
	{"problem-statement.md", "problem-statement-brief.md", "brief", "The problem statement of the idea owner, from the work area", true},
	{"vision.md", "architectural-vision-brief.md", "vision", "The vision brief of the idea owner, from the work area", false},
}

// runS06 is S06 (D4 of #90): it copies each brief byte for byte into
// docs/facts/, with its line in facts.sha256, its index row and its two
// record rows; the vision brief only when it exists.
func runS06(in Input) Outcome {
	target := filepath.Join(in.Dir, work.TargetPath)
	at, _ := in.Record.Value("S02", "pin.time")
	date := at[:min(10, len(at))]
	sums, _ := sys.show(target, "HEAD", factsSumsPath) // S04 wrote it
	if len(sums) > 0 && sums[len(sums)-1] != '\n' {
		sums = append(sums, '\n')
	}
	index, err := sys.show(target, "HEAD", factsDir+"/README.md")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: factsDir + "/README.md of the head: " + firstLine(err)}
	}
	text := string(index)
	var values [][]string
	for _, b := range briefs {
		data, err := sys.read(filepath.Join(in.Dir, "inputs", "briefs", b.input))
		if errors.Is(err, fs.ErrNotExist) && !b.required {
			continue
		}
		if err != nil {
			return Outcome{Kind: Fail, Evidence: "inputs/briefs/" + b.input + ": " + firstLine(err)}
		}
		path := factsDir + "/" + b.file
		if err := sys.write(targetFile(in.Dir, path), data); err != nil {
			return Outcome{Kind: Fail, Evidence: "S06: " + firstLine(err)}
		}
		sums = append(sums, sha(data)+"  "+path+"\n"...)
		if text, err = addIndexRow(text, "| ["+b.file+"]("+b.file+") | "+b.words+" | "+date+" | Raw |"); err != nil {
			return Outcome{Kind: Fail, Evidence: factsDir + "/README.md: " + err.Error()}
		}
		values = append(values, []string{b.name + ".copy", path, "computed", "the raw file name of a brief"},
			[]string{b.name + ".copy.sha256", sha(data), "computed", "sha256 " + path})
	}
	for p, data := range map[string][]byte{factsSumsPath: sums, factsDir + "/README.md": []byte(text)} {
		if err := sys.write(targetFile(in.Dir, p), data); err != nil {
			return Outcome{Kind: Fail, Evidence: "S06: " + firstLine(err)}
		}
	}
	return Outcome{Kind: Done, Evidence: "check facts", Commit: true, Values: values}
}

var proseCheck = map[string]string{"S07": "onboarding", "S08": "glossary", "S09": "guardrails"}

// runProse is S07, S08 or S09 (D5 of #90): it copies its named file from its
// input file, with its record row (K42).
func runProse(step string, in Input, paths []string) Outcome {
	var want []wanted
	for _, p := range paths {
		want = append(want, writeAsk(p))
	}
	rows, stops, err := copyInputs(step, in.Dir, want)
	switch {
	case err != nil:
		return Outcome{Kind: Fail, Evidence: step + ": " + err.Error()}
	case stops != nil:
		return Outcome{Kind: Stop, Stops: stops}
	}
	return Outcome{Kind: Done, Evidence: "check " + proseCheck[step], Commit: true, Values: rows}
}

// runS14 is S14 (D5 of #90, condition 2 of its plan review): it copies
// README.md, AGENTS.md, and each other file that check adapted flags on the
// head, less the files of S07 to S09, from their input files.
func runS14(c Calls, in Input) Outcome {
	if c.Flagged == nil {
		return Outcome{Kind: Fail, Evidence: "S14 has no call of the files that check adapted flags"}
	}
	flagged, err := c.Flagged(filepath.Join(in.Dir, work.TargetPath))
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the files that check adapted flags: " + firstLine(err)}
	}
	want := []wanted{writeAsk("README.md"), writeAsk("AGENTS.md")}
	for _, p := range flagged {
		if !slices.Contains(named, p) && !slices.ContainsFunc(want, func(w wanted) bool { return w.path == p }) {
			want = append(want, adaptAsk(p))
		}
	}
	rows, stops, err := copyInputs("S14", in.Dir, want)
	switch {
	case err != nil:
		return Outcome{Kind: Fail, Evidence: "S14: " + err.Error()}
	case stops != nil:
		return Outcome{Kind: Stop, Stops: stops}
	}
	return Outcome{Kind: Done, Evidence: "checks identity and adapted", Commit: true, Values: rows}
}

// A question of S10 is one file and marker text, with the first line where
// it occurs.
type question struct {
	id, file, text string
	line           int
}

// markerQuestions gives one question per file and marker text of marks, in
// the order of their first place.
func markerQuestions(marks []Marker) []question {
	var out []question
	for _, m := range marks {
		if !slices.ContainsFunc(out, func(q question) bool { return q.file == m.File && q.text == m.Text }) {
			out = append(out, question{MarkerID(m.File, m.Text), m.File, m.Text, m.Line})
		}
	}
	return out
}

// markerAsk gives the question of a marker of a file (D6 of #90).
func markerAsk(file, marker string) string {
	return "What is the value of " + marker + " in " + file + "? Answer gap to keep it as an open gap, with its question as question_text."
}

// runS10 is S10 (D6 of #90): one question M-<x8> per file and marker of the
// tree, with the first line of the marker, and one stop table for the markers
// with no answer; an M- answer to a marker that the tree does not hold is an
// input error.
func runS10(c Calls, in Input) Outcome {
	if c.Markers == nil {
		return Outcome{Kind: Fail, Evidence: "S10 has no call of the markers of the tree"}
	}
	marks, err := c.Markers(filepath.Join(in.Dir, work.TargetPath))
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the markers of the tree: " + firstLine(err)}
	}
	qs := markerQuestions(marks)
	var ids []string
	for _, q := range qs {
		ids = append(ids, q.id)
	}
	if err := CheckAsked(in.Answers, []string{"M-"}, ids); err != nil {
		return Outcome{Kind: Invalid, Evidence: err.Error()}
	}
	var stops []StopRow
	for _, q := range qs {
		if answerOf(in.Answers, q.id) == nil {
			stops = append(stops, StopRow{"S10", q.id, markerAsk(q.file, q.text), fmt.Sprintf("%s:%d %s", q.file, q.line, q.text)})
		}
	}
	if stops != nil {
		return Outcome{Kind: Stop, Stops: stops}
	}
	return Outcome{Kind: Done, Evidence: "every marker has an answer row"}
}

// answerOf gives the row of answers.tsv for the question id with a value, or
// nil.
func answerOf(a work.Answers, id string) []string {
	for _, r := range a {
		if r[0] == id && r[1] != "" {
			return r
		}
	}
	return nil
}

// runS11 is S11 (D8 of #90): it checks each answer and each marker, then fills
// each marker with a value at each of its places, keeps each gap with its row
// of open-gaps.tsv, writes one record row per place, and writes the second
// answers record with its index row and its line of facts.sha256.
func runS11(c Calls, in Input) Outcome {
	const evidence = "checks markers, sources and facts"
	if c.Markers == nil {
		return Outcome{Kind: Fail, Evidence: "S11 has no call of the markers of the tree"}
	}
	target := filepath.Join(in.Dir, work.TargetPath)
	marks, err := c.Markers(target)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the markers of the tree: " + firstLine(err)}
	}
	if len(marks) == 0 {
		return Outcome{Kind: Done, Evidence: evidence, Commit: true}
	}
	qs := markerQuestions(marks)
	for _, q := range qs { // each check before any write (condition 1 of the plan review)
		r := answerOf(in.Answers, q.id)
		switch {
		case r == nil:
			return Outcome{Kind: Fail, Evidence: q.id + ": no answer"}
		case r[1] == "gap" && strings.Contains(q.file, "\t"):
			return Outcome{Kind: Fail, Evidence: q.file + ": a file whose name holds a tab cannot have a row in " + work.OpenGapsPath}
		case r[1] != "gap" && strings.ContainsAny(r[1], "\u2039\u203a"):
			return Outcome{Kind: Fail, Evidence: q.id + ": the value holds an angle quote, so it would be a marker"}
		}
	}
	texts := map[string][]string{} // the lines of each file with a value to fill
	var values, gaps [][]string
	for _, m := range marks {
		r := answerOf(in.Answers, MarkerID(m.File, m.Text))
		name := fmt.Sprintf("marker:%s:%d", m.File, m.Line)
		if r[1] == "gap" {
			values = append(values, []string{name, m.Text, "gap", work.OpenGapsPath})
			if !slices.ContainsFunc(gaps, func(g []string) bool { return g[0] == m.File && g[1] == m.Text }) {
				gaps = append(gaps, []string{m.File, m.Text, r[4]})
			}
			continue
		}
		values = append(values, append([]string{name, r[1]}, sourceOf(r)...))
		if texts[m.File] == nil {
			data, err := sys.read(targetFile(in.Dir, m.File))
			if err != nil {
				return Outcome{Kind: Fail, Evidence: m.File + ": " + firstLine(err)}
			}
			texts[m.File] = strings.Split(string(data), "\n")
		}
	}
	for file, lines := range texts { // from the last place of a line to the first, so each column stays true
		for i := len(marks) - 1; i >= 0; i-- {
			m := marks[i]
			if m.File != file {
				continue
			}
			r := answerOf(in.Answers, MarkerID(m.File, m.Text))
			if r[1] == "gap" {
				continue
			}
			line := lines[m.Line-1]
			if m.Col+len(m.Text) > len(line) || line[m.Col:m.Col+len(m.Text)] != m.Text {
				return Outcome{Kind: Fail, Evidence: fmt.Sprintf("%s:%d: the marker %s is not at its column %d", m.File, m.Line, m.Text, m.Col)}
			}
			lines[m.Line-1] = line[:m.Col] + r[1] + line[m.Col+len(m.Text):]
		}
	}
	out := map[string][]byte{}
	for file, lines := range texts {
		out[file] = []byte(strings.Join(lines, "\n"))
	}
	if len(gaps) > 0 {
		list, _ := sys.show(target, "HEAD", work.OpenGapsPath)
		if len(list) > 0 && list[len(list)-1] != '\n' {
			list = append(list, '\n')
		}
		var b bytes.Buffer
		if err := tsv.Write(&b, work.OpenGapsSchema, gaps); err != nil {
			return Outcome{Kind: Fail, Evidence: work.OpenGapsPath + ": " + err.Error()}
		}
		out[work.OpenGapsPath] = append(list, b.Bytes()...)
	}
	at, _ := in.Record.Value("S02", "pin.time")
	date := at[:min(10, len(at))]
	n, err := nextNumber(target, "HEAD", factsDir, `^F-([0-9]{4})-[^/]*\.md$`)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	id := fmt.Sprintf("F-%04d", n)
	var asked []work.Question
	for _, q := range qs {
		asked = append(asked, work.Question{ID: q.id, Text: markerAsk(q.file, q.text)})
	}
	record := answersRecord(id, date, markerRecord, asked, in.Answers)
	recordPath := factsDir + "/" + id + "-marker-answers.md"
	sums, _ := sys.show(target, "HEAD", factsSumsPath)
	if len(sums) > 0 && sums[len(sums)-1] != '\n' {
		sums = append(sums, '\n')
	}
	index, err := sys.show(target, "HEAD", factsDir+"/README.md")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: factsDir + "/README.md of the head: " + firstLine(err)}
	}
	text, err := addIndexRow(string(index), fmt.Sprintf("| [%s](%s-marker-answers.md) | %s | %s | Raw |", id, id, markerRecord.title, date))
	if err != nil {
		return Outcome{Kind: Fail, Evidence: factsDir + "/README.md: " + err.Error()}
	}
	out[recordPath] = []byte(record)
	out[factsSumsPath] = append(sums, sha([]byte(record))+"  "+recordPath+"\n"...)
	out[factsDir+"/README.md"] = []byte(text)
	for p, data := range out {
		if err := sys.write(targetFile(in.Dir, p), data); err != nil {
			return Outcome{Kind: Fail, Evidence: "S11: " + firstLine(err)}
		}
	}
	values = append(values, []string{"marker.record", recordPath, "computed", "the next free ID of " + factsDir + "/"},
		[]string{"marker.record.sha256", sha([]byte(record)), "computed", "sha256 " + recordPath})
	return Outcome{Kind: Done, Evidence: evidence, Commit: true, Values: values}
}
