# T-6bq5 — the reading of the endpoints

`forge.md` ("The calls of M2a") says that the documentation of each endpoint is
read at the date of the build task. Read on 2026-10-08 (REST API version
`2022-11-28`), with the result for each call.

| Method | Page | Read | Result |
| ------ | ---- | ---- | ------ |
| `Installation` | [Apps](https://docs.github.com/en/rest/apps/apps?apiVersion=2022-11-28), "Get a repository installation for the authenticated app" | `GET`, a JWT; `200`; `id`, `permissions` | as the table |
| `Token` | the same page, "Create an installation access token" | `POST`, a JWT; `201`; `token`, `expires_at` | as the table |
| `Repository` | [Repositories](https://docs.github.com/en/rest/repos/repos?apiVersion=2022-11-28), "Get a repository"; [Branches](https://docs.github.com/en/rest/branches/branches?apiVersion=2022-11-28), "List branches" | `200`; `default_branch`, `visibility`; the branches: `200`, `per_page` up to 100 | as the table; the page does not say what an empty repository gives, so the step `forge` of row 25 shows it |
| `UserID` | [Users](https://docs.github.com/en/rest/users/users?apiVersion=2022-11-28), "Get a user" | `200`; `id`; an unknown login `404` | as the table; the page does not name a `<slug>[bot]` login, so the uat of row 27 shows it |
| `OpenIssue` | [Issues](https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28), "Create an issue" | `201`; `title` required; `number`; `410` when issues are off | as the table |
| `Comments` | [Issue comments](https://docs.github.com/en/rest/issues/comments?apiVersion=2022-11-28), "List issue comments" | `200`; `per_page` up to 100; `user` may be `null`; `performed_via_github_app.slug`; `created_at`, `updated_at` | a gap: `user` may be `null`; `forge.md` now gives the rule |
| (each) | [Rate limits for the REST API](https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2022-11-28) | `403` or `429`; `retry-after`, else `x-ratelimit-reset` when `x-ratelimit-remaining` is `0` | `forge.md` "Forge errors" now gives this reading (condition 1 of the plan review) |
