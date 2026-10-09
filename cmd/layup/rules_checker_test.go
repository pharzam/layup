package main

// The checker of the package rules of docs/spec/packages.md ("The test of the
// package rules"), and the scan of the input rule of docs/spec/README.md
// (Commands: Arguments). It is test code, so the binary holds none of it.

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
	m2aHeading   = "### The table of M2a"
	m2bHeading   = "### The table of M2b"
	adapterLine  = "The one package that imports them:"      // the line of rule 5
	engineLine   = "The packages of the engine checks:"      // the line of NFR-005 (session.md)
	phase1Cell   = "(the row of phase 1"                     // a cell of the table of M2a that stands for the cell of phase 1
	registerCell = "(the command of a harness register row)" // the words of a register row's cell Starts a program
)

var (
	network = []string{"crypto/tls", "net", "net/http"} // rule 5
	span    = regexp.MustCompile("`([^`]*)`")           // a code span of a cell
)

// A row is one row of the tables: the packages of the module that a package
// may import, the program that it starts ("" for no), whether it starts the
// command that a harness register row names instead (a register row), and the
// packages of rule 5 that it may depend on (its cell Connects).
type row struct {
	imports  []string
	program  string
	register bool
	connects []string
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

// cellsUnder gives the rows of cells of the first table under heading, with
// its header first and no separator row. A missing heading, header or
// separator, or a table with no row, is an error.
func cellsUnder(text, heading string) ([][]string, error) {
	lines := strings.Split(text, "\n")
	start := slices.IndexFunc(lines, func(l string) bool { return strings.TrimSpace(l) == heading })
	if start < 0 {
		return nil, fmt.Errorf("no heading %q", heading)
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
		return nil, fmt.Errorf("no header, separator and row under %q", heading)
	}
	for _, cells := range table[2:] {
		if len(cells) != len(table[0]) {
			return nil, fmt.Errorf("the row %q has %d cells, want %d", cells, len(cells), len(table[0]))
		}
	}
	return append(table[:1], table[2:]...), nil
}

// columns gives the index of each named column of a header, and an error for
// a missing one.
func columns(header []string, heading string, names ...string) (map[string]int, error) {
	col := map[string]int{}
	for i, c := range header {
		col[c] = i
	}
	for _, name := range names {
		if _, ok := col[name]; !ok {
			return nil, fmt.Errorf("%s: no column %q", heading, name)
		}
	}
	return col, nil
}

// readRow reads the cells Package, May import and Starts a program.
func readRow(cells []string, col map[string]int) (string, row, error) {
	pkg, rest := spans(cells[col["Package"]])
	if len(pkg) != 1 || rest != "" {
		return "", row{}, fmt.Errorf("cannot read the package %q", cells[col["Package"]])
	}
	var r row
	if imp := cells[col["May import"]]; imp != "—" {
		list, rest := spans(imp)
		if len(list) == 0 || rest != strings.Repeat(",", len(list)-1) {
			return "", row{}, fmt.Errorf("%s: cannot read May import %q", pkg[0], imp)
		}
		r.imports = list
	}
	if prog := cells[col["Starts a program"]]; prog != "no" {
		list, _ := spans(prog)
		if len(list) == 1 && prog == "`"+list[0]+"` "+registerCell {
			r.register = true // the span names the register, not a program
			return pkg[0], r, nil
		}
		if strings.Contains(prog, registerCell) {
			return "", row{}, fmt.Errorf("%s: a register cell is one code span and %s, and nothing else: %q", pkg[0], registerCell, prog)
		}
		if len(list) != 1 || !strings.HasPrefix(prog, "`") || len(strings.Fields(list[0])) == 0 {
			return "", row{}, fmt.Errorf("%s: cannot read Starts a program %q", pkg[0], prog)
		}
		r.program = strings.Fields(list[0])[0]
	}
	return pkg[0], r, nil
}

// readTable reads the table of phase 1 and the table of M2a (docs/spec/
// packages.md, The test of the package rules). A row of M2a whose cells Job,
// May import and Starts a program each start with "(the row of phase 1" adds
// its cell Connects to the package's row of phase 1; any other row of M2a is
// a package of its own. A cell that it cannot read, a missing heading or
// column, a table with no row, a row of M2a that mixes the two kinds of cell,
// a row of phase 1 that phase 1 lacks, and a package in two rows are errors,
// so the test never passes with nothing checked.
func readTable(text string) (map[string]row, error) {
	rows := map[string]row{}
	phase1, err := cellsUnder(text, tableHeading)
	if err != nil {
		return nil, err
	}
	col, err := columns(phase1[0], tableHeading, "Package", "May import", "Starts a program")
	if err != nil {
		return nil, err
	}
	for _, cells := range phase1[1:] {
		pkg, r, err := readRow(cells, col)
		if err != nil {
			return nil, err
		}
		if _, ok := rows[pkg]; ok {
			return nil, fmt.Errorf("%s has a second row", pkg)
		}
		rows[pkg] = r
	}
	m2a, err := cellsUnder(text, m2aHeading)
	if err != nil {
		return nil, err
	}
	col, err = columns(m2a[0], m2aHeading, "Package", "Job", "May import", "Starts a program", "Connects")
	if err != nil {
		return nil, err
	}
	for _, cells := range m2a[1:] {
		pkg, rest := spans(cells[col["Package"]])
		if len(pkg) != 1 || rest != "" {
			return nil, fmt.Errorf("cannot read the package %q", cells[col["Package"]])
		}
		var connects []string
		if c := cells[col["Connects"]]; c != "—" {
			list, rest := spans(c)
			if len(list) == 0 || rest != strings.Repeat(",", len(list)-1) {
				return nil, fmt.Errorf("%s: cannot read Connects %q", pkg[0], c)
			}
			connects = list
		}
		of1 := 0
		for _, name := range []string{"Job", "May import", "Starts a program"} {
			if strings.HasPrefix(cells[col[name]], phase1Cell) {
				of1++
			}
		}
		r, has := rows[pkg[0]]
		switch {
		case of1 == 3 && !has:
			return nil, fmt.Errorf("%s: a row of M2a stands for its row of phase 1, and phase 1 has none", pkg[0])
		case of1 == 3:
			r.connects = connects
			rows[pkg[0]] = r
		case of1 != 0:
			return nil, fmt.Errorf("%s: a row of M2a mixes cells of phase 1 and cells of its own", pkg[0])
		case has:
			return nil, fmt.Errorf("%s has a second row", pkg[0])
		default:
			_, r, err := readRow(cells, col)
			if err != nil {
				return nil, err
			}
			r.connects = connects
			rows[pkg[0]] = r
		}
	}
	m2b, err := cellsUnder(text, m2bHeading)
	if err != nil {
		return nil, err
	}
	col, err = columns(m2b[0], m2bHeading, "Package", "Job", "May import", "Starts a program", "Connects")
	if err != nil {
		return nil, err
	}
	for _, cells := range m2b[1:] {
		for _, name := range []string{"Job", "May import", "Starts a program"} {
			if strings.HasPrefix(cells[col[name]], phase1Cell) {
				return nil, fmt.Errorf("%s: a row of M2b is a package of its own; a package of another table keeps its one row", cells[col["Package"]])
			}
		}
		pkg, r, err := readRow(cells, col)
		if err != nil {
			return nil, err
		}
		if _, has := rows[pkg]; has {
			return nil, fmt.Errorf("%s has a second row", pkg)
		}
		if c := cells[col["Connects"]]; c != "—" {
			list, rest := spans(c)
			if len(list) == 0 || rest != strings.Repeat(",", len(list)-1) {
				return nil, fmt.Errorf("%s: cannot read Connects %q", pkg, c)
			}
			r.connects = list
		}
		rows[pkg] = r
	}
	return rows, nil
}

// readAdapter reads the one package that rule 5 lets import a network package:
// the one code span of the one line that starts with adapterLine. No such
// line, two such lines, or a line with no code span or with two is an error.
func readAdapter(text string) (string, error) {
	var found []string
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), adapterLine) {
			found = append(found, l)
		}
	}
	if len(found) != 1 {
		return "", fmt.Errorf("%d lines start with %q, want 1", len(found), adapterLine)
	}
	list, _ := spans(found[0])
	if len(list) != 1 {
		return "", fmt.Errorf("the line %q has %d code spans, want 1", strings.TrimSpace(found[0]), len(list))
	}
	return list[0], nil
}

// readEngine reads the packages of the engine checks: the code spans of the
// one line that starts with engineLine. No such line, two such lines, or a
// line with no code span is an error.
func readEngine(text string) ([]string, error) {
	var found []string
	for _, l := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(l), engineLine) {
			found = append(found, l)
		}
	}
	if len(found) != 1 {
		return nil, fmt.Errorf("%d lines start with %q, want 1", len(found), engineLine)
	}
	list, _ := spans(found[0])
	if len(list) == 0 {
		return nil, fmt.Errorf("the line %q has no code span", strings.TrimSpace(found[0]))
	}
	return list, nil
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
// go list, and from the imports of the sources, which add the files behind a
// build constraint; rule 3 and the column Starts a program from a scan of the
// sources.
func checkRules(rows map[string]row, adapter string, engine []string, m module) []string {
	var f []string
	add := func(format string, a ...any) { f = append(f, fmt.Sprintf(format, a...)) }
	deps := map[string][]string{} // the dependencies of each package of the listing
	for _, p := range m.packages {
		deps[p.ImportPath] = p.Deps
	}
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
		imports := append(slices.Clone(p.Imports), fileImports(m.sources[rel])...)
		slices.Sort(imports)
		imports = slices.Compact(imports)
		reach := slices.Clone(p.Deps)
		for _, imp := range imports {
			reach = append(append(reach, imp), deps[imp]...)
		}
		r, hasRow := rows[rel]
		if slices.Contains(engine, rel) { // the rule of the engine checks reads no cell
			for _, d := range append([]string{m.path + "/internal/session"}, network...) {
				if slices.Contains(reach, d) {
					add("engine checks: %s depends on %s", rel, strings.TrimPrefix(d, m.path+"/"))
				}
			}
		}
		for _, d := range network {
			switch {
			case slices.Contains(imports, d) && rel != adapter:
				add("rule 5: %s imports %s; only %s imports them", rel, d, adapter)
			case slices.Contains(reach, d) && !slices.Contains(r.connects, d):
				add("rule 5: %s depends on %s", rel, d)
			}
		}
		if !hasRow {
			add("table: %s has no row", rel)
		}
		for _, imp := range imports {
			if imp == "os/exec" && hasRow && r.program == "" && !r.register {
				add("rule 4: %s imports os/exec, and its row starts no program", rel)
			}
			if dep, ok := inModule(m.path, imp); ok && hasRow && !slices.Contains(r.imports, dep) {
				add("May import: %s imports %s, which its row does not allow", rel, dep)
			}
		}
		for _, s := range scan(rel, m.sources[rel]) {
			switch {
			case !strings.HasPrefix(s.call, "exec."):
				add("Starts a program: %s %s", s.at, s.call)
			case s.program == "git" && rel != "internal/git":
				add("rule 3: %s starts git outside internal/git", s.at)
			case !hasRow: // the finding "has no row" says it
			case r.register && s.program != "":
				add("Starts a program: %s starts %s; a register row starts only the command of its register", s.at, s.program)
			case r.register && s.call != "exec.Command" && s.call != "exec.CommandContext":
				add("Starts a program: %s starts a program that the scan cannot read", s.at)
			case r.register: // a call whose program is the register's command
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
// string literal (else ""); an exec.Cmd that the code makes itself, whose
// program the scan cannot read; or the words of another finding.
type start struct{ at, call, program string }

// starters is the names that start a program, by import path.
var starters = map[string][]string{"os/exec": {"Command", "CommandContext", "Cmd"}, "os": {"StartProcess"},
	"syscall": {"Exec", "ForkExec", "StartProcess"}}

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
		seen := map[ast.Node]bool{} // a name that an outer node already gave, or *exec.Cmd
		ast.Inspect(file, func(n ast.Node) bool {
			if c, ok := n.(*ast.CallExpr); ok {
				if fn := starter(c.Fun, local); fn == "exec.Command" || fn == "exec.CommandContext" {
					seen[c.Fun] = true
					out = append(out, start{at(c), fn, literal(c.Args, fn == "exec.CommandContext")})
				}
			}
			if p, ok := n.(*ast.StarExpr); ok && starter(p.X, local) == "exec.Cmd" {
				seen[p.X] = true // a pointer type starts nothing; exec.Command made its value
			}
			if fn := starter(n, local); fn != "" && !seen[n] {
				if !strings.HasPrefix(fn, "exec.") {
					fn = "uses " + fn
				}
				if strings.HasPrefix(fn, "exec.") {
					fn += " as a value" // not a call: a register row does not allow it either
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

// fileImports gives the import paths of the files; a file that does not
// parse gives none, and the scan reports it.
func fileImports(files map[string]string) []string {
	var out []string
	for name, text := range files {
		file, _ := parser.ParseFile(token.NewFileSet(), name, text, parser.ImportsOnly)
		for _, imp := range file.Imports {
			if path, err := strconv.Unquote(imp.Path.Value); err == nil {
				out = append(out, path)
			}
		}
	}
	return out
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

// readers is the names that read an environment variable or the standard
// input, by import path.
var readers = map[string][]string{"os": {"Environ", "ExpandEnv", "Getenv", "LookupEnv", "Stdin"},
	"syscall": {"Environ", "Getenv"}}

// allowedReader is the one function, by package, whose reads the rule
// allows: environ of internal/git gives PATH and TMPDIR of the host to git
// (docs/spec/packages.md, D3 of #79). That is not an input of a command.
var allowedReader = map[string]string{"internal/git": "environ"}

// checkInputs gives each read of an environment variable or of the standard
// input in the non-test Go files of m, except in allowedReader. A file that
// does not parse, and an import of os or syscall with a dot, are findings,
// so the scan never passes with a file that it could not read.
func checkInputs(m module) []string {
	var out []string
	for _, rel := range slices.Sorted(maps.Keys(m.sources)) {
		files := m.sources[rel]
		for _, name := range slices.Sorted(maps.Keys(files)) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, name, files[name], parser.SkipObjectResolution)
			if err != nil {
				out = append(out, fmt.Sprintf("input rule: %s/%s does not parse: %v", rel, name, err))
				continue
			}
			at := func(n ast.Node) string { return fmt.Sprintf("%s/%s:%d", rel, name, fset.Position(n.Pos()).Line) }
			local := map[string]string{} // the import path of a reader, by its name in the file
			for _, imp := range file.Imports {
				if path, _ := strconv.Unquote(imp.Path.Value); readers[path] != nil {
					n := path
					if imp.Name != nil {
						n = imp.Name.Name
					}
					if n == "." {
						out = append(out, fmt.Sprintf("input rule: %s imports %s with a dot, so the scan cannot read its uses", at(imp), path))
					}
					local[n] = path
				}
			}
			for _, d := range file.Decls {
				fn, _ := d.(*ast.FuncDecl)
				if fn != nil && fn.Recv == nil && fn.Name.Name == allowedReader[rel] {
					continue
				}
				ast.Inspect(d, func(n ast.Node) bool {
					if sel, ok := n.(*ast.SelectorExpr); ok {
						if id, ok := sel.X.(*ast.Ident); ok && slices.Contains(readers[local[id.Name]], sel.Sel.Name) {
							out = append(out, fmt.Sprintf("input rule: %s reads %s.%s", at(sel), local[id.Name], sel.Sel.Name))
						}
					}
					return true
				})
			}
		}
	}
	return out
}
