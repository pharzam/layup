# T-cr8j — the sources of the forge calls of M2e

Read on 2026-10-15 (UTC, 18:40 to 18:50), by `WebFetch` of the pages below, and, for the
GraphQL mutation, of GitHub's public schema file `https://docs.github.com/public/fpt/schema.docs.graphql`
(1,562,177 bytes, SHA-256 `3ed095654b22b1a8…`, saved locally and searched with `grep`). The test
fixtures of `internal/forge/github/github_integration_test.go` follow these facts; where a page
is silent the row says so, and `forge.md` marks it. The pull request fixtures are two (round 1 of `T-cr8j`, finding 1): `prJSON`, a pull request as `GET /pulls/{n}` gives it (`merged`, and a test merge commit in `merge_commit_sha` before a merge), and `listItemJSON`, a list item of `GET /pulls` (no `merged`, a `merged_at` that is null until the merge); the check runs and the statuses follow the two list pages; the GraphQL answer follows the schema's payload type.

| Call | Source | What it says (used by a fixture) | Not said |
| ---- | ------ | -------------------------------- | -------- |
| `SetStatus` | [REST: commit statuses](https://docs.github.com/en/rest/commits/statuses), "Create a commit status" | `POST /repos/{owner}/{repo}/statuses/{sha}`; body `state` (required: `error`, `failure`, `pending`, `success`), optional `target_url`, `description`, `context`; response `201`; 1,000 statuses per SHA and context; `creator` is null or a user | the length limit of `description` (the 140 characters are LAYUP's own bound, stricter than the page) |
| `Statuses` (statuses) | the same page, "List commit statuses for a reference" | `GET /repos/{owner}/{repo}/commits/{ref}/statuses`; newest first; `per_page` up to 100, default 30; each item has `creator`, null or a user | — |
| `Statuses` (check runs) | [REST: check runs](https://docs.github.com/en/rest/checks/runs), "List check runs for a Git reference" | `GET /repos/{owner}/{repo}/commits/{ref}/check-runs`; `filter` `latest` (default) or `all`; `per_page` up to 100; answer `{total_count, check_runs[]}`, each with `name`, `status`, `conclusion`, `app` (an object with `id` and `slug`, or null) | the permission (the page names none; the permissions page lists "Checks" read) |
| `OpenDraft` | [REST: pulls](https://docs.github.com/en/rest/pulls/pulls), "Create a pull request" | `POST /repos/{owner}/{repo}/pulls`; `draft` boolean; response `201` with the pull request schema; drafts exist in public repositories on GitHub Free, and in public and private repositories on Team and Enterprise Cloud | — |
| `PullRequest` (list) | the same page, "List pull requests" | `GET /repos/{owner}/{repo}/pulls`; `head` is `user:ref-name`; `state` `open`, `closed` or `all` (default `open`); `sort` default `created`, `direction` default `desc`: the first item is the newest | merged state of a list item (a list item has `merged_at`; the adapter reads `merged_at` too) |
| `PullRequest` (one) | the same page, "Get a pull request" | `number`, `node_id`, `state` (`open` or `closed`), `draft`, `merged`, `head.sha`, `merge_commit_sha` (before a merge a test merge commit, after it depending on the method) | — |
| `Merge` | the same page, "Merge a pull request" | `PUT .../pulls/{n}/merge`; body `sha` (the head SHA that must match), `merge_method` (`merge`, `squash`, `rebase`); `200` with `sha`, `merged`, `message`; `405` the merge cannot be performed; `409` a `sha` was given and does not match the head | — |
| `MarkReady` | the GraphQL schema file, `type Mutation` and its types | `markPullRequestReadyForReview(input: MarkPullRequestReadyForReviewInput!): MarkPullRequestReadyForReviewPayload`, "Marks a pull request ready for review"; the input has `pullRequestId: ID!` ("ID of the pull request to be marked as ready for review") and `clientMutationId`; the payload has `pullRequest: PullRequest` | the permission of the mutation (no page names it; taken as pull requests write, to confirm on the real forge in the demo) |
| the permissions | [Permissions required for GitHub Apps](https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps) | "Create a pull request": Pull requests write; **"Merge a pull request": listed under Contents write (not Pull requests)**; commit statuses write and read as above; check runs list: Checks read; the API names (`pull_requests`, `statuses`, `checks`) are not shown on the page | the API names of the permissions |

Facts that changed the specification (R10, in `forge.md`): the merge call needs Contents write, which the
App holds already (M2a); the description's 140 characters are LAYUP's own bound; the API names of the
three permissions are the ones that `GET /repos/{owner}/{repo}/installation` returns (`pull_requests`,
`statuses`, `checks`), to confirm on the real forge in the demo.
