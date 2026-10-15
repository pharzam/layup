# The forge interface and the GitHub adapter

Milestone `M2a` (task `T-zck8`, #123). Requirements: `NFR-001` (the records go
through the forge to the target's Git) and `NFR-007` rule 5 (only the adapter
connects). It derives from [`architecture.md`](../architecture.md) §1 (the forge
and its six capabilities, O-102), §3 (identities) and ADR-0013. The command that
uses it: [`run.md`](run.md). Conventions: [`README.md`](README.md).

The core engine names capabilities, never one platform's API.
`internal/forge` holds the interface; `internal/forge/github` is its one adapter
([`packages.md`](packages.md#the-table-of-m2a)). No other package of LAYUP
opens a connection to the forge; `git` reaches the Git remotes through
`internal/git`.

## The six capabilities

All six are required. The adapter declares the six; a stand-in adapter of a test
that lacks one makes the step `forge` of [`run.md`](run.md#the-steps-of-layup-run---new)
fail, naming it.

| Capability (§1) | Interface method | First used in |
| --------------- | ---------------- | ------------- |
| issues and comments, with the actor and whether an App made it | `OpenIssue`, `Comments(issue)`, `Comment(issue, body)` | `M2a` (steps 7 and 8 of [`run.md`](run.md#the-steps-of-layup-run---new)); `Comment` in `M2b` (the comment of a session) |
| pull requests with a draft state | `OpenDraft`, `MarkReady`, `PullRequest`, `Merge` | `M2e` ([The calls of M2e](#the-calls-of-m2e)) |
| commit statuses bound to a source | `SetStatus`, `Statuses` | `M2e` ([The calls of M2e](#the-calls-of-m2e)) |
| branch rules with bypass actors, read back | `EffectiveRules(branch)` | `M2d` ([`scaffold.md`](scaffold.md#the-read-back)) |
| the repository activity with its actors | `Activity(ref)` | `M2d` (the activity read of the forge check, [`scaffold.md`](scaffold.md#the-activity-read)); `layup audit` in `M2f` |
| an App identity for LAYUP with scoped permissions | `Token`, `Installation`, `UserID(login)`, `Repository` | `M2a` (step 1 of `run.md`) |

`M2a` specifies the calls of the rows it uses; each later milestone gives the
calls of its rows here. The interface of `M2a` holds the method by which an
adapter declares its capabilities and the six methods of
[The calls of M2a](#the-calls-of-m2a); the cell "First used in" of the first row
names the use of `OpenIssue` and `Comments`, and `M2b` specifies the call of
`Comment` ([The calls of M2b](#the-calls-of-m2b); task `T-fsjp`, the Operator's
O-189 of #165: `session.md` has each session of `M2b` post a comment, which
first named `M2c` here, task `T-6bq5`).

The permissions that `M2a` uses ([`run.md`](run.md#the-steps-of-layup-run---new),
step 1) are contents `write`, issues `write` and metadata `read`; a level
`write` holds `read`, and any other level of a permission gives neither (task
`T-6bq5`).

## The App identity

`layup run` alone holds the App's private key (§3). **Decided here:**

- The key file and the App ID come from the forge register
  ([`records.md`](records.md#nfr-001--the-records-of-start), K40). The run
  checks the file before any call: mode 0600, owned by its user,
  one PEM block of an RSA private key, with nothing before or after it but
  space (`RSA PRIVATE KEY`, PKCS #1, or `PRIVATE KEY`, PKCS #8 of the RSA
  algorithm; in both, the RSA key is of version 0, two primes; decided here,
  task `T-1g1q`). Any other state is exit 2, and its message names the file
  and, in each state but a missing file, its mode (the key-file row of
  [Input states](run.md#input-states)).
- The adapter makes a JSON Web Token (`JWT`) signed with that key (RS256): issued
  60 seconds in the past, to absorb a clock skew, valid for 9 minutes (GitHub's
  maximum is 10): `exp` is `iat` plus 9 minutes, with the App ID as issuer. It finds the installation of the App
  on `OWNER/NAME` and makes an **installation token** from it, which GitHub keeps
  valid for one hour. The adapter makes a new token when less than five minutes
  of the old one remain: `Token` gives the token it holds while five minutes or
  more remain, else a new one, so `layup run` calls it before each `git` call that
  needs the token. A token is never written to a file, a record or a log.
  **Decided here** (task `T-ax3r`; changed by task `T-mqty`, as rule 5 lets no
  package but the adapter import `net/http`): when its caller gives no HTTP
  client, the adapter makes one that follows no redirect, so a token goes to
  the `api` of the forge register only, and whose transport reads no proxy of
  the environment; `internal/cli` gives none, and gives the progress function
  of the wait of a rate limit. **Known limit:** the system certificate roots of
  `crypto/x509` read `SSL_CERT_FILE` and `SSL_CERT_DIR` on Linux, so a host
  with a private certificate authority sets them, as for `git`
  ([L-A7](../architecture.md#15-known-limits)).
- `git` gets the token for one call only, from `internal/git`, as an HTTP
  header in `GIT_CONFIG_COUNT`, `GIT_CONFIG_KEY_0` and `GIT_CONFIG_VALUE_0` of
  that call's environment (`http.<web>/.extraHeader`), not on its command line,
  where another user of the host could read it. The value is
  `Authorization: Basic` and the Base64 of `x-access-token:<token>`, the form
  in which GitHub reads an installation token over HTTPS (task `T-xhgz`). This
  adds to the fixed list of the environment ([`packages.md`](packages.md#the-calls-of-internalgit)); no
  value of the host is passed.

## The calls of M2a

The GitHub REST API, version header `X-GitHub-Api-Version: 2022-11-28`, at the
`api` of the forge register. The documentation of each endpoint is read at the
date of the build task; a difference from this table is a defect of this table.

| Method | Call | What it gives `layup run` |
| ------ | ---- | ------------------------- |
| `Installation` | `GET /repos/{owner}/{repo}/installation` (JWT) | the installation ID and its permissions |
| `Token` | `POST /app/installations/{id}/access_tokens` (JWT) | the installation token and its end time |
| `Repository` | `GET /repos/{owner}/{repo}`, then `GET /repos/{owner}/{repo}/branches`, every page | the default branch, the visibility and the names of the branches; the repository has a commit when the list is not empty (task `T-ax3r`) |
| `UserID` | `GET /users/{login}` | the numeric ID of a login; the bot's ID for `<slug>[bot]` |
| `OpenIssue` | `POST /repos/{owner}/{repo}/issues` | the issue number |
| `Comments` | `GET /repos/{owner}/{repo}/issues/{n}/comments`, every page | each comment: its ID, the author's ID and login, `performed_via_github_app` (the App's slug, or none), the times, the body |

`Comments` asks 100 a page and follows the `Link` header's `rel="next"` while it
names a page at the `api` of the forge register; a next page elsewhere is a
forge error. A comment whose `user` is `null` (a deleted account) has the author
ID `0` and the login `ghost`; no approver has the ID `0`, so it is never a
decision (task `T-6bq5`).

The read-back of the root commit and the push of the records branch are `git`
calls, not API calls ([`run.md`](run.md#the-steps-of-layup-run---new)).

## The calls of M2b

The same API and version header as [The calls of M2a](#the-calls-of-m2a)
(task `T-fsjp`, O-189).

| Method | Call | What it gives `layup run` |
| ------ | ---- | ------------------------- |
| `Comment` | `POST /repos/{owner}/{repo}/issues/{n}/comments`, the body `{"body": …}` | the comment's ID; a status other than 201 is a forge error naming `Comment` |

`M2b` posts the comment of a probe and of a task session on the control issue
([`session.md`](session.md#a-comment-for-a-session)).

## The calls of M2d

The same API and version header as [The calls of M2a](#the-calls-of-m2a)
(task `T-msrd`, #265; [`scaffold.md`](scaffold.md#req-003--scaffold-7-the-forge-check)).

| Method | Call | What it gives `layup run` |
| ------ | ---- | ------------------------- |
| `EffectiveRules` | `GET /repos/{owner}/{repo}/rules/branches/{branch}`, every page | each active rule that applies to the branch, also to one that does not exist yet: its `type`, its `parameters` and the ruleset that it is of (`ruleset_source_type`, `ruleset_source`, `ruleset_id`); a ruleset in `evaluate` or `disabled` gives none |
| `Repository` | the calls of `M2a`; it also gives the head of each branch, `commit.sha` of each item of the list of branches | the head of the default branch, for the tree check ([`scaffold.md`](scaffold.md#the-tree-check)) |
| `Activity` | `GET /repos/{owner}/{repo}/activity`, with `ref` the branch and one page | the activities of the ref: each with its `actor`, `ref`, `before`, `after`, `timestamp` and `activity_type`; `M2d` reads only that the call works |

The bypass list of a ruleset is not in either call, and the App's permissions
cannot read it (`architecture.md` §3): the Operator's apply output gives it.
The permission that each call needs, and the fields above, are read in GitHub's
documentation at the date of the build task; a difference from this table is a
defect of this table. The push of the probe is a `git` call, not an API call.

**Decided here** (task `T-jpkc`, #271, with the conditions of its plan review):

- `EffectiveRules` asks 100 a page and follows the `rel="next"` of the `Link`
  header as `Comments` does: a next page that is not at the `api` of the forge
  register is a forge error naming `EffectiveRules`. The branch is a multi-segment
  path parameter of GitHub (`x-multi-segment` in its description), so each segment
  of its name is escaped and its slashes stay. The endpoint answers by the name of the
  branch: a branch with no rule gives an empty list (GitHub answers 200 with `[]`), and
  a branch that does not exist yet, as `layup-probe` is before its probe, gives the
  rules of the rulesets whose conditions name it. The call asks for no list of branches
  and adds no filter of its own, as the endpoint gives the active rules only.
- A rule is its `type` and its `parameters`, decoded. A rule that GitHub gives with
  no `parameters`, or with an empty object of them, has none, so the read-back
  ([`scaffold.md`](scaffold.md#the-read-back)) compares the two forms alike, and reads the one parameter that GitHub omits
  at its default (`update_allows_fetch_and_merge` `false` of `update`) as equal, as that section says (task `T-j7c5`). A rule
  also names the ruleset that it is of, as GitHub gives it (**decided here**, task
  `T-gkvi`, #323): `ruleset_source_type` (`Repository`, `Organization`), `ruleset_source`
  and `ruleset_id`, each the zero value when the answer lacks it or gives JSON `null` (as
  decoding does); a `ruleset_id` that is anything else but an integer (a string, a fraction,
  a boolean, a list) is an error of the call. The read-back reads by branch and does not
  use them; the attribution of the probe's refusal does
  ([`scaffold.md`](scaffold.md#the-probe)).
- `Activity` sends the `ref` as it is given (GitHub takes `refs/heads/BRANCH_NAME` or
  `BRANCH_NAME`), escaped as a query value, with no other parameter, so it reads the
  first page of the default size, 30; it does not follow a `Link`. An activity whose
  `actor` is `null` has the actor ID `0` and the login `ghost`, as a comment of a
  deleted account has the author `0`.
- `Repository` also gives `Heads`, the head of each branch by its name, beside
  `Branches`, which does not change.
- A body of JSON `null` reads as an empty list, as for `Comments` and `Repository`; a
  body that is an object, or is not JSON, is a forge error naming the call and the
  status. The rules of [Forge errors](#forge-errors) hold for both calls with the
  handling that every call shares: one new token and one retry after a 401 of the
  token, a wait after a 403 or a 429 that names a reset time, and a failure at once
  when it names none.
- The permission that each call needs, in GitHub's documentation read on 2026-10-13
  (the page of each endpoint, and the REST description `github/rest-api-description`,
  commit `9492527246942bb85a2e05ad167881076c27835a`): `EffectiveRules`, `Metadata`
  read; `Activity`, `Contents` read. The App has `contents` `write`, which holds
  `read`, and `metadata` `read`, so `CheckPermissions` does not change.

## The calls of M2e

The same API and version header as [The calls of M2a](#the-calls-of-m2a) (task
`T-xs4y`, #333; [`loop.md`](loop.md#req-007--the-four-layup-statuses)).

| Method | Call | What it gives `layup run` |
| ------ | ---- | ------------------------- |
| `OpenDraft` | `POST /repos/{owner}/{repo}/pulls` with `draft` true, the head branch, the base, the title and the body | the pull request's number, its node ID and its head SHA; it requests no review |
| `MarkReady` | the GraphQL mutation `markPullRequestReadyForReview` with the node ID (the REST API has no call that leaves the draft state) | the pull request, no longer a draft; **to confirm in GitHub's documentation at the build** |
| `SetStatus` | `POST /repos/{owner}/{repo}/statuses/{sha}` with `state` (`success`, `failure`, `pending` or `error`), `context` and `description` of at most 140 characters | the status for that SHA and context; it is bound to the head SHA that the caller gives, and it is posted again only when its state or its description changes |
| `Statuses` | `GET /repos/{owner}/{repo}/commits/{sha}/statuses`, every page, and `GET /repos/{owner}/{repo}/commits/{sha}/check-runs`, every page | each status of the head with its context, state and creator (the App or another), and each check run with its name, `conclusion` and the App that made it (a job of GitHub Actions is a check run of the App with the ID 15368) |
| `PullRequest` | `GET /repos/{owner}/{repo}/pulls/{number}`, and `GET /repos/{owner}/{repo}/pulls?head={owner}:{branch}&state=all` to find the pull request of a branch | the pull request's state, whether it is merged, its merge commit's SHA, its head SHA, its draft flag and its node ID |
| `Merge` | `PUT /repos/{owner}/{repo}/pulls/{number}/merge` with `sha` the head SHA and `merge_method` `merge` | the merge commit; GitHub refuses (409) a head that moved |

**Decided here** (decision 5 and 6 of #333):

- **The permissions** beyond those of `M2a` are pull requests `write` (`OpenDraft`,
  `MarkReady`, `Merge`) and commit statuses `write` (`SetStatus`; both are in
  `architecture.md` §3); `Statuses` needs `read` of both and checks `read`, which ADR-0014 decision 3 does
  not list: it is **an amendment of that decision**, asked in #336 and granted by the
  Operator with the others (a draft pull request's `mergeable_state` is `draft`, never
  `clean`, so it cannot say that the native jobs are green before the pull request is
  marked ready). `CheckPermissions` reads the three, and a missing
  one is a `fail` of the step `forge` that names it. The Operator grants them
  (`gov-operator-setup-o112`).
- **`PullRequest`** also serves the restart and the tree rule ([`loop.md`](loop.md#the-tree-after-the-first-merge)): the merged state and the merge commit of a task's pull request, and the pull request of a branch.
- **A comparison reads, never trusts.** The ready and merge rule reads `Statuses` (a check run: the newest of each name, `filter=latest`)
  just before each of `MarkReady` and `Merge`; a context counts only when it is
  `success`, its creator is the App, and it is the newest status of that context at
  the head. A required check of another source with the same name does not count.
- **The stand-in** of `internal/standin` and the httptest forge each give the six
  methods, a status list that a test sets, a pull request whose head a test moves
  (to give the 409), and a status of another source.
- The permission that each call needs, and the fields above, are read in GitHub's
  documentation at the date of the build task; a difference from this table is a
  defect of this table.

## Forge errors

**Decided here** (FT1): a forge call that fails, gives a status other than the
one its call expects, or gives a body that the adapter cannot read, is `fail` on
its step, with the call, the status and the first line of the message; never a
pass and never a retry that hides it. A `401` of an installation token makes one
new token and one retry; a second `401` is `fail`. A rate limit (`403` or `429`
with a reset time) waits until the reset time and prints one progress line every
ten seconds; it is not a stall. The reset time is `retry-after` (seconds), or
`x-ratelimit-reset` (UTC epoch seconds) when `x-ratelimit-remaining` is `0`;
GitHub sends `x-ratelimit-reset` on each response, so a `403` or `429` with
neither is a `fail` at once (GitHub, "Rate limits for the REST API", read on
2026-10-08, task `T-6bq5`).

## The test of the adapter

An integration test of `internal/forge/github` starts an `httptest` server on
loopback that plays the calls above. It checks the JWT (its header, its claims,
its signature with the public key of a test key), the token request, the paging
of `Comments` and of `EffectiveRules`, the one page of `Activity`, the head of each
branch of `Repository`, and each error rule above, also for the calls of `M2d`. A server that lacks a permission or a
capability makes `forge` fail and name it: the permission, which `CheckPermissions`
names from the installation; or the endpoint of a call, and then that call fails with the call and the status. A
capability that an adapter does not declare is the test of `Missing`, with a
stand-in adapter (task `T-6bq5`). The test needs no secret and no
network.
