# T-xs4y — the technical specification of M2e, the task loop and the typed handoffs

Issue: [#333](https://github.com/pharzam/layup/issues/333). Milestone `M2e` of the
[implementation plan](../plan/README.md#milestones), which starts with its
specification task (O-114), taken by the Operator's directive ("go ahead with M2E
until finish it") after the demo of `M2d` was accepted. Serves `F-0003#51`
(`REQ-012`) for `REQ-005`, `REQ-007`, `REQ-012`, `REQ-003`, `NFR-004` and
`NFR-001`. Base `05988c5` (the merge of #332). Author: Claude Sonnet 5.5 on Claude
Code. Evidence: [`runs/T-xs4y/`](../../runs/T-xs4y/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #333. The plan
review (GPT-6.1 Sol on Codex CLI, effort `xhigh`, a fresh sandboxed session in a
clone at `05988c5`, 4 min 36 s; the harness is outside the routing rule of ADR-0026
and was allowed by the Operator on 2026-10-13) gave `approve-with-conditions`:
Budget maximum 4,200 lines added plus removed over 34 files, close-out inside;
Cycle cap 1; no panel. Its five conditions are applied: the deferrals named in S3;
the first-merge tree rule and the whole of F-22; the new blocks told from the built
ones, one owner each; the boundaries with `M2f`, `M2g` and the slicing task; the
corrections before the freeze.

## What was done

1. **S1, the public-solution search** ([`runs/T-xs4y/search/`](../../runs/T-xs4y/search/)):
   five lists, 24 recorded searches, two blind searchers (A: GPT-6 Sol on Devin;
   B: Claude Sonnet 5.5 on Claude Code). Both conclude that no public solution can
   be the base of `M2e`; what transfers is GitHub's status, check-run and merge API.
2. **S3, the decisions** (comment 6191483681, corrected in the answers to rounds 1
   to 4): the artifact that is a comment (#68), the task register and the
   transition table, independence, the four statuses, the merge and F-22, Shape,
   the bet, the full check, the parameters and the Retrospective.
3. **T:** `runs/T-xs4y/sections.sh` red on 32 headings before the files existed,
   33 green at the end; the block test red on ten new blocks, green with their
   names in `notYetBuilt` ([`test-runs.md`](../../runs/T-xs4y/test-runs.md)).
4. **D1 to D4:** [`docs/spec/loop.md`](../spec/loop.md) (new), ten blocks in
   `records.md`, the calls of `M2e` in `forge.md`, the table of `M2e` in
   `packages.md`, and the documents that named `M2e` (`architecture.md` §3, §6 and
   §8, `README.md`, `gate.md`, `setup.md`, `run.md`, `intake.md`, `session.md`,
   `scaffold.md`, the plan, `PRD-0001` §12, the lesson of `guardrails.md` §2).

**Revealed, off the path, as issues:** #336 (amend ADR-0017, ADR-0019 decision 2
and ADR-0014 decision 3 for the code-owner review and checks `read`) and #345 (the
activation batch against the native jobs that read the base, for `M2f`).

## Review rounds

The records are comments on #333. All five rounds were Claude Opus 5.5, effort
`xhigh`, a fresh read-only `claude -p` session in a clone, a model that wrote none
of the change.

| Round | Head | Cycle | Cap | Material | Then |
| ----- | ---- | ----- | --- | -------- | ---- |
| 1 | `cd2838f` | 0 | 1 | 9 | fixes; round 2 (cycle 1) is the last the cap allows |
| 2 | `7b3a4c3` | 1 | 1 | 6 | **O-210**: cap 2 |
| 3 | `757b1bf` | 2 | 2 | 6 | **O-211**: cap 3 |
| 4 | `c89bf1d` | 3 | 3 | 4 | **O-212**: cap 4 |
| 5 | `5974415` | 4 | 4 | 4 | **O-213**: settled as fixed |

Round 5 keeps its verdict, `not mergeable, findings recorded`. By **O-213**, the
Operator settled its findings as fixed by `03f8c87`, as O-207 (b) did for `T-msrd`,
and no round reads that commit. The cap records of rounds 1 to 4 are the `## Plan
review` comments that follow each.

## Verdict

Delivered: the specification of `M2e`, the task loop, in `docs/spec/loop.md`, with
ten blocks, the calls and the table of `M2e`, after a public-solution search by two
blind searchers.

- **How the review ended:** at cycle 4 of a cap that the Operator raised three times
  (O-210, O-211, O-212), by O-213 and not by decay: the material findings per round were 9, 6, 6, 4 and 4.
- **The size:** the diff against `05988c5` is inside 4,200 lines over 34 files.
- **Next:** the slicing task of `M2e`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15. Reviewers' and searchers'
tokens are as each harness reports them (Devin reports none); the author's are not
reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | not reported | not reported | 11:22 (the issue) to 11:38 (the plan's comment) |
| The plan review | reasoning | GPT-6.1 Sol | `xhigh` | 165,119 | 4 min 36 s |
| S1, the sources and the first tries | execution | Claude Sonnet 5.5 | not reported | not reported | 11:24 to 11:51 |
| Searcher A | reasoning | GPT-6 Sol (Devin) | — | not reported by Devin | 2 min 55 s |
| Searcher B | execution | Claude Sonnet 5.5 | — | 3,104,990 (USD 2.55) | 3 min 42 s |
| S3, T, D1 to D4 | execution | Claude Sonnet 5.5 | not reported | not reported | 11:55 to 12:21 |
| Round 1 | reasoning | Claude Opus 5.5 | `xhigh` | 9,302,130 (USD 5.56) | 18 min 12 s |
| The fix of round 1 | execution | Claude Sonnet 5.5 | not reported | not reported | 12:40 to 12:45 |
| Round 2 | reasoning | Claude Opus 5.5 | `xhigh` | 8,612,824 (USD 4.88) | 12 min 55 s |
| The fix of round 2 | execution | Claude Sonnet 5.5 | not reported | not reported | 12:59 to 13:08 |
| Round 3 | reasoning | Claude Opus 5.5 | `xhigh` | 10,353,779 (USD 5.85) | 17 min 46 s |
| The fix of round 3 | execution | Claude Sonnet 5.5 | not reported | not reported | 13:26 to 13:43 |
| Round 4 | reasoning | Claude Opus 5.5 | `xhigh` | 10,688,856 (USD 6.04) | 15 min 2 s |
| The fix of round 4 | execution | Claude Sonnet 5.5 | not reported | not reported | 13:59 to 14:03 |
| Round 5 | reasoning | Claude Opus 5.5 | `xhigh` | 8,576,090 (USD 5.25) | 16 min 14 s |
| The fix of round 5 | execution | Claude Sonnet 5.5 | not reported | not reported | 14:20 to 14:30 |
| The close-out | execution | Claude Sonnet 5.5 | not reported | not reported | 14:30 to 14:35 |
