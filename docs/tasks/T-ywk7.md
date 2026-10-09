# T-ywk7 — the technical specification of M2b, role sessions on registered harnesses

Issue: [#147](https://github.com/pharzam/layup/issues/147). Milestone `M2b` of the
[implementation plan](../plan/README.md#milestones), which starts with its
specification task (O-114). Serves `F-0003#51` (`REQ-012`) for `REQ-011`,
`REQ-005`, `REQ-003`, `REQ-013`, `NFR-005` and `NFR-001`. Base `ad3c400` (the
merge of #146). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-ywk7/`](../../runs/T-ywk7/).

## Plan and plan review

The plan (R12, comment 6075171228), its review and the author's answer
(6075289966) are comments on #147. The plan review (Claude Fable 5.1, effort
`xhigh`, on the Claude Code CLI with stream output, a fresh read-only session in
a clone at `ad3c400`, 8 min 5 s; comment 6075282781) gave
`approve-with-conditions`: Budget maximum 2,600 lines added plus removed over 26
files against `ad3c400`, close-out inside; Cycle cap 1; no panel. Its two
conditions (the columns of the built block `harness-register` in words; one row
per package) and seven notes are applied.

The Operator's decisions, each recorded on #147:

- **O-178** (b): searcher B is Claude Opus 5.5 in a fresh `claude -p` session,
  as OpenCode's login failed (401 on every model) and Devin's weekly quota ended.
- **O-179**: when Devin's or OpenCode's paid usage ends, their free models are
  used (Devin's SWE-2, recorded as an execution-tier model on a reasoning part).
- **O-180** (a): `layup run` of `M2b` runs the probe only; the developer session
  of the demo runs from the uat test.
- **O-181** (a): Budget maximum 3,200 lines over 28 files (comment 6075909417),
  as the search record came to 1,517 lines.
- **O-182** (a), **O-183** (a), **O-184** (a): one more cycle at each cap; the
  cycle cap rose from 1 to 4.

## What was done

1. **S1, the public-solution search** ([`runs/T-ywk7/search/`](../../runs/T-ywk7/search/), the first commit, `0a01743`): five curated lists chosen by a recorded GitHub search, read in full; 20 searches in the words of the facts that the requirements of `M2b` cite; two blind searchers, GPT-6 Sol on Devin and, by O-178 (b), Claude Opus 5.5 on Claude Code. Both conclude that no public solution can be the base of `M2b`; the summary gives what each part takes from them.
2. **S2, the survey** ([`survey-m2b.tsv`](../../runs/T-ywk7/survey-m2b.tsv), after the search): each of the 11 `M2b` rows of the survey of `T-w73g` read against its cited file at its commit, and taken (7) or rejected (4) with its reason and the heading that applies it.
3. **S3, the decisions** (comment 6075865294): K40, the models and the "not used" list, the probe, the usage report, the events of a task, the fixed variables, the check of `NFR-005`, the review of the release, the comments of a session, the Operator's inputs.
4. **D1, [`docs/spec/session.md`](../spec/session.md)** (new): a role session (its directory, environment and credential, rules only from the target, the start, the context, the limits, the end), the checks before a push (the fetch by SHA through a scratch repository, the workflow and rule-path refusals, the push and the bind), the result and the open attempt, the probe, admission and routing, the writer of the telemetry record, the input states, the check of `NFR-005`, the review of the release, "Not in M2b" and the acceptance tests.
5. **D2, [`docs/spec/records.md`](../spec/records.md#nfr-001--the-records-of-a-session):** eight schema blocks (`sessions`, `harnesses`, `routing`, `routing-register`, `models`, `events`, `result`, `probe-result`) in `notYetBuilt`; the columns that `M2b` adds to the harness register, in words (condition 1); the rule of `payloads/`.
6. **D3, [`docs/spec/packages.md`](../spec/packages.md#the-table-of-m2b):** the table of `M2b` (`internal/session`, `internal/ledger`, `internal/rules`), the rows of `M2a` and phase 1 changed in place (condition 2), the calls of `internal/git` of `M2b`; `gate.md`: the check of `NFR-005`, in the future tense.
7. **D4:** `docs/spec/README.md` (the version line, the files, the second read of `PATH`, the later phases), `run.md`, `docs/plan/README.md` (the row of `M2b`, K40 settled, `NFR-005`, #148), `PRD-0001` §12 (`REQ-012`, `NFR-005`), the glossary (role session, harness probe), and a lesson in `docs/guardrails.md` §2.
8. **Revealed, off the path:** [#148](https://github.com/pharzam/layup/issues/148), the release of `M2a` with no review for `REQ-015` and `REQ-017`, and `release-check.sh` that fails on `main`.

The rows of K23 and K38 of the defect register need no change (note 11 of round 1): this task applies K23 in `packages.md` (`internal/ledger` imports the schemas from `internal/records`) and K38 in `records.md` (`routing.tsv` and the admitted harnesses in `M2b`).

**Tests** ([`test-runs.md`](../../runs/T-ywk7/test-runs.md)): `sections.sh` failed for each of its 28 headings before `session.md` and passes at the head; `survey.sh` failed with no table, then for each `take` until its heading existed, and passes; the block test failed for each of the eight new blocks until they were listed. The search's commit comes before the first commit of `docs/spec/`.

**The rejected alternatives:** ACP's handshake as the probe (a protocol would change ADR-0015); a list column of models in the harness register; a variable passed on from `layup`'s own environment for the credential; the bind as a column of `sessions.tsv`; a fetch that runs `git upload-pack` in the session's clone (round 1, finding 3); a product step for the demo's developer session (O-180 (b)).

**Known limits:** a session can read its own credential (L-A1), and a copied login whose refresh token rotates can be spent by it; the probe's file list is the model's answer; a product whose code holds shell scripts changes them only in a rule batch; money is `unknown` for a session whose usage names more than one model, as the telemetry block holds one price ID; the usage formats of Codex, Gemini CLI and OpenCode wait for a registered harness that needs them.

## Review rounds

The records, the Fixes replies and the decisions at the caps are comments on
#147. Each round: Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with
stream output, a fresh read-only session in a clone at the frozen head, with
`ANTHROPIC_DEFAULT_HAIKU_MODEL` set to the same model; `modelUsage` named that
model only, and no command was denied.

| Round | Commit | Cycle | Verdict | Findings |
| ----- | ------ | ----- | ------- | -------- |
| 1 | `6ed1cf5` | 0 | `material` | 7 material (four parts with no test; a policy file fails the probe; the fetch reads the session's configuration; a refused push; the writer of `attempt`; `detail` of `events`; `vars` that can carry `GH_TOKEN`), 9 notes |
| 2 | `772cc69` | 1 | `material` (edited by O-182 (a)) | 2 material (a refused start before the session ID; a ref with a SHA of no commit), 4 notes |
| 3 | `cdd0432` | 2 | `material` (edited by O-183 (a)) | 1 material (a harness with no model of `use` `yes`), 5 notes |
| 4 | `80ab5b0` | 3 | `material` (edited by O-184 (a)) | 1 material (the size check of the prompt after the start row), 3 notes |
| 5 | `2784b94` | 4 | `nothing material in scope` | 2 notes |

Each material finding was fixed in the next cycle; each finding of rounds 2 and
3 was written by the fix before it. The two notes of round 5 are applied in the
close-out: the size check refuses a prompt over 131,071 bytes, as Linux's limit
of one argument counts its final zero byte; and the three counts of the step
`probe` are named. Note 11 of round 1 is recorded above: the rows of K23 and
K38 need no change.

## Verdict

Delivered: the technical specification of `M2b` in `docs/spec/` (`session.md`,
the records of a session, the table of `M2b`), after a public-solution search by
two blind searchers that found no public base, with K40 settled and the check of
`NFR-005` named. The review ended by decay at cycle 4 of a cap that the Operator
raised three times (O-182 to O-184). The diff against `ad3c400` is inside 3,200
lines over 28 files (O-181).

Next: the slicing task of `M2b`, which writes its rows in the plan and opens
their issues, as `T-zwke` did for `M2a` (O-164); and #148.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-09; tokens are
`modelUsage` of the `result` event of the Claude Code CLI (stream runs), summed
over the models; `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The issue, the sources, the plan | reasoning | Claude Opus 5.5 | max | not reported | to 05:51 (the issue 05:44, the plan 05:51) |
| The plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,220,399 (USD 6.41) | 8 min 5 s, from 05:52:37 |
| The answer; isolate; the sources of the search | — | Claude Opus 5.5 | max | not reported | 06:02 to 06:08 |
| Searcher A | reasoning | GPT-6 Sol (`gpt-6-sol-high`), Devin | not reported | not reported | 4 min 58 s, from 06:08:16 |
| Searcher B, OpenCode (failed: 401) and Devin (stopped: quota) | reasoning | Grok 4.7 | not reported | not reported | 9 s from 06:08:16; 2 min 53 s from 06:09:18 |
| Searcher B (O-178 b) | reasoning | Claude Opus 5.5, Claude Code CLI, five subagents of the same model; `claude-haiku-5-5` inside `WebFetch` | `xhigh` | 22,264,926 (USD 19.61), of them 5,237,499 of Haiku (USD 1.25) | 21 min 50 s, from 06:15:37 |
| Three test runs (the write rule; `ANTHROPIC_DEFAULT_HAIKU_MODEL` off and on); a test of SWE-2 | — | Claude Opus 5.5; SWE-2 High | low | not reported (USD 0.15 by their result events) | under 1 min each |
| The summary, the survey table, the decisions, the specification, the checks | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 06:37 to 07:05 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 1,306,545 (USD 7.41) | 14 min 6 s, from 07:07:02 |
| The fix of round 1; round 2 | execution; reasoning | Claude Opus 5.5; Claude Fable 5.1 | max; `xhigh` | not reported; 820,223 (USD 4.40) | 07:21 to 07:26; 7 min 27 s, from 07:26:42 |
| O-182; the fix of round 2; round 3 | execution; reasoning | Claude Opus 5.5; Claude Fable 5.1 | max; `xhigh` | not reported; 547,027 (USD 3.48) | 07:34 to 08:08; 6 min 37 s, from 08:08:59 |
| O-183; the fix of round 3; round 4 | execution; reasoning | Claude Opus 5.5; Claude Fable 5.1 | max; `xhigh` | not reported; 603,965 (USD 3.86) | 08:44 to 08:46; 6 min 34 s, from 08:46:18 |
| O-184; the fix of round 4; round 5 | execution; reasoning | Claude Opus 5.5; Claude Fable 5.1 | max; `xhigh` | not reported; 597,777 (USD 2.73) | 08:55 to 08:57; 5 min 9 s, from 08:57:34 |
| The close-out, with the notes of round 5 | — | Claude Opus 5.5 | max | not reported | from 09:03 |
