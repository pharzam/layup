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
type wanted struct {
	path, ask string
	adapt     bool // S14 adapts the file: its input keeps each place of each marker of the file (fix 2 of the first pilot, #97)
}

func writeAsk(p string) wanted {
	return wanted{p, "Write the text of " + p + " for the target, and give it as " + inputsDir + "/" + p + ".", false}
}
func adaptAsk(p string) wanted {
	return wanted{p, "Adapt " + p + ", which check adapted flags, to the target, and give it as " + inputsDir + "/" + p + ".", true}
}
func linkAsk(p string) wanted {
	return wanted{p, "Fix the links of " + p + " that the deletion of the baseline's history breaks, and give the file as " + inputsDir + "/" + p + ".", false}
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

// proseFiles gives the named file of each of S07, S08 and S09.
var proseFiles = map[string]string{"S07": onboardingPath, "S08": glossaryPath, "S09": guardrailsPath}

var proseEvidence = map[string]string{"S07": "check onboarding", "S08": "check glossary", "S09": "check guardrails", "S14": "checks identity and adapted"}

// proseWanted gives the files of each step of the prose step that is not done
// (O-123, D5 of #90 with condition 2 of its plan review and finding 1 of its
// review round 1): S07, S08 and S09 their named file; S14 README.md,
// AGENTS.md and each other file that check adapted flags on the head, less
// the named file of a step of S07 to S09 that is not done, so no question is
// twice in the group's table, and a named file that check adapted still flags
// after its step is S14's too.
func proseWanted(c Calls, in Input) (map[string][]wanted, error) {
	if c.Flagged == nil {
		return nil, errors.New("the prose step has no call of the files that check adapted flags")
	}
	flagged, err := c.Flagged(filepath.Join(in.Dir, work.TargetPath))
	if err != nil {
		return nil, errors.New("the files that check adapted flags: " + firstLine(err))
	}
	done := func(id string) bool { _, ok := in.Record.Value(id, "done"); return ok }
	out := map[string][]wanted{}
	for _, id := range []string{"S07", "S08", "S09"} {
		if !done(id) {
			out[id] = []wanted{writeAsk(proseFiles[id])}
		}
	}
	if !done("S14") {
		list := []wanted{writeAsk("README.md"), writeAsk("AGENTS.md")}
		for _, p := range flagged {
			if slices.ContainsFunc(list, func(w wanted) bool { return w.path == p }) {
				continue
			}
			if id := stepOf(p); id != "" && !done(id) {
				continue // the step of S07 to S09 asks for it
			}
			list = append(list, adaptAsk(p))
		}
		out["S14"] = list
	}
	return out, nil
}

// stepOf gives the step of S07 to S09 whose named file is p, or "".
func stepOf(p string) string {
	for id, f := range proseFiles {
		if f == p {
			return id
		}
	}
	return ""
}

// runProseStep is S07, S08, S09 or S14, the prose step (D5 of #90): the four
// steps are one unit of input. While an input of a step of the group that is
// not done is missing, no step copies a file: each gives the stop rows of its
// own missing inputs, and a step whose own inputs exist waits, with no row, so
// the run gives one table (O-123). Else the step copies its files, each with
// its record row (K42).
func runProseStep(step string, c Calls, in Input) Outcome {
	want, err := proseWanted(c, in)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	missing := false
	var mine []StopRow
	for _, id := range []string{"S07", "S08", "S09", "S14"} {
		for _, f := range want[id] {
			_, err := sys.read(inputFile(in.Dir, f.path))
			switch {
			case errors.Is(err, fs.ErrNotExist):
				missing = true
				if id == step {
					mine = append(mine, StopRow{step, "F-" + f.path, f.ask, f.path})
				}
			case err != nil:
				return Outcome{Kind: Fail, Evidence: fmt.Sprintf("%s/%s: %s", inputsDir, f.path, firstLine(err))}
			}
		}
	}
	if missing {
		return Outcome{Kind: Stop, Stops: mine}
	}
	if lost, err := lostMarkers(c, in, want[step]); err != nil || len(lost) > 0 {
		reason := ""
		switch {
		case err != nil:
			reason = err.Error()
		case len(lost) == 1:
			reason = lost[0]
		default:
			reason = fmt.Sprintf("%s (%d lost markers in all)", lost[0], len(lost))
		}
		return Outcome{Kind: Fail, Evidence: reason}
	}
	rows, _, err := copyInputs(step, in.Dir, want[step])
	if err != nil {
		return Outcome{Kind: Fail, Evidence: step + ": " + err.Error()}
	}
	return Outcome{Kind: Done, Evidence: proseEvidence[step], Commit: true, Values: rows}
}

// lostMarkers gives each marker that an input of an adapted file loses, place by
// place, from the file before it, by the call of check adapted (fix 2 of the first
// pilot, #97); a file that the step writes is not asked.
func lostMarkers(c Calls, in Input, files []wanted) ([]string, error) {
	var lost []string
	for _, f := range files {
		if !f.adapt {
			continue
		}
		if c.LostMarkers == nil {
			return nil, errors.New("the prose step has no call of the markers that an input loses")
		}
		before, err := sys.read(targetFile(in.Dir, f.path))
		if err != nil {
			return nil, fmt.Errorf("%s: %s", f.path, firstLine(err))
		}
		after, err := sys.read(inputFile(in.Dir, f.path))
		if err != nil {
			return nil, fmt.Errorf("%s/%s: %s", inputsDir, f.path, firstLine(err))
		}
		lost = append(lost, c.LostMarkers(f.path, before, after)...)
	}
	return lost, nil
}

// A question of S10 is one file and marker key (work.MarkerKey), with the
// first line where it occurs.
type question struct {
	id, file, text string
	line           int
}

// markerQuestions gives one question per file and marker text of marks, in
// the order of their first place.
func markerQuestions(marks []Marker) []question {
	var out []question
	for _, m := range marks {
		key := work.MarkerKey(m.Text)
		if !slices.ContainsFunc(out, func(q question) bool { return q.file == m.File && q.text == key }) {
			out = append(out, question{MarkerID(m.File, key), m.File, key, m.Line})
		}
	}
	return out
}

// markerAsk gives the question of a marker of a file (D6 of #90).
func markerAsk(file, marker string) string {
	return "What is the value of " + marker + " in " + file + "? Answer gap to keep it as an open gap, with its question as question_text."
}

// runS10 is S10 (D6 of #90): one question M-<x8> per file and marker key of
// the tree, at the line of its open quote, and one stop table for the markers
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

// fileOrder gives the file of each place, in the order of the places.
func fileOrder(marks []Marker) []string {
	var out []string
	for _, m := range marks {
		out = append(out, m.File)
	}
	return out
}

// lineStarts gives the byte offset of each line of s; index 0 is line 1.
func lineStarts(s string) []int {
	at := []int{0}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			at = append(at, i+1)
		}
	}
	return at
}

// fill replaces each place of file in text whose answer is a value with that
// value, the whole marker, also over more lines (fix 1 of the first pilot,
// #97), and gives the filled text and, for each place of file in marks (by
// its index), its byte offset in the filled text. A marker that is not at its
// line and column is an error.
func fill(text, file string, marks []Marker, answer func(Marker) []string) (string, map[int]int, error) {
	starts := lineStarts(text)
	var b strings.Builder
	at := map[int]int{}
	last, delta := 0, 0
	for i, m := range marks {
		if m.File != file {
			continue
		}
		if m.Line < 1 || m.Line > len(starts) || starts[m.Line-1]+m.Col+len(m.Text) > len(text) || text[starts[m.Line-1]+m.Col:starts[m.Line-1]+m.Col+len(m.Text)] != m.Text {
			return "", nil, fmt.Errorf("%s:%d: the marker %s is not at its column %d", m.File, m.Line, work.MarkerKey(m.Text), m.Col)
		}
		off := starts[m.Line-1] + m.Col
		at[i] = off + delta
		if r := answer(m); r[1] != "gap" {
			b.WriteString(text[last:off])
			b.WriteString(r[1])
			last = off + len(m.Text)
			delta += len(r[1]) - len(m.Text)
		}
	}
	b.WriteString(text[last:])
	return b.String(), at, nil
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
// each marker with a value at each of its places, the whole marker, also over
// more lines (fix 1 of the first pilot, #97), keeps each gap with its row of
// open-gaps.tsv by its key, writes one record row per place at its line in the
// tree that S11 writes, and writes the second answers record with its index
// row and its line of facts.sha256.
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
		case r[1] == "gap" && slices.ContainsFunc(marks, func(m Marker) bool {
			return m.File == q.file && work.MarkerKey(m.Text) == q.text && strings.ContainsAny(m.Text, "\t\r")
		}):
			return Outcome{Kind: Fail, Evidence: q.file + ": the marker " + q.text + " holds a tab or a carriage return, so it cannot have a row in " + work.OpenGapsPath}
		case r[1] != "gap" && strings.ContainsAny(r[1], "\u2039\u203a"):
			return Outcome{Kind: Fail, Evidence: q.id + ": the value holds an angle quote, so it would be a marker"}
		}
	}
	answer := func(m Marker) []string { return answerOf(in.Answers, MarkerID(m.File, work.MarkerKey(m.Text))) }
	out := map[string][]byte{}
	line := make([]int, len(marks)) // the line of each place in the tree that S11 writes
	starts := map[string][]int{}    // the byte offset of each line of a filled file at the scan; index 0 is line 1
	for i, m := range marks {
		line[i] = m.Line
	}
	for _, file := range slices.Compact(fileOrder(marks)) {
		if !slices.ContainsFunc(marks, func(m Marker) bool { return m.File == file && answer(m)[1] != "gap" }) {
			continue // a file with gaps only: S11 changes no line of it
		}
		data, err := sys.read(targetFile(in.Dir, file))
		if err != nil {
			return Outcome{Kind: Fail, Evidence: file + ": " + firstLine(err)}
		}
		starts[file] = lineStarts(string(data))
		filled, at, err := fill(string(data), file, marks, answer)
		if err != nil {
			return Outcome{Kind: Fail, Evidence: err.Error()}
		}
		for i, off := range at {
			line[i] = 1 + strings.Count(filled[:off], "\n")
		}
		out[file] = []byte(filled)
	}
	var values, gaps [][]string
	onLine := map[string]int{} // the places of each line of the tree that S11 writes: a line with more than one gives each row the column too, so each key of the record is unique
	first := map[string]int{}  // the first scan line of each such line: a filled marker over more lines joins scan lines
	for i, m := range marks {
		k := fmt.Sprintf("%s:%d", m.File, line[i])
		onLine[k]++
		if f, ok := first[k]; !ok || m.Line < f {
			first[k] = m.Line
		}
	}
	for i, m := range marks {
		r, key := answer(m), work.MarkerKey(m.Text)
		k := fmt.Sprintf("%s:%d", m.File, line[i])
		name := "marker:" + k
		if onLine[k] > 1 {
			col := m.Col // the byte column of the open quote at the scan, counted from the first scan line of its line
			if f := first[k]; m.Line > f {
				col += starts[m.File][m.Line-1] - starts[m.File][f-1]
			}
			name += fmt.Sprintf(":%d", col)
		}
		if r[1] == "gap" {
			values = append(values, []string{name, key, "gap", work.OpenGapsPath})
			if !slices.ContainsFunc(gaps, func(g []string) bool { return g[0] == m.File && g[1] == key }) {
				gaps = append(gaps, []string{m.File, key, r[4]})
			}
			continue
		}
		values = append(values, append([]string{name, r[1]}, sourceOf(r)...))
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
