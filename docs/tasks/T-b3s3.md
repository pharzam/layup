# T-b3s3 — the owner role of an accepted answer

Issue [#401](https://github.com/pharzam/layup/issues/401), PR [#403](https://github.com/pharzam/layup/pull/403), row 66m of the plan of `M2e`. Serves REQ-005/F-0003#45.
Author: Sonnet 5.5. Initial plan 6233490565, plan review 6233654484 (approve-with-conditions, SWE-2 Max), author
response 6233745458 before any edit. The Operator's "Approve #401 coupled scope" (comment 6233711446) keeps both goal
classes in the issue: A the owner-role predicate in `AnswerAccepted`, B the trusted data feed through `Position` and
`OpenQuestions`. Maximum 850 changed lines / 16 files (whole branch against current `origin/main`), cycle cap 1, no
panel. Evidence: [test runs and compiled mutants](../../runs/T-b3s3/test-runs.md).

## Delivered (before review)

`AnswerAccepted` accepts an answer only when the role of the session that wrote it, read from the start row of
`sessions.tsv` (`Roles`, built by `sessionRoles`; `layup run` writes it, no result can claim it), is the one owner role
of every question of the attempt it answers (`Owners`, the owner map, `RoutingRole` of its owner). A session with no
start row, a human question, a label of no owner and questions of more than one owner role fail closed. `Roles` and
`Owners` are two named map types, so a swap does not compile. `OpenQuestions` forwards them and `Position` reads
`sessions.tsv`. The rule has one operative home, `loop.md` ("A question during the work", decided here by this task);
the other sites point there.

## Reachability, as found

A `developer` session's `answer` row passed every reader (`CheckSession` takes `role` as free text, `ReadResult` checks
the form of the row, `RecordHandoff` reads no session role) and `AnswerAccepted` accepted it, so a wrong-role answer
could close a question and move a task downhill once #367 wires `Position`. The main-time counterexample is the red
run in the evidence.

## Alternatives rejected

A per-role rule on `answer` rows in `CheckRow`/`TaskSession` (the writer cannot see the question's kind, and the record
stays valid); "any" question's owner (an owner of one question would close another's); helper-only before #361 (it
duplicates `RoutingRole` and leaves the feed untested); widening #361.

## Known limits

- `RecordHandoff` reads no session's role: the `From`/`To` of a handoff are role-class labels of the transition table,
  so the handoff of a wrong-role answer is still `valid`; it closes no question and is left to the step loop (#367).
- Mixed-owner and human-routed attempts are never closed by a session's answer; who resolves them is #367's or later
  text's. #367 passes the confirmed `owners.tsv` (phase 3) to `Position`'s successor when it exists.
- The label and `needs_human` of a question are the session's own words in `M2e` (`loop.md`); the check is as strong
  as that input.

## Verdict and resource record

Fresh independent SWE-2 Max round1 on ef435d0, record6234079442:
**nothing material in scope**, cycle0/cap1. Its producing-record audit corrected
piped test-status claims with direct saved statuses, disclosed an invalid first
feed mutant, and distinguished two compiled assertion kills from one type-level
rejection. Own overflow/temporary scratch and byte-exact restored controls are
disclosed; no author source change followed the round.

Root full integration passed all21packages on ef435d0, direct0,
1000084.734ms. All12 required exact-head checks also pass on that pre-closeout
head, including CI-only end-to-end. Final closeout-head CI remains required.
The last commit changes only task detail/backlog/completed bookkeeping. The
complete current-main landing figure is recorded on the issue, within the
unchanged850/16 maximum; no overrun approval is inferred.

Notes disposition (no code or ordinary document patch after the review):
1. Retain records.md's existing answer-ID definition; it asserts no new role check.
2. Retain representative PRD test names; traceability registers all six new tests.
3. Retain the unconditional sessions.tsv read: corrupt records fail loudly.
4. Disclose the evidence's20-case wording as a typo: the new matrix has21 cases,
   four positive arms and17 refusals. Existing refusal matrix remains unchanged.
5. Retain the matching operator literals; a shared constant is optional.
6. Retain the awkward pointer comment: its operative target is correct.
7. Measure complete file count after the completed-log move, including all files.

### Actual resources

Reported subset only. Sonnet I/O/Cw/Cr are uncached input/output/cache creation/
read, with thinking inside output. SWE P/O/C are prompt/output/cache, C already
inside P. Retained metric rows are deduplicated across exports; missing elapsed
values are not inferred.

| Part | Expected tier | Model | Effort | Reported tokens | Elapsed |
| --- | --- | --- | --- | --- | --- |
| Plan | reasoning | Sonnet5.5 | not reported | I54/O28160/Cw105199/Cr1573593 | 288360ms |
| Independent plan review | reasoning | SWE-2 Max | not reported | P6114863/O60445/C5967040 | not reported |
| Implementation/tests/docs | execution | Sonnet5.5 | not reported | I54/O34129/Cw65439/Cr4105045 | 1000567ms |
| Evidence organization/correction | execution | Sonnet5.5 | not reported | I6/O2292/Cw5111/Cr542507 | 44017ms |
| Independent round1 including same-head record audit | reasoning | SWE-2 Max | not reported | P8327573/O65037/C8216704;72 unique rows | not reported |
| Root full integration | — | deterministic | — | not applicable | 1000084.734ms |
| Root producing audit/metadata/coordination/closeout | reasoning/execution | Codex GPT-6 family | not reported | not reported | not reported |
| **Total reported subset** | | | | 21029507 tokens; cache counted once | 2333028.734ms |

Requested models do not reach the repository's named reasoning tier; this routing
limit and model-level independence are explicit. The total is a known subset,
not whole-session usage or elapsed wall time. No cumulative-plus-increment double
counting and no elapsed estimate from transcript timestamps.
