# Report: ucsandman/Agnostic-AI

The report of one read-only agent of task `T-w73g` (#133), word for word below
the line. Model: Claude Opus 5.5, an Explore subagent on Claude Code, 2026-10-07,
20 tool calls, 100.6 s. Repository: <https://github.com/ucsandman/Agnostic-AI> at
`aab09ff0d9533967045e7f2b04e478ec4e05528e`, a shallow clone in a scratch
directory outside the checkout. No code from the repository was run.

---

I read the repo as untrusted data and ran none of its code. **License: MIT** (`/private/tmp/claude-501/-Users-farzam-projects-layup/6f4ece82-2504-4e0e-b120-e779fedc77da/scratchpad/repos/Agnostic-AI/LICENSE`, author Wes Sander). It is plain Node 18+ CommonJS with zero npm dependencies.

**What it is:** a personal harness for Claude Code: guard hooks, rules, subagents and workflow scripts. A capture/apply engine ports it to about 20 other AI coding clients. It is not a project orchestrator: there is no phase loop, no milestones and no forge.

All paths below are relative to that repo root.

| Item | Path | LAYUP concern | Verdict | Why |
|---|---|---|---|---|
| prove (mutation proof) | `tools/prove/prove.cjs` | Quality gate | Copy pattern | Runs a check, breaks the file it watches, expects red, restores, expects green. The restore is guaranteed. About 200 lines, easy to redo in Go. |
| task-contract | `tools/task-contract/{task-contract,parse}.mjs` | Acceptance, screening | Copy pattern | Exit codes: 0 pass, 1 fail, 2 `insufficient_spec` goes to a human, 3 malformed. A check that cannot run counts as failed. The `non_inferable` tag is written at spec time, never asked of a model. The strict-subset parser rejects anything outside the schema. |
| Frozen evaluator set | `tools/gates/freeze.cjs`, `gate-manifest.json`, `engine/hooks/gate-freeze.cjs` | Agents cannot change the gates | Copy pattern | A content-hash lock over the gate files. A `#key` entry hashes one JSON key. Edits are denied outside approved directories. Caveat: the lock file is per-machine and gitignored, which breaks LAYUP's Git-only state. Commit it instead. |
| wiredark | `tools/wiredark/wiredark.cjs` | Implement gate | Copy pattern | A new export with no non-test caller blocks the commit. Deterministic. Regex-based and limited to JS/TS/Python, so the check would need rewriting for Go. |
| Gate runner and ratchets | `tools/gates/gates.cjs`, `budgets.json` | Quality gate | Copy pattern | One function per check. Exit 0/1/2 (2 means the runner broke). Every check prints how many items it processed. Word budgets work as an anti-growth ratchet, and a reference ratchet catches dropped pointers. |
| repeat-tool-guard | `engine/hooks/repeat-tool-guard.cjs` | Stall detection | Copy pattern | Hashes calls with sorted keys. Bookkeeping tools neither count nor reset the streak. Thresholds 3/5/8 are logged so they can be tuned. Advisory only. |
| wakeup-guard | `engine/hooks/wakeup-guard.cjs` | Stall detection | Copy pattern | From the third consecutive machine wake-up with no human turn, it tells the session to stop. |
| fix-findings workflow | `workflows/fix-findings.js` | Review rounds | Copy pattern | Groups findings by the files they touch, with no two groups sharing a file. Fixer and reviewer use different models. Refix happens once, never in a loop. Ends with one fresh whole-diff convergence read. Bans git commands that change the shared tree. |
| adversarial-review workflow | `workflows/adversarial-review.js` | Blind review | Copy pattern | Finders per dimension, then a skeptic per finding that defaults to `real=false`, with an optional 2-of-3 vote. The finders and skeptics are LLMs (fine for review, not for gates). |
| subagent-budget-guard | `engine/hooks/subagent-budget-guard.cjs`, `tools/subagent-budget/calibrate.cjs` | Cost, routing learned from records | Copy pattern | Each spawn declares `# EST: n calls`. Declared vs actual is logged and the constants are re-fitted from data. A repeated identical denial is let through once (anti-thrash). |
| Capability-graph and model guards | `engine/hooks/{capability,agent-model}-guard.cjs` | Tier routing | Skip (idea only) | Downward-only delegation and an "advisor one rung up" are good ideas. The code is hard-wired to Claude Code payloads and Claude model names. |
| Harness port engine | `engine/harness/` (`index.cjs`, `targets/`, `common.cjs`) | Replaceable harnesses | Copy pattern | Guarded writer: generated-file marker, managed regions, backup before overwrite, skip hand-edited files, `--check` exits 1 on drift. Each target declares a `supports` capability table, and anything ported to a target lacking a feature is dropped with a recorded reason. It syncs config, not sessions, so it is not a runtime router. |
| spend | `tools/spend/spend.cjs` | Cost recording | Skip | Reads Claude transcripts after the fact. Prices are hard-coded and estimated. |
| distill/harvest | `engine/distill/`, `engine/harvest/` | Retrospective learning | Copy idea, skip code | Promotion ladder with a human approval step (`--approve`). State is gitignored JSONL. The code is weak (hash dedupe, day counting). |
| Jev skill-pick experiment | `labs/claude-mods/experiments/jev/FINDINGS.md` | Smart-if provider | Read | Measured data: Jev ranks poorly over 407 options but recall is 5x better on a 22-option list. Keep LAYUP's decision-point option sets small. |
| shim, secret guards, DashClaw | `engine/hooks/` | — | Skip | Tied to a client's hook dialect. `slopsquat-guard` makes live registry calls (non-deterministic). DashClaw is the author's own remote-approval service. |

**Doesn't fit LAYUP:**
- Runtime state lives in `~/.claude` and `storage/`, not in Git.
- Assumes Claude Code throughout (Fable/Opus/Sonnet/Haiku ladder, Claude transcript format).
- Several scheduled jobs are Windows Task Scheduler / PowerShell only.
- Review workflows put an LLM in the loop.

**README claims the code does not back:**
- "a turn cannot end without its verification": `forced-verify-stop-gate.cjs` only blocks when `STOP_GATE_MODE=enforce` is set, and it regex-matches words like "npm test" in the message. It never checks whether any edits happened.
- "pre-commit runs the secrets scan, wiredark, the guard freeze and the doc gates": that hook lives at `~/.claude/git-hooks/pre-commit` and is not in the repo. CI runs neither gates nor wiredark, and skips `reg-capabilities` and `reg-context`.
- Distill "three signals across two sessions, older signals counting half, a contradiction demoting": the code counts 3 distinct days with no decay, and a conflicting rule is converted to a few-shot example rather than demoted.
- "216 skills" in `skills/`: no skill files are tracked in the repo.
- Hook counts disagree: 39 in the table, 41 in the text, 42 files on disk.

**Top 3 recommendations:**
1. **Copy the task-contract verdicts** into LAYUP's acceptance and screening: pass, fail, human-routed `insufficient_spec`, malformed. A check that cannot run fails, and the "needs a human" tag comes from the spec, never from a model's self-doubt.
2. **Copy the frozen evaluator set plus prove** for the "agents cannot change the gates" concern. Keep the hash lock committed in the target repo, and require each new gate to be proven red before it counts.
3. **Copy fix-findings' structure** for review rounds: findings grouped by the files they touch, author and reviewer on different models, one refix, and a single fresh convergence read with no loop. Also borrow repeat-tool-guard's logged thresholds as raw data for stall detection at the limit.
