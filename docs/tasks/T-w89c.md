# T-w89c — CI on go1.26.9

Issue: [#150](https://github.com/pharzam/layup/issues/150), found by `T-ywk7`
(#147) in the CI of #149. Serves `F-0003#47` (verification on every change).
Base `ad3c400`. Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-w89c/`](../../runs/T-w89c/).

## Plan and plan review

The plan (R12, comment 6078105078), its review and the author's answer
(6078215407) are comments on #150. The plan review (Claude Fable 5.1, effort
`xhigh`, a fresh read-only session at `ad3c400`, 6 min 20 s; comment
6078215083) gave `approve-with-conditions`: Budget maximum 180 lines over 7
files against `ad3c400`, close-out inside; Cycle cap 2; no panel. Its four
conditions (the sentence of `setup.md`, the toolchain named in each run, the
Operator's step, the demo in one reading) are applied. **O-185 (a)** (#149):
this task first, then #149.

## What was done

`go.mod`: `go 1.26.9`, so `actions/setup-go` installs the toolchain that fixes
the ten standard-library advisories of #150. `docs/spec/setup.md`: the
catalog's `1.26` is the minor line of the directive; the catalog and its tests
stay at `1.26`. A lesson in `docs/guardrails.md` §2.

**Tests** ([`test-runs.md`](../../runs/T-w89c/test-runs.md)): govulncheck red at
`ad3c400` with `GOTOOLCHAIN=go1.26.8` (the ten advisories), green with
`go1.26.9`; the three test levels pass; in CI, run 37912054398, the job
`security` passes and setup-go logs `go1.26.9`.

**Rejected:** a `toolchain` line (not checked that setup-go reads it); a
version in `ci.yml` (a second home); `check-latest: true` (a toolchain that
changes under CI with no commit).

**Known limit:** the red returns at each security release of the 1.26 line,
until the directive moves.

## Review rounds

Round 1 (Claude Fable 5.1, `xhigh`, at `29abdee`, cycle 0; comment on #150):
`nothing material in scope`, three notes, applied in the close-out (the red
alone downloads a toolchain; the cut of the red's block named; the recurrence
in the lesson).

## Verdict

Delivered: CI builds on `go1.26.9`, and the required job `security` is green
on this task's pull request, whose only change to what CI builds is the `go`
directive. The review ended by decay at cycle 0 of a cap of 2, inside 180 lines
over 7 files.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-09; tokens are the
`modelUsage` of the `result` event of the Claude Code CLI; `not reported`
otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The issue, the plan | reasoning | Claude Opus 5.5 | max | not reported | 09:08 to 09:22 |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 573,538 (USD 3.37) | 6 min 20 s |
| The answer; the red, the change, the green; CI | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 09:29 to 09:35 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 474,873 (USD 3.00) | 8 min 13 s |
| The close-out | — | Claude Opus 5.5 | max | not reported | from 09:44 |
