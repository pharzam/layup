# 0017. Stop a stall at a counted limit and examine it with a fresh context

Date: 2026-09-25

## Status

Accepted

## Context

A stall is a task that does not reach its goal and does not fail cleanly, because role agents do not agree or because a step repeats without progress (`F-0001#37`); the PSB requires a procedure with a limit, an independent examination with a fresh context that gives a diagnosis, and the diagnosis and the outcome kept in the project repository (`F-0003#49`, REQ-009, REQ-010); a stall is a planned human decision point when its procedure reaches its limit, and the Operator may give the task to a different harness agent (`F-0001#14`). Question 5 of `PRD-0001` §11.

The panel (`runs/T-7qvc/`) compared: fixed limits — a count of identical-input attempts or a time without a progress event, a disagreement stopped at the review cycle cap (A S1, B 5A, C 5-A, 5-B); pre-registered per-step limits (A S2); a progress metric (B 5B); a cost limit from telemetry (C 5-C); and, for the examiner, a fresh session on a different harness, on the same harness with a different model, or the Operator. Every numeric value is a working hypothesis until the pilot's baseline (Invariant 4). The Operator decided (#66): **O-58**, N = 1 identical-input retry (the second identical attempt stops the step), T = 15 minutes without a recorded progress event, a review disagreement is a stall when the cycle cap is reached with a material finding open, and the pilot measures both for the Operator to reset in one batch; **O-59**, the examiner is a fresh session that did none of the work, with a different model; a different harness where one is available; the same harness allowed and recorded.

## Decision

We will give a stall an explicit state, counted limits, a fresh examiner, and a package for the Operator:

1. **The state.** A stall is an open row of `runs/<task>/stalls.tsv` (`id, kind, limit, limit_value, opened_sha, examiner_harness, examiner_model, session, diagnosis_path, outcome, closed_at`), `kind` being `no-progress` or `no-agreement`. While a row is open the App's check `layup/stall` ([ADR-0013](0013-run-stack-gates-through-a-layup-github-app.md)) posts `failure` on the task's pull request, and the task record's status is `stalled`. An open stall with no diagnosis is incomplete, never a pass (`F-0003#61`).
2. **Progress and the limits.** A **progress event** is a changed gate or check result, or a new evidence artifact under `runs/<task>/`; a repeated prompt, an empty commit, or a run with the same result is not progress. `no-progress` opens when the same fingerprint (gate, rule, input hash) fails a second time — **N = 1** retry — or when **T = 15 minutes** pass without a progress event (O-58). `no-agreement` opens when a review reaches its cycle cap with a material finding open. Both values are working hypotheses; the pilot measures them and the Operator resets them in one batch (Invariant 4).
3. **The examiner.** The engine opens the examination: a fresh session that did none of the stalled work, with a model different from every model in the task's `handoffs.tsv` ([ADR-0015](0015-roles-are-the-kit-gate-roles-and-a-handoff-is-a-tsv-row.md)), on a different harness where one is available and on the same harness otherwise, recorded (O-59). Its brief is the task's records — the trigger, the attempts, the opposing positions, the input commits and logs — and not the performers' reasoning. It writes `runs/<task>/stall-<id>.md` with the headings Trigger, Evidence, Diagnosis, Outcome and Operator package, and says what it was given and what it could not examine (the independence standard of Who may review).
4. **The outcome.** Either a verified recovery (a change that passes the gates, recorded as the row's outcome) or the **package for the Operator**: the stall record, the handoff rows and the linked commits, reported on the task's issue, with the one decision the Operator must make; the Operator may answer or give the task to a different harness agent (`F-0001#14`). The Operator's answer is a record in Git before the task continues. When no examiner is available, the package says `diagnosis pending`, and REQ-010 is not reported as met.

We rejected: **pre-registered per-step limits** — more values to set without evidence; its ceiling (never more than one identical-input retry) is kept in rule 2; **a progress metric alone** — a real refactor raises the failure count for a while and is flagged; its definition of a progress event is kept; **a cost limit from telemetry** — no baseline exists before the pilot; **the Operator as the examiner** — REQ-010 needs the fresh-context diagnosis before the Operator, and `F-0003#73` counts stalls that close without human input.

## Consequences

- REQ-009 and REQ-010 are testable: a seeded stall of each kind stops at the limit, a fresh examiner writes the diagnosis, and the Operator receives the package (`PRD-0001` §7.1).
- The time limit measures the harness's speed as well as progress (member C); a slow harness stalls sooner. The pilot reports stalls the examiner called premature, and T is reset from that.
- The engine gains `layup stall open|examine|close`; the check `layup/stall` is one more required check of the App.
- This repository's own review protocol (§ Bootstrap mode, rule 3) already stops at the cycle cap with a material finding; the `no-agreement` kind names that state for a target.
