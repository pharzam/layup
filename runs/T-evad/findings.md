# The findings of the first pilot

Task `T-evad` ([#97](https://github.com/pharzam/layup/issues/97)), row 20 of the
[implementation plan](../../docs/plan/README.md). The target is
[`pharzam/chat-orchestrator`](https://github.com/pharzam/chat-orchestrator), set up at
`cec749a` from the brief `PSB-CHAT-001`. The IDs F-1 to F-20 are the ones that the
comments on #97 use. The numbers of the pilot are in [`numbers.md`](numbers.md).

The class of each finding follows D10 of the plan, as the answer to the plan review
amends it (comment 5979731752, item 6). A defect on the pilot's path (a step of the
setup, `verify` or the gate run) that blocks the pilot is fixed here, test first. An
on-path defect that does not block is fixed here when it fits in the budget. Else it
is a finding in this file, and the issue of the ADR that ends bootstrap mode lists it.
A defect off the pilot's path opens an issue.

## Fixed in this task, test first

- **F-8: two defects of LAYUP that blocked the pilot.** The rehearsal of S10 to S15
  found them, and each is one Investigate of the pilot rule. (a) S11 gave two markers
  on one line the same record key, and the runner refused the record (`fa53cfb`: the
  key gets the column). (b) Check `sources` compared the column of a gap row with the
  final tree, but S11 fills the other markers of the line first (`39f1eb8`). The
  binary was built again, and the real run used it in a new work area.
- **F-18: the report of `rules-diff.sh`.** For a file that the prose step writes, it
  counted the baseline marker lines that the new text replaced as "marker lines
  filled" (`docs/guardrails.md`: 13, but S11 recorded one row there, a gap). The check
  was right; the report was wrong. Fixed in `f9fcc9d`: the report counts the marker
  rows of S11 (filled and gaps). Case 17 failed first, then passed; the 9 mutations
  are still caught; the setup head passes.

## Known limitations: the Operator's rulings

- **F-1: the catalog gates are not the gate script of R-TEST-05.** The Go catalog's
  `static` command (`gofmt -l .`, then `go vet ./...`) does not select the files of
  R-TEST-05 (the tracked Go files except `vendor/`). In 13 cases, the catalog form
  also fails on `vendor/`, on untracked files and on files in hidden directories, and
  it misses a tracked file whose name starts with a dot. The `test` kind has no
  `-race` and not the environment of R-TEST-05, and the gate jobs run with network.
  S12 does not read the answers of S10. No verdict of the pilot changed: the target
  has no Go file yet. The Operator's ruling (T2 revision 2): no catalog change in this
  task; never report the catalog gates as the full R-TEST-05 gate. The target's
  documents say so.
- **F-11: a file name hides a failure of the inactive lint command.** The approved
  answer `M-fc1f2462` adds a sentinel name when `git ls-files` fails. When the working
  tree holds a file of that name that `gofmt` accepts, the failure is lost (case c13).
  The command is a comment in the hook. The Operator's ruling (5995693012): keep the
  limit, and do not activate the command while it remains.

## On the path, not fixed: for the issue of the ADR

- **F-2:** LAYUP refuses a changed input of a done step (exit 2), by its
  specification. A changed prose file needs a new work area; the pilot made five
  (`rerun.sh`, about one minute each), and the commit IDs change with each run.
- **F-3:** S10 asks a marker that runs over two lines with its first line only (for
  example "for example money in currency, time in hours, counts such"). The only safe
  answer is a gap. A fix changes the question IDs of an approved table.
- **F-6:** `layup setup verify` before S12 exits 2, as `docs/gates.tsv` is missing.
- **F-13:** `commands.sh` of S13 is not fail-fast. Run with `sh`, after a refused push
  it pushed `layup-records`, applied the ruleset and exited 0 (a stand-in run). T3
  used `tools/t3.sh`: checks before the first command, one command at a time, a stop
  at the first failure, and a read-back before the ruleset (9 test cases). A fix
  changes S13 and gives a new target commit, so it is not in this task.
- **F-16:** a status cell of the baseline's enforcement table ("branch protection is
  your step") is stale after T3. The setup adapts only the flagged lines, and this
  line is not flagged. The target's issue
  [`pharzam/chat-orchestrator#3`](https://github.com/pharzam/chat-orchestrator/issues/3)
  holds it.
- **F-19:** GitHub adds two ruleset parameters that S13 does not set:
  `allowed_merge_methods` (merge, squash and rebase) and
  `require_extra_approval_for_unattributed_changes`. The first permits a squash,
  which the target's own rules forbid. The second did not block the pull request of
  D7.

## Found off the path of the pilot

- **F-15: the review-record check and the written review rule disagree.** The
  baseline's `review-record-lint.sh` reads one round per cycle (RR7: the k-th record
  carries cycle k-1). Its written rule wants more than one round on one frozen head
  ("One pass is never enough"). Two true rounds on one head fail the check. In the
  target this is a defect, with its issue
  [`pharzam/chat-orchestrator#4`](https://github.com/pharzam/chat-orchestrator/issues/4).
  LAYUP's copy of the check is the same, but bootstrap mode suspends that rule, so
  the check agrees with the rule in force. The ADR that ends bootstrap mode must make
  them agree before the rule comes back; its issue lists this finding.

## Errors in the author's inputs, found by a review or a check

Each is counted in [`numbers.md`](numbers.md). None is a defect of LAYUP.

- **F-5:** the first table of S10 had four faults that the Operator found and no
  check found: a command that is not the file selection of R-TEST-05, a status for
  the wrong check, questions cut in mid-sentence, and a filled command in a comment
  that reads like a running check.
- **F-7:** an example of the baseline (`F-0007#3`) in the glossary input; check
  `glossary` of S08 refused it, as it must.
- **F-9:** the first lint command lost a failure of `git ls-files` (the status of a
  pipe is that of its last command). The Operator found it. The fix was tested in
  13 cases, red then green, also with GNU `xargs`.
- **F-10:** the prose said that no job runs a file of `docs/ci/`; two workflows run
  its linter scripts. The Operator found it; `checks-claims.sh` now tests 40 claims.
- **F-12:** three faults of the prose, found by the Operator's review: a citation
  rule applied to every answer kind, a slot rule that mixed two kinds of wait, and
  the gate jobs called required while the ruleset was only prepared.
- **F-14:** the author's brief told the reviewer of D7 that round 1 was the only
  round, so the reviewer gave a terminal verdict to a material finding. O-144 (a)
  kept the finding and its words, with the verdict `material`.

## Notes

- **F-4:** S10 asked 107 questions: 5 values with a source, 11 conventions, 1
  status, 27 working rules, 26 unused workflow templates and 37 empty example fields.
  So 63 of 107 (59%) were inert text, and 27 were open decisions of the project.
- **F-17:** the Devin CLI with a new isolated home failed in 12 s ("Model not
  found"); a home that had run before worked.
- **F-20:** the touch points of the Operator; [`numbers.md`](numbers.md) counts them.
