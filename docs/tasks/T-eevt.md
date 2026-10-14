# T-eevt — row 58 of the plan, the e2e of the Scaffold

Issue: [#282](https://github.com/pharzam/layup/issues/282), row 58 of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#42` and `F-0003#43` through REQ-002, REQ-003 and NFR-001. Base `a880936`.
Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-eevt/`](../../runs/T-eevt/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #282. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 5 min 52 s) gave `approve-with-conditions`: 800 lines
over 14 files against `a880936`, close-out inside; cap 1; no panel. Its four conditions are
applied: the red runs read by their actual cause (the first red was the marker session before
the batch, not the setup), the surface assertions told apart (the details of the tables hold no
path of the host, the printed commands hold absolute ones, the later start's exceptions, the
lease at each exit), the record witness completed (the verification's checks in order,
`rule-paths.tsv`, the README, `setup/rulesets.txt` byte for byte, a new table in Git history), and
the Go environment of the real verification. The plan is reasoning-tier work that the session's
model, an execution-tier one, wrote under the Operator's routing message of 2026-10-13: a routing
mismatch by choice, as no reasoning-tier model was run for it.

## What was done

`TestTheStepsOfTheScaffold` (`cmd/layup/scaffold_e2e_test.go`): the built binary, a fake forge on
loopback, local bare repositories and the world's scripted harness, in five stages: the restart to
the batch on a scaffold baseline built into the pin of a second binary; the answers, then the
restart that runs the setup (the prose sessions) and the real verification, and waits; the test as
the Operator, running the printed `commands.sh` with a stand-in `gh`; the restart that records the
applies and runs the forge check (the fake forge gives the heads, the rules from the work area's
ruleset files and an empty activity; a hook refuses the probe with `GH013` and the ruleset's name);
a later start that passes and one that stops at a row that is not a pass. The world gains the
variant `newRunWorldWith`, the method `runWith` (the Go environment, a limit of three minutes) and
the endpoints of the fake. The specification ("Decided here" under the acceptance rows), the
traceability, the PRD cells and `guardrails.md` §2 (one pitfall) are updated.

**Tests:** [`test-runs.md`](../../runs/T-eevt/test-runs.md): the scenario run at each step, with the
cause of each red read from its output; the stages that earlier rows built passed on arrival and are
recorded as guards; 14 mutants, each alone on a copy of the tree, all caught; the whole gate on the head.
The scenario found no defect in the steps of the Scaffold.

**Known limits:** the scenario shows the binary's plumbing with stand-ins, not that GitHub holds the
rules or refuses the probe in those words (the demo of row 59b does). The `gh` stand-in answers an apply
with its input and the source of the repository, and is put first on the `PATH` of `commands.sh`, as the
`gh` of a host would otherwise answer. The real verification needs the Go toolchain and its build cache
on the host. The e2e level takes about 121 s here (84.6 s before) against CI's limit of 10 minutes.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`a631df1`, cycle 0 of the cap of 1; comment 6164621188): `material`: five findings of assertions that the scenario lacked or made
too weak (the whole sequence of rows, the leaks in the outputs, the gate rows, the published files and the README, the stand-in `gh`) and
three notes. The reply is comment 6164894885: each assertion was written first, run green, and shown red against a mutant that the old one
let pass (8 caught by the new scenario, 8 survived the old one).

Round 2 (`0ef5c1c`, cycle 1 of the cap of 1; comment 6165129972): `nothing material in scope`, with three notes: the direction of a reference
in `scaffold.md` and the count of a paragraph of the record (both applied), and the sandbox (no change).

## Verdict

Delivered: the e2e of the Scaffold, with its world, documents and runs. The review ended at `nothing material in scope` after two rounds, at
the cap of 1. The diff against `a880936` is 661 lines over 11 files with the close-out inside, against the plan review's 800 over 14. Next:
`T-5xjc` (the review of the release of M2d, row 59a).

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-14, from the clock; the tokens of GPT-6.1 Sol are the
"tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 07:12 to 07:14 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 143,986 | 5 min 52 s, 07:14:57 |
| The specification, the scenario, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 07:22 to 08:23 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 233,373 | 12 min 37 s, 08:23:00 |
| The fix of round 1, the mutants, the gate | execution | Claude Sonnet 5.5 | `max` | not reported | 08:35 to 09:09 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 192,238 | 13 min 8 s, 09:09:18 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 09:22 to 09:23 |
