# Backlog

Lean task list for this project. Two sections: **Now** (the current version) and
**Next** (deliberately deferred).

## How to keep this file readable

**One line per task — keep it that way.** Each entry is exactly an ID, a
one-sentence summary, and the link(s) that motivate or scope it. Nothing more. If
a task needs more than that — design notes, rejected alternatives, open questions,
reproduction detail — it goes in `tasks/<id>.md` and the entry links to it as
`[detail](<id>.md)`. Do **not** grow the entry itself; this file is an index, not
a design doc. (Multi-paragraph entries are prohibited.)

Each task has a stable ID assigned once and never reused or renumbered — an ID
stays with its task when promoted from Next to Now. Use a task-ID scheme (`T-` plus four random characters from `0-9 a-z` without `i l o u`): a
short, stable token per task. Prefer **random** IDs over a sequential counter — a
counter forces every session to agree on "the next number", so two people (or
agents) working in parallel both pick the same one and collide in filenames,
branches, and PRs. A random suffix needs no coordination. This project's scheme: `T-` plus four random characters from `0-9 a-z` without
`i l o u`; before using an ID, confirm `tasks/<id>.md` does not already exist
(Operator decision O-8 on [#8](https://github.com/pharzam/layup/issues/8)).

When a Now item is done, move its line to [completed.md](completed.md) — same ID,
same summary, dated — rather than deleting it or checking it off.

## Now

<!-- One line per task. Example shape:
- **‹ID›** — ‹one-sentence summary› ([‹ADR or doc link›](...); [detail](‹id›.md))
-->

- **T-meh2** — The PDR: `PRD-0001` with one requirement per PSB In-Scope item, traced to `F-0003#41`–`#52`, then the architecture and the plan ([#42](https://github.com/pharzam/layup/issues/42))
- **T-wjq4** — `PRD-0001`, the product requirements of LAYUP: full scope with four phases, an acceptance criterion per requirement, the traceability matrix ([#63](https://github.com/pharzam/layup/issues/63); child 3 of #42)
- **T-84r5** — `PRD-0001`, the successor of `T-wjq4` after its cycle cap: the eight criteria of round 3 tightened, one round, close-out ([#64](https://github.com/pharzam/layup/issues/64); child 3 of #42)
- **T-stfn** — Step 2, the core engine: `layup setup`, `layup gate`, the telemetry record and the stall record ([#29](https://github.com/pharzam/layup/issues/29); [ADR-0011](../adr/0011-structure-the-core-engine-as-a-go-cli-over-repository-files.md))

## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
