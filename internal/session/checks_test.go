package session

import (
	"errors"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// refused reports a test error unless err is a Refusal with the reason, and,
// when value is not "", that value.
func refused(t *testing.T, name string, err error, reason, value string) {
	t.Helper()
	var r Refusal
	switch {
	case !errors.As(err, &r):
		t.Errorf("%s: %v, want a Refusal %s", name, err, reason)
	case r.Reason != reason || (value != "" && r.Value != value):
		t.Errorf("%s: refused %q %q, want %q %q", name, r.Reason, r.Value, reason, value)
	}
}

func TestNewID(t *testing.T) {
	a, err := NewID()
	b, err2 := NewID()
	if err != nil || err2 != nil || !regexp.MustCompile(`^S-[0-9a-f]{8}$`).MatchString(a) || a == b {
		t.Errorf("NewID: %q %q (%v %v); want two IDs of the form S-xxxxxxxx", a, b, err, err2)
	}
}

func TestCheckContext(t *testing.T) {
	for _, c := range []struct {
		bytes, context, estimate int
		ok                       bool
	}{{400, 100, 100, true}, {401, 100, 101, false}, {5, 2, 2, true}, {5, 1, 2, false}, {0, 0, 0, true}} {
		est, err := CheckContext([]byte(strings.Repeat("x", c.bytes)), c.context)
		if est != c.estimate {
			t.Errorf("%d bytes: estimate %d, want %d", c.bytes, est, c.estimate)
		}
		if c.ok && err != nil {
			t.Errorf("%d bytes in %d: %v, want no refusal", c.bytes, c.context, err)
		}
		if !c.ok {
			refused(t, "an estimate over the context", err, "context", strings.Join([]string{itoa(c.estimate), itoa(c.context)}, " "))
		}
	}
}

func itoa(n int) string { return strconv.Itoa(n) }

func TestCheckPromptSize(t *testing.T) {
	fits, over := []byte(strings.Repeat("x", 131071)), []byte(strings.Repeat("x", 131072))
	if err := CheckPromptSize(fits, "arg"); err != nil {
		t.Errorf("131,071 bytes as one argument: %v", err)
	}
	refused(t, "131,072 bytes as one argument", CheckPromptSize(over, "arg"), "prompt", "131072")
	for _, mode := range []string{"file", "stdin"} {
		if err := CheckPromptSize(over, mode); err != nil {
			t.Errorf("131,072 bytes by %s: %v, want no limit", mode, err)
		}
	}
}

// gitFiles is a stand-in of the gitFiles of a .git directory: their content by
// path, a directory by "<dir>", a link by "<link>".
type gitFiles map[string]string

func (f gitFiles) Lstat(name string) (fs.FileInfo, error) {
	c, ok := f[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	mode := fs.FileMode(0o644)
	switch c {
	case "<dir>":
		mode = fs.ModeDir | 0o755
	case "<link>":
		mode = fs.ModeSymlink | 0o777
	}
	return info{name, mode}, nil
}

func (f gitFiles) ReadFile(name string) ([]byte, error) {
	c, ok := f[name]
	if !ok || c == "<dir>" || c == "<link>" {
		return nil, fs.ErrNotExist
	}
	return []byte(c), nil
}

type info struct {
	name string
	mode fs.FileMode
}

func (i info) Name() string       { return i.name }
func (i info) Size() int64        { return 0 }
func (i info) Mode() fs.FileMode  { return i.mode }
func (i info) ModTime() time.Time { return time.Time{} }
func (i info) IsDir() bool        { return i.mode.IsDir() }
func (i info) Sys() any           { return nil }

const sha = "0123456789abcdef0123456789abcdef01234567"

func TestHeadOfTheFiles(t *testing.T) {
	ref := "/s/repo/.git/refs/heads/task/T-ab12/1"
	packed := "/s/repo/.git/packed-refs"
	git := gitFiles{"/s/repo/.git": "<dir>"}
	with := func(extra gitFiles) gitFiles {
		f := gitFiles{}
		for k, v := range git {
			f[k] = v
		}
		for k, v := range extra {
			f[k] = v
		}
		return f
	}
	for name, f := range map[string]gitFiles{
		"a loose ref":                    with(gitFiles{ref: sha + "\n"}),
		"a packed ref":                   with(gitFiles{packed: "# pack-refs with: peeled fully-peeled sorted\n" + sha + " refs/heads/task/T-ab12/1\n"}),
		"a packed ref beside attempt 10": with(gitFiles{packed: strings.Replace(sha, "0", "f", 1) + " refs/heads/task/T-ab12/10\n" + sha + " refs/heads/task/T-ab12/1\n"}),
	} {
		if got, err := headOf(f, "/s/repo", "T-ab12", 1); err != nil || got != sha {
			t.Errorf("%s: %q, %v; want %s", name, got, err, sha)
		}
	}
	// The refusals whose rule another rule would hide, with the value that
	// names their own rule.
	values := map[string]string{
		".git that is a file":        "repo/.git is not a directory",
		".git that is a link":        "repo/.git is not a directory",
		"a loose ref that is a link": "refs/heads/task/T-ab12/1 is not a regular file",
	}
	for name, f := range map[string]gitFiles{
		".git that is a file":           {"/s/repo/.git": "gitdir: /elsewhere\n"},
		".git that is a link":           {"/s/repo/.git": "<link>"},
		"a loose ref that is a link":    with(gitFiles{ref: "<link>", packed: sha + " refs/heads/task/T-ab12/1\n"}),
		"a symbolic loose ref":          with(gitFiles{ref: "ref: refs/heads/main\n"}),
		"a short SHA":                   with(gitFiles{ref: sha[:39] + "\n"}),
		"no line feed":                  with(gitFiles{ref: sha}),
		"upper case":                    with(gitFiles{ref: strings.ToUpper(sha) + "\n"}),
		"two line feeds":                with(gitFiles{ref: sha + "\n\n"}),
		"no ref":                        with(gitFiles{}),
		"a packed line of another form": with(gitFiles{packed: "abc refs/heads/task/T-ab12/1\n"}),
		"only the ref of attempt 10":    with(gitFiles{packed: sha + " refs/heads/task/T-ab12/10\n"}),
	} {
		_, err := headOf(f, "/s/repo", "T-ab12", 1)
		refused(t, name, err, "branch", values[name])
	}
}
