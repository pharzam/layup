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
