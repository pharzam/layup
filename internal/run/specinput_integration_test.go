//go:build integration

package run

// The input of the full spec check on a real tree and real records (task
// T-cav3, #412, row 66p): the tree half through git, and the demo on the
// records of a target that went through Intake and the Scaffold.

import (
	"context"
	"io"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/verify"
)

// cavEmptyTree is the base of a head whose every path is a change: the empty
// tree, which git knows in any repository.
const cavEmptyTree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

const cavPRD = "# PRD\n\n## 6. Functional requirements\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| --- | --- | --- | --- | --- |\n| REQ-001 | Alpha. | Must | 1 | F-0007#1 |\n"

func cavTreeInput(t *testing.T, dir, head string) (spec.FullInput, error) {
	t.Helper()
	s := cavSessions(cavFiles(t))
	s.Clone = dir
	return s.specInput(context.Background(), cavEmptyTree, head, nil)
}

func TestTheTreeFieldsOfTheFullInputAreTheRegularFilesOfTheHead(t *testing.T) {
	t.Parallel()
	gitlink := strings.Repeat("1", 40)
	dir, head := treeOf(t,
		entry{"100644", "docs/spec/a.md", treeSection}, entry{"100644", "docs/spec/sub/b.md", "x\n"}, entry{"120000", "docs/spec/l.md", "a.md"},
		entry{"100644", "docs/prd/PRD-0001-x.md", cavPRD}, entry{"100644", "docs/prd/README.md", "r\n"}, entry{"100644", "docs/prd/PRD-0002-\nextra.md", "nl\n"}, entry{"100644", "docs/prd/sub/PRD-0002-y.md", "y\n"},
		entry{"100755", "docs/prd/PRD-0003-z.md", "z\n"}, entry{"160000", "docs/prd/PRD-0004-g.md", gitlink}, entry{"100644", "docs/prd/PRD-0005-x.txt", "t\n"},
		entry{"100644", "docs/facts/F-0003-setup-answers.md", "1. **x:** a\n"}, entry{"100644", "docs/facts/problem-statement-brief.md", "b\n"})
	in, err := cavTreeInput(t, dir, head)
	if err != nil {
		t.Fatal(err)
	}
	if got := slices.Sorted(maps.Keys(in.SpecFiles)); !reflect.DeepEqual(got, []string{"a.md"}) || string(in.SpecFiles["a.md"]) != treeSection {
		t.Errorf("SpecFiles %v; want only a.md with its bytes", got)
	}
	if got := slices.Sorted(maps.Keys(in.PRDs)); !reflect.DeepEqual(got, []string{"PRD-0001-x.md", "PRD-0002-\nextra.md"}) || string(in.PRDs["PRD-0001-x.md"]) != cavPRD || string(in.PRDs["PRD-0002-\nextra.md"]) != "nl\n" {
		t.Errorf("PRDs %v; want PRD-0001-x.md and the PRD whose name holds a line feed, with their bytes (not README.md, a nested, executable, linked or .txt file)", got)
	}
	if want := strings.Split(strings.TrimRight(gitOut(t, dir, "ls-tree", "-r", "-z", "--name-only", head), "\x00"), "\x00"); !reflect.DeepEqual(in.Changed, want) {
		t.Errorf("Changed %v; want the paths of the diff %v", in.Changed, want)
	}
	if !reflect.DeepEqual(in.PSBRecords, []string{"F-0004"}) || !reflect.DeepEqual(in.Records, []string{"F-0003"}) {
		t.Errorf("PSBRecords %v, Records %v; want the predicted ID F-0004 (no numbered facts file) and the setup record F-0003", in.PSBRecords, in.Records)
	}
	if string(in.PSB) != string(cavFiles(t)["start/problem-statement.md"]) {
		t.Errorf("PSB %d bytes; want the problem statement of the records", len(in.PSB))
	}
}

func TestTheNumberedFactsFileAtTheHeadIsTheProblemStatementRecord(t *testing.T) {
	t.Parallel()
	gitlink := strings.Repeat("1", 40)
	named := func(e ...entry) []entry { return append([]entry{{"100644", "docs/prd/PRD-0001-x.md", cavPRD}}, e...) }
	for _, c := range []struct {
		name    string
		entries []entry
		want    []string
		fail    bool
	}{
		{"none", named(), []string{"F-0001"}, false},
		{"one regular file", named(entry{"100644", "docs/facts/F-0007-numbered-facts.md", "1. **need:** one\n"}), []string{"F-0007"}, false},
		{"a link, an executable file, a gitlink, a nested file", named(entry{"120000", "docs/facts/F-0007-numbered-facts.md", "x"}, entry{"100755", "docs/facts/F-0008-numbered-facts.md", "1. x\n"},
			entry{"160000", "docs/facts/F-0009-numbered-facts.md", gitlink}, entry{"100644", "docs/facts/sub/F-0010-numbered-facts.md", "1. x\n"}), []string{"F-0001"}, false},
		{"two files", named(entry{"100644", "docs/facts/F-0007-numbered-facts.md", "1. x\n"}, entry{"100644", "docs/facts/F-0009-numbered-facts.md", "1. x\n"}), nil, true},
	} {
		dir, head := treeOf(t, c.entries...)
		in, err := cavTreeInput(t, dir, head)
		if c.fail {
			wantSource(t, c.name, err, "PSBRecords", "target:docs/facts@"+head)
			if !reflect.DeepEqual(in, spec.FullInput{}) {
				t.Errorf("%s: the input on an error is not the zero input", c.name)
			}
			continue
		}
		if err != nil || !reflect.DeepEqual(in.PSBRecords, c.want) {
			t.Errorf("%s: PSBRecords %v, %v; want %v", c.name, in.PSBRecords, err, c.want)
			continue
		}
		// The consumer: the fact 1 of F-0007 is the span fact 1, so the citation covers.
		rows := spec.CheckFull(in)
		if got := rows[0][1]; (c.want[0] == "F-0007") != (got == "pass") {
			t.Errorf("%s: covers is %v; want pass exactly with the record F-0007 (the citation F-0007#1)", c.name, rows[0])
		}
	}
}

func TestAFailedReadOfTheTreeIsASourceErrorOfTheTarget(t *testing.T) {
	t.Parallel()
	// A head the clone lacks fails at the diff.
	dir, _ := treeOf(t, entry{"100644", "docs/spec/a.md", treeSection})
	in, err := cavTreeInput(t, dir, strings.Repeat("a", 40))
	wantSource(t, "a head the clone lacks", err, "Changed", "target:diff@"+strings.Repeat("a", 40))
	if !reflect.DeepEqual(in, spec.FullInput{}) {
		t.Errorf("the input on an error is %+v; want the zero input", in)
	}
	// A blob that cannot be read, after every earlier read succeeded: only the Show fails.
	all := []entry{{"100644", "docs/spec/a.md", "spec\n"}, {"100644", "docs/prd/PRD-0001-x.md", cavPRD}, {"100644", "docs/facts/F-0007-numbered-facts.md", "1. x\n"}}
	for _, c := range []struct{ field, path string }{
		{"SpecFiles", "docs/spec/a.md"}, {"PRDs", "docs/prd/PRD-0001-x.md"}, {"PSBRecords", "docs/facts/F-0007-numbered-facts.md"},
	} {
		dir, head := treeOf(t, all...)
		blob := strings.TrimSpace(gitOut(t, dir, "rev-parse", head+":"+c.path))
		if err := os.Remove(filepath.Join(dir, ".git", "objects", blob[:2], blob[2:])); err != nil {
			t.Fatalf("%s: the blob is not a loose object: %v", c.path, err)
		}
		in, err := cavTreeInput(t, dir, head)
		wantSource(t, c.path, err, c.field, "target:"+map[string]string{"SpecFiles": "docs/spec", "PRDs": "docs/prd", "PSBRecords": "docs/facts"}[c.field]+"@"+head)
		if !reflect.DeepEqual(in, spec.FullInput{}) {
			t.Errorf("%s: the input on an error is not the zero input", c.path)
		}
	}
}

// The demo of row 66p: on the records of a target that went through Intake and
// the Scaffold, and a head that descends from the scaffolded main, each field
// of the FullInput is its record or its file, and CheckFull judges the head.
func TestTheFullInputOfAHeadAfterIntakeAndTheScaffold(t *testing.T) {
	t.Parallel()
	saved := verifyRun
	t.Cleanup(func() { verifyRun = saved })
	w := scaffoldWorld(t, proseFiles())
	verifyRun = func(string, func(int, int, string) func(), io.Writer) (verify.Table, error) {
		return verify.Table{Rows: []verify.Row{{Check: "pin", Result: "pass", Reason: "—"}}}, nil
	}
	toSetup(t, w, "go")
	area := filepath.Join(w.dir, "work", "acme", "target")
	gitOut(t, filepath.Join(area, "target"), "push", "-q", w.bare, "layup-setup:main")
	applies := filepath.Join(w.dir, "applies", "acme", "target")
	writeFile(t, filepath.Join(applies, "default.json"), "{\"source\":\"acme/target\",\"bypass_actors\":[]}\n")
	writeFile(t, filepath.Join(applies, "records.json"), "{\"source\":\"acme/target\",\"bypass_actors\":[{\"actor_id\":42}]}")
	rows := Restart(context.Background(), w.config("0.1.0-dev"))
	if r := rowOf(rows, "scaffold-push"); r.Result != "done" || !strings.Contains(r.Detail, "the apply of the rulesets recorded") {
		t.Fatalf("the Scaffold is not done: %+v\n%s", r, steps(rows))
	}
	// The register and the events of the tasks, in the records.
	files := cavFiles(t)
	pushed := map[string]string{}
	for _, p := range []string{"tasks.tsv", "tasks/T-aaaa/events.tsv", "tasks/T-bbbb/events.tsv"} {
		pushed[p] = string(files[p])
	}
	pushRecords(t, w.bare, pushed, "the register and the events of the tasks")

	clone := filepath.Join(t.TempDir(), "clone")
	gitOut(t, "/", "clone", "-q", "-b", "main", "file://"+w.bare, clone)
	sessions := func() *Sessions {
		gitOut(t, clone, "fetch", "-q", "origin")
		store := &recordsStore{dir: clone, url: "file://" + w.bare, read: strings.TrimSpace(gitOut(t, clone, "rev-parse", "origin/layup-records")),
			auth: func(context.Context) (git.Auth, error) { return git.Auth{}, nil },
			who: func() git.Identity {
				return git.Identity{Name: "layup-agent[bot]", Email: "9001+layup-agent[bot]@users.noreply.github.com", Time: time.Now()}
			}}
		return &Sessions{Store: store, RunID: "0123456789abcdef", Host: t.TempDir(), Target: "acme/target", Clone: clone, Branch: "main", Now: time.Now}
	}
	base := strings.TrimSpace(gitOut(t, clone, "rev-parse", "main"))
	recordOf := func(path string) []byte { return []byte(recordsFile(t, w, path)) }

	spans, err := records.ReadSpans(recordOf("spec/spans.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	fact := ""
	for _, r := range spans {
		if r[2] != "" {
			fact = r[2]
			break
		}
	}
	if fact == "" {
		t.Fatal("the records of Intake hold no span fact")
	}
	prd := strings.Replace(cavPRD, "F-0007#1", "F-0007#"+fact, 1)
	if err := os.MkdirAll(filepath.Join(clone, "docs", "spec"), 0o755); err != nil {
		t.Fatal(err)
	}
	for p, d := range map[string]string{"docs/spec/a.md": treeSection, "docs/prd/PRD-0001-x.md": prd, "docs/facts/F-0007-numbered-facts.md": "1. **need:** one\n"} {
		writeFile(t, filepath.Join(clone, filepath.FromSlash(p)), d)
	}
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "the specification head")
	head := strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))

	in, err := sessions().specInput(context.Background(), base, head, nil)
	if err != nil {
		t.Fatal(err)
	}
	questions, _ := records.ReadQuestions(recordOf("questions.tsv"))
	answers, _ := records.ReadAnswers(recordOf("intake/answers.tsv"))
	if want := []string{"docs/facts/F-0007-numbered-facts.md", "docs/prd/PRD-0001-x.md", "docs/spec/a.md"}; !reflect.DeepEqual(in.Changed, want) {
		t.Errorf("Changed %v; want %v", in.Changed, want)
	}
	if string(in.SpecFiles["a.md"]) != treeSection || string(in.PRDs["PRD-0001-x.md"]) != prd || !reflect.DeepEqual(in.PSBRecords, []string{"F-0007"}) || in.AnswersRecord != "" {
		t.Errorf("SpecFiles %d, PRDs %d, PSBRecords %v, AnswersRecord %q; want a.md, PRD-0001-x.md, [F-0007] and the empty answers record", len(in.SpecFiles), len(in.PRDs), in.PSBRecords, in.AnswersRecord)
	}
	if string(in.PSB) != string(recordOf("start/problem-statement.md")) || len(in.PSB) == 0 || !reflect.DeepEqual(in.Spans, spans) ||
		!reflect.DeepEqual(in.Questions, questions) || len(in.Questions) == 0 || !reflect.DeepEqual(in.Answers, answers) || len(in.Answers) == 0 {
		t.Errorf("PSB %d bytes, Spans %d, Questions %d, Answers %d; want each the record of the run", len(in.PSB), len(in.Spans), len(in.Questions), len(in.Answers))
	}
	if string(in.Register) != string(recordOf("tasks.tsv")) || !reflect.DeepEqual(in.Merged, []string{"T-aaaa"}) || in.Index != nil || in.Versions != nil || in.Copy != nil {
		t.Errorf("Register %d bytes, Merged %v, Index %v, Copy %q; want the register, [T-aaaa], none and none", len(in.Register), in.Merged, in.Index, in.Copy)
	}
	table := spec.CheckFull(in)
	if table[0][0] != "covers" || table[0][1] != "pass" || table[0][2] == "no specification yet" {
		t.Errorf("covers on the head: %v; want pass (the PRD cites F-0007#%s, the span fact)", table[0], fact)
	}

	// A head that changes no specification gets the clear table (the control is the table above).
	same, err := sessions().specInput(context.Background(), base, base, nil)
	if err != nil || len(same.Changed) != 0 {
		t.Fatalf("base..base: Changed %v, %v", same.Changed, err)
	}
	for _, r := range spec.CheckFull(same) {
		if r[1] != "clear" || r[2] != "no specification yet" {
			t.Errorf("the table of a head that changes no specification: %v; want every rule clear, no specification yet", r)
		}
	}

	// After a bet: the index, the exact versions, and the copy the caller passes.
	post := cavFiles(t)
	cavInventory(t, post, 1, "")
	pushRecords(t, w.bare, map[string]string{"spec/inventory.tsv": string(post["spec/inventory.tsv"]), "spec/inventory-1.tsv": string(post["spec/inventory-1.tsv"])}, "the first bet")
	s := sessions()
	for _, c := range []struct {
		name string
		copy []byte
		want string
	}{{"no copy", nil, "fail"}, {"the copy of the bet", []byte("REQ-001: Must, phase 1\n"), "pass"}} {
		in, err := s.specInput(context.Background(), base, head, c.copy)
		if err != nil || !reflect.DeepEqual(in.Copy, c.copy) || len(in.Index) != 1 || string(in.Versions[1]) != string(post["spec/inventory-1.tsv"]) {
			t.Fatalf("%s: Copy %q, Index %v, Versions %d, %v; want the caller's copy, the index and the exact version", c.name, in.Copy, in.Index, len(in.Versions), err)
		}
		table := spec.CheckFull(in)
		if got := table[4]; got[0] != "classes" || got[1] != c.want || (c.want == "fail" && got[2] != "did not run: the copy of the latest bet") {
			t.Errorf("classes with %s: %v; want %s", c.name, got, c.want)
		}
		if got := table[5]; got[0] != "inventory" || got[1] != "pass" {
			t.Errorf("inventory with %s: %v; want pass (the version is its exact bytes)", c.name, got)
		}
	}
}

// The specimen of the specification instruction, on a real tree and the records of a
// target, through the production input (specRecords and specTree: no PSBRecords and no
// AnswersRecord of the test's own): the full check passes every rule, and the pinned
// prd-lint accepts the head as a specification that waits for its bet.
func TestTheSpecimenPassesTheFullCheckAndThePrdLintThroughTheProductionInput(t *testing.T) {
	t.Parallel()
	prd := specimen(t, -1)
	dir, head := treeOf(t, entry{"100644", "docs/prd/PRD-0001-x.md", prd}, entry{"100644", "docs/spec/a.md", "## REQ-001 — Alpha\n\nThe text.\n"},
		entry{"100644", "docs/facts/F-0003-setup-answers.md", "1. **x:** a\n"})
	in, err := cavTreeInput(t, dir, head)
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range spec.CheckFull(in) {
		if row[1] == "fail" {
			t.Errorf("the rule %s fails on the specimen through the production input: %s", row[0], row[2])
		}
	}
	root := t.TempDir()
	for p, text := range map[string]string{"docs/prd/PRD-0001-x.md": prd, "docs/facts/F-0003-setup-answers.md": "1. **x:** a\n"} {
		os.MkdirAll(filepath.Dir(filepath.Join(root, p)), 0o755)
		os.WriteFile(filepath.Join(root, p), []byte(text), 0o644)
	}
	script, _ := filepath.Abs("../../docs/prd/prd-lint.sh")
	if ok, lines, err := spec.WaitingSpecAccepted(script, filepath.Join(root, "docs", "prd")); err != nil || !ok {
		t.Errorf("prd-lint: accepted %v, %v, %v", ok, lines, err)
	}
}

// loop.md, The instruction of a task step: the prompt of the specification and the input of the full
// check predict the facts records of one tree by one contract (the regular files of docs/facts/, as
// the numbered-facts reader takes them). The ID that the prompt tells the specification to cite
// is the one the input holds, the record that it names is among the input's records, and the
// citations of that prompt pass the rule covers; an executable file or a link is no record for either.
func TestThePromptAndTheFullInputPredictTheSameFactsRecords(t *testing.T) {
	t.Parallel()
	gitlink := strings.Repeat("1", 40)
	setup := entry{"100644", "docs/facts/F-0003-setup-answers.md", "1. **x:** a\n"}
	for _, c := range []struct {
		name         string
		extra        []entry
		record, next string
	}{
		{"regular files only (the control)", []entry{{"100644", "docs/facts/F-0008-extra.md", "x\n"}}, "F-0008", "F-0009"},
		{"an executable file", []entry{{"100755", "docs/facts/F-0008-extra.md", "x\n"}}, "F-0003", "F-0004"},
		{"a link", []entry{{"120000", "docs/facts/F-0008-extra.md", "F-0003-setup-answers.md"}}, "F-0003", "F-0004"},
		{"a gitlink", []entry{{"160000", "docs/facts/F-0008-extra.md", gitlink}}, "F-0003", "F-0004"},
	} {
		t.Run(c.name, func(t *testing.T) {
			entries := append([]entry{{"100644", "docs/prd/PRD-0001-x.md", cavPRD}, setup}, c.extra...)
			dir, head := treeOf(t, entries...)
			tree, err := targetTreeAt(dir, head)
			if err != nil {
				t.Fatal(err)
			}
			record, next, err := factsRecords(tree.Facts)
			if err != nil {
				t.Fatal(err)
			}
			if "F-"+record != c.record || "F-"+next != c.next {
				t.Errorf("the prompt: record F-%s, next F-%s; want %s and %s", record, next, c.record, c.next)
			}
			in, err := cavTreeInput(t, dir, head)
			if err != nil {
				t.Fatal(err)
			}
			if !slices.Equal(in.PSBRecords, []string{"F-" + next}) || !slices.Contains(in.Records, "F-"+record) || in.Records[len(in.Records)-1] != "F-"+record {
				t.Errorf("the full input: PSBRecords %v, Records %v; want the prompt's next F-%s and the record F-%s as the highest", in.PSBRecords, in.Records, next, record)
			}
			prd := strings.Replace(cavPRD, "F-0007#1", "F-"+next+"#1, F-"+record, 1)
			in.PRDs = map[string][]byte{"PRD-0001-x.md": []byte(prd)}
			for _, row := range spec.CheckFull(in) {
				if row[0] == "covers" && row[1] == "fail" {
					t.Errorf("the rule covers fails on the citations of the prompt: %s", row[2])
				}
			}
		})
	}
}
