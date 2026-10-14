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
// lacks the three fields gives their zero values, and a ruleset_id that is no integer is an error of the call.
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
