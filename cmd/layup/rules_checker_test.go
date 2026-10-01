package main

// The checker of the package rules of docs/spec/packages.md ("The test of the
// package rules"). It is test code, so the binary holds none of it.

import (
	"cmp"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"maps"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

const (
	modulePath   = "github.com/pharzam/layup" // the module of rule 1
	tableHeading = "### The table of phase 1"
)

var (
	network = []string{"crypto/tls", "net", "net/http"} // rule 5
	span    = regexp.MustCompile("`([^`]*)`")           // a code span of a cell
)

// A row is one row of the table: the packages of the module that a package
// may import, and the program that it starts ("" for no).
type row struct {
	imports []string
	program string
}

// A goPackage is one package, as go list -json gives it.
type goPackage struct {
	ImportPath, Name, Dir             string
	Standard                          bool
	Imports, Deps                     []string
	GoFiles, CgoFiles, IgnoredGoFiles []string
}

// A module is what the checker reads: the module path and the require lines
// of go.mod, the packages that go list -deps names, and the non-test Go files
// of the module's packages (name to text), by their path in the table.
type module struct {
	path     string
	requires []string
	packages []goPackage
	sources  map[string]map[string]string
}

// readTable reads the first table under tableHeading. A cell that it cannot
// read, a missing heading or column, and a table with no row are errors, so
// the test never passes with nothing checked.
func readTable(text string) (map[string]row, error) {
	lines := strings.Split(text, "\n")
	start := slices.IndexFunc(lines, func(l string) bool { return strings.TrimSpace(l) == tableHeading })
	if start < 0 {
		return nil, fmt.Errorf("no heading %q", tableHeading)
	}
	var table [][]string
	for _, l := range lines[start+1:] {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "#") || (table != nil && !strings.HasPrefix(l, "|")) {
			break
		}
		if strings.HasPrefix(l, "|") {
			cells := strings.Split(strings.TrimSuffix(strings.TrimPrefix(l, "|"), "|"), "|")
			for i := range cells {
				cells[i] = strings.TrimSpace(cells[i])
			}
			table = append(table, cells)
		}
	}
	if len(table) < 3 || strings.Trim(strings.Join(table[1], ""), "-:") != "" {
		return nil, fmt.Errorf("no header, separator and row under %q", tableHeading)
	}
	col := map[string]int{}
	for i, c := range table[0] {
		col[c] = i
	}
	for _, name := range []string{"Package", "May import", "Starts a program"} {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("no column %q", name)
		}
	}
	rows := map[string]row{}
	for _, cells := range table[2:] {
		if len(cells) != len(table[0]) {
			return nil, fmt.Errorf("the row %q has %d cells, want %d", cells, len(cells), len(table[0]))
		}
		pkg, rest := spans(cells[col["Package"]])
		if len(pkg) != 1 || rest != "" {
			return nil, fmt.Errorf("cannot read the package %q", cells[col["Package"]])
		}
		var r row
		if imp := cells[col["May import"]]; imp != "—" {
			list, rest := spans(imp)
			if len(list) == 0 || rest != strings.Repeat(",", len(list)-1) {
				return nil, fmt.Errorf("%s: cannot read May import %q", pkg[0], imp)
			}
			r.imports = list
		}
		if prog := cells[col["Starts a program"]]; prog != "no" {
			list, _ := spans(prog)
			if len(list) != 1 || !strings.HasPrefix(prog, "`") || len(strings.Fields(list[0])) == 0 {
				return nil, fmt.Errorf("%s: cannot read Starts a program %q", pkg[0], prog)
			}
			r.program = strings.Fields(list[0])[0]
		}
		if _, ok := rows[pkg[0]]; ok {
			return nil, fmt.Errorf("%s has a second row", pkg[0])
		}
		rows[pkg[0]] = r
	}
	return rows, nil
}

// spans gives the code spans of a cell, and the rest of the cell with no
// space.
func spans(cell string) ([]string, string) {
	var list []string
	for _, m := range span.FindAllStringSubmatch(cell, -1) {
		list = append(list, m[1])
	}
	return list, strings.ReplaceAll(span.ReplaceAllString(cell, ""), " ", "")
}

// checkRules gives the findings of the package rules in m, sorted; none means
// that m keeps them. Rules 1, 2, 4 and 5 and the column May import come from
// go list; rule 3 and the column Starts a program from a scan of the sources.
func checkRules(rows map[string]row, m module) []string {
	var f []string
	add := func(format string, a ...any) { f = append(f, fmt.Sprintf(format, a...)) }
	if m.path != modulePath {
		add("rule 1: the module is %s, want %s", m.path, modulePath)
	}
	for _, r := range m.requires {
		add("rule 2: go.mod requires %s", r)
	}
	binary := false
	for _, p := range m.packages {
		rel, ok := inModule(m.path, p.ImportPath)
		if !ok {
			if !p.Standard {
				add("rule 2: %s is in neither the standard library nor this module", p.ImportPath)
			}
			continue
		}
		switch {
		case p.Name == "main" && rel == "cmd/layup":
			binary = true
		case p.Name == "main":
			add("rule 1: %s is a second binary", rel)
		case !strings.HasPrefix(rel, "internal/"):
			add("rule 1: %s is not under internal/", rel)
		}
		for _, d := range p.Deps {
			if slices.Contains(network, d) {
				add("rule 5: %s depends on %s", rel, d)
			}
		}
		r, hasRow := rows[rel]
		if !hasRow {
			add("table: %s has no row", rel)
		}
		for _, imp := range p.Imports {
			if imp == "os/exec" && hasRow && r.program == "" {
				add("rule 4: %s imports os/exec, and its row starts no program", rel)
			}
			if dep, ok := inModule(m.path, imp); ok && hasRow && !slices.Contains(r.imports, dep) {
				add("May import: %s imports %s, which its row does not allow", rel, dep)
			}
		}
		for _, s := range scan(rel, m.sources[rel]) {
			switch {
			case s.call != "exec.Command" && s.call != "exec.CommandContext":
				add("Starts a program: %s %s", s.at, s.call)
			case s.program == "git" && rel != "internal/git":
				add("rule 3: %s starts git outside internal/git", s.at)
			case !hasRow: // the finding "has no row" says it
			case s.program == "":
				add("Starts a program: %s starts a program that the scan cannot read", s.at)
			case s.program != r.program:
				add("Starts a program: %s starts %s; its row says %s", s.at, s.program, cmp.Or(r.program, "no"))
			}
		}
	}
	if !binary {
		add("rule 1: no binary cmd/layup")
	}
	slices.Sort(f)
	return f
}

// inModule gives the path of pkg in the table, and whether pkg is in the
// module at path.
func inModule(path, pkg string) (string, bool) {
	if pkg != path && !strings.HasPrefix(pkg, path+"/") {
		return "", false
	}
	return strings.TrimPrefix(strings.TrimPrefix(pkg, path), "/"), true
}

// A start is a place in a source where a program can start: a call or another
// use of exec.Command or exec.CommandContext, with its program when that is a
// string literal (else ""), or the words of another finding.
type start struct{ at, call, program string }

// starters is the functions that start a program, by import path.
var starters = map[string][]string{"os/exec": {"Command", "CommandContext"}, "os": {"StartProcess"}, "syscall": {"Exec"}}

// scan gives each start of a program in the files of one package.
func scan(rel string, files map[string]string) []start {
	var out []start
	for _, name := range slices.Sorted(maps.Keys(files)) {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, name, files[name], parser.SkipObjectResolution)
		if err != nil {
			out = append(out, start{at: rel + "/" + name, call: "does not parse: " + err.Error()})
			continue
		}
		at := func(n ast.Node) string { return fmt.Sprintf("%s/%s:%d", rel, name, fset.Position(n.Pos()).Line) }
		local := map[string]string{} // the import path of a starter, by its name in the file
		for _, imp := range file.Imports {
			if path, _ := strconv.Unquote(imp.Path.Value); starters[path] != nil {
				n := path[strings.LastIndex(path, "/")+1:]
				if imp.Name != nil {
					n = imp.Name.Name
				}
				if n == "." {
					out = append(out, start{at: at(imp), call: "imports " + path + " with a dot, so the scan cannot read its calls"})
				}
				local[n] = path
			}
		}
		called := map[ast.Node]bool{}
		ast.Inspect(file, func(n ast.Node) bool {
			if c, ok := n.(*ast.CallExpr); ok {
				if fn := starter(c.Fun, local); fn == "exec.Command" || fn == "exec.CommandContext" {
					called[c.Fun] = true
					out = append(out, start{at(c), fn, literal(c.Args, fn == "exec.CommandContext")})
				}
			}
			if fn := starter(n, local); fn != "" && !called[n] {
				if !strings.HasPrefix(fn, "exec.") {
					fn = "uses " + fn
				}
				out = append(out, start{at: at(n), call: fn}) // not a call: the program is unknown
			}
			return true
		})
	}
	return out
}

// starter gives the function of starters that n names, as package.Function,
// or "".
func starter(n ast.Node, local map[string]string) string {
	sel, ok := n.(*ast.SelectorExpr)
	if !ok {
		return ""
	}
	id, ok := sel.X.(*ast.Ident)
	if !ok || !slices.Contains(starters[local[id.Name]], sel.Sel.Name) {
		return ""
	}
	path := local[id.Name]
	return path[strings.LastIndex(path, "/")+1:] + "." + sel.Sel.Name
}

// literal gives the program of a call when it is a string literal: the first
// argument, or the second after a context.
func literal(args []ast.Expr, afterContext bool) string {
	i := 0
	if afterContext {
		i = 1
	}
	if i < len(args) {
		if lit, ok := args[i].(*ast.BasicLit); ok && lit.Kind == token.STRING {
			s, _ := strconv.Unquote(lit.Value)
			return s
		}
	}
	return ""
}
