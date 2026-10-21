# T-wqrk: test runs

## Preparation cases (condition 3): red on the old preparer, green on the new

`runs/T-wqrk/prep-test.sh` runs the preparer on a round 1 config and on fix-round configs in a disposable clone (a local bare repository stands for the clone URL by git URL rewriting; a stub `gh` answers the forge reads). It runs the production `prep-round.py` as it is, with no test hook but the template path. Times are from `date -u` after each run.

**Red**, the live (old) `prep-round.py` and `round-brief.tmpl`, started 12:04:41Z, ended 12:04:43Z, exit 1:

```
PASS round1: prepares, no fix.diff, earlier-review claim kept
grep: /tmp/TMP/r2/layup/.review-in/fix.diff: No such file or directory
FAIL fix round: fix.diff and the search block
PASS unresolved @ in the text is refused
FAIL round 2 without fix_base is refused
FAIL fix_exclude outside runs/TASK/NAME.md is refused
FAIL an earlier verdict is refused without earlier_record
FAIL an earlier verdict handed over is stated in the brief
```

The fix-round case fails because the old preparer writes no `fix.diff` and no block. The round-2 refusals and the earlier-verdict cases fail because the old preparer has none. The round 1 case and the unresolved-`@` refusal pass on it, so those two behaviours are kept.

**Green**, `prep-round.py.new` with `round-brief.tmpl.new`, ended 12:04:45Z, exit 0:

```
PASS round1: prepares, no fix.diff, earlier-review claim kept
PASS fix round: fix.diff and the search block
PASS unresolved @ in the text is refused
PASS round 2 without fix_base is refused
PASS fix_exclude outside runs/TASK/NAME.md is refused
PASS an earlier verdict is refused without earlier_record
PASS an earlier verdict handed over is stated in the brief
```
