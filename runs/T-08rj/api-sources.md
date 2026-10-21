# T-08rj - the source of the edit of a pull request's body

Read on 2026-10-18 (UTC) by `WebFetch` of the two pages below. The fixtures of
`internal/forge/github/github_integration_test.go` for `EditBody` follow these facts; where a page is
silent the row says so, and `forge.md` marks it. No live-forge probe was made.

| Call | Source | What it says (used by a fixture) | Not said |
| ---- | ------ | -------------------------------- | -------- |
| `EditBody` | [REST: pulls](https://docs.github.com/en/rest/pulls/pulls), "Update a pull request" | `PATCH /repos/{owner}/{repo}/pulls/{pull_number}`; body parameters, all optional: `title`, `body` ("The contents of the pull request."), `state` (`open` or `closed`), `base`, `maintainer_can_modify`; response `200`; `403` forbidden; `422` validation failed or the endpoint spammed; the answer is the pull request schema of "Create a pull request", where `body` is "string or null" | the permission (the page names none for the endpoint); what an omitted field does (only the sent field is changed, as the fixture sends `body` alone); whether a closed or merged pull request can be edited |
| the permission | [Permissions required for GitHub Apps](https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps), "Pull requests" | "Update a pull request": `PATCH /repos/{owner}/{repo}/pulls/{pull_number}`, write access, with user and installation tokens, no additional permission | the API name of the permission (`pull_requests`, as `forge.LoopPermissions` already holds) |

Facts that changed the specification (R10, in `forge.md`): the edit needs pull requests `write`, which
the App holds already (row 64); the answer's `body` may be null, so the adapter gives `""`; the edit of a
closed or merged pull request is not stated by the page, so the stand-in edits it (as the forge is
expected to) and a caller that must refuse it reads the state first.
