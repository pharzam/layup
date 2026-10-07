# Report: kolezka/self-improvement-loop

The report of one read-only agent of task `T-w73g` (#133), word for word below
the line. Model: Claude Opus 5.5, an Explore subagent on Claude Code, 2026-10-07,
11 tool calls, 64.4 s. Repository: <https://github.com/kolezka/self-improvement-loop>
at `ba245adbb5be8075fa01044dd82aac636155870f`, a shallow clone in a scratch
directory outside the checkout. No code from the repository was run. The
Operator added this repository during the task (the scope amendment on #133).

---

**License.** `LICENSE` is PolyForm Noncommercial 1.0.0 with a "Required Notice" line: Copyright (c) 2026 Mariusz Rakus, and "Commercial use is NOT permitted". You may copy, modify and redistribute the code only for noncommercial purposes, and anything you distribute must carry the license or its URL plus the Required Notice line. `COMMERCIAL-LICENSE.md` says use by or for a for-profit company (internal tooling included), paid client work, or a product or service built on it is commercial use. That needs a separate written license from mariusz@raqz.pl. The license covers copyright and patents in the software, not ideas, so reimplementing the designs from scratch is not restricted. Copying the code, or porting it line by line from TypeScript to Go, is a derivative work. So for a commercial LAYUP: copy patterns, not code.

**What it is.** A Bun/TypeScript plugin for Claude Code, Codex and OpenClaw. Hooks log session events; an offline worker makes one model call per session to write a "reflection"; reflections cluster by pattern slug; when a pattern reaches a threshold it is drafted into a skill, hook, rule or agent on a Git branch, and a human accepts it.

| Item | Path | LAYUP concern | Verdict | Why |
|---|---|---|---|---|
| Watermark promotion planner | `packages/curriculum/src/plan.ts` | Retro lessons → learning loop | Copy pattern | Deterministic. A pattern needs `threshold` new items past `max(promoted_at_count, rejected_at_count)`. Because a rejection resets the bar too, a refused lesson can't come straight back. Its outcomes: promote / refine / escalate / observing / over-cap / retire-candidate / done. |
| Redraft policy (escalate when rewording doesn't help) | `packages/curriculum/src/plan.ts` (`redraftPolicy`) | Stall detection, escalate to human | Copy pattern | If the recurrence rate since the last revision is still at least `escalate_ratio` × the rate since promotion after `max_rewords` revisions, it escalates to a human ("change the artifact type"). No model involved. |
| Recurrence rate windows | `packages/feedback/src/rates.ts` | Checking whether a lesson worked | Copy pattern | Pure function: hit rate before promotion, since promotion and since revision, with a minimum session count. LAYUP can measure lesson effect per milestone the same way from its TSV records. |
| Scorecard proposals | `packages/feedback/src/index.ts` | Learned routing at retrospective | Copy pattern | Counts plus a rule: misfire+bad ≥2 more than last time and above helpful+good → refine; unused for N days → retire. The "misfire" and "helpful" inputs come from the model critic. |
| Artifact lint | `packages/curriculum/src/lint.ts` | Quality gate before any judge | Copy pattern | Deterministic and runs first: shape, slug, minimum length, secret scan, leftover placeholders, echoes of the drafting prompt, enough terms shared with the sources. |
| Verbatim-quote evidence check | `packages/curriculum/src/router.ts` (`substantiveQuote`) | Making model claims checkable | Copy pattern | A model claim only counts if it quotes the source verbatim: ≥30 chars, ≥5 words, ≥3 distinctive terms, word-aligned, with template lines stripped first. |
| Deterministic gate language plus corpus replay | `packages/nudges/src/gates.ts`, `packages/nudges/src/gate-runner.ts` | Hooks, rules agents can't loosen | Copy pattern | A small predicate tree (`all`/`any`/`not`, regex with a ReDoS lint). Drafted gates are replayed against recorded payloads in a separate process with a timeout. |
| Digest-bound accept | `packages/review/src/snapshot.ts`, `packages/review/src/index.ts` | Human approval, Git as system of record | Copy pattern | Accept takes the sha256 of the exact diff shown (branch, commit, base). If anything moved since, accept is refused. Staging happens in a scratch worktree, never the user's checkout. |
| One-hop alias map and Jaccard suggestions | `packages/store/src/aliases.ts`, `packages/store/src/alias-suggest.ts` | Deduplicating lesson keys | Copy pattern | Deterministic; it only proposes and a human applies. Writes that would create an alias chain are refused. |
| Critic (one model call per session) | `packages/critic/src/index.ts` | Lesson capture | Copy pattern partly | Fixed-schema JSON, strict parsing, unknown artifact names are never guessed. But the lesson itself comes from the model, and it is pushed to the inbox if the model's own confidence is ≥0.5. |
| Judge (5 typed questions, Jev/System One) | `packages/curriculum/src/prompts.ts` (`JUDGE_RULES`), `packages/providers/src/index.ts` | Smart-if provider | Copy pattern | Five named yes/no questions (contradicts, vague, unsupported, unsafe, unrelated) with a probability threshold. Fits as a smart-if decision point, not as an engine check. |
| Records outside Git | `~/.local/share`, `~/.local/state` (README "Configuration") | Git as system of record | Skip | Reflections, ledger, feedback and queue are JSON/JSONL under XDG dirs. Only promoted artifacts land in the Git target repo. |
| Vendor-specific delivery | `apps/hook/src/handlers.ts`, `hooks/hooks.json` | Replaceable harnesses | Skip | Built on Claude Code hook events; Codex sessions are not recorded at all. |

**README claims the code doesn't back:**
- "Every chat call runs at temperature 0": the `claude-cli` endpoint (`chatClaudeCli` in `packages/providers/src/index.ts`) sets no temperature.
- A lesson is "marked so it never repeats in a later session": the delivered list is kept per session, and a lesson is only archived after 5 deliveries (`LESSON_ARCHIVE_AT_DELIVERIES` in `packages/core/src/lessons.ts`, `apps/hook/src/lessons.ts`). So it can appear in up to 5 sessions.
- "Nothing reaches a shared remote without a human": `promotion.auto_merge` merges into the target's default branch without review. It is off by default and never pushes, but it does skip the human.

**Top 3 recommendations (reimplement in Go):**
1. **Watermark + redraft-policy + rate windows** as LAYUP's retro learning loop. All three are deterministic over records that can live as TSV in Git. "Rewording didn't help" becomes an escalation to a human.
2. **Gate order: deterministic lint and quote check, then a typed judge.** Lint and the verbatim-quote check are the gate; the five-question judge runs at a named smart-if decision point. That keeps model judgement out of the engine's own checks.
3. **Digest-bound accept with staging in a scratch worktree** for promoting lessons into the pinned template. Rejection counts toward the watermark just like promotion does.
