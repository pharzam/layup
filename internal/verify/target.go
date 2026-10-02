package verify

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/work"
)

// The four checks of a target's form (D2 to D9 of #89): facts, onboarding,
// glossary and guardrails. Each core is the part of its sh function that a
// target keeps, with the lines of that function, and the harness compares
// those lines on the shared fixtures; LAYUP's counts, headings and records do
// not apply (setup.md, The checks of layup setup verify).

const (
	sumsFile       = "docs/setup/facts.sha256"
	briefFile      = "docs/facts/problem-statement-brief.md"
	visionFile     = "docs/facts/architectural-vision-brief.md"
	factsIndexFile = "docs/facts/README.md"
	onboardingFile = "docs/onboarding-for-engineers.md"
	glossaryFile   = "docs/glossary.md"
	guardrailsFile = "docs/guardrails.md"
)

// isRegular reports whether path is a regular file of fsys, as [ -f ] does.
func isRegular(fsys fs.FS, path string) bool {
	info, err := fs.Stat(fsys, path)
	return err == nil && info.Mode().IsRegular()
}

// sumsLines gives the lines of the list of hashes, as read -r reads them: a
// last line with no line feed is a line too.
func sumsLines(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	return strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
}

// hashFields splits a line of the list as read -r splits it into two names:
// the hash, and the rest of the line, at spaces and tabs.
func hashFields(line string) (want, path string) {
	line = strings.TrimLeft(line, " \t")
	i := strings.IndexAny(line, " \t")
	if i < 0 {
		return line, ""
	}
	return line[:i], strings.Trim(line[i:], " \t")
}

// factsHashFindings is the hash loop of check_facts: each line of the list
// names a file whose SHA-256 it holds.
func factsHashFindings(fsys fs.FS) []string {
	if !isRegular(fsys, sumsFile) {
		return []string{"hash: " + sumsFile + " is absent"}
	}
	data, _ := fs.ReadFile(fsys, sumsFile)
	var out []string
	for i, line := range sumsLines(data) {
		want, path := hashFields(line)
		switch {
		case want == "":
		case path == "":
			out = append(out, fmt.Sprintf("hash: line %d of %s has no path", i+1, sumsFile))
		default:
			b, err := fs.ReadFile(fsys, path)
			if !isRegular(fsys, path) || err != nil || fmt.Sprintf("%x", sha256.Sum256(b)) != want {
				out = append(out, "hash: "+path+" does not match "+sumsFile)
			}
		}
	}
	return out
}

// An answersKind is one of the two answers records (D1, D2 of #89): its file
// name, its step, and the prefixes of its questions.
type answersKind struct {
	slug, step, other string
	prefixes          []string
}

var answersKinds = []answersKind{
	{"setup-answers", "S04", "S11", []string{"S01-", "Q-"}},
	{"marker-answers", "S11", "S04", []string{"M-"}},
}

// checkFacts is check facts of a target (D4 to D6 of #89, with condition 2 of
// its plan review): the hash loop; each brief and each answers record that
// exists is read in full (its line in the list, its index row, its question
// IDs); the done rows of the setup record decide only what must exist. The
// list comes with the record of S04 (O-124), so before S04 its absence is no
// finding.
func checkFacts(in input) []string {
	done := map[string]bool{}
	for _, r := range in.record {
		if r[1] == "done" {
			done[r[0]] = true
		}
	}
	var out []string
	if isRegular(in.fsys, sumsFile) || done["S04"] {
		out = factsHashFindings(in.fsys)
	}
	listed := map[string]bool{}
	sums, _ := fs.ReadFile(in.fsys, sumsFile)
	for _, line := range sumsLines(sums) {
		if want, path := hashFields(line); want != "" && path != "" {
			listed[path] = true
		}
	}
	for _, b := range []string{briefFile, visionFile} {
		switch {
		case isRegular(in.fsys, b) && !listed[b]:
			out = append(out, "listed: "+b+" is not in "+sumsFile)
		case !isRegular(in.fsys, b) && b == briefFile && done["S06"]:
			out = append(out, "brief: "+b+" is absent")
		}
	}
	for _, k := range answersKinds {
		paths, _ := fs.Glob(in.fsys, "docs/facts/F-[0-9][0-9][0-9][0-9]-"+k.slug+".md")
		must := done[k.step] && (k.step == "S04" || slices.ContainsFunc(in.answers, func(r []string) bool { return hasPrefix(r[0], k.prefixes) }))
		if len(paths) != 1 {
			if len(paths) > 1 || must {
				out = append(out, "record: expected one docs/facts/F-NNNN-"+k.slug+".md")
			}
			continue
		}
		p, id := paths[0], strings.TrimPrefix(paths[0], "docs/facts/")[:len("F-NNNN")]
		if !listed[p] {
			out = append(out, "listed: "+p+" is not in "+sumsFile)
		}
		out = append(out, answersFindings(in, p, id, k)...)
		index, _ := fs.ReadFile(in.fsys, factsIndexFile)
		if !regexp.MustCompile(`(?m)^\|[ \t\f\v\r]*\[?` + id + `[^0-9]`).Match(index) {
			out = append(out, "index: "+factsIndexFile+" has no row for "+id)
		}
	}
	return out
}

var factLine = regexp.MustCompile(`^([0-9]+)\. (.*)$`)

// answersFindings reads the facts of an answers record (D2 of #89): each
// names, in its first code span, a question of its kind that answers.tsv has,
// once; and each row of that kind has a fact.
func answersFindings(in input, path, id string, k answersKind) []string {
	data, _ := fs.ReadFile(in.fsys, path)
	other := answersKinds[0]
	if k.step == other.step {
		other = answersKinds[1]
	}
	asked := map[string]bool{}
	for _, r := range in.answers {
		asked[r[0]] = true
	}
	var out, order []string
	count := map[string]int{}
	for _, line := range strings.Split(string(data), "\n") {
		m := factLine.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		n, _ := strconv.Atoi(m[1])
		q := ""
		if rest, ok := strings.CutPrefix(m[2], "`"); ok {
			q, _, ok = strings.Cut(rest, "`")
			if !ok {
				q = ""
			}
		}
		switch {
		case q == "":
			out = append(out, fmt.Sprintf("answers: %s fact %d names no question ID", id, n))
			continue
		case hasPrefix(q, other.prefixes):
			out = append(out, fmt.Sprintf("answers: %s fact %d names %s, a question of the record of %s", id, n, q, other.step))
		case !hasPrefix(q, k.prefixes) || !asked[q]:
			out = append(out, fmt.Sprintf("answers: %s fact %d names %s, which %s does not have", id, n, q, work.AnswersPath))
		}
		if count[q] == 0 {
			order = append(order, q)
		}
		count[q]++
	}
	for _, q := range order {
		if count[q] > 1 {
			out = append(out, fmt.Sprintf("answers: %s holds %s %d times", id, q, count[q]))
		}
	}
	for _, r := range in.answers {
		if hasPrefix(r[0], k.prefixes) && count[r[0]] == 0 {
			out = append(out, fmt.Sprintf("answers: %s of %s is not a fact of %s", r[0], work.AnswersPath, id))
		}
	}
	return out
}

// hasPrefix reports whether q has one of the prefixes.
func hasPrefix(q string, prefixes []string) bool {
	return slices.ContainsFunc(prefixes, func(p string) bool { return strings.HasPrefix(q, p) })
}

var citation = regexp.MustCompile(`F-[0-9]{4}#[0-9]+`)

// citations gives each F-NNNN#n of the text once, in byte order (the sort -u
// of the sh functions).
func citations(data []byte) []string {
	var out []string
	for _, c := range citation.FindAll(data, -1) {
		out = append(out, string(c))
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// resolveFact gives the finding of a citation F-NNNN#n that does not resolve,
// or "": it resolves when docs/facts/ holds exactly one F-NNNN-*.md with a
// numbered fact n (D7 of #89). check_onboarding reads the first such record.
func resolveFact(fsys fs.FS, cite string) string {
	id, n, _ := strings.Cut(cite, "#")
	paths, _ := fs.Glob(fsys, "docs/facts/"+id+"-*.md")
	if len(paths) == 1 {
		data, err := fs.ReadFile(fsys, paths[0])
		if err == nil && regexp.MustCompile(`(?m)^0*`+regexp.QuoteMeta(n)+`\. `).Match(data) {
			return ""
		}
	}
	return "fact: " + cite + " is not a fact of the " + id + " record"
}

// citationFindings gives the finding of each citation of data that does not
// resolve.
func citationFindings(fsys fs.FS, data []byte) []string {
	var out []string
	for _, c := range citations(data) {
		if f := resolveFact(fsys, c); f != "" {
			out = append(out, f)
		}
	}
	return out
}

// onboardingFindings is check_onboarding in a target's form (D8 of #89): the
// file, no marker, the link to the problem statement, and each citation.
func onboardingFindings(fsys fs.FS) []string {
	if !isRegular(fsys, onboardingFile) {
		return []string{"missing: " + onboardingFile + " is absent"}
	}
	data, _ := fs.ReadFile(fsys, onboardingFile)
	var out []string
	if bytes.Contains(data, []byte("\u2039")) {
		out = append(out, "marker: "+onboardingFile+" holds a \u2039 character")
	}
	if !bytes.Contains(data, []byte("](facts/problem-statement-brief.md)")) {
		out = append(out, "link: "+onboardingFile+" has no link to facts/problem-statement-brief.md")
	}
	return append(out, citationFindings(fsys, data)...)
}

// glossaryFindings is check glossary in a target's form (D8 of #89): each
// citation resolves; any heading and any number of rows.
func glossaryFindings(fsys fs.FS) []string {
	if !isRegular(fsys, glossaryFile) {
		return []string{"missing: " + glossaryFile + " is absent"}
	}
	data, _ := fs.ReadFile(fsys, glossaryFile)
	return citationFindings(fsys, data)
}

var guardrailsEntry = regexp.MustCompile(`^- \*\*Inv-([0-9]+)\*\*`)

// guardrailsFindings is check guardrails in a target's form (D8 of #89): an
// entry is a bullet "- **Inv-N**" up to the next entry, a heading or the end
// of the file; its Check: value is read as check_guardrails reads it; and each
// citation of the file resolves.
func guardrailsFindings(fsys fs.FS) []string {
	if !isRegular(fsys, guardrailsFile) {
		return []string{"missing: " + guardrailsFile + " is absent"}
	}
	data, _ := fs.ReadFile(fsys, guardrailsFile)
	type entry struct {
		n     string
		value string // the first Check: value, after the last "Check: " of its line
		found bool
	}
	var entries []entry
	in := false
	for _, line := range strings.Split(string(data), "\n") {
		if m := guardrailsEntry.FindStringSubmatch(line); m != nil {
			entries, in = append(entries, entry{n: m[1]}), true
		} else if strings.HasPrefix(line, "#") {
			in = false
		}
		if !in {
			continue
		}
		e := &entries[len(entries)-1]
		if i := strings.LastIndex(line, "Check: "); i >= 0 && !e.found {
			e.value, e.found = strings.TrimRight(line[i+len("Check: "):], " \t\n\v\f\r"), true
		}
	}
	var out []string
	for _, e := range entries {
		switch v := e.value; {
		case v == "no check yet":
		case strings.Contains(v, " (") && strings.HasSuffix(v, ")"):
			path, gate := v[:strings.Index(v, " (")], strings.TrimSuffix(v[strings.LastIndex(v, " (")+2:], ")")
			if path == "" || !isRegular(fsys, path) {
				out = append(out, "check: Inv-"+e.n+" names \""+path+"\", which is not a file")
			}
			if gate != "hook" && !(strings.HasPrefix(gate, "ci:") && len(gate) > len("ci:")) {
				out = append(out, "check: Inv-"+e.n+" gate \""+gate+"\" is not hook or ci:<job>")
			}
		default:
			out = append(out, "check: Inv-"+e.n+" has no valid Check: value (\""+v+"\")")
		}
	}
	return append(out, citationFindings(fsys, data)...)
}

// checkOnboarding, checkGlossary and checkGuardrails are the three checks of
// a target: the cores alone.
func checkOnboarding(in input) []string { return onboardingFindings(in.fsys) }
func checkGlossary(in input) []string   { return glossaryFindings(in.fsys) }
func checkGuardrails(in input) []string { return guardrailsFindings(in.fsys) }
