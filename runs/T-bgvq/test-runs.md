# T-bgvq — test runs

The two lessons of the review rounds of `T-bng8` (#288). The test is
`runs/T-bgvq/check.sh`: red on `main`, where neither entry exists, green after the
two entries are in `docs/guardrails.md`.

## Red (2026-10-12, 12:23Z)

On `45f546c`, `sh runs/T-bgvq/check.sh`: exit 1, one line for each missing entry. The
first form of the check asked for each entry in four places and printed five lines;
it was fixed before this record, so that presence is decided once per title and the
other rules wait for it. The script sets `LC_ALL=C`, so a host's locale does not
change what it reads.

```
FAIL the first list of section 2 holds 0 entries titled: A review record without its heading.
FAIL the first list of section 2 holds 0 entries titled: An allow rule whose glob reaches a script with a write mode.
```
