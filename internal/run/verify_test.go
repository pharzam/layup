package run

import (
	"bytes"
	"errors"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/verify"
)

// The step verify (docs/spec/scaffold.md, The step verify and the records of the
// setup; the input state of a row that is not pass or clear; row 55 of the plan,
// task T-79zk, #278), with its seams as fakes: the loop of the runner, layup
// setup verify, the write of its table, the files of out/, the tree of the setup
// and the records commit. Each case breaks one rule, and a failure never calls
// the commit.

func goodTable() verify.Table {
	return verify.Table{Rows: []verify.Row{{Check: "pin", Result: "pass", Reason: "—"}, {Check: "link-lint", Result: "pass", Reason: "—"},
		{Check: "gate:layout", Result: "clear", Reason: "pending: fixture not run"}}}
}

var setupFiles = map[string][]byte{"setup/record.tsv": []byte("record\n"), "setup/verify.tsv": []byte("table\n"), "rule-paths.tsv": []byte("paths\n")}

// verifyFake is the world of a job: what each seam gives, and what the job asked of them.
type verifyFake struct {
	loops                                             []Step // the result of each call of the loop, in order
	table                                             verify.Table
	verifyErr, writeErr, filesErr, treeErr, commitErr error
	called                                            []string
	wrote                                             []byte
	commits                                           []map[string][]byte
}

func (f *verifyFake) job() verifyJob {
	return verifyJob{
		loop: func() Step {
			f.called = append(f.called, "loop")
			s := f.loops[0]
			f.loops = f.loops[1:]
			return s
		},
		verify: func() (verify.Table, error) { f.called = append(f.called, "verify"); return f.table, f.verifyErr },
		write:  func(b []byte) error { f.called = append(f.called, "write"); f.wrote = b; return f.writeErr },
		files: func() (map[string][]byte, error) {
			f.called = append(f.called, "files")
			return copyMap(setupFiles), f.filesErr
		},
		tree: func() ([]byte, error) { f.called = append(f.called, "tree"); return []byte("tree\n"), f.treeErr },
		commit: func(files map[string][]byte) error {
			f.called = append(f.called, "commit")
			f.commits = append(f.commits, files)
			return f.commitErr
		},
	}
}

func copyMap(m map[string][]byte) map[string][]byte {
	out := map[string][]byte{}
	for k, v := range m {
		out[k] = v
	}
	return out
}

func okLoop() Step { return Step{"setup", "done", setupDoneDetail} }

func TestTheStepVerifyCommitsTheFivePathsOfTheSetup(t *testing.T) {
	f := &verifyFake{loops: []Step{okLoop(), okLoop()}, table: goodTable()}
	got := f.job().do()
	if got != (Step{"verify", "done", verifyDoneDetail}) {
		t.Fatalf("the step %+v; want done", got)
	}
	if want := []string{"loop", "verify", "write", "loop", "files", "tree", "commit"}; !slices.Equal(f.called, want) {
		t.Errorf("the calls %q; want %q", f.called, want)
	}
	var table bytes.Buffer
	if err := goodTable().Write(&table); err != nil || !bytes.Equal(f.wrote, table.Bytes()) {
		t.Errorf("the table written to out/verify.tsv: %q, %v; want the table of layup setup verify, %q", f.wrote, err, table.String())
	}
	want := copyMap(setupFiles)
	want["setup/tree.tsv"], want["README.md"] = []byte("tree\n"), []byte(scaffoldReadme)
	if len(f.commits) != 1 || !reflect.DeepEqual(f.commits[0], want) {
		t.Errorf("the commits %q; want one commit of exactly the five paths %q", f.commits, want)
	}
}

func TestEachFailureOfTheStepVerifyPublishesNoCommit(t *testing.T) {
	boom := errors.New("boom\nthe second line")
	badTable := func(rows ...verify.Row) verify.Table { return verify.Table{Rows: append(goodTable().Rows, rows...)} }
	for _, c := range []struct {
		name   string
		f      verifyFake
		detail string
		calls  []string
	}{
		{"the loop before the verification does not end done", verifyFake{loops: []Step{{"setup", "fail", "S09: glossary: fail"}}},
			"S09: glossary: fail", []string{"loop"}},
		{"layup setup verify gives an error and no table", verifyFake{loops: []Step{okLoop()}, verifyErr: boom},
			"layup setup verify: boom", []string{"loop", "verify"}},
		{"layup setup verify gives an error with a table whose rows pass (the scratch tree is not removed)", verifyFake{loops: []Step{okLoop()}, table: goodTable(), verifyErr: &verify.CleanupError{Path: "/tmp/x", Err: boom}},
			"layup setup verify: the scratch tree /tmp/x is not removed: boom", []string{"loop", "verify"}},
		{"out/verify.tsv cannot be written", verifyFake{loops: []Step{okLoop()}, table: goodTable(), writeErr: boom},
			"out/verify.tsv: boom", []string{"loop", "verify", "write"}},
		{"a row is fail", verifyFake{loops: []Step{okLoop()}, table: badTable(verify.Row{Check: "adapted", Result: "fail", Reason: "docs/how-to.md:3: kit text"})},
			"layup setup verify: adapted is fail: docs/how-to.md:3: kit text", []string{"loop", "verify", "write"}},
		{"a row is not-active (the scratch tree could not be made): the first such row is named", verifyFake{loops: []Step{okLoop()}, table: badTable(verify.Row{Check: "markers", Result: "not-active", Reason: "scratch tree: add failed"}, verify.Row{Check: "sources", Result: "fail", Reason: "x"})},
			"layup setup verify: markers is not-active: scratch tree: add failed", []string{"loop", "verify", "write"}},
		{"the first of two rows that are not pass or clear is named", verifyFake{loops: []Step{okLoop()}, table: badTable(verify.Row{Check: "jobs", Result: "fail", Reason: "no job"}, verify.Row{Check: "sources", Result: "fail", Reason: "no source"})},
			"layup setup verify: jobs is fail: no job", []string{"loop", "verify", "write"}},
		{"the loop after the table does not end done", verifyFake{loops: []Step{okLoop(), {"setup", "fail", "S15: the row link-lint is fail: x"}}, table: goodTable()},
			"S15: the row link-lint is fail: x", []string{"loop", "verify", "write", "loop"}},
		{"S15 did not complete: the files of out/ are not there", verifyFake{loops: []Step{okLoop(), okLoop()}, table: goodTable(), filesErr: boom},
			"the files of out/: boom", []string{"loop", "verify", "write", "loop", "files"}},
		{"the tree of the setup cannot be made", verifyFake{loops: []Step{okLoop(), okLoop()}, table: goodTable(), treeErr: boom},
			"setup/tree.tsv: boom", []string{"loop", "verify", "write", "loop", "files", "tree"}},
		{"the records commit fails", verifyFake{loops: []Step{okLoop(), okLoop()}, table: goodTable(), commitErr: boom},
			"the records commit: boom", []string{"loop", "verify", "write", "loop", "files", "tree", "commit"}},
	} {
		f := c.f
		got := f.job().do()
		if got.Name != "verify" || got.Result != "fail" || got.Detail != c.detail || !slices.Equal(f.called, c.calls) {
			t.Errorf("%s: %+v after the calls %q; want fail %q after %q", c.name, got, f.called, c.detail, c.calls)
		}
		if !strings.Contains(strings.Join(f.called, ","), "commit") && len(f.commits) != 0 {
			t.Errorf("%s: a commit was made", c.name)
		}
	}
}

// A row that is clear is not a failure: the rows of a gate whose fixture is pending are clear.
func TestARowThatIsClearDoesNotFailTheStepVerify(t *testing.T) {
	f := &verifyFake{loops: []Step{okLoop(), okLoop()}, table: goodTable()}
	if got := f.job().do(); got.Result != "done" {
		t.Errorf("a table with a clear row: %+v; want done", got)
	}
}

// A table that the schema of out/verify.tsv refuses is a fail, before anything is written.
func TestATableThatItsSchemaRefusesFailsTheStepVerify(t *testing.T) {
	f := &verifyFake{loops: []Step{okLoop()}, table: verify.Table{Rows: []verify.Row{{Check: "pin", Result: "maybe", Reason: "—"}}}}
	if got := f.job().do(); got.Result != "fail" || !strings.HasPrefix(got.Detail, "layup setup verify: ") || !slices.Equal(f.called, []string{"loop", "verify"}) {
		t.Errorf("a table with the result maybe: %+v after %q; want fail after the loop and verify", got, f.called)
	}
}
