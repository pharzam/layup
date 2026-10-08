package run

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/forge"
)

func TestStartAndRestartCallStepAtTheStartOfEachStep(t *testing.T) {
	var calls []string
	step := func(i, n int, name string) { calls = append(calls, fmt.Sprintf("%d/%d %s", i, n, name)) }
	f := goodForge()
	f.repo.Visibility = "private" // the plan check fails at step 2
	cfg := startConfig(f)
	cfg.Step = step
	Start(context.Background(), cfg)
	if want := []string{"1/9 forge", "2/9 plan"}; !reflect.DeepEqual(calls, want) {
		t.Errorf("Start called Step %q; want %q", calls, want)
	}
	calls = nil
	f = goodForge()
	cfg = startConfig(f)
	cfg.Step = step
	Restart(context.Background(), cfg) // no layup-records: forge fails
	if want := []string{"1/5 forge"}; !reflect.DeepEqual(calls, want) {
		t.Errorf("Restart called Step %q; want %q", calls, want)
	}
	cfg.Step = nil // a nil Step is no call
	Restart(context.Background(), cfg)
	if got := stepNames(9); !reflect.DeepEqual(got, []string{"forge", "plan", "baseline", "root-push", "read-back", "records", "issues", "watch", "lease"}) {
		t.Errorf("the steps of Start %q; want those of run.md", got)
	}
	if got := stepNames(5); !reflect.DeepEqual(got, []string{"forge", "clone", "version", "lease", "phase"}) {
		t.Errorf("the steps of the restart %q; want those of run.md", got)
	}
}

func TestCheckGit(t *testing.T) {
	saved := gitVersion
	defer func() { gitVersion = saved }()
	for _, c := range []struct {
		version string
		err     error
		want    string
	}{
		{"", errors.New("git not found: exec: \"git\": executable file not found in $PATH"), "not found"},
		{"2.31.8", nil, "2.32.0 or newer"},
		{"2.32.0", nil, ""},
		{"2.54.0", nil, ""},
	} {
		gitVersion = func() (string, error) { return c.version, c.err }
		err := CheckGit()
		if (err == nil) != (c.want == "") || (err != nil && !strings.Contains(err.Error(), c.want)) {
			t.Errorf("git %q, %v: CheckGit = %v; want %q", c.version, c.err, err, c.want)
		}
	}
}

func TestCheckValues(t *testing.T) {
	good := []string{"free", "50.0,8.0", "5", "10"}
	if err := CheckValues(good[0], good[1], good[2], good[3]); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		at   int
		bad  string
		want string
	}{
		{"a plan of no list", 0, "gold", "--plan"},
		{"an intake cap of one value", 1, "50.0", "--intake-cap"},
		{"an intake cap of integers", 1, "50,8", "--intake-cap"},
		{"lease.H of 0", 2, "0", "--lease-h"},
		{"lease.H that is not a whole number", 2, "5.5", "--lease-h"},
		{"watch.T of 0", 3, "0", "--watch-t"},
	} {
		v := append([]string{}, good...)
		v[c.at] = c.bad
		if err := CheckValues(v[0], v[1], v[2], v[3]); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: CheckValues = %v; want an error that names %s", c.name, err, c.want)
		}
	}
	_ = forge.Capabilities
}
