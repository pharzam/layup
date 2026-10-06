# T5: the idea owner's acceptance of the first pilot, and of the result after the fix

Task `T-evad` ([#97](https://github.com/pharzam/layup/issues/97)). The idea owner's
answers, copied from their comments on #97 (the pilot rule of the plan review,
pass item 7, and condition 5). This copy writes each marker character as
`&lsaquo;` or `&rsaquo;`, as a marker in this repository is an open gap. The first
answer comes from the login `yaltunbicak`, and the second from the login `pharzam`;
this record does not say whether the two logins are one person.

## The verdict of the first pilot: Fail

By the pre-registered rule of the plan review (O-145 = a): the idea owner rejected
`REQ-002` and `NFR-003`; a value had no supporting source (V-061:
`docs/issue-workflow.md:243`); and the fix of the markers (fix 1) is a second
Investigate of S11, after F-8. The result after the fix is recorded as its own
result, below.

## The first answer (2026-10-05, comment 6002034171)

> ## D8 (a) and T5: my answers
>
> ### D8 (a): checked; V-040 and V-061 rejected
>
> I checked the table against the target at `cec749a`, the baseline at `a959655` and `layup-records` at `2a339bb`, with a mechanical read-back by Claude Opus 5.5 in a separate claude.ai session: each value at its place, the three hashes of `facts.sha256`, the pin tree (= the root tree of `242a205`), the eleven commits, and the diff counts (13 added, 29 changed, 31 deleted). I checked V-055 in the settings of the repository. The other 64 rows: accepted.
>
> - **V-040: rejected.** The marker at `docs/tasks/backlog.md:20` runs over lines 20 to 23 (`&lsaquo;State your exact scheme` … `does not already exist."&rsaquo;`). S11 replaced its first line only. Lines 21 to 23 keep the rest of the placeholder and its closing `&rsaquo;`. The question of S10 also showed the first line only (`F-0002`, answer 67). Cause: a marker ends at the line end when it does not close on its line.
> - **V-061: rejected.** The prose step removed markers from flagged lines with no answer, so S10, S11, `open-gaps.tsv` and this table never saw them:
>   - `docs/issue-workflow.md:243`: the branch-protection cell of `pr-link-lint` is now `require the check before merge`. The ruleset requires only the five gate jobs, and my answer `M-0690d23a` to the same marker on lines 245 to 247 is "not required: only the gate jobs of the Go stack are required checks". The cell is a value with no source, and it is false.
>   - `docs/issue-workflow.md:239`: `require a PR before merge`. It is true for the ruleset, but no record row names it.
>   - `docs/engineering-discipline.md:143-144`: `&lsaquo;name your reasoning-tier models&rsaquo;` and `&lsaquo;name your execution-tier models&rsaquo;` became prose. Their question is lost.
>   - `docs/engineering-discipline.md:883-884`, `docs/tests/test-levels.md:133`, `tests/README.md:34-35`: a marker lost its `&lsaquo;` and kept its `&rsaquo;`, or lost both.
>
>   Cause: the prose step runs before S10 (K16, O-123), and check `adapted` does not require a flagged line to keep its markers.
>
> ### T5
>
> | Requirement | Answer | Reason |
> | --- | --- | --- |
> | REQ-001 | accept | New finding: the one gap of the pilot is a false positive. Line 10 of the brief is `**Implementation stack:** Go; PostgreSQL durable adapter`; G1 reads only `technology stack:`. My answer to `Q-001` repeats line 10. |
> | REQ-002 | reject | The setup wrote a value with no source, and it replaced a marker in part (V-040, V-061). |
> | REQ-004 | accept | Record the finding on F-1 below. |
> | REQ-007 | accept | Record the finding on the gate jobs below. I do not need the blocked merge: it would show GitHub's required checks, not that finding. |
> | NFR-001 | accept | — |
> | NFR-002 | accept | — |
> | NFR-003 | reject | As REQ-002: "zero values without a source" does not hold on the target. |
> | NFR-004 | accept | — |
> | NFR-006 | accept | — |
> | REQ-016 | accept | I confirm that 5980424800 (O-141) and 5994748702 (revision 2 of T2) record my decisions as I gave them in the session. |
> | REQ-018 | accept | — |
> | REQ-009, REQ-011, NFR-005 | accept | — |
> | NFR-007, REQ-015, REQ-017 | accept | Only if the two checks pass at the freeze, as you wrote. |
>
> ### New findings to record
>
> - **The gate jobs (REQ-007).** Each gate job checks out the head of the pull request and runs its `.github/gates.sh`, which reads its `docs/gates.tsv`; on `pull_request`, the workflow file also comes from the pull request. The ruleset needs 0 approvals. So a pull request can weaken its own gate and merge. FT4 holds for `layup gate` only. The template is LAYUP's (`internal/catalog/go/files/.github/workflows/gates.yml.tmpl`), so each target has it. Until `layup/gates` is the required check (M2e, M2g): the jobs take the manifest, the config paths and `gates.sh` from the base, as FT4 does and as the baseline's `ci.yml` does for its own checks; and a change to `.github/**`, `docs/gates.tsv` or `docs/gates/**` needs my review as code owner. The test of M2e: a pull request that weakens its own gate is refused.
> - **F-1, for the first product task.** The three `pending` kinds fail each pull request that changes a `.go` path, so the target takes no Go code before M2f. `static` (`gofmt -l .`) also reads `vendor/`, but the brief selects `go mod vendor` (R-TEST-04) and leaves vendored code out of the format check (R-TEST-05). `test` has no `-race` and no PostgreSQL in its job, and its command does not read its config `docs/gates/coverage-floor.txt`. The first task of the target aligns the gates with R-TEST-05, in a pull request with no Go file.
>
> ### The fix (the pilot rule)
>
> 1. LAYUP, test first: a marker runs from `&lsaquo;` to its matching `&rsaquo;`, also across lines. S10 shows its whole text, and S11 replaces all of it. Fixture: the case of `backlog.md`.
> 2. LAYUP, test first: check `adapted` fails when an adapted file loses a marker of the baseline on a flagged line. The marker then goes to S10 like any other.
> 3. LAYUP: a check that fails on a `&lsaquo;` or a `&rsaquo;` with no pair in the target.
> 4. A rehearsal from the same inputs that shows the seven places right. Then a fix of the pushed target with a new freeze; the lost markers go to S10 for my answers.

## The decisions on the fix (2026-10-05, comment 6002406785)

> ## O-145, O-146 and the budget
>
> **O-145: (a).** Record the first pilot as a Fail by the pre-registered rule. The result after the fix is recorded as its own result.
>
> **O-146: (a)**, with two conditions:
> 1. The second run pins the same baseline commit, `a95965534b14b0bf14ad74da0c9a45b5f4aedf88`. If the newest commit of the baseline has changed, the run stops for me.
> 2. The target task shows, with a script whose output goes into `runs/T-evad/`, that the tree of `main` after its merge equals the tree of the head of `layup-setup-2`, except the four paths of `T-a0rt`.
>
> **The design of the fixes:** accepted, with two points.
> 1. A mention is a code span whose whole content is `&lsaquo;`, `&rsaquo;` or `&lsaquo;…&rsaquo;`. A marker with text inside a code span is a marker, also when the span runs over more than one line, as in `backlog.md:20-23`. Fixtures for both cases.
> 2. Fix 3 reads mentions by the same rule, so `AGENTS.md:132` (``Search for `&lsaquo;` ``) is not a finding. A fixture for it.
> 3. Fix 2: an adapted input keeps each marker of a flagged line byte for byte, not only its presence.
>
> **Budget:** approved: 3,200 lines over 42 files, for fixes 1 to 3, the second run and the target task. The template of the catalog (F-22) stays with M2e.

## The second answer (2026-10-06, comment 6023071943)

The answer to the second T5 (comment 6019305811), which asked for the check of the 67 values of the second run
and the acceptance of each requirement for the result after the fix.

> ## D8 (a) and T5: my answers
>
> ### D8 (a): checked, no row rejected
>
> I checked the table row by row against the target at `e2b402b`, the records at
> `layup-records` (`8bc105c`), the two source comments (5995217834, 6010563347),
> `open-gaps.tsv`, `facts.sha256` and `pin.time` — each value at its place, and
> each source in the column Reference. 67 of 67 hold; no row rejected.
>
> - **The eight whole-file digests (W-036 … W-043).** The record's ref is
>   `sha256 inputs/files/<path>`: the digest of the input prose, which lives on
>   the pilot host. I reconstructed that prose from the target itself —
>   reverting each recorded `S11` substitution (for `tests/README.md`,
>   restoring its two-line marker `&lsaquo;test` … `> directory&rsaquo;`, which the F-0002
>   question shows on one line) — and all eight digests reproduce byte for
>   byte. The prose is approved in comment 6010634567: "Approved: 33 changed
>   lines over 8 files, each of the 27 places back."
>
> - **W-044 (`answers.sha256`).** The value sits at `record.tsv:52`;
>   recomputing it needs `inputs/answers.tsv` on the pilot host, which the
>   target does not hold. I checked it at its sources instead: the routine is
>   written into `setup.md:656-663`, pass 2 reproduced the digest from
>   `sources/answers.tsv`, and its 115 `M-` rows match the 115 facts of
>   `F-0002`.
>
> - **Spec citations.** Every `docs/spec/setup.md` line cited in Reference
>   resolves with the claimed text on this branch (`T-evad` amends it — the
>   routine at 656-663 was written in this task); the citations land on `main`
>   with this pull request.
>
> ### T5
>
> | Requirement | Answer | Reason |
> | --- | --- | --- |
> | REQ-001 | accept | as before: the finding on G1 stands; the second run used the same brief and the same answers of S01 |
> | REQ-002 | accept | My two reject reasons are gone: the fix is landed, test first (`c68b390`, `cc31402`, `006fbee`); second run: `verify` exit 0, 15 `pass` and 3 `clear`; the seven places corrected; D8 (a): 67 of 67, no row rejected; on `main` by the target's own process (`T-vu2j`) |
> | REQ-004 | accept | as before (the finding on F-1 stands), plus a second gate run from outside, on `T-vu2j`: exit 0, one verdict for each of the five kinds; the change has no Go file |
> | REQ-007 | accept | as before (the finding on the gate jobs stands); on pharzam/chat-orchestrator#7 the five gate jobs ran and passed with the 8 other checks, `c3.sh ci` judged all 13, and the ruleset is unchanged (read back today) |
> | NFR-001 | accept | the records of the second run are on `layup-records` (`8bc105c`) since O-152; the first run's are on `layup-records-1` |
> | NFR-002 | accept | as before, plus: the commits of `T-vu2j` passed the target's own hooks; pharzam/chat-orchestrator#7 passed its 13 checks, and CI of `main` passed after the merge; the target's CI does not call LAYUP |
> | NFR-003 | accept | as REQ-002: D8 (a) of this run is 67 of 67 supported, none missing; check `sources`: `pass`; my check (1) completes this row |
> | NFR-004 | accept | as before |
> | NFR-006 | accept | the second run pinned the same commit `a959655` (O-146, condition 1); check `pin`: `pass` |
> | REQ-016 | accept | as before, plus O-145 to O-158 recorded on this issue as I gave them, the two approvals of R4 on chat-orchestrator#6, and my runs of `t6.sh`, `t7-1.sh`, `t7-2.sh` and `t7-3.sh` |
> | REQ-018 | accept | `rules-diff.sh` on `e2b402b`: PASS |
> | REQ-009, REQ-011, NFR-005 | accept | as before |
> | NFR-007, REQ-015, REQ-017 | accept, only if the two checks pass at the freeze | as before: the fix changed Go code, so `TestPackageRules` and `runs/T-efmy/release-check.sh` run on the head of this task's pull request at its freeze |

## The result after the fix

The idea owner checked the 67 values of the second run with no row rejected, and accepted each of the 17
requirements. So the result after the fix holds pass items (1) to (7) of the pre-registered rule. Pass item (8),
the end-to-end tests of the issue, and the two checks on which the acceptance of `NFR-007`, `REQ-015` and `REQ-017`
depends (`TestPackageRules` and `runs/T-efmy/release-check.sh`) run at the freeze of the pull request of this task;
the verdict of [`docs/tasks/T-evad.md`](../../docs/tasks/T-evad.md), which the close-out writes, gives their result.
The verdict of the first pilot stays a Fail (O-145).
