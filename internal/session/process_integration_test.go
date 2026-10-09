//go:build integration

package session

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// fake writes a fake harness program, a shell script of body, and gives its
// path.
func fake(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "fake")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"+body+"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestVersion(t *testing.T) {
	dir := t.TempDir()
	env := []string{"PATH=" + os.Getenv("PATH"), "V=2.1.295"}
	got, err := Version(context.Background(), fake(t, `printf '  %s (Claude Code) %s \nsecond\n' "$V" "$(basename "$(pwd)")"`)+" --version", env, dir)
	if want := "2.1.295 (Claude Code) " + filepath.Base(dir); err != nil || got != want {
		t.Errorf("the first line, trimmed, with the environment and in dir: %q, %v; want %q", got, err, want)
	}
	for name, body := range map[string]string{
		"a non-zero exit": "echo 1.0; exit 1",
		"no output":       "exit 0",
		"a blank line":    "printf ' \\n1.0\\n'",
	} {
		_, err := Version(context.Background(), fake(t, body), env, dir)
		refused(t, name, err, "version", "")
	}
}

// spec gives a Spec of fake with short waits and small caps, its outputs in a
// directory of the test, and M in its environment.
func spec(t *testing.T, fake, marker string, wall time.Duration) Spec {
	d := t.TempDir()
	return Spec{Words: []string{fake}, Env: []string{"PATH=" + os.Getenv("PATH"), "M=" + marker}, Dir: d,
		Stdout: filepath.Join(d, "stdout"), Stderr: filepath.Join(d, "stderr"), Wall: wall,
		IntWait: 400 * time.Millisecond, TermWait: 400 * time.Millisecond, StdoutCap: 1000, StderrCap: 1000, LineCap: 500}
}

func TestProcessGivesTheInputOfEachMode(t *testing.T) {
	prompt := filepath.Join(t.TempDir(), "prompt.md")
	os.WriteFile(prompt, []byte("the task\n"), 0o600)
	for _, stdin := range []string{"", prompt} {
		got := filepath.Join(t.TempDir(), "got")
		s := spec(t, fake(t, `cat > "$M"`), got, 5*time.Second)
		s.Stdin = stdin
		r, err := Process(s)
		data, _ := os.ReadFile(got)
		want := ""
		if stdin != "" {
			want = "the task\n"
		}
		if err != nil || r.StoppedBy != "" || r.Exit != 0 || string(data) != want {
			t.Errorf("stdin %q: %+v, %v, the input %q; want an end at once, with %q", stdin, r, err, data, want)
		}
	}
}

func TestProcessStopsAtWall(t *testing.T) {
	loop := "while :; do sleep 0.05; done"
	for _, c := range []struct {
		name, body string
		sent       []syscall.Signal
		exit       int
		signal     syscall.Signal
	}{
		{"a fake that ends on SIGINT", `trap 'exit 3' INT; trap 'echo TERM >> "$M"' TERM; ` + loop, []syscall.Signal{syscall.SIGINT}, 3, 0},
		{"a fake that ignores SIGINT", `trap '' INT; trap 'exit 4' TERM; ` + loop, []syscall.Signal{syscall.SIGINT, syscall.SIGTERM}, 4, 0},
		{"a fake that ignores both", `trap '' INT TERM; ` + loop, []syscall.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}, -1, syscall.SIGKILL},
	} {
		marker := filepath.Join(t.TempDir(), "m")
		begin := time.Now()
		r, err := Process(spec(t, fake(t, c.body), marker, 300*time.Millisecond))
		if err != nil || r.StoppedBy != "wall" || !slices.Equal(r.Sent, c.sent) || r.Exit != c.exit || r.Signal != c.signal {
			t.Errorf("%s: %+v, %v; want stopped at wall by %v, exit %d, signal %v", c.name, r, err, c.sent, c.exit, c.signal)
		}
		if took := time.Since(begin); took < 300*time.Millisecond {
			t.Errorf("%s: stopped after %v, before wall", c.name, took)
		}
		time.Sleep(500 * time.Millisecond)
		if data, err := os.ReadFile(marker); err == nil {
			t.Errorf("%s: a signal after its end: %q", c.name, data)
		}
	}
}

func TestProcessStopsTheWholeGroup(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "m")
	pid := filepath.Join(t.TempDir(), "pid")
	// The child runs in the foreground: a job of a shell's & starts with
	// SIGINT ignored, which it cannot trap. The leader's traps do nothing, so
	// the child does not inherit an ignored signal, and the leader waits.
	child := `trap : INT TERM
sh -c 'echo $$ > ` + pid + `; trap "echo INT >> \"$M\"" INT; trap "echo TERM >> \"$M\"" TERM; while :; do sleep 0.05; done'`
	r, err := Process(spec(t, fake(t, child), marker, 300*time.Millisecond))
	if err != nil || r.Signal != syscall.SIGKILL {
		t.Fatalf("%+v, %v; want the leader killed", r, err)
	}
	data, _ := os.ReadFile(marker)
	if got := strings.Fields(string(data)); !slices.Equal(got, []string{"INT", "TERM"}) {
		t.Errorf("the child saw %q, want INT and TERM", got)
	}
	p, _ := os.ReadFile(pid)
	n, _ := strconv.Atoi(strings.TrimSpace(string(p)))
	for i := 0; n > 0 && syscall.Kill(n, 0) == nil; i++ {
		if i == 40 {
			t.Fatalf("the child %d lives after SIGKILL", n)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestProcessStopsAtAnOutputCap(t *testing.T) {
	for _, c := range []struct {
		name, body, file string
		size             int64
		stopped          string
	}{
		{"stdout past its cap", "while :; do echo aaaaaaaaa; done", "stdout", 1000, "output"},
		{"stderr past its cap", "while :; do echo aaaaaaaaa >&2; done", "stderr", 1000, "output"},
		{"a line of stdout past its cap", "while :; do printf aaaaaaaaaa; done", "stdout", 500, "output"},
		{"stdout at its cap", "for i in $(seq 100); do echo aaaaaaaaa; done", "stdout", 1000, ""},
		{"a line at its cap", "head -c 500 /dev/zero | tr '\\000' a; echo; echo b", "stdout", 503, ""},
	} {
		s := spec(t, fake(t, c.body), "", 10*time.Second)
		r, err := Process(s)
		info, statErr := os.Stat(map[string]string{"stdout": s.Stdout, "stderr": s.Stderr}[c.file])
		if err != nil || statErr != nil || r.StoppedBy != c.stopped || info.Size() != c.size {
			t.Errorf("%s: %+v, %v %v; want stopped by %q, %s of %d bytes", c.name, r, err, statErr, c.stopped, c.file, c.size)
		}
	}
}

func TestProcessTimesTheFirstOutput(t *testing.T) {
	r, err := Process(spec(t, fake(t, "echo e >&2; sleep 0.3; echo x; sleep 0.2"), "", 5*time.Second))
	if err != nil || r.First.Sub(r.Start) < 300*time.Millisecond || r.End.Sub(r.First) < 200*time.Millisecond {
		t.Errorf("%+v, %v; want the first byte of stdout after 0.3 s, and the end 0.2 s later", r, err)
	}
	r, err = Process(spec(t, fake(t, "echo e >&2"), "", 5*time.Second))
	if err != nil || !r.First.IsZero() || r.End.Before(r.Start) {
		t.Errorf("%+v, %v; want no first output", r, err)
	}
}

func TestProcessOfAMissingProgram(t *testing.T) {
	s := spec(t, filepath.Join(t.TempDir(), "none"), "", time.Second)
	if _, err := Process(s); !errors.Is(err, ErrStart) {
		t.Errorf("%v, want ErrStart", err)
	}
}
