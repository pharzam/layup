# T-ax3r — the test runs

Evidence for row 25b (#142): each test red first, for the right reason, then green.

## `forge.Repository` with the branch names, `git.Clone` with `Auth` (2026-10-08T10:40Z)

Red: the adapter's test asks `Branches` (one page, then 101 branches over two pages, then none), and `TestEachCallRunsItsVerb` gets the row "clone with a token" (the three `GIT_CONFIG_*` values): neither compiled (`unknown field Branches`, `too many arguments in call to Clone`). Then `Repository` reads every page of the branches (at the register's `api` only), `Clone` takes `Auth` through `doAuth`, and `internal/setup` passes `git.Auth{}`: `go test` and `go test -tags=integration` of `internal/git`, `internal/forge/...` and `internal/setup` pass.
