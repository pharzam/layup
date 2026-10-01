# Test traceability

LAYUP's test-side view of the traceability line, in the form of
[`traceability-template.md`](traceability-template.md): one row per test, so a
reader goes from any requirement of [`PRD-0001`](../prd/PRD-0001-layup.md) to the
test that covers it. The [DoD checklist](dod-checklist.md) asks that the ID set of
this table equals the requirement set: each `REQ` and `NFR` has at least one row.

Task `T-55n2` ([#76](https://github.com/pharzam/layup/issues/76)) wrote this table
from the [implementation plan](../plan/README.md). Each row carries its status at
`7cdd346`: `green` for a test that exists and passes, `planned` for a test that a
task of the plan will write. A planned row names its test by a handle of the form
`<task>/<level>/<name>`; the delivering task puts the real test name in its place,
with its status, in its own pull request. A row of a later milestone names the
milestone, whose specification task names the test.

A discipline check (for example check `pin` of
[`setup-check.sh`](../setup/setup-check.sh) for LAYUP's own pin) is not a row here,
because a row's level is one of `unit`, `integration`, `e2e` and `uat`;
`PRD-0001` §12 names such a check in its Test column. The tests of `layup version`
and of the usage errors (`TestBinaryPrintsItsVersion`,
`TestVersionPrintsTheVersionAndExitsZero`, `TestNoSubcommandPrintsUsageAndExitsTwo`,
`TestUnknownSubcommandExitsTwo`) cover no requirement of `PRD-0001`: they test the
command frame of [`spec/README.md`](../spec/README.md), and row 3 of the plan
(`T-2yw7`) gives them rows. A row that covers more than one requirement gives the
fact and the ADR of the first.

| Test ID | Level | Covers (REQ/NFR) | Fact (F-NNNN#n) | Guardrail | ADR | Task | Status |
|---------|-------|------------------|-----------------|-----------|-----|------|--------|
| `TestGolden` (`internal/psb/check_test.go`) | unit | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05 | green |
| `TestPSBCheckExitCodes` (`internal/cli/cli_test.go`) | unit | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05 | green |
| `TestGoldenRealPSB` (`internal/psb/check_integration_test.go`) | integration | REQ-001 | F-0003#41 | — | ADR-0011 | T-dq05 | green |
| `T-5zmw/e2e/psb-check` | e2e | REQ-001 | F-0003#41 | — | ADR-0011 | T-5zmw | planned |
| `T-18v6/integration/schema-blocks` | integration | NFR-001 | F-0001#1 | guardrails.md §1.1 Inv-1 | ADR-0014 | T-18v6 | planned |
| `TestPackageRules` (`cmd/layup/rules_integration_test.go`) | integration | NFR-005, NFR-007 | F-0001#6 | guardrails.md §1.1 Inv-6 | ADR-0015 | T-2tc2 | green |
| `T-5sgt/e2e/gate-command` | e2e | REQ-004, REQ-007, NFR-004, NFR-005 | F-0003#44 | guardrails.md §1.1 Inv-5 | ADR-0016 | T-5sgt | planned |
| `T-c06a/integration/go-fixtures` | integration | REQ-004, REQ-007 | F-0003#44 | guardrails.md §1.1 Inv-7 | ADR-0016 | T-c06a | planned |
| `T-8vpw/e2e/verify-not-active` | e2e | NFR-004, REQ-002 | F-0001#5 | guardrails.md §1.1 Inv-5 | ADR-0011 | T-8vpw | planned |
| `T-7s0y/integration/target-pin` | integration | NFR-006, REQ-002 | F-0001#8 | guardrails.md §1.1 Inv-8 | ADR-0009 | T-7s0y | planned |
| `T-dep6/e2e/whole-setup` | e2e | REQ-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-4 | ADR-0016 | T-dep6 | planned |
| `T-dep6/e2e/records-in-git` | e2e | NFR-001 | F-0001#1 | guardrails.md §1.1 Inv-1 | ADR-0014 | T-dep6 | planned |
| `T-dep6/e2e/gate-without-layup` | e2e | NFR-002 | F-0001#2 | guardrails.md §1.1 Inv-2 | ADR-0013 | T-dep6 | planned |
| `T-tmhw/integration/telemetry-schema` | integration | REQ-011 | F-0003#50 | — | ADR-0024 | T-tmhw | planned |
| `T-dgy7/integration/stalls-schema` | integration | REQ-009 | F-0003#49 | — | ADR-0023 | T-dgy7 | planned |
| `T-efmy/uat/release-review` | uat | REQ-015, REQ-017 | F-0003#53 | — | — | T-efmy | planned |
| `T-evad/uat/first-pilot` | uat | REQ-001, REQ-002, REQ-004, REQ-007, REQ-016, REQ-018, NFR-001, NFR-002, NFR-003 | F-0003#42 | guardrails.md §1.1 Inv-7 | ADR-0012 | T-evad | planned |
| `M2f/uat/rule-protection` | uat | REQ-003 | F-0003#43 | guardrails.md §1.1 Inv-3 | ADR-0017 | M2f | planned |
| `M2e/uat/role-handoffs` | uat | REQ-005 | F-0003#45 | — | ADR-0019 | M2e | planned |
| `M3c/uat/clarification` | uat | REQ-006 | F-0003#46 | — | ADR-0019 | M3c | planned |
| `M3b/e2e/escalation` | e2e | REQ-008 | F-0003#48 | — | ADR-0022 | M3b | planned |
| `M3d/e2e/stall-procedure` | e2e | REQ-010 | F-0003#49 | — | ADR-0023 | M3d | planned |
| `M2c/uat/specification-trace` | uat | REQ-012 | F-0003#51 | — | ADR-0018 | M2c | planned |
| `M4a/uat/two-harnesses` | uat | REQ-013 | F-0003#52 | guardrails.md §1.1 Inv-9 | ADR-0020 | M4a | planned |
| `M4c/uat/two-stacks` | uat | REQ-014 | F-0003#67 | — | — | M4c | planned |
