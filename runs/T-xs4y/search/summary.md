# T-xs4y — the public-solution search of M2e: the result

Two blind searchers screened the same sources ([`sources.md`](sources.md)):
[`searcher-A.md`](searcher-A.md), GPT-6 Sol on Devin, and
[`searcher-B.md`](searcher-B.md), Claude Sonnet 5.5 on Claude Code. This file is
the author's reading of the two, written after both ended, and committed before
the decisions of S3.

**Both conclude: no public solution can be the base of `M2e`.** A: "No source
establishes the whole sequence"; B: "Typed handoffs in Go stdlib: none". What
transfers is GitHub's own API, patterns and formats.

| Part | Candidates (searcher) | What they bear on in `M2e` |
| ---- | --------------------- | -------------------------- |
| The typed handoff and the closed table of transitions | agent-context-store (A, B); agentctl (A); AgentPlane (B); agent-handoff (B); gh-aw "safe outputs" (B); TAKT (A) | A handoff record validated against a schema, with role-to-role rules in a file that code reads; gh-aw is the nearest statement that what an agent writes never makes the step valid. A says a Git-backed handoff file alone does not prove that the controller authorized the transition. |
| An artifact that is only an issue comment | none (A, B); Orbi, issues as the state store (B) | Nothing found that gives such a comment an immutable place and a hash. Decision 1 is made from the architecture. |
| The four statuses, ready and merge rule | GitHub commit-status, rulesets and merge APIs (A, B); cloudposse/github-status-updater, bulldozer (A, B); Orbi, isitdone (B) | `POST /repos/{o}/{r}/statuses/{sha}` with a `context` per class; the required check is pinned to the App as source; the merge call takes the head `sha` and rejects a moved head (A: compare and swap); GitHub accepts `neutral` and `skipped` for a required check, so "not active is not a pass" needs code's own rule (A); a required check that never reports leaves the page silent (B). Orbi and isitdone bind a verdict to a commit or tree hash (B). |
| The close-out commit after the verified head | Orbi "only the reviewed head can merge" (B); the merge `sha` (A) | The verdict is carried to the new head by code's own check of what the last commit changes; the forge does not do it. |
| The filter of options by code | adrkit, WorkSpec Decision Studio (A); Mneme (B) | Typed options and model-free evaluation exist; none drops an option with a failed, missing or evidence-less row and requires one survivor. |
| The brief to the idea owner, answered by one comment | none (A, B) | Nothing found. |
| Requirements with a trace and a full checker | reqmd.dev, supersigil, Spec Kit, OpenSpec (A); rtmx, Doorstop, Duvet, LOBSTER (B) | Formats of a trace graph; reqmd's missing links are warnings by default; Duvet quotes the stating text. None starts from a prose problem statement with the four failure rules together. |
| Controller records on a branch only the controller writes | grite, entireio/cli (B); AgentPlane (B) | Records in a dedicated ref or in the tree; none enforces a single writer, which a ruleset does (A: an App with no administration permission cannot set it up). |

**What the search gives the decisions of S3:** the commit status and merge APIs
fix decision 5 and 6 (a status per context, pinned to the App; merge at the
head `sha`); the closed table with code-owned validity, as the best of the
handoff tools do, fixes decision 3; the immutable place and hash of a comment
artifact (decision 1) is a gap that the architecture fills by copying; the filter
(decision 7), the brief (decision 8) and the full spec check (decision 9) have no
public base, so the specification keeps the architecture's rule.

## The limits of the search

- A ran 2 min 55 s and B about 4 min, both short; each file states its own
  limits (B's lists were read by sub-agents; neither read a candidate's code).
- Facts about GitHub that no page confirmed (B): how a required check that never
  reports is treated; whether a ruleset can exempt one ref by condition for an
  App. The build task reads GitHub's documentation at its date, as `forge.md`
  asks.
