# T-w73g — the survey of five public repositories, as an input of the later specification tasks

Issue: [#133](https://github.com/pharzam/layup/issues/133). Base `add60ef`.
Author: Claude Opus 5.5 on Claude Code. Serves `F-0003#52`, `F-0003#49` and
`F-0003#50` through the specification tasks that read the survey. Evidence:
[`runs/T-w73g/`](../../runs/T-w73g/).

On 2026-10-07 the Operator asked for four public repositories to be checked for
ready components and patterns for LAYUP. The patterns belong to milestones that
have no specification yet, so the Operator chose option A in the session: record
the findings as inputs, and write no code before those specification tasks.
During the work the Operator added a fifth repository (the scope amendment on
#133).

## Plan and plan review

The plan (R12) and its review are comments on #133. The plan review (Claude
Fable 5.1, a fresh subagent session on Claude Code) gave
`approve-with-conditions`: Budget maximum 900 lines over 12 files against
`add60ef`; Cycle cap 1. The author applied its seven conditions: the survey is
not the public-solution search; one citation form, `<repo>:<path>`; a red run for
the right reason and a mutation run; Jev as one row of `M3a`; the claims "Node or
TypeScript" and "state outside Git" with their evidence; the brief as evidence;
the plan names only a milestone with a pattern row. The scope amendment for the
fifth repository set the figure to 1,000 lines over 13 files.

## What was done

1. **The test first.** [`cites.sh`](../../runs/T-w73g/cites.sh) checks that each
   source path of the survey exists in its repository at the recorded commit; it
   fetches each repository at its commit (tree objects only) or reads local
   clones. The red, mutation and green runs are in
   [`test-runs.md`](../../runs/T-w73g/test-runs.md).
2. **The evidence.** Five read-only agents (Claude Opus 5.5) read one repository
   each; their briefs and reports are kept word for word
   ([`brief.md`](../../runs/T-w73g/brief.md), `report-*.md`). No code of a
   repository was run.
3. **The survey.** [`survey.md`](../../runs/T-w73g/survey.md): the five
   repositories with their commits and licenses, the patterns by milestone
   (`M2b`, `M2f`, `M2g`, `M3a`, `M3d`, the learning loop), what not to take, and
   the README claims that the code does not support.
4. **The plan.** [`docs/plan/README.md`](../plan/README.md), "Milestones": one
   paragraph names the survey as an input of those specification tasks and of
   the Operator's decision on the learning loop.
5. **The lesson.** [`docs/guardrails.md`](../guardrails.md) §2: "A citation check
   that proves the path, not the claim".

## Review rounds

| Round | Commit | Reviewer | Verdict |
| ----- | ------ | -------- | ------- |
| 1 | `bdb276f` | Claude Fable 5.1, a fresh subagent session on Claude Code | `nothing material in scope`, nine notes |

The close-out applies the nine notes: the count of the Rust files of `ruflo`; the
budget script of `ruflo` exits 1 at 100 %; receipts are appended, state files are
written atomically; two more cited files; prices with a URL or a marked estimate;
the fifth outcome `throttle`; "half of the hook commands"; the header of
`cites.sh` says that it copies the commits of the table; the issue's goal text
names the five milestones. `cites.sh` is green again on 67 paths
(`test-runs.md`, run 4).

## Verdict

Delivered. None of the five repositories gives a component that LAYUP can use as
it is: each is mainly Node or TypeScript, and each keeps the runtime state of its
agents outside Git. They give patterns, recorded with their sources, for the
specification tasks of `M2b`, `M2f`, `M2g`, `M3a` and `M3d` and for the
Operator's decision on the learning loop. `kolezka/self-improvement-loop` is
PolyForm Noncommercial: its rows are marked "idea only". The survey decides
nothing, and the order of the plan does not change.

## Known limits

- One model read each repository, and the review checked a sample of 21 rows;
  the other rows are the agents' reading, not checked by hand.
- The repositories were read, not run; no claim of behaviour is tested.
- `cites.sh` copies the five commits of the survey's table; an edit of one must
  be made in the other.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-07. Tokens of a
subagent are its task report; the main session's tokens are `not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The five survey agents (before the issue) | reasoning | Claude Opus 5.5, Explore subagents | not reported | 92,032; 73,785; 78,375; 100,833; 75,681 | 100.6 s; 139.2 s; 100.5 s; 120.6 s; 64.4 s |
| The plan, the answer, the scope amendment | reasoning | Claude Opus 5.5 | not reported | not reported | 10:30 to 10:54 |
| The plan review | reasoning | Claude Fable 5.1, a subagent | not reported | 119,668 | 6 min 54 s, from 10:44 |
| The test, the survey, the plan paragraph | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | not reported | not reported | 10:54 to 10:58 |
| Round 1 | reasoning | Claude Fable 5.1, a subagent | not reported | 163,459 | 9 min 56 s, from 10:59 |
| Isolate, guardrails, docs, close-out with the notes of round 1 | `—` | Claude Opus 5.5 | not reported | not reported | 11:09 to 11:20 |
| **Total** | | | | 703,833 reported, plus the main session (not reported) | |
