//go:build integration

package session

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The input of a session of a phase (docs/spec/intake.md, The sessions of
// Intake): input/ in the session's directory, each file byte for byte, mode
// 0600; a name that is not one plain file name is refused.
func TestMakeInput(t *testing.T) {
	d := dirOf(t.TempDir(), "S-1a2b3c4d")
	os.MkdirAll(d.Root, 0o700)
	if d.Input != filepath.Join(d.Root, "input") {
		t.Fatalf("Dir.Input %q", d.Input)
	}
	files := map[string][]byte{"problem-statement.md": []byte("# A\r\n\tB\n"), "spans.tsv": []byte("start\tend\n")}
	if err := MakeInput(d, files); err != nil {
		t.Fatal(err)
	}
	for name, want := range files {
		path := filepath.Join(d.Input, name)
		got, err := os.ReadFile(path)
		if err != nil || string(got) != string(want) {
			t.Errorf("input/%s: %q, %v; want %q", name, got, err, want)
		}
		if info, err := os.Stat(path); err != nil || info.Mode().Perm() != 0o600 {
			t.Errorf("input/%s: the mode %v, %v; want 0600", name, info.Mode().Perm(), err)
		}
	}
	for _, name := range []string{"", ".", "..", "a/b", "../x"} {
		e := dirOf(t.TempDir(), "S-1a2b3c4d")
		os.MkdirAll(e.Root, 0o700)
		if err := MakeInput(e, map[string][]byte{name: nil}); err == nil || !strings.Contains(err.Error(), "one file name") {
			t.Errorf("an input named %q: %v; want it refused", name, err)
		}
	}
}

// The output file of a session of a phase is read by the rules of a result
// file, with no block: missing, a link, not a regular file and over 1 MiB are
// refused; a file of exactly 1 MiB is read, byte for byte.
func TestReadOutput(t *testing.T) {
	d := dirOf(t.TempDir(), "S-1a2b3c4d")
	os.MkdirAll(d.Result, 0o700)
	out := filepath.Join(d.Result, "meaning.tsv")
	host := filepath.Join(t.TempDir(), "host.tsv")
	os.WriteFile(host, []byte("a\n"), 0o600)
	want := map[string]string{"missing": "no such file", "a link": "too many levels of symbolic links",
		"a directory": "is not a regular file", "over 1 MiB": "is over 1 MiB"}
	for name, set := range map[string]func(){
		"missing":     func() {},
		"a link":      func() { os.Symlink(host, out) },
		"a directory": func() { os.Mkdir(out, 0o700) },
		"over 1 MiB":  func() { os.WriteFile(out, []byte(strings.Repeat("x", 1<<20+1)), 0o600) },
	} {
		os.RemoveAll(out)
		set()
		if _, err := ReadOutput(d, "meaning.tsv"); err == nil || !strings.Contains(err.Error(), want[name]) {
			t.Errorf("an output file %s: %v, want an error with %q", name, err, want[name])
		}
	}
	os.RemoveAll(out)
	pad := strings.Repeat("x", 1<<20)
	os.WriteFile(out, []byte(pad), 0o600)
	if got, err := ReadOutput(d, "meaning.tsv"); err != nil || string(got) != pad {
		t.Errorf("an output file of 1 MiB: %d bytes, %v", len(got), err)
	}
	if _, err := ReadOutput(d, "../meaning.tsv"); err == nil || !strings.Contains(err.Error(), "one file name") {
		t.Errorf("an output named ../meaning.tsv: %v; want it refused", err)
	}
}
