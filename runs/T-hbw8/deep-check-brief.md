# T-hbw8 — the brief of the deep check (2026-09-28)

Asked by the Operator after the `layup psb check` overclaim: check the architecture and ADR-0013 to ADR-0020 against the PSB and the vision brief for defects of the same class. Two fresh reviewers, each in a disposable clone at `6368b57`, got this brief word for word; the author wrote an own pass before reading theirs.

--- begin
# Deep check of the LAYUP architecture against the PSB and the vision brief

You are a fresh, independent reviewer. The author of the change is Claude Opus 5.5; you are a different model. This directory is a disposable clone of `pharzam/layup` at branch `T-hbw8`, head `6368b57`. Do not change any tracked file.

## Read

- **The architecture under review:** `docs/architecture.md` and `docs/adr/0013-*.md` to `docs/adr/0020-*.md`. They amend `docs/adr/0011-*.md`.
- **The sources it must agree with:** the problem statement `docs/facts/problem-statement-brief.md` (numbered facts `docs/facts/F-0001-*.md`, `F-0003-*.md`, answers `F-0004-*.md`); the vision brief `docs/facts/architectural-vision-brief.md` (`F-0002`; a solution input, not requirements); the Operator's decisions `runs/T-hbw8/operator-decisions.md` (O-76 to O-84) and `runs/T-hbw8/inputs-from-pr-69.md` (O-66 to O-75 and the review of the previous architecture); the requirements `docs/prd/PRD-0001-layup.md`; the smart-if provider's facts `runs/T-hbw8/jev-sources.md`.
- The existing code: `internal/psb/check.go` (what `layup psb check` really does).

## The lens

The Operator found this defect: the architecture said `layup psb check` "finds the gaps of a problem statement", but it is five fixed text rules and cannot find a gap of meaning. Find **every defect of the same class and its neighbours**:

1. **Impossible or overclaimed mechanisms.** A component, command, check or provider that the text says does X, where the stated mechanism cannot do X (for example: a deterministic check asked to judge meaning; the smart-if asked to write text, count, compare dates or do arithmetic; a GitHub feature that does not exist or does not work on a repository owned by one user; a measurement that no record supports).
2. **Gaps against the PSB.** A PSB problem, In-Scope item (`F-0003#41`–`#52`), invariant (`F-0001#1`–`#9`), Human Decision Point, success criterion or measure (PSB §7) that the architecture does not really answer, or answers only in name. Check each one; do not sample.
3. **Contradictions** with the PSB, with the Operator's decisions O-66 to O-84, between two ADRs, or between an ADR and `docs/architecture.md`.
4. **Gaps against the vision brief** that O-67 keeps in scope (squads, Jev/Laya routing, blind panels, the learning loop, retrospectives), where the architecture claims coverage it does not have.
5. **Missing actors or steps:** who or what does each thing the phase loop needs (for example, who writes the product code, who opens and reviews the pull request, who starts a panel, who writes a stall diagnosis, where a role session's context comes from).

## Output

Write `.review-out/deep.md`, then print it as your final answer. ASD-STE100 Simplified Technical English. At most 250 lines. For each finding:

- **Where:** file and the exact sentence (quote it).
- **Source:** the PSB fact, vision item or O-decision it breaks, or the reason the mechanism cannot work (quote the fact or the documentation).
- **Why it matters:** a concrete failure in a pilot.
- **Fix:** one or two lines.
- **Severity:** `material` (the architecture is wrong or false) or `note`.

Sort by severity. End with a list of the PSB items you checked and found answered, one line each, so the author can see your coverage.
--- end
