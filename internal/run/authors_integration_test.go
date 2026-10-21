//go:build integration

package run

import (
	"context"
	"errors"
	"slices"
	"testing"

	"github.com/pharzam/layup/internal/records"
	"github.com/pharzam/layup/internal/route"
)

// lessonResult is a result whose lesson row names a harness: a field an agent
// writes, which makes nobody an author.
const lessonResult = resultHead + "lesson\t1\twritten by opencode, reviewed by codex\t—\t—\n"

// loop.md, The independence of the reviewers: the authors come from the
// records and the ancestry of the commits, never from a field an agent writes.
func TestTheAuthorsAreReadFromRealCommitsAndRealRecords(t *testing.T) {
	t.Parallel()
	l := newLoop(t, nil)
	ctx := context.Background()
	s := l.prod()
	l.plan(1, "codex", "the plan\n")
	one := l.dev(1, "claude", map[string]string{"one.txt": "one\n"}, lessonResult)
	l.attempt(2)
	two := l.dev(2, "devin", map[string]string{"two.txt": "two\n"}, okResult)
	head1, head2 := l.bound(one), l.bound(two)
	for _, c := range []struct {
		name string
		head string
		want []string
	}{
		{"the head of attempt 2", head2, []string{"devin"}},
		{"the head of attempt 1", head1, []string{"claude"}},
	} {
		got, err := s.Authors(ctx, "T-ab12", Verification, l.a.base, c.head)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: authors %v, %v; want %v (the commit of the other attempt is not in this diff; the lesson row names opencode and codex)", c.name, got, err, c.want)
		}
	}
	if got, err := s.Authors(ctx, "T-ab12", PlanReview, l.a.base, head2); err != nil || !slices.Equal(got, []string{"codex"}) {
		t.Errorf("the authors of the plan review: %v, %v; want [codex], the harness of the plan session", got, err)
	}
	// A commit that the base already holds is no author: the head is the base.
	if got, err := s.Authors(ctx, "T-ab12", Verification, head1, head1); err != nil || len(got) != 0 {
		t.Errorf("the authors when the base is the head of attempt 1: %v, %v; want none", got, err)
	}
	// A diff of no commit: the head is the base, and nobody is bound to a commit of it.
	if got, err := s.Authors(ctx, "T-ab12", Verification, l.a.base, l.a.base); err != nil || len(got) != 0 {
		t.Errorf("the authors of an empty diff: %v, %v; want none", got, err)
	}
}

// The writer of the frozen tests is an author of every later verification, and
// the harness of a later attempt is one by its commits.
func TestTheWriterOfTheFrozenTestsIsAnAuthorOfALaterAttempt(t *testing.T) {
	t.Parallel()
	l := newLoop(t, nil)
	l.frozenAt("claude")
	l.attempt(2)
	two := l.dev(2, "devin", map[string]string{"c.txt": "other\n"}, okResult)
	got, err := l.prod().Authors(context.Background(), "T-ab12", Verification, l.a.base, l.bound(two))
	if err != nil || !slices.Equal(got, []string{"claude", "devin"}) {
		t.Errorf("authors %v, %v; want [claude devin]: devin by its commits, claude by the frozen tests of tests.tsv", got, err)
	}
}

// routed commits a routing register of qa-engineer pairs and the probes of the
// harnesses, and gives the register rows that a run holds.
func (l *loopRun) routed(pairs ...string) ([][]string, [][]string) {
	l.t.Helper()
	routing := "role\ttier\tposition\tharness\tmodel\n"
	var probes, harnesses, models [][]string
	for i, h := range []string{"claude", "devin", "codex"} {
		probes = append(probes, []string{"S-b000000" + string(rune('1'+i)), h, "1.0", "m1", "passed", "—", "AGENTS.md", "m1", "2026-10-16T12:00:00Z"})
		harnesses = append(harnesses, []string{h, "", "1", "/bin/true", "stdin", "/bin/true", "", "", "AGENTS.md", "", "none", "subscription", ""})
		models = append(models, []string{h, "m1", "1000", "https://example.invalid", "2026-10-16T12:00:00Z", "yes", "—"})
	}
	for i, h := range pairs {
		routing += "qa-engineer\treasoning\t" + string(rune('1'+i)) + "\t" + h + "\tm1\n"
	}
	data, err := table(records.HarnessesSchema, probes)
	if err != nil {
		l.t.Fatal(err)
	}
	if err := l.r.commit(context.Background(), map[string][]byte{"harnesses.tsv": data, "routing.tsv": []byte(routing)}, "the routing"); err != nil {
		l.t.Fatal(err)
	}
	return harnesses, models
}

func TestTheReviewersPairComesFromTheRoutingOfTheRecordsOutsideTheAuthors(t *testing.T) {
	t.Parallel()
	l := newLoop(t, nil)
	ctx := context.Background()
	l.plan(1, "claude", "the plan\n")
	l.frozenAt("claude")
	l.attempt(2)
	two := l.dev(2, "devin", map[string]string{"c.txt": "other\n"}, okResult)
	head := l.bound(two)
	harnesses, models := l.routed("claude", "devin", "codex")
	s := l.prod()
	h, m, err := s.ReviewerPair(ctx, "T-ab12", Verification, l.a.base, head, harnesses, models)
	if err != nil || idOf(h) != "codex" || len(m) < 2 || m[1] != "m1" {
		t.Errorf("the verifier's pair = %v, %v, %v; want codex: claude wrote the tests, devin the commits", h, m, err)
	}
	h, _, err = s.ReviewerPair(ctx, "T-ab12", PlanReview, l.a.base, head, harnesses, models)
	if err != nil || idOf(h) != "devin" {
		t.Errorf("the plan reviewer's pair = %v, %v; want devin: claude wrote the plan", h, err)
	}
	// Every admitted harness is an author: not-active, and the refusal of a pair too.
	harnesses, models = l.routed("claude", "devin")
	_, _, err = s.ReviewerPair(ctx, "T-ab12", Verification, l.a.base, head, harnesses, models)
	if !errors.Is(err, ErrNoHarnessOutsideAuthors) || !errors.Is(err, route.ErrNoPair) || err.Error() == "" {
		t.Errorf("no harness outside the authors: %v; want ErrNoHarnessOutsideAuthors and route.ErrNoPair", err)
	}
	if ErrNoHarnessOutsideAuthors.Error() != "no harness outside the authors" {
		t.Errorf("the text of the error is %q; the status that posts it says `not-active: no harness outside the authors`", ErrNoHarnessOutsideAuthors.Error())
	}
}
