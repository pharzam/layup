# T-b3r1 — steps S05 to S11 and S14 of `layup setup`

Issue: [#90](https://github.com/pharzam/layup/issues/90), row 13 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77). Serves
`F-0003#42`. Base `b812943` (the merge of row 9, #113). Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-b3r1/`](../../runs/T-b3r1/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #90. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh session) gave `approve-with-conditions`: Budget maximum 3,600
lines added plus removed over 40 files against the base, close-out inside;
Cycle cap 1. Its three conditions are applied: a rerun gives the same files
and rows, as the runner puts the target back to its head before each step
(D11) and S11 checks before it writes (1); the list of S14 has no file of S07
to S09, so the group's one table has no question twice (2); the fixed texts of
the second answers record (3). Its notes 1 to 4 and 6 are applied (the counts
and the blockquote rule in the text of S05; the refusal of a brief with a
marker in the runner's inputs, with a known limit; the split rows of S06; the
records branch as a line of check `identity`, not a rule of S14; the union of
the input files in the known limit); notes 5 and 7 confirm the plan. The
author added D12 (no commit with no change).

## What was done

1. **S05:** the history of the baseline removed from the head of
   `layup-setup` (the two directories, each task file, each line of the task
   indexes that names a deleted task or links the baseline's repository, and
   each blockquote that holds one: K12), and the input file of each file whose
   links then break, with one stop table for the missing ones.
2. **S06:** each brief copied byte for byte into `docs/facts/`, with its line
   of `facts.sha256`, its index row and its two record rows.
3. **The prose step:** S07, S08 and S09 copy their named files, and S14
   `README.md`, `AGENTS.md` and each other file that check `adapted` flags;
   one stop table for the group; a record row for each copied file (K42); the
   evidence of S14 is checks `identity` and `adapted`.
4. **S10 and S11:** one question per file and marker, with its first line, in
   one stop table; then each marker filled at the places that the scanner
   found, each gap kept with its row of `open-gaps.tsv`, one record row per
   place, and the second answers record (O-124).
5. **The runner:** the reset of the target before each step from S04 to S14
   (D11, `git reset --hard`, the new call `ResetHard`) and no commit with no
   change (D12, the new call `Staged`).
6. **`internal/verify` and `internal/cli`:** the column of a marker, the link
   rule of check `kit-history` as a call, the markers of a text; check
   `identity` reads that `README.md` names the branch `layup-records` (§3);
   `internal/cli` gives the steps the calls of `internal/verify` and refuses a
   brief that holds a marker; the stand-in's `README.md` names the branch.
7. **The tests:** the unit tests of each step (a fake of `internal/git` and of
   the files), of the runner, of the calls and of the briefs; the integration
   tests on a stand-in baseline with real `git`, the baseline's own
   `link-lint.sh` and each step's checks as its evidence; the e2e test of the
   demo through the binary; 27 mutations, each detected.
8. **The documents:** `setup.md` (the runner's rules, the inputs, the rows
   S05 to S11 and S14, the rules of row 13 with the second answers record and
   the known limits, check `identity`), `packages.md`, the traceability rows,
   the `PRD-0001` cells with a §13 row, a lesson in `guardrails.md` §2,
   `README.md` and the onboarding document.

**The rejected alternatives:** S05 removes only "their lines", the lines of a
deleted task (check `kit-history` then fails on the real baseline); a link
parser of the engine (the baseline's `link-lint.sh` is the rule, through the
call of row 10); a second scanner of markers in `internal/setup` (the plan
keeps one, through `internal/cli`); S11 fills by a search of the text (it
would fill a marker in a place that the scanner does not read); a step that
writes the files of its parent again (D11 puts the target back once, before
each step); a fill of a marker in a raw brief (a raw fact never changes); a
README rule in S14 alone (check `identity` holds it, so `layup setup verify`
reads it later).

**Known limits** (in `setup.md`): a brief that uses the single angle quotation
marks as quotes is refused; a real setup at LAYUP's pin needs the union of the
files that check `adapted` flags, the files whose links S05 breaks, and the
five named files as input files.
