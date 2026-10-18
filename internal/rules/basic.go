package rules

// The basic rule of layup/rules (docs/spec/loop.md, layup/rules; row 63a of
// the plan, #352): a pull request that changes a path of the rule-path
// register fails unless the change is a rendered record, a file whose bytes
// equal code's rendering of the named copy. Judge is pure; Basic reads the
// repository through a Repo, from the sources that loop.md names: the register
// at a commit the caller names, never the head's; the paths that differ
// between the base and the head; what the head holds at a record's path; and,
// for a PRD, the file at the commit it is rendered from. The caller owns the
// copy: it gives the parsed values, and reads and checks nothing of the copy
// itself.

import (
	"bytes"
	"fmt"
	"path"
	"strings"

	"github.com/pharzam/layup/internal/git"
	"github.com/pharzam/layup/internal/spec"
)

// State is the state of a status that layup/rules gives.
type State string

// The three states this check gives; pending is the caller's.
const (
	Success State = "success"
	Failure State = "failure"
	Error   State = "error"
)

// The closed set of reasons.
const (
	ReasonNoRulePath      = "no rule path"
	ReasonRenderedRecord  = "rendered record"
	ReasonRulePath        = "rule path"
	ReasonNotEqual        = "not equal"
	ReasonNotARegularFile = "not a regular file"
	ReasonNoRendering     = "no rendering"
	ReasonRegister        = "register"
	ReasonInput           = "input"
	ReasonBase            = "base"
	ReasonGit             = "git"
)

// Result is the judgement of one diff.
type Result struct {
	State        State
	Reason, Path string
}

// Kind is the kind of a rendered record.
type Kind string

// The two kinds of rendered record.
const (
	PRD   Kind = "prd"
	Facts Kind = "facts"
)

// Record names one rendered record.
type Record struct {
	Path    string
	Kind    Kind
	Answer  spec.Answer
	Version [][]string
	Texts   map[string]string
}

// Input is what the caller names.
type Input struct {
	Base, Head, RegisterRev, PreRev string
	Records                         []Record
}

// Touched is one path that differs between the base and the head.
type Touched struct {
	Path       string
	Present    bool
	Mode, Type string
	Data       []byte
	Guardrails Change
}

// Repo is the six calls of git that Basic needs.
type Repo interface {
	RevParse(rev string) (string, error)
	Show(rev, path string) ([]byte, error)
	LsTree(rev, path string) ([]git.TreeEntry, error)
	DiffNames(base, head string) ([]string, error)
	DiffFile(base, head, path string) ([]byte, error)
	IsAncestor(base, head string) (bool, error)
}

type gitRepo struct{ dir string }

// Git is the Repo of the repository at dir, by internal/git.
func Git(dir string) Repo { return gitRepo{dir} }

func (g gitRepo) RevParse(rev string) (string, error)   { return git.RevParse(g.dir, rev+"^{commit}") }
func (g gitRepo) Show(rev, path string) ([]byte, error) { return git.Show(g.dir, rev, path) }
func (g gitRepo) LsTree(rev, path string) ([]git.TreeEntry, error) {
	return git.LsTree(g.dir, rev, path)
}
func (g gitRepo) DiffNames(base, head string) ([]string, error) {
	return git.DiffNames(g.dir, base, head)
}
func (g gitRepo) DiffFile(base, head, path string) ([]byte, error) {
	return git.DiffFile(g.dir, base, head, path)
}
func (g gitRepo) IsAncestor(base, head string) (bool, error) {
	return git.IsAncestor(g.dir, base, head)
}

const registerPath = "rule-paths.tsv"

// Render is code's rendering of a record: the PRD from its pre-image and the
// answer of the copy, or the numbered facts record from the version of the
// inventory and the text of each fact. It is the one rendering that Judge
// compares and that the handoff item rendering-equal reads. A kind that is
// none of the two gives one finding.
func Render(rec Record, pre []byte) ([]byte, []spec.Finding) {
	switch rec.Kind {
	case PRD:
		return spec.RenderPRD(pre, rec.Answer)
	case Facts:
		return spec.RenderFacts(rec.Version, rec.Texts)
	}
	return nil, []spec.Finding{{Rule: "kind", Detail: fmt.Sprintf("the kind %q is neither prd nor facts", rec.Kind)}}
}

// ValidPath reports whether p is a path of a tree that git can name by
// itself: relative, in its clean form, with no part that leaves the tree or
// reaches into .git. git show REV:../x reads relative to the working
// directory, so a spelling that a reader finds clean can still leave the tree.
func ValidPath(p string) bool {
	if p == "" || p != path.Clean(p) || strings.HasPrefix(p, "/") || strings.HasPrefix(p, ":") || strings.ContainsAny(p, "\\\x00") {
		return false
	}
	for _, part := range strings.Split(p, "/") {
		if part == "." || part == ".." || strings.EqualFold(part, ".git") {
			return false
		}
	}
	return true
}

func failure(reason, p string) Result { return Result{Failure, reason, p} }
func errorOf(reason, p string) Result { return Result{Error, reason, p} }

// Judge judges the paths that differ between a base and a head, in the order
// given (git's), against the register. Each path is asked of Check, which is
// the refusal of layer 2 before a push, so both layers read the register by
// one rule, the kind .sh and the exception of the guardrails included; a
// refusal of Check, whichever of its two reasons, is a rule path here. A rule
// path passes only as a rendered record: the caller named it, the head holds a
// regular file there (mode 100644), and its bytes are exactly Render's. pre
// holds, for each PRD record, the file that it is rendered from. A bad record
// is an error, never a pass, whether or not its path changed.
func Judge(reg []Entry, touched []Touched, records []Record, pre map[string][]byte) Result {
	byPath := map[string]Record{}
	for _, rec := range records {
		if _, hasPre := pre[rec.Path]; !recordInputOK(rec) || (rec.Kind == PRD && !hasPre) {
			return errorOf(ReasonInput, rec.Path)
		}
		if _, twice := byPath[rec.Path]; twice {
			return errorOf(ReasonInput, rec.Path)
		}
		byPath[rec.Path] = rec
	}
	rendered := 0
	for _, t := range touched {
		if reason, _ := Check([]string{t.Path}, reg, t.Guardrails); reason == "" {
			continue
		}
		rec, named := byPath[t.Path]
		if !named {
			return failure(ReasonRulePath, t.Path)
		}
		if reason := RecordRendered(rec, t, pre[t.Path]); reason != "" {
			return failure(reason, t.Path)
		}
		rendered++
	}
	if rendered > 0 {
		return Result{Success, ReasonRenderedRecord, ""}
	}
	return Result{Success, ReasonNoRulePath, ""}
}

// recordInputOK reports whether a record is one that can be judged: a path of
// a tree that is not the guardrails file, and a kind that is a PRD or facts.
func recordInputOK(rec Record) bool {
	return ValidPath(rec.Path) && rec.Path != GuardrailsPath && (rec.Kind == PRD || rec.Kind == Facts)
}

// RecordRendered is the per-record rule that Judge applies to a path the head
// changes and that the handoff item rendering-equal applies to a path the head
// names (task T-yeav, #404): the head holds a regular file there (mode 100644)
// whose bytes are exactly Render's, with no finding. It gives the closed reason
// of the failure, "" for none. pre is the file a PRD is rendered from.
func RecordRendered(rec Record, t Touched, pre []byte) string {
	switch {
	case !t.Present:
		return ReasonNotEqual
	case t.Type != "blob" || t.Mode != "100644":
		return ReasonNotARegularFile
	}
	want, findings := Render(rec, pre)
	switch {
	case len(findings) > 0:
		return ReasonNoRendering
	case !bytes.Equal(t.Data, want):
		return ReasonNotEqual
	}
	return ""
}

// fileAt gives what rev holds at p: whether there is an entry, its mode and
// type, and the bytes of a regular file.
func fileAt(r Repo, rev, p string) (Touched, error) {
	e, ok, err := entryAt(r, rev, p)
	if err != nil {
		return Touched{}, err
	}
	t := Touched{Path: p, Present: ok, Mode: e.Mode, Type: e.Type}
	if ok && e.Type == "blob" && e.Mode == "100644" {
		t.Data, err = r.Show(rev, p)
	}
	return t, err
}

// entryAt gives the entry at exactly p in rev, and whether there is one: the
// tree is asked, so a path below a link or a file has no entry.
func entryAt(r Repo, rev, p string) (git.TreeEntry, bool, error) {
	entries, err := r.LsTree(rev, p)
	if err != nil {
		return git.TreeEntry{}, false, err
	}
	for _, e := range entries {
		if e.Path == p {
			return e, true, nil
		}
	}
	return git.TreeEntry{}, false, nil
}

// Basic judges the change from in.Base to in.Head by the register at
// in.RegisterRev. It gives an error, never a success, for a check that could
// not run: a bad record, a register that cannot be used (none, a file that is
// no regular file, one that ReadRegister refuses, or one with no row, which
// would pass every path), a base that is not an ancestor of the head (the
// two-dot diff is the pull request's only then), a pre-image that is not
// between them, or a failing call of git.
func Basic(r Repo, in Input) Result {
	for _, rec := range in.Records {
		if !ValidPath(rec.Path) {
			return errorOf(ReasonInput, rec.Path)
		}
	}
	base, err := r.RevParse(in.Base)
	if err != nil {
		return errorOf(ReasonGit, "")
	}
	head, err := r.RevParse(in.Head)
	if err != nil {
		return errorOf(ReasonGit, "")
	}
	regRev, err := r.RevParse(in.RegisterRev)
	if err != nil || regRev == head {
		return errorOf(ReasonRegister, "")
	}
	if ok, err := r.IsAncestor(base, head); err != nil {
		return errorOf(ReasonGit, "")
	} else if !ok {
		return errorOf(ReasonBase, "")
	}
	reg, bad := readRegisterAt(r, regRev)
	if bad.State != "" {
		return bad
	}
	pre, bad := preImages(r, in, base, head)
	if bad.State != "" {
		return bad
	}
	names, err := r.DiffNames(base, head)
	if err != nil {
		return errorOf(ReasonGit, "")
	}
	named := map[string]bool{}
	for _, rec := range in.Records {
		named[rec.Path] = true
	}
	touched := make([]Touched, 0, len(names))
	for _, p := range names {
		t := Touched{Path: p}
		if p == GuardrailsPath || named[p] {
			if t, err = fileAt(r, head, p); err != nil {
				return errorOf(ReasonGit, "")
			}
		}
		if p == GuardrailsPath {
			if t.Guardrails.Diff, err = r.DiffFile(base, head, p); err != nil {
				return errorOf(ReasonGit, "")
			}
			// A head that lacks the file, or whose entry is no blob (a
			// gitlink), has an empty Head, which no added line passes.
			if t.Present && t.Type == "blob" {
				if t.Guardrails.Head, err = r.Show(head, p); err != nil {
					return errorOf(ReasonGit, "")
				}
			}
		}
		touched = append(touched, t)
	}
	return Judge(reg, touched, in.Records, pre)
}

// readRegisterAt reads rule-paths.tsv at the root of rev. The zero Result is
// no refusal.
func readRegisterAt(r Repo, rev string) ([]Entry, Result) {
	e, ok, err := entryAt(r, rev, registerPath)
	switch {
	case err != nil:
		return nil, errorOf(ReasonGit, "")
	case !ok || e.Type != "blob" || (e.Mode != "100644" && e.Mode != "100755"):
		return nil, errorOf(ReasonRegister, "")
	}
	data, err := r.Show(rev, registerPath)
	if err != nil {
		return nil, errorOf(ReasonGit, "")
	}
	reg, err := ReadRegister(data)
	if err != nil || len(reg) == 0 {
		return nil, errorOf(ReasonRegister, "")
	}
	return reg, Result{}
}

// preImages reads, for each PRD record, the file at in.PreRev that it is
// rendered from. The commit is the base or lies between the base and the
// head, and is never the head: rendering from the head's own bytes could not
// fail on a byte that is no rendering. An empty base is no link to the base. The
// zero Result is no refusal.
func preImages(r Repo, in Input, base, head string) (map[string][]byte, Result) {
	pre := map[string][]byte{}
	var prds []string
	for _, rec := range in.Records {
		if rec.Kind == PRD {
			prds = append(prds, rec.Path)
		}
	}
	if len(prds) == 0 {
		return pre, Result{}
	}
	if in.PreRev == "" {
		return nil, errorOf(ReasonInput, "")
	}
	rev, err := r.RevParse(in.PreRev)
	if err != nil || rev == head {
		return nil, errorOf(ReasonInput, "")
	}
	for _, link := range [][2]string{{base, rev}, {rev, head}} {
		if link[0] == "" { // a handoff has no base: the pre-image is the head's ancestor only
			continue
		}
		ok, err := r.IsAncestor(link[0], link[1])
		if err != nil {
			return nil, errorOf(ReasonGit, "")
		}
		if !ok {
			return nil, errorOf(ReasonInput, "")
		}
	}
	for _, p := range prds {
		e, ok, err := entryAt(r, rev, p)
		switch {
		case err != nil:
			return nil, errorOf(ReasonGit, "")
		case !ok || e.Type != "blob" || e.Mode != "100644":
			return nil, errorOf(ReasonInput, "")
		}
		if pre[p], err = r.Show(rev, p); err != nil {
			return nil, errorOf(ReasonGit, "")
		}
	}
	return pre, Result{}
}

// failing is a Repo that keeps the first error of any call, so that a failure
// that preImages maps to a reason is still an error to RenderingEqual.
type failing struct {
	Repo
	err error
}

func (f *failing) note(err error) error {
	if f.err == nil {
		f.err = err
	}
	return err
}

func (f *failing) RevParse(rev string) (string, error) {
	s, err := f.Repo.RevParse(rev)
	return s, f.note(err)
}

func (f *failing) Show(rev, p string) ([]byte, error) {
	b, err := f.Repo.Show(rev, p)
	return b, f.note(err)
}

func (f *failing) LsTree(rev, p string) ([]git.TreeEntry, error) {
	e, err := f.Repo.LsTree(rev, p)
	return e, f.note(err)
}

func (f *failing) IsAncestor(base, head string) (bool, error) {
	ok, err := f.Repo.IsAncestor(base, head)
	return ok, f.note(err)
}

// RenderingEqual reports whether the head holds, at each of the paths, the
// rendering of its record: the handoff item rendering-equal (loop.md, The items
// a caller computes; task T-yeav, #404). The paths and the records match one to
// one. A PRD record is rendered from the file at preRev, a revision that
// resolves, is not the head and is an ancestor of it; the bytes are never the
// caller's. Every refusal of the input or of the tree is (false, nil), the
// static ones before any call of git; a failed call of git, the resolution of
// the head or of preRev included, is (false, err).
func RenderingEqual(r Repo, head, preRev string, paths []string, records []Record) (bool, error) {
	if len(paths) == 0 || len(paths) != len(records) {
		return false, nil
	}
	byPath := map[string]Record{}
	for _, rec := range records {
		if !recordInputOK(rec) {
			return false, nil
		}
		byPath[rec.Path] = rec
	}
	seen := map[string]bool{}
	for _, p := range paths {
		if _, ok := byPath[p]; !ok || seen[p] {
			return false, nil
		}
		seen[p] = true
	}
	g := &failing{Repo: r}
	h, err := g.RevParse(head)
	if err != nil {
		return false, err
	}
	pre, bad := preImages(g, Input{PreRev: preRev, Records: records}, "", h)
	if g.err != nil {
		return false, g.err
	}
	if bad.State != "" {
		return false, nil
	}
	for _, p := range paths {
		t, err := fileAt(g, h, p)
		if err != nil {
			return false, err
		}
		if RecordRendered(byPath[p], t, pre[p]) != "" {
			return false, nil
		}
	}
	return true, nil
}
