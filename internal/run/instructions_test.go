package run

import (
	"embed"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/session"
	"github.com/pharzam/layup/internal/spec"
)

// The golden files are the oracle of the composed bytes of each step's prompt
// file (loop.md, The instruction of a task step). They are embedded, so the unit
// level reads no disk, and no flag rewrites one: a test that rewrites its own
// oracle defeats the freeze.
//
//go:embed testdata/prompts/*.txt
var promptGoldens embed.FS

var time0 = time.Date(2026, 10, 18, 0, 0, 0, 0, time.UTC)

const (
	sha40   = "0123456789abcdef0123456789abcdef01234567"
	sha64   = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	rootDir = "/s/S-aaaaaaaa"
)

var optionsTable = []byte("option\tconstraint\tresult\tevidence\npg\tR-001\tpass\tbench\nsqlite\tR-001\tpass\tbench\n")

// stepSrc is the source of the prompt file of a step of a task of the kind, with
// the inputs that inputsWanted says the step receives.
func stepSrc(t *testing.T, step, kind string) promptSrc {
	t.Helper()
	src := baseSrc(t, step)
	src.Kind, src.Issue, src.Base = kind, 701, sha40
	src.Extras.From = "fedcba9876543210fedcba9876543210fedcba98"
	src.Extras.Frozen = &kept{"tasks/T-ab12/tests.tsv", []byte("tests\n")}
	shape, options, ledger := inputsWanted(step, kind)
	if shape {
		src.Extras.Shape = []kept{{"spec/draft.tsv", []byte("draft\n")}, {"spec/spans.tsv", []byte("spans\n")},
			{"start/problem-statement.md", []byte("brief\n")}, {"intake/answers.tsv", []byte("answers\n")}}
	}
	if options {
		src.Extras.Options = []kept{{"spec/options/db.tsv", optionsTable}}
		src.Extras.Survivors = []survivor{{"db", []string{"pg", "sqlite"}}}
	}
	if ledger {
		src.Extras.Ledger = []kept{{"sessions.tsv", []byte("sessions\n")}, {"telemetry.tsv", []byte("telemetry\n")}}
	}
	if step == "answer-to-a-question" {
		src.Asking = []kept{{"tasks/T-ab12/results/S-cccccccc.tsv", []byte("asked\n")}}
		src.Questions = []string{"S-cccccccc.1", "S-cccccccc.2"}
	}
	return src
}

// The eight steps with an instruction, each with the kind of task it is shown for (row 66z adds answer-to-a-question).
var goldenCases = []struct{ step, kind string }{
	{"plan-of-a-task", "specification"}, {"plan-review", "architecture"}, {"specification", "specification"},
	{"shape-architecture", "architecture"}, {"implement", "build"}, {"verification", "build"}, {"close-out", "build"},
	{"answer-to-a-question", "specification"},
}

func promptText(t *testing.T, step, kind string) string {
	t.Helper()
	c, err := buildPrompt(stepSrc(t, step, kind))
	if err != nil {
		t.Fatalf("buildPrompt(%s): %v", step, err)
	}
	return string(c.Text(rootDir))
}

func TestThePromptFileOfEachStepIsItsGoldenFile(t *testing.T) {
	for _, c := range goldenCases {
		t.Run(c.step, func(t *testing.T) {
			file := "testdata/prompts/" + c.step + ".txt"
			want, err := promptGoldens.ReadFile(file)
			if err != nil {
				t.Fatal(err)
			}
			got := promptText(t, c.step, c.kind)
			if got != string(want) {
				i := 0
				for i < len(got) && i < len(want) && got[i] == want[i] {
					i++
				}
				t.Errorf("Text differs from %s: got %d bytes, want %d; first difference at byte %d: got %q, want %q",
					file, len(got), len(want), i, got[i:min(len(got), i+40)], string(want[i:min(len(want), i+40)]))
			}
		})
	}
}

var leftover = regexp.MustCompile(`\{[a-z-]+\}`)

func TestNoPlaceholderIsLeftInAnyPrompt(t *testing.T) {
	for _, c := range goldenCases {
		if m := leftover.FindString(promptText(t, c.step, c.kind)); m != "" {
			t.Errorf("the prompt of %s keeps the placeholder %s", c.step, m)
		}
	}
}

// loop.md, The instruction of a task step: the paths and the head are filled by code.
func TestTheVerifierIsToldTheHeadThatNothingMaterialReads(t *testing.T) {
	text := promptText(t, "verification", "build")
	for _, want := range []string{"\nhead: " + sha40 + "\n", " from fedcba9876543210fedcba9876543210fedcba98 to " + sha40 + ",",
		"the file " + rootDir + "/result/review-record.md,", "the result file " + rootDir + "/result/result.tsv:"} {
		if !strings.Contains(text, want) {
			t.Errorf("the prompt of the verifier does not hold %q", want)
		}
	}
}

func TestTheLabelsOfTheResultInstructionAreTheOwnerMaps(t *testing.T) {
	if got := labelList(map[string]string{"zeta-kind": "A", "alpha-kind": "B"}); got != "alpha-kind, zeta-kind" {
		t.Errorf("labelList = %q; want the sorted keys, joined by a comma and a space", got)
	}
	if got := labelList(DefaultOwners()); got != "architecture-boundary, domain-or-business, environment-or-infrastructure, interface-contract" {
		t.Errorf("labelList(DefaultOwners()) = %q", got)
	}
	if text := promptText(t, "implement", "build"); !strings.Contains(text, "one of "+labelList(DefaultOwners())+";") {
		t.Errorf("the result instruction does not name the labels of the owner map:\n%s", text)
	}
}

// formLines gives the lines of the form that a comment step shows: the lines
// after the one that ends "line feed:" for as long as they start as a field, a marker or a bullet.
var formLine = regexp.MustCompile(`^((size|open questions|new interfaces|new dependencies|verdict|head|tests|findings):( |$)|- )`)

func formLines(text string) []string {
	var out []string
	in := false
	for _, l := range strings.Split(text, "\n") {
		switch {
		case strings.HasSuffix(l, "line feed:"):
			in = true
		case in && formLine.MatchString(l):
			out = append(out, l)
		case in:
			return out
		}
	}
	return out
}

// shown fills the placeholders of a shown form with values of the right shape.
func shown(lines []string, verdict string, bullets bool) []byte {
	r := strings.NewReplacer("<small or large>", "small", "<one line, or —>", "—", "<approve, approve-with-conditions or reject>", "approve",
		"<material, nothing material in scope or not mergeable, findings recorded>", verdict, "<id>", "t1", "<path>", "internal/a.go", "<line>", "12", "<rule>", "gate", "<text>", "it fails")
	var b strings.Builder
	for _, l := range lines {
		if strings.HasPrefix(l, "- ") && !bullets {
			continue
		}
		b.WriteString(r.Replace(l) + "\n")
	}
	return []byte(b.String())
}

// The form that each instruction shows is read by the parser of the target
// (loop.md, The fields and the blocks of a comment): the instruction cannot drift from it.
func TestTheFormShownInEachInstructionParses(t *testing.T) {
	plan := formLines(promptText(t, "plan-of-a-task", "specification"))
	if len(plan) != 6 {
		t.Fatalf("the form of the plan: %q; want its four fields, tests: and one bullet", plan)
	}
	if p, tests, ok := ParsePlanComment(shown(plan, "", false)); !ok || p.Size != "small" || len(tests) != 0 {
		t.Errorf("the plan of a specification task (an empty tests: block) = %+v, %v, %v; want it to parse", p, tests, ok)
	}
	if _, tests, ok := ParsePlanComment(shown(plan, "", true)); !ok || len(tests) != 1 || tests[0].ID != "t1" {
		t.Errorf("the plan of a build task = %v, %v; want one test", tests, ok)
	}
	if _, _, ok := ParsePlanComment(shown(slices.Delete(slices.Clone(plan), 4, 5), "", false)); ok {
		t.Error("a plan with no tests: marker parses")
	}
	review := formLines(promptText(t, "plan-review", "build"))
	if v, ok := ParsePlanReview(shown(review, "", false)); !ok || v != "approve" {
		t.Errorf("the plan review = %q, %v; want approve", v, ok)
	}
	rec := formLines(promptText(t, "verification", "build"))
	if r, ok := ParseReviewRecord(shown(rec, "nothing material in scope", false)); !ok || len(r.Findings) != 0 || !NothingMaterialOK(shown(rec, "nothing material in scope", false), sha40) {
		t.Errorf("the review record with no finding = %+v, %v; want it to parse and name the head", r, ok)
	}
	if r, ok := ParseReviewRecord(shown(rec, "material", true)); !ok || len(r.Findings) != 1 || r.Findings[0].Line != 12 {
		t.Errorf("the review record with a finding = %+v, %v", r, ok)
	}
	// The words the model must write are in the text, not only in the code.
	text := promptText(t, "verification", "build") + promptText(t, "plan-review", "build")
	for _, w := range []string{"material", "nothing material in scope", "not mergeable, findings recorded", "approve", "approve-with-conditions", "reject", "positive integer", "| Field | Value |"} {
		if !strings.Contains(text, w) {
			t.Errorf("the instructions do not hold %q", w)
		}
	}
}

var quotedRow = regexp.MustCompile(`"([a-z]+<TAB>[^"]*)"`)

// resultOf follows the instruction of a step as a model would: the header and the
// rows that the text quotes, with their placeholders filled. Nothing is invented.
func followedResult(t *testing.T, text, status, sum string, withQuestion bool) []byte {
	t.Helper()
	rows := map[string]string{}
	for _, m := range quotedRow.FindAllStringSubmatch(text, -1) {
		if _, ok := rows[m[1][:strings.Index(m[1], "<TAB>")]]; !ok {
			rows[m[1][:strings.Index(m[1], "<TAB>")]] = m[1]
		}
	}
	r := strings.NewReplacer("<TAB>", "\t", "<n>", "1", "<status>", status, "<reason>", "done", "<sha256>", sha64,
		"<text>", "May it?", "<label>", "interface-contract", "<needs_human>", "no", "<record>", "docs/spec/forge.md")
	lines := []string{rows["kind"], rows["status"]}
	for i, path := range outputsOf(text) { // one artifact row for each required output, with the hash of its file
		row := strings.NewReplacer("<path>", path, "<sha256>", sumOf([]byte(outputContent(path))), "<n>", fmt.Sprint(i+1)).Replace(rows["artifact"])
		lines = append(lines, row)
	}
	if c, ok := rows["comment"]; ok {
		lines = append(lines, c)
	}
	if withQuestion {
		lines = append(lines, rows["question"])
	}
	out := regexp.MustCompile(`<the SHA-256 of [^>]*>`).ReplaceAllString(strings.Join(lines, "\n")+"\n", sum)
	return []byte(r.Replace(out))
}

var (
	outputsSentence = regexp.MustCompile(`The required outputs for the result file are ([^\n]*)`)
	docPath         = regexp.MustCompile(`docs/[A-Za-z0-9_./-]*[A-Za-z0-9_]\.md`)
)

// outputsOf gives the files that the text of a step names as its required outputs
// (the variable ones, "each file ... that you wrote", name no path).
func outputsOf(text string) []string {
	m := outputsSentence.FindStringSubmatch(text)
	if m == nil {
		return nil
	}
	return docPath.FindAllString(m[1], -1)
}

// outputContent is the bytes that a harness commits for a required output.
func outputContent(path string) string { return "the content of " + path + "\n" }

// checkOutputs says whether the artifact rows of a result are exactly the wanted
// paths, each with the hash of its file.
func checkOutputs(rows [][]string, want []string) error {
	var got []string
	for _, r := range rows {
		if r[0] != "artifact" {
			continue
		}
		got = append(got, r[2])
		if r[3] != sumOf([]byte(outputContent(r[2]))) {
			return fmt.Errorf("the artifact %s has the hash %s, not that of its file", r[2], r[3])
		}
	}
	if !slices.Equal(got, want) {
		return fmt.Errorf("the artifacts are %v; want %v", got, want)
	}
	return nil
}

// The result rows that the instructions quote are read by the reader of the block result.
func TestTheResultRowsThatTheInstructionQuotesAreReadByTheResultReader(t *testing.T) {
	for _, c := range goldenCases {
		text := promptText(t, c.step, c.kind)
		rows, err := records.ReadResult(followedResult(t, text, "completed", sha64, false))
		if err != nil || rows[0][0] != "status" {
			t.Errorf("%s: the result that follows the instruction: %v (%v)", c.step, err, rows)
		}
		rows, err = records.ReadResult(followedResult(t, text, "needs_context", sha64, true))
		qs := QuestionsOf(rows)
		if err != nil || len(qs) != 1 || qs[0].Label != "interface-contract" || qs[0].NeedsHuman != "no" || len(qs[0].Records) != 1 {
			t.Errorf("%s: a stopping question that follows the instruction: %v, %+v", c.step, err, qs)
		}
	}
}

func TestTheStepAndTheKindOfTheTaskAgree(t *testing.T) {
	for _, c := range []struct {
		step, kind string
		ok         bool
	}{
		{"specification", "specification", true}, {"specification", "architecture", false}, {"specification", "build", false},
		{"shape-architecture", "architecture", true}, {"shape-architecture", "specification", false},
		{"implement", "build", true}, {"implement", "specification", false}, {"implement", "architecture", false},
		{"plan-review", "build", true}, {"plan-review", "architecture", true}, {"close-out", "specification", true},
	} {
		err := checkStepKind(c.step, c.kind)
		if c.ok != (err == nil) {
			t.Errorf("checkStepKind(%s, %s) = %v; want ok %v", c.step, c.kind, err, c.ok)
		}
		if err != nil && !strings.Contains(err.Error(), "the step "+c.step+" is not a step of a task of kind "+c.kind) {
			t.Errorf("checkStepKind(%s, %s) = %q", c.step, c.kind, err)
		}
	}
}

// The inputs by step and kind (loop.md, The records each step receives).
func TestAShapeTaskReceivesTheRevisedDraftAndItsFactsAndAnArchitectureTaskItsOptions(t *testing.T) {
	for _, step := range []string{"plan-of-a-task", "plan-review", "specification", "shape-architecture", "verification", "implement", "close-out"} {
		for _, kind := range []string{"specification", "architecture", "build"} {
			shape, options, ledger := inputsWanted(step, kind)
			wantShape := step != "close-out" && (kind == "specification" || kind == "architecture")
			if shape != wantShape || options != (wantShape && kind == "architecture") || ledger != (step == "close-out") {
				t.Errorf("inputsWanted(%s, %s) = %v, %v, %v; want %v, %v, %v", step, kind, shape, options, ledger, wantShape, wantShape && kind == "architecture", step == "close-out")
			}
		}
	}
	for _, c := range goldenCases {
		items := itemsOf(mustPrompt(t, c.step, c.kind), session.KindRequirements)
		if has := slices.Contains(items, "spec/draft.tsv"); has != (c.step != "close-out" && c.kind != "build") {
			t.Errorf("%s of a %s task: requirements = %v", c.step, c.kind, items)
		}
	}
}

func tail(s []string, n int) []string { return s[max(0, len(s)-n):] }

func mustPrompt(t *testing.T, step, kind string) *PromptContext {
	t.Helper()
	c, err := buildPrompt(stepSrc(t, step, kind))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Each new input is in its own kind, after the items that exist, in the order of loop.md.
func TestTheNewInputsAreInTheirKindsAfterTheExistingItems(t *testing.T) {
	arch := mustPrompt(t, "shape-architecture", "architecture")
	if got := itemsOf(arch, session.KindRequirements); !slices.Equal(tail(got, 4), []string{"spec/draft.tsv", "spec/spans.tsv", "start/problem-statement.md", "intake/answers.tsv"}) {
		t.Errorf("requirements = %v; want the four Intake records last, in order", got)
	}
	if got := itemsOf(arch, session.KindDecisionRecords); !slices.Equal(tail(got, 1), []string{"spec/options/db.tsv"}) || len(got) != 4 {
		t.Errorf("decision-records = %v; want the ADRs then the options table", got)
	}
	if !strings.Contains(promptText(t, "shape-architecture", "architecture"), "the survivors of the filter:\n- db: pg sqlite\n\nhandoffs-and-answers:") {
		t.Error("the survivors are not the text under the decision records")
	}
	if got := itemsOf(mustPrompt(t, "verification", "build"), session.KindIssueAndPlan); !slices.Equal(got, []string{"tasks.tsv", "copies/701-1.md", "tasks/T-ab12/tests.tsv"}) {
		t.Errorf("issue-and-plan = %v; want the frozen list after the plan", got)
	}
	if got := itemsOf(mustPrompt(t, "close-out", "build"), session.KindHandoffsAndAnswers); !slices.Equal(tail(got, 2), []string{"sessions.tsv", "telemetry.tsv"}) {
		t.Errorf("handoffs-and-answers = %v; want the ledger after the answers", got)
	}
	src := stepSrc(t, "implement", "build")
	src.Extras.Frozen = nil
	c, _ := buildPrompt(src)
	if got := itemsOf(c, session.KindIssueAndPlan); !slices.Equal(got, []string{"tasks.tsv", "copies/701-1.md"}) {
		t.Errorf("a task with no frozen list: %v; want no item and no error", got)
	}
}

// Two records of one base name are refused, as before, for each new input.
func TestTwoNewInputsOfOneBaseNameAreRefused(t *testing.T) {
	for _, name := range []string{"draft", "spans", "answers", "tests", "tasks", "problem-statement"} {
		src := stepSrc(t, "shape-architecture", "architecture")
		src.Extras.Options = []kept{{"spec/options/" + name + ".tsv", optionsTable}}
		ext := ".tsv"
		if name == "problem-statement" {
			ext = ".md"
		}
		src.Extras.Options[0].Path = "spec/options/" + name + ext
		if _, err := buildPrompt(src); err == nil || !strings.Contains(err.Error(), "two records are named "+name+ext) {
			t.Errorf("an options table named %s%s: %v; want the refusal naming both", name, ext, err)
		}
	}
	src := stepSrc(t, "close-out", "build")
	src.Handoffs = append(src.Handoffs, kept{"x/sessions.tsv", []byte("x\n")})
	if _, err := buildPrompt(src); err == nil || !strings.Contains(err.Error(), "two records are named sessions.tsv") {
		t.Errorf("a handoff named sessions.tsv beside the ledger: %v", err)
	}
}

func resultOfOptions(names []string, survivors ...string) spec.Result {
	var r spec.Result
	for _, n := range names {
		r.Options = append(r.Options, spec.Option{Name: n})
	}
	for _, n := range survivors {
		r.Survivors = append(r.Survivors, spec.Option{Name: n})
	}
	return r
}

func TestTheResultsOfTheFilterAreCheckedAgainstTheirTables(t *testing.T) {
	read := func(p string) ([]byte, error) {
		switch p {
		case "spec/options/db.tsv":
			return optionsTable, nil
		case "spec/options/bad.tsv":
			return []byte("option\tconstraint\tresult\tevidence\nPG\tR-001\tpass\tx\n"), nil
		}
		return nil, fmt.Errorf("the records have no %s", p)
	}
	ok := resultOfOptions([]string{"pg", "sqlite"}, "pg", "sqlite")
	files, surv, err := optionsKept(map[string]spec.Result{"db": ok}, read)
	if err != nil || len(files) != 1 || files[0].Path != "spec/options/db.tsv" || len(surv) != 1 || !slices.Equal(surv[0].Options, []string{"pg", "sqlite"}) {
		t.Fatalf("a result that agrees with its table: %+v, %+v, %v", files, surv, err)
	}
	if _, surv, err := optionsKept(map[string]spec.Result{"db": resultOfOptions([]string{"pg", "sqlite"}, "sqlite")}, read); err != nil || !slices.Equal(surv[0].Options, []string{"sqlite"}) {
		t.Errorf("the chosen option as the only survivor (after the bet): %+v, %v", surv, err)
	}
	for _, c := range []struct {
		name, want string
		opts       map[string]spec.Result
	}{
		{"a missing table", "the records have no spec/options/cache.tsv", map[string]spec.Result{"cache": ok}},
		{"a renamed option", "decision db: the options of the result are [pg lite]", map[string]spec.Result{"db": resultOfOptions([]string{"pg", "lite"}, "pg")}},
		{"options out of order", "decision db: the options of the result are [sqlite pg]", map[string]spec.Result{"db": resultOfOptions([]string{"sqlite", "pg"}, "pg")}},
		{"no survivor", "decision db: the result has no survivor", map[string]spec.Result{"db": resultOfOptions([]string{"pg", "sqlite"})}},
		{"a survivor outside the table", "decision db: the survivor mongo is not an option of the table", map[string]spec.Result{"db": resultOfOptions([]string{"pg", "sqlite"}, "mongo")}},
		{"a table the reader refuses", "spec/options/bad.tsv: ", map[string]spec.Result{"bad": resultOfOptions([]string{"PG"}, "PG")}},
		{"a decision that is no path word", "decision ../db", map[string]spec.Result{"../db": ok}},
	} {
		if _, _, err := optionsKept(c.opts, read); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %v; want an error holding %q", c.name, err, c.want)
		}
	}
}

// loop.md, The specification output: the forms the specification instruction
// shows are the ones the existing checks read, so a waiting head is accepted.
func specimenBlocks(text string) [][]string {
	var blocks [][]string
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(l, "|") || strings.HasPrefix(l, "## ") && !strings.Contains(l, "—") {
			if n := len(blocks); n > 0 && blocks[n-1] != nil {
				blocks[n-1] = append(blocks[n-1], l)
				continue
			}
			blocks = append(blocks, []string{l})
		} else if n := len(blocks); n > 0 && blocks[n-1] != nil {
			blocks = append(blocks, nil)
		}
	}
	return slices.DeleteFunc(blocks, func(b []string) bool { return b == nil })
}

// specimen gives the PRD that follows the forms of the instruction, the blocks of
// index drop removed (-1 for none).
func specimen(t *testing.T, drop int) string {
	t.Helper()
	text := promptText(t, "specification", "specification")
	r := strings.NewReplacer("<NNN>", "001", "<the requirement>", "Alpha", "<the acceptance criterion>", "a check", "<the success criterion>", "faster")
	var out []string
	for i, b := range specimenBlocks(text) {
		if i == drop {
			continue
		}
		for _, l := range b {
			n := "1"
			if strings.Contains(l, "<the success criterion>") {
				n = "2"
			}
			out = append(out, strings.ReplaceAll(r.Replace(l), "<n>", n))
		}
		out = append(out, "")
	}
	return "# PRD-0001. A specimen\n\n" + strings.Join(out, "\n")
}

func specimenInput(t *testing.T, prd string) spec.FullInput {
	reg := &strings.Builder{}
	if err := WriteRegister(reg, []Task{{ID: "T-ab12", Issue: 3, Kind: "specification", Requirements: []string{"REQ-001"}, Opened: time0}}); err != nil {
		t.Fatal(err)
	}
	return spec.FullInput{Changed: []string{"docs/prd/PRD-0001-x.md"}, PSBRecords: []string{"F-0004"}, Records: []string{"F-0003"},
		PRDs: map[string][]byte{"PRD-0001-x.md": []byte(prd)}, SpecFiles: map[string][]byte{"loop.md": []byte("## REQ-001 — Alpha\n\nThe text.\n")},
		PSB: []byte(strings.Repeat("x", 100)), Spans: [][]string{{"0", "10", "1", "need", "", "s1"}, {"10", "20", "2", "measure", "", "s2"}},
		Questions: [][]string{}, Answers: [][]string{}, Register: []byte(reg.String())}
}

func TestTheSpecimenOfTheSpecificationInstructionPassesTheFullCheck(t *testing.T) {
	if len(specimenBlocks(promptText(t, "specification", "specification"))) != 4 {
		t.Fatalf("the specification instruction shows %d forms; want the requirements, the criteria, the success metrics and the matrix", len(specimenBlocks(promptText(t, "specification", "specification"))))
	}
	for _, row := range spec.CheckFull(specimenInput(t, specimen(t, -1))) {
		if row[1] == "fail" {
			t.Errorf("the rule %s fails on the specimen: %s", row[0], row[2])
		}
	}
	// A control: without its acceptance criteria the same specimen fails the rule criteria.
	var failed []string
	for _, row := range spec.CheckFull(specimenInput(t, specimen(t, 1))) {
		if row[1] == "fail" {
			failed = append(failed, row[0])
		}
	}
	if !slices.Contains(failed, "criteria") {
		t.Errorf("a specimen without its acceptance criteria fails %v; want criteria among them", failed)
	}
}

// The verifier reviews from the base of the last attempt, not the first.
// The ID of the next numbered facts record, and the records of docs/facts/ that prd-lint reads
// (its pattern: F-NNNN, anything, .md), have one home.
func TestTheFactsRecordsAreThoseOfPrdLintAndTheNextIDIsOneMore(t *testing.T) {
	names := []string{"F-0001-a.md", "F-0003-setup-answers.md", "README.md", "F-00044-x.md", "F-0002.md", "problem-statement-brief.md", "F-0009-x.txt", "f-0010-x.md"}
	ids := factRecordIDs(names)
	if want := []string{"F-0001", "F-0002", "F-0003", "F-0004"}; !slices.Equal(ids, want) {
		t.Errorf("factRecordIDs = %v; want %v", ids, want)
	}
	if got := nextFactsID(ids); got != "F-0005" {
		t.Errorf("nextFactsID = %s; want one more than the highest", got)
	}
	if got := nextFactsID(nil); got != "F-0001" {
		t.Errorf("nextFactsID of no record = %s; want F-0001", got)
	}
	in := spec.FullInput{}
	completeFactsRecords(&in, names)
	if !slices.Equal(in.Records, ids) || !slices.Equal(in.PSBRecords, []string{"F-0005"}) {
		t.Errorf("before the bet: Records %v, PSBRecords %v; want the records and the predicted ID", in.Records, in.PSBRecords)
	}
	in = spec.FullInput{PSBRecords: []string{"F-0007"}}
	completeFactsRecords(&in, append(names, "F-0007-numbered-facts.md"))
	if !slices.Equal(in.PSBRecords, []string{"F-0007"}) {
		t.Errorf("with a numbered facts file: PSBRecords %v; want it unchanged", in.PSBRecords)
	}
	in = spec.FullInput{Index: [][]string{{"1", sha64, "M2e"}}}
	completeFactsRecords(&in, names)
	if in.PSBRecords != nil {
		t.Errorf("after the first bet with no numbered facts file: PSBRecords %v; want none (the prediction is for the head before the bet)", in.PSBRecords)
	}
}

func TestTheChangeOfTheVerifierStartsAtTheBaseOfTheLastAttempt(t *testing.T) {
	events := [][]string{{"1", "attempt", "1", "", "aaaa", "", "", ""}, {"2", "bound", "1", "S-1", "aaaa", "bbbb", "", ""}, {"3", "attempt", "2", "", "cccc", "", "", ""}, {"4", "handoff", "2", "S-2", "", "", "", ""}}
	if got, ok := attemptBase(events); !ok || got != "cccc" {
		t.Errorf("attemptBase = %q, %v; want the base of attempt 2", got, ok)
	}
	if got, ok := attemptBase(events[1:2]); ok {
		t.Errorf("attemptBase of a task with no attempt = %q, %v; want none", got, ok)
	}
}

// Every required output of a step is a row of the result, with its real hash:
// the close-out has both of its files, and a result without one is refused.
func TestTheResultNamesEveryRequiredOutputWithItsHash(t *testing.T) {
	want := map[string][]string{"plan-of-a-task": nil, "plan-review": nil, "verification": nil, "specification": {"docs/tasks/T-ab12.md"},
		"shape-architecture": {"docs/tasks/T-ab12.md"}, "implement": {"docs/tasks/T-ab12.md"}, "close-out": {"docs/tasks/T-ab12.md", "docs/tasks/completed.md"}}
	for _, c := range goldenCases {
		text := promptText(t, c.step, c.kind)
		if got := outputsOf(text); !slices.Equal(got, want[c.step]) {
			t.Errorf("%s: the required outputs of the text are %v; want %v", c.step, got, want[c.step])
			continue
		}
		rows, err := records.ReadResult(followedResult(t, text, "completed", sha64, false))
		if err != nil {
			t.Fatal(err)
		}
		if err := checkOutputs(rows, want[c.step]); err != nil {
			t.Errorf("%s: %v", c.step, err)
		}
		for i := range want[c.step] { // the contract is not met by a result that lacks one output
			rest := slices.Concat(want[c.step][:i], want[c.step][i+1:])
			if checkOutputs(rows, rest) == nil {
				t.Errorf("%s: a result with the artifact %s is taken for one without it", c.step, want[c.step][i])
			}
		}
	}
}

// U9 (loop.md, The instruction of a task step; row 66z): the block instruction-answers follows the result text only when the prompt is given
// the answers that the attempt before has; the prompt without any is byte for byte the one of every existing golden file (the green negative).
func TestTheAnswersToCiteAreWrittenOnlyWhenTheLastAskingAttemptHasThem(t *testing.T) {
	src := stepSrc(t, "implement", "build")
	without, err := buildPrompt(src)
	if err != nil {
		t.Fatal(err)
	}
	want, err := promptGoldens.ReadFile("testdata/prompts/implement.txt")
	if err != nil || string(without.Text(rootDir)) != string(want) {
		t.Fatalf("the prompt with no answers differs from its golden file (%v)", err)
	}
	src.AnswerIDs = []string{"S-dddddddd.1", "S-dddddddd.2"}
	with, err := buildPrompt(src)
	if err != nil {
		t.Fatal(err)
	}
	text := string(with.Text(rootDir))
	for _, sub := range []string{"Their IDs are S-dddddddd.1 S-dddddddd.2.", "status row", "question row"} {
		if !strings.Contains(text, sub) {
			t.Errorf("the prompt with answers lacks %q", sub)
		}
	}
	if !strings.HasPrefix(text, string(without.Text(rootDir))) {
		t.Error("the prompt with answers does not start with the prompt without them: the block must follow the result text")
	}
	if m := leftover.FindString(text); m != "" {
		t.Errorf("the prompt with answers keeps the placeholder %s", m)
	}
}
