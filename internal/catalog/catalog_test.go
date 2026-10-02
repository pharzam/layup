package catalog

import (
	"slices"
	"strings"
	"testing"
	"testing/fstest"
)

const header = "kind\tstate\ttool\tversion\tcommand\tscope\tconfig\tfixture\tevidence\n"

// entry gives a catalog with the stack go: kinds.tsv with the rows, and the
// other files by their path in the entry.
func entry(rows string, files map[string]string) fstest.MapFS {
	fsys := fstest.MapFS{"go/kinds.tsv": {Data: []byte(header + rows)}}
	for p, text := range files {
		fsys["go/"+p] = &fstest.MapFile{Data: []byte(text)}
	}
	return fsys
}

const (
	staticRow = "static\tactive\tgo\t1.26\ttest -z \"$(gofmt -l .)\"\t./*.go\t.golangci.yml docs/rules.txt\tfixtures/static.patch\thttps://go.dev/doc\n"
	// A pending kind has — in each column but its kind, its state and its
	// scope (D1 of #91, with note 1 of its plan review).
	layoutRow = "layout\tpending\t—\t—\t—\t./*.go\t—\t—\t—\n"
	gapsFile  = "path\tmarker\tquestion\ndocs/floor.txt\tthe floor\tWhich floor?\n"
)

// good is an entry that keeps every rule.
func good() fstest.MapFS {
	return entry(staticRow+layoutRow, map[string]string{
		"files/go.mod.tmpl":                      "module {{module}}\n\ngo 1.26\n",
		"files/.github/workflows/gates.yml.tmpl": "name: gates\n",
		"files/layout/rules.txt.tmpl":            "{{Module}} {module} {{ module }} {{module}}x\n",
		"files/docs/floor.txt.tmpl":              "{{gap:the floor}}\n",
		"gaps.tsv":                               gapsFile,
		"fixtures/static.patch":                  "--- a\n+++ b\n",
	})
}

func read(t *testing.T, fsys fstest.MapFS) *Entry {
	t.Helper()
	e, err := Read(fsys, "go")
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestReadGivesTheKindsInTheOrderOfTheFile(t *testing.T) {
	var names []string
	for _, k := range read(t, good()).Kinds() {
		names = append(names, k.Name)
	}
	if !slices.Equal(names, []string{"static", "layout"}) {
		t.Fatalf("the kinds %q, want static and layout", names)
	}
	k := read(t, good()).Kinds()[1]
	if k.State != "pending" || k.Tool != "" || k.Version != "" || k.Command != "" || len(k.Config) != 0 || k.Fixture != "" || k.Evidence != "" || !slices.Equal(k.Scope, []string{"./*.go"}) {
		t.Fatalf("the layout kind %+v", k)
	}
}

func TestReadRefusesAnEntryThatBreaksARule(t *testing.T) {
	broken := func(change func(fstest.MapFS)) fstest.MapFS {
		fsys := good()
		change(fsys)
		return fsys
	}
	pending := func(row string) fstest.MapFS {
		return broken(func(f fstest.MapFS) { f["go/kinds.tsv"] = &fstest.MapFile{Data: []byte(header + staticRow + row)} })
	}
	active := func(row string) fstest.MapFS {
		return broken(func(f fstest.MapFS) { f["go/kinds.tsv"] = &fstest.MapFile{Data: []byte(header + row + layoutRow)} })
	}
	for _, c := range []struct {
		name string
		fsys fstest.MapFS
		want string
	}{
		{"no kinds.tsv", broken(func(f fstest.MapFS) { delete(f, "go/kinds.tsv") }), "go/kinds.tsv"},
		{"a bad state", entry(strings.Replace(staticRow, "active", "on", 1), map[string]string{"fixtures/static.patch": "x"}), `go/kinds.tsv: line 2, column "state"`},
		{"a missing column", entry("static\tactive\n", nil), `go/kinds.tsv: line 2`},
		{"no kind", entry("", nil), "go/kinds.tsv: no kind"},
		{"an active kind with no fixture", entry(strings.Replace(staticRow, "fixtures/static.patch", "—", 1), nil), "the active kind static names the fixture —"},
		{"an active kind with another fixture path", broken(func(f fstest.MapFS) {
			f["go/kinds.tsv"] = &fstest.MapFile{Data: []byte(header + strings.Replace(staticRow, "fixtures/static.patch", "fixtures/other.patch", 1) + layoutRow)}
		}), "the active kind static names the fixture fixtures/other.patch"},
		{"a missing fixture file", broken(func(f fstest.MapFS) { delete(f, "go/fixtures/static.patch") }), "go/fixtures/static.patch: the fixture of the active kind static does not exist"},
		{"a pending kind with a fixture", broken(func(f fstest.MapFS) {
			f["go/kinds.tsv"] = &fstest.MapFile{Data: []byte(header + staticRow + strings.Replace(layoutRow, "\t—\t—\n", "\tfixtures/layout.patch\t—\n", 1))}
		}), "the pending kind layout names the fixture fixtures/layout.patch"},
		{"a fixture of no kind", broken(func(f fstest.MapFS) { f["go/fixtures/extra.patch"] = &fstest.MapFile{Data: []byte("x")} }), "go/fixtures/extra.patch: the file is the fixture of no active kind (no active kind extra)"},
		{"a file of files/ without the suffix", broken(func(f fstest.MapFS) { f["go/files/go.mod"] = &fstest.MapFile{Data: []byte("x")} }), "go/files/go.mod: a file of files/ ends with .tmpl"},
		{"a file of files/ that is only the suffix", broken(func(f fstest.MapFS) { f["go/files/.tmpl"] = &fstest.MapFile{Data: []byte("x")} }), "go/files/.tmpl: a file of files/ ends with .tmpl"},
		// D1 of #91, with note 1 of its plan review: a pending kind has — in
		// its tool, version, command, config and evidence; an active kind has
		// a version, and an https URL as its evidence.
		{"a pending kind with a tool", pending("layout\tpending\tgo\t—\t—\t./*.go\t—\t—\t—\n"), "the pending kind layout names the tool go; a pending kind has —"},
		{"a pending kind with a version", pending("layout\tpending\t—\t1.26\t—\t./*.go\t—\t—\t—\n"), "the pending kind layout names the version 1.26; a pending kind has —"},
		{"a pending kind with a command", pending("layout\tpending\t—\t—\tgo test ./...\t./*.go\t—\t—\t—\n"), "the pending kind layout names the command go test ./...; a pending kind has —"},
		{"a pending kind with a config", pending("layout\tpending\t—\t—\t—\t./*.go\tlayout.txt\t—\t—\n"), "the pending kind layout names the config layout.txt; a pending kind has —"},
		{"a pending kind with an evidence", pending("layout\tpending\t—\t—\t—\t./*.go\t—\t—\thttps://go.dev/doc\n"), "the pending kind layout names the evidence https://go.dev/doc; a pending kind has —"},
		{"an active kind with no version", active(strings.Replace(staticRow, "\t1.26\t", "\t—\t", 1)), "the active kind static has no version"},
		{"an active kind whose evidence is no URL", active(strings.Replace(staticRow, "https://go.dev/doc", "go.dev/doc", 1)), "the evidence of the active kind static is not an https URL: go.dev/doc"},
		{"an active kind with no evidence", active(strings.Replace(staticRow, "https://go.dev/doc", "—", 1)), "the evidence of the active kind static is not an https URL: —"},
		// D7 of #91: a gap of the entry is a token in a file of files/ with
		// its row in gaps.tsv, and no file of the entry holds a marker.
		{"a gap token with no row", broken(func(f fstest.MapFS) { delete(f, "go/gaps.tsv") }), "go/files/docs/floor.txt.tmpl: the gap token {{gap:the floor}} has no row in go/gaps.tsv"},
		{"a gap row with no token", broken(func(f fstest.MapFS) { f["go/files/docs/floor.txt.tmpl"] = &fstest.MapFile{Data: []byte("80\n")} }), "go/gaps.tsv: line 2: the file docs/floor.txt holds the gap token {{gap:the floor}} 0 times; want 1"},
		{"a gap token twice", broken(func(f fstest.MapFS) {
			f["go/files/docs/floor.txt.tmpl"] = &fstest.MapFile{Data: []byte("{{gap:the floor}}\n{{gap:the floor}}\n")}
		}), "go/gaps.tsv: line 2: the file docs/floor.txt holds the gap token {{gap:the floor}} 2 times; want 1"},
		{"a gap row whose path is no file", broken(func(f fstest.MapFS) {
			f["go/gaps.tsv"] = &fstest.MapFile{Data: []byte(gapsFile + "docs/none.txt\tnone\tWhich?\n")}
		}), "go/gaps.tsv: line 3: docs/none.txt is not a file of files/"},
		{"a gaps.tsv of another form", broken(func(f fstest.MapFS) {
			f["go/gaps.tsv"] = &fstest.MapFile{Data: []byte("path\tmarker\ndocs/floor.txt\tthe floor\n")}
		}), "go/gaps.tsv: line 1"},
		{"a gap with an empty question", broken(func(f fstest.MapFS) {
			f["go/gaps.tsv"] = &fstest.MapFile{Data: []byte("path\tmarker\tquestion\ndocs/floor.txt\tthe floor\t—\n")}
		}), "go/gaps.tsv: line 2: the gap the floor has no question"},
		{"a marker character in a file of files/", broken(func(f fstest.MapFS) {
			f["go/files/README.md.tmpl"] = &fstest.MapFile{Data: []byte("a \u2039x\u203a b\n")}
		}), "go/files/README.md.tmpl: a marker character; an entry writes a gap token in its place"},
		{"a marker character in a fixture", broken(func(f fstest.MapFS) { f["go/fixtures/static.patch"] = &fstest.MapFile{Data: []byte("+\u203a\n")} }), "go/fixtures/static.patch: a marker character; an entry writes a gap token in its place"},
	} {
		_, err := Read(c.fsys, "go")
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want one that holds %q", c.name, err, c.want)
		}
	}
}

func TestReadRefusesANameThatIsNotAnEntry(t *testing.T) {
	for _, stack := range []string{"", ".", "go/files", "../go", "rust"} {
		if _, err := Read(good(), stack); err == nil {
			t.Errorf("Read(%q): no error", stack)
		}
	}
}

func TestFilesReplaceTheModuleAndNothingElse(t *testing.T) {
	files, err := read(t, good()).Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	want := []File{
		{".github/workflows/gates.yml", []byte("name: gates\n")},
		{"docs/floor.txt", []byte("\u2039the floor\u203a\n")},
		{"go.mod", []byte("module example.com/target\n\ngo 1.26\n")},
		{"layout/rules.txt", []byte("{{Module}} {module} {{ module }} example.com/targetx\n")},
	}
	if len(files) != len(want) {
		t.Fatalf("%d files, want %d: %v", len(files), len(want), files)
	}
	for i := range want {
		if files[i].Path != want[i].Path || string(files[i].Data) != string(want[i].Data) {
			t.Errorf("file %d: %s %q, want %s %q", i, files[i].Path, files[i].Data, want[i].Path, want[i].Data)
		}
	}
	for _, module := range []string{"", "a b", "a\nb", "a\tb"} {
		if _, err := read(t, good()).Files(module); err == nil {
			t.Errorf("the module path %q: no error", module)
		}
	}
}

func TestManifestIsTheKindsWithoutVersionFixtureAndEvidence(t *testing.T) {
	m, err := read(t, good()).Manifest()
	if err != nil {
		t.Fatal(err)
	}
	want := "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
		"static\tactive\tgo\ttest -z \"$(gofmt -l .)\"\t./*.go\t.golangci.yml docs/rules.txt\n" +
		"layout\tpending\t—\t—\t./*.go\t—\n"
	if string(m) != want {
		t.Fatalf("the manifest:\n%s\nwant:\n%s", m, want)
	}
}

func TestFixtureAndConfigOfAKind(t *testing.T) {
	e := read(t, good())
	if p, err := e.Fixture("static"); err != nil || string(p) != "--- a\n+++ b\n" {
		t.Errorf("Fixture(static) = %q, %v", p, err)
	}
	if _, err := e.Fixture("layout"); err == nil {
		t.Error("Fixture(layout), a pending kind: no error")
	}
	if c, ok := e.Config("static"); !ok || !slices.Equal(c, []string{".golangci.yml", "docs/rules.txt"}) {
		t.Errorf("Config(static) = %q, %v", c, ok)
	}
	if c, ok := e.Config("layout"); !ok || len(c) != 0 {
		t.Errorf("Config(layout) = %q, %v; want none", c, ok)
	}
	if _, ok := e.Config("test"); ok {
		t.Error("Config(test): the entry has no such kind")
	}
}

func TestHasNamesTheFilesOfTheEntryByTheirPathInIt(t *testing.T) {
	e := read(t, good())
	for ref, want := range map[string]bool{
		"go/kinds.tsv": true, "go/files/go.mod.tmpl": true, "go/files/.github/workflows/gates.yml.tmpl": true,
		"go/fixtures/static.patch": true, "go/gaps.tsv": true,
		"go/go.mod": false, "go/files/go.mod": false, "go/files/": false, "go/fixtures/layout.patch": false,
		"rust/kinds.tsv": false, "go": false, "kinds.tsv": false,
	} {
		if got := e.Has(ref); got != want {
			t.Errorf("Has(%q) = %v, want %v", ref, got, want)
		}
	}
}

func TestStacksListsTheDirectoriesWithAKindsFile(t *testing.T) {
	fsys := fstest.MapFS{
		"rust/kinds.tsv":  {Data: []byte(header)},
		"go/kinds.tsv":    {Data: []byte(header)},
		"notes/README.md": {Data: []byte("x")},
		"kinds.tsv":       {Data: []byte(header)},
	}
	if s, err := Stacks(fsys); err != nil || !slices.Equal(s, []string{"go", "rust"}) {
		t.Fatalf("Stacks = %q, %v; want go and rust", s, err)
	}
}

// Gaps gives each gap token of the files with its path in the target, its
// line, its marker and its question, by path and line (D7 of #91).
func TestGapsGiveEachGapAtItsLine(t *testing.T) {
	fsys := good()
	fsys["go/files/docs/a.txt.tmpl"] = &fstest.MapFile{Data: []byte("one\ntwo {{module}}\n{{gap:a}} and {{gap:b}}\n")}
	fsys["go/gaps.tsv"] = &fstest.MapFile{Data: []byte(gapsFile + "docs/a.txt\tb\tQ b?\ndocs/a.txt\ta\tQ a?\n")}
	got := read(t, fsys).Gaps()
	want := []Gap{
		{Path: "docs/a.txt", Line: 3, Marker: "\u2039a\u203a", Question: "Q a?"},
		{Path: "docs/a.txt", Line: 3, Marker: "\u2039b\u203a", Question: "Q b?"},
		{Path: "docs/floor.txt", Line: 1, Marker: "\u2039the floor\u203a", Question: "Which floor?"},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Gaps =\n%+v\nwant\n%+v", got, want)
	}
	files, err := read(t, fsys).Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if f.Path == "docs/a.txt" && string(f.Data) != "one\ntwo example.com/target\n\u2039a\u203a and \u2039b\u203a\n" {
			t.Errorf("docs/a.txt: %q", f.Data)
		}
	}
	if len(read(t, entry(staticRow, map[string]string{"fixtures/static.patch": "x"})).Gaps()) != 0 {
		t.Error("an entry with no gaps.tsv has a gap")
	}
}
