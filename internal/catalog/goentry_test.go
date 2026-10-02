package catalog

import (
	"regexp"
	"slices"
	"strings"
	"testing"
)

// The entries of the binary (D10 of #91): at least one, each read by the rules
// of Read, the Go entry among them, and never the test entry.
func TestTheEntriesOfTheBinary(t *testing.T) {
	stacks, err := Stacks(Embedded())
	if err != nil || len(stacks) == 0 {
		t.Fatalf("the binary embeds no entry: %q, %v", stacks, err)
	}
	for _, s := range stacks {
		if _, err := Read(Embedded(), s); err != nil {
			t.Errorf("the entry %s: %v", s, err)
		}
	}
	if !slices.Contains(stacks, "go") || slices.Contains(stacks, "test") {
		t.Errorf("the entries %q; want go, and not test", stacks)
	}
}

func goEntry(t *testing.T) *Entry {
	t.Helper()
	e, err := Read(Embedded(), "go")
	if err != nil {
		t.Fatal(err)
	}
	return e
}

// The five kinds of the Go entry (D1 to D3 of #91, K20): static and test are
// active, with the version and the documentation of go; layout, boundary and
// contract are pending, with — in each column but the scope.
func TestTheKindsOfTheGoEntry(t *testing.T) {
	const evidence = "https://pkg.go.dev/cmd/go@go1.26.0"
	want := []Kind{
		{Name: "static", State: "active", Tool: "go", Version: "1.26", Command: `out=$(gofmt -l .) && test -z "$out" && go vet ./...`,
			Scope: []string{"./*.go"}, Fixture: "fixtures/static.patch", Evidence: evidence},
		{Name: "layout", State: "pending", Scope: []string{"./*.go"}},
		{Name: "boundary", State: "pending", Scope: []string{"./*.go"}},
		{Name: "contract", State: "pending", Scope: []string{"./*.go"}},
		{Name: "test", State: "active", Tool: "go", Version: "1.26", Command: "go test -count=1 ./...",
			Scope: []string{"./*.go"}, Config: []string{"docs/gates/coverage-floor.txt"}, Fixture: "fixtures/test.patch", Evidence: evidence},
	}
	got := goEntry(t).Kinds()
	if len(got) != len(want) {
		t.Fatalf("%d kinds, want %d: %+v", len(got), len(want), got)
	}
	for i, w := range want {
		g := got[i]
		if g.Name != w.Name || g.State != w.State || g.Tool != w.Tool || g.Version != w.Version || g.Command != w.Command ||
			!slices.Equal(g.Scope, w.Scope) || !slices.Equal(g.Config, w.Config) || g.Fixture != w.Fixture || g.Evidence != w.Evidence {
			t.Errorf("kind %d:\n%+v\nwant\n%+v", i+1, g, w)
		}
	}
}

// The files of the Go entry: go.mod with the module path and the go line of
// LAYUP's own go.mod (D3), the workflow and the job script (D4, D5), and the
// file of the coverage floor, whose one line is its marker (D7, K24).
func TestTheFilesOfTheGoEntry(t *testing.T) {
	files, err := goEntry(t).Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, f := range files {
		paths = append(paths, f.Path)
		switch f.Path {
		case "go.mod":
			if string(f.Data) != "module example.com/target\n\ngo 1.26\n" {
				t.Errorf("go.mod: %q", f.Data)
			}
		case "docs/gates/coverage-floor.txt":
			if string(f.Data) != "\u2039the coverage floor of the test kind\u203a\n" {
				t.Errorf("docs/gates/coverage-floor.txt: %q", f.Data)
			}
		}
	}
	if !slices.Equal(paths, []string{".github/gates.sh", ".github/workflows/gates.yml", "docs/gates/coverage-floor.txt", "go.mod"}) {
		t.Fatalf("the files %q", paths)
	}
	want := []Gap{{Path: "docs/gates/coverage-floor.txt", Line: 1, Marker: "\u2039the coverage floor of the test kind\u203a",
		Question: "Which coverage floor, in percent of the statements, must the tests of the test kind reach, and what is its evidence?"}}
	if got := goEntry(t).Gaps(); !slices.Equal(got, want) {
		t.Errorf("the gaps %+v; want %+v", got, want)
	}
}

// file gives the text of the file at path of the Go entry, as the setup
// writes it.
func file(t *testing.T, path string) string {
	t.Helper()
	files, err := goEntry(t).Files("example.com/target")
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		if f.Path == path {
			return string(f.Data)
		}
	}
	t.Fatalf("the Go entry has no file %s", path)
	return ""
}

// The workflow of the Go entry (D4 of #91, K30, with notes 3 and 6 of its plan
// review): on pull_request, read-only; one job per kind, in the order of
// kinds.tsv, whose id and name are the kind; each job checks out the head of
// the pull request with its history, installs Go from go.mod, and runs the
// job script with its kind; no paths filter, no job-level if, no LAYUP program;
// each action at a version with its evidence. The test reads the text of the
// template; check jobs (row 15) owns the reader of the job names.
func TestTheWorkflowOfTheGoEntry(t *testing.T) {
	text := file(t, ".github/workflows/gates.yml")
	for _, want := range []string{"\non:\n  pull_request:\n", "\npermissions:\n  contents: read\n", "\njobs:\n"} {
		if strings.Count(text, want) != 1 {
			t.Errorf("the workflow holds %q %d times; want 1", want, strings.Count(text, want))
		}
	}
	kinds := goEntry(t).Kinds()
	_, jobs, _ := strings.Cut(text, "\njobs:\n")
	var ids []string
	for _, l := range strings.Split(jobs, "\n") {
		if regexp.MustCompile(`^  [a-z][a-z-]*:$`).MatchString(l) {
			ids = append(ids, strings.TrimSuffix(strings.TrimSpace(l), ":"))
		}
	}
	var names []string
	for _, k := range kinds {
		names = append(names, k.Name)
		job := "\n  " + k.Name + ":\n    name: " + k.Name + "\n    runs-on: ubuntu-latest\n"
		run := "      - run: sh .github/gates.sh " + k.Name + "\n"
		if strings.Count(text, job) != 1 || strings.Count(text, run) != 1 {
			t.Errorf("the job of the kind %s: %d heads and %d runs; want 1 of each", k.Name, strings.Count(text, job), strings.Count(text, run))
		}
	}
	if !slices.Equal(ids, names) {
		t.Errorf("the jobs %q; want one per kind, in its order: %q", ids, names)
	}
	for want, n := range map[string]int{
		"          ref: ${{ github.event.pull_request.head.sha }}\n":                        len(kinds),
		"          fetch-depth: 0\n":                                                        len(kinds),
		"          go-version-file: go.mod\n":                                               len(kinds),
		"          GATE_BASE: ${{ github.event.pull_request.base.sha }}\n":                  len(kinds),
		"          GATE_HEAD: ${{ github.event.pull_request.head.sha }}\n":                  len(kinds),
		"      - uses: actions/checkout@v4 # https://github.com/actions/checkout/tree/v4\n": len(kinds),
		"      - uses: actions/setup-go@v5 # https://github.com/actions/setup-go/tree/v5\n": len(kinds),
	} {
		if got := strings.Count(text, want); got != n {
			t.Errorf("the workflow holds %q %d times; want %d", want, got, n)
		}
	}
	uses := regexp.MustCompile(`^      - uses: ([a-z-]+/[a-z-]+)@(v[0-9]+) # https://github\.com/([a-z-]+/[a-z-]+)/tree/(v[0-9]+)$`)
	for i, l := range strings.Split(text, "\n") {
		if strings.Contains(l, "uses:") {
			m := uses.FindStringSubmatch(l)
			if m == nil || m[1] != m[3] || m[2] != m[4] {
				t.Errorf("line %d: %q is not an action at a version with its evidence", i+1, l)
			}
		}
		for _, bad := range []string{"paths:", "paths-ignore:", "layup", "setup-check", "&", "*", ": {", ": [", "\t"} {
			if strings.Contains(l, bad) {
				t.Errorf("line %d: %q holds %q", i+1, l, bad)
			}
		}
		if strings.HasPrefix(l, "    if:") {
			t.Errorf("line %d: %q is a job-level if", i+1, l)
		}
	}
}

// The job script of the Go entry starts no LAYUP program (D5 of #91,
// NFR-002): no line that is not a comment names it.
func TestTheJobScriptOfTheGoEntry(t *testing.T) {
	text := file(t, ".github/gates.sh")
	if !strings.HasPrefix(text, "#!/bin/sh\n") {
		t.Errorf("the script starts %q", text[:min(len(text), 20)])
	}
	for i, l := range strings.Split(text, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(l), "#") && strings.Contains(strings.ToLower(l), "layup") {
			t.Errorf("line %d: %q names layup", i+1, l)
		}
	}
}
