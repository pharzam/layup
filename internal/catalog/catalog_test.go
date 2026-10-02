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
	staticRow = "static\tactive\tgo\t1.26\ttest -z \"$(gofmt -l .)\"\t./*.go\t—\tfixtures/static.patch\thttps://go.dev/doc\n"
	layoutRow = "layout\tpending\tgo\t1.26\tgo test ./layout/\t./*.go\tlayout/layout_test.go layout/rules.txt\t—\t—\n"
)

// good is an entry that keeps every rule.
func good() fstest.MapFS {
	return entry(staticRow+layoutRow, map[string]string{
		"files/go.mod.tmpl":                      "module {{module}}\n\ngo 1.26\n",
		"files/.github/workflows/gates.yml.tmpl": "name: gates\n",
		"files/layout/rules.txt.tmpl":            "{{Module}} {module} {{ module }} {{module}}x\n",
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
	if k.State != "pending" || k.Fixture != "" || !slices.Equal(k.Config, []string{"layout/layout_test.go", "layout/rules.txt"}) || !slices.Equal(k.Scope, []string{"./*.go"}) {
		t.Fatalf("the layout kind %+v", k)
	}
}

func TestReadRefusesAnEntryThatBreaksARule(t *testing.T) {
	broken := func(change func(fstest.MapFS)) fstest.MapFS {
		fsys := good()
		change(fsys)
		return fsys
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
	if _, err := read(t, good()).Files(""); err == nil {
		t.Error("an empty module path: no error")
	}
}

func TestManifestIsTheKindsWithoutVersionFixtureAndEvidence(t *testing.T) {
	m, err := read(t, good()).Manifest()
	if err != nil {
		t.Fatal(err)
	}
	want := "kind\tstate\ttool\tcommand\tscope\tconfig\n" +
		"static\tactive\tgo\ttest -z \"$(gofmt -l .)\"\t./*.go\t—\n" +
		"layout\tpending\tgo\tgo test ./layout/\t./*.go\tlayout/layout_test.go layout/rules.txt\n"
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
	if c, ok := e.Config("layout"); !ok || !slices.Equal(c, []string{"layout/layout_test.go", "layout/rules.txt"}) {
		t.Errorf("Config(layout) = %q, %v", c, ok)
	}
	if c, ok := e.Config("static"); !ok || len(c) != 0 {
		t.Errorf("Config(static) = %q, %v; want none", c, ok)
	}
	if _, ok := e.Config("test"); ok {
		t.Error("Config(test): the entry has no such kind")
	}
}

func TestHasNamesTheFilesOfTheEntryByTheirPathInIt(t *testing.T) {
	e := read(t, good())
	for ref, want := range map[string]bool{
		"go/kinds.tsv": true, "go/files/go.mod.tmpl": true, "go/files/.github/workflows/gates.yml.tmpl": true,
		"go/fixtures/static.patch": true,
		"go/go.mod":                false, "go/files/go.mod": false, "go/files/": false, "go/fixtures/layout.patch": false,
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
