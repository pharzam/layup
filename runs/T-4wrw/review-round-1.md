## Review record — round 1

| Field | Value |
| ----- | ----- |
| Commit reviewed | `9655eb16af657d670a4fc6418406b86af42715a0` |
| Reviewer | Claude Fable 5.1 on Claude Code |
| Lens | correctness and acceptance criteria |
| Briefed on | `.review-in/brief.md`; `.review-in/issue-body.md` (#99, AC 1 to AC 7); `.review-in/issue-comments.md` (the plan, the plan review, the author's answer, the request for the approval, the note on rows 1 and 2, the Operator's comment); `.review-in/42.md`; `.review-in/forge.md`; the whole diff `git diff 3fb44d3 HEAD` (11 files); `git diff b48764f HEAD -- docs/prd/PRD-0001-layup.md`; `git diff b48764f 3fb44d3` for the approved documents; `docs/engineering-discipline.md` (Bootstrap mode; What a round records); `docs/pdr/PDR-0001.md`; `docs/tasks/T-4wrw.md`; `runs/T-4wrw/trace-check.py` and `trace-check.md`, at HEAD and at `0aca6a6`; `docs/prd/PRD-0001-layup.md` §1, §12, §13; `docs/prd/README.md`; `docs/setup/open-gaps.tsv`; `docs/tests/dod-checklist.md`; `docs/tasks/backlog.md`, `completed.md`, `T-18v6.md`, `T-2tc2.md`, `T-55n2.md` (O-122), `T-hbw8.md` (O-112); `docs/glossary.md` (the row "Preliminary Design Review"); `docs/onboarding-for-engineers.md` (the `prd/` and `pdr/` rows; "Where the project stands"); `docs/facts/F-0001` (#11) and `F-0004` (#1); `git log` of `main` from `b48764f` to `3fb44d3`. Not read: the forge itself (no network); the state of the issues and the approval comment come from `forge.md`. |
| Barred from | nothing; the comments of #99 are the plan and its record |
| Independence claimed | Context: a fresh read-only session in a clone at `9655eb1`, with no network; the author text read is the plan, its answer, the request, the note on rows 1 and 2 and the diff, as the brief names; no earlier review round exists for this head. Execution: this record. Model: Claude Fable 5.1; the author is Claude Opus 5.5 on every commit of the branch. Not reached: the plan reviewer was also Claude Fable 5.1, in a different session; method independence does not apply to a first round. |
| Cycle | 0 |
| Verdict | `nothing material in scope` |

### Raw findings

1. **Note.** `runs/T-4wrw/trace-check.md:17` and `:36`, the detail line `docs/pdr/PDR-0001.md is missing` of Red 1 and Red 2. The script prints the absolute path of the record (`trace-check.py:95`, `f"{REC} is missing"`, where `REC` is `ROOT / "docs/pdr/PDR-0001.md"` and `ROOT` is absolute). The script did not change after the red commit `0aca6a6` (`git diff 0aca6a6 HEAD -- runs/T-4wrw/trace-check.py` is empty), so the recorded line is a shortened copy of the real output. I reproduced Red 1 and Red 2 on copies under `/tmp/rv-T-4wrw/`: each FAIL line, each other detail line and the exit code agree with the record; only this path differs. Classification: `note` (the report is true in its result and its reason; one detail line is shortened). The author can say in the file that the path is shortened, or paste the line as printed.

2. **Note.** `docs/pdr/PDR-0001.md:74` and `docs/tasks/T-4wrw.md:48`, the quote `The PDR (#99): approved`. The body on the forge (`forge.md`, the approval comment) is ` The PDR (#99): approved`, with one leading space. The words are the same. Classification: `note`; no change is needed.

Checks with no finding.
- **The decision number.** O-129 is the next free number. In the tree, `git grep` finds O-127 in `docs/tasks/T-18v6.md:59`, O-128 in `docs/tasks/T-2tc2.md:64`, and O-129 only in the files of this task (`PDR-0001.md:7`, `:76`; `PRD-0001-layup.md:254`; `T-4wrw.md:31`, `:46`; `trace-check.md:48`); no O-130 or higher exists. On the forge (`forge.md`), the highest numbers outside this task are O-127 and O-128 (#78, #79), and O-127 on #99 is only the plan's placeholder (D4) and the plan review's note 8.
- **The approval copy.** `PDR-0001.md:70-74` and `T-4wrw.md:46-48` agree with `forge.md`: the link `issues/99#issuecomment-5946822016`, the time `2026-10-02T06:37:30Z`, the account `pharzam` (`author_association` OWNER, `user_type` User, `performed_via_github_app` null), and the words.
- **`PRD-0001` and the claim of `PDR-0001.md:35-37`.** `git diff b48764f HEAD -- docs/prd/PRD-0001-layup.md` changes exactly the Status line (`Draft` to `Accepted`), the §12 Test cells of `NFR-005` and `NFR-007`, and two §13 rows (`T-2tc2` and `T-4wrw`). The §13 row (`:254`) names `T-4wrw`, #99, O-129, `b48764f` and the record. `prd-lint` exits 0. The index row (`docs/prd/README.md:64`) and the onboarding `prd/` row (`:136`) say `Accepted`.
- **The commits.** `b48764f` is `Merge pull request #98` (`PDR-0001.md:21-22`); `git merge-base --is-ancestor b48764f origin/main` succeeds; `git cat-file -e` finds each of the four paths in it. The only first-parent commits of `main` after `b48764f` up to `3fb44d3` are the merges of #100 and #101 (`:31`). `git diff b48764f 3fb44d3` for the approved documents touches `docs/architecture.md` (the new limit L-A7 in §15), `docs/spec/` (`README.md`, `packages.md` with "The calls of `internal/git`" and "The test of the package rules", `psb-check.md`, `setup.md`) and the two §12 Test cells of `PRD-0001`, as `:32-34` says. No other approved document changes between `3fb44d3` and HEAD.
- **The counts.** The check reports 25 rows and 99 distinct tokens (`:41-42`). `open-gaps.tsv` has 18 rows: 16 for `docs/tests/dod-checklist.md`, one for `docs/issue-workflow.md` (O-7) and one for `docs/engineering-discipline.md` (`:57-63`). `NFR-007` cites `F-0004#1` only (`PRD-0001-layup.md:243`); `F-0001#11` is "Problem statement answers" (`:46-48`). The pilot's inputs (`:64-66`) are the words of O-122 (`T-55n2.md:47-49`). ADR-0013 to ADR-0025 were accepted at `b2e0ed4` under O-112 (`T-hbw8.md:131`).
- **The children of #42** (`T-4wrw.md:38-42`) agree with `forge.md`: #43 by PR #44, #45 by PR #62, #63 closed with no PR and #64 by PR #65, #66 `NOT_PLANNED` and #72 by PR #73, #74 by PR #75, #76 by PR #98; the body of #29 lists the plan's issues.
- **The plan review's conditions and notes**, as the author's answer states them: condition 1 (`PDR-0001.md:85-87`; the request comment says that each change becomes a new task); condition 3 (Red 1 holds the four seeded defects, one FAIL line per part); notes 1 (`:21-22`, `:35-37`), 2 (parts 2 and 3 and the `info` line of the script), 3 (`:46-48`), 4 (the docstring; `git merge-base` and `git cat-file` in part 4; the link form and a non-empty quote in part 5), 5 (the record at `0aca6a6` says `Pending: the Operator's comment on #99`, no marker), 6 (check `adapted` exits 0), 7 (`:53-66`), 8 (O-129), 9 (one file, three runs), 10 (`PRD-0001-layup.md:254`), 11 (`README.md:60`, `AGENTS.md:186`, the onboarding `pdr/` row and the sentence on #42), 13 (`T-4wrw.md:38-42`). Condition 2 (the `T-meh2` line) and the completed log are close-out, not in the head.
- **The budget.** `git diff --shortstat 3fb44d3 HEAD`: 11 files, 360 insertions, 6 deletions, so 366 lines over 11 files, against the maximum of 800 lines over 16 files. (`git diff --shortstat b48764f HEAD` gives 46 files and 4229 lines, but that figure holds the merged rows 1 and 2, which are not this task's change.)
- **The onboarding sentence** "closed by the PDR" (`docs/onboarding-for-engineers.md:164`) and the glossary example (`docs/glossary.md:147`) become true when this pull request merges and #42 closes; the whole change is read at the merge, so this is not a finding.

### Acceptance criteria

- AC 1: **met.** `docs/pdr/PDR-0001.md:24-29` names the four documents by path and by `b48764f`, a commit of `main` that holds each path.
- AC 2: **met.** `runs/T-4wrw/trace-check.py` runs from the clone root with 6 PASS lines and exit 0; its output is `runs/T-4wrw/trace-check.md`; the seeded copies fail (see Runs).
- AC 3: **met.** The approval is comment 5946822016 on #99 by `pharzam`, the Operator's own account with no App; `PDR-0001.md:70-74` copies it with its link, its date and time, and its words.
- AC 4: **met.** Status `Accepted` (`PRD-0001-layup.md:11`), the §13 row (`:254`), `prd-lint` exit 0, the index row (`docs/prd/README.md:64`) and the onboarding sentence (`docs/onboarding-for-engineers.md:136`).
- AC 5: **at close-out.** The head holds the check of #42's children against the forge (`T-4wrw.md:38-42`), which `forge.md` confirms; the ticked boxes and the close come with the pull request.
- AC 6: **met.** The seven local checks exit 0 (see Runs).
- AC 7: **met for the head; the completed log at close-out.** The glossary row names `pdr/PDR-0001.md` (`docs/glossary.md:147`); the backlog has the `T-4wrw` line (`docs/tasks/backlog.md:34`); the completed-log lines of `T-4wrw` and `T-meh2` are close-out.

### Runs

Each command ran from the clone root with `HOME=/tmp/rv-T-4wrw/home`; the copies are under `/tmp/rv-T-4wrw/`.

| Command | Result |
| ------- | ------ |
| `python3 runs/T-4wrw/trace-check.py` | `info  numbered facts: F-0001 39, F-0002 0, F-0003 75, F-0004 19`; 6 PASS lines; exit 0 |
| `python3 runs/T-4wrw/trace-check.py --record rec-pending.md` (the approval part reads `Pending: the Operator's comment on #99.`) | `FAIL  the record holds the approval` with `no link to the approval comment on #99`, `no quote of the approval`, `the approval is still pending`; 5 PASS; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --record rec-nolink.md` (the link removed) | `FAIL  the record holds the approval` with `no link to the approval comment on #99`; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --record rec-noquote.md` (the `> ` quote line removed) | `FAIL  the record holds the approval` with `no quote of the approval`; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --record rec-badsha2.md` (the `PRD-0001` row names `9655eb1`) | `FAIL  the record names each document by a commit of main that holds it` with `PRD-0001: 9655eb1 is not a commit of origin/main`; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --record rec-badsha3.md` (the specification row names `c0ffee1`) | `FAIL  the record names each document by a commit of main that holds it` with `the specification: c0ffee1 is not a commit of origin/main`; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --prd prd-draft.md` (a copy with the Status `Draft`) | `FAIL  PRD-0001 has the Status `Accepted`` with `Status is Draft`; exit 1 |
| `python3 runs/T-4wrw/trace-check.py --prd PRD-seeded.md --record none.md` (Red 1 reproduced: the Facts cell of `REQ-001` empty; `F-0003#76`, `F-0002#1`, `F-0099#1` added to `REQ-002`) | the six FAIL lines and the detail lines of `trace-check.md` Red 1, `102 distinct tokens`; the path line reads `/tmp/rv-T-4wrw/none.md is missing` (note 1); exit 1 |
| `python3 runs/T-4wrw/trace-check.py --prd prd-b48764f.md --record none.md` (Red 2 reproduced: `PRD-0001` as at `b48764f`, no record) | 3 PASS, then the three FAIL lines of `trace-check.md` Red 2 (`Status is Draft`); exit 1 |
| `sh docs/adr/adr-lint.sh` | `adr-lint: OK`; exit 0 |
| `sh docs/prd/prd-lint.sh` | `prd-lint: OK`; exit 0 |
| `sh docs/links/link-lint.sh` | `link-lint: OK  1284 links resolved`; exit 0 |
| `sh docs/tests/run-discipline-tests.sh` | `81 passed, 0 failed`; exit 0 |
| `sh docs/tests/nested-checkout-check.sh` | `nested-checkout-check: OK  10 cases behaved`; exit 0 |
| `sh docs/setup/setup-check.sh` | each check OK, `setup-check: kit-linters OK`; exit 0 |
| `git diff --check 3fb44d3 HEAD` | no output; exit 0 |
| `git diff --shortstat 3fb44d3 HEAD` | `11 files changed, 360 insertions(+), 6 deletions(-)`: 366 lines over 11 files, inside 800 over 16 |
| `git merge-base --is-ancestor b48764f origin/main`; `git cat-file -e b48764f:<path>` for the four paths | each exit 0 |
| `git grep -n -E 'O-12[5-9]\|O-13[0-9]' HEAD` | O-127 in `T-18v6.md`, O-128 in `T-2tc2.md`, O-129 only in this task's files, no O-130 or higher |
| `git diff 0aca6a6 HEAD -- runs/T-4wrw/trace-check.py` | empty: the script did not change after the red runs |
| `git status --porcelain` after the checks | clean |

*Posted for a reviewer session (Claude Fable 5.1 on Claude Code, `claude -p`, a fresh read-only session in a clone at `9655eb1`, from 09:54 +03, 6 min 25 s) through the `layup-agent` App. The text is the session's file, unchanged.*
