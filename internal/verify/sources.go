package verify

import (
	"bytes"
	"io/fs"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/catalog"
	"github.com/pharzam/layup/internal/work"
)

// catalogHas reports whether the catalog entry of stack in the binary has the
// file of ref (<stack>/<path>, catalog.Entry.Has; D10 of #91). A stack with
// no entry has no file.
var catalogHas = func(stack, ref string) bool {
	e, err := catalog.Read(catalog.Embedded(), stack)
	return err == nil && e.Has(ref)
}

// checkSources is check sources of a target (D4 of #87, with note 4 of its
// plan review): each value row of the setup record (each row but a done row)
// has a source that resolves.
func checkSources(in input) []string {
	var out []string
	stack, missing := value(in.record, "S01", "stack")
	for _, r := range in.record {
		step, name, val, source, ref := r[0], r[1], r[2], r[3], r[4]
		if name == "done" {
			continue
		}
		at := "source: " + step + " " + name + ": "
		switch source {
		case "answer":
			if !slices.ContainsFunc(in.answers, func(a []string) bool { return a[0] == ref }) {
				out = append(out, at+"the answer "+ref+" is not a row of "+work.AnswersPath)
			}
		case "catalog":
			if !catalogHas(stack, ref) {
				out = append(out, at+"the catalog entry "+stack+" has no file "+ref)
				if missing != "" {
					out = append(out, missing)
				}
			}
		case "fact":
			if f := resolveFact(in.fsys, ref); f != "" {
				out = append(out, at+f)
			}
		case "gap":
			out = append(out, gapFindings(in, at, name, val)...)
		case "computed":
			out = append(out, computedFindings(in, at, ref)...)
		case "step":
			out = append(out, at+"a value row with the source step")
		}
	}
	return out
}

// gapFindings checks a gap row of S11: its name is marker:<file>:<line>, the
// line of the file holds the marker of its value, and docs/setup/open-gaps.tsv
// has the row of the file and the marker.
func gapFindings(in input, at, name, marker string) []string {
	rest, ok := strings.CutPrefix(name, "marker:")
	i := strings.LastIndex(rest, ":")
	n, err := strconv.Atoi(rest[i+1:])
	if !ok || i < 0 || err != nil {
		return []string{at + "a gap row is not named marker:<file>:<line>"}
	}
	file := rest[:i]
	var out []string
	data, _ := fs.ReadFile(in.fsys, file)
	lines := bytes.Split(data, []byte("\n"))
	if n < 1 || n > len(lines) || !bytes.Contains(lines[n-1], []byte(marker)) {
		out = append(out, at+"line "+strconv.Itoa(n)+" of "+file+" does not hold "+marker)
	}
	gaps, _ := fs.ReadFile(in.fsys, work.OpenGapsPath)
	if !slices.ContainsFunc(openGapsRows(gaps), func(r []string) bool { return r != nil && r[0] == file && r[1] == marker }) {
		out = append(out, at+work.OpenGapsPath+" has no row for "+file+" "+marker)
	}
	return out
}

// computedFindings checks a computed row: its ref is not empty, and a ref
// sha256 <path> or sha256 <path> <prefix>... names a file of the tree or of
// the work area. The step that wrote a hash checks its value again (S06 the
// brief, the runner answers.sha256), so this check does not hash.
func computedFindings(in input, at, ref string) []string {
	if strings.TrimSpace(ref) == "" {
		return []string{at + "no ref"}
	}
	f := strings.Fields(ref)
	if f[0] != "sha256" {
		return nil
	}
	if len(f) < 2 || !isRegular(in.fsys, f[1]) && (in.area == nil || !isRegular(in.area, f[1])) {
		p := ""
		if len(f) > 1 {
			p = f[1]
		}
		return []string{at + "sha256 names " + p + ", which is not a file of the tree or of the work area"}
	}
	return nil
}
