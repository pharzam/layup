package run

// The input of the full spec check (docs/spec/loop.md, REQ-012, The full spec
// check; task T-cav3, #412, row 66p): each field of spec.FullInput is its
// record, on the in-memory records of a unit; the tree half is the integration
// level's, because the unit level touches no disk.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/spec"
	"github.com/pharzam/layup/internal/tsv"
	workrec "github.com/pharzam/layup/internal/work"
)

const (
	cavBase = "1111111111111111111111111111111111111111"
	cavHead = "2222222222222222222222222222222222222222"
	cavAt   = "2026-10-09T12:00:00Z"
)

func cavTable(t *testing.T, s tsv.Schema, rows [][]string) []byte {
	t.Helper()
	data, err := table(s, rows)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// cavFiles is the records of a target after Intake and before the first bet:
// three tasks (merged, only a session event, no events file), and the setup
// record whose two rows name a file and no numbered fact.
func cavFiles(t *testing.T) map[string][]byte {
	t.Helper()
	tasks := [][]string{
		{"T-aaaa", "1", "build", "REQ-012", "small", "", "1", cavAt},
		{"T-bbbb", "2", "build", "REQ-012", "small", "", "", cavAt},
		{"T-cccc", "3", "build", "REQ-012", "small", "", "", cavAt},
	}
	questions := fuQuestions()
	for _, q := range questions {
		if q[12] == "closed" {
			q[11] = cavAt
		}
	}
	return map[string][]byte{
		"start/problem-statement.md": []byte(strings.Repeat("x", 100)),
		"spec/spans.tsv": cavTable(t, records.SpansSchema, [][]string{
			{"0", "10", "1", "need", "", "s1"}, {"10", "20", "2", "measure", "", "s2"}, {"20", "30", "", "", "no fact", "s3"}}),
		"questions.tsv":      cavTable(t, records.QuestionsSchema, questions),
		"intake/answers.tsv": cavTable(t, records.AnswersSchema, fuAnswers("Q-001", "go", "Q-002", "ten")),
		"tasks.tsv":          cavTable(t, TasksSchema, tasks),
		"tasks/T-aaaa/events.tsv": cavTable(t, records.EventsSchema, [][]string{
			{"1", "session", "1", "S-0123abcd", "", "", "", cavAt}, {"2", "merged", "1", "", "", cavBase, "7", cavAt}}),
		"tasks/T-bbbb/events.tsv": cavTable(t, records.EventsSchema, [][]string{
			{"1", "session", "1", "S-0123abcd", "", "", "", cavAt}}),
		"setup/record.tsv": cavTable(t, workrec.RecordSchema, [][]string{
			{"S04", "answers.record", "docs/facts/F-0001-setup-answers.md", "computed", "the next free ID of docs/facts/"},
			{"S06", "brief.copy", "docs/facts/problem-statement-brief.md", "computed", "the raw file name of a brief"}}),
	}
}

// cavInventory adds an index of n versions to files, the last with the given
// bytes' hash replaced by hash when it is not empty.
func cavInventory(t *testing.T, files map[string][]byte, n int, hash string) {
	t.Helper()
	var index [][]string
	for i := 1; i <= n; i++ {
		v := cavTable(t, spec.InventorySchema, [][]string{{"1", "need", "in", "9"}, {"2", "measure", "in", "9"}, {"Q-001", "context", "in", "9"}})
		sum := fmt.Sprintf("%x", sha256.Sum256(v))
		if i == n && hash != "" {
			sum = hash
		}
		files[fmt.Sprintf("spec/inventory-%d.tsv", i)] = v
		index = append(index, []string{fmt.Sprint(i), sum, "M2e"})
	}
	files["spec/inventory.tsv"] = cavTable(t, spec.InventoryIndexSchema, index)
}

func cavSessions(files map[string][]byte) *Sessions {
	return newSessions(&standinStore{files: files})
}

func cavRecords(t *testing.T, s *Sessions) (spec.FullInput, error) {
	t.Helper()
	var in spec.FullInput
	err := s.specRecords(context.Background(), &in)
	return in, err
}

func wantSource(t *testing.T, name string, err error, field, source string) {
	t.Helper()
	var se *SourceError
	if !errors.As(err, &se) || se.Source != source || (field != "" && se.Field != field) {
		t.Errorf("%s: %v; want a SourceError with the source %q", name, err, source)
	}
}

func TestEachRecordFieldOfTheFullInputIsItsRecord(t *testing.T) {
	files := cavFiles(t)
	in, err := cavRecords(t, cavSessions(files))
	if err != nil {
		t.Fatal(err)
	}
	spans, _ := records.ReadSpans(files["spec/spans.tsv"])
	questions, _ := records.ReadQuestions(files["questions.tsv"])
	answers, _ := records.ReadAnswers(files["intake/answers.tsv"])
	if !bytes.Equal(in.PSB, files["start/problem-statement.md"]) || !reflect.DeepEqual(in.Spans, spans) || len(in.Spans) != 3 ||
		!reflect.DeepEqual(in.Questions, questions) || len(in.Questions) != 4 || !reflect.DeepEqual(in.Answers, answers) || len(in.Answers) != 2 ||
		!bytes.Equal(in.Register, files["tasks.tsv"]) || !reflect.DeepEqual(in.Merged, []string{"T-aaaa"}) {
		t.Errorf("the input: PSB %d bytes, Spans %v, Questions %d, Answers %v, Register %d bytes, Merged %v; want each its record",
			len(in.PSB), in.Spans, len(in.Questions), in.Answers, len(in.Register), in.Merged)
	}
	if in.Copy != nil || in.PSBRecords != nil || in.AnswersRecord != "" || in.Index != nil || in.Versions != nil {
		t.Errorf("the fields that the records do not hold: %+v; want the zero value", in)
	}
}

func TestAMissingRecordOfTheFullInputIsAnError(t *testing.T) {
	for _, path := range []string{"start/problem-statement.md", "spec/spans.tsv", "questions.tsv", "intake/answers.tsv", "tasks.tsv"} {
		files := cavFiles(t)
		delete(files, path)
		in, err := cavSessions(files).specInput(context.Background(), cavBase, cavHead, nil)
		wantSource(t, path, err, "", "records:"+path)
		if !reflect.DeepEqual(in, spec.FullInput{}) {
			t.Errorf("%s: the input on an error is %+v; want the zero input", path, in)
		}
	}
	if _, err := cavRecords(t, cavSessions(cavFiles(t))); err != nil {
		t.Errorf("the control, every record there: %v", err)
	}
}

func TestAFailedReadOfARecordIsAnErrorEvenWhereItsAbsenceIsNot(t *testing.T) {
	files := cavFiles(t)
	cavInventory(t, files, 2, "")
	for _, path := range []string{"start/problem-statement.md", "spec/spans.tsv", "tasks.tsv", "spec/inventory.tsv", "spec/inventory-2.tsv", "tasks/T-bbbb/events.tsv"} {
		s := cavSessions(files)
		s.Store = failRead{s.Store, path}
		in, err := s.specInput(context.Background(), cavBase, cavHead, nil)
		wantSource(t, path, err, "", "records:"+path)
		if !reflect.DeepEqual(in, spec.FullInput{}) {
			t.Errorf("%s: the input on an error is %+v; want the zero input", path, in)
		}
	}
	if _, err := cavRecords(t, cavSessions(files)); err != nil {
		t.Errorf("the control, every read there: %v", err)
	}
}

func TestARecordThatItsReaderRefusesIsAnError(t *testing.T) {
	files := cavFiles(t)
	cavInventory(t, files, 1, "")
	for _, path := range []string{"spec/spans.tsv", "questions.tsv", "intake/answers.tsv", "tasks.tsv", "spec/inventory.tsv", "tasks/T-aaaa/events.tsv"} {
		broken := map[string][]byte{}
		for p, d := range files {
			broken[p] = d
		}
		broken[path] = []byte(strings.Replace(string(files[path]), "\t", " ", 1))
		_, err := cavRecords(t, cavSessions(broken))
		wantSource(t, path, err, "", "records:"+path)
	}
	if _, err := cavRecords(t, cavSessions(files)); err != nil {
		t.Errorf("the control, no record broken: %v", err)
	}
}

func TestNoInventoryIndexIsBeforeTheFirstBet(t *testing.T) {
	in, err := cavRecords(t, cavSessions(cavFiles(t)))
	if err != nil || in.Index != nil || in.Versions != nil {
		t.Errorf("no spec/inventory.tsv: Index %v, Versions %v, %v; want nil, nil and no error (before the first bet)", in.Index, in.Versions, err)
	}
}

func TestEachVersionOfTheIndexIsItsExactBytes(t *testing.T) {
	files := cavFiles(t)
	cavInventory(t, files, 2, strings.Repeat("0", 64))
	in, err := cavRecords(t, cavSessions(files))
	if err != nil || len(in.Index) != 2 || !bytes.Equal(in.Versions[1], files["spec/inventory-1.tsv"]) || !bytes.Equal(in.Versions[2], files["spec/inventory-2.tsv"]) {
		t.Fatalf("Index %v, Versions %d, %v; want two rows and the exact bytes of both versions", in.Index, len(in.Versions), err)
	}
	// The builder passes a wrong hash through; the rule inventory fails it.
	in.Changed = []string{"docs/prd/PRD-0001-x.md"}
	if rows := spec.CheckFull(in); rows[5][0] != "inventory" || rows[5][1] != "fail" || !strings.Contains(rows[5][2], "version 2") {
		t.Errorf("the rule inventory on a wrong hash: %v; want fail naming version 2", rows[5])
	}
	delete(files, "spec/inventory-2.tsv")
	_, err = cavRecords(t, cavSessions(files))
	wantSource(t, "a version that the index names and the records lack", err, "Versions", "records:spec/inventory-2.tsv")
	files["spec/inventory.tsv"] = []byte("n\tsha256\tbet\n")
	_, err = cavRecords(t, cavSessions(files))
	wantSource(t, "an index with a header and no row", err, "Index", "records:spec/inventory.tsv")
}

func TestMergedIsTheTasksWithAMergedEvent(t *testing.T) {
	files := cavFiles(t)
	in, err := cavRecords(t, cavSessions(files))
	if err != nil || !reflect.DeepEqual(in.Merged, []string{"T-aaaa"}) {
		t.Errorf("Merged %v, %v; want [T-aaaa] (T-bbbb has only a session event, T-cccc no events file)", in.Merged, err)
	}
	files["tasks/T-cccc/events.tsv"] = cavTable(t, records.EventsSchema, [][]string{{"1", "merged", "1", "", "", cavBase, "9", cavAt}})
	if in, _ = cavRecords(t, cavSessions(files)); !reflect.DeepEqual(in.Merged, []string{"T-aaaa", "T-cccc"}) {
		t.Errorf("Merged %v; want [T-aaaa T-cccc] in the order of the register", in.Merged)
	}
}

func TestTheRegisterIsItsBytesAndNotJudged(t *testing.T) {
	files := cavFiles(t)
	files["tasks.tsv"] = cavTable(t, TasksSchema, [][]string{{"T-bbbb", "2", "build", "", "small", "", "", cavAt}})
	if _, err := ReadRegister(files["tasks.tsv"]); err == nil {
		t.Fatal("the control: ReadRegister accepts a task with no requirement")
	}
	in, err := cavRecords(t, cavSessions(files))
	if err != nil || !bytes.Equal(in.Register, files["tasks.tsv"]) {
		t.Fatalf("Register %q, %v; want the bytes and no error", in.Register, err)
	}
	in.Changed = []string{"docs/prd/PRD-0001-x.md"}
	if rows := spec.CheckFull(in); rows[8][0] != "tasks" || rows[8][1] != "fail" || !strings.Contains(rows[8][2], "the task T-bbbb names no requirement") {
		t.Errorf("the rule tasks: %v; want fail, the task T-bbbb names no requirement", rows[8])
	}
}

// The answers record of the setup numbers S01-stack first, so its fact 2 is not
// Q-002: no field of the input names it, and a citation of a bare Q-NNN
// resolves through the answers (decided here, task T-cav3).
func TestAnAnswersCitationIsAQuestionIDNotARecordLine(t *testing.T) {
	prd := func(fact string) map[string][]byte {
		return map[string][]byte{"PRD-0001-x.md": []byte("# PRD\n\n## 6. Functional requirements\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| --- | --- | --- | --- | --- |\n| REQ-001 | Alpha. | Must | 1 | " + fact + " |\n")}
	}
	in, err := cavRecords(t, cavSessions(cavFiles(t)))
	if err != nil {
		t.Fatal(err)
	}
	in.Changed = []string{"docs/prd/PRD-0001-x.md"}
	for _, c := range []struct{ fact, want string }{
		{"Q-002", "pass"},                  // the control: a question ID of the answers
		{"F-0001#2", "fail"},               // the line 2 of the setup record, which is S01-stack's neighbour and no Q-002
		{"F-0001-setup-answers#2", "fail"}, // no record is named by a file
	} {
		in.PRDs = prd(c.fact)
		rows := spec.CheckFull(in)
		if rows[0][0] != "covers" || rows[0][1] != c.want || (c.want == "fail" && !strings.Contains(rows[0][2], "cites "+c.fact+", which is no fact")) {
			t.Errorf("covers for %s: %v; want %s", c.fact, rows[0], c.want)
		}
	}
	in.PSBRecords = []string{"F-0001"}
	in.PRDs = prd("F-0001#2")
	if rows := spec.CheckFull(in); rows[0][1] != "pass" {
		t.Errorf("the control, F-0001 named as the record of the problem statement: %v; want covers pass on its fact 2", rows[0])
	}
}

func TestANonFullSHAIsRefusedBeforeAnyRead(t *testing.T) {
	for _, c := range []struct{ name, base, head, source string }{
		{"base", "main", cavHead, "base"}, {"head", cavBase, "HEAD", "head"}, {"upper case", cavBase, "ABCDEF0123456789ABCDEF0123456789ABCDEF01", "head"},
	} {
		reads := &countRead{SessionStore: &standinStore{files: cavFiles(t)}}
		s := newSessions(&standinStore{files: cavFiles(t)})
		s.Store = reads
		in, err := s.specInput(context.Background(), c.base, c.head, nil)
		wantSource(t, c.name, err, "", c.source)
		if reads.n != 0 || !reflect.DeepEqual(in, spec.FullInput{}) {
			t.Errorf("%s: %d reads, input %+v; want none and the zero input", c.name, reads.n, in)
		}
	}
}

// numberedFacts is the one file that PSBRecords reads: the name, the count.
func TestOnlyTheNumberedFactsFileNamesAProblemStatementRecord(t *testing.T) {
	for name, want := range map[string]bool{
		"F-0007-numbered-facts.md": true, "F-0007-setup-answers.md": false, "F-7-numbered-facts.md": false, "F-00007-numbered-facts.md": false,
		"f-0007-numbered-facts.md": false, "F-0007-numbered-facts.md.bak": false, "xF-0007-numbered-facts.md": false, "problem-statement-brief.md": false,
	} {
		if got := isNumberedFacts(name); got != want {
			t.Errorf("isNumberedFacts(%q) = %v; want %v", name, got, want)
		}
	}
	for _, c := range []struct {
		files   []string
		want    []string
		wantErr bool
	}{
		{nil, nil, false}, {[]string{"F-0007-numbered-facts.md"}, []string{"F-0007"}, false},
		{[]string{"F-0007-numbered-facts.md", "F-0009-numbered-facts.md"}, nil, true},
	} {
		files := map[string][]byte{}
		for _, f := range c.files {
			files[f] = []byte("1. **need:** one\n")
		}
		got, err := psbRecordsOf(files)
		if (err != nil) != c.wantErr || !reflect.DeepEqual(got, c.want) {
			t.Errorf("psbRecordsOf(%v) = %v, %v; want %v, error %v", c.files, got, err, c.want, c.wantErr)
		}
	}
}

// countRead counts the reads of a store.
type countRead struct {
	SessionStore
	n int
}

func (c *countRead) ReadFile(ctx context.Context, p string) ([]byte, error) {
	c.n++
	return c.SessionStore.ReadFile(ctx, p)
}
