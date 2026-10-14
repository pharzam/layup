package run

// This file holds the prose session of the step setup (docs/spec/scaffold.md,
// The prose rows; row 54 of the plan, task T-eahr, #277): one session of the role
// product-owner for each stop of the runner with F- rows, whose output is one
// file for each row, at result/files/<path>.

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/setup"
	"github.com/pharzam/layup/internal/verify"
)

// proseOutput is the directory of result/ that holds the files of the session.
const proseOutput = "files"

// proseNames gives the paths that the session writes: the path of each F- row of
// the stop, in the order of the table. A row that is no F- row of a path of the
// target below its root (the runner gives the path as the question's name and as
// its where), or a path twice, is a defect of the runner, and the step is fail.
func proseNames(rows []setup.StopRow) ([]string, error) {
	if len(rows) == 0 {
		return nil, errors.New("the stop has no row for a file")
	}
	var names []string
	for _, r := range rows {
		path, ok := strings.CutPrefix(r.Question, "F-")
		if !ok || r.Where != path || !fs.ValidPath(path) || path == "." {
			return nil, fmt.Errorf("the row %s of %s is not the F- row of a path of the target", r.Question, r.Step)
		}
		if slices.Contains(names, path) {
			return nil, fmt.Errorf("the path %s is in the stop twice", path)
		}
		names = append(names, path)
	}
	return names, nil
}

// proseKeep gives the Keep of the session for the paths that the stop named: the
// records of the files at setup/inputs/files/<path>. The reader of the output has
// checked the set; a file that no name covers or a missing one is refused here
// too, so a records commit never holds a file that the stop did not name.
func proseKeep(names []string) func(id string, files map[string][]byte) (map[string][]byte, error) {
	return func(id string, files map[string][]byte) (map[string][]byte, error) {
		if len(files) != len(names) {
			return nil, fmt.Errorf("%d files; the stop named %d", len(files), len(names))
		}
		kept := map[string][]byte{}
		for _, n := range names {
			data, ok := files[n]
			if !ok {
				return nil, fmt.Errorf("the file %s is missing", n)
			}
			kept["setup/inputs/files/"+n] = data
		}
		return kept, nil
	}
}

// stopTable is the stop table of the runner, the block setup-stop, as the
// session's input/stop.tsv.
func stopTable(rows []setup.StopRow) ([]byte, error) {
	var b bytes.Buffer
	if err := (setup.Result{Stops: rows}).Write(&b); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// proseText gives the prompt of the session: the paths of its input files, each
// on a line of its own (the vision brief only when the target has one), and of
// the directory of its output, then the rules of what it writes and of the
// checks that its files must pass, with the name of the target (task T-92lm,
// #326: the prompt gave the rules of the writing and none of the checks, and the
// step setup failed on the first real target).
func proseText(name string, input map[string]string, output string) []byte {
	var b strings.Builder
	b.WriteString("You write the prose files of a project's repository, for LAYUP. Change no file in this repository and make no commit.\n\n")
	b.WriteString("The problem statement: " + input["problem-statement.md"] + "\n")
	if v := input["vision.md"]; v != "" {
		b.WriteString("The vision brief: " + v + "\n")
	}
	b.WriteString("The stop table: " + input["stop.tsv"] + "\n")
	b.WriteString("The output directory: " + output + "\n\n")
	b.WriteString(proseRule)
	b.WriteString(proseChecks(name, output))
	return []byte(b.String())
}

// proseRule is the fixed text of the prompt after its paths. The ask of an
// adapted file begins Adapt (internal/setup), and S14 refuses an input that loses
// a place of a marker of the file (docs/spec/setup.md, S14).
const proseRule = `The stop table has one row for each file to write: its where column is the
path of the file in the project, and its ask column says what the file must be.
Read the problem statement, the vision brief when there is one, and the project
in this directory, which is as it is now. Write each file of the table in the
output directory, at its path below it (the file docs/glossary.md goes to
<output directory>/docs/glossary.md); make the directories that you need. Write
the files of the table and no other: a file that is missing, a file that the table
does not name, a link, a file over 1 MiB and a file that is not valid UTF-8 refuse
the whole output. Write plain text or Markdown in the language of the problem
statement, with no placeholder and no angle quote of your own. A row whose ask
says Adapt names a file that this directory already has: give that file with
your changes, and keep each marker that it holds (a text between the angle
quotes U+2039 and U+203A), in every place, byte for byte, with its line ends.

`

// proseChecksText is the fixed text of the rules of the checks, with {whole},
// {open}, {stems}, {exempt}, {cased}, {allowed}, {name}, {lines}, {output} and
// {aid}. Each rule says which files it governs; its home is the function of
// internal/verify that the test of the prompt names, and the phrases, the
// names, the files and the lines come from there.
const proseChecksText = `Code checks the files that you give when you end, and the first check that fails
refuses all of them, so check each file yourself before you end. These are the
rules of the checks, each with the files that it governs.

1. Words, in every file. No file may hold any of these phrases. The check reads each
   capital ASCII letter as the small one (the baseline's name below is the one
   exception), each line end as a space and each run of spaces and tabs as one space; a blank
   line ends a paragraph, and no phrase runs across it. A phrase is a hit unless an
   ASCII letter, an ASCII digit or an underscore stands just before it. A marker is
   read as one unit, so a word inside a marker is no hit, and you change no marker.
   Rewrite each sentence that holds a phrase so that it speaks of this project and of
   its people, in its own terms, and keep what the sentence says: do not drop a
   sentence for its words alone.
   Whole words (a hit unless an ASCII letter, an ASCII digit or an underscore follows: "kits" is a hit, "kitten" is not):
{whole}
   Open at the end (a hit whatever follows: "your projectile" is a hit):
{open}
   Stems (a hit as the start of any word of ASCII letters, with an angle quote U+2039 after it, in a marker or alone, before the next period: "fillable" is one):
{stems}
   The check passes the names {exempt}, and any word that starts with one of them, though they hold the word "kit": leave them as they are.
   The phrase {cased}, with its capital A, is read as it is, in its letter case. It may stand only in these files: {allowed}. In any other file write "the baseline".
2. Citations. The checks read them in docs/onboarding-for-engineers.md,
   docs/glossary.md and docs/guardrails.md; keep them right in every file. A
   citation is "F-", four digits, "#" and a number, as in F-0001#3. It must name a
   fact of this project: the directory docs/facts/ holds one file
   F-<four digits>-<name>.md for that number, and the file has a line that starts
   with the number, a full stop and a space. List docs/facts/ and read what is
   there. A citation that names no fact of this project is an unsupported
   attribution. Remove it, or point it to a fact of this project that supports the
   text. Never remove a rule, a step, a definition, an entry or a row of the
   baseline for that reason alone: keep its text and change only its attribution.
   Remove a sentence or a row only when it is an example that belongs to the
   baseline's own project and says nothing about this one.
3. README.md must hold the name of this repository, {name}, as text; must name the
   branch layup-records; and must hold the link ](docs/setup/armature.pin). Neither
   README.md nor AGENTS.md may hold any of these texts anywhere, also inside a
   longer line: {lines}.
4. docs/onboarding-for-engineers.md must hold no angle quote at all, not even a
   marker, and must hold the link ](facts/problem-statement-brief.md).
5. In docs/guardrails.md each entry is a bullet that starts with "- **Inv-", a
   number and "**", and runs to the next entry or heading. Each entry must have a
   line with "Check: ", and the text from the last "Check: " of the first such line
   to the end of that line must be exactly "no check yet", or the path of a file of
   this project that exists, a space and "(hook)", or that path, a space, "(ci:", a
   job name and ")". A file with no entry passes.
6. Markers: a row that asks you to adapt a file keeps each marker of it, as said
   above, and no file holds a marker of your own.

Before you end, search each file that you wrote for the phrases of rule 1. This is
an aid, not the check: it joins every line, also across a blank line; it takes a
letter that is not ASCII as part of a word, as the check does not; it shows a word
inside a marker and the two names that the check passes, which are no hits; and it
does not look for the baseline's name, which you search in its own letter case.
For example:

for f in $(find {output} -name '*.md'); do tr -s ' \t\n' ' ' < "$f" | grep -o -i -E '{aid}' | sed "s|^|$f: |"; done

A hit of a stem is a hit of the check only when an angle quote U+2039 follows it
before the next period: read each by eye, and rewrite its sentence when one does.
Repeat until every hit that is left is a stem with no angle quote U+2039 after it
in its sentence, a word inside a marker, or one of the two names. Then read your
files against the other rules.
Then end.
`

// proseChecks gives the rules of the checks for the target name and the output
// directory: the phrases of check adapted, the names that it passes, the files
// that may name the baseline and the lines of check identity are those of
// internal/verify, and the search command is built from the same phrases.
func proseChecks(name, output string) string {
	var whole, open, stems, wholeWords, openWords, stemWords, lines, exempt []string
	var cased string
	for _, p := range verify.AdaptedPhrases() {
		switch {
		case p.Cased:
			cased = p.Say
		case p.Tail == verify.TailStem:
			stems = append(stems, "   - words that start with \""+p.Stem()+"\" ("+p.Say+")")
			stemWords = append(stemWords, p.Stem())
		case p.Tail == verify.TailOpen:
			open = append(open, "   - "+p.Say)
			openWords = append(openWords, p.Words()...)
		default:
			whole = append(whole, "   - "+p.Say)
			wholeWords = append(wholeWords, p.Words()...)
		}
	}
	for _, l := range verify.KitPhrases() {
		lines = append(lines, `"`+l+`"`)
	}
	for _, n := range verify.AdaptedExempt() {
		exempt = append(exempt, `"`+n+`"`)
	}
	aid := `\<(` + strings.Join(wholeWords, "|") + `)\>|\<(` + strings.Join(openWords, "|") + `)|\<(` + strings.Join(stemWords, "|") + `)[a-z]*`
	return strings.NewReplacer("{whole}", strings.Join(whole, "\n"), "{open}", strings.Join(open, "\n"), "{stems}", strings.Join(stems, "\n"),
		"{exempt}", strings.Join(exempt, " and "), "{cased}", cased, "{allowed}", strings.Join(verify.AdaptedAllowed(), ", "),
		"{name}", name, "{lines}", strings.Join(lines, ", "), "{output}", output, "{aid}", aid).Replace(proseChecksText)
}
