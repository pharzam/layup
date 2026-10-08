# T-mqty — the test runs

Evidence for row 26 (#131): each test red first, for the right reason, then green.

## `internal/run` and `internal/records`: red, then green (2026-10-08T11:58Z)

`TestStartAndRestartCallStepAtTheStartOfEachStep`, `TestCheckGit`, `TestCheckValues` (`internal/run`) and `TestCheckValueGivesTheFormOfAName` (`internal/records`) before the code: neither package compiled (`cfg.Step undefined`, `undefined: CheckValue`). Then `Config.Step` (called at the start of each step, the defect of round 2 of #142), `CheckGit`, `CheckValues` and `records.CheckValue`: both packages pass.

## The adapter's own client: red, then green (2026-10-08T12:02Z)

`TestTheOwnClientFollowsNoRedirectAndNoProxy` (a server that answers `302` to another server; the transport's `Proxy`) on the adapter of `075624f`: `panic: runtime error: invalid memory address or nil pointer dereference` (no client). Then `New` makes a client with no redirect and a transport with no proxy: `go test -tags=integration ./internal/forge/github/` passes. A proxy of the environment is not shown by a request, as Go's proxy rule never takes a proxy for a loopback address; the test reads the transport.

## `internal/cli`: red, then green, and three mutations (2026-10-08T12:15Z)

`TestTheSelectingFlagAndTheOptionalFlag`, `TestANoteOfAWaitStartsTheBeatAgain`, and the tests of `run_test.go` (each input check gives exit 2 before any step, the `Config`, the table and its exit codes, a relative `--host`, Ctrl-C, the pin variables against `armature.pin`) before the code: `go vet ./internal/cli` did not compile (`unknown field selector in struct literal of type command`, `undefined: runStart`). Then the code: the unit tests passed at once, so each was shown to fail on a mutation of a backup copy, put back after (`cmp` equal):
- the selecting flag never read → `the row of selector "", arguments [], error unknown flag "--new"` (a first form of this mutation did not compile, an unused variable, and proved nothing; it was redone);
- `--host` not made absolute → `a relative --host: exit 0, Dir "001"`;
- `note` that does not start the beat again → `1 tickers started, want a ticker for step 2`.

`TestPackageRules` (`cmd/layup`) then failed for the right reason: `internal/cli imports internal/forge, …, internal/run, which its row does not allow`. The row of `internal/cli` in the table of phase 1 of `packages.md` gets the four packages (the change that the table of M2a gives to "the build task that needs" it), and it passes.

## The e2e test (2026-10-08T12:30Z)

`cmd/layup/run_e2e_test.go` was written after the command, so its red is shown on the commit before the command, `de455c8`, with the file copied into a work tree of it: `layup run --new: exit 2` and `stderr "layup: unknown command \"run\"…"` (both tests). On the head: `TestRunNewThenRestart` (Start in a world, the same arguments in a second world give the same standard output and exit code, the restart twice in the first world) and `TestRunUsageAndInputErrors` pass, in 26 s, the wait of step 4 being real (ten seconds a Start). Then `go test ./...`, `go test -tags=integration ./...` and `go test -tags=e2e ./...` pass.

## The documents (2026-10-08T12:40Z)

Red: `runs/T-mqty/docs.sh` on a work tree of the e2e commit `dd259bc` gave 19 `FAIL` lines of 20; the rule of the May import cell of `internal/cli` passed there, as the code commit `c6d3354` changed that cell for `TestPackageRules`. The check `adapted` of `setup-check.sh` refused the word "optional" in `README.md` and the glossary (its rule 3); the text says "a flag that may be left out", the words of `run.md`. Then 20 `ok`, exit 0.

## Fixes of round 1 (2026-10-08T12:45Z)

Finding 1: `TestANoteOfAWaitTakesThePlaceOfTheNextBeat` (a beat, a note, the wait's own ten seconds, a second note, two beats; one ticker for the step) on the code of `09a0377`: `panic: test timed out after 10s` at `progress_test.go:84`, a tick on the ticker that the old `note` had stopped (the old `note` started a new ticker, so the beat came just before each wait line). Then `note` marks the line and the next beat prints nothing, with no new ticker, and a beat counts from the start of the step (note 7 too). The first green run failed once on the test's own timing (a tick hands over before its beat has printed); the test now waits for each beat through a hook that is nil in the product, and three runs with `-race` pass.

Notes 3, 5, 6 and 9: a brief with a marker is read (unit); the e2e checks `pin.source` and `pin.commit` of `start.tsv`, so a name that `-X` did not find fails it; `TestMain` removes the build of the scenarios of `layup run`; the unit test of `Step` checks each step name of Start and of the restart. Note 4: one sentence of `README.md`. Then `go test`, `-tags=integration` and `-tags=e2e` pass.

## Fixes of round 2 (2026-10-08T12:55Z, O-176 a, cap 2)

Two new rules of `docs.sh` (the sentence of `run.md` on the beat after a line of a wait; the traceability row of `TestANoteOfAWaitTakesThePlaceOfTheNextBeat`) failed on a work tree of `eb913f5`, the two findings of round 2, and pass after the fix: 22 `ok`, exit 0. The record above keeps the old test name, as it records the run of that time. No Go file changed.
