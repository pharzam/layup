# Evaluator brief — T-hbw8, the hands-on evaluation of the reuse candidates

You evaluate public software that LAYUP could reuse. The Operator asked for a
**thorough, full, hands-on evaluation with no time bound** (decision O-86 in
`operator-decisions.md`). A README summary is not enough: read the code, install
it, run it where it runs, and record what you saw.

## Inputs (read these first, in this directory)

- `problem-statement-brief.md` — the approved PSB. §6 has the twelve **In-Scope
  items** (S1–S12, in the order listed there: S1 Problem Statement Quality,
  S2 Reproducible Discipline Setup, S3 Rule Protection, S4 Stack-Dependent Gates,
  S5 Role Handoffs, S6 Autonomous Clarification, S7 Verification on Every Change,
  S8 Human-on-the-Loop, S9 Stall Resolution, S10 Cost Visibility,
  S11 Specification Synthesis, S12 Harness-Agent Neutrality), the five
  **Human Decision Points**, and the nine **System Invariants** (I1–I9).
- `architectural-vision-brief.md` — the vision brief (squads, cross-verification,
  retrospectives, routing, blind panels, issue-centric communication, telemetry).
- `F-0004-psb-gap-answers.md` — fact 1: LAYUP is Go, **standard library only**,
  and calls Git as the `git` program. `0010-use-go-as-the-technology-stack.md`.
- `operator-decisions.md` — O-76 (LAYUP is only an external binary; a target has
  zero dependency on LAYUP), O-77 (one GitHub account; agents act through a
  GitHub App user token), O-78/O-79 (a pluggable "smart-if" decision provider;
  all thresholds are parameters), O-80..O-83 (parameters with defaults).
- `deep-check-sol.md` (findings 1–32, cite as `Sol-N`), `deep-check-fable.md`
  (M1–M23, cite as `Fable-MN`), `deep-check-author.md` (1–19, cite as
  `Author-N`) — about 50 defects of the withdrawn first-draft architecture. They
  are the checklist of the rewrite. For each candidate, say which of these
  findings a mechanism of the candidate answers, and how.
- `searcher-A.md`, `searcher-B.md`, `summary.md` — what the two blind searchers
  said about each candidate, from the READMEs. **Test their claims; do not copy
  them.** Where the code or a run disagrees with a searcher, say so.

## Where you work (safety)

- Work only under your own directory `EVAL_WORK/<candidate>/` (the path is in
  your task). Clone there, install there, make the test target there.
- **Never write anything under `/Users/farzam/projects/`.** Do not read the LAYUP
  checkout either; everything you need is in this brief directory.
- **No writes to any network service**: no push, no GitHub repo creation, no
  issue or comment, no PR, no npm publish, no sign-up. For a remote, use a local
  bare repository (`git init --bare`). Unset `GH_TOKEN`/`GITHUB_TOKEN` in the
  candidate's environment, and do not give a candidate the host's `gh` login.
  A read-only public fetch (clone, API GET, package download) is fine.
- Installing a package (brew, go install, npm, pip/uv, docker pull) is fine when
  you install it into your work directory or a container where you can. Record
  each install. Do not change the host's global config files (`~/.claude`,
  `~/.gitconfig`, shell rc files). If a candidate insists on writing into
  `$HOME`, run it with `HOME=EVAL_WORK/<candidate>/home`.
- A candidate that needs a coding agent: use Claude Code in print mode
  (`claude -p ... --model claude-haiku-4-5-20251001` or `--model claude-sonnet-5`)
  with a small task on a toy target, and a spend cap where the candidate or the
  harness offers one (`--max-budget-usd` if available). At most about five
  agent runs per candidate; the goal is to see the orchestration work, not to
  build software. Record each run's cost if shown.
- Kill every process you start (servers, daemons, tmux sessions, containers)
  before you finish. Record that you did.
- Do not print or store any secret, token or key in your output.

## The test target

Make a small Go module in your work directory (for example `toy/`, with a
`go.mod`, one package, one test, and a `Makefile` or a check script) and
`git init` it with a local bare remote. Use it wherever the candidate works on a
project. A Go target is the default because LAYUP's first pilots use it; say so
if a candidate cannot take it.

## What to find out, per candidate

1. **Pin it.** The commit SHA you cloned, the license from the LICENSE file (not
   the badge), the language, the dependencies it needs (server, database,
   daemon, tmux, cloud account), the last commit date.
2. **What it is, from the code.** Name the packages or files that implement each
   feature you rate, with `path:line` at the pinned SHA.
3. **What you ran.** The exact commands, and short output excerpts (a few lines
   each) that show the result. Say what failed and why. A step you could not
   run gets the reason and what would be needed.
4. **The twelve In-Scope items S1–S12:** `covers`, `partly` or `no`, each with
   evidence from the code or the run (not from the README alone). `covers`
   means the whole PSB item, as the PSB words it.
5. **The nine invariants I1–I9:** `conflicts`, `neutral` or `supports`, each with
   the reason. Key questions: where does its state live (I1: only in the
   project's Git?); can a human continue the target with it removed (I2); can
   the working agents change its rules or gates (I3); does it set values with
   no evidence (I4); does a check that did not run count as a pass (I5).
6. **The deep-check findings:** list each finding (`Sol-N`, `Fable-MN`,
   `Author-N`) that a mechanism of this candidate answers, and the mechanism in
   one line. Say "none" if none.
7. **Fit with LAYUP's constraints:** LAYUP is Go standard library only (F-0004
   fact 1), so a Go library cannot be a dependency unless an ADR changes that;
   LAYUP is an external binary (O-76). So "use" can mean: run the candidate as a
   separate program next to LAYUP, or depend on its file format / protocol, or
   (rarely) vendor code with its license. Say which, if any, works.
8. **Verdict per reusable part, and one overall verdict:** `use`, `borrow the
   pattern`, or `reject`, each with the reason and the evidence it rests on.

## Output

Write one Markdown file per candidate to `EVAL_OUT/<NN>-<candidate>.md` (the
path is in your task), in this form:

```
# Evaluation: <Name>

| Field | Value |
| ----- | ----- |
| Repository | <url> |
| Pinned commit | <full sha> (<date>) |
| License | <from the LICENSE file> |
| Language and needs | ... |
| Evaluated by | <your model> on Claude Code, <date> |
| Elapsed | <wall-clock minutes> |
| Agent runs and cost | <n runs, cost if shown, or "none"> |

## Verdict
<one line: use / borrow the pattern / reject — and why>

## What it is (from the code)
## What we ran
## In-Scope items S1–S12
(table: item | mark | evidence)
## Invariants I1–I9
(table: invariant | effect | reason)
## Deep-check findings it answers
## Parts and their verdicts
(table: part | verdict | how LAYUP would take it | reason)
## Where the searchers were wrong or incomplete
## Limits of this evaluation
```

**Writing rules.** Write in ASD-STE100 Simplified Technical English: short
sentences, active voice, one instruction or fact per sentence, common words. Do
not use the characters U+2039 or U+203A (the single angle quotation marks) at
all — a check fails on them; write `<...>` instead. Use absolute URLs for
external links (for a file at the pinned SHA, use
`https://github.com/<owner>/<repo>/blob/<sha>/<path>#L<n>`); do not write
relative links. Keep each file under about 400 lines.

When you finish, reply with: the file paths you wrote, the overall verdict per
candidate in one line each, and the processes you stopped.
