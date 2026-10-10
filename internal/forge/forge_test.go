package forge

import (
	"context"
	"reflect"
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
func (standIn) Comments(context.Context, int) ([]Comment, error)    { return nil, nil }
func (standIn) Comment(context.Context, int, string) (int64, error) { return 0, nil }

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
