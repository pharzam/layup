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
- `main` at `d5e199f` merged into the branch; conflicts and the checks after
  the merge are recorded below their resolution.
