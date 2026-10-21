# T-ym0k — row 68a of the plan, the review of the release of M2e

Issue [#369](https://github.com/pharzam/layup/issues/369), row 68a of [the tasks of M2e](../plan/README.md#the-tasks-of-m2e), opened by `T-57vx` (#346). Serves `F-0003#51` through REQ-015 and REQ-017. Goal count 1 (O-214). Evidence: [`runs/T-ym0k/`](../../runs/T-ym0k/).
Plan 6314742187 (Claude Opus 5.5, Plan agent); its plan review 6314912248 (GPT-6.1 Sol on Codex CLI: `approve-with-conditions`, seven conditions, D1 to D5, Q1 to Q4, 400 lines over 12 files against `main` at landing, close-out inside, ADR-0028, cap 2, no panel).

## The release

`01d7da3082333efadb227afda14fd98715edf5c4`, the merge of **PR #466** for **issue #368** (row 67, `T-3xvc`, the last build row): 149 paths and 38,708 lines (140 non-test Go files of `cmd/` and `internal/` out of `testdata/`, 38,297 lines; `go.mod`, 3; the eight files of `internal/catalog/go/`, 408). Since M2d the embedded files changed by +112/-1 (`gates.sh.tmpl` +73/-0, `gates.yml.tmpl` +39/-1). Every run names the full ID.

## Decisions (R10), with their reasons

- **D1.** The release is fixed at the full ID: check (1) of the script passes any name, such as `HEAD`.
- **D2.** Check (3) allows the exact call `exec.Command("sh", script, prdDir)` of `internal/spec/prdlint.go` (the pinned prd-lint, which `packages.md` allows as a process seam), not the file: a `-c` form (m1) passes a file-wide allowance and fails this one.
- **D3.** `RemoteHead` joins the regex and the list of check (7): `BranchHead` calls it with an installation token, and it was a silent omission from the listing, not a failure.
- **D4.** Check (6) counts every `.do(` occurrence and the rows that parse, lists the parsed rows only, and fails on a difference or on no occurrence (condition 2). What it does not prove is stated in the header: the form parses, not that the operation or method fits the call.
- **D5.** The comparator `runs/T-ym0k/coverage-check.sh` stays in the tree (counted in the budget) with its test, so a reader can rerun it; the precedent kept it in scratch. It builds both halves from COMMIT (a temporary index made by `git read-tree`), which is the answer to condition 3. It is a shell script, not a new `layup` command.

## State

S0 to S5 are done by the execution tier; S6 (the paid release review by a reviewer of a model that wrote none of the release, Claude Fable 5.1 first, then the eligible fallback of ADR-0026's order; GPT-6.1 Sol is not eligible: it wrote `go.mod:3`, `7f4e5aec`, T-xxr8) is not run, and condition 4 asks for an independent review of the changed check and the comparator first. Corrections to the plan's statements (condition 7): the embedded diff above; the trailerless blame group also holds `T-3q0b`; Gemini's attempt at `T-2vva` did no work (its authorship rests on the fixes of `T-h65e` and `T-r7b9`); in m3 the `grep` omits the variable-operation call, and the `sed` does not drop it.
Tests: [`test-runs.md`](../../runs/T-ym0k/test-runs.md). The verdict and the resource record are written at the close-out.
