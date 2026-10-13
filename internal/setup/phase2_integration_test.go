//go:build integration

package setup

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/standin"
	"github.com/pharzam/layup/internal/tsv"
	"github.com/pharzam/layup/internal/work"
)

// S01 to S04 of phase 2 through the runner on a real stand-in baseline, in the
// work area that layup run builds: a clone with main at the root commit, and
// the answers of S01-stack and two Q- rows. The source of the pin is a URL that
// no clone can reach, so a step that resolved it would fail: the run resolves
// nothing, takes the pin rows of Phase2, and makes the one commit of S04 on
// layup-setup. A second run changes nothing.
func TestPhase2S01ToS04OnARealTarget(t *testing.T) {
	tmp := t.TempDir()
	base := filepath.Join(tmp, "baseline")
	url, commit, err := standin.Baseline(base)
	if err != nil {
		t.Fatal(err)
	}
	tree := gitOut(t, base, "rev-parse", commit+"^{tree}")
	w := filepath.Join(tmp, "work")
	gitOut(t, tmp, "clone", "-q", url, filepath.Join(w, work.TargetPath))
	const at = "https://github.invalid/stand-in/issues/1#issuecomment-"
	var b bytes.Buffer
	if err := tsv.Write(&b, work.AnswersSchema, [][]string{
		{"S01-stack", "go", "operator", at + "1", ""}, {"Q-001", "Go", "idea-owner", at + "2", ""}, {"Q-002", "under 200 ms", "idea-owner", at + "3", ""},
	}); err != nil {
		t.Fatal(err)
	}
	if err := writeFile(filepath.Join(w, filepath.FromSlash(work.AnswersPath)), b.Bytes()); err != nil {
		t.Fatal(err)
	}
	p := Phase2{Name: standin.Name, Visibility: "public", PinSource: "file:///no/such/baseline", PinCommit: commit, PinTree: tree, PinTime: "2026-10-13T05:00:00Z", BriefSum: "0123",
		Questions: []work.Question{{ID: "Q-001", Text: "Which technology stack?"}, {ID: "Q-002", Text: "What does fast mean?"}}}
	steps := StepsPhase2(p, Calls{})
	for id, stub := range Stubs() { // a test of S01 to S04 ends at S05
		if id > "S04" {
			steps[id] = stub
		}
	}
	res, err := Run(w, steps, Who, noStep)
	if err != nil || !slices.Equal(res.Results()[:5], []string{Done, Done, Done, Done, NotActive}) || res.Steps[4].Evidence != "not built yet" {
		t.Fatalf("the run: %v, %+v; want S01 to S04 done and S05 not built yet", err, res.Steps)
	}
	if got := []string{res.Steps[1].Actor, res.Steps[2].Actor, res.Steps[0].Actor, res.Steps[3].Actor}; !slices.Equal(got, []string{"layup-run", "layup-run", "layup-setup", "layup-setup"}) {
		t.Errorf("the actors of S02, S03, S01 and S04: %q; want layup-run for S02 and S03", got)
	}
	target := filepath.Join(w, work.TargetPath)
	for _, c := range []struct{ step, name, want string }{
		{"S01", "name", standin.Name}, {"S01", "visibility", "public"}, {"S01", "baseline", p.PinSource}, {"S01", "brief.sha256", "0123"},
		{"S02", "pin.commit", commit}, {"S02", "pin.tree", tree}, {"S02", "pin.time", p.PinTime},
	} {
		if got := value(t, w, c.step, c.name); got != c.want {
			t.Errorf("the record row %s %s = %q; want %q", c.step, c.name, got, c.want)
		}
	}
	when, _ := time.Parse(timeForm, p.PinTime) // %at, not %aI: the text of the date depends on the version of git
	if got := gitOut(t, target, "log", "--format=%s|%an|%at", "main..layup-setup"); got != "chore: setup S04|"+Who.Name+"|"+strconv.FormatInt(when.Unix(), 10) {
		t.Errorf("layup-setup: %q; want the one commit of S04, by the identity of the run at pin.time", got)
	}
	if pin, err := os.ReadFile(filepath.Join(target, "docs", "setup", "armature.pin")); err != nil || string(pin) != pinText(p.PinSource, commit, tree, "2026-10-13") {
		t.Errorf("armature.pin: %q, %v; want the file of the pin rows", pin, err)
	}
	if cmds, err := os.ReadFile(filepath.Join(w, filepath.FromSlash(work.CommandsPath))); err != nil || len(cmds) != 0 {
		t.Errorf("commands.sh: %q, %v; want an empty file, as S03 gives no command in phase 2", cmds, err)
	}
	head := gitOut(t, target, "rev-parse", "layup-setup")
	again, err := Run(w, steps, Who, noStep)
	if err != nil || !slices.Equal(again.Steps, res.Steps) || gitOut(t, target, "rev-parse", "layup-setup") != head {
		t.Errorf("a second run: %v, %+v; want the same table and no commit", err, again.Steps)
	}
}
