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
	t.Parallel()
	dir := t.TempDir()
	cmd := regexp.MustCompile(`(?m)^for f in \$\(find .* -name '\*\.md'\); do .*; done$`).FindString(string(proseText(proseTarget, map[string]string{"stop.tsv": "/s/input/stop.tsv"}, dir)))
	if cmd == "" {
		t.Fatal("the prompt has no search command")
	}
	text, want, neg := "# Notes\n\nA line with a Kit.\nThe adopters, and this\ntemplate; skips this\nsection and your projectile.\n", []string{"kit", "adopters", "this template", "skips this section", "your project"}, ""
	for _, p := range verify.AdaptedPhrases() {
		for _, w := range p.Words() {
			switch {
			case p.Cased:
			case p.Tail == verify.TailStem:
				text += "Then " + strings.ToUpper(w) + "able to the marker.\n"
				want = append(want, w+"able")
			case p.Tail == verify.TailWord:
				text += "We keep " + strings.ToUpper(w) + " here.\n"
				want = append(want, w)
				if !strings.HasPrefix(w, "fill") { // "fill in" holds a stem, which the aid shows with no marker (the prompt says so)
					neg += w + "x, x" + w + ", " + w + "9 and _" + w + ".\n"
				}
			}
		}
	}
	writeFile(t, filepath.Join(dir, "a.md"), text)
	writeFile(t, filepath.Join(dir, "neg.md"), neg) // each whole word inside a longer word: no hit
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
	if strings.Contains(got, "neg.md") {
		t.Errorf("the search finds a whole word inside a longer word, which is no hit of the check:\n%s", out)
	}
	if strings.Contains(got, "b.txt") {
		t.Errorf("the search reads a file that is not .md:\n%s", out)
	}
}
