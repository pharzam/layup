# T-4wyg — test runs and evidence

The full spec check (issue #351, row 62c of M2e). This file records the red
run of the test-first order, the green run, the focused checks and the
mutation checks, each with its command, exit and UTC time.

## Red — 2026-10-16T13:04Z, stubs only (assertion failures, no compile errors)

`internal/spec/full.go` and `internal/spec/prdlint.go` carried the final API
with stub bodies (`CheckFull` nil, `ranRow` nil, `WriteSpecCheck` nil,
`RunPRDLint`/`WaitingSpecAccepted`/`waitingVerdict` error/zero), the real
schema values and `parseReqLine`, the `built` move of `task-spec-check`, and
all test files.

`go test ./internal/spec/` — exit 1 (21 failing tests). The failures are the
intended assertions: `no row covers in []` (CheckFull returned nil),
`waitingVerdict gives false, want true`, `line 1: no header row`
(WriteSpecCheck wrote nothing). Passing already: `TestParseReqLine`,
`TestEmptyCellFailureIsAnchored` (real regex and parser), and
`TestTheSchemasEqualTheirBlocks`/`TestTheOptionsSchemaEqualsItsBlock` —
`TaskSpecCheckSchema` and `tasksSchema` equal their blocks of records.md, and
`task-spec-check` sits in `built` (internal/tsv's block test: ok).

`go test -tags=integration ./internal/spec/` — exit 1; the same unit
failures plus the four prd-lint seam tests (real `sh`, real pinned
`docs/prd/prd-lint.sh` on fixtures in `t.TempDir()`):
`TestRunPRDLintAcceptsAWaitingSpecification`,
`TestRunPRDLintRefusesARealViolation`, `TestRunPRDLintReadsBothStreams`,
`TestRunPRDLintEdgeCases` — all fail on the stub's `stub` error, i.e. on the
assertion, not the compile.

`go test -tags=integration ./internal/tsv/` — ok (0.019s): the block list
holds `task-spec-check` in `built`, every block still listed once.

## The producer's real output (the pinned forms the classifier accepts)

`sh docs/prd/prd-lint.sh` on a fixture of one PRD whose MoSCoW and Phase
cells are all empty gives, on **stdout**, exit 1:

```
FAIL  PRD-0001-x.md: REQ-001 MoSCoW "" not in Must|Should|Could|Won't
FAIL  PRD-0001-x.md: REQ-001 has no phase value
```

The filename check goes to **stderr**: `FAIL  PRD-bad.md: filename must be
PRD-NNNN-kebab-case.md` — confirmed by `err()` writing to stderr while awk's
`print` goes to stdout (prd-lint.sh:36 vs :113–141). A badly named file is
not linted: the positional list holds only well-formed names (line 55).
A directory with no `PRD-*.md` prints `prd-lint: OK` and exits 0. Verified
by reading the pinned script; the integration tests assert the same streams.

## Green — 2026-10-16T13:11Z, implementation landed

`go test -count=1 ./internal/spec` — exit 0 at 2026-10-16T13:11:57Z
(`ok … 0.064s`): all ten rules, the modes, `ran`, the writer round-trip and
the pure classifier pass.

`go test -count=1 -tags=integration ./internal/spec ./internal/tsv` —
exit 0 at 2026-10-16T13:11:43Z (`ok internal/spec 0.574s`,
`ok internal/tsv 0.050s`): the real `sh` + real pinned `prd-lint.sh`
fixtures pass — accepted empty-cell specification, rejected real violation,
both streams exercised, bad-filename stderr rejection, empty-directory OK,
`—` MoSCoW rejected, dash/phase edge cases.

`go build ./...` — exit 0; `go vet ./...` — exit 0;
`go test ./...` (unit, no integration tag) — exit 0, all packages, run at
≈13:10Z before the interrupted launch below.

## Interrupted — 2026-10-16T13:11Z (recorded as failure, not a pass)

`go test -tags=integration ./...` — **INTERRUPTED, exit never produced.**
This launched a full-module integration run although the run order reserves
the sole full suite to the coordinator (issue #360 runs it serially once the
source is stable). The duplicate process group 2655167 was SIGTERM-stopped
around 13:11Z and the CLI received SIGINT. The command is not retried here
and must not be read as a green signal of any kind; the focused
package-named runs above are this task's integration evidence.

## Mutation checks — 2026-10-16T13:13Z–13:15Z, on a copy at /tmp/T-4wyg-mut

Each mutant was applied to a copy of the tree (`cp -a`), compile-checked
(`go build ./internal/spec`), and judged by `go test -count=1
./internal/spec` in the copy. "Killed" means the focused suite failed on the
intended assertion; a build that fails is invalid and uncredited — none was.

| Mutant | Compiled | Verdict |
| ------ | -------- | ------- |
| `waitingVerdict` drops the `fails == 0` refusal (vacuous exit-1 accepts) | yes | **killed** — TestWaitingVerdict |
| `waitingVerdict` judges stdout only (stderr stream lost) | yes | **killed** — TestWaitingVerdict |
| `touchesSpec` always true (mode precedence gone) | yes | **killed** — TestCheckFullNoSpecificationYet |
| `inventory` skips the sha256 comparison | yes | **killed** — TestInventoryVerifiesEachVersion |
| `resolveFact` resolves any `F-NNNN#n` (record binding dropped) | yes | **killed** — TestCoversBindsEachTokenToItsRecord, TestQuotesFollowsTheChain |
| `ranRow` treats a did-not-run row as pass | yes | **killed** — TestRanRowJudgesTheTable, TestCheckFullDidNotRunIsNeverAPass (+2) |
| `classes` drops the copy-names-a-missing-requirement direction | yes | **killed** — TestClassesComparesWithTheLatestCopy |
| a `clear` table emits no rows | yes | **killed** — TestCheckFullNoSpecificationYet |
| `needs` omitted from `fullRules` | yes | **killed** — seven tests incl. the literal order assertion |
| `covers`/`needs` swapped in `fullRules` | yes | **survived, then killed**: the first assertion compared rows against `fullRules` itself, so the mutant moved the oracle; `wantTenRows` now asserts the block's literal order and the mutant fails four tests |
| criteria prose feeds the success set via `factTokens` | yes | **survived** — the fixture's criterion token `(F-0001#3)` keeps its parens, which `factTokens` never unglues; a weak mutant |
| the same via `citeTok` (the real threat) | yes | **killed** — TestMeasuresReadsOnlyTheSuccessRows |

## Focused checks — 2026-10-16T13:16Z

`go build ./...` exit 0; `go vet ./...` exit 0; `go test -count=1
./internal/spec` exit 0 at 13:16:54Z; `go test -count=1 -tags=integration
./internal/spec ./internal/tsv` exit 0 (0.341s, 0.028s). The documented
readers: `adr-lint.sh` exit 0; `prd-lint.sh` exit 0; `link-lint.sh` exit 0
(2,732 links); `run-discipline-tests.sh` exit 0 (138/0); `git diff --check`
clean at 13:16:41Z.

## Size against `e16568c`

14 files, 1,898 insertions and 9 deletions = **1,907 changed lines**
(modified lines count twice by the add+remove measure), under the reviewed
maximum of 2,100 lines / 20 files.

## Coordinator resource audit and main merge

Main cfcfc1c was merged at9d500fd before the full integration run. The shared
built-list resolution retains both task-spec-check and retrospectives; the PRD
Task/Test cells retain both changes. The changelog retains #356 then #351.
The first independent code review must explicitly assess the resolutions.

Actual ATIF export metrics are retained outside git (issue351-export-resource-metadata.json
and issue351-resource-union-audit.json under /tmp/layup-m2e-resume/).
Plan reports prompt3186764/completion50035/cached3015704,36steps.
Interrupted build reports retained prompt969388/completion1216/cached935424,
14steps. Final resumed export reports prompt6847466/completion38553/
cached6549504,82steps. Every recorded step model is swe-2-max.
The final export contains all7 interrupted metric rows exactly, but none of the
plan's28 metric rows. Those exports are not three disjoint totals or a complete
monotonic conversation counter: compaction discarded earlier build history.
The first retained build metric is13:07:43Z; earlier build usage is not reported.
The unique retained plan-plus-final rows report prompt10034230/completion88588/
cached9565208 (101metricrows), plus missing portions. Cached is inside prompt,
not added again. The interrupted export is a subset and is not added again.
Elapsed is not reported; original system-message timestamps do not prove a
complete retained history. Private reasoning/transcripts stay outside git.

## Coordinator full-module integration

`go test -count=1 -tags=integration -timeout 25m ./...` passed directly,
all21packagesok, on6d3334e2357a4ac221e9918cfca2bb936e64c7a9 including the
current-main merge. Wrapper start2026-10-16T13:27:44.363674Z,
end2026-10-16T13:39:04.424635Z (680060.961ms), exit0.
Audit: /home/layup/.cache/layup-session/audits/issue351/integration-result.json.
Only this run-evidence entry follows; no Go or spec changes after that run.

## Round 1 fixes — 2026-10-16T14:12Z–14:15Z (cycle 0 of cap 1)

Review round 1 (comment 6217333827, a fresh Claude Sonnet 5.5 session) gave
`not mergeable, findings recorded`: F1–F4 material, notes N1–N5. The author's
answer (6217372050) accepted each before the edits. The coordinator's
integration entry above predates these fixes; a new head freezes for round 2.

- **F1** — red: `TestAShortCriteriaRowIsNoPanicAndNoCriterion` panics at
  `full.go:376` (`index out of range`), exit 1 at 14:12:45Z. Fix: the
  `criteria` case takes a `len(cells) >= 2` guard, so a truncated row records
  no criterion and the rule fails. Green at 14:12:55Z.
- **F3** — the two exemptions get their own tests
  (`TestMeasuresExemptsAnOutOfScopeFact`, `TestMeasuresExemptsAWaitingFact`),
  green on the correct code at 14:13:32Z. On the copy, the mutant
  `class == "measure" && !cited[f]` (both guards removed) compiles and is
  **killed**: both tests fail, 14:14:23Z; restored green.
- **N2** — a missing version's bytes is a `did not run: version <n>` (the
  `CheckFull` contract), and `ran` fails it; asserted at 14:14:09Z.
- **F2** — `docs/tests/traceability.md` gains the three rows (the full-check
  unit tests, the prd-lint unit tests, the real-`sh` integration tests with
  the `task-spec-check` and second-`tasks` schema compares).
- **F4** — the `packages.md` M2e bullet names the second Go value of `tasks`
  (C11's "and").
- **N1** — `loop.md` states the `waits for a bet` match (case-insensitive,
  any level, to the next heading, the citation token forms). **N4** — the task
  file restates the conditions under the review's own numbers.
- **N5** — declined: an empty `Index` is the documented pre-bet signal, a
  caller boundary the wiring task keeps.

Focused green after the fixes: `go test -count=1 ./internal/spec` exit 0 and
`go test -count=1 -tags=integration ./internal/spec` exit 0 at 14:14:09Z.

## Second merge of main (before round 2)

`origin/main` advanced to `f86f71c` (PR #388, `T-dx1h`) while the round-1
fixes landed; the branch takes it by merge, never rebase. Three files
conflicted and each was resolved by union — `docs/prd/PRD-0001-layup.md`
(the REQ-012 row keeps both `T-dx1h` and `T-4wyg` task cells and both test
lists), `docs/spec/records.md` (the `transitions` clause of `T-dx1h` kept,
the `clear`-table sentence of this task kept), and
`internal/tsv/blocks_integration_test.go` (`built` holds both `parameters`
and `task-spec-check`, `notYetBuilt` neither). Merge commit `64677ec`
("chore: T-4wyg merge current main f86f71c (T-dx1h) before round 2").

Post-merge focused verification, all exit 0 at 2026-10-16T14:23:57Z–14:24Z:
`go build ./...`, `go vet ./...`, `go test -count=1 ./internal/spec`
(0.070s), `go test -count=1 -tags=integration ./internal/spec
./internal/tsv`, and `git diff --check`; the commit's hooks (adr-lint,
prd-lint, run-discipline-tests 138/0, link-lint 2764 links, setup-check
adapted/markers) all pass. No full-module integration was run by this
author — that stays with the coordinator (#360).

Size against the new merge-base `f86f71c` (`git diff origin/main...HEAD`):
2,029 insertions + 10 deletions = **2,039 changed lines over 16 files**,
inside the approved 2,100-line / 20-file maximum.

## Coordinator resource audit before the repaired full-module run

SWE-2 Max fix export retains30 metric rows from14:19:02 to14:25:20UTC:
prompt2318995/completion10270/cached2266176, cached inside prompt. Those rows
are disjoint from the retained plan/build rows. Known union of131 metric rows:
prompt12353225/completion98858/cached11831384; cache is not added twice.
The export compacted earlier fix work before14:19 as well as the previously
reported missing build history; that usage is unavailable, not estimated.
Elapsed is not reported. Source: issue351-plan-build-fix-resource-union.json,
kept outside git with the private exports. The coordinator owns #351's repaired
full-module run; references above to coordinator #360 identify the serialized
slot previously occupied, not a claim that #360's run verifies this fix.
