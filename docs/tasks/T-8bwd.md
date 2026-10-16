# T-8bwd — the handoff items' rules and implementation ownership for M2e

Issue [#387](https://github.com/pharzam/layup/issues/387). Documentation only;
no production code changes. Author: SWE-2 Max (Devin CLI). Plan
[comment-6220167894](https://github.com/pharzam/layup/issues/387#issuecomment-6220167894);
independent plan review
[comment-6220326750](https://github.com/pharzam/layup/issues/387#issuecomment-6220326750):
`approve-with-conditions`, maximum 1,000 changed lines / 14 files against the
current `main` (`d5e199f` at confirmation, merged and then `d09b1c0` before the
stable finish), cycle cap 1, no panel, one goal class; the author's answer is
[comment-6220408093](https://github.com/pharzam/layup/issues/387#issuecomment-6220408093).
Evidence: [`runs/T-8bwd/`](../../runs/T-8bwd/) (`items.sh`, `baseline-red.txt`,
`green.txt`).

## What was done

1. **`docs/spec/loop.md`** — the new section "The items a caller computes"
   maps each of the fourteen handoff items to its input, its rule, the package
   and plan-table row that computes it, and its milestone; pins the field-line
   forms of the plan, the plan review and the review record; states that
   `internal/handoff` computes `schema` and `artifacts`, that `internal/run`
   computes `question-row` and `answer-row` from the result (`RecordHandoff`,
   row 66d) and that every other item is the caller's computed result, a
   missing result being false; records the question/answer comments as composed
   by `internal/run` from the `question`/`answer` rows and posted on the task's
   issue by the step-loop machinery; marks `patch-per-kind`, `register-rows`,
   `heading-ids` and the `tests-red-green` producer as `M2f`/`M2g` in the
   transition table and "Not in M2e"; keeps the `M2e` contract of
   `tests-red-green` as row 66e's typed per-test results with no `M2e`
   producer, so no `M2e` build-task handoff can be valid; adds the acceptance
   row "The item computers of the task loop".
2. **`docs/plan/README.md`** — new row 66l (`T-4pje`, #396, the item
   computers of the task loop), `After` `61, 62c, 63a, 66d`; rows 66f, 66k and
   67 take `66l` in `After`; row 66d's scope names `question-row` and
   `answer-row`; the `M2f`/`M2g` milestone cells name the items and producer
   their specification tasks take.
3. **`docs/spec/packages.md`** — the `internal/run` cell names row 66l's
   item computers and separates them from row 66k's wire.
4. **Registration** — `T-4pje`/#396 opened by the App after a duplicate
   check (`docs/tasks/` files, branch list, issue search); backlog lines for
   `T-8bwd` and `T-4pje`; #396 gets its own plan and plan review before any
   implementation.
5. **`docs/tests/traceability.md` and PRD-0001 §12** — deliberately not
   edited: their conventions register named tests and the tasks that deliver
   them (the acceptance-table row names of `loop.md` are not individually
   listed; `M2e/uat/role-handoffs` is a test ID, not a row name). Row 66l's own
   landing names its tests; this task lands none, so no row and no §12 Task
   cell is written here.

## Decisions and rejected alternatives

- **Owner of the item computers and the question/answer composer:** a new
  row 66l, the plan review's decision. Rejected: assigning the computers to
  the existing rows 66f, 66h and 66k — each would add an independent parser or
  builder to a row whose O-214 goal count of two is already spent on its own
  classes, and would conflate statuses, Shape and the wire with computing
  comment contents. Rejected: a second new row for the composer — it shares
  the comment-parse machinery and the records the computers read.
- **Comment form:** field lines — a lowercase key at column 0, `: `, the
  value to end of line, `—` the empty value — the form `set` and the answer
  lines already use. Rejected: a fenced TSV or a markdown table for the plan,
  plan review or review record — a second grammar for fields that
  `architecture.md` §9 and the loop already name; the verdict words are those
  `review-record-lint.sh` and `engineering-discipline.md` read, and a target's
  comment is a different artifact from LAYUP's own review record (a
  `| Field | Value |` table), which the note states.
- **Question/answer comment writer:** `internal/run` composes the body from
  the `question`/`answer` rows and the step loop posts it on the task's issue.
  Rejected: the session writing the body — two writers of one comment would
  fight the single-writer rule of the session comment files, and a restart
  could not recompose the same bytes from the records.
- **`tests-red-green`:** row 66e's typed per-test results stay the `M2e`
  contract; the producing gate run and the ID-to-file map are `M2g`'s. No
  `M2e` run claims a valid build-task handoff.
- **Contradictory wording:** `architecture.md` §8's "the checks parse fixed
  verdict fields" is met by the pinned `verdict:` field; the plan-review words
  are the check's own (`review-record-lint.sh`, RR8), so `loop.md` is the one
  operative home of the exact field names and values.

## Checks so far (pre-review)

- `runs/T-8bwd/items.sh`: baseline at `d5e199f` — 44 assertion failures, the
  mapping section absent ([`baseline-red.txt`](../../runs/T-8bwd/baseline-red.txt));
  after the edits — exit 0 with the one deliberate note on `tests-red-green`
  ([`green.txt`](../../runs/T-8bwd/green.txt)). One-off evidence; no hook, no
  CI job, no `run-discipline-tests` wiring.
- Text checks: `adr-lint.sh`, `prd-lint.sh`, `link-lint.sh` (2,781 links),
  `run-discipline-tests.sh` (138 passed), `setup-check.sh`, `git diff --check` — all OK.
- `go build ./...`, `go vet ./...` — clean.
- Unit: `go test ./internal/handoff ./internal/records ./internal/run ./internal/tsv` — ok.
- Document-reader integration (explicit `-run`, exit 0): `internal/run`'s
  `TestTheParameterTableEqualsTheSpecification`, `TestTheSchemaOfTasksEqualsItsBlock`,
  `TestTheSchemaOfParametersEqualsItsBlock`, `TestTheTextBlocksOfRunMd`;
  `internal/handoff`'s `TestTheSchemaEqualsItsBlock`,
  `TestTheDefaultRowsAreTheTableOfTheSpecification`,
  `TestTheListsAreThoseOfTheSpecification`; `internal/records`'
  `TestTheSchemasEqualTheirBlocks` and `TestEverySchemaBlockIsBuiltOrNotYetBuilt`;
  `internal/tsv`'s `TestEverySchemaBlockIsBuiltOrNotYetBuilt`; `cmd/layup`'s
  `TestPackageRules`.
- **Interrupted, not a result:** one `go test -tags=integration ./internal/run`
  run was started unfiltered, competed with the serialized full-integration
  run and was stopped before it reported; it is recorded here as interrupted
  and counts as neither pass nor fail. No full-module or e2e run was made by
  this task.

## The producing-document audit and its repairs

The pre-first-freeze audit
([comment-6220740712](https://github.com/pharzam/layup/issues/387#issuecomment-6220740712))
found three in-scope defects in the committed artifact `356f161`; the author's
answer is
[comment-6220818955](https://github.com/pharzam/layup/issues/387#issuecomment-6220818955):

1. The section's opening sentence assigned `question-row`/`answer-row` to
   `internal/handoff`; `internal/run.RecordHandoff` (handoffpost.go:287)
   computes them. Sentence corrected (R10); the table already named the right
   owners.
2. `items.sh` check 1 used membership joins, so a duplicated mapping row
   stayed green. Repaired to exact equality: a `uniq -d` duplicate diagnostic
   plus the per-item member diagnostics.
3. `items.sh` check 3 grepped the whole plan row, so an item appearing only
   in the Tests cell satisfied the owner anchor. Repaired to grep only the
   owner's Task and Demo cells (columns 4 and 5). The `tests-red-green`/66e
   exception (C1) is unchanged.

Controls in this worktree, run on the committed artifact and restored:
duplicating the `plan-fields` row — old exit 0, repaired exit 1 naming
`plan-fields` and its anchors; `plan-fields` only in row 66l's Tests cell —
old exit 0, repaired exit 1 naming owner row 66l; each restored copy exit 0.
Runs recorded in [`runs/T-8bwd/adversarial.txt`](../../runs/T-8bwd/adversarial.txt).

The #396 boundary note
([comment-6220700567](https://github.com/pharzam/layup/issues/396#issuecomment-6220700567))
is registered: the `findings-located` parse is the input the step loop's
caller boundary reads for the uphill rule of `architecture.md` §11 — the
parser is row 66l's, the wiring stays #367's; no O-214 count changes and
nothing is claimed implemented.

## Open

The independent review round, the verdict, the resource record, the
completed-log line and the merge are close-out and are not written yet.
