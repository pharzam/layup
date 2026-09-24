# T-q1x6 — store the Operator's routing policy as raw fact F-0005

Issue: [#47](https://github.com/pharzam/layup/issues/47) (parent [#46](https://github.com/pharzam/layup/issues/46), ADR-0012).

## Test runs

Run in the worktree of this task on 2026-09-24, base `0a20ee3`.

| # | State of the tree | Command | Exit | Output |
| - | ----------------- | ------- | ---- | ------ |
| 1 | Hash line added to `docs/setup/facts.sha256`; no source file | `sh docs/setup/setup-check.sh --only facts .` | 1 | `setup-check: facts FAIL hash: docs/facts/operator-routing-policy.md does not match docs/setup/facts.sha256` |
| 2 | Source file stored; the index loop names `F-0005`; no index row | the same | 1 | `setup-check: facts FAIL index: docs/facts/README.md has no row for F-0005` |
| 3 | Index row added | the same | 0 | `setup-check: facts OK` |
| 4 | A copy of the tree, with one byte of `docs/facts/operator-routing-policy.md` changed (line 84, `.` to `!`; `cmp -l` counts 1 byte) | the same, in the copy | 1 | `setup-check: facts FAIL hash: docs/facts/operator-routing-policy.md does not match docs/setup/facts.sha256` |
| 5 | Index loop extended, before the fixture change | `sh docs/setup/tests/run.sh` | 1 | `setup-check tests: 29 passed, 1 failed` (`frame/good-passthrough`: `facts FAIL index: … no row for F-0005`) |
| 6 | Stand-in row `F-0005` added to that fixture's index | the same | 0 | `setup-check tests: 30 passed, 0 failed` |
| 7 | Commits `5db2fc0` to `bc3974f`, before this file existed | `sh docs/setup/setup-check.sh` | 0 | every check `OK`; `run-discipline-tests: 81 passed, 0 failed`; `link-lint: OK` |
| 8 | Round 1 fix: the floor loop names `docs/facts/operator-routing-policy.md`; in a copy of the tree, its line removed from `facts.sha256` | `sh docs/setup/setup-check.sh --only facts .`, in the copy | 1 | `setup-check: facts FAIL listed: docs/facts/operator-routing-policy.md is not in docs/setup/facts.sha256` (before the fix the same removal gave `facts OK`, exit 0: round 1, check C9) |
| 9 | The floor extended, before the fixture change | `sh docs/setup/tests/run.sh` | 1 | `setup-check tests: 29 passed, 1 failed` (`frame/good-passthrough`) |
| 10 | A stand-in `operator-routing-policy.md` and its hash line added to that fixture | the same | 0 | `setup-check tests: 30 passed, 0 failed` |
| 11 | The round 1 fix head, with this file and no completed-log line yet | `sh docs/setup/setup-check.sh` | 1 | every check `OK` except `setup-check: kit-history FAIL orphan: docs/tasks/T-q1x6.md has no line with T-q1x6 in backlog.md or completed.md`, which the close-out commit ends (the completed-log line cannot exist before the rounds finish) |

## Verdict

Delivered: [`F-0005`](../facts/F-0005-operator-routing-policy.md), the Operator's
routing policy stored byte-identical (84 lines, SHA-256 `d5814496…da1fb1`) with
the Operator's attestation, as input for ADRs. Check `facts` now reads its index
row and lists its source in the floor, so a changed byte, a removed file, a
removed hash line and a removed index row each fail (rows 1–4 and 8). Round 1:
material (the floor did not list the source; row 7 named no tree), fixed in
`2981b01`. Round 2: nothing material in scope. Revealed: #48 (line endings on an
`autocrlf` checkout).

## Resource record

Recorded, not budgeted (ADR-0007). Author tokens: `not reported` (Claude Code
`/usage` was not read). AGY prints no token count.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5 (plan); Gemini 3.1 Pro on AGY (review) | xhigh (author); high (review) | not reported | 10 min (review: 239 s, and a 21 s run that produced no output) |
| The decay review rounds | reasoning | Claude Fable 5.1 (round 1); Gemini 3.1 Pro on AGY (round 2) | not reported; high | 105,752 (round 1); not reported (round 2) | 15 min (528 s + 391 s) |
| Writing the tests and the record | execution | Claude Opus 5.5 — a reasoning-tier model on an execution part: the author's session did the edits | xhigh | not reported | 8 min |
| Isolate, guardrails, docs, close-out | `—` | Claude Opus 5.5 | xhigh | not reported | 5 min |
| **Total** | | | | 105,752 reported + not reported | about 38 min |
