# The test runs of T-nxe4

The build ran on 2026-10-10 from 04:05Z (the claim of the branch) to 04:20Z, in
the order below; the steps inside it were not stamped one by one.

## Red 1: the pass rules and the routing register, before the code

`go vet ./internal/run` does not build `probe_test.go`, for want of the
function it names:

```
# [github.com/pharzam/layup/internal/run]
vet: internal/run/probe_test.go:37:13: undefined: probeReason
```

## The integration and e2e tests

`probe_integration_test.go` and the e2e of the step were written after
`probe.go`, so their red is the mutations below. Three expected changes
followed the new step: the restart's six steps in `command_test.go` and the
run's integration tests, `run.md`'s list of the restart, and the e2e world:

- a scripted harness in place of `claude` (condition 3);
- a bare repository made with `-b main` (the pitfall "A bare repository with no
  default branch" of `docs/guardrails.md` §2: with no `-b`, the run's clone
  had no local `main`, and the probe's `CloneLocal` failed);
- the comment endpoint of the fake forge.

`TestPackageRules` refused `internal/cli` importing `internal/records`
for the price list, so the reader moved into `internal/run` (`ReadPrices`).

## The mutations

A mutation of each rule, one at a time, run by a script that restores each
file from a copy after each; the unit tests for the pass rules and the routing
comparison, the integration tests for the rest. Each test line in full, cut at
400 characters:

```
== class: probeReason: an end that is not done passes (exit 1)
probe_test.go:38: an end that is not done: "", want "crash"
probe_test.go:44: no valid probe.tsv: "token", want no-result
== token-value: probeReason: another token passes (exit 1)
probe_test.go:38: another token: "", want "token"
== token-count: probeReason: no token passes (exit 1)
probe_test.go:38: no token: "", want "token"
== files: probeReason: no AGENTS.md passes (exit 1)
probe_test.go:38: no AGENTS.md: "", want "files"
== outside: probeReason: a file outside passes (exit 1)
probe_test.go:38: a file outside: "", want "outside"
probe_test.go:38: a relative path that leaves the session directory: "", want "outside"
== outside-policy: probeReason: a policy path is outside (exit 1)
probe_test.go:38: a policy path outside, of the row: "outside", want ""
== not-used: probeReason: a model of use no passes (exit 1)
probe_test.go:38: a model of use no: "", want "not-used"
== routing-order: routingDiffers: the order of the rows counts (exit 1)
probe_test.go:52: the same rows in another order: they differ; want the same table
== routing-copy: ProbeStep: the routing register never copied (exit 1)
probe_integration_test.go:63: git show layup-records:routing.tsv: exit status 128
== routing-always: ProbeStep: the routing register copied when equal (exit 1)
probe_integration_test.go:77: the step again: 0, 0, 0, git commit -m layup run: the routing register: exit status 1; want none probed and no commit
== skip: ProbeStep: a harness with no model of use yes probed (exit 1)
probe_integration_test.go:50: the step: 1, 0, 1, the harness none has no model of use yes; want 1 passed, 1 skipped
== probed: Probe: a version whose last probe passed probed again (exit 1)
probe_integration_test.go:77: the step again: 1, 0, 1, <nil>; want none probed and no commit
== admit-probe: Admit: a probe that fails admits the start (exit 1)
probe_integration_test.go:120: a probe that fails: <nil>, the refusals ["push-refused"]; want probe
== comment: endProbe: no comment (exit 1)
probe_integration_test.go:67: the comments []
probe_integration_test.go:96: the comments []
== prompt-token: Probe: the prompt holds no token (exit 1)
probe_integration_test.go:50: the step: 0, 1, 1, <nil>; want 1 passed, 1 skipped
probe_integration_test.go:116: a version with no probe that passes: refused: probe, ["probe"]
probe_integration_test.go:126: harnesses.tsv [["S-cf122b88" "fake" "2.0.0" "m1" "failed" "no-result" "—" "—" "2026-10-10T04:18:58Z"] ["S-4492d816" "fake" "3.0.0" "m1" "failed" "token" "AGENTS.md" "—" "2026-10-10T04:18:58Z"]]; want the probe of each version
```

A first run had four that were not caught. `outside-relative` and
`outside-clean` had no case that only they could fail; the cases "a relative
path into home/" and "an absolute path whose .. leaves the session directory"
were added. `start-noevent` had no assertion; `TestTheStepProbe` now asserts
that a probe's task has no `events.tsv`. `refused-version` was a redundant
clause (the version is already empty when the version check refused), so the
clause and its mutation are gone. The three reruns, in full:

```
== outside-relative: probeReason: a relative path read from the session directory (exit 1)
probe_test.go:41: a relative path into home/: "outside", want ""
== outside-clean: probeReason: a path not cleaned of .. (exit 1)
probe_test.go:41: an absolute path whose .. leaves the session directory: "", want "outside"
== start-noevent: startRow: a probe's start row with the event session (exit 1)
probe_integration_test.go:68: the probe has events:
```

## Green (2026-10-10T04:20Z)

Each with exit 0: `go build ./...`, `go vet ./...`, `gofmt -l internal cmd`
(empty), `go test ./...`, `go test -tags=integration ./...` (with
`TestPackageRules` and `TestInputRule`), `go test -tags=e2e ./...`;
`adr-lint`, `prd-lint`, `link-lint`, `setup-check`,
`run-discipline-tests`, `git diff --check`, `sh runs/T-fdaq/parts.sh`.

## The fix of round 1

Round 1 (comment 6094848498 of #168), notes 3 and 7. `TestTheStepProbe` now
asserts the `records` column of the probe's start row, the parent of the
commit that adds it, and that the probe that ended at its version check left
no directory. Red at 06:54Z on `918a2f4`'s code, with `ProbeStep` taking a
function read once per step (the copy of the routing register comes first):

```
--- FAIL: TestTheStepProbe (1.40s)
    probe_integration_test.go:65: the records column of the probe's start row: want ["2f8d08bff6bb6242cf0ef029e70cb4c128ca551a"]
```

Green after the read moved into the loop. The directory's mutation, at 06:55Z,
`if made && !errors.Is(err, errProbed)` in `TaskSession` (the early end keeps
its directory), is caught:

```
probe_integration_test.go:92: the session directories after the step again: [d S-4c293f8f/], <nil>; want none
```

Green at 06:58Z, each with exit 0: `go build ./...`, `go vet ./...`,
`gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration
./...`, `go test -tags=e2e ./...`; `adr-lint`, `prd-lint`, `link-lint`,
`setup-check`, `run-discipline-tests`, `git diff --check`; and `sh
runs/T-fdaq/parts.sh` after its two rows of row 38 name the new acceptance row
"The step probe".
