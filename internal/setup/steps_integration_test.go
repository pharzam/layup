//go:build integration

package setup

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// noGaps is what internal/cli hands over for a problem statement with no gap.
var noGaps = Brief{Gaps: []byte("id\trule\tline\texcerpt\tquestion\n"), Sum: "0123"}

// firstSteps gives the steps S01 to S04 of this version, and a stub of each
// later step, so a test of S01 to S04 ends at S05.
func firstSteps(b Brief, c Calls) map[string]Step {
	m, stubs := Steps(b, c), Stubs()
	for id := range m {
		if id > "S04" {
			m[id] = stubs[id]
		}
	}
	return m
}

// newWorkArea makes a work area at dir/work whose answers name the baseline at
// url, with no record.
func newWorkArea(t *testing.T, dir, url string) string {
	t.Helper()
	w := filepath.Join(dir, "work")
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-1"
	var b bytes.Buffer
	if err := tsv.Write(&b, work.AnswersSchema, [][]string{
		{"S01-stack", "go", "operator", at, ""}, {"S01-name", standin.Name, "operator", at, ""},
		{"S01-visibility", "public", "operator", at, ""}, {"S01-baseline", url, "operator", at, ""},
	}); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(filepath.Join(w, filepath.FromSlash(work.AnswersPath)), b.Bytes()); err != nil {
		t.Fatal(err)
	}
	return w
}

// value gives a value of the record of the work area at w.
func value(t *testing.T, w, step, name string) string {
	t.Helper()
	r, err := work.ReadRecord(w)
	if err != nil {
		t.Fatal(err)
	}
	v, _ := r.Value(step, name)
	return v
}

// S01 to S04 on a stand-in baseline by its file:// URL: S02 resolves the
// commit once and records it; the copy has no .git of the clone, and its root
// commit on main has the tree pin.tree, by the identity of the run at
// pin.time; S04 makes one commit on layup-setup. A second run, after a new
// commit of the baseline, resolves nothing and changes nothing (O-136).
func TestS01ToS04OnABaseline(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "baseline")
	url, commit, err := standin.Baseline(base)
	if err != nil {
		t.Fatal(err)
	}
	w := newWorkArea(t, tmp, url)
	res, err := Run(w, firstSteps(noGaps, Calls{}), Who, noStep)
	if err != nil || !slices.Equal(res.Results()[:5], []string{Done, Done, Done, Done, NotActive}) || res.Steps[4].Evidence != "not built yet" {
		t.Fatalf("the run: %v, %+v; want S01 to S04 done and S05 not built yet", err, res.Steps)
	}
	target := filepath.Join(w, work.TargetPath)
	tree := gitOut(t, base, "rev-parse", commit+"^{tree}")
	at := value(t, w, "S02", "pin.time")
	if _, err := time.Parse(timeForm, at); err != nil || value(t, w, "S02", "pin.commit") != commit || value(t, w, "S02", "pin.tree") != tree {
		t.Errorf("the pin rows: %q %q %q; want %s, %s and a time", value(t, w, "S02", "pin.commit"), value(t, w, "S02", "pin.tree"), at, commit, tree)
	}
	if got := gitOut(t, target, "log", "--format=%s|%an <%ae>|%aI|%cI|%T", "main"); got != "chore: the unmodified baseline at "+commit+"|"+Who.Name+" <"+Who.Email+">|"+at+"|"+at+"|"+tree {
		t.Errorf("main: %q; want the one root commit with the tree pin.tree", got)
	}
	if err := exec.Command("git", "-C", target, "cat-file", "-e", commit).Run(); err == nil {
		t.Errorf("the target holds the commit %s of the baseline: the .git of the clone stayed", commit)
	}
	if _, err := os.Stat(target + ".part"); err == nil {
		t.Errorf("target.part stayed")
	}
	if got := gitOut(t, target, "log", "--format=%s", "main..layup-setup"); got != "chore: setup S04" {
		t.Errorf("layup-setup: %q; want the one commit of S04", got)
	}
	head := gitOut(t, target, "rev-parse", "layup-setup")
	if err := os.WriteFile(filepath.Join(base, "new.md"), []byte("a new file\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitOut(t, base, "add", "new.md")
	gitOut(t, base, "-c", "user.name=t", "-c", "user.email=t@layup.invalid", "commit", "-q", "-m", "a new commit")
	again, err := Run(w, firstSteps(noGaps, Calls{}), Who, noStep)
	if err != nil || !slices.Equal(again.Steps, res.Steps) || value(t, w, "S02", "pin.commit") != commit || gitOut(t, target, "rev-parse", "layup-setup") != head {
		t.Errorf("a second run: %v, %+v; want the same table, the same pin and no commit", err, again.Steps)
	}
}

// A baseline whose URL asks for a login fails at S02 with no prompt (K31): git
// asks the server, and the run ends with S02 fail.
func TestALoginURLFailsWithNoPrompt(t *testing.T) {
	asked := make(chan bool, 16)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		asked <- true
		w.Header().Set("WWW-Authenticate", `Basic realm="layup"`)
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()
	w := newWorkArea(t, t.TempDir(), srv.URL+"/target.git")
	type ran struct {
		res Result
		err error
	}
	done := make(chan ran, 1)
	go func() { r, err := Run(w, firstSteps(noGaps, Calls{}), Who, noStep); done <- ran{r, err} }()
	select {
	case r := <-done:
		if r.err != nil || r.res.Steps[1].Result != Fail || !strings.HasPrefix(r.res.Steps[1].Evidence, "git ls-remote "+srv.URL+"/target.git HEAD: ") ||
			!strings.HasSuffix(r.res.Steps[1].Evidence, ": terminal prompts disabled") || len(asked) == 0 {
			t.Errorf("a login URL: %v, %+v, %d requests; want S02 fail after a request, with no prompt", r.err, r.res.Steps, len(asked))
		}
	case <-time.After(time.Minute):
		t.Fatal("the run waits after a minute, as for a prompt")
	}
	if _, err := os.Stat(filepath.Join(w, work.TargetPath)); err == nil {
		t.Errorf("a target exists after S02 failed")
	}
}

// A step that stopped in its middle (D8 of #86): S02 removes its own
// target.part; a target that exists while S02 is not done is an input error,
// and layup removes nothing; S03 takes its own root commit of a run that
// stopped, and a target with another history is an input error.
func TestAStepThatStoppedInItsMiddle(t *testing.T) {
	tmp := t.TempDir()
	url, _, err := standin.Baseline(filepath.Join(tmp, "baseline"))
	if err != nil {
		t.Fatal(err)
	}
	w := newWorkArea(t, tmp, url)
	if err := writeFile(filepath.Join(w, "target.part", "left.md"), []byte("x\n")); err != nil {
		t.Fatal(err)
	}
	m := firstSteps(noGaps, Calls{})
	s := m["S03"]
	s.Run = func(in Input) Outcome { // the run stops after the commit of S03
		if o := runS03(in); o.Kind != Done {
			return o
		}
		return Outcome{Kind: Fail, Evidence: "a stop after the commit"}
	}
	m["S03"] = s
	if res, err := Run(w, m, Who, noStep); err != nil || res.Steps[2].Result != Fail {
		t.Fatalf("the run that stops: %v, %+v", err, res.Steps)
	}
	target := filepath.Join(w, work.TargetPath)
	if _, err := os.Stat(filepath.Join(target, "left.md")); err == nil {
		t.Errorf("the file of target.part is in the copy")
	}
	root := gitOut(t, target, "rev-parse", "main")
	res, err := Run(w, firstSteps(noGaps, Calls{}), Who, noStep)
	if err != nil || res.Steps[2].Result != Done || res.Steps[3].Result != Done || gitOut(t, target, "rev-list", "--count", "main") != "1" || gitOut(t, target, "rev-parse", "main") != root {
		t.Errorf("the next run: %v, %+v; want S03 done with its own commit, and no second root commit", err, res.Steps)
	}

	w2 := newWorkArea(t, filepath.Join(tmp, "2"), url)
	if err := writeFile(filepath.Join(w2, "target", "mine.md"), []byte("mine\n")); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(w2, firstSteps(noGaps, Calls{}), Who, noStep); err == nil || err.Error() != "target of the work area exists, and S02 is not done: remove it, or start again in a new work area" {
		t.Errorf("a target before S02: %v; want the input error", err)
	}
	if _, err := os.Stat(filepath.Join(w2, "target", "mine.md")); err != nil {
		t.Errorf("layup removed a file of the target: %v", err)
	}

	for i, change := range [][]string{
		{"commit", "-q", "--allow-empty", "-m", "another commit"},  // a second commit
		{"commit", "-q", "--amend", "-m", "the baseline, by hand"}, // one commit with the tree pin.tree and another message
	} {
		w := newWorkArea(t, filepath.Join(tmp, fmt.Sprint(3+i)), url)
		if _, err := Run(w, m, Who, noStep); err != nil {
			t.Fatal(err)
		}
		gitOut(t, filepath.Join(w, "target"), append([]string{"-c", "user.name=t", "-c", "user.email=t@layup.invalid"}, change...)...)
		if _, err := Run(w, firstSteps(noGaps, Calls{}), Who, noStep); err == nil || err.Error() != "target of the work area has a history that S03 did not make: start again in a new work area" {
			t.Errorf("a target with another history (%q): %v; want the input error", change, err)
		}
	}
}

// The undo of a step whose evidence fails (D12 of #86): layup-setup goes back
// to the root commit, the record has no done row of S04, and the next run does
// S04 again on the branch that exists, with one commit.
func TestTheUndoOfAnEvidenceThatFails(t *testing.T) {
	tmp := t.TempDir()
	url, _, err := standin.Baseline(filepath.Join(tmp, "baseline"))
	if err != nil {
		t.Fatal(err)
	}
	w := newWorkArea(t, tmp, url)
	reason := "pin: fail: a broken pin"
	var names []string
	checks := func(dir string, n []string) string { names = n; return reason }
	res, err := Run(w, firstSteps(noGaps, Calls{Checks: checks}), Who, noStep)
	target := filepath.Join(w, work.TargetPath)
	if err != nil || res.Steps[3] != (StepRow{"S04", "layup-setup", Fail, reason}) || !slices.Equal(names, []string{"pin", "facts"}) {
		t.Fatalf("an evidence that fails: %v, %+v, the checks %q", err, res.Steps, names)
	}
	if gitOut(t, target, "rev-parse", "layup-setup") != gitOut(t, target, "rev-parse", "main") || value(t, w, "S04", "done") != "" {
		t.Errorf("after the undo: layup-setup is not at the root commit, or S04 has a done row")
	}
	reason = ""
	if res, err := Run(w, firstSteps(noGaps, Calls{Checks: checks}), Who, noStep); err != nil || res.Steps[3].Result != Done ||
		gitOut(t, target, "log", "--format=%s", "main..layup-setup") != "chore: setup S04" {
		t.Errorf("the next run: %v, %+v; want S04 done with one commit on layup-setup", err, res.Steps)
	}
}
