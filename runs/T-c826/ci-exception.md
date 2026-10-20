# T-c826 checker-only CI exception

Binding authorization: plan review
[6301028661](https://github.com/pharzam/layup/issues/449#issuecomment-6301028661)
and the supplied Operator decision
[6301084598](https://github.com/pharzam/layup/issues/449#issuecomment-6301084598).
They approve the narrowed R4 route. Removal belongs to
[T-szvp, #452](https://github.com/pharzam/layup/issues/452), which remains open.
The actual migration pull request is [#453](https://github.com/pharzam/layup/pull/453).

The Stage B workflow at `fc3a65a84e6f7b62353a8a1081769dda70266553` keeps
only `docs/tests/freeze-check.sh` when the event is `pull_request` and its number
is the literal 453. Both runners and the other six restore-list scripts always
restore from the default branch in place. The checker also restores for every
other PR and main push. There is no branch-name selector, pin, blob comparison,
fetch by SHA or restore replay. The nested-checkout driver retains its earlier,
separate exception. The workflow file itself still requires review.

## Observed red before implementation

[Run 39497659349, job 118887031319](https://github.com/pharzam/layup/actions/runs/39497659349/job/118887031319)
tested merge `1470b85f4f34cc30c40e804c7ab204e03d56c7ca` from Stage B head
`fc3a65a84e6f7b62353a8a1081769dda70266553`. The job concluded failure, with
`freeze-check: 7 passed, 74 failed` and process exit 1 at
`2026-10-20T21:07:22.6612549Z`. The restore and nested-checkout steps succeeded.
The decoded log and step summaries were read before implementing K5 in Stage C;
[test-runs.md](test-runs.md) gives the exact kept/restored lines and times.

## Local green and remaining hosted evidence

Stage C passes all 81 cases under both `sh` and `bash --posix`, the PR-numbered
task check, and setup-check including `ci`, with both runners unchanged.
[mutants.md](mutants.md) records independent attacks on the guards and restoration
assertions. Local results do not establish a hosted green. After committing and
pushing this stage, the coordinating session must record PR 453's green
`nested-checkout-check` run, its tested SHA and the same kept-checker/restored-runner
lines. The frozen-head review and final required CI checks remain later stages.
