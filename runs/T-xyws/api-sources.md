# T-xyws: the sources of the two close calls

Read on 2026-10-18 (UTC, 17:52 to 17:56) with `curl -sL` of the three pages below; the HTML was stripped to text
and searched with `grep`. Saved sizes and SHA-256 (first 16 hex): "Update an issue" page
`https://docs.github.com/en/rest/issues/issues` 1,146,489 bytes, `8bf860a842d1ea3c`; "Update a pull request" page
`https://docs.github.com/en/rest/pulls/pulls` 884,661 bytes, `6e182b278777e11f`; the permissions page
`https://docs.github.com/en/rest/authentication/permissions-required-for-github-apps` 2,005,553 bytes,
`c452acb6bd36b9ee`. The API version in the pages' samples is `2026-03-10`. Where a page is silent the row says so,
and `forge.md` marks the choice "decided here".

| Call | Source | What it says (used by a fixture) | Not said |
| ---- | ------ | -------------------------------- | -------- |
| `ClosePullRequest` | "Update a pull request" | `PATCH /repos/{owner}/{repo}/pulls/{pull_number}`; optional body `title`, `body`, `state` (`open` or `closed`), `base`, `maintainer_can_modify`; answers `200` (the pull request), `403`, `422` ("Validation failed, or the endpoint has been spammed") | A `404`, a `401`, a `410` (the page lists only `200`, `403`, `422`); what a merged pull request answers to `state: closed`; that closing a closed one is a `200` (the page does not say it is idempotent); that the answer's `state` is `closed` |
| `CloseIssue` | "Update an issue" | `PATCH /repos/{owner}/{repo}/issues/{issue_number}`; body `state` (`open` or `closed`); `state_reason` ("The reason for the state change. Ignored unless state is changed."), one of `completed`, `not_planned`, `duplicate`, `reopened`, `null`; answers `200`, `301`, `403`, `404`, `410`, `422`, `503`; the answer is the issue, with `state`, `state_reason` and, for a pull request, a `pull_request` object | What a second close does with a different `state_reason` (the page says only "ignored unless state is changed", so a repeat close ignores the reason); that `state_reason` is echoed in the answer of every close; `401` and `429` (the shared transport handles those) |
| The number space | the same page | The sample answer of "Update an issue" carries a `pull_request` object with `url`, `html_url`, `diff_url`, `patch_url`; the page does not say that an issue number can name a pull request, but the object shows that the endpoint answers for one | That `PATCH /issues/{n}` of a pull request number closes the pull request (the answer's shape implies it; no sentence says it) |
| A merged pull request | "Update a pull request" | Nothing | Everything: the page says nothing of a merged pull request and `state`. Silence does not establish a `422`. The plan's draft had the `Fake` refuse a merged pull request with a `422`; that proposal was rejected at S1. The policy chosen (`forge.md`, decided here) is a success: the close gives the pull request as the forge answers it, `Merged` kept, and the `Fake` models no `422`. GitHub was not observed on this case |
| The permissions | "Permissions required for GitHub Apps" | `PATCH /repos/{owner}/{repo}/pulls/{pull_number}`: Pull requests `write`. `PATCH /repos/{owner}/{repo}/issues/{issue_number}`: Issues `write`, and also listed under Pull requests `write` (the "Update an issue" page says the token must have at least one of the two) | — |

Facts that changed the specification (R10, in `forge.md`): no new permission is needed (pull requests `write` is in
`LoopPermissions`; issues `write` is in `Permissions` since `M2a`). The `Fake`'s `404` for an unknown number is the `Fake`'s
policy, since the pull request page lists no `404`. A merged pull request is not refused: the close succeeds and keeps
`Merged` (the plan's draft `422` was rejected; the choice is `forge.md`'s, decided here). Whether the real
forge accepts a close of a merged pull request is to be observed in the demo; the adapter does not depend on it,
because it reads the answer, and a `422` is a forge error either way.

Row 64's landed contracts (`docs/tasks/T-cr8j.md`, `docs/spec/forge.md`, "The calls of M2e") were read before code: the
`Forge` interface, `forge.Fake`'s lock and `Fail` order, `pullWire.pull()`, `do`, and the refusal-first pattern of
`Merge`; the four conditions of its plan review concern those calls and do not constrain these.
