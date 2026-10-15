//go:build integration

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

const (
	headSHA  = "1111111111111111111111111111111111111111"
	mergeSHA = "3333333333333333333333333333333333333333"
	pulls    = "/repos/acme/target/pulls"
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

func prJSON(number int, draft, merged bool) map[string]any {
	return map[string]any{"number": number, "node_id": fmt.Sprintf("PR_node%d", number), "state": "open", "draft": draft,
		"merged": merged, "merge_commit_sha": nil, "head": map[string]any{"sha": headSHA, "ref": "task/T-aaaa/1"}}
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
		io.WriteString(w, `{"data":{"markPullRequestReadyForReview":{"pullRequest":{"id":"PR_node5","number":5,"isDraft":false}}}}`)
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
	if pr.Number != 5 || pr.Draft || pr.NodeID != "PR_node5" {
		t.Errorf("MarkReady = %+v", pr)
	}
}

func TestMarkReadyGivesAForgeErrorForAnErrorsArrayOrAnEmptyAnswer(t *testing.T) {
	for name, body := range map[string]string{
		"errors": `{"errors":[{"message":"Could not resolve to a node\nwith the global id"}],"data":null}`,
		"empty":  `{"data":{"markPullRequestReadyForReview":null}}`,
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
	handle(f, "GET "+pulls, func(w http.ResponseWriter, r *http.Request) {
		queries = append(queries, r.URL.RawQuery)
		if r.URL.Query().Get("head") == "acme:task/T-aaaa/1" {
			json.NewEncoder(w).Encode([]map[string]any{prJSON(6, true, false), prJSON(5, false, true)})
			return
		}
		io.WriteString(w, `[]`)
	})
	ctx := context.Background()
	pr, found, err := a.PullRequest(ctx, forge.PullRequestQuery{Number: 5})
	if err != nil || !found || pr != (forge.PullRequest{Number: 5, NodeID: "PR_node5", State: "closed", Merged: true, HeadSHA: headSHA, MergeSHA: mergeSHA}) {
		t.Errorf("by number = %+v, %v, %v", pr, found, err)
	}
	pr, found, err = a.PullRequest(ctx, forge.PullRequestQuery{Branch: "task/T-aaaa/1"})
	if err != nil || !found || pr.Number != 6 || !pr.Draft {
		t.Errorf("by branch = %+v, %v, %v, want the first of the list, the newest", pr, found, err)
	}
	if len(queries) != 1 || !strings.Contains(queries[0], "state=all") {
		t.Errorf("the list request = %v, want head and state=all", queries)
	}
	if _, found, err := a.PullRequest(ctx, forge.PullRequestQuery{Branch: "task/none"}); found || err != nil {
		t.Errorf("a branch with no pull request: found = %v, err = %v, want a result and no error", found, err)
	}
}

func TestPullRequestRefusesAQueryThatNamesNeitherOrBothAndGivesAnErrorForANumberItDoesNotKnow(t *testing.T) {
	_, _, a, _ := setup(t)
	ctx := context.Background()
	for _, q := range []forge.PullRequestQuery{{}, {Number: 5, Branch: "b"}} {
		if _, _, err := a.PullRequest(ctx, q); err == nil {
			t.Errorf("PullRequest(%+v) gave no error", q)
		}
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
	if _, err := a.Merge(context.Background(), 5, "main"); err == nil {
		t.Errorf("Merge at a SHA that is none gave no error")
	}
	if _, err := a.Merge(context.Background(), 0, headSHA); err == nil {
		t.Errorf("Merge of the pull request 0 gave no error")
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
