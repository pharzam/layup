// Package catalog is the stack catalog of LAYUP: one entry per stack, with the
// gate kinds of the stack, the files that the setup writes into a target, and
// the known-bad fixture of each active kind (docs/spec/setup.md, The stack
// catalog). It reads an entry from an fs.FS whose root holds one directory per
// stack. It imports no package of this module other than internal/tsv.
package catalog

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/tsv"
)

const (
	suffix      = ".tmpl"      // the suffix of each file of files/ (K35)
	moduleToken = "{{module}}" // replaced with the module path in each file of files/
)

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
}

// Read reads the entry of stack from fsys and checks its rules: kinds.tsv
// matches the schema catalog-kinds and has a row; an active kind names
// fixtures/<kind>.patch, and that file exists; a pending kind has — as its
// fixture; each file of fixtures/ is the fixture of an active kind; each file
// of files/ ends with .tmpl. An error names the file.
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
		switch {
		case k.State == "pending" && k.Fixture != "":
			errs = append(errs, fmt.Errorf("catalog: %s: the pending kind %s names the fixture %q; a pending kind has —", at("kinds.tsv"), k.Name, k.Fixture))
		case k.State == "active" && k.Fixture != want:
			errs = append(errs, fmt.Errorf("catalog: %s: the active kind %s names the fixture %q; want %s", at("kinds.tsv"), k.Name, k.Fixture, want))
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
	if errs != nil {
		return nil, errors.Join(errs...)
	}
	return e, nil
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

// Kinds gives the kinds of the entry in the order of kinds.tsv.
func (e *Entry) Kinds() []Kind { return slices.Clone(e.kinds) }

// Files gives each file of files/ with its path in the target and its bytes,
// in the order of the paths, with each {{module}} replaced with module; no
// other byte changes. An empty module is an error.
func (e *Entry) Files(module string) ([]File, error) {
	if module == "" {
		return nil, errors.New("catalog: no module path")
	}
	var out []File
	for _, f := range e.files {
		data, err := fs.ReadFile(e.fsys, path.Join(e.stack, "files", f))
		if err != nil {
			return nil, fmt.Errorf("catalog: %w", err)
		}
		out = append(out, File{Path: target(f), Data: bytes.ReplaceAll(data, []byte(moduleToken), []byte(module))})
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
// directory: kinds.tsv, a file of files/ with its suffix, or a file of
// fixtures/.
func (e *Entry) Has(ref string) bool {
	stack, p, ok := strings.Cut(ref, "/")
	if !ok || stack != e.stack {
		return false
	}
	if p == "kinds.tsv" {
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
