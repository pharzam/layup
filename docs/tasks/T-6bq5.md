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
