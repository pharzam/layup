# T-evad — the first pilot

Issue: [#97](https://github.com/pharzam/layup/issues/97), row 20 of the
[implementation plan](../plan/README.md); child of the core engine (#29). Serves
`REQ-001`, `REQ-002`, `REQ-004`, `REQ-007`, `REQ-016`, `REQ-018`, `NFR-001`,
`NFR-002`, `NFR-003`, `NFR-004` and `NFR-006`. Base `c961652` (the merge of row 19,
#119). Authors: Claude Opus 5.5 (the plan, the answers, the work from T2 revision 2
on, the fix, the second run and the target task `T-vu2j`) and Claude Sonnet 5.5 (the
work from T1 to T2 revision 2, with four helper sessions of Sonnet 5.5 for S14), on
Claude Code. Evidence: [`runs/T-evad/`](../../runs/T-evad/README.md).

## Plan and plan review

The question of the inputs (comment 5970993385) got the Operator's answers O-139 and
O-140 (5979471031): the brief `PSB-CHAT-001` (Go), the target
`pharzam/chat-orchestrator`, the answers on its issue `Setup answers`, the prose
drafted by the author, the answers proposed by the author and approved by the
Operator in one comment, the baseline's newest commit, no App on the target, the
Operator runs `commands.sh`, and #33 and #77 close with this task. O-141 (the
session) closes #29 too.

The plan (5979571232) and its review by GPT-6 Sol (`xhigh`, on the Devin CLI, a
fresh read-only session; 5979726108) gave `approve-with-conditions`: Budget maximum
2,400 lines over 32 files against the base, close-out inside; Cycle cap 2, as
`rules-diff.sh` is a check script. The review's pilot rule replaced D9: Pass (1) to
(8), Investigate, Fail, with the idea owner's rejection a Fail. Its eight conditions
are applied (5979731752): the answers cite the comment that holds them, and the
prose cites the raw brief by its path (1); the gate runs from outside on the base
and head recorded before the merge (2); D8 (a) audits each value, with planted rows
and the Operator's check (3); `rules-diff.sh` compares the union of the paths, with
red cases for a deleted rule and a rule changed despite a row (4); a rejection is a
Fail (5); the classes of D10 (6); #29 by the Operator's answer (7); the cap (8).

**The Operator's decisions and rulings during the first run:** the four corrections
of T2 (5993893238) and the two of T2 revision 2 (in the session); F-1 a known
limitation, with no catalog change in this task; the three prose corrections and
F-11 a known limitation (5995693012); the prose approved for `cec749a` (5996154495);
O-142 (a), D7 by the target's own process (5996579611); O-143 (a), one round that
applies all the lenses (5996877039); O-144 (a), round 1 of `T-a0rt` recorded as
`material`, with the reviewer's words (5998876808).

**After the first T5:** O-145 (a), the first pilot is a Fail by the pre-registered
rule, and the result after the fix is its own result; O-146 (a), the three fixes, a
second run from the same baseline commit, and a target task that takes it into the
target's `main`, with a script that shows the tree; the budget rose to 3,200 lines
over 42 files (6002406785). O-147 (a), the target task closes out in its own pull
request (6011287620). The three plan reviews of the target task gave O-148 to O-151:
the copy writes the records of the second run, an exception to their immutability
that the Operator scoped (O-148, O-151); one round with the four lenses on each
frozen head, a workaround that two logins approve (O-149, O-155); the task's two files
(O-150). O-152, the branch names follow the text by a forced update after the merge
(6013939366). O-153, the second approval of R4; O-154 and O-155, its conditions, and the
logins as the identity of the approvers. O-156 and O-157, the cycle cap of the
target task rose from 2 to 3, then to 4. O-158, the last finding of its review is a
known limit.

## What was done

1. **The inputs** (D1 to D5). The Operator posted the S01 answers on the target
   (T1). The author wrote the 25 prose files of the prose step (S07, S08, S09, and
   S14: 20 adapted files by their flagged lines, with four helper sessions, and the
   identity text of `README.md` and `AGENTS.md`). The author proposed the 107
   answers of S10; after six corrections, the Operator posted them on the target
   (T2). The prose review asked three corrections; the Operator then approved the
   prose of `cec749a`.
2. **The setup** (D6): `layup setup` ran S01 to S15 in a work area, with its planned
   stops; `layup setup verify` gave exit 0, 15 `pass` and 3 `clear`. LAYUP refuses a
   changed input of a done step, so each change of an input made a new work area
   (five runs in all, F-2).
3. **T3:** the Operator ran [`t3.sh`](../../runs/T-evad/t3.sh), which runs the five
   commands of `commands.sh` one at a time and stops at the first failure (F-13).
   The author read back `main`, `layup-records` and the ruleset.
4. **D7:** the target's own task `T-a0rt` (its issue 2) changed its table "Which
   checks run", by the target's process (O-142): a plan and its review, two review
   rounds (the first `material`, fixed; the second `nothing material in scope`),
   the target's hooks, a pull request with 13 checks, and a plain merge under the
   ruleset (`7d7b387`). `layup gate` from outside on the recorded base `cec749a` and
   head `ea28ed4`: exit 0, one verdict per kind, the same bytes in a second run.
5. **D8:** (a) the inventory of 66 values, with their sources, and the audit by
   GPT-6 Sol: both planted rows caught, then 66 of 66 supported
   ([`value-audit.md`](../../runs/T-evad/value-audit.md)). (b)
   [`rules-diff.sh`](../../runs/T-evad/rules-diff.sh) on the setup head: PASS, with
   17 fixture cases and 9 mutations caught.
6. **D10:** two defects of LAYUP blocked the pilot and were fixed test first (F-8,
   `fa53cfb` and `39f1eb8`); the binary was built again and the real run used it.
   The report of `rules-diff.sh` was fixed (F-18, `f9fcc9d`). The other findings
   are in [`findings.md`](../../runs/T-evad/findings.md), with their class, and the
   numbers in [`numbers.md`](../../runs/T-evad/numbers.md).
7. **The first T5:** the idea owner accepted 15 requirements and rejected `REQ-002`
   and `NFR-003`: S11 filled a marker of four lines in part (V-040), and the prose
   step removed markers of flagged lines, so a cell got a value with no source
   (V-061) ([`acceptance.md`](../../runs/T-evad/acceptance.md)). So the first pilot
   is a Fail (O-145).
8. **The fix** (O-146), test first: one rule for a marker, also over more lines; S14
   refuses an input that loses a place of a marker; check `markers` names a quote
   with no pair (`c68b390`, `cc31402`, `006fbee`). Rehearsal A showed the refusal (27
   places lost), and rehearsal B the seven places right
   ([`test-runs.md`](../../runs/T-evad/test-runs.md)).
9. **The second run,** from the same baseline commit, with the corrected prose and
   14 new answers of S10: `verify` exit 0, 15 `pass` and 3 `clear`; `rules-diff.sh`:
   PASS. The Operator pushed `layup-setup-2` and `layup-records-2` (T6). D8 (a): 67
   of 67 values supported, after an inventory made from the diff (F-29); the
   specification now says how `answers.sha256` is computed (F-30, `c64f9df`).
10. **The target task `T-vu2j`** (condition 2 of O-146;
    `pharzam/chat-orchestrator#6`): by the target's own process, its `main` took the
    tree of the second run and kept the change of `T-a0rt`. Three plan reviews
    (`reject`, `reject`, `approve-with-conditions`), the two approvals of O-149, and
    five review rounds (4, 8, 3, 2 and 1 material findings; F-37, F-38); pull
    request `pharzam/chat-orchestrator#7`, merged as `907019f`.
    [`tree-equal.sh`](../../runs/T-evad/tree-equal.sh) passed on the frozen head,
    after the close-out and on the merge commit
    ([`tree-equal.txt`](../../runs/T-evad/tree-equal.txt)); `layup gate` from
    outside, with a clean build of `c961652` (F-34), gave exit 0 and five `clear`
    verdicts; the 13 checks of the pull request passed. The other scripts of that
    task stay on the pilot's host; [`test-runs.md`](../../runs/T-evad/test-runs.md)
    names them, with their SHA-256 and their tests.
11. **The second T5,** from the login `pharzam` (6023071943): the 67 values of the second
    run checked with no row rejected, and each of the 17 requirements accepted; three of
    them (`NFR-007`, `REQ-015`, `REQ-017`) if two checks pass at the freeze, and they pass.
12. **D11:** the traceability row of the pilot, the §12 Test cells of `PRD-0001` for
    the nine requirements that the plan's table "What phase 1 proves" gives to row
    20 (`NFR-004` and `NFR-006` keep the tests of rows 5, 7, 9, 10 and 15), a §13 row,
    and five pitfalls in `guardrails.md` §2.

**The rejected alternatives:** a D7 change that skips the target's process and
reads pass item (4) as "the required checks pass" (O-142 b); two rounds on one head,
which fail the target's `review-record` check (O-143 b, finding F-15); a first task
that changes that check (O-143 c); a fix of the catalog's `static` command in this
task (the Operator's ruling on F-1); a fix of `commands.sh` in this task (it changes
S13 and the approved target commit; F-13); a copy that keeps the records of the
first run (O-148 b); the findings of rounds 3 and 4 of `T-vu2j` settled as notes
(O-156 b, O-157 c); a sixth round on a new form of `t7-3.sh` (O-158 b).

**Known limits:** F-1 and F-11 (the Operator's rulings). The changes of D7 and of
`T-vu2j` have no Go file, so their gate runs show no active gate on product code,
and no merge was shown blocked by a failed job; the setup's known-bad fixtures and
the ruleset's read-back are the evidence of phase 1 (condition 2 of the plan
review). `rules-diff.sh` reads the setup head; the later commits of `main` are the
project's own (F-18). O-158: `t7-3.sh` can stop between the merge and the branch
moves; it did not stop. The binary of the second run, `layup5`, was built from
`006fbee` with `vcs.modified=true`. O-154: if the review of this task changes fixes 1 to
3, a setup run with the merged binary must give the tree of `layup-setup-2`, except the
values of `pin.time`, or `T-vu2j` opens again. Round 1 changed fix 1; a setup run with a
clean build of `c285d45`, whose Go code the merge keeps, gave that tree byte for byte
(`evidence/187-run3-setup.txt` on the pilot's host). The issue of the ADR that
ends bootstrap mode, #120 (`T-gd8q`), lists the findings that are not fixed.

**Lessons:** five pitfalls in `guardrails.md` §2: a one-line form that loses a
failure (F-9), a claim of a check that nothing runs (F-5, F-10, F-12), a command file
that goes on after a failure (F-13), a value audit that reads only the values it is
given (F-28, F-29), and a claim of a later state in the present tense (round 4 of
`T-vu2j`).

## Verdict

The first pilot is a Fail (O-145). The result after the fix, which O-145 records as its
own result, holds pass items (1) to (7) by the second T5 (6023071943), and item (8) at the
freezes `6dcd530` and `25c6567`, with `TestPackageRules` and `release-check.sh`, on which
the acceptance of `NFR-007`, `REQ-015` and `REQ-017` depends. Review: round 1 (`6dcd530`,
cycle 0) was `material`: S11 could give two places of a joined line one name; fixed test
first in `c285d45`, and the setup ran again with the fix (O-154). Round 2 (`25c6567`,
cycle 1) gave `nothing material in scope`, with 15 notes; the reply 6024285307 on #97
applies nine and declines five with their reasons. The diff is 3,188 lines over 40 files,
inside 3,200 over 42 (O-146).

## Resource record (ADR-0007: recorded, not budgeted)

Recorded, not budgeted. Times are UTC; a token count that a harness does not give is
`not reported`.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan and its review | reasoning | Claude Opus 5.5; GPT-6 Sol on the Devin CLI | max; `xhigh` | not reported | 34 min 45 s on 10-04; the review 5 min 29 s |
| The first run: the inputs, the setup, T3, D7 and D8 | execution; the reviews and audits reasoning | Claude Sonnet 5.5 (to T2 revision 2, with four helper sessions), then Claude Opus 5.5; GPT-6 Sol on the Devin CLI | not recorded, then max; `xhigh` | not reported | about 4 h 15 min of work on 10-05 ([`numbers.md`](../../runs/T-evad/numbers.md)) |
| The fix, the second run and its audit | execution; the audit reasoning | Claude Opus 5.5; GPT-6 Sol on the Devin CLI | max; `xhigh` | not reported | about 1 h 50 min on 10-05 and 10-06; the audit 6 min 15 s and 6 min 8 s |
| The target task `T-vu2j` | execution; the reviews reasoning | Claude Opus 5.5; GPT-6 Sol on the Devin CLI | max; `xhigh` | not reported | about 7 h 45 min on 10-06; three plan reviews 24 min 57 s, five rounds 62 min 54 s |
| The close-out and the freeze | execution | Claude Opus 5.5 | max | not reported | 18:45 to 19:02 on 10-06 |
| Review round 1, first try: each harness skipped | reasoning | GPT-6 Sol on the Devin CLI (no output in 300 s); Grok 4.7 on the OpenCode CLI ("Go usage limit exceeded"); Claude Fable 5.1 on Claude Code (no record in 900 s) | `xhigh` | not reported | 19:02:55 to 19:23:33 |
| Review round 1 | reasoning | Claude Fable 5.1 on Claude Code | `xhigh` | 1,893,437 (USD 5.87) | 6 min 36 s, from 19:24:54 |
| The fix of round 1, and the setup run of O-154 | execution | Claude Opus 5.5 | max | not reported | 19:32 to 19:42 |
| Review round 2 | reasoning | Claude Fable 5.1 on Claude Code, with two helper sessions | `xhigh` | 3,377,391 (USD 11.39) | 8 min 56 s, from 19:42:44 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | from 19:52 |
| **Total** | | | | not reported in full | 10-03 16:15 to 10-06 about 20:05: about 75 h 50 min; about 15 h 30 min of it was work, the rest waits for the Operator, with three nights |
