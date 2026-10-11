# T-gmey — test runs

The session of a phase (#215, row 42a).

## Red (2026-10-11)

The tests before the code, against stubs of the new names that give an error,
so each test fails at its assertion and none at the build.

`go test -tags=integration -run 'TestMakeInput|TestReadOutput' ./internal/session/` (04:58Z):

```
--- FAIL: TestMakeInput (0.00s)
    phase_integration_test.go:23: not built
--- FAIL: TestReadOutput (0.01s)
    phase_integration_test.go:64: an output file a directory: not built, want an error with "is not a regular file"
    phase_integration_test.go:64: an output file over 1 MiB: not built, want an error with "is over 1 MiB"
    phase_integration_test.go:64: an output file missing: not built, want an error with "no such file"
    phase_integration_test.go:64: an output file a link: not built, want an error with "too many levels of symbolic links"
    phase_integration_test.go:71: an output file of 1 MiB: 0 bytes, not built
    phase_integration_test.go:74: an output named ../meaning.tsv: not built; want it refused
```

`go test -tags=integration -run 'Phase|OutputFile|NoAdmittedPair' ./internal/run/` (04:59Z):

```
--- FAIL: TestASessionOfAPhaseCommitsItsRecordsAndItsTelemetryRow (1.27s)
    phase_integration_test.go:109: PhasePair: not built
--- FAIL: TestAnOutputFileThatIsRefusedIsNoRecord (5.90s)
    --- FAIL: TestAnOutputFileThatIsRefusedIsNoRecord/missing (1.64s)
        phase_integration_test.go:170: PhasePair: not built
    --- FAIL: TestAnOutputFileThatIsRefusedIsNoRecord/a_link (1.28s)
        phase_integration_test.go:170: PhasePair: not built
    --- FAIL: TestAnOutputFileThatIsRefusedIsNoRecord/over_1_MiB (1.51s)
        phase_integration_test.go:170: PhasePair: not built
    --- FAIL: TestAnOutputFileThatIsRefusedIsNoRecord/refused_by_its_check (1.47s)
        phase_integration_test.go:170: PhasePair: not built
--- FAIL: TestNoAdmittedPairStartsNoSession (2.57s)
    phase_integration_test.go:203: a role with no row: not built; want ErrNoPair
    phase_integration_test.go:203: a model of use no: not built; want ErrNoPair
```

`TestARefusedStartOfASessionOfAPhaseIsNoRecord` was written after the code, so
it has no red of its own; the mutation `phase-refusal-as-event` below is the
evidence that it can fail.

## Mutations (2026-10-11, 05:02Z to 05:07Z)

Each rule mutated alone, a condition cut with `&& false`, the six tests above
run: 16 of 16 caught. The first run built 15 and caught 15; the mutant
`keep-ignored` did not build (`kept` unused), was not counted, and was
rewritten (`maps.Copy(files, map[string][]byte{})`), then caught.

```
== keep-refusal: caught phase_integration_test.go:173: PhaseSession: "S-28a02019", <nil>; want a PhaseError of the session, refused file
== no-telemetry: caught phase_integration_test.go:122: the commit of the records changed ["intake/phase.tsv"]; want the records and the telemetry row in one
== no-comment: caught phase_integration_test.go:148: the comments []
== phase-error-dropped: caught phase_integration_test.go:173: PhaseSession: "S-c5f09449", <nil>; want a PhaseError of the session, no-result
== pair-unprobed: caught phase_integration_test.go:215: a harness never probed: <nil>; want ErrNoPair
== pair-model-use: caught phase_integration_test.go:203: a model of use no: <nil>; want ErrNoPair
== no-input: caught phase_integration_test.go:112: PhaseSession: "S-4c0acec1", the session S-4c0acec1: crash
== phase-prompt-absent: caught phase_integration_test.go:112: PhaseSession: "S-03da9fb0", the session S-03da9fb0: crash
== phase-refusal-as-event: caught phase_integration_test.go:231: PhaseSession: "S-c9cfcc9f", refused: version; want a PhaseError of the session, refused version
== phase-start-with-event: caught phase_integration_test.go:144: a session of a phase wrote tasks/T-flg6/events.tsv
== phase-attempt-check: caught phase_integration_test.go:112: PhaseSession: "S-122c6dfb", the session S-122c6dfb: refused attempt
== name-check: caught phase_integration_test.go:39: an input named "": open /home/layup/.cache/gotmp/TestMakeInput3147263142/002/sessions/S-1a2b3c4d/input: is a directory; want it refused
== input-mode: caught phase_integration_test.go:32: input/problem-statement.md: the mode -rw-r--r--, <nil>; want 0600
== no-size-limit: caught phase_integration_test.go:173: PhaseSession: "S-5be307af", <nil>; want a PhaseError of the session, no-result
== follow-link: caught phase_integration_test.go:173: PhaseSession: "S-dcbf7cca", <nil>; want a PhaseError of the session, no-result
== keep-ignored: caught phase_integration_test.go:117: git show layup-records:intake/phase.tsv: exit status 128
```

## Green

At `5d1ee05` (05:07Z to 05:10Z): `go build ./...`, `go vet ./...`, `go test
./...`, `go test -tags=integration ./...` (every package `ok`) and `go test
-tags=e2e ./cmd/...` (`ok`, 54 s) pass.

## The fix of round 1 (2026-10-11)

Note 5 of round 1: a comment that fails dropped the session's reason.
`TestACommentThatFailsKeepsTheReasonOfTheSession` first, red (05:23Z):

```
--- FAIL: TestACommentThatFailsKeepsTheReasonOfTheSession (1.90s)
    phase_integration_test.go:252: PhaseSession: "S-867dc94a", Comment: status 502: Bad Gateway; want the PhaseError of the session, refused file, with the error of the comment
```

Then `endPhase` joins the error of the comment with the `*PhaseError`; green
(05:24Z), with `go test -tags=integration ./internal/run/ ./internal/session/`.
