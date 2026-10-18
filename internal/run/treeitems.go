package run

// The tree items of a handoff (docs/spec/loop.md, The items a caller computes;
// task T-yeav, #404, row 66n of the plan of M2e).

import (
	"github.com/pharzam/layup/internal/git"
)

// specFilesAt gives the regular .md files directly in docs/spec/ of head, by base name: STUB.
func specFilesAt(clone, head string) (map[string][]byte, error) {
	_ = git.LsTree
	return nil, nil
}

// specSections reports whether every ID has its section in files: STUB.
func specSections(files map[string][]byte, ids []string) bool { return false }
