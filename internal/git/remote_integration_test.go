//go:build integration

package git

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync/atomic"
	"testing"
)

// RemoteHead against real bare repositories (task T-syxm): a branch that is there,
// one that is not (exit 2 of ls-remote), a branch that exists only as the tail of a
// longer ref, and a remote that cannot be reached (an error, never an absence).
func TestRemoteHeadOfABareRepository(t *testing.T) {
	home := isolate(t)
	bare := filepath.Join(t.TempDir(), "target.git")
	gitOK(t, "", plain(home), "init", "--bare", "-q", bare)
	url := "file://" + bare
	work := t.TempDir()
	first := commitTree(t, work, map[string]string{"a.txt": "a\n"})
	must(t, Push(work, url, first, "task/T-ab12/1", Auth{}))
	must(t, Push(work, url, first, "side/refs/heads/task/T-ab12/2", Auth{}))

	sha, found, err := RemoteHead(url, "task/T-ab12/1", Auth{})
	if err != nil || !found || sha != first {
		t.Errorf("a branch that is there: %q, %v, %v; want %s, true, nil", sha, found, err, first)
	}
	if sha, found, err = RemoteHead(url, "task/T-ab12/9", Auth{}); err != nil || found || sha != "" {
		t.Errorf("a branch that is absent: %q, %v, %v; want none, false, nil", sha, found, err)
	}
	// ls-remote matches the tail of a ref name: the branch task/T-ab12/2 exists only as such a tail.
	if sha, found, err = RemoteHead(url, "task/T-ab12/2", Auth{}); err != nil || found {
		t.Errorf("a branch that is only the tail of a longer ref: %q, %v, %v; want absent", sha, found, err)
	}
	_, found, err = RemoteHead("file://"+filepath.Join(t.TempDir(), "none.git"), "main", Auth{})
	var failed *FailedError
	if !errors.As(err, &failed) || failed.Code != 128 || found {
		t.Errorf("a remote that does not exist: %v, found %v; want a *FailedError of code 128", err, found)
	}
}

// With a token the request carries the header of forge.md (The App identity) and
// the token is in no argument; with none, no header.
func TestRemoteHeadSendsTheTokenAsAHeaderOnly(t *testing.T) {
	isolate(t)
	var got atomic.Value
	got.Store("")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header.Get("Authorization"); h != "" {
			got.Store(h)
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()
	_, _, err := RemoteHead(srv.URL+"/acme/target.git", "main", Auth{Web: srv.URL, Token: "ghs_testtoken"})
	if err == nil {
		t.Error("a server with no repository: no error")
	}
	if h := got.Load().(string); h != "Basic eC1hY2Nlc3MtdG9rZW46Z2hzX3Rlc3R0b2tlbg==" {
		t.Errorf("the server got Authorization %q, want Basic of x-access-token and the token", h)
	}
	got.Store("")
	RemoteHead(srv.URL+"/acme/target.git", "main", Auth{})
	if h := got.Load().(string); h != "" {
		t.Errorf("with no token the server got Authorization %q, want none", h)
	}
}
