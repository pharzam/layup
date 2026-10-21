# T-wqrk: test runs

## Chronology, stated plainly

The first draft of `prep-round.py.new` is dated 12:00:52Z, and the first red below (12:04:41Z, the first version of the tests) came after it: **the tests of the first head came after the first draft of the code, so that red is late evidence** and is kept only as history. After round 1 of #462 (five material findings) the tests were rewritten and run **before** the code was changed: the second section is that red, on the first head's preparer. The third section is the red on a clean copy of the OLD preparer with the final tests, and the green. Times are from `date -u` after each run.

## 1. First version of the tests (late evidence)

Old preparer: exit 1, 5 of 7 cases failed. New preparer (first draft): exit 0, 7 of 7 passed. Kept as it stood at commit 5b99be01.

## 2. The final tests, red on the first head's preparer (before the fixes of round 1)

`runs/T-wqrk/prep-test.sh` with the preparer and template of head 4323dd24 (the draft was copied before the edit), at 12:32:26Z:

```
PASS round1: prepares, no fix.diff, earlier-review claim kept
PASS fix round: fix.diff, the search block, nothing left out
PASS an unresolved @ in the text is refused
PASS round 2 without fix_base is refused
PASS a malformed fix_base is refused
FAIL a fix_base equal to the head is refused
FAIL a fix_base that is not an ancestor is refused
PASS an earlier verdict in extra is refused, with a valid fix_base
FAIL the earlier_record opt-in no longer lets a verdict in
FAIL a named passive record leaves the diff and the search, and is named
FAIL a passive path outside runs/TASK/NAME.md is refused
FAIL an executable record cannot be passive
FAIL a record that a script reads cannot be passive
```

These are the cases of the findings: M2 (the `earlier_record` opt-in still lets a verdict in), M3 (a fix_base equal to the head and one that is not an ancestor were not tested, so the guard could not fail a case), M4 (the named passive records, the executable one and the consumed one). The cases that passed already (the earlier-verdict case without the opt-in) passed for the reason round 1 gave: the guard of a missing `fix_base` came first in the old case; the new case has a valid `fix_base` and asserts the message.

## 3. The final tests: red on the OLD preparer, green on the new

Offline: a local bare repository stands for the clone URL by git URL rewriting, and a stub `gh` answers the forge reads. Every refusal case is an otherwise valid config and asserts its message.

**Red**, the live (old) `prep-round.py` and `round-brief.tmpl` (a clean copy, sha256 875cd66f8acd…), ended 12:33:30Z, exit 1:

```
PASS round1: prepares, no fix.diff, earlier-review claim kept
FAIL fix round: fix.diff, the search block, nothing left out
PASS an unresolved @ in the text is refused
FAIL round 2 without fix_base is refused
FAIL a malformed fix_base is refused
FAIL a fix_base equal to the head is refused
FAIL a fix_base that is not an ancestor is refused
FAIL an earlier verdict in extra is refused, with a valid fix_base
FAIL the earlier_record opt-in no longer lets a verdict in
FAIL a named passive record leaves the diff and the search, and is named
FAIL a passive path outside runs/TASK/NAME.md is refused
FAIL an executable record cannot be passive
FAIL a record that a script reads cannot be passive
```

The round 1 case and the unresolved-`@` refusal pass on the old preparer, so those behaviours are kept.

**Green**, `prep-round.py.new` with `round-brief.tmpl.new`, ended 12:33:32Z, exit 0:

```
PASS round1: prepares, no fix.diff, earlier-review claim kept
PASS fix round: fix.diff, the search block, nothing left out
PASS an unresolved @ in the text is refused
PASS round 2 without fix_base is refused
PASS a malformed fix_base is refused
PASS a fix_base equal to the head is refused
PASS a fix_base that is not an ancestor is refused
PASS an earlier verdict in extra is refused, with a valid fix_base
PASS the earlier_record opt-in no longer lets a verdict in
PASS a named passive record leaves the diff and the search, and is named
PASS a passive path outside runs/TASK/NAME.md is refused
PASS an executable record cannot be passive
PASS a record that a script reads cannot be passive
```

## 4. Each guard removed in a scratch copy (finding M3)

`mutants-T-wqrk.py` replaces the line of one guard by `pass` in a scratch copy of `prep-round.py.new` and runs the suite: the guard's own named case must fail. Ended 12:33:46Z, exit 0.

```
earlier    exit 1  own case fails: yes  all FAIL: an earlier verdict in extra is refused, with a valid fix_base; the earlier_record opt-in no longer lets a verdict in
missing    exit 1  own case fails: yes  all FAIL: round 2 without fix_base is refused
malformed  exit 1  own case fails: yes  all FAIL: a malformed fix_base is refused
shape      exit 1  own case fails: yes  all FAIL: a passive path outside runs/TASK/NAME.md is refused
self       exit 1  own case fails: yes  all FAIL: a fix_base equal to the head is refused
ancestor   exit 1  own case fails: yes  all FAIL: a fix_base that is not an ancestor is refused
exec       exit 1  own case fails: yes  all FAIL: an executable record cannot be passive
consumed   exit 1  own case fails: yes  all FAIL: a record that a script reads cannot be passive
```

Every guard has its own case, and removing it fails that case (the earlier-verdict guard also fails the opt-in case, which tests the same guard).
