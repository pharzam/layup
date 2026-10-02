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

// The demo of #84: the built binary on a stand-in work area prints the rows of
// the checks kit-history, pin, facts, onboarding, glossary and guardrails
// (#89), adapted (#88) and identity, each pass, and a not-active row for each
// check that it does not have yet, so it exits 1 (NFR-004); two runs give the
// same bytes (NFR-005).
func TestSetupVerifyOnAStandInWorkArea(t *testing.T) {
	w := standInWork(t)
	r := repeat(t, "setup", "verify", w)
	if r.code != 1 || !strings.HasPrefix(r.stdout, "check\tresult\treason\ndiscipline-tests\tnot-active\tnot built yet\npin\tpass\t—\nkit-history\tpass\t—\n") ||
		!strings.Contains(r.stdout, "\nkit-history\tpass\t—\nfacts\tpass\t—\nonboarding\tpass\t—\nglossary\tpass\t—\nguardrails\tpass\t—\n") ||
		!strings.Contains(r.stdout, "\nadapted\tpass\t—\nidentity\tpass\t—\n") || !strings.HasSuffix(r.stdout, "\ngate:static\tnot-active\tnot built yet\ngate:layout\tnot-active\tnot built yet\n") ||
		strings.Count(r.stdout, "\n") != 16 {
		t.Fatalf("exit %d, stdout:\n%s", r.code, r.stdout)
	}
	if !strings.Contains(r.stderr, "layup setup verify: [1/15] discipline-tests\n") || !strings.Contains(r.stderr, "layup setup verify: [15/15] gate:layout\n") {
		t.Errorf("stderr %q; want the progress line of each row", r.stderr)
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
