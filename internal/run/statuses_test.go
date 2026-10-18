package run

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/rules"
)

const (
	headSHA = "4444444444444444444444444444444444444444"
	botID   = 9001
)

// statusWorld is a Sessions on the stand-in store with the stand-in forge and the App's bot.
func statusWorld(t *testing.T) (*Sessions, *forge.Fake) {
	t.Helper()
	f := forge.NewFake()
	f.AppBot = forge.Status{CreatorID: botID, CreatorLogin: "layup-agent[bot]"}
	s := newSessions(&standinStore{holder: "0123456789abcdef", files: map[string][]byte{}})
	s.Forge, s.Bot = f, botID
	return s, f
}

func postedOn(t *testing.T, f *forge.Fake, sha string) []forge.Status {
	t.Helper()
	h, err := f.Statuses(context.Background(), sha)
	if err != nil {
		t.Fatal(err)
	}
	return h.Statuses
}

func TestADescriptionIsCutAtOneHundredFortyCharacters(t *testing.T) {
	for _, c := range []struct{ name, in, want string }{
		{"140 characters stay", strings.Repeat("a", 140), strings.Repeat("a", 140)},
		{"141 are cut to 139 and an ellipsis", strings.Repeat("a", 141), strings.Repeat("a", 139) + "…"},
		{"multi-byte characters are counted as characters", strings.Repeat("é", 141), strings.Repeat("é", 139) + "…"},
		{"140 multi-byte characters stay", strings.Repeat("é", 140), strings.Repeat("é", 140)},
		{"a CR or LF becomes a space", "a\r\nb\nc", "a  b c"},
		{"a short text stays", "pass", "pass"},
	} {
		if got := fit(c.in); got != c.want {
			t.Errorf("%s: fit = %q, want %q", c.name, got, c.want)
		}
	}
}

func TestOnlyTheNewestStatusOfTheAppCountsAndOtherSourcesNeverDo(t *testing.T) {
	app := func(state, d string) forge.Status {
		return forge.Status{Context: GatesContext, State: state, Description: d, CreatorID: botID}
	}
	other := forge.Status{Context: GatesContext, State: "success", Description: "pass", CreatorID: 7}
	checks := forge.HeadChecks{Statuses: []forge.Status{other, app("failure", "new"), {Context: SpecContext, State: "success", CreatorID: botID}, app("success", "old")}}
	got, ok := appStatus(checks, GatesContext, botID)
	if !ok || got.State != "failure" || got.Description != "new" {
		t.Errorf("appStatus = %+v, %v; want the newest status of the App (failure, new), not the other source's and not the oldest", got, ok)
	}
	if _, ok := appStatus(forge.HeadChecks{Statuses: []forge.Status{other}}, GatesContext, botID); ok {
		t.Error("a status of another source counted as the App's")
	}
	if got, ok := appStatus(forge.HeadChecks{Statuses: []forge.Status{app("success", "x")}}, GatesContext, botID); !ok || got.Description != "x" {
		t.Errorf("the control, the App's own status: %+v, %v", got, ok)
	}
	if _, ok := appStatus(checks, VerifyContext, botID); ok {
		t.Error("a status of another context counted")
	}
}

func TestAStatusIsPostedOnceForAnUnchangedStateAndAgainWhenItChanges(t *testing.T) {
	s, f := statusWorld(t)
	ctx := context.Background()
	st := Status{GatesContext, "success", "pass"}
	if ok, err := s.PostStatus(ctx, headSHA, st); err != nil || !ok {
		t.Fatalf("the first post: %v, %v; want posted", ok, err)
	}
	if ok, err := s.PostStatus(ctx, headSHA, st); err != nil || ok {
		t.Fatalf("the same state and description: %v, %v; want not posted", ok, err)
	}
	if n := len(postedOn(t, f, headSHA)); n != 1 {
		t.Fatalf("%d statuses after two calls; want 1", n)
	}
	for _, c := range []Status{{GatesContext, "success", "clear: no product path"}, {GatesContext, "failure", "clear: no product path"}, {GatesContext, "success", "pass"}} {
		if ok, err := s.PostStatus(ctx, headSHA, c); err != nil || !ok {
			t.Errorf("%+v after a different one: %v, %v; want posted again (the newest, not the oldest, is compared)", c, ok, err)
		}
	}
	if ok, err := s.PostStatus(ctx, headSHA, Status{SpecContext, "success", "pass"}); err != nil || !ok {
		t.Errorf("another context: %v, %v; want posted", ok, err)
	}
}

func TestAStatusOfAnotherSourceDoesNotCountAsPosted(t *testing.T) {
	s, f := statusWorld(t)
	f.AddStatus(headSHA, forge.Status{Context: GatesContext, State: "success", Description: "pass", CreatorID: 7, CreatorLogin: "mallory"})
	if ok, err := s.PostStatus(context.Background(), headSHA, Status{GatesContext, "success", "pass"}); err != nil || !ok {
		t.Errorf("a status of another creator: %v, %v; want the App's posted", ok, err)
	}
	if ok, err := s.PostStatus(context.Background(), headSHA, Status{GatesContext, "success", "pass"}); err != nil || ok {
		t.Errorf("the control, the App's own: %v, %v; want not posted again", ok, err)
	}
}

func TestADescriptionOverTheBoundIsFitAndPostedOnce(t *testing.T) {
	s, f := statusWorld(t)
	long := Status{GatesContext, "failure", strings.Repeat("x", 300)}
	for i, want := range []bool{true, false} {
		if ok, err := s.PostStatus(context.Background(), headSHA, long); err != nil || ok != want {
			t.Errorf("call %d: %v, %v; want posted=%v and no refusal", i+1, ok, err, want)
		}
	}
	if got := postedOn(t, f, headSHA); len(got) != 1 || len([]rune(got[0].Description)) != forge.MaxDescription {
		t.Errorf("posted %+v; want one status of 140 characters", got)
	}
}

func TestAStatusIsRefusedBeforeAnyRequestWhenItsInputIsBad(t *testing.T) {
	s, f := statusWorld(t)
	good := Status{GatesContext, "success", "pass"}
	for name, c := range map[string]struct {
		sha string
		st  Status
	}{
		"a short SHA": {"abc", good}, "an upper-case SHA": {strings.Repeat("A", 40), good},
		"a context of no status of ours": {headSHA, Status{"ci/other", "success", "pass"}},
		"a state of none of the four":    {headSHA, Status{GatesContext, "ok", "pass"}},
	} {
		f.Fail = "Statuses" // a request would be an error of the forge; the refusal must come first and differ
		if ok, err := s.PostStatus(context.Background(), c.sha, c.st); err == nil || ok || strings.Contains(err.Error(), "Bad Gateway") {
			t.Errorf("%s: %v, %v; want a refusal before the request", name, ok, err)
		}
	}
	f.Fail = ""
	s.Bot = 0
	if ok, err := s.PostStatus(context.Background(), headSHA, good); err == nil || ok {
		t.Errorf("an unknown bot ID: %v, %v; want an error, never a match of any creator", ok, err)
	}
	s.Bot = botID
	if ok, err := s.PostStatus(context.Background(), headSHA, good); err != nil || !ok {
		t.Errorf("the control, the same call with the bot: %v, %v", ok, err)
	}
}

func TestALostLeasePostsNoStatus(t *testing.T) {
	s, f := statusWorld(t)
	s.Store.(*standinStore).holder = "ffffffffffffffff"
	ok, err := s.PostStatus(context.Background(), headSHA, Status{GatesContext, "success", "pass"})
	var lost *LostError
	if !errors.As(err, &lost) || ok || len(postedOn(t, f, headSHA)) != 0 {
		t.Errorf("a lost lease: %v, %v, %d statuses; want a LostError and nothing posted", ok, err, len(postedOn(t, f, headSHA)))
	}
	s.Store.(*standinStore).holder = s.RunID
	if ok, err := s.PostStatus(context.Background(), headSHA, Status{GatesContext, "success", "pass"}); err != nil || !ok {
		t.Errorf("the control, the lease held: %v, %v", ok, err)
	}
}

func TestAFailedReadOrPostIsAnErrorAndAnotherCallPostsAfterIt(t *testing.T) {
	s, f := statusWorld(t)
	st := Status{GatesContext, "success", "pass"}
	for _, call := range []string{"Statuses", "SetStatus"} {
		f.Fail = call
		if ok, err := s.PostStatus(context.Background(), headSHA, st); err == nil || ok {
			t.Errorf("a failing %s: %v, %v; want an error", call, ok, err)
		}
	}
	if n := len(postedOn(t, f, headSHA)); n != 0 {
		t.Fatalf("%d statuses after two failures; want none", n)
	}
	f.Fail = ""
	if ok, err := s.PostStatus(context.Background(), headSHA, st); err != nil || !ok {
		t.Errorf("the call after the failures: %v, %v; want posted", ok, err)
	}
}

func TestTheRulesStatusIsTheResultOfTheBasicRule(t *testing.T) {
	for _, c := range []struct {
		r    rules.Result
		want Status
	}{
		{rules.Result{State: rules.Success, Reason: rules.ReasonNoRulePath}, Status{RulesContext, "success", "no rule path"}},
		{rules.Result{State: rules.Success, Reason: rules.ReasonRenderedRecord}, Status{RulesContext, "success", "rendered record"}},
		{rules.Result{State: rules.Failure, Reason: rules.ReasonRulePath, Path: "rule-paths.tsv"}, Status{RulesContext, "failure", "rule path: rule-paths.tsv"}},
		{rules.Result{State: rules.Failure, Reason: rules.ReasonNotEqual, Path: "docs/prd/PRD-0001.md"}, Status{RulesContext, "failure", "not equal: docs/prd/PRD-0001.md"}},
		{rules.Result{State: rules.Failure, Reason: rules.ReasonNoRendering}, Status{RulesContext, "failure", "no rendering"}},
		{rules.Result{State: rules.Error, Reason: rules.ReasonRegister}, Status{RulesContext, "error", "register"}},
		{rules.Result{State: rules.Error, Reason: rules.ReasonGit}, Status{RulesContext, "error", "git"}},
		{rules.Result{}, Status{RulesContext, "error", "error: no result"}},
	} {
		if got := RulesStatus(c.r); got != c.want {
			t.Errorf("RulesStatus(%+v) = %+v, want %+v", c.r, got, c.want)
		}
	}
}

// tClock is a clock for the constructor of the Sessions of a run.
type tClock struct{}

func (tClock) Now() time.Time                             { return sessionAt }
func (tClock) Sleep(context.Context, time.Duration) error { return nil }

func TestTheSessionsOfTheRunKnowTheBotOfTheApp(t *testing.T) {
	cfg := startConfig(forge.NewFake())
	cfg.Clock = tClock{}
	st := &state{cfg: cfg, botID: botID}
	if got := st.sessions().Bot; got != botID {
		t.Errorf("state.sessions().Bot = %d, want %d (the ID that step forge read)", got, botID)
	}
}

func TestTheRegisterCommitOfAHeadIsTheRecordsOfItsBoundSessionOrTheCallers(t *testing.T) {
	rec1, rec2 := strings.Repeat("5", 40), strings.Repeat("6", 40)
	row := func(id, task, records string) []string {
		r := sessionRow(id, task, "implement", "claude")
		r[sessionRecords] = records
		return r
	}
	sessions := [][]string{row("S-d0000001", "T-ab12", rec1), row("S-d0000002", "T-ab12", rec2), row("S-d0000003", "T-zz99", rec2)}
	bound := func(session, head string) []string { return boundEvent(1, session, head) }
	for _, c := range []struct {
		name   string
		events [][]string
		given  string
		want   string
		ok     bool
	}{
		{"the session bound to the head", [][]string{bound("S-d0000001", headA)}, "", rec1, true},
		{"the other session of the task", [][]string{bound("S-d0000002", headA)}, "", rec2, true},
		{"a head code produced, the caller's commit", nil, rec2, rec2, true},
		{"two sessions of one head that agree", [][]string{bound("S-d0000001", headA), bound("S-d0000001", headA)}, "", rec1, true},
		{"two sessions of one head that disagree", [][]string{bound("S-d0000001", headA), bound("S-d0000002", headA)}, "", "", false},
		{"no session and no commit", [][]string{bound("S-d0000001", headB)}, "", "", false},
		{"a session of another task", [][]string{bound("S-d0000003", headA)}, "", "", false},
		{"a commit of the caller for a head that has a bound session", [][]string{bound("S-d0000001", headA)}, rec2, "", false},
		{"a commit that is no full SHA", nil, "main", "", false},
		{"a bound event of a session with no start row", [][]string{bound("S-d0000009", headA)}, "", "", false},
	} {
		got, err := registerRevOf("T-ab12", headA, c.events, sessions, c.given)
		if (err == nil) != c.ok || got != c.want {
			t.Errorf("%s: %q, %v; want %q, ok=%v", c.name, got, err, c.want, c.ok)
		}
	}
}
