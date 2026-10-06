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
frozen head, a workaround that two operators approve (O-149); the task's two files
(O-150). O-152, the branch names follow the text by a forced update after the merge
(6013939366). O-154 and O-155, the conditions of the second approval, and the
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
11. **The second T5:** the idea owner checked the 67 values of the second run with no
    row rejected, and accepted each of the 17 requirements (6023071943).
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
`006fbee` with `vcs.modified=true` (O-154): if the review of this task changes fixes
1 to 3, a setup run with the merged binary must give the tree of `layup-setup-2`,
except the values of `pin.time`, or `T-vu2j` opens again. The issue of the ADR that
ends bootstrap mode, #120 (`T-gd8q`), lists the findings that are not fixed.

**Lessons:** five pitfalls in `guardrails.md` §2: a one-line form that loses a
failure (F-9), a claim of a check that nothing runs (F-5, F-12), a command file
that goes on after a failure (F-13), a value audit that reads only the values it is
given (F-28, F-29), and a claim of a later state in the present tense (round 4 of
`T-vu2j`).

## Verdict

Written after the freeze checks and the review round of this pull request.

## Resource record (ADR-0007: recorded, not budgeted)

Written after the review round of this pull request.
