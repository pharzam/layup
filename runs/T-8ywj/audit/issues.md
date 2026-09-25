# Audit of the open issues of pharzam/layup — 2026-09-25

Read-only. Sources: `gh issue view --comments`, the GitHub API (229 comments, all issues), `gh pr list --state all`, the branch list, and the clone's `T-7wqc` branch. Nothing was posted or changed. Scratch copies: `/private/tmp/claude-501/-Users-farzam-projects-layup/26cc4c98-384d-4636-8768-8ceed05e8160/scratchpad/audit/` (`issue-N.txt`, `body-N.txt`, `all-comments.json`, `all-issues.json`, `prs.txt`).

## 1. Open issues

| # | Task | Kind | Opened / last comment (UTC) | State | Blocked by | Consumed |
|---|------|------|------|-------|-----------|----------|
| 55 | T-gkps: store the revised routing policy as F-0006 | Process | 09-25 09:00 / 09-25 10:35 | Plan rev. 1 `reject`; plan rev. 2 `reject`. No branch, no PR. | Plan review 2, condition 1: "Resolve R11 before building. Revise the goal and slices or split the independent outcomes into child issues." | 2 plan reviews, 0 rounds, O-35–O-40 |
| 54 | T-7wqc: ADR-0012 (Proposed), successor of #53 | Process | 09-24 12:28 / 09-24 13:37 | Frozen head `30ca13f` (cycle 2, the cap); round 3 `not mergeable, findings recorded`. Paused. | O-34 "Pause here": "Leave #54 open as not mergeable, with the findings recorded, and stop for now." | 1 plan review, 3 rounds, O-32, O-33, O-34 |
| 52 | X2 of ADR-0012: stall protocol | Process | 09-24 08:58 / none (body edited 09-24 12:17) | Not started. | Child of the unmerged ADR-0012 (#54). Body: a change to `F-0001#14` "needs a PSB revision by the idea owner". | 0 |
| 51 | X1 of ADR-0012: review independence by harness | Process | 09-24 08:58 / none | Not started. | Child of the unmerged ADR-0012 (#54). | 0 |
| 49 | R11: how goal classes of a decision record are counted | Process | 09-24 07:19 / none | Not started. | Nothing named. O-28 and O-33 both say "#49 still fixes the R11 text". | 0 |
| 48 | Line-ending pin for raw facts and setup-check.sh | Process | 09-24 07:15 / none | Not started. Revealed by #47 round 1. | Nothing named; "does not block #47". | 0 |
| 45 | T-zmj6: gap batch + idea owner's answers (F-0004) | Product (PDR child 1) | 09-23 20:31 / 09-23 20:39 | Plan posted; plan review `approve-with-conditions` (Claude Fable 5.1); conditions answered; AC rewritten. No branch, no decision note, no PR. | Not named. The next step is the idea owner's answers to the 19 questions (Decision Point 2, O-14); none is recorded. | 1 plan review, 0 rounds, 0 new O |
| 42 | T-meh2: PDR (spec, architecture, plan) | Product (parent) | 09-23 19:59 / none | Child 2 (#43, T-ertw) merged via PR #44; child 1 (#45) as above; children 3–6 have no issue. | Waits for #45; the Operator's approval (O-14). | O-14, O-15 (#42's numbering) |
| 34 | ADR-0011 wording residuals; glossary 'rule path' | Process | 09-23 18:44 / none | Not started. | Nothing. | 0 |
| 33 | T-vk3k: Go stack gates run outside the target | Product (#29 child 5) | 09-23 18:44 / none | Not started; no AC yet. | Body: AC "to be set … when it starts, after `T-dq05` and `T-b97r`" (T-b97r has no issue). #29: "#33 and the planned `T-b97r` and `T-tmhw` are re-scoped by the PDR's plan." | 0 |
| 29 | T-stfn: Step 2, the core engine | Product (parent) | 09-23 17:29 / 09-23 19:59 | Children merged: #30 (ADR-0011, PR 31), #38 (Go skeleton, PR 39; replaces #32, PR 36 closed unmerged), #37 (`psb check`, PR 40), #35 (CI protection, PR 41). Remaining: T-b97r, #33, T-tmhw. | Comment "On hold for the PDR": "No new child of this issue starts; #33 and the planned `T-b97r` and `T-tmhw` are re-scoped by the PDR's plan." → #42 | O-11, O-12, O-13 recorded here |
| 24 | docs/ci/README.md claims about linter from PR head | Process | 09-23 14:53 / 09-23 16:55 | Not started; scope added twice. | Nothing. | 0 |
| 21 | markers check: quoted paths; stale T-nfh8 sentences | Process | 09-23 14:38 / 09-23 16:23 | Not started; scope added once. | Nothing. | 0 |
| 15 | Kept docs say docs/decisions/ still exists | Process | 09-23 12:07 / none | Not started. | Nothing. | 0 |

Product work: #29, #33, #42, #45 (4 of 14). None has started a build since 2026-09-23 20:39. Process work: 10 of 14.

### Branches on GitHub

| Branch | Head | Meaning |
|---|---|---|
| `main` | `43a05c8` | Merge of PR #50 (#47, F-0005). |
| `T-q1x6` | `919e2f0` | #47, merged. |
| `T-w79d` | `9c989ef` | #46, closed as split; no PR. |
| `T-9tgn` | `2592624` | #53, closed as split; no PR. |
| `T-7wqc` | `30ca13f` | #54, paused; no PR. |

No branch exists for #45, #55 or any other open issue.

### Pull requests

No open PR. All 19 PRs (#13–#50) are merged except #36 (closed unmerged, replaced by #39). Every PR closes a now-closed issue. `Refs #29`: PRs 31, 36, 39, 40, 41. `Refs #42`: PR 44. No PR was ever opened for #46, #53, #54 or #55 (each split or pause comment says so).

## 2. The ADR-0012 line

| When (UTC) | Issue | Event |
|---|---|---|
| 09-24 06:39 | #46 opened (T-w79d) | Operator pastes the routing policy (84 lines). Plan 1. |
| 06:49 | #46 | Plan review 1, GPT-6 Astra on Devin: `reject` ("at least 6" goal classes; wrong test command; #45 C5 dependency). |
| 06:50–07:37 | #47 (T-q1x6) | Split out to store the text as F-0005. Plan review (Gemini 3.1 Pro, `approve-with-conditions`), round 1 (Fable 5.1, `material`, 8 findings, 2 material), round 2 (Gemini, `nothing material`). PR #50 merged 07:37. **The only merge on the line.** |
| 07:03–07:16 | #46 | Plan rev. 2; plan review 2 (Astra): `reject` (R11 again; base not fixed). |
| 07:19 | #46 | **O-15**: "a decision record is one goal, as in #30, even with a check of its parts in the DoD." #49 opened to fix the R11 text. |
| 07:38–07:48 | #46 | Plan rev. 3; plan review 3 (Astra): `approve-with-conditions`, budget 1,100 lines / 12 files, cap 2. |
| 08:09 | #46 | Clause map: 13 model sessions, 1,019,459 tokens, 13.3 min; judge (Gemini on AGY) run 1 failed, run 2 settled 4 disputes. |
| 08:10–08:39 | #46 | Panel: A (Fable/Claude Code), B (Astra/Devin), C (Gemini/AGY) answered; D (Kimi/OpenCode, then Qwen) failed twice, dropped. **O-16, O-17, O-18.** |
| 08:45–09:23 | #46 | Freeze `6328f69`; round 1 (Sol/Devin) `material`, 8 findings; freeze `4e1d5b3`; round 2 (Astra) `material`, 5 findings; freeze `9c989ef`; round 3 (Sol) `not mergeable`, 4 findings. |
| 10:14 | #46 closed as split → #53 (T-9tgn) | **O-26** (successor takes all 4), **O-27** (no Fable 5.1, no Astra, no AGY this session). |
| 10:24 | #53 | Plan review 1 (Sol/Devin): `reject` — "O-15 expressly settled #46 only, not this successor." |
| 10:32–10:44 | #53 | **O-28** (extend O-15 to #53). Plan rev. 2. Devin refuses every model ("Upgrade to Pro"). **O-29** (Kimi K3 first, Sonnet 5 backup). |
| 11:15 | #53 | Plan review 2: Kimi no output in 1,200 s; Claude Sonnet 5 backup, `approve-with-conditions` (execution-tier model on a reasoning part, recorded). |
| 11:21–11:47 | #53 | Round 1 attempts: Kimi (stopped), Grok 4.7 (hung), Grok short prompt (watchdog stop), Sonnet (stopped) — none produced a review. **O-30** (model recommendations), **O-31** (stop and switch). Devin works again at 11:41. |
| 12:01–12:27 | #53 | Round 1 (Sol) `material`, 6; round 2 (Sol) `material`, 10; round 3 (Sol) `not mergeable`, 2 findings. |
| 12:28 | #53 closed as split → #54 (T-7wqc) | **O-32**, **O-33**. |
| 12:36–13:36 | #54 | Plan review (Sol) `approve-with-conditions`, 1,450 lines / 13 files; rounds 1 and 2 (Sol) `material`, 4 + 4; round 3: Sol quota exhausted after 545 s → Sonnet 5 backup, `not mergeable`, 2 findings. |
| 13:37 | #54 | **O-34** "Pause here". |
| 09-25 09:00 | #55 (T-gkps) | Operator pastes a revised policy (362 lines) that "supersedes my earlier routing answers where they conflict". **O-35, O-36, O-37.** |
| 10:22 | #55 | Plan review 1 (Sol): `reject`, "at least 6" goal classes. |
| 10:24–10:28 | #55 | **O-38** (revise as one goal), **O-39** (reviewers fixed), **O-40** (Devin for the re-review; Go rate-limited). |
| 10:35 | #55 | Plan review 2 (Sol, fresh session): `reject`, "at least 4" goal classes: "Calling their V1–V10 aggregator 'one test' does not meet R11's single mechanism and single test exception." |

### Totals for the line (#46, #47, #53, #54, #55)

- Plan reviews: **9** (5 `reject`, 4 `approve-with-conditions`). Five of the nine verdicts turn on the R11 goal count, which #49 was opened to settle and which is not started.
- Review rounds: **11** (#46: 3, #47: 2, #53: 3, #54: 3, #55: 0). Plus 4 panel dispatches (3 answered), 13 clause-map sessions, 2 judge runs.
- Frozen heads: 9 (`6328f69`, `4e1d5b3`, `9c989ef`; `4c03714`, `bb680c9`, `2592624`; `116612d`, `4f8f0d1`, `30ca13f`). Findings fixed inside cycles: 37 (8+5, 6+10, 4+4). Findings carried across splits: 4 then 2. Open now: 2.
- Operator decisions decided: **19** (O-15–O-18, O-26–O-40). Open questions inside ADR-0012: O-19–O-25 (7).
- Calendar days: 2 (09-24 06:39 → 13:37, about 7 h; 09-25 09:00 → 10:35, about 1.6 h).
- Reviewer machine time, summed from the run footers: about 3.4 h of runs that returned a record, plus about 1.4 h of runs that returned nothing (hangs, quota, refusals).
- Merged output of the line: F-0005 (#47). ADR-0012 exists only on branch `T-7wqc`.
- Effect of #55 on #54: F-0006 §2 says F-0005 "is incomplete: its Markdown fence does not close … I no longer have the missing original text", and §16 says "GPT-6 Sol is again permitted", "AGY remains removed", "The 70% and 90% quota bands remain removed". ADR-0012 at `30ca13f` cites F-0005 by line in all 25 clause rows and lists 70/90 in O-22; it does not cite F-0006.

### Open findings that stopped #54 (round 3, cycle 2, reviewer Claude Sonnet 5), word for word

1. "`docs/guardrails.md:298-301` mislabels two fail-expected self-test cases as a "valid" control." Basis: "the sentence reads 'the table ends only at a blank line or another block, so read rows by that rule and keep a valid pipe-less control (cases 2e, 1f, k1, k2).' But in `runs/T-w79d/clause-table-check.sh:289-292`, `2e` and `1f` are `mutate` calls (case `2e` expects `clause-table: FAIL P2 a row opens with "C25"`; case `1f` expects `clause-table: FAIL P1 no delimiter row under the clause table header`) — only `k1` and `k2` are `control` calls, which must print bare `clause-table: OK`." Classification: "in the change (introduced whole in commit `14ea66b`, this task's first fix commit; still present, unfixed, at `30ca13f`)."
2. "D1's stated application order and D3's own "eligible" definition give two different routes for a review step whose table lists a row with the author's own model." Basis: "D1 (`docs/adr/0012-*.md:50-53`) states the order as "the review levels (D7), the quota state (D6), the table order (D3)." D3 (`:59-65`) defines **eligible** using exactly three criteria — smoke-run window, "has not failed in this task (D4)", and "the state of each of its quota pools allows the step (D6)" — and never cites D7. D7 (`:90-93`) requires, "for every review, whatever its risk," a binding "whose model differs from the author's."" Classification: "in the change — the interaction is created by this task's own edit to D1 … and to D7 …; D3 itself was never updated to match either edit."

O-34 restates them: "(1) `docs/guardrails.md` §2 calls cases 2e and 1f "valid" controls; only k1 and k2 are controls. (2) D3's definition of "eligible" does not include D7, so D1's order (D7, D6, D3) and D3's predicate can route a review to the author's own model." The options offered and not chosen: "a successor that fixes both and folds the route into one definition in D3 (the author's recommendation); a successor with the two minimal fixes; landing with two child issues."

## 3. Operator decisions O-1 … O-40

| O | Issue, date | Decided |
|---|---|---|
| O-1 | #8, 09-23 | Worktree directory `.worktree/`. |
| O-2 | #8, 09-23 | Evidence store `runs/<task>/`. |
| O-3 | #8, 09-23 | **Model tiers**: reasoning Claude Opus 5.5, Claude Fable 5.1; execution Claude Sonnet 5, Claude Haiku 4.5. |
| O-4 | #8, 09-23 | Each new ADR gets a panel before selection. |
| O-5 | #8, 09-23 | Go test layout and levels. |
| O-6 | #8, 09-23 | Security scanners: govulncheck, go vet, gitleaks. |
| O-7 | #8, 09-23 | Coverage gate: none yet, open gap. |
| O-8 | #8, 09-23 | Task-ID scheme `T-` + 4 characters. |
| O-9 | #30, 09-23 | Rule protection (Invariant 3): decide later. |
| O-10 | #30, 09-23 | Nothing from LAYUP in a target beyond the Armature kit. |
| O-11 | #30/#29, 09-23 | "No LAYUP-written files" in the target; stack gates run outside it. |
| O-12 | #29, 09-23 | What LAYUP is: an orchestrator product, not a scaffold. |
| O-13 | #29, 09-23 | Clarification of O-12 (verbatim, long). |
| O-14 | #42, 09-23 | The Operator approves the PRD and PDR: "I approve it myself". |
| O-15 | #42, 09-23 | PRD scope: "full LAYUP scope with phases". **Collision:** #46 reuses **O-15** on 09-24 for "a decision record is one goal, as in #30, even with a check of its parts in the DoD"; #49, O-28 and O-33 cite the #46 meaning. |
| O-16 | #46, 09-24 | Keep the ADR-0012 boundary; reviews routed only as "a binding that Who may review accepts"; quota states inside; "'unknown' is never Normal". |
| O-17 | #46, 09-24 | ADR-0012 routes this repository's gate, not the product engine. |
| O-18 | #46, 09-24 | Verbatim: "for the later panel after current work done, max panel time should be 15 min, and each response timeout should not be more then 5min, blind harness can be same but model shuold be different , for fable 5.1 and ASTRA 6.1 in extream situation before the staleness status should be used". |
| O-19 | ADR-0012 body | **Open question**: harness names ("Davinci CLI" = `devin`? "OpenCoder CLI" = `opencode`?), the order criterion, where the table lives. |
| O-20 | ADR-0012 body | **Open question**: which bindings join which tier; is "Anthropic models exclusively" a rule; what "free tier" means; does a paid fallback need authorization. |
| O-21 | ADR-0012 body | **Open question**: do "Astra 6.1" and "Astra GPT / Astra 6" name one model or two; exact model and effort; harness change on AGY/OpenCode. |
| O-22 | ADR-0012 body | **Open question**: thresholds (70 %, 90 %), vendor quota window vs. reporting windows, figure source and age, steps allowed in `unknown`, smoke-run window of D3, dispatch time limit. |
| O-23 | ADR-0012 body | **Open question**: "panel arbitration" = selection over options, or a panel verdict. |
| O-24 | ADR-0012 body | **Open question**: AGY prompt preparation, scoping, option answering, scoring — Determinism first? |
| O-25 | ADR-0012 body | **Open question**: prompt packages committed in full, as a manifest, or not at all where sensitive. |
| O-26 | #53, 09-24 | Successor takes all 4 findings; #46 closes as split. |
| O-27 | #53, 09-24 | Verbatim: "For this session, do not select Fabel 5.1 or GPT 6 Astra , AGY for any task, review, panel, or fallback." |
| O-28 | #53, 09-24 | O-15 extended to #53. |
| O-29 | #53, 09-24 | Reviewers: "Kimi K3 on OpenCode (standalone) … If a run gives no answer in 20 minutes, Claude Sonnet 5 on Claude Code runs it instead, and the record names the tier mismatch." |
| O-30 | #53, 09-24 | Verbatim: "for opencode Recommendation: use Space Bunny Free for free private work. For frontier reasoning on OpenCode now, use Grok 4.7 on Go. Do not send confidential code to the free models that collect data. for Devin Recommendation: for reasoning reviews, use GPT-6 Sol (xhigh). For cheap reasoning work, use GPT-6 Luna. Use SWE-2 for coding while it is free." |
| O-31 | #53, 09-24 | Verbatim: "if kimi or GPT didn't work on DEVIN, stop the work in switdh to alternative". |
| O-32 | #54, 09-24 | Successor takes both #53 findings; GFM rows; D5 same-tier fallback first. |
| O-33 | #54, 09-24 | O-15 applies to every successor of #46's goal. |
| O-34 | #54, 09-24 | "Pause here": #54 stays open, not mergeable. |
| O-35 | #55, 09-25 | Capture the revised policy as F-0006 now; "ADR-0012 resolves points 1–4"; no labels created. |
| O-36 | #55, 09-25 | F-0005 record and index row: "Raw; partly superseded by F-0006 where they conflict"; source and hash unchanged. |
| O-37 | #55, 09-25 | "remove the mention lines"; then attestation "Yes, it is my text" for the 356-line file. |
| O-38 | #55, 09-25 | Revise #55 as one goal in R11's one-mechanism form. |
| O-39 | #55, 09-25 | Reviewers for #55 only: "Plan re-review: opencode-go/grok-4.7#xhigh. Round 1: Devin gpt-6-sol-xhigh. Later rounds: opencode-go/grok-4.7#high, then opencode-go/gpt-6-luna#high if Grok cannot run. Each model differs from mine (Claude Opus 5.5)." |
| O-40 | #55, 09-25 | Go rate-limited; plan re-review on "Devin gpt-6-sol-xhigh in a new fresh session"; if Go still limited at round 2, Devin gpt-6-sol-high. |

Note: panel member B's comment on #46 used O-16 to O-27 as its own candidate question numbers; those are not decisions. The ADR uses O-19 to O-25 for its open questions, and #53 uses O-26 and O-27 for real decisions.

Routing-related decisions in order: O-3 (Claude-only tiers) → O-18 (same harness, different model; Fable/Astra before stall) → O-27 (exclude Fable, Astra, AGY) → O-29 (Kimi, Sonnet backup) → O-30 (Grok on OpenCode; Sol on Devin) → O-31 (stop and switch) → F-0006 §16 (Sol permitted again, AGY removed, Fable/Astra last-resort Judge) → O-39/O-40 (#55 only). Four reviewer bindings changed in one day; O-27 is superseded by O-30 and F-0006 without a note that says so.

## 4. The product issues that have not started

**#29 T-stfn (Step 2, core engine).** Goal: "build the LAYUP core engine with no UI and with all state in Git … a CLI plus files in the project repository". Four of six children are merged (ADR-0011; the Go skeleton and gates; `layup psb check`; CI protection). It waits on one sentence, posted 2026-09-23 19:59: "The Operator asked (2026-09-23) for a PDR first: the specification (PRD), the architecture, and an implementation plan, approved by the Operator, then the core-engine build by that plan (#42). Tasks in flight finish: #37 (`layup psb check`, needed by the PDR) and #35 (protection). No new child of this issue starts; #33 and the planned `T-b97r` and `T-tmhw` are re-scoped by the PDR's plan." Its AC 3 ("no child keeps project state outside the repository") is not affected. Nothing has happened on it for 39 hours; every later hour went to the ADR-0012 line.

**#42 T-meh2 (PDR).** Goal: `docs/pdr/PDR-0001.md` records the Operator's approval of PRD-0001, the architecture and the plan; every requirement traces to a PSB fact. Six children in order: T-zmj6 (#45, gap answers), T-ertw (#43, merged 2026-09-23 20:31), T-wjq4 (PRD-0001), T-7qvc (architecture + ADRs, each after a panel per O-4), T-55n2 (plan + traceability), T-4wrw (PDR record). Children 3–6 have no issue. The chain is serial: the PRD needs #45's answers; #45 needs the idea owner to answer 19 questions in one batch (`internal/psb/testdata/psb.tsv`, pinned by `TestGoldenRealPSB`). No answer is recorded anywhere.

**#45 T-zmj6 in detail.** Plan posted 20:31:43; corrected 20:31:54 (step 1 would skip in silence where `go` is absent). Plan review at 20:39:20 (Claude Fable 5.1, fresh session, tree `0a20ee3`): `approve-with-conditions`, budget 320 lines / 18 files, cycle cap 2, five conditions (C1 one goal class via route (a), C2 one rule for an unanswered question, C3 check `facts` fails on an absent batch with fixtures, C4 name the docs, C5 "The Operator, who is the idea owner (O-14), approves the PR that adds `F-0004`. That approval is the only evidence that the words are the idea owner's"). The author answered all five at 20:39:43 and rewrote the AC in the body. No comment since. The concrete next step is an Operator action (answer the 19 questions), not agent work.

**#33 T-vk3k** cannot start: its AC are "to be set by this task's plan (R12) when it starts, after `T-dq05` and `T-b97r`", T-b97r has no issue, and #29 re-scopes it under the PDR.

## 5. Record defects found during the audit

1. **Three #46 comments hold only the footer.** The comments at 07:48:33 (plan review 3), 08:59:21 (review round 1) and 09:23:38 (review round 3) are 344–437 bytes: a `---` and "_Posted by the author … The text above is the reviewer's output byte for byte_", with no text above. The reviewer records were never posted. Round 1's 8 findings survive only in the author's "Fixes" reply; round 3's 4 findings only in O-26's question, #53's body and `T-9tgn.md`. `git grep` on `T-7wqc` finds no copy in `runs/T-w79d/` or `docs/tasks/`. #53 plan review 1 was "supplied the #46 O-15 and round-3 excerpt" from a source that does not exist on the issue.
2. **O-15 is two decisions** (#42 PRD scope; #46 goal count). The R11 issue #49 cites the #46 one.
3. **`docs/tasks/backlog.md` on `main` holds only the template line.** No open task (#29, #33, #42, #45, #54, #55) has a backlog line or a task file on `main`; the `T-7wqc` branch adds one line for itself.
4. **F-0006 and ADR-0012 disagree at the source.** F-0006 §2 declares F-0005 incomplete and §16 reverses O-27 and the 70/90 bands; ADR-0012 at `30ca13f` is built on F-0005 alone. Neither #54 nor #55 states which of the two texts ADR-0012 must cite after #55 lands; O-35 says only "ADR-0012 resolves points 1–4".
5. **No resource record exists for the line.** None of #46, #53, #54 reached close-out, so no ADR-0007 record (model, effort, tokens, time per gate part) is in Git. `T-7wqc.md` holds only seconds per role for #54; the clause-map comment holds the one token figure (1,019,459).
