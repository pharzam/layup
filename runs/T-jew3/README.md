# T-jew3 — row 65b of the plan of M2e, the code-owner rule

Issue [#356](https://github.com/pharzam/layup/issues/356); facts `F-0003#43` and `F-0003#47`. Base `8e26eca`. Author: Claude Sonnet 5.5 on Claude Code. The plan, its review (Devin SWE-2 Max,
comment 6214007155: `approve-with-conditions`, 1000 lines over 24 files, cap 2) and the author's answer (6214046545) are comments on #356. Files: [`test-runs.md`](test-runs.md) (the spike, the
red and the green runs), [`mutants.md`](mutants.md) and [`mutate.py`](mutate.py) (the mutation audit).

## One reading (R13) of the phrases that decide

| Phrase | The one reading, and the document that says it |
| ------ | ---------------------------------------------- |
| "the four path classes" | `.github/**`, `docs/gates.tsv`, `docs/gates/**` and each `config` path of the manifest at the head: one root-anchored line each, a path below an earlier line at a `/` has none ([`scaffold.md`](../../docs/spec/scaffold.md#the-code-owner-rule)) |
| "the Operator's account" | the login of the row of the role `operator` of `approvers.tsv`; no row is an error; the repository's owner is never read |
| "the recorded ID" | the positive integer `id` of the object under `== default.json` of `setup/rulesets.txt`; absent, unreadable or not positive is a `fail` that names the recovery |
| "no required approval" | `required_approving_review_count` 0, with `require_code_owner_review` `true`; whether GitHub binds the rule at 0 is unconfirmed and left to the demo of row 68b |
| "preserves existing protections" | the live ruleset is written back whole, reduced to the fields of a write, with only `require_code_owner_review` changed, and read back equal |
| "cannot run" (the landing) | the text that S13 writes starts with a command that exits 1 and names the open question; `landApplicable` is `false`; the child issue [#384](https://github.com/pharzam/layup/issues/384) tracks the tree check |

## The record that later closes this task (nothing is estimated here)

- The first build invocation of this session ended at the outer 30-minute limit (exit 124) before the documents were done; the resource figures of that invocation are **not reported**.
- The plan session's WebFetch ran `claude-haiku-5-5` (comment 6210869150); the attempts of the AGY reviewer gave no qualifying approval (comment 6213859944). Both go into the resource record.
- The budget (1000 lines, 24 files) is exceeded by the branch as measured at the freeze; that is a finding of the last round, not a new maximum.
