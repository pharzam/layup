## Panel member B — Git-native state, the forge and CI

Scope: this comment covers where state lives in Git, what the forge (GitHub) enforces, and what CI runs. I read ADR-0011, the PRD §4 and §11, F-0001 §1 and the terms, and `.github/workflows/ci.yml`. I did not select an option. Each argument has a test that can show it wrong ("Falsify:").

### Q1 — The runner for stack gates

**1A. A LAYUP-owned GitHub App that sets a required check run on the target.**
- **Mechanism:** A GitHub App is installed on the target. On each `pull_request` event it runs `layup gate` on a LAYUP-owned runner, against a checkout of the PR head and the gate code at a pinned LAYUP commit. It posts one check run for each gate kind (`layup/layout`, `layup/interfaces`, `layup/pyramid`). A `pass` gives `success`, a `fail` gives `failure`, and a `not-active` gives `failure`, because NFR-004 does not allow a neutral result. Branch protection on the target lists each check name as required. If no check run was reported, GitHub keeps the PR at "expected — waiting", so a gate that did not run blocks the merge (NFR-004).
- **For:** Nothing from LAYUP goes into the target tree (O-10, O-11). A required check with no report blocks the merge because of how the forge works, not because of a rule we wrote. Falsify: open a PR on a pilot target with the App stopped. If the PR can merge, this argument is wrong.
- **For:** The gate code comes from LAYUP, not from the PR, so a PR cannot weaken the gate that judges it. This closes the #84 hole *by construction*: today the restore step only reduces that hole. Falsify: a PR that edits any file in the target and turns a `fail` into a `pass`.
- **Against:** The branch-protection setting lives in the forge, not in Git (Invariant 1). A human cannot read it in the tree, and a different forge loses it (Invariant 2). Mitigation: `layup setup verify` compares the live protection with a `docs/setup/branch-protection.json` in the target. That file is a kit record, as in LAYUP today. Falsify: change the protection by hand. If `setup verify` does not go red, the mitigation fails.
- **Against:** The App needs a hosted service with a webhook endpoint. That strains ADR-0011 decision 1 ("no daemon, no network service") for the runner, though not for the engine. Falsify: the stated boundary holds only if the engine binary makes no network call, and the App is only a thin trigger.
- **Serves:** #3, #5, #6, #7. **Strains:** #1, #2 (the protection is forge state).
- **Main risk:** App credentials. A leaked App key can post `success` on any gate.

**1B. A workflow in LAYUP's own repository, started by `repository_dispatch` or a schedule, that posts a commit status to the target.**
- **Mechanism:** A workflow in the LAYUP repository gets the target PR's head SHA, runs `layup gate` from LAYUP's own checkout, and posts a commit status (`layup/<gate>`) to the target with a fine-grained token. Protection on the target requires those status contexts. The verdict row goes into `runs/<task>/gate-results.tsv` in the target through the PR's own branch, or into LAYUP's `runs/` (the Operator decides which).
- **For:** It uses only Actions and a token, with no hosted service. That fits ADR-0011 decision 1 best. Falsify: find one part that needs a server that is always on.
- **Against:** A commit status can be posted by *any* token with `repo:status` scope. The agents push with the same account (O-9), so an agent can post `success` itself. Check runs from 1A are tied to the App identity, but a commit status is not. Falsify: with the agent's token, run `gh api .../statuses/<sha> -f state=success -f context=layup/layout`. If protection refuses it, this argument is wrong.
- **Against:** The start is not tied to the event. Without a dispatch from the target, which means a LAYUP file in the target, the runner must poll, so there is a gap. Falsify: measure the time from PR open to status. The required status still blocks the merge while it waits, so the gap costs time but is not a hole.
- **Serves:** #1 (partly), #6, #7, #9. **Strains:** #3 (the shared account can forge a status), #2.
- **Main risk:** The status is forged under the shared account until O-9 is decided.

**1C. The engine runs as a local pre-merge step, and the merge goes only through `layup merge`.**
- **Mechanism:** The harness agent calls `layup gate <target> <pr>`, and then `layup merge` merges only on all-`pass`. The verdict is written as a signed row in the target's `runs/`.
- **For:** It needs no forge feature and no network, and it is fully Git-native (#1, #2).
- **Against:** Nothing stops a merge through the forge UI or `gh pr merge`. REQ-007 says "no PR with a failed gate … merges". This option cannot meet that without forge enforcement. Falsify: merge a PR with a `fail` row through `gh`. If it succeeds, REQ-007 fails.
- **Serves:** #1, #2, #9. **Strains:** #3, #5.
- **Main risk:** Enforcement is advisory only.

| | 1A App check run | 1B Actions + status | 1C local merge |
|---|---|---|---|
| Blocks a merge with no report (NFR-004) | yes (required check) | yes (required status) | no |
| Can the shared account forge a verdict? | no (App identity) | **yes** | yes |
| Hosted service needed | yes | no | no |
| Writes LAYUP code into the target | no | no (unless there is a dispatch hook) | no |

### Q2 — The rule-protection control

**2A. CODEOWNERS plus "require code-owner review" with a second, human-only account (the ADR-0011 named control).**
- **Mechanism:** The target has `.github/CODEOWNERS`, a kit record, that maps each rule path to `@operator-human`. Protection requires a code-owner approval and dismisses stale approvals. The agents use a separate account or an App. The author cannot approve their own PR.
- **For:** The forge enforces it, and CODEOWNERS is a file in Git (#1). Falsify: an agent PR that touches `docs/guardrails.md` merges with no approval from `@operator-human`.
- **Against:** It needs O-9: two accounts. With one account the author is the owner, and GitHub does not let you approve your own PR, so *every* rule-path PR blocks, including the Operator's own. Falsify: none, because this follows from how the forge works. The cost is a second account, which is a decision for the Operator.
- **Against:** An admin bypass ("do not enforce for admins" set off) opens it again. Falsify: `setup verify` must read `enforce_admins=true`.
- **Serves:** #3, #1. **Strains:** #2 (forge setting). **Main risk:** the agents get the human account's token.

**2B. The rule paths are judged from outside, by the App in 1A: a `layup/rule-guard` required check that fails on any diff to a rule path unless a signed approval exists.**
- **Mechanism:** The check lists the changed paths against a rule-path list pinned in LAYUP, not in the target. A match fails unless the PR head has a commit or tag signed with an Operator key that is not in the agents' environment (an SSH or GPG key, checked with `git verify-commit`).
- **For:** The key separates the human from the agents, not the forge account. So it works with one GitHub account (O-9 as it is today). Falsify: an agent that has only the harness environment makes a rule-path PR go green.
- **For:** The rule-path list is outside the target, so a PR cannot shrink it. Falsify: a PR that edits the list and passes.
- **Against:** It depends on 1A, and on the key staying off the agent's machine. If the Operator and the agents share one machine and one ssh-agent, the key is exposed. Falsify: run `ssh-add -l` in the harness session.
- **Serves:** #3, #6, #9. **Strains:** #2 (the guard lives in LAYUP). **Main risk:** key custody on a shared workstation.

**2C. Detect, do not prevent: a post-merge audit (the ADR-0011 complement).**
- **Mechanism:** `layup audit rules` lists each commit on `main` that touches a rule path, with its author, signer and PR. It writes `runs/audit/rule-changes.tsv`, and a non-zero count of unapproved changes opens an issue.
- **For:** It needs nothing new, and it gives the REQ-003 metric "audit shows zero agent writes". Falsify: a seeded agent rule commit that the audit misses.
- **Against:** It does not meet "is refused" (REQ-003 acceptance). Falsify: the seeded change merges, so the requirement fails.
- **Serves:** #1, #2, #6. **Strains:** #3. **Main risk:** it is taken as a control when it is only a detector.

| | 2A CODEOWNERS + 2nd account | 2B signed-approval check | 2C audit |
|---|---|---|---|
| Works with one account (O-9 today) | no | yes | yes |
| Refuses before merge (REQ-003) | yes | yes | no |
| Depends on the Q1 runner | no | yes (1A) | no |

### Q3 — The role model and the handoff schema

**3A. Roles from the kit's gate steps; one TSV event table.**
- **Roles:** `planner` (plan and plan review), `implementer`, `reviewer`, `examiner` (stall), `operator` (human).
- **Record:** `runs/<task>/handoffs.tsv`, append-only, with the columns `ts, task, from_role, to_role, kind, artifact_path, artifact_sha, harness, session, verdict`. `layup handoff check` validates the header, the `kind` against an enum from the kit (plan, review-record, verdict…), that `artifact_path` exists at `artifact_sha`, and the `from→to` pair against an allowed-transition table.
- **For:** It is the ADR-0011 decision 2 shape, so the table can be read with no tool. Falsify: REQ-005 needs "a schema that accepts everything does not pass". A fixture row with a legal header but a bad transition must fail.
- **Against:** A TSV row only points at the artifact, so it does not validate the artifact's content. Falsify: a row that points at an empty review record passes. That means 3A needs the kit linter for each kind (`review-record-lint.sh` already exists).
- **Serves:** #1, #2, #6, #9. **Strains:** none large. **Main risk:** the transition table becomes a rule path that is not protected (Q2).

**3B. A Markdown record for each handoff, with a front-matter schema.**
- **Record:** `runs/<task>/handoff-NN.md`, with fixed headings like the ADR and PRD records, linted by a `handoff-lint.sh` in the kit style.
- **For:** A human reads it directly, and the lint pattern already exists (adr-lint, prd-lint). Falsify: the kit copy of the lint must run in a target without LAYUP (#2). If it needs `layup`, this argument is wrong.
- **Against:** It is a new linter in the *kit copy*. That is kit content that LAYUP adds, and it strains O-10's "nothing beyond the kit". Falsify: the Operator reads O-10 as allowing kit-shape additions.
- **Serves:** #1, #2. **Strains:** O-10. **Main risk:** drift between the sh linter and the Go engine (ADR-0011 Consequences).

**3C. The forge as carrier: a handoff is a PR review or a comment with a fixed block, mirrored to Git.**
- **Against:** A record that lives only on the forge breaks #1 until it is mirrored. Falsify: delete the forge and check whether the handoff history survives.
- **Serves:** the Operator's visibility. **Strains:** #1, #2. **Main risk:** the mirror lags.

| | 3A TSV | 3B Markdown + lint | 3C forge |
|---|---|---|---|
| Git is the record | yes | yes | only if mirrored |
| Validates content, not only structure | with per-kind linters | yes | weak |
| Adds to the kit copy | no (record only) | yes (linter) | no |

### Q4 — The escalation rule

**4A. A declared-scope diff.**
- **Mechanism:** The rule reads the approved-intent files (the PSB fact, the PRD scope and priorities, and any budget fields) at the last Operator-signed commit. It escalates if a PR diff or a decision record touches them, or touches a `business_forking: yes` field. The result is a row in `runs/<task>/escalations.tsv` (`ts, task, trigger_path, rule_id, state=waiting-operator`) and a required check `layup/escalation` that stays `failure` until an Operator-signed resolution row exists.
- **For:** It is deterministic (#6), and it gives #26's "approved intent" a mechanical meaning: a change to the approved files. Falsify: a seeded scope-against-date trade-off made only in code, with no change to the intent files, is not caught.
- **Against:** Budget and legal changes often appear in no file. Falsify, as above: the rule misses what is not written.
- **Serves:** #1, #3, #6. **Strains:** coverage. **Main risk:** false negatives.

**4B. The agent declares the class, and the engine enforces the stop.**
- **Mechanism:** Every decision record (the task's plan, an ADR) must carry `class: business-forking|in-intent`, with a required reason. The engine stops on `business-forking`, and a second-model reviewer checks the class.
- **For:** It covers decisions that are not written in any file. Falsify: a sample audit counts misclassified records.
- **Against:** The judgement is an LLM's (#6 is strained), and an agent can under-declare. Falsify: a seeded budget decision labelled `in-intent` passes.
- **Serves:** #1, coverage. **Strains:** #3, #6. **Main risk:** under-declaration.

4A and 4B can stack: 4A as the deterministic floor, 4B as the declared ceiling. That is an option, not a choice.

| | 4A diff on intent files | 4B declared class |
|---|---|---|
| Deterministic | yes | no |
| Catches a decision not written in a file | no | partly |
| Where recorded | `escalations.tsv` + required check | the decision record + `escalations.tsv` |

### Q5 — The stall procedure

**5A. Counted limits in Git, a required `layup/stall` check, and the examiner is a different harness.**
- **State:** A `stalled` row in `runs/<task>/stalls.tsv`. The task file status becomes `stalled`, and the PR is blocked by a failed required check.
- **Limits:** A disagreement is stalled when the review cycle cap is reached with findings still open. The cap is 1, as in bootstrap part "cycle cap 1", so the cap already exists. A no-progress stall is N equal gate-failure fingerprints (gate, rule, file) in a row, or no new commit on the branch for T hours. N and T are Operator values (Invariant 4); they must not be guessed.
- **Examiner:** A fresh session on a different harness or model than the ones that did the stalled work. The engine picks it from `handoffs.tsv`, which records who did the work, so independence is checked deterministically. This is the same test as the current "Who may review".
- **Records:** `runs/<task>/stall-NN.md` with the headings Trigger, Evidence (SHAs, rows), Diagnosis, Outcome, and Operator package. The package is the file plus the linked SHAs, reported on the issue.
- **For:** The fingerprint count is deterministic. Falsify: a seeded loop of the same fail three times must stop at N, and a loop with a changing fingerprint that makes no progress escapes it. The time limit T is the backstop for that case.
- **Serves:** #1, #2, #6, #9. **Strains:** none. **Main risk:** a loop with a changing fingerprint that is only caught by the time limit T.

**5B. Measure progress, not repeats.**
- **Mechanism:** A step makes progress if the count of failing checks or open findings goes down within K attempts. If not, the step is stalled. The rest is the same as 5A.
- **For:** It catches the changing-fingerprint loop. Falsify: a seeded oscillation (one failure fixed, a new one added) never counts as progress.
- **Against:** A task that is fixed "properly" can raise the count for a while. Falsify: a real refactor task gets flagged.
- **Serves:** #6. **Main risk:** false stalls.

**5C. The examiner is the Operator's queue, with no machine examiner.**
- **Against:** REQ-010 needs a fresh-context diagnosis *before* the Operator. Falsify: the requirement text. So this option only fits as a fallback when no independent harness is available.

| | 5A fingerprint + T | 5B progress metric | 5C Operator only |
|---|---|---|---|
| Deterministic trigger | yes | yes | n/a |
| Catches an oscillating loop | only by T | yes | n/a |
| Meets REQ-010 examiner | yes | yes | no |

### Questions for the Operator

1. **O-9:** Can there be a second account that only a human uses (option 2A), or must the control work with one account through a key that the agents do not hold (option 2B)?
2. **A hosted service:** May LAYUP run a GitHub App and a webhook endpoint (1A), or must the runner stay inside Actions and tokens (1B, which the shared account can forge until O-9 is decided)?
3. **The protection setting is forge state, not Git state.** Is a verified mirror file (`branch-protection.json` in the target) enough for Invariant 1?
4. **Where do gate-result and stall rows live:** in the target's `runs/` (the kit rules, O-13) or in LAYUP's `runs/`?
5. **Numbers:** the retry count N, the no-progress time T, and the cycle cap for a disagreement. These are values for you to set (Invariant 4).
6. **Approved intent:** Which files hold the approved intent and the budget, so that option 4A can read them?
