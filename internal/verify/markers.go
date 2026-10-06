package verify

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/work"
)

// mkExempt is MK_EXEMPT of setup-check.sh: the paths that check markers and
// S10 do not read, embedded at the engine's version (D1 of #87);
// TestTheExemptionsOfMarkersEqualTheSh compares it with the sh file.
const mkExempt = `^(docs/(adr|ci|links|prd|setup)/tests/|\.githooks/tests/|docs/templates/)|^docs/[^/]+/template\.md$|^docs/tests/template-[^/]*\.md$|^docs/tests/traceability-template\.md$|^docs/adr/000[1-8]-[^/]*\.md$|^docs/links/link-lint\.sh$|^docs/prd/prd-lint\.sh$|^docs/setup/setup-check\.sh$|^docs/setup/open-gaps\.tsv$`

var mkExemptRE = regexp.MustCompile(mkExempt)

// A Marker is one occurrence of a marker in a tracked file: its file, the
// line of its open quote, and its text with the two angle quotes and each line
// end between them. S10 asks one question per file and key (work.MarkerKey),
// and check sources reads the line (D1 of #87; fix 1 of the first pilot, #97).
type Marker struct {
	File string
	Line int
	Col  int // the byte column of its open quote on its line (D7 of #90)
	Text string
}

// TextMarkers gives each marker of text with its line and column, and no
// file: internal/cli refuses a brief that holds one (D4 of #90).
func TextMarkers(text string) []Marker { return textMarkers("", text) }

// FirstUnpaired gives the line and the quote of the first angle quote of text
// with no pair, or 0: internal/cli refuses a brief that holds one, as it
// refuses a marker (fix 3 of the first pilot, #97).
func FirstUnpaired(text string) (int, string) {
	for _, it := range scanText(text) {
		if it.kind == aLoneOpen || it.kind == aLoneClose {
			return it.line, text[it.at:it.end]
		}
	}
	return 0, ""
}

// textMarkers gives the markers of the text of file.
func textMarkers(file, text string) []Marker {
	var out []Marker
	for _, it := range scanText(text) {
		if it.kind == aMarker {
			out = append(out, Marker{file, it.line, it.col, text[it.at:it.end]})
		}
	}
	return out
}

// Markers gives each marker of the work tree at tree, in the order of the
// files of git and of their lines: the one scanner of check markers and of
// S10, which gets the list through internal/cli (D1 of #87).
func Markers(tree string) ([]Marker, error) {
	files, err := git.LsFiles(tree)
	if err != nil {
		return nil, err
	}
	return scanMarkers(os.DirFS(tree), files), nil
}

// scanMarkers reads each tracked file outside MK_EXEMPT that is a regular
// file, as check_markers does.
func scanMarkers(fsys fs.FS, files []string) []Marker {
	marks, _ := scanFiles(fsys, files)
	return marks
}

// scanFiles gives the markers of the files that scanMarkers reads, and the
// finding of each angle quote of them with no pair.
func scanFiles(fsys fs.FS, files []string) ([]Marker, []string) {
	var marks []Marker
	var lone []string
	for _, p := range files {
		if mkExemptRE.MatchString(p) || !isRegular(fsys, p) {
			continue
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			continue
		}
		for _, it := range scanText(string(data)) {
			switch it.kind {
			case aMarker:
				marks = append(marks, Marker{p, it.line, it.col, string(data[it.at:it.end])})
			case aLoneOpen, aLoneClose:
				lone = append(lone, fmt.Sprintf("unpaired: %s:%d %s", p, it.line, data[it.at:it.end]))
			}
		}
	}
	return marks, lone
}

const (
	mkOpen  = "\u2039"
	mkClose = "\u203a"
	mkForm  = mkOpen + "\u2026" + mkClose // the convention's own text, not a marker
)

// An item is one place of an angle quote in a text, by the scan of scanText:
// the byte offsets of its first byte and of the byte after it, the line and
// the byte column of its first byte, and its kind.
type item struct {
	at, end, line, col int
	kind               itemKind
}

type itemKind int

const (
	aMarker    itemKind = iota // a marker: an open quote to the first close quote after it
	aForm                      // the text of the convention, the two quotes around an ellipsis: no marker
	aMention                   // an open or a close quote between two backticks
	aLoneOpen                  // an open quote with no close quote after it: no pair
	aLoneClose                 // a close quote that no open quote takes: no pair
)

// scanText gives each item of a text, in the order of the text (fix 1 and
// fix 3 of the first pilot, #97): an open quote starts a marker to the first
// close quote after it, also across lines; an open quote with no close quote
// after it starts none, and has no pair; a quote between two backticks is the
// mention; the text of the convention is no marker; a close quote that no
// open quote takes has no pair.
func scanText(text string) []item {
	var out []item
	line, from, start := 1, 0, 0 // the line of the byte at from, and the start of that line
	add := func(at, end int, kind itemKind) {
		for ; from < at; from++ {
			if text[from] == '\n' {
				line, start = line+1, from+1
			}
		}
		out = append(out, item{at, end, line, at - start, kind})
	}
	mention := func(at, n int) bool { return at > 0 && text[at-1] == '`' && at+n < len(text) && text[at+n] == '`' }
	for i := 0; i < len(text); {
		o, c := strings.Index(text[i:], mkOpen), strings.Index(text[i:], mkClose)
		if o < 0 && c < 0 {
			break
		}
		if c >= 0 && (o < 0 || c < o) { // a close quote that no open quote takes
			at := i + c
			if mention(at, len(mkClose)) {
				add(at, at+len(mkClose), aMention)
			} else {
				add(at, at+len(mkClose), aLoneClose)
			}
			i = at + len(mkClose)
			continue
		}
		at := i + o
		after := at + len(mkOpen)
		k := strings.Index(text[after:], mkClose)
		switch {
		case mention(at, len(mkOpen)):
			add(at, after, aMention)
			i = after
		case k < 0:
			add(at, after, aLoneOpen)
			i = after
		case text[at:after+k+len(mkClose)] == mkForm:
			add(at, after+k+len(mkClose), aForm)
			i = after + k + len(mkClose)
		default:
			add(at, after+k+len(mkClose), aMarker)
			i = after + k + len(mkClose)
		}
	}
	return out
}

// markersFindings is check_markers (D2 of #87): each marker of the tree needs
// a row in docs/setup/open-gaps.tsv by its key, each row needs its marker in
// the tree, each row needs a question that is not blank or the empty mark, and
// no angle quote is without its pair (fix 3 of the first pilot, #97). The file
// is read by tabs, by the columns of its block, as check_markers reads it; an
// absent file is no rows.
func markersFindings(fsys fs.FS, h history) []string {
	files, err := h.Files()
	if err != nil {
		return []string{"git: cannot list the tracked files"}
	}
	found := map[string]bool{}
	marks, lone := scanFiles(fsys, files)
	for _, m := range marks {
		found[m.File+"\t"+work.MarkerKey(m.Text)] = true
	}
	var out []string
	listed := map[string]bool{}
	data, _ := fs.ReadFile(fsys, work.OpenGapsPath)
	for i, row := range openGapsRows(data) {
		if row == nil {
			continue
		}
		listed[row[0]+"\t"+row[1]] = true
		if q := strings.Trim(row[2], " \t\r"); q == "" || q == "\u2014" {
			out = append(out, fmt.Sprintf("question: %s line %d has no question", work.OpenGapsPath, i+1))
		}
	}
	for _, k := range slices.Sorted(mapKeys(found)) {
		if !listed[k] {
			f, m, _ := strings.Cut(k, "\t")
			out = append(out, "unlisted: "+f+" "+m)
		}
	}
	for _, k := range slices.Sorted(mapKeys(listed)) {
		if !found[k] {
			f, m, _ := strings.Cut(k, "\t")
			out = append(out, "stale: "+f+" "+m+" is listed in "+work.OpenGapsPath+" but does not occur")
		}
	}
	return append(out, lone...)
}

// openGapsRows gives each line of the file as its three columns (file,
// marker, question), split at tabs; a missing column is empty, and an empty
// line is nil, as awk -F'\t' and cut read them.
func openGapsRows(data []byte) [][]string {
	if len(data) == 0 {
		return nil
	}
	var out [][]string
	for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
		if line == "" {
			out = append(out, nil)
			continue
		}
		f := append(strings.Split(line, "\t"), "", "")
		out = append(out, f[:3])
	}
	return out
}

// mapKeys gives the keys of a set, for slices.Sorted.
func mapKeys(m map[string]bool) func(func(string) bool) {
	return func(yield func(string) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}

// checkMarkers is check markers of a target: the core alone (setup.md:
// "check_markers, with LAYUP's MK_EXEMPT").
func checkMarkers(in input) []string { return markersFindings(in.fsys, in.repo) }
