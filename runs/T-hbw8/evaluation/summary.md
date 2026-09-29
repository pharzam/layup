# T-hbw8 — evaluation of the reuse candidates: the result

The method and the candidate list are in [`README.md`](README.md). Each
candidate has its own file with the pinned commit, the commands that ran, the
evidence at `path:line`, and a verdict per part. Written by the author (Claude
Opus 5.5) on 2026-09-28, from the eighteen files.

## In plain terms

No candidate can be LAYUP's base, and no candidate is a program that LAYUP should
run next to it. The hands-on runs confirm the two searches and add a reason that
the READMEs do not show: **most candidates break Invariant 5 in their own code** —
a check that did not run, an empty gate, or a cost that is not known counts as a
pass or as zero. But the candidates hold many **working mechanisms** that answer
about two thirds of the deep-check findings. The rewrite borrows those patterns
and writes them in Go with the standard library (`F-0004` fact 1). Two formats
can be used as they are, as options: Spec Kitty's handoff packet v1 and the
isitdone scanner as an advisory program.

## The verdicts

| No. | Candidate | Pinned | License | Verdict | The main reason |
| --- | --------- | ------ | ------- | ------- | --------------- |
| 01 | [Gas Town](01-gastown.md) | `649b832` | MIT | borrow the pattern | State in Dolt and town files (I1, I2); the merge gate and the stall judgement are LLM prose (I6); an empty gate passes (I5) |
| 02 | [Beads](02-beads.md) | `3d76cbd` | MIT | borrow the pattern | Default store is an embedded Dolt database that Git ignores; `bd init` replaces the target's `core.hooksPath` and commits harness files into it (O-76) |
| 03 | [Paperclip](03-paperclip.md) | `0f14d26` | MIT | borrow the pattern | All state in PostgreSQL (I1, I2); a keyless local call acts as the human (I3); subscription spend counts as 0, so the budget stop never fires (I5) |
| 04 | [Agent Orchestrator](04-agent-orchestrator.md) | `3ed3aff` | Apache-2.0 | borrow the pattern | State in SQLite under `~/.ao` (I1); no gates, no rule protection, no escalation rule; always squashes |
| 05 | [Omnigent](05-omnigent.md) | `63231a0` | Apache-2.0 | borrow the pattern | State in SQLite (I1); its supervisor agent widened a worker's permissions by itself (I3); unpriced cost passes the budget (I5) |
| 06 | [Spec Kitty](06-spec-kitty.md) | `af847be` | MIT | borrow the pattern; its handoff packet v1 is usable as a format | Writes about 630 files into `$HOME`; a forced move wrote an "approved" review that never ran; no token record; the trace to the brief is lost |
| 07 | [BMAD-METHOD](07-bmad-method.md) | `ec9a46e` | MIT, with a trademark notice | borrow the pattern | Every judgement is a prompt; a story that traces to a line that does not exist passes; its lint prints a failure and exits 0 |
| 08 | [MetaGPT](08-metagpt.md) | `11cdf46` | MIT | reject; three patterns | Wrote the same PRD 49 times with no stop; no requirement IDs, trace or criteria; filled a gap with a guess; a failed run exits 0 |
| 09 | [Springer](09-springer.md) | `902cd72` | MIT | borrow the pattern | Claude Code only (I9); copies its runtime into the target (O-76); no rule protection; the trace is free text |
| 10 | [GNAP](10-gnap.md) | `2571f55` | MIT | borrow the pattern | Two clones claimed one task and the race left no trace; the commit author proves nothing; no schema files exist |
| 11 | [AgentPlane](11-agentplane.md) | `81fc891` | MIT | borrow the pattern | Part of its state in `.git/agentplane/`, never pushed; puts hooks and `CLAUDE.md` into the target (O-76); "human" approval works from any shell |
| 12 | [Agent Context Store](12-agent-context-store.md) | `ae90049` | MIT | reject; two patterns | The producing agent set its own approval; one file edit gave a role a new right; validation passes empty text |
| 13 | [AI-SDLC](13-ai-sdlc.md) | `6e5b546` | Apache-2.0 | borrow the pattern | A forged envelope with no review verified as `independent`; rule protection is a hook that Bash passes; the scaffold gate runs `pnpm` on a Go target |
| 14 | [sdlc-gh](14-sdlc-gh.md) | `1ab4b6d` | MIT | borrow the pattern | The required checks run the PR's own copy of the gate scripts, so an edited gate passes; its own check fails on a Go repository |
| 15 | [AgentJury](15-agentjury.md) | `4bdad8c` | MIT | borrow the pattern | Public alpha in Python with API keys; the aggregator (about 60 lines) is the part worth having |
| 16 | [Loki Mode](16-loki-mode.md) | `b60ca0e` | BUSL-1.1 | reject | The license names LAYUP's use case as competing, also for a reimplementation from its documents; `loki verify` passed a Go change without `go vet` |
| 17 | [isitdone](17-isitdone.md) | `0051655` | MIT | borrow the pattern; usable as an advisory program | Blocked a build break, a `t.Skip` and two deleted assertions; missed one deleted assertion and a dead condition; warns by default |
| 18 | [no_human](18-no-human.md) | `9b50bf6` | MIT | borrow the pattern | Its fresh-session reviewer found all four planted defects; its tamper guard does not see `*_test.go`; the full loop keeps state in a database |

Loki Mode is the only candidate whose license blocks even the pattern:
LAYUP takes the blind-review idea from AgentJury and no_human (MIT) instead.

## Added on 2026-09-29: Shape Up (O-90)

[Shape Up](19-shape-up.md) is Basecamp's method, a book, not software; the
evaluator read all 24 pages in full. Verdict: **borrow the pattern**. It gives
four mechanisms that no software candidate gave: an **appetite** set by the idea
owner at Intake, so a budget exists before tasks do (Author-7, Fable-M11); the
**circuit breaker**, an absolute cap per milestone that no new record resets, with
the "all downhill" test before any extension (Sol-22, vision 3.5); a **hill
position computed from records** (uphill while unknowns are open; the top when the
acceptance tests are frozen; downhill as they pass), so a ping-pong round is not
progress (Sol-22, Fable-M4, Author-1); and **one bet per milestone** on a
one-screen brief in the pitch's five parts (HDP 1 and 3, Author-19, the approval
load). Rejected: every time value of the book (six weeks, two weeks, three days:
no evidence for agents, I4), "QA is not a gate" (S7, I5, I9) and "no backlogs"
(I1). The book's default "cancel at the cap" becomes "stop, then reshape or
escalate", because dropping a requirement is business-forking.

## What the rewrite takes, by concern

Each row is a mechanism that ran or was read at the pinned commit. The last
column names the deep-check findings it helps to answer
([`../deep-check-sol.md`](../deep-check-sol.md) `Sol-N`,
[`../deep-check-fable.md`](../deep-check-fable.md) `Fable-MN`,
[`../deep-check-author.md`](../deep-check-author.md) `Author-N`). A pattern is an
input to the design, not a decision; the rewrite decides.

| Concern | Pattern | From | Findings |
| ------- | ------- | ---- | -------- |
| One writer of the records | Sessions return a typed result file; only the supervisor admits it and commits; a result bound to an older state hash is refused | AgentPlane | Fable-M12, Sol-3, Sol-4 |
| Records in Git | One file per record, written with `encoding/json`; an append-only event log per task, state derived by code at read time, never stored | GNAP, Springer, AO, ACS | Fable-M12, Sol-3, Author-9 |
| Claims | A claim is a compare-and-set granted by the one writer; a lease per session with heartbeat and reclaim after a time-out | Paperclip, Beads | Sol-4, Author-10 |
| Identity and decisions | The actor comes from a verified credential, never from text the client writes or from the commit author; a decision is one typed record (pending, approved, rejected, ...), and a comment never changes it; no keyless trusted path | Paperclip, GNAP (as the failure), AgentPlane | Fable-M1, Fable-M3, Sol-5, Sol-30 |
| Handoffs | A transition table (from role, to role, required artifacts, required state computed by code); typed result statuses (`completed`, `blocked`, `needs_context`, `failed`) with a required reason; a plan dataflow check; per-kind JSON schemas | AgentPlane, ACS, Springer, MetaGPT | Author-8, Fable-M14, Sol-21 |
| Who acts at each step | A default table of step to role, prompt and harness, which the Operator can replace (O-81) | Spec Kitty, MetaGPT, Omnigent | Fable-M16, Sol-32 |
| Harness sessions | A table from harness to flags (permission, model, session id, prompt); sessions start with no user-level setup; hooks report start, stop and permission requests as facts; one `prime` command gives the context to any harness | AO, AgentPlane, Gas Town | Fable-M23, Fable-M16 |
| A question during a session | The session writes a question record and stops or waits; the owner answers; the answer goes back into the same session (ran in 34 s in AO) | AO, Omnigent, AI-SDLC | Author-2, Sol-9, Sol-28 |
| Pull request and merge | A draft pull request until the gates and the review pass, then ready; merge only at the verified head SHA with green required checks; a review record keyed by the head SHA with a verdict enum | AI-SDLC, AO, Gas Town | Sol-2, Sol-27, Fable-M9 |
| Gates | Native commands per stack from a catalog; "no pass without a run on the exact tree hash"; SKIP is a failure; the verifier never runs code from the checkout it checks and reads its rules from the base branch | sdlc-gh, isitdone, AI-SDLC | Sol-2, Fable-M19, Sol-26 |
| Rule protection | A ruleset applied by API and read back; code owners on every path that a required check runs; a hash of the rule files recorded at plan approval, and a task stopped when it changes; a drift manifest | sdlc-gh, AgentPlane | Sol-1, Fable-M2 |
| Counterpart verification | A fresh, single-turn, read-only reviewer on another harness or model, told to refute "done", with a `file:line` checklist, failing closed; the reviewer harness checked by code, not by a field an agent writes | no_human, AI-SDLC (as the failure), Omnigent | Sol-24, S12 |
| Stalls | A review-round limit between the same two roles, then escalation; attempt caps per task across resumes; wait states for a human kept apart from "no signal"; restart backoff and a crash-loop stop | Paperclip, Gas Town, AO, AgentPlane, no_human | Author-1, Fable-M4, Fable-M5, Sol-22 |
| Blind panel | N fresh sessions with the same sealed input, run in parallel; a deterministic aggregator with quorum, a provider floor and `insufficient_jury` when a member fails | AgentJury | Fable-M18, Fable-M6, Sol-24 |
| Cost | A ledger row per session: harness, provider, billing type, cost status (`reported`, `unpriced`), model, token classes, money, plus latency and duration; `token_usage` as `observed`, `partial` or `unavailable` with a reason — never a guessed zero | Paperclip, AgentPlane, GNAP | Sol-11, Sol-12 |
| Budget | Budgets on scopes that exist (project at Intake, milestone and task scopes when Plan makes them), a warn level, a hard stop that pauses and refuses the next start, an override as a typed decision; the stop kills the running session | Paperclip | Author-7, Fable-M11 |
| Specification | Numbered source lines as requirement IDs and a `covers` field on each item, checked byte for byte by code; every requirement maps to a task; an acceptance matrix where `pending` is not `pass`; a typed list of gaps from the spec role, each gap a question in the one batch | BMAD, Spec Kitty, MetaGPT, Springer | Sol-17 to Sol-20, Fable-M20, Fable-M21, Author-5, Author-18 |
| Retrospective | Generated by code from the event log; proposals are data, applied only by an approved step; rule changes need human promotion | Spec Kitty, BMAD, Springer | Fable-M17 (in part), Sol-25 (in part) |
| Parameters | Each default with its evidence line next to it; each change a record with its actor | no_human, Paperclip | Fable-M22, Invariant 4 |

## What no candidate answers

These findings have no public pattern in the eighteen candidates. The rewrite
designs them from the PSB and the Operator's decisions:

- **The escalation screen and the smart-if points** (Sol-6, Sol-7, Sol-8,
  Fable-M7, Fable-M13, Author-4): no candidate has a rule that selects
  business-forking decisions. Paperclip and AgentPlane let an agent decide by
  itself to ask.
- **Question routing by ambiguity kind, with a "needs a human" branch** (Sol-10,
  Fable-M8, Author-3).
- **The PSB measures** that need a population or a baseline (Sol-14, Sol-15,
  Sol-16, Sol-31, Fable-M15): no candidate counts unplanned human input, early
  questions or first-review acceptance.
- **Learning that changes routing and reaches the next project** (Fable-M17,
  Sol-25, Sol-23, Sol-29, Author-12): no candidate routes models by task
  complexity or measured reward; retrospectives stop at proposals.
- **The order of Intake before the target exists** (Author-6, Fable-M10) and the
  Armature baseline itself (Invariants 4, 7, 8): no candidate names Armature.

The entries that the searchers kept but did not shortlist include tools for some
of these parts (for example Claudexor for quota routing and cross-family review,
llm-panel for blind multi-CLI panels, fractal for cost and depth limits). If the
rewrite plan finds a part with no design, it reads those entries first.

## Failures that repeat across candidates

These are defects of the candidates, seen in the runs. The rewrite turns each
one into a test that LAYUP's own design must pass:

1. **A check that did not run counts as a pass** — Gas Town (empty gate),
   Loki Mode (`loki verify` without `go vet`), BMAD (lint exits 0 on failure),
   MetaGPT (a failed run exits 0), isitdone (exits 0 when it finds no checks).
2. **Unknown cost counts as zero** — Paperclip (subscription runs), Omnigent
   (unpriced actions pass the budget), MetaGPT.
3. **The actor comes from something an agent can write** — GNAP and Beads (the
   commit author), Paperclip (a keyless call), AgentPlane (`--by USER` from any
   shell), ACS (the agent sets its own approval), AI-SDLC (a marker file and a
   harness field).
4. **The gate runs the copy under review** — sdlc-gh and AI-SDLC's scaffold run
   the pull request's own gate scripts.
5. **The tool writes into the target or the host** — Beads, AgentPlane, Springer
   and ACS put their files into the target (against O-76); Spec Kitty, Loki Mode,
   ACS and Gas Town wrote into the user's home directory.
6. **State that a Git clone does not carry** — every orchestrator (a database,
   `.git/` internals, or files outside the repository).

## A correction to the second search

[`../search-v2/summary.md`](../search-v2/summary.md) says its table is the union
of the two shortlists; it left out seven candidates (see [`README.md`](README.md)).
This evaluation covers them. One of them, Agent Orchestrator, gave the working
mid-session question path and the merge preconditions above. Where a file says a
searcher was wrong (for example "Paperclip's database is not stated": it is
PostgreSQL), the section "Where the searchers were wrong or incomplete" of that
file gives the evidence.

## Incidents during the evaluation

The brief said: no writes to the host's global configuration. Three candidates
wrote there anyway, and one evaluator went past the run limit. Each was found
and reverted by the evaluator; the author checked the result.

1. **Loki Mode** installed a Claude Code plugin (`caveman`) at user scope and
   turned it on in `~/.claude/settings.json`, because the evaluator's wrapper
   gave Loki's calls the real `HOME`. The evaluator uninstalled it and moved its
   files to the work directory. The author found no trace of it afterwards. Two
   Claude Code sessions started while it was on and may have run its
   session-start hook.
2. **Spec Kitty** wrote 628 command and skill files for 13 harnesses into the
   real `$HOME` (`~/.claude/commands`, `~/.claude/skills`, `~/.gemini`,
   `~/.cursor`, `~/.github/prompts`, `~/.kittify`, ...) on its first calls. The
   evaluator listed them from a copy in a fake home and deleted exactly those
   files. The author checked: none of those directories remains, and the
   existing skills in `~/.claude/skills` are intact.
3. **AI-SDLC**'s `init -y` tried to apply branch protection on GitHub without
   being asked. It failed only because `gh` had no login in the isolated home.
4. **Gas Town** started 22 Claude Code sessions by itself for one 9-line task
   (about USD 3.25 at list price), over the brief's limit of about five runs;
   the evaluator stopped the town after the first merge. Gas Town and Dolt also
   wrote `~/.gt` and `~/.dolt`; the evaluator removed both, and the author
   checked.
5. Not explained: `~/.local/bin/agy` has a modification time of 15:46:58, in
   the window of incident 2. No evaluation log names it; the cause is not known.

The lesson for LAYUP is the same as failure 5 above: a harness or tool that
LAYUP starts must run with a home directory that LAYUP gives it.

## Effort and cost

- Seven evaluator sessions (Claude Opus 5.5 on Claude Code), in parallel, from
  about 15:30 to 17:00 on 2026-09-28; about 2.2 million tokens in total, as the
  sessions reported them.
- Candidate agent runs: about USD 36 as reported by Claude Code (Spec Kitty 11.20,
  Springer 11.98, BMAD 4.45, Gas Town about 3.27, AI-SDLC 2.70, the rest under
  USD 1 each). Some runs were not reported (no_human, two AO sessions, Loki
  Mode's reviewer calls). The host uses a subscription login, so these are list
  prices, not money spent.

## Limits

- No second harness ran: Codex and Gemini are not installed on this host, and the
  external harnesses were not used for candidate runs. Every cross-harness claim
  was read in the code or run with Claude on both sides.
- No candidate ran against GitHub (pull requests, rulesets, CI). Those parts were
  read in the code and, where possible, run with a local bare remote.
- Several runs stopped early (Spec Kitty after the first work package, Springer
  at the first merge gate, AgentPlane before a task reached done).
- The evaluators are one model, the author's. The plan review of the rewrite and
  the reviews of its slices use other models.
