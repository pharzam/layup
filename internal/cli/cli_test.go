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
