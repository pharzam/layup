# T-3q0b — the runs of the tests

Times from the clock of the host, UTC, 2026-10-16.

## Round-1 findings: red for the right reason

On the frozen head `72528bd`, 07:43:21Z–07:43:26Z, `go test -count=1 -run
'<the five tests>' ./internal/spec/` (direct exit 1):

- `TestCheckQuoteFacts` — `the answer's fact before the quote's fact: findings
  [{covers the row R-001: covers an answer's fact Q-001 and not the fact 1 of
  that answer's quote}]` (finding 2, the unordered covers).
- `TestTheOptionLinesAndTheVerdict` — `a decision with 1 survivors asks no
  option` and `with 0 survivors` each get `the decision stack has more than
  one survivor and no option line` (finding 3).
- `TestRenderPRDChangesOnlyThePermittedCells` — the portion reads `Out of
  scope by the bet: 5.`, not the union `2, 5, 10, Q-001, Q-002.` (finding 4).
- `TestTheAppetiteAgainstTheBand` — `the amount "1/2"`, `"0x10"`, `"1e2"`,
  `"+5"`, `"1_0"` each give no finding (finding 5; `"-5"` the sign guard
  already refused).
- `TestInventoryRefusals` — `panic: runtime error: index out of range [0]
  with length 0` at `inventory.go:108` (finding 7).

## Round-1 findings: green and the mutants of the fixes

07:44:52Z after the fixes, `go build ./...`, `go vet ./...` and `go test
-count=1 ./...` — all 21 packages `ok`.

A first mutation attempt before the recorded sweep was **discarded and is
not evidence**: its harness restored stale pre-fix copies from
`/tmp/*.go.bak` and silently reverted the five fixed files. The revert was
caught by reading the restored source, the post-fix `.bak2` copies (taken
07:45:29Z) were restored, the focused tests were re-run green, and the
recorded sweep then ran with backups taken fresh from the fixed files at
07:46:15Z (the `/tmp/*.go.bak` mtimes). The valid sweep — applied one at a
time, each restored after, 07:46:15Z–07:47:05Z, direct exits — gave the
six outcomes below: five assertion failures and the one expected panic the
last row names. It is not an "all assertions" sweep:

| Patch | Caught by |
|---|---|
| covers judged while filling | TestCheckQuoteFacts |
| the multi-survivor gate dropped | TestTheOptionLinesAndTheVerdict |
| the prior marks dropped | TestRenderPRDChangesOnlyThePermittedCells |
| the decimal grammar open | TestTheAppetiteAgainstTheBand |
| the empty-row guard off | TestInventoryRefusals (panic, the test aborts) |
| the index order check inverted | TestInventoryRefusals, TestTheNextInventory |

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

The mutants were applied one at a time to a copy of the one file and `go test
-count=1 ./internal/spec/` ran against each, the file restored after. The
author's run times recorded here earlier were wrong; the citable run is the
root's independent deterministic replay of the author's own stored mutation
calls at head `2e1e9b4`, **07:02:30Z–07:03:37Z** (audit
`.cache/layup-session/audits/T-3q0b.gjcAy0/`, `summary.json`):

- **68 unique patches** replayed.
- **59 caught-assertion**: the test binary built, ran and an assertion
  failed, with the killers below.
- **5 invalid-no-assertion**: build-broken re-expressions (an unused variable
  or a malformed splice, `go test` exit 1 with `[build failed]`, so no
  assertion evidence). Each covers a rule a caught re-expression then killed:
  `if r[c] == "" {` → `"?"`, the index-order `want+"?"` splice → `&&` for
  `||`, `if ss <= start {` → `end <= se+99`, `return nums[i] < nums[j]` →
  `x > y`, `findings = findings` → `if m[0] == "" { findings = findings }`.
  They are not counted as kills.
- **4 not-applicable**: `old fragment absent` — the patch never applied
  (harness assertion on the old text).
- **Zero valid survivors.**

Each table row names the rule broken and the replay's catching test.

`bet.go` — 17 caught:

| Patch | Caught by |
|---|---|
| limit accepts any | TestTheLineLimit |
| link map false means resolve | TestTheLinks |
| no-link finding dropped | TestTheLinks |
| batch compare dropped | TestTheApproveBatchLines |
| malformed appetite line passes | TestTheFixedLines |
| malformed arch line passes | TestTheFixedLines |
| malformed batch line passes | TestTheApproveBatchLines |
| money sum unchecked | TestTheAppetiteAgainstTheBand |
| hours sum unchecked | TestTheAppetiteAgainstTheBand |
| milestone-0 sum unchecked | TestTheAppetiteAgainstTheBand |
| unset brief.lines accepted | TestTheLineLimit |
| unknown heading accepted | TestEachHeadingIsRequiredOnce |
| ref link not refused (two re-expressions) | TestTheLinks |
| last line not counted | TestTheLineLimit |
| appetite section not required | TestTheFixedLines |
| arch section not required | TestTheFixedLines |

`answer.go` — 13 caught:

| Patch | Caught by |
|---|---|
| Won't phase rule dropped | TestTheRequirementLines |
| Won't needs — dropped | TestTheRequirementLines |
| phase tag unchecked | TestTheRequirementLines |
| dup requirement accepted | TestTheRequirementLines |
| unknown fact id accepted | TestTheFactLines |
| verdict not required | TestAnAnswerThatFollowsTheFormParses, TestTheRequirementLines, TestTheFactLines, TestTheOptionLinesAndTheVerdict |
| option survivor check dropped | TestTheOptionLinesAndTheVerdict |
| post-verdict line accepted | TestTheOptionLinesAndTheVerdict |
| option requirement dropped | TestTheOptionLinesAndTheVerdict |
| required req not checked | TestAnAnswerThatFollowsTheFormParses, TestTheRequirementLines, TestTheFactLines, TestTheOptionLinesAndTheVerdict |
| unknown req accepted | TestAnAnswerThatFollowsTheFormParses, TestTheRequirementLines, TestTheFactLines, TestTheOptionLinesAndTheVerdict |
| bet-gate on options dropped | TestTheOptionLinesAndTheVerdict |
| no-bet verdict unparsed | TestTheOptionLinesAndTheVerdict |

`inventory.go` — 15 caught:

| Patch | Caught by |
|---|---|
| copy 0 accepted | TestTheInventoryReaders |
| fact class unchecked | TestInventoryRefusals |
| out mark not applied | TestTheFirstInventory, TestTheNextInventory |
| version n miscounts | TestTheFirstInventory, TestTheNextInventory, TestInventoryRefusals |
| empty column accepted | TestTheInventoryReaders |
| index empty accepted | TestTheInventoryReaders |
| copyID not checked | TestInventoryRefusals |
| milestone not checked | TestInventoryRefusals |
| prior version not required | TestInventoryRefusals |
| answer fact change unchecked | TestInventoryRefusals |
| index n order unchecked | TestInventoryRefusals |
| answer class unchecked | TestInventoryRefusals |
| hash of wrong bytes | TestTheFirstInventory, TestTheNextInventory |
| prior rows not carried | TestTheNextInventory |
| prior index rows dropped | TestTheNextInventory |

`quote.go` — 6 caught:

| Patch | Caught by |
|---|---|
| quote bounds dropped | TestQuoteFactOfFindsTheOneContainingFact |
| paired offsets not required | TestCheckQuoteFacts |
| follow-up chain not walked | TestCheckQuoteFacts |
| answered check dropped | TestCheckQuoteFacts |
| quote fact cover not required | TestCheckQuoteFacts |
| containment one-sided | TestQuoteFactOfFindsTheOneContainingFact |

`render.go` — 8 caught:

| Patch | Caught by |
|---|---|
| non-6/7 rows edited | TestRenderPRDChangesOnlyThePermittedCells, TestRenderPRDRefusals |
| phase — for Won't dropped | TestRenderPRDChangesOnlyThePermittedCells |
| no-marker PRD corrupted | TestRenderPRDChangesOnlyThePermittedCells |
| out-of-scope suffix dropped | TestRenderFactsOfTheVersion |
| numeric sort inverted | TestRenderPRDChangesOnlyThePermittedCells |
| q sort inverted | TestRenderPRDChangesOnlyThePermittedCells |
| no row not found | TestRenderPRDRefusals |
| two rows accepted | TestRenderPRDRefusals |

## Green and the validation, at the content of the head

On the focused content of the implementation commit `f77762a`
(2026-10-16T06:51:43Z), during the session before the commit — and
re-run after the merge of `f12da6a` at `2e1e9b4` (06:52:38Z):

- `go test -count=1 ./internal/spec/ ./internal/records/ ./internal/run/` —
  all pass.
- `go build ./...`, `go vet ./...` — clean.
- `go test -count=1 ./...` — all 21 test packages pass.
- `sh docs/adr/adr-lint.sh` — OK; `sh docs/prd/prd-lint.sh` — OK;
  `sh docs/links/link-lint.sh` — OK, 2,717 links resolved after the merge;
  `sh docs/setup/setup-check.sh` — kit-linters OK;
  `sh docs/tests/run-discipline-tests.sh` — 138 passed, 0 failed.
- `gofmt -l .` — no output; `git diff --check` — clean.

(The focused checks of this row.)

## The full run

An author's attempt at `go test -count=1 -tags=integration -timeout 25m ./...`
was interrupted before it finished, and its harness piped the test through
`tail` so its `exit=$?` would have been `tail`'s status, not the suite's; it
proves nothing and is claimed as nothing. The citable full run is the root's
own subprocess capture in the audit directory
`.cache/layup-session/audits/T-3q0b.gjcAy0/` (`integration-result.json`), with
the measured start, end and exit of the run:

- Command: `go test -count=1 -tags=integration -timeout 25m ./...`,
  `TMPDIR=/home/layup/.cache/layup-session/tmp-green`.
- Head: `2e1e9b495d88ddae3c66e4661ad76101421b64b5`.
- Start `2026-10-16T07:08:36Z`, end `2026-10-16T07:27:36Z`, **exit 0** —
  all 21 packages `ok`, none `FAIL` (`internal/run` the long one,
  1,087 s; `internal/verify` 79 s).

After the round-1 fixes the root ran the same command, unchanged, on the
new frozen head `db3330f492e6e90ad4fe086a71a097253a6ee408` (audit
`.cache/layup-session/audits/T-3q0b-r1-fixes/`, `integration-result.json`):
start `2026-10-16T07:51:08.801645Z`, end `2026-10-16T08:03:48.871209Z`,
**exit 0** — all 21 packages `ok`, none `FAIL` (`internal/run` 705 s,
`internal/verify` 65 s).

The merge of `8e26eca` in `db3330f` hit one same-date PRD changelog
append; both dated entries were kept. An earlier report of this branch
described that keep-both resolution as "per O219"; that was wrong — the
coordinator's ruling (issue #350, comment 6210747940) is that the
changelog conflict is not one of O-219's named PRD §12 Task/Test-cell
exceptions, it is a proposed resolution only, and round 2 reviews it
explicitly. No waiver is claimed.
