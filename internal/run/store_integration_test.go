//go:build integration

package run

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// A records commit can remove files (docs/spec/scaffold.md, The stops of the
// runner: a fail after a prose session removes its files; task T-eahr): the
// same lock, work tree and push as a commit that writes, one commit for all the
// paths, and a path that the tree lacks or that leaves it is refused.
func TestARecordsCommitRemovesFiles(t *testing.T) {
	w, r, _, _ := sessionWorld(t)
	store := r.Store.(*recordsStore)
	ctx := context.Background()
	if err := store.Commit(ctx, map[string][]byte{"setup/inputs/files/a.md": []byte("a\n"), "setup/inputs/files/docs/b.md": []byte("b\n"), "keep.txt": []byte("k\n")}, "three files"); err != nil {
		t.Fatal(err)
	}
	if err := store.Remove(ctx, []string{"setup/inputs/files/a.md", "setup/inputs/files/docs/b.md"}, "the files of the session S-1a2b3c4d, removed: a check failed"); err != nil {
		t.Fatal(err)
	}
	files := strings.Fields(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records"))
	if !slices.Contains(files, "keep.txt") || slices.Contains(files, "setup/inputs/files/a.md") || slices.Contains(files, "setup/inputs/files/docs/b.md") {
		t.Errorf("the records after the removal: %q; want keep.txt and neither file", files)
	}
	if got := strings.TrimSpace(gitOut(t, w.bare, "log", "-1", "--format=%s", "layup-records")); got != "the files of the session S-1a2b3c4d, removed: a check failed" {
		t.Errorf("the message of the removal: %q", got)
	}
	if got := store.Base(); got != strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "layup-records")) {
		t.Errorf("the store's base %s is not the head of layup-records: the removal was not pushed through the store", got)
	}
	head := strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "layup-records"))
	// A path that leaves the work tree removes no file outside it: the work tree is made two
	// directories below TMPDIR, so ../../outside.txt is a file of TMPDIR that the test made.
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	outside := filepath.Join(tmp, "outside.txt")
	if err := os.WriteFile(outside, []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := store.Remove(ctx, []string{"../../outside.txt"}, "refused"); err == nil {
		t.Error("Remove of ../../outside.txt: no error; want a refusal")
	}
	if data, err := os.ReadFile(outside); err != nil || string(data) != "keep\n" {
		t.Errorf("the file outside the work tree after the refused removal: %q, %v; want it kept", data, err)
	}
	for _, paths := range [][]string{{"setup/inputs/files/nope.md"}, {"keep.txt", "setup/inputs/files/nope.md"}, {"../x"}, {"/etc/passwd"}, {"a//b"}, {""}} {
		if err := store.Remove(ctx, paths, "refused"); err == nil {
			t.Errorf("Remove(%q): no error; want a refusal", paths)
		}
	}
	if got := strings.TrimSpace(gitOut(t, w.bare, "rev-parse", "layup-records")); got != head {
		t.Errorf("a refused removal made a commit: %s; want the head %s", got, head)
	}
	if files := strings.Fields(gitOut(t, w.bare, "ls-tree", "-r", "--name-only", "layup-records")); !slices.Contains(files, "keep.txt") {
		t.Errorf("keep.txt went with a refused removal of it and a path that the tree lacks: %q", files)
	}
}

// The files under a directory of the records are read by their paths below it,
// and a link in that directory is refused: the work area copies them, and a link
// would reach a file of the host (task T-eahr).
func TestTheFilesOfTheRecordsAreReadByTheirPaths(t *testing.T) {
	w, r, _, _ := sessionWorld(t)
	store := r.Store.(*recordsStore)
	ctx := context.Background()
	if got, err := store.ReadFiles(ctx, "setup/inputs/files"); err != nil || len(got) != 0 {
		t.Errorf("a directory that the records lack: %q, %v; want no file", got, err)
	}
	if err := store.Commit(ctx, map[string][]byte{"setup/inputs/files/a.md": []byte("a\n"), "setup/inputs/files/docs/b.md": []byte("b\r\n"), "setup/other.txt": []byte("o\n")}, "files"); err != nil {
		t.Fatal(err)
	}
	got, err := store.ReadFiles(ctx, "setup/inputs/files")
	if want := map[string][]byte{"a.md": []byte("a\n"), "docs/b.md": []byte("b\r\n")}; err != nil || len(got) != 2 || string(got["a.md"]) != "a\n" || string(got["docs/b.md"]) != "b\r\n" {
		t.Errorf("the files: %q, %v; want %q, and not setup/other.txt", got, err, want)
	}
	// A link in the directory.
	clone := filepath.Join(t.TempDir(), "records")
	gitOut(t, "/", "clone", "-q", "--branch", "layup-records", "--single-branch", "file://"+w.bare, clone)
	if err := os.MkdirAll(filepath.Join(clone, "setup", "inputs", "files"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/etc/hostname", filepath.Join(clone, "setup", "inputs", "files", "link.md")); err != nil {
		t.Fatal(err)
	}
	gitOut(t, clone, "add", "-A")
	gitOut(t, clone, "commit", "-q", "-m", "a link")
	linked := &recordsStore{dir: clone, read: strings.TrimSpace(gitOut(t, clone, "rev-parse", "HEAD"))}
	if got, err := linked.ReadFiles(ctx, "setup/inputs/files"); err == nil || !strings.Contains(err.Error(), "link.md") {
		t.Errorf("a link in the directory: %q, %v; want an error that names link.md", got, err)
	}
}
