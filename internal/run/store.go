package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/tsv"
)

// records is the store of 25a on internal/git (docs/spec/run.md, The lease and
// fencing): the run's clone, the URL of the target and the token. It keeps the
// commit ID of its own last pushed records commit, and makes each records
// commit on that commit in a scratch work tree, so it never builds on a commit
// of another run; before its first push it builds on the commit of its last
// read. One lock covers each read and each commit with its push, so the
// heartbeat and the writes of the steps make one chain.
type recordsStore struct {
	mu     sync.Mutex
	dir    string
	url    string
	auth   func(context.Context) (git.Auth, error)
	who    func() git.Identity
	pushed string // the run's last pushed records commit
	read   string // the commit of its last read
}

const recordsRef = "refs/heads/layup-records"

// ReadLease fetches layup-records into the run's clone and reads lease.tsv
// there. The run never builds on that branch: its base is the commit it
// pushed last, so the fetch moves no base.
func (s *recordsStore) ReadLease(ctx context.Context) (LeaseRow, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, err := s.auth(ctx)
	if err != nil {
		return LeaseRow{}, err
	}
	if err := git.Fetch(s.dir, s.url, recordsRef, a); err != nil {
		return LeaseRow{}, err
	}
	id, err := git.RevParse(s.dir, recordsRef)
	if err != nil {
		return LeaseRow{}, err
	}
	row, err := s.leaseAt(id)
	if err == nil {
		s.read = id
	}
	return row, err
}

// leaseAt reads the lease row of a records commit.
func (s *recordsStore) leaseAt(rev string) (LeaseRow, error) {
	data, err := git.Show(s.dir, rev, "lease.tsv")
	if err != nil {
		return LeaseRow{}, err
	}
	rows, err := records.ReadLease(data)
	if err != nil {
		return LeaseRow{}, fmt.Errorf("lease.tsv: %w", err)
	}
	return leaseRow(rows[0])
}

// WriteLease writes the lease row with a records commit and its push.
func (s *recordsStore) WriteLease(ctx context.Context, row LeaseRow, message string) error {
	data, err := leaseTable(row)
	if err != nil {
		return err
	}
	return s.commit(ctx, map[string][]byte{"lease.tsv": data}, message, "")
}

// ReadFile gives the file at path of the run's records base, the commit it
// pushed last or else of its last read, or nil when that commit has none.
func (s *recordsStore) ReadFile(ctx context.Context, path string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	base := s.pushed
	if base == "" {
		base = s.read
	}
	if base == "" {
		return nil, errors.New("no records commit to read")
	}
	entries, err := git.LsTree(s.dir, base, path)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, nil
	}
	return git.Show(s.dir, base, path)
}

// ReadFiles gives each file at or under dir in the run's records base, by its
// path below dir (none when the base has no such directory). An entry that is
// not a regular file, a link or a submodule, is an error (task T-eahr: the work
// area copies these files, and a link would reach a file of the host).
func (s *recordsStore) ReadFiles(ctx context.Context, dir string) (map[string][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	base := s.pushed
	if base == "" {
		base = s.read
	}
	if base == "" {
		return nil, errors.New("no records commit to read")
	}
	entries, err := git.LsTree(s.dir, base, dir)
	if err != nil {
		return nil, err
	}
	files := map[string][]byte{}
	for _, e := range entries {
		if !regular(e) {
			return nil, fmt.Errorf("%s of the records is not a regular file (mode %s, %s)", e.Path, e.Mode, e.Type)
		}
		data, err := git.Show(s.dir, base, e.Path)
		if err != nil {
			return nil, err
		}
		files[strings.TrimPrefix(e.Path, dir+"/")] = data
	}
	return files, nil
}

// regular reports whether a tree entry is a regular file: a blob of mode 100644 or 100755.
func regular(e git.TreeEntry) bool {
	return e.Type == "blob" && (e.Mode == "100644" || e.Mode == "100755")
}

// Base gives the run's records base: the commit it pushed last, else of its
// last read.
func (s *recordsStore) Base() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pushed != "" {
		return s.pushed
	}
	return s.read
}

// LastChange gives the last commit of the records base that changed a path of paths (task T-36xw, #471).
func (s *recordsStore) LastChange(ctx context.Context, paths ...string) (string, error) {
	return git.LastChange(s.dir, s.Base(), paths...)
}

// Started gives the time of Start: the committer time of the root commit of
// the run's records base, the first commit of the records branch
// (docs/spec/intake.md, The budget of Intake; task T-ybd5).
func (s *recordsStore) Started(ctx context.Context) (time.Time, error) {
	base := s.Base()
	roots, err := git.RootCommits(s.dir, base)
	if err != nil {
		return time.Time{}, err
	}
	if len(roots) != 1 {
		return time.Time{}, fmt.Errorf("the records base %s has %d root commits; want 1", base, len(roots))
	}
	return git.CommitTime(s.dir, roots[0])
}

// PushHead pushes sha of the run's clone to branch of the target, never with
// force; git's refusal (its code 1) is errPushRefused. The store's dir must be
// the clone that Sessions.Clone names, into which the end fetched the head: a
// commit that dir lacks is git's code 1 too, a false refusal.
func (s *recordsStore) PushHead(ctx context.Context, sha, branch string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, err := s.auth(ctx)
	if err != nil {
		return err
	}
	if err := git.Push(s.dir, s.url, sha, branch, a); err != nil {
		var failed *git.FailedError
		if errors.As(err, &failed) && failed.Code == 1 {
			return errPushRefused
		}
		return err
	}
	return nil
}

// BranchHead gives the commit that branch of the target holds, found false for a
// branch it does not have: the same lock and the same token as PushHead (task
// T-syxm, packages.md, The calls of M2e).
func (s *recordsStore) BranchHead(ctx context.Context, branch string) (string, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, err := s.auth(ctx)
	if err != nil {
		return "", false, err
	}
	return git.RemoteHead(s.url, branch, a)
}

// FetchBranch fetches branch of the target into the same ref of the run's clone and gives its commit (task T-89xf).
func (s *recordsStore) FetchBranch(ctx context.Context, branch string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, err := s.auth(ctx)
	if err != nil {
		return "", err
	}
	ref := "refs/heads/" + branch
	if err := git.Fetch(s.dir, s.url, ref, a); err != nil {
		return "", err
	}
	return git.RevParse(s.dir, ref)
}

// Commit makes a records commit of files on the base and pushes it.
func (s *recordsStore) Commit(ctx context.Context, files map[string][]byte, message string) error {
	return s.commit(ctx, files, message, "")
}

// Remove makes a records commit that removes the files at paths, each a path of
// the records, and pushes it: the same lock, work tree and push as a commit that
// writes. A path that is not a path of the records, or not a regular file of the
// base's tree (a path below a link of it too), is an error, and no commit is made
// (task T-eahr).
func (s *recordsStore) Remove(ctx context.Context, paths []string, message string) error {
	for _, p := range paths {
		if !fs.ValidPath(p) || p == "." {
			return fmt.Errorf("%q is not a path of the records", p)
		}
	}
	return s.commitFiles(ctx, nil, paths, message, "")
}

// commit makes a records commit of files on the base and pushes it; orphan is
// the default branch of the first records commit, which has no parent.
// git's refusal of the push is ErrRefused.
func (s *recordsStore) commit(ctx context.Context, files map[string][]byte, message, orphan string) error {
	return s.commitFiles(ctx, files, nil, message, orphan)
}

// commitFiles is commit that also removes the files at remove.
func (s *recordsStore) commitFiles(ctx context.Context, files map[string][]byte, remove []string, message, orphan string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	base := s.pushed
	if base == "" {
		base = s.read
	}
	if orphan != "" {
		base = "refs/heads/" + orphan
	}
	if base == "" {
		return errors.New("no records commit to build on")
	}
	tmp, err := os.MkdirTemp("", "layup-run-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	tree := filepath.Join(tmp, "records")
	if err := git.WorktreeAdd(s.dir, tree, base); err != nil {
		return err
	}
	defer git.WorktreeRemove(s.dir, tree)
	if orphan != "" {
		if err := git.SwitchOrphan(tree, "layup-records"); err != nil {
			return err
		}
	}
	for p, data := range files {
		path := filepath.Join(tree, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return err
		}
	}
	for _, p := range remove {
		// The base's tree holds p as a regular file: a path below a link of the base has no entry, and
		// would take the removal out of the work tree.
		if entries, err := git.LsTree(s.dir, base, p); err != nil {
			return err
		} else if len(entries) != 1 || !regular(entries[0]) {
			return fmt.Errorf("the records base has no regular file %s to remove", p)
		}
		if err := os.Remove(filepath.Join(tree, filepath.FromSlash(p))); err != nil {
			return fmt.Errorf("the records base has no file %s to remove: %w", p, err)
		}
	}
	if err := git.Add(tree); err != nil {
		return err
	}
	if err := git.Commit(tree, message, s.who()); err != nil {
		return err
	}
	id, err := git.RevParse(tree, "HEAD")
	if err != nil {
		return err
	}
	a, err := s.auth(ctx)
	if err != nil {
		return err
	}
	if err := git.Push(s.dir, s.url, id, "layup-records", a); err != nil {
		var failed *git.FailedError
		if errors.As(err, &failed) && failed.Code == 1 {
			return ErrRefused
		}
		return err
	}
	s.pushed = id
	return nil
}

// The time form of the records (docs/spec/README.md, the type time).
const timeForm = "2006-01-02T15:04:05Z"

func stamp(t time.Time) string { return t.UTC().Format(timeForm) }

func leaseTable(r LeaseRow) ([]byte, error) {
	return table(records.LeaseSchema, [][]string{{r.Run, r.Host, r.Version, stamp(r.Started), strconv.Itoa(r.Heartbeat), r.State}})
}

func leaseRow(r []string) (LeaseRow, error) {
	started, err := time.Parse(timeForm, r[3])
	if err != nil {
		return LeaseRow{}, err
	}
	beat, err := strconv.Atoi(r[4])
	if err != nil {
		return LeaseRow{}, err
	}
	return LeaseRow{Run: r[0], Host: r[1], Version: r[2], Started: started, Heartbeat: beat, State: r[5]}, nil
}

// table gives the bytes of a record of schema s.
func table(s tsv.Schema, rows [][]string) ([]byte, error) {
	var b bytes.Buffer
	err := tsv.Write(&b, s, rows)
	return b.Bytes(), err
}
