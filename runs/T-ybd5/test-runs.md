# T-ybd5 — test runs

The budget of milestone 0 (#216, row 42b).

## Red (05:51Z, 2026-10-11)

The tests before the code, against stubs (`MilestoneCap` gave nil, `CommitTime`
the zero time), so each fails at its assertion. The `pre-commit` hook runs the
unit level, so the red is recorded here and the tests land with the code
(`docs/guardrails.md`, "The hook refuses a red commit").

```
--- FAIL: TestMilestoneCap (0.00s)
    budget_test.go:64: above: <nil>; want ErrMilestoneCap with "the known money 2.00 plus the session's cap 3.01 passes 5.00 USD"
    budget_test.go:64: the known money alone above: <nil>; want ErrMilestoneCap with "the known money 2.00 plus the session's cap 0 passes 1.99 USD"
    budget_test.go:64: the hours and one second: <nil>; want ErrMilestoneCap with "8h0m1s since Start passes 8.0 hours"
    budget_test.go:64: the hours of a decimal: <nil>; want ErrMilestoneCap with "31m0s since Start passes 0.5 hours"
    budget_test.go:64: a known row not in USD: <nil>; want ErrMilestoneCap with "the session S-00000004 has money in EUR, not USD"
--- FAIL: TestEachCallRunsItsVerb (0.00s)
    --- FAIL: TestEachCallRunsItsVerb/log_-1_--format=%ct (0.00s)
        git_test.go:107: 0 starts of git, want 1
--- FAIL: TestTheOutputIsRead (0.00s)
    git_test.go:178: log: the committer time, seconds since the epoch: "0001-01-01 00:00:00 +0000 UTC", <nil>; want "2026-01-01 00:00:00 +0000 UTC"
--- FAIL: TestCommitTime (0.02s)
    git_integration_test.go:828: CommitTime(HEAD) = 0001-01-01 00:00:00 +0000 UTC, <nil>; want 2026-01-01 00:00:00 +0000 UTC
    git_integration_test.go:831: CommitTime of a rev that names no commit: no error
--- FAIL: TestTheCapOfMilestoneZeroStartsNoSession (5.74s)
    --- FAIL: TestTheCapOfMilestoneZeroStartsNoSession/the_money (1.97s)
        phase_integration_test.go:281: PhaseSession: "S-7e3df3b5", <nil>; want ErrMilestoneCap with "the known money 0 plus the session's cap 50.01 passes 50.0 USD" and no session
    --- FAIL: TestTheCapOfMilestoneZeroStartsNoSession/the_hours (1.97s)
        phase_integration_test.go:281: PhaseSession: "S-0bb39227", <nil>; want ErrMilestoneCap with "since Start passes 8.0 hours" and no session
```

Before the code, the expected sums were changed from `2.00` and `0` to `2.0`
and `0.0`, the form of the type `decimal` that the ledger writes (`decimal`);
the rule each case tests is unchanged.

After the code, the six tests of row 42a went red too, each with "the cap of
milestone 0: 65h50m… since Start passes 8.0 hours": the clock of the test world
starts on 2026-10-08. `phaseWorld` now runs its sessions a minute after the
world's Start, which is evidence that the rule is asked before each session.

## Mutations (2026-10-11, 05:54Z to 05:55Z)

Each rule mutated alone, the five tests above run at the unit and integration
levels: 12 of 12 caught, each mutant built. The mutant `head-not-root` was
caught by an assertion added before the run, that the time of Start is the
first commit's, not the head's.

```
== money-equal-out: caught budget_test.go:62: equal is inside: the cap of milestone 0: the known money 2.0 plus the session's cap 3.00 passes 5.00 USD; want none
== money-cut: caught budget_test.go:64: above: <nil>; want ErrMilestoneCap with "the known money 2.0 plus the session's cap 3.01 passes 5.00 USD"
== hours-cut: caught budget_test.go:64: the hours and one second: <nil>; want ErrMilestoneCap with "8h0m1s since Start passes 8.0 hours"
== hours-equal-out: caught budget_test.go:62: the hours exactly: the cap of milestone 0: 8h0m0s since Start passes 8.0 hours; want none
== currency-cut: caught budget_test.go:64: a known row not in USD: <nil>; want ErrMilestoneCap with "the session S-00000004 has money in EUR, not USD"
== status-filter-cut: caught budget_test.go:62: an unknown row adds nothing: the cap of milestone 0: the session S-00000003 has money in , not USD; want none
== computed-dropped: caught budget_test.go:64: above: <nil>; want ErrMilestoneCap with "the known money 2.0 plus the session's cap 3.01 passes 5.00 USD"
== session-cap-dropped: caught budget_test.go:64: above: <nil>; want ErrMilestoneCap with "the known money 2.0 plus the session's cap 3.01 passes 5.00 USD"
== check-cut: caught phase_integration_test.go:287: PhaseSession: "S-1e642896", <nil>; want ErrMilestoneCap with "the known money 0.0 plus the session's cap 50.01 passes 50.0 USD" and no sess
== no-session-cap: caught phase_integration_test.go:287: PhaseSession: "S-c9ab0758", <nil>; want ErrMilestoneCap with "the known money 0.0 plus the session's cap 50.01 passes 50.0 USD" and no sess
== head-not-root: caught phase_integration_test.go:305: the time of Start 2026-10-11 05:55:36 +0000 UTC; want the first commit's, 1791460917 (the head's is 1791698136)
== format-ci: caught git_test.go:109: dir "r", args ["-c" "core.hooksPath=/dev/null" "-c" "core.attributesFile=/dev/null" "-c" "core.excludesFile=/dev/null" "-c" "core.autocrlf=false" "-c" "c
```

## Green

At `e3c48f7` (05:56Z to 05:58Z): `go build ./...`, `go vet ./...`, `go test
./...`, `go test -tags=integration ./...` and `go test -tags=e2e ./...` pass,
every package `ok`.

## The fix of round 1 (2026-10-11)

Finding 1: the hours refusal wrote the elapsed time into the detail. Finding
2: the demo had no unit test. The tests first, red (06:14Z): the two hours
cases of `TestMilestoneCap` now ask for the exact detail, and
`TestTheCapOfMilestoneZeroDrawsNoSession` (unit, the stand-in store) asks for
no session ID and no commit; its money case was green at once, as the check was
already in place, and its hours case red for the right reason:

```
--- FAIL: TestMilestoneCap (0.00s)
    budget_test.go:65: the hours and one second: the cap of milestone 0: 8h0m1s since Start passes 8.0 hours; want ErrMilestoneCap with "the cap of milestone 0: the time since Start passes 8.0 hours"
    budget_test.go:65: the hours of a decimal: the cap of milestone 0: 31m0s since Start passes 0.5 hours; want ErrMilestoneCap with "the cap of milestone 0: the time since Start passes 0.5 hours"
--- FAIL: TestTheCapOfMilestoneZeroDrawsNoSession (0.00s)
    phase_test.go:49: the hours: "", the cap of milestone 0: 9h0m0s since Start passes 8.0 hours; want ErrMilestoneCap, "the cap of milestone 0: the time since Start passes 8.0 hours"
```

Then the detail of the hours is fixed text; green (06:15Z). With the check of
`PhaseSession` cut (`&& false`), the unit test fails (the session goes on to
a hook the stand-in has not set, and the test panics): caught.
