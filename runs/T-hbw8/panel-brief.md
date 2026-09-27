# T-hbw8 — the panel brief

Each member got `COMMON.md` and one `DOMAIN-<member>.md`, word for word below.

## COMMON.md

--- begin
# Panel for the architecture of LAYUP (task T-hbw8, issue #72 of pharzam/layup)

You are one member of a panel of three (ADR-0006: `docs/adr/0006-convene-a-panel-to-generate-options.md`). A panel **generates and compares options; it does not vote and does not recommend one option.** Each member has a different domain (your domain is in `.panel-brief/DOMAIN.md`) and a different model. You do not see the other members' outputs, and there is no second round.

## What LAYUP is (the facts; read them)

- The problem statement (PSB): `docs/facts/problem-statement-brief.md`, with the numbered facts `docs/facts/F-0001-*.md` (facts 1–39: the invariants, the terms) and `docs/facts/F-0003-*.md` (facts 40–75: problems, the In-Scope items `F-0003#41`–`#52`, the measures). The idea owner's answers to the PSB gaps: `docs/facts/F-0004-*.md`.
- The vision brief (a solution input, not requirements): `docs/facts/architectural-vision-brief.md` and `F-0002`.
- The requirements: `docs/prd/PRD-0001-layup.md`.
- The decisions that stay: `docs/adr/0010-*.md` (Go), `docs/adr/0011-*.md` (a Go CLI over repository files), `docs/adr/0012-*.md` (bootstrap mode). The quality gate: `docs/engineering-discipline.md`.
- **The Operator's decisions O-66 to O-75 and the target design**, and the **review of the previous architecture** (findings A1–A3, B1–B9, the forge question, table C): `runs/T-hbw8/inputs-from-pr-69.md`. The Operator has accepted the target design (O-67 to O-71). Your options work inside it; where you think it is wrong, say so as a separate point with the fact that shows it.
- **Jev** (TypeSafe's System One model, the "smart-if"): `runs/T-hbw8/jev-sources.md` and the pages in `.panel-brief/jev/`. Use only facts from these pages; do not assume an API, a latency or a price that they do not state.

The previous architecture (PR #69, closed without a merge) is **not** an input. Do not look for it.

## The questions

For each question, give **two or three materially different options**. The map in brackets names the review findings each question must answer.

1. **The orchestrator** `layup run TARGET`: the phase state machine (Intake → Scaffold → Design → Architect → Plan → Implement → Accept → Retrospective → next milestone); how it starts and stops role sessions on a harness and reads their handoff records; where its state lives; how it detects a stall by the clock and a retry loop that never pushes. [B1, B2; PRD §11 question 5; table C 3.5]
2. **The decision component**: Jev at named decision points (escalation selection, question routing by ambiguity kind, the action for a stall, panel synthesis, model and harness routing, over budget inside the band); the record of each call; the fallback to a human when Jev is unavailable or below a threshold (O-71); the rule that a deterministic check decides first (Invariant 6). This decision replaces ADR-0011 decision 8 ("the engine makes no model call") in part: say which text it keeps. [A3, B8; table C 3.1, 3.2]
3. **Squads and routing**: the register of harnesses and models, the routing table, the roles (PRD §11 question 3), and verification of each change by a counterpart harness, not only a different model (Invariant 9). [B5; table C 2.2]
4. **The records in Git**: handoffs, questions and answers, requirement acceptance (Decision Point 3), telemetry per action with the requirement and the price; how a record that starts on the forge (an issue comment, a review) gets into Git. [B3, B6, B7; table C 3.4]
5. **Gates and rule protection**: where the stack gates of a target live and who writes them (the target must pass its gates without LAYUP, Invariant 2; ADR-0011 decisions 3 and 4 and the Operator decisions O-10 and O-11 put LAYUP's code outside the target); the runner that makes a gate block a merge; rule-path changes only at the retrospective (O-69). Keep true ADR-0012 part 6: the pilot "has run that target's gate from outside". [A1, A2, B9]
6. **Escalation and budget**: selecting a business-forking decision **before** the work (Decision Point 4), the tolerance band (O-68), a cost stop before the budget is spent. [B7, B8]
7. **The learning loop and the retrospective**: the reward from the records, how it changes the routing at each retrospective without a change to model weights, and the batch of rule changes (O-69). [table C 2.3, 3.3]
8. **Specification synthesis**: the component that derives the PRD, the specification and the phased plan (the Design, Architect and Plan phases), and the deterministic check of the trace from each requirement to the PSB text. [B4; `F-0003#51`, `F-0003#62`; table C 2.1]

Also: the **forge question** — GitHub is the only forge for the pilot (a planned known limit). Give an option only if you think that limit is wrong.

## For each option, write

- **Mechanism** — what exists (a command, a file with its columns, a process) and what it does, in three to six lines.
- **Rests on** — the PSB facts (`F-0001#n`, `F-0003#n`, `F-0004#n`) and O-decisions it serves.
- **Strains** — an invariant or fact it puts at risk, if any.
- **Falsifiable** — one observation in the pilot that would show the option wrong.
- **Cost** — infrastructure, model calls, human load.

Then, once, at the end: **Questions for the Operator** — only a choice that no fact and no O-decision settles, each with its options.

## Limits

- One pass. **At most 300 lines.** Plain Markdown, ASD-STE100 Simplified Technical English.
- Do not write the character pair used for placeholders (a single left and right angle quotation mark around text): a check fails on it. Use `<...>` if you need a shape.
- Do not change any file in this directory. Write your output only to `.panel-out/panel.md`, and then print it as your final answer.
--- end

## DOMAIN-A.md

--- begin
# Your domain: orchestration, state in Git, and the forge (member A)

You look at the questions as a systems engineer for a Go command-line program over repository files. Weigh most: process lifecycle and crash recovery of `layup run`; state that lives in Git and survives a restart; how a harness session is started, watched and stopped from Go; the forge (GitHub) events, the runner and required checks; what a second harness on a clean clone can read. Answer every question, and go deepest on 1, 4 and 5.
--- end

## DOMAIN-B.md

--- begin
# Your domain: decisions, routing, and learning (member B)

You look at the questions as an engineer of decision systems. Weigh most: which decision points need a semantic judgement and which a deterministic rule can settle; how each Jev question is shaped (Choice, Noul or Score, the state it gets, the threshold and its evidence), using only the Jev pages you have; routing of models and harnesses; squads and counterpart verification; a reward computed from the records and applied at a retrospective, with no change to model weights; what happens when Jev is unavailable. Answer every question, and go deepest on 2, 3 and 7.
--- end

## DOMAIN-C.md

--- begin
# Your domain: governance, invariants, and human decision points (member C)

You look at the questions as the auditor of the PSB. Weigh most: the nine System Invariants (`F-0001#1`–`#9`); the Human Decision Points and what counts as planned and unplanned human input (the Task Intervention Rate); the escalation rule and the budget; stalls; requirement acceptance; the trace from requirement to PSB text; which option would let a human become the message bus again (PSB Problem 1). Answer every question, and go deepest on 5, 6 and 8.
--- end
