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
is stale in the target after S04; and, from round 1, the warning of
`adr-lint.sh` for the new record, and a baseline whose `facts.sha256` lists a
file that S04 changes.

## Review round 1 and its fix (cycle 1)

Devin (its usage quota) and OpenCode ("Go usage limit exceeded") gave no
record. Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a
fresh session, on `5dadab0`, its record at 7 min 44 s; the record is on #86)
gave `material`, with one finding and eight notes. Its runs used the built
binary on a clone of LAYUP at its root commit `d2516fd` as the baseline (the
record `0009-pin-the-baseline.md`, `F-0001-setup-answers.md`, the two index
rows, `adr-lint.sh` exit 0). The fix has its red runs
([`test-runs.md`](../../runs/T-7s0y/test-runs.md)):

1. **The `source` of an answer kept its angle quotes** in the answers record,
   against `setup.md` ("each angle quote of a recorded text is written as
   `&lsaquo;` or `&rsaquo;`"), so check `markers` of the target would read a
   marker there. S04 now writes them as entities, as for the answer and the
   question.

Three defects that the author found in a self-review while the round ran are
in the same fix: a table of an index file at the end of the file with no line
feed got the new row on its last line; a list of hashes of the root commit
with no line feed at its end got the new line on its last line (note 2 (b));
a target on `main` whose branch `layup-setup` exists got a reason that named
`main` as a branch it is not on. The notes: 2 (a) (a baseline whose
`facts.sha256` lists an index file that S04 changes) and 5 (the warning of
`adr-lint.sh` for the new record, which no document links) are known limits
in `setup.md`; 3 is text (`setup.md` gives the order of the checks of a done
step's inputs, as the runner does them); 4 (`layup setup verify` on a setup
that is not finished gives exit 2, by its own rule) and 6 to 9 confirm the
change.

## Review round 2

Devin (its usage quota) gave no record, and OpenCode (Grok 4.7) gave no output
in five minutes (a smoke run of 120 s, then one of 300 s), so both are skipped
(Bootstrap mode rule 4). Round 2 (Claude Fable 5.1, effort `xhigh`, on the
Claude Code CLI, a fresh session with another prompt, on `8761a15`, its record
at 5 min 31 s; the record is on #86) gave `nothing material in scope`, with six
notes. It ran the binary on a clone of LAYUP at `d2516fd` as the baseline, with
an answer and a `source` that hold angle quotes (the record holds the entities
and no raw quote), and on baselines whose index files and list of hashes end
with no line feed, with a blank last line, and with CRLF lines. Its note 1 (a
list of hashes with CRLF lines fails the evidence of S04, as before the fix) is
a known limit in `setup.md` at the close-out; notes 2 to 6 confirm the fix.

## Verdict

Delivered: steps S01 to S04 of `layup setup`. S01 asks the four fixed
questions and the gaps of the problem statement in one stop table, checks the
answers, and records them with their sources; `internal/cli` hands it the gap
table of `internal/psb` (`psb-batch-api`), and a problem statement that changes
after S01 is exit 2. S02 resolves the baseline's commit once and copies it with
no `.git`; S03 makes the root commit with the tree `pin.tree` and writes the
remote and the push for the Operator; S04 writes, on `layup-setup`, the pin
file, the decision record of the pin and the record of the `S01-` and `Q-`
answers (O-124), with its evidence, checks `pin` and `facts`. The runner has the
outcome `input` and the evidence call of a step with its undo. A step that
stopped in its middle goes on, and a state that `layup` did not make is exit 2.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with three
conditions, applied; round 1 (`5dadab0`) gave `material`, one finding, fixed in
cycle 1 with three defects that the author found and notes 2, 3 and 5 as text
(`8761a15`); round 2 (`8761a15`) gave `nothing material in scope`, and its note
1 is a known limit. The records are on #86. At `8761a15`, `go build`, `go vet`
with each tag, `gofmt`, the three test levels, `go test -race` on six packages
(and on four with `-tags=integration`), `run.sh` and the discipline tests pass;
30 mutations are detected; at the head, all local checks pass, and
`review-record-lint` passes on the comments of #86 (2 rounds, cap 1). The diff
against `origin/main` is 2,620 lines added plus removed over 26 files with the close-out, inside the Budget maximum
of 3,000 lines over 40 files.

Next: row 13 of the plan (`T-b3r1`, #90), whose After cell (rows 9 to 12) is
then merged.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The first reading of row 9, and the question of the order (O-137) on #86 | reasoning | Claude Opus 5.5 | max | not reported | to 17:29 (its start not recorded) |
| The plan (D1 to D14) | reasoning | Claude Opus 5.5 | max | not reported | 20:45 to 20:51 |
| The plan review, first and second harness: skipped (Devin's usage quota; OpenCode "Go usage limit exceeded") | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 20:51:20 to 20:52:33 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 584,772 (USD 3.07) | 5 min 6 s, from 20:52:39 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 20:57 to 20:58 |
| The tests first, the code, the integration and e2e tests, the mutations, the documents and the evidence; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 20:58 to 21:40 |
| Review round 1, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 21:40:46 to 21:42:10 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,300,245 (USD 4.13) | 8 min 5 s, 21:42:21 to 21:50:26 |
| The author's self-review and the fixes of three defects, while round 1 ran | execution | Claude Opus 5.5 | max | not reported | 21:43 to 21:47 |
| The fix of round 1, test first, and the freeze | execution | Claude Opus 5.5 | max | not reported | 21:50 to 21:55 |
| Review round 2, first and second harness: skipped (Devin's usage quota; OpenCode no output in five minutes) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 21:55:56 to 22:04:26 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 995,112 (USD 3.45) | 6 min 0 s, 22:04:32 to 22:10:32 |
| The close-out, with note 1 of round 2 | execution | Claude Opus 5.5 | max | not reported | 22:10 to 22:16 |
