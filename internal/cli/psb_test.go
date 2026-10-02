package cli

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// standInFile puts a stand-in for the read of FILE in place until the test
// ends: it gives text, or err.
func standInFile(t *testing.T, text string, err error) {
	t.Helper()
	saved := readFile
	readFile = func(string) ([]byte, error) { return []byte(text), err }
	t.Cleanup(func() { readFile = saved })
}

// A FILE that is not valid UTF-8 gives exit 2 before any rule runs: nothing
// on standard output, and a diagnostic that names the file and the line of the
// first byte that is not valid (K32; D2 of #83).
func TestPSBCheckRefusesAFileThatIsNotUTF8(t *testing.T) {
	for _, c := range []struct {
		text string
		line int
	}{
		{"\xff", 1},
		{"Technology stack: Go\n\nThe café is f\xe9st.\n", 3},
		{"Technology stack: Go\nthe file ends inside a character \xe2\x80", 2},
		{"a surrogate half \xed\xa0\x80\n", 1},
		{"one\ntwo\n\x80\n\xff\n", 3},
	} {
		standInFile(t, c.text, nil)
		code, out, errOut := run("psb", "check", "in.md")
		if want := fmt.Sprintf("layup: in.md: line %d is not valid UTF-8\n", c.line); code != 2 || out != "" || errOut != want {
			t.Errorf("%q: exit %d, stdout %q, stderr %q; want 2, nothing and %q", c.text, code, out, errOut, want)
		}
	}
}

// refusingWriter is an output that refuses each write.
type refusingWriter struct{}

func (refusingWriter) Write([]byte) (int, error) {
	return 0, errors.New("the output refuses the write")
}

// A table that the command cannot write gives exit 2 and the error on standard
// error, with a gap and with no gap (condition 2 of the plan review of #83).
func TestPSBCheckGivesTwoWhenItCannotWriteTheTable(t *testing.T) {
	for _, text := range []string{"No stack is named here.\n", "Technology stack: Go\n"} {
		standInFile(t, text, nil)
		var errOut strings.Builder
		if code := Run([]string{"psb", "check", "in.md"}, refusingWriter{}, &errOut); code != 2 || errOut.String() != "layup: the output refuses the write\n" {
			t.Errorf("%q: exit %d, stderr %q; want 2 and the error of the output", text, code, errOut.String())
		}
	}
}
