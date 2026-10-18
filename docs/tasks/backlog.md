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
- **T-sdtb** — Row 59b of `M2d`: the demo of M2d (uat) ([#284](https://github.com/pharzam/layup/issues/284); [the plan](../plan/README.md#the-tasks-of-m2d))
- **T-q0kr** — Row 66f of `M2e`: The four statuses ([#362](https://github.com/pharzam/layup/issues/362); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-xvvd** — Row 66g of `M2e`: The pull request and the merge ([#363](https://github.com/pharzam/layup/issues/363); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-3f94** — Row 66h of `M2e`: The step shape ([#364](https://github.com/pharzam/layup/issues/364); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-8rm6** — Row 66i of `M2e`: The approvers and the revision of the draft ([#365](https://github.com/pharzam/layup/issues/365); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-2vva** — Row 66j of `M2e`: The step bet ([#366](https://github.com/pharzam/layup/issues/366); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-9nbf** — Row 66k of `M2e`: The step loop and the restart ([#367](https://github.com/pharzam/layup/issues/367); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-gyvo** — Row 66o of `M2e`: The question and answer comments ([#405](https://github.com/pharzam/layup/issues/405); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-4pje** — Row 66l of `M2e`: The item computers of the task loop ([#396](https://github.com/pharzam/layup/issues/396); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-3xvc** — Row 67 of `M2e`: The e2e of the specification task ([#368](https://github.com/pharzam/layup/issues/368); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-ym0k** — Row 68a of `M2e`: The review of the release of M2e ([#369](https://github.com/pharzam/layup/issues/369); [the plan](../plan/README.md#the-tasks-of-m2e))
- **T-8s13** — Row 68b of `M2e`: The demo of M2e ([#370](https://github.com/pharzam/layup/issues/370); [the plan](../plan/README.md#the-tasks-of-m2e))


## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
