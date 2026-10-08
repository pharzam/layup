package run

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
	"github.com/pharzam/layup/internal/tsv"
)

// standInForge is a stand-in forge of a unit test: each answer is a field, and
// fail names the call that gives a forge error.
type standInForge struct {
	declared []forge.Capability
	perms    map[string]string
	repo     forge.Repository
	repoErr  error
	ids      map[string]int64
	fail     string
}

func (f *standInForge) err(call string) error {
	if f.fail == call {
		return &forge.Error{Call: call, Status: 502, Message: "Bad Gateway"}
	}
	return nil
}
func (f *standInForge) Capabilities() []forge.Capability { return f.declared }
func (f *standInForge) Installation(context.Context) (forge.Installation, error) {
	return forge.Installation{ID: 7, Permissions: f.perms}, f.err("Installation")
}
func (f *standInForge) Token(context.Context) (forge.Token, error) {
	return forge.Token{Token: "ghs_x"}, f.err("Token")
}
func (f *standInForge) Repository(context.Context) (forge.Repository, error) {
	if f.repoErr != nil {
		return forge.Repository{}, f.repoErr
	}
	return f.repo, f.err("Repository")
}
func (f *standInForge) UserID(_ context.Context, login string) (int64, error) {
	if err := f.err("UserID"); err != nil {
		return 0, err
	}
	id, ok := f.ids[login]
	if !ok {
		return 0, &forge.Error{Call: "UserID", Status: 404, Message: "Not Found"}
	}
	return id, nil
}
func (f *standInForge) OpenIssue(context.Context, string, string) (int, error) {
	return 0, f.err("OpenIssue")
}
func (f *standInForge) Comments(context.Context, int) ([]forge.Comment, error) {
	return nil, f.err("Comments")
}

func goodForge() *standInForge {
	return &standInForge{
		declared: forge.Capabilities,
		perms:    map[string]string{"metadata": "read", "issues": "write", "contents": "write", "statuses": "read"},
		repo:     forge.Repository{DefaultBranch: "main", Visibility: "public"},
		ids:      map[string]int64{"pharzam": 101, "carol": 303, "layup-agent[bot]": 9001},
	}
}

func startConfig(f forge.Forge) Config {
	return Config{Owner: "acme", Name: "target", Operator: "pharzam", IdeaOwner: "carol", Plan: "free",
		Register: forge.Register{Forge: "github", AppSlug: "layup-agent"}, Forge: f}
}

func TestTheForgeStepOfStart(t *testing.T) {
	ctx := context.Background()
	f := goodForge()
	r := &state{cfg: startConfig(f)}
	if s := r.forgeStep(ctx, true); s.Result != "done" {
		t.Fatalf("the forge step of a good forge: %+v", s)
	}
	if r.permissions != "contents:write issues:write metadata:read statuses:read" {
		t.Errorf("app.permissions %q; want the pairs sorted by name", r.permissions)
	}
	if r.operatorID != 101 || r.ideaOwnerID != 303 || r.botID != 9001 || r.visibility != "public" {
		t.Errorf("the IDs %d %d %d, visibility %q", r.operatorID, r.ideaOwnerID, r.botID, r.visibility)
	}
	for _, c := range []struct {
		name   string
		change func(*standInForge)
		want   string
	}{
		{"a capability that the adapter does not declare", func(f *standInForge) { f.declared = forge.Capabilities[1:] }, string(forge.Capabilities[0])},
		{"a permission of M2a that is missing", func(f *standInForge) { delete(f.perms, "issues") }, "issues: write"},
		{"a repository that does not exist", func(f *standInForge) {
			f.repoErr = &forge.Error{Call: "Repository", Status: 404, Message: "Not Found"}
		}, "404"},
		{"--new on a repository with a commit (O-163)", func(f *standInForge) {
			f.repo.Branches, f.repo.HasCommit = []string{"main"}, true
		}, "commit"},
		{"--new on a repository with layup-records (O-163)", func(f *standInForge) {
			f.repo.Branches, f.repo.HasCommit = []string{"layup-records"}, true
		}, "layup-records"},
		{"an unknown login", func(f *standInForge) { delete(f.ids, "carol") }, "carol"},
		{"a forge error during the step", func(f *standInForge) { f.fail = "Installation" }, "502"},
	} {
		f := goodForge()
		c.change(f)
		s := (&state{cfg: startConfig(f)}).forgeStep(ctx, true)
		if s.Name != "forge" || s.Result != "fail" || !strings.Contains(s.Detail, c.want) {
			t.Errorf("%s: %+v; want forge: fail naming %q", c.name, s, c.want)
		}
	}
	// The restart: the repository must hold layup-records; the bot's ID is read.
	f = goodForge()
	f.repo.Branches, f.repo.HasCommit = []string{"main"}, true
	if s := (&state{cfg: startConfig(f)}).forgeStep(ctx, false); s.Result != "fail" || !strings.Contains(s.Detail, "layup-records") {
		t.Errorf("a restart on a repository with no layup-records: %+v", s)
	}
	f.repo.Branches = []string{"layup-records", "main"}
	delete(f.ids, "carol") // the restart looks up no approver login
	r = &state{cfg: startConfig(f)}
	if s := r.forgeStep(ctx, false); s.Result != "done" || r.botID != 9001 {
		t.Errorf("a restart: %+v, bot ID %d", s, r.botID)
	}
}

func TestThePlanCheck(t *testing.T) {
	for _, c := range []struct {
		plan, visibility string
		pass             bool
	}{
		{"free", "public", true}, {"pro", "public", true}, {"team", "public", true}, {"enterprise", "public", true},
		{"free", "private", false}, {"pro", "private", false}, {"team", "private", true}, {"enterprise", "private", true},
	} {
		r := &state{cfg: Config{Plan: c.plan}, visibility: c.visibility}
		s := r.planStep()
		if (s.Result == "done") != c.pass || s.Name != "plan" {
			t.Errorf("%s, %s: %+v; want pass %v", c.plan, c.visibility, s, c.pass)
		}
		if !c.pass && (!strings.Contains(s.Detail, c.plan) || !strings.Contains(s.Detail, c.visibility)) {
			t.Errorf("%s, %s: the detail %q does not name the plan and the visibility", c.plan, c.visibility, s.Detail)
		}
	}
}

func TestRunStepsIsATableOfItsBlock(t *testing.T) {
	var out strings.Builder
	if err := WriteSteps(&out, []Step{{"forge", "done", "the token"}, {"plan", "fail", "free, private"}}); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "step\tresult\tdetail\nforge\tdone\tthe token\nplan\tfail\tfree, private\n"; got != want {
		t.Errorf("run-steps:\n%q\nwant\n%q", got, want)
	}
	if _, err := tsv.Read([]byte("step\tresult\n"), RunStepsSchema); err == nil {
		t.Error("a table of run-steps with a wrong header was read")
	}
	if err := WriteSteps(&out, []Step{{"lunch", "done", ""}}); err == nil {
		t.Error("a step that is not one of the list was written")
	}
	var fe *forge.Error
	if !errors.As(error(&forge.Error{Call: "x"}), &fe) {
		t.Error("forge.Error")
	}
}
