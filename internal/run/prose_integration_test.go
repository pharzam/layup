//go:build integration

package run

import (
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/verify"
)

// The search command of the prompt of the prose session, run as printed by the real sh, find, tr, grep and sed (task T-92lm, #326,
// finding 1 of round 1): it finds each form of each phrase of check adapted in any letter case and across a line end, a stem with any
// ASCII letters after it and an open end with more letters, and no whole word inside a longer one, in .md files only.
func TestTheSearchCommandOfTheProsePromptFindsWhatCheckAdaptedFinds(t *testing.T) {
	dir := t.TempDir()
	cmd := regexp.MustCompile(`(?m)^for f in \$\(find .* -name '\*\.md'\); do .*; done$`).FindString(string(proseText(proseTarget, map[string]string{"stop.tsv": "/s/input/stop.tsv"}, dir)))
	if cmd == "" {
		t.Fatal("the prompt has no search command")
	}
	text, want, none := "# Notes\n\nA line with a Kit.\nThe adopters, and this\ntemplate; skips this\nsection and your projectile.\n", []string{"kit", "adopters", "this template", "skips this section", "your project"}, []string(nil)
	for _, p := range verify.AdaptedPhrases() {
		for _, w := range p.Words() {
			switch {
			case p.Cased:
			case p.Tail == verify.TailStem:
				text += "Then " + strings.ToUpper(w) + "able to the marker.\n"
				want = append(want, w+"able")
			case p.Tail == verify.TailWord:
				text += "We keep " + strings.ToUpper(w) + " here; " + w + "x and x" + w + " there.\n"
				want, none = append(want, w), append(none, w+"x", "x"+w)
			}
		}
	}
	writeFile(t, filepath.Join(dir, "a.md"), text)
	writeFile(t, filepath.Join(dir, "docs", "b.txt"), "the kit\n")
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()
	if err != nil {
		t.Fatalf("the search command: %v\n%s", err, out)
	}
	got := strings.ToLower(string(out))
	for _, w := range want {
		if !strings.Contains(got, strings.ToLower(w)) {
			t.Errorf("the search does not find %q:\n%s", w, out)
		}
	}
	for _, w := range none {
		if strings.Contains(got, ": "+strings.ToLower(w)+"\n") {
			t.Errorf("the search finds %q, which is no hit of the check:\n%s", w, out)
		}
	}
	if strings.Contains(got, "b.txt") {
		t.Errorf("the search reads a file that is not .md:\n%s", out)
	}
}
