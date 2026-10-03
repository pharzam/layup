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
- **<ID>** — <one-sentence summary> ([<ADR or doc link>](...); [detail](<id>.md))
-->

- **T-stfn** — Step 2 (the name that the title of its issue #29 gives it), the core engine: `layup gate`, `layup setup` and `layup setup verify`, the telemetry and stall records, the release review and the first pilot ([#29](https://github.com/pharzam/layup/issues/29); [plan](../plan/README.md))
- **T-b97r** — `layup setup` and `layup setup verify`, the parent of its child tasks ([#77](https://github.com/pharzam/layup/issues/77); [plan](../plan/README.md))
- **T-vk3k** — `layup gate`, the parent of its child tasks ([#33](https://github.com/pharzam/layup/issues/33); [plan](../plan/README.md))
- **T-d6q5** — Row 15 of the plan: Steps S12, S13 and S15, and the checks `jobs` and `gate:<kind>` ([#92](https://github.com/pharzam/layup/issues/92); [plan](../plan/README.md#the-tasks-of-phase-1))
- **T-dep6** — Row 16 of the plan: A whole setup, end to end, with no network ([#93](https://github.com/pharzam/layup/issues/93); [plan](../plan/README.md#the-tasks-of-phase-1))
- **T-tmhw** — Row 17 of the plan: The telemetry record: the schemas of `telemetry.tsv` and `prices.tsv` in code ([#94](https://github.com/pharzam/layup/issues/94); [plan](../plan/README.md#the-tasks-of-phase-1))
- **T-dgy7** — Row 18 of the plan: The stall record: the schema of `stalls.tsv` in code ([#95](https://github.com/pharzam/layup/issues/95); [plan](../plan/README.md#the-tasks-of-phase-1))
- **T-efmy** — Row 19 of the plan: The release review of phase 1 for the `Won't` rows ([#96](https://github.com/pharzam/layup/issues/96); [plan](../plan/README.md#the-tasks-of-phase-1))
- **T-evad** — Row 20 of the plan: The first pilot (O-122): set up one target from a Go problem statement, and run its gate from outside ([#97](https://github.com/pharzam/layup/issues/97); [plan](../plan/README.md#the-tasks-of-phase-1))

## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
