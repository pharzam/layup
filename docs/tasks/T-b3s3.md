# T-b3s3 — the owner role of an accepted answer

Issue [#401](https://github.com/pharzam/layup/issues/401), row 66m of the plan of `M2e`. Serves REQ-005/F-0003#45.
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

Pending the root audit and the independent review round.
