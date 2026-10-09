package session

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// A Refusal is a start or a result that the session refuses: the reason, the
// word of docs/spec/session.md (rules, context, prompt, branch, and the
// reasons of later rows), and its value, what the event refused writes after
// one space (row 36a writes it), or "" for none (row 33b, task T-6sbe, #160).
// A branch refusal has no value of the event (docs/spec/records.md, events):
// its Value is the error's text, for the log, and the event writes the word
// branch alone.
type Refusal struct{ Reason, Value string }

func (r Refusal) Error() string {
	if r.Value == "" {
		return "refused: " + r.Reason
	}
	return "refused: " + r.Reason + " " + r.Value
}

// NewID draws a session ID: S- and 8 random lowercase hexadecimal characters.
func NewID() (string, error) {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "S-" + hex.EncodeToString(b), nil
}

// CheckContext gives the estimate of the prompt, its bytes over four rounded
// up, and refuses an estimate over the model's context size: the value is the
// estimate, then the size.
func CheckContext(prompt []byte, context int) (int, error) {
	estimate := (len(prompt) + 3) / 4
	if estimate > context {
		return estimate, Refusal{"context", strconv.Itoa(estimate) + " " + strconv.Itoa(context)}
	}
	return estimate, nil
}

// maxArg is the largest prompt that one argument holds: Linux's limit of one
// argument is 131,072 bytes with its final zero byte.
const maxArg = 131071

// CheckPromptSize refuses, for a row whose prompt is arg, a prompt over
// 131,071 bytes; file and stdin have no limit here.
func CheckPromptSize(prompt []byte, mode string) error {
	if mode == "arg" && len(prompt) > maxArg {
		return Refusal{"prompt", strconv.Itoa(len(prompt))}
	}
	return nil
}

// CheckRuleFiles looks in each directory from the session directory up to
// the root for each rule-file name of the register row: one found refuses the
// start, with its path, as the harness would load it. It gives each policy
// path of the row that exists, for the start row (L-A5).
func CheckRuleFiles(d Dir, rules, policy []string) ([]string, error) {
	for dir := d.Root; ; dir = filepath.Dir(dir) {
		for _, name := range rules {
			p := filepath.Join(dir, name)
			if _, err := os.Lstat(p); err == nil {
				return nil, Refusal{"rules", p}
			}
		}
		if dir == filepath.Dir(dir) {
			break
		}
	}
	var found []string
	for _, p := range policy {
		if _, err := os.Stat(p); err == nil {
			found = append(found, p)
		}
	}
	return found, nil
}

// refFiles is what headOf reads: os, or a stand-in of the tests.
type refFiles interface {
	Lstat(name string) (fs.FileInfo, error)
	ReadFile(name string) ([]byte, error)
}

type osFiles struct{}

func (osFiles) Lstat(name string) (fs.FileInfo, error) { return os.Lstat(name) }
func (osFiles) ReadFile(name string) ([]byte, error)   { return os.ReadFile(name) }

// HeadOf gives the SHA of refs/heads/task/<task>/<attempt> from the files of
// repo/.git, with no command run in repo/ (docs/spec/session.md, The fetch by
// SHA).
func HeadOf(repo, task string, attempt int) (string, error) {
	return headOf(osFiles{}, repo, task, attempt)
}

var (
	looseForm = regexp.MustCompile(`^[0-9a-f]{40}\n$`)
	shaForm   = regexp.MustCompile(`^[0-9a-f]{40}$`)
)

// headOf reads the head: repo/.git a directory; the loose ref, a regular file
// of 40 lowercase hexadecimal characters and a line feed (a link, a symbolic
// ref or any other content refuses, and packed-refs is not read then); else
// the line of packed-refs whose ref field, the text after the first space, is
// the ref, and whose first field is 40 lowercase hexadecimal characters. Each
// refusal is Refusal{branch}.
func headOf(f refFiles, repo, task string, attempt int) (string, error) {
	gitDir := filepath.Join(repo, ".git")
	ref := "refs/heads/task/" + task + "/" + strconv.Itoa(attempt)
	if info, err := f.Lstat(gitDir); err != nil || !info.IsDir() {
		return "", Refusal{"branch", "repo/.git is not a directory"}
	}
	loose := filepath.Join(gitDir, filepath.FromSlash(ref))
	info, err := f.Lstat(loose)
	switch {
	case err == nil && !info.Mode().IsRegular():
		return "", Refusal{"branch", ref + " is not a regular file"}
	case err == nil:
		data, err := f.ReadFile(loose)
		if err != nil || !looseForm.Match(data) {
			return "", Refusal{"branch", ref + " is not a SHA and a line feed"}
		}
		return strings.TrimSuffix(string(data), "\n"), nil
	case !errors.Is(err, fs.ErrNotExist):
		return "", Refusal{"branch", err.Error()}
	}
	data, err := f.ReadFile(filepath.Join(gitDir, "packed-refs"))
	if err != nil {
		return "", Refusal{"branch", "no ref " + ref}
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		sha, name, ok := strings.Cut(sc.Text(), " ")
		if !ok || name != ref {
			continue
		}
		if !shaForm.MatchString(sha) {
			return "", Refusal{"branch", "the line of " + ref + " in packed-refs is of another form"}
		}
		return sha, nil
	}
	return "", Refusal{"branch", "no ref " + ref}
}
