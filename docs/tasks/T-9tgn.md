# T-9tgn — ADR-0012 (Proposed): successor of T-w79d after the issue split

Issue: [#53](https://github.com/pharzam/layup/issues/53). Source: [#46](https://github.com/pharzam/layup/issues/46), task [`T-w79d`](T-w79d.md), closed as split after its last round (`not mergeable, findings recorded`) at `9c989ef`. This task carries that work on branch `T-9tgn` and fixes the four findings of #46's round 3. Record: [ADR-0012](../adr/0012-route-work-across-harness-agents.md). Evidence: `runs/T-w79d/` — the directory keeps the name of the task that produced the check, the clause map and the inventory; this task adds to them.

## Operator decisions (the Operator's words, quoted; given in the session on 2026-09-24)

- **O-26 and O-27**, one answer: "Successor takes all 4 Close #46 as split and link it to the successor” states the required action directly ALSO , For this session, do not select Fabel 5.1 or GPT 6 Astra , AGY for any task, review, panel, or fallback." (The Operator's line break is written here as one space.) The author's reading: O-26, this task takes all four findings and #46 closes as split; O-27, no run of this session uses Claude Fable 5.1, GPT-6 Astra or the AGY harness. O-27 changes no text of ADR-0012.
- **O-28**: the Operator chose "Extend O-15 to #53 (Recommended)" — "Record a new decision (O-28): O-15 also applies to #53, the successor that carries #46's goal. …" So this task counts one goal, the ADR-0012 record.
- **O-29**: the Operator chose "Kimi K3 first, Sonnet 5 backup" — "Kimi K3 on OpenCode (standalone), a different provider from mine. If a run gives no answer in 20 minutes, Claude Sonnet 5 on Claude Code runs it instead, and the record names the tier mismatch." (Devin had refused every model: "Upgrade to Pro to access this model".)

## Which binding ran each role (O-29 names three runs)

| Role | First binding | Result | Ran on |
| -- | -- | -- | -- |
| Plan review (revision 2) | Kimi K3 on OpenCode, standalone | no output in 1,200 s (exit 124) | Claude Sonnet 5 on Claude Code, 603 s — an execution-tier model on a reasoning-tier part |

The plan review of revision 1 ran on GPT-6 Sol (xhigh) on Devin, before Devin refused it. The review rounds are recorded on #53 and in the resource record at close-out.

## Selection: what the route does when no binding fits (finding 3)

| Option | For | Against | Result |
| -- | -- | -- | -- |
| Stop and ask the Operator (#46's D5) | no reasoning step runs on a lighter model | contradicts ADR-0005 ("run the work on the tier you have, and name the tier you could not reach"); halts work a recorded limit allows | rejected |
| Use a binding of the other tier that D7, D6 and D3 allow, and record the limit | matches ADR-0005; the limit is visible in the resource record (ADR-0007) | a reasoning step can run on an execution-tier model | selected (D5) |
| Wait until a binding of the step's tier is eligible | no mismatch | no bound on the wait; `F-0001#14` has no such pause | rejected |

## Affected documents (step 5)

`git grep -n -E "empty output|Fable 5\.1|Astra|T-w79d" -- . ':!runs/' ':!docs/facts/'`, each hit classed:

- **Stale, fixed:** the guardrails §2 lesson said "empty output" where ADR-0012 D4 now says standard output — aligned; the backlog line named `T-w79d` — now `T-9tgn`.
- **Accurate history, kept:** ADR-0012's Context and clause rows (C07, C09, C13, O-21) and its D3 example describe `F-0005` and the inventory; O-3 in `engineering-discipline.md` and `setup/record-T-n1hp.md`; the resource records of earlier tasks; `T-w79d.md`, which is #46's record; the glossary rows; the `runs/T-w79d/` paths. Raw facts are not edited.
