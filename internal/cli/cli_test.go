package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func run(args ...string) (int, string, string) {
	var out, errOut bytes.Buffer
	code := Run(args, &out, &errOut)
	return code, out.String(), errOut.String()
}

func TestVersionPrintsTheVersionAndExitsZero(t *testing.T) {
	code, out, errOut := run("version")
	if code != 0 {
		t.Fatalf("exit %d, want 0 (stderr %q)", code, errOut)
	}
	if want := "layup " + Version + "\n"; out != want {
		t.Fatalf("stdout %q, want %q", out, want)
	}
	if Version != "0.1.0-dev" {
		t.Fatalf("Version %q, want the pinned 0.1.0-dev", Version)
	}
}

func TestNoSubcommandPrintsUsageAndExitsTwo(t *testing.T) {
	code, out, errOut := run()
	if code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
	if out != "" || !strings.Contains(errOut, "usage: layup") {
		t.Fatalf("stdout %q, stderr %q: want usage on stderr only", out, errOut)
	}
}

func TestUnknownSubcommandExitsTwo(t *testing.T) {
	code, _, errOut := run("frobnicate")
	if code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
	if !strings.Contains(errOut, `unknown command "frobnicate"`) {
		t.Fatalf("stderr %q: want the unknown command named", errOut)
	}
}

func TestPSBCheckExitCodes(t *testing.T) {
	dir := t.TempDir()
	clean := dir + "/clean.md"
	gaps := dir + "/gaps.md"
	writeFile(t, clean, "**Technology stack:** Go.\n")
	writeFile(t, gaps, "No stack is named here.\n")

	if code, out, _ := run("psb", "check", clean); code != 0 || !strings.HasPrefix(out, "id\trule\t") {
		t.Fatalf("clean: exit %d, stdout %q; want 0 and the header", code, out)
	}
	if code, out, _ := run("psb", "check", gaps); code != 1 || !strings.Contains(out, "\tG1\t") {
		t.Fatalf("gaps: exit %d, stdout %q; want 1 and a G1 row", code, out)
	}
	if code, _, errOut := run("psb", "check", dir+"/missing.md"); code != 2 || errOut == "" {
		t.Fatalf("unreadable: exit %d, stderr %q; want 2 and a message", code, errOut)
	}
	if code, _, _ := run("psb", "check"); code != 2 {
		t.Fatalf("no file: exit %d, want 2", code)
	}
	if code, _, _ := run("psb"); code != 2 {
		t.Fatalf("no psb subcommand: exit %d, want 2", code)
	}
}

// K34: `layup version` takes no argument.
func TestVersionWithAnArgumentIsAUsageError(t *testing.T) {
	if code, out, _ := run("version", "extra"); code != 2 || out != "" {
		t.Fatalf("exit %d, stdout %q; want 2 and nothing", code, out)
	}
}

func TestAUsageErrorPrintsTheReasonAndTheUsageOnStandardErrorOnly(t *testing.T) {
	code, out, errOut := run("psb", "check", "a.md", "b.md")
	if code != 2 || out != "" {
		t.Fatalf("exit %d, stdout %q; want 2 and nothing", code, out)
	}
	if want := "layup: extra argument \"b.md\"\n\nusage: layup "; !strings.HasPrefix(errOut, want) {
		t.Fatalf("stderr %q, want it to start with %q", errOut, want)
	}
}

func TestAnInputErrorPrintsTheReasonAndNoUsage(t *testing.T) {
	code, out, errOut := run("psb", "check", filepath.Join(t.TempDir(), "missing.md"))
	if code != 2 || out != "" || !strings.HasPrefix(errOut, "layup: ") || strings.Contains(errOut, "usage:") {
		t.Fatalf("exit %d, stdout %q, stderr %q; want 2, nothing, and the reason only", code, out, errOut)
	}
}

func TestTheUsageOfLayupListsItsCommands(t *testing.T) {
	_, _, errOut := run()
	for _, want := range []string{"\n  version ", "\n  psb check FILE "} {
		if !strings.Contains(errOut, want) {
			t.Errorf("the usage has no %q:\n%s", want, errOut)
		}
	}
}

// A check that did not run never gives 0 (NFR-004): not-active, a word that
// is not a result, and a table with no row give 1.
func TestExitCode(t *testing.T) {
	for _, c := range []struct {
		results []string
		want    int
	}{
		{[]string{"pass"}, 0},
		{[]string{"pass", "clear"}, 0},
		{[]string{"done", "operator", "done"}, 0},
		{[]string{"pass", "fail"}, 1},
		{[]string{"pass", "not-active"}, 1},
		{[]string{"not-active"}, 1},
		{nil, 1},
		{[]string{"pass", "skipped"}, 1},
		{[]string{"PASS"}, 1},
		{[]string{""}, 1},
	} {
		if got := exitCode(c.results); got != c.want {
			t.Errorf("exitCode(%q) = %d, want %d", c.results, got, c.want)
		}
	}
}
