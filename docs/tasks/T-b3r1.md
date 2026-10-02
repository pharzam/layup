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

## Review round 1 and its fix (cycle 1)

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record. Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a
fresh session, on `1df7566`, its record at 4 min 31 s; the record is on #90)
gave `material`, with one finding and six notes. It measured `verify.Flagged`
on a clone of LAYUP at `d2516fd`: 38 files, with the baseline's own
onboarding, glossary and guardrails files. The fix has its red runs
([`test-runs.md`](../../runs/T-b3r1/test-runs.md)):

1. **The prose step did not give its one table on the real baseline.** With
   the input of S07, S08 or S09 missing and each input of S14 present, S14
   committed and failed its evidence `adapted` on the baseline's file that the
   other step had not replaced (exit 1). The prose step is now one unit of
   input: while an input of a step of the group that is not done is missing, no
   step of the group copies a file, each step gives the rows of its own missing
   inputs, and a step whose own inputs exist waits, with no row. This changes
   the author's answer to condition 2 of the plan review ("A step of the group
   whose inputs exist is done and committed in the run that stops for another
   step"), as the finding asks. S14's list is README.md, AGENTS.md and each
   other flagged file less the named file of a step of S07 to S09 that is not
   done, so a named file that check `adapted` still flags after its step is
   S14's too, and a new input for it reaches the tree.

The notes: note 2 is fixed (the stop table and the answers record show a
marker with no carriage return at its end; a fill keeps the line end; a `gap`
for a marker with a tab or a carriage return is `fail`, a known limit of
`open-gaps.tsv`); note 3 is text (the ask of S10 has no code span); notes 4 to
7 confirm the change. The author also fixed the stale words "the record of
S06" in the row S11 of `setup.md`, and the count of the known limit (38 files,
row 11's measurement).

## Review round 2

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record again. Round 2 (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh session with another prompt, on `7f2e2ed`, its record at 4 min
13 s; the record is on #90) gave `nothing material in scope`, with six notes.
It read the group loop of the runner against the one unit of the prose step
(no input with which each step of the group waits and none gives a row; S10
does not run while a step of the group is not done), the list of S14, and the
carriage return of note 2; each note confirms the fix. Its note 3 names an
edge of the scanner of row 10 (an unclosed marker before a closed one on one
line is one span to the first close quote), which `check_markers` reads the
same way.

## Verdict

Delivered: steps S05 to S11 and S14 of `layup setup`. S05 removes the history
of the baseline (K12) and copies the input file of each file whose links then
break; S06 copies the briefs; the prose step (S07 to S09 and S14) is one unit
of input with one stop table, and copies the named files and each file that
check `adapted` flags, with a record row for each (K42); S10 asks one question
per marker in one table; S11 fills each marker at the place that the scanner
found, keeps each gap with its row of `open-gaps.tsv`, and writes the second
answers record (O-124). The runner starts each step from the head and makes
no commit when nothing changes; `internal/cli` hands the steps the calls of
`internal/verify` and refuses a brief that holds a marker; check `identity`
reads that `README.md` names the branch `layup-records`.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with three
conditions, applied; round 1 (`1df7566`) gave `material`, one finding, fixed in
cycle 1 with note 2, note 3 and the stale words of the row S11 (`7f2e2ed`);
round 2 (`7f2e2ed`) gave `nothing material in scope`. The records are on #90.
At `7f2e2ed`, `go build`, `go vet` with each tag, `gofmt`, the three test
levels, `go test -race` on six packages (and on four with
`-tags=integration`), `run.sh` and the discipline tests pass; 33 mutations
are detected; at the head, all local checks pass, and `review-record-lint`
passes on the comments of #90 (2 rounds, cap 1). The diff against
`origin/main` is 2,320 lines added plus removed over 32 files with the close-out, inside the Budget maximum of
3,600 lines over 40 files.

Next: row 15 of the plan (`T-d6q5`, #92), whose After cell (rows 13 and 14) is
then merged.

## Resource record

Recorded, not budgeted (ADR-0007). Times are 2026-10-02, UTC. A token count is
the `result` event of the Claude Code CLI (input, output, cache creation and
cache read tokens, and its cost) where that harness gave one; `not reported`
where the harness or the author's session gives none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (D1 to D10), with the measurements on `d2516fd` | reasoning | Claude Opus 5.5 | max | not reported | 22:15 to 22:22 |
| The plan review, first and second harness: skipped (Devin's usage quota; OpenCode no output in five minutes) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 22:22:30 to 22:28:50 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 784,795 (USD 3.34) | 5 min 45 s, from 22:28:57 |
| The answer to the plan review, with D11 and D12 | reasoning | Claude Opus 5.5 | max | not reported | 22:35 to 22:37 |
| The tests first, the code, the integration and e2e tests, the mutations, the documents and the evidence; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 22:37 to 23:06 |
| Review round 1, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 23:06:46 to 23:13:07 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 665,612 (USD 3.46) | 4 min 45 s, 23:13:13 to 23:17:58 |
| The author's self-review while round 1 ran (the stale words of the row S11) | execution | Claude Opus 5.5 | max | not reported | 23:13 to 23:18 |
| The fix of round 1, test first, and the freeze | execution | Claude Opus 5.5 | max | not reported | 23:18 to 23:27 |
| Review round 2, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 23:27:03 to 23:33:21 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,018,846 (USD 2.58) | 4 min 24 s, 23:33:26 to 23:37:50 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 23:38 to 23:42 |
