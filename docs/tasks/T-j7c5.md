# T-j7c5 — the read-back of the forge check on a parameter that GitHub omits at its default

Issue: [#330](https://github.com/pharzam/layup/issues/330), a defect that run 6 of the demo of `M2d` (row 59b, `T-sdtb`, #284) revealed, off its path, so a task of its own (gate step 5), as `T-gkvi` was.
Serves `F-0003#43` through REQ-003. Base `8802e6c`. Author: Claude Sonnet 5.5 on Claude Code. Evidence: [`runs/T-j7c5/`](../../runs/T-j7c5/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #330. The plan review (GPT-6.1 Sol on Codex CLI, `high`, 3 min 28 s) gave `approve-with-conditions`: 400 lines over 13 files, cap 1, no panel. Its four
conditions are applied: "red first" read as the changed behavior red against the base at the four levels and the rest as guards proven by mutants; an assertion on each fake's own answer, with a mutant that restores the
echo; the input-state and e2e edits of the specification named; the source claims corrected (a new case, the times, the first-hand answers of the forge). The plan is reasoning-tier work that an execution-tier model
wrote under the Operator's routing message of 2026-10-13: a routing mismatch by choice. The goal count is 1.

## What was done

The first forge check of the demo on the real forge passed seven rows of eight, the probe among them, and failed `rules:layup-records` with "update: update_allows_fetch_and_merge is missing": the file wrote the parameter `false`,
and GitHub returns the rule `update` with no parameters in the response of the apply, in the read of the ruleset and in the effective rules, while it echoes the other `false` parameters. Every fake of every test echoed the
file, so no test met the form.

1. **The comparison** (`internal/run/forgecheck.go`): a table `omittedAtDefault`, by rule type and parameter name, with one entry (`update`, `update_allows_fetch_and_merge`, `false`); a parameter of the file that the call lacks is
   equal when it is in the table and the file holds the default. An omitted `true`, another missing parameter, the same name on another type, a value that differs and a missing rule stay `fail`.
2. **The fakes** of the job, the world and the e2e give the rules of the records branch as GitHub does, each with an assertion on its own answer; the world has a case where the forge returns `true`.
3. **The documents:** `scaffold.md` (The read-back with a paragraph "Decided here", the comparison, the input states, the acceptance rows, stage 4), `forge.md`, the pitfall of `guardrails.md` §2, `traceability.md`, the cells of REQ-003.

**Tests:** [`test-runs.md`](../../runs/T-j7c5/test-runs.md): the first-hand answers, the red runs at the four levels, 10 mutants each alone on a copy (all caught), and the gate.
**Known limits:** the omission is observed on one repository at one time, and a fork may return the field; an omitted `false` cannot be told from a parameter that the forge never held (the apply output is the record that the file
was accepted); a new entry of the table needs its own first-hand answers.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`aad9829`, cycle 0 of the cap of 1; comment 6188450649): `material`, two findings in the change and a note. The job case that the answer to the plan review promised (the file wrote `true`, the forge returns `false`) was
missing from the fail table; and the record claimed that the mutants prove the guards of a returned value, though none changed that comparison. The note asked for the limit of the omitted default in the specification. The reply is
comment 6188467674: the case (a guard, run at 08:57:38Z), two mutants (12 in all, all caught), the limit in the specification.

Round 2 (`d844ae6`, cycle 1; comment 6188593930): `nothing material in scope`, no finding and no note.

## Verdict

Delivered: the read-back of the forge check passes on GitHub's real form of the records ruleset's rule `update` (no `update_allows_fetch_and_merge` for a file that wrote `false`), by a table of one entry, with the fakes of the job, the
world and the e2e in that form, each asserted, and the documents. The review ended at `nothing material in scope` after two rounds, at the cap of 1. Next: the demo of `M2d` (`T-sdtb`, #284), restarted with the binary of this merge.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-15, from the clock and the comments of #330; the tokens of GPT-6.1 Sol are the "tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (a routing mismatch) | Claude Sonnet 5.5 | `max` | not reported | 07:59 to 08:02 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 144,625 | 3 min 28 s, 08:02:12 |
| The answer, the tests, the code, the mutants, the documents, the gate | execution | Claude Sonnet 5.5 | `max` | not reported | 08:06 to 08:47 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 180,016 | 9 min 21 s, 08:47:35 |
| The fix of round 1, its mutants, the checks | execution | Claude Sonnet 5.5 | `max` | not reported | 08:56 to 09:04 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 153,934 | 6 min 36 s, 09:04:00 |
| The close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 09:11 to 09:12 |
