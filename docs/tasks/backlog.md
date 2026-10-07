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

- **T-esfe** — Row 22 of the plan: the forge interface, the two host registers, and the package rules of `M2a` ([#127](https://github.com/pharzam/layup/issues/127); [plan](../plan/README.md#the-tasks-of-m2a))
- **T-xhgz** — Row 23 of the plan: the calls `Fetch` and `Push` of `internal/git`, with the token of one call ([#128](https://github.com/pharzam/layup/issues/128); [plan](../plan/README.md#the-tasks-of-m2a))
- **T-6bq5** — Row 24 of the plan: the GitHub adapter ([#129](https://github.com/pharzam/layup/issues/129); [plan](../plan/README.md#the-tasks-of-m2a))
- **T-trej** — Row 25 of the plan: `internal/run`, Start and the restart ([#130](https://github.com/pharzam/layup/issues/130); [plan](../plan/README.md#the-tasks-of-m2a))
- **T-mqty** — Row 26 of the plan: the command `layup run` ([#131](https://github.com/pharzam/layup/issues/131); [plan](../plan/README.md#the-tasks-of-m2a))
- **T-fnsr** — Row 27 of the plan: the demo of `M2a` (uat), with the Operator's inputs ([#132](https://github.com/pharzam/layup/issues/132); [plan](../plan/README.md#the-tasks-of-m2a))


## Next

<!-- Deliberately deferred tasks, same one-line shape. -->
