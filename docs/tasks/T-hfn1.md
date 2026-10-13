# T-hfn1 — row 56 of the plan, the step scaffold-push

Issue: [#279](https://github.com/pharzam/layup/issues/279), row 56 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` through REQ-003 and NFR-001. Base `720a149`. Author: Claude Sonnet 5.5 on
Claude Code. Evidence: [`runs/T-hfn1/`](../../runs/T-hfn1/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #279. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 3 min 58 s) gave `approve-with-conditions`: 900 lines
over 22 files against `720a149`, close-out inside; cap 1; no panel. Its four conditions are
applied: the framing of `setup/rulesets.txt` (the output byte for byte, the marker line and
a separator line feed LAYUP's, an R10 correction in `scaffold.md` and `records.md` before
the tests), a missing work area rebuilt before either outcome, the whole rule of the Operator
printed, and the validation cases and the full gate. The plan was written by the session's
model, an execution-tier one: a choice, not an unavailability of the reasoning tier.

## What was done

The step `scaffold-push` of `layup run` (`scaffoldpush.go`), after `verify`: "done before" when
`setup/rulesets.txt` is in the records; else it makes sure of the work area (built again as
`setup` does, its head checked against the last commit of `setup/tree.tsv`, when it is missing),
reads the forge's head and the Operator's two apply outputs, and either waits, printing the
rule of the Operator and each line of `commands.sh` (the last row of the restart, which releases the
lease), or checks that each output is a JSON object of the repository and commits `setup/rulesets.txt`
in one commit through fencing. `readApplies`, `rulesetsRecord` and `scaffoldJob` (the seams, as
`verifyJob`) are its parts. The specification, the traceability, the PRD cells and `guardrails.md` §2
(one pitfall) are updated.

**Tests:** [`test-runs.md`](../../runs/T-hfn1/test-runs.md): red first for each part; 38 mutants,
each alone on a copy of the tree, all caught in the end; the whole gate on the head.

**Known limits:** the world runs the step with a stand-in `layup setup verify` that passes; the Operator's
real pushes and applies are first seen in the demo of row 59b (`T-sdtb`). The step `forge-check` (row 57a)
is not built: the step ends the restart. The integration level of `internal/run` took 573 s here against the
600 s of CI's limit for a package.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`6dd116e`, cycle 0; comment 6151696485): `material`, one finding in the change: condition 4 of the plan review
asked for an integration case of an apply output that cannot be read, and only a directory had one. Fixed (`5cb9c94`;
the answer is comment 6151968755).

Round 2 (`7d816c7`, cycle 1 of the cap of 1; comment 6152151770): `nothing material in scope`; one note, a comment of the
world of row 55 that still said a second restart builds no work area: applied.

## Verdict

Delivered: the step `scaffold-push` of `layup run`, with its tests, documents and runs. The review ended at `nothing material
in scope` after two rounds, at the cap of 1. The diff against `720a149` is 787 lines over 18 files with the close-out inside,
against the plan review's 900 over 22. Next: `T-wphp` (the step forge-check), which reads `setup/rulesets.txt` for the bypass lists.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13, from the clock; the tokens of GPT-6.1 Sol are the
"tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 16:45 to 16:49 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 123,282 | 3 min 58 s, 16:49:02 |
| The tests, the code, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 16:54 to 17:42 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 176,473 | 6 min 56 s, 17:42:40 |
| The fix of round 1, the runs | execution | Claude Sonnet 5.5 | `max` | not reported | 17:49 to 18:05 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 213,281 | 9 min 44 s, 18:05:18 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 18:15 to 18:16 |
