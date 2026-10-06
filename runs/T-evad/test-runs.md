# T-evad: the test runs

The test runs of `T-evad` (#97), the first pilot. Host: macOS, `go1.27.1`,
`git` 2.54.0, 2026-10-05 and 2026-10-06, UTC. The pilot run itself, with its
evidence, is in [`README.md`](README.md). A file named `evidence/` or `tools/` is on
the pilot's host.

## The check of D8 (b): `rules-diff.sh`, red first

[`rules-diff.sh`](rules-diff.sh) and its fixture test
[`rules-diff-test.sh`](rules-diff-test.sh) (step 1 of the plan, condition 4 of its
plan review). Each case builds a small target whose root commit is the baseline
and whose next commits are setup steps, and names one message that must be in the
output and, for a case that must fail, one that must not.

- **Red, 07:31Z.** The cases ran against a skeleton that passed each target: each
  case that must fail passed, so it failed the test.
- **Each new case red first.** The cases of a section replaced as a unit, a heading
  that is only renamed, the placeholder row of an index, a marker with a column, and
  a row hash at the commit of its step each failed first (07:36Z to 08:03Z).
- **Green, 08:07Z:** 16 of 16 cases (`6206335`).
- **Mutations, 08:11Z:** 9 mutations of the script, each caught by at least one
  case.
- **F-18, 18:36Z:** case 17 (`written-markers`) failed on the report "2 marker lines
  filled" for a written file with one gap row; green at 18:37Z, 17 of 17
  (`f9fcc9d`); the 9 mutations are still caught.

## The check of D8 (a): the planted rows

The red fixtures of the value audit (step 1 of the plan; the answer to the plan
review, item 3) are two planted rows in the table that the auditor read, which the
auditor did not know about: V-011, a value with no row and no source, and V-049, a
citation that does not support its value. Pass 1 of the audit (18:07Z to 18:15Z)
gave V-011 `no source` and V-049 `not supported`. `tools/value-inventory.py check`
finds each of the 66 values at its stated place in the target.

The audit of the second run planted two rows the same way in its pass 1: W-032, a
value with no source, and W-034, a citation that does not give its value. Pass 1
(2026-10-06, 06:24Z to 06:31Z) gave them `no source` and `not supported`.
`tools/value-inventory-2.py check` finds each of the 67 values of pass 2 at its place,
and no line that the diff of the two setup heads adds without a row
(`evidence/119-audit2-pass2-check.txt`).

## The two fixes of F-8, red first

- **S11, two markers on one line** (`fa53cfb`): `TestS11TwoMarkersOnOneLine` failed
  first (the three rows of one line had one name, and the record refused the
  repeated key), then passed.
- **Check `sources`, the column of a gap row** (`39f1eb8`):
  `TestTheSourcesOfARecord` failed first ("line 2 of docs/a.md does not hold" the
  marker "at column 9"), then passed.
- After both: the unit and integration levels and `go test -tags=e2e ./...` passed.
  The binary was built again from `6206335`, and the real run used it.

## The tools of the pilot

| Tool | Test | Result |
| ---- | ---- | ------ |
| `format-forms.sh` (F-1, F-9, F-11) | 13 cases of the lint command against the gate script of R-TEST-05 | Red with the first one-line form (4 mismatches); green with the sentinel form; green with GNU `xargs` and `dash` in a container; case c13 kept as a known limitation |
| [`t3.sh`](t3.sh) (F-13) | `tools/t3-test.sh`, a stand-in `gh` and local remotes | 9 of 9 |
| `t4-1.sh`, `t4-2.sh`, `t4-3.sh` (D7) | their tests, a stand-in `gh` and local remotes | 6 of 6, 5 of 5, 7 of 7 |
| `checks-claims.sh` (F-10) | 40 claims of the target's table "Which checks run" | All hold at the setup head, on GitHub's `main` after T3, and after the merge of D7 |
| `t6.sh` (T6) | `tools/t6-test.sh`, a stand-in `gh` and local remotes | 7 of 7 (`evidence/108-t6-test.txt`); its check mode passed on the second run (`evidence/116-t6-check.txt`) |
| `t7-1.sh` (the target task), SHA-256 `60f8d2b9fa6ce5ce88252f9ddb10d3f69f17cb2f29858cf5a76252dfa26c9689` | `tools/t7-1-test.sh`, a stand-in `gh` | 6 of 6 |
| `port.sh` (the copy of `T-vu2j`), `dd1a7f53d88b00114df0191cbc68dbec0c6c6d3f4b0f89eb0521826071825776` | `tools/port-test.sh`, 20 cases, with `tree-equal.sh` | 20 of 20 |
| `c3.sh` (the rules of criterion 3 of `T-vu2j`), `8bf349e976e37550005c3eba5655acc301c12dc0d607de6ad5d781ecdbae06a8` | `tools/c3-test.sh`, 35 cases | 35 of 35; each rule of the review rounds red first; `c3.sh gate` passed on the gate run of the close-out head, and `c3.sh ci` on `pharzam/chat-orchestrator#7` |
| `t7-2.sh` (the records and the pull request of `T-vu2j`), `54617cc7c692beae59ba3d702aadb0f4a7e54180fc313b0e9df5732e1ce06871` | `tools/t7-2-test.sh`, 21 cases, a stand-in `gh` and local remotes | 21 of 21; a dry run with each write refused passed its checks; then the Operator's run |
| `t7-3.sh` (the merge of `T-vu2j` and the branches of O-152), `3e9daa194d1289b2d6138a12a39ddf9905201c09f103cb800737568c004c9530` | `tools/t7-3-test.sh`, 27 cases, a stand-in `gh` and local remotes | 27 of 27; a dry run with the merge refused passed its checks; then the Operator's run |
| `hb.sh` (a progress line every 5 seconds; a stop reaches the command), `121546107ac0d982fae121d06f984a81d6308ec924dede132cc004db4ad47998` | `tools/hb-test.sh`, 5 cases | 5 of 5 |

The tests of the scripts of `T-vu2j` pass also under bash 5.3. Under `dash`, the two cases that send INT to a script in the
background cannot run, as `dash` cannot turn on job control without a terminal; the same cases pass with TERM.
| [`tree-equal.sh`](tree-equal.sh) (condition 2 of O-146) | [`tree-equal-test.sh`](tree-equal-test.sh), 5 cases | Red against a skeleton that passed each case, then 5 of 5 (`9226a7b`) |
| [`tree-equal.sh`](tree-equal.sh) with `--base` (O-148 to O-150, the plan review of the target task) | [`tree-equal-test.sh`](tree-equal-test.sh), 40 cases | Red against a skeleton that took the new options and did nothing with them: 27 of 40 failed. Each case that loses or changes a path of `T-a0rt`, or its line of the log, exited 0 where 1 is wanted (`evidence/124-tree-equal-base-red.txt`). Then 40 of 40 (`evidence/124-tree-equal-base-green.txt`), and 5 runs each under `dash` and bash 5.3 (`evidence/126-tree-equal-base-shells.txt`). 7 of 7 mutations caught, each by the case of its rule (`evidence/125-tree-equal-base-mutations.txt`) |
| [`tree-equal.sh`](tree-equal.sh): the form and the place of the task's line in the log, `--issue` (the second plan review of the target task) | [`tree-equal-test.sh`](tree-equal-test.sh), 47 cases | Red against a skeleton that took `--issue` and did nothing with it: the five cases of a line with no date, no link, another issue or a place below `T-a0rt` exited 0 where 1 is wanted (`evidence/130-tree-equal-line-red.txt`). Then 47 of 47, also under `dash` and bash 5.3, and 5 of 5 mutations caught; the fixture log has the example line of an HTML comment, as the target's log has (`evidence/131-tree-equal-line-green.txt`) |
| [`tree-equal.sh`](tree-equal.sh): the repository of the link, `--issue OWNER/NAME#N` (condition 1 of the third plan review of the target task) | [`tree-equal-test.sh`](tree-equal-test.sh), 52 cases | Red against a skeleton that took `OWNER/NAME#N` and kept only `N`, the rule of `e82d5fa`: the case of a line that links another repository exited 0 where 1 is wanted (`evidence/136-tree-equal-repo-red.txt`). Then 52 of 52, also under `dash` and bash 5.3. 3 of 4 mutations caught. The fourth removed the check that `--issue` holds a `#`; that check was redundant, as the check of the number (digits after the last `#`) and the check of the repository (`OWNER/NAME`) refuse each input without `#`, which the cases `issue-form` (`7a`) and `issue-number` (`7`) show, so the line was removed (`evidence/137-tree-equal-repo-green.txt`) |
| [`tree-equal.sh`](tree-equal.sh): a named file of the task is a file of mode 100644 (finding 4 of round 2 of the target task) | [`tree-equal-test.sh`](tree-equal-test.sh), 54 cases | Red: a task file that is a symlink, and an evidence file of mode 100755, each exited 0 where 1 is wanted. Then 54 of 54, also under `dash` and bash 5.3, and the mutation that drops the rule caught (`evidence/149-tree-equal-mode-red.txt`) |
| [`tree-equal.sh`](tree-equal.sh): the date of the task's line is a real calendar date (finding 2 of round 3 of the target task) | [`tree-equal-test.sh`](tree-equal-test.sh), 62 cases | Red: six dates that do not exist (`2026-99-99`, `2026-13-01`, `2026-02-30`, `2026-02-29`, `2026-00-10`, `2026-04-31`) each exited 0 where 1 is wanted. Then 62 of 62, with a real leap day, also under `dash` and bash 5.3, and two mutations caught: no check of the day, and no rule of the leap year (`evidence/159-tree-equal-date.txt`) |
| [`tree-equal.sh`](tree-equal.sh): the summary of the task's line starts with a character that is not a space (finding 1 of round 4 of the target task, O-157) | [`tree-equal-test.sh`](tree-equal-test.sh), 64 cases | Red: a summary of three spaces, and a summary of one tab, each exited 0 where 1 is wanted. Then 64 of 64, also under `dash` and bash 5.3, and the mutation that takes `.+` for the summary again caught; `port-test.sh` 20 of 20 with it (`evidence/172-tree-equal-summary.txt`) |

## The fixes of the markers (fixes 1 to 3), red first

The idea owner's fix of T5 (comment 6002034171), with the design points of O-146
(comment 6002406785). On 2026-10-05:

- **Red, 20:24Z and 20:25Z.** The new cases of nine tests of `internal/verify`
  (`TestAMarkerOverMoreLines`, `TestTheMarkersOfALine`, `TestTheMarkersOfAText`,
  `TestTheOpenGaps`, `TestTheSourcesOfARecord`, `TestTheHitsOfAdapted`,
  `TestTheTextOfAdapted`, `TestLostMarkers` and `TestUnpairedQuotes`) and of five tests
  of `internal/setup` (`TestS10`, `TestS11`, `TestS11AMarkerOverMoreLines`,
  `TestS11ChecksBeforeItWrites` and `TestS14RefusesAnInputThatLosesAMarker`) failed
  first (`evidence/fix123-verify-red.txt`, `evidence/fix123-setup-red.txt`). Then the
  code: the unit level passed.
- **The fixtures of `setup-check.sh`, 20:29Z.** `markers/bad-unpaired` and
  `adapted/bad-unit` failed `TestTheFixturesOfSetupCheck` and
  `TestTheShAndTheGoFormOfMarkersAgree`, until the script read a marker over more
  lines as the Go code does. The integration level passed at 20:36Z.
- **The end-to-end level, 20:37Z.** The stand-in of the whole setup, with a marker
  over two lines, failed five cases of `cmd/layup`. It passed at 20:45Z.
- **Green, 20:45Z:** `gofmt` clean; `go vet` with each build tag; the unit,
  integration and end-to-end levels; 81 discipline tests; 46 setup-check fixtures;
  `adr-lint`, `prd-lint`, `link-lint`, the nested-checkout check; 18 cases of
  `rules-diff-test.sh`; `git diff --check` clean (`c68b390`, `d11ef9c`;
  `evidence/fix123-ladder.txt`).
- **The count of places, 20:51Z.** A new case of `TestLostMarkers` (two places of a
  marker on a flagged line, one kept) failed first, then passed (`cc31402`).
- **Every place, 20:55Z.** The case of a marker on a line that check `adapted` does
  not flag changed from "no finding" to a finding, and the levels passed again
  (`006fbee`). With this rule, S14 finds the 27 places of the first pilot.
- **`rules-diff.sh`:** a marker over more lines, read by its key (18 cases, 11
  mutations caught, `d11ef9c`); a task index whose filled marker over more lines
  becomes one line (20 cases, 12 mutations caught, `beecf92`).

## The rehearsals and the second run

| Run | What | Result |
| --- | --- | --- |
| Rehearsal A | The first pilot's inputs, with the fixed binary | S14 refuses the inputs: 22 markers in 27 places lost (`evidence/100-rehearsal-fix-a.txt`) |
| Rehearsal B | The corrected prose inputs (the places kept) | S14 passes; S10 asks 115 questions (101 answered before, 14 new); `verify`: 15 `pass` and 3 `clear`; `rules-diff.sh`: PASS; the seven places of T5 right (`evidence/101` to `106`) |
| The second run | A new work area, the same baseline commit `a959655`, the approved prose and the 14 new answers | `verify`: 15 `pass` and 3 `clear` (`evidence/112`); `rules-diff.sh`: PASS (`evidence/114`); the seven places right (`evidence/115`); setup head `e2b402b` |
