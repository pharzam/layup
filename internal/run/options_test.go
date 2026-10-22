package run

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/spec"
)

// The options of the step shape (task T-3f94, #364; docs/spec/loop.md, REQ-012 — Shape).

const optHead = "option\tconstraint\tresult\tevidence\n"

// dbTable has the survivor pg; cacheTable has none (both options fail a row of the list Inv-1, NFR-001).
const (
	dbTable    = optHead + "pg\tInv-1\tpass\tdocs/a.md\npg\tNFR-001\tpass\thttps://x\nsqlite\tInv-1\tpass\tdocs/b.md\nsqlite\tNFR-001\tfail\t—\n"
	cacheTable = optHead + "redis\tInv-1\tfail\t—\nredis\tNFR-001\tpass\tdocs/c.md\nlocal\tInv-1\tpass\tdocs/d.md\n"
)

// out is the output of an options session: the header and a row of each table with its decision first.
func out(tables ...string) []byte {
	var b strings.Builder
	b.WriteString("decision\toption\tconstraint\tresult\tevidence\n")
	for i := 0; i < len(tables); i += 2 {
		for _, line := range strings.SplitAfter(strings.TrimPrefix(tables[i+1], optHead), "\n") {
			if line != "" {
				b.WriteString(tables[i] + "\t" + line)
			}
		}
	}
	return []byte(b.String())
}

func TestTheConstraintsThatApplyAreListedByCodeFromTheRecordsAndTheTarget(t *testing.T) {
	draft := [][]string{{"R-001", "requirement", "The loop", "1", "a check", "S-aaaaaaaa"}, {"R-002", "non-functional", "Fast", "1", "p99 under 200 ms", "S-aaaaaaaa"},
		{"R-003", "success-criterion", "Nine in ten", "1", "by the log", "S-aaaaaaaa"}}
	spans := [][]string{{"0", "9", "", "", "a heading", ""}, {"10", "20", "1", "need", "—", "x"}, {"21", "30", "2", "constraint", "—", "y"}, {"31", "40", "3", "constraint", "—", "z"}}
	adrs := map[string][]byte{
		"0001-a.md":    []byte("# 1. A\n\n## Status\n\nAccepted\n\n## Context\n"),
		"0002-b.md":    []byte("# 2. B\r\n\r\n## Status\r\n\r\nProposed\r\n"),
		"0003-c.md":    []byte("# 3. C\n\n## Status\n\nAccepted. Amended by ADR-0004\n"),
		"0005-d.md":    []byte("# 5. D\n\n## Status\n\nSuperseded by ADR-0006\n"),
		"README.md":    []byte("no status here"),
		"template.md":  []byte("no status here"),
		"12-short.md":  []byte("no status here"),
		"0006-e.txt":   []byte("no status here"),
		"0007-e.md.md": []byte("# 7. E\n\n## Status\n\nAccepted\n"),
	}
	cs, err := constraintsOf(draft, spans, []byte("- **Inv-1** a\n- **Inv-2** b\n"), adrs)
	if err != nil {
		t.Fatal(err)
	}
	want := []spec.Constraint{{ID: "NFR-002", Text: "p99 under 200 ms"}, {ID: "F-2"}, {ID: "F-3"}, {ID: "Inv-1"}, {ID: "Inv-2"}, {ID: "ADR-0001"}, {ID: "ADR-0003"}, {ID: "ADR-0007"}}
	if !reflect.DeepEqual(cs, want) {
		t.Errorf("constraintsOf = %v; want %v (the final ID of the non-functional row, the facts of class constraint, the invariants, the accepted and the amended records)", cs, want)
	}
	for name, adr := range map[string]string{"a record with no status line": "# 4. D\n\nno status\n", "a status heading with no line after it": "# 4. D\n\n## Status\n\n"} {
		_, err := constraintsOf(draft, spans, nil, map[string][]byte{"0004-x.md": []byte(adr)})
		if err == nil || !strings.Contains(err.Error(), "docs/adr/0004-x.md has no status line") {
			t.Errorf("%s: err = %v; want %q", name, err, "docs/adr/0004-x.md has no status line")
		}
	}
}

func TestTheOutputOfTheOptionsSessionIsCheckedBeforeItIsKept(t *testing.T) {
	got, err := splitOptions(out("db", dbTable, "cache", cacheTable))
	if err != nil || len(got) != 2 || string(got["db"]) != dbTable || string(got["cache"]) != cacheTable {
		t.Fatalf("splitOptions(two decisions) = %v, %v; want the two tables, byte for byte", got, err)
	}
	// F3: a name that the task prompts use for an input file is refused, each name alone, and db is the control.
	for _, name := range []string{"draft", "spans", "answers", "tests", "tasks"} {
		if _, err := splitOptions(out(name, dbTable)); err == nil || !strings.Contains(err.Error(), `"`+name+`" is reserved`) {
			t.Errorf("the decision %s: err = %v; want it refused as reserved", name, err)
		}
		if _, err := rewriteFiles(map[string][]byte{name: []byte(cacheTable)}, out(name, dbTable)); err == nil {
			t.Errorf("the rewrite of the decision %s: err = nil; want a refusal", name)
		}
	}
	if _, err := splitOptions(out("db", dbTable)); err != nil {
		t.Errorf("the decision db: %v; want kept", err)
	}
	one := optHead + "pg\tInv-1\tpass\tdocs/a.md\n"
	for _, c := range []struct{ name, data, want string }{
		{"one option", string(out("db", one)), "decision db has fewer than two options"},
		{"a bad name", string(out("DB", dbTable)), `"DB" is not a word of lowercase letters, digits and -`},
		{"no decision", "decision\toption\tconstraint\tresult\tevidence\n", "holds no decision"},
		{"another header", "option\tconstraint\tresult\tevidence\n", "the header"},
		{"a table its reader refuses", strings.Replace(string(out("db", dbTable)), "pass\tdocs/a.md", "maybe\tdocs/a.md", 1), `db: line 2, column "result"`},
		{"a row of four cells", "decision\toption\tconstraint\tresult\tevidence\ndb\tpg\tInv-1\tpass\n", "line 2"},
	} {
		if _, err := splitOptions([]byte(c.data)); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v; want %q", c.name, err, c.want)
		}
	}
}

func TestTheRoundOfADecisionIsReadFromTheRecords(t *testing.T) {
	current, round1 := optionTables(map[string][]byte{"db.tsv": []byte(dbTable), "cache.tsv": []byte(cacheTable), "round-1/db.tsv": []byte("old"), "notes.txt": []byte("x"), "deep/x/y.tsv": []byte("x")})
	if len(current) != 2 || string(current["db"]) != dbTable || string(current["cache"]) != cacheTable || len(round1) != 1 || string(round1["db"]) != "old" {
		t.Errorf("optionTables = %v, %v; want db and cache as current and db as round 1", current, round1)
	}
}

func TestEveryOptionDroppedAsksOneRewriteAndASecondDropStops(t *testing.T) {
	ids := []string{"Inv-1", "NFR-001"}
	res, err := filterAll(map[string][]byte{"db": []byte(dbTable), "cache": []byte(cacheTable)}, ids)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := res["db"].Chosen(); !ok || got != "pg" || len(res["cache"].Survivors) != 0 || len(res["cache"].Options) != 2 {
		t.Fatalf("filterAll = %+v; want pg the survivor of db and none of cache", res)
	}
	rewrite, stop := noSurvivors(nil, res)
	if !slices.Equal(rewrite, []string{"cache"}) || len(stop) != 0 {
		t.Errorf("round 1, cache dropped: rewrite %v, stop %v; want the rewrite of cache and no stop", rewrite, stop)
	}
	rewrite, stop = noSurvivors(map[string][]byte{"cache": []byte(cacheTable)}, res)
	if len(rewrite) != 0 || !slices.Equal(stop, []string{"cache"}) {
		t.Errorf("the rewrite dropped too: rewrite %v, stop %v; want no second rewrite and the stop", rewrite, stop)
	}
	if rewrite, stop = noSurvivors(map[string][]byte{"db": []byte(dbTable)}, map[string]spec.Result{"db": res["db"]}); len(rewrite)+len(stop) != 0 {
		t.Errorf("a decision with a survivor: rewrite %v, stop %v; want neither", rewrite, stop)
	}
	if _, err := filterAll(map[string][]byte{"db": []byte("option\tconstraint\n")}, ids); err == nil || !strings.Contains(err.Error(), "spec/options/db.tsv") {
		t.Errorf("a table its reader refuses: err = %v; want the path in it", err)
	}
}

func TestTheCommentWithNoSurvivorNamesEachDroppedOptionWithItsRows(t *testing.T) {
	ids := []string{"Inv-1", "NFR-001"}
	r1, _ := filterAll(map[string][]byte{"cache": []byte(cacheTable)}, ids)
	r2, _ := filterAll(map[string][]byte{"cache": []byte(optHead + "redis\tInv-1\tfail\t—\nredis\tNFR-001\tpass\tdocs/c.md\nlocal\tInv-1\tpass\t—\nlocal\tNFR-001\tpass\tdocs/d.md\n")}, ids)
	want := `LAYUP Shape: no survivor

## cache

Round 1:
- redis: Inv-1: fail
- local: NFR-001: missing

Round 2, the rewrite:
- redis: Inv-1: fail
- local: Inv-1: pass with no evidence
`
	if got := noSurvivorComment([]string{"cache"}, r1, r2); got != want {
		t.Errorf("noSurvivorComment =\n%s\nwant\n%s", got, want)
	}
}

func TestTheRewriteMustHoldExactlyTheDecisionsItReplaces(t *testing.T) {
	old := map[string][]byte{"cache": []byte(cacheTable)}
	files, err := rewriteFiles(old, out("cache", dbTable))
	if err != nil || len(files) != 2 || string(files["spec/options/cache.tsv"]) != dbTable || string(files["spec/options/round-1/cache.tsv"]) != cacheTable {
		t.Errorf("rewriteFiles = %v, %v; want the new table and the table of round 1", files, err)
	}
	for name, data := range map[string][]byte{"a decision added": out("cache", dbTable, "db", dbTable), "a decision omitted": out("db", dbTable), "one option": out("cache", optHead+"a\tInv-1\tpass\tdocs/a.md\n")} {
		if _, err := rewriteFiles(old, data); err == nil {
			t.Errorf("%s: err = nil; want a refusal", name)
		}
	}
}

func TestThePromptsListTheConstraintsAndTheFiles(t *testing.T) {
	cs := []spec.Constraint{{ID: "NFR-002", Text: "p99 under 200 ms"}, {ID: "Inv-1"}}
	in := map[string]string{"draft.tsv": "/s/input/draft.tsv", "spans.tsv": "/s/input/spans.tsv", "problem-statement.md": "/s/input/ps.md", "answers.tsv": "/s/input/answers.tsv", "dropped-cache.tsv": "/s/input/dropped-cache.tsv", "reasons.txt": "/s/input/reasons.txt"}
	first, second := string(optionsPrompt(cs)(in, "/s/result/options.tsv")), string(rewritePrompt(cs)(in, "/s/result/options.tsv"))
	for _, want := range []string{"The output file: /s/result/options.tsv\n", "The draft: /s/input/draft.tsv\n", "- NFR-002: p99 under 200 ms\n", "- Inv-1\n", "decision<TAB>option<TAB>constraint<TAB>result<TAB>evidence", "at least two options", "none of the names draft, spans, answers, tests and tasks"} {
		if !strings.Contains(first, want) {
			t.Errorf("the prompt of the options session lacks %q:\n%s", want, first)
		}
	}
	for _, want := range []string{"The output file: /s/result/options.tsv\n", "The reasons: /s/input/reasons.txt\n", "- NFR-002: p99 under 200 ms\n", "cache: /s/input/dropped-cache.tsv\n", "decision<TAB>option", "none of the names draft, spans, answers, tests and tasks"} {
		if !strings.Contains(second, want) {
			t.Errorf("the prompt of the rewrite session lacks %q:\n%s", want, second)
		}
	}
}

// A fail has the empty value in evidence, which the TSV rules write as — (U+2014). The prompts name that
// character, in quotes, and not in words: two sessions of the demo of #370 wrote "-" (task T-c3zv, #473).
func TestEachOptionsPromptNamesTheEmptyEvidenceOfAFail(t *testing.T) {
	cs := []spec.Constraint{{ID: "Inv-1"}}
	in := map[string]string{"draft.tsv": "/s/d", "spans.tsv": "/s/s", "problem-statement.md": "/s/p", "answers.tsv": "/s/a", "reasons.txt": "/s/r"}
	prompts := map[string]string{"options": string(optionsPrompt(cs)(in, "/s/o")), "rewrite": string(rewritePrompt(cs)(in, "/s/o"))}
	const lead, tail = `the single character "`, `" (U+2014)`
	for name, p := range prompts {
		if strings.Contains(p, "a dash") {
			t.Errorf("the %s prompt names the empty evidence in words (a dash)", name)
		}
		if !strings.Contains(p, "a hyphen - or any other text is refused") {
			t.Errorf("the %s prompt does not say that a hyphen or other text is refused:\n%s", name, p)
		}
		i := strings.Index(p, lead)
		if i < 0 {
			t.Errorf("the %s prompt lacks %q:\n%s", name, lead, p)
			continue
		}
		rest := p[i+len(lead):]
		j := strings.Index(rest, tail)
		if j < 0 {
			t.Errorf("the %s prompt lacks the closing %q after the character", name, tail)
			continue
		}
		token := rest[:j] // the bytes between the quotes, as the prompt gives them
		if _, err := spec.ReadOptions([]byte(optHead + "pg\tInv-1\tfail\t" + token + "\n")); err != nil {
			t.Errorf("the %s prompt quotes %q for the evidence of a fail, which the reader refuses: %v", name, token, err)
		}
	}
	if _, err := spec.ReadOptions([]byte(optHead + "pg\tInv-1\tfail\t-\n")); err == nil {
		t.Error("control: a fail with - in evidence is read; want a refusal")
	}
}
