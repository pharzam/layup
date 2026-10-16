package run

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/session"
)

func sum(text string) string {
	h := sha256.Sum256([]byte(text))
	return hex.EncodeToString(h[:])
}

func TestTheKeptPathIsTheTaskTheSessionAndTheFile(t *testing.T) {
	if got, want := KeptPath("T-ab12", "S-1a2b3c4d", "plan.md"), "tasks/T-ab12/results/S-1a2b3c4d/plan.md"; got != want {
		t.Errorf("KeptPath = %q, want %q", got, want)
	}
	if KeptPath("T-ab12", "S-1a2b3c4d", "plan.md") == KeptPath("T-cd34", "S-1a2b3c4d", "plan.md") ||
		KeptPath("T-ab12", "S-1a2b3c4d", "plan.md") == KeptPath("T-ab12", "S-5e6f7a8b", "plan.md") ||
		KeptPath("T-ab12", "S-1a2b3c4d", "plan.md") == KeptPath("T-ab12", "S-1a2b3c4d", "review.md") {
		t.Error("a task, a session or a file name does not change the path")
	}
}

func TestACommentBodyIsRefusedBeforeThePost(t *testing.T) {
	three := strings.Repeat("€", MaxCommentChars) // three bytes each: 196,608 bytes, 65,536 characters
	for _, c := range []struct {
		name string
		body string
		ok   bool
	}{
		{"a plan", "# The plan\n", true},
		{"no final newline", "text", true},
		{"crlf", "a\r\nb\r\n", true},
		{"exactly the limit in characters", strings.Repeat("a", MaxCommentChars), true},
		{"the limit in characters of three-byte letters (over it in bytes)", three, true},
		{"one character over the limit", strings.Repeat("a", MaxCommentChars+1), false},
		{"one three-byte letter over the limit", three + "€", false},
		{"empty", "", false},
		{"not valid UTF-8", "a\xffb", false},
	} {
		err := CheckCommentBody([]byte(c.body))
		var refusal HandoffRefusal
		switch {
		case c.ok && err != nil:
			t.Errorf("%s: refused: %v", c.name, err)
		case !c.ok && !errors.As(err, &refusal):
			t.Errorf("%s: %v; want a HandoffRefusal", c.name, err)
		case !c.ok && refusal.Reason != "handoff":
			t.Errorf("%s: the reason %q; want handoff", c.name, refusal.Reason)
		}
	}
}

// resultOf gives a session directory whose result/ holds the files by name.
func resultOf(t *testing.T, files map[string]string) session.Dir {
	t.Helper()
	d := session.Dir{Result: filepath.Join(t.TempDir(), "result")}
	if err := os.MkdirAll(d.Result, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, text := range files {
		if err := os.WriteFile(filepath.Join(d.Result, name), []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return d
}

func commentRow(n, file, sha string) []string {
	return []string{"comment", n, "result/" + file, sha, ""}
}

func TestTheKeptBytesAreTheBytesOfTheFile(t *testing.T) {
	for _, text := range []string{"a plan\n", "no final newline", "crlf\r\nlines\r\n", "trailing  \n  spaces \t\n", "ü 𝄞 non-ASCII\n", "\n", "a\x00nul"} {
		d := resultOf(t, map[string]string{"plan.md": text})
		files, bad := keepComments(d, "T-ab12", "S-1a2b3c4d", [][]string{{"status", "1", "completed", "", "done"}, commentRow("1", "plan.md", sum(text))})
		got, ok := files["tasks/T-ab12/results/S-1a2b3c4d/plan.md"]
		if bad || !ok || string(got) != text || len(files) != 1 {
			t.Errorf("%q: kept %q (found %v), bad %v; want the bytes alone at the kept path", text, got, ok, bad)
		}
	}
}

func TestACommentRowIsCheckedByItsHashAgainstTheResultFile(t *testing.T) {
	d := resultOf(t, map[string]string{"plan.md": "the plan\n", "review.md": "the review\n"})
	outside := filepath.Join(t.TempDir(), "x")
	os.WriteFile(outside, []byte("the plan\n"), 0o644)
	os.Symlink(outside, filepath.Join(d.Result, "link.md"))
	good := commentRow("1", "plan.md", sum("the plan\n"))
	for _, c := range []struct {
		name string
		rows [][]string
		bad  bool
		kept []string // the paths kept
	}{
		{"no comment row", [][]string{{"status", "1", "completed", "", "done"}}, false, nil},
		{"a row whose hash is the file's", [][]string{good}, false, []string{"plan.md"}},
		{"a hash that differs", [][]string{commentRow("1", "plan.md", sum("another\n"))}, true, nil},
		{"the hash of the text with its final newline trimmed", [][]string{commentRow("1", "plan.md", sum("the plan"))}, true, nil},
		{"a file that is missing", [][]string{commentRow("1", "none.md", sum("the plan\n"))}, true, nil},
		{"a file that is a link, of the right bytes", [][]string{commentRow("1", "link.md", sum("the plan\n"))}, true, nil},
		{"two good rows", [][]string{good, commentRow("2", "review.md", sum("the review\n"))}, false, []string{"plan.md", "review.md"}},
		{"one good row and one bad: the good one is kept, the failure is reported", [][]string{good, commentRow("2", "review.md", sum("no\n"))}, true, []string{"plan.md"}},
		{"two rows of one file", [][]string{good, commentRow("2", "plan.md", sum("the plan\n"))}, false, []string{"plan.md"}},
		{"an artifact row is not read here", [][]string{{"artifact", "1", "out.txt", sum("x"), ""}}, false, nil},
	} {
		files, bad := keepComments(d, "T-ab12", "S-1a2b3c4d", c.rows)
		if bad != c.bad || len(files) != len(c.kept) {
			t.Errorf("%s: bad %v, %d files; want %v, %d", c.name, bad, len(files), c.bad, len(c.kept))
			continue
		}
		for _, name := range c.kept {
			if _, ok := files[KeptPath("T-ab12", "S-1a2b3c4d", name)]; !ok {
				t.Errorf("%s: %s is not kept", c.name, name)
			}
		}
	}
}
