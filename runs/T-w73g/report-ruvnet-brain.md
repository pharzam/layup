# Report: stuinfla/ruvnet-brain

The report of one read-only agent of task `T-w73g` (#133), word for word below
the line. Model: Claude Opus 5.5, an Explore subagent on Claude Code, 2026-10-07,
20 tool calls, 120.6 s. Repository: <https://github.com/stuinfla/ruvnet-brain> at
`b4590d469f52937550d8d89afa86a394a1d72875`, a shallow clone in a scratch
directory outside the checkout. No code from the repository was run.

---

I surveyed ruvnet-brain read-only and ran none of its code. The verdict: almost nothing is worth taking as a ready-made component, but four or five patterns map directly onto LAYUP.

**License:** MIT (Stuart Kerr / Isovision.ai), in `LICENSE`. **Secrets:** none found committed. The only tracked file in `keys/` is `ruvnet-brain-signing.pub.pem`, a public key; `.gitignore` excludes the private key. A grep for common key formats only hit fake fixtures in tests. `.gitignore` mentions a past incident where live API keys were printed by a diagnostic, but I found no key values in the tracked tree. The clone has a single commit, so I could not check git history.

**What it is:** a Node/ESM Claude Code plugin with a searchable index ("brain") over rUv's RuvNet repos (vector store, cross-encoder reranking, MCP search tool). It is wrapped in many fail-open hooks: grounding, continuation, lessons, spend, and model routing.

| Item | Path | LAYUP concern | Verdict | Why |
|---|---|---|---|---|
| Frozen eval gate | `scripts/eval-brain.mjs`, `evals/held-out.json`, `evals/baseline.json` | evals, quality gate | copy pattern | The question set is hash-pinned and a test fails if it changes. Grading uses ground truth, no model. It fails if the Wilson lower bound drops or if no baseline exists, and baselines are only written by an explicit `--record`. |
| Citation resolver | `kb/verify-citation.mjs` | deterministic checks | copy pattern | A cited path counts only if it exists in the store. No model involved, and it parses by declared length so text inside a document can't fake a citation. |
| Route qualification | `scripts/model-weekly-qualification.mjs`, `scripts/model-routing-policy-promotion.mjs`, `config/model-router/qualification-contract.json` | routing learned from records, blind review | copy pattern (strongest) | A candidate model must pass exact-answer checks first. Then a separate reviewer grades incumbent and candidate outputs as anonymized, randomly ordered A/B. Only one role changes per run, there is a deadline, and promotion is atomic only if the policy file is unchanged (checked by hash). |
| Tiered router | `config/model-router/policy.default.mjs`, `routing-policy.template.json`, `routing-eval-cases.json` | routing by model tier | copy pattern | Five classes (fast / medium / substantial / hard / exceptional). The caller's structured `taskFacts` can raise the class but never lower it. "Exceptional" needs a named reason. If the route is unavailable it returns null instead of quietly using a weaker model. The eval cases set minimum and maximum class. The free-text classifier itself is fragile keyword matching, so skip it. |
| TriSmart review | `tri-smart-skill/tri-smart/scripts/review.mjs` | independent review by a different model, harness adapters | copy pattern | Per-CLI argument recipes for `claude`, `codex` and `grok` (read-only/plan mode, no tools, JSON output). Process-group timeouts, provider API-key variables removed from the environment, a scribe chosen by hash, every reviewer verifying the same synthesis digest, one bounded revision. |
| Lesson promotion | `plugin/scripts/lesson-promote.mjs`, `lesson-gate.mjs` | retro lessons feeding the learning loop | copy pattern | A lesson is promoted only after it was independently re-taught in two or more projects. That is a count, not a similarity score or model judgement. Each lesson has an enforcement level (advisory or block) and advisories have a show cap. |
| Single refusal point | `plugin/scripts/decision-gate.mjs` | quality gate, hooks | copy pattern | N policies feed one write gate, with ordered precedence and one combined reason. |
| Continuation ledger | `plugin/scripts/continuation-gate.mjs` | stall detection | skip (idea only) | A Stop hook compares a ledger of promised work against what happened. State lives in `~/.config`, and it fails open. |
| Spend guard | `plugin/scripts/spend-guard.mjs` | cost recording, budget cap | skip | It only blocks pay-per-token API keys on agent fleets. There is no per-session recording and no budget cap. |
| KB/RVF engine, MCP server | `kb/forge-*.mjs`, `plugin/mcp/server.mjs` | knowledge for agents | skip | Tied to `@ruvector/rvf`, `@xenova/transformers` and a large corpus. Lessons and memory live in AgentDB/ruflo SQLite, `~/.claude`, `~/.config`, not in Git. |
| Overlays, primer | `overlays/qwntik/`, `primer/` | setup from a pinned template | skip | Client-specific. The only reusable idea is check/apply that verifies hashes and keeps a rollback copy. |

**What does not fit LAYUP:**
- **Model in the loop:** TriSmart's acceptance is decided by an LLM emitting the text `MODEL_MESH_ACCEPT`, matched by regex. The qualification reviewer is also a model, though it sits behind exact-answer checks.
- **Not stored in Git:** AgentDB/SQLite, `~/.config/ruvnet-brain`, `~/.claude/metaharness/*.jsonl`.
- **Vendor lock-in:** Claude Code plugin hook semantics, the RVF format and ruflo. Model IDs are hardcoded.
- **Gates fail open by design.** That conflicts with LAYUP gates that agents cannot change.

**Claims the code or evals don't back:**
- README says the baseline was measured 2026-07-10 with routed 63/80 and abstain 18/20. `evals/baseline.json` is dated 2026-09-09 and shows 80/80 and 20/20. The jump on a set that is "never used for tuning" can't be checked from a one-commit clone.
- `PROOF.md` and `DESCRIBED-PROOF.md` are from 2026-06-30 but are presented next to v4.5.16. In `scripts/prove.mjs` the "threshold" defaults to -3, so row 20 (relevance -0.724) counts as a pass.
- Repo coverage disagrees across documents: the primer says 77 repos, README says 199 stores, SPEC says about 248.
- The TriSmart skill says accepted receipts are stored in AgentDB, but `review.mjs` only saves blocked checkpoints. The accepted receipt is just printed.
- The skill says the scribe is chosen from the hash of task plus source manifest; the code hashes only the task.
- Verification is not blind: reviewers see a synthesis built from their own proposals. `tri-smart/evals/evals.json` contains only prose expected outputs, with no grader.

**Top 3 recommendations:**
1. **Copy the route qualification and promotion pattern** for retro-driven routing changes. Have the engine compute exact-answer checks, use a reviewer from a different provider on blinded A/B outputs, change one role per milestone, and commit the policy only if its hash matches what was expected.
2. **Copy the frozen eval gate** for LAYUP's evals and quality gate: hash-pinned cases, Wilson lower-bound ratchet, fail on missing baseline, and record a baseline only deliberately.
3. **Lift TriSmart's CLI recipes and process handling** into LAYUP's harness adapters (argument lists, process-group timeouts, stripping API-key variables, cleaning output that may contain prompt injection). Replace its token-based LLM acceptance with LAYUP's own deterministic gate, and make reviewers blind to authorship.
