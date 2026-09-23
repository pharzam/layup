# Test levels

The fixed ladder of test kinds this project uses, from the cheapest and most
local to the most expensive and most whole-system. It is the reference the rest
of the [test section](README.md) points at: every template and checklist here
names one of these levels. The level definitions are the kit's; the commands
are this project's Go values, set by the Operator (decision O-5 on
[#8](https://github.com/pharzam/layup/issues/8)) and recorded with their
evidence in [`setup/record-T-n1hp.md`](../setup/record-T-n1hp.md).

## In plain terms

> A test is only useful if you know what it proves and how much it costs to run.
> This project sorts its tests into four kinds and runs the cheap ones on every
> commit and the expensive ones in the shared pipeline, so a broken change is
> caught in seconds locally and confirmed thoroughly before it merges.

## The ladder — cheap first

The levels are ordered so a failure stops the expensive work early. A change runs
the cheap levels in the [commit hook](../engineering-discipline.md#git-hooks) for
fast local feedback, and the whole ladder in
[CI](../engineering-discipline.md#continuous-integration-optional) as the
authority.

| Level | Proves | Scope | Speed | Runs in |
|-------|--------|-------|-------|---------|
| 1. Unit | one component behaves | one function/class/module, dependencies stubbed | fastest | hook + CI |
| 2. Integration | components work together | two or more units across a real seam | medium | CI |
| 3. End-to-end (E2E) | a whole user path works | the running system, front to back | slowest | CI |
| Discipline | the process stays honest | repo files, no product toolchain | fast | hook + CI |

The three numbered rungs — unit, integration, end-to-end — are the **product-test
levels** you tag by level and run cheap-first. **Discipline** tests are a separate,
process-level track (they lint the repo's own conventions); the table lists them for
the full picture, but they are not one of the tagged product levels.

## 1. Unit tests

A unit test exercises **one component in isolation** — the smallest piece of
behaviour that stands on its own — with its dependencies replaced by stand-ins.
When it fails, the fault is in that one component, not somewhere across a chain.
Unit tests touch no file, network, or clock, so they are fast and deterministic
and run first, on every commit.

- **Command:** `go test ./...`
- **Where:** the commit hook and CI.
- **Rule:** every component has at least one unit test (see
  [`template-unit.md`](template-unit.md)).

## 2. Integration tests

An integration test proves that **two or more components work together across a
real interface or workflow** — the seams a unit test stubs out. It uses the real
collaborator (a real datastore, a real adapter) rather than a stand-in, so it is
slower than a unit test and runs after it.

- **Command:** `go test -tags=integration ./...`
- **Where:** CI only (job `tests`).
- **Rule:** every interface or workflow has an integration test (see
  [`template-integration.md`](template-integration.md)).

## 3. End-to-end (E2E) tests

An E2E test walks a **whole user-facing path through the running system**, front
to back, the way a real user or caller would. It is the most expensive automated
level, so it usually runs in CI rather than the commit hook; a tiny smoke subset
may run locally to prove the wiring.

- **Command:** `go test -tags=e2e ./...`
- **Timeout:** `-timeout 10m` — an E2E test that hangs must fail, not stall the
  pipeline.
- **Where:** CI only (job `tests`).
- **Rule:** every user-facing scenario has an E2E test (see
  [`template-e2e.md`](template-e2e.md)).

**User acceptance (UAT) is a human layer on top of E2E.** A
[UAT](template-uat.md) scenario checks the *same* path a person cares about, but
a human runs or signs it off against plain Given/When/Then steps. It is judged by
a person, not asserted by a command, so it is not a rung of the automated ladder
— it is the acceptance step that rides on the E2E path.

## Discipline tests

A discipline test lints the **process rather than the product**: it checks the
repo's own conventions and needs no product toolchain, so it can be the project's
first test, before any product code exists. The kit ships five:
[`adr-lint.sh`](../adr/adr-lint.sh), [`prd-lint.sh`](../prd/prd-lint.sh) and
[`link-lint.sh`](../links/link-lint.sh)
read repo files and run in
both the hook and CI, while
[`pr-link-lint.sh`](../ci/pr-link-lint.sh) and
[`review-record-lint.sh`](../ci/review-record-lint.sh) read a forge
artifact absent at commit time — the pull-request body and the linked issue — so
they run in CI only. Add one whenever a
convention is worth enforcing by machine rather than by review.

These linters are themselves tested. [`run-discipline-tests.sh`](run-discipline-tests.sh)
runs each one — the five above **and** the [`commit-msg`](../../.githooks/commit-msg)
hook — against a fixture suite, asserting that a `good` fixture is accepted and a
`bad-*` fixture rejected. It is a *harness over* the linters, not a sixth linter,
and it runs in the hook and CI. See
[The discipline self-tests](README.md#the-discipline-self-tests).

## Security tests sit alongside the ladder

Security checks are not a fourth rung but a parallel track — a secret scan, a
dependency scan, and static analysis at minimum. In this project the hook runs
the static analysis (`go vet`); CI runs all three. They
have their own command placeholder and their own checklist:

- **Command:** `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./... && go vet ./... && gitleaks git --redact` runs the scans in [CI](../ci/) (the
  [hook](../../.githooks/pre-commit) runs `go vet` only) — and
  `govulncheck v1.8.0, go vet, and gitleaks` names the tool it drives.
- **Checklist:** [`security-checklist.md`](security-checklist.md).

## The placeholders this section uses

Fill these once, in your own copy, and every template here inherits them:

| Placeholder | Meaning |
|-------------|---------|
| `go test ./...` | Run the unit level. |
| `go test -tags=integration ./...` | Run the integration level. |
| `go test -tags=e2e ./...` | Run the E2E level. |
| `go run golang.org/x/vuln/cmd/govulncheck@v1.8.0 ./... && go vet ./... && gitleaks git --redact` | Run the security scan step. |
| `-timeout 10m` | The per-test (or per-suite) time limit before a hang is a failure. |
| `*_test.go beside the code; root tests/ for end-to-end fixtures` | Where product tests live — the root [`tests/`](../../tests/) drop-in, or your stack's convention. |
| `govulncheck v1.8.0, go vet, and gitleaks` | The tool that runs the security checks (secret scan, dependency scan, static analysis). |
