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
	posted      []string // the body of each comment posted

	rules           map[string][][]map[string]any // the pages of the rules of a branch, by its name as the server decodes it
	rulesBody       string                        // when set, the body of each answer of the rules
	rulesNext       string                        // when set, the next page that each answer of the rules names
	rulePaths       []string                      // the escaped path of each request of the rules
	activity        []map[string]any              // the items of the activity
	activityBody    string                        // when set, the body of the answer of the activity
	activityLink    bool                          // the answer of the activity names a next page
	activityQueries []string                      // the query of each request of the activity

	handlers map[string]func(http.ResponseWriter, *http.Request) // the calls of M2e, by "METHOD path"; the lock is held
}

// head gives the head that the fake gives to the branch of index i.
func (f *fakeGitHub) head(i int) string { return fmt.Sprintf("%040x", i+1) }

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
			list = append(list, map[string]any{"name": name, "commit": map[string]any{"sha": f.head(i)}})
		}
		if page*100 < f.branches {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/repos/acme/target/branches?per_page=100&page=%d>; rel="next"`, r.Host, page+1))
		}
		json.NewEncoder(w).Encode(list)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/repos/acme/target/rules/branches/"):
		branch := strings.TrimPrefix(r.URL.Path, "/repos/acme/target/rules/branches/")
		f.rulePaths = append(f.rulePaths, r.URL.EscapedPath())
		if r.URL.Query().Get("per_page") != "100" {
			http.Error(w, `{"message":"per_page is not 100"}`, 400)
			return
		}
		if f.rulesNext != "" {
			w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, f.rulesNext))
		}
		if f.rulesBody != "" {
			io.WriteString(w, f.rulesBody)
			return
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		pages := f.rules[branch]
		if page < len(pages) {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s%s?per_page=100&page=%d>; rel="next"`, r.Host, r.URL.EscapedPath(), page+1))
		}
		list := []map[string]any{}
		if page <= len(pages) {
			list = pages[page-1]
		}
		json.NewEncoder(w).Encode(list)
	case call == "GET /repos/acme/target/activity":
		f.activityQueries = append(f.activityQueries, r.URL.RawQuery)
		if f.activityLink {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s/repos/acme/target/activity?ref=main&per_page=30&page=2>; rel="next"`, r.Host))
		}
		if f.activityBody != "" {
			io.WriteString(w, f.activityBody)
			return
		}
		json.NewEncoder(w).Encode(f.activity)
	case call == "GET /users/acme-layup[bot]":
		json.NewEncoder(w).Encode(map[string]any{"id": 9001, "login": "acme-layup[bot]"})
	case call == "GET /users/nobody":
		http.Error(w, `{"message":"Not Found"}`, 404)
	case call == "POST /repos/acme/target/issues/3/comments":
		var in struct{ Body string }
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Body == "" {
			http.Error(w, `{"message":"no body"}`, 422)
			return
		}
		f.posted = append(f.posted, in.Body)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]any{"id": 7001, "body": in.Body})
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
	case f.handlers[call] != nil:
		f.handlers[call](w, r)
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
	if err != nil || !reflect.DeepEqual(repo, forge.Repository{DefaultBranch: "main", Visibility: "public", HasCommit: true, Branches: []string{"main"}, Heads: map[string]string{"main": f.head(0)}}) {
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
	body := "S-1a2b3c4d: developer session of T-ab12, attempt 1: done\n"
	if id, err := a.Comment(ctx, 3, body); err != nil || id != 7001 || len(f.posted) != 1 || f.posted[0] != body {
		t.Errorf("Comment = %d, %v, the bodies posted %q; want 7001 and the body", id, err, f.posted)
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
		{"a comment that the forge refuses", func(f *fakeGitHub) { f.reject["POST /repos/acme/target/issues/3/comments"] = []int{422} },
			func(a *Adapter) error { _, err := a.Comment(ctx, 3, "b"); return err }, 422, "rejected by the test"},
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

// The adapter's own client (forge.md, The App identity; task T-mqty): with no
// Client of its caller, it follows no redirect, so a token goes to the api
// only, and its transport reads no proxy of the environment.
func TestTheOwnClientFollowsNoRedirectAndNoProxy(t *testing.T) {
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("the adapter followed a redirect to another server, with %q", r.Header.Get("Authorization"))
	}))
	defer other.Close()
	moved := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+r.URL.Path, http.StatusFound)
	}))
	defer moved.Close()
	c := &clock{t: time.Date(2026, 10, 8, 12, 0, 0, 0, time.UTC)}
	a := New(Config{API: moved.URL, AppID: appID, Key: key(t), Owner: owner, Name: repoName, Progress: func(string) {}, Now: c.now, Sleep: c.sleep})
	_, err := a.Installation(context.Background())
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Status != http.StatusFound || fe.Call != "Installation" {
		t.Errorf("a redirect: %v; want a forge.Error of Installation, 302", err)
	}
	tr, ok := a.cfg.Client.Transport.(*http.Transport)
	if !ok || tr.Proxy != nil {
		t.Errorf("the transport of the own client is %T with a proxy function; want an *http.Transport with none", a.cfg.Client.Transport)
	}
}

// The calls of M2d (forge.md, The calls of M2d; task T-jpkc, #271), played on
// the same server: the rules of a branch, the activity of a ref, and the head of
// each branch. Each case breaks one rule only.

// ruleItem is an item of the answer of GET .../rules/branches/{branch}: the rule,
// with the fields of its ruleset that GitHub adds.
func ruleItem(typ string, parameters map[string]any) map[string]any {
	item := map[string]any{"type": typ, "ruleset_source_type": "Repository", "ruleset_source": "acme/target", "ruleset_id": 11}
	if parameters != nil {
		item["parameters"] = parameters
	}
	return item
}

// EffectiveRules gives each rule of the branch over every page, in order, with its
// type and its parameters as GitHub gives them, and the source type, the source and the
// ID of its ruleset (task T-gkvi, #323); a rule with no parameters, or with an empty
// object of them, has none.
func TestEffectiveRulesGivesEachActiveRuleOfABranchOverEveryPage(t *testing.T) {
	f, _, a, _ := setup(t)
	f.rules = map[string][][]map[string]any{"main": {
		{ruleItem("creation", nil), ruleItem("pull_request", map[string]any{"required_approving_review_count": 0, "allowed_merge_methods": []string{"merge"}, "required_review_thread_resolution": true})},
		{ruleItem("deletion", map[string]any{}), ruleItem("required_status_checks", map[string]any{"strict_required_status_checks_policy": true,
			"required_status_checks": []map[string]any{{"context": "static", "integration_id": 15368}}})},
	}}
	got, err := a.EffectiveRules(context.Background(), "main")
	want := []forge.Rule{
		{Type: "creation", RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11},
		{Type: "pull_request", Parameters: map[string]any{"required_approving_review_count": float64(0), "allowed_merge_methods": []any{"merge"}, "required_review_thread_resolution": true},
			RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11},
		{Type: "deletion", RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11},
		{Type: "required_status_checks", Parameters: map[string]any{"strict_required_status_checks_policy": true,
			"required_status_checks": []any{map[string]any{"context": "static", "integration_id": float64(15368)}}},
			RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("EffectiveRules over two pages = %+v, %v\nwant %+v", got, err, want)
	}
	if n := f.requests["GET /repos/acme/target/rules/branches/main"]; n != 2 {
		t.Errorf("%d requests, want 2 (each page once)", n)
	}
}

// The rules of rulesets of two sources: the repository's and an organization's, each with its own ID; an answer that
// lacks the three fields, or gives them as null, gives their zero values, and a ruleset_id that is anything else but an integer is an
// error of the call.
func TestEffectiveRulesGivesTheRulesetOfEachRuleAndZeroValuesWhenTheAnswerLacksIt(t *testing.T) {
	ctx := context.Background()
	f, _, a, _ := setup(t)
	org := ruleItem("creation", nil)
	org["ruleset_source_type"], org["ruleset_source"], org["ruleset_id"] = "Organization", "acme", 4242
	bare := map[string]any{"type": "deletion"}
	f.rules = map[string][][]map[string]any{"layup-probe": {{ruleItem("creation", nil), org, bare}}}
	got, err := a.EffectiveRules(ctx, "layup-probe")
	want := []forge.Rule{
		{Type: "creation", RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11},
		{Type: "creation", RulesetSourceType: "Organization", RulesetSource: "acme", RulesetID: 4242},
		{Type: "deletion"},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("EffectiveRules with rules of two sources and one with no ruleset = %+v, %v\nwant %+v", got, err, want)
	}
	// A null reads as the zero value, as an absent field does.
	f.rules = map[string][][]map[string]any{"layup-probe": {{{"type": "creation", "ruleset_id": nil, "ruleset_source": nil, "ruleset_source_type": nil}}}}
	if got, err := a.EffectiveRules(ctx, "layup-probe"); err != nil || !reflect.DeepEqual(got, []forge.Rule{{Type: "creation"}}) {
		t.Errorf("EffectiveRules with a null ruleset: %+v, %v; want the rule with the zero values", got, err)
	}
	for _, id := range []any{"11", 1.5, true, []int{11}} {
		f.rules = map[string][][]map[string]any{"layup-probe": {{{"type": "creation", "ruleset_id": id}}}}
		var fe *forge.Error
		if _, err := a.EffectiveRules(ctx, "layup-probe"); !errors.As(err, &fe) || fe.Call != "EffectiveRules" || fe.Status != 200 {
			t.Errorf("EffectiveRules with a ruleset_id of %v: %v; want a forge.Error of EffectiveRules, 200", id, err)
		}
	}
}

// A branch that has no rule gives no rule and no error: the server answers 200 with an empty
// list.
func TestEffectiveRulesOfABranchThatHasNone(t *testing.T) {
	_, _, a, _ := setup(t)
	got, err := a.EffectiveRules(context.Background(), "no-rule-here")
	if err != nil || len(got) != 0 {
		t.Errorf("EffectiveRules of a branch with no rule = %+v, %v; want none", got, err)
	}
}

// A branch that does not exist yet gives the rules of the rulesets that name it, as layup-probe
// does before its probe: the call asks no list of branches, and keeps each rule that the server gives.
func TestEffectiveRulesOfABranchThatDoesNotExistYetGivesTheRulesThatNameIt(t *testing.T) {
	f, _, a, _ := setup(t)
	f.rules = map[string][][]map[string]any{"layup-probe": {{ruleItem("creation", nil), ruleItem("deletion", nil)}}}
	got, err := a.EffectiveRules(context.Background(), "layup-probe")
	if want := []forge.Rule{{Type: "creation", RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11}, {Type: "deletion", RulesetSourceType: "Repository", RulesetSource: "acme/target", RulesetID: 11}}; err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("EffectiveRules of layup-probe, which no branch of the repository has = %+v, %v; want %+v", got, err, want)
	}
	if n := f.requests["GET /repos/acme/target/branches"]; n != 0 {
		t.Errorf("%d requests of the list of branches; want none", n)
	}
}

// The branch is a multi-segment path parameter of GitHub: each segment is
// escaped and the slashes stay.
func TestEffectiveRulesEscapesTheBranchSegmentBySegment(t *testing.T) {
	f, _, a, _ := setup(t)
	if _, err := a.EffectiveRules(context.Background(), "feature/x#y z"); err != nil {
		t.Fatal(err)
	}
	if want := []string{"/repos/acme/target/rules/branches/feature/x%23y%20z"}; !reflect.DeepEqual(f.rulePaths, want) {
		t.Errorf("the path of the request is %q; want %q", f.rulePaths, want)
	}
}

// A body of null reads as an empty list, as for Comments and Repository; a body that
// is an object, or is not JSON, is an error that names the call and the status.
func TestEffectiveRulesAndActivityReadABodyOfNullAsNoItem(t *testing.T) {
	ctx := context.Background()
	f, _, a, _ := setup(t)
	f.rulesBody, f.activityBody = "null", "null"
	if got, err := a.EffectiveRules(ctx, "main"); err != nil || len(got) != 0 {
		t.Errorf("EffectiveRules with a body of null = %+v, %v; want none", got, err)
	}
	if got, err := a.Activity(ctx, "main"); err != nil || len(got) != 0 {
		t.Errorf("Activity with a body of null = %+v, %v; want none", got, err)
	}
}

func TestEffectiveRulesAndActivityRefuseABodyThatIsNotAList(t *testing.T) {
	ctx := context.Background()
	for _, body := range []string{`{"message":"x"}`, "not json"} {
		f, _, a, _ := setup(t)
		f.rulesBody, f.activityBody = body, body
		for call, run := range map[string]func() error{
			"EffectiveRules": func() error { _, err := a.EffectiveRules(ctx, "main"); return err },
			"Activity":       func() error { _, err := a.Activity(ctx, "main"); return err },
		} {
			var fe *forge.Error
			if err := run(); !errors.As(err, &fe) || fe.Call != call || fe.Status != 200 {
				t.Errorf("%s with the body %q: %v; want a forge.Error of %s, 200", call, body, err, call)
			}
		}
	}
}

// A next page that is not at the api of the forge register is a forge error of the call.
func TestEffectiveRulesRefusesANextPageAtAnotherHost(t *testing.T) {
	f, _, a, _ := setup(t)
	f.rulesNext = "http://elsewhere.invalid/repos/acme/target/rules/branches/main?page=2"
	var fe *forge.Error
	_, err := a.EffectiveRules(context.Background(), "main")
	if !errors.As(err, &fe) || fe.Call != "EffectiveRules" || fe.Status != 200 || !strings.Contains(fe.Message, "not at the api of the forge register") {
		t.Errorf("EffectiveRules with a next page elsewhere: %v; want a forge.Error of EffectiveRules, 200, naming the api", err)
	}
}

// Activity gives the six fields of each item, a null actor as the ID 0 and the login
// ghost, with the ref of the query and one request, whatever the Link header says.
func TestActivityGivesTheEntriesOfARefInOnePage(t *testing.T) {
	f, _, a, _ := setup(t)
	f.activity = []map[string]any{
		{"id": 1, "node_id": "A_1", "before": strings.Repeat("a", 40), "after": strings.Repeat("b", 40), "ref": "refs/heads/main", "timestamp": "2026-10-13T08:00:00Z",
			"activity_type": "push", "actor": map[string]any{"id": 9001, "login": "acme-layup[bot]"}},
		{"id": 2, "node_id": "A_2", "before": strings.Repeat("b", 40), "after": strings.Repeat("c", 40), "ref": "refs/heads/main", "timestamp": "2026-10-13T09:30:00Z",
			"activity_type": "force_push", "actor": nil},
	}
	f.activityLink = true
	got, err := a.Activity(context.Background(), "main")
	want := []forge.ActivityEntry{
		{ActorID: 9001, ActorLogin: "acme-layup[bot]", Ref: "refs/heads/main", Before: strings.Repeat("a", 40), After: strings.Repeat("b", 40),
			Timestamp: time.Date(2026, 10, 13, 8, 0, 0, 0, time.UTC), Type: "push"},
		{ActorID: 0, ActorLogin: "ghost", Ref: "refs/heads/main", Before: strings.Repeat("b", 40), After: strings.Repeat("c", 40),
			Timestamp: time.Date(2026, 10, 13, 9, 30, 0, 0, time.UTC), Type: "force_push"},
	}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("Activity = %+v, %v\nwant %+v", got, err, want)
	}
	if !reflect.DeepEqual(f.activityQueries, []string{"ref=main"}) {
		t.Errorf("the queries are %q; want one request, with the ref and no other parameter", f.activityQueries)
	}
}

// The ref is sent as it is given, escaped as a query value.
func TestActivitySendsTheRefAsGiven(t *testing.T) {
	f, _, a, _ := setup(t)
	if _, err := a.Activity(context.Background(), "refs/heads/feature/x y"); err != nil {
		t.Fatal(err)
	}
	if want := []string{"ref=refs%2Fheads%2Ffeature%2Fx+y"}; !reflect.DeepEqual(f.activityQueries, want) {
		t.Errorf("the query is %q; want %q", f.activityQueries, want)
	}
}

// Repository gives the head of each branch over every page, and the names as before.
func TestRepositoryGivesTheHeadOfEachBranch(t *testing.T) {
	f, _, a, _ := setup(t)
	f.branches = 101
	repo, err := a.Repository(context.Background())
	if err != nil || len(repo.Heads) != 101 || len(repo.Branches) != 101 {
		t.Fatalf("Repository with 101 branches: %d heads, %d names, %v", len(repo.Heads), len(repo.Branches), err)
	}
	if repo.Heads["main"] != f.head(0) || repo.Heads["b050"] != f.head(50) || repo.Heads["b100"] != f.head(100) {
		t.Errorf("the heads of main, b050 and b100 are %q, %q, %q; want %q, %q, %q", repo.Heads["main"], repo.Heads["b050"], repo.Heads["b100"], f.head(0), f.head(50), f.head(100))
	}
	f.branches = 0
	if repo, err := a.Repository(context.Background()); err != nil || len(repo.Heads) != 0 {
		t.Errorf("Repository with no branch has the heads %v, %v; want none", repo.Heads, err)
	}
}

// Each new call names itself and the status in a forge error, with the first line of
// the message, after the handling that every call shares: a rate limit that names no
// reset time fails at once, and a 403 with a reset time but requests left is no limit.
func TestTheCallsOfM2dNameTheCallAndTheStatus(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		call   string
		path   string
		run    func(*Adapter) error
		status int
	}{
		{"EffectiveRules", "GET /repos/acme/target/rules/branches/main", func(a *Adapter) error { _, err := a.EffectiveRules(ctx, "main"); return err }, 422},
		{"Activity", "GET /repos/acme/target/activity", func(a *Adapter) error { _, err := a.Activity(ctx, "main"); return err }, 422},
	} {
		for _, status := range []int{c.status, 404, 500, 429} {
			f, clk, a, _ := setup(t)
			f.reject[c.path] = []int{status}
			var fe *forge.Error
			switch err := c.run(a); {
			case !errors.As(err, &fe):
				t.Errorf("%s with %d: %v; want a forge.Error", c.call, status, err)
			case fe.Call != c.call || fe.Status != status || fe.Message != "rejected by the test":
				t.Errorf("%s with %d: %+v; want the call, the status and the first line of the message", c.call, status, fe)
			case len(clk.sleeps) != 0:
				t.Errorf("%s with %d: the adapter waited %v", c.call, status, clk.sleeps)
			}
		}
		f, clk, a, _ := setup(t)
		f.reject[c.path] = []int{403}
		f.headers[c.path] = http.Header{"X-Ratelimit-Reset": {"1791460800"}, "X-Ratelimit-Remaining": {"5"}}
		var fe *forge.Error
		if err := c.run(a); !errors.As(err, &fe) || fe.Call != c.call || fe.Status != 403 || len(clk.sleeps) != 0 {
			t.Errorf("%s with a 403 that has requests left: %v (waited %v); want a forge.Error of the call, 403, and no wait", c.call, err, clk.sleeps)
		}
	}
}

// Each new call shares the handling of a 401 of the token (one new token and one
// retry, then the error) and of a rate limit (a wait until the reset time, by
// retry-after or by x-ratelimit-reset, with a progress line).
func TestTheCallsOfM2dShareTheHandlingOfA401AndOfARateLimit(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		call string
		path string
		run  func(*Adapter) error
	}{
		{"EffectiveRules", "GET /repos/acme/target/rules/branches/main", func(a *Adapter) error { _, err := a.EffectiveRules(ctx, "main"); return err }},
		{"Activity", "GET /repos/acme/target/activity", func(a *Adapter) error { _, err := a.Activity(ctx, "main"); return err }},
	} {
		f, _, a, _ := setup(t)
		f.reject[c.path] = []int{401}
		if err := c.run(a); err != nil || f.tokens != 2 {
			t.Errorf("%s after one 401: %v, %d tokens; want success and 2 tokens", c.call, err, f.tokens)
		}
		f, _, a, _ = setup(t)
		f.reject[c.path] = []int{401, 401, 401}
		var fe *forge.Error
		if err := c.run(a); !errors.As(err, &fe) || fe.Status != 401 || fe.Call != c.call || f.requests[c.path] != 2 {
			t.Errorf("%s after two 401: %v, %d requests; want a forge.Error of the call, 401, and one retry", c.call, err, f.requests[c.path])
		}
		f, clk, a, lines := setup(t)
		f.reject[c.path] = []int{429}
		f.headers[c.path] = http.Header{"Retry-After": {"25"}}
		if err := c.run(a); err != nil || len(clk.sleeps) != 3 || clk.sleeps[0] != 10*time.Second || clk.sleeps[2] != 5*time.Second || len(*lines) != 3 || !strings.HasPrefix((*lines)[0], c.call+": the rate limit") {
			t.Errorf("%s after a 429 with retry-after 25: %v, sleeps %v, lines %q; want 10 s, 10 s, 5 s and a line of the call each time", c.call, err, clk.sleeps, *lines)
		}
		f, clk, a, _ = setup(t)
		f.reject[c.path] = []int{403}
		f.headers[c.path] = http.Header{"X-Ratelimit-Remaining": {"0"}, "X-Ratelimit-Reset": {strconv.FormatInt(clk.now().Add(12*time.Second).Unix(), 10)}}
		if err := c.run(a); err != nil || len(clk.sleeps) != 2 || clk.sleeps[0] != 10*time.Second || clk.sleeps[1] != 2*time.Second {
			t.Errorf("%s after a 403 with x-ratelimit-reset in 12 s: %v, sleeps %v; want 10 s and 2 s", c.call, err, clk.sleeps)
		}
	}
}

const (
	headSHA      = "1111111111111111111111111111111111111111"
	mergeSHA     = "3333333333333333333333333333333333333333"
	testMergeSHA = "4444444444444444444444444444444444444444"
	pulls        = "/repos/acme/target/pulls"
)

// handle registers the handler of one call of M2e on the fake.
func handle(f *fakeGitHub, call string, h func(http.ResponseWriter, *http.Request)) {
	if f.handlers == nil {
		f.handlers = map[string]func(http.ResponseWriter, *http.Request){}
	}
	f.handlers[call] = h
}

func decode(t *testing.T, r *http.Request, v any) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		t.Errorf("the body of %s %s is not JSON: %v", r.Method, r.URL.Path, err)
	}
}

// prJSON is a pull request as GET /pulls/{number} gives it: merged is a boolean, and
// merge_commit_sha is, before a merge, a test merge commit (GitHub's page, "Get a
// pull request"; runs/T-cr8j/api-sources.md).
func prJSON(number int, draft, merged bool) map[string]any {
	return map[string]any{"number": number, "node_id": fmt.Sprintf("PR_node%d", number), "state": "open", "draft": draft,
		"merged": merged, "merge_commit_sha": testMergeSHA, "head": map[string]any{"sha": headSHA, "ref": "task/T-aaaa/1"}}
}

// listItemJSON is a pull request as GET /pulls gives it: no merged, a merged_at that
// is null until the merge (GitHub's page, "List pull requests").
func listItemJSON(number int, draft bool, mergedAt any, mergeSHA string) map[string]any {
	state := "open"
	if mergedAt != nil {
		state = "closed" // a merged pull request is closed
	}
	return map[string]any{"number": number, "node_id": fmt.Sprintf("PR_node%d", number), "state": state, "draft": draft,
		"merged_at": mergedAt, "merge_commit_sha": mergeSHA, "head": map[string]any{"sha": headSHA, "ref": "task/T-aaaa/1"}}
}

func TestOpenDraftPostsADraftPullRequestAndGivesIt(t *testing.T) {
	f, _, a, _ := setup(t)
	var got map[string]any
	handle(f, "POST "+pulls, func(w http.ResponseWriter, r *http.Request) {
		decode(t, r, &got)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(prJSON(5, true, false))
	})
	pr, err := a.OpenDraft(context.Background(), "task/T-aaaa/1", "main", "feat: T-aaaa x", "Closes #7")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"title": "feat: T-aaaa x", "head": "task/T-aaaa/1", "base": "main", "body": "Closes #7", "draft": true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the request body = %v, want %v (a draft, and no reviewer)", got, want)
	}
	if pr != (forge.PullRequest{Number: 5, NodeID: "PR_node5", State: "open", Draft: true, HeadSHA: headSHA}) {
		t.Errorf("OpenDraft = %+v", pr)
	}
}

func TestOpenDraftGivesAForgeErrorForAnyOtherStatus(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "POST "+pulls, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		io.WriteString(w, `{"message":"A pull request already exists\nfor acme:task"}`)
	})
	_, err := a.OpenDraft(context.Background(), "task/T-aaaa/1", "main", "t", "")
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "OpenDraft" || fe.Status != 422 || fe.Message != "A pull request already exists" {
		t.Errorf("OpenDraft = %v, want the call, the status 422 and the first line", err)
	}
}

func TestMarkReadySendsTheMutationWithTheNodeIDAndGivesThePullRequest(t *testing.T) {
	f, _, a, _ := setup(t)
	var in struct {
		Query     string
		Variables map[string]any
	}
	handle(f, "POST /graphql", func(w http.ResponseWriter, r *http.Request) {
		decode(t, r, &in)
		io.WriteString(w, `{"data":{"markPullRequestReadyForReview":{"pullRequest":{"id":"PR_node5","number":5,"isDraft":false,"state":"OPEN","headRefOid":"`+headSHA+`"}}}}`)
	})
	pr, err := a.MarkReady(context.Background(), "PR_node5")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(in.Query, "markPullRequestReadyForReview") || in.Variables["id"] != "PR_node5" {
		t.Errorf("the request = %+v, want the mutation and the node ID as a variable, never in the text", in)
	}
	if strings.Contains(in.Query, "PR_node5") {
		t.Errorf("the node ID is in the text of the mutation: %q", in.Query)
	}
	if pr != (forge.PullRequest{Number: 5, NodeID: "PR_node5", State: "open", HeadSHA: headSHA}) {
		t.Errorf("MarkReady = %+v, want the state and the head that the answer holds", pr)
	}
	closed := `{"data":{"markPullRequestReadyForReview":{"pullRequest":{"id":"PR_node5","number":5,"isDraft":false,"state":"CLOSED","headRefOid":"` + headSHA + `"}}}}`
	handle(f, "POST /graphql", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, closed) })
	if pr, err := a.MarkReady(context.Background(), "PR_node5"); err != nil || pr.State != "closed" || pr.Merged {
		t.Errorf("MarkReady of a closed pull request = %+v, %v, want the state closed", pr, err)
	}
	if !strings.Contains(in.Query, "state") || !strings.Contains(in.Query, "headRefOid") {
		t.Errorf("the mutation does not select the state and the head: %q", in.Query)
	}
}

func TestMarkReadyReadsTheMergedStateAndRefusesAnAnswerWithNoPullRequest(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "POST /graphql", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":{"markPullRequestReadyForReview":{"pullRequest":{"id":"PR_node5","number":5,"isDraft":false,"state":"MERGED","headRefOid":"`+headSHA+`"}}}}`)
	})
	if pr, err := a.MarkReady(context.Background(), "PR_node5"); err != nil || !pr.Merged || pr.State != "closed" {
		t.Errorf("MarkReady of a merged pull request = %+v, %v, want merged and closed", pr, err)
	}
}

func TestMarkReadyGivesAForgeErrorForAnErrorsArrayOrAnEmptyAnswer(t *testing.T) {
	for name, body := range map[string]string{
		"data null":  `{"data":null}`,
		"no payload": `{"data":{"markPullRequestReadyForReview":{"pullRequest":null}}}`,
		"errors":     `{"errors":[{"message":"Could not resolve to a node\nwith the global id"}],"data":null}`,
		"empty":      `{"data":{"markPullRequestReadyForReview":null}}`,
	} {
		f, _, a, _ := setup(t)
		handle(f, "POST /graphql", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, body) })
		_, err := a.MarkReady(context.Background(), "PR_node5")
		var fe *forge.Error
		if !errors.As(err, &fe) || fe.Call != "MarkReady" || fe.Status != 200 {
			t.Errorf("%s: MarkReady = %v, want a forge error of the call MarkReady", name, err)
		}
		if name == "errors" && fe != nil && fe.Message != "Could not resolve to a node" {
			t.Errorf("the message = %q, want the first line of the first error", fe.Message)
		}
	}
}

func TestSetStatusPostsToTheSHAWithItsStateContextAndDescription(t *testing.T) {
	f, _, a, _ := setup(t)
	var got map[string]any
	handle(f, "POST /repos/acme/target/statuses/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		decode(t, r, &got)
		w.WriteHeader(201)
		io.WriteString(w, `{"state":"success"}`)
	})
	if err := a.SetStatus(context.Background(), headSHA, "layup/gates", "success", "clear: no product path"); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{"state": "success", "context": "layup/gates", "description": "clear: no product path"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("the request body = %v, want %v", got, want)
	}
}

func TestSetStatusRefusesBeforeAnyRequestAWrongStateADescriptionTooLongOrASHAThatIsNone(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "POST /repos/acme/target/statuses/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		io.WriteString(w, `{}`)
	})
	ctx := context.Background()
	for _, c := range []struct{ sha, state, desc string }{
		{headSHA, "ok", "x"},
		{headSHA, "success", strings.Repeat("x", 141)},
		{"main", "success", "x"},
		{headSHA + "/../x", "success", "x"},
		{strings.ToUpper(headSHA[:39]) + "A", "success", "x"},
	} {
		if err := a.SetStatus(ctx, c.sha, "layup/gates", c.state, c.desc); err == nil {
			t.Errorf("SetStatus(%q, %q, %d chars) gave no error", c.sha, c.state, len(c.desc))
		}
	}
	if n := len(f.requests); n != 0 {
		t.Errorf("a refused status made %d requests: %v", n, f.requests)
	}
	if err := a.SetStatus(ctx, headSHA, "layup/gates", "success", strings.Repeat("é", 140)); err != nil {
		t.Errorf("140 characters of two bytes = %v: the limit counts characters", err)
	}
}

func TestSetStatusGivesAForgeErrorForAnyOtherStatus(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "POST /repos/acme/target/statuses/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		io.WriteString(w, `{"message":"Resource not accessible by integration"}`)
	})
	err := a.SetStatus(context.Background(), headSHA, "layup/gates", "success", "x")
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "SetStatus" || fe.Status != 403 {
		t.Errorf("SetStatus = %v, want a forge error of the call SetStatus, status 403", err)
	}
}

func TestStatusesGivesEveryPageOfStatusesAndOfCheckRuns(t *testing.T) {
	f, _, a, _ := setup(t)
	var statusQueries, runQueries []string
	handle(f, "GET /repos/acme/target/commits/"+headSHA+"/statuses", func(w http.ResponseWriter, r *http.Request) {
		statusQueries = append(statusQueries, r.URL.RawQuery)
		if r.URL.Query().Get("page") == "" {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s%s?per_page=100&page=2>; rel="next"`, r.Host, r.URL.Path))
			json.NewEncoder(w).Encode([]map[string]any{{"context": "layup/gates", "state": "success", "description": "pass",
				"creator": map[string]any{"id": 77, "login": "layup-agent[bot]"}}})
			return
		}
		json.NewEncoder(w).Encode([]map[string]any{{"context": "layup/gates", "state": "pending", "description": nil, "creator": nil}})
	})
	handle(f, "GET /repos/acme/target/commits/"+headSHA+"/check-runs", func(w http.ResponseWriter, r *http.Request) {
		runQueries = append(runQueries, r.URL.RawQuery)
		if r.URL.Query().Get("page") == "" {
			w.Header().Set("Link", fmt.Sprintf(`<http://%s%s?filter=latest&per_page=100&page=2>; rel="next"`, r.Host, r.URL.Path))
			json.NewEncoder(w).Encode(map[string]any{"total_count": 2, "check_runs": []map[string]any{{"name": "layout", "status": "completed",
				"conclusion": "success", "app": map[string]any{"id": 15368, "slug": "github-actions"}}}})
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"total_count": 2, "check_runs": []map[string]any{{"name": "boundary", "status": "in_progress",
			"conclusion": nil, "app": nil}}})
	})
	got, err := a.Statuses(context.Background(), headSHA)
	if err != nil {
		t.Fatal(err)
	}
	wantS := []forge.Status{
		{Context: "layup/gates", State: "success", Description: "pass", CreatorID: 77, CreatorLogin: "layup-agent[bot]"},
		{Context: "layup/gates", State: "pending", CreatorLogin: "ghost"},
	}
	wantR := []forge.CheckRun{
		{Name: "layout", Status: "completed", Conclusion: "success", AppID: 15368, AppSlug: "github-actions"},
		{Name: "boundary", Status: "in_progress"},
	}
	if !reflect.DeepEqual(got.Statuses, wantS) {
		t.Errorf("Statuses = %+v, want %+v (the order the forge gives, a null creator is ghost)", got.Statuses, wantS)
	}
	if !reflect.DeepEqual(got.CheckRuns, wantR) {
		t.Errorf("CheckRuns = %+v, want %+v", got.CheckRuns, wantR)
	}
	if len(statusQueries) != 2 || !strings.Contains(statusQueries[0], "per_page=100") {
		t.Errorf("the status requests = %v, want two pages of 100", statusQueries)
	}
	if len(runQueries) != 2 || !strings.Contains(runQueries[0], "filter=latest") || !strings.Contains(runQueries[0], "per_page=100") {
		t.Errorf("the check-run requests = %v, want two pages, the latest run of each name, 100 a page", runQueries)
	}
}

func TestStatusesRefusesANextPageAtAnotherHostOfEitherList(t *testing.T) {
	for _, which := range []string{"statuses", "check-runs"} {
		f, _, a, _ := setup(t)
		handle(f, "GET /repos/acme/target/commits/"+headSHA+"/statuses", func(w http.ResponseWriter, r *http.Request) {
			if which == "statuses" {
				w.Header().Set("Link", `<http://evil.example/x?page=2>; rel="next"`)
			}
			io.WriteString(w, `[]`)
		})
		handle(f, "GET /repos/acme/target/commits/"+headSHA+"/check-runs", func(w http.ResponseWriter, r *http.Request) {
			if which == "check-runs" {
				w.Header().Set("Link", `<http://evil.example/x?page=2>; rel="next"`)
			}
			io.WriteString(w, `{"total_count":0,"check_runs":[]}`)
		})
		_, err := a.Statuses(context.Background(), headSHA)
		var fe *forge.Error
		if !errors.As(err, &fe) || fe.Call != "Statuses" || !strings.Contains(fe.Message, "next page") {
			t.Errorf("%s: Statuses = %v, want a forge error that names the next page", which, err)
		}
	}
}

func TestPullRequestByNumberAndByBranch(t *testing.T) {
	f, _, a, _ := setup(t)
	var queries []string
	merged := prJSON(5, false, true)
	merged["state"], merged["merge_commit_sha"] = "closed", mergeSHA
	handle(f, "GET "+pulls+"/5", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(merged) })
	handle(f, "GET "+pulls+"/6", func(w http.ResponseWriter, r *http.Request) { json.NewEncoder(w).Encode(prJSON(6, true, false)) })
	handle(f, "GET "+pulls, func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		switch r.URL.Query().Get("head") {
		case "acme:task/T-aaaa/1": // the newest is open, with a test merge commit
			json.NewEncoder(w).Encode([]map[string]any{listItemJSON(6, true, nil, testMergeSHA), listItemJSON(5, false, "2026-10-15T12:00:00Z", mergeSHA)})
		case "acme:task/T-bbbb/1": // the newest is merged: only merged_at says so
			json.NewEncoder(w).Encode([]map[string]any{listItemJSON(7, false, "2026-10-15T12:00:00Z", mergeSHA)})
		default:
			io.WriteString(w, `[]`)
		}
	})
	ctx := context.Background()
	pr, found, err := a.PullRequest(ctx, forge.PullRequestQuery{Number: 5})
	if err != nil || !found || pr != (forge.PullRequest{Number: 5, NodeID: "PR_node5", State: "closed", Merged: true, HeadSHA: headSHA, MergeSHA: mergeSHA}) {
		t.Errorf("by number = %+v, %v, %v", pr, found, err)
	}
	pr, found, err = a.PullRequest(ctx, forge.PullRequestQuery{Branch: "task/T-aaaa/1"})
	if err != nil || !found || pr != (forge.PullRequest{Number: 6, NodeID: "PR_node6", State: "open", Draft: true, HeadSHA: headSHA}) {
		t.Errorf("by branch = %+v, %v, %v, want the first of the list, the newest, not merged, with no merge commit (the test merge commit is none)", pr, found, err)
	}
	pr, found, err = a.PullRequest(ctx, forge.PullRequestQuery{Branch: "task/T-bbbb/1"})
	if err != nil || !found || !pr.Merged || pr.MergeSHA != mergeSHA {
		t.Errorf("by branch, merged_at set = %+v, %v, %v, want merged with the merge commit", pr, found, err)
	}
	if pr, _, _ := a.PullRequest(ctx, forge.PullRequestQuery{Number: 6}); pr.Merged || pr.MergeSHA != "" {
		t.Errorf("an open pull request by number = %+v, want no merge commit: merge_commit_sha is a test merge commit before the merge", pr)
	}
	if len(queries) != 2 || !strings.Contains(queries[0], "state=all") {
		t.Errorf("the list request = %v, want head and state=all", queries)
	}
	if _, found, err := a.PullRequest(ctx, forge.PullRequestQuery{Branch: "task/none"}); found || err != nil {
		t.Errorf("a branch with no pull request: found = %v, err = %v, want a result and no error", found, err)
	}
}

func TestStatusesRefusesAHeadThatIsNoFullSHABeforeAnyRequest(t *testing.T) {
	f, _, a, _ := setup(t)
	for _, sha := range []string{"", "main", headSHA[:39], headSHA + "0", strings.ToUpper(headSHA[:39]) + "A"} {
		if _, err := a.Statuses(context.Background(), sha); err == nil {
			t.Errorf("Statuses(%q) gave no error", sha)
		}
	}
	if n := len(f.requests); n != 0 {
		t.Errorf("a refused head made %d requests: %v", n, f.requests)
	}
}

func TestPullRequestRefusesAQueryThatNamesNeitherOrBothAndGivesAnErrorForANumberItDoesNotKnow(t *testing.T) {
	f, _, a, _ := setup(t)
	ctx := context.Background()
	for _, q := range []forge.PullRequestQuery{{}, {Number: 5, Branch: "b"}, {Number: -1}} {
		if _, _, err := a.PullRequest(ctx, q); err == nil {
			t.Errorf("PullRequest(%+v) gave no error", q)
		}
	}
	if n := len(f.requests); n != 0 {
		t.Errorf("a refused query made %d requests: %v", n, f.requests)
	}
	_, _, err := a.PullRequest(ctx, forge.PullRequestQuery{Number: 404})
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "PullRequest" || fe.Status != 404 {
		t.Errorf("a number that is none: %v, want a forge error, status 404", err)
	}
}

func TestMergeSendsTheHeadSHAAndAMergeCommitAndGivesTheMergeCommit(t *testing.T) {
	f, _, a, _ := setup(t)
	var got map[string]any
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		decode(t, r, &got)
		json.NewEncoder(w).Encode(map[string]any{"sha": mergeSHA, "merged": true, "message": "Pull Request successfully merged"})
	})
	sha, err := a.Merge(context.Background(), 5, headSHA)
	if err != nil || sha != mergeSHA {
		t.Fatalf("Merge = %q, %v", sha, err)
	}
	if want := map[string]any{"sha": headSHA, "merge_method": "merge"}; !reflect.DeepEqual(got, want) {
		t.Errorf("the request body = %v, want %v (a plain merge, never a squash)", got, want)
	}
}

func TestMergeOfAHeadThatMovedIsAForgeErrorOfStatus409(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(409)
		io.WriteString(w, `{"message":"Head branch was modified. Review and try the merge again."}`)
	})
	_, err := a.Merge(context.Background(), 5, headSHA)
	var fe *forge.Error
	if !errors.As(err, &fe) || fe.Call != "Merge" || fe.Status != 409 || !strings.HasPrefix(fe.Message, "Head branch was modified") {
		t.Errorf("Merge = %v, want a forge error of the call Merge, status 409", err)
	}
}

func TestMergeThatAnswersNotMergedIsNeverASuccess(t *testing.T) {
	f, _, a, _ := setup(t)
	sha := ""
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"sha": sha, "merged": false, "message": "Pull Request is not mergeable"})
	})
	for _, sha = range []string{"", mergeSHA} { // merged false with a SHA is no merge either
		if got, err := a.Merge(context.Background(), 5, headSHA); err == nil {
			t.Errorf("Merge = %q, nil for an answer with merged false and the SHA %q", got, sha)
		}
	}
}

func TestMergeRefusesBeforeAnyRequestAPullRequestBelow1OrAHeadThatIsNoFullSHA(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"sha": mergeSHA, "merged": true}) // a call that got through would succeed
	})
	handle(f, "PUT "+pulls+"/0/merge", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"sha": mergeSHA, "merged": true})
	})
	for _, c := range []struct {
		n   int
		sha string
	}{{5, "main"}, {5, ""}, {5, headSHA[:39]}, {0, headSHA}, {-1, headSHA}} {
		if sha, err := a.Merge(context.Background(), c.n, c.sha); err == nil {
			t.Errorf("Merge(%d, %q) = %q, nil: it is a refusal", c.n, c.sha, sha)
		}
	}
	if n := len(f.requests); n != 0 {
		t.Errorf("a refused merge made %d requests: %v", n, f.requests)
	}
}

func TestTheCallsOfM2eShareTheHandlingOfA401AndOfARateLimit(t *testing.T) {
	f, c, a, lines := setup(t)
	handle(f, "POST /repos/acme/target/statuses/"+headSHA, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		io.WriteString(w, `{}`)
	})
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"sha": mergeSHA, "merged": true})
	})
	ctx := context.Background()
	if err := a.SetStatus(ctx, headSHA, "layup/gates", "success", "x"); err != nil { // a token first
		t.Fatal(err)
	}
	f.reject["POST /repos/acme/target/statuses/"+headSHA] = []int{401}
	if err := a.SetStatus(ctx, headSHA, "layup/gates", "success", "x"); err != nil {
		t.Errorf("SetStatus after a 401 = %v, want one new token and one retry", err)
	}
	if f.tokens != 2 {
		t.Errorf("tokens made = %d, want 2", f.tokens)
	}
	f.reject["PUT "+pulls+"/5/merge"] = []int{429}
	f.headers["PUT "+pulls+"/5/merge"] = http.Header{"Retry-After": []string{"25"}}
	if _, err := a.Merge(ctx, 5, headSHA); err != nil {
		t.Errorf("Merge after a rate limit = %v", err)
	}
	if len(c.sleeps) != 3 || len(*lines) != 3 || !strings.HasPrefix((*lines)[0], "Merge: the rate limit") {
		t.Errorf("the wait: sleeps %v, lines %v, want three waits of at most ten seconds", c.sleeps, *lines)
	}
}

func TestNoTokenIsInTheMessageOfAnErrorOfM2e(t *testing.T) {
	f, _, a, _ := setup(t)
	handle(f, "PUT "+pulls+"/5/merge", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		io.WriteString(w, `{"message":"boom"}`)
	})
	_, err := a.Merge(context.Background(), 5, headSHA)
	if err == nil {
		t.Fatal("no error")
	}
	noToken(t, err.Error())
}
