// Package session is a role session of milestone M2b (docs/spec/session.md,
// REQ-013 — A role session). This file holds the session's files: its
// directory under the host directory, its clone, home and scratch directory,
// and the environment of its process (row 33a of the plan, task T-vxdg,
// #159). It may not import internal/route, so it takes the values of a
// harness register row as its own type, which internal/run fills.
package session

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pharzam/layup/internal/git"
)

// Harness holds the values of a harness register row that a session's files
// need: the credential file and how the harness takes it, and the fixed
// variables (docs/spec/records.md, the block harness-register).
type Harness struct {
	Credential, CredentialTo string
	Vars                     []string
}

// Dir is the paths of a session's directory, DIR/sessions/<session>/.
type Dir struct{ Root, Repo, Home, Tmp, Prompt, Result string }

func dirOf(host, id string) Dir {
	root := filepath.Join(host, "sessions", id)
	return Dir{Root: root, Repo: filepath.Join(root, "repo"), Home: filepath.Join(root, "home"), Tmp: filepath.Join(root, "tmp"),
		Prompt: filepath.Join(root, "prompt.md"), Result: filepath.Join(root, "result")}
}

// Make makes the directory of session id for target (OWNER/NAME) under the
// host directory: the file target first; repo/, a clone of the one branch of
// the run's clone with no remote (git.CloneLocal), on task/<task>/<attempt>
// at base (git.SwitchCreate); home/ with .gitconfig and, for a credential of
// the route file:, its copy, mode 0600; tmp/ and result/ empty; prompt.md. A
// directory of that ID that exists is an error: a session ID is never reused.
func Make(host, target, id, runClone, branch, task string, attempt int, base string, prompt []byte, h Harness) (Dir, error) {
	host, err := filepath.Abs(host)
	if err != nil {
		return Dir{}, err
	}
	d := dirOf(host, id)
	if err := os.MkdirAll(filepath.Dir(d.Root), 0o700); err != nil {
		return Dir{}, err
	}
	if err := os.Mkdir(d.Root, 0o700); err != nil {
		return Dir{}, fmt.Errorf("the session directory %s: %w", d.Root, err)
	}
	if err := os.WriteFile(filepath.Join(d.Root, "target"), []byte(target+"\n"), 0o600); err != nil {
		return Dir{}, err
	}
	if err := git.CloneLocal(runClone, d.Repo, branch); err != nil {
		return Dir{}, err
	}
	if err := git.SwitchCreate(d.Repo, "task/"+task+"/"+strconv.Itoa(attempt), base); err != nil {
		return Dir{}, err
	}
	for _, dir := range []string{d.Home, d.Tmp, d.Result} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			return Dir{}, err
		}
	}
	cfg := "[user]\n\tname = layup session " + id + "\n\temail = " + id + "@sessions.layup.invalid\n"
	if err := os.WriteFile(filepath.Join(d.Home, ".gitconfig"), []byte(cfg), 0o600); err != nil {
		return Dir{}, err
	}
	if path, ok := strings.CutPrefix(h.CredentialTo, "file:"); ok {
		data, err := os.ReadFile(h.Credential)
		if err != nil {
			return Dir{}, fmt.Errorf("the credential %s: %w", h.Credential, err)
		}
		dst := filepath.Join(d.Home, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
			return Dir{}, err
		}
		if err := os.WriteFile(dst, data, 0o600); err != nil {
			return Dir{}, err
		}
	}
	return d, os.WriteFile(d.Prompt, prompt, 0o600)
}

// Sweep removes, before the next session of target starts, each session
// directory under the host directory that a stopped run of that target left:
// those whose file target is target, and those with no file target, which
// only a Make that stopped before its first write leaves. A directory of
// another target is kept: its own lease guards it.
func Sweep(host, target string) error {
	dir := filepath.Join(host, "sessions")
	list, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, e := range list {
		t, err := os.ReadFile(filepath.Join(dir, e.Name(), "target"))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err == nil && strings.TrimSuffix(string(t), "\n") != target {
			continue
		}
		if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// Environ gives the environment of a session's process, reading the
// credential file of the route var: from the disk.
func Environ(d Dir, h Harness) ([]string, error) { return environ(d.Home, d.Tmp, h, os.ReadFile) }

// environ builds the environment of a session's process: the named list and
// no other variable (docs/spec/session.md, The environment and the harness
// credential): PATH of the host, the one variable of the host that it reads
// (TestInputRule allows this function), LANG, HOME, TMPDIR,
// GIT_CONFIG_NOSYSTEM, the credential's variable for var:NAME (the file's
// content without its final line feed), and each fixed variable.
func environ(home, tmp string, h Harness, read func(string) ([]byte, error)) ([]string, error) {
	env := []string{"PATH=" + os.Getenv("PATH"), "LANG=C.UTF-8", "HOME=" + home, "TMPDIR=" + tmp, "GIT_CONFIG_NOSYSTEM=1"}
	if name, ok := strings.CutPrefix(h.CredentialTo, "var:"); ok {
		data, err := read(h.Credential)
		if err != nil {
			return nil, fmt.Errorf("the credential %s: %w", h.Credential, err)
		}
		env = append(env, name+"="+strings.TrimSuffix(string(data), "\n"))
	}
	return append(env, h.Vars...), nil
}
