package git

import (
	"errors"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

// exitCode gives the error of a process that exits with code, as git's does.
func exitCode(t *testing.T, code string) error {
	t.Helper()
	err := exec.Command("sh", "-c", "exit "+code).Run()
	if err == nil {
		t.Fatalf("exit %s: no error", code)
	}
	return err
}

// RemoteHead runs the verb of packages.md (The calls of M2e), with the token as a
// header of the environment when it has one (task T-syxm).
func TestRemoteHeadRunsItsVerb(t *testing.T) {
	url := "https://example.invalid/b.git"
	for _, tt := range []struct {
		name string
		auth Auth
		env  []string
	}{{"ls-remote of a branch", Auth{}, nil}, {"ls-remote of a branch with a token", testAuth, tokenEnv}} {
		t.Run(tt.name, func(t *testing.T) {
			calls := stub(t, fullID+"\trefs/heads/task/T-ab12/1\n", nil)
			sha, found, err := RemoteHead(url, "task/T-ab12/1", tt.auth)
			c := only(t, calls)
			if want := append(append([]string{}, wantConfig...), strings.Fields("ls-remote --exit-code -- "+url+" refs/heads/task/T-ab12/1")...); c.dir != "" || !reflect.DeepEqual(c.args, want) {
				t.Errorf("dir %q, args %q\nwant args %q", c.dir, c.args, want)
			}
			if want := environ(tt.env...); !reflect.DeepEqual(c.env, want) {
				t.Errorf("env %q\nwant the fixed list, then %q", c.env, tt.env)
			}
			if err != nil || !found || sha != fullID {
				t.Errorf("RemoteHead = %q, %v, %v; want %s, true, nil", sha, found, err, fullID)
			}
		})
	}
}

// A branch that is empty, holds ':' or starts with '-' is refused before git starts.
func TestRemoteHeadRefusesABranchBeforeGitStarts(t *testing.T) {
	for _, branch := range []string{"", "a:b", "-x"} {
		calls := stub(t, "", nil)
		_, found, err := RemoteHead("https://example.invalid/b.git", branch, Auth{})
		var failed *FailedError
		if !errors.As(err, &failed) || failed.Code != -1 || found || len(*calls) != 0 {
			t.Errorf("branch %q: found %v, err %v, %d starts; want a *FailedError of code -1 and no start", branch, found, err, len(*calls))
		}
	}
}

// ls-remote matches its pattern against the tail of a ref name: only the line of
// exactly refs/heads/BRANCH is the branch (Q5). Exit 2 is no such branch, another
// code is an error, and an ID that is no full object ID on the exact line is an
// error, never an absence.
func TestRemoteHeadReadsOnlyTheExactRefLine(t *testing.T) {
	const want = "refs/heads/task/T-ab12/1"
	other := strings.Repeat("ab", 20)
	for _, tt := range []struct {
		name, out string
		err       error
		sha       string
		found     bool
		code      int // of the error, -2 for none
	}{
		{"the exact line", fullID + "\t" + want + "\n", nil, fullID, true, -2},
		{"the exact line after a longer ref", other + "\trefs/heads/x/" + want + "\n" + fullID + "\t" + want + "\n", nil, fullID, true, -2},
		{"only a longer ref", other + "\trefs/heads/x/" + want + "\n", nil, "", false, -2},
		{"exit 2", "", exitCode(t, "2"), "", false, -2},
		{"exit 128", "", exitCode(t, "128"), "", false, 128},
		{"a short ID on the exact line", "abc123\t" + want + "\n", nil, "", false, 0},
		{"an upper-case ID on the exact line", strings.ToUpper(fullID) + "\t" + want + "\n", nil, "", false, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub(t, tt.out, tt.err)
			sha, found, err := RemoteHead("https://example.invalid/b.git", "task/T-ab12/1", Auth{})
			var failed *FailedError
			switch {
			case tt.code == -2 && err != nil, tt.code != -2 && (!errors.As(err, &failed) || failed.Code != tt.code):
				t.Errorf("err %v; want code %d", err, tt.code)
			}
			if sha != tt.sha || found != tt.found {
				t.Errorf("RemoteHead = %q, %v; want %q, %v", sha, found, tt.sha, tt.found)
			}
		})
	}
}
