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
| issues and comments, with the actor and whether an App made it | `OpenIssue`, `Comments(issue)`, `Comment(issue, body)` | `M2a` (Start 7 and 8) |
| pull requests with a draft state | `OpenDraft`, `MarkReady`, `Merge` | `M2e` |
| commit statuses bound to a source | `SetStatus` | `M2e` |
| branch rules with bypass actors, read back | `EffectiveRules(branch)` | `M2d` |
| the repository activity with its actors | `Activity(ref)` | `M2f` |
| an App identity for LAYUP with scoped permissions | `Token`, `Installation`, `UserID(login)`, `Repository` | `M2a` (Start 1) |

`M2a` specifies the calls of the rows it uses; each later milestone gives the
calls of its rows here.

## The App identity

`layup run` alone holds the App's private key (§3). **Decided here:**

- The key file and the App ID come from the forge register
  ([`records.md`](records.md#nfr-001--the-records-of-start), K40). The run
  checks the file before any call: mode 0600, owned by its user, a PEM block of
  an RSA private key. Any other state is exit 2.
- The adapter makes a JSON Web Token (`JWT`) signed with that key (RS256): issued
  60 seconds in the past, to absorb a clock skew, valid for 9 minutes (GitHub's
  maximum is 10), with the App ID as issuer. It finds the installation of the App
  on `OWNER/NAME` and makes an **installation token** from it, which GitHub keeps
  valid for one hour. The adapter makes a new token when less than five minutes
  of the old one remain. A token is never written to a file, a record or a log.
- `git` gets the token for one call only, from `internal/git`, as an HTTP
  header in `GIT_CONFIG_COUNT`, `GIT_CONFIG_KEY_0` and `GIT_CONFIG_VALUE_0` of
  that call's environment (`http.<web>/.extraHeader`), not on its command line,
  where another user of the host could read it. This adds to the fixed list of
  the environment ([`packages.md`](packages.md#the-calls-of-internalgit)); no
  value of the host is passed.

## The calls of M2a

The GitHub REST API, version header `X-GitHub-Api-Version: 2022-11-28`, at the
`api` of the forge register. The documentation of each endpoint is read at the
date of the build task; a difference from this table is a defect of this table.

| Method | Call | What it gives `layup run` |
| ------ | ---- | ------------------------- |
| `Installation` | `GET /repos/{owner}/{repo}/installation` (JWT) | the installation ID and its permissions |
| `Token` | `POST /app/installations/{id}/access_tokens` (JWT) | the installation token and its end time |
| `Repository` | `GET /repos/{owner}/{repo}` | the default branch, the visibility, whether the repository has a commit |
| `UserID` | `GET /users/{login}` | the numeric ID of a login; the bot's ID for `<slug>[bot]` |
| `OpenIssue` | `POST /repos/{owner}/{repo}/issues` | the issue number |
| `Comments` | `GET /repos/{owner}/{repo}/issues/{n}/comments`, every page | each comment: its ID, the author's ID and login, `performed_via_github_app` (the App's slug, or none), the times, the body |

The read-back of the root commit and the push of the records branch are `git`
calls, not API calls ([`run.md`](run.md#the-steps-of-layup-run---new)).

## Forge errors

**Decided here** (FT1): a forge call that fails, gives a status other than the
one its call expects, or gives a body that the adapter cannot read, is `fail` on
its step, with the call, the status and the first line of the message; never a
pass and never a retry that hides it. A `401` of an installation token makes one
new token and one retry; a second `401` is `fail`. A rate limit (`403` or `429`
with a reset time) waits until the reset time and prints one progress line every
ten seconds; it is not a stall.

## The test of the adapter

An integration test of `internal/forge/github` starts an `httptest` server on
loopback that plays the calls above. It checks the JWT (its header, its claims,
its signature with the public key of a test key), the token request, the paging
of `Comments`, and each error rule above. A server that lacks a permission or a
capability makes `forge` fail and name it. The test needs no secret and no
network.
