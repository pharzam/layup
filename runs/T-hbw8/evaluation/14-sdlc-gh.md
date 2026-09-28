# Evaluation: sdlc-gh

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/guilz-dev/sdlc-gh |
| Pinned commit | 1ab4b6d079002f27e93306cccf915924e4ce745e (2026-07-05, "Release 0.1.1 with npm pkg fix for bin path") |
| License | MIT (from the LICENSE file) |
| Language and needs | JavaScript (Node.js 22 or later, ES modules, no npm dependencies), bash, git. GitHub Actions and an authenticated `gh` for the GitHub part. GitHub Copilot Business or Enterprise for the coding-agent part. Optional Langfuse (docker-compose in `infra/`). No server, no database, no daemon of its own. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 20 minutes (15:31 to 15:50 local), in one session shared with AI-SDLC |
| Agent runs and cost | none (sdlc-gh starts no agent; its agents are GitHub Copilot coding agent, which this evaluation cannot run) |

## Verdict

**Borrow the pattern.** sdlc-gh is a template of CI files, a GitHub ruleset and Node scripts. Four patterns are useful to LAYUP: a stack catalog that selects one required check per stack, a ruleset applied through the API and read back by a doctor, "SKIP is a failure" in strict mode, and a hash manifest that finds drift from the template. The code itself does not fit: it is Node, it is Copilot-specific, and its rule protection has gaps that the runs below show.

## What it is (from the code)

- **Stack catalog.** `config/stacks.json` lists five stacks. Each stack has a marker file, a profile and a workflow (`go`: marker `go.mod`, workflow `product-ci-go.yml`). See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/config/stacks.json#L25.
- **Required checks per stack.** `buildMainRequiredContexts` returns `harness-static`, `diff-size`, `issue-spec-check` and `product-ci-<stack>`. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/lib/github-config.mjs#L42.
- **Ruleset.** `.github/ruleset.example.json` makes `main-protection`: a pull request is required, one approval, code-owner review, stale reviews dismissed, the four checks required. It has no `bypass_actors` field. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/ruleset.example.json#L1.
- **Applying the ruleset.** `scripts/setup-github.mjs` finds a ruleset with the same name and sends `PUT` or `POST repos/<repo>/rulesets` through `gh api`. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/setup-github.mjs#L148.
- **Reading it back.** `scripts/doctor.mjs` reads the ruleset from the API, checks `enforcement: active` and the four contexts. It does not check the code-owner rule or the bypass list. With `--strict`, a `SKIP` is a failure. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/doctor.mjs#L176.
- **CODEOWNERS.** Owners for `/.github/`, `/evals/`, `docs/telemetry-schema.md` and `docs/operations.md` only. `scripts/`, `config/`, `AGENTS.md` and `package.json` have no owner. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/CODEOWNERS#L1.
- **Harness CI.** `harness-ci.yml` runs on `pull_request`. Its required job `harness-static` runs `node scripts/validate-harness.mjs` from the PR's own checkout. `diff-size` runs `node scripts/check-diff-size.mjs` from the PR's own checkout. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/workflows/harness-ci.yml#L9 and #L105.
- **Go gate.** `product-ci-go.yml` runs `go vet ./...` and `go test ./...` with `go-version: "1.22"`. No layout, boundary or contract check. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/workflows/product-ci-go.yml#L1.
- **Diff size.** L1: 300 lines and 8 files, warn only unless `DIFF_SIZE_L1_HARD_FAIL=1`; L2: 120 lines and 4 files, hard fail. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/lib/diff-size.mjs#L4.
- **Retry limit.** `agent-retry-orchestrator.yml` counts retries in a `retry:N` PR label. It stops at `MAX_RETRIES = 3`, or when the same failure signature comes twice, and posts a comment that asks for a human. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/workflows/agent-retry-orchestrator.yml#L24.
- **Issue contract.** `check-issue-spec.mjs` with `lib/ccsd-contract.mjs` requires a linked issue with Goal, Non-goals, Constraints, Acceptance criteria and Rollback hints. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/lib/ccsd-contract.mjs#L8.
- **Hooks.** `.github/hooks/hooks.json` is a Copilot `preToolUse` hook. It blocks only `git push --force`, `rm -rf /`, `drop database` and `DROP TABLE`. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/.github/hooks/hooks.json#L1.
- **Telemetry.** `emit-telemetry-artifact.mjs` writes a JSON file that CI uploads as an Actions artifact. `cost`, `elapsed_time`, `model` and `tool_calls` are placeholders (`-1` or `n/a`) unless the caller sets them. See https://github.com/guilz-dev/sdlc-gh/blob/1ab4b6d079002f27e93306cccf915924e4ce745e/scripts/emit-telemetry-artifact.mjs#L57.
- **Drift report.** `harness-drift-report.mjs --against <repo>` hashes the harness files of the template and of a product, and marks each one that differs.
- **Identity.** The rule "cannot approve own PR" comes from GitHub's Copilot coding agent, not from sdlc-gh (`docs/auth-boundaries.md`).

## What we ran

All runs used `HOME=EVAL_WORK/sdlc-gh/home`, no `GH_TOKEN`, and a fake `gh` in `PATH` that logs each call and writes nothing to the network. The target was a Go module (`calc` package, one test, a Makefile) with a local bare remote.

1. **Wizard on the Go toy.** `node src/scripts/sdlc-gh-cli.mjs init --repo toy --stack go --codeowners @eval-owner --mode existing --skip-github --yes`.
   Result: 136 files copied (15,312 lines), `harness-ci.yml` trimmed to Go, a harness-only `package.json` added, CODEOWNERS owner set. Doctor: `PASS .harness-stack: stack is go`, `PASS CODEOWNERS: team placeholder replaced`, `SKIP GitHub labels`, then "doctor --strict reported failures".
2. **Ruleset dry run.** `node scripts/setup-github.mjs --dry-run --yes --github-repo eval/toy`. It printed the `main-protection` payload with `product-ci-go` added and no bypass actors. Note: the script calls `gh auth status` before the dry run, so a dry run also needs an authenticated `gh`.
3. **The required checks, locally.**
   - `node scripts/validate-harness.mjs` (the `harness-static` job) **failed, exit 1**: `ERROR: missing workflow for stack ts: product-ci-ts.yml`, then a crash on the missing file. The validator checks all five stacks of `config/stacks.json`, but the bootstrap copies one stack. So `harness-static` fails on every PR of a bootstrapped Go product, and the ruleset makes it required.
   - The hook, issue-spec and diff-size scenario tests passed. `go vet ./...` and `go test ./...` passed.
4. **Doctor and readiness.** `doctor.mjs --strict` exit 1 (SKIP counts); `check-l1-readiness.mjs --strict` exit 1, with "MANUAL Copilot coding agent entitlement".
5. **Rule-protection test (I3, PSB S3).** On a branch, an "agent" edited `scripts/check-diff-size.mjs` so that it never reports over-limit, and added 900 lines to `calc/calc.go`.
   - The gate as the PR runs it: `Autonomy: L1 | LOC: 903/300`, **exit 0**.
   - The same diff with the gate from `main`: `::error::Change size exceeds L1 limits`, **exit 1**.
   - Changed files against CODEOWNERS: `calc/calc.go -> none`, `scripts/check-diff-size.mjs -> none`. No code-owner review is required for the weakened gate.
   - `harness-drift-report.mjs --against toy-tamper` marked `scripts/check-diff-size.mjs | ❌ differs`. This is detection, run by hand or weekly from the template repository, not a PR check.
6. **Telemetry.** `emit-telemetry-artifact.mjs` wrote `harness-ci-pr1-run1.json` with `"cost": -1`, `"elapsed_time": -1` and `::notice::telemetry_placeholders=agent_type,model,tool_calls,cost,elapsed_time,review_outcome`.

Not run: anything on GitHub (rulesets, Actions, Copilot coding agent, gh-aw workflows). The brief does not allow repository writes, and Copilot coding agent needs a paid entitlement.

## How the rule protection would apply on GitHub (read, not run)

- `setup-github.mjs` creates the `main-protection` ruleset. With no bypass actors, no one, the owner included, can push to `main` or merge without an approving review, a code-owner review and the four green checks. A repository admin can still edit or delete the ruleset.
- The code-owner rule protects `/.github/` and `/evals/`. A change there needs a review by the owner in CODEOWNERS.
- The Copilot coding agent works on its own branch with a short-lived token and cannot approve its own PR. So when the agent is a separate identity, a person must approve every PR. That person is the real protection.
- **Under O-77** (agents act through the Operator's account with a GitHub App user token), the PR author is the Operator. GitHub does not let an author approve their own PR. The ruleset then blocks every agent PR, or the Operator must remove the rule, which is a change to a gate.
- The checks run the PR's own copy of the workflows and scripts (`pull_request` event). A PR that changes `scripts/` changes the gate that judges it. CODEOWNERS does not cover `scripts/` (run 5).
- The required context `product-ci-go` is the caller job's name. GitHub names a check from a called (reusable) workflow `<caller job> / <called job>`, here `product-ci-go / verify`. If so, the required context never reports, and the PR waits for ever. This was not tested on GitHub.

## In-Scope items S1–S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | The issue contract checks one task issue for five headings (`ccsd-contract.mjs#L8`). It does not read a problem statement or batch gap questions. |
| S2 Reproducible Discipline Setup | partly | Wizard, bootstrap and `doctor --strict` are repeatable (run 1). No baseline version pin, no evidence per value. The bootstrapped repo fails its own `harness-static` (run 3). |
| S3 Rule Protection | partly | Ruleset and CODEOWNERS protect `.github/` and `evals/` when a second human identity exists. The gate code in `scripts/` is not owned; a weakened gate passed (run 5). |
| S4 Stack-Dependent Gates | partly | `stacks.json` selects `product-ci-go`, which is only `go vet` and `go test`. No layout, boundary, contract or test-quality gate. |
| S5 Role Handoffs | partly | Triager, implementer and reviewer agent files; the issue contract is checked on each PR. No schema for other handoffs. |
| S6 Autonomous Clarification | no | No question routing. |
| S7 Verification on Every Change | partly | Four required checks on each PR. They cover diff size, issue link and `go vet`/`go test`, not layout, interfaces or a test pyramid. One of them fails in a product repo (run 3). |
| S8 Human-on-the-Loop | partly | "One human gate: PR review" (`docs/arch.md#L84`) and autonomy labels. No escalation rule for business-forking decisions. |
| S9 Stall Resolution | partly | Retry limit 3 and "same failure signature twice" stop, with a comment to a human. No fresh-context diagnosis; the record is PR labels and comments. |
| S10 Cost Visibility | no | The telemetry file has `cost: -1` and `elapsed_time: -1` (run 6), and it is an Actions artifact, not a file in the repository. |
| S11 Specification Synthesis | no | Nothing derives requirements. |
| S12 Harness-Agent Neutrality | no | Agents, hooks and instructions use Copilot formats; `AGENTS.md` is the only neutral file. No second-harness verification. |

## Invariants I1–I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the record | neutral | Rules and gates are files in the repository. Retry counts are PR labels, telemetry is an Actions artifact or Langfuse, both outside Git. |
| I2 Independent repository | supports | The gates are the target's own workflows and scripts. They run with the template absent. |
| I3 Agents cannot change rules | partly supports | Preventive for `.github/` and `evals/` only, and only with a second identity. `scripts/` is open (run 5). Admins can edit the ruleset. |
| I4 No value without evidence | conflicts | Fixed values with no cited evidence: 300 lines and 8 files, 3 retries, `go-version: "1.22"` (the toy asks for 1.23). |
| I5 Inactive check is not a pass | supports, with the strict flag | `doctor --strict` fails on SKIP (`doctor.mjs#L176`). A required check that does not report blocks the merge. Without `--strict`, SKIP passes, and the L1 diff-size limit only warns by default. |
| I6 Deterministic first | supports | "Walls are declarative and deterministic" (`docs/arch.md#L76`); all gates are scripts. |
| I7 Stack adds, never removes | neutral | The stack adds `product-ci-<stack>`. No mechanism keeps the baseline checks fixed. |
| I8 Pinned baseline | neutral | No template version is recorded in the product. The drift report can compare with a template checkout. |
| I9 Replaceable harness | conflicts | The agent integration needs Copilot Business or Enterprise, and the hooks and agent files are Copilot formats. |

## Deep-check findings it answers

- **Sol-1 and Fable-M2 (detection, not protection).** Partly. The pattern is a preventive ruleset on the default branch with no bypass actors, applied by API and read back by a doctor. It does not give the "only the orchestrator App bypasses" part (Fable-M1, M2), and it does not protect `scripts/`.
- **Sol-26 (the gate lives in the target).** Answered as a pattern: the stack gate is the target's own workflow file, selected at setup, and it runs without the template (O-76).
- **Fable-M9 (the target's own CI rejects every agent PR).** Not answered, but confirmed as a real risk: sdlc-gh's own required `harness-static` fails on a bootstrapped Go product (run 3). LAYUP must run the full required check set on a fresh target before it turns the checks on.
- **Sol-27 (PR creation and review step).** Partly: `issue-spec-check` refuses a PR with no linked issue contract. The PR lifecycle itself belongs to Copilot coding agent.
- **Stall trigger (no finding number).** The retry limit and the "same signature twice" rule are a pattern for the stall trigger (O-82).
- None for Fable-M3, Sol-5 (identity comes from Copilot, not from sdlc-gh), S12, Fable-M19, Fable-M23, Sol-2 (checks start when the PR opens, so a person can read a failing PR), or Sol-30.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Stack catalog (`config/stacks.json`) and one `product-ci-<stack>` check | borrow the pattern | A data file in LAYUP maps a stack to the native commands and the required context name that setup writes into the target. | Clear and small; run 1 showed it. Verify the context name that GitHub reports. |
| Ruleset payload, applied by API and read back | borrow the pattern | LAYUP's setup writes the ruleset with the orchestrator App as the only bypass actor, then reads it back and fails closed. | The read-back idea is good; the payload lacks bypass actors and the doctor skips the code-owner rule. |
| CODEOWNERS on rule paths | borrow the pattern, corrected | Own every path whose code a required check runs (`.github/`, the gate scripts, the ADR and guardrail files). | The template leaves `scripts/` open (run 5). Under O-77 it needs a second identity. |
| `doctor --strict`: SKIP is a failure | borrow the pattern | LAYUP's setup check treats an unverifiable item as a failure. | Matches I5 (run 4). |
| Retry limit and same-signature stop | borrow the pattern | A parameter N and a failure signature feed the stall trigger (O-82). | The mechanism is simple; the value 3 has no evidence. |
| Drift manifest (`harness-drift-report.mjs`) | borrow the pattern | Hash the pinned Armature files and compare in CI, as a complement to prevention. | It found the changed gate (run 5). |
| Diff-size gate | borrow the pattern | A native target check with limits set from evidence. | Limits are guesses (I4). |
| Telemetry artifact | reject | — | Placeholders only, stored outside Git (run 6). |
| Copilot agents, hooks, gh-aw workflows | reject | — | Copilot-only (I9); the hook does not protect rule files. |
| The scripts as code | reject | — | Node in a Go target; LAYUP is Go standard library only (F-0004 fact 1). MIT would allow vendoring, but nothing here needs it. |

## Where the searchers were wrong or incomplete

- Searcher A, S3 "partly ... only after real owner/ruleset setup": correct, but incomplete. Even after setup, CODEOWNERS leaves `scripts/` open, and the required checks run the PR's own scripts. A weakened gate passed (run 5).
- Searcher A, S10 "scores, cost, traces on PR": the code writes `cost: -1` unless a caller supplies it, and stores the file as an Actions artifact (run 6).
- Searcher A, S7 "required checks": A did not see that `harness-static` fails on a bootstrapped product repo (run 3), or that the required context name may not match the reusable-workflow check name.
- Searcher A, S9 "CI retry max 3": correct. The same-signature stop is also there.
- Searcher A, "last push 2026-09-05": correct for the API field `pushed_at`. The last commit on the default branch is 2026-07-05.
- Searcher A, I6 "Y": agreed.

## Limits of this evaluation

- Nothing ran on GitHub: no ruleset, no Actions run, no Copilot coding agent. The ruleset behaviour and the check-name mismatch are read from the code and GitHub's documented rules, not observed.
- The tamper test ran the gate script locally in the way CI runs it; it did not open a PR.
- The Copilot and gh-aw workflows (`nightly-harness-review`, `weekly-redteam`) were not read in full.
