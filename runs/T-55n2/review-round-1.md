## Review record — round 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `a88f99f81f616fb8d974803f432a33bdd432f224` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria |
| Briefed on | `.review-in/brief.md`; `.review-in/issue-76-body.md` (AC 1 to AC 12); `docs/tasks/T-55n2.md` (O-121 to O-126); the whole diff `git diff 7cdd346 HEAD` (20 files); `runs/T-55n2/` (`inventory.md`, `plan-check.py`, `plan-check.md`, `self-check.md`, `gen-issues.py`, `fill-hosts.py`, `test-runs.md`); the issue texts in `.review-in/issues/` (#29, #33, #77 to #97); `docs/spec/` (all five files); `PRD-0001` §6, §7.1, §8, §9, §12, §13; `docs/architecture.md` §5, §6, §8, §14; ADR-0012; ADR-0016; `docs/tests/` (`test-levels.md`, `traceability-template.md`, `dod-checklist.md`); `docs/engineering-discipline.md` § Bootstrap mode; `docs/issue-workflow.md` R2, R3, R11, R12; `docs/templates/github/ISSUE_TEMPLATE/task.md` |
| Barred from | the issue comments of #76 |
| Independence claimed | A fresh session, read-only. The reviewer's model (Claude Fable 5.1) differs from the authors' model (Claude Opus 5.5: the author's session, the inventory agents, the self-check critics and refuters). No comment of #76 was read. No `gh`, no network. The reviewer ran the checks and the scripts from the worktree. `gen-issues.py` wrote its output to a directory outside the repository. `docs/tasks/T-55n2.md` names Claude Fable 5.1 as the plan reviewer; that session wrote no line of the change, and this review did not read its comment. No file of the repository was changed. |
| Cycle | 0 |
| Verdict | `material` |

### Raw findings

1. **material** — AC 2 asks for one demo sentence per task; the plan and the issues have none.
   - Where: `.review-in/issue-76-body.md:42`; `docs/plan/README.md:94` (the header row of the task table) and `:96` (row 1); `.review-in/issues/78.md:5`; `docs/issue-workflow.md:114-116`.
   - Cited sentence (AC 2): "Each phase-1 task has a task ID, its issue, one demo sentence, its requirement IDs and In-Scope fact, its tests by level, a size class, its predecessors and its parent." The rule R11 defines the term: "The plan names the one demo — what a reader will be shown when this issue is done, in a single sentence, without an "and" joining two outcomes."
   - Basis: The task table has the columns `#`, `Task ID`, `Issue`, `Task`, `Parent`, `Items`, `Requirements`, `Fact`, `Tests`, `Size`, `Lines`, `After` and `Cap`. No column holds a demo sentence. The `Task` cell is a title that names the scope of the row. Ten of the twenty cells join two outcomes with "and" (rows 1, 2, 3, 7, 8, 9, 13, 14, 15 and 20). The issue texts that the code writes from the table carry the same title as their goal: "The goal: Records: `internal/tsv`, the field rule, and the test that reads every schema block of `docs/spec/`." (#78:5). No issue of #78 to #97 names a demo. The milestone table has a column "First demo" (`docs/plan/README.md:49`); the task table has no such column. The Goal of #76 also says that a reader finds each issue "with its demo, tests, size class and predecessors". AC 2 is not met: a failed acceptance criterion (Bootstrap mode rule 3).
   - Fix: Add a column "Demo" to the task table, with one sentence per row in the form of R11. Let `gen-issues.py` write it into the Goal of each issue, and edit the posted issues #78 to #97 with it.

2. **note** — The §13 row says "the Task column only", and the diff also changes a Test cell.
   - Where: `docs/prd/PRD-0001-layup.md:252` and `:243`.
   - Cited sentence: "REQ-001 to REQ-018, NFR-001 to NFR-007 (the Task column only)".
   - Basis: The diff changes the Test cell of NFR-007 from "CI job lint (gofmt, go vet)" to "planned: T-2tc2/integration/package-rules (...)". The change-log row and the change do not agree. The edit itself is correct (self-check `claims 24`: the old cell did not test the criterion). The §12 preamble (`:214`) says that each delivering task fills the Test column, and T-55n2 does not deliver NFR-007; the In-Scope item 4 of #76 names "§12 Task column and a §13 row". The §13 row is a record, not an operative sentence, so this is a note.
   - Fix at close-out: write "the Task column, and the Test cell of NFR-007" in the §13 row.

3. **note** — The §12 cell of REQ-004 names a task that the table "What phase 1 proves" does not name for REQ-004.
   - Where: `docs/prd/PRD-0001-layup.md:214` (the rule) and `:222` (the cell); `docs/plan/README.md:157`.
   - Cited sentence: "fills the Task column: in each row, the tasks and the milestones that its table "What phase 1 proves" names for that requirement".
   - Basis: The REQ-004 cell names `T-evad` (row 20). The REQ-004 row of the table names rows 5, 14 and 15 only. One of the two is incomplete. The pilot runs `layup gate` on its target (the Task cell of row 20; ADR-0012 part 6), so the table row is the one to extend.
   - Fix: add "the gate run from outside on the pilot's target (row 20)" to the REQ-004 row of the table.

4. **note** — The register preamble and the glossary say "K1 to K41"; the register holds K42.
   - Where: `docs/plan/README.md:216` and `:262`; `docs/glossary.md:148`.
   - Cited sentence: "Each conflict of the inventory (K1 to K41), with the task that settles it and the tasks that read the settled text."
   - Basis: K42 is a finding of the self-check (`:262`, "Found by the self-check"), not a conflict of the inventory. The preamble and the glossary row ("for each conflict K1 to K41 of the plan's inventory") do not name it. `plan-check.py` prints "42 register rows".
   - Fix: "K1 to K41 of the inventory, and K42 from the self-check".

5. **note** — The glossary names a path that no document fixes and that does not exist.
   - Where: `docs/glossary.md:144`.
   - Cited sentence: "its record is `docs/pdr/PDR-0001.md` (task `T-4wrw`)".
   - Basis: No other file of the tree names `docs/pdr/` or `PDR-0001` (a search over `docs/`, `README.md` and `AGENTS.md`). The directory does not exist. `T-4wrw` has no task file and no backlog line yet. `AGENTS.md` (Placeholders): "never invent a command, path or number".
   - Fix: "its record is the PDR record that task `T-4wrw` writes", and let `T-4wrw` fix the path.

6. **note** — Text written by code carries an item title from before O-124, and repeats the parent number.
   - Where: `.review-in/issues/90.md:20`; `.review-in/issues/94.md:28` and `:52` (the same in #95, #96 and #97); `.review-in/issues/97.md:25`.
   - Cited sentence (#90:20): "`setup-s06-facts` — S06: the briefs and the first answers record as raw facts".
   - Basis: By O-124 (`docs/tasks/T-55n2.md:51-55`), S04 writes the record of the `S01-` and `Q-` answers, and S06 keeps the briefs; row 9 (#86) holds that record. The title is the inventory's at `7cdd346`. K14 in the same issue says that row 9 writes the reading into the S04 and S06 rows, so the implementer can find it; one line under the item prevents a wrong read. In #94 to #97, "Related: #29, #29" and "Refs #29, #29" repeat the parent, because the parent is #29 (`gen-issues.py:127-128` and `:150`). In #97:25 the K12 line nests parentheses. The bodies agree with the plan: this review ran `gen-issues.py 77` into a temporary directory, and the 21 files equal the posted texts, except one trailing blank line per file.

7. **note** — The option texts of Q1 to Q4 are paraphrased, not quoted; this review could not verify the paraphrase.
   - Where: `docs/tasks/T-55n2.md:27-29`.
   - Cited sentence: "The options are those of the questions comment as edited after the plan review."
   - Basis: AC 9 asks that the Operator's answers are copied into Git. The answers are quoted ("Q-4 A , Q-3 A, Q-2 A , Q-1 A", `:27`). The text of each option (a) is in Git only as the author's prose under O-121 to O-124; the questions comment is the only home of the option texts. This review is barred from the comments, so it could not compare the prose with the comment.
   - Fix at close-out: quote each chosen option's text, as the comment shows it, under its O- number (Invariant 1).

8. **note** — The parent issues have no "Solution note (R3)" section.
   - Where: `.review-in/issues/77.md`; `.review-in/issues/33.md`; the form is `docs/templates/github/ISSUE_TEMPLATE/task.md:16`.
   - Cited sentence (AC 4): "The phase-1 issues exist, in the task-issue form".
   - Basis: #77 and #33 have Goal, Child tasks, Duplicate check (R2), Acceptance criteria and Notes. They have no Solution note. #29, the existing parent, has none either. A parent issue of R11 selects nothing itself; its children carry the note. So the form holds in substance. One line ("Solution note: the selections are the children's") makes the form complete.

9. **note** — The #42 children table of AC 12 could not be verified.
   - Where: `.review-in/issue-76-body.md:52`; `docs/tasks/T-55n2.md:58-59`.
   - Cited sentence: "Docs updated in the same PR (backlog lines of the new tasks; the #42 children table)."
   - Basis: The backlog lines are in the diff (`docs/tasks/backlog.md:34-57`). The #42 children table is on the forge and is not in `.review-in/`; this review did not use `gh`. `docs/tasks/T-55n2.md:58-59` says that O-125 and O-126 were copied to #76 and #42. This is a limit of this review, not a finding against the change.

**Checks of the brief with no finding.** (2) The plan is derived from the sources: each item of the task table names a section of `docs/spec/`: `gate.md` for row 5; `setup.md` "The steps", S01 to S15 complete, for rows 9, 13 and 15; `setup.md` "The checks of `layup setup verify`", all fifteen checks, for rows 7, 10, 11, 12 and 15; `setup.md` "The stack catalog" for rows 4 and 14; `records.md` for rows 17 and 18; `packages.md` for row 2; `README.md` "Commands" and "The schema block" for rows 1 and 3; `psb-check.md` for row 6. The milestones `M2a` to `M4c` follow the phases 2 to 4 of `PRD-0001` §9 and the sections §5, §6, §8, §10, §11 and §12 of the architecture. The order can be built: the "After" graph has no cycle, and the edges of the inventory are kept (checks 3, 9 and 10 of `plan-check.py`; the layers at `inventory.md:2199-2224`). (3) The §12 cells, the traceability rows and the glossary rows say true things, except findings 3, 4 and 5: the three `green` tests exist under the names and files given (`internal/psb/check_test.go:29`, `internal/cli/cli_test.go:48`, `internal/psb/check_integration_test.go:12`) and pass; the four tests that the traceability preamble names exist; `F-0004#11` and `PRD-0001` §8 say "The first pilot measures the baseline", as the glossary row says. (4) The issue texts agree with the plan: row, parent, items, requirements, predecessors, the defects settled and read, and the start rule (#78 and #79 start now; the others wait for #42), verified by a fresh run of `gen-issues.py` and by reading #29 and #33 against the parent table (`docs/plan/README.md:75-79`). (5) The plan takes no decision that O-121 to O-126 leave to the Operator: K25 and the option A of row 1 are marked as the plan's choice, which each task's plan review can change; the S05 line rule of row 13 is "decided here"; the learning loop waits for the Operator (`:212`). The diff is 4,071 lines added and 36 removed over 20 files, inside the budget of 4,800 lines over 36 files.

### Acceptance criteria

- AC 1: **met**. `docs/plan/README.md:47-65`: `M1` in full (the task table), `M2a` to `M4c` with their requirements, first demo, predecessors and "its specification task".
- AC 2: **not met**. No demo sentence per task; see finding 1.
- AC 3: **met**. `python3 runs/T-55n2/plan-check.py`, run by this review at `a88f99f`: twelve `PASS` lines, exit 0; the red and the green runs are in `plan-check.md`; checks 1 to 5 cover the four conditions of the criterion.
- AC 4: **met**. #77 to #97 are in the task-issue form; #78 and #79 "start now" (O-126), the others are "on hold until #42 closes"; #29's children table lists the parents in row order with #78 to #97; #33 is rewritten to ADR-0016. See finding 8 for the parents' form.
- AC 5: **met**. Each of the 25 Task cells names a task or a milestone; the §13 row is at `:252`; `prd-lint: OK` (run by this review).
- AC 6: **met**. The 25 rows of `docs/tests/traceability.md` cover REQ-001 to REQ-018 and NFR-001 to NFR-007 (check 7 of `plan-check.py`; counted by hand).
- AC 7: **met**. Six rows at `docs/glossary.md:144-148`; "First pilot" tells the phase-4 pilot, `PRD-0001` §8 and `F-0004#11` apart in its collision line.
- AC 8: **met**. Eight dispositions at `docs/plan/README.md:271-278` for #15, #21, #24, #34, #48, #49, #61 and #68.
- AC 9: **met** for the answers; the option texts are paraphrased, see finding 7.
- AC 10: **met**. This review finds the plan derived from `PRD-0001`, `docs/spec/` and the architecture; see "Checks of the brief with no finding", item (2).
- AC 11: **met**. Run by this review at `a88f99f`: `adr-lint`, `prd-lint`, `link-lint` (1231 links), `run-discipline-tests` (81 passed), `nested-checkout-check` (10 cases), `setup-check` and `git diff --check 7cdd346 HEAD`, each exit 0; `go build ./...`, `go vet ./...` and `go test -count=1 ./...` with no tag, `-tags=integration` and `-tags=e2e` pass. `runs/T-55n2/test-runs.md` agrees.
- AC 12: **met** for the backlog lines; the #42 children table was not verified, see finding 9.
