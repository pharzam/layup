package verify

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/pharzam/layup/internal/work"
)

// The marker scanner and check markers (D1 and D2 of #87, with the answer to
// its plan review; fix 1 of the first pilot, #97): on a line, an open quote
// starts a marker to the first close quote after it; an open quote with no
// close quote after it is no marker, but an unpaired quote; an open quote
// between two backticks is the mention; the text of the convention is no
// marker.
func TestTheMarkersOfALine(t *testing.T) {
	for _, c := range []struct {
		line string
		want []string
	}{
		{"a \u2039x\u203a b", []string{"\u2039x\u203a"}},
		{"a \u2039x b", nil},
		{"the `\u2039` mention", nil},
		{"the convention `\u2039\u2026\u203a`", nil},
		{"\u2039a\u203a and \u2039b\u203a", []string{"\u2039a\u203a", "\u2039b\u203a"}},
		{"`\u2039`\u2039x\u203a", []string{"\u2039x\u203a"}},
		{"\u2039`x`\u203a", []string{"\u2039`x`\u203a"}},
		{"two lines: `\u2039State one", nil},
		{"a \u2039x\r", nil},
		{"\u2039\u2039x\u203a", []string{"\u2039\u2039x\u203a"}},
		{"no marker", nil},
	} {
		var got []string
		for _, m := range TextMarkers(c.line) {
			got = append(got, m.Text)
		}
		if !slices.Equal(got, c.want) {
			t.Errorf("%q: %q, want %q", c.line, got, c.want)
		}
	}
}

// A marker runs from its open quote to the first close quote after it, also
// across lines (fix 1 of the first pilot, #97: a marker of
// docs/tasks/backlog.md runs over four lines, and S10 asked and S11 filled its
// first line only). Its line and column are those of its open quote, and its
// text holds its line ends; its key, which a question, a row of the record
// and a row of open-gaps.tsv show, has each line end as one space.
func TestAMarkerOverMoreLines(t *testing.T) {
	scheme := "\u2039State your exact scheme\nhere \u2014 for example: \"T-\" plus four characters drawn from 0-9 a-z minus the\nambiguous i l o u; before using an ID, confirm tasks/<id>.md does not already\nexist.\"\u203a"
	lead := "branches, and PRs. A random suffix needs no coordination. `"
	text := "agents) working in parallel\n" + lead + scheme + "`\n\n\u2039x\u203a after it\n"
	want := []Marker{{"", 2, len(lead), scheme}, {"", 7, 0, "\u2039x\u203a"}}
	if got := TextMarkers(text); !slices.Equal(got, want) {
		t.Errorf("TextMarkers\n got %+v\nwant %+v", got, want)
	}
	for _, c := range []struct{ text, key string }{
		{scheme, strings.ReplaceAll(scheme, "\n", " ")},
		{"\u2039one\r\ntwo\u203a", "\u2039one two\u203a"},
		{"\u2039x\u203a", "\u2039x\u203a"},
	} {
		if got := work.MarkerKey(c.text); got != c.key {
			t.Errorf("the key of %q: %q, want %q", c.text, got, c.key)
		}
	}
}

// The scanner reads each tracked file outside MK_EXEMPT and gives each
// occurrence with its file and line, so S10 and check sources read one list.
func TestScanMarkers(t *testing.T) {
	fsys := fstest.MapFS{
		"docs/a.md":                {Data: []byte("one \u2039x\u203a\n\ntwo \u2039x\u203a and \u2039y\u203a\n")},
		"docs/\u00e9.md":           {Data: []byte("\u2039q\u203a\n")},
		"docs/q\x7fx.md":           {Data: []byte("\u2039port\u203a\n")},
		"docs/templates/t.md":      {Data: []byte("\u2039t\u203a\n")},
		"docs/adr/0003-x.md":       {Data: []byte("\u2039adr\u203a\n")},
		"docs/setup/open-gaps.tsv": {Data: []byte("docs/a.md\t\u2039x\u203a\tq\n")},
		"docs/dir.md/f":            {Data: []byte("\u2039in a dir\u203a\n")},
	}
	files := []string{"docs/a.md", "docs/\u00e9.md", "docs/q\x7fx.md", "docs/templates/t.md", "docs/adr/0003-x.md", "docs/setup/open-gaps.tsv", "docs/dir.md", "gone.md"}
	want := []Marker{{"docs/a.md", 1, 4, "\u2039x\u203a"}, {"docs/a.md", 3, 4, "\u2039x\u203a"}, {"docs/a.md", 3, 16, "\u2039y\u203a"},
		{"docs/\u00e9.md", 1, 0, "\u2039q\u203a"}, {"docs/q\x7fx.md", 1, 0, "\u2039port\u203a"}}
	if got := scanMarkers(fsys, files); !slices.Equal(got, want) {
		t.Errorf("the markers\n got %+v\nwant %+v", got, want)
	}
}

// The markers of a text, each with its line and the byte column of its open
// quote, so S11 replaces a marker at the place that the scanner found and a
// mention in a code span stays (D4 and D7 of #90).
func TestTheMarkersOfAText(t *testing.T) {
	text := "a \u2039x\u203a\n`\u2039`\u2039y\u203a and `\u2039y\u203a`\n\u2039open to the end\n"
	want := []Marker{{"", 1, 2, "\u2039x\u203a"}, {"", 2, 5, "\u2039y\u203a"}, {"", 2, 18, "\u2039y\u203a"}}
	if got := TextMarkers(text); !slices.Equal(got, want) {
		t.Errorf("TextMarkers\n got %+v\nwant %+v", got, want)
	}
	for _, m := range want {
		line := strings.Split(text, "\n")[m.Line-1]
		if !strings.HasPrefix(line[m.Col:], m.Text) {
			t.Errorf("%+v: the line %q holds no marker at its column", m, line)
		}
	}
}

// Check markers: each marker needs its row in docs/setup/open-gaps.tsv, each
// row its marker, and each row a question; the file is read by tabs, as
// check_markers reads it, and an absent file is no rows.
func TestTheOpenGaps(t *testing.T) {
	const gaps = "docs/setup/open-gaps.tsv"
	marked := "\u2039x\u203a here\nand \u2039x\u203a again\n"
	for _, c := range []struct {
		name  string
		files map[string]string
		want  []string
	}{
		{"each marker listed", map[string]string{"docs/a.md": marked, gaps: "docs/a.md\t\u2039x\u203a\tWhich x?\n"}, nil},
		{"no marker and no file", map[string]string{"docs/a.md": "none\n"}, nil},
		{"an unlisted marker, once", map[string]string{"docs/a.md": marked}, []string{"unlisted: docs/a.md \u2039x\u203a"}},
		{"a stale row", map[string]string{"docs/a.md": "none\n", gaps: "docs/y.md\t\u2039bar\u203a\tWho?\n"},
			[]string{"stale: docs/y.md \u2039bar\u203a is listed in docs/setup/open-gaps.tsv but does not occur"}},
		{"a row with no question", map[string]string{"docs/a.md": marked, gaps: "docs/a.md\t\u2039x\u203a\n"},
			[]string{"question: docs/setup/open-gaps.tsv line 1 has no question"}},
		{"blank questions and the empty mark", map[string]string{"docs/a.md": marked + "\u2039b\u203a \u2039c\u203a \u2039d\u203a\n",
			gaps: "docs/a.md\t\u2039x\u203a\t \n\ndocs/a.md\t\u2039b\u203a\t\t\ndocs/a.md\t\u2039c\u203a\t\u2014\ndocs/a.md\t\u2039d\u203a\t\r\n"},
			[]string{"question: docs/setup/open-gaps.tsv line 1 has no question", "question: docs/setup/open-gaps.tsv line 3 has no question",
				"question: docs/setup/open-gaps.tsv line 4 has no question", "question: docs/setup/open-gaps.tsv line 5 has no question"}},
		{"the order: questions, unlisted, stale", map[string]string{"docs/b.md": "\u2039b\u203a\n", "docs/a.md": "\u2039a\u203a\n",
			gaps: "docs/z.md\t\u2039z\u203a\tq\ndocs/y.md\t\u2039y\u203a\n"},
			[]string{"question: docs/setup/open-gaps.tsv line 2 has no question", "unlisted: docs/a.md \u2039a\u203a", "unlisted: docs/b.md \u2039b\u203a",
				"stale: docs/y.md \u2039y\u203a is listed in docs/setup/open-gaps.tsv but does not occur",
				"stale: docs/z.md \u2039z\u203a is listed in docs/setup/open-gaps.tsv but does not occur"}},
		// a marker over two lines is listed by its key: each line end as one space (fix 1 of the first pilot, #97)
		{"a marker over two lines, listed by its key", map[string]string{"docs/a.md": "a `\u2039State one\nthing\u203a` b\n", gaps: "docs/a.md\t\u2039State one thing\u203a\tWhat?\n"}, nil},
		{"a marker over two lines, listed by its first line", map[string]string{"docs/a.md": "a `\u2039State one\nthing\u203a` b\n", gaps: "docs/a.md\t\u2039State one\tWhat?\n"},
			[]string{"unlisted: docs/a.md \u2039State one thing\u203a", "stale: docs/a.md \u2039State one is listed in docs/setup/open-gaps.tsv but does not occur"}},
	} {
		fsys := fstest.MapFS{}
		var paths []string
		for p, s := range c.files {
			fsys[p] = &fstest.MapFile{Data: []byte(s)}
			paths = append(paths, p)
		}
		slices.Sort(paths)
		same(t, c.name, markersFindings(fsys, stubHistory{files: paths}), c.want)
	}
	same(t, "a failed list", markersFindings(fstest.MapFS{}, stubHistory{err: errors.New("fatal")}), []string{"git: cannot list the tracked files"})
}

// Check markers names each angle quote with no pair, after the other findings
// (fix 3 of the first pilot, #97: a fill that lost its open quote left a
// close quote behind): an open quote with no close quote after it, and a close
// quote that no open quote before it takes. A quote between two backticks is
// the mention, and a file of MK_EXEMPT is not read.
func TestUnpairedQuotes(t *testing.T) {
	const gaps = "docs/setup/open-gaps.tsv"
	for _, c := range []struct {
		name  string
		files map[string]string
		want  []string
	}{
		{"the mentions and the convention", map[string]string{"docs/a.md": "the `\u2039` and `\u203a` quotes, and `\u2039\u2026\u203a`\n"}, nil},
		{"an open quote inside a marker", map[string]string{"docs/a.md": "\u2039\u2039x\u203a\n", gaps: "docs/a.md\t\u2039\u2039x\u203a\tWhich?\n"}, nil},
		{"an open quote with no close quote", map[string]string{"docs/a.md": "one\nthe \u2039rest of the file\nend\n"}, []string{"unpaired: docs/a.md:2 \u2039"}},
		{"a close quote with no open quote", map[string]string{"docs/a.md": "x\ntokens and elapsed time\u203a`; write\n"}, []string{"unpaired: docs/a.md:2 \u203a"}},
		{"a close quote after a marker", map[string]string{"docs/a.md": "\u2039x\u203a\u203a\n", gaps: "docs/a.md\t\u2039x\u203a\tWhich x?\n"}, []string{"unpaired: docs/a.md:1 \u203a"}},
		{"after the other findings", map[string]string{"docs/b.md": "\u203a\n", "docs/a.md": "\u2039a\u203a\n"},
			[]string{"unlisted: docs/a.md \u2039a\u203a", "unpaired: docs/b.md:1 \u203a"}},
		{"an exempt file", map[string]string{"docs/templates/t.md": "\u2039open\n"}, nil},
	} {
		fsys := fstest.MapFS{}
		var paths []string
		for p, s := range c.files {
			fsys[p] = &fstest.MapFile{Data: []byte(s)}
			paths = append(paths, p)
		}
		slices.Sort(paths)
		same(t, c.name, markersFindings(fsys, stubHistory{files: paths}), c.want)
	}
}
