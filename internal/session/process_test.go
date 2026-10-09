package session

import (
	"errors"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestWords(t *testing.T) {
	for _, c := range []struct {
		command, model, cap, prompt string
		want                        []string
	}{
		{"claude -p --model {model} --max-budget-usd {cap} {prompt}", "opus", "5", "a b {model}",
			[]string{"claude", "-p", "--model", "opus", "--max-budget-usd", "5", "a b {model}"}},
		{"h --file={prompt} --m={model}", "m1", "", "/s/prompt.md", []string{"h", "--file=/s/prompt.md", "--m=m1"}},
		{"codex exec", "m", "", "/s/prompt.md", []string{"codex", "exec"}},
	} {
		if got := Words(c.command, c.model, c.cap, c.prompt); !slices.Equal(got, c.want) {
			t.Errorf("Words(%q): %q, want %q", c.command, got, c.want)
		}
	}
}

func TestPromptOfEachMode(t *testing.T) {
	d := dirOf("/h", "S-1a2b3c4d")
	for mode, want := range map[string]string{"file": d.Prompt, "stdin": d.Prompt, "arg": "the text"} {
		if got := Prompt(mode, d, []byte("the text")); got != want {
			t.Errorf("Prompt(%s): %q, want %q", mode, got, want)
		}
	}
}

func TestNewSpec(t *testing.T) {
	d := dirOf("/h", "S-1a2b3c4d")
	for mode, stdin := range map[string]string{"file": "", "arg": "", "stdin": d.Prompt} {
		s := NewSpec(d, []string{"h"}, []string{"PATH=/bin"}, mode, 30*time.Minute)
		want := Spec{Words: []string{"h"}, Env: []string{"PATH=/bin"}, Dir: d.Repo, Stdin: stdin,
			Stdout: filepath.Join(d.Root, "stdout"), Stderr: filepath.Join(d.Root, "stderr"),
			Wall: 30 * time.Minute, IntWait: 10 * time.Second, TermWait: 10 * time.Second,
			StdoutCap: 64 << 20, StderrCap: 64 << 20, LineCap: 8 << 20}
		if !reflect.DeepEqual(s, want) {
			t.Errorf("%s: %+v, want %+v", mode, s, want)
		}
	}
}

// steps gives the hooks of a session that record their order, and refuse at
// the step named refuse.
func steps(order *[]string, refuse string, task bool) Steps {
	step := func(name string) func(string) error {
		return func(id string) error {
			*order = append(*order, name)
			if name == refuse {
				return Refusal{Reason: name}
			}
			return nil
		}
	}
	s := Steps{
		ID: func() (string, error) {
			*order = append(*order, "id")
			if refuse == "id" {
				return "", errors.New("no random bytes")
			}
			return "S-1a2b3c4d", nil
		},
		Version: step("version"), Context: step("context"), Prompt: step("prompt"), Rules: step("rules"),
		StartRow: step("start"),
		Process: func(string) (Run, error) {
			*order = append(*order, "process")
			return Run{StoppedBy: "wall"}, nil
		},
		End: func(_ string, r Run, err error) error {
			*order = append(*order, "end "+r.StoppedBy)
			if refuse == "end" {
				return errors.New("a records push refused")
			}
			return nil
		},
	}
	if task {
		s.Attempt = step("attempt")
		s.Push = func(_ string, r Run) error {
			*order = append(*order, "push")
			return nil
		}
	}
	return s
}

func TestCallRunsTheStepsInOrder(t *testing.T) {
	all := []string{"id", "attempt", "version", "context", "prompt", "rules", "start", "process", "end wall", "push"}
	var order []string
	if id, err := Call(steps(&order, "", true)); err != nil || id != "S-1a2b3c4d" || !slices.Equal(order, all) {
		t.Errorf("a task session: %q, %v, the steps %q; want %q", id, err, order, all)
	}
	order = nil
	probe := []string{"id", "version", "context", "prompt", "rules", "start", "process", "end wall"}
	if _, err := Call(steps(&order, "", false)); err != nil || !slices.Equal(order, probe) {
		t.Errorf("a probe: %v, the steps %q; want %q", err, order, probe)
	}
	// Each step that can refuse ends the call before the next: no process
	// starts after a refusal, and nothing runs after an end that fails.
	for i, name := range []string{"id", "attempt", "version", "context", "prompt", "rules", "start", "end"} {
		order = nil
		at := slices.Index(all, name)
		if name == "end" {
			at = slices.Index(all, "end wall")
		}
		id, err := Call(steps(&order, name, true))
		if err == nil || !slices.Equal(order, all[:at+1]) {
			t.Errorf("%d, a refusal at %s: %v, the steps %q; want %q", i, name, err, order, all[:at+1])
		}
		if name != "id" && id != "S-1a2b3c4d" {
			t.Errorf("a refusal at %s: the ID %q, want the session's", name, id)
		}
		var r Refusal
		if name != "id" && name != "end" && (!errors.As(err, &r) || r.Reason != name) {
			t.Errorf("a refusal at %s: %v, want its Refusal", name, err)
		}
	}
}

func TestCallEndsAProcessThatCannotStart(t *testing.T) {
	var order []string
	s := steps(&order, "", true)
	s.Process = func(string) (Run, error) { return Run{}, ErrStart }
	var got error
	s.End = func(_ string, _ Run, err error) error {
		got = err
		return nil
	}
	if _, err := Call(s); err != nil || !errors.Is(got, ErrStart) {
		t.Errorf("%v; End got %v, want ErrStart", err, got)
	}
}
