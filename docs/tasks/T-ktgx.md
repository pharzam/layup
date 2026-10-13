# T-ktgx — row 52a of the plan, the setup in phase 2

Issue: [#273](https://github.com/pharzam/layup/issues/273), row 52a of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#42` and `F-0001#8` through REQ-002, NFR-006 and NFR-003. Base `32bae7e`
(the merge of #302); the budget is read against it. Author: Claude Sonnet 5.5 on
Claude Code. Evidence: [`runs/T-ktgx/`](../../runs/T-ktgx/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #273. The plan
review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`32bae7e`, 9 min 46 s) gave `approve-with-conditions`: Budget maximum 1,200 lines
added plus removed over 18 files, close-out inside; Cycle cap 1; no panel. Its three
conditions are applied: `StepsPhase2` sets the actor `layup-run` on S02, S03 and S15
(`scaffold.md`, the paragraph after the table of phase 2), each rule that the plan
adds beyond the words of a row has a sentence in its row (S01 and S04; the form
checks of S02 were dropped, as `CheckStart` gives them), and S04 keeps its evidence
`pin` and `facts`. The goal count is 1, final by O-209 of #267.

## What was done

1. **The values:** `Phase2` in `internal/setup/phase2.go` holds what `layup run` gives
   in Go, in place of a person, the forge and the baseline's repository: the target's
   `OWNER/NAME`, `forge.visibility`, the four pin rows of `start.tsv`, `psb.sha256` and
   the questions of the batch.
2. **The steps:** `StepsPhase2` takes the steps of phase 1 and replaces five. S01 asks
   only `S01-stack` (it runs no `layup psb check`, and the `Q-` rows are taken as
   given), gives `name`, `visibility`, `baseline` and `brief.sha256` from `Phase2` and
   applies the checks of phase 1 to them. S02 resolves nothing and gives the four
   `pin.*` rows with the ref `records:start/start.tsv`. S03 makes no commit and no
   command, and checks that the tree of `main` is `pin.tree`. S04 is phase 1's body,
   shared as `runS04With`, with the questions that its raw fact records given by the
   phase: `S01-stack` and each question of the batch with the `question` of its row of
   `questions.tsv`; a `Q-` row that the batch does not hold is a `fail`. S15 is phase
   1's with no records commit and no command. The actor of S02, S03 and S15 is
   `layup-run`; S05 to S14 are phase 1's.
3. **The documents:** the rows S01, S02 and S04 of the table of phase 2 in
   `scaffold.md` (the rules that the plan reads beyond a row, and the evidence lines),
   the rule of `ref` and the gloss of `computed` in the block `setup-record`
   (`setup.md`), the Job cell of `internal/setup` in `packages.md`, the package
   comment, the rows of `traceability.md`, the Task and Test cells of REQ-002, NFR-003
   and NFR-006 and a §13 line of `PRD-0001`, and a sentence of the pitfall "Tests that
   pass for the wrong reason" in `guardrails.md` §2. One case of
   `TestTheSourcesOfARecord` in `internal/verify` guards the claim of `records.md`
   (lines 590-594) that a `computed` row may cite `records:start/start.tsv`; it passes
   at the base, so it has no red.

**Tests:** [`test-runs.md`](../../runs/T-ktgx/test-runs.md): red before the code
(`undefined: Phase2`); each rule of `phase2.go` and the two edges of the shared body of
`runS04` mutated alone, 36 of 36 caught (a first run of 36 left 1 alive, the refusal of
a gap table of another form, which the refactor moved into the shared body and no test
of phase 1 pinned; `TestS04OfPhase1RefusesABadGapTable` kills it); every check of the
gate exits 0 on the code, and again after the notes, at 06:56Z to 07:03Z (`TMPDIR` on the disk, so that `/tmp` is not the limit).

**Known limits:** a stray `S01-name`, `S01-visibility` or `S01-baseline` row of
`answers.tsv` is no input error: `checkRows` still takes the four `S01-` names before
any step, S01's `Reads` stay `S01-` and `Q-`, and S01 hashes those rows into
`answers.sha256`. Phase 2 reads none of the three, and only check `facts` finds the row
(each question of `answers.tsv` is a fact of one answers record); `layup run` writes
none. The `name` row cites `records:start/start.tsv`, a file that holds no `OWNER/NAME`
(`start.tsv` has no name row; the specification gives the ref twice, and check
`sources` takes any ref that is not empty): a stronger ref is a change of the
specification. `forge.visibility` and `pin.source` take any value in `records.CheckStart`,
so the checks of S01 for the visibility and the baseline are the first line and not a
second; the forms of the pin rows are `CheckStart`'s at each restart, and S02 checks
none. No code calls `StepsPhase2` yet: its caller is row 54 (`T-eahr`), the rulesets and
the commands of S13 are rows 52b and 52c, and the commit of `layup-records` is row 55's.

**The rejected alternatives:** a flag or a branch of phase 2 inside each step of phase 1
(a step cannot give up its `Commands` or `Records` by an input, and the tests of phase 2
would spread over the steps of phase 1); the inputs as a file of the work area
(`scaffold.md` rejects it); a copy of `Steps` for phase 2 (it repeats the wiring of S05
to S14); the binary as a program; form checks in S02 (`CheckStart` gives them).

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a clone at
`2017853`, cycle 0, the cap of 1; its record is comment 6139995839 on #273): `nothing
material in scope`, five notes: note 1 is not in the change, notes 2 to 5 are. The
`tests` job of CI passed on that head.

- **Applied** in `7169ac3`, a commit that no round read (two comments and documents, no
  code): note 2 (one evidence clause in each of the rows S01 and S02 of the table of
  phase 2), note 3 (the gloss of `computed` in `setup-record` names the values that
  `layup run` takes from the records of Start), note 4 (a sentence in the pitfall
  "Tests that pass for the wrong reason" of `guardrails.md`: a refusal that a refactor
  moves into a body that two callers share needs a case of the old caller) and note 5
  (the comments of `setup.go` and `phase2.go` name row 54, `T-eahr`, as the caller).
- **No change** for note 1: the integration level could not run in the review's clone
  because `/tmp` of this host, a 1.9 GB tmpfs, was full of the clones of earlier
  reviews. The clones of merged tasks were removed before 06:49Z, and every check of the
  gate passes again, the integration level included (06:56Z to 07:03Z). CI is the
  authority.

## Verdict

Delivered: the setup in phase 2 of `M2d`. `Phase2` and `StepsPhase2` of `internal/setup`
give the five steps that differ from phase 1: S01 asks `S01-stack` alone, S02 takes the
pin rows of `start.tsv` and resolves nothing, S03 makes no commit and checks the tree of
`main`, S04 writes the pin file and the raw fact of the batch, and S15 makes no records
commit. The tests were red first, each rule was mutated alone and caught (36 of 36), and
every check of the gate exits 0. The review ended by decay at cycle 0 of cap 1. The diff
against `32bae7e` is inside 1,200 lines over 18 files (798 lines added plus removed over 16 files,
measured with the close-out before the commit). Next: row 52b, `T-9v92`, writes the two rulesets of S13
on these steps, and row 54, `T-eahr`, is their first caller.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13; reviewers' tokens are `modelUsage`;
the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Sonnet 5.5 | `max` | not reported | 05:44 to 05:55 (the claim to the post) |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,112,062 (USD 6.34) | 9 min 46 s, 05:55 to 06:05 |
| The tests and the code, test first | execution | Claude Sonnet 5.5 | `max` | not reported | 06:08 to 06:33 (the answer to the pull request) |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,200,175 (USD 6.05) | 11 min 1 s, 06:33 to 06:44 |
| The notes | execution | Claude Sonnet 5.5 | `max` | not reported | 06:50 to 06:52 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 06:53, to the commit of this file |
