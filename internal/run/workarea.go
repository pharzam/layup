package run

// This file holds the work area of a target under layup run (docs/spec/scaffold.md,
// The work area; row 54 of the plan, task T-eahr, #277): DIR/work/OWNER/NAME/,
// removed and built again from the records, so nothing in it is the only copy of
// a value.

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/tsv"
	workpkg "github.com/pharzam/layup/internal/work"
)

// workSource is what the work area is built from: the problem statement and the
// vision brief of start/ (the vision is nil for none), the rows of
// inputs/answers.tsv, the prose texts of setup/inputs/files/ by their path
// below it, and setup/tree.tsv (nil while the step verify has not written it).
type workSource struct {
	Brief, Vision []byte
	Answers       [][]string
	Files         map[string][]byte
	Tree          []byte
}

// buildWorkArea removes the work area at dir and builds it again: the briefs,
// answers.tsv and the prose texts under inputs/, and target/, a clone of the
// branch of the run's clone with main at the root commit (the row of order 0 of
// setup/tree.tsv once it exists, else the one commit of the branch, which Start
// read back) and origin the URL of the forge. The runner reads refs/heads/main,
// so a default branch of another name gets main at the root commit as well
// (decided here, task T-eahr). A path of the prose texts that leaves inputs/files
// is refused before the disk changes, and so is a link at DIR/work/OWNER or at
// DIR/work, through which the removal would reach outside DIR.
func buildWorkArea(dir string, src workSource, clone, branch, url string) error {
	for p := range src.Files {
		if path.IsAbs(p) || path.Clean(p) != p || p == "." || p == ".." || strings.HasPrefix(p, "../") {
			return fmt.Errorf("setup/inputs/files/%s is not a path below inputs/files of the work area", p)
		}
	}
	for _, up := range []string{filepath.Dir(dir), filepath.Dir(filepath.Dir(dir))} {
		if info, err := os.Lstat(up); err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a link; the work area is DIR/work/OWNER/NAME, below no link", up)
		}
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	files := map[string][]byte{workpkg.BriefPath: src.Brief}
	if src.Vision != nil {
		files["inputs/briefs/vision.md"] = src.Vision
	}
	var answers bytes.Buffer
	if err := tsv.Write(&answers, workpkg.AnswersSchema, src.Answers); err != nil {
		return err
	}
	files[workpkg.AnswersPath] = answers.Bytes()
	for p, data := range src.Files {
		files["inputs/files/"+p] = data
	}
	for p, data := range files {
		dst := filepath.Join(dir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return err
		}
	}
	root, err := rootOf(src, clone, branch)
	if err != nil {
		return err
	}
	target := filepath.Join(dir, workpkg.TargetPath)
	if err := git.CloneLocal(clone, target, branch); err != nil {
		return err
	}
	if err := git.ResetHard(target); err != nil { // the clone has no checkout
		return err
	}
	tip, err := git.RevParse(target, "refs/heads/"+branch)
	if err != nil {
		return err
	}
	switch {
	case branch != "main":
		err = git.SwitchCreate(target, "main", root)
	case tip != root:
		if err = git.ResetSoft(target, root); err == nil {
			err = git.ResetHard(target)
		}
	}
	if err != nil {
		return err
	}
	return git.AddRemote(target, "origin", url)
}

// rootOf gives the root commit of the target: the first row of setup/tree.tsv,
// which the clone must have, else the one commit of the branch.
func rootOf(src workSource, clone, branch string) (string, error) {
	ref := "refs/heads/" + branch
	if src.Tree != nil {
		rows, err := ReadSetupTree(src.Tree)
		if err != nil {
			return "", fmt.Errorf("setup/tree.tsv: %w", err)
		}
		root := rows[0][1]
		if _, err := git.RevParse(clone, root+"^{commit}"); err != nil {
			return "", fmt.Errorf("the root commit %s of setup/tree.tsv is not in the clone of the default branch %s", root, branch)
		}
		return root, nil
	}
	tip, err := git.RevParse(clone, ref)
	if err != nil {
		return "", err
	}
	roots, err := git.RootCommits(clone, ref)
	if err != nil {
		return "", err
	}
	if len(roots) != 1 || roots[0] != tip {
		return "", fmt.Errorf("the default branch %s holds more than its root commit and the records have no setup/tree.tsv, so the root commit is not known", branch)
	}
	return tip, nil
}
