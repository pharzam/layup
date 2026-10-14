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

// flatPrompt gives the text of a prompt with each run of white space as one space, so a phrase is found across the line ends of the
// text as it is wrapped.
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

// The prompt gives the session the words that check adapted refuses, from the one home of them (task T-92lm, #326, the two
// failures of the demo of M2d: the session kept the words of the baseline in about twenty files): each phrase of internal/verify
// as that package says it, what the check does with letter case, line ends, blanks and a marker, and the files that may name the
// baseline (rule 2). Home of the rules: adaptedFiles, hitsOf and adHits in internal/verify.
func TestTheProsePromptGivesTheWordsOfCheckAdapted(t *testing.T) {
	text := proseTextOf()
	var cased []string
	for _, p := range verify.AdaptedPhrases() {
		if p.Cased {
			cased = append(cased, p.Say)
			continue
		}
		wantAll(t, text, "- "+p.Say+" ")
	}
	if len(cased) != 1 {
		t.Fatalf("%d phrases are read in their letter case; the prompt says one", len(cased))
	}
	wantAll(t, text,
		"Code checks the files that you give when you end, and the first check that fails refuses all of them, so check each file yourself before you end.",
		"1. Words, in every file. No file may hold any of these phrases.",
		"Letter case does not matter (the baseline's name below is the one exception), each line end counts as a space, a run of blanks counts as one space, and a letter, a digit or an underscore just before a phrase makes it another word.",
		"A marker is read as one unit, so a word inside a marker is no hit, and you change no marker.",
		"Rewrite each sentence that holds a phrase so that it speaks of this project and of its people, in its own terms, and keep what the sentence says: do not drop a sentence for its words alone.",
		"The phrase "+cased[0]+", with its capital A, is read as it is, in its letter case. It may stand only in these files: "+strings.Join(verify.AdaptedAllowed(), ", ")+`. In any other file write "the baseline".`)
}

// The prompt gives the rule of a citation with the files that read it, and keeps the rules of the baseline when a citation goes (the
// glossary of the first run held F-0007#3, an example of the baseline: S08 failed; condition 3 of the plan review of T-92lm: an
// unresolved citation is no leave to drop a rule). Home: resolveFact and citationFindings in internal/verify.
func TestTheProsePromptGivesTheRuleOfACitation(t *testing.T) {
	wantAll(t, proseTextOf(),
		"2. Citations. The checks read them in docs/onboarding-for-engineers.md, docs/glossary.md and docs/guardrails.md; keep them right in every file.",
		"A citation is \"F-\", four digits, \"#\" and a number, as in F-0001#3. It must name a fact of this project: the directory docs/facts/ holds one file F-<four digits>-<name>.md for that number, and the file has a line that starts with the number, a full stop and a space.",
		"List docs/facts/ and read what is there.",
		"A citation that names no fact of this project is an unsupported attribution. Remove it, or point it to a fact of this project that supports the text.",
		"Never remove a rule, a step, a definition, an entry or a row of the baseline for that reason alone: keep its text and change only its attribution.",
		"Remove a sentence or a row only when it is an example that belongs to the baseline's own project and says nothing about this one.")
}

// The prompt gives the rules of check identity with the name of the target of the run, the branch of the records and the lines of the
// baseline, each as internal/verify has it (the second failure of the demo: README.md held no name). Home: identityTarget and
// identityFindings in internal/verify.
func TestTheProsePromptGivesTheRulesOfCheckIdentity(t *testing.T) {
	var lines []string
	for _, l := range verify.KitPhrases() {
		lines = append(lines, `"`+l+`"`)
	}
	if len(lines) == 0 {
		t.Fatal("internal/verify has no line of the baseline")
	}
	wantAll(t, proseTextOf(),
		"3. README.md must hold the name of this repository, "+proseTarget+", as text; must name the branch layup-records; and must hold the link ](docs/setup/armature.pin).",
		"Neither README.md nor AGENTS.md may hold any of these lines: "+strings.Join(lines, ", ")+".")
}

// The prompt gives the rules of check onboarding and of check guardrails, with their exact forms and the one file each governs.
// Home: onboardingFindings and guardrailsFindings in internal/verify.
func TestTheProsePromptGivesTheRulesOfCheckOnboardingAndGuardrails(t *testing.T) {
	wantAll(t, proseTextOf(),
		"4. docs/onboarding-for-engineers.md must hold no angle quote at all, not even a marker, and must hold the link ](facts/problem-statement-brief.md).",
		`5. In docs/guardrails.md each entry is a bullet that starts with "- **Inv-", a number and "**".`,
		`The first "Check: " of an entry must be exactly "no check yet", or the path of a file of this project that exists, a space and "(hook)", or that path, a space, "(ci:", a job name and ")".`,
		"A file with no entry passes.",
		"6. Markers: a row that asks you to adapt a file keeps each marker of it, as said above, and no file holds a marker of your own.")
}

// The prompt asks the session to search its files before it ends, with a command that reads a phrase across a line end, from the
// forms of internal/verify, and says that it is an aid and not the check (condition 1 of the plan review).
func TestTheProsePromptAsksForASearchBeforeTheEnd(t *testing.T) {
	text := proseTextOf()
	var words []string
	for _, p := range verify.AdaptedPhrases() {
		if !p.Cased {
			for _, w := range p.Words() {
				words = append(words, "-e '"+w+"'")
			}
		}
	}
	wantAll(t, text,
		"Before you end, search each file that you wrote for each phrase of rule 1. This is an aid, not the check:",
		"it reads a phrase across a line end, it cannot read the two phrases about a marker (a form of fill or replace in a sentence that goes on to a marker before its full stop), which you read by eye, and it does not look for the baseline's name, which you search in its own letter case.",
		`for f in $(find /s/result/files -name '*.md'); do tr -s ' \t\n' ' ' < "$f" | grep -o -i -w -F `+strings.Join(words, " ")+` | sed "s|^|$f: |"; done`,
		"Read your files against the other rules, and repeat until nothing is left.")
	if !strings.HasSuffix(strings.TrimSpace(text), "Then end.") {
		t.Errorf("the prompt does not end with the end of the session:\n%s", text)
	}
}

// The name of the target is a parameter of the text: two targets give two prompts, each with its own name and with no trace of the
// other's.
func TestTheProsePromptHasTheNameOfItsTarget(t *testing.T) {
	input := map[string]string{"stop.tsv": "/s/input/stop.tsv"}
	a, b := flatPrompt(proseText("acme/target", input, "/s/o")), flatPrompt(proseText("pharzam/layup-uat-m2d", input, "/s/o"))
	if !strings.Contains(a, "the name of this repository, acme/target, as text") || strings.Contains(a, "layup-uat-m2d") {
		t.Errorf("the prompt of acme/target:\n%s", a)
	}
	if !strings.Contains(b, "the name of this repository, pharzam/layup-uat-m2d, as text") || strings.Contains(b, "acme/target") {
		t.Errorf("the prompt of pharzam/layup-uat-m2d:\n%s", b)
	}
}

// The prompt holds no angle quote of its own: it tells the session to write none, and says the two by their code points.
func TestTheProsePromptHoldsNoAngleQuote(t *testing.T) {
	if text := string(proseText(proseTarget, map[string]string{"stop.tsv": "/s/input/stop.tsv"}, "/s/result/files")); strings.ContainsAny(text, "\u2039\u203a") {
		t.Errorf("the prompt holds an angle quote:\n%s", text)
	}
}
