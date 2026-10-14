# T-92lm — the evidence run

The prompt of the prose session is text for a model, and no unit test can read it. This record is the evidence that a real session, given
this prompt and nothing else of the rules, wrote files that `layup run` accepts. It is evidence and not a test of the gate (conditions 4 and 5
of the plan review of #326, and round 1): it shows that this prompt produced acceptable files in the attempts below, and gives no rate of
success. The demo of `M2d` (`T-sdtb`, #284) is the second real sample.

## The protocol (written before any attempt)

1. **The snapshot** is the work area of `pharzam/layup-uat-m2d` that the production loop (`setupLoop.do`) leaves at the first prose stop,
   built by `buildWorkArea` from the target's records and never changed: the brief (sha256 `ecfd6cdb9b58c6594461bc04b626214a748095aafd340ebfa70af2aa753fbf80`,
   the `psb.sha256` of the run), `inputs/answers.tsv` (15 rows, `e3ff47bf33d21b8d66755aca0200ec2d0f2fd20a90844b3fec04cdab40b0c713`), no file in `inputs/files/`, and
   `target/` with `main` at the root commit `2fc9ad7502e8e89ec4f64d718bf1d7034eef54ef` (tree `8ffb250afd584da8b418bc220fb6d72e802924ce`, the `pin.tree`) and `layup-setup` at
   the commit of S06, `754055b591261019e69837dd46ab47be17645085` (tree `bc7c558d4e685d51ebbc0764b0392ddce760c708`), which is the `base` of the sessions of runs 3 and 4
   in `sessions.tsv`. Its stop table (`stop.tsv`, 25 rows, `2bbf5c6853c7863a312c9500962e20dffa82eb4195ab66293636a0fa09f416cc`) is the table that the loop gives
   its prose session; the harness below compares it with the table of the production loop on the same inputs and records that they are equal.
2. **An attempt** is a fresh session on a fresh copy of the snapshot, laid out as `internal/session.Make` lays it out (`input/` with the brief and `stop.tsv`;
   `repo/`, a clone of `target/` on `layup-setup` with no remote, on a task branch at its head; `home/` with `.gitconfig`; `tmp/`; `result/`), with the environment of
   `session.Environ` (`PATH`, `LANG`, `HOME`, `TMPDIR`, `GIT_CONFIG_NOSYSTEM`, the credential variable and the fixed variable of the row) and the claude row of the host's
   register: `claude -p --model claude-opus-5-5 --max-budget-usd 5.00 --output-format stream-json --verbose --permission-mode acceptEdits --allowedTools Bash,Edit,Write,Read`,
   the prompt on standard input, `timeout 1800` (its wall). The prompt is `proseText` of the head under test for `pharzam/layup-uat-m2d` and the paths of that session. No
   finding of a failed check is ever given to a session.
3. **The output** is read as `layup run` reads it, before anything else: `session.ReadOutputFiles` on `result/` with the paths of the stop table (the exact set, no link, no
   other file, 1 MiB, UTF-8), then the Keep of the prose session (`proseKeep`). A refusal ends the attempt as failed. The files are used byte for byte, with no hand fix.
4. **The engine's run** is the step `verify` of `layup run`, which holds the loop of the step `setup`. The harness `zz_evidence_test.go` (a throwaway test of a copy of the
   head under test, kept on the host with its sha256, never in the repository) calls the functions that production calls: `buildWorkArea` from the brief, the rows of
   `setupAnswers` (read from the records of the target, `layup-records` at `f029e1281551ee10c65f4d082f19bc090b58b1bb`) and the session's files; `setup.StepsPhase2` through
   `setupRunner` with the phase-2 values `Name pharzam/layup-uat-m2d`, `Visibility public`, `PinSource https://github.com/pharzam/armature`, `PinCommit a95965534b14b0bf14ad74da0c9a45b5f4aedf88`,
   `PinTree 8ffb250afd584da8b418bc220fb6d72e802924ce`, `PinTime 2026-10-14T13:32:22Z`, `BriefSum` the brief's sha256, `AppID 5118100`, `Host /home/layup/layup-host`,
   `DefaultBranch main` and the questions asked; `setupLoop.do` (the rows of the markers as `layup run` adds them); `verifyJob.do` with `verifyRun` (`layup setup verify` in
   process), `setupOutputs` and `setupTree`, and a commit that lists the files it would write instead of pushing them. An attempt **passes** when the step verify is `done`.
5. **Recorded for each attempt:** the head and the binary's sha256 (the engine is the code of the copy of the head), the exact prompt (a file kept on the host, with its sha256),
   the harness's version, the model, the cap and the wall, the manifest of the output with a sha256 for each file and an archive of it kept on the host with its sha256
   (`/home/layup/layup-host/evidence/T-92lm/`), the commands and each exit code, the step rows and the stop tables, the hits of check `adapted` that the output holds (the
   engine reports the first only), the cost, the turns and the time.
6. **Order.** The attempt runs after the last round that reads the head (condition 5): the tests, the code, the documents, the full ladder and the freeze come first. A prompt
   revision that an attempt demands is made tests first, re-frozen and read by a further round before another attempt, on a fresh session and a fresh copy; a failed attempt and
   a revision are recorded here, with whether the final prompt passed on its first attempt. If the attempt on the head that the last round read fails, the task ends `not mergeable,
   findings recorded`, or the Operator moves the cap.
7. **The tools** are kept in `/home/layup/layup-host/evidence/T-92lm/tools/`: `ev-session.sh` (the session directory and the harness command) `1a5841f0…f05ac`, `ev-prompt.py` (the prompt, by a throwaway test of the copy)
   `9fc51cad…cf8f`, `zz_evidence_test.go` (the harness of item 4, which also runs the production loop to the first prose stop and compares its table with `stop.tsv`) `da53a435…ee80`, and `ev-count.sh`
   (the hits of check `adapted` of the output, by LAYUP's own `setup-check.sh`) `765bfc91…5f8f`.

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
