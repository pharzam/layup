package run

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/records"
)

func jobRule(entries ...any) forge.Rule {
	return forge.Rule{Type: "required_status_checks", Parameters: map[string]any{"required_status_checks": entries}}
}

func job(context, id any) map[string]any {
	return map[string]any{"context": context, "integration_id": id}
}

// The required jobs of GitHub Actions are read from the rules of the default branch, and a rule that is malformed is an error
// whatever the other entries are (the plan review of #363, condition 5).
func TestTheRequiredJobsAreTheActionsEntriesOfTheRequiredChecksRules(t *testing.T) {
	for _, c := range []struct {
		name  string
		rules []forge.Rule
		want  []string
	}{
		{"the control, one job", []forge.Rule{jobRule(job("build", float64(15368)))}, []string{"build"}},
		{"integer types", []forge.Rule{jobRule(job("a", 15368), job("b", int64(15368)), job("c", json.Number("15368")))}, []string{"a", "b", "c"}},
		{"a check of another app is no job of Actions", []forge.Rule{jobRule(job("build", float64(15368)), job("external", float64(42)))}, []string{"build"}},
		{"two rules and a duplicate", []forge.Rule{jobRule(job("a", 15368)), {Type: "deletion"}, jobRule(job("a", 15368), job("b", 15368))}, []string{"a", "b"}},
		{"no rule of the type", []forge.Rule{{Type: "deletion"}}, nil},
		{"no rule", nil, nil},
	} {
		got, err := requiredJobs(c.rules)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: requiredJobs = %v, %v; want %v", c.name, got, err, c.want)
		}
	}
	for _, c := range []struct {
		name string
		rule forge.Rule
	}{
		{"no parameters", forge.Rule{Type: "required_status_checks"}},
		{"no list", forge.Rule{Type: "required_status_checks", Parameters: map[string]any{"strict": true}}},
		{"a list that is no list", forge.Rule{Type: "required_status_checks", Parameters: map[string]any{"required_status_checks": "build"}}},
		{"an entry that is no object", jobRule("build")},
		{"no context", jobRule(map[string]any{"integration_id": 15368})},
		{"an empty context", jobRule(job("", 15368))},
		{"a context that is no string", jobRule(job(7, 15368))},
		{"no integration_id", jobRule(map[string]any{"context": "build"})},
		{"an integration_id that is a string", jobRule(job("build", "15368"))},
		{"an integration_id of 0", jobRule(job("build", 0))},
		{"a negative integration_id", jobRule(job("build", -1))},
		{"a fractional integration_id", jobRule(job("build", 1.5))},
		{"a malformed entry beside a good one", jobRule(job("build", 15368), map[string]any{"context": "x"})},
	} {
		if got, err := requiredJobs([]forge.Rule{c.rule}); err == nil {
			t.Errorf("%s: requiredJobs = %v, no error", c.name, got)
		}
	}
}

// Each case changes the head of a world that is ready; the first blocking context or job is named.
func TestTheReadyRuleNamesTheFirstContextOrJobThatIsNotSuccess(t *testing.T) {
	app := func(c, state string) forge.Status { return forge.Status{Context: c, State: state, CreatorID: botID} }
	other := func(c, state string) forge.Status { return forge.Status{Context: c, State: state, CreatorID: 7} }
	runs := func(r ...forge.CheckRun) func(*prWorld) {
		return func(w *prWorld) { w.f.Fake.SetCheckRuns(head2, r) }
	}
	run := func(status, conclusion string, id int64) forge.CheckRun {
		return forge.CheckRun{Name: "build", Status: status, Conclusion: conclusion, AppID: id}
	}
	for _, c := range []struct {
		name string
		do   func(*prWorld)
		want string
	}{
		{"the control, four successes of the App and the job", func(*prWorld) {}, ""},
		{"an older failure of another source under the App's success", func(w *prWorld) {
			w.f.Fake = rebuilt(w, func(c string) bool { return c != SpecContext })
			w.f.Fake.AddStatus(head2, other(SpecContext, "failure"))
			w.f.Fake.AddStatus(head2, app(SpecContext, "success"))
		}, ""},
		{"no status of the gates", func(w *prWorld) { w.f.Fake = rebuilt(w, func(c string) bool { return c != GatesContext }) }, "layup/gates: no status"},
		{"a success of another source only", func(w *prWorld) {
			w.f.Fake = rebuilt(w, func(c string) bool { return c != SpecContext })
			w.f.Fake.AddStatus(head2, other(SpecContext, "success"))
		}, "layup/spec: no status"},
		{"an error of the App", func(w *prWorld) { w.f.Fake.AddStatus(head2, app(RulesContext, "error")) }, "layup/rules: error"},
		{"a pending of the App", func(w *prWorld) { w.f.Fake.AddStatus(head2, app(VerifyContext, "pending")) }, "layup/verify: pending"},
		{"a failure of the App newer than its success", func(w *prWorld) { w.f.Fake.AddStatus(head2, app(GatesContext, "failure")) }, "layup/gates: failure"},
		{"a failure of another source newer than the App's success", func(w *prWorld) { w.f.Fake.AddStatus(head2, other(SpecContext, "failure")) }, "layup/spec: a status of another source is newer"},
		{"a success of another source newer than the App's success", func(w *prWorld) { w.f.Fake.AddStatus(head2, other(SpecContext, "success")) }, "layup/spec: a status of another source is newer"},
		{"the first of two contexts that block is named", func(w *prWorld) {
			w.f.Fake.AddStatus(head2, app(VerifyContext, "failure"))
			w.f.Fake.AddStatus(head2, app(RulesContext, "failure"))
		}, "layup/rules: failure"},
		{"no check run", runs(), "job build: no check run"},
		{"a neutral job", runs(run("completed", "neutral", 15368)), "job build: completed neutral"},
		{"a skipped job", runs(run("completed", "skipped", 15368)), "job build: completed skipped"},
		{"a failure", runs(run("completed", "failure", 15368)), "job build: completed failure"},
		{"a job in progress", runs(run("in_progress", "", 15368)), "job build: in_progress"},
		{"a success of another app", runs(run("completed", "success", 42)), "job build: no check run"},
		{"a success of Actions beside a failure of Actions", runs(run("completed", "success", 15368), run("completed", "failure", 15368)), "job build: completed failure"},
	} {
		w := newPRWorld(t)
		w.green(head2)
		c.do(w)
		r, err := w.s.Ready(context.Background(), head2)
		if err != nil || r.Ready != (c.want == "") || r.Reason != c.want {
			t.Errorf("%s: Ready = %+v, %v; want reason %q", c.name, r, err, c.want)
		}
	}
}

// rebuilt gives the stand-in forge of w with the statuses of the head again, those of the contexts that keep says true.
func rebuilt(w *prWorld, keep func(string) bool) *forge.Fake {
	f := forge.NewFake()
	f.AppBot = w.f.Fake.AppBot
	f.Rules = w.f.Fake.Rules
	f.SetHead("task/T-ab12/2", head2)
	for _, c := range statusContexts {
		if keep(c) {
			f.AddStatus(head2, forge.Status{Context: c, State: "success", CreatorID: botID})
		}
	}
	f.SetCheckRuns(head2, []forge.CheckRun{{Name: "build", Status: "completed", Conclusion: "success", AppID: 15368}})
	w.f.Fake = f
	w.s.Forge = w.f
	return f
}

func TestARulesetWithNoActionsJobIsNotReadyAndAMalformedOneIsAnError(t *testing.T) {
	for name, rules := range map[string][]forge.Rule{
		"no rule":                         nil,
		"a rule that requires no check":   {{Type: "deletion"}},
		"a required check of another app": {jobRule(job("external", float64(42)))},
	} {
		w := newPRWorld(t)
		w.green(head2)
		w.f.Fake.Rules = map[string][]forge.Rule{"main": rules}
		r, err := w.s.Ready(context.Background(), head2)
		if err != nil || r.Ready || r.Reason != "no required job of GitHub Actions" {
			t.Errorf("%s: Ready = %+v, %v; want not ready, no required job of GitHub Actions", name, r, err)
		}
	}
	w := newPRWorld(t)
	w.green(head2)
	w.f.Fake.Rules = map[string][]forge.Rule{"main": {jobRule(job("", 15368))}}
	if r, err := w.s.Ready(context.Background(), head2); err == nil {
		t.Errorf("a malformed required check: Ready = %+v, no error", r)
	}
}

func TestTheReadRuleFailsClosedOnAnErrorOfTheForgeAndOnABadInput(t *testing.T) {
	ctx := context.Background()
	for _, call := range []string{"Statuses", "EffectiveRules"} {
		w := newPRWorld(t)
		w.green(head2)
		w.f.Fake.Fail = call
		if r, err := w.s.Ready(ctx, head2); err == nil || r.Ready {
			t.Errorf("%s fails: Ready = %+v, %v; want an error and not ready", call, r, err)
		}
	}
	w := newPRWorld(t)
	w.green(head2)
	if r, err := w.s.Ready(ctx, "main"); err == nil || len(w.f.log) != 0 {
		t.Errorf("a head that is no full SHA: %+v, %v with the calls %v; want an error before any call", r, err, w.f.log)
	}
	w.s.Bot = 0
	if r, err := w.s.Ready(ctx, head2); err == nil {
		t.Errorf("no ID of the App's bot: %+v, no error", r)
	}
}

func TestMarkReadyMarksOnlyAReadyDraftByTheRead(t *testing.T) {
	ctx := context.Background()
	w := newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	w.f.Fake.AddStatus(head2, forge.Status{Context: VerifyContext, State: "failure", CreatorID: botID})
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err != nil || r.Ready || w.f.count("MarkReady") != 0 {
		t.Fatalf("a failure: MarkReady = %+v, %v with %d marks; want not ready and no mark", r, err, w.f.count("MarkReady"))
	}
	if got, _, _ := w.f.Fake.PullRequest(ctx, forge.PullRequestQuery{Number: pr.Number}); !got.Draft {
		t.Fatal("the draft was marked although a context was a failure")
	}
	w.f.Fake.AddStatus(head2, forge.Status{Context: VerifyContext, State: "success", CreatorID: botID})
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err != nil || !r.Ready || w.f.count("MarkReady") != 1 {
		t.Fatalf("the four successes: MarkReady = %+v, %v with %d marks; want ready and one mark", r, err, w.f.count("MarkReady"))
	}
	if got, _, _ := w.f.Fake.PullRequest(ctx, forge.PullRequestQuery{Number: pr.Number}); got.Draft {
		t.Error("the draft is still a draft after the mark")
	}
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err != nil || !r.Ready || w.f.count("MarkReady") != 1 {
		t.Errorf("a pull request that is ready: MarkReady = %+v, %v with %d marks; want ready and no second mark", r, err, w.f.count("MarkReady"))
	}
}

func TestMarkReadyRefusesAMovedHeadAndAPullRequestThatIsNotOpenOrNotRegistered(t *testing.T) {
	ctx := context.Background()
	w := newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	w.f.Fake.MoveHead(pr.Number, head3)
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err != nil || r.Ready || r.Reason != "the head moved to "+head3 || w.f.count("MarkReady") != 0 {
		t.Errorf("a moved head: MarkReady = %+v, %v; want the head moved to %s and no mark", r, err, head3)
	}
	w = newPRWorld(t)
	w.green(head2)
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err == nil {
		t.Errorf("a task with no pull request in the register: %+v, no error", r)
	}
	w = newPRWorld(t)
	w.open(t)
	w.green(head2)
	closed := forge.PullRequest{Number: 1, State: "closed", Draft: true, HeadSHA: head2}
	w.f.found = &closed
	if r, err := w.s.MarkReady(ctx, "T-ab12", head2); err != nil || r.Ready || r.Reason != "the pull request is not open" || w.f.count("MarkReady") != 0 {
		t.Errorf("a closed pull request: %+v, %v; want not ready, the pull request is not open", r, err)
	}
}

// mergedRow reads the last event of the task.
func lastEvent(t *testing.T, w *prWorld) []string {
	t.Helper()
	rows := events(t, w.st, "T-ab12")
	return rows[len(rows)-1]
}

func TestTheMergeIsAtTheHeadItReadAndWritesTheMergedEvent(t *testing.T) {
	ctx := context.Background()
	w := newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	merge, r, err := w.s.MergeTask(ctx, "T-ab12", head2)
	if err != nil || !r.Ready || merge == "" {
		t.Fatalf("MergeTask = %q, %+v, %v; want a merge", merge, r, err)
	}
	if want := "Merge 1 " + head2; !slices.Contains(w.f.log, want) {
		t.Errorf("the calls %v; want %q: the merge at the head it read", w.f.log, want)
	}
	if got, _, _ := w.f.Fake.PullRequest(ctx, forge.PullRequestQuery{Number: pr.Number}); !got.Merged || got.MergeSHA != merge {
		t.Errorf("the pull request after the merge: %+v; want merged at %s", got, merge)
	}
	e := lastEvent(t, w)
	if e[1] != "merged" || e[2] != "2" || e[3] != "" || e[5] != merge || e[6] != "1" {
		t.Errorf("the last event %q; want merged, attempt 2 (the attempt of the bound head), no session, the merge commit and pull request 1", e)
	}
	// A second call finds the pull request merged and writes nothing more.
	n := len(events(t, w.st, "T-ab12"))
	if merge2, _, err := w.s.MergeTask(ctx, "T-ab12", head2); err != nil || merge2 != merge || len(events(t, w.st, "T-ab12")) != n || w.f.count("Merge ") != 1 {
		t.Errorf("the second call = %q, %v with %d events and %d merges; want the same merge, no event, no second merge", merge2, err, len(events(t, w.st, "T-ab12")), w.f.count("Merge "))
	}
}

func TestAHeadThatMovedBeforeTheMergeIsNotMerged(t *testing.T) {
	ctx := context.Background()
	w := newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	w.f.Fake.MoveHead(pr.Number, head3)
	merge, r, err := w.s.MergeTask(ctx, "T-ab12", head2)
	if err != nil || merge != "" || r.Ready || r.Reason != "the head moved to "+head3 {
		t.Fatalf("MergeTask = %q, %+v, %v; want the head moved to %s", merge, r, err, head3)
	}
	if w.f.count("Merge ") != 0 || w.f.count("MarkReady") != 0 || len(w.st.commits) != 0 {
		t.Errorf("the calls %v and %d commits; want no mutation and no event", w.f.log, len(w.st.commits))
	}
}

// The head moves between the read and the merge: the forge answers 409. The call reads the pull request, then the statuses of the new
// head, and sends no second merge from the old read (the plan review of #363, condition 6).
func TestAMergeRefusedForAMovedHeadReadsThePullRequestAndTheNewStatusesAndMergesNoMore(t *testing.T) {
	ctx := context.Background()
	w := newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	w.green(head3)
	// The last read before the merge is the second of the call (the draft is marked after the first); the head moves after it.
	w.f.onStatuses = func(n int) {
		if n == 2 {
			w.f.Fake.MoveHead(pr.Number, head3)
		}
	}
	merge, r, err := w.s.MergeTask(ctx, "T-ab12", head2)
	if err != nil || merge != "" || r.Ready || !strings.HasPrefix(r.Reason, "the head moved to "+head3) {
		t.Fatalf("MergeTask = %q, %+v, %v; want the head moved to %s", merge, r, err, head3)
	}
	i := slices.Index(w.f.log, "Merge 1 "+head2)
	if i < 0 || w.f.count("Merge ") != 1 {
		t.Fatalf("the calls %v; want exactly one merge, at %s", w.f.log, head2)
	}
	if after := w.f.log[i+1:]; !slices.Equal(after, []string{"PullRequest", "Statuses " + head3}) {
		t.Errorf("the calls after the refused merge: %v; want the pull request read, then the statuses of the new head, and nothing else", after)
	}
	if len(w.st.commits) != 0 {
		t.Error("an event was written for a merge that was refused")
	}
	// A 409 whose pull request still has the head that was sent is an unexplained refusal: an error, and no second merge.
	w = newPRWorld(t)
	w.open(t)
	w.green(head2)
	w.f.mergeErr = &forge.Error{Call: "Merge", Status: 409, Message: "Head branch was modified"}
	if _, r, err := w.s.MergeTask(ctx, "T-ab12", head2); err == nil || w.f.count("Merge ") != 1 {
		t.Errorf("a 409 with the same head: %+v, %v with %d merges; want an error and one merge", r, err, w.f.count("Merge "))
	}
}

func TestAFailedMergeWritesNoEvent(t *testing.T) {
	w := newPRWorld(t)
	w.open(t)
	w.green(head2)
	w.f.Fake.Fail = "Merge"
	merge, _, err := w.s.MergeTask(context.Background(), "T-ab12", head2)
	var fe *forge.Error
	if !errors.As(err, &fe) || merge != "" {
		t.Errorf("MergeTask = %q, %v; want a forge error", merge, err)
	}
	if len(w.st.commits) != 0 {
		t.Errorf("%d commits after a failed merge; want none", len(w.st.commits))
	}
}

func TestALostLeaseOpensMarksAndMergesNothing(t *testing.T) {
	ctx := context.Background()
	var lost *LostError
	w := newPRWorld(t)
	w.st.holder = "another-run"
	if _, err := w.s.OpenPullRequest(ctx, PullInput{Task: "T-ab12", Attempt: 1, Type: "feat", Description: "add the draft"}); !errors.As(err, &lost) || w.f.count("OpenDraft") != 0 {
		t.Errorf("OpenPullRequest: %v with %d drafts; want a *LostError and no draft", err, w.f.count("OpenDraft"))
	}
	w = newPRWorld(t)
	w.open(t)
	w.green(head2)
	w.st.holder = "another-run"
	if _, err := w.s.MarkReady(ctx, "T-ab12", head2); !errors.As(err, &lost) || w.f.count("MarkReady") != 0 {
		t.Errorf("MarkReady: %v with %d marks; want a *LostError and no mark", err, w.f.count("MarkReady"))
	}
	if _, _, err := w.s.MergeTask(ctx, "T-ab12", head2); !errors.As(err, &lost) || w.f.count("Merge ") != 0 {
		t.Errorf("MergeTask of a draft: %v with %d merges; want a *LostError and no merge", err, w.f.count("Merge "))
	}
	// A pull request that is already ready reaches the lease read just before the merge.
	w = newPRWorld(t)
	pr := w.open(t)
	w.green(head2)
	if _, err := w.f.Fake.MarkReady(ctx, pr.NodeID); err != nil {
		t.Fatal(err)
	}
	w.st.holder = "another-run"
	if _, _, err := w.s.MergeTask(ctx, "T-ab12", head2); !errors.As(err, &lost) || w.f.count("Merge ") != 0 {
		t.Errorf("MergeTask of a ready pull request: %v with %d merges; want a *LostError and no merge", err, w.f.count("Merge "))
	}
}

func TestAMergedPullRequestIsRecordedOnceAndTheAttemptIsNeverGuessed(t *testing.T) {
	ctx := context.Background()
	merged := func(n int) forge.PullRequest {
		return forge.PullRequest{Number: n, State: "closed", Merged: true, HeadSHA: head2, MergeSHA: head3}
	}
	w := newPRWorld(t)
	w.open(t)
	for i := 1; i <= 2; i++ {
		if err := w.s.RecordMerged(ctx, "T-ab12", merged(1)); err != nil {
			t.Fatalf("RecordMerged, call %d: %v", i, err)
		}
		n := 0
		for _, e := range events(t, w.st, "T-ab12") {
			if e[1] == "merged" {
				n++
			}
		}
		if n != 1 {
			t.Fatalf("%d events merged after %d calls; want 1", n, i)
		}
	}
	if e := lastEvent(t, w); e[2] != "2" || e[5] != head3 || e[6] != "1" {
		t.Errorf("the event %q; want attempt 2, the merge commit and pull request 1", e)
	}
	// No bound event of the head, a pull request with no merge commit and two attempts that bound the head are errors and write nothing.
	w = newPRWorld(t)
	w.open(t)
	other := merged(1)
	other.HeadSHA = head3
	nomerge := merged(1)
	nomerge.MergeSHA = ""
	w.events = append(w.events, []string{"7", "bound", "1", prSess, "", head2, "task/T-ab12/1", "2026-10-18T10:00:00Z"})
	w.write(t)
	for name, pr := range map[string]forge.PullRequest{"a head that no session bound": other, "no merge commit": nomerge, "a head bound by two attempts": merged(1)} {
		before := len(w.st.commits)
		if err := w.s.RecordMerged(ctx, "T-ab12", pr); err == nil || len(w.st.commits) != before {
			t.Errorf("%s: RecordMerged = %v with %d new commits; want an error and no event", name, err, len(w.st.commits)-before)
		}
	}
}

func TestTheEventMergedIsAcceptedByTheSchemaOfTheEvents(t *testing.T) {
	w := newPRWorld(t)
	w.open(t)
	w.green(head2)
	if _, _, err := w.s.MergeTask(context.Background(), "T-ab12", head2); err != nil {
		t.Fatal(err)
	}
	rows, err := records.ReadEvents(w.st.files[eventsPath("T-ab12")])
	if err != nil || rows[len(rows)-1][1] != "merged" {
		t.Errorf("events.tsv after the merge: %v, %v; want it to read back with the event merged last", rows, err)
	}
}
