# T-5xjc — test runs

The review of the release of `M2d` (#283, row 59a). The check is `runs/T-efmy/release-check.sh`, adapted (D3 of the plan, with conditions 1 and 4 of its plan review). Each run is from a copy of the script outside the clone, as its header asks, in a scratch clone of the release commit `0a104b8fc710f45cfe3e73778518616374bae1db`. Each mutant is its own commit of that clone, built with `go build ./...` before its run, so that each exit 1 has its own reason. Times are UTC, 2026-10-14, from `date -u`.

## A diagnostic run before the plan (09:36:31Z to 09:36:38Z)

The script as it stands, on the release commit, to learn whether it needed an adaptation: exit 1, on check (4) only (`FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:452: remote in AddRemote`); checks (1), (2), (3), (5) and (6) pass. It is the reason for the allowance of `AddRemote`. It is not the red slice below: the plan said that the script as it stands would pass the two push mutants, and that was wrong (condition 1 of the plan review: the same line fails them too).

## Red (09:48:59Z to 09:49:31Z, after the author's answer to the plan review)

**(i) The script as it stands** (the file at the base, `0a104b8`): the clean release and the two push mutants all exit 1 on the same line, so it cannot show that anything else does not bite.

```
orig clean: exit 1
   FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:452: remote in AddRemote
orig push-in-batch: exit 1
   FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:452: remote in AddRemote
orig push-in-forgecheckstep: exit 1
   FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:452: remote in AddRemote
```

**(ii) The control**: the script as it stands with one change, the allowance `remote in AddRemote` in check (4), and nothing else (a one-line `sed` of the script, kept in the scratch directory and not in the tree). The clean release passes it. Four mutants pass it as well: each is **red**, as the control has no check (7). The two verb mutants fail it on check (4), which shows that the allowance is narrow.

| Mutant | What it is | Control |
| ------ | ---------- | ------- |
| (clean release) | none | exit 0, `== PASS` |
| `push-in-batch` | `var _ = git.Push` in `internal/run/batch.go`, with its import | exit 0, `== PASS` (red) |
| `push-in-forgecheckstep` | the probe's push moved from `pushProbe` to `forgeCheckStep` | exit 0, `== PASS` (red) |
| `push-in-store-other-func` | a function `pushElsewhere` of `internal/run/store.go` that calls `git.Push` | exit 0, `== PASS` (red) |
| `alias-import` | `internal/run/aliasmut.go` imports `internal/git` as `gitx` and refers to `gitx.Push` | exit 0, `== PASS` (red) |
| `fetch-in-addremote` | a `"fetch"` literal in `AddRemote` | exit 1, `internal/git/git.go:453: fetch in AddRemote` |
| `remote-in-newfunc` | a function `SetRemoteURL` of `internal/git` that runs `git remote set-url` | exit 1, `internal/git/git.go:540: remote in SetRemoteURL` |

## The final script (09:50:27Z to 09:50:48Z)

Check (7) is written (09:49:31Z to 09:50:20Z), and the allowance is in check (4). The clean release exits 0 and each mutant exits 1, each on its own line:

```
final clean: exit 0, == PASS
final push-in-batch: exit 1
   FAIL: a reference to git.Push, out of its function: internal/run/batch.go:354: Push in no function
final push-in-forgecheckstep: exit 1
   FAIL: a reference to git.Push, out of its function: internal/run/forgecheck.go:403: Push in forgeCheckStep
final push-in-store-other-func: exit 1
   FAIL: a reference to git.Push, out of its function: internal/run/store.go:322: Push in pushElsewhere
final fetch-in-addremote: exit 1
   FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:453: fetch in AddRemote
final remote-in-newfunc: exit 1
   FAIL: a git verb that reaches a remote, out of its function: internal/git/git.go:540: remote in SetRemoteURL
final alias-import: exit 1
   FAIL: internal/git imported under a name other than git: internal/run/aliasmut.go:3:import gitx "github.com/pharzam/layup/internal/git"
```

The limits that the header of the script states, and that no mutant here shows: a second call of `git.Push` inside `PushHead`, `commitFiles` or `pushProbe` passes check (7), so the review reads those three functions in full; check (7) reads the selector `git.Push` of the package, and a push by a program that is started with the text `git push` is check (3)'s.

## Green: the release of `M2d` (09:51:00Z to 09:51:03Z)

`sh runs/T-efmy/release-check.sh 0a104b8fc710f45cfe3e73778518616374bae1db` in a clean clone: exit 0, `== PASS`. Its output is [`release-check.txt`](release-check.txt), 116 lines. Against the output at the release of `M2c` (`runs/T-8byd/release-check.txt`): the five sites of check (3) are the same; check (4) has one more hit, `remote in AddRemote`; the exported calls of `internal/git` gain `CommitTree` and `AddRemote`; check (6) gains `EffectiveRules` and `Activity`; and check (7) is new, with 14 references (three of `git.Push`, each in an allowed function). The same run in the clone of the reviewer, which has `runs/`, `docs/tasks/`, `docs/prd/` and `docs/plan/` out of its working tree by a sparse checkout (09:52:35Z to 09:53:15Z, with an empty Go cache): exit 0, the same output but for the time of check (2), and `git status` empty.

## The review (the uat): 09:53:22Z to 10:09:57Z

A fresh read-only session of GPT-6.1 Sol on Codex CLI (`codex exec`, effort `xhigh`, 16 min 35 s, 443,462 tokens, `model: gpt-6.1-sol` in its log, 194 commands, none denied) in a clone of the release commit with `runs/`, `docs/tasks/`, `docs/prd/` and `docs/plan/` out of its working tree (a sparse checkout; `git status` empty). Its brief was the issue, the criteria (verbatim extracts of `PRD-0001`, `scaffold.md`, `session.md` and ADR-0015), the list of the release (90 paths with their line counts), the script and its output, and a blank record form: no plan, no plan review, no answer and no earlier verdict (condition 2 of the plan review). The script that it read has the sha256 `7a24759e8a62c1647666759746f62d5f872ac4c877cb3d2db8c21d17812d65cd`, the file of this branch at that time; the line numbers of the record's citations of the script are those of that file.

It recorded `holds` for REQ-015 and for REQ-017, with 29 findings, none of which names a defect: findings 19 to 21, 23 and 25 to 29 are notes, each closed with "No change" ([`release-review.md`](release-review.md), posted on #283 as comment 6165976073). It read the whole release in chunks of at most 250 lines, ran the script again (exit 0, `TestPackageRules` 2.419 s), `go build ./...` and `go vet ./...` (each exit 0), and searched for each thing that its two rows name.

**The Coverage section against the commit (condition 3):** the 90 rows of the record are the 90 paths of the release in the order of the list of the brief (the 81 non-test Go files out of `testdata/`, `go.mod` and the eight files of `internal/catalog/go/`, from `git ls-files`), each with the line count that `wc -l` gives at the commit and the word `in full`: 90 of 90, no row missing, extra, out of order or with another count. The comparison is a script of the scratch directory, not part of the tree.

**The size of a chunk:** the brief said "about 250 lines at most", where the answer to the plan review said 200 (the author wrote 250 in the brief and did not note the difference). The aim of condition 3 held: the record says that each cut output was read again, and its 90 rows equal the commit.

**Sessions of the same model on this release:** the task files name GPT-6.1 Sol as a reviewer in the plan reviews and rounds of `T-4c5c`, `T-56dj`, `T-9v92`, `T-j4sd`, `T-jpkc`, `T-eahr`, `T-79zk`, `T-hfn1`, `T-wphp`, `T-2j13` and `T-eevt`, and in the plan review of this task (note 4 of it already reads the probe as a write of content). That is review, not a commit of the release, and the gate defines an author as the model whose session wrote the change (Who may review). The fresh session could not know it; the records of `T-4tjy` and `T-8byd` said the same of their reviewer's model.

**A miscount in the record, not a finding:** finding 15 says ten non-push references and its table of runs says thirteen references outside `internal/git`; check (7) lists fourteen, eleven of them not pushes (the two entries `in no function` of `internal/setup/setup.go` are two references), and finding 15 names each of the eleven. The record is posted as it was written.

## After round 1 (11:01Z to 11:03Z): notes 1 to 5 of the round

The round (comment 6166857401, `nothing material in scope`, nine notes) asked for words, and the script's logic is unchanged. The header of the script says what check (4) admits (the verb, not the subcommand, in `CloneLocal` and `AddRemote`) and what check (7) cannot see (a new caller of one of the three functions, a reference that a line break or a comment line hides, an import guard that passes when it finds no import at all); the heading of (7) and its refusal text say "name a remote verb" and "explicit name", where they said "reach a remote" and "a name other than git" (check (4) says that `AddRemote` and `CloneLocal` reach no remote). The header keeps its 197 lines, so the lines of the record's citations of the script do not move. The script now has the sha256 `b6f57973aa140cf1a37685f9a1098ae5f5f29a9fe34653cc34c594edd7be8eb6`; the one that the review read, `7a24759e8a62c1647666759746f62d5f872ac4c877cb3d2db8c21d17812d65cd`, is the file of `62fd959`.

The six mutants and the clean release again with the new script (11:02:44Z to 11:03:05Z): the clean release exits 0, and each mutant exits 1 on the same line as before; `alias-import` now says `FAIL: internal/git imported under an explicit name: internal/run/aliasmut.go:3:import gitx "github.com/pharzam/layup/internal/git"`. The output at the release (11:03:11Z to 11:03:14Z, exit 0) differs from the earlier one in the time of check (2) and the heading of (7) only; `release-check.txt` is that output.

## Criterion 3: the whole gate (10:14:02Z to 10:28:15Z, head `a0d759a`)

`go build ./...`, `go vet ./...`, `gofmt -l internal cmd` (empty), `go test ./...`, `go test -tags=integration ./...` (11 min 21 s on this host), `go test -tags=e2e ./...` (1 min 51 s), `adr-lint`, `prd-lint`, `link-lint`, `run-discipline-tests`, `nested-checkout-check`, `setup-check` and `git diff --check origin/main HEAD`: each exit 0. The branch changes no Go file, only `runs/` and `docs/`, so the commits that follow (the two pitfalls of `guardrails.md` §2 and this section) need the checks of the documents only, which the hook of each commit ran.
