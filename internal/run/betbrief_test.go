package run

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/ledger"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
)

// The brief of the first bet (docs/spec/loop.md, The brief of the first bet; task T-h65e, #418, row 66u): the parts that code writes,
// the checks of the brief, the appetite with the session's own spend, the prompt and the basis.

var (
	headSP = strings.Repeat("4f", 20)
	headAR = strings.Repeat("a1", 20)
	heads  = []headRef{{"T-sp01", "specification", headSP, 7}, {"T-ar01", "architecture", headAR, 8}}
	prUrls = []string{"https://github.com/acme/target/pull/7", "https://github.com/acme/target/pull/8"}
	blob   = "https://github.com/acme/target/blob/" + headAR + "/docs/spec/architecture.md"
)

// dbWeb are the results of the filter: db has two survivors and one drop; web has one survivor and one drop.
func dbWeb() map[string]spec.Result {
	opt := func(n string, reasons ...string) spec.Option {
		return spec.Option{Name: n, Dropped: reasons != nil, Reasons: reasons}
	}
	return map[string]spec.Result{
		"db":  {Options: []spec.Option{opt("pg"), opt("sqlite"), opt("mysql", "NFR-002: fail")}, Survivors: []spec.Option{opt("pg"), opt("sqlite")}},
		"web": {Options: []spec.Option{opt("go"), opt("rust", "Inv-1: fail", "ADR-0003: fail")}, Survivors: []spec.Option{opt("go")}},
	}
}

// qrow gives a row of questions.tsv as the reader gives it: the empty string for —.
func bq(id, source, ref, start, end, state string) []string {
	return []string{id, source, ref, "idea-owner", "before-delivery", source, start, end, "The question " + id, "2026-10-12T09:00:00Z", "layup run", "", state}
}

var spansFx = [][]string{
	{"0", "10", "1", "need", "", "The tool runs offline."},
	{"10", "20", "", "", "a title", "# Title"},
	{"20", "40", "2", "constraint", "", "The data stays in Postgres."},
}

func answersFx() [][]string {
	return [][]string{{"Q-001", "5.0", "801", "1", "2026-10-12T09:00:00Z"}, {"Q-004", "Postgres only", "802", "1", "2026-10-12T09:00:00Z"}}
}

func questionsFx() [][]string {
	return [][]string{bq("Q-001", "intent", "band.B", "", "", "answered"), bq("Q-004", "meaning", "M-001", "22", "30", "answered")}
}

func partsFx(t *testing.T) briefParts {
	t.Helper()
	prd := spec.PRDContent{Statements: map[string]string{"REQ-001": "Alpha", "NFR-002": "p99 under 200 ms"},
		Criteria: map[string]string{"REQ-001": "a check", "NFR-002": "measured at the edge"}, Metrics: []string{"faster"}}
	req, err := requirementsPart([]string{"REQ-001", "NFR-002"}, prd)
	if err != nil {
		t.Fatal(err)
	}
	open := [][]string{bq("Q-002", "intent", "band.U", "", "", "open"), bq("Q-003", "intent", "appetite.hours", "", "", "asked"), bq("Q-007", "meaning", "M-002", "", "", "asked")}
	return briefParts{requirements: req, solution: solutionPart("acme", "target", heads, dbWeb()), open: openPart(open), facts: factsPart(spansFx, questionsFx(), answersFx())}
}

const problemFx, appetiteFx, holesFx = "The problem in prose.\n", "appetite: USD 8, 12 hours\nWhy: small.\n", "No rabbit holes but one.\n"

func filesFx() map[string][]byte {
	return map[string][]byte{"problem.md": []byte(problemFx), "appetite.md": []byte(appetiteFx), "rabbit-holes.md": []byte(holesFx)}
}

func keeperFx(t *testing.T) briefKeeper {
	lines := 60
	return briefKeeper{parts: partsFx(t), multi: map[string][]string{"db": {"pg", "sqlite", "mysql"}}, lines: &lines,
		links: map[string]bool{prUrls[0]: true, prUrls[1]: true, blob: true}}
}

// U1: the solution part is code's bytes.
func TestTheSolutionPartIsCodesBytes(t *testing.T) {
	want := "### The verified heads\n" +
		"- T-sp01 (specification): [pull request #7](https://github.com/acme/target/pull/7), head " + headSP + "\n" +
		"- T-ar01 (architecture): [pull request #8](https://github.com/acme/target/pull/8), head " + headAR + "\n" +
		"### Decision db\n- pg (a survivor)\n- sqlite (a survivor)\n- mysql (dropped: NFR-002: fail)\n" +
		"### Decision web\n- go (the one survivor, decided at Shape)\n- rust (dropped: Inv-1: fail; ADR-0003: fail)\n"
	if got := solutionPart("acme", "target", heads, dbWeb()); got != want {
		t.Errorf("part = %q; want %q", got, want)
	}
}

// U2: every fact with its class and its text, the span facts first, then one fact per answer.
func TestTheFactsPartHoldsEveryFactWithItsClassAndText(t *testing.T) {
	want := "- 1 (need): The tool runs offline.\n- 2 (constraint): The data stays in Postgres.\n- Q-001 (context): 5.0\n- Q-004 (constraint): Postgres only\n"
	if got := factsPart(spansFx, questionsFx(), answersFx()); got != want {
		t.Errorf("part = %q; want %q", got, want)
	}
	if got := factsPart(spansFx, questionsFx(), [][]string{{"Q-001", "see [x](http://a) and [y][z]", "1", "1", "t"}}); !strings.Contains(got, "see [x] (http://a) and [y] [z]\n") {
		t.Errorf("a copied text with link syntax: %q; want no link in it", got)
	}
}

// U3: the class of an answer's fact is that of its quote, inherited by a follow-up and by Shape's row.
func TestTheClassOfAnAnswersFactIsThatOfItsQuote(t *testing.T) {
	questions := [][]string{
		bq("Q-001", "intent", "band.B", "", "", "answered"),
		bq("Q-002", "meaning", "M-001", "22", "30", "answered"),
		bq("Q-003", "follow-up", "Q-002", "", "", "answered"),
		append(bq("Q-004", "follow-up", "Q-002", "", "", "asked")[:5:5], "shape", "", "", "x", "t", "layup run", "", "asked"),
		bq("Q-005", "meaning", "M-002", "50", "55", "answered"),
	}
	for id, want := range map[string]string{"Q-002": "constraint", "Q-003": "constraint", "Q-004": "constraint", "Q-001": "context", "Q-005": "context", "Q-099": "context"} {
		if got := answerFactClass(id, questions, spansFx); got != want {
			t.Errorf("answerFactClass(%s) = %q; want %s", id, got, want)
		}
	}
}

// U4: the questions still open, among them one that Shape asked and that stays open.
func TestTheOpenPartNamesEachQuestionStillOpen(t *testing.T) {
	questions := [][]string{
		bq("Q-002", "intent", "band.U", "", "", "open"),
		bq("Q-003", "intent", "band.B", "", "", "asked"),
		bq("Q-004", "meaning", "M-001", "", "", "asked"),
		bq("Q-005", "meaning", "M-002", "", "", "closed"),
		bq("Q-006", "meaning", "M-003", "", "", "answered"),
		bq("Q-007", "meaning", "M-004", "", "", "asked"),
		append(bq("Q-008", "follow-up", "Q-003", "", "", "asked")[:5:5], "shape", "", "", "x", "t", "layup run", "", "asked"),
		bq("Q-009", "follow-up", "Q-004", "", "", "answered"),
	}
	want := "### The questions still open\n- Q-002: The question Q-002\n- Q-003: The question Q-003\n- Q-007: The question Q-007\n"
	if got := openPart(questions); got != want {
		t.Errorf("part = %q; want %q", got, want)
	}
	if got, want := openPart(nil), "### The questions still open\n- none\n"; got != want {
		t.Errorf("part with no question = %q; want %q", got, want)
	}
}

// U5: the requirements are the PRD's words at the verified head, for exactly the IDs of the register row.
func TestTheRequirementsPartUsesTheRegisteredFinalIDs(t *testing.T) {
	prd := func() spec.PRDContent {
		return spec.PRDContent{Statements: map[string]string{"REQ-001": "Alpha", "NFR-002": "p99 under 200 ms"},
			Criteria: map[string]string{"REQ-001": "a check", "NFR-002": "measured"}, Metrics: []string{"faster", "cheaper"}}
	}
	want := "### The requirements\n- REQ-001: Alpha Criterion: a check\n- NFR-002: p99 under 200 ms Criterion: measured\n### The success criteria\n- faster\n- cheaper\n"
	if got, err := requirementsPart([]string{"REQ-001", "NFR-002"}, prd()); err != nil || got != want {
		t.Errorf("part = %q, %v; want %q", got, err, want)
	}
	p := prd()
	p.Metrics = nil
	if got, err := requirementsPart([]string{"REQ-001", "NFR-002"}, p); err != nil || !strings.HasSuffix(got, "### The success criteria\n- none\n") {
		t.Errorf("part with no metric = %q, %v; want the line - none", got, err)
	}
	// The ID agrees and the draft's words differ: the part has the PRD's words, and the draft is no input of it. A copied text is no link.
	p = prd()
	p.Statements["REQ-001"] = "Alpha see [x](http://a)"
	if got, err := requirementsPart([]string{"REQ-001", "NFR-002"}, p); err != nil || !strings.Contains(got, "- REQ-001: Alpha see [x] (http://a) Criterion: a check\n") {
		t.Errorf("part with link syntax = %q, %v", got, err)
	}
	for name, c := range map[string]struct {
		ids  []string
		edit func(*spec.PRDContent)
		want string
	}{
		"a registered ID with no PRD row": {[]string{"REQ-001", "NFR-002", "REQ-009"}, func(*spec.PRDContent) {}, "the requirement REQ-009 of the register row has no row in the PRD at the verified head"},
		"a PRD row the register lacks":    {[]string{"REQ-001"}, func(*spec.PRDContent) {}, "the PRD at the verified head holds NFR-002, which the register row does not"},
		"a row with no criterion":         {[]string{"REQ-001", "NFR-002"}, func(p *spec.PRDContent) { delete(p.Criteria, "NFR-002") }, "the requirement NFR-002 has no acceptance criterion in the PRD at the verified head"},
		"an empty statement cell":         {[]string{"REQ-001", "NFR-002"}, func(p *spec.PRDContent) { p.Statements["REQ-001"] = "" }, "the requirement REQ-001 has no statement in the PRD at the verified head"},
		"a whitespace-only statement":     {[]string{"REQ-001", "NFR-002"}, func(p *spec.PRDContent) { p.Statements["REQ-001"] = "   " }, "the requirement REQ-001 has no statement in the PRD at the verified head"},
	} {
		p := prd()
		c.edit(&p)
		if _, err := requirementsPart(c.ids, p); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %q", name, err, c.want)
		}
	}
}

// U6: the assembled brief is the golden, and the checker accepts it.
func TestTheAssembledBriefIsTheGoldenAndPassesCheckBrief(t *testing.T) {
	want, err := os.ReadFile("testdata/briefs/first.md")
	if err != nil {
		t.Fatal(err)
	}
	k := keeperFx(t)
	got := assemble(k.parts, problemFx, appetiteFx, holesFx)
	if got != string(want) {
		t.Errorf("assemble differs from testdata/briefs/first.md: got %d bytes, want %d\n%s", len(got), len(want), got)
	}
	if _, findings := spec.CheckBrief([]byte(got), k.lines, k.links, nil); len(findings) != 0 {
		t.Errorf("CheckBrief findings = %v; want none", findings)
	}
}

// U7: each finding refuses the output, and the control is kept. The last case is the known limit of the option-word rule, which is
// not the requirement: it passes the rule and stays noncompliant.
func TestTheBriefIsRefusedForEachFinding(t *testing.T) {
	k := keeperFx(t)
	kept, err := k.keep("1", "S-1a2b3c4d", filesFx())
	if want := assemble(k.parts, problemFx, appetiteFx, holesFx); err != nil || len(kept) != 1 || string(kept["bets/1/attempts/S-1a2b3c4d.md"]) != want {
		t.Errorf("control: kept %d files, %v; want [bets/1/attempts/S-1a2b3c4d.md] and no error", len(kept), err)
	}
	with := func(name, text string) map[string][]byte {
		f := filesFx()
		f[name] = []byte(text)
		return f
	}
	for name, c := range map[string]struct {
		files map[string][]byte
		want  string
	}{
		"a heading line in problem.md":         {with("problem.md", "Text\n## The solution\n"), "file: problem.md has a line that starts with \"## \""},
		"a valid architecture line in problem": {with("problem.md", "architecture change: no\n"), "file: problem.md has a line that starts with \"architecture change:\""},
		"architecture yes in rabbit-holes":     {with("rabbit-holes.md", "architecture change: yes\n"), "file: rabbit-holes.md has a line that starts with \"architecture change:\""},
		"an approve batch line in appetite":    {with("appetite.md", appetiteFx+"approve batch x: y z\n"), "file: appetite.md has a line that starts with \"approve batch \""},
		"no appetite line":                     {with("appetite.md", "Why: small.\n"), "appetite: the line `appetite: USD <n>, <n> hours` is there 0 times, wants one"},
		"a malformed appetite line":            {with("appetite.md", "appetite: lots\n"), "appetite: the line \"appetite: lots\" is not `appetite: USD <n>, <n> hours`"},
		"a link outside the set":               {with("problem.md", "See [it](https://example.invalid/x).\n"), "link: the link to https://example.invalid/x does not resolve"},
		"a reference link":                     {with("problem.md", "See [it][ref].\n"), "link: the link [it][ref] is not of the inline form [text](target)"},
		"over brief.lines":                     {with("problem.md", strings.Repeat("x\n", 40)), "limit: the brief is "},
		"an option of a decision with two":     {with("rabbit-holes.md", "Use pg, it is fast.\n"), "option: the brief names the option pg of db outside code's part"},
		"an option in capitals":                {with("rabbit-holes.md", "Use SQLite.\n"), "option: the brief names the option sqlite of db outside code's part"},
		"over 65536 characters":                {with("problem.md", strings.Repeat("x", 70000)+"\n"), "handoff: the file of a comment has "},
	} {
		if _, err := k.keep("1", "S-1a2b3c4d", c.files); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %q", name, err, c.want)
		}
	}
	// Green negatives beside the cases above: an option of a decision with one survivor, and a word that only contains an option.
	for name, text := range map[string]string{"a one-survivor option": "Use go.\n", "a word that holds an option": "A pgsql note and sqlite3.\n"} {
		if _, err := k.keep("1", "S-1a2b3c4d", with("rabbit-holes.md", text)); err != nil {
			t.Errorf("%s: err = %v; want kept", name, err)
		}
	}
	// The known limit, characterised and not claimed as proof: prose that recommends without an option word passes the rule.
	if _, err := k.keep("1", "S-1a2b3c4d", with("rabbit-holes.md", "I recommend the first one.\n")); err != nil {
		t.Errorf("an indirect recommendation: err = %v; the rule cannot see it (a stated limit)", err)
	}
	// The recovery judges the stored bytes alone: the same findings, and the same pass.
	if findings := k.judge(assemble(k.parts, problemFx, appetiteFx, "Use pg.\n")); len(findings) != 1 || findings[0].Rule != "option" {
		t.Errorf("judge of a stored candidate that names pg = %v; want the option finding", findings)
	}
	if findings := k.judge(assemble(k.parts, problemFx, appetiteFx, holesFx)); len(findings) != 0 {
		t.Errorf("judge of a stored candidate that passes = %v; want none", findings)
	}
}

// telemetry gives a row of telemetry.tsv with the money of the test.
func briefTelemetry(session, money, currency, status string) []string {
	row := make([]string, len(records.TelemetrySchema.Columns))
	for i, c := range records.TelemetrySchema.Columns {
		switch c.Name {
		case "session":
			row[i] = session
		case "money":
			row[i] = money
		case "currency":
			row[i] = currency
		case "money_status":
			row[i] = status
		}
	}
	return row
}

// U8 (condition 4 of the plan review of #366): the appetite counts the brief session's own spend.
func TestTheAppetiteCountsTheSessionsOwnSpend(t *testing.T) {
	brief := func(usd, hours string) string {
		return "## The appetite\nappetite: USD " + usd + ", " + hours + " hours\n"
	}
	before := [][]string{briefTelemetry("S-00000001", "9.0", "USD", "reported")}
	withSession := append(slices.Clone(before), briefTelemetry("S-00000002", "2.0", "USD", "reported"), briefTelemetry("S-00000003", "", "", "unknown"))
	// The control beside the red case: without the session's row the same brief fits (49 of 50).
	if f, err := appetiteFindings(brief("40", "10"), before, "50.0", "40.0"); err != nil || len(f) != 0 {
		t.Errorf("without the session's row: %v, %v; want none (a green negative)", f, err)
	}
	f, err := appetiteFindings(brief("40", "10"), withSession, "50.0", "40.0")
	if want := "the money caps 40 plus the milestone-0 spend 11 are over band.U 50"; err != nil || len(f) != 1 || f[0].Rule != "money" || !strings.Contains(f[0].Detail, want) {
		t.Errorf("with the session's row: %v, %v; want the finding %q", f, err, want)
	}
	if f, err := appetiteFindings(brief("30", "10"), withSession, "50.0", "40.0"); err != nil || len(f) != 0 {
		t.Errorf("USD 30: %v, %v; want none", f, err)
	}
	if f, err := appetiteFindings(brief("5", "41"), withSession, "50.0", "40.0"); err != nil || len(f) != 1 || f[0].Rule != "hours" {
		t.Errorf("hours over: %v, %v; want the hours finding", f, err)
	}
	if f, err := appetiteFindings(brief("60", "10"), withSession, "50.0", "40.0"); err != nil || len(f) != 1 || f[0].Rule != "money" || !strings.Contains(f[0].Detail, "over band.U 50") {
		t.Errorf("a cap above band.U: %v, %v; want the money finding", f, err)
	}
	if f, err := appetiteFindings(brief("1", "10"), withSession, "50.0", "40.0"); err != nil || len(f) != 0 {
		t.Errorf("a cap below band.B: %v, %v; want none", f, err)
	}
	if _, err := appetiteFindings(brief("5", "10"), append(slices.Clone(withSession), briefTelemetry("S-00000004", "1.0", "EUR", "computed")), "50.0", "40.0"); !errors.Is(err, ledger.ErrMilestoneCap) {
		t.Errorf("a known row not in USD: err = %v; want ErrMilestoneCap", err)
	}
}

// U10: the parts of code alone over brief.lines wait before any session.
func TestThePartsOfCodeOverTheLimitWaitBeforeAnySession(t *testing.T) {
	k := keeperFx(t)
	for lines, want := range map[int]string{20: "the parts of code are 29 lines, over brief.lines 20", 29: "", 28: "the parts of code are 29 lines, over brief.lines 28"} {
		k.lines = &lines
		err := k.overLimit()
		switch {
		case want == "" && err != nil:
			t.Errorf("brief.lines %d: err = %v; want none", lines, err)
		case want != "" && (err == nil || err.Error() != want || !errors.Is(err, ErrBriefOverLimit)):
			t.Errorf("brief.lines %d: err = %v; want %q wrapping ErrBriefOverLimit", lines, err, want)
		}
	}
}

// U11: a milestone is one path part, and only the first bet has a brief.
func TestAMilestoneThatIsNoPathPartIsRefused(t *testing.T) {
	for _, m := range []string{"", ".", "..", "a/b", "a b", "1\n"} {
		if err := checkBriefRequest(m, "first"); err == nil || err.Error() != "the milestone "+quotedText(m)+" is not one path part" {
			t.Errorf("checkBriefRequest(%q) = %v; want the milestone refused", m, err)
		}
	}
	for _, m := range []string{"1", "M2e", "a.b-c_d"} {
		if err := checkBriefRequest(m, "first"); err != nil {
			t.Errorf("checkBriefRequest(%q, first) = %v; want none", m, err)
		}
	}
	if err := checkBriefRequest("1", "later"); err == nil || err.Error() != "the brief of a later bet is built by no row of M2e" {
		t.Errorf("kind later: err = %v; want the later bet refused", err)
	}
	if err := checkBriefRequest("1", "other"); err == nil || err.Error() != `the kind "other" is not first or later` {
		t.Errorf("kind other: err = %v; want the kind refused", err)
	}
}

// U12: the basis refuses each broken rule with its line, column and reason.
func basisFx() [][]string {
	return [][]string{
		{"requirement", "REQ-001", "R-001"}, {"requirement", "NFR-002", "R-002"},
		{"survivor", "db=pg", "—"}, {"survivor", "db=sqlite", "—"}, {"dropped", "db=mysql", "NFR-002: fail"}, {"survivor", "web=go", "—"},
		{"fact", "1", "need"}, {"fact", "2", "constraint"}, {"fact", "Q-004", "constraint"},
		{"open", "Q-002", "—"},
		{"head", "T-sp01", headSP}, {"head", "T-ar01", headAR},
	}
}

func TestTheBasisTableRefusesEachBrokenRule(t *testing.T) {
	read := func(edit func([][]string) [][]string) ([][]string, error) {
		data, err := table(BetBasisSchema, edit(basisFx()))
		if err != nil {
			return nil, err
		}
		return ReadBetBasis(data)
	}
	if rows, err := read(func(r [][]string) [][]string { return r }); err != nil || len(rows) != 12 {
		t.Errorf("control: rows = %d, %v; want 12", len(rows), err)
	}
	set := func(i, col int, v string) func([][]string) [][]string {
		return func(r [][]string) [][]string { r[i][col] = v; return r }
	}
	for name, c := range map[string]struct {
		edit func([][]string) [][]string
		want string
	}{
		"a — class":                     {set(6, 2, "—"), `line 8, column "value": a fact's class is need, constraint, measure or context`},
		"a decision with no survivor":   {func(r [][]string) [][]string { r[5] = []string{"dropped", "web=go", "Inv-1: fail"}; return r }, `line 7, column "id": the decision web has no survivor`},
		"a fact twice":                  {func(r [][]string) [][]string { return append(r, []string{"fact", "1", "need"}) }, `line 14: the key (kind "fact", id "1") repeats the key of line 8`},
		"a dropped option with —":       {set(4, 2, "—"), `line 6, column "value": a dropped option holds its reasons, not —`},
		"a survivor with a value":       {set(2, 2, "x"), `line 4, column "value": a survivor holds —`},
		"a head that is no SHA":         {set(10, 2, "abc"), `line 12, column "value": a head is 40 lowercase hexadecimal characters`},
		"a requirement of another row":  {set(0, 2, "R-009"), `line 2, column "value": the row of the draft for REQ-001 is R-001`},
		"an option with no decision":    {set(2, 1, "pg"), `line 4, column "id": a survivor or a dropped option is <decision>=<option>`},
		"an open question not a Q-NNN":  {set(9, 1, "5"), `line 11, column "id": an open question is a Q-NNN`},
		"no head":                       {func(r [][]string) [][]string { return r[:10] }, "the basis has 0 head rows, want two or more"},
		"no requirement":                {func(r [][]string) [][]string { return r[2:] }, "the basis has no requirement row"},
		"no fact":                       {func(r [][]string) [][]string { return append(r[:6:6], r[9:]...) }, "the basis has no fact row"},
		"exactly one head":              {func(r [][]string) [][]string { return r[:11] }, "the basis has 1 head rows, want two or more"},
		"a requirement that is no ID":   {set(0, 1, "REQ1"), `line 2, column "id": a requirement is REQ-NNN or NFR-NNN`},
		"a fact that is no number":      {set(6, 1, "x1"), `line 8, column "id": a fact is a number or a Q-NNN`},
		"a head of no task":             {set(10, 1, "T-xx"), `line 12, column "id": a head is of a task T-xxxx`},
		"an option that is empty":       {set(2, 1, "db="), `line 4, column "id": a survivor or a dropped option is <decision>=<option>`},
		"an open question with a value": {set(9, 2, "x"), `line 11, column "value": an open question holds —`},
	} {
		if _, err := read(c.edit); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %q", name, err, c.want)
		}
	}
	var te *tsv.Error
	if _, err := read(set(6, 2, "—")); !errors.As(err, &te) {
		t.Errorf("a — class: err = %T; want a *tsv.Error with the line and the column", err)
	}
}

// U13: the prompt holds the paths of the files and, after a refusal, the finding.
func TestTheRefusedPromptHoldsTheFinding(t *testing.T) {
	input := map[string]string{}
	for _, n := range []string{"problem-statement.md", "spans.tsv", "draft.tsv", "questions.tsv", "answers.tsv", "budget.tsv", "requirements.md", "solution.md", "open.md", "facts.md", "links.md"} {
		input[n] = "/s/input/" + n
	}
	first := string(briefText("1", 40, 11, "")(input, "/s/result/brief"))
	for _, want := range []string{"the bet of milestone 1, for LAYUP.", "- The requirements and the success criteria, which code writes: /s/input/requirements.md\n", "- The problem statement: /s/input/problem-statement.md\n",
		"in the directory /s/result/brief, and no other file there:", "The whole brief is at most 40 lines, and your three files together have 11 lines at most.", "Then end.\n"} {
		if !strings.Contains(first, want) {
			t.Errorf("the prompt lacks %q:\n%s", want, first)
		}
	}
	if strings.Contains(first, "{") || strings.Contains(first, "refused the last") {
		t.Errorf("the first prompt holds a placeholder or a finding:\n%s", first)
	}
	second := string(briefText("1", 40, 11, "limit: the brief is 45 lines, over the limit 40")(input, "/s/result/brief"))
	if want := "The check refused the last brief: limit: the brief is 45 lines, over the limit 40\nWrite the three files again.\nThen end.\n"; !strings.HasSuffix(second, want) {
		t.Errorf("the prompt after a refusal lacks %q:\n%s", want, second)
	}
}

// The link set is computed from the verified heads, never invented: the pull requests, and the regular files of the heads' changes.
func TestTheLinkSetIsThePullRequestsAndTheChangedFilesOfTheHeads(t *testing.T) {
	blobs := map[string][]string{"T-ar01": {"docs/spec/architecture.md", "docs/a b.md", "docs/x)y.md"}, "T-sp01": {"docs/prd/PRD-0001-x.md"}}
	want := map[string]bool{prUrls[0]: true, prUrls[1]: true, blob: true, "https://github.com/acme/target/blob/" + headSP + "/docs/prd/PRD-0001-x.md": true}
	if got := linkSet("acme", "target", heads, blobs); !sameSet(got, want) {
		t.Errorf("links = %v; want %v", got, want)
	}
}

func sameSet(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func quotedText(s string) string { return `"` + strings.NewReplacer("\n", `\n`).Replace(s) + `"` }

// F1 of round 1: a stored candidate is reused only if it holds the fixed parts of the inputs as they stand now, each in its place.
func TestACandidateWithoutTheFixedPartsOfTheInputsIsRefused(t *testing.T) {
	k := keeperFx(t)
	judged := func(k briefKeeper, edit func(briefParts) briefParts) []spec.Finding {
		return k.judge(assemble(edit(k.parts), problemFx, appetiteFx, holesFx))
	}
	same := func(p briefParts) briefParts { return p }
	if fs := judged(k, same); len(fs) != 0 {
		t.Fatalf("control: a candidate of the current parts: findings %v; want none", fs)
	}
	dropLine := func(s string) string { return s[strings.Index(s, "\n")+1:] }
	for name, edit := range map[string]func(briefParts) briefParts{
		"a requirement part removed": func(p briefParts) briefParts { p.requirements = ""; return p },
		"a requirement's content changed": func(p briefParts) briefParts {
			p.requirements = strings.Replace(p.requirements, "Alpha", "Gamma", 1)
			return p
		},
		"the facts removed":          func(p briefParts) briefParts { p.facts = ""; return p },
		"a fact removed":             func(p briefParts) briefParts { p.facts = dropLine(p.facts); return p },
		"the open questions removed": func(p briefParts) briefParts { p.open = ""; return p },
		"an open question removed":   func(p briefParts) briefParts { p.open = dropLine(p.open); return p },
		"the solution part changed":  func(p briefParts) briefParts { p.solution = strings.Replace(p.solution, "pg", "pgx", 1); return p },
	} {
		if fs := judged(k, edit); len(fs) == 0 || fs[0].Rule != "fixed" {
			t.Errorf("%s: findings %v; want the first of the rule fixed", name, fs)
		}
	}
	// The old candidate, the new facts: a fact came in after the candidate was made.
	k.parts.facts += "- 9 (need): A fact that came after the candidate.\n"
	if fs := k.judge(assemble(partsFx(t), problemFx, appetiteFx, holesFx)); len(fs) == 0 || fs[0].Rule != "fixed" {
		t.Errorf("an old candidate and the new facts: findings %v; want the first of the rule fixed", fs)
	}
}
