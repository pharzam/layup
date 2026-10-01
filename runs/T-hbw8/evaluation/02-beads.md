# Evaluation: Beads

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/gastownhall/beads (not moved; the Go module path is still `github.com/steveyegge/beads`) |
| Pinned commit | 3d76cbd4e91056c02bd724aca74ac5fcec0944ac (2026-09-28, `bd version 1.3.0`). Gas Town uses the tag v1.0.5 (6a3f515c, 2026-05-28); I built that one too. |
| License | MIT, "Copyright (c) 2025 Beads Contributors" (the LICENSE file); a `THIRD_PARTY_LICENSES` file lists the dependencies |
| Language and needs | Go 1.26, about 371,000 lines of non-test Go, 694 modules in `go list -m all`, a 196 MB binary. Storage is Dolt only: an embedded Dolt engine by default, or a `dolt sql-server`. Build tag `gms_pure_go`. Optional: the Anthropic API for `compact` and duplicate search. No daemon in embedded mode. |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 45 minutes of wall-clock time for Gas Town and Beads together (15:30 to 16:15) |
| Agent runs and cost | none for Beads alone (Beads is a tracker; the agent run is in the Gas Town report) |

## Verdict

**Borrow the pattern**: the dependency graph with a "ready" query, the atomic claim, the lease with heartbeat and reclaim, and the short hash IDs. **Reject** Beads as a program or a format that LAYUP depends on. By default the state is a binary Dolt database that Git ignores; `bd init` also changes the target's `core.hooksPath` and writes agent instruction files into it (Invariants 1 and 2, decision O-76).

## What it is (from the code)

- **Storage.** "Dolt is the default and only supported storage backend" ([cmd/bd/init.go#L334](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/init.go#L334)). The database is `.beads/embeddeddolt/`, which `.beads/.gitignore` excludes. Each write is a Dolt commit (seen in `dolt log`: `bd: close toy-q4p`, `bd: comment toy-q4p`).
- **JSONL.** Export is optional and off by default: "for viewers (bv), interchange, and issue-level migration; not backup. Cross-machine sync and backups use Dolt remotes" ([cmd/bd/init.go#L351](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/init.go#L351)).
- **Sync through Git.** Dolt pushes its chunks to the same Git remote under `refs/dolt/data`; `bd` probes that ref with `git ls-remote` ([cmd/bd/sync_git.go#L124](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/sync_git.go#L124)). `bd bootstrap` clones from it, or imports a tracked `.beads/issues.jsonl`.
- **Merge of two writers.** A cell conflict "is settled last-write-wins by `updated_at`" ([internal/storage/versioncontrolops/automerge.go#L440](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/internal/storage/versioncontrolops/automerge.go#L440)), with a notice ([#L579](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/internal/storage/versioncontrolops/automerge.go#L579)).
- **Claim.** `ClaimIssueInTx` sets the assignee and `in_progress` only if the issue is open and unassigned, or already the caller's; else `ErrAlreadyClaimed` ([internal/storage/issueops/claim.go#L36](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/internal/storage/issueops/claim.go#L36)). It also rewrites a row lock so that a racing claim conflicts.
- **Lease.** A claim takes a lease, `DefaultLeaseTTL = 5 * time.Minute`, kept alive by `bd heartbeat` ([internal/storage/issueops/lease.go#L26](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/internal/storage/issueops/lease.go#L26)). `bd reclaim` reverts stale leases to open; it is meant to run "from a supervisor on a timer" ([cmd/bd/reclaim.go#L17](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/reclaim.go#L17)). Leases and heartbeats are not versioned; claims are.
- **Gates.** `bd gate` gates are "async wait conditions that block workflow steps": `human`, `timer`, `gh:run`, `gh:pr`, `bead` ([cmd/bd/gate.go#L22](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/gate.go#L22)). They are waits, not checks of a change.
- **What `bd init` does to the target.** It writes `.claude/settings.json` (a SessionStart hook `bd prime --hook-json`), `.codex/`, `.cursor/`, `.agents/skills/beads/`, `AGENTS.md`, `CLAUDE.md` and `.beads/hooks/*`, and commits them ([cmd/bd/init.go#L3565](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/init.go#L3565)). It sets `core.hooksPath` to `.beads/hooks` ([cmd/bd/init_git_hooks.go#L378](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/cmd/bd/init_git_hooks.go#L378)). The hooks run `bd hooks run <hook>` only `if command -v bd`, so a clone without `bd` still commits.
- **LLM use.** `compact` calls the Anthropic API ([internal/compact/haiku.go#L14](https://github.com/gastownhall/beads/blob/3d76cbd4e91056c02bd724aca74ac5fcec0944ac/internal/compact/haiku.go#L14)). The core tracker does not.

## What we ran

All under `eval-work/beads/`, with `HOME=eval-work/beads/home`, `GIT_CONFIG_GLOBAL` there, no `GH_TOKEN`. `bd` 1.3.0 built from the pin; `dolt` 2.3.5 from the release binary.

1. **Init.** `bd init --non-interactive -p toy` in the toy Go repo. Output: `Backend: dolt  Mode: embedded`, `✓ Claude Code integration installed`, `✓ Codex native hooks installed`, `✓ Cursor integration installed`, `✓ Committed beads files to git`. The new commit `105fa39 bd init: initialize beads issue tracking` has 19 files and 990 lines, with no issue data. `.beads/config.yaml` in that commit holds `sync.remote:` with the absolute local path of my bare remote (a machine-specific value in Git).
2. **Issues and dependencies.** `bd create` x3, `bd dep add` x3. `bd ready` gave only `toy-q4p`; `bd blocked` gave `toy-6gu: Blocked by 2 open dependencies: [toy-9iz toy-q4p]`.
3. **Claim race.** `BD_ACTOR=agent-1 bd update toy-q4p --claim` passed. Then `BD_ACTOR=agent-2 ... --claim`: `Error updating toy-q4p: issue already claimed by agent-1`, exit 1.
4. **Close.** A comment, a code commit, `bd close toy-q4p --reason "Done in 0dcb8db"`. `bd ready` then gave `toy-9iz`. `git status` was clean: nothing of the issue work went into the working tree.
5. **Git push.** `git push origin main` sent only code. `bd dolt push` then added `refs/dolt/data` and `refs/heads/__dolt_remote_info__` to the bare remote. The tree of `refs/dolt/data` is binary files (`manifest`, a 466,847-byte `.darc` chunk file).
6. **A human with Git only.** A plain `git clone` (PATH without `bd` and `dolt`): no issue text anywhere in the checkout; the clone does not even fetch `refs/dolt/data`. After a manual `git fetch origin refs/dolt/data:refs/dolt/data`, `strings` on the blobs shows fragments (`Add Mul function`, `...go test passes`) but no structure.
7. **A human with Dolt only.** `dolt clone git+file://.../remote.git` and `dolt sql -q "select id,title,status,assignee,close_reason from issues"` gave the three issues and the `dependencies` rows. `dolt log` shows one commit per `bd` write.
8. **A second clone with `bd`.** `bd list` in a fresh clone: `Error: no beads database found`. `bd bootstrap`: `Synced database from .../remote.git`. `BD_ACTOR=bob bd update toy-9iz --claim` and `bd dolt push`; the first clone ran `bd dolt pull` and saw `Assignee: bob`.
9. **Two writers, one field.** Clone A set the title of `toy-6gu` to "(alice)" and pushed. Clone B set it to "(bob)"; its push was refused; `bd dolt pull` printed `Notice: auto-merged issue toy-6gu; title, updated_at settled last-write-wins (the older side's edit was superseded)`.
10. **JSONL path.** `bd export -o .beads/issues.jsonl` (one JSON object per issue with dependencies and comments), committed and pushed. A clone with no Dolt remote ran `bd bootstrap`: `Imported 4 issues from .../.beads/issues.jsonl`. The JSONL is readable with `grep` or `jq`.
11. **Lease.** `BD_ACTOR=dead-worker bd update toy-vtq --claim` at 16:06:07, no heartbeat. `bd reclaim --older-than 1s` at once: `No stale leases to reclaim`. At 16:11:39: `✓ Reclaimed 1 stale-lease issue(s): toy-vtq (was held by dead-worker)`; the issue was open again.
12. **With Gas Town.** Gas Town `main` with this `bd` 1.3.0 failed: `refusing to auto-apply 20 pending schema migrations to a shared server database (v49 -> v69)`. With `bd` v1.0.5 it worked (Gas Town report).

## In-Scope items S1–S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | A tracker; no problem-statement reading. |
| S2 Reproducible Discipline Setup | no | `bd init` sets up Beads, not a baseline; it writes a machine path into the committed config (run step 1). |
| S3 Rule Protection | no | Any actor can change any issue; `bd init` writes the agent instruction files and hooks that the agents read. |
| S4 Stack-Dependent Gates | no | Beads gates are waits (gate.go#L22), not checks of a change. |
| S5 Role Handoffs | partly | Typed issue fields, dependency types and comments; an export with a stable JSON shape. No role transitions and no validation of a handoff's content. |
| S6 Autonomous Clarification | no | Not in the code. |
| S7 Verification on Every Change | no | Not in the code. |
| S8 Human-on-the-Loop | partly | A `human` gate blocks a formula step until a person closes it. No escalation rule. |
| S9 Stall Resolution | partly | A lease, a heartbeat and `bd reclaim` recover a dead worker (run step 11). No disagreement or no-progress detection, no diagnosis. |
| S10 Cost Visibility | no | No token or time record per issue. |
| S11 Specification Synthesis | no | Not in the code. |
| S12 Harness-Agent Neutrality | partly | A plain CLI; `bd setup` writes integrations for Claude Code, Codex, Cursor and others; state is not in a harness format. It writes one file per harness into the target. |

## Invariants I1–I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | conflicts | The default record is `.beads/embeddeddolt/`, ignored by Git (run step 4). `bd dolt push` can put it in the same Git remote under `refs/dolt/data`, but as binary Dolt chunks that `git clone` does not fetch. A committed JSONL export is possible, but Beads calls it "not backup" (init.go#L351). |
| I2 The repository is independent | conflicts | With Git only, a human cannot read the issues (run step 6). With `dolt` alone, a human can read them by SQL (run step 7). With a committed JSONL, a human can read them and `bd` can restore them (run step 10). `bd init` also puts `bd`-specific hooks, `AGENTS.md`, `CLAUDE.md` and harness settings into the target; this is against O-76 if LAYUP set it up. |
| I3 Agents cannot change the rules | conflicts | `bd init` replaces `core.hooksPath` with `.beads/hooks` (init_git_hooks.go#L378). In an Armature target this displaces `.githooks`, the target's own gate hooks. |
| I4 No value without evidence | neutral | A fixed lease TTL of 5 minutes (lease.go#L26); a tracker default, not a project rule. |
| I5 An inactive check is not a pass | neutral | Beads has no checks of a change. Its hooks skip when `bd` is missing, but they are not gates. |
| I6 Deterministic before LLM | supports | Claim, ready, blocked and reclaim are deterministic queries; the LLM is only in `compact`. |
| I7 Domain changes content only | neutral | Not related. |
| I8 Pinned Armature | neutral | Not related. |
| I9 Harness replaceable | supports | One CLI for any harness; data in Dolt or JSONL, not in a harness format. |

## Deep-check findings it answers

- **Author-10 and Sol-4** (the dead-man job): partly, as a pattern. A claim holds a lease with a TTL; the worker renews it; a timer job reclaims stale leases (reclaim.go#L17). Seen in run step 11. The liveness data lives beside the claim, not in a separate heartbeat file.
- **Fable-M5** (the clock runs during a human wait): partly, as a pattern. A `human` gate is an explicit waiting state of a step (gate.go#L22). A watchdog can skip a step that waits on an open human gate.
- **Fable-M12 and Sol-3** (many writers): none for Git. Beads solves many writers with a database that merges cells and settles conflicts last-write-wins, which loses one edit without a stop (run step 9). This is a warning for LAYUP, not an answer.
- **Sol-22, Fable-M4, Fable-M6, Author-1, Sol-27, Fable-M9**: none. A lease shows that a worker is alive, not that it progresses or agrees.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Beads as the task store of a target | reject | - | I1, I2 (binary Dolt state, gitignored), I3 (`core.hooksPath`), O-76 (writes `bd` files into the target), 694 modules against the stdlib-only rule. |
| Dolt-in-Git (`refs/dolt/data`) | reject | - | In the Git remote, but not readable or diffable with Git; needs `dolt` to read. |
| JSONL issue export as a file format | borrow the pattern | LAYUP's own task file in the records, one JSON line per task with dependencies and comments; optionally a shape that `bd import` accepts, so a human can load it into `bd` | Readable with Git alone; `bd bootstrap` restored it (run step 10). |
| Ready and blocked queries over a dependency graph | borrow the pattern | Compute the ready set from the task table in stdlib Go | Seen: `bd ready` and `bd blocked` gave the correct sets. |
| Atomic claim (compare-and-set on assignee and status) | borrow the pattern | One writer (the orchestrator) grants claims, so a simple check is enough | Seen: the second claim failed with the owner's name. |
| Lease + heartbeat + reclaim | borrow the pattern | A lease per running role session, renewed by the session's hook; the orchestrator reclaims after TTL plus grace; TTL is a parameter (O-82) | Seen: reclaim after 5 minutes. |
| Short hash IDs (`toy-q4p`) | borrow the pattern | IDs that do not collide when two sessions create tasks | Useful for offline writers. |
| Cell merge with last-write-wins | reject | - | Loses an edit silently (run step 9); against a single-writer record. |
| Vendored code | reject | - | Deep dependency tree (Dolt engine, cobra, viper, OTEL, Anthropic SDK); LAYUP is stdlib only. |

## Where the searchers were wrong or incomplete

- Searcher B calls Beads a "Git-backed issue graph". By default the store is an embedded Dolt database that `.beads/.gitignore` excludes; `bd init` commits only config and agent files. "Git-backed" is true only in two opt-in forms: Dolt chunks under `refs/dolt/data`, or a JSONL export.
- Searcher A, "Dolt state not Git record", is right for the default. It misses that Dolt data can live in the same Git remote (binary) and that a committed JSONL can restore the tracker.
- Neither searcher said that `bd init` changes the target: `core.hooksPath`, `AGENTS.md`, `CLAUDE.md`, and harness settings for Claude Code, Codex and Cursor, all in one commit.
- Neither searcher said that concurrent edits merge last-write-wins, or that Beads `main` is not compatible with Gas Town `main`.

## Limits of this evaluation

- Only embedded mode was run. Server mode, the proxied server and federation were not run.
- No agent used Beads directly in this part; the Gas Town run used Beads v1.0.5 through `gt`.
- The `human`, `timer`, `gh:*` gates and formulas were read, not run.
- A `dolt version` call without a redirected HOME made `~/.dolt` (a version file and a lock) in the real HOME; I removed it.
- Processes: embedded mode starts no server. No process of mine remained after the runs (checked with `ps`).
