// Package github is the GitHub adapter of the forge interface
// (docs/spec/forge.md): the JWT of the App, the installation token, the calls
// of M2a and the rules of a forge error. It is the one package of LAYUP that
// imports a network package (rule 5 of docs/spec/packages.md). It takes the key
// that forge.CheckKeyFile gives, so it parses no key file.
package github

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pharzam/layup/internal/forge"
)

// Config is what the adapter needs: the api of the forge register, the App ID,
// the App's key, the repository OWNER/NAME, the HTTP client, the function that
// prints a progress line, and the clock and the sleep, which a test replaces.
type Config struct {
	API      string
	AppID    int64
	Key      *rsa.PrivateKey
	Owner    string
	Name     string
	Client   *http.Client
	Progress func(string)
	Now      func() time.Time
	Sleep    func(context.Context, time.Duration) error
}

// Adapter is the GitHub adapter. It holds the installation and its token.
type Adapter struct {
	cfg  Config
	inst int64
	tok  forge.Token
}

var _ forge.Forge = (*Adapter)(nil)

// New gives an adapter of cfg; an empty Now or Sleep is the real one. With no
// Client, it makes one that follows no redirect, so a token goes to the api
// of the forge register only, and whose transport reads no proxy of the
// environment (forge.md, The App identity; task T-mqty).
func New(cfg Config) *Adapter {
	if cfg.Client == nil {
		tr := http.DefaultTransport.(*http.Transport).Clone()
		tr.Proxy = nil
		cfg.Client = &http.Client{Transport: tr, CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	if cfg.Sleep == nil {
		cfg.Sleep = func(ctx context.Context, d time.Duration) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(d):
				return nil
			}
		}
	}
	cfg.API = strings.TrimSuffix(cfg.API, "/")
	return &Adapter{cfg: cfg}
}

// Capabilities declares the six (forge.md: "The adapter declares the six").
func (a *Adapter) Capabilities() []forge.Capability {
	return append([]forge.Capability{}, forge.Capabilities...)
}

func (a *Adapter) repo() string {
	return "/repos/" + url.PathEscape(a.cfg.Owner) + "/" + url.PathEscape(a.cfg.Name)
}

// Installation gives the installation of the App on the repository (JWT).
func (a *Adapter) Installation(ctx context.Context) (forge.Installation, error) {
	var out struct {
		ID          int64             `json:"id"`
		Permissions map[string]string `json:"permissions"`
	}
	if _, err := a.do(ctx, "Installation", "GET", a.cfg.API+a.repo()+"/installation", byJWT, nil, 200, &out); err != nil {
		return forge.Installation{}, err
	}
	a.inst = out.ID
	return forge.Installation{ID: out.ID, Permissions: out.Permissions}, nil
}

// Token gives the token it holds while five minutes or more of it remain, else
// a new one (forge.md, The App identity).
func (a *Adapter) Token(ctx context.Context) (forge.Token, error) {
	if a.tok.Token != "" && a.tok.Expires.Sub(a.cfg.Now()) >= 5*time.Minute {
		return a.tok, nil
	}
	if a.inst == 0 {
		if _, err := a.Installation(ctx); err != nil {
			return forge.Token{}, err
		}
	}
	var out struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	path := fmt.Sprintf("%s/app/installations/%d/access_tokens", a.cfg.API, a.inst)
	if _, err := a.do(ctx, "Token", "POST", path, byJWT, nil, 201, &out); err != nil {
		return forge.Token{}, err
	}
	if out.Token == "" {
		return forge.Token{}, &forge.Error{Call: "Token", Status: 201, Message: "the response holds no token"}
	}
	a.tok = forge.Token{Token: out.Token, Expires: out.ExpiresAt}
	return a.tok, nil
}

// Repository gives the default branch, the visibility, and the names of the
// branches, every page; the repository has a commit when the list is not
// empty.
func (a *Adapter) Repository(ctx context.Context) (forge.Repository, error) {
	var repo struct {
		DefaultBranch string `json:"default_branch"`
		Visibility    string `json:"visibility"`
	}
	if _, err := a.do(ctx, "Repository", "GET", a.cfg.API+a.repo(), byToken, nil, 200, &repo); err != nil {
		return forge.Repository{}, err
	}
	var names []string
	heads := map[string]string{}
	for next := a.cfg.API + a.repo() + "/branches?per_page=100"; next != ""; {
		var page []struct {
			Name   string `json:"name"`
			Commit struct {
				SHA string `json:"sha"`
			} `json:"commit"`
		}
		h, err := a.do(ctx, "Repository", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return forge.Repository{}, err
		}
		for _, b := range page {
			names = append(names, b.Name)
			heads[b.Name] = b.Commit.SHA
		}
		if next = nextPage(h.Get("Link")); next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return forge.Repository{}, &forge.Error{Call: "Repository", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	return forge.Repository{DefaultBranch: repo.DefaultBranch, Visibility: repo.Visibility, HasCommit: len(names) > 0, Branches: names, Heads: heads}, nil
}

// UserID gives the numeric ID of a login.
func (a *Adapter) UserID(ctx context.Context, login string) (int64, error) {
	var out struct {
		ID int64 `json:"id"`
	}
	if _, err := a.do(ctx, "UserID", "GET", a.cfg.API+"/users/"+url.PathEscape(login), byToken, nil, 200, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

// OpenIssue opens an issue and gives its number.
func (a *Adapter) OpenIssue(ctx context.Context, title, body string) (int, error) {
	in, err := json.Marshal(map[string]string{"title": title, "body": body})
	if err != nil {
		return 0, err
	}
	var out struct {
		Number int `json:"number"`
	}
	if _, err := a.do(ctx, "OpenIssue", "POST", a.cfg.API+a.repo()+"/issues", byToken, in, 201, &out); err != nil {
		return 0, err
	}
	return out.Number, nil
}

// Comment posts a comment on an issue and gives its ID.
func (a *Adapter) Comment(ctx context.Context, issue int, body string) (int64, error) {
	in, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return 0, err
	}
	var out struct {
		ID int64 `json:"id"`
	}
	if _, err := a.do(ctx, "Comment", "POST", fmt.Sprintf("%s%s/issues/%d/comments", a.cfg.API, a.repo(), issue), byToken, in, 201, &out); err != nil {
		return 0, err
	}
	return out.ID, nil
}

// Comments gives each comment of an issue, every page. A comment of a deleted
// account (user null) has the author ID 0 and the login ghost.
func (a *Adapter) Comments(ctx context.Context, issue int) ([]forge.Comment, error) {
	next := fmt.Sprintf("%s%s/issues/%d/comments?per_page=100", a.cfg.API, a.repo(), issue)
	var all []forge.Comment
	for next != "" {
		var page []struct {
			ID   int64 `json:"id"`
			User *struct {
				ID    int64  `json:"id"`
				Login string `json:"login"`
			} `json:"user"`
			App *struct {
				Slug string `json:"slug"`
			} `json:"performed_via_github_app"`
			Created time.Time `json:"created_at"`
			Updated time.Time `json:"updated_at"`
			Body    string    `json:"body"`
		}
		h, err := a.do(ctx, "Comments", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return nil, err
		}
		for _, c := range page {
			com := forge.Comment{ID: c.ID, AuthorLogin: "ghost", Created: c.Created, Updated: c.Updated, Body: c.Body}
			if c.User != nil {
				com.AuthorID, com.AuthorLogin = c.User.ID, c.User.Login
			}
			if c.App != nil {
				com.App = c.App.Slug
			}
			all = append(all, com)
		}
		next = nextPage(h.Get("Link"))
		if next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return nil, &forge.Error{Call: "Comments", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	return all, nil
}

// EffectiveRules gives each active rule that applies to a branch, every page
// (forge.md, The calls of M2d; tasks T-jpkc and T-gkvi): GET .../rules/branches/{branch},
// 100 a page, each rule with its type, its parameters and the ruleset that it is of
// (ruleset_source_type, ruleset_source and ruleset_id, the zero value when absent). The branch is a multi-segment parameter of GitHub, so each of its
// segments is escaped and its slashes stay. The endpoint answers by the name of
// the branch: a branch with no rule gives an empty list, and a branch that does not
// exist yet, as layup-probe before its probe, gives the rules of the rulesets whose
// conditions name it. The call asks for no list of branches and adds no filter of
// its own: the endpoint omits a ruleset in evaluate or disabled.
func (a *Adapter) EffectiveRules(ctx context.Context, branch string) ([]forge.Rule, error) {
	segments := strings.Split(branch, "/")
	for i, s := range segments {
		segments[i] = url.PathEscape(s)
	}
	next := a.cfg.API + a.repo() + "/rules/branches/" + strings.Join(segments, "/") + "?per_page=100"
	var all []forge.Rule
	for next != "" {
		var page []struct {
			Type              string         `json:"type"`
			Parameters        map[string]any `json:"parameters"`
			RulesetSourceType string         `json:"ruleset_source_type"`
			RulesetSource     string         `json:"ruleset_source"`
			RulesetID         int64          `json:"ruleset_id"`
		}
		h, err := a.do(ctx, "EffectiveRules", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return nil, err
		}
		for _, r := range page {
			rule := forge.Rule{Type: r.Type, RulesetSourceType: r.RulesetSourceType, RulesetSource: r.RulesetSource, RulesetID: r.RulesetID}
			if len(r.Parameters) > 0 {
				rule.Parameters = r.Parameters
			}
			all = append(all, rule)
		}
		next = nextPage(h.Get("Link"))
		if next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return nil, &forge.Error{Call: "EffectiveRules", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	return all, nil
}

// Activity gives the activities of a ref, one page of the default size
// (forge.md, The calls of M2d; task T-jpkc): GET .../activity with the ref as it
// is given, a branch name or a full ref, and no other parameter. A null actor is
// the ID 0 and the login ghost, as for Comments.
func (a *Adapter) Activity(ctx context.Context, ref string) ([]forge.ActivityEntry, error) {
	var page []struct {
		Actor *struct {
			ID    int64  `json:"id"`
			Login string `json:"login"`
		} `json:"actor"`
		Ref       string    `json:"ref"`
		Before    string    `json:"before"`
		After     string    `json:"after"`
		Timestamp time.Time `json:"timestamp"`
		Type      string    `json:"activity_type"`
	}
	if _, err := a.do(ctx, "Activity", "GET", a.cfg.API+a.repo()+"/activity?ref="+url.QueryEscape(ref), byToken, nil, 200, &page); err != nil {
		return nil, err
	}
	var all []forge.ActivityEntry
	for _, e := range page {
		entry := forge.ActivityEntry{ActorLogin: "ghost", Ref: e.Ref, Before: e.Before, After: e.After, Timestamp: e.Timestamp, Type: e.Type}
		if e.Actor != nil {
			entry.ActorID, entry.ActorLogin = e.Actor.ID, e.Actor.Login
		}
		all = append(all, entry)
	}
	return all, nil
}

// The calls of M2e (forge.md, The calls of M2e; task T-cr8j).

// pullWire is a pull request as the REST API gives it, in a list too (a list
// item has merged_at and no merged).
type pullWire struct {
	Number   int     `json:"number"`
	NodeID   string  `json:"node_id"`
	State    string  `json:"state"`
	Draft    bool    `json:"draft"`
	Merged   bool    `json:"merged"`
	Body     *string `json:"body"`
	MergedAt *string `json:"merged_at"`
	MergeSHA *string `json:"merge_commit_sha"`
	Head     struct {
		SHA string `json:"sha"`
	} `json:"head"`
}

func (w pullWire) pull() forge.PullRequest {
	pr := forge.PullRequest{Number: w.Number, NodeID: w.NodeID, State: w.State, Draft: w.Draft, HeadSHA: w.Head.SHA,
		Merged: w.Merged || w.MergedAt != nil}
	// Before a merge merge_commit_sha holds a test merge commit (GitHub's page), which
	// is no merge commit: the field is the merge commit's SHA only when it is merged.
	if pr.Merged && w.MergeSHA != nil {
		pr.MergeSHA = *w.MergeSHA
	}
	if w.Body != nil { // a null body is the empty string
		pr.Body = *w.Body
	}
	return pr
}

// OpenDraft opens a draft pull request of head onto base and requests no review.
func (a *Adapter) OpenDraft(ctx context.Context, head, base, title, body string) (forge.PullRequest, error) {
	in, err := json.Marshal(map[string]any{"title": title, "head": head, "base": base, "body": body, "draft": true})
	if err != nil {
		return forge.PullRequest{}, err
	}
	var out pullWire
	if _, err := a.do(ctx, "OpenDraft", "POST", a.cfg.API+a.repo()+"/pulls", byToken, in, 201, &out); err != nil {
		return forge.PullRequest{}, err
	}
	return out.pull(), nil
}

// MarkReady takes a pull request out of the draft state with the GraphQL
// mutation markPullRequestReadyForReview, which the REST API has no call for. The
// node ID is a variable, never part of the text of the mutation. An answer with an
// errors array, or with no pull request, is a forge error.
func (a *Adapter) MarkReady(ctx context.Context, nodeID string) (forge.PullRequest, error) {
	in, err := json.Marshal(map[string]any{
		"query":     "mutation($id: ID!) { markPullRequestReadyForReview(input: {pullRequestId: $id}) { pullRequest { id number isDraft state headRefOid } } }",
		"variables": map[string]string{"id": nodeID}})
	if err != nil {
		return forge.PullRequest{}, err
	}
	var out struct {
		Data *struct {
			Mark *struct {
				PullRequest *struct {
					ID         string `json:"id"`
					Number     int    `json:"number"`
					IsDraft    bool   `json:"isDraft"`
					State      string `json:"state"`
					HeadRefOid string `json:"headRefOid"`
				} `json:"pullRequest"`
			} `json:"markPullRequestReadyForReview"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if _, err := a.do(ctx, "MarkReady", "POST", graphqlURL(a.cfg.API), byToken, in, 200, &out); err != nil {
		return forge.PullRequest{}, err
	}
	if len(out.Errors) > 0 {
		return forge.PullRequest{}, &forge.Error{Call: "MarkReady", Status: 200, Message: firstLine([]byte(out.Errors[0].Message))}
	}
	if out.Data == nil || out.Data.Mark == nil || out.Data.Mark.PullRequest == nil {
		return forge.PullRequest{}, &forge.Error{Call: "MarkReady", Status: 200, Message: "the answer holds no pull request"}
	}
	p := out.Data.Mark.PullRequest
	state := "closed"
	if p.State == "OPEN" {
		state = "open"
	}
	return forge.PullRequest{Number: p.Number, NodeID: p.ID, State: state, Draft: p.IsDraft, HeadSHA: p.HeadRefOid, Merged: p.State == "MERGED"}, nil
}

// SetStatus posts a commit status on a SHA. The state, the length of the
// description and the form of the SHA are checked before any request.
func (a *Adapter) SetStatus(ctx context.Context, sha, context, state, description string) error {
	if !forge.FullSHA(sha) {
		return &forge.Error{Call: "SetStatus", Message: "the SHA is not a full object ID"}
	}
	if err := forge.ValidStatus(state, description); err != nil {
		return &forge.Error{Call: "SetStatus", Message: err.Error()}
	}
	in, err := json.Marshal(map[string]string{"state": state, "context": context, "description": description})
	if err != nil {
		return err
	}
	var out struct{}
	_, err = a.do(ctx, "SetStatus", "POST", a.cfg.API+a.repo()+"/statuses/"+sha, byToken, in, 201, &out)
	return err
}

// Statuses gives each status of a head, in the order of the forge (newest first),
// and each latest check run, every page of each. The next page of a list is at the
// api of the forge register only.
func (a *Adapter) Statuses(ctx context.Context, sha string) (forge.HeadChecks, error) {
	if !forge.FullSHA(sha) {
		return forge.HeadChecks{}, &forge.Error{Call: "Statuses", Message: "the SHA is not a full object ID"}
	}
	var out forge.HeadChecks
	next := a.cfg.API + a.repo() + "/commits/" + sha + "/statuses?per_page=100"
	for next != "" {
		var page []struct {
			Context     string  `json:"context"`
			State       string  `json:"state"`
			Description *string `json:"description"`
			Creator     *struct {
				ID    int64  `json:"id"`
				Login string `json:"login"`
			} `json:"creator"`
		}
		h, err := a.do(ctx, "Statuses", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return forge.HeadChecks{}, err
		}
		for _, s := range page {
			st := forge.Status{Context: s.Context, State: s.State, CreatorLogin: "ghost"}
			if s.Description != nil {
				st.Description = *s.Description
			}
			if s.Creator != nil {
				st.CreatorID, st.CreatorLogin = s.Creator.ID, s.Creator.Login
			}
			out.Statuses = append(out.Statuses, st)
		}
		if next = nextPage(h.Get("Link")); next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return forge.HeadChecks{}, &forge.Error{Call: "Statuses", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	next = a.cfg.API + a.repo() + "/commits/" + sha + "/check-runs?filter=latest&per_page=100"
	for next != "" {
		var page struct {
			CheckRuns []struct {
				Name       string  `json:"name"`
				Status     string  `json:"status"`
				Conclusion *string `json:"conclusion"`
				App        *struct {
					ID   int64  `json:"id"`
					Slug string `json:"slug"`
				} `json:"app"`
			} `json:"check_runs"`
		}
		h, err := a.do(ctx, "Statuses", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return forge.HeadChecks{}, err
		}
		for _, c := range page.CheckRuns {
			run := forge.CheckRun{Name: c.Name, Status: c.Status}
			if c.Conclusion != nil {
				run.Conclusion = *c.Conclusion
			}
			if c.App != nil {
				run.AppID, run.AppSlug = c.App.ID, c.App.Slug
			}
			out.CheckRuns = append(out.CheckRuns, run)
		}
		if next = nextPage(h.Get("Link")); next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return forge.HeadChecks{}, &forge.Error{Call: "Statuses", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	return out, nil
}

// PullRequest gives a pull request by its number (GET /pulls/{n}), or the
// newest one of a branch (GET /pulls?head=OWNER:BRANCH&state=all, the first item
// of the list); a branch with none is found false and no error. A query that
// names neither or both is refused before any request.
func (a *Adapter) PullRequest(ctx context.Context, q forge.PullRequestQuery) (forge.PullRequest, bool, error) {
	if (q.Number != 0) == (q.Branch != "") || q.Number < 0 {
		return forge.PullRequest{}, false, &forge.Error{Call: "PullRequest", Message: "the query names a number or a branch, one of the two"}
	}
	if q.Number != 0 {
		var out pullWire
		if _, err := a.do(ctx, "PullRequest", "GET", fmt.Sprintf("%s%s/pulls/%d", a.cfg.API, a.repo(), q.Number), byToken, nil, 200, &out); err != nil {
			return forge.PullRequest{}, false, err
		}
		return out.pull(), true, nil
	}
	var list []pullWire
	target := a.cfg.API + a.repo() + "/pulls?state=all&head=" + url.QueryEscape(a.cfg.Owner+":"+q.Branch)
	if _, err := a.do(ctx, "PullRequest", "GET", target, byToken, nil, 200, &list); err != nil {
		return forge.PullRequest{}, false, err
	}
	if len(list) == 0 {
		return forge.PullRequest{}, false, nil
	}
	return list[0].pull(), true, nil
}

// Merge merges a pull request at its head SHA by a merge commit, never a squash,
// and gives the merge commit's SHA. A head that moved is the forge's 409, a forge
// error; an answer with merged false is one too.
func (a *Adapter) Merge(ctx context.Context, number int, sha string) (string, error) {
	if number < 1 || !forge.FullSHA(sha) {
		return "", &forge.Error{Call: "Merge", Message: "the pull request number or the head SHA is not valid"}
	}
	in, err := json.Marshal(map[string]string{"sha": sha, "merge_method": "merge"})
	if err != nil {
		return "", err
	}
	var out struct {
		SHA     string `json:"sha"`
		Merged  bool   `json:"merged"`
		Message string `json:"message"`
	}
	if _, err := a.do(ctx, "Merge", "PUT", fmt.Sprintf("%s%s/pulls/%d/merge", a.cfg.API, a.repo(), number), byToken, in, 200, &out); err != nil {
		return "", err
	}
	if !out.Merged || out.SHA == "" {
		msg := firstLine([]byte(out.Message))
		if out.Message == "" {
			msg = "the pull request is not merged"
		}
		return "", &forge.Error{Call: "Merge", Status: 200, Message: msg}
	}
	return out.SHA, nil
}

// EditBody replaces the body of a pull request (PATCH /pulls/{n} with the field
// body alone, so the title, the state and the base cannot change) and gives the
// pull request. A number below 1 is refused before any request.
func (a *Adapter) EditBody(ctx context.Context, number int, body string) (forge.PullRequest, error) {
	if number < 1 {
		return forge.PullRequest{}, &forge.Error{Call: "EditBody", Message: "the pull request number is not valid"}
	}
	in, err := json.Marshal(map[string]string{"body": body})
	if err != nil {
		return forge.PullRequest{}, err
	}
	var out pullWire
	if _, err := a.do(ctx, "EditBody", "PATCH", fmt.Sprintf("%s%s/pulls/%d", a.cfg.API, a.repo(), number), byToken, in, 200, &out); err != nil {
		return forge.PullRequest{}, err
	}
	return out.pull(), nil
}

// graphqlURL gives the URL of the GraphQL API of a REST api: <api>/graphql, and
// for a GitHub Enterprise Server api that ends in /api/v3, /api/graphql.
func graphqlURL(api string) string {
	if strings.HasSuffix(api, "/api/v3") {
		return strings.TrimSuffix(api, "/v3") + "/graphql"
	}
	return api + "/graphql"
}

// nextPage gives the URL of rel="next" of a Link header, or "".
func nextPage(link string) string {
	for _, part := range strings.Split(link, ",") {
		u, params, ok := strings.Cut(strings.TrimSpace(part), ";")
		if ok && strings.Contains(params, `rel="next"`) {
			return strings.Trim(strings.TrimSpace(u), "<>")
		}
	}
	return ""
}

type auth int

const (
	byJWT auth = iota
	byToken
)

// do makes one call: it sends the request, retries once after a 401 of an
// installation token with a new token, waits out a rate limit, and gives a
// forge.Error for any other status than want or a body it cannot read
// (forge.md, Forge errors).
func (a *Adapter) do(ctx context.Context, call, method, target string, by auth, body []byte, want int, out any) (http.Header, error) {
	retried := false
	for {
		var bearer string
		if by == byJWT {
			jwt, err := makeJWT(a.cfg.Key, a.cfg.AppID, a.cfg.Now())
			if err != nil {
				return nil, &forge.Error{Call: call, Message: "the JWT cannot be made: " + err.Error()}
			}
			bearer = jwt
		} else {
			tok, err := a.Token(ctx)
			if err != nil {
				return nil, err
			}
			bearer = tok.Token
		}
		req, err := http.NewRequestWithContext(ctx, method, target, bytes.NewReader(body))
		if err != nil {
			return nil, &forge.Error{Call: call, Message: err.Error()}
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
		req.Header.Set("Authorization", "Bearer "+bearer)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}
		resp, err := a.cfg.Client.Do(req)
		if err != nil {
			return nil, &forge.Error{Call: call, Message: strings.SplitN(err.Error(), "\n", 2)[0]}
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, &forge.Error{Call: call, Status: resp.StatusCode, Message: "the body cannot be read: " + err.Error()}
		}
		switch {
		case resp.StatusCode == want:
			if err := json.Unmarshal(data, out); err != nil {
				return nil, &forge.Error{Call: call, Status: resp.StatusCode, Message: "the body cannot be read: " + err.Error()}
			}
			return resp.Header, nil
		case resp.StatusCode == 401 && by == byToken && !retried:
			retried = true
			a.tok = forge.Token{}
			continue
		case resp.StatusCode == 403 || resp.StatusCode == 429:
			if until, ok := resetTime(resp.Header, a.cfg.Now()); ok {
				if err := a.wait(ctx, call, until); err != nil {
					return nil, &forge.Error{Call: call, Status: resp.StatusCode, Message: "the wait for the rate limit stopped: " + err.Error()}
				}
				continue
			}
		}
		return nil, &forge.Error{Call: call, Status: resp.StatusCode, Message: firstLine(data)}
	}
}

// resetTime gives the end of a rate limit: retry-after, or x-ratelimit-reset
// when x-ratelimit-remaining is 0 (forge.md, Forge errors). Any other 403 or
// 429 is not a rate limit.
func resetTime(h http.Header, now time.Time) (time.Time, bool) {
	if s, err := strconv.Atoi(h.Get("Retry-After")); err == nil && s >= 0 {
		return now.Add(time.Duration(s) * time.Second), true
	}
	if h.Get("X-Ratelimit-Remaining") == "0" {
		if epoch, err := strconv.ParseInt(h.Get("X-Ratelimit-Reset"), 10, 64); err == nil {
			return time.Unix(epoch, 0), true
		}
	}
	return time.Time{}, false
}

// wait sleeps until until, with a progress line every ten seconds; it is not
// a stall.
func (a *Adapter) wait(ctx context.Context, call string, until time.Time) error {
	for {
		left := until.Sub(a.cfg.Now())
		if left <= 0 {
			return nil
		}
		a.cfg.Progress(fmt.Sprintf("%s: the rate limit of the forge; waiting until %s, %s left",
			call, until.UTC().Format(time.RFC3339), left.Round(time.Second)))
		if err := a.cfg.Sleep(ctx, min(left, 10*time.Second)); err != nil {
			return err
		}
	}
}

// makeJWT gives the JWT of the App (RS256): issued 60 seconds before now, to
// absorb a clock skew, valid for 9 minutes from then, with the App ID as issuer.
func makeJWT(key *rsa.PrivateKey, appID int64, now time.Time) (string, error) {
	enc := base64.RawURLEncoding
	iat := now.Unix() - 60
	header := enc.EncodeToString([]byte(`{"alg":"RS256","typ":"JWT"}`))
	claims, err := json.Marshal(map[string]any{"iat": iat, "exp": iat + 9*60, "iss": strconv.FormatInt(appID, 10)})
	if err != nil {
		return "", err
	}
	signed := header + "." + enc.EncodeToString(claims)
	sum := sha256.Sum256([]byte(signed))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, sum[:])
	if err != nil {
		return "", err
	}
	return signed + "." + enc.EncodeToString(sig), nil
}

// firstLine gives the first line of the message of an error body: its JSON
// field message, else its first line.
func firstLine(body []byte) string {
	var e struct {
		Message string `json:"message"`
	}
	text := string(body)
	if json.Unmarshal(body, &e) == nil && e.Message != "" {
		text = e.Message
	}
	line := strings.TrimSpace(strings.SplitN(text, "\n", 2)[0])
	if line == "" {
		return "(no message)"
	}
	return line
}
