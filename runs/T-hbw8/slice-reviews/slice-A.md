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
