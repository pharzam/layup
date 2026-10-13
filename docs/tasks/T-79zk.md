# T-79zk — row 55 of the plan, the step verify

Issue: [#278](https://github.com/pharzam/layup/issues/278), row 55 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#42` and `F-0001#1` through REQ-002 and NFR-001. Base `6dec793`. Author: Claude
Sonnet 5.5 on Claude Code. Evidence: [`runs/T-79zk/`](../../runs/T-79zk/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #278. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 6 min 5 s) gave `approve-with-conditions`: 1,000 lines
over 22 files against `6dec793`, close-out inside; cap 1; no panel. Its three conditions
are applied: the failure contract (an error of `verify.Run`, a `CleanupError` with passing
rows too, fails the step and publishes no commit; each outcome is a unit test with fakes),
the commit's wording (it writes five paths and keeps the other records) and the gate and
close-out in the size. The plan was written by the session's model, an execution-tier one: a
choice, not an unavailability of the reasoning tier.

## What was done

The step `verify` of `layup run` (`verify.go`), after `setup`: it builds the work area again
from the records, runs the loop of the runner, `layup setup verify` in process, S15, and
commits `setup/record.tsv`, `setup/verify.tsv`, `rule-paths.tsv`, `setup/tree.tsv` and the README
of the Scaffold in one records commit through fencing; `setupStep` and `verifyStep` share
`setupRecorded` and `setupArea`. The real runner found a defect of S12 in phase 2 (its record row
`module` cited an answer that phase 2 has not), fixed in `internal/setup` with its specification
row. The specification, the traceability, the PRD cells and `guardrails.md` §2 (one pitfall) are updated.

**Tests:** [`test-runs.md`](../../runs/T-79zk/test-runs.md): red first for each part; 42 mutants,
each alone on a copy of the tree, all caught in the end; the whole gate on the head.

**Known limits:** `layup setup verify` runs in the world of the demo on a fixture baseline; a real
target's verify is first seen in the demo of row 59b (`T-sdtb`). A restart that verifies runs S01
to S14 twice (about 2.6 s each in the world); the integration level of `internal/run` took 510 s,
52 s more than at the end of T-eahr, against the 600 s of CI's limit for a package.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`37358f3`, cycle 0; comment 6149719082): `material`, three findings in the change: the guard of S15 omitted
the rows `verify.sha256` and `rule-paths.sha256`, the hashes and an empty register that the accepted Condition 1 named;
the five-path contract had two readings on recovery; no case of a row `not-active`; and a note on a red-first claim. All
fixed (`07292cf`, `41867b0`; the answer is comment 6150123673).

Round 2 (`0e1cd25`, cycle 1 of the cap of 1; comment 6150387549): `nothing material in scope`; one note, that the record gave
no separate timing of the second run of S01 to S14, which the answer to the plan review promised: applied (2.5 to 2.7 s).

## Verdict

Delivered: the step `verify` of `layup run`, with its tests, documents and runs, and the fix of S12 in phase 2 that its real
run found. The review ended at `nothing material in scope` after two rounds, at the cap of 1. The diff against `6dec793` is
991 lines over 21 files with the close-out inside, against the plan review's 1,000 over 22. Next: `T-hfn1` (the step
scaffold-push), which builds a missing work area again with `setupArea`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13, from the clock; the tokens of GPT-6.1 Sol are the
"tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 14:43 to 14:48 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 151,454 | 6 min 5 s, 14:47:57 |
| The tests, the code, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 14:55 to 15:47 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 216,489 | 13 min 13 s, 15:47:48 |
| The fixes of round 1, the runs | execution | Claude Sonnet 5.5 | `max` | not reported | 16:01 to 16:23 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 201,679 | 13 min 36 s, 16:23:05 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 16:37 to 16:39 |
