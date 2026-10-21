# T-3xvc: the compiled mutants

Code head `c1aa9505` (round 2 of #368: C1 expects exactly zero merge requests; the product code is that of `357a757c`). The records against `21a3ae79` and `dd02f40a` are superseded: C1 was killed by the row mismatch, then by a count that always expected two (review finding of round 2). The tree was copied (`tar` without `.git` and `runs/`) to `/home/layup/.cache/layup-session/t3/scratch`, outside the worktree, by `/home/layup/.cache/layup-session/t3/mutants.sh`; each mutant is one replacement of one text, built (`go build ./...`), vetted (`go vet -tags=e2e ./cmd/layup/ ./internal/run/`), run with `go test -count=1 -tags=e2e -run 'TestTheSpecificationTaskEndToEnd' ./cmd/layup/` (the test builds the binary from the mutated tree, so the binary carries the mutant), restored by copying the original file back, and its hash read again. The complete output is `/home/layup/.cache/layup-session/t3/mutants-out.txt` (rerun for round 1).

**Result: 5 mutants built and vetted (E1, E2, E3, C1a, C1b); 5 killed; every restore matched its hash: True; the control before (`ok 21.344s`) and after (`ok 21.422s`) was green: True.**

| # | Mutant (the one change) | File | Killed by | First failing assertion | sha256 before | during | after | Restore ok |
| - | ----------------------- | ---- | --------- | ----------------------- | ------------- | ------ | ----- | ---------- |
| E1 | `afterScaffold` drops `r.loopStep` | `internal/run/steps.go` | the row sequence of stage 6 | `loop_e2e_test.go:464: stage 6: exit 0, the rows [... "shape" "bet"]; want ... "shape" "bet" "loop"` | b3ee6748df8dab62 | f78919d00ca4ff16 | b3ee6748df8dab62 | True |
| E2 | `case adv.Merged:` returns `done` after the first merge | `internal/run/loop.go` | the merge count (1, want 2), then the row `loop` of stage 6 | `loop_e2e_test.go:444: stage 6: 1 merge requests in the log; want 2, one for each task of Shape` (then `:447`, the rows map[... "bet":["done" "the bet 1: the rendered commit 5d23e2cbec63 ...` (the row `loop` is not `the loop merged <S>, <A>`) | a9d9aaadd137843d | b1267380ec56abad | a9d9aaadd137843d | True |
| E3 | the SHA that `MergeTask` passes to the forge is `0000...` (40 zeros) in place of `head`, **after** the head comparison (condition 5: a mutation of `ready.go:213`, not of the head the statuses were read at) | `internal/run/ready.go` | the fake's 409 and the row `loop` | `loop_e2e_test.go:464: stage 6: exit 1 ... loop fail the merge of the pull request 1 at 983dbaa0... was refused and its head did not move: Merge: status 409` | 17e36140730b0505 | 32f274eaf008ea34 | 17e36140730b0505 | True |
| C1a | the fake answers the required job `static` with `conclusion: failure` (`const failedJob = "static"`; fixture control, NFR-004): the scenario then expects exactly zero merge requests | `cmd/layup/loop_e2e_test.go` | the zero assertion **passes** (0 requests, executed before the endpoint comparison); the happy-path row comparison fails, as the scenario is not the happy path | `loop_e2e_test.go:473: stage 6: the rows map[... "bet":["done" "the bet 1: the rendered commit 6af75812c110 ...` | 1fbf7ff9330cc7f3 | 04109b5534f95090 | 1fbf7ff9330cc7f3 | True |
| C1b | C1a with a product that ignores the conclusion of a required job (`internal/run/ready.go`: `if run.Status != "completed" {`) | `internal/run/ready.go` (test file as in C1a) | **the zero assertion** | `loop_e2e_test.go:470: stage 6: 2 merge requests in the log; want exactly 0` | 17e36140730b0505 | 0f171ec6155a9b95 | 17e36140730b0505 (test file back to 1fbf7ff9330cc7f3) | True |

The patches (unified diff, from `mutants-out.txt`):

```diff
E1  steps.go:  -	return []func(context.Context) Step{r.atBoundary("shape", r.shapeStep), r.atBoundary("bet", r.betStep), r.loopStep}
               +	return []func(context.Context) Step{r.atBoundary("shape", r.shapeStep), r.atBoundary("bet", r.betStep)}
E2  loop.go:   -			merged, stuck = append(merged, t.ID), 0
               -			continue
               +			merged, stuck = append(merged, t.ID), 0
               +			return done(name, "the loop merged "+strings.Join(merged, ", "))
E3  ready.go:  -	merge, err := s.Forge.Merge(ctx, pr.Number, head)
               +	merge, err := s.Forge.Merge(ctx, pr.Number, strings.Repeat("0", 40))
C1a loop_e2e_test.go: const failedJob = ""  ->  const failedJob = "static"
C1b ready.go:  -		if run.Status != "completed" || run.Conclusion != "success" {  (with C1a)
               +		if run.Status != "completed" {
```

The original E3 of the plan (the verified head in place of the head of the statuses) was rewritten as condition 5 asks: passing the verified head to `MergeTask` stops at the head comparison of `markReady` before any merge request, so it cannot reach the fake's 409.

The negative check of the zero assertion is the probe `TestTheCountOfMergeRequestsIsExact` (red first against a stub that accepted every count: `red-count.txt`; green with `mergeCountFinding`): zero is accepted and one or two refused when zero is expected, two accepted and one or zero refused otherwise. The line `restored test file ... (before ...)` of `mutants-out.txt` prints the hash of C1b's own step as its `before`; the test file's hash is 1fbf7ff9330cc7f3 before and after.
