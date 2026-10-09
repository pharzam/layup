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

- **T-cht1** — Row 30b of `M2b`: admission and the routing order ([#156](https://github.com/pharzam/layup/issues/156); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-vxdg** — Row 33a of `M2b`: the directory and the environment of a session ([#159](https://github.com/pharzam/layup/issues/159); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-6sbe** — Row 33b of `M2b`: the checks before the start of a session ([#160](https://github.com/pharzam/layup/issues/160); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-5pxd** — Row 34a of `M2b`: the process of a session and its stop ([#161](https://github.com/pharzam/layup/issues/161); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-bpxg** — Row 34b of `M2b`: the end of a session ([#162](https://github.com/pharzam/layup/issues/162); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-d8t9** — Row 36a of `M2b`: the start of a task session and its refusals ([#164](https://github.com/pharzam/layup/issues/164); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-fsjp** — Row 36b of `M2b`: the result and the end of a task session ([#165](https://github.com/pharzam/layup/issues/165); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-z027** — Row 37a of `M2b`: the checks before a push ([#166](https://github.com/pharzam/layup/issues/166); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-e3sy** — Row 37b of `M2b`: the push and the bind ([#167](https://github.com/pharzam/layup/issues/167); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-nxe4** — Row 38 of `M2b`: the step probe ([#168](https://github.com/pharzam/layup/issues/168); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-4tjy** — Row 39a of `M2b`: the review of the release of M2b ([#169](https://github.com/pharzam/layup/issues/169); [the plan](../plan/README.md#the-tasks-of-m2b))
- **T-x7cs** — Row 39b of `M2b`: the demo of M2b (uat) ([#170](https://github.com/pharzam/layup/issues/170); [the plan](../plan/README.md#the-tasks-of-m2b))

## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
