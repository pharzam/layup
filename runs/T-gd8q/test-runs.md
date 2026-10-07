# The test runs of T-gd8q

Each run of [`checks.sh`](checks.sh), from the root of the worktree.

## Red, at base `34dd858` with only `checks.sh` added (2026-10-07T05:35Z)

```
FAIL C1 ADR-0026 is missing or not Accepted
FAIL C2 the Status of ADR-0012 does not name ADR-0026
FAIL C3 docs/engineering-discipline.md still says 'One pass is never enough'
FAIL C4 runs/T-gd8q/mentions.tsv is missing
ok   C5 docs/adr/adr-lint.sh exits 0
ok   C5 docs/links/link-lint.sh exits 0
ok   C5 docs/setup/setup-check.sh exits 0
exit 1
```

## A failing run before the freeze

With ADR-0026 staged, C4 listed its own lines (the file name holds the word) and C5 failed on check `adapted` (the words "kit" and "Armature" in one line of the ADR). Fixed: `mentions.tsv` classes ADR-0026 as the new rule, and the line was reworded.

The mutation of C4: a line `bootstrap test line` added at the end of `docs/glossary.md` gave `FAIL C4 mentions with no class in runs/T-gd8q/mentions.tsv: docs/glossary.md:172`. Its undo by `git checkout` also lost the unstaged glossary edits, which were made again (the new pitfall of `docs/guardrails.md` §2).

## Green, on the head before the freeze (2026-10-07T05:48Z)

```
ok   C1 ADR-0026 is Accepted (docs/adr/0026-keep-the-bootstrap-review-rules-as-the-standing-gate.md)
ok   C2 the Status of ADR-0012 is Superseded by ADR-0026
ok   C3 'One pass is never enough' is gone
ok   C4 each live mention of 'bootstrap' is classed in runs/T-gd8q/mentions.tsv
ok   C5 docs/adr/adr-lint.sh exits 0
ok   C5 docs/links/link-lint.sh exits 0
ok   C5 docs/setup/setup-check.sh exits 0
exit 0
```
