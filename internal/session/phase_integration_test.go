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

// The output files of a session of a phase (docs/spec/scaffold.md, The prose
// rows; task T-eahr): result/<dir>/ holds exactly the files that the step named,
// at their paths, and the whole set is read, and validated, before any file is
// given. A missing file, a file that no name covers (also in a directory, and an
// empty directory), a link as a file or as a directory above one, an entry that
// is not a regular file, a file over 1 MiB or not valid UTF-8, and a name that is
// not a clean relative path refuse the output.
func TestReadOutputFiles(t *testing.T) {
	names := []string{"README.md", "docs/glossary.md"}
	host := t.TempDir()
	os.WriteFile(filepath.Join(host, "glossary.md"), []byte("host\n"), 0o600)
	for _, c := range []struct {
		name  string
		names []string
		make  func(dir string) // writes the output under dir, after the good one
		want  string           // a part of the error; "" for none
	}{
		{"the files that were named", names, func(string) {}, ""},
		{"an empty file", names, func(dir string) { os.WriteFile(filepath.Join(dir, "README.md"), nil, 0o600) }, ""},
		{"a file of 1 MiB", names, func(dir string) {
			os.WriteFile(filepath.Join(dir, "README.md"), []byte(strings.Repeat("x", 1<<20)), 0o600)
		}, ""},
		{"a missing file", names, func(dir string) { os.Remove(filepath.Join(dir, "README.md")) }, "README.md is missing"},
		{"a file that no name covers", names, func(dir string) { os.WriteFile(filepath.Join(dir, "extra.md"), []byte("x"), 0o600) }, "extra.md, which the stop does not name"},
		{"a nested file that no name covers", names, func(dir string) { os.WriteFile(filepath.Join(dir, "docs", "extra.md"), []byte("x"), 0o600) }, "docs/extra.md, which the stop does not name"},
		{"an empty directory that no name covers", names, func(dir string) { os.Mkdir(filepath.Join(dir, "empty"), 0o700) }, "empty, which the stop does not name"},
		{"a link as a file", names, func(dir string) {
			os.Remove(filepath.Join(dir, "docs", "glossary.md"))
			os.Symlink(filepath.Join(host, "glossary.md"), filepath.Join(dir, "docs", "glossary.md"))
		}, "docs/glossary.md is a link"},
		{"a link as a directory above a file", names, func(dir string) {
			os.RemoveAll(filepath.Join(dir, "docs"))
			os.Symlink(host, filepath.Join(dir, "docs"))
		}, "docs is a link"},
		{"a directory where a file is named", names, func(dir string) {
			os.Remove(filepath.Join(dir, "README.md"))
			os.Mkdir(filepath.Join(dir, "README.md"), 0o700)
		}, "README.md is not a regular file"},
		{"a file over 1 MiB", names, func(dir string) {
			os.WriteFile(filepath.Join(dir, "README.md"), []byte(strings.Repeat("x", 1<<20+1)), 0o600)
		}, "README.md is over 1 MiB"},
		{"a file that is not valid UTF-8", names, func(dir string) { os.WriteFile(filepath.Join(dir, "docs", "glossary.md"), []byte("a\xffb"), 0o600) }, "docs/glossary.md is not valid UTF-8"},
		{"a name that leaves the directory", []string{"../README.md"}, func(string) {}, "not a clean relative path"},
		{"an absolute name", []string{"/etc/passwd"}, func(string) {}, "not a clean relative path"},
		{"a name with an empty part", []string{"docs//glossary.md"}, func(string) {}, "not a clean relative path"},
		{"a name that is the directory itself", []string{"."}, func(string) {}, "not a clean relative path"},
		{"a file two directories deep", []string{"README.md", "a/b/c.md"}, func(dir string) {
			os.RemoveAll(filepath.Join(dir, "docs"))
			os.MkdirAll(filepath.Join(dir, "a", "b"), 0o700)
			os.WriteFile(filepath.Join(dir, "a", "b", "c.md"), []byte("deep\n"), 0o600)
		}, ""},
		{"a file where a directory is named", names, func(dir string) {
			os.RemoveAll(filepath.Join(dir, "docs"))
			os.WriteFile(filepath.Join(dir, "docs"), []byte("x"), 0o600)
		}, "a file at docs, where a directory is needed"},
		{"an empty name", []string{""}, func(string) {}, "not a clean relative path"},
	} {
		d := dirOf(t.TempDir(), "S-1a2b3c4d")
		out := filepath.Join(d.Result, "files")
		os.MkdirAll(filepath.Join(out, "docs"), 0o700)
		os.WriteFile(filepath.Join(out, "README.md"), []byte("# Readme\n"), 0o600)
		os.WriteFile(filepath.Join(out, "docs", "glossary.md"), []byte("# Glossary\n"), 0o600)
		c.make(out)
		got, err := ReadOutputFiles(d, "files", c.names)
		switch {
		case c.want == "" && err != nil:
			t.Errorf("%s: %v; want the files", c.name, err)
		case c.want == "" && len(got) != len(c.names):
			t.Errorf("%s: %d files; want %d", c.name, len(got), len(c.names))
		case c.want != "" && (err == nil || !strings.Contains(err.Error(), c.want)):
			t.Errorf("%s: %v; want an error with %q", c.name, err, c.want)
		case c.want != "" && got != nil:
			t.Errorf("%s: the files %q were given with the error; want none, as the whole set is validated first", c.name, got)
		}
	}
	// The bytes are those of the files.
	d := dirOf(t.TempDir(), "S-1a2b3c4d")
	os.MkdirAll(filepath.Join(d.Result, "files", "docs"), 0o700)
	os.WriteFile(filepath.Join(d.Result, "files", "docs", "glossary.md"), []byte("# G\r\n\tx\n"), 0o600)
	if got, err := ReadOutputFiles(d, "files", []string{"docs/glossary.md"}); err != nil || string(got["docs/glossary.md"]) != "# G\r\n\tx\n" {
		t.Errorf("the bytes of docs/glossary.md: %q, %v", got, err)
	}
	// The directory itself is missing, or is a link, or is named by a path.
	for name, set := range map[string]func(d Dir){
		"no output directory": func(Dir) {},
		"a link as the output directory": func(d Dir) {
			os.Symlink(host, filepath.Join(d.Result, "files"))
		},
		"a link as the result directory, to a directory that holds the output": func(d Dir) {
			os.MkdirAll(filepath.Join(host, "files"), 0o700)
			os.WriteFile(filepath.Join(host, "files", "README.md"), []byte("outside\n"), 0o600)
			os.RemoveAll(d.Result)
			os.Symlink(host, d.Result)
		},
	} {
		e := dirOf(t.TempDir(), "S-1a2b3c4d")
		os.MkdirAll(e.Result, 0o700)
		set(e)
		if _, err := ReadOutputFiles(e, "files", []string{"README.md"}); err == nil || (name != "no output directory" && !strings.Contains(err.Error(), "link")) {
			t.Errorf("%s: %v; want a refusal that names the link", name, err)
		}
	}
	if _, err := ReadOutputFiles(d, "../files", []string{"docs/glossary.md"}); err == nil || !strings.Contains(err.Error(), "one file name") {
		t.Errorf("an output directory named ../files: %v; want it refused", err)
	}
}
