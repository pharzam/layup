package verify

import (
	"bytes"
	"fmt"
	"io/fs"
	"regexp"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/work"
)

// The cores of the checks: each reads a repository root and gives its
// findings, each the text after "FAIL " of a line of setup-check.sh, in the
// order of its sh function (D2 of #84). No finding is a pass. The harness runs
// them on the fixtures of docs/setup/tests/.

// history is the git part of a repository that check pin reads.
type history interface {
	IsShallow() (bool, error)
	RootCommits() ([]string, error) // of HEAD
	Tree(commit string) (string, error)
}

const pinPath = "docs/setup/armature.pin"

var fullCommit = regexp.MustCompile(`^[0-9a-f]{40}$`)

// pinValues gives the lines of the pin file by key, in the order of the file:
// the text after "key=" of each line that starts with it.
func pinValues(data []byte) map[string][]string {
	values := map[string][]string{}
	for _, line := range strings.Split(string(data), "\n") {
		if key, value, ok := strings.Cut(line, "="); ok {
			values[key] = append(values[key], value)
		}
	}
	return values
}

// first gives the first value of key, as `sed -n "s/^key=//p" | head -1`.
func first(values map[string][]string, key string) string {
	if v := values[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

// pinFindings is check_pin: the pin names the commit and the tree of the
// baseline, and the tree is that of the one root commit.
func pinFindings(fsys fs.FS, h history) []string {
	data, err := fs.ReadFile(fsys, pinPath)
	if err != nil {
		return []string{"missing: docs/setup/armature.pin is absent"}
	}
	var out []string
	values := pinValues(data)
	for _, key := range []string{"source", "commit", "tree", "method", "date"} {
		if n := len(values[key]); n != 1 {
			out = append(out, fmt.Sprintf("key: %s appears %d times, expected 1", key, n))
		}
	}
	commit, tree := first(values, "commit"), first(values, "tree")
	if !fullCommit.MatchString(commit) {
		out = append(out, "commit: not 40 hexadecimal characters: "+commit)
	}
	if shallow, err := h.IsShallow(); err == nil && shallow {
		return append(out, "shallow: the clone is shallow, so the root commit is not visible (fetch the full history)")
	}
	roots, _ := h.RootCommits()
	if len(roots) != 1 {
		return append(out, fmt.Sprintf("tree: expected one root commit, found %d", len(roots)))
	}
	if rootTree, _ := h.Tree(roots[0]); tree != rootTree {
		out = append(out, fmt.Sprintf("tree: pin names %s, root commit has %s", tree, rootTree))
	}
	return out
}

// kitHistoryFindings is check_kit_history: the kit's own history is gone, each
// task file has a line in a task index, and no task index links the baseline,
// whose repository is repo ("" reads no link). A link is repo followed by the
// end of the text or by a character that is not a letter, a digit, "-" or
// "_", so a link to the repository itself counts, and a link to a repository
// whose name only starts with repo does not (round 1 of #84, finding 1); the
// sh function reads only a link under repo.
func kitHistoryFindings(fsys fs.FS, repo string) []string {
	var out []string
	for _, dir := range []string{"decisions", "audit"} {
		if _, err := fs.Stat(fsys, "docs/"+dir); err == nil {
			out = append(out, fmt.Sprintf("%s: docs/%s/ exists (kit step 4 deletes it)", dir, dir))
		}
	}
	backlog, _ := fs.ReadFile(fsys, "docs/tasks/backlog.md")
	completed, _ := fs.ReadFile(fsys, "docs/tasks/completed.md")
	indexes := append(append([]byte{}, backlog...), completed...)
	tasks, _ := fs.Glob(fsys, "docs/tasks/T-*.md")
	for _, t := range tasks {
		if info, err := fs.Stat(fsys, t); err != nil || !info.Mode().IsRegular() {
			continue
		}
		id := strings.TrimSuffix(strings.TrimPrefix(t, "docs/tasks/"), ".md")
		if !bytes.Contains(indexes, []byte(id)) {
			out = append(out, fmt.Sprintf("orphan: docs/tasks/%s.md has no line with %s in backlog.md or completed.md", id, id))
		}
	}
	link := regexp.MustCompile(regexp.QuoteMeta(repo) + `($|[^A-Za-z0-9_-])`)
	for _, index := range []string{"backlog", "completed"} {
		data, err := fs.ReadFile(fsys, "docs/tasks/"+index+".md")
		if repo != "" && err == nil && link.Match(data) {
			out = append(out, fmt.Sprintf("kit-link: docs/tasks/%s.md links %s/ (a kit task or note)", index, repo))
		}
	}
	return out
}

// kitPhrases are the phrases with which the kit's entry files say that the
// repository is the kit.
var kitPhrases = []string{"Agent context for **Armature**", "This repository is a generic **template**", "A domain-free **template**, not a product"}

// identityFindings is check_identity: the entry files do not say that the
// repository is the kit, and README.md links the pin.
func identityFindings(fsys fs.FS) []string {
	var out []string
	for _, f := range []string{"README.md", "AGENTS.md"} {
		data, err := fs.ReadFile(fsys, f)
		if err != nil {
			continue
		}
		for _, p := range kitPhrases {
			if bytes.Contains(data, []byte(p)) {
				out = append(out, fmt.Sprintf("kit: %s says the repository is the Armature kit (\"%s\")", f, p))
			}
		}
	}
	if readme, err := fs.ReadFile(fsys, "README.md"); err != nil || !bytes.Contains(readme, []byte("](docs/setup/armature.pin)")) {
		out = append(out, "pin: README.md has no link to docs/setup/armature.pin")
	}
	return out
}

// The target's part of the three checks (D3 of #84): the frame adds it to the
// core, with the values of the setup record. A row that a check needs and
// that the record does not have, or has with no value, is a finding of that
// check, so a check of a step reads the record as it stands at that step.

// value gives the value of the record row of step and name, or a finding.
func value(r work.Record, step, name string) (string, string) {
	if v, ok := r.Value(step, name); ok && v != "" {
		return v, ""
	}
	return "", fmt.Sprintf("record: no value at %s %s", step, name)
}

// timeForm is the form of pin.time, which S02 writes: UTC, to the second.
const timeForm = "2006-01-02T15:04:05Z"

// pinTarget is the target's part of check pin: the source, the commit and the
// tree of the pin are the record rows of S02, its date is the date of
// pin.time, and its method is the text of NFR-006. A missing pin file gives
// nothing here: the core says it.
func pinTarget(fsys fs.FS, r work.Record) []string {
	data, err := fs.ReadFile(fsys, pinPath)
	if err != nil {
		return nil
	}
	values := pinValues(data)
	var out []string
	rec := map[string]string{}
	for _, key := range []string{"source", "commit", "tree"} {
		v, missing := value(r, "S02", "pin."+key)
		switch {
		case missing != "":
			out = append(out, missing)
		case first(values, key) != v:
			out = append(out, fmt.Sprintf("%s: the pin names %s, the record row pin.%s is %s", key, first(values, key), key, v))
		}
		rec[key] = v
	}
	switch t, missing := value(r, "S02", "pin.time"); {
	case missing != "":
		out = append(out, missing)
	case !validTime(t):
		out = append(out, "record: pin.time is not of the form YYYY-MM-DDTHH:MM:SSZ: "+t)
	case first(values, "date") != t[:10]:
		out = append(out, fmt.Sprintf("date: the pin names %s, the date of pin.time is %s", first(values, "date"), t[:10]))
	}
	if rec["source"] != "" && rec["commit"] != "" {
		want := "git clone " + rec["source"] + ", checkout " + rec["commit"] + ", .git removed"
		if m := first(values, "method"); m != want {
			out = append(out, fmt.Sprintf("method: the pin names %s, not %s", m, want))
		}
	}
	return out
}

func validTime(t string) bool {
	_, err := time.Parse(timeForm, t)
	return err == nil && len(t) == len(timeForm)
}

// kitLink gives the repository of the kit-link rule for a baseline: its URL
// with no scheme, and no "/" and no ".git" at its end. "" gives "".
func kitLink(source string) string {
	if _, rest, ok := strings.Cut(source, "://"); ok {
		source = rest
	}
	return strings.TrimRight(strings.TrimSuffix(strings.TrimRight(source, "/"), ".git"), "/")
}

// identityTarget is the target's part of check identity: README.md holds the
// name of the target, the record row name of S01.
func identityTarget(fsys fs.FS, r work.Record) []string {
	name, missing := value(r, "S01", "name")
	if missing != "" {
		return []string{missing}
	}
	if readme, err := fs.ReadFile(fsys, "README.md"); err != nil || !bytes.Contains(readme, []byte(name)) {
		return []string{"name: README.md does not hold the name of the target, " + name}
	}
	return nil
}

// checkPin is check pin of a target: the core, then the target's part.
func checkPin(in input) []string {
	return append(pinFindings(in.fsys, in.repo), pinTarget(in.fsys, in.record)...)
}

// checkKitHistory is check kit-history of a target: the core, with the link
// text of the target's baseline.
func checkKitHistory(in input) []string {
	source, missing := value(in.record, "S02", "pin.source")
	out := kitHistoryFindings(in.fsys, kitLink(source))
	if missing != "" {
		out = append(out, missing)
	}
	return out
}

// checkIdentity is check identity of a target: the core, then the name.
func checkIdentity(in input) []string {
	return append(identityFindings(in.fsys), identityTarget(in.fsys, in.record)...)
}
