# Completed

Append-only log of finished tasks, backlog-listed or not, most recent first.

- **2026-09-26** — **T-0kn4** — Landed the software architecture (`docs/architecture.md`) and ADR-0013 to ADR-0017 (the LAYUP App's check runs, the rule guard, the kit's roles and the handoff row, the escalation floor and class, the stall limits) after the panel of `T-7qvc` and the Operator's decisions O-52 to O-65 ([#67](https://github.com/pharzam/layup/issues/67); [detail](T-0kn4.md))
- **2026-09-26** — **T-7qvc** — Ran the panel of three on the five open questions of `PRD-0001` §11, drafted ADR-0013 to ADR-0017 and `docs/architecture.md`, and ended at the cycle cap with `not mergeable, findings recorded`; split to `T-0kn4` (O-63) ([#66](https://github.com/pharzam/layup/issues/66); [detail](T-7qvc.md))
## How to keep this file readable

**One line per task — keep it that way.** When a task in [backlog.md](backlog.md)
is done, move its line here unchanged except for a leading completion date:
`**YYYY-MM-DD** — **<ID>** — <summary> (<links>)`. Any detail worth keeping — the
design, the bugs caught, the verification story — lives in `tasks/<id>.md`, linked
as `[detail](<id>.md)`; it does **not** go inline here. This is a dated index, not
a changelog narrative. The git history and commit messages hold the blow-by-blow.

**A task that never had a backlog line** has no line to move: write its entry
here directly, in the shape above.

## Log

<!-- Most recent first. Example shape:
- **YYYY-MM-DD** — **‹ID›** — ‹one-sentence summary of what the task found or delivered› ([‹link›](...); [detail](‹id›.md))
-->

- **2026-09-25** — **T-84r5** — Landed `PRD-0001`, the product requirements of LAYUP (Draft; 18 REQ, 7 NFR, an acceptance criterion each, four phases, the §12 matrix) as the successor of `T-wjq4`, with the `prd-lint` fix for a fact list of two or more records ([#64](https://github.com/pharzam/layup/issues/64); [detail](T-84r5.md))
- **2026-09-25** — **T-wjq4** — Wrote `PRD-0001` and fixed the PRD linter (test first), and ended at the cycle cap `not mergeable, findings recorded`; split to `T-84r5` ([#63](https://github.com/pharzam/layup/issues/63); [detail](T-wjq4.md))
- **2026-09-25** — **T-zmj6** — Stored the idea owner's answers to the 19 gap questions of the PSB as raw fact `F-0004` (accepted in one batch, Decision Point 2) and made check `facts` hold the answers record to one fact per question; the first product task under bootstrap mode ([#45](https://github.com/pharzam/layup/issues/45); [detail](T-zmj6.md))
- **2026-09-25** — **T-7sbn** — Landed ADR-0012, "Build LAYUP in bootstrap mode" (Accepted; O-41 to O-44), as the successor of `T-8ywj`: the operative `## Bootstrap mode` section with the Operator's materiality test and one home per rule, the revert of PR #50 (`F-0005`), the Operator's revised routing text as evidence, and the product path #45 → #42 → #29 in the backlog ([#59](https://github.com/pharzam/layup/issues/59); [detail](T-7sbn.md))
- **2026-09-25** — **T-8ywj** — Wrote ADR-0012 and the diagnosis of the routing-line stall, reverted `F-0005`, and ended at the cycle cap `not mergeable, findings recorded`; split to `T-7sbn` ([#56](https://github.com/pharzam/layup/issues/56); [detail](T-8ywj.md))
- **2026-09-23** — **T-ertw** — Added `F-0003`, 75 numbered facts for the rest of the PSB, and extended check `facts` to it ([#43](https://github.com/pharzam/layup/issues/43); [detail](T-ertw.md))
- **2026-09-23** — **T-bhsf** — Required the three Go CI jobs on `main` (12 required checks in all) and recorded the read-back ([#35](https://github.com/pharzam/layup/issues/35); [detail](T-bhsf.md))
- **2026-09-23** — **T-dq05** — Added `layup psb check`, a deterministic gap check of a problem statement (rules G1–G5, one TSV batch); on LAYUP's own PSB it asks 19 questions ([#37](https://github.com/pharzam/layup/issues/37); [detail](T-dq05.md))
- **2026-09-23** — **T-mtb9** — Added the first Go code (`layup version`) and turned on the Go gates in CI and the hook; successor of `T-t8qp` (#32) after the issue split ([#38](https://github.com/pharzam/layup/issues/38); [detail](T-mtb9.md))
- **2026-09-23** — **T-edtd** — Decided the core engine's architecture in ADR-0011 (one Go CLI over repository files; LAYUP's machinery never in a target) after a three-member panel and Operator decisions O-9 to O-13 ([#30](https://github.com/pharzam/layup/issues/30); [detail](T-edtd.md))
- **2026-09-23** — **T-9mmm** — Recorded the setup procedure (15 steps, `docs/setup/README.md` and `steps.tsv`) and the baseline summary, replaced the kit adaptation section, and added check `procedure` ([#11](https://github.com/pharzam/layup/issues/11); [detail](T-9mmm.md))
- **2026-09-23** — **T-6rg3** — Made `README.md` and `AGENTS.md` describe LAYUP (an Armature project pinned at `a959655`) and added check `identity` ([#10](https://github.com/pharzam/layup/issues/10); [detail](T-6rg3.md))
- **2026-09-23** — **T-afa5** — Required all nine CI jobs on `main` by branch protection, kept the body in `docs/setup/branch-protection.json`, and added check `protection` ([#12](https://github.com/pharzam/layup/issues/12); [detail](T-afa5.md))
- **2026-09-23** — **T-fvng** — Made `pr-link.yml` and `review-record.yml` restore their lint script from `main`, and added cause `restore` to check `ci` ([#23](https://github.com/pharzam/layup/issues/23); [detail](T-fvng.md))
- **2026-09-23** — **T-q344** — Replaced the kit's workflow headers, added CI job `setup-check` (full history, checks restored from `main`) and check `ci` ([#9](https://github.com/pharzam/layup/issues/9); [detail](T-q344.md))
- **2026-09-23** — **T-nfh8** — Filled each adopter marker that has a source (Operator batch O-1 to O-8, the Go command reference), listed the rest as open gaps, added ADR-0010 (Go) and check `markers` ([#8](https://github.com/pharzam/layup/issues/8); [detail](T-nfh8.md))
- **2026-09-23** — **T-7ndb** — Encoded the 9 System Invariants into `docs/guardrails.md` §1.1, wrote back three lessons to §2, and added check `guardrails` ([#7](https://github.com/pharzam/layup/issues/7); [detail](T-7ndb.md))
- **2026-09-23** — **T-xgz4** — Merged the 25 PSB §8 terms into `docs/glossary.md` §1, each citing its `F-0001` fact, and added check `glossary` ([#6](https://github.com/pharzam/layup/issues/6); [detail](T-xgz4.md))
- **2026-09-23** — **T-vpty** — Bound `docs/onboarding-for-engineers.md` to the PSB (`F-0001`) and added check `onboarding` ([#5](https://github.com/pharzam/layup/issues/5); [detail](T-vpty.md))
- **2026-09-23** — **T-fvwj** — Stored the PSB (`F-0001`, 39 numbered verbatim facts) and the vision brief (`F-0002`, a solution document) byte-identical, with hashes and check `facts` ([#4](https://github.com/pharzam/layup/issues/4); [detail](T-fvwj.md))
- **2026-09-23** — **T-vbwc** — Removed the kit's own history (kit step 4: `docs/decisions/`, `docs/audit/`, 19 kit task files, kit log and backlog lines) and added check `kit-history` ([#3](https://github.com/pharzam/layup/issues/3); [detail](T-vbwc.md))
- **2026-09-23** — **T-r7zg** — Recorded the Armature pin (`a959655`, tree `8ffb250a`) and started `docs/setup/setup-check.sh` with check `pin` and the kit-linter pass-through ([#2](https://github.com/pharzam/layup/issues/2); [ADR-0009](../adr/0009-pin-armature-at-a-recorded-commit.md); [detail](T-r7zg.md))
