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
	for next := a.cfg.API + a.repo() + "/branches?per_page=100"; next != ""; {
		var page []struct {
			Name string `json:"name"`
		}
		h, err := a.do(ctx, "Repository", "GET", next, byToken, nil, 200, &page)
		if err != nil {
			return forge.Repository{}, err
		}
		for _, b := range page {
			names = append(names, b.Name)
		}
		if next = nextPage(h.Get("Link")); next != "" && !strings.HasPrefix(next, a.cfg.API+"/") {
			return forge.Repository{}, &forge.Error{Call: "Repository", Status: 200, Message: "the next page is not at the api of the forge register"}
		}
	}
	return forge.Repository{DefaultBranch: repo.DefaultBranch, Visibility: repo.Visibility, HasCommit: len(names) > 0, Branches: names}, nil
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
