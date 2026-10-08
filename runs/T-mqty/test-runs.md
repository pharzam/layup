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
