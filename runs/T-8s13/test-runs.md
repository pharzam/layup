# T-8s13: test runs of the bound script (step T)

The test `brief-bound-test.sh` (nine cases) pins the terms and `L` that `brief-bound.sh` prints, on small records it builds in a temporary
plain git repository (`refs/remotes/origin/layup-records`). The formula is the one of `internal/run/betbrief.go` (`assemble`, with
`requirementsPart`, `solutionPart`, `briefOpenPart` and `briefFactsPart`) and the plan review (condition 1):
`L = 12 + R + max(M,1) + sum(1+O_d) + Q_io + E + S + A_io + C_op`.

## Red: the stub prints wrong terms

Command (from the root of the worktree at `76dfb684` plus these two files): `sh runs/T-8s13/brief-bound-test.sh`

- Start 2026-10-21T16:59:18Z, end 2026-10-21T16:59:18Z (read from `date -u` before and after the command).
- Exit 1. Nine cases, nine `FAIL`, no `ok`. Each `FAIL` is an assertion: the stub `brief-bound.sh` prints
  `R=0 M=0 S=0 A_io=0 C_op=0 Q_io=0 E=0 L0=0` and the case wants its own terms (the first: `R=2 M=0 S=2 A_io=0 C_op=1 Q_io=1 E=0 L0=19`
  and the line `G1 floor=22 limit=80 stop=no`). The reason is the intended one: the stub does not count.

## Green: the script counts

Runs of the same command while implementing: 2026-10-21T17:00:03Z and 17:00:08Z, exit 1, one `FAIL` each (the follow-up case: a quote compared as a string in awk, then a lookup in the
wrong array, both defects of the implementation); 17:00:12Z, exit 0. The fixtures were then made valid records (see the cross-check below) and the test run again:

Final run: start 2026-10-21T17:01:13Z, end 2026-10-21T17:01:13Z (`date -u`), exit 0, nine cases, nine `ok`, no `FAIL`:

```text
ok only an idea-owner question open: E=0, Q_io=1
ok no question open: E=1 for the line - none
ok an Operator and an idea-owner question open: E=1, Q_io=1
ok answered follow-up closes its question, asked one keeps it open; C_op by quote class
ok idea-owner question with an unanswered follow-up is open
ok G2 at L=80: room 20 passes
ok G2 at L=81: room 19 fails
ok G1 floor 79 (76 plus one decision of two options) passes
ok G1 floor 81 stops
```

The cases pin: the open-question case of plan condition 1 (E=1 when an Operator question is still open or none is open, E=0 when only idea-owner questions are open;
"still open" by `stillOpen` of `internal/run/shapequestions.go:25`, an answered follow-up closes its question and an asked one does not); C_op by the quote class of the
question an answer asks (a follow-up takes the class of the question it asks again); the 80/81 boundary of G2 (L=80 gives room 20, pass; L=81 gives room 19, fail); and the
G1 floor of 79 (L0=76 plus one decision of two options), with 81 stopping.

### Mutants (each a copy of the script, the test run on it)

- E always 1: 2 cases fail. C_op always 1: 7 cases fail. G2 room `> 20` instead of `>= 20`: the L=80 case fails.

### Cross-check against `internal/run/betbrief.go` (a throwaway Go test, not committed, deleted after the run)

The same fixtures (kept with `KEEP=DIR`) were read by `records.ReadDraft`, `ReadSpans`, `ReadQuestions` and `ReadAnswers` (which also proved them valid records), and
the brief's parts were built by `requirementsPart`, `solutionPart` (the two heads of `betbrief_test.go`, `filterAll` on the options table), `briefOpenPart` and `briefFactsPart`;
`strings.Count(assemble(parts, "", "", ""), "\n")` is the `L` that `overLimit` compares with `brief.lines`. The Go count and the script's `L` (with HEAD) or `L0` agree on every case:
followups 20, g1-79 76, g1-81 78, g2-80 80, g2-81 81, io-follow 17, io-open 19, none-open 20, op-open 19.

## Supplementary: the M2d records (read-only)

`git clone https://github.com/pharzam/layup-uat-m2d.git`, `origin/layup-records` at `d25f4df469d21d345c278af1eae73a6224efc92a`; `sh runs/T-8s13/brief-bound.sh CLONE` at 2026-10-21T17:00:51Z, exit 0:

```text
R=18 M=3 S=25 A_io=14 C_op=3 Q_io=0 E=1 L0=76
G1 floor=79 limit=80 stop=no
```

76 before decisions, as the plan predicted (R=18, M=3, S=25, A_io=14, C_op=3, Q_io=0, E=1); the floor is 79, so G1 passes with room 21. The Go functions above, run on the
same four files of that branch, also count 76. This is a comparison, not the test: the target of the demo is `pharzam/layup-uat-m2e`.

## Round 1 (F1, F2): three cases added, red first, then green (2026-10-23)

Round 1 of the review found that `brief-bound.sh` hid git failures behind pipelines (F1) and read a table header with tabs as no header (F2). The test got three cases:
the header `|\tMetric\t|\t...Fact\t|` on the 81-line fixture (wants `M_prd=2 L=81 room=19 G2=fail`, as the space form), an unreadable head (`no-such-verified-head`) and a missing
records ref; the last two want a non-zero exit and neither `G2=pass` nor `stop=no` in the output.

- Red: `sh runs/T-8s13/brief-bound-test.sh` on the script of 7a64e43, 2026-10-23T06:28:59Z to 06:29:00Z (from `date -u`). Nine `ok` and two `FAIL`. The tab case printed `M_prd=0 L=80 room=20 G2=pass`.
  The unreadable head printed `fatal: Not a valid object name`, then `G2=pass`, and the script exited 0. The third case (missing ref) was already green: its first read, under `set -e`, stopped the script before any line.
- Fix: every git read goes through one function that checks the status before the output is parsed (exit 2 and a message); the head is verified first; all reads run before the first line is printed;
  the header test trims spaces and tabs, as `tableCells` of `internal/spec/full.go` does.
- Green: the same command, 2026-10-23T06:29:08Z, exit 0, twelve cases, twelve `ok`.

## The binaries and the release checks of the demo (added at the close-out, 2026-10-23)

The host binary is `/home/layup/layup-host/bin/layup`. It was rebuilt before a run whenever a defect of the demo was fixed on `main` (a task of its own, defects.txt), with
`go build -o /home/layup/layup-host/bin/layup ./cmd/layup` (go1.26.10) at the head of this branch after it merged `main`, and the release check `sh runs/T-efmy/release-check.sh <full SHA>` ran on that commit
(outputs `T-8s13-release-check*.txt` in `/home/layup/.cache/layup-session/`, 141 lines each, every one ending `== PASS`). The size of the change against the release,
`git diff --stat 01d7da3082333efadb227afda14fd98715edf5c4 COMMIT -- cmd internal go.mod` (tests included), is the last column.

| Release check | Built from (full SHA) | SHA-256 of the binary | Diff against the release | In place for the runs |
| ------------- | --------------------- | --------------------- | ------------------------ | --------------------- |
| `T-8s13-release-check.txt` | 20c7cc3772a8d81f6c6c4a97161003327fa95744 | 478351653b13f38b788c523110fea623175b0ccb2c490e11442c003a0456607a (recorded at the build, `T-8s13-binary.txt`, built 2026-10-21T17:07:21Z) | 8 files, +823 -14 | r0 to r5b |
| `-2` | 48d0a8b5917f2a2d9a31a9e33ad667bcb2c6d2be | 5141367097df5e48d2dd2b295e6b977959dcb6576cab3b7fd3e5325126287dac | 10 files, +861 -17 | r5c |
| `-3` | 314e5f82dc8adda4b57683abdd96b74e35cd6e68 | 002c08039a6f3f3b29fd8d5e3975b95cfce5c60681fe2b32a7352f424a647922 | 21 files, +1367 -47 | r5d |
| `-4` | 03d729e967e2dcd35e4f103cec53f53c08ca61ee | 7a0d5cecb67bac9cab4100c2faf5ef0a5c9625a5fc8ea45f7a54661ad54996b1 | 34 files, +1698 -59 | r5e |
| `-5` | 6f249ac3caa228a3b53562c0c2b6e418d331555d | 01a9bd112ec44536e01a426c3e92ac9d341ee7dd7ea2e93ea19000db0d9b2256 | 48 files, +2256 -87 | r5f |
| `-6` | b3c87dbc9015b2fdca4288099fbb032d0afa29c8 | 2ed30140b984857a849cbed87d469d0c0b85dde9340675e1ac0ce234f31d16dd | 60 files, +3363 -122 | r5g, e2-r0 to e2-r2 |
| `-7` | ca02feb89bb4a89de5991dbec5293595a1c67c10 | e2ec40e25226005552c3ebcd8c4e814c754e832b3dc891448a5bed4fd3e7f731 | 61 files, +3414 -122 | e2-r3 to e2-r5 |
| `-8` | e4a79cda1b4dc42c22e3e6edcc9ee988cfa606be | 6ab75316a571d9e9ba5a2c679de3fe88ae80fd3e2a18d67cb88d6d470ffe6f04 | 65 files, +4188 -135 | e2-r6 |
| `-9` | e3d0318213a6e33db1a771dc28015a00753f3e64 | a398b85004867e8393c5bc95f8f765b0cae6e0a6dfd8532446f465aecb11d922 (the installed file, read 2026-10-23) | 70 files, +5086 -150 | e2-r7 |

How the hashes are tied to the commits, stated so that nothing is taken for more than it is. Only the first row was written down at the build. For the others the coordinator copied each fresh build to
`layup.before-T-<next task>` within one to two seconds of the release check of that build (file times in `/home/layup/layup-host/bin/` against the release-check files: `-2` 09:13:42 and
`layup.before-T-qjjd` 09:13:41; `-3` 12:25:51 and `before-T-0sj8` 12:25:50; `-4` 14:04:11 and `before-T-yngf` 14:04:10; `-5` 15:12:54 and `before-T-t94t` 15:12:53; `-6` 17:10:50 and
`before-T-paan` 17:10:48; `-7` 18:53:04 and `before-T-dtw6` 18:53:03; `-8` 21:00:04 and `before-T-0xr7` 21:00:03; `-9` 08:50:10 and the installed `layup` 08:50:09, host local time, UTC+3). The first row
checks this reading: `layup.before-T-c3zv` is 478351653b13…, the recorded build of 20c7cc3. The column "In place for the runs" is the build whose release check precedes the run's start in `runs.txt`
(for example r5c starts 2026-10-22T06:13:57Z, after the `-2` check at 06:13:42Z). A rebuild of 20c7cc3 in a second worktree gives 5ee307aa5e95…, not 478351653b13…, so a build is
tied to its path and the hashes cannot be re-derived from a commit; they are the files' own. The old binary before the demo is `layup.before-T-8s13`,
1b4b96d8de1ccf5227b5d8d13a5e199b5980534a78b2abfbb76e4d11201f2cf1 (the build of T-sdtb at 8b76fdc, M2d).

The test of the bound script above ran nine cases, nine `ok`, at the step-T head; round 1 added three (the section above), twelve in all. The G1 figures of both targets are in `rules.txt`; `brief-bound.sh` and its test count as executable files, and everything else in this directory
(runs.txt, records.txt, rules.txt, defects.txt, this file's tables) is evidence.
