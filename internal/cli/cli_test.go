package cli

import (
	"bytes"
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
