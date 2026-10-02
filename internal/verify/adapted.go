package verify

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/git"
)

// adExclude is AD_EXCLUDE of setup-check.sh: the paths that check adapted
// does not read. O-123 keeps the rule of the check as it is, so a target's
// check reads LAYUP's lists (D1 of #88); TestTheListsOfAdaptedEqualTheSh
// compares them with the sh file.
const adExclude = `^(docs/facts/|runs/|docs/tasks/T-[^/]*\.md$|docs/tasks/completed\.md$|docs/adr/00(0[1-9]|1[0-2])-[^/]*\.md$|docs/setup/record-[^/]*\.md$|docs/(adr|ci|links|prd|setup)/tests/|\.githooks/tests/|internal/psb/testdata/)`

var adExcludeRE = regexp.MustCompile(adExclude)

// adAllowed is ad_allowed of setup-check.sh: each file that may name
// Armature (rule 2), with its reason.
var adAllowed = [][2]string{
	{"README.md", "the pinned baseline (docs/setup/armature.pin, ADR-0009)"},
	{"AGENTS.md", "the pinned baseline (docs/setup/armature.pin, ADR-0009)"},
	{"docs/setup/README.md", "the setup of the pinned baseline and its record"},
	{"docs/glossary.md", "the term Armature and the PSB terms that name it"},
	{"docs/prd/PRD-0001-layup.md", "the product parts that work on Armature"},
	{"docs/adr/README.md", "the index lists the title of ADR-0009"},
	{"docs/guardrails.md", "Inv-8 restates F-0001#8 (check guardrails reads that entry)"},
	{"docs/onboarding-for-engineers.md", "Armature is LAYUP's baseline (F-0001#19)"},
	{"docs/engineering-discipline.md", "How this project was set up: the pinned baseline"},
}

// An adHit is one of the 16 patterns of check_adapted: rule 2 reads the text
// as it is, each other pattern the text in lower case. Group 1 is the
// boundary before the word.
type adHit struct {
	rule, name string
	re         *regexp.Regexp
	cased      bool
}

var adHits = []adHit{
	{"rule-1", "kit", regexp.MustCompile(`(^|[^a-z0-9_])kits?([^a-z0-9_]|$)`), false},
	{"rule-1", "adopter", regexp.MustCompile(`(^|[^a-z0-9_])adopters?([^a-z0-9_]|$)`), false},
	{"rule-1", "the template", regexp.MustCompile(`(^|[^a-z0-9_])(the|this) templates?([^a-z0-9_-]|$)`), false},
	{"rule-2", "Armature", regexp.MustCompile(`(^|[^A-Za-z0-9_])Armature([^A-Za-z0-9_]|$)`), true},
	{"rule-3", "optional", regexp.MustCompile(`(^|[^a-z0-9_])optional([^a-z0-9_]|$)`), false},
	{"rule-3", "skip this section", regexp.MustCompile(`(^|[^a-z0-9_])skips? this section`), false},
	{"rule-3", "fill \u2039", regexp.MustCompile(`(^|[^a-z0-9_#-])fill[a-z]*[^.\x01]*\x01`), false},
	{"rule-3", "replace \u2039", regexp.MustCompile(`(^|[^a-z0-9_#-])replac[a-z]*[^.\x01]*\x01`), false},
	{"rule-3", "fill in", regexp.MustCompile(`(^|[^a-z0-9_#-])fill in([^a-z0-9_-]|$)`), false},
	{"rule-3", "delete this", regexp.MustCompile(`(^|[^a-z0-9_])delete this([^a-z0-9_]|$)`), false},
	{"rule-3", "delete the one you do not use", regexp.MustCompile(`(^|[^a-z0-9_])delete the ones? you do not use`), false},
	{"rule-3", "adapt", regexp.MustCompile(`(^|[^a-z0-9_])adapts?([^a-z0-9_]|$)`), false},
	{"rule-3", "your project", regexp.MustCompile(`(^|[^a-z0-9_])your project`), false},
	{"rule-3", "your forge", regexp.MustCompile(`(^|[^a-z0-9_])your forge`), false},
	{"rule-3", "your stack", regexp.MustCompile(`(^|[^a-z0-9_])your stack`), false},
	{"rule-3", "you use", regexp.MustCompile(`(^|[^a-z0-9_])you use([^a-z0-9_]|$)`), false},
}

// adaptedFindings is check_adapted (D2 to D5 of #88): the hits of the tracked
// .md files of the tree, each finding once, in byte order.
func adaptedFindings(fsys fs.FS, h history) []string {
	files, err := h.Files()
	if err != nil {
		return []string{"git: cannot list the tracked files"}
	}
	var out []string
	for _, f := range adaptedFiles(fsys, files) {
		out = append(out, f.findings...)
	}
	slices.Sort(out)
	return slices.Compact(out)
}

// Flagged gives the path of each file of the work tree at tree that check
// adapted flags, once, in byte order: the files for which the prose step
// stops (O-123, K12; D6 of #88).
func Flagged(tree string) ([]string, error) {
	files, err := git.LsFiles(tree)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, f := range adaptedFiles(os.DirFS(tree), files) {
		out = append(out, f.path)
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}

// A flaggedFile is a file with its findings.
type flaggedFile struct {
	path     string
	findings []string
}

// adaptedFiles gives each file that check adapted reads and flags: a tracked
// path that ends in .md, that AD_EXCLUDE does not match, and that is a
// regular file when links are followed, as [ -f ] reads it.
func adaptedFiles(fsys fs.FS, files []string) []flaggedFile {
	var out []flaggedFile
	for _, p := range files {
		if !strings.HasSuffix(p, ".md") || adExcludeRE.MatchString(p) {
			continue
		}
		info, err := fs.Stat(fsys, p)
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			continue
		}
		allowed := slices.ContainsFunc(adAllowed, func(a [2]string) bool { return a[0] == p })
		if f := fileHits(p, data, allowed); len(f) > 0 {
			out = append(out, flaggedFile{p, f})
		}
	}
	return out
}

// fileHits gives the findings of one file, paragraph by paragraph, as the awk
// program of check_adapted does in the C locale (D3 of #88): bytes, not
// characters.
func fileHits(path string, data []byte, allowed bool) []string {
	var out []string
	var para []byte
	var starts, lines []int // the start of each line in para, and its number
	flush := func() {
		if len(starts) == 0 {
			return
		}
		out = append(out, paragraphHits(path, string(para), starts, lines, allowed)...)
		para, starts, lines = nil, nil, nil
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSuffix(line, "\r")
		line = strings.Trim(blanks.ReplaceAllString(line, " "), " ")
		if line == "" {
			flush()
			continue
		}
		var mk strings.Builder // a marker is one unit, the byte 0x01
		for {
			j := strings.Index(line, "\u2039")
			if j < 0 {
				break
			}
			rest := line[j+len("\u2039"):]
			mk.WriteString(line[:j])
			mk.WriteByte(1)
			line = ""
			if k := strings.Index(rest, "\u203a"); k >= 0 {
				line = rest[k+len("\u203a"):]
			}
		}
		line = mk.String() + line
		if len(starts) > 0 {
			para = append(para, ' ')
		}
		starts, lines = append(starts, len(para)), append(lines, i+1)
		para = append(para, line...)
	}
	flush()
	return out
}

var blanks = regexp.MustCompile(`[ \t]+`)

// paragraphHits gives the hits of the 16 patterns in one paragraph, each at
// the line where its word starts (D4 of #88).
func paragraphHits(path, para string, starts, lines []int, allowed bool) []string {
	low := []byte(para) // the lower case changes only A to Z, so a position stays
	for i, b := range low {
		if 'A' <= b && b <= 'Z' {
			low[i] = b + 'a' - 'A'
		}
	}
	var out []string
	for _, h := range adHits {
		text := string(low)
		if h.cased {
			if allowed {
				continue
			}
			text = para
		}
		for off := 0; off < len(text); {
			loc := h.re.FindStringSubmatchIndex(text[off:])
			if loc == nil {
				break
			}
			word := off + loc[3] // the end of the boundary: the first letter of the word
			if h.name != "kit" || !(bytes.HasPrefix(low[word:], []byte("kit-history")) || bytes.HasPrefix(low[word:], []byte("kit-linters"))) {
				line := lines[0]
				for k, s := range starts {
					if s <= word {
						line = lines[k]
					}
				}
				out = append(out, fmt.Sprintf("%s %s: %s:%d", h.rule, h.name, path, line))
			}
			off += loc[0] + 1
		}
	}
	return out
}
