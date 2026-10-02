package gate

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"syscall"

	"github.com/pharzam/layup/internal/git"
)

// gitAPI is the calls of internal/git that a run makes.
type gitAPI interface {
	Version() (string, error)
	RevParse(dir, rev string) (string, error)
	Show(dir, rev, path string) ([]byte, error)
	LsTree(dir, rev, path string) ([]git.TreeEntry, error)
	DiffNames(dir, base, head string) ([]string, error)
	WorktreeAdd(dir, path, rev string) error
	WorktreeRemove(dir, path string) error
}

type realGit struct{}

func (realGit) Version() (string, error)                     { return git.Version() }
func (realGit) RevParse(dir, rev string) (string, error)     { return git.RevParse(dir, rev) }
func (realGit) Show(dir, rev, p string) ([]byte, error)      { return git.Show(dir, rev, p) }
func (realGit) DiffNames(dir, b, h string) ([]string, error) { return git.DiffNames(dir, b, h) }
func (realGit) WorktreeAdd(dir, p, rev string) error         { return git.WorktreeAdd(dir, p, rev) }
func (realGit) WorktreeRemove(dir, p string) error           { return git.WorktreeRemove(dir, p) }
func (realGit) LsTree(dir, rev, p string) ([]git.TreeEntry, error) {
	return git.LsTree(dir, rev, p)
}

// A baseFile is a file of a config path at the base: its path and its mode.
type baseFile struct {
	path string
	exec bool
}

// baseFiles gives, for each config path of the kinds, the files that the base
// has at or under it; a path that the base does not have gives none. A
// symbolic link or a submodule at a config path is an input error.
func baseFiles(repo, base string, kinds []Kind) (map[string][]baseFile, error) {
	out := map[string][]baseFile{}
	for _, k := range kinds {
		for _, c := range k.Config {
			if _, ok := out[c]; ok {
				continue
			}
			entries, err := repoAPI.LsTree(repo, base, c)
			if err != nil {
				return nil, &InputError{fmt.Errorf("the config path %s of the kind %s at the base: %v", c, k.Name, firstLine(err))}
			}
			files := []baseFile{}
			for _, e := range entries {
				if e.Type != "blob" || (e.Mode != "100644" && e.Mode != "100755") {
					return nil, &InputError{fmt.Errorf("the config path %s of the kind %s: %s at the base is a %s of mode %s, not a file", c, k.Name, e.Path, e.Type, e.Mode)}
				}
				files = append(files, baseFile{path: e.Path, exec: e.Mode == "100755"})
			}
			out[c] = files
		}
	}
	return out, nil
}

// A scratch is the scratch work tree of the head, in a new temporary
// directory outside the repository. Each write goes through an os.Root, so a
// symbolic link of the head cannot send it out of the tree.
type scratch struct {
	repo, dir, tree string
	root            *os.Root
}

// newScratch adds the work tree of head, on no branch.
func newScratch(repo, head string) (*scratch, error) {
	dir, err := os.MkdirTemp("", "layup-gate-")
	if err != nil {
		return nil, err
	}
	s := &scratch{repo: repo, dir: dir, tree: filepath.Join(dir, "tree")}
	if err := repoAPI.WorktreeAdd(repo, s.tree, head); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}
	if s.root, err = os.OpenRoot(s.tree); err != nil {
		s.remove()
		return nil, err
	}
	return s, nil
}

// overlay writes the manifest of the base into the tree, and puts each config
// path as the base has it: its files with their modes, and no other file. The
// paths go in their sorted order, and a path under a file of the tree is
// absent, so two runs on one input give one tree (review round 1, finding 1).
func (s *scratch) overlay(base string, manifest []byte, files map[string][]baseFile) error {
	if err := s.write("docs/gates.tsv", manifest, false); err != nil {
		return err
	}
	for _, c := range slices.Sorted(maps.Keys(files)) {
		if err := s.removeAll(c); err != nil {
			return err
		}
		for _, f := range files[c] {
			data, err := repoAPI.Show(s.repo, base, f.path)
			if err != nil {
				return err
			}
			if err := s.write(f.path, data, f.exec); err != nil {
				return err
			}
		}
	}
	return nil
}

// write writes a file of the tree with its mode.
func (s *scratch) write(p string, data []byte, executable bool) error {
	mode := fs.FileMode(0o644)
	if executable {
		mode = 0o755
	}
	if err := s.root.MkdirAll(path.Dir(p), 0o755); err != nil {
		return err
	}
	if err := s.removeAll(p); err != nil {
		return err
	}
	if err := s.root.WriteFile(p, data, mode); err != nil {
		return err
	}
	return s.root.Chmod(p, mode)
}

// removeAll removes the path from the tree; a path under a file of the tree
// is absent already. A path that is the file .git of the tree, by any name (a
// file system that folds case, a hard link), is never removed: the tree would
// lose its record in the repository (review round 2, finding 1).
func (s *scratch) removeAll(p string) error {
	if git, err := s.root.Lstat(".git"); err == nil {
		if info, err := s.root.Lstat(p); err == nil && os.SameFile(git, info) {
			return fmt.Errorf("%s is the file .git of the scratch tree", p)
		}
	}
	if err := s.root.RemoveAll(p); err != nil && !errors.Is(err, syscall.ENOTDIR) {
		return err
	}
	return nil
}

// files gives the path of each file of the tree, without .git.
func (s *scratch) files() ([]string, error) {
	var out []string
	err := fs.WalkDir(s.root.FS(), ".", func(p string, d fs.DirEntry, err error) error {
		switch {
		case err != nil:
			return err
		case p == ".git":
			if d.IsDir() {
				return fs.SkipDir
			}
		case !d.IsDir():
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

// remove removes the work tree and its directory.
func (s *scratch) remove() error {
	if s.root != nil {
		s.root.Close()
	}
	err := repoAPI.WorktreeRemove(s.repo, s.tree)
	return errors.Join(err, os.RemoveAll(s.dir))
}

// shell runs a gate command with sh -c in dir and gives how it ended, and its
// standard output and standard error together, in their order.
func shell(dir, command string) (outcome, []byte) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	var exit *exec.ExitError
	switch {
	case err == nil:
		return outcome{}, out
	case errors.As(err, &exit):
		if ws, ok := exit.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return outcome{signal: ws.Signal().String()}, out
		}
		return outcome{code: exit.ExitCode()}, out
	}
	return outcome{code: -1}, append(out, []byte(err.Error()+"\n")...)
}
