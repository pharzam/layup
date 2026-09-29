# Slice A review record

Reviewer: Claude Fable 5.1 (`claude-fable-5-1`), a fresh `claude -p` session in a disposable clone at `7080ae8`, read-only tools, no forge token, `--setting-sources project,local`. Brief: the common slice-review brief of plan v2 plus the slice A part. Not a gate round (plan-review note 10).

## Round 1

2026-09-29T06:42:11Z to 06:47:25Z (5 min 14 s); 23,809 output tokens; USD 3.68 at list price. The output follows word for word; its line numbers point into `7080ae8`.

~~~~text
# Slice A review — round 1

Reviewer: claude-fable-5-1. Head: `7080ae8`. Verdict: `material findings: 8`.

## Material findings

### M1. A linked work tree lets a session write the records branch and the hooks
- Where: `docs/architecture.md:153` "`work/`: a Git work tree of LAYUP's clone of the target", and `:83` "Only `layup run` commits to the records branch". ADR-0014 Consequences: "An agent session cannot write a record or a decision with any token it holds."
- Source: K01, FT3. `git worktree` shares refs, config and hooks with the main clone.
- Why it matters: A session can move `refs/heads/layup-records` without any token (`git update-ref`). `layup run` then pushes that commit as the App. A session can also set `core.hooksPath` in the shared config, and `npm install` with husky does this by accident. The hook then runs inside `layup run`'s own commits and pushes, and it can see a token passed with `git -c` through `GIT_CONFIG_PARAMETERS`. L-A1 covers a session that searches the host. It does not cover a session that only uses its own work tree.
- Fix: Give each session its own clone, not a linked work tree. Keep the records in a clone that no session can reach. Push only by explicit SHA (`<sha>:refs/heads/...`), and check that the task branch descends from the base commit. Run LAYUP's own `git` with hooks turned off.

### M2. The token that pushes agent branches can change CI workflows
- Where: `architecture.md:98-99` "and workflows (write), because the setup writes the target's CI job", and `:173-174` "`layup run` pushes the session's branch with the App's token".
- Source: Invariant 3 (`F-0001#3`); FT4. GitHub refuses a push that changes `.github/workflows/*` from an App token without `workflows`.
- Why it matters: A session commits an edited workflow, and LAYUP pushes it. A `push`-triggered workflow then runs the edited copy on the task branch, with the repository's secrets. Without this permission, GitHub would refuse the push for free.
- Fix: Ask for installation tokens scoped by request (the `permissions` body of the access-token call). Use `workflows` only for the setup push. Or refuse any session commit under `.github/` before the push.

### M3. The session gets rules from outside the work tree (K04 only half answered)
- Where: `architecture.md:166-167` "The session reads its rules only from the work tree".
- Source: K04, Fable-M23 (its fix asked for a probe). code.claude.com/docs/en/memory says two things. First, Claude Code loads `CLAUDE.md` "from your current working directory and every directory above it". Second, the managed file `/Library/Application Support/ClaudeCode/CLAUDE.md` "cannot be excluded".
- Why it matters: An empty `HOME` does not remove a `CLAUDE.md` in an ancestor of the LAYUP work area, or a managed policy file. Claude sessions then get rules that Codex sessions do not, and S12 fails in silence.
- Fix: Put session directories under a path with no ancestor rule files, and check this in code at each start. Add a per-harness probe that records the instruction files it loaded. Or record the managed-policy case as a known limit.

### M4. The "effective rules" read cannot see who may bypass
- Where: `architecture.md:120-123` "reads the effective rules of the default branch and of the records branch … and stops … when they are not the expected ones"; ADR-0014 decision 6.
- Source: docs.github.com/en/rest/repos/rules. "Get rules for a branch" returns rules only. On "Get a repository ruleset", "the bypass_actors property is only returned if the user making the API request has write access to the ruleset". The App has "no administration permission" (`:99-100`).
- Why it matters: The key property is "only the LAYUP App bypasses" (selection-v2 row 0017). An added bypass actor, such as the admin role or a deploy key, passes the start check. "The forge must enforce this, or LAYUP stops" is then not true.
- Fix: Say which token reads `bypass_actors`, and when (for example the Operator's, at setup, copied into the records). Or record this as a known limit, with the audit as detection after the fact.

### M5. `layup audit` fails on every correctly set-up target
- Where: `architecture.md:124-126` "fails when … an update of the default branch was not a merge of a pull request by the App".
- Source: The GitHub activity API types are `push`, `branch_creation`, `pr_merge` and others. A pull request needs a base branch that already exists.
- Why it matters: The first commit of the default branch, and the setup output, are a `branch_creation` or `push`, not a `pr_merge`. So the audit's documented result (pass on a clean target) differs from its actual exit.
- Fix: Scope the rule from the setup commit that the records name, or list the allowed setup updates by SHA.

### M6. A "live lease" is not defined
- Where: `architecture.md:60-62` "A second run that finds a live lease … stops"; ADR-0013 decision 1.
- Source: FT6, L-A3. The evaluation pattern "a lease per session with heartbeat and reclaim after a time-out" (summary, Claims: Paperclip, Beads) was not taken.
- Why it matters: A lease row holds a run ID, a host and a start time, and gets an end time only on a clean stop (W-12 step 12). After a crash, one of two things happens. Either every later run stops for ever, or code guesses and two hosts run at once.
- Fix: Add a heartbeat or expiry field and the rule that reads it. Or name the Operator's takeover command, and point to section 11 with `later: slice F`.

### M7. The installation token changes O-77 and O-92, and nothing records the change
- Where: `architecture.md:96-97` "acts with its installation token, as the App's bot"; ADR-0014 "reads O-77 this way" (it names only the session credentials).
- Source: O-77 (approved by O-84): "through a GitHub App's user access token … The Operator stays the responsible actor". O-92 (3): "goes through the LAYUP GitHub App's user access token".
- Why it matters: The actor becomes `layup-agent[bot]`, not the Operator's account with a badge. Author IDs and the Operator's responsibility change. The Operator decided the other token.
- Fix: State the change and the reason. A user token carries the Operator's rights; the bot is a separate actor with no human rights. Get the Operator's confirmation on #72 (R13), or use the user access token.

### M8. Section 14 breaks its own coverage rule
- Where: `architecture.md:185-186` "a row with no walkthrough fails". Rows `:195` (`NFR-007`) and `:198` (#69 forge question) have "—".
- Source: runs/T-hbw8/root-cause-missed-solution.md fix 5 (the rule cited at `:186`).
- Fix: State the exception (a named build check, or a known limit). Name the exact command and the rule `go list -deps` applies (only standard-library and module packages). Or give each row a W-12 step.

## Notes
- N1. `architecture.md:52` "each reads files": `layup audit` reads the forge API, not files. Say so.
- N2. Step 10 should say that the base and records commits come from the session start row, not from the result file that the model writes.
- N3. The register's "session or resume flag" (`:161-162`) resumes from harness state in `home/`, which a clone does not carry (FT6). Step 9 says "no memory of H1". State when a resume is allowed.
- N4. W-12 step 7 parses a human's free-text comment ("give `T-7` to H2") and is tagged `code`. Slice F must set a fixed form, or not tag it `code`.
- N5. `:100-101` "the App can be `layup-agent`": as O-92 records it, that App has no workflows permission and is installed on `pharzam/layup` only. List the changes, which the Operator must approve.
- N6. After LAYUP stops (Part 2), no human can push the records branch under its ruleset. Say that the records freeze, or say how a human adds to them.
- N7. The answer to ADR-0011's "outside the tree a reviewer reads" is a link in a PR body, which is forge text. This is acceptable, but it is weak.
- N8. L-A1 should name the likely route. A session that uses the Operator's `gh` login posts a comment with an empty `performed_via_github_app` (`gh` is an OAuth app), and that comment passes as a human decision.

## Checklist rows
- K01: not answered in full (M1). K02: answered. K03: answered, with L-A1; see M2. K04: not answered in full (M3). K05, K06, K07: answered.
- P04, P09, P21: answered. P13: known limit L-A2 (acceptable).
- S12: answered, subject to M3. R07, R08, R11: answered. R13: answered by a build check, but the coverage rule is broken (M8).
- I1, I2: answered. I9 (in part): answered, subject to M3.
- C2: answered. C6: answered (dead-man job later). C7 (in part): answered.
- D05, D06: answered. D07 (in part): answered. D08: not answered (M7). D17: answered, with L-A1.
- FT3: answered, subject to M1. FT5: answered. FT6: answered, subject to M6 and N3.
- L-A1, L-A2, L-A3: acceptable. L-A1 does not cover M1.

## Existing solutions
- GitHub Copilot coding agent: an ephemeral sandbox, pushes only to its own branch prefix, and workflows wait for a human's approval. This is the pattern for M1, M2 and L-A1.
- Dagger container-use and dev containers: one container and one branch per agent, which closes M1 and M3.
- Spec Kitty's coordination branch (ev06) is a precedent for `layup-records`. Paperclip and Beads leases with heartbeat and reclaim answer M6, and the design ignores them.~~~~

### The author's answer to round 1

All eight material findings are accepted and fixed in one commit; seven of the eight notes too.

| Finding | Fix |
| ------- | --- |
| M1 | Each session gets its own clone, not a linked work tree; LAYUP fetches the session branch by SHA with hooks off, checks its ancestry, and pushes that SHA (§4; ADR-0015 d4, d5; W-12 steps 2, 8, 10). |
| M2, M7 | LAYUP now acts with the App's **user access token**, as O-77 and O-92 decide, not an installation token. The App keeps the O-92 permissions plus commit statuses: no workflows, so the forge refuses a workflow change that LAYUP pushes; the approver merges a batch that changes `.github/workflows/` (O-93). Records-branch bypass: the admin role (§3; ADR-0014 d3). |
| M3 | Code refuses a harness rule file above the session and records system-wide policy files; the probe records loaded instruction files; the rest is known limit L-A5 (§4). |
| M4 | The effective-rules read plus a push probe with the App's token; the bypass list is recorded from the Operator's apply command at setup; a later change is known limit L-A4 (§3). |
| M5 | The audit starts from the setup commits named in the records, allowed by SHA (§3). |
| M6 | The lease has a heartbeat renewed every `lease.H` and a takeover after `3 × lease.H` (§2; ADR-0013 d1). |
| M8 | The coverage rule allows a named check or a known limit; `NFR-007` names its command, the forge row names L-A2 (§14). |
| N1, N2, N3, N5, N6, N8 | Applied (§2 audit reads the API; bases from the start row; no resume across attempts; the App changes listed; the records freeze; L-A1 names the `gh` route). |
| N4 | The comment form is fixed by slice F; W-12 step 5 says so. |
| N7 | Declined: the link is a convenience; the records branch itself is in every clone. |

## Round 2

2026-09-29, head `a72bf49`; 6 min 6 s; 26,952 output tokens; USD 2.90 at list price. Same reviewer and setup. Word for word:

~~~~text
# Slice A review — round 2
Reviewer: claude-fable-5-1. Head: `a72bf49`. Verdict: `material findings: 4`.

Round 1 status:
- **Closed:** M1 (separate clones, fetch by SHA with hooks off), M3 (ancestor check plus L-A5), M5 (the audit counts from the setup SHAs), M7 (user access token, as O-77 and O-92 decide) and M8.
- **Closed, with a new defect:**
  - M2 is closed by the forge refusing workflow pushes, but the refusal path is missing (new M4).
  - M4 is closed as L-A4, but the new probe writes to the default branch (new M1).
  - M6 is closed for lease writes, but not for other records pushes (new M2).

## Material findings

### M1. A probe that the forge accepts puts a commit on the default branch
- Where: `docs/architecture.md:135-137` "pushes an empty probe commit to the default branch with the App's token and stops unless the forge refuses it". Also ADR-0014 decision 6.
- Source: FT5 and `architecture.md:27` "Only three things go into a target". An accepted probe is none of the three. It also breaks §3's own reason at `:82-84`: a commit on the default branch "would put every open pull request out of date". And it breaks the Safety limits: an irreversible write with no authorization.
- Why it matters: the Operator adds the admin role as a bypass actor. At the next start, the probe lands on `main`. Push-triggered CI runs on it, and every open pull request goes out of date. It cannot be removed without a force-push. The check does its damage in exactly the case it exists to detect.
- Fix: make the default-branch ruleset of section 6 also target a probe ref, for example `refs/heads/layup-probe`. Bypass is per ruleset, so pushing to that ref proves the same thing, and step 1 already reads that ref's effective rules. Otherwise, state what happens after an accepted probe: record it, and have the Operator revert it through a pull request.

### M2. A run whose lease was taken over keeps writing
- Where: `architecture.md:65-66` "Each lease write is a push that is not forced, so of two runs that write at the same time, one is refused and stops". ADR-0013 d1 has the same gap: "A second run that … or whose push is refused, stops".
- Source: `architecture.md:60` "One run per target", FT6, and the lease pattern the design cites (Paperclip, Beads).
- Why it matters: host A sleeps (laptop lid closed) for longer than `3 × lease.H`, and B takes over. A wakes in the middle of step 10. Its result push is refused, but that is not a lease write. The rule only covers the second run, and records are append-only tables, so a natural implementation fetches, rebases and pushes again. A and B then both push task branches, post comments and merge.
- Fix: "Any refused push of the records branch stops the run; it may continue only after it re-reads the lease and still holds it." Also state that each forge write comes after the records push that announces it. Step 6 and step 10 already follow this order, so it acts as a fencing token.

### M3. "No longer current for its task" is undefined, and the heartbeat makes one reading refuse every result
- Where: `architecture.md:95` "when they are no longer the current ones for the task, the result is refused". Also ADR-0014 d2 "a result whose commits are no longer current for its task is refused", and W-12 step 10 "checks that they are still current for `T-7`".
- Source: Bootstrap rule 3 (operative ambiguity). It interacts with `architecture.md:61`, the heartbeat renewed every `lease.H`, which is a records commit.
- Why it matters: read as "the records commit is still the head of `layup-records`", the rule fails any session longer than `lease.H`. Heartbeats alone move the head, so every result is refused. Read loosely, it guards nothing. Two implementers will build two different rules.
- Fix: define "current" by a record: "no row for this task after the session start row changes its assignment, attempt or base commit", and name the table that holds such rows. Or define it as "the attempt in the start row is still the task's open attempt".

### M4. A task-branch push the forge refuses has no outcome, and the records already bind its SHA
- Where: `architecture.md:214-215` "and it pushes that SHA to `task/<task>/<attempt>` with the App's token". W-12 step 10 commits the result and "the commit SHAs bound to the session" before the push. W-12 step 8 starts from "the head of attempt 1, which `layup run` pushed".
- Source: the round-1 M2 fix, which relies on the forge refusing a push that "changes `.github/workflows/`" (the author's answer). GitHub docs: an App needs the Workflows permission to "edit Actions files in the `.github/workflows` directory". Also O-93 ("a batch that changes `.github/workflows/` is merged by its human approver") and FI1.
- Why it matters: a session edits `ci.yml`. The records say the SHA is bound to the session, the push is refused, and nothing says what happens next. Step 8 (attempt 2) then starts from a head that is not on the forge. The same limit means neither the O-93 workflows batch nor the setup's CI job (§1 item 1) can reach the forge through `layup run`, and no route is named for either.
- Fix: before the push, refuse a branch whose diff from the base touches `.github/workflows/`. Record that as a failed result with a reason, and push only after that check. Say in §3 that LAYUP never delivers a workflow change, and that the Operator pushes or opens such a branch (slice B names the step).

## Notes
- N1. The audit (`:146-148`) accepts "a push … by the Operator's account" on the default branch, but its ruleset allows no push. On that branch, accept only `pr_merge` and `merge_queue_merge` (plus the setup SHAs). A `push` there means a bypass, which closes part of L-A4. The docs say `actor` "can be null": fail on it (FT1).
- N2. O-93 says "its human approver". If the idea owner is that approver, the audit's "by the Operator's account" fails a legitimate merge. Make the two agree (for example, "an account in `approvers.tsv`").
- N3. The lease compares a heartbeat time from host A's clock with host B's clock. Kubernetes leader election avoids clock skew by timing locally from when it observes a change. Adopt that, or state a bound on the skew.
- N4. A user access token refresh makes "that refresh token and the old user access token" unusable (GitHub docs). A takeover on another host needs a new device-flow authorization. Say so under L-A3.
- N5. Say that `layup-records` is an orphan branch. If it forks from `main`, it carries `.github/workflows`, and every heartbeat push runs the target's push-triggered CI.
- N6. "The one variable that carries the harness's own credential" does not fit Codex with a ChatGPT login. That login lives in `auth.json` under `CODEX_HOME`, which also loads a global `AGENTS.md`. The register row needs a file route, and the rule-file check must cover it.
- N7. Make each session clone with `--no-local` (or `--no-hardlinks`). Otherwise its object files are hard links into LAYUP's clone (this belongs to the L-A1 class).
- N8. The rule-suites API (`result` = `bypass`) could detect a later bypass after the fact under L-A4. Check which permission it needs before relying on it.

## Checklist rows
- K01: answered, subject to M2 and M3.
- K02: answered.
- K03: answered, with L-A1.
- K04: answered, with L-A5.
- K05, K06, K07: answered.
- P04, P09, P21: answered.
- P13: known limit L-A2 (acceptable).
- S12: answered, subject to N6.
- R07, R11: answered.
- R08: answered (step 13 is later: slice B).
- R13: answered by its `PRD-0001` build check.
- I1, I2: answered.
- I9 (in part): answered, subject to N6.
- C2: answered, subject to M2.
- C6, C7 (in part): answered.
- D05, D06, D08, D17: answered.
- D07 (in part): answered, subject to M1 and M4.
- FT3: answered.
- FT5: answered, subject to M1.
- FT6: answered, subject to M2.
- L-A1, L-A2, L-A3, L-A5: acceptable.
- L-A4: acceptable, but N1 narrows it for free.

## Existing solutions
- Kubernetes Lease and client-go leader election use a renew time, a lease duration and locally observed time. Kleppmann's fencing tokens are the standard answer to M2.
- The GitHub Copilot coding agent pushes only to its own branches, and workflow runs wait for a human's approval. This is a precedent for M4's route.
- Spec Kitty's coordination branch (ev06) is a precedent for `layup-records`. GitHub's rule-suites API is an existing detector for bypasses (N8).~~~~

### The author's answer to round 2

All four material findings are fixed in the next commit; notes N1 to N7 are applied, N8 is recorded under L-A4. This was the second round, the limit of plan v2, so the fixes of round 2 are not reviewed in this slice; the author asks the Operator on #72 how to review them (plan v2: "A material finding that stays open comes to you before the next slice").

| Finding | Fix |
| ------- | --- |
| M1 | The probe pushes to the ref `layup-probe`, which the default branch's ruleset also covers (§3; ADR-0014 d6; ADR-0017 d3). |
| M2 | Fencing: every forge write comes after the records push that announces it; any refused records push stops the run until it re-reads the lease and still holds it (§2; ADR-0013 d1). |
| M3 | "Current" is defined by a record: the attempt of the session start row is still the task's open attempt (§3; ADR-0014 d2; W-12 step 10). |
| M4 | A branch that touches `.github/workflows/` is refused before the push; its change becomes a proposal for an approved batch that the human approver pushes and merges (O-93); a SHA is bound only after the forge accepted the push (§4; ADR-0015 d5). |
| N1, N2 | The audit accepts only pull-request merges by an account in `approvers.tsv` on the default branch; a push or an empty actor fails. |
| N3 | The lease is timed by the watcher's own clock from the last change it saw. |
| N4 | L-A3: a takeover on another host needs a new authorization of the App. |
| N5 | `layup-records` is an orphan branch. |
| N6 | A harness credential can be a file that the register row names, copied into `home/`. |
| N7 | Session clones are made with `--no-local`. |
| N8 | Recorded under L-A4 as a possible narrowing, not relied on. |
