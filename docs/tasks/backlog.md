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
- **T-hpid** — Row 41b of `M2c`: the command layup spec check --facts ([#214](https://github.com/pharzam/layup/issues/214); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-1u2v** — Row 43b of `M2c`: the review of meaning ([#218](https://github.com/pharzam/layup/issues/218); [the plan](../plan/README.md#the-tasks-of-m2c); [detail](T-1u2v.md))
- **T-hiv9** — Row 44 of `M2c`: the steps numbering and spans ([#219](https://github.com/pharzam/layup/issues/219); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-2et0** — Row 45a of `M2c`: the draft ([#220](https://github.com/pharzam/layup/issues/220); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-0mey** — Row 45b of `M2c`: the completeness review ([#221](https://github.com/pharzam/layup/issues/221); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-8r3u** — Row 46 of `M2c`: the batch ([#222](https://github.com/pharzam/layup/issues/222); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-hxb9** — Row 47a of `M2c`: the answers ([#223](https://github.com/pharzam/layup/issues/223); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-vwu3** — Row 47b of `M2c`: the review of the answers and the follow-up ([#224](https://github.com/pharzam/layup/issues/224); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-tmm2** — Row 47c of `M2c`: the answers to the follow-up ([#225](https://github.com/pharzam/layup/issues/225); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-mo9w** — Row 48 of `M2c`: the e2e of the steps of Intake ([#226](https://github.com/pharzam/layup/issues/226); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-8byd** — Row 49a of `M2c`: the review of the release of M2c ([#227](https://github.com/pharzam/layup/issues/227); [the plan](../plan/README.md#the-tasks-of-m2c))
- **T-ekh5** — Row 49b of `M2c`: the demo of M2c (uat) ([#228](https://github.com/pharzam/layup/issues/228); [the plan](../plan/README.md#the-tasks-of-m2c))


## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
