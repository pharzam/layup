# T-6bq5 — row 24 of the plan, the forge interface and the GitHub adapter

Issue: [#129](https://github.com/pharzam/layup/issues/129), row 24 of [the tasks
of M2a](../plan/README.md#the-tasks-of-m2a), opened by `T-zwke` (#124); O-171 of
#137 moved the forge interface here from row 22b. Serves `F-0003#42` through
`NFR-001` and `NFR-007`. Base `fc3e0d4` (the merge of #139). Author: Claude Opus
5.5 on Claude Code. Evidence: [`runs/T-6bq5/`](../../runs/T-6bq5/).

## Plan and plan review

The plan (R12, comment 6055969539) kept the scope of O-171: the interface and
the adapter that implements it, one goal by the Operator's count. Its review
(Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream output, a
fresh read-only session in a clone at `fc3e0d4`, 7 min 18 s; comment
6056099124): `approve-with-conditions`, found no wider scope; three conditions
(the reading of a rate limit; the contract of `Token`; the one reading of the
acceptance row of `run.md`) and eight notes. The author's answer (6056099433)
took all three and notes 2 to 6, and added the rule for a comment of a deleted
account, found in the endpoint reading. Budget maximum 1,400 lines added plus
removed over 18 files against `fc3e0d4`, close-out inside; Cycle cap 1; no
panel.

## What was done

1. **`internal/forge`** (the work that O-171 moved): `Capability` and the six
   capabilities in the order of `forge.md`; the interface `Forge` of `M2a`
   (`Capabilities`, the declaration, and the six methods of "The calls of
   M2a"; `Comment` is specified by `M2c`); the types `Installation`, `Token`,
   `Repository`, `Comment` (both times) and `Error` (the call, the status, the
   first line); `Missing`; `Permissions` and `CheckPermissions` (`write` holds
   `read`). No network package.
2. **`internal/forge/github`**: the adapter. The JWT (RS256, `iat` now − 60 s,
   `exp` `iat` + 9 min, `iss` the App ID) with `crypto/rsa`; `Token` keeps the
   token while 5 minutes or more remain; a `401` of a token makes one new token
   and one retry; `Comments` follows `rel="next"` at the register's `api` only;
   a comment of a deleted account has the ID `0` and the login `ghost`; a forge
   error is a `forge.Error`; a rate limit is `retry-after`, or
   `x-ratelimit-reset` with `x-ratelimit-remaining` `0`, and its wait prints a
   progress line every 10 s.
3. **The documents:** `forge.md` (the interface of `M2a` and `Comment`, the
   permissions and "`write` holds `read`", `exp`, `Token`, the paging and the
   deleted account, the reading of a rate limit, the one reading of the test of
   the adapter); the acceptance row of `run.md`; a pitfall in `guardrails.md`
   §2; the Job cell of `internal/forge`; the traceability rows; the Test cells
   of `NFR-001` and `NFR-007`, the Task cell of `NFR-007` and a §13 line;
   [`endpoints.md`](../../runs/T-6bq5/endpoints.md), the endpoint reading of
   2026-10-08. [`docs.sh`](../../runs/T-6bq5/docs.sh) checks them.

**Tests:** [`test-runs.md`](../../runs/T-6bq5/test-runs.md): red, then green,
for each package; three mutations of the adapter, each caught (one after the
test of the `401` was made stronger); `docs.sh` red, then green. One 2,048-bit
key per test binary, made at run time: no key file enters the tree.

**The rejected alternatives:** the interface in `internal/forge/github` (O-171);
a JWT or GitHub library (`NFR-007`); a retry of any other status (FT1);
`Comment` in the interface of `M2a`.

## Review rounds

Round 1 (Claude Fable 5.1, effort `xhigh`, on the Claude Code CLI with stream
output, a fresh read-only session in a clone at `e762c77`, 10 min 7 s; comment
6056403036): `nothing material in scope`, seven notes. The reviewer could not
run the mutations (its permission mode refused a copy outside the clone) and
checked them by reading. No note is applied to the code after the last round,
so the reviewed code is the code that lands. The notes stay as known limits:

- Note 1: a `401` of a JWT call has no test; the code retries only a token call.
- Note 2: the reset time of the `403` case equals the fake clock, so its sleep
  check cannot bite; the case catches the mutation by its error.
- Note 3: `UserID` accepts the ID `0` from a body with no `id`.
- Note 4: for a `403` or `429` with neither header, GitHub's page says to wait
  one minute; `forge.md` makes it a `fail` at once (FT1). `endpoints.md` says so.
- Notes 5 and 6, for the caller of rows 25 and 26: `Client` and `Progress` are
  required (no default); the client that the caller passes should refuse
  redirects, so the token stays at the `api`.
- Note 7: one `t.Fatal` runs in a server handler of a test.

## Verdict

Delivered: against a loopback server, the adapter plays each call of `M2a` with
an installation token that it made from a JWT; the forge interface of `M2a`, its
types, the declaration of the six capabilities, `Missing` and
`CheckPermissions` are in `internal/forge`, with no network package. Conditions
4 to 6 and note 5 of the review of #137 are met. The review ended by decay at
cycle 0. The diff against `fc3e0d4` is inside 1,400 lines over 18 files.

Next: row 23 (`T-xhgz`, #128), then row 25 (`T-trej`, #130), whose After cell
holds rows 21, 22b, 23 and 24.

## Resource record

Recorded, not budgeted (ADR-0007). Times are UTC on 2026-10-08; tokens are the
`result` event of the Claude Code CLI (stream runs); `not reported` otherwise.

| Part | Expected tier | Model | Effort | Tokens | Elapsed |
| ---- | ------------- | ----- | ------ | ------ | ------- |
| The plan | reasoning | Claude Opus 5.5 | max | not reported | 08:28 to 08:32 |
| Its plan review | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 589,465 (USD 3.94) | 7 min 18 s, from 08:32 |
| The endpoint reading; the answer | reasoning | Claude Opus 5.5 | max | not reported | 08:33 to 08:40 |
| The work, test first; the documents | execution | Claude Opus 5.5, a reasoning-tier model on an execution part | max | not reported | 08:40 to 08:48 |
| Round 1 | reasoning | Claude Fable 5.1, Claude Code CLI | `xhigh` | 893,350 (USD 4.80) | 10 min 7 s, from 08:48 |
| The close-out | execution | Claude Opus 5.5 | max | not reported | 08:59 to 09:05 |
