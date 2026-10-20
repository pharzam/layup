# T-7r93: the runs of the checks

Task `T-7r93` (#441): an ADR, proposed and then accepted by the Operator, [ADR-0028](../../docs/adr/0028-delegate-two-bounded-landing-decisions.md).

## The acceptance check (condition 1 of the plan review 6292149530)

The command fails while the record or its registration is absent:

```
sh -c 'set -- docs/adr/0028-*.md; test -f "$1" && grep -q "^Proposed$" "$1" && grep -q "^| \[0028\](0028-" docs/adr/README.md'
```

- **Red:** at `2866a3febf99130a074a5608df98b32509a70536` (`main` at the claim), 2026-10-20T13:36:51Z: exit 1.
- **Green:** after S1, on the tree of `main` `1ef5980` (which holds ADR-0027 of #443) with the record and its index row, 2026-10-20T14:34:35Z: exit 0.

## The checks of S2

- `sh docs/adr/adr-lint.sh`: start 2026-10-20T14:35:00Z, exit 0
- `sh docs/prd/prd-lint.sh`: start 2026-10-20T14:35:01Z, exit 0
- `sh docs/links/link-lint.sh`: start 2026-10-20T14:35:01Z, exit 0
- `sh docs/tests/run-discipline-tests.sh`: start 2026-10-20T14:35:13Z, exit 0
- `sh docs/setup/setup-check.sh`: start 2026-10-20T14:35:16Z, exit 0
- `git diff --check`: exit 0.

`adr-lint` warns that nothing outside `docs/adr/` links the record yet. The warning is not fatal, and the task file of the close-out links it.

## The Operator's decision (6296168934), and the record made Accepted

After the merge of `origin/main` (`0f0890b`), the record's Status is `Accepted`, its figures are those of the Operator (20 %, at least 60 and at most 500 lines, plus 4 files), and its index row is `Accepted`. ADR-0026's Status line is left to the follow-up task that changes the operative text, as the plan put it: editing it here would be a sixth file against the plan review's maximum of five.
- `sh docs/adr/adr-lint.sh`: 2026-10-20T16:47:17Z, exit 0
- `sh docs/prd/prd-lint.sh`: 2026-10-20T16:47:18Z, exit 0
- `sh docs/links/link-lint.sh`: 2026-10-20T16:47:18Z, exit 0
- `sh docs/tests/run-discipline-tests.sh`: 2026-10-20T16:47:31Z, exit 0
- `sh docs/setup/setup-check.sh`: 2026-10-20T16:47:34Z, exit 0

## Round 1 (6296550840) and its fixes

- **F1.** The Consequences named "the budget sentence of ADR-0007". An accepted body is immutable, so they now name the Status lines of ADR-0007 and ADR-0026 (`Accepted. Amended by ADR-0028`) with their index rows, and leave the operative text to `engineering-discipline.md` and the glossary.
- **F2.** The pre-freeze check needs the conditions table of this task's plan review: `runs/T-7r93/conditions.md` holds its five conditions.

The acceptance check of condition 1 tested `^Proposed$`. Since the Operator accepted the record, it fails in that form. Its accepted form, `grep -q "^Accepted$"` with the index row ending `| Accepted |`, is the check from here. At 2026-10-20T17:05:50Z on the fix tree, the Proposed form exited `1` and the Accepted form exited `0`.
