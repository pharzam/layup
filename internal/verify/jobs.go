package verify

import (
	"io/fs"
	"regexp"
	"strings"

	"github.com/pharzam/layup/internal/gate"
)

// checkJobs is check jobs (task T-d6q5, D2 of #92), the evidence of S12: each
// kind of docs/gates.tsv has a job of a workflow .github/workflows/*.yml or
// *.yaml whose check name is the kind, as a ruleset names a required check.
func checkJobs(in input) []string {
	data, err := fs.ReadFile(in.fsys, "docs/gates.tsv")
	if err != nil {
		return []string{"docs/gates.tsv: no manifest"}
	}
	kinds, err := gate.ReadManifest(data)
	if err != nil {
		return []string{"docs/gates.tsv: " + err.Error()}
	}
	names := map[string]bool{}
	for _, pattern := range []string{".github/workflows/*.yml", ".github/workflows/*.yaml"} {
		files, _ := fs.Glob(in.fsys, pattern)
		for _, f := range files {
			if text, err := fs.ReadFile(in.fsys, f); err == nil {
				for _, n := range jobNames(string(text)) {
					names[n] = true
				}
			}
		}
	}
	var out []string
	for _, k := range kinds {
		if !names[k.Name] {
			out = append(out, "no CI job for the kind "+k.Name)
		}
	}
	return out
}

// trailingComment is a comment at the end of a line of YAML.
var trailingComment = regexp.MustCompile(`[ \t]+#.*$`)

// jobNames gives the check name of each job of a workflow: its name:, else its
// id, as the reader of check_protection of LAYUP's setup-check.sh reads a job
// (D2 of #92), with no YAML library (NFR-007). Under the key jobs: at the start
// of a line, a line at the indent of the first line is a job, its id the text
// before its colon; a name: at the indent of the job's first child line is
// its name. A blank line and a comment are skipped, a carriage return at a
// line end is not read, and a comment and a quote around a name are not part
// of it. The reader also takes the quotes off an id, and the white space off
// the end of a name.
func jobNames(text string) []string {
	var out []string
	in, jobIndent, childIndent := false, -1, -1
	job, name := "", ""
	emit := func() {
		switch {
		case name != "":
			out = append(out, name)
		case job != "":
			out = append(out, job)
		}
	}
	for _, l := range strings.Split(text, "\n") {
		l = strings.TrimSuffix(l, "\r")
		if strings.HasPrefix(l, "jobs:") {
			in, jobIndent = true, -1
			continue
		}
		if in && l != "" && l[0] != ' ' && l[0] != '#' {
			in = false
		}
		if t := strings.TrimLeft(l, " \t"); !in || t == "" || t[0] == '#' {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, " "))
		if jobIndent < 0 {
			jobIndent = indent
		}
		if indent == jobIndent {
			emit()
			id, _, _ := strings.Cut(l[indent:], ":")
			job, name, childIndent = unquote(id), "", -1
			continue
		}
		if childIndent < 0 {
			childIndent = indent
		}
		if v, ok := strings.CutPrefix(l[indent:], "name:"); ok && indent == childIndent {
			name = unquote(trailingComment.ReplaceAllString(strings.TrimLeft(v, " \t"), ""))
		}
	}
	emit()
	return out
}

// unquote takes the white space off the end of s, then one quote off each
// end, as check_protection takes the quotes off a name.
func unquote(s string) string {
	s = strings.TrimRight(s, " \t")
	if s != "" && (s[0] == '"' || s[0] == '\'') {
		s = s[1:]
	}
	if s != "" && (s[len(s)-1] == '"' || s[len(s)-1] == '\'') {
		s = s[:len(s)-1]
	}
	return s
}
