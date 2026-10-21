//go:build integration

package git

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The shared indexes (T-5ve2, #461) are the repository's own files, read from
// the root of the tree, so that the test pins the real .gitattributes and not
// a copy of it.
const (
	completedPath    = "docs/tasks/completed.md"
	traceabilityPath = "docs/tests/traceability.md"
	prdPath          = "docs/prd/PRD-0001-layup.md"
)

var indexFixtures = []string{".gitattributes", completedPath, traceabilityPath, prdPath,
	"docs/tasks/backlog.md", "docs/guardrails.md", "internal/tsv/blocks_integration_test.go"}

// indexRepo makes a repository in a temporary directory that holds, on main,
// the real fixtures of indexFixtures, and gives it.
func indexRepo(t *testing.T, home string) (string, []string) {
	t.Helper()
	dir, env := t.TempDir(), plain(home)
	gitOK(t, dir, env, "init", "-q", "-b", "main")
	files := map[string]string{}
	for _, name := range indexFixtures {
		b, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(name)))
		must(t, err)
		files[name] = string(b)
	}
	write(t, dir, files)
	gitOK(t, dir, env, "add", "-A")
	gitOK(t, dir, env, "commit", "-q", "-m", "chore: base")
	return dir, env
}

// edit rewrites one file of the repository and commits it on the current
// branch.
func edit(t *testing.T, dir string, env []string, name, msg string, f func(string) string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(name))
	b, err := os.ReadFile(path)
	must(t, err)
	must(t, os.WriteFile(path, []byte(f(string(b))), 0o644))
	gitOK(t, dir, env, "commit", "-q", "-a", "-m", msg)
}

// insertAfter puts text on the line after the first match of anchor.
func insertAfter(t *testing.T, anchor, text string) func(string) string {
	re := regexp.MustCompile(anchor)
	return func(s string) string {
		loc := re.FindStringIndex(s)
		if loc == nil {
			t.Fatalf("no match of %q", anchor)
		}
		return s[:loc[1]] + text + "\n" + s[loc[1]:]
	}
}

const (
	logAnchor   = `(?m)^## Log\n\n`
	tableAnchor = `(?m)^\|-+\|.*\n`
)

// branchAdding makes a branch from main that adds one record to each shared
// index, the way the documented rules say: a line directly under "## Log", a
// row directly below the table's separator. It leaves the branch checked out.
func branchAdding(t *testing.T, dir string, env []string, branch, id string) (string, string) {
	t.Helper()
	line := "- **2026-10-21** — **" + id + "** — a task that closes out"
	row := "| `Test" + id + "` (`internal/x/x_test.go`) | unit | — | — | — | — | " + id + " | green |"
	gitOK(t, dir, env, "checkout", "-q", "-b", branch, "main")
	edit(t, dir, env, completedPath, "docs: "+id+" the completed line", insertAfter(t, logAnchor, line))
	edit(t, dir, env, traceabilityPath, "docs: "+id+" the traceability row", insertAfter(t, tableAnchor, row))
	return line, row
}

func hasConflictMarker(s string) bool {
	return regexp.MustCompile(`(?m)^(<{7}|>{7}) `).MatchString(s)
}

// The two shared indexes merge with no conflict when two branches each add a
// record, the merge keeps all four records, and the later-landing branch's
// record is above the earlier one in both files ("most recent first", B before
// A). The files that no merge attribute may reach keep git's ordinary merge: two
// edits of one PRD row still conflict.
func TestTheSharedIndexesMergeWithNoConflictInLandingOrder(t *testing.T) {
	home := isolate(t)

	t.Run("the two indexes", func(t *testing.T) {
		dir, env := indexRepo(t, home)
		aLine, aRow := branchAdding(t, dir, env, "T-aaaa", "T-aaaa")
		gitOK(t, dir, env, "checkout", "-q", "main")
		bLine, bRow := branchAdding(t, dir, env, "T-bbbb", "T-bbbb")
		gitOK(t, dir, env, "checkout", "-q", "main")
		gitOK(t, dir, env, "merge", "-q", "--no-ff", "-m", "Merge A", "T-aaaa") // A lands first.
		gitOK(t, dir, env, "checkout", "-q", "T-bbbb")
		oldHead := strings.TrimSpace(gitOK(t, dir, env, "rev-parse", "HEAD"))

		// B, now behind main, merges it: B is "ours".
		out, err := gitIn(dir, env, "merge", "-q", "--no-ff", "-m", "Merge main", "main")
		if err != nil {
			t.Fatalf("git merge main into B: %v\n%s", err, out)
		}
		merged := strings.TrimSpace(gitOK(t, dir, env, "rev-parse", "HEAD"))

		for _, c := range []struct{ path, a, b, anchor string }{
			{completedPath, aLine, bLine, logAnchor}, {traceabilityPath, aRow, bRow, tableAnchor}} {
			b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(c.path)))
			must(t, err)
			text := string(b)
			if hasConflictMarker(text) {
				t.Errorf("%s holds a conflict marker", c.path)
			}
			for _, rec := range []string{c.a, c.b} {
				if n := strings.Count(text, rec+"\n"); n != 1 {
					t.Errorf("%s holds %q %d times, want 1", c.path, rec, n)
				}
			}
			// B before A, directly under the anchor: most recent first.
			loc := regexp.MustCompile(c.anchor).FindStringIndex(text)
			if loc == nil || !strings.HasPrefix(text[loc[1]:], c.b+"\n"+c.a+"\n") {
				t.Errorf("%s: B's record is not directly above A's under its anchor", c.path)
			}
			// The numstat from the old head: the merge adds A's record and removes no line.
			if got := strings.TrimSpace(gitOK(t, dir, env, "diff", "--numstat", oldHead, merged, "--", c.path)); got != "1\t0\t"+c.path {
				t.Errorf("git diff --numstat %s: %q, want one added line and none removed", c.path, got)
			}
		}
		// Every line of the table, from its header to the end, is a table row.
		tb, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(traceabilityPath)))
		must(t, err)
		lines := strings.Split(strings.TrimRight(string(tb), "\n"), "\n")
		in := false
		for _, l := range lines {
			if strings.HasPrefix(l, "| Test ID ") {
				in = true
			}
			if in && !strings.HasPrefix(l, "|") {
				t.Errorf("a line inside the table is not a row: %q", l)
			}
		}
		if got := gitOK(t, dir, env, "diff", "--check", oldHead, merged); got != "" {
			t.Errorf("git diff --check %s %s: %q", oldHead[:8], merged[:8], got)
		}

		// Then B lands: a merge into main with no conflict, and each record still once.
		gitOK(t, dir, env, "checkout", "-q", "main")
		gitOK(t, dir, env, "merge", "-q", "--no-ff", "-m", "Merge B", "T-bbbb")
		for _, name := range []string{completedPath, traceabilityPath} {
			b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(name)))
			must(t, err)
			for _, rec := range []string{aLine, bLine, aRow, bRow} {
				if hasConflictMarker(string(b)) || strings.Count(string(b), rec+"\n") > 1 {
					t.Errorf("%s: a marker, or %q more than once, after B lands", name, rec)
				}
			}
		}
	})

	t.Run("the files that stay under git's ordinary merge", func(t *testing.T) {
		dir, env := indexRepo(t, home)
		for _, name := range []string{"docs/tasks/backlog.md", prdPath, "docs/guardrails.md", "internal/tsv/blocks_integration_test.go"} {
			got := strings.TrimSpace(gitOK(t, dir, env, "check-attr", "merge", "--", name))
			if want := name + ": merge: unspecified"; got != want {
				t.Errorf("git check-attr merge: %q, want %q", got, want)
			}
		}
		// Two edits of one row of PRD section 12 conflict, and are not kept as two versions.
		editRow := func(tag string) func(string) string {
			return func(s string) string {
				i := strings.Index(s, "\n| REQ-005 | F-0003#45, F-0003#59 ")
				if i < 0 {
					t.Fatal("no REQ-005 matrix row in the PRD")
				}
				j := i + 1 + strings.Index(s[i+1:], "\n")
				return s[:j] + " " + tag + s[j:]
			}
		}
		gitOK(t, dir, env, "checkout", "-q", "-b", "T-aaaa")
		edit(t, dir, env, prdPath, "docs: A edits the row", editRow("TestA"))
		gitOK(t, dir, env, "checkout", "-q", "main")
		gitOK(t, dir, env, "checkout", "-q", "-b", "T-bbbb")
		edit(t, dir, env, prdPath, "docs: B edits the row", editRow("TestB"))
		out, err := gitIn(dir, env, "merge", "-q", "--no-ff", "-m", "Merge A", "T-aaaa")
		if err == nil || !strings.Contains(out, "CONFLICT") {
			t.Errorf("the merge of two edits of one PRD row: %v\n%s, want a conflict", err, out)
		}
	})
}
