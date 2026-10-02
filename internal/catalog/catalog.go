// Package catalog is the stack catalog of LAYUP: one entry per stack, with the
// gate kinds of the stack, the files that the setup writes into a target, and
// the known-bad fixture of each active kind (docs/spec/setup.md, The stack
// catalog). It reads an entry from an fs.FS whose root holds one directory per
// stack; Embedded is the root of the entries of the binary. It imports no
// package of this module other than internal/tsv.
package catalog

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

const (
	suffix      = ".tmpl"      // the suffix of each file of files/ (K35)
	moduleToken = "{{module}}" // replaced with the module path in each file of files/
	// The two quotes of a marker (docs/spec/README.md). No file of an entry
	// holds one, because LAYUP's own check markers would read it: a gap of an
	// entry is a gap token, which becomes a marker when the setup writes the
	// file (D7 of #91).
	markOpen, markClose = "\u2039", "\u203a"
	gapPrefix           = "{{gap:"
)

// gapToken is a gap token of a file of files/: {{gap:<text>}}, whose text
// becomes the text of the marker.
var gapToken = regexp.MustCompile(`\{\{gap:([^{}\n]+)\}\}`)

// KindsSchema is the form of kinds.tsv: the block catalog-kinds of
// docs/spec/setup.md.
var KindsSchema = tsv.Schema{Name: "catalog-kinds", Location: "layup:internal/catalog/<stack>/kinds.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "id(<word>)", Key: true},
	{Name: "state", Type: "enum(active|pending)"},
	{Name: "tool", Type: "text"},
	{Name: "version", Type: "text"},
	{Name: "command", Type: "text"},
	{Name: "scope", Type: "list(text)"},
	{Name: "config", Type: "list(path)"},
	{Name: "fixture", Type: "text"},
	{Name: "evidence", Type: "text"},
}}

// ManifestSchema is the form of the manifest that an entry gives: the block
// gate-manifest of docs/spec/gate.md. internal/catalog may not import
// internal/gate, so the two packages each hold it, and each test compares it
// with the one block.
var ManifestSchema = tsv.Schema{Name: "gate-manifest", Location: "target:docs/gates.tsv", Columns: []tsv.Column{
	{Name: "kind", Type: "id(<word>)", Key: true},
	{Name: "state", Type: "enum(active|pending)"},
	{Name: "tool", Type: "text"},
	{Name: "command", Type: "text"},
	{Name: "scope", Type: "list(text)"},
	{Name: "config", Type: "list(path)"},
}}

// GapsSchema is the form of gaps.tsv, the gaps of an entry: the block
// catalog-gaps of docs/spec/setup.md (D7 of #91).
var GapsSchema = tsv.Schema{Name: "catalog-gaps", Location: "layup:internal/catalog/<stack>/gaps.tsv", Columns: []tsv.Column{
	{Name: "path", Type: "path", Key: true},
	{Name: "marker", Type: "text", Key: true},
	{Name: "question", Type: "text"},
}}

// A Kind is one row of kinds.tsv. Scope and Config hold the values of their
// lists; an empty value (—) is "" or no value.
type Kind struct {
	Name, State, Tool, Version, Command string
	Scope, Config                       []string
	Fixture, Evidence                   string
}

// A File is a file that the setup writes into the target: its path in the
// target and its bytes.
type File struct {
	Path string
	Data []byte
}

// An Entry is the catalog entry of one stack, read and checked by Read.
type Entry struct {
	fsys     fs.FS
	stack    string
	kinds    []Kind
	files    []string // the paths under files/, with the suffix, in path order
	fixtures []string // the paths under fixtures/, in path order
	hasGaps  bool     // the entry has gaps.tsv
	gaps     []Gap    // the gaps of gaps.tsv, by path, line and marker
}

// Read reads the entry of stack from fsys and checks its rules: kinds.tsv
// matches the schema catalog-kinds and has a row; an active kind names
// fixtures/<kind>.patch, and that file exists, has a version, and has an
// https URL as its evidence; a pending kind has — in each column but its kind,
// its state and its scope (D1 of #91, with note 1 of its plan review); each
// file of fixtures/ is the fixture of an active kind; each file of files/ ends
// with .tmpl; gaps.tsv, when there is one, matches the schema catalog-gaps,
// each of its rows names a file of files/ that holds its gap token once, with
// a question, and each gap token of a file has its row; no file of the entry
// holds a marker character (D7). An error names the file.
func Read(fsys fs.FS, stack string) (*Entry, error) {
	if !fs.ValidPath(stack) || strings.Contains(stack, "/") || stack == "." {
		return nil, fmt.Errorf("catalog: %q is not the name of an entry", stack)
	}
	at := func(p string) string { return path.Join(stack, p) }
	data, err := fs.ReadFile(fsys, at("kinds.tsv"))
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	rows, err := tsv.Read(data, KindsSchema)
	if err != nil {
		return nil, fmt.Errorf("catalog: %s: %w", at("kinds.tsv"), err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("catalog: %s: no kind", at("kinds.tsv"))
	}
	e := &Entry{fsys: fsys, stack: stack}
	for _, r := range rows {
		e.kinds = append(e.kinds, Kind{Name: r[0], State: r[1], Tool: r[2], Version: r[3], Command: r[4],
			Scope: strings.Fields(r[5]), Config: strings.Fields(r[6]), Fixture: r[7], Evidence: r[8]})
	}
	if e.files, err = filesUnder(fsys, at("files")); err != nil {
		return nil, err
	}
	if e.fixtures, err = filesUnder(fsys, at("fixtures")); err != nil {
		return nil, err
	}
	var errs []error
	for _, f := range e.files {
		if !strings.HasSuffix(f, suffix) || target(f) == "" {
			errs = append(errs, fmt.Errorf("catalog: %s: a file of files/ ends with %s", at("files/"+f), suffix))
		}
	}
	named := map[string]bool{}
	for _, k := range e.kinds {
		want := "fixtures/" + k.Name + ".patch"
		if k.State == "pending" {
			for _, c := range [][2]string{{"tool", k.Tool}, {"version", k.Version}, {"command", k.Command}, {"config", strings.Join(k.Config, " ")}, {"evidence", k.Evidence}} {
				if c[1] != "" {
					errs = append(errs, fmt.Errorf("catalog: %s: the pending kind %s names the %s %s; a pending kind has —", at("kinds.tsv"), k.Name, c[0], c[1]))
				}
			}
		}
		if k.State == "active" && k.Version == "" {
			errs = append(errs, fmt.Errorf("catalog: %s: the active kind %s has no version", at("kinds.tsv"), k.Name))
		}
		if k.State == "active" && (!strings.HasPrefix(k.Evidence, "https://") || len(k.Evidence) == len("https://") || strings.Contains(k.Evidence, " ")) {
			errs = append(errs, fmt.Errorf("catalog: %s: the evidence of the active kind %s is not an https URL: %s", at("kinds.tsv"), k.Name, cmp.Or(k.Evidence, "—")))
		}
		switch {
		case k.State == "pending" && k.Fixture != "":
			errs = append(errs, fmt.Errorf("catalog: %s: the pending kind %s names the fixture %s; a pending kind has —", at("kinds.tsv"), k.Name, k.Fixture))
		case k.State == "active" && k.Fixture != want:
			errs = append(errs, fmt.Errorf("catalog: %s: the active kind %s names the fixture %s; want %s", at("kinds.tsv"), k.Name, cmp.Or(k.Fixture, "—"), want))
		case k.State == "active" && !slices.Contains(e.fixtures, k.Name+".patch"):
			errs = append(errs, fmt.Errorf("catalog: %s: the fixture of the active kind %s does not exist", at(want), k.Name))
		case k.State == "active":
			named[k.Name+".patch"] = true
		}
	}
	for _, f := range e.fixtures {
		if !named[f] {
			errs = append(errs, fmt.Errorf("catalog: %s: the file is the fixture of no active kind (no active kind %s)", at("fixtures/"+f), strings.TrimSuffix(f, ".patch")))
		}
	}
	gapErrs, err := e.readGaps()
	if err != nil {
		return nil, err
	}
	if errs = append(errs, gapErrs...); errs != nil {
		return nil, errors.Join(errs...)
	}
	return e, nil
}

// readGaps reads gaps.tsv, when the entry has one, and checks the gap rules
// and the marker rule of Read (D7 of #91). It gives the broken rules, or an
// error when a file cannot be read.
func (e *Entry) readGaps() ([]error, error) {
	at := func(p string) string { return path.Join(e.stack, p) }
	var errs []error
	texts := map[string]string{} // the text of each file of files/, by its path in the target
	paths := []string{"kinds.tsv", "gaps.tsv"}
	for _, f := range e.files {
		paths = append(paths, "files/"+f)
	}
	for _, f := range e.fixtures {
		paths = append(paths, "fixtures/"+f)
	}
	var gaps []byte
	for _, p := range paths {
		data, err := fs.ReadFile(e.fsys, at(p))
		if p == "gaps.tsv" && errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("catalog: %w", err)
		}
		if bytes.Contains(data, []byte(markOpen)) || bytes.Contains(data, []byte(markClose)) {
			errs = append(errs, fmt.Errorf("catalog: %s: a marker character; an entry writes a gap token in its place", at(p)))
		}
		switch f, ok := strings.CutPrefix(p, "files/"); {
		case p == "gaps.tsv":
			e.hasGaps, gaps = true, data
		case ok:
			texts[target(f)] = string(data)
			if strings.Count(string(data), gapPrefix) != len(gapToken.FindAllIndex(data, -1)) {
				errs = append(errs, fmt.Errorf("catalog: %s: a gap token of another form; a gap token is {{gap:<text>}} on one line", at(p)))
			}
		}
	}
	var rows [][]string
	if e.hasGaps {
		var err error
		if rows, err = tsv.Read(gaps, GapsSchema); err != nil {
			return append(errs, fmt.Errorf("catalog: %s: %w", at("gaps.tsv"), err)), nil
		}
	}
	listed := map[[2]string]bool{}
	for i, r := range rows {
		line := i + 2
		listed[[2]string{r[0], r[1]}] = true
		text, ok := texts[r[0]]
		if !ok {
			errs = append(errs, fmt.Errorf("catalog: %s: line %d: %s is not a file of files/", at("gaps.tsv"), line, r[0]))
			continue
		}
		if r[2] == "" {
			errs = append(errs, fmt.Errorf("catalog: %s: line %d: the gap %s has no question", at("gaps.tsv"), line, r[1]))
		}
		token := gapPrefix + r[1] + "}}"
		if n := strings.Count(text, token); n != 1 {
			errs = append(errs, fmt.Errorf("catalog: %s: line %d: the file %s holds the gap token %s %d times; want 1", at("gaps.tsv"), line, r[0], token, n))
			continue
		}
		e.gaps = append(e.gaps, Gap{Path: r[0], Line: strings.Count(text[:strings.Index(text, token)], "\n") + 1, Marker: markOpen + r[1] + markClose, Question: r[2]})
	}
	for _, f := range e.files {
		for _, m := range gapToken.FindAllStringSubmatch(texts[target(f)], -1) {
			if !listed[[2]string{target(f), m[1]}] {
				errs = append(errs, fmt.Errorf("catalog: %s: the gap token %s has no row in %s", at("files/"+f), m[0], at("gaps.tsv")))
			}
		}
	}
	slices.SortFunc(e.gaps, func(a, b Gap) int {
		return cmp.Or(strings.Compare(a.Path, b.Path), cmp.Compare(a.Line, b.Line), strings.Compare(a.Marker, b.Marker))
	})
	return errs, nil
}

// filesUnder gives the path of each file under dir, relative to dir, in path
// order; a missing dir gives none.
func filesUnder(fsys fs.FS, dir string) ([]string, error) {
	var out []string
	err := fs.WalkDir(fsys, dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == dir && errors.Is(err, fs.ErrNotExist) {
				return fs.SkipDir
			}
			return err
		}
		if !d.IsDir() {
			out = append(out, strings.TrimPrefix(p, dir+"/"))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	return out, nil
}

// target gives the path in the target of a file of files/: its path without
// the suffix, or "" when that is not a path of a file.
func target(f string) string {
	t := strings.TrimSuffix(f, suffix)
	if !fs.ValidPath(t) || t == "." || path.Base(t) == "" || strings.HasSuffix(t, "/") {
		return ""
	}
	return t
}

// A Gap is a gap of an entry (D7 of #91): the path of a file in the target,
// the line of its marker, the marker, and the question for open-gaps.tsv. The
// setup writes the row of open-gaps.tsv and the record row
// marker:<path>:<line> of each (S12, row 15 of the plan).
type Gap struct {
	Path     string
	Line     int
	Marker   string
	Question string
}

// Gaps gives each gap of the entry, by path, line and marker.
func (e *Entry) Gaps() []Gap { return slices.Clone(e.gaps) }

// Kinds gives the kinds of the entry in the order of kinds.tsv.
func (e *Entry) Kinds() []Kind { return slices.Clone(e.kinds) }

// Files gives each file of files/ with its path in the target and its bytes,
// in the order of the paths, with each {{module}} replaced with module and
// each gap token with its marker; no other byte changes. A module that is
// empty or holds white space is an error: a module path has none, and a line
// feed would move the line of a gap.
func (e *Entry) Files(module string) ([]File, error) {
	if module == "" || strings.ContainsAny(module, " \t\r\n") {
		return nil, fmt.Errorf("catalog: the module path %q is empty or holds white space", module)
	}
	var out []File
	for _, f := range e.files {
		data, err := fs.ReadFile(e.fsys, path.Join(e.stack, "files", f))
		if err != nil {
			return nil, fmt.Errorf("catalog: %w", err)
		}
		data = bytes.ReplaceAll(data, []byte(moduleToken), []byte(module))
		out = append(out, File{Path: target(f), Data: gapToken.ReplaceAll(data, []byte(markOpen+"${1}"+markClose))})
	}
	slices.SortFunc(out, func(a, b File) int { return strings.Compare(a.Path, b.Path) })
	return out, nil
}

// Fixture gives the bytes of the fixture of an active kind.
func (e *Entry) Fixture(kind string) ([]byte, error) {
	for _, k := range e.kinds {
		if k.Name == kind && k.State == "active" {
			data, err := fs.ReadFile(e.fsys, path.Join(e.stack, k.Fixture))
			if err != nil {
				return nil, fmt.Errorf("catalog: %w", err)
			}
			return data, nil
		}
	}
	return nil, fmt.Errorf("catalog: %s has no active kind %q", e.stack, kind)
}

// Manifest gives the bytes of the manifest (docs/gates.tsv in the target):
// the table of kinds.tsv without the columns version, fixture and evidence.
func (e *Entry) Manifest() ([]byte, error) {
	var rows [][]string
	for _, k := range e.kinds {
		scope, err := tsv.JoinList(k.Scope)
		if err != nil {
			return nil, err
		}
		config, err := tsv.JoinList(k.Config)
		if err != nil {
			return nil, err
		}
		rows = append(rows, []string{k.Name, k.State, k.Tool, k.Command, scope, config})
	}
	var b bytes.Buffer
	if err := tsv.Write(&b, ManifestSchema, rows); err != nil {
		return nil, fmt.Errorf("catalog: the manifest of %s: %w", e.stack, err)
	}
	return b.Bytes(), nil
}

// Config gives the config paths of a kind in their order, and whether the
// entry has the kind.
func (e *Entry) Config(kind string) ([]string, bool) {
	for _, k := range e.kinds {
		if k.Name == kind {
			return slices.Clone(k.Config), true
		}
	}
	return nil, false
}

// Has reports whether ref, the ref of a catalog value of the setup record
// (<stack>/<path>), names a file of the entry by its path in the entry's
// directory: kinds.tsv, gaps.tsv, a file of files/ with its suffix, or a file
// of fixtures/.
func (e *Entry) Has(ref string) bool {
	stack, p, ok := strings.Cut(ref, "/")
	if !ok || stack != e.stack {
		return false
	}
	if p == "kinds.tsv" || p == "gaps.tsv" && e.hasGaps {
		return true
	}
	if f, ok := strings.CutPrefix(p, "files/"); ok {
		return slices.Contains(e.files, f)
	}
	if f, ok := strings.CutPrefix(p, "fixtures/"); ok {
		return slices.Contains(e.fixtures, f)
	}
	return false
}

// Stacks gives the name of each entry of fsys: each directory at its root
// that holds a kinds.tsv, in name order.
func Stacks(fsys fs.FS) ([]string, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("catalog: %w", err)
	}
	var out []string
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		if _, err := fs.Stat(fsys, path.Join(d.Name(), "kinds.tsv")); err == nil {
			out = append(out, d.Name())
		}
	}
	return out, nil
}
