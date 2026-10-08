//go:build integration

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
)

// fakeGitHub plays the calls of M2a (forge.md, The test of the adapter) on a
// loopback server. Each field changes one answer.
type fakeGitHub struct {
	t   *testing.T
	now func() time.Time

	mu          sync.Mutex
	tokens      int               // the installation tokens made
	requests    map[string]int    // the requests, by "METHOD path"
	permissions map[string]string // of the installation
	noComments  bool              // the server lacks the endpoint of Comments
	reject      map[string][]int  // "METHOD path" → the statuses to give first
	headers     map[string]http.Header
	branches    int
	userNull    bool
}

const (
	appID    = 42
	instID   = 7
	owner    = "acme"
	repoName = "target"
)

func newFake(t *testing.T, now func() time.Time) (*fakeGitHub, *httptest.Server) {
	f := &fakeGitHub{t: t, now: now, requests: map[string]int{}, reject: map[string][]int{},
		headers:     map[string]http.Header{},
		permissions: map[string]string{"contents": "write", "issues": "write", "metadata": "read"},
		branches:    1}
	s := httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(s.Close)
	return f, s
}

func (f *fakeGitHub) token(n int) string { return fmt.Sprintf("ghs_test%d", n) }

func (f *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	call := r.Method + " " + r.URL.Path
	f.requests[call]++
	if r.Header.Get("X-GitHub-Api-Version") != "2022-11-28" {
		http.Error(w, `{"message":"no version header"}`, 400)
		return
	}
	if q := f.reject[call]; len(q) > 0 {
		f.reject[call] = q[1:]
		for k, v := range f.headers[call] {
			w.Header()[k] = v
		}
		w.WriteHeader(q[0])
		io.WriteString(w, `{"message":"rejected by the test\nsecond line"}`)
		return
	}
	auth := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	jwtCall := call == "GET /repos/acme/target/installation" || call == "POST /app/installations/7/access_tokens"
	if jwtCall {
		if msg := checkJWT(auth, &key(f.t).PublicKey, appID, f.now()); msg != "" {
			http.Error(w, `{"message":"`+msg+`"}`, 401)
			return
		}
	} else if auth != f.token(f.tokens) || f.tokens == 0 {
		http.Error(w, `{"message":"Bad credentials"}`, 401)
		return
	}
	switch {
	case call == "GET /repos/acme/target/installation":
		json.NewEncoder(w).Encode(map[string]any{"id": instID, "permissions": f.permissions})
	case call == "POST /app/installations/7/access_tokens":
		f.tokens++
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"token": f.token(f.tokens), "expires_at": f.now().Add(time.Hour).Format(time.RFC3339)})
	case call == "GET /repos/acme/target":
		json.NewEncoder(w).Encode(map[string]any{"default_branch": "main", "visibility": "public"})
	case call == "GET /repos/acme/target/branches":
		if r.URL.Query().Get("per_page") != "100" {
			http.Error(w, `{"message":"per_page is not 100"}`, 400)
			return
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		list := []map[string]any{}
		for i := (page - 1) * 100; i < f.branches && i < page*100; i++ {
			name := "main"
			if i > 0 {
				name = fmt.Sprintf("b%03d", i)
			}
			list = append(list, map[string]any{"name": name})
		}
		if page*100 < f.branches {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/repos/acme/target/branches?per_page=100&page=%d>; rel="next"`, r.Host, page+1))
		}
		json.NewEncoder(w).Encode(list)
	case call == "GET /users/acme-layup[bot]":
		json.NewEncoder(w).Encode(map[string]any{"id": 9001, "login": "acme-layup[bot]"})
	case call == "GET /users/nobody":
		http.Error(w, `{"message":"Not Found"}`, 404)
	case call == "POST /repos/acme/target/issues":
		var in struct{ Title, Body string }
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Title == "" {
			http.Error(w, `{"message":"no title"}`, 422)
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"number": 3, "title": in.Title})
	case call == "GET /repos/acme/target/issues/3/comments" && !f.noComments:
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		if r.URL.Query().Get("per_page") != "100" {
			http.Error(w, `{"message":"per_page is not 100"}`, 400)
			return
		}
		if page < 3 {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/repos/acme/target/issues/3/comments?per_page=100&page=%d>; rel="next", <http://%s/repos/acme/target/issues/3/comments?per_page=100&page=3>; rel="last"`, r.Host, page+1, r.Host))
		}
		user := any(map[string]any{"id": 100 + page, "login": fmt.Sprintf("user%d", page)})
		if f.userNull && page == 2 {
			user = nil
		}
		app := any(nil)
		if page == 3 {
			app = map[string]any{"slug": "layup-watch"}
		}
		json.NewEncoder(w).Encode([]map[string]any{{
			"id": 900 + page, "user": user, "performed_via_github_app": app,
			"created_at": "2026-10-08T12:00:00Z", "updated_at": "2026-10-08T12:05:00Z",
			"body": fmt.Sprintf("comment %d", page)}})
	default:
		http.Error(w, `{"message":"Not Found"}`, 404)
	}
}

// clock is a fake clock; its sleep moves it and records each wait.
type clock struct {
	mu     sync.Mutex
	t      time.Time
	sleeps []time.Duration
}

func (c *clock) now() time.Time { c.mu.Lock(); defer c.mu.Unlock(); return c.t }
func (c *clock) sleep(_ context.Context, d time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.sleeps = append(c.sleeps, d)
	c.t = c.t.Add(d)
	return nil
}

// setup gives the fake server, the clock, the adapter and its progress lines.
func setup(t *testing.T) (*fakeGitHub, *clock, *Adapter, *[]string) {
	c := &clock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	f, s := newFake(t, c.now)
	var lines []string
	a := New(Config{API: s.URL, AppID: appID, Key: key(t), Owner: owner, Name: repoName,
		Client: s.Client(), Progress: func(l string) { lines = append(lines, l) }, Now: c.now, Sleep: c.sleep})
	return f, c, a, &lines
}

// noToken fails the test if text holds an installation token.
func noToken(t *testing.T, text string) {
	t.Helper()
	if strings.Contains(text, "ghs_") {
		t.Errorf("a token is in %q", text)
	}
}

func TestTheAdapterPlaysEachCallOfM2a(t *testing.T) {
	f, c, a, _ := setup(t)
	ctx := context.Background()
	if got := forge.Missing(a.Capabilities()); len(got) != 0 {
		t.Errorf("the adapter does not declare %v", got)
	}
	inst, err := a.Installation(ctx)
	if err != nil || inst.ID != instID || !reflect.DeepEqual(inst.Permissions, f.permissions) {
		t.Fatalf("Installation = %+v, %v", inst, err)
	}
	tok, err := a.Token(ctx)
	if err != nil || tok.Token != "ghs_test1" || !tok.Expires.Equal(c.now().Add(time.Hour)) {
		t.Fatalf("Token = %v, %v", tok.Expires, err)
	}
	repo, err := a.Repository(ctx)
	if err != nil || !reflect.DeepEqual(repo, forge.Repository{DefaultBranch: "main", Visibility: "public", HasCommit: true, Branches: []string{"main"}}) {
		t.Errorf("Repository = %+v, %v", repo, err)
	}
	f.branches = 101 // two pages
	if repo, err := a.Repository(ctx); err != nil || len(repo.Branches) != 101 || repo.Branches[0] != "main" || repo.Branches[100] != "b100" {
		t.Errorf("Repository with 101 branches: %d branches, %v; want each page", len(repo.Branches), err)
	}
	f.branches = 0
	if repo, err := a.Repository(ctx); err != nil || repo.HasCommit || len(repo.Branches) != 0 {
		t.Errorf("Repository with no branch = %+v, %v; want no commit", repo, err)
	}
	if id, err := a.UserID(ctx, "acme-layup[bot]"); err != nil || id != 9001 {
		t.Errorf("UserID = %d, %v", id, err)
	}
	if n, err := a.OpenIssue(ctx, "Intake", "the body"); err != nil || n != 3 {
		t.Errorf("OpenIssue = %d, %v", n, err)
	}
	f.userNull = true
	got, err := a.Comments(ctx, 3)
	if err != nil {
		t.Fatal(err)
	}
	created := time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)
	want := []forge.Comment{
		{ID: 901, AuthorID: 101, AuthorLogin: "user1", App: "", Created: created, Updated: created.Add(5 * time.Minute), Body: "comment 1"},
		{ID: 902, AuthorID: 0, AuthorLogin: "ghost", App: "", Created: created, Updated: created.Add(5 * time.Minute), Body: "comment 2"},
		{ID: 903, AuthorID: 103, AuthorLogin: "user3", App: "layup-watch", Created: created, Updated: created.Add(5 * time.Minute), Body: "comment 3"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Comments over three pages, a deleted account on page 2:\n got %+v\nwant %+v", got, want)
	}
	if f.tokens != 1 {
		t.Errorf("%d tokens made, want 1", f.tokens)
	}
}

func TestTokenIsRenewedWhenLessThanFiveMinutesRemain(t *testing.T) {
	f, c, a, _ := setup(t)
	ctx := context.Background()
	if _, err := a.Token(ctx); err != nil {
		t.Fatal(err)
	}
	c.t = c.t.Add(54 * time.Minute) // six minutes left
	if tok, err := a.Token(ctx); err != nil || tok.Token != "ghs_test1" {
		t.Fatalf("Token with six minutes left = %v; want the one it holds", err)
	}
	c.t = c.t.Add(2 * time.Minute) // four minutes left
	if tok, err := a.Token(ctx); err != nil || tok.Token != "ghs_test2" {
		t.Fatalf("Token with four minutes left: %v; want a new one", err)
	}
	if _, err := a.UserID(ctx, "acme-layup[bot]"); err != nil {
		t.Fatal(err)
	}
	if f.tokens != 2 || f.requests["POST /app/installations/7/access_tokens"] != 2 {
		t.Errorf("%d tokens made, want 2", f.tokens)
	}
}

func TestA401MakesOneNewTokenAndOneRetry(t *testing.T) {
	f, _, a, _ := setup(t)
	ctx := context.Background()
	f.reject["GET /users/acme-layup[bot]"] = []int{401}
	if id, err := a.UserID(ctx, "acme-layup[bot]"); err != nil || id != 9001 {
		t.Fatalf("UserID after one 401 = %d, %v; want the ID", id, err)
	}
	if f.tokens != 2 {
		t.Errorf("%d tokens made, want 2 (the first, and one after the 401)", f.tokens)
	}
	f.reject["GET /users/acme-layup[bot]"] = []int{401, 401, 401}
	before := f.requests["GET /users/acme-layup[bot]"]
	_, err := a.UserID(ctx, "acme-layup[bot]")
	if n := f.requests["GET /users/acme-layup[bot]"] - before; n != 2 {
		t.Errorf("%d requests after two 401, want 2 (one retry)", n)
	}
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Status != 401 || fe.Call != "UserID" {
		t.Fatalf("UserID after two 401 = %v; want a forge.Error of UserID, 401", err)
	}
	noToken(t, err.Error())
}

func TestEachForgeErrorNamesTheCallAndTheStatus(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name   string
		break_ func(*fakeGitHub)
		call   func(*Adapter) error
		status int
		inMsg  string
	}{
		{"an unexpected status", func(f *fakeGitHub) { f.reject["POST /repos/acme/target/issues"] = []int{410} },
			func(a *Adapter) error { _, err := a.OpenIssue(ctx, "t", "b"); return err }, 410, "rejected by the test"},
		{"an unknown login", func(*fakeGitHub) {},
			func(a *Adapter) error { _, err := a.UserID(ctx, "nobody"); return err }, 404, "Not Found"},
		{"a server with no endpoint of Comments", func(f *fakeGitHub) { f.noComments = true },
			func(a *Adapter) error { _, err := a.Comments(ctx, 3); return err }, 404, "Not Found"},
		{"a 403 with a reset time but requests left", func(f *fakeGitHub) {
			f.reject["GET /repos/acme/target"] = []int{403}
			f.headers["GET /repos/acme/target"] = http.Header{"X-Ratelimit-Reset": {"1791460800"}, "X-Ratelimit-Remaining": {"5"}}
		}, func(a *Adapter) error { _, err := a.Repository(ctx); return err }, 403, "rejected by the test"},
		{"a 429 with no reset time", func(f *fakeGitHub) { f.reject["GET /repos/acme/target"] = []int{429} },
			func(a *Adapter) error { _, err := a.Repository(ctx); return err }, 429, "rejected by the test"},
	} {
		f, clk, a, _ := setup(t)
		c.break_(f)
		err := c.call(a)
		var fe *forge.Error
		switch {
		case !errors.As(err, &fe):
			t.Errorf("%s: %v; want a forge.Error", c.name, err)
		case fe.Status != c.status || fe.Message != c.inMsg || fe.Call == "":
			t.Errorf("%s: %+v; want the call, status %d and %q", c.name, fe, c.status, c.inMsg)
		case len(clk.sleeps) != 0:
			t.Errorf("%s: the adapter waited %v", c.name, clk.sleeps)
		}
		if err != nil {
			noToken(t, err.Error())
		}
	}
}

func TestABodyThatCannotBeReadOrNoServerIsAnError(t *testing.T) {
	ctx := context.Background()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "not json") }))
	defer s.Close()
	c := &clock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	a := New(Config{API: s.URL, AppID: appID, Key: key(t), Owner: owner, Name: repoName, Client: s.Client(), Progress: func(string) {}, Now: c.now, Sleep: c.sleep})
	_, err := a.Installation(ctx)
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "Installation" || fe.Status != 200 {
		t.Errorf("a body that is not JSON: %v; want a forge.Error of Installation, 200", err)
	}
	s.Close()
	_, err = a.Installation(ctx)
	if !errors.As(err, &fe) || fe.Call != "Installation" || fe.Status != 0 {
		t.Errorf("a closed server: %v; want a forge.Error of Installation with no status", err)
	}
}

func TestARateLimitWaitsUntilTheResetWithAProgressLineEveryTenSeconds(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		name    string
		status  int
		headers http.Header
		wait    time.Duration
	}{
		{"429 with retry-after", 429, http.Header{"Retry-After": {"25"}}, 25 * time.Second},
		{"403 with no request left", 403, http.Header{"X-Ratelimit-Remaining": {"0"},
			"X-Ratelimit-Reset": {strconv.FormatInt(time.Date(2026, 10, 8, 12, 0, 30, 0, time.UTC).Unix(), 10)}}, 30 * time.Second},
	} {
		f, clk, a, lines := setup(t)
		f.reject["GET /repos/acme/target"] = []int{c.status}
		f.headers["GET /repos/acme/target"] = c.headers
		if _, err := a.Repository(ctx); err != nil {
			t.Errorf("%s: %v; want the call after the wait", c.name, err)
			continue
		}
		var total time.Duration
		for _, d := range clk.sleeps {
			if d > 10*time.Second {
				t.Errorf("%s: a wait of %v, more than ten seconds with no progress line", c.name, d)
			}
			total += d
		}
		if total < c.wait || len(*lines) != len(clk.sleeps) {
			t.Errorf("%s: waited %v in %d sleeps with %d progress lines; want %v, a line per sleep", c.name, total, len(clk.sleeps), len(*lines), c.wait)
		}
		for _, l := range *lines {
			noToken(t, l)
		}
		if f.requests["GET /repos/acme/target"] != 2 {
			t.Errorf("%s: %d requests, want 2", c.name, f.requests["GET /repos/acme/target"])
		}
	}
}

func TestAMissingPermissionOfTheInstallationIsNamed(t *testing.T) {
	f, _, a, _ := setup(t)
	delete(f.permissions, "issues")
	inst, err := a.Installation(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got := forge.CheckPermissions(inst.Permissions); !reflect.DeepEqual(got, []string{"issues: write"}) {
		t.Errorf("CheckPermissions = %v, want [issues: write]", got)
	}
}

func TestALinkToAnotherHostIsRefused(t *testing.T) {
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("the adapter sent a request to another host: %s", r.Header.Get("Authorization"))
	}))
	defer other.Close()
	c := &clock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	f, s := newFake(t, c.now)
	_ = f
	wrap := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/comments") {
			w.Header().Set("Link", "<"+other.URL+"/next>; rel=\"next\"")
			io.WriteString(w, "[]")
			return
		}
		r2, _ := http.NewRequest(r.Method, s.URL+r.URL.RequestURI(), r.Body)
		r2.Header = r.Header
		resp, err := s.Client().Do(r2)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}))
	defer wrap.Close()
	a := New(Config{API: wrap.URL, AppID: appID, Key: key(t), Owner: owner, Name: repoName, Client: wrap.Client(), Progress: func(string) {}, Now: c.now, Sleep: c.sleep})
	_, err := a.Comments(context.Background(), 3)
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "Comments" {
		t.Errorf("a next page on another host: %v; want a forge.Error of Comments", err)
	}
}
