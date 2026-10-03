//go:build e2e

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/standin"
)

// standInWork makes a stand-in work area (K11), with the git of the test
// isolated from the host's configuration.
func standInWork(t *testing.T) string {
	t.Helper()
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1"} {
		t.Setenv(k, v)
	}
	w, err := standin.Make(t.TempDir(), standin.Options{})
	if err != nil {
		t.Fatal(err)
	}
	return w.Dir
}

// The demos of #84 and #92: the built binary on a stand-in work area prints
// a row of each check, each pass, then a row of each kind of the manifest:
// gate:static passes, as the gate gives clear on the setup head and fail on
// the commit of its fixture, and gate:layout is pending. So a correct setup
// exits 0 (NFR-004), and two runs give the same bytes (NFR-005).
func TestSetupVerifyOnAStandInWorkArea(t *testing.T) {
	w := standInWork(t)
	r := repeatWith(t, goEnv(t), "setup", "verify", w)
	want := "check\tresult\treason\n"
	for _, c := range []string{"discipline-tests", "pin", "kit-history", "facts", "onboarding", "glossary", "guardrails", "markers", "adapted", "identity",
		"link-lint", "sources", "jobs", "gate:static"} {
		want += c + "\tpass\t\u2014\n"
	}
	if want += "gate:layout\tclear\tpending: fixture not run\n"; r.code != 0 || r.stdout != want {
		t.Fatalf("exit %d, stdout:\n%s\nwant 0 and\n%s", r.code, r.stdout, want)
	}
	if !strings.Contains(r.stderr, "layup setup verify: [1/15] discipline-tests\n") || !strings.Contains(r.stderr, "layup setup verify: [15/15] gate:layout\n") ||
		!strings.Contains(r.stderr, "layup setup verify: the fixture run of static on ") {
		t.Errorf("stderr %q; want the progress line of each row, and the line of each run of the gate", r.stderr)
	}
}

// An input error gives exit 2, the reason with no usage, and nothing on
// standard output.
func TestSetupVerifyInputErrors(t *testing.T) {
	w := standInWork(t)
	if err := os.RemoveAll(filepath.Join(w, "target")); err != nil {
		t.Fatal(err)
	}
	if r := layup(t, "setup", "verify", w); r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: target: no commit at the branch layup-setup: ") || strings.Contains(r.stderr, "usage:") {
		t.Errorf("no target: exit %d, stdout %q, stderr %q; want 2, nothing and the reason", r.code, r.stdout, r.stderr)
	}
	for reason, args := range map[string][]string{"missing argument WORK": {"setup", "verify"},
		`extra argument "x"`: {"setup", "verify", w, "x"}, `unknown flag "--x"`: {"setup", "verify", w, "--x", "y"}} {
		if r := layup(t, args...); r.code != 2 || r.stdout != "" || !strings.HasPrefix(r.stderr, "layup: "+reason+"\n\nusage: layup ") {
			t.Errorf("layup %q: exit %d, stdout %q, stderr %q; want 2, nothing and %q with the usage", args, r.code, r.stdout, r.stderr, reason)
		}
	}
}

// A baseline script that is not there gives its row not-active and exit 1
// (#87, NFR-004): a check that does not run cannot pass.
func TestSetupVerifyWithNoBaselineScript(t *testing.T) {
	for k, v := range map[string]string{"HOME": t.TempDir(), "XDG_CONFIG_HOME": t.TempDir(), "GIT_CONFIG_NOSYSTEM": "1"} {
		t.Setenv(k, v)
	}
	w, err := standin.Make(t.TempDir(), standin.Options{Files: map[string]string{"docs/tests/run-discipline-tests.sh": ""}})
	if err != nil {
		t.Fatal(err)
	}
	r := layup(t, "setup", "verify", w.Dir)
	if r.code != 1 || !strings.HasPrefix(r.stdout, "check\tresult\treason\ndiscipline-tests\tnot-active\tmissing: docs/tests/run-discipline-tests.sh\n") {
		t.Errorf("exit %d, stdout:\n%s", r.code, r.stdout)
	}
}
