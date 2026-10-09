package session

// This file holds the process of a session: the version check, the words of
// the command, the process under its limits and its stop, and the one call
// that orders the steps of a session (docs/spec/session.md, REQ-013 — A role
// session, The start of a session, steps 2 and 5, and The limit of a
// session; row 34a of the plan, task T-5pxd, #161).

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// Version runs the version command of a harness register row, its words split
// at each space, with env in dir, and gives the first line of its standard
// output with no white space at either end. A command that exits non-zero or
// whose first line is empty is Refusal{version}. env is the caller's: the
// session's environment with no credential (step 2), which rows 36a and 38
// build, as a credential of file: is already in home/ once Make has run.
func Version(ctx context.Context, command string, env []string, dir string) (string, error) {
	words := strings.Split(command, " ")
	cmd := exec.CommandContext(ctx, words[0], words[1:]...)
	cmd.Env, cmd.Dir = env, dir
	out, err := cmd.Output()
	line, _, _ := bytes.Cut(out, []byte("\n"))
	v := strings.TrimSpace(string(line))
	if err != nil || v == "" {
		return "", Refusal{Reason: "version"}
	}
	return v, nil
}

// Words gives the words of a harness register row's command, split at each
// space, with {model}, {cap} and {prompt} replaced in each word in one pass,
// so a value that holds a placeholder or a space stays as it is, in its word.
func Words(command, model, cap, prompt string) []string {
	r := strings.NewReplacer("{model}", model, "{cap}", cap, "{prompt}", prompt)
	words := strings.Split(command, " ")
	for i, w := range words {
		words[i] = r.Replace(w)
	}
	return words
}

// Prompt gives the value of {prompt} for the row's prompt: the text, as one
// word, for arg; the path of prompt.md for file and for stdin.
func Prompt(mode string, d Dir, text []byte) string {
	if mode == "arg" {
		return string(text)
	}
	return d.Prompt
}

// Spec is what Process starts: the words in Dir with Env, the file Stdin as
// the standard input ("" for an empty input), the two outputs to the files
// Stdout and Stderr, the limits and the two waits of the stop.
type Spec struct {
	Words, Env                    []string
	Dir, Stdin, Stdout, Stderr    string
	Wall, IntWait, TermWait       time.Duration
	StdoutCap, StderrCap, LineCap int64
}

// NewSpec gives the Spec of a session's process in d: in repo/, prompt.md as
// the standard input for the mode stdin, the files stdout and stderr, wall
// from the start, the outputs cut at 64 MiB and a line of stdout at 8 MiB,
// and 10 s before each later signal of the stop.
func NewSpec(d Dir, words, env []string, mode string, wall time.Duration) Spec {
	s := Spec{Words: words, Env: env, Dir: d.Repo, Stdout: filepath.Join(d.Root, "stdout"), Stderr: filepath.Join(d.Root, "stderr"),
		Wall: wall, IntWait: 10 * time.Second, TermWait: 10 * time.Second, StdoutCap: 64 << 20, StderrCap: 64 << 20, LineCap: 8 << 20}
	if mode == "stdin" {
		s.Stdin = d.Prompt
	}
	return s
}

// Run is what a process did: its start, the first byte of its stdout (zero
// for none), its end, its exit code (-1 for a signal) or the signal that
// ended it, what stopped it ("wall", "output" or "" for none), and the
// signals that Process sent its group, in order.
type Run struct {
	Start, First, End time.Time
	Exit              int
	Signal            syscall.Signal
	StoppedBy         string
	Sent              []syscall.Signal
}

// ErrStart is a program that is missing or cannot run.
var ErrStart = errors.New("the program cannot start")

// Process starts s in a process group of its own and waits for its end: the
// leader has exited and both outputs are closed, so a child that holds an
// output runs under the same limits. At s.Wall from the start, or when an
// output passes its cap, comes the stop: SIGINT to the group, SIGTERM
// s.IntWait later, SIGKILL s.TermWait after that, each only while the
// process has not ended. An output is cut at its cap, and the rest is read
// and dropped.
func Process(s Spec) (Run, error) {
	files := make([]*os.File, 0, 6)
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()
	open := func(path string, flag int) (*os.File, error) {
		f, err := os.OpenFile(path, flag, 0o600)
		if err == nil {
			files = append(files, f)
		}
		return f, err
	}
	stdout, err := open(s.Stdout, os.O_WRONLY|os.O_CREATE|os.O_EXCL)
	if err != nil {
		return Run{}, err
	}
	stderr, err := open(s.Stderr, os.O_WRONLY|os.O_CREATE|os.O_EXCL)
	if err != nil {
		return Run{}, err
	}
	cmd := exec.Command(s.Words[0], s.Words[1:]...)
	cmd.Dir, cmd.Env = s.Dir, s.Env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if s.Stdin != "" {
		if cmd.Stdin, err = open(s.Stdin, os.O_RDONLY); err != nil {
			return Run{}, err
		}
	}
	outR, outW, err := os.Pipe()
	if err != nil {
		return Run{}, err
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		outR.Close()
		outW.Close()
		return Run{}, err
	}
	files = append(files, outR, errR)
	cmd.Stdout, cmd.Stderr = outW, errW
	r := Run{Start: time.Now()}
	err = cmd.Start()
	outW.Close()
	errW.Close()
	if err != nil {
		return r, fmt.Errorf("%w: %v", ErrStart, err)
	}

	var (
		mu       sync.Mutex
		writeErr error
		over     = make(chan struct{})
		overOnce sync.Once
		readers  sync.WaitGroup
	)
	fail := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if writeErr == nil {
			writeErr = err
		}
	}
	first := func() {
		mu.Lock()
		defer mu.Unlock()
		if r.First.IsZero() {
			r.First = time.Now()
		}
	}
	readers.Add(2)
	go func() {
		defer readers.Done()
		fail(capped(outR, stdout, s.StdoutCap, s.LineCap, first, func() { overOnce.Do(func() { close(over) }) }))
	}()
	go func() {
		defer readers.Done()
		fail(capped(errR, stderr, s.StderrCap, 0, nil, func() { overOnce.Do(func() { close(over) }) }))
	}()
	exited := make(chan struct{})
	go func() {
		exitedUnreaped(cmd.Process.Pid)
		close(exited)
	}()
	ended := make(chan struct{})
	go func() {
		<-exited
		readers.Wait()
		close(ended)
	}()

	wall := time.NewTimer(s.Wall)
	defer wall.Stop()
	select {
	case <-ended:
	case <-wall.C:
		r.StoppedBy = "wall"
	case <-over:
		r.StoppedBy = "output"
	}
	if r.StoppedBy != "" {
		stop(cmd.Process.Pid, &r, ended, exited, []time.Duration{s.IntWait, s.TermWait}, outR, errR)
	}
	<-ended
	mu.Lock()
	defer mu.Unlock()
	r.End = time.Now()
	cmd.Wait() // reaps the leader only now; its exit is read from cmd.ProcessState
	r.Exit = cmd.ProcessState.ExitCode()
	if ws, ok := cmd.ProcessState.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		r.Signal = ws.Signal()
	}
	return r, writeErr
}

// exitedUnreaped waits until the process pid has exited and leaves it
// unreaped (waitid with WNOWAIT), so its PID, the ID of its group, stays taken
// and a signal to the group reaches no other process until Process reaps it.
func exitedUnreaped(pid int) {
	const pPID, wExited, wNoWait = 1, 4, 0x1000000 // P_PID, WEXITED, WNOWAIT of Linux
	var info [128]byte                             // a siginfo_t, which the call fills and the code does not read
	for {
		_, _, e := syscall.Syscall6(syscall.SYS_WAITID, pPID, uintptr(pid), uintptr(unsafe.Pointer(&info[0])), wExited|wNoWait, 0, 0)
		if e != syscall.EINTR {
			return
		}
	}
}

// stop sends SIGINT, SIGTERM and SIGKILL to the group pgid, each after its
// wait while the process has not ended. After SIGKILL it closes the two
// outputs once the leader is reaped, so a program that left the group and
// holds an output cannot hold the stop.
func stop(pgid int, r *Run, ended, exited <-chan struct{}, waits []time.Duration, outputs ...*os.File) {
	for i, sig := range []syscall.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL} {
		syscall.Kill(-pgid, sig)
		r.Sent = append(r.Sent, sig)
		if i < len(waits) {
			select {
			case <-ended:
				return
			case <-time.After(waits[i]):
			}
		}
	}
	<-exited
	for _, f := range outputs {
		f.Close()
	}
}

// capped copies from to to until from ends, at most max bytes, and, when line
// is not 0, at most line bytes in a line before its line feed. first is
// called at the first byte read, and over when a cap is passed; the rest is
// read and dropped.
func capped(from io.Reader, to io.Writer, max, line int64, first, over func()) error {
	var written, inLine int64
	cut := false
	buf := make([]byte, 64<<10)
	for {
		n, err := from.Read(buf)
		if n > 0 && first != nil {
			first()
			first = nil
		}
		if n > 0 && !cut {
			keep := min(int64(n), max-written)
			if line > 0 {
				for i, b := range buf[:keep] {
					if b == '\n' {
						inLine = 0
						continue
					}
					if inLine++; inLine > line {
						keep = int64(i)
						break
					}
				}
			}
			if _, werr := to.Write(buf[:keep]); werr != nil {
				return werr
			}
			written += keep
			if keep < int64(n) {
				cut = true
				over()
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
				return nil
			}
			return err
		}
	}
}

// Steps is the hooks of the one call of a session, which its caller gives (a
// task session from row 36a, a probe from row 38): the session ID; the attempt
// check (nil for a probe); the version, context, prompt-size and rule-file
// checks; the start row; the process; its end, with what Process gave; and
// the checks before a push, the push and the bind (nil for a probe; rows 37a
// and 37b).
type Steps struct {
	ID                                       func() (string, error)
	Attempt, Version, Context, Prompt, Rules func(id string) error
	StartRow                                 func(id string) error
	Process                                  func(id string) (Run, error)
	End                                      func(id string, r Run, err error) error
	Push                                     func(id string, r Run) error
}

// Call runs the steps of a session in their order (docs/spec/session.md,
// REQ-013 — A role session): a step that refuses ends the call with its
// error before the next step runs, so no process starts after a refusal. A
// process that cannot start is still ended: its error goes to End. It gives
// the session ID, or "" when none was drawn.
func Call(s Steps) (string, error) {
	id, err := s.ID()
	if err != nil {
		return "", err
	}
	checks := []func(string) error{s.Version, s.Context, s.Prompt, s.Rules, s.StartRow}
	if s.Attempt != nil { // a task session's; a probe has no attempt
		checks = append([]func(string) error{s.Attempt}, checks...)
	}
	for _, check := range checks {
		if err := check(id); err != nil {
			return id, err
		}
	}
	r, err := s.Process(id)
	if err := s.End(id, r, err); err != nil {
		return id, err
	}
	if s.Push == nil {
		return id, nil
	}
	return id, s.Push(id, r)
}
