package run

// This file holds the prompt file of a step of a task (docs/spec/loop.md, The
// tier and the context of a task; docs/architecture.md §9, The context of a
// session; task T-6zj3, #361, row 66e of the plan of M2e): the list of links,
// in the five kinds of the step table in their order, each a path with the
// SHA-256 of the bytes that code read, built from the records and the target at
// the base. The records items are copied under the session's input/, as the
// session's clone holds no records branch. The text form of the plan, and so the
// tier it gives, are the caller's (tier.go); the issue's own text is no record
// in M2e, so the issue is the row of tasks.tsv (known limit, loop.md).

import (
	"context"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// An Item is one link of a prompt file: its kind, whether it is a file of the
// target at the base or of the records, its path there, and the SHA-256 of its
// bytes.
type Item struct {
	Kind        session.Kind
	Where, Path string
	SHA         string
}

// PromptIn is what the prompt file of a step is built for: the task, the step
// of the step table, the base commit of the target and the findings of the last
// attempt (none for a first attempt).
type PromptIn struct {
	Task, Step, Base string
	Findings         []Finding
	// Options are the results of the filter of each decision of an architecture
	// task (loop.md, The instruction of a task step).
	Options map[string]spec.Result
}

// A PromptContext is the prompt file of a step and the input files it names.
// Input holds the records items by their file name, for session.MakeInput.
type PromptContext struct {
	Items    []Item
	Input    map[string][]byte
	task     string
	step     string
	findings []Finding
	kind     string
	issue    int
	base     string
	facts    string
	record   string
	extras   taskExtras
}

// A kept is a file of the records, by its path.
type kept struct {
	Path string
	Data []byte
}

// targetTree is the target at the base: its paths, in order, and a read of one.
type targetTree struct {
	Paths []string
	Read  func(path string) ([]byte, error)
}

// promptSrc is what the prompt file of a step is built from.
type promptSrc struct {
	Step         session.StepRow
	Task         string
	Register     kept
	Plan         *kept
	Handoffs     []kept
	Answers      []kept
	Requirements []string
	Rules        []rules.Entry
	Tree         targetTree
	Findings     []Finding
	Kind         string
	Issue        int
	Base         string
	Extras       taskExtras
}

const (
	prdDir  = "docs/prd/"
	specDir = "docs/spec/"
	adrDir  = "docs/adr/"
)

var adrName = regexp.MustCompile(`ADR-([0-9]{4})`)

func alnum(b byte) bool { return b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' }

// holdsID reports whether line holds the requirement ID as a whole word: no
// letter or digit before or after it (REQ-0050 does not hold REQ-005).
func holdsID(line, id string) bool {
	for from := 0; ; {
		i := strings.Index(line[from:], id)
		if i < 0 {
			return false
		}
		i += from
		end := i + len(id)
		if (i == 0 || !alnum(line[i-1])) && (end == len(line) || !alnum(line[end])) {
			return true
		}
		from = i + 1
	}
}

// headingLevel is the level of a Markdown heading line, 0 for any other line.
func headingLevel(line string) int {
	n := 0
	for n < len(line) && line[n] == '#' {
		n++
	}
	if n == 0 || n > 6 || len(line) <= n || line[n] != ' ' {
		return 0
	}
	return n
}

// sectionsOf gives the text of each section of a Markdown file whose heading
// holds id: from the heading down to the next heading of its level or above.
func sectionsOf(text, id string) []string {
	lines := strings.Split(text, "\n")
	var out []string
	for i, l := range lines {
		level := headingLevel(l)
		if level == 0 || !holdsID(l, id) {
			continue
		}
		j := i + 1
		for j < len(lines) && (headingLevel(lines[j]) == 0 || headingLevel(lines[j]) > level) {
			j++
		}
		out = append(out, strings.Join(lines[i:j], "\n"))
	}
	return out
}

// hasRow reports whether a table row of the Markdown file holds the id.
func hasRow(text, id string) bool {
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), "|") && holdsID(l, id) {
			return true
		}
	}
	return false
}

// buildPrompt builds the prompt file from the sources: the items of each kind
// of the step in the order of the step table, each file once (the first kind that
// names it keeps it).
func buildPrompt(src promptSrc) (*PromptContext, error) {
	if len(src.Step.Kinds) == 0 {
		return nil, fmt.Errorf("the step %s has no record kinds: it has no prompt file", src.Step.Step)
	}
	if _, ok := stepInstructions[src.Step.Step]; !ok {
		return nil, noInstruction(src.Step.Step)
	}
	c := &PromptContext{Input: map[string][]byte{}, task: src.Task, step: src.Step.Step, findings: src.Findings,
		kind: src.Kind, issue: src.Issue, base: src.Base, extras: src.Extras}
	if src.Step.Step == "specification" {
		var err error
		if c.record, c.facts, err = factsRecords(src.Tree.Paths); err != nil {
			return nil, err
		}
	}
	seen := map[string]bool{}
	names := map[string]string{} // file name under input/ -> the records path
	record := func(kind session.Kind, k kept) error {
		if seen["records:"+k.Path] {
			return nil
		}
		name := path.Base(k.Path)
		if other, ok := names[name]; ok && other != k.Path {
			return fmt.Errorf("two records are named %s: %s and %s", name, other, k.Path)
		}
		names[name], seen["records:"+k.Path] = k.Path, true
		c.Input[name] = k.Data
		c.Items = append(c.Items, Item{kind, "records", k.Path, digest(k.Data)})
		return nil
	}
	target := func(kind session.Kind, p string) error {
		if seen["target:"+p] {
			return nil
		}
		data, err := src.Tree.Read(p)
		if err != nil {
			return err
		}
		seen["target:"+p] = true
		c.Items = append(c.Items, Item{kind, "target", p, digest(data)})
		return nil
	}
	underDir := func(dir string) []string {
		var out []string
		for _, p := range src.Tree.Paths {
			if strings.HasPrefix(p, dir) && strings.HasSuffix(p, ".md") {
				out = append(out, p)
			}
		}
		return out
	}
	for _, kind := range src.Step.Kinds {
		var err error
		switch kind {
		case session.KindIssueAndPlan:
			if err = record(kind, src.Register); err == nil && src.Plan != nil {
				err = record(kind, *src.Plan)
			}
			if err == nil && src.Extras.Frozen != nil {
				err = record(kind, *src.Extras.Frozen)
			}
		case session.KindRequirements:
			for _, id := range src.Requirements {
				for _, p := range underDir(prdDir) {
					data, rerr := src.Tree.Read(p)
					if rerr != nil {
						return nil, rerr
					}
					if hasRow(string(data), id) {
						if err := target(kind, p); err != nil {
							return nil, err
						}
					}
				}
				for _, p := range underDir(specDir) {
					data, rerr := src.Tree.Read(p)
					if rerr != nil {
						return nil, rerr
					}
					if len(sectionsOf(string(data), id)) > 0 {
						if err := target(kind, p); err != nil {
							return nil, err
						}
					}
				}
			}
			for _, k := range src.Extras.Shape {
				if err = record(kind, k); err != nil {
					break
				}
			}
		case session.KindDecisionRecords:
			if err = decisionRecords(src, underDir, target, kind); err == nil {
				for _, k := range src.Extras.Options {
					if err = record(kind, k); err != nil {
						break
					}
				}
			}
		case session.KindHandoffsAndAnswers:
			for _, k := range slices.Concat(src.Handoffs, src.Answers, src.Extras.Ledger) {
				if err = record(kind, k); err != nil {
					break
				}
			}
		case session.KindRuleFiles:
			for _, p := range src.Tree.Paths {
				if _, ok := rules.Match(src.Rules, p); ok {
					if err := target(kind, p); err != nil {
						return nil, err
					}
				}
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// attemptBase gives the base of the last event attempt of a task: the commit that the
// change of the verifier starts from.
func attemptBase(events [][]string) (string, bool) {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i][1] == "attempt" {
			return events[i][4], true
		}
	}
	return "", false
}

// noInstruction is the refusal of a step that has no instruction in M2e.
func noInstruction(step string) error {
	return fmt.Errorf("the step %s has no instruction in M2e", step)
}

// factsRecords gives the numbers NNNN of the highest record of docs/facts/ in the paths and of
// the next one: the record that the Facts cells of the specification cite because it exists, and
// the one that the rendered commit of the first bet will write.
func factsRecords(paths []string) (record, facts string, err error) {
	var names []string
	for _, p := range paths {
		if path.Dir(p) == "docs/facts" {
			names = append(names, path.Base(p))
		}
	}
	ids := factRecordIDs(names)
	if len(ids) == 0 {
		return "", "", fmt.Errorf("the target has no record in docs/facts/ at the base: the Facts cells of the specification would cite none")
	}
	return ids[len(ids)-1][len("F-"):], nextFactsID(ids)[len("F-"):], nil
}

// decisionRecords adds the decision record of each ADR-NNNN that a section of the
// specification of the task's requirements names; one that the target lacks is an
// error naming it.
func decisionRecords(src promptSrc, underDir func(string) []string, target func(session.Kind, string) error, kind session.Kind) error {
	var nums []string
	for _, id := range src.Requirements {
		for _, p := range underDir(specDir) {
			data, err := src.Tree.Read(p)
			if err != nil {
				return err
			}
			for _, sec := range sectionsOf(string(data), id) {
				for _, m := range adrName.FindAllStringSubmatch(sec, -1) {
					if !slices.Contains(nums, m[1]) {
						nums = append(nums, m[1])
					}
				}
			}
		}
	}
	slices.Sort(nums)
	for _, n := range nums {
		i := slices.IndexFunc(underDir(adrDir), func(p string) bool { return strings.HasPrefix(p, adrDir+n+"-") })
		if i < 0 {
			return fmt.Errorf("ADR-%s, which a section of the specification names, is not in %s of the target", n, adrDir)
		}
		if err := target(kind, underDir(adrDir)[i]); err != nil {
			return err
		}
	}
	return nil
}

// Text gives the prompt file of the session whose directory is root: for each
// kind its items, each with its path and SHA-256, a records item at its place
// under root/input/ and a target item at its path in the clone; then, under the
// handoffs, the findings of the last attempt, and no diff.
func (c *PromptContext) Text(root string) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "Prompt file of the step %s of the task %s.\n", c.step, c.task)
	b.WriteString("Read each file at its path: its SHA-256 follows it. A path that starts with / is a copy of a file of the run's records; any other path is a file of the target, relative to your working directory.\n")
	for _, kind := range session.ContextKinds() {
		fmt.Fprintf(&b, "\n%s:\n", kind)
		n := 0
		for _, it := range c.Items {
			if it.Kind != kind {
				continue
			}
			p := it.Path
			if it.Where == "records" {
				p = path.Join(root, "input", path.Base(it.Path))
			}
			fmt.Fprintf(&b, "- %s %s\n", p, it.SHA)
			n++
		}
		if n == 0 {
			b.WriteString("- none\n")
		}
		if kind == session.KindDecisionRecords && len(c.extras.Survivors) > 0 {
			b.WriteString("the survivors of the filter:\n")
			for _, sv := range c.extras.Survivors {
				fmt.Fprintf(&b, "- %s: %s\n", sv.Decision, strings.Join(sv.Options, " "))
			}
		}
		if kind == session.KindHandoffsAndAnswers && len(c.findings) > 0 {
			b.WriteString("findings of the last attempt:\n")
			for _, f := range c.findings {
				where := ""
				if f.File != "" {
					where = " " + f.File
					if f.Line > 0 {
						where += fmt.Sprintf(":%d", f.Line)
					}
				}
				fmt.Fprintf(&b, "- %s%s: %s\n", f.Rule, where, f.Text)
			}
		}
	}
	v := map[string]string{"task": c.task, "kind": c.kind, "issue": strconv.Itoa(c.issue), "result": path.Join(root, "result", "result.tsv"),
		"comment": path.Join(root, "result", commentFiles[c.step]), "head": c.base, "from": c.extras.From,
		"labels": labelList(DefaultOwners()), "facts": c.facts, "record": c.record}
	b.WriteString("\nThe work of this step:\n" + fill(stepInstructions[c.step], v) + "\n" + fill(resultInstruction, v))
	return []byte(b.String())
}

// readRecord reads a file of the records that must be there.
func (s *Sessions) readRecord(ctx context.Context, p string) (kept, error) {
	data, err := s.Store.ReadFile(ctx, p)
	if err != nil {
		return kept{}, err
	}
	if data == nil {
		return kept{}, fmt.Errorf("the records have no %s", p)
	}
	return kept{p, data}, nil
}

// PromptFile builds the prompt file of a step of a task from the records and
// the target at in.Base. The plan is the kept copy of the last valid handoff
// with a comment from a session of the plan role; the handoffs are the copies of
// every valid comment handoff, and the answers the result files that hold an
// answer row; the requirement IDs are those of the task's row of tasks.tsv, and
// the rule files those that the rule-path register matches.
func (s *Sessions) PromptFile(ctx context.Context, in PromptIn) (*PromptContext, error) {
	var row session.StepRow
	for _, r := range session.StepTable() {
		if r.Step == in.Step {
			row = r
		}
	}
	if row.Step == "" {
		return nil, fmt.Errorf("the step table has no step %q", in.Step)
	}
	if _, ok := stepInstructions[in.Step]; !ok {
		return nil, noInstruction(in.Step) // before every read
	}
	register, err := s.readRecord(ctx, tasksPath)
	if err != nil {
		return nil, err
	}
	tasks, err := ReadRegister(register.Data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tasksPath, err)
	}
	src := promptSrc{Step: row, Task: in.Task, Register: register, Findings: in.Findings, Base: in.Base}
	i := slices.IndexFunc(tasks, func(t Task) bool { return t.ID == in.Task })
	if i < 0 {
		return nil, fmt.Errorf("%s has no row of %s", tasksPath, in.Task)
	}
	src.Requirements, src.Kind, src.Issue = tasks[i].Requirements, tasks[i].Kind, tasks[i].Issue
	if err := checkStepKind(in.Step, src.Kind); err != nil {
		return nil, fmt.Errorf("%w: the task %s is of kind %s", err, in.Task, src.Kind)
	}
	if _, takes, _ := inputsWanted(in.Step, src.Kind); takes != (len(in.Options) > 0) {
		return nil, fmt.Errorf("the options of the task %s: the step %s of a task of kind %s takes %s, and %d decisions are given", in.Task, in.Step, src.Kind,
			map[bool]string{true: "the options of its decisions", false: "none"}[takes], len(in.Options))
	}
	events, err := s.readTable(ctx, eventsPath(in.Task), records.ReadEvents)
	if err != nil {
		return nil, err
	}
	if in.Step == "verification" {
		from, ok := attemptBase(events)
		if !ok {
			return nil, fmt.Errorf("the task %s has no attempt", in.Task)
		}
		src.Extras.From = from
	}
	sessions, err := s.readTable(ctx, "sessions.tsv", records.ReadSessions)
	if err != nil {
		return nil, err
	}
	role := map[string]string{}
	for _, r := range sessions {
		role[r[sessionID]] = r[sessionRole]
	}
	for _, e := range events {
		switch _, comment := validHandoff(e[6]); {
		case e[1] == "handoff" && comment:
			k, err := s.readRecord(ctx, "copies/"+strings.Fields(e[6])[2]+".md")
			if err != nil {
				return nil, err
			}
			src.Handoffs = append(src.Handoffs, k)
			if role[e[3]] == stepRole("plan-of-a-task") {
				src.Plan = &k
			}
		case e[1] == "result" && e[6] == "done":
			p := "tasks/" + in.Task + "/results/" + e[3] + ".tsv"
			data, err := s.Store.ReadFile(ctx, p)
			if err != nil || data == nil {
				if err != nil {
					return nil, err
				}
				continue
			}
			rows, err := records.ReadResult(data)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", p, err)
			}
			if slices.ContainsFunc(rows, func(r []string) bool { return r[0] == "answer" }) {
				src.Answers = append(src.Answers, kept{p, data})
			}
		}
	}
	extras, err := s.taskInputs(ctx, in, src.Kind)
	if err != nil {
		return nil, err
	}
	extras.From = src.Extras.From
	src.Extras = extras
	regData, err := s.readRecord(ctx, "rule-paths.tsv")
	if err != nil {
		return nil, err
	}
	if src.Rules, err = rules.ReadRegister(regData.Data); err != nil {
		return nil, fmt.Errorf("rule-paths.tsv: %w", err)
	}
	entries, err := git.LsTree(s.Clone, in.Base, ".")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.Type == "blob" {
			src.Tree.Paths = append(src.Tree.Paths, e.Path)
		}
	}
	src.Tree.Read = func(p string) ([]byte, error) { return git.Show(s.Clone, in.Base, p) }
	return buildPrompt(src)
}

// taskInputs reads the inputs that come with the task, as inputsWanted says for
// the step (loop.md, The records each step receives): the frozen list when it
// exists, the revised draft and its facts, the options with their survivors and
// the ledger, each by its reader, which fails closed.
func (s *Sessions) taskInputs(ctx context.Context, in PromptIn, kind string) (taskExtras, error) {
	var x taskExtras
	if data, err := s.Store.ReadFile(ctx, testsPath(in.Task)); err != nil {
		return x, err
	} else if data != nil {
		if _, err := ReadTests(data); err != nil {
			return x, fmt.Errorf("%s: %w", testsPath(in.Task), err)
		}
		x.Frozen = &kept{testsPath(in.Task), data}
	}
	shape, options, ledger := inputsWanted(in.Step, kind)
	readers := map[string]func([]byte) ([][]string, error){"spec/draft.tsv": records.ReadDraft, "spec/spans.tsv": records.ReadSpans,
		"start/problem-statement.md": nil, "intake/answers.tsv": records.ReadAnswers, "sessions.tsv": records.ReadSessions, "telemetry.tsv": records.ReadTelemetry}
	var paths []string
	if shape {
		paths = []string{"spec/draft.tsv", "spec/spans.tsv", "start/problem-statement.md", "intake/answers.tsv"}
	}
	if ledger {
		paths = append(paths, "sessions.tsv", "telemetry.tsv")
	}
	for _, p := range paths {
		k, err := s.readRecord(ctx, p)
		if err != nil {
			return x, err
		}
		if read := readers[p]; read != nil {
			if _, err := read(k.Data); err != nil {
				return x, fmt.Errorf("%s: %w", p, err)
			}
		}
		if ledger {
			x.Ledger = append(x.Ledger, k)
		} else {
			x.Shape = append(x.Shape, k)
		}
	}
	if options {
		var err error
		x.Options, x.Survivors, err = optionsKept(in.Options, func(p string) ([]byte, error) {
			k, err := s.readRecord(ctx, p)
			return k.Data, err
		})
		if err != nil {
			return x, err
		}
	}
	return x, nil
}

// factRecordName is the name of a record of docs/facts/ as the pinned prd-lint reads it:
// F-NNNN, anything, .md.
var factRecordName = regexp.MustCompile(`^(F-[0-9]{4}).*\.md$`)

// factRecordIDs gives the record IDs (F-NNNN) of the base names of the direct entries of
// docs/facts/, sorted, each once.
func factRecordIDs(names []string) []string {
	var ids []string
	for _, n := range names {
		if m := factRecordName.FindStringSubmatch(n); m != nil && !slices.Contains(ids, m[1]) {
			ids = append(ids, m[1])
		}
	}
	slices.Sort(ids)
	return ids
}

// nextFactsID is the ID that follows the highest of ids: the record that the rendered commit
// of the first bet writes (loop.md, The instruction of a task step). It is the one home of
// that rule: the prompt of the specification and the input of the full check call it.
func nextFactsID(ids []string) string {
	top := 0
	for _, id := range ids {
		n, _ := strconv.Atoi(id[len("F-"):])
		top = max(top, n)
	}
	return fmt.Sprintf("F-%04d", top+1)
}

// completeFactsRecords sets the records of docs/facts/ in the input of the full check and,
// before the first bet (no index) with no numbered facts file, the predicted problem
// statement record, so that the citations of a specification that waits for its bet resolve.
func completeFactsRecords(in *spec.FullInput, names []string) {
	in.Records = factRecordIDs(names)
	if len(in.Index) == 0 && len(in.PSBRecords) == 0 {
		in.PSBRecords = []string{nextFactsID(in.Records)}
	}
}
