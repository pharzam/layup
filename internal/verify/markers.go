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

// A Marker is one occurrence of a marker in a tracked file: its file, its
// line and its text with the two angle quotes. S10 asks one question per file
// and text, and check sources reads the line (D1 of #87).
type Marker struct {
	File string
	Line int
	Text string
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
// file, line by line, as check_markers does.
func scanMarkers(fsys fs.FS, files []string) []Marker {
	var out []Marker
	for _, p := range files {
		if mkExemptRE.MatchString(p) || !isRegular(fsys, p) {
			continue
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			continue
		}
		for i, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
			for _, m := range lineMarkers(line) {
				out = append(out, Marker{p, i + 1, m})
			}
		}
	}
	return out
}

const (
	mkOpen  = "\u2039"
	mkClose = "\u203a"
	mkForm  = mkOpen + "\u2026" + mkClose // the convention's own text, not a marker
)

// lineMarkers gives the markers of one line, as the awk program of
// check_markers reads it: each open quote starts a marker to the first close
// quote after it, or to the line end; an open quote between two backticks is
// the mention and is skipped; the text of the convention is not a marker.
func lineMarkers(line string) []string {
	var out []string
	for i := 0; ; {
		j := strings.Index(line[i:], mkOpen)
		if j < 0 {
			return out
		}
		at := i + j
		after := at + len(mkOpen)
		if at > 0 && line[at-1] == '`' && after < len(line) && line[after] == '`' {
			i = after
			continue
		}
		m := line[at:]
		if k := strings.Index(line[after:], mkClose); k >= 0 {
			m = line[at : after+k+len(mkClose)]
		}
		if m != mkForm {
			out = append(out, m)
		}
		i = at + len(m)
	}
}

// markersFindings is check_markers (D2 of #87): each marker of the tree needs
// a row in docs/setup/open-gaps.tsv, each row needs its marker in the tree,
// and each row needs a question that is not blank or the empty mark. The file
// is read by tabs, by the columns of its block, as check_markers reads it; an
// absent file is no rows.
func markersFindings(fsys fs.FS, h history) []string {
	files, err := h.Files()
	if err != nil {
		return []string{"git: cannot list the tracked files"}
	}
	found := map[string]bool{}
	for _, m := range scanMarkers(fsys, files) {
		found[m.File+"\t"+m.Text] = true
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
	return out
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
