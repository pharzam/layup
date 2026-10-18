package run

// The tree items of a handoff (docs/spec/loop.md, The items a caller computes;
// task T-yeav, #404, row 66n of the plan of M2e): spec-sections and
// rendering-equal are computed from the tree at a head that code holds, in the
// row of the pair that lists them, and a value of the caller's is never read.

import (
	"context"
	"path"
	"slices"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/handoff"
	"github.com/pharzam/layup/internal/rules"
	"github.com/pharzam/layup/internal/spec"
)

// specFilesAt gives the regular .md files directly in docs/spec/ of head, by
// base name: not a file below it, a link, a gitlink, an executable or a file
// of another extension (git ls-tree is recursive and lists a link as a blob).
func specFilesAt(clone, head string) (map[string][]byte, error) {
	return filesAt(clone, head, "docs/spec", func(base string) bool { return strings.HasSuffix(base, ".md") })
}

// filesAt gives the regular files (type blob, mode 100644) directly in dir of
// head whose base name keep accepts, by base name, with their exact bytes
// (task T-cav3, #412: the one selector of the files of a directory at a head).
func filesAt(clone, head, dir string, keep func(base string) bool) (map[string][]byte, error) {
	entries, err := git.LsTree(clone, head, dir)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, e := range entries {
		if e.Type != "blob" || e.Mode != "100644" || path.Dir(e.Path) != dir || !keep(path.Base(e.Path)) {
			continue
		}
		if files[path.Base(e.Path)], err = git.Show(clone, head, e.Path); err != nil {
			return nil, err
		}
	}
	return files, nil
}

// specSections reports whether each of the IDs has its section in files; no ID is false.
func specSections(files map[string][]byte, ids []string) bool {
	for _, id := range ids {
		if !spec.HasSection(files, id) {
			return false
		}
	}
	return len(ids) > 0
}

// treeItems sets spec-sections and rendering-equal: false unless the row lists
// the item, the handoff has a head and the tree at that head satisfies it. The
// head is the SHA of the bound event of the session, or the Head that code
// sets on a handoff that code wrote, and never the clone's HEAD or a Head of a
// session's caller. A failed read is an error, and nothing is written.
func (s *Sessions) treeItems(ctx context.Context, in HandoffInput, row handoff.Row, events [][]string, items map[string]bool) error {
	items["spec-sections"], items["rendering-equal"] = false, false
	wantSpec, wantRender := slices.Contains(row.Items, "spec-sections"), slices.Contains(row.Items, "rendering-equal")
	head := in.Head
	if !in.CodeWritten {
		head = boundHead(events, in.Session)
	}
	if head == "" || (!wantSpec && !wantRender) {
		return nil
	}
	if wantSpec {
		task, err := s.taskRow(ctx, in.Task)
		if err != nil {
			return err
		}
		files, err := specFilesAt(s.Clone, head)
		if err != nil {
			return err
		}
		items["spec-sections"] = specSections(files, task.Requirements)
	}
	if wantRender {
		var paths []string
		for _, a := range in.Artifacts {
			if a.Kind == "rendering" {
				paths = append(paths, a.Path)
			}
		}
		ok, err := rules.RenderingEqual(rules.Git(s.Clone), head, in.PreRev, paths, in.Render)
		if err != nil {
			return err
		}
		items["rendering-equal"] = ok
	}
	return nil
}
