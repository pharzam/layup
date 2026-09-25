# T-wjq4 — PRD-0001: the product requirements of LAYUP, full scope with phases

Issue: [#63](https://github.com/pharzam/layup/issues/63) (parent
[#42](https://github.com/pharzam/layup/issues/42), the PDR, child 3). Serves
`F-0003#51`. Evidence: [`runs/T-wjq4/test-runs.txt`](../../runs/T-wjq4/test-runs.txt).

## Test runs

Run in the worktree of this task on 2026-09-25, base `9017d99`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | The template copied to `docs/prd/PRD-0001-layup.md`, unfilled | `sh docs/prd/prd-lint.sh` | 1 | `awk: newline in string F-0001 F-0002 F-0003... at source line 1` — the linter itself: it hands `awk -v` the fact list newline-separated, which macOS awk rejects once the facts directory holds two or more records and a PRD exists (revealed on the task's path; fixed here under rule 1, O-48) |
| 2 | A second fixture fact `docs/prd/tests/facts/F-0002-sample.md` added; the old linter | `sh docs/tests/run-discipline-tests.sh` | 1 | `run-discipline-tests: 79 passed, 2 failed` (`prd-lint/good`, `prd-lint/good-path with space`: wanted exit 0, got 1) |
| 3 | The linter fixed: the list space-separated, membership tested with spaces | the same | 0 | `run-discipline-tests: 81 passed, 0 failed` |
| 4 | `PRD-0001` written (18 REQ, 7 NFR, the criteria table, the §12 matrix) | `sh docs/prd/prd-lint.sh` | 0 | `prd-lint: OK` |
| 5 | Copies of the tree with one mutation each: `REQ-014` renamed to `REQ-013`; `REQ-015` (`Won't`) given phase 1; `REQ-014` citing `F-0099`; the `NFR-007` matrix row removed | the same, in each copy | 1 each | `duplicate requirement id REQ-013` and `matrix lists REQ-014 which is not a requirement`; `REQ-015 is Won't but Phase is not —`; `REQ-014 cites no resolvable F-NNNN fact`; `requirement NFR-007 missing from traceability matrix` |
| 6 | The index row, the README note removed, the onboarding and glossary rows; this file not yet written | `sh docs/adr/adr-lint.sh`; `sh docs/links/link-lint.sh`; `sh docs/setup/setup-check.sh`; `git diff --check` | 0 each | `adr-lint: OK`; `link-lint: OK  852 links resolved`; every check `OK` (the `markers` check accepts the removed slug note of the PRD README); nothing |
| 7 | Frozen head `5badccb`, round 1 (cycle 0), first named harness | `devin -p --prompt-file .review-brief/PROMPT.md --model gpt-6-sol-xhigh --permission-mode dangerous --respect-workspace-trust false`, in a disposable clone | 0 | the review record, 549 s: `material`, 5 findings and 4 notes (posted on #63). Finding 1: at that head `sh docs/setup/setup-check.sh` exits 1 — `kit-history` (this file had no backlog line) and `markers` (row 6 quoted the slug marker with its brackets) — while row 6 and the evidence, recorded before this file existed, said every check `OK`. Findings 2–5: four sentences of the PRD (REQ-008, REQ-014, the criteria of REQ-012 and NFR-002) |
| 8 | The fixes of round 1 present: the backlog line, row 6 reworded, the four PRD sentences, the four notes | `sh docs/prd/prd-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `sh docs/setup/setup-check.sh`; `sh docs/links/link-lint.sh`; `git diff --check` | 0 each | `prd-lint: OK`; `81 passed, 0 failed`; every check `OK`; `853 links resolved`; nothing (in `runs/T-wjq4/test-runs.txt`) |
| 9 | Fix head `9ce41cb`, round 2 (cycle 1 of cap 2), first named harness | the same `devin -p …` command, in a new disposable clone | 0 | the review record, 566 s: `material`, 13 findings and 3 notes (posted on #63): one class — an acceptance criterion that can pass while part of its requirement fails — plus two stale matrix fact cells and two statements narrower or wider than their facts (REQ-008's `F-0001#28` exception; REQ-015's boundary) |
| 10 | The fixes of round 2 present: thirteen criteria or statements tightened, the two matrix cells, the glossary list, the §10 link | the five checks of row 8 | see `runs/T-wjq4/test-runs.txt` | recorded there at the fix commit; expected: every check `OK` |
| 11 | Fix head `ec57e91`, round 3 (cycle 2, the cap), first named harness | the same `devin -p …` command, in a new disposable clone | 0 | the review record, 3,294 s: `not mergeable, findings recorded`, 8 findings and 2 notes (posted on #63). The 900 s watchdog (`perl -e 'alarm 900; exec …'`) did not stop the process; the record was accepted as returned, and rule 4's fifteen-minute limit was not applied — a limit of this run, recorded |

## Verdict

Ended at the cycle cap (2) with `not mergeable, findings recorded`; split to the
successor [`T-84r5`](T-84r5.md) ([#64](https://github.com/pharzam/layup/issues/64)),
which carries the branch at `ec57e91` with the eight findings and two notes of
round 3 applied. The Operator ruled the findings material (O-49). Delivered on
the branch: [`PRD-0001`](../prd/PRD-0001-layup.md) (Draft; 18 REQ, 7 NFR, a
criterion each, four phases, the §12 matrix), the linter fix of
`docs/prd/prd-lint.sh` with its fixture (O-48, test first), and the registration.
Round 1 (GPT-6 Sol on Devin): `material`, 5 findings and 4 notes, fixed in
`9ce41cb`. Round 2 (GPT-6 Sol): `material`, 13 findings and 3 notes, fixed in
`ec57e91`. Round 3 (GPT-6 Sol): 8 findings and 2 notes, all of one class — a
criterion that admits a narrower gap — fixed in the successor. Revealed: none
beyond the linter defect fixed here.

## Resource record

Recorded, not budgeted (ADR-0007). Token counts: `not reported` (Devin prints
none; the author's session count was not read).

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The issue, the plan and its review | reasoning | Claude Fable 5.1 on Claude Code (author); the Operator (review) | not reported | — | about 15 min |
| The linter fix (test first) and the PRD | execution | Claude Fable 5.1 — a reasoning-tier model on an execution part: the author's session did the edits | not reported | not reported | about 40 min |
| Round 1 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 549 s |
| The fixes of round 1 | execution | Claude Fable 5.1 | not reported | not reported | about 15 min |
| Round 2 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 566 s |
| The fixes of round 2 | execution | Claude Fable 5.1 | not reported | not reported | about 20 min |
| Round 3 | reasoning | GPT-6 Sol on Devin | xhigh | not reported | 3,294 s |
| **Total** | | | | not reported | about 2 h 45 min |
