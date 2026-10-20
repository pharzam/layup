//go:build integration

package run

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/pharzam/layup/internal/git"
)

// The findings of layup/spec and layup/rules start a new attempt, through the Sessions of state.sessions(), the httptest forge, real git and
// scripted harnesses (task T-tnxn, #427, row 66aa; loop.md, The attempts and The driver of a Shape task).

func (d *drvWorld) attempts(task string) []string {
	var out []string
	for _, e := range d.eventsOf(task) {
		if e[1] == "attempt" {
			out = append(out, e[4])
		}
	}
	return out
}

// headStates gives "<context> <state>" of every status of a head, newest first.
func headStates(d *drvWorld, sha string) []string {
	var out []string
	for _, s := range statusesOf(d.w, sha) {
		out = append(out, s["context"].(string)+" "+s["state"].(string))
	}
	return out
}

// seededRules rewrites the records commit of the developer's sessions of the attempt to a register that makes docs/tasks/ a rule path: a rule-path
// change cannot be pushed, so layup/rules cannot fail on a head of the loop's own making and this is the seeded state of the review of #415, finding 2
// (plan review condition 6 of #427: it is no evidence that a developer session can push a rule path, and the push refusal stays).
func (d *drvWorld) seededRules(task, attempt string, withRegister bool) func(*Sessions) {
	files := map[string]string{"rule-paths.tsv": register + "docs/tasks/\t—\tarchitecture\n"}
	if !withRegister {
		files = map[string]string{"README.md": "no register\n"}
	}
	rev := sideCommit(d.t, d.r.target, files)
	return func(s *Sessions) {
		s.Store = &registerAfterBound{SessionStore: s.Store, task: task, rev: rev, attempt: attempt}
	}
}

// I1, the demo: a head that fails layup/spec starts a new attempt at the same base, whose developer's prompt file holds the finding; the red
// control is the stop of #415, which this row replaces.
func TestASpecFailureDuringShapeStartsANewAttemptWithItsFindings(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("spec-fail", "")
	spec, _ := d.open()
	d.mustVerify(spec)
	bases := d.attempts(spec)
	if len(bases) != 2 || bases[0] != bases[1] {
		t.Fatalf("the attempts start at %v; want two attempts at one base", bases)
	}
	// The prompts are those of the plan, the plan review, the developer of attempt 1, the developer of attempt 2 and the verifier.
	prompts := d.promptsOf(spec)
	if len(prompts) != 5 || !strings.Contains(prompts[3], "spec:criteria") || strings.Contains(prompts[2], "spec:criteria") {
		t.Errorf("%d prompts; want the finding spec:criteria in the developer's prompt of attempt 2 and not in attempt 1's", len(prompts))
	}
	bound := d.boundOf(spec)
	if len(bound) != 2 || bound[0] == bound[1] {
		t.Fatalf("the bound heads %v; want two", bound)
	}
	if got := headStates(d, bound[0]); !slices.Contains(got, SpecContext+" failure") {
		t.Errorf("the statuses of the first head = %v; want a failure of layup/spec", got)
	}
	if got := headStates(d, bound[1]); !slices.Contains(got, SpecContext+" success") {
		t.Errorf("the statuses of the second head = %v; want a success of layup/spec", got)
	}
}

// I2 (plan review condition 6): a failed layup/rules starts a new attempt whose prompt file holds the exact finding with its full path.
func TestARulesFailureStartsANewAttemptWithItsPath(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	got, err := d.advance(spec, d.seededRules(spec, "1", true))
	if err != nil || !got.Verified {
		t.Fatalf("AdvanceShapeTask = %+v, %v; want the verified head after the second attempt", got, err)
	}
	if n := len(d.attempts(spec)); n != 2 {
		t.Errorf("%d attempts; want 2", n)
	}
	path := "docs/tasks/" + spec + ".md"
	prompts := d.promptsOf(spec)
	if len(prompts) != 5 || !strings.Contains(prompts[3], "- rules:rule path "+path+": rule path: "+path+"\n") || strings.Contains(prompts[2], "rules:") {
		t.Errorf("%d prompts; want the finding %q in the developer's prompt of attempt 2 and not in attempt 1's", len(prompts), "rules:rule path "+path)
	}
	if st := d.boundOf(spec); len(st) != 2 || !slices.Contains(headStates(d, st[0]), RulesContext+" failure") {
		t.Errorf("the statuses of the first head: want a failure of layup/rules")
	}
}

// I3: a restart between the post of the failed statuses and the start of the attempt starts exactly one attempt, posts no status twice at
// the failed head and context, and a run after that adds none (plan review note 2).
func TestARestartAfterTheFailedStatusesStartsExactlyOneAttempt(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	d.ctlWrite("spec-fail", "")
	spec, _ := d.open()
	refuse := func(s *Sessions) {
		s.Store = &faultStore{SessionStore: s.Store, fail: func(m string) error {
			if strings.Contains(m, "layup run: attempt 2 of "+spec) {
				return &LostError{Holder: "0000000000000000"}
			}
			return nil
		}}
	}
	var head string
	posts := -1
	for i := 0; i < 2; i++ {
		var lost *LostError
		if _, err := d.advance(spec, refuse); !errors.As(err, &lost) {
			t.Fatalf("call %d = %v; want the refusal of the attempt commit", i+1, err)
		}
		if n := len(d.attempts(spec)); n != 1 {
			t.Fatalf("call %d: %d attempts; want 1 after the refusal", i+1, n)
		}
		head = d.boundOf(spec)[0]
		if n := d.requests("POST /repos/acme/target/statuses/" + head); posts >= 0 && n != posts {
			t.Errorf("call %d: %d status posts at the failed head; want %d: no second post", i+1, n, posts)
		} else {
			posts = n
		}
	}
	d.mustVerify(spec)
	if n := len(d.attempts(spec)); n != 2 {
		t.Errorf("%d attempts after the restart; want exactly 2", n)
	}
	if n := d.requests("POST /repos/acme/target/statuses/" + head); n != posts {
		t.Errorf("%d status posts at the failed head after the restart; want %d", n, posts)
	}
	got := headStates(d, head)
	for _, ctx := range []string{GatesContext, SpecContext, RulesContext} {
		if n := len(slices.DeleteFunc(slices.Clone(got), func(s string) bool { return !strings.HasPrefix(s, ctx+" ") })); n != 1 {
			t.Errorf("the statuses of the failed head = %v; want one of %s", got, ctx)
		}
	}
	d.mustVerify(spec)
	if n := len(d.attempts(spec)); n != 2 {
		t.Errorf("%d attempts after a third run; want 2", n)
	}
}

// I4, a green negative (its red control is I1): the specification that waits for its bet, whose pinned prd-lint is accepted by code (an empty
// MoSCoW or Phase cell only), starts no attempt, and no layup/ status fails. The gate table of the head causes no retry.
func TestTheSpecificationThatWaitsForItsBetStartsNoAttempt(t *testing.T) {
	d := newDrvWorld(t, loopPerms)
	spec, _ := d.open()
	d.mustVerify(spec)
	if n := len(d.attempts(spec)); n != 1 {
		t.Errorf("%d attempts; want 1", n)
	}
	head := d.boundOf(spec)[0]
	for _, ctx := range []string{GatesContext, SpecContext, RulesContext} {
		if got := headStates(d, head); !slices.Contains(got, ctx+" success") || slices.Contains(got, ctx+" failure") {
			t.Errorf("the statuses of the head = %v; want %s success", got, ctx)
		}
	}
	base, err := git.RevParse(d.r.target, "refs/remotes/origin/main")
	if err != nil {
		t.Fatal(err)
	}
	if ok, lines, err := d.r.sessions().waitingSpecAccepted(base, head); err != nil || !ok || len(lines) != 0 {
		t.Errorf("waitingSpecAccepted = %v, %q, %v; want the head accepted by the pinned prd-lint with no line refused", ok, lines, err)
	}
}

// Plan review condition 3: a stop comes before the retry, whatever else failed. Each case fails layup/spec once, and then a status is an error.
func TestAStatusErrorBesideASpecFailureStopsAndStartsNoAttempt(t *testing.T) {
	for _, c := range []struct {
		name, want, gates string
		rules             bool
	}{
		{"layup/gates error (a green negative at the stub: the gates are judged first)", ": layup/gates is error: ", "kind\tstate\ttool\tcommand\tscope\tconfig\n", false},
		{"layup/rules error", ": layup/rules is error: ", gatesTSV, true},
	} {
		d := newDrvWorldWith(t, loopPerms, c.gates)
		d.ctlWrite("spec-fail", "")
		spec, _ := d.open()
		var tweak func(*Sessions)
		if c.rules {
			tweak = d.seededRules(spec, "1", false)
		}
		got, err := d.advance(spec, tweak)
		if err != nil || !strings.HasPrefix(got.Stop, spec) || !strings.Contains(got.Stop, c.want) {
			t.Errorf("%s: %+v, %v; want a stop holding %q", c.name, got, err, c.want)
		}
		if n := len(d.attempts(spec)); n != 1 {
			t.Errorf("%s: %d attempts; want one: a stop starts no attempt", c.name, n)
		}
	}
}
