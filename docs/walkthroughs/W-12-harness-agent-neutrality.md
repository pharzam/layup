# W-12 — Harness-Agent Neutrality

**Item.** S12, Harness-Agent Neutrality (`F-0003#52`): "The rules, the context,
and the task state stay in a form that does not belong to one harness agent. A
second harness agent can do the work, and can check the work of the first
one." Also Invariants 1, 2 and 9, and `NFR-002`.

**The case.** Task `T-7` of a Go target starts on harness H1 (Claude Code). It
stalls, and the Operator gives it to harness H2 (Codex). Later the Operator
stops LAYUP for good, and a human continues the target without it.

Sections are those of [`architecture.md`](../architecture.md).

## Part 1 — the task moves from H1 to H2

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | `layup run` | the task register; the routing register | picks `T-7`; admits H1 with its model for the developer role and the task's tier; the first admitted pair by weight or order (§9) | a routing row | `code` | §9; ADR-0020 |
| 2 | `layup run` | the target at the base commit; the register row of H1 | makes the session directory: a separate clone on `task/T-7/1`, an empty `home/`, the prompt file, an environment from the named list with H1's credential variable and no forge credential; checks that no rule file of H1 is in a directory above it; starts H1 with the command line of its row | a session start row: session ID, task, attempt, harness, model, base commit, records commit | `code` | §4; ADR-0015 |
| 3 | H1 session | `prompt.md` and `repo/` | works; reads its rules from `AGENTS.md` in its clone; commits in `repo/`; writes its typed result | the result file in `result/`; commits in `repo/` (not pushed by the session) | `model` | §4; ADR-0015 |
| 4 | `layup run` | the stalled attempt | later: slice F | later: slice F | `code` | — |
| 5 | the Operator | the stall package on the task's issue | a comment in the fixed form that the stall package gives (later: slice F), asking to give `T-7` to H2 (Decision Point 5) | the comment on the forge | `human` | §3 |
| 6 | `layup run` | the comment, read by the forge API | copies it before acting: body, author ID, comment ID, `performed_via_github_app` (empty), time, SHA-256; the author ID is in `approvers.tsv` and the App field is empty, so it is a human decision | a row of the human-input table, a payload with the body | `code` | §3; ADR-0014 |
| 7 | `layup run` | the parsed comment | later: slice F | later: slice F | `code` | — |
| 8 | `layup run` | the base commit of `T-7`; the diff and findings of attempt 1 (payloads); the records; the register row of H2 | makes a new session directory with a separate clone on `task/T-7/2` from the base commit (§8, Attempts), with the diff of attempt 1 as a payload, with a prompt file built by code from the records (the task, its handoffs, the stall diagnosis) and the target's rule files; starts H2 as in step 2 | a session start row for attempt 2 | `code` | §4; ADR-0015 |
| 9 | H2 session | the same kind of prompt file; no memory of H1 | works from the files alone; `AGENTS.md` is the same file that H1 read; H2's entry file in the target is a pointer to it | its result file and commits | `model` | §4; ADR-0015 |
| 10 | `layup run` | the result of H2 | checks the result's schema; takes the attempt and base commit from its own session start row and checks that attempt 2 is still the open attempt of `T-7` (else it refuses the result); fetches the branch by SHA into its own clone with hooks turned off; checks that it descends from the base commit and touches nothing under `.github/workflows/`; commits the result to the records branch with a push that is not forced; pushes that SHA to `task/T-7/2` with the App's token; then records the SHA as bound to the session | the result, the session end row, the commit SHAs bound to the session | `code` | §3; ADR-0014 |
| 11 | a verifier session | the change | on a harness that wrote no commit of the change by the ledger rows: attempt 2 started from the base, so its only author is H2, and H1 is admitted; refutes "done" (W-07 steps 7, 8) | its record and `layup/verify` | `model` | §8, §9; ADR-0020 |

## Part 2 — a human continues without LAYUP

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 12 | the Operator | — | stops `layup run`; the records stop growing | the last lease row, with its end time | `human` | §2; ADR-0013 |
| 13 | the Operator | the target's rulesets | later: slice B | later: slice B | `human` | — |
| 14 | a human | a fresh `git clone` of the target | reads the default branch (the product, the task files, `AGENTS.md`) and `origin/layup-records` (the tables and payloads, plain text), with no LAYUP tool; the README names the branch | — | `human` | §3; ADR-0014 |
| 15 | the target's CI | a pull request of that human | later: slice B | later: slice B | `code` | — |
| 16 | a human | the open task `T-7`: its issue, its branch `task/T-7/2`, its records | continues the work and lands it through a pull request under the target's own gates; LAYUP is not needed for any step | the target's own records of the task (its issue and pull request) | `human` | §1; ADR-0013 |

## Checklist rows

K01 to K07; P04, P09, P13, P21; S12; R07, R08, R11, R13; I1, I2, I9; C2, C6, C7;
D05 to D08, D17; FT3, FT5, FT6. Known limits: L-A1 (the shared host), L-A2 (one
forge), L-A3 (one host during delivery), L-A4 (the bypass list is read once), L-A5 (a policy file of the host).
