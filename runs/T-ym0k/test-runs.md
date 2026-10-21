# T-ym0k — test runs

The review of the release of `M2e` (#369, row 68a). The check is `runs/T-efmy/release-check.sh`, adapted (D2 to D4 of the plan, with conditions 1, 2 and 7 of its plan review). The release commit is `01d7da3082333efadb227afda14fd98715edf5c4` (the merge of PR #466, issue #368, row 67). Each run is from a copy of the script outside a scratch clone of that commit (`~/.cache/layup-session/t-ym0k/clone`, Go 1.26.10), as the header of the script asks. Each mutant is its own commit of that clone, built with `go build ./...` first (one attempt at m6 did not build, `errors` was not imported: not run; it was rebuilt and then run). Times are UTC, 2026-10-21, from `date -u`. The red slice below is recorded before any document or script of the branch changed (condition 1).

## The diagnostic (11:36:33 to 11:36:35)

`sh old.sh 01d7da3082333efadb227afda14fd98715edf5c4` in the clean clone (`old.sh` is `runs/T-efmy/release-check.sh` at `01d7da3`): exit 1, on check (3) only: `FAIL: a call that starts a program, out of the list: internal/spec/prdlint.go:32:	cmd := exec.Command("sh", script, prdDir)`. Checks (1), (2), (4) to (7) pass. It is the plan's prediction. Check (6) lists 21 rows for 21 `.do(` call sites; check (7) lists 16 references and has no `RemoteHead`.

## Red 1: the check (11:37:54 to 11:38:14)

The mutants (the commit of each in the clone): m1 `sh-c-in-prdlint` (`exec.Command("sh", "-c", script)` in `internal/spec/prdlint.go`, `fe6c7640`), m2 `sh-in-spec-full` (an `exec.Command("sh", s)` in `internal/spec/full.go`, `e35cb3fb`), m3 `do-variable-op` (a method of `github.go` that calls `a.do(ctx, op, "PUT", …)` with `op` a variable, `4a1bca0d`), m4 `push-in-run-push` (a function of `internal/run/push.go` that refers to `git.Push`, `39846c56`). Two more, for the claim that condition 2 asks to bound: m5 `do-variable-method` (`a.do(ctx, "Merge", m, …)`, `98f0e8b1`) and m6 `do-two-on-a-line` (two calls on one line, `ee174ec4`). The loop is `sh run-all.sh SCRIPT LABEL` (a scratch script; it checks out each commit, runs the script with its full ID and prints the `FAIL` and `== ` lines). Each `FAIL` line was read, not the exit code.

**(i) The script as it stands** (`old.sh`): every run exits 1, and every run has the same line `internal/spec/prdlint.go:32` (check (3)), so it cannot show that anything else bites. m2 adds its own line (`internal/spec/full.go:834`), and m4 adds `internal/run/push.go:193: Push in mutantPush`.

**(ii) The control** (`control.sh`: the old script plus one line, a file-level allowance `internal/spec/prdlint.go: exec.Command("sh", `, in the style of `scratch.go`; kept in the scratch directory):

| Run | Control |
| --- | ------- |
| the clean release | exit 0, `== PASS` |
| m1 | exit 0, `== PASS` (**red**: the allowance is file-wide, so a `-c` form passes) |
| m2 | exit 1, `FAIL … internal/spec/full.go:834` (the allowance is file-narrow; a regression check) |
| m3 | exit 0, `== PASS` (**red**): the grep of check (6) does not match `a.do(ctx, op,`, so the list has 21 rows and the call is not in it |
| m4 | exit 1, `FAIL: a reference to git.Push, out of its function: internal/run/push.go:193: Push in mutantPush` (a regression check of check (7)) |
| m5 | exit 0, `== PASS` (**red**): the grep matches the line and the `sed` does not match it, so the raw source line is printed unchanged in the list (line 771 without the two-space indent): it is neither a row nor a failure |
| m6 | exit 0, `== PASS` (**red**): the greedy `sed` keeps the last call of the line (`PUT`) and drops the first (`GET`); a count of lines would be equal |

**(iii) The listing of D3** (the output assertion of condition 1): the control's output of the clean release has no line `RemoteHead in BranchHead` in check (7) (`grep -c 'RemoteHead in BranchHead' out-control-clean.txt` gives 0; `RemoteHead` appears once, in the list of the exported calls). The final script must list `internal/run/store.go:202: RemoteHead in BranchHead`. It is a listing change, so it has no mutant.

**(iv) The empty adapter** (condition 2, "two empty counts must not establish coverage"; 11:40:14): section (6) of the control run alone, in a scratch repository whose `internal/forge/github/github.go` holds `package github` (`sh sec6.sh control.sh empty6`): exit 0, no row, no `FAIL`. Red.

## Red 2: the comparator, before its code (11:39:50)

`runs/T-ym0k/coverage-check-test.sh` is written first, with `COVERAGE_CHECK=PATH` to test another script. A missing script is not the right red, so two incomplete comparators of the scratch directory were run instead:

- `naive1.sh`: the set of the paths of the release at the commit, without count, order or `Read`, and no guard on an empty list. It misses 9 of the 17 cases: the count off by one, the swapped rows, a `Read` that is not `in full`, a duplicate row, a row of two cells, the empty record table, the missing heading, and both empty-list cases (exit 0 for the empty list with an empty table: **two empty lists establish coverage**).
- `naive2.sh`: exact rows, but read from the working index and tree (what `git ls-files` and `wc -l` give in the clone at the time), without the `Read` check. It misses 13 of the 17: it fails the good record of the commit when `HEAD` is another commit (condition 3: the index and COMMIT are not the same tree), and it cannot name a count, a duplicate, a `Read`, a row of two cells, or an empty list.

Each is a red for the test file: the cases that expect exit 1 or a named line and do not get them are listed in the output (kept in the scratch directory).
