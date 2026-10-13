# T-jpkc — row 51a of the plan, the forge calls of M2d

Issue: [#271](https://github.com/pharzam/layup/issues/271), row 51a of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` through REQ-003 and NFR-001. Base `9144e4e` (the plan review's), read
after four merges of `main` against `5e40952`. Author: Claude Sonnet 5.5 on Claude
Code. Evidence: [`runs/T-jpkc/`](../../runs/T-jpkc/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #271. The plan
review (GPT-6.1 Sol on Codex CLI, `high`, 3 min 58 s) gave `approve-with-conditions`:
850 lines over 14 files against `9144e4e`, close-out inside; cap 1; no panel. Its five
conditions are applied: both stand-ins of `forge.Forge` gain the two methods; the
permission of each call is read again from its page (`Metadata` read for the rules,
`Contents` read for the activity: the plan had taken the sentence of the next endpoint)
and the sources are in the record of the runs; the error contract is written, a body
of `null` reading as an empty list; the cases of a branch that does not exist yet, of
the active rules only and of the shared error rules on the new calls are named; and a
validation step comes before the freeze. The goal count is 1, final by O-209 of #267.

## What was done

`forge.go` gains `Rule` (the type and the parameters, nil when GitHub gives none),
`ActivityEntry` (the six fields, a null actor being ID 0 and login `ghost`),
`Repository.Heads` and the methods `EffectiveRules` and `Activity`. The GitHub adapter
reads the rules of a branch 100 a page over every page, the branch escaped segment by
segment (the parameter is multi-segment), an empty list for a branch that no ruleset
names; the activity of a ref, one page of the default size, the ref as given; and the
`commit.sha` of each branch. A body of `null` is an empty list and a body that is no
list is a forge error; the 401 and the rate limit are the handling that every call
shares. The two stand-ins gain the methods. Documents: `forge.md`, `packages.md`, the
traceability row, `PRD-0001`, and the pitfall "A merge of `main` moves the base of the
budget" in `guardrails.md` §2.

**Tests:** [`test-runs.md`](../../runs/T-jpkc/test-runs.md): red for the new methods,
types and field; eleven new tests against the `httptest` forge; each rule mutated alone,
35 of 35 caught after a first run that left one alive and one that did not build; the
fix of round 1 checked by a mutant; the whole gate on the head after each merge of `main`.

**Known limits:** what GitHub gives for a branch that no ruleset names, or that does not
exist yet, is read from its documentation and the read-back's need, and the tests use a
fake server: the first apply on a real target, row 59b (`T-sdtb`), observes it. `Activity`
reads one page of 30. The reviewer's sandbox has no local socket, so the integration
level is the author's run and CI's.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`a775a2e`, cycle 0; comment 6143289886): `material`. Finding 1: the sentence of
`forge.md` and the comment of `EffectiveRules` equated a branch that does not exist yet
with one that has no rules; fixed in `6178bc5`, with a test of its own. Finding 2: the
budget at its named base, answered by the comments on the base of the budget and the
reply 6144221123. A note: the sandbox.

Round 2 (`7a25b0e`, cycle 1 of the cap of 1; comment 6144740724): `not mergeable, findings
recorded`. Finding 1, `material`, stays: read against `9144e4e` the branch is 4,823 lines
over 61 files, as the branch took `main` four times; against `5e40952` it is 700 over 11.
The Operator ruled it a known limit (comment 6144813606): the budget is read against `5e40952`, and the
figure at landing is approved as 797 lines over 14 files. Notes 2 and 3 are applied in `dcb39e4` (the limit of the API evidence; the
heading of the base in the record and the figures of the pitfall); note 4 is the sandbox.

## Verdict

Delivered: `EffectiveRules`, `Activity` and the heads of the branches, on the interface
and on the adapter, with their tests and documents. The review ended at the cap with the
budget as its one open finding. The diff against `5e40952` is 797 lines added plus removed over 14 files with the close-out
inside, against the 850 over 14 of the plan review. Next: `T-wphp` (the step forge-check),
which uses these calls.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13, from the clock; GPT-6.1 Sol tokens are
the "tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 09:05 to 09:12 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 223,785 | 3 min 58 s, 09:11:59 |
| The tests, the code, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 09:20 to 09:54 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 196,076 | 13 min 24 s, 09:54:18 |
| The fix, four merges of `main`, the runs | execution | Claude Sonnet 5.5 | `max` | not reported | 10:00 to 11:30 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 165,704 | 9 min 16 s, 11:30:37 |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 11:41 to 11:46 |
