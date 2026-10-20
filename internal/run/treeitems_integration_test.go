//go:build integration

package run

// The tree items of a handoff, on a real tree (task T-yeav, #404): the files
// of docs/spec/ at the head, the rendering of a record at a head and at an
// earlier revision, the head of a session and of code, through the production
// Sessions and the records.

import (
	"context"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/spec"
)

const (
	treeSection = "## REQ-005 — the handoff\nthe body of the section\n"
	treeHeading = "## REQ-005 — the handoff\n\n"
	treePRDPath = "docs/prd/PRD.md"
	treeFacts   = "docs/facts/numbered.md"
	treePRDPre  = "# PRD\n\n**Non-goals.** None.\n\n## 6. Functional requirements\n\n| REQ | Statement | MoSCoW | Phase | Facts |\n| REQ-001 | The app does one thing. |  |  | F-0001#1 |\n"
	treeFactsAt = "1. **need:** One.\n"
)

var (
	treeAnswer = spec.Answer{Reqs: map[string]spec.ReqAnswer{"REQ-001": {Class: "Must", Phase: "1"}}}
	treeRecs   = []rules.Record{
		{Path: treeFacts, Kind: rules.Facts, Version: [][]string{{"1", "need", "in", "7"}}, Texts: map[string]string{"1": "One."}},
		{Path: treePRDPath, Kind: rules.PRD, Answer: treeAnswer},
	}
)

// commitIn commits the files in the run's clone, on its branch, and gives the commit.
func commitIn(l *loopRun, files map[string]string) string {
	l.t.Helper()
	for p, d := range files {
		full := filepath.Join(l.r.Clone, p)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			l.t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(d), 0o644); err != nil {
			l.t.Fatal(err)
		}
	}
	gitOut(l.t, l.r.Clone, "add", "-A")
	gitOut(l.t, l.r.Clone, "commit", "-q", "-m", "work")
	return strings.TrimSpace(gitOut(l.t, l.r.Clone, "rev-parse", "HEAD"))
}

// specWorld runs a developer session whose commit holds docs/spec/loop.md with
// text, and gives the world, the production Sessions, the session and its bound head.
func specWorld(t *testing.T, text string) (*loopRun, *Sessions, string, string) {
	t.Helper()
	l := newLoop(t, nil)
	id := l.dev(1, "fake", map[string]string{"docs/spec/loop.md": text}, okResult)
	head := l.bound(id)
	if head == "" {
		t.Fatal("the session has no bound head")
	}
	return l, l.prod(), id, head
}

func specInput(id, head string, claim bool) HandoffInput {
	return HandoffInput{Task: "T-ab12", Session: id, From: "developer (specification)", To: "verifier", Items: map[string]bool{"spec-sections": claim},
		Artifacts: []handoff.Artifact{{Kind: "commits", Named: head, Found: head}, {Kind: "task-file", Path: "docs/tasks/T-ab12.md", Named: hashA, Found: hashA}}}
}

// renderWorld commits the PRD before the bet, then a head with the PRD and the
// facts record rendered (the PRD with the forged text, when given), and gives
// the production Sessions, the two commits.
func renderWorld(t *testing.T, forged string) (*loopRun, *Sessions, string, string) {
	t.Helper()
	l := newLoop(t, nil)
	pre := commitIn(l, map[string]string{treePRDPath: treePRDPre})
	prd, finds := spec.RenderPRD([]byte(treePRDPre), treeAnswer)
	if len(finds) != 0 {
		t.Fatalf("the fixture does not render: %v", finds)
	}
	if forged != "" {
		prd = []byte(strings.Replace(string(prd), "does one thing", forged, 1))
	}
	return l, l.prod(), pre, commitIn(l, map[string]string{treePRDPath: string(prd), treeFacts: treeFactsAt})
}

func renderInput(head, pre string, claim bool) HandoffInput {
	var arts []handoff.Artifact
	for _, p := range []string{treeFacts, treePRDPath} {
		h := strings.Repeat("a", 64)
		arts = append(arts, handoff.Artifact{Kind: "rendering", Path: p, Named: h, Found: h})
	}
	return HandoffInput{Task: "T-ab12", Attempt: 1, From: "code (rendered task)", To: "verifier", CodeWritten: true, Head: head, PreRev: pre,
		Render: treeRecs, Artifacts: arts, Items: map[string]bool{"rendering-equal": claim}}
}

func verdict(t *testing.T, s *Sessions, in HandoffInput) string {
	t.Helper()
	v, err := s.RecordHandoff(context.Background(), in)
	if err != nil {
		t.Fatalf("RecordHandoff: %v", err)
	}
	return v.Detail
}

// entry is a path of a tree, of the mode, with data (the object of a gitlink).
type entry struct{ mode, path, data string }

// treeOf makes a repository with one commit that holds the entries.
func treeOf(t *testing.T, entries ...entry) (string, string) {
	t.Helper()
	dir := t.TempDir()
	gitOut(t, dir, "init", "-q", "-b", "main")
	for _, e := range entries {
		object := e.data
		if e.mode != "160000" {
			src := filepath.Join(t.TempDir(), "blob")
			if err := os.WriteFile(src, []byte(e.data), 0o644); err != nil {
				t.Fatal(err)
			}
			object = strings.TrimSpace(gitOut(t, dir, "hash-object", "-w", src))
		}
		gitOut(t, dir, "update-index", "--add", "--cacheinfo", e.mode+","+object+","+e.path)
	}
	tree := strings.TrimSpace(gitOut(t, dir, "write-tree"))
	return dir, strings.TrimSpace(gitOut(t, dir, "commit-tree", tree, "-m", "x"))
}

func TestTheSpecFilesAreTheRegularMdBlobsOfDocsSpecAtTheHead(t *testing.T) {
	t.Parallel()
	gitlink := strings.Repeat("1", 40)
	for _, c := range []struct {
		name string
		e    entry
		want bool
	}{
		{"a regular .md file (the control)", entry{"100644", "docs/spec/a.md", treeSection}, true},
		{"a file below docs/spec/", entry{"100644", "docs/spec/sub/b.md", treeSection}, false},
		{"a link whose target text is the section", entry{"120000", "docs/spec/c.md", treeSection}, false},
		{"a gitlink", entry{"160000", "docs/spec/d.md", gitlink}, false},
		{"an executable file", entry{"100755", "docs/spec/e.md", treeSection}, false},
		{"a file of another extension", entry{"100644", "docs/spec/f.txt", treeSection}, false},
		{"a sibling directory", entry{"100644", "docs/specs/g.md", treeSection}, false},
	} {
		dir, head := treeOf(t, c.e)
		files, err := specFilesAt(dir, head)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got := specSections(files, []string{"REQ-005"}); got != c.want {
			t.Errorf("%s: sections %v, want %v (files %v)", c.name, got, c.want, slices.Sorted(maps.Keys(files)))
		}
	}
	dir, head := treeOf(t, entry{"100644", "docs/spec/a.md", treeSection}, entry{"100644", "docs/spec/sub/b.md", "x\n"}, entry{"120000", "docs/spec/c.md", "a.md"})
	files, err := specFilesAt(dir, head)
	if err != nil || len(files) != 1 || string(files["a.md"]) != treeSection {
		t.Errorf("files %v, %v; want only a.md with its bytes", slices.Sorted(maps.Keys(files)), err)
	}
	if _, err := specFilesAt(dir, strings.Repeat("a", 40)); err == nil {
		t.Error("a head the clone lacks reads as no file")
	}
}

// The demo of row 66n (group (b)): a handoff that lists a tree item is valid
// or invalid by the tree, whatever the caller says, for each item and each tree.
func TestAHandoffOfTheTreeItemsIsValidOrInvalidByTheTreeNotByTheCaller(t *testing.T) {
	t.Parallel()
	for _, c := range []struct{ name, text, want string }{
		{"the section is at the bound head", treeSection, "valid"},
		{"a heading with no body", treeHeading, "invalid spec-sections"},
	} {
		_, s, id, head := specWorld(t, c.text)
		for _, claim := range []bool{false, true} {
			if got := verdict(t, s, specInput(id, head, claim)); got != c.want {
				t.Errorf("%s, the caller says %v: %s; want %s", c.name, claim, got, c.want)
			}
		}
	}
	for _, c := range []struct{ name, forged, want string }{
		{"the rendering is at the head", "", "valid"},
		{"a byte of the head is forged", "does two things", "invalid rendering-equal"},
	} {
		_, s, pre, head := renderWorld(t, c.forged)
		for _, claim := range []bool{false, true} {
			if got := verdict(t, s, renderInput(head, pre, claim)); got != c.want {
				t.Errorf("%s, the caller says %v: %s; want %s", c.name, claim, got, c.want)
			}
		}
	}
}

func TestASessionHandoffIgnoresASuppliedHeadAndTheCloneHead(t *testing.T) {
	t.Parallel()
	// The section is at the bound head only: a later commit of the clone, and a
	// Head that the caller supplies, hold none.
	l, s, id, head := specWorld(t, treeSection)
	later := commitIn(l, map[string]string{"docs/spec/loop.md": treeHeading})
	in := specInput(id, head, false)
	in.Head = later
	if got := verdict(t, s, in); got != "valid" {
		t.Errorf("the section at the bound head, a supplied Head without it: %s; want valid", got)
	}
	// The section is at the supplied Head and at the clone's HEAD only.
	l, s, id, head = specWorld(t, treeHeading)
	later = commitIn(l, map[string]string{"docs/spec/loop.md": treeSection})
	in = specInput(id, head, true)
	if got := verdict(t, s, in); got != "invalid spec-sections" {
		t.Errorf("the section at the clone's HEAD only: %s; want invalid spec-sections", got)
	}
	in.Head = later
	if got := verdict(t, s, in); got != "invalid spec-sections" {
		t.Errorf("the section at the supplied Head only: %s; want invalid spec-sections", got)
	}
}

func TestARenderedTaskHandoffNeedsAPreImageThatIsNotTheHead(t *testing.T) {
	t.Parallel()
	_, s, pre, head := renderWorld(t, "")
	if got := verdict(t, s, renderInput(head, head, true)); got != "invalid rendering-equal" {
		t.Errorf("the head as its own pre-image: %s; want invalid rendering-equal", got)
	}
	if got := verdict(t, s, renderInput(head, "", true)); got != "invalid rendering-equal" {
		t.Errorf("no PreRev: %s; want invalid rendering-equal", got)
	}
	if got := verdict(t, s, renderInput(head, pre, false)); got != "valid" {
		t.Errorf("the control, the earlier revision: %s; want valid", got)
	}
}

func TestAFailedReadOfTheTreeIsAnErrorAndWritesNoEvent(t *testing.T) {
	t.Parallel()
	// A head the clone lacks.
	l, s, _, head := renderWorld(t, "")
	n := len(l.events())
	in := renderInput(strings.Repeat("a", 40), head, true)
	if v, err := s.RecordHandoff(context.Background(), in); err == nil || v.Valid || len(l.events()) != n {
		t.Errorf("a head the clone lacks: %+v, %v, %d new events; want an error and none", v, err, len(l.events())-n)
	}
	// A PreRev that is no revision.
	if v, err := s.RecordHandoff(context.Background(), renderInput(head, "nope", true)); err == nil || v.Valid || len(l.events()) != n {
		t.Errorf("a PreRev that is no revision: %+v, %v, %d new events; want an error and none", v, err, len(l.events())-n)
	}
	// A register the store cannot read.
	l, s, id, head := specWorld(t, treeSection)
	n = len(l.events())
	good := s.Store
	s.Store = failRead{good, tasksPath}
	if v, err := s.RecordHandoff(context.Background(), specInput(id, head, true)); err == nil || v.Valid || len(l.events()) != n {
		t.Errorf("a register that fails: %+v, %v, %d new events; want an error and none", v, err, len(l.events())-n)
	}
	s.Store = good
	if got := verdict(t, s, specInput(id, head, false)); got != "valid" {
		t.Errorf("the control, the store read again: %s; want valid", got)
	}
}

func TestASpecBlobThatCannotBeReadIsAnErrorAndWritesNoEvent(t *testing.T) {
	t.Parallel()
	l, s, id, head := specWorld(t, treeSection)
	n := len(l.events())
	blob := strings.TrimSpace(gitOut(t, l.r.Clone, "rev-parse", head+":docs/spec/loop.md"))
	loose := filepath.Join(l.r.Clone, ".git", "objects", blob[:2], blob[2:])
	if _, err := os.Stat(loose); err != nil {
		t.Fatalf("the blob is not a loose object (%v); the test cannot remove it alone", err)
	}
	if err := os.Remove(loose); err != nil {
		t.Fatal(err)
	}
	// The listing of the tree still succeeds: only the read of the blob fails.
	if v, err := s.RecordHandoff(context.Background(), specInput(id, head, true)); err == nil || v.Valid || len(l.events()) != n {
		t.Errorf("a blob that cannot be read: %+v, %v, %d new events; want an error and none", v, err, len(l.events())-n)
	}
}
