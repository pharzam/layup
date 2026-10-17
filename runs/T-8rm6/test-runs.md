# T-8rm6 — test runs

Evidence for #365 (the approvers and the success criteria; the revision of the
draft). Every run in the worktree `.worktree/T-8rm6` unless noted. The unit
suite is `go test ./...`; the tagged tests are
`go test -tags=integration ./internal/run -run '<names>'`.

## Red before code (test first)

- `go test ./internal/run ./internal/records` — the new tests each failed for
  the named reason before any production change: `TestTheCommentOfTheBatch`
  (the `success-criteria` question and the point lines missing from the
  comment), `TestThePointBindings` / `TestTheApproverRows` /
  `TestTheApproverOfEachPoint` / `TestTheSuccessCriteriaOfTheAnswer` (no such
  functions), `TestCheckRevision` / `TestKeepOfRevision` / `TestFactOfQuote`
  (no such functions), `TestFollowUpRows` / `TestTheCommentOfTheFollowUp`
  (criteria field absent), `TestQuestionRefs` in `internal/records`
  (`success-criteria` refused as a ref).

## Audit-fix reds (after comment 6219977060)

- `TestCheckRevision` — red on `a context fact newly covered` and `a context
  fact dropped`: `the class of the fact 4 changed, and no row that covers it
  names a fact of the answers` (the audit's C6 falsifier, reproduced by my own
  production-reader fixtures).
- `TestAnAccountARowAlreadyServesIsNamedByAnotherSpelling` — red: a second row
  `[101 approver PHARZAM ...]` was written for the recorded Operator's account,
  and `activation` resolved to the idea owner's row 303 (the audit's C4
  falsifier, reproduced with distinct Operator/idea-owner IDs 101/303).

## Green

- `go test ./...` — all packages ok after the implementation; again after the
  C4/C6 fixes and the cap revert.
- `go build ./...`, `go vet ./...` — clean.
- The five checks: `docs/adr/adr-lint.sh`, `docs/prd/prd-lint.sh`,
  `docs/links/link-lint.sh`, `docs/tests/run-discipline-tests.sh`,
  `docs/setup/setup-check.sh`; `git diff --check` clean.

## Compiled mutants (each applied, killed, restored)

First round, the base implementation — 14 mutants, all killed (quote-fact
clause, class-change check, criteria-question refusal, criteria append, role
check of `pointBindings`, first-wins `bound`, latest-wins `named`, 404
handling, idea-owner fallback, need/constraint coverage, measure criterion,
`answered` check — the last needed a new case where the unrowed Q-NNN's quote
fact is covered).

Second round, the audit fixes — 7 mutants, all killed:

| Mutant | Killed by |
| ------ | --------- |
| `rowed` keyed by `a[2]` (login) not `a[0]` (ID) | `TestAnAccountARowAlreadyServesIsNamedByAnotherSpelling` — duplicate row of account 101 |
| `if rowed[key]` skip dropped | same — duplicate rows of 101 and 707 |
| `rowed[key] = true` mark dropped | same — second row of 707 for `DAVE` |
| `strings.EqualFold` back to `==` in `approverFor` | same — `activation` fell back to row 303 |
| `len(old[n]) > 0` dropped from the class-change guard | `a context fact newly covered` |
| `len(kinds[n]) > 0` dropped | `a context fact dropped` |
| guard removed entirely (`!maps.Equal` alone) | both cases |

## Integration (scripted harnesses, real `Restart`/`PhaseSession`)

- `go test -tags=integration -timeout 480s ./internal/run -run
  'TestTheRevision|TestTheStepForgeCheck'` — ok, 107.7s, with the world cap
  back at `50.0,8.0`. `TestTheRevisionOfTheDraft` (criteria rows appended,
  approver row tied to its point by `source`, done-before on restart, batch
  repost on a lost comment), `TestTheRevisionRefusedThenKept` (refusal →
  finding in the next prompt → kept, one telemetry row per attempt),
  `TestTheRevisionStopsAtTheCap` (the clock jumped past 8.0h stops the retry),
  and `TestTheStepForgeCheckChecksTheForgeAndCommitsItsTableAtEachRestart`
  (the world whose harness exercises the default `revision.tsv` fallback —
  the run that had looped ~577s on the malformed six-column output before the
  fixture emitted the four-column session form).

## Not run here

The full-module `-tags=integration` run and the e2e suite are the root
coordinator's, serialized after this head is stable. An earlier combined
invocation during development was interrupted at exit 143 inside the invalid
retry loop the malformed fixture caused; it is not counted as a pass anywhere.

## Post-audit additions (comment 6220264944)

- `TestTheCapOfMilestoneZeroStartsNoSession` — the revert of the provisional
  cap had missed its clock jump (still `96 * time.Hour`); restored to
  `8 * time.Hour`, re-run green in 9.1s. The file is again byte-identical to
  `a3a4cdd` for its remaining lines.
- `docs/architecture.md` §5 step 5 qualified: a named login is added to
  `approvers.tsv` only when its account holds no row already (the key is the
  numeric ID); operative detail stays in `loop.md`.
- `docs/tests/traceability.md`: three rows for the new unit/integration tests
  of this task, and `T-8rm6` appended to the rows of the tests it extended.
- `main` at `d5e199f` merged into the branch (merge commit `0334197`, strategy
  `ort`): no conflicts; `docs/prd/PRD-0001-layup.md`, `docs/spec/loop.md` and
  `docs/tests/traceability.md` auto-merged, keeping both sides. After the
  merge: `go build ./...`, `go vet ./...` and `go test ./...` all clean.

## The root's full integration on `773bcf9` — FAILED (direct exit 1)

- Command (the root's): the complete `-tags=integration` module run on
  `773bcf91ce32d38328fadeb45909eab8e7bcb321`,
  2026-10-16T16:52:02.899492Z – 17:10:03.037731Z, direct exit 1. Exactly one
  test failed: `TestTheFollowUpAsksTheMarkersAgain`
  (`internal/run/markers_integration_test.go`) — bug sub-issue #397.
- Cause (fixture, not production): the test's idea-owner answer map supplied
  only `band.B`, `band.U` and `appetite.hours`. With `success-criteria` a
  fourth intent question, the map returned `""` for it, the follow-up asked it
  again, and the frozen assertions (3 marker questions, Operator's block
  waiting) broke: 4 questions, both blocks waiting.
- Reproduced in this worktree at `773bcf9`: `go test -tags=integration
  -timeout 300s ./internal/run -run 'TestTheFollowUpAsksTheMarkersAgain'`
  fails identically (11.1s).
- Repair: the map gains `"success-criteria": "it does its work"`; every
  original assertion unchanged. Green at 11.8s; removing the correction
  restores the failure (the "4 questions" detail); restoring it is green
  again.
- `main` at `d09b1c0` merged after the repair (merge commit `7e621d7`,
  `ort`, no conflicts; `PRD-0001`, `loop.md`, `traceability.md` and
  `forgecheck_integration_test.go` auto-merged, both sides kept — the
  codeowner work of `T-4cep` and this branch's step). Post-merge:
  `go build`, `go vet`, `go test ./...` clean.
- Focused intersection run (explicit, the markers/forgecheck/revision worlds
  on the merged head), `go test -tags=integration -timeout 480s
  ./internal/run -run 'TestTheRevision|TestTheStepForgeCheck|
  TestTheStepMarkers|TestTheFollowUpAsksTheMarkersAgain|
  TestTheStepBatchAsksTheSetupQuestions|TestThePinnedTreeOfTheMarkers'`:
  11 tests, all PASS, 211.5s — `TestTheStepForgeCheck...` 89.5s exercises the
  codeowner row with `revision` in the world, `TestTheFollowUpAsksTheMarkersAgain`
  11.1s passes with the repaired fixture.
- The next complete `-tags=integration` run and the final combined-head CI are
  the root's; #397 is fixed pending them, not closed.

## Second C4 finding — account-ID serving (audit comment 6221307548)

The producing audit found a second instance of the login-as-key error, at the
other end of the same helper pair: `approverRows` correctly declines a second
row for an account already recorded, but `approverFor` picked the point's
serving row by case-folded login *text*. A point naming a login that resolves
to an existing row's numeric ID under a spelling beyond case fell through to
the idea-owner row.

- Author response before the fix: `#365` comment 6221341023.
- Regression `TestThePointIsServedByTheResolvedAccount`
  (`internal/run/answers_test.go`, author's own): `bet` names
  `operator-today`, resolving to the Operator's recorded ID 101; `audit`
  names `carol-renamed`, resolving to the idea owner's ID 303; `setup` names
  a login the forge answers 404. RED at first run — `approverFor("bet")`
  returned `[303 idea-owner carol …]` and `pointsServed` omitted `bet` —
  the audit's exact counterexample (distinct IDs 101/303 so the fallback
  could not mask it).
- Repair (bounded, existing caller boundary): `approverFor` and
  `pointsServed` take `(ctx, idOf)` — the same resolver shape
  `approverRows` uses — and the bound login resolves to a numeric ID
  before the row is chosen: the row whose `a[0]` is the resolved ID serves,
  whatever role it holds and however its stored login is spelled; a 404 is
  input only and falls to the idea owner as before; any other lookup error
  is the call's error and chooses nothing. No new record kind, column,
  point, role or authority; `approverRows` unchanged. #366's caller must
  pass a real `Forge.UserID`-shaped resolver and decide its own error
  policy — the helpers now return errors for it.
- Compiled controls (each applied, killed, restored):
  M1 text-match in place of the ID match → `bet` returned 303, killed;
  M2 non-404 errors treated as unknown → the idea owner answered a 503
  lookup, killed;
  M3 ID match dropped → every point fell to the idea owner, killed;
  M4 `pointsServed` swallowing the error → no error returned, killed.
  Restored `internal/run/answers.go` sha256
  `420b7d26e6924785ceff880c653a18d176fc02a8c757b4b4af77c409322c733a`;
  `go test -count=1 ./internal/...` green, `go build`, `go vet` clean.
- Doc sync (R10): `loop.md` Decided-here gains the serving-resolution and
  fail-closed clauses; `traceability.md` names the new test.

## Coordinator full rerun and current-main integration, 2026-10-17

The repaired source a79d26e5bb1a805b9d123c0b099e42dcca0d782f passed
`go test -count=1 -tags=integration -timeout 25m ./...`, all21 packages, direct0.
Recorded wrapper start10:41:17.688780UTC/end10:57:17.780785UTC,960092.005ms.
The original failed run remains recorded as failure; this pass includes #397's
fixture repair and the account-ID serving fix. Current main3c9c4be (including
#351) was then merged before first code review. Its shared REQ-012 conflict
is the additive union of both Task/Test cells, with other cells unchanged.
The new account-ID regression name is also added to that PRD Test cell;
traceability already names it. No production code changed in this integration.
Final exact combined-head CI is still required; the full run above was on a79d26e.
