# The test runs of T-9et9 (#352, row 63a of M2e, the basic layup/rules)

Times are UTC, read from `date -u` around each command; the exit code is the command's own (no pipe into another command). The branch is `T-9et9` on `cfcfc1c`; the integration commands name `./internal/rules/` or `./cmd/layup/ -run TestPackageRules`, and no unfiltered module-wide integration or end-to-end command was run (the root runs the module's integration level).

## Red: the tests against a stub (before any rule)

`internal/rules/basic.go` was first a stub: the types, the adapter over `internal/git`, `Judge` and `Basic` that give `success` with `no rule path`, and `Render` that gives nothing, so that the tests compile and fail by assertion (a compile error is not the red).

| Level | Start | End | Exit | Result |
| ----- | ----- | --- | ---- | ------ |
| `go test ./internal/rules/` | 13:36:05Z | 13:36:08Z | 1 | 19 of the 21 new unit tests fail; no panic; each failing line is an assertion such as `a link with the bytes: {State:success Reason:no rule path Path:}, want failure "not a regular file" "docs/facts/numbered.md"` |
| `go test -tags=integration ./internal/rules/` | 13:37:27Z | 13:37:32Z | 1 | all 10 new integration tests fail by assertion (the same stub answer), with the 19 unit tests |

The two unit tests that pass against the stub pass for a stated reason: `TestJudgeANoRulePathChangeIsASuccess` expects the stub's own answer (the mutant M25 is among the tests it fails, below), and `TestBasicReadsNoCopyAndChecksNoHash` reads the source text, which the stub does not mention (no mutant of the table touches it; a copy of the tree with the line `// the sha256 of a copy` appended to `basic.go` fails it with `basic.go mentions "sha256"`, and passes without the line). The old tests (`TestReadRegister`, `TestMatch`, `TestGuardrailsAdditions`, `TestCheck`, `TestTheSchemaEqualsItsBlock`) pass throughout and `rules.go`, `rules_test.go` and `rules_integration_test.go` are unchanged except the comment lines of the package comment of `rules.go`.

## Green, and the three corrections of the tests before the first commit

1. 13:38:11Z, the unit run failed once: the case `docs/facts/:x` of `TestJudgeBadInputIsAnErrorNeverASuccess` asked for a refusal that D7 does not give (D7 refuses a path that **starts** with `:`). The case was removed; the source was right.
2. 13:38:20Z to 13:38:26Z, the integration run failed twice, both defects of the test helper: `git add -A` undid the entries made by `update-index --cacheinfo` (a link, an executable, a submodule), and a case used the base as both the head and the register's commit, which the code refuses on purpose. The helper commits the index as it is after a `cache` call, and the case gives the register its own commit. The source was right.
3. Both then passed (the exits were read, the times not: the run of the table below is the timed one).

The commit `feat: T-9et9 the basic layup/rules` went through the pre-commit hook (the lints, `gofmt`, `go vet`, the unit level): 13:39:54Z to 13:40:38Z, exit 0.

## The compiled mutants (`runs/T-9et9/mutate.py`)

Each mutant is one exact text replacement of `internal/rules/basic.go` on a **copy** of the tree (the worktree is never mutated); the control, the unmutated copy, is green before and after; the file is written back from the bytes read before the run and its SHA-256 is checked after each mutant. A mutant counts only when `go vet -tags=integration` is clean; a build or vet failure is "invalid" and never a kill; a kill is a failing assertion, and the table names the first failing test of the run (the output of each run was read in full, and is not kept).

**Run 1** (13:40:51Z to 13:42:17Z, exit 0): 33 killed, 2 survived, 3 invalid of 38. Invalid: M15 (the `path` import became unused), M29 and M33 (a variable became unused); each mutant was rewritten with `&& false` to keep the code building. M17 was also mis-composed (it dropped the `..` check too, and was "killed" by `../x`); it was rewritten to change only the dot check. Run 1 counts nothing for those four.

**Run 2** (13:42:31Z to 13:43:56Z, exit 0, controls green before and after): **35 killed, 3 survived, 0 invalid of 38.**

| # | The mutant | Verdict |
| - | ---------- | ------- |
| M01 | the register read at the head | killed by unit: `TestBasicReadsTheRegisterAtTheBase` (+2 more) |
| M02 | the empty register refusal removed | killed by unit: `TestBasicRefusesARegisterItCannotUse` |
| M03 | equality after TrimSpace | killed by unit: `TestJudgeARecordThatDiffersByAByteFails` (+1 more) |
| M04 | equality on a prefix | killed by unit: `TestJudgeARecordThatDiffersByAByteFails` (+2 more) |
| M05 | the mode and type check removed | killed by unit: `TestJudgeARecordEntryThatIsNoRegularFileFails` (+1 more) |
| M06 | the type alone dropped from that check | survived |
| M07 | an absent entry not refused as absent | killed by unit: `TestJudgeARecordEntryThatIsNoRegularFileFails` (+1 more) |
| M08 | return at the first passing record | killed by unit: `TestJudgeAPassingRecordDoesNotHideALaterRulePath` (+1 more) |
| M09 | the first record covers any path | killed by unit: `TestJudgeAPassingRecordDoesNotHideALaterRulePath` (+2 more) |
| M10 | a rendering finding ignored | killed by unit: `TestJudgeARenderingThatFailsIsNoPass` |
| M11 | the PRD pre-image read at the head | killed by unit: `TestBasicPreImageAndBaseChains` |
| M12 | the base ancestor check removed | killed by unit: `TestBasicPreImageAndBaseChains` |
| M13 | the chain of the pre-image removed | killed by unit: `TestBasicPreImageAndBaseChains` |
| M14 | the pre-image may be the head | killed by unit: `TestBasicPreImageAndBaseChains` |
| M15 | a clean-form check removed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M16 | a leading slash allowed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M17 | a dot part allowed | survived |
| M18 | a dot-dot part allowed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` (+2 more) |
| M19 | a .git part allowed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M20 | a backslash or NUL allowed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M21 | a leading colon allowed | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M22 | an error mapped to a success | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` (+2 more) |
| M23 | the base error mapped to a success | killed by unit: `TestBasicPreImageAndBaseChains` |
| M24 | the guardrails exception without its diff | killed by unit: `TestBasicJudgesWhatTheHeadHolds` |
| M25 | the exception never given to Check | killed by unit: `TestJudgeANoRulePathChangeIsASuccess` (+2 more) |
| M26 | the diff of the guardrails taken backwards | killed by integration: `TestTheGuardrailsExceptionInARealRepo` |
| M27 | the diff names taken backwards | survived |
| M28 | the register's commit may be the head | killed by unit: `TestBasicRefusesARegisterItCannotUse` |
| M29 | a bad path read before it is refused | killed by unit: `TestBasicRefusesABadRecordBeforeItReadsAnything` |
| M30 | a PRD with no pre-image accepted | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M31 | a record twice accepted | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M32 | the workflow reason not a rule path | killed by unit: `TestJudgeAWorkflowPathFailsAsARulePath` |
| M33 | a register that is no regular file accepted | killed by unit: `TestBasicRefusesARegisterItCannotUse` |
| M34 | the reason of a rendered record | killed by unit: `TestJudgeAnEqualRenderedFactsRecordPasses` (+2 more) |
| M35 | the record's bytes never read | killed by unit: `TestBasicJudgesWhatTheHeadHolds` (+1 more) |
| M36 | the guardrails head never read | killed by unit: `TestBasicJudgesWhatTheHeadHolds` |
| M37 | the kind not checked | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |
| M38 | the guardrails path may be a record | killed by unit: `TestJudgeBadInputIsAnErrorNeverASuccess` |

The three survivors are equivalent mutants, recorded and not credited: **M06** (the type `blob` is implied by the mode `100644`: git lists a mode other than `100644` for every entry that is no blob, a submodule being `160000`); **M17** (`path.Clean` already removes a `.` part, so a path with one differs from its clean form and is refused before the part check); **M27** (`git diff --name-only` gives the same names and order for the two directions, so the swap changes nothing; the same swap on `DiffFile` is killed by M26). No other survivor, no mutant repeated beyond the four rewritten ones.

## The timed green run

| Step | Start | End | Exit | Command |
| ---- | ----- | --- | ---- | ------- |
| build | 2026-10-16T13:44:17Z | 2026-10-16T13:44:19Z | 0 | `go build ./...` |
| vet | 2026-10-16T13:44:19Z | 2026-10-16T13:44:20Z | 0 | `go vet ./...` |
| gofmt | 2026-10-16T13:44:20Z | 2026-10-16T13:44:20Z | 0 | `sh -c test -z "$(gofmt -l cmd internal)"` |
| unit-all | 2026-10-16T13:44:20Z | 2026-10-16T13:44:28Z | 0 | `go test -count=1 ./...` |
| int-rules | 2026-10-16T13:44:28Z | 2026-10-16T13:44:34Z | 0 | `go test -count=1 -tags=integration -timeout 10m ./internal/rules/` |
| int-pkgrules | 2026-10-16T13:44:34Z | 2026-10-16T13:44:36Z | 0 | `go test -count=1 -tags=integration -timeout 10m -run TestPackageRules ./cmd/layup/` |
| adr | 2026-10-16T13:44:36Z | 2026-10-16T13:44:38Z | 0 | `sh docs/adr/adr-lint.sh` |
| prd | 2026-10-16T13:44:38Z | 2026-10-16T13:44:38Z | 0 | `sh docs/prd/prd-lint.sh` |
| link | 2026-10-16T13:44:38Z | 2026-10-16T13:44:54Z | 0 | `sh docs/links/link-lint.sh` |
| discipline | 2026-10-16T13:44:54Z | 2026-10-16T13:45:02Z | 0 | `sh docs/tests/run-discipline-tests.sh` |
| setup | 2026-10-16T13:45:02Z | 2026-10-16T13:45:33Z | 0 | `sh docs/setup/setup-check.sh` |
| diffcheck | 2026-10-16T13:45:33Z | 2026-10-16T13:45:33Z | 0 | `git diff --check` |

The head of the run is the commit above plus the uncommitted documents (`loop.md`, `packages.md`, `PRD-0001`, `traceability.md`, the glossary, `guardrails.md`) and this directory; the documents were committed after it.

## The measurements behind the issues opened by this task

On a scratch copy of the tree, not in a branch (2026-10-16 about 13:44Z): `ReadRegister` of a header-only file gives 0 entries and no error, and `Check` over it gives no refusal for `AGENTS.md`, `.github/CODEOWNERS`, `docs/tests/test-levels.md` and `go.mod` (#390). `Match` gives no entry for `docs/facts`, `.github` and `docs/ci`, and `Basic` over the real `git` gives `success`, `no rule path` for a link at `docs/facts`, `.github` or `docs/tests` when the base lacks the directory (#391). Reading the register twice, by `push.go` and by `Basic` (#392).

## Coordinator full integration and pre-review merge

Full-module `go test -count=1 -tags=integration -timeout 25m ./...` at
f5693500d5a9d76b202a24fa54f2d7b6d89f18d7 passed: direct exit0,21 packages;
2026-10-16T14:12:19.645234Z–14:26:39.749756Z,860104.522ms.
Current main f86f71c (#358) then merged without textual conflicts. First blind
code review must assess semantic agreement of all five overlapping documents:
guardrails, PRD, loop.md, packages.md and traceability. Final combined full-SHA
CI is required; the earlier integration is not claimed to cover new main Go.

Actual Sonnet5.5 CLI metadata, effort high, no other model reported:
plan I68/O69116/cache creation201709/read4590433,585946ms; build increment
I196/O101557/cache creation172629/read31070610. Thinking34887 in the build
increment is already included in output. The build's main turn reports909884ms,
and its completed-monitor notification reports2131ms separately; no idle is
inferred. End cumulative counters count plan and build once, not twice.
SWE-2 Max independent plan-review retained metrics: prompt2043525/output35830/
cached1801766,32 unique rows; the initial export is a subset of the corrected
export and counted once. Cache is inside prompt. Elapsed is not reported.
Metadata audits issue352-resource-metadata-complete.json and
issue352-plan-review-resource-union.json stay outside git with private exports.
