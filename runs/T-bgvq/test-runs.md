# T-bgvq — test runs

The two lessons of the review rounds of `T-bng8` (#267), as task `T-bgvq`
(#288). The test is `runs/T-bgvq/check.sh`: red on `main`, where neither entry exists, green after the
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

## Green (2026-10-12, 12:24Z)

The two entries are the last two of the first list of section 2 of
`docs/guardrails.md`, after the entry of `T-w89c` and before "Writing a lesson back"
(30 lines). `sh runs/T-bgvq/check.sh`: exit 0, "ok   the two lessons of T-bng8 are the
last entries of the first list of section 2, in the form of their neighbours".

## The mutants (2026-10-12, 12:25Z)

On copies of the file, through `GUARDRAILS`, never on the worktree. A first run
built two of them wrongly (a shell read the backquotes of the script as commands, so
one mutant was never made; and the swap joined two entries on one line, so it read
"0 entries"), and was thrown away. They were built again from a script file, each
compared with the real file (1 to 31 changed lines), and these are the runs of those.
Each exits 1:

- **a.** The `**The check:**` of the first entry written without its stars: two lines,
  the missing clause and the order of the parts.

  ```
  FAIL the entry "A review record without its heading." does not hold: **The check:**
  FAIL the entry "A review record without its heading." does not hold its parts in the order: silent, check, Learned in
  ```

- **b.** The second entry moved to the end of section 2, out of the first list: one line.

  ```
  FAIL the first list of section 2 holds 0 entries titled: An allow rule whose glob reaches a script with a write mode.
  ```

- **c.** The second entry without its last line, "Learned in": two lines.

  ```
  FAIL the entry "An allow rule whose glob reaches a script with a write mode." does not hold: Learned in `T-bng8`
  FAIL the entry "An allow rule whose glob reaches a script with a write mode." does not hold its parts in the order: silent, check, Learned in
  ```

- **d.** The two entries swapped: one line.

  ```
  FAIL the two lessons are not the last two entries of the first list of section 2
  ```

- **e.** One line of the first entry not indented: one line.

  ```
  FAIL the entry "A review record without its heading." has 1 lines that are not indented by two spaces
  ```

## The checks (2026-10-12, 12:25Z)

`adr-lint`, `prd-lint`, `link-lint` (2,329 links), `setup-check` (`guardrails`,
`markers` and `adapted` each OK: the entries hold none of the words of rule 3),
`run-discipline-tests` (102 passed) and `sh runs/T-bgvq/check.sh`: exit 0.
`git diff --cached --check`, after `git add` of the new files: exit 0.
