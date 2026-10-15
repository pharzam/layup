package forge

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

// standIn is a stand-in adapter of a test: it declares the capabilities it is
// given and makes no call.
type standIn struct{ declared []Capability }

func (s standIn) Capabilities() []Capability { return s.declared }
func (standIn) Installation(context.Context) (Installation, error) {
	return Installation{}, nil
}
func (standIn) Token(context.Context) (Token, error)           { return Token{}, nil }
func (standIn) Repository(context.Context) (Repository, error) { return Repository{}, nil }
func (standIn) UserID(context.Context, string) (int64, error)  { return 0, nil }
func (standIn) OpenIssue(context.Context, string, string) (int, error) {
	return 0, nil
}
func (standIn) Comments(context.Context, int) ([]Comment, error)          { return nil, nil }
func (standIn) Comment(context.Context, int, string) (int64, error)       { return 0, nil }
func (standIn) EffectiveRules(context.Context, string) ([]Rule, error)    { return nil, nil }
func (standIn) Activity(context.Context, string) ([]ActivityEntry, error) { return nil, nil }

func (standIn) OpenDraft(context.Context, string, string, string, string) (PullRequest, error) {
	return PullRequest{}, nil
}
func (standIn) MarkReady(context.Context, string) (PullRequest, error) { return PullRequest{}, nil }
func (standIn) SetStatus(context.Context, string, string, string, string) error {
	return nil
}
func (standIn) Statuses(context.Context, string) (HeadChecks, error) { return HeadChecks{}, nil }
func (standIn) PullRequest(context.Context, PullRequestQuery) (PullRequest, bool, error) {
	return PullRequest{}, false, nil
}
func (standIn) Merge(context.Context, int, string) (string, error) { return "", nil }

var _ Forge = standIn{}

func TestMissingNamesEachUndeclaredCapability(t *testing.T) {
	if len(Capabilities) != 6 {
		t.Fatalf("Capabilities has %d, want the six of forge.md", len(Capabilities))
	}
	five := append([]Capability{}, Capabilities[:2]...)
	five = append(five, Capabilities[3:]...)
	if got := Missing(standIn{five}.Capabilities()); !reflect.DeepEqual(got, []Capability{Capabilities[2]}) {
		t.Errorf("a stand-in that declares five: Missing = %v, want [%s]", got, Capabilities[2])
	}
	if got := Missing(standIn{Capabilities}.Capabilities()); len(got) != 0 {
		t.Errorf("a stand-in that declares six: Missing = %v, want none", got)
	}
	if got := Missing(nil); !reflect.DeepEqual(got, Capabilities) {
		t.Errorf("no declaration: Missing = %v, want the six in order", got)
	}
}

func TestCheckPermissionsNamesEachMissingPermission(t *testing.T) {
	for _, c := range []struct {
		name string
		got  map[string]string
		want []string
	}{
		{"the full set", map[string]string{"contents": "write", "issues": "write", "metadata": "read"}, nil},
		{"write holds read", map[string]string{"contents": "write", "issues": "write", "metadata": "write"}, nil},
		{"more than M2a uses", map[string]string{"contents": "write", "issues": "write", "metadata": "read", "statuses": "write"}, nil},
		{"no issues", map[string]string{"contents": "write", "metadata": "read"}, []string{"issues: write"}},
		{"contents read, not write", map[string]string{"contents": "read", "issues": "write", "metadata": "read"}, []string{"contents: write"}},
		{"none", nil, []string{"contents: write", "issues: write", "metadata: read"}},
		{"an unknown level", map[string]string{"contents": "admin", "issues": "write", "metadata": "read"}, []string{"contents: write"}},
	} {
		if got := CheckPermissions(c.got); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: CheckPermissions = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestCheckLoopPermissionsNamesEachMissingPermission(t *testing.T) {
	for _, c := range []struct {
		name string
		got  map[string]string
		want []string
	}{
		{"the full set", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "read"}, nil},
		{"write holds read", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "write"}, nil},
		{"more than the loop uses", map[string]string{"pull_requests": "write", "statuses": "write", "checks": "read", "contents": "write"}, nil},
		{"no checks", map[string]string{"pull_requests": "write", "statuses": "write"}, []string{"checks: read"}},
		{"statuses read, not write", map[string]string{"pull_requests": "write", "statuses": "read", "checks": "read"}, []string{"statuses: write"}},
		{"pull requests read, not write", map[string]string{"pull_requests": "read", "statuses": "write", "checks": "read"}, []string{"pull_requests: write"}},
		{"none", nil, []string{"pull_requests: write", "statuses: write", "checks: read"}},
		{"an unknown level", map[string]string{"pull_requests": "admin", "statuses": "write", "checks": "read"}, []string{"pull_requests: write"}},
	} {
		if got := CheckLoopPermissions(c.got); !reflect.DeepEqual(got, c.want) {
			t.Errorf("%s: CheckLoopPermissions = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestTheLoopPermissionsLeaveTheM2aOnesAlone(t *testing.T) {
	if got := CheckPermissions(map[string]string{"contents": "write", "issues": "write", "metadata": "read"}); got != nil {
		t.Errorf("CheckPermissions with the M2a set = %v, want none: the loop's three are not M2a's", got)
	}
}

func TestValidStatusRefusesAWrongStateOrATooLongDescription(t *testing.T) {
	long := strings.Repeat("é", MaxDescription) // 140 characters, 280 bytes
	for _, c := range []struct {
		state, desc string
		ok          bool
	}{
		{"success", "pass", true},
		{"failure", "", true},
		{"pending", "waits for its turn", true},
		{"error", long, true},
		{"success", long + "x", false},
		{"ok", "pass", false},
		{"", "pass", false},
		{"Success", "pass", false},
	} {
		err := ValidStatus(c.state, c.desc)
		if (err == nil) != c.ok {
			t.Errorf("ValidStatus(%q, %d characters) = %v, ok = %v", c.state, len([]rune(c.desc)), err, c.ok)
		}
	}
	if len(StatusStates) != 4 {
		t.Errorf("StatusStates = %v, want the four of forge.md", StatusStates)
	}
}
