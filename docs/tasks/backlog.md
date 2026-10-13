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
- **T-jpkc** — Row 51a of `M2d`: the forge calls of M2d ([#271](https://github.com/pharzam/layup/issues/271); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-j4sd** — Row 51b of `M2d`: the call CommitTree ([#272](https://github.com/pharzam/layup/issues/272); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-56dj** — Row 52c of `M2d`: the commands of S13 ([#275](https://github.com/pharzam/layup/issues/275); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-eahr** — Row 54 of `M2d`: the step setup ([#277](https://github.com/pharzam/layup/issues/277); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-79zk** — Row 55 of `M2d`: the step verify ([#278](https://github.com/pharzam/layup/issues/278); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-hfn1** — Row 56 of `M2d`: the step scaffold-push ([#279](https://github.com/pharzam/layup/issues/279); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-wphp** — Row 57a of `M2d`: the step forge-check ([#280](https://github.com/pharzam/layup/issues/280); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-2j13** — Row 57b of `M2d`: the forge check at each later start ([#281](https://github.com/pharzam/layup/issues/281); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-eevt** — Row 58 of `M2d`: the e2e of the Scaffold ([#282](https://github.com/pharzam/layup/issues/282); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-5xjc** — Row 59a of `M2d`: the review of the release of M2d ([#283](https://github.com/pharzam/layup/issues/283); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-sdtb** — Row 59b of `M2d`: the demo of M2d (uat) ([#284](https://github.com/pharzam/layup/issues/284); [the plan](../plan/README.md#the-tasks-of-m2d))


## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
