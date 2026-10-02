# T-7s0y — steps S01 to S04 of `layup setup`

Issue: [#86](https://github.com/pharzam/layup/issues/86), row 9 of the
[implementation plan](../plan/README.md), after row 14 by O-137 (#86); child of
`layup setup` (#77). Serves `F-0003#42`. Base `94d1d71` (the merge of row 14,
#112). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-7s0y/`](../../runs/T-7s0y/).

## Plan and plan review

The author's question of the order (row 14 before row 9), the Operator's
decision O-137, the plan (R12), its review and the author's answer are
comments on #86. The reviewer order is the Operator's (Devin, then OpenCode,
then Claude). Devin (its usage quota) and OpenCode (Grok 4.7: "Go usage limit
exceeded") gave no record. The plan review (Claude Fable 5.1, effort `xhigh`,
on the Claude Code CLI, a fresh session) gave `approve-with-conditions`:
Budget maximum 3,000 lines added plus removed over 40 files against the base,
close-out inside; Cycle cap 1. Its three conditions are applied: S04 takes a
branch `layup-setup` at the root commit, and the undo runs only from the
commit that the step made (1); each value of the header of the answers record,
with the date of `pin.time` (2); the integration test of the hand-off of the
gap table (3). Its notes 1 and 3 to 7 are applied (check `identity` reads the
whole `OWNER/NAME`; the rows S01 and S06 name the check of the problem
statement; the fixed texts are in `setup.md`, and the stale prose of the
baseline's ADR index is a known limit; S04 writes the rows `pin.adr`,
`answers.record` and `answers.record.sha256`; GitHub is named as the forge
once; another visibility is a `fail`); notes 2 and 8 to 10 confirm the plan.

## What was done

1. **S01** (`internal/setup`): the four fixed questions (`internal/work`), the
   `Q-` questions of the gap table, one stop table for each missing answer,
   an input error for an answer to no gap, the checks of the answers (the
   catalog entry of the stack, `OWNER/NAME`, `public` or `private`, a URL that
   `internal/git` reads with no login), and the record rows with the source of
   D3 and `brief.sha256`.
2. **The hand-off** (`internal/cli`): the problem statement read once per run,
   the gap table of `internal/psb` and the SHA-256 handed to the steps, and
   read in `internal/setup` by its own Go value of the block `psb-gaps`.
3. **S02 and S03:** the commit resolved once, the clone into `target.part`
   renamed at the end, `.git` removed, the root commit by the identity of the
   run at `pin.time`, and its own commit of a run that stopped taken (its tree
   and its message); the remote `origin` and the push in `commands.sh`.
4. **S04:** on `layup-setup` from the root commit, the pin file, the decision
   record of the pin with its index row, and the answers record with its index
   row and its line of `facts.sha256`, each from the files of the root commit
   and the date of `pin.time`; three record rows.
5. **The runner:** the outcome `input` (exit 2), the check of a done step's
   inputs before any step (the problem statement, D6), and the evidence call
   of a step after its commit, with the undo `git reset --soft` only from the
   commit that the step made; two calls of `internal/git` (`ResetSoft`,
   `Message`).
6. **The tests:** the unit tests (a fake of `internal/git` and a fixed clock
   for S02 to S04); the integration tests on a stand-in baseline by its
   `file://` URL with real `git` (the pin resolved once, no `.git` of the
   clone, a login URL with no prompt, a step that stopped in its middle, the
   undo, the hand-off against `layup psb check`, checks `pin` and `facts`, the
   baseline's own `adr-lint.sh`); the e2e test of the demo through the binary;
   26 mutations, each detected.
7. **The documents:** `setup.md` (the runner's rules, the rows S01 to S06 with
   the reading of O-124, K14, the rules of S01 to S04 with the fixed texts of
   the two records, check `identity`), `packages.md`, the After cell of row 9
   (O-137), the traceability rows (the planned row named by its tests), the
   `PRD-0001` cells with a §13 row, two lessons in `guardrails.md` §2,
   `README.md` and the onboarding document (two stale sentences of the
   command, and of the checks of `layup setup verify`).

**The rejected alternatives:** `S01-name` as the name alone, with the owner
from another input (S01 has four fixed questions, and S12 and S13 need the
owner); `internal/setup` that imports `internal/psb` (the package rule of
`packages.md`); a runner that calls `internal/verify` itself (a step must not
depend on the checks that judge it); an undo by `git reset --hard` or by a new
commit (the first loses the step's files, the second leaves a commit with no
`done` row); a clone straight into `WORK/target` (a stop in its middle leaves
a directory that S03 cannot tell from a whole copy); a pin time from the
commit's date (the plan's D7 and §5 Start 2: the time of the resolve).

**Known limits** (in `setup.md`): a version of `layup` with other rules of
`layup psb check` gives another gap table for the same problem statement; the
prose of the baseline's ADR index ("the next constitutional ADR is `0009`")
is stale in the target after S04.
