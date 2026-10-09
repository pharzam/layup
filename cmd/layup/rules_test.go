package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// goodTable is a table of phase 1, a table of M2a and a table of M2b in the
// cell form of docs/spec/packages.md, with the line of rule 5 and the line of
// the engine checks; ' stands for a backtick.
var goodTable = strings.ReplaceAll(`# The packages

5. No package depends on 'net', 'net/http' or 'crypto/tls' unless its row says so.
   The one package that imports them: 'internal/forge/github'.

### The table of phase 1

The test of the package rules reads this table.

| Package | Job | May import | Starts a program |
| ------- | --- | ---------- | ---------------- |
| 'cmd/layup' | main | 'internal/cli' | no |
| 'internal/cli' | the command frame | 'internal/psb', 'internal/gate' | no |
| 'internal/psb' | the rules G1 to G5 | — | no |
| 'internal/git' | the one caller of git | — | 'git' |
| 'internal/gate' | the gates | 'internal/git' | 'sh -c': the gate commands |
| 'internal/setup' | not built yet | 'internal/git' | no |

The text after the table.

### The table of M2a

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| 'internal/forge/github' | the adapter | — | no | 'net', 'net/http', 'crypto/tls' |
| 'cmd/layup' | (the row of phase 1) | (the row of phase 1) | (the row of phase 1) | 'net/http' |
| 'internal/cli' | (the row of phase 1) | (the row of phase 1, with the change below) | (the row of phase 1) | 'net/http' |

The text after the table of M2a.

### The table of M2b

| Package | Job | May import | Starts a program | Connects |
| ------- | --- | ---------- | ---------------- | -------- |
| 'internal/session' | a role session | 'internal/git' | 'harness' (the command of a harness register row) | — |

The packages of the engine checks: 'internal/psb', 'internal/gate'.
`, "'", "`")

// The packages of goodModule, by their index.
const iCmd, iCLI, iPSB, iGate = 0, 1, 2, 4

// engine is the line of the engine checks of goodTable.
var engine = []string{"internal/psb", "internal/gate"}

// goodModule is a module that keeps every rule of goodTable.
func goodModule() module {
	pkg := func(rel, name string, imports ...string) goPackage {
		return goPackage{ImportPath: modulePath + "/" + rel, Name: name, Imports: imports, Deps: append([]string{}, imports...)}
	}
	return module{
		path: modulePath,
		packages: []goPackage{
			pkg("cmd/layup", "main", modulePath+"/internal/cli", "os"),
			pkg("internal/cli", "cli", modulePath+"/internal/psb", modulePath+"/internal/gate", "fmt"),
			pkg("internal/psb", "psb", "strings"),
			pkg("internal/git", "git", "os/exec"),
			pkg("internal/gate", "gate", modulePath+"/internal/git", "os/exec"),
			pkg("internal/session", "session", modulePath+"/internal/git", "os/exec"),
			{ImportPath: "fmt", Standard: true}, {ImportPath: "os", Standard: true},
			{ImportPath: "os/exec", Standard: true}, {ImportPath: "strings", Standard: true},
		},
		sources: map[string]map[string]string{
			"cmd/layup":    {"main.go": "package main\n\nfunc main() {}\n"},
			"internal/git": {"git.go": "package git\n\nimport \"os/exec\"\n\nvar c = exec.Command(\"git\", \"--version\")\n"},
			"internal/gate": {"gate.go": "package gate\n\nimport x \"os/exec\"\n\nvar c = x.CommandContext(nil, \"sh\", \"-c\", \"go vet\")\n\n" +
				"func run(c *x.Cmd) error { return c.Run() }\n"},
			"internal/session": {"session.go": "package session\n\nimport \"os/exec\"\n\nfunc start(command []string) *exec.Cmd { return exec.Command(command[0], command[1:]...) }\n"},
		},
	}
}

func TestReadTable(t *testing.T) {
	rows, err := readTable(goodTable)
	want := map[string]row{
		"cmd/layup":             {imports: []string{"internal/cli"}, connects: []string{"net/http"}},
		"internal/cli":          {imports: []string{"internal/psb", "internal/gate"}, connects: []string{"net/http"}},
		"internal/psb":          {},
		"internal/git":          {program: "git"},
		"internal/gate":         {imports: []string{"internal/git"}, program: "sh"},
		"internal/setup":        {imports: []string{"internal/git"}},
		"internal/forge/github": {connects: []string{"net", "net/http", "crypto/tls"}},
		"internal/session":      {imports: []string{"internal/git"}, register: true},
	}
	if err != nil || !reflect.DeepEqual(rows, want) {
		t.Fatalf("rows %+v, %v\nwant %+v", rows, err, want)
	}
}

// The reader refuses, and so checks nothing on, a table that it cannot read.
func TestReadTableRefusesWhatItCannotRead(t *testing.T) {
	under := "### The table of phase 1\n\n| Package | Job | May import | Starts a program |\n"
	head := under + "| - | - | - | - |\n"
	for name, text := range map[string]string{
		"no heading":                                    strings.Replace(goodTable, "### The table of phase 1", "### The table", 1),
		"no table in its section":                       "### The table of phase 1\n\nText.\n\n## Next\n\n" + strings.TrimPrefix(head, "### The table of phase 1\n\n") + "| `cmd/layup` | main | — | no |\n",
		"zero rows":                                     head + "\nText.\n",
		"no separator row":                              under + "| `cmd/layup` | main | — | no |\n| `internal/psb` | rules | — | no |\n",
		"a missing column":                              "### The table of phase 1\n\n| Package | Job | May import |\n| - | - | - |\n| `cmd/layup` | main | — |\n",
		"a short row":                                   head + "| `cmd/layup` | main | — |\n",
		"a package with no span":                        head + "| cmd/layup | main | — | no |\n",
		"two packages in a cell":                        head + "| `cmd/layup`, `cmd/x` | main | — | no |\n",
		"an import list in prose":                       head + "| `cmd/layup` | main | `internal/cli` and `internal/psb` | no |\n",
		"an empty import cell":                          head + "| `cmd/layup` | main |  | no |\n",
		"a program cell in prose":                       head + "| `internal/gate` | gates | — | the gate commands, with `sh -c` |\n",
		"a program cell with no span":                   head + "| `internal/gate` | gates | — | yes |\n",
		"two programs in a cell":                        head + "| `internal/gate` | gates | — | `sh`, or `bash` |\n",
		"a package in two rows":                         head + "| `cmd/layup` | main | — | no |\n| `cmd/layup` | main | — | no |\n",
		"no table of M2a":                               strings.Replace(goodTable, "### The table of M2a", "### The other table", 1),
		"no column Connects":                            strings.Replace(goodTable, "| Starts a program | Connects |", "| Starts a program | Other |", 1),
		"a Connects cell in prose":                      strings.Replace(goodTable, "| no | `net`, `net/http`, `crypto/tls` |", "| no | `net` and `net/http` |", 1),
		"a row of phase 1 that phase 1 lacks":           strings.Replace(goodTable, "| `cmd/layup` | (the row of phase 1) |", "| `internal/none` | (the row of phase 1) |", 1),
		"a row that mixes the two kinds of cell":        strings.Replace(goodTable, "| `cmd/layup` | (the row of phase 1) | (the row of phase 1) |", "| `cmd/layup` | main | (the row of phase 1) |", 1),
		"a full row of M2a for a package of phase 1":    strings.Replace(goodTable, "| `cmd/layup` | (the row of phase 1) | (the row of phase 1) | (the row of phase 1) |", "| `cmd/layup` | main | `internal/cli` | no |", 1),
		"no table of M2b":                               strings.Replace(goodTable, "### The table of M2b", "### The table of M2c", 1),
		"no column Connects in the table of M2b":        strings.Replace(goodTable, "| Package | Job | May import | Starts a program | Connects |\n| ------- | --- | ---------- | ---------------- | -------- |\n| `internal/session`", "| Package | Job | May import | Starts a program | Other |\n| ------- | --- | ---------- | ---------------- | -------- |\n| `internal/session`", 1),
		"a row of M2b that stands for a row of phase 1": strings.Replace(goodTable, "| `internal/session` | a role session |", "| `internal/session` | (the row of phase 1) |", 1),
		"a row of M2b for a package of another table":   strings.Replace(goodTable, "| `internal/session` | a role session |", "| `internal/psb` | a role session |", 1),
		"a register cell and a second span":             strings.Replace(goodTable, "(the command of a harness register row) |", "(the command of a harness register row), or `sh` |", 1),
		"a register cell with more text":                strings.Replace(goodTable, "(the command of a harness register row) |", "(the command of a harness register row) and more |", 1),
	} {
		if rows, err := readTable(text); err == nil {
			t.Errorf("%s: no error, %d rows; want an error", name, len(rows))
		}
	}
}

// The line of rule 5 names the one package that imports a network package
// (docs/spec/packages.md, The test of the package rules).
func TestReadAdapter(t *testing.T) {
	if a, err := readAdapter(goodTable); err != nil || a != "internal/forge/github" {
		t.Fatalf("readAdapter = %q, %v; want internal/forge/github", a, err)
	}
	line := "The one package that imports them: `internal/forge/github`."
	for name, text := range map[string]string{
		"no line":        strings.Replace(goodTable, line, "No line here.", 1),
		"two lines":      strings.Replace(goodTable, line, line+"\n   "+line, 1),
		"no code span":   strings.Replace(goodTable, line, "The one package that imports them: internal/forge/github.", 1),
		"two code spans": strings.Replace(goodTable, line, "The one package that imports them: `internal/forge/github`, `internal/x`.", 1),
	} {
		if a, err := readAdapter(text); err == nil {
			t.Errorf("%s: %q, no error; want an error", name, a)
		}
	}
}

// The line of the engine checks names their packages (docs/spec/packages.md,
// The test of the package rules; docs/spec/session.md, NFR-005).
func TestReadEngine(t *testing.T) {
	if e, err := readEngine(goodTable); err != nil || !slices.Equal(e, []string{"internal/psb", "internal/gate"}) {
		t.Fatalf("readEngine = %q, %v; want internal/psb, internal/gate", e, err)
	}
	line := "The packages of the engine checks: `internal/psb`, `internal/gate`."
	for name, text := range map[string]string{
		"no line":      strings.Replace(goodTable, line, "No line here.", 1),
		"two lines":    strings.Replace(goodTable, line, line+"\n"+line, 1),
		"no code span": strings.Replace(goodTable, line, "The packages of the engine checks: internal/psb.", 1),
	} {
		if e, err := readEngine(text); err == nil {
			t.Errorf("%s: %q, no error; want an error", name, e)
		}
	}
}

func TestAGoodModuleKeepsTheRules(t *testing.T) {
	rows, err := readTable(goodTable)
	if err != nil {
		t.Fatal(err)
	}
	if f := checkRules(rows, "internal/forge/github", engine, goodModule()); len(f) != 0 {
		t.Fatalf("findings on a good module:\n%s", strings.Join(f, "\n"))
	}
}

// One breach of each rule and each column; each gives exactly its findings,
// sorted, one per line.
func TestEachRuleAndColumnFindsItsBreach(t *testing.T) {
	// source puts a file in a package: its imports on line 3, its code on line 5.
	source := func(rel, imports, code string) func(*module) {
		return func(m *module) {
			m.sources[rel] = map[string]string{"x.go": fmt.Sprintf("package x\n\nimport %s\n\n%s\n", imports, code)}
		}
	}
	add := func(p goPackage) func(*module) { return func(m *module) { m.packages = append(m.packages, p) } }
	withRow := strings.Replace(goodTable, "\n\nThe text after", "\n| `internal/tool` | a tool | — | no |\n| `pkg/util` | helpers | — | no |\n\nThe text after", 1)
	bad := "package x\n\nfunc (\n"
	_, perr := parser.ParseFile(token.NewFileSet(), "x.go", bad, parser.SkipObjectResolution)
	for _, c := range []struct {
		name, table string
		breach      func(*module)
		want        string
	}{
		{"rule 1, the module path", goodTable, func(m *module) {
			*m = module{path: "example.com/other", packages: []goPackage{{ImportPath: "example.com/other/cmd/layup", Name: "main"}}}
		}, "rule 1: the module is example.com/other, want github.com/pharzam/layup"},
		{"rule 1, a second binary", withRow, add(goPackage{ImportPath: modulePath + "/internal/tool", Name: "main"}),
			"rule 1: internal/tool is a second binary"},
		{"rule 1, a package outside internal/", withRow, add(goPackage{ImportPath: modulePath + "/pkg/util", Name: "util"}),
			"rule 1: pkg/util is not under internal/"},
		{"rule 1, no binary", goodTable, func(m *module) { m.packages = m.packages[iCmd+1:] }, "rule 1: no binary cmd/layup"},
		{"rule 2, a require line", goodTable, func(m *module) { m.requires = []string{"example.com/lib v1.0.0"} },
			"rule 2: go.mod requires example.com/lib v1.0.0"},
		{"rule 2, a package of another module", goodTable, add(goPackage{ImportPath: "example.com/lib"}),
			"rule 2: example.com/lib is in neither the standard library nor this module"},
		{"rule 3", goodTable, source("internal/gate", `"os/exec"`, `var c = exec.Command("git", "status")`),
			"rule 3: internal/gate/x.go:5 starts git outside internal/git"},
		{"rule 4", goodTable, func(m *module) { m.packages[iCLI].Imports = append(m.packages[iCLI].Imports, "os/exec") },
			"rule 4: internal/cli imports os/exec, and its row starts no program"},
		{"rule 5, a dependency that Connects does not name", goodTable, func(m *module) { m.packages[iCmd].Deps = append(m.packages[iCmd].Deps, "net") },
			"rule 5: cmd/layup depends on net"},
		{"rule 5, a dependency with no Connects", goodTable, func(m *module) { m.packages[iPSB].Deps = append(m.packages[iPSB].Deps, "net/http") },
			"engine checks: internal/psb depends on net/http\nrule 5: internal/psb depends on net/http"},
		{"rule 5, a dependency that Connects names", goodTable, func(m *module) { m.packages[iCmd].Deps = append(m.packages[iCmd].Deps, "net/http") }, ""},
		{"rule 5, an own import outside the adapter", goodTable, func(m *module) {
			m.packages[iCLI].Imports = append(m.packages[iCLI].Imports, "net/http")
			m.packages[iCLI].Deps = append(m.packages[iCLI].Deps, "net/http")
		}, "rule 5: internal/cli imports net/http; only internal/forge/github imports them"},
		{"rule 5, the adapter imports it", goodTable, add(goPackage{ImportPath: modulePath + "/internal/forge/github", Name: "github",
			Imports: []string{"net/http"}, Deps: []string{"crypto/tls", "net", "net/http"}}), ""},
		{"May import", goodTable, func(m *module) {
			m.packages[iPSB].Imports = append(m.packages[iPSB].Imports, modulePath+"/internal/git")
		},
			"May import: internal/psb imports internal/git, which its row does not allow"},
		{"a package with no row", goodTable, add(goPackage{ImportPath: modulePath + "/internal/extra", Name: "extra"}),
			"table: internal/extra has no row"},
		{"Starts a program, another program", goodTable, source("internal/gate", `"os/exec"`, `var c = exec.Command("make")`),
			"Starts a program: internal/gate/x.go:5 starts make; its row says sh"},
		{"Starts a program, a row that says no", goodTable, source("cmd/layup", `e "os/exec"`, `func main() { e.Command("sh").Run() }`),
			"Starts a program: cmd/layup/x.go:5 starts sh; its row says no\nrule 4: cmd/layup imports os/exec, and its row starts no program"},
		{"Starts a program, not a string literal", goodTable, source("internal/git", `"os/exec"`, `const name = "git"; var c = exec.Command(name)`),
			"Starts a program: internal/git/x.go:5 starts a program that the scan cannot read"},
		{"Starts a program, exec.Command as a value", goodTable, source("internal/git", `"os/exec"`, `var start = exec.Command`),
			"Starts a program: internal/git/x.go:5 starts a program that the scan cannot read"},
		{"Starts a program, os.StartProcess", goodTable, source("internal/git", `"os"`, `var _, _ = os.StartProcess("/usr/bin/git", nil, nil)`),
			"Starts a program: internal/git/x.go:5 uses os.StartProcess"},
		{"Starts a program, syscall.Exec", goodTable, source("internal/psb", `"syscall"`, `var _ = syscall.Exec("/bin/sh", nil, nil)`),
			"Starts a program: internal/psb/x.go:5 uses syscall.Exec"},
		{"Starts a program, a dot import", goodTable, source("internal/gate", `. "os/exec"`, `var c = Command("sh")`),
			"Starts a program: internal/gate/x.go:3 imports os/exec with a dot, so the scan cannot read its calls"},
		{"Starts a program, syscall.ForkExec", goodTable, source("internal/psb", `"syscall"`, `var _, _ = syscall.ForkExec("/usr/bin/git", nil, nil)`),
			"Starts a program: internal/psb/x.go:5 uses syscall.ForkExec"},
		{"Starts a program, syscall.StartProcess", goodTable, source("internal/psb", `"syscall"`, `var _, _, _ = syscall.StartProcess("/bin/sh", nil, nil)`),
			"Starts a program: internal/psb/x.go:5 uses syscall.StartProcess"},
		{"Starts a program, an exec.Cmd made in place", goodTable, source("internal/gate", `"os/exec"`, `var _ = (&exec.Cmd{Path: "/usr/bin/git"}).Run()`),
			"Starts a program: internal/gate/x.go:5 starts a program that the scan cannot read"},
		{"a file that does not parse", goodTable, func(m *module) { m.sources["internal/psb"] = map[string]string{"x.go": bad} },
			"Starts a program: internal/psb/x.go does not parse: " + fmt.Sprint(perr)},
		// A file behind a build constraint: go list gives none of its imports.
		{"rule 4, a file that go list leaves out", goodTable, source("internal/cli", `"os/exec"`, `var _ = exec.ErrNotFound`),
			"rule 4: internal/cli imports os/exec, and its row starts no program"},
		{"rule 5, a file that go list leaves out", goodTable, func(m *module) {
			source("internal/psb", `"expvar"`, `var _ = expvar.NewInt`)(m)
			m.packages = append(m.packages, goPackage{ImportPath: "expvar", Standard: true, Deps: []string{"net/http"}})
		}, "engine checks: internal/psb depends on net/http\nrule 5: internal/psb depends on net/http"},
		{"May import, a file that go list leaves out", goodTable, source("internal/psb", `"`+modulePath+`/internal/git"`, `var _ = git.Version`),
			"May import: internal/psb imports internal/git, which its row does not allow"},
		// The rule of the engine checks reads each dependency, and no cell.
		{"engine checks, a dependency on internal/session", goodTable, func(m *module) {
			m.packages[iGate].Deps = append(m.packages[iGate].Deps, modulePath+"/internal/session")
		}, "engine checks: internal/gate depends on internal/session"},
		{"engine checks, a dependency on a package of rule 5", goodTable, func(m *module) { m.packages[iPSB].Deps = append(m.packages[iPSB].Deps, "net") },
			"engine checks: internal/psb depends on net\nrule 5: internal/cli depends on net\nrule 5: internal/psb depends on net"},
		{"engine checks, a Connects cell does not lift the rule", strings.Replace(goodTable, "| `cmd/layup` | (the row of phase 1) |", "| `internal/psb` | (the row of phase 1) | (the row of phase 1) | (the row of phase 1) | `net` |\n| `cmd/layup` | (the row of phase 1) |", 1),
			func(m *module) { m.packages[iPSB].Deps = append(m.packages[iPSB].Deps, "net") }, "engine checks: internal/psb depends on net\nrule 5: internal/cli depends on net"},
		{"engine checks, a package off the line", goodTable, func(m *module) {
			m.packages[iCLI].Deps = append(m.packages[iCLI].Deps, modulePath+"/internal/session")
		}, ""},
		// A register row starts the command that a register row names.
		{"a register row, a start by a string literal", goodTable, source("internal/session", `"os/exec"`, `var c = exec.Command("claude")`),
			"Starts a program: internal/session/x.go:5 starts claude; a register row starts only the command of its register"},
		{"a register row, exec.Command as a value", goodTable, source("internal/session", `"os/exec"`, `var start = exec.Command`),
			"Starts a program: internal/session/x.go:5 starts a program that the scan cannot read"},
		{"a register row, os.StartProcess", goodTable, source("internal/session", `"os"`, `var _, _ = os.StartProcess("/usr/bin/claude", nil, nil)`),
			"Starts a program: internal/session/x.go:5 uses os.StartProcess"},
	} {
		t.Run(c.name, func(t *testing.T) {
			rows, err := readTable(c.table)
			if err != nil {
				t.Fatal(err)
			}
			m := goodModule()
			c.breach(&m)
			if f := checkRules(rows, "internal/forge/github", engine, m); strings.Join(f, "\n") != c.want {
				t.Errorf("findings\n%s\nwant exactly\n%s", strings.Join(f, "\n"), c.want)
			}
		})
	}
}

func TestCheckInputsFindsEachReadOfTheEnvironmentOrTheStandardInput(t *testing.T) {
	m := module{sources: map[string]map[string]string{
		"internal/cli": {
			"a.go": "package cli\n\nimport \"os\"\n\nfunc f() string { return os.Getenv(\"X\") }\n",
			"b.go": "package cli\n\nimport o \"os\"\n\nfunc g() { _, _ = o.LookupEnv(\"Y\"); _ = o.Stdin }\n",
		},
		"internal/psb":  {"c.go": "package psb\n\nimport \"syscall\"\n\nvar v, _ = syscall.Getenv(\"Z\")\n"},
		"internal/gate": {"d.go": "package gate\n\nimport . \"os\"\n\nfunc h() []string { return Environ() }\n"},
		"internal/git": {"git.go": "package git\n\nimport \"os\"\n\nfunc environ() { os.LookupEnv(\"PATH\") }\n\n" +
			"func other() string { return os.Getenv(\"HOME\") }\n"},
		"internal/tsv": {"e.go": "package tsv\n\nimport \"os\"\n\nfunc environ() []string { return os.Environ() }\n"},
		"internal/session": {"session.go": "package session\n\nimport \"os\"\n\nfunc environ() string { return os.Getenv(\"PATH\") }\n\n" +
			"func Make() string { return os.Getenv(\"HOME\") }\n"},
	}}
	want := []string{
		"input rule: internal/cli/a.go:5 reads os.Getenv",
		"input rule: internal/cli/b.go:5 reads os.LookupEnv",
		"input rule: internal/cli/b.go:5 reads os.Stdin",
		"input rule: internal/gate/d.go:3 imports os with a dot, so the scan cannot read its uses",
		"input rule: internal/git/git.go:7 reads os.Getenv",
		"input rule: internal/psb/c.go:5 reads syscall.Getenv",
		"input rule: internal/session/session.go:7 reads os.Getenv",
		"input rule: internal/tsv/e.go:5 reads os.Environ",
	}
	if got := checkInputs(m); !slices.Equal(got, want) {
		t.Fatalf("the findings:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

func TestCheckInputsPassesAModuleThatReadsNoInput(t *testing.T) {
	m := module{sources: map[string]map[string]string{
		"cmd/layup": {"main.go": "package main\n\nimport \"os\"\n\nfunc main() { os.Exit(len(os.Args)) }\n"},
	}}
	if got := checkInputs(m); len(got) != 0 {
		t.Fatalf("the findings %q, want none", got)
	}
}
