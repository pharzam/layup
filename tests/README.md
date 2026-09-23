# tests/

The home for this project's cross-package **end-to-end test fixtures**. LAYUP is
written in Go, so unit and integration tests are `*_test.go` files beside the code
they test, and only end-to-end fixtures that span packages live here (Operator
decision O-5 on [#8](https://github.com/pharzam/layup/issues/8)). It is empty
until the first Go code exists.

## In plain terms

> Unit and integration tests sit beside the Go code they test. This folder is for
> end-to-end fixtures that span packages; it is empty because no Go code exists
> yet. The rules for what to write and how live in `docs/tests/`.

## Where the conventions live

This directory holds end-to-end fixtures; the **conventions** for writing them live in
[`docs/tests/`](../docs/tests/):

- [`docs/tests/test-levels.md`](../docs/tests/test-levels.md) — the test levels and
  their Go commands.
- [`docs/tests/template-unit.md`](../docs/tests/template-unit.md),
  [`template-integration.md`](../docs/tests/template-integration.md),
  [`template-e2e.md`](../docs/tests/template-e2e.md) — a pattern to copy per level.
- [`docs/tests/traceability-template.md`](../docs/tests/traceability-template.md) —
  the row that ties each test back to the requirement it proves.

The kit's own [discipline tests](../docs/tests/test-levels.md) — the ADR, PRD,
link, PR-link and review-record linters — are not product tests and do
**not** live here; they stay beside the conventions they enforce, under
[`docs/`](../docs/).
