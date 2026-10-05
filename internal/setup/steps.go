package setup

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/catalog"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// BriefPath is the problem statement of a work area, from its root, which
// internal/cli reads for the steps, and VisionPath the vision brief, when
// there is one (D4 of #90).
const (
	BriefPath  = work.BriefPath
	VisionPath = "inputs/briefs/vision.md"
)

// A Brief is what internal/cli hands to the steps (D5 of #86): the gap table
// of the problem statement, as internal/psb writes it, and the SHA-256 of the
// bytes that the rules read.
type Brief struct {
	Gaps []byte
	Sum  string
}

// Checks runs the named checks of layup setup verify on the work area at dir,
// and gives the reason of the first row that is not pass or clear, or "": the
// one-check call of docs/spec/setup.md, which internal/cli makes.
type Checks func(dir string, names []string) string

// A Marker is one place of a marker in a tracked file of the target: its file,
// its line, the byte column of its open quote on that line, and its text with
// its angle quotes, as the scanner of internal/verify gives it (D7 of #90).
type Marker struct {
	File      string
	Line, Col int
	Text      string
}

// Calls are the calls of internal/verify that internal/cli gives the steps
// (D9 of #90), as internal/setup does not import internal/verify: the
// one-check call of the evidence (D12 of #86); the markers of a tree (S10,
// S11); the files whose links break (S05); the files that check adapted flags
// (S14); and the link rule of check kit-history (S05). A nil call gives a step
// no evidence call, or makes the step that needs it fail.
type Calls struct {
	Checks        Checks
	Markers       func(tree string) ([]Marker, error)
	BrokenLinks   func(tree string) ([]string, error)
	Flagged       func(tree string) ([]string, error)
	LinksBaseline func(source, line string) bool
	LostMarkers   func(path string, before, after []byte) []string
}

// Steps gives the steps of this version of layup: S01 to S04 (task T-7s0y,
// #86), S05 to S11 and S14 (task T-b3r1, #90), and S12, S13 and S15 (task
// T-d6q5, #92).
func Steps(b Brief, c Calls) map[string]Step {
	m := Stubs()
	set := func(id string, run func(Input) Outcome, checks ...string) {
		s := m[id]
		s.Run = run
		if c.Checks != nil && len(checks) > 0 {
			s.Evidence = func(dir string) string { return c.Checks(dir, checks) }
		}
		m[id] = s
	}
	set("S01", func(in Input) Outcome { return runS01Step(b, in) })
	s := m["S01"]
	s.Unchanged = func(r work.Record) error { return briefUnchanged(b, r) }
	m["S01"] = s
	set("S02", runS02)
	set("S03", runS03)
	s = m["S03"]
	s.Commands = commandsS03
	m["S03"] = s
	set("S04", func(in Input) Outcome { return runS04(b, in) }, "pin", "facts")
	set("S05", func(in Input) Outcome { return runS05(c, in) }, "kit-history", "link-lint")
	set("S06", runS06, "facts")
	set("S07", func(in Input) Outcome { return runProseStep("S07", c, in) }, "onboarding")
	set("S08", func(in Input) Outcome { return runProseStep("S08", c, in) }, "glossary")
	set("S09", func(in Input) Outcome { return runProseStep("S09", c, in) }, "guardrails")
	set("S14", func(in Input) Outcome { return runProseStep("S14", c, in) }, "identity", "adapted")
	set("S10", func(in Input) Outcome { return runS10(c, in) })
	set("S11", func(in Input) Outcome { return runS11(c, in) }, "markers", "sources", "facts")
	set("S12", runS12, "jobs", "gates")
	set("S13", runS13)
	s = m["S13"]
	s.Commands = commandsS13
	m["S13"] = s
	set("S15", runS15)
	s = m["S15"]
	s.Commands, s.Records = commandsS15, recordsS15
	m["S15"] = s
	return m
}

// GapsSchema is the form of the gap table that internal/cli hands to S01: the
// block psb-gaps of docs/spec/psb-check.md. internal/psb writes the table, and
// this package reads it with internal/tsv (D5 of #86), so the two packages
// each hold the schema, and each test compares it with the one block.
var GapsSchema = tsv.Schema{Name: "psb-gaps", Location: "stdout", Columns: []tsv.Column{
	{Name: "id", Type: "id(Q-NNN)", Key: true},
	{Name: "rule", Type: "enum(G1|G2|G3|G4|G5)"},
	{Name: "line", Type: "int"},
	{Name: "excerpt", Type: "text"},
	{Name: "question", Type: "text"},
}}

// A Gap is a row of the gap table: a question for the idea owner.
type Gap struct {
	ID       string
	Line     int
	Question string
}

// readGaps reads the gap table by GapsSchema.
func readGaps(table []byte) ([]Gap, error) {
	rows, err := tsv.Read(table, GapsSchema)
	if err != nil {
		return nil, fmt.Errorf("the gap table of %s: %w", work.BriefPath, err)
	}
	gaps := []Gap{}
	for _, r := range rows {
		line, _ := strconv.Atoi(r[2])
		gaps = append(gaps, Gap{ID: r[0], Line: line, Question: r[4]})
	}
	return gaps, nil
}

var (
	ownerForm  = regexp.MustCompile(`^[A-Za-z0-9](?:-?[A-Za-z0-9])*$`)
	nameForm   = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)
	urlForm    = regexp.MustCompile(`^(https?|git|file)://[^\s]+$`)
	factSource = regexp.MustCompile(`^F-[0-9]{4}#[0-9]+$`)
)

// ownerName reports whether s is OWNER/NAME of a repository on GitHub (D1 of
// #86): OWNER of letters, digits and single -, not at an end, at most 39
// characters; NAME of letters, digits, ., _ and -, at most 100, not . or ...
func ownerName(s string) bool {
	owner, name, ok := strings.Cut(s, "/")
	return ok && len(owner) <= 39 && ownerForm.MatchString(owner) &&
		len(name) <= 100 && nameForm.MatchString(name) && name != "." && name != ".."
}

// baselineURL reports whether s is a URL whose scheme internal/git reads with
// no login: https, http, git or file (K31; no ssh).
func baselineURL(s string) bool { return urlForm.MatchString(s) }

// runS01Step is S01 (D1 to D4 of #86): one stop table for each missing answer,
// of both kinds; an input error for an answer to no gap; the checks of the
// answers; and the record rows.
func runS01Step(b Brief, in Input) Outcome {
	gaps, err := readGaps(b.Gaps)
	if err != nil {
		return Outcome{Kind: Invalid, Evidence: err.Error()}
	}
	asked := []string{}
	for _, g := range gaps {
		asked = append(asked, g.ID)
	}
	if err := CheckAsked(in.Answers, []string{"Q-"}, asked); err != nil {
		return Outcome{Kind: Invalid, Evidence: err.Error()}
	}
	answer := map[string][]string{}
	for _, r := range in.Answers {
		if r[1] != "" {
			answer[r[0]] = r
		}
	}
	var stops []StopRow
	for _, q := range work.S01Questions {
		if answer[q.ID] == nil {
			stops = append(stops, StopRow{"S01", q.ID, q.Text, ""})
		}
	}
	for _, g := range gaps {
		if answer[g.ID] == nil {
			where := ""
			if g.Line > 0 {
				where = fmt.Sprintf("%s:%d", work.BriefPath, g.Line)
			}
			stops = append(stops, StopRow{"S01", g.ID, g.Question, where})
		}
	}
	if stops != nil {
		return Outcome{Kind: Stop, Stops: stops}
	}
	value := func(q string) string { return answer[q][1] }
	switch stack, name, vis, url := value("S01-stack"), value("S01-name"), value("S01-visibility"), value("S01-baseline"); {
	case !hasEntry(stack):
		return Outcome{Kind: Fail, Evidence: "the stack " + stack + " has no entry in the catalog"}
	case !ownerName(name):
		return Outcome{Kind: Fail, Evidence: "S01-name: " + name + " is not OWNER/NAME"}
	case vis != "public" && vis != "private":
		return Outcome{Kind: Fail, Evidence: "S01-visibility: " + vis + " is not public or private"}
	case !baselineURL(url):
		return Outcome{Kind: Fail, Evidence: "S01-baseline: " + url + " is not a URL with the scheme https, http, git or file"}
	}
	var values [][]string
	for _, q := range []struct{ name, id string }{{"stack", "S01-stack"}, {"name", "S01-name"}, {"visibility", "S01-visibility"}, {"baseline", "S01-baseline"}} {
		values = append(values, append([]string{q.name, value(q.id)}, sourceOf(answer[q.id])...))
	}
	values = append(values, []string{"brief.sha256", b.Sum, "computed", "sha256 " + work.BriefPath})
	return Outcome{Kind: Done, Evidence: "every answer present; the stack has a catalog entry", Values: values}
}

// hasEntry reports whether the catalog of the binary has an entry for stack.
var hasEntry = func(stack string) bool {
	_, err := catalog.Read(catalog.Embedded(), stack)
	return err == nil
}

// sourceOf gives the source and the ref of the record row of an answer (D3 of
// #86): fact and the citation when the answer's source is F-NNNN#n, else
// answer and the question ID.
func sourceOf(r []string) []string {
	if factSource.MatchString(r[3]) {
		return []string{"fact", r[3]}
	}
	return []string{"answer", r[0]}
}

// briefUnchanged is the check of a done S01 (D6 of #86): the problem statement
// that the rules read now has the SHA-256 of the row S01 brief.sha256.
func briefUnchanged(b Brief, r work.Record) error {
	v, ok := r.Value("S01", "brief.sha256")
	switch {
	case !ok:
		return fmt.Errorf("%s: S01 is done and has no row brief.sha256", work.RecordPath)
	case v != b.Sum:
		return fmt.Errorf("%s: an input that changed after a step read it: S01 read it with the SHA-256 %s", work.BriefPath, v)
	}
	return nil
}

// timeForm is the form of pin.time (task T-6x75).
const timeForm = "2006-01-02T15:04:05Z"

// runS02 is S02 (D7, D8 of #86): it resolves the commit of the baseline once,
// clones it into WORK/target.part, checks it out, removes .git, and renames the
// directory to WORK/target at its end.
func runS02(in Input) Outcome {
	url, _ := in.Record.Value("S01", "baseline")
	src, ref := "answer", "S01-baseline"
	for _, r := range in.Record {
		if r[0] == "S01" && r[1] == "baseline" {
			src, ref = r[3], r[4]
		}
	}
	target := filepath.Join(in.Dir, work.TargetPath)
	part := target + ".part"
	if sys.exists(target) {
		return Outcome{Kind: Invalid, Evidence: work.TargetPath + " of the work area exists, and S02 is not done: remove it, or start again in a new work area"}
	}
	if err := sys.removeAll(part); err != nil {
		return Outcome{Kind: Fail, Evidence: "the clone: " + firstLine(err)}
	}
	commit, err := sys.lsRemote(url, "HEAD")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "git ls-remote " + url + " HEAD: " + firstLine(err)}
	}
	at := sys.now().UTC().Format(timeForm)
	if err := sys.clone(url, part); err != nil {
		return Outcome{Kind: Fail, Evidence: "git clone " + url + ": " + firstLine(err)}
	}
	if err := sys.checkoutDetach(part, commit); err != nil {
		return Outcome{Kind: Fail, Evidence: "git checkout " + commit + ": " + firstLine(err)}
	}
	tree, err := sys.revParse(part, commit+"^{tree}")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the tree of " + commit + ": " + firstLine(err)}
	}
	if err := sys.removeAll(filepath.Join(part, ".git")); err != nil {
		return Outcome{Kind: Fail, Evidence: "the clone: " + firstLine(err)}
	}
	if err := sys.rename(part, target); err != nil {
		return Outcome{Kind: Fail, Evidence: "the clone: " + firstLine(err)}
	}
	return Outcome{Kind: Done, Evidence: "the commit and the tree", Values: [][]string{
		{"pin.source", url, src, ref},
		{"pin.commit", commit, "computed", "git ls-remote " + url + " HEAD"},
		{"pin.tree", tree, "computed", "git rev-parse " + commit + "^{tree}"},
		{"pin.time", at, "computed", "the clock of the LAYUP host"},
	}}
}

// runS03 is S03 (D8 of #86): the root commit of the unmodified copy on main, by
// the identity of the run at pin.time; its tree must be pin.tree. A target
// whose main is one commit with the tree pin.tree and the message of S03 is
// the commit of a run that stopped, and S03 takes it.
func runS03(in Input) Outcome {
	target := filepath.Join(in.Dir, work.TargetPath)
	commit, _ := in.Record.Value("S02", "pin.commit")
	tree, _ := in.Record.Value("S02", "pin.tree")
	at, _ := in.Record.Value("S02", "pin.time")
	when, err := time.Parse(timeForm, at)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the record has no pin.time of the form " + timeForm}
	}
	message := "chore: the unmodified baseline at " + commit
	if sys.exists(filepath.Join(target, ".git")) {
		head, err := sys.revParse(target, "refs/heads/main^{commit}")
		roots, rerr := sys.rootCommits(target, "refs/heads/main")
		got, terr := sys.revParse(target, head+"^{tree}")
		msg, merr := sys.message(target, head)
		if err != nil || rerr != nil || terr != nil || merr != nil || !slices.Equal(roots, []string{head}) || got != tree || msg != message {
			return Outcome{Kind: Invalid, Evidence: work.TargetPath + " of the work area has a history that S03 did not make: start again in a new work area"}
		}
		return Outcome{Kind: Done, Evidence: "root tree = pin.tree"}
	}
	who := in.Who
	who.Time = when
	if err := sys.initRepo(target); err != nil {
		return Outcome{Kind: Fail, Evidence: "git init: " + firstLine(err)}
	}
	if err := sys.commit(target, message, who); err != nil {
		return Outcome{Kind: Fail, Evidence: "git commit: " + firstLine(err)}
	}
	got, err := sys.revParse(target, "HEAD^{tree}")
	if err != nil || got != tree {
		return Outcome{Kind: Fail, Evidence: fmt.Sprintf("the root tree %s is not pin.tree %s", got, tree)}
	}
	return Outcome{Kind: Done, Evidence: "root tree = pin.tree"}
}

// layupPin is the commit of LAYUP's own pin (docs/setup/armature.pin), which
// the comment of the push names (D9 of #86, section 5 Start 2 of the
// architecture).
const layupPin = "a95965534b14b0bf14ad74da0c9a45b5f4aedf88"

// commandsS03 gives the commands of S03 (D9 of #86), run from the work area:
// the remote origin of the target on GitHub, and the push of the root commit.
func commandsS03(r work.Record) []Command {
	name, _ := r.Value("S01", "name")
	commit, _ := r.Value("S02", "pin.commit")
	pin := "not the commit of LAYUP's own pin " + layupPin[:7]
	if commit == layupPin {
		pin = "the commit of LAYUP's own pin " + layupPin[:7]
	}
	return []Command{
		{Order: 1, Comment: "the remote of the target: its empty repository on GitHub", Text: "git -C " + work.TargetPath + " remote add origin " + quote("https://github.com/"+name+".git")},
		{Order: 1, Comment: "the push of the root commit: the unmodified baseline at " + commit + ", " + pin, Text: "git -C " + work.TargetPath + " push origin main"},
	}
}

// quote gives s as one word of sh.
func quote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

// The paths that S04 writes.
const (
	pinPath       = "docs/setup/armature.pin"
	adrDir        = "docs/adr"
	factsDir      = "docs/facts"
	factsSumsPath = "docs/setup/facts.sha256"
)

// runS04 is S04 (D10 of #86, with condition 1 of its plan review): on the
// branch layup-setup from the root commit, the pin file, the decision record
// of the pin with its index row, and the answers record of the S01- and Q-
// answers with its index row and its line of facts.sha256 (O-124). Each file
// comes from the files of the root commit, so a rerun writes the same files.
func runS04(b Brief, in Input) Outcome {
	target := filepath.Join(in.Dir, work.TargetPath)
	pinned := map[string]string{}
	for _, n := range []string{"pin.source", "pin.commit", "pin.tree", "pin.time"} {
		v, ok := in.Record.Value("S02", n)
		if !ok {
			return Outcome{Kind: Fail, Evidence: "the record has no row S02 " + n}
		}
		pinned[n] = v
	}
	date := pinned["pin.time"][:min(10, len(pinned["pin.time"]))]
	root, err := sys.revParse(target, "refs/heads/main^{commit}")
	if err != nil {
		return Outcome{Kind: Fail, Evidence: "the root commit: " + firstLine(err)}
	}
	if o, ok := onSetupBranch(target, root); !ok {
		return o
	}
	readme := func(dir string) (string, error) {
		data, err := sys.show(target, root, dir+"/README.md")
		if err != nil {
			return "", errors.New(dir + "/README.md of the baseline: " + firstLine(err))
		}
		return string(data), nil
	}
	n, err := nextNumber(target, root, adrDir, `^([0-9]{4})-[^/]*\.md$`)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	f, err := nextNumber(target, root, factsDir, `^F-([0-9]{4})-[^/]*\.md$`)
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	id := fmt.Sprintf("F-%04d", f)
	adrIndex, err := readme(adrDir)
	if err == nil {
		adrIndex, err = addIndexRow(adrIndex, fmt.Sprintf("| [%04d](%s) | Pin the baseline | Accepted |", n, adrFile(n)))
	}
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	factsIndex, err := readme(factsDir)
	if err == nil {
		factsIndex, err = addIndexRow(factsIndex, fmt.Sprintf("| [%s](%s-setup-answers.md) | The answers to the questions of the setup | %s | Raw |", id, id, date))
	}
	if err != nil {
		return Outcome{Kind: Fail, Evidence: err.Error()}
	}
	gaps, err := readGaps(b.Gaps)
	if err != nil {
		return Outcome{Kind: Invalid, Evidence: err.Error()}
	}
	asked := slices.Clone(work.S01Questions)
	for _, g := range gaps {
		asked = append(asked, work.Question{ID: g.ID, Text: g.Question})
	}
	record := answersRecord(id, date, setupRecord, asked, in.Answers)
	recordPath := factsDir + "/" + id + "-setup-answers.md"
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(record)))
	sums, err := sys.show(target, root, factsSumsPath)
	if err != nil {
		sums = nil // the baseline at LAYUP's pin has no docs/setup/
	}
	if len(sums) > 0 && sums[len(sums)-1] != '\n' {
		sums = append(sums, '\n') // a last line with no line feed stays a line
	}
	files := map[string]string{
		pinPath:                   pinText(pinned["pin.source"], pinned["pin.commit"], pinned["pin.tree"], date),
		adrDir + "/" + adrFile(n): adrText(n, date, pinned["pin.source"], pinned["pin.commit"], pinned["pin.tree"]),
		adrDir + "/README.md":     adrIndex,
		recordPath:                record,
		factsDir + "/README.md":   factsIndex,
		factsSumsPath:             string(sums) + sum + "  " + recordPath + "\n",
	}
	for p, text := range files {
		if err := sys.write(filepath.Join(target, filepath.FromSlash(p)), []byte(text)); err != nil {
			return Outcome{Kind: Fail, Evidence: "S04: " + firstLine(err)}
		}
	}
	return Outcome{Kind: Done, Evidence: "checks pin and facts", Commit: true, Values: [][]string{
		{"pin.adr", adrDir + "/" + adrFile(n), "computed", "the next free number of " + adrDir + "/"},
		{"answers.record", recordPath, "computed", "the next free ID of " + factsDir + "/"},
		{"answers.record.sha256", sum, "computed", "sha256 " + recordPath},
	}}
}

// onSetupBranch puts the target on the branch layup-setup from the root commit
// (condition 1 of the plan review of #86): it makes the branch when there is
// none, and takes it when the target is on it at the root commit; another
// state is an input error.
func onSetupBranch(target, root string) (Outcome, bool) {
	branch, err := sys.branch(target)
	if err != nil {
		return Outcome{Kind: Invalid, Evidence: work.TargetPath + " of the work area is on no branch"}, false
	}
	_, noBranch := sys.revParse(target, "refs/heads/layup-setup^{commit}")
	switch {
	case branch == "refs/heads/layup-setup":
		if head, err := sys.head(target); err != nil || head != root {
			return Outcome{Kind: Invalid, Evidence: "the branch layup-setup of " + work.TargetPath + " is not at the root commit, and S04 is not done: start again in a new work area"}, false
		}
	case branch == "refs/heads/main" && noBranch != nil:
		if err := sys.switchCreate(target, "layup-setup", root); err != nil {
			return Outcome{Kind: Fail, Evidence: "git switch -c layup-setup: " + firstLine(err)}, false
		}
	case branch == "refs/heads/main":
		return Outcome{Kind: Invalid, Evidence: work.TargetPath + " of the work area is on main, and its branch layup-setup exists: put it on layup-setup, or start again in a new work area"}, false
	default:
		return Outcome{Kind: Invalid, Evidence: work.TargetPath + " of the work area is on " + branch + ", not on main or layup-setup"}, false
	}
	return Outcome{}, true
}

// nextNumber gives the next number after the highest number of a file of dir
// at commit whose name matches form (its first group is the number); 1 when
// none does.
func nextNumber(target, commit, dir, form string) (int, error) {
	entries, err := sys.lsTree(target, commit, dir)
	if err != nil {
		return 0, errors.New(dir + "/ of the baseline: " + firstLine(err))
	}
	re := regexp.MustCompile(form)
	high := 0
	for _, e := range entries {
		if m := re.FindStringSubmatch(strings.TrimPrefix(e.Path, dir+"/")); m != nil {
			v, _ := strconv.Atoi(m[1])
			high = max(high, v)
		}
	}
	return high + 1, nil
}

// pinText gives the pin file of the target, in the form of NFR-006.
func pinText(source, commit, tree, date string) string {
	return fmt.Sprintf("source=%s\ncommit=%s\ntree=%s\nmethod=git clone %s, checkout %s, .git removed\ndate=%s\n",
		source, commit, tree, source, commit, date)
}

// adrFile gives the file name of the decision record of the pin.
func adrFile(n int) string { return fmt.Sprintf("%04d-pin-the-baseline.md", n) }

// adrText gives the decision record of the pin, from a fixed text in the
// baseline's form (D10 of #86; docs/spec/setup.md holds the text): no
// reference to an issue, and the pin file named in a code span.
func adrText(n int, date, source, commit, tree string) string {
	return fmt.Sprintf(`# %04d. Pin the baseline

Date: %s

## Status

Accepted

## Context

This repository started as a copy of a baseline, the repository at
`+"`%s`"+`. A copy that names no version of its source cannot be compared
with that source later, and a later change of the source would change what
"the baseline" means.

## Decision

We will pin the baseline at the commit `+"`%s`"+`, whose tree is
`+"`%s`"+`: the latest commit of its default branch when this repository was
set up. The pin file `+"`docs/setup/armature.pin`"+` holds the source, the commit, the tree,
the method and the date. The root commit of the default branch is the unchanged
copy, and its tree is the pinned tree.

We reject a copy with no recorded version, which no one can compare with its
source.

## Consequences

- A later version of the baseline comes into this repository only by a change
  that names the new commit.
- The pin file and the root commit show the version that this repository
  started from.
`, n, date, source, commit, tree)
}

var (
	indexHeading = regexp.MustCompile(`(?m)^## Index[ \t]*$`)
	placeholder  = regexp.MustCompile(`^\|[ \t]*_none yet_[ \t]*(\|[ \t]*)+$`)
)

// addIndexRow puts row into the table under the heading "## Index" of an
// index file: in place of its one row when that row is the placeholder
// _none yet_, else after its last row.
func addIndexRow(readme, row string) (string, error) {
	loc := indexHeading.FindStringIndex(readme)
	if loc == nil {
		return "", errors.New("the index file has no heading ## Index")
	}
	lines := strings.SplitAfter(readme[loc[1]:], "\n")
	first, last := -1, -1
	for i, l := range lines {
		t := strings.TrimRight(l, "\n")
		if strings.HasPrefix(t, "|") {
			if first < 0 {
				first = i
			}
			last = i
		} else if first >= 0 {
			break
		}
	}
	if first < 0 || last < first+2 {
		return "", errors.New("the index file has no table under ## Index")
	}
	if last == first+2 && placeholder.MatchString(strings.TrimRight(lines[last], "\n")) {
		lines[last] = row + "\n"
	} else {
		if !strings.HasSuffix(lines[last], "\n") { // a table at the end of a file with no line feed
			lines[last] += "\n"
		}
		lines = slices.Insert(lines, last+1, row+"\n")
	}
	return readme[:loc[1]] + strings.Join(lines, ""), nil
}

// A recordKind is one of the two answers records: its title, the text of its
// Source, and the step that writes it (K15; condition 3 of the plan review of
// #90).
type recordKind struct{ title, source, step string }

var (
	setupRecord  = recordKind{"The answers to the questions of the setup", "The answers of `" + work.AnswersPath + "` to the questions of S01 and to the gaps of the problem statement", "S04"}
	markerRecord = recordKind{"The answers to the markers of the setup", "The answers of `" + work.AnswersPath + "` to the markers that S10 listed", "S11"}
)

// answersRecord gives an answers record (K15, K17; condition 2 of the plan
// review of #86, condition 3 of #90): the header table, and one fact per
// answer of a question of asked, in the order of asked, with each angle quote
// of a recorded text (the answer, its source, the question) as an entity.
func answersRecord(id, date string, k recordKind, asked []work.Question, a work.Answers) string {
	entity := strings.NewReplacer("\u2039", "&lsaquo;", "\u203a", "&rsaquo;")
	var b strings.Builder
	fmt.Fprintf(&b, "# %s. %s\n\n| Field | Value |\n| ------------ | ----- |\n"+
		"| Fact ID | `%s` |\n"+
		"| Source | %s |\n"+
		"| Collected by | `layup setup`, step %s |\n| Date collected | %s |\n"+
		"| Origin | `%s` of the work area |\n| Status | `Raw` |\n\n## Facts as collected\n\n", id, k.title, id, k.source, k.step, date, work.AnswersPath)
	n := 0
	for _, q := range asked {
		for _, r := range a {
			if r[0] == q.ID {
				n++
				fmt.Fprintf(&b, "%d. `%s` %s — by %s; source %s; the question: %s\n", n, q.ID, entity.Replace(r[1]), r[2], entity.Replace(r[3]), entity.Replace(q.Text))
			}
		}
	}
	b.WriteString("\n## Notes on capture\n\nEach angle quote of a recorded text is written as `&lsaquo;` or `&rsaquo;`.\n")
	return b.String()
}

// firstLine gives the first line of the text of err.
func firstLine(err error) string {
	s, _, _ := strings.Cut(err.Error(), "\n")
	return s
}
