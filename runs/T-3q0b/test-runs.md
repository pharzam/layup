# T-3q0b — the runs of the tests

Times from the clock of the host, UTC, 2026-10-16.

## Red for the right reason

- 06:25:44Z, before the code of row 62b existed (`bet.go`, `answer.go`,
  `inventory.go`, `render.go`, `quote.go`), on `go test ./internal/spec/
  ./internal/records/ ./internal/run/`:
  - `internal/spec` fails to build with `undefined: Answer`, `Finding`,
    `ParseAnswer`, `ReqAnswer`, `VerifiedBatch`, `Brief` (the tests of the bet
    checks, the answer, the inventory, the rendering and the quotes).
  - `internal/records` fails `TestValidIntakeRecordsAreRead` with
    `draft: line 4, column "covers": "1 Q-001" is not a value of list(int)` —
    the column still rejects the Q-NNN token it is to accept.
  - `internal/run` fails `TestKeepOfDraft` with the same `list(int)` refusal —
    the M2c regression case is red on the record read, where it belongs once
    the column is `list(text)`.

## Each rule's test fails when that rule alone is broken

The mutants, applied one at a time to a copy of the one file, `go test
-count=1 ./internal/spec/` run against each, the file restored after. Runs of
07:10Z–07:49Z. Each line names the rule broken and the test that killed it; a
survivor meant the test named the rule's own detail, then the mutant died.

`bet.go`:

| Rule broken | Killed by |
|---|---|
| an unresolved link passes | TestTheLinks |
| a linkless brief passes | TestTheLinks |
| a reference link passes | TestTheLinks |
| a batch the caller did not verify passes | TestTheApproveBatchLines |
| a malformed appetite line passes | TestTheFixedLines |
| a malformed architecture line passes | TestTheFixedLines |
| the appetite line counts outside `## The appetite` | TestTheFixedLines |
| the architecture line counts outside `## The rabbit holes and the no-gos` | TestTheFixedLines |
| the money cap is unchecked | TestTheAppetiteAgainstTheBand |
| the hour cap is unchecked | TestTheAppetiteAgainstTheBand |
| the milestone-0 spend is unchecked | TestTheAppetiteAgainstTheBand |
| an unset `brief.lines` is accepted | TestTheLineLimit |
| a last line without a line feed is not counted | TestTheLineLimit |
| an unknown heading is accepted | TestEachHeadingIsRequiredOnce |
| a missing heading is not named | TestEachHeadingIsRequiredOnce |
| a repeated heading is not named | TestEachHeadingIsRequiredOnce |

`answer.go`:

| Rule broken | Killed by |
|---|---|
| a requirement of the brief is not required | TestTheRequirementLines |
| a requirement the brief did not name passes | TestTheRequirementLines |
| a requirement twice passes | TestTheRequirementLines |
| `phase —` off a Won't passes | TestTheRequirementLines |
| a phase off the target's tags passes | TestTheRequirementLines |
| a fact the records do not know passes | TestTheFactLines |
| a fact twice passes | TestTheFactLines |
| a class of no four passes | TestTheFactLines |
| an option that is no survivor passes | TestTheOptionLinesAndTheVerdict |
| a decision of two survivors needs no option line | TestTheOptionLinesAndTheVerdict |
| a `no bet` needs option lines | TestTheOptionLinesAndTheVerdict |
| the verdict is not required | TestAnAnswerThatFollowsTheFormParses, TestTheOptionLinesAndTheVerdict |
| a line after the verdict passes | TestTheOptionLinesAndTheVerdict |
| `no bet` is no verdict | TestTheOptionLinesAndTheVerdict |

`inventory.go`:

| Rule broken | Killed by |
|---|---|
| a `—` column of a version passes | TestTheInventoryReaders |
| a `—` column of the index passes | TestTheInventoryReaders |
| a copy of 0 passes | TestTheInventoryReaders |
| a fact class of no four passes | TestInventoryRefusals |
| a copy ID that is no positive comment ID passes | TestInventoryRefusals |
| a missing bet milestone passes | TestInventoryRefusals |
| an index out of order passes | TestInventoryRefusals |
| a second version needs no first | TestInventoryRefusals |
| an answer's class of no four passes | TestInventoryRefusals |
| an out mark is not applied | TestTheFirstInventory, TestTheNextInventory |
| version n miscounts | TestTheFirstInventory, TestTheNextInventory |
| the prior rows are not carried forward | TestTheNextInventory |
| the prior index rows are dropped | TestTheNextInventory |
| the hash is of other bytes | TestTheFirstInventory, TestTheNextInventory |

`quote.go`:

| Rule broken | Killed by |
|---|---|
| the quote's bounds are not checked | TestQuoteFactOfFindsTheOneContainingFact |
| the span holds one side only | TestQuoteFactOfFindsTheOneContainingFact |
| one quote offset passes | TestCheckQuoteFacts |
| the follow-up chain is not walked | TestCheckQuoteFacts |
| an answer the answers lack passes | TestCheckQuoteFacts |
| the quote's fact is not required in covers | TestCheckQuoteFacts |

`render.go`:

| Rule broken | Killed by |
|---|---|
| a row outside §6/§7 is edited | TestRenderPRDChangesOnlyThePermittedCells |
| the number order is not numeric | TestRenderPRDChangesOnlyThePermittedCells |
| the answer-fact order is inverted | TestRenderPRDChangesOnlyThePermittedCells |
| a Won't gets no `—` | TestRenderPRDChangesOnlyThePermittedCells |
| a PRD with no non-goals is written into | TestRenderPRDChangesOnlyThePermittedCells |
| an out-of-scope fact loses its mark | TestRenderFactsOfTheVersion |
| a requirement with no row passes | TestRenderPRDRefusals |
| a requirement of two rows passes | TestRenderPRDRefusals |

## Green and the validation, at the content of the head

07:50Z, on the focused content before the merge of `d24fae3` (row 62a's
integration, the landing base at the time):

- `go test -count=1 ./internal/spec/ ./internal/records/ ./internal/run/` —
  all pass.
- `go build ./...`, `go vet ./...` — clean.
- `go test -count=1 ./...` — all 21 test packages pass.
- `sh docs/adr/adr-lint.sh` — OK; `sh docs/prd/prd-lint.sh` — OK;
  `sh docs/links/link-lint.sh` — OK, 2,702 links resolved;
  `sh docs/setup/setup-check.sh` — kit-linters OK;
  `sh docs/tests/run-discipline-tests.sh` — 138 passed, 0 failed.
- `gofmt -l .` — no output; `git diff --check` — clean.

(The focused checks of this row; the full-module run is the root's.)
