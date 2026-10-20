//go:build integration

package run

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// The pull requests of the stand-in forge, through the adapter (row 66g, task T-xvvd, #363): the draft, the list by branch, the
// read by number, the GraphQL mutation that leaves the draft state, and the merge, which makes a real merge commit in the bare
// repository and answers 409 for a head that is not the pull request's.
type ghPull struct {
	number       int
	node, branch string
	base, title  string
	body, head   string
	draft        bool
	merged       bool
	closed       bool // closed by a person or by a PATCH of the state (task T-89xf)
	mergeSHA     string
}

type pullRoutes struct {
	list      []*ghPull
	afterMark func() // runs when a pull request leaves the draft state, under the lock of the fake
	noDraft   bool   // a request to open a pull request that is not a draft arrived
	denyPatch int    // the status that an edit of a pull request gets instead of the edit (row 66q)
}

func bareGit(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_NOSYSTEM=1", "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func (p *ghPull) wire() map[string]any {
	state := "open"
	if p.closed {
		state = "closed"
	}
	w := map[string]any{"number": p.number, "node_id": p.node, "state": state, "draft": p.draft, "merged": p.merged, "merged_at": nil,
		"merge_commit_sha": "9999999999999999999999999999999999999999", "head": map[string]any{"sha": p.head}, "body": p.body}
	if p.merged {
		w["state"], w["merged_at"], w["merge_commit_sha"] = "closed", "2026-10-18T12:00:00Z", p.mergeSHA
	}
	return w
}

// move moves the head of a pull request, as a push to its branch would.
func (f *fakeGitHub) movePull(number int, sha string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pulls.list[number-1].head = sha
}

func (f *fakeGitHub) pullOf(number int) *ghPull {
	f.mu.Lock()
	defer f.mu.Unlock()
	if number < 1 || number > len(f.pulls.list) {
		f.t.Fatalf("the forge holds no pull request %d", number)
	}
	return f.pulls.list[number-1]
}

func (p *pullRoutes) serve(f *fakeGitHub, w http.ResponseWriter, r *http.Request) bool {
	const base = "/repos/acme/target/pulls"
	reply := func(code int, v any) {
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(v)
	}
	switch {
	case r.Method == "POST" && r.URL.Path == base:
		var in struct {
			Title, Head, Base, Body string
			Draft                   bool
		}
		json.NewDecoder(r.Body).Decode(&in)
		head, err := bareGit(f.bare, "rev-parse", "refs/heads/"+in.Head)
		if err != nil {
			reply(422, map[string]string{"message": "head does not exist"})
			return true
		}
		if !in.Draft {
			p.noDraft = true
		}
		pr := &ghPull{number: len(p.list) + 1, branch: in.Head, base: in.Base, title: in.Title, body: in.Body, head: head, draft: in.Draft}
		pr.node = "PR_" + strconv.Itoa(pr.number)
		p.list = append(p.list, pr)
		reply(201, pr.wire())
	case r.Method == "GET" && r.URL.Path == base:
		branch := strings.TrimPrefix(r.URL.Query().Get("head"), "acme:")
		out := []map[string]any{}
		for i := len(p.list) - 1; i >= 0; i-- {
			if p.list[i].branch == branch {
				out = append(out, p.list[i].wire())
			}
		}
		reply(200, out)
	case r.Method == "GET" && strings.HasPrefix(r.URL.Path, base+"/") && !strings.HasSuffix(r.URL.Path, "/merge"):
		n, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, base+"/"))
		if n < 1 || n > len(p.list) {
			reply(404, map[string]string{"message": "Not Found"})
			return true
		}
		reply(200, p.list[n-1].wire())
	case r.Method == "PATCH" && strings.HasPrefix(r.URL.Path, base+"/"):
		n, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, base+"/"))
		var in struct{ Body, State *string }
		json.NewDecoder(r.Body).Decode(&in)
		switch {
		case p.denyPatch != 0:
			reply(p.denyPatch, map[string]string{"message": "Resource not accessible by integration"})
		case n < 1 || n > len(p.list) || (in.Body == nil && in.State == nil):
			reply(404, map[string]string{"message": "Not Found"})
		default:
			if in.Body != nil {
				p.list[n-1].body = *in.Body
			}
			if in.State != nil && *in.State == "closed" {
				p.list[n-1].closed = true
			}
			reply(200, p.list[n-1].wire())
		}
	case r.Method == "POST" && r.URL.Path == "/graphql":
		var in struct {
			Variables struct{ ID string }
		}
		json.NewDecoder(r.Body).Decode(&in)
		for _, pr := range p.list {
			if pr.node == in.Variables.ID {
				pr.draft = false
				if p.afterMark != nil {
					p.afterMark()
				}
				reply(200, map[string]any{"data": map[string]any{"markPullRequestReadyForReview": map[string]any{"pullRequest": map[string]any{
					"id": pr.node, "number": pr.number, "isDraft": false, "state": "OPEN", "headRefOid": pr.head}}}})
				return true
			}
		}
		reply(200, map[string]any{"errors": []map[string]string{{"message": "Could not resolve to a node"}}})
	case r.Method == "PUT" && strings.HasPrefix(r.URL.Path, base+"/") && strings.HasSuffix(r.URL.Path, "/merge"):
		n, _ := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, base+"/"), "/merge"))
		var in struct {
			SHA         string `json:"sha"`
			MergeMethod string `json:"merge_method"`
		}
		json.NewDecoder(r.Body).Decode(&in)
		switch {
		case n < 1 || n > len(p.list):
			reply(404, map[string]string{"message": "Not Found"})
		case in.MergeMethod != "merge":
			reply(422, map[string]string{"message": "a plain merge only"})
		case p.list[n-1].draft || p.list[n-1].merged:
			reply(405, map[string]string{"message": "Pull Request is not mergeable"})
		case in.SHA != p.list[n-1].head:
			reply(409, map[string]string{"message": "Head branch was modified. Review and try the merge again."})
		default:
			pr := p.list[n-1]
			tree, e1 := bareGit(f.bare, "rev-parse", pr.head+"^{tree}")
			main, e2 := bareGit(f.bare, "rev-parse", "refs/heads/main")
			merge, e3 := bareGit(f.bare, "commit-tree", tree, "-p", main, "-p", pr.head, "-m", fmt.Sprintf("Merge pull request #%d", n))
			if _, e4 := bareGit(f.bare, "update-ref", "refs/heads/main", merge); e1 != nil || e2 != nil || e3 != nil || e4 != nil {
				reply(500, map[string]string{"message": "the merge commit could not be made"})
				return true
			}
			pr.merged, pr.mergeSHA = true, merge
			reply(200, map[string]any{"sha": merge, "merged": true, "message": "Pull Request successfully merged"})
		}
	default:
		return false
	}
	return true
}
