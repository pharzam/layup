# T-92lm — the evidence run

The prompt of the prose session is text for a model, and no unit test can read it. This record is the evidence that a real session, given this prompt and nothing else of the rules, wrote
files that `layup run` accepts. It is evidence and not a test of the gate (conditions 4 and 5 of the plan review of #326): it shows that this prompt produced acceptable files in the
attempts below, and gives no rate of success. The demo of `M2d` (`T-sdtb`, #284) is the second real sample.

## The protocol (written before any attempt)

1. **The snapshot** is the work area of `pharzam/layup-uat-m2d` that the production loop (`setupLoop.do`) leaves at the first prose stop, built by `buildWorkArea` from the target's records by the
   harness of item 4 (`TestZZStopTable`) and kept unchanged on the host (`/home/layup/layup-host/evidence/T-92lm/snapshot.tar`, sha256 `7b18082e8a5ef415420b15b3834a9f167a316823ab8dc718b2f5e547b868aa96`).
   Its **manifest hash** is the sha256 of the sorted list of the sha256 of every file but the `.git` of `target/`: `c4652dea77c719b674b5a1b927f2961e262c3da89155ba8b72d163b13e385af2` (built twice, the same).
   It holds the brief (sha256 `ecfd6cdb9b58c6594461bc04b626214a748095aafd340ebfa70af2aa753fbf80`, the `psb.sha256` of the run), `inputs/answers.tsv` (15 rows,
   `e3ff47bf33d21b8d66755aca0200ec2d0f2fd20a90844b3fec04cdab40b0c713`), no file in `inputs/files/`, `out/record.tsv` (22 lines), and `target/` with `main` at the root commit
   `2fc9ad7502e8e89ec4f64d718bf1d7034eef54ef` (tree `8ffb250afd584da8b418bc220fb6d72e802924ce`, the `pin.tree`) and `layup-setup` at S06, `754055b591261019e69837dd46ab47be17645085`
   (tree `bc7c558d4e685d51ebbc0764b0392ddce760c708`), the `base` of the sessions of runs 3 and 4 in `sessions.tsv`. Its stop table `stop.tsv` (25 rows, 4,485 bytes,
   `2bbf5c6853c7863a312c9500962e20dffa82eb4195ab66293636a0fa09f416cc`) is equal, byte for byte, to the table of the production loop (`TestZZStopTable`, 20:25Z and 20:40Z).
2. **An attempt** is a fresh session on a fresh copy of the snapshot, laid out as `internal/session.Make` lays it out (`input/` with the brief and `stop.tsv`; `repo/`, a clone of `target/` on
   `layup-setup` with no remote, on a task branch at its head; `home/` with `.gitconfig`; `tmp/`; `result/`), with the environment of `session.Environ` and the claude row of the host's register:
   `claude -p --model claude-opus-5-5 --max-budget-usd 5.00 --output-format stream-json --verbose --permission-mode acceptEdits --allowedTools Bash,Edit,Write,Read`, the prompt on standard input,
   `timeout 1800` (its wall). The prompt is `proseText` of the head under test for `pharzam/layup-uat-m2d` and the paths of that session. The session script recomputes the manifest hash of
   the snapshot before the session and the record says that it equals the one above. No finding of a failed check is ever given to a session.
3. **The end and the output** are accepted as `endPhase` accepts them, in its order: the output is read by `session.ReadOutputFiles` on `result/` with the paths of the stop table (the exact set, no
   link, no other file, 1 MiB, UTF-8); the class of the end is `session.Class` of the exit status of the process (the exit 124 of `timeout` is a stop by the wall) and of the reader's error, and
   must be `done`; then the Keep of the prose session (`proseKeep`). Any refusal, or a class other than `done`, ends the attempt as failed. The files are used byte for byte, with no hand fix.
4. **The engine's run** is the step `verify` of `layup run`, which holds the loop of the step `setup`, by the harness `zz_evidence_test.go` (a throwaway test of a copy of the head under test, never in the
   repository) through the functions that production calls: `buildWorkArea` from the snapshot's brief and root commit, the rows of `setupAnswers` (from the records of the target, `layup-records` at
   `f029e1281551ee10c65f4d082f19bc090b58b1bb`) and the session's files; `setup.StepsPhase2` through `setupRunner` with the phase-2 values `Name pharzam/layup-uat-m2d`, `Visibility public`,
   `PinSource https://github.com/pharzam/armature`, `PinCommit a95965534b14b0bf14ad74da0c9a45b5f4aedf88`, `PinTree 8ffb250afd584da8b418bc220fb6d72e802924ce`, `PinTime 2026-10-14T13:32:22Z`, `BriefSum`
   the brief's sha256, `AppID 5118100`, `Host /home/layup/layup-host`, `DefaultBranch main` and the questions asked; `setupLoop.do` (the rows of the markers as `layup run` adds them); `verifyJob.do` with
   `verifyRun` (`layup setup verify` in process), `setupOutputs` and `setupTree`, and a commit that lists the files it would write instead of pushing them. An attempt **passes** when the end and the output
   are accepted (item 3) and the step verify is `done`.
5. **Recorded for each attempt:** the head and the sha256 of the binary built from it, the exact prompt (kept on the host, with its sha256), the harness version, the model, the cap and the wall; the manifest
   hash of the snapshot before the session; the exit status, the class of the end and the final result line of the stream (`subtype`, `is_error`, `terminal_reason`); the manifest of the output with a sha256 for
   each file, and an archive of it kept on the host with its sha256 (`/home/layup/layup-host/evidence/T-92lm/`); the results of the reader, the Keep and the engine's run (the step rows and the stop tables); the
   hits of check `adapted` that the output holds (the engine reports the first only); the cost, the turns and the time.
6. **Order.** The attempt runs after the last round that reads the head (condition 5): the tests, the code, the documents, the full ladder and the freeze come first. A prompt revision that an attempt
   demands is made tests first, re-frozen and read by a further round before another attempt, on a fresh session and a fresh copy; a failed attempt and a revision are recorded here, with whether the final
   prompt passed on its first attempt. If the attempt on the head that the last round read fails, the task ends `not mergeable, findings recorded`, or the Operator moves the cap.
7. **The tools** are kept in `/home/layup/layup-host/evidence/T-92lm/tools/`: `ev-session.sh` (the session directory, the manifest and the harness command) `e3afdd16…c24f`, `ev-prompt.py` (the prompt, by a
   throwaway test of the copy) `9fc51cad…cf8f`, `zz_evidence_test.go` (the harness of items 1, 3 and 4) `eb1b88cf…0a80`, and `ev-count.sh` (the hits of check `adapted`, by LAYUP's own `setup-check.sh`) `765bfc91…5f8f`.

## The control: the old prompt on the same snapshot (not an attempt)

It runs merged code, so it needs no review; it shows that this layout reproduces the failure of runs 3 and 4. The prompt is the one of `b7f3f7d`, 1,395 bytes with the
paths of a real session as in `sessions.tsv`. Session `ctl1`, `claude` 2.1.300, `claude-opus-5-5`, cap USD 5.00, wall 1,800 s; 18:47:27Z to 19:03:13Z (15 min 46 s), 56 turns,
USD 3.76 on the list basis (the billing is a subscription), 25 files written, exit 0. The prompt file has sha256 `203e17496682e374d219c93adc2b38554b1feb63789b1a8224b31b177e7ad8ea`;
the manifest of the output `50704aee905ed72aa777e7d6ef97557b0a7bf3bf6f26a0d377f27dd46a68dc88`; the archive `/home/layup/layup-host/evidence/T-92lm/ctl1.tar.gz`
`31e00063e9c5c9db751f1f4f9073a004d2a8eeb8ad1f56c123f1f9e296d37035`.

The engine (`layup` sha256 `7f6379d6d7cf1d22daa3a35f19e07ef2f0a72d8d2e25332cb9e9bd6f0db49a1d`, built at `236c9af`), on the snapshot with the control's files in `inputs/files/`: S01 to S07
done, **S08 `fail` with `glossary: fail: fact: F-0007#3 is not a fact of the F-0007 record`**, the same text as runs 3 and 4; S09 to S15 not run. Check `adapted` (LAYUP's own `setup-check.sh`
with the lists of the engine, which reports every hit) holds **154 hits in 24 of the 25 files**: kit 28, optional 26, adapt 16, fill-before-a-marker 15, replace-before-a-marker 11, delete this 11,
the template 10, your stack 9, adopter 8, your project 6, fill in 6, your forge 3, you use 2, skip this section 2, Armature 1.

## Attempt 1 — the head `4fbfc06`, before round 3 (02:07:16Z to 02:16:47Z)

Run while the Operator's answer at the cap was awaited, on the head that rounds 1 and 2 had read (order of condition 5 kept: the code was reviewed before the paid session), as the protocol says; the final head
differs from it by four sentences of the prompt (the stem context, the stopping rule, spaces and tabs, the passed words) and is attempt 2's. **Passed, on its first attempt.**

| What | Value |
| ---- | ----- |
| Binary and prompt | `layup` built from the copy of `4fbfc06`, sha256 `cec1058a82df45d7b864d2b6499ce35754752a98715d27f2b9c0441074ff0ada`; the prompt 7,226 bytes, sha256 `3a8011b0b1bc2e7e57c726cd812c6cb6920346de572a568d233df6e93dda9602` |
| Session | `claude` 2.1.300, `claude-opus-5-5`, cap USD 5.00, wall 1,800 s, snapshot manifest `c4652dea…85af2` before it (equal); exit 0, result `success`, `terminal_reason completed`, `is_error` false; 48 turns, 568.9 s in the stream, **USD 3.49** (list basis; subscription); `repo/` left unchanged |
| End and output | class `done`; `session.ReadOutputFiles` and `proseKeep` accept 25 files (manifest sha256 `2336155a…102d`; archive `a1.tar.gz`, sha256 `8f3d321c…441e`, in `/home/layup/layup-host/evidence/T-92lm/`) |
| Engine | step verify `done` ("the setup verified and its records committed"), the 18 rows of `verify.tsv` `pass` or `clear` (its sha256 `368a24d4…47c8`, that of the mechanically fixed control files too); the commit would hold `README.md`, `rule-paths.tsv`, `setup/record.tsv`, `setup/tree.tsv` and `setup/verify.tsv` |
| Hits of check `adapted` | 0 (the control had 154); 210 markers kept in 19 files; glossary citations all of `F-0001`; `README.md` holds `pharzam/layup-uat-m2d` and `layup-records` |
