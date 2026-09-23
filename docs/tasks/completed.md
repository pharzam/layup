# Completed

Append-only log of finished tasks, backlog-listed or not, most recent first.

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

- **2026-09-23** — **T-nfh8** — Filled each adopter marker that has a source (Operator batch O-1 to O-8, the Go command reference), listed the rest as open gaps, added ADR-0010 (Go) and check `markers` ([#8](https://github.com/pharzam/layup/issues/8); [detail](T-nfh8.md))
- **2026-09-23** — **T-7ndb** — Encoded the 9 System Invariants into `docs/guardrails.md` §1.1, wrote back three lessons to §2, and added check `guardrails` ([#7](https://github.com/pharzam/layup/issues/7); [detail](T-7ndb.md))
- **2026-09-23** — **T-xgz4** — Merged the 25 PSB §8 terms into `docs/glossary.md` §1, each citing its `F-0001` fact, and added check `glossary` ([#6](https://github.com/pharzam/layup/issues/6); [detail](T-xgz4.md))
- **2026-09-23** — **T-vpty** — Bound `docs/onboarding-for-engineers.md` to the PSB (`F-0001`) and added check `onboarding` ([#5](https://github.com/pharzam/layup/issues/5); [detail](T-vpty.md))
- **2026-09-23** — **T-fvwj** — Stored the PSB (`F-0001`, 39 numbered verbatim facts) and the vision brief (`F-0002`, a solution document) byte-identical, with hashes and check `facts` ([#4](https://github.com/pharzam/layup/issues/4); [detail](T-fvwj.md))
- **2026-09-23** — **T-vbwc** — Removed the kit's own history (kit step 4: `docs/decisions/`, `docs/audit/`, 19 kit task files, kit log and backlog lines) and added check `kit-history` ([#3](https://github.com/pharzam/layup/issues/3); [detail](T-vbwc.md))
- **2026-09-23** — **T-r7zg** — Recorded the Armature pin (`a959655`, tree `8ffb250a`) and started `docs/setup/setup-check.sh` with check `pin` and the kit-linter pass-through ([#2](https://github.com/pharzam/layup/issues/2); [ADR-0009](../adr/0009-pin-armature-at-a-recorded-commit.md); [detail](T-r7zg.md))
