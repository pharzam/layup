# W-01 — Problem Statement Quality

**Item.** S1, Problem Statement Quality (`F-0003#41`): "Detection of the gaps in
a problem statement before delivery starts, so that the questions for the idea
owner come in one batch and not one at a time." Also Decision Points 1 and 2
(`F-0001#10`, `#11`) and `REQ-001`.

**The case.** A problem statement for a small invoicing service. Its §7 table has
the metric "Invoice accuracy ≥ 99 %" with no measurement method, and its text
uses "customer" for both the payer and the account holder. It names no budget.

Sections are those of [`architecture.md`](../architecture.md).

| # | Actor | Input | Mechanism | Record | Tag | Where |
| - | ----- | ----- | --------- | ------ | --- | ----- |
| 1 | the Operator | an empty repository on the forge | installs the LAYUP App on it; runs `layup run --new OWNER/NAME --psb brief.md --operator LOGIN --idea-owner LOGIN`; pushes the unmodified pinned copy as the root commit with the command that `layup run` prints | the root commit of the default branch | `human` | §5 |
| 2 | `layup run` | the repository; `brief.md`; the two logins | reads back the default branch and its root tree; pushes the orphan records branch: `brief.md` byte for byte with its SHA-256, `approvers.tsv` with the two numeric user IDs, the lease row; opens the Intake issue | the records branch; the Intake issue | `code` | §5; ADR-0014 |
| 3 | `layup psb check` | `brief.md` | rules G1 to G5 (`internal/psb`); G2 fires on the metric row "Invoice accuracy ≥ 99 %", which gives no measurement method | a gap table on standard output, committed by `layup run` | `code` | §5 |
| 4 | a review session | `brief.md`; the review prompt with the gap kinds of §5 | finds that "customer" names two parties (a term with two readings), and that no budget is named (a missing intent value); writes each as a typed row with a byte-exact quote; code checks each quote against the file | a gap table in the session result | `model` | §5; ADR-0015 |
| 5 | spec sessions | `brief.md` | later: slice C | later: slice C | `model` | — |
| 6 | `layup run` | the gap tables of steps 3 to 5; the questions of setup steps S01 and S10 from the pinned commit | merges the rows; drops a repeat of the same quote and kind; gives each question an ID | the project question table, each row "before delivery" | `code` | §5 |
| 7 | `layup run` | the batch | posts one comment on the Intake issue: the questions and the intent form for the idea owner; the setup questions and the parameters for the Operator | the comment ID in the records | `code` | §5; ADR-0014 |
| 8 | the idea owner, the Operator | the Intake comment | each answers in one comment, one line per question ID (Decision Points 1 and 2) | the comments on the forge | `human` | §5 |
| 9 | `layup run` | the answer comments | copies each comment (§3); takes as answers only the comments whose author ID is in `approvers.tsv`; checks that every ID has an answer and that each number and login parses | the answers in the records; the intent values, the approvers of each planned point and the parameters as registers | `code` | §3, §5; ADR-0014 |
| 10 | a review session | the questions and the answers | lists each answer that leaves its gap open, for example "customer means the client" with no word on the account holder | a row per open answer | `model` | §5 |
| 11 | `layup run` | the missing and the open answers | asks again for them in one comment; each follow-up is one more row "before delivery" | the question table | `code` | §5 |

## Checklist rows

S1, R01, K08, K12, D06. The step table of §9 (the review session's role and harness) is later: slice D.
