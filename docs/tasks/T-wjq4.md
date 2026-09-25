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
| 8 | The fixes of round 1 present: the backlog line, row 6 reworded, the four PRD sentences, the four notes | `sh docs/prd/prd-lint.sh`; `sh docs/tests/run-discipline-tests.sh`; `sh docs/setup/setup-check.sh`; `sh docs/links/link-lint.sh`; `git diff --check` | see `runs/T-wjq4/test-runs.txt` | recorded there at the fix commit; expected: every check `OK` |
