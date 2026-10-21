# 0027. Bound the integration level at 20 minutes in CI

Date: 2026-10-20

## Status

Accepted

## Context

[ADR-0010](0010-use-go-as-the-technology-stack.md) sets the per-test-binary timeout to the Go default, `-timeout 10m`, and [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml) passes `-timeout 10m` to the unit, integration and end-to-end steps of job `tests`.

In CI, the `internal/run` integration test binary passed in 519.198 s at `0246d9c` (run 39389141724) and 493.632 s on PR #436's head `51d949b` (run 39421233368), then timed out at 600.012 s (`panic: test timed out after 10m0s`, with no assertion failure) on `main` at `2866a3f` (run 39422675099, the merge of #436) and at 600.014 s on PR #437's head `614bef4` (run 39425720339), which completed in 482.275 s on the session host ([`runs/T-jc45/test-runs.md`](../../runs/T-jc45/test-runs.md)). A timeout at 600 s is an aborted run, not a completed runtime. The unit and end-to-end levels stay well under 10 minutes (the end-to-end level took 65 s on the session host in [`runs/T-r7b9/test-runs.md`](../../runs/T-r7b9/test-runs.md)).

Remaining M2e tasks add integration tests to `internal/run`, while task `T-05f4` ([#439](https://github.com/pharzam/layup/issues/439)) parallelizes `internal/run`'s integration tests as the lasting fix. On 2026-10-20, the Operator selected option (c) on [#442](https://github.com/pharzam/layup/issues/442): raise CI's integration bound to 20 minutes as a stopgap while #439 proceeds in parallel, serving [`F-0003#47`](../facts/F-0003-layup-psb-numbered-facts-part-2.md) (Verification on Every Change). An accepted ADR is immutable outside its `Status` line ([`README.md`](README.md)), and [R10](../issue-workflow.md#r10--sync-with-governance) requires the governance documents and CI to agree.

## Decision

We will amend [ADR-0010](0010-use-go-as-the-technology-stack.md) so that CI's integration level uses a **per-test-binary** limit of **20 minutes** (`go test -tags=integration -timeout 20m ./...`). As with `-timeout 10m` (`go help testflag`), `-timeout 20m` bounds each compiled package test binary in the integration invocation; it is neither a whole-job limit nor an individual-test limit.

The unit (`go test -timeout 10m ./...`) and end-to-end (`go test -tags=e2e -timeout 10m ./...`) invocations in CI and the [`pre-commit` hook](../engineering-discipline.md#git-hooks) keep `-timeout 10m`. A hanging integration test binary still fails, after 20 minutes.

We reject:
- **Editing ADR-0010's body:** an accepted ADR is immutable except its `Status` line.
- **Changing `.github/workflows/ci.yml` alone:** leaves ADR-0010 and the testing and guardrail documents contradicting CI ([R10](../issue-workflow.md#r10--sync-with-governance)).
- **Splitting `internal/run` across CI steps by `-run` patterns:** a new test whose name misses the patterns can fall between them and never run in CI.
- **Waiting for #439 alone:** leaves `main` red until #439 lands and contradicts the Operator's selection of option (c) on [#442](https://github.com/pharzam/layup/issues/442).
- **Widening the unit or end-to-end bounds too:** both run well under 10 minutes, so a wider bound would only delay failure on a hang.

No panel was convened: this record changes CI test configuration, not product architecture ([ADR-0026](0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md)).

## Consequences

- A hanging integration test binary costs up to 20 minutes of CI before failing.
- ADR-0010's `Status` becomes `Accepted. Amended by ADR-0027`, and the restatements in [`tests/test-levels.md`](../tests/test-levels.md), [`tests/scaling-checklist.md`](../tests/scaling-checklist.md), [`engineering-discipline.md`](../engineering-discipline.md) and [`guardrails.md`](../guardrails.md) link to this record as the home of the integration bound ([one home per rule](../engineering-discipline.md#one-home-per-rule)).
- This is a finite stopgap, not evidence of capacity for later rows. Once [#439](https://github.com/pharzam/layup/issues/439) parallelizes `internal/run` and CI measures `internal/run` well under 10 minutes, a later ADR may restore `-timeout 10m`.
