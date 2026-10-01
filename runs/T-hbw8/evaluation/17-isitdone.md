# Evaluation: isitdone

| Field | Value |
| ----- | ----- |
| Repository | https://github.com/raimondasl/isitdone |
| Pinned commit | 0051655756ff5b15f21ede78d1727c4a3184dfbb (2026-09-27T15:05:20-04:00, v0.8.2) |
| License | MIT (LICENSE file, "Copyright (c) 2026 raimondasl") |
| Language and needs | TypeScript, Node.js 20 or later, `git`. No runtime dependencies (only dev dependencies: TypeScript, esbuild, vitest). No server, no database, no daemon, no network, no LLM. State in `.isitdone/` in the target (receipt, HMAC key, decision log). |
| Evaluated by | Claude Opus 5.5 on Claude Code, 2026-09-28 |
| Elapsed | about 25 minutes of the session (shared with three other candidates) |
| Agent runs and cost | 1 agent run (Claude Code, `claude-haiku-4-5-20251001`, print mode, spend cap USD 0.50): cost USD 0.057 |

## Verdict

Borrow the pattern. The done-gate and the weakened-test scan work on a Go target, but the Go scan is shallow, it warns by default, and one deleted assertion never blocks; LAYUP's own `layup gate` already runs the native commands. The receipt bound to a tree hash, the attempt cap and the gate-weakening rules are good patterns. LAYUP can also run it as an optional, advisory separate program.

## What it is (from the code)

isitdone is a CLI and a Stop hook. When the agent ends its turn, the hook runs the repository's own check commands on the working tree, scans the diff for weakened tests, and refuses the stop until both pass.

- **Check detection.** For a Go module it adds `go vet ./...` (a "lite" check) and `go test ./...` (a "full" check), from `go.mod`: [src/detect.ts#L344-L349](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/detect.ts#L344-L349). It did not add the toy's `Makefile` target `check`.
- **Weakened-test scan.** A line and regex scan over diff hunks, with no AST ([src/integrity.ts#L1-L7](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1-L7)). How it handles Go:
  - A Go test file is `*_test.go` or a file under `test/`, `tests/` and similar ([src/integrity.ts#L101-L103](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L101-L103)).
  - Go counters ([src/integrity.ts#L254-L258](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L254-L258)): a test is `func Test|Example|Benchmark|Fuzz...`, plus each `t.Run(` ([#L598](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L598)). An assertion is a call to `t.Error*`, `t.Fatal*`, `t.Fail*`, or to testify `assert.*` / `require.*` or `is.*`. A skip is `t.Skip*` or `testing.Short()`.
  - The condition of an `if` that guards `t.Errorf` is not an assertion to the scanner. A change to the condition (for example `false && got != 5`) is not visible to it.
  - Detectors that apply to Go: `skip-added` (severity high, [#L1160-L1173](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1160-L1173)); `assertions-removed` (1 lost assertion is `medium`; 3 lost, or 2 lost and more than 30 %, is `high`, [#L1457](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1457)); `tests-removed`; `test-file-deleted`; `test-case-removed` for rows of a table-driven test; `recover()` added in a test (`medium`, [#L1345](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1345)); a tautology only in testify form ([#L664](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L664)); `go test -run` or `-short` added to a command line (`high`, [#L1488](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1488)).
  - A probable defect: the Go entry of `STRONG_MATCHER` is a copy of the JavaScript matcher (`.toBe(`, `.toEqual(` and so on, [#L650-L653](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L650-L653)). The "a weaker assertion replaced a stronger one" rule therefore cannot match Go code.
- **Gate-weakening rules** for configuration and CI files: `|| true`, `continue-on-error`, `if: false`, `--no-verify`, a removed test step ([src/integrity.ts#L1484-L1541](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1484-L1541), [#L1595-L1604](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L1595-L1604)).
- **Default mode is warn.** Integrity findings block only in `strict` or `--ci` mode ([src/verify.ts#L128](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/verify.ts#L128), [src/config.ts#L34-L35](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/config.ts#L34-L35)). Only unsuppressed `high` and `critical` findings block. `--ci` also turns each new `isitdone: allow` suppression into a `high` finding ([src/integrity.ts#L636-L643](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/integrity.ts#L636-L643)).
- **Receipt.** A JSON receipt bound to the `git write-tree` hash of the working tree, with an HMAC. The key is a local file `.isitdone/key` ([src/receipt.ts#L46](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/receipt.ts#L46), [#L60-L66](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/receipt.ts#L60-L66), [#L158](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/receipt.ts#L158)). The agent can read that key, so the HMAC stops a hand edit, not a determined agent.
- **Attempt cap.** The hook blocks at most 3 times per turn, then lets the agent stop with a message ([src/hook.ts#L18](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/hook.ts#L18), [#L247-L251](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/hook.ts#L247-L251), [#L379-L380](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/hook.ts#L379-L380)). When it cannot write its state, it allows the stop "without enforcement" ([#L375-L377](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/hook.ts#L375-L377)).
- **No checks.** With no detected check, `isitdone` exits 0 (`res.ok` is true, [src/cli.ts#L203](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/cli.ts#L203)); the JSON says `done: false`, `state: NONE`. The hook allows every stop ([src/cli.ts#L416](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/cli.ts#L416)). The MCP tool says "NOT VERIFIED" ([src/mcp.ts#L704-L705](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/mcp.ts#L704-L705)).
- **Hosts.** Hook adapters for 12 harnesses ([src/hosts.ts](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/src/hosts.ts)) and a GitHub Action whose `strict: true` runs `--ci` ([action.yml#L116](https://github.com/raimondasl/isitdone/blob/0051655756ff5b15f21ede78d1727c4a3184dfbb/action.yml#L116)).

## What we ran

Install, in the work directory only: `npm ci` and `npm run build` in the clone (npm cache in the work directory). Wrapper: `node src/dist/isitdone.js`. `HOME` was the work directory's `home/`. The toy target: a Go module with package `calc` (`Add`, `Div`, `Sub`), one test file, a `Makefile`, and a local bare remote.

| Case | Command | Result |
| ---- | ------- | ------ |
| Detect | `isitdone detect` | `vet go vet ./... lite`, `test go test ./... full` |
| Clean change (add `Sub` and `TestSub`) | `isitdone --strict` | `Tests 2 -> 3 Assertions 4 -> 5 ... nothing weakened`, `DONE`, exit 0. Not blocked. |
| 1 assertion deleted (`Div(1,0)` check) | default, `--strict`, `--ci` | `assertion count dropped 5 -> 4 with no test removed  medium`; `DONE`, exit 0 in all three modes. **Not blocked.** |
| 2 assertions deleted | default | `high`, "warn only", `DONE`, exit 0 |
| 2 assertions deleted | `--strict` | `NOT DONE checks passed but 1 test-integrity finding blocks`, exit 1. **Blocked.** |
| `t.Skip("flaky")` added | default | `test skipped or marked expected-failure  high`, "warn only", `DONE`, exit 0 |
| `t.Skip("flaky")` added | `--strict` | `NOT DONE`, exit 1. **Blocked.** |
| `t.Skip` with `// isitdone: allow flaky on CI` | `--strict` / `--ci` | strict: `(allowed ...)`, `DONE`; ci: `suppression added in CI mode`, `NOT DONE` |
| Assertion made dead (`if ...; false && got != 5`) | `--strict` | `nothing weakened`, `DONE`. **Evasion not detected.** |
| Build break (`return a + c`) | default | `go vet ./... FAIL ... undefined: c`, `go test ./... SKIP skipped because a lite check failed`, `NOT DONE`, exit 1. **Blocked.** |
| Repository with no checks | default | `no checks detected`, exit 0; JSON `"done": false, "state": "NONE"`; synthetic Stop payload allowed with no output |
| Stop hook, failing test, synthetic payloads | `hook --host claude` five times | 3 x `{"decision":"block" ... attempt 1/3 .. 3/3}`, then `checks still failing after 3 attempts; allowing the agent to stop`. `.isitdone/decisions.jsonl` records `blocked` x3, `gave-up` x2. |
| Real agent run | `init --agent claude --command "node .../isitdone.js hook --host claude"`, then `claude -p` (haiku) told to add `Mul` and not to run commands, on a tree with a seeded bug in `Add` | The hook blocked the first stop (`kind: blocked, failed: 1`). The agent fixed `Add`; the second stop passed (`kind: passed`). `go test ./...` passed after the run. 32 s, USD 0.057. |

`init` wrote `.claude/settings.json` (Stop and PostToolUse hooks) and a `.gitignore` line into the toy target. `doctor` then reported `synthetic "tests pass" stop was blocked in 0.3s`.

## In-Scope items S1-S12

| Item | Mark | Evidence |
| ---- | ---- | -------- |
| S1 Problem Statement Quality | no | No problem-statement input. |
| S2 Reproducible Discipline Setup | no | `init` installs one hook; no baseline. |
| S3 Rule Protection | partly | It detects diffs that weaken gates or tests (integrity.ts#L1484-L1604). It does not prevent them. Its own config `.isitdone.json` and suppression comments are in the agent's workspace. |
| S4 Stack-Dependent Gates | partly | It selects checks by stack (go.mod gives `go vet`, `go test`). No layout, interface or contract gates. |
| S5 Role Handoffs | no | None. |
| S6 Autonomous Clarification | no | None. |
| S7 Verification on Every Change | partly | Deterministic: it ran the native checks on the exact tree and blocked the build break and a failing test. It has no layout or interface-boundary gate. It blocks the agent's "done", not a PR. |
| S8 Human-on-the-Loop | no | None. |
| S9 Stall Resolution | partly | A bounded block count (3 per turn), then release with a message. No diagnosis, no record in the repository. |
| S10 Cost Visibility | no | Check durations only. |
| S11 Specification Synthesis | no | None. |
| S12 Harness-Agent Neutrality | partly | Hook adapters for 12 harnesses and one CLI; the receipt is harness-independent. |

## Invariants I1-I9

| Invariant | Effect | Reason |
| --------- | ------ | ------ |
| I1 Git is the system of record | neutral | The receipt and `decisions.jsonl` are in `.isitdone/`, which `init` puts in `.gitignore`. The data is disposable. LAYUP would copy the result into its own records. |
| I2 Independent repository | supports | Run from outside (`--cwd`), it writes only `.isitdone/`. `init` writes hook entries into `.claude/settings.json` of the target, which O-76 forbids for LAYUP files. Run it, do not `init` it. |
| I3 Agents cannot change the gates | conflicts (as installed) | The agent can edit `.isitdone.json` (`"integrity": "off"`), add `isitdone: allow` comments (they pass in `--strict`), and read `.isitdone/key`. Run by LAYUP from outside with `--ci` and explicit checks, this is neutral. |
| I4 No value without evidence | neutral | The check commands come from the target's own files. The severity thresholds (1 lost assertion is `medium`) are fixed constants with no stated evidence. |
| I5 An inactive check is not a pass | conflicts (partly) | No checks gives exit 0. A lite failure marks the full check `SKIP`, not `PASS` (good). After 3 blocks, or with unwritable state, the stop is allowed; the receipt stays `FAIL`. |
| I6 Deterministic first | supports | Zero LLM calls; exit codes. |
| I7 Stack adds, never removes | supports | Stack detection only adds checks. |
| I8 Pinned Armature | neutral | Not related. |
| I9 Replaceable harness | supports | 12 hosts; same CLI for each. |

## Deep-check findings it answers

- **Sol-2** (a PR with a failed gate reaches review): pattern only. The check runs before the agent may claim done, so a red tree never becomes a "done" claim. LAYUP can run the target's gates before it opens the PR, the same way. The GitHub Action runs on an open PR, so it does not answer Sol-2 by itself.
- **Sol-22** (a new record resets the stall clock): pattern. The hook counts attempts per turn and releases after a fixed cap, whatever the diff did. That is an absolute attempt cap that artifact changes do not reset.
- **Fable-M19** (a wrong gate stops the milestone with no exit): partly, as a pattern. The block message tells the agent: "if a check is wrong for this repo, say so explicitly to the user", and the cap ends the block. A LAYUP stall record should carry the same "the gate may be wrong" branch.
- **Sol-1** (rule protection is detection): it is detection too. The gate-weakening rules are a usable detector list; they do not prevent a change.
- **S7 and S12** (verification and neutrality): partly, above.
- **Fable-M18, Sol-24, Fable-M6**: none. It has no panel and no diagnosis.

## Parts and their verdicts

| Part | Verdict | How LAYUP would take it | Reason |
| ---- | ------- | ----------------------- | ------ |
| Done-gate (run the native checks on the exact tree) | borrow the pattern | `layup gate` already runs the native commands; keep the rule "no pass without a run on the exact tree hash" | Same function; no gain from a Node dependency. |
| Weakened-test scan | use (optional, advisory) or borrow the pattern | Separate program: `npx @aivolution/isitdone --ci --json --base <sha> --cwd <checkout>`; LAYUP reads `integrity.findings` | MIT, no network. For Go it caught `t.Skip` and 2 deleted assertions; it missed 1 deleted assertion (never blocks) and a dead condition. Not enough alone for S7. |
| Gate-weakening rules (`|| true`, `continue-on-error`, removed test step) | borrow the pattern | Port the list to a Go check in `layup gate` | Simple regexes; useful against Sol-1 as detection. |
| Receipt bound to `git write-tree` | borrow the pattern | Record the tree hash in each gate record | Makes a stale pass visible. The HMAC key in the workspace adds little. |
| Attempt cap with release message | borrow the pattern | Fixed cap per step, then a stall record | Answers part of Sol-22. |
| Hook adapters and `docs/agent-hooks.json` | borrow the pattern | Reference for harness Stop hooks | Only if LAYUP uses in-session hooks; O-76 keeps them out of the target. |
| Overall | borrow the pattern | See above | Evidence: the table in "What we ran". |

## Where the searchers were wrong or incomplete

- Searcher B: "blocks a coding agent's done until the repo's real test, typecheck and lint commands pass" is true. "Scans the diff for weakened tests" is true, but incomplete: the scan **warns by default** and blocks only in `strict` or `--ci`. One deleted Go assertion does not block in any mode.
- Searcher B names isitdone a pattern for Rule Protection. The code detects weakening in the diff; the agent can switch the scan off in `.isitdone.json` in its own workspace.
- Searcher A: "Git-bound deterministic completion" is true for the receipt (tree hash). The receipt is not in Git; `init` puts `.isitdone/` in `.gitignore`.
- Neither searcher noted that the tool exits 0 when it finds no checks.

## Limits of this evaluation

- One agent run (Claude Code, haiku). The other 11 host adapters were not run.
- The GitHub Action and the MCP server were not run.
- `isitdone history` reads `~/.claude/projects`; it was not run.
- The claim that the Go `STRONG_MATCHER` never matches Go code comes from reading the code, not from a failing case.
- `init` wrote a `.claude/settings.json` into the toy target only; the host's own settings were not touched by isitdone.
