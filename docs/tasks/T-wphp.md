# T-wphp — row 57a of the plan, the step forge-check

Issue: [#280](https://github.com/pharzam/layup/issues/280), row 57a of [the tasks
of M2d](../plan/README.md#the-tasks-of-m2d), opened by `T-bng8` (#267). Serves
`F-0003#43` and `F-0003#64` through REQ-003 and NFR-001. Base `be28498`. Author: Claude
Sonnet 5.5 on Claude Code. Evidence: [`runs/T-wphp/`](../../runs/T-wphp/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #280. The plan review
(GPT-6.1 Sol on Codex CLI, `high`, 4 min 35 s) gave `approve-with-conditions`: 1,400 lines
over 24 files against `be28498`, close-out inside; cap 1; no panel. Its four conditions are
applied: the comparison of the rules by one returned rule of the type (an R10 reading in
`scaffold.md` before the tests), what ends the step with no table, the boundary tests of the
work area, and the slices. The plan was written by the session's model, an execution-tier one:
a choice, not an unavailability of the reasoning tier.

## What was done

The step `forge-check` of `layup run` (`forgecheck.go`), after `scaffold-push` when that step is
done and not waiting: the eight checks of `setup/forge-check.tsv`, each with its row whatever an
earlier one gave (the head against `setup/tree.tsv`; the rules of three branches read back against
the two ruleset files of the work area; the bypass lists of `setup/rulesets.txt`; the probe, a pass
only for a refusal with `GH013` and the ruleset's name; the activity), and the table in one records
commit through fencing in either outcome. `forgeCheckJob` (the seams), `readBack`, `readApplied`,
`probeVerdict` and `pushProbe` are its parts. The specification, the traceability, the PRD cells and
`guardrails.md` §2 (three pitfalls) are updated.

**Tests:** [`test-runs.md`](../../runs/T-wphp/test-runs.md): red first for the tests of the Tests cell;
110 mutants, each alone on a copy of the tree, all caught in the end; the whole gate on the head, twice.

**Known limits:** the refusal text that the pass needs is that of the tests' stand-in, and no official page of
GitHub confirms it: the demo of row 59b (`T-sdtb`) reads it (a pitfall says what to do). The step runs at the
tail of the restart; its place at each later start is row 57b (`T-2j13`). A probe that the forge accepts leaves
`layup-probe` for the Operator to delete. The scenarios of conditions 2 and 3 and some exact-reason cases were
written after the code, each red against the mutant of its rule (round 2, finding 1). A restart after a failed
rebuild uses the half-built work area (#318, from `scaffold-push`). CI's integration run of `internal/run` took
282 s against its limit of 600 s for a package (645 s here): row 58 (`T-eevt`) adds a world and should watch it.

## Review rounds

Both rounds: GPT-6.1 Sol on Codex CLI (`xhigh`; ADR-0026 does not name it, the Operator allowed it).

Round 1 (`e7e818c`, cycle 0; comment 6154043779): `material`, three findings and a note: a required check with no
`integration_id` in the answer was equal to one whose file says 0 (fixed test first, `3adb2b8`); the cases written after
the code; stale test names and a comment (fixed); the sandbox has no loopback socket. The CI of the first push had failed
in the probe test (the `git` of CI makes no remote-tracking ref for an accepted push): fixed with them. The answer is
comment 6154934153; #318 was opened for the defect that the change revealed.

Round 2 (`b94b1c6`, cycle 1 of the cap of 1; comment 6155102401): `not mergeable, findings recorded`: finding 1, the
order of the tests, is still material at the cap. Its notes: the wording of run 8 (applied), the sandbox, and #318, which
the reviewer agrees is revealed. The author recommended a known limit; the Operator chose one more cycle (Telegram message
227), and the cap was raised to 2 (comment 6160511807).

Round 3 (`3613865`, cycle 2 of the cap of 2; comment 6160637217): `not mergeable, findings recorded`, the same finding: the
reviewer says that no author-only fix inside a cycle settles it, only the Operator's disposition at the cap.

## Verdict

Delivered: the step `forge-check` of `layup run`, with its tests, documents and runs. The last round ended `not mergeable,
findings recorded` at the cap of 2 on one finding, the order of the tests that conditions 2 and 3 of the plan review named
(finding 1 of rounds 2 and 3), which the Operator ruled a known limit (comment 6161941436 of #280, 2026-10-14: end 2 of the cap);
its other findings are notes, and each finding of round 1 is settled. The diff against `be28498` is 1,374 lines over 16 files with
the close-out inside, against the plan review's 1,400 over 24. Next: `T-2j13` (the forge check at each later start, row 57b).

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-13 unless a row says 2026-10-14, from the clock; the tokens of GPT-6.1 Sol are the
"tokens used" of Codex CLI; the author's are not reported.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning (not reached) | Claude Sonnet 5.5 | `max` | not reported | 18:23 to 18:27 |
| The plan review | reasoning | GPT-6.1 Sol, Codex CLI | `high` | 182,139 | 4 min 35 s, 18:27:44 |
| The tests, the code, the mutations, the documents | execution | Claude Sonnet 5.5 | `max` | not reported | 18:35 to 19:46 |
| Round 1 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 207,174 | 9 min 37 s, 19:46:55 |
| The fix of round 1, the mutants, the gate | execution | Claude Sonnet 5.5 | `max` | not reported | 19:57 to 20:50 |
| Round 2 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 232,861 | 9 min 0 s, 20:50:57 |
| The cap raised to 2, the record of the choice | execution | Claude Sonnet 5.5 | `max` | not reported | 03:49 to 03:55 (2026-10-14) |
| Round 3 | reasoning | GPT-6.1 Sol, Codex CLI | `xhigh` | 186,126 | 10 min 12 s, 03:54:51 (2026-10-14) |
| The notes and the close-out | execution | Claude Sonnet 5.5 | `max` | not reported | 21:02 (2026-10-13), and 05:56 to 05:58 (2026-10-14) |
