package verify

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

// adapted runs the core of check adapted on an in-memory tree that holds the
// files, each tracked.
func adapted(files map[string]string) []string {
	fsys := fstest.MapFS{}
	var paths []string
	for p, text := range files {
		fsys[p] = &fstest.MapFile{Data: []byte(text)}
		paths = append(paths, p)
	}
	slices.Sort(paths)
	return adaptedFindings(fsys, stubHistory{files: paths})
}

// Each of the 16 patterns of check_adapted has a hit, and a near miss at each
// of its boundaries (D4 of #88).
func TestTheHitsOfAdapted(t *testing.T) {
	const m = "\u2039x\u203a" // a marker
	for _, c := range []struct {
		text string
		want []string // the rule and the name of each hit, at line 1 of docs/a.md
	}{
		{"This kit is copied.", []string{"rule-1 kit"}},
		{"The Kits, kit-wide.", []string{"rule-1 kit"}},
		{"A toolkit, a skit, kit_x, kit9 and kitten.", nil},
		{"The checks kit-history and kit-linters run.", nil},
		{"A kit-lint run.", []string{"rule-1 kit"}},
		{"An adopter and the adopters.", []string{"rule-1 adopter"}},
		{"It was adopted.", nil},
		{"Copy the template, this templates and The Template.", []string{"rule-1 the template"}},
		{"A template; the template-based file; the templated text.", nil},
		{"Armature is the source.", []string{"rule-2 Armature"}},
		{"The armature.pin, Armatures and Armature_x.", nil},
		{"This step is optional.", []string{"rule-3 optional"}},
		{"It is optionally done.", nil},
		{"It skips this section.", []string{"rule-3 skip this section"}},
		{"Fill each `" + m + "` marker.", []string{"rule-3 fill \u2039"}},
		{"Filling the gap: see `" + m + "`.", []string{"rule-3 fill \u2039"}},
		{"Fill the gap. Then see `" + m + "`.", nil},
		{"See [it](#fill-in-x) and -fill `" + m + "`.", nil},
		{"Then replace the `" + m + "` value.", []string{"rule-3 replace \u2039"}},
		{"Replace it. Then see `" + m + "`.", nil},
		{"Fill in the table.", []string{"rule-3 fill in"}},
		{"The fill-in step, a fill ins and a fill in-x.", nil},
		{"Delete this section.", []string{"rule-3 delete this"}},
		{"Delete thistle.", nil},
		{"Delete the ones you do not use.", []string{"rule-3 delete the one you do not use"}},
		{"Adapt the list; it adapts.", []string{"rule-3 adapt"}},
		{"The text is adapted, an adaptation.", nil},
		{"Use what your project, your forge and your stack need.", []string{"rule-3 your forge", "rule-3 your project", "rule-3 your stack"}},
		{"Keep the ones you use.", []string{"rule-3 you use"}},
		{"The ones you used.", nil},
		{"See `\u2039the kit\u203a` and `\u2039an adopter\u203a`.", nil},
		// an open quote with no close quote after it is no marker, so the words after it are read (fix 1 of the first pilot, #97)
		{"See `\u2039the kit\u203a` and `\u2039an adopter`.", []string{"rule-1 adopter"}},
		{"Fill the `\u2039scheme\nhere\u203a` value.", []string{"rule-3 fill \u2039"}},
	} {
		var want []string
		for _, w := range c.want {
			want = append(want, w+": docs/a.md:1")
		}
		slices.Sort(want)
		if got := adapted(map[string]string{"docs/a.md": c.text + "\n"}); !slices.Equal(got, want) {
			t.Errorf("%q:\n got %q\nwant %q", c.text, got, want)
		}
	}
}

// A paragraph is one text: a hit is at the line where its word starts; a
// line of spaces ends the paragraph; spaces and tabs are one space, and a
// line loses the space at its start and at its end (condition 1 of the plan
// review), and its carriage return (D3 of #88). A marker runs to the first
// close quote after it, also across lines, and is one unit; an open quote
// with no close quote after it, or between two backticks, starts no unit (fix 1
// of the first pilot, #97).
func TestTheTextOfAdapted(t *testing.T) {
	for _, c := range []struct {
		name, text string
		want       []string
	}{
		{"a phrase across a line end", "A project that has no CI skips this\nsection.\n", []string{"rule-3 skip this section: docs/a.md:1"}},
		{"the word on the next line", "One line.\nThe kit.\n", []string{"rule-1 kit: docs/a.md:2"}},
		{"an empty line ends the paragraph", "the\n\ntemplate\n", nil},
		{"a line of spaces ends the paragraph", "the\n \t \ntemplate\n", nil},
		{"two lines are one text", "Copy the\ntemplate.\n", []string{"rule-1 the template: docs/a.md:1"}},
		{"spaces and tabs", "the \t  template\n", []string{"rule-1 the template: docs/a.md:1"}},
		{"CRLF", "a\r\nkit\r\n", []string{"rule-1 kit: docs/a.md:2"}},
		{"no last line end", "the kit", []string{"rule-1 kit: docs/a.md:1"}},
		{"an open quote that does not close", "see \u2039the kit and more\nkit\n", []string{"rule-1 kit: docs/a.md:1", "rule-1 kit: docs/a.md:2"}},
		{"a marker over two lines is one unit", "see \u2039the kit and\nmore kit\u203a here\nkit\n", []string{"rule-1 kit: docs/a.md:3"}},
		{"the mention starts no unit", "the `\u2039` kit\n", []string{"rule-1 kit: docs/a.md:1"}},
		{"a hit twice on one line is one finding", "kit and kit\n", []string{"rule-1 kit: docs/a.md:1"}},
		{"a space at the end of a line inside a phrase", "Fill \nin the table.\n", []string{"rule-3 fill in: docs/a.md:1"}},
		{"a tab at the end of a line inside a phrase", "Fill\t\nin the table.\n", []string{"rule-3 fill in: docs/a.md:1"}},
		{"a space at the start of a line", "Fill\n in the table.\n", []string{"rule-3 fill in: docs/a.md:1"}},
		{"hits at the end of a line and at the start of the next share a space", "The kit\nkit set.\n", []string{"rule-1 kit: docs/a.md:1", "rule-1 kit: docs/a.md:2"}},
		{"the eleven bytes from the hit decide the exemption", "The kit-historyx and kit-linters2 lines; the kit-histor line.\n", []string{"rule-1 kit: docs/a.md:1"}},
		{"a boundary that is not ASCII", "\u00e9kit\n", []string{"rule-1 kit: docs/a.md:1"}},
		{"kit-history after a boundary that is not ASCII", "\u00e9kit-history and \u00e9kit-linters\n", nil},
		{"upper case outside ASCII that is shorter in lower case", "\u0130\u0130\u0130\u0130\nkit\n", []string{"rule-1 kit: docs/a.md:2"}},
		{"upper case outside ASCII that is longer in lower case", "\u023a\u023a\u023a\u023a kit\nx\n", []string{"rule-1 kit: docs/a.md:1"}},
	} {
		if got := adapted(map[string]string{"docs/a.md": c.text}); !slices.Equal(got, c.want) {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// Check adapted reads each tracked .md file that AD_EXCLUDE does not match;
// rule 2 does not read a file of ad_allowed; the findings are in byte order,
// each once; a failed list is a finding (D1, D2 and D5 of #88).
func TestTheFilesOfAdapted(t *testing.T) {
	const text = "The kit; Armature.\n"
	var excluded = []string{"docs/facts/F-0001.md", "runs/T-1/a.md", "docs/tasks/T-1.md", "docs/tasks/completed.md",
		"docs/adr/0001-a.md", "docs/adr/0012-a.md", "docs/setup/record-T-1.md", "docs/adr/tests/a.md",
		"docs/ci/tests/a.md", "docs/links/tests/a.md", "docs/prd/tests/a.md", "docs/setup/tests/a/b.md",
		".githooks/tests/a.md", "internal/psb/testdata/a.md", "docs/a.txt", "docs/a.MD"}
	files := map[string]string{"docs/tasks/backlog.md": text, "docs/adr/0013-a.md": text, "README.md": text, "AGENTS.md": text}
	for _, p := range excluded {
		files[p] = text
	}
	want := []string{"rule-1 kit: AGENTS.md:1", "rule-1 kit: README.md:1", "rule-1 kit: docs/adr/0013-a.md:1",
		"rule-1 kit: docs/tasks/backlog.md:1", "rule-2 Armature: docs/adr/0013-a.md:1", "rule-2 Armature: docs/tasks/backlog.md:1"}
	if got := adapted(files); !slices.Equal(got, want) {
		t.Errorf("the files:\n got %q\nwant %q", got, want)
	}
	for p := range files {
		if p != "docs/a.MD" && p != "docs/a.txt" && !adExcludeRE.MatchString(p) && !slices.Contains([]string{"docs/tasks/backlog.md", "docs/adr/0013-a.md", "README.md", "AGENTS.md"}, p) {
			t.Errorf("%s is not excluded; the test lists it as excluded", p)
		}
	}
	// A listed path that is not a regular file of the tree is not read.
	fsys := fstest.MapFS{"d.md/a.txt": {Data: []byte("x")}, "docs/b.md": {Data: []byte("the kit\n")}}
	if got := adaptedFindings(fsys, stubHistory{files: []string{"d.md", "gone.md", "docs/b.md"}}); !slices.Equal(got, []string{"rule-1 kit: docs/b.md:1"}) {
		t.Errorf("a directory and a missing file: %q", got)
	}
	if got := adaptedFindings(fstest.MapFS{}, stubHistory{err: errors.New("fatal: not a git repository")}); !slices.Equal(got, []string{"git: cannot list the tracked files"}) {
		t.Errorf("a failed list: %q", got)
	}
}

// S14 refuses an adapted input that loses a marker that the file before it
// has on a line that check adapted flags (fix 2 of the first pilot, #97: the
// prose step lost seven of them, so S10 never asked them). A marker over more
// lines counts when one of its lines is flagged; the finding names the line of
// its open quote and its key. A marker that the input keeps somewhere byte for
// byte, or a marker of a line that check adapted does not flag, is no finding.
func TestLostMarkers(t *testing.T) {
	const before = "# The kit\n\nName the models: `\u2039name your reasoning-tier models\u203a` for your project.\n\nThe value \u2039x\u203a stays.\n\nReport `\u2039model, effort,\ntokens and elapsed time\u203a` as your project does.\n"
	const kept = "# The project\n\nName the models: `\u2039name your reasoning-tier models\u203a`.\n\nThe value \u2039x\u203a stays.\n\nReport `\u2039model, effort,\ntokens and elapsed time\u203a` here.\n"
	for _, c := range []struct {
		name, after string
		want        []string
	}{
		{"each marker kept", kept, nil},
		{"a marker of a flagged line lost", strings.Replace(kept, "`\u2039name your reasoning-tier models\u203a`", "the reasoning-tier models", 1),
			[]string{"docs/e.md:3: the input loses the marker \u2039name your reasoning-tier models\u203a"}},
		{"a marker over two lines that lost its open quote", strings.Replace(kept, "`\u2039model, effort,", "`model, effort,", 1),
			[]string{"docs/e.md:7: the input loses the marker \u2039model, effort, tokens and elapsed time\u203a"}},
		{"a marker of a line that check adapted does not flag", strings.Replace(kept, "\u2039x\u203a", "8080", 1), nil},
		{"a marker kept on another line", "Name the models.\n\n`\u2039name your reasoning-tier models\u203a` and \u2039x\u203a, `\u2039model, effort,\ntokens and elapsed time\u203a`.\n", nil},
		// the input keeps each marker byte for byte, with its line ends (point 3 of the Operator's comment 6002406785)
		{"a marker kept with other line ends", strings.Replace(kept, "\u2039model, effort,\ntokens and elapsed time\u203a", "\u2039model, effort, tokens\nand elapsed time\u203a", 1),
			[]string{"docs/e.md:7: the input loses the marker \u2039model, effort, tokens and elapsed time\u203a"}},
	} {
		same(t, c.name, LostMarkers("docs/e.md", []byte(before), []byte(c.after)), c.want)
	}
}
