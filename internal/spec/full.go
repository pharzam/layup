package spec

// The full spec check (docs/spec/loop.md, REQ-012 — The full spec check; the
// block task-spec-check of docs/spec/records.md; task T-4wyg, #351, row 62c):
// ten rules, in the block's order, each a row of the table. The function is
// pure: every input is a record or a file the caller has already read (and,
// for the inventory, already given the exact bytes the index's sha256 names).
// The one process seam is the pinned prd-lint of the target (prdlint.go).

import (
	"crypto/sha256"
	"fmt"
	"io"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

// TaskSpecCheckSchema is the form of the table the full check writes: the
// block task-spec-check, which this package owns.
var TaskSpecCheckSchema = tsv.Schema{Name: "task-spec-check", Location: "records:tasks/<task>/spec-check.tsv", Columns: []tsv.Column{
	{Name: "head", Type: "sha1", Key: true},
	{Name: "rule", Type: "enum(covers|needs|measures|criteria|classes|inventory|quotes|sections|tasks|ran)", Key: true},
	{Name: "result", Type: "enum(pass|fail|clear)"},
	{Name: "detail", Type: "text"},
}}

// tasksSchema is a second Go value of the block tasks, whose owner is
// internal/run; this package reads the register with it, and its own test
// compares it with the block (the spans precedent of records.md).
var tasksSchema = tsv.Schema{Name: "tasks", Location: "records:tasks.tsv", Columns: []tsv.Column{
	{Name: "task", Type: "id(T-xxxx)", Key: true},
	{Name: "issue", Type: "int"},
	{Name: "kind", Type: "enum(specification|architecture|batch|rendered|plan|build)"},
	{Name: "requirements", Type: "list(text)"},
	{Name: "size", Type: "enum(small|large)"},
	{Name: "predecessors", Type: "list(id(T-xxxx))"},
	{Name: "pr", Type: "int"},
	{Name: "opened", Type: "time"},
}}

// The nine judged rules in the block's order; the tenth row, ran, judges the
// nine.
var fullRules = []string{"covers", "needs", "measures", "criteria", "classes", "inventory", "quotes", "sections", "tasks"}

const (
	resPass  = "pass"
	resFail  = "fail"
	resClear = "clear"

	notRun = "did not run: "
)

// FullInput is everything the full check reads, supplied by the caller —
// internal/run in the loop, after the first bet — as records already read and
// files already taken. A record the caller cannot read is a reason not to
// call, not an input to fake: an empty Index is the one honest signal of
// "before the first bet", never of a read that failed.
type FullInput struct {
	Changed       []string          // the paths the head changes against the base, repo-relative
	Records       []string          // the record IDs of the direct docs/facts/ entries at the head that prd-lint reads (F-NNNN*.md); a bare token naming one is no failure of covers
	PSBRecords    []string          // the record IDs that number the problem statement's facts (LAYUP's own are F-0001 and F-0003, two records of the one file)
	AnswersRecord string            // the record ID of the answers, whose fact N answers Q-00N (LAYUP's is F-0004)
	PRDs          map[string][]byte // the head's docs/prd/PRD-*.md, by base name
	SpecFiles     map[string][]byte // the head's docs/spec/*.md, by base name
	PSB           []byte            // the problem statement the span offsets index
	Spans         [][]string        // the rows of spec/spans.tsv, as Read gives them
	Questions     [][]string        // the rows of questions.tsv
	Answers       [][]string        // the rows of intake/answers.tsv
	Index         [][]string        // the rows of spec/inventory.tsv, as ReadInventoryIndex gives them; empty is before the first bet
	Versions      map[int][]byte    // the exact bytes of spec/inventory-<n>.tsv, for each n of Index
	Copy          []byte            // the text of the latest bet's or decision's copy (the answer's lines)
	Register      []byte            // the bytes of tasks.tsv
	Merged        []string          // the IDs of the register's tasks whose pull request merged
}

// CheckFull gives the ten rows of the table tasks/<task>/spec-check.tsv:
// rule, result, detail, in the block's order (the head column is the writer's,
// not the check's). Before the first bet — an empty Index — a head that
// touches docs/prd/ or docs/spec/ is judged against the numbered facts of
// Intake and the answers, with classes and inventory clear "before the first
// bet"; a head that touches neither is clear for the whole table with the
// reason "no specification yet". After the first bet every rule is pass or
// fail, whatever it touches (loop.md: `no specification yet` is only before
// the first bet). A rule whose input is absent or unreadable did not run: its
// row is a fail, and ran fails it again (NFR-004).
func CheckFull(in FullInput) [][]string {
	if len(in.Index) == 0 && !touchesSpec(in.Changed) {
		return clearTable("no specification yet")
	}
	c := fullCheck{in: in, prd: parsePRDs(in.PRDs), u: universeOf(in)}
	var rows [][]string
	for _, rule := range fullRules {
		if len(in.Index) == 0 && (rule == "classes" || rule == "inventory") {
			rows = append(rows, []string{rule, resClear, "before the first bet"})
			continue
		}
		v := c.judge(rule)
		rows = append(rows, []string{rule, v.result, v.detail})
	}
	return append(rows, ranRow(rows))
}

func touchesSpec(changed []string) bool {
	return slices.ContainsFunc(changed, func(p string) bool {
		return strings.HasPrefix(p, "docs/prd/") || strings.HasPrefix(p, "docs/spec/")
	})
}

func clearTable(reason string) [][]string {
	rows := make([][]string, 10)
	for i, rule := range append(slices.Clone(fullRules), "ran") {
		rows[i] = []string{rule, resClear, reason}
	}
	return rows
}

// WriteSpecCheck writes the table of the full check, rows as CheckFull gives
// them, with the head in the first column, by TaskSpecCheckSchema. It writes
// nothing on a row the schema refuses.
func WriteSpecCheck(w io.Writer, head string, rows [][]string) error {
	full := make([][]string, len(rows))
	for i, r := range rows {
		full[i] = append([]string{head}, r...)
	}
	return tsv.Write(w, TaskSpecCheckSchema, full)
}

// ranRow gives the row of the rule ran over the nine rows before it: fail
// when a rule's row is missing or is a did-not-run (NFR-004 — a check that
// did not run never counts as passed), pass otherwise.
func ranRow(rows [][]string) []string {
	seen := map[string]bool{}
	for _, r := range rows {
		seen[r[0]] = true
	}
	for _, rule := range fullRules {
		if !seen[rule] {
			return []string{"ran", resFail, fmt.Sprintf("the rule %s has no row", rule)}
		}
	}
	for _, r := range rows {
		if r[1] == resFail && strings.HasPrefix(r[2], notRun) {
			return []string{"ran", resFail, notRun + r[0]}
		}
	}
	return []string{"ran", resPass, ""}
}

// verdict is one rule's result and detail ("" is the pass detail, which the
// writer writes as —; a clear row's detail is its reason).
type verdict struct{ result, detail string }

func pass() verdict          { return verdict{resPass, ""} }
func failv(d string) verdict { return verdict{resFail, d} }
func noRun(what string) verdict {
	return failv(notRun + what)
}

// fullCheck is the parsed input the ten rules read.
type fullCheck struct {
	in  FullInput
	prd prdData
	u   universe
}

func (c *fullCheck) judge(rule string) verdict {
	switch rule {
	case "covers":
		return c.covers()
	case "needs":
		return c.needs()
	case "measures":
		return c.measures()
	case "criteria":
		return c.criteria()
	case "classes":
		return c.classes()
	case "inventory":
		return c.inventory()
	case "quotes":
		return c.quotes()
	case "sections":
		return c.sections()
	case "tasks":
		return c.tasks()
	}
	return failv("no rule " + rule)
}

// universe is the set of fact IDs the PRD may cite — the numbered facts of
// the problem statement and the answers' Q-NNNs — with the class and scope of
// each. Before the first bet it is spec/spans.tsv and the answers; after, the
// latest version of the confirmed inventory (the version the index's
// greatest n names) and the answers. A latest version that is absent or
// unreadable leaves notWhy set: the rules that need the universe did not run.
type universe struct {
	facts  map[string]bool
	class  map[string]string // the universe facts' classes; an answer fact has none
	out    map[string]bool   // the facts the inventory marks out of scope
	notWhy string
}

func universeOf(in FullInput) universe {
	u := universe{facts: map[string]bool{}, class: map[string]string{}, out: map[string]bool{}}
	for _, a := range in.Answers {
		u.facts[a[aQuestion]] = true
	}
	if len(in.Index) == 0 {
		if in.Spans == nil {
			u.notWhy = "spec/spans.tsv"
			return u
		}
		for _, s := range in.Spans {
			if s[colFact] != "" {
				u.facts[s[colFact]] = true
				u.class[s[colFact]] = s[colClass]
			}
		}
		return u
	}
	latest := 0
	for _, r := range in.Index {
		if n, err := strconv.Atoi(r[idxN]); err == nil && n > latest {
			latest = n
		}
	}
	b, ok := in.Versions[latest]
	if !ok || b == nil {
		u.notWhy = fmt.Sprintf("version %d", latest)
		return u
	}
	rows, err := ReadInventory(b)
	if err != nil {
		u.notWhy = fmt.Sprintf("version %d", latest)
		return u
	}
	for _, r := range rows {
		u.facts[r[invFact]] = true
		u.class[r[invFact]] = r[invClass]
		if r[invScope] == "out" {
			u.out[r[invFact]] = true
		}
	}
	return u
}

// resolveFact resolves one citation token to a fact ID of the universe
// (review condition C7): F-<a record of the problem statement's facts>#n is
// the numbered fact n; F-<the answers' record>#N is Q-00N, the question its
// fact N answers; a bare Q-NNN is itself; a bare record ID of one of those
// records resolves to no fact ("", true) — the record itself, which the
// prd-lint fact rule already requires to exist. Anything else is ("", false).
var (
	factCite = regexp.MustCompile(`^(F-[0-9]{4})(?:#([0-9]+))?$`)
	citeTok  = regexp.MustCompile(`\bF-[0-9]{4}(#[0-9]+)?\b|\bQ-[0-9]{3,}\b`)
)

func (c *fullCheck) resolveFact(tok string) (id string, ok bool) {
	if qIDForm.MatchString(tok) {
		return tok, c.u.facts[tok]
	}
	m := factCite.FindStringSubmatch(tok)
	if m == nil {
		return "", false
	}
	rec, num := m[1], m[2]
	if num == "" {
		if rec == c.in.AnswersRecord || slices.Contains(c.in.PSBRecords, rec) || slices.Contains(c.in.Records, rec) {
			return "", true
		}
		return "", false
	}
	switch {
	case slices.Contains(c.in.PSBRecords, rec):
		return num, c.u.facts[num]
	case rec == c.in.AnswersRecord:
		n, err := strconv.Atoi(num)
		if err != nil {
			return "", false
		}
		q := fmt.Sprintf("Q-%03d", n)
		return q, c.u.facts[q]
	}
	return "", false
}

// The parsed content of the head's PRDs: every five-cell REQ/NFR row wherever
// it stands (as the pinned producer reads them), the criteria table, the
// success rows and the waits-for-a-bet extent.
type prdData struct {
	ok         bool
	reqs       []prdReq
	criteria   map[string]string // requirement ID → its acceptance criterion
	statements map[string]string // requirement ID → the Statement cell of its row
	metrics    []string          // the first cell of each row of a success-metrics table
	success    map[string]bool   // raw Fact tokens of the success-criterion rows
	waits      map[string]bool   // raw cite tokens under a `waits for a bet` heading
	cited      []citedTok        // every fact token cited, with its site
}

type citedTok struct{ tok, site string }

type prdReq struct {
	id, mo, phase, file string
	facts               []string
}

var (
	headingLine = regexp.MustCompile(`^#{1,6} `)
	reqIDCell   = regexp.MustCompile(`^(REQ|NFR)-[0-9]{3}$`)
	sepCell     = regexp.MustCompile(`^:?-+:?$`)
)

// PRDContent is the words of the requirement rows of PRD files (task T-h65e, #418): the Statement cell of each row of the requirements
// table by its ID, the acceptance criterion of each row of the criteria table by its requirement, and the first cell of each row of a
// success-metrics table, in the order of the files by name and of the rows.
type PRDContent struct {
	Statements, Criteria map[string]string
	Metrics              []string
}

// ContentOfPRDs reads the words of the PRD files, by the same parse as the full check.
func ContentOfPRDs(prds map[string][]byte) PRDContent {
	p := parsePRDs(prds)
	return PRDContent{Statements: p.statements, Criteria: p.criteria, Metrics: p.metrics}
}

func parsePRDs(prds map[string][]byte) prdData {
	p := prdData{ok: len(prds) > 0, criteria: map[string]string{}, statements: map[string]string{}, success: map[string]bool{}, waits: map[string]bool{}}
	for _, name := range slices.Sorted(maps.Keys(prds)) {
		p.parse(name, string(prds[name]))
	}
	return p
}

func (p *prdData) parse(name, text string) {
	var header []string
	var pending []string // the row right above a separator is the table's header
	var table string     // "criteria" | "success" | "" — the open table's kind
	factCol := -1
	inWaits, inMetrics := false, false
	for _, raw := range strings.Split(text, "\n") {
		l := strings.TrimRight(raw, " \r")
		if headingLine.MatchString(l) {
			level := len(l) - len(strings.TrimLeft(l, "#"))
			h := strings.TrimSpace(l[level:])
			inWaits = strings.EqualFold(h, "waits for a bet")
			if level <= 2 {
				inMetrics = level == 2 && strings.Contains(h, "Success metrics")
			}
			header, pending, table = nil, nil, ""
			continue
		}
		if inWaits {
			for _, tok := range citeTok.FindAllString(l, -1) {
				p.waits[tok] = true
				p.cited = append(p.cited, citedTok{tok, "a requirement marked waits for a bet"})
			}
		}
		if !strings.HasPrefix(strings.TrimSpace(l), "|") {
			header, pending, table = nil, nil, ""
			continue
		}
		cells := tableCells(l)
		if isSepRow(cells) {
			header, table, factCol = pending, "", -1
			if len(header) == 2 && header[0] == "Requirement" && header[1] == "Acceptance criterion" {
				table = "criteria"
			} else if inMetrics {
				for i, h := range header {
					if h == "Fact" {
						table, factCol = "success", i
					}
				}
			}
			continue
		}
		pending = cells
		if len(cells) == 5 && reqIDCell.MatchString(cells[0]) {
			r := prdReq{id: cells[0], mo: cells[2], phase: cells[3], facts: factTokens(cells[4]), file: name}
			p.reqs = append(p.reqs, r)
			p.statements[r.id] = cells[1]
			for _, tok := range r.facts {
				p.cited = append(p.cited, citedTok{tok, "the Facts cell of " + r.id})
			}
		}
		switch table {
		case "criteria":
			// A malformed row (fewer than two cells) records no criterion:
			// the requirement it truncates then fails the rule, never a panic.
			if len(cells) >= 2 {
				p.criteria[strings.Trim(cells[0], "`")] = cells[1]
			}
		case "success":
			if factCol >= 0 && factCol < len(cells) {
				p.metrics = append(p.metrics, cells[0])
				for _, tok := range factTokens(cells[factCol]) {
					p.success[tok] = true
					p.cited = append(p.cited, citedTok{tok, "a success-criterion row"})
				}
			}
		}
	}
}

// tableCells splits a Markdown table row into its cells, as the pinned
// producer's split_row does: \| escapes a bar inside a cell.
func tableCells(l string) []string {
	l = strings.ReplaceAll(l, `\|`, "\x01")
	parts := strings.Split(l, "|")
	if len(parts) < 3 {
		return nil
	}
	cells := make([]string, 0, len(parts)-2)
	for _, c := range parts[1 : len(parts)-1] {
		cells = append(cells, strings.ReplaceAll(strings.TrimSpace(c), "\x01", "|"))
	}
	return cells
}

func isSepRow(cells []string) bool {
	if len(cells) == 0 {
		return false
	}
	for _, c := range cells {
		if !sepCell.MatchString(c) {
			return false
		}
	}
	return true
}

// factTokens gives the tokens of a Facts or Fact cell: cite forms only;
// anything else in the cell is a token that resolves to no fact.
func factTokens(cell string) []string {
	var toks []string
	for _, t := range strings.FieldsFunc(cell, func(r rune) bool { return r == ',' || r == ' ' }) {
		if t = strings.Trim(t, "`"); t != "" {
			toks = append(toks, t)
		}
	}
	return toks
}

// factLess orders fact IDs for a stable first-broken-place: numbers
// numerically, other IDs lexically.
func factLess(a, b string) int {
	ai, ae := strconv.Atoi(a)
	bi, be := strconv.Atoi(b)
	if ae == nil && be == nil {
		return ai - bi
	}
	return strings.Compare(a, b)
}

// citedSet is the set of fact IDs the requirements' Facts cells cover.
func (c *fullCheck) citedSet() map[string]bool {
	covered := map[string]bool{}
	for _, r := range c.prd.reqs {
		for _, tok := range r.facts {
			if id, ok := c.resolveFact(tok); ok && id != "" {
				covered[id] = true
			}
		}
	}
	return covered
}

// covers: every fact token the PRD cites — a requirement's Facts cell, a
// success-criterion row's Fact cell, a waits-for-a-bet line — resolves to a
// fact of the universe.
func (c *fullCheck) covers() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if c.u.notWhy != "" {
		return noRun(c.u.notWhy)
	}
	for _, ct := range c.prd.cited {
		if _, ok := c.resolveFact(ct.tok); !ok {
			return failv(fmt.Sprintf("%s cites %s, which is no fact of the inventory or the answers", ct.site, ct.tok))
		}
	}
	return pass()
}

// needs: an in-scope need or constraint fact has a requirement that covers
// it, is out of scope, or is named by a requirement marked waits for a bet.
func (c *fullCheck) needs() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if c.u.notWhy != "" {
		return noRun(c.u.notWhy)
	}
	covered := c.citedSet()
	var facts []string
	for f, class := range c.u.class {
		if (class == "need" || class == "constraint") && !c.u.out[f] && !covered[f] && !c.waitsCover(f) {
			facts = append(facts, f)
		}
	}
	if len(facts) > 0 {
		slices.SortFunc(facts, factLess)
		return failv(fmt.Sprintf("the %s fact %s has no requirement", c.u.class[facts[0]], facts[0]))
	}
	return pass()
}

// measures: an in-scope measure fact has a success-criterion row that cites
// its ID — a row of a table under a `## … Success metrics` heading whose
// header holds a Fact column; a cite in a requirement's acceptance cell is no
// success row (review condition C8), is out of scope, or is named by a
// requirement marked waits for a bet.
func (c *fullCheck) measures() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if c.u.notWhy != "" {
		return noRun(c.u.notWhy)
	}
	cited := map[string]bool{}
	for tok := range c.prd.success {
		if id, ok := c.resolveFact(tok); ok && id != "" {
			cited[id] = true
		}
	}
	var facts []string
	for f, class := range c.u.class {
		if class == "measure" && !c.u.out[f] && !cited[f] && !c.waitsCover(f) {
			facts = append(facts, f)
		}
	}
	if len(facts) > 0 {
		slices.SortFunc(facts, factLess)
		return failv(fmt.Sprintf("the measure fact %s has no success-criterion row", facts[0]))
	}
	return pass()
}

// waitsCover reports whether a waits-for-a-bet section names the fact.
func (c *fullCheck) waitsCover(f string) bool {
	for tok := range c.prd.waits {
		if id, ok := c.resolveFact(tok); ok && id == f {
			return true
		}
	}
	return false
}

// criteria: every requirement has a row of the criteria table — a table whose
// header is Requirement | Acceptance criterion — with a non-empty criterion.
func (c *fullCheck) criteria() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	for _, r := range c.prd.reqs {
		if c.prd.criteria[r.id] == "" {
			return failv(fmt.Sprintf("the requirement %s has no criterion", r.id))
		}
	}
	return pass()
}

// classes: each requirement's MoSCoW and Phase equal the latest copy's —
// where a copy line is a REQ/NFR-headed line, parsed by the one grammar the
// bet's answer uses (parseReqLine, C6). A PRD row with an empty MoSCoW and an
// empty or — Phase and no copy line is a requirement waiting for its next bet
// and agrees (C5); any other row without a copy line fails.
func (c *fullCheck) classes() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if len(c.in.Copy) == 0 {
		return noRun("the copy of the latest bet")
	}
	copy := map[string]ReqAnswer{}
	for _, raw := range strings.Split(string(c.in.Copy), "\n") {
		l := strings.TrimSpace(raw)
		if l == "" || !reqHead.MatchString(l) {
			continue
		}
		id, a, ok := parseReqLine(l)
		if !ok {
			return failv(fmt.Sprintf("the copy's line %q is no requirement line", l))
		}
		copy[id] = a
	}
	have := map[string]bool{}
	for _, r := range c.prd.reqs {
		have[r.id] = true
		a, ok := copy[r.id]
		if !ok {
			if r.mo == "" && (r.phase == "" || r.phase == "—") {
				continue
			}
			return failv(fmt.Sprintf("the requirement %s of %s has MoSCoW %q and Phase %q, and the copy has no line for it", r.id, r.file, r.mo, r.phase))
		}
		phase := r.phase
		if phase == "—" {
			phase = ""
		}
		if a.Class != r.mo || a.Phase != phase {
			return failv(fmt.Sprintf("the requirement %s of %s has MoSCoW %q Phase %q; the copy has %q phase %q", r.id, r.file, r.mo, r.phase, a.Class, a.Phase))
		}
	}
	for _, id := range slices.Sorted(maps.Keys(copy)) {
		if !have[id] {
			return failv(fmt.Sprintf("the copy names the requirement %s, which is no row of the PRD", id))
		}
	}
	return pass()
}

// inventory: each version's exact bytes hash to the sha256 of its index row
// and read as an inventory.
func (c *fullCheck) inventory() verdict {
	for i, r := range c.in.Index {
		n, err := strconv.Atoi(r[idxN])
		if err != nil {
			return failv(fmt.Sprintf("the index's row %d: n %q is no number", i+1, r[idxN]))
		}
		b, ok := c.in.Versions[n]
		if !ok {
			// An absent required input did not run the rule; ran fails it
			// again (the CheckFull contract, and the universe's own signal).
			return noRun(fmt.Sprintf("version %d", n))
		}
		if sha := fmt.Sprintf("%x", sha256.Sum256(b)); sha != r[idxSHA] {
			return failv(fmt.Sprintf("version %d: the sha256 of its bytes differs from the index", n))
		}
		if _, err := ReadInventory(b); err != nil {
			return failv(fmt.Sprintf("version %d: %v", n, err))
		}
	}
	return pass()
}

// quotes: a requirement that covers an answer's fact — a bare Q-NNN or
// F-<answers>#N — also names the fact of that answer's quote, resolved
// through the questions' chain to the one span fact that holds it.
func (c *fullCheck) quotes() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if c.u.notWhy != "" {
		return noRun(c.u.notWhy)
	}
	byID := map[string][]string{}
	for _, q := range c.in.Questions {
		byID[q[qID]] = q
	}
	for _, r := range c.prd.reqs {
		covered := map[string]bool{}
		for _, tok := range r.facts {
			if id, ok := c.resolveFact(tok); ok {
				covered[id] = true
			}
		}
		for _, tok := range r.facts {
			id, ok := c.resolveFact(tok)
			if !ok || !qIDForm.MatchString(id) {
				continue
			}
			if c.in.Questions == nil || c.in.Spans == nil || c.in.PSB == nil {
				return noRun("the quote's chain (questions.tsv, spec/spans.tsv, the problem statement)")
			}
			origin, f := quoteQuestion(id, byID)
			if f != nil {
				return failv(fmt.Sprintf("the requirement %s: %s", r.id, f.Detail))
			}
			s, e := origin[qStart], origin[qEnd]
			if s == "" && e == "" {
				continue // genuinely no quote
			}
			start, e1 := strconv.Atoi(s)
			end, e2 := strconv.Atoi(e)
			if s == "" || e == "" || e1 != nil || e2 != nil {
				return failv(fmt.Sprintf("the requirement %s: the question %s has quote offsets that are no pair of numbers", r.id, id))
			}
			fact, err := QuoteFactOf(start, end, c.in.Spans, c.in.PSB)
			if err != nil {
				return failv(fmt.Sprintf("the requirement %s: %v", r.id, err))
			}
			if !covered[fact] {
				return failv(fmt.Sprintf("the requirement %s covers the answer's fact %s and not the fact %s of its quote", r.id, id, fact))
			}
		}
	}
	return pass()
}

// The register's columns this check reads.
const (
	tID   = 0
	tReqs = 3
)

// sections: a delivered requirement — one every register task that names it
// is merged, or one no task names — has a non-empty docs/spec/ section whose
// heading holds its ID.
func (c *fullCheck) sections() verdict {
	if !c.prd.ok {
		return noRun("no PRD")
	}
	if len(c.in.Register) == 0 {
		return noRun("tasks.tsv")
	}
	tasks, err := tsv.Read(c.in.Register, tasksSchema)
	if err != nil {
		return noRun("tasks.tsv")
	}
	naming := map[string][]string{}
	for _, r := range tasks {
		for _, req := range strings.Fields(r[tReqs]) {
			naming[req] = append(naming[req], r[tID])
		}
	}
	for _, r := range c.prd.reqs {
		delivered := true
		for _, task := range naming[r.id] {
			if !slices.Contains(c.in.Merged, task) {
				delivered = false
			}
		}
		if delivered && !HasSection(c.in.SpecFiles, r.id) {
			return failv(fmt.Sprintf("the delivered requirement %s has no non-empty docs/spec/ section whose heading holds its ID", r.id))
		}
	}
	return pass()
}

// HasSection reports whether a docs/spec/ file of files, by name, holds a
// section whose heading contains the ID, with at least one non-blank line
// before the next heading. It is the one predicate of sections and of the
// handoff item spec-sections (task T-yeav, #404); no ID is no section.
func HasSection(files map[string][]byte, id string) bool {
	if id == "" {
		return false
	}
	re := regexp.MustCompile(`(^|[^0-9A-Za-z])` + regexp.QuoteMeta(id) + `([^0-9A-Za-z]|$)`)
	for _, name := range slices.Sorted(maps.Keys(files)) {
		open, body := false, false
		for _, l := range strings.Split(string(files[name]), "\n") {
			if headingLine.MatchString(l) {
				if open && body {
					return true
				}
				open, body = re.MatchString(l), false
				continue
			}
			if open && strings.TrimSpace(l) != "" {
				body = true
			}
		}
		if open && body {
			return true
		}
	}
	return false
}

// tasks: every task of the register names a requirement — the column
// requirements is never —.
func (c *fullCheck) tasks() verdict {
	if len(c.in.Register) == 0 {
		return noRun("tasks.tsv")
	}
	rows, err := tsv.Read(c.in.Register, tasksSchema)
	if err != nil {
		return noRun("tasks.tsv")
	}
	for _, r := range rows {
		if strings.TrimSpace(r[tReqs]) == "" {
			return failv(fmt.Sprintf("the task %s names no requirement", r[tID]))
		}
	}
	return pass()
}
