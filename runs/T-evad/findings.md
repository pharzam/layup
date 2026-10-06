# The findings of the first pilot

Task `T-evad` ([#97](https://github.com/pharzam/layup/issues/97)), row 20 of the
[implementation plan](../../docs/plan/README.md). The target is
[`pharzam/chat-orchestrator`](https://github.com/pharzam/chat-orchestrator), set up at
`cec749a` from the brief `PSB-CHAT-001`, and again at `e2b402b` after the fix. The IDs
F-1 to F-30 are the ones that the comments on #97 use. The numbers of the pilot are in [`numbers.md`](numbers.md).

The class of each finding follows D10 of the plan, as the answer to the plan review
amends it (comment 5979731752, item 6). A defect on the pilot's path (a step of the
setup, `verify` or the gate run) that blocks the pilot is fixed here, test first. An
on-path defect that does not block is fixed here when it fits in the budget. Else it
is a finding in this file, and the issue of the ADR that ends bootstrap mode lists it.
A defect off the pilot's path opens an issue.

**The verdict of the first pilot is a Fail** by the pre-registered rule (O-145 = a, the
Operator's comment 6002406785): the idea owner rejected `REQ-002` and `NFR-003` at T5
(F-24, F-25), a value had no supporting source (`docs/issue-workflow.md:243`), and fix 1
is a second Investigate of S11, after F-8. The second run after the fixes has its own
result, in [`acceptance.md`](acceptance.md).

## Fixed in this task, test first

- **F-8: two defects of LAYUP that blocked the pilot.** The rehearsal of S10 to S15
  found them, and each is one Investigate of the pilot rule. (a) S11 gave two markers
  on one line the same record key, and the runner refused the record (`fa53cfb`: the
  key gets the column). (b) Check `sources` compared the column of a gap row with the
  final tree, but S11 fills the other markers of the line first (`39f1eb8`). The
  binary was built again, and the real run used it in a new work area.
- **F-24 and F-25: the markers that the setup lost** (V-040 and V-061, which the
  Operator rejected at T5). (F-24) A marker over more than one line was read to the end
  of its first line, so S10 asked it cut and S11 replaced its first line only: the
  rest of the placeholder stayed in `docs/tasks/backlog.md` (F-3 with its effect). (F-25)
  The prose step turned the markers of 27 places into plain text, so S10 never asked
  them (F-26: the Operator listed 7; the check of fix 2 finds all 27). Fixed by the
  Operator's three fixes (comment 6002034171): one rule for a marker, also over more
  lines, in S10, S11 and the checks (`c68b390`); S14 refuses an input that holds fewer
  places of a marker than the file before it, byte for byte (`c68b390`, `cc31402`,
  `006fbee`); check `markers` names each quote with no pair (`c68b390`). LAYUP's own
  `setup-check.sh` follows the same rule.
- **F-27:** the rule of a task index in `rules-diff.sh` did not allow a filled marker
  over more lines, whose lines become one (`beecf92`).
- **F-18: the report of `rules-diff.sh`.** For a file that the prose step writes, it
  counted the baseline marker lines that the new text replaced as "marker lines
  filled" (`docs/guardrails.md`: 13, but S11 recorded one row there, a gap). The check
  was right; the report was wrong. Fixed in `f9fcc9d`: the report counts the marker
  rows of S11 (filled and gaps). Case 17 failed first, then passed; the 9 mutations
  are still caught; the setup head passes.
- **F-30: the specification did not say how `answers.sha256` is computed.** It said
  "the SHA-256 of those rows as the step read them", so the auditor of the second run
  could not compute the recorded value from `answers.tsv` (pass 1, W-044). The
  specification now gives the routine of the code: each row its fields joined by a
  tab, the empty mark read as an empty field, the rows sorted in byte order and joined
  by a line feed, with no line feed at the end (`c64f9df`). The routine gives the
  recorded value, and pass 2 of the audit computed the same value.

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
- **F-21:** the one gap of S01 is a false positive (the Operator, T5): rule G1 reads
  only `technology stack:`, and line 10 of the brief says `**Implementation stack:**`.
- **F-22: a pull request can weaken its own gate and merge** (the Operator, T5). Each gate
  job runs the head's `.github/gates.sh` with the head's `docs/gates.tsv`, the workflow
  file of a `pull_request` comes from the pull request too, and the ruleset asks for 0
  approvals. `layup gate` (FT4) reads the base's gate files; the jobs do not. The
  template is LAYUP's (`internal/catalog/go/files/.github/workflows/gates.yml.tmpl`), so
  each target has it. The Operator's direction, for `M2e`: the jobs read the manifest,
  the config paths and `gates.sh` from the base; a change to `.github/**`,
  `docs/gates.tsv` or `docs/gates/**` needs the Operator's review as code owner; and the
  test of `M2e` is a pull request that weakens its own gate and is refused. The template
  stays as it is in this task (comment 6002406785).
- **F-23: F-1 for the first product task of the target** (the Operator, T5). The three
  `pending` kinds fail each pull request that changes a `.go` path, so the target takes
  no Go code before `M2f`. `static` also reads `vendor/`, while the brief selects `go mod
  vendor` (R-TEST-04) and leaves vendored code out of the format check (R-TEST-05).
  `test` has no `-race` and no PostgreSQL in its job, and its command does not read its
  config `docs/gates/coverage-floor.txt`. The first task of the target aligns the gates
  with R-TEST-05, in a pull request with no Go file.
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
- **F-28: the audit of D8 (a) missed V-040 and V-061.** Both passes of the model audit,
  and the author's own check, read each value at its place. None checked that a filled
  marker was gone whole, or that each marker of the baseline was kept. The Operator's
  row-by-row check found both. The checks of fixes 2 and 3 now cover both cases.
- **F-29: the first inventory of the second run was not complete.** It matched the
  places of the two runs by their values, so it put 15 new places on old lines, where
  the same value stood in the first run. It also did not list the values that no row
  of S11 names: the 7 date places, 2 lines of `facts.sha256` and the 14 new facts of
  `F-0002`. The completeness check of the audit found it (pass 1). The inventory of
  pass 2 reads the new places from the diff of the two setup heads, and its own check
  finds no added line without a row.

## Notes

- **F-4:** S10 asked 107 questions: 5 values with a source, 11 conventions, 1
  status, 27 working rules, 26 unused workflow templates and 37 empty example fields.
  So 63 of 107 (59%) were inert text, and 27 were open decisions of the project.
- **F-17:** the Devin CLI with a new isolated home failed in 12 s ("Model not
  found"); a home that had run before worked.
- **F-20:** the touch points of the Operator; [`numbers.md`](numbers.md) counts them.
