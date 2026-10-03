# T-d6q5 — steps S12, S13 and S15, and the checks `jobs` and `gate:<kind>`

Issue: [#92](https://github.com/pharzam/layup/issues/92), row 15 of the
[implementation plan](../plan/README.md); child of `layup setup` (#77); with
the wording of #34. Serves `F-0003#42` and `F-0003#47`. Base `25fcf8c` (the
merge of row 13, #114). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-d6q5/`](../../runs/T-d6q5/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #92. The
reviewer order is the Operator's (Devin, then OpenCode, then Claude). Devin
(its usage quota) and OpenCode (Grok 4.7: no output in five minutes) gave no
record. The plan review (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session in a clone at `25fcf8c`) gave
`approve-with-conditions`: Budget maximum 3,600 lines added plus removed over
40 files against the base, close-out inside; Cycle cap 1 (2 if the branch
changes `ci.yml` or `.githooks/`; it changes neither). Its three conditions
are applied: a fixture that does not apply is `not-active` (1); the required
checks of S13 are the gate jobs only, and the row S13 of `setup.md` names the
decision that changes `steps.tsv` S13, with its reason and a note for the
Operator (2); the fixed form of the two JSON bodies, with the evidence of the
app ID, and a hash row of each file (3). Notes 1 to 7 and 9 are applied; note
8 (a Status pointer on ADR-0011) is declined: an amendment of an accepted
decision is a decision of its own, not wording.

**For the Operator** (condition 2; no stop, as the first pilot, row 20,
applies the ruleset): with this rule, the baseline's own CI jobs (its linters)
run on each pull request of a target and block no merge. An objection before
row 20 changes the row S13.

## What was done

1. **S12:** the files of the catalog entry of the stack, with the module path
   `github.com/OWNER/NAME`, the manifest `docs/gates.tsv`, the row of
   `open-gaps.tsv` of each gap of the entry, and a record row of each file,
   gap and the module path; a path of the entry that the head holds with
   other bytes is `fail` (`REQ-018`). Its evidence is checks `jobs` and
   `gates`.
2. **S13:** `docs/setup/branch-protection.json` (the keys of LAYUP's own file)
   and `out/ruleset-default.json` (the body of a repository ruleset), each
   check pinned to the GitHub Actions app 15368, with a hash row of each file;
   its commands are the push of `layup-setup` and the apply of the ruleset (a
   hand-off).
3. **S15:** the stop `O-verify` with no `out/verify.tsv`; a table of another
   form is exit 2, and a row that is not `pass` or `clear` is `fail`; the
   rule-path register with its sources; and, through a new hook of the
   runner, the first commit of the orphan branch `layup-records` with the four
   files, in a scratch work tree, so `main` and `layup-setup` do not move. A
   records branch of a run that stopped is taken only when it is S15's own
   commit.
4. **`internal/verify`:** check `jobs` (the reader of `check_protection`, with
   no YAML library), and the rows `gate:<kind>`: the clean run of
   `internal/gate` on the setup head, and one run on a commit of each active
   kind's fixture, made in a scratch work tree on no ref; the first rule that
   matches decides. Each check of the table is built now, so a correct setup
   exits 0.
5. **`internal/work`:** the schema of the table of `layup setup verify` (K9).
6. **The tests:** the unit tests of each step, of the register, of the records
   commit and of the hook (a fake of `internal/git`), of the job reader and of
   the decision table of the gate rows (a stand-in of the gate); the
   integration tests of the gate rows on a real target, and of S01 to S15
   through `internal/cli` on a stand-in baseline with real `git` and the Go
   toolchain, with `layup setup verify` and a rerun; the e2e tests of the two
   commands through the binary; 33 mutations, each detected. The stand-in of
   the verify tests got a `go.mod`, the command of the Go entry and a
   workflow, so its fixture fails for its own reason.
7. **The documents:** `setup.md` (the rules of S12, S13 and S15 and of the two
   checks, the text of the records README, the sources of the register, the
   hook of a step, the known limits), `packages.md`, `records.md`, the
   glossary row `Rule path`, the traceability rows, the `PRD-0001` cells with
   a §13 row, two lessons in `guardrails.md` §2.
8. **#34:** point 1, the glossary says that LAYUP orchestrates the role
   agents, which a harness agent runs; point 3, the backlog line of `T-stfn`
   says where the name "Step 2" comes from; point 5, the glossary has the row
   `Rule path`. Point 4 needs no change: S04 writes the decision record of the
   pin for each target since row 9, the sentence "for ADR-0011 to decide" is
   not in `docs/setup/README.md` any more, and the list of the Invariant 3
   complement is the rule-path register that S15 now writes (its output,
   `layup/rules`, is phase 2's, milestone `M2f`). Points 2 and 5's ADR-0011
   parts are known limits (note 8 above): ADR-0011 is accepted, and
   `docs/spec/README.md` gives a register such as `open-gaps.tsv` the form
   `no-header`.

**The rejected alternatives:** a YAML library for check `jobs` (`NFR-007`; the
reader of `check_protection` is the rule); `git commit-tree` and `update-ref`
for the records commit (two new calls of `internal/git`; a scratch work tree
uses the calls that exist); the records commit in the work tree of the target
(it would move `HEAD` and the files of `layup-setup`); one gate run per kind
with a kind filter (a change of the interface of `internal/gate`; the run of
all kinds on a fixture commit costs seconds); a fixture that does not apply as
`fail` (the fixture run did not run, so it is `not-active`, `NFR-004`); the
required checks of the baseline's own jobs in the ruleset (condition 2).

**Known limits** (in `setup.md`): the app ID, `strict` and the review settings
of S13 have no record row of their own (`NFR-003` item 1; the two hash rows
hold the files); a workflow in another form of YAML gives a false `fail` of
check `jobs`; `verify.tsv` names no head and no work area, and judged the
record before the rows of S15.

## Review round 1

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record. Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI, a
fresh read-only session in a clone at `a7a64c9`, its record at 4 min 51 s;
the record is on #92) gave `nothing material in scope`, with six notes; its
verdict cell now reads `material`, with the reason, as CI then found a bug on
that head (below). It
ran the three test levels, the race tests and the local checks on the head,
and read the brief's other inputs against the code and the unit cases. The
notes:

1. K6 had no sentence in `docs/spec/`: applied, as a "decided here" of D8 in
   `setup.md`.
2. The job reader takes a comment off a quoted name before its quotes, so a
   quoted name that holds ` #` loses the text after it: a known limit in
   `setup.md` (a false `fail`, never a false pass; a kind holds no space).
3. S15 does not refuse a `TMPDIR` in the work area, as `layup setup verify`
   does: a known limit in `setup.md` (the scratch tree is removed after the
   commit, and the commit is the same).
4. The green runs of the evidence are on `dde2e66`: the evidence now has the
   ladder on the freeze head `a7a64c9` too.
5. The budget: 2,623 lines over 32 files before the close-out.
6. The limit of the round: the inputs of the brief that it read and did not
   run by hand.

## The CI failure and its fix (cycle 1)

CI of PR #115 failed in the job `tests` on `8a3a565` (the close-out of round
1; Linux, `go1.26.8`, `git` 2.55.0): `TestEachFindingOfATarget` could not
remove a work area (`unlinkat …/.git/objects: directory not empty`), as a
maintenance of `git` that a commit had started in the background (a repack,
on `git` 2.55.0) wrote into it. The same test passed on the LAYUP host
(`git` 2.54.0) and in round 1. A bug is material, so by the rule of O-120 (the pitfall
"A cycle cap raised after a last-round verdict" of `guardrails.md` §2) the
verdict cell of round 1 is edited to `material`, with the reason, and the fix
is cycle 1, inside the cap of 1:

1. Each call of `internal/git` starts with `-c maintenance.auto=false`
   (`packages.md`), so no call of `layup` leaves a process of `git` in a
   repository after it ends; the two test helpers that commit with a plain
   `git` (`gitOut`, `fixtureGit`) and the control of
   `TestAHostileHostChangesNothing` do the same.
2. `TestNoCallStartsTheMaintenance` (integration): a repository whose own
   configuration asks for the loose-objects task at once, in the foreground,
   keeps its loose objects after a `Commit`; its control, a plain
   `git commit`, packs them.
3. The lesson "A background maintenance of git" in `guardrails.md` §2.

The red and the green runs are in the evidence. Round 2 reviews the fix.

## Review round 2

Devin (its usage quota) and OpenCode (no output in five minutes) gave no
record again. Round 2 (Claude Fable 5.1, effort `xhigh`, on the Claude Code
CLI, a fresh read-only session with another lens, cleanup, determinism and
leftovers, on `cd3daa8`, its record at 6 min 22 s; the record is on #92) gave
`nothing material in scope`, with seven notes. It read each automatic start
of `git maintenance` against the fix, and ran the three test levels, the race
tests and hand runs of `git` 2.54.0. The notes:

1. The version of the automatic repack: on `git` 2.54.0 a commit's automatic
   maintenance did not repack in its hand run, so "since 2.54" was not shown;
   applied: the documents and the two comments name the CI run on 2.55.0 and
   say that the first version is not shown.
2. Two helpers and one control, not "four test helpers": applied.
3. "CI of #92" and "CI of #115" named one run: applied ("CI of the pull
   request #115").
4. When the removal of the scratch work tree of S15 fails, the entry of the
   work tree stays in the target until `git worktree prune`; the step fails
   and says it: a known limit in `setup.md`.
5. LAYUP's own `nested-checkout-check.sh` commits with a plain `git` in a
   temporary directory, under the thresholds of the maintenance; it is off the
   path of this task, and `layup` does not run it.
6. The budget: 2,802 lines over 37 files before the close-out.
7. The limit of the round: the work tree and the refs of a stand-in after a
   run rest on the assertions of the tests.

## Verdict

Delivered: steps S12, S13 and S15 of `layup setup`, and the checks `jobs` and
`gate:<kind>` of `layup setup verify`. S12 writes the files of the catalog
entry of the stack with the module path of the target, the manifest and the
gap of the coverage floor, with a record row of each; S13 writes the
protection file and the ruleset of the default branch, each required check a
gate job pinned to GitHub Actions, and hands their commands to the Operator;
S15 stops for the table of `layup setup verify`, writes the rule-path register,
and makes the first commit of the orphan branch `layup-records` through a new
hook of the runner, so the committed record is the record of the run. Check
`jobs` reads the job names as `check_protection` does; each row
`gate:<kind>` runs the gate on the setup head and on a commit of the kind's
fixture, on no ref. So a correct setup of phase 1 runs from S01 to S15 and
`layup setup verify` exits 0. No call of `internal/git` starts the maintenance
of `git` any more. The wording of #34 is done, with its ADR-0011 parts as
known limits.

The plan review (Claude Fable 5.1) gave `approve-with-conditions`, with three
conditions, applied. Round 1 (`a7a64c9`) named no material finding, with six
notes (note 1 applied, notes 2 and 3 known limits); CI of PR #115 then failed
on a background maintenance of `git`, so by O-120 the verdict of round 1 reads
`material`, and the fix is cycle 1 (`5eb2a76`, `cd3daa8`). Round 2 (`cd3daa8`)
gave `nothing material in scope`, with seven notes: notes 1 to 4 applied as
text, notes 5 to 7 recorded. The records are on #92. At `cd3daa8`, the local
checks, `go build`, `go vet` with each tag, `gofmt`, the three test levels,
`go test -race` on six packages (and on four with `-tags=integration`),
`run.sh` and the discipline tests pass, and the job `tests` of CI passes on
Linux with `git` 2.55.0; 33 mutations of the code of row 15 are detected;
`review-record-lint` passes on the comments of #92 (2 rounds, cap 1). The
close-out commit changes the text of documents and two comments of Go files,
and no code. The diff against `origin/main` is 2,908 lines added plus removed
over 37 files with the close-out, inside the Budget maximum of 3,600 lines
over 40 files.

Next: row 16 of the plan (`T-dep6`, #93), whose After cell (row 15) is then
merged.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC, 2026-10-02 to 2026-10-03. A
token count is the `result` event of the Claude Code CLI (input, output, cache
creation and cache read tokens, and its cost) where that harness gave one;
`not reported` where the harness or the author's session gives none. The
author's session was summarized once in the code part; that changes no part.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan (D1 to D10) | reasoning | Claude Opus 5.5 | max | not reported | 23:42 to 23:46 |
| The plan review, first and second harness: skipped (Devin's usage quota; OpenCode no output in five minutes) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 23:46:28 to 23:52:46 |
| The plan review | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 758,232 (USD 3.79) | 5 min 34 s, from 23:52:52 |
| The answer to the plan review | reasoning | Claude Opus 5.5 | max | not reported | 23:58 to 23:59 |
| The tests first, the code, the integration and e2e tests, the mutations, the documents and the evidence; the freeze | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 23:59 to 00:50 |
| Review round 1, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 00:50:35 to 00:56:53 |
| Review round 1 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 1,150,935 (USD 3.53) | 5 min 0 s, 00:57:01 to 01:02:01 |
| The close-out of round 1, and the pull request | execution | Claude Opus 5.5 | max | not reported | 01:02 to 01:06 |
| The CI failure: its cause, the fix test first, the freeze | execution | Claude Opus 5.5 | max | not reported | 01:07 to 01:17 |
| Review round 2, first and second harness: skipped (the same) | reasoning | GPT-6 Sol on the Devin CLI; Grok 4.7 on the OpenCode CLI | `xhigh`; `xhigh` | not reported | 01:17:06 to 01:23:25 |
| Review round 2 | reasoning | Claude Fable 5.1 on the Claude Code CLI | `xhigh` | 738,533 (USD 3.52) | 6 min 39 s, 01:23:30 to 01:30:09 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 01:30 to 01:40 |
