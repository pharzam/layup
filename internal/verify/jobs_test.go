package verify

import (
	"reflect"
	"testing"
	"testing/fstest"
)

// The job reader reads the check name of each job as check_protection does: a
// key at the indent of the first line under jobs: is a job, its name: at the
// indent of its first child line is its name, else its id (D2 of #92).
func TestJobNames(t *testing.T) {
	for _, c := range []struct {
		name, text string
		want       []string
	}{
		{"an id and a name", "on: push\njobs:\n  static:\n    name: static\n    runs-on: x\n  test:\n    runs-on: x\n", []string{"static", "test"}},
		{"a name of another job", "jobs:\n  lint:\n    name: static\n", []string{"static"}},
		{"quotes and comments", "jobs:\n  # the gates\n  \"a\": # x\n    name: 'static' # the kind\n\n  b:\n    name: \"test\"\n", []string{"static", "test"}},
		{"CR line ends", "jobs:\r\n  static:\r\n    name: static  \r\n    steps:\r\n      - name: other\r\n", []string{"static"}},
		{"a deeper name", "jobs:\n  a:\n    steps:\n      - name: static\n    name: b\n", []string{"b"}},
		{"a name of a deeper map", "jobs:\n  a:\n    with:\n      name: static\n    runs-on: x\n", []string{"a"}},
		{"four spaces", "jobs:\n    static:\n        name: x\n    test:\n        runs-on: y\n", []string{"x", "test"}},
		{"a key after jobs", "jobs:\n  a:\n    runs-on: x\nenv:\n  b:\n    name: c\n", []string{"a"}},
		{"no jobs", "name: gates\non: push\n", nil},
		{"jobs inside a line", "x: |\n  jobs:\n    a:\n", nil},
	} {
		if got := jobNames(c.text); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: %q; want %q", c.name, got, c.want)
		}
	}
}

// Check jobs gives a finding for each kind of the manifest with no job of a
// workflow .yml or .yaml (D2 of #92).
func TestCheckJobs(t *testing.T) {
	tree := fstest.MapFS{
		"docs/gates.tsv":              {Data: []byte(twoKinds)},
		".github/workflows/a.yml":     {Data: []byte("jobs:\n  static:\n    runs-on: x\n")},
		".github/workflows/b.yaml":    {Data: []byte("jobs:\n  x:\n    name: layout\n")},
		".github/workflows/c.json":    {Data: []byte("jobs:\n  y:\n    name: other\n")},
		".github/workflows/sub/d.yml": {Data: []byte("jobs:\n  z:\n    name: other\n")},
	}
	if got := checkJobs(input{fsys: tree}); got != nil {
		t.Errorf("a job per kind: %q; want no finding", got)
	}
	delete(tree, ".github/workflows/b.yaml")
	if got, want := checkJobs(input{fsys: tree}), []string{"no CI job for the kind layout"}; !reflect.DeepEqual(got, want) {
		t.Errorf("no job of layout: %q; want %q", got, want)
	}
	for _, c := range []struct {
		manifest string
		want     string
	}{
		{"", "docs/gates.tsv: no manifest"},
		{"kind\tstate\n", "docs/gates.tsv: line 1"},
	} {
		tree := fstest.MapFS{}
		if c.manifest != "" {
			tree["docs/gates.tsv"] = &fstest.MapFile{Data: []byte(c.manifest)}
		}
		if got := checkJobs(input{fsys: tree}); len(got) != 1 || len(got[0]) < len(c.want) || got[0][:len(c.want)] != c.want {
			t.Errorf("the manifest %q: %q; want one finding that starts with %q", c.manifest, got, c.want)
		}
	}
}
