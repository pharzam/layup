# The test runs of T-zwke

## Red, at base `add60ef` (2026-10-07T10:31Z)

```
FAIL row 21 is not a row of The tasks of M2a
FAIL row 22 is not a row of The tasks of M2a
FAIL row 23 is not a row of The tasks of M2a
FAIL row 24 is not a row of The tasks of M2a
FAIL row 25 is not a row of The tasks of M2a
FAIL row 26 is not a row of The tasks of M2a
exit 1
```

## D2 red, before the issues (2026-10-07T10:32Z)

```
FAIL runs/T-zwke/issues.tsv is missing: no issue is open
exit 1
```

## Green, on the head before the freeze (2026-10-07T10:34Z)

```
ok   each part is in one row of The tasks of M2a, or named with its reason
ok   row 21: #126 (T-8kqn) by layup-agent[bot]
ok   row 22: #127 (T-esfe) by layup-agent[bot]
ok   row 23: #128 (T-xhgz) by layup-agent[bot]
ok   row 24: #129 (T-6bq5) by layup-agent[bot]
ok   row 25: #130 (T-trej) by layup-agent[bot]
ok   row 26: #131 (T-mqty) by layup-agent[bot]
ok   row 27: #132 (T-fnsr) by layup-agent[bot]
```

Both exit 0.

## Green, on the fix of round 1 (2026-10-07T10:45Z)

```
ok   each part is in one row of The tasks of M2a, or named with its reason
ok   row 21: #126 (T-8kqn) by layup-agent[bot]
ok   row 22: #127 (T-esfe) by layup-agent[bot]
ok   row 23: #128 (T-xhgz) by layup-agent[bot]
ok   row 24: #129 (T-6bq5) by layup-agent[bot]
ok   row 25: #130 (T-trej) by layup-agent[bot]
ok   row 26: #131 (T-mqty) by layup-agent[bot]
ok   row 27: #132 (T-fnsr) by layup-agent[bot]
```

Both exit 0. The bodies of #126, #130, #131 and #132 were edited by the App to match the generator; the body of #132 equals the output of `gen-issues.py` for row 27.

## The fix of O-166 (2026-10-07T11:11Z)

`go run github.com/zricethezav/gitleaks/v8@v8.30.1 git --redact`: 477 commits, `no leaks found`. The same on `add60ef..HEAD`. A first try with `--gitleaks-ignore-path` set to an empty file gave `no leaks found`: the option did not stop the read of the repository's own `.gitleaksignore`, so that run proved nothing. The test that does: in a clone of `b29bfea` with `.gitleaksignore` removed, the run on `add60ef..HEAD` gives `leaks found: 1`, the finding of `592c398` that CI reported. So the line of `.gitleaksignore` is what clears it, and no other finding is hidden.
