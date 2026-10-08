# T-trej — row 25a of the plan, the rules of a run

Issue: [#130](https://github.com/pharzam/layup/issues/130), opened by `T-zwke`
(#124) as row 25 of [the tasks of M2a](../plan/README.md#the-tasks-of-m2a);
the Operator's decision O-173 (a) split it into 25a (this task) and 25b
(`T-ax3r`, #142). Serves `F-0003#42` through `NFR-001`. Base `db99f4c` (the
merge of #141). Author: Claude Opus 5.5 on Claude Code. Started under O-169.
Evidence: [`runs/T-trej/`](../../runs/T-trej/).

## O-173 and the plan

Before the plan, the author counted about five goal classes in row 25 and asked
the Operator first (O-173, comment 6057109037). **O-173 (a)** (the Operator's
answer in the session, recorded in comment 6057190187): split into 25a, the
rules of a run (the lease and fencing, a human decision, copy before read), and
25b, Start and the restart. The Operator asked why the slicing did not split
it; the cause, the author's (one row per package, no goal count at the
slicing, note 12 of the review of #124 declined for a reason about tests), is
the pitfall that this task writes into `guardrails.md` §2. By the same count,
rows 26 (the command, one artifact and its registration) and 27 (the uat) are
one goal each.

The plan (R12, comment 6057214118). Its review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session
in a clone at `db99f4c`, 11 min 44 s; comment 6057420690):
`approve-with-conditions`, the scope of O-173 (a) kept; five conditions (the
order of the two rules of the wait; the acceptance row of `run.md`; the
heartbeat, the takeover and the release under fencing; the owner of the period
of the heartbeat; each part of row 25 in 25a or 25b). The author's answer
(6057421011) took all five and notes 1 to 8 and 10. Budget maximum 1,000 lines
added plus removed over 16 files against `db99f4c`, close-out inside; Cycle cap
1; no panel.

## What was done

1. **`internal/run`** (new): `LeaseRow`, `Store` (read the lease row; write it
   with a records commit and its push, `ErrRefused` on a refusal) and `Clock`;
   `Take` (a released lease at once; else a read every ten seconds, the takeover
   after `3 × H` with no change since the last change seen, checked first, then
   `HeldError` after `3 × H` from the start of the wait; a progress line at each
   read; a refused takeover is a `LostError`); `Fenced` (a refusal: read the
   lease again; another holder is a `LostError`; still held, one more try);
   `Beat`, `Release` and `Heartbeat` (a beat each `H` until the context ends)
   through `Fenced`; `Decision` (a first copy, no App, an approver in the role);
   `Copy` (a new row and its body, the next `seen` for an edit, none for an
   unchanged body; `created` is the last edit, else the creation).
2. **The documents:** `run.md` (the order of the two rules of the wait and a
   refused takeover; one more try of a refused push, and the heartbeat and the
   release under it; a comment edited before its first copy; the acceptance row
   with a stand-in records store); the plan rows 25a and 25b, row 26 After 25b,
   the sentence of O-173; the pitfall in `guardrails.md` §2; the traceability
   rows; the Test cell of `NFR-001` and a §13 line; the backlog lines of 25a and
   25b. [`docs.sh`](../../runs/T-trej/docs.sh) checks them. The issue of 25b
   (#142) lists each part of row 25 that 25a does not take, and what this task
   carries to it.

**Tests:** [`test-runs.md`](../../runs/T-trej/test-runs.md): red 1 and 2,
green after one fix of the stand-in clock, two mutations, each caught;
`docs.sh` red, then green after two fixes of the check.

**The rejected alternatives:** the rules inside 25b's steps (no test with no
`git`); a lock file or a forge lock for the lease (§2 keeps it in the records);
a decision over a forge comment rather than its copy.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `6926ca6`, 8 min 18 s; comment
6057687327): `nothing material in scope`, five notes. Note 3 is applied on the
issue only: the After of #142 now says 22b, 23, 25a, as the plan row. The other
notes are kept as known limits, so the reviewed text is the text that lands:
note 1, the acceptance row "A human decision" names a review, a review comment,
a commit and a reaction, which are never rows of `copies.tsv` and so are not
in the test; note 2, two lines of `run.md` run past the wrap; note 4, the doc
comment of `Take` says a line at each read (the first read and the last give
none), and that of `Heartbeat` says nil when the context ends (a context that
ends inside a write gives that write's error); note 5, the reviewer could not
run mutations and read the record.

## Verdict

Delivered: with a stand-in clock and a stand-in records store, a second run
takes the lease only after the heartbeat has not moved for `3 × lease.H` by its
own clock; fencing, a human decision and copy before read are in
`internal/run` for row 25b. Row 25 is split into 25a and 25b (O-173 a), and the
issue of 25b is #142. The review ended by decay at cycle 0. The diff against
`db99f4c` is inside 1,000 lines over 16 files.

Next: row 25b (`T-ax3r`, #142), whose After cell holds 22b, 23 and 25a.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The reading of row 25; O-173 | reasoning | Claude Opus 5.5 | max | not reported | 09:42 to 09:45 |
| The answer to O-173; the plan | reasoning | Claude Opus 5.5 | max | not reported | 09:45 to 09:48 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,155,533 (USD 6.06) | 11 min 44 s, from 09:48 |
| The answer; the work, test first; the issue of 25b; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 10:00 to 10:08 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 940,739 (USD 4.13) | 8 min 18 s, from 10:08 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 10:17 to 10:20 |

## The body of #130 before the split

> ## Goal
>
> Task `T-trej`: row 25 of the [implementation plan](https://github.com/pharzam/layup/blob/main/docs/plan/README.md#the-tasks-of-m2a), internal/run: Start and the restart. Milestone `M2a`; the specification is [`docs/spec/run.md`](https://github.com/pharzam/layup/blob/main/docs/spec/run.md), [`forge.md`](https://github.com/pharzam/layup/blob/main/docs/spec/forge.md), [`records.md`](https://github.com/pharzam/layup/blob/main/docs/spec/records.md#nfr-001--the-records-of-start) and [`packages.md`](https://github.com/pharzam/layup/blob/main/docs/spec/packages.md#the-table-of-m2a). Fact: `F-0003#42`. Opened by `T-zwke` (#124). Refs #123, Refs #29.
>
> **The demo** (R11): Against the `httptest` forge and a local bare repository, Start makes the first records commit with the pin, the briefs, `approvers.tsv`, `start.tsv` and the lease row.
>
> ## Scope
>
> - **The parts of the specification:** `run.md`: "The steps of `layup run --new`", "The README of a target that Start makes", "The restart", "The lease and fencing", "A human decision", "Copy before read", "The Intake and control issues", rows 6 to 11 of "Input states", the table `run-steps`, NFR-006, NFR-002, REQ-002.
> - **Packages:** `internal/run`. **Requirements:** NFR-001, NFR-002, NFR-006, REQ-002. **Items:** `later-p2-run-start`.
> - **Tests:** unit (the lease with a stand-in clock and `git`; the two rules; the input states); integration (Start, and a restart with another LAYUP version); the block test (`run-steps` to `built`).
> - **Size:** large, about 1200 lines (an estimate; the plan review sets the budget). **After:** 21, 22, 23, 24. **Expected cap:** 1.
>
> ## Acceptance criteria
>
> - [ ] The demo above holds, by the tests of the Tests cell, written red first.
> - [ ] Each part above is built as the specification says; a difference is a defect of the specification, fixed in the same pull request (R10).
> - [ ] `go build ./...`, `go vet ./...`, `go test ./...` and `go test -tags=integration ./...` pass; `adr-lint`, `prd-lint`, `link-lint`, `setup-check` and `run-discipline-tests` exit 0.
> - [ ] Docs updated in the same PR (`docs/tests/traceability.md`, the Test cells of `PRD-0001` §12).
>
> *Written by Claude Opus 5.5, the author of T-zwke (#124).*
>
