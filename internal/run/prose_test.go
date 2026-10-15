package run

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/verify"
)

// The prose session of the step setup (docs/spec/scaffold.md, The prose rows;
// task T-eahr, #277): its stop table, the paths that it writes, its prompt, and
// the records that its files become.

func proseRows() []setup.StopRow {
	return []setup.StopRow{
		{Step: "S07", Question: "F-docs/onboarding-for-engineers.md", Ask: "Write the text of docs/onboarding-for-engineers.md for the target, and give it as inputs/files/docs/onboarding-for-engineers.md.", Where: "docs/onboarding-for-engineers.md"},
		{Step: "S07", Question: "F-docs/glossary.md", Ask: "Write the text of docs/glossary.md for the target, and give it as inputs/files/docs/glossary.md.", Where: "docs/glossary.md"},
	}
}

// The files that the session writes are the paths of the F- rows, in the order
// of the stop table; a row that is no F- row of its own path is a defect of the
// runner, and the step is fail.
func TestTheProseNamesAreThePathsOfTheFRows(t *testing.T) {
	names, err := proseNames(proseRows())
	if want := []string{"docs/onboarding-for-engineers.md", "docs/glossary.md"}; err != nil || !reflect.DeepEqual(names, want) {
		t.Errorf("the names %q, %v; want %q", names, err, want)
	}
	for name, rows := range map[string][]setup.StopRow{
		"no row":                       nil,
		"an M- row":                    {{Step: "S10", Question: "M-aaaaaaaa", Ask: "q", Where: "NOTES.md:2 <x>"}},
		"an M- row that names itself":  {{Step: "S10", Question: "M-aaaaaaaa", Ask: "q", Where: "M-aaaaaaaa"}},
		"the root as a path":           {{Step: "S07", Question: "F-.", Ask: "q", Where: "."}},
		"an F- row of another path":    {{Step: "S07", Question: "F-docs/a.md", Ask: "q", Where: "docs/b.md"}},
		"a path that leaves the tree":  {{Step: "S07", Question: "F-../a.md", Ask: "q", Where: "../a.md"}},
		"an absolute path":             {{Step: "S07", Question: "F-/a.md", Ask: "q", Where: "/a.md"}},
		"the same path twice":          {proseRows()[0], proseRows()[0]},
		"an F- row with no path":       {{Step: "S07", Question: "F-", Ask: "q", Where: ""}},
		"an F- row of an unclean path": {{Step: "S07", Question: "F-docs//a.md", Ask: "q", Where: "docs//a.md"}},
	} {
		if names, err := proseNames(rows); err == nil {
			t.Errorf("%s: %q; want an error", name, names)
		}
	}
}

// The files go to the records at setup/inputs/files/<path>, and only the files
// that the stop named: a file that no name covers, or that is missing, is a
// refusal of the output (the session reader has checked it; this is the second line).
func TestTheProseKeepGivesTheFilesAtTheirRecordsPaths(t *testing.T) {
	keep := proseKeep([]string{"docs/glossary.md", "README.md"})
	got, err := keep("S-1a2b3c4d", map[string][]byte{"docs/glossary.md": []byte("# G\n"), "README.md": []byte("# R\n")})
	want := map[string][]byte{"setup/inputs/files/docs/glossary.md": []byte("# G\n"), "setup/inputs/files/README.md": []byte("# R\n")}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("the records %q, %v; want %q", got, err, want)
	}
	for name, files := range map[string]map[string][]byte{
		"a file that no name covers": {"docs/glossary.md": []byte("g"), "README.md": []byte("r"), "extra.md": []byte("x")},
		"a missing file":             {"README.md": []byte("r")},
		"a file in place of a named": {"docs/glossary.md": []byte("g"), "other.md": []byte("x")},
		"no file":                    {},
	} {
		if got, err := keep("S-1a2b3c4d", files); err == nil || got != nil {
			t.Errorf("%s: %q, %v; want an error and no record", name, got, err)
		}
	}
}

// The prompt names, one on each line, the problem statement, the vision brief
// when there is one, the stop table and the directory of the output files.
func TestTheProsePromptNamesItsPaths(t *testing.T) {
	input := map[string]string{"problem-statement.md": "/s/input/problem-statement.md", "vision.md": "/s/input/vision.md", "stop.tsv": "/s/input/stop.tsv"}
	text := string(proseText(proseTarget, input, "/s/result/files"))
	for _, line := range []string{"The problem statement: /s/input/problem-statement.md", "The vision brief: /s/input/vision.md", "The stop table: /s/input/stop.tsv", "The output directory: /s/result/files"} {
		if !strings.Contains("\n"+text, "\n"+line+"\n") {
			t.Errorf("the prompt has no line %q:\n%s", line, text)
		}
	}
	delete(input, "vision.md")
	if text := string(proseText(proseTarget, input, "/s/result/files")); strings.Contains(text, "The vision brief:") || strings.Contains(text, "vision.md") {
		t.Errorf("the prompt of a target with no vision brief names its file:\n%s", text)
	}
}

// A file that the runner asks to adapt keeps each marker that it holds: S14 refuses an input that
// loses a place of one (docs/spec/setup.md), so the ban of angle quotes has that exception.
func TestTheProsePromptLetsAnAdaptedFileKeepItsMarkers(t *testing.T) {
	text := flatPrompt(proseText(proseTarget, map[string]string{"stop.tsv": "/s/input/stop.tsv"}, "/s/result/files"))
	for _, want := range []string{"A row whose ask says Adapt", "keep each marker that it holds", "U+2039 and U+203A", "in every place, byte for byte, with its line ends"} {
		if !strings.Contains(text, want) {
			t.Errorf("the prompt does not say %q:\n%s", want, text)
		}
	}
}

// The stop table of the session is the table that the runner printed, the block
// setup-stop, byte for byte.
func TestTheStopTableOfTheSessionIsTheRunnersTable(t *testing.T) {
	got, err := stopTable(proseRows())
	var want strings.Builder
	if werr := (setup.Result{Stops: proseRows()}).Write(&want); err != nil || werr != nil || string(got) != want.String() || !strings.HasPrefix(string(got), "step\tquestion\task\twhere\n") {
		t.Errorf("the stop table %q, %v; want %q, %v", got, err, want.String(), werr)
	}
}

// proseTarget is the name of the target of the prompts of these tests.
const proseTarget = "acme/target"

// flatPrompt gives a prompt with each run of white space as one space, so a sentence is found across the wrapped lines.
func flatPrompt(b []byte) string { return strings.Join(strings.Fields(string(b)), " ") }

// proseTextOf gives the prompt of the target proseTarget with the output directory of these tests.
func proseTextOf() string {
	return flatPrompt(proseText(proseTarget, map[string]string{"problem-statement.md": "/s/input/problem-statement.md", "stop.tsv": "/s/input/stop.tsv"}, "/s/result/files"))
}

// wantAll reports each sentence that the prompt does not hold.
func wantAll(t *testing.T, text string, sentences ...string) {
	t.Helper()
	for _, want := range sentences {
		if !strings.Contains(text, want) {
			t.Errorf("the prompt does not say %q:\n%s", want, text)
		}
	}
}

// The prompt gives the words that check adapted refuses (task T-92lm, #326: the demo of M2d failed on them), each phrase of
// internal/verify in the group of its tail, how the check reads case, line ends, blanks and a marker, the names it passes and the
// files that may name the baseline. Round 1 found the first text looser than the check. Home: hitsOf, paragraphHits, adHits.
func TestTheProsePromptGivesTheWordsOfCheckAdapted(t *testing.T) {
	text := proseTextOf()
	const (
		hWhole = `Whole words (a hit unless an ASCII letter, an ASCII digit or an underscore follows: "kits" is a hit, "kitten" is not):`
		hOpen  = `Open at the end (a hit whatever follows: "your projectile" is a hit):`
		hStems = `Stems (a hit as the letters at the start of a word, then any text that has no period in it, then an angle quote U+2039, in a marker or alone: "fillable" and "fill2" are the start of one):`
		hEnd   = "The check passes the names"
	)
	seg := func(from, to string) string { // the group between its heading and the next
		i, j := strings.Index(text, from), strings.Index(text, to)
		if i < 0 || j < i {
			t.Fatalf("the prompt has no group between %q and %q:\n%s", from, to, text)
		}
		return text[i:j]
	}
	groups := map[string]string{verify.TailWord: seg(hWhole, hOpen), verify.TailOpen: seg(hOpen, hStems), verify.TailStem: seg(hStems, hEnd)}
	var cased []string
	for _, p := range verify.AdaptedPhrases() {
		if p.Cased {
			cased = append(cased, p.Say)
			continue
		}
		line := "- " + p.Say + " "
		if p.Tail == verify.TailStem {
			line = `- words that start with "` + p.Stem() + `" (` + p.Say + `)`
		}
		for tail, g := range groups {
			if has := strings.Contains(g, line); has != (tail == p.Tail) {
				t.Errorf("the group %q holds %q: %v; its phrase is of the group %q", tail, line, has, p.Tail)
			}
		}
	}
	if len(cased) != 1 {
		t.Fatalf("%d phrases are read in their letter case; the prompt says one", len(cased))
	}
	var passed []string
	for _, n := range verify.AdaptedExempt() {
		passed = append(passed, `"`+n+`"`)
	}
	wantAll(t, text, "Code checks the files that you give when you end, and the first check that fails refuses all of them, so check each file yourself before you end.",
		"1. Words, in every file. No file may hold any of these phrases.",
		"The check reads each capital ASCII letter as the small one (the baseline's name below is the one exception), each line end as a space and each run of spaces and tabs as one space; a blank line ends a paragraph, and no phrase runs across it.",
		"A phrase is a hit unless an ASCII letter, an ASCII digit or an underscore stands just before it.",
		"A marker is read as one unit, so a word inside a marker is no hit, and you change no marker.",
		"Rewrite each sentence that holds a phrase so that it speaks of this project and of its people, in its own terms, and keep what the sentence says: do not drop a sentence for its words alone.",
		hWhole, hOpen, hStems, "The check passes the names "+strings.Join(passed, " and ")+`, and any word that starts with one of them, though they hold the word "kit": leave them as they are.`,
		"The phrase "+cased[0]+", with its capital A, is read as it is, in its letter case. It may stand only in these files: "+strings.Join(verify.AdaptedAllowed(), ", ")+`. In any other file write "the baseline".`)
}

// The prompt gives the rules of citations (and keeps the baseline's rules when one goes: the plan review's condition 3), of check
// identity (the name of the run's target, the branch, the texts anywhere in a line), of onboarding and of guardrails (the exact forms
// of Check:), each with its files. Home: resolveFact, identityTarget, identityFindings, onboardingFindings, guardrailsFindings.
func TestTheProsePromptGivesTheRulesOfTheOtherChecks(t *testing.T) {
	var lines []string
	for _, l := range verify.KitPhrases() {
		lines = append(lines, `"`+l+`"`)
	}
	wantAll(t, proseTextOf(),
		"2. Citations. The checks read them in docs/onboarding-for-engineers.md, docs/glossary.md and docs/guardrails.md; keep them right in every file.",
		"A citation is \"F-\", four digits, \"#\" and a number, as in F-0001#3. It must name a fact of this project: the directory docs/facts/ holds one file F-<four digits>-<name>.md for that number, and the file has a line that starts with the number, a full stop and a space.",
		"List docs/facts/ and read what is there. A citation that names no fact of this project is an unsupported attribution. Remove it, or point it to a fact of this project that supports the text.",
		"Never remove a rule, a step, a definition, an entry or a row of the baseline for that reason alone: keep its text and change only its attribution.",
		"Remove a sentence or a row only when it is an example that belongs to the baseline's own project and says nothing about this one.",
		"3. README.md must hold the name of this repository, "+proseTarget+", as text; must name the branch layup-records; and must hold the link ](docs/setup/armature.pin).",
		"Neither README.md nor AGENTS.md may hold any of these texts anywhere, also inside a longer line: "+strings.Join(lines, ", ")+".",
		"4. docs/onboarding-for-engineers.md must hold no angle quote at all, not even a marker, and must hold the link ](facts/problem-statement-brief.md).",
		`5. In docs/guardrails.md each entry is a bullet that starts with "- **Inv-", a number and "**", and runs to the next entry or heading.`,
		`Each entry must have a line with "Check: ", and the text from the last "Check: " of the first such line to the end of that line must be exactly "no check yet", or the path of a file of this project that exists, a space and "(hook)", or that path, a space, "(ci:", a job name and ")", with no space followed by an opening parenthesis in the path or in the job name (the check splits the value at the first of them for the path and at the last for the gate). A file with no entry passes.`,
		"6. Markers: a row that asks you to adapt a file keeps each marker of it, as said above, and no file holds a marker of your own.")
}

// The prompt asks for a search before the end, with a command built from the phrases of internal/verify by their tails, and says
// what the aid is not (finding 1 of round 1: it differs from the check).
func TestTheProsePromptAsksForASearchBeforeTheEnd(t *testing.T) {
	text := proseTextOf()
	var whole, open, stems []string
	for _, p := range verify.AdaptedPhrases() {
		switch {
		case p.Cased:
		case p.Tail == verify.TailStem:
			stems = append(stems, p.Stem())
		case p.Tail == verify.TailOpen:
			open = append(open, p.Words()...)
		default:
			whole = append(whole, p.Words()...)
		}
	}
	regex := `\<(` + strings.Join(whole, "|") + `)\>|\<(` + strings.Join(open, "|") + `)|\<(` + strings.Join(stems, "|") + `)[a-z]*`
	wantAll(t, text, "Before you end, search each file that you wrote for the phrases of rule 1. This is an aid, not the check:",
		"it joins every line, also across a blank line; it takes a letter that is not ASCII as part of a word, as the check does not; it shows a word inside a marker and the two names that the check passes, which are no hits; and it does not look for the baseline's name, which you search in its own letter case.",
		"A hit of a stem is a hit of the check only when an angle quote U+2039 follows it before the next period: read each by eye, and rewrite its sentence when one does.",
		"Repeat until every hit that is left is a stem with no angle quote U+2039 after it before the next period, a word inside a marker, or one of the two names. Then read your files against the other rules.",
		`for f in $(find /s/result/files -name '*.md'); do tr -s ' \t\n' ' ' < "$f" | grep -o -i -E '`+regex+`' | sed "s|^|$f: |"; done`)
	if !strings.HasSuffix(strings.TrimSpace(text), "Then end.") {
		t.Errorf("the prompt does not end with the end of the session:\n%s", text)
	}
}

// The name of the target is a parameter of the text: two targets give two prompts, each with its own name and no trace of the other's.
func TestTheProsePromptHasTheNameOfItsTarget(t *testing.T) {
	input := map[string]string{"stop.tsv": "/s/input/stop.tsv"}
	a, b := flatPrompt(proseText("acme/target", input, "/s/o")), flatPrompt(proseText("pharzam/layup-uat-m2d", input, "/s/o"))
	if !strings.Contains(a, "the name of this repository, acme/target, as text") || strings.Contains(a, "layup-uat-m2d") ||
		!strings.Contains(b, "the name of this repository, pharzam/layup-uat-m2d, as text") || strings.Contains(b, "acme/target") {
		t.Errorf("the prompts of two targets:\n%s\n%s", a, b)
	}
}

// The prompt holds no angle quote: it tells the session to write none, and names the two by their code points.
func TestTheProsePromptHoldsNoAngleQuote(t *testing.T) {
	if text := string(proseText(proseTarget, map[string]string{"stop.tsv": "/s/input/stop.tsv"}, "/s/result/files")); strings.ContainsAny(text, "\u2039\u203a") {
		t.Errorf("the prompt holds an angle quote:\n%s", text)
	}
}
