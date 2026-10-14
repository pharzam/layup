# T-92lm — the evidence run

The prompt of the prose session is text for a model, and no unit test can read it. This record is the evidence that a real session, given the
prompt and nothing else of the rules, wrote files that the engine accepts. It is evidence and not a test of the gate (conditions 4 and 5 of
the plan review of #326): it shows that this prompt produced acceptable files in the attempts below, and gives no rate of success. The demo of
`M2d` (`T-sdtb`, #284) is the second real sample.

## The protocol (written before any attempt)

1. **The snapshot** is the work area of `pharzam/layup-uat-m2d` as the first prose stop finds it, made once and kept unchanged: `inputs/answers.tsv`
   (15 rows), `inputs/briefs/problem-statement.md`, no file in `inputs/files/`, `out/record.tsv` with the rows of S01 to S06, and `target/` with
   `layup-setup` at the commit of S06 (`754055b`, the `base` of the sessions of runs 3 and 4 in `sessions.tsv`). It was made from the work area of
   run 4 by removing `inputs/files/` and `out/`, putting back the rows of S01 to S04 and `layup-setup` at S04 (`60072e8`), and running the engine
   once (`layup setup .`): S05 and S06 ran, and it stopped with the 25 rows of the prose step (exit 3), the table that every session gets as `stop.tsv`.
   Identity: `layup-setup` `754055b591261019e69837dd46ab47be17645085`, its tree `bc7c558d4e685d51ebbc0764b0392ddce760c708`; `main`, the root commit,
   `2fc9ad7502e8e89ec4f64d718bf1d7034eef54ef`, its tree `8ffb250afd584da8b418bc220fb6d72e802924ce` (the `pin.tree` of the run). sha256 of `inputs/answers.tsv`
   `e3ff47bf33d21b8d66755aca0200ec2d0f2fd20a90844b3fec04cdab40b0c713`, of the brief `ecfd6cdb9b58c6594461bc04b626214a748095aafd340ebfa70af2aa753fbf80`
   (the `brief.sha256` of the run), of `out/record.tsv` `9e1703f651c5b90c6d005640d9268adda0a533ba7b38673496ba886c05590cb5`, of `stop.tsv`
   `2bbf5c6853c7863a312c9500962e20dffa82eb4195ab66293636a0fa09f416cc`; the sha256 of the sorted list of the sha256 of every file of the snapshot but
   the `.git` of `target/`, `0f4f9e20841e3d46470e82365bf1fc2cd8cc97ace843cd8e7ddd927c0464f59e`. The records of the target, for the rows of the markers,
   are `layup-records` at `f029e1281551ee10c65f4d082f19bc090b58b1bb`.
2. **An attempt** is a fresh session on a fresh copy of the snapshot, laid out as `internal/session.Make` lays it out (`input/` with the brief and
   `stop.tsv`; `repo/`, a clone of `target/` on `layup-setup` with no remote, on a task branch at its head; `home/` with `.gitconfig`; `tmp/`; `result/`),
   with the environment of `session.Environ` (`PATH`, `LANG`, `HOME`, `TMPDIR`, `GIT_CONFIG_NOSYSTEM`, the credential variable and the fixed variable of the
   row) and the claude row of the host's register: `claude -p --model claude-opus-5-5 --max-budget-usd 5.00 --output-format stream-json --verbose
   --permission-mode acceptEdits --allowedTools Bash,Edit,Write,Read`, the prompt on standard input, `timeout 1800` (its wall). The prompt is
   `proseText` of the head under test for `pharzam/layup-uat-m2d` and the paths of that session. No finding of a failed check is ever given to a session.
3. **The output** is `result/files/`, copied byte for byte to `inputs/files/` of the copy with no hand fix.
4. **The engine's run** on the copy, with the binary built from the head under test, does what `layup run` does (`setupLoop.do`): `layup setup .`; on a
   stop with `M-` rows it adds the rows that the batch gave, read from the records of the target (`layup-records`, as `setupAnswers` gives them), and
   runs again; any other stop ends the run; at `O-verify` it runs `layup setup verify . > out/verify.tsv` and `layup setup .` again. An attempt
   **passes** when the last run exits 0 and every row of `verify.tsv` is `pass` or `clear`.
5. **Recorded for each attempt:** the head, the binary's sha256, the exact prompt (a file kept on the host, with its sha256), the harness version, the
   model, the cap and the wall, the manifest of the output with a sha256 for each file, and an archive of the output kept on the host with its sha256
   (`/home/layup/layup-host/evidence/T-92lm/`), the commands, each exit code, the step and stop tables, the hits of check `adapted` that the output
   holds (the engine reports the first only), the cost, the turns and the time.
6. **Order.** The attempt runs after round 1 on the head (condition 5): the tests, the code and the documents, the full ladder and the freeze come
   first. A prompt revision that an attempt demands is made tests first, re-frozen and read by round 2 (the last of the cap of 1) before another attempt, on a
   fresh session and a fresh copy; a failed attempt and a revision are recorded here, and the record says whether the final prompt passed on its first
   attempt. If the attempt on the head that the last round read fails, the task ends `not mergeable, findings recorded`, or the Operator moves the cap.

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
