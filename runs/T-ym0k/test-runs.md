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

## Green 1: the final check (11:42:23 to 11:42:30, and 11:42:39 to 11:42:40)

`runs/T-efmy/release-check.sh` is changed (D2 to D4): the allowance of check (3) is the exact call `exec.Command("sh", script, prdDir)` of `internal/spec/prdlint.go`; `RemoteHead` is in the regex and the list of check (7); check (6) counts every `.do(` occurrence (`git grep -o`) and the rows that parse, lists the rows it parsed, and fails on a call in another form, on a line of two calls, and on no occurrence at all. The claim, as the header states it: every `.do(` of the adapter is a row of the form `a.do(ctx, "Operation", "METHOD", target,` on one line, and the list is made of those rows; it does not claim that the operation or method fits what the call does, nor that `do` is the only way to the network.

The clean release and each mutant (`sh run-all.sh final.sh final`), each `FAIL` line read:

| Run | Final |
| --- | ----- |
| the clean release | exit 0, `== PASS` |
| m1 | exit 1, `FAIL … internal/spec/prdlint.go:32: … exec.Command("sh", "-c", script)` (was red against the control) |
| m2 | exit 1, `FAIL … internal/spec/full.go:834` (regression check) |
| m3 | exit 1, `FAIL: a .do( call that is not in the form … internal/forge/github/github.go:771: … a.do(ctx, op, "PUT", …` and `the .do( occurrences (22) and the rows parsed (21) differ` (was red) |
| m4 | exit 1, `FAIL: a reference to git.Push, out of its function: internal/run/push.go:193: Push in mutantPush` (regression check) |
| m5 | exit 1, the same two lines of check (6), for line 771 `a.do(ctx, "Merge", m, …` (was red) |
| m6 | exit 1, `FAIL: a line with two .do( calls: internal/forge/github/github.go:771` and `(23) … (22) differ` (was red) |

The empty adapter (section (6) alone, 11:42:33): exit 1, `FAIL: no .do( call found in internal/forge/github`, `0 .do( occurrences, 0 rows parsed`. The output assertion of D3: the output of the clean release lists `internal/run/store.go:202: RemoteHead in BranchHead` (the control's does not).

## Green 2: the release of `M2e` (11:42:39 to 11:42:40)

`sh runs/T-efmy/release-check.sh 01d7da3082333efadb227afda14fd98715edf5c4` from a clean clone of the commit (Go 1.26.10): exit 0, `== PASS`. Its output is [`release-check.txt`](release-check.txt), 140 lines. Against the output at the release of `M2d` (`runs/T-5xjc/release-check.txt`, 116 lines): the sites of check (3) are the five and the sixth, `internal/spec/prdlint.go:32` (allowed by its exact form); `internal/handoff` and `internal/spec` join the packages of `go list -deps`; check (4) has the same six hits; the exported calls of `internal/git` and the references of check (7) grow (17 references, with `RemoteHead in BranchHead` (the control lists 16), three of `git.Push` each in an allowed function); check (6) lists 21 rows for 21 occurrences (ten at M2d), the seven new writes among them (`OpenDraft`, `MarkReady`, `SetStatus`, `Merge`, `EditBody`, `ClosePullRequest`, `CloseIssue`).

## Green 3: the comparator (11:41:38 to 11:41:39)

`sh runs/T-ym0k/coverage-check-test.sh` with the comparator `runs/T-ym0k/coverage-check.sh`: the 17 cases pass (`all ok`), including the ones that `naive1.sh` and `naive2.sh` missed. `sh runs/T-ym0k/coverage-check.sh --list 01d7da3082333efadb227afda14fd98715edf5c4` from the worktree (HEAD is the task branch, not the release) gives 149 paths and 38,708 lines, the numbers of the plan; the same list, byte for byte, from the clone with `HEAD` at m3.

## S5: the brief and the clone of the reviewer, prepared and not run (11:45:29 to 11:45:31)

The clone is `~/.cache/layup-session/reviews/issue369-release/layup`, a fresh clone checked out at `01d7da3082333efadb227afda14fd98715edf5c4` with `git sparse-checkout init --no-cone` and `set --no-cone '/*' '!/runs/' '!/docs/tasks/' '!/docs/prd/' '!/docs/plan/'`; `.git/info/exclude` lists `.review-in/`, `.gocache/` and `.tmp/`. `git status` is empty, and `sh .review-in/release-check.sh 01d7da3082333efadb227afda14fd98715edf5c4` in it exits 0 with the output of `release-check.txt` but for the time of check (2). `.review-in/` holds `issue.md` (the body of #369, and nothing else), `criteria.md` (the rows REQ-015 and REQ-017 and their acceptance rows of `PRD-0001`, `loop.md` "The review of the release", `packages.md` "The table of M2e", `forge.md` "The calls of M2e", ADR-0015, each verbatim), `release-files.txt` (149 rows, 38,708 lines), the script and its output, `record-form.md` (its Coverage table has the path and the count of each row and an empty `Read` column) and `brief.md` (what may be read; no plan, plan review, answer or earlier verdict). The comparator on the form as it stands: exit 1, 149 lines `FAIL: unparsed row`; on the form with `in full` in each `Read` cell: exit 0, `ok: 149 rows … each in full`. The files are made by a scratch script (`prepare-review.sh`) from the files of this branch, so they are rebuilt from the head that the review reads.

## After round 1 of #369 (12:02:27 to 12:03:14): two holes in the check

Round 1 found that check (3) exempted a whole grep line that held the allowed call, and that check (6) matched `a.do` without a boundary before `a`. Three mutants, each a commit of the scratch clone built first: m7 `cmd := exec.Command("sh", script, prdDir); _ = exec.Command("sh", "-c", script)` (`fb182403`), m8 `cmd := exec.Command("sh", "-c", script) // exec.Command("sh", script, prdDir)` (`39e44ffc`), m9 the receiver of `Installation` and its references renamed from `a` to `aa` (`8757c392`).

**Red (12:02:27 to 12:02:30):** the script of the head `6df1389` and the control each exit 0, `== PASS`, on m7, m8 and m9.

**The fix:** check (3) allows a `prdlint.go` line only when it holds the exact call `exec.Command("sh", script, prdDir)` once, no `//` before it, and no other `exec.Command`, `os.StartProcess` or `syscall.` text in the rest of the line. Check (6) puts an identifier boundary (not a letter, digit, `_` or `.`) before `a.do` in the parsing and the diagnostic expressions.

**Green:** the script over the clean release and m1 to m9 (each `FAIL` line read): the clean release exits 0; m1 to m6 fail on the same lines as before; m7 and m8 exit 1 with `FAIL: a call that starts a program, out of the list: internal/spec/prdlint.go:32:` and the mutated line; m9 exits 1 with `FAIL: a .do( call that is not in the form …: internal/forge/github/github.go:96: … aa.do(ctx, "Installation", …` and `occurrences (21) and the rows parsed (20) differ`. The clean release again (12:03:13 to 12:03:14): exit 0, `== PASS`, 140 lines, 17 references in check (7); `release-check.txt` is that output. The 17 cases of the comparator stay green.

## After round 2 of #369 (12:17:45 to 12:18:29): the receiver `αa`

Round 2 found that the boundary before `a.do` excluded only ASCII letters, digits, `_` and `.`, so a receiver named `αa` (a valid Go identifier) passed as `a.do`. Mutant m10 renames the receiver of `Installation` and its references from `a` to `αa` (`110b2403`, built first).

**Red (12:17:45):** the script of the head `24f3db6` and the control each exit 0, `== PASS`, on m10.

**The fix:** a positive boundary. The byte before `a.do` must be a space, a tab or one of `( , ; = ! & | { :`, or `a.do` starts the line; the parsing `sed` and the diagnostic `grep` use the same set, both under `LC_ALL=C` so that bytes are compared, and the header says so and states what is a refusal.

**Green (12:18:01 to 12:18:20), each `FAIL` line read:** the clean release exit 0; m1 to m9 fail on the same lines as before; m10 exit 1 with `FAIL: a .do( call that is not in the form …: internal/forge/github/github.go:96: … αa.do(ctx, "Installation", …` and `occurrences (21) and the rows parsed (20) differ`.

| Run | Control | Final |
| --- | ------- | ----- |
| clean release | exit 0 | exit 0 |
| m1 `sh-c-in-prdlint` | exit 0 (red) | exit 1, (3) |
| m2 `sh-in-spec-full` | exit 1, (3) | exit 1, (3) |
| m3 `do-variable-op` | exit 0 (red) | exit 1, (6) |
| m4 `push-in-run-push` | exit 1, (7) | exit 1, (7) |
| m5 `do-variable-method` | exit 0 (red) | exit 1, (6) |
| m6 `do-two-on-a-line` | exit 0 (red) | exit 1, (6) |
| m7 `allowed-call-and-sh-c` | exit 0 (red) | exit 1, (3) |
| m8 `sh-c-with-allowed-text-in-a-comment` | exit 0 (red) | exit 1, (3) |
| m9 receiver `aa` | exit 0 (red) | exit 1, (6) |
| m10 receiver `αa` | exit 0 (red) | exit 1, (6) |

The clean release again (12:18:28 to 12:18:29): exit 0, `== PASS`, 140 lines; `release-check.txt` is that output. The 17 cases of the comparator stay green.

## After the review: the comparator accepts a plain `in full` (13:05:30 to 13:05:38)

The review record wrote its `Read` cells as plain `in full`, where the comparator required the backticks (the brief said "with `in full`"). **Red (13:05:30):** two cases added to `coverage-check-test.sh` first, a record whose `Read` cells are plain `in full` must pass (`FAIL … exit 1, wanted 0`: the rows were `unparsed row`), and a plain `partly` must be refused by name (`FAIL … no line with: not `in full`: internal/a/a.go`). **Fix:** the third cell is any non-empty text, with or without backticks, and the `in full` comparison is unchanged. **Green (13:05:31):** the 19 cases pass, the 17 earlier ones included; the unfilled `record-form.md` (an empty `Read` cell) is still 149 `unparsed row` lines.

`sh runs/T-ym0k/coverage-check.sh 01d7da3082333efadb227afda14fd98715edf5c4 /home/layup/.cache/layup-session/reviews/issue369-release/review.md` on the unedited record (13:05:38):

```
ok: 149 rows of /home/layup/.cache/layup-session/reviews/issue369-release/review.md are the release at 01d7da3082333efadb227afda14fd98715edf5c4, each in full
exit 0
```

## The review (the uat): 2026-10-21T12:30:04Z to 12:53:29Z

Claude Fable 5.1 on OpenCode (`anthropic/claude-fable-5-1`, variant `max`), chosen by the reviewer selector among five candidates (GPT-6.1 Sol and Claude Opus 5.5 excluded as author families by blame; Fable first eligible; its probe 12:29:49Z to 12:29:56Z, exit 0, pass: "tool call seen, marker and count present"; [`selection.json`](selection.json)). Exit 0, USD 41.26, 106,552,552 tokens with cache. It ran in the prepared clone at `01d7da3` with `.review-in/`. The `release-check.sh` it read has the sha256 `2713c529cec981f16a8855e9752040893f8631bdbc436b7217703c2861dfb8bf` (the file of this branch at that time, byte for byte). The record is [`release-review.md`](release-review.md), copied unchanged, posted on #369 as comment 6316749856: `holds` for REQ-015 and REQ-017, 11 findings, none a defect.

**The Coverage table against the commit (S7):** first run, exit 1 on form only (the `Read` cells were plain); after the fix above, exit 0 on the unedited record, 149 rows, none missing, extra, out of order or with another count. The table is not edited by the author.
