//go:build integration

package session

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// The rule names of these tests are their own, so no file of the host's own
// directories above the test's can match them.
const ruleName = "LAYUP-TEST-RULES.md"

func TestCheckRuleFiles(t *testing.T) {
	top := t.TempDir()
	d := dirOf(filepath.Join(top, "a", "host"), "S-1a2b3c4d")
	os.MkdirAll(d.Repo, 0o700)
	os.WriteFile(filepath.Join(d.Repo, ruleName), []byte("# the target's own\n"), 0o644)
	policy := filepath.Join(top, "policy.json")
	os.WriteFile(policy, []byte("{}\n"), 0o644)
	got, err := CheckRuleFiles(d, []string{ruleName}, []string{policy, filepath.Join(top, "none.json")})
	if err != nil || !slices.Equal(got, []string{policy}) {
		t.Fatalf("no rule file above, a file inside repo/: %q, %v; want the one policy path that exists", got, err)
	}
	for _, where := range []string{d.Root, filepath.Join(top, "a", "host", "sessions"), filepath.Join(top, "a")} {
		p := filepath.Join(where, ruleName)
		os.WriteFile(p, []byte("# a rule file above\n"), 0o644)
		_, err := CheckRuleFiles(d, []string{ruleName}, nil)
		refused(t, "a rule file at "+where, err, "rules", p)
		os.Remove(p)
	}
}

// The two forms that the real git writes: a loose ref, and the same ref after
// git pack-refs --all.
func TestHeadOfARealClone(t *testing.T) {
	src, base := runClone(t)
	d, err := Make(t.TempDir(), "acme/target", "S-1a2b3c4d", src, "main", "T-ab12", 1, base, nil, Harness{})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := HeadOf(d.Repo, "T-ab12", 1); err != nil || got != base {
		t.Errorf("the loose ref: %q, %v; want %s", got, err, base)
	}
	run(t, d.Repo, "pack-refs", "--all")
	if _, err := os.Stat(filepath.Join(d.Repo, ".git", "refs", "heads", "task", "T-ab12", "1")); err == nil {
		t.Fatal("pack-refs left the loose ref")
	}
	if got, err := HeadOf(d.Repo, "T-ab12", 1); err != nil || got != base {
		t.Errorf("the packed ref: %q, %v; want %s", got, err, base)
	}
}
