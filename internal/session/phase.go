package session

// This file holds the files of a session of a phase (docs/spec/intake.md, The
// sessions of Intake; row 42a of the plan, task T-gmey, #215): its input/,
// and the read of its output file by the rules of a result file.

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"unicode/utf8"
)

// fileName refuses a name that is not one plain file name: empty, . or .., or
// with a path separator.
func fileName(name string) error {
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name {
		return fmt.Errorf("%q is not one file name", name)
	}
	return nil
}

// MakeInput makes input/ in the directory d of a session of a phase and
// writes each file in it by its name, byte for byte, mode 0600. It runs right
// after Make (whose last write is prompt.md) and before the start row.
func MakeInput(d Dir, files map[string][]byte) error {
	for name := range files {
		if err := fileName(name); err != nil {
			return fmt.Errorf("an input file: %w", err)
		}
	}
	if err := os.Mkdir(d.Input, 0o700); err != nil {
		return err
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(d.Input, name), data, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// ReadOutput reads the output file result/<name> of a session of a phase by
// the rules of a result file, with no block: the caller reads it by its own.
func ReadOutput(d Dir, name string) ([]byte, error) {
	if err := fileName(name); err != nil {
		return nil, fmt.Errorf("the output file: %w", err)
	}
	return readResult(filepath.Join(d.Result, name))
}

// ReadOutputFiles reads the output files of a session of a phase, whose output
// is a directory (docs/spec/scaffold.md, The prose rows; task T-eahr): result/<dir>
// must hold the files at the paths in names, and nothing else. The whole set is
// read and checked before any file is given: a name that is not a clean relative
// path; a missing file; an entry that no name covers (a file, a directory with no
// name below it, also a nested one); a link, as a file, as a directory above
// one, or as result/ or result/<dir> itself; an entry that is no regular file;
// a file over 1 MiB; and a file that is not valid UTF-8 refuse the output, and
// the files read are not given.
func ReadOutputFiles(d Dir, dir string, names []string) (map[string][]byte, error) {
	if err := fileName(dir); err != nil {
		return nil, fmt.Errorf("the output directory: %w", err)
	}
	wanted, above := map[string]bool{}, map[string]bool{}
	for _, n := range names {
		if !fs.ValidPath(n) || n == "." {
			return nil, fmt.Errorf("the output file name %q is not a clean relative path", n)
		}
		wanted[n] = true
		for p := path.Dir(n); p != "."; p = path.Dir(p) {
			above[p] = true
		}
	}
	root := filepath.Join(d.Result, dir)
	for _, p := range []string{d.Result, root} { // a link at either takes the read out of the session's directory
		if info, err := os.Lstat(p); err != nil {
			return nil, fmt.Errorf("the output directory: %w", err)
		} else if !info.IsDir() {
			return nil, fmt.Errorf("the output directory %s is a link, or not a directory", p)
		}
	}
	seen := map[string]bool{}
	err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil || rel == "." {
			return err
		}
		rel = filepath.ToSlash(rel)
		switch {
		case e.Type()&fs.ModeSymlink != 0:
			return fmt.Errorf("the output file %s is a link", rel)
		case wanted[rel]:
			seen[rel] = true
		case above[rel]:
			if !e.IsDir() {
				return fmt.Errorf("the output has a file at %s, where a directory is needed", rel)
			}
		default:
			return fmt.Errorf("the output holds %s, which the stop does not name", rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, n := range names {
		if !seen[n] {
			return nil, fmt.Errorf("the output file %s is missing", n)
		}
		data, err := readResult(filepath.Join(root, filepath.FromSlash(n)))
		if err != nil {
			return nil, err
		}
		if !utf8.Valid(data) {
			return nil, fmt.Errorf("the output file %s is not valid UTF-8", n)
		}
		files[n] = data
	}
	return files, nil
}
