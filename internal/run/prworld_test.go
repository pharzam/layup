package run

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

// The world of the pull request and the merge on the stand-in store (row 66g, task T-xvvd, #363): T-ab12 on issue 3,
// two attempts, each with a developer session that bound a head, and the stand-in forge that records what it was asked.
const (
	recA   = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	recB   = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	head1  = "1111111111111111111111111111111111111111"
	head2  = "2222222222222222222222222222222222222222"
	head3  = "3333333333333333333333333333333333333333"
	prWeb  = "https://github.example"
	prSess = "S-d0000001"
)

// A spyForge is the stand-in forge that logs each call of the pull request and the merge, keeps the title and the body that
// forge.Fake does not, and lets a test act between two reads.
type spyForge struct {
	*forge.Fake
	log        []string
	title      string
	body       string
	statusN    int
	onStatuses func(n int)        // runs before the n-th read of a head's statuses
	found      *forge.PullRequest // what PullRequest gives instead of the stand-in's
	mergeErr   error              // what Merge gives instead of the stand-in's
}

func (f *spyForge) OpenDraft(ctx context.Context, head, base, title, body string) (forge.PullRequest, error) {
	f.log = append(f.log, "OpenDraft "+head)
	f.title, f.body = title, body
	return f.Fake.OpenDraft(ctx, head, base, title, body)
}

func (f *spyForge) PullRequest(ctx context.Context, q forge.PullRequestQuery) (forge.PullRequest, bool, error) {
	f.log = append(f.log, "PullRequest")
	if f.found != nil {
		return *f.found, true, nil
	}
	return f.Fake.PullRequest(ctx, q)
}

func (f *spyForge) Statuses(ctx context.Context, sha string) (forge.HeadChecks, error) {
	f.statusN++
	f.log = append(f.log, "Statuses "+sha)
	if f.onStatuses != nil {
		f.onStatuses(f.statusN)
	}
	return f.Fake.Statuses(ctx, sha)
}

func (f *spyForge) MarkReady(ctx context.Context, nodeID string) (forge.PullRequest, error) {
	f.log = append(f.log, "MarkReady")
	return f.Fake.MarkReady(ctx, nodeID)
}

func (f *spyForge) Merge(ctx context.Context, number int, sha string) (string, error) {
	f.log = append(f.log, "Merge "+strconv.Itoa(number)+" "+sha)
	if f.mergeErr != nil {
		return "", f.mergeErr
	}
	return f.Fake.Merge(ctx, number, sha)
}

func (f *spyForge) count(prefix string) int {
	n := 0
	for _, l := range f.log {
		if len(l) >= len(prefix) && l[:len(prefix)] == prefix {
			n++
		}
	}
	return n
}

type prWorld struct {
	s      *Sessions
	st     *standinStore
	f      *spyForge
	events [][]string
}

func newPRWorld(t *testing.T) *prWorld {
	t.Helper()
	fake := forge.NewFake()
	fake.AppBot = forge.Status{CreatorID: botID, CreatorLogin: "layup-agent[bot]"}
	fake.SetHead("task/T-ab12/1", head1)
	fake.SetHead("task/T-ab12/2", head2)
	st := &standinStore{holder: "0123456789abcdef", files: map[string][]byte{}}
	w := &prWorld{s: newSessions(st), st: st, f: &spyForge{Fake: fake}}
	w.s.Forge, w.s.Bot, w.s.Web = w.f, botID, prWeb
	w.register(t, 0)
	var rows [][]string
	for i, a := range []struct{ id, records, head string }{{prSess, recA, head1}, {"S-d0000002", recB, head2}} {
		row := sessionRow(a.id, "T-ab12", "implement", "claude")
		row[sessionAttempt], row[sessionRecords] = strconv.Itoa(i+1), a.records
		rows = append(rows, row)
	}
	data, err := table(records.SessionsSchema, rows)
	if err != nil {
		t.Fatal(err)
	}
	st.files["sessions.tsv"] = data
	at := "2026-10-18T10:00:00Z"
	for i, a := range []struct{ id, head string }{{prSess, head1}, {"S-d0000002", head2}} {
		n := strconv.Itoa(i + 1)
		w.events = append(w.events,
			[]string{strconv.Itoa(len(w.events) + 1), "attempt", n, "", shaOld, "", "", at},
			[]string{strconv.Itoa(len(w.events) + 2), "session", n, a.id, "", "", "", at},
			[]string{strconv.Itoa(len(w.events) + 3), "bound", n, a.id, "", a.head, "task/T-ab12/" + n, at})
	}
	w.write(t)
	return w
}

func (w *prWorld) write(t *testing.T) {
	t.Helper()
	data, err := table(records.EventsSchema, w.events)
	if err != nil {
		t.Fatal(err)
	}
	w.st.files[eventsPath("T-ab12")] = data
}

// register writes tasks.tsv with T-ab12 on the pull request pr.
func (w *prWorld) register(t *testing.T, pr int) {
	t.Helper()
	var b bytes.Buffer
	err := WriteRegister(&b, []Task{{ID: "T-ab12", Issue: 3, Kind: "build", Requirements: []string{"REQ-005"}, Size: "small", PR: pr, Opened: time.Date(2026, 10, 18, 9, 0, 0, 0, time.UTC)}})
	if err != nil {
		t.Fatal(err)
	}
	w.st.files[tasksPath] = b.Bytes()
}

// open opens the draft of attempt 2 in the stand-in directly and registers its number: a pull request with the head head2.
func (w *prWorld) open(t *testing.T) forge.PullRequest {
	t.Helper()
	pr, err := w.f.Fake.OpenDraft(context.Background(), "task/T-ab12/2", "main", "feat: x", "b")
	if err != nil {
		t.Fatal(err)
	}
	w.register(t, pr.Number)
	return pr
}

func (w *prWorld) prOf(t *testing.T) int {
	t.Helper()
	tasks, err := ReadRegister(w.st.files[tasksPath])
	if err != nil {
		t.Fatal(err)
	}
	return tasks[0].PR
}

var actionsJob = forge.Rule{Type: "required_status_checks", Parameters: map[string]any{"required_status_checks": []any{
	map[string]any{"context": "build", "integration_id": float64(15368)}}}}

// green gives the head the four successes of the App, the job build of GitHub Actions and the rule that requires it.
func (w *prWorld) green(sha string) {
	for _, c := range statusContexts {
		w.f.Fake.AddStatus(sha, forge.Status{Context: c, State: "success", CreatorID: botID})
	}
	w.f.Fake.SetCheckRuns(sha, []forge.CheckRun{{Name: "build", Status: "completed", Conclusion: "success", AppID: 15368}})
	w.f.Fake.Rules = map[string][]forge.Rule{"main": {actionsJob}}
}
