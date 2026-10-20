# T-7r93: the runs of the checks

Task `T-7r93` (#441): a Proposed ADR, [ADR-0028](../../docs/adr/0028-delegate-two-bounded-landing-decisions.md).

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
