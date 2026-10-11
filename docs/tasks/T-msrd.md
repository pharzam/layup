# T-msrd — the technical specification of M2d, the Scaffold

Issue: [#265](https://github.com/pharzam/layup/issues/265). Milestone `M2d` of the
[implementation plan](../plan/README.md#milestones), which starts with its
specification task (O-114), taken by the Operator's directive ("go ahead, I
need to finish layup and reach first version") while row 49b of `M2c` (#228)
waited for the Operator's inputs. Serves `F-0003#51` (`REQ-012`) for `REQ-002`,
`REQ-003`, `NFR-001`, `NFR-002`, `NFR-003` and `NFR-006`. Base `4ccbb12` (the
merge of #263). Author: Claude Opus 5.5 on Claude Code. Evidence:
[`runs/T-msrd/`](../../runs/T-msrd/).

## Plan and plan review

The plan (R12), its review and the author's answer are comments on #265. The
plan review (Claude Fable 5.1, effort `xhigh`, a fresh read-only session in a
clone at `4ccbb12`, 10 min 4 s) gave `approve-with-conditions`: Budget maximum
3,400 lines added plus removed over 30 files, close-out inside; Cycle cap 1; no
panel. Its five conditions are applied:

1. a section for `NFR-003`;
2. the four sections of the inventory that D4 had left out;
3. the six Tests rows of the inventory in the table of acceptance tests;
4. a built block changed in words only;
5. the role of the prose rows decided here.

Of its notes, note 4 settled decision 2 inside ADR-0014 and ADR-0017, and note 5
gave the `integration_id` of the `layup/` checks from `app_id` of the forge
register, so no marker was needed.

## What was done

1. **S1, the public-solution search**
   ([`runs/T-msrd/search/`](../../runs/T-msrd/search/), commit `5cd2888`,
   before S3):
   - five curated lists, read in full, and 20 GitHub searches with a result
     (eight long queries gave none and were run again in short forms, recorded);
   - two blind searchers: GPT-6 Sol on Devin, a model of the reasoning tier, as
     Devin's paid models answered again, and Claude Sonnet 5.5 on Claude Code
     (OpenCode's login failed).

   Both conclude that no public solution can be the base of `M2d`: what
   transfers is GitHub's rules and activity API, with the limits that its
   documentation states. Searcher A's finding that a push to `layup-probe` is a
   creation gave the ruleset its rule `creation`.
2. **S2:** `runs/T-w73g/survey.md` has no section `M2d`, so no row is taken or
   rejected.
3. **S3, the decisions** (comment 6110224778 of #265):
   - **decided by the Operator:** decision 1, the adoption of a phase-1 target,
     **O-200** (a): none;
   - **decided here, with their reasons:**
     - the records-branch ruleset applied at Scaffold 6;
     - the setup steps in phase 2;
     - `layup run` driving the step runner in its own process, with a work area
       built from the records;
     - the setup questions and the proposed marker sources in the one batch,
       with the step `markers`;
     - the prose rows by a `product-owner` session;
     - the two rulesets, with F-13 and F-19;
     - the forge check and its repeat at each later start;
     - the findings of the first pilot;
     - the demo and the release review as rows of the slicing task.
4. **T, the tests first** ([`test-runs.md`](../../runs/T-msrd/test-runs.md)):
   - `sections.sh` red on each of its 31 headings before the files existed;
   - the block test red with the three new names in `notYetBuilt`;
   - both green after D1 and D2.
5. **D1:** [`docs/spec/scaffold.md`](../spec/scaffold.md) (new):
   - the setup questions in the one batch;
   - the setup under `layup run`, its stops and its prose rows;
   - the rulesets of Scaffold 6 and the Operator's commands, fail-fast and safe
     to run again;
   - Scaffold 7, the forge check, and its repeat at each later start;
   - `NFR-002`, `NFR-003` and `NFR-006` at the Scaffold;
   - the adoption (none);
   - the findings of the first pilot;
   - the steps, the input states and the acceptance tests.
6. **D2:** [`records.md`](../spec/records.md):
   - the blocks `markers`, `setup-tree` and `forge-check`, in `notYetBuilt`;
   - the layout rows changed in place;
   - the list of the built blocks that `M2d` changes in words;
   - NFR-001 "Not in phase 1" and NFR-002 item 3.
7. **D3:** [`packages.md`](../spec/packages.md):
   - the table of `M2d`, with no new package;
   - the May import cell of `internal/run`;
   - the job cells of `internal/setup` and of the GitHub adapter;
   - the call `CommitTree` of `internal/git`.

   [`forge.md`](../spec/forge.md): the calls of `M2d`, and the first use of
   `Activity` moved to `M2d`.
8. **D4:**
   - `README.md`: the file, the exit code 3 and the milestone;
   - `setup.md`: the boundary rows, S13, `commands.sh`, the verify row and "Not
     in phase 1";
   - `intake.md`: the Operator's block, the raw fact, the steps and "Not in M2c";
   - `run.md`;
   - the plan: the `M2d` row, K41, and F-2, F-6, F-13, F-16, F-19, F-31 and
     F-35;
   - the Task cells of `NFR-001`, `NFR-002`, `NFR-003` and `NFR-006` in
     `PRD-0001` §12.

**The open questions of `later-p2-scaffold`**, settled:

- the adoption path, by the Operator: none, O-200;
- the time of the records-branch ruleset: Scaffold 6, decided here inside
  ADR-0014 decision 6 and ADR-0017 decision 3.

## Review rounds

The records are comments on #265. Round 1 was Claude Fable 5.1, effort `xhigh`.
Fable then reached its usage limit, so rounds 2 to 5 were GPT-6 Sol on Devin, by
the order of reviewer selection. A one-word probe of Fable before rounds 3, 4
and 5 answered "You've reached your Fable limit".

GPT-6 Sol wrote `runs/T-msrd/search/searcher-A.md`, a file of this change, word
for word. Each of its records says so: its review is independent of the author
of the specification, not of every model whose output the branch holds.

| Round | Head | Cycle | Cap | Material | The findings | Then |
| ----- | ---- | ----- | --- | -------- | ------------ | ---- |
| 1 | `46baf93` | 0 | 1 | 4 | the first command compared the forge's head with `pin.commit`; one ruleset name for both files, unlike phase 1; a failed prose text that no restart wrote again; parts with no acceptance test | **O-202** (a): cap 2 |
| 2 | `8b454dd` | 1 | 2 | 2 | an `M-` row for a marker that the setup deleted, which the runner refuses; a work area rebuilt after the push at the setup head | **O-204** (a): cap 3 |
| 3 | `fdca09b` | 2 | 3 | 3 | a marker text twice in a file; no result file for the proposals; a probe that did not ask for the ruleset's name | **O-205** (a): cap 4 |
| 4 | `97856bd` | 3 | 4 | 2 | `commands.sh` could not be run again after a partial apply; the phase-2 committer of `setup/verify.tsv` | **O-206** (a): cap 5 |
| 5 | `1c91410` | 4 | 5 | 1 | an output file left by a failed apply | **O-207** (b) |

Further detail:

- Round 1 also had five notes, all fixed.
- Round 2 is GPT-6 Sol's first; Fable's try at it ended on the usage limit after
  12 min 47 s, with no record.
- The verdict of each of rounds 1 to 4 at its cap was edited to `material` after
  the Operator raised the cap.

Round 5 keeps its verdict, `not mergeable, findings recorded`. By **O-207** (b),
the Operator settled its finding as fixed by `3171d79`, as O-158 did for
`T-vu2j`, and no round reads that commit.

## Verdict

Delivered: the specification of `M2d`, the Scaffold driven by `layup run`, in
`docs/spec/scaffold.md`, with the three blocks of its records, the calls of
`M2d`, the table of `M2d` and the documents it changed, after a public-solution
search by two blind searchers.

- **How the review ended:** at cycle 4 of a cap that the Operator raised four
  times (O-202, O-204, O-205, O-206). Its end is the Operator's decision
  (O-207), not decay: the material findings per round were 4, 2, 3, 2 and 1.
- **The size:** the diff against `4ccbb12`, the branch's base, is inside 3,400
  lines over 30 files.
- **Next:** the slicing task of `M2d`.

## Resource record

Recorded, not budgeted (ADR-0007). UTC, 2026-10-11. Reviewers' and searchers'
tokens are as each harness reports them; the author's are not reported, and
Devin reports none.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | not reported | not reported | 14:24 (the issue) to 14:25 (the plan's comment); the reading of the sources before the issue is not timed |
| The plan review | reasoning | Claude Fable 5.1 | `xhigh` | 1,914,208 (USD 7.45) | 10 min 4 s |
| S1, the sources | execution | Claude Opus 5.5 | not reported | not reported | 14:36 (the answer to the review) to 14:40 (the searchers' start) |
| Searcher A | reasoning | GPT-6 Sol | — | not reported by Devin | 3 min 20 s |
| Searcher B | reasoning | Claude Sonnet 5.5 (neither tier list) | — | 1,733,482 (USD 1.61) | 2 min 49 s |
| The summary, S3, T, D1 to D4 | execution | Claude Opus 5.5 | not reported | not reported | 14:44 to 14:59 |
| The adoption by O-200 | execution | Claude Opus 5.5 | not reported | not reported | 15:17 to 15:19 |
| Round 1 | reasoning | Claude Fable 5.1 | `xhigh` | 1,338,274 (USD 7.18) | 16 min 4 s |
| The fix of round 1 | execution | Claude Opus 5.5 | not reported | not reported | 15:36 to 15:38 |
| Round 2, Fable's try | reasoning | Claude Fable 5.1 | `xhigh` | no record: the usage limit | 12 min 47 s |
| Round 2 | reasoning | GPT-6 Sol | — | not reported by Devin | 6 min 14 s |
| The fix of round 2 | execution | Claude Opus 5.5 | not reported | not reported | 16:14 to 16:15 |
| Round 3 | reasoning | GPT-6 Sol | — | not reported by Devin | 6 min 8 s |
| The fix of round 3 | execution | Claude Opus 5.5 | not reported | not reported | 20:29 to 20:30 |
| Round 4 | reasoning | GPT-6 Sol | — | not reported by Devin | 7 min 0 s |
| The fix of round 4 | execution | Claude Opus 5.5 | not reported | not reported | 20:42 to 20:43 |
| Round 5 | reasoning | GPT-6 Sol | — | not reported by Devin | 6 min 55 s |
| The fix of round 5 | execution | Claude Opus 5.5 | not reported | not reported | 20:53 |
| The close-out | execution | Claude Opus 5.5 | not reported | not reported | 20:59 to 21:00 |
