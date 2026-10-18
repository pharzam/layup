# T-gyvo — the question and answer comment composers

Issue [#405](https://github.com/pharzam/layup/issues/405). Serves REQ-005/NFR-001/NFR-004/F-0003#45.
Group (c) of the Operator's three-group ruling (comment 6235098876): the comment items are
`T-4pje` (#396) and the tree items `T-yeav` (#404). Author: Sonnet 5.5.
Plan comment 6235211104; plan review 6246374389 (GPT-6.1 Sol, approve-with-conditions);
the first review 6233740566 supplies C1-C12 (C7 is this group's). Maximum 800 changed
lines/15 files against main at landing, cycle cap 1, no panel.
Evidence: [runs](../../runs/T-gyvo/test-runs.md) (red, green, mutants, gates, merge of main).

## Delivered and boundaries

`QuestionComment(start, row)` and `AnswerComment(start, row, questionID)` in
`internal/run/qacomment.go` give the body `first line + "\n\n" + row text`, no final
newline. The first line is `firstLinePrefix(start)` (one home in `end.go`, also used by
`End`) plus `question <S>.<n>` or `answer <S>.<n> of <q-id>`; `n` is the result column
`n`. A start row or row off the reader's schema, types or rules (read as a one-row table,
no repair), a row of another kind, an ID off `answerForm`, or a body that
`CheckCommentBody` refuses gives `(nil, err)`. Pure functions: no post, no record, no
event (posting is #367's, row 66k); no new artifact, kind, column or role.

## Tests mapped to the criteria

| Criterion | Test (red on the stubs, semantic assertion) |
| --- | --- |
| Exact bytes, separator, no final newline | `TestTheQuestionCommentIsItsFirstLineABlankLineAndTheText`, `TestTheAnswerCommentNamesTheAnswerAndTheQuestionItAnswers` |
| ID from column `n`, never text or reason (C7a) | `TestTheIDUsesTheColumnNNeverTheTextNorTheReason` |
| One prefix home shared with `End` (C7b) | `TestTheFirstLinePrefixIsTheOneOfTheEndAndTheComposers`, the existing end tests |
| `questionID` shape-checked, not resolved (C7c) | `TestAnAnswersQuestionIDIsCheckedByItsShapeAndNeverResolved` |
| Malformed input, schema/type boundary, one-token controls (condition 2) | `TestAMalformedRowGivesAnErrorAndNoBytes` (72 subtests) |
| Body refusal, limit in characters, no bytes (C7d) | `TestABodyThatCheckCommentBodyRefusesGivesAnErrorAndNoBytes` |
| Determinism, inputs unchanged | `TestTheSameRowsGiveTheSameBytesAndTheInputsAreNotChanged` |
| Demo on the rows of a real end | `TestTheQuestionAndAnswerCommentsAreComposedFromTheRowsOfARealEnd`, `TestARowOfARealResultThatIsNoQuestionOrAnswerIsRefused` |

## Alternatives

| Decision | Selected | Rejected, and why |
| --- | --- | --- |
| Question ID | session ID + result column `n` | the row's value (the text) or reason: neither is the number |
| Answer's `questionID` | explicit parameter checked by `answerForm` | scanned from the reason (a legal row need not cite it) or resolved against records (#367's attempt semantics) |
| Input validation | `tsv.Read` of the row as a one-row table, equal to the input, then the row rules | the row checkers alone (admit a tab, `02`, `bad`); `tsv.Write` (repairs bytes silently) |
| Separator and end | one blank line, no final newline | `\n` (renders as one paragraph); a final newline (not in the spec) |
| Prefix | `firstLinePrefix` in `end.go`, used by `End` and both composers | a copy in each composer |
| Evidence files | one `runs/T-gyvo/test-runs.md` with the mutant table | a second `mutants.md`: counted against the file maximum |

## Verification and verdict

Own test-first unit tests, the integration demo and 22 valid compiled mutants (all
killed; none that failed to build was counted) cover the bytes, column `n`, the shared
prefix, the `questionID` form, the reader boundary and the body refusal. The whole
integration level passed on the merged tree (exit 0, 09:02:25Z to 09:28:19Z on a loaded
host); build, vet, unit and the five text checks pass. Final exact-head CI remains a
landing condition, with its result on the issue.

Round 1 (GPT-6.1 Sol, xhigh, on Codex): **material**, one finding. F1: invalid UTF-8 in
the row text returned a wrapped `*tsv.Error`, not the `HandoffRefusal` that the plan
and plan-review condition 2 promise. Fixed at c4222e5, test first (red recorded in the
evidence), mutant M22 killed. Round 2 (comment 6248469200, Sol xhigh): **nothing
material in scope**. Harness exception: ADR-0026 names `claude`, `devin` and
`opencode`; Codex was used on the Operator's message of 2026-10-13, for both rounds.
Main (cf2b319, with #407 and #365) was merged with six index and document conflicts
resolved by keeping both sides.

Complete landing figure against main is recorded on the issue at the merge.

## Resource record

ADR-0007, recorded, not budgeted. Tokens are as the harness reported them.

| Part | Tier | Model | Effort | Tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Claude Sonnet 5.5 | not recorded | not recorded | not recorded |
| Plan review | reasoning | GPT-6.1 Sol | xhigh | 173,860 | 2026-10-18T05:58:41Z to 06:04:18Z |
| Implementation | execution | Claude Sonnet 5.5 | not recorded | 171,537 | 45 min |
| Merge of main | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 29 min |
| Review round 1 | reasoning | GPT-6.1 Sol | xhigh | 185,729 | 09:30:18Z to 09:41:35Z |
| Fix of round 1 | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min |
| Review round 2 | reasoning | GPT-6.1 Sol | xhigh | 152,038 | 09:47:22Z to 09:54:13Z |
| Close-out | execution | Claude Sonnet 5.5 | not recorded | not recorded | about 5 min (to 09:58Z) |
| Orchestration | reasoning | Claude Opus 5.5 | not recorded | not separately recorded | not separately recorded |
| Whole integration level | deterministic | | | | 26 min (merged tree) |
